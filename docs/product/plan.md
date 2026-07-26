# PLAN Lumyn v3.1: Provider-Originated API Update Delivery

**Date:** 2026-07-25
**Source of truth:** `docs/product/prd.md`
**Status:** Active v3.1 planning contract. A v3.1 source update has landed in
the external Factory profile, but this compiled control generation still
records Factory profile/runtime compatibility as unqualified. factoryd task
dispatch remains paused until the complete control set is regenerated and its
bundle/runtime and exact mission are qualified and explicitly unpaused.
Separately approved attended tasks remain available
**Scope:** Build and validate one provider-originated TypeScript/Node API
sunset campaign in which a provider publishes one confirmed change, authorized
consumer installations produce repository-specific verified updates, and
Lumyn opens tested draft PRs under consumer-controlled policy. Services are
the initial onboarding and GTM motion, not the product identity.

---

## Global Decisions (Locked)

1. The product thesis is a trusted provider-to-consumer API update-delivery
   channel; the main commercial thesis is a provider-paid API or SDK
   deprecation campaign.
2. A consumer-paid urgent upgrade sprint is secondary revenue. Before M6 it
   produces paid workflow/problem evidence only; after a qualified M6 run it
   may prove engine value. It never validates provider demand or provider-led
   distribution.
3. The API Provider is the campaign buyer and sponsor. The API Consumer
   Organization is the repository, Agent Runner selection, execution,
   model-data, credential, disclosure, and merge authority.
4. Provider payment never grants consumer repository or model-context access.
5. Initial onboarding is services-assisted through a local CLI and
   consumer-controlled execution environment, not hosted SaaS.
6. The implementation core remains Go `1.26.5` at module
   `github.com/Clyra-AI/lumyn`.
7. The first target is one official TypeScript/Node npm SDK, one migration
   from a defined source version to a defined target version, and one
   explicitly selected package root.
8. Pinned public docs, OpenAPI descriptions, SDK releases, migration guides,
   and licensed fixtures are sufficient for M0 through M4 engineering.
9. Public evidence does not prove provider endorsement, demand, or
   repository-specific value.
10. Repository-specific value requires a real consenting consumer repository.
11. The consumer-local repository impact inventory, represented internally as
    an integration graph, includes dependency state, imports, wrappers,
    adapters, call sites, mappings, relevant configuration, tests, mocks,
    cassettes, fixtures, exclusions, and uncertainty.
12. Every affected item routes to `deterministic`, `agent_assisted`, `manual`,
    `needs_input`, or `blocked`.
13. Known safe transformations are deterministic. Repository-specific
    adaptation may use a bounded coding agent from the first MVP.
14. Agent output is an untrusted candidate. It cannot self-verify.
15. `agent_execution_policy` defaults to `disabled`. When it is
    `configured`, the Consumer Maintainer selects one qualified exact Agent
    Runner adapter/version/executable digest and approves its auth mode and
    entitlement class, Agent Runner Vendor, actual Model Provider or local
    route, context boundary, data egress, native agent configuration, tools,
    commands, Agent Runner/model network and credentials, time, token, cost,
    attempt, file, and diff budgets.
16. The default maximum is three agent or repair attempts. A higher limit
    requires new approval.
17. API Provider evidence and repository content are untrusted data, not agent
    instructions. Embedded prompt injection cannot widen the approved plan.
18. Determinism governs pinned inputs, deterministic transforms, routing,
    budgets, verification, evidence, and status. Agent source output need not
    be byte-identical.
19. Impact analysis and migration planning are read-only.
20. Candidate generation runs in an isolated worktree or consumer-approved
    equivalent and stays within explicit path and diff budgets.
21. Pre-existing repository failures are measured before mutation.
22. Deterministic and agent-assisted candidates pass the same verification
    ladder and proof labels.
23. The complete product loop is provider event to installed consumer policy,
    repository-specific impact, verified candidate, tested draft PR, and
    consented rollout status; code generation alone is not the product.
24. `lumyn export` supports evidence plus patch, local branch, or PR-ready
    bundle as a safe assisted fallback.
25. At least one tested Lumyn-opened draft PR using short-lived,
    least-privilege authorization is mandatory first-campaign product proof;
    manual delivery cannot substitute for it.
26. Lumyn never writes to the default branch or auto-merges.
27. The API Provider receives only event-bound, consumer-consented status or
    aggregates. Silence is unknown, merge does not imply retirement, and raw
    source, diffs, prompts, responses, logs, traces, and credentials are not
    API Provider-visible.
28. `agent_execution_policy` defaults to `disabled`; notify-only, scan-only,
    and deterministic-only installations require no Agent Runner credential.
    If a plan needs `agent_assisted`, Lumyn pauses for explicit
    `configured` policy and authorization. For configured execution,
    `consumer_managed` is the default funding mode: the consumer owns the
    qualifying agent account, subscription, BYOK credential, or local runtime
    and third-party usage billing. Optional
    `provider_sponsored_lumyn_managed` execution makes approved agent/model
    usage Lumyn campaign COGS without giving the API Provider code, context,
    session, or credential access.
29. Production credentials and production mutation are prohibited.
30. Repository tests run without network and secrets by default. Registry,
    dependency lifecycle scripts, sandbox network, payload disclosure, and
    sandbox credentials are independently approved.
31. The initial provider campaign price hypothesis is `$25k–$50k`; one
    provider must clear at least `$25k` in non-refundable prepaid funds before
    M5 begins.
32. The provider must identify five reachable Eligible Consumer Units across
    five distinct API Consumer Organizations and commit to lead their
    distribution and onboarding before M5 begins.
33. The campaign must be compared with migration guide, vendor codemod or
    skill, plus the same capable coding engine under matched runner, model,
    auth, funding, context, tool, command, attempt, token, time, and cost
    controls.
    Unmatched engine comparisons are descriptive, not causal.
34. No annual connected-repository contract, hosted coordinator, universal
    registry, public-changelog monitor, elaborate provider PKI, or
    connection-receipt billing system is an MVP requirement.
35. Public fixtures prove engineering behavior only. Pre-M6 consumer-paid
    work proves paid workflow/problem evidence; qualified post-M6 work may
    prove engine value. Neither proves the provider campaign thesis.
36. Historical planning artifacts remain immutable. This rebaseline compiles
    v3 into `.factory/artifacts/prd-to-plan/lumyn-migration-mvp/`; the compiled
    set remains non-dispatchable while its external Factory profile or
    factoryd compatibility lock disagrees.
37. This PRD, plan, operating-doc, and compiled-control rebaseline is separate
    from M0 runtime implementation. No code task is authorized merely because
    the planning artifacts are current.
38. The campaign must clear a material provider-outcome threshold frozen
    before invitations, in addition to beating the maintainer-time and quality
    baseline.
39. A Provider Change Contract is confirmed once and reused for the exact
    cohort and version. Each provider event references it and cannot execute
    code or grant repository authority.
40. A revocable Consumer Installation binds provider or channel, repository,
    package root, audience or version selectors, allowed actions,
    `agent_execution_policy`, model-data, GitHub, retention, and disclosure
    policy. A configured agent policy additionally binds the exact Agent Runner
    route and funding/credential/billing ownership. Each event-specific
    authorization may narrow but never widen that policy.
41. The generic coding agent is a replaceable implementation adapter and the
    explicit status-quo comparator, not Lumyn's proprietary moat.
42. Launch Agent Runner targets are Codex and Claude Code. Each pinned adapter
    version must pass the same conformance suite and an explicitly approved
    live canary before advertised support. Cursor remains deferred behind that
    identical gate.
43. Every attempt resolves an approved executable by canonical path and digest,
    then starts a clean ephemeral Agent Runner session from neutral home and
    config roots. Lumyn never accepts repository-local PATH shadowing or resumes
    a personal or unrelated conversation.
44. Native user/project rules, memories, plugins, and configuration are ignored
    unless explicitly selected, digest-bound, and treated as untrusted context
    that cannot widen Lumyn authority.
45. Lumyn never silently changes Agent Runner adapter/version, Model Provider,
    model, endpoint, credential owner, or usage-billing owner. An unavailable
    or unqualified route blocks agent execution; a separately valid
    deterministic route may still proceed.
46. Agent Runner Vendor and Model Provider are separate recorded roles even
    when one company supplies both or the runner brokers the model route.
47. Agent-reported tests are generation evidence only. Lumyn independently
    executes the approved verification commands from the exact candidate head
    outside the agent session.
48. The first MVP provides controlled installation-time selection among
    qualified launch adapters, not an arbitrary adapter marketplace or dynamic
    mid-run switching.

---

## Current Baseline (Observed)

Implemented:

- Go CLI, configuration, result envelope, and exit-code foundation.
- `lumyn init` and `lumyn check`.
- OpenAPI and local-doc parsing, fingerprints, structured refs, deprecation
  findings, and concrete source locations.
- Executable schemas for workflows, evidence, cassettes, traces, proof,
  boundaries, redaction, and command results.
- CI, coverage, CodeQL, branch policy, CODEOWNERS, required checks, review,
  Factory planning, commit/push, and post-merge governance.

Not implemented:

- Provider Change Contract and event normalization;
- consumer installation and event-specific authorization;
- API or SDK semantic diffing;
- TypeScript repository impact inventory;
- repository impact analysis;
- routed migration planning;
- deterministic migration transforms;
- customer-selected Agent Runner adapter, conformance, or execution policy;
- bounded agent and agent-assisted repair loop;
- Agent Runner/model provenance, funding, credential/billing, context-policy,
  token, and cost evidence;
- repository verification orchestration;
- replay, mock, or live sandbox verification runtime;
- evidence plus patch, branch, or PR-bundle export;
- short-lived GitHub draft-PR delivery;
- event-bound consented provider status projections;
- migration outcome ingestion.

M0 dispatch baseline:

- `record`, `verify`, `trace`, `demo`, `share`, and `eval` are recognized by
  the command dispatcher even though they have no implementation. M0 replaces
  their baseline generic pass with typed `command_not_implemented` failures.
- Result contract v1.0 retains the legacy serialized `provider_metadata` key
  for Model Provider metadata. M0 makes new results explicit with
  `semantic_role: model_provider`, keeps legacy payloads valid, and stops
  setting evaluation metadata on `init` and `check`.
- The repo-local v3 contract, task packets, and acceptance ledger are
  rebaselined. A v3.1 source update has landed in Factory
  `profiles/lumyn.yaml`, but this compiled generation still records Factory
  profile/runtime compatibility as unqualified. Regeneration, factoryd
  bundle/runtime qualification, and unpause remain incomplete.
- Historical task evidence proves only the exact foundation it recorded.

No line in this plan represents an unimplemented surface as shipped.

---

## Acceptance Ownership

The v3 PRD defines 53 item-level acceptance units.

| Category | IDs | Count | Primary owner |
|---|---|---:|---|
| Retained foundation | `BASE-001`–`BASE-005` | 5 | M0 |
| Provider Change Contract | `PACK-001`–`PACK-004` | 4 | M3 |
| Provider Change Event | `EVENT-001`–`EVENT-002` | 2 | M3 |
| Consumer Installation | `INSTALL-001`–`INSTALL-002` | 2 | M5 |
| Impact and repository inventory | `IMPACT-001`–`IMPACT-005` | 5 | M4 |
| Plan and routing | `PLAN-001`–`PLAN-003` | 3 | M5 |
| Bounded hybrid execution | `AGENT-001`–`AGENT-007` | 7 | M6 |
| Verification | `VER-001`–`VER-006` | 6 | M7, with M8 conditional |
| Export and delivery | `EXP-001`–`EXP-004` | 4 | M9 |
| Trust and privacy | `TRUST-001`–`TRUST-004` | 4 | M2 |
| Design-partner qualification | `DISC-001`–`DISC-003` | 3 | M2.5 |
| Provider campaign pilot | `PILOT-001`–`PILOT-008` | 8 | M10 |
| **Total** |  | **53** |  |

