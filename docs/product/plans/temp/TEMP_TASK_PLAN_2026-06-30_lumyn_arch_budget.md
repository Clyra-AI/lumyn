# TEMP Task Plan: Lumyn Architecture Budget And Early Source Split

Date: 2026-06-30
Status: queued after Factory/factoryd architecture budget contracts
Repo: Lumyn

## Objective

Adopt architecture budget rules early so Lumyn does not repeat Relia's monolith
pattern as product behavior expands.

## Current Finding

Lumyn already has a thin `cmd/lumyn` entrypoint and several implementation
packages under `internal/`. The immediate risk is not the CLI entrypoint; it is
that `internal/source` has already accumulated broad responsibility and large
tests.

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
3. Split `source_test.go` into tests near the moved package responsibilities.
4. Preserve public CLI behavior, command-result JSON, schemas, and examples.

## Candidate Package Boundaries

- `internal/source`: public source-facing orchestration and shared types.
- `internal/source/parse`: source parsing and raw input normalization.
- `internal/source/validate`: deterministic source validation.
- `internal/source/evidence`: evidence mapping and proof-honesty fields.
- `internal/source/report`: report-facing projections.

## Validation

- `make lint-fast`
- `make test-fast`
- `make test-coverage`
- `make test-contracts`
- targeted package tests for split behavior
- `make prepush-full`

## Acceptance Criteria

- Lumyn keeps the thin CLI entrypoint shape.
- `internal/source` no longer carries every source responsibility in one file.
- Test files are split around package behavior.
- Future tasks have deterministic architecture budget feedback before PR
  review.

