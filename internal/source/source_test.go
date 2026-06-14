package source

import (
	"encoding/json"
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

func TestCheckProjectFailsClosedForInvalidConfig(t *testing.T) {
	root := t.TempDir()
	_, err := CheckProject(filepath.Join(root, "missing-lumyn.yaml"))
	if err == nil {
		t.Fatal("CheckProject should fail when config is missing")
	}
}

func TestCheckProjectPassesCompleteYAMLOpenAPIAndDocs(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(completeOpenAPIYAML), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
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
	if report.Status != "pass" {
		t.Fatalf("report status = %q, want pass; findings=%#v", report.Status, report.Findings)
	}
	if len(report.Findings) != 0 {
		t.Fatalf("findings = %#v, want none", report.Findings)
	}
	if report.SurfaceFingerprint == "" {
		t.Fatal("surface fingerprint should be populated")
	}
}

func TestCheckProjectDoesNotTreatYAMLResponseDescriptionAsOperationDescription(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOperationMissingDescriptionWithResponseDescription), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if !hasFindingKind(report.Findings, "docs_api_ambiguity") {
		t.Fatalf("expected docs_api_ambiguity when only response description exists; findings=%#v", report.Findings)
	}
}

func TestCheckProjectDoesNotTreatYAMLErrorResponseSchemaAsSuccessSchema(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlSuccessResponseWithoutSchemaAndErrorSchema), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if !hasFindingKind(report.Findings, "proof_gap") {
		t.Fatalf("expected proof_gap when only non-2xx response has schema; findings=%#v", report.Findings)
	}
}

func TestCheckProjectDoesNotTreatYAMLSchemaMethodPropertyAsOperation(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlSchemaPropertyNamedLikeHTTPMethod), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if report.Status != "pass" {
		t.Fatalf("report status = %q, want pass; findings=%#v", report.Status, report.Findings)
	}
}

func TestCheckProjectKeepsParsingAfterNestedYAMLPathsSchemaProperty(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlNestedPathsSchemaPropertyBeforeLaterOperation), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if !hasFindingKind(report.Findings, "wrong_tool_or_endpoint") {
		t.Fatalf("later operation after nested schema.paths property should still be parsed; findings=%#v", report.Findings)
	}
}

func TestCheckProjectAllowsBodylessDeleteWithoutRequestSchema(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.json"), []byte(openAPIWithBodylessDelete), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.json",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if hasFindingKind(report.Findings, "validator_coverage_gap") {
		t.Fatalf("body-less DELETE should not require a request schema; findings=%#v", report.Findings)
	}
}

func TestCheckProjectReportsEmptyOAuthScopeDescriptions(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.json"), []byte(openAPIWithEmptyOAuthScopeDescription), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.json",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if !hasFindingKind(report.Findings, "auth_confusion") {
		t.Fatalf("empty OAuth scope descriptions should be reported; findings=%#v", report.Findings)
	}
}

func TestCheckProjectReportsEmptyYAMLOAuthScopeDescriptions(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithEmptyOAuthScopeDescription), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if !hasFindingKind(report.Findings, "auth_confusion") {
		t.Fatalf("empty YAML OAuth scope descriptions should be reported; findings=%#v", report.Findings)
	}
}

func TestCheckProjectReportsEmptyYAMLOAuthScopeDescriptionsWhenTypeFollowsFlows(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithLateOAuthType), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if !hasFindingKind(report.Findings, "auth_confusion") {
		t.Fatalf("YAML OAuth scopes should be reported even when type follows flows; findings=%#v", report.Findings)
	}
}

func TestCheckProjectReportsEmptyInlineYAMLOAuthScopes(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithInlineOAuthScopes), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if !hasFindingKind(report.Findings, "auth_confusion") {
		t.Fatalf("inline YAML OAuth scope descriptions should be reported; findings=%#v", report.Findings)
	}
}

func TestCheckProjectReportsEmptyFlowStyleYAMLOAuthScopes(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithFlowStyleOAuthSecuritySchemes), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if !hasFindingKind(report.Findings, "auth_confusion") {
		t.Fatalf("flow-style YAML OAuth scope descriptions should be reported; findings=%#v", report.Findings)
	}
}

func TestCheckProjectRejectsUnsupportedSwagger2(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "swagger.json"), []byte(swagger2JSONFixture), 0o644); err != nil {
		t.Fatalf("write Swagger fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./swagger.json",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	first, ok := FirstFinding(report.Findings)
	if report.Status != "fail" || !ok || first.Kind != "command_error" || !containsString(first.Message, "swagger 2.0 is not supported") {
		t.Fatalf("expected unsupported Swagger command_error; status=%q first=%#v findings=%#v", report.Status, first, report.Findings)
	}
}

func TestCheckProjectIgnoresNestedYAMLSwaggerPropertyForVersion(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithNestedSwaggerProperty), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if hasFindingKind(report.Findings, "command_error") {
		t.Fatalf("nested YAML swagger schema property should not be treated as Swagger 2.0; findings=%#v", report.Findings)
	}
}

func TestCheckProjectRequiresRootYAMLOpenAPIVersion(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithNestedOpenAPIPropertyOnly), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	first, ok := FirstFinding(report.Findings)
	if report.Status != "fail" || !ok || first.Kind != "command_error" || !containsString(first.Message, "missing openapi version") {
		t.Fatalf("expected missing root version command_error; status=%q first=%#v findings=%#v", report.Status, first, report.Findings)
	}
}

