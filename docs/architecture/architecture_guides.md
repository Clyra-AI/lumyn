# Lumyn Architecture Guide

Status: v3.1 executable contract architecture; product runtime not implemented

## Architecture Objective

Lumyn is the provider-to-consumer application layer for consequential API and
SDK changes. It carries one provider-confirmed change event through a
revocable Consumer Installation into repository-specific impact, a verified
candidate, a tested draft PR, and consented rollout status. Provider
sponsorship never grants access to consumer code, and bounded-agent generation
never grants a model authority over repository scope, verification, GitHub
delivery, or merge.

The architecture optimizes for:

- services-assisted campaign onboarding without ambient operator authority;
- reusable Provider Change Contracts and authenticated change events;
- consumer-installed, event-specific authorization that cannot widen;
- consumer-local or consumer-controlled CI execution;
- deterministic transforms for exact mappings;
- bounded-agent generation for approved repository-specific changes;
- deterministic, baseline-aware verification;
- exact Agent Runner and model egress, credential, network, disclosure, and
  cost controls;
- patch and PR-bundle fallback without requiring GitHub access;
- separately authorized short-lived remote branch and required pilot draft PR;
- event-bound, consented provider status with no inference from silence;
- proof-honest evidence and human merge authority.

This guide defines the active boundaries. M2 implements their executable
artifact and semantic-policy contracts, while later milestones still own the
v3.1 runtime. The compiled v3.1 control set is repo-local planning and
validation authority; neither it nor the M2 contracts authorize a live product
action.

A v3.1 source update has landed in the external Factory
`profiles/lumyn.yaml` profile, but the active compiled generation still records
Factory profile/runtime compatibility as unqualified. factoryd dispatch
remains paused. Separate reviewed work must regenerate the controls and qualify
the factoryd bundle/runtime against the exact active mission before any task is
unpaused.

## Trust And Data Planes

### Provider Change And Campaign Plane

Owns:

- provider identity and accountable campaign operator;
- the paid sunset objective and compatibility window;
- provider-confirmed source/target semantics in a reusable Provider Change
  Contract;
- signed versioned change event at a pinned provider-controlled HTTPS URL,
  embedded or exact-URL Provider Change Contract delivery, retrieved-byte
  digest, audience, monotonic sequence, freshness, deadline, severity, and
  supersession/withdrawal state;
- migration guidance, canary information, and rollback guidance;
- the invited cohort and commercial campaign decision.

Does not own:

- consumer repository access;
- consumer commands, Agent Runner/model credentials, or execution;
- raw source, diffs, logs, traces, prompts, or responses;
- consumer branch, PR, review, or merge authority.

The Provider Change Contract is authoritative change data when accountably
confirmed. It and its event cannot execute code or widen policy. Duplicate,
stale, conflicting, superseded, withdrawn, wrong-audience, or unauthenticated
events fail closed. V3.1 defers elaborate provider PKI, a universal event
network, and receipt-backed billing.
The first transport uses one campaign signing key and a pinned provider HTTPS
origin. Attended file import is recovery and cannot prove provider-channel
delivery or authorize installed-preauthorization writes.

The v2 `provider-authenticated consumer receipt-key bindings` and
provider-signed acknowledgement design remains historical context, not an
active v3 dependency.

### Consumer Execution Plane

Owns:

- revocable provider/channel installation, repository/package binding,
  selectors, action ceilings, `per_event_approval` or
  `installed_preauthorization`, GitHub token-issuance policy, expiry, and
  disclosure policy;
- immutable event-specific authorization derived without widening;
- repository and selected package root;
- read/write path scopes;
- exact per-event migration-plan approval or installed-policy plan evaluation;
- isolated workspace and local branch;
- command allowlist and host-isolation policy;
- `agent_execution_policy` set to `disabled` or `configured`; only a configured
  policy binds the exact qualified Agent Runner adapter/version,
  execution-funding mode, credential and usage-billing owners,
  clean-session/native-configuration policy, and Agent Runner/model
  disclosure, endpoint, credential, tool, and budget policy;
- dependency and package-registry policy;
- private impact, plan, prompt, response, patch, verification, and PR-bundle
  evidence;
- short-lived remote branch and draft-PR authorization;
- review and merge.

Consumer-private runtime state lives outside the checkout and public source
repository with explicit retention, deletion, and evidence ownership.

### Agent Runner Plane

