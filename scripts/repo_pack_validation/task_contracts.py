"""Validate Lumyn v3 task ownership and fail-closed execution boundaries."""

from __future__ import annotations

import hashlib
import json
from typing import Any

from repo_pack_validation.authority import manual_preflight_scope_digest


TASK_DEPENDENCIES = {
    "M0": [],
    "M1": ["M0"],
    "M2": ["M0"],
    "M2.5": [],
    "M3": ["M1", "M2"],
    "M4": ["M3"],
    "M5": ["M2", "M2.5", "M4"],
    "M6": ["M5"],
    "M7": ["M6"],
    "M8": ["M7"],
    "M9": ["M7"],
    "M10": ["M2.5", "M7", "M9"],
}

PRIMARY_ACCEPTANCE = {
    "M0": {f"BASE-{number:03d}" for number in range(1, 6)},
    "M2": {f"TRUST-{number:03d}" for number in range(1, 5)},
    "M2.5": {f"DISC-{number:03d}" for number in range(1, 4)},
    "M3": {f"PACK-{number:03d}" for number in range(1, 5)},
    "M4": {f"IMPACT-{number:03d}" for number in range(1, 6)},
    "M5": {f"PLAN-{number:03d}" for number in range(1, 4)},
    "M6": {f"AGENT-{number:03d}" for number in range(1, 8)},
    "M7": {f"VER-{number:03d}" for number in range(1, 7)},
    "M9": {f"EXP-{number:03d}" for number in range(1, 5)},
    "M10": {f"PILOT-{number:03d}" for number in range(1, 9)},
}

CANONICAL_WORKER_ORDER = [
    "task-executor",
    "validation-gate",
    "code-review",
    "holdout-evaluator",
    "trace-grader",
    "evidence-attestor",
    "commit-push",
    "post-merge-monitor",
]

MODEL_FIELDS = {
    "model_provider",
    "model_endpoint",
    "model_version",
    "model_parameters",
    "system_policy_digest",
    "prompt_digest",
    "tool_definition_digest",
    "context_selection",
    "request_disclosure_classes",
    "provider_retention_and_region",
    "credential_environment",
    "credential_scopes",
    "network_endpoint",
    "network_operations",
    "read_paths",
    "write_paths",
    "file_budget",
    "line_budget",
    "diff_budget",
    "turn_budget",
    "token_budget",
    "time_budget",
    "retry_budget",
    "concurrency_budget",
    "attempt_budget",
    "cost_budget",
    "request_response_tool_usage_patch_provenance",
}

UNTRUSTED_AGENT_INPUTS = {
    "provider_evidence",
    "repository_source",
    "repository_comments",
    "tests",
    "retrieved_context",
    "tool_output",
    "model_output",
}

CARRY_FORWARD_PROOF = {
    "BASE-001": (
        {"T1", "T3", "M0"},
        {
            ".factory/artifacts/task-runs/T1/validation-report.json",
            ".factory/artifacts/task-runs/T1/work-proof-marker.json",
            ".factory/artifacts/task-runs/T3/validation-report.json",
            ".factory/artifacts/task-runs/T3/work-proof-marker.json",
        },
    ),
    "BASE-002": (
        {"T2", "T2.7", "M0"},
        {
            ".factory/artifacts/task-runs/T2/validation-report.json",
            ".factory/artifacts/task-runs/T2.7/validation-report.json",
        },
    ),
    "BASE-004": (
        {"T3", "M0"},
        {
            ".factory/artifacts/task-runs/T3/validation-report.json",
            ".factory/artifacts/task-runs/T3/work-proof-marker.json",
        },
    ),
}


def _require(condition: bool, message: str) -> None:
    if not condition:
        raise AssertionError(message)


def validate_carry_forward_proof(
    ledger_by_id: dict[str, dict[str, Any]],
) -> None:
    """Require evidence whose task semantics prove each retained v3 claim."""

    for item_id, (task_refs, evidence_refs) in CARRY_FORWARD_PROOF.items():
        item = ledger_by_id[item_id]
        _require(
            set(item.get("task_refs", [])) == task_refs,
            f"{item_id} carry-forward task refs differ from its proof set",
        )
        _require(
            set(item.get("evidence_refs", [])) == evidence_refs,
            f"{item_id} carry-forward evidence is semantically incomplete",
        )


