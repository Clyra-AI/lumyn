# Factory Artifacts

Status: v3 repo-local controls compiled; factoryd dispatch paused

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

The active compiled v3 planning and control set is:

```text
.factory/artifacts/prd-to-plan/lumyn-migration-mvp/
```

The checked-in contents were regenerated from the v3 PRD and plan as one
repo-local context, risk, execution, task, validation, acceptance, mapping, and
closure control set. They are planning and validation authority, not product
runtime or factoryd execution authority.

factoryd dispatch remains paused because the external Factory
`profiles/lumyn.yaml` profile and the factoryd bundle/runtime have not been
requalified against this exact generation. A separate reviewed change must
prove that compatibility and explicitly unpause one bounded task.

The prior plan remains:

```text
.factory/artifacts/prd-to-plan/lumyn-mvp/
```

It and its task, pilot, lifecycle, PR, and exception artifacts are historical
evidence of the superseded product direction. They remain immutable records;
do not reinterpret or rewrite them.

## V3 Control Set

The compiled control set represents:

- provider-paid, services-led API or SDK sunset campaigns;
- provider-confirmed intent, with a signed declarative packet as optional
  authoritative input;
- consumer-local deterministic or bounded-agent patch generation;
- exact model disclosure, endpoint, credential, network, tool, token, time,
  retry, cost, path, and diff controls;
- deterministic verification independent from generation mode;
- patch artifact, optional local branch, and PR bundle without GitHub access;
- separately authorized optional remote branch and draft PR;
- no default-branch write or auto-merge;
- no provider access to consumer code;
- deferred elaborate provider PKI, status, receipt, acknowledgement, and
  receipt-backed billing.

Every new acceptance item remains `planned` until direct evidence closes it.
Retained evidence may carry forward only when its recorded semantics genuinely
prove the new item.

## Authority

Factory workers use only `approval`, `credentials`, and `network`. Exact Lumyn
product grants remain private and action-specific:

- repository read/write;
- command execution;
- model request disclosure, network, and credential;
- package registry;
- optional sandbox;
- optional remote branch and draft PR;
- optional `provider_reporting`;
- retention and deletion.

Empty checked-in grants authorize nothing. The v3 rebaseline grants no product
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
