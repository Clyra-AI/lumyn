# Lumyn Factory Integration

Factory artifacts for Lumyn live under `.factory/artifacts/`.

- PRD-to-plan artifacts: `.factory/artifacts/prd-to-plan/lumyn-mvp/`
- Task-run evidence: `.factory/artifacts/task-runs/`
- PR lifecycle evidence: `.factory/artifacts/pr-lifecycle/`
- Downstream pilot evidence: `.factory/artifacts/pilot/lumyn-mvp-slice/`

The canonical product input is:

```text
docs/product/prd.md
```

This path is repo-relative so Factory profiles and downstream workers do not depend on machine-local paths.

The current downstream pilot evidence package intentionally marks the full MVP
as blocked, not complete. It closes only the bootstrap/planning baseline and
routes the first missing product slice to `T3-repair-001`.
