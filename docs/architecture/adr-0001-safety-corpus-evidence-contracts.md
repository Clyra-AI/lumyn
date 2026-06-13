# ADR-0001: Local Safety And Corpus-Ready Evidence Contracts

## Status

Accepted for T2.7 foundation work.

## Context

Lumyn must preserve normalized result and failure evidence before downstream
record, verify, boundary, report, CI, share, and eval surfaces start emitting
artifacts. The contract has to make proof strength, boundary status, safety
relevance, fix target, surface identity, eval mode, and provider/model context
machine-readable without introducing hosted telemetry or corpus upload behavior.

## Decision

`command-result.schema.json`, `result-axes.schema.json`, and
`evidence-event.schema.json` now require local safety/corpus-ready fields:
`finding_kind`, `proof_strength`, `action_boundary_status`,
`security_relevance`, `fix_target`, `surface_fingerprint`, `eval_mode`,
`provider_metadata`, and `corpus_eligible`.

`corpus_eligible` is required and must be `false` for the MVP. Provider/model
metadata is represented as a required object with an `applicable` flag; when it
is not applicable, provider and model are explicitly `not_applicable`.

Boundary violations that involve forbidden endpoints, scope escalation, data
exposure risk, or out-of-policy actions can be represented as
safety/security-relevant findings. This is a workflow completion safety
classification, not a claim that the MVP is a broad security platform.

## Systems Map

State owner: schema-backed local artifacts under `schemas/`, command-result
JSON, canonical evidence events, cassettes, traces, and reports.

Feedback source: `make test-contracts`, schema compilation, representative
valid fixtures, invalid fail-closed fixtures, and CLI schema validation tests.

Blast radius: downstream product tasks must emit the normalized fields whenever
they create result axes, evidence events, or command results.

Rollback path: revert the schema fields, Go envelope fields, and corresponding
tests before any downstream artifact depends on them.

Performance and cost impact: local JSON Schema validation only; no network,
model call, telemetry, upload, or provider endpoint is introduced.

Reliability posture: missing normalized fields and `corpus_eligible: true`
fail schema validation before closure.
