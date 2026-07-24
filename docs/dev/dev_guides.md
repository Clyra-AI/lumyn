# Lumyn Developer Guide

## Toolchain Pins

| Tool | Version |
|---|---:|
| Go | `1.26.5` |

Module path: `github.com/Clyra-AI/lumyn`.

The Go core remains authoritative for Lumyn artifact, authorization, impact,
patch, verification, and delivery orchestration. A future TypeScript parser or
tooling subprocess requires a pinned dependency, a bounded interface, an ADR
when it crosses the runtime boundary, and deterministic fixtures.

The supported npm migration slice must record exact Node and npm versions
before `package-lock.json` mutation is enabled. The migration plan and evidence
also bind the approved registry or offline snapshot, package-integrity inputs,
and toolchain digest. Floating versions, an unpinned registry view, lifecycle
scripts, or a missing integrity value make deterministic lockfile mutation
ineligible; there is no implicit local-machine fallback.

## Dependency Pins

- `github.com/santhosh-tekuri/jsonschema/v5 v5.3.1`: executable JSON Schema
  validation.

New dependencies must be pinned, justified in task evidence, scanner-covered,
and exercised by a failing test or fixture before implementation.

## Validation Matrix

- `make lint-fast`: repo contract, layout, policy, and Go vet.
- `make test-fast`: Go unit tests.
- `make test-coverage`: first-party Go coverage gate.
- `make test-contracts`: Go tests, schema tests, active migration-plan
  validation, historical pilot-evidence validation, and repo-pack self-tests.
- `make prepush-full`: full local gate before PR or merge.
- `make audit-remote-protection`: networked GitHub audit of `main` protection
  and the `protect-main-from-direct-push` ruleset.

The active migration control set contains exactly 62 acceptance items across
12 task packets (`M0` through `M10`, including `M2.5`). Task completion and
delivery slices are coverage lenses only; validation closes individual
acceptance item IDs.

## CI Lane Mapping

- Fast: `make lint-fast`, `make test-fast`.
- Core: `make test-contracts`, `make prepush-full`.
- Acceptance: item status in
  `.factory/artifacts/prd-to-plan/lumyn-migration-mvp/acceptance-ledger.json`
  and `scope-closure-map.json`.
- Cross-platform: reserved until standalone binary packaging.
- Risk: `CodeQL analyze`, plus targeted security/architecture review for
  parser, patch, external-call, credential, GitHub, data-sharing, and release
  surfaces.
- Release: reserved until supported binary packaging.
- Cross-system: explicit task-scoped GitHub or provider-sandbox checks.

## 12-Level Test Matrix

| Tier | Status | Migration-MVP evidence |
|---|---|---|
| Tier 1 Unit | Active | Go units through `make test-fast` |
| Tier 2 Integration | Planned/active | Schema, artifact, parser, impact, patch, and verification integration |
| Tier 3 End-to-End | Planned | `campaign kit create`, `change publish`, `campaign invite create`, `provider enroll`, `campaign accept`, `trust refresh`, `authorization issue/revoke/validate`, `campaign receipt issue/submit/acknowledge/ack import`, `canary run --offline`, `impact`, `migrate`, `verify`, `trace`, `artifacts gc`, and `pr` command flows |
| Tier 4 Acceptance | Active planning | Item-level active ledger and closure map |
| Tier 5 Hardening | Planned | Path escape, stale input, retry, cleanup, redaction, idempotency, and crash recovery |
| Tier 6 Chaos | Reserved | Controlled sandbox, GitHub, filesystem, and command-runner failures |
| Tier 7 Performance | Planned | Impact and PR-preparation budgets over fixed fixtures |
| Tier 8 Soak | Reserved | Repeated deterministic and campaign-idempotency runs |
| Tier 9 Contract | Active | JSON Schemas, typed exits, artifact compatibility, negative fixtures |
| Tier 10 UAT | Planned | Consumer-maintainer authorization and review workflow |
| Tier 11 Scenario | Planned | Gold, held-out, unsupported, and false-verification migration corpus |
| Tier 12 Cross-System Integration | Blocked until approved | M2.5-qualified provider/cohort, independently authorized M8 provider sandbox and M9 GitHub draft PR, and M10 real design-partner campaign |

