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
		"docs/architecture/adr-0003-services-led-bounded-agent-migration-execution.md",
		"docs/architecture/adr-0004-provider-originated-api-update-delivery.md",
		"docs/architecture/adr-0005-customer-selected-agent-runners.md",
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
	if !strings.Contains(prd, "Provider-Originated API Update Delivery") {
		t.Fatal("prd should name the v3.1 provider-originated API update-delivery MVP")
	}
	if !strings.Contains(prd, "Primary API Provider Job") {
		t.Fatal("prd should name the provider-side job and buyer")
	}
	if !strings.Contains(prd, "API Consumer Job") {
		t.Fatal("prd should preserve the consumer-side job and authority")
	}
	if !strings.Contains(prd, "53 item-level closure units") {
		t.Fatal("prd should state the v3 item-level acceptance count")
	}
	if !strings.Contains(prd, "Codex and Claude Code") ||
		!strings.Contains(prd, "provider_sponsored_lumyn_managed") ||
		!strings.Contains(prd, "agent_execution_policy") ||
		!strings.Contains(prd, "real participating repository") {
		t.Fatal("prd should bind conditional launch runners, funding, and real agent proof")
	}
	if strings.Contains(prd, "/"+("Users")+"/") {
		t.Fatal("prd should not contain machine-local user paths")
	}
}
