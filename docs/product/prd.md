# Lumyn MVP — Provider-Sponsored Verified API Migrations

| Field | Value |
|---|---|
| Version | 2.0 |
| Status | Ready for execution after product rebaseline |
| Owner | Product and Engineering |
| Last Updated | 2026-07-23 |
| Primary Audience | Lumyn builders, API-provider design partners, API-consumer maintainers, and technical investors |
| Source Task Plan | `docs/product/plan.md` |
| Active Factory Plan | `.factory/artifacts/prd-to-plan/lumyn-migration-mvp/` |

---

## Purpose

This document is the product source of truth for the Lumyn MVP.

Lumyn enables an API provider to sponsor a migration campaign that helps
participating API consumers move from an old API or SDK version to a new one.
For each customer-authorized repository, Lumyn identifies affected code,
prepares a bounded migration, verifies the repository and relevant business
workflow, and opens an evidence-backed draft pull request.

The product promise is:

> When an API provider changes an API or SDK, Lumyn finds the affected code in
> authorized customer repositories, opens the fix, and shows exactly what was
> and was not verified.

This version supersedes the previous product direction in which generic API
workflow verification and live agent evaluation were the headline product.
Workflow contracts, evidence schemas, redaction rules, and proof semantics
remain valuable, but they now serve as the verification foundation for an API
migration rather than the product's primary job.

The prior Factory plan under
`.factory/artifacts/prd-to-plan/lumyn-mvp/` is retained as historical evidence.
It is not an active execution plan. The active plan is
`.factory/artifacts/prd-to-plan/lumyn-migration-mvp/`.

---

## Executive Summary

API providers publish migration guides, changelogs, SDK releases, and
occasionally codemods. The consumer still has to determine whether its
repository is affected, translate generic guidance into its own abstractions,
change the code, run the right checks, and decide whether the integration still
performs the required business job.

Lumyn closes that loop:

```text
signed provider change intent
-> customer-authorized repository impact
-> reviewable migration plan
-> bounded patch
-> repository and workflow verification
-> evidence-backed draft PR
-> customer review and merge
```

Lumyn is provider-sponsored and customer-controlled:

- The API provider is the initial economic buyer and campaign sponsor.
- The API consumer owns the repository, credentials, execution environment,
  disclosure choices, and merge decision.
- Provider sponsorship never grants the provider access to customer code.
- The consumer explicitly approves read, write, command, network, credential,
  PR, and reporting scopes.
- Lumyn never auto-merges in the MVP.

The commercial outcome is not merely a generated PR. It is faster completion
of a provider migration campaign: more customers safely moved off an old
version, less support work, fewer migration incidents, and a shorter legacy
compatibility tail.

---

## Product Thesis

Dependency bots can update version numbers. Codemods can rewrite known syntax.
Coding agents can propose broader edits. None of those mechanisms alone proves
that:

- every relevant call site was found;
- the patch stayed within the authorized scope;
- the repository still builds and passes its selected tests;
- the integration's important business workflow still completes;
- unsupported or ambiguous cases were disclosed instead of guessed.

Lumyn's differentiation is the combination of:

```text
authoritative change intent
+ repository-specific impact analysis
+ bounded migration
+ proof-honest workflow verification
+ two-party authorization
```

The durable value is a migration evidence chain, not access to a particular
model.

---

## Roles And Terminology

Use these terms consistently. Avoid the unqualified word `customer` where it is
unclear which party is meant.

### API Provider

The company that owns and sells the API or official SDK. It sponsors the
migration campaign and supplies authoritative change intent.

### API Consumer Organization

The organization whose application depends on the provider. It owns the
repository, credentials, test environment, and integration risk.

### Consumer Maintainer

The engineer authorized by the API Consumer Organization to review impact,
approve execution, inspect the patch, and merge or reject the draft PR.

### Provider Operator

The provider-side DX, SDK, API-platform, solutions-engineering, or customer
engineering person who creates and runs the campaign.

### Lumyn Operator

A member of the Lumyn team who assists with campaign setup, onboarding,
operations, or support. Lumyn Operator time is product delivery COGS and is
measured separately from Provider Operator and Consumer Maintainer time.

### Provider Change Packet

A versioned, signed, declarative artifact that describes a migration from a
source API/SDK version to a target version. The signature must verify against a
consumer-pinned API-provider trust root and a verified provider-to-package
ownership binding. A packet is provider input, not executable code and not
proof that a generated patch is correct.

### Provider Enrollment Bundle

A separately obtained, signed bootstrap artifact that binds the API Provider
organization, official package, root-key fingerprint, rotation policy, and
revocation/recovery endpoints. It also pins the provider status-signing key,
maximum acceptable status age, and either an exact status endpoint or the
contract for an offline signed status snapshot, plus the provider
receipt-acknowledgement signing key and permitted online or offline receipt
exchange classes. A campaign invitation narrows these pinned values to one
campaign and eligible-repository unit; it cannot introduce an untrusted key or
endpoint. The Consumer Maintainer obtains the bundle and expected fingerprint
through an authenticated provider admin console, named security contact, or
another independently verified channel—not from the campaign invitation alone.

### Migration Campaign

One provider change packet plus an invited cohort of API Consumer
Organizations, their independent authorization states, acknowledged minimal
connection units, and separately consented richer outcome attestations.

### Verified

`Verified` must always name its evidence boundary. The canonical labels are
`static_verified`, `repo_verified`, `workflow_contract_replay_passed`,
`workflow_verified_replay`, `workflow_verified_mock`, and
`workflow_verified_sandbox`. A workflow-verified label requires an approved
entrypoint executed from the exact patched repository head plus observed
interaction and outcome evidence in the named environment. Independent
contract or cassette replay cannot exceed `repo_verified`.

---

## Jobs To Be Done

### API Provider Job

> When we need customers off version X by date Y, help authorized customer
> teams identify and merge the required code changes so we can retire the old
> version without support spikes, outages, or an indefinite compatibility tail.

### API Consumer Job

> When an API provider requires a consequential migration, give me a minimal,
> explainable patch that fits my repository and prove the important integration
> still works, without giving the provider uncontrolled access to my code or
> production environment.

Both jobs are required. Provider willingness to pay without consumer
authorization does not create a functioning product.

---

## Initial Segment

The first design-partner segment is API-first B2B providers with:

- a consequential migration planned within six months;
- roughly 20 to 500 named, managed customer integrations;
- an official TypeScript/Node SDK published as an npm package;
- versioned OpenAPI, SDK releases, migration guidance, or equivalent source
  artifacts;
- a non-production sandbox with reliable read-back available for the first
  commercial pilot; each consumer may still choose exact-patched-head mock
  proof instead, and sandbox use requires its own grants;
- a DX, SDK, API-platform, solutions-engineering, or customer-engineering team
  that can recruit consumer maintainers;
- measurable support or legacy-version cost.
- either at least two consequential migrations or deprecations expected in the
  next 12 months, or at least 20 named managed integrations that could support
  an annual connected-repository program.

The first campaign should be a provider-managed cohort, not an anonymous
self-serve ecosystem.

Avoid initially:

- providers that cannot identify or contact integration owners;
- low-code or generated-client integrations with little application code;
- providers with no non-production sandbox or reliable outcome signal;
- migrations dominated by production-only data or irreversible actions;
- providers that cannot share prerelease intent with participating consumers;
- integrations that require broad production credentials;
- providers whose consequential changes are too rare to support recurring
  value.

### Development Entry Boundary

Lumyn does not need a contracted API Provider to begin M0-M4 engineering.
Publicly available API documentation, OpenAPI descriptions, published SDK
releases, migration guides, and license-compatible synthetic or historical
fixtures are sufficient to build and test the contracts, corpus, semantic
change intake, and read-only impact engine. Every input is pinned, provenance-
and license-checked, and treated as untrusted.

Those public inputs do not constitute a provider-authored signed packet,
provider endorsement, prerelease intent, sandbox authority, reachable customer
cohort, or evidence of demand. M5 migration-plan implementation and the
commercial pilot remain gated by M2.5's direct provider commitment,
supported-class signed canary, and frozen qualified cohort. Public artifacts
accelerate engineering; they do not counterfeit the sell-side relationship.

### Buyer And Champion

- Primary economic buyer: the VP or Head of Engineering accountable for the API
  platform and legacy compatibility/support cost. The startup substitute is
  the CTO.
- Primary champion/operator: Head of DX or SDK. API-platform, customer
  engineering, solutions engineering, support, and security are secondary
  stakeholders.
- Required second principal: Consumer Maintainer with authority over the target
  repository.

---

## Commercial Model

The initial offer is a paid, services-assisted migration campaign:

- The API Provider pays.
- API Consumer Organizations participate without a seat charge.
- The provider supplies a real migration, an accountable operator, a
  non-production sandbox, a supported-class canary packet, and a reachable
  cohort.
- Lumyn may assist with the first change packet and campaign operation.
- The consumer runs Lumyn locally or in its own CI environment by default.
- Before the pilot begins, the provider names the decision owner and criteria
  for converting to an annual connected-repository program or a second named
  migration campaign.

The intended paid continuation is an annual provider platform fee tied to
connected repositories, with separately priced active campaigns when useful.
A second one-off services purchase is insufficient unless it is a named step
into repeatable platform use. Per-consumer-seat pricing is explicitly avoided
because consumer participation should have minimal friction.

The consumer-controlled local runner and inspectable source are the trust and
adoption wedge. Lumyn does not describe the current repository or distribution
as open source until an explicit license, security policy, contribution policy,
support boundary, and release-integrity process exist. Design-partner
distribution uses explicit pilot terms plus a named security and support
contact. A managed provider coordinator may later handle invitations and
consented status, but it must not require provider access to consumer source
code.

---

## System Under Test

The system under test is the migration across three surfaces:

```text
provider change intent
+ consumer repository
+ verification environment
```

The API surface alone is not the complete system under test. A coding model is
not the primary system under test.

The result must preserve separate evidence about:

- what the provider says changed;
- what Lumyn found in the authorized repository;
- what Lumyn changed;
- what repository checks ran;
- what workflow environment ran;
- what outcome was observed;
- what remains unsupported, ambiguous, or unverified.

---

## Product Principles

### 1. Two Principals, Two Authorities

The provider is authoritative about intended API/SDK semantics. The consumer is
authoritative about repository access, execution, disclosure, and merge.
Neither authority implies the other.

### 2. Read Before Write

Impact analysis is read-only. The Consumer Maintainer sees an impact report and
migration plan before authorizing file changes.

### 3. Declarative Provider Input

Provider change packets may declare mappings, constraints, defaults, and
verification references. They may not carry arbitrary provider-supplied
scripts in the MVP.

### 4. Deterministic Before Model-Assisted

The MVP supports bounded deterministic transformations first. Model-assisted
patching is a future, separately approved execution mode. Any future model
output remains untrusted until the same evidence gates pass.

### 5. Proof Is Multidimensional

Impact coverage, patch provenance, repository validation, workflow evidence,
cleanup, boundaries, and residual risk remain separate axes. Lumyn never
compresses them into an unsupported green result.

### 6. Customer-Controlled Execution

Repository analysis, patching, tests, and credentials run in a
consumer-controlled local or CI environment by default. Raw code, diffs, logs,
traces, prompts, responses, and credentials are not shared with the provider
unless the consumer explicitly consents.

### 7. Fail Closed

Unknown wrappers, dynamic call construction, unsupported package managers,
ambiguous mappings, missing business values, stale packets, unsafe redaction,
or uncertain verification yield `needs_input`, `unsupported`, `uncertain`, or
`blocked`—never a speculative patch or false verified result.

### 8. Human Merge Authority

Lumyn opens draft PRs only. It never writes to the default branch and never
auto-merges during the MVP.

---

## MVP Product Flow

### 1. Provider Authors A Change Packet

The Provider Operator starts from a versioned Lumyn campaign kit. The kit
contains the packet template, canary manifest, signed-invitation contract,
consumer-connect instructions, consent/disclosure fields, and receipt
contracts. `lumyn campaign kit create` materializes it. `lumyn change publish`
validates and signs canonical packet bytes through a configured signer and
emits an immutable published packet. `lumyn campaign invite create` binds that
packet to the campaign, audience, expiry, confidentiality posture, and
consumer-visible authority request and emits the signed invitation. Signing
secrets are never embedded in the kit, packet, invitation, or repository.

The Provider Operator creates a packet containing:

- API-provider identity and packet issuer;
- verified provider-organization and official-package ownership binding;
- signing-key ID and algorithm;
- packet ID and lifecycle status;
- issue time, expiry, nonce, authorized audience, and confidentiality class;
- source and target API/SDK versions;
- immutable source references, tags, and digests;
- official npm package and supported version range;
- typed semantic changes;
- explicit operation, symbol, request-field, or response-field mappings;
- applicability conditions;
- deterministic derivation or default rules where applicable;
- unsupported and human-input-required conditions;
- safe file and transformation boundaries;
- migration deadline and compatibility window;
- sandbox or mock verification references;
- rollback and support guidance;
- key-rotation, revocation, and packet-withdrawal references;
- confidentiality posture for prerelease changes;
- a signature over the canonical packet bytes and immutable input manifest.

Packets use the lifecycle:

```text
draft -> published -> superseded | withdrawn
```

A published packet is immutable for its authorized audience; `published` does
not mean publicly disclosed. A correction creates a new packet version.
Expired, replayed, wrong-audience, withdrawn, superseded, revoked-key,
signature-invalid, or freshness-unverifiable packets cannot start or continue
mutation. Provider enrollment pins the consumer trust root and binds the
provider organization to the official SDK package. Initial enrollment requires
`lumyn provider enroll` with a provider enrollment bundle and an expected
fingerprint confirmed through an independent authenticated channel. An
invitation can reference a root or enrollment location but cannot authenticate
its own root. Registry/package provenance may corroborate the binding but is
not sufficient by itself. Normal rotation must be signed by the currently
pinned root. Emergency recovery requires explicit out-of-band re-enrollment,
freezes open campaigns, and invalidates prior unexecuted approvals. Rotation
and revocation state must be fresh enough for the packet policy before every
side effect. Freshness comes from either a signed, unexpired offline provider
status snapshot or an exact-endpoint read under a separate
`provider_trust_status_read` grant. The status response is signed by the pinned
provider status key, is bound to the provider, package, campaign, packet
digest, issue time, expiry, rotation epoch, revocation and withdrawal state,
and carries an anti-replay nonce. Online requests disclose no repository,
consumer, or migration evidence. If neither a current snapshot nor an
authorized status read is available, Lumyn blocks.

### 2. Provider Proves The Packet On Canary Fixtures

Before customer distribution, Lumyn validates the packet against provider-owned
or licensed canary repositories with:

- pinned source inputs;
- annotated affected call sites;
- expected patches;
- expected verification outcomes;
- negative and unsupported cases.

### 3. Provider Invites A Cohort

The provider identifies a bounded set of API Consumer Organizations and sends a
signed campaign invitation produced from the standard kit. The invitation
binds the provider, package, campaign, packet digest, non-authoritative
trust-enrollment instructions, expiry, confidentiality posture,
consumer-visible authority request, one eligible-repository unit nonce, the
enrollment-pinned provider receipt-acknowledgement signing-key reference, and
either an enrollment-approved receipt endpoint or offline receipt-exchange
contract. It contains no repository token, trusted root, or implicit grant. An
invitation grants no repository access, cannot self-authenticate its signing
key, and can acknowledge at most one opaque repository unit.

Before the Lumyn implementation worker collects, stores, or discloses
identifiable design-partner evidence, a narrow manual privacy/legal preflight
under Factory task-scoped `approval` must fix the allowed private fields,
participant consent, approved external private storage boundary, TTL,
expiry/revocation deletion, deletion-receipt and orphan ownership, and the
minimal connection-receipt and separately consented aggregate/hash-only public
fields. It identifies provider export and public commit as irreversible
disclosures. This preflight governs evidence handling only. It is not Lumyn
runtime product authority and cannot authorize repository, sandbox, GitHub, or
provider-attestation actions. The active M2.5 Factory approval grant must cite
both the preflight evidence and the canonical digest of this exact scope; a
generic approval or a digest from an earlier field/storage/TTL policy cannot
release evidence collection.

### 4. Consumer Enrolls The Provider And Authorizes A Repository

Before accepting an invitation from a new provider, the Consumer Maintainer
runs `lumyn provider enroll` with the separately obtained enrollment bundle and
out-of-band confirmed fingerprint. Enrollment stores the pinned root and
provider/package binding in the configured consumer-private state root outside
the repository. The supported activation runner accepts no repository
argument, starts from a neutral directory, and mounts only the enrollment input
and private state root; the checkout is unavailable to the process.

The Consumer Maintainer then runs `lumyn campaign accept` to verify the signed
invitation, provider/package binding, packet digest, expiry, and already pinned
trust root. Acceptance writes campaign state only to that external private root
and renders a reviewable, still-unapproved authorization-request template. It
uses the same activation isolation boundary and grants no capability by itself.

After reviewing the request, the Consumer Maintainer runs `lumyn authorization
issue` with an exact grant manifest and a consumer-controlled configured
signer. The command emits a signed, time-bounded private authorization bundle
outside the checkout. `lumyn authorization revoke` records a signed revocation
for one grant or the whole bundle. Issuance and revocation do not perform the
authorized side effect; every later action revalidates the current bundle at
its own boundary.

The grant manifest explicitly approves:

- repository and campaign;
- readable paths;
- writable paths;
- branch namespace;
- allowed commands;
- network destinations;
- sandbox credential scope;
- dependency-install posture;
- PR creation;
- independent `artifact_retention` and `artifact_deletion` grants;
- the minimal connection-receipt disclosure, if joining the sponsored program;
- exact richer status metadata, if any, that may be shared with the provider.

Authorization is per repository and per campaign, revocable, and time-bounded.
Tokens must expire. Production credentials and production mutations are out of
scope.

For a provider-sponsored program, the consumer may separately issue a minimal
signed campaign connection receipt. Before acknowledgement, the provider
authenticates the consumer organization through its existing customer-admin
channel or a documented out-of-band ownership check and signs a
`consumer-receipt-key-binding` over the invitation, opaque organization ID,
consumer receipt public-key fingerprint, verification method, verifier, and
expiry. The opt-in receipt contains only the campaign and packet digest,
eligible-repository unit and invitation nonces, opaque per-campaign
organization and repository IDs, connection event and time,
tool/schema versions, consent-policy and key-binding digests, issuer public-key
fingerprint, anti-replay nonce, audience, expiry, and signature. It contains no
repository name, source, diff, logs, test output, or credentials.

`lumyn campaign receipt submit` either sends the canonical receipt to the exact
pinned endpoint under the `campaign_receipt` grant or writes an offline export
bundle for provider import. The provider verifies the authenticated key
binding, signature, invitation, packet, audience, expiry, consent, and unit
cardinality, then returns a signed `campaign-receipt-ack` bound to the receipt
digest, provider acknowledgement key, decision, time, and deduplication key.
The consumer verifies and imports that acknowledgement. One invitation unit
can acknowledge only one opaque repository ID; an identical receipt digest is
idempotent, while a different repository or receipt for the same unit nonce
conflicts and cannot count. A provider counts and bills only a valid,
provider-acknowledged, non-replayed unit. A consumer that declines this
disclosure may still use Lumyn privately, but is not counted as a connected
repository in the sponsored program.

### 5. Lumyn Produces A Read-Only Impact Report

Lumyn detects:

- installed SDK package and version;
- package manifest and supported lockfile state;
- direct SDK imports;
- statically resolvable one-hop local wrappers;
- affected operations, symbols, request fields, and response fields;
- excluded generated or vendored code;
- ambiguous, dynamic, or unsupported usage.

