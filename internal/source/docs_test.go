package source

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitProjectWritesConfigAndIntakeArtifact(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.json"), []byte(validOpenAPIWithMissingMetadata), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}

	report, err := InitProject(InitOptions{
		ConfigPath:  filepath.Join(root, "lumyn.yaml"),
		OpenAPIPath: "./openapi.json",
		DocsPath:    "./docs",
	})
	if err != nil {
		t.Fatalf("InitProject: %v", err)
	}
	if report.Status != "pass" {
		t.Fatalf("report status = %q, want pass", report.Status)
	}
	if report.ReportPath != filepath.Join("runs", "source-checks", "source-intake.json") {
		t.Fatalf("report path = %q, want source intake path", report.ReportPath)
	}
	configBytes, err := os.ReadFile(filepath.Join(root, "lumyn.yaml"))
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if string(configBytes) == "" || !containsString(string(configBytes), "openapi") {
		t.Fatalf("generated config is missing sources:\n%s", configBytes)
	}
	if _, err := os.Stat(filepath.Join(root, report.ReportPath)); err != nil {
		t.Fatalf("source intake artifact missing: %v", err)
	}
}

func TestInitProjectRebasesSourcesWhenConfigLivesInSubdirectory(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(completeOpenAPIYAML), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	t.Chdir(root)

	configPath := filepath.Join(".lumyn", "lumyn.yaml")
	initReport, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	})
	if err != nil {
		t.Fatalf("InitProject: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, initReport.ReportPath)); err != nil {
		t.Fatalf("source intake artifact should be written to reported project-root path: %v", err)
	}
	projectConfig, err := ReadProjectConfig(configPath)
	if err != nil {
		t.Fatalf("ReadProjectConfig: %v", err)
	}
	if got, want := projectConfig.Sources.OpenAPI[0].Path, "../openapi.yaml"; got != want {
		t.Fatalf("OpenAPI path = %q, want %q", got, want)
	}
	if got, want := projectConfig.Sources.Docs[0].Path, "../docs"; got != want {
		t.Fatalf("docs path = %q, want %q", got, want)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, report.ReportPath)); err != nil {
		t.Fatalf("source check artifact should be written to reported project-root path: %v", err)
	}
	if hasFindingKind(report.Findings, "context_missing") || hasFindingKind(report.Findings, "command_error") {
		t.Fatalf("rebased init config should resolve generated source paths; findings=%#v", report.Findings)
	}
}

func TestCheckProjectReportsConcreteWorkflowRelevantFinding(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.json"), []byte(validOpenAPIWithMissingMetadata), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte("See [missing](missing.md)\n"), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	if _, err := InitProject(InitOptions{
		ConfigPath:  filepath.Join(root, "lumyn.yaml"),
		OpenAPIPath: "./openapi.json",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(filepath.Join(root, "lumyn.yaml"))
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if report.Status != "warning" {
		t.Fatalf("report status = %q, want warning", report.Status)
	}
	if len(report.Findings) == 0 {
		t.Fatal("expected at least one source finding")
	}
	finding := report.Findings[0]
	if finding.Kind != "docs_api_ambiguity" {
		t.Fatalf("finding kind = %q, want docs_api_ambiguity", finding.Kind)
	}
	if finding.Reference.Path == "" {
		t.Fatalf("finding reference missing path: %#v", finding.Reference)
	}
	if finding.Reference.JSONPointer == "" && finding.Reference.Line == 0 && finding.Reference.Object == "" {
		t.Fatalf("finding lacks concrete source reference: %#v", finding.Reference)
	}
	if finding.WorkflowRelevance == "" {
		t.Fatalf("finding lacks workflow relevance: %#v", finding)
	}
	if _, err := os.Stat(filepath.Join(root, report.ReportPath)); err != nil {
		t.Fatalf("source check artifact missing: %v", err)
	}
}

