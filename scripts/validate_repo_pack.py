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
REPAIR_TASK_PACKETS = [
    ROOT / ".factory" / "artifacts" / "pilot" / "lumyn-mvp-slice" / "repair-loop" / "task-packet.json"
]
TEST_MATRIX_SOURCE_BASE = "docs/dev/dev_guides.md"
ARCHITECTURE_GUIDE_BASE = "docs/architecture/architecture_guides.md"

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

REQUIRED_PLAN_SKILL_REFS = [
    "factory://skills/prd-to-plan",
    "factory://skills/execution-compiler",
]

REQUIRED_PLAN_LEVEL_FIELDS = [
    "planning_skill_alignment",
    "public_api_and_contract_map",
    "docs_and_oss_readiness_baseline",
    "test_matrix_wiring",
    "minimum_now_sequence",
    "explicit_non_goals",
    "definition_of_done",
]

REQUIRED_TASK_PLANNING_FIELDS = [
    "planning_skill_refs",
    "runtime_pins",
    "slice_rationale",
    "changelog",
    "contract_impact",
    "versioning_migration_impact",
    "architecture_constraints",
    "adr_required",
    "tdd_first_failing_tests",
    "cost_perf_impact",
    "chaos_failure_hypothesis",
    "semantic_invariants",
]

REQUIRED_RUNTIME_PIN_FIELDS = [
    "language",
    "go_version",
    "module_path",
    "dependency_policy",
    "distribution_target",
]

REQUIRED_CHANGELOG_FIELDS = [
    "impact",
    "section",
    "draft_entry",
    "semver_marker_override",
]

