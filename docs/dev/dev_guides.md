# Lumyn Developer Guide

## Validation Matrix

- `make lint-fast`: repo operating pack and layout checks.
- `make test-fast`: Python unit tests.
- `make test-contracts`: unit tests plus Factory planning artifact presence.
- `make prepush-full`: full local gate before PR or merge.

## Bootstrap Rules

- Deterministic bootstrap must not require network, sandbox credentials, or model keys.
- Live sandbox and eval work require explicit human approval before credentials are introduced.
- Tests should be added before implementation when practical.
- Evidence artifacts must use repo-relative paths.

