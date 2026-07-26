package redaction

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestValidateProviderPayloadAcceptsTypedProjection(t *testing.T) {
	payload := validProviderPayload()
	if _, err := BuildProviderPayload(payload, nil); err != nil {
		t.Fatalf("typed provider projection rejected: %v", err)
	}
}

func TestValidateProviderPayloadMatchesContractFixtures(t *testing.T) {
	tests := map[string]struct {
		fixture string
		valid   bool
	}{
		"valid":   {fixture: "valid.json", valid: true},
		"invalid": {fixture: "invalid.json", valid: false},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			payload := readProviderProjectionFixture(t, test.fixture)
			_, err := BuildProviderPayload(payload, nil)
			if test.valid && err != nil {
				t.Fatalf("valid contract fixture rejected: %v", err)
			}
			if !test.valid && err == nil {
				t.Fatal("invalid contract fixture accepted")
			}
		})
	}
}

func TestValidateProviderPayloadAcceptsFixedExternalConsentCeiling(t *testing.T) {
	payload := validProviderPayload()
	consentedFields := []string{
		"projection_id",
		"api_provider_id",
		"provider_change_event_id",
		"provider_change_contract_digest",
		"consumer_installation_pseudonym",
		"bindings",
		"status",
		"provenance",
		"evidence_refs",
		"observed_at",
	}
	if _, err := BuildProviderPayload(payload, consentedFields); err != nil {
		t.Fatalf("fixed consent ceiling rejected: %v", err)
	}
}

func TestBuildProviderPayloadEmitsOnlyDeclaredProviderFields(t *testing.T) {
	emitted, err := BuildProviderPayload(validProviderPayload(), nil)
	if err != nil {
		t.Fatalf("build provider payload: %v", err)
	}
	fields := decodeProviderPayloadFields(t, emitted)
	for _, field := range []string{
		"object_type", "schema_version", "consent", "provider_visible_fields",
		"interpretation_rules", "privacy", "artifact_digest",
	} {
		if _, present := fields[field]; present {
			t.Fatalf("consumer-private structural/control field %q crossed provider boundary", field)
		}
	}
	declared := validProviderPayload()["provider_visible_fields"].([]any)
	if len(fields) != len(declared) {
		t.Fatalf("emitted %d fields, want %d declared fields", len(fields), len(declared))
	}
	for _, rawField := range declared {
		field := rawField.(string)
		if _, present := fields[field]; !present {
			t.Fatalf("declared provider-visible field %q was not emitted", field)
		}
	}
}

func TestBuildProviderPayloadNarrowsToExternalConsentSubset(t *testing.T) {
	emitted, err := BuildProviderPayload(validProviderPayload(), []string{"status", "observed_at"})
	if err != nil {
		t.Fatalf("build narrowed provider payload: %v", err)
	}
	fields := decodeProviderPayloadFields(t, emitted)
	if len(fields) != 2 {
		t.Fatalf("emitted %d fields, want exact two-field consent subset", len(fields))
	}
	if string(fields["status"]) != `"merged"` || string(fields["observed_at"]) != `"2026-07-26T14:31:00Z"` {
		t.Fatalf("narrowed provider payload has unexpected values: %s", emitted)
	}
}

func TestBuildProviderPayloadRejectsExternalConsentWidening(t *testing.T) {
	payload := validProviderPayload()
	payload["provider_visible_fields"] = []any{"status"}
	if _, err := BuildProviderPayload(payload, []string{"status", "observed_at"}); err == nil {
		t.Fatal("expected external consent that widens declared provider fields to fail closed")
	}
}

func TestBuildProviderPayloadRejectsStructuralOrControlMetadataConsent(t *testing.T) {
	for _, field := range []string{
		"object_type", "schema_version", "consent", "provider_visible_fields",
		"interpretation_rules", "privacy", "artifact_digest",
	} {
		t.Run(field, func(t *testing.T) {
			if _, err := BuildProviderPayload(validProviderPayload(), []string{"status", field}); err == nil {
				t.Fatalf("expected consent attempt for consumer-private field %q to fail closed", field)
			}
		})
	}
}