M1 creates benchmark prerequisites for `PACK-001`, `EVENT-*`, `INSTALL-*`,
`IMPACT-005`, `AGENT-*`, `VER-006`, and `PILOT-005`; it does not close those
items by itself. M2 defines executable contract prerequisites for `EVENT-*`
and `INSTALL-*`; M3 and M5 own their runtime closure respectively. M8 is
conditional and strengthens verification when sandbox use is selected. An
acceptance item may have supporting tasks, but it has one primary closure
owner.

---

## Exit Criteria

### Technical MVP

The execution engine is technically complete when:

- every applicable `BASE`, `PACK`, `EVENT`, `INSTALL`, `IMPACT`, `PLAN`,
  `AGENT`, `VER`, `EXP`, and `TRUST` item has current evidence;
- public and held-out benchmarks cover deterministic, agent-assisted, and
  blocked routes;
- unimplemented commands fail closed;
- no negative case receives a false verified label;
- agent execution stays inside every declared boundary;
- a full offline canary reaches a verified local export without hidden
  network, credentials, API Provider reporting, or GitHub writes;
- the full local validation gate and required GitHub checks are green.

### First Commercial MVP

Commercial completion additionally requires all `DISC` and `PILOT` items:

- one provider clears at least `$25,000` in non-refundable prepaid funds;
- the provider recruits five Eligible Consumer Units across five distinct API
  Consumer Organizations;
- three consumer-local impact scans complete within 14 days;
- two tested, reviewable outcomes are produced;
- at least one qualifying outcome contains an organically agent-assisted plan
  item on a consumer-selected qualified runner, passes independent exact-head
  verification, becomes a Lumyn-opened draft PR, and sends the consented
  same-run status projection to the API Provider;
- one verified outcome is accepted or merged;
- Lumyn materially beats the frozen guide or codemod plus generic-agent
  baseline;
- the frozen primary provider outcome clears its material threshold;
- actual payment, COGS, effort, corrections, and the provider outcome are
  measured;
- the campaign receives a `pass` or `fail` verdict at its frozen deadline;
  `reframe` is allowed only as a recorded post-failure disposition.

A second campaign or annual order is a scale signal, not a first-MVP exit
criterion.

---

## Public API And Contract Map

| Surface | Status | Contract |
|---|---|---|
| `lumyn init` | Retained | Repo-local configuration |
| `lumyn check` | Retained and reframed | Non-mutating source, repository, verification, and conditional Agent Runner readiness preflight |
| `lumyn pack` | New | Build or validate a versioned Provider Change Contract |
| `lumyn install` | New | Bind provider channel, repository, action, optional Agent Runner/model and funding, GitHub, and disclosure policy |
| `lumyn update --event` | New composed surface | Process one event through installed policy, stopping at the selected action mode |
| `lumyn impact` | New | Read-only repository impact inventory |
| `lumyn plan` | New | No-write routed migration plan |
| `lumyn apply` | New | Bounded deterministic and agent-assisted candidate generation |
| `lumyn candidate import --manual` | New | Bind an approved manual diff to the exact base, pack, plan, and route |
| `lumyn verify` | New/reused semantics | Non-mutating baseline-aware candidate verification |
| `lumyn repair` | New | Separately authorized agent-assisted repair that creates a new candidate |
| `lumyn export` | New | Evidence plus patch, branch, or PR-ready bundle fallback |
| `lumyn trace` | New runtime over retained schema | Local evidence rendering |
| `lumyn outcome record` | New | Append an authorized, exact-candidate-bound adoption or remediation outcome |
| `lumyn pr create --draft` | New | Required short-lived-token pilot delivery proof |
| `lumyn.command_result` `1.0` | Compatibility surface | Existing envelope with corrected terminology |
| v3 migration artifacts | New | Versioned schemas and typed results |
| workflow/evidence schemas | Retained | Verification substrate |

Compatibility rules:

- Exit codes `0` through `9` remain stable.
- Exit code `6` remains reserved.
- Existing bare `provider_metadata` continues to mean Model Provider metadata
  during the compatibility window.
- New v1.0 writers emit optional `semantic_role: model_provider`; historical
  v1.0 payloads without it remain valid. V1.0 schemas preserve their prior
  unknown-extension openness; an extension value cannot override the key's
  legacy Model Provider meaning, and current writers emit no other role.
- API Provider identity uses `api_provider_id` and `change_authority`.
- New Model Provider fields use `model_provider_metadata`.
- Renaming the v1.0 key requires a v2.0 schema, frozen v1.0 validators,
  version-discriminated readers, v2.0 invalid mixed-role fixtures, and a new
  provenance-linked artifact rather than an in-place historical rewrite.
- Persisted schema changes require a new version, valid and invalid fixtures,
  compatibility notes, and migration behavior.
- Unimplemented commands return typed nonzero errors.

Deferred command families:

- provider enrollment beyond one configured channel and trust refresh;
- signed campaign invitation;
- consumer and provider receipt exchange;
- cryptographic connected-repository billing;
- hosted coordinator or dashboard;
- automatic merge.

---

## Docs And Distribution Readiness

The first screen of `README.md` must eventually communicate:

```text
migration evidence
-> Provider Change Contract and event
-> installed consumer policy
-> consumer-local repository impact inventory
-> deterministic or bounded-agent candidate
-> repository/workflow verification
-> tested draft PR
-> consented provider rollout status
```

Before implementation dispatch, docs must:

- define the provider and consumer jobs;
- distinguish the API Provider, Agent Runner Vendor, and Model Provider;
- explain the cloud-model, BYOK, and local-model data boundaries;
- state that provider sponsorship does not grant source access;
- state actual implementation status;
- distinguish benchmark proof, consumer-sprint proof, and provider commercial
  proof;
- document agent budgets, verification labels, and fail-closed behavior;
- document patch, branch, and PR-bundle fallback plus mandatory pilot
  draft-PR proof honestly;
- avoid advertising deferred PKI, receipt, dashboard, or annual billing
  surfaces.

The current repository and design-partner distribution are not represented as
open source until an explicit license, security, contribution, support,
release-integrity, and vulnerability-response gate closes.

---

## Validation And Test Matrix

Every first-party code task runs the repo-local fast tests, coverage, contract
tests, and `make prepush-full`. Risk-sensitive tasks add the relevant
architecture, security, holdout, integration, or live-environment review.

| Task | Fast/core | Scenario and acceptance | Risk lane | External dependency |
|---|---|---|---|---|
| M0 | CLI, schema, full gate | `BASE` | compatibility, CodeQL | none |
| M1 | fixture and contract | benchmark prerequisites | provenance, holdout | deterministic fake only; live adapters deferred to M6/M10 |
| M2 | schema and contract | `TRUST`; `EVENT`/`INSTALL` prerequisites | security, privacy, architecture | none |
| M2.5 | evidence validation | `DISC` | product, privacy, commercial | provider and cohort |
| M3 | unit and integration | `PACK`, `EVENT` | parser, transport, and provenance | pinned provider channel |
| M4 | parser and scenario | `IMPACT` | holdout accuracy | none |
| M5 | unit and contract | `INSTALL`, `PLAN` | write/model-policy boundary | none |
| M6 | unit, scenario, mutation | `AGENT` | runner conformance, prompt injection, scope, cost | selected qualified Agent Runner |
| M7 | unit and integration | `VER` | untrusted commands, false green | offline canary |
| M8 | mock and optional live | conditional `VER` | credentials, network, privacy | approved sandbox |
| M9 | unit and integration | `EXP` | GitHub and disclosure | approved short-lived GitHub grant |
| M10 | evidence validation | `PILOT` | product, privacy, economics | provider and consumers |

Implementation workers must not inspect held-out answers or self-attest
commercial evidence. The detailed lifecycle worker chain remains governed by
the repo-local Factory contract and is not a customer product flow.

---

## Planning Rebaseline Gate

This documentation and control-set change completes the repo-local v3.1 planning
rebaseline. It regenerates the 53-item acceptance ledger, mapping, execution
plan, task packets, validation contract, and closure map without changing
runtime product behavior. Historical v2 evidence remains immutable.

A v3.1 source update has landed in the external Factory profile, but the active
compiled control set remains the authority for factoryd and still records
Factory profile/runtime compatibility as unqualified. The checked-in factoryd
mission therefore stays paused until the complete control set is regenerated,
its bundle/runtime qualifies the exact active mission, and a bounded task is
explicitly unpaused. That pause blocks factoryd dispatch only. A separately
approved attended task may use the same packet and full lifecycle gates while
recording current external dependencies separately; it does not claim factoryd
readiness. No implementation task is approved by this rebaseline itself.

---

## Epic 0 — Fail-Closed Compatibility

### M0: Correct command and result foundations

**Priority:** P0
**Risk class:** Medium
**Blocked by:** none
**Pre-dispatch gate:** approved v3 compiled control set, aligned external
Factory profile, and an explicitly authorized implementation task; factoryd
profile/runtime compatibility is additionally required only when factoryd is
the selected execution path
**Primary acceptance IDs:** `BASE-001`–`BASE-005`

#### Goal

Make the existing runtime foundation honest without reinterpreting historical
evidence.

#### Tasks

- Make every recognized but unimplemented command return a typed nonzero
  result, or remove it from the command registry.
- Stop setting evaluation-oriented metadata on `init` and `check`.
- Separate API Provider, Agent Runner Vendor, and Model Provider terminology in
  docs, results, fixtures, and compatibility notes.
- Preserve exit-code compatibility, including reserved exit code `6`.
- Define the versioned migration path for command-result and evidence schemas.
- Preserve the implemented OpenAPI/docs parser and retained verification
  schemas.

#### Expected repo areas

- `cmd/lumyn/`
- `internal/result/`
- `internal/exitcode/`
- `schemas/`
- `tests/`
- `docs/`
- `CHANGELOG.md`

#### Required tests

- Red-first tests for every recognized unimplemented command.
- JSON envelope tests for `init`, `check`, unknown, and unimplemented commands.
- API Provider versus Agent Runner Vendor versus Model Provider terminology
  fixtures.
- Exit-code stability tests.
- Schema compatibility tests.
- Repo-pack validation proving M0 was dispatched from the approved v3
  generation.

#### Completion criteria

- No unimplemented command reports success.
- Existing valid foundation artifacts remain valid or have an explicit
  migration.
- M0 evidence binds the approved v3 task packet and exact repository head.
- No v2 deterministic/PKI-first task was selected as v3 implementation.

#### Stop conditions

- A compatibility change would silently reinterpret historical evidence.
- The v3 compiled set is stale or contradictory.
- factoryd is selected while its external Factory profile or runtime
  compatibility posture remains unqualified.

**ADR impact:** compatibility and product-direction ADR required before
implementation.
**Changelog impact:** required.
**Cost/performance:** low.

---

## Epic 1 — Migration Corpus And Status-Quo Baseline

