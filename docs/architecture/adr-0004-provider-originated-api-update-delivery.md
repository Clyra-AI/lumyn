# ADR-0004: Provider-Originated API Update Delivery

## Status

Accepted as the v3.1 product-direction decision. This decision changes planning
and acceptance semantics only. It authorizes no product runtime, customer
repository access, model call, credential, network action, branch write, or PR
write.

ADR-0003 remains authoritative only for bounded-agent execution, deterministic
verification, consumer privacy, and human merge control. This ADR supersedes
ADR-0003's services-led product identity, manual-first proof, and
optional-draft-PR product-proof clauses.
ADR-0005 refines the replaceable Agent Runner seam, customer selection,
adapter qualification, execution funding, credential ownership, and
independent-verification boundary.

## Date

2026-07-24

## Context

The YC Fall 2026 request for startups identifies a missing application layer:
API Providers should be able to carry an important API change into consenting
customers' codebases and open the fixing PR, instead of relying on changelogs
that customers must discover and interpret.

Lumyn v3.0 established useful safety and execution primitives, but framed the
product as a services-led migration utility and made automated draft-PR
delivery optional. That framing underweighted the distinctive cross-company
value and left the core YC loop unproven. A capable generic coding agent can
already attempt a migration when a maintainer asks it to. Lumyn must add value
before and after that edit:

- provider-originated distribution and reusable authoritative intent;
- durable, consumer-owned installation and least-privilege policy;
- repository-specific impact and planning without provider source access;
- independent verification around a replaceable coding agent;
- delivery into the customer's normal draft-PR workflow; and
- consented provider visibility into cohort progress and legacy retirement.

## Decision

Lumyn will build a neutral consumer-local runner that can power one
provider-specific update channel at a time.

The first complete product loop is:

```text
provider-confirmed change contract
-> provider-originated change event
-> consumer installation and event-specific authorization
-> consumer-local repository impact inventory
-> reviewable plan
-> deterministic transforms plus replaceable bounded coding agent
-> independent baseline-aware verification
-> tested draft PR
-> consented provider rollout status
```

The initial wedge remains one consequential TypeScript/Node REST/npm SDK
deprecation or sunset. Feature-adoption updates, additional languages, and
neutral multi-provider discovery are expansion paths, not first-MVP scope.

Lumyn operators may assist provider setup and the first consumer cohort. That
is the go-to-market and learning motion, not the product identity. The API
Provider remains the economic buyer. The API Consumer Organization remains the
authority over repository access, model data, commands, credentials, GitHub,
disclosure, review, and merge.

## Provider Change Contract

The user-facing contract records stable change identity, source and target
versions, intended audience and applicability, typed semantics, known
mappings, ambiguity and unsupported cases, verification guidance, provenance,
provider confirmation, and supersession or withdrawal state.

The internal artifact identifier remains `migration-pack` during the v3
compatibility window. The contract is declarative, non-executable, and grants
no consumer repository authority. An accountable Provider Operator confirms
one exact version once; it may then be reused across the invited, consenting
cohort.

## Provider Change Event

A provider event references one exact Provider Change Contract and declares
issuer, API or SDK, audience, deadline, severity, and supersession or
withdrawal state. Duplicate, stale, conflicting, superseded, withdrawn,
wrong-audience, and unauthenticated events fail closed. Unattended write
actions require an authenticated origin.

The MVP may use a bounded configured provider channel rather than a universal
registry or arbitrary changelog monitoring service. The first channel is a
signed, versioned JSON manifest at an exact provider-controlled HTTPS URL. The
Consumer Installation pins that origin and one campaign public key; Lumyn
verifies the detached signature, monotonic sequence, issued/expiry times,
audience, and lifecycle state. The manifest embeds the Provider Change
Contract or names its exact provider-controlled HTTPS URL, and Lumyn verifies
the retrieved bytes against the declared contract digest. The Provider
Operator publishes the manifest. Lumyn may assist setup but cannot originate
it as the provider. An attended file import is recovery and cannot prove
channel delivery or authorize installed-preauthorization writes.

## Consumer Installation

A Consumer Installation is a revocable, expiring consumer-owned policy binding:

- provider or channel;
- repository and selected package root;
- audience and version selectors;
- action ceiling: `notify_only`, `scan_only`, `prepare_patch`, or
  `open_draft_pr`;
- authorization mode: `per_event_approval` or
  `installed_preauthorization`;
- readable and writable paths and commands;
- exact qualified Agent Runner adapter/version, execution-funding mode,
  credential and usage-billing owners, native agent configuration, and
  Agent Runner/model-data/network/credential posture;
- remote-branch and draft-PR scopes;
- provider-visible fields, retention, and deletion.

Every update run freezes an immutable event-specific authorization snapshot
derived from the installation. A provider event may narrow installed authority
but can never widen it. Under `per_event_approval`, the Consumer Maintainer
approves the exact no-write plan. Under `installed_preauthorization`, an
authenticated in-policy event may proceed without a new human approval only
when its plan, paths, commands, model policy, budgets, verification, and action
all satisfy installed policy. Action labels alone grant no side effect. The
installation stores no GitHub token; an approved local or CI broker mints the
short-lived installation token at the remote-delivery step. Anything outside
policy pauses for explicit consumer approval.

