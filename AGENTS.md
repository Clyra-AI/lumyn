# AGENTS.md — Lumyn Repository Contract

Version: 2.0
Status: Normative
Scope: This repository only.

## 1. Scope And Intent

- Build Lumyn as a provider-sponsored, customer-controlled API migration
  system.
- Treat `docs/product/prd.md` as the product source of truth.
- Treat `docs/product/plan.md` as the human-readable active implementation
  plan.
- Treat `.factory/artifacts/prd-to-plan/lumyn-migration-mvp/` as the active
  compiled planning and acceptance contract.
- Keep the task/control artifacts under
  `.factory/artifacts/prd-to-plan/lumyn-mvp/` immutable as historical evidence
  from the superseded workflow/eval product direction. Its README may carry
  only a non-operative dispatch tombstone; it cannot change historical claims.
- Keep Factory run evidence under `.factory/artifacts/`.
- Keep independent review, holdout, trace-grade, attestation, shipping, and PR
  lifecycle evidence lifecycle-owned. `task-executor` may not write
  `.factory/artifacts/lifecycle-evidence/` or
  `.factory/artifacts/pr-lifecycle/`.
- Keep Factory scratch and daemon state under `.factory/tmp/` and `.factoryd/`.
- Keep consumer-private product runtime artifacts in an explicitly configured,
  non-committable root outside the consumer checkout and every public source
  repository.

## 2. North Star

Every product change should improve one or more of:

- signed, trust-bound provider change intent;
- low-friction provider publication and consumer campaign acceptance through
  standard signed artifacts;
- customer-authorized repository impact coverage;
- bounded, explainable migration patches;
- baseline-aware repository verification;
- proof-honest workflow outcome evidence;
- consumer-controlled repository, command, credential, network, PR, and
  disclosure permissions;
- draft-PR reviewability and merge confidence;
- provider migration-campaign completion;
- fail-closed handling of unsupported or ambiguous integrations.

## 3. Product Authorities: Two Principals, Two Authorities

Keep the two principals separate:

- The API Provider owns API/SDK intent, packet publication, campaign cohort,
  compatibility window, sandbox semantics, and rollback guidance.
- The API Consumer Organization owns repository access, commands, credentials,
  execution, disclosure, branch policy, review, and merge.

Provider payment or campaign sponsorship never grants consumer repository
authority. Consumer participation never lets Lumyn rewrite provider intent.

Use explicit terms:

- `api_provider` or `change_authority` for the API seller;
- `model_provider` for any future model endpoint;
- `api_consumer_organization` for the repository-owning organization;
- `consumer_maintainer` for the human with review and merge authority.

Do not use bare `provider` where the meaning could be ambiguous.

## 4. Non-Negotiable Product Constraints

- Analyze only explicitly authorized repositories.
- Never claim coverage of all downstream integrations.
- Provider change packets are declarative and cannot execute arbitrary
  provider-supplied scripts.
- Published packets are immutable for their authorized audience and are trusted
  only when their signature verifies against a consumer-pinned API-provider
  trust root and verified provider-to-package ownership binding. Key,
  timestamp, audience, expiry, rotation, revocation, withdrawal, and replay
  checks must pass immediately before every local, sandbox, or remote side
  effect. A prior trust decision must not be cached across a write boundary.
  Current status comes only from a signed offline snapshot inside the pinned
  maximum age or an exact endpoint read under
  `provider_trust_status_read`; missing, stale, replayed, or undeclared status
  access blocks.
- Initial provider enrollment requires a separately authenticated enrollment
  bundle and out-of-band confirmed fingerprint. Never let an invitation
  authenticate its own root. Normal rotation is signed by the active root;
  emergency recovery requires explicit re-enrollment and freezes open
  campaigns. Enrollment also pins the provider status signer,
  receipt-acknowledgement signer, and permitted status/receipt exchange classes;
  an invitation may narrow but never replace them.
