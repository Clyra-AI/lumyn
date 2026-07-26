package authorization

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/Clyra-AI/lumyn/schemas"
)

// ArtifactBinding is an exact persisted artifact identity and digest.
type ArtifactBinding struct {
	ArtifactID     string `json:"artifact_id"`
	ArtifactDigest string `json:"artifact_digest"`
}

// RepositoryBinding pins authorization to one repository, package root, and
// immutable base commit.
type RepositoryBinding struct {
	RepositoryID string `json:"repository_id"`
	PackageRoot  string `json:"package_root"`
	BaseCommit   string `json:"base_commit"`
}

// PolicyDigestBindings names every installation-policy projection referenced
// by an Event Authorization. The expected values come from the trusted policy
// compiler, not from the authorization being checked.
type PolicyDigestBindings struct {
	PathPolicy              string `json:"path_policy_digest"`
	CommandPolicy           string `json:"command_policy_digest"`
	ModelDataPolicy         string `json:"model_data_policy_digest"`
	BudgetPolicy            string `json:"budget_policy_digest"`
	GitHubPolicy            string `json:"github_policy_digest"`
	ProviderReportingPolicy string `json:"provider_reporting_policy_digest"`
}

// ExpectedAuthorizationBindings supplies the independently known inputs that
// a persisted Event Authorization must bind. Requiring this object prevents a
// self-consistent but substituted event/plan/manifest document from becoming
// runtime authority merely because it is schema-valid.
type ExpectedAuthorizationBindings struct {
	Event                           ArtifactBinding      `json:"event"`
	MigrationPack                   ArtifactBinding      `json:"migration_pack"`
	Repository                      RepositoryBinding    `json:"repository"`
	Plan                            ArtifactBinding      `json:"plan"`
	ExecutionManifest               ArtifactBinding      `json:"execution_manifest"`
	Policies                        PolicyDigestBindings `json:"policies"`
	VerificationConfigurationDigest string               `json:"verification_configuration_digest"`
	CredentialIssuancePolicyDigest  string               `json:"credential_issuance_policy_digest"`
	AgentRouteDigest                string               `json:"agent_route_digest,omitempty"`
	ExcludedPaths                   []string             `json:"excluded_paths"`
	ProviderReportingFields         []string             `json:"provider_reporting_fields"`
	Budgets                         Budgets              `json:"budgets"`
}

type persistedEventAuthorization struct {
	ObjectType                   string                    `json:"object_type"`
	SchemaVersion                string                    `json:"schema_version"`
	AuthorizationID              string                    `json:"authorization_id"`
	AuthorizationVersion         int                       `json:"authorization_version"`
	AuthorizationDigest          string                    `json:"authorization_digest"`
	APIConsumerOrganizationID    string                    `json:"api_consumer_organization_id"`
	AuthorizationOwnerRole       string                    `json:"authorization_owner_role"`
	InstallationBinding          ArtifactBinding           `json:"installation_binding"`
	EventBinding                 ArtifactBinding           `json:"event_binding"`
	MigrationPackBinding         ArtifactBinding           `json:"migration_pack_binding"`
	RepositoryBinding            RepositoryBinding         `json:"repository_binding"`
	PlanBinding                  ArtifactBinding           `json:"plan_binding"`
	ExecutionManifestBinding     ArtifactBinding           `json:"execution_manifest_binding"`
	AuthorizationMode            AuthorizationMode         `json:"authorization_mode"`
	EventProvenance              persistedEventProvenance  `json:"event_provenance"`
	SelectedAction               ActionMode                `json:"selected_action"`
	GenerationMode               string                    `json:"generation_mode"`
	AgentExecutionPolicy         persistedEventAgentPolicy `json:"agent_execution_policy"`
	SelectedCapabilities         []string                  `json:"selected_capabilities"`
	Scope                        persistedEventScope       `json:"scope"`
	Derivation                   persistedDerivation       `json:"derivation"`
	Approval                     persistedApproval         `json:"approval"`
	VerificationRequirement      persistedVerification     `json:"verification_requirement"`
	CredentialIssuance           persistedCredentialPolicy `json:"credential_issuance"`
	State                        string                    `json:"state"`
	IssuedAt                     time.Time                 `json:"issued_at"`
	ExpiresAt                    time.Time                 `json:"expires_at"`
	ActionLabelGrantsSideEffect  bool                      `json:"action_label_alone_grants_side_effect"`
	APIProviderMayApprove        bool                      `json:"api_provider_may_approve"`
	ProductionCredentialsAllowed bool                      `json:"production_credentials_allowed"`
	ProductionMutationAllowed    bool                      `json:"production_mutation_allowed"`
}

