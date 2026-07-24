# PLAN Lumyn Migration MVP: Provider-Sponsored Verified API Migrations

**Date:** 2026-07-23
**Source of truth:** `docs/product/prd.md`
**Scope:** Rebaseline Lumyn around one provider-sponsored, customer-controlled
TypeScript/Node migration campaign and deliver evidence-backed draft PRs for
authorized GitHub repositories.

---

## Global Decisions (Locked)

1. The API Provider is the economic buyer and campaign sponsor.
2. The API Consumer Organization owns code, credentials, execution, disclosure,
   and merge authority.
3. Provider payment never grants repository access.
4. Consumer execution is local or consumer-CI by default.
5. Provider change packets are signed, immutable for their authorized audience
   after publication, and declarative. The consumer pins a provider trust root
   and verified provider-to-package ownership binding; arbitrary provider
   scripts are prohibited.
6. The implementation core remains Go `1.26.5` at module
   `github.com/Clyra-AI/lumyn`.
7. The first target is one official TypeScript/Node npm SDK in one explicitly
   selected package root.
8. Only three deterministic change classes are patchable: method/operation
   rename, request-property rename/relocation, and response-property
   rename/relocation.
9. Impact analysis is read-only. Migration planning precedes write approval.
10. Patch execution uses an isolated workspace and explicit path/diff budgets.
11. Generic model-assisted patching and live agent-evaluation product surfaces
    are outside the MVP.
12. Repository tests have network disabled by default. Registry network,
    lifecycle scripts, sandbox network, transmitted sandbox payload classes,
    remote branch write, and PR write require independent approval.
13. Automatic `package-lock.json` mutation pins the exact Node/npm toolchain,
    registry or immutable snapshot, package integrity inputs, and toolchain
    digest.
14. Private runtime and pilot evidence lives outside the checkout and public
    source repository. TTL and revocation deletion are enforced automatically
    on creation, read, startup, and the next run; cleanup emits receipts or
    orphan reports, with `lumyn artifacts gc` as retry-only operator recovery.
15. Lumyn opens draft PRs only and never writes to the default branch or
    auto-merges.
16. Provider-visible and public disclosure are separate. A provider
    attestation contains only the exact fields allowed by that consumer; public
    evidence requires separate consent and is aggregate/hash-only.
17. Public fixtures prove engineering behavior, not commercial demand.
18. Historical Factory artifacts under
    `.factory/artifacts/prd-to-plan/lumyn-mvp/` remain immutable and
    non-active.
19. The active plan and closure source live under
    `.factory/artifacts/prd-to-plan/lumyn-migration-mvp/`.
20. Provider activation uses a standard campaign kit, configured signer,
    separately authenticated provider-enrollment bundle and out-of-band
    confirmed fingerprint, signed invitation, consumer trust-enrollment
    receipt, and no-authority authorization request. Invitations cannot
    self-authenticate their root.
21. A synthetic offline end-to-end canary must pass before any provider
    demonstration; generic unimplemented `demo` behavior is never an
    acceptable substitute.
22. An independent holdout owner provisions and freezes held-out inputs and
    answer keys in an evaluator-controlled private root unavailable to
    implementation workers. Only a non-resolving manifest and aggregate,
    commitment-bound lifecycle results are committable.
23. Draft-PR delivery requires M7 repository/workflow proof, but neither live
    sandbox access nor provider reporting. Sandbox evidence and provider
    attestation remain independently authorized capabilities.
24. M10 may use exact-patched-head mock proof or separately authorized sandbox
    proof per repository. M8 and the three sandbox grants are conditional on an
    actual sandbox action, not campaign-global prerequisites.
25. M0-M4 engineering may start from pinned, license-compatible public API
    docs, OpenAPI descriptions, SDK releases, migration guides, and synthetic
    fixtures. Those inputs are untrusted engineering evidence, not a signed
    provider packet, provider endorsement, sandbox authority, customer cohort,
    or demand proof; M5 and the pilot remain gated by M2.5.
26. `campaign accept` creates only a no-authority request. A consumer-controlled
    signer must explicitly issue or revoke the private authorization bundle;
    issuance never performs the authorized side effect.
27. Packet freshness comes from either a signed offline provider-status
    snapshot inside the pinned maximum age or an exact-endpoint read under
    `provider_trust_status_read`. Missing or undeclared freshness blocks.
28. Repository commands require a supported fail-closed host-isolation backend
    with exact mounts, sanitized environment, closed host credential/socket/
    descriptor surfaces, and inherited child-process restrictions.
29. A sponsored connection is counted only from a provider-acknowledged,
    consumer-signed minimal connection receipt. Richer provider reporting is a
    separate opt-in; private local use remains available without either.
30. Provider export and public commit are irreversible disclosures. Revocation
    stops future sharing and deletes Lumyn-controlled private copies but cannot
    recall recipient copies or public history. The current repository is not
    described as OSS; public OSS/self-serve distribution has a separate
    license, security, contribution, support, and release-integrity gate.

---

## Current Baseline (Observed)

Implemented:

- Go CLI/config/result/exit-code foundation.
- `lumyn init` and `lumyn check`.
- OpenAPI and local-doc parsing, fingerprints, structured refs, deprecation
  findings, and concrete source locations.
- Executable schemas for workflows, evidence, cassettes, traces, proof,
  boundaries, redaction, and command results.
- CI, coverage, CodeQL, branch-policy, CODEOWNERS, passive-review, Factory
  planning, commit/push, and post-merge governance.

Not implemented:

- provider change packets;
- provider enrollment and provider-status refresh;
- explicit authorization issuance and revocation;
- provider-authenticated connection-receipt issuance, acknowledgement, and
  import;
- fail-closed host command isolation;
- historical migration corpus;
- API/SDK semantic diffing;
- TypeScript repository impact analysis;
- migration planning or patching;
- repository verification orchestration;
- replay or live workflow verification runtime;
- evidence-backed PR delivery;
- richer provider campaign attestation;
- migration outcome ingestion.

Known correctness debt:

- `record`, `verify`, `trace`, `demo`, `share`, and `eval` are recognized by the
  command dispatcher even though they have no implementation and can return a
  generic pass result.
- current result contracts use bare `provider_metadata` for model-provider
  metadata and set eval-oriented values on non-eval commands.
- the old 88-item plan requires live agent evaluation and is incompatible with
  the new product.

Carried evidence is accepted only for the exact implemented foundation it
proves. Unstarted T4–T12 work from the historical plan is superseded, not
relabelled.

---

## Exit Criteria

The technical migration engine is complete when:

- the supported packet, authorization, impact, plan, patch, verification, and
  PR artifact contracts are schema-backed and fail closed;
- the fixed benchmark and two-sided activation path meet every `CORPUS`, `ACT`,
  `IMP`, `PATCH`, `VER`, `EVD`, and `PR` acceptance item;
- every held-out positive and negative workflow case receives its exact
  canonical label with zero false positives or false negatives;
- every `workflow_verified_*` result causally executes an approved entrypoint
  from the exact patched repository head;
- customer-private execution and provider-visible attestation are separated;
- the full local gate and required GitHub checks are green.

Migration-plan implementation beyond contracts and impact analysis is gated on
`DISC-001` and `DISC-002`. Commercial MVP completion additionally requires all
`PILOT` acceptance items, including actual paid evidence.

---

## Public API And Contract Map

| Surface | Status | Contract |
|---|---|---|
| `lumyn init` | Retained | Repo-local configuration initialization |
| `lumyn check` | Retained and reframed | Source and prerequisite validation |
| `lumyn campaign kit create` | New | Standard provider packet, canary, invitation, consent, and receipt kit |
| `lumyn change publish` | New | Configured-signing publication of canonical packet bytes |
| `lumyn campaign invite create` | New | Signed campaign/audience invitation bound to a published packet |
| `lumyn provider enroll` | New | Out-of-band-confirmed provider root/package enrollment outside the checkout |
| `lumyn campaign accept` | New | Invitation verification against an already pinned root and no-authority external-state authorization request |
| `lumyn trust refresh` | New | Import a signed provider-status snapshot or use an exact authorized status endpoint without repository or consumer disclosure |
| `lumyn authorization issue` | New | Consumer-signed creation of an exact private, time-bounded product-authority bundle; performs no granted action |
| `lumyn authorization revoke` | New | Consumer-signed revocation of one grant or the entire private bundle |
| `lumyn change validate` | New | Provider-change validation |
| `lumyn impact` | New | Read-only repository impact |
| `lumyn authorization validate` | New diagnostic/closure surface | Exact private product-authority bundle validation; live actions revalidate internally at the side-effect boundary |
| `lumyn campaign receipt issue/submit` | New optional sponsored-program surface | Issue the minimal receipt under an authenticated consumer key binding, then submit only to the pinned endpoint or bounded offline exchange |
| `lumyn campaign receipt acknowledge/ack import` | New provider/consumer exchange | Provider verifies, deduplicates, and signs the unit acknowledgement; consumer verifies and imports it |
| `lumyn migrate plan` | New | No-write migration plan |
| `lumyn migrate apply` | New | Approved isolated patch |
| `lumyn verify --migration` | New/reused semantics | Repository and workflow verification ladder |
| `lumyn canary run --offline` | New | Receipt-backed synthetic end-to-end flow through a local draft-PR preview |
| `lumyn trace` | New runtime over retained schema | Local evidence rendering |
| `lumyn artifacts gc [--dry-run]` | New recovery surface | Enforce private-artifact TTL/revocation deletion and retry or preview orphan cleanup |
| `lumyn pr create --draft` | New | Explicit draft-PR delivery |
| `lumyn.command_result` `1.0` | Compatibility surface | Existing envelope; legacy eval/model fields become non-applicable |
| migration artifact schemas | New | Versioned JSON Schema contracts |
| workflow/evidence schemas | Retained | Verification substrate |

Compatibility rules:

- Exit codes `0` through `9` remain stable.
- Exit code `6` is reserved and is not reassigned.
- Existing bare `provider_metadata` continues to mean model-provider metadata
  during the compatibility window.
- API-provider identity uses `api_provider_id` and `change_authority`.
- Persisted schema changes require versioning, valid/invalid fixtures, and
  migration notes.
- Programmatic failures use typed machine-readable errors with concrete source
  or artifact references.

