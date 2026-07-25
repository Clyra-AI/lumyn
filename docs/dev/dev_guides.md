# Lumyn Developer Guide

Status: v3.1 engineering planning contract; product runtime not implemented

Engineering work targets provider-originated API update delivery launched
through services-assisted, provider-paid sunset campaigns while execution and
authority remain consumer-local.

## Toolchain Pins

| Tool | Version |
|---|---:|
| Go | `1.26.5` |

Module path: `github.com/Clyra-AI/lumyn`.

The Go core remains authoritative for artifact, authorization, impact,
planning, agent orchestration, patch, verification, and delivery contracts.
Any TypeScript parser, model client, tool subprocess, or SDK requires a pinned
dependency, bounded interface, license/security review, and task evidence.

Exact Node/npm, registry or immutable snapshot, package-integrity, and
toolchain pins are required before `package-lock.json` mutation.

## Dependency Pins

- `github.com/santhosh-tekuri/jsonschema/v5 v5.3.1`: executable JSON Schema
  validation.

New dependencies are pinned, justified, scanner-covered, and exercised by a
failing test or fixture before implementation.

No model SDK or endpoint is implicitly approved by this planning rebaseline.

## Validation Matrix

- `make lint-fast`: repo contract, layout, policy, and Go vet.
- `make test-fast`: Go unit tests.
- `make test-coverage`: first-party Go coverage gate.
- `make test-contracts`: Go tests, schema tests, active-plan validation,
  historical evidence validation, and repo-pack self-tests.
- `make prepush-full`: full local gate before PR or merge.
- `make audit-remote-protection`: networked GitHub protection audit.

The checked-in compiled set is regenerated from the approved v3.1 PRD and plan
and is the repo-local source for acceptance counts, task mappings, validation,
and closure. It authorizes no product runtime implementation. A v3.1 source
update has landed in the external Factory profile, but the active compiled
controls still record Factory profile/runtime compatibility as unqualified.
factoryd dispatch remains paused pending full control regeneration,
bundle/runtime qualification, and an explicit bounded-task unpause.

## CI Lane Mapping

- Fast: `make lint-fast`, `make test-fast`.
- Core: `make test-contracts`, `make prepush-full`.
- Acceptance: item-level active ledger and closure map.
- Cross-platform: reserved until supported packaging.
- Risk: `CodeQL analyze` plus targeted security/architecture review for
  parser, agent, patch, model egress, credentials, external calls, GitHub, and
  disclosure.
- Release: reserved until supported distribution.
- Cross-system: separately approved Agent Runner, model, sandbox, or GitHub
  checks.

## 12-Level Test Matrix

| Tier | Status | V3 evidence |
|---|---|---|
| Tier 1 Unit | Active | Go units through `make test-fast` |
| Tier 2 Integration | Planned/active | Schema, parser, impact, plan, agent, patch, verification, and bundle integration |
| Tier 3 End-to-End | Planned | Provider event and installed policy to deterministic verification and tested draft PR; local bundle fallback separately |
| Tier 4 Acceptance | Active planning | Compiled v3.1 item-level ledger, mapping, and closure map |
| Tier 5 Hardening | Planned | Path escape, prompt injection, stale input, budget, retry, cleanup, redaction, idempotency |
| Tier 6 Chaos | Reserved | Model, filesystem, command, sandbox, and GitHub failure injection |
| Tier 7 Performance | Planned | Impact, generation, verification, PR-bundle, token, cost, and operator budgets |
| Tier 8 Soak | Reserved | Repeated deterministic verification and bounded-agent campaign runs |
| Tier 9 Contract | Active | JSON Schemas, typed exits, compatibility, negative fixtures |
| Tier 10 UAT | Planned | Consumer authorization, plan approval, review, and handoff |
| Tier 11 Scenario | Planned | Deterministic and agent-eligible gold, holdout, unsupported, injection, and false-verification corpus |
| Tier 12 Cross-System Integration | Blocked until approved | Qualified Codex/Claude Agent Runner, exact model route, optional provider sandbox, and required short-lived remote branch/draft-PR pilot path |

