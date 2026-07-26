package evidence

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestValidateMigrationBindingsAcceptsOneExactChain(t *testing.T) {
	bindings := coherentBindings()
	if err := ValidateMigrationBindings(bindings); err != nil {
		t.Fatalf("coherent bindings rejected: %v", err)
	}
}

func TestValidateMigrationBindingsRejectsCrossArtifactDigestMismatch(t *testing.T) {
	tests := map[string]func(*MigrationBindings){
		"authorization to installation": func(value *MigrationBindings) { value.AuthorizationInstallationDigest = "sha256:other" },
		"event to pack":                 func(value *MigrationBindings) { value.EventContractDigest = "sha256:other" },
		"authorization to event":        func(value *MigrationBindings) { value.AuthorizationEventDigest = "sha256:other" },
		"plan to authorization":         func(value *MigrationBindings) { value.PlanAuthorizationDigest = "sha256:other" },
		"plan to impact":                func(value *MigrationBindings) { value.PlanImpactDigest = "sha256:other" },
		"candidate to plan":             func(value *MigrationBindings) { value.CandidatePlanDigest = "sha256:other" },
		"verification to candidate":     func(value *MigrationBindings) { value.VerificationCandidateDigest = "sha256:other" },
		"export to candidate":           func(value *MigrationBindings) { value.ExportCandidateDigest = "sha256:other" },
		"export to verification":        func(value *MigrationBindings) { value.ExportVerificationDigest = "sha256:other" },
		"projection to event":           func(value *MigrationBindings) { value.ProjectionEventDigest = "sha256:other" },
		"projection to run":             func(value *MigrationBindings) { value.ProjectionRunID = "run-other" },
		"projection to installation": func(value *MigrationBindings) {
			value.ProjectionInstallationDigest = "sha256:other"
		},
		"projection to authorization": func(value *MigrationBindings) {
			value.ProjectionAuthorizationDigest = "sha256:other"
		},
		"projection to plan": func(value *MigrationBindings) { value.ProjectionPlanDigest = "sha256:other" },
		"projection to candidate": func(value *MigrationBindings) {
			value.ProjectionCandidateDigest = "sha256:other"
		},
		"projection to verification": func(value *MigrationBindings) {
			value.ProjectionVerificationDigest = "sha256:other"
		},
		"projection to delivery": func(value *MigrationBindings) { value.ProjectionDeliveryDigest = "sha256:other" },
		"projection evidence to verification": func(value *MigrationBindings) {
			value.ProjectionEvidenceDigest = "sha256:other"
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			value := coherentBindings()
			mutate(&value)
			if err := ValidateMigrationBindings(value); err == nil {
				t.Fatal("expected digest mismatch to fail")
			}
		})
	}
}

func TestProjectionBindingInputsInvalidateEveryBoundArtifact(t *testing.T) {
	prior := ProjectionBindingInputs(coherentBindings())
	tests := map[string]func(*MigrationBindings){
		"run":           func(value *MigrationBindings) { value.RunID = "run-changed" },
		"installation":  func(value *MigrationBindings) { value.InstallationDigest = "sha256:changed" },
		"authorization": func(value *MigrationBindings) { value.AuthorizationDigest = "sha256:changed" },
		"plan":          func(value *MigrationBindings) { value.PlanDigest = "sha256:changed" },
		"candidate":     func(value *MigrationBindings) { value.CandidateDigest = "sha256:changed" },
		"verification":  func(value *MigrationBindings) { value.VerificationDigest = "sha256:changed" },
		"delivery":      func(value *MigrationBindings) { value.ExportDigest = "sha256:changed" },
	}
	for key, mutate := range tests {
		t.Run(key, func(t *testing.T) {
			bindings := coherentBindings()
			mutate(&bindings)
			if got, want := ChangedBindingInputs(prior, ProjectionBindingInputs(bindings)), []string{key}; !reflect.DeepEqual(got, want) {
				t.Fatalf("changed inputs = %#v, want %#v", got, want)
			}
		})
	}
}

