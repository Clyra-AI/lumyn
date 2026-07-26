// Package authorization defines the fail-closed, non-widening policy contract
// shared by Consumer Installations and immutable event authorizations.
package authorization

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/Clyra-AI/lumyn/internal/installation"
)

type ActionMode string

const (
	NotifyOnly   ActionMode = "notify_only"
	ScanOnly     ActionMode = "scan_only"
	PreparePatch ActionMode = "prepare_patch"
	OpenDraftPR  ActionMode = "open_draft_pr"
)

type AuthorizationMode string

const (
	PerEventApproval          AuthorizationMode = "per_event_approval"
	InstalledPreauthorization AuthorizationMode = "installed_preauthorization"
)

type AgentExecutionState string

const (
	AgentExecutionDisabled   AgentExecutionState = "disabled"
	AgentExecutionConfigured AgentExecutionState = "configured"
)

type ExecutionFundingMode string

const (
	ConsumerManaged               ExecutionFundingMode = "consumer_managed"
	ProviderSponsoredLumynManaged ExecutionFundingMode = "provider_sponsored_lumyn_managed"
)

type Route string

const (
	Deterministic Route = "deterministic"
	AgentAssisted Route = "agent_assisted"
	Manual        Route = "manual"
)

type Scope struct {
	ReadPaths     []string `json:"read_paths"`
	WritePaths    []string `json:"write_paths"`
	ExcludedPaths []string `json:"excluded_paths"`
	// Commands contains canonical command IDs. Exact program/argument policy is
	// separately bound by CommandPolicyDigest on persisted event authority.
	Commands                []string `json:"commands"`
	AgentRunnerNetwork      bool     `json:"agent_runner_network"`
	AgentRunnerCredential   bool     `json:"agent_runner_credential"`
	ModelRequestDisclosure  bool     `json:"model_request_disclosure"`
	ModelNetwork            bool     `json:"model_network"`
	ModelCredential         bool     `json:"model_credential"`
	Registry                bool     `json:"registry"`
	LifecycleScripts        bool     `json:"lifecycle_scripts"`
	SandboxRequestData      bool     `json:"sandbox_request_data"`
	SandboxNetwork          bool     `json:"sandbox_network"`
	SandboxCredentials      bool     `json:"sandbox_credentials"`
	RemoteBranch            bool     `json:"remote_branch"`
	DraftPR                 bool     `json:"draft_pr"`
	ProviderReporting       bool     `json:"provider_reporting"`
	ProviderReportingFields []string `json:"provider_reporting_fields"`
	Retention               bool     `json:"retention"`
	Deletion                bool     `json:"deletion"`
}

type Budgets struct {
	MaxChangedFiles    int `json:"max_changed_files"`
	MaxDiffLines       int `json:"max_diff_lines"`
	MaxDiffBytes       int `json:"max_diff_bytes"`
	MaxTurns           int `json:"max_turns"`
	MaxAttempts        int `json:"max_attempts"`
	MaxTokens          int `json:"max_tokens"`
	MaxCostCents       int `json:"max_cost_cents"`
	MaxDurationSeconds int `json:"max_duration_seconds"`
}

type ManagedCredentialPolicy struct {
	BrokerIssuer                 string `json:"broker_issuer"`
	MaximumTTLSeconds            int    `json:"maximum_ttl_seconds"`
	OneTimeRedemption            bool   `json:"one_time_redemption"`
	RefreshAllowed               bool   `json:"refresh_allowed"`
	CrossAttemptReuseAllowed     bool   `json:"cross_attempt_reuse_allowed"`
	HardTokenQuota               bool   `json:"hard_token_quota"`
	HardCostQuota                bool   `json:"hard_cost_quota"`
	Revocation                   bool   `json:"revocation"`
	UsageReconciliation          bool   `json:"usage_reconciliation"`
	VendorNativeOrEnforcingProxy bool   `json:"vendor_native_or_enforcing_proxy"`
}