func TestBuildProviderPayloadRejectsEmptyOrDuplicateExternalConsent(t *testing.T) {
	for name, consent := range map[string][]string{
		"empty":     {},
		"duplicate": {"status", "status"},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := BuildProviderPayload(validProviderPayload(), consent); err == nil {
				t.Fatal("expected invalid external consent to fail closed")
			}
		})
	}
}

func TestValidateProviderPayloadRejectsArbitraryConsentedExtension(t *testing.T) {
	payload := cloneWith(validProviderPayload(), "notes", "const apiKey = process.env.SECRET")
	if _, err := BuildProviderPayload(payload, []string{"notes"}); err == nil {
		t.Fatal("expected arbitrary consented extension to fail closed")
	}
}

func TestValidateProviderPayloadRejectsUnknownNestedExtension(t *testing.T) {
	payload := validProviderPayload()
	payload["evidence_refs"] = []any{map[string]any{
		"evidence_kind":   "merge_outcome",
		"evidence_digest": testDigest("d"),
		"provenance":      "consumer_reported",
		"details":         "raw consumer implementation",
	}}
	if _, err := BuildProviderPayload(payload, nil); err == nil {
		t.Fatal("expected unknown nested extension to fail closed")
	}
}

func TestValidateProviderPayloadRejectsMalformedNestedShapes(t *testing.T) {
	tests := map[string]func(map[string]any){
		"evidence ref missing provenance": func(payload map[string]any) {
			payload["evidence_refs"] = []any{map[string]any{
				"evidence_kind": "merge_outcome", "evidence_digest": testDigest("d"),
			}}
		},
		"consent missing event binding": func(payload map[string]any) {
			payload["consent"] = map[string]any{
				"consumer_consent_id": "consent-001", "consumer_consent_digest": testDigest("c"),
				"fields_enumerated": true, "active": true, "checked_at": "2026-07-26T14:30:00Z",
			}
		},
		"privacy has unsafe constant": func(payload map[string]any) {
			payload["privacy"].(map[string]any)["raw_consumer_code_included"] = true
		},
		"wrong field type": func(payload map[string]any) {
			payload["provider_visible_fields"] = "status"
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			payload := validProviderPayload()
			mutate(payload)
			if _, err := BuildProviderPayload(payload, nil); err == nil {
				t.Fatal("expected malformed typed projection to fail closed")
			}
		})
	}
}

func TestValidateProviderPayloadRejectsPrivateMaterialRecursively(t *testing.T) {
	tests := map[string]func(map[string]any){
		"raw diff field": func(payload map[string]any) {
			payload["privacy"].(map[string]any)["raw_diff"] = "-old +new"
		},
		"credential field": func(payload map[string]any) {
			payload["consent"].(map[string]any)["credential"] = "secret"
		},
		"agent session field": func(payload map[string]any) {
			payload["evidence_refs"] = []any{map[string]any{
				"evidence_kind": "merge_outcome", "evidence_digest": testDigest("d"),
				"provenance": "consumer_reported", "agent_session": map[string]any{"id": "session-1"},
			}}
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			payload := validProviderPayload()
			mutate(payload)
			if _, err := BuildProviderPayload(payload, nil); err == nil {
				t.Fatal("expected consumer-private material to fail closed")
			}
		})
	}
}

func TestValidateProviderPayloadRejectsCommonSecretPatternsInAllowedFields(t *testing.T) {
	tests := map[string]string{
		// Build credential-shaped test inputs at runtime so push protection does
		// not mistake inert negative-test fixtures for usable credentials.
		"stripe live key": strings.Join([]string{"sk", "live", "51N1234567890abcdefghijklmnop"}, "_"),
		"github token":    strings.Join([]string{"ghp", "1234567890abcdefghijklmnopqrstuvwxyz"}, "_"),
		"private key":     strings.Join([]string{"-----BEGIN", "PRIVATE KEY-----"}, " "),
	}
	for name, secret := range tests {
		t.Run(name, func(t *testing.T) {
			payload := validProviderPayload()
			payload["projection_id"] = secret
			if _, err := BuildProviderPayload(payload, nil); err == nil {
				t.Fatal("expected common secret pattern to fail closed")
			}
		})
	}
}

