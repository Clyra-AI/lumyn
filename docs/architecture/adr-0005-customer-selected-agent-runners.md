# ADR-0005: Customer-Selected Agent Runners

## Status

Accepted as a v3.1 planning and trust-boundary refinement. This decision
authorizes no runtime implementation, customer repository access, agent or
model call, credential, network action, branch write, or PR write.

ADR-0004 remains authoritative for the provider-originated product loop.
ADR-0003 remains authoritative for bounded generation, independent
verification, privacy, and human merge control. This ADR makes the replaceable
Agent Runner seam, customer choice, funding, credential ownership, and
qualification rules exact.

## Date

2026-07-25

## Context

Lumyn's value is the trusted application layer around a coding engine, not a
proprietary model. A migration still needs repository-specific reasoning, so
the product must invoke a capable coding agent inside the consumer execution
boundary. Customers may already standardize on Codex or Claude Code, and may
later ask for Cursor or another runner.

Calling the engine “replaceable” is insufficient without answering:

- who selects and authorizes it;
- whether the runner and model vendors are the same party;
- whose account or credential is used and who pays usage;
- which local rules, memories, plugins, or session history can influence it;
- how Lumyn proves that different adapters obey the same boundary;
- what happens when an adapter, model, or credential route is unavailable; and
- which evidence Lumyn requires before independently verifying the candidate.

If those choices remain implicit, an adapter can inherit ambient authority,
change its downstream model route, reuse personal context, report incomplete
provenance, or silently fall back to a different billing or data-processing
boundary. That would make consumer consent and campaign economics ambiguous.

## Decision

Lumyn will implement one adapter-neutral Agent Runner v1 contract. Each
Consumer Installation sets `agent_execution_policy` to `disabled` or
`configured`. `disabled` is the least-privilege default and is sufficient for
notify-only, scan-only, and deterministic-only operation; it grants no Agent
Runner or model authority. If a routed plan later needs `agent_assisted`,
Lumyn pauses until the API Consumer Organization explicitly configures and
authorizes one exact qualified adapter and version. That selection is frozen
into each event-specific authorization snapshot.

Launch adapter targets are:

1. `codex`
2. `claude_code`

Neither is advertised or used live merely because its CLI or SDK exists. Each
exact adapter version and executable digest from an approved source must pass
the common conformance suite and an explicitly approved live canary. Its auth
mode and entitlement class must permit the intended non-interactive local or CI
use under the consumer's vendor agreement and organization policy. `cursor` is
deferred until it passes the same gate.

The first MVP provides controlled installation-time choice among qualified
adapters. It does not provide an arbitrary adapter marketplace, dynamic
mid-run switching, or a public model-provider panel.

## Roles

- The **API Provider** funds the campaign and owns change intent. It cannot
  select an agent for the consumer or access consumer code, context, session,
  or credentials.
- The **API Consumer Organization** owns Agent Runner selection, repository and
  command authority, context disclosure, execution, credential consent,
  revocation, review, and merge.
- The **Agent Runner Vendor** supplies the selected agent harness.
- The **Model Provider** supplies the actual model endpoint or local runtime.
  The Agent Runner Vendor and Model Provider may be the same company, but
  Lumyn records them as separate roles and disclosure boundaries.
- **Lumyn** composes the approved context and tools, enforces the outer
  workspace and authority boundary, normalizes adapter evidence, and performs
  independent verification.

## Agent-Execution, Funding, And Credential Modes

Only an installation with `agent_execution_policy=configured` chooses a
funding and credential mode. A disabled policy requires no runner account,
subscription, credential, or billing owner and cannot be widened by an API
Provider event or plan.

### `consumer_managed`

This is the default configured mode.

- The consumer owns and authorizes its qualifying agent account, enterprise
  subscription, API credential, or local runtime.
- The route exposes the actual Model Provider, endpoint, model, and version and
  permits non-interactive automation. Opaque or changing downstream routing
  does not qualify; the consumer may instead choose a qualifying BYOK, local,
  or Lumyn-managed route.