Runner-ready packets cite each applicable tier or an approved non-applicable
reason.

## Coverage Gates

| Scope | Threshold | Enforcement |
|---|---:|---|
| Go first-party packages overall | `>= 75%` | `make test-coverage` and CI |
| Stable command or core packages | `>= 85%` | `make test-coverage` |

Coverage is not a substitute for schema fixtures, held-out scoring,
prompt-injection tests, golden deterministic patches, proof scorecards, CodeQL,
or cross-system evidence.

## Architecture Budget Gate

Source files warn at 1200 lines and fail at 2500 lines for supported source
extensions. Generated runtime, dependencies, caches, and build outputs are
excluded.

Keep these responsibilities separate:

- provider event and Provider Change Contract intake;
- Consumer Installation and event-specific authorization;
- product authorization;
- TypeScript analysis;
- migration planning;
- deterministic transformation;
- bounded-agent execution;
- workspace and command execution;
- verification;
- PR-bundle rendering;
- optional sandbox and short-lived GitHub delivery.

## CI And PR Lifecycle

The canonical lifecycle is:

1. `task-executor`
2. `validation-gate`
3. `code-review` when required
4. `holdout-evaluator` when selected
5. `trace-grader` when selected
6. `evidence-attestor` when selected
7. `commit-push`
8. `post-merge-monitor`
9. `repair-feedback` on failure

Independent lifecycle evidence must be task-bound, current, passing, and
outside the implementation worker's writable scope.

Passive Codex review settle is required before merge. Green CI alone is not
merge-ready. Do not merge manually through `gh pr merge`, the GitHub UI, or a
connector before the latest-head terminal review signal. A merge without that
evidence is a process escape and requires recorded repair or exception.

GitHub `main` remains protected by branch protection and the
`protect-main-from-direct-push` ruleset. Use `make audit-remote-protection`.

## Security Scanner Enforcement

CodeQL and risk review are required for:

- dependency additions;
- parser or generated-code intake;
- Agent Runner/model clients and tool execution;
- prompt construction and context selection;
- patch generation and filesystem writes;
- command execution;
- external network or API calls;
- credential, redaction, retention, or data-sharing behavior;
- GitHub integration;
- release-sensitive work.

Scanner failure blocks closure without a scoped approved exception.

## Bootstrap Rules

- Planning and public-fixture work uses no consumer repository, model key,
  external credential, live network, provider sandbox, or GitHub write.
- Test-first or fixture-first development is expected.
- Consumer-private runtime and identifiable campaign evidence lives outside the
  checkout and public source repository.
- Factory worker grants use only `approval`, `credentials`, and `network`.
  Exact Lumyn product grants are separate private artifacts.
- Conditional Factory grants require one frozen task/action mode, exact sorted
  selected capability set, common activation evidence/digest/expiry, and
  complete-set validation.
- Historical task, pilot, lifecycle, and closure evidence is immutable.
- Structured artifact changes include valid and invalid fixtures.
- Behavior, command, schema, artifact, permission, and evidence changes update
  docs and active Factory planning together.
- Runner-ready packets preserve acceptance IDs, paths, commands, risk,
  lifecycle gates, evidence, proof, capability, budget, stop, changelog, and
  semantic-invariant fields.
- This rebaseline authorizes no product runtime implementation.
- factoryd execution remains blocked until the complete control set is
  regenerated, its bundle/runtime is qualified against the exact active v3
  mission, and a bounded task is explicitly unpaused. The landed external
  source-profile update is not runtime proof.

## Docs Parity

User-facing sources:

- `README.md`
- `AGENTS.md`
- `WORKFLOW.md`
- `docs/product/prd.md`
- `docs/product/plan.md`
- `docs/dev/dev_guides.md`
- `docs/architecture/architecture_guides.md`
- relevant ADRs

