// Package redaction enforces disclosure boundaries before persistence or
// provider sharing.
package redaction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
)

var providerVisibleFields = map[string]struct{}{
	"projection_id": {}, "api_provider_id": {},
	"provider_change_event_id": {}, "provider_change_contract_digest": {},
	"consumer_installation_pseudonym": {}, "bindings": {}, "status": {}, "provenance": {},
	"evidence_refs": {}, "observed_at": {},
}

var prohibitedProviderFields = map[string]struct{}{
	"raw_source": {}, "source": {}, "source_code": {}, "code": {},
	"raw_diff": {}, "diff": {}, "raw_patch": {}, "patch": {},
	"prompt": {}, "raw_prompt": {}, "response": {}, "model_response": {},
	"tool_trace": {}, "tool_calls": {}, "log": {}, "logs": {}, "trace": {},
	"session": {}, "agent_session": {}, "credential": {}, "credentials": {},
	"secret": {}, "secrets": {}, "token": {}, "access_token": {}, "refresh_token": {},
}

var (
	identifierPattern    = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{2,127}$`)
	digestPattern        = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
	commonSecretPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)\b(?:sk|rk)_live_[A-Za-z0-9]{16,}\b`),
		regexp.MustCompile(`\bgh[pousr]_[A-Za-z0-9]{20,}\b`),
		regexp.MustCompile(`\b(?:AKIA|ASIA)[0-9A-Z]{16}\b`),
		regexp.MustCompile(`(?i)\bsk-(?:proj-)?[A-Za-z0-9_-]{16,}\b`),
		regexp.MustCompile(`(?i)\bbearer\s+[A-Za-z0-9._~+/=-]{16,}\b`),
		regexp.MustCompile(`(?i)\b(?:password|passwd|api[_-]?key|access[_-]?token|refresh[_-]?token)\s*[:=]\s*\S+`),
		regexp.MustCompile(`-----BEGIN(?: [A-Z0-9]+)? PRIVATE KEY-----`),
	}
)

type providerStatusProjection struct {
	ObjectType                    string                       `json:"object_type"`
	SchemaVersion                 string                       `json:"schema_version"`
	ProjectionID                  string                       `json:"projection_id"`
	APIProviderID                 string                       `json:"api_provider_id"`
	ProviderChangeEventID         string                       `json:"provider_change_event_id"`
	ProviderChangeContractDigest  string                       `json:"provider_change_contract_digest"`
	ConsumerInstallationPseudonym string                       `json:"consumer_installation_pseudonym"`
	Bindings                      *providerBindings            `json:"bindings"`
	Status                        string                       `json:"status"`
	Provenance                    string                       `json:"provenance"`
	EvidenceRefs                  []providerEvidenceRef        `json:"evidence_refs"`
	Consent                       *providerConsent             `json:"consent"`
	ProviderVisibleFields         []string                     `json:"provider_visible_fields"`
	InterpretationRules           *providerInterpretationRules `json:"interpretation_rules"`
	Privacy                       *providerPrivacy             `json:"privacy"`
	ObservedAt                    string                       `json:"observed_at"`
	ArtifactDigest                string                       `json:"artifact_digest"`
}

type providerEvidenceRef struct {
	EvidenceKind   string `json:"evidence_kind"`
	EvidenceDigest string `json:"evidence_digest"`
	Provenance     string `json:"provenance"`
}

type providerBindings struct {
	RunID                      string         `json:"run_id"`
	ConsumerInstallationDigest string         `json:"consumer_installation_digest"`
	EventAuthorizationDigest   string         `json:"event_authorization_digest"`
	MigrationPlanDigest        nullableDigest `json:"migration_plan_digest"`
	CandidateDigest            nullableDigest `json:"candidate_digest"`
	VerificationDigest         nullableDigest `json:"verification_digest"`
	DeliveryDigest             nullableDigest `json:"delivery_digest"`
}

type nullableDigest struct {
	Present bool
	Value   *string
}

