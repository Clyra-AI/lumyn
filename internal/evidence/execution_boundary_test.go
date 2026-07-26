package evidence

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestExecutionManifestBoundaryFixtures(t *testing.T) {
	host := HostBoundary{
		HomePath:               "/Users/customer",
		OSCredentialStorePaths: []string{"/Users/customer/Library/Keychains"},
	}
	normal := readExecutionManifestFixture(t, "valid.json")
	if err := ValidateExecutionManifestBoundary(normal, host); err != nil {
		t.Fatalf("normal execution manifest rejected: %v", err)
	}
	sandbox := readExecutionManifestFixture(t, "valid-sandbox.json")
	if err := ValidateExecutionManifestBoundary(sandbox, host); err != nil {
		t.Fatalf("sandbox execution manifest rejected: %v", err)
	}
}

func TestExecutionManifestBoundaryRejectsHostHomeDisguisedAsRepository(t *testing.T) {
	value := readExecutionManifestFixture(t, "semantic-invalid-host-home-source.json")
	err := ValidateExecutionManifestBoundary(value, HostBoundary{HomePath: "/Users/customer"})
	if err == nil {
		t.Fatal("expected a source-class label to be unable to hide a host-home mount")
	}
}

func TestExecutionManifestBoundaryRejectsProgramOutsideRoot(t *testing.T) {
	value := readExecutionManifestFixture(t, "semantic-invalid-executable-root.json")
	err := ValidateExecutionManifestBoundary(value, HostBoundary{HomePath: "/Users/customer"})
	if err == nil {
		t.Fatal("expected a command outside its declared executable root to fail")
	}
}

func TestExecutionManifestBoundaryRejectsExecutableRootOutsideMount(t *testing.T) {
	value := readExecutionManifestFixture(t, "valid.json")
	value.ExecutableRoots[0].SourcePath = "/unmounted/toolchain/bin"
	err := ValidateExecutionManifestBoundary(value, HostBoundary{HomePath: "/Users/customer"})
	if err == nil {
		t.Fatal("expected an executable root outside its source mount to fail")
	}
}

func TestExecutionManifestBoundaryRejectsSandboxCrossBindingDrift(t *testing.T) {
	tests := map[string]func(*ExecutionManifestBoundary){
		"writable candidate mount": func(value *ExecutionManifestBoundary) {
			for index := range value.Mounts {
				if value.Mounts[index].MountID == value.SandboxEntrypointProfile.CandidateMountID {
					value.Mounts[index].Mode = "read_write"
				}
			}
		},
		"wrong credential": func(value *ExecutionManifestBoundary) {
			value.SandboxEntrypointProfile.SandboxCredentialGrantID = "credential.sandbox.other"
		},
		"wrong network operations": func(value *ExecutionManifestBoundary) {
			value.SandboxEntrypointProfile.NetworkOperations = []string{"payment_method.read"}
		},
		"entrypoint outside candidate": func(value *ExecutionManifestBoundary) {
			value.SandboxEntrypointProfile.EntrypointPath = "/tmp/entrypoint.js"
		},
		"neutral root overlaps host home": func(value *ExecutionManifestBoundary) {
			value.SandboxEntrypointProfile.NeutralHomeRoot = "/Users/customer"
		},
		"credential has extra scope": func(value *ExecutionManifestBoundary) {
			value.CredentialGrants[0].ScopeIDs = append(value.CredentialGrants[0].ScopeIDs, "sandbox_profile.other")
		},
		"duplicate grant identity": func(value *ExecutionManifestBoundary) {
			value.CredentialGrants[0].GrantID = value.NetworkGrants[0].GrantID
		},
		"missing resource quota": func(value *ExecutionManifestBoundary) {
			value.SandboxEntrypointProfile.ResourceQuotas.MaxOpenFiles = 0
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			value := readExecutionManifestFixture(t, "valid-sandbox.json")
			mutate(&value)
			if err := ValidateExecutionManifestBoundary(value, HostBoundary{HomePath: "/Users/customer"}); err == nil {
				t.Fatal("expected sandbox binding drift to fail")
			}
		})
	}
}

