# Lumyn MVP — Verified API Migration Execution

| Field | Value |
|---|---|
| Version | 3.0 |
| Status | Active v3 planning contract; no runtime implementation authorized by this document alone |
| Owner | Product and Engineering |
| Last Updated | 2026-07-24 |
| Primary Audience | Lumyn builders, API-provider design partners, API-consumer maintainers, and technical investors |
| Source Task Plan | `docs/product/plan.md` |
| Compiled Factory Plan | Repo-local v3 set regenerated; factoryd dispatch paused pending external Factory profile and runtime alignment |

---

## Purpose

This document is the product source of truth for the Lumyn MVP.

Lumyn is a services-led execution product for consequential API and SDK
migrations. It translates authoritative migration evidence into a
repository-specific, tested change that an API Consumer Organization can
review and adopt. The initial commercial motion is a provider-paid deprecation
or sunset campaign. A consumer-paid urgent upgrade sprint is a secondary proof
and revenue offer, not a second marketplace side.

The product promise is:

> When an API or SDK must change, Lumyn determines how an authorized consumer
> repository is affected, produces the smallest bounded migration it can
> justify, verifies the result against that repository's available evidence,
> and delivers a reviewable patch, branch, PR bundle, or optional draft PR.

Lumyn is not differentiated by opening a pull request. Generic coding agents
can already propose code changes. Lumyn is differentiated by combining
provider change semantics, consumer-local repository understanding, bounded
hybrid execution, deterministic verification, and proof-honest outcome
evidence across independent customer environments.

This version supersedes PRD v2.0's deterministic-first patch engine,
cryptographic two-sided activation, connection-receipt billing, and
draft-PR-centered product shape. Those mechanisms are not active MVP
requirements. Historical artifacts remain historical evidence and must not be
rewritten to claim v3 behavior.

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
official/public/provider change evidence
-> normalized migration pack
-> consumer-local integration graph
-> reviewable migration plan
-> deterministic transforms + bounded coding agent
-> baseline-aware compile/test/mock/sandbox repair loop
-> verified migration outcome
-> patch, branch, PR bundle, or optional draft PR
```

The operating model is provider-funded and consumer-controlled:

- The API Provider is the main economic buyer and campaign sponsor.
- The API Consumer Organization owns repository, execution, model-data,
  credential, disclosure, review, and merge authority.
- Provider payment never grants access to consumer source code.
- The coding agent runs through a consumer-local or explicitly
  consumer-approved execution boundary.
- The API Provider never receives raw code, diffs, prompts, responses, logs,
  traces, or credentials by default.
- Lumyn never writes to the default branch or auto-merges in the MVP.

The commercial outcome is a completed and evidenced migration, not a generated
PR. For an API Provider, value is a shorter legacy-version tail, fewer support
hours, lower migration risk, and clearer cohort readiness. For an API Consumer
Organization, value is less engineering effort and a smaller, better-evidenced
change to review.

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
authoritative migration evidence
+ consumer-local integration graph
+ deterministic and bounded-agent routing
+ baseline-aware verification and repair
+ consumer-controlled delivery
+ consented campaign outcome measurement
```

Every real migration must be compared with the status-quo baseline:

> official migration guide, vendor skill, or codemod plus a capable generic
> coding agent such as Codex or Copilot.

If Lumyn does not materially reduce maintainer effort or materially improve
verified completion without increasing correction risk, it has not
demonstrated product value.

The intended durable position is cross-company API migration campaign
execution, not code transformation alone.

---

## Roles And Terminology

Use these terms consistently. Avoid the bare word `provider` where it could
mean the API Provider or Model Provider.

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

A member of the Lumyn team who prepares migration packs, assists onboarding,
operates the first campaigns, or supports a consumer upgrade sprint. Lumyn
Operator time is measured as delivery COGS, separate from product development.

### Model Provider

The company or local runtime that supplies the model used by the coding agent.
It is never referred to as the API Provider. Cloud-model data handling is a
separate consumer decision and disclosure boundary.

### Agent Runner

