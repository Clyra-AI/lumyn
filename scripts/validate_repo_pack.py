#!/usr/bin/env python3
from __future__ import annotations

import json
import re
import sys
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parents[1]
PLAN_DIR = ROOT / ".factory" / "artifacts" / "prd-to-plan" / "lumyn-mvp"
EXECUTION_PLAN = PLAN_DIR / "execution-plan.json"
TASK_PACKETS = PLAN_DIR / "task-packets.json"
VALIDATION_CONTRACT = PLAN_DIR / "validation-contract.json"

REQUIRED_GUIDES = [
    "docs/dev/dev_guides.md",
    "docs/architecture/architecture_guides.md",
]

REQUIRED_TASK_FIELDS = [
    "ci_lane_refs",
    "test_matrix_refs",
    "security_scanner_gates",
    "engineering_policy_refs",
    "architecture_guidance_refs",
]

STOP_CONDITION_CATEGORIES = {
    "test_matrix": ["test-matrix", "test matrix", "test_matrix"],
    "ci_lanes": ["ci lane", "ci/status", "status check"],
    "scanner": ["scanner", "security"],
    "engineering_policies": ["docs parity", "output contract", "release integrity", "provenance"],
    "architecture_policies": ["architecture", "systems-thinking", "systems thinking", "adr", "fail-closed"],
}


def fail(message: str) -> None:
    raise AssertionError(message)


def load_json(path: Path) -> dict[str, Any]:
    if not path.exists():
        fail(f"missing JSON artifact: {path.relative_to(ROOT)}")
    try:
        payload = json.loads(path.read_text())
    except Exception as exc:
        fail(f"{path.relative_to(ROOT)} is not valid JSON: {exc}")
    if not isinstance(payload, dict):
        fail(f"{path.relative_to(ROOT)} must contain a JSON object")
    return payload


def require_existing(relative_path: str) -> None:
    if not (ROOT / relative_path).exists():
        fail(f"missing required repo-pack file: {relative_path}")


def validate_guides() -> None:
    for relative_path in REQUIRED_GUIDES:
        require_existing(relative_path)
    dev_guide = (ROOT / "docs/dev/dev_guides.md").read_text()
    tiers = set(re.findall(r"\|\s*Tier\s+(\d+)\b", dev_guide))
    expected = {str(index) for index in range(1, 13)}
    if tiers != expected:
        fail(f"docs/dev/dev_guides.md must preserve all 12 test tiers; found {sorted(tiers)}")
    arch_guide = (ROOT / "docs/architecture/architecture_guides.md").read_text().lower()
    for token in ["systems thinking", "tdd", "adr", "performance", "reliability", "fail-closed"]:
        if token not in arch_guide:
            fail(f"docs/architecture/architecture_guides.md missing architecture token {token!r}")


def task_id(task: dict[str, Any]) -> str:
    value = task.get("task_id")
    return value if isinstance(value, str) else ""


def depends_on(task_id_value: str, baseline_task_id: str, tasks_by_id: dict[str, dict[str, Any]], seen: set[str] | None = None) -> bool:
    if task_id_value == baseline_task_id:
        return True
    seen = seen or set()
    if task_id_value in seen:
        return False
    seen.add(task_id_value)
    task = tasks_by_id.get(task_id_value)
    if not task:
        return False
    blocked_by = task.get("blocked_by", [])
    if not isinstance(blocked_by, list):
        fail(f"{task_id_value}.blocked_by must be a list")
    for dependency in [str(value) for value in blocked_by]:
        if dependency == baseline_task_id or depends_on(dependency, baseline_task_id, tasks_by_id, seen):
            return True
    return False


