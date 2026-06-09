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
- Source of truth for MVP scope is `docs/product/prd.md`; source of truth for governed delivery state is `.factory/artifacts/prd-to-plan/lumyn-mvp/`.
- Feedback lives in command-result JSON, schema validation, tests, CodeQL status, validation reports, PR lifecycle reports, and scope closure.
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
- `.factory/artifacts/pr-lifecycle/`: Factory delivery evidence tying PR validation, CI/status checks, review, shipping, merge, and post-merge monitoring together.

Use Go `1.26.4`. T1 stays standard-library-only. T2 introduces the pinned
`github.com/santhosh-tekuri/jsonschema/v5` validator for executable schema
tests.

The T2.5 lifecycle baseline is delivery infrastructure for Factory-governed
work. It does not replace the later product-facing T9 work, which remains
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
definition of done, changelog/versioning intent, contract impact, ADR posture,
TDD-first checks, cost/perf posture, failure hypotheses, and semantic
invariants.