---

## Docs And Distribution Readiness Baseline

The first screen of `README.md` must communicate:

```text
provider change
-> authorized repository impact
-> bounded patch
-> repository/workflow evidence
-> draft PR
```

Documentation must:

- define both provider and consumer jobs;
- state that the provider cannot access consumer code by default;
- distinguish engineering benchmark proof from customer demand;
- state actual implementation status;
- avoid advertising unimplemented commands as working;
- link the PRD, plan, architecture decision, and active Factory artifacts;
- preserve install and validation guidance;
- document data ownership, authorization, revocation, and evidence boundaries
  before hosted coordination is introduced.

The current repository and design-partner artifact are not represented as open
source. Before any consumer receives a pilot binary or source bundle, M2 must
provide explicit pilot/evaluation terms, a named security contact, a support
and incident route, signed artifact provenance, checksums, and install-integrity
verification. Public OSS or self-serve distribution is a separate release gate
that additionally requires an approved `LICENSE`, `SECURITY.md`, contribution
and support policies, release documentation, and maintained vulnerability
response. The plan cannot close a public-release boundary until those surfaces
exist.

---

## Test Matrix Wiring

| Task | Fast | Core | Acceptance | Risk | Cross-system |
|---|---|---|---|---|---|
| M0 | unit, schema | full gate | REB | CodeQL | none |
| M1 | fixture/unit | contract | CORPUS | license/provenance review | none |
| M2 | schema/unit | contract | CHG, AUTH, EVD, ACT | security/architecture review | configured signer fixture |
| M2.5 | evidence validation | product signal | DISC | privacy/commercial review | qualified provider and cohort |
| M3 | unit/integration | full gate | CHG | parser integrity | none |
| M4 | unit/scenario | full gate | IMP | held-out precision/recall | none |
| M5 | unit/contract | full gate | PLAN, AUTH | write-boundary review | none |
| M6 | unit/scenario | full gate | PATCH | diff/supply-chain review | none |
| M7 | unit/integration/scenario | full gate | VER, EVD, ACT | untrusted-command review | deterministic replay and offline canary |
| M8 | unit/integration | full gate | VER | credentials/network/security | approved sandbox |
| M9 | unit/integration | full gate | PR | GitHub permission/security | approved GitHub |
| M10 | local checks | evidence validation | DISC, AUTH, PILOT | privacy/product review | real provider and consumer repos |

All first-party code tasks use `make test-coverage`. All tasks run
`make prepush-full` before PR. CodeQL is required for dependency, generated
code, external-call, credential, data-exposure, workflow, and release-sensitive
changes.

When selected by task policy, independent lifecycle review runs after
`code-review` and before `commit-push` in this order:
`holdout-evaluator`, `trace-grader`, `evidence-attestor`. These are
independent or human-operated gates, not implementation-worker self-claims.
Shipping must verify their schema-valid, task-bound, current-run, marker-digest
and candidate-digest-bound passing artifacts under
`.factory/artifacts/lifecycle-evidence/<task>/` before creating a commit or PR.
The implementation worker may not write that root or
`.factory/artifacts/pr-lifecycle/`.

---

## Epic 0 — Rebaseline And Fail-Closed Compatibility

### M0: Correct the command and result foundations

**Priority:** P0
**Risk class:** Medium
**Blocked by:** none

**Acceptance item IDs:** `BASE-001`, `BASE-002`, `BASE-003`, `BASE-004`,
`BASE-005`, `REB-001`, `REB-002`, `REB-003`

**Tasks:**

- Make unimplemented commands return a typed nonzero result or remove them from
  the recognized command registry.
- Stop setting eval-oriented metadata on `init` and `check`.
- Introduce explicit API-provider versus model-provider terminology.
- Define the versioned migration path for command-result and evidence schemas.
- Preserve exit-code compatibility, including reserved exit code `6`.
- Revalidate the current Factory profile, task packets, validation contract,
  acceptance ledger/mapping, and closure map with the canonical Factory
  schemas and repo-pack validator.
- Decompose `scripts/validate_repo_pack.py` along stable validation seams while
  preserving or reducing the line ceiling in
  `.factory/artifacts/exceptions/architecture-debt-lumyn-migration-rebaseline.json`.
- Close `BASE-003` only with current-head CI, review, shipping, and post-merge
  evidence; the blocked historical PR lifecycle report is not carry-forward
  proof.
- Make the historical plan explicitly non-dispatchable and preserve historical
  evidence and approved exceptions without mutation.

**Repo paths:**

- `cmd/lumyn/`
- `internal/result/`
- `internal/exitcode/`
- `schemas/`
- `scripts/validate_repo_pack.py`
- `scripts/repo_pack_validation/`
- `docs/`
- `CHANGELOG.md`

**Run commands:**

- `make lint-fast`
- `make test-fast`
- `make test-coverage`
- `make test-contracts`
- `python3 "$FACTORY_REPO/scripts/factory_run_mission.py" validate-repo-pack
  --target-repo . --profile "$FACTORY_REPO/profiles/lumyn.yaml"
  --execution-plan
  .factory/artifacts/prd-to-plan/lumyn-migration-mvp/execution-plan.json
  --task-packets
  .factory/artifacts/prd-to-plan/lumyn-migration-mvp/task-packets.json
  --validation-contract
  .factory/artifacts/prd-to-plan/lumyn-migration-mvp/validation-contract.json
  --json`
- `factoryd doctor --config .factory/factoryd.example.json --repo lumyn --json`
- `make prepush-full`

**Test requirements:**

- Red-first test that every unimplemented command returns nonzero with a typed
  error.
- JSON-envelope tests for `init`, `check`, unknown, and unimplemented commands.
- Schema compatibility fixtures for old and new terminology.
- Exit-code stability tests.
- Canonical Factory schema, profile, and repo-pack validation.
- Active-config capability-grant validation and a negative selection test for
  the historical plan.

**Matrix wiring:** Tiers 1, 3, 4, 9, and CodeQL risk lane.

**Acceptance criteria:**

- No unimplemented command reports `status: pass`.
- API-provider and model-provider terms cannot be confused.
- Existing valid command-result `1.0` fixtures remain valid or have an explicit
  versioned migration.
- Every runner-ready task inherits the developer/architecture propagation
  contract and 12-level test matrix.
- The current plan is the only selectable plan and historical evidence is
  unchanged.

**Changelog impact:** required
**Changelog section:** Fixed
**Draft changelog entry:** Make unimplemented command paths fail closed and separate API-provider identity from model-provider metadata.
**Semver marker override:** `[semver:minor]`
**Contract/API impact:** Corrects CLI behavior and versions result/evidence terminology.
**Versioning/migration impact:** Requires documented compatibility for existing command-result and evidence artifacts.
**Architecture constraints:** Preserve stable exits, schema-backed JSON, local privacy, and fail-closed defaults.
**ADR required:** yes; extend ADR-0002 or add a focused compatibility decision.
**TDD first failing test(s):** Unimplemented `verify` returns nonzero; `init` does not claim eval mode.
**Cost/perf impact:** low
**Chaos/failure hypothesis:** Unknown or stubbed command paths can look green unless every dispatch branch is tested.
**Semantic invariants:** No command reports success without implemented behavior; API-provider identity never occupies a model-provider field.

---

## Epic 1 — Deterministic Migration Benchmark

### M1: Build the pinned gold and negative corpus

**Priority:** P0
**Risk class:** Medium
**Blocked by:** `M0`

**Acceptance item IDs:** `CORPUS-001`, `CORPUS-002`, `CORPUS-003`

**Tasks:**

- Create at least three pinned historical or synthetic provider changes, one
  for each supported deterministic change class.
- Create at least nine controlled TypeScript consumer fixtures, three per
  class. Each class contains at least 20 annotated affected sites and 20
  annotated non-affected candidates.
- Have an independent holdout owner, operating as `holdout-evaluator`,
  provision and freeze a separate evaluator-controlled suite with at least one
  repository per class. Before implementation scoring, it emits a
  lifecycle-owned M1 `holdout_result` binding the frozen suite commitment.
- Use M1 `holdout_policy.mode: provision` with only an opaque namespace and
  `hmac-sha256` algorithm plus `holdout_provisioning_required`. Do not invent a
  suite commitment in the static plan or claim current-candidate evaluation.
  M4, M6, and M7 use `holdout_evaluation_required` and evaluate-mode policies
  that resolve and byte-bind the trusted M1 result before scoring.
- Permit the implementation worker to commit only opaque case IDs,
  non-resolving provenance class and license posture, a frozen suite
  commitment, and encrypted or HMAC artifact commitments in
  `examples/holdout-manifest.json`. Prohibit source URLs, repository or package
  identifiers, and plaintext content digests.
- Store held-out inputs, answer keys, expected patches, expected labels, and
  raw traces only under private `LUMYN_HOLDOUT_ROOT`; never expose that root to
  `task-executor`.
- Store visible development expected patches and verification outcomes.
- Add license, attribution, source digest, and redistribution posture.
- Add negative cases for ambiguous, stale, unsupported, out-of-boundary,
  semantic-non-equivalence, and false-verification behavior.
- Include at least two positive and two negative patched-repository executions
  for each deterministic workflow-verification label the MVP emits.
- Split visible development fixtures from held-out evaluation fixtures.

**Repo paths:**

- `examples/provider-changes/`
- `examples/consumer-repos/`
- `examples/impact-reports/`
- `examples/patches/`
- `examples/negative/`
- `examples/holdout-manifest.json`
- `tests/`
- `docs/`

**Run commands:**

- `make lint-fast`
- `make test-fast`
- `make test-contracts`
- `make prepush-full`

**Test requirements:**

- Fixture manifest and digest validation.
- License/provenance completeness checks.
- Ground-truth call-site and expected-patch completeness checks.
- Negative fixture rejection tests.
- Negative repo-pack fixtures proving source URLs, repository or package
  identifiers, plaintext content digests, held-out inputs, expected labels,
  answer keys, patches, or traces are rejected.
- An independent provisioning/freeze test proving the M1 `holdout_result`
  matches the committed frozen suite commitment before later evaluation.
- Mutation tests proving a swapped M1 provisioning-result ref or changed result
  bytes invalidate M4/M6/M7 evaluation evidence.
