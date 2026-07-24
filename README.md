# Lumyn

Lumyn is building provider-sponsored, customer-controlled API migration
automation.

When an API provider introduces a consequential API or SDK change, Lumyn is
designed to:

```text
validate the provider's signed change intent
-> find affected code in an authorized consumer repository
-> prepare a bounded migration
-> verify repository and workflow behavior
-> open an evidence-backed draft PR
```

The API Provider is the initial economic buyer. The API Consumer Organization
retains control of repository access, commands, credentials, disclosure,
review, and merge. Provider sponsorship does not grant access to consumer code,
and Lumyn never auto-merges in the MVP.

The canonical product contract is [docs/product/prd.md](docs/product/prd.md).
The human-readable implementation plan is
[docs/product/plan.md](docs/product/plan.md). The governing architecture
decision is
[ADR-0002](docs/architecture/adr-0002-provider-sponsored-customer-controlled-migrations.md).
The active Factory execution artifacts live under:

```text
.factory/artifacts/prd-to-plan/lumyn-migration-mvp/
```

The prior plan under `.factory/artifacts/prd-to-plan/lumyn-mvp/` is frozen
historical evidence and is not active.

The active PRD and compiled planning package enumerate exactly 62 item-level
acceptance controls.

## Product Wedge

The first supported product is one provider-led migration campaign across
participating GitHub repositories that use one official TypeScript/Node npm
SDK.

The activation wedge is deliberately two-sided and artifact-first: a Provider
Operator creates a standard campaign kit, publishes a signed packet and signed
invitation, and a Consumer Maintainer first enrolls the provider root/package
binding from a separately authenticated bundle and out-of-band confirmed
fingerprint. The maintainer then verifies and accepts the invitation locally.
Enrollment and acceptance run without the checkout mounted, use external
consumer-private state, and produce a reviewable authorization request but no
repository or execution authority. The maintainer must explicitly issue or
revoke a signed authorization bundle. A provider-sponsored connection is
counted only from a separate consumer-signed minimal receipt whose signer is
provider-authenticated and whose one-repository invitation unit receives a
provider-signed, deduplicated acknowledgement; private local use remains
possible without that disclosure. A fully synthetic `lumyn canary run
--offline` must exercise the flow through explicit authorization and a local
draft-PR preview before any provider demonstration can be called working.

The initial deterministic patch classes are:

- SDK method or API-operation rename with an explicit mapping;
- request-property rename or relocation with an explicit mapping;
- response-property rename or relocation with an explicit mapping.

Impact analysis is read-only. A Consumer Maintainer reviews the migration plan
before authorizing changes. Patch execution occurs in an isolated workspace,
and GitHub delivery is a separate explicit draft-PR action.

Every successful verification claim uses a canonical evidence-bound label:
`static_verified`, `repo_verified`, `workflow_contract_replay_passed`,
`workflow_verified_replay`, `workflow_verified_mock`, or
`workflow_verified_sandbox`. A `workflow_verified_*` result requires an
approved entrypoint executed from the exact patched repository head plus
observed interaction and outcome evidence in that environment. An independent
contract or cassette replay is `workflow_contract_replay_passed` and cannot
exceed `repo_verified`.

Unsupported wrappers, dynamic calls, missing business values, unapproved
commands, unsafe redaction, stale evidence, or uncertain verification fail
closed.

## Current Implementation Status

Implemented:

- Go CLI/config/result/exit-code foundation;
- `lumyn init`;
- `lumyn check`;
- OpenAPI and local-doc source parsing, refs, fingerprints, and findings;
- executable workflow, evidence, cassette, trace, proof, boundary, redaction,
  and command-result schemas;
- local validation, coverage, CodeQL, branch-policy, review, and Factory
  delivery controls.

Still planned:

- provider change packets and migration corpus;
- standard campaign-kit creation, configured signing/publication, signed
  invitations, consumer-side `campaign accept`, explicit authorization
  issue/revoke, and minimal signed connection receipts;