- `provider enroll` and `campaign accept` accept no repository argument, run
  with the checkout unavailable, and store only consumer-private state. They
  create no authority. `authorization issue` and `authorization revoke` are
  explicit consumer-signed actions and perform no authorized side effect.
- Impact analysis is read-only.
- A Consumer Maintainer must review the migration plan before write approval.
- Repository read, local write, host-isolated command execution,
  provider-status read, package-registry network, sandbox request disclosure,
  sandbox network, sandbox credential, remote branch write, PR write, minimal
  campaign receipt, richer provider reporting, retention, and deletion scopes
  remain independent.
- Apply patches only in an isolated workspace within approved path and diff
  budgets.
- Do not infer missing business values.
- Unsupported or ambiguous cases fail closed as `needs_input`, `unsupported`,
  `uncertain`, or `blocked`.
- Repository tests run without network by default.
- Dependency lifecycle scripts require separate consumer approval.
- Automatic `package-lock.json` mutation requires exact Node and npm versions,
  a registry endpoint or immutable snapshot, package-integrity inputs, a
  toolchain digest, disabled lifecycle scripts, and separate registry-network
  authorization.
- Production credentials and production mutations are prohibited in the MVP.
- Sandbox credentials are isolated from general build and test commands.
- Sandbox consent names transmitted payload classes, synthetic or approved
  non-sensitive test data, provider logging, retention and deletion terms,
  budgets, namespace, cleanup, and orphan handling. Production customer data,
  PII, credentials, and secrets are prohibited.
- Redact secrets before persistence or sharing.
- Raw consumer code, diffs, logs, traces, prompts, responses, and credentials
  are not provider-visible by default.
- Private runtime artifacts enforce authorization TTL and deletion on expiry or
  revocation automatically on creation, read, startup, and the next run.
  `lumyn artifacts gc` is the operator recovery path for retrying deletion and
  inspecting orphan records; cleanup must not rewrite historical closure.
  Provider-visible attestations contain only exact
  consumer-consented fields. Public evidence contains only consented redacted
  aggregates or hashes.
- Open draft PRs only; never write to the default branch or auto-merge.
- Use only the canonical successful-verification labels `static_verified`,
  `repo_verified`, `workflow_contract_replay_passed`,
  `workflow_verified_replay`, `workflow_verified_mock`, and
  `workflow_verified_sandbox`. A `workflow_verified_*` label requires an
  approved entrypoint executed from the exact patched repository head plus
  observed interaction and outcome evidence in that environment. Independent
  contract or cassette replay cannot exceed `repo_verified`.
- Unimplemented commands must return a typed nonzero result.
- The standard offline activation canary must exercise campaign-kit creation,
  packet publication, invitation acceptance, explicit synthetic
  consumer-signed authorization issuance, impact, planning, patching,
  host-isolated verification, evidence, and a local draft-PR preview. A missing
  stage is a typed nonzero result, never a demonstration success.
- M8 sandbox verification and M9 draft-PR delivery are independent after M7.
  Missing sandbox or `provider_attestation` authority cannot block an
  otherwise authorized M9 draft PR; it disables only the corresponding
  optional action.
- Do not make generic live agent eval, model-provider panels, or agent readiness
  required MVP scope.

## 5. Initial MVP Boundary

- GitHub repositories.
- One explicitly selected package root.
- TypeScript source discoverable through `tsconfig.json`.
- One official npm SDK dependency.
- `package-lock.json` is the first writable lockfile, subject to the exact
  toolchain, registry/snapshot, package-integrity, lifecycle-script, and
  registry-network controls above.
- Direct imports are patchable.
- One-hop wrappers are detectable but may require maintainer input.
- Exactly three deterministic change classes:
  - method or operation rename;
  - request-property rename or relocation;
  - response-property rename or relocation.
- Authentication, webhook/event, GraphQL, gRPC, generated-client,
  cross-language, and production migrations are out of scope.
- Model-assisted patching is out of scope until deterministic behavior and the
  same verification gates are proven.

## 6. Required Boundaries

