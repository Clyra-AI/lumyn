# ADR-0003: Services-Led Bounded-Agent Migration Execution

## Status

Partially superseded by ADR-0004. No product runtime implementation, model
route, consumer repository access, credential, command, network grant, branch
write, or PR write is authorized by this decision.

ADR-0004 supersedes this ADR's services-led product identity, manual-first
proof, and optional-draft-PR clauses. ADR-0005 refines Agent Runner selection,
qualification, funding, credentials, and session isolation. This ADR remains
authoritative only for bounded-agent execution, deterministic verification,
consumer privacy, and human merge control where those later decisions do not
provide a more specific rule.

## Date

2026-07-24

## Context

ADR-0002 established the useful commercial and trust foundation that an API
Provider can fund a migration campaign while each API Consumer Organization
retains authority over its repository, credentials, execution, disclosure,
review, and merge.

Its initial implementation plan was too narrow in two ways:

1. It treated three deterministic transforms as the complete MVP patch surface,
   even though consequential API and SDK sunsets commonly require
   repository-specific reasoning.
2. It made a standard packet PKI, continuous provider-status protocol,
   connection receipts, and receipt-backed recurring billing prerequisites for
   learning whether Lumyn can deliver verified migrations.

The near-term business is a provider-paid managed sunset campaign. Lumyn can
operate the engagement as a service while preserving consumer-local execution,
bounded authority, deterministic verification, and reusable product evidence.

Coding agents can cover more repository-specific changes than fixed recipes,
but they introduce a separate model egress, credential, privacy, cost,
provenance, prompt-injection, and non-determinism boundary. Agent generation
cannot be treated as verification.

## Decision

Lumyn v3.0 selected services-led, provider-paid API and SDK sunset campaigns
using hybrid patch generation:

- deterministic transforms for exact supported mappings;
- a consumer-local bounded agent for approved repository-specific changes;
- deterministic repository and workflow verification for every candidate;
- patch artifact and PR bundle as recovery and review surfaces;
- local branch, remote branch, and draft PR under progressively separate
  consumer authorization;
- human review and merge.

For active v3.1, ADR-0004 makes the provider-originated installed channel the
product identity and at least one composed tested draft PR mandatory product
proof. Services remain the onboarding motion.

The v3 PRD, product plan, operating docs, repo-local compiled controls, and
external Factory profile are aligned as planning authority. That alignment
does not implement the product or qualify factoryd execution. factoryd remains
paused until its bundle/runtime is qualified against the exact active mission.

## Commercial Motion

The API Provider pays for a managed sunset campaign with:

- a consequential migration and target retirement outcome;
- an accountable provider operator;
- provider-confirmed source/target semantics;
- an intended consumer cohort;
- a commercial scope and decision date;
- consented campaign measurements.

A Lumyn Operator may normalize provider materials, coordinate onboarding,
prepare reviewable planning artifacts, and support campaign execution.
Operator participation is measured as delivery cost and learning. It grants no
ambient access to consumer repositories, model credentials, commands, GitHub,
or merge controls.

The initial commercial proof is a paid, instrumented engagement that produces
useful verified migration outcomes and a credible repeatability path. Annual
connected-repository billing, receipt-backed metering, and software-like gross
margin are later productization hypotheses, not prerequisites for the first v3
campaign.

## Provider Intent

Provider-confirmed migration intent must name:

- provider and accountable confirmer;
- source and target API or SDK versions;
- change semantics and migration guidance;
- deadlines, deprecation/sunset posture, and rollback guidance;
- unresolved questions and unsupported cases;
- provenance and immutable source refs.

A signed declarative provider packet is authoritative when supplied and
confirmed. It remains data and cannot execute code, grant repository authority,
or widen an agent's policy.

V3 defers mandatory:

- consumer-pinned provider root enrollment;
- continuous provider-status endpoint resolution;
- receipt-key binding and provider-signed acknowledgement;
- connection-receipt billing.

Those controls may return when distribution, campaign automation, or billing
requires them. Deferral does not permit unconfirmed or executable provider
input.

