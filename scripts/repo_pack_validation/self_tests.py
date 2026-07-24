"""Negative tests for the Lumyn v3 repo-pack validator."""

from __future__ import annotations

import copy
from collections.abc import Callable
from typing import Any

from repo_pack_validation.authority import (
    conditional_activation_scope_digest,
    manual_preflight_scope_digest,
    validate_authority_grants,
)
from repo_pack_validation.markdown_refs import _markdown_anchors
from repo_pack_validation.task_contracts import _policy_digest


Payload = dict[str, dict[str, Any]]
ValidateLoaded = Callable[..., dict[str, dict[str, Any]]]


def _require(condition: bool, message: str) -> None:
    if not condition:
        raise AssertionError(message)


def _expect_failure(
    base: Payload,
    mutate: Callable[[Payload], Any],
    expected: str,
    validate_loaded: ValidateLoaded,
) -> None:
    candidate = copy.deepcopy(base)
    mutate(candidate)
    try:
        validate_loaded(candidate, validate_configs=False)
    except AssertionError as exc:
        _require(
            expected.lower() in str(exc).lower(),
            f"self-test expected {expected!r}, got {exc!r}",
        )
        return
    raise AssertionError(f"self-test mutation did not fail: {expected}")


def _task(payload: Payload, task_id: str) -> dict[str, Any]:
    return next(
        task
        for task in payload["packets"]["tasks"]
        if task.get("task_id") == task_id
    )


def _ledger_item(payload: Payload, item_id: str) -> dict[str, Any]:
    return next(
        item
        for item in payload["ledger"]["items"]
        if item.get("acceptance_item_id") == item_id
    )


def _closure_item(payload: Payload, item_id: str) -> dict[str, Any]:
    return next(
        item
        for item in payload["closure"]["items"]
        if item.get("scope_item_id") == item_id
    )


def _remove_task_acceptance(
    payload: Payload,
    task_id: str,
    acceptance_item_id: str,
) -> None:
    task = _task(payload, task_id)
    task["acceptance_item_ids"].remove(acceptance_item_id)
    task["acceptance_checks"] = [
        check
        for check in task["acceptance_checks"]
        if not str(check).startswith(f"{acceptance_item_id}:")
    ]


def _move_task_acceptance(
    payload: Payload,
    source_task_id: str,
    target_task_id: str,
    acceptance_item_id: str,
) -> None:
    _remove_task_acceptance(payload, source_task_id, acceptance_item_id)
    target = _task(payload, target_task_id)
    source_text = next(
        item["source_text"]
        for item in payload["ledger"]["items"]
        if item.get("acceptance_item_id") == acceptance_item_id
    )
    target["acceptance_item_ids"].append(acceptance_item_id)
    target["acceptance_checks"].append(
        f"{acceptance_item_id}: {source_text}"
    )


def _remove_holdout_prohibition(payload: Payload, field: str) -> None:
    policy = _task(payload, "M1")["holdout_policy"]
    policy["prohibited_committed_fields"].remove(field)
    policy["policy_digest"] = _policy_digest(policy)


def _change_holdout_baseline(payload: Payload, value: str) -> None:
    policy = _task(payload, "M1")["holdout_policy"]
    policy["comparison_baseline"] = value
    policy["policy_digest"] = _policy_digest(policy)


def _weaken_preflight_consent(payload: Payload) -> None:
    preflight = _task(payload, "M2.5")["manual_external_evidence_preflight"]
    preflight["participant_consent_required"] = False
    preflight["approval_scope_digest"] = manual_preflight_scope_digest(preflight)


def _conditional_grants(
    task_id: str,
    action_mode: str,
    selected_capabilities: list[str],
) -> list[dict[str, Any]]:
    selected = sorted(selected_capabilities)
    activation = {
        "task_id": task_id,
        "action_mode": action_mode,
        "selected_capabilities": selected,
        "evidence_ref": f"private:authorizations/{task_id}-activation.json",
        "expires_at": "2099-01-01T00:00:00Z",
    }
    activation["scope_digest"] = conditional_activation_scope_digest(activation)
    grants: list[dict[str, Any]] = []
    for capability in selected:
        grant: dict[str, Any] = {
            "task_id": task_id,
            "capability": capability,
            "approved": True,
            "evidence_ref": (
                f"private:authorizations/{task_id}-{capability}.json"
            ),
            "expires_at": activation["expires_at"],
            "conditional_activation": activation,
        }
        if capability == "credentials":
            grant["credential_scopes"] = ["bounded-worker"]
            grant["credential_environment"] = "consumer-controlled"
        elif capability == "network":
            grant["network_allowlist"] = ["api.example.com:443"]
        grants.append(grant)
    return grants


