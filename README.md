# Lumyn

Lumyn is planning the provider-to-consumer application layer for important API
and SDK updates.

An API Provider publishes one confirmed change and pays Lumyn to carry it into
an invited customer cohort. Each API Consumer Organization installs the exact
channel and actions it accepts and keeps control of its repository, commands,
model egress, credentials, disclosure, review, and merge.

The v3 product path is:

```text
provider-originated change event
-> reusable Provider Change Contract
-> customer-installed channel and action policy
-> customer-authorized read-only repository impact
-> reviewable migration plan
-> deterministic transform or customer-selected consumer-local bounded agent
-> deterministic repository and workflow verification
-> tested customer-authorized draft PR
-> customer review and merge
-> consented provider rollout status
```

Provider sponsorship does not grant access to consumer code. Lumyn never
writes the default branch or auto-merges.

The canonical product contract is [docs/product/prd.md](docs/product/prd.md).
The human-readable plan is [docs/product/plan.md](docs/product/plan.md).
[ADR-0004](docs/architecture/adr-0004-provider-originated-api-update-delivery.md)
governs the v3.1 provider-to-consumer product direction.
[ADR-0005](docs/architecture/adr-0005-customer-selected-agent-runners.md)
governs Agent Runner selection, qualification, credentials, funding, and
verification.
[ADR-0003](docs/architecture/adr-0003-services-led-bounded-agent-migration-execution.md)
remains the bounded-agent execution and trust decision.
[ADR-0002](docs/architecture/adr-0002-provider-sponsored-customer-controlled-migrations.md)
is retained as historical context for the v2 deterministic-first rebaseline.

The active compiled Factory planning and control set is:

```text
.factory/artifacts/prd-to-plan/lumyn-migration-mvp/
```

It is regenerated from the v3.1 PRD and plan as one source-aligned set. It
provides repo-local planning and validation authority, not product runtime or
factoryd execution authority. factoryd dispatch remains paused because the
external Factory `profiles/lumyn.yaml` profile and factoryd bundle/runtime have
not yet been requalified against this v3.1 generation.

The earlier `.factory/artifacts/prd-to-plan/lumyn-mvp/` package is immutable
historical evidence and is not active.

## Product Wedge

The product wedge is a per-provider update channel for an official
TypeScript/Node SDK. The initial offer is a services-assisted, provider-paid
sunset campaign: Lumyn helps the provider normalize and confirm one Provider
Change Contract and event, helps consumer maintainers install bounded local
policy, and opens evidence-backed tested draft PRs.

A signed event manifest at a pinned provider-controlled HTTPS URL is the first
channel. It embeds the Provider Change Contract or pins the exact
provider-controlled URL whose retrieved bytes must match its digest. The
provider publishes it; Lumyn may assist setup. It remains data, never
executable code. An attended import is recovery, not channel or
unattended-delivery proof. V3 does not make elaborate provider PKI, continuous
provider-status resolution, connection receipts, or receipt-backed billing
prerequisites for the first managed campaigns.
The v2 provider-authenticated receipt flow is historical compatibility context,
not active v3 authority or billing policy.

Impact and planning are always no-write. The Consumer Maintainer either
approves the exact event plan or installs narrowly bounded preauthorization;
any plan outside the installed policy pauses before mutation. Each authorized
plan item routes to one of:

- a deterministic transformation for an exact supported mapping;
- a bounded agent for repository-specific reasoning;
- `needs_input`, `unsupported`, `uncertain`, or `blocked`.

Agent execution is disabled unless the consumer configures it; notify-only,
scan-only, and deterministic-only installations need no agent account or
credential. When a routed plan needs agent assistance, the bounded agent runs
through a customer-selected Agent Runner in a consumer-local isolated
workspace. Launch adapters are Codex and Claude Code once each passes the
same conformance suite and live canary; Cursor remains deferred behind that
gate. Before use, the consumer approves the exact adapter and version,
executable source/digest, auth mode and entitlement class, Agent Runner Vendor,
actual Model Provider and model route, credential and usage-billing owner,
source/context disclosure, network allowlists, logging/training/retention
posture, tools, native agent configuration, writable paths, file/line/diff
limits, turns, tokens, time, retries, and cost. Lumyn uses a neutral home/config
root, rejects repository-local executable shadowing, and never resumes a
personal agent session or silently falls back to another adapter or model.