- The consumer owns third-party usage billing.
- Lumyn receives no reusable credential.
- A required credential is injected only into the task-scoped runner or model
  boundary and is absent from repository commands, evidence, and API
  Provider-visible output.

### `provider_sponsored_lumyn_managed`

This is an optional campaign route, not a fallback.

- The API Provider pays Lumyn for the campaign.
- Lumyn owns the approved Agent Runner/model usage billing.
- The consumer still selects, consents to, and can revoke the exact route.
- Lumyn brokers a short-lived, task-scoped credential only into the isolated
  consumer-local or consumer-CI execution boundary.
- The approved credential broker binds issuer, installation, event, plan,
  attempt, Agent Runner/model audience, and maximum one-hour TTL. One-time
  redemption creates one attempt-scoped session. Multiple in-attempt calls are
  allowed only within hard token/cost quotas; refresh, post-attempt replay, and
  cross-attempt reuse are forbidden. Revocation and post-run usage
  reconciliation are required. If the vendor cannot issue those bounds, Lumyn
  must enforce them through an approved budget-enforcing proxy; otherwise this
  funding mode is unavailable.
- The reusable source credential is never persisted in the repository,
  installation, candidate, logs, or public evidence.
- The API Provider receives no credential, code, prompt, response, tool trace,
  or agent-session access.

Credential owner and usage-billing owner are separate recorded fields. Payment
by the API Provider never creates execution authority.

Every configured action also freezes one `agent_route_topology`:

- `local_runtime` selects no external runner or model scope;
- `runner_mediated` requires Agent Runner network and credential plus model
  request disclosure;
- `direct_model` requires model network, credential, and request disclosure;
  and
- `hybrid` requires both remote minimum sets.

The selected topology's minimum scopes are mandatory. Other topology scopes
are not authorized, and package-registry access remains separately conditional
on the approved plan.

## Adapter Contract

Each adapter must provide or allow Lumyn to derive normalized evidence for:

- adapter identifier and exact version;
- canonical executable path, approved install source, and executable digest;
- conformance-suite version and passing-result digest;
- auth mode and entitlement class without account or secret values;
- Agent Runner Vendor;
- actual Model Provider, endpoint, model route, and model version;
- clean-session identity and attempt identity;
- system policy, prompt, tool-definition, context-selection, and native
  configuration digests;
- startup, structured output, tool calls, file edits, usage, errors,
  cancellation, and exit;
- token, time, retry, attempt, and cost use;
- resulting patch digest; and
- credential owner, usage-billing owner, and execution-funding mode without
  secret values.

An opaque or changing downstream model route, or an auth/subscription mode
that does not permit the intended non-interactive automation, does not qualify
for MVP agent execution.

## Session And Native Configuration Policy

Every attempt resolves the approved executable by canonical path and digest,
then starts a clean, ephemeral session with neutral home and configuration
roots. Lumyn rejects repository-local PATH shadowing and never resumes a
personal, unrelated, or prior migration conversation.

Static user/project instruction rules and memories are disabled or ignored by
default. A consumer may explicitly select a supported static subset. Lumyn
then:

- records its identity and digest;
- treats its content as untrusted repository context;
- subjects its network, tool, credential, and disclosure needs to the same
  grants; and
- rejects any attempt to widen paths, commands, tools, network, credentials,
  budgets, task intent, or approval.

Executable plugins, MCP servers, hooks, and equivalent native extension
surfaces are prohibited for the MVP. A later contract may admit one only after
its executable/source digest, protocol, process isolation, mounts, egress,
credentials, tools, lifecycle, and conformance tests are separately pinned,
qualified, sandboxed, and authorized. Treating executable extension code as
untrusted text is not sufficient.

## Authority And Isolation

Lumyn, not the adapter, controls:

- readable and writable paths;
- isolated workspace and process boundary;
- tool and command allowlists;
- Agent Runner and model network endpoints;
- Agent Runner and model credentials;
- context and disclosure classes;
- file, line, diff, turn, token, time, retry, concurrency, attempt, and cost
  budgets; and
- cancellation and fail-closed cleanup.

The adapter cannot self-approve, self-verify, push a remote branch, create a
PR, merge, or share provider status.