### M1: Build deterministic, agent-assisted, blocked, and generic-agent benchmarks

**Priority:** P0
**Risk class:** Medium
**Blocked by:** M0, M2
**Supports:** `PACK-001`, `EVENT-001`–`EVENT-002`,
`INSTALL-001`–`INSTALL-002`, `IMPACT-005`, `AGENT-001`–`AGENT-007`,
`VER-006`, `PILOT-005`

#### Goal

Create a frozen, provenance-backed benchmark and one thin local walking
skeleton that test Lumyn's actual v3.1 claim against the status quo before
deep implementation.

#### Tasks

- Create at least three pinned API or SDK migration scenarios covering:
  - method or operation rename;
  - request-property rename or relocation;
  - response-property rename or relocation.
- Create at least three agent-assisted scenarios involving wrappers, adapters,
  signature or type adaptation, and related test repair without new business
  decisions.
- Create blocked scenarios for missing business values, auth redesign, event
  semantics, ambiguous evidence, production-only behavior, scope escape, and
  inadequate verification.
- Include at least six visible TypeScript consumer fixtures across direct and
  wrapper-heavy repository shapes.
- Annotate affected, unaffected, uncertain, unsupported, expected-route,
  expected-edit, expected-command, and expected-verification ground truth.
- Record source digest, license, attribution, redistribution, and
  public-derived versus synthetic provenance.
- Freeze a private held-out suite before M4, M6, or M7 scoring. Implementation
  workers receive only opaque case identifiers and aggregate results.
- Define a fair status-quo baseline using the same repository snapshot,
  authoritative migration evidence, selected Agent Runner adapter, version,
  and executable, actual Model Provider/model version, auth/entitlement and
  execution-funding route, credential and usage-billing owners, context-access
  ceiling, tools, commands, engineer role, attempt/token/time/cost budgets, and
  verification opportunity: migration guide, vendor codemod or skill, plus the
  same capable coding engine without Lumyn's orchestration, impact/routing,
  boundary enforcement, independent verification, delivery, or status loop.
  Any unmatched engine comparison is descriptive and cannot prove causal
  advantage.
- Freeze the primary comparison metric. Default: at least 30% lower median
  Consumer Maintainer hands-on time with no worse substantive-correction,
  revert, or false-verification rate.
- Record model/tool versions, prompts or instruction policy, attempts, token
  and cost budgets for the baseline without committing consumer-private data.
- Using public or synthetic fixtures, prove the thin sequence `event ->
  installation -> impact -> plan -> candidate -> verify -> PR bundle` through
  the common Agent Runner contract, with a deterministic fake for the agent
  and every external write. This walking skeleton proves contract composition
  only; it cannot close adapter qualification, automated GitHub delivery,
  real-repository, or commercial acceptance.
- Keep M1 offline and deterministic-fake-only. It must not launch Codex,
  Claude Code, Cursor, or any live model, consume Agent Runner/model
  credentials, or receive external network grants. Freeze the common harness
  and matched-baseline protocol here; M6 owns qualified adapter conformance and
  live canaries, and M10 owns the matched live generic-agent comparison.

#### Expected repo areas

- `examples/migration-packs/`
- `examples/consumer-repos/`
- `examples/integration-graphs/`
- `examples/candidates/`
- `examples/negative/`
- `examples/holdout-manifest.json`
- `tests/`
- `docs/`

#### Required tests

- Fixture digest and manifest validation.
- License and provenance completeness.
- Ground-truth completeness.
- Negative fixture rejection.
- Holdout access-isolation tests.
- Mutation tests proving changed ground truth or suite binding invalidates
  results.
- Baseline reproducibility checks.
- Walking-skeleton contract, event replay, and fake-delivery tests.
- Adapter-neutral fake conformance, clean-session, normalized output,
  cancellation, and no-silent-fallback tests.
- Negative tests proving M1 cannot launch a live Agent Runner/model or accept
  Agent Runner/model credentials or network grants.

#### Completion criteria

- Deterministic, agent-assisted, and blocked routes all have visible and
  held-out coverage.
- Ground truth and baseline method are frozen before implementation scoring.
- No fixture implies provider endorsement or customer proof.
- The generic-agent comparison is fair enough to falsify the Lumyn thesis.

#### Stop conditions

- Required source evidence cannot be licensed or pinned.
- The baseline gives Lumyn materially more context, permissions, or time than
  the comparator.
- Held-out answers are visible to the implementation worker.

**ADR impact:** none unless the holdout storage boundary changes.
**Changelog impact:** required.
**Cost/performance:** low to medium.

---

## Epic 2 — V3.1 Product And Trust Contracts

### M2: Define update, installation, agent, verification, delivery, and privacy contracts

**Priority:** P0
**Risk class:** High
**Blocked by:** M0
**Primary acceptance IDs:** `TRUST-001`–`TRUST-004`
**Supports:** `EVENT-001`–`EVENT-002`, `INSTALL-001`–`INSTALL-002`
**Implementation status:** Executable contract layer implemented; runtime
acceptance remains owned by M3 and later milestones.

#### Goal

Define the smallest executable contracts needed for the provider-originated,
consumer-installed update loop without building a universal registry or
hosted control plane.

#### Tasks

- Define schemas and validators for:
  - provider change event;
  - Provider Change Contract using internal artifact id `migration-pack`;
  - Consumer Installation and event-specific authorization snapshot;
  - integration graph;
  - impact report;
  - migration plan;
  - candidate manifest;
  - Agent Runner adapter manifest and conformance result;
  - agent attempt;
  - migration verification;
  - export result;
  - campaign summary;
  - provider status projection;
  - remediation outcome.
- Define separate status axes for impact, route, candidate, verification, and
  delivery.
- Define the consumer execution manifest for repository roots, readable and
  writable paths, commands, dependency posture, network, credentials, model
  mode, and `agent_execution_policy` (`disabled` or `configured`). Require the
  exact Agent Runner adapter/version/executable source and digest, auth mode and
  entitlement class, Agent Runner Vendor, actual Model Provider/model route,
  execution-funding mode, credential and usage-billing owners, native agent
  configuration, data egress, context policy, retention, and budgets only when
  agent execution is configured.
- Define independent scopes for repository read, local write, commands,
  Agent Runner network and credential, model-data egress, model network and
  credential, registry, lifecycle scripts, sandbox request data, sandbox
  network, sandbox credentials, remote branch, PR, retention, deletion, and
  API Provider reporting.
- Define the first provider channel as a signed, versioned JSON manifest at an
  exact provider-controlled HTTPS URL pinned by the installation. Bind a
  campaign key, monotonic sequence, issued/expiry times, embedded Provider
  Change Contract or exact provider-controlled contract URL, retrieved-byte
  digest verification, replay protection, audience matching, supersession,
  withdrawal, and an attended import recovery mode that cannot prove
  provider-channel delivery or authorize installed-preauthorization writes.
- Define provider-confirmation metadata without requiring elaborate root
  enrollment, continuous status refresh, or receipt exchange.
- Define installation action modes `notify_only`, `scan_only`,
  `prepare_patch`, and `open_draft_pr` as ceilings; authorization modes
  `per_event_approval` and `installed_preauthorization`; expiry and
  revocation; no stored GitHub token; runtime short-lived token issuance; and
  the rule that an event-specific snapshot may narrow but never widen
  installation authority. Default agent execution to disabled, allow
  notify-only, scan-only, and deterministic-only operation without a runner,
  and require an explicit configuration/authorization pause before any later
  agent-assisted route.
- Define `consumer_managed` as the default configured funding/credential mode and
  `provider_sponsored_lumyn_managed` as the optional campaign-funded mode.
  Bind exact credential owner and usage-billing owner, prohibit reusable
  credential storage, preserve consumer consent/revocation, and grant the API
  Provider no code, context, session, credential, or adapter-selection access.
- Define the managed-credential broker contract: approved issuer; exact
  installation, event, plan, attempt, runner, and model audience; maximum
  one-hour TTL; one-time redemption into one attempt-scoped session; multiple
  in-attempt calls only within hard token and cost quotas; no refresh,
  post-attempt replay, or cross-attempt reuse; revocation; post-run usage
  reconciliation; and no API Provider access. Require a vendor-native bounded
  credential or an approved budget-enforcing proxy; otherwise the managed mode
  is unavailable.
- Define the adapter-neutral contract for clean-session startup, exact adapter
  version, resolved executable source/digest, conformance digest, auth mode and
  entitlement class, structured lifecycle output, tool-call/edit/usage/error
  and cancellation provenance, native configuration policy, neutral home/config
  isolation, and fail-closed termination. Define Codex and Claude Code as
  launch targets only after common conformance plus live canary; defer Cursor
  behind the same gate. Prohibit repository-local executable shadowing, silent
  fallback, or dynamic mid-run switching.
- Define the runner host-isolation contract: explicit read-only and writable
  mounts, no host home or OS credential store, no ambient service sockets or
  unrelated inherited descriptors, child-process restriction inheritance,
  host-enforced egress, a bound backend/version/configuration and qualification
  digest plus host platform, hard CPU, memory, PID, process-tree-depth, disk,
  and open-file quotas, adversarial resource-exhaustion tests, fail-closed
  cleanup evidence, and no executable plugins, MCP servers, or hooks in the
  MVP.
- Define `lumyn check` as a non-mutating onboarding preflight. When agent
  execution is configured, it resolves the selected executable by canonical
  path, validates source/version/digest and current conformance, confirms auth
  mode and entitlement without collecting secrets, and confirms
  non-interactive actual-model-route identity without making a model call.
  When agent execution is disabled, runner readiness is not required.
- Define an allowlisted provider projection with explicit observed,
  consumer-reported, and unknown provenance. Silence is not `unaffected`;
  merge is not deployed or retired.
- Define local private-state storage outside the checkout and an explicit
  operator-managed retention/deletion policy for the services pilot.
- Define redaction before model transfer, persistence, export, or API Provider
  summary.
- Define agent provenance fields: Agent Runner Vendor, exact adapter version
  and executable source/digest, conformance digest, auth mode and entitlement
  class, clean-session identity, execution-funding mode, credential and
  usage-billing owners, actual Model Provider/model route, native
  configuration state/digest, policy digests, tools, commands, attempts,
  tokens, cost, changed files, and human input.
- Define proof labels and evidence invalidation.
- Define CLI grammar and typed error compatibility for all v3 commands.
- Add valid and invalid fixtures for every trust boundary.
- Document that prompt or instructions embedded in provider evidence,
  repository source, comments, tests, or generated output cannot override the
  approved plan.

#### Explicitly deferred

- Provider root-key enrollment.
- Universal event registry and signed invitation network.
- Provider status signer, rotation, revocation, and freshness service.
- Consumer connection receipts and provider acknowledgements.
- Connected-repository billing.
- Hosted campaign coordinator or dashboard.
- Product-owned universal host-isolation runtime.

Safety is not deferred: the consumer execution environment, least privilege,
redaction, no production access, and no auto-merge remain required.

#### Expected repo areas

- `schemas/`
- `internal/config/`
- future `internal/pack/`
- future `internal/change/`
- future `internal/installation/`
- future `internal/status/`
- future `internal/authorization/`
- future `internal/agent/`
- future `internal/evidence/`
- future `internal/redaction/`
- `examples/`
- `docs/product/`
- `docs/dev/`
- `docs/architecture/`
- `CHANGELOG.md`

#### Required tests

