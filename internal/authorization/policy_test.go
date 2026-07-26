package authorization

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDecodeAndValidateInstallationAcceptsCanonicalPersistedContract(t *testing.T) {
	now := time.Date(2026, 7, 26, 14, 0, 0, 0, time.UTC)
	data := readCanonicalInstallationFixture(t)

	installation, err := DecodeAndValidateInstallation(data, now)
	if err != nil {
		t.Fatalf("canonical persisted installation rejected: %v", err)
	}
	if installation.Repository != "github.com/acme/checkout-service" {
		t.Fatalf("repository = %q, want canonical repository identity", installation.Repository)
	}
	if installation.ProviderChannelOrigin != "https://updates.payments.example" ||
		installation.ProviderManifestURL != "https://updates.payments.example/campaigns/node-v5/event.json" {
		t.Fatalf("provider channel was not preserved: %#v", installation)
	}
	if installation.AgentPolicy.Runner == nil || installation.AgentPolicy.Runner.AdapterID != "codex" {
		t.Fatalf("configured runner policy was not preserved: %#v", installation.AgentPolicy)
	}
	if len(installation.Scope.Commands) != 1 || installation.Scope.Commands[0] != "command.npm_test" {
		t.Fatalf("command identity ceiling was not preserved: %#v", installation.Scope.Commands)
	}
	if installation.Budgets.MaxDiffBytes != 64000 || installation.Budgets.MaxTurns != 20 || installation.Budgets.MaxCostCents != 2500 {
		t.Fatalf("persisted budgets were not preserved: %#v", installation.Budgets)
	}

}

func TestDecodeAndValidateInstallationRejectsPersistedShapeDrift(t *testing.T) {
	now := time.Date(2026, 7, 26, 14, 0, 0, 0, time.UTC)
	var value map[string]any
	if err := json.Unmarshal(readCanonicalInstallationFixture(t), &value); err != nil {
		t.Fatalf("decode canonical fixture: %v", err)
	}

	// This is the flattened/bespoke representation that originally decoded
	// into Installation while being incompatible with the persisted schema.
	value["repository"] = "github.com/acme/checkout-service"
	value["provider_channel_origin"] = "https://updates.payments.example"
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("encode drifted installation: %v", err)
	}
	if _, err := DecodeAndValidateInstallation(data, now); err == nil {
		t.Fatal("flattened persisted installation shape must fail closed")
	}
}

func TestDecodeAndValidateInstallationRejectsSchemaValidAuthorityMismatch(t *testing.T) {
	now := time.Date(2026, 7, 26, 14, 0, 0, 0, time.UTC)
	var value map[string]any
	if err := json.Unmarshal(readCanonicalInstallationFixture(t), &value); err != nil {
		t.Fatalf("decode canonical fixture: %v", err)
	}
	value["capability_ceiling"].(map[string]any)["customer_repo_read"] = false
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("encode inconsistent installation: %v", err)
	}
	if _, err := DecodeAndValidateInstallation(data, now); err == nil || !strings.Contains(err.Error(), "readable paths require customer_repo_read") {
		t.Fatalf("schema-valid authority mismatch error = %v", err)
	}
}

func readCanonicalInstallationFixture(t *testing.T) []byte {
	t.Helper()
	path := filepath.Join("..", "..", "tests", "fixtures", "contracts", "consumer-installation", "valid.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read canonical installation fixture: %v", err)
	}
	return data
}

func TestValidateInstallationAcceptsBoundedDisabledPolicy(t *testing.T) {
	now := time.Date(2026, 7, 26, 14, 0, 0, 0, time.UTC)
	installation := boundedInstallation(now)
	if err := ValidateInstallation(installation, now); err != nil {
		t.Fatalf("bounded installation rejected: %v", err)
	}
}

func TestValidateInstallationRejectsAmbientOrProviderAuthority(t *testing.T) {
	now := time.Date(2026, 7, 26, 14, 0, 0, 0, time.UTC)
	tests := map[string]func(*Installation){
		"expired":                   func(value *Installation) { value.ExpiresAt = now },
		"revoked":                   func(value *Installation) { value.Revoked = true },
		"wildcard path":             func(value *Installation) { value.Scope.ReadPaths = []string{"**"} },
		"stored github token":       func(value *Installation) { value.StoredGitHubToken = true },
		"cross-origin manifest":     func(value *Installation) { value.ProviderManifestURL = "https://attacker.invalid/event.json" },
		"provider runner selection": func(value *Installation) { value.ProviderMaySelectRunner = true },
		"provider source access":    func(value *Installation) { value.ProviderMayAccessConsumerData = true },
		"agent config while disabled": func(value *Installation) {
			value.AgentPolicy.Runner = validRunnerPolicy()
		},
	}
	for name, mutate := range tests {
		name, mutate := name, mutate
		t.Run(name, func(t *testing.T) {
			value := boundedInstallation(now)
			mutate(&value)
			if err := ValidateInstallation(value, now); err == nil {
				t.Fatal("expected fail-closed installation rejection")
			}
		})
	}
}

