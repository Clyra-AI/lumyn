# Lumyn MVP PRD-To-Plan

This directory contains the initial Factory planning artifacts derived from:

```text
docs/product/prd.md
```

The artifacts use repo-relative paths so they can be consumed from any machine or Factory worker.

Initial implementation should start with task `T1` in `task-packets.json`.

The plan treats the historical `v0.0`, `v0.1`, and `v0.2` labels as required
internal MVP capability slices, not optional releases:

- `v0.0`: record, contract, replay, and report (`T1`, `T2`, `T2.7`, `T3`, `T4.1`, `T4.2`, `T4.3`, `T5`, `T6`, `T10`)
- `v0.1`: live known-path verify, basic boundaries, CI, and share (`T7`, `T8`, `T9`, `T10`)
- `v0.2`: live agent eval (`T11`, `T12`)

Full MVP closure requires all three slices to close or carry explicit approved
delivery debt. The validator enforces this mapping across the execution plan,
task packets, validation contract, acceptance mapping, and scope-closure map.

The repo-pack validator enforces both guide propagation and Factory planning-skill alignment. In practice, changes to these artifacts must keep `scripts/validate_repo_pack.py` green so task packets preserve the Lumyn dev guide, architecture guide, `prd-to-plan`, and `execution-compiler` requirements before implementation continues.

The plan also carries the MVP safety/corpus-ready evidence requirement as
itemized acceptance work, not as a positioning note. `T2.7` establishes the
local normalized failure/result contract before `T3`; `T6`, `T8`, `T9`, `T11`,
and `T12` then surface those fields through reports, boundaries, CI/PR output,
and eval diagnostics. The MVP remains local-first: normalized evidence defaults
to `corpus_eligible: false`, and hosted telemetry, shared failure databases, and
community registries remain post-MVP.

T4 was split into recorder capture/redaction (`T4.1`), workflow/cassette drafting (`T4.2`), and recorder quality measurement (`T4.3`) so the 70 percent REC-QUALITY-001 gate has its own deterministic fixture corpus and closure evidence.

T6 carries first-session smoke evidence for ACT-001, ACT-002, and ACT-003. The smoke report must capture elapsed times and stay offline/local by default.
