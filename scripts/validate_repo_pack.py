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

REQUIRED_CI_LANES = [
    "fast",
    "core",
    "acceptance",
    "cross_platform",
    "risk",
    "release",
]

REQUIRED_ENGINEERING_POLICIES = [
    "docs_parity",
    "output_contracts",
    "release_integrity",
    "provenance_evidence",
]

REQUIRED_ARCHITECTURE_POLICIES = [
    "systems_thinking",
    "tdd",
    "adr_triggers",
    "performance",
    "reliability",
    "failure_semantics",
]

STOP_CONDITION_CATEGORIES = {
    "test_matrix": ["test-matrix", "test matrix", "test_matrix"],
    "ci_lanes": ["ci lane", "ci/status", "status check"],
    "scanner": ["scanner", "security"],
    "engineering_policies": ["docs parity", "output contract", "release integrity", "provenance"],
    "architecture_policies": ["architecture", "systems-thinking", "systems thinking", "adr", "fail-closed"],
}

TASK_ORDER_RE = re.compile(r"^T(?P<version>\d+(?:\.\d+)*)(?:[A-Za-z].*)?$", re.IGNORECASE)


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


def task_order_key(value: Any) -> tuple[int, ...] | None:
    if not isinstance(value, str):
        return None
    match = TASK_ORDER_RE.match(value.strip())
    if not match:
        return None
    return tuple(int(part) for part in match.group("version").split("."))


def version_gte(candidate: tuple[int, ...], baseline: tuple[int, ...]) -> bool:
    width = max(len(candidate), len(baseline))
    return candidate + (0,) * (width - len(candidate)) >= baseline + (0,) * (width - len(baseline))


def at_or_after_baseline(task: dict[str, Any], baseline_task_id: str) -> bool:
    baseline_key = task_order_key(baseline_task_id)
    if baseline_key is None:
        return False
    for value in [task_id(task), task.get("phase")]:
        candidate_key = task_order_key(value)
        if candidate_key is not None and version_gte(candidate_key, baseline_key):
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
    security_scanning = propagation.get("security_scanning")
    if not isinstance(security_scanning, dict):
        fail("dev_architecture_propagation.security_scanning must be an object")
    if not isinstance(security_scanning.get("required"), bool):
        fail("dev_architecture_propagation.security_scanning.required must be boolean")
    if not isinstance(security_scanning.get("scanner"), str) or not security_scanning["scanner"].strip():
        fail("dev_architecture_propagation.security_scanning.scanner must be non-empty")
    if security_scanning.get("required") is True and not any(
        isinstance(security_scanning.get(key), str) and security_scanning[key].strip()
        for key in ["workflow_ref", "status_check", "exception_policy", "exception_ref"]
    ):
        fail("required dev_architecture_propagation.security_scanning needs workflow/status/exception evidence")
    ci_lanes = propagation.get("ci_lanes")
    if not isinstance(ci_lanes, dict):
        fail("dev_architecture_propagation.ci_lanes must be an object")
    missing_lanes = [
        lane
        for lane in REQUIRED_CI_LANES
        if not isinstance(ci_lanes.get(lane), list) or not ci_lanes[lane]
    ]
    if missing_lanes:
        fail(f"dev_architecture_propagation.ci_lanes missing non-empty lanes: {missing_lanes}")
    engineering = propagation.get("engineering_policies")
    if not isinstance(engineering, dict):
        fail("dev_architecture_propagation.engineering_policies must be an object")
    missing_engineering = [
        policy
        for policy in REQUIRED_ENGINEERING_POLICIES
        if not isinstance(engineering.get(policy), str) or not engineering[policy].strip()
    ]
    if missing_engineering:
        fail(f"dev_architecture_propagation.engineering_policies missing {missing_engineering}")
    architecture = propagation.get("architecture_policies")
    if not isinstance(architecture, dict):
        fail("dev_architecture_propagation.architecture_policies must be an object")
    missing_architecture = [
        policy
        for policy in REQUIRED_ARCHITECTURE_POLICIES
        if not isinstance(architecture.get(policy), str) or not architecture[policy].strip()
    ]
    if missing_architecture:
        fail(f"dev_architecture_propagation.architecture_policies missing {missing_architecture}")
    return baseline_task_id


def validate_task_packets(packets: dict[str, Any], baseline_task_id: str) -> None:
    tasks = packets.get("tasks")
    if not isinstance(tasks, list):
        fail("task-packets.json must contain tasks list")
    tasks_by_id = {task_id(task): task for task in tasks if isinstance(task, dict) and task_id(task)}
    scoped_tasks = []
    for candidate_id, task in tasks_by_id.items():
        depends_on_baseline = depends_on(candidate_id, baseline_task_id, tasks_by_id)
        ordered_after_baseline = at_or_after_baseline(task, baseline_task_id)
        if ordered_after_baseline and not depends_on_baseline:
            fail(f"{candidate_id} is at or after propagation baseline {baseline_task_id} but does not depend on it")
        if depends_on_baseline or ordered_after_baseline:
            scoped_tasks.append(task)
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


def propagated_task(task_id_value: str, blocked_by: list[str]) -> dict[str, Any]:
    return {
        "task_id": task_id_value,
        "blocked_by": blocked_by,
        "test_matrix_refs": [{"tier": "Tier 1 Unit", "source_ref": "docs/dev/dev_guides.md#12-level-test-matrix"}],
        "ci_lane_refs": [
            {
                "lane": "core",
                "source_ref": "docs/dev/dev_guides.md#ci-lane-mapping",
                "command_refs": ["make test-contracts"],
            }
        ],
        "security_scanner_gates": {
            "required": True,
            "scanner": "CodeQL",
            "workflow_ref": ".github/workflows/codeql.yml",
            "status_check": "CodeQL analyze",
        },
        "engineering_policy_refs": [
            {
                "policy": "docs_parity",
                "source_ref": "docs/dev/dev_guides.md#docs-parity",
                "rule": "docs move with behavior",
            }
        ],
        "architecture_guidance_refs": [
            {"source_ref": "docs/architecture/architecture_guides.md#systems-thinking-map", "rule": "record state and feedback"}
        ],
    }


def run_self_test() -> int:
    valid_packets = {"tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]}
    validate_task_packets(valid_packets, "T2.6")

    disconnected_packets = {"tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", [])]}
    try:
        validate_task_packets(disconnected_packets, "T2.6")
    except AssertionError as exc:
        if "does not depend on it" not in str(exc):
            raise
    else:
        fail("self-test expected disconnected T3 task to fail")

    placeholder_packets = {"tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]}
    placeholder_packets["tasks"][1]["ci_lane_refs"] = [{}]
    try:
        validate_task_packets(placeholder_packets, "T2.6")
    except AssertionError as exc:
        if "ci_lane_refs" not in str(exc):
            raise
    else:
        fail("self-test expected placeholder refs to fail")

    print("repo-pack validator self-test passed")
    return 0


def main() -> int:
    if sys.argv[1:] == ["--self-test"]:
        try:
            return run_self_test()
        except AssertionError as exc:
            print(f"repo-pack validator self-test failed: {exc}", file=sys.stderr)
            return 2
    if sys.argv[1:]:
        print("usage: validate_repo_pack.py [--self-test]", file=sys.stderr)
        return 2
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