- `docs/product/`: product requirements, plan, and item-level scope source.
- `docs/dev/`: repo-local engineering and validation guidance.
- `docs/architecture/`: architecture, trust boundaries, ADRs, and findings.
- `.factory/artifacts/prd-to-plan/lumyn-migration-mvp/`: active planning control
  truth.
- `.factory/artifacts/prd-to-plan/lumyn-mvp/`: historical non-active plan.
- `.factory/artifacts/task-runs/`: task-owned validation and work proof.
- `.factory/artifacts/pr-lifecycle/`: validation, CI, review, ship, merge, and
  post-merge evidence.
- `.factory/tmp/`: ignored local scratch space.
- `.factoryd/`: ignored daemon state, claims, worktrees, events, and run reports.
- `schemas/`: versioned executable artifact contracts.
- `cmd/lumyn/`: CLI entrypoint and process result.
- `internal/source/`: source parsing only; do not add unrelated migration
  domains.
- Future `internal/change/`, `internal/trust/`, `internal/impact/`,
  `internal/patch/`, `internal/verify/`, `internal/github/`, and
  `internal/attestation/`: distinct product boundaries.
- `workflows/`, `cassettes/`, and `runs/`: retained synthetic or licensed
  project fixtures only; consumer-private runtime instances stay outside the
  checkout.
- Future `changes/`, `authorizations/`, `impacts/`, and `migrations/`:
  schema-backed product artifact types introduced only by their planned tasks;
  private instances stay outside the checkout.

## 7. Trust And Capability Gates

Deterministic public-fixture work defaults to:

- no ambient secrets;
- no live network;
- no customer repositories;
- no provider sandbox;
- no GitHub writes.

Live or external work requires task-scoped grants:

- `customer_repo_read`: repository, readable paths, expiry, independently
  consented retention and deletion, and evidence owner;
- `customer_repo_write`: approved plan digest, writable paths, local isolated
  workspace, diff budget, expiry, and rollback; it does not authorize a remote
  branch;
- `command_execution`: exact commands, lifecycle-script posture, timeout,
  output budget, network posture, exact read/write mounts, neutral home/temp
  roots, executable roots, environment classes, local-socket and inherited-FD
  policy, process-tree inheritance, OS-credential denial, and a supported
  fail-closed isolation backend;
- `provider_trust_status_read`: exact provider/package, pinned status endpoint,
  request shape, response budget, maximum age, expiry, and proof that no
  repository or consumer data is transmitted; a current signed offline
  snapshot needs no network grant;
- `package_registry_read`: exact registry endpoint or immutable snapshot,
  package and integrity inputs, exact Node/npm versions, toolchain digest,
  lifecycle-script posture, expiry, and read-only budget;
- `sandbox_request_disclosure`: exact payload classes, synthetic or approved
  non-sensitive test-data classes, provider logging, retention and deletion
  terms, and evidence;
- `sandbox_network`: API-provider identity, exact non-production endpoint and
  operation allowlist, namespace, request/write budgets, cleanup, orphan
  handling, and expiry; it does not imply credentials or payload disclosure;
- `sandbox_credential`: non-production credential class and scopes, isolated
  injection stage, expiry, and revocation; it does not imply network or payload
  authority;
- `github_branch_write`: repository, authorized non-default branch namespace,
  base commit, token expiry, and rollback;
- `github_pr_write`: repository, authorized remote branch, base branch,
  draft-only posture, token expiry, idempotency key, and approved plan/evidence
  refs; it does not imply remote branch write;
- `campaign_receipt`: minimal connection-receipt field allowlist, invitation
  and packet binding, eligible-repository unit nonce, opaque per-campaign
  organization/repository IDs, provider-authenticated consumer signer binding,
  provider acknowledgement key, pinned endpoint or offline exchange,
  one-unit cardinality, idempotent deduplication, audience, consent-policy
  digest, nonce, expiry, and provider-signed acknowledgement;
- `provider_attestation`: exact richer consumer-consented fields, consumer
  signer, audience, expiry, anti-replay binding, and retention;