Active planning sources:

- `.factory/artifacts/prd-to-plan/lumyn-migration-mvp/`
- Factory `profiles/lumyn.yaml`

Behavior, status, generation modes, verification labels, artifact paths,
authority, Agent Runner/model policy, budgets, and implementation claims must
agree.
A v3.1 source update has landed in the external Factory profile, but the active
repo-local compilation still records Factory profile/runtime compatibility as
unqualified and remains authoritative for factoryd. Dispatch stays paused
until complete control regeneration and factoryd bundle/runtime and exact
active-mission qualification.

## Structured Data Policy

OpenAPI, JSON, YAML, manifests, lockfiles, TypeScript ASTs, model/tool events,
coverage, GitHub responses, and logs use structured parsers or stable APIs.

Structured outputs:

- declare artifact type and schema version;
- use stable enums;
- preserve unknown, unsupported, and stale states;
- include concrete source refs;
- bind freshness-sensitive inputs by digest;
- separate API Provider, Agent Runner Vendor, Model Provider, and consumer
  identity;
- avoid machine-local paths and secret values;
- fail on malformed or ambiguous input.

## Result And Evidence Compatibility

- `lumyn.command_result` remains at schema version `1.0` for M0.
- The v1.0 serialized `provider_metadata` key is legacy Model Provider metadata
  only. New writers add `semantic_role: model_provider`; historical payloads
  without it and payloads using previously allowed unknown extensions remain
  schema-valid. Unknown extensions do not override the key's Model Provider
  meaning and are not current writer output.
- API Provider identity uses `api_provider_id` or `change_authority`. Agent
  Runner Vendor identity uses `agent_runner_vendor_metadata`.
- A future serialized rename to `model_provider_metadata` requires schema
  version `2.0`, frozen executable v1.0 schemas, dual version-discriminated
  readers, current and legacy fixtures, v2.0 invalid mixed-role fixtures, and a
  new provenance-linked artifact rather than historical in-place rewriting.
- ADR-0006 is the normative migration decision.

## Agent-Native CLI Policy

State-returning commands:

- support stable JSON;
- remain machine-readable when piped or non-interactive;
- preserve status, evidence refs, typed errors, and exit code;
- return nonzero for unimplemented behavior;
- never emit a generic success placeholder.

Help and docs must not advertise v3 commands as working before their end-to-end
acceptance passes.

## Migration Corpus Policy

Every development fixture includes:

- fixture, campaign, and change IDs;
- pinned source/target refs and digests;
- provenance, license, attribution, and redistribution posture;
- official SDK package/version;
- annotated impacted and unaffected sites;
- provider-confirmed intent and unresolved questions;
- expected allowed paths and residual risk;
- deterministic expected patch when applicable;
- agent-eligible outcome constraints when applicable;
- expected verification stage and outcome;
- unsupported, injection, or negative classification where applicable.

Visible fixtures and held-out scoring remain separate. The implementation
worker never receives held-out inputs, expected patches, answer keys, or raw
traces.

For deterministic mode, score exact patch and stable output. For agent mode,
score scope, semantic constraints, unrelated edits, repository/workflow
verification, provenance completeness, budget compliance, and human correction.
Do not require byte-identical agent output.

Public fixtures demonstrate engineering behavior only.

## Provider Change Contract And Event Policy

- An accountable Provider Operator confirms one exact source/target Provider
  Change Contract for reuse across the invited cohort.
- The first channel is a signed, versioned JSON manifest at an exact
  provider-controlled HTTPS URL pinned with the campaign key by the
  installation. The manifest embeds the Provider Change Contract or pins its
  exact provider-controlled HTTPS URL. Verify origin, key, sequence, freshness,
  retrieved-byte contract digest, audience, and lifecycle state before policy
  evaluation.
