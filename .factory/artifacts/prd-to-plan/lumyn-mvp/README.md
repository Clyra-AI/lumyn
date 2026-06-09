# Lumyn MVP PRD-To-Plan

This directory contains the initial Factory planning artifacts derived from:

```text
docs/product/prd.md
```

The artifacts use repo-relative paths so they can be consumed from any machine or Factory worker.

Initial implementation should start with task `T1` in `task-packets.json`.

The repo-pack validator enforces both guide propagation and Factory planning-skill alignment. In practice, changes to these artifacts must keep `scripts/validate_repo_pack.py` green so task packets preserve the Lumyn dev guide, architecture guide, `prd-to-plan`, and `execution-compiler` requirements before implementation continues.