- `artifact_retention`: exact artifact classes, storage boundary, TTL, expiry
  behavior, and evidence owner;
- `artifact_deletion`: exact revocation/expiry triggers, deletion scope,
  receipt owner, retry posture, and orphan-report route.

Generic network, credential, or repository approval does not satisfy a more
specific grant. Retention and deletion are independently consented wherever
data is produced or disclosed. Wildcard grants are prohibited for customer
repository, provider-status, registry, sandbox, GitHub write,
campaign-receipt, provider-attestation, retention, and deletion work.
Provider export and public commit are irreversible disclosure boundaries:
revocation prevents future sharing and deletes Lumyn-controlled private copies,
but cannot promise recall from recipients, Git history, clones, or caches.
Factory credential scopes and network allowlists likewise reject semantic
wildcards such as `all`, `any`, or `default`, case-insensitive duplicates,
wildcard hosts, CIDRs, and unspecified addresses.

These are Lumyn product authorizations, stored as private schema-backed
artifacts under the configured consumer-controlled state root. They are not
`.factory/factoryd.json` capability names. Factory uses its closed
`approval`, `credentials`, and `network` vocabulary to govern the
implementation worker; the `approval` evidence cites a validated opaque
product-authorization bundle. Factory grants never widen or replace Lumyn
product authority. For live product behavior, Lumyn revalidates the current
packet and exact applicable bundle inside the product immediately before every
side effect and retry. `lumyn authorization validate` is closure and diagnostic
proof; Factory dispatch does not implement this product gate and cannot confer
its authority.

## 8. Evidence And Proof Rules

- Keep impact, patch, verification, delivery, permission, and residual-risk
  axes separate.
- Bind migration evidence to the packet digest, base/head commits, plan digest,
  patch digest, environment, and artifact hashes.
- Invalidate evidence when a bound input changes.
- Capture pre-existing repository failures before patching.
- Label static, repository, replay, mock, sandbox, and production evidence
  separately.
- Use the canonical verification labels in Section 4 and bind any
  `workflow_verified_*` result to causal execution from the exact patched head.
- Production evidence is outside the MVP.
- Do not call the current repository or pilot distribution OSS. Design-partner
  delivery requires explicit terms, security/support routes, signed provenance,
  checksums, and install integrity; public OSS/self-serve additionally requires
  an approved license and contribution/vulnerability-response policies.
- Cleanup failure, orphan evidence, boundary violation, redaction uncertainty,
  stale evidence, or missing read-back prevents workflow verification.
- Large logs and traces are cited by a stable artifact reference and hash, not
  duplicated into planning or PR metadata. Factory refs are repo-relative;
  private runtime refs are relative to the configured private root.
- Consumer-private runtime evidence remains outside the checkout and public
  source. Provider-visible attestations contain only exact consumer-consented
  fields; public evidence is aggregate/hash-only. Provider export and public
  commit are irreversible, so revocation applies only to future sharing and
  Lumyn-controlled private copies.
- Only a non-resolving opaque holdout manifest may be committed: opaque case
  IDs, provenance class and license posture, a frozen suite commitment, and
  encrypted or HMAC artifact commitments. Resolving provenance, plaintext
  content digests, held-out inputs, answer keys, expected patches or labels,
  and raw traces live in `LUMYN_HOLDOUT_ROOT`, which is provisioned by an
  independent holdout owner, available only to `holdout-evaluator`, and absent
  from task-executor mounts, prompts, and environment.
- M1 uses provision-mode holdout policy and creates the keyed suite commitment;
  M4/M6/M7 evaluate-mode policy resolves and byte-binds M1's trusted result.
  Never fabricate a future holdout commitment in a static task packet.
- Independent lifecycle artifacts schema-validate and bind the exact task, work
  item, lifecycle run, current validation run, candidate digest, and work-proof
  marker digest. Their worker provenance must be independent of
  `task-executor`; stale, replayed, mismatched, or worker-authored evidence
  blocks `commit-push`.
