# AGENTS.md - Lumyn Repository Contract

Version: 0.1
Status: Normative
Scope: This repository only.

## 1. Scope And Intent

- Build Lumyn OSS as a local, repo-native CLI for workflow recording, verification, evidence, reporting, CI integration, and live agent eval.
- Treat [docs/product/prd.md](docs/product/prd.md) as the product source of truth for the full MVP.
- Keep Factory run artifacts under `.factory/artifacts/`.
- Keep transient local runtime material under `.factory/tmp/`.

## 2. North Star

Every change should improve one or more of:

- workflow contract executability
- deterministic replay verification
- live known-path verification
- proof-honest evidence and reports
- redaction before persistence or sharing
- action-boundary enforcement
- CI adoption
- live eval honesty and traceability

## 3. Non-Negotiable Constraints

- Do not claim strong proof unless a reliable read-back confirms the business state.
- Do not persist unredacted secrets in traces, cassettes, reports, logs, or shared artifacts.
- Do not use live network or model calls in deterministic bootstrap tests.
- Do not require a hosted Lumyn account for MVP workflows.
- Do not treat stochastic eval as a default hard CI gate.
- Do not broaden MVP scope into MCP, event assertions, hosted dashboards, runtime enforcement, production trace import, or multi-provider panels.

## 4. Required Boundaries

- `docs/product/`: product requirements and scope closure source.
- `docs/dev/`: repo-local developer guidance.
- `docs/architecture/`: repo-local architecture guidance.
- `.factory/artifacts/`: durable Factory planning, validation, closure, and handoff artifacts.
- `.factory/tmp/`: ignored local execution scratch space.
- `schemas/`: versioned schemas.
- `src/lumyn/`: implementation.
- `tests/`: automated tests.
- `examples/`: deterministic examples and fixtures.
- `workflows/`, `cassettes/`, `baselines/`, `runs/`: Lumyn product artifact surfaces.

## 5. Required Validation

For normal changes, run:

- `make lint-fast`
- `make test-fast`
- `make test-contracts`

Before PR or merge, run:

- `make prepush-full`

If any command is skipped, the validation report must record why.

