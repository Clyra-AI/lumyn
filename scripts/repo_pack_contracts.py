#!/usr/bin/env python3
from __future__ import annotations

import re
from typing import Any


REQUIRED_ACCEPTANCE_ITEM_IDS = {
    "FDN-001",
    "FDN-002",
    "FDN-003",
    "FDL-001",
    "FDAP-001",
    "CLI-AGENT-001",
    "REC-QUALITY-001",
    *{f"FR{index}" for index in range(1, 26)},
    *{f"NFR{index}" for index in range(1, 15)},
    *{f"RCRR-{index:03d}" for index in range(1, 13)},
    *{f"LVCIS-{index:03d}" for index in range(1, 11)},
    *{f"EVAL-{index:03d}" for index in range(1, 12)},
    *{f"ACT-{index:03d}" for index in range(1, 5)},
    *{f"PULL-{index:03d}" for index in range(1, 6)},
}

ACCEPTANCE_LEDGER_REF = ".factory/artifacts/prd-to-plan/lumyn-mvp/acceptance-ledger.json"
ACCEPTANCE_MAPPING_REF = ".factory/artifacts/prd-to-plan/lumyn-mvp/acceptance-mapping.json"
SCOPE_CLOSURE_MAP_REF = ".factory/artifacts/prd-to-plan/lumyn-mvp/scope-closure-map.json"
REQUIRED_LIVE_EVAL_DISPATCH_GATES = {"PULL-001", "PULL-004"}

REQUIRED_MVP_VERSION_SLICES = {
    "v0.0": {
        "capability_group_id": "record_contract_replay_report",
        "task_refs": {"T1", "T2", "T2.7", "T3", "T4.1", "T4.2", "T4.3", "T5.1", "T5.2", "T6.1", "T6.2", "T10"},
    },
    "v0.1": {
        "capability_group_id": "live_verify_boundary_ci_share",
        "task_refs": {"T7", "T8", "T9", "T10"},
    },
    "v0.2": {
        "capability_group_id": "live_agent_eval",
        "task_refs": {"T11.1", "T11.2", "T12.1", "T12.2"},
    },
}

TASK_VERSION_SLICE_REFS = {
    task_ref: {slice_id}
    for slice_id, spec in REQUIRED_MVP_VERSION_SLICES.items()
    for task_ref in spec["task_refs"]
}
TASK_VERSION_SLICE_REFS["T10"] = {"v0.0", "v0.1"}
DOTTED_TASK_PARENT_SLICE_EXEMPTIONS = {"T2.5", "T2.6"}

REQUIRED_ACCEPTANCE_TASK_REFS = {
    "FDN-003": {"T2.7"},
    "RCRR-012": {"T6.2"},
    "LVCIS-009": {"T8"},
    "LVCIS-010": {"T9"},
    "EVAL-011": {"T11.2", "T12.2"},
    "FR14": {"T3", "T4.1", "T4.2", "T4.3", "T5.2", "T6.2", "T7", "T8", "T9", "T10", "T11.2", "T12.2"},
    "CLI-AGENT-001": {"T3"},
    "NFR9": {"T3", "T4.1", "T4.2", "T4.3", "T5.2", "T6.1", "T7", "T8", "T9", "T10", "T11.2", "T12.2"},
    "NFR12": {"T4.1", "T4.2", "T4.3", "T5.2", "T6.2", "T7", "T8", "T9", "T10", "T11.2", "T12.2"},
    "FR9": {"T4.2", "T4.3", "T5.1", "T6.2", "T7", "T8", "T9", "T10", "T11.1", "T12.2"},
    "FR2": {"T3", "T4.1", "T4.2", "T4.3", "T5.2", "T6.1", "T10", "T11.2", "T12.2"},
    "NFR6": {"T4.1", "T4.2", "T4.3", "T5.2", "T6.1", "T7", "T10", "T11.1", "T12.2"},
}

REQUIRED_MVP_EVAL_PROVIDERS = [
    "openai_compatible_http_adapter",
    "anthropic_messages_http_adapter",
]

REQUIRED_MVP_EVAL_ADAPTERS = [
    "openai_compatible_http",
    "anthropic_messages_http",
]

LEGACY_PROVIDER_FIELD = "first_eval_provider"
REQUIRED_PROVIDER_DECISION_ID = "mvp_eval_providers"


def fail(message: str) -> None:
    raise AssertionError(message)


def has_nonempty_list(value: Any) -> bool:
    return isinstance(value, list) and any(isinstance(item, str) and item.strip() for item in value)


def has_nonempty_string(value: Any) -> bool:
    return isinstance(value, str) and bool(value.strip())


def has_nonempty_collection(value: Any) -> bool:
    return (isinstance(value, dict) and bool(value)) or (isinstance(value, list) and bool(value))


def has_required_string_refs(value: Any, expected_refs: list[str]) -> bool:
    if not isinstance(value, list):
        return False
    present = {item for item in value if isinstance(item, str) and item.strip()}
    return all(expected in present for expected in expected_refs)