type persistedEventProvenance struct {
	TransportMode            string `json:"transport_mode"`
	Authenticated            bool   `json:"authenticated"`
	ChannelDeliveryQualified bool   `json:"channel_delivery_qualified"`
	AudienceMatched          bool   `json:"audience_matched"`
	Fresh                    bool   `json:"fresh"`
	SequenceAccepted         bool   `json:"sequence_accepted"`
	LifecycleActive          bool   `json:"lifecycle_active"`
}

type persistedEventAgentPolicy struct {
	State                 AgentExecutionState  `json:"state"`
	AdapterID             string               `json:"adapter_id"`
	AdapterVersion        string               `json:"adapter_version"`
	ExecutableDigest      string               `json:"executable_digest"`
	ConformanceDigest     string               `json:"conformance_digest"`
	AgentRouteDigest      string               `json:"agent_route_digest"`
	FundingMode           ExecutionFundingMode `json:"funding_mode"`
	CredentialOwnerRole   string               `json:"credential_owner_role"`
	UsageBillingOwnerRole string               `json:"usage_billing_owner_role"`
}

type persistedEventScope struct {
	ReadablePaths                 []string `json:"readable_paths"`
	WritablePaths                 []string `json:"writable_paths"`
	CommandIDs                    []string `json:"command_ids"`
	PathPolicyDigest              string   `json:"path_policy_digest"`
	CommandPolicyDigest           string   `json:"command_policy_digest"`
	ModelDataPolicyDigest         string   `json:"model_data_policy_digest"`
	BudgetPolicyDigest            string   `json:"budget_policy_digest"`
	GitHubPolicyDigest            string   `json:"github_policy_digest"`
	ProviderReportingPolicyDigest string   `json:"provider_reporting_policy_digest"`
}

type persistedDerivation struct {
	InstallationPolicyDigest           string `json:"installation_policy_digest"`
	ComparisonResult                   string `json:"comparison_result"`
	NonWidening                        bool   `json:"non_widening"`
	ProviderInputMayWiden              bool   `json:"provider_input_may_widen"`
	RouteSelectedCapabilityUnionFrozen bool   `json:"route_selected_capability_union_frozen"`
	AggregateCampaignAuthorityGranted  bool   `json:"aggregate_campaign_authority_granted"`
}

type persistedApproval struct {
	Mode                       AuthorizationMode `json:"mode"`
	ApprovedByRole             string            `json:"approved_by_role"`
	ConsumerMaintainerID       string            `json:"consumer_maintainer_id"`
	ApprovedPlanDigest         string            `json:"approved_plan_digest"`
	ApprovedAt                 *time.Time        `json:"approved_at"`
	EvaluatedBy                string            `json:"evaluated_by"`
	EventInPolicy              bool              `json:"event_in_policy"`
	PlanInPolicy               bool              `json:"plan_in_policy"`
	CommandsInPolicy           bool              `json:"commands_in_policy"`
	PathsInPolicy              bool              `json:"paths_in_policy"`
	ModelPolicyInPolicy        bool              `json:"model_policy_in_policy"`
	BudgetsInPolicy            bool              `json:"budgets_in_policy"`
	VerificationInPolicy       bool              `json:"verification_in_policy"`
	CredentialIssuanceInPolicy bool              `json:"credential_issuance_in_policy"`
	EvaluatedAt                *time.Time        `json:"evaluated_at"`
}

type persistedVerification struct {
	ConfigurationDigest                  string `json:"configuration_digest"`
	ExactCandidateHeadRequired           bool   `json:"exact_candidate_head_required"`
	IndependentFromGeneration            bool   `json:"independent_from_generation"`
	AgentOrModelCredentialsAllowed       bool   `json:"agent_or_model_credentials_allowed"`
	GenerationEvidenceWriteHandleAllowed bool   `json:"generation_evidence_write_handle_allowed"`
}

