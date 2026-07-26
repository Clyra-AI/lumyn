# ADR-0007: Executable V3 Migration Contracts

## Status

Accepted for M2. This decision implements artifact and policy contracts only.
It authorizes no provider fetch, consumer repository access, Agent Runner or
model call, credential issuance, command execution, branch write, PR creation,
provider report, or hosted control-plane behavior.

## Context

The provider-originated update loop crosses independent authorities: provider
change intent, consumer repository policy, optional Agent Runner and model
processing, deterministic verification, GitHub delivery, and consented status
reporting. Prose alone cannot prevent one artifact, actor, status, or payment
relationship from silently widening another.

M2 therefore needs executable contracts before M3 and later milestones build
runtime behavior. JSON Schema can validate artifact shape and local
conditionals, but it cannot by itself prove that an event authorization is a
subset of the installation, that digests agree across independently persisted
artifacts, or that a filesystem path resolves outside the checkout.

## Decision

Every v3 artifact family has a self-contained JSON Schema Draft 2020-12
contract at schema version `1.0`, with one representative valid fixture and at
least one material invalid trust-boundary fixture. Top-level unknown fields are
rejected so misspelled authority, identity, digest, status, privacy, and
provenance fields cannot degrade into ignored metadata. Extension points are
explicit and never grant authority.

The contract families are:

- provider change event and `migration-pack`;
- Consumer Installation, consumer execution manifest, immutable event
  authorization, and managed-credential grant;
- integration graph, impact report, migration plan, and candidate manifest;
- Agent Runner manifest, conformance result, and attempt;
- migration verification, export result, campaign summary, provider status
  projection, and remediation outcome.

The schemas keep impact, route, candidate, verification, delivery, and provider
projection states independent. A stronger value on one axis never implies a
stronger value on another. Provider projection states bind one exact run,
event, Consumer Installation, event authorization, plan, candidate,
verification, delivery, and status-specific evidence chain as far as the
reported state permits. They carry `observed`, `consumer_reported`, or
`unknown` provenance, require explicit evidence for every non-unknown state,
and require a distinct `retirement_confirmation` rather than inferring
`retired` from merge evidence.

Provider inputs remain non-executable data. The first channel contract pins an
exact provider HTTPS origin and campaign key, sequence, issue and expiry time,
audience, retrieved contract-byte digest, signature evidence, and lifecycle
state. Semantic validation requires both the manifest URL and any separately
retrieved Provider Change Contract URL to use that exact normalized origin;
redirects or cross-origin substitutions cannot inherit trust. Attended import
is recovery-only and cannot prove channel delivery or authorize installed
preauthorization.

Consumer installation action modes are ceilings. Persisted notify-only,
scan-only, and patch-preparation installations cannot carry latent command,
mutation, remote-delivery, GitHub-token, Agent Runner, model, registry,
sandbox, or lifecycle-script authority outside their action ceiling. An event
authorization is an immutable snapshot that must bind the installation, event,
and exact plan and may only narrow action, paths, commands, capabilities,
disclosure, and budgets. The Go authorization validator performs these
cross-object subset checks. It also rejects expiry, revocation, wildcard or
traversal scope, stored GitHub or runner/model credentials, provider runner
selection, provider consumer-data access, agent use while disabled,
topology-scope mismatch, and incomplete managed-credential bounds.

Configured Agent Runner policy binds one launch-qualified adapter and version,
canonical executable source/path/digest, conformance digest, permitted
non-interactive auth and entitlement, Agent Runner Vendor, actual Model
Provider and route, funding mode, credential and usage-billing owners, native
configuration digest, clean neutral session, route topology, and no-silent-
fallback posture. Codex and Claude Code are contract targets; neither becomes
qualified until every declared conformance case passes, summary counts are
derived from that exact case set, and an approved live canary carries evidence.
Cursor is rejected until the same gate is separately satisfied.

Consumer execution manifests bind every mount by exact normalized source path,
source digest, source class, target, and mode. Every command names a canonical
absolute program, program digest, and an immutable read-only executable root;
semantic validation rejects programs outside the referenced root. Selecting a
sandbox route additionally requires a separate credential-bearing entrypoint
profile bound to the exact candidate head and read-only candidate mount, exact
entrypoint and working directory, neutral roots, sole sandbox credential,
exact endpoint/operation grant, inherited child and resource limits, and
teardown, cleanup, and orphan evidence. Host-home or OS-credential mounts are
never represented as safe merely by assigning them an approved source class.

Managed credentials bind an approved broker issuer and the exact installation,
event, plan, attempt, runner, and model audience. TTL is at most one hour;
redemption is one-time and attempt-scoped; refresh and cross-attempt reuse are
forbidden; hard token and cost quotas, revocation, reconciliation, and a
vendor-native bound credential or approved enforcing proxy are mandatory.

Consumer-private artifacts resolve through symlinks to an operator-approved
private directory outside the checkout. The initial validator also requires
that directory to deny group and other permissions. Provider-visible payloads
use a strict typed status-projection vocabulary with no arbitrary extension or
free-text surface; unknown nested fields, contradictory status claims, and
secret-shaped values fail closed. Raw source, diffs, patches, prompts,
responses, tool traces, logs, sessions, and credentials are never valid
API-provider projection fields.

The v3 CLI grammar and typed error map are frozen in
`docs/dev/cli-v3-contract.md`. M2 defines artifact and error compatibility; it
does not make an unimplemented command successful. Until its runtime milestone
lands, a command remains unavailable and returns a typed nonzero result under
the retained exit-code contract.

## Compatibility And Validation

- Existing v1 workflow/evidence schemas and historical bytes remain unchanged.
- Each new schema compiles and validates its positive fixture while rejecting
  its negative fixture.
- Cross-contract tests bind event, installation, plan, candidate,
  verification, delivery, and status digests and reject mismatches.
- Semantic tests cover non-widening authorization, disabled/configured agent
  policy, action-ceiling authority, exact provider origin, funding and
  credential ownership, private-root resolution, typed provider-payload
  redaction, complete Agent Runner conformance, exact mount and executable-root
  binding, credential-bearing sandbox isolation, status honesty, prompt-
  instruction fields, default-branch write, and auto-merge denial.
- Promotion evidence is emitted by the trusted validation runner and binds the
  exact base Git SHA, canonical candidate tree digest, declared command,
  runner identity, timestamps, exit status, and hashed stdout/stderr logs.
- `make prepush-full` is the local promotion gate; high-risk M2 additionally
  requires independent security review, CI, passive latest-head Codex review,
  and post-merge monitoring.

## Consequences

The runtime milestones can parse stable, reviewable contracts without gaining
implicit authority. Some invariants deliberately exist in both schema and Go
validation: schema checks protect every persisted artifact, while semantic
checks protect relationships that JSON Schema cannot compare. New contract
versions require compatibility fixtures and migration rules rather than
in-place reinterpretation.

The contract layer is broader than the first runtime slice, but it does not
claim that provider event intake, repository analysis, generation,
verification, delivery, or reporting works end to end. Those claims require
their own acceptance items and direct runtime evidence.
