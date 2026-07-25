# Factory Artifacts

Status: v3.1 repo-local controls compiled; factoryd dispatch paused

- `.factory/artifacts/`: committed, source-safe Factory planning, validation,
  closure, and handoff artifacts. It is not a consumer-private runtime store.
- `.factory/artifacts/lifecycle-evidence/`: independent evaluator, grader,
  reviewer, and attestor results; never writable by `task-executor`.
- `.factory/artifacts/pr-lifecycle/`: validation, CI, review, shipping, merge,
  and post-merge evidence.
- `.factory/tmp/`: ignored local scratch.
- `.factoryd/`: ignored daemon state, claims, worktrees, events, and reports.
- `.factory/factoryd.example.json`: mission-paused attended template.
- `.factory/factoryd.autoship.example.json`: mission-paused future full-loop
  template; it authorizes no shipping.

The active compiled v3.1 planning and control set is:

```text
.factory/artifacts/prd-to-plan/lumyn-migration-mvp/
```

The checked-in contents were regenerated from the v3.1 PRD and plan as one
repo-local context, risk, execution, task, validation, acceptance, mapping, and
closure control set. They are planning and validation authority, not product
runtime or factoryd execution authority.

factoryd dispatch remains paused because the external Factory
`profiles/lumyn.yaml` profile and the factoryd bundle/runtime have not been
requalified against this exact generation. A separate reviewed change must
prove that compatibility and explicitly unpause one bounded task.
This pause is scoped to factoryd. A separately approved attended task may use
the same packet and full lifecycle gates without claiming factoryd readiness;
this planning generation itself approves no such task.

The prior plan remains:

```text
.factory/artifacts/prd-to-plan/lumyn-mvp/
```

It and its task, pilot, lifecycle, PR, and exception artifacts are historical
evidence of the superseded product direction. They remain immutable records;
do not reinterpret or rewrite them.

## V3.1 Control Set

The compiled control set represents:

- provider-originated API update delivery launched through a provider-paid,
  services-assisted sunset campaign;
- one reusable provider-confirmed Provider Change Contract and exact
  provider-originated event;
- revocable Consumer Installation and event-specific authorization;
- consumer-local deterministic or bounded-agent patch generation, with agent
  execution disabled unless explicitly configured;
- conditional customer-selected qualified Codex or Claude Code Agent Runner,
  with Cursor deferred behind the same conformance gate;
- exact Agent Runner/model disclosure, endpoint, credential, network, funding,
  usage-billing, clean-session, native-configuration, tool, token, time, retry,
  cost, path, and diff controls;
- exact route-selected product-action capability unions frozen before action;
- managed-credential broker, agent-route topology minimums, pinned
  backend/resource isolation, and separate sandbox-entrypoint contracts,
  including no executable plugins, MCP servers, or hooks for MVP;
- deterministic verification in an independent fresh exact-head process
  without runner/model credentials or a generation-owned evidence writer;
- agent-only repair under a configured, explicitly authorized route;
- patch artifact, optional local branch, and PR-bundle fallback without GitHub
  access;
- one same-run first-campaign proof from authenticated provider event and
  installed preauthorization through an organically agent-assisted item on
  the consumer-selected qualified runner, independent exact-head verification,
  a separately authorized tested Lumyn-opened draft PR, and a consented
  provider-received status projection;
- event-bound, consented provider status with no inference from silence;
- no default-branch write or auto-merge;
- no provider access to consumer code;
- deferred universal event registry, elaborate provider PKI, hosted status
  service, receipt acknowledgement, and receipt-backed billing.

Every new acceptance item remains `planned` until direct evidence closes it.
Retained evidence may carry forward only when its recorded semantics genuinely
prove the new item.

## Authority

Factory workers use only `approval`, `credentials`, and `network`. Exact Lumyn
product grants remain private and action-specific:

- repository read/write;
- command execution;
- Agent Runner network and credential;
- model request disclosure, network, and credential;
- package registry;
- optional sandbox;
- short-lived remote branch and draft PR;
- event-bound, consumer-consented `provider_reporting`;
- retention and deletion.

Empty checked-in grants authorize nothing. The v3.1 rebaseline grants no product
implementation, customer repository, model, credential, command, network,
branch, PR, or merge authority.

## Evidence

Product workers emit task-scoped validation and work proof. They cannot mutate
the canonical active plan, acceptance ledger, mapping, validation contract,
closure map, independent lifecycle evidence, or PR-lifecycle evidence.

Only source-safe control evidence and separately consented aggregate/hash
evidence may be committed. Consumer-private intent, authorization, prompts,
responses, patches, traces, verification, and PR bundles live in an explicitly
configured external private root and are referenced by opaque ID and digest.

Independent `code-review`, `holdout-evaluator`, `trace-grader`, and
`evidence-attestor` gates run when selected by policy. Shipping verifies their
current task/work bindings before `commit-push`; an implementation worker cannot
self-review, self-grade, or self-attest.

Historical evidence proves only its recorded behavior.