def _policy_digest(policy: dict[str, Any]) -> str:
    canonical = json.dumps(
        {key: value for key, value in policy.items() if key != "policy_digest"},
        sort_keys=True,
        separators=(",", ":"),
    ).encode()
    return f"sha256:{hashlib.sha256(canonical).hexdigest()}"


def _validate_paused_runtime(task: dict[str, Any]) -> None:
    task_id = str(task["task_id"])
    _require(
        task.get("dispatch_status") == "paused_factory_profile_and_runtime_unqualified",
        f"{task_id} must remain dispatch-paused",
    )
    runtime = task.get("factoryd_runtime")
    _require(isinstance(runtime, dict), f"{task_id}.factoryd_runtime is required")
    _require(runtime.get("worker_type") == "codex_cli", f"{task_id} worker type drifted")
    _require(runtime.get("dispatch_enabled") is False, f"{task_id} dispatch must be disabled")
    _require(runtime.get("fail_closed") is True, f"{task_id} runtime must fail closed")
    _require(
        runtime.get("profile_compatibility_status")
        == "blocked_v3_profile_update_required",
        f"{task_id} must record the incompatible v2 Factory profile",
    )
    _require(
        runtime.get("runtime_qualification_status")
        == "blocked_factoryd_v3_qualification_required",
        f"{task_id} must record missing factoryd v3 qualification",
    )
    control = runtime.get("runtime_control")
    _require(isinstance(control, dict), f"{task_id}.runtime_control is required")
    _require(control.get("mission_paused") is True, f"{task_id} mission must be paused")
    _require(
        control.get("launch_request", {}).get("expected_decision") == "deny",
        f"{task_id} launch request must be denied",
    )
    _require(
        control.get("conflict_behavior") == "fail_closed",
        f"{task_id} runtime conflict behavior must fail closed",
    )
    grants = runtime.get("capability_grants")
    _require(isinstance(grants, list), f"{task_id} capability grants must be a list")
    _require(
        all(grant.get("approved") is False for grant in grants),
        f"{task_id} planning grants must remain unapproved",
    )
    required_capabilities = set(task.get("requires_capabilities", []))
    _require(
        {grant.get("capability") for grant in grants} == required_capabilities,
        f"{task_id} seed grants must match declared Factory capabilities",
    )


def _validate_worker_chain(task: dict[str, Any]) -> None:
    task_id = str(task["task_id"])
    chain = task.get("required_worker_chain")
    _require(isinstance(chain, list), f"{task_id} worker chain must be a list")
    ordered = [worker for worker in CANONICAL_WORKER_ORDER if worker in chain]
    _require(chain == ordered, f"{task_id} worker chain is not canonical")
    _require(
        chain[:2] == ["task-executor", "validation-gate"],
        f"{task_id} must begin with implementation and validation",
    )
    _require(
        chain[-2:] == ["commit-push", "post-merge-monitor"],
        f"{task_id} must preserve shipping and post-merge controls",
    )
    _require("ship-pr" not in chain, f"{task_id} uses the deprecated ship-pr alias")
    gates = task.get("lifecycle_gates")
    _require(isinstance(gates, dict), f"{task_id} lifecycle gates are required")
    for field in (
        "local_validation_required",
        "ci_required",
        "codex_review_required",
        "commit_push_required",
        "post_merge_monitor_required",
        "pr_lifecycle_report_required",
    ):
        _require(gates.get(field) is True, f"{task_id}.{field} must be true")
    for worker, gate in (
        ("code-review", "code_review_required"),
        ("holdout-evaluator", "holdout_evaluation_required"),
        ("trace-grader", "trace_grading_required"),
        ("evidence-attestor", "evidence_attestation_required"),
    ):
        if gates.get(gate) is True:
            _require(worker in chain, f"{task_id} must include {worker}")


