#!/usr/bin/env python3
from __future__ import annotations

from typing import Any

from repo_pack_contracts import (
    ACCEPTANCE_LEDGER_REF,
    ACCEPTANCE_MAPPING_REF,
    REQUIRED_ACCEPTANCE_ITEM_IDS,
    REQUIRED_ACCEPTANCE_TASK_REFS,
    REQUIRED_MVP_EVAL_PROVIDERS,
    REQUIRED_MVP_VERSION_SLICES,
    SCOPE_CLOSURE_MAP_REF,
    fail,
    has_nonempty_list,
    has_nonempty_string,
    has_required_string_refs,
    validate_delivery_slice_coverage,
    validate_mvp_version_slice_coverage,
    validate_no_legacy_provider_fields,
)


def validate_acceptance_ledger_coverage(value: Any, label: str) -> None:
    if not isinstance(value, dict):
        fail(f"{label} must be an object")
    expected_refs = {
        "ledger_ref": ACCEPTANCE_LEDGER_REF,
        "acceptance_mapping_ref": ACCEPTANCE_MAPPING_REF,
        "scope_closure_map_ref": SCOPE_CLOSURE_MAP_REF,
    }
    for field, expected in expected_refs.items():
        if value.get(field) != expected:
            fail(f"{label}.{field} must cite {expected}")
    if value.get("coverage_unit") != "acceptance_item":
        fail(f"{label}.coverage_unit must be acceptance_item")
    if value.get("group_only_refs_allowed") is not False:
        fail(f"{label}.group_only_refs_allowed must be false")
    if value.get("required_item_count") != len(REQUIRED_ACCEPTANCE_ITEM_IDS):
        fail(f"{label}.required_item_count must match acceptance-ledger item count")
    if value.get("status") != "mapped":
        fail(f"{label}.status must be mapped")
    if not has_nonempty_list(value.get("required_groups")):
        fail(f"{label}.required_groups must be non-empty")
    required_slices = value.get("required_version_slices")
    if not isinstance(required_slices, list) or set(str(item) for item in required_slices) != set(REQUIRED_MVP_VERSION_SLICES):
        fail(f"{label}.required_version_slices must list v0.0, v0.1, and v0.2")


def validate_acceptance_ledger(ledger: dict[str, Any]) -> set[str]:
    validate_no_legacy_provider_fields(ledger, "acceptance-ledger.json")
    if ledger.get("artifact_type") != "acceptance_ledger":
        fail("acceptance-ledger.json artifact_type must be acceptance_ledger")
    if ledger.get("source_prd_ref") != "docs/product/prd.md":
        fail("acceptance-ledger.json source_prd_ref must point at docs/product/prd.md")
    policy = ledger.get("coverage_policy")
    if not isinstance(policy, dict):
        fail("acceptance-ledger.json coverage_policy must be an object")
    if policy.get("enumerated_items_required") is not True:
        fail("acceptance-ledger.json coverage_policy.enumerated_items_required must be true")
    if policy.get("group_only_refs_allowed") is not False:
        fail("acceptance-ledger.json coverage_policy.group_only_refs_allowed must be false")
    if policy.get("closure_unit") != "acceptance_item":
        fail("acceptance-ledger.json coverage_policy.closure_unit must be acceptance_item")
    items = ledger.get("items")
    if not isinstance(items, list) or not items:
        fail("acceptance-ledger.json must contain items")
    seen: set[str] = set()
    task_refs_by_item_id: dict[str, set[str]] = {}
    for index, item in enumerate(items):
        if not isinstance(item, dict):
            fail(f"acceptance-ledger.json items[{index}] must be an object")
        item_id = str(item.get("acceptance_item_id", "")).strip()
        if not item_id:
            fail(f"acceptance-ledger.json items[{index}] missing acceptance_item_id")
        if item_id in seen:
            fail(f"acceptance-ledger.json duplicate acceptance_item_id {item_id}")
        seen.add(item_id)
        for key in ["group_id", "source_ref", "source_text", "kind", "evidence_mode", "status", "risk_class"]:
            if not has_nonempty_string(item.get(key)):
                fail(f"acceptance-ledger.json {item_id} missing {key}")
        if not has_nonempty_list(item.get("closure_required_for")):
            fail(f"acceptance-ledger.json {item_id} missing closure_required_for")
        if not has_nonempty_list(item.get("task_refs")):
            fail(f"acceptance-ledger.json {item_id} missing task_refs")
        task_refs_by_item_id[item_id] = {str(value) for value in item.get("task_refs", [])}
        if item.get("status") == "implemented" and not has_nonempty_list(item.get("validation_refs")):
            fail(f"acceptance-ledger.json {item_id} implemented item missing validation_refs")
    missing = sorted(REQUIRED_ACCEPTANCE_ITEM_IDS - seen)
    if missing:
        fail(f"acceptance-ledger.json missing required item ids: {missing}")
    if "REC-QUALITY-001" not in seen:
        fail("acceptance-ledger.json must include recorder 70 percent quality gate")
    recorder_quality = next(
        (item for item in items if isinstance(item, dict) and item.get("acceptance_item_id") == "REC-QUALITY-001"),
        None,
    )
    if not isinstance(recorder_quality, dict):
        fail("acceptance-ledger.json missing REC-QUALITY-001")
    if recorder_quality.get("source_ref") != "docs/product/prd.md#phase-1-recorder-spike":
        fail("REC-QUALITY-001 source_ref must point at the Phase 1 recorder spike PRD anchor")
    for item_id, required_task_refs in REQUIRED_ACCEPTANCE_TASK_REFS.items():
        actual_task_refs = task_refs_by_item_id.get(item_id, set())
        missing_task_refs = sorted(required_task_refs - actual_task_refs)
        if missing_task_refs:
            fail(f"acceptance-ledger.json {item_id} missing required task_refs: {missing_task_refs}")
    return seen


