# Lumyn Factory Integration

Status: v3.1 repo-local controls compiled; factoryd dispatch paused

Source-safe Factory control and lifecycle artifacts live under
`.factory/artifacts/`. Consumer-private runtime artifacts and identifiable
campaign evidence do not.

- Active compiled v3.1 control set:
  `.factory/artifacts/prd-to-plan/lumyn-migration-mvp/`
- Historical agent-readiness plan:
  `.factory/artifacts/prd-to-plan/lumyn-mvp/`
- Task-run evidence: `.factory/artifacts/task-runs/`
- Independent lifecycle evidence:
  `.factory/artifacts/lifecycle-evidence/`
- PR lifecycle evidence: `.factory/artifacts/pr-lifecycle/`
- Post-PRD findings: `.factory/artifacts/post-prd/`
- Daemon scratch/state: `.factory/tmp/` and `.factoryd/`

The checked-in `lumyn-migration-mvp` directory contains the control set
regenerated from the v3.1 PRD and plan. It is one repo-local planning, task,
validation, acceptance, and closure authority, not product runtime or factoryd
execution authority.

A v3.1 source update has landed in the external Factory
`profiles/lumyn.yaml` profile, but the active compiled control generation still
records Factory profile/runtime compatibility as unqualified. factoryd
dispatch remains paused until the complete controls are regenerated and its
bundle/runtime and exact active mission are qualified. A source-profile update
alone authorizes no product implementation or live product action.
A compiled task packet's `factory_compatibility` block records the external
posture observed when that generation was produced and remains authoritative
for factoryd. Attended dispatch evidence may verify current external
dependencies separately; it must not hand-edit a compiled snapshot or
reinterpret the source update as factoryd qualification.
A separately approved attended task may use the same packet and full lifecycle
gates without claiming factoryd readiness. The checked-in pause therefore
blocks factoryd dispatch, not every possible human-run implementation path.

The former `lumyn-mvp` package and its task, pilot, lifecycle, PR, and exception
artifacts are immutable historical records. They must not be rewritten to imply
v3.1 behavior.

## V3.1 Compiled Control Boundary

The active compiled generation covers:

- provider-originated API update delivery launched through a provider-paid,
  services-assisted sunset campaign;
- one reusable provider-confirmed Provider Change Contract and exact event;
- revocable Consumer Installation and event-specific authorization;
- consumer-local deterministic or bounded-agent generation, with agent
  execution disabled until explicitly configured;
- conditional customer-selected qualified Codex or Claude Code Agent Runner,
  with exact executable/auth/entitlement, funding, clean-session, and
  no-fallback controls;
- exact Agent Runner/model disclosure, endpoint, credential, network, tool,
  provenance, and budget controls;
- exact route-selected product-action capability unions frozen before action;
- managed-credential broker, agent-route topology minimums, pinned
  backend/resource isolation, and separate sandbox-entrypoint contracts, with
  executable plugins, MCP servers, and hooks prohibited for MVP;
- deterministic repository and workflow verification in a fresh exact-head
  process without runner/model credentials or a generation-owned evidence
  writer;
- agent-only repair under a configured, explicitly authorized route;
- patch artifact, optional local branch, and PR bundle as fallback;
- one same-run first-campaign proof from authenticated provider event and
  installed preauthorization through an organically agent-assisted item on
  the consumer-selected qualified runner, independent exact-head verification,
  a short-lived tested Lumyn-opened draft PR, and an event-bound, consented
  provider-received status projection;
- consumer review and merge;
- no provider access to code, default-branch write, or auto-merge.

Universal event distribution, elaborate provider PKI, hosted status,
connection receipts, acknowledgements, and receipt-backed billing are deferred
concepts, not active v3.1 prerequisites.

This compiled rebaseline is planning and control work only. M2 separately
implements executable artifact and semantic-policy contracts, but neither the
compiled controls nor those contracts authorize a model endpoint, consumer
repository, credential, command, network, sandbox, GitHub, or merge action.

## Operator Flow

Use the checked-in `make` and validator commands for repo-local structural and
contract validation. Do not invoke factoryd to select or execute a Lumyn task
while the mission is paused.

Before a future factoryd dispatch, a separate reviewed change must prove:

- the external Factory `profiles/lumyn.yaml` profile still describes the exact
  v3.1 runtime, validation, evidence, and capability posture;
- the exact factoryd binary, bundle, schemas, semantic validators, and
  active-mission resolution support this compiled generation;
- the repo-local PRD, plan, ADR, operating docs, validators, examples, and
  compiled controls still agree;
- the selected task has bounded paths, commands, tools, credentials, network,
  evidence, and positive budgets; and
- the paused runtime control is explicitly changed for that task.

The checked-in configs contain no grants and remain paused. A future operator
may put approved, task-scoped Factory grants only in gitignored
`.factory/factoryd.json`; those grants still do not confer Lumyn product
authority.

For `conditional_factory_capabilities`, the runtime grant set must bind one
frozen task/action mode, the exact sorted selected capabilities, one evidence
ref and scope digest, and a common expiry. It must reject partial, extra,
stale, wildcard, or cross-task activation.

## Product Authority Is Not Factory Authority

Factory uses the closed worker capability vocabulary:

- `approval`
- `credentials`
- `network`

These govern the Factory implementation worker. They do not grant Lumyn product
authority.

V3 product grants remain private, schema-backed, expiring, and revocable:

- `customer_repo_read`
- `customer_repo_write`
- `command_execution`
- `model_request_disclosure`
- `agent_runner_network`
- `agent_runner_credential`
- `model_network`
- `model_credential`
- `package_registry_read`
- `sandbox_request_disclosure`
- `sandbox_network`
- `sandbox_credential`
- `github_branch_write`
- `github_pr_write`
- `provider_reporting`
- `artifact_retention`
- `artifact_deletion`

Important separations:

- provider payment is not repository consent;
- Lumyn Operator service work is not consumer authority;
- plan approval is not local-write approval;
- model disclosure is not provider disclosure;
- Agent Runner Vendor is not necessarily the Model Provider;
- Agent Runner network/credential is not model network/credential;
- provider-sponsored usage is not API Provider agent or repository authority;
- model disclosure, network, and credential are independent;
- patch and PR-bundle creation require no GitHub credential;
- local branch is not remote branch;
- remote branch write is not draft-PR write;
- draft-PR write is not merge authority.

Factory dispatch cannot validate or confer private Lumyn product authority.
Product code must enforce the exact applicable bundle at each side-effect
boundary.

## Bounded-Agent Factory Contract

Any task that implements or evaluates agent-assisted generation records:

- explicit `agent_execution_policy`;
- selected model provider, endpoint, model/version, and parameters;
- prompt/system/tool definition refs or digests;
- context disclosure and redaction;
- credential environment and scopes;
- exact network endpoint;
- read/write paths and tool allowlist;
- file, line, diff, turn, token, time, retry, concurrency, and cost budgets;
- isolated workspace and cleanup;
- exact managed-credential broker and runner-host isolation evidence when
  applicable;
- request, response, tool-call, usage, attempt, and patch provenance;
- prompt-injection and scope-widening negative tests;
- deterministic verification from a fresh exact-head view in an independent
  process without runner/model credentials;
- independent holdout/review and human approval requirements.

The implementation worker cannot access evaluator-controlled answer material or
write independent evidence. Model output cannot satisfy its own acceptance or
verification.

## Active Control Truth

The following files form one generated control set:

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

They were regenerated together for the v3 rebaseline. Regenerate the full set
again whenever an authoritative input changes. Do not hand-edit one artifact
to bypass a task, capability, acceptance, validation, or closure gate.

The active set must derive its acceptance count and task mapping from the v3
PRD and plan. Do not preserve the v2 count merely to reduce diff size.

Carried evidence remains valid only for unchanged semantics. Every new
agent/model/services item begins planned until direct evidence closes it.

## Validation And Lifecycle

Required local lanes:

```bash
make lint-fast
make test-fast
make test-coverage
make test-contracts
make prepush-full
```

The normal chain is:

1. `task-executor`
2. `validation-gate`
3. `code-review` when risk requires it
4. `holdout-evaluator` when selected
5. `trace-grader` when selected
6. `evidence-attestor` when selected
7. `commit-push`
8. `post-merge-monitor`
9. `repair-feedback` on failure

Required independent evidence is lifecycle-owned, task-bound, current, and
passing before `commit-push`. Product workers cannot write the lifecycle or
PR-lifecycle namespaces.

Passive latest-head Codex review settle, `validate`, `CodeQL analyze`, branch
protection, conversation resolution, shipping evidence, and post-merge
monitoring remain merge gates. Green CI alone is not merge-ready.

## Evidence Boundaries

Committed Factory artifacts contain only source-safe control, validation,
lifecycle, or separately consented aggregate/hash evidence.

Consumer-private provider intent, repository authorization, prompts, responses,
tool traces, patches, verification, and PR bundles live in an explicitly
configured root outside the checkout and public source repository. Committed
records use opaque IDs and digests.

Provider disclosure and model-provider disclosure are separate. Raw source,
diffs, prompts, responses, agent sessions, logs, traces, credentials, and
private evidence are never API-provider-visible. Only enumerated,
consumer-consented campaign status or aggregate fields may cross that
boundary.

## Runtime State And Supervision

`.factory/factoryd.json` is gitignored attended configuration.
`.factory/factoryd.example.json` and
`.factory/factoryd.autoship.example.json` are paused templates, not executable
authority. The autoship shape is retained only for future requalification.

Before any future live task, validate:

- exact task and current plan generation;
- allowed/forbidden paths;
- commands and host isolation;
- Agent Runner executable/auth/entitlement/funding, network/credential, and
  model disclosure/network/credential and budgets when applicable;
- optional sandbox grants;
- separately approved, short-lived remote-branch and draft-PR grants;
- expiry, revocation, retention, deletion, and evidence refs.

Empty grants prove only that offline planning and fixture validation can remain
safe. They do not prove product execution or factoryd readiness.

## Post-PRD Findings

Save material audits, reviews, pilot findings, and recommendations as governed
post-PRD artifacts. A human explicitly promotes a finding into product scope,
then regenerates every affected authored and compiled surface.

Do not mutate the PRD or canonical control set from a task-run worker.
