# Lumyn Workflow Contract

Version: 2.0
Status: Normative

## Work Signal

Lumyn accepts governed work from:

- `docs/product/prd.md`;
- `docs/product/plan.md`;
- active Factory task packets under
  `.factory/artifacts/prd-to-plan/lumyn-migration-mvp/`;
- governed post-PRD findings under `.factory/artifacts/post-prd/`;
- GitHub issues or pull requests when they map to the active product contract.

The prior `.factory/artifacts/prd-to-plan/lumyn-mvp/` package and its task-run,
pilot, and PR-lifecycle evidence are historical. They are not an active work
signal.

## Product Delivery Flow

The product workflow must preserve these state boundaries:

```text
standard provider campaign kit
-> configured signing and immutable packet publication
-> signed provider packet and signed invitation
-> trust-root/package-binding/lifecycle validation
-> packet/canary validation
-> provider invitation
-> independently authenticated provider enrollment and fingerprint confirmation
-> consumer campaign acceptance and no-authority authorization request
-> current provider-status resolution from a signed snapshot or exact authorized read
-> explicit consumer-signed repository authorization issuance
-> [consumer-signed minimal connection receipt, provider-signed acknowledgement, optional]
-> read-only impact
-> no-write migration plan
-> consumer local-write/command/required-network-and-credential approval
-> current packet-trust and product-authority revalidation at write boundary
-> isolated patch
-> repository verification
-> [approved sandbox verification, independently optional]
-> consumer remote-branch authorization
-> consumer PR authorization
-> draft PR
-> consumer review and merge
-> [separately consented richer provider attestation, optional]
```

No downstream state implies an earlier missing authorization. A provider
invitation is not repository consent. An approved migration plan is not PR
authorization. A passing repository check is not workflow proof.

An invitation cannot authenticate its own provider root. Initial
`provider enroll` requires a separately obtained enrollment bundle and expected
fingerprint from an authenticated provider admin/security channel. Enrollment
and `campaign accept` write only to external consumer-private state and leave
the checkout and Git index unchanged.

## Normal Factory Chain

1. `scout-context`
2. `execution-compiler`
3. `task-executor`
4. `validation-gate`
5. `code-review` when required by risk
6. `holdout-evaluator` when selected by policy
7. `trace-grader` when selected by policy
8. `evidence-attestor` when selected by policy
9. `commit-push`
10. `post-merge-monitor`
11. `repair-feedback` when validation, review, shipping, or closure fails

The independent evaluators are external or human-operated lifecycle gates.
They run after implementation validation and before shipping, and each writes
a task-bound, current-work-proof, passing artifact in the trusted Factory
evidence root. The implementation worker cannot self-grade or self-attest.

The active PRD-derived ledger contains exactly 62 item-level acceptance
controls. Acceptance item IDs remain the closure source; delivery slices are
coverage and sequencing metadata only.

## Approval Gates

### Product Planning

- Plan approval is required before implementation work begins.
- Runtime, distribution, trust-boundary, artifact namespace, active-plan path,
  or initial-segment changes require aligned PRD, plan, ADR, Factory profile,
  task packets, validation contract, acceptance ledger, mapping, and closure
  updates.
- M2.5 must close `DISC-001` and `DISC-002` with direct external evidence for
  one qualified API Provider commitment, a supported-class canary, and a
  frozen cohort of five distinct eligible repositories across at least three
  Consumer Organizations before M5 migration-plan implementation begins. The
  same provider must document either a second consequential
  migration/deprecation within 12 months or an annual connected-repository
  opportunity covering at least 20 named managed integrations. Qualification
  also requires an operational signed provider-status channel, a pinned
  provider receipt-acknowledgement key plus endpoint or offline exchange, and a
  supported OS/architecture with enforceable host isolation for every candidate
  environment.
- The M2.5 and M10 evidence validators must validate the actual private
  attestation plus its committed aggregate/hash-only public manifest.
  `--self-test` proves validator behavior only and cannot close a product
  signal. An independent `evidence-attestor` record is required.
