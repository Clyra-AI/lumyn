# Lumyn Developer Guide

## Toolchain Pins

| Tool | Version |
|---|---:|
| Go | `1.26.4` |

## Dependency Pins

- `github.com/santhosh-tekuri/jsonschema/v5 v5.3.1`: executable JSON Schema validation for T2 and later schema/artifact work.

## Validation Matrix

- `make lint-fast`: repo operating pack and layout checks.
- `make test-fast`: Go unit tests.
- `make test-coverage`: Go coverage gate over first-party CLI/schema packages.
- `make test-contracts`: unit tests, Factory planning artifact presence, downstream pilot evidence validation, repo-pack guide propagation validation, and required schema-file presence.
- `make prepush-full`: full local gate before PR or merge.
- `make audit-remote-protection`: networked GitHub audit for live `main` branch protection and the `protect-main-from-direct-push` ruleset.

## CI Lane Mapping

- Fast lane: `make lint-fast`, `make test-fast`.
- Core lane: `make test-contracts`, `make prepush-full`.
- Acceptance lane: PRD item coverage in `.factory/artifacts/prd-to-plan/lumyn-mvp/acceptance-ledger.json` and item status in `.factory/artifacts/prd-to-plan/lumyn-mvp/scope-closure-map.json`.
- Cross-platform lane: reserved until standalone binary packaging starts.
- Risk lane: GitHub Actions `CodeQL analyze`, plus future hardening/chaos/perf checks as replay, live verify, share, and eval mature.
- Release lane: reserved until standalone binary release packaging starts.

## 12-Level Test Matrix

Lumyn preserves the Factory 12-level test matrix even when a tier is planned,
reserved, or blocked by credential/network approval.

| Tier | Status | Current command, check, or evidence |
|---|---|---|
| Tier 1 Unit | Active | `go test ./... -count=1` through `make test-fast` |
| Tier 2 Integration | Planned | `make test-contracts`; grows with schema/workflow integration tests |
| Tier 3 End-to-End | Planned | CLI command invocation tests as `lumyn init`, `check`, `record`, `verify`, `trace`, and `eval` mature |
| Tier 4 Acceptance | Planned | `.factory/artifacts/prd-to-plan/lumyn-mvp/acceptance-ledger.json` and `.factory/artifacts/prd-to-plan/lumyn-mvp/scope-closure-map.json` against PRD acceptance item IDs |
| Tier 5 Hardening | Planned | fail-closed, redaction, stale-cassette, cleanup, retry, and orphan-evidence tests |
| Tier 6 Chaos | Reserved | controlled failure injection after live verify/retry boundaries exist |
| Tier 7 Performance | Reserved | runtime, cost, and duration budgets after replay/eval paths exist |
| Tier 8 Soak | Reserved | repeated replay/eval stability after deterministic replay exists |
| Tier 9 Contract | Active | `make test-contracts`; schemas, command-result envelope, exit-code stability |
| Tier 10 UAT | Reserved | install/distribution acceptance after first standalone binary packaging |
| Tier 11 Scenario | Planned | planted-flaw workflow matrix and proof-honesty scenario coverage |
| Tier 12 Cross-System Integration | Blocked until approved | live sandbox/model-provider checks after explicit credential and network approval |

Future task packets must cite the applicable tiers or record an approved
non-applicable reason.

## Coverage Gates

Lumyn inherits the org-wide Factory coverage policy derived from Wrkr's launch
bar.

| Scope | Threshold | Enforcement |
|---|---:|---|
| Go first-party packages overall (`cmd/`, `internal/`, `schemas/`) | `>= 75%` | `make test-coverage`, included in `make prepush-full` and CI |
| Go stable command or core packages | `>= 85%` | Required once the package contains stable product-bearing logic; exceptions must name owner, reason, expiry or follow-up task, and compensating validation |