## Consumer Authority

The Consumer Organization owns:

- repository read and local write;
- plan approval;
- command execution;
- model request disclosure;
- model endpoint network;
- model credential;
- package registry;
- optional sandbox disclosure/network/credential;
- local and remote branch;
- draft PR;
- provider-visible reporting;
- retention and deletion;
- review and merge.

These are independent, expiring, revocable grants. Provider payment, Lumyn
operator involvement, and plan preparation confer none of them.

## Bounded-Agent Boundary

Agent-assisted generation requires an approved plan item and exact:

- model provider, endpoint, model/version, and parameters;
- prompt, system policy, and tool-definition versions or digests;
- context selection and request-disclosure classes;
- provider logging, training, retention, deletion, and regional posture;
- credential environment, scopes, injection, expiry, and revocation;
- network endpoint and operations;
- read paths, writable paths, and tools;
- file, line, diff, turn, token, time, retry, concurrency, and cost budgets;
- isolated workspace, cancellation, and cleanup;
- request, response, tool-call, usage, attempt, and patch provenance.

Repository content, provider guidance, retrieved context, tool output, and model
output are untrusted. They cannot:

- change system policy;
- widen paths, tools, network, credentials, disclosure, or budget;
- approve a migration plan or authorization;
- access evaluator-controlled holdouts;
- label a patch verified;
- push a remote branch;
- create or update a PR;
- merge.

Budget exhaustion, malformed tool use, prompt injection, missing business
input, redaction uncertainty, or scope ambiguity fails closed.

Agent output is not promised to be byte-identical across runs. Lumyn records
exact provenance and bounded attempts rather than making a false determinism
claim.

## Deterministic Verification

Generation and verification remain independent.

Every deterministic, agent-assisted, or manually adjusted candidate binds:

- provider-confirmed intent;
- repository base and candidate head;
- plan digest;
- generation mode and provenance;
- patch digest;
- pinned verification commands, toolchain, fixtures, and environment;
- evidence artifact hashes.

Verification captures the pre-patch baseline and runs from the exact candidate
head. Canonical labels remain `static_verified`, `repo_verified`,
`workflow_contract_replay_passed`, `workflow_verified_replay`,
`workflow_verified_mock`, and `workflow_verified_sandbox`.

A model completion, agent trace, human glance, or independent contract replay
cannot independently produce a stronger label. Required independent
holdout-evaluator, trace-grader, evidence-attestor, code-review, and human
approval gates remain.

## Delivery Ladder

Delivery states are separate:

1. patch artifact;
2. optional local branch;
3. PR bundle;
4. remote branch when the installed action requires it;
5. draft PR when the installed action requires it;
6. consumer review and merge.

The PR bundle is a consumer-local artifact containing:

- provider-confirmed intent ref;
- plan and base/candidate identities;
- patch and generation provenance;
- deterministic verification evidence;
- excluded, unsupported, and uncertain scope;
- model disclosure and cost summary when applicable;
- residual risks;
- suggested PR title and body.

Patch, local branch, and PR-bundle creation require no GitHub credential.
Remote branch and draft-PR actions require separate exact grants. These
surfaces remain optional per installation, but ADR-0004 requires the first
provider campaign to prove at least one composed tested draft PR. Lumyn never
writes the default branch or auto-merges.

## Data And Disclosure

API-provider disclosure and model-provider disclosure are separate boundaries.

Raw consumer code, diffs, logs, traces, prompts, responses, agent sessions,
credentials, and private evidence are never API-provider-visible. Only
enumerated, consumer-consented campaign status or aggregate fields may cross
that boundary.

Only the exact consumer-authorized request classes may reach the model
endpoint. Secrets, credentials, PII, and production data are prohibited unless
a later approved contract explicitly permits a narrower class.

Consumer-private artifacts live outside the checkout and public source
repository with bounded retention, deletion, and orphan recovery. Committed
evidence uses opaque IDs and digests.

## Runtime And Distribution