- Access-isolation tests proving only `holdout-evaluator`, never the M1/M4/M6/M7
  implementation workspace or environment, receives `LUMYN_HOLDOUT_ROOT`.

**Matrix wiring:** Tiers 1, 2, 4, 9, and 11; risk review for provenance.

**Acceptance criteria:**

- Every visible development fixture is reproducible offline from committed
  permitted artifacts.
- Ground truth is fixed before impact/patch scoring.
- Every supported class meets the visible fixture, affected-site,
  non-affected-candidate, evaluator-controlled holdout, and
  workflow-positive/negative denominators in `CORPUS-001` and `CORPUS-002`.
- The M1 lifecycle-owned `holdout_result` proves independent provisioning,
  freeze, and commitment integrity. Later holdout results contain aggregate
  counts and commitments only; answer material and resolving provenance never
  enter the repository or implementation context.
- No fixture description implies provider endorsement or customer demand.

**Changelog impact:** required
**Changelog section:** Added
**Draft changelog entry:** Add a pinned, provenance-backed TypeScript API migration benchmark with held-out and negative fixtures.
**Semver marker override:** `[semver:minor]`
**Contract/API impact:** Introduces benchmark manifest and fixture conventions.
**Versioning/migration impact:** Fixture revisions require new IDs or versions; ground truth is immutable after scoring begins.
**Architecture constraints:** Deterministic, offline, licensed, and separated from private pilot data.
**ADR required:** no; ADR-0002 governs the benchmark role.
**TDD first failing test(s):** Missing provenance and altered fixture digest are rejected.
**Cost/perf impact:** low
**Chaos/failure hypothesis:** Upstream artifacts disappear or change, making an unpinned benchmark non-reproducible.
**Semantic invariants:** Public fixture evidence proves engineering only; held-out ground truth cannot be changed to improve scores.

---

## Epic 2 — Migration And Authorization Contracts

### M2: Define executable migration contracts

**Priority:** P0
**Risk class:** High
**Blocked by:** `M0`

**Acceptance item IDs:** `CHG-001`, `CHG-002`, `CHG-003`, `CHG-004`,
`AUTH-001`, `AUTH-002`, `AUTH-003`, `AUTH-004`, `AUTH-005`, `EVD-001`,
`EVD-002`, `EVD-003`, `ACT-001`, `ACT-002`

**Tasks:**

- Define schemas for provider enrollment, provider change, repository
  authorization, impact report, migration plan, patch manifest, migration
  verification, migration PR result, campaign kit, signed invitation,
  provider-status snapshot, authorization request, authorization revocation,
  consumer receipt-key binding, campaign connection receipt,
  provider-signed receipt acknowledgement, campaign attestation,
  command-isolation profile, and remediation outcome.
- Implement `lumyn campaign kit create`, configured-signing `lumyn change
  publish`, signed `lumyn campaign invite create`, out-of-band-confirmed `lumyn
  provider enroll`, consumer-side `lumyn campaign accept`, `lumyn trust
  refresh`, explicit consumer-signed `lumyn authorization issue` and
  `authorization revoke`, and opt-in `lumyn campaign receipt issue`, `submit`,
  provider-side `acknowledge`, and consumer `ack import`. Enrollment and
  acceptance run with the checkout unavailable, write only to the
  consumer-private state root, and produce no product authority. Authorization
  issuance performs no authorized side effect.
- Define packet lifecycle, immutability, canonical signing bytes, consumer
  trust-root pinning, provider/package ownership binding, first-pin bootstrap,
  key rotation, emergency re-enrollment/recovery, revocation, audience, expiry,
  replay, provenance, withdrawal, provider-status signer, maximum snapshot age,
  exact status endpoint, signed offline snapshot, provider
  receipt-acknowledgement signer, permitted receipt exchange classes, and
  anti-replay semantics. An invitation-supplied root, package binding,
  acknowledgement key, or endpoint is never self-authenticating.
- Define independent read/write/host-isolated-command/
  `provider_trust_status_read`/registry-network/sandbox-network/credential/
  remote-branch-write/PR-write/`campaign_receipt`/`provider_attestation`/
  `artifact_retention`/`artifact_deletion` scopes.
- Define the command-isolation contract: exact read-only and writable mounts,
  neutral home/temp roots, executable roots, environment classes, local socket
  and inherited-descriptor policy, process-tree inheritance, OS credential
  denial, resource budgets, and fail-closed backend detection.
- Keep those exact Lumyn product grants in private, schema-backed
  authorization artifacts. Factory's closed `approval`, `credentials`, and
  `network` grants authorize only the implementation worker and cite the
  validated product bundle; they never substitute for product authority.
- Define private runtime storage outside the checkout, TTL deletion, revocation
  receipts, the minimal signed connection-receipt allowlist,
  provider-authenticated consumer key binding, provider acknowledgement key,
  endpoint or offline export/import exchange, one-invitation-unit cardinality
  and idempotent deduplication protocol, and public aggregate/hash-only pilot
  evidence. Provider export and public commit are irreversible; revocation
  stops future disclosure and deletes only Lumyn-controlled private copies.
- Define sandbox payload classes, approved non-sensitive test-data policy, and
  provider logging/retention/deletion disclosure.
- Define separate impact, patch, verification, and delivery axes.
- Define freshness binding to packet, base/head, plan, and artifact hashes.
- Add pilot/evaluation distribution terms, a named security and support route,
  signed artifact provenance, checksums, and install-integrity verification.
- Add valid and invalid fixtures for every trust boundary.

**Repo paths:**

- `cmd/lumyn/`
- `internal/campaign/`
- `internal/invitation/`
- `internal/trust/`
- `internal/authorization/`
- `internal/isolation/`
- `internal/receipt/`
- `internal/attestation/`
- `schemas/`
- `examples/`
- `docs/product/`
- `docs/dev/`
- `docs/architecture/`
- `CHANGELOG.md`

**Run commands:**

- `make test-contracts`
- `make test-coverage`
- `make prepush-full`

**Test requirements:**

- Schema compilation and round-trip tests.
- Command tests for deterministic campaign-kit creation, configured signer
  failure, canonical packet bytes, invitation signature/expiry/audience,
  out-of-band first pin, active-root-signed rotation, emergency re-enrollment,
  no-authority consumer authorization-request receipt, explicit signed
  authorization issuance/revocation, and signed minimal connection receipts.
- Activation-isolation tests run `provider enroll` and `campaign accept` from a
  neutral directory with no checkout mount and only the exact input plus private
  state roots. Repository/Git before-and-after digests prove no mutation only;
  the unavailable mount and deny audit prove the commands cannot read the
  checkout.
- Provider-status tests cover signed offline snapshots, exact-endpoint online
  refresh, maximum age, response signature, audience, nonce/replay, rotation,
  withdrawal, endpoint mismatch, undeclared egress, and absence of repository
  or consumer data in requests.
- Command-isolation fixtures attempt to read host home, SSH/GPG/cloud credential
  stores and keychain, connect to agent/Docker/unrelated local sockets, inherit
  extra file descriptors, access undeclared mounts or OS credentials, and
  escape through child processes; every attempt fails or the command is not
  launched.
- Connection-receipt tests verify consumer issuer, invitation and packet
  binding, provider-authenticated consumer signer, opaque IDs, consent-policy
  and key-binding digests, audience, expiry, nonce, pinned online or offline
  exchange, provider acknowledgement signature, one-invitation-unit
  cardinality, idempotent same-digest retry, conflicting-unit and replay
  rejection, and the exact prohibited-field list.
- Negative fixtures for arbitrary scripts, invalid/untrusted/replayed packets,
  invitation-supplied bootstrap roots, unconfirmed fingerprints, unsigned
  rotation, missing provenance, overbroad authorization, unsafe runtime roots,
  stale/missing status, sandbox payload leakage, auto-merge, raw provider data
  sharing, recallable-disclosure claims, and stale evidence.
- Cross-contract tests proving a product capability cannot be placed directly
  in `.factory/factoryd.json` and a generic Factory grant cannot satisfy a
  missing Lumyn authorization artifact.
- Compatibility tests for retained workflow/evidence contracts.

**Matrix wiring:** Tiers 1, 2, 4, 9, and security/architecture risk lanes.

**Acceptance criteria:**

- All planned migration artifacts are executable contracts.
- The provider and consumer can complete the standardized activation handoff
  without bespoke artifact authoring, shared secrets, or implied repository
  authority.
- First enrollment is authenticated independently of the invitation; normal
  rotation preserves key continuity, emergency recovery requires explicit
  re-enrollment, and enrollment/acceptance run without checkout access.
- Consumer authorization has an explicit signed issue/revoke path; a request
  alone never becomes authority, and the 60-minute activation path includes
  issuance.
- Packet status is provably current through a signed snapshot or exact
  authorized status read with no undeclared egress.
- A provider can verify and acknowledge the minimal sponsored-program meter
  using an authenticated consumer signer and provider-signed deduplicated
  acknowledgement without receiving repository identity or raw evidence.
- Pilot distribution has explicit terms, security/support contacts, and signed
  install-integrity evidence without claiming the repository is OSS.
- Provider campaign authority cannot imply consumer repository authority.
- Published packet mutation, signature/trust/freshness failure, overbroad
  capability scope, unsafe private storage, default raw-data sharing, and
  auto-merge are schema-invalid.

**Changelog impact:** required
**Changelog section:** Added
**Draft changelog entry:** Add provider-change, provider-status, explicit authorization issue/revoke, host-isolation, connection-receipt, migration, verification, PR-result, attestation, and outcome contracts.
**Semver marker override:** `[semver:minor]`
**Contract/API impact:** Adds public versioned artifact schemas and trust states.
**Versioning/migration impact:** Existing eval-centric evidence fields receive an explicit compatibility path.
**Architecture constraints:** Declarative input, dual-principal authority, immutable provenance, authorized freshness transport, host isolation, and separate private, receipt, provider-visible, and public artifacts.
**ADR required:** yes; ADR-0002 plus a schema-version decision if compatibility requires it.
**TDD first failing test(s):** Packet with executable script, authorization with implicit write, and verification missing environment must fail.
**Cost/perf impact:** low
**Chaos/failure hypothesis:** A permissive schema could turn the provider packet into a code-execution or data-exfiltration channel.
**Semantic invariants:** Provider input is untrusted; consumer consent is explicit; evidence axes never collapse into a false roll-up.

