# Lumyn Architecture Guide

## Initial Boundaries

- CLI command surface
- Go module and command package layout
- configuration loading
- schema and artifact models
- redaction pipeline
- evidence normalization
- workflow contract generation
- replay verification
- live known-path verification
- action-boundary checks
- report rendering
- live eval harness

## Early Architecture Rules

- Keep deterministic replay independent of network access.
- Keep redaction before persistence.
- Keep validators over normalized evidence, not raw transcripts.
- Keep live eval diagnostic by default.
- Keep provider adapters behind a small interface.
- Treat ambiguous proof, redaction, credential, stale cassette, undeclared control, or scanner states as fail-closed.
- Keep product artifact namespaces separate from Factory evidence namespaces.

## Systems Thinking Map

- State lives in repo-local config, workflow contracts, cassettes, run artifacts, baselines, schemas, and Factory evidence artifacts.
- Source of truth for MVP scope is `docs/product/prd.md`; source of truth for itemized PRD closure is `.factory/artifacts/prd-to-plan/lumyn-mvp/acceptance-ledger.json`; source of truth for governed delivery state is `.factory/artifacts/prd-to-plan/lumyn-mvp/`.
- Feedback lives in command-result JSON, schema validation, tests, coverage gates, CodeQL status, validation reports, PR lifecycle reports, and scope closure.
- Deleting schemas breaks artifact validation; deleting cassettes breaks replay; deleting redaction rules risks unsafe persistence; deleting Factory evidence breaks governed closure.
- Medium/high risk tasks must record blast radius and rollback or deletion checks in task evidence or review evidence.

## TDD And Red-First Expectations

- Behavior changes should add or update a failing test, fixture, schema example, or validator expectation before implementation when practical.
- Schema/output changes should make invalid output fail before accepting the new implementation.
- If test-first is not practical, the validation report must carry a structured skipped reason.
- Flash/spike work is disposable and cannot satisfy implementation, review, scanner, or validation gates by itself.

## ADR And Decision Triggers

Require an ADR or decision note when a task changes:

- CLI or JSON output contracts
- schema compatibility
- artifact namespace ownership
- redaction, share, live verify, or eval failure semantics
- provider/network/credential posture
- release/distribution posture
- major performance or reliability tradeoffs

## Performance And Cost Triggers

- Replay, report rendering, and eval loops are future hot paths.
- Model-provider eval work must record run count, pass rate, cost, duration, and model snapshot uncertainty.
- Performance validation is reserved until measurable hot paths exist, but tasks touching fan-out, repeated runs, or model calls must declare expected cost impact.

## Reliability And Recovery Triggers

- Retrying, cleanup, persistence, filesystem mutation, redaction, live network, and provider code require hardening checks.
- Live verify must report cleanup success or orphan evidence.
- Replay must remain deterministic without network.
- Remote dependency failures must surface through stable exit codes and command-result JSON.

## Trust-Mode Posture

- Deterministic bootstrap has no network, no ambient secrets, and no live credentials.
- Live sandbox, model-key, and network work is blocked until human approval.
- Approved live work must use explicit config and record credential/network posture in evidence.

## Runtime Shape

