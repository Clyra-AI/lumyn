package tests

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/Clyra-AI/lumyn/internal/authorization"
	"github.com/Clyra-AI/lumyn/internal/evidence"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

const m1ManifestPath = "examples/migration-packs/benchmark-manifest.json"

type m1Manifest struct {
	SchemaVersion string       `json:"schema_version"`
	FrozenAt      string       `json:"frozen_at"`
	Scenarios     []m1Scenario `json:"scenarios"`
	Minimums      m1Minimums   `json:"minimums"`
}

type m1Minimums struct {
	Deterministic  int `json:"deterministic"`
	AgentAssisted  int `json:"agent_assisted"`
	VisibleRepos   int `json:"visible_repositories"`
	BlockedClasses int `json:"blocked_classes"`
}

type m1Scenario struct {
	ID                   string             `json:"id"`
	Route                string             `json:"route"`
	ChangeClass          string             `json:"change_class"`
	RepositoryShape      string             `json:"repository_shape"`
	PackPath             string             `json:"pack_path"`
	BeforePath           string             `json:"before_path"`
	ExpectedPath         string             `json:"expected_path"`
	IntegrationGraphPath string             `json:"integration_graph_path"`
	Digests              map[string]string  `json:"digests"`
	SupportingFiles      []m1SupportingFile `json:"supporting_files,omitempty"`
}

type m1SupportingFile struct {
	BeforePath     string `json:"before_path"`
	ExpectedPath   string `json:"expected_path"`
	BeforeDigest   string `json:"before_digest"`
	ExpectedDigest string `json:"expected_digest"`
}

type m1Pack struct {
	SchemaVersion  string        `json:"schema_version"`
	ScenarioID     string        `json:"scenario_id"`
	Title          string        `json:"title"`
	Route          string        `json:"route"`
	RouteRationale string        `json:"route_rationale,omitempty"`
	ChangeClass    string        `json:"change_class"`
	Operation      m1Operation   `json:"operation"`
	Provenance     m1Provenance  `json:"provenance"`
	GroundTruth    m1GroundTruth `json:"ground_truth"`
}

type m1Operation struct {
	Kind string `json:"kind"`
	From string `json:"from"`
	To   string `json:"to"`
}

type m1Provenance struct {
	SourceKind          string `json:"source_kind"`
	SourceID            string `json:"source_id"`
	SourcePath          string `json:"source_path"`
	SourceDigest        string `json:"source_digest"`
	License             string `json:"license"`
	Attribution         string `json:"attribution"`
	Redistribution      bool   `json:"redistribution_allowed"`
	ProviderEndorsement bool   `json:"provider_endorsement"`
	CustomerProof       bool   `json:"customer_proof"`
}

type m1GroundTruth struct {
	Affected             []string `json:"affected"`
	Unaffected           []string `json:"unaffected"`
	Uncertain            []string `json:"uncertain"`
	Unsupported          []string `json:"unsupported"`
	ExpectedRoute        string   `json:"expected_route"`
	ExpectedEdits        []string `json:"expected_edits"`
	ExpectedCommands     []string `json:"expected_commands"`
	ExpectedVerification []string `json:"expected_verification"`
}

type m1IntegrationGraph struct {
	SchemaVersion string        `json:"schema_version"`
	ScenarioID    string        `json:"scenario_id"`
	Nodes         []m1GraphNode `json:"nodes"`
}

type m1GraphNode struct {
	SiteID         string   `json:"site_id"`
	Path           string   `json:"path"`
	Line           int      `json:"line"`
	Classification string   `json:"classification"`
	Symbol         string   `json:"symbol"`
	Via            []string `json:"via"`
}

type m1NegativeSuite struct {
	SchemaVersion string           `json:"schema_version"`
	Cases         []m1NegativeCase `json:"cases"`
}

type m1NegativeCase struct {
	ID                      string `json:"id"`
	Category                string `json:"category"`
	InputSummary            string `json:"input_summary"`
	ExpectedStatus          string `json:"expected_status"`
	ReasonCode              string `json:"reason_code"`
	SpeculativePatchAllowed bool   `json:"speculative_patch_allowed"`
}

type m1BaselineProtocol struct {
	SchemaVersion                   string           `json:"schema_version"`
	FrozenAt                        string           `json:"frozen_at"`
	Comparator                      string           `json:"comparator"`
	LumynTreatment                  string           `json:"lumyn_treatment"`
	ControlVariables                []string         `json:"control_variables"`
	ControlBinding                  m1ControlBinding `json:"control_binding"`
	PrimaryMetric                   m1PrimaryMetric  `json:"primary_metric"`
	Guardrails                      []string         `json:"guardrails"`
	UnmatchedEngineComparisonCausal bool             `json:"unmatched_engine_comparison_is_causal"`
	LiveExecutionAllowed            bool             `json:"live_execution_allowed"`
}

type m1ControlBinding struct {
	Timing             string `json:"timing"`
	ValueRecord        string `json:"value_record"`
	EqualityRule       string `json:"equality_rule"`
	MissingValuePolicy string `json:"missing_value_policy"`
}

type m1PrimaryMetric struct {
	Name                        string `json:"name"`
	ImprovementThresholdPercent int    `json:"improvement_threshold_percent"`
}

type m1HoldoutManifest struct {
	SchemaVersion         string   `json:"schema_version"`
	Status                string   `json:"status"`
	SuiteNamespace        string   `json:"suite_namespace"`
	CommitmentAlgorithm   string   `json:"commitment_algorithm"`
	ProvisioningResultRef string   `json:"provisioning_result_ref"`
	CommittedFields       []string `json:"committed_fields"`
	ProhibitedFields      []string `json:"prohibited_fields"`
}

type m1WalkingSkeleton struct {
	SchemaVersion       string `json:"schema_version"`
	ScenarioID          string `json:"scenario_id"`
	EventDigest         string `json:"event_digest"`
	InstallationDigest  string `json:"installation_digest"`
	ImpactDigest        string `json:"impact_digest"`
	AuthorizationDigest string `json:"authorization_digest"`
	PlanDigest          string `json:"plan_digest"`
	CandidateDigest     string `json:"candidate_digest"`
	VerificationDigest  string `json:"verification_digest"`
	PRBundleDigest      string `json:"pr_bundle_digest"`
	AgentRoute          string `json:"agent_route"`
	ExternalWriteMode   string `json:"external_write_mode"`
	ContractBinding     string `json:"contract_binding"`
	Status              string `json:"status"`
}

type m1AgentRequest struct {
	ScenarioID        string
	PlanDigest        string
	Path              string
	Input             string
	Operation         m1Operation
	ExpectedCandidate string
}

type m1AgentResult struct {
	SchemaVersion string   `json:"schema_version"`
	Status        string   `json:"status"`
	SessionMode   string   `json:"session_mode"`
	Runner        string   `json:"runner"`
	Candidate     string   `json:"candidate"`
	ChangedPaths  []string `json:"changed_paths"`
}

type m1AgentRunner interface {
	Run(context.Context, m1AgentRequest) (m1AgentResult, error)
}

type m1DeterministicFakeAgent struct {
	Route m1FakeRoute
}

var _ m1AgentRunner = m1DeterministicFakeAgent{}

func TestM1Corpus(t *testing.T) {
	root := repoRoot(t)
	manifest := readM1JSON[m1Manifest](t, filepath.Join(root, m1ManifestPath))
	if manifest.SchemaVersion != "lumyn.m1.benchmark-manifest/v1" || manifest.FrozenAt == "" {
		t.Fatal("M1 manifest identity and freeze time are required")
	}
	if manifest.Minimums != (m1Minimums{Deterministic: 3, AgentAssisted: 3, VisibleRepos: 6, BlockedClasses: 7}) {
		t.Fatalf("M1 minimums drifted: %+v", manifest.Minimums)
	}
	if len(manifest.Scenarios) < manifest.Minimums.VisibleRepos {
		t.Fatalf("visible scenarios = %d, want at least %d", len(manifest.Scenarios), manifest.Minimums.VisibleRepos)
	}

	seen := map[string]bool{}
	routeCounts := map[string]int{}
	shapeCounts := map[string]int{}
	changeClasses := map[string]bool{}
	for _, scenario := range manifest.Scenarios {
		if scenario.ID == "" || seen[scenario.ID] {
			t.Fatalf("scenario id must be non-empty and unique: %q", scenario.ID)
		}
		seen[scenario.ID] = true
		routeCounts[scenario.Route]++
		shapeCounts[scenario.RepositoryShape]++
		changeClasses[scenario.ChangeClass] = true
		validateM1Scenario(t, root, scenario)
	}
	if routeCounts["deterministic"] < manifest.Minimums.Deterministic || routeCounts["agent_assisted"] < manifest.Minimums.AgentAssisted {
		t.Fatalf("route coverage drifted: %+v", routeCounts)
	}
	if shapeCounts["direct"] == 0 || shapeCounts["wrapper_heavy"] < 3 {
		t.Fatalf("repository-shape coverage drifted: %+v", shapeCounts)
	}
	for _, class := range []string{
		"operation_rename",
		"request_property_relocation",
		"response_property_relocation",
		"wrapper_adaptation",
		"signature_type_adaptation",
		"related_test_repair",
	} {
		if !changeClasses[class] {
			t.Fatalf("missing M1 change class %s", class)
		}
	}
}

func TestM1VisibleRepositoriesRunPinnedOfflineVerification(t *testing.T) {
	root := repoRoot(t)
	requireM1PinnedRuntime(t)
	manifest := readM1JSON[m1Manifest](t, filepath.Join(root, m1ManifestPath))
	for _, scenario := range manifest.Scenarios {
		scenario := scenario
		t.Run(scenario.ID, func(t *testing.T) {
			repoDir := filepath.Join(root, "examples", "consumer-repos", scenario.ID)
			packageJSON := readM1JSON[struct {
				Name           string            `json:"name"`
				Private        bool              `json:"private"`
				Version        string            `json:"version"`
				Type           string            `json:"type"`
				Engines        map[string]string `json:"engines"`
				PackageManager string            `json:"packageManager"`
				Scripts        map[string]string `json:"scripts"`
			}](t, filepath.Join(repoDir, "package.json"))
			if !packageJSON.Private || packageJSON.Type != "module" || packageJSON.Engines["node"] != "22.16.0" ||
				packageJSON.Engines["npm"] != "11.4.1" || packageJSON.PackageManager != "npm@11.4.1" ||
				packageJSON.Scripts["test"] != fmt.Sprintf(
					"node --no-warnings --experimental-strip-types --experimental-vm-modules ../../verification/run.mjs %s ../../candidates/%s/src/client.ts",
					scenario.ID, scenario.ID,
				) {
				t.Fatalf("scenario %s runtime or npm test contract drifted: %+v", scenario.ID, packageJSON)
			}
			command := m1OfflineNPMTestCommand(t, repoDir)
			output, err := command.CombinedOutput()
			if err != nil {
				t.Fatalf("scenario %s npm test: %v\n%s", scenario.ID, err, output)
			}
			want := scenario.ID + ": red-before/green-candidate offline behavior verified"
			if !strings.Contains(string(output), want) {
				t.Fatalf("scenario %s npm test lacks proof line %q:\n%s", scenario.ID, want, output)
			}
		})
	}
}

func TestM1VerifierRejectsAmbientAuthority(t *testing.T) {
	t.Setenv("LUMYN_M1_HOST_SECRET", "must-not-cross")
	root := repoRoot(t)
	repoDir := filepath.Join(root, "examples/consumer-repos/det-operation-rename")
	maliciousCandidate := filepath.Join(t.TempDir(), "client.ts")
	if err := os.WriteFile(maliciousCandidate, []byte(`
export function createCharge() {
  return expect["con" + "structor"]("return pro" + "cess")().version;
}

export function health() {
  return "ok";
}
`), 0o400); err != nil {
		t.Fatal(err)
	}
	nodePath, err := exec.LookPath("node")
	if err != nil {
		t.Fatalf("resolve pinned node: %v", err)
	}
	npmPath, err := exec.LookPath("npm")
	if err != nil {
		t.Fatalf("resolve pinned npm: %v", err)
	}
	command := exec.Command(
		nodePath,
		"--no-warnings",
		"--experimental-strip-types",
		"--experimental-vm-modules",
		filepath.Join(root, "examples/verification/run.mjs"),
		"det-operation-rename",
		maliciousCandidate,
	)
	command.Dir = repoDir
	command.Env, _ = m1OfflineCommandEnv(t, nodePath, npmPath, "")
	output, err := command.CombinedOutput()
	if err == nil || !strings.Contains(string(output), "Code generation from strings disallowed") {
		t.Fatalf("restricted verifier accepted computed host-realm constructor escape: %v\n%s", err, output)
	}
	for _, entry := range command.Env {
		if strings.HasPrefix(entry, "LUMYN_M1_HOST_SECRET=") {
			t.Fatalf("restricted verifier inherited an ambient test credential: %s", entry)
		}
	}
}