- Before M2.5 collects, stores, or discloses identifiable external evidence, a
  narrow manual privacy/legal preflight under the Factory implementation
  worker's task-scoped `approval` must name allowed private fields,
  participant consent, the approved external private root, TTL,
  expiry/revocation deletion, deletion-receipt and orphan ownership, and the
  minimal connection-receipt, authenticated consumer signer binding,
  provider-signed acknowledgement/cardinality policy, and separately consented
  aggregate/hash-only public fields. Consent states that provider export and
  public commit cannot be recalled. It is an evidence-handling preflight, not
  Lumyn runtime product authority. The active M2.5 approval must cite that
  preflight and exactly match its canonical scope digest; generic approval, an
  older digest, or a changed field/storage/TTL policy fails closed.
- M2.5 canary inspection qualifies the opportunity but cannot satisfy the M3
  signature, trust-root, provider/package-binding, lifecycle, replay, and
  declarative-only runtime gate.
- Pilot distribution must carry explicit terms, a security/support route,
  signed provenance, checksums, and install-integrity instructions. Do not call
  the current repository or design-partner package OSS.

### Consumer Authorization Issuance

`campaign accept` produces only a reviewable no-authority request. The
activation process accepts no repository argument and runs with the checkout
unavailable.

`authorization issue` requires an exact grant manifest, consumer-controlled
configured signer, campaign/repository binding, issue time, expiry, nonce, and
revocation route. It writes the private bundle only and performs none of the
granted actions. `authorization revoke` records a signed grant or bundle
revocation. Every later side effect independently revalidates the current
signed state.

### Provider Trust Status

Every side effect consumes either a signed offline provider-status snapshot
inside the enrollment policy's maximum age or an exact endpoint read under
`provider_trust_status_read`. The grant names provider/package, endpoint,
request shape, response budget, maximum age, expiry, and proof that no
repository or consumer data is transmitted. The signed response binds packet,
provider, package, campaign, audience, issue time, expiry, rotation epoch,
revocation/withdrawal state, and anti-replay nonce. Missing, stale, replayed,
unsigned, wrong-endpoint, or undeclared status access blocks.

The same independently authenticated provider enrollment pins the
receipt-acknowledgement signing key and allowed online/offline exchange
classes. An invitation may narrow those values for a campaign and unit but
cannot introduce a new key or destination.

### Customer Repository Read

`customer_repo_read` must name:

- campaign and packet;
- repository;
- readable paths and exclusions;
- token/approval expiry;
- retention scope;
- deletion scope and receipt owner;
- consumer evidence owner.

Impact analysis cannot start without this grant.

### Customer Repository Write

`customer_repo_write` must additionally name:

- approved plan digest;
- base commit;
- writable paths;
- file/line/diff budgets;
- isolated workspace;
- rollback posture.

Read authorization never implies local write authorization. Local write
authorization never implies remote branch write.

### Command Execution

`command_execution` must name:

- exact commands;
- working directory;
- exact read-only and writable mounts;
- neutral home and temp roots;
- executable and toolchain roots;
- timeout and output budgets;
- dependency lifecycle-script posture;
- network posture;
- environment variables by class, never secret values;
- local socket and inherited-file-descriptor policy;
- process-tree limits and mandatory child-process inheritance;
- OS credential access, which defaults to denied;
- the supported fail-closed isolation backend.

Repository commands default to no network, no dependency lifecycle scripts,
no host home or credential stores, no SSH/GPG/cloud/keychain access, no
agent/Docker/unrelated local-service sockets, and no inherited descriptors
beyond standard streams. If the backend cannot enforce the exact boundary, the
command does not run.

### Package Registry

`package_registry_read` is independent from command execution and must name:

- exact Node and npm versions and toolchain digest;
- exact package/version and package-integrity inputs;
- registry endpoint or immutable snapshot;
- read-only network allowlist and budget;
- disabled lifecycle scripts unless separately approved;
- expiry and evidence.

Automatic `package-lock.json` mutation is prohibited when any of these pins or
the registry-network grant is missing.

### Provider Sandbox

`sandbox_request_disclosure` must separately name:

- exact transmitted payload classes;
- synthetic or approved non-sensitive test-data classes;
- prohibited production customer data, PII, credentials, and secrets;
- what the provider may log;
- provider retention and deletion terms;
- consumer-approved evidence of disclosure and deletion.

`sandbox_network` must separately name:

- API-provider identity and environment;
- network allowlist;
- operation allowlist;
- resource namespace;
- request and write budgets;
- settle/retry posture;
- cleanup and orphan evidence;
- approval evidence and expiry.

`sandbox_credential` must separately name the non-production credential class,
scopes, isolated injection stage, expiry, revocation, and evidence.

