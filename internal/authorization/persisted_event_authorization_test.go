package authorization

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestValidateSnapshotFromPersistedInstallationAcceptsCanonicalContracts(t *testing.T) {
	now := time.Date(2026, 7, 25, 17, 30, 0, 0, time.UTC)
	snapshot, err := ValidateSnapshotFromPersistedInstallation(
		readCanonicalInstallationFixture(t),
		readCanonicalEventAuthorizationFixture(t),
		canonicalAuthorizationExpectations(),
		now,
	)
	if err != nil {
		t.Fatalf("canonical persisted authority contracts rejected: %v", err)
	}
	if snapshot.InstallationID != "installation.acme.payments_node" ||
		snapshot.InstallationDigest != "sha256:cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc" {
		t.Fatalf("installation binding not preserved: %#v", snapshot)
	}
	if snapshot.EventID != "event.payments.node_v5.0001" || snapshot.PlanID != "migration_plan.acme.payments_node_v5.0001" ||
		snapshot.ExecutionManifestID != "execution_manifest.acme.payments_node_v5.0001" {
		t.Fatalf("event/plan/execution bindings not preserved: %#v", snapshot)
	}
	if snapshot.RepositoryID != "github.acme.checkout-service" || snapshot.PackageRoot != "." ||
		snapshot.BaseCommit != "1234567890abcdef1234567890abcdef12345678" {
		t.Fatalf("repository binding not preserved: %#v", snapshot)
	}
	if snapshot.PathPolicyDigest != canonicalAuthorizationExpectations().Policies.PathPolicy ||
		snapshot.VerificationConfigurationDigest != canonicalAuthorizationExpectations().VerificationConfigurationDigest ||
		snapshot.CredentialIssuancePolicyDigest != canonicalAuthorizationExpectations().CredentialIssuancePolicyDigest {
		t.Fatalf("policy bindings not preserved: %#v", snapshot)
	}
	if len(snapshot.Scope.Commands) != 1 || snapshot.Scope.Commands[0] != "command.npm_test" {
		t.Fatalf("command selection not preserved: %#v", snapshot.Scope.Commands)
	}
	if strings.Join(snapshot.Scope.ExcludedPaths, ",") != ".git,node_modules" {
		t.Fatalf("installed exclusions were not retained: %#v", snapshot.Scope.ExcludedPaths)
	}
	if strings.Join(snapshot.Scope.ProviderReportingFields, ",") != "event_id,status,status_provenance,updated_at,verification_label,draft_pr_open" {
		t.Fatalf("provider reporting field ceiling not preserved: %#v", snapshot.Scope.ProviderReportingFields)
	}
}

func TestValidateSnapshotFromPersistedInstallationRejectsEventShapeDrift(t *testing.T) {
	value := readEventAuthorizationObject(t)
	value["plan_binding"] = "sha256:eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
	value["plan_digest"] = "sha256:eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"

	_, err := ValidateSnapshotFromPersistedInstallation(
		readCanonicalInstallationFixture(t), mustMarshalEventAuthorization(t, value),
		canonicalAuthorizationExpectations(), time.Date(2026, 7, 25, 17, 30, 0, 0, time.UTC),
	)
	if err == nil {
		t.Fatal("flattened Event Authorization shape must fail closed")
	}
}

func TestValidateSnapshotFromPersistedInstallationRejectsBindingSubstitution(t *testing.T) {
	tests := map[string]func(map[string]any){
		"installation": func(value map[string]any) {
			value["installation_binding"].(map[string]any)["artifact_digest"] = testContractDigest("0")
		},
		"event": func(value map[string]any) {
			value["event_binding"].(map[string]any)["artifact_digest"] = testContractDigest("1")
		},
		"migration pack": func(value map[string]any) {
			value["migration_pack_binding"].(map[string]any)["artifact_digest"] = testContractDigest("2")
		},
		"repository": func(value map[string]any) {
			value["repository_binding"].(map[string]any)["base_commit"] = strings.Repeat("a", 40)
		},
		"plan": func(value map[string]any) {
			value["plan_binding"].(map[string]any)["artifact_digest"] = testContractDigest("3")
		},
		"execution manifest": func(value map[string]any) {
			value["execution_manifest_binding"].(map[string]any)["artifact_digest"] = testContractDigest("4")
		},
		"policy digest": func(value map[string]any) {
			value["scope"].(map[string]any)["budget_policy_digest"] = testContractDigest("5")
		},
		"verification policy": func(value map[string]any) {
			value["verification_requirement"].(map[string]any)["configuration_digest"] = testContractDigest("6")
		},
		"credential policy": func(value map[string]any) {
			value["credential_issuance"].(map[string]any)["policy_digest"] = testContractDigest("7")
		},
		"Agent Runner": func(value map[string]any) {
			value["agent_execution_policy"].(map[string]any)["adapter_id"] = "claude_code"
		},
		"Agent Runner route digest": func(value map[string]any) {
			value["agent_execution_policy"].(map[string]any)["agent_route_digest"] = testContractDigest("9")
		},
	}
	for name, mutate := range tests {
		name, mutate := name, mutate
		t.Run(name, func(t *testing.T) {
			value := readEventAuthorizationObject(t)
			mutate(value)
			_, err := ValidateSnapshotFromPersistedInstallation(
				readCanonicalInstallationFixture(t), mustMarshalEventAuthorization(t, value),
				canonicalAuthorizationExpectations(), time.Date(2026, 7, 25, 17, 30, 0, 0, time.UTC),
			)
			if err == nil {
				t.Fatalf("schema-valid %s substitution must fail closed", name)
			}
		})
	}
}