func requireM1PinnedRuntime(t *testing.T) map[string]string {
	t.Helper()
	observed := map[string]string{}
	for _, runtime := range []struct {
		program string
		args    []string
		want    string
	}{
		{program: "node", args: []string{"--version"}, want: "v22.16.0"},
		{program: "npm", args: []string{"--version"}, want: "11.4.1"},
	} {
		output, err := exec.Command(runtime.program, runtime.args...).CombinedOutput()
		if err != nil || strings.TrimSpace(string(output)) != runtime.want {
			t.Fatalf("M1 requires exact %s %s; got %q (%v)", runtime.program, runtime.want, strings.TrimSpace(string(output)), err)
		}
		observed[runtime.program] = strings.TrimSpace(string(output))
	}
	return observed
}

func m1OfflineNPMTestCommand(t *testing.T, repoDir string) *exec.Cmd {
	t.Helper()
	command, _ := m1OfflineNPMTestCommandWithCleanup(t, repoDir, "")
	return command
}

func m1OfflineNPMTestCommandWithCleanup(t *testing.T, repoDir string, target string) (*exec.Cmd, func() string) {
	t.Helper()
	npmPath, err := exec.LookPath("npm")
	if err != nil {
		t.Fatalf("resolve pinned npm: %v", err)
	}
	nodePath, err := exec.LookPath("node")
	if err != nil {
		t.Fatalf("resolve pinned node: %v", err)
	}
	command := exec.Command(npmPath, "test", "--silent")
	command.Dir = repoDir
	environment, cleanup := m1OfflineCommandEnv(t, nodePath, npmPath, target)
	command.Env = environment
	return command, cleanup
}

func m1OfflineCommandEnv(t *testing.T, nodePath string, npmPath string, target string) ([]string, func() string) {
	t.Helper()
	isolationRoot, err := os.MkdirTemp("", "lumyn-m1-command-")
	if err != nil {
		t.Fatalf("create isolated command root: %v", err)
	}
	removed := false
	cleanupProof := digestM1Canonical(map[string]any{"isolated_command_root_removed": true})
	cleanup := func() string {
		if removed {
			return cleanupProof
		}
		removed = true
		if err := os.RemoveAll(isolationRoot); err != nil {
			t.Fatalf("remove isolated command root: %v", err)
		}
		if _, err := os.Stat(isolationRoot); !os.IsNotExist(err) {
			t.Fatalf("isolated command root still exists after cleanup: %v", err)
		}
		return cleanupProof
	}
	t.Cleanup(func() { cleanup() })
	userNPMConfig := filepath.Join(isolationRoot, "user.npmrc")
	globalNPMConfig := filepath.Join(isolationRoot, "global.npmrc")
	for _, configPath := range []string{userNPMConfig, globalNPMConfig} {
		if err := os.WriteFile(configPath, nil, 0o400); err != nil {
			t.Fatalf("create isolated npm config: %v", err)
		}
	}
	pathEntries := []string{filepath.Dir(nodePath), filepath.Dir(npmPath), "/usr/bin", "/bin"}
	seenPathEntries := map[string]bool{}
	cleanPathEntries := make([]string, 0, len(pathEntries))
	for _, entry := range pathEntries {
		if entry != "" && !seenPathEntries[entry] {
			seenPathEntries[entry] = true
			cleanPathEntries = append(cleanPathEntries, entry)
		}
	}
	environment := []string{
		"CI=1",
		"HOME=" + isolationRoot,
		"USERPROFILE=" + isolationRoot,
		"TMPDIR=" + isolationRoot,
		"TMP=" + isolationRoot,
		"TEMP=" + isolationRoot,
		"PATH=" + strings.Join(cleanPathEntries, string(os.PathListSeparator)),
		"NODE_NO_WARNINGS=1",
		"NO_COLOR=1",
		"npm_config_audit=false",
		"npm_config_cache=" + filepath.Join(isolationRoot, "npm-cache"),
		"npm_config_engine_strict=true",
		"npm_config_fund=false",
		"npm_config_globalconfig=" + globalNPMConfig,
		"npm_config_ignore_scripts=true",
		"npm_config_offline=true",
		"npm_config_registry=https://invalid.invalid",
		"npm_config_update_notifier=false",
		"npm_config_userconfig=" + userNPMConfig,
	}
	if target != "" {
		environment = append(environment, "LUMYN_M1_VERIFICATION_TARGET="+target)
	}
	if target == "candidate" {
		environment = append(environment, "LUMYN_M1_CANDIDATE_PATH=src/client.ts")
	}
	return environment, cleanup
}

func TestM1BlockedCorpus(t *testing.T) {
	root := repoRoot(t)
	suite := readM1JSON[m1NegativeSuite](t, filepath.Join(root, "examples/negative/blocked-scenarios.json"))
	if suite.SchemaVersion != "lumyn.m1.blocked-suite/v1" {
		t.Fatalf("blocked suite version = %q", suite.SchemaVersion)
	}
	want := map[string]bool{
		"missing_business_value":  false,
		"auth_redesign":           false,
		"event_semantics":         false,
		"ambiguous_evidence":      false,
		"production_only":         false,
		"scope_escape":            false,
		"inadequate_verification": false,
	}
	if len(suite.Cases) != len(want) {
		t.Fatalf("blocked cases = %d, want exactly %d", len(suite.Cases), len(want))
	}
	seenIDs := map[string]bool{}
	for _, item := range suite.Cases {
		if _, ok := want[item.Category]; !ok {
			t.Fatalf("unknown blocked category %q", item.Category)
		}
		if item.ID == "" || item.InputSummary == "" || item.ReasonCode == "" {
			t.Fatalf("blocked case is incomplete: %+v", item)
		}
		if seenIDs[item.ID] || want[item.Category] {
			t.Fatalf("blocked case id or category is duplicated: %+v", item)
		}
		seenIDs[item.ID] = true
		if item.ExpectedStatus != "blocked" && item.ExpectedStatus != "needs_input" {
			t.Fatalf("blocked case %s has speculative status %q", item.ID, item.ExpectedStatus)
		}
		if item.SpeculativePatchAllowed {
			t.Fatalf("blocked case %s permits a speculative patch", item.ID)
		}
		want[item.Category] = true
	}
	for category, covered := range want {
		if !covered {
			t.Fatalf("missing blocked category %s", category)
		}
	}
}

func TestM1BaselineProtocol(t *testing.T) {
	root := repoRoot(t)
	protocol := readM1JSON[m1BaselineProtocol](t, filepath.Join(root, "examples/migration-packs/generic-agent-baseline-v1.json"))
	if protocol.SchemaVersion != "lumyn.m1.generic-agent-baseline/v1" || protocol.FrozenAt == "" {
		t.Fatal("baseline protocol identity and freeze time are required")
	}
	wantControls := []string{
		"actual_model_provider",
		"actual_model_version",
		"agent_runner_adapter",
		"agent_runner_executable_digest",
		"agent_runner_version",
		"attempt_budget",
		"auth_mode",
		"authoritative_migration_evidence",
		"commands",
		"context_access_ceiling",
		"cost_budget",
		"credential_owner",
		"engineer_role",
		"entitlement_class",
		"execution_funding_mode",
		"repository_snapshot",
		"time_budget",
		"token_budget",
		"tools",
		"usage_billing_owner",
		"verification_commands_and_time",
	}
	if !reflect.DeepEqual(protocol.ControlVariables, wantControls) {
		t.Fatalf("baseline controls drifted:\n got %v\nwant %v", protocol.ControlVariables, wantControls)
	}
	if protocol.Comparator != "same_capable_engine_without_lumyn_orchestration" ||
		protocol.LumynTreatment != "orchestration_impact_routing_boundary_enforcement_independent_verification_delivery_and_status" {
		t.Fatal("baseline must isolate Lumyn orchestration as the treatment")
	}
	if protocol.ControlBinding != (m1ControlBinding{
		Timing:             "before_each_matched_trial",
		ValueRecord:        "immutable_trial_manifest",
		EqualityRule:       "exact_match_between_arms",
		MissingValuePolicy: "block_trial",
	}) {
		t.Fatalf("baseline control binding drifted: %+v", protocol.ControlBinding)
	}
	if protocol.PrimaryMetric.Name != "consumer_maintainer_hands_on_time" || protocol.PrimaryMetric.ImprovementThresholdPercent != 30 {
		t.Fatalf("primary metric drifted: %+v", protocol.PrimaryMetric)
	}
	if len(protocol.Guardrails) < 3 || protocol.UnmatchedEngineComparisonCausal || protocol.LiveExecutionAllowed {
		t.Fatal("baseline guardrail or offline causal boundary drifted")
	}
}

func TestM1HoldoutIsolation(t *testing.T) {
	root := repoRoot(t)
	path := filepath.Join(root, "examples/holdout-manifest.json")
	manifest := readM1JSON[m1HoldoutManifest](t, path)
	if manifest.SchemaVersion != "lumyn.m1.holdout-manifest/v1" ||
		manifest.Status != "provisioning_required" ||
		manifest.SuiteNamespace != "private://lumyn/m1/v1" ||
		manifest.CommitmentAlgorithm != "hmac-sha256" {
		t.Fatalf("holdout boundary drifted: %+v", manifest)
	}
	if manifest.ProvisioningResultRef != ".factory/artifacts/lifecycle-evidence/M1/holdout-result.json" {
		t.Fatalf("holdout provisioning result ref drifted: %q", manifest.ProvisioningResultRef)
	}
	wantProhibited := []string{
		"answer_key",
		"expected_labels",
		"expected_patches",
		"inputs",
		"plaintext_content_digest",
		"raw_traces",
		"repository_url",
	}
	wantCommitted := []string{
		"opaque_case_count",
		"license_posture",
		"frozen_suite_commitment",
		"encrypted_or_hmac_artifact_commitments",
	}
	if !reflect.DeepEqual(manifest.CommittedFields, wantCommitted) {
		t.Fatalf("holdout commitments drifted: %v", manifest.CommittedFields)
	}
	if !reflect.DeepEqual(manifest.ProhibitedFields, wantProhibited) {
		t.Fatalf("holdout prohibitions drifted: %v", manifest.ProhibitedFields)
	}
	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]any
	if err := json.Unmarshal(payload, &raw); err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range append(wantProhibited, "suite_commitment", "case_ids", "opaque_case_ids") {
		if containsM1JSONKey(raw, forbidden) {
			t.Fatalf("source-safe holdout manifest exposes %s", forbidden)
		}
	}
}

func TestM1WalkingSkeleton(t *testing.T) {
	root := repoRoot(t)
	manifest := readM1JSON[m1Manifest](t, filepath.Join(root, m1ManifestPath))
	var scenario m1Scenario
	for _, candidate := range manifest.Scenarios {
		if candidate.ID == "det-operation-rename" {
			scenario = candidate
			break
		}
	}
	if scenario.ID == "" {
		t.Fatal("walking-skeleton scenario is missing")
	}
	observedAt := time.Now().UTC().Truncate(time.Second)
	first := runM1WalkingSkeleton(t, root, scenario, observedAt)
	second := runM1WalkingSkeleton(t, root, scenario, observedAt)
	left, err := json.Marshal(first)
	if err != nil {
		t.Fatal(err)
	}
	right, err := json.Marshal(second)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(left, right) {
		t.Fatalf("walking skeleton is not byte-stable:\n%s\n%s", left, right)
	}
	if first.Status != "static_verified_local_pr_bundle" || first.AgentRoute != "disabled" || first.ExternalWriteMode != "manual_pr_bundle_only" {
		t.Fatalf("walking-skeleton boundary drifted: %+v", first)
	}
}

