# Lumyn

Lumyn is planning a provider-sponsored, customer-controlled service for
verified API and SDK sunset migrations.

An API Provider pays Lumyn to help move an invited customer cohort off a legacy
surface. Each API Consumer Organization keeps control of its repository,
commands, model egress, credentials, disclosure, review, and merge.

The v3 product path is:

```text
provider-confirmed sunset intent
-> customer-authorized read-only impact
-> reviewable migration plan
-> deterministic transform or consumer-local bounded agent
-> deterministic repository and workflow verification
-> patch + optional local branch + PR bundle
-> optional customer-authorized remote branch and draft PR
-> customer review and merge
```

Provider sponsorship does not grant access to consumer code. Lumyn never
writes the default branch or auto-merges.

The canonical product contract is [docs/product/prd.md](docs/product/prd.md).
The human-readable plan is [docs/product/plan.md](docs/product/plan.md).
[ADR-0003](docs/architecture/adr-0003-services-led-bounded-agent-migration-execution.md)
governs the v3 execution model.
[ADR-0002](docs/architecture/adr-0002-provider-sponsored-customer-controlled-migrations.md)
is retained as historical context for the v2 deterministic-first rebaseline.

The active compiled Factory planning and control set is:

```text
.factory/artifacts/prd-to-plan/lumyn-migration-mvp/
```

It is regenerated from the v3 PRD and plan as one source-aligned set. It
provides repo-local planning and validation authority, not product runtime or
factoryd execution authority. factoryd dispatch remains paused because the
external Factory `profiles/lumyn.yaml` profile and factoryd bundle/runtime have
not yet been requalified against this v3 generation.

The earlier `.factory/artifacts/prd-to-plan/lumyn-mvp/` package is immutable
historical evidence and is not active.

## Product Wedge

The initial offer is a services-led, provider-paid sunset campaign for
participating GitHub repositories that use an official TypeScript/Node SDK.
Lumyn helps the provider normalize and confirm the change intent, works with
consumer maintainers to authorize local execution, and returns evidence-backed
migration handoff artifacts.

A signed declarative provider packet is authoritative when available and
confirmed. It remains data, never executable code. V3 does not make elaborate
provider PKI, continuous provider-status resolution, connection receipts, or
receipt-backed billing prerequisites for the first managed campaigns.
The v2 provider-authenticated receipt flow is historical compatibility context,
not active v3 authority or billing policy.

The Consumer Maintainer reviews a no-write impact report and migration plan
before any change. Each approved plan item routes to one of:

- a deterministic transformation for an exact supported mapping;
- a bounded agent for repository-specific reasoning;
- `needs_input`, `unsupported`, `uncertain`, or `blocked`.

The bounded agent runs in a consumer-local isolated workspace. Before model use,
the consumer approves the exact model provider, endpoint, model/version,
credential environment, source/context disclosure, network allowlist,
logging/training/retention posture, tools, writable paths, file/line/diff
limits, turns, tokens, time, retries, and cost.

Repository and provider content cannot widen those boundaries. Model output is
an untrusted patch candidate and never approves or verifies itself.

Verification is deterministic with respect to the pinned repository head,
commands, fixtures, toolchain, environment, and evidence policy. Generation
mode and verification strength remain separate. Canonical successful labels
are `static_verified`, `repo_verified`,
`workflow_contract_replay_passed`, `workflow_verified_replay`,
`workflow_verified_mock`, and `workflow_verified_sandbox`.

The default handoff requires no GitHub credential:

- patch artifact;
- optional local branch;
- reviewable PR bundle with plan, diff, verification, provenance, and residual
  risk.

Remote branch creation and draft-PR creation are separate optional actions with
separate consumer authorization. Merge always remains human-controlled.

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

- services-led campaign intake and provider-confirmed sunset intent;
- migration corpus and bounded-agent holdouts;
- consumer repository and model execution authorization;
- TypeScript repository impact analysis;
- reviewable migration planning;
- deterministic and bounded-agent patch generation;
- repository and workflow verification runtime;
- patch, local-branch, and PR-bundle handoff;
- optional remote branch and draft-PR delivery;
- private-artifact retention/deletion and operator recovery;
- design-partner campaign measurement.

Only `init` and `check` currently have product behavior. Other command names in
the early dispatcher are compatibility placeholders and must return typed
nonzero results until their implementation acceptance closes.

Public API docs, OpenAPI descriptions, SDK releases, migration guides, and
synthetic fixtures may support planning and engineering. They do not prove
provider demand, customer authorization, or product readiness.

## Two-Party Trust Model

- The API Provider owns and confirms the intended API or SDK change.
- The API Consumer Organization owns repository access, model disclosure,
  credentials, commands, branches, PRs, review, and merge.
- The Lumyn Operator coordinates the paid service but gains no ambient
  repository or credential authority.
- The model provider is a separate egress and data-processing boundary. Model
  access never inherits API-provider, Lumyn-operator, or consumer authority.
- Repository read, local write, command execution, model disclosure, model
  network, model credential, package-registry access, sandbox access, remote
  branch write, and draft-PR write are independent grants.
- Repository tests run without network and through fail-closed host isolation
  by default.
- Raw consumer code, diffs, logs, traces, prompts, responses, and credentials
  are not API-provider-visible by default.
- Consumer-private artifacts remain outside the checkout and public source
  repository with bounded retention and deletion.
- Independent held-out evaluation material is unavailable to the
  implementation worker.
- Lumyn opens a draft PR only when separately authorized and never auto-merges.

The current repository and design-partner distribution are not represented as
open source. Public OSS/self-serve distribution requires a separate approved
license, security, contribution, support, vulnerability-response, and release-
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
- `.factory/artifacts/prd-to-plan/lumyn-migration-mvp/`: v3 compiled-plan target
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
evidence contracts. The active v3 control set is compiled planning authority,
but the checked-in factoryd templates remain paused. This rebaseline grants no
product implementation, consumer repository, model, credential, command,
network, GitHub, or merge authority.

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
