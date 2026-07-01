package source

import (
	"encoding/json"
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
