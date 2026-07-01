package markdownlinks

import (
	"fmt"
	"testing"
)

func TestTargetsParsesLocalMarkdownLinks(t *testing.T) {
	tests := []struct {
		name string
		line string
		want []string
	}{
		{
			name: "space target in brackets",
			line: `Read [auth](<auth guide.md>).`,
			want: []string{"<auth guide.md>"},
		},
		{
			name: "balanced parentheses",
			line: `Read [setup](setup(v1).md).`,
			want: []string{"setup(v1).md"},
		},
		{
			name: "multiple links",
			line: `[one](one.md) and [two](two.md#install)`,
			want: []string{"one.md", "two.md#install"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Targets(test.line)
			if fmt.Sprint(got) != fmt.Sprint(test.want) {
				t.Fatalf("Targets() = %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestTargetNormalizationHelpers(t *testing.T) {
	if got := CleanTarget(`<auth guide.md>`); got != "auth guide.md" {
		t.Fatalf("CleanTarget() = %q", got)
	}
	if got := LocalPath("missing.md?api_key=secret#install"); got != "missing.md" {
		t.Fatalf("LocalPath() = %q", got)
	}
	if got := FindingTarget("missing.md?api_key=secret#install"); got != "missing.md [query-or-fragment-redacted]" {
		t.Fatalf("FindingTarget() = %q", got)
	}
}

func TestIsFenceDelimiter(t *testing.T) {
	for _, line := range []string{"```", "```md", "~~~", "~~~text"} {
		if !IsFenceDelimiter(line) {
			t.Fatalf("IsFenceDelimiter(%q) = false", line)
		}
	}
}