Sandbox credentials remain isolated from build and repository-test commands.
Generic network, credential, or payload approval does not satisfy any gate.

All grants in this section are private, schema-backed Lumyn product
authorizations. `.factory/factoryd.json` retains Factory's closed
`approval`, `credentials`, and `network` vocabulary for the implementation
worker. Its approval evidence may cite an opaque validated product bundle, but
no Factory grant substitutes for or broadens a Lumyn product authorization.
Factory selection and dispatch do not evaluate Lumyn product authority. Lumyn
must revalidate the current packet trust state and exact applicable product
bundle inside the product immediately before every local, sandbox, GitHub, or
attestation side effect and every retry; no earlier decision may be cached.
`lumyn authorization validate` supplies diagnostic and closure proof, not the
authoritative side-effect gate.

### GitHub Remote Branch

`github_branch_write` must name:

- repository;
- authorized non-default branch namespace;
- base commit;
- installation-token expiry;
- rollback and idempotency posture.

It grants no authority to create or update a pull request.

### GitHub Draft PR

`github_pr_write` must name:

- repository;
- authorized remote branch;
- base branch;
- draft-only posture;
- installation-token expiry;
- idempotency key;
- approved plan/evidence refs.

It grants no authority to create the remote branch. The MVP never writes the
default branch or auto-merges.

### Provider Attestation

`provider_attestation` lists the exact fields the consumer permits Lumyn to
share plus their retention and deletion terms. Raw source, diffs, logs, traces,
prompts, responses, credentials, and secret values are excluded by default.
It is optional for M9 and M10 reporting actions; absence of this grant must
silence reporting, not block an otherwise authorized draft PR.

### Campaign Connection Receipt

`campaign_receipt` authorizes only the minimal signed sponsored-program meter:
campaign and packet digest, eligible-repository unit and invitation nonces,
opaque per-campaign organization/repository IDs, connection event/time,
tool/schema versions, consent-policy and provider-signed consumer-key-binding
digests, consumer issuer public-key fingerprint, audience, expiry, anti-replay
nonce, and signature. Before acknowledgement, the provider authenticates the
consumer organization through an existing account-admin or documented
out-of-band ownership channel and signs the receipt-key binding.

The invitation pins the provider acknowledgement key and either the exact
receipt endpoint or offline exchange contract. Online submission sends only
the canonical receipt under the exact grant; offline submission uses a bounded
export/import bundle. The provider-signed acknowledgement binds receipt digest,
invitation unit, key-binding digest, decision, time, and deduplication key.
One invitation unit can acknowledge only one opaque repository; an identical
digest retry is idempotent, while conflicting reuse blocks. The consumer
verifies and imports the acknowledgement before the unit is treated as
connected. Declining this grant leaves private local use and draft-PR delivery
available but excludes the repository from sponsored-program counts.

### Artifact Retention And Deletion

`artifact_retention` names exact artifact classes, their private or disclosed
storage boundary, TTL, expiry behavior, and evidence owner.

`artifact_deletion` independently names revocation/expiry triggers, deletion
scope, receipt owner, retry posture, and orphan-report route. A data-producing
or disclosure grant must reference both authorities; neither is implied by
repository, command, sandbox, GitHub, or provider-attestation approval.

The private-artifact owner enforces these grants automatically on creation,
read, process startup, and the next run. `lumyn artifacts gc` retries partial
deletion and reports unresolved orphans; it cannot extend TTL, revive revoked
authority, or rewrite historical closure.

Provider export and public commit are irreversible. Revocation blocks future
sharing and deletes Lumyn-controlled private copies; it cannot erase provider
records, Git history, clones, or caches. Consent and rollback evidence must say
so explicitly.

### Merge

- Merge approval remains with the Consumer Maintainer for product-generated
  migration PRs.
- Merge approval for Lumyn's own repository follows the Factory lifecycle.

## Artifact Rules

- Product source of truth: `docs/product/prd.md`.
- Human task plan: `docs/product/plan.md`.
- Active compiled plan:
  `.factory/artifacts/prd-to-plan/lumyn-migration-mvp/`.
- Historical compiled plan:
  `.factory/artifacts/prd-to-plan/lumyn-mvp/`.
- Task evidence: `.factory/artifacts/task-runs/<task_id>/`.
- Independent lifecycle evidence:
  `.factory/artifacts/lifecycle-evidence/<task_id>/`; this namespace is
  lifecycle-owned and forbidden to `task-executor`.