type RunnerPolicy struct {
	AdapterID                 string                   `json:"adapter_id"`
	AdapterVersion            string                   `json:"adapter_version"`
	ExecutablePath            string                   `json:"executable_path"`
	ExecutableSource          string                   `json:"executable_source"`
	ExecutableDigest          string                   `json:"executable_digest"`
	ConformanceDigest         string                   `json:"conformance_digest"`
	AuthMode                  string                   `json:"auth_mode"`
	EntitlementClass          string                   `json:"entitlement_class"`
	AgentRunnerVendor         string                   `json:"agent_runner_vendor"`
	ActualModelProvider       string                   `json:"actual_model_provider"`
	ActualModelRoute          string                   `json:"actual_model_route"`
	RouteTopology             string                   `json:"route_topology"`
	FundingMode               ExecutionFundingMode     `json:"funding_mode"`
	CredentialOwner           string                   `json:"credential_owner"`
	UsageBillingOwner         string                   `json:"usage_billing_owner"`
	NativeConfigurationDigest string                   `json:"native_configuration_digest"`
	CleanSession              bool                     `json:"clean_session"`
	NeutralHome               bool                     `json:"neutral_home"`
	SilentFallback            bool                     `json:"silent_fallback"`
	ReusableCredentialStored  bool                     `json:"reusable_credential_stored"`
	ManagedCredential         *ManagedCredentialPolicy `json:"managed_credential,omitempty"`
}

type AgentExecutionPolicy struct {
	State  AgentExecutionState `json:"state"`
	Runner *RunnerPolicy       `json:"runner,omitempty"`
}

type Installation struct {
	InstallationID                string               `json:"installation_id"`
	InstallationDigest            string               `json:"installation_digest"`
	APIConsumerOrganizationID     string               `json:"api_consumer_organization_id"`
	ConsumerMaintainerID          string               `json:"consumer_maintainer_id"`
	Repository                    string               `json:"repository"`
	RepositoryID                  string               `json:"repository_id"`
	PackageRoot                   string               `json:"package_root"`
	ProviderChannelOrigin         string               `json:"provider_channel_origin"`
	ProviderManifestURL           string               `json:"provider_manifest_url"`
	ProviderChannelKeyID          string               `json:"provider_channel_key_id"`
	ActionMode                    ActionMode           `json:"action_mode"`
	AuthorizationMode             AuthorizationMode    `json:"authorization_mode"`
	ExpiresAt                     time.Time            `json:"expires_at"`
	Revoked                       bool                 `json:"revoked"`
	AgentPolicy                   AgentExecutionPolicy `json:"agent_execution_policy"`
	StoredGitHubToken             bool                 `json:"stored_github_token"`
	ProviderMaySelectRunner       bool                 `json:"provider_may_select_runner"`
	ProviderMayAccessConsumerData bool                 `json:"provider_may_access_consumer_data"`
	ProviderReportingConsent      bool                 `json:"provider_reporting_consent"`
	Scope                         Scope                `json:"scope"`
	Budgets                       Budgets              `json:"budgets"`
}