func (digest *nullableDigest) UnmarshalJSON(data []byte) error {
	digest.Present = true
	if bytes.Equal(data, []byte("null")) {
		digest.Value = nil
		return nil
	}
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("digest must be a string or null: %w", err)
	}
	digest.Value = &value
	return nil
}

type providerConsent struct {
	ConsumerConsentID     string `json:"consumer_consent_id"`
	ConsumerConsentDigest string `json:"consumer_consent_digest"`
	EventBound            *bool  `json:"event_bound"`
	FieldsEnumerated      *bool  `json:"fields_enumerated"`
	Active                *bool  `json:"active"`
	CheckedAt             string `json:"checked_at"`
}

type providerInterpretationRules struct {
	Silence                               string `json:"silence"`
	MergedImpliesRetired                  *bool  `json:"merged_implies_retired"`
	NotApplicableRequiresExplicitEvidence *bool  `json:"not_applicable_requires_explicit_evidence"`
	UnaffectedMayBeInferred               *bool  `json:"unaffected_may_be_inferred"`
}

type providerPrivacy struct {
	Classification                string `json:"classification"`
	RedactionPolicyDigest         string `json:"redaction_policy_digest"`
	RedactionStatus               string `json:"redaction_status"`
	SafeToShareWithAPIProvider    *bool  `json:"safe_to_share_with_api_provider"`
	RawConsumerCodeIncluded       *bool  `json:"raw_consumer_code_included"`
	RawDiffIncluded               *bool  `json:"raw_diff_included"`
	RawLogsOrTracesIncluded       *bool  `json:"raw_logs_or_traces_included"`
	RawPromptsOrResponsesIncluded *bool  `json:"raw_prompts_or_responses_included"`
	AgentSessionDataIncluded      *bool  `json:"agent_session_data_included"`
	CredentialsIncluded           *bool  `json:"credentials_included"`
}

// ValidateProviderPayload accepts only the fixed, typed
// lumyn.provider_status_projection vocabulary. consentedFields is an optional
// external ceiling over provider_visible_fields; it cannot add extension
// fields. The secret scan recognizes common high-signal credential formats but
// is not a general content classifier. The typed vocabulary deliberately has
// no provider-visible free-text field in which raw consumer material can hide.
func ValidateProviderPayload(payload map[string]any, consentedFields []string) error {
	if payload == nil {
		return fmt.Errorf("provider projection is required")
	}
	if err := inspectProviderValue(payload, "projection"); err != nil {
		return err
	}

	serialized, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("serialize provider projection: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(serialized))
	decoder.DisallowUnknownFields()
	var projection providerStatusProjection
	if err := decoder.Decode(&projection); err != nil {
		return fmt.Errorf("decode typed provider projection: %w", err)
	}
	if err := ensureJSONEOF(decoder); err != nil {
		return err
	}
	if err := validateProviderProjection(projection); err != nil {
		return err
	}
	return validateExternalConsent(projection.ProviderVisibleFields, consentedFields)
}

func inspectProviderValue(value any, path string) error {
	switch typed := value.(type) {
	case map[string]any:
		for field, child := range typed {
			if _, prohibited := prohibitedProviderFields[normalizeField(field)]; prohibited {
				return fmt.Errorf("consumer-private field %q is prohibited at %s", field, path)
			}
			if err := inspectProviderValue(child, path+"."+field); err != nil {
				return err
			}
		}
	case []any:
		for index, child := range typed {
			if err := inspectProviderValue(child, fmt.Sprintf("%s[%d]", path, index)); err != nil {
				return err
			}
		}
	case string:
		for _, pattern := range commonSecretPatterns {
			if pattern.MatchString(typed) {
				return fmt.Errorf("common secret pattern is prohibited at %s", path)
			}
		}
	case nil:
		// Later-artifact binding digests use explicit null before the artifact
		// exists. Typed validation below rejects null everywhere else.
	case bool, float64, json.Number:
		// These are the only non-container values produced by decoding JSON.
	default:
		return fmt.Errorf("non-JSON value of type %T is prohibited at %s", value, path)
	}
	return nil
}

func ensureJSONEOF(decoder *json.Decoder) error {
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return fmt.Errorf("provider projection contains trailing JSON")
		}
		return fmt.Errorf("decode provider projection trailer: %w", err)
	}
	return nil
}

