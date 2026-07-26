package evidence

import (
	"errors"
	"fmt"
	"sort"
)

// MigrationBindings contains only cross-artifact digests. The artifacts retain
// their own schemas; this structure makes relationships executable without
// creating aggregate runtime authority.
type MigrationBindings struct {
	RunID                           string
	InstallationDigest              string
	AuthorizationInstallationDigest string
	AuthorizationDigest             string
	PlanAuthorizationDigest         string
	PackDigest                      string
	EventDigest                     string
	EventContractDigest             string
	AuthorizationEventDigest        string
	ImpactDigest                    string
	PlanImpactDigest                string
	PlanDigest                      string
	CandidatePlanDigest             string
	CandidateDigest                 string
	VerificationCandidateDigest     string
	VerificationDigest              string
	ExportCandidateDigest           string
	ExportVerificationDigest        string
	ExportDigest                    string
	ProjectionRunID                 string
	ProjectionEventDigest           string
	ProjectionInstallationDigest    string
	ProjectionAuthorizationDigest   string
	ProjectionPlanDigest            string
	ProjectionCandidateDigest       string
	ProjectionVerificationDigest    string
	ProjectionDeliveryDigest        string
	ProjectionEvidenceDigest        string
}

func ValidateMigrationBindings(value MigrationBindings) error {
	if err := requireBindings(value); err != nil {
		return err
	}
	checks := []struct {
		name        string
		left, right string
	}{
		{"authorization installation", value.AuthorizationInstallationDigest, value.InstallationDigest},
		{"event contract", value.EventContractDigest, value.PackDigest},
		{"authorization event", value.AuthorizationEventDigest, value.EventDigest},
		{"plan authorization", value.PlanAuthorizationDigest, value.AuthorizationDigest},
		{"plan impact", value.PlanImpactDigest, value.ImpactDigest},
		{"candidate plan", value.CandidatePlanDigest, value.PlanDigest},
		{"verification candidate", value.VerificationCandidateDigest, value.CandidateDigest},
		{"export candidate", value.ExportCandidateDigest, value.CandidateDigest},
		{"export verification", value.ExportVerificationDigest, value.VerificationDigest},
		{"provider projection run", value.ProjectionRunID, value.RunID},
		{"provider projection event", value.ProjectionEventDigest, value.EventDigest},
		{"provider projection installation", value.ProjectionInstallationDigest, value.InstallationDigest},
		{"provider projection authorization", value.ProjectionAuthorizationDigest, value.AuthorizationDigest},
		{"provider projection plan", value.ProjectionPlanDigest, value.PlanDigest},
		{"provider projection candidate", value.ProjectionCandidateDigest, value.CandidateDigest},
		{"provider projection verification", value.ProjectionVerificationDigest, value.VerificationDigest},
		{"provider projection delivery", value.ProjectionDeliveryDigest, value.ExportDigest},
		{"provider projection evidence", value.ProjectionEvidenceDigest, value.VerificationDigest},
	}
	for _, check := range checks {
		if check.left != check.right {
			return fmt.Errorf("%s digest binding mismatch", check.name)
		}
	}
	return nil
}

func requireBindings(value MigrationBindings) error {
	fields := map[string]string{
		"run": value.RunID, "installation": value.InstallationDigest,
		"authorization installation": value.AuthorizationInstallationDigest,
		"authorization":              value.AuthorizationDigest, "plan authorization": value.PlanAuthorizationDigest,
		"pack": value.PackDigest, "event": value.EventDigest,
		"event contract": value.EventContractDigest, "authorization event": value.AuthorizationEventDigest,
		"impact": value.ImpactDigest, "plan impact": value.PlanImpactDigest,
		"plan": value.PlanDigest, "candidate plan": value.CandidatePlanDigest,
		"candidate": value.CandidateDigest, "verification candidate": value.VerificationCandidateDigest,
		"verification": value.VerificationDigest, "export candidate": value.ExportCandidateDigest,
		"export verification": value.ExportVerificationDigest, "projection event": value.ProjectionEventDigest,
		"export": value.ExportDigest, "projection run": value.ProjectionRunID,
		"projection installation":  value.ProjectionInstallationDigest,
		"projection authorization": value.ProjectionAuthorizationDigest,
		"projection plan":          value.ProjectionPlanDigest, "projection candidate": value.ProjectionCandidateDigest,
		"projection verification": value.ProjectionVerificationDigest,
		"projection delivery":     value.ProjectionDeliveryDigest, "projection evidence": value.ProjectionEvidenceDigest,
	}
	for name, digest := range fields {
		if digest == "" {
			return errors.New(name + " digest binding is required")
		}
	}
	return nil
}

// ProjectionBindingInputs returns the complete provider-status dependency set.
// Any changed value makes the projection stale; values are opaque identities
// or digests and do not expose consumer source or private execution evidence.
func ProjectionBindingInputs(value MigrationBindings) map[string]string {
	return map[string]string{
		"run":           value.RunID,
		"installation":  value.InstallationDigest,
		"authorization": value.AuthorizationDigest,
		"plan":          value.PlanDigest,
		"candidate":     value.CandidateDigest,
		"verification":  value.VerificationDigest,
		"delivery":      value.ExportDigest,
	}
}

// ChangedBindingInputs returns a stable list of any added, removed, or changed
// binding. A non-empty result invalidates dependent evidence.
func ChangedBindingInputs(previous, current map[string]string) []string {
	keys := make(map[string]struct{}, len(previous)+len(current))
	for key := range previous {
		keys[key] = struct{}{}
	}
	for key := range current {
		keys[key] = struct{}{}
	}
	changed := make([]string, 0)
	for key := range keys {
		if previous[key] != current[key] {
			changed = append(changed, key)
		}
	}
	sort.Strings(changed)
	return changed
}
