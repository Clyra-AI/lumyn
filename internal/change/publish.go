package change

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/Clyra-AI/lumyn/internal/installation"
	"github.com/Clyra-AI/lumyn/internal/pack"
	"github.com/Clyra-AI/lumyn/internal/source"
)

type EventArtifact struct {
	ObjectType                        string              `json:"object_type"`
	SchemaVersion                     string              `json:"schema_version"`
	EventID                           string              `json:"event_id"`
	EventVersion                      int                 `json:"event_version"`
	EventDigest                       string              `json:"event_digest"`
	APIProviderID                     string              `json:"api_provider_id"`
	Issuer                            EventIssuer         `json:"issuer"`
	APIOrSDK                          pack.APIOrSDK       `json:"api_or_sdk"`
	AudienceID                        string              `json:"audience_id"`
	Deadline                          time.Time           `json:"deadline"`
	Severity                          string              `json:"severity"`
	Sequence                          uint64              `json:"sequence"`
	IssuedAt                          time.Time           `json:"issued_at"`
	ExpiresAt                         time.Time           `json:"expires_at"`
	Transport                         EventTransport      `json:"transport"`
	Authentication                    EventAuthentication `json:"authentication"`
	ContractDelivery                  ContractDelivery    `json:"contract_delivery"`
	Lifecycle                         pack.Lifecycle      `json:"lifecycle"`
	ChannelDeliveryQualified          bool                `json:"channel_delivery_qualified"`
	InstalledPreauthorizationEligible bool                `json:"installed_preauthorization_eligible"`
	NonExecutable                     bool                `json:"non_executable"`
	GrantsConsumerAuthority           bool                `json:"grants_consumer_authority"`
	ProductionCredentialsAllowed      bool                `json:"production_credentials_allowed"`
	ProductionMutationAllowed         bool                `json:"production_mutation_allowed"`
}

type EventIssuer struct {
	Role               string `json:"role"`
	ProviderOperatorID string `json:"provider_operator_id"`
}

type EventTransport struct {
	Mode          string `json:"mode"`
	ManifestURL   string `json:"manifest_url"`
	PinnedOrigin  string `json:"pinned_origin"`
	PublisherRole string `json:"publisher_role"`
}

type EventAuthentication struct {
	Method            string `json:"method"`
	KeyID             string `json:"key_id"`
	DetachedSignature string `json:"detached_signature"`
	SignatureVerified bool   `json:"signature_verified"`
	OriginVerified    bool   `json:"origin_verified"`
	AudienceVerified  bool   `json:"audience_verified"`
	FreshnessVerified bool   `json:"freshness_verified"`
	SequenceVerified  bool   `json:"sequence_verified"`
	ReplayCheck       string `json:"replay_check"`
}

type ContractDelivery struct {
	Mode                 string `json:"mode"`
	MigrationPackID      string `json:"migration_pack_id"`
	ContractVersion      int    `json:"contract_version"`
	ContractURL          string `json:"contract_url,omitempty"`
	RetrievedBytesDigest string `json:"retrieved_bytes_digest"`
	DigestVerified       bool   `json:"digest_verified"`
}

type PublishRequest struct {
	EventID            string
	EventVersion       int
	ProviderOperatorID string
	KeyID              string
	AudienceID         string
	Deadline           time.Time
	Severity           string
	Sequence           uint64
	IssuedAt           time.Time
	ExpiresAt          time.Time
	PinnedOrigin       string
	ManifestURL        string
	ContractURL        string
}

type PublishKit struct {
	ContractBytes []byte
	EventBytes    []byte
}

