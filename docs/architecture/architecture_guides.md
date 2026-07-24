# Lumyn Architecture Guide

## Architecture Objective

Lumyn turns a signed API-provider change, verified against a consumer-pinned
provider trust root and provider-to-package ownership binding, into a bounded,
verified draft PR inside a consumer-authorized repository without granting the
provider implicit access to consumer code.

The product begins with a standard campaign kit and signed invitation, not an
ad hoc operator workflow. Before invitation acceptance, the consumer enrolls
the provider root/package binding from a separately authenticated bundle and
out-of-band confirmed fingerprint. Enrollment and acceptance persist only
outside the checkout and produce a no-authority authorization request. A
synthetic offline canary must prove the complete local path through a draft-PR
preview before any live integration is presented as working.

The architecture optimizes for:

- explicit authority;
- local or consumer-CI execution;
- deterministic, explainable impact and patching;
- proof-honest verification;
- fail-closed mutation and disclosure;
- evidence bound to exact inputs.

## Trust And Data Planes

### Provider Campaign Plane

Owns:

- API-provider identity;
- provider change packet and lifecycle;
- source/target versions and provenance;
- invited cohort;
- compatibility deadline;
- sandbox semantics;
- provider-authenticated consumer receipt-key bindings and provider-signed,
  deduplicated minimal connection-receipt acknowledgements;
- optional consumer-consented richer status attestations.

Does not own:

- consumer repository access;
- consumer commands or credentials;
- raw source, diffs, logs, traces, prompts, or responses;
- consumer branch or merge authority.

### Consumer Execution Plane

Owns:

- repository and selected package root;
- read/write path scopes;
- branch namespace;
- command allowlist;
- command host-isolation policy and backend;
- provider-status snapshot and exact online-refresh policy;
- dependency lifecycle-script posture;
- network allowlist;
- non-production sandbox credentials;
- private impact, plan, patch, and verification evidence;
- an explicitly configured private runtime root outside the checkout and any
  public source repository;
- automatic private-record TTL and expiry/revocation deletion on creation,
  read, startup, and the next run, plus deletion receipts or orphan reports;
- draft-PR approval;
- review and merge.

### Optional Future Model Plane

Model-assisted patching is outside the MVP. If introduced, it is a third trust
boundary with independent endpoint, credential, data-disclosure, budget,
retention, and output-verification policy. It never inherits API-provider or
consumer authority.

## Authority Flow

```text
API Provider
  publishes signed declarative packet
          |
          v
Provider campaign plane
  distributes immutable packet; invites; authenticates receipt signers;
  signs deduplicated receipt acknowledgements and consented richer status only
          |
          | invitation (no repo authority)
          v
Consumer Maintainer
  pins provider trust root -> accepts no-authority request
  -> explicitly issues signed read authorization
  -> optionally signs and submits connection receipt
  -> verifies and imports provider-signed acknowledgement
  -> reviews plan
  -> grants local write/commands -> grants registry if needed
  -> grants sandbox and payload disclosure if needed
  -> grants remote branch -> grants draft PR
          |
          v
Consumer execution plane
  resolves signed current provider status -> verifies signature/package binding
  -> impact -> patch -> host-isolated verify -> draft PR
```

Each arrow is an explicit state transition with independent evidence.

## Product State Machines

### Provider Packet

```text
draft -> published -> superseded | withdrawn
```

- Published content is immutable.
- Published means immutable for the authorized audience, not publicly
  disclosed.
- A new digest creates a new version.
- Signature, issuer-to-package binding, audience, issue time, expiry, key
  rotation, revocation, withdrawal, and replay checks must all pass.
- Superseded or withdrawn packets cannot start new patch runs.

### Consumer Authorization

```text
request_rendered -> issued -> active -> expired | revoked
```

A request has no authority. Issuance and revocation are explicit
consumer-signed private-state actions and perform no granted side effect.
Repository read, local write, host-isolated command execution, provider-status
read, package-registry network, sandbox network, sandbox credential, sandbox
request disclosure, remote branch write, PR write, minimal campaign receipt,
richer provider reporting, retention, and deletion grants have separate state,
scope, expiry, and revocation. No grant implies another.

### Migration

