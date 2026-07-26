package installation

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestValidateProviderChannelRejectsCrossOriginManifest(t *testing.T) {
	channel := ProviderChannel{
		PinnedOrigin: "https://updates.payments.example",
		ManifestURL:  "https://attacker.invalid/campaign/event.json",
	}
	if err := ValidateProviderChannel(channel); err == nil {
		t.Fatal("cross-origin provider manifest must fail closed")
	}
}

func TestValidateProviderChannelAcceptsSameOriginManifest(t *testing.T) {
	channel := ProviderChannel{
		PinnedOrigin: "https://updates.payments.example",
		ManifestURL:  "https://updates.payments.example/campaign/event.json",
	}
	if err := ValidateProviderChannel(channel); err != nil {
		t.Fatalf("same-origin provider manifest rejected: %v", err)
	}
}

func TestValidateProviderChannelRejectsPersistedCrossOriginFixture(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve installation test path")
	}
	path := filepath.Join(filepath.Dir(filename), "..", "..", "tests", "fixtures", "contracts", "consumer-installation", "semantic-invalid-cross-origin.json")
	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read semantic fixture: %v", err)
	}
	var channel ProviderChannel
	if err := json.Unmarshal(payload, &channel); err != nil {
		t.Fatalf("decode semantic fixture: %v", err)
	}
	if err := ValidateProviderChannel(channel); err == nil {
		t.Fatal("cross-origin persisted installation channel must fail closed")
	}
}