- PR lifecycle:
  `.factory/artifacts/pr-lifecycle/<work_item_id>/pr-lifecycle-report.json`.
- Factory scratch: `.factory/tmp/`.
- Daemon state: `.factoryd/`.
- Retained synthetic or licensed project fixtures: `workflows/`, `cassettes/`,
  and `runs/`; these roots do not authorize consumer-private runtime storage.
- Planned migration artifacts: `changes/`, `authorizations/`, `impacts/`, and
  `migrations/`, introduced only as schema-backed artifact types.
- Consumer-private runtime artifacts: an explicitly configured,
  non-committable root outside the consumer checkout and every public source
  repository. It holds private workflow, cassette, run, change, authorization,
  impact, migration, and PR-result instances.
- Public pilot evidence:
  `.factory/artifacts/pilot/lumyn-migration-mvp/public/`, containing only
  consented redacted aggregates and evidence hashes.
- Holdout inputs and answer material: evaluator-controlled
  `LUMYN_HOLDOUT_ROOT` outside the checkout and every task-executor mount,
  prompt, and environment. An independent holdout owner provisions and freezes
  it. Only a non-resolving manifest of opaque IDs, provenance class and license
  posture, a frozen suite commitment, and encrypted or HMAC commitments may be
  committed; source locations, package identifiers, and plaintext content
  digests may not.
- M1 provision mode declares only the opaque namespace and keyed-commitment
  algorithm. Its trusted result creates the suite ref and commitment;
  M4/M6/M7 evaluate mode resolves and byte-binds that result before scoring.

Factory refs are repo-relative; private runtime refs are relative to the
configured private root. Large output is cited by hash and stable artifact
reference. Private artifacts enforce authorization TTL and deletion on expiry
or revocation, with deletion receipts or orphan reports. Consumer-private,
provider-visible, and public artifacts are separate. Provider-visible
attestations contain only exact consumer-consented fields; public evidence is
aggregate/hash-only. Provider/public copies are non-recallable after export or
commit; private deletion evidence must never claim otherwise.

Historical evidence is immutable and may prove only its recorded semantics.

## Bootstrap Validation Lanes

- Fast: `make lint-fast`, `make test-fast`
- Coverage: `make test-coverage`
- Contract: `make test-contracts`
- Full: `make prepush-full`
- Risk: GitHub Actions `CodeQL analyze`
- Acceptance: active item-level ledger and scope-closure map
- Cross-system: task-scoped approved sandbox or GitHub integration

## PR Lifecycle Baseline

- Local validation gate: `make prepush-full`.
- GitHub validation check: `validate`.
- Security scanner: `CodeQL analyze`.
- Required-check manifest: `.github/required-checks.json`.
- Owner-review policy: `.github/CODEOWNERS`.
- Action-ref posture: `.github/action-ref-exceptions.yaml`.
- Structured review: `code-review` for trust, credential, external-call,
  patch, GitHub, schema, policy, or high-risk behavior.
- Independent lifecycle evidence: `holdout_result`, `trace_grade_report`, and
  `attestation_record` when selected by task policy, all verified before
  `commit-push`. Each artifact must schema-validate and bind the exact task,
  work item, lifecycle run, current validation run, candidate digest, and work
  proof marker digest, with passing status and independent worker provenance.
- Task-level review type and reviewer class are authoritative and must exactly
  match the inherited validation contract. Shipping rejects a missing,
  mismatched, stale, or implementation-authored review.
- Shipping evidence: `commit-push` ship packet.
- Post-merge evidence: `post-merge-monitor`.
- Unavailable lifecycle surfaces require an explicit approved exception.

Passive Codex review settle is required before merge. Green CI alone is not merge-ready.
Do not merge manually through `gh pr merge`, the GitHub UI, or a
connector before the latest-head terminal Codex review signal. A merge without
that evidence is a process escape and requires recorded follow-up.

GitHub `main` must remain protected by branch protection and the
`protect-main-from-direct-push` ruleset. Audit live controls with:

```bash
make audit-remote-protection
```

## Validation And Proof Rules

- Capture repository baseline before mutation.
- Keep impact, patch, verification, delivery, permission, and residual-risk
  axes separate.