- signed offline or exactly authorized online provider-status freshness;
- semantic source/target diffing;
- TypeScript repository impact analysis;
- migration planning and patching;
- repository and workflow verification runtime, including local-only
  `lumyn trace` evidence rendering and the receipt-backed offline canary;
- automatic private-artifact retention/deletion plus `lumyn artifacts gc`
  recovery;
- evidence-backed GitHub draft-PR delivery;
- optional, separately consented provider campaign attestations;
- design-partner campaign validation.

Only `init` and `check` currently have product behavior. Other command names in
the early dispatcher are compatibility placeholders and must not be treated as
working product surfaces until their fail-closed rebaseline task lands.
Migration-plan implementation is gated by M2.5: one qualified API Provider
commitment, a supported-class canary, and a frozen cohort of five distinct
eligible repositories across at least three Consumer Organizations must close
`DISC-001` and `DISC-002`. The provider must also show recurring-value
potential through a second consequential migration/deprecation or an eligible
annual connected-repository program; a one-off bespoke-services opportunity is
not enough.

M0-M4 engineering does not wait for that partner. Pinned, license-compatible
public docs, OpenAPI descriptions, SDK releases, migration guides, and
synthetic fixtures can build the contracts, benchmark, semantic intake, and
read-only impact engine. They remain untrusted engineering inputs and do not
stand in for a signed provider packet, endorsement, sandbox authorization,
reachable cohort, or commercial demand.

## Two-Party Trust Model

- Provider change packets are signed, declarative, and immutable after
  publication for their authorized audience. Immediately before every local,
  sandbox, or remote side effect, the current packet must verify against a
  consumer-pinned API-provider trust root and verified provider-to-package
  binding, including key, timestamp, audience, expiry, rotation, revocation,
  withdrawal, and replay checks. A prior decision is not cached across a write.
  Current status comes only from a signed offline snapshot inside the pinned
  maximum age or an exact endpoint read under `provider_trust_status_read`;
  missing or undeclared freshness blocks.
- First-time provider enrollment cannot trust root material supplied only by
  the invitation. The consumer confirms the root fingerprint through a
  separately authenticated provider channel; normal rotation is signed by the
  active root, and emergency recovery requires explicit re-enrollment. The
  enrollment also pins the provider status signer, receipt-acknowledgement
  signer, and permitted status/receipt exchange classes; invitations may narrow
  but never replace them.
- Repository access is per campaign, per repository, scoped, revocable, and
  time-bounded.
- Repository read, local write, host-isolated command execution,
  provider-status read, package-registry network, sandbox request disclosure,
  sandbox network, sandbox credential, remote branch write, PR write, minimal
  campaign receipt, richer provider attestation, retention, and deletion
  permissions are independent.
- Repository analysis and patching run locally or in consumer-controlled CI by
  default.
- Repository tests run without network and inside a fail-closed isolation
  backend by default. Exact mounts, sanitized environment, child-process
  inheritance, and denial of host home, credential stores, OS credentials,
  agent/Docker/local-service sockets, and extra inherited descriptors are
  required; otherwise the command does not run.
- Automatic `package-lock.json` mutation pins exact Node and npm versions, the
  registry endpoint or immutable snapshot, package-integrity inputs, and the
  toolchain digest; lifecycle scripts remain disabled, and registry network is
  separately authorized.
- Sandbox credentials are non-production and isolated from build/test commands.
  Sandbox consent separately names endpoint and payload classes, synthetic or
  approved non-sensitive test data, provider logging, retention, deletion,
  budgets, namespace, cleanup, and orphan handling.
- Raw code, diffs, logs, traces, prompts, responses, and credentials are not
  provider-visible by default.
- An independent holdout owner provisions and freezes held-out repositories,
  inputs, answer keys, expected patches, expected labels, and raw traces in an
  evaluator-controlled root outside the checkout. The committed manifest is
  non-resolving: opaque case IDs, provenance class and license posture, a
  frozen suite commitment, and encrypted or HMAC artifact commitments only.
  It contains no source location, repository or package identifier, or
  plaintext content digest. M1 creates the commitment in independent provision
  mode; M4/M6/M7 resolve and byte-bind that trusted result in evaluate mode.
  Only `holdout-evaluator` receives the private root.