type AuthorizationSnapshot struct {
	AuthorizationID                     string     `json:"authorization_id"`
	AuthorizationDigest                 string     `json:"authorization_digest"`
	APIConsumerOrganizationID           string     `json:"api_consumer_organization_id"`
	InstallationID                      string     `json:"installation_id"`
	InstallationDigest                  string     `json:"installation_digest"`
	EventID                             string     `json:"event_id"`
	EventDigest                         string     `json:"event_digest"`
	MigrationPackID                     string     `json:"migration_pack_id"`
	MigrationPackDigest                 string     `json:"migration_pack_digest"`
	RepositoryID                        string     `json:"repository_id"`
	PackageRoot                         string     `json:"package_root"`
	BaseCommit                          string     `json:"base_commit"`
	PlanID                              string     `json:"plan_id"`
	PlanDigest                          string     `json:"plan_digest"`
	ExecutionManifestID                 string     `json:"execution_manifest_id"`
	ExecutionManifestDigest             string     `json:"execution_manifest_digest"`
	PathPolicyDigest                    string     `json:"path_policy_digest"`
	CommandPolicyDigest                 string     `json:"command_policy_digest"`
	ModelDataPolicyDigest               string     `json:"model_data_policy_digest"`
	BudgetPolicyDigest                  string     `json:"budget_policy_digest"`
	GitHubPolicyDigest                  string     `json:"github_policy_digest"`
	ProviderReportingPolicyDigest       string     `json:"provider_reporting_policy_digest"`
	VerificationConfigurationDigest     string     `json:"verification_configuration_digest"`
	CredentialIssuancePolicyDigest      string     `json:"credential_issuance_policy_digest"`
	AgentRouteDigest                    string     `json:"agent_route_digest,omitempty"`
	ActionMode                          ActionMode `json:"action_mode"`
	Route                               Route      `json:"route"`
	ExactPlanApproved                   bool       `json:"exact_plan_approved"`
	InstalledPolicySatisfied            bool       `json:"installed_policy_satisfied"`
	ProviderChannelDeliveryProven       bool       `json:"provider_channel_delivery_proven"`
	ShortLivedCredentialPolicySatisfied bool       `json:"short_lived_credential_policy_satisfied"`
	Scope                               Scope      `json:"scope"`
	Budgets                             Budgets    `json:"budgets"`
	IssuedAt                            time.Time  `json:"issued_at"`
	ExpiresAt                           time.Time  `json:"expires_at"`
}

func ValidateInstallation(value Installation, now time.Time) error {
	if strings.TrimSpace(value.InstallationID) == "" || strings.TrimSpace(value.Repository) == "" {
		return errors.New("installation and repository identity are required")
	}
	if err := validateRelativePath(value.PackageRoot); err != nil {
		return fmt.Errorf("package root: %w", err)
	}
	if err := validatePinnedHTTPSOrigin(value.ProviderChannelOrigin); err != nil {
		return fmt.Errorf("provider channel: %w", err)
	}
	if err := installation.ValidateURLAtPinnedOrigin(value.ProviderChannelOrigin, value.ProviderManifestURL); err != nil {
		return fmt.Errorf("provider channel manifest: %w", err)
	}
	if strings.TrimSpace(value.ProviderChannelKeyID) == "" {
		return errors.New("provider channel authentication key is required")
	}
	if _, ok := actionRanks[value.ActionMode]; !ok {
		return fmt.Errorf("unknown action mode %q", value.ActionMode)
	}
	if value.AuthorizationMode != PerEventApproval && value.AuthorizationMode != InstalledPreauthorization {
		return fmt.Errorf("unknown authorization mode %q", value.AuthorizationMode)
	}
	if value.Revoked {
		return errors.New("installation is revoked")
	}
	if value.ExpiresAt.IsZero() || !value.ExpiresAt.After(now) {
		return errors.New("installation is expired")
	}
	if value.StoredGitHubToken {
		return errors.New("stored GitHub token is prohibited")
	}
	if value.ProviderMaySelectRunner || value.ProviderMayAccessConsumerData {
		return errors.New("API Provider cannot select the runner or access consumer data")
	}
	if value.Scope.ProviderReporting && !value.ProviderReportingConsent {
		return errors.New("provider reporting requires explicit consent")
	}
	if err := validateScope(value.Scope); err != nil {
		return err
	}
	if err := validateBudgets(value.Budgets); err != nil {
		return err
	}
	if err := validateActionScope(value.ActionMode, value.Scope); err != nil {
		return err
	}
	return validateAgentPolicy(value.AgentPolicy, value.Scope)
}

