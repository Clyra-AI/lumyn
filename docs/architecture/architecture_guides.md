# Lumyn Architecture Guide

## Initial Boundaries

- CLI command surface
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