- Schema compilation and round trip.
- Every schema has valid and invalid fixtures.
- Cross-contract status and digest binding.
- Negative cases for implicit permission, model-data ambiguity, wildcard
  scope, prompt injection, secret persistence, provider data leakage,
  production credentials, default-branch write, and auto-merge.
- Adapter contract fixtures for Codex, Claude Code, and the deterministic fake;
  reject unqualified versions, unapproved executable source/digest,
  repository-local PATH shadowing, disallowed or ambiguous auth/entitlement,
  stale conformance, personal-session reuse, malformed or partial lifecycle
  output, undisclosed native configuration, opaque downstream model routing,
  owner-ambiguous billing, reusable credentials, and silent fallback. Cursor
  remains a negative/deferred fixture.
- Duplicate, stale, conflicting, superseded, withdrawn, wrong-audience, and
  unauthenticated provider-event denial; signed-manifest transport, sequence,
  expiry, pinned-origin, and attended-import non-proof tests.
- Installation expiry, revocation, event-policy widening, authorization-mode
  mismatch, stored-token rejection, disclosure overreach,
  silence-as-not-applicable-or-unaffected, and merge-as-retired denial.
- Compatibility tests for retained workflow and evidence contracts.

#### Completion criteria

- All v3 artifacts are executable contracts before runtime implementation.
- No product contract requires provider PKI, receipt billing, or hosted SaaS.
- The API Provider cannot gain repository or model-context access from
  sponsorship.
- Agent Runner/model-data posture, funding, credential/billing ownership,
  native configuration, and budgets are reviewable before execution.
- Provider evidence and repository content cannot act as control instructions.
- The same event and Provider Change Contract can be reused across authorized
  installations without sharing consumer-private artifacts.

#### Stop conditions

- A schema conflates API Provider, Agent Runner Vendor, or Model Provider.
- A generic scope implies repository mutation, model egress, credentials,
  network, remote delivery, or reporting.
- A provider event can widen installed authority or an installation becomes
  ambient cross-repository authority.
- The consumer-private root can resolve inside the checkout.
- Redaction, Agent Runner/model-data processing, credential ownership, or
  usage-billing ownership is ambiguous.

**ADR impact:** v3 execution and model-data boundary ADR required.
**Changelog impact:** required.
**Cost/performance:** low.

### M2.5: Pre-sell and qualify one provider campaign

**Priority:** P0 commercial gate
**Risk class:** High
**Blocked by:** none for `DISC-001` and `DISC-002`; `DISC-003` cannot close
until M2 contracts are approved
**Primary acceptance IDs:** `DISC-001`–`DISC-003`

#### Goal

Prove that an API Provider will pay for the outcome and can activate a real
cohort before Lumyn builds the expensive repository-specific execution stages.

#### Tasks

- Offer the `$7,500–$15,000` sunset-readiness sprint immediately and make it
  optionally creditable toward the campaign. Record it as paid discovery
  only: it cannot close `DISC-001`, prove consumer activation, or authorize
  repository work unless the separate `$25,000` campaign prepayment gate
  clears. Credited readiness funds count only after a signed campaign
  conversion makes them non-refundable campaign consideration and total
  cleared campaign funds reach at least `$25,000`.
- Sell a provider campaign and receive at least `$25,000` in cleared,
  non-refundable prepaid funds. Sell it as a bounded five-unit campaign
  attempt with a success floor of three valid scans and two tested outcomes,
  subject to consumer consent—not a promise of three to five completed
  migrations.
- Define the additional-repository billable event in the order form. Default:
  a repository beyond the included five-unit cohort reaches an independently
  verified Lumyn-generated candidate and consumer-authorized tested draft PR.
  A specifically contracted verified local-bundle fallback may be billable but
  cannot count as automated-delivery or MVP proof.
- Record the economic buyer, Provider Operator, hard 90-to-180-day deadline,
  source and target versions, business risk, and purchasing decision process.
- Obtain authoritative public or private migration evidence, provider
  commitment to confirm the Provider Change Contract once, and a named
  provider-owned distribution and onboarding motion for the cohort.
- Prequalify five reachable Eligible Consumer Units across five distinct API
  Consumer Organizations, each with one designated primary repository, an
  accountable maintainer, TypeScript/Node shape, plausible affected usage,
  and useful verification signals.
- Prequalify at least one real consumer repository for an agent-assisted run:
  record organizational/runtime feasibility for a consumer-selected Codex or
  Claude Code family (or explicit consent to a feasible managed route):
  available executable/version/source/digest, permitted non-interactive auth
  mode and entitlement, observable model-route identity, supported
  host-isolation posture, execution-funding mode, credential owner, and
  usage-billing owner without collecting secret values. Do not require or
  claim Lumyn adapter conformance before M6 implements the adapter; current
  conformance and live-canary evidence are mandatory immediately before the
  first M6 or M10 agent run. Scan-only and deterministic-only consumers may
  keep `agent_execution_policy=disabled`.
- For that repository, record a plausible naturally agent-eligible migration
  hypothesis grounded in its known shape, such as wrapper or adapter changes,
  signature or type adaptation, or related test repair. The hypothesis must
  require repository-specific reasoning but no new business decision. It is a
  prequalification signal only: it cannot force an agent route, and M4 impact
  plus M5 planning must confirm or reject it from actual authorized evidence.
- Confirm that at least one prequalified consumer is willing to install the
  `open_draft_pr` action with `installed_preauthorization` and runtime
  short-lived branch/draft-PR token issuance; record constraints without
  collecting repository credentials.
- Confirm that the same intended qualifying consumer reviews the exact
  allowlisted event-bound provider-status projection and is willing to
  authorize its transmission if the run reaches that state. This is
  prequalification evidence, not durable consent; execution still requires
  the installation's revocable `provider_reporting` grant.
- Before M2 privacy contracts are approved, `DISC-002` may be substantiated
  only from provider-controlled records reviewed under the provider's existing
  authority or by a privacy-approved non-identifying attestation. Lumyn does
  not copy repository or participant identities into its evidence system
  before the approved consent and storage boundary exists.
- After M2 approves the privacy, model-data, authorization, and evidence
  contracts, freeze before invitations:
  - cohort and eligibility;
  - campaign price and payment evidence;
  - source and target evidence;
  - model-data and privacy protocol;
  - agent-execution policy, qualifying Agent Runner/model route,
    auth/entitlement, execution-funding mode, credential owner, and
    usage-billing owner for the required agent-assisted run;
  - allowed private and API Provider-visible fields;
  - baseline method and material maintainer threshold;
  - correction and revert rubric;
  - Lumyn COGS boundary;
  - minimum contribution margin or maximum Lumyn Operator hours per
    reviewable outcome;
  - invitation, scan, outcome, and observation windows;
  - absolute campaign judgment deadline;
  - one primary provider outcome, source, denominator, comparator, and
    material pass threshold.
- Obtain consent for private evidence handling before collecting identifiable
  repository or participant data.
- Store private evidence outside the repository. Commit only consented,
  redacted aggregates or digests.
- Offer patch or PR-bundle fallback for assisted recovery, while making clear
  that at least one Lumyn-opened draft PR is required for pilot success. Do
  not require hosted dashboard, elaborate PKI enrollment, or receipt exchange.
- Record any consumer-paid sprint's price, real-repo evidence, effort, and
  outcome separately. Label a pre-M6 sprint as workflow/problem evidence, not
  engine proof; only a qualified post-M6 run may prove engine value. Neither
  can close `DISC-001` or `DISC-002`.
- A pre-M5 consumer sprint may use a customer-approved native Codex or Claude
  Code workflow as a manual service under the consumer's boundary. It is not a
  Lumyn-qualified adapter run, implemented Lumyn runtime, or release of the
  provider roadmap gate.
- Trigger product reframe review after two genuinely qualified provider
  opportunities fail to prepay or recruit rather than silently broadening the
  product.

#### Expected write boundaries

- approved private external evidence system;
- `.factory/artifacts/product-signals/M2.5/` for consented, redacted public
  evidence only;
- no `docs/product/` write; any finding that changes product truth must use a
  separate reviewed rebaseline and regenerate the active control set;
- validation scripts and schemas created only in later implementation tasks.

#### Evidence checks

- Cleared non-refundable funds and paid-invoice evidence.
- Readiness-sprint payment, if any, separately traced and reconciled. Any
  amount counted toward the campaign gate has signed conversion and allocation
  evidence making it non-refundable campaign consideration, and total cleared
  campaign funds still reach at least `$25,000`.
- Named buyer and operator.
- Five distinct Eligible Consumer Units across five distinct API Consumer
  Organizations.
- Provider-led distribution/onboarding commitment and one qualified
  short-lived draft-PR authorization path.
- One intended qualifying consumer's recorded willingness for the exact
  allowlisted event-bound provider-status projection to reach the API
  Provider, without treating willingness as runtime consent.
- One real consumer agent-assisted route with organizational/runtime
  feasibility recorded: runner family, available version/executable
  provenance, non-interactive auth/entitlement, observable model-route
  identity, host-isolation posture, funding mode, credential owner, and
  usage-billing owner, without secrets or a premature Lumyn-conformance claim.
- One plausible naturally agent-eligible migration hypothesis for that
  repository, grounded in wrapper, adapter, signature/type, or related-test
  evidence, without forcing the eventual M4/M5 route.
- Pre-outcome timestamps for the frozen protocol.
- Private/public field separation.
- Baseline and threshold fixed before repository execution.
- No consumer or provider evidence is fabricated from public fixtures.

#### Completion criteria

- `DISC-001`, `DISC-002`, and `DISC-003` have direct external evidence;
  `DISC-003` cites the approved M2 contract versions it instantiates.
- The economic buyer records that clearing the provider threshold would
  justify a retirement or paid-continuation decision and missing it makes the
  campaign fail.
- M5 is released only after the provider prepayment and cohort gates close.
- The first campaign can operate without a hosted control-plane build while
  still proving one automated draft-PR delivery.
- Failure is recorded honestly and triggers stop or reframe.

#### Stop conditions

- Payment is refundable, uncleared, contingent only, nominal, or an informal
  expression of interest.
- The provider cannot recruit five Eligible Consumer Units across five
  distinct API Consumer Organizations.
- The migration is outside the supported TypeScript/Node boundary.
- Cohort, threshold, or deadline would be chosen after outcomes.
- Private evidence lacks an approved storage and consent boundary.

**ADR impact:** none unless a new data or execution boundary is requested.
**Changelog impact:** none for evidence-only work.
**Cost/performance:** medium sales and operator effort.

---

## Epic 3 — Change Understanding And Repository Impact

### M3: Normalize provider intent into a reusable Provider Change Contract

**Priority:** P0
**Risk class:** High
**Blocked by:** M1, M2
**Primary acceptance IDs:** `PACK-001`–`PACK-004`,
`EVENT-001`–`EVENT-002`

#### Goal

Turn pinned official, public, and provider-confirmed evidence into one
provider-confirmed change contract and versioned provider event without
letting the parser or coding agent invent provider intent.

#### Tasks

- Extend source intake from one API surface to source and target OpenAPI,
  documentation, and SDK artifacts.
- Record source digests, versions, provenance, license posture, confirmation
  status, and concrete source locations.
- Assign stable change identity, intended audience and applicability, deadline,
  severity, and supersession or withdrawal state.
- Normalize supported method, request-property, response-property, type, and
  signature changes.
- Preserve provider-declared and Lumyn-detected unsupported, ambiguous, and
  needs-input conditions.