The consumer-approved process or harness that selects model context, invokes
the Model Provider or local model, exposes bounded tools, applies budgets, and
records attempt provenance. The Agent Runner operates inside the consumer
execution boundary and cannot grant itself broader repository or network
authority.

### Migration Pack

A versioned, digest-bound, declarative artifact assembled from official
documentation, OpenAPI descriptions, SDK changes, examples, migration
guidance, provider clarification, and Lumyn analysis. It contains change
semantics, applicability, known transformations, unsupported conditions,
verification guidance, provenance, and confirmation status.

A migration pack is classified as:

- `public_derived`: derived from pinned public evidence without provider
  endorsement;
- `provider_confirmed`: reviewed and confirmed by an accountable Provider
  Operator.

A migration pack is untrusted input, cannot execute arbitrary code, and is not
proof that a candidate is correct.

### Integration Graph

A consumer-local model of where the selected API or SDK appears in the
repository, including dependency state, imports, aliases, wrappers, adapters,
call sites, request and response mappings, relevant configuration, tests,
mocks, cassettes, fixtures, and explicit uncertainty.

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

One migration pack, one API Provider, a bounded cohort of consenting API
Consumer Organizations, and measured migration outcomes. It is a services-led
engagement in the MVP, not a two-sided marketplace or cryptographic billing
network.

### Consumer Upgrade Sprint

A consumer-paid engagement for one consequential migration across one to
three repositories. It can prove the execution engine and generate revenue,
but it does not prove provider willingness to pay, provider-led distribution,
or recurring campaign demand.

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
> consenting customers complete and verify their repository-specific
> migrations without our team inspecting their proprietary code, and show us
> which integrations are actually ready.

### API Consumer Job

> When an API dependency changes, show exactly how my repository is affected
> and produce the smallest tested migration I can safely review and adopt.

Both roles are required for a provider-sponsored campaign. They are not both
economic buyers. A consumer-paid sprint may operate from public and official
evidence without making the API Provider a Lumyn customer.

---

## Initial Segment

The primary design-partner segment is an API-first B2B provider with:

- a hard deprecation, sunset, or migration deadline in the next 90 to 180
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

- migration-pack intake and provenance;
- semantic change normalization;
- a TypeScript integration graph;
- deterministic transforms;
- bounded-agent contracts plus a deterministic fake and status-quo baseline
  harness;
- verification and evidence contracts;
- status-quo baseline experiments.

Public inputs do not establish provider endorsement, prerelease authority,
demand, a reachable cohort, or repository-specific value. Lumyn may not claim
customer-specific value until it runs in at least one real, consenting
consumer repository.

Substantial provider-specific campaign automation remains gated by a prepaid
qualified campaign. A consumer-paid sprint may prove the engine but does not
close the provider commercial gate.

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
| Sunset-readiness sprint | API Provider | Migration-pack preparation, cohort qualification, baseline, and campaign plan | `$7.5k–$15k` |
| Provider migration campaign | API Provider | One migration across three to five consenting repositories | `$25k–$50k` |
| Additional completed repository | API Provider | One additional repository in the same supported campaign | `$2k–$5k` |
| Urgent consumer upgrade sprint | API Consumer Organization | One migration across one to three repositories | `$10k–$25k` |

The main company thesis is the provider-paid migration campaign. The
consumer-paid sprint is secondary revenue and engine proof.

For the first campaigns:

- Lumyn prepares the migration pack with provider input.
- The API Provider recruits the cohort.
- The API Consumer Organization runs Lumyn locally or in consumer-controlled
  CI.
- Lumyn-assisted model use is included in campaign COGS, or the consumer uses
  an approved BYOK/local model path.
- Provider reporting may be a manually prepared, consumer-consented aggregate.
- API Consumer Organizations pay no campaign seat fee.
- No annual connected-repository contract is assumed.

Annual platform pricing, hosted coordination, campaign subscriptions, or
connected-repository billing may be tested only after repeatable evidence from
a second paid campaign or an executed annual purchase order. A successful
one-off service engagement does not prove recurring SaaS demand.

---

## System Under Test

The system under test is:

```text
migration evidence
+ consumer-local integration graph
+ generated candidate
+ repository verification signals
```

The coding agent is an implementation mechanism, but its context selection,
tool use, attempts, cost, output, and repair behavior are part of the evaluated
system.

Every result preserves separate evidence about:

- what the migration pack says changed;
- what Lumyn found and did not find;
- what route each change used;
- what files and dependencies changed;
- which agent, model, tools, and budgets were used;
- what baseline and post-change commands ran;
- what behavior was observed;
- what remains unsupported, ambiguous, or unverified;
- what delivery and provider-reporting actions occurred.

---

## Product Principles

### 1. Two Principals, Two Authorities

The API Provider is authoritative about intended API or SDK semantics. The API
Consumer Organization is authoritative about repository, execution,
model-data, disclosure, and merge. Neither authority implies the other.

### 2. Read Before Write

Impact analysis is read-only. The Consumer Maintainer sees an impact report and
migration plan before authorizing mutation.

### 3. Provider Evidence Is Declarative

Migration packs may contain mappings, constraints, examples, and verification
references. They may not execute arbitrary provider-supplied scripts.

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

### 1. Qualify The Migration

Lumyn and the Provider Operator confirm the deadline, source and target
versions, expected cohort, evidence sources, buyer, operator, verification
signals, and campaign price. The first provider campaign does not proceed
without the commercial gates in this PRD.

### 2. Build The Migration Pack

`lumyn pack` normalizes pinned docs, OpenAPI or SDK diffs, examples, and
provider clarification. Conflicts and ambiguity remain visible. A
`public_derived` pack may drive engineering fixtures or a consumer sprint; a
provider-funded campaign requires accountable provider confirmation before a
real consumer mutation.

### 3. Establish The Consumer-Local Boundary

The Consumer Maintainer selects:

- repository and package root;
- readable and writable paths;
- approved commands;
- dependency-install policy;
- network and credential posture;
- model mode, Model Provider, data-egress policy, and context boundary;
- agent token, cost, time, attempt, file, and diff budgets;
- verification environments;
- output form and optional provider-visible status fields.

No selected scope implies another. Production credentials and production
mutation are prohibited.

### 4. Build The Integration Graph And Impact Report

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
- verification stages;
- delivery mode;
- residual risk and required human input.

Planning performs no write. Approval binds the exact pack, repository base,
scope, route, model policy, budgets, and commands.

### 6. Produce A Bounded Candidate

`lumyn apply` runs in an isolated worktree or consumer-approved equivalent:

1. deterministic codemods apply known safe mappings;
2. the bounded coding agent receives only the approved migration context and
   repository context;
3. a Consumer Maintainer or Lumyn Operator may instead make approved manual
   edits in that worktree, after which `lumyn candidate import --manual`
   validates and binds the diff to the exact repository base, candidate head,
   migration pack, plan, and route;
4. all writes remain within approved files and diff budgets;
5. every edit maps to a migration-pack item and repository evidence;
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
the Consumer Maintainer may separately authorize `lumyn repair`. The repair
authorization binds the failed candidate and evidence, exact repair intent,
remaining write and model-data permissions, and remaining time, token, cost,
attempt, file, and diff budgets. Repair creates a new candidate head,
invalidates prior verification evidence, and requires a fresh `lumyn verify`
run. It cannot expand scope. A non-diagnostic failure, missing business input,
or exhausted budget becomes `needs_input`, `blocked`, or `failed`.

Repository tests are untrusted code. They run in a consumer-approved execution
environment, without network or secrets by default. Dependency lifecycle
scripts, registry access, and sandbox credentials require distinct approval.

### 8. Export The Outcome

`lumyn export` produces the local evidence bundle plus one of:

- patch;
- local branch;
- PR-ready bundle for manual push;
- optional draft PR after short-lived GitHub authorization is implemented and
  separately approved.

A manually pushed branch or PR bundle is acceptable for the first services
pilot. It does not count as proof of automated PR delivery.

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
aggregates. Source code and raw execution evidence are not provider-visible by
default.

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

No roll-up may hide a weaker axis.

---

## Evidence Contract

