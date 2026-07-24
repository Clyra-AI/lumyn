# ADR-0002: Provider-Sponsored, Customer-Controlled API Migrations

## Status

Accepted for the Lumyn migration MVP rebaseline.

## Date

2026-07-23

## Context

Lumyn's first product thesis treated API documentation and workflows as a
surface to evaluate for agent readiness. The repository produced useful
foundations for stable command results, executable evidence contracts, local
OpenAPI and documentation intake, redaction, and governed delivery. It did not
implement the recorder, replay verifier, live verifier, report renderer,
GitHub delivery, migration analysis, patching, or live agent-evaluation runtime
described by the historical plan.

The stronger commercial problem is a consequential API or SDK migration that an
API provider must complete across customer-owned integrations. The provider can
fund and coordinate that campaign, but it cannot grant access to a customer's
repository, credentials, tests, or merge controls. The customer bears the code
and production risk and must retain authority over every sensitive action.

This creates two principals with different authority:

- the API provider owns authoritative change intent and sponsors the campaign;
- the API consumer owns repository access, execution, disclosure, review, and
  merge.

## Decision

Lumyn will build provider-sponsored, customer-controlled verified API
migrations.

The initial product accepts a signed, immutable-for-its-authorized-audience,
declarative provider change packet; analyzes an explicitly authorized
TypeScript/Node repository; shows a read-only impact report and migration plan;
applies only supported, bounded transformations after separate approval; runs
evidence-scoped verification; and may open an idempotent draft GitHub PR under
customer authority. It never writes the default branch or auto-merges in the
MVP.

The activation contract is also part of the product: a Provider Operator uses
a standard campaign kit and configured signer to publish the packet and signed
invitation; a Consumer Maintainer first enrolls the provider root/package
binding from a separately authenticated bundle and out-of-band confirmed
fingerprint, then verifies and accepts the invitation locally. Enrollment and
acceptance run without the checkout mounted, persist only to private state, and
produce a no-authority authorization request. The maintainer must explicitly
issue or revoke a signed authorization bundle; issuance performs no granted
side effect. A sponsored connection counts only from a separate
consumer-signed minimal receipt whose signer is authenticated to the invited
organization and whose one-repository invitation unit receives a
provider-signed, deduplicated acknowledgement. The invitation cannot
authenticate its own root. A synthetic offline canary must exercise that flow
through host-isolated verification and a local draft-PR preview, returning a
typed nonzero result if any implementation stage is missing.

The deployment model has two data planes:

1. The provider campaign plane owns change packets, invitations,
   authenticated consumer receipt-key bindings, provider-signed minimal
   connection-receipt acknowledgements, and optional customer-consented richer
   status attestations.
2. The consumer execution plane owns source code, credentials, commands,
   patches, raw logs, private evidence, and remote-branch, PR, and merge
   authorization.

Provider sponsorship conveys no customer-data authority. Provider-supplied
packets are treated as untrusted supply-chain input until their signatures
verify against a consumer-pinned API-provider trust root and verified
provider-to-package ownership binding and their key, timestamp, audience,
expiry, rotation, revocation, withdrawal, replay, immutability, schema, and
canary checks pass. `Published` means immutable for the authorized audience,
not publicly disclosed. Raw code, diffs, logs, traces, and credentials do not
cross from the consumer plane to the provider plane by default.

First-time trust enrollment cannot authenticate itself from invitation
material. A consumer obtains the provider enrollment bundle and expected root
fingerprint through a separately authenticated provider admin/security
channel, then explicitly pins the provider root and provider/package binding
outside the checkout before accepting an invitation. Normal rotation is signed
by the active root. Emergency recovery requires explicit re-enrollment,
freezes affected campaigns, and invalidates prior approvals.
The enrollment also pins a provider-status signer, maximum status age, and an
exact endpoint or offline snapshot contract, plus the provider
receipt-acknowledgement signer and permitted online/offline receipt exchange
classes. Every side effect consumes either a signed current offline snapshot
or an exact endpoint read under `provider_trust_status_read`; missing or
undeclared freshness blocks. Invitations may narrow enrolled receipt settings
but cannot introduce a new signer or destination.

Consumer-private runtime artifacts live in an explicitly configured,
non-committable root outside the consumer checkout and every public source
repository. They enforce authorization TTL and deletion on expiry or
revocation automatically at creation, read, process startup, and the next run,
with deletion receipts or orphan reports. `lumyn artifacts gc` retries partial
deletion and surfaces unresolved orphans without extending authority or
rewriting historical closure. Provider-visible and public campaign evidence is
stored separately. Provider-visible attestations contain only exact
consumer-consented fields; public campaign evidence contains only consented
redacted aggregates or hashes. Export and public commit are irreversible:
revocation blocks future sharing and deletes Lumyn-controlled private copies,
but cannot recall provider records, Git history, clones, or caches.

