# Lumyn MVP PRD-To-Plan

This directory contains the initial Factory planning artifacts derived from:

```text
docs/product/prd.md
```

The artifacts use repo-relative paths so they can be consumed from any machine or Factory worker.

Initial implementation should start with task `T1` in `task-packets.json`.

The plan treats the historical `v0.0`, `v0.1`, and `v0.2` labels as required
internal MVP capability slices, not optional releases:

- `v0.0`: record, contract, replay, and report (`T1`, `T2`, `T3`, `T4`, `T5`, `T6`, `T10`)
- `v0.1`: live known-path verify, basic boundaries, CI, and share (`T7`, `T8`, `T9`, `T10`)
- `v0.2`: live agent eval (`T11`, `T12`)

Full MVP closure requires all three slices to close or carry explicit approved
delivery debt. The validator enforces this mapping across the execution plan,
task packets, validation contract, acceptance mapping, and scope-closure map.

The repo-pack validator enforces both guide propagation and Factory planning-skill alignment. In practice, changes to these artifacts must keep `scripts/validate_repo_pack.py` green so task packets preserve the Lumyn dev guide, architecture guide, `prd-to-plan`, and `execution-compiler` requirements before implementation continues.
