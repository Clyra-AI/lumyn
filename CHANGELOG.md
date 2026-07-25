# Changelog

## [Unreleased]

### Added

- Added ADR-0004 and a v3.1 product/engineering planning contract for
  provider-originated API update delivery launched through services-assisted,
  provider-paid API and SDK sunset campaigns.
- Added planned Provider Change Event, Consumer Installation, event-specific
  authorization, and consented provider-status projection boundaries.
- Added a concrete first provider transport: signed versioned manifests at a
  pinned provider-controlled HTTPS origin with sequence, freshness, audience,
  contract-digest, and lifecycle checks.
- Added hybrid patch routing: deterministic transforms for exact mappings and
  a consumer-local bounded coding agent for approved repository-specific work.
- Added exact model request-disclosure, endpoint, credential, tool, path,
  token, time, retry, cost, provenance, retention, and deletion boundaries.
- Added ADR-0005 and an adapter-neutral customer-selected Agent Runner
  contract: Codex and Claude Code launch targets behind common conformance and
  live-canary gates, with Cursor deferred behind the same gate.
- Added explicit `consumer_managed` and
  `provider_sponsored_lumyn_managed` execution-funding modes, separate
  credential and usage-billing ownership, clean ephemeral sessions,
  native-agent-configuration policy, and no-silent-fallback controls.
- Added executable managed-credential broker bounds for issuer, exact audience,
  one-time redemption into a quota-bound attempt session, no refresh or
  cross-attempt replay, revocation, and usage reconciliation; managed mode is
  unavailable without vendor-native enforcement or an approved
  budget-enforcing proxy.
- Added runner-host isolation contracts for mounts, host credentials, sockets,
  descriptors, child processes, egress, cleanup, and malicious-child/tool
  tests; executable plugins, MCP servers, and hooks are prohibited for MVP.
- Added an independent verifier boundary with a fresh exact-head view,
  separate process, frozen command/configuration digests, no runner/model
  credentials, and an evidence writer unavailable to generation.
- Added exact product-action route capability unions for impact, candidate,
  verification, and repair actions, plus fail-closed agent-only repair
  authorization.
- Added machine-enforced repository-command isolation across M6–M10, including
  exact commands, mounts, neutral roots, environment, credential, socket, and
  descriptor denial, pinned qualified backend identity, hard resource quotas,
  inherited child limits, offline/lifecycle defaults, and fail-closed cleanup.
- Added agent-route topology minimum scope sets and a separate
  credential-bearing sandbox-entrypoint isolation contract with endpoint-only
  egress, teardown, cleanup, and orphan evidence.
- Added separate Agent Runner Vendor versus Model Provider roles and
  `agent_runner_network`/`agent_runner_credential` product capabilities.
- Added patch artifact and PR bundle as the no-GitHub fallback, with local
  branch, short-lived remote branch, and tested draft PR as separate delivery
  states.
- Added a source-aligned v3.1 compiled Factory control set covering context,
  risk, execution, task packets, validation, acceptance, mapping, and closure.
- Added explicit factoryd mission pause and compatibility gates for the
  external Factory profile, factoryd bundle/runtime, and exact active mission.
- Retained the Go CLI/config/result foundation, `lumyn init`, `lumyn check`,
  local source parsing, executable schemas, validation, coverage, CodeQL,
  branch protection, and Factory lifecycle controls.

### Changed

- Reframed Lumyn from generic agent-readiness evaluation and the v2
  deterministic/receipt-first migration plan to provider-originated API update
  delivery.
- Made the API Provider the initial campaign buyer while preserving API
  Consumer Organization authority over repository access, execution, model
  egress, credentials, disclosure, review, and merge.
- Made services-assisted local or consumer-CI execution the initial onboarding
  and GTM motion rather than the product identity; hosted SaaS is not required
  for the first campaign.
- Made at least one Lumyn-opened tested draft PR mandatory first-campaign
  product proof; manual patch and PR-bundle handoff remain fallback and cannot
  close automated-delivery acceptance.
- Made Provider Change Event and Consumer Installation semantics four direct
  acceptance units, bringing the active ledger to 53 items, and required a
  composed provider-event-to-verified-draft-PR proof that excludes imported
  manual candidates and standalone PR creation.
- Defined installation action modes as ceilings and separated exact per-event
  approval from bounded installed preauthorization; durable installations
  store token-issuance policy, never GitHub tokens.
- Made the coding agent a replaceable adapter and explicit status-quo
  comparator rather than Lumyn's proprietary differentiation.
- Made the API Consumer Organization select the exact qualified Agent Runner
  adapter/version and funding route. Agent output and self-reported tests
  remain generation evidence; independent exact-head verification remains
  Lumyn-owned.
- Defaulted `agent_execution_policy` to `disabled` so notify-only, scan-only,
  and deterministic-only installations need no runner credential; any
  agent-assisted route now pauses for explicit configured authorization.
- Required one same-run first-campaign proof from authenticated provider event
  and installed preauthorization through an organically agent-assisted item on
  the consumer-selected qualified runner, independent exact-head verification,
  a tested Lumyn-opened draft PR, and a consented provider-received status
  projection; separate agent, delivery, or reporting runs do not qualify.
- Split remote branch push from draft-PR creation into separate atomic M9
  actions, removed aggregate generation/delivery routes, and made M9/M10
  compositions dereference and exactly match their M4/M6/M7/M8/M9 source
  route contracts.
- Made raw code, diffs, prompts, responses, agent sessions, tool traces, logs,
  and credentials never API-provider-visible; provider disclosure is limited
  to enumerated, consumer-consented status or aggregates.
- Kept provider-confirmed migration intent authoritative and signed
  declarative packets supported, while deferring mandatory packet PKI,
  continuous provider status, connection receipts, acknowledgements, and
  receipt-backed billing.
- Made generation provenance independent from deterministic, exact-candidate
  verification strength.
- Replaced checked-in v2 active-control claims under
  `lumyn-migration-mvp` with the regenerated v3.1 compiled control set.
- Preserved `.factory/artifacts/prd-to-plan/lumyn-mvp/`, ADR-0002, and their
  lifecycle evidence as immutable historical records.
- Clarified that the repo-local v3.1 compilation is planning and validation
  authority only. It authorizes no product runtime implementation or live
  action.
- Kept factoryd dispatch paused until the external Factory
  `profiles/lumyn.yaml` profile and factoryd bundle/runtime are separately
  requalified and a bounded task is explicitly unpaused.
- Scoped that pause to factoryd dispatch while retaining a separately approved
  attended execution path through the same task packets and lifecycle gates.
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

- Made recognized but unimplemented `record`, `verify`, `trace`, `demo`,
  `share`, and `eval` commands return typed `command_not_implemented` results
  with exit code `2` instead of false-green success envelopes.
- Stopped assigning evaluation-mode metadata to `lumyn init` and
  `lumyn check`.
- Clarified the retained result/evidence schema v1.0 `provider_metadata` key as
  Model Provider metadata only, added an optional `model_provider` role
  discriminator for new output, preserved legacy payload validation, and
  documented the provenance-preserving v2.0 rename path in ADR-0006.
- Aligned product, workflow, developer, architecture, Factory-integration, and
  repository-agent documents with the compiled v3.1 control state.
- Removed stale language that called the active compiled directory a checked-in
  v2 or “next” generation.
- Made factoryd readiness and product implementation explicitly separate from
  successful repo-local planning compilation.
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