Impact status is one of:

```text
unaffected
affected_supported
affected_needs_input
unsupported
uncertain
```

`unaffected` is allowed only when the analyzed scope and coverage are explicit.

### 6. Consumer Reviews The Migration Plan

The migration plan names:

- every planned file change;
- every affected and excluded call site;
- dependency and lockfile changes;
- transformation recipe and provenance;
- expected commands;
- workflow verification level;
- network and credential requirements;
- residual risks and unsupported cases.

No source file is changed before this plan is approved.

### 7. Lumyn Applies A Bounded Patch

The patch runs in an isolated worktree or equivalent disposable workspace.
Writes are restricted to approved paths and a declared diff budget. Identical
pinned inputs and the same deterministic recipe must produce the same patch.
Immediately before each file or lockfile write, Lumyn revalidates the current
packet bytes, digest, trust root, provider/package binding, lifecycle, audience,
expiry, rotation, revocation, withdrawal, supersession, and replay state. A
plan-time trust decision is never cached across the mutation boundary.

### 8. Lumyn Runs The Verification Ladder

Verification proceeds in order:

1. packet, plan, and patch integrity;
2. dependency and lockfile integrity;
3. repository baseline comparison;
4. compile and typecheck;
5. consumer-allowlisted tests;
6. independent contract/cassette replay, labeled
   `workflow_contract_replay_passed`;
7. deterministic replay or mock execution from an approved entrypoint at the
   exact patched repository head;
8. provider sandbox read-back from that patched head when separately approved;
9. cleanup, action-boundary, and redaction checks.

Before Lumyn is demonstrated to a provider as an end-to-end product, `lumyn
canary run --offline` must take the standard synthetic campaign kit through
signed invitation acceptance, explicit synthetic consumer-signed authorization
issuance, impact, plan, bounded patch, host-isolated deterministic replay or
mock verification, evidence rendering, and a local draft-PR preview. It uses no
live credential, sandbox, network, provider-reporting, or remote write
capability. Any unimplemented stage returns a typed nonzero result.

Repository tests are untrusted code. They run with network disabled by default.
Package lifecycle scripts are disabled by default and require separate consumer
approval. Sandbox credentials are isolated from general build and test
commands.

### 9. Lumyn Opens A Draft PR

PR creation is a separate explicit write action. The PR is bound to:

- provider change packet digest;
- repository base commit;
- generated head commit;
- migration plan digest;
- evidence artifact digests.

Draft-PR delivery depends on repository and deterministic workflow proof, not
on live sandbox availability. When separately authorized sandbox evidence
exists, Lumyn includes it without upgrading weaker evidence. Provider reporting
is a different optional action and never blocks PR creation.

The PR becomes stale when its base/head or packet changes.
Immediately before every repository, host-command, provider-status,
package-registry, sandbox, minimal campaign-receipt, richer
provider-attestation, remote-branch, PR, retention, or deletion side effect and
retry, Lumyn revalidates current packet trust and the exact applicable private
product grants. Factory worker `approval`, `credentials`, and `network` grants
govern implementation work only; Factory selection and dispatch do not
validate or confer Lumyn product authority.

### 10. Consumer Reviews And Merges

The Consumer Maintainer retains normal branch-protection, review, CI, and merge
controls.

### 11. Provider Receives Separately Consented Richer Campaign Status

After the distinct minimal connection-receipt flow, the consumer may
separately authorize a richer campaign attestation such as:

```text
invited -> authorized -> impact_found -> plan_approved
-> draft_pr_open -> merged | closed | blocked
```

Source code, raw diffs, logs, traces, prompts, responses, and credentials are
not provider-visible by default.

### 12. Independent Promotion Evidence

For tasks selected by risk and evidence policy, independent or human-operated
`holdout-evaluator`, `trace-grader`, and `evidence-attestor` gates run in that
order after code review and before commit/push. They emit passing artifacts
bound to the exact task and current validation work proof in the trusted
Factory evidence root. The implementation worker cannot self-grade or
self-attest, and shipping fails before commit or PR creation when required
independent evidence is missing, stale, malformed, self-authored, or
non-passing.

---

## Initial Technical Scope

### Supported Repository Shape

- GitHub repository.
- TypeScript source discoverable from a checked-in `tsconfig.json`.
- Direct dependency on one official npm SDK in `package.json`.
- One SDK major migration per campaign.
- Direct imports are patchable.
- Statically resolvable one-hop wrappers are detectable but may return
  `affected_needs_input`.
- `package-lock.json` is the first automatically writable lockfile. Automatic
  mutation requires the exact Node and npm versions, registry endpoint or
  immutable snapshot, package integrity inputs, toolchain digest, disabled
  lifecycle scripts, and a separately authorized registry-network capability.
- `pnpm-lock.yaml` and `yarn.lock` are detectable but impact-only until their
  update behavior is explicitly implemented and tested.
- Generated, vendored, minified, and build-output paths are excluded.
- Monorepos are supported only when one package root is explicitly selected.

### Supported Change Classes

The first deterministic patch engine supports exactly:

1. SDK method or API-operation rename with an explicit one-to-one mapping.
2. Request-property rename or relocation with an explicit mapping and no new
   business value.
3. Response-property rename or relocation with an explicit mapping and
   statically identifiable consumer access.

Dependency version and import updates needed by those changes are included.

A newly required field may be patched only after a later contract explicitly
defines a deterministic derivation or constant default. Otherwise it returns
`affected_needs_input`.

### Explicitly Unsupported In The MVP

- Authentication or authorization migrations.
- Webhook and event-semantics migrations.
- GraphQL, gRPC, JSON-RPC, and body-dispatched operation migration.
- Arbitrary business-logic changes.
- Dynamic reflection or runtime-generated method names.
- More than one selected package root.
- Generated-client regeneration.
- Production credentials or production mutations.
- Arbitrary provider-supplied scripts.
- Automatic merge.

---

## Product Status Axes

Lumyn must preserve these axes independently.

### Impact

- `not_analyzed`
- `unaffected`
- `affected_supported`
- `affected_needs_input`
- `unsupported`
- `uncertain`

### Patch

- `not_attempted`
- `planned`
- `generated`
- `failed`
- `stale`

Patch provenance is separately:

- `deterministic`
- `model_assisted`
- `manual`

Only `deterministic` is an MVP generation mode.

### Verification

- `not_run`
- `static_verified`
- `repo_verified`
- `workflow_contract_replay_passed`
- `workflow_verified_replay`
- `workflow_verified_mock`
- `workflow_verified_sandbox`
- `partial`
- `failed`
- `gap`
- `stale`

### Delivery

- `not_requested`
- `draft_pr_open`
- `merged`
- `closed`
- `blocked`
- `superseded`

No single roll-up status may hide a weaker axis.

---

## Evidence Contract

Every migration result must report:

- API-provider identity and change authority;
- campaign and packet IDs;
- source and target versions;
- packet digest and source provenance;
- repository identity without exposing it beyond consent;
- base and head commit;
- analyzed paths and explicit exclusions;
- affected, unaffected, unsupported, and uncertain call sites;
- impact precision/recall evidence when running a benchmark;
- migration recipe and patch provenance;
- changed files and dependency/lockfile delta;
- exact commands and results;
- pre-existing baseline failures;
- verification environment: static, repository, replay, mock, or sandbox;
- business-outcome/read-back evidence;
- cleanup, boundary, and redaction results;
- permissions, network destinations, and credential classes used;
- residual unsupported call sites and risks;
- artifact hashes and freshness;
- rollback guidance;
- reviewer checklist.

Evidence is invalidated when the packet, repository base/head, plan, patch, or
verification inputs change.

The existing proof-strength labels remain useful inside workflow evidence, but
they do not replace the migration axes above.

---

## Command Model

`lumyn` remains the primary local command surface.

### Required MVP Commands

| Command | Purpose |
|---|---|
| `lumyn init` | Initialize repo-local Lumyn configuration |
| `lumyn check` | Validate configured source paths and local prerequisites |
| `lumyn campaign kit create` | Create the standard provider campaign kit, canary manifest, invitation template, and receipt contracts |
| `lumyn change publish --packet <draft> --signer <signer-ref>` | Validate, sign through a configured signer, and publish immutable packet bytes without exporting signing secrets |
| `lumyn campaign invite create --campaign <config> --packet <published> --signer <signer-ref>` | Bind the published packet to a campaign/audience and emit a signed, expiring invitation with a no-authority request |
| `lumyn provider enroll --bundle <bundle> --fingerprint <expected>` | Pin a provider root and provider/package binding from a separately authenticated bootstrap bundle; never trust invitation-supplied root material |
| `lumyn campaign accept --invitation <invitation>` | Verify the invitation against an already pinned root, initialize campaign state outside the repository, and render a no-authority authorization request |
| `lumyn trust refresh --provider <provider-id> [--snapshot <signed-snapshot> \| --online]` | Import a signed offline status snapshot or, under an exact `provider_trust_status_read` grant, fetch one from the pinned endpoint without disclosing consumer or repository data |
| `lumyn authorization issue --request <request> --grant <grant-manifest> --signer <consumer-signer-ref>` | Explicitly issue a signed, private, time-bounded product-authority bundle after review; issuance performs no repository, command, network, sandbox, or GitHub action |
| `lumyn authorization revoke --authorization <id> --reason <reason> --signer <consumer-signer-ref>` | Record a signed grant or bundle revocation in consumer-private state |
| `lumyn change validate <packet>` | Validate packet schema, provenance, lifecycle, and pinned inputs |
| `lumyn impact --change <packet> --repo <path>` | Produce a read-only impact report |
| `lumyn authorization validate --bundle <private-ref> --task <task>` | Validate an exact private product-authority bundle for diagnostics and closure evidence |
| `lumyn campaign receipt issue --authorization <id> --binding <provider-signed-key-binding> --consent <policy> --signer <consumer-signer-ref> --out <receipt>` | Optionally join a provider-sponsored program by emitting the minimal signed connection receipt; private local use does not require this disclosure |
| `lumyn campaign receipt submit --receipt <receipt> [--online \| --out <offline-bundle>]` | Under the exact `campaign_receipt` grant, submit only the canonical receipt to the invitation-pinned endpoint or create the bounded offline exchange bundle |
| `lumyn campaign receipt acknowledge --receipt <receipt> --binding <verified-binding> --signer <provider-ack-signer> --out <ack>` | Provider-side verification and idempotent acknowledgement bound to authenticated consumer signer, invitation unit, receipt digest, and deduplication key |
| `lumyn campaign receipt ack import --ack <provider-signed-ack>` | Verify the provider acknowledgement against the invitation-pinned acknowledgement key and record the sponsored-program unit locally |
| `lumyn migrate plan --impact <report>` | Produce a reviewable no-write migration plan |
| `lumyn migrate apply --plan <plan>` | Apply an approved plan in an isolated workspace |
| `lumyn verify --migration <result>` | Run the declared verification ladder |
| `lumyn canary run --kit <kit> --repo <synthetic-repo> --offline --pr-preview <path>` | Run the receipt-backed synthetic end-to-end activation and migration canary without live credentials, network, sandbox, reporting, or remote writes |
| `lumyn trace <evidence>` | Render local evidence without changing state |
| `lumyn artifacts gc [--dry-run]` | Retry or preview TTL/revocation deletion and report unresolved private-artifact orphans |
| `lumyn pr create --draft <result>` | Open an explicitly authorized draft PR |

