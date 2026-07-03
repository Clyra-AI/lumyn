# TEMP Finding: Lumyn Architecture Budget And Early Source Split

Date: 2026-06-30
Status: closed source finding
Repo: Lumyn

## Boundary

This file is repo-local source evidence retained for architecture audit history.
It is not a generated execution contract, task packet, acceptance ledger, or
scope-closure artifact. The finding was promoted through the governed Factory
path and closed by the systemic architecture/review-convergence ledger in
Factory PR #240, with Lumyn evidence synchronized through PR #58.

## Objective

Adopt architecture budget rules early so Lumyn does not repeat Relia's monolith
pattern as product behavior expands.

## Closure

Closed. Lumyn's repo operating pack now records the architecture budget and
exception path, `internal/source` has been split below the warning threshold,
and the validator orchestration surface is below its tracked shrink-only
ceiling:

- `internal/source/source.go`: 277 lines.
- `internal/source/source_test.go`: 1,029 lines.
- `scripts/validate_repo_pack.py`: 879 lines.

## Current Finding

Lumyn already has a thin `cmd/lumyn` entrypoint and several implementation
packages under `internal/`. The immediate risk is not the CLI entrypoint. The
first `internal/source` decomposition pass has brought source package files
below the warning threshold. The repo-pack validator has also been split below
the warning threshold, with `scripts/validate_repo_pack.py` still tracked under
a shrink-only orchestration ceiling. Source tests and validator self-tests now
live in focused files so future behavior can be split without rebuilding a
single test monolith.

## Workstream A: Repo-Pack Adoption

1. Update `AGENTS.md`, `docs/dev/dev_guides.md`, and
   `docs/architecture/architecture_guides.md` with explicit architecture
   budgets.
2. Add the shared architecture fitness gate once Factory/factoryd exposes it.
3. Require task packets touching oversized files to shrink, split, or carry an
   approved architecture-debt exception.
4. Keep artifact namespaces and deterministic replay boundaries separate from
   Factory delivery evidence.

## Workstream B: Early Source Split

1. Split `internal/source` by responsibility before more product work lands.
2. Keep parsing, normalization, validation, evidence mapping, and report-facing
   logic separate.
3. Keep `source_test.go` below budget by moving tests near the affected package
   responsibilities.
4. Preserve public CLI behavior, command-result JSON, schemas, and examples.

Current progress:

- `internal/source/fingerprint.go` now owns source hashing, generated-artifact
  skip rules for docs hashing, and source-surface fingerprint assembly,
  lowering `internal/source/source.go` while preserving source-check behavior.
- `internal/source/config.go` now owns source config read/write/defaulting and
  repo-local config validation, further lowering `internal/source/source.go`
  while preserving init/check behavior.
- `internal/source/openapi.go` now owns OpenAPI source parsing, operation
  metadata checks, schema/reference resolution, and parser-facing findings,
  bringing `internal/source/source.go` below the warning threshold while
  preserving source-check behavior.
- `internal/source/openapi_auth.go` now owns OpenAPI security-scheme and OAuth
  scope description checks so the OpenAPI split remains below the warning
  threshold.
- `scripts/repo_pack_architecture.py` now owns architecture-budget policy,
  exception, inventory, and line-count checks.
- `scripts/repo_pack_self_test.py` now owns the repo-pack validator self-test
  fixture harness, bringing `scripts/validate_repo_pack.py` below the fail
  threshold while preserving validation behavior.
- `scripts/repo_pack_contracts.py` now owns shared repo-pack contract constants,
  slice coverage helpers, and reusable predicates used by validator modules.
- `scripts/repo_pack_acceptance.py` now owns acceptance ledger, acceptance
  mapping, scope-closure, and acceptance-coverage checks, lowering
  `scripts/validate_repo_pack.py` to 1,972 lines while preserving validation
  behavior.
- `scripts/repo_pack_model_provider.py` now owns model-provider endpoint grant
  extraction and live-eval provider gate checks, lowering
  `scripts/validate_repo_pack.py` to 1,848 lines while preserving validation
  behavior.
- `scripts/repo_pack_factoryd.py` now owns factoryd runtime, portable path, and
  repo config checks, lowering `scripts/validate_repo_pack.py` to 1,665 lines
  while preserving validation behavior.
- `scripts/repo_pack_ci.py` now owns guide coverage, coverage-policy refs, and
  GitHub CI/control-file checks, lowering `scripts/validate_repo_pack.py` to
  1,582 lines while preserving validation behavior.
- `scripts/repo_pack_safety.py` now owns safety/corpus-ready plan checks and
  model-provider risk classification checks, lowering
  `scripts/validate_repo_pack.py` to 1,478 lines while preserving validation
  behavior.
- `scripts/repo_pack_task_specials.py` now owns recorder split and
  first-session smoke task checks, lowering `scripts/validate_repo_pack.py` to
  1,419 lines while preserving validation behavior.
- `scripts/repo_pack_planning.py` now owns context brief, execution plan,
  validation contract, and shared planning-policy checks, lowering
  `scripts/validate_repo_pack.py` to 879 lines while preserving validation
  behavior.

## Candidate Package Boundaries

- `internal/source`: public source-facing orchestration and shared types.
- `internal/source/config.go`: source config read/write/defaulting and
  repo-local config validation.
- `internal/source/fingerprint.go`: source hashing, docs hashing, generated
  source-directory skips, and surface fingerprint assembly.
- `internal/source/openapi.go`: OpenAPI source parsing, operation metadata
  checks, schema/reference resolution, and parser-facing findings.
- `internal/source/openapi_auth.go`: OpenAPI security-scheme and OAuth scope
  description checks.
- `internal/source/docs.go`: docs source walking, operational-guidance checks,
  and broken local Markdown reference findings.
- `internal/source/docs_test.go`: source-check docs/init/link behavior tests.
- `internal/source/markdownlinks`: local Markdown link target parsing, fence
  detection, and missing-reference target normalization.
- `internal/source/parse`: source parsing and raw input normalization.
- `internal/source/validate`: deterministic source validation.
- `internal/source/evidence`: evidence mapping and proof-honesty fields.
- `internal/source/report.go`: report persistence, status projection, and
  finding classification helpers.
- `internal/source/yaml_helpers.go`: YAML scalar rendering, inline-flow
  parsing, OpenAPI component reference helpers, and JSON pointer escaping for
  source checks.
- `internal/source/source_config_report_test.go`: config parsing, report
  persistence, and finding helper tests.
- `internal/source/source_fixtures_*_test.go`: shared OpenAPI/YAML/docs
  fixtures for source tests.
- `internal/source/source_helpers_test.go`: shared source-test assertion
  helpers.
- `internal/source/source_parameters_test.go`: OpenAPI parameter metadata
  coverage tests.

## Required Promotion

- Source kind: review finding / architecture finding.
- Candidate mission: `systemic-architecture-budget`.
- Required command before implementation: `factoryd ingest --kind review` or
  the equivalent governed Factory planning path.
- Required validation after materialization: `make prepush-full`.
