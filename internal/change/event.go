// Package change defines provider change-event intake contracts. It performs
// no network access and grants no consumer authority.
package change

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/Clyra-AI/lumyn/internal/installation"
)

type ProviderEvent struct {
	EventID                       string    `json:"event_id"`
	EventVersion                  string    `json:"event_version"`
	EventDigest                   string    `json:"event_digest"`
	Issuer                        string    `json:"issuer"`
	APIOrSDK                      string    `json:"api_or_sdk"`
	Audience                      []string  `json:"audience"`
	Deadline                      time.Time `json:"deadline"`
	Severity                      string    `json:"severity"`
	Sequence                      uint64    `json:"sequence"`
	IssuedAt                      time.Time `json:"issued_at"`
	ExpiresAt                     time.Time `json:"expires_at"`
	TransportMode                 string    `json:"transport_mode"`
	TransportOrigin               string    `json:"transport_origin"`
	ManifestURL                   string    `json:"manifest_url"`
	Authenticated                 bool      `json:"authenticated"`
	SignatureValid                bool      `json:"signature_valid"`
	SignatureProvenance           string    `json:"signature_provenance"`
	ContractDigest                string    `json:"contract_digest"`
	RetrievedContractDigest       string    `json:"retrieved_contract_digest"`
	ContractDeliveryMode          string    `json:"contract_delivery_mode"`
	ContractURL                   string    `json:"contract_url,omitempty"`
	ProviderChannelDeliveryProven bool      `json:"provider_channel_delivery_proven"`
	Lifecycle                     string    `json:"lifecycle"`
	Executable                    bool      `json:"executable"`
}

type EventContext struct {
	PinnedOrigin string
	Audience     string
	LastSequence uint64
	SeenEvents   map[string]string
	Now          time.Time
	MaximumAge   time.Duration
}

// ProviderEventNetworkBindings is the decoded network-bearing subset of a
// persisted provider-change-event artifact. JSON Schema validates its shape;
// this validator compares the URLs that schema cannot relate.
type ProviderEventNetworkBindings struct {
	Transport struct {
		Mode         string `json:"mode"`
		ManifestURL  string `json:"manifest_url"`
		PinnedOrigin string `json:"pinned_origin"`
	} `json:"transport"`
	ContractDelivery struct {
		Mode        string `json:"mode"`
		ContractURL string `json:"contract_url"`
	} `json:"contract_delivery"`
}

func ValidateProviderEventNetworkBindings(value ProviderEventNetworkBindings) error {
	switch value.Transport.Mode {
	case "pinned_provider_https":
		if err := installation.ValidateURLAtPinnedOrigin(value.Transport.PinnedOrigin, value.Transport.ManifestURL); err != nil {
			return fmt.Errorf("provider event manifest: %w", err)
		}
		if value.ContractDelivery.Mode == "exact_provider_https_url" {
			if err := installation.ValidateURLAtPinnedOrigin(value.Transport.PinnedOrigin, value.ContractDelivery.ContractURL); err != nil {
				return fmt.Errorf("Provider Change Contract URL: %w", err)
			}
		}
	case "attended_import":
		// Recovery input carries no authenticated provider-origin proof.
	default:
		return fmt.Errorf("unsupported provider transport mode %q", value.Transport.Mode)
	}
	return nil
}

func ValidateProviderEvent(event ProviderEvent, context EventContext) error {
	for name, value := range map[string]string{
		"event id": event.EventID, "event version": event.EventVersion,
		"event digest": event.EventDigest, "issuer": event.Issuer,
		"API or SDK": event.APIOrSDK, "signature provenance": event.SignatureProvenance,
		"contract digest": event.ContractDigest, "retrieved contract digest": event.RetrievedContractDigest,
	} {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s is required", name)
		}
	}
	if event.Executable {
		return errors.New("provider event and contract must be non-executable")
	}
	if len(event.Audience) == 0 || event.Deadline.IsZero() || strings.TrimSpace(event.Severity) == "" {
		return errors.New("provider event audience, deadline, and severity are required")
	}
	if context.Now.IsZero() {
		return errors.New("validation time is required")
	}
	if !event.Deadline.After(context.Now) {
		return errors.New("provider event deadline has passed")
	}
	if event.Sequence <= context.LastSequence {
		return errors.New("provider event is replayed or stale")
	}
	if priorDigest, seen := context.SeenEvents[event.EventID]; seen {
		if priorDigest == event.EventDigest {
			return errors.New("provider event is a duplicate")
		}
		return errors.New("provider event identity conflicts with prior bytes")
	}
	if event.IssuedAt.IsZero() || event.IssuedAt.After(context.Now) {
		return errors.New("provider event issue time is invalid or in the future")
	}
	if event.ExpiresAt.IsZero() || !event.ExpiresAt.After(context.Now) || !event.ExpiresAt.After(event.IssuedAt) {
		return errors.New("provider event is expired or has invalid freshness")
	}
	if context.MaximumAge > 0 && context.Now.Sub(event.IssuedAt) > context.MaximumAge {
		return errors.New("provider event is stale")
	}
	if !slices.Contains(event.Audience, context.Audience) {
		return errors.New("provider event audience does not match installation")
	}
	switch event.TransportMode {
	case "pinned_provider_https":
		if event.TransportOrigin != context.PinnedOrigin {
			return errors.New("provider event origin does not match pinned installation origin")
		}
	case "attended_import":
		if event.ProviderChannelDeliveryProven {
			return errors.New("attended import cannot count as provider-channel delivery")
		}
	default:
		return fmt.Errorf("unsupported provider transport mode %q", event.TransportMode)
	}
	bindings := ProviderEventNetworkBindings{}
	bindings.Transport.Mode = event.TransportMode
	bindings.Transport.ManifestURL = event.ManifestURL
	bindings.Transport.PinnedOrigin = context.PinnedOrigin
	bindings.ContractDelivery.Mode = event.ContractDeliveryMode
	bindings.ContractDelivery.ContractURL = event.ContractURL
	if err := ValidateProviderEventNetworkBindings(bindings); err != nil {
		return err
	}
	if !event.Authenticated || !event.SignatureValid {
		return errors.New("provider event authentication or signature is invalid")
	}
	if event.ContractDigest != event.RetrievedContractDigest {
		return errors.New("retrieved Provider Change Contract bytes do not match event digest")
	}
	if event.Lifecycle != "active" {
		return fmt.Errorf("provider event lifecycle %q is not actionable", event.Lifecycle)
	}
	switch event.ContractDeliveryMode {
	case "embedded":
		if event.ContractURL != "" {
			return errors.New("embedded Provider Change Contract cannot carry a fetch URL")
		}
		if event.TransportMode == "pinned_provider_https" && !event.ProviderChannelDeliveryProven {
			return errors.New("provider channel delivery is not proven")
		}
	case "exact_provider_https_url":
		if !event.ProviderChannelDeliveryProven {
			return errors.New("provider channel delivery is not proven")
		}
	default:
		return fmt.Errorf("unsupported Provider Change Contract delivery mode %q", event.ContractDeliveryMode)
	}
	return nil
}
