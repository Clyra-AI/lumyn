# ADR-0006: Command-Result v1 Provider Terminology

## Status

Accepted for M0.

## Context

The retained `lumyn.command_result` `1.0` contract and its evidence-family
schemas serialize model execution identity under `provider_metadata`. The key
predates Lumyn's separation of the API Provider, Agent Runner Vendor, and Model
Provider roles. Reinterpreting the existing key as the API Provider would
silently change historical evidence; renaming it in place would invalidate
persisted v1.0 artifacts.

## Decision

Version `1.0` keeps the serialized `provider_metadata` key. In every v1.0
command result, result-axes artifact, evidence event, cassette event, and trace
event, that key means **Model Provider metadata only**. Its nested `provider`
value is the Model Provider identifier. It never identifies the API Provider,
change authority, Agent Runner Vendor, Lumyn Operator, or API Consumer
Organization.

New M0 command-result writers emit the optional discriminator
`semantic_role: model_provider`. Historical v1.0 artifacts that omit the
discriminator remain valid. Because v1.0 previously allowed unknown object
properties, its schemas must also continue accepting historical payloads whose
unknown extensions happen to use that field with another value. Such an
extension is not an authoritative role claim: the v1.0 `provider_metadata` key
still means Model Provider metadata, and current writers emit only
`model_provider`. Go identifiers use `ModelProviderMetadata` and
`ModelProvider` even though the compatibility JSON keys remain unchanged.

API Provider identity uses `api_provider_id` or `change_authority`. Agent
Runner Vendor identity uses `agent_runner_vendor_metadata`. Neither role may
be written into the v1.0 `provider_metadata` object.

## Versioned Migration Path

A future serialized rename is a breaking contract change and therefore uses
version `2.0`, not an in-place change to v1.0:

1. Freeze executable v1.0 schema snapshots before introducing v2.0.
2. Add v2.0 schemas whose canonical key is `model_provider_metadata` and whose
   API Provider and Agent Runner Vendor fields remain separate.
3. Accept v1.0 and v2.0 through version-discriminated readers and normalize
   both into the explicit internal Model Provider type.
4. Make v2.0 writers emit only `model_provider_metadata`; v2.0 schemas reject
   payloads that mix legacy and current keys or assign conflicting roles.
5. Preserve historical v1.0 bytes and digests. Migration creates a new artifact
   with provenance linking the source digest; it never rewrites or backfills
   historical evidence in place.
6. Version containing evidence artifacts coherently. A cassette or trace may
   not silently reinterpret a nested v1.0 evidence event as v2.0.

Every versioned change ships schema files, valid current and legacy fixtures,
invalid v2.0 mixed-role fixtures, reader/writer compatibility tests, and
migration notes together.

## Consequences

- Existing v1.0 artifacts and integrations remain valid.
- New output is self-describing without claiming a breaking schema migration.
- New result writers make Model Provider intent explicit without narrowing the
  accepted v1.0 language or treating unknown historical extensions as
  authority.
- The eventual clean JSON rename has an explicit major-version boundary and
  provenance-preserving migration rule.

## Validation

- CLI and Go envelope tests require the Model Provider discriminator on new
  command results.
- Schema fixtures prove current v1.0, historical v1.0, and previously permitted
  unknown-extension payloads remain valid.
- Writer tests prove current output emits only the Model Provider role.
- `make test-contracts` compiles and exercises all retained schema families.