Coverage output is written to `.factory/tmp/coverage.out`. The gate is a
false-green guard, not a replacement for command-result schema tests, acceptance
closure, replay scenarios, CodeQL, or future live eval validation.

Future task packets that touch first-party code, tests, CI, generated code, or
package boundaries must cite `coverage_policy_refs` or record an approved
coverage exception.

## Architecture Budget Gate

Lumyn uses the Factory default architecture budget: warn at `1200` source lines
and fail at `2500` source lines for `.go`, `.py`, `.ts`, `.tsx`, `.js`, and
`.jsx` files, excluding generated runtime, dependency, cache, and build
directories. `factoryd doctor` must emit an `architecture_budget_report` before
daemon work. The remaining warning-level validator orchestration file is
`scripts/validate_repo_pack.py`, covered by
`.factory/artifacts/exceptions/architecture-debt-lumyn-source.json`; work that
touches it must avoid net growth unless it is shrink-only decomposition.
Architecture-budget logic lives in `scripts/repo_pack_architecture.py`, and
self-test fixture construction lives in `scripts/repo_pack_self_test.py`.
`internal/source` has been split below the warning threshold and future work
must preserve those smaller responsibility boundaries.

## CI And PR Lifecycle

- GitHub Actions workflow: `.github/workflows/validate.yml`.
- Required check name: `validate`.
- CI command: `make prepush-full`.
- Required-check manifest: `.github/required-checks.json`, expected checks `validate` and `CodeQL analyze`.
- Owner-review coverage: `.github/CODEOWNERS`, covering workflow, policy, Factory artifact, schema, CLI, and core implementation paths.
- Workflow hardening: validate and CodeQL workflows declare least-privilege permissions, concurrency cancellation, job timeouts, and toolchain setup from pinned repo files.
- Action-ref posture: `.github/action-ref-exceptions.yaml` records audited exceptions for non-SHA action refs with owner, reason, scope, expiry, and review command.
- PR lifecycle report path: `.factory/artifacts/pr-lifecycle/<work_item_id>/pr-lifecycle-report.json`.
- Lifecycle-gated tasks require local validation, CI/status evidence, review evidence when required, ship evidence, post-merge evidence, and a PR lifecycle report or an explicit approved exception.
- Passive Codex review settle is required before merge when the repository review integration is enabled.
- Green CI alone is not merge-ready. The latest PR head must have Codex approval, thumbs-up, actionable-resolved, carry-forward, or an approved exception before merge.
- Do not merge manually through `gh pr merge`, the GitHub UI, or a connector before passive Codex review settles. A configured `factoryd` autoship run may use its `github_cli` provider only after required CI, passive Codex review, merge, post-merge, and semantic scope-closure gates pass.
- A PR merged without latest-head Codex review evidence is a process escape and must be recorded in PR lifecycle evidence with a follow-up fix or blocker.
- GitHub `main` must be protected by branch protection plus the `protect-main-from-direct-push` ruleset.
- Required live controls: pull requests required, strict `validate` and `CodeQL analyze` status checks, admin enforcement, conversation resolution, no force pushes, no branch deletion, and no current-user ruleset bypass.
- Verify live remote controls with `make audit-remote-protection` when GitHub credentials are available.

## Security Scanner Enforcement

- Scanner: CodeQL.
- Workflow: `.github/workflows/codeql.yml`.
- Status source: GitHub Actions `CodeQL analyze`.
- Required for: dependency additions, generated-code intake, CI/workflow changes, redaction/share/live/eval/provider code, external calls, data exposure, and release-sensitive work.
- Exception behavior: blocked unless the task packet, validation report, and PR lifecycle report record an approved scanner exception with compensating validation.

## Bootstrap Rules

