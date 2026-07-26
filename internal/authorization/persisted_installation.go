package authorization

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"strings"
	"time"

	"github.com/Clyra-AI/lumyn/schemas"
)

// persistedInstallation mirrors the security-relevant shape of
// lumyn.consumer_installation. Fields that do not participate in authority
// derivation remain RawMessage values after the canonical JSON Schema has
// validated their complete shape.
type persistedInstallation struct {
	ObjectType                           string                   `json:"object_type"`
	SchemaVersion                        string                   `json:"schema_version"`
	InstallationID                       string                   `json:"installation_id"`
	InstallationVersion                  int                      `json:"installation_version"`
	InstallationDigest                   string                   `json:"installation_digest"`
	APIProviderID                        string                   `json:"api_provider_id"`
	APIConsumerOrganizationID            string                   `json:"api_consumer_organization_id"`
	ConsumerMaintainerID                 string                   `json:"consumer_maintainer_id"`
	AuthorityOwnerRole                   string                   `json:"authority_owner_role"`
	ProviderChannel                      persistedProviderChannel `json:"provider_channel"`
	Repository                           persistedRepository      `json:"repository"`
	Selectors                            json.RawMessage          `json:"selectors"`
	ActionCeiling                        ActionMode               `json:"action_ceiling"`
	AuthorizationMode                    AuthorizationMode        `json:"authorization_mode"`
	PathPolicy                           persistedPathPolicy      `json:"path_policy"`
	Commands                             []persistedCommand       `json:"commands"`
	CapabilityCeiling                    persistedCapabilities    `json:"capability_ceiling"`
	AgentExecutionPolicy                 persistedAgentPolicy     `json:"agent_execution_policy"`
	ModelDataPolicy                      json.RawMessage          `json:"model_data_policy"`
	Budgets                              persistedBudgets         `json:"budgets"`
	GitHubTokenIssuancePolicy            persistedGitHubPolicy    `json:"github_token_issuance_policy"`
	ProviderReporting                    persistedReporting       `json:"provider_reporting"`
	RetentionAndDeletion                 json.RawMessage          `json:"retention_and_deletion"`
	ExpiresAt                            time.Time                `json:"expires_at"`
	Revocation                           persistedRevocation      `json:"revocation"`
	ActionLabelGrantsSideEffect          bool                     `json:"action_label_grants_side_effect"`
	StoresReusableAgentOrModelCredential bool                     `json:"stores_reusable_agent_or_model_credential"`
	StoresGitHubToken                    bool                     `json:"stores_github_token"`
	ProductionCredentialsAllowed         bool                     `json:"production_credentials_allowed"`
	ProductionMutationAllowed            bool                     `json:"production_mutation_allowed"`
	APIProviderAccess                    persistedProviderAccess  `json:"api_provider_access"`
}

type persistedProviderChannel struct {
	TransportID          string `json:"transport_id"`
	ManifestURL          string `json:"manifest_url"`
	PinnedOrigin         string `json:"pinned_origin"`
	CampaignKeyID        string `json:"campaign_key_id"`
	CampaignPublicKey    string `json:"campaign_public_key"`
	AuthenticationMethod string `json:"authentication_method"`
}

type persistedRepository struct {
	Host          string `json:"host"`
	Owner         string `json:"owner"`
	Name          string `json:"name"`
	PackageRoot   string `json:"package_root"`
	DefaultBranch string `json:"default_branch"`
}

type persistedPathPolicy struct {
	ReadablePaths []string `json:"readable_paths"`
	WritablePaths []string `json:"writable_paths"`
	ExcludedPaths []string `json:"excluded_paths"`
}

type persistedCommand struct {
	CommandID               string   `json:"command_id"`
	Program                 string   `json:"program"`
	Arguments               []string `json:"arguments"`
	WorkingDirectory        string   `json:"working_directory"`
	NetworkMode             string   `json:"network_mode"`
	LifecycleScriptsAllowed bool     `json:"lifecycle_scripts_allowed"`
}

