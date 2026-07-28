# Lumyn Workflow Contract

Version: 3.1
Status: Normative v3.1 planning contract; product runtime not implemented

## Work Signal

Lumyn accepts governed work from:

- `docs/product/prd.md`;
- `docs/product/plan.md`;
- the source-aligned compiled v3 control set under
  `.factory/artifacts/prd-to-plan/lumyn-migration-mvp/`;
- governed post-PRD findings under `.factory/artifacts/post-prd/`;
- GitHub issues or pull requests mapped to the approved product contract.

The active compiled set is repo-local planning and validation authority. It is
not factoryd execution authority. factoryd dispatch remains paused until the
external Factory profile and factoryd bundle/runtime are requalified against
the exact v3.1 mission.

The prior `.factory/artifacts/prd-to-plan/lumyn-mvp/` package and its task,
pilot, and PR-lifecycle evidence are immutable historical records.

## Product Delivery Flow

The provider-originated v3.1 workflow preserves these state boundaries:

```text
provider-paid sunset campaign
-> provider-confirmed Provider Change Contract
-> provider-originated change event
-> consumer campaign invitation
-> consumer installation of provider/channel and action policy
-> event-specific authorization derived without widening
-> customer-authorized repository analysis
-> no-write migration plan
-> exact per-event approval or installed-preauthorization policy evaluation
-> explicit action-specific consumer execution grants
-> [deterministic transformation | bounded-agent generation | needs input]
-> isolated local patch
-> deterministic repository and workflow verification
-> patch artifact + optional local branch + PR bundle fallback
-> consumer-authorized non-default remote branch
-> consumer-authorized tested draft PR
-> consumer review and merge
-> separately consented, event-bound provider rollout status
```

A Lumyn Operator may coordinate campaign intake, provider confirmation,
consumer onboarding, and review. Operator participation grants no ambient
repository, command, model, credential, branch, PR, or merge authority.

A Provider Change Contract is authoritative when accountably confirmed. It and
its provider event remain untrusted, non-executable input and cannot widen a
Consumer Installation. Duplicate, stale, conflicting, superseded, withdrawn,
wrong-audience, and unauthenticated events fail closed. Elaborate PKI,
universal event distribution, connection receipts, provider acknowledgements,
and receipt-backed billing are deferred from the v3.1 campaign prerequisite.
The first concrete channel is a signed versioned JSON manifest at an exact
provider-controlled HTTPS URL pinned with its campaign key during
installation. It embeds the Provider Change Contract or names the exact
provider-controlled HTTPS URL whose retrieved bytes must match the event
digest. Attended import is recovery, not provider-channel or
installed-preauthorization proof.

No downstream state implies an earlier permission. Provider payment is not
repository consent. Plan approval is not file-write approval. A patch is not a
branch. A local branch is not a remote branch. A PR bundle is not a GitHub
write. A passing repository check is not merge authority.

This document describes planned behavior. M0 must replace every current
generic-success placeholder with a typed nonzero result before implementation
acceptance can close; a placeholder is not an implemented command.

## Normal Factory Chain

1. `scout-context`
2. `execution-compiler`
3. `task-executor`
4. `validation-gate`
5. `code-review` when required by risk
6. `holdout-evaluator` when selected by policy
7. `trace-grader` when selected by policy
8. `evidence-attestor` when selected by policy
9. `commit-push`
10. `post-merge-monitor`
11. `repair-feedback` when validation, review, shipping, or closure fails

Independent lifecycle workers run after implementation validation and before
shipping. They emit task-bound, current-work-proof, passing artifacts in the
trusted Factory evidence root. The implementation worker cannot self-review,
self-grade, self-attest, or write PR-lifecycle evidence.
For historical M1 the exact completed order was validation, independent
`code-review`, then independent holdout provisioning. Each report remains
schema- and candidate-bound; review cannot cite the later holdout result, and
`make lifecycle-evidence` rejects missing, replayed, widened, self-produced,
out-of-order, non-promoting, or later-workspace-reinterpreted evidence. The
consumed M1 authorization and task-specific review exception cannot be reused.

Acceptance item IDs remain the closure source. Tasks, waves, and delivery slices
are sequencing and coverage lenses only.

## Approval Gates

### Product Planning

- Plan approval is required before implementation or product execution.
- The PRD, plan, ADRs, operating docs, Factory profile, active compiled
  artifacts, validators, examples, and fixtures must agree before v3.1 dispatch.