The Agent Runner plane is an optional selected adapter inside the Consumer
Execution Plane, not an independent authority. A Consumer Installation
defaults `agent_execution_policy` to `disabled`; notify-only, scan-only, and
deterministic-only routes require no runner credential. An agent-assisted route
pauses until the Consumer Maintainer explicitly configures it.

Launch targets are pinned Codex and Claude Code adapters after each exact
version and executable digest from an approved source passes the common
conformance suite and an approved live canary. The selected auth mode and
entitlement class must permit the intended automation. Cursor remains deferred
behind that gate. Each attempt resolves the executable by canonical path,
rejects repository-local PATH shadowing, and starts a clean ephemeral session
with neutral home/config roots; personal history is never resumed. Native user
or project static rules and memories are ignored unless explicitly selected
and digest-bound as untrusted context. Executable plugins, MCP servers, and
hooks are prohibited for the MVP.

When configured, the installation chooses default
`consumer_managed` or optional `provider_sponsored_lumyn_managed` execution.
Credential owner and usage-billing owner remain explicit. A consumer-managed
subscription qualifies only when it permits non-interactive automation and
exposes the actual downstream model route; otherwise execution blocks or uses
a qualifying BYOK, local, or managed route. The API Provider never receives
agent selection, credential, code, context, or session authority.

The managed route uses a credential broker with exact issuer,
installation/event/plan/attempt and runner/model audience and maximum one-hour
TTL. One-time redemption creates one attempt-scoped session. Multiple
in-attempt calls are allowed only within hard token/cost quotas; refresh,
post-attempt replay, and cross-attempt reuse are forbidden. Revocation and
reconciliation require a vendor-native bounded credential or approved
budget-enforcing proxy; otherwise the route is unavailable.

The runner process uses explicit read-only/writable mounts, no host home or OS
credential store, no ambient service sockets or unrelated inherited
descriptors, inherited child-process limits, host-enforced egress, and
evidence-backed cleanup. Failure to enforce any boundary blocks launch.

The plane uses separate grants for:

- `agent_runner_network`;
- `agent_runner_credential`;
- `model_request_disclosure`;
- `model_network`;
- `model_credential`.

The Agent Runner Vendor and downstream Model Provider are recorded separately.
No adapter/version, model route, endpoint, credential owner, or billing owner
may change through silent fallback.

### Bounded Model Plane

The model plane is active v3.1 planned scope and a separate trust boundary.

It receives only consumer-authorized, minimized request classes through:

- `model_request_disclosure`;
- `model_network`;
- `model_credential`.

It is pinned to an exact endpoint, model/version, parameters, prompt/tool
definitions, and turn/token/time/retry/cost budgets. Provider material,
repository content, tool output, and model output are all untrusted data.

The model plane:

- cannot widen tools, paths, network, credentials, or budget;
- cannot approve a plan or product grant;
- cannot invoke undeclared shell or network actions;
- cannot label a result verified;
- cannot write a remote branch or PR;
- cannot expose consumer material to the API Provider;
- cannot self-grade against evaluator-controlled holdouts.

## Authority Flow

```text
API Provider
  funds campaign, confirms one change contract, publishes event
        |
        v
Provider Change And Campaign Plane
  supplies non-executable event and contract
        |
        | no repository authority
        v
Consumer Maintainer
  installs provider channel and bounded actions
  -> authorizes event-specific read-only impact
  -> chooses exact-plan approval or installed-policy evaluation
  -> [configures exact qualified Agent Runner and funding route if needed]
  -> authorizes exact local write/command/runner/model boundaries
        |
        v
Consumer Execution Plane
  deterministic transform or bounded-agent generation
  -> deterministic verification
  -> patch + optional local branch + PR bundle fallback
        |
        +--> short-lived remote branch grant
        +--> short-lived draft-PR grant
        |
        v
Consumer review and merge
        |
        v
Consented event-bound provider status projection
```

Lumyn Operator assistance is an operating role, not an authority plane.

## Product State Machines

### Provider Change Contract And Event

```text
contract_draft
-> provider_confirmed
-> event_published
-> received | superseded | withdrawn | expired
```

Duplicate, stale, conflicting, unconfirmed, unauthenticated, wrong-audience,
superseded, withdrawn, or executable intent blocks planning or mutation.
The first event transport is a signed versioned JSON manifest fetched from the
exact provider-controlled HTTPS URL and verified against the campaign key
pinned by the installation. It embeds the Provider Change Contract or names
its exact provider-controlled HTTPS URL; retrieved bytes must match the event's
contract digest.

### Consumer Installation And Authorization