def field_has_evidence(task: dict[str, Any], field: str) -> bool:
    value = task.get(field)
    if field == "security_scanner_gates":
        if not isinstance(value, dict):
            return False
        if isinstance(value.get("exception_ref"), str) and value["exception_ref"].strip():
            return True
        scanner = value.get("scanner")
        if not isinstance(scanner, str) or not scanner.strip():
            return False
        if value.get("required") is False:
            return True
        return any(
            isinstance(value.get(key), str) and value[key].strip()
            for key in ["workflow_ref", "status_check", "evidence_ref"]
        )
    if not isinstance(value, list) or not value:
        return False
    if field == "test_matrix_refs":
        return all(
            isinstance(item, dict)
            and isinstance(item.get("tier"), str)
            and item["tier"].strip()
            and isinstance(item.get("source_ref"), str)
            and item["source_ref"].strip()
            for item in value
        )
    if field == "ci_lane_refs":
        return all(
            isinstance(item, dict)
            and isinstance(item.get("lane"), str)
            and item["lane"].strip()
            and isinstance(item.get("source_ref"), str)
            and item["source_ref"].strip()
            and (
                bool(item.get("command_refs"))
                or bool(item.get("status_check_refs"))
                or (isinstance(item.get("exception_ref"), str) and item["exception_ref"].strip())
            )
            for item in value
        )
    if field == "engineering_policy_refs":
        return all(
            isinstance(item, dict)
            and isinstance(item.get("policy"), str)
            and item["policy"].strip()
            and isinstance(item.get("source_ref"), str)
            and item["source_ref"].strip()
            and (
                isinstance(item.get("rule"), str)
                and item["rule"].strip()
                or isinstance(item.get("exception_ref"), str)
                and item["exception_ref"].strip()
            )
            for item in value
        )
    if field == "architecture_guidance_refs":
        return all(
            isinstance(item, dict)
            and isinstance(item.get("source_ref"), str)
            and item["source_ref"].strip()
            and isinstance(item.get("rule"), str)
            and item["rule"].strip()
            for item in value
        )
    return False


def validate_execution_plan(plan: dict[str, Any]) -> str:
    propagation = plan.get("dev_architecture_propagation")
    if not isinstance(propagation, dict):
        fail("execution plan missing dev_architecture_propagation")
    if propagation.get("status") != "proven":
        fail("dev_architecture_propagation.status must be proven before product implementation")
    baseline_task_id = propagation.get("baseline_task_ref")
    if not isinstance(baseline_task_id, str) or not baseline_task_id:
        fail("dev_architecture_propagation.baseline_task_ref is required")
    if propagation.get("test_matrix_source_ref") != "docs/dev/dev_guides.md#12-level-test-matrix":
        fail("dev_architecture_propagation.test_matrix_source_ref must point at docs/dev/dev_guides.md")
    if propagation.get("architecture_guide_ref") != "docs/architecture/architecture_guides.md":
        fail("dev_architecture_propagation.architecture_guide_ref must point at docs/architecture/architecture_guides.md")
    requirements = propagation.get("task_packet_requirements")
    if not isinstance(requirements, list):
        fail("dev_architecture_propagation.task_packet_requirements must be a list")
    missing = [field for field in REQUIRED_TASK_FIELDS if field not in requirements]
    if missing:
        fail(f"dev_architecture_propagation.task_packet_requirements missing {missing}")
    return baseline_task_id


def validate_task_packets(packets: dict[str, Any], baseline_task_id: str) -> None:
    tasks = packets.get("tasks")
    if not isinstance(tasks, list):
        fail("task-packets.json must contain tasks list")
    tasks_by_id = {task_id(task): task for task in tasks if isinstance(task, dict) and task_id(task)}
    scoped_tasks = [
        task
        for candidate_id, task in tasks_by_id.items()
        if depends_on(candidate_id, baseline_task_id, tasks_by_id)
    ]
    if not scoped_tasks:
        fail(f"no task packets are at or after propagation baseline {baseline_task_id}")
    for task in scoped_tasks:
        missing = [field for field in REQUIRED_TASK_FIELDS if not field_has_evidence(task, field)]
        if missing:
            fail(f"{task_id(task)} missing guide propagation fields: {', '.join(missing)}")


def validate_validation_contract(contract: dict[str, Any]) -> None:
    stop_conditions = contract.get("stop_conditions")
    if not isinstance(stop_conditions, list):
        fail("validation-contract.json must contain stop_conditions list")
    combined = "\n".join(str(value).lower() for value in stop_conditions)
    missing = [
        category
        for category, tokens in STOP_CONDITION_CATEGORIES.items()
        if not any(token in combined for token in tokens)
    ]
    if missing:
        fail(f"validation-contract.json stop_conditions missing guide categories: {missing}")


def main() -> int:
    try:
        validate_guides()
        plan = load_json(EXECUTION_PLAN)
        packets = load_json(TASK_PACKETS)
        contract = load_json(VALIDATION_CONTRACT)
        baseline_task_id = validate_execution_plan(plan)
        validate_task_packets(packets, baseline_task_id)
        validate_validation_contract(contract)
    except AssertionError as exc:
        print(f"repo-pack validation failed: {exc}", file=sys.stderr)
        return 2
    print("repo-pack validation passed")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