func TestM1WalkingSkeletonRejectsSubstitutedCandidate(t *testing.T) {
	root := repoRoot(t)
	manifest := readM1JSON[m1Manifest](t, filepath.Join(root, m1ManifestPath))
	var scenario m1Scenario
	for _, candidate := range manifest.Scenarios {
		if candidate.ID == "det-operation-rename" {
			scenario = candidate
			break
		}
	}
	if scenario.ID == "" {
		t.Fatal("walking-skeleton scenario is missing")
	}
	expected := readM1File(t, filepath.Join(root, scenario.ExpectedPath))
	mutated := strings.Replace(expected, "paymentIntents", "paymentIntentsUnexpected", 1)
	if mutated == expected {
		t.Fatal("walking-skeleton substitution mutation did not change the generated candidate")
	}
	view := createM1ReadOnlyVerificationView(t, root, scenario)
	t.Cleanup(func() {
		if !view.Removed {
			removeM1VerificationView(t, &view, scenario.ID)
		}
	})
	promoteM1VerificationViewCandidate(t, root, scenario, mutated, &view)
	command, cleanup := m1OfflineNPMTestCommandWithCleanup(t, view.RepoDir, "candidate")
	output, err := command.CombinedOutput()
	cleanup()
	if err == nil {
		t.Fatalf("walking skeleton accepted a substituted generated candidate:\n%s", output)
	}
	if strings.Contains(string(output), scenario.ID+": exact candidate target contract verified") {
		t.Fatalf("walking skeleton fell back to checked-in expected bytes:\n%s", output)
	}
	removeM1VerificationView(t, &view, scenario.ID)
}

func TestM1BaselineRejectsInfrastructureFailures(t *testing.T) {
	root := repoRoot(t)
	manifest := readM1JSON[m1Manifest](t, filepath.Join(root, m1ManifestPath))
	var scenario m1Scenario
	for _, candidate := range manifest.Scenarios {
		if candidate.ID == "det-operation-rename" {
			scenario = candidate
			break
		}
	}
	if scenario.ID == "" {
		t.Fatal("baseline infrastructure-failure scenario is missing")
	}
	for _, testCase := range []struct {
		name   string
		mutate func(string) error
	}{
		{name: "missing source", mutate: os.Remove},
		{name: "syntax error", mutate: func(path string) error {
			if err := os.Chmod(path, 0o644); err != nil {
				return err
			}
			return os.WriteFile(path, []byte("export function broken(\n"), 0o444)
		}},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			view := createM1ReadOnlyVerificationView(t, root, scenario)
			t.Cleanup(func() {
				if !view.Removed {
					removeM1VerificationView(t, &view, scenario.ID)
				}
			})
			unlockM1VerificationView(t, &view)
			baselinePath := filepath.Join(view.RepoDir, filepath.FromSlash(m1ScenarioRelativePath(t, scenario, scenario.BeforePath)))
			if err := testCase.mutate(baselinePath); err != nil {
				t.Fatal(err)
			}
			lockM1VerificationView(t, &view)
			command, cleanup := m1OfflineNPMTestCommandWithCleanup(t, view.RepoDir, "baseline")
			output, err := command.CombinedOutput()
			cleanup()
			if err == nil {
				t.Fatalf("baseline infrastructure failure unexpectedly passed:\n%s", output)
			}
			if strings.Contains(string(output), scenario.ID+": baseline target contract rejected") {
				t.Fatalf("baseline infrastructure failure was misclassified as an expected target mismatch:\n%s", output)
			}
			removeM1VerificationView(t, &view, scenario.ID)
		})
	}
}

func TestM1CanonicalArtifactDigestsBindSemanticFields(t *testing.T) {
	root := repoRoot(t)
	manifest := readM1JSON[m1Manifest](t, filepath.Join(root, m1ManifestPath))
	var scenario m1Scenario
	for _, candidate := range manifest.Scenarios {
		if candidate.ID == "det-operation-rename" {
			scenario = candidate
			break
		}
	}
	if scenario.ID == "" {
		t.Fatal("canonical digest scenario is missing")
	}
	pack := readM1JSON[m1Pack](t, filepath.Join(root, scenario.PackPath))
	before := readM1File(t, filepath.Join(root, scenario.BeforePath))
	candidate, err := applyM1Operation(pack.Operation, before)
	if err != nil {
		t.Fatal(err)
	}
	observation := observeM1RepositoryVerification(t, root, scenario, candidate)
	chain := buildM1CanonicalContractChain(t, root, scenario, candidate, observation, time.Date(2026, 7, 28, 12, 0, 0, 0, time.UTC))
	wantDigests := map[string]string{
		"provider-change-event":  chain.EventDigest,
		"migration-pack":         chain.PackDigest,
		"consumer-installation":  chain.InstallationDigest,
		"event-authorization":    chain.AuthorizationDigest,
		"impact-report":          chain.ImpactDigest,
		"migration-plan":         chain.PlanDigest,
		"candidate-manifest":     chain.CandidateDigest,
		"migration-verification": chain.VerificationDigest,
		"export-result":          chain.ExportDigest,
	}
	mutations := map[string]func(map[string]any){
		"provider-change-event": func(value map[string]any) { value["audience_id"] = "audience.synthetic_other" },
		"migration-pack": func(value map[string]any) {
			m1ObjectAt(t, value, "semantic_changes", 0)["intent"] = "A different schema-valid provider intent."
		},
		"consumer-installation": func(value map[string]any) {
			m1Object(t, value, "budgets")["max_tokens"] = float64(49_999)
		},
		"event-authorization": func(value map[string]any) {
			m1Object(t, value, "approval")["evaluated_at"] = "2026-07-25T17:06:00Z"
		},
		"impact-report": func(value map[string]any) {
			m1ObjectAt(t, value, "findings", 0)["reason"] = "A different schema-valid impact reason."
		},
		"migration-plan": func(value map[string]any) {
			m1ObjectAt(t, value, "items", 0)["rationale"] = "A different schema-valid route rationale."
		},
		"candidate-manifest": func(value map[string]any) {
			m1ObjectAt(t, value, "changes", 0)["rationale"] = "A different schema-valid edit rationale."
		},
		"migration-verification": func(value map[string]any) {
			value["residual_risks"] = append(value["residual_risks"].([]any), "A different schema-valid residual risk.")
		},
		"export-result": func(value map[string]any) {
			m1Object(t, value, "delivery")["source_branch"] = "lumyn/synthetic-other"
		},
	}
	for artifactName, mutate := range mutations {
		mutated := cloneM1JSONObject(t, chain.Artifacts[artifactName])
		mutate(mutated)
		validateM1Schema(t, root, artifactName+".schema.json", mutated)
		if got := digestM1CanonicalArtifact(t, artifactName, mutated); got == wantDigests[artifactName] {
			t.Fatalf("%s semantic mutation did not change its canonical artifact digest", artifactName)
		}
	}
	mutatedHead := cloneM1JSONObject(t, chain.Artifacts["candidate-manifest"])
	mutatedHead["candidate_head"] = strings.Repeat("f", 40)
	m1Object(t, mutatedHead, "independent_verification")["candidate_head"] = strings.Repeat("f", 40)
	validateM1Schema(t, root, "candidate-manifest.schema.json", mutatedHead)
	if got := digestM1CanonicalArtifact(t, "candidate-manifest", mutatedHead); got == chain.CandidateDigest {
		t.Fatal("real Git candidate-head mutation did not change the candidate-manifest digest")
	}
}

func TestM1MutationAndReplayReject(t *testing.T) {
	root := repoRoot(t)
	manifest := readM1JSON[m1Manifest](t, filepath.Join(root, m1ManifestPath))
	scenario := manifest.Scenarios[0]
	pack := readM1JSON[m1Pack](t, filepath.Join(root, scenario.PackPath))
	before := readM1File(t, filepath.Join(root, scenario.BeforePath))
	expected := readM1File(t, filepath.Join(root, scenario.ExpectedPath))
	generated, err := applyM1Operation(pack.Operation, before)
	if err != nil {
		t.Fatal(err)
	}
	if generated != expected {
		t.Fatal("baseline fixture must verify before mutation")
	}
	mutated := strings.Replace(expected, pack.Operation.To, pack.Operation.To+"Unexpected", 1)
	if mutated == expected || generated == mutated {
		t.Fatal("mutated ground truth did not invalidate verification")
	}
	currentPackDigest := digestM1Bytes(readM1FileBytes(t, filepath.Join(root, scenario.PackPath)))
	if err := validateM1ReplayBinding(currentPackDigest, strings.Repeat("0", 64)); err == nil {
		t.Fatal("cross-event pack replay unexpectedly validated")
	}
}

func TestM1NoLiveRoute(t *testing.T) {
	for _, route := range []m1FakeRoute{
		{Runner: "codex"},
		{Runner: "claude_code"},
		{Runner: "cursor"},
		{Network: true},
		{Credential: true},
		{ExternalWrite: true},
	} {
		if err := route.validate(); err == nil {
			t.Fatalf("live or external route unexpectedly accepted: %+v", route)
		}
	}
	if err := (m1FakeRoute{}).validate(); err != nil {
		t.Fatalf("deterministic fake route rejected: %v", err)
	}
}

func TestM1AdapterConformanceCancellationReject(t *testing.T) {
	root := repoRoot(t)
	manifest := readM1JSON[m1Manifest](t, filepath.Join(root, m1ManifestPath))
	scenario := manifest.Scenarios[3]
	pack := readM1JSON[m1Pack](t, filepath.Join(root, scenario.PackPath))
	before := readM1File(t, filepath.Join(root, scenario.BeforePath))
	expected := readM1File(t, filepath.Join(root, scenario.ExpectedPath))
	request := m1AgentRequest{
		ScenarioID:        scenario.ID,
		PlanDigest:        digestM1Canonical(map[string]any{"scenario_id": scenario.ID}),
		Path:              scenario.BeforePath,
		Input:             before,
		Operation:         pack.Operation,
		ExpectedCandidate: expected,
	}
	adapter := m1DeterministicFakeAgent{}
	result, err := adapter.Run(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if result.SchemaVersion != "lumyn.agent-runner-result/v1" || result.Status != "candidate_ready" ||
		result.SessionMode != "clean_ephemeral" || result.Runner != "deterministic_fake" ||
		result.Candidate != expected || !reflect.DeepEqual(result.ChangedPaths, []string{scenario.BeforePath}) {
		t.Fatalf("fake adapter output is not normalized: %+v", result)
	}

	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := adapter.Run(canceled, request); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled fake adapter run = %v, want context.Canceled", err)
	}
	if _, err := (m1DeterministicFakeAgent{Route: m1FakeRoute{Runner: "codex"}}).Run(context.Background(), request); err == nil {
		t.Fatal("fake adapter silently fell back to a live runner")
	}
}

