# Lumyn Migration MVP Factory Plan

This is the active Factory planning generation for the Lumyn
provider-sponsored, customer-controlled migration MVP.

Authored sources:

- `docs/product/prd.md`
- `docs/product/plan.md`
- `docs/architecture/adr-0002-provider-sponsored-customer-controlled-migrations.md`

Control files:

- `context-brief.json`
- `risk-classification.json`
- `execution-plan.json`
- `task-packets.json`
- `validation-contract.json`
- `acceptance-ledger.json`
- `acceptance-mapping.json`
- `scope-closure-map.json`

These files are one generation and must be updated atomically when product
scope, authority, runtime pins, task dependencies, validation, acceptance, or
Factory compatibility changes. Product workers may read them but may not update
canonical closure directly.

The previous `.factory/artifacts/prd-to-plan/lumyn-mvp/` task and control
artifacts are immutable historical evidence for the superseded agent-readiness
plan. Its README contains only a non-operative dispatch tombstone. It is not an
active task source. Carry-forward references in this generation prove only the
exact CLI, schema, source-intake, planning, and delivery foundations they name.

Current closure: 5 of 62 acceptance
items have direct carry-forward or rebaseline evidence. All remaining items are
planned and must close through item-level evidence.