```text
not_analyzed
-> impact_complete
-> plan_pending_approval
-> plan_approved
-> patch_generated
-> verification_complete
-> draft_pr_authorized
-> draft_pr_open
-> merged | closed | blocked | superseded
```

Stale packet, plan, base/head, or evidence moves the migration to `blocked` or
`superseded`; it does not silently resume.

## Component Boundaries

| Component | Responsibility | Must not own |
|---|---|---|
| CLI | command parsing, config, JSON envelope, exits | product inference |
| Source intake | pinned OpenAPI/docs/SDK refs, parsing, fingerprints | call-site patching |
| Change normalizer | supported semantic change entries | consumer repo writes |
| Packet verifier | signature, trust root, API-provider/package binding, lifecycle, provenance, digest, declarative policy | packet authoring authority |
| Provider-status resolver | pinned signer, signed offline snapshots, exact authorized endpoint reads, freshness and replay | repository or consumer disclosure |
| Authorization engine | no-authority requests, consumer-signed issue/revoke, exact grants, expiry, revocation | provider campaign intent or side-effect execution |
| TypeScript analyzer | SDK/version, imports, call sites, exclusions | file mutation |
| Impact engine | applicability and coverage classification | patch application |
| Migration planner | complete no-write plan | filesystem writes |
| Workspace manager | isolated workspace and safe paths | semantic transformations |
| Patch engine | three deterministic recipes and manifest | arbitrary business inference |
| Command runner | allowlisted baseline/build/test execution through exact fail-closed host isolation | ambient host access or sandbox credential distribution |
| Replay verifier | deterministic workflow evidence | live outcome claims |
| Sandbox verifier | approved live read-back, budgets, cleanup | production access |
| Evidence engine | normalized axes, hashes, freshness, redaction | unsupported roll-up claims |
| Retention manager | automatic TTL/revocation deletion, receipts, retry, orphan recovery, `artifacts gc` | authority extension or historical-closure rewriting |
| PR adapter | idempotent draft branch/PR delivery | merge |
| Receipt exchange | provider-authenticated consumer signer binding, minimal receipt issue/submit, provider verification/deduplication/signing, consumer acknowledgement import | repository identity or richer outcome status |
| Attestation exporter | consumer-signed, receipt/evidence-bound richer campaign status | raw consumer evidence |
| Outcome ingester | merged/closed/reverted/corrected provenance | silent recipe mutation |

Keep these boundaries in separate packages or small interfaces. Do not add
impact, patch, GitHub, or attestation behavior to `internal/source`.
M3 implements the packet verifier in `internal/trust/`; schema-only M2 work and
commercial canary inspection cannot substitute for that runtime boundary.
M7 exposes the evidence engine through local-only `lumyn trace` rendering;
rendering does not cross the provider-attestation or network boundary.

## Initial Architecture Spine

```text
provider source/target artifacts
-> standard campaign kit
-> configured signing and immutable provider change packet
-> signed invitation
-> packet validation
-> independently authenticated provider enrollment
-> consumer campaign acceptance
-> explicit consumer-signed authorization issuance
-> [optional authenticated signer binding, signed minimal receipt, and provider-signed acknowledgement]
-> signed current provider-status resolution
-> TypeScript repository analysis
-> impact report
-> migration plan
-> plan approval
-> isolated deterministic patch
-> host-isolated repository baseline and post-patch checks
-> replay/mock workflow evidence
-> [independently authorized sandbox evidence]
-> migration evidence packet
-> explicit draft-PR delivery
-> [optional consented campaign attestation]
-> outcome evidence
```

Sandbox verification and draft-PR delivery are independent successors to the
offline verification spine. M9 cannot inherit or require M8 sandbox authority,
and absence of provider-reporting consent suppresses only that optional output.

## Artifact Ownership

### Provider-Controlled Inputs

- provider change packet;
- public or prerelease source/target refs;
- canary fixtures;
- sandbox and rollback guidance.

### Consumer-Private Artifacts

- repository authorization;
- signed authorization revocation;
- provider enrollment and signed status snapshots;
- provider-authenticated consumer receipt-key bindings, signed campaign
  connection receipts, and provider-signed acknowledgements;
