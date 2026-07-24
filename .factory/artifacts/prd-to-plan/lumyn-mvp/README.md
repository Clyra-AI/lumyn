# Lumyn MVP PRD-To-Plan

> **Superseded historical plan — do not dispatch.** This directory is retained
> only as immutable planning history. The active plan is
> `.factory/artifacts/prd-to-plan/lumyn-migration-mvp/`; its first selectable
> task is `M0`. Every instruction below describes the former product direction
> and is non-operative.

This directory contains the initial Factory planning artifacts derived from:

```text
docs/product/prd.md
```

The artifacts use repo-relative paths so they can be consumed from any machine or Factory worker.

The historical implementation sequence began with task `T1` in
`task-packets.json`; that instruction must not be dispatched now.

The plan treats the historical `v0.0`, `v0.1`, and `v0.2` labels as required
internal MVP capability slices, not optional releases:

- `v0.0`: record, contract, replay, and report (`T1`, `T2`, `T2.7`, `T3`, `T4.1`, `T4.2`, `T4.3`, `T5.1`, `T5.2`, `T6.1`, `T6.2`, `T10`)
- `v0.1`: live known-path verify, basic boundaries, CI, and share (`T7`, `T8`, `T9`, `T10`)
- `v0.2`: live agent eval (`T11.1`, `T11.2`, `T12.1`, `T12.2`)

Full MVP closure requires all three slices to close or carry explicit approved
delivery debt. The validator enforces this mapping across the execution plan,
task packets, validation contract, acceptance mapping, and scope-closure map.

The repo-pack validator enforces both guide propagation and Factory planning-skill alignment. In practice, changes to these artifacts must keep `scripts/validate_repo_pack.py` green so task packets preserve the Lumyn dev guide, architecture guide, `prd-to-plan`, and `execution-compiler` requirements before implementation continues.

The plan also carries the MVP safety/corpus-ready evidence requirement as
itemized acceptance work, not as a positioning note. `T2.7` establishes the
local normalized failure/result contract before `T3`; `T6.1`, `T6.2`, `T8`,
`T9`, `T11.1`, `T11.2`, `T12.1`, and `T12.2` then surface those fields through
reports, boundaries, CI/PR output, and eval diagnostics. The MVP remains
local-first: normalized evidence defaults to `corpus_eligible: false`, and
hosted telemetry, shared failure databases, and community registries remain
post-MVP.

T4 was split into recorder capture/redaction (`T4.1`), workflow/cassette drafting (`T4.2`), and recorder quality measurement (`T4.3`) so the 70 percent REC-QUALITY-001 gate has its own deterministic fixture corpus and closure evidence.

T6 was split into local report/demo rendering (`T6.1`) and first-session smoke
evidence (`T6.2`). The smoke report must capture ACT-001, ACT-002, and ACT-003
elapsed times and stay offline/local by default.
