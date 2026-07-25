# Lumyn MVP — Provider-Originated API Update Delivery

| Field | Value |
|---|---|
| Version | 3.1 |
| Status | Active v3.1 planning contract; no runtime implementation authorized by this document alone |
| Owner | Product and Engineering |
| Last Updated | 2026-07-25 |
| Primary Audience | Lumyn builders, API-provider design partners, API-consumer maintainers, and technical investors |
| Source Task Plan | `docs/product/plan.md` |
| Compiled Factory Plan | Repo-local v3.1 set regenerated; factoryd dispatch paused pending external Factory profile and runtime alignment |

---

## Purpose

This document is the product source of truth for the Lumyn MVP.

Lumyn is the provider-to-consumer application layer for consequential API and
SDK changes. An API Provider publishes one confirmed change; each authorized
consumer installation determines whether its repository is affected and,
within its selected action mode, prepares and verifies the repository-specific
update or opens a tested draft PR. The initial go-to-market motion is a
services-assisted, provider-paid deprecation or sunset campaign. Services
accelerate onboarding and learning; they are not the product identity.

The product promise is:

> When an API Provider publishes an important change, Lumyn carries that
> authoritative intent into each consenting customer's repository, determines
> the affected usages, produces and verifies the smallest justified update,
> and opens a tested draft PR without taking review or merge authority away
> from the customer.

Lumyn is not differentiated by owning a better coding model or merely opening
a pull request. Generic coding agents can already propose code changes. Lumyn
is differentiated by the installed, trusted delivery loop around that
replaceable engine: provider-originated change identity, reusable confirmed
semantics, consumer-local repository understanding, bounded authorization,
independent verification, draft-PR delivery, and consented cohort progress
across otherwise disconnected organizations.

This version refines v3.0 without reviving v2.0's cryptographic activation,
connection-receipt billing, or deterministic-only engine. Draft-PR delivery is
now required product proof because it tests the application layer named in the
YC thesis; manual bundles remain a safe fallback but do not prove that loop.
Historical artifacts remain historical evidence and must not be rewritten to
claim v3.1 behavior.

---

## Executive Summary

API providers publish changelogs, migration guides, SDK releases, OpenAPI
descriptions, examples, and occasionally codemods or agent instructions. An
API consumer still has to determine where its integration lives, understand
local wrappers and data mappings, change the code, repair repository-specific
failures, run the right checks, and decide whether the integration still
performs its business job.

Lumyn closes that loop:

```text
provider-originated change event
-> reusable Provider Change Contract
-> authorized consumer installation
-> consumer-local repository impact inventory
-> reviewable migration plan
-> deterministic transforms + replaceable bounded coding agent
-> independent baseline-aware compile/test/mock/sandbox verification
-> optional separately authorized agent-assisted repair and re-verification
-> verified candidate
-> tested draft PR
-> consented provider rollout status
```

The operating model is provider-funded and consumer-controlled:

- The API Provider is the main economic buyer and campaign sponsor.
- The API Consumer Organization owns repository, execution, model-data,
  credential, disclosure, review, and merge authority.
- Provider payment never grants access to consumer source code.
- Agent execution is disabled unless configured. When needed, the coding agent
  is a customer-selected, qualified adapter and runs through a consumer-local
  or explicitly consumer-approved execution boundary.
- Codex and Claude Code are the launch Agent Runner targets after each passes
  the same conformance gate. Cursor remains deferred behind that gate.
- For configured execution, the default funding route uses the consumer's
  qualifying account, enterprise subscription, API credential, or local
  runtime. A provider-sponsored, Lumyn-managed usage route is optional and
  never gives the API Provider code or agent access.
- The API Provider never receives raw code, diffs, prompts, responses, agent
  sessions, tool traces, logs, or credentials.
- Lumyn never writes to the default branch or auto-merges in the MVP.

The commercial outcome is an affected customer cohort moving safely off the
targeted version, not code generation in isolation. The product proof begins
with Lumyn-opened tested draft PRs and ends with consumer-controlled adoption
and provider-visible, consented rollout evidence. For an API Provider, value
is a shorter legacy-version tail, fewer support hours, lower migration risk,
and actionable cohort readiness. For an API Consumer Organization, value is
less discovery and migration work and a smaller, better-evidenced change to
review.

---

## Product Thesis

Dependency bots update versions. Vendor codemods rewrite known syntax. Generic
coding agents can inspect a repository and attempt broader changes. Provider
SDK and documentation tools explain what changed. None of those mechanisms
alone owns the cross-company migration outcome:

- provider evidence may not describe consumer abstractions;
- a generic agent may not know which change semantics are authoritative;
- syntactically correct changes may be behaviorally wrong;
- weak or missing tests can create false confidence;
- the provider cannot inspect private consumer code;
- the consumer may not trust a provider-controlled mutation path;
- opening a PR does not prove adoption or legacy-version retirement.

Lumyn's initial differentiation is:

```text
provider-originated distribution
+ reusable authoritative change intent
+ installed consumer authorization
+ consumer-local repository impact inventory
+ deterministic and bounded-agent routing
+ baseline-aware verification and repair
+ tested draft-PR delivery
+ consented cohort progress
```

Every real migration must be compared with the status-quo baseline:

> official migration guide, vendor skill, or codemod plus the same selected
> qualified runner—Codex or Claude Code at launch.

If Lumyn does not materially reduce maintainer effort or materially improve
verified completion without increasing correction risk, it has not
demonstrated product value.

The intended durable position is the trusted update-delivery channel between
API Providers and their customers' codebases, not code transformation alone.

---

## Roles And Terminology

Use these terms consistently. Avoid the bare word `provider` where it could
mean the API Provider, Agent Runner Vendor, or Model Provider.

### API Provider

The company that owns and sells the API or official SDK. It supplies or
confirms migration semantics, recruits the sponsored cohort, and pays for the
initial provider campaign.

### API Consumer Organization

The organization whose application depends on the API Provider. It owns the
repository, tests, credentials, execution environment, model-data policy, and
integration risk.

### Consumer Maintainer

The engineer authorized by the API Consumer Organization to approve analysis,
review the migration plan, authorize execution, inspect the result, and merge
or reject it.

### Provider Operator

The provider-side DX, SDK, API-platform, customer-engineering, or
solutions-engineering person who supplies change context and recruits
participants.

### Lumyn Operator

A member of the Lumyn team who prepares Provider Change Contracts, assists
onboarding, operates the first campaigns, or supports a consumer upgrade
sprint. Lumyn Operator time is measured as delivery COGS, separate from
product development.

### Model Provider

The company or local runtime that supplies the model used by the coding agent.
It is never referred to as the API Provider. Cloud-model data handling is a
separate consumer decision and disclosure boundary.

### Agent Runner Vendor

The company that supplies the selected coding-agent harness. It may also be the
Model Provider, or it may broker a separately identified downstream Model
Provider. Lumyn records both roles and their data-processing boundaries instead
of treating “agent” and “model” as interchangeable.

### Agent Runner

The exact consumer-selected process or harness that executes Lumyn-composed
context and tool requests, invokes the Model Provider or local model, returns
candidate edits and normalized lifecycle events, and records attempt
provenance. Lumyn's orchestrator and host-isolation boundary select context,
expose tools, and enforce paths, credentials, network, disclosure, and budgets.
The Agent Runner operates as a clean, ephemeral session inside that envelope
and cannot grant itself broader authority. Lumyn never resumes a personal or
unrelated agent conversation.

The first adapter targets are Codex and Claude Code. Each exact adapter version
and executable digest from an approved source must pass one common conformance
contract and an approved live canary before Lumyn advertises it as supported.
The selected auth mode and entitlement class must permit the intended
non-interactive use; Lumyn does not assume a personal subscription is valid for
consumer CI. Cursor is a later candidate behind the same gate, not a launch
commitment. Lumyn never silently falls back between adapters, versions, models,
endpoints, credential owners, or usage-billing owners.

Each Consumer Installation sets `agent_execution_policy` to `disabled` or
`configured`. `disabled` is the least-privilege default and is valid for
`notify_only`, `scan_only`, and deterministic-only mutation. It grants no
runner or model authority. If a routed plan later needs `agent_assisted`,
Lumyn pauses until the Consumer Maintainer explicitly configures and authorizes
an exact qualified route; it never upgrades the policy implicitly.

When agent execution is `configured`, the consumer chooses one of two funding
and credential modes:

- `consumer_managed` is the default configured mode. The API Consumer
  Organization owns and authorizes the agent account, qualifying enterprise
  subscription, API credential, or local runtime and owns third-party usage
  billing. The route must expose the actual Model Provider, endpoint, model,
  and version and permit non-interactive automation; an opaque or changing
  subscription route does not qualify. Lumyn receives no reusable credential.
- `provider_sponsored_lumyn_managed` is optional. The API Provider funds the
  campaign, while Lumyn owns the approved agent/model usage billing and injects
  a task-scoped brokered credential only into the same consumer-authorized
  local or CI boundary. The broker binds issuer; installation, event, plan,
  attempt, runner, and model audience; and maximum one-hour TTL. One-time
  redemption creates one attempt-scoped session. Multiple in-attempt calls are
  allowed only within hard token/cost quotas; refresh, post-attempt replay, and
  cross-attempt reuse are forbidden. Require revocation and reconciliation. A
  budget-enforcing proxy is required when the vendor cannot issue those
  bounds; otherwise this mode is unavailable. The API Provider receives no
  credential, repository context, prompt, response, or agent-session access.

Every configured agent action also selects exactly one authorization topology:

- `local_runtime` permits no external runner or model egress;
- `runner_mediated` requires Agent Runner network and credential plus model
  request-disclosure authorization;
- `direct_model` requires direct model network, credential, and
  request-disclosure authorization; or
- `hybrid` requires both runner and direct-model minimum sets.

Package-registry read remains a separate plan-selected capability. Lumyn
freezes the topology and its minimum scope set before launch; a route cannot
omit a required capability or gain an unselected topology's authority.

### Provider Change Contract

A versioned, digest-bound, non-executable artifact assembled from official
documentation, OpenAPI descriptions, SDK changes, examples, migration
guidance, provider clarification, and Lumyn analysis. It tells Lumyn what
changed once, without prescribing repository-specific code. It contains a
stable change identity, source and target versions, intended audience,
applicability, typed semantics, known mappings, unsupported or ambiguous
conditions, verification guidance, provenance, provider confirmation, and
supersession or withdrawal state.

