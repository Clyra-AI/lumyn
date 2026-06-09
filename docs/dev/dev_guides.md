# Lumyn Developer Guide

## Toolchain Pins

| Tool | Version |
|---|---:|
| Go | `1.26.4` |

## Dependency Pins

- `github.com/santhosh-tekuri/jsonschema/v5 v5.3.1`: executable JSON Schema validation for T2 and later schema/artifact work.

## Validation Matrix

- `make lint-fast`: repo operating pack and layout checks.
- `make test-fast`: Go unit tests.
- `make test-contracts`: unit tests, Factory planning artifact presence, and required schema-file presence.
- `make prepush-full`: full local gate before PR or merge.

## CI And PR Lifecycle

- GitHub Actions workflow: `.github/workflows/validate.yml`.
- Required check name: `validate`.
- CI command: `make prepush-full`.
- PR lifecycle report path: `.factory/artifacts/pr-lifecycle/<task_id>/pr-lifecycle-report.json`.
- Lifecycle-gated tasks require local validation, CI/status evidence, review evidence when required, ship evidence, post-merge evidence, and a PR lifecycle report or an explicit approved exception.

## Bootstrap Rules

- Deterministic bootstrap must not require network, sandbox credentials, or model keys.
- Live sandbox and eval work require explicit human approval before credentials are introduced.
- Tests should be added before implementation when practical.
- Evidence artifacts must use repo-relative paths.
- T1 must use the Go standard library only.
- Any new dependency must be pinned in `go.mod`, justified in the task evidence, and covered by validation.
- Schema/artifact changes must include representative validation coverage in `schemas/`.
- Changes to CI, review, shipping, or post-merge workflow must update `WORKFLOW.md`, this guide, and the relevant Factory planning artifacts in the same branch.

## Distribution Pins

- Primary: standalone binary.
- Secondary: Homebrew.
- Not primary: PyPI.

## Provider Pins

- First eval provider: OpenAI-compatible HTTP adapter.
- Provider config shape: `provider`, `model`, `temperature`, `base_url`, `api_key_env`.
- Eval provider work is blocked until deterministic replay foundation passes and model-key/network posture is approved.