The command grammar may be implemented incrementally, but unimplemented
commands must return a typed nonzero result. They must never return a generic
successful command envelope.

### Deferred Commands

Generic `record`, `eval`, `demo`, and `share` product surfaces are not required
for the migration MVP. The narrowly specified `canary run` command is required
and must not route through the generic unimplemented `demo` surface. Workflow
recording may return later as a verification authoring aid. Model-provider
evaluation is not part of the customer-facing MVP.

### Stable Exit Codes

- `0`: success
- `1`: general or internal error
- `2`: invalid usage, invalid input, parse error, or local configuration error
- `3`: source or provider-change completeness failure in strict mode
- `4`: provider-change, authorization, migration-plan, or workflow-contract
  validation failure
- `5`: migration application or verification failure
- `6`: reserved for the pre-2.0 live-eval threshold contract; not emitted by
  migration MVP commands
- `7`: authorization, credential, auth, or environment error
- `8`: dependency, API provider, GitHub, sandbox, or network error
- `9`: packet, impact, plan, patch, trace, cassette, evidence, or replay
  integrity failure

### Command Result Compatibility

The existing `lumyn.command_result` version `1.0` remains supported while the
migration result schemas are introduced. Existing `provider_metadata` means
model-provider metadata and must not be reused for the API Provider. New
contracts use explicit names such as:

- `api_provider_id`
- `change_authority`
- `model_provider_metadata`

Migration commands must not set `eval_mode: surface_only`. During the
compatibility window, non-applicable legacy fields remain explicitly
`not_applicable` until a versioned result-schema migration removes them.

---

## Artifact Model

### Retained Verification Artifacts

The retained executable schema inventory is:

- `workflow-contract.schema.json`
- `expected-outcome.schema.json`
- `validator.schema.json`
- `action-boundary.schema.json`
- `human-annotation.schema.json`
- `required-context.schema.json`
- `state-binding.schema.json`
- `canonical-trace.schema.json`
- `evidence-event.schema.json`
- `cassette.schema.json`
- `result-axes.schema.json`
- `proof-strength.schema.json`
- `command-result.schema.json`
- `redaction-config.schema.json`

### Planned Migration Artifacts

- `campaign-kit`
- `provider-enrollment`
- `provider-status-snapshot`
- `provider-change`
- `campaign-invitation`
- `authorization-request`
- `repository-authorization`
- `authorization-revocation`
- `consumer-receipt-key-binding`
- `campaign-connection-receipt`
- `campaign-receipt-ack`
- `impact-report`
- `migration-plan`
- `patch-manifest`
- `migration-verification`
- `migration-pr-result`
- `campaign-attestation`
- `offline-canary-receipt`
- `remediation-outcome`

These planned artifacts do not become executable contracts until their schema,
valid/invalid fixtures, compatibility posture, and validator tests ship in the
same task.

### Product Artifact Layout

```text
consumer checkout (committable)
  lumyn.yaml
  examples/                 # synthetic/public, licensed fixtures only
  examples/holdout-manifest.json  # opaque IDs and digests; no held-out inputs or answers

consumer-private runtime root (outside the checkout)
  providers/<api-provider>/enrollment.json
  providers/<api-provider>/status-snapshots/<snapshot-id>.json
  campaigns/<campaign-id>/acceptance.json
  campaigns/<campaign-id>/consumer-receipt-key-bindings/<binding-id>.json
  campaigns/<campaign-id>/connection-receipts/<receipt-id>.json
  campaigns/<campaign-id>/connection-receipt-acks/<receipt-id>.json
  changes/<api-provider>/<change-id>.json
  authorizations/<campaign-id>/<repository-id>.json
  revocations/<authorization-id>.json
  impacts/<run-id>.json
  migrations/<run-id>/
    plan.json
    patch-manifest.json
    verification.json
    pr-result.json
  workflows/<workflow-id>.yaml
  cassettes/<workflow-id>.json
  runs/<run-id>/
    trace.json
    report.html
  retention/
    deletion-receipts/<artifact-id>.json
    orphans/<artifact-id>.json

evaluator-controlled holdout root (outside checkout and task-executor access)
  suites/<suite-id>/
    inputs/
    answer-key/
    expected-patches/

public pilot evidence
  .factory/artifacts/pilot/lumyn-migration-mvp/public/
    aggregate.json
    evidence-hashes.json
```

The private runtime root is configured explicitly, cannot resolve inside the
consumer checkout or any public source repository, is excluded from source
control, and enforces authorization TTL plus deletion on expiry or revocation
at creation, read, process startup, and the next run. Deletion is retry-safe
and produces a receipt or durable orphan record. `lumyn artifacts gc` is the
operator recovery surface for retrying partial cleanup and inspecting orphans;
it cannot extend TTL, revive revoked authority, or rewrite historical closure.
Repository-root ignore rules are defense in depth for a misconfigured or legacy
runtime, not permission to store private artifacts there. Only synthetic or
licensed fixtures may be committed under `examples/`. Provider-visible
attestations and public pilot evidence are separate disclosure products:
provider attestations may contain only the exact fields authorized by that
consumer, while public pilot evidence may contain only separately consented,
redacted aggregates or hashes. Neither may contain unconsented private
evidence. Export to a provider and commit to public history are irreversible
disclosure boundaries: revocation stops future export and deletes only
Lumyn-controlled private copies. It cannot promise recall from a recipient,
clone, cache, or Git history.

---

## Trust And Authorization Contract

### Provider Change Channel

The provider change channel is a software supply-chain boundary.

- Packets must identify their issuer.
- Published packets must be immutable and digest-bound.
- The signature must verify against a consumer-pinned provider trust root,
  including verified organization/package ownership, issuer key, issue time,
  audience, expiry, rotation, revocation, withdrawal, and replay checks.
- The consumer must be able to inspect the packet before authorization.
- Packet mappings cannot execute arbitrary code.
- A withdrawn or superseded packet blocks new mutation.
- Packet trust is revalidated immediately before every local, sandbox, or
  remote side effect and retry; no prior validation decision is reusable.
- Revalidation consumes either a signed offline status snapshot within the
  enrollment policy's maximum age or a freshly fetched snapshot under an exact
  `provider_trust_status_read` grant. The grant binds one provider, package,
  status endpoint, request shape, response budget, and expiry; it cannot carry
  consumer or repository data. Missing, stale, unsigned, replayed, or
  wrong-endpoint status blocks.
- Packet provenance is necessary context, not correctness proof.

### Repository Access

- Impact analysis begins read-only.
- Read and write permissions are separate.
- Allowed paths and excluded paths are explicit.
- Patch execution occurs outside the default working tree.
- Default-branch writes are prohibited.
- Remote branch write and PR write are separate permissions.
- Installation tokens are scoped and time-bounded.
- Authorization can be revoked.

### Command Execution

- Commands are allowlisted and displayed before execution.
- Repository commands run only through a supported fail-closed isolation
  backend. The grant binds exact read-only and writable mounts, a neutral home
  and temp root, executable/toolchain roots, sanitized environment classes,
  inherited file descriptors, local socket policy, process-tree limits, and OS
  credential access.
- The default isolation profile exposes no host home, SSH/GPG/cloud credential
  stores, keychain, agent or Docker sockets, unrelated local services, or
  inherited descriptors beyond standard streams. Child processes inherit the
  same restrictions. If the backend cannot enforce them, the command does not
  run.
- Repository tests run without network by default.
- Dependency lifecycle scripts are disabled by default.
- Registry network, sandbox network, and lifecycle-script enablement require
  separate consumer approval.
- Lockfile mutation is bound to exact Node/npm versions, registry/snapshot,
  package-integrity inputs, and the toolchain digest.
- Pre-existing failing checks are captured before mutation.
- Commands must have time and output budgets.

### Credentials And Network

- Production credentials are prohibited.
- Sandbox credentials are limited to sandbox verification.
- Build and test commands do not receive sandbox credentials.
- Network destinations are allowlisted.
- Sandbox writes use namespaces, request budgets, cleanup, and orphan reporting.
- Sandbox access separately authorizes transmitted payload classes. It uses
  synthetic or approved non-sensitive test data; production customer data,
  PII, credentials, and secrets are prohibited.
- The grant records what the provider can log, the retention period, deletion
  terms, and the consumer-approved payload disclosure.

### Data Sharing

The provider receives no raw repository data by default.