func validateM1Scenario(t *testing.T, root string, scenario m1Scenario) {
	t.Helper()
	if scenario.Route != "deterministic" && scenario.Route != "agent_assisted" {
		t.Fatalf("scenario %s has invalid route %q", scenario.ID, scenario.Route)
	}
	if scenario.RepositoryShape != "direct" && scenario.RepositoryShape != "wrapper_heavy" {
		t.Fatalf("scenario %s has invalid repository shape %q", scenario.ID, scenario.RepositoryShape)
	}
	paths := map[string]string{
		"pack":              scenario.PackPath,
		"before":            scenario.BeforePath,
		"expected":          scenario.ExpectedPath,
		"integration_graph": scenario.IntegrationGraphPath,
	}
	if len(scenario.Digests) != len(paths) {
		t.Fatalf("scenario %s digest set must bind exactly pack, before, expected, and integration graph", scenario.ID)
	}
	for label, relativePath := range paths {
		if relativePath == "" || filepath.IsAbs(relativePath) || strings.Contains(relativePath, "..") {
			t.Fatalf("scenario %s has unsafe %s path %q", scenario.ID, label, relativePath)
		}
		payload := readM1FileBytes(t, filepath.Join(root, relativePath))
		if got, want := digestM1Bytes(payload), scenario.Digests[label]; got != want {
			t.Fatalf("scenario %s %s digest = %s, want %s", scenario.ID, label, got, want)
		}
	}
	for index, supporting := range scenario.SupportingFiles {
		for label, item := range map[string]struct {
			path   string
			digest string
		}{
			"before":   {supporting.BeforePath, supporting.BeforeDigest},
			"expected": {supporting.ExpectedPath, supporting.ExpectedDigest},
		} {
			if item.path == "" || filepath.IsAbs(item.path) || strings.Contains(item.path, "..") {
				t.Fatalf("scenario %s supporting file %d has unsafe %s path %q", scenario.ID, index, label, item.path)
			}
			if got := digestM1Bytes(readM1FileBytes(t, filepath.Join(root, item.path))); got != item.digest {
				t.Fatalf("scenario %s supporting file %d %s digest = %s, want %s", scenario.ID, index, label, got, item.digest)
			}
		}
	}

	pack := readM1JSON[m1Pack](t, filepath.Join(root, scenario.PackPath))
	if pack.SchemaVersion != "lumyn.m1.migration-pack/v1" || pack.ScenarioID != scenario.ID || pack.Route != scenario.Route || pack.ChangeClass != scenario.ChangeClass {
		t.Fatalf("scenario %s pack identity drifted: %+v", scenario.ID, pack)
	}
	if scenario.Route == "deterministic" {
		if pack.Operation.Kind != "replace_exact_once" || pack.Operation.From == "" || pack.Operation.To == "" || pack.Operation.From == pack.Operation.To {
			t.Fatalf("scenario %s deterministic operation is not bounded: %+v", scenario.ID, pack.Operation)
		}
	} else if pack.Operation.Kind != "agent_required" || pack.Operation.From != "" || pack.Operation.To != "" || strings.TrimSpace(pack.RouteRationale) == "" {
		t.Fatalf("scenario %s agent route exposes a deterministic recipe or lacks rationale: %+v", scenario.ID, pack)
	}
	if pack.Provenance.SourceKind != "synthetic" || pack.Provenance.SourceID == "" ||
		pack.Provenance.License != "CC0-1.0" || pack.Provenance.Attribution == "" ||
		!pack.Provenance.Redistribution || pack.Provenance.ProviderEndorsement || pack.Provenance.CustomerProof {
		t.Fatalf("scenario %s provenance is incomplete or overclaims proof: %+v", scenario.ID, pack.Provenance)
	}
	if pack.Provenance.SourcePath == "" || filepath.IsAbs(pack.Provenance.SourcePath) || strings.Contains(pack.Provenance.SourcePath, "..") ||
		pack.Provenance.SourcePath != filepath.ToSlash(filepath.Join(filepath.Dir(scenario.PackPath), "source.md")) {
		t.Fatalf("scenario %s source evidence path is unsafe or not colocated: %q", scenario.ID, pack.Provenance.SourcePath)
	}
	sourcePayload := readM1FileBytes(t, filepath.Join(root, pack.Provenance.SourcePath))
	if digestM1Bytes(sourcePayload) != pack.Provenance.SourceDigest {
		t.Fatalf("scenario %s source evidence digest drifted", scenario.ID)
	}
	if len(pack.GroundTruth.Affected) == 0 || len(pack.GroundTruth.Unaffected) == 0 ||
		pack.GroundTruth.Uncertain == nil || pack.GroundTruth.Unsupported == nil ||
		len(pack.GroundTruth.ExpectedEdits) == 0 || len(pack.GroundTruth.ExpectedCommands) == 0 ||
		len(pack.GroundTruth.ExpectedVerification) < 2 || pack.GroundTruth.ExpectedRoute != scenario.Route {
		t.Fatalf("scenario %s ground truth is incomplete: %+v", scenario.ID, pack.GroundTruth)
	}
	wantEdits := map[string]bool{
		m1ScenarioRelativePath(t, scenario, scenario.BeforePath): true,
	}
	for _, supporting := range scenario.SupportingFiles {
		wantEdits[m1ScenarioRelativePath(t, scenario, supporting.BeforePath)] = true
	}
	if len(pack.GroundTruth.ExpectedEdits) != len(wantEdits) {
		t.Fatalf("scenario %s expected edits do not exactly cover changed fixture files", scenario.ID)
	}
	for _, path := range pack.GroundTruth.ExpectedEdits {
		if !wantEdits[path] {
			t.Fatalf("scenario %s expected edit is outside the changed fixture set: %s", scenario.ID, path)
		}
	}
	if !reflect.DeepEqual(pack.GroundTruth.ExpectedCommands, []string{"npm test --silent"}) {
		t.Fatalf("scenario %s expected commands drifted: %v", scenario.ID, pack.GroundTruth.ExpectedCommands)
	}
	if scenario.Route == "agent_assisted" && (len(scenario.SupportingFiles) == 0 || len(pack.GroundTruth.ExpectedEdits) < 2) {
		t.Fatalf("scenario %s agent route is not a coordinated multi-file migration", scenario.ID)
	}

	before := readM1File(t, filepath.Join(root, scenario.BeforePath))
	sourceLines := map[string][]string{
		m1ScenarioRelativePath(t, scenario, scenario.BeforePath): strings.Split(before, "\n"),
	}
	for _, supporting := range scenario.SupportingFiles {
		sourceLines[m1ScenarioRelativePath(t, scenario, supporting.BeforePath)] = strings.Split(
			readM1File(t, filepath.Join(root, supporting.BeforePath)), "\n",
		)
	}
	graph := readM1JSON[m1IntegrationGraph](t, filepath.Join(root, scenario.IntegrationGraphPath))
	if graph.SchemaVersion != "lumyn.m1.integration-graph/v1" || graph.ScenarioID != scenario.ID || len(graph.Nodes) < 2 {
		t.Fatalf("scenario %s integration graph is incomplete", scenario.ID)
	}
	classifications := map[string]bool{}
	groundTruthRefs := map[string]map[string]bool{
		"affected":   {},
		"unaffected": {},
	}
	for _, ref := range pack.GroundTruth.Affected {
		groundTruthRefs["affected"][ref] = true
	}
	for _, ref := range pack.GroundTruth.Unaffected {
		groundTruthRefs["unaffected"][ref] = true
	}
	seenGraphRefs := map[string]bool{}
	graphPaths := map[string]bool{}
	graphNodes := map[string]m1GraphNode{}
	for _, node := range graph.Nodes {
		if node.SiteID == "" || node.Path == "" || node.Line < 1 || node.Symbol == "" || len(node.Via) == 0 {
			t.Fatalf("scenario %s graph node is incomplete: %+v", scenario.ID, node)
		}
		if node.Classification != "affected" && node.Classification != "unaffected" {
			t.Fatalf("scenario %s graph node has unsupported visible classification: %+v", scenario.ID, node)
		}
		lines, exists := sourceLines[node.Path]
		if !exists || node.Line > len(lines) || !strings.Contains(lines[node.Line-1], node.Symbol) {
			t.Fatalf("scenario %s graph node does not resolve to its source symbol: %+v", scenario.ID, node)
		}
		if _, exists := graphNodes[node.SiteID]; exists {
			t.Fatalf("scenario %s graph site id is duplicated: %s", scenario.ID, node.SiteID)
		}
		graphNodes[node.SiteID] = node
		graphPaths[node.Path] = true
		ref := fmt.Sprintf("%s:%d", node.Path, node.Line)
		if seenGraphRefs[ref] || !groundTruthRefs[node.Classification][ref] {
			t.Fatalf("scenario %s graph node is duplicated or not bound to ground truth: %+v", scenario.ID, node)
		}
		seenGraphRefs[ref] = true
		classifications[node.Classification] = true
	}
	if !classifications["affected"] || !classifications["unaffected"] {
		t.Fatalf("scenario %s graph lacks affected/unaffected evidence", scenario.ID)
	}
	if len(seenGraphRefs) != len(pack.GroundTruth.Affected)+len(pack.GroundTruth.Unaffected) {
		t.Fatalf("scenario %s graph does not exactly cover visible ground truth", scenario.ID)
	}
	if scenario.Route == "agent_assisted" {
		if len(graphPaths) < 2 {
			t.Fatalf("scenario %s wrapper-heavy graph does not cross files: %v", scenario.ID, graphPaths)
		}
		expectedSource := readM1File(t, filepath.Join(root, scenario.ExpectedPath))
		for _, supporting := range scenario.SupportingFiles {
			moduleRef := "./" + filepath.Base(supporting.BeforePath)
			if !strings.Contains(before, `from "`+moduleRef+`"`) || !strings.Contains(expectedSource, `from "`+moduleRef+`"`) {
				t.Fatalf("scenario %s primary before/candidate must import supporting module %s", scenario.ID, moduleRef)
			}
		}
		chains := map[string]map[string]bool{}
		for _, node := range graph.Nodes {
			containsSelf := false
			chainKey := strings.Join(node.Via, "\x00")
			if chains[chainKey] == nil {
				chains[chainKey] = map[string]bool{}
			}
			for _, siteID := range node.Via {
				linkedNode, exists := graphNodes[siteID]
				if !exists {
					t.Fatalf("scenario %s graph via references unknown site %s", scenario.ID, siteID)
				}
				chains[chainKey][linkedNode.Path] = true
				containsSelf = containsSelf || siteID == node.SiteID
			}
			if !containsSelf {
				t.Fatalf("scenario %s graph node %s is not a member of its declared chain", scenario.ID, node.SiteID)
			}
		}
		for chain, paths := range chains {
			if len(paths) < 2 {
				t.Fatalf("scenario %s graph chain %q does not cross an imported file boundary", scenario.ID, chain)
			}
		}
	}
	expected := readM1File(t, filepath.Join(root, scenario.ExpectedPath))
	if scenario.Route == "deterministic" {
		actual, err := applyM1Operation(pack.Operation, before)
		if err != nil {
			t.Fatalf("scenario %s apply operation: %v", scenario.ID, err)
		}
		if actual != expected {
			t.Fatalf("scenario %s candidate mismatch:\n--- actual\n%s\n--- expected\n%s", scenario.ID, actual, expected)
		}
		if strings.Count(before, pack.Operation.From) != 1 || strings.Contains(actual, pack.Operation.From) {
			t.Fatalf("scenario %s transform is not exact-once", scenario.ID)
		}
	} else {
		if _, err := applyM1Operation(pack.Operation, before); err == nil {
			t.Fatalf("scenario %s agent route is solvable by the deterministic transformer", scenario.ID)
		}
		if before == expected {
			t.Fatalf("scenario %s agent ground truth contains no change", scenario.ID)
		}
		recipes := map[string]bool{}
		from, to := m1MinimalReplacement(before, expected)
		recipes[from+"\x00"+to] = true
		for _, supporting := range scenario.SupportingFiles {
			supportingBefore := readM1File(t, filepath.Join(root, supporting.BeforePath))
			supportingExpected := readM1File(t, filepath.Join(root, supporting.ExpectedPath))
			if supportingBefore == supportingExpected {
				t.Fatalf("scenario %s supporting agent edit contains no change: %s", scenario.ID, supporting.BeforePath)
			}
			from, to := m1MinimalReplacement(supportingBefore, supportingExpected)
			recipes[from+"\x00"+to] = true
		}
		if len(recipes) < 2 {
			t.Fatalf("scenario %s does not require distinct repository-context edits", scenario.ID)
		}
	}
}

func m1MinimalReplacement(before string, expected string) (string, string) {
	prefix := 0
	for prefix < len(before) && prefix < len(expected) && before[prefix] == expected[prefix] {
		prefix++
	}
	suffix := 0
	for suffix < len(before)-prefix && suffix < len(expected)-prefix &&
		before[len(before)-1-suffix] == expected[len(expected)-1-suffix] {
		suffix++
	}
	return before[prefix : len(before)-suffix], expected[prefix : len(expected)-suffix]
}

func m1ScenarioRelativePath(t *testing.T, scenario m1Scenario, path string) string {
	t.Helper()
	marker := "/" + scenario.ID + "/"
	normalized := "/" + filepath.ToSlash(path)
	index := strings.Index(normalized, marker)
	if index < 0 {
		t.Fatalf("scenario %s path is outside its repository fixture: %s", scenario.ID, path)
	}
	return normalized[index+len(marker):]
}

type m1FakeRoute struct {
	Runner        string
	Network       bool
	Credential    bool
	ExternalWrite bool
}

