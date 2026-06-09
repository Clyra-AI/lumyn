#!/usr/bin/env python3
from __future__ import annotations

import json
import sys
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parents[1]
PILOT_DIR = ROOT / ".factory" / "artifacts" / "pilot" / "lumyn-mvp-slice"
REQUIRED_FILES = [
    "README.md",
    "work-proof-marker.json",
    "validation-report.json",
    "mission-event-log.json",
    "blockers.json",
    "scope-closure-report.json",
    "repair-loop/task-packet.json",
    "repair-loop/repair-report.json",
    "review-report.json",
    "ship-packet.json",
]

REQUIRED_WORK_PROOF_CHANGED_PATHS = {
    "Makefile",
    "README.md",
    "docs/dev/dev_guides.md",
    "docs/factory/README.md",
    "scripts/validate_factory_pilot_evidence.py",
    ".factory/artifacts/pilot/lumyn-mvp-slice/README.md",
    ".factory/artifacts/pilot/lumyn-mvp-slice/work-proof-marker.json",
    ".factory/artifacts/pilot/lumyn-mvp-slice/validation-report.json",
    ".factory/artifacts/pilot/lumyn-mvp-slice/mission-event-log.json",
    ".factory/artifacts/pilot/lumyn-mvp-slice/blockers.json",
    ".factory/artifacts/pilot/lumyn-mvp-slice/scope-closure-report.json",
    ".factory/artifacts/pilot/lumyn-mvp-slice/repair-loop/task-packet.json",
    ".factory/artifacts/pilot/lumyn-mvp-slice/repair-loop/repair-report.json",
    ".factory/artifacts/pilot/lumyn-mvp-slice/review-report.json",
    ".factory/artifacts/pilot/lumyn-mvp-slice/ship-packet.json",
}


def fail(message: str) -> None:
    raise AssertionError(message)


def load_json(relative_path: str) -> dict[str, Any]:
    path = PILOT_DIR / relative_path
    if not path.exists():
        fail(f"missing pilot evidence file: {path.relative_to(ROOT)}")
    try:
        payload = json.loads(path.read_text())
    except Exception as exc:
        fail(f"{path.relative_to(ROOT)} is not valid JSON: {exc}")
    if not isinstance(payload, dict):
        fail(f"{path.relative_to(ROOT)} must contain a JSON object")
    return payload


def require_ref(ref: str, label: str) -> None:
    if not isinstance(ref, str) or not ref.strip():
        fail(f"{label}: expected non-empty repo-relative ref")
    path_text = ref.split("#", 1)[0]
    ref_path = Path(path_text)
    path = ROOT / path_text
    if ref_path.is_absolute() or ".." in ref_path.parts:
        fail(f"{label}: ref must be repo-relative")
    if not path.exists():
        fail(f"{label}: ref target does not exist: {path_text}")


def require_ref_list(values: Any, label: str, *, allow_empty: bool = False) -> None:
    if not isinstance(values, list):
        fail(f"{label}: expected list")
    if not allow_empty and not values:
        fail(f"{label}: expected at least one ref")
    for index, value in enumerate(values):
        require_ref(value, f"{label}[{index}]")


def validate_scope_closure(report: dict[str, Any]) -> None:
    if report.get("overall_status") != "blocked":
        fail("scope-closure-report.overall_status must be blocked until the full MVP closes")
    if report.get("rerun_required") is not True:
        fail("scope-closure-report.rerun_required must be true while repair tasks remain")
    require_ref(report.get("mission_contract_ref", ""), "scope-closure-report.mission_contract_ref")
    require_ref(report.get("intent_source_ref", ""), "scope-closure-report.intent_source_ref")
    require_ref(report.get("event_log_ref", ""), "scope-closure-report.event_log_ref")
    require_ref_list(report.get("work_proof_marker_refs"), "scope-closure-report.work_proof_marker_refs")
    require_ref_list(report.get("validation_report_refs"), "scope-closure-report.validation_report_refs")
    require_ref_list(report.get("blocker_refs"), "scope-closure-report.blocker_refs")

    items = report.get("closure_items")
    if not isinstance(items, list) or len(items) < 4:
        fail("scope-closure-report.closure_items must cover baseline, current PRD slice, and live-work groups")
    implemented = [item for item in items if isinstance(item, dict) and item.get("classification") == "implemented"]
    gaps = [item for item in items if isinstance(item, dict) and item.get("classification") in {"not_implemented", "blocked"}]
    if not implemented:
        fail("scope-closure-report must contain implemented current-slice evidence")
    if not gaps:
        fail("scope-closure-report must contain the current incomplete PRD scope")

    repair_gap = None
    for item in items:
        if not isinstance(item, dict):
            fail("scope-closure-report.closure_items entries must be objects")
        for key in ["implementation_refs", "evidence_refs", "work_proof_marker_refs", "repair_task_packet_refs", "repair_report_refs"]:
            require_ref_list(item.get(key, []), f"scope-closure-report.{item.get('intent_item_id')}.{key}", allow_empty=True)
        if item.get("intent_item_id") == "record_contract_replay_report":
            repair_gap = item
    if repair_gap is None:
        fail("scope-closure-report must include record_contract_replay_report closure item")
    if repair_gap.get("classification") != "not_implemented":
        fail("record_contract_replay_report must stay not_implemented until T3-T6 close")
    require_ref_list(repair_gap.get("repair_task_packet_refs"), "record_contract_replay_report.repair_task_packet_refs")
    require_ref_list(repair_gap.get("repair_report_refs"), "record_contract_replay_report.repair_report_refs")


