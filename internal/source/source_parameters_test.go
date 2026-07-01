package source

import (
	"os"
	"path/filepath"
	"testing"
)

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