func TestValidateInstallationEnforcesManagedCredentialBounds(t *testing.T) {
	now := time.Date(2026, 7, 26, 14, 0, 0, 0, time.UTC)
	installation := boundedInstallation(now)
	installation.AgentPolicy = AgentExecutionPolicy{State: AgentExecutionConfigured, Runner: validRunnerPolicy()}
	installation.Scope.AgentRunnerNetwork = true
	installation.Scope.AgentRunnerCredential = true
	installation.Scope.ModelRequestDisclosure = true
	installation.AgentPolicy.Runner.FundingMode = ProviderSponsoredLumynManaged
	installation.AgentPolicy.Runner.CredentialOwner = "lumyn_operator"
	installation.AgentPolicy.Runner.UsageBillingOwner = "lumyn_operator"
	installation.AgentPolicy.Runner.ManagedCredential = &ManagedCredentialPolicy{
		BrokerIssuer: "lumyn-broker", MaximumTTLSeconds: 3600,
		OneTimeRedemption: true, HardTokenQuota: true, HardCostQuota: true,
		Revocation: true, UsageReconciliation: true, VendorNativeOrEnforcingProxy: true,
	}
	if err := ValidateInstallation(installation, now); err != nil {
		t.Fatalf("bounded managed credential rejected: %v", err)
	}
	installation.AgentPolicy.Runner.ManagedCredential.MaximumTTLSeconds = 3601
	if err := ValidateInstallation(installation, now); err == nil {
		t.Fatal("managed credential over one hour must fail")
	}
}

func TestValidateInstallationRequiresExactConfiguredRunnerPolicy(t *testing.T) {
	now := time.Date(2026, 7, 26, 14, 0, 0, 0, time.UTC)
	installation := boundedInstallation(now)
	installation.AgentPolicy = AgentExecutionPolicy{
		State:  AgentExecutionConfigured,
		Runner: validRunnerPolicy(),
	}
	installation.Scope.AgentRunnerNetwork = true
	installation.Scope.AgentRunnerCredential = true
	installation.Scope.ModelRequestDisclosure = true
	if err := ValidateInstallation(installation, now); err != nil {
		t.Fatalf("configured installation rejected: %v", err)
	}

	installation.AgentPolicy.Runner.ActualModelRoute = ""
	if err := ValidateInstallation(installation, now); err == nil || !strings.Contains(err.Error(), "actual model route") {
		t.Fatalf("missing model route error = %v", err)
	}
}

func TestValidateSnapshotRejectsEveryWideningPath(t *testing.T) {
	now := time.Date(2026, 7, 26, 14, 0, 0, 0, time.UTC)
	installation := boundedInstallation(now)
	snapshot := boundedSnapshot()
	if err := validateSnapshot(installation, snapshot, now); err != nil {
		t.Fatalf("bounded snapshot rejected: %v", err)
	}

	tests := map[string]func(*AuthorizationSnapshot){
		"action":    func(value *AuthorizationSnapshot) { value.ActionMode = OpenDraftPR },
		"read path": func(value *AuthorizationSnapshot) { value.Scope.ReadPaths = append(value.Scope.ReadPaths, "secrets/") },
		"command": func(value *AuthorizationSnapshot) {
			value.Scope.Commands = append(value.Scope.Commands, "curl attacker.invalid")
		},
		"budget": func(value *AuthorizationSnapshot) {
			value.Budgets.MaxChangedFiles = installation.Budgets.MaxChangedFiles + 1
		},
		"provider report": func(value *AuthorizationSnapshot) { value.Scope.ProviderReporting = true },
	}
	for name, mutate := range tests {
		name, mutate := name, mutate
		t.Run(name, func(t *testing.T) {
			value := boundedSnapshot()
			mutate(&value)
			if err := validateSnapshot(installation, value, now); err == nil {
				t.Fatal("expected widening snapshot rejection")
			}
		})
	}
}