func (adapter m1DeterministicFakeAgent) Run(ctx context.Context, request m1AgentRequest) (m1AgentResult, error) {
	if err := ctx.Err(); err != nil {
		return m1AgentResult{}, err
	}
	if err := adapter.Route.validate(); err != nil {
		return m1AgentResult{}, err
	}
	if request.ScenarioID == "" || request.PlanDigest == "" || request.Path == "" {
		return m1AgentResult{}, errors.New("fake agent request identity, plan, and path are required")
	}
	if request.Operation.Kind != "agent_required" || request.ExpectedCandidate == "" || request.ExpectedCandidate == request.Input {
		return m1AgentResult{}, errors.New("fake agent requires non-deterministic visible ground truth")
	}
	return m1AgentResult{
		SchemaVersion: "lumyn.agent-runner-result/v1",
		Status:        "candidate_ready",
		SessionMode:   "clean_ephemeral",
		Runner:        "deterministic_fake",
		Candidate:     request.ExpectedCandidate,
		ChangedPaths:  []string{request.Path},
	}, nil
}

func (route m1FakeRoute) validate() error {
	if route.Runner != "" && route.Runner != "deterministic_fake" {
		return fmt.Errorf("live agent runner %q is forbidden", route.Runner)
	}
	if route.Network || route.Credential || route.ExternalWrite {
		return errors.New("network, credentials, and external writes are forbidden")
	}
	return nil
}

type m1ObservedVerification struct {
	CommandsDigest         string
	BaselineCommandDigest  string
	BaselineOutputDigest   string
	BaselineExitCode       int
	CandidateCommandDigest string
	CandidateOutputDigest  string
	CandidateExitCode      int
	EvidenceDigest         string
	ToolchainDigest        string
	FixturesDigest         string
	CleanupEvidenceDigest  string
	RepositoryBase         string
	CandidateHead          string
}

type m1VerificationView struct {
	Root          string
	RepoDir       string
	Directories   []string
	BaseCommit    string
	CandidateHead string
	Removed       bool
}

func runM1WalkingSkeleton(t *testing.T, root string, scenario m1Scenario, observedAt time.Time) m1WalkingSkeleton {
	t.Helper()
	pack := readM1JSON[m1Pack](t, filepath.Join(root, scenario.PackPath))
	before := readM1File(t, filepath.Join(root, scenario.BeforePath))
	expected := readM1File(t, filepath.Join(root, scenario.ExpectedPath))
	candidate, err := applyM1Operation(pack.Operation, before)
	if err != nil {
		t.Fatal(err)
	}
	if candidate != expected {
		t.Fatal("walking-skeleton candidate failed independent exact-output verification")
	}
	observation := observeM1RepositoryVerification(t, root, scenario, candidate)
	chain := buildM1CanonicalContractChain(t, root, scenario, candidate, observation, observedAt)
	runID := "run:" + scenario.ID
	bindings := evidence.MigrationBindings{
		RunID:                           runID,
		InstallationDigest:              chain.InstallationDigest,
		AuthorizationInstallationDigest: chain.InstallationDigest,
		AuthorizationDigest:             chain.AuthorizationDigest,
		PlanAuthorizationDigest:         chain.AuthorizationDigest,
		PackDigest:                      chain.PackDigest,
		EventDigest:                     chain.EventDigest,
		EventContractDigest:             chain.PackDigest,
		AuthorizationEventDigest:        chain.EventDigest,
		ImpactDigest:                    chain.ImpactDigest,
		PlanImpactDigest:                chain.ImpactDigest,
		PlanDigest:                      chain.PlanDigest,
		CandidatePlanDigest:             chain.PlanDigest,
		CandidateDigest:                 chain.CandidateDigest,
		VerificationCandidateDigest:     chain.CandidateDigest,
		VerificationDigest:              chain.VerificationDigest,
		ExportCandidateDigest:           chain.CandidateDigest,
		ExportVerificationDigest:        chain.VerificationDigest,
		ExportDigest:                    chain.ExportDigest,
		ProjectionRunID:                 runID,
		ProjectionEventDigest:           chain.EventDigest,
		ProjectionInstallationDigest:    chain.InstallationDigest,
		ProjectionAuthorizationDigest:   chain.AuthorizationDigest,
		ProjectionPlanDigest:            chain.PlanDigest,
		ProjectionCandidateDigest:       chain.CandidateDigest,
		ProjectionVerificationDigest:    chain.VerificationDigest,
		ProjectionDeliveryDigest:        chain.ExportDigest,
		ProjectionEvidenceDigest:        chain.VerificationDigest,
	}
	if err := evidence.ValidateMigrationBindings(bindings); err != nil {
		t.Fatalf("walking-skeleton M2 binding contract: %v", err)
	}
	return m1WalkingSkeleton{
		SchemaVersion:       "lumyn.m1.walking-skeleton/v1",
		ScenarioID:          scenario.ID,
		EventDigest:         chain.EventDigest,
		InstallationDigest:  chain.InstallationDigest,
		ImpactDigest:        chain.ImpactDigest,
		AuthorizationDigest: chain.AuthorizationDigest,
		PlanDigest:          chain.PlanDigest,
		CandidateDigest:     chain.CandidateDigest,
		VerificationDigest:  chain.VerificationDigest,
		PRBundleDigest:      chain.ExportDigest,
		AgentRoute:          "disabled",
		ExternalWriteMode:   "manual_pr_bundle_only",
		ContractBinding:     "canonical_m2_schemas_and_authorization_validated",
		Status:              "static_verified_local_pr_bundle",
	}
}

func observeM1RepositoryVerification(t *testing.T, root string, scenario m1Scenario, candidate string) m1ObservedVerification {
	t.Helper()
	observedRuntime := requireM1PinnedRuntime(t)
	view := createM1ReadOnlyVerificationView(t, root, scenario)
	t.Cleanup(func() {
		if !view.Removed {
			removeM1VerificationView(t, &view, scenario.ID)
		}
	})

	baselineCommand, baselineCleanup := m1OfflineNPMTestCommandWithCleanup(t, view.RepoDir, "baseline")
	baselineOutput, baselineErr := baselineCommand.CombinedOutput()
	baselineCleanupDigest := baselineCleanup()
	baselineExitCode := m1CommandExitCode(baselineErr)
	if baselineErr == nil || baselineExitCode == 0 {
		t.Fatalf("walking-skeleton baseline unexpectedly passed:\n%s", baselineOutput)
	}
	baselineWant := scenario.ID + ": baseline target contract rejected"
	if !strings.Contains(string(baselineOutput), baselineWant) {
		t.Fatalf("walking-skeleton baseline lacks %q:\n%s", baselineWant, baselineOutput)
	}
	promoteM1VerificationViewCandidate(t, root, scenario, candidate, &view)

	candidateCommand, candidateCleanup := m1OfflineNPMTestCommandWithCleanup(t, view.RepoDir, "candidate")
	candidateOutput, candidateErr := candidateCommand.CombinedOutput()
	candidateCleanupDigest := candidateCleanup()
	if candidateErr != nil {
		t.Fatalf("walking-skeleton exact-candidate verification: %v\n%s", candidateErr, candidateOutput)
	}
	candidateWant := scenario.ID + ": exact candidate target contract verified"
	if !strings.Contains(string(candidateOutput), candidateWant) {
		t.Fatalf("walking-skeleton exact-candidate verification lacks %q:\n%s", candidateWant, candidateOutput)
	}

	baselineCommandDigest := m1VerificationCommandDigest(scenario, "baseline")
	candidateCommandDigest := m1VerificationCommandDigest(scenario, "candidate")
	commandsDigest := digestM1Canonical([]any{baselineCommandDigest, candidateCommandDigest})
	baselineOutputDigest := digestM1Bytes(baselineOutput)
	candidateOutputDigest := digestM1Bytes(candidateOutput)
	toolchainDigest := digestM1Canonical(observedRuntime)
	fixturesDigest := m1FixtureDigest(t, root, scenario, candidate)
	cleanupEvidenceDigest := digestM1Canonical(map[string]any{
		"baseline_command_cleanup_digest":  baselineCleanupDigest,
		"candidate_command_cleanup_digest": candidateCleanupDigest,
		"verification_view_cleanup_digest": removeM1VerificationView(t, &view, scenario.ID),
	})
	evidenceDigest := digestM1Canonical(map[string]any{
		"baseline_command_digest":  baselineCommandDigest,
		"baseline_exit_code":       baselineExitCode,
		"baseline_output_digest":   baselineOutputDigest,
		"candidate_command_digest": candidateCommandDigest,
		"candidate_digest":         digestM1Bytes([]byte(candidate)),
		"candidate_exit_code":      0,
		"candidate_output_digest":  candidateOutputDigest,
		"cleanup_evidence_digest":  cleanupEvidenceDigest,
		"commands_digest":          commandsDigest,
		"fixtures_digest":          fixturesDigest,
		"toolchain_digest":         toolchainDigest,
		"repository_base":          view.BaseCommit,
		"candidate_head":           view.CandidateHead,
	})
	return m1ObservedVerification{
		CommandsDigest:        commandsDigest,
		BaselineCommandDigest: baselineCommandDigest, BaselineOutputDigest: baselineOutputDigest,
		BaselineExitCode:       baselineExitCode,
		CandidateCommandDigest: candidateCommandDigest, CandidateOutputDigest: candidateOutputDigest,
		CandidateExitCode: 0, EvidenceDigest: evidenceDigest, ToolchainDigest: toolchainDigest,
		FixturesDigest: fixturesDigest, CleanupEvidenceDigest: cleanupEvidenceDigest,
		RepositoryBase: view.BaseCommit, CandidateHead: view.CandidateHead,
	}
}

func m1VerificationCommandDigest(scenario m1Scenario, target string) string {
	environment := []any{
		"CI=1",
		"HOME=<isolated>",
		"USERPROFILE=<isolated>",
		"TMPDIR=<isolated>",
		"TMP=<isolated>",
		"TEMP=<isolated>",
		"PATH=<pinned-node-npm-system-bins>",
		"NODE_NO_WARNINGS=1",
		"NO_COLOR=1",
		"npm_config_audit=false",
		"npm_config_cache=<isolated>",
		"npm_config_engine_strict=true",
		"npm_config_fund=false",
		"npm_config_globalconfig=<empty>",
		"npm_config_ignore_scripts=true",
		"npm_config_offline=true",
		"npm_config_registry=https://invalid.invalid",
		"npm_config_update_notifier=false",
		"npm_config_userconfig=<empty>",
		"LUMYN_M1_VERIFICATION_TARGET=" + target,
	}
	if target == "candidate" {
		environment = append(environment, "LUMYN_M1_CANDIDATE_PATH=src/client.ts")
	}
	return digestM1Canonical(map[string]any{
		"argv":                []any{"npm", "test", "--silent"},
		"cwd":                 "examples/consumer-repos/" + scenario.ID,
		"ambient_environment": "cleared",
		"module_execution":    "node_vm_no_host_globals_local_relative_imports_only",
		"verification_target": target,
		"environment":         environment,
	})
}

func m1CommandExitCode(err error) int {
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		return exitError.ExitCode()
	}
	if err == nil {
		return 0
	}
	return 255
}

func createM1ReadOnlyVerificationView(t *testing.T, root string, scenario m1Scenario) m1VerificationView {
	t.Helper()
	viewRoot, err := os.MkdirTemp("", "lumyn-m1-verification-")
	if err != nil {
		t.Fatal(err)
	}
	paths := []string{
		"examples/verification/run.mjs",
		filepath.ToSlash(filepath.Join("examples", "consumer-repos", scenario.ID, "package.json")),
		scenario.BeforePath,
	}
	for _, supporting := range scenario.SupportingFiles {
		paths = append(paths, supporting.BeforePath)
	}
	for _, relativePath := range paths {
		destination := filepath.Join(viewRoot, filepath.FromSlash(relativePath))
		if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(destination, readM1FileBytes(t, filepath.Join(root, relativePath)), 0o444); err != nil {
			t.Fatal(err)
		}
	}
	repoDir := filepath.Join(viewRoot, "examples", "consumer-repos", scenario.ID)
	runM1Git(t, viewRoot, repoDir, "init", "--quiet", "--initial-branch=main")
	runM1Git(t, viewRoot, repoDir, "add", "-A")
	runM1Git(t, viewRoot, repoDir, "-c", "commit.gpgsign=false", "-c", "user.name=Lumyn M1", "-c", "user.email=m1@example.invalid", "commit", "--quiet", "--no-verify", "-m", "baseline")
	baseCommit := runM1Git(t, viewRoot, repoDir, "rev-parse", "HEAD")
	view := m1VerificationView{Root: viewRoot, RepoDir: repoDir, BaseCommit: baseCommit}
	lockM1VerificationView(t, &view)
	return view
}

