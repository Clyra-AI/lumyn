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
		".factory/artifacts/prd-to-plan/lumyn-mvp/context-brief.json",
		".factory/artifacts/prd-to-plan/lumyn-mvp/execution-plan.json",
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
	if !strings.Contains(prd, "Lumyn OSS MVP") {
		t.Fatal("prd should name Lumyn OSS MVP")
	}
	if strings.Contains(prd, "/"+("Users")+"/") {
		t.Fatal("prd should not contain machine-local user paths")
	}
}
