# Factory Artifacts

- `.factory/artifacts/`: committed, source-safe Factory planning, validation,
  closure, and handoff artifacts. It is not the Lumyn consumer-private runtime
  or identifiable pilot-evidence store.
- `.factory/artifacts/lifecycle-evidence/`: trusted independent evaluator,
  grader, and attestor results; never writable by `task-executor`.
- `.factory/tmp/`: ignored local scratch space.
- `.factory/factoryd.json`: local active safe attended daemon configuration
  copied from `.factory/factoryd.example.json`; gitignored.
- `.factory/factoryd.example.json`: safe repo-local daemon configuration template.
- `.factory/factoryd.autoship.example.json`: explicit full-loop daemon configuration template for protected GitHub execution.
- `.factoryd/`: ignored local daemon state, worktrees, claims, events, and run reports.
- `.factory/artifacts/prd-to-plan/lumyn-migration-mvp/acceptance-ledger.json`:
  active itemized migration-MVP acceptance and product-signal closure source.
- `.factory/artifacts/pilot/lumyn-migration-mvp/public/`: the only planned
  committable pilot namespace; it may contain explicitly public, consented,
  redacted aggregates and evidence hashes only.
- `.factory/artifacts/prd-to-plan/lumyn-mvp/`: historical plan for the
  superseded agent-readiness thesis. Its task/control, PR, exception, and pilot
  artifacts remain immutable records and are not current product scope; the
  README carries only a non-operative dispatch tombstone.

The active PRD-to-plan artifacts are under:

```text
.factory/artifacts/prd-to-plan/lumyn-migration-mvp/
```

This control set contains exactly 62 acceptance items across 12 task packets
(`M0` through `M10`, including `M2.5`). `M2.5` freezes the qualified
design partner, supported-class canary, five-repository/three-organization
cohort, operational provider-status and receipt-acknowledgement channels,
supported consumer OS/architecture and host-isolation backend, consent,
confidentiality, retention, economics, and measurement protocol before
migration planning; its identifiable-evidence collection first
requires a narrow manual privacy/legal preflight under Factory `approval`.
That preflight freezes the minimal signed connection-receipt fields,
provider-authenticated consumer signer binding, provider-signed
acknowledgement/cardinality policy, separate public fields, and the
irreversibility of provider/public disclosure. The approval must match that
preflight's exact canonical scope digest. `M10` runs the governed campaign.

Product workers may read but must not directly mutate the active plan,
acceptance ledger, mapping, validation contract, or canonical closure map.
They emit task-scoped evidence; trusted Factory lifecycle workers update
closure.

When selected by policy, independent `holdout-evaluator`, `trace-grader`, and
`evidence-attestor` gates run after code review and before commit/push.
Shipping verifies their lifecycle-owned artifacts are passing, task-bound, and
schema-valid, task/work/lifecycle-run-bound, linked to the current validation
run, candidate digest, and work-proof marker digest, and independently
authored; the implementation worker cannot write the lifecycle namespace,
self-grade, or self-attest.

Only a non-resolving opaque holdout manifest may be committed: opaque case IDs,
provenance class and license posture, a frozen suite commitment, and encrypted
or HMAC artifact commitments. An independent holdout owner provisions and
freezes `LUMYN_HOLDOUT_ROOT`; resolving provenance, plaintext content digests,
held-out inputs, answer keys, expected patches or labels, and raw traces remain
there and are unavailable to implementation workers.
M1's provision-mode policy creates the suite commitment at independent
evaluation time; M4/M6/M7 evaluate-mode policies resolve and byte-bind that
trusted result. The static plan never fabricates a future commitment.

Private Lumyn runtime and identifiable pilot artifacts live in an explicitly
configured consumer-controlled state root outside the checkout and any public
source repository. That store enforces authorization TTL and deletion on
expiry or revocation at creation, read, process startup, and the next run, and
emits a deletion receipt or orphan report. `lumyn artifacts gc` retries partial
cleanup without extending authority or rewriting historical closure. Committed
Factory artifacts refer to private evidence by opaque identifier and digest,
never by raw content or machine-local path. The repository ignore rules for
legacy runtime roots are defense in depth, not authorization to store private
data in the checkout.
Provider export and public commit are non-recallable; revocation stops future
sharing and deletes only Lumyn-controlled private copies.

Factory worker grants use only Factory's closed `approval`, `credentials`, and
`network` vocabulary. Exact Lumyn product authorizations remain private,
schema-backed artifacts; a Factory approval may cite their validated opaque
bundle but cannot substitute for them. Factory selection and dispatch do not
validate the private product bundle. Lumyn revalidates current packet trust and
exact consumer-signed product authority internally, immediately before every
side effect and retry. Current packet status comes from a signed offline
snapshot or exact `provider_trust_status_read` grant. Repository commands
require the declared fail-closed host-isolation backend, and sponsored
connections count only from valid minimal signed receipts whose consumer
signer is provider-authenticated and whose one-repository unit has a
provider-signed, deduplicated acknowledgement.
