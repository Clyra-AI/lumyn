"""Validate immutable M1 lifecycle closure and landed-content evidence."""

from __future__ import annotations

import hashlib
import json
import subprocess
import sys
from pathlib import Path
from typing import Any


M1_VALIDATION_REPORT_REF = (
    ".factory/artifacts/task-runs/M1-IMPLEMENTATION/validation-report.json"
)
M1_SCORECARD_REF = (
    ".factory/artifacts/task-runs/M1-IMPLEMENTATION/"
    "proof-of-behavior-scorecard.json"
)
M1_REVIEW_REF = ".factory/artifacts/lifecycle-evidence/M1/review-report.json"
M1_HOLDOUT_REF = ".factory/artifacts/lifecycle-evidence/M1/holdout-result.json"
M1_IMPLEMENTATION_MARKER_REF = (
    ".factory/artifacts/pr-lifecycle/lumyn-v3-m1/implementation/"
    "work-proof-marker.json"
)
M1_PR_LIFECYCLE_REF = (
    ".factory/artifacts/pr-lifecycle/lumyn-v3-m1/pr-lifecycle-report.json"
)
M1_POST_MERGE_REF = (
    ".factory/artifacts/pr-lifecycle/lumyn-v3-m1/post-merge-report.json"
)
M1_SCOPE_CLOSURE_REF = (
    ".factory/artifacts/pr-lifecycle/lumyn-v3-m1/scope-closure-report.json"
)
M1_DEBT_REF = (
    ".factory/artifacts/pr-lifecycle/lumyn-v3-m1/delivery-debt-record.json"
)
M1_EXCEPTION_REF = (
    ".factory/artifacts/pr-lifecycle/lumyn-v3-m1/process-exception-handoff.json"
)
M1_AUTHORIZATION_REF = (
    ".factory/artifacts/pr-lifecycle/lumyn-v3-m1/lifecycle-authorization.json"
)
M1_SHIP_PACKET_REF = (
    ".factory/artifacts/pr-lifecycle/lumyn-v3-m1/ship-packet.json"
)
M1_EVENT_LOG_REF = (
    ".factory/artifacts/pr-lifecycle/lumyn-v3-m1/mission-event-log.json"
)
M1_POST_MERGE_MARKER_REF = (
    ".factory/artifacts/pr-lifecycle/lumyn-v3-m1/post-merge-work-proof/"
    "work-proof-marker.json"
)
M1_RETAINED_BUNDLE_REF = (
    ".factory/artifacts/pr-lifecycle/lumyn-v3-m1/implementation/"
    "pr74-original-head.bundle"
)
M1_BINDING_VERIFIER_REF = (
    ".factory/artifacts/pr-lifecycle/lumyn-v3-m1/implementation/verify_binding.py"
)
M1_EXPECTED_VALIDATION_RUN = (
    "validation:2026-07-28T17:05:04Z:lumyn-m1-implementation-r2"
)
M1_EXPECTED_CANDIDATE = (
    "sha256:971af9dd40e66ea47b220bcb4a817c62705e9764926eeb84ab4787564b8b2ed4"
)
M1_EXPECTED_BASE = "d7ae311c391775a2517d56add6d57148d5891ef3"
M1_EXPECTED_VALIDATION_CHECKOUT = "bc17b58f128bc15c4e715aba15c505c20b224e35"
M1_EXPECTED_PR_HEAD = "656c1f0bbb61cd558500c2f1b91a5a8f084f4f29"
M1_EXPECTED_LANDED_HEAD = "702b5be8d53b46a8c2a394f0b00770f626a8bbdd"
M1_EXPECTED_ACCEPTANCE_ITEMS = {
    "PACK-001",
    "EVENT-001",
    "EVENT-002",
    "INSTALL-001",
    "INSTALL-002",
    "IMPACT-005",
    *(f"AGENT-{number:03d}" for number in range(1, 8)),
    "VER-006",
    "PILOT-005",
}
M1_EXPECTED_REMAINING_TASK_REFS = {
    "PACK-001": ["M3"],
    "EVENT-001": ["M3"],
    "EVENT-002": ["M3"],
    "INSTALL-001": ["M5"],
    "INSTALL-002": ["M5"],
    "IMPACT-005": ["M4"],
    **{f"AGENT-{number:03d}": ["M6"] for number in range(1, 8)},
    "VER-006": ["M7"],
    "PILOT-005": ["M10"],
}
M1_PARTIAL_IMPLEMENTED_TASK_REFS = {
    "PACK-001": ["M1"],
    "EVENT-001": ["M1", "M2"],
    "EVENT-002": ["M1", "M2"],
    "INSTALL-001": ["M1", "M2"],
    "INSTALL-002": ["M1", "M2"],
    "IMPACT-005": ["M1"],
    **{f"AGENT-{number:03d}": ["M1"] for number in range(1, 8)},
    "VER-006": ["M1"],
    "PILOT-005": ["M1"],
}
M1_CLOSURE_JSON_REFS = (
    M1_VALIDATION_REPORT_REF,
    M1_SCORECARD_REF,
    M1_REVIEW_REF,
    M1_HOLDOUT_REF,
    M1_IMPLEMENTATION_MARKER_REF,
    M1_PR_LIFECYCLE_REF,
    M1_POST_MERGE_REF,
    M1_SCOPE_CLOSURE_REF,
    M1_DEBT_REF,
    M1_EXCEPTION_REF,
    M1_AUTHORIZATION_REF,
    M1_SHIP_PACKET_REF,
    M1_EVENT_LOG_REF,
    M1_POST_MERGE_MARKER_REF,
)

