package evidence

import (
	"errors"
	"fmt"
	"path"
	"sort"
	"strings"
)

// HostBoundary supplies runtime-resolved host paths that must never be
// exposed to an untrusted command or credential-bearing sandbox process.
type HostBoundary struct {
	HomePath               string
	OSCredentialStorePaths []string
}

type ExecutionRepository struct {
	RepositoryID             string `json:"repository_id"`
	CheckoutPath             string `json:"checkout_path"`
	CheckoutDigest           string `json:"checkout_digest"`
	CandidateWorkspaceID     string `json:"candidate_workspace_id"`
	CandidateWorkspacePath   string `json:"candidate_workspace_path"`
	CandidateWorkspaceDigest string `json:"candidate_workspace_digest"`
	CandidateHead            string `json:"candidate_head"`
}

type PrivateStateBinding struct {
	PrivateStateID string `json:"private_state_id"`
	SourcePath     string `json:"source_path"`
	SourceDigest   string `json:"source_digest"`
}

type MountBinding struct {
	MountID         string `json:"mount_id"`
	SourceClass     string `json:"source_class"`
	SourceBindingID string `json:"source_binding_id"`
	SourcePath      string `json:"source_path"`
	SourceDigest    string `json:"source_digest"`
	TargetPath      string `json:"target_path"`
	Mode            string `json:"mode"`
}

type ExecutableRoot struct {
	ExecutableRootID string `json:"executable_root_id"`
	SourceMountID    string `json:"source_mount_id"`
	SourcePath       string `json:"source_path"`
	SourceDigest     string `json:"source_digest"`
	RootPath         string `json:"root_path"`
	ReadOnly         bool   `json:"read_only"`
}

type CommandBoundary struct {
	CommandID        string `json:"command_id"`
	Program          string `json:"program"`
	ProgramDigest    string `json:"program_digest"`
	ExecutableRootID string `json:"executable_root_id"`
	WorkingDirectory string `json:"working_directory"`
}

type SandboxCapabilities struct {
	RequestDisclosure bool `json:"sandbox_request_disclosure"`
	Network           bool `json:"sandbox_network"`
	Credential        bool `json:"sandbox_credential"`
}

type NetworkGrant struct {
	GrantID    string   `json:"grant_id"`
	Capability string   `json:"capability"`
	Endpoints  []string `json:"endpoints"`
	Operations []string `json:"operations"`
}

type CredentialGrant struct {
	GrantID             string   `json:"grant_id"`
	Kind                string   `json:"kind"`
	EnvironmentVariable string   `json:"environment_variable"`
	ScopeIDs            []string `json:"scope_ids"`
	InjectionBoundary   string   `json:"injection_boundary"`
	Reusable            bool     `json:"reusable"`
	Production          bool     `json:"production"`
}

type SandboxEntrypointProfile struct {
	ProfileID                     string          `json:"profile_id"`
	ProfileDigest                 string          `json:"profile_digest"`
	CandidateHead                 string          `json:"candidate_head"`
	CandidateMountID              string          `json:"candidate_mount_id"`
	CandidateMountTarget          string          `json:"candidate_mount_target"`
	EntrypointPath                string          `json:"entrypoint_path"`
	EntrypointDigest              string          `json:"entrypoint_digest"`
	Program                       string          `json:"program"`
	ProgramDigest                 string          `json:"program_digest"`
	ExecutableRootID              string          `json:"executable_root_id"`
	WorkingDirectory              string          `json:"working_directory"`
	NeutralHomeRoot               string          `json:"neutral_home_root"`
	NeutralTempRoot               string          `json:"neutral_temp_root"`
	SanitizedEnvironment          bool            `json:"sanitized_environment"`
	HostHomeMounted               bool            `json:"host_home_mounted"`
	OSCredentialStoreAvailable    bool            `json:"os_credential_store_available"`
	AmbientServiceSocketsAllowed  bool            `json:"ambient_service_sockets_allowed"`
	UnrelatedDescriptorsAllowed   bool            `json:"unrelated_inherited_descriptors_allowed"`
	SandboxCredentialGrantID      string          `json:"sandbox_credential_grant_id"`
	CredentialEnvironmentVariable string          `json:"credential_environment_variable"`
	CredentialInjectionBoundary   string          `json:"credential_injection_boundary"`
	SandboxNetworkGrantID         string          `json:"sandbox_network_grant_id"`
	NetworkEndpoints              []string        `json:"network_endpoints"`
	NetworkOperations             []string        `json:"network_operations"`
	NetworkDefault                string          `json:"network_default"`
	ChildProcessLimitsInherited   bool            `json:"child_process_limits_inherited"`
	ResourceQuotas                ResourceQuotas  `json:"resource_quotas"`
	Teardown                      SandboxTeardown `json:"teardown"`
}

