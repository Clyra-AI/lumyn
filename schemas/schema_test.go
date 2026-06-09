package schemas_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

var requiredSchemas = []string{
	"workflow-contract.schema.json",
	"expected-outcome.schema.json",
	"validator.schema.json",
	"action-boundary.schema.json",
	"human-annotation.schema.json",
	"required-context.schema.json",
	"state-binding.schema.json",
	"canonical-trace.schema.json",
	"evidence-event.schema.json",
	"cassette.schema.json",
	"result-axes.schema.json",
	"proof-strength.schema.json",
	"command-result.schema.json",
	"redaction-config.schema.json",
}

func repoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("resolve repo root: %v", err)
	}
	return root
}

func TestRequiredSchemasCompile(t *testing.T) {
	root := repoRoot(t)
	for _, name := range requiredSchemas {
		name := name
		t.Run(name, func(t *testing.T) {
			schemaPath := filepath.Join(root, "schemas", name)
			if _, err := os.Stat(schemaPath); err != nil {
				t.Fatalf("missing schema: %v", err)
			}
			if _, err := jsonschema.Compile(schemaPath); err != nil {
				t.Fatalf("compile schema: %v", err)
			}
		})
	}
}

func TestRequiredSchemaNamesMatchPRD(t *testing.T) {
	root := repoRoot(t)
	payload, err := os.ReadFile(filepath.Join(root, "docs/product/prd.md"))
	if err != nil {
		t.Fatalf("read PRD: %v", err)
	}
	prd := string(payload)
	for _, name := range requiredSchemas {
		if !strings.Contains(prd, name) {
			t.Fatalf("PRD does not name required schema %s", name)
		}
	}
}

func TestRepresentativeArtifactsValidate(t *testing.T) {
	root := repoRoot(t)
	for schemaName, sample := range representativeSamples() {
		schemaName := schemaName
		sample := sample
		t.Run(schemaName, func(t *testing.T) {
			schema, err := jsonschema.Compile(filepath.Join(root, "schemas", schemaName))
			if err != nil {
				t.Fatalf("compile schema: %v", err)
			}
			if err := schema.Validate(sample); err != nil {
				t.Fatalf("validate representative sample: %v", err)
			}
		})
	}
}

func TestInvalidCommandResultFails(t *testing.T) {
	root := repoRoot(t)
	schema, err := jsonschema.Compile(filepath.Join(root, "schemas", "command-result.schema.json"))
	if err != nil {
		t.Fatalf("compile command result schema: %v", err)
	}
	invalid := map[string]any{
		"object_type":    "lumyn.command_result",
		"schema_version": "1.0",
		"command":        "verify",
		"mode":           "verify",
		"metadata":       map[string]any{},
	}
	if err := schema.Validate(invalid); err == nil {
		t.Fatal("expected missing status and evidence fields to fail validation")
	}
}