A Provider Change Contract is classified as:

- `public_derived`: derived from pinned public evidence without provider
  endorsement;
- `provider_confirmed`: reviewed and confirmed by an accountable Provider
  Operator.

The internal artifact identifier remains `migration-pack` during the v3
compatibility window. A Provider Change Contract is untrusted data, cannot
execute arbitrary code or grant repository authority, and is not proof that a
candidate is correct.

### Repository Impact Inventory

A consumer-local model, represented internally as an `integration-graph`, of
where the selected API or SDK appears in the
repository, including dependency state, imports, aliases, wrappers, adapters,
call sites, request and response mappings, relevant configuration, tests,
mocks, cassettes, fixtures, and explicit uncertainty.

### Provider Change Event

A versioned provider-originated envelope that references one exact Provider
Change Contract and declares issuer, API or SDK, audience, deadline, severity,
and supersession or withdrawal state. The event contains no executable
instructions. The first transport is a versioned JSON manifest at an exact
provider-controlled HTTPS URL pinned by the Consumer Installation. The
manifest either embeds the Provider Change Contract or names its exact
provider-controlled HTTPS URL; the retrieved bytes must match the declared
contract digest. The manifest's monotonic sequence, issued/expiry times, and
detached signature are verified against one campaign key enrolled during setup.
Provider Operators publish the manifest; Lumyn may assist setup but cannot
publish as the provider. An attended file import of the same envelope is a
recovery path, not proof of provider-channel delivery. Duplicate, replayed,
stale, conflicting, withdrawn, wrong-audience, expired, or unauthenticated
events fail closed.

### Consumer Installation

A durable, consumer-owned local policy binding an API Provider or channel to
an exact repository and package root, version or audience selectors, allowed
actions, paths, commands, model-data posture, GitHub scopes, disclosure fields,
expiry, revocation, and authorization mode. Action modes are ceilings:
`notify_only`, `scan_only`, `prepare_patch`, or `open_draft_pr`. The consumer
chooses either `per_event_approval` or narrowly bounded
`installed_preauthorization`. The latter permits an authenticated in-policy
event to proceed without a new human approval, but only through the installed
paths, commands, model policy, budgets, verification, and action ceiling.
Each update run freezes an immutable event-specific authorization derived from
the installation; the event may narrow that authority but can never widen it.
An installation stores no GitHub token. An approved local or CI credential
broker issues the short-lived token only when a qualifying event reaches its
remote-delivery step.

### Migration Candidate

A proposed set of code and dependency changes. A candidate may be
`deterministic`, `agent_assisted`, or `manual`. It is not a verified outcome
until the verification ladder passes at the declared evidence level.

### Verified Migration Outcome

A migration candidate plus evidence bound to the exact repository base and
candidate head, with explicit coverage, commands, results, residual risk, and
proof boundary. `Verified` never means production correctness unless
production evidence exists; production evidence is outside the MVP.

### Migration Campaign

One Provider Change Contract and event, one API Provider, a bounded cohort of
consenting API Consumer Organizations, and measured update outcomes. Lumyn
operates the first campaigns as services-assisted product onboarding; it is
not a two-sided marketplace or cryptographic billing network.

### Consumer Upgrade Sprint

A consumer-paid engagement for one consequential migration across one to
three repositories. Before M6, it produces paid workflow/problem evidence and
revenue but is not Lumyn adapter/runtime proof. After M6 qualification, it may
prove the execution engine. It never proves provider willingness to pay,
provider-led distribution, or recurring campaign demand.

### Eligible Consumer Unit

The canonical campaign activation and funnel unit: one distinct API Consumer
Organization plus one designated primary eligible repository and one
accountable Consumer Maintainer. Additional repositories from the same
organization are tracked separately and do not increase the five-unit
activation denominator.

### Verified

`Verified` must always name its evidence boundary. The canonical labels remain:

- `static_verified`
- `repo_verified`
- `workflow_contract_replay_passed`
- `workflow_verified_replay`
- `workflow_verified_mock`
- `workflow_verified_sandbox`

A `workflow_verified_*` label requires an approved entrypoint executed from
the exact candidate head plus observed interaction and outcome evidence in the
named environment. Independent contract or cassette replay cannot exceed
`repo_verified`.

---

## Jobs To Be Done

### Primary API Provider Job

> When we must retire an API or SDK version by a fixed deadline, help
> us publish the change once, deliver tested fixes into consenting customers'
> normal review workflows without inspecting proprietary code, and show us
> which affected integrations are progressing toward retirement.

### API Consumer Job

> When an API dependency changes, use only the authority I installed, show
> exactly how my repository is affected, use the coding agent and credential
> route my organization already approves when agent help is needed, and open
> the smallest tested draft PR I can safely review, change, and merge.

Both roles are required for a provider-sponsored campaign. They are not both
economic buyers. A consumer-paid sprint may operate from public and official
evidence without making the API Provider a Lumyn customer.

---

## Initial Segment

The primary design-partner segment is an API-first B2B provider with:

- a consequential deprecation, sunset, or migration deadline in the next 90 to 180
  days;
- roughly 20 to 500 identifiable managed customer integrations;
- at least five reachable Eligible Consumer Units across five distinct API
  Consumer Organizations;
- material support cost, compatibility cost, revenue risk, or retirement risk;
- an accountable executive buyer and named Provider Operator;
- an official TypeScript/Node npm SDK used from a REST API integration;
- versioned docs, OpenAPI descriptions, SDK releases, migration guidance, or
  equivalent evidence;
- repository-level compile, typecheck, or test signals for participating
  consumers;
- preferably a non-production sandbox, reliable mock, or read-back signal.

The first campaign is a named, provider-managed cohort. It is not anonymous
self-service.

Avoid initially:

- providers without a hard deadline or meaningful legacy cost;
- providers unable to identify or recruit five distinct API Consumer
  Organizations with one designated repository each;
- migrations dominated by production-only state or irreversible actions;
- integrations requiring broad production credentials;
- providers unable to clarify ambiguous change semantics;
- low-code or generated-client integrations with little application code;
- migrations dominated by auth redesign, webhook semantics, GraphQL, gRPC, or
  new business values;
- consumer repositories with no usable compile, typecheck, test, mock, replay,
  or sandbox signal.

The secondary consumer-paid segment is an organization with an urgent,
consequential migration, one to three TypeScript/Node repositories, useful
verification signals, and an engineering buyer willing to fund a bounded
upgrade sprint.

### Development Entry Boundary

Lumyn does not need a contracted API Provider to begin the foundational
engineering work. Pinned public documentation, OpenAPI descriptions, SDK
releases, migration guides, license-compatible historical examples, and
synthetic fixtures are sufficient to develop:

- Provider Change Contract intake and provenance;
- provider-event and consumer-installation contracts;
- semantic change normalization;
- a TypeScript repository impact inventory;
- deterministic transforms;
- bounded-agent contracts plus a deterministic fake and status-quo baseline
  harness;
- verification and evidence contracts;
- status-quo baseline experiments.

Public inputs do not establish provider endorsement, prerelease authority,
demand, a reachable cohort, or repository-specific value. Lumyn may not claim
customer-specific value until it runs in at least one real, consenting
consumer repository.

Public fixtures may prove a complete local walking skeleton before a provider
signs only through the deterministic Agent Runner fake. A live generic-agent
public-fixture canary follows M6 implementation and qualification. Live cohort
onboarding, provider distribution, consumer repository access, automated
GitHub delivery, and provider status projection remain gated by a qualified
paid campaign and the relevant consumer authority. A consumer-paid sprint may
produce paid workflow/problem evidence before M6 or engine proof after a
qualified M6 run, but it does not close the provider commercial gate.

### Buyer And Champion

- Primary economic buyer: GM or VP accountable for the API business, VP/Head
  of Platform or Engineering, or CTO at a smaller provider.
- Primary champion/operator: Head of DX or SDK, customer engineering, or
  solutions engineering.
- Required consumer authority: Consumer Maintainer with repository and merge
  authority.
- Consumer-sprint buyer: CTO, VP Engineering, Head of Platform, or engineering
  leader accountable for the urgent migration.

---

## Commercial Model

The initial offers are pricing hypotheses to test, not established list
prices:

| Offer | Buyer | Scope | Price hypothesis |
|---|---|---|---|
| Sunset-readiness sprint | API Provider | Provider Change Contract preparation, cohort qualification, baseline, and campaign plan | `$7.5k–$15k` |
| Provider migration campaign | API Provider | One bounded five-unit campaign attempt; success floor is three valid scans and two tested outcomes, subject to consumer consent | `$25k–$50k` |
| Additional completed repository | API Provider | One repository beyond the included five-unit cohort reaching its contracted verified-delivery event | `$2k–$5k` |
| Urgent consumer upgrade sprint | API Consumer Organization | One migration across one to three repositories | `$10k–$25k` |

The main company thesis is the provider-paid migration campaign. A pre-M6
consumer-paid sprint is secondary revenue and paid workflow/problem evidence;
only a qualified post-M6 run can prove the Lumyn execution engine.

The sunset-readiness sprint is the fastest paid entry point. It may be sold
before the full campaign and credited toward a later campaign invoice, but it
proves only budget, urgency, change definition, and cohort access. Readiness
funds count toward `DISC-001` only after a signed campaign conversion allocates
them as non-refundable campaign consideration and total cleared campaign funds
reach at least `$25,000`; the sprint alone never closes the gate or validates
automated provider-to-consumer delivery.

The default billable event for an additional repository is a Lumyn-generated
candidate that passes independent verification and is opened as the
consumer-authorized tested draft PR. An order form may instead name a verified
local bundle when GitHub delivery is declined; that fallback may be billable
but never counts as automated-delivery or MVP proof.

For the first campaigns:

- Lumyn prepares the Provider Change Contract and signed event manifest with
  provider input; the Provider Operator publishes it at the pinned
  provider-controlled HTTPS URL.
- The API Provider recruits the cohort.
- Each API Consumer Organization installs the provider channel and allowed
  actions locally or in consumer-controlled CI.