func validateProviderProjection(projection providerStatusProjection) error {
	if projection.ObjectType != "lumyn.provider_status_projection" {
		return fmt.Errorf("object_type must be lumyn.provider_status_projection")
	}
	if projection.SchemaVersion != "1.0" {
		return fmt.Errorf("schema_version must be 1.0")
	}
	for field, value := range map[string]string{
		"projection_id":                   projection.ProjectionID,
		"api_provider_id":                 projection.APIProviderID,
		"provider_change_event_id":        projection.ProviderChangeEventID,
		"consumer_installation_pseudonym": projection.ConsumerInstallationPseudonym,
	} {
		if !identifierPattern.MatchString(value) {
			return fmt.Errorf("%s is not a valid identifier", field)
		}
	}
	for field, value := range map[string]string{
		"provider_change_contract_digest": projection.ProviderChangeContractDigest,
		"artifact_digest":                 projection.ArtifactDigest,
	} {
		if !digestPattern.MatchString(value) {
			return fmt.Errorf("%s is not a valid SHA-256 digest", field)
		}
	}
	if !oneOf(projection.Status, "unknown", "received", "not_applicable", "affected", "needs_input", "candidate_ready", "verified", "draft_pr_open", "accepted", "merged", "retired") {
		return fmt.Errorf("status %q is not supported", projection.Status)
	}
	if !oneOf(projection.Provenance, "observed", "consumer_reported", "unknown") {
		return fmt.Errorf("provenance %q is not supported", projection.Provenance)
	}
	if projection.EvidenceRefs == nil {
		return fmt.Errorf("evidence_refs is required")
	}
	if err := validateEvidenceRefs(projection.EvidenceRefs); err != nil {
		return err
	}
	if err := validateBindings(projection.Bindings, projection.Status); err != nil {
		return err
	}
	if err := validateConsent(projection.Consent); err != nil {
		return err
	}
	if err := validateProviderVisibleFields(projection.ProviderVisibleFields); err != nil {
		return err
	}
	if err := validateInterpretationRules(projection.InterpretationRules); err != nil {
		return err
	}
	if err := validatePrivacy(projection.Privacy); err != nil {
		return err
	}
	if err := validateTimestamp("observed_at", projection.ObservedAt); err != nil {
		return err
	}

	if projection.Status == "unknown" {
		if projection.Provenance != "unknown" || len(projection.EvidenceRefs) != 0 {
			return fmt.Errorf("unknown status requires unknown provenance and no evidence")
		}
	} else if projection.Provenance == "unknown" || len(projection.EvidenceRefs) == 0 {
		return fmt.Errorf("non-unknown status requires observed or consumer_reported provenance and evidence")
	}
	statusEvidenceKind := map[string]string{
		"received": "event_receipt", "not_applicable": "explicit_not_applicable",
		"affected": "impact_outcome", "needs_input": "consumer_input_request",
		"candidate_ready": "candidate_outcome", "verified": "verification_outcome",
		"draft_pr_open": "draft_pr_outcome", "accepted": "consumer_acceptance",
		"merged": "merge_outcome", "retired": "retirement_confirmation",
	}
	if requiredKind := statusEvidenceKind[projection.Status]; requiredKind != "" && !hasEvidenceKind(projection.EvidenceRefs, requiredKind) {
		return fmt.Errorf("%s status requires %s evidence", projection.Status, requiredKind)
	}
	return nil
}

