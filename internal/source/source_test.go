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