```text
not_installed
-> notify_only | scan_only | prepare_patch | open_draft_pr
-> per_event_approval | installed_preauthorization
-> event_received
-> event_authorization_frozen
-> read_authorized
-> plan_ready
-> exact_plan_approved | installed_policy_satisfied
-> local_execution_authorized
-> [agent_runner_authorized]
-> [model_authorized]
-> [remote_branch_authorized]
-> [draft_pr_authorized]
-> expired | revoked
```

Action modes are ceilings and each transition is separately policy-checked. An
event can narrow but never widen the installed policy. The installation stores
no GitHub token; an approved local or CI broker mints a short-lived GitHub App
installation token only after the bound runtime delivery policy is satisfied.

### Migration

```text
not_analyzed
-> impacted
-> planned
-> deterministic_candidate | agent_candidate | needs_input | unsupported
-> locally_patched
-> verified | partial | failed | stale
-> patch_delivered
-> [local_branch_prepared]
-> pr_bundle_ready
-> [remote_branch_pushed]
-> draft_pr_open
-> merged | closed | reverted
```

Generation mode and verification strength remain independent.

## Component Boundaries

| Component | Responsibility | Must not own |
|---|---|---|
| CLI | parsing, config, JSON envelope, exits | product inference |
| Campaign intake | provider-paid scope, distribution, confirmed intent | repository authority |
| Change-event intake | issuer, audience, deadline, contract binding, supersession | executable instructions or consumer authority |
| Installation engine | provider/channel, repository, actions, expiry, revocation, disclosure | cross-repository or ambient authority |
| Source intake | pinned provider materials and SDK refs | consumer writes |
| Intent normalizer | typed changes and unresolved questions | arbitrary execution |
| Authorization engine | exact consumer grants, expiry, revocation | side-effect execution |
| TypeScript analyzer | package, imports, call sites, exclusions | file mutation |
| Impact engine | applicability and coverage | patch application |
| Migration planner | complete no-write route and budgets | approval or writes |
| Workspace manager | isolated workspace, safe paths, local branch | semantic decisions |
| Deterministic transformer | exact supported mappings | repository inference |
| Agent Runner selector | disabled/configured policy; when configured, exact qualified adapter/version, funding, credentials, billing, native configuration | consumer authority, implicit enablement, or silent fallback |
| Agent Runner adapter | normalized clean-session model/tool loop and provenance | verification, approval, GitHub |
| Command runner | exact host-isolated commands | ambient host access |
| Verification engine | deterministic checks and evidence | patch generation |
| Sandbox verifier | optional approved read-back | production access |
| Evidence engine | axes, hashes, freshness, redaction | unsupported roll-up |
| PR-bundle renderer | reviewable offline handoff | remote write |
| GitHub adapter | short-lived remote branch and tested draft PR | merge |
| Status projector | event-bound consented provider status and provenance | raw evidence or inferred retirement |

Keep components behind small interfaces. Do not add impact, agent, patch,
verification, or GitHub behavior to `internal/source`.

## Initial Architecture Spine

```text
provider-paid campaign scope
-> provider-confirmed Provider Change Contract
-> provider-originated change event
-> Consumer Installation
-> event-specific read authorization
-> TypeScript repository impact inventory
-> no-write migration plan
-> exact plan approval or installed-preauthorization policy evaluation
-> local write/command authorization
-> [Agent Runner network/credential authorization when needed]
-> [model disclosure/network/credential authorization when needed]
-> isolated deterministic or bounded-agent patch candidate
-> deterministic repository and workflow verification
-> patch artifact
-> optional local branch
-> PR bundle
-> short-lived remote branch
-> tested draft PR
-> consumer review and merge
-> consented event-bound provider status
```

## Artifact Ownership

### Provider-Controlled Inputs

- confirmed sunset objective and deadline;
- source/target API or SDK artifacts;
- migration guidance and semantic intent;
- Provider Change Contract and exact provider event;
- sandbox semantics and rollback guidance.

These inputs cannot execute code or grant consumer authority.

### Consumer-Private Artifacts

- Consumer Installation and event-specific authorization;
- repository impact inventory and migration plan;
- Agent Runner network/credential and model disclosure/network/credential
  grants;
- prompts, responses, tool traces, token/cost records;
- workspace, patch, local branch, and verification;
- PR bundle and GitHub result;
- retention/deletion evidence.

### Provider-Visible By Explicit Consent

- campaign-level status fields;
- opaque repository status;
- verification boundary;
- merge/close outcome when consented.

