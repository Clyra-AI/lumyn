# Lumyn Workflow Contract

Version: 0.1
Status: Normative

## Work Signal

Lumyn accepts work from:

- [docs/product/prd.md](docs/product/prd.md)
- Factory task packets under `.factory/artifacts/`
- GitHub issues or pull requests once the repo is active

## Normal Factory Chain

1. `scout-context`
2. `execution-compiler`
3. `task-executor`
4. `validation-gate`
5. `code-review`
6. `ship-pr`
7. `post-merge-monitor`

## Approval Gates

- Plan approval is required before implementation work starts.
- Merge approval is required before protected branch updates.
- Human approval is required before any live sandbox credential, model key, or non-deterministic eval is used.
- Human approval is required before changing the pinned runtime, module path, or primary distribution target.

## Artifact Rules

- Durable Factory artifacts go under `.factory/artifacts/`.
- The canonical MVP PRD path is `docs/product/prd.md`.
- Factory references must use repo-relative paths, never machine-local absolute paths.
- Validation evidence must include command, status, artifact refs, and skipped-command reasons when applicable.

## Bootstrap Validation Lanes

- Fast lane: `make lint-fast`, `make test-fast`
- Contract lane: `make test-contracts`
- Full lane: `make prepush-full`

## Runtime And Distribution Pins

- Go version: `1.26.4`
- Module path: `github.com/Clyra-AI/lumyn`
- Primary distribution: standalone binary
- Secondary distribution: Homebrew
- Non-primary distribution: PyPI
- First eval provider: OpenAI-compatible HTTP adapter

## Stop Conditions

Stop and request human decision if:

- a task requires live credentials before credential posture is approved
- a task requires network access during deterministic bootstrap
- a proposed change writes outside declared allowed paths
- redaction confidence is unknown for persisted artifacts
- implementation would satisfy a command while violating explicit PRD scope exclusions
