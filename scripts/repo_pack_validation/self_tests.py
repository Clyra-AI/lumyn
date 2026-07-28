"""Negative tests for the Lumyn v3.1 repo-pack validator."""

from __future__ import annotations

import copy
from collections.abc import Callable
from typing import Any

from repo_pack_validation.authority import (
    conditional_activation_scope_digest,
    manual_preflight_scope_digest,
    validate_authority_grants,
)
from repo_pack_validation.acceptance_text import validate_acceptance_text
from repo_pack_validation.markdown_refs import _markdown_anchors
from repo_pack_validation.task_contracts import (
    M2_IMPLEMENTATION_MARKER_REF,
    M2_PR_LIFECYCLE_REF,
    M2_REVIEW_REF,
    M2_SCORECARD_REF,
    M2_SCOPE_CLOSURE_REF,
    M2_VALIDATION_REPORT_REF,
    M2_VALIDATION_SUMMARY_REF,
    QUALIFYING_SAME_RUN_EVIDENCE_FIELDS,
    _policy_digest,
    validate_delegated_route_refs,
)


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


def _expect_evidence_failure(
    base: Payload,
    evidence_payloads: dict[str, dict[str, Any]],
    mutate: Callable[[dict[str, dict[str, Any]]], Any],
    expected: str,
    validate_loaded: ValidateLoaded,
) -> None:
    candidate = copy.deepcopy(evidence_payloads)
    mutate(candidate)
    try:
        validate_loaded(
            copy.deepcopy(base),
            validate_configs=False,
            evidence_payloads=candidate,
        )
    except AssertionError as exc:
        _require(
            expected.lower() in str(exc).lower(),
            f"evidence self-test expected {expected!r}, got {exc!r}",
        )
        return
    raise AssertionError(f"evidence self-test mutation did not fail: {expected}")


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
    if target_task_id == "M1":
        target["validation_contract_inheritance"][
            "acceptance_item_ids"
        ].append(acceptance_item_id)
        target["acceptance_result_requirements"].append(
            {
                "acceptance_item_id": acceptance_item_id,
                "allowed_statuses": ["partial", "missing", "blocked"],
                "evidence_mode": "automated",
                "closure_evidence": "validation_ref",
                "evidence_required": ["validation_report"],
            }
        )


def _remove_holdout_prohibition(payload: Payload, field: str) -> None:
    policy = _task(payload, "M1")["holdout_policy"]
    policy["prohibited_committed_fields"].remove(field)
    policy["policy_digest"] = _policy_digest(policy)


def _change_holdout_baseline(payload: Payload, value: str) -> None:
    policy = _task(payload, "M1")["holdout_policy"]
    policy["comparison_baseline"] = value
    policy["policy_digest"] = _policy_digest(policy)