- The repo-local v3.1 compiled set does not close the separate external-profile
  and factoryd-runtime qualification gate.
- Services-assisted discovery may collect provider intent and anonymized planning
  evidence under an approved privacy scope. It grants no product runtime
  authority.
- Before collecting identifiable provider or consumer evidence, approve exact
  fields, participant consent, storage, retention, deletion, disclosure,
  confidentiality, and evidence ownership.
- A provider-paid engagement must identify the sunset objective, accountable
  provider operator, intended consumer cohort, provider-confirmed semantics,
  commercial scope, and success decision.
- Public fixtures prove engineering feasibility only. They do not close product
  demand, repository consent, or paid campaign evidence.
- The current rebaseline authorizes no product implementation, consumer
  repository, model endpoint, external credential, command, network, branch,
  or PR action.
- Task- and campaign-level product-authority arrays are capability universes,
  not grants. Every product side effect requires one named route whose exact
  required plus conditionally selected union is frozen before action; composed
  runs reuse validated routes rather than authorizing their aggregate union.

### Customer Repository Read

`customer_repo_read` names:

- campaign and provider-confirmed change intent;
- repository and selected package root;
- readable paths and exclusions;
- expiry and revocation;
- retention, deletion, and evidence owner.

Impact analysis cannot start without this grant and cannot mutate repository
state.

### Migration Plan

The no-write plan names:

- every affected, excluded, unsupported, and uncertain site;
- every intended file, dependency, and lockfile change;
- the route for each change: deterministic, bounded agent, manual input, or
  blocked;
- expected commands and verification levels;
- model disclosure, endpoint, credential, network, tool, and budget needs;
- delivery targets: patch, local branch, PR bundle, remote branch, or draft PR;
- residual risk.

The Consumer Maintainer approves the exact plan digest. A changed plan requires
new approval.

### Customer Repository Write

`customer_repo_write` names:

- approved plan digest;
- base commit;
- writable paths;
- file, line, and diff budgets;
- isolated local workspace;
- rollback posture;
- expiry and revocation.

Read authorization never implies write authorization. Local write never implies
remote branch write.

### Bounded Agent

Agent-assisted generation requires all of:

- an authorized plan item explicitly routed to `agent_assisted`;
- a consumer-selected, qualified exact Agent Runner adapter and version;
- an ephemeral clean session that does not resume personal or unrelated agent
  history;
- exact model provider, endpoint, model/version, and parameters;
- exact prompt, system policy, and tool-definition versions or digests;
- allowed read paths, writable paths, and tool calls;
- file, line, diff, turn, token, time, retry, and cost budgets;
- isolated workspace and fail-closed cancellation;
- separate Agent Runner network and credential grants when the adapter requires
  them;
- separate model request disclosure, network, and credential grants;
- provenance for every request, response, tool call, attempt, and patch digest;
- independent verification and human review.

Repository files, provider material, tool output, and model output are untrusted
data. They cannot alter system policy, permissions, budget, network allowlists,
tools, or writable paths. A model may propose a patch but cannot approve its
plan, change its authorization, label itself verified, open a remote branch, or
create a PR.

Do not claim deterministic patch reproduction for agent mode. Deterministic
verification remains required for the exact candidate head.

### Agent Runner Selection And Funding

The first supported adapters are `codex` and `claude_code`, but neither is
advertised or used live until its pinned adapter version and executable digest
from an approved source pass the common contract suite and an explicitly
approved live canary. Its auth mode and entitlement class must allow the
intended non-interactive local or CI use. `cursor` is deferred until it passes
the same gate. `agent_execution_policy` defaults to `disabled`; notify-only,
scan-only, and deterministic-only runs require no Agent Runner credential. An
agent-assisted route pauses until the Consumer Maintainer explicitly configures
and authorizes one exact qualified route. Lumyn does not silently change
adapter, version, Model Provider, model, endpoint, credential owner, or
usage-billing owner.

Every agent-enabled run declares one execution-funding mode:

- `consumer_managed` is the default configured mode. The API Consumer
  Organization owns and authorizes a qualifying agent account, enterprise
  subscription, API credential, or local runtime and owns third-party usage
  billing. The route must expose actual downstream model identity and permit
  non-interactive automation. Lumyn receives no reusable credential.