Future task packets cite applicable tiers or an approved non-applicable reason.

## Coverage Gates

| Scope | Threshold | Enforcement |
|---|---:|---|
| Go first-party packages overall | `>= 75%` | `make test-coverage` and CI |
| Stable command or core packages | `>= 85%` | Hard gate calculated by `make test-coverage` |

Coverage output goes to `.factory/tmp/coverage.out`. Coverage is not a
substitute for schema fixtures, held-out impact scoring, golden patches,
proof-of-behavior scorecards, CodeQL, or cross-system evidence.

## Architecture Budget Gate

Source files warn at `1200` lines and fail at `2500` lines for supported source
extensions. Generated runtime, dependency, cache, and build directories are
excluded.

The repo-pack validation rebaseline must reduce or stay below its existing
shrink-only ceiling. Feature work keeps these responsibilities separate:

- source/change intake;
- TypeScript repository analysis;
- migration planning;
- patch application;
- command execution;
- workflow verification;
- GitHub delivery;
- provider attestation.

Do not turn `internal/source` or a single validator file into the new product
monolith.

## CI And PR Lifecycle

- Validation workflow: `.github/workflows/validate.yml`.
- Required check: `validate`.
- Security workflow: `.github/workflows/codeql.yml`.
- Required security check: `CodeQL analyze`.
- Required-check manifest: `.github/required-checks.json`.
- Owner-review map: `.github/CODEOWNERS`.
- Action-ref exceptions: `.github/action-ref-exceptions.yaml`.
- PR lifecycle report:
  `.factory/artifacts/pr-lifecycle/<work_item_id>/pr-lifecycle-report.json`.

Lifecycle-gated tasks require local validation, CI, review where required,
shipping, post-merge, and item-level closure evidence.

When task policy selects an independent gate, the canonical pre-shipping order
is `holdout-evaluator`, `trace-grader`, then `evidence-attestor`, after
`code-review` and before `commit-push`. These are external or human-operated
lifecycle reviews. Shipping verifies their lifecycle-owned artifacts are
schema-valid and passing, bind the exact task, work item, lifecycle run, current
validation run, candidate digest, and work-proof marker digest, and carry
independent worker provenance. Implementation-worker self-attestation is
invalid. The task-level review lens and reviewer class must exactly match the
inherited validation contract and current review artifact.

Passive Codex review settle is required before merge. Green CI alone is not merge-ready.
Do not merge manually through `gh pr merge`, the GitHub UI, or a
connector before the configured latest-head terminal review signal. A merge
without that evidence is a process escape and requires a recorded repair or
exception.

GitHub `main` remains protected by branch protection and the
`protect-main-from-direct-push` ruleset. Use `make audit-remote-protection` to
verify the live state.

## Security Scanner Enforcement

CodeQL is required for:

- dependency additions;
- generated-code or fixture generators;
- CI and workflow changes;
- structured parser boundaries;
- patch generation and filesystem writes;
- command execution;
- external network or API calls;
- credential, redaction, or data-sharing behavior;
- GitHub integration;
- release-sensitive work.

Scanner failure blocks closure unless a scoped, approved exception names the
owner, reason, expiry/follow-up, and compensating validation.

## Bootstrap Rules

- Deterministic benchmark work uses no network, sandbox credential, customer
  repository, GitHub write, or model key.
- Test-first or fixture-first development is expected.
- Committed Factory control and lifecycle evidence uses repo-relative paths.
  Consumer-private runtime and identifiable pilot evidence lives in an
  explicitly configured consumer-controlled state root outside the checkout
  and any public source repository; committed artifacts refer to it only by
  opaque identifier and digest, never a machine-local path.
