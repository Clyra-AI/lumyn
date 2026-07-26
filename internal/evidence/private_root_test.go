package evidence

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidatePrivateRootAcceptsSiblingDirectory(t *testing.T) {
	root := t.TempDir()
	checkout := filepath.Join(root, "checkout")
	privateRoot := filepath.Join(root, "private")
	for _, path := range []string{checkout, privateRoot} {
		if err := os.Mkdir(path, 0o700); err != nil {
			t.Fatalf("mkdir %s: %v", path, err)
		}
	}
	resolved, err := ValidatePrivateRoot(checkout, privateRoot)
	if err != nil {
		t.Fatalf("sibling private root rejected: %v", err)
	}
	expected, err := filepath.EvalSymlinks(privateRoot)
	if err != nil {
		t.Fatal(err)
	}
	if resolved != expected {
		t.Fatalf("resolved private root = %q, want %q", resolved, expected)
	}
}

func TestValidatePrivateRootRejectsCheckoutDescendantAndSymlink(t *testing.T) {
	root := t.TempDir()
	checkout := filepath.Join(root, "checkout")
	inside := filepath.Join(checkout, "private")
	outsideLink := filepath.Join(root, "private-link")
	if err := os.MkdirAll(inside, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(inside, outsideLink); err != nil {
		t.Fatal(err)
	}
	for name, privateRoot := range map[string]string{
		"descendant":            inside,
		"symlink into checkout": outsideLink,
		"checkout itself":       checkout,
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := ValidatePrivateRoot(checkout, privateRoot); err == nil {
				t.Fatal("expected unsafe private root to fail")
			}
		})
	}
}
