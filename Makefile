.PHONY: lint-fast test-fast test-contracts prepush-full

lint-fast:
	test -f AGENTS.md
	test -f WORKFLOW.md
	test -f README.md
	test -f docs/product/prd.md
	test -d .factory/artifacts
	test -d src/lumyn
	test -d tests
	! grep -RIn "TODO\\|TBD\\|FIXME" AGENTS.md WORKFLOW.md README.md docs src tests pyproject.toml

test-fast:
	python3 -m unittest discover -s tests

test-contracts:
	python3 -m unittest discover -s tests
	test -f .factory/artifacts/prd-to-plan/lumyn-mvp/context-brief.json
	test -f .factory/artifacts/prd-to-plan/lumyn-mvp/execution-plan.json
	test -f .factory/artifacts/prd-to-plan/lumyn-mvp/task-packets.json
	test -f .factory/artifacts/prd-to-plan/lumyn-mvp/validation-contract.json
	test -f .factory/artifacts/prd-to-plan/lumyn-mvp/scope-closure-map.json

prepush-full: lint-fast test-fast test-contracts