### M2.5: Qualify the design partner and freeze the pilot protocol

**Priority:** P0 product gate
**Risk class:** High
**Blocked by:** `M0`, `M2`

**Acceptance item IDs:** `DISC-001`, `DISC-002`

**Tasks:**

- Obtain written commitment from one qualified provider for a consequential
  migration within six months, including the named economic buyer, operator,
  authoritative source/target artifacts, prerelease-sharing authority,
  non-production sandbox, operational signed provider-status snapshot or exact
  endpoint channel with pinned key and maximum-age policy, provider
  receipt-acknowledgement key and pinned endpoint or offline exchange,
  paid-pilot price, and decision process.
- Require recurring-value evidence: at least two consequential migrations or
  deprecations expected within 12 months, or at least 20 named managed
  integrations eligible for an annual connected-repository program. Name the
  decision owner, date, and criteria for the post-pilot annual-platform or
  second-campaign purchase.
- Use the M2 standard campaign kit, signed packet, signed invitation,
  authorization issue/revoke flow, minimal signed connection receipt, and
  provider-authenticated consumer receipt-key binding plus provider-signed
  acknowledgement/import protocol; do not qualify a bespoke artifact bundle
  that the product cannot reproduce.
- Independently inspect a provider-signed canary against authoritative source
  and target artifacts and prove that it contains at least one supported
  change class. This is qualification evidence only: it does not close
  `CHG-003` or `CHG-004`, authorize packet execution, or replace the M3 runtime
  verifier.
- Prequalify a frozen cohort of at least five distinct repository IDs across
  at least three Consumer Organizations, with one accountable maintainer per
  repository. Each repository must match GitHub + TypeScript/Node + direct
  official SDK + one selected package root + `package-lock.json` and be
  plausibly affected by the source version, operation, or field. Each candidate
  environment must also prove a supported OS/architecture and an enforceable
  fail-closed host-isolation backend before it can enter the frozen cohort.
- Before observing results, preregister the invitation and measurement windows,
  absolute campaign judgment deadline, baseline comparator, loaded labor
  rates, the boundary between one-time product development and Lumyn campaign
  COGS, correction rubric, one primary provider outcome and its material
  threshold, paid-pilot price, contribution threshold, recurring gross-margin
  and operator-hour thresholds, buyer-TCO treatment, continuation evidence,
  confidentiality, retention, the exact minimal connection-receipt fields,
  provider-authenticated consumer key-binding method, provider acknowledgement
  key and exchange, one-invitation-unit cardinality and deduplication rule,
  richer provider-visible fields, public evidence fields, the irreversibility
  of provider/public disclosure, supported consumer host-isolation matrix, and
  consented data/payload protocol.
- Before collecting, storing, or disclosing identifiable external evidence,
  complete a narrow manual privacy/legal preflight under the Factory
  implementation worker's task-scoped `approval`. The preflight names allowed
  private fields, the approved private storage boundary, participant consent,
  TTL, expiry/revocation deletion, deletion-receipt and orphan ownership, and
  minimal connection-receipt, authenticated signer-binding, provider
  acknowledgement/cardinality policy, and separately consented
  aggregate/hash-only public fields. Consent must state that provider export
  and public commit cannot be recalled. This preflight authorizes evidence
  handling for M2.5 only; it is not Lumyn runtime product authority and cannot
  authorize a repository or sandbox action. The active approval grant must bind
  the canonical digest of this exact preflight scope; a generic approval or a
  stale digest cannot release collection.
- Verify the M2 pilot distribution package has explicit evaluation/commercial
  terms, a named security and support route, signed provenance, checksums, and
  install-integrity instructions. Do not describe it as OSS.
- Store identifiable provider and consumer evidence only in the approved
  private evidence system. Commit only redacted aggregate evidence and hashes
  under `.factory/artifacts/pilot/lumyn-migration-mvp/public/`.
- Stop and trigger product reframe review after two qualified provider attempts
  fail the gate.

**Repo paths:**

- `.factory/artifacts/pilot/lumyn-migration-mvp/public/`
- `scripts/`
- `schemas/`
- `tests/`
- `docs/product/`
- approved private external evidence records

**Run commands:**

- `make test-contracts`
- `python3 scripts/validate_design_partner_evidence.py --self-test`
- `python3 scripts/validate_design_partner_evidence.py --attestation "$LUMYN_PRIVATE_DESIGN_PARTNER_ATTESTATION" --public-manifest .factory/artifacts/pilot/lumyn-migration-mvp/public/design-partner-manifest.json`
- `make prepush-full`

**Test requirements:**

- The `evidence-attestor` independently verifies the actual private
  attestation and its public aggregate/hash-only manifest; self-tests alone
  cannot close `DISC-001` or `DISC-002`.
- The attestation is a lifecycle-owned artifact produced independently after
  implementation validation and before `commit-push`; the implementation
  worker cannot attest its own product-signal evidence.
- Manual-preflight negative cases cover missing participant consent, an
  unapproved private field or storage boundary, absent TTL/deletion ownership,
  an unapproved connection-receipt or signer-binding field, an unsigned
  acknowledgement or ambiguous unit cardinality, a promise to recall
  external/public copies, and public disclosure without separate consent.
- Approval-binding cases prove only an active M2.5 approval carrying the exact
  canonical preflight-scope digest can release collection. Missing, generic,
  stale, or wrong-scope approvals fail; changing allowed fields, private
  storage, TTL/deletion rules, or public fields changes the digest.
- Privacy review proving that committed output is aggregate or hash-only.
- Distribution-package tests verify explicit terms, security/support contacts,
  signer identity, checksums, install integrity, and absence of an OSS claim.
- Qualification tests prove the provider status channel produces a valid
  signed, in-age snapshot or exact authorized endpoint response and that every
  candidate OS/architecture can enforce the declared host-isolation profile.
- Protocol completeness and timestamp check proving the cohort and thresholds
  were frozen before outcomes.
- Negative fixtures for a one-organization cohort, duplicate repository IDs,
  unsupported migrations, nominal price, post-outcome thresholds, and missing
  judgment deadline.

**Matrix wiring:** Tiers 4, 10, and 11; product, privacy, and
evidence-attestation gates.

**Acceptance criteria:**

- `DISC-001` and `DISC-002` have direct external evidence.
- `DISC-001` includes an annual connected-repository or second-campaign
  decision backed by the required frequency or managed-integration evidence.
- The denominator is five distinct eligible repositories across at least three
  organizations, with one funnel unit per repository.
- The canary proves at least one supported change class and every cohort member
  has recorded plausible exposure to that migration.
- The primary provider outcome and material threshold are frozen before any
  invitation.
- The minimal connection-receipt schema, authenticated consumer signer binding,
  provider acknowledgement key/exchange, and one-unit deduplication rule are
  frozen, and no repository is counted or billed without a valid
  provider-signed acknowledgement.
- The provider status channel is operational and every candidate environment
  proves a supported OS/architecture with enforceable host isolation.
- The supported pilot distribution package is security- and support-routable,
  integrity-verifiable, and honestly licensed under explicit pilot terms.
- The committed repository contains no identifiable pilot evidence.
- M5 remains blocked until both items close.

**Changelog impact:** none for evidence-only work
**Contract/API impact:** none; any product behavior requested by the partner becomes governed follow-up work.
**Versioning/migration impact:** protocol revisions create a new preregistration version before invitations.
**Architecture constraints:** private evidence remains outside the checkout; connection receipts are minimal and signed; public evidence is aggregate/hash-only and irreversibly disclosed.
**ADR required:** no unless qualification introduces a new trust boundary.
**TDD first failing test(s):** qualification record missing price, cohort support shape, consent fields, or pre-outcome timestamp fails validation.
**Cost/perf impact:** medium
**Chaos/failure hypothesis:** A friendly but unqualified provider or retrospectively selected cohort creates false product-market evidence.
**Semantic invariants:** no threshold changes after outcomes; private data never becomes repo evidence; contracts and impact can proceed, but migration planning cannot bypass this gate.

---

## Epic 3 — Change Understanding And Repository Impact

### M3: Ingest pinned change sources and classify semantic changes

**Priority:** P0
**Risk class:** High
**Blocked by:** `M1`, `M2`

**Acceptance item IDs:** `BASE-004`, `CHG-001`, `CHG-002`, `CHG-003`,
`CHG-004`

**Tasks:**

- Extend source intake from one API surface to pinned source and target
  OpenAPI/SDK artifacts.
- Verify canonical packet signing bytes against the consumer-pinned provider
  trust root and provider-to-package ownership binding.
- Enforce issuer key, issue time, audience, expiry, rotation, revocation,
  withdrawal, replay, publication immutability, and digest checks before
  classification.
- Reject executable hooks, arbitrary scripts, and undeclared external
  references; packet content remains declarative.
- Validate package/version/digest consistency.
- Normalize only the three supported semantic change classes.
- Preserve provider-declared unsupported and needs-input conditions.
- Emit typed change entries with concrete source references.

**Repo paths:**

- `internal/source/`
- `internal/change/`
- `internal/trust/`
- `schemas/`
- `examples/provider-changes/`
- `tests/`
- `docs/`

**Run commands:**

- `make lint-fast`
- `make test-fast`
- `make test-coverage`
- `make test-contracts`
- `make prepush-full`

**Test requirements:**

- JSON/YAML normalization equivalence.
- Old/new artifact digest mismatch and stale packet tests.
- Invalid signature, unpinned trust root, wrong provider/package binding,
  wrong audience, expired/revoked/withdrawn/replayed packet, changed published
  bytes, and arbitrary-script denial tests.
- Supported versus unsupported semantic-change classification tests.
- External reference handling remains offline unless explicitly approved.

**Matrix wiring:** Tiers 1, 2, 4, 9, and 11.

**Acceptance criteria:**

- Every normalized change cites source and target evidence.
- No classification occurs until the complete `CHG-003` and `CHG-004` runtime
  trust gate passes.
- Unsupported change classes remain explicit.
- Parser ambiguity blocks downstream patching.