- Detect conflicts among docs, OpenAPI, SDK types, examples, and migration
  guidance.
- Prohibit executable hooks and arbitrary provider scripts.
- Render a reviewable Provider Change Contract for one-time Provider Operator
  confirmation and a non-executable event referencing its exact digest.
- Implement publish-kit output for the Provider Operator and consumer-side
  fetch/verify or attended-import intake for the pinned signed-manifest
  channel. The publish kit embeds the Provider Change Contract or pins its
  exact provider-controlled HTTPS URL; the consumer verifies retrieved bytes
  against the event digest. Only the provider-controlled HTTPS path counts as
  channel delivery; attended import is labeled recovery.
- Enforce issuer key, origin, sequence, issued/expiry time, audience, contract
  digest, duplicate/replay, supersession, and withdrawal checks before an
  event can reach installation policy.
- Ensure `public_derived` contracts remain usable for engineering and consumer
  sprints without implying provider endorsement.

#### Expected repo areas

- `internal/source/`
- future `internal/pack/`
- future `internal/change/`
- `schemas/`
- `examples/migration-packs/`
- `tests/`
- `docs/`

#### Required tests

- JSON and YAML normalization equivalence.
- Source and target digest mismatch.
- Conflicting docs/spec/SDK semantics.
- Missing or stale source evidence.
- Public-derived versus provider-confirmed status.
- Event replay, audience, supersession, withdrawal, and contract-digest
  binding.
- Pinned HTTPS origin, detached signature, monotonic sequence, expiry, and
  attended-import recovery labeling.
- Arbitrary script and prompt-injection denial.
- Supported, agent-assisted, blocked, and unsupported classification fixtures.
- Offline determinism for pinned sources.

#### Completion criteria

- Every normalized change cites source and target evidence.
- Public-derived and provider-confirmed states cannot be confused.
- Ambiguity blocks affected routes.
- No contract or event field executes code or widens consumer authority.
- `EVENT-001` and `EVENT-002` have direct runtime evidence, not schema-only
  closure.

#### Stop conditions

- Source and target semantics conflict without accountable clarification.
- Required evidence is not pinned or licensed.
- A provider asks to ship executable migration logic in the contract.

**ADR impact:** no if source, event, and contract boundaries remain separate.
**Changelog impact:** required.
**Cost/performance:** low.

### M4: Build the consumer-local TypeScript repository impact inventory

**Priority:** P0
**Risk class:** High
**Blocked by:** M3
**Primary acceptance IDs:** `IMPACT-001`–`IMPACT-005`

#### Goal

Find how the selected API or SDK is actually used in a repository without
modifying the repository or claiming coverage Lumyn cannot prove.

#### Tasks

- Detect the selected npm SDK, installed version, package root, manifest, and
  lockfile state.
- Parse TypeScript project configuration using an AST or compiler-backed
  integration; text matching alone is insufficient.
- Enforce the authorized read root across path normalization, symlinks,
  project references, config extension, and package traversal.
- Build the internal integration-graph nodes and edges for:
  - imports and aliases;
  - wrappers and adapters;
  - SDK calls;
  - request and response mappings;
  - relevant configuration;
  - tests, mocks, cassettes, and fixtures;
  - generated, vendored, excluded, dynamic, and uncertain regions.
- Record why each call site is affected, unaffected, uncertain, or
  unsupported.
- Never label the repository `unaffected` without explicit coverage and
  exclusions.
- Score visible and independently held-out fixtures.
- Freeze M1 thresholds before scoring. Initial target: at least 95% recall,
  at least 90% precision, and zero false `unaffected` results for supported
  held-out sites.
- Measure analysis time and graph size.
- Bind the `impact_read_only` product-action route to the exact union of
  `customer_repo_read`, `artifact_retention`, and `artifact_deletion`.
  Command, write, runner/model, network, and credential scopes are prohibited.

#### Expected repo areas

- future `internal/impact/`
- future `internal/typescript/`
- `schemas/`
- `examples/consumer-repos/`
- `tests/`
- `docs/`

#### Required tests

- Direct and aliased imports.
- Wrapper and adapter chains within the supported graph budget.
- Dynamic access and reflection uncertainty.
- Request/response mapping.
- Test and fixture linkage.
- Generated and vendored exclusion.
- Monorepo, multi-root, and multi-version rejection.
- Symlink and config-reference escape denial.
- Repeated-run graph determinism.
- Held-out precision, recall, and false-unaffected report.

#### Completion criteria

- Every impact result names analyzed scope, coverage, exclusions, and
  uncertainty.
- The inventory provides enough local context for M5 routing.
- Held-out thresholds pass without implementation-worker access to answers.
- Public-fixture performance is not represented as real-repository proof.

#### Stop conditions

- The parser cannot enforce read scope.
- Wrapper depth or dynamic use makes coverage non-diagnostic.
- The target threshold is changed after scoring.

**ADR impact:** required if a TypeScript/Node helper crosses the Go process
boundary.
**Changelog impact:** required.
**Cost/performance:** medium.

---

## Epic 4 — Reviewable Routing And Hybrid Execution

### M5: Produce a no-write migration plan and consumer approval boundary

**Priority:** P0
**Risk class:** High
**Blocked by:** M2, M2.5, M4
**Primary acceptance IDs:** `INSTALL-001`–`INSTALL-002`,
`PLAN-001`–`PLAN-003`

#### Goal

Implement the consumer-owned installation boundary and turn provider change
intent, installed policy, and repository evidence into a complete,
reviewable plan before any file or model-assisted mutation.

#### Tasks

- Implement `lumyn install` for the exact provider-channel key/origin,
  repository/package root, selectors, action ceiling, authorization mode,
  paths, commands, `agent_execution_policy`, model-data policy and budgets,
  GitHub token-issuance policy, reporting, retention/deletion, expiry, and
  revocation. When agent execution is configured, additionally bind the
  qualified exact Agent Runner adapter/version/executable source/digest, auth
  mode and entitlement class, execution-funding mode, credential and
  usage-billing owners, native agent configuration, and Agent Runner/model
  policy and budgets. Reject stored reusable Agent Runner, model, or GitHub
  credentials.
- Route each affected item to `deterministic`, `agent_assisted`, `manual`,
  `needs_input`, or `blocked`.
- List every proposed or conditional file and dependency change.
- Include the exact Provider Change Contract item and repository-impact
  evidence for
  every route.
- Define agent context selection; do not send the entire repository by default.
- List the explicit disabled or configured `agent_execution_policy`. For a
  configured policy, list Agent Runner Vendor and exact qualified
  adapter/version/executable source/digest/conformance digest, auth mode and
  entitlement class, clean-session and native-configuration posture,
  execution-funding mode, credential and usage-billing owners, actual Model
  Provider or local route, data-egress and retention posture, tools, commands,
  Agent Runner/model network and credentials, and all budgets.
- List baseline and post-change verification stages and expected proof level.
- List installed action mode, delivery mode, and exact API Provider-visible
  fields, if any.
- Derive an immutable event-specific authorization from the Consumer
  Installation and reject every attempted widening.
- Under `per_event_approval`, bind exact Consumer Maintainer approval to event,
  contract, installation and authorization digests, repository base, selected
  package root, paths, routes, `agent_execution_policy`, tools, commands,
  budgets, verification, delivery, and disclosure. When configured, also bind
  the exact Agent Runner/model policy and funding/credential/billing ownership.
- Under `installed_preauthorization`, evaluate the same values against the
  installed policy and pause on any mismatch; action-mode labels alone do not
  authorize side effects.
- Invalidate or pause authorization when any bound input changes.
- Prove planning performs no repository or Git mutation.

#### Expected repo areas

- future `internal/migrationplan/`
- future `internal/authorization/`
- future `internal/installation/`
- `schemas/`
- `tests/`
- `docs/`

#### Required tests

- Filesystem and Git before/after immutability.
- Stable plan output for identical inputs.
- Missing business value, ambiguous semantics, and unsupported route.
- Changed event, contract, installation, base, graph,
  `agent_execution_policy`, path, command, or budget invalidates approval.
  For configured execution, changed Agent Runner adapter,
  version/executable/conformance/auth/entitlement, funding/credential/billing
  owner, model policy, or native configuration also invalidates approval.
- Installation expiry, revocation, wrong audience, unauthorized action mode,
  authorization-mode mismatch, stored-token input, and event-policy widening
  fail closed.
- Installed-preauthorization tests prove an in-policy plan can proceed without
  per-event approval and any out-of-policy plan pauses before mutation.
- Plan cannot infer Agent Runner/model network or credentials, model egress,
  funding/billing ownership, remote delivery, or API Provider reporting.
- Provider instructions and source comments cannot alter the control plan.

#### Completion criteria

- The Consumer Maintainer can review the complete mutation, model, command,
  verification, delivery, and disclosure boundary.
- Every change is routed explicitly.
- No write or model call occurs during plan mode.
- Approval cannot widen implicitly.
- `INSTALL-001` and `INSTALL-002` close from executable installation and
  authorization evidence, not schemas alone.

#### Stop conditions

- M2.5 provider prepayment or cohort evidence is missing.
- A planned agent route lacks explicit semantics or verification.
- Agent Runner/model-data handling or credential/billing ownership is
  ambiguous.

**ADR impact:** no if M2 boundaries remain unchanged.
**Changelog impact:** required.
**Cost/performance:** low.

### M6: Implement deterministic transforms and the bounded coding agent

**Priority:** P0
**Risk class:** High
**Blocked by:** M5
**Primary acceptance IDs:** `AGENT-001`–`AGENT-007`

#### Goal

Produce a repository-specific migration candidate using deterministic code
where safe and a tightly bounded coding agent where local adaptation is
necessary.

#### Tasks

- Implement deterministic transforms for:
  - dependency and import update;
  - method or operation rename;
  - request-property rename or relocation;
  - response-property rename or relocation.
- Implement one adapter-neutral Agent Runner v1 seam plus pinned Codex and
  Claude Code adapters. Resolve each executable by canonical path, approved
  source, and digest. Advertise or use an adapter live only after that exact
  executable/version passes the common conformance suite and an approved live
  canary and its auth mode/entitlement class permits the intended automation.
  Keep Cursor deferred behind the identical gate. Do not build a proprietary
  model layer, arbitrary adapter marketplace, or dynamic mid-run switching.
- Start every attempt as a clean ephemeral session with neutral home/config
  roots and no repository-local PATH shadowing. Ignore personal history,
  memories, plugins, and user/project rules by default. If native
  configuration is explicitly selected, bind its digest, treat it as untrusted
  context, and prevent it from widening Lumyn policy.
- Require `agent_execution_policy=configured` before any agent attempt. Support
  default configured mode `consumer_managed` and optional
  `provider_sponsored_lumyn_managed` execution. Inject credentials only at the
  task-scoped runner/model boundary, record credential and usage-billing
  owners, persist no reusable credential, and give the API Provider no agent
  or repository access.
- Select exactly one agent authorization topology per attempt:
  `local_runtime` with no external runner/model egress; `runner_mediated` with
  runner network/credential and model-disclosure minimums; `direct_model` with
  model network/credential/disclosure minimums; or `hybrid` with both remote
  minimum sets. Enforce the selected minimum set before launch and keep
  package-registry read independently conditional.