// BuildPublishKit emits only inert JSON bytes. The private key signs the
// event in memory and is never serialized into the kit.
func BuildPublishKit(
	contract pack.Contract,
	confirmation pack.ConfirmationRecord,
	confirmationPublicKey ed25519.PublicKey,
	request PublishRequest,
	privateKey ed25519.PrivateKey,
) (PublishKit, error) {
	if err := pack.VerifyConfirmation(contract, confirmation, confirmationPublicKey); err != nil {
		return PublishKit{}, fmt.Errorf("publish Provider Change Contract: %w", err)
	}
	if request.ProviderOperatorID != *contract.ProviderConfirmation.ProviderOperatorID {
		return PublishKit{}, errors.New("event issuer does not match contract confirmer")
	}
	if request.EventVersion < 1 || request.Sequence < 1 || len(privateKey) != ed25519.PrivateKeySize {
		return PublishKit{}, errors.New("event version, sequence, and Ed25519 signing key are required")
	}
	for label, raw := range map[string]string{
		"pinned origin": request.PinnedOrigin, "manifest URL": request.ManifestURL, "contract URL": request.ContractURL,
	} {
		parsed, err := url.Parse(raw)
		if err != nil || parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil ||
			parsed.RawQuery != "" || parsed.Fragment != "" {
			return PublishKit{}, fmt.Errorf("%s must be a credential-free HTTPS URL", label)
		}
	}
	if err := installation.ValidateURLAtPinnedOrigin(request.PinnedOrigin, request.ManifestURL); err != nil {
		return PublishKit{}, fmt.Errorf("provider event manifest: %w", err)
	}
	if err := installation.ValidateURLAtPinnedOrigin(request.PinnedOrigin, request.ContractURL); err != nil {
		return PublishKit{}, fmt.Errorf("Provider Change Contract URL: %w", err)
	}
	if request.EventID == "" || request.KeyID == "" || request.AudienceID != contract.Audience.AudienceID {
		return PublishKit{}, errors.New("event identity, key identity, and exact contract audience are required")
	}
	switch request.Severity {
	case "informational", "recommended", "breaking", "sunset_critical":
	default:
		return PublishKit{}, fmt.Errorf("unsupported provider event severity %q", request.Severity)
	}
	if request.IssuedAt.IsZero() || request.ExpiresAt.IsZero() || !request.ExpiresAt.After(request.IssuedAt) ||
		request.Deadline.IsZero() || !request.Deadline.After(request.IssuedAt) {
		return PublishKit{}, errors.New("event issue, expiry, and deadline times are invalid")
	}
	contractBytes, err := pack.CanonicalBytes(contract)
	if err != nil {
		return PublishKit{}, err
	}
	event := EventArtifact{
		ObjectType: "lumyn.provider_change_event", SchemaVersion: pack.SchemaVersion,
		EventID: request.EventID, EventVersion: request.EventVersion,
		APIProviderID: contract.APIProviderID,
		Issuer:        EventIssuer{Role: "provider_operator", ProviderOperatorID: request.ProviderOperatorID},
		APIOrSDK:      contract.APIOrSDK, AudienceID: request.AudienceID,
		Deadline: request.Deadline.UTC(), Severity: request.Severity, Sequence: request.Sequence,
		IssuedAt: request.IssuedAt.UTC(), ExpiresAt: request.ExpiresAt.UTC(),
		Transport: EventTransport{
			Mode: "pinned_provider_https", ManifestURL: request.ManifestURL,
			PinnedOrigin: request.PinnedOrigin, PublisherRole: "provider_operator",
		},
		Authentication: EventAuthentication{
			Method: "detached_signature", KeyID: request.KeyID,
			SignatureVerified: true, OriginVerified: true, AudienceVerified: true,
			FreshnessVerified: true, SequenceVerified: true, ReplayCheck: "unseen",
		},
		ContractDelivery: ContractDelivery{
			Mode: "exact_provider_https_url", MigrationPackID: contract.MigrationPackID,
			ContractVersion: contract.ContractVersion, ContractURL: request.ContractURL,
			RetrievedBytesDigest: source.DigestBytes(contractBytes), DigestVerified: true,
		},
		Lifecycle:                pack.Lifecycle{State: "active"},
		ChannelDeliveryQualified: true, InstalledPreauthorizationEligible: true,
		NonExecutable: true, GrantsConsumerAuthority: false,
		ProductionCredentialsAllowed: false, ProductionMutationAllowed: false,
	}
	event.EventDigest = computeEventDigest(event)
	unsigned, _ := eventSignatureBytes(event)
	event.Authentication.DetachedSignature = base64.RawURLEncoding.EncodeToString(ed25519.Sign(privateKey, unsigned))
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return PublishKit{}, err
	}
	return PublishKit{ContractBytes: contractBytes, EventBytes: eventBytes}, nil
}

