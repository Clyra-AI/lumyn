# Lumyn Factory Integration

Source-safe Factory control and lifecycle artifacts for Lumyn live under
`.factory/artifacts/`. Consumer-private runtime artifacts and identifiable
pilot evidence do not.

- Active migration plan:
  `.factory/artifacts/prd-to-plan/lumyn-migration-mvp/`
- Historical agent-readiness plan:
  `.factory/artifacts/prd-to-plan/lumyn-mvp/`
- Task-run evidence: `.factory/artifacts/task-runs/`
- Independent lifecycle evidence:
  `.factory/artifacts/lifecycle-evidence/` (never writable by
  `task-executor`)
- PR lifecycle evidence: `.factory/artifacts/pr-lifecycle/`
- Task-supervisor intake evidence:
  `.factory/artifacts/task-supervisor-runs/`
- Autoship supervisor evidence: `.factory/artifacts/supervisor-runs/`
- Historical downstream pilot:
  `.factory/artifacts/pilot/lumyn-mvp-slice/`
- Public migration-pilot summaries:
  `.factory/artifacts/pilot/lumyn-migration-mvp/public/` (consented, redacted
  aggregates and evidence hashes only)
- Local attended daemon config: `.factory/factoryd.json` (gitignored; copy
  from `.factory/factoryd.example.json`)
- Safe config template: `.factory/factoryd.example.json`
- Explicit autoship template: `.factory/factoryd.autoship.example.json`
- Local daemon state: `.factoryd/` (gitignored)

The authored product truth is:

```text
docs/product/prd.md
docs/product/plan.md
```

The prior `lumyn-mvp` task/control artifacts, task runs, PR lifecycle records,
exceptions, and pilot evidence remain immutable historical records; its README
contains only a non-operative dispatch tombstone. They may establish the
existence of the Go CLI, schema/evidence foundation, source intake, and delivery
controls. They must not be cited as proof of migration impact analysis,
patching, verification execution, GitHub delivery, or product demand.

## Operator Flow

Use Factory to change shared contracts and the Lumyn profile. Use `factoryd` to
execute the active Lumyn task packets. Use Lumyn commits and PRs to change
product code, docs, CI, and product evidence.

Start with non-mutating proof:

```text
export FACTORY_REPO=../factory
cp .factory/factoryd.example.json .factory/factoryd.json
factoryd doctor --config .factory/factoryd.json --repo lumyn --json
factoryd run --config .factory/factoryd.json --repo lumyn --dry-run --json
```

Run one implementation task without remote shipping only after its packet has
resolved allowed and forbidden paths, validation commands, evidence,
dependencies, customer authority, lifecycle gates, and stop conditions:

```text
export FACTORY_REPO=../factory
factoryd run --config .factory/factoryd.json --repo lumyn --once --json
```

Use autoship only after branch protection, required `validate` and
`CodeQL analyze` checks, passive Codex review settle, merge policy, post-merge
monitoring, and semantic scope-closure evidence are proven:

```text
export FACTORY_REPO=../factory
factoryd run --config .factory/factoryd.autoship.example.json --repo lumyn --loop --max-tasks 1 --json
```

The one-task loop is intentional. It keeps PRs small and produces
task-scoped work proof, validation, PR lifecycle, post-merge, scope-closure, and
repair evidence.

## Product Authority Is Not Factory Authority

Factory shipping authority governs changes to the Lumyn repository. It does
not grant Lumyn product authority over a consumer repository.

Migration tasks use distinct product capability grants:

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

Read-only impact may proceed only with an approved repository-read grant.
Migration planning never implies write permission. Patch application, command
execution through exact host isolation, provider-status read, package-registry
access, sandbox network access, sandbox payload disclosure, remote branch
creation/update, PR creation/update, minimal connection receipt, and richer
provider-visible attestation each require their own approved scope. A registry
or status grant does not authorize arbitrary network access; a branch-write
grant does not authorize PR write; a PR-write grant does not authorize branch
or default-branch mutation. No product grant authorizes auto-merge.

The checked-in daemon configs are templates. Empty grants prove only that
deterministic offline work can be selected safely. Operators record live,
task-scoped Factory worker grants in the gitignored `.factory/factoryd.json`;
they do not edit generated task packets to bypass a gate. Exact Lumyn product
grants live in the configured private runtime root as schema-backed
authorization artifacts. A Factory `approval` grant cites their opaque
validated bundle reference; it never embeds private grant details or replaces
product authority.