type persistedCredentialPolicy struct {
	PolicyDigest              string `json:"policy_digest"`
	ShortLivedOnly            bool   `json:"short_lived_only"`
	ReusableCredentialAllowed bool   `json:"reusable_credential_allowed"`
	IssuedAtRuntimeOnly       bool   `json:"issued_at_runtime_only"`
}

// ValidateSnapshotFromPersistedInstallation is the production authority
// entry point. Both persisted artifacts are schema-validated and strictly
// decoded, and every independently expected binding is checked before the
// resulting non-widening snapshot is returned.
func ValidateSnapshotFromPersistedInstallation(
	installationData []byte,
	authorizationData []byte,
	expected ExpectedAuthorizationBindings,
	now time.Time,
) (AuthorizationSnapshot, error) {
	installation, err := DecodeAndValidateInstallation(installationData, now)
	if err != nil {
		return AuthorizationSnapshot{}, err
	}
	return decodeAndValidateEventAuthorization(authorizationData, installation, expected, now)
}

func decodeAndValidateEventAuthorization(
	data []byte,
	installation Installation,
	expected ExpectedAuthorizationBindings,
	now time.Time,
) (AuthorizationSnapshot, error) {
	if err := validateExpectedBindings(expected); err != nil {
		return AuthorizationSnapshot{}, err
	}
	var schemaValue any
	if err := json.Unmarshal(data, &schemaValue); err != nil {
		return AuthorizationSnapshot{}, fmt.Errorf("decode event authorization JSON: %w", err)
	}
	if err := schemas.ValidateEventAuthorization(schemaValue); err != nil {
		return AuthorizationSnapshot{}, err
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var persisted persistedEventAuthorization
	if err := decoder.Decode(&persisted); err != nil {
		return AuthorizationSnapshot{}, fmt.Errorf("decode typed event authorization: %w", err)
	}
	if err := requireEventAuthorizationEOF(decoder); err != nil {
		return AuthorizationSnapshot{}, err
	}
	if err := persisted.validateBindings(installation, expected, now); err != nil {
		return AuthorizationSnapshot{}, err
	}

	snapshot, err := persisted.toAuthorizationSnapshot(installation, expected)
	if err != nil {
		return AuthorizationSnapshot{}, err
	}
	if err := validateSnapshot(installation, snapshot, now); err != nil {
		return AuthorizationSnapshot{}, fmt.Errorf("validate event authorization semantics: %w", err)
	}
	return snapshot, nil
}

func (value persistedEventAuthorization) validateBindings(
	installation Installation,
	expected ExpectedAuthorizationBindings,
	now time.Time,
) error {
	if value.ObjectType != "lumyn.event_authorization" || value.SchemaVersion != "1.0" {
		return fmt.Errorf("unsupported event authorization contract %q version %q", value.ObjectType, value.SchemaVersion)
	}
	if value.APIConsumerOrganizationID != installation.APIConsumerOrganizationID {
		return fmt.Errorf("event authorization binds a different consumer organization")
	}
	if err := equalArtifactBinding("installation", value.InstallationBinding, ArtifactBinding{installation.InstallationID, installation.InstallationDigest}); err != nil {
		return err
	}
	if value.Derivation.InstallationPolicyDigest != installation.InstallationDigest {
		return fmt.Errorf("derivation binds a different installation policy digest")
	}
	if value.AuthorizationMode != installation.AuthorizationMode || value.Approval.Mode != installation.AuthorizationMode {
		return fmt.Errorf("authorization mode differs from the Consumer Installation")
	}
	if err := equalArtifactBinding("event", value.EventBinding, expected.Event); err != nil {
		return err
	}
	if err := equalArtifactBinding("migration pack", value.MigrationPackBinding, expected.MigrationPack); err != nil {
		return err
	}
	if err := equalArtifactBinding("plan", value.PlanBinding, expected.Plan); err != nil {
		return err
	}
	if err := equalArtifactBinding("execution manifest", value.ExecutionManifestBinding, expected.ExecutionManifest); err != nil {
		return err
	}
	if value.RepositoryBinding != expected.Repository {
		return fmt.Errorf("event authorization binds a substituted repository")
	}
	if value.RepositoryBinding.RepositoryID != installation.RepositoryID || value.RepositoryBinding.PackageRoot != installation.PackageRoot {
		return fmt.Errorf("event authorization repository differs from the Consumer Installation")
	}
	if value.Scope.policyBindings() != expected.Policies {
		return fmt.Errorf("event authorization binds substituted installation policy digests")
	}
	if value.VerificationRequirement.ConfigurationDigest != expected.VerificationConfigurationDigest {
		return fmt.Errorf("event authorization binds a substituted verification configuration")
	}
	if value.CredentialIssuance.PolicyDigest != expected.CredentialIssuancePolicyDigest {
		return fmt.Errorf("event authorization binds a substituted credential-issuance policy")
	}
	if value.AuthorizationMode == PerEventApproval {
		if value.Approval.ConsumerMaintainerID != installation.ConsumerMaintainerID || value.Approval.ApprovedPlanDigest != value.PlanBinding.ArtifactDigest {
			return fmt.Errorf("per-event approval does not bind the installed maintainer and exact plan")
		}
	}
	if value.IssuedAt.After(now) || !value.ExpiresAt.After(now) || !value.ExpiresAt.After(value.IssuedAt) {
		return fmt.Errorf("event authorization is not active at the validation time")
	}
	if value.ExpiresAt.After(installation.ExpiresAt) {
		return fmt.Errorf("event authorization outlives the Consumer Installation")
	}
	return value.validateAgentBinding(installation, expected.AgentRouteDigest)
}

func (value persistedEventAuthorization) validateAgentBinding(installation Installation, expectedRouteDigest string) error {
	if value.AgentExecutionPolicy.State == AgentExecutionDisabled {
		if expectedRouteDigest != "" {
			return fmt.Errorf("disabled event authorization cannot bind an Agent Runner route")
		}
		return nil
	}
	if strings.TrimSpace(expectedRouteDigest) == "" || value.AgentExecutionPolicy.AgentRouteDigest != expectedRouteDigest {
		return fmt.Errorf("event authorization binds a substituted Agent Runner route digest")
	}
	runner := installation.AgentPolicy.Runner
	if installation.AgentPolicy.State != AgentExecutionConfigured || runner == nil {
		return fmt.Errorf("event authorization configures an uninstalled Agent Runner")
	}
	if value.AgentExecutionPolicy.AdapterID != runner.AdapterID ||
		value.AgentExecutionPolicy.AdapterVersion != runner.AdapterVersion ||
		value.AgentExecutionPolicy.ExecutableDigest != runner.ExecutableDigest ||
		value.AgentExecutionPolicy.ConformanceDigest != runner.ConformanceDigest ||
		value.AgentExecutionPolicy.FundingMode != runner.FundingMode ||
		value.AgentExecutionPolicy.CredentialOwnerRole != runner.CredentialOwner ||
		value.AgentExecutionPolicy.UsageBillingOwnerRole != runner.UsageBillingOwner {
		return fmt.Errorf("event authorization substitutes the installed Agent Runner route")
	}
	return nil
}

func (value persistedEventAuthorization) toAuthorizationSnapshot(installation Installation, expected ExpectedAuthorizationBindings) (AuthorizationSnapshot, error) {
	scope, err := scopeFromSelectedCapabilities(value.SelectedCapabilities)
	if err != nil {
		return AuthorizationSnapshot{}, err
	}
	scope.ReadPaths = append([]string(nil), value.Scope.ReadablePaths...)
	scope.WritePaths = append([]string(nil), value.Scope.WritablePaths...)
	scope.Commands = append([]string(nil), value.Scope.CommandIDs...)
	// Event Authorization carries policy digests for fields that are not
	// repeated in its JSON. Resolve those fields from the independently trusted
	// compiled-policy expectation, never from the authorization itself.
	scope.ExcludedPaths = append([]string(nil), expected.ExcludedPaths...)
	if scope.ProviderReporting {
		scope.ProviderReportingFields = append([]string(nil), expected.ProviderReportingFields...)
	}
	if len(scope.ReadPaths) > 0 && !containsCapability(value.SelectedCapabilities, "customer_repo_read") {
		return AuthorizationSnapshot{}, fmt.Errorf("readable paths require selected customer_repo_read capability")
	}
	if len(scope.WritePaths) > 0 && !containsCapability(value.SelectedCapabilities, "customer_repo_write") {
		return AuthorizationSnapshot{}, fmt.Errorf("writable paths require selected customer_repo_write capability")
	}
	if len(scope.Commands) > 0 && !containsCapability(value.SelectedCapabilities, "command_execution") {
		return AuthorizationSnapshot{}, fmt.Errorf("command IDs require selected command_execution capability")
	}

	route, err := routeFromGenerationMode(value.GenerationMode)
	if err != nil {
		return AuthorizationSnapshot{}, err
	}
	policy := value.Scope.policyBindings()
	return AuthorizationSnapshot{
		AuthorizationID:                     value.AuthorizationID,
		AuthorizationDigest:                 value.AuthorizationDigest,
		APIConsumerOrganizationID:           value.APIConsumerOrganizationID,
		InstallationID:                      value.InstallationBinding.ArtifactID,
		InstallationDigest:                  value.InstallationBinding.ArtifactDigest,
		EventID:                             value.EventBinding.ArtifactID,
		EventDigest:                         value.EventBinding.ArtifactDigest,
		MigrationPackID:                     value.MigrationPackBinding.ArtifactID,
		MigrationPackDigest:                 value.MigrationPackBinding.ArtifactDigest,
		RepositoryID:                        value.RepositoryBinding.RepositoryID,
		PackageRoot:                         value.RepositoryBinding.PackageRoot,
		BaseCommit:                          value.RepositoryBinding.BaseCommit,
		PlanID:                              value.PlanBinding.ArtifactID,
		PlanDigest:                          value.PlanBinding.ArtifactDigest,
		ExecutionManifestID:                 value.ExecutionManifestBinding.ArtifactID,
		ExecutionManifestDigest:             value.ExecutionManifestBinding.ArtifactDigest,
		PathPolicyDigest:                    policy.PathPolicy,
		CommandPolicyDigest:                 policy.CommandPolicy,
		ModelDataPolicyDigest:               policy.ModelDataPolicy,
		BudgetPolicyDigest:                  policy.BudgetPolicy,
		GitHubPolicyDigest:                  policy.GitHubPolicy,
		ProviderReportingPolicyDigest:       policy.ProviderReportingPolicy,
		VerificationConfigurationDigest:     value.VerificationRequirement.ConfigurationDigest,
		CredentialIssuancePolicyDigest:      value.CredentialIssuance.PolicyDigest,
		AgentRouteDigest:                    value.AgentExecutionPolicy.AgentRouteDigest,
		ActionMode:                          value.SelectedAction,
		Route:                               route,
		ExactPlanApproved:                   value.AuthorizationMode == PerEventApproval,
		InstalledPolicySatisfied:            value.AuthorizationMode == InstalledPreauthorization,
		ProviderChannelDeliveryProven:       value.EventProvenance.TransportMode == "pinned_provider_https" && value.EventProvenance.Authenticated && value.EventProvenance.ChannelDeliveryQualified,
		ShortLivedCredentialPolicySatisfied: value.CredentialIssuance.ShortLivedOnly && !value.CredentialIssuance.ReusableCredentialAllowed && value.CredentialIssuance.IssuedAtRuntimeOnly,
		Scope:                               scope,
		Budgets:                             expected.Budgets,
		IssuedAt:                            value.IssuedAt,
		ExpiresAt:                           value.ExpiresAt,
	}, nil
}

func scopeFromSelectedCapabilities(capabilities []string) (Scope, error) {
	scope := Scope{}
	for _, capability := range capabilities {
		switch capability {
		case "customer_repo_read", "customer_repo_write", "command_execution":
			// These capabilities are represented by their exact path/command lists.
		case "agent_runner_network":
			scope.AgentRunnerNetwork = true
		case "agent_runner_credential":
			scope.AgentRunnerCredential = true
		case "model_request_disclosure":
			scope.ModelRequestDisclosure = true
		case "model_network":
			scope.ModelNetwork = true
		case "model_credential":
			scope.ModelCredential = true
		case "package_registry_read":
			scope.Registry = true
		case "dependency_lifecycle_scripts":
			scope.LifecycleScripts = true
		case "sandbox_request_disclosure":
			scope.SandboxRequestData = true
		case "sandbox_network":
			scope.SandboxNetwork = true
		case "sandbox_credential":
			scope.SandboxCredentials = true
		case "github_branch_write":
			scope.RemoteBranch = true
		case "github_pr_write":
			scope.DraftPR = true
		case "provider_reporting":
			scope.ProviderReporting = true
		case "artifact_retention":
			scope.Retention = true
		case "artifact_deletion":
			scope.Deletion = true
		default:
			return Scope{}, fmt.Errorf("unknown selected capability %q", capability)
		}
	}
	return scope, nil
}

func routeFromGenerationMode(mode string) (Route, error) {
	switch mode {
	case "none":
		return Manual, nil
	case "deterministic":
		return Deterministic, nil
	case "agent_assisted":
		return AgentAssisted, nil
	default:
		return "", fmt.Errorf("unknown generation mode %q", mode)
	}
}

func (value persistedEventScope) policyBindings() PolicyDigestBindings {
	return PolicyDigestBindings{
		PathPolicy: value.PathPolicyDigest, CommandPolicy: value.CommandPolicyDigest,
		ModelDataPolicy: value.ModelDataPolicyDigest, BudgetPolicy: value.BudgetPolicyDigest,
		GitHubPolicy: value.GitHubPolicyDigest, ProviderReportingPolicy: value.ProviderReportingPolicyDigest,
	}
}

func validateExpectedBindings(value ExpectedAuthorizationBindings) error {
	for label, binding := range map[string]ArtifactBinding{
		"event": value.Event, "migration pack": value.MigrationPack,
		"plan": value.Plan, "execution manifest": value.ExecutionManifest,
	} {
		if strings.TrimSpace(binding.ArtifactID) == "" || strings.TrimSpace(binding.ArtifactDigest) == "" {
			return fmt.Errorf("expected %s binding is required", label)
		}
	}
	if strings.TrimSpace(value.Repository.RepositoryID) == "" || strings.TrimSpace(value.Repository.PackageRoot) == "" || strings.TrimSpace(value.Repository.BaseCommit) == "" {
		return fmt.Errorf("expected repository binding is required")
	}
	for label, digest := range map[string]string{
		"path": value.Policies.PathPolicy, "command": value.Policies.CommandPolicy,
		"model data": value.Policies.ModelDataPolicy, "budget": value.Policies.BudgetPolicy,
		"GitHub": value.Policies.GitHubPolicy, "provider reporting": value.Policies.ProviderReportingPolicy,
		"verification configuration": value.VerificationConfigurationDigest,
		"credential issuance":        value.CredentialIssuancePolicyDigest,
	} {
		if strings.TrimSpace(digest) == "" {
			return fmt.Errorf("expected %s policy digest is required", label)
		}
	}
	if err := validateUniqueStrings(value.ExcludedPaths, validateExcludedPath); err != nil {
		return fmt.Errorf("expected excluded path: %w", err)
	}
	if err := validateUniqueStrings(value.ProviderReportingFields, validatePolicyIdentifier); err != nil {
		return fmt.Errorf("expected provider reporting field: %w", err)
	}
	if err := validateBudgets(value.Budgets); err != nil {
		return fmt.Errorf("expected budget policy: %w", err)
	}
	return nil
}

func equalArtifactBinding(label string, actual, expected ArtifactBinding) error {
	if actual != expected {
		return fmt.Errorf("event authorization binds a substituted %s artifact", label)
	}
	return nil
}

func containsCapability(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}

func requireEventAuthorizationEOF(decoder *json.Decoder) error {
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return fmt.Errorf("event authorization contains trailing JSON")
		}
		return fmt.Errorf("decode event authorization trailer: %w", err)
	}
	return nil
}
