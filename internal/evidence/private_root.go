// Package evidence owns consumer-private artifact storage boundaries.
package evidence

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ValidatePrivateRoot resolves both paths through symlinks and rejects a
// consumer-private root that is the checkout, lives below it, is not a private
// directory, or resolves to a filesystem root.
func ValidatePrivateRoot(checkout, privateRoot string) (string, error) {
	resolvedCheckout, err := resolveDirectory(checkout)
	if err != nil {
		return "", fmt.Errorf("resolve checkout: %w", err)
	}
	resolvedPrivate, err := resolveDirectory(privateRoot)
	if err != nil {
		return "", fmt.Errorf("resolve private root: %w", err)
	}
	if resolvedPrivate == filepath.VolumeName(resolvedPrivate)+string(filepath.Separator) {
		return "", fmt.Errorf("private root cannot be a filesystem root")
	}
	relative, err := filepath.Rel(resolvedCheckout, resolvedPrivate)
	if err != nil {
		return "", fmt.Errorf("compare checkout and private root: %w", err)
	}
	if relative == "." || (relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))) {
		return "", fmt.Errorf("private root must resolve outside the checkout")
	}
	info, err := os.Stat(resolvedPrivate)
	if err != nil {
		return "", fmt.Errorf("stat private root: %w", err)
	}
	if info.Mode().Perm()&0o077 != 0 {
		return "", fmt.Errorf("private root must not grant group or other permissions")
	}
	return resolvedPrivate, nil
}

func resolveDirectory(path string) (string, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%s is not a directory", resolved)
	}
	return filepath.Clean(resolved), nil
}