- The mandatory first-campaign product proof is one same run from authenticated
  provider event and installed preauthorization through an organically
  agent-assisted item on a qualified consumer-selected Agent Runner,
  independent exact-head verification, a short-lived tested Lumyn-opened draft
  PR, and the bound consumer-consented provider status projection. Separate
  agent, delivery, reporting, or manual-bundle evidence does not qualify.
- For that run the consumer chooses `consumer_managed` access or an approved
  `provider_sponsored_lumyn_managed` route. Credential owner,
  usage-billing owner, and actual Agent Runner/model cost remain explicit.
  Agent execution remains disabled for other work where it is unnecessary.
- The qualifying status projection is machine-bound to that run. Additional
  consumer-consented cohort rollups may be manually prepared, but cannot
  substitute for the qualifying projection.
- API Consumer Organizations pay no Lumyn campaign seat fee. In
  `consumer_managed` mode they remain responsible for their own existing agent
  subscription or metered usage; in `provider_sponsored_lumyn_managed` mode
  that approved usage is Lumyn campaign COGS.
- No annual connected-repository contract is assumed.

Annual platform pricing, hosted coordination, campaign subscriptions, or
connected-repository billing may be tested only after repeatable evidence from
a second paid campaign or an executed annual purchase order. A successful
one-off service engagement does not prove recurring SaaS demand.

---

## System Under Test

The system under test is the complete application-layer handoff:

```text
provider change event and contract
+ consumer installation policy
+ consumer-local repository impact inventory
+ generated candidate
+ repository verification signals
+ draft-PR delivery
+ consented provider status projection
```

The coding agent is an implementation mechanism, but its context selection,
tool use, attempts, cost, output, and repair behavior are part of the evaluated
system.

Every result preserves separate evidence about:

- what the Provider Change Contract says changed;
- which event and installed consumer policy authorized the run;
- what Lumyn found and did not find;
- what route each change used;
- what files and dependencies changed;
- the explicit `agent_execution_policy` and, when configured, which Agent
  Runner adapter/version, execution-funding mode, credential and usage-billing
  owners, actual model route, tools, and budgets were used;
- what baseline and post-change commands ran;
- what behavior was observed;
- what remains unsupported, ambiguous, or unverified;
- what delivery and provider-reporting actions occurred;
- whether status is observed, consumer-reported, inferred, or unknown.

---

## Product Principles

### 1. Two Principals, Two Authorities

The API Provider is authoritative about intended API or SDK semantics. The API
Consumer Organization is authoritative about repository, execution,
model-data, disclosure, and merge. Neither authority implies the other.

### 2. Read Before Write

Impact analysis and planning are read-only and always precede mutation. Under
`per_event_approval`, the Consumer Maintainer approves the exact plan before
any write. Under `installed_preauthorization`, the plan is evaluated
automatically against the installed policy and attached to the resulting
candidate or draft PR; any mismatch pauses for explicit approval.

### 3. Provider Evidence Is Declarative

Provider Change Contracts and events may contain mappings, constraints,
examples, and verification references. They may not execute arbitrary
provider-supplied scripts or widen a Consumer Installation.

### 4. Deterministic Controls, Hybrid Generation

Known safe changes use deterministic transforms. Repository-specific
adaptation may use a bounded coding agent. Missing semantics, required business
input, unsafe access, or inadequate proof is blocked.

Determinism governs pinned inputs, codemods, scope, budgets, commands,
verification, evidence, and status. It does not require byte-identical
agent-generated source.

### 5. Agent Output Is An Untrusted Candidate

The agent cannot declare its own result verified. The same independent
verification ladder applies to deterministic and agent-assisted candidates.
When `agent_execution_policy` is `configured`, the Consumer Maintainer selects
one qualified exact adapter/version and funding route; each attempt starts a
clean session, and Lumyn never silently falls back to another adapter, model,
credential owner, or billing owner. A disabled policy cannot be widened by a
provider event or routed plan.

### 6. Proof Is Multidimensional

Impact coverage, generation provenance, repository validation, workflow
evidence, boundaries, cleanup, and residual risk remain separate axes. No
single green status hides a weaker axis.

### 7. Consumer-Controlled Execution And Disclosure

Repository analysis, candidate generation, commands, credentials, and raw
evidence run inside a consumer-local or explicitly consumer-approved boundary.
Any cloud-model context transfer is disclosed separately and never implies
disclosure to the API Provider.

### 8. Fail Closed

Lumyn may retain an explicitly unverified candidate for diagnosis. It may not
label or deliver it as verified when scope, authorization, redaction, baseline,
commands, evidence, or outcome checks are incomplete or failed.

### 9. Human Adoption Authority

Lumyn never writes to the default branch or auto-merges. The Consumer
Maintainer retains normal review, branch protection, CI, and merge control.

---

## MVP Product Flow

### 1. Qualify The Change And Distribution

Lumyn and the Provider Operator confirm the deadline, source and target
versions, expected cohort, provider-led onboarding commitment, evidence
sources, buyer, operator, verification signals, and campaign price. The first
provider campaign does not proceed without the commercial gates in this PRD.

### 2. Publish The Change Event And Contract

`lumyn pack` normalizes pinned docs, OpenAPI or SDK diffs, examples, and
provider clarification. Conflicts and ambiguity remain visible. A
`public_derived` contract may drive engineering fixtures or a consumer sprint;
a provider-funded campaign requires accountable provider confirmation before
a real consumer mutation. The Provider Operator then publishes one versioned
signed JSON manifest at the campaign's pinned provider-controlled HTTPS URL,
embedding the Provider Change Contract or its exact provider-controlled HTTPS
URL and referencing the exact retrieved-byte digest, audience, sequence,
deadline, and supersession state. Reuse across the cohort is allowed only while
the event and contract remain current. An attended import can recover a
campaign but does not prove the provider-originated channel.

### 3. Install The Consumer-Local Boundary

`lumyn install` records the Consumer Maintainer's durable provider-channel
policy. The Consumer Maintainer selects:

- provider or channel, accepted audience, and version selectors;
- repository and package root;
- readable and writable paths;
- approved commands;
- dependency-install policy;
- network and credential posture;
- `agent_execution_policy`: `disabled` or `configured`;
- when agent execution is configured, the exact qualified Agent Runner
  adapter/version/executable, execution-funding mode, credential owner,
  usage-billing owner, native agent-configuration policy, Agent Runner Vendor,
  actual Model Provider/model route, data-egress and context boundary, and
  agent token, cost, time, attempt, file, and diff budgets;
- verification environments;
- allowed action mode: `notify_only`, `scan_only`, `prepare_patch`, or
  `open_draft_pr`;
- authorization mode: `per_event_approval` or
  `installed_preauthorization`;
- short-lived remote-branch and draft-PR token-issuance policy when selected,
  never a stored token;
- output form and provider-visible status fields.

Each `lumyn update --event` run freezes an event-specific authorization
snapshot derived from the installation. An event may narrow installed
authority but cannot widen it. `open_draft_pr` is an action ceiling, not an
ambient grant: a run proceeds only when its event, plan, commands, paths,
model-data posture, budgets, verification result, and runtime token issuance
all satisfy the selected authorization mode. Anything outside policy pauses
for explicit approval. No selected scope implies another. Production
credentials and production mutation are prohibited.

`lumyn check` performs a non-mutating onboarding preflight. For a configured
agent route it resolves the selected Codex or Claude Code executable by
canonical path, verifies approved source/version/digest and current
conformance, confirms auth mode and entitlement without exposing secrets, and
proves that actual model-route identity and non-interactive use are available.
It performs no model call. A disabled agent policy passes without runner
credentials; a later `agent_assisted` plan pauses for configuration.

### 4. Build The Repository Impact Inventory

`lumyn impact` runs read-only and detects:

- installed SDK package and version;
- package and lockfile state;
- direct and aliased imports;
- local wrappers and adapters;
- affected call sites and data mappings;
- relevant configuration, tests, mocks, cassettes, and fixtures;
- generated, vendored, or excluded paths;
- dynamic, ambiguous, and unsupported usage.

Impact status is:

```text
unaffected
affected_supported
affected_needs_input
unsupported
uncertain
```

`unaffected` is allowed only when analyzed scope, coverage, and limitations are
explicit.

### 5. Produce And Approve The Migration Plan

`lumyn plan` names:

- every planned or conditionally planned file;
- every affected, excluded, and uncertain call site;
- dependency and lockfile intent;
- the route for each item: `deterministic`, `agent_assisted`, `manual`,
  `needs_input`, or `blocked`;
- agent context, tools, commands, model policy, and budgets;
- exact Agent Runner route, funding/credential/billing ownership, native
  configuration posture, and no-fallback policy;
- verification stages;
- delivery mode;
- residual risk and required human input.

Planning performs no write. Consumer authorization—either exact per-event
approval or installed preauthorization—binds the exact event, Provider Change
Contract, Consumer Installation and derived authorization, repository base,
scope, route, model policy, budgets, commands, verification, delivery, and
disclosure. Installed preauthorization pauses instead of writing when any
bound value falls outside policy.

### 6. Produce A Bounded Candidate

`lumyn apply` runs in an isolated worktree or consumer-approved equivalent:

1. deterministic codemods apply known safe mappings;
2. the bounded coding agent starts a clean session through the selected,
   qualified exact adapter and receives only the approved migration and
   repository context;
3. a Consumer Maintainer or Lumyn Operator may instead make approved manual
   edits in that worktree, after which `lumyn candidate import --manual`
   validates and binds the diff to the exact repository base, candidate head,
   Provider Change Contract, plan, and route;
4. all writes remain within approved files and diff budgets;
5. every edit maps to a Provider Change Contract item and repository-impact
   evidence;
6. no new network, credential, command, path, or task scope is inferred.

The default maximum is three agent or repair attempts. A higher limit requires
new Consumer Maintainer approval.

### 7. Verify, Then Separately Authorize Repair

`lumyn verify` is non-mutating with respect to the candidate. It executes:

1. migration-pack, plan, and candidate integrity;
2. pre-existing repository baseline comparison;
3. dependency and lockfile integrity;
4. compile and typecheck;
5. consumer-allowlisted tests;
6. optional contract or cassette replay;
7. optional exact-head mock or replay execution;
8. optional separately approved provider sandbox read-back;
9. boundary, redaction, cleanup, and evidence checks.

