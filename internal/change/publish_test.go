package change

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/Clyra-AI/lumyn/internal/pack"
	"github.com/Clyra-AI/lumyn/internal/source"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

func TestPublishKitVerifiesPinnedChannelAndLabelsAttendedRecovery(t *testing.T) {
	kit, publicKey, context := validPublishKit(t)
	result, err := VerifyPublishKit(kit, publicKey, context, IntakeObservation{
		Mode: "pinned_provider_https", ObservedManifestURL: "https://updates.payments.example/events/v5.json",
	})
	if err != nil {
		t.Fatalf("verify pinned publish kit: %v", err)
	}
	if !result.ChannelDeliveryQualified || !result.PreauthorizationEligible {
		t.Fatalf("pinned provider delivery lost qualification: %#v", result)
	}

	recovery, err := VerifyPublishKit(kit, publicKey, context, IntakeObservation{Mode: "attended_import"})
	if err != nil {
		t.Fatalf("verify attended recovery: %v", err)
	}
	if recovery.ChannelDeliveryQualified || recovery.PreauthorizationEligible {
		t.Fatalf("attended import must not prove channel delivery or preauthorization: %#v", recovery)
	}
}

func TestPublishKitRejectsTamperReplayAndWrongAudience(t *testing.T) {
	kit, publicKey, context := validPublishKit(t)
	var authorityWidened EventArtifact
	if err := json.Unmarshal(kit.EventBytes, &authorityWidened); err != nil {
		t.Fatalf("decode event: %v", err)
	}
	authorityWidened.GrantsConsumerAuthority = true
	if err := validateEventArtifact(authorityWidened); err == nil {
		t.Fatal("even provider-signed event content cannot grant consumer authority")
	}

	tampered := kit
	tampered.ContractBytes = append([]byte(nil), kit.ContractBytes...)
	tampered.ContractBytes[len(tampered.ContractBytes)-2] ^= 1
	if _, err := VerifyPublishKit(tampered, publicKey, context, IntakeObservation{
		Mode: "pinned_provider_https", ObservedManifestURL: "https://updates.payments.example/events/v5.json",
	}); err == nil {
		t.Fatal("tampered contract bytes must fail")
	}

	replay := context
	replay.LastSequence = 41
	if _, err := VerifyPublishKit(kit, publicKey, replay, IntakeObservation{
		Mode: "pinned_provider_https", ObservedManifestURL: "https://updates.payments.example/events/v5.json",
	}); err == nil {
		t.Fatal("replayed sequence must fail")
	}

	wrongAudience := context
	wrongAudience.Audience = "audience.other"
	if _, err := VerifyPublishKit(kit, publicKey, wrongAudience, IntakeObservation{
		Mode: "pinned_provider_https", ObservedManifestURL: "https://updates.payments.example/events/v5.json",
	}); err == nil {
		t.Fatal("wrong audience must fail")
	}

	tamperedEvent := kit
	tamperedEvent.EventBytes = append([]byte(nil), kit.EventBytes...)
	tamperedEvent.EventBytes[len(tamperedEvent.EventBytes)/2] ^= 1
	if _, err := VerifyPublishKit(tamperedEvent, publicKey, context, IntakeObservation{
		Mode: "pinned_provider_https", ObservedManifestURL: "https://updates.payments.example/events/v5.json",
	}); err == nil {
		t.Fatal("tampered event must fail")
	}
}