def validate_acceptance_mapping(mapping: dict[str, Any], ledger_ids: set[str], contract: dict[str, Any]) -> None:
    validate_no_legacy_provider_fields(mapping, "acceptance-mapping.json")
    if mapping.get("acceptance_ledger_ref") != ACCEPTANCE_LEDGER_REF:
        fail("acceptance-mapping.json must cite acceptance-ledger.json")
    validate_mvp_version_slice_coverage(
        mapping.get("mvp_required_version_slices"),
        "acceptance-mapping.json.mvp_required_version_slices",
    )
    validate_delivery_slice_coverage(
        mapping.get("delivery_slices"),
        "acceptance-mapping.json.delivery_slices",
    )
    groups = mapping.get("groups")
    if not isinstance(groups, list):
        fail("acceptance-mapping.json must contain groups list")
    mapped_ids: set[str] = set()
    mapped_groups: set[str] = set()
    for group in groups:
        if not isinstance(group, dict):
            continue
        group_id = str(group.get("group_id", "")).strip()
        if not group_id:
            fail("acceptance-mapping.json group missing group_id")
        mapped_groups.add(group_id)
        item_ids = group.get("acceptance_item_ids")
        if not isinstance(item_ids, list) or not item_ids:
            fail(f"acceptance group {group.get('group_id')} missing acceptance_item_ids")
        unknown = sorted(str(value) for value in item_ids if str(value) not in ledger_ids)
        if unknown:
            fail(f"acceptance group {group.get('group_id')} references unknown acceptance item ids: {unknown}")
        mapped_ids.update(str(value) for value in item_ids)
    live_eval = next((group for group in groups if isinstance(group, dict) and group.get("group_id") == "live_agent_eval"), None)
    if not isinstance(live_eval, dict):
        fail("acceptance-mapping.json missing live_agent_eval group")
    if not has_required_string_refs(live_eval.get("provider_adapter_coverage"), REQUIRED_MVP_EVAL_PROVIDERS):
        fail("live_agent_eval.provider_adapter_coverage must include both MVP provider adapters")
    approvals = "\n".join(str(value).lower() for value in live_eval.get("requires_human_approval", []))
    if "openai" not in approvals or "anthropic" not in approvals:
        fail("live_agent_eval.requires_human_approval must name both provider credential postures")
    for required_id in ["EVAL-001", "PULL-001", "PULL-004"]:
        if required_id not in set(str(value) for value in live_eval.get("acceptance_item_ids", [])):
            fail(f"live_agent_eval acceptance mapping missing {required_id}")
    contract_groups = contract.get("acceptance_groups")
    if not isinstance(contract_groups, list) or not contract_groups:
        fail("validation-contract.json must declare acceptance_groups")
    missing_contract_groups = sorted(str(group_id) for group_id in contract_groups if str(group_id) not in mapped_groups)
    if missing_contract_groups:
        fail(f"acceptance-mapping.json missing validation-contract groups: {missing_contract_groups}")
    nfr_group = next((group for group in groups if isinstance(group, dict) and group.get("group_id") == "nonfunctional_requirements"), None)
    if not isinstance(nfr_group, dict):
        fail("acceptance-mapping.json missing nonfunctional_requirements group")
    if nfr_group.get("source_ref") != "docs/product/prd.md#non-functional-requirements":
        fail("nonfunctional_requirements source_ref must match PRD heading anchor")
    missing = sorted(ledger_ids - mapped_ids)
    if missing:
        fail(f"acceptance-mapping.json does not map ledger ids: {missing}")


