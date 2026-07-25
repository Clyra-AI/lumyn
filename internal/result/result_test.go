package result

import (
	"encoding/json"
	"testing"
)

func TestCommandResultJSONEnvelope(t *testing.T) {
	payload := CommandResult{
		ObjectType:           "lumyn.command_result",
		SchemaVersion:        "1.0",
		Metadata:             map[string]any{"lumyn_version": "0.0.0"},
		Command:              "verify",
		Status:               "pass",
		Mode:                 "verify",
		Warnings:             []string{},
		Errors:               []CommandError{},
		Artifacts:            []ArtifactRef{},
		DurationMS:           0,
		RedactionStatus:      "applied",
		FindingKind:          "none",
		ProofStrength:        "unknown",
		ActionBoundaryStatus: "not_configured",
		SecurityRelevance:    "none",
		FixTarget:            "not_applicable",
		SurfaceFingerprint:   "not_applicable",
		EvalMode:             "not_applicable",
		ModelProviderMetadata: ModelProviderMetadata{
			Applicable:    false,
			SemanticRole:  "model_provider",
			ModelProvider: "not_applicable",
			Model:         "not_applicable",
		},
		CorpusEligible: false,
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal command result: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("unmarshal command result: %v", err)
	}
	for _, key := range []string{
		"object_type",
		"schema_version",
		"metadata",
		"command",
		"status",
		"mode",
		"warnings",
		"errors",
		"artifacts",
		"duration_ms",
		"redaction_status",
		"finding_kind",
		"proof_strength",
		"action_boundary_status",
		"security_relevance",
		"fix_target",
		"surface_fingerprint",
		"eval_mode",
		"provider_metadata",
		"corpus_eligible",
	} {
		if _, ok := decoded[key]; !ok {
			t.Fatalf("missing envelope key %s", key)
		}
	}
	modelProviderMetadata, ok := decoded["provider_metadata"].(map[string]any)
	if !ok {
		t.Fatalf("provider_metadata = %T, want object", decoded["provider_metadata"])
	}
	if modelProviderMetadata["semantic_role"] != "model_provider" {
		t.Fatalf("provider_metadata.semantic_role = %v, want model_provider", modelProviderMetadata["semantic_role"])
	}
}