func validPublishKit(t *testing.T) (PublishKit, ed25519.PublicKey, EventContext) {
	t.Helper()
	now := time.Date(2026, 7, 28, 12, 0, 0, 0, time.UTC)
	target := "paymentIntents.confirm(id, { payment_method })"
	normalized, err := pack.Normalize(pack.NormalizeRequest{
		MigrationPackID: "migration_pack.payments.v5", ContractVersion: 1,
		APIProviderID: "api_provider.payments",
		APIOrSDK:      pack.APIOrSDK{Kind: "sdk", Name: "Payments Node SDK", Package: "@payments/sdk"},
		SourceVersion: "4.x", TargetVersion: "5.x",
		Audience: pack.Audience{AudienceID: "audience.node_v4", PackageVersionSelector: ">=4 <5"},
		Evidence: []pack.SourceEvidenceInput{
			changeTestEvidence("evidence.sdk_v4", source.VersionSource, "4.x", "https://updates.payments.example/sdk/v4/index.d.ts", "v4 signature", now),
			changeTestEvidence("evidence.sdk_v5", source.VersionTarget, "5.x", "https://updates.payments.example/sdk/v5/index.d.ts", "v5 signature", now),
		},
		Declarations: []pack.Declaration{{SemanticChange: pack.SemanticChange{
			ChangeID: "change.confirm_options", Kind: "signature_change",
			SourceSymbol: "paymentIntents.confirm(id, paymentMethod)", TargetSymbol: &target,
			Intent:        "Move the payment method into the options object.",
			Applicability: "Direct and statically resolved wrapper calls.",
			Mapping: pack.Mapping{Mode: "repository_specific_reasoning", DeclarativeSteps: []string{
				"Preserve the existing payment method expression in the payment_method field.",
			}},
			SourceEvidenceIDs: []string{"evidence.sdk_v4"}, TargetEvidenceIDs: []string{"evidence.sdk_v5"},
		}}},
		VerificationGuidance: []string{"Run the consumer-approved typecheck and tests."},
		NormalizedAt:         now, Lifecycle: pack.Lifecycle{State: "active"},
	})
	if err != nil {
		t.Fatalf("normalize contract: %v", err)
	}
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	confirmed, record, err := pack.Confirm(normalized.Contract, "provider_operator.release_team", now.Add(time.Minute), privateKey)
	if err != nil {
		t.Fatalf("confirm contract: %v", err)
	}
	if err := pack.VerifyConfirmation(confirmed, record, publicKey); err != nil {
		t.Fatalf("verify contract confirmation: %v", err)
	}
	kit, err := BuildPublishKit(confirmed, record, publicKey, PublishRequest{
		EventID: "event.payments.v5.0001", EventVersion: 1,
		ProviderOperatorID: "provider_operator.release_team", KeyID: "campaign_key.payments.v5",
		AudienceID: "audience.node_v4", Deadline: now.Add(90 * 24 * time.Hour),
		Severity: "breaking", Sequence: 41, IssuedAt: now.Add(-time.Minute),
		ExpiresAt: now.Add(time.Hour), PinnedOrigin: "https://updates.payments.example",
		ManifestURL: "https://updates.payments.example/events/v5.json",
		ContractURL: "https://updates.payments.example/contracts/v5.json",
	}, privateKey)
	if err != nil {
		t.Fatalf("build publish kit: %v", err)
	}
	validateGeneratedArtifact(t, "migration-pack.schema.json", kit.ContractBytes)
	validateGeneratedArtifact(t, "provider-change-event.schema.json", kit.EventBytes)
	return kit, publicKey, EventContext{
		PinnedOrigin: "https://updates.payments.example", Audience: "audience.node_v4",
		LastSequence: 40, SeenEvents: map[string]string{}, Now: now, MaximumAge: 24 * time.Hour,
	}
}

func validateGeneratedArtifact(t *testing.T, schemaName string, data []byte) {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve M3 test path")
	}
	schema, err := jsonschema.Compile(filepath.Join(filepath.Dir(filename), "..", "..", "schemas", schemaName))
	if err != nil {
		t.Fatalf("compile %s: %v", schemaName, err)
	}
	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatalf("decode generated %s artifact: %v", schemaName, err)
	}
	if err := schema.Validate(value); err != nil {
		t.Fatalf("validate generated %s artifact: %v", schemaName, err)
	}
}

func changeTestEvidence(id string, role source.VersionRole, version, location, data string, now time.Time) pack.SourceEvidenceInput {
	bytes := []byte(data)
	return pack.SourceEvidenceInput{
		Artifact: source.PinnedArtifact{
			ID: id, Kind: source.ArtifactSDKTypes, VersionRole: role, Version: version,
			Location: location, License: "Apache-2.0", MediaType: "text/typescript",
			Digest: source.DigestBytes(bytes), RetrievedAt: now.Add(-time.Minute),
			FreshUntil: now.Add(time.Hour), Data: bytes,
		},
		SourceLocator: "PaymentIntentsResource.confirm",
	}
}