**Changelog impact:** required
**Changelog section:** Added
**Draft changelog entry:** Add pinned source/target API and SDK change intake for supported deterministic migration classes.
**Semver marker override:** `[semver:minor]`
**Contract/API impact:** Extends source configuration and adds typed change output.
**Versioning/migration impact:** Existing single-surface source config remains readable during migration.
**Architecture constraints:** Keep parsing separate from repository analysis and patching; no network in deterministic tests.
**ADR required:** no if ADR-0002 boundaries hold.
**TDD first failing test(s):** Mismatched target digest and unsupported auth change block normalization.
**Cost/perf impact:** low
**Chaos/failure hypothesis:** OpenAPI and SDK releases disagree, producing a plausible but wrong mapping.
**Semantic invariants:** Change classification never invents provider intent; every entry is provenance-bound.

### M4: Analyze TypeScript consumer impact

**Priority:** P0
**Risk class:** High
**Blocked by:** `M3`

**Acceptance item IDs:** `AUTH-001`, `AUTH-003`, `AUTH-004`, `IMP-001`,
`IMP-002`, `IMP-003`, `IMP-004`, `IMP-005`

**Tasks:**

- Detect the selected official npm SDK and version.
- Parse `tsconfig.json`, `package.json`, and supported lockfile state.
- Enforce the authorized read root across path normalization, symlinks,
  `tsconfig extends`, and TypeScript project references; any escape blocks.
- Find direct imports and statically identifiable affected call sites.
- Detect one-hop wrapper uncertainty without speculative patching.
- Exclude generated, vendored, build-output, dynamic, multi-version, and
  non-selected package-root cases.
- Score visible and held-out fixtures with separate precision and recall.
- Require an independent `holdout-evaluator` to score the frozen held-out
  repository set after implementation validation and before `commit-push`.
  Only that evaluator receives `LUMYN_HOLDOUT_ROOT`; the implementation
  workspace, prompt, environment, and mount set do not.

**Repo paths:**

- `internal/impact/`
- `internal/typescript/`
- `schemas/`
- `examples/consumer-repos/`
- `tests/`
- `docs/`

**Run commands:**

- `make test-fast`
- `make test-coverage`
- `make test-contracts`
- `make prepush-full`

**Test requirements:**

- AST/parser-backed call-site tests; text matching alone is insufficient.
- Direct import, aliased import, one-hop wrapper, dynamic access, generated
  path, monorepo, multi-version, and lockfile fixtures.
- Read-root escape tests for symlinks, `tsconfig extends`, project references,
  and package-root traversal.
- Held-out precision/recall report.
- A lifecycle-owned `holdout_result` bound to M4 and the current validation
  work proof; the implementation worker cannot inspect or self-grade the
  held-out answer key.
- The result exposes aggregate counts, frozen suite/candidate digests, and
  failing opaque case IDs only; it never commits held-out source, answers,
  expected labels, patches, or raw traces.
- Deterministic repeated-run output.

**Matrix wiring:** Tiers 1, 2, 4, 7, 9, and 11; CodeQL when parser dependencies
are added.

**Acceptance criteria:**

- Every supported class has `100%` held-out recall, at least `95%` held-out
  precision, and zero false `unaffected` results.
- `uncertain`, `unsupported`, or missed on a known-supported affected site is a
  false negative. Report per-class plus micro/macro aggregate confusion
  matrices.
- Unsupported and uncertain cases remain visible.

**Changelog impact:** required
**Changelog section:** Added
**Draft changelog entry:** Add read-only TypeScript/npm repository impact analysis with explicit coverage and uncertainty.
**Semver marker override:** `[semver:minor]`
**Contract/API impact:** Adds `lumyn impact` behavior and impact-report output.
**Versioning/migration impact:** New output starts at schema version `1.0`.
**Architecture constraints:** Parser-backed, read-only, deterministic, bounded to one selected package root.
**ADR required:** yes if a TypeScript parser runtime crosses the Go boundary.
**TDD first failing test(s):** Aliased import is found; dynamic wrapper is uncertain; generated code is excluded.
**Cost/perf impact:** medium
**Chaos/failure hypothesis:** A high apparent precision hides missed wrapper call sites and yields false `unaffected`.
**Semantic invariants:** Impact never mutates; coverage always names scope and exclusions; uncertainty cannot become unaffected.

---

## Epic 4 — Reviewable Plan And Bounded Patch

### M5: Produce a no-write migration plan and approval gate

**Priority:** P0
**Risk class:** High
**Blocked by:** `M2`, `M2.5`, `M4`

**Acceptance item IDs:** `AUTH-001`, `AUTH-002`, `AUTH-004`, `PLAN-001`,
`PLAN-002`

**Tasks:**

- Convert supported impact entries into a deterministic migration plan.
- List every proposed file, recipe, command, permission, network destination,
  credential class, verification stage, and residual risk.
- Render the exact follow-on grant manifest needed for patch, command,
  registry, GitHub, receipt, or richer-reporting actions. The maintainer uses
  the M2 `authorization issue` command to sign it; the plan cannot self-approve
  or widen the active bundle.
- Implement `lumyn authorization validate` as a fail-closed local gate over the
  private bundle's task ID, exact capabilities, scopes, packet/plan/base
  bindings, expiry, revocation, retention, and deletion authorities. It emits
  no secret or private evidence.
- Prove plan mode does not modify repository or Git state.

**Repo paths:**

- `internal/migrationplan/`
- `internal/authorization/`
- `schemas/`
- `tests/`
- `docs/`

**Run commands:**

- `make test-fast`
- `make test-coverage`
- `make test-contracts`
- `make prepush-full`

**Test requirements:**

- Filesystem and Git before/after immutability checks.
- Approval expiry and revocation tests.
- Request-versus-issued-bundle tests prove that an unissued template, changed
  grant manifest, missing consumer signature, or superseded authorization has
  no authority.
- Missing bundle, wrong task, partial capability set, stale binding, expired or
  revoked grant, and absent retention/deletion authority tests.
- Missing business-value and unsupported-change tests.
- Stable JSON and typed exit tests.

**Matrix wiring:** Tiers 1, 2, 3, 4, 9, and security/architecture risk lanes.

**Acceptance criteria:**

- The maintainer can review the entire mutation and execution boundary before
  authorizing it.
- Plan mode is demonstrably read-only.
- Approval cannot silently widen after issuance.
- The standard accept-to-impact path includes explicit authorization issuance
  and needs no manually fabricated private bundle.
- Every live task packet runs the validator as both a task validation command
  and a final validation command; a generic Factory approval cannot replace
  the private bundle.
- The M5 task packet carries runner-enforced `gated_by_acceptance_items` entries
  for `DISC-001` and `DISC-002`; `blocked_by: M2.5` alone is not sufficient to
  release the task from the queue.

**Changelog impact:** required
**Changelog section:** Added
**Draft changelog entry:** Add reviewable no-write migration plans and explicit consumer authorization gates.
**Semver marker override:** `[semver:minor]`
**Contract/API impact:** Adds migration-plan and authorization behavior.
**Versioning/migration impact:** Authorization and plan revisions invalidate prior approval.
**Architecture constraints:** Separate planning from mutation; immutable plan digest; independent capability scopes.
**ADR required:** no if M2 contract is unchanged.
**TDD first failing test(s):** Plan changes a file or approval survives a plan digest change.
**Cost/perf impact:** low
**Chaos/failure hypothesis:** A time-of-check/time-of-use plan drift widens the eventual patch.
**Semantic invariants:** No write before approved plan; approval binds exact plan and capability scopes.

### M6: Apply the three deterministic migration recipes

**Priority:** P0
**Risk class:** High
**Blocked by:** `M5`

**Acceptance item IDs:** `CHG-002`, `CHG-004`, `AUTH-002`, `AUTH-004`,
`PATCH-001`, `PATCH-002`, `PATCH-003`, `PATCH-004`, `PATCH-005`, `PATCH-006`

**Tasks:**

- Apply method/operation, request-property, and response-property mappings.
- Check every recipe's semantic-equivalence preconditions and emit no edit when
  any precondition fails.
- Update dependency and supported lockfile state only with the exact pinned
  Node/npm versions, registry or immutable snapshot, integrity inputs,
  toolchain digest, disabled lifecycle scripts, and approved registry-network
  grant.
- Run in an isolated worktree or equivalent workspace.
- Immediately before every filesystem or lockfile write, revalidate the exact
  current packet bytes, digest, trust root, provider/package binding, lifecycle
  state, audience, expiry, rotation, revocation, withdrawal, supersession, and
  replay status from either a still-current signed offline status snapshot or
  an exact-endpoint read under `provider_trust_status_read`. A decision made
  during planning cannot be cached for writes, and no undeclared status egress
  is allowed.
- Run any package-manager or repository command through the M2 host-isolation
  profile; an unavailable isolation backend blocks rather than falling back to
  the host shell.
- Enforce path, file-count, line-count, and diff-content budgets.
- Produce a patch manifest that maps each edit to change and recipe IDs.
- Refuse ambiguous or needs-input changes.

**Repo paths:**

- `internal/patch/`
- `internal/workspace/`
- `internal/trust/`
- `internal/isolation/`
- `schemas/`
- `examples/patches/`
- `tests/`
- `docs/`

**Run commands:**

- `make test-fast`
- `make test-coverage`
- `make test-contracts`
- `make prepush-full`

**Test requirements:**

- Golden patch tests and byte-stable repeat runs.
- Out-of-boundary, symlink, generated-file, lockfile, conflict, and stale-base
  tests.
- Time-of-check/time-of-use tests revoke, withdraw, supersede, rotate, replay,
  expire, or change the packet digest after plan approval but before each
  write; every case produces no edit.
- No-op/idempotent rerun tests.
- Expected-patch comparison over held-out fixtures.
- An independent `holdout-evaluator` verifies golden output against the
  held-out patch set before `commit-push`. Only that evaluator receives
  `LUMYN_HOLDOUT_ROOT`; the implementation worker receives neither the root nor
  held-out answer material.

**Matrix wiring:** Tiers 1, 2, 3, 4, 5, 7, 9, and 11; CodeQL risk lane.

**Acceptance criteria:**

- Every edit has provenance.
- Identical pinned inputs produce the same patch.
- No unrelated or unauthorized edit is produced.
- Ambiguity produces no patch.