- `provider_sponsored_lumyn_managed` is optional. The API Provider funds the
  campaign, while Lumyn owns approved agent/model usage billing and injects
  only a task-scoped brokered credential into the consumer-authorized local or
  CI boundary. The broker binds issuer, installation/event/plan/attempt and
  runner/model audience and maximum one-hour TTL. One-time redemption creates
  one attempt-scoped session; multiple calls are allowed only within hard token
  and cost quotas, with no refresh, post-attempt replay, or cross-attempt
  reuse. Revocation and reconciliation require a vendor-native bounded
  credential or approved budget-enforcing proxy; otherwise the managed route
  is unavailable. The API Provider receives no credential, code, context, or
  agent-session access.

Static native user or project rules and memories are ignored unless explicitly
selected in the installation. If selected, their identity and digest are
recorded, they are treated as untrusted context, and they cannot widen Lumyn's
tool, path, network, credential, disclosure, or budget authority. Executable
plugins, MCP servers, and hooks are prohibited for the MVP. Every attempt
resolves the approved executable by canonical path and digest, uses neutral
home/config roots, and rejects repository-local PATH shadowing.

The runner process receives only explicit read-only and writable mounts, no
host home or OS credential store, no ambient service sockets or unrelated
inherited descriptors, inherited child-process limits, host-enforced egress,
and evidence-backed cleanup. If the host cannot enforce that boundary, the
runner does not launch.

### Agent Runner Network

`agent_runner_network` names the exact Agent Runner Vendor control-plane
endpoint and operations. It is distinct from `model_network`, including when
the selected runner brokers or multiplexes access to a downstream Model
Provider. Opaque or changing downstream model routes do not qualify for the
MVP.

### Agent Runner Credential

`agent_runner_credential` names credential owner, auth mode, entitlement class,
injection environment, scopes, expiry, revocation, prohibited persistence and
logging, and evidence. A consumer account credential and a direct Model
Provider credential are not interchangeable. Missing, ambient, reusable,
owner-ambiguous, or entitlement-invalid credentials block execution.

### Model Request Disclosure

`model_request_disclosure` names:

- exact source, metadata, and context classes permitted to leave the consumer
  plane;
- prohibited code, secrets, credentials, PII, and production data;
- redaction and minimization;
- Agent Runner Vendor and downstream Model Provider processing, logging,
  training, retention, deletion, and regional posture;
- prompt/response private-artifact policy;
- evidence and expiry.

API-provider disclosure and model-provider disclosure are distinct. Neither
implies the other.

### Model Network

`model_network` names:

- exact endpoint and operation allowlist;
- request/response size;
- token, time, retry, concurrency, and cost budgets;
- model/version pin and failure behavior;
- expiry and evidence.

No generic internet access, fallback endpoint, silent model or adapter upgrade,
or undeclared tool network is allowed.

### Model Credential

`model_credential` names:

- credential environment and exact scopes;
- isolated injection stage;
- expiry and revocation;
- prohibited persistence and logging;
- evidence.

Model credentials remain absent from repository commands, build/test stages,
provider artifacts, and committed evidence.

### Command Execution

`command_execution` names:

- exact commands and working directory;
- exact read-only and writable mounts;
- neutral home and temp roots;
- executable and toolchain roots;
- timeout and output budgets;
- dependency lifecycle-script posture;
- network posture;
- environment variable classes, never secret values;
- local-socket and inherited-descriptor policy;
- process-tree limits and child-process inheritance;
- exact isolation backend, version, configuration digest, qualification
  digest, and host platform;
- hard CPU-time, memory, PID, process-tree-depth, disk, and open-file quotas;
- absolute denial of host OS credential stores and credentials;
- only explicit task-scoped credential injection or broker handles authorized
  by the named credential grant;
- a supported fail-closed isolation backend.

Repository commands use no host home or host credential stores, no
agent/Docker/unrelated service sockets, and no extra inherited descriptors.
Network and lifecycle scripts remain disabled unless their separate exact
route grants select them. If the boundary cannot be enforced, the command does
not run.

An agent action also freezes exactly one `agent_route_topology`. Local runtime
grants no external egress. Runner-mediated execution requires Agent Runner
network/credential and model-disclosure scopes; direct-model execution
requires model network/credential/disclosure scopes; hybrid requires both.
Package-registry read is independently conditional.

### Package Registry

`package_registry_read` independently names exact Node/npm and toolchain
versions, package/version/integrity inputs, registry or immutable snapshot,
read-only network allowlist and budget, lifecycle-script posture, expiry, and
evidence.

### Provider Sandbox

