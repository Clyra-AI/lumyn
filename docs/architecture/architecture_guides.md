# Lumyn Architecture Guide

Status: v3 planning architecture; product runtime not implemented

## Architecture Objective

Lumyn turns provider-confirmed API or SDK sunset intent into a bounded,
evidence-backed migration handoff inside a consumer-authorized repository.
Provider sponsorship never grants access to consumer code, and bounded-agent
generation never grants a model authority over repository scope, verification,
GitHub delivery, or merge.

The architecture optimizes for:

- services-led campaign operation without ambient operator authority;
- consumer-local or consumer-controlled CI execution;
- deterministic transforms for exact mappings;
- bounded-agent generation for approved repository-specific changes;
- deterministic, baseline-aware verification;
- exact model egress, credential, network, disclosure, and cost controls;
- patch and PR-bundle handoff without requiring GitHub access;
- optional separately authorized remote branch and draft PR;
- proof-honest evidence and human merge authority.

This guide defines planned boundaries. The compiled v3 control set is
repo-local planning and validation authority; it does not implement the v3
runtime or authorize a live product action.

factoryd dispatch remains paused. A separate reviewed change must rebaseline
the external Factory `profiles/lumyn.yaml` profile and qualify the factoryd
bundle/runtime against the exact active mission before any task is unpaused.

## Trust And Data Planes

### Provider Campaign Plane

Owns:

- provider identity and accountable campaign operator;
- the paid sunset objective and compatibility window;
- provider-confirmed source/target semantics;
- migration guidance, canary information, and rollback guidance;
- a signed declarative provider packet when one is available and confirmed;
- the invited cohort and commercial campaign decision.

Does not own:

- consumer repository access;
- consumer commands, model credentials, or execution;
- raw source, diffs, logs, traces, prompts, or responses;
- consumer branch, PR, review, or merge authority.

The provider packet is authoritative change data when supplied and confirmed.
It cannot execute code or widen policy. V3 defers elaborate provider PKI,
continuous provider-status resolution, and receipt-backed billing.

The v2 `provider-authenticated consumer receipt-key bindings` and
provider-signed acknowledgement design remains historical context, not an
active v3 dependency.

### Consumer Execution Plane

Owns:

- repository and selected package root;
- read/write path scopes;
- migration-plan approval;
- isolated workspace and local branch;
- command allowlist and host-isolation policy;
- model disclosure, endpoint, credential, tool, and budget policy;
- dependency and package-registry policy;
- private impact, plan, prompt, response, patch, verification, and PR-bundle
  evidence;
- optional remote branch and draft-PR authorization;
- review and merge.

Consumer-private runtime state lives outside the checkout and public source
repository with explicit retention, deletion, and evidence ownership.

### Bounded Model Plane

The model plane is active v3 planned scope and a separate trust boundary.

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
  funds campaign and confirms sunset intent
        |
        v
Provider Campaign Plane
  supplies guidance and optional signed declarative packet
        |
        | no repository authority
        v
Consumer Maintainer
  authorizes read-only impact
  -> approves no-write migration plan
  -> authorizes exact local write/command/model boundaries
        |
        v
Consumer Execution Plane
  deterministic transform or bounded-agent generation
  -> deterministic verification
  -> patch + optional local branch + PR bundle
        |
        +--> optional remote branch grant
        +--> optional draft-PR grant
        |
        v
Consumer review and merge
```

Lumyn Operator assistance is an operating role, not an authority plane.

## Product State Machines

### Provider Intent

```text
draft
-> provider_confirmed
-> selected_for_campaign
-> superseded | withdrawn | completed
```

If a signed declarative packet is used:

```text
draft
-> signed
-> confirmed
-> selected_for_campaign
-> superseded | withdrawn | expired
```

Stale, conflicting, unconfirmed, or executable intent blocks planning or
mutation.

### Consumer Authorization

```text
not_requested
-> read_requested
-> read_authorized
-> plan_ready
-> plan_approved
-> local_execution_authorized
-> [model_authorized]
-> [remote_branch_authorized]
-> [draft_pr_authorized]
-> expired | revoked
```

Each transition is action-specific. Approval at one state does not imply the
next.

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
-> [draft_pr_open]
-> merged | closed | reverted
```