- Historical evidence may prove only the exact old behavior it recorded; do not
  rewrite old artifacts to fit the new product.

## 9. Required Validation

For normal changes, run:

- `make lint-fast`
- `make test-fast`
- `make test-coverage`
- `make test-contracts`

Before PR or merge, run:

- `make prepush-full`

If any command is skipped, the validation report records the reason.

GitHub Actions `validate` runs the same full gate on pull requests and `main`.
CodeQL Go analysis is the security scanner risk lane. Coverage-gated work
requires `make test-coverage` evidence or an approved scoped exception.

Passive Codex review settle is required before merge. Green CI alone is not
merge-ready when Codex review is enabled. Do not merge manually through
`gh pr merge`, the GitHub UI, or a connector until the latest PR head has the
configured terminal Codex review evidence.

GitHub `main` must be protected by branch protection plus the
`protect-main-from-direct-push` ruleset. Required controls include pull
requests, strict `validate` and `CodeQL analyze` checks, admin enforcement,
conversation resolution, and no force-push or deletion bypass. Verify the live
state with `make audit-remote-protection` when GitHub credentials are available.

The PR lifecycle report path is:

```text
.factory/artifacts/pr-lifecycle/<work_item_id>/pr-lifecycle-report.json
```

## 10. Runtime And Distribution Pins

- Language: Go.
- Go version: `1.26.5`.
- Module: `github.com/Clyra-AI/lumyn`.
- Initial distribution: standalone local runner.
- Secondary distribution: public OSS/self-serve and Homebrew only after the
  separate approved license, security, contribution, support,
  vulnerability-response, and release-integrity gate.
- Consumer execution: local or consumer-controlled CI by default.
- Target consumer ecosystem: TypeScript/Node and one official npm SDK.
- Deterministic transforms first.
- No required model-provider endpoint for MVP.
- Factory artifact namespace: `.factory/artifacts/`.
- Committable public product artifact namespaces remain separate and
  repo-relative. Consumer-private product references are relative to the
  configured private root and cross disclosure boundaries only as opaque IDs
  and digests.

Changing runtime, module path, primary distribution, execution plane, target
language, provider/consumer authority, credential posture, or active planning
path requires an ADR or explicit decision update before implementation.

## 11. Factory And factoryd Operation

Factory is the shared contract source. `factoryd` may execute active task packets
from `.factory/artifacts/prd-to-plan/lumyn-migration-mvp/`.

M2.5 is a product-signal gate, not an engineering placeholder. Migration-plan
implementation beyond contracts and read-only impact analysis must not proceed
until `DISC-001` and `DISC-002` have direct external evidence for one qualified
API Provider commitment, a supported-class canary, and a frozen cohort of five
distinct eligible repositories across at least three Consumer Organizations.
The provider must also demonstrate recurring-value potential through a second
consequential migration/deprecation within 12 months or at least 20 named
managed integrations eligible for an annual connected-repository program.
Qualification also requires an operational signed provider-status channel, a
pinned provider receipt-acknowledgement key plus endpoint or offline exchange,
and a supported OS/architecture with enforceable host isolation for every
candidate environment. Canary inspection qualifies the opportunity; it does
not close runtime packet trust. M3 owns executable `CHG-003` and `CHG-004`
enforcement before any packet can drive classification or mutation.

Before M2.5 collects, stores, or discloses identifiable external evidence, a
narrow manual privacy/legal preflight must be approved for the Factory
implementation worker. It records allowed private fields, participant consent,
the approved external private root, TTL, expiry/revocation deletion,
deletion-receipt and orphan ownership, the minimal connection-receipt
allowlist, provider-authenticated consumer signer binding,
provider-signed acknowledgement/cardinality policy, and a separate public
field allowlist with separate disclosure consent. Consent states that provider
export and public commit cannot be recalled. This preflight authorizes M2.5
evidence handling only. It is not a Lumyn runtime product grant. The active
Factory approval must cite the preflight and exactly match its canonical scope
digest; a generic approval or
stale digest is invalid.