The runner process uses explicit read-only and writable mounts, no host home or
OS credential store, no ambient agent/Docker/service sockets, no unrelated
inherited file descriptors, and host-enforced network egress. It binds the
isolation backend, version, configuration and qualification digests, host
platform, and hard CPU-time, memory, PID, process-tree-depth, disk, and
open-file quotas. Every child process inherits the same limits. Cleanup emits
evidence for process
termination, workspace removal, credential revocation, and mount/socket
absence. If the host cannot enforce that boundary, agent execution fails
closed. Malicious child-process, tool, socket, credential-access, fork-bomb,
and resource-exhaustion tests are part of adapter qualification.

No silent fallback is allowed between adapter, adapter version, Agent Runner
Vendor, Model Provider, model, endpoint, credential owner, or usage-billing
owner. Unavailability, stale conformance, authentication failure, malformed or
partial output, or contract violation blocks the agent route. A separately
authorized deterministic route may still proceed because it does not use the
failed agent authority.

## Qualification

The common suite includes:

- deterministic-fake contract tests;
- exact executable source/path/digest, version, and conformance-digest binding;
- permitted non-interactive auth mode and entitlement class;
- clean-session and no-history-reuse proof;
- normalized structured output and lifecycle events;
- tool-call, edit, usage, error, cancellation, and exit provenance;
- path, command, network, credential, disclosure, budget, and scope denial;
- native-configuration disabled/selected behavior and non-widening;
- malformed, partial, timed-out, and canceled run handling;
- no silent fallback or opaque model routing;
- repository-local executable shadowing and unapproved auto-update denial;
- reusable-credential non-persistence;
- credential-owner and usage-billing-owner attribution; and
- one separately approved live canary for each advertised adapter version.

Passing one adapter does not qualify another. A new version invalidates the
prior qualification unless the adapter policy explicitly permits and proves a
compatible digest-bound range.

## Independent Verification

Agent output, self-reported tests, and runner success status are generation
evidence only. Lumyn's deterministic verifier executes the approved commands
from the exact candidate head in a separate verification view. It has no Agent
Runner or model credential. A future verification mode that requires either
credential needs a separate reviewed contract and cannot satisfy this MVP's
independent-verifier requirement.

The independent verifier starts in a fresh process and verification view,
receives frozen command and verification-configuration digests, and cannot
write through the generation session. Agent Runner/model credentials and
generation-owned evidence handles are absent. Verification evidence is written
only by the verifier/evidence boundary and binds the exact candidate head.

The adapter conformance result proves contract behavior, not migration
correctness. Migration correctness remains bound to the exact candidate and
independent repository/workflow evidence.

## Consequences

Positive:

- notify-only, scan-only, and deterministic-only consumers avoid unnecessary
  agent credentials and vendor onboarding;
- consumers can use a familiar authorized agent without granting the API
  Provider repository access;
- Lumyn stays differentiated by orchestration, authority, evidence, and
  verification rather than model ownership;
- campaign COGS and customer-paid usage remain attributable;
- adapter failures cannot silently change privacy, billing, or model posture;
- the same product contract can qualify additional runners later.

Costs:

- two launch adapters require separate live qualification and ongoing version
  compatibility;
- consumer-managed account and enterprise-policy differences increase
  onboarding test coverage;
- the optional Lumyn-managed route requires a credential broker and clear COGS
  accounting;
- some agent-native configuration must be disabled or explicitly modeled.

## Alternatives Rejected

### Codex-only implementation

Faster initially, but it would turn an implementation choice into an accidental
product dependency and fail the agreed customer-choice posture.

### Arbitrary adapter or model marketplace

Too broad for the first campaign. It multiplies qualification, support,
security, disclosure, and billing combinations before product value is proven.

### Reuse the consumer's existing personal session

Rejected because hidden history, memories, native configuration, and ambient
authority would make runs irreproducible and consent ambiguous.

### API-Provider-owned agent installed into customer repositories

Rejected because provider payment and change authority do not grant repository,
model-context, credential, or execution authority.

### Trust agent-reported tests as verification

Rejected because generation and verification must remain independent and bound
to the exact candidate head.