- For `provider_sponsored_lumyn_managed`, require the approved broker to bind
  issuer, installation/event/plan/attempt and runner/model audience, maximum
  one-hour TTL, and one-time redemption into one attempt-scoped session.
  Permit multiple in-attempt calls only within hard token/cost quotas; forbid
  refresh, post-attempt replay, and cross-attempt reuse; require revocation and
  reconciliation. Use a budget-enforcing proxy when the vendor cannot issue
  those bounds; otherwise block the managed mode.
- Run the Agent Runner process with explicit mounts, no host home or OS
  credentials, no ambient service sockets or unrelated inherited descriptors,
  inherited child-process restrictions, host-enforced egress, a pinned and
  qualified isolation backend identity, hard CPU, memory, PID,
  process-tree-depth, disk, and open-file quotas, and cleanup evidence.
  Prohibit executable plugins, MCP servers, and hooks for the MVP.
- Normalize startup, structured output, tool calls, edits, usage, errors,
  cancellation, exit, and provenance across adapters. Do not silently fall
  back between adapter, version, model route, endpoint, credential owner, or
  billing owner.
- Implement `lumyn candidate import --manual` so an approved human-edited diff
  is checked against the exact base, pack, plan, route, paths, and diff budget
  before it can enter verification.
- Pass only the approved Provider Change Contract entries, repository-impact
  evidence,
  selected repository files, relevant tests, and repair diagnostics.
- Treat provider documents, repository source, comments, test output, and
  generated content as untrusted data that cannot alter system policy.
- Enforce:
  - readable and writable paths;
  - file and diff budgets;
  - tool and command allowlists;
  - separate Agent Runner and model network and credential policy;
  - context and data-egress policy;
  - time, token, cost, and attempt budgets;
  - default maximum of three attempts.
- Freeze the exact product-action capability union before each route:
  a pure parser/AST deterministic transform requires repository read/write,
  retention, and deletion only; an explicit package or lockfile-tool route
  additionally selects command execution and, only when needed,
  package-registry read. Agent-assisted candidates select command execution
  plus the exact Agent Runner/model disclosure, network, and credential
  scopes. No optional capability universe grants ambient authority.
- Run in an isolated worktree or consumer-approved equivalent.
- Require every edit to map to a pack item, graph evidence, route, and
  rationale.
- Update dependency and `package-lock.json` only through an approved,
  reproducible toolchain posture with lifecycle scripts disabled by default.
- Stop on missing business values, auth or event redesign, ambiguous
  semantics, production-only behavior, non-diagnostic failure, or scope
  expansion.
- Record Agent Runner Vendor, exact adapter version, executable source/digest
  and conformance digest, auth mode and entitlement class, clean-session
  identity, execution-funding mode, credential and usage-billing owners,
  actual Model Provider/model route, native-configuration state/digest, policy
  digests, tools, commands, attempts, tokens, cost, changed files, and human
  input.
- Do not persist raw prompts or responses by default.

#### Expected repo areas

- future `internal/patch/`
- future `internal/agent/`
- future `internal/workspace/`
- future `internal/authorization/`
- `schemas/`
- `examples/candidates/`
- `tests/`
- `docs/`

#### Required tests

- Byte-stable deterministic golden patches.
- Common Agent Runner contract tests with a deterministic fake.
- Pinned Codex and Claude Code conformance tests for clean-session startup,
  executable resolution/integrity, allowed auth/entitlement, normalized
  structured output, tool/edit/usage provenance, cancellation, error mapping,
  and fail-closed exit. Cursor remains a deferred negative fixture.
- Manual-candidate import tests for base, pack, plan, route, path, and diff
  binding, including stale and out-of-scope rejection.
- One explicitly approved live-agent canary per advertised adapter/version
  before pilot use.
- Prompt injection through docs, source, comments, tests, and tool output.
- Path, symlink, file, diff, command, network, credential, token, cost, time,
  and attempt budget enforcement.
- Host-home, OS-credential, socket, inherited-descriptor, mount, egress, and
  cleanup denial, including malicious child-process and tool fixtures.
- Executable plugin, MCP server, and hook rejection.
- Managed-credential audience, TTL, one-time redemption, in-attempt quota,
  no-refresh, no-replay/cross-attempt-reuse, revocation, reconciliation, and
  unsupported-vendor/proxy failure cases.
- Unrelated-edit and scope-expansion rejection.
- Agent unavailable, unqualified version, authentication failure, model
  unavailable, BYOK failure, invalid entitlement, untrusted or shadowed
  executable, malformed/partial structured output, stale conformance, and
  cancellation.
- Personal-session reuse, undisclosed native configuration, native-rule scope
  widening, opaque model routing, credential persistence, owner-ambiguous
  billing, and silent adapter/model fallback rejection.
- Both `consumer_managed` and `provider_sponsored_lumyn_managed` contract tests;
  live proof may select one approved mode and need not exercise both.
- Raw prompt, response, secret, and credential non-persistence.
- Missing business value and unsupported semantic stop.
- Candidate staleness on changed pack, base, plan, policy, or graph.
- Mutation tests over budget and route checks.

#### Completion criteria

- Deterministic cases are repeatable.
- Both launch adapters satisfy one contract, and agent-assisted cases remain
  inside every approved boundary.
- Every edit has provenance and a migration rationale.
- Approved manual candidates have explicit provenance and enter the same
  independent verification ladder without being relabeled deterministic or
  agent-assisted.
- The agent cannot self-approve, self-verify, or broaden the task.
- Cost, attempt, credential-owner, and usage-billing evidence is available for
  campaign COGS and consumer-managed spend.

#### Stop conditions

- The execution environment cannot enforce the approved boundaries.
- Required code context would violate the consumer's Agent Runner/model-data
  policy.
- The selected adapter/version lacks current conformance, exact model-route
  disclosure, or enforceable clean-session and no-fallback behavior.
- The agent needs production credentials or a new business decision.
- Candidate generation cannot satisfy the installed scope, provenance, cost,
  or independent-verification contract. Patch similarity to a generic agent
  is not by itself a failure; end-to-end substitution is judged in M10.

**ADR impact:** Agent Runner, model-data, and workspace isolation ADR required.
**Changelog impact:** required.
**Cost/performance:** medium to high; measured per attempt.

---

## Epic 5 — Baseline-Aware Verification And Repair

### M7: Verify deterministic, agent-assisted, and manual candidates

**Priority:** P0
**Risk class:** High
**Blocked by:** M6
**Primary acceptance IDs:** `VER-001`–`VER-006`

#### Goal

Produce proof-honest repository and workflow evidence, and allow only bounded
diagnostic repair.

#### Tasks

- Capture pre-candidate dependency, compile, typecheck, and selected-test
  baseline.
- Run deterministic, agent-assisted, and imported manual candidates through
  the same verification ladder and proof-label rules.
- Keep `lumyn verify` non-mutating with respect to the candidate. Run untrusted
  repository commands in a disposable verification view or detect and reject
  any candidate-tree or index mutation.
- Start verification in a fresh process and verification view with frozen
  command and verification-configuration digests. Keep Agent Runner/model
  credentials and generation-owned evidence handles absent; only the
  verifier/evidence boundary may write verification evidence.
- Run candidate integrity, dependency integrity, compile, typecheck, and
  consumer-allowlisted tests.
- Run repository commands in the consumer-approved environment with explicit
  command and working-directory allowlists, read-only/writable mounts, neutral
  home and temporary roots, executable roots, sanitized environment, inherited
  child-process restrictions, timeout/output budgets, and socket, descriptor,
  OS-credential, and ambient-secret denial. Bind the exact qualified
  backend, version, configuration, host platform, hard CPU, memory, PID,
  process-tree-depth, disk, and open-file quotas, and cleanup evidence; fail
  closed when any limit is unenforceable.
- Keep tests offline and lifecycle scripts disabled by default. Any network or
  lifecycle-script exception requires its own exact route grant and approval.
- Reuse retained workflow, cassette, trace, proof, boundary, and redaction
  schemas where their semantics remain valid.
- Keep independent contract replay separate from exact-head replay, mock, and
  sandbox evidence.
- Require exact candidate-head causal execution for every
  `workflow_verified_*` label.
- Implement `lumyn repair` as a separate command. Do not mutate from
  `lumyn verify`; feed actionable diagnostics back to M6 only after a separate
  Consumer Maintainer repair authorization. `lumyn repair` is agent-assisted
  only and requires `agent_execution_policy=configured`. For a failed
  agent-assisted candidate, bind and reuse its exact Agent Runner
  adapter/version/executable, actual model route, execution-funding mode,
  credential owner, and usage-billing owner unless a new route is explicitly
  authorized. For a failed deterministic or imported-manual candidate, require
  a newly configured, explicitly authorized exact agent route or return
  `needs_input`/`blocked`. Bind the failed candidate and evidence, exact repair
  intent, remaining permissions, and remaining attempt, time, token, cost,
  file, and diff budgets. Every repair is a new attempt and candidate.
- Every repair creates a new candidate head, invalidates prior verification
  evidence, and requires a fresh full verification run.
- Stop on non-diagnostic failures, missing business input, boundary failure,
  stale evidence, redaction uncertainty, or exhausted budget.
- Bind evidence to pack, graph, plan, base, candidate, route,
  `agent_execution_policy`, commands, environment, and artifact hashes. When
  configured, also bind Agent Runner, exact adapter/version and conformance,
  execution-funding mode, credential and usage-billing owners, and actual model
  policy.
- Bind `verify` to the exact product-action union of `customer_repo_read`,
  `command_execution`, `artifact_retention`, and `artifact_deletion`.
  `repair_agent_assisted` adds `customer_repo_write` plus the exact
  route-selected Agent Runner/model disclosure, network, and credential
  scopes. Freeze that union in the private authorization; infer no scope from
  the task's optional capability universe.
- Implement `lumyn trace` as a local, no-network evidence renderer.
- Implement an offline canary:

  ```text
  pack -> impact -> plan -> apply -> verify -> export preview
  ```

  The canary uses public or synthetic fixtures, a fake or explicitly approved
  agent, no live credentials, no API Provider reporting, and no remote write.

#### Expected repo areas

- `cmd/lumyn/`
- future `internal/verify/`
- future `internal/replay/`
- future `internal/evidence/`
- future `internal/redaction/`
- future `internal/authorization/`
- `schemas/`
- `workflows/`
- `cassettes/`
- `runs/`
- `tests/`
- `docs/`

#### Required tests

- Pre-existing failure attribution.
- Same verification ladder for deterministic, agent-assisted, and imported
  manual candidates.
- Candidate-tree and Git-index immutability across verification, including a
  mutating-test negative fixture.
- Fresh verification view/process, absent runner/model credentials,
  generation-unwritable evidence, and frozen command/config digest tests.
- Missing, stale, or scope-widening repair authorization; new-candidate and
  evidence-invalidation behavior after an approved repair.
- Network, secret, and lifecycle-script denial.
- Flaky, timeout, non-diagnostic, and exhausted-repair cases.
- Exact-head versus wrong-head causal binding.
- Static, repo, contract replay, replay, mock, and sandbox label separation.
- Boundary, redaction, stale evidence, and false-green cases.
- Negative suite with zero false verified outcomes.
- Evidence digest invalidation.
- Offline canary golden and stage-failure paths.
- Local trace rendering with no network or provider report.

#### Completion criteria

- Pre-existing and migration-attributable failures remain separate.
- Verification never mutates the candidate, and repair never occurs implicitly
  from a verification command.
