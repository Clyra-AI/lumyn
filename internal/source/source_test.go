package source

import (
	"os"
	"path/filepath"
	"testing"
)

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

func TestCheckProjectRejectsNullJSONMediaSchemas(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.json"), []byte(openAPIWithNullMediaSchemas), 0o644); err != nil {
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
	if !hasFindingKind(report.Findings, "validator_coverage_gap") || !hasFindingKind(report.Findings, "proof_gap") {
		t.Fatalf("null media schemas should not satisfy request/response coverage; findings=%#v", report.Findings)
	}
}

func TestCheckProjectRejectsNullYAMLMediaSchemas(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithNullMediaSchemas), 0o644); err != nil {
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
		t.Fatalf("null YAML media schemas should not satisfy request/response coverage; findings=%#v", report.Findings)
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

func TestCheckProjectReportsEmptyBlockFlowStyleYAMLOAuthScopes(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithBlockFlowStyleOAuthSecurityScheme), 0o644); err != nil {
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
		t.Fatalf("block flow-style YAML OAuth scope descriptions should be reported; findings=%#v", report.Findings)
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

func TestCheckProjectRequiresUsableYAMLSecuritySchemes(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithEmptyNamedSecurityScheme), 0o644); err != nil {
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
		t.Fatalf("empty YAML security scheme should not satisfy auth scheme coverage; findings=%#v", report.Findings)
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

func TestCheckProjectResolvesFlowStyleYAMLPathItems(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithFlowStylePathItem), 0o644); err != nil {
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
		t.Fatalf("flow-style path item should parse as operations; findings=%#v", report.Findings)
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

func TestCheckProjectResolvesFlowStyleYAMLComponentGroups(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "openapi.yaml"), []byte(yamlOpenAPIWithFlowStyleComponentGroups), 0o644); err != nil {
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
		t.Fatalf("flow-style YAML component groups should satisfy request/response schema checks; findings=%#v", report.Findings)
	}
}