The planned runtime remains a Go orchestration core with bounded adapters for:

- provider-intent intake;
- TypeScript analysis;
- deterministic transformation;
- model generation;
- isolated commands;
- deterministic verification;
- PR-bundle rendering;
- optional sandbox and GitHub actions.

Consumer execution is local or in consumer-controlled CI. This decision does
not introduce a hosted consumer-source plane.

Design-partner distribution remains explicitly licensed and integrity-signed.
Public OSS/self-serve and Homebrew require a separate approved license,
security, contribution, support, vulnerability-response, and release-integrity
gate.

## Factory And Lifecycle

The implementation-to-merge chain remains:

1. `task-executor`
2. `validation-gate`
3. `code-review` when required
4. `holdout-evaluator` when selected
5. `trace-grader` when selected
6. `evidence-attestor` when selected
7. `commit-push`
8. `post-merge-monitor`
9. `repair-feedback` on failure

Factory uses only `approval`, `credentials`, and `network` for its workers.
Those grants do not substitute for Lumyn product authority.

The repo-local compiled v3 control set and aligned external profile do not by
themselves authorize an implementation task. factoryd dispatch remains paused
until its runtime is qualified and a bounded task is explicitly authorized and
unpaused.

## Consequences

Positive:

- Lumyn can test a valuable provider-paid service before building elaborate
  campaign infrastructure.
- The bounded agent expands migration coverage without weakening consumer
  authority or proof standards.
- Patch and PR-bundle handoff preserves a no-GitHub recovery path. It does not
  satisfy automated-delivery or commercial-MVP success.
- Provider intent, generation provenance, verification, and delivery remain
  independently inspectable.
- Operator work produces explicit productization evidence instead of being
  hidden.

Costs and risks:

- model egress creates a third-party privacy and credential boundary;
- stochastic generation complicates reproducibility and evaluation;
- services delivery can hide poor product leverage unless operator effort and
  correction are measured;
- prompt injection and tool misuse add new failure modes;
- consumer approval of model disclosure may be harder than local deterministic
  execution;
- broader patch coverage increases review and corpus requirements.

## Rejected Alternatives

### Deterministic-Only MVP

Rejected as the sole v3 execution mode because it cannot cover enough
repository-specific change to test consequential sunset campaigns. Exact
deterministic recipes remain preferred where applicable.

### Unbounded General Coding Agent

Rejected because general repository access, ambient tools, silent model
routing, and self-reported success conflict with Lumyn's authority and evidence
model.

### Provider-Controlled Repository Automation

Rejected because provider payment cannot authorize consumer code, credentials,
commands, branches, PRs, or merge.

### Mandatory GitHub Delivery For Every Installation

Rejected because a verified patch and PR bundle must remain available without
remote credentials. Superseded in part by ADR-0004: draft PR remains optional
per installation, but at least one composed tested draft PR is mandatory for
first-campaign product proof.

### PKI And Receipt Infrastructure First

Deferred because it does not prove migration quality or services demand. Signed
declarative packets remain supported as authoritative confirmed input.

### Automatic Merge

Rejected. Consumer review and merge remain part of the safety contract.

## Rollout And Rollback

Rollout begins with aligned planning documents, this ADR, the regenerated
repo-local v3 Factory control set, and the aligned external Factory profile. No
product implementation or live campaign action occurs in this planning
change. The checked-in factoryd templates remain mission-paused until
bundle/runtime and exact-mission qualification land in separate reviewed work.

Initial implementation uses pinned public/synthetic fixtures, mocked model
routes, negative disclosure and injection fixtures, and evaluator-controlled
holdouts. Live model, consumer repository, sandbox, remote branch, and draft-PR
work require separate approved tasks and exact grants.

Rollback:

- stop v3 dispatch;
- revoke model, repository, command, sandbox, and GitHub grants;
- delete TTL-bound Lumyn-controlled private artifacts;
- preserve historical evidence;
- retain patch/verification evidence only under its approved policy;
- never claim to recall data already sent to an approved external endpoint.
