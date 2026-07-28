"""Validate Lumyn v3 task ownership and fail-closed execution boundaries."""

from __future__ import annotations

import hashlib
import json
from typing import Any

from repo_pack_validation.authority import manual_preflight_scope_digest
from repo_pack_validation.m1_closure import validate_m1_closure_evidence


TASK_DEPENDENCIES = {
    "M0": [],
    "M1": ["M0", "M2"],
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
    "M3": (
        {f"PACK-{number:03d}" for number in range(1, 5)}
        | {f"EVENT-{number:03d}" for number in range(1, 3)}
    ),
    "M4": {f"IMPACT-{number:03d}" for number in range(1, 6)},
    "M5": (
        {f"INSTALL-{number:03d}" for number in range(1, 3)}
        | {f"PLAN-{number:03d}" for number in range(1, 4)}
    ),
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
    "agent_execution_policy",
    "agent_runner_vendor",
    "agent_runner_adapter",
    "agent_runner_version",
    "agent_runner_executable_source",
    "agent_runner_executable_path",
    "agent_runner_executable_digest",
    "agent_runner_conformance_digest",
    "agent_runner_auth_mode",
    "agent_runner_entitlement_class",
    "clean_session_identity",
    "execution_funding_mode",
    "credential_owner",
    "usage_billing_owner",
    "native_configuration_digest",
    "agent_runner_credential_environment",
    "agent_runner_credential_scopes",
    "agent_runner_network_endpoint",
    "agent_runner_network_operations",
    "actual_model_route",
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

GENERIC_AGENT_CONTROL_FIELDS = {
    "repository_snapshot",
    "authoritative_migration_evidence",
    "agent_runner_adapter",
    "agent_runner_version",
    "agent_runner_executable_digest",
    "actual_model_provider",
    "actual_model_version",
    "auth_mode",
    "entitlement_class",
    "execution_funding_mode",
    "credential_owner",
    "usage_billing_owner",
    "context_access_ceiling",
    "tools",
    "commands",
    "engineer_role",
    "verification_commands_and_time",
    "attempt_budget",
    "token_budget",
    "time_budget",
    "cost_budget",
}
GENERIC_AGENT_TREATMENT = (
    "lumyn_orchestration_impact_routing_boundary_enforcement_"
    "independent_verification_delivery_and_status"
)

RUNNER_PREFLIGHT_FIELDS = {
    "agent_execution_policy",
    "agent_runner_adapter",
    "agent_runner_version",
    "agent_runner_executable_source",
    "agent_runner_executable_digest",
    "agent_runner_auth_mode",
    "agent_runner_entitlement_class",
    "actual_model_provider",
    "actual_model_endpoint",
    "actual_model_version",
    "execution_funding_mode",
    "credential_owner",
    "usage_billing_owner",
}

QUALIFYING_SAME_RUN_ROUTE_SEQUENCE = [
    "M4/impact_read_only",
    "M6/agent_assisted_candidate",
    "M7/verify",
    "M9/remote_branch_push",
    "M9/draft_pr_create",
    "M9/provider_status_transmit",
]

QUALIFYING_SAME_RUN_EVIDENCE_FIELDS = [
    "run_id",
    "provider_change_event",
    "consumer_installation_authorization",
    "organic_agent_assisted_plan_item",
    "candidate_head",
    "verification_evidence_digest",
    "remote_branch",
    "draft_pr",
    "provider_status_projection",
]

ISOLATION_BACKEND_IDENTITY_FIELDS = [
    "backend",
    "version",
    "configuration_digest",
    "qualification_digest",
    "host_platform",
]

ISOLATION_RESOURCE_QUOTA_FIELDS = [
    "cpu_time",
    "memory_bytes",
    "pids",
    "process_tree_depth",
    "disk_bytes",
    "open_files",
]

MANAGED_CREDENTIAL_CONTRACT = {
    "mode": "provider_sponsored_lumyn_managed",
    "approved_broker_issuer_bound": True,
    "audience_fields": [
        "consumer_installation",
        "provider_change_event",
        "migration_plan",
        "agent_attempt",
        "agent_runner",
        "model_route",
    ],
    "maximum_ttl_seconds": 3600,
    "broker_redemption_mode": "one_time_into_single_agent_attempt_session",
    "in_attempt_multiple_calls_allowed_within_hard_quotas": True,
    "replay_after_attempt_allowed": False,
    "cross_attempt_credential_reuse_allowed": False,
    "refresh_allowed": False,
    "hard_token_quota_required": True,
    "hard_cost_quota_required": True,
    "revocation_required": True,
    "post_run_usage_reconciliation_required": True,
    "vendor_native_or_budget_enforcing_proxy_required": True,
    "unsupported_vendor_behavior": "managed_mode_unavailable",
    "reusable_credential_persistence_allowed": False,
    "api_provider_credential_access": False,
}

RUNNER_HOST_ISOLATION_CONTRACT = {
    "backend_identity_fields": ISOLATION_BACKEND_IDENTITY_FIELDS,
    "backend_qualification_required": True,
    "hard_resource_quota_fields": ISOLATION_RESOURCE_QUOTA_FIELDS,
    "hard_resource_quotas_required": True,
    "explicit_read_only_and_writable_mounts_required": True,
    "host_home_mounted": False,
    "os_credential_access_allowed": False,
    "ambient_service_sockets_allowed": False,
    "unrelated_inherited_descriptors_allowed": False,
    "child_process_restrictions_inherited": True,
    "host_enforced_egress_required": True,
    "cleanup_evidence_required": [
        "process_tree_terminated",
        "workspace_removed",
        "credential_revoked",
        "mount_and_socket_absence",
    ],
    "unsupported_host_behavior": "block_before_agent_launch",
    "executable_extensions_mode": "plugins_mcp_servers_and_hooks_prohibited_for_mvp",
    "malicious_child_and_tool_negative_tests_required": True,
    "fork_bomb_and_resource_exhaustion_negative_tests_required": True,
}

REPOSITORY_COMMAND_ISOLATION_CONTRACT = {
    "backend_identity_fields": ISOLATION_BACKEND_IDENTITY_FIELDS,
    "backend_qualification_required": True,
    "hard_resource_quota_fields": ISOLATION_RESOURCE_QUOTA_FIELDS,
    "hard_resource_quotas_required": True,
    "exact_command_allowlist_required": True,
    "exact_working_directory_required": True,
    "explicit_read_only_and_writable_mounts_required": True,
    "neutral_home_and_temp_roots_required": True,
    "explicit_executable_roots_required": True,
    "sanitized_environment_required": True,
    "ambient_secrets_allowed": False,
    "host_home_mounted": False,
    "os_credential_access_allowed": False,
    "ambient_service_sockets_allowed": False,
    "unrelated_inherited_descriptors_allowed": False,
    "child_process_restrictions_inherited": True,
    "network_default": "disabled",
    "network_enablement_requires_exact_route_grant": True,
    "lifecycle_scripts_default": "disabled",
    "lifecycle_script_enablement_requires_separate_approval": True,
    "host_isolation_backend_required": True,
    "unsupported_host_behavior": "block_before_command_execution",
    "timeout_and_output_budgets_required": True,
    "ambient_credentials_allowed": False,
    "route_selected_credential_injection_only": True,
    "agent_runner_model_and_sandbox_credentials_absent": True,
    "pre_and_post_patch_results_separate": True,
    "cleanup_evidence_required": [
        "process_tree_terminated",
        "workspace_or_mount_cleanup",
        "socket_and_descriptor_absence",
    ],
    "fork_bomb_and_resource_exhaustion_negative_tests_required": True,
}

SANDBOX_ENTRYPOINT_ISOLATION_CONTRACT = {
    "backend_identity_fields": ISOLATION_BACKEND_IDENTITY_FIELDS,
    "backend_qualification_required": True,
    "hard_resource_quota_fields": ISOLATION_RESOURCE_QUOTA_FIELDS,
    "hard_resource_quotas_required": True,
    "exact_candidate_head_read_only_mount_required": True,
    "exact_entrypoint_and_working_directory_required": True,
    "neutral_home_and_temp_roots_required": True,
    "sanitized_environment_required": True,
    "credential_injection_mode": "sole_task_scoped_sandbox_credential",
    "unrelated_credentials_allowed": False,
    "agent_runner_model_and_github_credentials_absent": True,
    "egress_policy": "exact_sandbox_endpoint_and_operations_only",
    "ambient_service_sockets_allowed": False,
    "unrelated_inherited_descriptors_allowed": False,
    "child_process_restrictions_inherited": True,
    "timeout_and_output_budgets_required": True,
    "teardown_required": True,
    "cleanup_evidence_required": [
        "process_tree_terminated",
        "sandbox_credential_revoked",
        "temporary_workspace_removed",
        "sandbox_state_cleaned_or_orphan_recorded",
        "mount_socket_and_descriptor_absence",
    ],
    "orphan_evidence_required_when_cleanup_incomplete": True,
    "unsupported_host_behavior": "block_before_sandbox_entrypoint",
    "fork_bomb_and_resource_exhaustion_negative_tests_required": True,
}

IMPLEMENTED_PROOF = {
    "BASE-001": (
        {"T1", "T3", "M0"},
        {
            ".factory/artifacts/task-runs/T1/validation-report.json",
            ".factory/artifacts/task-runs/T1/work-proof-marker.json",
            ".factory/artifacts/task-runs/T3/validation-report.json",
            ".factory/artifacts/task-runs/T3/work-proof-marker.json",
            ".factory/artifacts/task-runs/M0/validation-report.json",
            ".factory/artifacts/task-runs/M0/proof-of-behavior-scorecard.json",
            ".factory/artifacts/pr-lifecycle/lumyn-v3-m0/implementation/work-proof-marker.json",
            ".factory/artifacts/pr-lifecycle/lumyn-v3-m0/scope-closure-report.json",
        },
    ),
    "BASE-002": (
        {"T2", "T2.7", "M0"},
        {
            ".factory/artifacts/task-runs/T2/validation-report.json",
            ".factory/artifacts/task-runs/T2.7/validation-report.json",
            ".factory/artifacts/task-runs/M0/validation-report.json",
            ".factory/artifacts/lifecycle-evidence/M0/review-report.json",
            ".factory/artifacts/pr-lifecycle/lumyn-v3-m0/implementation/work-proof-marker.json",
            ".factory/artifacts/pr-lifecycle/lumyn-v3-m0/scope-closure-report.json",
        },
    ),
    "BASE-003": (
        {"M0"},
        {
            ".factory/artifacts/task-runs/M0/validation-report.json",
            ".factory/artifacts/lifecycle-evidence/M0/review-report.json",
            ".factory/artifacts/pr-lifecycle/lumyn-v3-m0/post-merge-report.json",
            ".factory/artifacts/pr-lifecycle/lumyn-v3-m0/process-exception-handoff.json",
            ".factory/artifacts/pr-lifecycle/lumyn-v3-m0/delivery-debt-record.json",
            ".factory/artifacts/pr-lifecycle/lumyn-v3-m0/scope-closure-report.json",
        },
    ),
    "BASE-004": (
        {"T3", "M0"},
        {
            ".factory/artifacts/task-runs/T3/validation-report.json",
            ".factory/artifacts/task-runs/T3/work-proof-marker.json",
            ".factory/artifacts/task-runs/M0/validation-report.json",
            ".factory/artifacts/pr-lifecycle/lumyn-v3-m0/implementation/work-proof-marker.json",
            ".factory/artifacts/pr-lifecycle/lumyn-v3-m0/scope-closure-report.json",
        },
    ),
    "BASE-005": (
        {"M0"},
        {
            ".factory/artifacts/task-runs/M0/validation-report.json",
            ".factory/artifacts/task-runs/M0/proof-of-behavior-scorecard.json",
            ".factory/artifacts/lifecycle-evidence/M0/review-report.json",
            ".factory/artifacts/pr-lifecycle/lumyn-v3-m0/implementation/work-proof-marker.json",
            ".factory/artifacts/pr-lifecycle/lumyn-v3-m0/scope-closure-report.json",
        },
    ),
    "TRUST-001": (
        {"M2"},
        {
            ".factory/artifacts/task-runs/M2/validation-report.json",
            ".factory/artifacts/task-runs/M2/validation-run-summary.json",
            ".factory/artifacts/task-runs/M2/proof-of-behavior-scorecard.json",
            ".factory/artifacts/lifecycle-evidence/M2/review-report.json",
            ".factory/artifacts/pr-lifecycle/lumyn-v3-m2/implementation/work-proof-marker.json",
            ".factory/artifacts/pr-lifecycle/lumyn-v3-m2/post-merge-report.json",
            ".factory/artifacts/pr-lifecycle/lumyn-v3-m2/scope-closure-report.json",
            ".factory/artifacts/pr-lifecycle/lumyn-v3-m2/delivery-debt-record.json",
        },
    ),
    "TRUST-002": (
        {"M2"},
        {
            ".factory/artifacts/task-runs/M2/validation-report.json",
            ".factory/artifacts/task-runs/M2/validation-run-summary.json",
            ".factory/artifacts/task-runs/M2/proof-of-behavior-scorecard.json",
            ".factory/artifacts/lifecycle-evidence/M2/review-report.json",
            ".factory/artifacts/pr-lifecycle/lumyn-v3-m2/implementation/work-proof-marker.json",
            ".factory/artifacts/pr-lifecycle/lumyn-v3-m2/post-merge-report.json",
            ".factory/artifacts/pr-lifecycle/lumyn-v3-m2/scope-closure-report.json",
        },
    ),
    "TRUST-003": (
        {"M2"},
        {
            ".factory/artifacts/task-runs/M2/validation-report.json",
            ".factory/artifacts/task-runs/M2/validation-run-summary.json",
            ".factory/artifacts/task-runs/M2/proof-of-behavior-scorecard.json",
            ".factory/artifacts/lifecycle-evidence/M2/review-report.json",
            ".factory/artifacts/pr-lifecycle/lumyn-v3-m2/implementation/work-proof-marker.json",
            ".factory/artifacts/pr-lifecycle/lumyn-v3-m2/post-merge-report.json",
            ".factory/artifacts/pr-lifecycle/lumyn-v3-m2/scope-closure-report.json",
        },
    ),
    "TRUST-004": (
        {"M2"},
        {
            ".factory/artifacts/task-runs/M2/validation-report.json",
            ".factory/artifacts/task-runs/M2/validation-run-summary.json",
            ".factory/artifacts/task-runs/M2/proof-of-behavior-scorecard.json",
            ".factory/artifacts/lifecycle-evidence/M2/review-report.json",
            ".factory/artifacts/pr-lifecycle/lumyn-v3-m2/implementation/work-proof-marker.json",
            ".factory/artifacts/pr-lifecycle/lumyn-v3-m2/post-merge-report.json",
            ".factory/artifacts/pr-lifecycle/lumyn-v3-m2/scope-closure-report.json",
        },
    ),
}

M2_VALIDATION_REPORT_REF = ".factory/artifacts/task-runs/M2/validation-report.json"
M2_VALIDATION_SUMMARY_REF = (
    ".factory/artifacts/task-runs/M2/validation-run-summary.json"
)
M2_SCORECARD_REF = ".factory/artifacts/task-runs/M2/proof-of-behavior-scorecard.json"
M2_REVIEW_REF = ".factory/artifacts/lifecycle-evidence/M2/review-report.json"
M2_IMPLEMENTATION_MARKER_REF = (
    ".factory/artifacts/pr-lifecycle/lumyn-v3-m2/implementation/"
    "work-proof-marker.json"
)
M2_PR_LIFECYCLE_REF = (
    ".factory/artifacts/pr-lifecycle/lumyn-v3-m2/pr-lifecycle-report.json"
)
M2_POST_MERGE_REF = (
    ".factory/artifacts/pr-lifecycle/lumyn-v3-m2/post-merge-report.json"
)
M2_SCOPE_CLOSURE_REF = (
    ".factory/artifacts/pr-lifecycle/lumyn-v3-m2/scope-closure-report.json"
)
M2_DEBT_REF = (
    ".factory/artifacts/pr-lifecycle/lumyn-v3-m2/delivery-debt-record.json"
)
M2_EXCEPTION_REF = (
    ".factory/artifacts/pr-lifecycle/lumyn-v3-m2/process-exception-handoff.json"
)
M2_RETAINED_BUNDLE_REF = (
    ".factory/artifacts/pr-lifecycle/lumyn-v3-m2/implementation/"
    "pr72-original-head.bundle"
)
M2_EXPECTED_VALIDATION_RUN = (
    "validation:2026-07-26T16:32:25Z:50773146792248c5907105189e125e8a"
)
M2_EXPECTED_CANDIDATE = (
    "sha256:8509fbe0980c17d6e6ad3d9cd2b8f04e122d9d6699cb7162b33492198e040324"
)
M2_EXPECTED_BASE = "7609e5c49c0776c1028c1aeb3e2e2ee942b613b6"
M2_EXPECTED_PR_HEAD = "9345f3392ec98eb0e10345fe7941fd9d1450e55b"
M2_EXPECTED_LANDED_HEAD = "f89bc82490ffb6df908df6f8572054ee051ed6c6"
M2_EXPECTED_TRUST_ITEMS = {f"TRUST-{number:03d}" for number in range(1, 5)}
M2_EXPECTED_SCOPE_ITEMS = M2_EXPECTED_TRUST_ITEMS | {
    "EVENT-001",
    "EVENT-002",
    "INSTALL-001",
    "INSTALL-002",
}
M2_CLOSURE_JSON_REFS = (
    M2_VALIDATION_REPORT_REF,
    M2_VALIDATION_SUMMARY_REF,
    M2_SCORECARD_REF,
    M2_REVIEW_REF,
    M2_IMPLEMENTATION_MARKER_REF,
    M2_PR_LIFECYCLE_REF,
    M2_POST_MERGE_REF,
    M2_SCOPE_CLOSURE_REF,
    M2_DEBT_REF,
    M2_EXCEPTION_REF,
)


def _require(condition: bool, message: str) -> None:
    if not condition:
        raise AssertionError(message)


def _validate_repository_command_isolation(
    task: dict[str, Any],
    task_id: str,
) -> None:
    _require(
        task.get("repository_command_isolation_contract")
        == REPOSITORY_COMMAND_ISOLATION_CONTRACT,
        f"{task_id} repository-command isolation contract drifted",
    )


def _validate_sandbox_entrypoint_isolation(
    task: dict[str, Any],
    task_id: str,
) -> None:
    _require(
        task.get("sandbox_entrypoint_isolation_contract")
        == SANDBOX_ENTRYPOINT_ISOLATION_CONTRACT,
        f"{task_id} sandbox-entrypoint isolation contract drifted",
    )


def _m2_evidence(
    evidence_payloads: dict[str, dict[str, Any]],
    ref: str,
) -> dict[str, Any]:
    payload = evidence_payloads.get(ref)
    _require(isinstance(payload, dict), f"M2 closure evidence is missing {ref}")
    return payload


def _require_m2_identity(payload: dict[str, Any], label: str) -> None:
    _require(payload.get("task_id") == "M2", f"{label} must bind task M2")
    _require(
        payload.get("work_item_id") == "lumyn-v3-m2",
        f"{label} must bind work item lumyn-v3-m2",
    )


def _validate_m2_closure_evidence(
    evidence_payloads: dict[str, dict[str, Any]],
) -> None:
    validation = _m2_evidence(evidence_payloads, M2_VALIDATION_REPORT_REF)
    summary = _m2_evidence(evidence_payloads, M2_VALIDATION_SUMMARY_REF)
    scorecard = _m2_evidence(evidence_payloads, M2_SCORECARD_REF)
    review = _m2_evidence(evidence_payloads, M2_REVIEW_REF)
    marker = _m2_evidence(evidence_payloads, M2_IMPLEMENTATION_MARKER_REF)
    lifecycle = _m2_evidence(evidence_payloads, M2_PR_LIFECYCLE_REF)
    post_merge = _m2_evidence(evidence_payloads, M2_POST_MERGE_REF)
    closure = _m2_evidence(evidence_payloads, M2_SCOPE_CLOSURE_REF)
    debt = _m2_evidence(evidence_payloads, M2_DEBT_REF)
    exception = _m2_evidence(evidence_payloads, M2_EXCEPTION_REF)

    for label, payload in (
        ("M2 validation report", validation),
        ("M2 validation summary", summary),
        ("M2 proof scorecard", scorecard),
        ("M2 review report", review),
        ("M2 PR lifecycle report", lifecycle),
        ("M2 post-merge report", post_merge),
        ("M2 scope-closure report", closure),
        ("M2 delivery-debt record", debt),
    ):
        _require_m2_identity(payload, label)

    for label, payload in (
        ("M2 validation report", validation),
        ("M2 validation summary", summary),
        ("M2 proof scorecard", scorecard),
    ):
        _require(
            payload.get("validation_run_id") == M2_EXPECTED_VALIDATION_RUN,
            f"{label} validation run binding drifted",
        )
        _require(
            payload.get("candidate_digest") == M2_EXPECTED_CANDIDATE,
            f"{label} candidate binding drifted",
        )

    _require(
        validation.get("result") == "pass"
        and validation.get("promotion_decision") == "ready_for_pr"
        and validation.get("checks")
        and all(check.get("status") == "pass" for check in validation["checks"]),
        "M2 validation report must be a passing ready-for-PR result",
    )
    binding = summary.get("candidate_binding", {})
    _require(
        summary.get("status") == "pass"
        and binding.get("base_git_sha") == M2_EXPECTED_BASE
        and binding.get("candidate_digest") == M2_EXPECTED_CANDIDATE
        and len(binding.get("changed_paths", [])) == 105,
        "M2 validation summary binding or passing status drifted",
    )
    levels = {
        level.get("level"): level.get("status")
        for level in scorecard.get("levels", [])
        if isinstance(level, dict)
    }
    _require(
        scorecard.get("overall_status") == "pass"
        and scorecard.get("required_proof_level") == "workflow_behavior"
        and levels.get("workflow_behavior") == "pass",
        "M2 proof scorecard is not passing at its required proof level",
    )

    current_work = review.get("current_work", {})
    _require(
        review.get("verdict") == "approved"
        and review.get("review_type") == "security"
        and current_work.get("validation_run_id") == M2_EXPECTED_VALIDATION_RUN
        and current_work.get("candidate_digest") == M2_EXPECTED_CANDIDATE
        and review.get("findings")
        and all(
            finding.get("status") == "resolved"
            for finding in review.get("findings", [])
        ),
        "M2 independent review is not approved, resolved, and current-work bound",
    )

    landed_proof = marker.get("landed_binding_proof", {})
    _require(
        marker.get("execution_status") == "pass"
        and marker.get("exit_code") == 0
        and marker.get("git_sha") == M2_EXPECTED_LANDED_HEAD
        and landed_proof.get("base_head") == M2_EXPECTED_BASE
        and landed_proof.get("original_pr_head") == M2_EXPECTED_PR_HEAD
        and landed_proof.get("landed_main_head") == M2_EXPECTED_LANDED_HEAD
        and landed_proof.get("validation_run_id") == M2_EXPECTED_VALIDATION_RUN
        and landed_proof.get("candidate_digest") == M2_EXPECTED_CANDIDATE
        and landed_proof.get("changed_path_count") == 105
        and landed_proof.get("retained_bundle_ref") == M2_RETAINED_BUNDLE_REF
        and str(landed_proof.get("retained_bundle_sha256", "")).startswith(
            "sha256:"
        ),
        "M2 landed work-proof marker is not passing or fully bound",
    )

    local_validation = lifecycle.get("local_validation", {})
    codex_review = lifecycle.get("codex_review", {})
    merge = lifecycle.get("merge", {})
    _require(
        lifecycle.get("status") == "complete"
        and local_validation.get("status") == "pass"
        and local_validation.get("validation_run_id") == M2_EXPECTED_VALIDATION_RUN
        and local_validation.get("candidate_digest") == M2_EXPECTED_CANDIDATE
        and lifecycle.get("code_review", {}).get("status") == "approved"
        and codex_review.get("configured_enabled") is True
        and codex_review.get("enabled") is False
        and codex_review.get("status") == "disabled"
        and codex_review.get("unresolved_actionable_comments") == 0
        and codex_review.get("late_review_after_merge") is False
        and merge.get("status") == "merged"
        and merge.get("merged_commit") == M2_EXPECTED_LANDED_HEAD
        and lifecycle.get("post_merge", {}).get("status") == "healthy"
        and any(
            gate.get("gate") == "codex_review"
            and gate.get("approved_exception_ref") == M2_EXCEPTION_REF
            for gate in lifecycle.get("skipped_gates", [])
        ),
        "M2 PR lifecycle report is incomplete or misstates its review exception",
    )

    _require(
        post_merge.get("status") == "healthy"
        and post_merge.get("merge_reference") == M2_EXPECTED_LANDED_HEAD
        and post_merge.get("merge_performed") is True
        and post_merge.get("default_branch_checked") is True
        and post_merge.get("checks_observed")
        and all(
            check.get("status") == "pass"
            for check in post_merge.get("checks_observed", [])
        )
        and post_merge.get("regressions_detected") == []
        and post_merge.get("late_codex_review", {}).get(
            "late_review_after_merge"
        )
        is False,
        "M2 post-merge report is not healthy on the landed commit",
    )

    closure_items = {
        item.get("intent_item_id"): item
        for item in closure.get("closure_items", [])
        if isinstance(item, dict)
    }
    _require(
        closure.get("overall_status") == "closed_with_debt"
        and closure.get("rerun_required") is False
        and closure.get("blocker_refs") == []
        and set(closure_items) == M2_EXPECTED_SCOPE_ITEMS
        and closure_items.get("TRUST-001", {}).get("classification")
        == "implemented_with_debt"
        and all(
            closure_items[item_id].get("classification") == "implemented"
            for item_id in M2_EXPECTED_TRUST_ITEMS - {"TRUST-001"}
        )
        and closure.get("delivery_debt_refs") == [M2_DEBT_REF],
        "M2 scope-closure report is incomplete or overstates debt-free closure",
    )
    _require(
        debt.get("debt_type") == "skipped_required_latest_head_codex_review"
        and debt.get("severity") == "high"
        and debt.get("follow_up_action")
        == "retain_as_non_reusable_process_exception_and_require_clean_m1_lifecycle",
        "M2 delivery debt does not preserve the skipped review gate",
    )
    _require(
        exception.get("artifact_type") == "handoff_record"
        and exception.get("work_item_id") == "lumyn-v3-m2"
        and exception.get("current_state") == "exception_approved"
        and exception.get("blockers") == []
        and exception.get("approval_card", {}).get("decision_status") == "approved"
        and exception.get("merge_authority_grant", {}).get("task_id") == "M2",
        "M2 process exception is missing, blocked, or not task scoped",
    )


def validate_implemented_proof(
    ledger_by_id: dict[str, dict[str, Any]],
    evidence_payloads: dict[str, dict[str, Any]],
) -> None:
    """Require exact evidence whose task semantics prove each completed item."""

    for item_id, (task_refs, evidence_refs) in IMPLEMENTED_PROOF.items():
        item = ledger_by_id[item_id]
        _require(
            set(item.get("task_refs", [])) == task_refs,
            f"{item_id} implemented task refs differ from its proof set",
        )
        _require(
            set(item.get("evidence_refs", [])) == evidence_refs,
            f"{item_id} implemented evidence is semantically incomplete",
        )
    validate_m1_closure_evidence(evidence_payloads)
    _validate_m2_closure_evidence(evidence_payloads)


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
        task.get("dispatch_status")
        == "attended_task_requires_explicit_approval_factoryd_paused",
        f"{task_id} must require attended approval while factoryd is paused",
    )
    runtime = task.get("factoryd_runtime")
    _require(isinstance(runtime, dict), f"{task_id}.factoryd_runtime is required")
    _require(runtime.get("worker_type") == "codex_cli", f"{task_id} worker type drifted")
    _require(runtime.get("dispatch_enabled") is False, f"{task_id} dispatch must be disabled")
    _require(runtime.get("fail_closed") is True, f"{task_id} runtime must fail closed")
    _require(
        runtime.get("profile_compatibility_status")
        == "aligned_lumyn_v3_1",
        f"{task_id} must bind the aligned Lumyn v3.1 profile",
    )
    _require(
        runtime.get("runtime_qualification_status")
        == "blocked_factoryd_v3_1_qualification_required",
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
        ("holdout-evaluator", "holdout_provisioning_required"),
        ("holdout-evaluator", "holdout_evaluation_required"),
        ("trace-grader", "trace_grading_required"),
        ("evidence-attestor", "evidence_attestation_required"),
    ):
        if gates.get(gate) is True:
            _require(worker in chain, f"{task_id} must include {worker}")
    _require(
        not (
            gates.get("holdout_provisioning_required") is True
            and gates.get("holdout_evaluation_required") is True
        ),
        f"{task_id} cannot provision and evaluate a holdout in one task packet",
    )