Optional sandbox verification independently requires:

- `sandbox_request_disclosure` for exact payload and data classes plus provider
  logging, retention, and deletion;
- `sandbox_network` for the exact non-production endpoint, operations,
  namespace, request/write budgets, retry, cleanup, and orphan handling;
- `sandbox_credential` for non-production scopes, injection, expiry, and
  revocation.

Production data, credentials, PII, and irreversible production mutations are
prohibited.

### Patch, Branch, PR Bundle, And Draft PR

The delivery ladder is:

1. patch artifact;
2. optional local branch;
3. PR bundle;
4. separately authorized remote branch;
5. separately authorized draft PR;
6. consumer review and merge.

An installation may remain scan-only or use manual fallback, but the first
provider campaign does not pass until at least one tested outcome reaches step
5 through Lumyn's short-lived delivery path.

The PR bundle contains the provider-confirmed intent ref, plan digest, base/head
identity, patch provenance, model provenance when applicable, deterministic
verification evidence, excluded/unsupported scope, residual risk, and suggested
title/body. It requires no GitHub credential.

`github_branch_write` names the repository, authorized non-default namespace,
base commit, token expiry, idempotency, and rollback. It grants no PR authority.
The durable Consumer Installation stores only token-issuance policy; an
approved local or CI broker mints the short-lived GitHub App installation
token at runtime.

`github_pr_write` names the repository, authorized remote branch, base branch,
draft-only posture, token expiry, idempotency key, and approved plan/evidence
refs. It grants no remote-branch or merge authority.

The product never writes the default branch or auto-merges.

`lumyn update --event` is the composed product path. `EXP-003` requires one
provider-channel event to traverse installation, impact, plan, Lumyn
deterministic or agent-assisted generation, independent verification, remote
branch, draft PR, and local event-bound status projection or explicit
reporting decline. Standalone PR creation, attended event import, imported
manual candidates, and manual bundles remain valid recovery surfaces but
cannot close that proof. Provider transmission is optional for technical
delivery; the qualifying pilot run must transmit its consented projection.
`PILOT-003` further requires that this qualifying composed run itself contain
an organically `agent_assisted` plan item on a consumer-selected qualified
runner, pass independent exact-head verification, open the Lumyn draft PR, and
transmit the bound status projection. A separate agent run plus deterministic
delivery run does not qualify, and deterministic work cannot be rerouted to
manufacture proof.

### Optional Provider Reporting

`provider_reporting` lists only the exact fields the consumer permits Lumyn to
share. Raw source, diffs, logs, traces, prompts, responses, agent sessions, and
credentials are never shareable with the API Provider. Missing reporting
consent does not block local patch, branch, PR-bundle, or otherwise authorized
draft-PR delivery.

### Artifact Retention And Deletion

`artifact_retention` names exact artifact classes, storage boundary, TTL,
expiry, and evidence owner.

`artifact_deletion` names revocation/expiry triggers, deletion scope, receipt
owner, retry, and orphan route. Cleanup cannot extend authority or rewrite
historical closure.

### Merge

- Consumer Maintainers own review and merge for product-generated migration
  work.
- Lumyn repository merges follow the Factory lifecycle.
- No product or Factory grant authorizes auto-merge.

## Artifact Rules

- Product source: `docs/product/prd.md`.
- Human plan: `docs/product/plan.md`.
- Compiled v3 target:
  `.factory/artifacts/prd-to-plan/lumyn-migration-mvp/`.
- Historical plan: `.factory/artifacts/prd-to-plan/lumyn-mvp/`.
- Task evidence: `.factory/artifacts/task-runs/<task_id>/`.
- Independent evidence:
  `.factory/artifacts/lifecycle-evidence/<task_id>/`.
- PR lifecycle:
  `.factory/artifacts/pr-lifecycle/<work_item_id>/pr-lifecycle-report.json`.
- Consumer-private runtime artifacts: configured root outside the checkout and
  public source repository.
- Public evidence: separately consented, redacted aggregate/hash evidence only.

Historical evidence is immutable and proves only its recorded semantics.

## Bootstrap Validation Lanes

- Fast: `make lint-fast`, `make test-fast`
- Coverage: `make test-coverage`
- Contract: `make test-contracts`
- Full: `make prepush-full`
- Historical lifecycle evidence for immutable M1 markers, reports, retained
  original-head bundle, landed-content proof, and exact-main monitoring:
  `make lifecycle-evidence`