Every outcome reports:

- API Provider and migration-pack identity;
- public-derived or provider-confirmed provenance;
- source and target versions and source digests;
- repository base and candidate head;
- analyzed paths, exclusions, and integration-graph uncertainty;
- affected, unaffected, unsupported, and uncertain call sites;
- deterministic, agent-assisted, and manual route by change item;
- changed files and dependency or lockfile delta;
- Agent Runner, adapter, model class, Model Provider or local mode, and
  version;
- model-data policy and context-policy digest;
- tool, command, path, network, credential, time, token, cost, attempt, file,
  and diff budgets;
- attempt and repair history without persisting raw prompts by default;
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

Evidence becomes stale when the pack, repository base or head, plan, model
policy, candidate, commands, or verification inputs change.

Raw prompts, model responses, source, diffs, logs, traces, and credentials are
consumer-private by default. Provider-visible summaries contain only
separately consented fields.

---

## Command Model

`lumyn` remains the primary local surface.

| Command | Purpose |
|---|---|
| `lumyn init` | Initialize repo-local Lumyn configuration |
| `lumyn check` | Validate configured sources and local prerequisites |
| `lumyn pack` | Build or validate a migration pack from pinned evidence |
| `lumyn impact` | Produce a read-only integration graph and impact report |
| `lumyn plan` | Produce a no-write routed migration plan |
| `lumyn apply` | Produce a bounded deterministic or agent-assisted candidate |
| `lumyn candidate import --manual` | Validate and bind an approved manual candidate to the base, pack, plan, and route |
| `lumyn verify` | Non-mutating baseline-aware candidate verification |
| `lumyn repair` | Separately authorized bounded repair that creates a new candidate for fresh verification |
| `lumyn export` | Export evidence plus patch, branch, or PR-ready bundle |
| `lumyn trace` | Render local evidence without changing state |
| `lumyn outcome record` | Record durable consumer acceptance, merge, closure, correction, or reversion evidence |
| `lumyn pr create --draft` | Optional later draft-PR delivery using short-lived authorization |

The command grammar may be implemented incrementally. Unimplemented commands
return a typed nonzero result and never a generic successful envelope.

Stable exit codes `0` through `9` remain reserved according to the existing
compatibility contract. Exit code `6` remains reserved. API Provider identity
uses `api_provider_id` or `change_authority`; Model Provider metadata uses
`model_provider_metadata`.

---

## Artifact Model

The MVP introduces only:

- `migration-pack`
- `integration-graph`
- `impact-report`
- `migration-plan`
- `candidate-manifest`
- `agent-attempt`
- `migration-verification`
- `export-result`
- `campaign-summary`
- `remediation-outcome`

An artifact becomes executable only when its schema, valid and invalid
fixtures, compatibility posture, and validator tests ship together.

Product artifacts are separated as:

```text
consumer checkout
  lumyn.yaml
  public or synthetic fixtures only

consumer-private state root outside the checkout
  packs/
  graphs/
  impacts/
  plans/
  candidates/
  attempts/
  verification/
  exports/

consumer-consented API Provider summary
  campaign status or aggregate only
```

The private root cannot resolve inside the checkout or public source
repository. Private artifacts have explicit retention and deletion policy.
The first services pilot may use documented operator-managed cleanup; a
productized deletion-receipt and orphan-recovery control plane is deferred.

---

## Trust, Authorization, And Data Boundaries

### Migration Evidence

- Migration packs are declarative, versioned, digest-bound, and inspectable.
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
- Remote branch and PR actions are separate from local mutation.
- Default-branch write and auto-merge are prohibited.

### Coding Agent And Model

- The Consumer Maintainer approves the agent runner, model mode, Model
  Provider, context boundary, data egress, retention posture, tools, and
  budgets.
- The agent receives no ambient secrets.
- Raw repository content is not transmitted to a cloud Model Provider unless
  the consumer explicitly approves that mode.
- BYOK and local-model use are supported as policy choices, not required public
  product integrations in the first pilot.
- The API Provider never receives agent context by virtue of sponsorship.

### Commands, Credentials, And Network