func validateSnapshot(installation Installation, snapshot AuthorizationSnapshot, now time.Time) error {
	if err := ValidateInstallation(installation, now); err != nil {
		return fmt.Errorf("installation: %w", err)
	}
	if strings.TrimSpace(snapshot.AuthorizationID) == "" || strings.TrimSpace(snapshot.EventID) == "" ||
		strings.TrimSpace(snapshot.EventDigest) == "" || strings.TrimSpace(snapshot.PlanDigest) == "" {
		return errors.New("authorization, event, and exact plan bindings are required")
	}
	if snapshot.InstallationID != installation.InstallationID {
		return errors.New("authorization binds a different installation")
	}
	if snapshot.InstallationDigest != "" && snapshot.InstallationDigest != installation.InstallationDigest {
		return errors.New("authorization binds a different installation digest")
	}
	if _, ok := actionRanks[snapshot.ActionMode]; !ok {
		return fmt.Errorf("unknown authorization action mode %q", snapshot.ActionMode)
	}
	if actionRanks[snapshot.ActionMode] > actionRanks[installation.ActionMode] {
		return errors.New("event authorization widens the installed action ceiling")
	}
	if snapshot.Route != Deterministic && snapshot.Route != AgentAssisted && snapshot.Route != Manual {
		return fmt.Errorf("unknown route %q", snapshot.Route)
	}
	if snapshot.Route == AgentAssisted && installation.AgentPolicy.State != AgentExecutionConfigured {
		return errors.New("agent-assisted route requires configured agent execution")
	}
	if snapshot.Route != AgentAssisted && (snapshot.Scope.AgentRunnerNetwork || snapshot.Scope.AgentRunnerCredential ||
		snapshot.Scope.ModelRequestDisclosure || snapshot.Scope.ModelNetwork || snapshot.Scope.ModelCredential) {
		return errors.New("non-agent route cannot receive Agent Runner or model authority")
	}
	if err := validateScope(snapshot.Scope); err != nil {
		return fmt.Errorf("authorization scope: %w", err)
	}
	if !scopeIsSubset(snapshot.Scope, installation.Scope) {
		return errors.New("event authorization widens installed scope")
	}
	if err := validateActionScope(snapshot.ActionMode, snapshot.Scope); err != nil {
		return fmt.Errorf("authorization action ceiling: %w", err)
	}
	if snapshot.Route == AgentAssisted {
		if err := validateTopology(installation.AgentPolicy.Runner.RouteTopology, snapshot.Scope); err != nil {
			return fmt.Errorf("authorization agent route: %w", err)
		}
	}
	if err := validateBudgets(snapshot.Budgets); err != nil {
		return fmt.Errorf("authorization budgets: %w", err)
	}
	if !budgetsWithin(snapshot.Budgets, installation.Budgets) {
		return errors.New("event authorization widens installed budgets")
	}
	switch installation.AuthorizationMode {
	case PerEventApproval:
		if !snapshot.ExactPlanApproved {
			return errors.New("per-event approval requires the exact plan approval")
		}
	case InstalledPreauthorization:
		if !snapshot.InstalledPolicySatisfied {
			return errors.New("installed preauthorization requires policy satisfaction")
		}
		if !snapshot.ProviderChannelDeliveryProven {
			return errors.New("installed preauthorization requires authenticated provider-channel delivery")
		}
	}
	if (snapshot.Scope.RemoteBranch || snapshot.Scope.DraftPR) && !snapshot.ShortLivedCredentialPolicySatisfied {
		return errors.New("remote delivery requires satisfied short-lived credential policy")
	}
	return nil
}

var actionRanks = map[ActionMode]int{
	NotifyOnly: 0, ScanOnly: 1, PreparePatch: 2, OpenDraftPR: 3,
}

