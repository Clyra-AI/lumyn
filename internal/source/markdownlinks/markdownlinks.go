package markdownlinks

import "strings"

func IsFenceDelimiter(line string) bool {
	return strings.HasPrefix(line, "```") || strings.HasPrefix(line, "~~~")
}

func Targets(line string) []string {
	targets := []string{}
	offset := 0
	for offset < len(line) {
		start := strings.Index(line[offset:], "](")
		if start < 0 {
			break
		}
		targetStart := offset + start + 2
		targetEnd := targetEnd(line, targetStart)
		if targetEnd < 0 {
			offset = targetStart
			continue
		}
		targets = append(targets, line[targetStart:targetEnd])
		offset = targetEnd + 1
	}
	return targets
}

func targetEnd(line string, start int) int {
	depth := 0
	escaped := false
	for index := start; index < len(line); index++ {
		char := line[index]
		switch {
		case escaped:
			escaped = false
		case char == '\\':
			escaped = true
		case char == '(':
			depth++
		case char == ')':
			if depth == 0 {
				return index
			}
			depth--
		}
	}
	return -1
}

func LocalPath(target string) string {
	targetPath := strings.SplitN(target, "#", 2)[0]
	targetPath = strings.SplitN(targetPath, "?", 2)[0]
	return targetPath
}

func FindingTarget(target string) string {
	targetPath := LocalPath(target)
	if targetPath == "" {
		return "<local-fragment>"
	}
	if targetPath != target {
		return targetPath + " [query-or-fragment-redacted]"
	}
	return targetPath
}

func CleanTarget(target string) string {
	target = strings.TrimSpace(target)
	if strings.HasPrefix(target, "<") {
		if end := strings.Index(target, ">"); end > 0 {
			return strings.Trim(strings.TrimSpace(target[1:end]), `"'`)
		}
	}
	if space := strings.IndexByte(target, ' '); space >= 0 {
		target = target[:space]
	}
	return strings.Trim(strings.Trim(target, "<>"), `"'`)
}