func TestCheckProjectResolvesRootRelativeMarkdownLinksFromProjectRoot(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(completeOpenAPIYAML), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte("Read [auth](/docs/auth.md).\n"), 0o644); err != nil {
		t.Fatalf("write docs guide: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "auth.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs auth: %v", err)
	}
	if _, err := InitProject(InitOptions{
		ConfigPath:  filepath.Join(root, "lumyn.yaml"),
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(filepath.Join(root, "lumyn.yaml"))
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if hasFindingKind(report.Findings, "context_missing") {
		t.Fatalf("root-relative docs link should resolve from project root; findings=%#v", report.Findings)
	}
}

func TestCheckProjectResolvesBracketedMarkdownLinksWithSpaces(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(completeOpenAPIYAML), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte("Read [auth](<auth guide.md>).\n"), 0o644); err != nil {
		t.Fatalf("write docs guide: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "auth guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs auth: %v", err)
	}
	if _, err := InitProject(InitOptions{
		ConfigPath:  filepath.Join(root, "lumyn.yaml"),
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(filepath.Join(root, "lumyn.yaml"))
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if hasFindingKind(report.Findings, "context_missing") {
		t.Fatalf("bracketed docs link should resolve with spaces; findings=%#v", report.Findings)
	}
}

func TestCheckProjectResolvesMarkdownLinksWithParentheses(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(completeOpenAPIYAML), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte("Read [setup](setup(v1).md).\n"), 0o644); err != nil {
		t.Fatalf("write docs guide: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "setup(v1).md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs setup: %v", err)
	}
	if _, err := InitProject(InitOptions{
		ConfigPath:  filepath.Join(root, "lumyn.yaml"),
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(filepath.Join(root, "lumyn.yaml"))
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if hasFindingKind(report.Findings, "context_missing") {
		t.Fatalf("docs link with balanced parentheses should resolve; findings=%#v", report.Findings)
	}
}

func TestCheckProjectIgnoresMarkdownLinksInFencedCode(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(completeOpenAPIYAML), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	body := "Example only:\n```md\n[placeholder](missing.md)\n```\n\n" + completeDocs
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(body), 0o644); err != nil {
		t.Fatalf("write docs guide: %v", err)
	}
	if _, err := InitProject(InitOptions{
		ConfigPath:  filepath.Join(root, "lumyn.yaml"),
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(filepath.Join(root, "lumyn.yaml"))
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if hasFindingKind(report.Findings, "context_missing") {
		t.Fatalf("docs links in fenced code blocks should be ignored; findings=%#v", report.Findings)
	}
}

func TestCheckProjectResolvesURLEncodedMarkdownLinks(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(completeOpenAPIYAML), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte("Read [auth](auth%20guide.md).\n"), 0o644); err != nil {
		t.Fatalf("write docs guide: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "auth guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs auth: %v", err)
	}
	if _, err := InitProject(InitOptions{
		ConfigPath:  filepath.Join(root, "lumyn.yaml"),
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(filepath.Join(root, "lumyn.yaml"))
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if hasFindingKind(report.Findings, "context_missing") {
		t.Fatalf("URL-encoded docs link should resolve to local path; findings=%#v", report.Findings)
	}
}

func TestCheckProjectRedactsSecretBearingBrokenMarkdownLinkTargets(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(completeOpenAPIYAML), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte("Read [setup](missing.md?api_key=super-secret-token#install).\n"+completeDocs), 0o644); err != nil {
		t.Fatalf("write docs guide: %v", err)
	}
	if _, err := InitProject(InitOptions{
		ConfigPath:  filepath.Join(root, "lumyn.yaml"),
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(filepath.Join(root, "lumyn.yaml"))
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if !hasFindingKind(report.Findings, "context_missing") {
		t.Fatalf("expected missing local reference finding; findings=%#v", report.Findings)
	}
	for _, finding := range report.Findings {
		if containsString(finding.Message, "api_key") || containsString(finding.Message, "super-secret-token") || containsString(finding.Message, "#install") {
			t.Fatalf("finding message leaked secret-bearing link target: %#v", finding)
		}
	}
}