func validateAgentPolicy(value AgentExecutionPolicy, scope Scope) error {
	switch value.State {
	case AgentExecutionDisabled:
		if value.Runner != nil || scope.AgentRunnerNetwork || scope.AgentRunnerCredential ||
			scope.ModelRequestDisclosure || scope.ModelNetwork || scope.ModelCredential {
			return errors.New("disabled agent execution grants no runner or model authority")
		}
		return nil
	case AgentExecutionConfigured:
		if value.Runner == nil {
			return errors.New("configured agent execution requires an exact runner policy")
		}
	default:
		return fmt.Errorf("unknown agent execution state %q", value.State)
	}
	runner := value.Runner
	required := map[string]string{
		"adapter id": runner.AdapterID, "adapter version": runner.AdapterVersion,
		"executable path": runner.ExecutablePath, "executable source": runner.ExecutableSource,
		"executable digest": runner.ExecutableDigest, "conformance digest": runner.ConformanceDigest,
		"auth mode": runner.AuthMode, "entitlement class": runner.EntitlementClass,
		"Agent Runner Vendor": runner.AgentRunnerVendor, "actual Model Provider": runner.ActualModelProvider,
		"actual model route": runner.ActualModelRoute, "route topology": runner.RouteTopology,
		"credential owner": runner.CredentialOwner, "usage-billing owner": runner.UsageBillingOwner,
		"native configuration digest": runner.NativeConfigurationDigest,
	}
	for name, field := range required {
		if strings.TrimSpace(field) == "" {
			return fmt.Errorf("%s is required", name)
		}
	}
	if runner.AdapterID != "codex" && runner.AdapterID != "claude_code" {
		return fmt.Errorf("agent adapter %q is not a qualified launch target", runner.AdapterID)
	}
	if !filepath.IsAbs(runner.ExecutablePath) {
		return errors.New("runner executable path must be canonical and absolute")
	}
	if !runner.CleanSession || !runner.NeutralHome || runner.SilentFallback || runner.ReusableCredentialStored {
		return errors.New("runner requires a clean neutral session with no fallback or reusable credential")
	}
	if err := validateTopology(runner.RouteTopology, scope); err != nil {
		return err
	}
	switch runner.FundingMode {
	case ConsumerManaged:
		if runner.ManagedCredential != nil {
			return errors.New("consumer-managed funding cannot attach a Lumyn-managed credential")
		}
		if runner.CredentialOwner != "api_consumer_organization" || runner.UsageBillingOwner != "api_consumer_organization" {
			return errors.New("consumer-managed funding requires consumer credential and usage-billing ownership")
		}
	case ProviderSponsoredLumynManaged:
		if runner.CredentialOwner != "lumyn_operator" || runner.UsageBillingOwner != "lumyn_operator" {
			return errors.New("Lumyn-managed funding requires Lumyn credential and usage-billing ownership")
		}
		// The durable Consumer Installation selects funding and ownership but
		// stores no reusable credential. If a broker ceiling is attached by an
		// attended caller, validate it here; the attempt-scoped executable grant
		// remains a separate managed-credential artifact.
		if runner.ManagedCredential != nil {
			if err := validateManagedCredential(runner.ManagedCredential); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unknown execution funding mode %q", runner.FundingMode)
	}
	return nil
}

func validateManagedCredential(value *ManagedCredentialPolicy) error {
	if value == nil || strings.TrimSpace(value.BrokerIssuer) == "" {
		return errors.New("managed funding requires an approved broker issuer")
	}
	if value.MaximumTTLSeconds < 1 || value.MaximumTTLSeconds > 3600 {
		return errors.New("managed credential TTL must be at most one hour")
	}
	if !value.OneTimeRedemption || value.RefreshAllowed || value.CrossAttemptReuseAllowed ||
		!value.HardTokenQuota || !value.HardCostQuota || !value.Revocation ||
		!value.UsageReconciliation || !value.VendorNativeOrEnforcingProxy {
		return errors.New("managed credential is not one-time, quota-bound, revocable, and reconcilable")
	}
	return nil
}

func validateTopology(topology string, scope Scope) error {
	switch topology {
	case "local_runtime":
		if scope.AgentRunnerNetwork || scope.AgentRunnerCredential || scope.ModelNetwork || scope.ModelCredential {
			return errors.New("local runtime cannot receive remote runner or model scope")
		}
	case "runner_mediated":
		if !scope.AgentRunnerNetwork || !scope.AgentRunnerCredential || !scope.ModelRequestDisclosure {
			return errors.New("runner-mediated route requires runner network, credential, and model disclosure")
		}
	case "direct_model":
		if !scope.ModelRequestDisclosure || !scope.ModelNetwork || !scope.ModelCredential {
			return errors.New("direct-model route requires model disclosure, network, and credential")
		}
	case "hybrid":
		if !scope.AgentRunnerNetwork || !scope.AgentRunnerCredential || !scope.ModelRequestDisclosure || !scope.ModelNetwork || !scope.ModelCredential {
			return errors.New("hybrid route requires both runner and model scope sets")
		}
	default:
		return fmt.Errorf("unknown agent route topology %q", topology)
	}
	return nil
}

func validateScope(value Scope) error {
	for label, paths := range map[string][]string{"read path": value.ReadPaths, "write path": value.WritePaths} {
		if err := validateUniqueStrings(paths, validateRelativePath); err != nil {
			return fmt.Errorf("%s: %w", label, err)
		}
	}
	if err := validateUniqueStrings(value.ExcludedPaths, validateExcludedPath); err != nil {
		return fmt.Errorf("excluded path: %w", err)
	}
	if err := validateUniqueStrings(value.Commands, validateCommand); err != nil {
		return fmt.Errorf("command: %w", err)
	}
	if err := validateUniqueStrings(value.ProviderReportingFields, validatePolicyIdentifier); err != nil {
		return fmt.Errorf("provider reporting field: %w", err)
	}
	if !value.ProviderReporting && len(value.ProviderReportingFields) > 0 {
		return errors.New("provider reporting fields require provider_reporting capability")
	}
	return nil
}

func validatePolicyIdentifier(value string) error {
	if strings.TrimSpace(value) == "" || strings.ContainsAny(value, "\x00\r\n") {
		return fmt.Errorf("unsafe identifier %q", value)
	}
	return nil
}

func validateExcludedPath(value string) error {
	if value == "" || strings.ContainsAny(value, "\x00\\*") || filepath.IsAbs(value) {
		return fmt.Errorf("unsafe path %q", value)
	}
	cleaned := filepath.Clean(value)
	if cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return fmt.Errorf("unsafe path %q", value)
	}
	return nil
}

