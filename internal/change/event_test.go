package change

import (
	"crypto/ed25519"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestValidateProviderEventNetworkBindingsRejectsPersistedCrossOriginFixture(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve provider event test path")
	}
	path := filepath.Join(filepath.Dir(filename), "..", "..", "tests", "fixtures", "contracts", "provider-change-event", "semantic-invalid-cross-origin.json")
	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read semantic fixture: %v", err)
	}
	var bindings ProviderEventNetworkBindings
	if err := json.Unmarshal(payload, &bindings); err != nil {
		t.Fatalf("decode semantic fixture: %v", err)
	}
	if err := ValidateProviderEventNetworkBindings(bindings); err == nil {
		t.Fatal("cross-origin persisted provider event URLs must fail closed")
	}
}

func TestValidateProviderEventAcceptsFreshAuthenticatedPinnedEvent(t *testing.T) {
	now := time.Date(2026, 7, 26, 14, 0, 0, 0, time.UTC)
	event := validProviderEvent(now)
	if err := ValidateProviderEvent(event, EventContext{
		PinnedOrigin:      "https://updates.example.com",
		Audience:          "sdk:stripe-go:v80",
		ExpectedKeyID:     "campaign_key:stripe-2026",
		ExpectedPublicKey: make(ed25519.PublicKey, ed25519.PublicKeySize),
		LastSequence:      41,
		SeenEvents:        map[string]string{},
		Now:               now,
	}); err != nil {
		t.Fatalf("valid event rejected: %v", err)
	}
}

func TestValidateProviderEventRejectsReplayConflictLifecycleAndRecoveryProof(t *testing.T) {
	now := time.Date(2026, 7, 26, 14, 0, 0, 0, time.UTC)
	tests := map[string]func(*ProviderEvent, *EventContext){
		"duplicate": func(event *ProviderEvent, context *EventContext) {
			context.SeenEvents[event.EventID] = event.EventDigest
		},
		"conflict": func(event *ProviderEvent, context *EventContext) {
			context.SeenEvents[event.EventID] = "sha256:other"
		},
		"replay":         func(event *ProviderEvent, context *EventContext) { event.Sequence = context.LastSequence },
		"expired":        func(event *ProviderEvent, _ *EventContext) { event.ExpiresAt = now },
		"past deadline":  func(event *ProviderEvent, _ *EventContext) { event.Deadline = now },
		"future issued":  func(event *ProviderEvent, _ *EventContext) { event.IssuedAt = now.Add(time.Minute) },
		"wrong audience": func(event *ProviderEvent, _ *EventContext) { event.Audience = []string{"sdk:other"} },
		"wrong origin": func(event *ProviderEvent, _ *EventContext) {
			event.TransportOrigin = "https://mirror.invalid"
		},
		"cross-origin manifest": func(event *ProviderEvent, _ *EventContext) {
			event.ManifestURL = "https://mirror.invalid/event.json"
		},
		"cross-origin contract": func(event *ProviderEvent, _ *EventContext) {
			event.ContractURL = "https://mirror.invalid/migration-pack.json"
		},
		"unauthenticated": func(event *ProviderEvent, _ *EventContext) { event.Authenticated = false },
		"wrong campaign key": func(event *ProviderEvent, _ *EventContext) {
			event.SignatureProvenance = "campaign_key:attacker"
		},
		"digest mismatch": func(event *ProviderEvent, _ *EventContext) { event.RetrievedContractDigest = "sha256:other" },
		"withdrawn":       func(event *ProviderEvent, _ *EventContext) { event.Lifecycle = "withdrawn" },
		"executable":      func(event *ProviderEvent, _ *EventContext) { event.Executable = true },
		"attended import as channel proof": func(event *ProviderEvent, _ *EventContext) {
			event.TransportMode = "attended_import"
			event.ProviderChannelDeliveryProven = true
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			event := validProviderEvent(now)
			context := EventContext{
				PinnedOrigin: "https://updates.example.com",
				Audience:     "sdk:stripe-go:v80", ExpectedKeyID: "campaign_key:stripe-2026",
				ExpectedPublicKey: make(ed25519.PublicKey, ed25519.PublicKeySize),
				LastSequence:      41,
				SeenEvents:        map[string]string{}, Now: now,
			}
			mutate(&event, &context)
			if err := ValidateProviderEvent(event, context); err == nil {
				t.Fatal("expected provider event rejection")
			}
		})
	}
}

func validProviderEvent(now time.Time) ProviderEvent {
	return ProviderEvent{
		EventID: "evt-stripe-2026-07", EventVersion: "1.0", EventDigest: "sha256:event",
		Issuer: "stripe", APIOrSDK: "stripe-go", Audience: []string{"sdk:stripe-go:v80"}, Sequence: 42,
		Deadline: now.Add(30 * 24 * time.Hour), Severity: "breaking",
		IssuedAt: now.Add(-time.Minute), ExpiresAt: now.Add(time.Hour),
		TransportMode:   "pinned_provider_https",
		TransportOrigin: "https://updates.example.com",
		ManifestURL:     "https://updates.example.com/stripe/v1.json",
		Authenticated:   true, SignatureValid: true, SignatureProvenance: "campaign_key:stripe-2026",
		ContractDigest: "sha256:pack", RetrievedContractDigest: "sha256:pack",
		ContractDeliveryMode:          "exact_provider_https_url",
		ContractURL:                   "https://updates.example.com/stripe/migration-pack.json",
		ProviderChannelDeliveryProven: true,
		Lifecycle:                     "active", Executable: false,
	}
}