**Changelog impact:** required
**Changelog section:** Added
**Draft changelog entry:** Add isolated, provenance-bound deterministic TypeScript migration patches for three supported change classes.
**Semver marker override:** `[semver:minor]`
**Contract/API impact:** Adds migration application and patch-manifest behavior.
**Versioning/migration impact:** Recipe versions are immutable and recorded in patch evidence.
**Architecture constraints:** Disposable workspace, bounded writes, no arbitrary scripts, deterministic transforms.
**ADR required:** yes for workspace isolation and patch-boundary enforcement.
**TDD first failing test(s):** Path escape, stale base, and missing business value produce no patch.
**Cost/perf impact:** medium
**Chaos/failure hypothesis:** Symlinks, lockfile tooling, or stale source can escape the intended write boundary.
**Semantic invariants:** Every write is planned, authorized, provenance-bound, and reversible.

---

## Epic 5 — Repository And Workflow Verification

### M7: Verify repository and deterministic workflow behavior

**Priority:** P0
**Risk class:** High
**Blocked by:** `M6`, `M2`

**Acceptance item IDs:** `AUTH-002`, `AUTH-004`, `VER-001`, `VER-002`,
`VER-003`, `VER-005`, `VER-006`, `VER-007`, `EVD-001`, `EVD-002`, `EVD-003`,
`ACT-003`

**Tasks:**

- Capture pre-patch dependency, compile/typecheck, and selected-test baseline.
- Run dependency integrity, compile/typecheck, and allowlisted tests after patch.
- Execute every repository-defined command only through the supported
  fail-closed isolation backend with exact mounts, neutral home/temp roots,
  sanitized environment, no host credential stores or OS credentials, no
  agent/Docker/unrelated local sockets, no extra inherited descriptors, and
  process-tree inheritance of the same restrictions.
- Implement independent contract/cassette replay over retained
  workflow/evidence contracts and label it
  `workflow_contract_replay_passed`.
- Execute approved entrypoints from the exact patched repository head for
  deterministic replay and mock verification; record observed interaction and
  outcome evidence.
- Use the canonical static, repo, contract-replay,
  `workflow_verified_replay`, and `workflow_verified_mock` labels.
- Bind evidence to packet, commits, plan, patch, and environment.
- Add proof-of-behavior scorecards and freshness invalidation.
- Implement the private-artifact retention owner in `internal/retention/`.
  Creation, read, startup, and next-run recovery enforce TTL and
  expiry/revocation deletion automatically; deletion is retry-safe and emits a
  receipt or durable orphan report without rewriting historical closure.
- Implement `lumyn artifacts gc` as the explicit operator recovery path for
  retrying failed deletion and inspecting unresolved orphan records. It does
  not extend TTL or revive revoked authority.
- Implement `lumyn trace` as a local renderer over the bound evidence axes. It
  performs no network call and never implies provider disclosure.
- Implement `lumyn canary run --offline` over the standard synthetic campaign
  kit. It must verify invitation acceptance, explicitly issue a synthetic
  consumer-signed authorization, run impact, plan, bounded patch,
  host-isolated deterministic replay or mock proof, render evidence, and emit a
  local draft-PR preview plus receipts without live credentials, sandbox,
  network, provider reporting, or remote writes.

**Repo paths:**

- `cmd/lumyn/`
- `internal/verify/`
- `internal/replay/`
- `internal/isolation/`
- `internal/evidence/`
- `internal/redaction/`
- `internal/retention/`
- `schemas/`
- `workflows/`
- `cassettes/`
- `runs/`
- `tests/`
- `docs/`

**Run commands:**

- `make test-fast`
- `make test-coverage`
- `go run ./cmd/lumyn --json canary run --kit examples/canary/campaign-kit.json --repo examples/canary/consumer-repo --offline --pr-preview .factory/artifacts/task-runs/M7/canary-pr-preview.md`
- `go run ./cmd/lumyn --json artifacts gc --dry-run`
- `make test-contracts`
- `make prepush-full`

**Test requirements:**

- Pre-existing failure attribution.
- Network-disabled command execution.
- Package-lifecycle-script denial.
- Adversarial host-isolation tests attempt undeclared mount and host-home reads,
  SSH/GPG/cloud/keychain credential access, agent/Docker/local-service socket
  access, inherited-descriptor use, OS credential access, and child-process
  escape. Every attempt is denied, and an unavailable enforcement backend
  blocks before command launch.
- Replay determinism and stale evidence.
- Negative causal-binding tests proving independent replay or execution from a
  different head cannot produce a workflow-verified label.
- Boundary, cleanup, redaction, proof-gap, and false-verification scenarios.
- Offline canary golden-path and stage-failure tests prove every stage receipt
  is causally bound, no prohibited capability is requested, and any
  unimplemented stage exits nonzero instead of falling through to `demo`.
- Retention tests cover creation, read, process restart, next-run sweep, TTL
  expiry, revocation, partial deletion, retry, crash recovery, deletion
  receipts, orphan reporting, and preservation of historical closure claims.
- Trace rendering tests prove every evidence axis is visible, stale bindings
  remain visibly stale, and no network or provider-attestation path is called.
- Independent lifecycle-owned `holdout_result` and `trace_grade_report`
  artifacts bind to M7 and the current validation work proof and pass before
  `commit-push`.

**Matrix wiring:** Tiers 1, 2, 3, 4, 5, 7, 8, 9, and 11; CodeQL risk lane.

**Acceptance criteria:**

- Repository and workflow evidence cannot be conflated.
- Pre-existing failures are reported separately.
- Every held-out workflow case receives its exact expected label with zero
  false positives or false negatives.
- The synthetic two-sided activation canary reaches a local evidence-backed
  draft-PR preview through explicit authorization issuance and host-isolated
  commands without hidden network, credentials, provider export, or remote
  writes.
- Evidence becomes stale when any bound input changes.
- `lumyn trace` renders the exact local evidence bundle and fails nonzero on
  missing, invalid, or stale inputs.

**Changelog impact:** required
**Changelog section:** Added
**Draft changelog entry:** Add baseline-aware repository checks, deterministic workflow replay evidence, the offline two-sided activation canary, and local `lumyn trace` rendering for migration patches.
**Semver marker override:** `[semver:minor]`
**Contract/API impact:** Adds verification ladder, migration-verification artifacts, and the receipt-backed `lumyn canary run --offline` surface.
**Versioning/migration impact:** Retained workflow/evidence schemas remain compatible or receive explicit versions.
**Architecture constraints:** Untrusted commands isolated; network disabled by default; normalized evidence precedes validators.
**ADR required:** yes for command sandbox and evidence binding.
**TDD first failing test(s):** Existing failing test is misattributed; replay-only evidence reports sandbox verification.
**Cost/perf impact:** medium
**Chaos/failure hypothesis:** Flaky tests or stale evidence create a false regression or false green result.
**Semantic invariants:** Evidence boundary is explicit; stale or incomplete evidence never verifies.

### M8: Add approved sandbox read-back verification

**Priority:** P1
**Risk class:** High
**Blocked by:** `M7`; approved, independent `command_execution`,
`sandbox_network`, `sandbox_credential`, `sandbox_request_disclosure`,
`artifact_retention`, and `artifact_deletion` product grants; and the Factory
worker's task-scoped `approval`, `credentials`, and `network` grants

**Acceptance item IDs:** `CHG-002`, `CHG-004`, `AUTH-002`, `AUTH-003`,
`AUTH-004`, `AUTH-005`, `VER-003`, `VER-004`, `VER-005`, `VER-006`,
`VER-007`

**Tasks:**

- Add allowlisted provider sandbox execution.
- Require `command_execution` for the exact patched-head entrypoint; sandbox
  network, credential, or payload authority never implies command authority.
  The entrypoint runs under the M2 host-isolation profile.
- Require an explicit sandbox-payload disclosure grant that allows only
  synthetic or approved non-sensitive test data and records provider logging,
  retention, and deletion terms.
- Isolate non-production credentials from build/test commands.
- Inside Lumyn, immediately before every sandbox side effect and every retry,
  revalidate the current packet trust state from a current signed offline
  snapshot or an exact `provider_trust_status_read` grant, plus every exact
  action grant. Do not cache an earlier validation decision. Factory dispatch
  and its worker grants do not perform or confer this product authorization.
- Add namespaces, idempotency, request/write budgets, settle/retry, cleanup, and
  orphan evidence.
- Preserve sandbox-versus-production limitations in the result.

**Repo paths:**

- `internal/live/`
- `internal/verify/`
- `internal/authorization/`
- `internal/isolation/`
- `internal/redaction/`
- `internal/trust/`
- `schemas/`
- `tests/`
- `docs/`

**Run commands:**

- `make test-fast`
- `make test-coverage`
- `make test-contracts`
- `go run ./cmd/lumyn --json authorization validate --bundle "$LUMYN_PRIVATE_PRODUCT_AUTHORITY_BUNDLE" --task M8`
- `make prepush-full`
- approved design-partner sandbox command recorded in task evidence

**Test requirements:**

- Deterministic mock tests for timeout, retry, cleanup, orphan, auth, budget,
  network allowlist, and redaction.
- Live test only after task-scoped human approval.
- Credential non-leakage between command stages.
- Host-isolation inheritance and provider-status snapshot/authorized-refresh
  tests at the sandbox side-effect boundary.
- Production-data/PII/secret payload denial and provider-retention disclosure
  tests.
- Exact patched-head causal execution and wrong-head rejection.
- Time-of-check/time-of-use tests revoke or expire a product grant, or revoke,
  withdraw, supersede, rotate, replay, or mutate the packet, after validation
  but before each sandbox call and retry; every live action is blocked.

**Matrix wiring:** Tiers 1, 2, 4, 5, 9, 11, and approved Tier 12; CodeQL and
security review.

**Acceptance criteria:**

- Lumyn sandbox execution cannot begin unless all six exact product grants are
  present, current, mutually consistent, and validated at the action boundary.
  A Factory-run M8 live integration test additionally requires the
  implementation worker's three task-scoped Factory grants.
- The `authorization validate` command is closure and diagnostic proof. The
  authoritative live-action check runs inside Lumyn at the side-effect
  boundary; it is not a Factory pre-dispatch product gate.
