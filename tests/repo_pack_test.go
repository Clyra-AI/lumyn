package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("resolve repo root: %v", err)
	}
	return root
}

func TestOperatingPackExists(t *testing.T) {
	root := repoRoot(t)
	for _, relativePath := range []string{
		"AGENTS.md",
		"WORKFLOW.md",
		"README.md",
		"docs/product/prd.md",
		"docs/product/plan.md",
		"docs/architecture/adr-0002-provider-sponsored-customer-controlled-migrations.md",
		".factory/artifacts/prd-to-plan/lumyn-migration-mvp/context-brief.json",
		".factory/artifacts/prd-to-plan/lumyn-migration-mvp/execution-plan.json",
	} {
		if _, err := os.Stat(filepath.Join(root, relativePath)); err != nil {
			t.Fatalf("%s: %v", relativePath, err)
		}
	}
}

func TestPRDReferencesAreRepoRelative(t *testing.T) {
	root := repoRoot(t)
	payload, err := os.ReadFile(filepath.Join(root, "docs/product/prd.md"))
	if err != nil {
		t.Fatalf("read prd: %v", err)
	}
	prd := string(payload)
	if !strings.Contains(prd, "Provider-Sponsored Verified API Migrations") {
		t.Fatal("prd should name the provider-sponsored verified API migration MVP")
	}
	if !strings.Contains(prd, "API Consumer Job") {
		t.Fatal("prd should preserve the customer-side job and authority")
	}
	if strings.Contains(prd, "/"+("Users")+"/") {
		t.Fatal("prd should not contain machine-local user paths")
	}
}
