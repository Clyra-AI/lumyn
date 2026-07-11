#!/usr/bin/env python3
from __future__ import annotations

from typing import Any

from repo_pack_contracts import fail


def validate_source_parser_review(tasks_by_id: dict[str, dict[str, Any]]) -> None:
    task = tasks_by_id.get("T3.1")
    if not isinstance(task, dict):
        fail("task-packets.json missing T3.1 source-parser repair task")
    review = task.get("required_review")
    if not isinstance(review, dict) or review.get("required") is not True:
        fail("T3.1.required_review.required must be true")
    if review.get("review_type") != "architecture":
        fail("T3.1.required_review.review_type must be architecture")
    if review.get("reviewer_class") != "peer_agent":
        fail("T3.1.required_review.reviewer_class must be peer_agent")
    lifecycle = task.get("lifecycle_gates") or {}
    if lifecycle.get("code_review_required") is not True:
        fail("T3.1.lifecycle_gates.code_review_required must be true")


def validate_recorder_task_split(tasks_by_id: dict[str, dict[str, Any]]) -> None:
    if "T4" in tasks_by_id:
        fail("task-packets.json must split broad T4 into T4.1, T4.2, and T4.3 before recorder dispatch")
    required = {
        "T4.1": {"RCRR-003", "RCRR-009", "FR8", "NFR5"},
        "T4.2": {"RCRR-004", "ACT-004", "FR5", "FR7", "NFR7", "NFR13"},
        "T4.3": {"REC-QUALITY-001", "NFR2"},
    }
    for task_id_value, item_ids in required.items():
        task = tasks_by_id.get(task_id_value)
        if not isinstance(task, dict):
            fail(f"task-packets.json missing recorder split task {task_id_value}")
        actual_ids = {str(value) for value in task.get("acceptance_item_ids", [])}
        missing = sorted(item_ids - actual_ids)
        if missing:
            fail(f"{task_id_value}.acceptance_item_ids missing recorder split ids: {missing}")
        if "v0.0" not in {str(value) for value in task.get("delivery_slice_refs", [])}:
            fail(f"{task_id_value}.delivery_slice_refs must include v0.0")
    t43_ids = {str(value) for value in tasks_by_id["T4.3"].get("acceptance_item_ids", [])}
    t43_acceptance_checks = "\n".join(str(value) for value in tasks_by_id["T4.3"].get("acceptance_checks", []))
    if "ACT-003" in t43_ids or "ACT-003" in t43_acceptance_checks:
        fail("T4.3 must not own ACT-003 because replay-verified timing belongs to T5/T6")
    t42_deps = {str(value) for value in tasks_by_id["T4.2"].get("blocked_by", [])}
    if "T4.1" not in t42_deps:
        fail("T4.2 must depend on T4.1")
    t43_deps = {str(value) for value in tasks_by_id["T4.3"].get("blocked_by", [])}
    if "T4.2" not in t43_deps:
        fail("T4.3 must depend on T4.2")
    t5_deps = {str(value) for value in tasks_by_id["T5.1"].get("blocked_by", [])}
    if "T4.3" not in t5_deps or "T4" in t5_deps:
        fail("T5.1 must depend on T4.3, not the removed broad T4 packet")
    quality_checks = "\n".join(str(value).lower() for value in tasks_by_id["T4.3"].get("acceptance_checks", []))
    for token in ["strong-spec", "weak-spec", "non-claims", "70 percent"]:
        if token not in quality_checks:
            fail(f"T4.3 acceptance_checks must define REC-QUALITY-001 measurement token {token!r}")


def validate_first_session_smoke_task(tasks_by_id: dict[str, dict[str, Any]]) -> None:
    task = tasks_by_id.get("T6.2")
    if not isinstance(task, dict):
        fail("task-packets.json missing T6.2")
    smoke = task.get("first_session_smoke")
    if not isinstance(smoke, dict):
        fail("T6.2.first_session_smoke is required")
    if smoke.get("required") is not True:
        fail("T6.2.first_session_smoke.required must be true")
    if smoke.get("command") != "make smoke-first-session":
        fail("T6.2.first_session_smoke.command must be make smoke-first-session")
    if smoke.get("report_ref") != ".factory/artifacts/task-runs/T6.2/first-session-smoke.json":
        fail("T6.2.first_session_smoke.report_ref must be task-scoped")
    smoke_ids = {str(value) for value in smoke.get("acceptance_item_ids", [])}
    if not {"ACT-001", "ACT-002", "ACT-003"}.issubset(smoke_ids):
        fail("T6.2.first_session_smoke.acceptance_item_ids must cover ACT-001, ACT-002, and ACT-003")
    if "first_session_smoke_report" not in {str(value) for value in task.get("evidence_required", [])}:
        fail("T6.2.evidence_required must include first_session_smoke_report")
    if "make smoke-first-session" not in {str(value) for value in task.get("validation_commands", [])}:
        fail("T6.2.validation_commands must include make smoke-first-session")
    allowed_paths = {str(value) for value in task.get("allowed_paths", [])}
    missing_paths = sorted({"Makefile", "scripts/"} - allowed_paths)
    if missing_paths:
        fail(f"T6.2.allowed_paths must include smoke target implementation paths: {missing_paths}")