def base_task_id(value: Any) -> str:
    if not isinstance(value, str):
        return ""
    match = re.match(r"^(T\d+(?:\.\d+)*)", value.strip(), re.IGNORECASE)
    candidate = match.group(1) if match else value.strip()
    if candidate in TASK_VERSION_SLICE_REFS or candidate in DOTTED_TASK_PARENT_SLICE_EXEMPTIONS:
        return candidate
    if "." in candidate:
        parent = candidate.split(".", 1)[0]
        if parent in TASK_VERSION_SLICE_REFS:
            return parent
    return candidate


def expected_task_version_slices(task_id_value: str) -> set[str]:
    return set(TASK_VERSION_SLICE_REFS.get(base_task_id(task_id_value), set()))


def validate_mvp_version_slice_coverage(value: Any, label: str) -> None:
    if not isinstance(value, list) or not value:
        fail(f"{label} must be a non-empty list")
    by_id = {
        str(item.get("slice_id")): item
        for item in value
        if isinstance(item, dict) and str(item.get("slice_id", "")).strip()
    }
    missing_slices = sorted(set(REQUIRED_MVP_VERSION_SLICES) - set(by_id))
    if missing_slices:
        fail(f"{label} missing required MVP version slices: {missing_slices}")
    for slice_id, spec in REQUIRED_MVP_VERSION_SLICES.items():
        item = by_id[slice_id]
        if item.get("required_for_full_mvp") is not True:
            fail(f"{label}.{slice_id}.required_for_full_mvp must be true")
        if item.get("public_release_boundary") is not False:
            fail(f"{label}.{slice_id}.public_release_boundary must be false")
        if item.get("capability_group_id") != spec["capability_group_id"]:
            fail(f"{label}.{slice_id}.capability_group_id must be {spec['capability_group_id']}")
        if not has_nonempty_string(item.get("source_ref")):
            fail(f"{label}.{slice_id}.source_ref is required")
        task_refs = {str(task_ref) for task_ref in item.get("task_refs", [])}
        missing_task_refs = sorted(spec["task_refs"] - task_refs)
        if missing_task_refs:
            fail(f"{label}.{slice_id}.task_refs missing {missing_task_refs}")
        group_refs = {str(group_ref) for group_ref in item.get("acceptance_group_refs", [])}
        if spec["capability_group_id"] not in group_refs:
            fail(f"{label}.{slice_id}.acceptance_group_refs must include {spec['capability_group_id']}")


def validate_delivery_slice_coverage(value: Any, label: str) -> None:
    if not isinstance(value, list) or not value:
        fail(f"{label} must be a non-empty list")
    by_id = {
        str(item.get("slice_id")): item
        for item in value
        if isinstance(item, dict) and str(item.get("slice_id", "")).strip()
    }
    missing_slices = sorted(set(REQUIRED_MVP_VERSION_SLICES) - set(by_id))
    if missing_slices:
        fail(f"{label} missing required delivery slices: {missing_slices}")
    for slice_id, spec in REQUIRED_MVP_VERSION_SLICES.items():
        item = by_id[slice_id]
        if item.get("required_for_completion") is not True:
            fail(f"{label}.{slice_id}.required_for_completion must be true")
        if item.get("public_release_boundary") is not False:
            fail(f"{label}.{slice_id}.public_release_boundary must be false")
        if not has_nonempty_string(item.get("source_ref")):
            fail(f"{label}.{slice_id}.source_ref is required")
        task_refs = {str(task_ref) for task_ref in item.get("task_refs", [])}
        missing_task_refs = sorted(spec["task_refs"] - task_refs)
        if missing_task_refs:
            fail(f"{label}.{slice_id}.task_refs missing {missing_task_refs}")
        group_refs = {str(group_ref) for group_ref in item.get("acceptance_group_refs", [])}
        if spec["capability_group_id"] not in group_refs:
            fail(f"{label}.{slice_id}.acceptance_group_refs must include {spec['capability_group_id']}")
        item_ids = item.get("acceptance_item_ids")
        if not isinstance(item_ids, list) or not item_ids:
            fail(f"{label}.{slice_id}.acceptance_item_ids must preserve item-level closure refs")


def iter_key_paths(value: Any, target_key: str, path: str = "$") -> list[str]:
    paths: list[str] = []
    if isinstance(value, dict):
        for key, item in value.items():
            child_path = f"{path}.{key}"
            if key == target_key:
                paths.append(child_path)
            paths.extend(iter_key_paths(item, target_key, child_path))
    elif isinstance(value, list):
        for index, item in enumerate(value):
            paths.extend(iter_key_paths(item, target_key, f"{path}[{index}]"))
    return paths


def validate_no_legacy_provider_fields(value: Any, label: str) -> None:
    paths = iter_key_paths(value, LEGACY_PROVIDER_FIELD)
    if paths:
        fail(f"{label} uses legacy {LEGACY_PROVIDER_FIELD} fields: {', '.join(paths)}")