type IntakeObservation struct {
	Mode                string
	ObservedManifestURL string
}

type IntakeResult struct {
	Event                    EventArtifact
	Contract                 pack.Contract
	IntakeMode               string
	ChannelDeliveryQualified bool
	PreauthorizationEligible bool
}

// VerifyPublishKit validates exact bytes and detached signature. The observed
// intake mode is external evidence: copied bytes are recovery input even when
// their signed manifest names a provider channel.
func VerifyPublishKit(kit PublishKit, publicKey ed25519.PublicKey, context EventContext, observation IntakeObservation) (IntakeResult, error) {
	if len(publicKey) != ed25519.PublicKeySize {
		return IntakeResult{}, errors.New("provider event Ed25519 public key is invalid")
	}
	var event EventArtifact
	if err := json.Unmarshal(kit.EventBytes, &event); err != nil {
		return IntakeResult{}, fmt.Errorf("decode provider event: %w", err)
	}
	if err := validateEventArtifact(event); err != nil {
		return IntakeResult{}, err
	}
	signature, err := base64.RawURLEncoding.DecodeString(event.Authentication.DetachedSignature)
	if err != nil {
		return IntakeResult{}, errors.New("provider event detached signature is not base64url")
	}
	unsigned, _ := eventSignatureBytes(event)
	if !ed25519.Verify(publicKey, unsigned, signature) {
		return IntakeResult{}, errors.New("provider event detached signature is invalid")
	}
	if event.EventDigest != computeEventDigest(event) {
		return IntakeResult{}, errors.New("provider event digest does not bind immutable event content")
	}
	if event.ContractDelivery.RetrievedBytesDigest != source.DigestBytes(kit.ContractBytes) {
		return IntakeResult{}, errors.New("retrieved Provider Change Contract bytes do not match event digest")
	}
	var contract pack.Contract
	if err := json.Unmarshal(kit.ContractBytes, &contract); err != nil {
		return IntakeResult{}, fmt.Errorf("decode Provider Change Contract: %w", err)
	}
	if err := pack.ValidateContract(contract); err != nil {
		return IntakeResult{}, err
	}
	if event.APIProviderID != contract.APIProviderID || event.APIOrSDK != contract.APIOrSDK ||
		event.ContractDelivery.MigrationPackID != contract.MigrationPackID ||
		event.ContractDelivery.ContractVersion != contract.ContractVersion ||
		event.AudienceID != contract.Audience.AudienceID {
		return IntakeResult{}, errors.New("provider event does not identify the exact Provider Change Contract and audience")
	}
	if contract.ProviderConfirmation.ProviderOperatorID == nil ||
		event.Issuer.ProviderOperatorID != *contract.ProviderConfirmation.ProviderOperatorID {
		return IntakeResult{}, errors.New("provider event issuer does not match the confirmed contract operator")
	}
	runtimeEvent := ProviderEvent{
		EventID: event.EventID, EventVersion: fmt.Sprintf("%d", event.EventVersion),
		EventDigest: event.EventDigest, Issuer: event.Issuer.ProviderOperatorID,
		APIOrSDK: event.APIOrSDK.Package, Audience: []string{event.AudienceID},
		Deadline: event.Deadline, Severity: event.Severity, Sequence: event.Sequence,
		IssuedAt: event.IssuedAt, ExpiresAt: event.ExpiresAt,
		TransportMode: event.Transport.Mode, TransportOrigin: event.Transport.PinnedOrigin,
		ManifestURL: event.Transport.ManifestURL, Authenticated: true, SignatureValid: true,
		SignatureProvenance:     event.Authentication.KeyID,
		ContractDigest:          event.ContractDelivery.RetrievedBytesDigest,
		RetrievedContractDigest: source.DigestBytes(kit.ContractBytes),
		ContractDeliveryMode:    event.ContractDelivery.Mode, ContractURL: event.ContractDelivery.ContractURL,
		ProviderChannelDeliveryProven: true, Lifecycle: event.Lifecycle.State,
		Executable: !event.NonExecutable,
	}
	switch observation.Mode {
	case "pinned_provider_https":
		if observation.ObservedManifestURL != event.Transport.ManifestURL {
			return IntakeResult{}, errors.New("observed provider manifest URL does not match signed event")
		}
		if err := ValidateProviderEvent(runtimeEvent, context); err != nil {
			return IntakeResult{}, err
		}
		return IntakeResult{
			Event: event, Contract: contract, IntakeMode: observation.Mode,
			ChannelDeliveryQualified: true, PreauthorizationEligible: true,
		}, nil
	case "attended_import":
		// Revalidate all semantic controls while deliberately removing the
		// unobserved channel and fetch claims.
		runtimeEvent.TransportMode = "attended_import"
		runtimeEvent.ProviderChannelDeliveryProven = false
		runtimeEvent.ContractDeliveryMode = "embedded"
		runtimeEvent.ContractURL = ""
		if err := ValidateProviderEvent(runtimeEvent, context); err != nil {
			return IntakeResult{}, err
		}
		return IntakeResult{
			Event: event, Contract: contract, IntakeMode: observation.Mode,
			ChannelDeliveryQualified: false, PreauthorizationEligible: false,
		}, nil
	default:
		return IntakeResult{}, fmt.Errorf("unsupported observed intake mode %q", observation.Mode)
	}
}

