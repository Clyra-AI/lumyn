"""Validate Lumyn-specific task ownership and external-evidence closure."""

from __future__ import annotations

import hashlib
import json
from typing import Any

from repo_pack_validation.authority import manual_preflight_scope_digest

CANONICAL_WORKER_ORDER = [
    "task-executor",
    "validation-gate",
    "code-review",
    "holdout-evaluator",
    "trace-grader",
    "evidence-attestor",
    "commit-push",
    "post-merge-monitor",
]


def _require(condition: bool, message: str) -> None:
    if not condition:
        raise AssertionError(message)


def _commands(task: dict[str, Any], field: str) -> str:
    value = task.get(field)
    _require(isinstance(value, list) and value, f"{task.get('task_id')}.{field} must be non-empty")
    return json.dumps(value)


def _holdout_policy_digest(policy: dict[str, Any]) -> str:
    canonical = json.dumps(
        {key: value for key, value in policy.items() if key != "policy_digest"},
        sort_keys=True,
        separators=(",", ":"),
    ).encode()
    return f"sha256:{hashlib.sha256(canonical).hexdigest()}"


def _validate_attested_signal(task: dict[str, Any], validator: str) -> None:
    task_id = str(task["task_id"])
    commands = _commands(task, "validation_commands")
    final_commands = _commands(task, "final_validation_commands")
    for blob, label in [(commands, "validation"), (final_commands, "final validation")]:
        _require(validator in blob and "--attestation" in blob and "--public-manifest" in blob, f"{task_id} {label} must validate actual private and public evidence")
    attestation = task.get("evidence_attestation")
    _require(isinstance(attestation, dict) and attestation.get("required") is True, f"{task_id} evidence attestation must be required")
    _require(attestation.get("required_worker") == "evidence-attestor", f"{task_id} must require evidence-attestor")
    _require(attestation.get("failure_behavior") == "block_commit_push_and_acceptance_closure", f"{task_id} attestation failure must block shipping and closure")
    _require(task.get("lifecycle_gates", {}).get("evidence_attestation_required") is True, f"{task_id} lifecycle must gate on evidence attestation")
    for field in ["evidence_required", "worker_evidence_required"]:
        _require("attestation_record" not in task.get(field, []), f"{task_id}.{field} must not own attestation_record")
    _require("attestation_record" in task.get("lifecycle_evidence_required", []), f"{task_id} lifecycle must own attestation_record")
    _require(task.get("lifecycle_evidence_refs", {}).get("attestation_record") == attestation.get("output_ref"), f"{task_id} attestation output and shipping ref must match")
    _require("evidence-attestor" in task.get("required_worker_chain", []), f"{task_id} chain must include evidence-attestor")


def _validate_product_authority(task: dict[str, Any]) -> None:
    task_id = str(task["task_id"])
    expected = task.get("product_authority_requirements")
    _require(isinstance(expected, list) and expected, f"{task_id} product authority must be explicit")
    command_token = f"authorization validate --bundle \\\"$LUMYN_PRIVATE_PRODUCT_AUTHORITY_BUNDLE\\\" --task {task_id}"
    for field in ["validation_commands", "final_validation_commands"]:
        _require(command_token in _commands(task, field), f"{task_id}.{field} must enforce its private product-authority bundle")
    gate = task.get("product_live_action_gate")
    _require(isinstance(gate, dict) and gate.get("required") is True, f"{task_id} product live-action gate must be required")
    _require(gate.get("exact_capabilities") == expected, f"{task_id} product authority gate must preserve exact capabilities")
    _require(gate.get("failure_behavior") == "block_product_live_action", f"{task_id} missing product authority must block the live action")
    _require(gate.get("cached_authority_decision_allowed") is False, f"{task_id} must revalidate authority at each live action")
    _require("immediately before every product live side effect" in gate.get("enforcement_point", ""), f"{task_id} must enforce authority at the side-effect boundary")
    _require("Factory dispatch does not validate" in gate.get("factory_dispatch_semantics", ""), f"{task_id} must not misstate Factory authority enforcement")
    _require("product_authority_validation_required" not in task.get("lifecycle_gates", {}), f"{task_id} must not invent a Factory product-authority gate")


