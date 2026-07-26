// Package installation validates semantic relationships in persisted Consumer
// Installation contracts that JSON Schema cannot compare.
package installation

import (
	"fmt"
	"net/url"
	"strings"
)

// ProviderChannel is the decoded provider_channel portion of a Consumer
// Installation. A runtime must validate it after schema validation and before
// any network request.
type ProviderChannel struct {
	PinnedOrigin string `json:"pinned_origin"`
	ManifestURL  string `json:"manifest_url"`
}

func ValidateProviderChannel(channel ProviderChannel) error {
	if err := ValidateURLAtPinnedOrigin(channel.PinnedOrigin, channel.ManifestURL); err != nil {
		return fmt.Errorf("provider manifest URL: %w", err)
	}
	return nil
}

// ValidateURLAtPinnedOrigin canonicalizes both values and rejects a fetch URL
// outside the exact installed HTTPS origin. It deliberately grants no network
// authority and performs no request.
func ValidateURLAtPinnedOrigin(pinnedOrigin, fetchURL string) error {
	pinned, err := canonicalHTTPSOrigin(pinnedOrigin, true)
	if err != nil {
		return fmt.Errorf("pinned origin: %w", err)
	}
	fetch, err := canonicalHTTPSOrigin(fetchURL, false)
	if err != nil {
		return fmt.Errorf("fetch URL: %w", err)
	}
	if fetch != pinned {
		return fmt.Errorf("origin %q does not match pinned origin %q", fetch, pinned)
	}
	return nil
}

func canonicalHTTPSOrigin(value string, originOnly bool) (string, error) {
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil || parsed.Fragment != "" {
		return "", fmt.Errorf("must be an absolute HTTPS URL without user info or fragment")
	}
	if originOnly && ((parsed.Path != "" && parsed.Path != "/") || parsed.RawQuery != "") {
		return "", fmt.Errorf("must contain only scheme and authority")
	}
	host := strings.ToLower(parsed.Hostname())
	if host == "" {
		return "", fmt.Errorf("host is required")
	}
	port := parsed.Port()
	if port == "443" {
		port = ""
	}
	if port != "" {
		host += ":" + port
	}
	return "https://" + host, nil
}