func TestValidateProviderPayloadRejectsStatusEvidenceContradictions(t *testing.T) {
	tests := map[string]func(map[string]any){
		"unknown has observed provenance": func(payload map[string]any) {
			payload["status"] = "unknown"
			payload["provenance"] = "observed"
			payload["evidence_refs"] = []any{}
		},
		"unknown has evidence": func(payload map[string]any) {
			payload["status"] = "unknown"
			payload["provenance"] = "unknown"
		},
		"not applicable lacks explicit evidence": func(payload map[string]any) {
			payload["status"] = "not_applicable"
		},
		"merged implies retired": func(payload map[string]any) {
			payload["interpretation_rules"].(map[string]any)["merged_implies_retired"] = true
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			payload := validProviderPayload()
			mutate(payload)
			if _, err := BuildProviderPayload(payload, nil); err == nil {
				t.Fatal("expected contradictory status projection to fail closed")
			}
		})
	}
}

func validProviderPayload() map[string]any {
	return map[string]any{
		"object_type": "lumyn.provider_status_projection", "schema_version": "1.0",
		"projection_id": "projection-001", "api_provider_id": "provider-001",
		"provider_change_event_id": "event-001", "provider_change_contract_digest": testDigest("a"),
		"consumer_installation_pseudonym": "consumer-unit-7", "status": "merged", "provenance": "consumer_reported",
		"bindings": map[string]any{
			"run_id": "run-001", "consumer_installation_digest": testDigest("1"),
			"event_authorization_digest": testDigest("2"), "migration_plan_digest": testDigest("3"),
			"candidate_digest": testDigest("4"), "verification_digest": testDigest("5"),
			"delivery_digest": testDigest("6"),
		},
		"evidence_refs": []any{map[string]any{
			"evidence_kind": "merge_outcome", "evidence_digest": testDigest("d"), "provenance": "consumer_reported",
		}},
		"consent": map[string]any{
			"consumer_consent_id": "consent-001", "consumer_consent_digest": testDigest("c"),
			"event_bound": true, "fields_enumerated": true, "active": true,
			"checked_at": "2026-07-26T14:30:00Z",
		},
		"provider_visible_fields": []any{
			"projection_id", "api_provider_id", "provider_change_event_id",
			"provider_change_contract_digest", "consumer_installation_pseudonym",
			"bindings", "status", "provenance", "evidence_refs", "observed_at",
		},
		"interpretation_rules": map[string]any{
			"silence": "unknown", "merged_implies_retired": false,
			"not_applicable_requires_explicit_evidence": true, "unaffected_may_be_inferred": false,
		},
		"privacy": map[string]any{
			"classification": "provider_consented_projection", "redaction_policy_digest": testDigest("b"),
			"redaction_status": "passed", "safe_to_share_with_api_provider": true,
			"raw_consumer_code_included": false, "raw_diff_included": false,
			"raw_logs_or_traces_included": false, "raw_prompts_or_responses_included": false,
			"agent_session_data_included": false, "credentials_included": false,
		},
		"observed_at": "2026-07-26T14:31:00Z", "artifact_digest": testDigest("f"),
	}
}

func decodeProviderPayloadFields(t *testing.T, payload ProviderPayload) map[string]json.RawMessage {
	t.Helper()
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(payload, &fields); err != nil {
		t.Fatalf("decode emitted provider payload: %v", err)
	}
	return fields
}

func readProviderProjectionFixture(t *testing.T, name string) map[string]any {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve provider projection test path")
	}
	fixturePath := filepath.Join(filepath.Dir(filename), "..", "..", "tests", "fixtures", "contracts", "provider-status-projection", name)
	contents, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("read provider projection fixture: %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(contents, &payload); err != nil {
		t.Fatalf("decode provider projection fixture: %v", err)
	}
	return payload
}

func testDigest(character string) string {
	return "sha256:" + repeat(character, 64)
}

func repeat(value string, count int) string {
	result := ""
	for range count {
		result += value
	}
	return result
}

func cloneWith(source map[string]any, key string, value any) map[string]any {
	clone := make(map[string]any, len(source)+1)
	for name, item := range source {
		clone[name] = item
	}
	clone[key] = value
	return clone
}