def _validate_independent_workers(task: dict[str, Any], expected: list[str]) -> None:
    task_id = str(task["task_id"])
    chain = task.get("required_worker_chain")
    _require(isinstance(chain, list), f"{task_id} required worker chain must be a list")
    ordered = [worker for worker in CANONICAL_WORKER_ORDER if worker in chain]
    _require(chain == ordered, f"{task_id} required worker chain must use canonical lifecycle order")
    for worker in expected:
        _require(worker in chain, f"{task_id} must include {worker}")
    evidence_by_worker = {
        "holdout-evaluator": ("holdout_result", "holdout-result.json"),
        "trace-grader": ("trace_grade_report", "trace-grade-report.json"),
        "evidence-attestor": ("attestation_record", "attestation-record.json"),
    }
    for worker in expected:
        evidence, filename = evidence_by_worker[worker]
        _require(evidence in task.get("lifecycle_evidence_required", []), f"{task_id} must require lifecycle-owned {evidence}")
        ref = task.get("lifecycle_evidence_refs", {}).get(evidence)
        expected_root = f".factory/artifacts/lifecycle-evidence/{task_id}/"
        _require(ref == f"{expected_root}{filename}", f"{task_id} must declare the exact trusted {evidence} ref")
        _require(expected_root in task.get("forbidden_paths", []), f"{task_id} must forbid implementation writes to lifecycle evidence")
        _require(
            all(not str(path).startswith(expected_root) for path in task.get("allowed_paths", [])),
            f"{task_id} implementation paths must not include lifecycle evidence",
        )


def _validate_manual_external_evidence_preflight(task: dict[str, Any]) -> None:
    preflight = task.get("manual_external_evidence_preflight")
    _require(isinstance(preflight, dict) and preflight.get("required") is True, "M2.5 manual external-evidence preflight must be required")
    _require(preflight.get("factory_approval_capability") == "approval", "M2.5 preflight must use the canonical Factory approval capability")
    _require(preflight.get("product_runtime_authority") is False, "M2.5 preflight must not masquerade as product runtime authority")
    _require(preflight.get("failure_behavior") == "block_external_evidence_collection", "M2.5 preflight failure must block collection")
    for field in [
        "allowed_private_fields",
        "public_fields",
        "approved_private_storage_boundary",
        "retention_ttl_and_expiry_required",
        "deletion_on_revocation_required",
        "deletion_receipt_and_orphan_owner_required",
        "public_disclosure_requires_separate_consent",
        "participant_consent_required",
    ]:
        _require(preflight.get(field), f"M2.5 preflight missing {field}")
    private_scope = " ".join(preflight["allowed_private_fields"])
    for token in [
        "provider-status channel",
        "receipt-acknowledgement key and exchange",
        "OS/architecture and fail-closed host-isolation compatibility",
        "provider-authenticated consumer signer binding",
        "provider-signed acknowledgement",
        "one-invitation-unit cardinality",
    ]:
        _require(token in private_scope, f"M2.5 preflight must bind {token}")
    _require(
        "irreversibility acknowledgement" in private_scope,
        "M2.5 preflight must bind the irreversible external-disclosure acknowledgement",
    )
    grants = task.get("factoryd_runtime", {}).get("capability_grants", [])
    approval = [grant for grant in grants if grant.get("capability") == "approval"]
    _require(len(approval) == 1 and approval[0].get("evidence_ref") == preflight.get("approval_evidence_ref"), "M2.5 Factory approval must cite the manual preflight")
    digest = preflight.get("approval_scope_digest")
    _require(
        isinstance(digest, str) and digest.startswith("sha256:") and len(digest) == 71,
        "M2.5 preflight must carry a canonical SHA-256 approval scope digest",
    )
    expected_digest = manual_preflight_scope_digest(preflight)
    _require(digest == expected_digest, "M2.5 preflight scope digest must match its canonical content")
    _require(approval[0].get("approval_scope_digest") == digest, "M2.5 Factory approval must bind the exact manual preflight scope digest")