Before dispatch, `factoryd doctor`, dry-run selection, and repo-pack validation
must reject product capability names placed directly in Factory config,
wildcard or duplicate Factory grants, missing expiry/evidence, and a Factory
grant whose task or capability does not match the selected packet. Factory
selection checks that the task declares the exact Lumyn capabilities and live
gate contract, but it does not read or validate the private product bundle.
Inside Lumyn, immediately before every local, sandbox, GitHub, or attestation
side effect and retry, the product must revalidate the exact applicable signed
bundle and current packet trust from a signed offline status snapshot or exact
`provider_trust_status_read` grant. `campaign accept` remains a no-authority
request; only explicit consumer-signed `authorization issue`/`revoke` changes
private authority state. `lumyn authorization validate` supplies task
validation and closure proof; it is not a Factory pre-dispatch product gate.
Negative product tests must prove that M8, M9, and M10 perform no live action
when any required product grant or current-status proof is missing, expired,
revoked, stale, or inconsistent.

M8 and M9 are independent successors to M7. M9 requires exact remote-branch and
draft-PR authority, but it does not require an M8 sandbox run or
`provider_attestation`. A sponsored repository counts only after a
consumer-signed minimal `campaign_receipt` uses a provider-authenticated
consumer signer binding and receives a provider-signed, one-unit, deduplicated
acknowledgement; richer provider reporting is an optional, separately
authorized action.

Codex CLI authentication used by a Factory worker is not a Lumyn product
credential. Customer and provider credentials must remain task-scoped,
least-privilege, isolated by capability, and absent from committed artifacts.

## Active Control Truth

The following files are one generated control set and must remain aligned:

```text
.factory/artifacts/prd-to-plan/lumyn-migration-mvp/context-brief.json
.factory/artifacts/prd-to-plan/lumyn-migration-mvp/risk-classification.json
.factory/artifacts/prd-to-plan/lumyn-migration-mvp/execution-plan.json
.factory/artifacts/prd-to-plan/lumyn-migration-mvp/task-packets.json
.factory/artifacts/prd-to-plan/lumyn-migration-mvp/validation-contract.json
.factory/artifacts/prd-to-plan/lumyn-migration-mvp/acceptance-ledger.json
.factory/artifacts/prd-to-plan/lumyn-migration-mvp/acceptance-mapping.json
.factory/artifacts/prd-to-plan/lumyn-migration-mvp/scope-closure-map.json
```

The active set contains exactly 62 acceptance items across 12 task packets
(`M0` through `M10`, including `M2.5`). `M2.5` is the preimplementation
provider/supported-canary and five-repository/three-organization cohort,
operational provider-status and receipt-acknowledgement channels, supported
consumer OS/architecture and host-isolation backend, confidentiality, consent,
retention, economics, and measurement gate;
`M10` is the governed real-campaign task. Task completion and delivery slices
are coverage lenses, never substitutes for item-level closure.

M2 establishes the standard provider campaign kit, configured packet
publication, signed invitation, no-authority `campaign accept`, explicit
consumer-signed authorization issue/revoke, provider-status freshness, host
isolation, and provider-authenticated/signed minimal connection-receipt
acknowledgement flow. M7 must then pass the synthetic
receipt-backed offline canary through a local draft-PR preview before any live
sandbox or GitHub work.

Product workers must not edit these planning artifacts directly. They emit
task-scoped evidence; trusted lifecycle workers update closure. If the PRD,
runtime pins, authority model, task graph, validation requirements, acceptance
items, or Factory profile changes, regenerate and validate the full set before
dispatch.

Validate guide propagation and item-level closure with:

```text
python3 scripts/validate_repo_pack.py --self-test
python3 scripts/validate_repo_pack.py
python3 scripts/validate_factory_pilot_evidence.py
```

The pilot validator covers only the frozen historical pilot package. A new
migration campaign receives a new evidence namespace and may not inherit the
historical pilot's product claims.

M2.5 and M10 each require a task-specific validator invocation against the
actual private attestation and its aggregate/hash-only public manifest.
Validator `--self-test` output is engineering proof, not external outcome
proof. The task packet also requires an independent `evidence-attestor` record
before any `DISC` or `PILOT` item can close.

