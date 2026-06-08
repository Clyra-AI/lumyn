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

## Runtime Shape

- `cmd/lumyn/`: command entrypoint and process exit behavior.
- `internal/result/`: command-result envelope.
- `internal/exitcode/`: stable exit-code constants matching `docs/product/prd.md`.
- `internal/config/`: config discovery and repo-relative path handling.
- `internal/version/`: version metadata.

Use Go `1.26.4` and standard-library-only implementation for T1.
