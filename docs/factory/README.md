# Lumyn Factory Integration

Factory artifacts for Lumyn live under `.factory/artifacts/`.

- PRD-to-plan artifacts: `.factory/artifacts/prd-to-plan/lumyn-mvp/`
- Task-run evidence: `.factory/artifacts/task-runs/`
- PR lifecycle evidence: `.factory/artifacts/pr-lifecycle/`
- Downstream pilot evidence: `.factory/artifacts/pilot/lumyn-mvp-slice/`
- Daemon config template: `.factory/factoryd.example.json`
- Local daemon runtime state: `.factoryd/` (gitignored)

The canonical product input is:

```text
docs/product/prd.md
```

This path is repo-relative so Factory profiles and downstream workers do not depend on machine-local paths.

The current downstream pilot evidence package intentionally marks the full MVP
as blocked, not complete. It closes only the bootstrap/planning baseline and
routes the first missing product slice to `T3-repair-001`.

Repo-pack guide propagation is enforced locally by:

```text
python3 scripts/validate_repo_pack.py
```

That check verifies the Lumyn dev guide, architecture guide, execution plan,
task packets, and validation contract keep CI lanes, the 12-level test matrix,
scanner posture, engineering policy refs, and architecture guidance refs
connected before T3+ product implementation continues.

`factoryd` runtime state is operational state, not Lumyn product evidence.
Durable evidence remains under `.factory/artifacts/`; claims, worktrees, daemon
events, and single-task run reports remain under `.factoryd/` unless promoted
into a committed Factory artifact.

Codex CLI authentication used by a Factory daemon worker is not Lumyn product
credential access. Lumyn task packets still default to no ambient secrets and
offline product/runtime network posture until a specific live sandbox or model
provider task is approved.

Approved live approval, credential, or network work must be represented in
`.factory/factoryd.example.json` through task-scoped `capability_grants`.
Do not edit PRD-derived task packets just to bypass a daemon gate.

Autonomous shipping is disabled in `.factory/factoryd.example.json`. Enabling
it requires explicit command hooks for remote lifecycle phases such as push, PR,
CI/status wait, passive Codex review settle, merge, post-merge monitoring, and
scope-closure mutation. Missing required hooks must block rather than being
treated as successful delivery.
