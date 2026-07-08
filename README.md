# Lumyn

Lumyn is a local-first OSS CLI for proving whether API workflows can be completed, verified, bounded, and explained with durable evidence.

Lumyn is not a one-shot "agent readiness" generator. Its wedge is continuous
workflow-completion proof: record a workflow, prove the outcome, keep that proof
fresh in local reports and CI, and fail closed when source references,
read-backs, action boundaries, or evidence are not trustworthy.

The full MVP product contract is [docs/product/prd.md](docs/product/prd.md).
Factory planning artifacts map the full MVP through acceptance item IDs plus
generic `delivery_slices` / task `delivery_slice_refs` for the internal
`v0.0`, `v0.1`, and `v0.2` capability slices. The slice refs are coverage
metadata, not separate release boundaries or replacements for item-level scope
closure.

## MVP Scope

The MVP includes:

- workflow recording, contract generation, replay verification, and proof-labeled reports
- live known-path verification, sandbox cleanup, action boundaries, and CI integration
- live agent eval with OpenAI-compatible provider adapters, including custom `base_url` local endpoints, Anthropic provider adapters, repeated runs, cost/duration reporting, and baseline comparison

Post-MVP exclusions include MCP recording, event assertions, hosted dashboards, runtime enforcement, production trace import, and multi-provider eval panels. OpenAI-compatible and Anthropic adapters are MVP scope; comparative multi-provider panels are not. Local open-source model servers are supported through the OpenAI-compatible `base_url` path; Lumyn does not bundle model weights or local inference-runtime payloads.
Live eval provider endpoints, including custom `base_url` endpoints, stay
blocked until the repo records a complete `model_provider_endpoint` grant with
provider identity, provider model, endpoint or `base_url`, credential
environment, budget posture, redaction posture, and allowlist; generic network
or credential approval is not enough.

## Repository Layout

- `cmd/lumyn/`: CLI entrypoint
- `internal/`: Go implementation packages
- `tests/`: tests
- `schemas/`: versioned schemas
- `examples/`: deterministic fixtures and examples
- `workflows/`: workflow contracts
- `cassettes/`: replay cassettes
- `baselines/`: eval baselines
- `runs/`: Lumyn run outputs
- `.factory/artifacts/`: Factory planning and evidence artifacts
- `.factory/artifacts/pr-lifecycle/`: Factory PR lifecycle evidence for validation, CI, review, ship, merge, and post-merge gates
- `.factory/artifacts/pilot/lumyn-mvp-slice/`: downstream Factory pilot evidence, current PRD closure status, and bounded T3 repair task
- `.github/required-checks.json`: required status-check contract for branch protection and `commit-push`
- `.github/CODEOWNERS`: owner-review coverage for workflow, policy, Factory artifact, schema, CLI, and core implementation paths
- `.github/action-ref-exceptions.yaml`: audited exception records for non-SHA GitHub Action refs
- `.github/workflows/codeql.yml`: CodeQL Go security scanner risk lane
- `docs/dev/dev_guides.md`: toolchain, CI lanes, 12-level test matrix, scanner, docs parity, output contract, release, and provenance rules
- `docs/architecture/architecture_guides.md`: boundaries, systems-thinking, TDD, ADR, performance, reliability, trust-mode, and fail-closed rules

## Validation

```bash
make lint-fast
make test-fast
make test-coverage
make test-contracts
make prepush-full
```

`make test-contracts` also runs
`python3 scripts/validate_factory_pilot_evidence.py` to ensure the downstream
Factory pilot proof chain remains coherent, and
`python3 scripts/validate_repo_pack.py` to ensure task packets keep the dev and
architecture guide requirements plus Factory compatibility, alignment, drift,
scope, worker-chain, and lifecycle requirements propagated into executable
planning artifacts.

GitHub Actions runs `make prepush-full` through the `validate` workflow on pull
requests and pushes to `main`.
GitHub Actions also runs CodeQL Go analysis through the `codeql` workflow.
The required check manifest lists `validate` and `CodeQL analyze`, and the
workflow files declare least-privilege permissions, concurrency cancellation,
job timeouts, and toolchain setup from `go.mod`.
Coverage-gated changes require `make test-coverage` evidence or an approved
scoped exception.
Scanner-gated changes require CodeQL status evidence or an approved exception.
GitHub `main` branch protection and the `protect-main-from-direct-push` ruleset
can be audited with `make audit-remote-protection` when GitHub credentials are
available.

## Source Checks And Reports

`lumyn check` must separate infrastructure-tier source facts from
workflow/domain-tier proof claims. Missing, unprovable, or hallucinated source
paths fail source-evidence closure instead of becoming weak findings. Recorder
and report tasks should carry a grounded workflow insight or business-job field
so reports explain why the workflow matters, not only which endpoints appeared.

## Runtime Pins

- Go `1.26.4`
- Module path `github.com/Clyra-AI/lumyn`
- Standalone binary first, Homebrew next
- MVP eval providers: OpenAI-compatible HTTP and Anthropic Messages HTTP adapters
- Live sandbox, model-key, and network work remains blocked until deterministic replay passes and human approval unlocks credential/network posture

## Factory And factoryd Operation

Factory is the contract source for Lumyn's planning artifacts, task packets,
validation contracts, repo-pack requirements, and worker-chain expectations.
`factoryd` is the executable runtime that may consume those artifacts to select
tasks, run Codex or another bounded worker, validate evidence, ship PRs, update
scope closure, and generate repair tasks. Lumyn remains the product source of
truth for the PRD, code, CI, branch policy, and product evidence.

Safe operator path:

```bash
export FACTORY_REPO=../factory
factoryd doctor --config .factory/factoryd.json --repo lumyn --json
factoryd run --config .factory/factoryd.json --repo lumyn --dry-run --json
```

Use the autoship config only after branch protection, required `validate` and
`CodeQL analyze` checks, passive Codex review settle, merge policy, post-merge
monitoring, itemized acceptance-ledger coverage, and semantic scope-closure
evidence are all verified:

```bash
export FACTORY_REPO=../factory
factoryd run --config .factory/factoryd.autoship.example.json --repo lumyn --loop --max-tasks 1 --json
```

Post-PRD audit or review findings become governed work through `factoryd ingest`:

```bash
export FACTORY_REPO=../factory
factoryd ingest --config .factory/factoryd.example.json --repo lumyn --kind audit --input product/audits/<mission>.finding-list.json --mission <mission> --json
factoryd ingest --config .factory/factoryd.example.json --repo lumyn --kind review --input product/reviews/<mission>.finding-list.json --mission <mission> --json
```

Those commands create `.factory/artifacts/post-prd/<mission>/`; that directory
becomes the execution contract for the follow-up mission.
Use Factory `task-supervisor` for guided audit, review, recommendation, or idea
intake that records `task_supervisor_report` evidence before selecting an
autoship task.
