package source

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Clyra-AI/lumyn/internal/source/markdownlinks"
)

func checkDocs(root string, entry SourceEntry) []Finding {
	docsPath := cleanSlashPath(entry.Path)
	absPath := resolveProjectPath(root, entry.Path)
	info, err := os.Stat(absPath)
	if err != nil {
		return []Finding{{
			Kind:     "context_missing",
			Severity: "warning",
			Message:  fmt.Sprintf("docs source %s could not be read", docsPath),
			Reference: Reference{
				Path:   docsPath,
				Object: "sources.docs." + entry.ID,
			},
			FixTarget:         "sources.docs.path",
			WorkflowRelevance: "Agents need local docs context for auth, retries, pagination, and proof guidance.",
		}}
	}

	findings := []Finding{}
	readableFiles := 0
	var combined strings.Builder
	visitFile := func(path string, d fs.DirEntry) {
		if d.IsDir() {
			return
		}
		if !isDocsFile(path) {
			return
		}
		relPath := relativePath(root, path)
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			findings = append(findings, Finding{
				Kind:     "context_missing",
				Severity: "warning",
				Message:  fmt.Sprintf("docs file %s could not be read", relPath),
				Reference: Reference{
					Path:   relPath,
					Object: "docs_file",
				},
				FixTarget:         "docs_file",
				WorkflowRelevance: "Unreadable docs can hide auth, retry, or validation instructions from agent probes.",
			})
			return
		}
		readableFiles++
		combined.Write(bytes.ToLower(data))
		combined.WriteByte('\n')
		findings = append(findings, brokenLocalReferenceFindings(root, path, data)...)
	}

	if info.IsDir() {
		walkErr := filepath.WalkDir(absPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				findings = append(findings, Finding{
					Kind:     "context_missing",
					Severity: "warning",
					Message:  fmt.Sprintf("docs path %s could not be inspected: %v", relativePath(root, path), err),
					Reference: Reference{
						Path:   relativePath(root, path),
						Object: "docs_path",
					},
					FixTarget:         "docs_path",
					WorkflowRelevance: "Unreadable docs can hide workflow constraints from source checks.",
				})
				return nil
			}
			if d.IsDir() && path != absPath && shouldSkipGeneratedSourceDir(d.Name()) {
				return filepath.SkipDir
			}
			visitFile(path, d)
			return nil
		})
		if walkErr != nil {
			findings = append(findings, Finding{
				Kind:     "context_missing",
				Severity: "warning",
				Message:  fmt.Sprintf("docs source %s could not be walked: %v", docsPath, walkErr),
				Reference: Reference{
					Path:   docsPath,
					Object: "sources.docs." + entry.ID,
				},
				FixTarget:         "sources.docs.path",
				WorkflowRelevance: "Agents need readable local docs for workflow-relevant source context.",
			})
		}
	} else {
		visitFile(absPath, fileEntry{name: filepath.Base(absPath)})
	}

	if readableFiles == 0 {
		findings = append(findings, Finding{
			Kind:     "context_missing",
			Severity: "warning",
			Message:  fmt.Sprintf("docs source %s contains no readable docs files", docsPath),
			Reference: Reference{
				Path:   docsPath,
				Object: "sources.docs." + entry.ID,
			},
			FixTarget:         "docs_content",
			WorkflowRelevance: "Agent probes need readable docs for setup, auth, retry, and validation instructions.",
		})
		return findings
	}

	if missingOperationalGuidance(combined.String()) {
		findings = append(findings, Finding{
			Kind:     "docs_api_ambiguity",
			Severity: "warning",
			Message:  "local docs do not mention retry, rate-limit, pagination, or idempotency guidance",
			Reference: Reference{
				Path:   docsPath,
				Object: "docs_guidance",
			},
			FixTarget:         "operational_docs",
			WorkflowRelevance: "Retries, rate limits, pagination, and idempotency affect agent workflow stability and write safety.",
		})
	}
	return findings
}

func brokenLocalReferenceFindings(root, docPath string, data []byte) []Finding {
	findings := []Finding{}
	lines := bytes.Split(data, []byte("\n"))
	inFence := false
	for index, line := range lines {
		trimmedLine := strings.TrimSpace(string(line))
		if markdownlinks.IsFenceDelimiter(trimmedLine) {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		for _, rawTarget := range markdownlinks.Targets(string(line)) {
			target := markdownlinks.CleanTarget(rawTarget)
			if target == "" || isExternalReference(target) {
				continue
			}
			targetPath := markdownlinks.LocalPath(target)
			if targetPath == "" {
				continue
			}
			if unescaped, err := url.PathUnescape(targetPath); err == nil {
				targetPath = unescaped
			}
			resolved := resolveMarkdownLinkTarget(root, docPath, targetPath)
			if _, err := os.Stat(resolved); err == nil {
				continue
			}
			findings = append(findings, Finding{
				Kind:     "context_missing",
				Severity: "warning",
				Message:  fmt.Sprintf("docs file %s links to missing local reference %s", relativePath(root, docPath), markdownlinks.FindingTarget(target)),
				Reference: Reference{
					Path: relativePath(root, docPath),
					Line: index + 1,
				},
				FixTarget:         "docs_local_reference",
				WorkflowRelevance: "Broken local docs links can hide workflow setup, auth, or validation instructions from agents.",
			})
		}
	}
	return findings
}

func missingOperationalGuidance(lowerDocs string) bool {
	hasRetry := strings.Contains(lowerDocs, "retry")
	hasRateLimit := strings.Contains(lowerDocs, "rate limit") || strings.Contains(lowerDocs, "rate-limit") || strings.Contains(lowerDocs, "429")
	hasPagination := strings.Contains(lowerDocs, "pagination") || strings.Contains(lowerDocs, "page ")
	hasIdempotency := strings.Contains(lowerDocs, "idempotency") || strings.Contains(lowerDocs, "idempotent")
	return !(hasRetry && hasRateLimit && hasPagination && hasIdempotency)
}

func isExternalReference(target string) bool {
	lower := strings.ToLower(target)
	return strings.HasPrefix(lower, "http://") ||
		strings.HasPrefix(lower, "https://") ||
		strings.HasPrefix(lower, "mailto:") ||
		strings.HasPrefix(lower, "tel:") ||
		strings.HasPrefix(target, "#")
}

func isDocsFile(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".md", ".mdx", ".txt", ".rst":
		return true
	default:
		return false
	}
}

func resolveMarkdownLinkTarget(root, docPath, targetPath string) string {
	if strings.HasPrefix(targetPath, "/") {
		return filepath.Clean(filepath.Join(root, filepath.FromSlash(strings.TrimPrefix(targetPath, "/"))))
	}
	return filepath.Clean(filepath.Join(filepath.Dir(docPath), filepath.FromSlash(targetPath)))
}

type fileEntry struct {
	name string
}

func (f fileEntry) Name() string               { return f.name }
func (f fileEntry) IsDir() bool                { return false }
func (f fileEntry) Type() fs.FileMode          { return 0 }
func (f fileEntry) Info() (fs.FileInfo, error) { return nil, errors.New("file info unavailable") }