type ResourceQuotas struct {
	CPUSeconds          int `json:"cpu_seconds"`
	MemoryBytes         int `json:"memory_bytes"`
	MaxPIDs             int `json:"max_pids"`
	MaxProcessTreeDepth int `json:"max_process_tree_depth"`
	DiskBytes           int `json:"disk_bytes"`
	MaxOpenFiles        int `json:"max_open_files"`
}

type SandboxTeardown struct {
	TeardownRequired     bool   `json:"teardown_required"`
	CleanupRequired      bool   `json:"cleanup_required"`
	OrphanCheckRequired  bool   `json:"orphan_check_required"`
	EvidencePolicyDigest string `json:"evidence_policy_digest"`
}

type ExecutionManifestBoundary struct {
	ManifestDigest           string                    `json:"manifest_digest"`
	Repository               ExecutionRepository       `json:"repository"`
	PrivateState             PrivateStateBinding       `json:"private_state"`
	Mounts                   []MountBinding            `json:"mounts"`
	ExecutableRoots          []ExecutableRoot          `json:"executable_roots"`
	Commands                 []CommandBoundary         `json:"commands"`
	Capabilities             SandboxCapabilities       `json:"capabilities"`
	NetworkGrants            []NetworkGrant            `json:"network_grants"`
	CredentialGrants         []CredentialGrant         `json:"credential_grants"`
	SandboxEntrypointProfile *SandboxEntrypointProfile `json:"sandbox_entrypoint_profile"`
}

// ValidateExecutionManifestBoundary enforces relationships JSON Schema cannot
// compare: exact mount sources, program-to-root membership, and the complete
// credential-bearing sandbox profile.
func ValidateExecutionManifestBoundary(value ExecutionManifestBoundary, host HostBoundary) error {
	if err := validateHostBoundary(host); err != nil {
		return err
	}
	mounts, err := validateMounts(value, host)
	if err != nil {
		return err
	}
	roots, err := validateExecutableRoots(value.ExecutableRoots, mounts)
	if err != nil {
		return err
	}
	if err := validateCommands(value.Commands, roots, mounts); err != nil {
		return err
	}
	return validateSandboxProfile(value, mounts, roots, host)
}

func validateHostBoundary(host HostBoundary) error {
	if err := requireNormalizedAbsolute("host home", host.HomePath); err != nil {
		return err
	}
	for _, credentialPath := range host.OSCredentialStorePaths {
		if err := requireNormalizedAbsolute("OS credential store", credentialPath); err != nil {
			return err
		}
	}
	return nil
}

