# Changelog

## [Unreleased]

### Added

- Initial Factory-compatible repository operating pack.
- T2 executable JSON Schema contracts for workflow, evidence, cassette, proof, command-result, redaction, and related artifact models.
- Local safety/corpus-ready normalized result and failure evidence fields with `corpus_eligible: false` validation.
- A provider-sponsored, customer-controlled API migration PRD, architecture
  decision, implementation plan, and active Factory control set.
- Planned `lumyn artifacts gc [--dry-run]` recovery for private-artifact
  TTL/revocation deletion, receipts, and orphan reporting.
- Planned standard campaign-kit creation, configured packet publication, signed
  invitations, no-authority consumer `campaign accept`, explicit signed
  authorization issue/revoke, signed provider-status snapshots, minimal signed
  connection receipts, and a receipt-backed synthetic offline canary through a
  local draft-PR preview.

### Changed

- Reframed Lumyn from generic agent-readiness evaluation to verified migration
  campaigns for API providers and authorized TypeScript/Node customer
  repositories.
- Made the API provider the campaign sponsor and economic buyer while keeping
  repository reads, writes, commands, credentials, disclosure, PR creation,
  review, and merge under explicit customer authority.
- Narrowed the MVP to declarative provider change packets, read-only impact,
  three deterministic transformations, evidence-scoped verification, and
  customer-authorized draft GitHub PRs.
- Reclassified the prior `lumyn-mvp` plan and evidence as immutable historical
  records; the active plan is now `lumyn-migration-mvp`.
- Preregistered commercial validation around five distinct eligible
  repositories across at least three organizations, supported-change canary
  evidence, fixed deadlines, post-merge observation, non-negative campaign
  contribution, and executed paid continuation.
- Added real independent holdout, trace-grade, and evidence-attestation
  promotion gates before commit/push for the tasks that rely on those proofs.
- Split M8 sandbox verification from M9 draft-PR delivery and made provider
  reporting optional, so consumers can authorize a verified draft PR without
  granting sandbox or provider-attestation access.
- Made design-partner sandbox proof action-specific: the campaign can use
  exact-patched-head mock proof without a sandbox grant, while separately
  authorized repositories may add sandbox evidence.
- Added invitation-to-impact elapsed time and total consumer security,
  privacy, platform, and maintainer labor to activation, so the pilot cannot
  hide pre-acceptance onboarding friction.
- Added recurring-value qualification plus preregistered annual/second-campaign
  margin and Lumyn operator-effort gates; provider and consumer labor remains
  separate buyer total cost.
- Defined acknowledged, consumer-signed minimal connection receipts as the
  only sponsored connected-repository meter, with provider-authenticated
  consumer signer binding, one-invitation-unit cardinality, pinned online or
  offline exchange, and provider-signed deduplicated acknowledgements, while
  keeping richer provider reporting separately optional.
- Replaced the premature OSS claim with explicit design-partner distribution
  terms and a separate license/security/support/release-integrity gate for any
  future public OSS or self-serve release.

### Deprecated

- Generic live agent evaluation, model-provider panels, public API teardown
  content, and buy-side monitoring of every vendor as mandatory MVP scope.

### Removed

### Fixed

- Planning now records that unimplemented `record`, `verify`, `trace`, `demo`,
  `share`, and `eval` commands must fail closed before migration execution work
  begins.
- Product documentation no longer describes recorder, replay, live
  verification, reporting, GitHub delivery, migration patching, or live agent
  evaluation as implemented.

### Security

- Defined separate provider and consumer data planes, signed immutable
  provider packets with a consumer-pinned trust root and provider/package
  binding, independently scoped capabilities, private artifact
  TTL/deletion, no provider access to raw customer repository data by default,
  bounded writes, sandbox isolation, draft-only PRs, and human merge authority.
- Separated Factory worker `approval`/`credentials`/`network` grants from exact
  private Lumyn product authorizations, and separated field-allowlisted
  provider attestations from aggregate/hash-only public evidence.
- Required current packet trust and exact product authority to be revalidated
  inside Lumyn immediately before every local, sandbox, GitHub, or attestation
  side effect and retry; Factory dispatch and standalone validation output do
  not confer runtime authority.
- Required every trust refresh to use either a signed offline provider-status
  snapshot inside the pinned maximum age or an exact authorized endpoint read,
  and required repository commands to use a fail-closed host-isolation backend
  that denies undeclared mounts, credentials, sockets, descriptors, and child
  escape.
- Required design-partner qualification to prove an operational provider-status
  channel and an enforceable supported host-isolation backend for each
  candidate environment before migration-plan work can begin.
- Made provider export and public commit explicit irreversible disclosure
  boundaries: revocation blocks future sharing and deletes Lumyn-controlled
  private copies but cannot claim to recall recipient copies or public history.
- Added a narrow manual privacy/legal preflight before identifiable
  design-partner evidence collection and automatic private-artifact
  TTL/revocation enforcement with retry-safe deletion receipts and orphan
  recovery.
- Bound that preflight to an exact scope digest, rejected semantic-wildcard
  Factory grants, made first-time provider trust enrollment independently
  authenticated, and restricted the committed holdout manifest to
  non-resolving provenance and keyed/encrypted commitments. Held-out inputs and
  answer keys stay hidden from implementation workers, and independent
  lifecycle evidence is schema-valid,
  current-candidate-bound, provenance-backed, and unwritable by the executor.
- Split holdout policy into non-circular provision and evaluate modes: M1
  creates the keyed suite commitment independently, while later scoring tasks
  resolve and byte-bind the trusted M1 result instead of fabricating a future
  commitment in the static plan.