def _validate_holdout_contracts(tasks: dict[str, dict[str, Any]]) -> None:
    provision = tasks["M1"].get("holdout_policy")
    _require(isinstance(provision, dict), "M1 holdout policy is required")
    _require(provision.get("mode") == "provision", "M1 must provision the holdout")
    _require(
        provision.get("task_executor_access") == "forbidden",
        "M1 holdout answers must be hidden from task-executor",
    )
    _require(
        provision.get("comparison_baseline")
        == "same_snapshot_evidence_commands_role_and_time_budget_generic_agent",
        "M1 must freeze a fair generic-agent baseline",
    )
    _require(
        provision.get("policy_digest") == _policy_digest(provision),
        "M1 holdout policy digest drifted",
    )
    prohibited = set(provision.get("prohibited_committed_fields", []))
    _require(
        {
            "inputs",
            "answer_key",
            "expected_patches",
            "expected_labels",
            "raw_traces",
            "repository_url",
            "plaintext_content_digest",
        }.issubset(prohibited),
        "M1 holdout policy exposes answer or resolving material",
    )
    for task_id in ("M4", "M6", "M7"):
        policy = tasks[task_id].get("holdout_policy")
        _require(
            isinstance(policy, dict)
            and policy.get("mode") == "evaluate"
            and policy.get("provisioning_result_ref")
            == ".factory/artifacts/lifecycle-evidence/M1/holdout-result.json",
            f"{task_id} must evaluate the independently frozen holdout",
        )
        _require(
            policy.get("policy_digest") == _policy_digest(policy),
            f"{task_id} holdout policy digest drifted",
        )
    for task_id, task in tasks.items():
        if task_id not in {"M1", "M4", "M6", "M7"}:
            _require("holdout_policy" not in task, f"{task_id} must not access holdouts")


def _validate_preflight(task: dict[str, Any]) -> None:
    preflight = task.get("manual_external_evidence_preflight")
    _require(isinstance(preflight, dict), "M2.5 preflight is required")
    _require(preflight.get("required") is True, "M2.5 preflight must be required")
    _require(
        preflight.get("product_runtime_authority") is False,
        "M2.5 preflight is not product runtime authority",
    )
    _require(
        preflight.get("failure_behavior") == "block_external_evidence_collection",
        "M2.5 preflight must fail closed",
    )
    _require(
        preflight.get("factory_approval_capability") == "approval",
        "M2.5 preflight must use the canonical Factory approval capability",
    )
    _require(
        isinstance(preflight.get("allowed_private_fields"), list)
        and bool(preflight["allowed_private_fields"]),
        "M2.5 preflight must enumerate allowed private fields",
    )
    _require(
        isinstance(preflight.get("public_fields"), list)
        and bool(preflight["public_fields"]),
        "M2.5 preflight must enumerate source-safe public fields",
    )
    storage = str(preflight.get("approved_private_storage_boundary", "")).lower()
    _require(
        "external" in storage and "outside" in storage and "repository" in storage,
        "M2.5 private storage must be explicitly external to the repository",
    )
    for field in (
        "participant_consent_required",
        "retention_ttl_and_expiry_required",
        "deletion_on_revocation_required",
        "deletion_receipt_and_orphan_owner_required",
        "public_disclosure_requires_separate_consent",
    ):
        _require(
            preflight.get(field) is True,
            f"M2.5 preflight must require {field}",
        )
    prohibited = " ".join(preflight.get("prohibited_actions", [])).lower()
    for token in ("provider code access", "production access", "fabricated demand"):
        _require(token in prohibited, f"M2.5 preflight must prohibit {token}")
    _require(
        str(preflight.get("approval_evidence_ref", "")).startswith("pending:"),
        "M2.5 preflight must name a pending private approval-evidence reference",
    )
    private_scope = " ".join(preflight.get("allowed_private_fields", [])).lower()
    for token in (
        "provider payment",
        "provider operator",
        "eligible consumer",
        "privacy and model-data protocol",
        "generic-agent baseline",
        "campaign cogs",
        "lumyn operator hours",
        "consumer maintainer time",
        "absolute judgment deadline",
    ):
        _require(token in private_scope, f"M2.5 preflight must bind {token}")
    digest = preflight.get("approval_scope_digest")
    _require(
        digest == manual_preflight_scope_digest(preflight),
        "M2.5 preflight scope digest drifted",
    )


def _validate_model_contract(task: dict[str, Any]) -> None:
    task_id = str(task["task_id"])
    contract = task.get("bounded_agent_contract")
    _require(isinstance(contract, dict), f"{task_id} bounded-agent contract is required")
    _require(
        set(contract.get("exact_fields", [])) == MODEL_FIELDS,
        f"{task_id} model control fields are incomplete",
    )
    _require(
        set(contract.get("untrusted_inputs", [])) == UNTRUSTED_AGENT_INPUTS,
        f"{task_id} untrusted agent inputs are incomplete",
    )
    for field in (
        "cannot_widen_scope",
        "cannot_self_approve",
        "cannot_self_verify",
        "cannot_push_or_open_pr",
        "cannot_merge",
        "raw_prompt_response_persistence_default",
    ):
        expected = False if field == "raw_prompt_response_persistence_default" else True
        _require(
            contract.get(field) is expected,
            f"{task_id} bounded-agent {field} posture drifted",
        )