func TestPersistedProviderProjectionRejectsCrossRunSubstitution(t *testing.T) {
	payload, err := os.ReadFile(filepath.Join("..", "..", "tests", "fixtures", "contracts", "provider-status-projection", "cross-run-substitution.json"))
	if err != nil {
		t.Fatalf("read persisted cross-run fixture: %v", err)
	}
	var fixture struct {
		ProjectionBindings struct {
			RunID               string `json:"run_id"`
			InstallationDigest  string `json:"consumer_installation_digest"`
			AuthorizationDigest string `json:"event_authorization_digest"`
			PlanDigest          string `json:"migration_plan_digest"`
			CandidateDigest     string `json:"candidate_digest"`
			VerificationDigest  string `json:"verification_digest"`
			DeliveryDigest      string `json:"delivery_digest"`
		} `json:"projection_bindings"`
	}
	if err := json.Unmarshal(payload, &fixture); err != nil {
		t.Fatalf("decode persisted cross-run fixture: %v", err)
	}
	bindings := coherentBindings()
	bindings.ProjectionRunID = fixture.ProjectionBindings.RunID
	bindings.ProjectionInstallationDigest = fixture.ProjectionBindings.InstallationDigest
	bindings.ProjectionAuthorizationDigest = fixture.ProjectionBindings.AuthorizationDigest
	bindings.ProjectionPlanDigest = fixture.ProjectionBindings.PlanDigest
	bindings.ProjectionCandidateDigest = fixture.ProjectionBindings.CandidateDigest
	bindings.ProjectionVerificationDigest = fixture.ProjectionBindings.VerificationDigest
	bindings.ProjectionDeliveryDigest = fixture.ProjectionBindings.DeliveryDigest
	if err := ValidateMigrationBindings(bindings); err == nil {
		t.Fatal("persisted evidence from another run unexpectedly satisfied the provider projection")
	}
}

func TestChangedBindingInputsInvalidatesEvidenceDeterministically(t *testing.T) {
	prior := map[string]string{
		"event": "sha256:event", "pack": "sha256:pack", "base": "sha256:base",
		"plan": "sha256:plan", "candidate": "sha256:candidate", "commands": "sha256:commands",
	}
	current := map[string]string{
		"event": "sha256:event", "pack": "sha256:pack-v2", "base": "sha256:base",
		"plan": "sha256:plan", "candidate": "sha256:candidate-v2", "commands": "sha256:commands",
	}
	if got, want := ChangedBindingInputs(prior, current), []string{"candidate", "pack"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("changed inputs = %#v, want %#v", got, want)
	}
}

func coherentBindings() MigrationBindings {
	return MigrationBindings{
		RunID: "run-001", ProjectionRunID: "run-001",
		InstallationDigest: "sha256:installation", AuthorizationInstallationDigest: "sha256:installation",
		AuthorizationDigest: "sha256:authorization", PlanAuthorizationDigest: "sha256:authorization",
		PackDigest: "sha256:pack", EventDigest: "sha256:event", EventContractDigest: "sha256:pack",
		AuthorizationEventDigest: "sha256:event", ImpactDigest: "sha256:impact",
		PlanImpactDigest: "sha256:impact", PlanDigest: "sha256:plan",
		CandidatePlanDigest: "sha256:plan", CandidateDigest: "sha256:candidate",
		VerificationCandidateDigest: "sha256:candidate", VerificationDigest: "sha256:verification",
		ExportCandidateDigest: "sha256:candidate", ExportVerificationDigest: "sha256:verification",
		ExportDigest: "sha256:delivery", ProjectionEventDigest: "sha256:event",
		ProjectionInstallationDigest:  "sha256:installation",
		ProjectionAuthorizationDigest: "sha256:authorization", ProjectionPlanDigest: "sha256:plan",
		ProjectionCandidateDigest: "sha256:candidate", ProjectionVerificationDigest: "sha256:verification",
		ProjectionDeliveryDigest: "sha256:delivery", ProjectionEvidenceDigest: "sha256:verification",
	}
}