- Commands are allowlisted and displayed before execution.
- Repository commands run in a consumer-approved environment with explicit
  mounts, environment, process, timeout, and output limits.
- Tests run without network and secrets by default.
- Dependency lifecycle scripts and registry access require separate approval.
- Sandbox network, transmitted payload classes, and non-production credentials
  require separate approval.
- Production credentials and production mutation are prohibited.

### Sharing

The API Provider receives no raw consumer source or private evidence by
default. Consumer consent is required before sharing repository identity,
impact counts, candidate state, verification state, failure reason, merge
state, or retirement state.

Provider reporting may initially be manual. A hosted coordinator, signed
status exchange, and cryptographic billing receipts are deferred until demand
and operational need are proven.

---

## Functional Requirements

### FR1: Migration-Pack Intake

Lumyn must normalize pinned official or public migration evidence into a
versioned, provenance-visible, non-executable migration pack.

### FR2: Consumer-Local Integration Graph

Lumyn must identify dependency state, imports, wrappers, call sites, mappings,
relevant tests, and explicit uncertainty inside the authorized scope.

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
model policy, attempts, time, token, cost, file, and diff budgets.

### FR7: Repair Loop

Lumyn may use diagnostic verification failures for bounded repair without
expanding the approved plan. Exhausted or non-diagnostic cases stop.

### FR8: Baseline-Aware Verification

Lumyn must distinguish pre-existing failures from migration-attributable
failures and preserve the exact evidence boundary.

### FR9: Evidence-Bound Outcome

Every outcome must bind pack, base, candidate, plan, route, agent provenance,
commands, verification, residual risk, and artifact hashes.

### FR10: Multi-Form Export

Lumyn must export evidence plus a patch, branch, or PR-ready bundle without
requiring a hosted service.

### FR11: Optional Draft PR

Draft-PR delivery must use short-lived, least-privilege authorization and
remain distinct from code generation and verification.

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
plane except for explicitly disclosed Model Provider context and
consumer-consented provider summaries.

### NFR4: Bounded Autonomy

An agent cannot widen paths, tools, commands, network, credentials, model
policy, task intent, or budgets.

### NFR5: Least Privilege

Read, write, command, model-data, registry, sandbox, credential, remote branch,
PR, retention, deletion, and provider-reporting scopes are independent.

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

The v3 acceptance ledger contains 49 item-level closure units. Group headings,
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
   Provider versus Model Provider terminology is unambiguous.

### Migration Pack

1. `PACK-001`: Pinned public docs, OpenAPI descriptions, SDK releases, and
   migration guidance can produce a versioned `public_derived` migration pack
   with source digests and license/provenance records.
2. `PACK-002`: A pack records source and target versions, typed change
   semantics, applicability, known mappings, verification guidance, ambiguity,
   and unsupported cases.
3. `PACK-003`: Provider confirmation is independently recorded for a sponsored
   live campaign; public-derived evidence is never represented as provider
   endorsement, and provider input cannot execute code.
4. `PACK-004`: Conflicting, incomplete, stale, or semantically ambiguous
   evidence blocks affected routes instead of inviting the agent to infer
   provider intent.

### Impact And Integration Graph

1. `IMPACT-001`: Impact analysis is read-only and remains inside the exact
   authorized repository and package-root scope.
2. `IMPACT-002`: Lumyn detects the selected official npm SDK, installed
   version, package manifest, supported lockfile, direct imports, and aliases.
3. `IMPACT-003`: The integration graph represents statically traceable
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
   command, model policy, permission, budget, verification stage, and residual
   risk without modifying repository or Git state.
2. `PLAN-002`: Every affected item is explicitly classified as
   `deterministic`, `agent_assisted`, `manual`, `needs_input`, or `blocked`.
3. `PLAN-003`: Consumer approval binds the exact pack, repository base, paths,
   route, model-data policy, tools, commands, budgets, and verification intent;
   any change invalidates approval.

### Bounded Hybrid Execution

1. `AGENT-001`: Deterministic transforms produce byte-stable candidates for
   identical pinned inputs.