type persistedCapabilities struct {
	CustomerRepoRead          bool `json:"customer_repo_read"`
	CustomerRepoWrite         bool `json:"customer_repo_write"`
	CommandExecution          bool `json:"command_execution"`
	AgentRunnerNetwork        bool `json:"agent_runner_network"`
	AgentRunnerCredential     bool `json:"agent_runner_credential"`
	ModelRequestDisclosure    bool `json:"model_request_disclosure"`
	ModelNetwork              bool `json:"model_network"`
	ModelCredential           bool `json:"model_credential"`
	PackageRegistryRead       bool `json:"package_registry_read"`
	DependencyLifecycleScript bool `json:"dependency_lifecycle_scripts"`
	SandboxRequestDisclosure  bool `json:"sandbox_request_disclosure"`
	SandboxNetwork            bool `json:"sandbox_network"`
	SandboxCredential         bool `json:"sandbox_credential"`
	GitHubBranchWrite         bool `json:"github_branch_write"`
	GitHubPRWrite             bool `json:"github_pr_write"`
	ProviderReporting         bool `json:"provider_reporting"`
	ArtifactRetention         bool `json:"artifact_retention"`
	ArtifactDeletion          bool `json:"artifact_deletion"`
}

type persistedAgentPolicy struct {
	State                        AgentExecutionState       `json:"state"`
	SelectionAuthorityRole       string                    `json:"selection_authority_role"`
	AdapterID                    string                    `json:"adapter_id"`
	AdapterVersion               string                    `json:"adapter_version"`
	AgentRunnerVendor            string                    `json:"agent_runner_vendor"`
	Executable                   persistedExecutable       `json:"executable"`
	ConformanceDigest            string                    `json:"conformance_digest"`
	LiveCanaryQualified          bool                      `json:"live_canary_qualified"`
	AuthMode                     string                    `json:"auth_mode"`
	EntitlementClass             string                    `json:"entitlement_class"`
	RouteTopology                string                    `json:"route_topology"`
	ExecutionFundingMode         ExecutionFundingMode      `json:"execution_funding_mode"`
	CredentialOwnerRole          string                    `json:"credential_owner_role"`
	UsageBillingOwnerRole        string                    `json:"usage_billing_owner_role"`
	ActualModelRoute             persistedActualModelRoute `json:"actual_model_route"`
	NativeConfiguration          persistedNativeConfig     `json:"native_configuration"`
	CleanSessionRequired         bool                      `json:"clean_session_required"`
	NeutralHomeAndConfigRequired bool                      `json:"neutral_home_and_config_required"`
	ExecutableExtensionsAllowed  bool                      `json:"executable_extensions_allowed"`
	SilentFallbackAllowed        bool                      `json:"silent_fallback_allowed"`
}

type persistedExecutable struct {
	CanonicalPath                  string `json:"canonical_path"`
	ApprovedSource                 string `json:"approved_source"`
	Digest                         string `json:"digest"`
	RepositoryPathShadowingAllowed bool   `json:"repository_path_shadowing_allowed"`
}

type persistedActualModelRoute struct {
	ModelProvider string `json:"model_provider"`
	Endpoint      string `json:"endpoint"`
	Model         string `json:"model"`
	ModelVersion  string `json:"model_version"`
	Observable    bool   `json:"observable"`
}

type persistedNativeConfig struct {
	Mode                      string `json:"mode"`
	Digest                    string `json:"digest"`
	TreatedAsUntrustedContext bool   `json:"treated_as_untrusted_context"`
}

type persistedBudgets struct {
	MaxFiles           int         `json:"max_files"`
	MaxChangedLines    int         `json:"max_changed_lines"`
	MaxDiffBytes       int         `json:"max_diff_bytes"`
	MaxTurns           int         `json:"max_turns"`
	MaxTokens          int         `json:"max_tokens"`
	MaxCostUSD         json.Number `json:"max_cost_usd"`
	MaxDurationSeconds int         `json:"max_duration_seconds"`
	MaxAttempts        int         `json:"max_attempts"`
}