func validateMounts(value ExecutionManifestBoundary, host HostBoundary) (map[string]MountBinding, error) {
	mounts := make(map[string]MountBinding, len(value.Mounts))
	for _, mount := range value.Mounts {
		if mount.MountID == "" || mount.SourceBindingID == "" || mount.SourceDigest == "" {
			return nil, errors.New("mount requires exact identity, source binding, and source digest")
		}
		if _, exists := mounts[mount.MountID]; exists {
			return nil, fmt.Errorf("duplicate mount id %q", mount.MountID)
		}
		if err := requireNormalizedAbsolute("mount source", mount.SourcePath); err != nil {
			return nil, err
		}
		if err := requireNormalizedAbsolute("mount target", mount.TargetPath); err != nil {
			return nil, err
		}
		if pathContains(mount.SourcePath, host.HomePath) {
			return nil, fmt.Errorf("mount %q exposes the host home", mount.MountID)
		}
		for _, credentialPath := range host.OSCredentialStorePaths {
			if pathsOverlap(mount.SourcePath, credentialPath) {
				return nil, fmt.Errorf("mount %q overlaps an OS credential store", mount.MountID)
			}
		}
		switch mount.SourceClass {
		case "consumer_repository":
			if mount.SourceBindingID != value.Repository.RepositoryID || mount.SourcePath != value.Repository.CheckoutPath || mount.SourceDigest != value.Repository.CheckoutDigest || mount.Mode != "read_only" {
				return nil, errors.New("consumer repository mount does not match its exact read-only repository binding")
			}
		case "candidate_workspace":
			if mount.SourceBindingID != value.Repository.CandidateWorkspaceID || mount.SourcePath != value.Repository.CandidateWorkspacePath || mount.SourceDigest != value.Repository.CandidateWorkspaceDigest {
				return nil, errors.New("candidate mount does not match its exact workspace binding")
			}
		case "consumer_private_state":
			if mount.SourceBindingID != value.PrivateState.PrivateStateID || mount.SourcePath != value.PrivateState.SourcePath || mount.SourceDigest != value.PrivateState.SourceDigest {
				return nil, errors.New("private-state mount does not match its exact private-state binding")
			}
		case "toolchain_readonly":
			if mount.Mode != "read_only" {
				return nil, errors.New("toolchain mount must be read-only")
			}
		default:
			return nil, fmt.Errorf("unsupported mount source class %q", mount.SourceClass)
		}
		mounts[mount.MountID] = mount
	}
	return mounts, nil
}

func validateExecutableRoots(values []ExecutableRoot, mounts map[string]MountBinding) (map[string]ExecutableRoot, error) {
	roots := make(map[string]ExecutableRoot, len(values))
	for _, root := range values {
		if root.ExecutableRootID == "" || root.SourceDigest == "" || !root.ReadOnly {
			return nil, errors.New("executable root requires identity, digest, and read-only posture")
		}
		if _, exists := roots[root.ExecutableRootID]; exists {
			return nil, fmt.Errorf("duplicate executable root id %q", root.ExecutableRootID)
		}
		if err := requireNormalizedAbsolute("executable source", root.SourcePath); err != nil {
			return nil, err
		}
		if err := requireNormalizedAbsolute("executable root", root.RootPath); err != nil {
			return nil, err
		}
		mount, ok := mounts[root.SourceMountID]
		if !ok || mount.SourceClass != "toolchain_readonly" || mount.Mode != "read_only" {
			return nil, fmt.Errorf("executable root %q must reference a read-only toolchain mount", root.ExecutableRootID)
		}
		if !pathContains(mount.SourcePath, root.SourcePath) || !pathContains(mount.TargetPath, root.RootPath) {
			return nil, fmt.Errorf("executable root %q escapes its source mount", root.ExecutableRootID)
		}
		roots[root.ExecutableRootID] = root
	}
	return roots, nil
}

func validateCommands(values []CommandBoundary, roots map[string]ExecutableRoot, mounts map[string]MountBinding) error {
	seen := make(map[string]struct{}, len(values))
	for _, command := range values {
		if command.CommandID == "" || command.ProgramDigest == "" {
			return errors.New("command requires exact identity and program digest")
		}
		if _, exists := seen[command.CommandID]; exists {
			return fmt.Errorf("duplicate command id %q", command.CommandID)
		}
		seen[command.CommandID] = struct{}{}
		if err := requireNormalizedAbsolute("command program", command.Program); err != nil {
			return err
		}
		if err := requireNormalizedAbsolute("command working directory", command.WorkingDirectory); err != nil {
			return err
		}
		root, ok := roots[command.ExecutableRootID]
		if !ok || !pathContains(root.RootPath, command.Program) {
			return fmt.Errorf("command %q is outside its declared executable root", command.CommandID)
		}
		if !withinAnyMount(command.WorkingDirectory, mounts) {
			return fmt.Errorf("command %q working directory is outside declared mounts", command.CommandID)
		}
	}
	return nil
}