Repository read, local write, host-isolated command execution, provider-status
read, package-registry network, sandbox request disclosure, sandbox network,
sandbox credential, remote branch write, PR write, minimal campaign receipt,
richer provider attestation, retention, and deletion are independent,
time-bounded, revocable authorities. In particular, remote branch write does
not imply PR write, and PR write does not imply remote branch write. Repository
commands expose only exact mounts and a sanitized environment; host home,
credential stores, OS credentials, agent/Docker/local-service sockets, and
extra inherited descriptors are denied, child processes inherit the boundary,
and an unenforceable backend blocks.
Retention and deletion use explicit `artifact_retention` and
`artifact_deletion` grants; every data-producing or disclosure grant references
both.

Sandbox verification and draft-PR delivery are independent after offline
verification. M9 does not wait for M8 and does not require
`provider_attestation`; reporting is an optional, separately consented action.

These Lumyn product authorizations are private schema-backed artifacts, not
free-form Factory capabilities. Factory retains its closed worker
`approval`, `credentials`, and `network` vocabulary. A Factory approval can
cite a validated opaque product-authorization bundle, but cannot replace or
broaden any consumer grant. Factory selection and dispatch do not validate the
private product bundle. Lumyn revalidates current packet trust and the exact
applicable bundle internally immediately before every local, sandbox, GitHub,
or attestation side effect and retry, using current signed provider-status
evidence; a separate validation command is diagnostic and closure evidence, not
cached runtime authority.

M2.5 external-evidence collection has a narrower boundary. Before collecting,
storing, or disclosing identifiable partner evidence, the Factory
implementation worker requires a manual privacy/legal preflight naming allowed
fields, participant consent, the approved external private root, TTL,
expiry/revocation deletion, deletion-receipt and orphan ownership, and
the minimal connection-receipt, provider-authenticated consumer signer binding,
provider-signed acknowledgement/cardinality policy, and separately consented
aggregate/hash-only public fields. Consent states that provider/public copies
cannot be recalled.
This preflight does not grant Lumyn product runtime authority. Its active
Factory approval must cite the evidence and exactly match a canonical digest
of that scope.

Independent promotion evidence is also explicit. When selected by task policy,
`holdout-evaluator`, `trace-grader`, and `evidence-attestor` run in that order
after code review and before commit/push. They are external or human-operated
principals. Shipping verifies their task-bound, current-work-proof, passing
artifacts before creating a commit or PR. Every artifact schema-validates and
binds the exact task, work item, lifecycle run, current validation run,
candidate digest, work-proof marker digest, and independent worker provenance;
the implementation worker cannot access held-out answer material, write the
lifecycle namespace, self-grade, or self-attest.

An independent holdout owner provisions and freezes the benchmark private
suite. The repository contains only opaque case IDs, non-resolving provenance
class and license posture, a frozen suite commitment, and encrypted or HMAC
artifact commitments; it contains no source location, repository or package
identifier, plaintext content digest, or answer material. M1's
`holdout_result` binds the frozen commitment. Only task-selected workers run:
M10 uses `evidence-attestor` for privacy-approved campaign calculations and
never receives or reuses the benchmark holdout root.

The MVP is deliberately narrow:

- GitHub repositories;
- TypeScript/Node;
- direct imports from one official npm SDK;
- `package-lock.json` as the first writable lockfile, with exact Node and npm
  versions, registry endpoint or immutable snapshot, package-integrity inputs,
  toolchain digest, disabled lifecycle scripts, and separately authorized
  registry network;
- three deterministic change classes defined in the PRD;
- local or customer-CI execution by default;
- deterministic transformations before any model-assisted fallback;
- compile/typecheck and allowlisted repository tests;
- canonical `static_verified`, `repo_verified`,
  `workflow_contract_replay_passed`, `workflow_verified_replay`,
  `workflow_verified_mock`, and `workflow_verified_sandbox` labels;
- `workflow_verified_*` only after an approved entrypoint executes from the
  exact patched repository head and produces observed interaction and outcome
  evidence in that environment; independent contract/cassette replay is
  `workflow_contract_replay_passed` and cannot exceed `repo_verified`;
- provider-sandbox read-back only after separate endpoint, non-production
  credential, request-payload, synthetic or approved non-sensitive test-data,
  provider-logging, retention, deletion, budget, cleanup, and orphan consent;
- draft PR delivery only.

Public Stripe, GitHub, OpenAI, or other published API artifacts may seed pinned
engineering fixtures. They prove technical feasibility, not provider demand,
change authority, customer consent, or permission to operate a migration
campaign.

## Retained Foundations

The following existing work remains valid and must be reused without
overstating what is implemented:

- Go CLI, config, result-envelope, and typed-exit foundations;
- executable workflow, evidence, cassette, proof, redaction, and action-boundary
  schemas;
