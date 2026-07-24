"""Negative tests for the Lumyn repo-pack validator."""

from __future__ import annotations

import copy
from collections.abc import Callable
from typing import Any

from repo_pack_validation.authority import validate_authority_grants


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


def run_repo_pack_self_tests(
    base: Payload,
    *,
    validate_loaded: ValidateLoaded,
    validate_config_payload: Callable[..., None],
    validate_active_config: Callable[[dict[str, Any], dict[str, dict[str, Any]]], None],
    historical_plan_rel: str,
    expected_capabilities: dict[str, set[str]],
) -> None:
    """Prove the validator rejects drift in authority, evidence, and scope."""

    tasks = validate_loaded(base, validate_configs=False)
    mutations: list[tuple[Callable[[Payload], Any], str]] = [
        (
            lambda value: value["ledger"]["items"].pop(),
            "compiled acceptance text differs from PRD",
        ),
        (
            lambda value: value["packets"]["tasks"][9][
                "required_worker_chain"
            ].__setitem__(2, "ship-pr"),
            "worker chain",
        ),
        (
            lambda value: value["packets"]["tasks"][0][
                "acceptance_item_ids"
            ].remove("REB-001"),
            "inherited acceptance IDs",
        ),
        (
            lambda value: value["packets"]["tasks"][9].__setitem__(
                "auto_merge", True
            ),
            "forbidden product capability",
        ),
        (
            lambda value: value["packets"]["tasks"][9]["factoryd_runtime"][
                "capability_grants"
            ][0].__setitem__("approved", True),
            "planning-time capability grants",
        ),
        (
            lambda value: value["mapping"]["groups"][0][
                "acceptance_item_ids"
            ].pop(),
            "exact PRD acceptance set",
        ),
        (
            lambda value: next(
                task
                for task in value["packets"]["tasks"]
                if task["task_id"] == "M4"
            )["allowed_paths"].append(
                ".factory/artifacts/lifecycle-evidence/M4/"
            ),
            "implementation writes to lifecycle-owned evidence",
        ),
        (
            lambda value: next(
                task
                for task in value["packets"]["tasks"]
                if task["task_id"] == "M7"
            )["validation_contract_inheritance"]["required_review"].__setitem__(
                "review_type", "security"
            ),
            "inherited review requirement",
        ),
        (
            lambda value: next(
                task
                for task in value["packets"]["tasks"]
                if task["task_id"] == "M9"
            )["blocked_by"].append("M8"),
            "independently deliverable",
        ),
        (
            lambda value: next(
                task
                for task in value["packets"]["tasks"]
                if task["task_id"] == "M1"
            )["holdout_suite_policy"]["prohibited_committed_fields"].remove(
                "answer_key"
            ),
            "resolving provenance and answer material",
        ),
        (
            lambda value: next(
                task
                for task in value["packets"]["tasks"]
                if task["task_id"] == "M1"
            )["holdout_suite_policy"]["prohibited_committed_fields"].remove(
                "source_url"
            ),
            "resolving provenance and answer material",
        ),
        (
            lambda value: next(
                task
                for task in value["packets"]["tasks"]
                if task["task_id"] == "M1"
            )["holdout_policy"].__setitem__("mode", "evaluate"),
            "without a fabricated pre-existing commitment",
        ),
        (
            lambda value: next(
                task
                for task in value["packets"]["tasks"]
                if task["task_id"] == "M1"
            )["lifecycle_gates"].update(
                {
                    "holdout_provisioning_required": False,
                    "holdout_evaluation_required": True,
                }
            ),
            "must not claim current-candidate holdout evaluation",
        ),
        (
            lambda value: next(
                task
                for task in value["packets"]["tasks"]
                if task["task_id"] == "M4"
            )["holdout_policy"].__setitem__(
                "provisioning_result_ref",
                ".factory/artifacts/lifecycle-evidence/M4/holdout-result.json",
            ),
            "independently provisioned M1 holdout result",
        ),
        (
            lambda value: next(
                task
                for task in value["packets"]["tasks"]
                if task["task_id"] == "M10"
            )["blocked_by"].append("M8"),
            "without forcing sandbox proof",
        ),
        (
            lambda value: next(
                task
                for task in value["packets"]["tasks"]
                if task["task_id"] == "M10"
            ).__setitem__(
                "holdout_policy",
                copy.deepcopy(
                    next(
                        task
                        for task in value["packets"]["tasks"]
                        if task["task_id"] == "M4"
                    )["holdout_policy"]
                ),
            ),
            "must not declare holdout policy without holdout-evaluator",
        ),
        (
            lambda value: next(
                task
                for task in value["packets"]["tasks"]
                if task["task_id"] == "M10"
            )["product_authority_requirements"].remove("campaign_receipt"),
            "product authority requirements differ",
        ),
        (
            lambda value: next(
                task
                for task in value["packets"]["tasks"]
                if task["task_id"] == "M8"
            )["optional_product_action_capabilities"].remove(
                "provider_trust_status_read"
            ),
            "optional product authority requirements differ",
        ),
        (
            lambda value: next(
                task
                for task in value["packets"]["tasks"]
                if task["task_id"] == "M2.5"
            )["manual_external_evidence_preflight"].__setitem__(
                "approved_private_storage_boundary",
                "a changed boundary",
            ),
            "scope digest must match its canonical content",
        ),
        (
            lambda value: next(
                task
                for task in value["packets"]["tasks"]
                if task["task_id"] == "M2.5"
            )["manual_external_evidence_preflight"].__setitem__(
                "allowed_private_fields",
                [
                    field
                    for field in next(
                        task
                        for task in value["packets"]["tasks"]
                        if task["task_id"] == "M2.5"
                    )["manual_external_evidence_preflight"][
                        "allowed_private_fields"
                    ]
                    if "provider-authenticated consumer signer binding"
                    not in field
                ],
            ),
            "preflight must bind provider-authenticated consumer signer binding",
        ),
        (
            lambda value: value["plan"].__setitem__(
                "rollback_or_deletion_test",
                [
                    "Revoke everything and verify no public or provider-facing copy survives."
                ],
            ),
            "irreversible external-disclosure semantics",
        ),
    ]
    for mutate, expected in mutations:
        _expect_failure(base, mutate, expected, validate_loaded)

    historical_config = copy.deepcopy(base["config"])
    historical_config["repos"]["lumyn"][
        "task_packets"
    ] = f"{historical_plan_rel}/task-packets.json"
    try:
        validate_config_payload(
            historical_config, "self-test historical config", autoship=False
        )
    except AssertionError as exc:
        _require(
            "task_packets ref is stale" in str(exc),
            f"historical-plan self-test failed unexpectedly: {exc}",
        )
    else:
        raise AssertionError("historical-plan config remained selectable")

    active_config = copy.deepcopy(base["config"])
    active_config["repos"]["lumyn"]["capability_grants"] = [
        {
            "task_id": "M9",
            "capability": "approval",
            "approved": True,
            "evidence_ref": "private:authorizations/M9.json",
            "expires_at": "2099-01-01T00:00:00Z",
        }
    ]
    try:
        validate_active_config(active_config, tasks)
    except AssertionError as exc:
        _require(
            "missing exact Factory capabilities" in str(exc),
            f"active-grant self-test failed unexpectedly: {exc}",
        )
    else:
        raise AssertionError("partial live capability grant remained selectable")

    product_config = copy.deepcopy(base["config"])
    product_config["repos"]["lumyn"]["capability_grants"] = [
        {
            "task_id": "M9",
            "capability": "github_pr_write",
            "approved": True,
            "evidence_ref": "private:authorizations/M9.json",
            "expires_at": "2099-01-01T00:00:00Z",
        }
    ]
    try:
        validate_active_config(product_config, tasks)
    except AssertionError as exc:
        _require(
            "unknown or product capability" in str(exc),
            f"product-capability separation self-test failed unexpectedly: {exc}",
        )
    else:
        raise AssertionError(
            "product authority was accepted as a Factory grant"
        )

    duplicate_config = copy.deepcopy(base["config"])
    duplicate_grant = {
        "task_id": "M2.5",
        "capability": "approval",
        "approved": False,
    }
    duplicate_config["repos"]["lumyn"]["capability_grants"] = [
        duplicate_grant,
        copy.deepcopy(duplicate_grant),
    ]
    try:
        validate_active_config(duplicate_config, tasks)
    except AssertionError as exc:
        _require(
            "duplicates grant" in str(exc),
            f"duplicate-grant self-test failed unexpectedly: {exc}",
        )
    else:
        raise AssertionError("duplicate Factory grant remained selectable")

    preflight = tasks["M2.5"]["manual_external_evidence_preflight"]
    wrong_scope_grant = [
        {
            "task_id": "M2.5",
            "capability": "approval",
            "approved": True,
            "evidence_ref": preflight["approval_evidence_ref"],
            "approval_scope_digest": "sha256:" + ("0" * 64),
            "expires_at": "2099-01-01T00:00:00Z",
        }
    ]
    try:
        validate_authority_grants(
            wrong_scope_grant, tasks, expected_capabilities
        )
    except AssertionError as exc:
        _require(
            "exactly bind the current manual preflight scope" in str(exc),
            f"approval-scope self-test failed unexpectedly: {exc}",
        )
    else:
        raise AssertionError(
            "a generic or stale M2.5 approval remained selectable"
        )

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
            f"semantic-wildcard self-test failed unexpectedly: {exc}",
        )
    else:
        raise AssertionError(
            "a semantic-wildcard network grant remained selectable"
        )

    duplicate_network_grants = copy.deepcopy(broad_network_grants)
    duplicate_network_grants[-1]["network_allowlist"] = [
        "api.stripe.com:443",
        "API.STRIPE.COM:443",
    ]
    try:
        validate_authority_grants(
            duplicate_network_grants, tasks, expected_capabilities
        )
    except AssertionError as exc:
        _require(
            "case-insensitive duplicates" in str(exc),
            f"network-duplicate self-test failed unexpectedly: {exc}",
        )
    else:
        raise AssertionError(
            "case-insensitive duplicate endpoints remained selectable"
        )

    unspecified_network_grants = copy.deepcopy(broad_network_grants)
    unspecified_network_grants[-1]["network_allowlist"] = ["0.0.0.0"]
    try:
        validate_authority_grants(
            unspecified_network_grants, tasks, expected_capabilities
        )
    except AssertionError as exc:
        _require(
            "unspecified address" in str(exc),
            f"unspecified-address self-test failed unexpectedly: {exc}",
        )
    else:
        raise AssertionError(
            "an unspecified network address remained selectable"
        )

    wildcard_network_grants = copy.deepcopy(broad_network_grants)
    wildcard_network_grants[-1]["network_allowlist"] = ["*.stripe.com:443"]
    try:
        validate_authority_grants(
            wildcard_network_grants, tasks, expected_capabilities
        )
    except AssertionError as exc:
        _require(
            "wildcard-free" in str(exc),
            f"wildcard-host self-test failed unexpectedly: {exc}",
        )
    else:
        raise AssertionError("a wildcard network host remained selectable")