Before M2.5 collects, stores, or discloses identifiable external evidence, its
Factory `approval` evidence must cite a narrow manual privacy/legal preflight
that fixes allowed private fields, participant consent, the approved external
private root, TTL, expiry/revocation deletion, deletion-receipt and orphan
ownership, the minimal connection-receipt allowlist, authenticated consumer
signer binding, provider-signed acknowledgement/cardinality policy, separately
consented aggregate/hash-only public fields, and the irreversibility of
provider/public disclosure. This is an implementation-worker evidence-handling
approval, not Lumyn runtime product authority. The active grant must also
exactly match the preflight's canonical scope digest; a generic approval or
stale digest cannot release collection.

Tasks selected for independent evaluation place `holdout-evaluator`,
`trace-grader`, and `evidence-attestor`, in that order, after `code-review` and
before `commit-push`. These are external or human-operated gates. Shipping
must verify their lifecycle-owned, task-bound, current-work-proof, passing
artifacts from the trusted Factory evidence root; the implementation worker
cannot synthesize its own approval.

Only a non-resolving opaque holdout manifest is committed: opaque case IDs,
provenance class and license posture, a frozen suite commitment, and encrypted
or HMAC artifact commitments. An independent holdout owner provisions and
freezes `LUMYN_HOLDOUT_ROOT`; resolving provenance, plaintext content digests,
held-out inputs, answer keys, expected patches or labels, and raw traces stay
there and are unavailable to implementation workers.
M1's provision-mode policy creates the suite commitment at independent
evaluation time; M4/M6/M7 evaluate-mode policies resolve and byte-bind that
trusted result. The static plan never fabricates a future commitment.
Lifecycle artifacts must schema-validate and bind the exact task, work item,
lifecycle run, current validation run, candidate digest, work-proof marker
digest, passing result, and independent worker provenance. Task-level review
lens and reviewer class must exactly match inherited review requirements.

## Runtime State And Supervision

`factoryd` operational state lives under `.factoryd/`. Committed, source-safe
Factory planning and delivery evidence lives under `.factory/artifacts/`.
Consumer-private Lumyn runtime artifacts and identifiable pilot evidence live
in an explicitly configured consumer-controlled state root outside the checkout
and any public source repository. The private store enforces authorization TTL,
deletion on expiry or revocation at creation, read, process startup, and the
next run, and a deletion receipt or orphan report. `lumyn artifacts gc` is the
operator recovery path for retrying partial deletion and inspecting unresolved
orphans; it cannot extend TTL or rewrite historical closure.
Factory artifacts may reference private evidence by opaque identifier and
digest only.

Provider-visible and public are separate disclosure decisions. A
provider-attestation grant authorizes only its field allowlist; publishing even
that redacted attestation requires separate explicit public consent. The
committed migration-pilot namespace may contain only consented aggregates and
evidence hashes, never repository identity, raw source, diffs, logs, traces,
prompts, responses, credentials, or identifiable participant records.
Provider export and public commit are irreversible: revocation blocks future
sharing and deletes Lumyn-controlled private copies, but cannot erase recipient
records, Git history, clones, or caches.

The active daemon config carries `runtime_control` and must fail closed for
pause, cancel, freeze, disabled adapter, read-only, write-scope, credential, or
network conflicts. Claimed runs emit a canonical mission event log and rerun
repo validation after task changes exist in the worker worktree.

For attended daemon-led delivery, invoke Factory `autoship-supervisor` from the
Lumyn repository root and name exactly one task. `factoryd` owns implementation,
validation, PR creation, CI/review polling, merge, post-merge monitoring, and
routine repair. The supervisor intervenes only for typed blockers, explicit
human acceptance, non-convergent repair, stale plan scope, or systemic
Factory/factoryd follow-up, and records:

```text
.factory/artifacts/supervisor-runs/<task_id>/<timestamp>.json
```

## Post-PRD Findings

Material `app-audit`, `repo-audit`, and `code-review` findings become governed
work only after they are saved as structured repo-local finding-list JSON and
ingested. Use Factory `task-supervisor` for guided intake:

```text
FACTORY_REPO=/path/to/factory factoryd ingest --config .factory/factoryd.example.json --repo lumyn --kind audit --input product/audits/<mission>.finding-list.json --mission <mission> --json
FACTORY_REPO=/path/to/factory factoryd ingest --config .factory/factoryd.example.json --repo lumyn --kind review --input product/reviews/<mission>.finding-list.json --mission <mission> --json
```

Generated `.factory/artifacts/post-prd/<mission>/` artifacts become that
mission's execution contract. Do not mutate `docs/product/prd.md` unless a human
explicitly promotes the finding into product scope.