- Use only the canonical successful-verification labels `static_verified`,
  `repo_verified`, `workflow_contract_replay_passed`,
  `workflow_verified_replay`, `workflow_verified_mock`, and
  `workflow_verified_sandbox`.
- A `workflow_verified_*` result requires an approved entrypoint executed from
  the exact patched repository head plus observed interaction and outcome
  evidence in that named environment. Independent contract or cassette replay
  is `workflow_contract_replay_passed` and cannot exceed `repo_verified`.
- Production access is out of scope.
- Bind evidence to packet, plan, base/head commits, patch, environment, and
  artifact hashes.
- Stale evidence cannot close acceptance.
- Boundary, cleanup, orphan, redaction, authorization, or proof failures block
  stronger verification labels.
- Unimplemented commands must return typed nonzero results.
- Public fixtures prove engineering behavior, not demand.

## Factory And factoryd

Safe attended path:

```bash
export FACTORY_REPO=../factory
factoryd doctor --config .factory/factoryd.example.json --repo lumyn --json
factoryd run --config .factory/factoryd.example.json --repo lumyn --dry-run --json
```

Use `.factory/factoryd.autoship.example.json` only when:

- branch protection and required checks are verified;
- task paths, commands, evidence, and stop conditions are bounded;
- passive review, merge, post-merge, and semantic item closure are enforced;
- every required customer-repository, command, GitHub, sandbox, credential,
  network, reporting, retention, and deletion capability grant is task-scoped
  and active;
- package-registry network, sandbox request disclosure, sandbox access,
  remote-branch write, and PR write remain separately authorized.

Product workers may emit task-scoped evidence. They must not mutate active
PRD-derived control truth to bypass a gate.

## Post-PRD Findings

Material audit, review, pilot, or recommendation findings are saved as
structured repo-local finding lists and ingested before implementation.
Generated post-PRD artifacts become the execution contract for the follow-up
mission.

Do not alter the PRD from a task-run worker. A human explicitly promotes a
finding into product scope and then rebaselines all affected planning artifacts.

## Stop Conditions

Stop and request a human decision if:

- provider and consumer authority are ambiguous;
- a provider packet can execute code, lacks a valid consumer-pinned trust root
  or provider-to-package binding, or fails key, timestamp, audience, expiry,
  rotation, revocation, withdrawal, or replay checks;
- an exact capability grant is missing or expired;
- read-only analysis would mutate state;
- the approved plan, repository base, or packet digest has drifted;
- current packet trust or exact product authority cannot be revalidated
  immediately before a local, sandbox, GitHub, or attestation side effect;
- a local write, command, registry-network, sandbox-network, credential,
  payload-disclosure, remote-branch, PR, reporting, retention, or deletion
  boundary is unclear;
- required business input is missing;
- a task requires production access;
- a repository command needs unapproved network or lifecycle scripts;
- a repository command lacks enforceable host isolation for mounts,
  environment, credentials, sockets, descriptors, or child processes;
- packet freshness lacks a current signed offline snapshot or exact authorized
  provider-status read;
- lockfile mutation lacks exact Node/npm versions, registry/snapshot,
  package-integrity inputs, toolchain digest, or registry-network approval;
- sandbox execution lacks payload/test-data, provider logging, retention,
  deletion, cleanup, or orphan consent;
- redaction or evidence freshness is uncertain;
- private runtime evidence would enter a checkout or public repository, or its
  automatic TTL/revocation deletion, receipt ownership, or orphan recovery
  cannot be enforced;
- evidence would be labelled above its proven boundary;
- provider-visible data exceeds consumer consent;
- a sponsored connection is counted without an acknowledged valid signed
  minimal receipt;
- deletion or rollback evidence claims to recall an already exported or
  publicly committed copy;
- remote branch write or PR write is inferred from the other;
- default-branch write or auto-merge is requested;
- M5 or later migration-plan work would proceed before M2.5 closes
  `DISC-001` and `DISC-002`;
- an implementation worker can read held-out answer material or write
  lifecycle-owned review, holdout, trace-grade, attestation, shipping, or PR
  lifecycle evidence;
- M9 draft-PR delivery is made dependent on M8 sandbox access or optional
  provider reporting;
- the active PRD, plan, Factory profile, task packets, acceptance ledger,
  mapping, validation contract, and scope closure disagree;
- required validation, coverage, CodeQL, review, shipping, post-merge, or
  item-level evidence is missing without an approved exception.