- An attended file import is labeled recovery and cannot prove
  provider-channel delivery or authorize installed-preauthorization writes.
- Every provider event binds issuer, contract digest, audience, deadline,
  severity, and supersession or withdrawal state.
- Provider contracts, events, and artifacts are data, never executable scripts.
- Duplicate, replayed, stale, conflicting, unconfirmed, unauthenticated,
  wrong-audience, superseded, withdrawn, malformed, or executable intent fails
  closed.
- The managed v3.1 wedge does not require a universal registry, elaborate root
  enrollment, connection receipts, or receipt-backed billing.
- Provider payment, contract, and event presence grant no consumer repository,
  model, command, branch, PR, or disclosure authority.
- Valid and invalid fixtures cover every active intent boundary.

## Consumer Installation Policy

- Bind provider/channel, repository and package root, audience/version
  selectors, action ceiling, authorization mode, paths, commands,
  `agent_execution_policy`, Agent Runner/model data and network posture, GitHub
  token-issuance policy, reporting, retention, deletion, disclosure, expiry,
  and revocation. Require the exact qualified Agent Runner adapter/version,
  execution-funding mode, credential and usage-billing owners, and native agent
  configuration only when `agent_execution_policy=configured`.
- Supported action modes are `notify_only`, `scan_only`, `prepare_patch`, and
  `open_draft_pr`; supported authorization modes are `per_event_approval` and
  `installed_preauthorization`.
- Agent execution defaults to `disabled`. Notify-only, scan-only, and
  deterministic-only routes need no runner credential; a later
  `agent_assisted` route pauses for explicit configuration and authorization.
- `lumyn check` performs a non-mutating runner preflight only when configured:
  canonical executable source/path/version/digest, current conformance,
  permitted non-interactive auth/entitlement, and actual downstream model-route
  identity. It collects no secret and performs no model call.
- Store no reusable Agent Runner, model, or GitHub credential. An approved
  local or CI credential broker issues a task-scoped credential only at its
  qualifying runner/model or delivery boundary.
- Managed runner/model credentials bind issuer; installation, event, plan,
  attempt, runner, and model audience; and maximum one-hour TTL. One-time
  broker redemption creates one attempt-scoped session. Multiple in-attempt
  calls are allowed only within hard token/cost quotas; refresh, post-attempt
  replay, and cross-attempt reuse are forbidden. Require revocation and
  reconciliation through a vendor-native bounded credential or approved
  budget-enforcing proxy; otherwise the managed route is unavailable.
- Freeze an immutable authorization snapshot for each event.
- Treat task- and campaign-level product-authority arrays as capability
  universes only. Freeze one named action route and its exact required plus
  conditionally selected union before each side effect; a composed campaign
  reuses validated routes and never grants their aggregate union to every
  installation.
- An event may narrow but never widen installed authority.
- Expiry, revocation, wrong audience, action mismatch, cross-repository reuse,
  authorization-mode mismatch, stored-token input, and attempted policy
  widening fail closed.

## TypeScript Impact Policy

- Use a parser/AST or comparably structured representation.
- Select and canonicalize one package/read root explicitly.
- Resolve real paths before reading.
- Reject traversal, symlink escape, out-of-root references, ambiguous roots,
  and multiple package roots.
- Detect direct imports, aliases, and wrapper uncertainty.
- Report dynamic/reflection use as uncertain.
- Exclude generated, vendored, cache, and build output by default.
- Report package-manager and lockfile posture.
- Score precision and recall separately.
- Never label uncertain scope as unaffected.

## Patch And Filesystem Policy

- No patch before a no-write plan and exact local-write authorization. Under
  `per_event_approval`, authorization includes approval of that exact plan;
  under `installed_preauthorization`, every bound plan value must satisfy the
  installed policy or the run pauses.
- Use an isolated worktree or equivalent disposable workspace.
- Bind the provider event, Provider Change Contract, Consumer Installation,
  event-specific authorization, plan digest, and base commit.