- Cleanup success or orphan evidence is mandatory.
- Sandbox proof is never represented as production guarantee.

**Changelog impact:** required
**Changelog section:** Added
**Draft changelog entry:** Add explicitly approved provider-sandbox read-back verification with budgets, cleanup, and orphan evidence.
**Semver marker override:** `[semver:minor]`
**Contract/API impact:** Adds live verification configuration and evidence.
**Versioning/migration impact:** Sandbox evidence records environment identity and freshness.
**Architecture constraints:** No production credentials; task-scoped allowlist; isolated secrets; fail-closed cleanup.
**ADR required:** yes for credential and network posture.
**TDD first failing test(s):** Missing grant, leaked sandbox credential, and cleanup failure all block verification.
**Cost/perf impact:** medium
**Chaos/failure hypothesis:** Sandbox drift or partial cleanup appears as a valid business outcome.
**Semantic invariants:** Live access is explicit, non-production, budgeted, isolated, and never silently downgraded.

---

## Epic 6 — Customer-Controlled Draft PR Delivery

### M9: Produce migration evidence and open an idempotent draft PR

**Priority:** P0
**Risk class:** High
**Blocked by:** `M7`; approved `github_branch_write`, `github_pr_write`,
`artifact_retention`, and `artifact_deletion` product grants; and the Factory
worker's task-scoped `approval`, `credentials`, and `network` grants.
`provider_attestation` is optional and required only if the consumer separately
authorizes a provider-reporting action. M8 sandbox proof is optional for PR
delivery and included only when available.

**Acceptance item IDs:** `CHG-002`, `CHG-004`, `AUTH-002`, `AUTH-003`,
`AUTH-004`, `PR-001`, `PR-002`, `PR-003`, `PR-004`, `PR-005`

**Tasks:**

- Render the migration evidence packet and reviewer checklist.
- Require independent remote-branch-write and PR-write authorization after plan
  approval.
- Immediately before each branch or PR read-modify-write and each retry,
  revalidate current packet trust from a current signed offline status snapshot
  or an exact `provider_trust_status_read` grant, plus the exact
  product-authority bundle inside Lumyn. An earlier validation result is never
  reusable, undeclared status egress is forbidden, and Factory dispatch does
  not authorize the remote write.
- Create or update one authorized branch and one draft PR idempotently.
- Bind the PR to packet, base/head, plan, patch, and evidence hashes.
- Mark stale, superseded, withdrawn, closed, and duplicate states.
- If—and only if—the consumer grants `provider_attestation`, emit the
  separately consented campaign status as an independent action. It is signed
  by the consumer-controlled configured signer, bound to the acknowledged
  minimal connection receipt, packet, exact evidence commitments, audience,
  consent-policy digest, issue time, expiry, and anti-replay nonce. Missing
  richer reporting authority does not block branch or PR creation. Export is
  irreversible and revocation stops only future attestations.
- Include M8 sandbox evidence when present and current, but preserve the exact
  weaker M7 evidence label otherwise.

**Repo paths:**

- `internal/report/`
- `internal/github/`
- `internal/attestation/`
- `internal/trust/`
- `schemas/`
- `.github/`
- `tests/`
- `docs/`

**Run commands:**

- `make test-fast`
- `make test-coverage`
- `make test-contracts`
- `go run ./cmd/lumyn --json authorization validate --bundle "$LUMYN_PRIVATE_PRODUCT_AUTHORITY_BUNDLE" --task M9`
- `make prepush-full`
- approved GitHub integration test recorded in task evidence

**Test requirements:**

- Mock GitHub tests for permission denial, default-branch write, duplicate PR,
  stale base, superseded packet, and idempotent update.
- Time-of-check/time-of-use tests revoke or expire branch, PR, retention, or
  deletion authority—or revoke, withdraw, supersede, rotate, replay, or mutate
  the packet—between validation and each remote write; no branch or PR mutation
  occurs. Separately revoke reporting authority and prove the PR still opens
  while attestation does not.
- PR body golden tests covering every required evidence section.
- Provider-attestation redaction and allowlist tests.
- Consumer-issuer signature, receipt/packet/evidence binding, audience,
  consent-policy digest, issue-time/expiry, acknowledgement, and replay tests
  for the minimal connection receipt and any richer provider attestation.
- Disclosure tests prove revocation prevents future export and deletes
  Lumyn-controlled private copies without claiming to recall provider copies.
- A no-sandbox/no-reporting integration test produces the correctly labeled
  evidence-backed draft PR; optional current sandbox evidence is included
  without becoming a prerequisite.

**Matrix wiring:** Tiers 1, 2, 3, 4, 5, 9, 11, and approved Tier 12; CodeQL and
security review.

**Acceptance criteria:**

- PR creation is an explicit action distinct from impact and patching.
- Only a draft PR on a non-default branch is possible.
- Provider status contains no unconsented source or private evidence.
- A sponsored connection is billable only from an acknowledged valid minimal
  connection receipt; richer status is independently signed and authorized.
- Repeated delivery cannot create duplicates.
- The validation command is closure and diagnostic proof; the authoritative
  product gate executes inside Lumyn immediately before each remote side
  effect.

**Changelog impact:** required
**Changelog section:** Added
**Draft changelog entry:** Add explicit, idempotent GitHub draft-PR delivery with complete migration evidence, signed minimal connection receipts, and separately optional consented provider attestations.
**Semver marker override:** `[semver:minor]`
**Contract/API impact:** Adds draft-PR and campaign-attestation surfaces.
**Versioning/migration impact:** PR/evidence binding invalidates on packet or base/head change.
**Architecture constraints:** Separate consumer data plane and provider control plane; least-privilege GitHub permissions.
**ADR required:** yes for GitHub and data-sharing boundaries.
**TDD first failing test(s):** Default-branch target, duplicate run, and raw-log attestation are rejected.
**Cost/perf impact:** medium
**Chaos/failure hypothesis:** Retries or stale branches create duplicate/conflicting PRs or leak private evidence.
**Semantic invariants:** Consumer controls PR creation and merge; provider receives only consented status.

---

## Epic 7 — Provider Campaign Pilot And Outcome Learning

### M10: Run one qualified design-partner campaign

**Priority:** P0 for product validation
**Risk class:** High
**Blocked by:** `M2.5`, `M9`, consumer consent, every exact Lumyn product grant
required by each selected action, and the Factory worker's task-scoped
`approval`, `credentials`, and `network` grants. M8 plus `sandbox_network`,
`sandbox_credential`, and `sandbox_request_disclosure` are conditional on a
repository selecting sandbox proof. A distinct `campaign_receipt` grant is
required for a repository to count in the provider-sponsored cohort.
`provider_attestation` is action-specific and optional; the campaign may use
exact-patched-head mock proof and consumer-private outcome evidence when richer
reporting is not consented.

**Acceptance item IDs:** `DISC-001`, `DISC-002`, `AUTH-002`, `AUTH-003`,
`AUTH-004`, `AUTH-005`, `PILOT-001`, `PILOT-002`, `PILOT-003`, `PILOT-004`,
`PILOT-005`, `PILOT-006`, `PILOT-007`, `PILOT-008`, `PILOT-009`

**Tasks:**

- Advance the provider and frozen cohort qualified in M2.5 without changing
  thresholds after outcomes are visible.
- Before every repository, host-command, provider-status, sandbox, minimal
  campaign-receipt, richer provider-attestation, GitHub, retention, or deletion
  side effect, have Lumyn revalidate the current packet trust state from a
  current signed offline status snapshot or exact
  `provider_trust_status_read` grant and every exact product grant required by
  that action. Factory worker
  approval, credentials, and network access govern only the
  implementation/pilot worker and never substitute for this live-action gate.
  Do not request or require sandbox or `provider_attestation` grants when the
  corresponding action is not performed.
- Author and validate the provider packet and canary fixtures.
- Invite at least five distinct prequalified eligible repositories across at
  least three Consumer Organizations, counting one funnel unit per repository.
- Record the full repository-based authorization-to-merge funnel.
- Require at least three organizations to explicitly issue authorization and
  import a provider-signed acknowledgement for a consumer-signed minimal
  connection receipt whose signer is provider-authenticated and whose
  invitation unit maps to one distinct repository each within 30 days, three
  draft PRs for distinct
  repositories within 14 calendar days of authorization with no paused clock,
  and three merges from distinct repositories and organizations within 45 days
  of PR creation.
- Across connected repositories, require median calendar time from invitation
  receipt to the first valid impact report to be at most seven days and median
  total Consumer Organization security, privacy, platform, and maintainer
  hands-on labor over that interval to be at most two hours.
- Require median Consumer Maintainer in-product hands-on time from starting
  `campaign accept` through explicit authorization issuance to the first valid
  impact report to be at most 60 minutes. Record Lumyn and provider assistance
  separately; bespoke product changes or provider source access invalidate the
  onboarding observation.
- Require at least two of the first three merges to need no substantive manual
  correction from PR creation through the fixed 30-day post-merge observation
  window; a migration-attributable revert or fix counts as correction. At
  least two of those PRs must also reach `workflow_verified_mock` or
  `workflow_verified_sandbox` from the exact patched head before merge.
- Measure setup, onboarding, support, repository conversion, automation, lead
  time, contribution, cost-per-merge, and actual legacy-retirement duration
  against the frozen baseline.
- Pass the one preregistered primary provider outcome: at least 20% lower
  support hours per merged repository, at least 20% lower
  authorization-to-merge lead time, or at least 60% of the frozen eligible
  cohort retired from the targeted legacy version by day 120. Missing or
  `not_measurable` primary evidence fails the campaign.
- Record merged, closed, reverted, and corrected outcomes with provenance.
- By the frozen absolute campaign judgment deadline, obtain both a paid pilot
  invoice at least equal to frozen Lumyn campaign COGS and an executed annual
  connected-repository or second-named-migration purchase. Preregister at least
  60% projected recurring gross margin and at most four Lumyn operator hours
  per merged repository after campaign setup. Nominal payments, LOIs, a second
  bespoke services order, and informal willingness to pay do not qualify.

**Repo paths:**

- `.factory/artifacts/pilot/lumyn-migration-mvp/public/`
- `scripts/`
- `schemas/`
- `tests/`
- `docs/product/`
- customer-private evidence paths approved per participant
- no committed private source, credentials, raw logs, or traces