_EVIDENCE_CACHE: dict[str, dict[str, Any]] | None = None


def _require(condition: bool, message: str) -> None:
    if not condition:
        raise AssertionError(message)


def _sha256_bytes(value: bytes) -> str:
    return f"sha256:{hashlib.sha256(value).hexdigest()}"


def _load_json(root: Path, ref: str) -> dict[str, Any]:
    path = root / ref
    _require(path.is_file(), f"M1 closure evidence is missing {ref}")
    try:
        payload = json.loads(path.read_text())
    except Exception as exc:
        raise AssertionError(f"M1 closure evidence is invalid JSON at {ref}: {exc}")
    _require(isinstance(payload, dict), f"M1 closure evidence must be an object: {ref}")
    return payload


def load_m1_closure_evidence(root: Path) -> dict[str, dict[str, Any]]:
    """Load evidence and reproduce the retained-head-to-landed binding."""

    global _EVIDENCE_CACHE
    if _EVIDENCE_CACHE is not None:
        return _EVIDENCE_CACHE

    payloads = {ref: _load_json(root, ref) for ref in M1_CLOSURE_JSON_REFS}
    marker = payloads[M1_IMPLEMENTATION_MARKER_REF]
    proof = marker.get("landed_binding_proof", {})
    _require(
        proof.get("retained_bundle_ref") == M1_RETAINED_BUNDLE_REF,
        "M1 retained PR-head bundle reference drifted",
    )
    bundle_path = root / M1_RETAINED_BUNDLE_REF
    _require(bundle_path.is_file(), "M1 retained PR-head bundle is missing")
    _require(
        _sha256_bytes(bundle_path.read_bytes())
        == proof.get("retained_bundle_sha256"),
        "M1 retained PR-head bundle digest differs from its work-proof marker",
    )

    verifier_path = root / M1_BINDING_VERIFIER_REF
    _require(verifier_path.is_file(), "M1 landed-binding verifier is missing")
    _require(
        _sha256_bytes(verifier_path.read_bytes()) == proof.get("verifier_sha256"),
        "M1 landed-binding verifier digest differs from its work-proof marker",
    )
    result = subprocess.run(
        [sys.executable, str(verifier_path)],
        cwd=root,
        check=False,
        capture_output=True,
    )
    _require(
        result.returncode == 0,
        "M1 retained candidate-to-landed binding verification failed",
    )
    try:
        verifier_result = json.loads(result.stdout)
    except Exception as exc:
        raise AssertionError(
            f"M1 landed-binding verifier returned invalid JSON: {exc}"
        )
    _require(
        isinstance(verifier_result, dict),
        "M1 landed-binding verifier must return an object",
    )
    assertions = verifier_result.get("assertions", {})
    _require(
        verifier_result.get("status") == "pass"
        and verifier_result.get("task_id") == "M1"
        and verifier_result.get("original_pr_head") == M1_EXPECTED_PR_HEAD
        and verifier_result.get("landed_main_head") == M1_EXPECTED_LANDED_HEAD
        and verifier_result.get("candidate_digest") == M1_EXPECTED_CANDIDATE
        and verifier_result.get("changed_path_count") == 81
        and verifier_result.get("retained_bundle_sha256")
        == proof.get("retained_bundle_sha256")
        and assertions
        and all(value is True for value in assertions.values()),
        "M1 landed-binding verifier output is incomplete or stale",
    )

    logs = marker.get("raw_log_refs", [])
    _require(
        isinstance(logs, list) and len(logs) == 2,
        "M1 landed-binding marker must bind stdout and stderr logs",
    )
    _require(
        all(isinstance(ref, str) and (root / ref).is_file() for ref in logs),
        "M1 landed-binding raw logs are missing",
    )
    stdout_bytes = (root / logs[0]).read_bytes()
    stderr_bytes = (root / logs[1]).read_bytes()
    _require(
        _sha256_bytes(stdout_bytes) == marker.get("stdout_digest")
        and _sha256_bytes(stderr_bytes) == marker.get("stderr_digest")
        and stdout_bytes == result.stdout
        and stderr_bytes == result.stderr,
        "M1 landed-binding logs do not reproduce the committed passing marker",
    )
    post_marker = payloads[M1_POST_MERGE_MARKER_REF]
    post_logs = post_marker.get("raw_log_refs", [])
    _require(
        isinstance(post_logs, list)
        and len(post_logs) == 2
        and all(
            isinstance(ref, str) and (root / ref).is_file()
            for ref in post_logs
        ),
        "M1 post-merge marker must bind existing stdout and stderr logs",
    )
    post_stdout = (root / post_logs[0]).read_bytes()
    post_stderr = (root / post_logs[1]).read_bytes()
    _require(
        _sha256_bytes(post_stdout) == post_marker.get("stdout_digest")
        and _sha256_bytes(post_stderr) == post_marker.get("stderr_digest"),
        "M1 post-merge raw-log digests differ from their marker",
    )
    try:
        observed_runs = json.loads(post_stdout)
    except Exception as exc:
        raise AssertionError(f"M1 post-merge stdout is invalid JSON: {exc}")
    _require(
        isinstance(observed_runs, list)
        and {run.get("workflowName") for run in observed_runs} == {"validate", "codeql"}
        and all(
            run.get("headSha") == M1_EXPECTED_LANDED_HEAD
            and run.get("status") == "completed"
            and run.get("conclusion") == "success"
            for run in observed_runs
        ),
        "M1 post-merge raw evidence must contain successful exact-main workflows",
    )
    _EVIDENCE_CACHE = payloads
    return payloads