When a diagnostic failure is actionable and still within the approved plan,
the Consumer Maintainer may separately authorize `lumyn repair`.
`lumyn repair` is an agent-assisted action and requires
`agent_execution_policy=configured`. For a failed agent-assisted candidate, the
repair authorization reuses and binds its exact runner, model, funding,
credential, and billing route. For a failed deterministic or imported-manual
candidate, repair requires a newly configured, explicitly authorized exact
agent route; without it, Lumyn returns `needs_input` or `blocked`. Changing an
existing route likewise requires new explicit authorization. The
authorization also binds the failed candidate and evidence, exact repair
intent, remaining write and model-data permissions, and remaining time, token,
cost, attempt, file, and diff budgets. Every repair is a new attempt, creates a
new candidate head, invalidates prior verification evidence, and requires a
fresh `lumyn verify` run. It cannot expand semantic or file scope. A
non-diagnostic failure, missing business input, or exhausted budget becomes
`needs_input`, `blocked`, or `failed`.

Repository tests are untrusted code. They run in a consumer-approved execution
environment, without network or secrets by default. Dependency lifecycle
scripts, registry access, and sandbox credentials require distinct approval.

### 8. Deliver A Tested Draft PR

`lumyn export` always produces the local evidence bundle plus one of:

- patch;
- local branch;
- PR-ready bundle for manual push;
- a draft PR after short-lived GitHub authorization is implemented and
  separately approved.

A manually pushed branch or PR bundle is an acceptable assisted fallback. It
does not close automated-delivery acceptance, count as a Lumyn-opened draft
PR, or prove the self-maintaining application layer.

### 9. Record Adoption And Campaign Outcome

The Consumer Maintainer may record `consumer_accepted`, merged, closed,
reverted, or corrected outcomes. `consumer_accepted` requires a durable
consumer-acceptance artifact containing the Consumer Maintainer identity and
authority, API Consumer Organization, repository, exact candidate head,
verification-evidence digest, adoption decision, and timestamp. Plan approval,
candidate generation, PR creation, or an informal acknowledgement is not
consumer acceptance. Later correction or reversion appends a new outcome and
does not rewrite the historical acceptance record.

The API Provider receives only consumer-consented campaign status or
aggregates tied to the exact event and evidence boundary. Allowed states are
`received`, `not_applicable`, `affected`, `needs_input`, `candidate_ready`,
`verified`, `draft_pr_open`, `accepted`, `merged`, and `retired`. Silence is
`unknown`; neither `not_applicable` nor `unaffected` may be inferred without
explicit observed or consumer-reported evidence. `merged` is not inferred to
mean deployed or retired. Source code and raw execution evidence are not
provider-visible.

---

## Initial Technical Scope

### Supported Repository Shape

- Git repository; GitHub is the first optional remote-delivery target.
- TypeScript source discoverable from checked-in project configuration.
- One explicitly selected package root.
- One official npm SDK dependency and one source-to-target migration.
- Direct imports, aliases, local wrappers, adapters, and statically traceable
  call sites.
- `package-lock.json` is the first automatically writable lockfile.
- `pnpm-lock.yaml` and `yarn.lock` are impact-only until separately
  implemented and tested.
- Generated, vendored, minified, and build-output paths are excluded.
- A monorepo is supported only when one package root is explicitly selected.

### Supported Routing Classes

#### Deterministic

- SDK dependency and import update.
- Method or operation rename with explicit semantic equivalence.
- Request-property rename or relocation with no new business value.
- Response-property rename or relocation with statically identifiable access.

#### Agent Assisted

- Bounded updates to local wrappers, adapters, call sites, and directly
  related tests when the migration semantics are explicit.
- Type or signature adaptation that introduces no new business decision.
- Repository-specific repair after compile, typecheck, or allowlisted test
  failure when the repair remains within the approved migration scope.

#### Blocked

- A newly required business value without deterministic derivation or human
  input.
- Authentication or authorization redesign.
- Webhook or event-semantics changes.
- Production-only behavior or irreversible actions.
- Ambiguous provider semantics.
- Missing repository verification signal.
- Scope that exceeds approved paths, commands, credentials, network, model
  policy, or budgets.

### Explicitly Unsupported In The MVP

- GraphQL, gRPC, JSON-RPC, and body-dispatched operation migrations.
- Generated-client regeneration.
- Cross-language migrations.
- More than one selected package root.
- Broad production credentials or production mutation.
- Arbitrary provider-supplied executable scripts.
- Autonomous merge.
- General-purpose repository refactoring unrelated to the migration.

---

## Product Status Axes

### Impact

- `not_analyzed`
- `unaffected`
- `affected_supported`
- `affected_needs_input`
- `unsupported`
- `uncertain`

### Route

- `not_routed`
- `deterministic`
- `agent_assisted`
- `manual`
- `blocked`

### Candidate

- `not_attempted`
- `planned`
- `candidate_generated`
- `repairing`
- `needs_input`
- `failed`
- `stale`

### Verification

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

### Delivery

- `not_requested`
- `patch_exported`
- `local_branch_ready`
- `pr_bundle_ready`
- `remote_branch_pushed`
- `draft_pr_open`
- `consumer_accepted`
- `merged`
- `closed`
- `blocked`
- `superseded`

### Provider Rollout Projection

- `unknown`
- `received`
- `not_applicable`
- `affected`
- `needs_input`
- `candidate_ready`
- `verified`
- `draft_pr_open`
- `accepted`
- `merged`
- `retired`

Every projected state names its exact event and evidence plus `observed`,
`consumer_reported`, or `unknown` provenance. Silence is `unknown`; `merged`
does not imply `retired`.

No roll-up may hide a weaker axis.

---

## Evidence Contract

Every outcome reports:

- API Provider, change-event, and Provider Change Contract identity;
- consumer-installation and event-specific authorization digests;
- public-derived or provider-confirmed provenance;
- source and target versions and source digests;
- repository base and candidate head;
- analyzed paths, exclusions, and integration-graph uncertainty;
- affected, unaffected, unsupported, and uncertain call sites;
- deterministic, agent-assisted, and manual route by change item;
- changed files and dependency or lockfile delta;
- explicit `agent_execution_policy`;
- when agent execution is configured, Agent Runner Vendor, exact adapter
  version, resolved executable source/digest and conformance digest, auth mode
  and entitlement class, execution-funding mode, credential owner,
  usage-billing owner, native configuration/rules digest, and actual Model
  Provider, endpoint, model route, and version, including any runner-brokered
  downstream route;
- model-data policy and context-policy digest;
- tool, command, path, network, credential, time, token, cost, attempt, file,
  and diff budgets;
- attempt and repair history without persisting raw prompts or responses by
  default;
- human input and approval events;
- exact baseline and post-change commands and results;
- verification environment and observed outcome;
- pre-existing failures, residual risk, and unsupported items;
- redaction, boundary, and cleanup results;
- delivery mode and outcome;
- consumer-acceptance actor, authority, candidate/evidence binding, decision,
  and timestamp when `consumer_accepted` is claimed;
- artifact hashes and freshness;
- model/tool cost, elapsed time, Consumer Maintainer time, and Lumyn Operator
  time;
- rollback guidance and reviewer checklist.

Evidence becomes stale when the pack, repository base or head, plan,
`agent_execution_policy`, Agent Runner adapter/version/executable/conformance,
auth or entitlement, execution-funding mode, credential or usage-billing
owner, native-configuration posture, actual model route, model policy,
candidate, commands, or verification inputs change.

Raw prompts, model responses, source, diffs, logs, traces, agent sessions, and
credentials are consumer-private and never API-provider-visible.
Provider-visible output is limited to separately consented, enumerated
campaign status or aggregate fields.

---

## Command Model

`lumyn` remains the primary local surface.

| Command | Purpose |
|---|---|
| `lumyn init` | Initialize repo-local Lumyn configuration |
| `lumyn check` | Non-mutating source, repository, verification, and conditional Agent Runner readiness preflight |
| `lumyn pack` | Build or validate a Provider Change Contract from pinned evidence |
| `lumyn install` | Bind a provider channel to consumer-owned repository, action, optional Agent Runner/model and funding, GitHub, and disclosure policy |
| `lumyn update --event` | Process one event through installed policy, stopping at the selected action mode |
| `lumyn impact` | Produce a read-only repository impact inventory |
| `lumyn plan` | Produce a no-write routed migration plan |
| `lumyn apply` | Produce a bounded deterministic or agent-assisted candidate |
| `lumyn candidate import --manual` | Validate and bind an approved manual candidate to the base, pack, plan, and route |
| `lumyn verify` | Non-mutating baseline-aware candidate verification |
| `lumyn repair` | Separately authorized agent-assisted repair that creates a new candidate for fresh verification |
| `lumyn export` | Export evidence plus patch, branch, or PR-ready bundle fallback |
| `lumyn trace` | Render local evidence without changing state |
| `lumyn outcome record` | Record durable consumer acceptance, merge, closure, correction, or reversion evidence |
| `lumyn pr create --draft` | Required pilot proof using short-lived authorization |

The command grammar may be implemented incrementally. Unimplemented commands
return a typed nonzero result and never a generic successful envelope.

Stable exit codes `0` through `9` remain reserved according to the existing
compatibility contract. Exit code `6` remains reserved. API Provider identity
uses `api_provider_id` or `change_authority`; Agent Runner Vendor metadata uses
`agent_runner_vendor_metadata`; Model Provider metadata uses
`model_provider_metadata`.

---

## Artifact Model

The MVP introduces only:

- `provider-change-event`
- `migration-pack`
- `consumer-installation`
- `integration-graph`
- `impact-report`
- `migration-plan`
- `candidate-manifest`
- `agent-runner-manifest`
- `agent-runner-conformance-result`
- `agent-attempt`
- `migration-verification`
- `export-result`
- `campaign-summary`
- `provider-status-projection`
- `remediation-outcome`

An artifact becomes executable only when its schema, valid and invalid
fixtures, compatibility posture, and validator tests ship together.

Product artifacts are separated as:

```text
consumer checkout
  lumyn.yaml
  provider-channel and action policy
  public or synthetic fixtures only

consumer-private state root outside the checkout
  events/
  installations/
  packs/
  graphs/
  impacts/
  plans/
  candidates/
  attempts/
  verification/
  exports/

consumer-consented API Provider summary
  event-bound rollout status or aggregate only
```