func representativeSamples() map[string]any {
	metadata := map[string]any{}
	return map[string]any{
		"command-result.schema.json": map[string]any{
			"object_type":      "lumyn.command_result",
			"schema_version":   "1.0",
			"command":          "verify",
			"status":           "pass",
			"mode":             "replay",
			"warnings":         []any{},
			"errors":           []any{},
			"artifacts":        []any{map[string]any{"path": "runs/run_123/trace.json", "type": "canonical_trace"}},
			"duration_ms":      12,
			"redaction_status": "applied",
			"metadata":         metadata,
		},
		"proof-strength.schema.json": map[string]any{
			"object_type":            "lumyn.proof_strength",
			"schema_version":         "1.0",
			"level":                  "strong",
			"weakest_required_level": "strong",
			"required_validators":    []any{"validator_customer_active"},
			"advisory_validators":    []any{},
			"summary":                "Required read-back validator passed.",
			"metadata":               metadata,
		},
		"evidence-event.schema.json": evidenceEvent(metadata),
		"cassette.schema.json": map[string]any{
			"object_type":      "lumyn.cassette",
			"schema_version":   "1.0",
			"cassette_id":      "cas_create_customer",
			"workflow_id":      "create_customer_with_readback",
			"recorded_at":      "2026-06-07T00:00:00Z",
			"lumyn_version":    "0.0.0-dev",
			"source_refs":      []any{map[string]any{"path": "docs/openapi.yaml", "hash": "sha256:test"}},
			"redaction_status": "applied",
			"evidence_events":  []any{evidenceEvent(metadata)},
			"state_bindings":   map[string]any{"customer_id": "cus_test_123"},
			"validator_inputs": map[string]any{"validator_customer_active": map[string]any{"customer_id": "cus_test_123"}},
			"replay_integrity": map[string]any{
				"source_hashes": map[string]any{"docs/openapi.yaml": "sha256:test"},
				"cassette_hash": "sha256:cassette",
				"stale_policy":  "warn_normal_fail_strict",
			},
			"metadata": metadata,
		},
		"redaction-config.schema.json": map[string]any{
			"object_type":              "lumyn.redaction_config",
			"schema_version":           "1.0",
			"rules":                    []any{map[string]any{"id": "auth_header", "target": "headers.authorization", "action": "redact"}},
			"default_redactions":       []any{"authorization", "cookie", "set-cookie"},
			"fail_closed":              true,
			"safe_to_persist_required": true,
			"safe_to_share_required":   true,
			"metadata":                 metadata,
		},
		"workflow-contract.schema.json": map[string]any{
			"object_type":    "lumyn.workflow_contract",
			"schema_version": "1.0",
			"id":             "create_customer_with_readback",
			"version":        1,
			"goal":           "Create a customer and verify active status.",
			"expected_outcome": map[string]any{
				"type": "action_completed",
			},
			"context":        map[string]any{"sources": []any{"public_api"}},
			"constraints":    map[string]any{"max_requests": 20},
			"state_bindings": map[string]any{"customer_id": "steps.create_customer.response.body.id"},
			"steps":          []any{map[string]any{"id": "create_customer", "intent": "Create a customer."}},
			"validators":     []any{map[string]any{"type": "api_state"}},
			"cleanup":        []any{map[string]any{"method": "DELETE", "path": "/customers/{customer_id}"}},
			"metadata":       metadata,
		},
		"expected-outcome.schema.json": map[string]any{
			"object_type":         "lumyn.expected_outcome",
			"schema_version":      "1.0",
			"type":                "action_completed",
			"description":         "Customer is created and active.",
			"required_validators": []any{"validator_customer_active"},
			"metadata":            metadata,
		},
		"validator.schema.json": map[string]any{
			"object_type":  "lumyn.validator",
			"schema_version": "1.0",
			"validator_id": "validator_customer_active",
			"type":         "api_state",
			"required":     true,
			"proof_cap":    "strong",
			"expect":       map[string]any{"status": "active"},
			"metadata":     metadata,
		},
		"action-boundary.schema.json": map[string]any{
			"object_type":           "lumyn.action_boundary",
			"schema_version":        "1.0",
			"boundary_id":           "customer_write_boundary",
			"allowed_paths":         []any{"/customers"},
			"forbidden_paths":       []any{"/admin"},
			"allowed_operations":    []any{"createCustomer"},
			"forbidden_operations":  []any{"deleteAccount"},
			"classification_policy": "fail_closed_on_uncertain_write",
			"metadata":              metadata,
		},
		"human-annotation.schema.json": map[string]any{
			"object_type":    "lumyn.human_annotation",
			"schema_version": "1.0",
			"annotation_id":  "ann_123",
			"run_id":         "run_123",
			"author":         "reviewer",
			"verdict":        "accepted",
			"notes":          "Looks correct.",
			"metadata":       metadata,
		},
		"required-context.schema.json": map[string]any{
			"object_type":    "lumyn.required_context",
			"schema_version": "1.0",
			"sources":        []any{"openapi", "docs"},
			"required":       []any{"customer lifecycle states"},
			"missing_policy": "fail_strict_warn_normal",
			"metadata":       metadata,
		},
		"state-binding.schema.json": map[string]any{
			"object_type":    "lumyn.state_binding",
			"schema_version": "1.0",
			"binding_id":     "binding_customer_id",
			"name":           "customer_id",
			"from":           "steps.create_customer.response.body.id",
			"required":       true,
			"confidence":     "high",
			"metadata":       metadata,
		},
		"canonical-trace.schema.json": map[string]any{
			"object_type":      "lumyn.canonical_trace",
			"schema_version":   "1.0",
			"trace_id":         "trace_123",
			"run_id":           "run_123",
			"workflow_id":      "create_customer_with_readback",
			"lumyn_version":    "0.0.0-dev",
			"started_at":       "2026-06-07T00:00:00Z",
			"finished_at":      "2026-06-07T00:00:01Z",
			"redaction_status": "applied",
			"events":           []any{evidenceEvent(metadata)},
			"proof_strength":   map[string]any{"level": "strong"},
			"metadata":         metadata,
		},
		"result-axes.schema.json": map[string]any{
			"object_type":     "lumyn.result_axes",
			"schema_version":  "1.0",
			"workflow_result": "pass",
			"proof_strength":  "strong",
			"freshness":       "fresh",
			"redaction":       "applied",
			"boundary":        "in_bounds",
			"metadata":        metadata,
		},
	}
}

func evidenceEvent(metadata map[string]any) map[string]any {
	return map[string]any{
		"object_type":      "lumyn.evidence_event",
		"schema_version":   "1.0",
		"event_id":         "evt_001",
		"run_id":           "run_123",
		"timestamp":        "2026-06-07T00:00:00Z",
		"source":           "http",
		"kind":             "http_request",
		"redaction_status": "applied",
		"raw_refs":         []any{"redacted://request/evt_001"},
		"classification":   map[string]any{"action_type": "write", "confidence": "high"},
		"operation": map[string]any{
			"method":                    "POST",
			"path":                      "/customers",
			"operation_id":              "createCustomer",
			"action_type":               "write",
			"classification_confidence": "high",
		},
		"request":    map[string]any{"headers_redacted": true, "body_ref": "redacted://request/evt_001"},
		"response":   map[string]any{"status_code": 200, "body_ref": "redacted://response/evt_001"},
		"bindings":   map[string]any{"customer_id": "cus_test_123"},
		"confidence": "high",
		"metadata":   metadata,
	}
}
