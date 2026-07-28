// Package pack normalizes pinned provider evidence into a declarative,
// non-executable Provider Change Contract.
package pack

import (
	"time"

	"github.com/Clyra-AI/lumyn/internal/source"
)

const SchemaVersion = "1.0"

type APIOrSDK struct {
	Kind    string `json:"kind"`
	Name    string `json:"name"`
	Package string `json:"package"`
}

type Audience struct {
	AudienceID             string `json:"audience_id"`
	PackageVersionSelector string `json:"package_version_selector"`
}

type SourceEvidence struct {
	EvidenceID    string `json:"evidence_id"`
	Kind          string `json:"kind"`
	VersionRole   string `json:"version_role"`
	Version       string `json:"version"`
	Location      string `json:"location"`
	Digest        string `json:"digest"`
	License       string `json:"license"`
	SourceLocator string `json:"source_locator"`
}

type Mapping struct {
	Mode             string   `json:"mode"`
	DeclarativeSteps []string `json:"declarative_steps"`
}

type SemanticChange struct {
	ChangeID          string   `json:"change_id"`
	Kind              string   `json:"kind"`
	SourceSymbol      string   `json:"source_symbol"`
	TargetSymbol      *string  `json:"target_symbol"`
	Intent            string   `json:"intent"`
	Applicability     string   `json:"applicability"`
	Mapping           Mapping  `json:"mapping"`
	SourceEvidenceIDs []string `json:"source_evidence_ids"`
	TargetEvidenceIDs []string `json:"target_evidence_ids"`
}

type Ambiguity struct {
	AmbiguityID        string   `json:"ambiguity_id"`
	Description        string   `json:"description"`
	RequiredResolution string   `json:"required_resolution"`
	AffectedChangeIDs  []string `json:"affected_change_ids"`
}

type UnsupportedCase struct {
	CaseID      string `json:"case_id"`
	Description string `json:"description"`
	Reason      string `json:"reason"`
}

type VerificationGuidance struct {
	RequiredBehaviors       []string `json:"required_behaviors"`
	ProviderCommandsAllowed bool     `json:"provider_commands_allowed"`
}

type Provenance struct {
	Class         string    `json:"class"`
	SourceURLs    []string  `json:"source_urls"`
	SourceDigests []string  `json:"source_digests"`
	NormalizedAt  time.Time `json:"normalized_at"`
}

type ProviderConfirmation struct {
	Status                  string     `json:"status"`
	ProviderOperatorID      *string    `json:"provider_operator_id"`
	ConfirmedContractDigest *string    `json:"confirmed_contract_digest"`
	ConfirmedAt             *time.Time `json:"confirmed_at"`
}

type Lifecycle struct {
	State            string  `json:"state"`
	SupersededBy     *string `json:"superseded_by"`
	WithdrawalReason *string `json:"withdrawal_reason"`
}

type Contract struct {
	ObjectType                   string               `json:"object_type"`
	SchemaVersion                string               `json:"schema_version"`
	MigrationPackID              string               `json:"migration_pack_id"`
	ContractVersion              int                  `json:"contract_version"`
	ContractDigest               string               `json:"contract_digest"`
	APIProviderID                string               `json:"api_provider_id"`
	APIOrSDK                     APIOrSDK             `json:"api_or_sdk"`
	SourceVersion                string               `json:"source_version"`
	TargetVersion                string               `json:"target_version"`
	Audience                     Audience             `json:"audience"`
	SourceEvidence               []SourceEvidence     `json:"source_evidence"`
	SemanticChanges              []SemanticChange     `json:"semantic_changes"`
	Ambiguities                  []Ambiguity          `json:"ambiguities"`
	UnsupportedCases             []UnsupportedCase    `json:"unsupported_cases"`
	VerificationGuidance         VerificationGuidance `json:"verification_guidance"`
	Provenance                   Provenance           `json:"provenance"`
	ProviderConfirmation         ProviderConfirmation `json:"provider_confirmation"`
	Lifecycle                    Lifecycle            `json:"lifecycle"`
	NonExecutable                bool                 `json:"non_executable"`
	GrantsConsumerAuthority      bool                 `json:"grants_consumer_authority"`
	ProductionCredentialsAllowed bool                 `json:"production_credentials_allowed"`
	ProductionMutationAllowed    bool                 `json:"production_mutation_allowed"`
}

// Declaration is a typed, source-bound semantic assertion. It is data, not a
// script or prompt, and every change must cite both migration sides.
type Declaration struct {
	SemanticChange
}

type NormalizeRequest struct {
	MigrationPackID      string
	ContractVersion      int
	APIProviderID        string
	APIOrSDK             APIOrSDK
	SourceVersion        string
	TargetVersion        string
	Audience             Audience
	Evidence             []SourceEvidenceInput
	Declarations         []Declaration
	Ambiguities          []Ambiguity
	UnsupportedCases     []UnsupportedCase
	VerificationGuidance []string
	NormalizedAt         time.Time
	Lifecycle            Lifecycle
}

type SourceEvidenceInput struct {
	Artifact      source.PinnedArtifact
	SourceLocator string
}

type Finding struct {
	Code        string   `json:"code"`
	Message     string   `json:"message"`
	EvidenceIDs []string `json:"evidence_ids"`
	ChangeIDs   []string `json:"change_ids"`
}

type NormalizationReport struct {
	Status     string    `json:"status"`
	Actionable bool      `json:"actionable"`
	Findings   []Finding `json:"findings"`
}

type Result struct {
	Contract Contract
	Report   NormalizationReport
}