ADR_CONTRACT_TOKENS = [
    "public",
    "cli",
    "command",
    "schema",
    "artifact",
    "output",
    "json",
    "contract",
    "ci",
    "review",
    "redaction",
    "share",
    "eval",
    "proof",
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

ARCHITECTURE_POLICY_TOKENS = {
    "systems_thinking": ["systems-thinking", "systems thinking"],
    "tdd": ["tdd", "red-first"],
    "adr_triggers": ["adr", "decision"],
    "performance": ["performance", "cost"],
    "reliability": ["reliability", "recovery"],
    "failure_semantics": ["fail-closed", "failure", "trust-mode"],
}

STOP_CONDITION_CATEGORIES = {
    "test_matrix": ["test-matrix", "test matrix", "test_matrix"],
    "ci_lanes": ["ci lane", "ci/status", "status check"],
    "scanner": ["scanner", "security"],
    "engineering_policies": ["docs parity", "output contract", "release integrity", "provenance"],
    "architecture_policies": ["architecture", "systems-thinking", "systems thinking", "adr", "fail-closed"],
    "planning_skill": ["prd-to-plan", "execution-compiler", "planning-skill", "planning skill"],
    "contract_discipline": ["changelog", "contract/api", "semantic invariants", "semantic_invariants"],
}

TASK_ORDER_RE = re.compile(r"^T(?P<version>\d+(?:\.\d+)*)(?:[^.\d].*)?$", re.IGNORECASE)


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


def source_ref_base(value: Any) -> str:
    return value.split("#", 1)[0] if isinstance(value, str) else ""


def refs_include_base(task: dict[str, Any], field: str, expected_base: str) -> bool:
    value = task.get(field)
    if not isinstance(value, list):
        return False
    return any(isinstance(item, dict) and source_ref_base(item.get("source_ref")) == expected_base for item in value)


def has_nonempty_list(value: Any) -> bool:
    return isinstance(value, list) and any(isinstance(item, str) and item.strip() for item in value)


def has_nonempty_string(value: Any) -> bool:
    return isinstance(value, str) and bool(value.strip())


def has_nonempty_dict(value: Any) -> bool:
    return isinstance(value, dict) and bool(value)


def has_nonempty_collection(value: Any) -> bool:
    return (isinstance(value, dict) and bool(value)) or (isinstance(value, list) and bool(value))


def has_required_string_refs(value: Any, expected_refs: list[str]) -> bool:
    if not isinstance(value, list):
        return False
    present = {item for item in value if isinstance(item, str) and item.strip()}
    return all(expected in present for expected in expected_refs)


def contains_machine_local_path(value: Any) -> bool:
    if isinstance(value, str):
        return "/Users/" in value
    if isinstance(value, list):
        return any(contains_machine_local_path(item) for item in value)
    if isinstance(value, dict):
        return any(contains_machine_local_path(item) for item in value.values())
    return False


def missing_ci_lane_refs(task: dict[str, Any]) -> list[str]:
    value = task.get("ci_lane_refs")
    if not isinstance(value, list):
        return list(REQUIRED_CI_LANES)
    present = {item.get("lane") for item in value if isinstance(item, dict)}
    return [lane for lane in REQUIRED_CI_LANES if lane not in present]


def missing_engineering_policy_refs(task: dict[str, Any]) -> list[str]:
    value = task.get("engineering_policy_refs")
    if not isinstance(value, list):
        return list(REQUIRED_ENGINEERING_POLICIES)
    present = {item.get("policy") for item in value if isinstance(item, dict)}
    return [policy for policy in REQUIRED_ENGINEERING_POLICIES if policy not in present]


def missing_architecture_policy_refs(task: dict[str, Any]) -> list[str]:
    value = task.get("architecture_guidance_refs")
    if not isinstance(value, list):
        return list(REQUIRED_ARCHITECTURE_POLICIES)
    combined = "\n".join(
        f"{item.get('source_ref', '')} {item.get('rule', '')}".lower()
        for item in value
        if isinstance(item, dict)
    )
    return [
        policy
        for policy, tokens in ARCHITECTURE_POLICY_TOKENS.items()
        if not any(token in combined for token in tokens)
    ]


def at_or_after_baseline(task: dict[str, Any], baseline_task_id: str) -> bool:
    baseline_key = task_order_key(baseline_task_id)
    if baseline_key is None:
        return False
    for value in [task_id(task), task.get("phase")]:
        candidate_key = task_order_key(value)
        if candidate_key is not None and version_gte(candidate_key, baseline_key):
            return True
    return False


def validate_task_guide_sources(task: dict[str, Any]) -> None:
    task_id_value = task_id(task)
    if not refs_include_base(task, "test_matrix_refs", TEST_MATRIX_SOURCE_BASE):
        fail(f"{task_id_value} test_matrix_refs must include source {TEST_MATRIX_SOURCE_BASE}")
    if not refs_include_base(task, "architecture_guidance_refs", ARCHITECTURE_GUIDE_BASE):
        fail(f"{task_id_value} architecture_guidance_refs must include source {ARCHITECTURE_GUIDE_BASE}")
    missing_lanes = missing_ci_lane_refs(task)
    if missing_lanes:
        fail(f"{task_id_value} ci_lane_refs missing: {', '.join(missing_lanes)}")
    missing_engineering = missing_engineering_policy_refs(task)
    if missing_engineering:
        fail(f"{task_id_value} engineering_policy_refs missing: {', '.join(missing_engineering)}")
    missing_architecture = missing_architecture_policy_refs(task)
    if missing_architecture:
        fail(f"{task_id_value} architecture_guidance_refs missing: {', '.join(missing_architecture)}")


def task_planning_field_has_evidence(task: dict[str, Any], field: str) -> bool:
    value = task.get(field)
    if field == "planning_skill_refs":
        return has_required_string_refs(value, REQUIRED_PLAN_SKILL_REFS)
    if field == "runtime_pins":
        if not has_nonempty_dict(value):
            return False
        return all(has_nonempty_string(value.get(pin)) for pin in REQUIRED_RUNTIME_PIN_FIELDS)
    if field == "slice_rationale":
        return (
            has_nonempty_dict(value)
            and has_nonempty_string(value.get("slice_type"))
            and has_nonempty_string(value.get("why_this_shape"))
        )
    if field == "changelog":
        if not has_nonempty_dict(value):
            return False
        return all(has_nonempty_string(value.get(changelog_field)) for changelog_field in REQUIRED_CHANGELOG_FIELDS)
    if field in ["contract_impact", "versioning_migration_impact"]:
        return has_nonempty_string(value)
    if field in ["architecture_constraints", "tdd_first_failing_tests", "semantic_invariants"]:
        return has_nonempty_list(value)
    if field == "adr_required":
        return isinstance(value, bool)
    if field == "cost_perf_impact":
        return (
            has_nonempty_dict(value)
            and has_nonempty_string(value.get("level"))
            and has_nonempty_string(value.get("measurement_expectation"))
        )
    if field == "chaos_failure_hypothesis":
        return (
            has_nonempty_dict(value)
            and has_nonempty_string(value.get("hypothesis"))
            and has_nonempty_string(value.get("expected_fail_closed_behavior"))
        )
    return False


def validate_task_planning_skill_fields(task: dict[str, Any]) -> None:
    task_id_value = task_id(task)
    missing = [
        field for field in REQUIRED_TASK_PLANNING_FIELDS if not task_planning_field_has_evidence(task, field)
    ]
    if missing:
        fail(f"{task_id_value} missing planning-skill fields: {', '.join(missing)}")
    if task.get("slice_type") != "vertical" and not has_nonempty_string(task.get("non_vertical_justification")):
        fail(f"{task_id_value} non-vertical task requires non_vertical_justification")
    contract_impact = str(task.get("contract_impact", "")).lower()
    if any(token in contract_impact for token in ADR_CONTRACT_TOKENS) and task.get("adr_required") is not True:
        fail(f"{task_id_value} public or executable contract impact requires adr_required=true")
    if contains_machine_local_path(task):
        fail(f"{task_id_value} contains a machine-local /Users/ path")


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
            return isinstance(value.get("exception_ref"), str) and bool(value["exception_ref"].strip())
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
                has_nonempty_list(item.get("command_refs"))
                or has_nonempty_list(item.get("status_check_refs"))
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
    for field in REQUIRED_PLAN_LEVEL_FIELDS:
        value = plan.get(field)
        if not has_nonempty_collection(value):
            fail(f"execution plan missing required planning-skill section {field}")
    alignment = plan.get("planning_skill_alignment")
    if not isinstance(alignment, dict):
        fail("execution plan planning_skill_alignment must be an object")
    if alignment.get("status") != "aligned":
        fail("execution plan planning_skill_alignment.status must be aligned")
    if not has_required_string_refs(alignment.get("source_refs"), REQUIRED_PLAN_SKILL_REFS):
        fail("execution plan planning_skill_alignment.source_refs must include Factory planning skills")
    if contains_machine_local_path(plan):
        fail("execution plan contains a machine-local /Users/ path")

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
    missing_planning = [field for field in REQUIRED_TASK_PLANNING_FIELDS if field not in requirements]
    if missing_planning:
        fail(f"dev_architecture_propagation.task_packet_requirements missing planning fields {missing_planning}")
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
    if baseline_task_id not in tasks_by_id:
        fail(f"task-packets.json missing propagation baseline task {baseline_task_id}")
    scoped_tasks = []
    all_task_objects = []
    baseline_has_order_key = task_order_key(baseline_task_id) is not None
    baseline_seen = False
    for task in tasks:
        if not isinstance(task, dict):
            continue
        candidate_id = task_id(task)
        if not candidate_id:
            continue
        all_task_objects.append(task)
        if candidate_id == baseline_task_id:
            baseline_seen = True
        depends_on_baseline = depends_on(candidate_id, baseline_task_id, tasks_by_id)
        ordered_after_baseline = at_or_after_baseline(task, baseline_task_id) or (
            not baseline_has_order_key and baseline_seen
        )
        if ordered_after_baseline and not depends_on_baseline:
            fail(f"{candidate_id} is at or after propagation baseline {baseline_task_id} but does not depend on it")
        if depends_on_baseline or ordered_after_baseline:
            scoped_tasks.append(task)
    if not scoped_tasks:
        fail(f"no task packets are at or after propagation baseline {baseline_task_id}")
    for task in all_task_objects:
        missing = [field for field in REQUIRED_TASK_FIELDS if not field_has_evidence(task, field)]
        if missing:
            fail(f"{task_id(task)} missing guide propagation fields: {', '.join(missing)}")
        validate_task_guide_sources(task)
        validate_task_planning_skill_fields(task)


def validate_standalone_task_packet(packet: dict[str, Any], baseline_task_id: str) -> None:
    task_id_value = task_id(packet)
    if not task_id_value:
        fail("standalone task packet missing task_id")
    if not at_or_after_baseline(packet, baseline_task_id):
        return
    blocked_by = packet.get("blocked_by", [])
    if not isinstance(blocked_by, list) or baseline_task_id not in [str(value) for value in blocked_by]:
        fail(f"{task_id_value} is at or after propagation baseline {baseline_task_id} but does not depend on it")
    missing = [field for field in REQUIRED_TASK_FIELDS if not field_has_evidence(packet, field)]
    if missing:
        fail(f"{task_id_value} missing guide propagation fields: {', '.join(missing)}")
    validate_task_guide_sources(packet)
    validate_task_planning_skill_fields(packet)


def validate_validation_contract(contract: dict[str, Any]) -> None:
    alignment = contract.get("planning_skill_alignment")
    if not isinstance(alignment, dict):
        fail("validation-contract.json missing planning_skill_alignment")
    if not has_required_string_refs(alignment.get("source_refs"), REQUIRED_PLAN_SKILL_REFS):
        fail("validation-contract.json planning_skill_alignment.source_refs must include Factory planning skills")
    required_plan_sections = alignment.get("required_plan_sections")
    if not isinstance(required_plan_sections, list) or not all(
        section in required_plan_sections for section in REQUIRED_PLAN_LEVEL_FIELDS
    ):
        fail("validation-contract.json planning_skill_alignment.required_plan_sections is incomplete")
    required_task_fields = alignment.get("required_task_fields")
    if not isinstance(required_task_fields, list) or not all(
        field in required_task_fields for field in REQUIRED_TASK_PLANNING_FIELDS
    ):
        fail("validation-contract.json planning_skill_alignment.required_task_fields is incomplete")
    if contains_machine_local_path(contract):
        fail("validation-contract.json contains a machine-local /Users/ path")
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
                "lane": "fast",
                "source_ref": "docs/dev/dev_guides.md#ci-lane-mapping",
                "command_refs": ["make lint-fast", "make test-fast"],
            },
            {
                "lane": "core",
                "source_ref": "docs/dev/dev_guides.md#ci-lane-mapping",
                "command_refs": ["make test-contracts"],
            },
            {
                "lane": "acceptance",
                "source_ref": "docs/dev/dev_guides.md#ci-lane-mapping",
                "command_refs": [".factory/artifacts/prd-to-plan/lumyn-mvp/scope-closure-map.json"],
            },
            {
                "lane": "cross_platform",
                "source_ref": "docs/dev/dev_guides.md#ci-lane-mapping",
                "exception_ref": ".factory/artifacts/exceptions/cross-platform-deferred.json",
            },
            {
                "lane": "risk",
                "source_ref": "docs/dev/dev_guides.md#ci-lane-mapping",
                "status_check_refs": ["CodeQL analyze"],
            },
            {
                "lane": "release",
                "source_ref": "docs/dev/dev_guides.md#ci-lane-mapping",
                "exception_ref": ".factory/artifacts/exceptions/release-deferred.json",
            },
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
            },
            {
                "policy": "output_contracts",
                "source_ref": "docs/dev/dev_guides.md#output-contracts",
                "rule": "output contracts stay stable",
            },
            {
                "policy": "release_integrity",
                "source_ref": "docs/dev/dev_guides.md#release-integrity",
                "rule": "release integrity evidence is required",
            },
            {
                "policy": "provenance_evidence",
                "source_ref": "docs/dev/dev_guides.md#provenance-evidence",
                "rule": "provenance evidence stays repo-relative",
            },
        ],
        "architecture_guidance_refs": [
            {"source_ref": "docs/architecture/architecture_guides.md#systems-thinking-map", "rule": "record state and feedback"},
            {
                "source_ref": "docs/architecture/architecture_guides.md#tdd-and-red-first-expectations",
                "rule": "use TDD and red-first evidence",
            },
            {
                "source_ref": "docs/architecture/architecture_guides.md#adr-and-decision-triggers",
                "rule": "record ADR decision triggers",
            },
            {
                "source_ref": "docs/architecture/architecture_guides.md#performance-and-cost-triggers",
                "rule": "track performance and cost impact",
            },
            {
                "source_ref": "docs/architecture/architecture_guides.md#reliability-and-recovery-triggers",
                "rule": "record reliability and recovery evidence",
            },
            {
                "source_ref": "docs/architecture/architecture_guides.md#trust-mode-posture",
                "rule": "fail-closed trust-mode posture",
            },
        ],
        "planning_skill_refs": list(REQUIRED_PLAN_SKILL_REFS),
        "runtime_pins": {
            "language": "go",
            "go_version": "1.26.4",
            "module_path": "github.com/Clyra-AI/lumyn",
            "dependency_policy": "standard library first; pinned dependencies only when task-required",
            "distribution_target": "standalone_binary",
        },
        "slice_type": "vertical",
        "slice_rationale": {
            "slice_type": "vertical",
            "why_this_shape": "self-test task preserves a vertically scoped implementation contract",
        },
        "changelog": {
            "impact": "required_when_implemented",
            "section": "Unreleased",
            "draft_entry": "Self-test task changelog entry.",
            "semver_marker_override": "pre_1_0_minor_candidate",
        },
        "contract_impact": "Self-test task changes only its declared contract surface.",
        "versioning_migration_impact": "Pre-1.0 changes must preserve explicit migration notes before release.",
        "architecture_constraints": ["record state owner, feedback source, and fail-closed behavior"],
        "adr_required": True,
        "tdd_first_failing_tests": ["add a failing test or fixture before implementation when practical"],
        "cost_perf_impact": {"level": "low", "measurement_expectation": "no material cost increase expected"},
        "chaos_failure_hypothesis": {
            "hypothesis": "invalid evidence must fail closed",
            "expected_fail_closed_behavior": "do not mark the task complete",
        },
        "semantic_invariants": ["evidence remains repo-relative", "closure cannot claim missing PRD scope"],
    }