func validateRelativePath(value string) error {
	if value == "" || strings.ContainsAny(value, "\x00\\*") || filepath.IsAbs(value) {
		return fmt.Errorf("unsafe path %q", value)
	}
	cleaned := filepath.Clean(value)
	if cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) || cleaned == ".git" || strings.HasPrefix(cleaned, ".git"+string(filepath.Separator)) {
		return fmt.Errorf("unsafe path %q", value)
	}
	return nil
}

func validateCommand(value string) error {
	if strings.TrimSpace(value) == "" || strings.ContainsAny(value, "\x00\r\n") {
		return fmt.Errorf("unsafe command %q", value)
	}
	return nil
}

func validateUniqueStrings(values []string, validate func(string) error) error {
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if err := validate(value); err != nil {
			return err
		}
		if _, exists := seen[value]; exists {
			return fmt.Errorf("duplicate value %q", value)
		}
		seen[value] = struct{}{}
	}
	return nil
}

func validatePinnedHTTPSOrigin(value string) error {
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil || parsed.Fragment != "" ||
		(parsed.Path != "" && parsed.Path != "/") || parsed.RawQuery != "" {
		return fmt.Errorf("origin must be an exact HTTPS URL")
	}
	return nil
}

func validateBudgets(value Budgets) error {
	if value.MaxChangedFiles < 0 || value.MaxDiffLines < 0 || value.MaxDiffBytes < 0 || value.MaxTurns < 0 || value.MaxAttempts < 0 ||
		value.MaxTokens < 0 || value.MaxCostCents < 0 || value.MaxDurationSeconds < 0 {
		return errors.New("budgets cannot be negative")
	}
	return nil
}