- Resolve and validate real paths before writes.
- Enforce allowed/forbidden paths and file/line/diff budgets.
- Reject symlink/path traversal escape.
- Map every edit to a plan item and generation mode.
- Record deterministic recipe provenance or model/prompt/tool provenance.
- Preserve deterministic output only for deterministic mode.
- Do not execute provider scripts or infer missing business values.
- Leave the default branch untouched.
- Record rollback, cleanup, and residual risk.

## Bounded Agent Policy

Agent mode requires:

- an approved plan item routed to `agent_assisted`;
- `agent_execution_policy=configured`;
- one consumer-selected exact Agent Runner adapter/version with a current
  executable source/digest, permitted auth mode and entitlement class, and
  current common-conformance digest;
- Codex or Claude Code for launch, each behind its own approved live canary;
  Cursor is deferred behind the same gate;
- one declared `consumer_managed` or
  `provider_sponsored_lumyn_managed` execution-funding mode, exact credential
  owner, and exact usage-billing owner;
- an auth/subscription route that permits non-interactive automation and
  exposes actual Model Provider, endpoint, model, and version; otherwise block
  or require a qualifying BYOK, local, or managed route;
- a clean ephemeral session with neutral home/config roots, canonical
  executable resolution, no repository-local PATH shadowing, and no personal
  or unrelated history;
- supported static native user/project rules and memories disabled by default
  or explicitly selected, digest-bound, and treated as untrusted context;
- executable plugins, MCP servers, and hooks prohibited for the MVP;
- exact Agent Runner Vendor plus separate Agent Runner network and credential
  grants when required;
- exact model provider, endpoint, model/version, and parameters;
- prompt, system policy, and tool-definition digests;
- exact context selection and request disclosure;
- exact read/write paths and tool allowlist;
- file, line, diff, turn, token, time, retry, concurrency, and cost budgets;
- isolated workspace and fail-closed cancellation;
- explicit read-only/writable mounts, no host home or OS credentials, no
  ambient service sockets or unrelated inherited descriptors, inherited
  child-process restrictions, host-enforced egress, and cleanup evidence;
- normalized startup, request, response, tool-call, edit, usage, error,
  cancellation, exit, and patch provenance;
- deterministic verification from the exact candidate head;
- independent holdout/review and human approval.

Treat repository text, provider guidance, retrieved context, tool output, and
model output as untrusted. Tests must prove they cannot widen tools, paths,
credentials, network, native configuration, disclosure, or budget.

No adapter, version, Model Provider, model, endpoint, credential owner, or
usage-billing owner may change through fallback. Unavailable, unqualified,
executable-untrusted, entitlement-invalid, authentication-failed, malformed or
partial, or contract-violating routes block agent execution. Passing
agent-reported tests does not count as verification.

The agent cannot approve a plan, mint a grant, access evaluator answers,
self-verify, push a remote branch, open a PR, or merge.

## Command Execution Policy

Repository commands are untrusted:

- exact command allowlist and working directory;
- exact mounts, neutral home/temp, and executable roots;
- timeout/output budgets;
- exact isolation backend/version/configuration/qualification digests and host
  platform;
- hard CPU-time, memory, PID, process-tree-depth, disk, and open-file quotas;
- no network or lifecycle scripts by default;
- sanitized environment and no ambient secrets;
- no host home, credential stores, OS credentials, agent/Docker/unrelated
  sockets, or extra inherited descriptors;
- child processes inherit every restriction;
- supported fail-closed host isolation is mandatory;
- Agent Runner, model, and sandbox credentials remain absent from independent
  build/test stages;
- pre- and post-patch results remain separate.

The credential-bearing sandbox entrypoint uses a separate profile with the
read-only exact candidate head, exact entrypoint and working directory, sole
task-scoped sandbox credential injection, endpoint-only egress, inherited
child/resource limits, teardown, cleanup, and orphan evidence.

