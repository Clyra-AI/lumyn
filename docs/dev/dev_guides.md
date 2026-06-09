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
- `make test-contracts`: unit tests, Factory planning artifact presence, and required schema-file presence.
- `make prepush-full`: full local gate before PR or merge.

## CI Lane Mapping

- Fast lane: `make lint-fast`, `make test-fast`.
- Core lane: `make test-contracts`, `make prepush-full`.
- Acceptance lane: PRD scope closure in `.factory/artifacts/prd-to-plan/lumyn-mvp/scope-closure-map.json`.
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
| Tier 3 End-to-End | Planned | CLI command invocation tests as `lumyn init`, `check`, `record`, `verify`, `report`, and `eval` mature |
| Tier 4 Acceptance | Planned | `.factory/artifacts/prd-to-plan/lumyn-mvp/scope-closure-map.json` against PRD acceptance groups |
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

## CI And PR Lifecycle

- GitHub Actions workflow: `.github/workflows/validate.yml`.
- Required check name: `validate`.
- CI command: `make prepush-full`.
- PR lifecycle report path: `.factory/artifacts/pr-lifecycle/<task_id>/pr-lifecycle-report.json`.
- Lifecycle-gated tasks require local validation, CI/status evidence, review evidence when required, ship evidence, post-merge evidence, and a PR lifecycle report or an explicit approved exception.
- Passive Codex review settle is required before merge when the repository review integration is enabled.
- Green CI alone is not merge-ready. The latest PR head must have Codex approval, thumbs-up, actionable-resolved, carry-forward, or an approved exception before merge.
- Do not merge manually through `gh pr merge`, the GitHub UI, or a connector before passive Codex review settles.
- A PR merged without latest-head Codex review evidence is a process escape and must be recorded in PR lifecycle evidence with a follow-up fix or blocker.

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
- T1 must use the Go standard library only.
- Any new dependency must be pinned in `go.mod`, justified in the task evidence, and covered by validation.
- Schema/artifact changes must include representative validation coverage in `schemas/`.
- Changes to CI, review, shipping, or post-merge workflow must update `WORKFLOW.md`, this guide, and the relevant Factory planning artifacts in the same branch.

## Docs Parity

- Behavior, flags, output shape, exit codes, artifact paths, install paths, and workflow semantics must update docs in the same change.
- User-facing docs: `README.md`, `docs/product/prd.md`, `docs/dev/dev_guides.md`, and `docs/architecture/architecture_guides.md`.
- Factory planning docs: `.factory/artifacts/prd-to-plan/lumyn-mvp/`.
- Quickstart or example changes must be backed by a command, schema test, fixture, or explicit planned-lane note.

## Output Contracts

- The command-result JSON envelope is a public contract.
- Exit-code constants are a public contract and must match `docs/product/prd.md`.
- Schemas under `schemas/` are executable contracts.
- Product artifacts in `workflows/`, `cassettes/`, `runs/`, `baselines/`, and `examples/` must stay repo-relative and schema-backed when a schema exists.
- Output contract changes require tests or fixtures before implementation when practical.

## Release Integrity

- Primary distribution is a standalone binary.
- Release work is reserved until packaging starts, but any release-sensitive task must define version, changelog, install, artifact-integrity, and UAT evidence before merge.
- Homebrew follows the first binary release; PyPI is not primary.

## Provenance Evidence

- Validation report path: `.factory/artifacts/task-runs/<task_id>/validation-report.json`.
- Work-proof marker path: `.factory/artifacts/task-runs/<task_id>/work-proof-marker.json`.
- PR lifecycle report path: `.factory/artifacts/pr-lifecycle/<task_id>/pr-lifecycle-report.json`.
- Codex review evidence source: latest-head PR review, thumbs-up, actionable-resolved, carry-forward, or approved exception before merge.
- Scanner evidence source: GitHub Actions `CodeQL analyze`.
- Scope closure source: `.factory/artifacts/prd-to-plan/lumyn-mvp/scope-closure-map.json`.
- Evidence must use repo-relative paths and record skipped-command reasons.

## Distribution Pins

- Primary: standalone binary.
- Secondary: Homebrew.
- Not primary: PyPI.

## Provider Pins

- First eval provider: OpenAI-compatible HTTP adapter.
- Provider config shape: `provider`, `model`, `temperature`, `base_url`, `api_key_env`.
- Eval provider work is blocked until deterministic replay foundation passes and model-key/network posture is approved.
