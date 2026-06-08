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

## Validation

```bash
make lint-fast
make test-fast
make test-contracts
make prepush-full
```

## Runtime Pins

- Go `1.26.4`
- Module path `github.com/Clyra-AI/lumyn`
- Standalone binary first, Homebrew next
