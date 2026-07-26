package status

import "testing"

func TestValidateProviderProjectionAcceptsObservedEvidenceBoundState(t *testing.T) {
	projection := ProviderProjection{
		EventID: "event-001", EventDigest: "sha256:event", RunID: "run-001",
		InstallationDigest: "sha256:installation", AuthorizationDigest: "sha256:authorization",
		PlanDigest: "sha256:plan", CandidateDigest: "sha256:candidate",
		VerificationDigest: "sha256:verification", EvidenceKinds: []string{"verification_outcome"},
		State: "verified", Provenance: "observed", Consented: true,
	}
	if err := ValidateProviderProjection(projection); err != nil {
		t.Fatalf("valid projection rejected: %v", err)
	}
}

func TestValidateProviderProjectionRejectsSilenceAndMergeInference(t *testing.T) {
	tests := map[string]ProviderProjection{
		"silence as affected": {
			EventID: "event-001", EventDigest: "sha256:event", RunID: "run-001",
			InstallationDigest: "sha256:installation", AuthorizationDigest: "sha256:authorization",
			State: "affected", Provenance: "unknown", Consented: true, EvidenceKinds: []string{"impact_outcome"},
		},
		"not applicable without evidence": {
			EventID: "event-001", EventDigest: "sha256:event", RunID: "run-001",
			InstallationDigest: "sha256:installation", AuthorizationDigest: "sha256:authorization",
			State: "not_applicable", Provenance: "consumer_reported", Consented: true,
		},
		"merge inferred as retired": {
			EventID: "event-001", EventDigest: "sha256:event", RunID: "run-001",
			InstallationDigest: "sha256:installation", AuthorizationDigest: "sha256:authorization",
			PlanDigest: "sha256:plan", CandidateDigest: "sha256:candidate",
			VerificationDigest: "sha256:verification", DeliveryDigest: "sha256:delivery",
			State: "retired", Provenance: "observed", Consented: true, EvidenceKinds: []string{"merge_outcome"},
		},
		"no consent": {
			EventID: "event-001", EventDigest: "sha256:event", RunID: "run-001",
			InstallationDigest: "sha256:installation", AuthorizationDigest: "sha256:authorization",
			PlanDigest: "sha256:plan", CandidateDigest: "sha256:candidate",
			VerificationDigest: "sha256:verification", EvidenceKinds: []string{"verification_outcome"},
			State: "verified", Provenance: "observed", Consented: false,
		},
	}
	for name, projection := range tests {
		t.Run(name, func(t *testing.T) {
			if err := ValidateProviderProjection(projection); err == nil {
				t.Fatal("expected dishonest provider projection to fail")
			}
		})
	}
}

func TestValidateProviderProjectionRequiresStatusSpecificEvidence(t *testing.T) {
	tests := []struct {
		state        string
		evidenceKind string
	}{
		{"received", "event_receipt"},
		{"not_applicable", "explicit_not_applicable"},
		{"affected", "impact_outcome"},
		{"needs_input", "consumer_input_request"},
		{"candidate_ready", "candidate_outcome"},
		{"verified", "verification_outcome"},
		{"draft_pr_open", "draft_pr_outcome"},
		{"accepted", "consumer_acceptance"},
		{"merged", "merge_outcome"},
		{"retired", "retirement_confirmation"},
	}
	for _, test := range tests {
		t.Run(test.state, func(t *testing.T) {
			projection := projectionForState(test.state, "merge_outcome")
			if test.evidenceKind != "merge_outcome" {
				if err := ValidateProviderProjection(projection); err == nil {
					t.Fatalf("%s accepted evidence kind merge_outcome, want %s", test.state, test.evidenceKind)
				}
			}
			projection.EvidenceKinds = []string{test.evidenceKind}
			if err := ValidateProviderProjection(projection); err != nil {
				t.Fatalf("%s rejected exact evidence kind %s: %v", test.state, test.evidenceKind, err)
			}
		})
	}
}

func TestValidateProviderProjectionKeepsEarlyBindingsNullable(t *testing.T) {
	evidenceKinds := map[string]string{
		"received": "event_receipt", "not_applicable": "explicit_not_applicable", "affected": "impact_outcome",
	}
	for _, state := range []string{"unknown", "received", "not_applicable", "affected"} {
		t.Run(state, func(t *testing.T) {
			projection := ProviderProjection{
				EventID: "event-001", EventDigest: "sha256:event", RunID: "run-001",
				InstallationDigest: "sha256:installation", AuthorizationDigest: "sha256:authorization",
				State: state, Consented: true,
			}
			if state == "unknown" {
				projection.Provenance = "unknown"
			} else {
				projection.Provenance = "observed"
				projection.EvidenceKinds = []string{evidenceKinds[state]}
			}
			if err := ValidateProviderProjection(projection); err != nil {
				t.Fatalf("early state %s rejected nullable later-artifact bindings: %v", state, err)
			}
		})
	}
}

func TestValidateProviderProjectionRejectsMissingExactBinding(t *testing.T) {
	tests := map[string]func(*ProviderProjection){
		"run":           func(value *ProviderProjection) { value.RunID = "" },
		"installation":  func(value *ProviderProjection) { value.InstallationDigest = "" },
		"authorization": func(value *ProviderProjection) { value.AuthorizationDigest = "" },
		"plan":          func(value *ProviderProjection) { value.PlanDigest = "" },
		"candidate":     func(value *ProviderProjection) { value.CandidateDigest = "" },
		"verification":  func(value *ProviderProjection) { value.VerificationDigest = "" },
		"delivery":      func(value *ProviderProjection) { value.DeliveryDigest = "" },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			projection := projectionForState("retired", "retirement_confirmation")
			mutate(&projection)
			if err := ValidateProviderProjection(projection); err == nil {
				t.Fatalf("retired projection without exact %s binding unexpectedly validated", name)
			}
		})
	}
}

func projectionForState(state, evidenceKind string) ProviderProjection {
	projection := ProviderProjection{
		EventID: "event-001", EventDigest: "sha256:event", RunID: "run-001",
		InstallationDigest: "sha256:installation", AuthorizationDigest: "sha256:authorization",
		State: state, Provenance: "observed", Consented: true, EvidenceKinds: []string{evidenceKind},
	}
	switch state {
	case "candidate_ready":
		projection.PlanDigest = "sha256:plan"
		projection.CandidateDigest = "sha256:candidate"
	case "verified":
		projection.PlanDigest = "sha256:plan"
		projection.CandidateDigest = "sha256:candidate"
		projection.VerificationDigest = "sha256:verification"
	case "draft_pr_open", "accepted", "merged", "retired":
		projection.PlanDigest = "sha256:plan"
		projection.CandidateDigest = "sha256:candidate"
		projection.VerificationDigest = "sha256:verification"
		projection.DeliveryDigest = "sha256:delivery"
	}
	return projection
}

func TestMigrationAxesRemainIndependent(t *testing.T) {
	axes := MigrationAxes{
		Impact: "affected_supported", Route: "agent_assisted", Candidate: "candidate_generated",
		Verification: "failed", Delivery: "draft_pr_open",
	}
	if err := ValidateMigrationAxes(axes); err != nil {
		t.Fatalf("independent axes rejected: %v", err)
	}
	if axes.Verification != "failed" {
		t.Fatalf("stronger delivery label rewrote verification to %q", axes.Verification)
	}
	axes.Verification = "green"
	if err := ValidateMigrationAxes(axes); err == nil {
		t.Fatal("unknown axis value must fail closed")
	}
}