- Factory worker grants use the closed `approval`, `credentials`, and `network`
  vocabulary. Exact Lumyn product grants are private schema-backed artifacts;
  task validation records diagnostic/closure proof over their opaque referenced
  bundle, and no Factory grant substitutes for product authority. Factory
  dispatch does not implement the Lumyn gate: product code revalidates current
  packet trust and exact product authority immediately before every side effect
  and retry.
- M2.5 identifiable-evidence handling requires a narrow manual privacy/legal
  preflight under Factory `approval`, including allowed fields, participant
  consent, external private storage, TTL, expiry/revocation deletion,
  deletion-receipt/orphan ownership, the minimal connection-receipt field
  allowlist, separate public disclosure consent, and the irreversibility of
  provider/public disclosure. This is not Lumyn runtime product authority. The
  active approval must cite this preflight and exactly match its canonical
  scope digest.
- Existing historical task-run and pilot evidence is immutable.
- T1 remains standard-library-only; later dependencies require task approval.
- Structured artifact changes include valid and invalid schema fixtures.
- Behavior, command, schema, artifact, permission, and evidence changes update
  docs and active Factory planning together.
- Runner-ready packets preserve acceptance IDs, paths, commands, risk,
  lifecycle gates, evidence, proof level, capability requirements, stop
  conditions, changelog/versioning intent, and semantic invariants.

## Docs Parity

User-facing sources:

- `README.md`
- `AGENTS.md`
- `WORKFLOW.md`
- `docs/product/prd.md`
- `docs/product/plan.md`
- `docs/dev/dev_guides.md`
- `docs/architecture/architecture_guides.md`
- relevant ADRs

Active planning sources:

- `.factory/artifacts/prd-to-plan/lumyn-migration-mvp/`
- Factory `profiles/lumyn.yaml`

Behavior, flags, status axes, exits, artifact paths, trust boundaries, and
implementation status must agree across these surfaces.

## Structured Data Policy

OpenAPI, JSON, YAML, package manifests, lockfiles, TypeScript ASTs, schemas,
coverage, GitHub responses, and logs use structured parsers or stable APIs.
Regex and text search may assist discovery but cannot be the only evidence for
a supported call-site or mutation.

Structured outputs:

- declare object type and schema version;
- use stable enum values;
- preserve unknown/unsupported states;
- include concrete source references;
- avoid machine-local paths;
- bind to input hashes where freshness matters;
- fail on ambiguous or malformed input.

## Agent-Native CLI Policy

State-returning commands:

- support stable JSON;
- remain machine-readable when piped or non-interactive;
- preserve status, evidence refs, typed errors, and exit code in quiet/compact
  modes;
- return nonzero for unimplemented behavior;
- never use a generic pass envelope as a placeholder.

Help and docs must not advertise a command as working before its end-to-end
acceptance passes.

## Migration Corpus Policy

Every benchmark fixture includes:

- fixture and change IDs;
- pinned source/target refs and digests;
- license, attribution, and redistribution posture;
- official SDK package/version;
- annotated impacted and unaffected call sites;
- expected patch;
- expected verification stage and outcome;
- unsupported/negative classification where applicable.

Visible development fixtures and held-out scoring fixtures remain separate.
Ground truth is versioned before scoring. Public fixtures demonstrate
engineering behavior only.

Only a non-resolving opaque holdout manifest containing opaque case IDs,
provenance class and license posture, a frozen suite commitment, and encrypted
or HMAC artifact commitments may be committed. An independent holdout owner
provisions and freezes `LUMYN_HOLDOUT_ROOT`. Source URLs, repository or package
identifiers, plaintext content digests, held-out repositories, inputs, answer
keys, expected patches or labels, and raw traces remain there, outside every
task-executor mount, prompt, and environment.

M1 uses provision-mode holdout policy and `holdout_provisioning_required` with
an opaque private namespace and `hmac-sha256` algorithm; the plan never invents
a commitment before the suite exists or labels provisioning as hidden
evaluation of M1. M4, M6, and M7 use `holdout_evaluation_required` and
evaluate-mode policy that resolves the trusted M1 result. The evaluator and
shipping gate bind the exact referenced result bytes so replacement or replay
invalidates the evaluation.