Generation mode and verification strength remain independent.

## Component Boundaries

| Component | Responsibility | Must not own |
|---|---|---|
| CLI | parsing, config, JSON envelope, exits | product inference |
| Campaign intake | provider-paid scope and confirmed intent | repository authority |
| Source intake | pinned provider materials and SDK refs | consumer writes |
| Intent normalizer | typed changes and unresolved questions | arbitrary execution |
| Authorization engine | exact consumer grants, expiry, revocation | side-effect execution |
| TypeScript analyzer | package, imports, call sites, exclusions | file mutation |
| Impact engine | applicability and coverage | patch application |
| Migration planner | complete no-write route and budgets | approval or writes |
| Workspace manager | isolated workspace, safe paths, local branch | semantic decisions |
| Deterministic transformer | exact supported mappings | repository inference |
| Agent runner | bounded model/tool loop and provenance | verification, approval, GitHub |
| Command runner | exact host-isolated commands | ambient host access |
| Verification engine | deterministic checks and evidence | patch generation |
| Sandbox verifier | optional approved read-back | production access |
| Evidence engine | axes, hashes, freshness, redaction | unsupported roll-up |
| PR-bundle renderer | reviewable offline handoff | remote write |
| GitHub adapter | optional remote branch and draft PR | merge |
| Attestation exporter | optional consented provider status | raw consumer evidence |

Keep components behind small interfaces. Do not add impact, agent, patch,
verification, or GitHub behavior to `internal/source`.

## Initial Architecture Spine

```text
provider-paid campaign scope
-> provider-confirmed source/target intent
-> [optional signed declarative provider packet]
-> consumer read authorization
-> TypeScript impact report
-> no-write migration plan
-> plan approval
-> local write/command authorization
-> [model disclosure/network/credential authorization when needed]
-> isolated deterministic or bounded-agent patch candidate
-> deterministic repository and workflow verification
-> patch artifact
-> optional local branch
-> PR bundle
-> [optional remote branch]
-> [optional draft PR]
-> consumer review and merge
-> [optional consented provider status]
```

## Artifact Ownership

### Provider-Controlled Inputs

- confirmed sunset objective and deadline;
- source/target API or SDK artifacts;
- migration guidance and semantic intent;
- optional signed declarative packet;
- sandbox semantics and rollback guidance.

These inputs cannot execute code or grant consumer authority.

### Consumer-Private Artifacts

- repository authorization;
- impact and migration plan;
- model disclosure/network/credential grants;
- prompts, responses, tool traces, token/cost records;
- workspace, patch, local branch, and verification;
- PR bundle and optional GitHub result;
- retention/deletion evidence.

### Provider-Visible By Explicit Consent

- campaign-level status fields;
- opaque repository status;
- verification boundary;
- merge/close outcome when consented.

Raw source, diffs, prompts, responses, logs, traces, and credentials are private
by default.

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

The active repo-local v3 control set does not by itself qualify factoryd.
Until the external profile and factoryd runtime are requalified, the
mission-paused configs are an enforced stop rather than an executable
implementation path.

### Independent Promotion Evidence

When task policy selects `code-review`, `holdout-evaluator`, `trace-grader`, or
`evidence-attestor`, each writes current, task-bound evidence outside the
implementation worker's writable scope. Shipping fails before `commit-push`
when required independent evidence is absent, stale, self-authored, or
non-passing.

## Evidence Model

Evidence preserves separate axes for:

- provider-confirmed intent;
- impact coverage;
- generation mode and provenance;
- patch scope;
- repository baseline and checks;
- workflow execution;
- model disclosure and cost;
- delivery state;
- permission state;
- residual risk.