def validate_repair_loop(task_packet: dict[str, Any], repair_report: dict[str, Any]) -> None:
    expected_gap_ref = ".factory/artifacts/pilot/lumyn-mvp-slice/scope-closure-report.json#/closure_items/1"
    if task_packet.get("task_id") != "T3-repair-001":
        fail("repair task packet must target T3 as the first missing MVP slice")
    if task_packet.get("closure_gap_ref") != expected_gap_ref:
        fail("repair task packet must reference the record/replay/report closure gap")
    if repair_report.get("closure_gap_ref") != expected_gap_ref:
        fail("repair report must reference the same closure gap")
    if repair_report.get("repair_task_packet_ref") != ".factory/artifacts/pilot/lumyn-mvp-slice/repair-loop/task-packet.json":
        fail("repair report must link the bounded repair task packet")
    if repair_report.get("terminal_outcome") != "ready_for_rerun":
        fail("repair report must leave the loop ready for rerun, not falsely closed")
    if repair_report.get("retry_budget_remaining") != 1:
        fail("repair report must preserve a retry budget")
    require_ref(task_packet.get("closure_gap_ref", ""), "repair task closure_gap_ref")
    require_ref(repair_report.get("rerun_scope_closure_report_ref", ""), "repair report rerun_scope_closure_report_ref")


def validate_review_and_ship(review_report: dict[str, Any], ship_packet: dict[str, Any]) -> None:
    if review_report.get("verdict") != "approved":
        fail("review-report.verdict must approve the evidence package before shipping")
    if review_report.get("approval_effect", {}).get("blocks_promotion") is not False:
        fail("review-report.approval_effect.blocks_promotion must be false")
    require_ref_list(review_report.get("evidence_refs"), "review-report.evidence_refs")
    if ship_packet.get("scope_closure_report_ref") != ".factory/artifacts/pilot/lumyn-mvp-slice/scope-closure-report.json":
        fail("ship-packet must reference the pilot scope closure report")
    commit_set = ship_packet.get("commit_set")
    if not isinstance(commit_set, list) or not commit_set or "pending-pr-head" in commit_set:
        fail("ship-packet.commit_set must record concrete commit refs")
    if ship_packet.get("merge_readiness") != "ready_after_ci_and_codex_review":
        fail("ship-packet.merge_readiness must wait for CI and Codex review")
    require_ref(ship_packet.get("validation_reference", ""), "ship-packet.validation_reference")
    require_ref_list(ship_packet.get("work_proof_marker_refs"), "ship-packet.work_proof_marker_refs")


def validate_work_proof(work_proof: dict[str, Any]) -> None:
    if work_proof.get("execution_status") != "completed":
        fail("work-proof-marker.execution_status must be completed")
    changed_paths = work_proof.get("changed_paths")
    if not isinstance(changed_paths, list):
        fail("work-proof-marker.changed_paths must be a list")
    changed_path_set = {str(path) for path in changed_paths}
    missing = sorted(REQUIRED_WORK_PROOF_CHANGED_PATHS - changed_path_set)
    if missing:
        fail(f"work-proof-marker.changed_paths missing pilot proof paths: {missing}")
    forbidden_paths = work_proof.get("forbidden_paths_touched")
    if not isinstance(forbidden_paths, list) or forbidden_paths:
        fail("work-proof-marker.forbidden_paths_touched must be an empty list")


def main() -> int:
    try:
        for relative_path in REQUIRED_FILES:
            path = PILOT_DIR / relative_path
            if not path.exists():
                fail(f"missing pilot evidence file: {path.relative_to(ROOT)}")

        closure = load_json("scope-closure-report.json")
        repair_task = load_json("repair-loop/task-packet.json")
        repair_report = load_json("repair-loop/repair-report.json")
        review_report = load_json("review-report.json")
        ship_packet = load_json("ship-packet.json")
        validation_report = load_json("validation-report.json")
        work_proof = load_json("work-proof-marker.json")

        if validation_report.get("status") != "passed":
            fail("validation-report.status must be passed")
        validate_work_proof(work_proof)
        validate_scope_closure(closure)
        validate_repair_loop(repair_task, repair_report)
        validate_review_and_ship(review_report, ship_packet)
    except AssertionError as exc:
        print(f"factory pilot evidence validation failed: {exc}", file=sys.stderr)
        return 1

    print(json.dumps({"status": "pass", "pilot_evidence": "lumyn-mvp-slice"}, sort_keys=True))
    return 0


if __name__ == "__main__":
    sys.exit(main())