def _remove_holdout_control(payload: Payload, field: str) -> None:
    policy = _task(payload, "M1")["holdout_policy"]
    policy["comparison_control_variables"].remove(field)
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
    evidence_payloads: dict[str, dict[str, Any]],
    validate_loaded: ValidateLoaded,
    validate_config_payload: Callable[..., None],
    validate_active_config: Callable[[dict[str, Any], dict[str, dict[str, Any]]], None],
    historical_plan_rel: str,
    expected_capabilities: dict[str, set[str]],
) -> None:
    """Prove drift cannot enable dispatch, authority, false closure, or weak proof."""

    tasks = validate_loaded(
        base,
        validate_configs=False,
        evidence_payloads=evidence_payloads,
    )
    _require(
        _markdown_anchors(
            "# API [Trust](https://example.com)\n"
            "# `Code` *Mode*\n"
        )
        == {"api-trust", "code-mode"},
        "Markdown heading slugs must use rendered inline text",
    )
    try:
        validate_acceptance_text(
            "1. `TEST-001`: native-\n configuration\n---\n",
            {
                "items": [
                    {
                        "acceptance_item_id": "TEST-001",
                        "source_text": "native- configuration",
                    }
                ]
            },
        )
    except AssertionError as exc:
        _require(
            "hard-wrap fracture" in str(exc),
            f"hard-wrap self-test failed unexpectedly: {exc}",
        )
    else:
        raise AssertionError("hard-wrapped acceptance text remained valid")

    delegated_tasks = copy.deepcopy(tasks)
    delegated_tasks["M10"]["product_action_route_contract"][
        "impact_read_only"
    ]["required_capabilities"].append("command_execution")
    try:
        validate_delegated_route_refs(delegated_tasks)
    except AssertionError as exc:
        _require(
            "differs from M4/impact_read_only" in str(exc),
            f"delegated-route self-test failed unexpectedly: {exc}",
        )
    else:
        raise AssertionError("stale copied M10 route remained valid")

    mutations: list[tuple[Callable[[Payload], Any], str]] = [
        (
            lambda value: _task(value, "M1").pop("baseline_commands"),
            "M1 parent control",
        ),
        (
            lambda value: _task(value, "M1")[
                "acceptance_result_requirements"
            ][0]["allowed_statuses"].append("implemented"),
            "must not claim terminal acceptance closure",
        ),
        (
            lambda value: value["m1_implementation"].pop(
                "mission_contract_ref"
            ),
            "M1 attended implementation packet missing",
        ),
        (
            lambda value: value["m1_implementation"]["lifecycle_gates"].__setitem__(
                "commit_push_required", True
            ),
            "authorize local validation only",
        ),
        (
            lambda value: value["m1_implementation"].__setitem__(
                "lifecycle_evidence_required", ["ship_packet"]
            ),
            "must not select lifecycle or holdout work",
        ),
        (
            lambda value: value["m1_implementation"]["factory_contract_binding"].__setitem__(
                "task_packet_schema_digest", "sha256:" + "0" * 64
            ),
            "Factory schema or semantic binding drifted",
        ),
        (
            lambda value: value["plan"]["alignment_gate"].__setitem__(
                "authorized_attended_task_refs", []
            ),
            "authorize only attended M1",
        ),
        (
            lambda value: value["m1_implementation"]["validation_contract_inheritance"][
                "acceptance_criteria"
            ].__setitem__(0, {}),
            "inheritance must be complete and local-only",
        ),
        (
            lambda value: value["m1_implementation"].pop("blocked_by"),
            "M1 attended implementation packet missing",
        ),
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
            "implemented evidence is semantically incomplete",
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
            lambda value: value["contract"][
                "conditional_acceptance_rules"
            ].__setitem__(
                "INSTALL-001",
                "requires exact Agent Runner route for every installation",
            ),
            "conditional Agent Runner setup",
        ),
        (
            lambda value: value["contract"][
                "conditional_acceptance_rules"
            ].__setitem__(
                "PILOT-003",
                "requires only one composed installed-preauthorization draft PR",
            ),
            "both agent and composed-delivery proofs",
        ),
        (
            lambda value: value["contract"][
                "model_control_requirements"
            ].pop("managed_credential"),
            "model controls missing managed_credential",
        ),
        (
            lambda value: value["contract"][
                "authority_requirements"
            ].__setitem__(
                "campaign_aggregate_union_authorized_per_installation",
                True,
            ),
            "exact per-action route authority",
        ),
        (
            lambda value: value["contract"][
                "authority_requirements"
            ].__setitem__("aggregate_composed_action_route_allowed", True),
            "exact per-action route authority",
        ),
        (
            lambda value: value["contract"][
                "model_control_requirements"
            ].pop("repository_command_isolation"),
            "model controls missing repository_command_isolation",
        ),
        (
            lambda value: value["contract"][
                "model_control_requirements"
            ].pop("sandbox_entrypoint_isolation"),
            "model controls missing sandbox_entrypoint_isolation",
        ),
        (
            lambda value: value["contract"][
                "verification_requirements"
            ].__setitem__("separate_verifier_process", False),
            "separate_verifier_process must be true",
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
            ].append("M0"),
            "cannot retain completed task scope",
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
            lambda value: _task(value, "M2.5")[
                "consumer_runner_prequalification"
            ].__setitem__("secret_values_collected", True),
            "prequalify one feasible",
        ),
        (
            lambda value: _task(value, "M2.5")[
                "consumer_runner_prequalification"
            ].__setitem__(
                "lumyn_adapter_conformance_required_at_prequalification",
                True,
            ),
            "premature conformance",
        ),
        (
            lambda value: _task(value, "M2.5")[
                "consumer_runner_prequalification"
            ].__setitem__(
                "plausible_organic_agent_item_hypothesis_required",
                False,
            ),
            "plausible organic agent-eligible item",
        ),
        (
            lambda value: _task(value, "M2.5")[
                "readiness_sprint_contract"
            ].__setitem__("may_close_disc_001", True),
            "readiness sprint must remain paid discovery",
        ),
        (
            lambda value: _task(value, "M2.5")[
                "readiness_sprint_contract"
            ].__setitem__(
                "credited_funds_count_only_after_signed_campaign_conversion",
                False,
            ),
            "readiness sprint must remain paid discovery",
        ),
        (
            lambda value: _task(value, "M2.5")[
                "campaign_offer_contract"
            ].__setitem__("minimum_tested_reviewable_outcomes", 5),
            "campaign offer must align",
        ),
        (
            lambda value: _task(value, "M2.5")[
                "consumer_delivery_prequalification"
            ].__setitem__(
                "minimum_provider_status_transmit_willing_consumers", 0
            ),
            "draft-PR and provider-status willingness",
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
            lambda value: _remove_holdout_control(
                value, "agent_runner_executable_digest"
            ),
            "baseline must isolate Lumyn",
        ),
        (
            lambda value: _task(value, "M1")[
                "walking_skeleton_contract"
            ]["stages"].remove("consumer_installation"),
            "walking skeleton stages",
        ),
        (
            lambda value: _task(value, "M1")[
                "walking_skeleton_contract"
            ].__setitem__("live_agent_execution_allowed", True),
            "remain offline",
        ),
        (
            lambda value: _task(value, "M1")[
                "conditional_factory_capabilities"
            ].append("approval"),
            "conditional Factory capabilities",
        ),
        (
            lambda value: _task(value, "M2")[
                "update_channel_contract"
            ].__setitem__("event_may_widen_installation", True),
            "Consumer Installation",
        ),
        (
            lambda value: _task(value, "M2")[
                "update_channel_contract"
            ]["provider_channel_transport"].__setitem__(
                "detached_signature_required", False
            ),
            "provider-channel transport",
        ),
        (
            lambda value: _task(value, "M2")[
                "update_channel_contract"
            ]["provider_channel_transport"].__setitem__(
                "contract_retrieved_bytes_digest_verified", False
            ),
            "provider-channel transport",
        ),
        (
            lambda value: _task(value, "M2")[
                "update_channel_contract"
            ].__setitem__("stored_github_token_allowed", True),
            "installation authorization mode",
        ),
        (
            lambda value: _task(value, "M2")[
                "update_channel_contract"
            ].__setitem__("provider_may_select_agent_runner", True),
            "Agent Runner selection",
        ),
        (
            lambda value: _task(value, "M2")[
                "update_channel_contract"
            ]["consumer_installation_fields"].remove(
                "agent_execution_policy"
            ),
            "Consumer Installation scope",
        ),
        (
            lambda value: _task(value, "M2")[
                "update_channel_contract"
            ].__setitem__("agent_execution_policy_default", "configured"),
            "Consumer Installation scope",
        ),
        (
            lambda value: _task(value, "M2")[
                "update_channel_contract"
            ]["provider_status"].__setitem__("merge_is_not_retired", False),
            "status projection",
        ),
        (
            lambda value: _task(value, "M2")[
                "update_channel_contract"
            ]["provider_status"].__setitem__(
                "raw_consumer_data_never_provider_visible", False
            ),
            "status projection",
        ),
        (
            lambda value: _task(value, "M2")[
                "managed_credential_contract"
            ].__setitem__("refresh_allowed", True),
            "managed credential broker contract",
        ),
        (
            lambda value: _task(value, "M4")[
                "product_action_route_contract"
            ]["impact_read_only"]["required_capabilities"].remove(
                "customer_repo_read"
            ),
            "routes must cover its exact product capability universe",
        ),
        (
            lambda value: _task(value, "M8").__setitem__(
                "product_action_route_contract", {}
            ),
            "must define at least one exact product action route",
        ),
        (
            lambda value: _task(value, "M9")[
                "product_action_route_contract"
            ]["local_export"]["required_capabilities"].append(
                "github_pr_write"
            ),
            "exact product action route",
        ),
        (
            lambda value: _task(value, "M9")[
                "product_action_route_contract"
            ]["remote_branch_push"][
                "conditionally_selected_capabilities"
            ].append(
                "github_branch_write"
            ),
            "route capability classes must be disjoint",
        ),
        (
            lambda value: _task(value, "M9")[
                "delivery_route_composition_contract"
            ].__setitem__(
                "aggregate_cross_action_scope_union_authorized",
                True,
            ),
            "atomic delegated actions",
        ),
        (
            lambda value: _task(value, "M10")[
                "campaign_route_composition_contract"
            ]["delegated_route_refs"].remove("M9/provider_status_transmit"),
            "compose exact delegated routes",
        ),
        (
            lambda value: _task(value, "M10")[
                "campaign_route_composition_contract"
            ]["required_same_run_route_sequence"].remove(
                "M9/provider_status_transmit"
            ),
            "one bound run",
        ),
        (
            lambda value: _task(value, "M10")[
                "campaign_route_composition_contract"
            ]["qualifying_same_run_evidence_binding"].__setitem__(
                "cross_run_evidence_allowed", True
            ),
            "same-run evidence binding",
        ),
        (
            lambda value: _task(value, "M10")[
                "paid_campaign_contract"
            ].__setitem__(
                "qualifying_same_run_evidence_binding_ref",
                "unbound",
            ),
            "reference the qualifying same-run",
        ),
        (
            lambda value: _task(value, "M6")["bounded_agent_contract"][
                "exact_fields"
            ].remove("cost_budget"),
            "model control fields",
        ),
        (
            lambda value: _task(value, "M6")["bounded_agent_contract"][
                "exact_fields"
            ].remove("agent_runner_adapter"),
            "model control fields",
        ),
        (
            lambda value: _task(value, "M6")["bounded_agent_contract"][
                "exact_fields"
            ].remove("agent_runner_executable_digest"),
            "model control fields",
        ),
        (
            lambda value: _task(value, "M6")["bounded_agent_contract"][
                "exact_fields"
            ].remove("agent_execution_policy"),
            "model control fields",
        ),
        (
            lambda value: _task(value, "M6")["bounded_agent_contract"].__setitem__(
                "clean_session_required", False
            ),
            "clean_session_required",
        ),
        (
            lambda value: _task(value, "M6")["bounded_agent_contract"].__setitem__(
                "silent_fallback_allowed", True
            ),
            "silent_fallback_allowed",
        ),
        (
            lambda value: _task(value, "M6")["bounded_agent_contract"].__setitem__(
                "repository_path_shadowing_allowed", True
            ),
            "repository_path_shadowing_allowed",
        ),
        (
            lambda value: _task(value, "M6")["bounded_agent_contract"][
                "untrusted_inputs"
            ].remove("repository_source"),
            "untrusted agent inputs",
        ),
        (
            lambda value: _task(value, "M6")[
                "runner_host_isolation_contract"
            ].__setitem__("ambient_service_sockets_allowed", True),
            "runner host-isolation contract",
        ),
        (
            lambda value: _task(value, "M6")[
                "runner_host_isolation_contract"
            ]["hard_resource_quota_fields"].remove("pids"),
            "runner host-isolation contract",
        ),
        (
            lambda value: _task(value, "M6")[
                "product_action_route_contract"
            ]["agent_assisted_candidate"][
                "authorization_topology_contract"
            ]["minimum_capability_sets"]["runner_mediated"].remove(
                "model_request_disclosure"
            ),
            "agent authorization topology",
        ),
        (
            lambda value: _task(value, "M6")[
                "managed_credential_contract"
            ].__setitem__("maximum_ttl_seconds", 86400),
            "managed credential broker contract",
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
            ].__setitem__("agent_runner_and_model_credentials_absent", False),
            "agent_runner_and_model_credentials_absent",
        ),
        (
            lambda value: _task(value, "M7")[
                "repository_command_isolation_contract"
            ].__setitem__("network_default", "enabled"),
            "repository-command isolation contract",
        ),
        (
            lambda value: _task(value, "M8")[
                "repository_command_isolation_contract"
            ].__setitem__(
                "agent_runner_model_and_sandbox_credentials_absent",
                False,
            ),
            "repository-command isolation contract",
        ),
        (
            lambda value: _task(value, "M8")[
                "sandbox_entrypoint_isolation_contract"
            ].__setitem__(
                "credential_injection_mode",
                "ambient_sandbox_credential",
            ),
            "sandbox-entrypoint isolation contract",
        ),
        (
            lambda value: _task(value, "M8")[
                "sandbox_entrypoint_isolation_contract"
            ]["backend_identity_fields"].remove("version"),
            "sandbox-entrypoint isolation contract",
        ),
        (
            lambda value: _task(value, "M7")[
                "deterministic_verification_contract"
            ]["repair_authorization"].__setitem__(
                "repair_route", "deterministic_or_agent_assisted"
            ),
            "repair must be agent-assisted",
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
            ]["repair_authorization"].__setitem__(
                "route_change_requires_new_explicit_authorization", False
            ),
            "route_change_requires_new_explicit_authorization",
        ),
        (
            lambda value: _task(value, "M7")[
                "deterministic_verification_contract"
            ]["candidate_modes"].remove("imported_manual"),
            "imported manual candidates",
        ),
        (
            lambda value: _task(value, "M7")["allowed_paths"].remove(
                "internal/authorization/"
            ),
            "M7 must own the future internal/authorization path",
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
            lambda value: _task(value, "M9")["delivery_contract"].__setitem__(
                "manual_delivery_cannot_close_exp_003", False
            ),
            "cannot waive automated draft-PR proof",
        ),
        (
            lambda value: _task(value, "M9")["delivery_contract"].__setitem__(
                "long_lived_token_allowed", True
            ),
            "short-lived",
        ),
        (
            lambda value: _task(value, "M9")["delivery_contract"][
                "provider_status_contract"
            ].__setitem__("silence_is_unknown", False),
            "provider status",
        ),
        (
            lambda value: _task(value, "M9")[
                "composed_update_contract"
            ].__setitem__("standalone_pr_create_qualifies", True),
            "composed provider-event-to-draft-PR proof",
        ),
        (
            lambda value: _task(value, "M9")[
                "composed_update_contract"
            ].__setitem__("pilot_same_run_provider_projection_required", False),
            "status projection must stay bound",
        ),
        (
            lambda value: _task(value, "M9")[
                "product_authority_requirements"
            ].append("provider_reporting"),
            "product authority requirements",
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
            ].__setitem__("api_provider_agent_access", True),
            "without provider access",
        ),
        (
            lambda value: _task(value, "M10")[
                "paid_campaign_contract"
            ].__setitem__("provider_raw_consumer_data_access", True),
            "without provider access",
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
            lambda value: _task(value, "M10")[
                "paid_campaign_contract"
            ].__setitem__("minimum_lumyn_opened_tested_draft_prs", 0),
            "composed Lumyn-generated",
        ),
        (
            lambda value: _task(value, "M10")[
                "paid_campaign_contract"
            ].__setitem__("minimum_provider_status_projections", 0),
            "real provider status projection",
        ),
        (
            lambda value: _task(value, "M10")[
                "paid_campaign_contract"
            ].__setitem__(
                "provider_status_projection_must_bind_composed_pr", False
            ),
            "real provider status projection",
        ),
        (
            lambda value: _task(value, "M10")[
                "paid_campaign_contract"
            ].__setitem__(
                "qualifying_composed_pr_requires_organic_agent_assisted_plan_item",
                False,
            ),
            "organic agent-assisted item",
        ),
        (
            lambda value: _task(value, "M10")[
                "paid_campaign_contract"
            ]["campaign_verdict_values"].append("reframe"),
            "verdict must be pass/fail",
        ),
        (
            lambda value: _task(value, "M10")[
                "paid_campaign_contract"
            ].__setitem__("campaign_economics_threshold_pass_required", False),
            "contribution-margin or automation threshold",
        ),
        (
            lambda value: _task(value, "M10")[
                "paid_campaign_contract"
            ].__setitem__("minimum_real_agent_assisted_outcomes", 0),
            "real consumer-selected",
        ),
        (
            lambda value: _task(value, "M10")[
                "paid_campaign_contract"
            ].__setitem__("unmatched_engine_comparison_is_causal", True),
            "generic-agent baseline must be fair",
        ),
        (
            lambda value: _task(value, "M3")["factoryd_runtime"].__setitem__(
                "dispatch_enabled", True
            ),
            "dispatch must be disabled",
        ),
        (
            lambda value: _task(value, "M2")[
                "validation_contract_inheritance"
            ]["required_review"].__setitem__("review_type", "code"),
            "M2 review type must remain security",
        ),
        (
            lambda value: _task(value, "M2").__setitem__(
                "risk_class", "medium"
            ),
            "M2 risk class must remain high",
        ),
        (
            lambda value: _task(value, "M2")[
                "validation_contract_inheritance"
            ]["required_review"].__setitem__(
                "reviewer_class", "validation_gate_or_review"
            ),
            "M2 reviewer class must remain independent_or_human",
        ),
        (
            lambda value: value["contract"]["required_review"].__setitem__(
                "review_type", "code"
            ),
            "validation contract must require the security review lens",
        ),
        (
            lambda value: value["contract"]["required_review"].__setitem__(
                "reviewer_class", "validation_gate_or_review"
            ),
            "validation contract reviewer class must remain independent_or_human",
        ),
        (
            lambda value: value["plan"]["alignment_gate"].__setitem__(
                "implementation_may_start", False
            ),
            "authorize only attended M1",
        ),
        (
            lambda value: value["mapping"].__setitem__(
                "generated_at", "2026-07-28T15:13:04Z"
            ),
            "exact regenerated M1 control set",
        ),
        (
            lambda value: value["plan"]["alignment_gate"][
                "completed_task_refs"
            ].remove("M2"),
            "preserve M0/M2 closure",
        ),
        (
            lambda value: value["risk"]["current_contract_state"].__setitem__(
                "M2", "implementation_and_local_validation_complete_lifecycle_pending"
            ),
            "preserve M0/M2 closure",
        ),
        (
            lambda value: _ledger_item(value, "TRUST-001")[
                "evidence_refs"
            ].remove(
                ".factory/artifacts/pr-lifecycle/lumyn-v3-m2/delivery-debt-record.json"
            ),
            "implemented evidence is semantically incomplete",
        ),
        (
            lambda value: next(
                group
                for group in value["mapping"]["groups"]
                if group.get("group_id") == "trust_and_privacy"
            )["implementation_progress"]["lifecycle_pending_task_refs"].append(
                "M2"
            ),
            "trust group must bind terminal M2 evidence",
        ),
        (
            lambda value: _task(value, "M2")["execution_state"].__setitem__(
                "state",
                "implementation_and_local_validation_complete_lifecycle_pending",
            ),
            "M2 task packet must bind closed lifecycle evidence",
        ),
        (
            lambda value: value["contract"]["m2_contract_implementation"].__setitem__(
                "state",
                "implementation_and_local_validation_complete_lifecycle_pending",
            ),
            "validation contract must bind closed M2 contract evidence",
        ),
        (
            lambda value: _closure_item(value, "TRUST-002").__setitem__(
                "implemented_task_refs", ["M0"]
            ),
            "TRUST-002 implemented closure must bind to M2",
        ),
        (
            lambda value: _closure_item(value, "EVENT-001")[
                "remaining_task_refs"
            ].append("M2"),
            "retain every unclosed active v3 task",
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
    for field in QUALIFYING_SAME_RUN_EVIDENCE_FIELDS:
        mutations.append(
            (
                lambda value, field=field: _task(value, "M10")[
                    "campaign_route_composition_contract"
                ]["qualifying_same_run_evidence_binding"][
                    "binding_fields"
                ].remove(field),
                "same-run evidence binding",
            )
        )
    mutations.append(
        (
            lambda value: value["plan"]["locked_decisions"].append(
                value["plan"]["locked_decisions"][0]
            ),
            "locked decisions must be unique",
        )
    )
    for mutate, expected in mutations:
        _expect_failure(base, mutate, expected, validate_loaded)

    evidence_mutations = [
        (
            lambda value: value[M2_SCORECARD_REF].__setitem__("task_id", "M1"),
            "proof scorecard must bind task M2",
        ),
        (
            lambda value: value[M2_IMPLEMENTATION_MARKER_REF].__setitem__(
                "execution_status", "fail"
            ),
            "work-proof marker is not passing",
        ),
        (
            lambda value: value[M2_VALIDATION_REPORT_REF].__setitem__(
                "candidate_digest", "sha256:" + "0" * 64
            ),
            "validation report candidate binding drifted",
        ),
        (
            lambda value: value[M2_VALIDATION_SUMMARY_REF].__setitem__(
                "status", "fail"
            ),
            "validation summary binding or passing status drifted",
        ),
        (
            lambda value: value[M2_REVIEW_REF]["current_work"].__setitem__(
                "validation_run_id", "validation:wrong-run"
            ),
            "independent review is not approved, resolved, and current-work bound",
        ),
        (
            lambda value: value[M2_PR_LIFECYCLE_REF].__setitem__(
                "status", "incomplete"
            ),
            "PR lifecycle report is incomplete",
        ),
        (
            lambda value: value[M2_SCOPE_CLOSURE_REF].__setitem__(
                "overall_status", "closed"
            ),
            "scope-closure report is incomplete",
        ),
    ]
    for mutate, expected in evidence_mutations:
        _expect_evidence_failure(
            base,
            evidence_payloads,
            mutate,
            expected,
            validate_loaded,
        )

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

    for task_id, action_mode in (("M6", "bounded_agent"),):
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
        "M6",
        "bounded_agent",
        list(tasks["M6"]["conditional_factory_capabilities"]),
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