type persistedGitHubPolicy struct {
	Enabled              bool    `json:"enabled"`
	BrokerID             *string `json:"broker_id"`
	ShortLived           bool    `json:"short_lived"`
	MaximumTTLSeconds    int     `json:"maximum_ttl_seconds"`
	NonDefaultBranchOnly bool    `json:"non_default_branch_only"`
	DraftPROnly          bool    `json:"draft_pr_only"`
	AutoMergeAllowed     bool    `json:"auto_merge_allowed"`
	StoredTokenAllowed   bool    `json:"stored_token_allowed"`
}

type persistedReporting struct {
	Enabled                    bool     `json:"enabled"`
	AllowedFields              []string `json:"allowed_fields"`
	RawConsumerEvidenceAllowed bool     `json:"raw_consumer_evidence_allowed"`
	ConsentRequired            bool     `json:"consent_required"`
}

type persistedRevocation struct {
	State         string     `json:"state"`
	RevokedAt     *time.Time `json:"revoked_at"`
	RevokedByRole *string    `json:"revoked_by_role"`
}

type persistedProviderAccess struct {
	Source         bool `json:"source"`
	Context        bool `json:"context"`
	AgentSession   bool `json:"agent_session"`
	Credentials    bool `json:"credentials"`
	AgentSelection bool `json:"agent_selection"`
}