- `cmd/lumyn/`: command entrypoint and process exit behavior.
- `internal/result/`: command-result envelope.
- `internal/exitcode/`: stable exit-code constants matching `docs/product/prd.md`.
- `internal/config/`: config discovery and repo-relative path handling.
- `internal/version/`: version metadata.
- `schemas/`: executable JSON Schema contracts for workflow, evidence, cassette, proof, result, redaction, and related artifact models.
- `.github/workflows/validate.yml`: repository delivery gate that runs `make prepush-full` for PR and `main` validation.
- `.github/workflows/codeql.yml`: repository risk-lane security scanner for CodeQL Go analysis.
- `.github/required-checks.json`: branch-protection and `commit-push` status-check contract.
- `.github/CODEOWNERS`: owner-review map for workflow, policy, schema, Factory artifact, CLI, and core implementation paths.
- `.github/action-ref-exceptions.yaml`: audited exception records for non-SHA action refs.
- `.factory/tmp/coverage.out`: ignored local/CI coverage profile emitted by `make test-coverage`.
- `.factory/artifacts/pr-lifecycle/`: Factory delivery evidence tying PR validation, CI/status checks, review, shipping, merge, and post-merge monitoring together.
- `internal/source/docs.go`: docs source walking, operational-guidance checks, and broken local Markdown reference findings.
- `internal/source/docs_test.go`: source-check docs/init/link behavior tests.
- `internal/source/markdownlinks/`: local Markdown link target parsing, fence detection, and missing-reference target normalization for docs source checks.
- `internal/source/report.go`: source-check report persistence and finding-status helpers.
- `internal/source/source_fixtures_*_test.go`: shared OpenAPI/YAML/docs fixtures for source tests.
- `internal/source/source_helpers_test.go`: shared source-test assertion helpers.

## Architecture Budget And Decomposition

Lumyn follows the Factory architecture budget gate: source files warn at `1200`
lines and fail at `2500` lines. The inventory excludes daemon state,
dependencies, caches, and build output, but it does include product source and
tests. The approved current over-budget source surfaces are
`internal/source/source.go`, `internal/source/source_test.go`, and
`scripts/validate_repo_pack.py`, recorded in
`.factory/artifacts/exceptions/architecture-debt-lumyn-source.json` and backed
by `docs/architecture/findings/TEMP_FINDING_2026-06-30_lumyn_arch_budget.md`.

Until that exception is closed, product work that touches `internal/source` or
the repo-pack validator must either reduce file size, split coherent behavior
into smaller packages or files, or record why the change is shrink-neutral with
compensating validation. New feature work must not add unrelated product
domains to the source-ingestion package or validator.

Current source decomposition has started by extracting Markdown link target
parsing and local target normalization into `internal/source/markdownlinks`,
docs-source validation into `internal/source/docs.go`, and report persistence
and finding-status helpers into `internal/source/report.go`. Source tests are
also split so docs/init/link behavior, shared fixtures, and shared assertion
helpers no longer live in one oversized test file. `internal/source` should
keep the project-level source-check orchestration, filesystem resolution, and
shared types while smaller files or internal packages own pure parsing,
validation, normalization, and reporting responsibilities.

Use Go `1.26.4`. T1 stays standard-library-only. T2 introduces the pinned
`github.com/santhosh-tekuri/jsonschema/v5` validator for executable schema
tests.

The T2.5 lifecycle baseline is delivery infrastructure for Factory-governed
work. It includes required-check metadata, workflow hardening, owner-review
coverage, action-ref posture, branch-protection expectations, CI/status,
review, shipping, post-merge, and PR lifecycle evidence. It does not replace
the later product-facing T9 work, which remains
responsible for Lumyn's own GitHub Action/JUnit/PR-comment behavior as an MVP
feature.

The T2.6 dev/architecture propagation baseline is also delivery
infrastructure. It maps the 12-level test matrix, CodeQL scanner posture,
CI lanes, docs parity, output contracts, release integrity, provenance evidence,
architecture boundaries, systems-thinking prompts, TDD/ADR posture,
performance/reliability triggers, failure rules, and evidence requirements into
the execution plan and future task packets before T3 product implementation
starts. It also preserves the Factory planning-skill contract: public API and
contract maps, docs/OSS readiness, minimum-now sequencing, explicit non-goals,
definition of done, Factory compatibility, explicit scope exclusions, alignment
gate refs, plan-drift refs, expanded runtime pins, changelog/versioning intent,
contract impact, ADR posture, TDD-first checks, cost/perf posture, failure
hypotheses, semantic invariants, canonical worker chains, and lifecycle gates.
