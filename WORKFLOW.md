# Lumyn Workflow Contract

Version: 0.1
Status: Normative

## Work Signal

Lumyn accepts work from:

- [docs/product/prd.md](docs/product/prd.md)
- Factory task packets under `.factory/artifacts/`
- GitHub issues or pull requests once the repo is active

## Normal Factory Chain

1. `scout-context`
2. `execution-compiler`
3. `task-executor`
4. `validation-gate`
5. `commit-push`
6. `post-merge-monitor`

## Approval Gates

- Plan approval is required before implementation work starts.
- Merge approval is required before protected branch updates.
- Human approval is required before any live sandbox credential, model key, or non-deterministic eval is used.
- Human approval is required before changing the pinned runtime, module path, or primary distribution target.

## Artifact Rules

- Durable Factory artifacts go under `.factory/artifacts/`.
- Safe `factoryd` config template lives at `.factory/factoryd.example.json`.
- Explicit full-loop `factoryd` config template lives at `.factory/factoryd.autoship.example.json`.
- Ignored daemon runtime state goes under `.factoryd/`.
- The canonical MVP PRD path is `docs/product/prd.md`.
- Factory references must use repo-relative paths, never machine-local absolute paths.
- Validation evidence must include command, status, artifact refs, and skipped-command reasons when applicable.

## Bootstrap Validation Lanes

- Fast lane: `make lint-fast`, `make test-fast`
- Coverage lane: `make test-coverage`
- Contract lane: `make test-contracts`
- Full lane: `make prepush-full`

## PR Lifecycle Baseline

- Local validation gate: `make prepush-full`.
- CI status check: GitHub Actions workflow `.github/workflows/validate.yml`, check name `validate`, runs `make prepush-full`.
- Security scanner: GitHub Actions workflow `.github/workflows/codeql.yml`, status source `CodeQL analyze`, required for dependency additions, generated-code intake, CI/workflow changes, external calls, redaction/share/live/eval/provider code, data exposure, and release-sensitive work.
- Coverage gate: `make test-coverage`, included in `make prepush-full`, required for first-party code, tests, CI, generated code, or package-boundary work unless an approved scoped exception is recorded.
- Required-check manifest: `.github/required-checks.json` must name `validate` and `CodeQL analyze`.
- Owner-review policy: `.github/CODEOWNERS` must cover workflow, policy, schema, Factory artifact, CLI, and core implementation paths.
- Workflow hardening: GitHub Actions workflows must declare least-privilege permissions, concurrency cancellation, job timeouts, and toolchain setup from pinned repo files.
- Action-ref posture: non-SHA GitHub Action refs require audited exception records in `.github/action-ref-exceptions.yaml`.
- Structured review: `code-review` must produce review evidence only when risk, workflow policy, validation findings, pre-release posture, or explicit review policy require it.
- Shipping evidence: `commit-push` must produce or reference a ship packet before merge.
- Post-merge monitoring: default branch health must be checked after merge and recorded when the task requires lifecycle evidence.
- PR lifecycle report path: `.factory/artifacts/pr-lifecycle/<work_item_id>/pr-lifecycle-report.json`.
- Unavailable CI or review gates require an explicit approved exception; they are not silently treated as passed.
- Passive Codex review settle is required before merge when the repository review integration is enabled.
- Green CI alone is not merge-ready. The latest PR head must have Codex approval, thumbs-up, actionable-resolved, carry-forward, or an approved exception before merge.
- Do not merge manually through `gh pr merge`, the GitHub UI, or a connector before passive Codex review settles. A configured `factoryd` autoship run may use its `github_cli` provider only after required CI, passive Codex review, merge, post-merge, and semantic scope-closure gates pass.
- A PR merged without latest-head Codex review evidence is a process escape and must be recorded in PR lifecycle evidence with a follow-up fix or blocker.
- Remote branch protection: GitHub `main` must be protected by branch protection plus the `protect-main-from-direct-push` ruleset.
- Required live controls: pull requests required, strict `validate` and `CodeQL analyze` status checks, admin enforcement, conversation resolution, no force pushes, no branch deletion, and no current-user ruleset bypass.
- Remote audit command: `make audit-remote-protection` when GitHub credentials are available.

## Runtime And Distribution Pins

- Go version: `1.26.4`
- Module path: `github.com/Clyra-AI/lumyn`
- Primary distribution: standalone binary
- Secondary distribution: Homebrew
- Non-primary distribution: PyPI
- MVP eval providers: OpenAI-compatible HTTP and Anthropic Messages HTTP adapters

## Stop Conditions

Stop and request human decision if:

- a task requires live credentials before credential posture is approved
- a task requires network access during deterministic bootstrap
- a proposed change writes outside declared allowed paths
- redaction confidence is unknown for persisted artifacts
- implementation would satisfy a command while violating explicit PRD scope exclusions
- lifecycle gates are required but ship, CI/status, passive Codex review, post-merge, or PR lifecycle evidence is missing without an approved exception
- structured `code-review` is required by risk, pre-release posture, validation findings, or review policy but review evidence is missing without an approved exception
- a PR is merge-ready by CI but lacks latest-head passive Codex review settle evidence
- GitHub `main` branch protection or the `protect-main-from-direct-push` ruleset is missing, disabled, or bypassable
- required-check metadata, owner-review coverage, workflow permissions,
  concurrency, timeouts, toolchain-pin setup, or action-ref exception posture is
  missing for CI/workflow/policy/lifecycle work
- scanner-gated work lacks CodeQL status evidence or an approved scanner exception
- a product task omits required test-matrix refs, scanner gates, or architecture guidance refs inherited from the repo operating pack
- a product task omits required CI lane refs, docs parity, output contract, release integrity, provenance, systems-thinking, TDD, ADR, performance, reliability, or fail-closed refs inherited from the repo operating pack
- a product task omits Factory `prd-to-plan` / `execution-compiler` fields for Factory compatibility, explicit scope exclusions, alignment gate refs, plan-drift refs, runtime pins, slice rationale, changelog/versioning intent, contract impact, ADR posture, TDD-first evidence, cost/perf impact, failure hypotheses, semantic invariants, canonical worker chains, or lifecycle gates
- a runner-ready product task omits `worker_type`, `factoryd_runtime`, `validation_commands`, `evidence_required`, or `stop_conditions`