func TestCheckProjectRejectsJSONTrailingData(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.json"), []byte(openAPIWithBodylessDelete+"\n{}"), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.json",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	first, ok := FirstFinding(report.Findings)
	if report.Status != "fail" || !ok || first.Kind != "command_error" || !containsString(first.Message, "trailing data") {
		t.Fatalf("expected trailing data command_error; status=%q first=%#v findings=%#v", report.Status, first, report.Findings)
	}
}

func TestCheckProjectRequiresUsableJSONSecuritySchemes(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.json"), []byte(openAPIWithEmptySecurityScheme), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.json",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if !hasFindingKind(report.Findings, "auth_confusion") {
		t.Fatalf("empty JSON security scheme should not satisfy auth scheme coverage; findings=%#v", report.Findings)
	}
}

func TestCheckProjectRequiresNonEmptyYAMLSecuritySchemes(t *testing.T) {
	cases := []struct {
		name string
		body string
	}{
		{name: "empty components securitySchemes", body: yamlEmptySecuritySchemes},
		{name: "nested schema property named securitySchemes", body: yamlNestedSecuritySchemesProperty},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(tc.body), 0o644); err != nil {
				t.Fatalf("write OpenAPI fixture: %v", err)
			}
			if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
				t.Fatalf("create docs dir: %v", err)
			}
			if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
				t.Fatalf("write docs fixture: %v", err)
			}
			configPath := filepath.Join(root, "lumyn.yaml")
			if _, err := InitProject(InitOptions{
				ConfigPath:  configPath,
				OpenAPIPath: "./openapi.yaml",
				DocsPath:    "./docs",
			}); err != nil {
				t.Fatalf("InitProject: %v", err)
			}

			report, err := CheckProject(configPath)
			if err != nil {
				t.Fatalf("CheckProject: %v", err)
			}
			if !hasFindingKind(report.Findings, "auth_confusion") {
				t.Fatalf("expected auth_confusion for missing usable security scheme; findings=%#v", report.Findings)
			}
		})
	}
}

func TestCheckProjectDoesNotTreatNestedYAMLSecuritySchemesPropertyAsAuthScheme(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlNestedSecuritySchemesSchemaPropertyOnly), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if !hasFindingKind(report.Findings, "auth_confusion") {
		t.Fatalf("nested securitySchemes schema property should not satisfy auth scheme coverage; findings=%#v", report.Findings)
	}
}

func TestCheckProjectDoesNotTreatYAMLResponseHeaderSchemaAsBodySchema(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlResponseHeaderSchemaOnly), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if !hasFindingKind(report.Findings, "proof_gap") {
		t.Fatalf("expected proof_gap when only 2xx header schema exists; findings=%#v", report.Findings)
	}
}

func TestCheckProjectRequiresDirectYAMLRequestMediaSchema(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlRequestBodyNestedExampleSchemaOnly), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if !hasFindingKind(report.Findings, "validator_coverage_gap") {
		t.Fatalf("nested example schema should not satisfy request schema coverage; findings=%#v", report.Findings)
	}
}

func TestCheckProjectRequiresDirectFlowStyleMediaSchemas(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlFlowStyleMediaNestedExampleSchemaOnly), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if !hasFindingKind(report.Findings, "validator_coverage_gap") || !hasFindingKind(report.Findings, "proof_gap") {
		t.Fatalf("nested flow-style example schemas should not satisfy direct media schema checks; findings=%#v", report.Findings)
	}
}

func TestCheckProjectResolvesInlineYAMLRequestAndResponseContentSchemas(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithInlineRequestAndResponseContentSchemas), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if hasFindingKind(report.Findings, "validator_coverage_gap") || hasFindingKind(report.Findings, "proof_gap") {
		t.Fatalf("inline request/response content schemas should satisfy schema checks; findings=%#v", report.Findings)
	}
}

func TestCheckProjectResolvesFlowStyleYAMLOperationSchemas(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithFlowStyleOperationSchemas), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if hasFindingKind(report.Findings, "validator_coverage_gap") || hasFindingKind(report.Findings, "proof_gap") {
		t.Fatalf("flow-style operation schemas should satisfy schema checks; findings=%#v", report.Findings)
	}
}

func TestCheckProjectResolvesInlineYAMLOperationObjects(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithInlineOperationObject), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if hasFindingKind(report.Findings, "wrong_tool_or_endpoint") ||
		hasFindingKind(report.Findings, "docs_api_ambiguity") ||
		hasFindingKind(report.Findings, "validator_coverage_gap") ||
		hasFindingKind(report.Findings, "proof_gap") {
		t.Fatalf("inline operation object should satisfy metadata and schema checks; findings=%#v", report.Findings)
	}
}

func TestCheckProjectResolvesRootFlowStyleYAMLPaths(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithRootFlowStylePaths), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if hasFindingKind(report.Findings, "command_error") ||
		hasFindingKind(report.Findings, "wrong_tool_or_endpoint") ||
		hasFindingKind(report.Findings, "docs_api_ambiguity") ||
		hasFindingKind(report.Findings, "proof_gap") {
		t.Fatalf("root flow-style paths should parse as operations; findings=%#v", report.Findings)
	}
}

func TestCheckProjectHonorsYAMLDeprecatedReplacementHint(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithDeprecatedReplacementHint), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if hasFindingFixTarget(report.Findings, "deprecated_operation_guidance") {
		t.Fatalf("x-deprecated-replacement should satisfy deprecated guidance; findings=%#v", report.Findings)
	}
}

func TestCheckProjectResolvesLocalComponentRequestAndResponseRefs(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.json"), []byte(openAPIWithComponentRequestAndResponseRefs), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.json",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if hasFindingKind(report.Findings, "validator_coverage_gap") || hasFindingKind(report.Findings, "proof_gap") {
		t.Fatalf("component refs should satisfy request/response schema checks; findings=%#v", report.Findings)
	}
}

