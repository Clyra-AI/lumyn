// Package status validates the deliberately narrow API-provider rollout view.
package status

import (
	"errors"
	"fmt"
	"strings"
)

type ProviderProjection struct {
	EventID             string   `json:"event_id"`
	EventDigest         string   `json:"event_digest"`
	RunID               string   `json:"run_id"`
	InstallationDigest  string   `json:"consumer_installation_digest"`
	AuthorizationDigest string   `json:"event_authorization_digest"`
	PlanDigest          string   `json:"migration_plan_digest,omitempty"`
	CandidateDigest     string   `json:"candidate_digest,omitempty"`
	VerificationDigest  string   `json:"verification_digest,omitempty"`
	DeliveryDigest      string   `json:"delivery_digest,omitempty"`
	EvidenceKinds       []string `json:"evidence_kinds"`
	State               string   `json:"state"`
	Provenance          string   `json:"provenance"`
	Consented           bool     `json:"consented"`
}

var providerStates = map[string]struct{}{
	"unknown": {}, "received": {}, "not_applicable": {}, "affected": {},
	"needs_input": {}, "candidate_ready": {}, "verified": {},
	"draft_pr_open": {}, "accepted": {}, "merged": {}, "retired": {},
}

var projectionProvenance = map[string]struct{}{
	"observed": {}, "consumer_reported": {}, "unknown": {},
}

var requiredEvidenceKind = map[string]string{
	"received":        "event_receipt",
	"not_applicable":  "explicit_not_applicable",
	"affected":        "impact_outcome",
	"needs_input":     "consumer_input_request",
	"candidate_ready": "candidate_outcome",
	"verified":        "verification_outcome",
	"draft_pr_open":   "draft_pr_outcome",
	"accepted":        "consumer_acceptance",
	"merged":          "merge_outcome",
	"retired":         "retirement_confirmation",
}

// MigrationAxes deliberately has no aggregate status field. Callers must
// render every axis so a strong delivery or candidate label cannot hide a
// weaker impact, route, or verification result.
type MigrationAxes struct {
	Impact       string `json:"impact"`
	Route        string `json:"route"`
	Candidate    string `json:"candidate"`
	Verification string `json:"verification"`
	Delivery     string `json:"delivery"`
}

var migrationAxisValues = map[string]map[string]struct{}{
	"impact":       set("not_analyzed", "unaffected", "affected_supported", "affected_needs_input", "unsupported", "uncertain"),
	"route":        set("not_routed", "deterministic", "agent_assisted", "manual", "blocked"),
	"candidate":    set("not_attempted", "planned", "candidate_generated", "repairing", "needs_input", "failed", "stale"),
	"verification": set("not_run", "static_verified", "repo_verified", "workflow_contract_replay_passed", "workflow_verified_replay", "workflow_verified_mock", "workflow_verified_sandbox", "partial", "failed", "gap", "stale"),
	"delivery":     set("not_requested", "patch_exported", "local_branch_ready", "pr_bundle_ready", "remote_branch_pushed", "draft_pr_open", "consumer_accepted", "merged", "closed", "blocked", "superseded"),
}

func ValidateMigrationAxes(value MigrationAxes) error {
	values := map[string]string{
		"impact": value.Impact, "route": value.Route, "candidate": value.Candidate,
		"verification": value.Verification, "delivery": value.Delivery,
	}
	for axis, state := range values {
		if _, ok := migrationAxisValues[axis][state]; !ok {
			return fmt.Errorf("unknown %s status %q", axis, state)
		}
	}
	return nil
}

func set(values ...string) map[string]struct{} {
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		result[value] = struct{}{}
	}
	return result
}

func ValidateProviderProjection(value ProviderProjection) error {
	if strings.TrimSpace(value.EventID) == "" || strings.TrimSpace(value.EventDigest) == "" {
		return errors.New("provider projection requires exact event identity and digest")
	}
	if strings.TrimSpace(value.RunID) == "" || strings.TrimSpace(value.InstallationDigest) == "" || strings.TrimSpace(value.AuthorizationDigest) == "" {
		return errors.New("provider projection requires exact run, installation, and authorization bindings")
	}
	if _, ok := providerStates[value.State]; !ok {
		return fmt.Errorf("unknown provider projection state %q", value.State)
	}
	if _, ok := projectionProvenance[value.Provenance]; !ok {
		return fmt.Errorf("unknown provider projection provenance %q", value.Provenance)
	}
	if !value.Consented {
		return errors.New("provider projection requires consumer consent")
	}
	if value.Provenance == "unknown" && value.State != "unknown" {
		return errors.New("silence or unknown provenance cannot imply a rollout state")
	}
	if value.State == "unknown" {
		if len(value.EvidenceKinds) != 0 {
			return errors.New("unknown provider projection cannot carry inferred evidence")
		}
		if value.PlanDigest != "" || value.CandidateDigest != "" || value.VerificationDigest != "" || value.DeliveryDigest != "" {
			return errors.New("unknown provider projection cannot claim later-artifact bindings")
		}
		return nil
	}
	if required := requiredEvidenceKind[value.State]; !contains(value.EvidenceKinds, required) {
		return fmt.Errorf("provider projection state %q requires %s evidence", value.State, required)
	}
	if err := validateBindingDependencies(value); err != nil {
		return err
	}
	if err := requireStatusBindings(value); err != nil {
		return err
	}
	return nil
}

func requireStatusBindings(value ProviderProjection) error {
	require := func(name, digest string) error {
		if strings.TrimSpace(digest) == "" {
			return fmt.Errorf("provider projection state %q requires exact %s binding", value.State, name)
		}
		return nil
	}
	switch value.State {
	case "received", "not_applicable", "affected":
		if value.PlanDigest != "" || value.CandidateDigest != "" || value.VerificationDigest != "" || value.DeliveryDigest != "" {
			return fmt.Errorf("provider projection state %q cannot claim later-artifact bindings", value.State)
		}
	case "candidate_ready":
		if err := require("plan", value.PlanDigest); err != nil {
			return err
		}
		if err := require("candidate", value.CandidateDigest); err != nil {
			return err
		}
		if value.VerificationDigest != "" || value.DeliveryDigest != "" {
			return errors.New("candidate_ready cannot claim verification or delivery bindings")
		}
		return nil
	case "verified":
		for name, digest := range map[string]string{
			"plan": value.PlanDigest, "candidate": value.CandidateDigest, "verification": value.VerificationDigest,
		} {
			if err := require(name, digest); err != nil {
				return err
			}
		}
		if value.DeliveryDigest != "" {
			return errors.New("verified cannot claim a delivery binding")
		}
	case "draft_pr_open", "accepted", "merged", "retired":
		for name, digest := range map[string]string{
			"plan": value.PlanDigest, "candidate": value.CandidateDigest,
			"verification": value.VerificationDigest, "delivery": value.DeliveryDigest,
		} {
			if err := require(name, digest); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateBindingDependencies(value ProviderProjection) error {
	if value.CandidateDigest != "" && value.PlanDigest == "" {
		return errors.New("candidate binding requires an exact plan binding")
	}
	if value.VerificationDigest != "" && (value.PlanDigest == "" || value.CandidateDigest == "") {
		return errors.New("verification binding requires exact plan and candidate bindings")
	}
	if value.DeliveryDigest != "" && (value.PlanDigest == "" || value.CandidateDigest == "" || value.VerificationDigest == "") {
		return errors.New("delivery binding requires exact plan, candidate, and verification bindings")
	}
	return nil
}

func contains(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}