The complete provider-status artifact is consumer-private control state. Only
the serialized output of the provider-payload projector may be transmitted;
it contains exactly the artifact's declared provider-visible fields narrowed
by any current external consent ceiling. Artifact type/version, consent,
field-policy, interpretation, privacy, and integrity metadata do not cross the
boundary merely because they were validated.

Raw source, diffs, prompts, responses, agent sessions, tool traces, logs, and
credentials are never API-provider-visible.

### Storage And Disclosure Boundary

Private product state lives in a configured consumer-controlled root outside
the checkout and public source repository. Provider disclosure and model
disclosure are separate. Both require exact field/request policies.

Committed Factory artifacts contain source-safe control, lifecycle, or
aggregate/hash evidence only. They refer to private evidence by opaque ID and
digest.

### Factory Worker Versus Product Authority

Factory uses only its closed `approval`, `credentials`, and `network`
capabilities for implementation workers. Lumyn product grants remain private
and action-specific. Factory dispatch cannot confer repository, model, branch,
PR, or merge authority.

Task- and campaign-level product-authority arrays describe the complete
capability universe that implementation may need; they are not runtime grants.
Each side effect selects a named route, freezes that route's exact required
plus conditionally selected capability union, and leaves every unselected
capability unauthorized. A campaign composes the validated impact, candidate,
verification, optional sandbox, delivery, and reporting routes per
installation/event/run; it never grants their aggregate union to every
participant.

The active repo-local v3 control set remains authoritative for factoryd and
still records Factory profile/runtime compatibility as unqualified. The landed
external source-profile update does not change that state. Until the controls
are regenerated and the bundle/runtime and exact mission are qualified, the
mission-paused configs are an enforced stop rather than an executable
implementation path.

### Independent Promotion Evidence

When task policy selects `code-review`, `holdout-evaluator`, `trace-grader`, or
`evidence-attestor`, each writes current, task-bound evidence outside the
implementation worker's writable scope. Shipping fails before `commit-push`
when required independent evidence is absent, stale, self-authored, or
non-passing.

Independent repository verification runs in a fresh process and view with
frozen command and verification-configuration digests, no Agent Runner/model
credentials, and no generation-owned evidence handle. Only the
verifier/evidence boundary can persist verification results for the exact
candidate head.

## Evidence Model

Evidence preserves separate axes for:

- provider-confirmed intent;
- impact coverage;
- generation mode and provenance;
- patch scope;
- repository baseline and checks;
- workflow execution;
- Agent Runner/model disclosure, funding, credential/billing ownership, and
  cost;
- delivery state;
- permission state;
- residual risk.

Bind evidence to:

- intent and source/target digests;
- repository base and candidate heads;
- plan digest;
- explicit `agent_execution_policy`, plus deterministic recipe provenance or,
  when configured, exact Agent Runner/model/prompt/tool provenance;
- patch digest;
- command and environment identity;
- workflow/cassette/sandbox identity;
- artifact hashes.

Changing any bound input invalidates dependent evidence.

Canonical verification labels remain `not_run`, `static_verified`,
`repo_verified`, `workflow_contract_replay_passed`,
`workflow_verified_replay`, `workflow_verified_mock`,
`workflow_verified_sandbox`, `partial`, `failed`, `gap`, and `stale`.

Agent completion is never a verification label.

## Structured Parser Boundaries

- OpenAPI, JSON, YAML, manifests, lockfiles, schemas, CI results, and GitHub
  responses use structured parsers or stable APIs.
- Supported TypeScript impact uses an AST or comparably structured parser.
- Text search can seed discovery but cannot prove patchability.
- The consumer selects one canonical package/read root.
- Symlink escape, path traversal, out-of-root references, and ambiguous roots
  fail before analysis.
- External refs remain blocked in deterministic tests.

## Patch Safety Boundary

Every patch mode must:

- start from an immutable approved plan;
- bind current provider-confirmed intent and repository base;
- run in an isolated worktree or equivalent;
- canonicalize and validate real paths;
- reject traversal and symlink escape;
- exclude generated, vendored, cache, and build output unless explicitly
  approved;
- enforce writable paths and file/line/diff budgets;
- map each edit to a plan item;
- produce generation provenance and rollback evidence;
- leave the default branch untouched.

Deterministic mode additionally records the recipe and produces the same patch
for identical pinned inputs.