A provider-sponsored connection meter is a separate opt-in disclosure. Its
minimal signed receipt is schema-validated, audience/expiry/nonce checked,
bound to the invitation, packet, opaque organization/repository IDs, consent
policy, provider-authenticated consumer key binding, and consumer-controlled
issuer key. The provider-signed acknowledgement binds the receipt digest,
invitation unit, verified key-binding digest, decision, time, and
deduplication key. One invitation unit acknowledges at most one opaque
repository ID; same-digest retry is idempotent and conflicting reuse blocks.
Only acknowledged valid receipts count as connected repositories. Declining
the receipt leaves private local use available but excludes that repository
from sponsored-program counting and billing.

Consumer consent is required before sharing:

- repository identity;
- affected-file counts;
- patch or PR state;
- validation status;
- workflow status;
- failure reason;
- merge outcome.

Raw source, diffs, logs, traces, prompts, responses, credentials, and secret
values are excluded from default attestations.

Private runtime and pilot records expire according to the authorization and are
deleted on revocation automatically at creation, read, process startup, and the
next run. A deletion receipt or orphan report is required, and operator
recovery cannot rewrite historical closure.
Provider visibility does not imply public disclosure. Public evidence requires
separate consent and is aggregate/hash-only. Consent for provider export or
public commit must disclose that the copy is not recallable. Revocation blocks
future disclosure and deletes Lumyn-controlled private copies; it cannot erase
a provider recipient's records, public Git history, clones, or caches.

---

## Functional Requirements

### FR1: Provider Change Contract

Lumyn must validate versioned, immutable, declarative provider change packets
with provenance, lifecycle, semantic mappings, applicability, unsupported
conditions, verification references, and rollback guidance, and must revalidate
their current trust state at every side-effect boundary.

### FR2: Two-Party Authorization

Lumyn must keep provider campaign authority separate from consumer repository,
execution, credential, disclosure, and merge authority. Lumyn product code,
not Factory dispatch, enforces the private product bundle immediately before
each live action.

### FR3: Read-Only Impact

Lumyn must produce an impact report without modifying source files, dependency
files, lockfiles, branches, or PRs.

### FR4: TypeScript SDK Analysis

Lumyn must detect the selected npm SDK/version, direct imports, affected call
sites, one-hop wrapper uncertainty, generated/vendored exclusions, and
unsupported repository shapes.

### FR5: Explicit Impact Coverage

Lumyn must report analyzed scope, exclusions, uncertainty, and residual
unsupported call sites. It must not infer global downstream coverage from
participating repositories.

### FR6: Reviewable Migration Plan

Lumyn must create a no-write plan that identifies every proposed file,
transformation, command, permission, network destination, verification level,
and residual risk.

### FR7: Bounded Deterministic Patch

Lumyn must implement only declared supported transformations, remain within
approved files and diff budgets, and produce reproducible output for identical
pinned inputs. It must revalidate packet trust immediately before each write.

### FR8: Repository Baseline

Lumyn must capture pre-existing dependency, compile, typecheck, and selected
test status before applying the patch so existing failures are not attributed
to the migration.

### FR9: Verification Ladder

Lumyn must execute verification in declared order and record skipped, blocked,
failed, stale, and passed stages separately.

### FR10: Workflow Outcome Proof

Lumyn must reuse schema-backed workflow, replay, evidence, boundary, cleanup,
redaction, and read-back semantics to distinguish static, repository, replay,
mock, and sandbox proof. A workflow-verified label requires causal execution
from the exact patched repository head; independent contract replay uses
`workflow_contract_replay_passed` and cannot exceed `repo_verified`.

### FR11: Evidence-Bound Result

Lumyn must bind impact, patch, and verification evidence to the packet digest,
base commit, head commit, plan digest, and artifact hashes.

### FR12: Draft PR Delivery

Lumyn must require explicit PR-write authorization, create a draft PR only,
include the required evidence and reviewer checklist, and behave idempotently
for the same campaign/repository/base state.

### FR13: Consented Connection Receipt And Richer Campaign Attestation

Lumyn must distinguish the minimal signed campaign connection receipt from
optional richer provider reporting. A consumer joins provider-sponsored
counting only by explicitly issuing the receipt; private local use and
draft-PR delivery remain available without it. The provider must verify the
receipt's provider-authenticated consumer signer binding, invitation and packet
binding, audience, expiry, nonce, signature, and one-unit cardinality before
acknowledging or billing it. The acknowledgement is provider-signed, bound to
the receipt digest and deduplication key, and portable through a pinned endpoint
or bounded offline export/import exchange. Any richer campaign attestation
requires a separate field allowlist and may return only those approved fields.

### FR14: Outcome Feedback

Lumyn must record merged, closed, reverted, and substantively corrected
outcomes with provenance. Outcome feedback may inform future remediation
recipes, but cannot silently change an active packet or patch rule.

### FR15: Stable Machine Interface

Every state-returning command must support stable JSON, typed errors,
non-interactive execution, and the stable exit-code contract. Factory and
committable public artifact references are repo-relative. Consumer-private
runtime references are relative to the explicitly configured private state
root and appear in public or provider-visible records only as opaque IDs and
digests.

### FR16: Private-Artifact Lifecycle

Lumyn must automatically enforce private-artifact TTL and
expiry/revocation deletion at creation, read, process startup, and the next
run; emit deletion receipts or durable orphan reports; retry safely after
partial deletion or crash; expose `lumyn artifacts gc` for operator recovery;
and never rewrite historical closure to hide deleted evidence.

### FR17: Two-Sided Activation

Lumyn must provide a standard provider campaign kit, configured-signing packet
publication, a signed invitation, out-of-band-authenticated provider
enrollment, consumer campaign acceptance plus an authorization-request receipt,
explicit signed authorization issuance and revocation, an optional signed
campaign connection receipt with provider-authenticated signer binding and a
provider-signed acknowledgement exchange, and a synthetic offline end-to-end
canary.
Invitations cannot authenticate their own root. Enrollment and acceptance run
without the checkout mounted; no activation request implies repository or
product authority, and no provider demonstration may use a false-green or
partially implemented command path.

---

## Non-Functional Requirements

### NFR1: Fail-Closed Honesty

Unsupported, ambiguous, stale, unauthorized, unredactable, or unverified states
must never report a successful migration or verified workflow.

### NFR2: Determinism

Pinned public fixtures and deterministic patch recipes must produce stable
impact, patch, and evidence results on repeated runs.

### NFR3: Consumer Privacy

Consumer code and private evidence remain in the consumer-controlled execution
plane, outside the checkout and public source repository, by default. Runtime
artifacts enforce TTL and deletion on expiry or revocation at every lifecycle
entry point; only explicitly consented redacted aggregates or hashes may be
public. Provider export and public commit are disclosed as irreversible:
revocation stops future sharing and removes Lumyn-controlled private copies,
but cannot promise recall from recipients, Git history, clones, or caches.

### NFR4: Least Privilege

Read, write, command, provider-status read, registry, credential, PR, campaign
receipt, richer reporting, retention, and deletion permissions are
independently bounded and evidenced.

### NFR5: Redaction Before Persistence

Secrets must be redacted before persistence or sharing. Redaction uncertainty
blocks the affected artifact.

### NFR6: Explainability

Every changed file and call site must map to a provider change entry and
migration recipe. Unsupported cases must include a concrete reason and next
action.

### NFR7: Artifact Stability

Persisted artifacts are versioned, schema-backed, digest-bound, and
migration-aware. Factory and committable public artifacts are repo-relative;
consumer-private artifacts are private-root-relative and cross a disclosure
boundary only through opaque IDs and digests.

### NFR8: Recovery And Idempotency

Interrupted or repeated runs must not duplicate branches, draft PRs, sandbox
resources, connection-receipt units or acknowledgements, richer campaign
attestations, deletion receipts, or orphan records.

### NFR9: Bounded Performance

On the fixed MVP benchmark, median read-only impact analysis should complete in
under five minutes and median draft-PR preparation in under twenty minutes,
excluding repository-defined test duration.

### NFR10: No Production Dependency

The deterministic benchmark and default consumer run require no production
credentials, production traffic, or production mutation.

### NFR11: Host Isolation

Repository-defined commands run only when a supported isolation backend can
enforce exact mounts, sanitized environment, closed local-socket and inherited
descriptor policy, process-tree inheritance, resource budgets, and denial of
host home, credential stores, agents, keychains, and container-control sockets.
An unavailable or unverifiable isolation backend blocks execution.

---

## Failure Taxonomy

Required failure classes include:

- `packet_invalid`
- `packet_untrusted`
- `packet_superseded`
- `packet_withdrawn`
- `source_target_mismatch`
- `spec_sdk_mapping_incomplete`
- `authorization_missing`
- `authorization_expired`
- `authorization_revoked`
- `provider_status_stale`
- `provider_status_read_not_allowed`
- `read_scope_exceeded`
- `write_scope_exceeded`
- `command_not_allowed`
- `command_isolation_unavailable`
- `command_isolation_violation`
- `network_not_allowed`
- `credential_scope_invalid`
- `unsupported_repository_shape`
- `unsupported_package_manager`
- `multiple_sdk_versions`
- `dynamic_usage_uncertain`
- `wrapper_usage_needs_input`
- `generated_or_vendored_code`
- `impact_uncertain`
- `required_business_value_missing`
- `unsupported_change_class`
- `patch_conflict`
- `patch_stale`
- `diff_budget_exceeded`
- `baseline_already_failing`
- `dependency_integrity_failed`
- `compile_failed`
- `typecheck_failed`
- `tests_failed`
- `tests_flaky`
- `replay_failed`
- `sandbox_unavailable`
- `sandbox_drift`
- `workflow_proof_gap`
- `cleanup_failed`
- `orphan_detected`
- `redaction_uncertain`
- `evidence_stale`
- `duplicate_pr`
- `campaign_receipt_invalid`
- `campaign_receipt_signer_untrusted`
- `campaign_receipt_unit_conflict`
- `campaign_receipt_ack_invalid`
- `campaign_receipt_not_consented`
- `provider_attestation_not_consented`