- Deterministic bootstrap must not require network, sandbox credentials, or model keys.
- Live sandbox and eval work require explicit human approval before credentials are introduced.
- Tests should be added before implementation when practical.
- Evidence artifacts must use repo-relative paths.
- Use Factory `task-supervisor` for guided post-PRD, audit, review,
  recommendation, or idea intake before implementation. The report path is
  `.factory/artifacts/task-supervisor-runs/<mission>/<timestamp>.json`; it
  records source validation, ingest, doctor, dry-run, alignment gates, and the
  recommended next task.
- Runner-ready task packets must preserve `semantic_invariants` so workers,
  supervisors, and review repairs keep behavior boundaries intact across
  decomposition and follow-up tasks.
- T1 must use the Go standard library only.
- Any new dependency must be pinned in `go.mod`, justified in the task evidence, and covered by validation.
- Schema/artifact changes must include representative validation coverage in `schemas/`.
- Changes to CI, review, shipping, or post-merge workflow must update `WORKFLOW.md`, this guide, and the relevant Factory planning artifacts in the same branch.
- Changes to CI, workflow, policy, or lifecycle gates must keep
  `.github/required-checks.json`, `.github/CODEOWNERS`,
  `.github/action-ref-exceptions.yaml`, `scripts/validate_repo_pack.py`, and
  the workflow files aligned.
- Changes that affect T3+ task planning must keep `scripts/validate_repo_pack.py`
  green so product task packets preserve CI lanes, 12-level test matrix refs,
  scanner gates, engineering policy refs, architecture guidance refs, and the
  Factory `prd-to-plan` / `execution-compiler` fields for Factory
  compatibility, explicit scope exclusions, alignment gate refs, plan-drift
  refs, acceptance-ledger refs, acceptance item IDs, expanded runtime pins,
  changelog intent, contract impact, ADR posture, TDD-first evidence,
  cost/perf impact, failure hypotheses, semantic invariants, canonical worker
  chains, and lifecycle gates.

## Docs Parity

- Behavior, flags, output shape, exit codes, artifact paths, install paths, and workflow semantics must update docs in the same change.
- User-facing docs: `README.md`, `docs/product/prd.md`, `docs/dev/dev_guides.md`, and `docs/architecture/architecture_guides.md`.
- Factory planning docs: `.factory/artifacts/prd-to-plan/lumyn-mvp/`.
- Quickstart or example changes must be backed by a command, schema test, fixture, or explicit planned-lane note.

## Output Contracts

- The command-result JSON envelope is a public contract.
- Exit-code constants are a public contract and must match `docs/product/prd.md`.
- Schemas under `schemas/` are executable contracts.
- Normalized result and failure contracts must carry `finding_kind`,
  `proof_strength`, `action_boundary_status`, `security_relevance`,
  `fix_target`, `surface_fingerprint`, `eval_mode`, `provider_metadata`, and
  `corpus_eligible: false` as local-only evidence fields.
- Product artifacts in `workflows/`, `cassettes/`, `runs/`, `baselines/`, and `examples/` must stay repo-relative and schema-backed when a schema exists.
- Output contract changes require tests or fixtures before implementation when practical.

## Agent-Native CLI Policy

- Agent-facing commands must support stable JSON output mode.
- Commands should emit machine-readable output when stdout is not a TTY unless
  explicitly human-only with an approved exception.
- `--quiet` and `--compact` must preserve status, evidence refs, typed errors,
  and the command-result envelope.
- T3+ task packets must carry acceptance checks for JSON mode, piped or
  non-interactive behavior, quiet/compact output posture, typed exits, and
  machine-readable errors when the task touches CLI behavior.

## Source-Check Proof Rubric

`lumyn check` must stay distinct from generic generation-time readiness scoring.
It verifies whether source evidence can support later workflow-completion proof.
Source checks use two tiers:

- Infrastructure tier: parser/schema/readability facts such as path existence,
  operation IDs, auth declaration shape, response schema availability, and
  validator input paths.
- Workflow/domain tier: evidence that a source surface can support a concrete
  workflow job, expected outcome, read-back, cleanup, boundary, or fix target.