def _evidence(
    evidence_payloads: dict[str, dict[str, Any]],
    ref: str,
) -> dict[str, Any]:
    payload = evidence_payloads.get(ref)
    _require(isinstance(payload, dict), f"M1 closure evidence is missing {ref}")
    return payload


def _require_identity(payload: dict[str, Any], label: str) -> None:
    _require(payload.get("task_id") == "M1", f"{label} must bind task M1")
    _require(
        payload.get("work_item_id") == "lumyn-v3-m1",
        f"{label} must bind work item lumyn-v3-m1",
    )


def validate_m1_mapping_progress(by_group: dict[str, dict[str, Any]]) -> None:
    """Require group-level M1 progress without implying terminal acceptance."""

    for group_id, runtime_task, progress_items in (
        ("migration_pack", "M3", {"PACK-001"}),
        ("impact_and_integration_graph", "M4", {"IMPACT-005"}),
        (
            "bounded_hybrid_execution",
            "M6",
            {f"AGENT-{number:03d}" for number in range(1, 8)},
        ),
        ("verification", "M7", {"VER-006"}),
        ("provider_campaign_pilot", "M10", {"PILOT-005"}),
    ):
        progress = by_group[group_id].get("implementation_progress", {})
        _require(
            progress.get("completed_prerequisite_task_refs") == ["M1"]
            and progress.get("runtime_closure_task_ref") == runtime_task
            and progress.get("acceptance_status") == "planned"
            and set(progress.get("nonterminal_progress_item_ids", []))
            == progress_items
            and M1_SCOPE_CLOSURE_REF
            in progress.get("progress_evidence_refs", []),
            f"mapping {group_id} must preserve closed M1 prerequisite "
            "without terminal acceptance",
        )


