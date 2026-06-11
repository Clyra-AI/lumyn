# Lumyn Factory Integration

Factory artifacts for Lumyn live under `.factory/artifacts/`.

- PRD-to-plan artifacts: `.factory/artifacts/prd-to-plan/lumyn-mvp/`
- Task-run evidence: `.factory/artifacts/task-runs/`
- PR lifecycle evidence: `.factory/artifacts/pr-lifecycle/`
- Downstream pilot evidence: `.factory/artifacts/pilot/lumyn-mvp-slice/`
- Active safe attended daemon config: `.factory/factoryd.json`
- Safe daemon config template: `.factory/factoryd.example.json`
- Explicit autoship daemon config template: `.factory/factoryd.autoship.example.json`
- Local daemon runtime state: `.factoryd/` (gitignored)

The canonical product input is:

```text
docs/product/prd.md
```

This path is repo-relative so Factory profiles and downstream workers do not depend on machine-local paths.

## Operator Flow

Use Factory to change shared contracts and planning rules. Use `factoryd` to
execute Lumyn task packets. Use Lumyn commits and PRs to change product code,
CI, docs, and product evidence.

Start with non-mutating proof:

```text
export FACTORY_REPO=../factory
factoryd doctor --config .factory/factoryd.json --repo lumyn --json
factoryd run --config .factory/factoryd.json --repo lumyn --dry-run --json
```

Run one implementation task without remote shipping only after the selected
task packet has allowed paths, forbidden paths, validation commands, evidence
requirements, lifecycle gates, and stop conditions:

```text
export FACTORY_REPO=../factory
factoryd run --config .factory/factoryd.json --repo lumyn --once --json
```

Use autoship only after branch protection, required `validate` and
`CodeQL analyze` checks, passive Codex review settle, merge policy, post-merge
monitoring, and semantic scope-closure evidence are all verified:

```text
export FACTORY_REPO=../factory
factoryd run --config .factory/factoryd.autoship.example.json --repo lumyn --loop --max-tasks 1 --json
```

The one-task loop is intentional. It keeps PRs small, allows CI and passive
review to settle per task, and produces task-scoped work-proof, validation,
PR lifecycle, post-merge, scope-closure, and repair evidence.

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

The active daemon configs carry a restricted `runtime_control` block. `factoryd`
must honor that block before worker dispatch, including pause/cancel/freeze,
disabled-adapter, read-only, and write-scope gates. Claimed task runs also emit
a canonical `.factoryd/runs/<run>/mission-event-log.json` and rerun repo-level
validation commands in the worker worktree after task changes are present.

PRD-derived control truth under `.factory/artifacts/prd-to-plan/lumyn-mvp/` is
trusted runtime/planning state. Product workers must not edit the context brief,
execution plan, task packets, validation contract, acceptance mapping, or
scope-closure map directly. They emit task-scoped evidence; `factoryd` updates
scope closure only through the configured semantic shipping phase.

Codex CLI authentication used by a Factory daemon worker is not Lumyn product
credential access. Lumyn task packets still default to no ambient secrets and
offline product/runtime network posture until a specific live sandbox or model
provider task is approved.

Approved live approval, credential, or network work must be represented in
the active `.factory/factoryd*.json` config through task-scoped
`capability_grants`.
Do not edit PRD-derived task packets just to bypass a daemon gate.

Autonomous shipping is disabled in `.factory/factoryd.example.json`. The
explicit `.factory/factoryd.autoship.example.json` template may be used only
after branch protection, `validate`, CodeQL, passive Codex review, merge,
post-merge, and semantic scope-closure gates are proven. The autoship template
uses the `github_cli` provider for remote lifecycle phases and must block rather
than treating missing or failed required phases as successful delivery.

The expected one-task proof command is:

```text
FACTORY_REPO=/path/to/factory factoryd run --config .factory/factoryd.autoship.example.json --repo lumyn --loop --max-tasks 1 --json
```

## Post-PRD Audit Or Review Findings

Material findings from `app-audit` and `code-review` become governed work only
after they are saved as a repo-local source list and ingested:

```text
FACTORY_REPO=/path/to/factory factoryd ingest --config .factory/factoryd.example.json --repo lumyn --kind audit --input product/audits/<mission>.md --mission <mission> --json
FACTORY_REPO=/path/to/factory factoryd ingest --config .factory/factoryd.example.json --repo lumyn --kind review --input product/reviews/<mission>.md --mission <mission> --json
```

The generated `.factory/artifacts/post-prd/<mission>/` artifacts become the
mission contract for execution. Do not mutate `docs/product/prd.md` unless a
human explicitly promotes the finding into product scope.
