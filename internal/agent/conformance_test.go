package agent

import "testing"

func TestValidateConformanceResultRejectsFalseGreenSummary(t *testing.T) {
	result := ConformanceResult{
		QualificationStatus: "passed",
		LaunchEligible:      true,
		TestCases: []ConformanceTestCase{
			{CaseID: "isolation", Status: "passed", EvidenceDigest: "sha256:one"},
			{CaseID: "auth", Status: "failed", EvidenceDigest: "sha256:two"},
		},
		Summary:    ConformanceSummary{Passed: 2, Failed: 0, Blocked: 0},
		LiveCanary: LiveCanary{Approved: true, Status: "passed", EvidenceDigest: "sha256:canary"},
	}
	if err := ValidateConformanceResult(result); err == nil {
		t.Fatal("failed test hidden by a passing summary must fail closed")
	}
}

func TestValidateConformanceResultAcceptsDerivedPassingResult(t *testing.T) {
	result := ConformanceResult{
		QualificationStatus: "passed",
		LaunchEligible:      true,
		TestCases: []ConformanceTestCase{
			{CaseID: "isolation", Status: "passed", EvidenceDigest: "sha256:one"},
			{CaseID: "auth", Status: "passed", EvidenceDigest: "sha256:two"},
		},
		Summary:    ConformanceSummary{Passed: 2, Failed: 0, Blocked: 0},
		LiveCanary: LiveCanary{Approved: true, Status: "passed", EvidenceDigest: "sha256:canary"},
	}
	if err := ValidateConformanceResult(result); err != nil {
		t.Fatalf("derived passing conformance rejected: %v", err)
	}
}