def _validate_holdout_contracts(tasks: dict[str, dict[str, Any]]) -> None:
    provision = tasks["M1"].get("holdout_policy")
    provisioning = tasks["M1"].get("holdout_provisioning_contract")
    _require(isinstance(provision, dict), "M1 holdout policy is required")
    _require(
        isinstance(provisioning, dict),
        "M1 holdout provisioning contract is required",
    )
    _require(provision.get("mode") == "provision", "M1 must provision the holdout")
    m1_gates = tasks["M1"].get("lifecycle_gates", {})
    _require(
        m1_gates.get("holdout_provisioning_required") is True
        and m1_gates.get("holdout_evaluation_required") is False,
        "M1 must provision, not evaluate, the hidden holdout",
    )
    _require(
        provisioning.get("task_executor_access") == "forbidden",
        "M1 holdout answers must be hidden from task-executor",
    )
    _require(
        provisioning.get("comparison_baseline")
        == "matched_agent_engine_and_budget_status_quo",
        "M1 must freeze a fair generic-agent baseline",
    )
    _require(
        set(provisioning.get("comparison_control_variables", []))
        == GENERIC_AGENT_CONTROL_FIELDS
        and provisioning.get("difference_under_test") == GENERIC_AGENT_TREATMENT
        and provisioning.get("unmatched_engine_comparison_is_causal") is False,
        "M1 baseline must isolate Lumyn from runner, model, auth, funding, and budget effects",
    )
    _require(
        provision.get("policy_digest") == _policy_digest(provision),
        "M1 holdout policy digest drifted",
    )
    prohibited = set(provisioning.get("prohibited_committed_fields", []))
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
        "agent execution policy",
        "runner family",
        "observable model route",
        "host-isolation posture",
        "auth and entitlement",
        "credential and usage-billing ownership",
        "generic-agent baseline",
        "campaign cogs",
        "lumyn operator hours",
        "consumer maintainer time",
        "absolute judgment deadline",
    ):
        _require(token in private_scope, f"M2.5 preflight must bind {token}")
    runner = task.get("consumer_runner_prequalification", {})
    _require(
        runner.get("minimum_agent_assisted_feasible_consumers") == 1
        and runner.get(
            "scan_or_deterministic_consumers_may_disable_agent_execution"
        ) is True
        and runner.get("feasible_runner_families") == ["codex", "claude_code"]
        and runner.get(
            "consumer_selected_or_managed_route_explicitly_consented"
        ) is True
        and set(runner.get("exact_fields", [])) == RUNNER_PREFLIGHT_FIELDS
        and runner.get("noninteractive_automation_required") is True
        and runner.get("opaque_model_route_allowed") is False
        and runner.get("secret_values_collected") is False
        and runner.get("host_isolation_feasibility_required") is True
        and runner.get(
            "managed_route_broker_feasibility_required_when_selected"
        ) is True
        and runner.get(
            "lumyn_adapter_conformance_required_at_prequalification"
        ) is False
        and runner.get(
            "lumyn_adapter_conformance_required_before_first_agent_run"
        ) is True
        and runner.get(
            "live_canary_required_before_first_agent_run"
        ) is True
        and runner.get(
            "plausible_organic_agent_item_hypothesis_required"
        ) is True
        and runner.get("hypothesis_examples")
        == [
            "wrapper_adaptation",
            "adapter_adaptation",
            "signature_or_type_adaptation",
            "related_test_repair",
        ]
        and runner.get("new_business_decision_required") is False
        and runner.get("confirmation_stage") == "M4_M5_impact_and_plan"
        and runner.get("prequalification_may_force_agent_route") is False
        and runner.get("frozen_before_invitations") is True,
        "M2.5 must prequalify one feasible, consented, noninteractive Agent Runner route and plausible organic agent-eligible item without forcing the route, premature conformance, or secrets",
    )
    digest = preflight.get("approval_scope_digest")
    _require(
        digest == manual_preflight_scope_digest(preflight),
        "M2.5 preflight scope digest drifted",
    )
    readiness = task.get("readiness_sprint_contract", {})
    _require(
        readiness.get("minimum_price_usd") == 7500
        and readiness.get("maximum_price_usd") == 15000
        and readiness.get("creditable_toward_campaign") is True
        and readiness.get("may_close_disc_001") is False
        and readiness.get("may_authorize_repository_work") is False
        and readiness.get(
            "credited_funds_count_only_after_signed_campaign_conversion"
        ) is True
        and readiness.get(
            "converted_funds_are_non_refundable_campaign_consideration"
        ) is True
        and readiness.get("total_cleared_campaign_funds_minimum_usd") == 25000
        and readiness.get("readiness_sprint_alone_closes_disc_001") is False
        and readiness.get("proof_boundary")
        == "paid_discovery_not_provider_to_consumer_delivery",
        "M2.5 readiness sprint must remain paid discovery, not campaign proof",
    )
    offer = task.get("campaign_offer_contract", {})
    _require(
        offer.get("included_eligible_consumer_units") == 5
        and offer.get("minimum_valid_scans") == 3
        and offer.get("minimum_tested_reviewable_outcomes") == 2
        and offer.get("scope_posture")
        == "bounded_campaign_attempt_not_completed_migration_guarantee"
        and offer.get("consumer_consent_dependency_disclosed") is True
        and offer.get("additional_repository_default_billable_event")
        == "independently_verified_lumyn_candidate_and_authorized_tested_draft_pr"
        and offer.get(
            "contracted_verified_local_bundle_may_be_billable"
        ) is True
        and offer.get(
            "local_bundle_counts_as_automated_delivery_or_mvp_proof"
        ) is False,
        "M2.5 campaign offer must align price, success floor, and billable event",
    )
    delivery = task.get("consumer_delivery_prequalification", {})
    _require(
        delivery.get(
            "minimum_installed_preauthorization_draft_pr_willing_consumers"
        ) == 1
        and delivery.get(
            "minimum_provider_status_transmit_willing_consumers"
        ) == 1
        and delivery.get(
            "same_intended_qualifying_consumer_required"
        ) is True
        and delivery.get(
            "exact_allowlisted_event_bound_projection_reviewed"
        ) is True
        and delivery.get("willingness_is_runtime_consent") is False
        and delivery.get(
            "runtime_provider_reporting_grant_still_required"
        ) is True,
        "M2.5 must prequalify draft-PR and provider-status willingness without treating it as consent",
    )


