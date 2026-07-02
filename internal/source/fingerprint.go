package source

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func hashPath(root, path string) string {
	resolved := resolveProjectPath(root, path)
	info, err := os.Stat(resolved)
	if err != nil {
		return ""
	}
	if !info.IsDir() {
		return hashFile(resolved)
	}
	fileHashes := []string{}
	_ = filepath.WalkDir(resolved, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		hash := hashFile(path)
		if hash == "" {
			return nil
		}
		fileHashes = append(fileHashes, relativePath(root, path)+"="+hash)
		return nil
	})
	sort.Strings(fileHashes)
	digest := sha256.Sum256([]byte(strings.Join(fileHashes, "\n")))
	return "sha256:" + hex.EncodeToString(digest[:])
}

func hashDocsPath(root, path string) string {
	resolved := resolveProjectPath(root, path)
	info, err := os.Stat(resolved)
	if err != nil {
		return ""
	}
	if !info.IsDir() {
		return hashFile(resolved)
	}
	fileHashes := []string{}
	_ = filepath.WalkDir(resolved, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if shouldSkipGeneratedSourceDir(d.Name()) && path != resolved {
				return filepath.SkipDir
			}
			return nil
		}
		if !isDocsFile(path) {
			return nil
		}
		hash := hashFile(path)
		if hash == "" {
			return nil
		}
		fileHashes = append(fileHashes, relativePath(root, path)+"="+hash)
		return nil
	})
	sort.Strings(fileHashes)
	digest := sha256.Sum256([]byte(strings.Join(fileHashes, "\n")))
	return "sha256:" + hex.EncodeToString(digest[:])
}

func shouldSkipGeneratedSourceDir(name string) bool {
	switch name {
	case ".git", ".factory", ".factoryd", "runs":
		return true
	default:
		return false
	}
}

func hashFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	digest := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(digest[:])
}

func surfaceFingerprint(refs []SourceRef) string {
	parts := make([]string, 0, len(refs))
	for _, ref := range refs {
		parts = append(parts, ref.Kind+":"+ref.ID+":"+ref.Path+":"+ref.Hash)
	}
	sort.Strings(parts)
	digest := sha256.Sum256([]byte(strings.Join(parts, "\n")))
	return "sha256:" + hex.EncodeToString(digest[:])
}