def validate_m1_closure_evidence(
    evidence_payloads: dict[str, dict[str, Any]],
) -> None:
    """Enforce task closure without upgrading product acceptance."""

    validation = _evidence(evidence_payloads, M1_VALIDATION_REPORT_REF)
    scorecard = _evidence(evidence_payloads, M1_SCORECARD_REF)
    review = _evidence(evidence_payloads, M1_REVIEW_REF)
    holdout = _evidence(evidence_payloads, M1_HOLDOUT_REF)
    marker = _evidence(evidence_payloads, M1_IMPLEMENTATION_MARKER_REF)
    lifecycle = _evidence(evidence_payloads, M1_PR_LIFECYCLE_REF)
    post_merge = _evidence(evidence_payloads, M1_POST_MERGE_REF)
    closure = _evidence(evidence_payloads, M1_SCOPE_CLOSURE_REF)
    debt = _evidence(evidence_payloads, M1_DEBT_REF)
    exception = _evidence(evidence_payloads, M1_EXCEPTION_REF)
    authorization = _evidence(evidence_payloads, M1_AUTHORIZATION_REF)
    ship_packet = _evidence(evidence_payloads, M1_SHIP_PACKET_REF)
    event_log = _evidence(evidence_payloads, M1_EVENT_LOG_REF)
    post_marker = _evidence(evidence_payloads, M1_POST_MERGE_MARKER_REF)

    _require(
        validation.get("task_id") == "M1-IMPLEMENTATION"
        and validation.get("work_item_id")
        == "lumyn-v3-m1-attended-implementation"
        and validation.get("validation_run_id") == M1_EXPECTED_VALIDATION_RUN,
        "M1 validation report identity or validation-run binding drifted",
    )
    candidate_binding = validation.get("candidate_binding", {})
    _require(
        validation.get("result") == "pass"
        and validation.get("promotion_decision") == "blocked"
        and validation.get("checks")
        and all(
            check.get("status") == "pass"
            for check in validation.get("checks", [])
        )
        and candidate_binding.get("base_git_sha") == M1_EXPECTED_BASE
        and candidate_binding.get("candidate_digest") == M1_EXPECTED_CANDIDATE,
        "M1 validation must remain passing, candidate-bound, and lifecycle-blocked",
    )

    levels = {
        level.get("level"): level.get("status")
        for level in scorecard.get("levels", [])
        if isinstance(level, dict)
    }
    _require(
        scorecard.get("task_id") == "M1-IMPLEMENTATION"
        and scorecard.get("work_item_id")
        == "lumyn-v3-m1-attended-implementation"
        and scorecard.get("candidate_digest") == M1_EXPECTED_CANDIDATE
        and scorecard.get("overall_status") == "pass"
        and scorecard.get("required_proof_level") == "workflow_behavior"
        and levels.get("workflow_behavior") == "pass",
        "M1 proof scorecard is not passing at its required proof level",
    )

    for label, payload in (
        ("M1 review report", review),
        ("M1 holdout result", holdout),
        ("M1 PR lifecycle report", lifecycle),
        ("M1 post-merge report", post_merge),
        ("M1 scope-closure report", closure),
        ("M1 delivery-debt record", debt),
    ):
        _require_identity(payload, label)

    review_work = review.get("current_work", {})
    holdout_work = holdout.get("current_work", {})
    _require(
        review.get("verdict") == "approved"
        and review.get("review_type") == "code"
        and review_work.get("validation_run_id") == M1_EXPECTED_VALIDATION_RUN
        and review_work.get("candidate_digest") == M1_EXPECTED_CANDIDATE
        and review.get("findings")
        and all(
            finding.get("status") == "resolved"
            for finding in review.get("findings", [])
        ),
        "M1 independent review is not approved, resolved, and current-work bound",
    )
    _require(
        holdout.get("policy_mode") == "provision"
        and holdout.get("promotion_decision") == "pass"
        and holdout.get("cases_run") == 6
        and holdout.get("failing_cases") == []
        and holdout_work == review_work,
        "M1 holdout provisioning is not passing or bound to the reviewed work",
    )

    landed_proof = marker.get("landed_binding_proof", {})
    _require(
        marker.get("execution_status") == "pass"
        and marker.get("exit_code") == 0
        and marker.get("git_sha") == M1_EXPECTED_LANDED_HEAD
        and landed_proof.get("base_head") == M1_EXPECTED_BASE
        and landed_proof.get("validation_checkout_sha")
        == M1_EXPECTED_VALIDATION_CHECKOUT
        and landed_proof.get("original_pr_head") == M1_EXPECTED_PR_HEAD
        and landed_proof.get("landed_main_head") == M1_EXPECTED_LANDED_HEAD
        and landed_proof.get("validation_run_id") == M1_EXPECTED_VALIDATION_RUN
        and landed_proof.get("candidate_digest") == M1_EXPECTED_CANDIDATE
        and landed_proof.get("changed_path_count") == 81
        and landed_proof.get("retained_bundle_ref") == M1_RETAINED_BUNDLE_REF
        and str(landed_proof.get("retained_bundle_sha256", "")).startswith(
            "sha256:"
        ),
        "M1 landed work-proof marker is not passing or fully bound",
    )

    local_validation = lifecycle.get("local_validation", {})
    codex_review = lifecycle.get("codex_review", {})
    merge = lifecycle.get("merge", {})
    _require(
        lifecycle.get("status") == "complete"
        and local_validation.get("status") == "pass"
        and local_validation.get("validation_run_id")
        == M1_EXPECTED_VALIDATION_RUN
        and local_validation.get("candidate_digest") == M1_EXPECTED_CANDIDATE
        and lifecycle.get("code_review", {}).get("status") == "approved"
        and lifecycle.get("holdout", {}).get("status") == "pass"
        and lifecycle.get("holdout", {}).get("evaluation_claimed") is False
        and codex_review.get("configured_enabled") is True
        and codex_review.get("enabled") is False
        and codex_review.get("status") == "disabled"
        and codex_review.get("unresolved_actionable_comments") == 0
        and codex_review.get("late_review_after_merge") is False
        and merge.get("status") == "merged"
        and merge.get("merged_commit") == M1_EXPECTED_LANDED_HEAD
        and lifecycle.get("post_merge", {}).get("status") == "healthy"
        and any(
            gate.get("gate") == "codex_review"
            and gate.get("approved_exception_ref") == M1_EXCEPTION_REF
            for gate in lifecycle.get("skipped_gates", [])
        ),
        "M1 PR lifecycle report is incomplete or misstates its review exception",
    )

    _require(
        post_merge.get("status") == "healthy"
        and post_merge.get("merge_reference") == M1_EXPECTED_LANDED_HEAD
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
        "M1 post-merge report is not healthy on the landed commit",
    )
    post_proof = post_marker.get("post_merge_check_proof", {})
    _require(
        post_marker.get("execution_status") == "pass"
        and post_marker.get("exit_code") == 0
        and post_marker.get("git_sha") == M1_EXPECTED_LANDED_HEAD
        and post_proof.get("merge_commit") == M1_EXPECTED_LANDED_HEAD
        and post_proof.get("checks_observed")
        and all(
            check.get("status") == "pass"
            for check in post_proof.get("checks_observed", [])
        ),
        "M1 post-merge work-proof marker is not passing on the landed commit",
    )

    _require(
        ship_packet.get("artifact_type") == "ship_packet"
        and ship_packet.get("task_id") == "M1"
        and ship_packet.get("work_item_id") == "lumyn-v3-m1"
        and ship_packet.get("merge_ref") == M1_EXPECTED_LANDED_HEAD
        and ship_packet.get("merge_readiness")
        == "landed_with_approved_non_reusable_process_exception"
        and ship_packet.get("process_exception_ref") == M1_EXCEPTION_REF
        and ship_packet.get("gate_status", {}).get("closure")
        == "closed_with_debt_acceptance_nonterminal",
        "M1 ship packet must bind the landed task-specific closure disposition",
    )

    events = event_log.get("events", [])
    previous_digest = ""
    _require(
        event_log.get("append_only") is True
        and event_log.get("immutability_posture") == "append_only_no_mutation"
        and isinstance(events, list)
        and len(events) == 4,
        "M1 mission event log must be a four-event append-only chain",
    )
    for sequence, event in enumerate(events, start=1):
        _require(
            event.get("sequence") == sequence
            and event.get("previous_digest") == previous_digest,
            "M1 mission event-log sequence or previous digest drifted",
        )
        canonical = json.dumps(
            {key: value for key, value in event.items() if key != "event_digest"},
            sort_keys=True,
            separators=(",", ":"),
        ).encode()
        expected_digest = _sha256_bytes(canonical)
        _require(
            event.get("event_digest") == expected_digest,
            "M1 mission event-log digest chain is invalid",
        )
        previous_digest = expected_digest

    closure_items = {
        item.get("intent_item_id"): item
        for item in closure.get("closure_items", [])
        if isinstance(item, dict)
    }
    acceptance_results = {
        item.get("acceptance_item_id"): item
        for item in closure.get("acceptance_item_results", [])
        if isinstance(item, dict)
    }
    _require(
        closure.get("overall_status") == "closed_with_debt"
        and closure.get("acceptance_closure_claimed") is False
        and closure.get("rerun_required") is False
        and closure.get("blocker_refs") == []
        and set(closure_items)
        == {
            "M1-CORPUS-001",
            "M1-SKELETON-001",
            "M1-AGENT-BENCH-001",
            "M1-HOLDOUT-001",
            "M1-LIFECYCLE-001",
        }
        and closure_items["M1-LIFECYCLE-001"].get("classification")
        == "implemented_with_debt"
        and all(
            closure_items[item_id].get("classification") == "implemented"
            for item_id in set(closure_items) - {"M1-LIFECYCLE-001"}
        )
        and closure.get("delivery_debt_refs") == [M1_DEBT_REF],
        "M1 scope-closure report is incomplete or overstates debt-free closure",
    )
    _require(
        set(acceptance_results) == M1_EXPECTED_ACCEPTANCE_ITEMS
        and all(
            result.get("status") == "partial"
            and result.get("terminal") is False
            and result.get("remaining_task_refs")
            == M1_EXPECTED_REMAINING_TASK_REFS[item_id]
            for item_id, result in acceptance_results.items()
        ),
        "M1 scope closure must keep every linked product acceptance item nonterminal",
    )
    _require(
        debt.get("debt_type") == "skipped_required_latest_head_codex_review"
        and debt.get("severity") == "high"
        and debt.get("follow_up_action")
        == "retain_as_non_reusable_process_exception_and_require_clean_later_task_lifecycles",
        "M1 delivery debt does not preserve the skipped review gate",
    )
    _require(
        exception.get("artifact_type") == "handoff_record"
        and exception.get("work_item_id") == "lumyn-v3-m1"
        and exception.get("current_state") == "exception_approved"
        and exception.get("blockers") == []
        and exception.get("approval_card", {}).get("decision_status")
        == "approved"
        and exception.get("merge_authority_grant", {}).get("task_id") == "M1",
        "M1 process exception is missing, blocked, or not task scoped",
    )
    _require(
        authorization.get("artifact_type") == "lifecycle_authorization"
        and authorization.get("task_id") == "M1"
        and authorization.get("work_item_id") == "lumyn-v3-m1"
        and authorization.get("status") == "consumed"
        and authorization.get("consumed_at")
        and set(authorization.get("closure_evidence_refs", []))
        == {M1_PR_LIFECYCLE_REF, M1_POST_MERGE_REF, M1_SCOPE_CLOSURE_REF},
        "M1 lifecycle authorization must be consumed and closure-bound",
    )
