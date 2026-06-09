# Lumyn

Lumyn is a local-first OSS CLI for proving whether API workflows can be completed, verified, bounded, and explained with durable evidence.

The full MVP product contract is [docs/product/prd.md](docs/product/prd.md).

## MVP Scope

The MVP includes:

- workflow recording, contract generation, replay verification, and proof-labeled reports
- live known-path verification, sandbox cleanup, action boundaries, and CI integration
- live agent eval with one provider adapter, repeated runs, cost/duration reporting, and baseline comparison

Post-MVP exclusions include MCP recording, event assertions, hosted dashboards, runtime enforcement, production trace import, and multi-provider eval panels.

## Repository Layout

- `cmd/lumyn/`: CLI entrypoint
- `internal/`: Go implementation packages
- `tests/`: tests
- `schemas/`: versioned schemas
- `examples/`: deterministic fixtures and examples
- `workflows/`: workflow contracts
- `cassettes/`: replay cassettes
- `baselines/`: eval baselines
- `runs/`: Lumyn run outputs
- `.factory/artifacts/`: Factory planning and evidence artifacts
- `.factory/artifacts/pr-lifecycle/`: Factory PR lifecycle evidence for validation, CI, review, ship, merge, and post-merge gates
- `.github/workflows/codeql.yml`: CodeQL Go security scanner risk lane
- `docs/dev/dev_guides.md`: toolchain, CI lanes, 12-level test matrix, scanner, docs parity, output contract, release, and provenance rules
- `docs/architecture/architecture_guides.md`: boundaries, systems-thinking, TDD, ADR, performance, reliability, trust-mode, and fail-closed rules

## Validation

```bash
make lint-fast
make test-fast
make test-contracts
make prepush-full
```

GitHub Actions runs `make prepush-full` through the `validate` workflow on pull
requests and pushes to `main`.
GitHub Actions also runs CodeQL Go analysis through the `codeql` workflow.
Scanner-gated changes require CodeQL status evidence or an approved exception.
GitHub `main` branch protection and the `protect-main-from-direct-push` ruleset
can be audited with `make audit-remote-protection` when GitHub credentials are
available.

## Runtime Pins

- Go `1.26.4`
- Module path `github.com/Clyra-AI/lumyn`
- Standalone binary first, Homebrew next