// DecodeAndValidateInstallation decodes the canonical persisted Consumer
// Installation contract, rejects structural drift, converts it without
// widening authority, and applies semantic validation before returning an
// Installation usable by ValidateSnapshot.
func DecodeAndValidateInstallation(data []byte, now time.Time) (Installation, error) {
	var schemaValue any
	if err := json.Unmarshal(data, &schemaValue); err != nil {
		return Installation{}, fmt.Errorf("decode consumer installation JSON: %w", err)
	}
	if err := schemas.ValidateConsumerInstallation(schemaValue); err != nil {
		return Installation{}, err
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	decoder.UseNumber()
	var persisted persistedInstallation
	if err := decoder.Decode(&persisted); err != nil {
		return Installation{}, fmt.Errorf("decode typed consumer installation: %w", err)
	}
	if err := requireInstallationEOF(decoder); err != nil {
		return Installation{}, err
	}

	installation, err := persisted.toAuthorizationInstallation()
	if err != nil {
		return Installation{}, err
	}
	if err := ValidateInstallation(installation, now); err != nil {
		return Installation{}, fmt.Errorf("validate consumer installation semantics: %w", err)
	}
	return installation, nil
}

func (value persistedInstallation) toAuthorizationInstallation() (Installation, error) {
	if value.ObjectType != "lumyn.consumer_installation" || value.SchemaVersion != "1.0" {
		return Installation{}, fmt.Errorf("unsupported consumer installation contract %q version %q", value.ObjectType, value.SchemaVersion)
	}
	if value.AuthorityOwnerRole != "api_consumer_organization" {
		return Installation{}, fmt.Errorf("consumer installation authority must be owned by api_consumer_organization")
	}
	if value.ActionLabelGrantsSideEffect || value.ProductionCredentialsAllowed || value.ProductionMutationAllowed {
		return Installation{}, fmt.Errorf("consumer installation cannot grant implicit or production authority")
	}
	if value.agentExtensionsEnabled() {
		return Installation{}, fmt.Errorf("consumer installation cannot enable executable agent extensions")
	}
	if len(value.PathPolicy.ReadablePaths) > 0 && !value.CapabilityCeiling.CustomerRepoRead {
		return Installation{}, fmt.Errorf("readable paths require customer_repo_read capability")
	}
	if len(value.PathPolicy.WritablePaths) > 0 && !value.CapabilityCeiling.CustomerRepoWrite {
		return Installation{}, fmt.Errorf("writable paths require customer_repo_write capability")
	}
	if len(value.Commands) > 0 && !value.CapabilityCeiling.CommandExecution {
		return Installation{}, fmt.Errorf("commands require command_execution capability")
	}
	if (value.CapabilityCeiling.GitHubBranchWrite || value.CapabilityCeiling.GitHubPRWrite) && !value.GitHubTokenIssuancePolicy.Enabled {
		return Installation{}, fmt.Errorf("remote delivery capability requires GitHub token issuance policy")
	}

	commands := make([]string, 0, len(value.Commands))
	for _, command := range value.Commands {
		commands = append(commands, command.CommandID)
	}
	costCents, err := dollarsToCents(value.Budgets.MaxCostUSD)
	if err != nil {
		return Installation{}, fmt.Errorf("max_cost_usd: %w", err)
	}
	providerReportingFields := []string(nil)
	if value.CapabilityCeiling.ProviderReporting {
		providerReportingFields = append(providerReportingFields, value.ProviderReporting.AllowedFields...)
	}

	installation := Installation{
		InstallationID:                value.InstallationID,
		InstallationDigest:            value.InstallationDigest,
		APIConsumerOrganizationID:     value.APIConsumerOrganizationID,
		ConsumerMaintainerID:          value.ConsumerMaintainerID,
		Repository:                    strings.Join([]string{value.Repository.Host, value.Repository.Owner, value.Repository.Name}, "/"),
		RepositoryID:                  strings.Join([]string{"github", value.Repository.Owner, value.Repository.Name}, "."),
		PackageRoot:                   value.Repository.PackageRoot,
		ProviderChannelOrigin:         value.ProviderChannel.PinnedOrigin,
		ProviderManifestURL:           value.ProviderChannel.ManifestURL,
		ProviderChannelKeyID:          value.ProviderChannel.CampaignKeyID,
		ActionMode:                    value.ActionCeiling,
		AuthorizationMode:             value.AuthorizationMode,
		ExpiresAt:                     value.ExpiresAt,
		Revoked:                       value.Revocation.State == "revoked",
		AgentPolicy:                   value.AgentExecutionPolicy.toAuthorizationPolicy(value.StoresReusableAgentOrModelCredential),
		StoredGitHubToken:             value.StoresGitHubToken || value.GitHubTokenIssuancePolicy.StoredTokenAllowed,
		ProviderMaySelectRunner:       value.APIProviderAccess.AgentSelection,
		ProviderMayAccessConsumerData: value.APIProviderAccess.Source || value.APIProviderAccess.Context || value.APIProviderAccess.AgentSession || value.APIProviderAccess.Credentials,
		ProviderReportingConsent:      value.ProviderReporting.Enabled && value.ProviderReporting.ConsentRequired,
		Scope: Scope{
			ReadPaths:               value.PathPolicy.ReadablePaths,
			WritePaths:              value.PathPolicy.WritablePaths,
			ExcludedPaths:           value.PathPolicy.ExcludedPaths,
			Commands:                commands,
			AgentRunnerNetwork:      value.CapabilityCeiling.AgentRunnerNetwork,
			AgentRunnerCredential:   value.CapabilityCeiling.AgentRunnerCredential,
			ModelRequestDisclosure:  value.CapabilityCeiling.ModelRequestDisclosure,
			ModelNetwork:            value.CapabilityCeiling.ModelNetwork,
			ModelCredential:         value.CapabilityCeiling.ModelCredential,
			Registry:                value.CapabilityCeiling.PackageRegistryRead,
			LifecycleScripts:        value.CapabilityCeiling.DependencyLifecycleScript,
			SandboxRequestData:      value.CapabilityCeiling.SandboxRequestDisclosure,
			SandboxNetwork:          value.CapabilityCeiling.SandboxNetwork,
			SandboxCredentials:      value.CapabilityCeiling.SandboxCredential,
			RemoteBranch:            value.CapabilityCeiling.GitHubBranchWrite,
			DraftPR:                 value.CapabilityCeiling.GitHubPRWrite,
			ProviderReporting:       value.CapabilityCeiling.ProviderReporting,
			ProviderReportingFields: providerReportingFields,
			Retention:               value.CapabilityCeiling.ArtifactRetention,
			Deletion:                value.CapabilityCeiling.ArtifactDeletion,
		},
		Budgets: Budgets{
			MaxChangedFiles:    value.Budgets.MaxFiles,
			MaxDiffLines:       value.Budgets.MaxChangedLines,
			MaxDiffBytes:       value.Budgets.MaxDiffBytes,
			MaxTurns:           value.Budgets.MaxTurns,
			MaxAttempts:        value.Budgets.MaxAttempts,
			MaxTokens:          value.Budgets.MaxTokens,
			MaxCostCents:       costCents,
			MaxDurationSeconds: value.Budgets.MaxDurationSeconds,
		},
	}
	return installation, nil
}

func (value persistedInstallation) agentExtensionsEnabled() bool {
	return value.AgentExecutionPolicy.ExecutableExtensionsAllowed
}

func (value persistedAgentPolicy) toAuthorizationPolicy(storesReusableCredential bool) AgentExecutionPolicy {
	policy := AgentExecutionPolicy{State: value.State}
	if value.State != AgentExecutionConfigured {
		return policy
	}
	nativeConfigurationDigest := "disabled"
	if value.NativeConfiguration.Mode == "selected_static" {
		nativeConfigurationDigest = value.NativeConfiguration.Digest
	}
	policy.Runner = &RunnerPolicy{
		AdapterID:                 value.AdapterID,
		AdapterVersion:            value.AdapterVersion,
		ExecutablePath:            value.Executable.CanonicalPath,
		ExecutableSource:          value.Executable.ApprovedSource,
		ExecutableDigest:          value.Executable.Digest,
		ConformanceDigest:         value.ConformanceDigest,
		AuthMode:                  value.AuthMode,
		EntitlementClass:          value.EntitlementClass,
		AgentRunnerVendor:         value.AgentRunnerVendor,
		ActualModelProvider:       value.ActualModelRoute.ModelProvider,
		ActualModelRoute:          fmt.Sprintf("%s@%s via %s", value.ActualModelRoute.Model, value.ActualModelRoute.ModelVersion, value.ActualModelRoute.Endpoint),
		RouteTopology:             value.RouteTopology,
		FundingMode:               value.ExecutionFundingMode,
		CredentialOwner:           value.CredentialOwnerRole,
		UsageBillingOwner:         value.UsageBillingOwnerRole,
		NativeConfigurationDigest: nativeConfigurationDigest,
		CleanSession:              value.CleanSessionRequired,
		NeutralHome:               value.NeutralHomeAndConfigRequired,
		SilentFallback:            value.SilentFallbackAllowed,
		ReusableCredentialStored:  storesReusableCredential,
	}
	return policy
}

func dollarsToCents(value json.Number) (int, error) {
	amount, ok := new(big.Rat).SetString(value.String())
	if !ok || amount.Sign() < 0 {
		return 0, fmt.Errorf("must be a non-negative decimal amount")
	}
	amount.Mul(amount, big.NewRat(100, 1))
	if !amount.IsInt() || !amount.Num().IsInt64() {
		return 0, fmt.Errorf("must resolve to whole cents within the supported range")
	}
	cents := amount.Num().Int64()
	converted := int(cents)
	if int64(converted) != cents {
		return 0, fmt.Errorf("is outside the supported range")
	}
	return converted, nil
}

func requireInstallationEOF(decoder *json.Decoder) error {
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return fmt.Errorf("consumer installation contains trailing JSON")
		}
		return fmt.Errorf("decode consumer installation trailer: %w", err)
	}
	return nil
}
