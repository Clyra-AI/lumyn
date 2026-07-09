#!/usr/bin/env python3
from __future__ import annotations

import re
from pathlib import Path
from typing import Any

from repo_pack_acceptance import validate_acceptance_ledger_coverage
from repo_pack_contracts import (
    ACCEPTANCE_LEDGER_REF,
    LEGACY_PROVIDER_FIELD,
    REQUIRED_ACCEPTANCE_ITEM_IDS,
    REQUIRED_MVP_EVAL_ADAPTERS,
    REQUIRED_MVP_EVAL_PROVIDERS,
    REQUIRED_PROVIDER_DECISION_ID,
    fail,
    has_nonempty_collection,
    has_nonempty_list,
    has_nonempty_string,
    has_required_string_refs,
    validate_delivery_slice_coverage,
    validate_mvp_version_slice_coverage,
    validate_no_legacy_provider_fields,
)
from repo_pack_ci import validate_coverage_policy_refs as validate_coverage_policy_refs_for_root
from repo_pack_factoryd import contains_machine_local_path, validate_factoryd_runtime


ROOT = Path(__file__).resolve().parents[1]

REQUIRED_TASK_FIELDS = [
    "ci_lane_refs",
    "test_matrix_refs",
    "coverage_policy_refs",
    "security_scanner_gates",
    "engineering_policy_refs",
    "architecture_guidance_refs",
    "factory_compatibility",
    "scope_exclusions",
    "alignment_gate_ref",
    "plan_drift_policy_ref",
    "required_worker_chain",
    "lifecycle_gates",
    "allowed_paths",
    "forbidden_paths",
    "architecture_target_paths",
    "path_planning_method",
    "worker_type",
    "factoryd_runtime",
    "validation_commands",
    "max_iterations",
    "evidence_required",
    "stop_conditions",
    "acceptance_ledger_ref",
    "acceptance_item_ids",
    "required_proof_level",
    "artifact_budget_refs",
    "redaction_posture",
]

REQUIRED_RUNNER_READY_FIELDS = [
    "worker_type",
    "factoryd_runtime",
    "validation_commands",
    "max_iterations",
    "evidence_required",
    "stop_conditions",
    "allowed_paths",
    "forbidden_paths",
    "architecture_target_paths",
    "path_planning_method",
    "semantic_invariants",
    "required_worker_chain",
    "lifecycle_gates",
    "scope_exclusions",
    "acceptance_ledger_ref",
    "acceptance_item_ids",
    "required_proof_level",
    "artifact_budget_refs",
    "redaction_posture",
]

REQUIRED_PLAN_SKILL_REFS = [
    "factory://skills/prd-to-plan",
    "factory://skills/execution-compiler",
]

DEPRECATED_ACTIVE_WORKERS = {
    "ship-pr": "commit-push",
}

REQUIRED_PLAN_LEVEL_FIELDS = [
    "planning_skill_alignment",
    "factory_compatibility",
    "runtime_pins",
    "alignment_gate",
    "plan_drift_policy",
    "acceptance_ledger_coverage",
    "mvp_required_version_slices",
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
    "toolchain_version",
    "module_path",
    "module_or_package_path",
    "dependency_policy",
    "distribution_target",
    "provider_policy",
    "artifact_namespace",
    "live_work_policy",
]

REQUIRED_FACTORY_COMPATIBILITY_FIELDS = [
    "factory_contract_version",
    "profile_ref",
    "skill_vocabulary_version",
    "skill_inventory_ref",
    "generated_by",
    "generated_at",
    "deprecated_worker_policy",
    "deprecated_worker_aliases",
]

REQUIRED_PLAN_DRIFT_UPDATES = [
    "context_brief",
    "execution_plan",
    "task_packets",
    "validation_contract",
    "factory_compatibility",
    "acceptance_ledger",
    "acceptance_mapping",
    "scope_closure_map",
]

REQUIRED_CHANGELOG_FIELDS = [
    "impact",
    "section",
    "draft_entry",
    "semver_marker_override",
]