## API-Provider Change Packet Trust Policy

- API-provider change packets are signed declarative data, never executable
  scripts.
- Initial enrollment uses a provider-enrollment bundle and expected root
  fingerprint obtained through a separately authenticated provider
  admin/security channel. Invitation-supplied root material is never
  self-authenticating.
- `provider enroll` and `campaign accept` accept no repository argument and run
  from a neutral directory with the checkout unavailable. They persist only to
  the configured consumer-private root and create no authority.
- `authorization issue` and `authorization revoke` are explicit,
  consumer-signed private-state actions. An authorization request is never a
  grant, and issuance never performs a granted side effect.
- Canonical signing bytes and the immutable packet digest are versioned.
- The consumer pins the provider trust root and verifies provider-organization
  ownership of the named npm package.
- Issuer key, issue time, audience, expiry, rotation, revocation, withdrawal,
  and replay checks must all pass before repository analysis or mutation.
- Current lifecycle state comes from either a signed offline provider-status
  snapshot inside the enrollment policy's maximum age or an exact endpoint
  read under `provider_trust_status_read`. The request carries no repository or
  consumer data. Missing, stale, replayed, unsigned, wrong-endpoint, or
  undeclared status blocks.
- `published` means immutable for the authorized audience; it does not mean the
  packet or prerelease migration is public.
- Unknown issuer/package binding, stale or replayed packet, invalid signature,
  revoked key, unsigned rotation, unconfirmed first-pin fingerprint,
  withdrawn/superseded packet, or executable content fails closed. Normal
  rotation is signed by the active root; emergency recovery requires explicit
  re-enrollment and invalidates open approvals.
- Valid and invalid fixtures cover every trust and lifecycle boundary.

## TypeScript Impact Policy

- Use a parser/AST or comparably structured representation.
- Select and canonicalize one package/read root explicitly.
- Resolve real paths before reading. Source files, manifests, lockfiles,
  resolved modules, `tsconfig` `extends`, and TypeScript project references
  must remain inside the selected root.
- Reject path traversal, symlink escape, out-of-root `tsconfig` dependencies,
  ambiguous roots, and multiple package roots before analysis.
- Detect direct imports, aliases, and one-hop wrapper uncertainty.
- Report dynamic/reflection use as uncertain.
- Exclude generated, vendored, minified, cache, and build output.
- Report package-manager and lockfile posture.
- Score precision and recall separately by supported class.
- Never label uncertain scope as unaffected.

## Patch And Filesystem Policy

- No patch before exact plan approval.
- Revalidate the current packet bytes, digest, provider trust root and package
  binding, lifecycle, audience, expiry, rotation, revocation, withdrawal,
  supersession, and replay state immediately before every file or lockfile
  write. Never cache the plan-time trust decision across a write boundary.
- Use an isolated worktree or equivalent disposable workspace.
- Resolve and validate real paths before writes.
- Enforce allowed/forbidden paths and diff budgets.
- Reject symlink/path traversal escape.
- Map each edit to change and recipe IDs.
- Preserve deterministic output for pinned inputs.
- Mutate `package-lock.json` only with exact Node/npm versions, a pinned
  registry or offline snapshot, recorded package-integrity inputs, lifecycle
  scripts disabled, and a bound toolchain digest.
- Do not execute arbitrary provider scripts.
- Do not infer business values.
- Record rollback and deletion checks.

## Command Execution Policy

Repository commands are untrusted:

- exact command allowlist;
- explicit working directory;
- exact read-only and writable mounts;
- neutral home and temp roots;
- explicit executable/toolchain roots;
- timeout and output budget;
- no network by default;
- package-registry access requires a separate `package_registry_read` grant
  constrained to the approved registry or snapshot and package set;
- dependency lifecycle scripts disabled by default;
- no lockfile mutation through an ambient registry, floating Node/npm version,
  or missing integrity input;
