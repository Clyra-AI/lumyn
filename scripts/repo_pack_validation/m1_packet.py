"""Focused validation for the attended M1 implementation packet."""

from __future__ import annotations

import hashlib
from pathlib import Path
from typing import Any, Callable, Mapping

from repo_pack_validation import factory_schema_core, m1_evidence


M1_IMPLEMENTATION_REQUIRED_FIELDS = set(
    """mission_contract_ref shared_constraints_ref validation_contract_inheritance
    blocked_by worker_type allowed_paths forbidden_paths evidence_ownership architecture_target_paths
    generated_by generated_at planning_prerequisite_paths planning_prerequisite_owner
    planning_prerequisite_rule path_planning_method baseline_commands red_first_commands
    validation_commands final_validation_commands required_worker_chain lifecycle_gates
    evidence_required worker_evidence_required stop_conditions max_iterations retry_budget
    input_artifacts upstream_debt_refs commands acceptance_checks scope_exclusions runtime_pins
    alignment_gate_ref plan_drift_policy_ref factory_compatibility factory_contract_binding
    factoryd_runtime ci_lane_refs ci_control_refs test_matrix_refs coverage_policy_refs
    security_scanner_gates engineering_policy_refs architecture_guidance_refs
    structured_parser_policy_refs agent_native_cli_checks artifact_budget_refs
    redaction_posture required_proof_level proof_scorecard_required semantic_invariants
    changelog_intent versioning_impact migration_impact docs_sync_refs""".split()
)

VENDORED_TASK_PACKET_SCHEMA_REF = ".factory/contracts/factory/task-packet.schema.json"
VENDORED_RUNTIME_CONTROL_SCHEMA_REF = ".factory/contracts/factory/runtime-control.schema.json"
VENDORED_SCHEMA_VALIDATOR_REF = "scripts/repo_pack_validation/factory_schema_core.py"
FACTORY_CONTRACT_BINDING = {
    "task_packet_schema_ref": "factory://schemas/artifacts/task-packet.schema.json",
    "task_packet_schema_digest": "sha256:c46a8ff085f5dd0e88181cc70341bfa7ddcad890695a322499e9648dc34b93ab",
    "vendored_task_packet_schema_ref": VENDORED_TASK_PACKET_SCHEMA_REF,
    "runtime_control_schema_ref": "factory://schemas/artifacts/runtime-control.schema.json",
    "runtime_control_schema_digest": "sha256:288fa14a1e3e40b926b0cae19e9753c4415cd8da7c5eb8c1650f4008820f4d9a",
    "vendored_runtime_control_schema_ref": VENDORED_RUNTIME_CONTROL_SCHEMA_REF,
    "vendored_schema_validator_ref": VENDORED_SCHEMA_VALIDATOR_REF,
    "vendored_schema_validator_digest": "sha256:a0748a9d44bf967f7a73dc0026bef31db50b7b47dfc715d049492096319b1a86",
    "semantic_validator_ref": "factory://scripts/contract_validators/task_packet_shape.py",
    "semantic_validator_digest": "sha256:e5348377293ffaa73b7a2f31f0fe363bb3f21455cab1d270d2cbc773cdf3f63d",
}


