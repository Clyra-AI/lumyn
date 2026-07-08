#!/usr/bin/env python3
from __future__ import annotations

import json
import re
import sys
from pathlib import Path
from typing import Any

from repo_pack_acceptance import (
    validate_acceptance_ledger,
    validate_acceptance_mapping,
    validate_scope_closure_map,
)
from repo_pack_contracts import (
    ACCEPTANCE_LEDGER_REF,
    REQUIRED_ACCEPTANCE_ITEM_IDS,
    REQUIRED_LIVE_EVAL_DISPATCH_GATES,
    expected_task_version_slices,
    has_nonempty_list,
    has_nonempty_string,
    has_required_string_refs,
    validate_no_legacy_provider_fields,
)
from repo_pack_ci import (
    ref_file_exists as ref_file_exists_for_root,
    require_existing as require_existing_for_root,
    validate_ci_control_set as validate_ci_control_set_for_root,
    validate_coverage_policy_refs as validate_coverage_policy_refs_for_root,
    validate_guides as validate_guides_for_root,
)
from repo_pack_factoryd import (
    contains_machine_local_path,
    is_valid_factoryd_runtime,
    validate_factoryd_config as validate_factoryd_config_for_root,
    validate_factoryd_runtime,
)
from repo_pack_model_provider import (
    factoryd_config_capability_grants as collect_factoryd_config_capability_grants,
    validate_model_provider_gate as validate_model_provider_gate_with_grants,
)
from repo_pack_planning import (
    REQUIRED_ARCHITECTURE_POLICIES,
    REQUIRED_CHANGELOG_FIELDS,
    REQUIRED_CI_LANES,
    REQUIRED_ENGINEERING_POLICIES,
    REQUIRED_PLAN_SKILL_REFS,
    REQUIRED_TASK_FIELDS,
    REQUIRED_TASK_PLANNING_FIELDS,
    has_adr_contract_token,
    has_factory_compatibility,
    has_runtime_pins,
    validate_context_brief,
    validate_execution_plan,
    validate_factory_compatibility,
    validate_mvp_eval_provider_adapters,
    validate_no_deprecated_active_workers,
    validate_runtime_pins,
    validate_validation_contract,
)
from repo_pack_safety import (
    validate_risk_classification as validate_risk_classification_with_policy,
    validate_safety_corpus_ready_plan as validate_safety_corpus_ready_plan_with_prd,
)
from repo_pack_task_specials import validate_first_session_smoke_task, validate_recorder_task_split
from repo_pack_task_paths import validate_architecture_target_paths