def run_repo_pack_self_tests(
    base: Payload,
    *,
    validate_loaded: ValidateLoaded,
    validate_config_payload: Callable[..., None],
    validate_active_config: Callable[[dict[str, Any], dict[str, dict[str, Any]]], None],
    historical_plan_rel: str,
    expected_capabilities: dict[str, set[str]],
) -> None:
    """Prove drift cannot enable dispatch, authority, false closure, or weak proof."""

    tasks = validate_loaded(base, validate_configs=False)
    _require(
        _markdown_anchors(
            "# API [Trust](https://example.com)\n"
            "# `Code` *Mode*\n"
        )
        == {"api-trust", "code-mode"},
        "Markdown heading slugs must use rendered inline text",
    )

    mutations: list[tuple[Callable[[Payload], Any], str]] = [
        (
            lambda value: value["ledger"]["items"].pop(),
            "compiled acceptance text differs from PRD",
        ),
        (
            lambda value: _ledger_item(value, "BASE-001")[
                "evidence_refs"
            ].remove(
                ".factory/artifacts/task-runs/T3/validation-report.json"
            ),
            "carry-forward evidence is semantically incomplete",
        ),
        (
            lambda value: _task(value, "M8")["required_worker_chain"].__setitem__(
                2, "ship-pr"
            ),
            "worker chain",
        ),
        (
            lambda value: _move_task_acceptance(
                value, "M0", "M1", "BASE-005"
            ),
            "primary acceptance ownership",
        ),
        (
            lambda value: _task(value, "M9").__setitem__("auto_merge", True),
            "forbidden product capability",
        ),
        (
            lambda value: _task(value, "M8")["factoryd_runtime"][
                "capability_grants"
            ][0].__setitem__("approved", True),
            "planning grants",
        ),
        (
            lambda value: value["mapping"]["groups"][0][
                "acceptance_item_ids"
            ].pop(),
            "exact PRD acceptance set",
        ),
        (
            lambda value: _task(value, "M9")["blocked_by"].append("M8"),
            "dependency graph",
        ),
        (
            lambda value: _task(value, "M2.5")["blocked_by"].append("M0"),
            "dependency graph",
        ),
        (
            lambda value: _task(value, "M2.5")["allowed_paths"].append(
                "docs/product/"
            ),
            "M2.5 writes must remain limited",
        ),
        (
            lambda value: _closure_item(value, "BASE-001")[
                "remaining_task_refs"
            ].clear(),
            "retain every active v3 task",
        ),
        (
            lambda value: _task(value, "M2.5")["gated_acceptance_items"][0].__setitem__(
                "required_milestone", "M1"
            ),
            "DISC-003",
        ),
        (
            _weaken_preflight_consent,
            "participant_consent_required",
        ),
        (
            lambda value: _task(value, "M0")["allowed_paths"].append(
                "scripts/repo_pack_validation/"
            ),
            "M0 allowed paths",
        ),
        (
            lambda value: _remove_holdout_prohibition(value, "answer_key"),
            "answer or resolving material",
        ),
        (
            lambda value: _change_holdout_baseline(
                value, "unbounded_generic_agent"
            ),
            "generic-agent baseline",
        ),
        (
            lambda value: _task(value, "M6")["bounded_agent_contract"][
                "exact_fields"
            ].remove("cost_budget"),
            "model control fields",
        ),
        (
            lambda value: _task(value, "M6")["bounded_agent_contract"][
                "untrusted_inputs"
            ].remove("repository_source"),
            "untrusted agent inputs",
        ),
        (
            lambda value: _task(value, "M7")[
                "deterministic_verification_contract"
            ].__setitem__("independent_from_generation", False),
            "verification independent_from_generation",
        ),
        (
            lambda value: _task(value, "M7")[
                "deterministic_verification_contract"
            ].__setitem__("verify_mutates_candidate", True),
            "explicitly non-mutating",
        ),
        (
            lambda value: _task(value, "M7")[
                "deterministic_verification_contract"
            ]["repair_authorization"].__setitem__(
                "prior_verification_evidence_invalidated", False
            ),
            "invalidate prior proof",
        ),
        (
            lambda value: _task(value, "M7")[
                "deterministic_verification_contract"
            ]["candidate_modes"].remove("imported_manual"),
            "imported manual candidates",
        ),
        (
            lambda value: _task(value, "M6")[
                "manual_candidate_contract"
            ].__setitem__("command", "git apply"),
            "manual candidate import",
        ),
        (
            lambda value: _task(value, "M9")["delivery_contract"].__setitem__(
                "auto_merge", True
            ),
            "forbidden product capability",
        ),
        (
            lambda value: _task(value, "M9")[
                "outcome_record_contract"
            ].__setitem__("append_only", False),
            "append-only",
        ),
        (
            lambda value: _task(value, "M9")[
                "outcome_record_contract"
            ]["durable_outcomes"].remove("reverted"),
            "reversion",
        ),
        (
            lambda value: _task(value, "M10")[
                "paid_campaign_contract"
            ].__setitem__("provider_code_access", True),
            "forbidden product capability",
        ),
        (
            lambda value: _task(value, "M10")[
                "paid_campaign_contract"
            ].__setitem__("payment_posture", "refundable_intent"),
            "cleared non-refundable",
        ),
        (
            lambda value: _task(value, "M10")[
                "paid_campaign_contract"
            ].__setitem__("material_provider_outcome_pass_required", False),
            "provider-outcome metric",
        ),
        (
            lambda value: _task(value, "M3")["factoryd_runtime"].__setitem__(
                "dispatch_enabled", True
            ),
            "dispatch must be disabled",
        ),
        (
            lambda value: value["plan"]["alignment_gate"].__setitem__(
                "implementation_may_start", True
            ),
            "implementation must remain blocked",
        ),
        (
            lambda value: value["context"]["alignment_decisions"]["resolved"][
                0
            ].__setitem__(
                "evidence_ref",
                "docs/product/prd.md#missing-planning-source",
            ),
            "missing Markdown anchor",
        ),
    ]
    for mutate, expected in mutations:
        _expect_failure(base, mutate, expected, validate_loaded)

    historical_config = copy.deepcopy(base["config"])
    historical_config["repos"]["lumyn"][
        "task_packets"
    ] = f"{historical_plan_rel}/task-packets.json"
    try:
        validate_config_payload(historical_config, "historical config")
    except AssertionError as exc:
        _require(
            "task_packets ref is stale" in str(exc),
            f"historical-plan self-test failed unexpectedly: {exc}",
        )
    else:
        raise AssertionError("historical plan remained selectable")

    enabled_config = copy.deepcopy(base["config"])
    enabled_config["repos"]["lumyn"]["runtime_control"]["mission_paused"] = False
    enabled_config["repos"]["lumyn"]["runtime_control"]["launch_request"][
        "expected_decision"
    ] = "allow"
    try:
        validate_config_payload(enabled_config, "enabled config")
    except AssertionError as exc:
        _require(
            "mission_paused" in str(exc),
            f"dispatch-pause self-test failed unexpectedly: {exc}",
        )
    else:
        raise AssertionError("an enabled v3 config remained selectable")

    product_grant = [
        {
            "task_id": "M9",
            "capability": "github_pr_write",
            "approved": True,
            "evidence_ref": "private:authorizations/M9.json",
            "expires_at": "2099-01-01T00:00:00Z",
        }
    ]
    try:
        validate_authority_grants(product_grant, tasks, expected_capabilities)
    except AssertionError as exc:
        _require(
            "unknown or product capability" in str(exc),
            f"product-capability separation failed unexpectedly: {exc}",
        )
    else:
        raise AssertionError("product authority was accepted as a Factory grant")

    for task_id, action_mode in (
        ("M6", "bounded_agent"),
        ("M9", "automated_draft_pr"),
    ):
        selected = list(tasks[task_id]["conditional_factory_capabilities"])
        validate_authority_grants(
            _conditional_grants(task_id, action_mode, selected),
            tasks,
            expected_capabilities,
        )

    missing_activation = _conditional_grants(
        "M6",
        "bounded_agent",
        list(tasks["M6"]["conditional_factory_capabilities"]),
    )
    missing_activation[0].pop("conditional_activation")
    try:
        validate_authority_grants(
            missing_activation, tasks, expected_capabilities
        )
    except AssertionError as exc:
        _require(
            "conditional_activation is required" in str(exc),
            f"missing conditional activation failed unexpectedly: {exc}",
        )
    else:
        raise AssertionError("conditional grant lacked activation evidence")

    partial_activation = _conditional_grants(
        "M6",
        "bounded_agent",
        list(tasks["M6"]["conditional_factory_capabilities"]),
    )
    partial_activation.pop()
    try:
        validate_authority_grants(
            partial_activation, tasks, expected_capabilities
        )
    except AssertionError as exc:
        _require(
            "exact selected conditional capabilities" in str(exc),
            f"partial conditional activation failed unexpectedly: {exc}",
        )
    else:
        raise AssertionError("partial conditional activation remained valid")

    wrong_task_activation = _conditional_grants(
        "M6",
        "bounded_agent",
        list(tasks["M6"]["conditional_factory_capabilities"]),
    )
    wrong_task_activation[0]["conditional_activation"]["task_id"] = "M9"
    try:
        validate_authority_grants(
            wrong_task_activation, tasks, expected_capabilities
        )
    except AssertionError as exc:
        _require(
            "must bind task M6" in str(exc),
            f"wrong-task conditional activation failed unexpectedly: {exc}",
        )
    else:
        raise AssertionError("conditional activation accepted the wrong task")

    extra_activation = _conditional_grants(
        "M9",
        "automated_draft_pr",
        list(tasks["M9"]["conditional_factory_capabilities"]),
    )
    activation = extra_activation[0]["conditional_activation"]
    activation["selected_capabilities"].append("github_pr_write")
    activation["selected_capabilities"].sort()
    activation["scope_digest"] = conditional_activation_scope_digest(activation)
    try:
        validate_authority_grants(
            extra_activation, tasks, expected_capabilities
        )
    except AssertionError as exc:
        _require(
            "selects undeclared capabilities" in str(exc),
            f"extra conditional capability failed unexpectedly: {exc}",
        )
    else:
        raise AssertionError("conditional activation accepted an extra capability")

    broad_network_grants = [
        {
            "task_id": "M8",
            "capability": "approval",
            "approved": True,
            "evidence_ref": "private:authorizations/M8-approval.json",
            "expires_at": "2099-01-01T00:00:00Z",
        },
        {
            "task_id": "M8",
            "capability": "credentials",
            "approved": True,
            "evidence_ref": "private:authorizations/M8-credentials.json",
            "expires_at": "2099-01-01T00:00:00Z",
            "credential_scopes": ["sandbox.read"],
            "credential_environment": "provider-sandbox",
        },
        {
            "task_id": "M8",
            "capability": "network",
            "approved": True,
            "evidence_ref": "private:authorizations/M8-network.json",
            "expires_at": "2099-01-01T00:00:00Z",
            "network_allowlist": ["ALL"],
        },
    ]
    try:
        validate_authority_grants(
            broad_network_grants, tasks, expected_capabilities
        )
    except AssertionError as exc:
        _require(
            "semantic wildcard" in str(exc),
            f"network wildcard failed unexpectedly: {exc}",
        )
    else:
        raise AssertionError("a wildcard Factory network grant remained selectable")

    active_config = copy.deepcopy(base["config"])
    active_config["repos"]["lumyn"]["runtime_control"]["mission_paused"] = False
    try:
        validate_active_config(active_config, tasks)
    except AssertionError as exc:
        _require(
            "mission_paused" in str(exc),
            f"active-config pause failed unexpectedly: {exc}",
        )
    else:
        raise AssertionError("an enabled v3 active config remained selectable")