def validate_m1_implementation_packet(
    packet: dict[str, Any],
    *,
    root: Path,
    plan_rel: str,
    artifact_refs: Mapping[str, str],
    expected_generation_at: str,
    require: Callable[[bool, str], None],
    fail: Callable[[str], None],
    load_json: Callable[[Path], dict[str, Any]],
    list_of_strings: Callable[..., bool],
    contains_machine_local_path: Callable[[Any], bool],
) -> None:
    """Validate the portable attended M1 implementation-only runner packet."""

    missing = M1_IMPLEMENTATION_REQUIRED_FIELDS - set(packet)
    require(not missing, f"M1 attended implementation packet missing {sorted(missing)}")
    require(
        packet.get("artifact_type") == "task_packet"
        and packet.get("artifact_version") == "1.0"
        and packet.get("task_id") == "M1-IMPLEMENTATION"
        and packet.get("worker_type") == "task-executor",
        "M1 attended implementation packet identity or worker drifted",
    )
    require(
        packet.get("generated_at") == expected_generation_at
        and packet.get("generated_by")
        == "prd-to-plan+execution-compiler+plan-verify-repair",
        "M1 attended implementation packet is outside the exact regenerated control set",
    )
    require(
        packet.get("blocked_by") == ["M0", "M2"],
        "M1 attended implementation dependencies must remain M0 and M2",
    )
    require(
        packet.get("mission_contract_ref") == artifact_refs["plan"]
        and packet.get("shared_constraints_ref") == artifact_refs["contract"]
        and packet.get("alignment_gate_ref") == f"{artifact_refs['plan']}#/alignment_gate"
        and packet.get("plan_drift_policy_ref") == f"{artifact_refs['plan']}#/plan_drift_policy",
        "M1 attended implementation packet control refs drifted",
    )
    inherited = packet.get("validation_contract_inheritance", {})
    require(
        inherited.get("source_validation_contract_ref") == artifact_refs["contract"]
        and isinstance(inherited.get("acceptance_criteria"), list)
        and len(inherited["acceptance_criteria"]) >= 5
        and all(isinstance(item, str) and item.strip() for item in inherited["acceptance_criteria"])
        and inherited.get("required_review")
        == {"required": False, "review_type": "code", "reviewer_class": "none"},
        "M1 attended implementation inheritance must be complete and local-only",
    )
    require(
        "acceptance_ledger_ref" not in packet
        and "acceptance_item_ids" not in packet
        and "acceptance_result_requirements" not in packet,
        "M1 attended implementation packet must not claim item-level acceptance",
    )
    require(
        packet.get("baseline_commands") == ["make lint-fast", "make test-fast"]
        and packet.get("red_first_commands")
        == [
            "LUMYN_M1_VERIFICATION_TARGET=baseline npm --prefix "
            "examples/consumer-repos/det-operation-rename test --silent"
        ]
        and packet.get("final_validation_commands") == ["make prepush-full"]
        and packet.get("max_iterations") == 2
        and packet.get("retry_budget") == 1,
        "M1 attended implementation command or retry contract drifted",
    )
    require(
        packet.get("required_worker_chain") == ["task-executor", "validation-gate"],
        "M1 attended implementation worker chain must remain local-only",
    )
    gates = packet.get("lifecycle_gates", {})
    lifecycle_fields = (
        "ci_required", "code_review_required", "codex_review_required",
        "holdout_provisioning_required", "holdout_evaluation_required",
        "trace_grading_required", "evidence_attestation_required",
        "commit_push_required", "post_merge_monitor_required",
        "pr_lifecycle_report_required",
    )
    require(
        gates.get("local_validation_required") is True
        and all(gates.get(field) is False for field in lifecycle_fields),
        "M1 attended implementation packet must authorize local validation only",
    )
    require(
        "lifecycle_evidence_required" not in packet
        and "lifecycle_evidence_refs" not in packet
        and "holdout_policy" not in packet,
        "M1 attended implementation packet must not select lifecycle or holdout work",
    )
    allowed = set(packet.get("allowed_paths", []))
    required_allowed = {
        "examples/migration-packs/", "examples/consumer-repos/",
        "examples/integration-graphs/", "examples/candidates/",
        "examples/negative/", "examples/verification/",
        "examples/holdout-manifest.json", ".factory/contracts/",
        ".github/action-ref-exceptions.yaml", ".github/workflows/validate.yml",
        ".tool-versions", "Makefile", "scripts/repo_pack_ci.py",
        "scripts/repo_pack_validation/", "tests/", "docs/dev/dev_guides.md",
        "docs/architecture/architecture_guides.md", "CHANGELOG.md",
    }
    require(
        required_allowed.issubset(allowed)
        and "docs/product/prd.md" not in allowed
        and "docs/product/plan.md" not in allowed
        and ".factory/artifacts/task-runs/M1-IMPLEMENTATION/" not in allowed
        and "docs/" not in allowed,
        "M1 task-executor paths must be exact and exclude validation evidence writes",
    )
    require(
        all(not Path(path).is_absolute() and ".." not in Path(path).parts for path in allowed),
        "M1 attended implementation paths must be repo-relative",
    )
    forbidden = set(packet.get("forbidden_paths", []))
    require(
        {
            ".git/", f"{plan_rel}/", ".factory/artifacts/lifecycle-evidence/M1/",
            ".factory/artifacts/pr-lifecycle/lumyn-v3-m1/",
        }.issubset(forbidden),
        "M1 attended implementation must forbid control and lifecycle writes",
    )
    planning_prerequisites = [
        "AGENTS.md", "WORKFLOW.md", "docs/product/prd.md", "docs/product/plan.md",
        f"{plan_rel}/",
    ]
    require(
        packet.get("planning_prerequisite_paths") == planning_prerequisites
        and packet.get("planning_prerequisite_owner") == "plan-verify-repair"
        and "task-executor cannot write them" in packet.get("planning_prerequisite_rule", "")
        and set(planning_prerequisites).isdisjoint(allowed),
        "M1 attended implementation planning prerequisites or ownership drifted",
    )
    require(
        not any(lane.get("lane") == "lifecycle-evidence" for lane in packet.get("ci_lane_refs", [])),
        "M1 attended implementation child must not own lifecycle-evidence execution",
    )
    m1_evidence.validate_evidence_ownership(packet, require)
    require(
        packet.get("requires_capabilities") == []
        and packet.get("conditional_factory_capabilities") == []
        and packet.get("product_authority_requirements") == []
        and packet.get("optional_product_action_capabilities") == []
        and packet.get("requires_credentials") is False
        and packet.get("requires_network") is False
        and packet.get("sovereignty_mode") == "local_only",
        "M1 attended implementation must remain offline and authority-free",
    )
    runtime = packet.get("factoryd_runtime", {})
    control = runtime.get("runtime_control", {})
    launch = control.get("launch_request", {})
    require(
        runtime.get("network_posture") == "offline"
        and runtime.get("sovereignty_mode") == "local_only"
        and runtime.get("capability_grants") == []
        and control.get("mission_paused") is True
        and control.get("max_token_budget") == 0
        and control.get("max_write_scope_paths") == []
        and launch.get("requested_action") == "continue_mission"
        and launch.get("expected_decision") == "block"
        and list_of_strings(launch.get("requested_write_paths")),
        "M1 attended implementation runtime must remain schema-valid and paused",
    )
    binding = packet.get("factory_contract_binding", {})
    require(
        all(binding.get(key) == value for key, value in FACTORY_CONTRACT_BINDING.items()),
        "M1 attended implementation Factory schema or semantic binding drifted",
    )
    for ref, digest_key in (
        (VENDORED_TASK_PACKET_SCHEMA_REF, "task_packet_schema_digest"),
        (VENDORED_RUNTIME_CONTROL_SCHEMA_REF, "runtime_control_schema_digest"),
        (VENDORED_SCHEMA_VALIDATOR_REF, "vendored_schema_validator_digest"),
    ):
        path = root / ref
        require(path.is_file(), f"M1 vendored Factory schema is missing: {ref}")
        actual = f"sha256:{hashlib.sha256(path.read_bytes()).hexdigest()}"
        require(actual == binding.get(digest_key), f"M1 vendored Factory schema digest drifted: {ref}")
    vendored_root = root / ".factory/contracts/factory"
    task_packet_schema = load_json(root / VENDORED_TASK_PACKET_SCHEMA_REF)
    factory_schema_core.validate_schema(
        task_packet_schema,
        packet,
        "M1 attended implementation packet",
        task_packet_schema,
        root=vendored_root,
        fail=fail,
        load_json=load_json,
        validation_error_type=AssertionError,
    )
    state = packet.get("execution_state", {})
    require(
        state.get("state") == "attended_local_implementation_authorized"
        and state.get("acceptance_closure_claimed") is False
        and state.get("lifecycle_authorized") is False
        and state.get("factoryd_authorized") is False,
        "M1 attended implementation state must not imply lifecycle or closure authority",
    )
    require(
        not contains_machine_local_path(packet),
        "M1 attended implementation packet contains a machine-local path",
    )