The private root cannot resolve inside the checkout or public source
repository. Private artifacts have explicit retention and deletion policy.
The first services pilot may use documented operator-managed cleanup; a
productized deletion-receipt and orphan-recovery control plane is deferred.

---

## Trust, Authorization, And Data Boundaries

### Migration Evidence

- Provider Change Contracts and events are declarative, versioned,
  digest-bound, and inspectable.
- Public-derived evidence is labeled honestly.
- A provider-funded campaign requires accountable provider confirmation before
  a real consumer mutation.
- Conflicting or incomplete semantics block affected routes.
- Provider-supplied executable migration code is prohibited.

### Repository

- Analysis requires explicit read scope and performs no mutation.
- Write scope is separate and binds the exact plan and base.
- Candidate generation runs in an isolated worktree or consumer-approved
  equivalent.
- Consumer Installation, event-specific authorization, remote branch, and PR
  actions are separate from local mutation.
- Default-branch write and auto-merge are prohibited.

### Coding Agent And Model

- `agent_execution_policy` defaults to `disabled`. Notify-only, scan-only, and
  deterministic-only installations need no runner account or credential. If a
  plan needs agent assistance, execution pauses until the Consumer Maintainer
  explicitly changes the policy to `configured`.
- A configured policy selects one qualified exact Agent Runner adapter and
  version plus `consumer_managed` or
  `provider_sponsored_lumyn_managed` execution funding. The selection binds
  credential owner, usage-billing owner, Agent Runner Vendor, Model Provider or
  local mode, model route, context boundary, data egress, retention posture,
  tools, native agent configuration, and budgets.
- Codex and Claude Code are launch targets only after independent conformance
  and live-canary qualification. Cursor remains deferred until it passes the
  identical gate.
- Every attempt resolves an approved executable by canonical path and digest,
  then starts a clean ephemeral session with a neutral home/config root.
  Repository-local PATH shadowing, personal sessions, unrelated history,
  ambient memories, and user/project rules are not reused. Supported static
  native instructions may be used only when explicitly selected, digest-bound,
  treated as untrusted context, and unable to widen Lumyn authority.
  Executable plugins, MCP servers, and hooks are prohibited for the MVP.
- The selected auth mode and entitlement class must allow the intended local
  or CI automation under the consumer's vendor agreement and organization
  policy and expose the actual downstream provider, endpoint, model, and
  version. Unsupported entitlement or opaque/changing routing blocks
  execution; the consumer may instead approve a qualifying BYOK, local, or
  Lumyn-managed route.
- No selected adapter, version, model route, endpoint, credential owner, or
  usage-billing owner may change through fallback. Unavailability or
  conformance failure blocks the agent route; a separately valid deterministic
  route may still proceed.
- The agent receives no ambient secrets.
- Raw repository content is not transmitted to a cloud Model Provider unless
  the consumer explicitly approves that mode.
- Consumer-managed agent accounts, enterprise subscriptions, BYOK, and local
  models are supported policy choices. The optional Lumyn-managed route uses
  only task-scoped brokered credentials with exact audience, TTL, one-time
  redemption into a quota-bound attempt session, no refresh or cross-attempt
  replay/reuse, revocation, and reconciliation controls and remains
  consumer-authorized. If neither the vendor nor an approved budget-enforcing
  proxy can enforce them, that route is unavailable.
- The API Provider never receives agent context by virtue of sponsorship.
- Agent output and self-reported test results are generation evidence only.
  Lumyn's independent verifier runs the approved commands from the exact
  candidate head in a fresh process and view outside the agent session, with
  frozen command/config digests, no runner/model credentials, and no
  generation-owned evidence handle.

### Commands, Credentials, And Network

- Commands are allowlisted and displayed before execution.
- Repository commands run in a consumer-approved environment with explicit
  mounts, environment, process, timeout, and output limits.
- The Agent Runner process uses explicit mounts, no host home or OS credential
  store, no ambient service sockets or unrelated inherited descriptors,
  inherited child-process restrictions, host-enforced egress, and recorded
  cleanup. An unenforceable boundary blocks execution.
- Tests run without network and secrets by default.
- Dependency lifecycle scripts and registry access require separate approval.
- Sandbox network, transmitted payload classes, and non-production credentials
  require separate approval.
- Production credentials and production mutation are prohibited.

### Sharing

The API Provider never receives raw consumer source, diffs, prompts, model
responses, agent sessions, tool traces, logs, or credentials. Consumer consent
is required before sharing any enumerated campaign-status or aggregate field,
including repository identity, impact counts, candidate state, verification
state, failure reason, merge state, or retirement state.

Provider reporting may initially be a locally produced, consented projection.
A hosted coordinator, universal registry, public changelog monitor,
cryptographic status exchange, and billing receipts are deferred until demand
and operational need are proven.

---

## Functional Requirements

### FR1: Provider Change Intake

Lumyn must normalize pinned official or public migration evidence into a
versioned, provenance-visible, non-executable Provider Change Contract and
bind a signed event from the pinned provider-controlled channel to its exact
version, transport identity, sequence, freshness, audience, and authentication
evidence.

### FR2: Consumer Installation And Repository Impact

Lumyn must derive each run's authority from a consumer-owned installation, then
identify dependency state, imports, wrappers, call sites, mappings, relevant
tests, and explicit uncertainty inside that authorized scope. Installation
action and authorization modes are ceilings enforced before every side effect,
not ambient grants. Task- or campaign-level authority arrays are capability
universes only: before every product action, Lumyn freezes the exact selected
route and capability union. A composed campaign reuses those action contracts
instead of granting their aggregate union to every installation.

### FR3: Read-Only Impact

Lumyn must produce an impact report without modifying repository, dependency,
lockfile, branch, or remote state.

### FR4: Routed Migration Plan

Lumyn must produce a reviewable no-write plan that classifies every affected
item as deterministic, agent-assisted, manual, needs-input, or blocked.

### FR5: Deterministic Execution

Known safe transforms must be repeatable for identical pinned inputs and stay
within approved boundaries.

### FR6: Bounded Agent Execution

Agent-assisted work must enforce approved context, paths, tools, commands,
qualified Agent Runner adapter/version, funding/credential/billing ownership,
native-configuration posture, Agent Runner and model policy, attempts, time,
token, cost, file, and diff budgets. It must expose normalized lifecycle and
provenance evidence, must not reuse personal sessions, and must not silently
fall back to another route. The runner process must enforce explicit mounts,
credential/socket/descriptor denial, inherited child-process restrictions,
egress, and cleanup; executable plugins, MCP servers, and hooks are prohibited
for the MVP.

### FR7: Repair Loop

Lumyn may use diagnostic verification failures for separately authorized,
agent-assisted bounded repair without expanding the approved plan. Repair
requires configured agent execution; a deterministic or manual candidate
requires a newly authorized exact agent route, while an agent candidate reuses
its bound route unless a new route is explicitly authorized. Exhausted,
unauthorized, or non-diagnostic cases stop.

### FR8: Baseline-Aware Verification

Lumyn must distinguish pre-existing failures from migration-attributable
failures and preserve the exact evidence boundary.

### FR9: Evidence-Bound Outcome

Every outcome must bind pack, base, candidate, plan, route,
`agent_execution_policy`, commands, verification, residual risk, and artifact
hashes. Configured agent execution must additionally bind exact Agent
Runner/model/funding/credential/billing provenance.

### FR10: Fallback Export

Lumyn must export evidence plus a patch, branch, or PR-ready bundle without
requiring a hosted service.

### FR11: Tested Draft PR

Draft-PR delivery must use short-lived, least-privilege authorization and
remain distinct from code generation and verification. The composed
`lumyn update --event` path must bind provider event, installation, impact,
plan, Lumyn-generated candidate, verification, branch, and PR evidence. At
least one Lumyn-opened draft PR is required to prove the first provider
campaign. Its immutable event-specific authorization snapshot must expose and
revalidate the union of every product action exercised by that run; completed
milestones and Factory worker grants do not delegate consumer authority.

### FR12: Consented Campaign Reporting

Provider-visible reporting must contain only consumer-consented fields and
must never be a prerequisite for private consumer execution.

### FR13: Outcome Feedback

Lumyn must record accepted, merged, closed, reverted, and substantively
corrected outcomes with provenance.

### FR14: Stable Machine Interface

State-returning commands must support stable JSON, typed errors,
non-interactive execution, and versioned artifacts.

---

## Non-Functional Requirements

### NFR1: Fail-Closed Honesty

Unsupported, ambiguous, unauthorized, unredactable, out-of-budget, or
unverified states never report a verified migration.

### NFR2: Controlled Determinism

Pinned inputs, deterministic transforms, routing, verification commands,
evidence, and status are repeatable. Agent source output need not be
byte-identical.

### NFR3: Consumer Privacy

Consumer code and private evidence remain in the approved consumer execution
plane except for context explicitly disclosed to the Agent Runner Vendor and
downstream Model Provider. The API Provider receives only enumerated,
consumer-consented status or aggregate fields and never raw code, diffs,
prompts, responses, agent sessions, tool traces, logs, or credentials.

### NFR4: Bounded Autonomy

An agent cannot widen paths, tools, commands, network, credentials, native
configuration, Agent Runner/model policy, task intent, or budgets.

### NFR5: Least Privilege

Read, write, command, Agent Runner network and credential, model-data, model
network and credential, registry, sandbox, remote branch, PR, retention,
deletion, and provider-reporting scopes are independent.

### NFR6: Redaction Before Persistence

Secrets are redacted before persistence or sharing. Redaction uncertainty
blocks the affected artifact.

### NFR7: Explainability

Every edit maps to a migration-pack item, repository evidence, route, and
rationale. Unsupported cases include a reason and next action.

### NFR8: Artifact Stability

Persisted artifacts are versioned, schema-backed, digest-bound, and
migration-aware.

### NFR9: Recovery And Idempotency

Interrupted or repeated runs do not duplicate candidates, branches, PRs,
sandbox resources, or API Provider campaign summaries.

### NFR10: Bounded Performance And Cost

On the fixed benchmark, median read-only impact completes in under five
minutes. Every agent run has explicit time, token, cost, attempt, file, and
diff budgets.

### NFR11: No Production Dependency

The benchmark and default consumer run require no production credentials,
traffic, or mutation.