- detailed impact report;
- migration plan;
- patch manifest and diff;
- command logs;
- workflow traces;
- credentials and secret-bearing runtime state;
- detailed verification report.
- evaluator-controlled held-out repositories, inputs, answer keys, expected
  patches or labels, and raw traces. An independent holdout owner provisions
  and freezes these before scoring; they remain outside the checkout and every
  implementation-worker mount, prompt, and environment.

### Provider-Visible By Explicit Consent

- minimal provider-acknowledged connection unit containing only campaign and
  packet binding, eligible-repository unit and invitation nonces, opaque
  per-campaign organization/repository IDs, connection event/time, tool/schema
  versions, consent-policy and authenticated key-binding digests, consumer
  issuer fingerprint, audience, expiry, nonce, signature, and the
  provider-signed receipt-digest/deduplication acknowledgement;
- richer campaign/repository pseudonymous status only under a separate grant;
- authorization state;
- impact found/no impact/blocked state;
- draft PR opened;
- merged/closed;
- consumer-approved failure category;
- timing and aggregate outcome fields.

Provider-visible schemas use field allowlists; they are not redacted copies of a
larger private object.

### Storage And Disclosure Boundary

- Consumer-private runtime and identifiable pilot artifacts live in an
  explicitly configured consumer-controlled state root outside the checkout
  and any public source repository.
- The runtime rejects a private root that resolves inside either boundary,
  including through a symlink. Repository ignore rules for legacy roots are
  defense in depth, not storage authorization.
- Private artifacts carry authorization-bound TTL and are deleted on expiry or
  revocation automatically at creation, read, process startup, and the next
  run. Cleanup emits a deletion receipt or an orphan report.
- `lumyn artifacts gc` is an explicit recovery surface for retrying partial
  deletion and inspecting unresolved orphans. It cannot extend TTL, revive
  revoked authority, or alter historical closure claims.
- The committed holdout manifest may contain only opaque case IDs,
  non-resolving provenance class and license posture, a frozen suite
  commitment, and encrypted or HMAC artifact commitments. It contains no
  source URL, repository or package identifier, plaintext content digest, raw
  private content, or machine-local path. The M1 lifecycle-owned
  `holdout_result` binds independent provisioning and that frozen commitment.
  M1 declares provision mode, `holdout_provisioning_required`, an opaque
  private namespace, and the HMAC algorithm; that gate creates future-use
  evidence and is not current-candidate evaluation. M4/M6/M7 declare
  `holdout_evaluation_required`; their evaluate-mode policies resolve and bind
  the exact trusted M1 result bytes. Static planning never fabricates the
  future commitment.
- Provider-visible and public are separate consent decisions. The only public
  pilot artifacts are explicitly public, consented, redacted aggregates and
  evidence hashes under
  `.factory/artifacts/pilot/lumyn-migration-mvp/public/`.
- Provider export and public commit are irreversible disclosure boundaries.
  Revocation stops future sharing and deletes Lumyn-controlled private copies;
  it cannot recall provider records, Git history, clones, or caches.

### Factory Worker Versus Product Authority

Factory controls repository implementation work with its closed
`approval`, `credentials`, and `network` capabilities. Lumyn controls consumer
operations with separate private, schema-backed product grants such as
`customer_repo_read`, `provider_trust_status_read`, `command_execution`,
`sandbox_request_disclosure`, `github_pr_write`, `campaign_receipt`,
`provider_attestation`, `artifact_retention`, and `artifact_deletion`.
Factory approval may cite the opaque digest of a validated product
authorization bundle, but the bundle is evaluated by Lumyn's authorization
contract and no Factory grant conveys product authority. Factory selection and
dispatch validate implementation-worker controls and task declarations, not
the private product bundle. Lumyn revalidates current packet trust and the
exact applicable bundle internally, immediately before every local, sandbox,
GitHub, or attestation side effect and retry. A standalone
`lumyn authorization validate` result is diagnostic and closure evidence, not
cached runtime authority.

The narrower M2.5 evidence-handling approval binds the canonical digest of its
exact private fields, storage boundary, consent, TTL, deletion, orphan,
minimal-receipt, public-disclosure, and disclosure-irreversibility scope. A
generic approval or a digest from an earlier scope cannot release collection.

### Independent Promotion Evidence