def _validate_verification(task: dict[str, Any]) -> None:
    contract = task.get("deterministic_verification_contract")
    _require(isinstance(contract, dict), "M7 deterministic verification contract is required")
    for field in (
        "independent_from_generation",
        "same_ladder_for_generation_modes",
        "preexisting_baseline_required",
        "exact_candidate_head_required",
        "causal_workflow_execution_required",
        "failed_or_stale_evidence_blocks_verified",
        "model_self_verification_forbidden",
    ):
        _require(contract.get(field) is True, f"M7 verification {field} must be true")
    _require(
        contract.get("verify_command") == "lumyn verify"
        and contract.get("verify_mutates_candidate") is False,
        "M7 verify must be explicitly non-mutating",
    )
    _require(
        contract.get("repair_is_separate_action") is True
        and contract.get("repair_produces_new_candidate") is True,
        "M7 repair must be a separate bounded action producing a new candidate",
    )
    repair = contract.get("repair_authorization")
    _require(isinstance(repair, dict), "M7 repair authorization contract is required")
    for field in (
        "failed_candidate_and_evidence",
        "exact_repair_intent",
        "remaining_permissions",
        "remaining_model_data_permissions",
        "remaining_time_budget",
        "remaining_token_budget",
        "remaining_cost_budget",
        "remaining_attempt_budget",
        "remaining_file_budget",
        "remaining_diff_budget",
    ):
        _require(
            repair.get(field) is True,
            f"M7 repair authorization must bind {field}",
        )
    _require(
        repair.get("scope_expansion_allowed") is False
        and repair.get("prior_verification_evidence_invalidated") is True
        and repair.get("fresh_full_verify_required") is True,
        "M7 repair must preserve scope, invalidate prior proof, and reverify fully",
    )
    _require(
        set(contract.get("candidate_modes", []))
        == {"deterministic", "agent_assisted", "imported_manual"},
        "M7 must run the same verification ladder for deterministic, agent-assisted, and imported manual candidates",
    )
    labels = set(contract.get("canonical_labels", []))
    _require(
        {
            "static_verified",
            "repo_verified",
            "workflow_contract_replay_passed",
            "workflow_verified_replay",
            "workflow_verified_mock",
            "workflow_verified_sandbox",
        }.issubset(labels),
        "M7 verification labels are incomplete",
    )


def _validate_export(task: dict[str, Any]) -> None:
    contract = task.get("delivery_contract")
    _require(isinstance(contract, dict), "M9 delivery contract is required")
    _require(
        contract.get("baseline_forms")
        == ["evidence_bundle", "patch", "optional_local_branch", "pr_ready_bundle"],
        "M9 must preserve the multi-form local export baseline",
    )
    _require(
        contract.get("optional_remote_forms") == ["remote_branch", "draft_pr"],
        "M9 optional remote forms drifted",
    )
    _require(
        contract.get("remote_branch_and_pr_separate_grants") is True,
        "M9 must separate branch and PR grants",
    )
    _require(
        contract.get("default_branch_write") is False
        and contract.get("auto_merge") is False,
        "M9 must forbid default-branch write and auto-merge",
    )
    _require(
        contract.get("exp_003_applicability")
        == "binary_frozen_before_execution",
        "M9 EXP-003 applicability must be frozen and binary",
    )
    automated = " ".join(contract.get("automated_path_requires", [])).lower()
    manual = " ".join(contract.get("not_applicable_path_requires", [])).lower()
    for token in (
        "dated frozen campaign protocol",
        "short-lived least-privilege token",
        "non-default branch",
        "draft pr",
        "idempotency evidence",
        "no auto-merge evidence",
    ):
        _require(token in automated, f"M9 automated EXP-003 evidence missing {token}")
    for token in (
        "dated frozen campaign protocol",
        "manual-only delivery",
        "patch, local branch, or pr-bundle evidence",
        "no automated-delivery claim",
    ):
        _require(token in manual, f"M9 manual EXP-003 evidence missing {token}")
    _require(
        contract.get("missing_protocol_behavior") == "EXP-003 remains open",
        "M9 cannot infer EXP-003 applicability without protocol evidence",
    )
    outcome = task.get("outcome_record_contract")
    _require(isinstance(outcome, dict), "M9 outcome record contract is required")
    _require(
        outcome.get("command") == "lumyn outcome record"
        and outcome.get("append_only") is True
        and outcome.get("exact_candidate_binding_required") is True
        and outcome.get("verification_evidence_digest_required") is True,
        "M9 outcome recording must be append-only and exact-candidate bound",
    )
    _require(
        {"consumer_accepted", "remediation"}.issubset(
            set(outcome.get("required_evidence_classes", []))
        ),
        "M9 outcome recording must preserve consumer acceptance and remediation evidence",
    )
    _require(
        {"accepted", "merged", "closed", "corrected", "reverted"}.issubset(
            set(outcome.get("durable_outcomes", []))
        ),
        "M9 outcome recording must preserve acceptance, merge, closure, correction, and reversion",
    )
    _require(
        "internal/outcome/" in task.get("allowed_paths", []),
        "M9 must own the future internal/outcome path",
    )