2. `AGENT-002`: Agent-assisted execution occurs only in a consumer-local or
   explicitly consumer-approved environment and discloses Model Provider,
   context, data-egress, and retention posture before execution.
3. `AGENT-003`: The agent enforces approved path, tool, command, network,
   credential, time, token, cost, attempt, file, and diff budgets.
4. `AGENT-004`: Every agent edit maps to a migration-pack item, integration
   graph evidence, and recorded rationale; unrelated edits fail.
5. `AGENT-005`: The repair loop performs no more than the approved attempts,
   cannot expand scope, and stops on non-diagnostic failure or exhausted
   budget.
6. `AGENT-006`: Missing business values, auth redesign, event semantics,
   production-only behavior, and ambiguous provider intent return
   `needs_input` or `blocked` without speculative completion.
7. `AGENT-007`: Agent evidence records runner, adapter, model class/version,
   model-data and context-policy digests, tools, commands, attempts, token and
   cost use, changed files, and human input without persisting secrets or raw
   prompts by default.

### Verification

1. `VER-001`: Lumyn records pre-existing dependency, compile, typecheck, and
   selected-test failures before candidate generation.
2. `VER-002`: Deterministic, agent-assisted, and imported manual candidates run
   the same dependency-integrity, compile/typecheck, and
   consumer-allowlisted test ladder in the approved execution environment.
3. `VER-003`: Static, repository, independent contract replay, exact-head
   replay, mock, and sandbox evidence use distinct canonical labels.
4. `VER-004`: Every `workflow_verified_*` result causally executes an approved
   entrypoint from the exact candidate head and records observed interaction
   and outcome evidence.
5. `VER-005`: Failed, missing, stale, out-of-boundary, unredactable, or
   inconclusive evidence cannot produce a verified label; repair failures
   remain visible.
6. `VER-006`: The frozen negative suite has zero false verified outcomes, and
   all evidence binds pack, plan, base, candidate, route, agent provenance,
   commands, environment, and artifact hashes.

### Export And Delivery

1. `EXP-001`: Lumyn exports complete local evidence plus a patch, local branch,
   or PR-ready bundle without requiring a hosted service.
2. `EXP-002`: A manual branch or PR bundle is labeled manual delivery and
   cannot count as automated PR delivery.
3. `EXP-003`: When the frozen campaign protocol selects automated delivery,
   it uses a short-lived, least-privilege token, writes only a non-default
   branch and draft PR, is idempotent, and never auto-merges. Otherwise the
   item closes `not_applicable` only with dated protocol evidence selecting
   manual patch, branch, or PR-bundle delivery and no automated-delivery claim.
4. `EXP-004`: Provider-visible status contains only consumer-consented fields
   and excludes raw source, diffs, prompts, responses, logs, traces, and
   credentials.

### Trust And Privacy

1. `TRUST-001`: Repository read, local write, commands, model-data egress,
   registry, sandbox, credentials, remote branch, PR, retention, deletion, and
   provider reporting are independently approved.
2. `TRUST-002`: Migration packs and provider artifacts cannot execute code;
   production credentials and production mutation are prohibited.
3. `TRUST-003`: Cloud-model context transfer, BYOK, and local-model posture are
   explicit; secrets and prohibited data are redacted or blocked before model
   or artifact persistence.
4. `TRUST-004`: Consumer-private artifacts remain outside the checkout and
   public repository, and provider sponsorship never grants source or model
   context access.

### Design-Partner Qualification

1. `DISC-001`: One qualified API Provider clears at least `$25,000` in
   non-refundable prepaid funds for a defined migration campaign and names an
   economic buyer, Provider Operator, hard deadline, source and target
   versions, and decision process.
2. `DISC-002`: Before campaign-specific migration execution, the API Provider
   identifies at least five reachable Eligible Consumer Units across five
   distinct API Consumer Organizations, each with an accountable maintainer
   and useful verification signals.