func validateSandboxProfile(value ExecutionManifestBoundary, mounts map[string]MountBinding, roots map[string]ExecutableRoot, host HostBoundary) error {
	if err := validateGrantIDs(value.NetworkGrants, value.CredentialGrants); err != nil {
		return err
	}
	selected := value.Capabilities.RequestDisclosure || value.Capabilities.Network || value.Capabilities.Credential
	if selected && !(value.Capabilities.RequestDisclosure && value.Capabilities.Network && value.Capabilities.Credential) {
		return errors.New("sandbox disclosure, network, and credential capabilities must be selected together")
	}
	sandboxNetworks := filterSandboxNetworks(value.NetworkGrants)
	sandboxCredentials := filterSandboxCredentials(value.CredentialGrants)
	if !selected {
		if value.SandboxEntrypointProfile != nil || len(sandboxNetworks) != 0 || len(sandboxCredentials) != 0 {
			return errors.New("unselected sandbox route cannot carry profile, network, or credential grants")
		}
		return nil
	}
	profile := value.SandboxEntrypointProfile
	if profile == nil {
		return errors.New("selected sandbox route requires a dedicated entrypoint profile")
	}
	if len(sandboxNetworks) != 1 || len(sandboxCredentials) != 1 {
		return errors.New("sandbox entrypoint requires exactly one sandbox network and credential grant")
	}
	if profile.ProfileID == "" || profile.ProfileDigest == "" || profile.EntrypointDigest == "" || profile.ProgramDigest == "" {
		return errors.New("sandbox entrypoint requires exact profile, entrypoint, and program digests")
	}
	if profile.CandidateHead != value.Repository.CandidateHead {
		return errors.New("sandbox profile candidate head does not match the execution manifest")
	}
	mount, ok := mounts[profile.CandidateMountID]
	if !ok || mount.SourceClass != "candidate_workspace" || mount.Mode != "read_only" || mount.TargetPath != profile.CandidateMountTarget {
		return errors.New("sandbox profile must reference the exact read-only candidate mount")
	}
	if !pathContains(mount.TargetPath, profile.EntrypointPath) || !pathContains(mount.TargetPath, profile.WorkingDirectory) {
		return errors.New("sandbox entrypoint and working directory must remain inside the candidate mount")
	}
	root, ok := roots[profile.ExecutableRootID]
	if !ok || !pathContains(root.RootPath, profile.Program) {
		return errors.New("sandbox program is outside its declared executable root")
	}
	credential := sandboxCredentials[0]
	if credential.GrantID != profile.SandboxCredentialGrantID || credential.EnvironmentVariable != profile.CredentialEnvironmentVariable || credential.InjectionBoundary != "sandbox_verification" || profile.CredentialInjectionBoundary != "sandbox_verification" || credential.Reusable || credential.Production || len(credential.ScopeIDs) != 1 || credential.ScopeIDs[0] != profile.ProfileID {
		return errors.New("sandbox profile does not bind its sole task-scoped credential grant")
	}
	network := sandboxNetworks[0]
	if network.GrantID != profile.SandboxNetworkGrantID || !sameStrings(network.Endpoints, profile.NetworkEndpoints) || !sameStrings(network.Operations, profile.NetworkOperations) {
		return errors.New("sandbox profile does not bind its exact endpoint and operation grant")
	}
	if err := requireNormalizedAbsolute("sandbox neutral home", profile.NeutralHomeRoot); err != nil {
		return err
	}
	if err := requireNormalizedAbsolute("sandbox neutral temp", profile.NeutralTempRoot); err != nil {
		return err
	}
	for _, neutralRoot := range []string{profile.NeutralHomeRoot, profile.NeutralTempRoot} {
		if pathsOverlap(neutralRoot, host.HomePath) {
			return errors.New("sandbox neutral roots must not overlap the runtime host home")
		}
		for _, credentialPath := range host.OSCredentialStorePaths {
			if pathsOverlap(neutralRoot, credentialPath) {
				return errors.New("sandbox neutral roots must not overlap an OS credential store")
			}
		}
		for _, mount := range mounts {
			if pathsOverlap(neutralRoot, mount.TargetPath) {
				return errors.New("sandbox neutral roots must remain outside declared mount targets")
			}
		}
	}
	if pathsOverlap(profile.NeutralHomeRoot, profile.NeutralTempRoot) {
		return errors.New("sandbox neutral home and temporary roots must not overlap")
	}
	if profile.HostHomeMounted || profile.OSCredentialStoreAvailable || profile.AmbientServiceSocketsAllowed || profile.UnrelatedDescriptorsAllowed || !profile.SanitizedEnvironment || !profile.ChildProcessLimitsInherited || profile.NetworkDefault != "offline" {
		return errors.New("sandbox profile weakens the required isolation boundary")
	}
	if !profile.ResourceQuotas.Valid() {
		return errors.New("sandbox profile requires positive CPU, memory, PID, process-tree, disk, and open-file quotas")
	}
	if !profile.Teardown.TeardownRequired || !profile.Teardown.CleanupRequired || !profile.Teardown.OrphanCheckRequired || profile.Teardown.EvidencePolicyDigest == "" {
		return errors.New("sandbox profile requires teardown, cleanup, orphan, and evidence controls")
	}
	return nil
}

