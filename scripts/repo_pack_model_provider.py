#!/usr/bin/env python3
from __future__ import annotations

from typing import Any

from repo_pack_contracts import fail


def factoryd_config_capability_grants(config: dict[str, Any], repo_key: str) -> list[dict[str, Any]]:
    grants: list[dict[str, Any]] = []
    repos = config.get("repos")
    if isinstance(repos, dict):
        repo = repos.get(repo_key)
        if isinstance(repo, dict) and isinstance(repo.get("capability_grants"), list):
            grants.extend(grant for grant in repo["capability_grants"] if isinstance(grant, dict))
    elif isinstance(config.get("capability_grants"), list):
        grants.extend(grant for grant in config["capability_grants"] if isinstance(grant, dict))
    return grants


def missing_grant_value(value: Any) -> bool:
    if value is None or value == []:
        return True
    if isinstance(value, str):
        return not value.strip()
    return False


def task_id(task: dict[str, Any]) -> str:
    value = task.get("task_id")
    return value if isinstance(value, str) else ""


def validate_model_provider_gate(task: dict[str, Any], active_grants: list[dict[str, Any]] | None = None) -> None:
    active_grants = active_grants or []
    task_id_value = task_id(task)
    if task.get("requires_model_provider_endpoint") is not True:
        fail(f"{task_id_value}.requires_model_provider_endpoint must be true for live eval provider work")
    requirements = task.get("model_provider_requirements")
    if not isinstance(requirements, dict) or requirements.get("required_grant") != "model_provider_endpoint":
        fail(f"{task_id_value}.model_provider_requirements must require model_provider_endpoint")
    provider_surfaces = {str(value) for value in requirements.get("provider_surfaces", [])}
    missing_surfaces = {"openai_compatible_http", "anthropic_messages_http"} - provider_surfaces
    if missing_surfaces:
        fail(f"{task_id_value}.model_provider_requirements.provider_surfaces missing {sorted(missing_surfaces)}")
    required_fields = {str(value) for value in requirements.get("required_fields", [])}
    expected_required_fields = {
        "provider_identity",
        "provider_model",
        "provider_endpoint_or_base_url",
        "credential_environment",
        "budget_posture",
        "redaction_posture",
        "network_allowlist",
    }
    missing_required_fields = sorted(expected_required_fields - required_fields)
    if missing_required_fields:
        fail(f"{task_id_value}.model_provider_requirements.required_fields missing {missing_required_fields}")
    if task.get("requires_human_approval") is not False:
        fail(f"{task_id_value}.requires_human_approval must be false; model-only approval is represented by model_provider_endpoint grant")
    required_grant_fields = [
        "evidence_ref",
        "network_allowlist",
        "provider_identity",
        "provider_model",
        "credential_environment",
        "budget_posture",
        "redaction_posture",
    ]
    required_string_fields = [
        "evidence_ref",
        "provider_identity",
        "provider_model",
        "credential_environment",
        "budget_posture",
        "redaction_posture",
    ]

    def validate_provider_grant_fields(candidate: dict[str, Any], label: str) -> str:
        missing = [
            field for field in required_grant_fields
            if field not in candidate or missing_grant_value(candidate[field])
        ]
        if missing:
            fail(f"{label} missing fields: {missing}")
        non_string_fields = [
            field for field in required_string_fields
            if field in candidate and not isinstance(candidate.get(field), str)
        ]
        if non_string_fields:
            fail(f"{label} fields must be non-empty strings: {non_string_fields}")
        allowlist = candidate.get("network_allowlist")
        if not isinstance(allowlist, list) or not all(isinstance(item, str) and item.strip() for item in allowlist):
            fail(f"{label} network_allowlist must be a non-empty string list")
        provider_endpoint_value = candidate.get("provider_endpoint", "")
        base_url_value = candidate.get("base_url", "")
        if provider_endpoint_value not in (None, "") and not isinstance(provider_endpoint_value, str):
            fail(f"{label} provider_endpoint must be a string")
        if base_url_value not in (None, "") and not isinstance(base_url_value, str):
            fail(f"{label} base_url must be a string")
        provider_endpoint = provider_endpoint_value.strip() if isinstance(provider_endpoint_value, str) else ""
        base_url = base_url_value.strip() if isinstance(base_url_value, str) else ""
        provider_endpoint_or_base_url = provider_endpoint or base_url
        if not provider_endpoint_or_base_url:
            fail(f"{label} must include provider_endpoint or base_url")
        return provider_endpoint_or_base_url

    seed_grants = ((task.get("factoryd_runtime") or {}).get("capability_grants")) or []
    active_wildcard_grants = [
        grant for grant in active_grants
        if isinstance(grant, dict)
        and str(grant.get("task_id", "")).strip() == "*"
        and grant.get("capability") == "model_provider_endpoint"
    ]
    if active_wildcard_grants:
        fail(f"{task_id_value}.active model_provider_endpoint grants must be task-scoped, not wildcard")
    seed_matching = [
        grant for grant in seed_grants
        if isinstance(grant, dict)
        and str(grant.get("task_id", "")).strip() in {"*", task_id_value}
        and grant.get("capability") == "model_provider_endpoint"
    ]
    if any(grant.get("approved") is True for grant in seed_matching):
        fail(f"{task_id_value}.seed model_provider_endpoint grants must stay approved false; active approvals belong in .factory/factoryd.json")
    for seed_grant in seed_matching:
        validate_provider_grant_fields(seed_grant, f"{task_id_value}.seed model_provider_endpoint grant")
    active_matching = [
        grant for grant in active_grants
        if isinstance(grant, dict)
        and str(grant.get("task_id", "")).strip() == task_id_value
        and grant.get("capability") == "model_provider_endpoint"
    ]
    matching = [*seed_matching, *active_matching]
    if not matching:
        fail(f"{task_id_value} must include one seed wildcard or task-scoped model_provider_endpoint grant in factoryd_runtime.capability_grants, or one task-scoped active .factory/factoryd.json config grant")
    grant = next((candidate for candidate in matching if candidate.get("approved") is True), matching[0])
    approved = grant.get("approved")
    if approved not in (False, True):
        fail(f"{task_id_value}.model_provider_endpoint grant approved flag must be true or false")
    provider_endpoint_or_base_url = validate_provider_grant_fields(grant, f"{task_id_value}.model_provider_endpoint grant")
    if approved is True:
        checked_values = [
            grant.get("provider_identity"),
            grant.get("provider_model"),
            provider_endpoint_or_base_url,
            grant.get("credential_environment"),
            grant.get("budget_posture"),
            grant.get("redaction_posture"),
            *list(grant.get("network_allowlist") or []),
        ]
        if any("pending-approved" in str(value).lower() or str(value).lower().startswith("pending-") for value in checked_values):
            fail(f"{task_id_value}.approved model_provider_endpoint grant must not use pending placeholders")
    if "model_provider_endpoint" not in str(grant.get("evidence_ref")):
        fail(f"{task_id_value}.model_provider_endpoint grant evidence_ref must cite the alignment decision")
    joined_stop_conditions = "\n".join(str(value) for value in task.get("stop_conditions", []))
    if "model_provider_endpoint grant" not in joined_stop_conditions:
        fail(f"{task_id_value}.stop_conditions must fail closed without model_provider_endpoint grant")