---

## Failure Taxonomy

Required failure classes include:

- `migration_pack_invalid`
- `migration_evidence_ambiguous`
- `source_target_mismatch`
- `provider_confirmation_missing`
- `authorization_missing`
- `authorization_expired`
- `read_scope_exceeded`
- `write_scope_exceeded`
- `unsupported_repository_shape`
- `unsupported_package_manager`
- `multiple_sdk_versions`
- `repo_context_insufficient`
- `dynamic_usage_uncertain`
- `impact_uncertain`
- `required_business_value_missing`
- `unsupported_change_class`
- `test_signal_insufficient`
- `plan_stale`
- `candidate_conflict`
- `candidate_stale`
- `diff_budget_exceeded`
- `agent_not_authorized`
- `agent_runner_unqualified`
- `agent_runner_unavailable`
- `agent_runner_auth_failed`
- `agent_runner_entitlement_invalid`
- `agent_runner_executable_untrusted`
- `agent_runner_contract_violation`
- `agent_runner_fallback_not_authorized`
- `usage_billing_owner_ambiguous`
- `agent_scope_exceeded`
- `agent_budget_exhausted`
- `repair_limit_reached`
- `model_unavailable`
- `model_data_policy_blocked`
- `human_input_required`
- `baseline_already_failing`
- `dependency_integrity_failed`
- `compile_failed`
- `typecheck_failed`
- `tests_failed`
- `tests_flaky`
- `verification_non_diagnostic`
- `replay_failed`
- `sandbox_unavailable`
- `workflow_proof_gap`
- `redaction_uncertain`
- `evidence_stale`
- `export_failed`
- `duplicate_pr`
- `provider_reporting_not_consented`

Failures cite concrete artifact references. Agent diagnosis is labeled as
interpretation until independently verified.

---

## Acceptance Tests

The v3 acceptance ledger contains 53 item-level closure units. Group headings,
milestones, phases, and an overall MVP label cannot substitute for evidence
against the applicable item IDs.

### Retained Foundation

1. `BASE-001`: The Go CLI, configuration, stable result envelope, exit codes,
   `lumyn init`, and `lumyn check` remain functional.
2. `BASE-002`: Existing workflow, evidence, cassette, trace, proof, boundary,
   redaction, and command-result schemas remain executable and versioned.
3. `BASE-003`: CI, coverage, CodeQL, CODEOWNERS, required checks, review,
   shipping, and post-merge governance remain enforced.
4. `BASE-004`: Existing OpenAPI and local-doc intake continues to produce
   structured refs, fingerprints, deprecation findings, and concrete source
   locations.
5. `BASE-005`: Unimplemented commands return typed nonzero results, and API
   Provider, Agent Runner Vendor, and Model Provider terminology is
   unambiguous.

### Provider Change Contract

1. `PACK-001`: Pinned public docs, OpenAPI descriptions, SDK releases, and
   migration guidance can produce a versioned `public_derived` Provider Change
   Contract with source digests and license/provenance records.
2. `PACK-002`: A Provider Change Contract records stable change identity,
   source and target versions, audience and applicability, typed semantics,
   known mappings, verification guidance, ambiguity, unsupported cases, and
   supersession or withdrawal state.
3. `PACK-003`: Accountable provider confirmation is independently recorded
   once and is reusable across the invited, consenting cohort for the exact
   contract version; public-derived evidence is never represented as provider
   endorsement, and provider input cannot execute code or grant repository
   authority.
4. `PACK-004`: Conflicting, incomplete, stale, or semantically ambiguous
   evidence blocks affected routes instead of inviting the agent to infer
   provider intent.

### Provider Change Event

1. `EVENT-001`: Lumyn accepts a Provider Change Event only when its immutable
   identity and version, provider-controlled pinned transport or
   attended-recovery provenance, issuer and authentication, monotonic sequence
   and freshness, embedded or exact-URL contract delivery and retrieved-byte
   digest, API or SDK, audience, deadline, severity, and supersession or
   withdrawal state validate; the event and Provider Change Contract remain
   non-executable.
2. `EVENT-002`: Duplicate, replayed, stale, conflicting, expired, superseded,
   withdrawn, wrong-audience, or unauthenticated events fail closed; an
   attended import cannot count as provider-channel delivery or authorize an
   installed-preauthorization write.

### Consumer Installation

1. `INSTALL-001`: A revocable Consumer Installation binds the exact
   provider-controlled channel and authentication key, repository and package
   root, audience or version selectors, action ceiling, authorization mode,
   paths, commands, `agent_execution_policy` (`disabled` or `configured`),
   model-data policy and budgets, GitHub token-issuance policy,
   provider-reporting fields, retention and deletion, expiry, and revocation.
   A configured policy additionally binds the exact Agent Runner
   adapter/version/executable qualification and auth/entitlement policy,
   execution-funding mode, credential and usage-billing ownership, and native
   agent-configuration policy without storing a reusable Agent Runner, model,
   or GitHub credential; a disabled policy grants none of those capabilities.
2. `INSTALL-002`: Every event-specific authorization is an immutable,
   non-widening derivation of the installation; `per_event_approval` binds the
   exact plan, while `installed_preauthorization` permits side effects only
   when event, plan, commands, paths, model policy, budgets, verification, and
   short-lived credential issuance all satisfy the installed policy. An
   action-mode label alone grants no side effect, and an `agent_assisted` route
   pauses unless the installation explicitly configures and authorizes it.

### Impact And Repository Inventory

1. `IMPACT-001`: Impact analysis is read-only and remains inside the exact
   authorized repository and package-root scope.
2. `IMPACT-002`: Lumyn detects the selected official npm SDK, installed
   version, package manifest, supported lockfile, direct imports, and aliases.
3. `IMPACT-003`: The repository impact inventory represents statically traceable
   wrappers, adapters, call sites, request/response mappings, relevant
   configuration, tests, mocks, cassettes, and fixtures.
4. `IMPACT-004`: Generated, vendored, dynamic, multi-root, multi-version, and
   unsupported cases receive explicit exclusion, uncertainty, or failure
   status; no result implies global downstream coverage.
5. `IMPACT-005`: On frozen held-out fixtures, supported affected sites meet
   the preregistered recall and precision thresholds with zero false
   `unaffected` results; real-repository value is not claimed from fixtures
   alone.

### Plan And Routing

1. `PLAN-001`: The plan names every proposed or conditional file, route,
   command, explicit `agent_execution_policy`, permission, budget,
   verification stage, and residual risk without modifying repository or Git
   state; a configured agent policy additionally names the exact Agent
   Runner/model and funding/credential/billing policy and native-configuration
   posture.
2. `PLAN-002`: Every affected item is explicitly classified as
   `deterministic`, `agent_assisted`, `manual`, `needs_input`, or `blocked`.
3. `PLAN-003`: Consumer authorization, whether exact per-event approval or
   installed preauthorization, binds the event, Provider Change Contract,
   Consumer Installation and derived authorization, repository base, paths,
   route, `agent_execution_policy`, model-data policy, tools, commands, budgets,
   disclosure, delivery, and verification intent. A configured agent policy
   additionally binds Agent Runner adapter/version/executable/conformance and
   auth/entitlement policy, execution funding, credential and usage-billing
   owners, and native configuration; any out-of-policy change invalidates or
   pauses the run.

### Bounded Hybrid Execution

1. `AGENT-001`: Deterministic transforms produce byte-stable candidates for
   identical pinned inputs.
2. `AGENT-002`: Agent-assisted execution requires
   `agent_execution_policy=configured` and one Consumer Maintainer-selected,
   qualified exact Agent Runner adapter/version/executable digest, allowed auth
   mode and entitlement class, and execution-funding mode; it starts a clean
   ephemeral session with a neutral home/config root only in a consumer-local
   or explicitly consumer-approved environment and discloses Agent Runner
   Vendor, actual Model Provider or local route, credential owner,
   usage-billing owner, context, data egress, retention, and native
   configuration posture before execution. Codex and Claude Code are launch
   targets only after passing the common gate; Cursor is deferred behind it.
3. `AGENT-003`: Lumyn's orchestrator and host-isolation boundary enforce
   approved path, tool, command, Agent Runner and model network, Agent Runner
   and model credential, time, token, cost, attempt, file, and diff budgets;
   native agent configuration is disabled or digest-bound as untrusted
   static context; explicit mounts, OS-credential denial, socket and inherited
   descriptor denial, child-process restriction inheritance, egress, and
   cleanup evidence are required; the exact isolation backend, version,
   configuration and qualification digests, host platform, and hard CPU,
   memory, PID, process-tree-depth, disk, and open-file quotas are bound and
   resource-exhaustion tested; executable plugins, MCP servers, and hooks are
   prohibited, and no adapter, model, credential, or billing-owner fallback
   occurs.
4. `AGENT-004`: Every agent edit maps to a Provider Change Contract item,
   repository-impact evidence, and recorded rationale; unrelated edits fail.
5. `AGENT-005`: The repair loop performs no more than the approved attempts,
   requires configured agent execution, reuses a failed agent candidate's bound
   Agent Runner/model/funding/credential/billing route unless a new explicit
   authorization creates a new attempt and candidate, requires a newly
   authorized exact agent route for a deterministic or manual candidate,
   cannot expand scope, and stops on missing authorization, non-diagnostic
   failure, or exhausted budget.
6. `AGENT-006`: Missing business values, auth redesign, event semantics,
   production-only behavior, and ambiguous provider intent return
   `needs_input` or `blocked` without speculative completion.
7. `AGENT-007`: Agent evidence records Agent Runner Vendor, exact adapter
   version, executable source/digest and conformance digest, auth mode and
   entitlement class, clean-session identity, execution-funding mode,
   credential and usage-billing owners, actual Model Provider, model class and
   version, native-configuration state/digest, model-data and context-policy
   digests, tools, commands, attempts, token and cost use, changed files, and
   human input, never persists reusable credentials or secrets, and does not
   persist raw prompts or responses by default. Managed credentials
   additionally record broker issuer, exact audience, TTL,
   one-time-redemption/attempt-session
   posture, no refresh or cross-attempt replay/reuse, hard quota, revocation,
   and reconciliation without secret values.

### Verification

1. `VER-001`: Lumyn records pre-existing dependency, compile, typecheck, and
   selected-test failures before candidate generation.