func validateEventArtifact(event EventArtifact) error {
	if event.ObjectType != "lumyn.provider_change_event" || event.SchemaVersion != pack.SchemaVersion {
		return errors.New("unsupported provider event type or schema version")
	}
	if event.EventID == "" || event.EventVersion < 1 || event.EventDigest == "" ||
		event.APIProviderID == "" || event.AudienceID == "" || event.Sequence < 1 ||
		event.Issuer.Role != "provider_operator" || event.Issuer.ProviderOperatorID == "" ||
		event.APIOrSDK.Name == "" || event.APIOrSDK.Package == "" {
		return errors.New("provider event identity is incomplete")
	}
	if !event.NonExecutable || event.GrantsConsumerAuthority ||
		event.ProductionCredentialsAllowed || event.ProductionMutationAllowed {
		return errors.New("provider event must remain non-executable and authority-free")
	}
	if event.Transport.Mode != "pinned_provider_https" ||
		event.Transport.PublisherRole != "provider_operator" ||
		event.ContractDelivery.Mode != "exact_provider_https_url" ||
		event.ContractDelivery.MigrationPackID == "" ||
		event.ContractDelivery.ContractVersion < 1 ||
		event.ContractDelivery.RetrievedBytesDigest == "" ||
		!event.ContractDelivery.DigestVerified {
		return errors.New("provider event transport or contract delivery is incomplete")
	}
	authentication := event.Authentication
	if authentication.Method != "detached_signature" || authentication.KeyID == "" ||
		authentication.DetachedSignature == "" || !authentication.SignatureVerified ||
		!authentication.OriginVerified || !authentication.AudienceVerified ||
		!authentication.FreshnessVerified || !authentication.SequenceVerified ||
		authentication.ReplayCheck != "unseen" {
		return errors.New("provider event authentication claims are incomplete")
	}
	if event.Lifecycle.State != "active" || event.Lifecycle.SupersededBy != nil ||
		event.Lifecycle.WithdrawalReason != nil ||
		!event.ChannelDeliveryQualified || !event.InstalledPreauthorizationEligible {
		return errors.New("provider event is not an active channel-delivery candidate")
	}
	return nil
}

func computeEventDigest(event EventArtifact) string {
	event.EventDigest = ""
	event.Authentication.DetachedSignature = ""
	data, _ := json.Marshal(event)
	return source.DigestBytes(data)
}

func eventSignatureBytes(event EventArtifact) ([]byte, error) {
	event.Authentication.DetachedSignature = ""
	return json.Marshal(event)
}
