#!/usr/bin/env python3
from __future__ import annotations

import json
import re
from pathlib import Path
from typing import Any

from repo_pack_contracts import fail


def load_json(root: Path, path: Path) -> dict[str, Any]:
    if not path.exists():
        fail(f"missing JSON artifact: {path.relative_to(root)}")
    try:
        payload = json.loads(path.read_text(encoding="utf-8"))
    except Exception as exc:
        fail(f"{path.relative_to(root)} is not valid JSON: {exc}")
    if not isinstance(payload, dict):
        fail(f"{path.relative_to(root)} must contain a JSON object")
    return payload


def require_existing(root: Path, relative_path: str) -> None:
    if not (root / relative_path).exists():
        fail(f"missing required repo-pack file: {relative_path}")


def ref_file_exists(root: Path, ref: Any) -> bool:
    if not isinstance(ref, str) or not ref.strip():
        return False
    path_part = ref.split("#", 1)[0]
    return bool(path_part) and (root / path_part).exists()


def validate_coverage_policy_refs(root: Path, value: Any, label: str) -> None:
    if not isinstance(value, dict):
        fail(f"{label} must be an object")
    for key in ["exception_ref", "evidence_ref"]:
        if isinstance(value.get(key), str) and value[key].strip() and not ref_file_exists(root, value[key]):
            fail(f"{label}.{key} points to missing file {value[key]}")
    minimums = value.get("minimums")
    if isinstance(minimums, list):
        for index, item in enumerate(minimums):
            if not isinstance(item, dict):
                continue
            if isinstance(item.get("exception_ref"), str) and item["exception_ref"].strip() and not ref_file_exists(root, item["exception_ref"]):
                fail(f"{label}.minimums[{index}].exception_ref points to missing file {item['exception_ref']}")


def validate_guides(root: Path, required_guides: list[str]) -> None:
    for relative_path in required_guides:
        require_existing(root, relative_path)
    dev_guide = (root / "docs/dev/dev_guides.md").read_text(encoding="utf-8")
    dev_guide_lower = dev_guide.lower()
    tiers = set(re.findall(r"\|\s*Tier\s+(\d+)\b", dev_guide))
    expected = {str(index) for index in range(1, 13)}
    if tiers != expected:
        fail(f"docs/dev/dev_guides.md must preserve all 12 test tiers; found {sorted(tiers)}")
    for token in ["coverage gates", "make test-coverage", ">= 75%"]:
        if token not in dev_guide_lower:
            fail(f"docs/dev/dev_guides.md missing coverage token {token!r}")
    makefile = (root / "Makefile").read_text(encoding="utf-8")
    for token in ["test-coverage:", "check_go_coverage.py", "prepush-full: fmt lint-fast test-fast test-coverage"]:
        if token not in makefile:
            fail(f"Makefile missing coverage gate token {token!r}")
    arch_guide = (root / "docs/architecture/architecture_guides.md").read_text(encoding="utf-8").lower()
    for token in ["systems thinking", "tdd", "adr", "performance", "reliability", "fail-closed", "coverage gates"]:
        if token not in arch_guide:
            fail(f"docs/architecture/architecture_guides.md missing architecture token {token!r}")


def validate_ci_control_set(
    root: Path,
    required_checks_path: Path,
    codeowners_path: Path,
    action_ref_exceptions_path: Path,
    validate_workflow_path: Path,
    codeql_workflow_path: Path,
    required_status_checks: list[str],
    required_action_refs: list[str],
) -> None:
    for path in [
        required_checks_path,
        codeowners_path,
        action_ref_exceptions_path,
        validate_workflow_path,
        codeql_workflow_path,
    ]:
        if not path.exists():
            fail(f"missing CI control file: {path.relative_to(root)}")

    required_checks = load_json(root, required_checks_path).get("required_checks")
    if not isinstance(required_checks, list):
        fail(".github/required-checks.json.required_checks must be a list")
    missing_checks = [check for check in required_status_checks if check not in required_checks]
    if missing_checks:
        fail(f".github/required-checks.json missing required checks: {missing_checks}")

    validate_workflow = validate_workflow_path.read_text(encoding="utf-8")
    validate_tokens = [
        "pull_request:",
        "push:",
        "branches:",
        "- main",
        "permissions:",
        "contents: read",
        "concurrency:",
        "cancel-in-progress: true",
        "timeout-minutes:",
        "actions/checkout@v6.0.2",
        "actions/setup-go@v6.3.0",
        "go-version-file: go.mod",
        "check-latest: false",
        "cache: true",
        "make prepush-full",
    ]
    for token in validate_tokens:
        if token not in validate_workflow:
            fail(f".github/workflows/validate.yml missing CI control token {token!r}")

    codeql_workflow = codeql_workflow_path.read_text(encoding="utf-8")
    codeql_tokens = [
        "pull_request:",
        "push:",
        "branches:",
        "- main",
        "permissions:",
        "security-events: write",
        "contents: read",
        "concurrency:",
        "cancel-in-progress: true",
        "timeout-minutes:",
        "actions/checkout@v6.0.2",
        "actions/setup-go@v6.3.0",
        "go-version-file: go.mod",
        "check-latest: false",
        "github/codeql-action/init@v4",
        "github/codeql-action/autobuild@v4",
        "github/codeql-action/analyze@v4",
        "languages: go",
    ]
    for token in codeql_tokens:
        if token not in codeql_workflow:
            fail(f".github/workflows/codeql.yml missing CI control token {token!r}")

    codeowners = codeowners_path.read_text(encoding="utf-8")
    for token in ["*", "/.github/**", "/schemas/**", "/cmd/**", "/internal/**"]:
        if token not in codeowners:
            fail(f".github/CODEOWNERS missing owner token {token!r}")

    action_exceptions = action_ref_exceptions_path.read_text(encoding="utf-8")
    for token in required_action_refs + ["owner:", "reason:", "scope:", "expires:", "review_command:"]:
        if token not in action_exceptions:
            fail(f".github/action-ref-exceptions.yaml missing action-ref token {token!r}")
