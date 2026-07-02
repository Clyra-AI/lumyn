#!/usr/bin/env python3
from __future__ import annotations

from pathlib import Path
from typing import Any

from repo_pack_contracts import fail, validate_no_legacy_provider_fields


def validate_safety_corpus_ready_plan(
    prd_path: Path,
    plan: dict[str, Any],
    packets: dict[str, Any],
    contract: dict[str, Any],
    ledger: dict[str, Any],
    mapping: dict[str, Any],
    scope: dict[str, Any],
) -> None:
    prd_text = prd_path.read_text(encoding="utf-8").lower()
    for token in [
        "safety and corpus-ready evidence",
        "corpus_eligible: false",
        "boundary violations",
    ]:
        if token not in prd_text:
            fail(f"docs/product/prd.md missing safety/corpus-ready token: {token}")

    required_by_task = {
        "T2.7": {"FDN-003"},
        "T6.2": {"RCRR-012"},
        "T8": {"LVCIS-009"},
        "T9": {"LVCIS-010"},
        "T11.2": {"EVAL-011"},
        "T12.2": {"EVAL-011"},
    }
    required_ids = set().union(*required_by_task.values())
    ledger_items = {
        str(item.get("acceptance_item_id")): item
        for item in ledger.get("items", [])
        if isinstance(item, dict)
    }
    missing_ledger_ids = sorted(required_ids - set(ledger_items))
    if missing_ledger_ids:
        fail(f"acceptance-ledger.json missing safety/corpus ids: {missing_ledger_ids}")
    for item_id in required_ids:
        source_ref = str(ledger_items[item_id].get("source_ref", ""))
        if source_ref != "docs/product/prd.md#safety-and-corpus-ready-evidence":
            fail(f"acceptance-ledger.json {item_id} must cite the safety/corpus PRD section")

    criteria = "\n".join(str(value) for value in contract.get("acceptance_criteria", []))
    for item_id in sorted(required_ids):
        if item_id not in criteria:
            fail(f"validation-contract.json acceptance_criteria missing {item_id}")

    def delivery_slice_ids(document: dict[str, Any], slice_id: str) -> set[str]:
        slice_item = next(
            (
                item
                for item in document.get("delivery_slices", [])
                if isinstance(item, dict) and item.get("slice_id") == slice_id
            ),
            None,
        )
        if not isinstance(slice_item, dict):
            fail(f"{slice_id} delivery slice missing")
        return {str(value) for value in slice_item.get("acceptance_item_ids", [])}

    for label, document in [
        ("execution-plan.json", plan),
        ("validation-contract.json", contract),
        ("acceptance-mapping.json", mapping),
        ("scope-closure-map.json", scope),
    ]:
        if not {"FDN-003", "RCRR-012"}.issubset(delivery_slice_ids(document, "v0.0")):
            fail(f"{label} v0.0 delivery slice missing safety/corpus record/report ids")
        if not {"LVCIS-009", "LVCIS-010"}.issubset(delivery_slice_ids(document, "v0.1")):
            fail(f"{label} v0.1 delivery slice missing boundary/CI safety ids")
        if "EVAL-011" not in delivery_slice_ids(document, "v0.2"):
            fail(f"{label} v0.2 delivery slice missing eval failure-event id")

    tasks = {
        str(task.get("task_id")): task
        for task in packets.get("tasks", [])
        if isinstance(task, dict)
    }
    if "T2.7" not in tasks:
        fail("task-packets.json missing T2.7 safety/corpus contract task")
    if "T2.7" not in set(str(value) for value in tasks["T3"].get("blocked_by", [])):
        fail("T3 must depend on T2.7 before product implementation resumes")
    for task_id_value, ids in required_by_task.items():
        task = tasks.get(task_id_value)
        if not isinstance(task, dict):
            fail(f"task-packets.json missing {task_id_value}")
        task_ids = {str(value) for value in task.get("acceptance_item_ids", [])}
        missing = sorted(ids - task_ids)
        if missing:
            fail(f"{task_id_value}.acceptance_item_ids missing safety/corpus ids: {missing}")
        checks = "\n".join(str(value).lower() for value in task.get("acceptance_checks", []))
        if not any(token in checks for token in ["safety", "corpus", "boundary"]):
            fail(f"{task_id_value}.acceptance_checks must describe safety/corpus-ready evidence")
        refs = task.get("safety_corpus_ready_evidence_refs")
        if not isinstance(refs, list) or not refs:
            fail(f"{task_id_value}.safety_corpus_ready_evidence_refs must be non-empty")
        ref_text = "\n".join(
            f"{item.get('source_ref', '')} {item.get('rule', '')}".lower()
            for item in refs
            if isinstance(item, dict)
        )
        if "safety-and-corpus-ready-evidence" not in ref_text:
            fail(f"{task_id_value}.safety_corpus_ready_evidence_refs must cite the PRD safety/corpus section")

    plan_text = "\n".join(
        str(value).lower()
        for value in plan.get("definition_of_done", []) + plan.get("explicit_non_goals", [])
    )
    for token in ["corpus_eligible", "boundary", "hosted corpus"]:
        if token not in plan_text:
            fail(f"execution-plan.json missing safety/corpus plan token: {token}")


def validate_risk_classification(risk: dict[str, Any]) -> None:
    validate_no_legacy_provider_fields(risk, "risk-classification.json")
    rules = risk.get("risk_rules")
    if not isinstance(rules, list):
        fail("risk-classification.json must contain risk_rules list")
    high = next((rule for rule in rules if isinstance(rule, dict) and rule.get("risk_class") == "high"), None)
    if not isinstance(high, dict):
        fail("risk-classification.json missing high risk rule")
    applies = "\n".join(str(value).lower() for value in high.get("applies_to", []))
    if "openai-compatible" not in applies or "anthropic" not in applies:
        fail("high risk rule must name both OpenAI-compatible and Anthropic provider key surfaces")
