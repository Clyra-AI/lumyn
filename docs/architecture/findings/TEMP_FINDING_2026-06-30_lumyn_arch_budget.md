# TEMP Finding: Lumyn Architecture Budget And Early Source Split

Date: 2026-06-30
Status: source finding; not dispatchable
Repo: Lumyn

## Boundary

This file is repo-local source evidence for future Factory/factoryd planning.
It is not a generated execution contract, task packet, acceptance ledger, or
scope-closure artifact. Before implementation starts, this finding must be
ingested or promoted through the governed Factory path so runner-ready task
packets, validation commands, lifecycle evidence, and acceptance refs are
materialized.

## Objective

Adopt architecture budget rules early so Lumyn does not repeat Relia's monolith
pattern as product behavior expands.

## Current Finding

Lumyn already has a thin `cmd/lumyn` entrypoint and several implementation
packages under `internal/`. The immediate risk is not the CLI entrypoint. The
first `internal/source` decomposition pass has brought source package files
below the warning threshold. The repo-pack validator has been split below the
fail threshold, with `scripts/validate_repo_pack.py` still tracked as
warning-level orchestration debt. Source tests and validator self-tests now live
in focused files so future behavior can be split without rebuilding a single
test monolith.

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