Bind evidence to:

- intent and source/target digests;
- repository base and candidate heads;
- plan digest;
- deterministic recipe or model/prompt/tool provenance;
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

Agent mode additionally enforces exact model, endpoint, prompt/tool, disclosure,
network, credential, turn/token/time/retry/cost, and tool-call boundaries. It
does not claim byte-identical patch determinism.

## Command Execution Boundary

Repository commands are untrusted code:

- exact allowlist and working directory;
- exact read-only/writable mounts and neutral home/temp;
- explicit executable/toolchain roots;
- timeout and output budgets;
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
verification and optional draft-PR delivery.

## GitHub Boundary

- Patch and PR-bundle delivery require no GitHub credential.
- Local branch creation remains inside the consumer execution plane.
- Remote branch write and draft-PR write are separate grants.
- PRs are draft-only.
- Idempotency binds repository, base, head, and plan/evidence digests.
- Default-branch writes and auto-merge are prohibited.
- Provider payment never authorizes a GitHub action.

## Systems Thinking Map

| State | Owner | Feedback | Failure signal | Recovery |
|---|---|---|---|---|
| Campaign scope | API Provider + Lumyn Operator | confirmed sunset decision | unclear intent or no accountable buyer | re-scope or stop |
| Repository authority | Consumer | grant/revoke | missing or stale grant | request exact authority |
| Model authority | Consumer | request/cost/provenance evidence | disclosure, credential, endpoint, or budget drift | stop and reauthorize |
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
- model provider, endpoint, credential, disclosure, or budget posture;
- agent tools or isolation;
- public command or schema contracts;
- patch/branch/PR-bundle ownership;
- verification semantics;
- GitHub permissions;
- hosted campaign coordination;
- release/distribution posture;
- major performance, cost, or reliability tradeoffs.

ADR-0003 governs services-led bounded-agent execution. ADR-0002 remains
historical context for the v2 provider-sponsored deterministic-first rebaseline.

## Performance And Cost Triggers

- Measure impact analysis, generation, verification, and PR-bundle duration
  separately.
- Record repository size, file count, AST memory, command duration, artifact
  size, and GitHub calls.
- Agent mode records prompt/response size, turns, tokens, retries, wall time,
  model cost, tool calls, and operator interventions.
- Budget exhaustion fails closed; it never silently switches model or widens
  scope.
- Provider and consumer labor remain visible beside Lumyn operator time.

## Reliability And Recovery Triggers

Test:

- interrupted workspace and partial patch;
- stale base, intent, plan, or model policy;
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

- exact disclosure, endpoint, credential, tools, and budgets;
- isolated workspace;
- no ambient authority;
- untrusted output;
- independent deterministic verification.

### Live Sandbox

- optional non-production verification;
- separate disclosure, network, and credential grants;
- cleanup and evidence.

### Draft PR

- optional remote branch plus separate draft-PR grant;
- evidence-bound idempotency;
- no merge authority.

### Provider Attestation

- optional exact field allowlist;
- no raw consumer evidence by default.

## Runtime Shape

```text
Go orchestration core
  -> structured provider-intent intake
  -> TypeScript impact adapter
  -> migration planner
  -> deterministic transformer
  -> bounded-agent adapter
  -> isolated workspace and command runner
  -> deterministic verification engine
  -> evidence and PR-bundle renderer
  -> optional sandbox adapter
  -> optional GitHub adapter
```

Keep model adapters behind a narrow interface. No model endpoint, SDK, or
hosted control plane becomes an implicit dependency.

## Architecture Budget And Decomposition

Source files warn at 1200 lines and fail at 2500 lines under the repository
architecture-budget policy. Decompose campaign intake, authorization, impact,
planning, agent execution, patching, verification, PR-bundle rendering, and
GitHub delivery rather than creating a product monolith.