ROOT = Path(__file__).resolve().parents[1]
PLAN_DIR = ROOT / ".factory" / "artifacts" / "prd-to-plan" / "lumyn-mvp"
CONTEXT_BRIEF = PLAN_DIR / "context-brief.json"
EXECUTION_PLAN = PLAN_DIR / "execution-plan.json"
TASK_PACKETS = PLAN_DIR / "task-packets.json"
VALIDATION_CONTRACT = PLAN_DIR / "validation-contract.json"
ACCEPTANCE_LEDGER = PLAN_DIR / "acceptance-ledger.json"
ACCEPTANCE_MAPPING = PLAN_DIR / "acceptance-mapping.json"
SCOPE_CLOSURE_MAP = PLAN_DIR / "scope-closure-map.json"
RISK_CLASSIFICATION = PLAN_DIR / "risk-classification.json"
FACTORYD_CONFIG = ROOT / ".factory" / "factoryd.example.json"
FACTORYD_ACTIVE_CONFIG = ROOT / ".factory" / "factoryd.json"
FACTORYD_AUTOSHIP_CONFIG = ROOT / ".factory" / "factoryd.autoship.example.json"
FACTORYD_REPO_KEY = "lumyn"
REQUIRED_CHECKS = ROOT / ".github" / "required-checks.json"
CODEOWNERS = ROOT / ".github" / "CODEOWNERS"
ACTION_REF_EXCEPTIONS = ROOT / ".github" / "action-ref-exceptions.yaml"
VALIDATE_WORKFLOW = ROOT / ".github" / "workflows" / "validate.yml"
CODEQL_WORKFLOW = ROOT / ".github" / "workflows" / "codeql.yml"
REPAIR_TASK_PACKETS = [
    ROOT / ".factory" / "artifacts" / "pilot" / "lumyn-mvp-slice" / "repair-loop" / "task-packet.json"
]
PRD = ROOT / "docs" / "product" / "prd.md"
TEST_MATRIX_SOURCE_BASE = "docs/dev/dev_guides.md"
COVERAGE_POLICY_SOURCE_BASE = "docs/dev/dev_guides.md"
ARCHITECTURE_GUIDE_BASE = "docs/architecture/architecture_guides.md"
RUNTIME_CONTROL_ALLOWED_RE = re.compile(
    r"^\.factory/artifacts/?$"
    r"|^\.factory/artifacts/(prd-to-plan|post-prd)(/.*)?$"
)
RUNTIME_CONTROL_FORBIDDEN_PATHS = [
    ".factory/artifacts/prd-to-plan/lumyn-mvp/",
    ".factory/artifacts/prd-to-plan/lumyn-mvp/context-brief.json",
    ".factory/artifacts/prd-to-plan/lumyn-mvp/execution-plan.json",
    ".factory/artifacts/prd-to-plan/lumyn-mvp/task-packets.json",
    ".factory/artifacts/prd-to-plan/lumyn-mvp/validation-contract.json",
    ".factory/artifacts/prd-to-plan/lumyn-mvp/acceptance-ledger.json",
    ".factory/artifacts/prd-to-plan/lumyn-mvp/acceptance-mapping.json",
    ".factory/artifacts/prd-to-plan/lumyn-mvp/scope-closure-map.json",
]
REQUIRED_GUIDES = [
    "docs/dev/dev_guides.md",
    "docs/architecture/architecture_guides.md",
]


def normalize_repo_path(value: object) -> str:
    path = str(value).strip().replace("\\", "/")
    parts: list[str] = []
    for part in path.split("/"):
        if part in {"", "."}:
            continue
        if part == "..":
            if parts:
                parts.pop()
            else:
                parts.append(part)
            continue
        parts.append(part)
    return "/".join(parts)


REQUIRED_PROOF_LEVELS = {
    "syntax",
    "source_evidence",
    "workflow_behavior",
    "user_visible_behavior",
}
BEHAVIORAL_PROOF_LEVELS = {
    "workflow_behavior",
    "user_visible_behavior",
}
PROOF_SCORECARD_ARTIFACT = "proof-of-behavior-scorecard"
REDACTION_RECURSIVE_TERMS = {
    "owner",
    "credential",
    "secret",
    "endpoint",
    "path",
}

DEFAULT_REQUIRED_WORKER_CHAIN = [
    "task-executor",
    "validation-gate",
    "commit-push",
    "post-merge-monitor",
]

REVIEW_REQUIRED_WORKER_CHAIN = [
    "task-executor",
    "validation-gate",
    "code-review",
    "commit-push",
    "post-merge-monitor",
]

REQUIRED_STATUS_CHECKS = [
    "validate",
    "CodeQL analyze",
]

REQUIRED_ACTION_REFS = [
    "actions/checkout@v6.0.2",
    "actions/setup-go@v6.3.0",
    "github/codeql-action/init@v4",
    "github/codeql-action/autobuild@v4",
    "github/codeql-action/analyze@v4",
]