Unprovable, missing, or hallucinated paths are hard failures for source-evidence
closure. A check may report `needs_user_input` or `coverage_gap`, but it must not
count those as proof claims or as accepted recorder-quality claims.

T3.1 is the parser/source-proof repair guard. Recorder-heavy tasks T4.1 through
T4.3 must consume only source-check outputs that passed the structured parser,
proof-tier, and hallucinated-reference checks or record an approved exception.

## Workflow Insight Framing

Recorder and report work should capture the workflow's business job in addition
to endpoint mechanics. Use an optional field such as `workflow_insight` or
`business_job` to explain what the workflow enables, why the evidence matters,
and which fix target unlocks completion. This field is evidence framing, not
marketing copy; it must be grounded in the workflow goal, expected outcome,
validators, trace, or source references.

## Structured Data, Proof, Budgets, And Redaction

- OpenAPI, workflow, cassette, evidence, trace, baseline, and report data must
  be read through parsers, schemas, or stable APIs.
- Source checks must distinguish syntax proof, source-evidence proof, workflow
  proof, and user-visible proof; workflow or user-visible closure requires a
  proof-of-behavior scorecard or approved exception.
- Large logs, traces, reports, and generated evidence must be cited by artifact
  ref with full-output hashes and truncation metadata instead of duplicated
  payloads.
- `lumyn share` and any customer-safe artifact must recursively redact nested
  owner, credential, endpoint, secret, and machine-local path fields.

## Release Integrity

- Primary distribution is a standalone binary.
- Release work is reserved until packaging starts, but any release-sensitive task must define version, changelog, install, artifact-integrity, and UAT evidence before merge.
- Homebrew follows the first binary release; PyPI is not primary.

## Provenance Evidence

- Validation report path: `.factory/artifacts/task-runs/<task_id>/validation-report.json`.
- Work-proof marker path: `.factory/artifacts/task-runs/<task_id>/work-proof-marker.json`.
- PR lifecycle report path: `.factory/artifacts/pr-lifecycle/<work_item_id>/pr-lifecycle-report.json`.
- Codex review evidence source: latest-head PR review, thumbs-up, actionable-resolved, carry-forward, or approved exception before merge.
- Scanner evidence source: GitHub Actions `CodeQL analyze`.
- Branch protection evidence path: `.factory/artifacts/repo-controls/main-branch-protection.json`.
- Pilot closure path: `.factory/artifacts/pilot/lumyn-mvp-slice/scope-closure-report.json`.
- Pilot repair task path: `.factory/artifacts/pilot/lumyn-mvp-slice/repair-loop/task-packet.json`.
- Acceptance item source: `.factory/artifacts/prd-to-plan/lumyn-mvp/acceptance-ledger.json`.
- Scope closure source: `.factory/artifacts/prd-to-plan/lumyn-mvp/scope-closure-map.json`.
- Repo-pack propagation validator: `scripts/validate_repo_pack.py`.
- Evidence must use repo-relative paths and record skipped-command reasons.

## Distribution Pins

- Primary: standalone binary.
- Secondary: Homebrew.
- Not primary: PyPI.

## Provider Pins

- MVP eval providers: OpenAI-compatible HTTP and Anthropic Messages HTTP adapters.
- Provider config shape: `provider`, `model`, `temperature`, `base_url`, `api_key_env`.
- Local open-source model servers are represented as OpenAI-compatible
  `base_url` endpoints. They are provider endpoints, not bundled model
  artifacts, and they require the same explicit `model_provider_endpoint` grant
  before live eval closure.
- Eval provider work is blocked until deterministic replay foundation passes and
  the task has a complete `model_provider_endpoint` grant naming provider
  identity, provider model, endpoint or `base_url`, credential environment,
  budget posture, redaction posture, and allowlist. Generic network or
  credential approval does not satisfy that model-specific gate.