Failure summaries must cite concrete artifact references. Unsupported model
diagnosis must not be presented as fact.

---

## Acceptance Tests

The active acceptance ledger enumerates exactly 62 item-level closure units.
Group headings, epics, waves, and an overall MVP label cannot substitute for
item-level evidence.

### Retained Foundation

1. `BASE-001`: The Go CLI, config discovery, stable result envelope, exit-code
   constants, and local validation baseline remain functional.
2. `BASE-002`: Existing workflow, evidence, cassette, trace, proof, boundary,
   redaction, and command-result schemas remain executable and versioned.
3. `BASE-003`: CI, coverage, CodeQL, CODEOWNERS, required checks, passive review,
   shipping, and post-merge governance remain enforced.
4. `BASE-004`: OpenAPI/docs intake produces structured refs, fingerprints,
   deprecation findings, and concrete source locations.
5. `BASE-005`: The dev/architecture planning policies and 12-level test matrix
   remain propagated into every runner-ready migration task.

### Product Rebaseline

1. `REB-001`: Unimplemented commands return a typed nonzero result instead of a
   generic pass.
2. `REB-002`: API-provider identity and model-provider metadata are
   unambiguous in schemas, CLI results, docs, and fixtures.
3. `REB-003`: Historical task evidence remains immutable and the old plan is
   explicitly non-active.

### Provider Change And Benchmark

1. `CHG-001`: A valid packet includes pinned source/target artifacts, typed
   semantic mappings, an immutable digest, and declared audience and
   confidentiality, and passes schema and semantic validation.
2. `CHG-002`: Draft, superseded, withdrawn, expired, wrong-audience,
   revoked-key, replayed, mutated-after-publish, or provenance-mismatched
   packets block mutation, including when the state changes after planning but
   before a local, sandbox, or remote side effect.
3. `CHG-003`: A provider packet cannot execute arbitrary scripts.
4. `CHG-004`: A published packet is trusted only when its signature verifies
   against a consumer-pinned API-provider trust root and its issuer-to-package
   binding, key, timestamp, audience, expiry, rotation, revocation, and replay
   checks pass immediately before every side effect and retry; no prior trust
   result is cached. Current status comes from either a signed offline snapshot
   within the pinned maximum age or an exact-endpoint read under
   `provider_trust_status_read`; missing, stale, replayed, unsigned, or
   undeclared status access blocks. Published means immutable for its
   authorized audience, not publicly disclosed.
5. `CORPUS-001`: At least three provider-change packets, one per supported
   change class, and at least nine visible development consumer fixtures,
   three per class, contain frozen ground truth. Each class has at least 20
   annotated affected sites and 20 annotated non-affected candidate sites. An
   independent holdout owner provisions and freezes a separate
   evaluator-controlled holdout containing at least one repository per class.
   Only opaque case IDs, non-resolving provenance class and license posture, a
   frozen suite commitment, and encrypted or HMAC artifact commitments are
   committed. Source URLs, repository or package identifiers, plaintext
   content digests, held-out inputs, and answer keys are never committed,
   mounted, passed, or exposed to the task executor.
6. `CORPUS-002`: Negative development fixtures cover unsupported, ambiguous,
   stale, out-of-boundary, semantic-non-equivalence, and false-verification
   cases. The evaluator-controlled held-out workflow suite has at least two
   positive and two negative patched-repository executions for every
   deterministic verification environment label emitted by the MVP, and
   returns aggregate results and digests without answer material.
7. `CORPUS-003`: Public fixtures are described as engineering benchmarks, not
   evidence of provider demand.

### Two-Sided Activation

1. `ACT-001`: A Provider Operator can create a standard campaign kit, author
   and validate a declarative change packet, sign and publish immutable packet
   bytes through a configured signer, and then create a signed campaign
   invitation without embedding signing secrets or granting repository access.
2. `ACT-002`: A Consumer Maintainer can enroll the provider from a separately
   authenticated enrollment bundle and out-of-band confirmed fingerprint,
   pin the provider root and package binding, and then verify and accept the
   signed invitation. Enrollment and acceptance store state only in the
   consumer-private root with the checkout unavailable and produce a reviewable
   authorization request. They grant no authority. The maintainer can then
   explicitly issue and revoke a signed, time-bounded private authorization
   bundle and, independently, opt into a signed minimal campaign connection
   receipt without bespoke artifact authoring or repository disclosure.
3. `ACT-003`: A synthetic offline canary runs the published packet and
   invitation through consumer acceptance, explicit authorization issuance,
   impact, plan, bounded patch, host-isolated replay or mock verification,
   evidence rendering, and a local draft-PR preview with receipts, while using
   no live credentials, sandbox, provider export, network, or remote write.
   Any missing implementation stage returns a typed nonzero result, so provider
   demonstrations cannot false-green.

### Authorization And Impact

1. `AUTH-001`: Repository analysis requires per-campaign, per-repository,
   time-bounded read authorization.
2. `AUTH-002`: Local write, host-isolated command execution, provider trust
   status read, registry network, sandbox network, credential, remote branch
   write, PR write, minimal campaign receipt, richer provider attestation,
   retention, and deletion are separate, time-bounded scopes. Issuance and
   revocation are signed consumer actions; no grant implies another.
3. `AUTH-003`: The API Provider and the public cannot retrieve raw consumer
   code, private evidence, or identifiable pilot records through the default
   campaign interface.
4. `AUTH-004`: Private runtime and pilot artifacts live outside the consumer
   checkout and public source repository by default, are non-committable,
   and automatically enforce TTL plus deletion on revocation at creation, read,
   process startup, and the next run, with a deletion receipt or durable orphan
   report and retry-safe operator recovery that does not rewrite historical
   closure. Provider attestations expose only exact consumer-consented fields;
   public evidence exposes only separately consented redacted aggregates or
   hashes. Consent identifies provider export and public commit as irreversible;
   revocation stops future disclosure and deletes Lumyn-controlled private
   copies but cannot promise recall from recipients or public history.
5. `AUTH-005`: Sandbox verification separately authorizes transmitted payload
   classes, uses synthetic or approved non-sensitive test data, forbids
   production customer data and secrets, and records provider logging,
   retention, and deletion terms.
6. `IMP-001`: Lumyn detects the selected official npm SDK and current version.
7. `IMP-002`: Lumyn identifies all annotated directly imported affected call
   sites in the fixed supported corpus.
8. `IMP-003`: Lumyn detects one-hop wrappers but does not patch uncertain
   wrapper usage without input.
9. `IMP-004`: Generated, vendored, dynamic, multi-root, multi-version, and
   unsupported lockfile cases receive explicit exclusion or failure status.
10. `IMP-005`: On held-out fixtures, every supported change class has `100%`
    recall, at least `95%` precision, and zero false `unaffected` results.
    `uncertain`, `unsupported`, or missed on a known-supported affected site
    counts as a false negative. Per-class and aggregate confusion matrices are
    reported.

### Plan And Patch

1. `PLAN-001`: Impact analysis performs no repository mutation.
2. `PLAN-002`: The Consumer Maintainer can inspect every proposed file,
   transformation, command, permission, verification level, and residual risk
   before write approval.
3. `PATCH-001`: A method or operation rename is patchable only when the packet
   declares a one-to-one mapping with unchanged call signature, request and
   response meaning, side effects, and error semantics.
4. `PATCH-002`: A request-property rename or relocation is patchable only when
   value, type, optionality, cardinality, units, and business meaning are
   unchanged and no new value is inferred.
5. `PATCH-003`: A response-property rename or relocation modifies only
   statically identified accesses and is patchable only when value, type,
   nullability, cardinality, units, and business meaning are unchanged.
6. `PATCH-004`: Identical packet, base commit, toolchain digest,
   package-integrity inputs, and registry snapshot produce the same
   deterministic patch.
7. `PATCH-005`: Patches remain inside approved paths and diff budgets with no
   unrelated edits.
8. `PATCH-006`: Any semantic-equivalence precondition failure, missing business
   value, unsupported change, or ambiguous mapping returns
   `affected_needs_input` or `unsupported` without edits.

### Verification And Evidence

1. `VER-001`: Lumyn records pre-existing repository failures before mutation.
2. `VER-002`: Supported patches pass dependency integrity, compile/typecheck,
   and declared allowlisted tests through a supported fail-closed host-isolation
   backend or report the exact failed stage. Tests cannot access undeclared
   mounts, host home or credential stores, agent/container sockets, inherited
   descriptors, OS credentials, or network, and child processes inherit the
   same boundary.
3. `VER-003`: Independent contract or cassette replay, patched-repository
   replay, patched-repository mock execution, and live sandbox execution use
   separate canonical labels.
4. `VER-004`: Sandbox verification requires approved endpoint and payload
   allowlists, synthetic or approved non-sensitive data, a non-production
   credential and namespace, provider logging and retention terms, budgets,
   cleanup, and an orphan report.
5. `VER-005`: A final-state pass with a missing patched-head execution or any
   boundary, cleanup, redaction, freshness, or evidence-integrity failure
   cannot return a `workflow_verified_replay`, `workflow_verified_mock`, or
   `workflow_verified_sandbox` label.
6. `VER-006`: On the frozen held-out workflow suite, every positive and
   negative case receives its exact expected label, with zero false-positive
   and zero false-negative workflow-verification results.
7. `VER-007`: A `workflow_verified_replay`, `workflow_verified_mock`, or
   `workflow_verified_sandbox` result requires an approved entrypoint from the
   exact patched repository head and observed interaction and outcome evidence
   in that environment. Independent contract or cassette replay is
   `workflow_contract_replay_passed` and cannot exceed `repo_verified`.
8. `EVD-001`: Results preserve separate impact, patch, verification, delivery,
   permission, and residual-risk axes, and `lumyn trace` renders those bound
   axes locally without network access or provider disclosure.
9. `EVD-002`: Evidence is bound to packet digest, base/head commits, plan
   digest, artifact hashes, and verification environment.