When a task requires independent review, the lifecycle order after
`code-review` is `holdout-evaluator`, `trace-grader`, then
`evidence-attestor`, before `commit-push`. These evaluators are external or
human-operated trust principals, not implementation-worker modes. Each emits a
schema-valid passing artifact bound to the exact task, work item, lifecycle run,
current validation run, candidate digest, and work-proof marker digest in the
trusted Factory evidence root. The lifecycle namespace is not writable by the
implementation worker. Shipping fails before commit or PR creation if evidence
is missing, stale, replayed, self-authored, malformed, or non-passing, or if its
review lens/reviewer class differs from the task-level requirement.
Only workers selected by the task run: M1 provisions and freezes its suite
through `holdout-evaluator`, while M10 uses only `evidence-attestor` to review
privacy-approved campaign calculations and never receives the benchmark
holdout root.

## Evidence Model

Migration evidence has orthogonal axes:

- impact coverage;
- patch state and provenance;
- repository verification;
- workflow environment and outcome;
- cleanup, boundary, and redaction;
- delivery;
- permission state;
- residual risk.

Evidence binds to:

- packet digest;
- source and target artifact digests;
- repository base and head commits;
- plan digest;
- patch/recipe digests;
- command and environment identity;
- workflow/cassette/sandbox identity;
- evidence artifact hashes.

Any bound input change invalidates dependent evidence.

Verification uses the canonical labels `not_run`, `static_verified`,
`repo_verified`, `workflow_contract_replay_passed`,
`workflow_verified_replay`, `workflow_verified_mock`,
`workflow_verified_sandbox`, `partial`, `failed`, `gap`, and `stale`.
Independent contract or cassette replay is
`workflow_contract_replay_passed` and cannot exceed `repo_verified`.
`workflow_verified_replay`, `workflow_verified_mock`, and
`workflow_verified_sandbox` require an approved entrypoint executed from the
exact patched repository head and observed interaction and outcome evidence in
that named environment. Results copied from the base commit, another head, or
an independently executed contract/cassette are not causal patched-head proof.
A missing causal execution or any boundary, cleanup, redaction, freshness, or
evidence-integrity failure prevents a workflow-verified label.

## Structured Parser Boundaries

- OpenAPI, JSON, YAML, manifests, lockfiles, schemas, CI results, and GitHub
  responses use structured parsers or APIs.
- TypeScript supported impact uses an AST or comparably structured parser.
- Text search may seed discovery but cannot alone prove a patchable call site.
- The consumer selects one canonical package/read root. All manifests,
  lockfiles, TypeScript sources, `tsconfig` `extends`, project references, and
  resolved module paths must remain inside that root after real-path
  resolution.
- Symlink escape, path traversal, an out-of-root `tsconfig` reference, or an
  ambiguous/multiple package root fails closed before analysis.
- External refs remain blocked in deterministic tests.
- Parser cycles, missing refs, and source/SDK disagreement fail closed.

## Patch Safety Boundary

Patch application must:

- start from an approved immutable plan;
- revalidate the current packet bytes, digest, provider trust root and package
  binding, lifecycle, audience, expiry, rotation, revocation, withdrawal,
  supersession, and replay state immediately before every write; plan-time
  trust is never cached across a mutation boundary;
- run in an isolated worktree or disposable equivalent;
- canonicalize and validate real paths;
- reject path traversal and symlink escape;
- exclude generated, vendored, minified, cache, and build output;
- enforce writable paths, file count, line count, and diff content budgets;
- map every edit to packet change and recipe IDs;
- remain deterministic for identical pinned inputs;
- mutate `package-lock.json` only with exact Node and npm versions, a pinned
  registry or offline snapshot, recorded package-integrity inputs, lifecycle
  scripts disabled, and a bound toolchain digest;
- produce rollback evidence;
- leave the default branch untouched.

## Command Execution Boundary

Repository commands are untrusted code:

- exact allowlist and working directory;
- exact read-only/writable mounts plus neutral home and temp roots;
- explicit executable and toolchain roots;
- timeout and output budgets;
- network disabled by default;
- package-registry access requires a separate `package_registry_read` grant
  restricted to the approved registry or snapshot and package set;
- dependency lifecycle scripts disabled by default;
- exact Node/npm versions, registry identity, package-integrity inputs, and
  toolchain digest recorded for any lockfile mutation;