def _validate_campaign(task: dict[str, Any]) -> None:
    contract = task.get("paid_campaign_contract")
    _require(isinstance(contract, dict), "M10 paid campaign contract is required")
    for field in (
        "provider_prepayment_minimum_usd",
        "frozen_eligible_cohort_minimum",
        "minimum_valid_scans",
        "minimum_reviewable_outcomes",
        "minimum_accepted_or_merged_outcomes",
        "generic_agent_baseline",
        "material_advantage_threshold",
        "campaign_cogs",
        "lumyn_operator_hours",
        "consumer_maintainer_time",
        "absolute_judgment_deadline",
    ):
        _require(field in contract, f"M10 campaign contract missing {field}")
    _require(
        contract.get("generic_agent_baseline")
        == "same_snapshot_evidence_commands_role_and_time_budget",
        "M10 generic-agent baseline must be fair",
    )
    _require(
        contract.get("provider_prepayment_minimum_usd") == 25000
        and contract.get("payment_posture")
        == "cleared_non_refundable_prepaid_funds",
        "M10 must require cleared non-refundable prepaid funds of at least $25,000",
    )
    _require(
        contract.get("frozen_eligible_cohort_minimum") == 5
        and contract.get("eligible_consumer_units") == 5
        and contract.get("distinct_consumer_organizations") == 5,
        "M10 must bind five Eligible Consumer Units across five distinct organizations",
    )
    _require(
        contract.get("minimum_valid_scans") == 3
        and contract.get("minimum_reviewable_outcomes") == 2
        and contract.get("minimum_accepted_or_merged_outcomes") == 1,
        "M10 technical campaign minima must remain 3 scans, 2 outcomes, and 1 accepted or merged outcome",
    )
    advantage = str(contract.get("material_advantage_threshold", "")).lower()
    _require(
        "30 percent lower median" in advantage
        and "no worse" in advantage
        and "correction" in advantage
        and "revert" in advantage
        and "false-verification" in advantage
        and "frozen" in advantage,
        "M10 maintainer advantage must preserve the 30 percent and quality guardrail or an equally material frozen alternative",
    )
    _require(
        contract.get("material_provider_outcome_metric")
        and contract.get("material_provider_outcome_threshold")
        and contract.get("material_provider_outcome_pass_required") is True,
        "M10 needs a separate material provider-outcome metric, threshold, and pass",
    )
    for field in (
        "material_provider_outcome_source",
        "material_provider_outcome_denominator",
        "material_provider_outcome_comparator",
        "material_provider_outcome_economic_buyer_approval",
    ):
        _require(contract.get(field), f"M10 provider outcome missing {field}")
    _require(
        "before" in str(contract.get("absolute_judgment_deadline", "")).lower()
        and "invitation" in str(contract.get("absolute_judgment_deadline", "")).lower(),
        "M10 absolute judgment deadline must be frozen before invitations",
    )
    _require(
        contract.get("provider_code_access") is False,
        "M10 must not grant the provider code access",
    )


