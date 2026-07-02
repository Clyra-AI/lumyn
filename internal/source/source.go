package source

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Clyra-AI/lumyn/internal/config"
)

const (
	schemaVersion = "1.0"

	intakeReportPath = "runs/source-checks/source-intake.json"
	checkReportPath  = "runs/source-checks/source-check.json"
)

type InitOptions struct {
	ConfigPath  string
	OpenAPIPath string
	DocsPath    string
}

type ProjectConfig struct {
	Version int          `json:"version"`
	Sources SourceConfig `json:"sources"`
	Env     EnvConfig    `json:"env"`
	Redact  RedactConfig `json:"redaction"`
	Verify  VerifyConfig `json:"verify"`
}

type SourceConfig struct {
	OpenAPI []SourceEntry `json:"openapi"`
	Docs    []SourceEntry `json:"docs"`
}

type SourceEntry struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

type EnvConfig struct {
	BaseURL string `json:"base_url"`
	APIKey  string `json:"api_key"`
}

type RedactConfig struct {
	Headers []string `json:"headers"`
	Fields  []string `json:"fields"`
}

type VerifyConfig struct {
	DefaultStrategy string       `json:"default_strategy"`
	Replay          ReplayConfig `json:"replay"`
}

type ReplayConfig struct {
	AllowNetwork bool `json:"allow_network"`
}

type Report struct {
	ObjectType         string      `json:"object_type"`
	SchemaVersion      string      `json:"schema_version"`
	Status             string      `json:"status"`
	ConfigPath         string      `json:"config_path"`
	ReportPath         string      `json:"report_path"`
	SourceRefs         []SourceRef `json:"source_refs"`
	Findings           []Finding   `json:"findings"`
	SurfaceFingerprint string      `json:"surface_fingerprint"`
}

type SourceRef struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
	Path string `json:"path"`
	Hash string `json:"hash,omitempty"`
}

type Finding struct {
	Kind              string    `json:"kind"`
	Severity          string    `json:"severity"`
	Message           string    `json:"message"`
	Reference         Reference `json:"reference"`
	FixTarget         string    `json:"fix_target"`
	WorkflowRelevance string    `json:"workflow_relevance"`
}

type Reference struct {
	Path        string `json:"path"`
	JSONPointer string `json:"json_pointer,omitempty"`
	Line        int    `json:"line,omitempty"`
	Object      string `json:"object,omitempty"`
}

func InitProject(options InitOptions) (Report, error) {
	options = normalizedInitOptions(options)
	root := projectRoot(options.ConfigPath)
	configPathWasRelative := !filepath.IsAbs(options.ConfigPath)
	projectConfig := ProjectConfig{
		Version: 1,
		Sources: SourceConfig{
			OpenAPI: []SourceEntry{{ID: "public_api", Path: sourcePathForConfigRoot(root, options.OpenAPIPath, configPathWasRelative)}},
			Docs:    []SourceEntry{{ID: "docs", Path: sourcePathForConfigRoot(root, options.DocsPath, configPathWasRelative)}},
		},
		Env: EnvConfig{
			BaseURL: "${API_BASE_URL}",
			APIKey:  "${LUMYN_API_KEY}",
		},
		Redact: RedactConfig{
			Headers: []string{"authorization", "x-api-key"},
			Fields:  []string{"access_token", "refresh_token", "secret", "password"},
		},
		Verify: VerifyConfig{
			DefaultStrategy: "replay",
			Replay:          ReplayConfig{AllowNetwork: false},
		},
	}

	if err := writeProjectConfig(options.ConfigPath, projectConfig); err != nil {
		return Report{}, err
	}

	report := newReport(cleanSlashPath(options.ConfigPath), intakeReportPath, "pass")
	report.SourceRefs = sourceRefs(root, projectConfig)
	report.SurfaceFingerprint = surfaceFingerprint(report.SourceRefs)
	if err := writeReport(reportRoot(options.ConfigPath), report); err != nil {
		return Report{}, err
	}
	return report, nil
}