- sanitized environment classes and no ambient secrets;
- no host home, SSH/GPG/cloud credential stores, keychain or OS credentials;
- no agent, Docker, or unrelated local-service sockets;
- no inherited file descriptors beyond standard streams;
- child processes inherit mounts, environment, socket, descriptor, credential,
  network, and resource restrictions;
- a supported fail-closed isolation backend is mandatory and an unavailable or
  unverifiable backend blocks before launch;
- sandbox credentials absent from build/test stages;
- pre-patch baseline separated from post-patch result;
- output stored by bounded artifact ref and hash;
- cancellation and crash recovery leave a typed result.

## Live Sandbox Boundary

Sandbox verification requires:

- task-scoped API-provider and environment identity;
- non-production credentials;
- a `sandbox_network` grant naming the destination and operation allowlists;
- an independent `sandbox_credential` grant naming the non-production
  credential class, scopes, injection stage, expiry, and revocation;
- a separate `sandbox_request_disclosure` grant naming transmitted payload
  classes;
- synthetic or explicitly approved non-sensitive test data only; production
  customer data, PII, credentials, and secrets are prohibited;
- declared provider logging behavior, retention period, and deletion terms;
- namespace and idempotency key;
- request/write budgets;
- settle/retry limits;
- cleanup;
- orphan evidence;
- explicit statement that sandbox proof is not production proof.
- current packet trust and every exact product grant revalidated inside Lumyn
  immediately before each sandbox call and retry from a current signed offline
  status snapshot or exact `provider_trust_status_read` grant; Factory worker
  grants do not confer this authority.

Production access is outside the MVP.

## GitHub Boundary

- Repository read, `github_branch_write`, and `github_pr_write` are separate;
  none implies another.
- Installation tokens are repository-scoped and short-lived.
- Only the authorized branch namespace may be used.
- Default-branch write is prohibited.
- PRs are draft-only.
- Idempotency binds campaign, repository, packet, and base state.
- Current packet trust plus branch, PR, reporting, retention, and deletion
  authority are revalidated inside Lumyn from a current signed offline status
  snapshot or exact `provider_trust_status_read` grant immediately before every
  remote read-modify-write and retry.
- Lumyn does not merge.

## Systems Thinking Map

State owners:

- provider packet: API Provider;
- API-provider trust root and authorization: API Consumer Organization;
- active plan and acceptance: Lumyn repository/Factory;
- patch and detailed evidence: consumer execution plane;
- attestation: consumer-consented export;
- merge outcome: consumer repository.

Feedback sources:

- schema and parser validation;
- corpus precision/recall;
- golden patch comparison;
- compile/typecheck/tests;
- replay and sandbox read-back;
- cleanup/orphan evidence;
- PR correction and merge outcome;
- provider support and migration metrics;
- invite-to-consent conversion.

Deletion blast radius:

- deleting packet provenance invalidates migration trust;
- deleting authorization invalidates further execution;
- expiry or revocation triggers private-record deletion and a deletion receipt
  or orphan report;
- revocation blocks future provider/public disclosure but cannot recall copies
  already exported or committed;
- deleting plan or patch hashes invalidates evidence;
- deleting workflow/cassette invalidates replay proof;
- deleting private logs may preserve a bounded attestation but prevents detailed
  audit;
- deleting active Factory artifacts prevents governed dispatch;
- deleting historical evidence must not be used to rewrite old closure claims.

Medium/high-risk tasks record state owner, feedback source, deletion impact,
rollback/deletion test, and source-of-truth changes.

## TDD And Red-First Expectations

- Behavior changes start with a failing unit, fixture, contract, golden patch,
  scenario, or permission test when practical.
- Schemas add invalid fixtures before accepting new shapes.
- Impact tasks fix ground truth before scoring.
- Patch tasks fix expected patches before implementation.
- Verification tasks include false-green cases first.
- GitHub and sandbox tasks use mocks before live access.
- If red-first is impractical, validation evidence records the exact reason and
  compensating proof.

## ADR And Decision Triggers

Require an ADR or decision update for:

- provider/consumer authority;
- execution-plane or data-sharing changes;
- command or JSON contract changes;
- schema compatibility;
- parser runtime boundary;
- patch isolation or filesystem ownership;
- credential or network posture;
- GitHub permissions;
- workflow proof semantics;
- model-assisted patching;
- hosted campaign coordination;
- release/distribution posture;
- major performance or reliability tradeoffs.