func TestCheckProjectResolvesChainedLocalComponentRefs(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.json"), []byte(openAPIWithChainedComponentRequestAndResponseRefs), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.json",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if hasFindingKind(report.Findings, "validator_coverage_gap") || hasFindingKind(report.Findings, "proof_gap") {
		t.Fatalf("chained component refs should satisfy request/response schema checks; findings=%#v", report.Findings)
	}
}

func TestCheckProjectReportsMissingDescriptionForJSONParameterRef(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.json"), []byte(openAPIWithUndescribedParameterRef), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.json",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if !hasFindingKind(report.Findings, "source_missing_metadata") {
		t.Fatalf("referenced parameter without description should be reported; findings=%#v", report.Findings)
	}
}

func TestCheckProjectReportsMissingDescriptionForChainedJSONParameterRef(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.json"), []byte(openAPIWithChainedUndescribedParameterRef), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.json",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if !hasFindingKind(report.Findings, "source_missing_metadata") {
		t.Fatalf("chained referenced parameter without description should be reported; findings=%#v", report.Findings)
	}
}

func TestCheckProjectHonorsOperationParameterOverrides(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.json"), []byte(openAPIWithOperationParameterOverride), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.json",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if hasFindingKind(report.Findings, "source_missing_metadata") {
		t.Fatalf("operation parameter override should replace stale path metadata; findings=%#v", report.Findings)
	}
}

func TestCheckProjectReportsMissingDescriptionForYAMLFlowParameter(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithUndescribedFlowParameter), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if !hasFindingKind(report.Findings, "source_missing_metadata") {
		t.Fatalf("flow-style parameter without description should be reported; findings=%#v", report.Findings)
	}
}

func TestCheckProjectReportsMissingDescriptionForYAMLFlowParameterRef(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithUndescribedFlowParameterRef), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if !hasFindingKind(report.Findings, "source_missing_metadata") {
		t.Fatalf("flow-style parameter ref without description should be reported; findings=%#v", report.Findings)
	}
}

func TestCheckProjectReportsMissingDescriptionForYAMLInlineParameterArrays(t *testing.T) {
	cases := []struct {
		name string
		body string
	}{
		{name: "operation parameters", body: yamlOpenAPIWithInlineOperationParameterArray},
		{name: "path parameters", body: yamlOpenAPIWithInlinePathParameterArray},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(tc.body), 0o644); err != nil {
				t.Fatalf("write OpenAPI fixture: %v", err)
			}
			if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
				t.Fatalf("create docs dir: %v", err)
			}
			if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
				t.Fatalf("write docs fixture: %v", err)
			}
			configPath := filepath.Join(root, "lumyn.yaml")
			if _, err := InitProject(InitOptions{
				ConfigPath:  configPath,
				OpenAPIPath: "./openapi.yaml",
				DocsPath:    "./docs",
			}); err != nil {
				t.Fatalf("InitProject: %v", err)
			}

			report, err := CheckProject(configPath)
			if err != nil {
				t.Fatalf("CheckProject: %v", err)
			}
			if !hasFindingKind(report.Findings, "source_missing_metadata") {
				t.Fatalf("inline YAML parameter array without descriptions should be reported; findings=%#v", report.Findings)
			}
		})
	}
}

func TestCheckProjectReportsMissingDescriptionForYAMLParameterRef(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithUndescribedParameterRef), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if !hasFindingKind(report.Findings, "source_missing_metadata") {
		t.Fatalf("referenced YAML parameter without description should be reported; findings=%#v", report.Findings)
	}
}

func TestReadProjectConfigAcceptsJSONConfig(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "lumyn.json")
	configBody := map[string]any{
		"version": 1,
		"sources": map[string]any{
			"openapi": []map[string]string{{"id": "public_api", "path": "./openapi.json"}},
			"docs":    []map[string]string{{"id": "docs", "path": "./docs"}},
		},
	}
	encoded, err := json.Marshal(configBody)
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}
	if err := os.WriteFile(configPath, encoded, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	parsed, err := ReadProjectConfig(configPath)
	if err != nil {
		t.Fatalf("ReadProjectConfig: %v", err)
	}
	if parsed.Sources.OpenAPI[0].Path != "./openapi.json" {
		t.Fatalf("openapi path = %q, want ./openapi.json", parsed.Sources.OpenAPI[0].Path)
	}
}

func TestReadProjectConfigAcceptsYAMLSourcesInEitherOrder(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "lumyn.yaml")
	configBody := `version: 1
sources:
  docs:
    - id: docs
      path: ./docs
  openapi:
    - id: public_api
      path: ./openapi.yaml
`
	if err := os.WriteFile(configPath, []byte(configBody), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	parsed, err := ReadProjectConfig(configPath)
	if err != nil {
		t.Fatalf("ReadProjectConfig: %v", err)
	}
	if parsed.Sources.OpenAPI[0].Path != "./openapi.yaml" || parsed.Sources.Docs[0].Path != "./docs" {
		t.Fatalf("parsed sources = %#v", parsed.Sources)
	}
}

func TestReadProjectConfigAcceptsYAMLSourceEntryKeysInAnyOrderAndInlineComments(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "lumyn.yaml")
	configBody := `version: 1
sources:
  openapi:
    - path: ./openapi.yaml # local API contract
      id: public_api
  docs:
    - path: "./docs" # local docs
      id: docs
`
	if err := os.WriteFile(configPath, []byte(configBody), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	parsed, err := ReadProjectConfig(configPath)
	if err != nil {
		t.Fatalf("ReadProjectConfig: %v", err)
	}
	if parsed.Sources.OpenAPI[0].ID != "public_api" || parsed.Sources.OpenAPI[0].Path != "./openapi.yaml" {
		t.Fatalf("openapi source = %#v", parsed.Sources.OpenAPI[0])
	}
	if parsed.Sources.Docs[0].ID != "docs" || parsed.Sources.Docs[0].Path != "./docs" {
		t.Fatalf("docs source = %#v", parsed.Sources.Docs[0])
	}
}

