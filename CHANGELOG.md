# Changelog

## [Unreleased]

### Added

- Added ADR-0003 and a v3 product/engineering planning contract for
  services-led, provider-paid API and SDK sunset campaigns.
- Added hybrid patch routing: deterministic transforms for exact mappings and
  a consumer-local bounded coding agent for approved repository-specific work.
- Added exact model request-disclosure, endpoint, credential, tool, path,
  token, time, retry, cost, provenance, retention, and deletion boundaries.
- Added patch artifact and PR bundle as the no-GitHub baseline handoff, with
  optional local branch, remote branch, and draft PR as separate delivery
  states.
- Added a source-aligned v3 compiled Factory control set covering context,
  risk, execution, task packets, validation, acceptance, mapping, and closure.
- Added explicit factoryd mission pause and compatibility gates for the
  external Factory profile, factoryd bundle/runtime, and exact active mission.
- Retained the Go CLI/config/result foundation, `lumyn init`, `lumyn check`,
  local source parsing, executable schemas, validation, coverage, CodeQL,
  branch protection, and Factory lifecycle controls.

### Changed

- Reframed Lumyn from generic agent-readiness evaluation and the v2
  deterministic/receipt-first migration plan to verified API migration
  execution.
- Made the API Provider the initial campaign buyer while preserving API
  Consumer Organization authority over repository access, execution, model
  egress, credentials, disclosure, review, and merge.
- Made services-led local or consumer-CI execution the initial form factor;
  hosted SaaS is not required for the first campaign.
- Kept provider-confirmed migration intent authoritative and signed
  declarative packets supported, while deferring mandatory packet PKI,
  continuous provider status, connection receipts, acknowledgements, and
  receipt-backed billing.
- Made generation provenance independent from deterministic, exact-candidate
  verification strength.
- Replaced checked-in v2 active-control claims under
  `lumyn-migration-mvp` with the regenerated v3 compiled control set.
- Preserved `.factory/artifacts/prd-to-plan/lumyn-mvp/`, ADR-0002, and their
  lifecycle evidence as immutable historical records.
- Clarified that the repo-local v3 compilation is planning and validation
  authority only. It authorizes no product runtime implementation or live
  action.
- Kept factoryd dispatch paused until the external Factory
  `profiles/lumyn.yaml` profile and factoryd bundle/runtime are separately
  requalified and a bounded task is explicitly unpaused.
- Replaced the premature OSS claim with explicitly licensed,
  integrity-signed design-partner distribution and a separate gate for any
  future public OSS, self-serve, or Homebrew release.

### Deprecated

- Deprecated the v2 deterministic-only execution boundary and
  PKI/receipt-first activation and billing path as active v3 requirements.
- Deprecated generic live-agent evaluation, model-provider panels, public API
  teardown content, and buy-side monitoring of every vendor as mandatory MVP
  scope.

### Fixed

- Aligned product, workflow, developer, architecture, Factory-integration, and
  repository-agent documents with the compiled v3 control state.
- Removed stale language that called the active compiled directory a checked-in
  v2 or “next” generation.
- Made factoryd readiness and product implementation explicitly separate from
  successful repo-local planning compilation.
- Recorded current generic-success placeholders for unimplemented `record`,
  `verify`, `trace`, `demo`, `share`, `eval`, and migration-runtime commands as
  M0 blockers; this planning rebaseline does not claim the typed fail-closed
  command behavior is implemented.
- Stopped describing recorder, replay, live verification, reporting, GitHub
  delivery, migration patching, or bounded-agent execution as implemented.

### Security

- Separated API-provider disclosure from model-provider disclosure and
  required exact consumer authorization for model endpoint, model/version,
  credentials, network, logging/training/retention, tools, paths, and resource
  budgets.
- Treated repository content, provider material, tool output, and model output
  as untrusted data that cannot widen policy, authority, tools, writable paths,
  network, disclosure, credentials, or budgets.
- Prohibited agent self-approval, self-verification, default-branch writes, and
  automatic merge.
- Required exact-candidate deterministic verification, human review and merge,
  fail-closed host isolation for repository commands, and separate grants for
  repository read/write, commands, models, registries, sandboxes, remote
  branches, and draft PRs.
- Kept raw consumer code, diffs, logs, traces, prompts, responses, credentials,
  and private evidence outside public source and API-provider visibility by
  default.
- Kept Factory worker `approval`, `credentials`, and `network` grants separate
  from private Lumyn product authorization.
- Prevented factoryd dispatch while its mission is paused or its external
  profile, runtime bundle, schemas, or active-mission semantics are
  unqualified.

### Historical Unreleased Rebaseline

Before v3, this unreleased repository briefly carried a deterministic-first,
packet-PKI, provider-status, connection-receipt, and receipt-backed campaign
plan with 62 acceptance items across `M0` through `M10`. That design is retained
in ADR-0002, Git history, and immutable historical Factory evidence so its
decisions and proofs remain auditable. It is not current product, billing,
dispatch, validation, or rollout authority.