- sanitized environment classes and no ambient secrets;
- no host home, SSH/GPG/cloud credential stores, keychain, OS credentials,
  agent/Docker/unrelated local-service sockets, or extra inherited file
  descriptors;
- child processes inherit the same mount, environment, socket, descriptor,
  credential, network, and resource boundary;
- a supported fail-closed isolation backend is mandatory;
- sandbox credentials never exposed to build/test stages;
- pre- and post-patch results kept separate;
- large logs stored by artifact ref and hash.

## Proof-Of-Behavior Policy

Product verification state uses exactly:

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

`workflow_contract_replay_passed` means an independent contract or cassette
replay and cannot exceed `repo_verified`. A `workflow_verified_replay`,
`workflow_verified_mock`, or `workflow_verified_sandbox` result requires an
approved entrypoint executed from the exact patched repository head plus
observed interaction and outcome evidence in the named environment. A result
copied from the base commit or another head, or an independent replay that did
not execute the patched repository, is not causal patched-head proof. Missing
causal execution or any boundary, cleanup, redaction, freshness, or
evidence-integrity failure blocks a workflow-verified label.

Artifact syntax and grounded source references remain engineering evidence
strengths, not alternate product verification labels. Consumer/provider
outcomes from the real pilot are user-visible evidence, not a substitute for
repository or workflow verification. No lower boundary closes a higher-level
item; workflow and pilot closure require proof-of-behavior scorecards.

## Redaction And Evidence Budgets

- Redact before persistence and before sharing.
- Redaction uncertainty blocks the artifact.
- Provider-visible attestation has an explicit field allowlist.
- Raw source, diffs, logs, traces, prompts, responses, and credentials are
  private by default and persist only in the approved private state root
  outside the checkout and public source repository.
- Private records carry authorization-bound TTL and are deleted on expiry or
  revocation automatically on creation, read, process startup, and the next
  run; cleanup produces a deletion receipt or orphan report.
- `lumyn artifacts gc` retries partial deletion and reports unresolved orphans.
  It cannot extend TTL, revive revoked authority, or rewrite historical
  closure evidence.
- Provider-visible does not mean public. Public release requires separate
  explicit consent and is limited to redacted aggregates and evidence hashes
  under `.factory/artifacts/pilot/lumyn-migration-mvp/public/`.
- Large output is referenced by path, digest, row/event count, and truncation
  metadata inside private evidence; shareable artifacts use an opaque
  identifier and digest rather than a private or machine-local path.
- Machine-local paths are removed from shareable artifacts.

## Capability Grants

Live tasks use exact, task-scoped grants:

- `customer_repo_read`
- `customer_repo_write`
- `command_execution`
- `provider_trust_status_read`
- `package_registry_read`
- `sandbox_network`
- `sandbox_credential`
- `sandbox_request_disclosure`
- `github_branch_write`
- `github_pr_write`
- `campaign_receipt`
- `provider_attestation`
- `artifact_retention`
- `artifact_deletion`

`customer_repo_write` authorizes bounded local workspace mutation, not a remote
branch. `provider_trust_status_read` authorizes only the pinned status endpoint,
request shape, response budget, maximum age, and expiry; a current signed
offline snapshot needs no egress. `package_registry_read` authorizes only the
named registry/snapshot and package set, not general network access.
`sandbox_network` authorizes only the named non-production endpoint and
operation allowlist; `sandbox_credential` independently authorizes the
non-production credential class and scopes.
`sandbox_request_disclosure`
separately names allowed payload classes, synthetic or approved non-sensitive
test data, and provider logging, retention, and deletion terms. Production
customer data, PII, credentials, and secrets are never eligible payloads.
`github_branch_write` and `github_pr_write` are independent and neither
authorizes default-branch write or merge. `campaign_receipt` authorizes only
the signed minimal sponsored-program meter and requires a
provider-authenticated consumer signer binding, invitation/packet and
eligible-repository unit binding, opaque IDs, consent-policy and key-binding
digests, pinned endpoint or offline exchange, audience, expiry, nonce,
signature, provider acknowledgement key, provider-signed acknowledgement,
one-unit cardinality, idempotent same-digest retry, and conflicting-unit/replay
rejection.
`provider_attestation` authorizes only its richer field allowlist and does not
authorize public disclosure. It is optional for M9/M10 richer reporting and
cannot be made a prerequisite for an otherwise authorized M9 draft PR.
`artifact_retention` and `artifact_deletion` independently name the exact
artifact classes, storage boundary, TTL, expiry/revocation triggers, deletion
scope, receipt owner, and orphan route. Every data-producing or disclosure
grant references both.