def validate_scope_closure_map(scope: dict[str, Any], ledger_ids: set[str]) -> None:
    validate_no_legacy_provider_fields(scope, "scope-closure-map.json")
    if scope.get("acceptance_ledger_ref") != ACCEPTANCE_LEDGER_REF:
        fail("scope-closure-map.json must cite acceptance-ledger.json")
    validate_mvp_version_slice_coverage(
        scope.get("mvp_required_version_slices"),
        "scope-closure-map.json.mvp_required_version_slices",
    )
    validate_delivery_slice_coverage(
        scope.get("delivery_slices"),
        "scope-closure-map.json.delivery_slices",
    )
    items = scope.get("items")
    if not isinstance(items, list):
        fail("scope-closure-map.json must contain items list")
    closure_ids: set[str] = set()
    for item in items:
        if not isinstance(item, dict):
            continue
        item_ids = item.get("acceptance_item_ids")
        statuses = item.get("acceptance_item_statuses")
        if not isinstance(item_ids, list) or not item_ids:
            fail(f"scope item {item.get('scope_item')} missing acceptance_item_ids")
        if not isinstance(statuses, list) or not statuses:
            fail(f"scope item {item.get('scope_item')} missing acceptance_item_statuses")
        scope_item_id = str(item.get("scope_item_id", "")).strip()
        for slice_id, spec in REQUIRED_MVP_VERSION_SLICES.items():
            if scope_item_id != spec["capability_group_id"]:
                continue
            actual_slices = {str(value) for value in item.get("mvp_required_version_slices", [])}
            if slice_id not in actual_slices:
                fail(f"scope item {scope_item_id} missing mvp_required_version_slices entry {slice_id}")
            delivery_slices = {str(value) for value in item.get("delivery_slice_refs", [])}
            if slice_id not in delivery_slices:
                fail(f"scope item {scope_item_id} missing delivery_slice_refs entry {slice_id}")
        status_ids = {str(status.get("acceptance_item_id")) for status in statuses if isinstance(status, dict)}
        missing_status = sorted(str(value) for value in item_ids if str(value) not in status_ids)
        if missing_status:
            fail(f"scope item {item.get('scope_item')} missing item statuses for {missing_status}")
        closure_ids.update(str(value) for value in item_ids)
    live_eval = next((item for item in items if isinstance(item, dict) and item.get("scope_item") == "Live agent eval"), None)
    if not isinstance(live_eval, dict):
        fail("scope-closure-map.json missing Live agent eval item")
    if not has_required_string_refs(live_eval.get("provider_adapter_coverage"), REQUIRED_MVP_EVAL_PROVIDERS):
        fail("Live agent eval provider_adapter_coverage must include both MVP provider adapters")
    blockers = "\n".join(str(value).lower() for value in live_eval.get("blockers", []))
    if "openai" not in blockers or "anthropic" not in blockers:
        fail("Live agent eval blockers must name both provider credential postures")
    missing = sorted(ledger_ids - closure_ids)
    if missing:
        fail(f"scope-closure-map.json does not cover ledger ids: {missing}")
    nfr_item = next((item for item in items if isinstance(item, dict) and item.get("scope_item_id") == "nonfunctional_requirements"), None)
    if not isinstance(nfr_item, dict):
        fail("scope-closure-map.json missing nonfunctional_requirements scope item")
    if nfr_item.get("source_ref") != "docs/product/prd.md#non-functional-requirements":
        fail("nonfunctional_requirements scope source_ref must match PRD heading anchor")
