package pack

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/Clyra-AI/lumyn/internal/source"
)

var (
	stableIDPattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._:-]{0,127}$`)
	digestPattern   = regexp.MustCompile(`^sha256:[a-f0-9]{64}$`)
)

var allowedChangeKinds = map[string]bool{
	"rename": true, "signature_change": true, "request_shape_change": true,
	"response_shape_change": true, "endpoint_change": true,
	"behavior_change": true, "removal": true,
}

var allowedMappingModes = map[string]bool{
	"exact": true, "repository_specific_reasoning": true, "manual_input_required": true,
}

// BlockedError preserves the reviewable contract and explains why no affected
// route may proceed.
type BlockedError struct {
	Findings []Finding
}

func (err *BlockedError) Error() string {
	return "Provider Change Contract normalization blocked by ambiguous or conflicting evidence"
}

func Normalize(request NormalizeRequest) (Result, error) {
	result := Result{Report: NormalizationReport{Status: "blocked", Findings: []Finding{}}}
	if err := validateRequestEnvelope(request); err != nil {
		return result, err
	}

	evidenceByID := make(map[string]SourceEvidenceInput, len(request.Evidence))
	sourceCount, targetCount := 0, 0
	for _, item := range request.Evidence {
		normalized, err := validateEvidence(item, request)
		if err != nil {
			return result, err
		}
		item = normalized
		if _, exists := evidenceByID[item.Artifact.ID]; exists {
			return result, fmt.Errorf("duplicate source evidence id %q", item.Artifact.ID)
		}
		evidenceByID[item.Artifact.ID] = item
		switch item.Artifact.VersionRole {
		case source.VersionSource:
			sourceCount++
		case source.VersionTarget:
			targetCount++
		}
	}
	if sourceCount == 0 || targetCount == 0 {
		return result, errors.New("Provider Change Contract requires pinned source and target evidence")
	}

	changes := map[string]SemanticChange{}
	detectedAmbiguities := []Ambiguity{}
	for _, declaration := range request.Declarations {
		change := declaration.SemanticChange
		if err := validateChange(change, evidenceByID); err != nil {
			return result, err
		}
		if prior, exists := changes[change.ChangeID]; exists && !reflect.DeepEqual(prior, change) {
			result.Report.Findings = append(result.Report.Findings, Finding{
				Code: "conflicting_semantics", Message: "evidence declarations disagree for one stable change identity",
				EvidenceIDs: sortedUnique(append(append([]string{}, prior.SourceEvidenceIDs...), change.SourceEvidenceIDs...)),
				ChangeIDs:   []string{change.ChangeID},
			})
			detectedAmbiguities = append(detectedAmbiguities, Ambiguity{
				AmbiguityID:        "ambiguity.conflict." + stableFragment(change.ChangeID),
				Description:        "Evidence declarations disagree for change " + change.ChangeID + ".",
				RequiredResolution: "provider_clarification",
				AffectedChangeIDs:  []string{change.ChangeID},
			})
			continue
		}
		changes[change.ChangeID] = change
	}
	if len(changes) == 0 {
		return result, errors.New("Provider Change Contract requires at least one source-bound semantic change")
	}

	semanticChanges := make([]SemanticChange, 0, len(changes))
	for _, change := range changes {
		change.SourceEvidenceIDs = sortedUnique(change.SourceEvidenceIDs)
		change.TargetEvidenceIDs = sortedUnique(change.TargetEvidenceIDs)
		semanticChanges = append(semanticChanges, change)
	}
	sort.Slice(semanticChanges, func(i, j int) bool { return semanticChanges[i].ChangeID < semanticChanges[j].ChangeID })

	evidence := make([]SourceEvidence, 0, len(request.Evidence))
	urls := make([]string, 0, len(request.Evidence))
	digests := make([]string, 0, len(request.Evidence))
	for _, item := range request.Evidence {
		evidence = append(evidence, SourceEvidence{
			EvidenceID: item.Artifact.ID, Kind: string(item.Artifact.Kind),
			VersionRole: string(item.Artifact.VersionRole), Version: item.Artifact.Version,
			Location: item.Artifact.Location, Digest: item.Artifact.Digest,
			License: item.Artifact.License, SourceLocator: item.SourceLocator,
		})
		urls = append(urls, item.Artifact.Location)
		digests = append(digests, item.Artifact.Digest)
	}
	sort.Slice(evidence, func(i, j int) bool { return evidence[i].EvidenceID < evidence[j].EvidenceID })

	ambiguities := append([]Ambiguity{}, request.Ambiguities...)
	ambiguities = append(ambiguities, detectedAmbiguities...)
	sort.Slice(ambiguities, func(i, j int) bool { return ambiguities[i].AmbiguityID < ambiguities[j].AmbiguityID })
	for _, ambiguity := range ambiguities {
		if !stableIDPattern.MatchString(ambiguity.AmbiguityID) || strings.TrimSpace(ambiguity.Description) == "" {
			return result, errors.New("ambiguity requires a stable id and description")
		}
		switch ambiguity.RequiredResolution {
		case "provider_clarification", "consumer_business_input", "unsupported":
		default:
			return result, fmt.Errorf("unsupported ambiguity resolution %q", ambiguity.RequiredResolution)
		}
		result.Report.Findings = append(result.Report.Findings, Finding{
			Code: "semantic_ambiguity", Message: ambiguity.Description,
			ChangeIDs: sortedUnique(ambiguity.AffectedChangeIDs), EvidenceIDs: []string{},
		})
	}
	sort.Slice(result.Report.Findings, func(i, j int) bool {
		return result.Report.Findings[i].Code+"|"+strings.Join(result.Report.Findings[i].ChangeIDs, ",") <
			result.Report.Findings[j].Code+"|"+strings.Join(result.Report.Findings[j].ChangeIDs, ",")
	})

	unsupported := append([]UnsupportedCase{}, request.UnsupportedCases...)
	sort.Slice(unsupported, func(i, j int) bool { return unsupported[i].CaseID < unsupported[j].CaseID })
	contract := Contract{
		ObjectType: "lumyn.migration_pack", SchemaVersion: SchemaVersion,
		MigrationPackID: request.MigrationPackID, ContractVersion: request.ContractVersion,
		APIProviderID: request.APIProviderID, APIOrSDK: request.APIOrSDK,
		SourceVersion: request.SourceVersion, TargetVersion: request.TargetVersion,
		Audience: request.Audience, SourceEvidence: evidence, SemanticChanges: semanticChanges,
		Ambiguities: ambiguities, UnsupportedCases: unsupported,
		VerificationGuidance: VerificationGuidance{
			RequiredBehaviors:       sortedUnique(request.VerificationGuidance),
			ProviderCommandsAllowed: false,
		},
		Provenance: Provenance{
			Class: "public_derived", SourceURLs: sortedUnique(urls),
			SourceDigests: sortedUnique(digests), NormalizedAt: request.NormalizedAt.UTC(),
		},
		ProviderConfirmation: ProviderConfirmation{Status: "not_confirmed"},
		Lifecycle:            request.Lifecycle,
		NonExecutable:        true, GrantsConsumerAuthority: false,
		ProductionCredentialsAllowed: false, ProductionMutationAllowed: false,
	}
	contract.ContractDigest = ComputeDigest(contract)
	result.Contract = contract
	if len(result.Report.Findings) > 0 {
		return result, &BlockedError{Findings: result.Report.Findings}
	}
	if err := ValidateContract(contract); err != nil {
		return result, fmt.Errorf("validate normalized Provider Change Contract: %w", err)
	}
	result.Report.Status = "ready_for_provider_review"
	result.Report.Actionable = true
	return result, nil
}

func validateRequestEnvelope(request NormalizeRequest) error {
	for name, value := range map[string]string{
		"migration pack id": request.MigrationPackID, "API provider id": request.APIProviderID,
		"audience id": request.Audience.AudienceID,
	} {
		if !stableIDPattern.MatchString(value) {
			return fmt.Errorf("%s is not a stable identifier", name)
		}
	}
	if request.ContractVersion < 1 || request.NormalizedAt.IsZero() {
		return errors.New("contract version and normalization time are required")
	}
	if request.SourceVersion == "" || request.TargetVersion == "" || request.SourceVersion == request.TargetVersion {
		return errors.New("distinct source and target versions are required")
	}
	if request.APIOrSDK.Kind != "api" && request.APIOrSDK.Kind != "sdk" {
		return errors.New("API or SDK kind must be api or sdk")
	}
	if request.APIOrSDK.Name == "" || request.APIOrSDK.Package == "" || request.Audience.PackageVersionSelector == "" {
		return errors.New("API or SDK identity and audience selector are required")
	}
	if request.Lifecycle.State != "active" || request.Lifecycle.SupersededBy != nil || request.Lifecycle.WithdrawalReason != nil {
		return errors.New("new normalized contract must begin active and unsuperseded")
	}
	if len(request.VerificationGuidance) == 0 {
		return errors.New("verification guidance is required")
	}
	return nil
}

func validateEvidence(item SourceEvidenceInput, request NormalizeRequest) (SourceEvidenceInput, error) {
	if !stableIDPattern.MatchString(item.Artifact.ID) || !digestPattern.MatchString(item.Artifact.Digest) {
		return SourceEvidenceInput{}, errors.New("source evidence requires a pinned SHA-256 digest")
	}
	pinned, err := source.PinArtifact(item.Artifact, item.Artifact.Digest, request.NormalizedAt)
	if err != nil {
		return SourceEvidenceInput{}, fmt.Errorf("pin source evidence %s: %w", item.Artifact.ID, err)
	}
	item.Artifact = pinned
	switch item.Artifact.VersionRole {
	case source.VersionSource:
		if item.Artifact.Version != request.SourceVersion {
			return SourceEvidenceInput{}, fmt.Errorf("source evidence %s version does not match contract source version", item.Artifact.ID)
		}
	case source.VersionTarget:
		if item.Artifact.Version != request.TargetVersion {
			return SourceEvidenceInput{}, fmt.Errorf("target evidence %s version does not match contract target version", item.Artifact.ID)
		}
	case source.VersionGuidance:
	default:
		return SourceEvidenceInput{}, fmt.Errorf("unsupported source evidence version role %q", item.Artifact.VersionRole)
	}
	if strings.TrimSpace(item.SourceLocator) == "" {
		return SourceEvidenceInput{}, fmt.Errorf("source evidence %s requires a concrete source locator", item.Artifact.ID)
	}
	return item, nil
}

func validateChange(change SemanticChange, evidence map[string]SourceEvidenceInput) error {
	if !stableIDPattern.MatchString(change.ChangeID) || !allowedChangeKinds[change.Kind] {
		return errors.New("semantic change requires a stable id and supported kind")
	}
	if change.SourceSymbol == "" || change.Intent == "" || change.Applicability == "" {
		return fmt.Errorf("semantic change %s is incomplete", change.ChangeID)
	}
	if !allowedMappingModes[change.Mapping.Mode] || len(change.Mapping.DeclarativeSteps) == 0 {
		return fmt.Errorf("semantic change %s requires a supported declarative mapping", change.ChangeID)
	}
	for _, step := range change.Mapping.DeclarativeSteps {
		if strings.TrimSpace(step) == "" || containsExecutableDirective(step) {
			return fmt.Errorf("semantic change %s contains an executable or empty mapping step", change.ChangeID)
		}
	}
	if len(change.SourceEvidenceIDs) == 0 || len(change.TargetEvidenceIDs) == 0 {
		return fmt.Errorf("semantic change %s must cite both source and target evidence", change.ChangeID)
	}
	for _, id := range change.SourceEvidenceIDs {
		item, ok := evidence[id]
		if !ok || item.Artifact.VersionRole != source.VersionSource {
			return fmt.Errorf("semantic change %s has invalid source evidence %s", change.ChangeID, id)
		}
	}
	for _, id := range change.TargetEvidenceIDs {
		item, ok := evidence[id]
		if !ok || item.Artifact.VersionRole != source.VersionTarget {
			return fmt.Errorf("semantic change %s has invalid target evidence %s", change.ChangeID, id)
		}
	}
	return nil
}

func containsExecutableDirective(value string) bool {
	lower := strings.ToLower(value)
	for _, token := range []string{"#!", "<script", "javascript:", "sh -c ", "bash -c ", "powershell ", "curl |", "wget |", "$("} {
		if strings.Contains(lower, token) {
			return true
		}
	}
	return false
}

func sortedUnique(values []string) []string {
	result := append([]string(nil), values...)
	sort.Strings(result)
	return slices.Compact(result)
}

func CanonicalBytes(contract Contract) ([]byte, error) {
	return json.Marshal(contract)
}

// ComputeDigest hashes the semantic contract with self-referential digest and
// separately recorded provider confirmation fields cleared.
func ComputeDigest(contract Contract) string {
	contract.ContractDigest = ""
	contract.ProviderConfirmation = ProviderConfirmation{Status: "not_confirmed"}
	data, _ := json.Marshal(contract)
	return source.DigestBytes(data)
}

func ValidateContract(contract Contract) error {
	if contract.ObjectType != "lumyn.migration_pack" || contract.SchemaVersion != SchemaVersion {
		return errors.New("unsupported Provider Change Contract type or schema version")
	}
	for name, value := range map[string]string{
		"migration pack id": contract.MigrationPackID,
		"API provider id":   contract.APIProviderID,
		"audience id":       contract.Audience.AudienceID,
	} {
		if !stableIDPattern.MatchString(value) {
			return fmt.Errorf("%s is not a stable identifier", name)
		}
	}
	if contract.ContractVersion < 1 || contract.SourceVersion == "" || contract.TargetVersion == "" ||
		contract.SourceVersion == contract.TargetVersion || contract.Audience.PackageVersionSelector == "" {
		return errors.New("Provider Change Contract versions and audience selector are invalid")
	}
	if (contract.APIOrSDK.Kind != "api" && contract.APIOrSDK.Kind != "sdk") ||
		contract.APIOrSDK.Name == "" || contract.APIOrSDK.Package == "" {
		return errors.New("Provider Change Contract API or SDK identity is invalid")
	}
	if !contract.NonExecutable || contract.GrantsConsumerAuthority || contract.ProductionCredentialsAllowed || contract.ProductionMutationAllowed {
		return errors.New("Provider Change Contract must remain non-executable and authority-free")
	}
	evidenceByID := make(map[string]SourceEvidenceInput, len(contract.SourceEvidence))
	sourceCount, targetCount := 0, 0
	urls := make([]string, 0, len(contract.SourceEvidence))
	digests := make([]string, 0, len(contract.SourceEvidence))
	for _, evidence := range contract.SourceEvidence {
		if !stableIDPattern.MatchString(evidence.EvidenceID) || !digestPattern.MatchString(evidence.Digest) ||
			evidence.Version == "" || evidence.License == "" || evidence.SourceLocator == "" {
			return errors.New("Provider Change Contract source evidence is incomplete")
		}
		if _, duplicate := evidenceByID[evidence.EvidenceID]; duplicate {
			return fmt.Errorf("duplicate Provider Change Contract evidence id %s", evidence.EvidenceID)
		}
		parsed, err := url.Parse(evidence.Location)
		if err != nil || parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil ||
			parsed.RawQuery != "" || parsed.Fragment != "" {
			return fmt.Errorf("Provider Change Contract evidence %s has invalid provenance URL", evidence.EvidenceID)
		}
		kind := source.ArtifactKind(evidence.Kind)
		switch kind {
		case source.ArtifactOpenAPI, source.ArtifactDocumentation, source.ArtifactSDKTypes, source.ArtifactMigrationGuide:
		default:
			return fmt.Errorf("Provider Change Contract evidence %s has unsupported kind", evidence.EvidenceID)
		}
		role := source.VersionRole(evidence.VersionRole)
		switch role {
		case source.VersionSource:
			if evidence.Version != contract.SourceVersion {
				return fmt.Errorf("Provider Change Contract evidence %s has wrong source version", evidence.EvidenceID)
			}
			sourceCount++
		case source.VersionTarget:
			if evidence.Version != contract.TargetVersion {
				return fmt.Errorf("Provider Change Contract evidence %s has wrong target version", evidence.EvidenceID)
			}
			targetCount++
		case source.VersionGuidance:
		default:
			return fmt.Errorf("Provider Change Contract evidence %s has unsupported version role", evidence.EvidenceID)
		}
		evidenceByID[evidence.EvidenceID] = SourceEvidenceInput{
			Artifact:      source.PinnedArtifact{ID: evidence.EvidenceID, Kind: kind, VersionRole: role},
			SourceLocator: evidence.SourceLocator,
		}
		urls = append(urls, evidence.Location)
		digests = append(digests, evidence.Digest)
	}
	if sourceCount == 0 || targetCount == 0 || len(contract.SemanticChanges) == 0 {
		return errors.New("Provider Change Contract requires source, target, and semantic-change evidence")
	}
	seenChanges := map[string]bool{}
	for _, change := range contract.SemanticChanges {
		if seenChanges[change.ChangeID] {
			return fmt.Errorf("duplicate Provider Change Contract change id %s", change.ChangeID)
		}
		seenChanges[change.ChangeID] = true
		if err := validateChange(change, evidenceByID); err != nil {
			return err
		}
	}
	if contract.VerificationGuidance.ProviderCommandsAllowed || len(contract.VerificationGuidance.RequiredBehaviors) == 0 {
		return errors.New("Provider Change Contract verification guidance is invalid")
	}
	for _, unsupported := range contract.UnsupportedCases {
		if !stableIDPattern.MatchString(unsupported.CaseID) || unsupported.Description == "" || unsupported.Reason == "" {
			return errors.New("Provider Change Contract unsupported case is invalid")
		}
	}
	if contract.Provenance.NormalizedAt.IsZero() ||
		!reflect.DeepEqual(contract.Provenance.SourceURLs, sortedUnique(urls)) ||
		!reflect.DeepEqual(contract.Provenance.SourceDigests, sortedUnique(digests)) {
		return errors.New("Provider Change Contract provenance does not exactly bind source evidence")
	}
	if contract.ContractDigest != ComputeDigest(contract) {
		return errors.New("Provider Change Contract digest does not bind its semantic content")
	}
	switch contract.Provenance.Class {
	case "public_derived":
		if contract.ProviderConfirmation.Status != "not_confirmed" ||
			contract.ProviderConfirmation.ProviderOperatorID != nil ||
			contract.ProviderConfirmation.ConfirmedContractDigest != nil ||
			contract.ProviderConfirmation.ConfirmedAt != nil {
			return errors.New("public-derived contract must not imply provider confirmation")
		}
	case "provider_confirmed":
		confirmation := contract.ProviderConfirmation
		if confirmation.Status != "provider_confirmed" || confirmation.ProviderOperatorID == nil ||
			confirmation.ConfirmedContractDigest == nil || confirmation.ConfirmedAt == nil ||
			*confirmation.ConfirmedContractDigest != contract.ContractDigest {
			return errors.New("provider-confirmed contract has invalid confirmation binding")
		}
	default:
		return errors.New("unsupported Provider Change Contract provenance class")
	}
	if contract.Lifecycle.State != "active" {
		return fmt.Errorf("Provider Change Contract lifecycle %q is not actionable", contract.Lifecycle.State)
	}
	if len(contract.Ambiguities) > 0 {
		return errors.New("Provider Change Contract contains unresolved ambiguity")
	}
	return nil
}

type ConfirmationRecord struct {
	ObjectType            string    `json:"object_type"`
	SchemaVersion         string    `json:"schema_version"`
	MigrationPackID       string    `json:"migration_pack_id"`
	ContractVersion       int       `json:"contract_version"`
	ContractDigest        string    `json:"contract_digest"`
	ProviderOperatorID    string    `json:"provider_operator_id"`
	ConfirmedAt           time.Time `json:"confirmed_at"`
	ConfirmationSignature string    `json:"confirmation_signature"`
}

func Confirm(contract Contract, providerOperatorID string, confirmedAt time.Time, privateKey ed25519.PrivateKey) (Contract, ConfirmationRecord, error) {
	if contract.Provenance.Class != "public_derived" || contract.ProviderConfirmation.Status != "not_confirmed" {
		return Contract{}, ConfirmationRecord{}, errors.New("only an unconfirmed public-derived contract can be confirmed")
	}
	if err := ValidateContract(contract); err != nil {
		return Contract{}, ConfirmationRecord{}, fmt.Errorf("confirm Provider Change Contract: %w", err)
	}
	if !stableIDPattern.MatchString(providerOperatorID) || confirmedAt.IsZero() || len(privateKey) != ed25519.PrivateKeySize {
		return Contract{}, ConfirmationRecord{}, errors.New("provider confirmation identity, time, and Ed25519 key are required")
	}
	if confirmedAt.Before(contract.Provenance.NormalizedAt) {
		return Contract{}, ConfirmationRecord{}, errors.New("provider confirmation cannot predate normalization")
	}
	if len(contract.Ambiguities) > 0 {
		return Contract{}, ConfirmationRecord{}, errors.New("ambiguous contract cannot be provider-confirmed")
	}
	contract.Provenance.Class = "provider_confirmed"
	contract.ContractDigest = ComputeDigest(contract)
	digest := contract.ContractDigest
	at := confirmedAt.UTC()
	contract.ProviderConfirmation = ProviderConfirmation{
		Status: "provider_confirmed", ProviderOperatorID: &providerOperatorID,
		ConfirmedContractDigest: &digest, ConfirmedAt: &at,
	}
	record := ConfirmationRecord{
		ObjectType: "lumyn.provider_confirmation", SchemaVersion: SchemaVersion,
		MigrationPackID: contract.MigrationPackID, ContractVersion: contract.ContractVersion,
		ContractDigest: contract.ContractDigest, ProviderOperatorID: providerOperatorID, ConfirmedAt: at,
	}
	unsigned, _ := confirmationBytes(record)
	record.ConfirmationSignature = base64.RawURLEncoding.EncodeToString(ed25519.Sign(privateKey, unsigned))
	return contract, record, nil
}

func VerifyConfirmation(contract Contract, record ConfirmationRecord, publicKey ed25519.PublicKey) error {
	if len(publicKey) != ed25519.PublicKeySize {
		return errors.New("provider confirmation Ed25519 public key is invalid")
	}
	if err := ValidateContract(contract); err != nil {
		return err
	}
	if record.MigrationPackID != contract.MigrationPackID || record.ContractVersion != contract.ContractVersion ||
		record.ContractDigest != contract.ContractDigest || contract.ProviderConfirmation.ProviderOperatorID == nil ||
		record.ProviderOperatorID != *contract.ProviderConfirmation.ProviderOperatorID ||
		contract.ProviderConfirmation.ConfirmedAt == nil || !record.ConfirmedAt.Equal(*contract.ProviderConfirmation.ConfirmedAt) {
		return errors.New("provider confirmation record does not bind the exact contract")
	}
	signature, err := base64.RawURLEncoding.DecodeString(record.ConfirmationSignature)
	if err != nil {
		return errors.New("provider confirmation signature is not base64url")
	}
	unsigned, _ := confirmationBytes(record)
	if !ed25519.Verify(publicKey, unsigned, signature) {
		return errors.New("provider confirmation signature is invalid")
	}
	return nil
}

func confirmationBytes(record ConfirmationRecord) ([]byte, error) {
	record.ConfirmationSignature = ""
	return json.Marshal(record)
}