**Run commands:**

- `make prepush-full`
- `go run ./cmd/lumyn --json authorization validate --bundle "$LUMYN_PRIVATE_PRODUCT_AUTHORITY_BUNDLE" --task M10`
- `python3 scripts/validate_pilot_evidence.py --self-test`
- `python3 scripts/validate_pilot_evidence.py --attestation "$LUMYN_PRIVATE_PILOT_ATTESTATION" --public-manifest .factory/artifacts/pilot/lumyn-migration-mvp/public/pilot-manifest.json`
- approved provider/consumer commands recorded without secrets

**Test requirements:**

- Human approval and external evidence for each product-signal item, followed
  by independent `evidence-attestor` verification of the actual private
  attestation and public aggregate/hash-only manifest.
- Independent `evidence-attestor` verification of the frozen cohort, source
  bindings, consumer-signed authorization and acknowledged connection receipts,
  funnel, workflow, correction, economics, and payment calculations inside the
  privacy-approved pilot attestation. No benchmark holdout root is reused for
  campaign data.
- Privacy review of every provider-visible field.
- Per-repository proof-of-behavior scorecard.
- Aggregate funnel and economics calculations with source refs.
- Negative fixtures for duplicate repository units, organization concentration,
  missing/invalid/unacknowledged/replayed connection receipts, forged consumer
  issuer, receipt fields outside the allowlist, recallable-disclosure claims,
  invitation-to-impact time above seven days, total consumer activation labor
  above two hours, in-product maintainer time above 60 minutes, excessive
  Lumyn/provider onboarding time, missed or paused deadlines,
  late/reverted outcomes, missing patched-head workflow proof, negative
  contribution, recurring gross margin below 60%, missing or unmeasurable
  primary provider outcomes, nominal payment, one-off services continuation,
  and missing paid continuation.

**Matrix wiring:** Tiers 4, 10, 11, and approved Tier 12; product, privacy,
security, and evidence-attestation gates.

**Acceptance criteria:**

- Every `PILOT` item has direct evidence; frozen thresholds cannot be
  rebaselined after the first invitation, and abandonment or timeout counts as
  failure at the absolute campaign judgment deadline.
- Technical success is not substituted for consent, merge, or payment.
- No repository counts as connected or billable without an acknowledged valid
  minimal connection receipt; richer reporting remains separately optional.
- The pilot has non-negative campaign contribution, at least 60% projected
  recurring gross margin, bounded Lumyn operator effort, and executed
  repeatable paid continuation—not a nominal or bespoke services payment.
- The preregistered primary provider outcome clears its material threshold;
  measurement-only or `not_measurable` evidence does not pass.
- Private consumer artifacts are not committed or shared with the provider by
  default.
- Provider/public disclosures are consented as irreversible; revocation blocks
  future sharing and deletes Lumyn-controlled private copies without promising
  recall.
- Validation-command output is closure proof, not cached runtime authority;
  every live action is authorized at its side-effect boundary.

**Changelog impact:** not required for evidence-only operation; required for any product behavior changed during pilot
**Changelog section:** none unless behavior changes
**Draft changelog entry:** none
**Semver marker override:** none
**Contract/API impact:** No new behavior may be invented inside the pilot; findings return through governed follow-up work.
**Versioning/migration impact:** Packet and recipe changes require new versions and revalidation.
**Architecture constraints:** Manual coordination is allowed; privacy, authorization, proof, and evidence contracts are not.
**ADR required:** no unless the pilot requires a new trust boundary.
**TDD first failing test(s):** Not applicable to external signals; any code repair follows its own red-first task.
**Cost/perf impact:** high
**Chaos/failure hypothesis:** A technically successful pilot masks poor consent conversion or service-only willingness to pay.
**Semantic invariants:** Product signals remain external evidence; no threshold is silently waived; outcome learning cannot mutate active rules.

The preregistered pilot uses these exact definitions:

- `substantive_manual_correction`: a human edit from PR creation through the
  fixed 30-day post-merge observation window that changes the migrated API/SDK
  invocation, request/response mapping, error handling, workflow behavior, or
  Lumyn-generated semantic edit. A migration-attributable post-merge revert or
  fix counts. Formatting, deterministic lockfile normalization, comments, and
  unrelated pre-existing CI repair do not count, but remain recorded.
- `campaign_setup_hours`: Lumyn plus provider time for packet authoring, canary
  fixtures, and cohort preparation before the first invitation.
- `median_invitation_to_first_impact_days`: median calendar time from a
  Consumer Organization's receipt of the invitation to its first valid impact
  report.
- `consumer_activation_hands_on_hours`: total Consumer Organization security,
  privacy, platform, and maintainer hands-on labor from invitation receipt to
  the first valid impact report.
- `consumer_maintainer_in_product_hours`: Consumer Maintainer hands-on time
  from starting `campaign accept` through explicit authorization issuance to
  the first valid impact report.
- `lumyn_onboarding_hours`: Lumyn operator time from invitation to the first
  valid impact report, excluding core product development.
- `provider_onboarding_hours`: provider operator time from invitation to the
  first valid impact report.
- `provider_support_hours_per_merged_repo`: provider DX/support/solutions hours
  in the measurement window divided by merged PRs.
- `median_authorization_to_merge_days`: median calendar days from repository
  authorization to merge for the frozen cohort.
- `automation_rate`: merged PRs without substantive correction divided by
  `affected_supported` repositories.
- `eligible_to_merge_conversion`: distinct merged eligible repositories divided
  by the frozen eligible-repository cohort, with at most one unit per
  repository.
- `acknowledged_connected_repositories`: distinct signature-valid,
  non-replayed minimal connection receipts acknowledged by the provider after
  verifying invitation, packet, opaque organization/repository IDs,
  consent-policy digest, issuer, audience, expiry, and nonce. This is the only
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
  minus projected recurring Lumyn COGS, divided by that revenue.
- `lumyn_operator_hours_per_merged_repo`: Lumyn onboarding, operations, and
  support hours after campaign setup divided by merged PRs.
- `legacy_retirement_delta`: current actual retirement duration from packet
  publication or migration announcement to actual legacy retirement, minus the
  comparable prior migration's actual duration. If the current surface has not
  retired by the frozen judgment deadline or no reliable comparator exists,
  record `not_measurable` and claim no improvement.
- `legacy_cohort_retirement_rate`: distinct frozen eligible repositories whose
  resolved dependency state no longer includes the targeted legacy SDK/API
  version by day 120, divided by the frozen eligible-repository cohort.
- `primary_provider_outcome`: exactly one of
  `provider_support_hours_per_merged_repo`,
  `median_authorization_to_merge_days`, or
  `legacy_cohort_retirement_rate`, selected before the first invitation. The
  first two require at least 20% improvement against a comparable frozen
  baseline; the third requires at least 60%. A missing comparator makes the
  first two ineligible, and a missing or `not_measurable` selected outcome
  fails `PILOT-007`.

---

## Minimum-Now Sequence

### Wave 1

- M0 rebaseline compatibility first.
- After M0, run M1 benchmark corpus and M2 migration/authorization plus
  two-sided-activation contracts in parallel.

### Wave 2

- After M2, run M2.5 provider/cohort qualification with the standard campaign
  kit while M3 source/target semantic change intake proceeds after M1 and M2.
- M4 TypeScript impact analysis.

### Wave 3

- M5 reviewable migration plan only after both M4 and M2.5 close.
- M6 deterministic patch engine.

### Wave 4

- M7 repository/replay verification and offline end-to-end canary.

### Wave 5

- Run M8 approved sandbox verification and M9 evidence-backed draft PR delivery
  as independent branches after M7. M9 does not wait for M8.

### Wave 6

- M10 qualified provider campaign and commercial evidence after M9. M8 evidence
  is included only for repositories that separately authorize sandbox proof.

Do not build a substantial hosted campaign dashboard before a provider commits
a real change and reachable cohort.

The exact dependency graph is:

```text
M0    -> none
M1    -> M0
M2    -> M0
M2.5  -> M0, M2
M3    -> M1, M2
M4    -> M3
M5    -> M2, M2.5, M4
M6    -> M5
M7    -> M2, M6
M8    -> M7
M9    -> M7
M10   -> M2.5, M9
```

---

## Explicit Non-Goals

- Generic buy-side monitoring of every API dependency.
- Anonymous scanning of downstream repositories.
- Model-provider panels or customer-facing live agent evaluation.
- Arbitrary model-assisted patching.
- Authentication, webhook, event, GraphQL, gRPC, or production migrations.
- Multi-language or broad package-manager support.
- Generated-client regeneration.
- Provider-supplied executable migration scripts.
- Production credentials or mutations.
- Default-branch writes or automatic merge.
- Provider access to raw consumer data without explicit consent.
- Claiming external/public disclosures can be recalled after export or commit.
- Calling the current repository or pilot distribution OSS before the explicit
  public-release gate closes.
- Hosted dashboard as a prerequisite.
- Treating benchmark performance as product-market proof.

---

## Definition Of Done

The plan is done only when:

- every required PRD acceptance item is represented in the active acceptance
  ledger and mapped to a bounded task or explicit external evidence gate;
- all technical items have command, fixture, schema, or proof-of-behavior
  evidence;
- all consumer authorization and data-sharing items pass;
- required CI, coverage, CodeQL, review, commit/push, PR lifecycle, and
  post-merge evidence exists;
- no active task uses a deprecated worker alias;
- no historical item is silently reinterpreted;
- all 62 PRD acceptance items are represented exactly once in the active
  ledger and mapped to bounded work or external evidence;
- the fixed corpus has zero false-positive and zero false-negative canonical
  workflow-verification results;
- the design-partner pilot meets its activation, consent, PR, merge,
  patched-head workflow, correction, recurring-economic, and paid-continuation
  gates;
- every counted sponsored connection has an acknowledged valid minimal signed
  receipt, every repository command used a supported host-isolation backend,
  and every packet side effect had authorized current-status evidence;
- pilot distribution has explicit terms, security/support contacts, and
  signed install-integrity evidence without an unsupported OSS claim;
- README, PRD, plan, architecture, developer guidance, Factory profile, active
  task packets, validation contract, acceptance mapping, and scope closure
  agree.