For configured agent execution, the default `consumer_managed` mode uses the
consumer's own qualifying agent account, enterprise subscription, API
credential, or local runtime. A route qualifies only when its actual model
identity is observable and non-interactive automation is permitted. An
optional `provider_sponsored_lumyn_managed` route lets the API Provider fund
the campaign and Lumyn pay approved agent/model usage while execution and
consent remain in the consumer boundary. The API Provider receives no
repository, prompt, session, or credential access in either mode.

Managed credentials use one-time broker redemption into one attempt-scoped
session. That session may make multiple calls only within hard token/cost
quotas and cannot refresh, replay after the attempt, or carry into another
attempt. It is revocable and reconciled through a vendor-native credential or
approved budget-enforcing proxy; otherwise that route is unavailable. Agent
processes receive only explicit mounts and no host home, OS credentials,
ambient sockets, or unrelated descriptors. Child-process, egress, and cleanup
controls are host-enforced under a pinned qualified backend with hard CPU,
memory, PID/process-tree, disk, and open-file quotas. Each agent action freezes
one local, runner-mediated, direct-model, or hybrid authorization topology so
its minimum network, credential, and disclosure scopes cannot be omitted.
Executable plugins, MCP servers, and hooks are outside MVP.

Repository and provider content cannot widen those boundaries. Model output is
an untrusted patch candidate and never approves or verifies itself.

Verification is deterministic with respect to the pinned repository head,
commands, fixtures, toolchain, environment, and evidence policy. It runs from
a fresh exact-head view in a separate process without Agent Runner or model
credentials and writes evidence through a boundary unavailable to generation.
Generation mode and verification strength remain separate. `lumyn repair` is
an agent-assisted action only and requires a configured, explicitly authorized
agent route. Canonical successful labels are `static_verified`, `repo_verified`,
`workflow_contract_replay_passed`, `workflow_verified_replay`,
`workflow_verified_mock`, and `workflow_verified_sandbox`.

Every run retains a no-GitHub fallback:

- patch artifact;
- optional local branch;
- reviewable PR bundle with plan, diff, verification, provenance, and residual
  risk.

Remote branch creation and draft-PR creation are separate consumer-authorized
actions. At least one Lumyn-opened draft PR is required to prove the first
provider campaign; manual fallback does not prove automated delivery. The
qualifying commercial proof is one same run from authenticated provider event
and installed preauthorization through an organically agent-assisted item on
the consumer-selected qualified runner, independent exact-head verification,
the Lumyn-opened draft PR, and a consented provider-received status
projection. Separate agent, delivery, or reporting runs do not qualify. Merge
always remains human-controlled.

## Current Implementation Status

Implemented:

- Go CLI/config/result/exit-code foundation;
- `lumyn init`;
- `lumyn check`;
- OpenAPI and local-document source parsing, refs, fingerprints, and findings;
- executable workflow, evidence, cassette, trace, proof, boundary, redaction,
  and command-result schemas;
- local validation, coverage, CodeQL, branch-policy, review, and Factory
  delivery controls.

Planning-only, not implemented:

- provider event, Provider Change Contract, and services-assisted campaign intake;
- consumer installation and event-specific authorization;
- migration corpus and bounded-agent holdouts;
- consumer repository, Agent Runner, and model execution authorization;
- TypeScript repository impact analysis;
- reviewable migration planning;
- deterministic and bounded-agent patch generation;
- repository and workflow verification runtime;
- patch, local-branch, and PR-bundle handoff;
- short-lived remote branch and draft-PR delivery;
- event-bound consented provider status projection;
- private-artifact retention/deletion and operator recovery;
- design-partner campaign measurement.