func TestValidateSnapshotFromPersistedInstallationRejectsScopeWidening(t *testing.T) {
	tests := map[string]func(map[string]any){
		"path": func(value map[string]any) {
			scope := value["scope"].(map[string]any)
			scope["readable_paths"] = append(scope["readable_paths"].([]any), "secrets")
		},
		"capability": func(value map[string]any) {
			value["selected_capabilities"] = append(value["selected_capabilities"].([]any), "model_network")
		},
	}
	for name, mutate := range tests {
		name, mutate := name, mutate
		t.Run(name, func(t *testing.T) {
			value := readEventAuthorizationObject(t)
			mutate(value)
			_, err := ValidateSnapshotFromPersistedInstallation(
				readCanonicalInstallationFixture(t), mustMarshalEventAuthorization(t, value),
				canonicalAuthorizationExpectations(), time.Date(2026, 7, 25, 17, 30, 0, 0, time.UTC),
			)
			if err == nil {
				t.Fatalf("schema-valid %s widening must fail closed", name)
			}
		})
	}
}

func TestValidateSnapshotFromPersistedInstallationRejectsResolvedPolicyWidening(t *testing.T) {
	tests := map[string]func(*ExpectedAuthorizationBindings){
		"dropped exclusion": func(value *ExpectedAuthorizationBindings) {
			value.ExcludedPaths = nil
		},
		"budget": func(value *ExpectedAuthorizationBindings) {
			value.Budgets.MaxChangedFiles++
		},
		"provider reporting field": func(value *ExpectedAuthorizationBindings) {
			value.ProviderReportingFields = append(value.ProviderReportingFields, "merged")
		},
	}
	for name, mutate := range tests {
		name, mutate := name, mutate
		t.Run(name, func(t *testing.T) {
			expected := canonicalAuthorizationExpectations()
			mutate(&expected)
			_, err := ValidateSnapshotFromPersistedInstallation(
				readCanonicalInstallationFixture(t), readCanonicalEventAuthorizationFixture(t),
				expected, time.Date(2026, 7, 25, 17, 30, 0, 0, time.UTC),
			)
			if err == nil {
				t.Fatalf("resolved %s policy widening must fail closed", name)
			}
		})
	}
}

func TestValidateSnapshotFromPersistedInstallationRequiresIndependentExpectations(t *testing.T) {
	_, err := ValidateSnapshotFromPersistedInstallation(
		readCanonicalInstallationFixture(t), readCanonicalEventAuthorizationFixture(t),
		ExpectedAuthorizationBindings{}, time.Date(2026, 7, 25, 17, 30, 0, 0, time.UTC),
	)
	if err == nil || !strings.Contains(err.Error(), "expected ") || !strings.Contains(err.Error(), " required") {
		t.Fatalf("missing independent expectations error = %v", err)
	}
}

func canonicalAuthorizationExpectations() ExpectedAuthorizationBindings {
	return ExpectedAuthorizationBindings{
		Event:         ArtifactBinding{"event.payments.node_v5.0001", testContractDigest("b")},
		MigrationPack: ArtifactBinding{"migration_pack.payments.2027-01", testContractDigest("a")},
		Repository: RepositoryBinding{
			RepositoryID: "github.acme.checkout-service", PackageRoot: ".",
			BaseCommit: "1234567890abcdef1234567890abcdef12345678",
		},
		Plan:              ArtifactBinding{"migration_plan.acme.payments_node_v5.0001", testContractDigest("e")},
		ExecutionManifest: ArtifactBinding{"execution_manifest.acme.payments_node_v5.0001", testContractDigest("f")},
		Policies: PolicyDigestBindings{
			PathPolicy: testContractDigest("1"), CommandPolicy: testContractDigest("2"),
			ModelDataPolicy: testContractDigest("3"), BudgetPolicy: testContractDigest("4"),
			GitHubPolicy: testContractDigest("5"), ProviderReportingPolicy: testContractDigest("6"),
		},
		VerificationConfigurationDigest: testContractDigest("7"),
		CredentialIssuancePolicyDigest:  testContractDigest("8"),
		AgentRouteDigest:                "sha256:" + strings.Repeat("ab", 32),
		ExcludedPaths:                   []string{".git", "node_modules"},
		ProviderReportingFields: []string{
			"event_id", "status", "status_provenance", "updated_at", "verification_label", "draft_pr_open",
		},
		Budgets: Budgets{
			MaxChangedFiles: 12, MaxDiffLines: 500, MaxDiffBytes: 64000, MaxTurns: 20,
			MaxAttempts: 2, MaxTokens: 50000, MaxCostCents: 2500, MaxDurationSeconds: 1800,
		},
	}
}

func readCanonicalEventAuthorizationFixture(t *testing.T) []byte {
	t.Helper()
	path := filepath.Join("..", "..", "tests", "fixtures", "contracts", "event-authorization", "valid.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read canonical Event Authorization fixture: %v", err)
	}
	return data
}

func readEventAuthorizationObject(t *testing.T) map[string]any {
	t.Helper()
	var value map[string]any
	if err := json.Unmarshal(readCanonicalEventAuthorizationFixture(t), &value); err != nil {
		t.Fatalf("decode canonical Event Authorization fixture: %v", err)
	}
	return value
}

func mustMarshalEventAuthorization(t *testing.T, value map[string]any) []byte {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("encode Event Authorization fixture: %v", err)
	}
	return data
}

func testContractDigest(character string) string {
	return "sha256:" + strings.Repeat(character, 64)
}