3. `DISC-003`: After the v3 privacy, model-data, authorization, and evidence
   contracts are approved, and before the first invitation, the parties freeze
   cohort, eligibility, price, evidence sources, privacy and model-data
   protocol, baseline method, material maintainer comparison threshold,
   correction rubric, campaign COGS boundary, measurement windows, one
   material provider-outcome metric and threshold, and absolute judgment
   deadline. A provider threshold is material only when the economic buyer
   records that meeting it would justify a retirement or paid-continuation
   decision and missing it makes the campaign fail.

### Provider Campaign Pilot

1. `PILOT-001`: At least five prequalified Eligible Consumer Units across five
   distinct API Consumer Organizations are invited without changing the
   frozen cohort after outcomes are visible.
2. `PILOT-002`: At least three repositories complete valid consumer-local
   impact scans within 14 calendar days of their invitations.
3. `PILOT-003`: At least two distinct repositories produce tested, reviewable
   migration outcomes with explicit evidence boundaries by the frozen pilot
   deadline.
4. `PILOT-004`: At least one verified migration outcome is accepted or merged
   by its Consumer Maintainer; closed, rejected, reverted, and corrected
   outcomes remain visible.
5. `PILOT-005`: Against the frozen guide/codemod plus generic-agent baseline,
   Lumyn reduces median Consumer Maintainer hands-on time by at least 30%
   without a worse substantive-correction, revert, or false-verification rate,
   or clears another equally material threshold frozen before execution.
6. `PILOT-006`: Actual cleared payment, model/tool COGS, Lumyn Operator hours,
   Consumer Maintainer time, support effort, and the frozen primary provider
   outcome are measured from source evidence, and the provider outcome clears
   its frozen material threshold.
7. `PILOT-007`: Every migration records accepted, merged, closed, reverted,
   corrected, blocked, and residual-risk outcomes through the frozen
   observation window.
8. `PILOT-008`: The campaign receives a pass or fail at the frozen judgment
   deadline. Failure to prepay, recruit, activate, outperform the generic
   baseline, clear the material provider-outcome threshold, or produce
   verified outcomes triggers a documented stop or reframe; the experiment
   cannot remain open indefinitely.

---

## Success Metrics

### North Star

`verified migration outcomes accepted per active provider campaign`

Each outcome carries its evidence boundary. A generated candidate or opened PR
does not count.

### Provider Outcome

Freeze exactly one primary provider outcome and material pass threshold before
the first invitation. The default is the share of the frozen target cohort
that is off the deprecated API or SDK version by the migration deadline.
Provider support hours per accepted migration or invitation-to-accepted-
migration lead time may be selected instead when the economic buyer records
that the threshold is material to retirement or paid continuation. The cohort
denominator and legacy-version retirement share remain reported whenever they
are measurable.

### Funnel

```text
eligible
-> invited
-> consented
-> impacted
-> migration planned
-> candidate produced
-> repo verified
-> consumer accepted
-> merged or deployed
-> legacy version retired
```

Report conversion and elapsed time between every stage.

### Technical Quality

- Impact recall and precision by supported route.
- False `unaffected` rate.
- Deterministic transform repeatability.
- Agent scope-violation rate.
- Agent-assisted completion rate.
- Human-input and blocked rate.
- Verification pass and false-verification rate.
- Substantive correction and revert rate.
- Unrelated-edit rate.
- Repair attempts per accepted outcome.
- Median impact, generation, verification, and export time.

### Economics And Effort

- Provider campaign revenue and actual payment.
- Model, tool, infrastructure, and Lumyn Operator COGS.
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

The initial provider motion is direct and services-led:

1. Sell an active sunset or deprecation campaign to an accountable provider
   buyer.
2. Prepay and qualify the cohort before substantial campaign-specific
   automation.
3. Lumyn prepares the migration pack with the Provider Operator.
4. The API Provider recruits participating Consumer Maintainers.
5. Consumers run impact, planning, generation, and verification locally or in
   consumer-controlled CI.
6. Lumyn exports consumer-controlled outcomes and manually aggregates
   consented campaign status.
7. Productize repeated operational friction only after the pilot.

The secondary motion is a paid urgent consumer upgrade sprint. It is useful
for engine proof and cash but does not validate provider-led distribution.