func validateBindings(bindings *providerBindings, status string) error {
	if bindings == nil {
		return fmt.Errorf("bindings is required")
	}
	if !identifierPattern.MatchString(bindings.RunID) {
		return fmt.Errorf("bindings.run_id is not a valid identifier")
	}
	for field, value := range map[string]string{
		"consumer_installation_digest": bindings.ConsumerInstallationDigest,
		"event_authorization_digest":   bindings.EventAuthorizationDigest,
	} {
		if !digestPattern.MatchString(value) {
			return fmt.Errorf("bindings.%s is not a valid SHA-256 digest", field)
		}
	}
	for field, digest := range map[string]nullableDigest{
		"migration_plan_digest": bindings.MigrationPlanDigest,
		"candidate_digest":      bindings.CandidateDigest,
		"verification_digest":   bindings.VerificationDigest,
		"delivery_digest":       bindings.DeliveryDigest,
	} {
		if !digest.Present {
			return fmt.Errorf("bindings.%s is required", field)
		}
		if digest.Value != nil && !digestPattern.MatchString(*digest.Value) {
			return fmt.Errorf("bindings.%s is not a valid SHA-256 digest", field)
		}
	}

	planSet := bindings.MigrationPlanDigest.Value != nil
	candidateSet := bindings.CandidateDigest.Value != nil
	verificationSet := bindings.VerificationDigest.Value != nil
	deliverySet := bindings.DeliveryDigest.Value != nil
	if candidateSet && !planSet {
		return fmt.Errorf("candidate binding requires migration plan binding")
	}
	if verificationSet && (!planSet || !candidateSet) {
		return fmt.Errorf("verification binding requires migration plan and candidate bindings")
	}
	if deliverySet && (!planSet || !candidateSet || !verificationSet) {
		return fmt.Errorf("delivery binding requires complete plan, candidate, and verification bindings")
	}

	switch status {
	case "unknown", "received", "not_applicable", "affected":
		if planSet || candidateSet || verificationSet || deliverySet {
			return fmt.Errorf("%s status requires early bindings with no later-artifact digests", status)
		}
	case "candidate_ready":
		if !planSet || !candidateSet || verificationSet || deliverySet {
			return fmt.Errorf("candidate_ready status requires plan and candidate bindings only")
		}
	case "verified":
		if !planSet || !candidateSet || !verificationSet || deliverySet {
			return fmt.Errorf("verified status requires plan, candidate, and verification bindings only")
		}
	case "draft_pr_open", "accepted", "merged", "retired":
		if !planSet || !candidateSet || !verificationSet || !deliverySet {
			return fmt.Errorf("%s status requires complete bindings", status)
		}
	}
	return nil
}

func validateEvidenceRefs(refs []providerEvidenceRef) error {
	seen := make(map[providerEvidenceRef]struct{}, len(refs))
	for index, ref := range refs {
		if !oneOf(ref.EvidenceKind, "event_receipt", "explicit_not_applicable", "impact_outcome", "consumer_input_request", "candidate_outcome", "verification_outcome", "draft_pr_outcome", "consumer_acceptance", "merge_outcome", "retirement_confirmation") {
			return fmt.Errorf("evidence_refs[%d].evidence_kind is not supported", index)
		}
		if !digestPattern.MatchString(ref.EvidenceDigest) {
			return fmt.Errorf("evidence_refs[%d].evidence_digest is not a valid SHA-256 digest", index)
		}
		if !oneOf(ref.Provenance, "observed", "consumer_reported") {
			return fmt.Errorf("evidence_refs[%d].provenance is not supported", index)
		}
		if _, duplicate := seen[ref]; duplicate {
			return fmt.Errorf("evidence_refs[%d] duplicates an earlier reference", index)
		}
		seen[ref] = struct{}{}
	}
	return nil
}

func validateConsent(consent *providerConsent) error {
	if consent == nil {
		return fmt.Errorf("consent is required")
	}
	if !identifierPattern.MatchString(consent.ConsumerConsentID) {
		return fmt.Errorf("consent.consumer_consent_id is not a valid identifier")
	}
	if !digestPattern.MatchString(consent.ConsumerConsentDigest) {
		return fmt.Errorf("consent.consumer_consent_digest is not a valid SHA-256 digest")
	}
	if !requiredTrue(consent.EventBound) || !requiredTrue(consent.FieldsEnumerated) || !requiredTrue(consent.Active) {
		return fmt.Errorf("consent must be event-bound, field-enumerated, and active")
	}
	return validateTimestamp("consent.checked_at", consent.CheckedAt)
}