- A failed or partial result cannot become verified.
- The repair loop never exceeds the approved plan or attempts.
- Every held-out negative case avoids false verification.
- The offline canary reaches a local verified export preview without hidden
  authority.

#### Stop conditions

- Repository commands cannot run inside the approved environment.
- A weak evidence level would be promoted to a stronger label.
- The agent requests scope expansion to repair a failure.
- Verification is too weak to support a customer-specific claim.

**ADR impact:** verification, repair, and execution environment ADR required.
**Changelog impact:** required.
**Cost/performance:** medium; command and repair duration measured.

### M8: Add optional provider sandbox read-back

**Priority:** P1 conditional
**Risk class:** High
**Blocked by:** M7 and explicit consumer approval
**Supports:** `VER-003`, `VER-004`, `VER-005`

#### Goal

Add stronger non-production outcome evidence when a real provider sandbox and
consumer permission exist. M8 is not required for M9 or M10 when repository or
mock evidence is sufficient.

#### Tasks

- Add exact allowlisted sandbox endpoints and operations.
- Require separate command, payload disclosure, network, credential,
  retention, and cleanup approval.
- Use only synthetic or approved non-sensitive test data.
- Run build and test commands under the same fail-closed repository-command
  isolation contract as M7: exact commands/directories/mounts/executable roots,
  neutral home/temp, sanitized environment, inherited child restrictions,
  timeout/output budgets, no ambient secrets, sockets, descriptors, host home,
  or OS credentials, and network/lifecycle scripts disabled by default.
  Isolate sandbox credentials from those commands.
- Execute the approved entrypoint from the exact candidate head.
- Run that credential-bearing entrypoint in a separate qualified isolation
  profile with a read-only exact-head mount, exact entrypoint and working
  directory, neutral roots and sanitized environment, only the task-scoped
  sandbox credential, endpoint/operation-only egress, inherited child limits,
  hard CPU/memory/PID/process-tree/disk/open-file quotas, teardown, cleanup,
  and orphan evidence.
- Enforce namespace, request/write budget, idempotency, retries, read-back,
  cleanup, and orphan reporting.
- Preserve sandbox-versus-production limitations.
- Record provider logging, retention, and deletion terms.
- Bind the single `sandbox_read_back` action to the exact fixed union of
  `customer_repo_read`, `command_execution`, sandbox request disclosure,
  sandbox network, sandbox credential, artifact retention, and artifact
  deletion. The task-level authority array is a capability universe, not an
  ambient grant.

#### Expected repo areas

- future `internal/live/`
- future `internal/verify/`
- future `internal/authorization/`
- future `internal/redaction/`
- `schemas/`
- `tests/`
- `docs/`

#### Required tests

- Mock timeout, retry, budget, auth, namespace, cleanup, orphan, and redaction
  cases.
- Missing or expired approval.
- Credential non-leakage across stages.
- Production data, PII, secret, and unapproved payload denial.
- Wrong-head and stale-candidate denial.
- One live integration only after explicit approval.

#### Completion criteria

- Sandbox use cannot begin without every required scope.
- Cleanup succeeds or produces explicit orphan evidence.
- Sandbox verification is never presented as production proof.
- Absence of sandbox evidence does not weaken an otherwise honest repository
  outcome label.

#### Stop conditions

- Production credentials or customer data are required.
- The provider cannot state logging and retention behavior.
- Cleanup cannot be made safe.

**ADR impact:** required before live use.
**Changelog impact:** required when implemented.
**Cost/performance:** medium.

---

## Epic 6 — Consumer-Controlled Delivery

### M9: Export evidence and open a tested draft PR

**Priority:** P0
**Risk class:** High
**Blocked by:** M7; M8 is optional
**Primary acceptance IDs:** `EXP-001`–`EXP-004`

#### Goal

Own the composed installed-event-to-draft-PR path and deliver its verified
outcome through the consumer's normal review workflow while retaining local
fallback and never taking merge authority.

#### Tasks

- Render a complete local evidence bundle and reviewer checklist.
- Implement `lumyn update --event` as the composed orchestrator from verified
  provider-channel intake and event-specific installation authorization
  through impact, plan, Lumyn deterministic or agent-assisted candidate,
  independent verification, non-default branch, and tested draft PR. It stops
  at the installed action ceiling and never relabels a manual candidate as
  product-loop proof. It then records a local event-bound provider-status
  projection or explicit reporting decline; provider transmission remains a
  separate optional grant and cannot block draft-PR delivery.
- Capture and revalidate on that run's immutable authorization snapshot the
  exact union for each product action it exercises. M9 retains only delivery,
  reporting, artifact, and local repository/command scopes; it references the
  already validated M4, M6, and M7 action authorizations rather than copying or
  aggregating their Agent Runner, model, registry, or verification authority.
  Provider reporting remains separately optional. Completed M6 or M8 work and
  Factory approval, credential, or network grants do not delegate consumer
  permissions.
- Treat M9's task-level authority arrays as a capability universe only. Freeze
  a separate exact union before each action:
  - `local_export` uses repository read, command, retention, and deletion only;
  - `local_branch` additionally uses repository write;
  - `remote_branch_push` uses repository read, command, retention, deletion,
    and only the remote-branch grant;
  - `draft_pr_create` uses retention, deletion, and only the draft-PR grant;
  - `provider_status_decline` records a local decline without
    `provider_reporting`;
  - `provider_status_transmit` uses the separate consented
    `provider_reporting` grant without repository, GitHub, runner, or model
    access;
  - the composed deterministic and agent-assisted sequences reference exact
    M4/M6/M7 actions, then `remote_branch_push`, then `draft_pr_create`; no
    aggregate composed route or cross-action authority union exists.
- Export:
  - patch;
  - local branch;
  - PR-ready bundle with suggested title and body.
- Label manual push and manual PR creation honestly.
- Preserve evidence bindings across export and mark stale exports.
- Produce an event-bound, consumer-consented provider status projection
  without raw private artifacts and with provenance for observed,
  consumer-reported, and unknown state.
- Implement `lumyn outcome record` as a local, append-only command for
  authorized consumer acceptance, merge, closure, correction, and reversion
  evidence bound to the exact candidate head and verification-evidence digest.
- Measure which delivery steps actually create pilot friction.
- Implement `lumyn pr create --draft` after the GitHub App installation and
  short-lived installation-token issuance path is designed and approved; the
  App installation may persist, but the token may not.
- Keep remote branch and draft-PR permissions separate.
- Use only a non-default branch and draft posture. Bind idempotency to the
  event, Provider Change Contract, Consumer Installation authorization,
  repository, base and candidate heads, plan, and verification evidence.
- Never auto-merge.
- Include current M8 sandbox evidence when available without making it a
  prerequisite.

#### Expected repo areas

- future `internal/export/`
- future `internal/report/`
- future `internal/outcome/`
- future `internal/status/`
- future `internal/github/`
- `schemas/`
- `tests/`
- `docs/`

#### Required tests

- Patch, branch, and PR-bundle golden output.
- Export staleness and idempotency.
- Manual delivery labeling.
- Provider summary field allowlist and redaction.
- Status-event/evidence binding, explicit provenance, silence-as-unknown, and
  merge-not-retired tests.
- No-provider-reporting export.
- Outcome authority, exact candidate/evidence binding, append-only history, and
  correction/reversion tests; plan approval, PR creation, and informal
  acknowledgement cannot produce `consumer_accepted`.
- GitHub tests for short-lived token, permission denial, default-branch
  denial, duplicate PR, stale base, and draft-only behavior.
- Prove a manual bundle does not emit automated-delivery success.
- Prove `EXP-003` cannot close `not_applicable` or from manual delivery
  evidence.
- Composed-path tests prove `EXP-003` cannot close from standalone
  `lumyn pr create --draft`, an imported manual candidate, attended event
  import, or disconnected component evidence.
- Status tests bind any projection to the same event, installation
  authorization, candidate, verification evidence, and draft PR; missing
  reporting consent records a decline without blocking delivery.

#### Completion criteria

- A complete local fallback remains available if GitHub delivery fails.
- Provider reporting remains optional and consented.
- Durable outcome recording remains consumer-controlled and does not disclose
  private evidence to the API Provider.
- `EXP-003` closes only with automated draft-PR evidence.
- Automated PR delivery uses short-lived authorization and cannot write the
  default branch or merge.
- Repeated export or delivery cannot create conflicting state.

#### Stop conditions

- Local export is blocked on hosted infrastructure.
- A long-lived broad GitHub credential is required.
- Provider reporting would disclose raw consumer evidence.
- Manual delivery is proposed as proof of automated delivery.

**ADR impact:** export and GitHub App installation-token ADRs required.
**Changelog impact:** required.
**Cost/performance:** low for local export, medium for GitHub integration.

---

## Epic 7 — Provider Campaign And Outcome Learning

### M10: Run one prepaid provider-originated update campaign

**Priority:** P0 product validation
**Risk class:** High
**Blocked by:** M2.5, M7, M9, consumer consent; M8 optional
**Primary acceptance IDs:** `PILOT-001`–`PILOT-008`

#### Goal

Determine whether one provider can publish a confirmed change once and Lumyn
can carry it through installed consumer policy to tested draft PRs and
measurable, consented rollout progress better than the status quo.

#### Tasks

- Advance the provider, migration, cohort, baseline, and measurement protocol
  frozen in M2.5.
- Confirm the Provider Change Contract once and have the Provider Operator
  publish one signed versioned manifest on the pinned provider-controlled
  channel. Any semantic, audience, deadline, supersession, or withdrawal
  change creates a new version and invalidates affected plans.
- Invite five prequalified Eligible Consumer Units across five distinct API
  Consumer Organizations.
- Install the exact provider channel, action, agent-execution policy,
  model-data, GitHub, and disclosure policy in each consenting repository.
  Bind a qualified Agent Runner route and execution-funding/credential/billing
  ownership only for agent-enabled installations.
- Compose, do not union, each installation's run from the validated M4 impact,
  M6 deterministic or agent-assisted candidate, M7 verify or separately
  authorized repair, optional M8 sandbox, and M9 export, draft-PR, reporting,
  or reporting-decline routes. Freeze the ordered route sequence and exact
  capability union for every action before it runs. Campaign-level authority
  arrays are aggregate maxima, never a grant to every installation.
- Keep campaign proof minima separate from per-consumer authority: the
  same qualifying composed run must include one organically agent-assisted
  plan item, one remote draft-PR route, and one reporting-transmit route, but
  no installation receives those scopes unless its own policy selects and
  authorizes them.
- Run impact locally in each consenting repository.
- Require at least three valid impact scans within 14 calendar days of
  invitation.
- Produce at least two tested, reviewable migration outcomes by the frozen
  deadline.
- Produce one qualifying real-repository outcome that starts at the
  authenticated provider channel, uses `installed_preauthorization`, contains
  at least one plan item organically routed `agent_assisted` because
  repository-specific reasoning is necessary, executes it through a
  consumer-selected qualified Codex or Claude Code route without bespoke Lumyn
  Operator code edits, passes independent exact-head verification, and reaches
  a Lumyn-opened draft PR plus provider-received consented status projection
  through the composed `lumyn update --event` flow. Do not reroute
  deterministic work merely to manufacture proof. Imported manual candidates,
  standalone PR creation, and a separate agent run plus deterministic composed
  run cannot satisfy `PILOT-003`.