func promoteM1VerificationViewCandidate(t *testing.T, root string, scenario m1Scenario, candidate string, view *m1VerificationView) {
	t.Helper()
	unlockM1VerificationView(t, view)
	primaryPath := filepath.Join(view.RepoDir, filepath.FromSlash(m1ScenarioRelativePath(t, scenario, scenario.BeforePath)))
	if err := os.Chmod(primaryPath, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(primaryPath, []byte(candidate), 0o444); err != nil {
		t.Fatal(err)
	}
	if digestM1Bytes(readM1FileBytes(t, primaryPath)) != digestM1Bytes([]byte(candidate)) {
		t.Fatal("verification view did not preserve exact generated candidate bytes")
	}
	for _, supporting := range scenario.SupportingFiles {
		destination := filepath.Join(view.RepoDir, filepath.FromSlash(m1ScenarioRelativePath(t, scenario, supporting.BeforePath)))
		if err := os.Chmod(destination, 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(destination, readM1FileBytes(t, filepath.Join(root, supporting.ExpectedPath)), 0o444); err != nil {
			t.Fatal(err)
		}
	}
	runM1Git(t, view.Root, view.RepoDir, "add", "-A")
	runM1Git(t, view.Root, view.RepoDir, "-c", "commit.gpgsign=false", "-c", "user.name=Lumyn M1", "-c", "user.email=m1@example.invalid", "commit", "--quiet", "--no-verify", "-m", "candidate")
	view.CandidateHead = runM1Git(t, view.Root, view.RepoDir, "rev-parse", "HEAD")
	if view.CandidateHead == view.BaseCommit {
		t.Fatal("candidate Git head did not advance from the baseline commit")
	}
	runM1Git(t, view.Root, view.RepoDir, "diff", "--quiet", "HEAD", "--")
	runM1Git(t, view.Root, view.RepoDir, "cat-file", "-e", view.CandidateHead+"^{commit}")
	lockM1VerificationView(t, view)
}

func runM1Git(t *testing.T, viewRoot string, repoDir string, arguments ...string) string {
	t.Helper()
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("resolve git for M1 exact-head fixture: %v", err)
	}
	home := filepath.Join(viewRoot, "git-home")
	if err := os.MkdirAll(home, 0o700); err != nil {
		t.Fatal(err)
	}
	command := exec.Command(gitPath, arguments...)
	command.Dir = repoDir
	command.Env = []string{
		"HOME=" + home,
		"PATH=" + filepath.Dir(gitPath) + string(os.PathListSeparator) + "/usr/bin" + string(os.PathListSeparator) + "/bin",
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_AUTHOR_DATE=2026-07-28T12:00:00Z",
		"GIT_COMMITTER_DATE=2026-07-28T12:00:00Z",
		"LC_ALL=C",
	}
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("M1 exact-head git %v: %v\n%s", arguments, err, output)
	}
	return strings.TrimSpace(string(output))
}

func unlockM1VerificationView(t *testing.T, view *m1VerificationView) {
	t.Helper()
	for _, directory := range view.Directories {
		if err := os.Chmod(directory, 0o755); err != nil && !os.IsNotExist(err) {
			t.Fatal(err)
		}
	}
}