func validateProviderVisibleFields(fields []string) error {
	if len(fields) == 0 {
		return fmt.Errorf("provider_visible_fields must contain at least one fixed field")
	}
	seen := make(map[string]struct{}, len(fields))
	for _, field := range fields {
		if _, allowed := providerVisibleFields[field]; !allowed {
			return fmt.Errorf("provider_visible_fields contains unsupported field %q", field)
		}
		if _, duplicate := seen[field]; duplicate {
			return fmt.Errorf("provider_visible_fields contains duplicate field %q", field)
		}
		seen[field] = struct{}{}
	}
	return nil
}

func validateInterpretationRules(rules *providerInterpretationRules) error {
	if rules == nil {
		return fmt.Errorf("interpretation_rules is required")
	}
	if rules.Silence != "unknown" || !requiredFalse(rules.MergedImpliesRetired) ||
		!requiredTrue(rules.NotApplicableRequiresExplicitEvidence) || !requiredFalse(rules.UnaffectedMayBeInferred) {
		return fmt.Errorf("interpretation_rules must preserve unknown silence, explicit not-applicable evidence, and no inferred retirement or unaffected status")
	}
	return nil
}

func validatePrivacy(privacy *providerPrivacy) error {
	if privacy == nil {
		return fmt.Errorf("privacy is required")
	}
	if privacy.Classification != "provider_consented_projection" ||
		privacy.RedactionStatus != "passed" || !digestPattern.MatchString(privacy.RedactionPolicyDigest) ||
		!requiredTrue(privacy.SafeToShareWithAPIProvider) ||
		!requiredFalse(privacy.RawConsumerCodeIncluded) ||
		!requiredFalse(privacy.RawDiffIncluded) ||
		!requiredFalse(privacy.RawLogsOrTracesIncluded) ||
		!requiredFalse(privacy.RawPromptsOrResponsesIncluded) ||
		!requiredFalse(privacy.AgentSessionDataIncluded) ||
		!requiredFalse(privacy.CredentialsIncluded) {
		return fmt.Errorf("privacy must be a passed provider-consented projection with every private-material flag false")
	}
	return nil
}

func validateExternalConsent(visibleFields, consentedFields []string) error {
	if consentedFields == nil {
		return nil
	}
	consented := make(map[string]struct{}, len(consentedFields))
	for _, field := range consentedFields {
		if _, allowed := providerVisibleFields[field]; !allowed {
			return fmt.Errorf("external consent contains unsupported field %q", field)
		}
		if _, duplicate := consented[field]; duplicate {
			return fmt.Errorf("external consent contains duplicate field %q", field)
		}
		consented[field] = struct{}{}
	}
	if len(consented) != len(visibleFields) {
		return fmt.Errorf("external consent does not exactly match provider_visible_fields")
	}
	for _, field := range visibleFields {
		if _, ok := consented[field]; !ok {
			return fmt.Errorf("provider-visible field %q is absent from external consent", field)
		}
	}
	return nil
}

func validateTimestamp(field, value string) error {
	if _, err := time.Parse(time.RFC3339, value); err != nil {
		return fmt.Errorf("%s is not RFC3339: %w", field, err)
	}
	return nil
}

func hasEvidenceKind(refs []providerEvidenceRef, kind string) bool {
	for _, ref := range refs {
		if ref.EvidenceKind == kind {
			return true
		}
	}
	return false
}

func requiredTrue(value *bool) bool {
	return value != nil && *value
}

func requiredFalse(value *bool) bool {
	return value != nil && !*value
}

func oneOf(value string, options ...string) bool {
	for _, option := range options {
		if value == option {
			return true
		}
	}
	return false
}

func normalizeField(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