ADR_CONTRACT_TOKENS = [
    "public",
    "api",
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

NO_CONTRACT_IMPACT_BREAKERS = [
    " but ",
    " except ",
    " however ",
    " although ",
    " changes ",
    " adds ",
    " removes ",
    " modifies ",
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
    "coverage_policy": ["coverage", "test-coverage"],
    "engineering_policies": ["docs parity", "output contract", "release integrity", "provenance"],
    "architecture_policies": ["architecture", "systems-thinking", "systems thinking", "adr", "fail-closed"],
    "planning_skill": ["prd-to-plan", "execution-compiler", "planning-skill", "planning skill"],
    "contract_discipline": ["changelog", "contract/api", "semantic invariants", "semantic_invariants"],
}


def declares_no_contract_impact(value: str) -> bool:
    collapsed = re.sub(r"\s+", " ", value.lower()).strip(" .")
    normalized = f" {collapsed} "
    if normalized.strip() in {"none", "n/a", "not applicable"}:
        return True
    if not normalized.startswith(" no "):
        return False
    if any(breaker in normalized for breaker in NO_CONTRACT_IMPACT_BREAKERS):
        return False
    return any(token in normalized for token in [" impact ", " change ", " changes ", " effect ", " effects "])


def has_adr_contract_token(value: str) -> bool:
    if declares_no_contract_impact(value):
        return False
    normalized = value.lower()
    return any(re.search(rf"\b{re.escape(token)}\b", normalized) for token in ADR_CONTRACT_TOKENS)


def iter_required_worker_chains(value: Any, path: str = "$") -> list[tuple[str, list[Any]]]:
    chains: list[tuple[str, list[Any]]] = []
    if isinstance(value, dict):
        for key, item in value.items():
            child_path = f"{path}.{key}"
            if key == "required_worker_chain" and isinstance(item, list):
                chains.append((child_path, item))
            chains.extend(iter_required_worker_chains(item, child_path))
    elif isinstance(value, list):
        for index, item in enumerate(value):
            chains.extend(iter_required_worker_chains(item, f"{path}[{index}]"))
    return chains


def validate_no_deprecated_active_workers(value: dict[str, Any], label: str) -> None:
    for path, chain in iter_required_worker_chains(value):
        for index, worker in enumerate(chain):
            if isinstance(worker, str) and worker in DEPRECATED_ACTIVE_WORKERS:
                replacement = DEPRECATED_ACTIVE_WORKERS[worker]
                fail(
                    f"{label}{path[1:]}[{index}] uses deprecated worker {worker!r}; "
                    f"use {replacement!r} in active required_worker_chain values"
                )


def has_factory_compatibility(value: Any) -> bool:
    if not isinstance(value, dict):
        return False
    if not all(field in value for field in REQUIRED_FACTORY_COMPATIBILITY_FIELDS):
        return False
    if value.get("deprecated_worker_policy") != "block_active_aliases":
        return False
    aliases = value.get("deprecated_worker_aliases")
    if not isinstance(aliases, list) or not aliases:
        return False
    return any(
        isinstance(alias, dict)
        and alias.get("deprecated") == "ship-pr"
        and alias.get("replacement") == "commit-push"
        for alias in aliases
    )


def validate_factory_compatibility(value: Any, label: str) -> None:
    if not has_factory_compatibility(value):
        fail(f"{label} must include current Factory compatibility metadata and ship-pr -> commit-push alias policy")


def has_runtime_pins(value: Any) -> bool:
    return (
        isinstance(value, dict)
        and all(has_nonempty_string(value.get(field)) for field in REQUIRED_RUNTIME_PIN_FIELDS)
        and has_required_string_refs(value.get("mvp_eval_providers"), REQUIRED_MVP_EVAL_PROVIDERS)
    )


def validate_runtime_pins(value: Any, label: str) -> None:
    if not isinstance(value, dict):
        fail(f"{label} must be an object")
    missing = [field for field in REQUIRED_RUNTIME_PIN_FIELDS if not has_nonempty_string(value.get(field))]
    if missing:
        fail(f"{label} missing runtime pin fields: {', '.join(missing)}")
    if not has_required_string_refs(value.get("mvp_eval_providers"), REQUIRED_MVP_EVAL_PROVIDERS):
        fail(
            f"{label}.mvp_eval_providers must include "
            f"{', '.join(REQUIRED_MVP_EVAL_PROVIDERS)}"
        )


def validate_mvp_eval_provider_adapters(value: Any, label: str) -> None:
    if not isinstance(value, dict):
        fail(f"{label} must be an object")
    if not has_required_string_refs(value.get("adapters"), REQUIRED_MVP_EVAL_ADAPTERS):
        fail(f"{label}.adapters must include {', '.join(REQUIRED_MVP_EVAL_ADAPTERS)}")
    if not has_required_string_refs(value.get("adapter_ids"), REQUIRED_MVP_EVAL_PROVIDERS):
        fail(f"{label}.adapter_ids must include {', '.join(REQUIRED_MVP_EVAL_PROVIDERS)}")
    if not has_required_string_refs(value.get("config_fields"), ["provider", "model", "temperature", "base_url", "api_key_env"]):
        fail(f"{label}.config_fields must include provider/model/temperature/base_url/api_key_env")


def has_alignment_gate(value: Any) -> bool:
    return (
        isinstance(value, dict)
        and value.get("status") == "resolved"
        and has_nonempty_string(value.get("source_context_brief_ref"))
        and has_nonempty_list(value.get("blocking_decisions"))
        and value.get("implementation_may_start") is True
    )


def validate_alignment_gate(value: Any, label: str) -> None:
    if not has_alignment_gate(value):
        fail(f"{label} must be resolved, cite the context brief, list blocking decisions, and allow implementation")


def has_plan_drift_policy(value: Any) -> bool:
    if not isinstance(value, dict):
        return False
    updates = value.get("required_updates")
    return (
        has_nonempty_list(value.get("drift_triggers"))
        and isinstance(updates, list)
        and all(required in updates for required in REQUIRED_PLAN_DRIFT_UPDATES)
        and value.get("continuation_behavior") == "block_until_artifacts_updated"
    )


def validate_plan_drift_policy(value: Any, label: str) -> None:
    if not has_plan_drift_policy(value):
        fail(
            f"{label} must require context brief, execution plan, task packets, validation contract, "
            "factory_compatibility, acceptance_ledger, acceptance_mapping, and scope_closure_map updates before continuing"
        )


def validate_semantic_invariant_policy_text(value: Any, label: str) -> None:
    text = str(value or "").lower()
    if "semantic_invariants" not in text and "semantic invariants" not in text:
        fail(f"{label} must include semantic_invariants in runner-ready dispatch policy")


def validate_context_brief(context: dict[str, Any]) -> None:
    validate_no_legacy_provider_fields(context, "context-brief.json")
    validate_factory_compatibility(context.get("factory_compatibility"), "context-brief.json.factory_compatibility")
    validate_alignment_gate(context.get("alignment_gate"), "context-brief.json.alignment_gate")
    validate_plan_drift_policy(context.get("plan_drift_policy"), "context-brief.json.plan_drift_policy")
    questions = context.get("alignment_questions")
    if not isinstance(questions, list):
        fail("context-brief.json.alignment_questions must be a list")
    question_ids = {question.get("id") for question in questions if isinstance(question, dict)}
    if LEGACY_PROVIDER_FIELD in question_ids:
        fail(f"context-brief.json must not use legacy alignment question id {LEGACY_PROVIDER_FIELD!r}")
    if REQUIRED_PROVIDER_DECISION_ID not in question_ids:
        fail(f"context-brief.json must include alignment question id {REQUIRED_PROVIDER_DECISION_ID!r}")
    decision_points = context.get("decision_points")
    if not isinstance(decision_points, list):
        fail("context-brief.json.decision_points must be a list")
    if LEGACY_PROVIDER_FIELD in decision_points:
        fail(f"context-brief.json must not use legacy decision point {LEGACY_PROVIDER_FIELD!r}")
    if REQUIRED_PROVIDER_DECISION_ID not in decision_points:
        fail(f"context-brief.json must include decision point {REQUIRED_PROVIDER_DECISION_ID!r}")
    decisions = context.get("alignment_decisions")
    if not isinstance(decisions, dict):
        fail("context-brief.json missing alignment_decisions")
    validate_factoryd_runtime(
        decisions.get("factoryd_runtime"),
        "context-brief.json.alignment_decisions.factoryd_runtime",
    )
    validate_semantic_invariant_policy_text(
        decisions.get("factoryd_dispatch_policy"),
        "context-brief.json.alignment_decisions.factoryd_dispatch_policy",
    )
    validate_mvp_eval_provider_adapters(
        decisions.get("mvp_eval_provider_adapters"),
        "context-brief.json.alignment_decisions.mvp_eval_provider_adapters",
    )
    if contains_machine_local_path(context):
        fail("context-brief.json contains a machine-local absolute path")


def validate_execution_plan(plan: dict[str, Any]) -> str:
    validate_no_legacy_provider_fields(plan, "execution-plan.json")
    validate_no_deprecated_active_workers(plan, "execution-plan.json")
    validate_factory_compatibility(plan.get("factory_compatibility"), "execution-plan.json.factory_compatibility")
    validate_runtime_pins(plan.get("runtime_pins"), "execution-plan.json.runtime_pins")
    validate_mvp_eval_provider_adapters(
        plan.get("mvp_eval_provider_adapters"),
        "execution-plan.json.mvp_eval_provider_adapters",
    )
    validate_factoryd_runtime(plan.get("factoryd_runtime"), "execution-plan.json.factoryd_runtime")
    validate_alignment_gate(plan.get("alignment_gate"), "execution-plan.json.alignment_gate")
    validate_plan_drift_policy(plan.get("plan_drift_policy"), "execution-plan.json.plan_drift_policy")
    validate_acceptance_ledger_coverage(
        plan.get("acceptance_ledger_coverage"),
        "execution-plan.json.acceptance_ledger_coverage",
    )
    validate_mvp_version_slice_coverage(
        plan.get("mvp_required_version_slices"),
        "execution-plan.json.mvp_required_version_slices",
    )
    validate_delivery_slice_coverage(
        plan.get("delivery_slices"),
        "execution-plan.json.delivery_slices",
    )
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
    runner_ready_rules = [
        str(rule)
        for rule in alignment.get("validated_rules") or []
        if "runner-ready" in str(rule).lower()
    ]
    if not runner_ready_rules:
        fail("execution plan planning_skill_alignment.validated_rules must include runner-ready dispatch policy")
    if not any(
        "semantic_invariants" in rule.lower() or "semantic invariants" in rule.lower()
        for rule in runner_ready_rules
    ):
        fail("execution plan runner-ready dispatch policy must include semantic_invariants")
    if contains_machine_local_path(plan):
        fail("execution plan contains a machine-local absolute path")

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
    coverage_policy = propagation.get("coverage_policy")
    if not isinstance(coverage_policy, dict):
        fail("dev_architecture_propagation.coverage_policy must be an object")
    if coverage_policy.get("source_ref") != "docs/dev/dev_guides.md#coverage-gates":
        fail("dev_architecture_propagation.coverage_policy.source_ref must point at docs/dev/dev_guides.md#coverage-gates")
    if coverage_policy.get("required") is not True:
        fail("dev_architecture_propagation.coverage_policy.required must be true")
    if not has_nonempty_list(coverage_policy.get("command_refs")):
        fail("dev_architecture_propagation.coverage_policy.command_refs must be non-empty")
    validate_coverage_policy_refs_for_root(ROOT, coverage_policy, "dev_architecture_propagation.coverage_policy")
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
    task_supervision_policy = plan.get("task_supervision_policy")
    if not isinstance(task_supervision_policy, dict):
        fail("execution plan missing task_supervision_policy")
    if task_supervision_policy.get("skill_ref") != "factory://skills/task-supervisor":
        fail("task_supervision_policy.skill_ref must be factory://skills/task-supervisor")
    if task_supervision_policy.get("evidence_artifact_type") != "task_supervisor_report":
        fail("task_supervision_policy.evidence_artifact_type must be task_supervisor_report")
    if task_supervision_policy.get("evidence_path") != ".factory/artifacts/task-supervisor-runs/<mission>/<timestamp>.json":
        fail("task_supervision_policy.evidence_path must point at task-supervisor-runs")
    return baseline_task_id


def validate_validation_contract(contract: dict[str, Any]) -> None:
    validate_no_legacy_provider_fields(contract, "validation-contract.json")
    validate_factory_compatibility(contract.get("factory_compatibility"), "validation-contract.json.factory_compatibility")
    validate_runtime_pins(contract.get("runtime_pins"), "validation-contract.json.runtime_pins")
    if contract.get("acceptance_ledger_ref") != ACCEPTANCE_LEDGER_REF:
        fail("validation-contract.json must cite acceptance-ledger.json")
    if contract.get("acceptance_item_count") != len(REQUIRED_ACCEPTANCE_ITEM_IDS):
        fail("validation-contract.json acceptance_item_count must match acceptance-ledger item count")
    if not has_nonempty_list(contract.get("acceptance_criteria")):
        fail("validation-contract.json must include itemized acceptance_criteria")
    validate_mvp_version_slice_coverage(
        contract.get("mvp_required_version_slices"),
        "validation-contract.json.mvp_required_version_slices",
    )
    validate_delivery_slice_coverage(
        contract.get("delivery_slices"),
        "validation-contract.json.delivery_slices",
    )
    validate_mvp_eval_provider_adapters(
        contract.get("mvp_eval_provider_adapters"),
        "validation-contract.json.mvp_eval_provider_adapters",
    )
    validate_plan_drift_policy(contract.get("plan_drift_policy"), "validation-contract.json.plan_drift_policy")
    validate_coverage_policy_refs_for_root(ROOT, contract.get("coverage_policy"), "validation-contract.json.coverage_policy")
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
    missing_execution_fields = [field for field in REQUIRED_TASK_FIELDS if field not in required_task_fields]
    if missing_execution_fields:
        fail(
            "validation-contract.json planning_skill_alignment.required_task_fields missing "
            f"execution-compiler fields: {missing_execution_fields}"
        )
    factoryd_requirements = contract.get("factoryd_runtime_requirements")
    if not isinstance(factoryd_requirements, dict):
        fail("validation-contract.json missing factoryd_runtime_requirements")
    missing_runner_ready = [
        field
        for field in REQUIRED_RUNNER_READY_FIELDS
        if field not in factoryd_requirements.get("runner_ready_fields", [])
    ]
    if missing_runner_ready:
        fail(f"validation-contract.json.factoryd_runtime_requirements.runner_ready_fields missing {missing_runner_ready}")
    validate_factoryd_runtime(
        factoryd_requirements.get("runtime"),
        "validation-contract.json.factoryd_runtime_requirements.runtime",
    )
    if contains_machine_local_path(contract):
        fail("validation-contract.json contains a machine-local absolute path")
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