- Risk: GitHub Actions `CodeQL analyze`
- Acceptance: item-level ledger and closure map
- Cross-system: separately approved model, sandbox, or GitHub integration

## PR Lifecycle Baseline

- Local gates: `make prepush-full`, then `make lifecycle-evidence` to verify
  M1's immutable historical candidate, lifecycle, exact-main, and non-reusable
  exception bindings without comparing later work to M1's frozen diff.
- GitHub validation: `validate`.
- Security scanner: `CodeQL analyze`.
- Owner review: `.github/CODEOWNERS`.
- Remote protection audit: `make audit-remote-protection`.
- Structured `code-review` for model, agent, patch, credential, external-call,
  GitHub, schema, policy, or other high-risk behavior.
- Independent holdout, trace-grade, and attestation evidence when selected.
- Shipping through `commit-push`.
- Monitoring through `post-merge-monitor`.

Passive Codex review settle is required before merge.
Green CI alone is not merge-ready. Do not merge manually through
`gh pr merge`, the GitHub UI, or a
connector before the configured terminal latest-head review signal. A merge
without that evidence is a process escape and requires recorded repair.

## Validation And Proof Rules

- Capture the repository baseline before mutation.
- Keep generation mode separate from verification strength.
- Bind verification to the exact candidate head and pinned environment.
- Run verification in a fresh process and view with frozen command and
  verification-configuration digests, no Agent Runner/model credentials, and
  no generation-owned verification-evidence write handle.
- Record agent/model provenance without exposing private prompts or responses.
- Treat pre-existing failures separately.
- Stale evidence cannot close acceptance.
- Boundary, cleanup, redaction, authorization, provenance, or proof failure
  blocks stronger labels.
- Public fixtures prove engineering behavior, not demand.
- Unimplemented commands must return typed nonzero results.

## Factory And factoryd

The checked-in Factory configs are templates with no active grants. Factory
worker `approval`, `credentials`, and `network` capabilities never substitute
for Lumyn product authorization.

If a reviewed task declares `conditional_factory_capabilities`, a future
runtime may activate only one frozen task/action mode and its exact sorted
selected capability set. Every selected grant must carry the same
task-bound evidence ref, scope digest, and expiry; partial, extra, stale,
wildcard, or cross-task activation fails closed.

The factoryd templates remain mission-paused. Do not select or execute v3 work
through factoryd until a separate reviewed change:

- rebaselines the external Factory `profiles/lumyn.yaml` profile;
- qualifies the factoryd bundle/runtime against the exact active mission and
  schemas;
- preserves repo-local validator and operating-doc agreement; and
- explicitly unpauses one bounded task with positive budgets and required
  grants.

The current compiled set remains usable for planning and deterministic
repo-local validation. Future product workers emit task-scoped evidence and
cannot mutate canonical planning truth.

## Post-PRD Findings

Material audit, review, pilot, or recommendation findings are saved as governed
post-PRD artifacts. A human explicitly promotes a finding into product scope
and regenerates every affected planning surface.

## Stop Conditions

Stop and request a human decision if:

- provider, Lumyn operator, model provider, and consumer authority are unclear;
- provider intent is unconfirmed or provider material would execute code;
- the approved plan, base commit, or intent binding has drifted;
- read-only impact would mutate state;
- model disclosure, endpoint, credential, tools, network, retention, or budgets
  are absent or ambiguous;
- untrusted content attempts to widen authority or tool access;
- agent generation would self-approve or self-verify;
- a local write, command, model, registry, sandbox, branch, PR, retention, or
  deletion boundary is unclear;
- host isolation cannot be enforced;
- required business input is missing;
- production access would be required;
- redaction or evidence freshness is uncertain;
- private runtime data would enter a checkout or public repository;
- held-out answer material would be visible to an implementation worker;
- evidence would be labelled above its proven boundary;
- provider-visible data exceeds consent;
- a patch, local branch, or PR bundle is treated as remote-write authority;
- default-branch write or auto-merge is requested;
- a task treats the v3 planning rebaseline as product implementation
  authorization;
- factoryd dispatch is attempted while the mission is paused or its external
  profile, runtime, bundle, or active-mission semantics are unqualified;
- the PRD, plan, Factory profile, task packets, acceptance ledger, validation
  contract, mapping, and closure disagree;
- required validation, coverage, CodeQL, review, `commit-push`, post-merge, or
  item-level evidence is missing without an approved exception.