type SandboxVerificationEvidence struct {
	EvidenceKind                   string `json:"evidence_kind"`
	CandidateHead                  string `json:"candidate_head"`
	SandboxEntrypointProfileDigest string `json:"sandbox_entrypoint_profile_digest,omitempty"`
}

type MigrationVerificationBoundary struct {
	ConsumerExecutionManifestDigest string                        `json:"consumer_execution_manifest_digest"`
	CandidateHead                   string                        `json:"candidate_head"`
	VerificationLabel               string                        `json:"verification_label"`
	SandboxEntrypointProfileDigest  string                        `json:"sandbox_entrypoint_profile_digest,omitempty"`
	ObservedEvidence                []SandboxVerificationEvidence `json:"observed_evidence_refs"`
}

// ValidateMigrationVerificationBoundary prevents a sandbox label from being
// carried by unrelated exact-head or sandbox evidence.
func ValidateMigrationVerificationBoundary(value MigrationVerificationBoundary, manifest ExecutionManifestBoundary) error {
	if value.ConsumerExecutionManifestDigest == "" || value.ConsumerExecutionManifestDigest != manifest.ManifestDigest {
		return errors.New("verification does not bind the exact consumer execution manifest")
	}
	if value.CandidateHead == "" || value.CandidateHead != manifest.Repository.CandidateHead {
		return errors.New("verification candidate head does not match the execution manifest")
	}
	for _, item := range value.ObservedEvidence {
		if item.CandidateHead == "" || item.CandidateHead != value.CandidateHead {
			return fmt.Errorf("%s evidence does not bind the exact candidate head", item.EvidenceKind)
		}
	}
	if value.VerificationLabel != "workflow_verified_sandbox" {
		if value.SandboxEntrypointProfileDigest != "" {
			return errors.New("non-sandbox verification cannot bind a sandbox profile")
		}
		for _, item := range value.ObservedEvidence {
			if item.EvidenceKind == "sandbox_observation" {
				return errors.New("non-sandbox verification cannot carry sandbox observation evidence")
			}
			if item.SandboxEntrypointProfileDigest != "" {
				return errors.New("non-sandbox evidence cannot bind a sandbox profile")
			}
		}
		if value.VerificationLabel == "workflow_verified_replay" || value.VerificationLabel == "workflow_verified_mock" {
			if !hasEvidenceKind(value.ObservedEvidence, "exact_head_workflow_outcome") {
				return errors.New("workflow replay or mock verification requires exact-head outcome evidence")
			}
		}
		return nil
	}
	profile := manifest.SandboxEntrypointProfile
	if profile == nil || value.SandboxEntrypointProfileDigest == "" || value.SandboxEntrypointProfileDigest != profile.ProfileDigest {
		return errors.New("sandbox verification does not bind the exact sandbox entrypoint profile")
	}
	required := map[string]bool{"exact_head_workflow_outcome": false, "sandbox_observation": false}
	for _, item := range value.ObservedEvidence {
		if _, ok := required[item.EvidenceKind]; !ok {
			continue
		}
		if item.CandidateHead != value.CandidateHead || item.SandboxEntrypointProfileDigest != profile.ProfileDigest {
			return fmt.Errorf("%s evidence is stale or profile-mismatched", item.EvidenceKind)
		}
		required[item.EvidenceKind] = true
	}
	for kind, present := range required {
		if !present {
			return fmt.Errorf("sandbox verification requires %s evidence", kind)
		}
	}
	return nil
}