- Consumer-private runtime artifacts live in an explicitly configured,
  non-committable root outside the checkout and any public source repository.
  They enforce authorization TTL and deletion on expiry or revocation at
  creation, read, startup, and the next run; deletion emits a receipt or orphan
  report, and `lumyn artifacts gc` retries partial cleanup.
  Provider attestations contain only exact consumer-consented fields; public
  evidence requires separate consent and contains only redacted aggregates or
  hashes. Provider export and public commit are irreversible disclosures:
  revocation stops future sharing and removes Lumyn-controlled private copies,
  but cannot recall recipient copies, Git history, clones, or caches.
- Lumyn opens draft PRs only; the Consumer Maintainer owns review and merge.

The current repository and design-partner package are not represented as open
source. Public OSS/self-serve distribution requires a separate approved
license, security, contribution, support, vulnerability-response, and signed
release-integrity gate.

## Repository Layout

- `cmd/lumyn/`: CLI entrypoint
- `internal/`: Go implementation packages
- `tests/`: automated tests
- `schemas/`: versioned executable artifact schemas
- `examples/`: deterministic examples and future migration fixtures
- `workflows/`, `cassettes/`, and `runs/`: retained synthetic or licensed
  project fixtures only; consumer-private runtime instances live outside the
  checkout
- `docs/product/prd.md`: product source of truth
- `docs/product/plan.md`: human-readable execution plan
- `docs/dev/dev_guides.md`: repo engineering and validation contract
- `docs/architecture/architecture_guides.md`: architecture and trust boundaries
- `.factory/artifacts/prd-to-plan/lumyn-migration-mvp/`: active compiled plan
- `.factory/artifacts/prd-to-plan/lumyn-mvp/`: frozen historical plan
- `.factory/artifacts/pr-lifecycle/`: Factory PR lifecycle evidence
- `.factory/artifacts/lifecycle-evidence/`: independent evaluator/reviewer
  evidence; implementation workers cannot write this namespace
- `.github/required-checks.json`: required status-check contract
- `.github/CODEOWNERS`: owner-review coverage

Planned product artifact types such as `changes`, `authorizations`, `impacts`,
and `migrations` are introduced only by their schema-backed implementation
tasks. Their consumer-private instances stay outside the checkout.

## Validation

```bash
make lint-fast
make test-fast
make test-coverage
make test-contracts
make prepush-full
```

`make prepush-full` is the required local gate before PR and merge. GitHub
Actions runs the same gate through `validate` and runs CodeQL through
`CodeQL analyze`.

## Factory Operation

Factory supplies the planning, task-packet, validation, review, shipping, and
evidence contracts. `factoryd` may execute the active Lumyn task graph after
the product rebaseline lands and only when its runtime supports the packet's
independent lifecycle-evidence gates.

Factory's `approval`, `credentials`, and `network` grants govern its
implementation worker. They do not validate or confer Lumyn product authority.
Lumyn enforces the current private product bundle inside the product at each
side-effect boundary; the standalone authorization-validation command is
diagnostic and closure proof.

Safe attended checks:

```bash
export FACTORY_REPO=../factory
factoryd doctor --config .factory/factoryd.example.json --repo lumyn --json
factoryd run --config .factory/factoryd.example.json --repo lumyn --dry-run --json
```

Autoship remains gated by branch protection, required checks, passive Codex
review, merge policy, post-merge monitoring, item-level acceptance evidence,
and every exact Lumyn repository-read, local-write, command, package-registry,
provider-status-read, sandbox-disclosure, sandbox-network, sandbox-credential,
remote-branch, PR-write, campaign-receipt, provider-attestation, retention, and
deletion authorization actually required by the task. Factory's closed worker
`approval`, `credentials`, and `network` grants are separate and never
substitute for those product authorizations.

M8 sandbox verification and M9 draft-PR delivery are independent successors to
M7. A consumer may receive a repository- and mock-verified draft PR without
granting sandbox or provider-reporting authority. Provider reporting remains an
optional, separately authorized action and never blocks PR creation.