- Persist that proof as one qualifying-run record with one `run_id` and exact
  references for the provider event, installation authorization, organic agent
  plan item, candidate head, verification-evidence digest, remote branch,
  draft PR, and provider-status projection. Every artifact binds the same run,
  event, and installation authorization; cross-run substitution or a missing
  reference leaves `PILOT-003` open.
- Obtain at least one `consumer_accepted` outcome or merge.
- Use `lumyn outcome record` to record `consumer_accepted` only with a durable
  artifact naming the Consumer Maintainer and authority, API Consumer
  Organization, repository, exact candidate head, verification-evidence
  digest, adoption decision, and timestamp; plan approval, PR creation, or
  informal acknowledgement does not count.
- Record manual, deterministic, and agent-assisted contribution separately.
  Bespoke Lumyn Operator edits are manual and cannot count as automation.
- Run the frozen guide or codemod plus generic-agent baseline under the same
  repository snapshot, authoritative migration evidence, selected Agent Runner
  adapter/version/executable, actual Model Provider/model version,
  auth/entitlement and execution-funding route, credential and usage-billing
  owners, context-access ceiling, tools, commands, role, verification
  opportunity, and attempt/token/time/cost budgets. Lumyn's orchestration,
  impact/routing, boundary enforcement, independent verification, delivery,
  and status loop are the treatment. Unmatched engine comparisons remain
  descriptive.
- Pass the frozen material maintainer comparison threshold:
  - default: at least 30% lower median Consumer Maintainer hands-on time;
  - alternative: another equally material status-quo comparison frozen under
    `DISC-003` before execution;
  - guardrail: no worse substantive-correction, revert, or false-verification
    rate.
- Clear the frozen material provider-outcome threshold.
- Clear the frozen contribution-margin or maximum-operator-hours threshold so
  services labor cannot hide failed product automation.
- Record every accepted, merged, closed, rejected, blocked, reverted, and
  corrected outcome.
- Measure:
  - actual cleared provider payment;
  - Agent Runner, Model Provider, tool, and infrastructure spend by
    execution-funding mode and usage-billing owner, separating Lumyn campaign
    COGS from consumer-paid usage;
  - Agent Runner installation, qualification, authentication, fallback-block,
    and time-to-first-run friction by adapter/version;
  - Lumyn Operator hours;
  - Consumer Maintainer hands-on time;
  - provider setup and support time;
  - cost per verified and accepted migration;
  - funnel conversion and time;
  - one frozen primary provider outcome.
- Keep raw consumer evidence private. Provider reporting is limited to
  event-bound consented status or aggregates. Silence remains unknown and no
  status beyond observed evidence is inferred.
- Require each participating installation to emit or explicitly decline the
  frozen allowlisted event-bound status projection, and require at least one
  consented projection to reach the provider. The required projection must
  come from the same qualifying composed run used for draft-PR product proof.
- Judge the campaign at the frozen absolute deadline. Record only `pass` or
  `fail`; abandonment and timeout are `fail`, and `reframe` is only a
  post-failure disposition.
- If the campaign fails because the provider cannot recruit or generic-agent
  performance is equivalent, stop provider-specific automation instead of
  broadening scope.
- Record any consumer-paid sprint separately; it cannot satisfy provider
  pilot items.

#### Recommended primary provider outcomes

Choose exactly one before the first invitation:

- reduction in provider support hours per accepted migration;
- reduction in invitation-to-accepted-migration lead time;
- share of the frozen cohort off the targeted legacy version by the deadline.

The metric, source, denominator, comparator, and material threshold are frozen
before execution.

#### Expected evidence locations

- approved consumer-private evidence roots;
- approved provider-commercial evidence system;
- redacted aggregate or digest-only pilot evidence in the future compiled
  pilot artifact area;
- no committed source, credentials, prompts, responses, raw logs, or traces.

#### Evidence checks

- Cleared non-refundable funds, paid invoice, and the recorded accounting
  treatment for the prepaid engagement.
- Five Eligible Consumer Unit invitations across five distinct API Consumer
  Organizations.
- Three valid scans within 14 days.
- Two tested, reviewable outcomes.
- One qualifying run that carries one authenticated provider-channel event
  through installed preauthorization, an organically agent-assisted plan item
  on the consumer-selected qualified runner, independent exact-head
  verification, a tested Lumyn-opened draft PR, and a provider-received
  consented status projection. Its evidence binds the same event,
  installation authorization, plan item, candidate, verification digest,
  non-default branch, draft PR, and projection, including token, idempotency,
  draft-only, runner-route provenance, and no-bespoke-edit evidence.
- One valid `consumer_accepted` artifact or merged outcome.
- Fair generic-agent baseline.
- Substantive correction and revert review.
- Agent Runner/model/tool spend by funding mode and usage-billing owner, with
  Lumyn campaign COGS and consumer-paid usage separated.
- Frozen contribution-margin or operator-automation threshold result.
- Consumer Maintainer time.
- Primary provider outcome and material-threshold result.
- Provider-visible field consent.
- Explicit reporting declines for every participating installation that does
  not authorize the qualifying provider-status projection.
- Frozen-deadline `pass` or `fail`, plus any post-failure `reframe`
  disposition.

#### Completion criteria

- Every `PILOT` item has direct source evidence.
- Technical output cannot substitute for payment, consumer activation,
  acceptance, or material baseline advantage.
- Campaign success requires both the maintainer baseline and the material
  provider-outcome threshold.
- No unverified candidate counts.
- At least one installed-preauthorization composed run with an organically
  agent-assisted plan item, independent verification, Lumyn-opened draft PR,
  and provider-received same-run projection proves the application-layer
  handoff. Manual-only, deterministic-only, separate agent/delivery, or
  standalone-delivery evidence fails `PILOT-003`.
- Operator assistance and manual edits remain visible.
- The campaign is closed by its deadline.

#### Stop conditions

- The cohort or threshold would be changed after results.
- The provider requests access to raw consumer code.
- Consumer model-data or repository consent is missing.
- Production access is required.
- The guide/codemod-plus-generic-agent baseline is materially equivalent end
  to end on maintainer effort, verified completion, correction risk, provider
  rollout evidence, and provider outcome.
- The primary provider outcome misses its frozen material threshold.
- Fewer than three scans or two outcomes remain mathematically possible by the
  frozen deadline.

**ADR impact:** none unless the pilot requests a new runtime boundary.
**Changelog impact:** only if product behavior changes.
**Cost/performance:** high and measured as commercial evidence.

---

## Minimum-Now Sequence

### Planning rebaseline

- This change aligns the PRD, plan, operating docs, ADRs, compiled 53-item
  control set, validators, and paused factoryd templates.
- No runtime implementation task is authorized by this planning change.
- A v3.1 source update has landed in the external Factory profile, but the
  active compiled generation still records Factory profile/runtime
  compatibility as unqualified. factoryd execution remains paused until the
  complete controls are regenerated and its bundle/runtime and exact mission
  are qualified and explicitly unpaused. An attended implementation path
  outside factoryd still requires explicit task approval and every repo-local
  lifecycle gate; it does not prove factoryd readiness.

### Wave 1

- Sales may begin the readiness-sprint and campaign-prepayment work
  immediately, but M2.5 evidence writes remain an explicitly approved attended
  task; `DISC-003` cannot close until M2.
- After an attended M0 task is explicitly approved—or after factoryd
  qualification when factoryd is selected—M0 corrects false-green runtime
  behavior.
- M2 defines the executable event, installation, and authorization contracts
  after M0.
- M1 corpus work may be prepared in parallel, but its task cannot dispatch or
  complete the walking skeleton until M2 closes; the compiled dependency is
  therefore M0 plus M2.

### Wave 2

- M3 Provider Change Contract and event intake after M1 and M2.
- M4 TypeScript repository impact inventory after M3.

### Wave 3

- M5 no-write plan only after M2.5 commercial gates and M4.
- M6 deterministic plus bounded-agent execution after M5.

### Wave 4

- M7 baseline-aware verification, agent-assisted repair, and offline canary.

### Wave 5

- M9 local fallback export and automated draft-PR delivery.
- M8 sandbox verification may proceed independently when a real sandbox and
  consent exist.

### Wave 6

- M10 prepaid provider campaign after M2.5, M7, and complete M9 delivery.
- M10 does not wait for M8.

The dependency graph is:

```text
M0    -> none
M1    -> M0, M2
M2    -> M0
M2.5 start -> none
M2.5 DISC-003 closure -> M2
M3    -> M1, M2
M4    -> M3
M5    -> M2, M2.5, M4
M6    -> M5
M7    -> M6
M8    -> M7
M9    -> M7
M10   -> M2.5, M7, M9
```

This sequence starts sales qualification before runtime implementation, while
preventing privacy and model-data protocol freeze until the governing M2
contracts exist. Expensive agent execution still waits for cleared payment and
the qualified five-organization cohort.

---

## Explicit Non-Goals

- Generic buy-side monitoring of every API dependency or arbitrary public
  changelog scraping.
- Anonymous scanning of downstream repositories.
- A generic coding-agent product, proprietary foundation model, or separate
  provider-specific coding-agent codebase.
- An arbitrary Agent Runner marketplace, unqualified adapter, opaque model
  route, or dynamic mid-run adapter/model switching. Controlled
  installation-time selection between qualified Codex and Claude Code
  adapters is in scope; Cursor remains deferred.
- Authentication, webhook/event, GraphQL, gRPC, generated-client,
  cross-language, or production migrations.
- Broad package-manager and monorepo support.
- Provider-supplied executable migration scripts.
- Production credentials or mutation.
- Default-branch writes or automatic merge.
- Provider access to raw consumer source, model context, prompts, responses,
  logs, traces, or credentials.
- Hosted dashboard or coordinator as a prerequisite.
- Long-lived GitHub credential as a prerequisite.
- Elaborate provider PKI, a universal event network, key rotation service, or
  connection-receipt billing.
- Annual connected-repository pricing before repeated paid evidence.
- Calling manual bundle delivery automated PR delivery.
- Closing the first provider campaign without at least one Lumyn-opened,
  tested draft PR.
- Closing Agent Runner product proof without a real consumer-selected,
  independently verified agent-assisted run.
- Calling public fixtures customer proof.
- Calling consumer-paid work provider-demand proof.
- Calling the current repository or pilot distribution OSS before the
  explicit release gate closes.

---

## Definition Of Done

This plan is complete only when:

- all 53 PRD acceptance items are represented in the active compiled acceptance
  ledger and mapped to their primary closure owner;
- conditional items and optional M8 behavior are marked conditional rather
  than silently required; automated draft-PR delivery remains required;
- every technical item has schema, fixture, command, test, or
  proof-of-behavior evidence;
- all agent acceptance items have explicit model-data, scope, budget,
  provenance, repair, and prompt-injection coverage;
- held-out answers remain unavailable to implementation workers;
- all consumer privacy, no-production, provider non-disclosure, and
  no-auto-merge constraints pass;
- required CI, coverage, CodeQL, review, shipping, PR lifecycle, and
  post-merge evidence exists;
- historical artifacts are not reinterpreted;
- the provider campaign satisfies or explicitly fails its frozen commercial
  and outcome gates;
- payment, COGS, operator time, consumer effort, corrections, and baseline
  comparison are measured;
- README, PRD, plan, `AGENTS.md`, architecture and developer guidance, Factory
  profile, compiled task packets, validation contract, acceptance mapping, and
  scope closure agree before implementation dispatch.

No runtime implementation is performed or authorized by this planning rewrite.