func filterSandboxNetworks(values []NetworkGrant) []NetworkGrant {
	result := make([]NetworkGrant, 0, 1)
	for _, value := range values {
		if value.Capability == "sandbox_network" {
			result = append(result, value)
		}
	}
	return result
}

func filterSandboxCredentials(values []CredentialGrant) []CredentialGrant {
	result := make([]CredentialGrant, 0, 1)
	for _, value := range values {
		if value.Kind == "sandbox_credential" {
			result = append(result, value)
		}
	}
	return result
}

func validateGrantIDs(networks []NetworkGrant, credentials []CredentialGrant) error {
	seen := make(map[string]struct{}, len(networks)+len(credentials))
	for _, network := range networks {
		if network.GrantID == "" {
			return errors.New("network grant identity is required")
		}
		if _, exists := seen[network.GrantID]; exists {
			return fmt.Errorf("duplicate grant id %q", network.GrantID)
		}
		seen[network.GrantID] = struct{}{}
	}
	for _, credential := range credentials {
		if credential.GrantID == "" {
			return errors.New("credential grant identity is required")
		}
		if _, exists := seen[credential.GrantID]; exists {
			return fmt.Errorf("duplicate grant id %q", credential.GrantID)
		}
		seen[credential.GrantID] = struct{}{}
	}
	return nil
}

func (value ResourceQuotas) Valid() bool {
	return value.CPUSeconds > 0 && value.MemoryBytes > 0 && value.MaxPIDs > 0 &&
		value.MaxProcessTreeDepth > 0 && value.DiskBytes > 0 && value.MaxOpenFiles > 0
}

func hasEvidenceKind(values []SandboxVerificationEvidence, expected string) bool {
	for _, value := range values {
		if value.EvidenceKind == expected {
			return true
		}
	}
	return false
}

func requireNormalizedAbsolute(label, value string) error {
	if value == "" || !path.IsAbs(value) || path.Clean(value) != value || value == "/" {
		return fmt.Errorf("%s must be a normalized absolute non-root path", label)
	}
	return nil
}

func pathContains(root, candidate string) bool {
	return candidate == root || strings.HasPrefix(candidate, root+"/")
}

func pathsOverlap(left, right string) bool {
	return pathContains(left, right) || pathContains(right, left)
}

func withinAnyMount(candidate string, mounts map[string]MountBinding) bool {
	for _, mount := range mounts {
		if pathContains(mount.TargetPath, candidate) {
			return true
		}
	}
	return false
}

func sameStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	leftCopy := append([]string(nil), left...)
	rightCopy := append([]string(nil), right...)
	sort.Strings(leftCopy)
	sort.Strings(rightCopy)
	for index := range leftCopy {
		if leftCopy[index] != rightCopy[index] {
			return false
		}
	}
	return true
}