Grants name target, scope, expiry, revocation, evidence, and failure behavior.
Wildcard customer-repository, registry, sandbox, request-disclosure, GitHub,
provider-status, campaign-receipt, provider-attestation, retention, and
deletion grants are invalid.
Factory credential scopes and network allowlists also reject semantic
wildcards (`all`, `any`, `default`), case-insensitive duplicates, wildcard or
unspecified hosts, and CIDR-wide access.

Product-signal tasks validate the actual private attestation and its
aggregate/hash-only public manifest with a task-specific validator. A passing
`--self-test` is not outcome evidence. An independent `evidence-attestor`
record is required before `DISC` or `PILOT` closure.

Provider export and public commit are irreversible disclosure boundaries.
Revocation blocks future sharing and deletes Lumyn-controlled private copies;
it cannot erase provider records, Git history, clones, or caches. Tests and
rollback evidence must not claim otherwise.

The current repository and pilot distribution are not described as OSS.
Design-partner artifacts require explicit terms, a security/support route,
signed provenance, checksums, and install-integrity instructions. Public
OSS/self-serve adds an approved license plus security, contribution, support,
and vulnerability-response policy gates.

M1, M4, and M6 require an independent `holdout_result`; M7 requires both
`holdout_result` and `trace_grade_report`; M10 requires an
`attestation_record`. Each is lifecycle-owned and must pass before shipping.
M1's result binds independent provisioning and the frozen suite commitment;
M10's attestation reviews private campaign calculations without receiving or
reusing the benchmark holdout root.

## Release Integrity

Primary design-partner distribution is an explicitly licensed,
integrity-signed local or CI binary/source package. Public OSS/self-serve and
Homebrew wait for the separate approved license, security, contribution,
support, vulnerability-response, and release-integrity gate. Release work
requires version, changelog, install, provenance, artifact-integrity, supported
platform, and UAT evidence. Planned commands are not release claims.

## Provenance Evidence

- Task validation:
  `.factory/artifacts/task-runs/<task_id>/validation-report.json`
- Work proof:
  `.factory/artifacts/task-runs/<task_id>/work-proof-marker.json`
- Independent lifecycle evidence:
  `.factory/artifacts/lifecycle-evidence/<task_id>/`
- PR lifecycle:
  `.factory/artifacts/pr-lifecycle/<work_item_id>/pr-lifecycle-report.json`
- Active acceptance:
  `.factory/artifacts/prd-to-plan/lumyn-migration-mvp/acceptance-ledger.json`
- Active closure:
  `.factory/artifacts/prd-to-plan/lumyn-migration-mvp/scope-closure-map.json`
- Historical plan:
  `.factory/artifacts/prd-to-plan/lumyn-mvp/`
- Historical pilot:
  `.factory/artifacts/pilot/lumyn-mvp-slice/`
- Public migration-pilot summaries:
  `.factory/artifacts/pilot/lumyn-migration-mvp/public/` (consented, redacted
  aggregates and evidence hashes only)
- Remote controls:
  `.factory/artifacts/repo-controls/main-branch-protection.json`

These repository paths contain only source-safe Factory control, lifecycle, or
public aggregate/hash evidence. Consumer-private impact, plan, patch,
verification, runtime, and identifiable pilot records live in the configured
external private state root and are referenced here only by opaque identifier
and digest. Evidence records command, status, artifact refs, hashes where
applicable, and skipped-command reasons.