func TestCheckProjectReportsInvalidOpenAPIAndMissingDocs(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.json"), []byte(`{"info": {}}`), 0o644); err != nil {
		t.Fatalf("write invalid OpenAPI fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.json",
		DocsPath:    "./missing-docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if report.Status != "fail" {
		t.Fatalf("report status = %q, want fail; findings=%#v", report.Status, report.Findings)
	}
	if !HasErrorFindings(report.Findings) {
		t.Fatalf("expected error finding: %#v", report.Findings)
	}
	first, ok := FirstFinding(report.Findings)
	if !ok || first.Kind != "command_error" {
		t.Fatalf("first finding = %#v, %v; want command_error", first, ok)
	}
}

func TestCheckProjectResolvesYAMLComponentResponseRef(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithComponentResponseRef), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if hasFindingKind(report.Findings, "proof_gap") {
		t.Fatalf("YAML component response ref should satisfy response schema checks; findings=%#v", report.Findings)
	}
}

func TestCheckProjectResolvesInlineYAMLComponentRefs(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithInlineComponentRefs), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if hasFindingKind(report.Findings, "validator_coverage_gap") || hasFindingKind(report.Findings, "proof_gap") {
		t.Fatalf("inline YAML component refs should satisfy schema checks; findings=%#v", report.Findings)
	}
}

func TestCheckProjectResolvesQuotedInlineYAMLComponentRefs(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithQuotedInlineComponentRefs), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if hasFindingKind(report.Findings, "validator_coverage_gap") || hasFindingKind(report.Findings, "proof_gap") {
		t.Fatalf("quoted inline YAML component refs should satisfy schema checks; findings=%#v", report.Findings)
	}
}

func TestCheckProjectResolvesJSONStyleFlowYAMLComponentsAtRoot(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithJSONStyleFlowComponentsAndNestedComponentsProperty), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if hasFindingKind(report.Findings, "validator_coverage_gap") || hasFindingKind(report.Findings, "proof_gap") {
		t.Fatalf("JSON-style flow refs and root components should satisfy schema checks; findings=%#v", report.Findings)
	}
}

func TestCheckProjectResolvesFlowStyleYAMLComponentSchemas(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithFlowStyleComponentSchemas), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if hasFindingKind(report.Findings, "validator_coverage_gap") || hasFindingKind(report.Findings, "proof_gap") {
		t.Fatalf("flow-style YAML component schemas should satisfy schema checks; findings=%#v", report.Findings)
	}
}

func TestCheckProjectPreservesYAMLPathKeysWithColons(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithColonPathKey), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	found := false
	for _, finding := range report.Findings {
		if finding.Kind == "docs_api_ambiguity" && finding.Reference.Object == "GET /v1/books:batchGet" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected finding to preserve colon path key; findings=%#v", report.Findings)
	}
}

func TestCheckProjectFlagsDuplicateAgentNamesOnSamePath(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.json"), []byte(openAPIWithDuplicateSamePathNames), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.json",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if !hasFindingKind(report.Findings, "docs_api_ambiguity") {
		t.Fatalf("duplicate same-path agent-facing names should be reported; findings=%#v", report.Findings)
	}
}

func TestCheckProjectFlagsDuplicateOperationIDsWithDistinctSummaries(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.json"), []byte(openAPIWithDuplicateOperationIDsDistinctSummaries), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.json",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	found := false
	for _, finding := range report.Findings {
		if finding.FixTarget == "operation_id_disambiguation" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("duplicate operationIds should be reported independently of summaries; findings=%#v", report.Findings)
	}
}

func TestCheckProjectCarriesYAMLPathItemParameters(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithUndescribedPathParameter), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if !hasFindingKind(report.Findings, "source_missing_metadata") {
		t.Fatalf("shared YAML path parameter without description should be reported; findings=%#v", report.Findings)
	}
}

func TestCheckProjectCarriesYAMLPathItemParametersAfterMethods(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithLateUndescribedPathParameter), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    "./docs",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if !hasFindingKind(report.Findings, "source_missing_metadata") {
		t.Fatalf("late shared YAML path parameter without description should be reported; findings=%#v", report.Findings)
	}
}

func TestCheckProjectFingerprintIgnoresGeneratedReports(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(completeOpenAPIYAML), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    ".",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	first, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("first CheckProject: %v", err)
	}
	second, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("second CheckProject: %v", err)
	}
	if first.SurfaceFingerprint != second.SurfaceFingerprint {
		t.Fatalf("surface fingerprint changed after generated report write: first=%s second=%s", first.SurfaceFingerprint, second.SurfaceFingerprint)
	}
}