Agent mode additionally enforces exact qualified adapter/version, executable
source/digest and conformance digest, auth mode/entitlement class, clean
session, funding/credential/billing ownership, native configuration, Agent
Runner/model endpoint, prompt/tool, disclosure, network, credential, and
turn/token/time/retry/cost boundaries. It rejects executable
shadowing and silent fallback and does not claim byte-identical patch
determinism.

## Command Execution Boundary

Repository commands are untrusted code:

- exact allowlist and working directory;
- exact read-only/writable mounts and neutral home/temp;
- explicit executable/toolchain roots;
- timeout and output budgets;
- exact backend/version/configuration/qualification identity and host platform;
- hard CPU-time, memory, PID, process-tree-depth, disk, and open-file quotas;
- network disabled by default;
- lifecycle scripts disabled by default;
- sanitized environment and no ambient secrets;
- no host home, credential stores, OS credentials, agent/Docker/unrelated
  sockets, or extra inherited descriptors;
- child processes inherit every restriction;
- an unenforceable host-isolation backend blocks before launch;
- pre-patch baseline remains separate from post-patch result.

Model tools do not bypass this boundary.

## Live Sandbox Boundary

Optional sandbox verification requires separate request disclosure, network,
and credential grants. It uses non-production data and credentials, exact
endpoint/operations, request/write budgets, idempotency, cleanup, retention, and
orphan evidence. Sandbox proof is independent from local deterministic
verification and required pilot draft-PR delivery.

The sandbox entrypoint runs in its own qualified isolation profile with a
read-only exact-head mount, exact entrypoint/directory, neutral roots,
sanitized environment, only the task-scoped sandbox credential,
endpoint/operation-only egress, inherited child and hard resource limits,
teardown, cleanup, and orphan evidence.

## GitHub Boundary

- Patch and PR-bundle delivery require no GitHub credential.
- Local branch creation remains inside the consumer execution plane.
- The GitHub App installation may persist; its installation token is
  short-lived, least-privilege, issued at runtime, and never stored in the
  Consumer Installation.
- Remote branch write and draft-PR write are separate exact grants.
- A manual fallback cannot close automated-delivery acceptance.
- PRs are draft-only.
- Idempotency binds event, Provider Change Contract, Consumer Installation
  authorization, repository, base, head, plan, and verification-evidence
  digests.
- `EXP-003` requires the composed provider-channel event through
  installation, impact, plan, Lumyn-generated candidate, verification,
  branch, draft-PR, and local status-projection path. Transmission is optional
  for technical delivery. For `PILOT-003`, that same qualifying run must
  contain an organically agent-assisted plan item on a consumer-selected
  qualified runner, pass independent exact-head verification, open the Lumyn
  draft PR, and transmit the bound provider projection. Separate
  agent/delivery runs or deterministic rerouting do not qualify.
- Default-branch writes and auto-merge are prohibited.
- Provider payment never authorizes a GitHub action.

## Systems Thinking Map

| State | Owner | Feedback | Failure signal | Recovery |
|---|---|---|---|---|
| Campaign scope | API Provider + Lumyn Operator | confirmed sunset decision | unclear intent or no accountable buyer | re-scope or stop |
| Repository authority | Consumer | grant/revoke | missing or stale grant | request exact authority |
| Agent Runner authority | Consumer | conformance, session, auth, usage, and cost evidence | adapter/version, native config, credential/billing owner, or fallback drift | stop and reauthorize |
| Model authority | Consumer | request/cost/provenance evidence | provider, model route, disclosure, credential, endpoint, or budget drift | stop and reauthorize |
| Patch candidate | Consumer execution plane | diff and provenance | scope escape or unsupported inference | discard workspace |
| Verification | Consumer execution plane | pinned command evidence | stale head or failed check | repair or report |
| PR bundle | Consumer | review feedback | incomplete evidence or residual risk | regenerate bundle |
| Remote branch/PR | Consumer | GitHub state | stale base, duplicate, or denied grant | refresh or remain offline |
| Merge | Consumer Maintainer | CI/review/outcome | close, revert, or correction | human remediation |

Do not optimize generation quality while hiding consent, activation, review, or
operator cost.

## TDD And Red-First Expectations

- Add a failing unit, integration, scenario, permission, prompt-injection, or
  budget test before behavior when practical.
- Freeze corpus ground truth before scoring.
- Fix expected deterministic patches before implementing recipes.
- For agent mode, freeze outcome constraints and evaluators without exposing
  held-out answers.
- Add false-green verification cases first.
- Use mocks before live model, sandbox, or GitHub access.

## ADR And Decision Triggers

Require an ADR or decision update for:

- authority, execution-plane, or data-sharing changes;
- Agent Runner adapter/version/executable, auth/entitlement, conformance,
  funding, credential/billing, native-configuration, or fallback posture;
- Model Provider, endpoint, credential, disclosure, or budget posture;
- agent tools, session, or isolation;
- public command or schema contracts;
- patch/branch/PR-bundle ownership;
- verification semantics;
- GitHub permissions;
- hosted campaign coordination;
- release/distribution posture;
- major performance, cost, or reliability tradeoffs.

ADR-0004 governs provider-originated API update delivery. ADR-0005 governs
customer-selected Agent Runners. ADR-0007 governs their executable M2 artifact,
authorization, privacy, and compatibility contracts. ADR-0003 governs the
remaining bounded-agent execution and trust substrate. ADR-0002 remains
historical context for the v2 provider-sponsored deterministic-first
rebaseline.

## Performance And Cost Triggers

- Measure impact analysis, generation, verification, and PR-bundle duration
  separately.
- Record repository size, file count, AST memory, command duration, artifact
  size, and GitHub calls.
- Agent mode records adapter/version, conformance, funding mode, credential and
  usage-billing owners, prompt/response size, turns, tokens, retries, wall time,
  Agent Runner/model cost, tool calls, and operator interventions.
- Budget exhaustion or runner failure fails closed; it never silently switches
  adapter, model, credential/billing owner, or scope.
- Provider and consumer labor remain visible beside Lumyn operator time.

## Reliability And Recovery Triggers

Test:

- interrupted workspace and partial patch;
- stale base, intent, plan, Agent Runner conformance, or model policy;
- runner executable integrity/shadowing, authentication/entitlement,
  clean-session, native-configuration, cancellation, malformed/partial output,
  silent-fallback, and credential-persistence failures;
- model timeout, malformed tool call, prompt injection, budget exhaustion, and
  partial response;
- command timeout and flaky/pre-existing tests;
- sandbox drift and cleanup failure;
- PR-bundle regeneration;
- GitHub retry, stale branch, and duplicate PR;
- revoked or expired authority;
- redaction and private-artifact deletion failure.

Retries preserve the same authorization and idempotency identity.

## Trust-Mode Posture

### Public Fixture Planning

- no consumer repository;
- no live model or external credential;
- no GitHub write;
- pinned, licensed, source-safe fixtures only.

### Consumer Repository Read

- exact task-scoped read grant;
- read-only impact;
- no provider visibility.

### Consumer Mutation

- approved plan digest;
- scoped isolated write;
- exact command boundary;
- no network by default.

### Bounded Agent

- consumer-selected exact qualified adapter/version/executable,
  auth/entitlement, and funding route;
- clean ephemeral session and explicit native-configuration posture;
- separate Agent Runner/model disclosure, endpoints, credentials, tools, and
  budgets;
- isolated workspace;
- no ambient authority;
- no silent fallback;
- untrusted output;
- independent deterministic verification.

### Live Sandbox

- optional non-production verification;
- separate disclosure, network, and credential grants;
- cleanup and evidence.

### Draft PR

- short-lived remote branch plus separate draft-PR grant;
- evidence-bound idempotency;
- no merge authority.

### Provider Status Projection

- exact consumer-consented field allowlist and event/evidence binding;
- observed, consumer-reported, and unknown provenance;
- no `not_applicable` or `unaffected` from silence and no `retired` from merge;
- no raw consumer code, diffs, prompts, responses, agent sessions, tool traces,
  logs, or credentials.

## Runtime Shape

```text
Go orchestration core
  -> structured Provider Change Contract and event intake
  -> Consumer Installation and event-authorization engine
  -> TypeScript impact adapter
  -> migration planner
  -> deterministic transformer
  -> [configured Agent Runner selector and common conformance contract]
  -> [pinned Codex or Claude Code adapter]
  -> isolated workspace and command runner
  -> deterministic verification engine
  -> evidence and PR-bundle renderer
  -> optional sandbox adapter
  -> short-lived GitHub adapter
  -> consented status projector
```

Keep Agent Runner and model adapters behind narrow interfaces. No runner
account, native configuration, model endpoint, SDK, or hosted control plane
becomes an implicit dependency.

## Architecture Budget And Decomposition

Source files warn at 1200 lines and fail at 2500 lines under the repository
architecture-budget policy. Decompose change/event intake, installation,
authorization, impact, planning, agent execution, patching, verification,
PR-bundle rendering, GitHub delivery, and status projection rather than
creating a product monolith.