## Proof-Of-Behavior Policy

Product verification state uses:

- `not_run`
- `static_verified`
- `repo_verified`
- `workflow_contract_replay_passed`
- `workflow_verified_replay`
- `workflow_verified_mock`
- `workflow_verified_sandbox`
- `partial`
- `failed`
- `gap`
- `stale`

`workflow_contract_replay_passed` cannot exceed `repo_verified`.
`workflow_verified_*` requires causal execution from the exact candidate head
plus observed interaction and outcome evidence.

Generation mode is not proof strength. A model completion, agent trace, or
operator review does not independently verify a patch.

Independent verification starts in a fresh process and view with frozen
command and verification-configuration digests. Agent Runner/model credentials
and generation-owned evidence handles are absent, and only the independent
verifier evidence boundary may persist exact-head verification results.

## Redaction And Evidence Budgets

- Redact before persistence, model egress, or sharing.
- Redaction uncertainty blocks the action.
- Provider disclosure and Agent Runner/model disclosure use separate
  allowlists.
- Prompts, responses, raw source, diffs, logs, traces, agent sessions, and
  credentials are never API-provider-visible. Model disclosure remains a
  separate exact allowlist; provider disclosure is limited to enumerated,
  consented campaign status or aggregates.
- Private artifacts carry bounded retention and deletion rules.
- Large output is referenced by opaque ID, digest, count, and truncation
  metadata.
- Machine-local paths and secrets are removed from shareable evidence.
- Record Agent Runner/model tokens and cost, funding mode, credential and
  usage-billing owners, retries, tool calls, and operator intervention.

## Capability Grants

Live product work uses exact private grants:

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

Every grant names target, scope, expiry, revocation, evidence, and failure
behavior. Agent Runner network and credential grants and model disclosure,
network, and credential grants are independent. Patch, local branch, and
PR-bundle creation imply no GitHub grant. Remote branch and draft-PR grants
are independent and neither authorizes merge.
`provider_reporting` is optional for M9 delivery and cannot block an otherwise
authorized draft PR; M10 campaign proof separately requires at least one
consumer-consented event-bound projection.

Automated-delivery acceptance additionally requires short-lived credentials,
non-default branch, draft-only posture, tested-candidate binding, idempotency,
and negative proof for default-branch write and auto-merge. Manual fallback
cannot satisfy that acceptance. `EXP-003` also requires the composed
provider-channel event -> installation -> impact -> plan -> Lumyn-generated
candidate -> independent verification -> branch -> draft PR -> local
status-projection path; standalone PR creation, attended event import, or an
imported manual candidate cannot close it. Provider transmission remains
optional for `EXP-003`, but M10 requires the qualifying run's consented
projection.

Wildcard targets, endpoints, paths, credentials, or budgets are invalid.

## Release Integrity

Primary design-partner distribution is an explicitly licensed,
integrity-signed local or consumer-CI package. Public OSS/self-serve and
Homebrew wait for the separate approved license, security, contribution,
support, vulnerability-response, and release-integrity gate.

Planned commands are not release claims.

## Provenance Evidence

- Task validation:
  `.factory/artifacts/task-runs/<task_id>/validation-report.json`
- Work proof:
  `.factory/artifacts/task-runs/<task_id>/work-proof-marker.json`
- Independent lifecycle evidence:
  `.factory/artifacts/lifecycle-evidence/<task_id>/`
- PR lifecycle:
  `.factory/artifacts/pr-lifecycle/<work_item_id>/pr-lifecycle-report.json`
- Compiled v3.1 target:
  `.factory/artifacts/prd-to-plan/lumyn-migration-mvp/`
- Historical plan:
  `.factory/artifacts/prd-to-plan/lumyn-mvp/`

Committed evidence remains source-safe and repo-relative. Consumer-private
intent, authorization, prompt, response, patch, verification, and PR-bundle
instances live in the configured external private root.
