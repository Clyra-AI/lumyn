// Package agent validates Agent Runner qualification contracts without
// launching a runner or model.
package agent

import (
	"fmt"
	"strings"
)

type ConformanceTestCase struct {
	CaseID         string `json:"case_id"`
	Status         string `json:"status"`
	EvidenceDigest string `json:"evidence_digest"`
}

type ConformanceSummary struct {
	Passed  int `json:"passed"`
	Failed  int `json:"failed"`
	Blocked int `json:"blocked"`
}

type LiveCanary struct {
	Approved       bool   `json:"approved"`
	Status         string `json:"status"`
	EvidenceDigest string `json:"evidence_digest"`
}

type ConformanceResult struct {
	QualificationStatus string                `json:"qualification_status"`
	LaunchEligible      bool                  `json:"launch_eligible"`
	TestCases           []ConformanceTestCase `json:"test_cases"`
	Summary             ConformanceSummary    `json:"summary"`
	LiveCanary          LiveCanary            `json:"live_canary"`
}

// ValidateConformanceResult derives the summary from the exact test-case set.
// A declared summary cannot make a failed or blocked case launch-eligible.
func ValidateConformanceResult(result ConformanceResult) error {
	if len(result.TestCases) == 0 {
		return fmt.Errorf("at least one conformance test case is required")
	}
	seen := make(map[string]struct{}, len(result.TestCases))
	derived := ConformanceSummary{}
	for _, testCase := range result.TestCases {
		if strings.TrimSpace(testCase.CaseID) == "" || strings.TrimSpace(testCase.EvidenceDigest) == "" {
			return fmt.Errorf("conformance case identity and evidence digest are required")
		}
		if _, duplicate := seen[testCase.CaseID]; duplicate {
			return fmt.Errorf("duplicate conformance case %q", testCase.CaseID)
		}
		seen[testCase.CaseID] = struct{}{}
		switch testCase.Status {
		case "passed":
			derived.Passed++
		case "failed":
			derived.Failed++
		case "blocked":
			derived.Blocked++
		default:
			return fmt.Errorf("unknown conformance case status %q", testCase.Status)
		}
	}
	if result.Summary != derived {
		return fmt.Errorf("conformance summary does not match exact test-case results")
	}
	switch result.QualificationStatus {
	case "passed":
		if derived.Failed != 0 || derived.Blocked != 0 || !result.LaunchEligible {
			return fmt.Errorf("passing qualification requires an all-passed test set and launch eligibility")
		}
		if !result.LiveCanary.Approved || result.LiveCanary.Status != "passed" || strings.TrimSpace(result.LiveCanary.EvidenceDigest) == "" {
			return fmt.Errorf("passing qualification requires approved passing live-canary evidence")
		}
	case "failed", "deferred", "stale":
		if result.LaunchEligible {
			return fmt.Errorf("non-passing qualification cannot be launch-eligible")
		}
	default:
		return fmt.Errorf("unknown qualification status %q", result.QualificationStatus)
	}
	return nil
}