func lockM1VerificationView(t *testing.T, view *m1VerificationView) {
	t.Helper()
	var directories []string
	if err := filepath.WalkDir(view.Root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			directories = append(directories, path)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	for index := len(directories) - 1; index >= 0; index-- {
		if err := os.Chmod(directories[index], 0o555); err != nil {
			t.Fatal(err)
		}
	}
	view.Directories = directories
}

func removeM1VerificationView(t *testing.T, view *m1VerificationView, scenarioID string) string {
	t.Helper()
	if view.Removed {
		return digestM1Canonical(map[string]any{"scenario_id": scenarioID, "verification_view_removed": true})
	}
	for _, directory := range view.Directories {
		if err := os.Chmod(directory, 0o755); err != nil && !os.IsNotExist(err) {
			t.Fatalf("make verification view removable: %v", err)
		}
	}
	if err := os.RemoveAll(view.Root); err != nil {
		t.Fatalf("remove verification view: %v", err)
	}
	if _, err := os.Stat(view.Root); !os.IsNotExist(err) {
		t.Fatalf("verification view still exists after cleanup: %v", err)
	}
	view.Removed = true
	return digestM1Canonical(map[string]any{"scenario_id": scenarioID, "verification_view_removed": true})
}

func m1FixtureDigest(t *testing.T, root string, scenario m1Scenario, candidate string) string {
	t.Helper()
	supporting := make([]any, 0, len(scenario.SupportingFiles))
	for _, file := range scenario.SupportingFiles {
		supporting = append(supporting, map[string]any{
			"before_path": file.BeforePath, "before_digest": file.BeforeDigest,
			"expected_path": file.ExpectedPath, "expected_digest": file.ExpectedDigest,
		})
	}
	return digestM1Canonical(map[string]any{
		"scenario_id":                scenario.ID,
		"declared_digests":           scenario.Digests,
		"supporting":                 supporting,
		"package_json_digest":        digestM1Bytes(readM1FileBytes(t, filepath.Join(root, "examples", "consumer-repos", scenario.ID, "package.json"))),
		"verification_runner_digest": digestM1Bytes(readM1FileBytes(t, filepath.Join(root, "examples/verification/run.mjs"))),
		"generated_candidate_digest": digestM1Bytes([]byte(candidate)),
	})
}

type m1CanonicalContractChain struct {
	PackDigest          string
	EventDigest         string
	InstallationDigest  string
	ImpactDigest        string
	AuthorizationDigest string
	PlanDigest          string
	CandidateDigest     string
	PatchDigest         string
	VerificationDigest  string
	ExportDigest        string
	Artifacts           map[string]map[string]any `json:"-"`
}

func buildM1CanonicalContractChain(
	t *testing.T,
	root string,
	scenario m1Scenario,
	candidate string,
	observation m1ObservedVerification,
	observedAt time.Time,
) m1CanonicalContractChain {
	t.Helper()
	objects := map[string]map[string]any{}
	for _, name := range []string{
		"provider-change-event", "migration-pack", "consumer-installation", "event-authorization",
		"impact-report", "migration-plan", "candidate-manifest", "migration-verification", "export-result",
	} {
		objects[name] = readM1JSONObject(t, filepath.Join(root, "tests", "fixtures", "contracts", name, "valid.json"))
	}

	packDigest := scenario.Digests["pack"]
	eventDigest := digestM1Canonical(map[string]any{"scenario_id": scenario.ID, "pack_digest": packDigest, "sequence": 1})
	installationDigest := digestM1Canonical(map[string]any{"scenario_id": scenario.ID, "action": "prepare_patch", "agent": "disabled"})
	impactDigest := scenario.Digests["integration_graph"]
	authorizationDigest := digestM1Canonical(map[string]any{
		"scenario_id": scenario.ID, "event_digest": eventDigest, "installation_digest": installationDigest,
	})
	planDigest := digestM1Canonical(map[string]any{
		"scenario_id": scenario.ID, "authorization_digest": authorizationDigest, "impact_digest": impactDigest,
	})
	patchDigest := digestM1Bytes([]byte(candidate))
	candidateDigest := digestM1Canonical(map[string]any{
		"scenario_id": scenario.ID, "plan_digest": planDigest, "patch_digest": patchDigest,
	})
	candidateHead := observation.CandidateHead
	repositoryBindingDigest := digestM1Canonical(map[string]any{
		"repository_id": "github.acme.checkout-service", "package_root": ".", "base_commit": observation.RepositoryBase,
	})
	verificationDigest := digestM1Canonical(map[string]any{
		"scenario_id": scenario.ID, "plan_digest": planDigest, "candidate_digest": candidateDigest,
		"patch_digest": patchDigest, "evidence_digest": observation.EvidenceDigest, "outcome": "verified",
	})
	exportDigest := digestM1Canonical(map[string]any{
		"scenario_id": scenario.ID, "candidate_digest": candidateDigest, "verification_digest": verificationDigest, "mode": "pr_bundle",
	})

	eventID := "event.synthetic." + strings.ReplaceAll(scenario.ID, "-", "_") + ".0001"
	packID := "migration_pack.synthetic." + strings.ReplaceAll(scenario.ID, "-", "_")
	planID := "migration_plan.synthetic." + strings.ReplaceAll(scenario.ID, "-", "_")
	event := objects["provider-change-event"]
	event["event_id"], event["event_digest"] = eventID, eventDigest
	m1Object(t, event, "contract_delivery")["migration_pack_id"] = packID
	m1Object(t, event, "contract_delivery")["retrieved_bytes_digest"] = packDigest

	pack := objects["migration-pack"]
	pack["migration_pack_id"], pack["contract_digest"] = packID, packDigest
	m1Object(t, pack, "provider_confirmation")["confirmed_contract_digest"] = packDigest
	changeItemID := "change." + strings.ReplaceAll(scenario.ID, "-", "_")
	semanticChange := m1ObjectAt(t, pack, "semantic_changes", 0)
	semanticChange["change_id"] = changeItemID
	semanticChange["kind"] = "rename"
	semanticChange["source_symbol"] = "client.charges.create(request)"
	semanticChange["target_symbol"] = "client.paymentIntents.create(request)"
	semanticChange["intent"] = "Rename the supported SDK operation without changing the request payload."
	semanticChange["applicability"] = "Statically resolved direct SDK calls in the authorized package root."
	semanticChange["mapping"] = map[string]any{
		"mode": "exact", "declarative_steps": []any{"Apply the exact-once operation rename."},
	}

	installation := objects["consumer-installation"]
	installation["installation_digest"] = installationDigest
	installation["action_ceiling"] = "prepare_patch"
	installation["agent_execution_policy"] = map[string]any{"state": "disabled"}
	commandID := "command.npm_test_silent"
	commands := installation["commands"].([]any)
	command := commands[0].(map[string]any)
	command["command_id"] = commandID
	command["program"] = "npm"
	command["arguments"] = []any{"test", "--silent"}
	command["working_directory"] = "."
	command["network_mode"] = "offline"
	command["lifecycle_scripts_allowed"] = false
	capabilities := m1Object(t, installation, "capability_ceiling")
	for _, key := range []string{
		"agent_runner_network", "agent_runner_credential", "model_request_disclosure", "model_network", "model_credential",
		"package_registry_read", "dependency_lifecycle_scripts", "sandbox_request_disclosure", "sandbox_network", "sandbox_credential",
		"github_branch_write", "github_pr_write", "provider_reporting", "artifact_retention", "artifact_deletion",
	} {
		capabilities[key] = false
	}
	githubPolicy := m1Object(t, installation, "github_token_issuance_policy")
	githubPolicy["enabled"], githubPolicy["broker_id"] = false, nil
	providerReporting := m1Object(t, installation, "provider_reporting")
	providerReporting["enabled"] = false
	providerReporting["allowed_fields"] = []any{}
	providerReporting["raw_consumer_evidence_allowed"] = false
	providerReporting["consent_required"] = true
	pathPolicyDigest := digestM1Canonical(installation["path_policy"])
	commandPolicyDigest := digestM1Canonical(installation["commands"])
	modelDataPolicyDigest := digestM1Canonical(installation["model_data_policy"])
	budgetPolicyDigest := digestM1Canonical(installation["budgets"])
	githubPolicyDigest := digestM1Canonical(installation["github_token_issuance_policy"])
	providerReportingPolicyDigest := digestM1Canonical(installation["provider_reporting"])
	verificationConfigurationDigest := digestM1Canonical(map[string]any{
		"commands_digest": observation.CommandsDigest,
		"fixtures_digest": observation.FixturesDigest,
		"offline":         true,
	})

	authorizationObject := objects["event-authorization"]
	authorizationObject["authorization_digest"] = authorizationDigest
	m1Object(t, authorizationObject, "installation_binding")["artifact_digest"] = installationDigest
	m1Object(t, authorizationObject, "event_binding")["artifact_id"] = eventID
	m1Object(t, authorizationObject, "event_binding")["artifact_digest"] = eventDigest
	m1Object(t, authorizationObject, "migration_pack_binding")["artifact_id"] = packID
	m1Object(t, authorizationObject, "migration_pack_binding")["artifact_digest"] = packDigest
	m1Object(t, authorizationObject, "plan_binding")["artifact_id"] = planID
	m1Object(t, authorizationObject, "plan_binding")["artifact_digest"] = planDigest
	m1Object(t, authorizationObject, "repository_binding")["base_commit"] = observation.RepositoryBase
	authorizationObject["selected_action"] = "prepare_patch"
	authorizationObject["generation_mode"] = "deterministic"
	authorizationObject["agent_execution_policy"] = map[string]any{"state": "disabled"}
	authorizationObject["selected_capabilities"] = []any{
		"customer_repo_read", "customer_repo_write", "command_execution",
	}
	authorizationScope := m1Object(t, authorizationObject, "scope")
	authorizationScope["command_ids"] = []any{commandID}
	authorizationScope["path_policy_digest"] = pathPolicyDigest
	authorizationScope["command_policy_digest"] = commandPolicyDigest
	authorizationScope["model_data_policy_digest"] = modelDataPolicyDigest
	authorizationScope["budget_policy_digest"] = budgetPolicyDigest
	authorizationScope["github_policy_digest"] = githubPolicyDigest
	authorizationScope["provider_reporting_policy_digest"] = providerReportingPolicyDigest
	m1Object(t, authorizationObject, "verification_requirement")["configuration_digest"] = verificationConfigurationDigest
	m1Object(t, authorizationObject, "derivation")["installation_policy_digest"] = installationDigest

	impact := objects["impact-report"]
	impact["event_id"], impact["contract_digest"] = eventID, packDigest
	impact["installation_id"] = installation["installation_id"]
	impact["authorization_snapshot_digest"] = authorizationDigest
	impact["repository_base"] = observation.RepositoryBase
	impact["integration_graph_id"], impact["integration_graph_digest"] = "graph."+strings.ReplaceAll(scenario.ID, "-", "_"), impactDigest
	impact["report_digest"] = impactDigest
	finding := m1ObjectAt(t, impact, "findings", 0)
	finding["change_item_id"], finding["path"] = changeItemID, "src/client.ts"

	plan := objects["migration-plan"]
	plan["plan_id"], plan["plan_digest"] = planID, planDigest
	plan["event_id"], plan["contract_digest"] = eventID, packDigest
	plan["installation_id"], plan["authorization_snapshot_digest"] = installation["installation_id"], authorizationDigest
	plan["repository_base"] = observation.RepositoryBase
	plan["integration_graph_digest"], plan["impact_report_digest"] = impactDigest, impactDigest
	planItem := m1ObjectAt(t, plan, "items", 0)
	planItem["change_item_id"], planItem["route"] = changeItemID, "deterministic"
	planItem["planned_files"] = []any{"src/client.ts"}
	delete(planItem, "agent_route")
	planItem["deterministic_recipe"] = map[string]any{
		"recipe_id": "recipe." + strings.ReplaceAll(scenario.ID, "-", "_"), "recipe_version": "1.0.0", "recipe_digest": packDigest,
	}
	m1Object(t, plan, "delivery")["mode"] = "pr_bundle"
	m1Object(t, plan, "approval")["plan_digest"] = planDigest

	candidateManifest := objects["candidate-manifest"]
	candidateManifest["candidate_digest"], candidateManifest["patch_digest"] = candidateDigest, patchDigest
	candidateManifest["candidate_head"] = candidateHead
	candidateManifest["repository_base"] = observation.RepositoryBase
	m1Object(t, candidateManifest, "independent_verification")["candidate_head"] = candidateHead
	candidateManifest["event_id"], candidateManifest["contract_digest"] = eventID, packDigest
	candidateManifest["installation_id"], candidateManifest["authorization_snapshot_digest"] = installation["installation_id"], authorizationDigest
	candidateManifest["plan_id"], candidateManifest["plan_digest"] = planID, planDigest
	generation := m1ObjectAt(t, candidateManifest, "generation_records", 0)
	generation["route"] = "deterministic"
	delete(generation, "agent")
	generation["recipe"] = map[string]any{
		"recipe_id": "recipe." + strings.ReplaceAll(scenario.ID, "-", "_"), "recipe_version": "1.0.0", "recipe_digest": packDigest,
	}
	change := m1ObjectAt(t, candidateManifest, "changes", 0)
	change["change_item_id"], change["route"], change["path"], change["edit_digest"] = changeItemID, "deterministic", "src/client.ts", patchDigest

	verification := objects["migration-verification"]
	verification["provider_change_event_id"], verification["provider_change_contract_digest"] = eventID, packDigest
	verification["consumer_installation_id"], verification["event_authorization_digest"] = installation["installation_id"], authorizationDigest
	verification["migration_plan_digest"], verification["artifact_digest"] = planDigest, verificationDigest
	verificationCandidate := m1Object(t, verification, "candidate_binding")
	verificationCandidate["candidate_manifest_digest"], verificationCandidate["patch_digest"] = candidateDigest, patchDigest
	verificationCandidate["candidate_head"] = candidateHead
	verificationCandidate["repository_binding_digest"] = repositoryBindingDigest
	verificationCandidate["repository_base_commit"] = observation.RepositoryBase
	provenance := m1Object(t, verification, "generation_provenance")
	provenance["generation_mode"], provenance["generation_evidence_digest"] = "deterministic", candidateDigest
	delete(provenance, "agent_attempt_id")
	delete(provenance, "agent_attempt_digest")
	executionBinding := m1Object(t, verification, "execution_binding")
	executionBinding["agent_execution_policy"], executionBinding["generation_mode"] = "disabled", "deterministic"
	executionBinding["route_digest"], executionBinding["authorization_snapshot_digest"] = packDigest, authorizationDigest
	for _, key := range []string{"agent_runner_binding_digest", "model_route_binding_digest", "execution_funding_mode", "credential_owner", "usage_billing_owner"} {
		delete(executionBinding, key)
	}
	frozenInputs := m1Object(t, verification, "frozen_inputs")
	frozenInputs["commands_digest"] = observation.CommandsDigest
	frozenInputs["toolchain_digest"] = observation.ToolchainDigest
	frozenInputs["fixtures_digest"] = observation.FixturesDigest
	frozenInputs["environment_digest"] = digestM1Canonical(map[string]any{
		"network": "offline", "lifecycle_scripts": false, "engine_strict": true,
	})
	verification["baseline_results"] = []any{map[string]any{
		"command_id": "red-before-target-contract", "command_digest": observation.BaselineCommandDigest,
		"exit_code": observation.BaselineExitCode, "result": "failed", "output_digest": observation.BaselineOutputDigest, "pre_existing_failure": true,
	}}
	verification["candidate_results"] = []any{map[string]any{
		"command_id": "green-candidate-target-contract", "command_digest": observation.CandidateCommandDigest,
		"exit_code": observation.CandidateExitCode, "result": "passed", "output_digest": observation.CandidateOutputDigest, "pre_existing_failure": false,
	}}
	verification["verification_label"], verification["outcome"] = "static_verified", "verified"
	verification["observed_evidence_refs"] = []any{map[string]any{
		"evidence_kind": "repository_commands", "evidence_digest": observation.EvidenceDigest,
		"candidate_head": candidateHead,
	}}
	verification["verified_at"] = observedAt.Format(time.RFC3339)
	m1Object(t, verification, "cleanup")["evidence_digest"] = observation.CleanupEvidenceDigest
	verification["residual_risks"] = []any{"Public synthetic fixture evidence is not real-consumer proof."}

	export := objects["export-result"]
	export["provider_change_event_id"], export["provider_change_contract_digest"] = eventID, packDigest
	export["consumer_installation_id"], export["event_authorization_digest"] = installation["installation_id"], authorizationDigest
	export["migration_plan_digest"], export["artifact_digest"] = planDigest, exportDigest
	exportCandidate := m1Object(t, export, "candidate_binding")
	exportCandidate["candidate_manifest_digest"], exportCandidate["patch_digest"] = candidateDigest, patchDigest
	exportCandidate["candidate_head"] = candidateHead
	exportCandidate["repository_binding_digest"] = repositoryBindingDigest
	exportCandidate["repository_base_commit"] = observation.RepositoryBase
	exportVerification := m1Object(t, export, "verification_binding")
	exportVerification["verification_digest"], exportVerification["candidate_manifest_digest"], exportVerification["patch_digest"] = verificationDigest, candidateDigest, patchDigest
	exportVerification["candidate_head"] = candidateHead
	exportVerification["verification_label"], exportVerification["fresh"] = "static_verified", true
	delivery := m1Object(t, export, "delivery")
	delivery["mode"], delivery["state"], delivery["delivery_authority"] = "pr_bundle", "pr_bundle_ready", "manual_handoff"
	delivery["github_actions"], delivery["counts_as_automated_pr_delivery"] = []any{}, false
	delivery["artifact_ref"] = "consumer-private://exports/" + scenario.ID + "/pr-bundle.json"
	delivery["source_branch"] = "lumyn/" + scenario.ID
	delete(delivery, "github_authorization")
	delete(delivery, "draft")
	export["evidence_bundle_digest"] = observation.EvidenceDigest
	export["exported_at"] = observedAt.Format(time.RFC3339)

	// Finalize exact artifact identities only after every semantic field is
	// populated. Self-digest fields are normalized, and the authorization's
	// later plan binding is normalized to break the deliberate authorization ↔
	// plan cycle; the populated cross-binding is validated separately below.
	packDigest = digestM1CanonicalArtifact(t, "migration-pack", pack)
	pack["contract_digest"] = packDigest
	m1Object(t, pack, "provider_confirmation")["confirmed_contract_digest"] = packDigest

	m1Object(t, event, "contract_delivery")["retrieved_bytes_digest"] = packDigest
	eventDigest = digestM1CanonicalArtifact(t, "provider-change-event", event)
	event["event_digest"] = eventDigest

	installationDigest = digestM1CanonicalArtifact(t, "consumer-installation", installation)
	installation["installation_digest"] = installationDigest

	m1Object(t, authorizationObject, "installation_binding")["artifact_digest"] = installationDigest
	m1Object(t, authorizationObject, "event_binding")["artifact_digest"] = eventDigest
	m1Object(t, authorizationObject, "migration_pack_binding")["artifact_digest"] = packDigest
	m1Object(t, authorizationObject, "derivation")["installation_policy_digest"] = installationDigest
	authorizationDigest = digestM1CanonicalArtifact(t, "event-authorization", authorizationObject)
	authorizationObject["authorization_digest"] = authorizationDigest

	impact["contract_digest"] = packDigest
	impact["authorization_snapshot_digest"] = authorizationDigest
	impactDigest = digestM1CanonicalArtifact(t, "impact-report", impact)
	impact["report_digest"] = impactDigest

	plan["contract_digest"] = packDigest
	plan["authorization_snapshot_digest"] = authorizationDigest
	plan["impact_report_digest"] = impactDigest
	m1Object(t, planItem, "deterministic_recipe")["recipe_digest"] = packDigest
	planDigest = digestM1CanonicalArtifact(t, "migration-plan", plan)
	plan["plan_digest"] = planDigest
	m1Object(t, plan, "approval")["plan_digest"] = planDigest
	m1Object(t, authorizationObject, "plan_binding")["artifact_digest"] = planDigest
	if rebound := digestM1CanonicalArtifact(t, "event-authorization", authorizationObject); rebound != authorizationDigest {
		t.Fatalf("later plan binding changed normalized authorization digest: %s != %s", rebound, authorizationDigest)
	}

	candidateManifest["contract_digest"] = packDigest
	candidateManifest["authorization_snapshot_digest"] = authorizationDigest
	candidateManifest["plan_digest"] = planDigest
	m1Object(t, generation, "recipe")["recipe_digest"] = packDigest
	candidateManifest["candidate_head"] = candidateHead
	m1Object(t, candidateManifest, "independent_verification")["candidate_head"] = candidateHead
	candidateDigest = digestM1CanonicalArtifact(t, "candidate-manifest", candidateManifest)
	candidateManifest["candidate_digest"] = candidateDigest

	verification["provider_change_contract_digest"] = packDigest
	verification["event_authorization_digest"] = authorizationDigest
	verification["migration_plan_digest"] = planDigest
	verificationCandidate["candidate_manifest_digest"] = candidateDigest
	verificationCandidate["candidate_head"] = candidateHead
	provenance["generation_evidence_digest"] = candidateDigest
	executionBinding["route_digest"] = packDigest
	executionBinding["authorization_snapshot_digest"] = authorizationDigest
	frozenInputs["verification_config_digest"] = verificationConfigurationDigest
	frozenInputs["isolation_backend_id"] = "m1_unqualified_developer_test_process"
	frozenInputs["isolation_backend_version"] = "1"
	frozenInputs["isolation_configuration_digest"] = digestM1Canonical(map[string]any{
		"candidate_vm": "node_vm_not_security_boundary", "environment": "allowlisted", "network_enforcement": "not_host_enforced",
	})
	frozenInputs["isolation_qualification_digest"] = digestM1Canonical(map[string]any{
		"status": "unqualified_for_product_repository_execution",
	})
	frozenInputs["mounts_digest"] = digestM1Canonical(map[string]any{
		"filesystem_view": "chmod_read_only", "host_enforced_mounts": false,
	})
	frozenInputs["resource_limits_digest"] = digestM1Canonical(map[string]any{
		"hard_host_resource_quotas": false,
	})
	verification["verification_label"], verification["outcome"] = "static_verified", "verified"
	m1ObjectAt(t, verification, "observed_evidence_refs", 0)["evidence_kind"] = "static_analysis"
	m1ObjectAt(t, verification, "observed_evidence_refs", 0)["candidate_head"] = candidateHead
	verification["residual_risks"] = []any{
		"Public synthetic fixture evidence is not real-consumer proof.",
		"The developer test process is not a qualified host-isolation backend and cannot earn repo_verified or authorize product repository command execution.",
	}
	verificationDigest = digestM1CanonicalArtifact(t, "migration-verification", verification)
	verification["artifact_digest"] = verificationDigest

	export["provider_change_contract_digest"] = packDigest
	export["event_authorization_digest"] = authorizationDigest
	export["migration_plan_digest"] = planDigest
	exportCandidate["candidate_manifest_digest"] = candidateDigest
	exportCandidate["candidate_head"] = candidateHead
	exportVerification["verification_digest"] = verificationDigest
	exportVerification["candidate_manifest_digest"] = candidateDigest
	exportVerification["candidate_head"] = candidateHead
	exportVerification["verification_label"] = "static_verified"
	exportDigest = digestM1CanonicalArtifact(t, "export-result", export)
	export["artifact_digest"] = exportDigest

	for name, value := range objects {
		validateM1Schema(t, root, name+".schema.json", value)
	}
	if candidateManifest["candidate_head"] != observation.CandidateHead ||
		m1Object(t, candidateManifest, "independent_verification")["candidate_head"] != observation.CandidateHead ||
		verificationCandidate["candidate_head"] != observation.CandidateHead ||
		exportCandidate["candidate_head"] != observation.CandidateHead ||
		exportVerification["candidate_head"] != observation.CandidateHead {
		t.Fatal("canonical artifacts are not all bound to the verified Git candidate head")
	}
	if candidateManifest["repository_base"] != observation.RepositoryBase ||
		verificationCandidate["repository_base_commit"] != observation.RepositoryBase ||
		exportCandidate["repository_base_commit"] != observation.RepositoryBase {
		t.Fatal("canonical artifacts are not all bound to the verified Git repository base")
	}
	installationBytes := mustM1Marshal(t, installation)
	authorizationBytes := mustM1Marshal(t, authorizationObject)
	expectedBindings := authorization.ExpectedAuthorizationBindings{
		Event:         authorization.ArtifactBinding{ArtifactID: eventID, ArtifactDigest: eventDigest},
		MigrationPack: authorization.ArtifactBinding{ArtifactID: packID, ArtifactDigest: packDigest},
		Repository: authorization.RepositoryBinding{
			RepositoryID: "github.acme.checkout-service", PackageRoot: ".", BaseCommit: observation.RepositoryBase,
		},
		Plan:              authorization.ArtifactBinding{ArtifactID: planID, ArtifactDigest: planDigest},
		ExecutionManifest: authorization.ArtifactBinding{ArtifactID: "execution_manifest.acme.payments_node_v5.0001", ArtifactDigest: "sha256:" + strings.Repeat("f", 64)},
		Policies: authorization.PolicyDigestBindings{
			PathPolicy: pathPolicyDigest, CommandPolicy: commandPolicyDigest,
			ModelDataPolicy: modelDataPolicyDigest, BudgetPolicy: budgetPolicyDigest,
			GitHubPolicy: githubPolicyDigest, ProviderReportingPolicy: providerReportingPolicyDigest,
		},
		VerificationConfigurationDigest: verificationConfigurationDigest,
		CredentialIssuancePolicyDigest:  "sha256:" + strings.Repeat("8", 64),
		ExcludedPaths:                   []string{".git", "node_modules"},
		ProviderReportingFields:         []string{},
		Budgets: authorization.Budgets{
			MaxChangedFiles: 12, MaxDiffLines: 500, MaxDiffBytes: 64000, MaxTurns: 20,
			MaxAttempts: 2, MaxTokens: 50000, MaxCostCents: 2500, MaxDurationSeconds: 1800,
		},
	}
	snapshot, err := authorization.ValidateSnapshotFromPersistedInstallation(
		installationBytes, authorizationBytes, expectedBindings, time.Date(2026, 7, 25, 17, 30, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("canonical deterministic authorization: %v", err)
	}
	if snapshot.Route != authorization.Deterministic || snapshot.ActionMode != authorization.PreparePatch || snapshot.AgentRouteDigest != "" ||
		snapshot.Scope.AgentRunnerNetwork || snapshot.Scope.AgentRunnerCredential || snapshot.Scope.ModelRequestDisclosure || snapshot.Scope.RemoteBranch || snapshot.Scope.DraftPR {
		t.Fatalf("canonical authorization widened deterministic local route: %+v", snapshot)
	}

	return m1CanonicalContractChain{
		PackDigest: packDigest, EventDigest: eventDigest, InstallationDigest: installationDigest,
		ImpactDigest: impactDigest, AuthorizationDigest: authorizationDigest, PlanDigest: planDigest,
		CandidateDigest: candidateDigest, PatchDigest: patchDigest,
		VerificationDigest: verificationDigest, ExportDigest: exportDigest, Artifacts: objects,
	}
}

func readM1JSONObject(t *testing.T, path string) map[string]any {
	t.Helper()
	payload := readM1FileBytes(t, path)
	var value map[string]any
	if err := json.Unmarshal(payload, &value); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	return value
}

func m1Object(t *testing.T, parent map[string]any, key string) map[string]any {
	t.Helper()
	value, ok := parent[key].(map[string]any)
	if !ok {
		t.Fatalf("%s is %T, want object", key, parent[key])
	}
	return value
}

func m1ObjectAt(t *testing.T, parent map[string]any, key string, index int) map[string]any {
	t.Helper()
	values, ok := parent[key].([]any)
	if !ok || index < 0 || index >= len(values) {
		t.Fatalf("%s is %T with unavailable index %d", key, parent[key], index)
	}
	value, ok := values[index].(map[string]any)
	if !ok {
		t.Fatalf("%s[%d] is %T, want object", key, index, values[index])
	}
	return value
}

func mustM1Marshal(t *testing.T, value any) []byte {
	t.Helper()
	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return payload
}

func validateM1Schema(t *testing.T, root, schemaName string, value any) {
	t.Helper()
	schema, err := jsonschema.Compile(filepath.Join(root, "schemas", schemaName))
	if err != nil {
		t.Fatalf("compile %s: %v", schemaName, err)
	}
	if err := schema.Validate(value); err != nil {
		t.Fatalf("validate canonical %s: %v", schemaName, err)
	}
}

func validateM1ReplayBinding(currentDigest, eventDigest string) error {
	if currentDigest != eventDigest {
		return errors.New("provider event migration-pack digest mismatch")
	}
	return nil
}

func applyM1Operation(operation m1Operation, input string) (string, error) {
	if operation.Kind != "replace_exact_once" {
		return "", fmt.Errorf("unsupported operation %q", operation.Kind)
	}
	if strings.Count(input, operation.From) != 1 {
		return "", fmt.Errorf("operation source occurrence count = %d, want 1", strings.Count(input, operation.From))
	}
	return strings.Replace(input, operation.From, operation.To, 1), nil
}

func readM1JSON[T any](t *testing.T, path string) T {
	t.Helper()
	payload := readM1FileBytes(t, path)
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()
	var value T
	if err := decoder.Decode(&value); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	if decoder.Decode(&struct{}{}) == nil {
		t.Fatalf("decode %s: trailing JSON value", path)
	}
	return value
}

func readM1File(t *testing.T, path string) string {
	t.Helper()
	return string(readM1FileBytes(t, path))
}

func readM1FileBytes(t *testing.T, path string) []byte {
	t.Helper()
	payload, err := os.ReadFile(path)
	if err != nil {
		display := filepath.ToSlash(path)
		if relative, relativeErr := filepath.Rel(repoRoot(t), path); relativeErr == nil && !strings.HasPrefix(relative, "..") {
			display = filepath.ToSlash(relative)
		}
		reason := "read failed"
		if os.IsNotExist(err) {
			reason = "file does not exist"
		} else if os.IsPermission(err) {
			reason = "permission denied"
		}
		t.Fatalf("read %s: %s", display, reason)
	}
	return payload
}

func digestM1Bytes(payload []byte) string {
	digest := sha256.Sum256(payload)
	return "sha256:" + hex.EncodeToString(digest[:])
}

func digestM1Canonical(value any) string {
	payload, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return digestM1Bytes(payload)
}

type m1DigestNormalization struct {
	Path  []string
	Value any
}

func digestM1CanonicalArtifact(t *testing.T, artifactName string, value map[string]any) string {
	t.Helper()
	digestPlaceholder := "sha256:" + strings.Repeat("0", 64)
	normalizations := map[string][]m1DigestNormalization{
		"provider-change-event": {{Path: []string{"event_digest"}, Value: digestPlaceholder}},
		"migration-pack": {
			{Path: []string{"contract_digest"}, Value: digestPlaceholder},
			{Path: []string{"provider_confirmation", "confirmed_contract_digest"}, Value: digestPlaceholder},
		},
		"consumer-installation": {{Path: []string{"installation_digest"}, Value: digestPlaceholder}},
		"event-authorization": {
			{Path: []string{"authorization_digest"}, Value: digestPlaceholder},
			// The plan binds this authorization, so normalizing the later-artifact
			// digest breaks the deliberate cross-artifact cycle. MigrationBindings
			// verifies the populated plan binding separately.
			{Path: []string{"plan_binding", "artifact_digest"}, Value: digestPlaceholder},
		},
		"impact-report": {{Path: []string{"report_digest"}, Value: digestPlaceholder}},
		"migration-plan": {
			{Path: []string{"plan_digest"}, Value: digestPlaceholder},
			{Path: []string{"approval", "plan_digest"}, Value: digestPlaceholder},
		},
		"candidate-manifest": {
			{Path: []string{"candidate_digest"}, Value: digestPlaceholder},
		},
		"migration-verification": {{Path: []string{"artifact_digest"}, Value: digestPlaceholder}},
		"export-result":          {{Path: []string{"artifact_digest"}, Value: digestPlaceholder}},
	}
	rules, ok := normalizations[artifactName]
	if !ok {
		t.Fatalf("missing M1 canonical digest normalization for %s", artifactName)
	}
	clone := cloneM1JSONObject(t, value)
	for _, rule := range rules {
		current := clone
		for _, segment := range rule.Path[:len(rule.Path)-1] {
			next, ok := current[segment].(map[string]any)
			if !ok {
				t.Fatalf("M1 canonical digest normalization path is invalid for %s: %v", artifactName, rule.Path)
			}
			current = next
		}
		current[rule.Path[len(rule.Path)-1]] = rule.Value
	}
	return digestM1Canonical(clone)
}

func cloneM1JSONObject(t *testing.T, value map[string]any) map[string]any {
	t.Helper()
	payload := mustM1Marshal(t, value)
	var clone map[string]any
	if err := json.Unmarshal(payload, &clone); err != nil {
		t.Fatalf("clone M1 JSON object: %v", err)
	}
	return clone
}

func containsM1JSONKey(value any, target string) bool {
	switch typed := value.(type) {
	case map[string]any:
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			if key == target || containsM1JSONKey(typed[key], target) {
				return true
			}
		}
	case []any:
		for _, item := range typed {
			if containsM1JSONKey(item, target) {
				return true
			}
		}
	}
	return false
}
