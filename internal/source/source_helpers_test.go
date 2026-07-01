package source

func containsString(haystack, needle string) bool {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}

func hasFindingKind(findings []Finding, kind string) bool {
	for _, finding := range findings {
		if finding.Kind == kind {
			return true
		}
	}
	return false
}

func hasFindingFixTarget(findings []Finding, fixTarget string) bool {
	for _, finding := range findings {
		if finding.FixTarget == fixTarget {
			return true
		}
	}
	return false
}