2. `VER-002`: Deterministic, agent-assisted, and imported manual candidates run
   the same dependency-integrity, compile/typecheck, and
   consumer-allowlisted test ladder in a fresh verification process and view
   with frozen command/config digests, no Agent Runner/model credentials, and
   no generation-owned verification-evidence write handle. Repository commands
   bind a qualified isolation backend identity, exact commands, mounts and
   environment, hard CPU, memory, PID/process-tree, disk and open-file quotas,
   inherited child restrictions, offline/lifecycle defaults, and cleanup
   evidence.
3. `VER-003`: Static, repository, independent contract replay, exact-head
   replay, mock, and sandbox evidence use distinct canonical labels.
4. `VER-004`: Every `workflow_verified_*` result causally executes an approved
   entrypoint from the exact candidate head and records observed interaction
   and outcome evidence. A sandbox entrypoint uses a separate qualified
   isolation profile with a read-only exact-head mount, sole task-scoped
   sandbox credential injection, endpoint-only egress, hard resource quotas,
   inherited child restrictions, teardown, cleanup, and orphan evidence.
5. `VER-005`: Failed, missing, stale, out-of-boundary, unredactable, or
   inconclusive evidence cannot produce a verified label; repair failures
   remain visible.
6. `VER-006`: The frozen negative suite has zero false verified outcomes, and
   all evidence binds event, Provider Change Contract, Consumer Installation
   and derived authorization, plan, base, candidate, route,
   `agent_execution_policy`, commands, environment, and artifact hashes;
   configured agent execution additionally binds exact Agent Runner, model,
   funding, credential, and billing provenance.

### Export And Delivery

1. `EXP-001`: Lumyn exports complete local evidence plus a patch, local branch,
   or PR-ready bundle without requiring a hosted service.
2. `EXP-002`: A manual branch or PR bundle is labeled manual delivery and
   cannot count as automated PR delivery.
3. `EXP-003`: From a valid provider-channel event accepted by an
   `open_draft_pr` Consumer Installation, `lumyn update --event` composes
   impact, plan, a Lumyn deterministic or agent-assisted candidate,
   independent verification, non-default-branch write, and tested draft-PR
   creation using a short-lived least-privilege token and evidence-bound
   idempotency, with no auto-merge, then records a local event-bound status
   projection or explicit reporting decline. Provider transmission remains
   separately consented and cannot block the draft PR. Standalone PR creation,
   an imported manual candidate, or a manual patch, branch, or PR bundle cannot
   close this item.
4. `EXP-004`: Provider-visible status is bound to the exact event and evidence,
   contains only consumer-consented fields, distinguishes observed,
   consumer-reported, and unknown state, never infers `not_applicable` or
   `unaffected` from silence or `retired` from merge, and excludes raw source,
   diffs, prompts, responses, agent sessions, tool traces, logs, and
   credentials.

### Trust And Privacy

1. `TRUST-001`: A revocable Consumer Installation binds provider or channel,
   repository and package root, audience or version selectors, allowed
   actions, paths, commands, `agent_execution_policy`, model-data posture,
   registry, sandbox, remote branch, PR, retention, deletion, provider
   reporting, expiry, and disclosure. A configured agent policy additionally
   binds Agent Runner selection, executable, auth, entitlement, funding, native
   configuration, Agent Runner/model egress, network and credentials, and
   credential and usage-billing owners; each event-specific authorization is
   derived without widening that installed policy.
2. `TRUST-002`: Provider Change Contracts, events, and provider artifacts
   cannot execute code; production credentials and production mutation are
   prohibited.
3. `TRUST-003`: When agent execution is configured, Agent Runner Vendor and
   Model Provider processing, cloud-model context transfer, qualifying
   consumer-managed account/subscription/BYOK/local posture, optional
   provider-sponsored Lumyn-managed usage, actual model-route identity,
   credential owner, and usage-billing owner are explicit; the API Provider
   gains no code, context, session, or credential access, and secrets or
   prohibited data are redacted or blocked before runner/model egress or
   artifact persistence.
4. `TRUST-004`: Consumer-private artifacts remain outside the checkout and
   public repository, and provider sponsorship never grants source, Agent
   Runner/model context, session, or credential access.

### Design-Partner Qualification

1. `DISC-001`: One qualified API Provider clears at least `$25,000` in
   non-refundable prepaid funds for a defined migration campaign and names an
   economic buyer, Provider Operator, hard deadline, source and target
   versions, and decision process.
2. `DISC-002`: Before campaign-specific migration execution, the API Provider
   commits to distribute the update and lead onboarding for at least five
   reachable Eligible Consumer Units across five distinct API Consumer
   Organizations, each with an accountable maintainer and useful verification
   signals; at least one consumer prequalifies a consumer-selected Codex or
   Claude Code route, or explicitly consents to a qualifying managed route, for
   a real agent-assisted run without disclosing a reusable credential. That
   repository also has a plausible naturally agent-eligible migration
   hypothesis grounded in a wrapper, adapter, signature/type adaptation, or
   related test repair that needs repository-specific reasoning but no new
   business decision. This prequalification cannot force the route; M4 impact
   and M5 planning confirm or reject it from authorized evidence. The same
   intended qualifying consumer reviews the exact allowlisted event-bound
   provider-status projection and records willingness to authorize
   transmission if the run qualifies; willingness is not runtime consent.
3. `DISC-003`: After the v3 privacy, model-data, authorization, and evidence
   contracts are approved, and before the first invitation, the parties freeze
   cohort, eligibility, price, evidence sources, privacy and combined Agent
   Runner and model-data protocol, permitted execution-funding modes,
   credential and usage-billing ownership, baseline method, material
   maintainer comparison threshold, correction rubric, campaign COGS boundary,
   minimum contribution margin or maximum Lumyn Operator hours per reviewable
   outcome, measurement
   windows, one material provider-outcome metric and threshold, and absolute
   judgment deadline. A provider threshold is material only when the economic
   buyer records that meeting it would justify a retirement or paid
   continuation decision and missing it makes the campaign fail.

### Provider Campaign Pilot

1. `PILOT-001`: At least five prequalified Eligible Consumer Units across five
   distinct API Consumer Organizations are invited without changing the
   frozen cohort after outcomes are visible.
2. `PILOT-002`: At least three repositories complete valid consumer-local
   impact scans within 14 calendar days of their invitations.
3. `PILOT-003`: At least two distinct repositories produce tested, reviewable
   migration outcomes with explicit evidence boundaries by the frozen pilot
   deadline. One qualifying outcome itself starts from the authenticated
   provider channel, uses `installed_preauthorization`, and contains at least
   one plan item organically classified `agent_assisted` because
   repository-specific reasoning is necessary. It runs that item through a
   consumer-selected qualified Agent Runner without bespoke Lumyn Operator
   code edits, passes
   independent exact-head verification, and is delivered by the composed
   `lumyn update --event` flow as a Lumyn-opened draft PR. A
   consumer-consented provider status projection bound to that same event,
   installation authorization, candidate, verification evidence, and draft PR
   reaches the API Provider. Deterministic work may not be rerouted merely to
   manufacture agent proof; a campaign with no organically agent-eligible
   qualifying run fails this item.
4. `PILOT-004`: At least one verified migration outcome is accepted or merged
   by its Consumer Maintainer; closed, rejected, reverted, and corrected
   outcomes remain visible.
5. `PILOT-005`: Against the frozen guide/codemod plus generic-agent baseline
   using the same repository snapshot, authoritative migration evidence,
   selected Agent Runner adapter/version/executable, actual Model Provider and
   model version, auth/entitlement and execution-funding route, credential and
   usage-billing owners, context-access ceiling, tools, commands, and attempt,
   token, time, and cost budgets, Lumyn reduces median Consumer Maintainer
   hands-on time by at least 30% without a worse substantive-correction,
   revert, or false-verification rate, or clears another equally material
   threshold frozen before execution. Unmatched engine comparisons are
   descriptive, not causal proof of Lumyn's advantage.
6. `PILOT-006`: Actual cleared payment, Agent Runner/model/tool spend by
   execution-funding mode and usage-billing owner, the subset attributable to
   Lumyn campaign COGS, Lumyn Operator hours, Consumer Maintainer time, support
   effort, and the frozen primary provider outcome are measured from source
   evidence, and the provider outcome clears its frozen material threshold
   while the campaign clears its frozen contribution-margin or
   operator-automation threshold.
7. `PILOT-007`: Every migration records accepted, merged, closed, reverted,
   corrected, blocked, and residual-risk outcomes through the frozen
   observation window; each participating installation emits or explicitly
   declines its allowlisted event-bound status projection, and the projection
   from the qualifying `PILOT-003` run reaches the API Provider.
8. `PILOT-008`: The campaign receives a pass or fail at the frozen judgment
   deadline. Failure to prepay, recruit, activate, outperform the generic
   baseline, clear the material provider-outcome threshold, or produce
   verified outcomes triggers a documented stop or reframe; the experiment
   cannot remain open indefinitely.

---

## Success Metrics

### North Star

`share of all Eligible Consumer Units in the frozen targeted cohort with an
accepted or merged verified Lumyn PR by the provider deadline`

The denominator is every unit in the frozen targeted cohort, including units
that are not invited, do not install, decline authorization, remain
unknown-impact, or become blocked; those units do not disappear from the
provider-level metric. Report invitation-to-install conversion,
installation-to-authorized conversion, affected share among valid scans, and
accepted-or-merged success among affected installations separately. Each
outcome carries its evidence boundary. A candidate or opened PR does not count
as accepted or merged, but draft-PR creation is reported as the required
delivery-stage proof.

### Provider Outcome

Freeze exactly one primary provider outcome and material pass threshold before
the first invitation. The default is the share of the frozen target cohort
that is off the deprecated API or SDK version by the migration deadline.
Provider support hours per accepted migration or
invitation-to-accepted-migration lead time may be selected instead when the
economic buyer records that the threshold is material to retirement or paid
continuation. The cohort denominator and legacy-version retirement share
remain reported whenever they are measurable.

### Funnel

```text
eligible
-> invited
-> installed
-> consented
-> impacted
-> migration planned
-> candidate produced
-> repo verified
-> draft PR open
-> consumer accepted
-> merged
-> deployed when explicitly observed or consumer-reported
-> legacy version retired
```

