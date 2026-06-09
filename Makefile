GO ?= go
PKG_LIST := ./cmd/... ./internal/... ./schemas/...

.PHONY: fmt lint-fast test-fast test-contracts build prepush-full

lint-fast:
	test -f AGENTS.md
	test -f WORKFLOW.md
	test -f README.md
	test -f docs/product/prd.md
	test -f .tool-versions
	test -f go.mod
	test -f .github/workflows/validate.yml
	test -f .github/workflows/codeql.yml
	test -d .factory/artifacts
	test -d cmd/lumyn
	test -d internal
	test -d tests
	! grep -RIn "TODO\\|TBD\\|FIXME" AGENTS.md WORKFLOW.md README.md docs cmd internal schemas tests
	grep -q '^golang 1.26.4$$' .tool-versions
	grep -q '^go 1.26.4$$' go.mod
	grep -q 'go-version: "1.26.4"' .github/workflows/validate.yml
	grep -q 'make prepush-full' .github/workflows/validate.yml
	grep -q 'github/codeql-action/init@v3' .github/workflows/codeql.yml
	grep -q 'languages: go' .github/workflows/codeql.yml
	$(GO) vet $(PKG_LIST)

fmt:
	gofmt -w $$(find cmd internal tests -name '*.go' -type f)

test-fast:
	$(GO) test ./... -count=1

test-contracts:
	$(GO) test ./... -count=1
	test -f .factory/artifacts/prd-to-plan/lumyn-mvp/context-brief.json
	test -f .factory/artifacts/prd-to-plan/lumyn-mvp/execution-plan.json
	test -f .factory/artifacts/prd-to-plan/lumyn-mvp/task-packets.json
	test -f .factory/artifacts/prd-to-plan/lumyn-mvp/validation-contract.json
	test -f .factory/artifacts/prd-to-plan/lumyn-mvp/scope-closure-map.json
	test -f .factory/artifacts/task-runs/T2.5/validation-report.json
	test -f .factory/artifacts/task-runs/T2.5/work-proof-marker.json
	test -f .factory/artifacts/task-runs/T2.6/validation-report.json
	test -f .factory/artifacts/task-runs/T2.6/work-proof-marker.json
	test -f schemas/workflow-contract.schema.json
	test -f schemas/expected-outcome.schema.json
	test -f schemas/validator.schema.json
	test -f schemas/action-boundary.schema.json
	test -f schemas/human-annotation.schema.json
	test -f schemas/required-context.schema.json
	test -f schemas/state-binding.schema.json
	test -f schemas/canonical-trace.schema.json
	test -f schemas/evidence-event.schema.json
	test -f schemas/cassette.schema.json
	test -f schemas/result-axes.schema.json
	test -f schemas/proof-strength.schema.json
	test -f schemas/command-result.schema.json
	test -f schemas/redaction-config.schema.json

build:
	$(GO) build -o .factory/tmp/lumyn ./cmd/lumyn

prepush-full: fmt lint-fast test-fast test-contracts build