The MVP does not require hosted SaaS, a long-lived GitHub installation,
automatic PR delivery, a campaign dashboard, an annual connected-repository
program, or public OSS distribution.

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

- migration-pack normalization;
- TypeScript integration graph and impact analysis;
- deterministic migration transforms;
- coding-agent adapter and consumer-approved execution boundary;
- bounded agent and repair budgets;
- hybrid migration planning and routing;
- repository verification orchestration;
- replay, mock, or sandbox verification runtime;
- agent provenance and cost evidence;
- patch, branch, and PR-bundle export;
- optional GitHub draft-PR delivery;
- consented provider campaign reporting;
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

1. **Generic-agent substitution:** Guide plus Codex or Copilot may perform
   equivalently.
2. **Provider activation:** A willing provider may fail to recruit consumers.
3. **Consumer trust:** Repository or model-data policies may block execution.
4. **Weak tests:** A repository may not provide enough evidence for a useful
   verified result.
5. **Semantic ambiguity:** Provider materials may omit business decisions.
6. **Repository complexity:** Wrappers, dynamic use, monorepos, generated code,
   and multiple SDK versions may dominate.
7. **Agent reliability:** The bounded agent may produce unrelated, brittle, or
   expensive candidates.
8. **Services economics:** Model, support, and operator COGS may prevent
   attractive margins.
9. **Sandbox mismatch:** Non-production behavior may not represent production.
10. **False confidence:** A weak verification stage may be mistaken for
    business-outcome proof.
11. **Provider visibility:** Aggregate status may not prove deployed retirement.
12. **Episodic demand:** Consequential migrations may not recur often enough
    for SaaS.
13. **Data leakage:** Model context, logs, traces, diffs, or summaries may
    expose consumer information.

---

## Falsification And Reframe Gates

Stop provider-specific control-plane or SaaS investment when:

- no qualified provider clears at least `$25,000` in non-refundable prepaid
  funds;
- the provider cannot recruit five Eligible Consumer Units across five
  distinct API Consumer Organizations;
- fewer than three impact scans complete within 14 days;
- fewer than two tested, reviewable outcomes result by the frozen deadline;
- no Consumer Maintainer accepts a verified outcome;
- guide or vendor tooling plus a generic agent performs materially
  equivalently;
- the frozen primary provider outcome misses its material threshold;
- model, support, or operator COGS makes the campaign uneconomic;
- most real repositories are blocked by unsupported semantics or inadequate
  verification;
- the provider cannot measure a meaningful retirement, support, or lead-time
  outcome.

A successful consumer-paid sprint may justify continued engine development. It
does not justify claims about provider demand or recurring campaign software.
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
5. build a universal API-change monitor;
6. infer missing business values or provider semantics;
7. support every language, SDK, package manager, or API style;
8. support auth, webhook/event, GraphQL, gRPC, generated-client, or production
   migrations;
9. execute arbitrary provider-supplied scripts;
10. require a hosted SaaS coordinator or dashboard;
11. require a long-lived GitHub installation;
12. claim manual bundle delivery as automated PR delivery;
13. build provider PKI, signed invitations, status channels, or cryptographic
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
- at least one verified outcome is accepted or merged;
- Lumyn materially beats the frozen guide/codemod plus generic-agent baseline
  without worse correction, revert, or false-verification risk;
- the frozen primary provider outcome clears its material threshold;
- no candidate is falsely labeled verified;
- campaign revenue, model/tool COGS, Lumyn Operator time, Consumer Maintainer
  time, and the provider outcome are measured;
- the campaign receives an explicit pass, fail, or reframe decision by its
  frozen deadline.

A second paid campaign or annual purchase order is evidence for repeatability
and a later SaaS decision. It is not required to close the first MVP
experiment.

Implementation sequence, ownership, and validation are defined in
`docs/product/plan.md`. The compiled Factory artifacts and repo-local operating
contracts are reconciled for v3 planning. This PRD, plan, and compilation do
not authorize runtime implementation or live product action; a later
implementation task requires explicit approval, and factoryd use additionally
requires its paused external profile/runtime gate to clear.