func TestValidateSnapshotDistinguishesApprovalModesAndChannelProof(t *testing.T) {
	now := time.Date(2026, 7, 26, 14, 0, 0, 0, time.UTC)
	installation := boundedInstallation(now)
	snapshot := boundedSnapshot()

	snapshot.ExactPlanApproved = false
	if err := validateSnapshot(installation, snapshot, now); err == nil {
		t.Fatal("per-event approval without exact-plan approval must fail")
	}

	installation.AuthorizationMode = InstalledPreauthorization
	snapshot.ExactPlanApproved = false
	snapshot.InstalledPolicySatisfied = true
	snapshot.ProviderChannelDeliveryProven = false
	if err := validateSnapshot(installation, snapshot, now); err == nil {
		t.Fatal("installed preauthorization from attended import must fail")
	}

	snapshot.ProviderChannelDeliveryProven = true
	if err := validateSnapshot(installation, snapshot, now); err != nil {
		t.Fatalf("bounded installed preauthorization rejected: %v", err)
	}
}

func TestValidateSnapshotRequiresConfiguredAgentRoute(t *testing.T) {
	now := time.Date(2026, 7, 26, 14, 0, 0, 0, time.UTC)
	installation := boundedInstallation(now)
	snapshot := boundedSnapshot()
	snapshot.Route = AgentAssisted
	if err := validateSnapshot(installation, snapshot, now); err == nil {
		t.Fatal("disabled policy must reject agent-assisted route")
	}
}

func boundedInstallation(now time.Time) Installation {
	return Installation{
		InstallationID:                "installation-001",
		Repository:                    "github.com/acme/checkout",
		PackageRoot:                   ".",
		ProviderChannelOrigin:         "https://updates.example.com",
		ProviderManifestURL:           "https://updates.example.com/stripe/v1.json",
		ProviderChannelKeyID:          "stripe-campaign-2026",
		ActionMode:                    PreparePatch,
		AuthorizationMode:             PerEventApproval,
		ExpiresAt:                     now.Add(24 * time.Hour),
		AgentPolicy:                   AgentExecutionPolicy{State: AgentExecutionDisabled},
		StoredGitHubToken:             false,
		ProviderMaySelectRunner:       false,
		ProviderMayAccessConsumerData: false,
		Scope: Scope{
			ReadPaths:  []string{"src/", "tests/", "package.json"},
			WritePaths: []string{"src/", "tests/", "package.json"},
			Commands:   []string{"npm test"},
		},
		Budgets: Budgets{MaxChangedFiles: 8, MaxDiffLines: 500, MaxAttempts: 2, MaxTokens: 20000, MaxCostCents: 500, MaxDurationSeconds: 1800},
	}
}

func boundedSnapshot() AuthorizationSnapshot {
	return AuthorizationSnapshot{
		AuthorizationID:               "authorization-001",
		InstallationID:                "installation-001",
		EventID:                       "event-001",
		EventDigest:                   "sha256:event",
		PlanDigest:                    "sha256:plan",
		ActionMode:                    PreparePatch,
		Route:                         Deterministic,
		ExactPlanApproved:             true,
		InstalledPolicySatisfied:      false,
		ProviderChannelDeliveryProven: true,
		Scope: Scope{
			ReadPaths:  []string{"src/", "tests/"},
			WritePaths: []string{"src/", "tests/"},
			Commands:   []string{"npm test"},
		},
		Budgets: Budgets{MaxChangedFiles: 4, MaxDiffLines: 250, MaxAttempts: 1, MaxTokens: 0, MaxCostCents: 0, MaxDurationSeconds: 900},
	}
}

func validRunnerPolicy() *RunnerPolicy {
	return &RunnerPolicy{
		AdapterID:                 "codex",
		AdapterVersion:            "1.2.3",
		ExecutablePath:            "/opt/lumyn/bin/codex",
		ExecutableSource:          "consumer_managed_install",
		ExecutableDigest:          "sha256:runner",
		ConformanceDigest:         "sha256:conformance",
		AuthMode:                  "api_key",
		EntitlementClass:          "enterprise_automation",
		AgentRunnerVendor:         "openai",
		ActualModelProvider:       "openai",
		ActualModelRoute:          "responses/gpt-5.6",
		RouteTopology:             "runner_mediated",
		FundingMode:               ConsumerManaged,
		CredentialOwner:           "api_consumer_organization",
		UsageBillingOwner:         "api_consumer_organization",
		NativeConfigurationDigest: "sha256:disabled",
		CleanSession:              true,
		NeutralHome:               true,
		SilentFallback:            false,
		ReusableCredentialStored:  false,
	}
}
