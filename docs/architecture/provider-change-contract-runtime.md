# Provider Change Contract runtime

Status: implemented M3 boundary

## Purpose

M3 turns pinned provider evidence into two non-executable artifacts:

1. a versioned Provider Change Contract (lumyn.migration_pack); and
2. a signed Provider Change Event that identifies the exact contract bytes,
   audience, sequence, freshness window, lifecycle, and provider origin.

Neither artifact grants access to a consumer repository, selects an Agent
Runner, authorizes mutation, or carries executable provider code.

## Source intake and normalization

The internal/source PinArtifact function accepts OpenAPI descriptions,
documentation, SDK type or release artifacts, and migration guidance. Each
artifact must have:

- immutable bytes and their matching SHA-256 digest;
- source, target, or guidance role and an exact version;
- a credential-free HTTPS provenance location;
- retrieval and optional freshness timestamps;
- explicit license posture.

The internal/source SnapshotOpenAPI function emits the same sorted semantic
view for the supported JSON and YAML OpenAPI subset. The internal/pack
DeriveOpenAPIChanges function compares stable operation identities. A removed
or renamed operation remains ambiguous until an accountable reviewer supplies
an explicit source-to-target mapping. The current bounded snapshot does not
flatten request or response schema properties, required fields, or types, so
every unreviewed schema-bearing operation also fails closed as ambiguous
rather than being reported unchanged.

The internal/pack Normalize function accepts only pinned artifacts and typed
declarations. Every semantic change cites at least one source-side and one
target-side evidence identity and concrete locator. Conflicting declarations
and unresolved ambiguity return a reviewable contract plus a blocking error.
Raw documentation and provider prose remain untrusted bytes; they are never
interpreted as commands, and mapping fields reject executable directives.

Public evidence produces provenance class public_derived and provider
confirmation status not_confirmed. It must never be displayed as provider
endorsement. The internal/pack Confirm function creates a new
provider_confirmed contract and a separate Ed25519-signed confirmation record.
VerifyConfirmation binds the operator, contract identity, version, semantic
digest, and confirmation time, so the same confirmed contract can be reused
across the exact invited cohort without repeated consumer-specific
interpretation.

## Digest conventions

The contract digest is the SHA-256 digest of canonical JSON with the
self-referential digest and separate confirmation fields cleared. It binds the
normalized semantics, evidence, provenance class, audience, lifecycle, and
safety flags.

The Provider Change Event separately records the SHA-256 digest of the exact
canonical contract bytes delivered. This byte digest detects serialization or
confirmation-record changes in transit. The event digest is computed with its
event-digest and detached-signature values cleared; the Ed25519 signature then
covers the complete event including that digest.

## Publish and intake

The internal/change BuildPublishKit function requires a provider-confirmed
active contract, its independently signed confirmation record, and the
confirmation public key before it emits canonical contract and event JSON. It
validates that manifest and contract URLs are at the same pinned,
credential-free provider HTTPS origin. Signing keys are used in memory and are
never serialized.

The internal/change VerifyPublishKit function independently checks:

- event digest and detached Ed25519 signature;
- exact enrolled campaign key identity and public key from the Consumer
  Installation boundary;
- exact contract-byte digest and semantic contract digest;
- provider, API or SDK, audience, contract identity, and version bindings;
- pinned origin and manifest and contract URLs;
- monotonic sequence, duplicate or conflicting identity, issue and expiry
  time, maximum age, deadline, lifecycle, and audience.

The caller must supply the observed intake mode. Only bytes actually observed
through the pinned provider HTTPS route qualify as provider-channel delivery
and can reach installed-preauthorization policy. Identical signed bytes
received through attended import are accepted only as recovery input and
return both channel-delivery and preauthorization eligibility as false.

## Verification

Focused tests:

    go test ./internal/source ./internal/pack ./internal/change ./schemas

Repository gate:

    make prepush-full

The generated runtime artifacts are validated against the migration-pack and
provider-change-event schemas. Tests cover JSON/YAML equivalence, digest
mismatch, staleness, missing licensing, conflicts, ambiguity, public-derived
versus provider-confirmed state, executable mapping denial, tampering, replay,
wrong audience, pinned delivery, and attended recovery labeling.