ADR-0002 records the product and trust reframe. ADR-0001 remains historical
context for the retained evidence foundation.

## Performance And Cost Triggers

- Impact analysis target: median under five minutes on the fixed corpus.
- Draft-PR preparation target: median under twenty minutes excluding
  repository-defined test duration.
- Parser startup, repository size, file count, AST memory, test duration,
  sandbox requests, artifact size, and GitHub calls are explicit budgets.
- Fan-out across repositories is not introduced before one-repository
  determinism and authorization are proven.
- Model cost is not an MVP budget because model-assisted patching is deferred.

## Reliability And Recovery Triggers

Test:

- interrupted workspace creation;
- stale base or packet;
- partial patch;
- command timeout;
- flaky/pre-existing tests;
- sandbox timeout and drift;
- cleanup and orphan failure;
- GitHub retry and duplicate PR;
- revoked or expired authorization;
- revocation, packet rotation, or packet withdrawal between validation and a
  local, sandbox, or remote side effect;
- partial private-artifact deletion, process restart, retry, deletion-receipt
  loss, and orphan recovery;
- redaction failure;
- attestation retry/idempotency.

Retries never widen permissions or change idempotency identity. Recovery cannot
reuse stale evidence.

## Trust-Mode Posture

### Deterministic Benchmark

- no network;
- no ambient secrets;
- no customer repository;
- no GitHub write;
- no provider sandbox;
- committed permitted fixtures only.

### Consumer Repository Read

- task-scoped read grant;
- explicit consumer-signed issuance after a no-authority request;
- current signed offline provider-status snapshot or exact authorized online
  refresh;
- no mutation;
- no provider visibility.

### Consumer Mutation

- approved plan digest;
- scoped isolated write;
- exact commands through the fail-closed host-isolation backend;
- no network by default.
- package-registry access, when required, is a separate destination- and
  package-scoped grant.

### Live Sandbox

- separate non-production credential/network and request-disclosure grants;
- synthetic or approved non-sensitive payloads only;
- declared provider logging, retention, and deletion terms;
- cleanup and orphan evidence.

### Draft PR

- separate remote-branch-write and PR-write grants;
- draft-only non-default branch.

### Provider Attestation

- exact field-level consumer consent.

## Runtime Shape

Current:

- `cmd/lumyn/`: command entry and process exit.
- `internal/config/`: repo-local config.
- `internal/result/`: stable command-result envelope.
- `internal/exitcode/`: stable exit codes.
- `internal/source/`: OpenAPI/docs parsing and findings.
- `internal/version/`: version metadata.
- `schemas/`: current executable contracts.

Planned boundaries:

- `internal/change/`
- `internal/trust/`
- `internal/authorization/`
- `internal/isolation/`
- `internal/receipt/`
- `internal/typescript/`
- `internal/impact/`
- `internal/migrationplan/`
- `internal/workspace/`
- `internal/patch/`
- `internal/verify/`
- `internal/replay/`
- `internal/live/`
- `internal/evidence/`
- `internal/retention/`
- `internal/report/`
- `internal/github/`
- `internal/attestation/`
- `internal/outcome/`

Private product artifact roots are added only outside the checkout and public
source repository, after their schema, explicit configuration, TTL, cleanup
ownership, deletion receipt/orphan behavior, and reference contract are
implemented. Only synthetic/licensed fixtures and consented aggregate/hash-only
public evidence are committable.

The current repository and pilot package are not called OSS. Design-partner
distribution requires explicit terms, security/support routes, signed
provenance, checksums, and install-integrity instructions. Public OSS/self-serve
requires an approved license and security, contribution, support, and
vulnerability-response policies.

## Architecture Budget And Decomposition

Source files warn at `1200` lines and fail at `2500` lines. New product domains
must not be appended to the source parser or repo-pack orchestration.

The rebaseline should shrink the plan validator below its previous
shrink-only ceiling. The existing architecture-debt exception remains valid
only while its recorded line ceiling matches the file and expires according to
policy. Remove it in a dedicated validated change when the validator no longer
needs an exception.

`internal/source` is already decomposed. Preserve that progress by placing
change normalization, TypeScript analysis, impact, patch, verification, GitHub,
and attestation in their own bounded packages.