10. `EVD-003`: A packet or repository-head change invalidates stale evidence.

### Draft PR Delivery

1. `PR-001`: PR creation requires a separate explicit authorization after the
   migration plan is reviewed.
2. `PR-002`: Lumyn writes only a draft PR on an authorized non-default branch.
3. `PR-003`: The PR includes provider packet provenance, impact scope, changed
   files, dependency delta, generation method, exact commands and results,
   baseline failures, workflow environment, permissions used, unsupported
   cases, residual risk, rollback guidance, and a reviewer checklist.
4. `PR-004`: Repeated delivery for the same campaign/repository/base state is
   idempotent and does not create duplicate PRs.
5. `PR-005`: A provider-sponsored connection counts only after the provider
   signs an acknowledgement for a valid consumer-signed, invitation-bound,
   non-replayed minimal connection receipt whose signer is authenticated to the
   invited organization and whose unit nonce maps to one opaque repository.
   Richer provider-visible campaign status remains separately authorized and
   contains only consumer-consented fields.

### Design-Partner Qualification

1. `DISC-001`: Before migration-plan implementation proceeds, one qualified API
   Provider supplies written commitment for a consequential migration within
   six months, a named economic buyer and operator, authoritative source and
   target artifacts, prerelease-sharing authority, a non-production sandbox,
   a signed canary packet proving at least one supported change class, an
   operational signed provider-status channel with its key and maximum-age
   policy, and a provider receipt-acknowledgement key and endpoint or offline
   exchange contract, plus an agreed paid-pilot price and decision process.
   The provider also documents recurring-value potential through at least two
   consequential migrations or deprecations expected within 12 months, or at
   least 20 named managed integrations eligible for an annual
   connected-repository program, and names the post-pilot annual-platform or
   second-campaign purchase decision.
2. `DISC-002`: Before migration-plan implementation proceeds, a frozen candidate
   cohort contains at least five distinct eligible repository IDs across at
   least three Consumer Organizations, with one accountable maintainer per
   repository. Every repository matches the supported GitHub and
   TypeScript/Node shape and is prequalified as plausibly using an affected
   source version, operation, or field. One funnel unit is counted per
   repository. The invitation window, absolute campaign judgment deadline,
   baseline, one primary provider outcome and its material threshold,
   correction rubric, payment and contribution evidence, confidentiality,
   retention, and consented-data protocol are frozen before outcomes are
   observed. Before identifiable evidence is collected, stored, or disclosed,
   a manual privacy/legal preflight fixes allowed private fields, participant
   consent, the approved external private root, TTL, expiry/revocation
   deletion, receipt/orphan ownership, the minimal signed connection-receipt
   fields, provider-authenticated consumer signer binding,
   provider-signed acknowledgement and one-invitation-unit cardinality, and
   separately consented public fields. Every candidate environment also proves
   a supported OS/architecture and fail-closed host-isolation backend before
   qualification. The protocol states that provider export and public commit
   cannot be recalled.

### Design-Partner Pilot

1. `PILOT-001`: The `DISC-001` provider advances with the committed migration
   and preregistered protocol without changing thresholds after outcomes are
   visible.
2. `PILOT-002`: At least five distinct prequalified eligible repositories from
   the frozen `DISC-002` cohort, spanning at least three Consumer
   Organizations, are invited. A repository contributes at most one funnel
   unit.
3. `PILOT-003`: At least three distinct Consumer Organizations each explicitly
   issue repository authorization and import a provider-signed acknowledgement
   for one valid minimal connection receipt bound to an authenticated consumer
   signer and one eligible-repository invitation unit within 30 days of
   invitation. Across those connected repositories, median calendar time from
   invitation receipt to the first valid impact report is at most seven days,
   median total Consumer Organization security, privacy, platform, and
   maintainer hands-on labor over that interval is at most two hours, and
   median Consumer Maintainer in-product hands-on time from starting `campaign
   accept` through authorization issuance to the first valid impact report is
   at most 60 minutes. Lumyn and provider assistance are measured separately;
   custom product code or provider access to source invalidates the
   observation.
4. `PILOT-004`: At least three evidence-backed draft PRs are generated for
   three distinct eligible repositories within 14 calendar days of repository
   authorization; the clock has no test-duration or other pause.
5. `PILOT-005`: At least three PRs from three distinct eligible repositories
   across at least three Consumer Organizations merge within 45 days of PR
   creation. This threshold is frozen before the first invitation and cannot
   be rebaselined after outcomes are observed.
6. `PILOT-006`: At least two of the first three merged PRs require no
   substantive manual correction from PR creation through a fixed 30-day
   post-merge observation window. A migration-attributable revert or
   corrective fix counts as substantive correction. At least two of those
   three PRs also achieve `workflow_verified_mock` or
   `workflow_verified_sandbox` from the exact patched head before merge; a
   missing approved entrypoint or weaker evidence does not pass this gate.
7. `PILOT-007`: Campaign setup, onboarding, support, repository-based
   conversion, automation, lead-time, contribution, cost-per-merge, and actual
   legacy-retirement duration are measured against a baseline frozen before
   the first invitation and judged by the preregistered absolute campaign
   deadline. Campaign contribution is non-negative, and the one preregistered
   primary provider outcome achieves either at least 20% improvement in
   support hours per merged repository, at least 20% improvement in
   authorization-to-merge lead time, or retirement of the targeted legacy
   version from at least 60% of the frozen eligible-repository cohort by day
   120. A missing or `not_measurable` primary outcome fails this gate.
8. `PILOT-008`: The provider pays the preregistered pilot invoice at a price at
   least equal to frozen Lumyn campaign COGS and executes a paid renewal,
   expansion, or purchase order before the campaign judgment deadline.
   Letters of intent, compliments, nominal payments, and informal
   willingness-to-pay statements do not qualify.
9. `PILOT-009`: Paid continuation covers repeatable platform use: an annual
   connected-repository program or a second named migration campaign.
   Preregistered recurring delivery economics show at least 60% projected gross
   margin after Lumyn onboarding, support, infrastructure, and external-tool
   COGS, excluding separately itemized one-time core product development, and
   no more than four Lumyn operator hours per merged repository after campaign
   setup. Provider and consumer labor are reported separately as buyer total
   cost and never netted into Lumyn gross margin.

---

## Success Metrics

### North-Star Metric

`verified migration PRs merged per active provider campaign`

This is accompanied by the evidence boundary; a PR is not counted as verified
when required verification stages were skipped or stale.

### Funnel

Track:

```text
eligible
-> invited
-> consented
-> analyzed
-> affected
-> plan approved
-> patch generated
-> repo verified
-> draft PR opened
-> merged
```

Report conversion and time between every stage. A strong patch engine with poor
invite-to-consent conversion is a failed provider-led product.

### Technical Quality

- Impact recall and precision by supported change class.
- Supported holdout patches reaching `repo_verified`.
- Substantive manual-correction rate.
- Unrelated-edit rate.
- False `unaffected` rate.
- False-positive and false-negative rate for each canonical
  `workflow_verified_*` label.
- Deterministic repeatability.
- Median impact and PR-preparation duration.

Held-out scoring is frozen before M4 implementation scoring and uses:

```text
recall = TP / (TP + FN)
precision = TP / (TP + FP)
false_unaffected_rate =
  affected sites labeled unaffected / all affected sites
```

An `uncertain`, `unsupported`, or missed result on a known-supported affected
site is a false negative. Report each class independently and the micro/macro
aggregate. A verifier that always returns `gap` fails the positive cases.
Independent contract replay that does not execute patched code is not a
workflow-positive case.

Only a non-resolving opaque holdout manifest is committed. It contains opaque
case IDs, provenance class and license posture, a frozen suite commitment, and
encrypted or HMAC artifact commitments. It contains no source URL, repository
or package identifier, plaintext content digest, input, answer key, expected
label, expected patch, or raw trace. An independent holdout owner provisions
and freezes the private suite, and the lifecycle-owned M1 `holdout_result`
binds its frozen commitment before implementation scoring. Held-out material
lives under an evaluator-controlled private root such as
`LUMYN_HOLDOUT_ROOT`. That root is not mounted into, passed to, named in the
prompt for, or exposed through the environment of `task-executor`, M4, M6, or
M7. Only `holdout-evaluator` receives it, and the repository receives aggregate
counts, commitments, and the lifecycle-owned result—not answer material.

The static task contract does not fabricate that commitment before the private
suite exists. M1 uses `holdout_policy.mode: provision` and
`holdout_provisioning_required` with an opaque private namespace and the
`hmac-sha256` commitment algorithm. Provisioning proves future-suite creation;
it does not claim that M1's current candidate passed hidden evaluation. Its
independent result creates the opaque suite ref and keyed commitment. M4, M6,
and M7 use `holdout_policy.mode: evaluate` and
`holdout_evaluation_required`, resolving that trusted M1 result; each
evaluation result binds the exact provisioning-result ref and bytes used for
the run. Replacing the referenced result invalidates evaluation evidence.

### Provider Economics

- `campaign_setup_hours`: Lumyn plus provider time for packet authoring, canary
  fixtures, and cohort preparation before the first invitation.
- `median_invitation_to_first_impact_days`: Median calendar time from a
  Consumer Organization's receipt of the invitation to its first valid impact
  report.
- `consumer_activation_hands_on_hours`: Total Consumer Organization security,
  privacy, platform, and maintainer hands-on labor from invitation receipt to
  the first valid impact report.
- `consumer_maintainer_in_product_hours`: Consumer Maintainer hands-on time
  from starting `campaign accept` through authorization issuance to the first
  valid impact report.
- `lumyn_onboarding_hours`: Lumyn operator time from invitation to the first
  valid impact report, excluding core product development.
- `provider_onboarding_hours`: provider operator time from invitation to the
  first valid impact report.
- `provider_support_hours_per_merged_repo`: provider DX, support, and solutions
  hours during the measurement window divided by merged PRs.
