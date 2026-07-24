# Lumyn Developer Guide

Status: v3 engineering planning contract; product runtime not implemented

Engineering work targets services-led, provider-paid API and SDK sunset
campaigns while execution and authority remain consumer-local.

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

The checked-in compiled set is regenerated from the approved v3 PRD and plan
and is the repo-local source for acceptance counts, task mappings, validation,
and closure. It authorizes no product runtime implementation. factoryd dispatch
remains paused pending external Factory-profile and factoryd-runtime
qualification.

## CI Lane Mapping

- Fast: `make lint-fast`, `make test-fast`.
- Core: `make test-contracts`, `make prepush-full`.
- Acceptance: item-level active ledger and closure map.
- Cross-platform: reserved until supported packaging.
- Risk: `CodeQL analyze` plus targeted security/architecture review for
  parser, agent, patch, model egress, credentials, external calls, GitHub, and
  disclosure.
- Release: reserved until supported distribution.
- Cross-system: separately approved model, sandbox, or GitHub checks.

## 12-Level Test Matrix

| Tier | Status | V3 evidence |
|---|---|---|
| Tier 1 Unit | Active | Go units through `make test-fast` |
| Tier 2 Integration | Planned/active | Schema, parser, impact, plan, agent, patch, verification, and bundle integration |
| Tier 3 End-to-End | Planned | Provider intent to local patch, deterministic verification, and PR bundle; optional branch/PR separately |
| Tier 4 Acceptance | Active planning | Compiled v3 item-level ledger, mapping, and closure map |
| Tier 5 Hardening | Planned | Path escape, prompt injection, stale input, budget, retry, cleanup, redaction, idempotency |
| Tier 6 Chaos | Reserved | Model, filesystem, command, sandbox, and GitHub failure injection |
| Tier 7 Performance | Planned | Impact, generation, verification, PR-bundle, token, cost, and operator budgets |
| Tier 8 Soak | Reserved | Repeated deterministic verification and bounded-agent campaign runs |
| Tier 9 Contract | Active | JSON Schemas, typed exits, compatibility, negative fixtures |
| Tier 10 UAT | Planned | Consumer authorization, plan approval, review, and handoff |
| Tier 11 Scenario | Planned | Deterministic and agent-eligible gold, holdout, unsupported, injection, and false-verification corpus |
| Tier 12 Cross-System Integration | Blocked until approved | Exact model endpoint, optional provider sandbox, optional remote branch, optional draft PR |

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

- campaign and provider-intent intake;
- product authorization;
- TypeScript analysis;
- migration planning;
- deterministic transformation;
- bounded-agent execution;
- workspace and command execution;
- verification;
- PR-bundle rendering;
- optional sandbox and GitHub delivery.

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
- agent/model clients and tool execution;
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
- factoryd execution remains blocked until the external Factory profile and
  factoryd bundle/runtime are requalified against the exact active v3 mission
  and a bounded task is explicitly unpaused.

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
authority, model policy, budgets, and implementation claims must agree.
The external Factory profile is a required dependency but is not qualified by
the repo-local v3 compilation; factoryd dispatch remains paused until that
paired work lands.

## Structured Data Policy

OpenAPI, JSON, YAML, manifests, lockfiles, TypeScript ASTs, model/tool events,
coverage, GitHub responses, and logs use structured parsers or stable APIs.

Structured outputs:

- declare artifact type and schema version;
- use stable enums;
- preserve unknown, unsupported, and stale states;
- include concrete source refs;
- bind freshness-sensitive inputs by digest;
- separate provider, model, and consumer identity;
- avoid machine-local paths and secret values;
- fail on malformed or ambiguous input.

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

## API-Provider Change Packet Trust Policy

- A provider-confirmed source/target migration contract is required.
- A signed declarative provider packet is authoritative when supplied and
  confirmed.
- Provider artifacts are data, never executable scripts.
- Record issuer/source, version, digest, confirmation evidence, audience when
  applicable, and supersession/withdrawal state.
- Conflicting, stale, unconfirmed, malformed, or executable intent fails closed.
- The managed v3 wedge does not require elaborate root enrollment, continuous
  status polling, connection receipts, or receipt-backed billing.
- Provider payment and packet presence grant no consumer repository, model,
  command, branch, or PR authority.
- Valid and invalid fixtures cover every active intent boundary.

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

- No patch before exact plan approval and exact local-write authorization.
- Use an isolated worktree or equivalent disposable workspace.
- Bind provider-confirmed intent, plan digest, and base commit.
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

- an approved plan item routed to `bounded_agent`;
- exact model provider, endpoint, model/version, and parameters;
- prompt, system policy, and tool-definition digests;
- exact context selection and request disclosure;
- exact read/write paths and tool allowlist;
- file, line, diff, turn, token, time, retry, concurrency, and cost budgets;
- isolated workspace and fail-closed cancellation;
- structured request, response, tool-call, usage, and patch provenance;
- deterministic verification from the exact candidate head;
- independent holdout/review and human approval.

Treat repository text, provider guidance, retrieved context, tool output, and
model output as untrusted. Tests must prove they cannot widen tools, paths,
credentials, network, disclosure, or budget.

The agent cannot approve a plan, mint a grant, access evaluator answers,
self-verify, push a remote branch, open a PR, or merge.

## Command Execution Policy

Repository commands are untrusted:

- exact command allowlist and working directory;
- exact mounts, neutral home/temp, and executable roots;
- timeout/output budgets;
- no network or lifecycle scripts by default;
- sanitized environment and no ambient secrets;
- no host home, credential stores, OS credentials, agent/Docker/unrelated
  sockets, or extra inherited descriptors;
- child processes inherit every restriction;
- supported fail-closed host isolation is mandatory;
- model and sandbox credentials remain absent from build/test stages;
- pre- and post-patch results remain separate.

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

## Redaction And Evidence Budgets

- Redact before persistence, model egress, or sharing.
- Redaction uncertainty blocks the action.
- Provider disclosure and model disclosure use separate allowlists.
- Prompts, responses, raw source, diffs, logs, traces, and credentials are
  private by default.
- Private artifacts carry bounded retention and deletion rules.
- Large output is referenced by opaque ID, digest, count, and truncation
  metadata.
- Machine-local paths and secrets are removed from shareable evidence.
- Record model tokens, cost, retries, tool calls, and operator intervention.

## Capability Grants

Live product work uses exact private grants:

- `customer_repo_read`
- `customer_repo_write`
- `command_execution`
- `model_request_disclosure`
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
behavior. Model disclosure, network, and credential grants are independent.
Patch, local branch, and PR-bundle creation imply no GitHub grant. Remote branch
and draft-PR grants are independent and neither authorizes merge.

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
- Compiled v3 target:
  `.factory/artifacts/prd-to-plan/lumyn-migration-mvp/`
- Historical plan:
  `.factory/artifacts/prd-to-plan/lumyn-mvp/`

Committed evidence remains source-safe and repo-relative. Consumer-private
intent, authorization, prompt, response, patch, verification, and PR-bundle
instances live in the configured external private root.