Only `init` and `check` currently have product behavior. Other command names in
the early dispatcher are compatibility placeholders and must return typed
nonzero results until their implementation acceptance closes.

Public API docs, OpenAPI descriptions, SDK releases, migration guides, and
synthetic fixtures may support planning and engineering. They do not prove
provider demand, customer authorization, or product readiness.

## Two Authorities And External Processors

- The API Provider owns and confirms the intended API or SDK change.
- The API Consumer Organization owns repository access, model disclosure,
  credentials, commands, branches, PRs, review, and merge.
- The Lumyn Operator coordinates the paid service but gains no ambient
  repository or credential authority.
- The Agent Runner Vendor and Model Provider are distinct egress and
  data-processing roles even when one vendor supplies both. Their access never
  inherits API-provider, Lumyn-operator, or consumer authority.
- Repository read, local write, command execution, model disclosure, model
  network, model credential, Agent Runner network, Agent Runner credential,
  package-registry access, sandbox access, remote branch write, and draft-PR
  write are independent grants.
- Repository tests run without network and through fail-closed host isolation
  by default.
- Raw consumer code, diffs, logs, traces, prompts, responses, agent sessions,
  and credentials are never API-provider-visible. Only enumerated, consented
  campaign status or aggregate fields may cross that boundary.
- Consumer-private artifacts remain outside the checkout and public source
  repository with bounded retention and deletion.
- Independent held-out evaluation material is unavailable to the
  implementation worker.
- Lumyn opens a draft PR only when separately authorized and never auto-merges.

The current repository and design-partner distribution are not represented as
open source. Public OSS/self-serve distribution requires a separate approved
license, security, contribution, support, vulnerability-response, and release
integrity gate.

## Repository Layout

- `cmd/lumyn/`: CLI entrypoint
- `internal/`: Go implementation packages
- `tests/`: automated tests
- `schemas/`: versioned executable artifact schemas
- `examples/`: source-safe examples and migration fixtures
- `workflows/`, `cassettes/`, and `runs/`: retained synthetic or licensed
  fixtures only
- `docs/product/prd.md`: product source of truth
- `docs/product/plan.md`: human-readable plan
- `docs/dev/dev_guides.md`: engineering and validation contract
- `docs/architecture/architecture_guides.md`: architecture and trust boundaries
- `.factory/artifacts/prd-to-plan/lumyn-migration-mvp/`: v3.1 compiled-plan target
- `.factory/artifacts/prd-to-plan/lumyn-mvp/`: immutable historical plan
- `.factory/artifacts/lifecycle-evidence/`: independent lifecycle evidence
- `.factory/artifacts/pr-lifecycle/`: Factory PR lifecycle evidence

## Validation

```bash
make lint-fast
make test-fast
make test-coverage
make test-contracts
make prepush-full
```

`make prepush-full` is the required local gate before PR and merge. GitHub
Actions runs the same gate through `validate` and runs CodeQL through
`CodeQL analyze`.

## Factory Operation

Factory supplies planning, task-packet, validation, review, shipping, and
evidence contracts. The active v3.1 control set is compiled planning authority,
but the checked-in factoryd templates remain paused. This rebaseline grants no
product implementation, consumer repository, model, credential, command,
network, GitHub, or merge authority.
Separately approved attended tasks may use the same packets and lifecycle
gates without claiming factoryd readiness.

Factory's `approval`, `credentials`, and `network` grants govern its
implementation worker. They do not validate or confer Lumyn product authority.

Do not use factoryd to select or execute Lumyn work until a separate reviewed
change rebaselines the external Factory profile, proves exact active-mission
and bundle/runtime compatibility, and explicitly unpauses a bounded task with
positive budgets. Structural repository validation may continue through the
checked-in `make` and validator commands.

Any later autoship path remains gated by branch protection, required checks,
passive Codex review, `commit-push`, post-merge monitoring, item-level
acceptance evidence, and every exact product authorization required by the
task.
