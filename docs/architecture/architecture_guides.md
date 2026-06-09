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
- `schemas/`: executable JSON Schema contracts for workflow, evidence, cassette, proof, result, redaction, and related artifact models.
- `.github/workflows/validate.yml`: repository delivery gate that runs `make prepush-full` for PR and `main` validation.
- `.factory/artifacts/pr-lifecycle/`: Factory delivery evidence tying PR validation, CI/status checks, review, shipping, merge, and post-merge monitoring together.

Use Go `1.26.4`. T1 stays standard-library-only. T2 introduces the pinned
`github.com/santhosh-tekuri/jsonschema/v5` validator for executable schema
tests.

The T2.5 lifecycle baseline is delivery infrastructure for Factory-governed
work. It does not replace the later product-facing T9 work, which remains
responsible for Lumyn's own GitHub Action/JUnit/PR-comment behavior as an MVP
feature.
