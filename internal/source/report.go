package source

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func writeReport(root string, report Report) error {
	path := resolveProjectPath(root, report.ReportPath)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create report directory %s: %w", cleanSlashPath(filepath.Dir(path)), err)
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("encode source report %s: %w", report.ReportPath, err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write source report %s: %w", report.ReportPath, err)
	}
	return nil
}

func statusFromFindings(findings []Finding) string {
	status := "pass"
	for _, finding := range findings {
		if finding.Severity == "error" {
			return "fail"
		}
		if finding.Severity == "warning" {
			status = "warning"
		}
	}
	return status
}

func HasErrorFindings(findings []Finding) bool {
	for _, finding := range findings {
		if finding.Severity == "error" {
			return true
		}
	}
	return false
}

func FirstErrorFinding(findings []Finding) (Finding, bool) {
	for _, finding := range findings {
		if finding.Severity == "error" {
			return finding, true
		}
	}
	return Finding{}, false
}

func FirstFinding(findings []Finding) (Finding, bool) {
	if len(findings) == 0 {
		return Finding{}, false
	}
	return findings[0], true
}

func sourceFindingHasKind(findings []Finding, kind string) bool {
	for _, finding := range findings {
		if finding.Kind == kind {
			return true
		}
	}
	return false
}

func ContainsProofGap(findings []Finding) bool {
	return sourceFindingHasKind(findings, "proof_gap")
}

func ContainsSecurityRelevantFinding(findings []Finding) bool {
	for _, finding := range findings {
		switch finding.Kind {
		case "auth_confusion", "data_exposure_risk", "forbidden_endpoint_call", "scope_escalation", "unexpected_write_action":
			return true
		}
	}
	return false
}