def validate_migration_task_contracts(
    tasks: dict[str, dict[str, Any]],
) -> None:
    """Validate task-specific ownership not expressible in generic Factory schemas."""

    _require(
        set(tasks) == set(TASK_DEPENDENCIES),
        "active task set differs from the M0-M10 authored plan",
    )
    for task_id, expected in TASK_DEPENDENCIES.items():
        _require(
            tasks[task_id].get("blocked_by") == expected,
            f"{task_id} dependency graph differs from docs/product/plan.md",
        )
    for task_id, required in PRIMARY_ACCEPTANCE.items():
        _require(
            required.issubset(set(tasks[task_id].get("acceptance_item_ids", []))),
            f"{task_id} primary acceptance ownership drifted",
        )
    _require(
        "rebaseline" not in str(tasks["M0"].get("objective", "")).lower(),
        "M0 is runtime foundation work and must not own the planning rebaseline",
    )
    m0_paths = set(tasks["M0"].get("allowed_paths", []))
    _require(
        m0_paths.issubset(
            {
                "cmd/lumyn/",
                "internal/result/",
                "internal/exitcode/",
                "schemas/",
                "tests/",
                "docs/",
                "CHANGELOG.md",
            }
        ),
        "M0 allowed paths must remain limited to runtime foundations",
    )
    _require(
        not {
            ".factory/artifacts/prd-to-plan/lumyn-migration-mvp/",
            "scripts/validate_repo_pack.py",
            "scripts/repo_pack_validation/",
        }.intersection(m0_paths),
        "M0 must not include planning-control compilation paths",
    )
    _require(
        {"PACK-001", "IMPACT-005", "AGENT-001", "VER-006", "PILOT-005"}.issubset(
            set(tasks["M1"].get("acceptance_item_ids", []))
        ),
        "M1 benchmark must cover pack, impact, hybrid agent, verification, and pilot comparison",
    )
    _require(
        {"VER-003", "VER-004", "VER-005"}.issubset(
            set(tasks["M8"].get("acceptance_item_ids", []))
        ),
        "M8 optional sandbox verification scope drifted",
    )
    _require(
        tasks["M5"].get("gated_by_acceptance_items")
        == [
            {"acceptance_item_id": "DISC-001", "required_status": "implemented"},
            {"acceptance_item_id": "DISC-002", "required_status": "implemented"},
            {"acceptance_item_id": "DISC-003", "required_status": "implemented"},
        ],
        "M5 must wait for paid campaign qualification",
    )
    _validate_holdout_contracts(tasks)
    _validate_preflight(tasks["M2.5"])
    _require(
        tasks["M2.5"].get("blocked_by") == [],
        "M2.5 DISC-001 and DISC-002 must remain startable without a milestone dependency",
    )
    disc_gates = tasks["M2.5"].get("gated_acceptance_items")
    _require(
        isinstance(disc_gates, list) and len(disc_gates) == 1,
        "M2.5 must declare exactly one DISC-003 closure gate",
    )
    disc_gate = disc_gates[0]
    _require(
        disc_gate.get("acceptance_item_id") == "DISC-003"
        and disc_gate.get("required_milestone") == "M2"
        and disc_gate.get("required_status") == "implemented"
        and set(disc_gate.get("required_contracts", []))
        == {"privacy", "model_data", "authorization", "evidence"},
        "DISC-003 must cite approved M2 privacy, model-data, authorization, and evidence contracts",
    )
    for task_id in ("M6", "M7", "M10"):
        _validate_model_contract(tasks[task_id])
    manual = tasks["M6"].get("manual_candidate_contract")
    _require(
        isinstance(manual, dict)
        and manual.get("command") == "lumyn candidate import --manual"
        and manual.get("label") == "manual"
        and manual.get("counts_as_agent_or_deterministic_automation") is False,
        "M6 must support proof-honest manual candidate import",
    )
    _require(
        {
            "migration_pack_digest",
            "integration_graph_digest",
            "plan_digest",
            "repository_base",
            "route",
            "candidate_digest",
        }.issubset(set(manual.get("binds", []))),
        "M6 manual candidate import must bind the full candidate provenance",
    )
    _validate_verification(tasks["M7"])
    _validate_export(tasks["M9"])
    _validate_campaign(tasks["M10"])
    for task in tasks.values():
        _validate_paused_runtime(task)
        _validate_worker_chain(task)