## Coding Agent And Verification

The coding agent is a replaceable implementation adapter and the explicit
status-quo comparator, not Lumyn's proprietary moat. Under ADR-0005, the
installation defaults agent execution to disabled. Notify-only, scan-only, and
deterministic-only routes require no runner credential. An agent-assisted route
pauses until the consumer configures one exact qualified adapter/version and
funding route. Codex and Claude Code are launch targets behind one conformance
contract; Cursor is deferred behind the same gate. Every attempt starts a clean
session, and no adapter, model, credential-owner, or billing-owner fallback is
silent.

Agent output remains an untrusted candidate. Deterministic and agent-assisted
candidates use the same exact-head, baseline-aware verification ladder defined
by ADR-0003. Missing business input, ambiguous semantics, inadequate evidence,
scope escape, or exhausted budgets fail closed.

## Delivery And Status

Local evidence, patch, branch, and PR-bundle export remain required recovery
and review surfaces. They do not prove automated delivery.

The first provider campaign must open at least one tested draft PR using a
short-lived, least-privilege token, a non-default branch, evidence-bound
idempotency, and no auto-merge. Manual-only delivery cannot close `EXP-003` or
`PILOT-003`. `EXP-003` requires the composed provider-channel event through
installation, impact, plan, Lumyn deterministic or agent-assisted generation,
independent verification, branch, and draft-PR path; standalone PR creation,
attended event import, and imported manual candidates do not qualify.
The composed path then records an event-bound local status projection or
explicit reporting decline. Provider transmission is separately consented and
cannot block the PR. For `PILOT-003`, the same qualifying composed run must
contain an organically `agent_assisted` plan item on a consumer-selected
qualified runner, pass independent exact-head verification, open the Lumyn
draft PR, and transmit the projection bound to that event, installation
authorization, candidate, verification evidence, and draft PR. Separate agent
and deterministic-delivery runs or deterministic rerouting do not qualify.

A narrowly scoped GitHub App installation may persist at the consumer's
discretion. Its installation tokens are short-lived and least-privilege; a
durable token or broad organization-wide grant is prohibited.

Provider-visible status is event-bound and limited to consumer-consented
fields. Supported rollout states may include `received`, `not_applicable`,
`affected`, `needs_input`, `candidate_ready`, `verified`, `draft_pr_open`,
`accepted`, `merged`, and `retired`. Each state records whether it is directly
observed or consumer-reported. Silence is `unknown`; neither `not_applicable`
nor `unaffected` may be inferred from it, and merge is not inferred to mean
deployed or retired.

## Commercial Proof

The first campaign is provider-paid and services-assisted. The provider must:

- clear the frozen non-refundable prepayment gate;
- name an economic buyer and Provider Operator;
- publish or confirm one consequential change;
- lead distribution and onboarding for five prequalified consumer
  organizations; and
- accept a material rollout, retirement, lead-time, or support-workload
  outcome frozen before execution.

Pilot proof requires three valid repository scans, two tested reviewable
outcomes, at least one consumer-accepted or merged verified outcome, a fair
guide/codemod-plus-generic-agent comparison, and the frozen provider outcome
without worse correction, revert, or false verification. One qualifying
outcome must itself begin at the authenticated provider channel, use installed
preauthorization, contain an organically agent-assisted item on a
consumer-selected qualified runner without bespoke operator edits, pass
independent exact-head verification, become a Lumyn-opened draft PR, and emit
the consented provider-received status projection. A deterministic-only
campaign can prove the provider channel but not `PILOT-003`.

The north star is the share of all Eligible Consumer Units in the frozen
targeted cohort with an accepted or merged verified Lumyn PR by the provider
deadline. Units that do not install, decline authorization, or remain unknown
or blocked stay in the denominator; installation and success-among-installed
are reported separately.

## Deferred

- hosted multitenant coordination or dashboard;
- universal provider registry or arbitrary public changelog monitoring;
- elaborate provider PKI and a public event network;
- long-lived broad GitHub credentials;
- automatic merge or deployment;
- production verification or credentials;
- cross-language and multi-package-root support;
- provider-specific coding-agent codebases;
- annual connected-repository or receipt-backed billing.

## Consequences

Positive:

- Lumyn occupies a distinct provider-to-consumer delivery position rather than
  competing as another migration prompt or coding agent.
- Provider distribution and cohort visibility create value unavailable to a
  one-off maintainer prompt.
- Consumer installation and short-lived delivery grants preserve trust.
- A thin CLI/services-assisted wedge can reach paid learning before hosted
  SaaS.

Costs and risks:

- provider events need authenticity, replay, audience, and supersession rules;
- durable installations risk ambient authority if not narrowly bound;
- draft-PR delivery introduces a real GitHub security boundary;
- provider status can leak consumer information or overstate rollout;
- service labor can mask missing product automation;
- generic agents may still make the code-change layer economically
  substitutable.

## Revisit Triggers

Revisit this decision if providers will not prepay or lead cohort onboarding,
consumers will not install the bounded channel, Lumyn cannot safely open draft
PRs, the generic-agent baseline is materially equivalent end to end, the
provider outcome is immaterial, or delivery COGS prevents a credible
productization path.