Report conversion and elapsed time between every stage.

### Technical Quality

- Impact recall and precision by supported route.
- False `unaffected` rate.
- Deterministic transform repeatability.
- Agent scope-violation rate.
- Agent-assisted completion rate.
- Agent Runner install, authentication, conformance, and fallback-block rate by
  adapter/version and execution-funding mode.
- Time from installation to first qualified agent attempt.
- Human-input and blocked rate.
- Verification pass and false-verification rate.
- Substantive correction and revert rate.
- Unrelated-edit rate.
- Repair attempts per accepted outcome.
- Median impact, generation, verification, and export time.

### Economics And Effort

- Provider campaign revenue and actual payment.
- Agent Runner, model, tool, and infrastructure spend by execution-funding mode
  and usage-billing owner, with Lumyn campaign COGS kept separate from
  consumer-paid usage.
- Campaign contribution.
- Consumer Maintainer hands-on time.
- Lumyn Operator hours per verified and accepted migration.
- Provider setup and support hours.
- Cost per verified and accepted migration.
- Guide/codemod plus generic-agent baseline delta.
- Legacy cohort retirement rate.

`substantive_manual_correction` means a human edit after Lumyn export that
changes migrated API invocation, request or response mapping, error handling,
workflow behavior, or another Lumyn-generated semantic edit. Formatting,
comments, deterministic lockfile normalization, and unrelated pre-existing CI
repair do not count, but remain recorded.

---

## Distribution

The initial product motion is a per-provider update channel with
services-assisted onboarding:

1. Sell an active sunset or deprecation campaign to an accountable provider
   buyer.
2. Prepay and qualify the cohort before substantial campaign-specific
   automation.
3. Lumyn prepares one Provider Change Contract and signed event manifest with
   the Provider Operator; the provider publishes it on the pinned configured
   channel.
4. The API Provider recruits participating Consumer Maintainers and asks them
   to install its Lumyn update channel.
5. Consumers install exact action and disclosure policy, then run impact,
   planning, generation, and verification locally or in consumer-controlled
   CI.
6. Lumyn opens a tested draft PR under short-lived consumer authorization and
   projects only consented, event-bound status to the provider.
7. Lumyn operators assist the first cohort, measure COGS, and productize the
   repeated provider-to-consumer loop.

The secondary motion is a paid urgent consumer upgrade sprint. Before M6 it
provides paid workflow/problem evidence and cash; after M6 qualification it
may prove engine value. It does not validate provider-led distribution.

The MVP does not require hosted SaaS, a campaign dashboard, an annual
connected-repository program, or public OSS distribution. It may use a
narrowly scoped GitHub App installation that persists at the consumer's
discretion, but never a long-lived token or broad ambient grant. It requires
short-lived automated draft-PR delivery for at least one pilot repository.

Public fixtures are engineering benchmarks, not provider endorsement or
customer proof.

---

## Current Repository Baseline

Implemented today:

- Go CLI, configuration, exit-code, and command-result foundation.
- `lumyn init`.
- `lumyn check`.
- OpenAPI and local-doc parsing, refs, fingerprints, and findings.
- Executable workflow, evidence, cassette, proof, boundary, redaction, and
  command-result schemas.
- CI, coverage, CodeQL, branch policy, review, Factory planning, shipping, and
  post-merge governance.

Designed, schema-backed, or planned but not implemented as the v3 runtime:

- Provider Change Contract and event normalization;
- consumer installation and event-specific authorization;
- TypeScript repository impact inventory;
- deterministic migration transforms;
- customer-selected Agent Runner contract, Codex and Claude Code adapters, and
  consumer-approved execution boundary;
- bounded agent and repair budgets;
- hybrid migration planning and routing;
- repository verification orchestration;
- replay, mock, or sandbox verification runtime;
- agent provenance and cost evidence;
- patch, branch, and PR-bundle export;
- short-lived GitHub draft-PR delivery;
- event-bound consented provider status projection;
- migration outcome ingestion.

The current command dispatcher recognizes several unimplemented commands and
can return a base success envelope. M0 must correct that false-green behavior;
this planning rebaseline does not implement the correction.

Historical v2 planning and task evidence remains immutable. The current
repo-local compiled Factory set was regenerated from this v3 PRD and
`docs/product/plan.md`. It is planning and validation authority, not product
runtime or factoryd execution authority. factoryd remains paused until the
external Factory profile and factoryd runtime are separately aligned and
approved.

---

## Risks

1. **Generic-agent substitution:** The guide plus the same selected qualified
   runner—Codex or Claude Code at launch—may perform the code edit equivalently,
   leaving insufficient value in the installed provider-to-consumer loop.
2. **Provider activation:** A willing provider may fail to recruit consumers.
3. **Consumer trust:** Repository or model-data policies may block execution.
4. **Weak tests:** A repository may not provide enough evidence for a useful
   verified result.
5. **Semantic ambiguity:** Provider materials may omit business decisions.
6. **Repository complexity:** Wrappers, dynamic use, monorepos, generated code,
   and multiple SDK versions may dominate.
7. **Agent reliability and portability:** A selected adapter may drift,
   authenticate differently, omit needed structured evidence, or produce
   unrelated, brittle, or expensive candidates.
8. **Services economics:** Agent Runner, model, support, and operator COGS may
   prevent attractive margins, especially under the Lumyn-managed route.
9. **Sandbox mismatch:** Non-production behavior may not represent production.
10. **False confidence:** A weak verification stage may be mistaken for
    business-outcome proof.
11. **Provider visibility:** Aggregate status may not prove deployed retirement.
12. **Episodic demand:** Consequential migrations may not recur often enough
    for SaaS.
13. **Data leakage:** Model context, logs, traces, diffs, or summaries may
    expose consumer information.
14. **Event authenticity and replay:** A spoofed, stale, duplicated,
    superseded, or withdrawn provider event may trigger the wrong work.
15. **Ambient installation authority:** A durable installation may be treated
    as broader authority than the consumer intended.
16. **Status side channels:** Cohort reporting may expose consumer identity or
    integration posture, or may mistake silence and merge for stronger states.

---

## Falsification And Reframe Gates

Stop provider-specific control-plane or SaaS investment when:

- no qualified provider clears at least `$25,000` in non-refundable prepaid
  funds;
- the provider cannot recruit five Eligible Consumer Units across five
  distinct API Consumer Organizations;
- fewer than three impact scans complete within 14 days;
- fewer than two tested, reviewable outcomes result by the frozen deadline;
- Lumyn cannot open at least one tested draft PR with short-lived,
  least-privilege consumer authorization;
- no real participating repository completes a consumer-selected, qualified,
  independently verified agent-assisted run;
- no Consumer Maintainer accepts a verified outcome;
- guide or vendor tooling plus a generic agent performs materially
  equivalently end to end on maintainer effort, verified completion,
  correction risk, provider rollout evidence, and provider outcome;
- the frozen primary provider outcome misses its material threshold;
- model, support, or operator COGS misses the frozen contribution-margin or
  operator-automation threshold;
- most real repositories are blocked by unsupported semantics or inadequate
  verification;
- the provider cannot measure a meaningful retirement, support, or lead-time
  outcome.

A successful post-M6 qualified consumer-paid sprint may justify continued
engine development. A pre-M6 sprint justifies only workflow/problem learning.
Neither justifies claims about provider demand or recurring campaign software.
If consumer sprints repeatedly pay while providers refuse to prepay or recruit,
explicitly consider a consumer-services reframe rather than hiding that result.

Do not change cohorts, thresholds, correction definitions, or deadlines after
outcomes are visible. Abandonment and timeout count as failure.

---

## Non-Goals

The MVP does not:

1. scan repositories without explicit authorization;
2. claim coverage of every downstream integration;
3. grant API Providers access to consumer code or model context;
4. become a generic coding agent or dependency updater;
5. build a universal public-changelog monitor or neutral API registry;
6. infer missing business values or provider semantics;
7. support every language, SDK, package manager, or API style;
8. support auth, webhook/event, GraphQL, gRPC, generated-client, or production
   migrations;
9. execute arbitrary provider-supplied scripts;
10. require a hosted SaaS coordinator or dashboard;
11. require a long-lived GitHub token or broad organization-wide App grant;
12. claim manual bundle delivery as automated PR delivery;
13. build elaborate provider PKI, a universal event network, or cryptographic
    connection receipts before operational need is proven;
14. assume annual connected-repository pricing;
15. use production credentials or mutate production;
16. write to the default branch or auto-merge;
17. train shared models on private consumer artifacts;
18. present public fixtures or consumer-paid work as provider commercial
    validation.

---

## Definition Of MVP Success

Lumyn reaches first commercial MVP success only when:

- the applicable technical acceptance items pass with current evidence;
- one qualified API Provider clears at least `$25,000` in non-refundable
  prepaid funds;
- the provider identifies and invites five Eligible Consumer Units across five
  distinct API Consumer Organizations;
- three real repositories complete consumer-local impact analysis within 14
  days;
- at least two repositories produce tested, reviewable migration outcomes;
- at least one qualifying outcome contains an organically agent-assisted plan
  item on a consumer-selected qualified runner, passes independent exact-head
  verification, becomes a Lumyn-opened draft PR under short-lived authority,
  and sends the consented same-run status projection to the API Provider;
- at least one verified outcome is accepted or merged;
- Lumyn materially beats the frozen guide/codemod plus generic-agent baseline
  without worse correction, revert, or false-verification risk;
- the frozen primary provider outcome clears its material threshold;
- no candidate is falsely labeled verified;
- campaign revenue, model/tool COGS, Lumyn Operator time, Consumer Maintainer
  time, and the provider outcome are measured;
- the campaign receives an explicit `pass` or `fail` verdict by its frozen
  deadline; `reframe` may be recorded only as a post-failure disposition.

A second paid campaign or annual purchase order is evidence for repeatability
and a later SaaS decision. It is not required to close the first MVP
experiment.

Implementation sequence, ownership, and validation are defined in
`docs/product/plan.md`. The compiled Factory artifacts and repo-local operating
contracts are reconciled for v3 planning. This PRD, plan, and compilation do
not authorize runtime implementation or live product action; a later
implementation task requires explicit approval, and factoryd use additionally
requires its paused external profile/runtime gate to clear.
