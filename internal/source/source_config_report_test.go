package source

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

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