func TestCheckProjectIgnoresGeneratedDocsFindings(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(completeOpenAPIYAML), 0o644); err != nil {
		t.Fatalf("write OpenAPI fixture: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "guide.md"), []byte(completeDocs), 0o644); err != nil {
		t.Fatalf("write docs fixture: %v", err)
	}
	generatedDir := filepath.Join(root, ".factory", "artifacts")
	if err := os.MkdirAll(generatedDir, 0o755); err != nil {
		t.Fatalf("create generated artifact dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(generatedDir, "generated.md"), []byte("Generated [missing](missing.md).\n"), 0o644); err != nil {
		t.Fatalf("write generated docs fixture: %v", err)
	}
	configPath := filepath.Join(root, "lumyn.yaml")
	if _, err := InitProject(InitOptions{
		ConfigPath:  configPath,
		OpenAPIPath: "./openapi.yaml",
		DocsPath:    ".",
	}); err != nil {
		t.Fatalf("InitProject: %v", err)
	}

	report, err := CheckProject(configPath)
	if err != nil {
		t.Fatalf("CheckProject: %v", err)
	}
	if hasFindingKind(report.Findings, "context_missing") {
		t.Fatalf("generated docs should not contribute source findings; findings=%#v", report.Findings)
	}
}

func TestFindingHelpersClassifyProofAndSecurity(t *testing.T) {
	findings := []Finding{
		{Kind: "proof_gap", Severity: "warning"},
		{Kind: "auth_confusion", Severity: "warning"},
	}
	if !ContainsProofGap(findings) {
		t.Fatal("ContainsProofGap should detect proof_gap")
	}
	if !ContainsSecurityRelevantFinding(findings) {
		t.Fatal("ContainsSecurityRelevantFinding should detect auth_confusion")
	}
	if HasErrorFindings(findings) {
		t.Fatal("HasErrorFindings should be false for warnings")
	}
	findings = append(findings, Finding{Kind: "command_error", Severity: "error"})
	if !HasErrorFindings(findings) {
		t.Fatal("HasErrorFindings should detect error severity")
	}
}

const validOpenAPIWithMissingMetadata = `{
  "openapi": "3.0.3",
  "info": {"title": "Fixture API", "version": "1.0.0"},
  "paths": {
    "/customers": {
      "post": {
        "operationId": "createCustomer",
        "responses": {
          "201": {"description": "created"}
        }
      }
    }
  }
}`

const openAPIWithBodylessDelete = `{
  "openapi": "3.0.3",
  "info": {"title": "Fixture API", "version": "1.0.0"},
  "paths": {
    "/customers/{id}": {
      "delete": {
        "operationId": "deleteCustomer",
        "summary": "Delete customer",
        "description": "Delete one customer.",
        "parameters": [
          {"name": "id", "in": "path", "description": "Customer ID"}
        ],
        "responses": {
          "204": {
            "description": "Deleted.",
            "content": {
              "application/json": {
                "schema": {"type": "object"}
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "apiKeyAuth": {"type": "apiKey", "in": "header", "name": "X-API-Key"}
    }
  }
}`

const openAPIWithEmptyOAuthScopeDescription = `{
  "openapi": "3.0.3",
  "info": {"title": "Fixture API", "version": "1.0.0"},
  "paths": {
    "/customers": {
      "get": {
        "operationId": "listCustomers",
        "summary": "List customers",
        "description": "List customers.",
        "responses": {
          "200": {
            "description": "Customers.",
            "content": {
              "application/json": {
                "schema": {"type": "object"}
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "oauth": {
        "type": "oauth2",
        "flows": {
          "authorizationCode": {
            "authorizationUrl": "https://example.com/auth",
            "tokenUrl": "https://example.com/token",
            "scopes": {
              "customers:read": ""
            }
          }
        }
      }
    }
  }
}`

const openAPIWithEmptySecurityScheme = `{
  "openapi": "3.0.3",
  "info": {"title": "Fixture API", "version": "1.0.0"},
  "paths": {
    "/customers": {
      "get": {
        "operationId": "listCustomers",
        "summary": "List customers",
        "description": "List customers.",
        "responses": {
          "200": {
            "description": "Customers.",
            "content": {
              "application/json": {
                "schema": {"type": "object"}
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "ApiKeyAuth": {}
    }
  }
}`

const yamlOpenAPIWithEmptyOAuthScopeDescription = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    get:
      operationId: listCustomers
      summary: List customers
      description: List customers.
      responses:
        "200":
          description: Customers.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    oauth:
      type: oauth2
      flows:
        authorizationCode:
          authorizationUrl: https://example.com/auth
          tokenUrl: https://example.com/token
          scopes:
            customers:read: ""
`

const yamlOpenAPIWithLateOAuthType = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    get:
      operationId: listCustomers
      summary: List customers
      description: List customers.
      responses:
        "200":
          description: Customers.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    oauth:
      flows:
        authorizationCode:
          authorizationUrl: https://example.com/auth
          tokenUrl: https://example.com/token
          scopes:
            customers:read: ""
      type: oauth2
`

const yamlOpenAPIWithInlineOAuthScopes = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    get:
      operationId: listCustomers
      summary: List customers
      description: List customers.
      responses:
        "200":
          description: Customers.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    oauth:
      type: oauth2
      flows:
        authorizationCode:
          authorizationUrl: https://example.com/auth
          tokenUrl: https://example.com/token
          scopes: {customers:read: ""}
`

const yamlOpenAPIWithFlowStyleOAuthSecuritySchemes = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    get:
      operationId: listCustomers
      summary: List customers
      description: List customers.
      responses:
        "200":
          description: Customers.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes: {oauth: {type: oauth2, flows: {authorizationCode: {authorizationUrl: https://example.com/auth, tokenUrl: https://example.com/token, scopes: {customers:read: ""}}}}}
`

const openAPIWithUndescribedParameterRef = `{
  "openapi": "3.0.3",
  "info": {"title": "Fixture API", "version": "1.0.0"},
  "paths": {
    "/customers/{id}": {
      "parameters": [
        {"$ref": "#/components/parameters/CustomerId"}
      ],
      "get": {
        "operationId": "getCustomer",
        "summary": "Get customer",
        "description": "Get one customer.",
        "responses": {
          "200": {
            "description": "Customer.",
            "content": {
              "application/json": {
                "schema": {"type": "object"}
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "apiKeyAuth": {"type": "apiKey", "in": "header", "name": "X-API-Key"}
    },
    "parameters": {
      "CustomerId": {
        "name": "id",
        "in": "path"
      }
    }
  }
}`

const openAPIWithChainedUndescribedParameterRef = `{
  "openapi": "3.0.3",
  "info": {"title": "Fixture API", "version": "1.0.0"},
  "paths": {
    "/customers/{id}": {
      "parameters": [
        {"$ref": "#/components/parameters/CustomerIdAlias"}
      ],
      "get": {
        "operationId": "getCustomer",
        "summary": "Get customer",
        "description": "Get one customer.",
        "responses": {
          "200": {
            "description": "Customer.",
            "content": {
              "application/json": {
                "schema": {"type": "object"}
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "apiKeyAuth": {"type": "apiKey", "in": "header", "name": "X-API-Key"}
    },
    "parameters": {
      "CustomerIdAlias": {"$ref": "#/components/parameters/CustomerId"},
      "CustomerId": {
        "name": "id",
        "in": "path"
      }
    }
  }
}`

const openAPIWithOperationParameterOverride = `{
  "openapi": "3.0.3",
  "info": {"title": "Fixture API", "version": "1.0.0"},
  "paths": {
    "/customers/{id}": {
      "parameters": [
        {"name": "id", "in": "path"}
      ],
      "get": {
        "operationId": "getCustomer",
        "summary": "Get customer",
        "description": "Get one customer.",
        "parameters": [
          {"name": "id", "in": "path", "description": "Customer ID"}
        ],
        "responses": {
          "200": {
            "description": "Customer.",
            "content": {
              "application/json": {
                "schema": {"type": "object"}
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "apiKeyAuth": {"type": "apiKey", "in": "header", "name": "X-API-Key"}
    }
  }
}`

const openAPIWithDuplicateSamePathNames = `{
  "openapi": "3.0.3",
  "info": {"title": "Fixture API", "version": "1.0.0"},
  "paths": {
    "/customers": {
      "get": {
        "operationId": "listCustomers",
        "summary": "Customers",
        "description": "List customers.",
        "responses": {
          "200": {
            "description": "Customers.",
            "content": {
              "application/json": {
                "schema": {"type": "object"}
              }
            }
          }
        }
      },
      "post": {
        "operationId": "createCustomer",
        "summary": "Customers",
        "description": "Create customer.",
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {"type": "object"}
            }
          }
        },
        "responses": {
          "201": {
            "description": "Created.",
            "content": {
              "application/json": {
                "schema": {"type": "object"}
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "apiKeyAuth": {"type": "apiKey", "in": "header", "name": "X-API-Key"}
    }
  }
}`

const openAPIWithDuplicateOperationIDsDistinctSummaries = `{
  "openapi": "3.0.3",
  "info": {"title": "Fixture API", "version": "1.0.0"},
  "paths": {
    "/customers": {
      "get": {
        "operationId": "syncCustomer",
        "summary": "List customers",
        "description": "List customers.",
        "responses": {
          "200": {
            "description": "Customers.",
            "content": {
              "application/json": {
                "schema": {"type": "object"}
              }
            }
          }
        }
      }
    },
    "/customers/{id}/sync": {
      "post": {
        "operationId": "syncCustomer",
        "summary": "Sync one customer",
        "description": "Sync one customer.",
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {"type": "object"}
            }
          }
        },
        "responses": {
          "200": {
            "description": "Synced.",
            "content": {
              "application/json": {
                "schema": {"type": "object"}
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "apiKeyAuth": {"type": "apiKey", "in": "header", "name": "X-API-Key"}
    }
  }
}`

const yamlOpenAPIWithUndescribedParameterRef = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers/{id}:
    parameters:
      - $ref: "#/components/parameters/CustomerId"
    get:
      operationId: getCustomer
      summary: Get customer
      description: Get one customer.
      responses:
        "200":
          description: Customer.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
  parameters:
    CustomerId:
      name: id
      in: path
`

const yamlOpenAPIWithUndescribedFlowParameter = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers/{id}:
    get:
      operationId: getCustomer
      summary: Get customer
      description: Get one customer.
      parameters:
        - { name: id, in: path }
      responses:
        "200":
          description: Customer.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlOpenAPIWithUndescribedFlowParameterRef = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers/{id}:
    get:
      operationId: getCustomer
      summary: Get customer
      description: Get one customer.
      parameters:
        - { $ref: "#/components/parameters/CustomerId" }
      responses:
        "200":
          description: Customer.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
  parameters:
    CustomerId:
      name: id
      in: path
`

const completeOpenAPIYAML = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    get:
      operationId: listCustomers
      summary: List customers
      description: List customers with pagination.
      parameters:
        - name: page
          in: query
          description: Page number.
      responses:
        "200":
          description: Customers.
          content:
            application/json:
              schema:
                type: object
    post:
      operationId: createCustomer
      summary: Create customer
      description: Create one customer.
      requestBody:
        content:
          application/json:
            schema:
              type: object
      responses:
        "201":
          description: Created.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlOperationMissingDescriptionWithResponseDescription = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    get:
      operationId: listCustomers
      responses:
        "200":
          description: Response object description should not describe the operation.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlSuccessResponseWithoutSchemaAndErrorSchema = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    get:
      operationId: getCustomer
      summary: Get customer
      description: Get one customer.
      responses:
        "200":
          description: Customer read-back without schema.
        "404":
          description: Missing customer error.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlSchemaPropertyNamedLikeHTTPMethod = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /reports:
    get:
      operationId: getReport
      summary: Get report
      description: Get one report.
      responses:
        "200":
          description: Report response.
          content:
            application/json:
              schema:
                type: object
                properties:
                  get:
                    type: string
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlNestedPathsSchemaPropertyBeforeLaterOperation = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    get:
      operationId: listCustomers
      summary: List customers
      description: List customers.
      responses:
        "200":
          description: Customers.
          content:
            application/json:
              schema:
                type: object
                properties:
                  paths:
                    type: array
  /orders:
    get:
      summary: List orders
      description: List orders.
      responses:
        "200":
          description: Orders.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const swagger2JSONFixture = `{
  "swagger": "2.0",
  "info": {"title": "Fixture API", "version": "1.0.0"},
  "paths": {
    "/customers": {
      "get": {
        "operationId": "listCustomers",
        "summary": "List customers",
        "description": "List customers.",
        "responses": {
          "200": {
            "description": "Customers.",
            "schema": {"type": "object"}
          }
        }
      }
    }
  }
}`

const yamlOpenAPIWithNestedSwaggerProperty = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    get:
      operationId: listCustomers
      summary: List customers
      description: List customers.
      responses:
        "200":
          description: Customers.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
  schemas:
    Customer:
      type: object
      properties:
        swagger:
          type: string
`

const yamlOpenAPIWithNestedOpenAPIPropertyOnly = `info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    get:
      operationId: listCustomers
      summary: List customers
      description: List customers.
      responses:
        "200":
          description: Customers.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
  schemas:
    Metadata:
      type: object
      properties:
        openapi:
          type: string
`

const yamlEmptySecuritySchemes = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    get:
      operationId: listCustomers
      summary: List customers
      description: List customers.
      responses:
        "200":
          description: Customers.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes: {}
`

const yamlNestedSecuritySchemesProperty = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    get:
      operationId: listCustomers
      summary: List customers
      description: List customers.
      responses:
        "200":
          description: Customers.
          content:
            application/json:
              schema:
                type: object
                properties:
                  securitySchemes:
                    type: string
`

const yamlNestedSecuritySchemesSchemaPropertyOnly = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    get:
      operationId: listCustomers
      summary: List customers
      description: List customers.
      responses:
        "200":
          description: Customers.
          content:
            application/json:
              schema:
                type: object
components:
  schemas:
    Agent:
      type: object
      properties:
        securitySchemes:
          type: object
          properties:
            apiKeyAuth:
              type: string
`

const yamlResponseHeaderSchemaOnly = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    get:
      operationId: listCustomers
      summary: List customers
      description: List customers.
      responses:
        "200":
          description: Customers.
          headers:
            X-Request-ID:
              schema:
                type: string
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlRequestBodyNestedExampleSchemaOnly = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    post:
      operationId: createCustomer
      summary: Create customer
      description: Create one customer.
      requestBody:
        content:
          application/json:
            examples:
              one:
                value:
                  schema:
                    type: object
      responses:
        "201":
          description: Created.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlFlowStyleMediaNestedExampleSchemaOnly = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    post:
      operationId: createCustomer
      summary: Create customer
      description: Create one customer.
      requestBody:
        content:
          application/json: { examples: { e: { value: { schema: {} } } } }
      responses:
        "200":
          description: Customer.
          content:
            application/json: { examples: { e: { value: { schema: {} } } } }
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlOpenAPIWithInlineRequestAndResponseContentSchemas = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    post:
      operationId: createCustomer
      summary: Create customer
      description: Create one customer.
      requestBody: {content: {application/json: {schema: {type: object}}}}
      responses:
        "201": {description: Created, content: {application/json: {schema: {type: object}}}}
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlOpenAPIWithFlowStyleOperationSchemas = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    post:
      operationId: createCustomer
      summary: Create customer
      description: Create one customer.
      requestBody:
        content:
          application/json: { schema: { type: object } }
      responses:
        "201":
          description: Created.
          content:
            application/json: { schema: { type: object } }
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlOpenAPIWithInlineOperationObject = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    post: {operationId: createCustomer, summary: Create customer, description: Create one customer., requestBody: {content: {application/json: {schema: {type: object}}}}, responses: {"201": {description: Created., content: {application/json: {schema: {type: object}}}}}}
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlOpenAPIWithRootFlowStylePaths = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths: {/customers: {get: {operationId: listCustomers, summary: List customers, description: List customers., responses: {"200": {description: Customers., content: {application/json: {schema: {type: object}}}}}}}}
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlOpenAPIWithDeprecatedReplacementHint = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /v1/customers:
    get:
      operationId: listLegacyCustomers
      summary: Legacy customers
      description: Legacy endpoint.
      deprecated: true
      x-deprecated-replacement: listCustomers
      responses:
        "200":
          description: Customers.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const openAPIWithComponentRequestAndResponseRefs = `{
  "openapi": "3.0.3",
  "info": {"title": "Fixture API", "version": "1.0.0"},
  "paths": {
    "/customers": {
      "post": {
        "operationId": "createCustomer",
        "summary": "Create customer",
        "description": "Create one customer.",
        "requestBody": {"$ref": "#/components/requestBodies/CustomerWrite"},
        "responses": {
          "201": {"$ref": "#/components/responses/CustomerRead"}
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "apiKeyAuth": {"type": "apiKey", "in": "header", "name": "X-API-Key"}
    },
    "requestBodies": {
      "CustomerWrite": {
        "content": {
          "application/json": {
            "schema": {"type": "object"}
          }
        }
      }
    },
    "responses": {
      "CustomerRead": {
        "description": "Customer read-back.",
        "content": {
          "application/json": {
            "schema": {"type": "object"}
          }
        }
      }
    }
  }
}`

const openAPIWithChainedComponentRequestAndResponseRefs = `{
  "openapi": "3.0.3",
  "info": {"title": "Fixture API", "version": "1.0.0"},
  "paths": {
    "/customers": {
      "post": {
        "operationId": "createCustomer",
        "summary": "Create customer",
        "description": "Create one customer.",
        "requestBody": {"$ref": "#/components/requestBodies/CustomerWriteAlias"},
        "responses": {
          "201": {"$ref": "#/components/responses/CustomerReadAlias"}
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "apiKeyAuth": {"type": "apiKey", "in": "header", "name": "X-API-Key"}
    },
    "requestBodies": {
      "CustomerWriteAlias": {"$ref": "#/components/requestBodies/CustomerWrite"},
      "CustomerWrite": {
        "content": {
          "application/json": {
            "schema": {"type": "object"}
          }
        }
      }
    },
    "responses": {
      "CustomerReadAlias": {"$ref": "#/components/responses/CustomerRead"},
      "CustomerRead": {
        "description": "Customer read-back.",
        "content": {
          "application/json": {
            "schema": {"type": "object"}
          }
        }
      }
    }
  }
}`

const yamlOpenAPIWithComponentResponseRef = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    get:
      operationId: listCustomers
      summary: List customers
      description: List customers.
      responses:
        "200":
          $ref: "#/components/responses/CustomerList"
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
  responses:
    CustomerList:
      description: Customer list.
      content:
        application/json:
          schema:
            type: object
`

const yamlOpenAPIWithInlineComponentRefs = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    post:
      operationId: createCustomer
      summary: Create customer
      description: Create one customer.
      requestBody: { $ref: "#/components/requestBodies/CustomerWrite" }
      responses:
        "201": { $ref: "#/components/responses/CustomerRead" }
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
  requestBodies:
    CustomerWrite:
      content:
        application/json:
          schema:
            type: object
  responses:
    CustomerRead:
      description: Customer read-back.
      content:
        application/json:
          schema:
            type: object
`

const yamlOpenAPIWithQuotedInlineComponentRefs = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    post:
      operationId: createCustomer
      summary: Create customer
      description: Create one customer.
      requestBody: { "$ref": "#/components/requestBodies/CustomerWrite" }
      responses:
        "201": { "$ref": "#/components/responses/CustomerRead" }
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
  requestBodies:
    CustomerWrite:
      content:
        application/json:
          schema:
            type: object
  responses:
    CustomerRead:
      description: Customer read-back.
      content:
        application/json:
          schema:
            type: object
`

const yamlOpenAPIWithJSONStyleFlowComponentsAndNestedComponentsProperty = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    post:
      operationId: createCustomer
      summary: Create customer
      description: Create one customer.
      requestBody: {"$ref":"#/components/requestBodies/CreateCustomer"}
      responses:
        "200": {"$ref":"#/components/responses/Customer"}
components:
  schemas:
    Customer:
      type: object
      properties:
        components:
          type: string
  requestBodies:
    CreateCustomer:
      content:
        application/json: {"schema":{"type":"object"}}
  responses:
    Customer:
      description: Customer.
      content:
        application/json: {"schema":{"type":"object"}}
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlOpenAPIWithFlowStyleComponentSchemas = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    post:
      operationId: createCustomer
      summary: Create customer
      description: Create one customer.
      requestBody:
        $ref: "#/components/requestBodies/CustomerWrite"
      responses:
        "201":
          $ref: "#/components/responses/CustomerRead"
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
  requestBodies:
    CustomerWrite:
      content:
        application/json: { schema: { type: object } }
  responses:
    CustomerRead:
      description: Customer read-back.
      content:
        application/json: { schema: { type: object } }
`

const yamlOpenAPIWithColonPathKey = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /v1/books:batchGet:
    get:
      operationId: batchGetBooks
      responses:
        "200":
          description: Books.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlOpenAPIWithUndescribedPathParameter = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers/{id}:
    parameters:
      - name: id
        in: path
    get:
      operationId: getCustomer
      summary: Get customer
      description: Get one customer.
      responses:
        "200":
          description: Customer.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlOpenAPIWithLateUndescribedPathParameter = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers/{id}:
    get:
      operationId: getCustomer
      summary: Get customer
      description: Get one customer.
      responses:
        "200":
          description: Customer.
          content:
            application/json:
              schema:
                type: object
    parameters:
      - name: id
        in: path
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlOpenAPIWithInlineOperationParameterArray = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers/{id}:
    get:
      operationId: getCustomer
      summary: Get customer
      description: Get one customer.
      parameters: [{name: id, in: path}]
      responses:
        "200":
          description: Customer.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlOpenAPIWithInlinePathParameterArray = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers/{id}:
    parameters: [{name: id, in: path}]
    get:
      operationId: getCustomer
      summary: Get customer
      description: Get one customer.
      responses:
        "200":
          description: Customer.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const completeDocs = `# API Guide

Use the API key header for auth.
Retry 429 responses with backoff.
Respect the rate limit.
Use pagination with the page parameter.
Create requests are idempotent when an idempotency key is supplied.
`

func containsString(haystack, needle string) bool {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}

func hasFindingKind(findings []Finding, kind string) bool {
	for _, finding := range findings {
		if finding.Kind == kind {
			return true
		}
	}
	return false
}

func hasFindingFixTarget(findings []Finding, fixTarget string) bool {
	for _, finding := range findings {
		if finding.FixTarget == fixTarget {
			return true
		}
	}
	return false
}