ARCHITECTURE_POLICY_TOKENS = {
    "systems_thinking": ["systems-thinking", "systems thinking"],
    "tdd": ["tdd", "red-first"],
    "adr_triggers": ["adr", "decision"],
    "performance": ["performance", "cost"],
    "reliability": ["reliability", "recovery"],
    "failure_semantics": ["fail-closed", "failure", "trust-mode"],
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


def factoryd_config_capability_grants() -> list[dict[str, Any]]:
    config = load_json(FACTORYD_ACTIVE_CONFIG) if FACTORYD_ACTIVE_CONFIG.exists() else {}
    return collect_factoryd_config_capability_grants(config, FACTORYD_REPO_KEY)


def validate_model_provider_gate(
    task: dict[str, Any],
    active_grants: list[dict[str, Any]] | None = None,
) -> None:
    grants = factoryd_config_capability_grants() if active_grants is None else active_grants
    validate_model_provider_gate_with_grants(task, grants)


def require_existing(relative_path: str) -> None:
    require_existing_for_root(ROOT, relative_path)


def ref_file_exists(ref: Any) -> bool:
    return ref_file_exists_for_root(ROOT, ref)


def validate_coverage_policy_refs(value: Any, label: str) -> None:
    validate_coverage_policy_refs_for_root(ROOT, value, label)


def validate_guides() -> None:
    validate_guides_for_root(ROOT, REQUIRED_GUIDES)


def validate_ci_control_set() -> None:
    validate_ci_control_set_for_root(
        ROOT,
        REQUIRED_CHECKS,
        CODEOWNERS,
        ACTION_REF_EXCEPTIONS,
        VALIDATE_WORKFLOW,
        CODEQL_WORKFLOW,
        REQUIRED_STATUS_CHECKS,
        REQUIRED_ACTION_REFS,
    )


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


def object_source_ref_base(value: Any) -> str:
    if not isinstance(value, dict):
        return ""
    return source_ref_base(value.get("source_ref"))


def has_nonempty_dict(value: Any) -> bool:
    return isinstance(value, dict) and bool(value)


def task_slice_type(task: dict[str, Any]) -> str:
    rationale = task.get("slice_rationale")
    nested = str(rationale["slice_type"]) if isinstance(rationale, dict) and has_nonempty_string(rationale.get("slice_type")) else ""
    top_level = task.get("slice_type")
    top_level = top_level if isinstance(top_level, str) else ""
    if nested and top_level and nested != top_level:
        fail(f"{task_id(task)} has conflicting slice_type declarations")
    return nested or top_level


def has_lifecycle_gates(value: Any) -> bool:
    required_true = [
        "local_validation_required",
        "ci_required",
        "codex_review_required",
        "commit_push_required",
        "post_merge_monitor_required",
        "pr_lifecycle_report_required",
    ]
    if not isinstance(value, dict):
        return False
    exception_ref = value.get("exception_ref")
    has_exception = isinstance(exception_ref, str) and bool(exception_ref.strip())
    review_gate_is_declared = isinstance(value.get("code_review_required"), bool)
    return review_gate_is_declared and all(value.get(field) is True or has_exception for field in required_true)


def expected_required_worker_chain(task: dict[str, Any]) -> list[str]:
    gates = task.get("lifecycle_gates")
    if isinstance(gates, dict) and gates.get("code_review_required") is True:
        return REVIEW_REQUIRED_WORKER_CHAIN
    return DEFAULT_REQUIRED_WORKER_CHAIN


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


def is_live_eval_dispatch_task(task: dict[str, Any]) -> bool:
    for value in [task_id(task), task.get("phase")]:
        key = task_order_key(value)
        if key is not None and key[:1] in {(11,), (12,)}:
            return True
    return False


def validate_task_guide_sources(task: dict[str, Any]) -> None:
    task_id_value = task_id(task)
    if not refs_include_base(task, "test_matrix_refs", TEST_MATRIX_SOURCE_BASE):
        fail(f"{task_id_value} test_matrix_refs must include source {TEST_MATRIX_SOURCE_BASE}")
    if object_source_ref_base(task.get("coverage_policy_refs")) != COVERAGE_POLICY_SOURCE_BASE:
        fail(f"{task_id_value} coverage_policy_refs must include source {COVERAGE_POLICY_SOURCE_BASE}")
    validate_coverage_policy_refs(task.get("coverage_policy_refs"), f"{task_id_value}.coverage_policy_refs")
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
        return has_runtime_pins(value)
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
    if task_slice_type(task) != "vertical" and not has_nonempty_string(task.get("non_vertical_justification")):
        fail(f"{task_id_value} non-vertical task requires non_vertical_justification")
    contract_impact = str(task.get("contract_impact", ""))
    if has_adr_contract_token(contract_impact) and task.get("adr_required") is not True:
        fail(f"{task_id_value} public or executable contract impact requires adr_required=true")
    if contains_machine_local_path(task):
        fail(f"{task_id_value} contains a machine-local absolute path")


def validate_task_execution_compiler_fields(task: dict[str, Any]) -> None:
    task_id_value = task_id(task)
    validate_factory_compatibility(task.get("factory_compatibility"), f"{task_id_value}.factory_compatibility")
    validate_runtime_pins(task.get("runtime_pins"), f"{task_id_value}.runtime_pins")
    if not has_nonempty_list(task.get("scope_exclusions")):
        fail(f"{task_id_value}.scope_exclusions must preserve explicit PRD non-goals")
    if task.get("alignment_gate_ref") != ".factory/artifacts/prd-to-plan/lumyn-mvp/execution-plan.json#/alignment_gate":
        fail(f"{task_id_value}.alignment_gate_ref must cite the execution-plan alignment gate")
    if task.get("plan_drift_policy_ref") != ".factory/artifacts/prd-to-plan/lumyn-mvp/execution-plan.json#/plan_drift_policy":
        fail(f"{task_id_value}.plan_drift_policy_ref must cite the execution-plan drift policy")
    if task.get("acceptance_ledger_ref") != ACCEPTANCE_LEDGER_REF:
        fail(f"{task_id_value}.acceptance_ledger_ref must cite {ACCEPTANCE_LEDGER_REF}")
    item_ids = task.get("acceptance_item_ids")
    if not isinstance(item_ids, list) or not item_ids:
        fail(f"{task_id_value}.acceptance_item_ids must be non-empty")
    unknown_item_ids = sorted(str(value) for value in item_ids if str(value) not in REQUIRED_ACCEPTANCE_ITEM_IDS)
    if unknown_item_ids:
        fail(f"{task_id_value}.acceptance_item_ids references unknown ids: {unknown_item_ids}")
    validate_acceptance_item_gates(task)
    inherited = task.get("validation_contract_inheritance")
    if isinstance(inherited, dict):
        if inherited.get("acceptance_ledger_ref") != ACCEPTANCE_LEDGER_REF:
            fail(f"{task_id_value}.validation_contract_inheritance.acceptance_ledger_ref must cite {ACCEPTANCE_LEDGER_REF}")
        inherited_ids = inherited.get("acceptance_item_ids")
        if not isinstance(inherited_ids, list) or not set(str(value) for value in item_ids).issubset({str(value) for value in inherited_ids}):
            fail(f"{task_id_value}.validation_contract_inheritance.acceptance_item_ids must include task acceptance_item_ids")
    if not has_lifecycle_gates(task.get("lifecycle_gates")):
        fail(f"{task_id_value}.lifecycle_gates must enable local, CI, Codex review, ship, post-merge, and PR lifecycle gates, and explicitly declare code_review_required true only when review policy requires it")
    if task.get("required_worker_chain") != expected_required_worker_chain(task):
        fail(f"{task_id_value}.required_worker_chain must match the lifecycle gates: default validation/commit-push chain, or validation/code-review/commit-push chain when code_review_required=true")
    validate_architecture_target_paths(task)
    allowed_paths = [str(value).strip() for value in task.get("allowed_paths", [])]
    bad_allowed = [path for path in allowed_paths if RUNTIME_CONTROL_ALLOWED_RE.match(normalize_repo_path(path))]
    if bad_allowed:
        fail(f"{task_id_value}.allowed_paths includes runtime-owned control artifact paths: {bad_allowed}")
    forbidden_paths = set(str(value).strip() for value in task.get("forbidden_paths", []))
    missing_forbidden = [path for path in RUNTIME_CONTROL_FORBIDDEN_PATHS if path not in forbidden_paths]
    if missing_forbidden:
        fail(f"{task_id_value}.forbidden_paths missing runtime-owned control paths: {missing_forbidden}")


def validate_acceptance_item_gates(task: dict[str, Any]) -> None:
    task_id_value = task_id(task)
    gates = task.get("gated_by_acceptance_items")
    if gates is None:
        return
    if not isinstance(gates, list) or not gates:
        fail(f"{task_id_value}.gated_by_acceptance_items must be a non-empty list when present")
    seen: set[str] = set()
    for index, gate in enumerate(gates):
        if not isinstance(gate, dict):
            fail(f"{task_id_value}.gated_by_acceptance_items[{index}] must be an object")
        item_id_value = str(gate.get("acceptance_item_id", "")).strip()
        if not item_id_value:
            fail(f"{task_id_value}.gated_by_acceptance_items[{index}].acceptance_item_id is required")
        if item_id_value not in REQUIRED_ACCEPTANCE_ITEM_IDS:
            fail(f"{task_id_value}.gated_by_acceptance_items[{index}] references unknown acceptance item {item_id_value}")
        if item_id_value in seen:
            fail(f"{task_id_value}.gated_by_acceptance_items contains duplicate gate {item_id_value}")
        seen.add(item_id_value)
        if gate.get("required_status") not in {"implemented", "deferred_with_approval"}:
            fail(f"{task_id_value}.gated_by_acceptance_items[{index}].required_status must be implemented or deferred_with_approval")
        if not has_nonempty_string(gate.get("reason")):
            fail(f"{task_id_value}.gated_by_acceptance_items[{index}].reason is required")


def validate_live_eval_dispatch_gates(task: dict[str, Any]) -> None:
    task_id_value = task_id(task)
    gates = task.get("gated_by_acceptance_items")
    if not isinstance(gates, list):
        fail(f"{task_id_value}.gated_by_acceptance_items must gate live eval dispatch")
    by_id = {
        str(gate.get("acceptance_item_id")): gate
        for gate in gates
        if isinstance(gate, dict)
    }
    missing = sorted(REQUIRED_LIVE_EVAL_DISPATCH_GATES - set(by_id))
    if missing:
        fail(f"{task_id_value}.gated_by_acceptance_items missing live eval pull gates: {missing}")
    for required_id in REQUIRED_LIVE_EVAL_DISPATCH_GATES:
        gate = by_id[required_id]
        if gate.get("required_status") != "implemented":
            fail(f"{task_id_value}.gated_by_acceptance_items[{required_id}].required_status must be implemented")
        if gate.get("evidence_mode") != "product_signal":
            fail(f"{task_id_value}.gated_by_acceptance_items[{required_id}].evidence_mode must be product_signal")


def validate_task_version_slice_refs(task: dict[str, Any]) -> None:
    task_id_value = task_id(task)
    expected = expected_task_version_slices(task_id_value)
    if not expected:
        return
    actual = task.get("mvp_required_version_slices")
    if not isinstance(actual, list) or not actual:
        fail(f"{task_id_value}.mvp_required_version_slices must map task to required MVP version slices")
    missing = sorted(expected - {str(value) for value in actual})
    if missing:
        fail(f"{task_id_value}.mvp_required_version_slices missing {missing}")
    delivery_refs = task.get("delivery_slice_refs")
    if not isinstance(delivery_refs, list) or not delivery_refs:
        fail(f"{task_id_value}.delivery_slice_refs must map task to generic delivery slices")
    delivery_missing = sorted(expected - {str(value) for value in delivery_refs})
    if delivery_missing:
        fail(f"{task_id_value}.delivery_slice_refs missing {delivery_missing}")
    unexpected = sorted({str(value) for value in delivery_refs} - {str(value) for value in actual})
    if unexpected:
        fail(f"{task_id_value}.delivery_slice_refs has refs not present in mvp_required_version_slices: {unexpected}")


def field_has_evidence(task: dict[str, Any], field: str) -> bool:
    value = task.get(field)
    if field == "factory_compatibility":
        return has_factory_compatibility(value)
    if field == "scope_exclusions":
        return has_nonempty_list(value)
    if field in ["alignment_gate_ref", "plan_drift_policy_ref", "acceptance_ledger_ref", "path_planning_method"]:
        return has_nonempty_string(value)
    if field == "required_worker_chain":
        return value == expected_required_worker_chain(task)
    if field == "lifecycle_gates":
        return has_lifecycle_gates(value)
    if field in ["allowed_paths", "forbidden_paths", "architecture_target_paths"]:
        return has_nonempty_list(value)
    if field == "worker_type":
        return value == "codex_cli"
    if field == "factoryd_runtime":
        return is_valid_factoryd_runtime(value)
    if field == "max_iterations":
        return isinstance(value, int) and not isinstance(value, bool) and value > 0
    if field in ["validation_commands", "evidence_required", "stop_conditions"]:
        return has_nonempty_list(value)
    if field == "required_proof_level":
        if value not in REQUIRED_PROOF_LEVELS:
            return False
        if value in BEHAVIORAL_PROOF_LEVELS:
            evidence_required = task.get("evidence_required")
            return (
                task.get("proof_scorecard_required") is True
                and task.get("proof_scorecard_artifact") == PROOF_SCORECARD_ARTIFACT
                and isinstance(evidence_required, list)
                and PROOF_SCORECARD_ARTIFACT in evidence_required
            )
        return True
    if field == "artifact_budget_refs":
        return has_nonempty_list(value)
    if field == "redaction_posture":
        if not isinstance(value, dict):
            return False
        if value.get("classification") not in {"internal", "customer_safe", "public"}:
            return False
        if not isinstance(value.get("customer_safe"), bool):
            return False
        if value.get("classification") in {"customer_safe", "public"} or value.get("customer_safe") is True:
            policy = value.get("recursive_policy")
            if not isinstance(policy, str) or not policy.strip():
                return False
            normalized_policy = policy.lower()
            if not all(term in normalized_policy for term in REDACTION_RECURSIVE_TERMS):
                return False
        return True
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
    if field == "coverage_policy_refs":
        if not isinstance(value, dict):
            return False
        if value.get("required") is False and isinstance(value.get("exception_ref"), str) and value["exception_ref"].strip():
            return True
        if value.get("required") is not True:
            return False
        if has_nonempty_list(value.get("command_refs")) or has_nonempty_list(value.get("evidence_refs")):
            return True
        if isinstance(value.get("exception_ref"), str) and value["exception_ref"].strip():
            return True
        minimums = value.get("minimums")
        return isinstance(minimums, list) and any(
            isinstance(item, dict) and has_nonempty_list(item.get("command_refs"))
            for item in minimums
        )
    if not isinstance(value, list) or not value:
        return False
    if field == "acceptance_item_ids":
        return all(has_nonempty_string(item) for item in value)
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


def validate_task_packets(
    packets: dict[str, Any],
    baseline_task_id: str,
    active_model_provider_grants: list[dict[str, Any]] | None = None,
) -> None:
    artifact_type = packets.get("artifact_type")
    if artifact_type is not None and artifact_type != "task_packets":
        fail("task-packets.json artifact_type must be task_packets")
    validate_no_legacy_provider_fields(packets, "task-packets.json")
    validate_no_deprecated_active_workers(packets, "task-packets.json")
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
        current_task_id = task_id(task)
        missing = [field for field in REQUIRED_TASK_FIELDS if not field_has_evidence(task, field)]
        if missing:
            fail(f"{current_task_id} missing guide propagation fields: {', '.join(missing)}")
        validate_task_guide_sources(task)
        validate_task_planning_skill_fields(task)
        validate_task_execution_compiler_fields(task)
        validate_task_version_slice_refs(task)
        item_count = len(task.get("acceptance_item_ids", []))
        if item_count > 15:
            fail(f"{current_task_id}.acceptance_item_ids has {item_count} items; split runner-ready tasks at 15 or fewer acceptance items")
        if current_task_id == "T11.1":
            validate_mvp_eval_provider_adapters(
                task.get("mvp_eval_provider_adapters"),
                "T11.1.mvp_eval_provider_adapters",
            )
            checks = "\n".join(str(value).lower() for value in task.get("acceptance_checks", []))
            if "openai-compatible" not in checks or "anthropic" not in checks:
                fail("T11.1 acceptance_checks must name both OpenAI-compatible and Anthropic adapter coverage")
            validate_model_provider_gate(task, active_model_provider_grants)
        if current_task_id == "T11.2":
            validate_model_provider_gate(task, active_model_provider_grants)
        if is_live_eval_dispatch_task(task):
            validate_live_eval_dispatch_gates(task)
    if any(task_ref in tasks_by_id for task_ref in ["T4", "T4.1", "T4.2", "T4.3"]):
        validate_recorder_task_split(tasks_by_id)
    if "T6.2" in tasks_by_id:
        validate_first_session_smoke_task(tasks_by_id)


def validate_standalone_task_packet(packet: dict[str, Any], baseline_task_id: str) -> None:
    validate_no_legacy_provider_fields(packet, "standalone task packet")
    validate_no_deprecated_active_workers(packet, "standalone task packet")
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
    validate_task_execution_compiler_fields(packet)
    validate_task_version_slice_refs(packet)
    if is_live_eval_dispatch_task(packet):
        validate_live_eval_dispatch_gates(packet)


def validate_safety_corpus_ready_plan(
    plan: dict[str, Any],
    packets: dict[str, Any],
    contract: dict[str, Any],
    ledger: dict[str, Any],
    mapping: dict[str, Any],
    scope: dict[str, Any],
) -> None:
    validate_safety_corpus_ready_plan_with_prd(PRD, plan, packets, contract, ledger, mapping, scope)


def validate_factoryd_config(config: dict[str, Any], active_config: dict[str, Any], autoship_config: dict[str, Any]) -> None:
    validate_factoryd_config_for_root(ROOT, config, active_config, autoship_config, FACTORYD_REPO_KEY)

def validate_risk_classification(risk: dict[str, Any]) -> None:
    validate_risk_classification_with_policy(risk)


def main() -> int:
    if sys.argv[1:] == ["--self-test"]:
        try:
            from repo_pack_self_test import run_self_test

            return run_self_test()
        except AssertionError as exc:
            print(f"repo-pack validator self-test failed: {exc}", file=sys.stderr)
            return 2
    if sys.argv[1:]:
        print("usage: validate_repo_pack.py [--self-test]", file=sys.stderr)
        return 2
    try:
        validate_guides()
        require_existing(".factory/artifacts/supervisor-runs/.gitkeep")
        require_existing(".factory/artifacts/task-supervisor-runs/.gitkeep")
        validate_ci_control_set()
        context = load_json(CONTEXT_BRIEF)
        plan = load_json(EXECUTION_PLAN)
        packets = load_json(TASK_PACKETS)
        contract = load_json(VALIDATION_CONTRACT)
        factoryd_config = load_json(FACTORYD_CONFIG)
        factoryd_active_config = load_json(FACTORYD_ACTIVE_CONFIG) if FACTORYD_ACTIVE_CONFIG.exists() else {}
        factoryd_autoship_config = load_json(FACTORYD_AUTOSHIP_CONFIG)
        acceptance_ledger = load_json(ACCEPTANCE_LEDGER)
        acceptance_mapping = load_json(ACCEPTANCE_MAPPING)
        scope_closure_map = load_json(SCOPE_CLOSURE_MAP)
        risk_classification = load_json(RISK_CLASSIFICATION)
        active_model_provider_grants = collect_factoryd_config_capability_grants(
            factoryd_active_config,
            FACTORYD_REPO_KEY,
        )
        validate_context_brief(context)
        baseline_task_id = validate_execution_plan(plan)
        ledger_ids = validate_acceptance_ledger(acceptance_ledger)
        if packets.get("artifact_type") != "task_packets":
            fail("task-packets.json artifact_type must be task_packets")
        validate_task_packets(packets, baseline_task_id, active_model_provider_grants)
        for packet_path in REPAIR_TASK_PACKETS:
            validate_standalone_task_packet(load_json(packet_path), baseline_task_id)
        validate_validation_contract(contract)
        validate_factoryd_config(factoryd_config, factoryd_active_config, factoryd_autoship_config)
        validate_acceptance_mapping(acceptance_mapping, ledger_ids, contract)
        validate_scope_closure_map(scope_closure_map, ledger_ids)
        validate_safety_corpus_ready_plan(
            plan,
            packets,
            contract,
            acceptance_ledger,
            acceptance_mapping,
            scope_closure_map,
        )
        validate_risk_classification(risk_classification)
    except AssertionError as exc:
        print(f"repo-pack validation failed: {exc}", file=sys.stderr)
        return 2
    print("repo-pack validation passed")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