func CheckProject(configPath string) (Report, error) {
	if configPath == "" {
		configPath = config.DefaultConfigPath
	}
	projectConfig, err := ReadProjectConfig(configPath)
	if err != nil {
		return Report{}, err
	}

	root := projectRoot(configPath)
	report := newReport(cleanSlashPath(configPath), checkReportPath, "pass")
	report.SourceRefs = sourceRefs(root, projectConfig)
	for _, entry := range projectConfig.Sources.OpenAPI {
		report.Findings = append(report.Findings, checkOpenAPI(root, entry)...)
	}
	for _, entry := range projectConfig.Sources.Docs {
		report.Findings = append(report.Findings, checkDocs(root, entry)...)
	}
	report.Status = statusFromFindings(report.Findings)
	report.SurfaceFingerprint = surfaceFingerprint(report.SourceRefs)
	if err := writeReport(reportRoot(configPath), report); err != nil {
		return Report{}, err
	}
	return report, nil
}

func newReport(configPath, reportPath, status string) Report {
	return Report{
		ObjectType:    "lumyn.source_check",
		SchemaVersion: schemaVersion,
		Status:        status,
		ConfigPath:    configPath,
		ReportPath:    cleanSlashPath(reportPath),
		SourceRefs:    []SourceRef{},
		Findings:      []Finding{},
	}
}

func sourceRefs(root string, projectConfig ProjectConfig) []SourceRef {
	refs := make([]SourceRef, 0, len(projectConfig.Sources.OpenAPI)+len(projectConfig.Sources.Docs))
	for _, entry := range projectConfig.Sources.OpenAPI {
		refs = append(refs, SourceRef{
			ID:   entry.ID,
			Kind: "openapi",
			Path: cleanSlashPath(entry.Path),
			Hash: hashPath(root, entry.Path),
		})
	}
	for _, entry := range projectConfig.Sources.Docs {
		refs = append(refs, SourceRef{
			ID:   entry.ID,
			Kind: "docs",
			Path: cleanSlashPath(entry.Path),
			Hash: hashDocsPath(root, entry.Path),
		})
	}
	return refs
}

func projectRoot(configPath string) string {
	dir := filepath.Dir(configPath)
	if dir == "" {
		return "."
	}
	return dir
}

func reportRoot(configPath string) string {
	if configPath == "" {
		return "."
	}
	if !filepath.IsAbs(configPath) {
		return "."
	}
	dir := filepath.Dir(configPath)
	if strings.HasPrefix(filepath.Base(dir), ".") {
		parent := filepath.Dir(dir)
		if parent == "" {
			return "."
		}
		return parent
	}
	return dir
}

func resolveProjectPath(root, path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Clean(filepath.Join(root, filepath.FromSlash(path)))
}

func sourcePathForConfigRoot(root, path string, preferWorkingDirectory bool) string {
	if path == "" {
		return ""
	}
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return cleanSlashPath(path)
	}
	sourcePath := filepath.FromSlash(filepath.Clean(path))
	sourceAbs := sourcePath
	if !filepath.IsAbs(sourceAbs) {
		cwd, cwdErr := os.Getwd()
		cwdCandidate := sourcePath
		if cwdErr == nil {
			cwdCandidate = filepath.Join(cwd, sourcePath)
		}
		rootCandidate := filepath.Join(rootAbs, sourcePath)
		switch {
		case preferWorkingDirectory && cwdErr == nil && pathExists(cwdCandidate):
			sourceAbs = cwdCandidate
		case pathExists(rootCandidate):
			sourceAbs = rootCandidate
		case cwdErr == nil:
			sourceAbs = cwdCandidate
		default:
			return cleanSlashPath(path)
		}
	}
	rel, err := filepath.Rel(rootAbs, sourceAbs)
	if err != nil {
		return cleanSlashPath(path)
	}
	return cleanSlashPath(rel)
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func relativePath(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil || strings.HasPrefix(rel, "..") {
		return cleanSlashPath(path)
	}
	return cleanSlashPath(rel)
}

func cleanSlashPath(path string) string {
	if path == "" {
		return ""
	}
	return filepath.ToSlash(filepath.Clean(path))
}