def validate_migration_task_contracts(tasks: dict[str, dict[str, Any]]) -> None:
    """Validate task-specific ownership that generic Factory schemas cannot infer."""

    _require({"REB-001", "REB-002"}.issubset(tasks["M0"]["acceptance_item_ids"]), "M0 must fix false success and provider terminology")
    _require(".factory/artifacts/exceptions/architecture-debt-lumyn-migration-rebaseline.json" in tasks["M0"]["upstream_debt_refs"], "M0 must own the active architecture-debt route")
    _require({"scripts/validate_repo_pack.py", "scripts/repo_pack_validation/"}.issubset(tasks["M0"]["allowed_paths"]), "M0 must be allowed to decompose repo-pack validation")
    _require({"CORPUS-001", "CORPUS-002", "CORPUS-003"} == set(tasks["M1"]["acceptance_item_ids"]), "M1 corpus scope drifted")
    policy = tasks["M1"].get("holdout_suite_policy", {})
    _require(policy.get("private_root_env") == "LUMYN_HOLDOUT_ROOT", "M1 holdout root must be evaluator-controlled")
    _require(policy.get("task_executor_access") == "forbidden", "M1 task executor must not access holdouts")
    _require(policy.get("evaluator_access") == "holdout-evaluator-only", "M1 holdouts must be evaluator-only")
    _require(policy.get("provisioning_owner") == "independent_holdout_owner", "M1 holdout must have an independent provisioning owner")
    _require(policy.get("provisioning_worker") == "holdout-evaluator", "M1 holdout provisioning must use holdout-evaluator")
    _require(
        policy.get("freeze_result_ref") == ".factory/artifacts/lifecycle-evidence/M1/holdout-result.json",
        "M1 holdout freeze must bind the lifecycle-owned holdout result",
    )
    _require(policy.get("frozen_suite_commitment_must_match_result") is True, "M1 frozen suite commitment must match its result")
    _require(
        set(policy.get("committed_fields", []))
        == {
            "opaque_case_ids",
            "non_resolving_provenance_class",
            "license_posture",
            "frozen_suite_commitment",
            "encrypted_or_hmac_artifact_commitments",
        },
        "M1 committed holdout fields must be opaque and non-resolving",
    )
    _require(
        {
            "inputs",
            "source",
            "source_url",
            "repository_url",
            "package_name",
            "answer_key",
            "expected_labels",
            "expected_patches",
            "raw_traces",
            "plaintext_content_digest",
            "resolvable_provenance",
        }
        .issubset(set(policy.get("prohibited_committed_fields", []))),
        "M1 committed holdout manifest must prohibit resolving provenance and answer material",
    )
    holdout = tasks["M1"].get("holdout_evaluation", {})
    _require(holdout.get("purpose") == "provision_freeze_and_integrity_validate", "M1 holdout result must prove provisioning and freeze")
    _require(holdout.get("frozen_suite_commitment_required") is True, "M1 holdout result must bind a frozen suite commitment")
    provision_policy = tasks["M1"].get("holdout_policy", {})
    _require(
        provision_policy.get("mode") == "provision"
        and provision_policy.get("suite_namespace") == "private://lumyn-migration-mvp/holdouts"
        and provision_policy.get("commitment_algorithm") == "hmac-sha256",
        "M1 must independently provision its holdout without a fabricated pre-existing commitment",
    )
    _require(
        provision_policy.get("policy_digest") == _holdout_policy_digest(provision_policy),
        "M1 holdout provisioning policy digest drifted",
    )
    _require(
        tasks["M1"]["lifecycle_gates"].get("holdout_provisioning_required") is True
        and tasks["M1"]["lifecycle_gates"].get("holdout_evaluation_required") is False,
        "M1 provisioning evidence must not claim current-candidate holdout evaluation",
    )
    for task_id in ["M4", "M6", "M7"]:
        evaluation_policy = tasks[task_id].get("holdout_policy", {})
        _require(
            evaluation_policy.get("mode") == "evaluate"
            and evaluation_policy.get("provisioning_result_ref")
            == ".factory/artifacts/lifecycle-evidence/M1/holdout-result.json",
            f"{task_id} must evaluate the independently provisioned M1 holdout result",
        )
        _require(
            evaluation_policy.get("policy_digest")
            == _holdout_policy_digest(evaluation_policy),
            f"{task_id} holdout evaluation policy digest drifted",
        )
        _require(
            tasks[task_id]["lifecycle_gates"].get("holdout_provisioning_required")
            is False
            and tasks[task_id]["lifecycle_gates"].get("holdout_evaluation_required")
            is True,
            f"{task_id} must evaluate, not provision, the frozen holdout",
        )
    for task_id, task in tasks.items():
        if task_id not in {"M1", "M4", "M6", "M7"}:
            _require(
                "holdout_policy" not in task,
                f"{task_id} must not declare holdout policy without holdout-evaluator",
            )
            _require(
                task["lifecycle_gates"].get("holdout_provisioning_required") is False,
                f"{task_id} must not require holdout provisioning",
            )
    trust_items = {"CHG-001", "CHG-002", "CHG-003", "CHG-004"}
    _require((trust_items | {"AUTH-001", "AUTH-002", "AUTH-003", "AUTH-004", "AUTH-005", "ACT-001", "ACT-002"}).issubset(tasks["M2"]["acceptance_item_ids"]), "M2 trust and activation contracts are incomplete")
    _require(
        {"internal/authorization/", "internal/isolation/", "internal/receipt/", "internal/attestation/"}.issubset(
            tasks["M2"]["allowed_paths"]
        ),
        "M2 must own authorization, host-isolation, receipt, and richer-attestation contracts",
    )
    _require(trust_items.issubset(tasks["M3"]["acceptance_item_ids"]) and "internal/trust/" in tasks["M3"]["allowed_paths"], "M3 must own executable packet trust enforcement")
    _require(set(tasks["M2.5"]["acceptance_item_ids"]) == {"DISC-001", "DISC-002"}, "M2.5 design-partner qualification drifted")
    _require(set(tasks["M2.5"]["blocked_by"]) == {"M0", "M2"}, "M2.5 must wait for the standard campaign and invitation contracts")
    _validate_manual_external_evidence_preflight(tasks["M2.5"])
    _require({f"IMP-{number:03d}" for number in range(1, 6)}.issubset(tasks["M4"]["acceptance_item_ids"]), "M4 impact scope drifted")
    _require({"PLAN-001", "PLAN-002"}.issubset(tasks["M5"]["acceptance_item_ids"]) and "M2.5" in tasks["M5"]["blocked_by"], "M5 plan scope or qualification dependency drifted")
    _require(
        "consumer-signed product-authority bundle" in tasks["M5"]["objective"]
        and "cannot self-approve" in tasks["M5"]["objective"],
        "M5 must consume explicit issuance and cannot mint its own authority",
    )
    gates = tasks["M5"].get("gated_by_acceptance_items", [])
    _require(
        {(gate.get("acceptance_item_id"), gate.get("required_status")) for gate in gates}
        == {("DISC-001", "implemented"), ("DISC-002", "implemented")},
        "M5 must carry runner-enforced DISC-001 and DISC-002 acceptance gates",
    )
    _require({f"PATCH-{number:03d}" for number in range(1, 7)}.issubset(tasks["M6"]["acceptance_item_ids"]), "M6 patch scope drifted")
    _require("internal/isolation/" in tasks["M6"]["allowed_paths"], "M6 package and repository commands must use host isolation")
    for task_id in ["M6", "M8", "M9"]:
        _require({"CHG-002", "CHG-004"}.issubset(tasks[task_id]["acceptance_item_ids"]), f"{task_id} must revalidate packet lifecycle and trust at write time")
        _require("internal/trust/" in tasks[task_id]["allowed_paths"], f"{task_id} must include the trust boundary")
    _require({"VER-001", "VER-002", "VER-003", "VER-005", "VER-006", "VER-007", "EVD-001", "ACT-003"}.issubset(tasks["M7"]["acceptance_item_ids"]), "M7 verification, trace, and offline-canary scope drifted")
    _require({"cmd/lumyn/", "internal/isolation/", "internal/retention/", "workflows/", "cassettes/", "runs/"}.issubset(tasks["M7"]["allowed_paths"]), "M7 runtime, isolation, retention, and fixture paths drifted")
    _require("explicit authorization issuance" in tasks["M7"]["objective"], "M7 offline canary must issue explicit synthetic authority")
    _require("internal/isolation/" in tasks["M8"]["allowed_paths"], "M8 sandbox entrypoint must use host isolation")
    _require(set(tasks["M9"]["blocked_by"]) == {"M7"} and tasks["M9"].get("conditional_blocked_by") == [], "M9 must be independently deliverable after offline verification")
    _require("provider_attestation" not in tasks["M9"].get("product_authority_requirements", []), "M9 draft PR delivery must not require provider reporting")
    _require("provider_attestation" in tasks["M9"].get("optional_product_action_capabilities", []), "M9 provider reporting must remain optional and separately authorized")
    for task_id in ["M8", "M9", "M10"]:
        _require(
            "provider_trust_status_read"
            in tasks[task_id].get("optional_product_action_capabilities", []),
            f"{task_id} must allow an exact status read without requiring it over a signed offline snapshot",
        )
    _require({f"PILOT-{number:03d}" for number in range(1, 10)}.issubset(tasks["M10"]["acceptance_item_ids"]), "M10 pilot and recurring-economics gates drifted")
    activation = next(
        check for check in tasks["M10"]["acceptance_checks"] if check.startswith("PILOT-003:")
    )
    for required_text in [
        "explicitly issue repository authorization",
        "provider-signed acknowledgement",
        "authenticated consumer signer",
        "eligible-repository invitation unit",
        "invitation receipt",
        "seven days",
        "security, privacy, platform, and maintainer",
        "two hours",
        "in-product",
        "60 minutes",
    ]:
        _require(required_text in activation, f"M10 consumer activation criterion missing {required_text}")
    _require(set(tasks["M10"]["blocked_by"]) == {"M2.5", "M9"}, "M10 must require qualification and delivery without forcing sandbox proof")
    _require(
        "campaign_receipt" in tasks["M10"].get("product_authority_requirements", []),
        "M10 sponsored-program counting must require the minimal campaign-receipt grant",
    )
    _require(
        {"provider_trust_status_read", "sandbox_network", "sandbox_credential", "sandbox_request_disclosure", "provider_attestation"}
        == set(tasks["M10"].get("optional_product_action_capabilities", [])),
        "M10 status-read alternative, sandbox, and richer reporting capabilities must remain optional and action-specific",
    )
    _validate_attested_signal(tasks["M2.5"], "validate_design_partner_evidence.py")
    _validate_attested_signal(tasks["M10"], "validate_pilot_evidence.py")
    for task_id, workers in {
        "M1": ["holdout-evaluator"],
        "M2.5": ["evidence-attestor"],
        "M4": ["holdout-evaluator"],
        "M6": ["holdout-evaluator"],
        "M7": ["holdout-evaluator", "trace-grader"],
        "M10": ["evidence-attestor"],
    }.items():
        _validate_independent_workers(tasks[task_id], workers)
    for task in tasks.values():
        task_id = str(task["task_id"])
        _require(
            task.get("required_review") == task.get("validation_contract_inheritance", {}).get("required_review"),
            f"{task_id} task and inherited review requirements must be identical",
        )
        _require(f".factory/artifacts/pr-lifecycle/{task_id}/" in task.get("forbidden_paths", []), f"{task_id} must forbid implementation writes to PR lifecycle evidence")
        control = task.get("factoryd_runtime", {}).get("runtime_control", {})
        _require(control.get("launch_request", {}).get("requested_write_paths") == control.get("max_write_scope_paths"), f"{task['task_id']} launch write request must cover every bounded implementation path")
    for task_id in ["M8", "M9", "M10"]:
        _validate_product_authority(tasks[task_id])