def _validate_update_delivery_contracts(
    m1: dict[str, Any], m2: dict[str, Any]
) -> None:
    skeleton = m1.get("walking_skeleton_contract")
    _require(isinstance(skeleton, dict), "M1 walking skeleton contract is required")
    _require(
        skeleton.get("input_mode") == "public_or_synthetic_only"
        and skeleton.get("external_write_mode") == "manual_pr_bundle_only_no_external_write"
        and skeleton.get("generic_agent_harness")
        == "separate_common_agent_runner_contract_conformance_with_deterministic_fake",
        "M1 walking skeleton must be public-fixture, agent-free, no-write, and separate from fake-adapter conformance",
    )
    _require(
        skeleton.get("stages")
        == [
            "provider_change_event",
            "consumer_installation",
            "repository_impact",
            "migration_plan",
            "candidate",
            "verification",
            "pr_bundle",
        ],
        "M1 walking skeleton stages drifted",
    )
    _require(
        {
            "EVENT-001", "EVENT-002", "INSTALL-001", "INSTALL-002",
            "DISC-001", "DISC-002", "EXP-003", "PILOT-002", "PILOT-003",
            "PILOT-004",
        } == set(skeleton.get("cannot_close_acceptance_items", [])),
        "M1 walking skeleton must not claim commercial, real-repo, or automated-PR proof",
    )
    _require(
        skeleton.get("live_agent_execution_allowed") is False
        and skeleton.get("live_agent_execution_deferred_to") == ["M6", "M10"]
        and skeleton.get("no_grant_fallback")
        == "agent_route_disabled_no_adapter_invocation",
        "M1 must remain offline and defer live Agent Runner execution",
    )
    _require(
        skeleton.get("future_launch_agent_runner_targets")
        == ["codex", "claude_code"]
        and skeleton.get("deferred_agent_runner_targets") == ["cursor"]
        and set(skeleton.get("offline_fake_contract_assertions", []))
        == {
            "clean_session",
            "executable_integrity",
            "auth_and_entitlement_shape",
            "neutral_home_and_config",
            "repository_path_shadowing_rejected",
            "silent_fallback_rejected",
        },
        "M1 must test only the offline common clean-session contract",
    )
    channel = m2.get("update_channel_contract")
    _require(isinstance(channel, dict), "M2 update channel contract is required")
    _require(
        channel.get("provider_change_contract_confirmation")
        == "once_per_exact_version"
        and channel.get("provider_event_non_executable") is True
        and channel.get("authenticated_origin_required_for_unattended_writes") is True,
        "M2 provider change contract and event posture drifted",
    )
    _require(
        {
            "event_id", "event_version", "issuer", "api_or_sdk",
            "contract_digest", "contract_location", "audience", "deadline",
            "severity", "sequence", "issued_at", "expires_at", "transport_origin",
            "signature_provenance", "supersession_or_withdrawal",
        }
        == set(channel.get("provider_event_fields", []))
        and {
            "duplicate", "replayed", "stale", "conflicting", "expired",
            "superseded", "withdrawn", "wrong_audience", "origin_mismatch",
            "signature_invalid", "unauthenticated",
        }
        == set(channel.get("rejected_event_states", [])),
        "M2 provider event fields or fail-closed states drifted",
    )
    transport = channel.get("provider_channel_transport", {})
    _require(
        transport.get("transport_id")
        == "pinned_provider_https_signed_manifest_v1"
        and transport.get("publisher") == "provider_operator"
        and transport.get("exact_https_url_pinned") is True
        and transport.get("campaign_public_key_pinned") is True
        and transport.get("detached_signature_required") is True
        and transport.get("monotonic_sequence_required") is True
        and transport.get("issued_and_expiry_required") is True
        and set(transport.get("contract_delivery_modes", []))
        == {"embedded", "exact_provider_https_url"}
        and transport.get("contract_retrieved_bytes_digest_verified") is True
        and transport.get("attended_import_mode") == "recovery_only"
        and transport.get("attended_import_counts_as_channel_delivery") is False
        and transport.get(
            "attended_import_may_authorize_installed_preauthorization"
        ) is False,
        "M2 first provider-channel transport or recovery boundary drifted",
    )
    _require(
        {
            "provider_or_channel", "channel_origin_and_authentication_key",
            "repository", "package_root", "audience_or_version_selectors",
            "allowed_actions", "authorization_mode", "paths", "commands",
            "agent_execution_policy",
            "agent_runner_adapter_version_executable_and_conformance_policy",
            "agent_runner_auth_and_entitlement_policy",
            "agent_runner_execution_funding_mode",
            "agent_runner_credential_owner",
            "agent_runner_usage_billing_owner",
            "agent_runner_native_configuration_policy",
            "model_data_and_budgets", "github_token_issuance_policy",
            "provider_reporting", "retention_and_deletion", "disclosure",
            "expiry", "revocation",
        }
        == set(channel.get("consumer_installation_fields", []))
        and {
            "notify_only", "scan_only", "prepare_patch", "open_draft_pr",
        }
        == set(channel.get("consumer_installation_action_modes", []))
        and set(channel.get("consumer_installation_authorization_modes", []))
        == {"per_event_approval", "installed_preauthorization"}
        and channel.get("event_specific_authorization_snapshot_required") is True
        and channel.get("event_may_widen_installation") is False
        and channel.get("agent_execution_policy_states")
        == ["disabled", "configured"]
        and channel.get("agent_execution_policy_default") == "disabled"
        and channel.get("runner_fields_required_only_when_configured") is True
        and channel.get("disabled_policy_grants_agent_authority") is False
        and channel.get(
            "agent_assisted_route_requires_configured_policy"
        ) is True
        and channel.get(
            "deterministic_route_may_proceed_with_disabled_policy"
        ) is True,
        "M2 Consumer Installation scope or derived authorization drifted",
    )
    _require(
        channel.get("agent_runner_selection_consumer_owned") is True
        and channel.get("provider_may_select_agent_runner") is False
        and channel.get(
            "stored_reusable_agent_or_model_credential_allowed"
        ) is False
        and set(channel.get("execution_funding_modes", []))
        == {"consumer_managed", "provider_sponsored_lumyn_managed"}
        and channel.get("launch_agent_runner_targets")
        == ["codex", "claude_code"]
        and channel.get("deferred_agent_runner_targets") == ["cursor"]
        and channel.get("common_adapter_conformance_required") is True
        and channel.get("executable_integrity_required") is True
        and channel.get("auth_and_entitlement_qualification_required") is True
        and channel.get("neutral_home_and_config_required") is True
        and channel.get("repository_path_shadowing_allowed") is False
        and channel.get("clean_session_required") is True
        and channel.get("silent_fallback_allowed") is False
        and channel.get("configured_route_preflight_required") is True
        and channel.get("configured_route_preflight_no_model_call") is True
        and channel.get("configured_route_actual_model_identity_required") is True,
        "M2 Agent Runner selection, funding, credential, or fallback posture drifted",
    )
    _require(
        channel.get("action_modes_are_ceilings") is True
        and channel.get("stored_github_token_allowed") is False
        and channel.get("per_event_approval_binds_exact_plan") is True
        and channel.get(
            "installed_preauthorization_requires_all_bound_values_in_policy"
        ) is True
        and channel.get("action_mode_alone_grants_side_effect") is False,
        "M2 installation authorization mode is ambient or under-specified",
    )
    status = channel.get("provider_status", {})
    _require(
        status.get("exact_event_and_evidence_binding") is True
        and status.get("consumer_consent_required") is True
        and set(status.get("provenance_labels", []))
        == {"observed", "consumer_reported", "unknown"}
        and status.get("silence_is_unknown") is True
        and status.get("merge_is_not_retired") is True
        and status.get("not_applicable_requires_explicit_evidence") is True
        and status.get("unaffected_requires_explicit_evidence") is True
        and status.get("raw_consumer_evidence") is False
        and status.get("raw_consumer_data_never_provider_visible") is True,
        "M2 provider status projection is not proof-honest",
    )
    _require(
        m2.get("managed_credential_contract") == MANAGED_CREDENTIAL_CONTRACT,
        "M2 managed credential broker contract drifted",
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
        "adapter_conformance_required",
        "live_canary_per_advertised_version",
        "clean_session_required",
        "executable_integrity_required",
        "auth_and_entitlement_qualification_required",
        "neutral_home_and_config_required",
        "native_configuration_cannot_widen",
        "raw_prompt_response_persistence_default",
        "personal_session_reuse_allowed",
        "native_configuration_default_enabled",
        "repository_path_shadowing_allowed",
        "silent_fallback_allowed",
        "agent_reported_verification_qualifies",
        "requires_configured_agent_execution_policy",
        "actual_model_route_identity_required",
        "noninteractive_automation_entitlement_required",
    ):
        expected = field not in {
            "raw_prompt_response_persistence_default",
            "personal_session_reuse_allowed",
            "native_configuration_default_enabled",
            "repository_path_shadowing_allowed",
            "silent_fallback_allowed",
            "agent_reported_verification_qualifies",
        }
        _require(
            contract.get(field) is expected,
            f"{task_id} bounded-agent {field} posture drifted",
        )
    _require(
        contract.get("launch_adapters") == ["codex", "claude_code"]
        and contract.get("deferred_adapters") == ["cursor"]
        and set(contract.get("execution_funding_modes", []))
        == {"consumer_managed", "provider_sponsored_lumyn_managed"},
        f"{task_id} Agent Runner launch or funding modes drifted",
    )
    _require(
        task.get("runner_host_isolation_contract")
        == RUNNER_HOST_ISOLATION_CONTRACT,
        f"{task_id} runner host-isolation contract drifted",
    )
    _require(
        task.get("managed_credential_contract")
        == MANAGED_CREDENTIAL_CONTRACT,
        f"{task_id} managed credential broker contract drifted",
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
        "fresh_verification_view_required",
        "separate_verifier_process_required",
        "agent_runner_and_model_credentials_absent",
        "generation_session_cannot_write_verification_evidence",
        "verification_command_digest_frozen",
        "verification_configuration_digest_frozen",
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
        "bound_agent_execution_policy",
        "bound_agent_runner_adapter",
        "bound_agent_runner_version",
        "bound_agent_runner_executable_digest",
        "bound_actual_model_route",
        "bound_execution_funding_mode",
        "bound_credential_owner",
        "bound_usage_billing_owner",
        "route_change_requires_new_explicit_authorization",
        "route_change_requires_new_attempt",
        "route_change_requires_new_candidate",
        "agent_execution_policy_configured_required",
        "failed_agent_candidate_reuses_bound_route",
        "deterministic_or_manual_candidate_requires_new_agent_route_authorization",
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
        repair.get("repair_route") == "agent_assisted_only"
        and repair.get("disabled_policy_without_new_authorization_behavior")
        == "needs_input_or_blocked",
        "M7 repair must be agent-assisted and fail closed without configured authorization",
    )
    _require(
        contract.get("verification_evidence_writer")
        == "independent_verifier_evidence_boundary",
        "M7 verification evidence writer must be independent from generation",
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
        contract.get("fallback_forms")
        == ["evidence_bundle", "patch", "optional_local_branch", "pr_ready_bundle"],
        "M9 must preserve the multi-form local export fallback",
    )
    _require(
        contract.get("required_remote_forms") == ["remote_branch", "draft_pr"],
        "M9 required remote forms drifted",
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
        contract.get("short_lived_least_privilege_token") is True
        and contract.get("github_app_installation_may_persist") is True
        and contract.get("long_lived_token_allowed") is False
        and contract.get("broad_organization_grant_allowed") is False
        and contract.get("non_default_branch_only") is True
        and contract.get("idempotency_evidence_required") is True
        and contract.get("tested_candidate_evidence_required") is True,
        "M9 automated delivery must be short-lived, non-default, tested, and idempotent",
    )
    _require(
        contract.get("manual_delivery_cannot_close_exp_003") is True
        and contract.get("missing_automated_evidence_behavior")
        == "EXP-003 remains open",
        "M9 cannot waive automated draft-PR proof",
    )
    status = contract.get("provider_status_contract", {})
    _require(
        status.get("exact_event_and_evidence_binding") is True
        and status.get("consumer_consented_fields_only") is True
        and status.get("silence_is_unknown") is True
        and status.get("merge_is_not_retired") is True
        and status.get("not_applicable_requires_explicit_evidence") is True
        and status.get("unaffected_requires_explicit_evidence") is True
        and status.get("raw_consumer_evidence") is False,
        "M9 provider status must be consented, evidence-bound, and proof-honest",
    )
    _require(
        set(contract.get("idempotency_binding_fields", []))
        == {
            "provider_change_event", "provider_change_contract",
            "consumer_installation_authorization", "repository",
            "repository_base", "candidate_head", "migration_plan",
            "verification_evidence",
        },
        "M9 idempotency must bind the full event-to-evidence identity",
    )
    composed = task.get("composed_update_contract", {})
    _require(
        composed.get("command") == "lumyn update --event"
        and composed.get("required_stages")
        == [
            "provider_change_event", "consumer_installation",
            "repository_impact", "migration_plan", "lumyn_candidate",
            "independent_verification", "remote_branch", "draft_pr",
            "provider_status_projection",
        ]
        and set(composed.get("qualifying_candidate_modes", []))
        == {"deterministic", "agent_assisted"}
        and composed.get("installed_action_ceiling_required") == "open_draft_pr"
        and set(composed.get("authorization_modes_supported", []))
        == {"per_event_approval", "installed_preauthorization"}
        and composed.get("imported_manual_candidate_qualifies") is False
        and composed.get("standalone_pr_create_qualifies") is False
        and composed.get("attended_event_import_qualifies") is False
        and composed.get("pilot_requires_installed_preauthorization") is True,
        "M9 composed provider-event-to-draft-PR proof drifted",
    )
    _require(
        composed.get("provider_status_projection_generated_locally") is True
        and composed.get("provider_transmission_optional_for_exp_003") is True
        and composed.get("pilot_same_run_provider_projection_required") is True,
        "M9 status projection must stay bound to the composed proof run",
    )
    route_composition = task.get("delivery_route_composition_contract", {})
    expected_route_refs = {
        "M4/impact_read_only",
        "M6/deterministic_candidate",
        "M6/deterministic_package_tool_candidate",
        "M6/agent_assisted_candidate",
        "M7/verify",
        "M9/remote_branch_push",
        "M9/draft_pr_create",
        "M9/provider_status_decline",
        "M9/provider_status_transmit",
    }
    _require(
        isinstance(route_composition, dict)
        and set(route_composition.get("delegated_route_refs", []))
        == expected_route_refs
        and route_composition.get("deterministic_sequence")
        == [
            "M4/impact_read_only",
            "M6/deterministic_candidate",
            "M7/verify",
            "M9/remote_branch_push",
            "M9/draft_pr_create",
        ]
        and route_composition.get("deterministic_package_tool_sequence")
        == [
            "M4/impact_read_only",
            "M6/deterministic_package_tool_candidate",
            "M7/verify",
            "M9/remote_branch_push",
            "M9/draft_pr_create",
        ]
        and route_composition.get("agent_assisted_sequence")
        == [
            "M4/impact_read_only",
            "M6/agent_assisted_candidate",
            "M7/verify",
            "M9/remote_branch_push",
            "M9/draft_pr_create",
        ]
        and route_composition.get("terminal_provider_status_routes")
        == [
            "M9/provider_status_decline",
            "M9/provider_status_transmit",
        ]
        and route_composition.get("pilot_terminal_provider_status_route")
        == "M9/provider_status_transmit"
        and route_composition.get(
            "each_delegated_action_freezes_own_exact_union"
        ) is True
        and route_composition.get(
            "aggregate_cross_action_scope_union_authorized"
        ) is False
        and route_composition.get(
            "branch_push_and_pr_create_separate_actions"
        ) is True,
        "M9 delivery must compose atomic delegated actions without aggregate authority",
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


def _validate_repair_path(task: dict[str, Any]) -> None:
    _require(
        "internal/authorization/" in task.get("allowed_paths", []),
        "M7 must own the future internal/authorization path",
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
        == "matched_agent_engine_and_budget_status_quo",
        "M10 generic-agent baseline must be fair",
    )
    _require(
        set(contract.get("generic_agent_control_variables", []))
        == GENERIC_AGENT_CONTROL_FIELDS
        and contract.get("difference_under_test") == GENERIC_AGENT_TREATMENT
        and contract.get("unmatched_engine_comparison_is_causal") is False,
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
    _require(
        contract.get("provider_led_distribution_and_onboarding_required") is True
        and contract.get("minimum_lumyn_opened_tested_draft_prs") == 1
        and contract.get("minimum_composed_lumyn_draft_prs") == 1
        and contract.get(
            "composed_pr_requires_installed_preauthorization"
        ) is True
        and set(contract.get("composed_pr_candidate_modes", []))
        == {"deterministic", "agent_assisted"}
        and contract.get("bespoke_operator_edits_qualify") is False
        and contract.get("standalone_pr_creation_qualifies") is False
        and contract.get("manual_only_delivery_fails_campaign") is True,
        "M10 must prove one composed Lumyn-generated installed-event draft PR",
    )
    _require(
        contract.get(
            "qualifying_composed_pr_requires_organic_agent_assisted_plan_item"
        ) is True
        and contract.get("deterministic_reroute_for_agent_proof_allowed")
        is False
        and contract.get(
            "separate_agent_and_deterministic_composed_runs_qualify"
        ) is False,
        "M10 qualifying composed proof must itself contain an organic agent-assisted item",
    )
    _require(
        contract.get("provider_status_exact_event_and_evidence_binding") is True
        and contract.get("provider_status_consumer_consent_required") is True
        and contract.get("silence_is_unknown") is True
        and contract.get("merge_is_not_retired") is True,
        "M10 rollout status must be consented, event-bound, and proof-honest",
    )
    _require(
        contract.get("minimum_provider_status_projections") == 1
        and contract.get(
            "provider_status_projection_must_bind_composed_pr"
        ) is True
        and contract.get("explicit_status_decline_record_required") is True,
        "M10 must prove at least one real provider status projection",
    )
    _require(
        contract.get("campaign_verdict_values") == ["pass", "fail"]
        and contract.get("reframe_is_post_failure_disposition") is True
        and contract.get("abandonment_or_timeout_verdict") == "fail",
        "M10 verdict must be pass/fail with reframe only after failure",
    )
    _require(
        contract.get("campaign_economics_threshold")
        and contract.get("campaign_economics_threshold_pass_required") is True,
        "M10 must freeze and pass a contribution-margin or automation threshold",
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
    _require(
        contract.get("minimum_real_agent_assisted_outcomes") == 1
        and contract.get("real_consumer_repository_required") is True
        and contract.get(
            "agent_assisted_requires_consumer_selected_qualified_runner"
        ) is True
        and contract.get(
            "agent_assisted_independent_exact_head_verification_required"
        ) is True
        and contract.get("agent_assisted_bespoke_operator_edits_qualify")
        is False
        and contract.get(
            "deterministic_only_closes_agent_runner_product_proof"
        ) is False,
        "M10 must prove one real consumer-selected agent-assisted component of the qualifying same-run outcome",
    )
    _require(
        set(contract.get("execution_funding_modes", []))
        == {"consumer_managed", "provider_sponsored_lumyn_managed"}
        and contract.get(
            "funding_mode_selected_per_agent_enabled_installation"
        ) is True
        and contract.get(
            "runner_configuration_required_only_for_agent_enabled_installations"
        ) is True
        and contract.get("credential_owner_recorded") is True
        and contract.get("usage_billing_owner_recorded") is True
        and contract.get("api_provider_agent_access") is False
        and contract.get("provider_raw_consumer_data_access") is False
        and "Agent Runner" in str(contract.get("campaign_cogs", "")),
        "M10 must attribute Agent Runner funding, credentials, billing, and COGS without provider access",
    )
    composition = task.get("campaign_route_composition_contract")
    _require(
        isinstance(composition, dict),
        "M10 campaign route composition contract is required",
    )
    _require(
        set(composition.get("delegated_route_refs", []))
        == {
            "M4/impact_read_only",
            "M6/deterministic_candidate",
            "M6/deterministic_package_tool_candidate",
            "M6/agent_assisted_candidate",
            "M7/verify",
            "M7/repair_agent_assisted",
            "M8/sandbox_read_back",
            "M9/local_export",
            "M9/local_branch",
            "M9/remote_branch_push",
            "M9/draft_pr_create",
            "M9/provider_status_decline",
            "M9/provider_status_transmit",
        }
        and composition.get(
            "per_installation_event_run_route_sequence_frozen"
        ) is True
        and composition.get("exact_union_frozen_before_each_action") is True
        and composition.get("cross_route_aggregate_union_authorized") is False
        and composition.get(
            "campaign_authority_arrays_are_aggregate_maxima_not_grants"
        ) is True,
        "M10 must compose exact delegated routes without ambient campaign authority",
    )
    _require(
        set(composition.get("required_campaign_proof_routes", []))
        == {
            "M6/agent_assisted_candidate",
            "M9/remote_branch_push",
            "M9/draft_pr_create",
            "M9/provider_status_transmit",
        }
        and composition.get(
            "campaign_minimums_grant_per_consumer_authority"
        ) is False,
        "M10 proof minima must not grant every consumer agent, GitHub, or reporting authority",
    )
    _require(
        composition.get("required_same_run_route_sequence")
        == QUALIFYING_SAME_RUN_ROUTE_SEQUENCE
        and composition.get("required_same_run_id_cardinality") == 1
        and composition.get("separate_runs_may_satisfy_sequence") is False,
        "M10 qualifying proof routes must execute in one bound run",
    )
    same_run = composition.get("qualifying_same_run_evidence_binding", {})
    _require(
        isinstance(same_run, dict)
        and same_run.get("binding_fields")
        == QUALIFYING_SAME_RUN_EVIDENCE_FIELDS
        and same_run.get("all_fields_required") is True
        and same_run.get("one_value_per_field") is True
        and same_run.get("all_artifacts_bind_same_run_id") is True
        and same_run.get(
            "all_artifacts_bind_same_event_and_installation_authorization"
        ) is True
        and same_run.get(
            "projection_binds_candidate_verification_and_draft_pr"
        ) is True
        and same_run.get("cross_run_evidence_allowed") is False
        and same_run.get("missing_or_mismatched_binding_result")
        == "PILOT-003 remains open",
        "M10 qualifying same-run evidence binding is incomplete",
    )
    _require(
        contract.get("qualifying_same_run_route_sequence_ref")
        == (
            "campaign_route_composition_contract."
            "required_same_run_route_sequence"
        )
        and contract.get("qualifying_same_run_evidence_binding_ref")
        == (
            "campaign_route_composition_contract."
            "qualifying_same_run_evidence_binding"
        ),
        "M10 paid campaign must reference the qualifying same-run route and evidence contracts",
    )


def _resolve_product_route_ref(
    tasks: dict[str, dict[str, Any]],
    route_ref: str,
) -> dict[str, Any]:
    parts = route_ref.split("/", 1)
    _require(
        len(parts) == 2 and all(parts),
        f"invalid delegated product route ref {route_ref!r}",
    )
    task_id, route_id = parts
    _require(
        task_id in tasks,
        f"delegated product route ref names unknown task {task_id}",
    )
    routes = tasks[task_id].get("product_action_route_contract", {})
    _require(
        isinstance(routes, dict) and route_id in routes,
        f"delegated product route ref does not resolve: {route_ref}",
    )
    route = routes[route_id]
    _require(
        isinstance(route, dict),
        f"delegated product route is not an object: {route_ref}",
    )
    return route


def validate_delegated_route_refs(
    tasks: dict[str, dict[str, Any]],
) -> None:
    """Prove composed route references resolve to the exact source contracts."""

    for task_id in ("M9", "M10"):
        key = (
            "delivery_route_composition_contract"
            if task_id == "M9"
            else "campaign_route_composition_contract"
        )
        composition = tasks[task_id].get(key, {})
        for route_ref in composition.get("delegated_route_refs", []):
            _resolve_product_route_ref(tasks, route_ref)

    m10_routes = tasks["M10"].get("product_action_route_contract", {})
    _require(
        isinstance(m10_routes, dict) and bool(m10_routes),
        "M10 delegated product route contract is required",
    )
    for route_id, delegated_route in m10_routes.items():
        _require(
            isinstance(delegated_route, dict),
            f"M10 delegated route {route_id} must be an object",
        )
        route_ref = delegated_route.get("delegated_route_ref")
        _require(
            isinstance(route_ref, str),
            f"M10 delegated route {route_id} must name its source route",
        )
        source_route = _resolve_product_route_ref(tasks, route_ref)
        target_route = {
            key: value
            for key, value in delegated_route.items()
            if key != "delegated_route_ref"
        }
        _require(
            target_route == source_route,
            f"M10 delegated route {route_id} differs from {route_ref}",
        )
    _require(
        {
            route["delegated_route_ref"]
            for route in m10_routes.values()
        }
        == set(
            tasks["M10"]["campaign_route_composition_contract"][
                "delegated_route_refs"
            ]
        ),
        "M10 campaign composition refs must equal its delegated action routes",
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
    m25_paths = set(tasks["M2.5"].get("allowed_paths", []))
    _require(
        m25_paths == {".factory/artifacts/product-signals/M2.5/"},
        "M2.5 writes must remain limited to task-scoped product-signal evidence",
    )
    _require(
        "docs/product/" in set(tasks["M2.5"].get("forbidden_paths", [])),
        "M2.5 must forbid writes to product source-of-truth documents",
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
    _validate_update_delivery_contracts(tasks["M1"], tasks["M2"])
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
    _validate_repair_path(tasks["M7"])
    _validate_export(tasks["M9"])
    _validate_campaign(tasks["M10"])
    for task_id in ("M6", "M7", "M8", "M9", "M10"):
        _validate_repository_command_isolation(tasks[task_id], task_id)
    for task_id in ("M8", "M10"):
        _validate_sandbox_entrypoint_isolation(tasks[task_id], task_id)
    validate_delegated_route_refs(tasks)
    for task in tasks.values():
        _validate_paused_runtime(task)
        _validate_worker_chain(task)