func validateActionScope(mode ActionMode, scope Scope) error {
	if mode == NotifyOnly && len(scope.ReadPaths) > 0 {
		return errors.New("notify-only action cannot read the repository")
	}
	if actionRanks[mode] < actionRanks[PreparePatch] && scopeHasMutation(scope) {
		return errors.New("action ceiling does not permit mutation or execution scope")
	}
	if actionRanks[mode] < actionRanks[OpenDraftPR] && (scope.RemoteBranch || scope.DraftPR) {
		return errors.New("action ceiling does not permit remote delivery")
	}
	return nil
}

func scopeHasMutation(value Scope) bool {
	return len(value.WritePaths) > 0 || len(value.Commands) > 0 || value.AgentRunnerNetwork ||
		value.AgentRunnerCredential || value.ModelRequestDisclosure || value.ModelNetwork ||
		value.ModelCredential || value.Registry || value.LifecycleScripts ||
		value.SandboxRequestData || value.SandboxNetwork || value.SandboxCredentials ||
		value.RemoteBranch || value.DraftPR
}

func scopeIsSubset(child, parent Scope) bool {
	return stringSubset(child.ReadPaths, parent.ReadPaths) &&
		stringSubset(child.WritePaths, parent.WritePaths) &&
		stringSubset(parent.ExcludedPaths, child.ExcludedPaths) &&
		stringSubset(child.Commands, parent.Commands) &&
		(!child.AgentRunnerNetwork || parent.AgentRunnerNetwork) &&
		(!child.AgentRunnerCredential || parent.AgentRunnerCredential) &&
		(!child.ModelRequestDisclosure || parent.ModelRequestDisclosure) &&
		(!child.ModelNetwork || parent.ModelNetwork) &&
		(!child.ModelCredential || parent.ModelCredential) &&
		(!child.Registry || parent.Registry) &&
		(!child.LifecycleScripts || parent.LifecycleScripts) &&
		(!child.SandboxRequestData || parent.SandboxRequestData) &&
		(!child.SandboxNetwork || parent.SandboxNetwork) &&
		(!child.SandboxCredentials || parent.SandboxCredentials) &&
		(!child.RemoteBranch || parent.RemoteBranch) &&
		(!child.DraftPR || parent.DraftPR) &&
		(!child.ProviderReporting || parent.ProviderReporting) &&
		stringSubset(child.ProviderReportingFields, parent.ProviderReportingFields) &&
		(!child.Retention || parent.Retention) &&
		(!child.Deletion || parent.Deletion)
}

func stringSubset(child, parent []string) bool {
	allowed := make(map[string]struct{}, len(parent))
	for _, value := range parent {
		allowed[value] = struct{}{}
	}
	for _, value := range child {
		if _, ok := allowed[value]; !ok {
			return false
		}
	}
	return true
}

func budgetsWithin(child, parent Budgets) bool {
	return child.MaxChangedFiles <= parent.MaxChangedFiles && child.MaxDiffLines <= parent.MaxDiffLines &&
		child.MaxDiffBytes <= parent.MaxDiffBytes && child.MaxTurns <= parent.MaxTurns &&
		child.MaxAttempts <= parent.MaxAttempts && child.MaxTokens <= parent.MaxTokens &&
		child.MaxCostCents <= parent.MaxCostCents && child.MaxDurationSeconds <= parent.MaxDurationSeconds
}