def run_self_test() -> int:
    valid_packets = {"tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]}
    validate_task_packets(valid_packets, "T2.6")

    disconnected_packets = {"tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3-repair-001", [])]}
    try:
        validate_task_packets(disconnected_packets, "T2.6")
    except AssertionError as exc:
        if "does not depend on it" not in str(exc):
            raise
    else:
        fail("self-test expected disconnected T3 task to fail")

    slug_baseline_packets = {
        "tasks": [
            propagated_task("task-context", []),
            propagated_task("task-dev-architecture-propagation", ["task-context"]),
            propagated_task("feature-local-check", []),
        ]
    }
    try:
        validate_task_packets(slug_baseline_packets, "task-dev-architecture-propagation")
    except AssertionError as exc:
        if "does not depend on it" not in str(exc):
            raise
    else:
        fail("self-test expected disconnected task after slug baseline to fail")

    placeholder_packets = {"tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]}
    placeholder_packets["tasks"][1]["ci_lane_refs"] = [{}]
    try:
        validate_task_packets(placeholder_packets, "T2.6")
    except AssertionError as exc:
        if "ci_lane_refs" not in str(exc):
            raise
    else:
        fail("self-test expected placeholder refs to fail")

    disabled_scanner_packets = {"tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]}
    disabled_scanner_packets["tasks"][1]["security_scanner_gates"] = {"required": False, "scanner": "CodeQL"}
    try:
        validate_task_packets(disabled_scanner_packets, "T2.6")
    except AssertionError as exc:
        if "security_scanner_gates" not in str(exc):
            raise
    else:
        fail("self-test expected disabled scanner without exception to fail")

    missing_policy_packets = {"tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]}
    missing_policy_packets["tasks"][1]["engineering_policy_refs"] = missing_policy_packets["tasks"][1][
        "engineering_policy_refs"
    ][:1]
    try:
        validate_task_packets(missing_policy_packets, "T2.6")
    except AssertionError as exc:
        if "engineering_policy_refs missing" not in str(exc):
            raise
    else:
        fail("self-test expected missing engineering policy refs to fail")

    missing_ci_packets = {"tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]}
    missing_ci_packets["tasks"][1]["ci_lane_refs"] = [
        item for item in missing_ci_packets["tasks"][1]["ci_lane_refs"] if item["lane"] == "core"
    ]
    try:
        validate_task_packets(missing_ci_packets, "T2.6")
    except AssertionError as exc:
        if "ci_lane_refs missing" not in str(exc):
            raise
    else:
        fail("self-test expected missing CI lane refs to fail")

    blank_ci_packets = {"tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]}
    blank_ci_packets["tasks"][1]["ci_lane_refs"][1]["command_refs"] = [""]
    try:
        validate_task_packets(blank_ci_packets, "T2.6")
    except AssertionError as exc:
        if "missing guide propagation fields" not in str(exc):
            raise
    else:
        fail("self-test expected blank CI lane evidence to fail")

    missing_planning_packets = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    del missing_planning_packets["tasks"][1]["changelog"]
    try:
        validate_task_packets(missing_planning_packets, "T2.6")
    except AssertionError as exc:
        if "missing planning-skill fields" not in str(exc):
            raise
    else:
        fail("self-test expected missing planning-skill fields to fail")

    foundation_without_justification = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    foundation_without_justification["tasks"][1]["slice_type"] = "foundation"
    foundation_without_justification["tasks"][1]["slice_rationale"]["slice_type"] = "foundation"
    try:
        validate_task_packets(foundation_without_justification, "T2.6")
    except AssertionError as exc:
        if "non_vertical_justification" not in str(exc):
            raise
    else:
        fail("self-test expected non-vertical task without justification to fail")

    validate_standalone_task_packet(propagated_task("T3-repair-001", ["T2.6"]), "T2.6")

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
        for packet_path in REPAIR_TASK_PACKETS:
            validate_standalone_task_packet(load_json(packet_path), baseline_task_id)
        validate_validation_contract(contract)
    except AssertionError as exc:
        print(f"repo-pack validation failed: {exc}", file=sys.stderr)
        return 2
    print("repo-pack validation passed")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