func TestMigrationVerificationBoundary(t *testing.T) {
	manifest := readExecutionManifestFixture(t, "valid-sandbox.json")
	profileDigest := manifest.SandboxEntrypointProfile.ProfileDigest
	value := MigrationVerificationBoundary{
		ConsumerExecutionManifestDigest: manifest.ManifestDigest,
		CandidateHead:                   manifest.Repository.CandidateHead,
		VerificationLabel:               "workflow_verified_sandbox",
		SandboxEntrypointProfileDigest:  profileDigest,
		ObservedEvidence: []SandboxVerificationEvidence{
			{EvidenceKind: "exact_head_workflow_outcome", CandidateHead: manifest.Repository.CandidateHead, SandboxEntrypointProfileDigest: profileDigest},
			{EvidenceKind: "sandbox_observation", CandidateHead: manifest.Repository.CandidateHead, SandboxEntrypointProfileDigest: profileDigest},
		},
	}
	if err := ValidateMigrationVerificationBoundary(value, manifest); err != nil {
		t.Fatalf("valid sandbox verification rejected: %v", err)
	}

	tests := map[string]func(*MigrationVerificationBoundary){
		"wrong manifest": func(value *MigrationVerificationBoundary) {
			value.ConsumerExecutionManifestDigest = "sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
		},
		"wrong profile": func(value *MigrationVerificationBoundary) {
			value.SandboxEntrypointProfileDigest = "sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
		},
		"missing observation": func(value *MigrationVerificationBoundary) {
			value.ObservedEvidence = value.ObservedEvidence[:1]
		},
		"stale evidence head": func(value *MigrationVerificationBoundary) {
			value.ObservedEvidence[1].CandidateHead = "ffffffffffffffffffffffffffffffffffffffff"
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			candidate := value
			candidate.ObservedEvidence = append([]SandboxVerificationEvidence(nil), value.ObservedEvidence...)
			mutate(&candidate)
			if err := ValidateMigrationVerificationBoundary(candidate, manifest); err == nil {
				t.Fatal("expected stale sandbox verification binding to fail")
			}
		})
	}
}

func TestNonSandboxVerificationRejectsSandboxEvidence(t *testing.T) {
	manifest := readExecutionManifestFixture(t, "valid.json")
	value := MigrationVerificationBoundary{
		ConsumerExecutionManifestDigest: manifest.ManifestDigest,
		CandidateHead:                   manifest.Repository.CandidateHead,
		VerificationLabel:               "repo_verified",
	}
	if err := ValidateMigrationVerificationBoundary(value, manifest); err != nil {
		t.Fatalf("ordinary verification rejected: %v", err)
	}
	value.ObservedEvidence = []SandboxVerificationEvidence{{EvidenceKind: "sandbox_observation"}}
	if err := ValidateMigrationVerificationBoundary(value, manifest); err == nil {
		t.Fatal("expected non-sandbox verification to reject sandbox evidence")
	}
}

func TestWorkflowMockVerificationRequiresCurrentHeadEvidence(t *testing.T) {
	manifest := readExecutionManifestFixture(t, "valid.json")
	value := MigrationVerificationBoundary{
		ConsumerExecutionManifestDigest: manifest.ManifestDigest,
		CandidateHead:                   manifest.Repository.CandidateHead,
		VerificationLabel:               "workflow_verified_mock",
		ObservedEvidence: []SandboxVerificationEvidence{
			{EvidenceKind: "exact_head_workflow_outcome", CandidateHead: manifest.Repository.CandidateHead},
		},
	}
	if err := ValidateMigrationVerificationBoundary(value, manifest); err != nil {
		t.Fatalf("valid exact-head mock verification rejected: %v", err)
	}
	value.ObservedEvidence[0].CandidateHead = "ffffffffffffffffffffffffffffffffffffffff"
	if err := ValidateMigrationVerificationBoundary(value, manifest); err == nil {
		t.Fatal("expected stale mock evidence to fail")
	}
	value.ObservedEvidence = nil
	if err := ValidateMigrationVerificationBoundary(value, manifest); err == nil {
		t.Fatal("expected missing exact-head mock evidence to fail")
	}
}

func readExecutionManifestFixture(t *testing.T, name string) ExecutionManifestBoundary {
	t.Helper()
	path := filepath.Join("..", "..", "tests", "fixtures", "contracts", "consumer-execution-manifest", name)
	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	var value ExecutionManifestBoundary
	if err := json.Unmarshal(payload, &value); err != nil {
		t.Fatalf("decode fixture: %v", err)
	}
	return value
}