The active PRD-derived acceptance ledger contains exactly 62 item-level closure
units. Tasks, waves, and an overall MVP label cannot substitute for direct
evidence against each applicable acceptance item ID.

Runner-ready packets must include:

- exact acceptance item IDs;
- dependencies and gates;
- allowed and forbidden paths;
- architecture target paths;
- validation commands;
- worker and lifecycle chain;
- actual private/public evidence-validation commands and an
  independent `evidence-attestor` record for product-signal closure;
- required `holdout_result`, `trace_grade_report`, and `attestation_record`
  lifecycle refs when their independent gates apply;
- evidence requirements;
- proof level;
- runtime pins;
- capability requirements;
- stop conditions;
- changelog, versioning, migration, and docs intent;
- semantic invariants.

Product workers may write task-scoped evidence but must not mutate active
PRD-derived task DAG, acceptance ledger, mapping, validation contract, or scope
closure to bypass a gate. They also must not write independent lifecycle
evidence or PR-lifecycle evidence.

The canonical implementation-to-merge chain is:

1. `task-executor`
2. `validation-gate`
3. `code-review` when risk requires it
4. `holdout-evaluator` when selected by policy
5. `trace-grader` when selected by policy
6. `evidence-attestor` when selected by policy
7. `commit-push`
8. `post-merge-monitor`
9. `repair-feedback` when a gate fails

The independent workers are external or human-operated lifecycle gates. They
must produce task-bound, current-work-proof, passing artifacts in the trusted
Factory evidence root before `commit-push`; an implementation worker cannot
self-grade or self-attest them.

Task-level `required_review` is canonical. Its review lens and reviewer class
must exactly equal the inherited validation contract and the schema-valid
current review artifact checked by shipping.

Do not use deprecated lifecycle aliases in active chains.

## 12. Stop Conditions

Stop and request a human decision when:

- provider and consumer authority are conflated;
- a provider packet would execute code, lacks a valid pinned trust root or
  provider-to-package binding, or fails key, lifecycle, audience, expiry,
  revocation, withdrawal, or replay checks;
- repository access lacks a specific active grant;
- a read-only phase would mutate state;
- the approved plan no longer matches the patch inputs;
- current packet trust or product authority cannot be revalidated immediately
  before a local, sandbox, or remote side effect;
- neither a current signed provider-status snapshot nor an exact authorized
  status read is available;
- a path, diff, command, network, credential, or PR boundary is ambiguous;
- a repository command lacks enforceable host isolation for mounts,
  environment, credentials, sockets, descriptors, or child processes;
- required business input is missing;
- production access would be required;
- repository tests require unapproved network or lifecycle scripts;
- lockfile mutation lacks the exact toolchain, registry/snapshot,
  package-integrity inputs, or registry-network grant;
- sandbox execution lacks endpoint, payload/test-data, credential, logging,
  retention, deletion, cleanup, or orphan consent;
- redaction confidence is unknown;
- private runtime data would enter a checkout or public source repository, or
  automatic TTL/revocation deletion, receipt ownership, or orphan recovery
  cannot be enforced;
- held-out inputs or answer material would be visible to an implementation
  worker;
- an implementation worker could write independent lifecycle or PR-lifecycle
  evidence;
- evidence freshness cannot be proved;
- a weak evidence level would be labelled as stronger proof;
- provider-visible status exceeds consumer consent;
- a sponsored connection lacks a valid acknowledged minimal signed receipt;
- rollback or revocation would claim to recall an already exported or publicly
  committed copy;
- remote branch write or PR write is inferred from the other;
- a PR target is the default branch or auto-merge is requested;
- a task depends on product-signal evidence that does not exist;
- M5 or later migration-plan work would proceed before M2.5 closes
  `DISC-001` and `DISC-002`;
- active planning artifacts disagree with the PRD or plan;
- required CI, coverage, CodeQL, review, shipping, post-merge, or item-level
  closure evidence is missing without an approved exception.