- local OpenAPI and documentation parsing, fingerprints, references,
  deprecations, and replacement-hint extraction;
- Factory lifecycle, validation, ownership, review, and provenance controls.

Existing bare `provider_metadata` continues to mean model-provider metadata
until a versioned result-contract migration introduces
`model_provider_metadata`. New migration contracts use `api_provider_id` and
`change_authority`; they must not overload `provider`.

## Superseded Plan

The former plan under `.factory/artifacts/prd-to-plan/lumyn-mvp/` is not active
product scope. Its task/control, PR, exception, and pilot artifacts are
immutable historical records and must not be rewritten to imply migration
behavior. The directory README carries only a non-operative dispatch tombstone.

The active plan is
`.factory/artifacts/prd-to-plan/lumyn-migration-mvp/`. It maps retained
foundation evidence explicitly and gives new task IDs to work that has not
been implemented. Its ledger maps the PRD's exactly 62 item-level acceptance
controls.

## Consequences

Positive consequences:

- the economic buyer and business outcome are explicit;
- customer authorization is designed into the product rather than left as an
  integration detail;
- the payer receives a cryptographically verifiable, privacy-minimal connection
  meter without gaining repository access;
- existing source and evidence work becomes a useful migration-verification
  kernel;
- the MVP can prove value with a bounded campaign instead of attempting generic
  coverage of every vendor and language;
- verification claims name their exact evidence boundary.

Costs and risks:

- adoption is two-sided even though monetization is provider-led;
- a provider change packet is a privileged supply-chain channel and requires
  signing, a consumer-pinned trust root, provider-to-package binding,
  lifecycle/freshness enforcement, immutable references, canary proof, and
  fail-closed parsing;
- real migrations may require business context that cannot be inferred safely;
- customer repository authorization and sandbox credentials may be harder to
  obtain than provider budget;
- host-isolation portability and provider-status availability may block a run
  on otherwise supported repositories;
- a services-assisted first campaign is expected while packet authoring and
  customer onboarding are learned.
- a one-off successful campaign is insufficient commercial proof; the provider
  must support an annual connected-repository program or a second named
  migration with preregistered recurring margin and operator-effort bounds.

## Rejected Alternatives

### Buy-Side Monitoring Of Every Vendor

Rejected as the initial wedge because public documentation rarely contains
enough authoritative semantic intent to make consequential patches safely, and
the approach lacks a natural distribution channel into customer repositories.
Read-only monitoring may become a later acquisition or discovery surface.

### Provider-Controlled Repository Automation

Rejected because provider payment or campaign sponsorship cannot authorize
access to customer code, credentials, commands, or merge controls.

### Generic Model-Generated Migration Patches

Rejected for the MVP because model availability is not the defensible input and
cannot replace authoritative change semantics, bounded transforms, or workflow
evidence. A model-assisted fallback may be considered later behind an explicit
consumer-controlled boundary and separate ADR.

### Automatic Merge

Rejected for the MVP. Branch protection, customer review, and human merge
authority remain part of the safety contract.

## Rollout And Rollback

Rollout begins with the product/plan rebaseline and a fail-closed correction to
unimplemented commands. Benchmark, contract, and read-only impact work may
advance while M2.5 qualifies one API Provider, validates a supported-class
canary, and freezes five distinct eligible repositories across at least three
Consumer Organizations plus the pilot protocol. Qualification includes an
operational signed provider-status channel, the provider acknowledgement key
and receipt exchange, authenticated consumer receipt-signer binding, and a
supported OS/architecture with enforceable host isolation for every candidate
environment. No-write migration-plan implementation does not begin until M2.5
closes `DISC-001` and `DISC-002` with direct external evidence and independent
attestation. M5 also carries explicit
runner-enforced item gates for both IDs; a task dependency alone cannot release
it. The patch engine, verification ladder, sandbox boundary, separate
remote-branch and draft-PR delivery, and design-partner campaign then advance
as separate evidence-backed tasks. M2 first establishes the standard campaign
kit and invitation contracts; after M7, M8 sandbox proof and M9 draft-PR
delivery may proceed independently.

If provider recruitment, customer consent, deterministic coverage, or
willingness to pay fails the PRD's falsification gates, stop campaign-plane
investment. The retained source-ingestion and workflow-evidence foundations may
still support a narrower customer-side migration analyzer. Historical artifacts
remain unchanged so rollback never requires rewriting evidence. Rollback also
revokes active product capabilities and deletes TTL-bound private artifacts;
failed deletion produces an orphan report rather than a false clean state.
It stops future provider/public disclosure but does not claim to recall copies
already exported or committed. The current repository and pilot distribution
are not called OSS; public OSS/self-serve requires a separate approved license,
security, contribution, support, vulnerability-response, and release-integrity
gate.
