# Lumyn Workflow Contract

Version: 3.0
Status: Normative v3 planning contract; product runtime not implemented

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
the exact v3 mission.

The prior `.factory/artifacts/prd-to-plan/lumyn-mvp/` package and its task,
pilot, and PR-lifecycle evidence are immutable historical records.

## Product Delivery Flow

The services-led v3 workflow preserves these state boundaries:

```text
provider-paid sunset campaign
-> provider-confirmed migration intent
-> [signed declarative provider packet, optional authoritative input]
-> consumer campaign invitation
-> read-only impact authorization
-> customer-authorized repository analysis
-> no-write migration plan
-> consumer plan approval
-> explicit action-specific consumer execution grants
-> [deterministic transformation | bounded-agent generation | needs input]
-> isolated local patch
-> deterministic repository and workflow verification
-> patch artifact + optional local branch + PR bundle
-> [consumer-authorized remote branch, optional]
-> [consumer-authorized draft PR, optional]
-> consumer review and merge
-> [separately consented provider status, optional]
```

A Lumyn Operator may coordinate campaign intake, provider confirmation,
consumer onboarding, and review. Operator participation grants no ambient
repository, command, model, credential, branch, PR, or merge authority.

A signed declarative provider packet is authoritative when available and
confirmed. It remains untrusted executable input and cannot widen authority.
Elaborate packet PKI, continuous status resolution, connection receipts,
provider acknowledgements, and receipt-backed billing are deferred from the v3
managed-campaign prerequisite.

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

Acceptance item IDs remain the closure source. Tasks, waves, and delivery slices
are sequencing and coverage lenses only.

## Approval Gates

### Product Planning

- Plan approval is required before implementation or product execution.
- The PRD, plan, ADRs, operating docs, Factory profile, active compiled
  artifacts, validators, examples, and fixtures must agree before v3 dispatch.
- The repo-local v3 compiled set does not close the separate external-profile
  and factoryd-runtime qualification gate.
- Services-led discovery may collect provider intent and anonymized planning
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

- an approved plan item explicitly routed to `bounded_agent`;
- exact model provider, endpoint, model/version, and parameters;
- exact prompt, system policy, and tool-definition versions or digests;
- allowed read paths, writable paths, and tool calls;
- file, line, diff, turn, token, time, retry, and cost budgets;
- isolated workspace and fail-closed cancellation;
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

### Model Request Disclosure

`model_request_disclosure` names:

- exact source, metadata, and context classes permitted to leave the consumer
  plane;
- prohibited code, secrets, credentials, PII, and production data;
- redaction and minimization;
- model-provider logging, training, retention, deletion, and regional posture;
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

No generic internet access, fallback endpoint, silent model upgrade, or
undeclared tool network is allowed.

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
- OS credential access, denied by default;
- a supported fail-closed isolation backend.

Repository commands default to no network, no lifecycle scripts, no host home
or credential stores, no agent/Docker/unrelated service sockets, and no extra
inherited descriptors. If the boundary cannot be enforced, the command does not
run.

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
4. optional remote branch;
5. optional draft PR;
6. consumer review and merge.

The PR bundle contains the provider-confirmed intent ref, plan digest, base/head
identity, patch provenance, model provenance when applicable, deterministic
verification evidence, excluded/unsupported scope, residual risk, and suggested
title/body. It requires no GitHub credential.

`github_branch_write` names the repository, authorized non-default namespace,
base commit, token expiry, idempotency, and rollback. It grants no PR authority.

`github_pr_write` names the repository, authorized remote branch, base branch,
draft-only posture, token expiry, idempotency key, and approved plan/evidence
refs. It grants no remote-branch or merge authority.

The product never writes the default branch or auto-merges.

### Optional Provider Reporting

`provider_reporting` lists only the exact fields the consumer permits Lumyn to
share. Raw source, diffs, logs, traces, prompts, responses, and credentials are
excluded by default. Missing reporting consent does not block local patch,
branch, PR-bundle, or otherwise authorized draft-PR delivery.

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
- Risk: GitHub Actions `CodeQL analyze`
- Acceptance: item-level ledger and closure map
- Cross-system: separately approved model, sandbox, or GitHub integration

## PR Lifecycle Baseline

- Local gate: `make prepush-full`.
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