- `median_authorization_to_merge_days`: median calendar days from the
  repository authorization timestamp to merge for the frozen cohort.
- `automation_rate`: merged PRs without substantive correction divided by
  `affected_supported` repositories.
- `eligible_to_merge_conversion`: distinct merged eligible repositories divided
  by the frozen eligible-repository cohort; every repository contributes at
  most one numerator and denominator unit.
- `acknowledged_connected_repositories`: distinct non-replayed minimal
  connection receipts whose invitation, packet, opaque organization/repository
  IDs, provider-authenticated consumer key binding, issuer, audience, expiry,
  nonce, consent-policy digest, signature, one-unit cardinality, and
  provider-signed acknowledgement were verified. This is the only
  connected-repository billing denominator.
- `median_invite_to_merge_days`.
- `buyer_total_cost_per_merged_pr`: loaded provider and consumer labor plus the
  provider fee and provider-borne external spend, divided by merged PRs.
- `lumyn_campaign_cogs`: frozen campaign-specific Lumyn onboarding, operations,
  support, infrastructure, and external-tool cost, excluding separately
  itemized one-time core product development.
- `campaign_contribution`: recognized pilot revenue minus
  `lumyn_campaign_cogs`.
- `projected_recurring_gross_margin`: annual-program or second-campaign revenue
  minus projected recurring Lumyn COGS, divided by that revenue. Provider and
  consumer labor is buyer total cost, not Lumyn COGS.
- `lumyn_operator_hours_per_merged_repo`: Lumyn onboarding, operations, and
  support hours after campaign setup divided by merged PRs.
- `legacy_retirement_delta`: current actual retirement duration, measured from
  the current packet publication or migration announcement to actual legacy
  retirement, minus the comparable prior migration's actual duration. If the
  current legacy surface has not retired by the frozen judgment deadline or no
  reliable comparator exists, report `not_measurable`; do not claim
  improvement.
- `legacy_cohort_retirement_rate`: distinct frozen eligible repositories whose
  resolved dependency state no longer includes the targeted legacy SDK/API
  version by day 120, divided by the frozen eligible-repository cohort.
- `primary_provider_outcome`: exactly one of
  `provider_support_hours_per_merged_repo`,
  `median_authorization_to_merge_days`, or
  `legacy_cohort_retirement_rate`, selected before the first invitation. The
  first two pass only at 20% or greater improvement against a comparable
  frozen baseline; the third passes only at 60% or greater. A missing
  comparator makes the first two ineligible for selection, and a missing or
  `not_measurable` selected outcome fails `PILOT-007`.
- Paid pilot receipt, campaign contribution, an annual connected-repository or
  second-migration continuation, projected recurring gross margin, and
  executed purchase evidence.

`substantive_manual_correction` means a human edit from PR creation through the
fixed 30-day post-merge observation window that changes the migrated API/SDK
invocation, request/response mapping, error handling, workflow behavior, or a
Lumyn-generated semantic edit. A migration-attributable post-merge revert or
corrective fix counts. Formatting, deterministic lockfile normalization,
comments, and unrelated pre-existing CI repair do not count, but must still be
recorded.

Before the first invitation, freeze the distinct eligible-repository cohort,
organization IDs, one accountable maintainer per repository, invitation
window, 30-day consent window, minimal connection-receipt schema,
provider-authenticated consumer signer binding, acknowledgement signing key,
endpoint or offline exchange, one-invitation-unit cardinality and deduplication
rule, supported consumer OS/architecture/isolation-backend matrix, provider
status channel and maximum-age policy, unpaused 14-calendar-day PR window,
45-day merge window,
30-day post-merge observation window, an absolute campaign judgment deadline
no later than 120 days after the first invitation, comparable prior migration,
one primary provider outcome and its material threshold, loaded labor rates,
the boundary between one-time core product development and Lumyn campaign
COGS, pilot price, contribution threshold, recurring gross-margin and
operator-hour thresholds, provider and consumer buyer-TCO treatment,
correction rubric, exact provider-visible fields, public evidence fields, the
irreversibility of provider/public disclosure, confidentiality, retention, and
consented-data protocol.

---

## Distribution

The initial distribution motion is provider-led and services-assisted:

1. Recruit a provider with an active migration, recurring-value potential, and
   a named cohort.
2. Create the standard campaign kit, signed packet, signed invitation, and
   provider canary.
3. Prove the synthetic offline end-to-end canary before demonstrating the
   product.
4. Let the provider invite consumer maintainers through the standard artifact.
5. Use consumer-side `campaign accept`, explicit authorization issuance, and
   the local/CI runner to earn trust.
6. Count a sponsored connection only from the acknowledged minimal signed
   receipt; return richer status only under separate consent.
7. Convert successful campaign evidence into an annual connected-repository
   program or second named migration campaign.

The current repository and design-partner distribution are not represented as
open source. A public OSS or self-serve release is a separate gate requiring an
approved license, `SECURITY.md`, contribution and support policies, signed
release provenance, install-integrity verification, and a maintained
vulnerability-response route.

Public Stripe, GitHub, OpenAI, or synthetic fixtures are engineering
benchmarks. They are not customer proof and must not be represented as provider
endorsement.

---

## Current Repository Baseline

Implemented today:

- Go CLI, config, exit-code, and command-result foundation.
- `lumyn init`.
- `lumyn check`.
- OpenAPI and local-doc source parsing, refs, fingerprints, and findings.
- Executable workflow/evidence/cassette/proof/boundary/redaction schemas.
- CI, coverage, CodeQL, branch-policy, review, Factory planning, and shipping
  governance.

Designed or schema-backed but not implemented as a runtime:

- workflow recording;
- replay verification;
- live sandbox verification;
- HTML report rendering;
- GitHub product integration;
- migration impact analysis;
- patch generation;
- migration PR delivery;
- provider enrollment and provider-status refresh;
- explicit authorization issuance and revocation;
- provider-authenticated connection-receipt issuance, acknowledgement, and
  import;
- fail-closed host command isolation;
- live agent evaluation.

The current command dispatcher recognizes several unimplemented commands and
can return a base success envelope. The first implementation task must correct
that false-green behavior.

Existing task-run, PR-lifecycle, pilot, and closure evidence remains historical
and immutable. New acceptance items may cite old evidence only where the old
semantics genuinely prove the new foundation.

---

## Risks

1. **Two-sided activation:** Providers may pay but fail to recruit consumer
   maintainers.
2. **Repository trust:** Consumers may reject the authorization or execution
   model.
3. **Supply-chain risk:** A compromised provider packet could attempt to steer
   customer code changes.
4. **Codemod ceiling:** Simple mappings may look interchangeable with vendor
   codemods unless impact coverage and workflow proof add clear value.
5. **Semantic ambiguity:** Real migrations may require business values that no
   provider packet can infer.
6. **Repository complexity:** Wrappers, monorepos, generated code, multiple SDK
   versions, and dynamic use may dominate.
7. **Sandbox mismatch:** Sandbox behavior may not represent production.
8. **False confidence:** A weak verification stage may be mistaken for business
   outcome proof.
9. **Episodic demand:** Consequential migrations may be too infrequent for the
   intended recurring model.
10. **Provider reachability:** Providers may not know which repositories or
    maintainers depend on them.
11. **Data leakage:** Logs, traces, diffs, or attestations could expose consumer
    information without strict separation and redaction.

---

## Falsification And Reframe Gates

Reconsider the provider-led thesis when either two qualified provider
commitment attempts fail `DISC-001` or `DISC-002`, or the first qualified
campaign misses any required `PILOT` gate by its frozen absolute judgment
deadline. Judge earlier when a gate becomes mathematically impossible.
Abandonment, no invitation, insufficient authorization, late PR creation,
insufficient merges, poor correction-free or workflow-verified performance,
excessive onboarding or operator effort, negative contribution, recurring
gross margin below the frozen threshold, a missed or unmeasurable primary
provider outcome, missing repeatable paid continuation, or a
migration-attributable post-merge revert all count as failure; they cannot keep
the experiment open indefinitely. Do not select only successful applicants,
change the cohort after seeing outcomes, or rebaseline a threshold after the
first invitation.

Also record, without hiding it behind a broader scope, when:

- most affected repositories require unsupported business input or dynamic
  analysis;
- providers cannot supply stable prerelease artifacts or verification
  semantics;
- migration frequency cannot support recurring provider value.

These are product-learning gates, not criteria to conceal through broader
scope.

---

## Non-Goals

The MVP does not:

1. scan repositories that have not been explicitly authorized;
2. claim coverage of every downstream integration;
3. grant API Providers access to consumer code by virtue of payment;
4. execute arbitrary provider-supplied scripts;
5. infer missing business values;
6. support every language, SDK, package manager, or API style;
7. support auth, webhook, event, GraphQL, gRPC, or production migrations;
8. regenerate arbitrary clients;
9. become a generic dependency updater or coding agent;
10. sell generic agent readiness or agent evaluation;
11. require OpenAI or Anthropic model-provider adapters;
12. make a hosted dashboard a prerequisite for consumer execution;
13. use production credentials or mutate production;
14. write to the default branch;
15. auto-merge;
16. train shared models on private consumer artifacts;
17. present public benchmark fixtures as commercial validation.

---

## Definition Of MVP Success

Lumyn reaches commercial MVP success only when:

- all required technical acceptance items for the supported change classes are
  implemented with evidence;
- no required trust or authorization item is deferred without explicit
  approval;
- the fixed benchmark has zero false workflow-verification results;
- the standard provider kit, signed invitation, consumer-connect flow, and
  offline synthetic end-to-end canary pass without hidden operator authority;
- one qualified provider campaign completes the design-partner pilot gates;
- provider and consumer participants both accept the evidence and control
  model;
- paid continuation covers an annual connected-repository program or second
  named migration, and meets the frozen recurring gross-margin and
  operator-effort gates.

The implementation sequence and bounded task packets are defined in
`docs/product/plan.md` and
`.factory/artifacts/prd-to-plan/lumyn-migration-mvp/`.
