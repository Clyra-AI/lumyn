#!/usr/bin/env python3
from __future__ import annotations

import json
from datetime import datetime, timezone
from pathlib import Path
from tempfile import TemporaryDirectory
from typing import Any

import validate_repo_pack as validator
from repo_pack_architecture import (
    ARCHITECTURE_BUDGET_EXCEPTION_REFS,
    FACTORY_ARCHITECTURE_BUDGET_POLICY_REF,
    architecture_budget_unexcepted_failures,
    architecture_debt_exception_expiry_error,
    validate_architecture_debt_exception_expiry,
)

ACCEPTANCE_LEDGER_REF = validator.ACCEPTANCE_LEDGER_REF
FACTORYD_REPO_KEY = validator.FACTORYD_REPO_KEY
REQUIRED_PLAN_SKILL_REFS = validator.REQUIRED_PLAN_SKILL_REFS
RUNTIME_CONTROL_FORBIDDEN_PATHS = validator.RUNTIME_CONTROL_FORBIDDEN_PATHS
contains_machine_local_path = validator.contains_machine_local_path
expected_task_version_slices = validator.expected_task_version_slices
fail = validator.fail
validate_factoryd_config = validator.validate_factoryd_config
validate_model_provider_gate = validator.validate_model_provider_gate
validate_standalone_task_packet = validator.validate_standalone_task_packet
validate_task_execution_compiler_fields = validator.validate_task_execution_compiler_fields
validate_task_packets = validator.validate_task_packets


def propagated_task(task_id_value: str, blocked_by: list[str]) -> dict[str, Any]:
    return {
        "task_id": task_id_value,
        "blocked_by": blocked_by,
        "mvp_required_version_slices": sorted(expected_task_version_slices(task_id_value)),
        "delivery_slice_refs": sorted(expected_task_version_slices(task_id_value)),
        "factory_compatibility": {
            "factory_contract_version": "1.0",
            "profile_ref": "profiles/lumyn.yaml",
            "skill_vocabulary_version": "2026-06-09",
            "skill_inventory_ref": "skills/README.md",
            "generated_by": "prd-to-plan+execution-compiler",
            "generated_at": "2026-06-09T00:00:00Z",
            "deprecated_worker_policy": "block_active_aliases",
            "deprecated_worker_aliases": [
                {
                    "deprecated": "ship-pr",
                    "replacement": "commit-push",
                    "status": "deprecated",
                    "migration_behavior": "block active task packets until required_worker_chain is migrated",
                }
            ],
        },
        "scope_exclusions": [
            "MCP recording",
            "event assertions",
            "hosted dashboard",
            "runtime enforcement",
        ],
        "alignment_gate_ref": ".factory/artifacts/prd-to-plan/lumyn-mvp/execution-plan.json#/alignment_gate",
        "plan_drift_policy_ref": ".factory/artifacts/prd-to-plan/lumyn-mvp/execution-plan.json#/plan_drift_policy",
        "acceptance_ledger_ref": ACCEPTANCE_LEDGER_REF,
        "acceptance_item_ids": ["RCRR-001"],
        "validation_contract_inheritance": {
            "acceptance_ledger_ref": ACCEPTANCE_LEDGER_REF,
            "acceptance_item_ids": ["RCRR-001"],
        },
        "required_worker_chain": [
            "task-executor",
            "validation-gate",
            "commit-push",
            "post-merge-monitor",
        ],
        "lifecycle_gates": {
            "local_validation_required": True,
            "ci_required": True,
            "code_review_required": False,
            "codex_review_required": True,
            "commit_push_required": True,
            "post_merge_monitor_required": True,
            "pr_lifecycle_report_required": True,
            "skip_policy": "approved_exception_required",
        },
        "allowed_paths": [
            "cmd/lumyn/",
            "internal/result/",
            "internal/source/",
            "schemas/",
            "tests/",
            "docs/",
            f".factory/artifacts/task-runs/{task_id_value}/",
            f".factory/artifacts/pr-lifecycle/{task_id_value}/",
        ],
        "architecture_target_paths": [
            "cmd/lumyn/",
            "internal/result/",
            "internal/source/",
        ],
        "path_planning_method": "architecture_target_paths_v1",
        "forbidden_paths": [".git/", ".factory/tmp/", ".factoryd/", *RUNTIME_CONTROL_FORBIDDEN_PATHS],
        "worker_type": "codex_cli",
        "factoryd_runtime": {
            "state_dir": ".factoryd/",
            "workspace_root": ".factoryd/workspaces/",
            "branch_prefix": "codex",
            "worker_type": "codex_cli",
            "worker_command": "",
            "approval_posture": "human approval required for live credentials, high-risk tasks, and merge",
            "credential_posture": "no ambient secrets during deterministic MVP bootstrap",
            "network_posture": "offline by default until live sandbox/model work is approved",
            "capability_grants": [],
        },
        "validation_commands": ["make prepush-full"],
        "max_iterations": 2,
        "evidence_required": ["validation_report", "work_proof_marker", "factoryd_run_once_report"],
        "required_proof_level": "source_evidence",
        "artifact_budget_refs": ["docs/dev/dev_guides.md#structured-data-proof-and-evidence-budgets"],
        "redaction_posture": {
            "classification": "internal",
            "customer_safe": False,
            "redaction_notes": "Self-test task-run evidence remains internal.",
            "recursive_policy": "nested owner, credential, secret, endpoint, and machine-local path fields must be redacted before sharing",
        },
        "stop_conditions": [
            "missing runner-ready task contract",
            "changed path outside allowed_paths",
            "forbidden path changed",
            "required validation command failed",
            "plan drift detected",
        ],
        "test_matrix_refs": [{"tier": "Tier 1 Unit", "source_ref": "docs/dev/dev_guides.md#12-level-test-matrix"}],
        "ci_lane_refs": [
            {
                "lane": "fast",
                "source_ref": "docs/dev/dev_guides.md#ci-lane-mapping",
                "command_refs": ["make lint-fast", "make test-fast"],
            },
            {
                "lane": "core",
                "source_ref": "docs/dev/dev_guides.md#ci-lane-mapping",
                "command_refs": ["make test-contracts"],
            },
            {
                "lane": "acceptance",
                "source_ref": "docs/dev/dev_guides.md#ci-lane-mapping",
                "command_refs": [".factory/artifacts/prd-to-plan/lumyn-mvp/scope-closure-map.json"],
            },
            {
                "lane": "cross_platform",
                "source_ref": "docs/dev/dev_guides.md#ci-lane-mapping",
                "exception_ref": ".factory/artifacts/exceptions/cross-platform-deferred.json",
            },
            {
                "lane": "risk",
                "source_ref": "docs/dev/dev_guides.md#ci-lane-mapping",
                "status_check_refs": ["CodeQL analyze"],
            },
            {
                "lane": "release",
                "source_ref": "docs/dev/dev_guides.md#ci-lane-mapping",
                "exception_ref": ".factory/artifacts/exceptions/release-deferred.json",
            },
        ],
        "coverage_policy_refs": {
            "required": True,
            "source_ref": "docs/dev/dev_guides.md#coverage-gates",
            "minimums": [
                {
                    "scope": "go_first_party_packages_overall",
                    "threshold": 75,
                    "command_refs": ["make test-coverage"],
                }
            ],
            "command_refs": ["make test-coverage"],
        },
        "security_scanner_gates": {
            "required": True,
            "scanner": "CodeQL",
            "workflow_ref": ".github/workflows/codeql.yml",
            "status_check": "CodeQL analyze",
        },
        "engineering_policy_refs": [
            {
                "policy": "docs_parity",
                "source_ref": "docs/dev/dev_guides.md#docs-parity",
                "rule": "docs move with behavior",
            },
            {
                "policy": "output_contracts",
                "source_ref": "docs/dev/dev_guides.md#output-contracts",
                "rule": "output contracts stay stable",
            },
            {
                "policy": "release_integrity",
                "source_ref": "docs/dev/dev_guides.md#release-integrity",
                "rule": "release integrity evidence is required",
            },
            {
                "policy": "provenance_evidence",
                "source_ref": "docs/dev/dev_guides.md#provenance-evidence",
                "rule": "provenance evidence stays repo-relative",
            },
        ],
        "architecture_guidance_refs": [
            {"source_ref": "docs/architecture/architecture_guides.md#systems-thinking-map", "rule": "record state and feedback"},
            {
                "source_ref": "docs/architecture/architecture_guides.md#tdd-and-red-first-expectations",
                "rule": "use TDD and red-first evidence",
            },
            {
                "source_ref": "docs/architecture/architecture_guides.md#adr-and-decision-triggers",
                "rule": "record ADR decision triggers",
            },
            {
                "source_ref": "docs/architecture/architecture_guides.md#performance-and-cost-triggers",
                "rule": "track performance and cost impact",
            },
            {
                "source_ref": "docs/architecture/architecture_guides.md#reliability-and-recovery-triggers",
                "rule": "record reliability and recovery evidence",
            },
            {
                "source_ref": "docs/architecture/architecture_guides.md#trust-mode-posture",
                "rule": "fail-closed trust-mode posture",
            },
        ],
        "planning_skill_refs": list(REQUIRED_PLAN_SKILL_REFS),
        "runtime_pins": {
            "language": "go",
            "go_version": "1.26.4",
            "toolchain_version": "1.26.4",
            "module_path": "github.com/Clyra-AI/lumyn",
            "module_or_package_path": "github.com/Clyra-AI/lumyn",
            "dependency_policy": "standard library first; pinned dependencies only when task-required",
            "distribution_target": "standalone_binary",
            "provider_policy": "OpenAI-compatible HTTP and Anthropic Messages HTTP adapters; no model key or network in deterministic bootstrap",
            "mvp_eval_providers": [
                "openai_compatible_http_adapter",
                "anthropic_messages_http_adapter",
            ],
            "artifact_namespace": ".factory/artifacts/ for Factory evidence",
            "live_work_policy": "blocked until deterministic replay foundation passes and human approval unlocks live work",
        },
        "slice_type": "vertical",
        "slice_rationale": {
            "slice_type": "vertical",
            "why_this_shape": "self-test task preserves a vertically scoped implementation contract",
        },
        "changelog": {
            "impact": "required_when_implemented",
            "section": "Unreleased",
            "draft_entry": "Self-test task changelog entry.",
            "semver_marker_override": "pre_1_0_minor_candidate",
        },
        "contract_impact": "Self-test task changes only its declared contract surface.",
        "versioning_migration_impact": "Pre-1.0 changes must preserve explicit migration notes before release.",
        "architecture_constraints": ["record state owner, feedback source, and fail-closed behavior"],
        "adr_required": True,
        "tdd_first_failing_tests": ["add a failing test or fixture before implementation when practical"],
        "cost_perf_impact": {"level": "low", "measurement_expectation": "no material cost increase expected"},
        "chaos_failure_hypothesis": {
            "hypothesis": "invalid evidence must fail closed",
            "expected_fail_closed_behavior": "do not mark the task complete",
        },
        "semantic_invariants": ["evidence remains repo-relative", "closure cannot claim missing PRD scope"],
    }


def model_provider_gate_task(task_id_value: str = "T11.1", grant_task_id: str = "*") -> dict[str, Any]:
    task = propagated_task(task_id_value, ["T11"])
    task.update(
        {
            "requires_model_provider_endpoint": True,
            "requires_human_approval": False,
            "model_provider_requirements": {
                "required_grant": "model_provider_endpoint",
                "provider_surfaces": ["openai_compatible_http", "anthropic_messages_http"],
                "required_fields": [
                    "provider_identity",
                    "provider_model",
                    "provider_endpoint_or_base_url",
                    "credential_environment",
                    "budget_posture",
                    "redaction_posture",
                    "network_allowlist",
                ],
            },
        }
    )
    task["factoryd_runtime"]["capability_grants"] = [
        {
            "task_id": grant_task_id,
            "capability": "model_provider_endpoint",
            "approved": False,
            "evidence_ref": ".factory/artifacts/approvals/model_provider_endpoint.md",
            "network_allowlist": ["pending-approved-provider-host"],
            "provider_identity": "pending-approved-provider",
            "provider_model": "pending-approved-model",
            "provider_endpoint": "pending-approved-provider-endpoint",
            "credential_environment": "pending-approved-credential-environment",
            "budget_posture": "pending-approved-budget",
            "redaction_posture": "pending-approved-redaction",
        }
    ]
    task["stop_conditions"].append("missing model_provider_endpoint grant")
    return task


def run_self_test() -> int:
    valid_packets = {"tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]}
    validate_task_packets(valid_packets, "T2.6")
    with TemporaryDirectory() as temp_dir:
        temp_root = Path(temp_dir)
        oversized = temp_root / "internal" / "source" / "new.py"
        oversized.parent.mkdir(parents=True)
        oversized.write_text("line\n" * 2501, encoding="utf-8")
        sample_budget = {
            "source_extensions": [".py"],
            "excluded_dirs": [".git", ".factoryd", ".factory/tmp", "workspaces"],
            "fail_line_threshold": 2500,
        }
        failures = architecture_budget_unexcepted_failures(temp_root, sample_budget, set(), {})
        if not failures or "internal/source/new.py" not in failures[0]:
            fail("architecture budget self-test expected unexcepted oversized source to fail")
        scratch = temp_root / ".factory" / "tmp" / "scratch.py"
        scratch.parent.mkdir(parents=True)
        scratch.write_text("line\n" * 2501, encoding="utf-8")
        if any(".factory/tmp/scratch.py" in failure for failure in architecture_budget_unexcepted_failures(temp_root, sample_budget, set(), {})):
            fail("architecture budget self-test expected .factory/tmp scratch to be excluded")
        workspace_source = temp_root / "workspaces" / "demo" / "generated.py"
        workspace_source.parent.mkdir(parents=True)
        workspace_source.write_text("line\n" * 2501, encoding="utf-8")
        if any("workspaces/demo/generated.py" in failure for failure in architecture_budget_unexcepted_failures(temp_root, sample_budget, set(), {})):
            fail("architecture budget self-test expected workspaces scratch to be excluded")
        if architecture_budget_unexcepted_failures(
            temp_root,
            sample_budget,
            {"internal/source/new.py"},
            {"internal/source/new.py": 2501},
        ):
            fail("architecture budget self-test expected exception-scoped source to pass")
        ceiling_failures = architecture_budget_unexcepted_failures(
            temp_root,
            sample_budget,
            {"internal/source/new.py"},
            {"internal/source/new.py": 2500},
        )
        if not ceiling_failures or "approved ceiling" not in ceiling_failures[0]:
            fail("architecture budget self-test expected exception growth over ceiling to fail")
        prefix_root = temp_root / "prefix-check"
        first_party_build = prefix_root / "internal" / "build" / "big.py"
        first_party_build.parent.mkdir(parents=True)
        first_party_build.write_text("line\n" * 2501, encoding="utf-8")
        generated_build = prefix_root / "build" / "generated.py"
        generated_build.parent.mkdir(parents=True)
        generated_build.write_text("line\n" * 2501, encoding="utf-8")
        nested_dependency = prefix_root / "packages" / "web" / "node_modules" / "dep" / "generated.py"
        nested_dependency.parent.mkdir(parents=True)
        nested_dependency.write_text("line\n" * 2501, encoding="utf-8")
        prefix_budget = {**sample_budget, "excluded_dirs": [*sample_budget["excluded_dirs"], "build", "node_modules"]}
        prefix_failures = architecture_budget_unexcepted_failures(prefix_root, prefix_budget, set(), {})
        if not any("internal/build/big.py" in failure for failure in prefix_failures):
            fail("architecture budget self-test expected first-party internal/build source to fail")
        if any("build/generated.py" in failure for failure in prefix_failures):
            fail("architecture budget self-test expected root build output to be excluded")
        if any("node_modules/dep/generated.py" in failure for failure in prefix_failures):
            fail("architecture budget self-test expected nested node_modules dependency to be excluded")
    validate_architecture_debt_exception_expiry(
        "self-test-valid",
        {"expires_at": "2099-01-01T00:00:00Z"},
        datetime(2026, 7, 1, tzinfo=timezone.utc),
    )
    if not architecture_debt_exception_expiry_error(
        "self-test-expired",
        {"expires_at": "2026-06-30T00:00:00Z"},
        datetime(2026, 7, 1, tzinfo=timezone.utc),
    ):
        fail("architecture debt exception self-test expected expired evidence to fail")
    try:
        validate_task_packets(
            {
                "artifact_type": "task_packet_set",
                "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])],
            },
            "T2.6",
        )
    except AssertionError as exc:
        if "artifact_type" not in str(exc):
            raise
    else:
        fail("self-test expected legacy task packet set artifact_type to fail")
    if expected_task_version_slices("T11.1") != {"v0.2"}:
        fail("self-test expected dotted live-eval task to inherit parent v0.2 delivery slice")
    if expected_task_version_slices("T2.5"):
        fail("self-test expected T2.5 lifecycle baseline task to remain delivery-slice exempt")
    if contains_machine_local_path(".factory/artifacts/prd-to-plan/lumyn-mvp/execution-plan.json#/alignment_gate"):
        fail("self-test expected repo-relative JSON pointer to remain portable")

    validate_model_provider_gate(model_provider_gate_task())

    original_config_grants = validator.factoryd_config_capability_grants
    original_example_config = validator.FACTORYD_CONFIG
    original_active_config = validator.FACTORYD_ACTIVE_CONFIG
    original_autoship_config = validator.FACTORYD_AUTOSHIP_CONFIG
    active_config_grant = {
        "task_id": "T11.1",
        "capability": "model_provider_endpoint",
        "approved": True,
        "evidence_ref": ".factory/artifacts/approvals/model_provider_endpoint.md",
        "network_allowlist": ["api.example.com"],
        "provider_identity": "example-provider",
        "provider_model": "example-model",
        "provider_endpoint": "https://api.example.com/v1",
        "credential_environment": "LUMYN_PROVIDER_API_KEY",
        "budget_posture": "capped",
        "redaction_posture": "redacted",
    }
    with TemporaryDirectory() as temp_dir:
        temp_root = Path(temp_dir)
        example_config = temp_root / "factoryd.example.json"
        active_config = temp_root / "factoryd.json"
        autoship_config = temp_root / "factoryd.autoship.example.json"
        config_payload = {"repos": {FACTORYD_REPO_KEY: {"capability_grants": [active_config_grant]}}}
        empty_active_payload = {"repos": {FACTORYD_REPO_KEY: {"capability_grants": []}}}
        example_config.write_text(json.dumps(config_payload), encoding="utf-8")
        autoship_config.write_text(json.dumps(config_payload), encoding="utf-8")
        active_config.write_text(json.dumps(empty_active_payload), encoding="utf-8")
        try:
            validator.FACTORYD_CONFIG = example_config
            validator.FACTORYD_ACTIVE_CONFIG = active_config
            validator.FACTORYD_AUTOSHIP_CONFIG = autoship_config
            if validator.factoryd_config_capability_grants():
                fail("self-test expected example and autoship config grants to be ignored")
            active_config.write_text(json.dumps(config_payload), encoding="utf-8")
            if validator.factoryd_config_capability_grants() != [active_config_grant]:
                fail("self-test expected active factoryd.json grants to be visible")
        finally:
            validator.FACTORYD_CONFIG = original_example_config
            validator.FACTORYD_ACTIVE_CONFIG = original_active_config
            validator.FACTORYD_AUTOSHIP_CONFIG = original_autoship_config

    base_runtime = {
        "repo_path": "..",
        "acceptance_ledger": ACCEPTANCE_LEDGER_REF,
        "task_packets": ".factory/artifacts/prd-to-plan/lumyn-mvp/task-packets.json",
        "scope_closure_map": ".factory/artifacts/prd-to-plan/lumyn-mvp/scope-closure-map.json",
        "validation_contract": ".factory/artifacts/prd-to-plan/lumyn-mvp/validation-contract.json",
        "state_dir": "../.factoryd",
        "workspace_root": "../.factoryd/workspaces",
        "validation_commands": [
            "python3 scripts/validate_repo_pack.py --self-test",
            "python3 scripts/validate_repo_pack.py",
        ],
        "branch_prefix": "codex",
        "worker_type": "codex_cli",
        "worker_command": "",
        "approval_posture": "human approval required for live credentials, high-risk tasks, and merge",
        "credential_posture": "no ambient secrets during deterministic MVP bootstrap",
        "network_posture": "offline by default until live sandbox/model work is approved",
        "capability_grants": [],
        "architecture_budget": {
            "enabled": True,
            "policy_ref": FACTORY_ARCHITECTURE_BUDGET_POLICY_REF,
            "warn_line_threshold": 1200,
            "fail_line_threshold": 2500,
            "source_extensions": [".go", ".py", ".ts", ".tsx", ".js", ".jsx"],
            "excluded_dirs": [".git", ".factoryd", ".factory/tmp", "workspaces", ".venv", "__pycache__", "build", "dist", "node_modules", "vendor"],
            "exception_refs": ARCHITECTURE_BUDGET_EXCEPTION_REFS,
        },
        "auto_ship": False,
        "shipping": {
            "enabled": False,
            "push_required": False,
            "pr_required": False,
            "ci_required": False,
            "codex_review_required": False,
            "merge_required": False,
            "post_merge_required": False,
            "scope_closure_required": False,
            "push_command": "",
            "open_pr_command": "",
            "ci_command": "",
            "codex_review_command": "",
            "merge_command": "",
            "post_merge_command": "",
            "scope_closure_command": "",
        },
    }
    base_config = {
        "factory": {"repo_path": "${FACTORY_REPO}", "profile_path": "profiles/lumyn.yaml"},
        "repos": {"lumyn": base_runtime},
    }
    active_runtime = json.loads(json.dumps(base_runtime))
    active_runtime["capability_grants"] = [active_config_grant]
    active_config_with_grant = {
        "factory": {"repo_path": "${FACTORY_REPO}", "profile_path": "profiles/lumyn.yaml"},
        "repos": {"lumyn": active_runtime},
    }
    autoship_runtime = json.loads(json.dumps(base_runtime))
    autoship_runtime["auto_ship"] = True
    autoship_runtime["shipping"].update(
        {
            "enabled": True,
            "provider": "github_cli",
            "push_required": True,
            "pr_required": True,
            "ci_required": True,
            "codex_review_required": True,
            "merge_required": True,
            "post_merge_required": True,
            "scope_closure_required": True,
            "scope_closure_mode": "semantic",
        }
    )
    autoship_config = {
        "factory": {"repo_path": "${FACTORY_REPO}", "profile_path": "profiles/lumyn.yaml"},
        "repos": {"lumyn": autoship_runtime},
    }
    validate_factoryd_config(base_config, active_config_with_grant, autoship_config)
    bad_active_config = json.loads(json.dumps(active_config_with_grant))
    bad_active_config["repos"]["lumyn"]["validation_commands"] = ["go test ./..."]
    try:
        validate_factoryd_config(base_config, bad_active_config, autoship_config)
    except AssertionError as exc:
        if "must run validate_repo_pack.py" not in str(exc):
            raise
    else:
        fail("self-test expected active config without repo-pack validation to fail")

    active_wildcard_task = model_provider_gate_task()
    active_wildcard_task["factoryd_runtime"]["capability_grants"] = []
    active_wildcard_grant = {
        "task_id": "*",
        "capability": "model_provider_endpoint",
        "approved": True,
        "evidence_ref": ".factory/artifacts/approvals/model_provider_endpoint.md",
        "network_allowlist": ["api.example.com"],
        "provider_identity": "example-provider",
        "provider_model": "example-model",
        "provider_endpoint": "https://api.example.com/v1",
        "credential_environment": "LUMYN_PROVIDER_API_KEY",
        "budget_posture": "capped",
        "redaction_posture": "redacted",
    }
    try:
        validator.factoryd_config_capability_grants = lambda: [active_wildcard_grant]
        try:
            validate_model_provider_gate(active_wildcard_task)
        except AssertionError as exc:
            if "active model_provider_endpoint grants must be task-scoped" not in str(exc):
                raise
        else:
            fail("self-test expected active wildcard model-provider grant to fail")
    finally:
        validator.factoryd_config_capability_grants = original_config_grants

    pending_base_url_task = model_provider_gate_task("T11.1", "T11.1")
    pending_base_url_task["factoryd_runtime"]["capability_grants"] = []
    pending_base_url_grant = {
        "task_id": "T11.1",
        "capability": "model_provider_endpoint",
        "evidence_ref": ".factory/artifacts/approvals/model_provider_endpoint.md",
    }
    pending_base_url_grant.update(
        {
            "approved": True,
            "network_allowlist": ["api.example.com"],
            "provider_identity": "example-provider",
            "provider_model": "example-model",
            "provider_endpoint": "   ",
            "base_url": "pending-approved-base-url",
            "credential_environment": "LUMYN_PROVIDER_API_KEY",
            "budget_posture": "capped",
            "redaction_posture": "redacted",
        }
    )
    try:
        validator.factoryd_config_capability_grants = lambda: [pending_base_url_grant]
        try:
            validate_model_provider_gate(pending_base_url_task)
        except AssertionError as exc:
            if "pending placeholders" not in str(exc):
                raise
        else:
            fail("self-test expected whitespace endpoint with pending base_url to fail")
    finally:
        validator.factoryd_config_capability_grants = original_config_grants

    approved_seed_task = model_provider_gate_task("T11.1", "T11.1")
    approved_seed_grant = approved_seed_task["factoryd_runtime"]["capability_grants"][0]
    approved_seed_grant.update(
        {
            "approved": True,
            "network_allowlist": ["api.example.com"],
            "provider_identity": "example-provider",
            "provider_model": "example-model",
            "provider_endpoint": "https://api.example.com/v1",
            "credential_environment": "LUMYN_PROVIDER_API_KEY",
            "budget_posture": "capped",
            "redaction_posture": "redacted",
        }
    )
    try:
        validate_model_provider_gate(approved_seed_task)
    except AssertionError as exc:
        if "seed model_provider_endpoint grants must stay approved false" not in str(exc):
            raise
    else:
        fail("self-test expected approved seed model-provider grant to fail")

    non_string_allowlist_task = model_provider_gate_task("T11.1", "T11.1")
    non_string_allowlist_grant = non_string_allowlist_task["factoryd_runtime"]["capability_grants"][0]
    non_string_allowlist_grant["network_allowlist"] = [{"host": "api.example.com"}]
    try:
        validate_model_provider_gate(non_string_allowlist_task)
    except AssertionError as exc:
        if "network_allowlist must be a non-empty string list" not in str(exc):
            raise
    else:
        fail("self-test expected non-string provider allowlist to fail")

    non_string_metadata_task = model_provider_gate_task("T11.1", "T11.1")
    non_string_metadata_task["factoryd_runtime"]["capability_grants"] = []
    non_string_metadata_grant = dict(active_config_grant)
    non_string_metadata_grant["provider_model"] = {"name": "example-model"}
    try:
        validator.factoryd_config_capability_grants = lambda: [non_string_metadata_grant]
        try:
            validate_model_provider_gate(non_string_metadata_task)
        except AssertionError as exc:
            if "fields must be non-empty strings" not in str(exc):
                raise
        else:
            fail("self-test expected non-string provider metadata to fail")
    finally:
        validator.factoryd_config_capability_grants = original_config_grants

    missing_requirement_field_task = model_provider_gate_task("T11.1", "T11.1")
    missing_requirement_field_task["model_provider_requirements"]["required_fields"].remove("network_allowlist")
    try:
        validate_model_provider_gate(missing_requirement_field_task)
    except AssertionError as exc:
        if "required_fields missing" not in str(exc):
            raise
    else:
        fail("self-test expected missing model-provider requirement field to fail")

    missing_seed_metadata_task = model_provider_gate_task("T11.1", "T11.1")
    missing_seed_metadata_task["factoryd_runtime"]["capability_grants"] = [
        {
            "task_id": "T11.1",
            "capability": "model_provider_endpoint",
            "approved": False,
        }
    ]
    try:
        validator.factoryd_config_capability_grants = lambda: [active_config_grant]
        try:
            validate_model_provider_gate(missing_seed_metadata_task)
        except AssertionError as exc:
            if "seed model_provider_endpoint grant missing fields" not in str(exc):
                raise
        else:
            fail("self-test expected missing seed provider metadata to fail")
    finally:
        validator.factoryd_config_capability_grants = original_config_grants

    unknown_gate_task = propagated_task("T3", ["T2.6"])
    unknown_gate_task["gated_by_acceptance_items"] = [
        {
            "acceptance_item_id": "PULL-999",
            "required_status": "implemented",
            "reason": "intentional invalid gate for self-test",
        }
    ]
    try:
        validate_task_execution_compiler_fields(unknown_gate_task)
    except AssertionError as exc:
        if "unknown acceptance item" not in str(exc):
            raise
    else:
        fail("self-test expected unknown acceptance item gate to fail")

    missing_standalone_live_eval_gate = propagated_task("T11", ["T2.6"])
    try:
        validate_standalone_task_packet(missing_standalone_live_eval_gate, "T2.6")
    except AssertionError as exc:
        if "must gate live eval dispatch" not in str(exc):
            raise
    else:
        fail("self-test expected standalone live-eval packet without pull gates to fail")

    missing_repair_live_eval_gate = propagated_task("T11-repair-001", ["T2.6"])
    try:
        validate_standalone_task_packet(missing_repair_live_eval_gate, "T2.6")
    except AssertionError as exc:
        if "must gate live eval dispatch" not in str(exc):
            raise
    else:
        fail("self-test expected live-eval repair packet without pull gates to fail")

    nested_slice_packets = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    del nested_slice_packets["tasks"][1]["slice_type"]
    validate_task_packets(nested_slice_packets, "T2.6")

    conflicting_slice_packets = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    conflicting_slice_packets["tasks"][1]["slice_type"] = "foundation"
    conflicting_slice_packets["tasks"][1]["slice_rationale"]["slice_type"] = "vertical"
    try:
        validate_task_packets(conflicting_slice_packets, "T2.6")
    except AssertionError as exc:
        if "conflicting slice_type declarations" not in str(exc):
            raise
    else:
        fail("self-test expected conflicting slice type declarations to fail")

    missing_delivery_slice_packets = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    missing_delivery_slice_packets["tasks"][1].pop("delivery_slice_refs", None)
    try:
        validate_task_packets(missing_delivery_slice_packets, "T2.6")
    except AssertionError as exc:
        if "delivery_slice_refs" not in str(exc):
            raise
    else:
        fail("self-test expected missing delivery_slice_refs to fail")

    missing_dotted_delivery_slice_packets = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T11.1", ["T2.6", "T11"])]
    }
    missing_dotted_delivery_slice_packets["tasks"][1].pop("delivery_slice_refs", None)
    try:
        validate_task_packets(missing_dotted_delivery_slice_packets, "T2.6")
    except AssertionError as exc:
        if "delivery_slice_refs" not in str(exc):
            raise
    else:
        fail("self-test expected dotted task missing delivery_slice_refs to fail")

    deprecated_worker_packets = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    deprecated_worker_packets["tasks"][1]["required_worker_chain"] = ["task-executor", "ship-pr"]
    try:
        validate_task_packets(deprecated_worker_packets, "T2.6")
    except AssertionError as exc:
        if "deprecated worker" not in str(exc):
            raise
    else:
        fail("self-test expected deprecated active worker to fail")

    control_allowed_packets = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    control_allowed_packets["tasks"][1]["allowed_paths"].append(
        "./.factory/artifacts/prd-to-plan/lumyn-mvp/scope-closure-map.json"
    )
    control_allowed_packets["tasks"][1]["allowed_paths"].append(
        ".factory//artifacts/prd-to-plan/lumyn-mvp/scope-closure-map.json"
    )
    control_allowed_packets["tasks"][1]["allowed_paths"].append(
        ".factory/artifacts/prd-to-plan/lumyn-mvp/scope-closure-map.json"
    )
    control_allowed_packets["tasks"][1]["allowed_paths"].append(
        ".factory/artifacts/prd-to-plan/lumyn-mvp/risk-classification.json"
    )
    try:
        validate_task_packets(control_allowed_packets, "T2.6")
    except AssertionError as exc:
        if "runtime-owned control artifact" not in str(exc):
            raise
    else:
        fail("self-test expected runtime-owned control artifact allowed path to fail")

    traversal_allowed_packets = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    traversal_allowed_packets["tasks"][1]["allowed_paths"].append(
        ".factory/artifacts/task-runs/T3/../../prd-to-plan/lumyn-mvp/scope-closure-map.json"
    )
    try:
        validate_task_packets(traversal_allowed_packets, "T2.6")
    except AssertionError as exc:
        if "must not traverse outside the repository" not in str(exc):
            raise
    else:
        fail("self-test expected traversing allowed path to fail")

    missing_runtime_pin_packets = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    del missing_runtime_pin_packets["tasks"][1]["runtime_pins"]["provider_policy"]
    try:
        validate_task_packets(missing_runtime_pin_packets, "T2.6")
    except AssertionError as exc:
        if "runtime_pins" not in str(exc) and "runtime pin fields" not in str(exc):
            raise
    else:
        fail("self-test expected missing runtime pin to fail")

    boolean_iteration_budget_packets = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    boolean_iteration_budget_packets["tasks"][1]["max_iterations"] = True
    try:
        validate_task_packets(boolean_iteration_budget_packets, "T2.6")
    except AssertionError as exc:
        if "max_iterations" not in str(exc):
            raise
    else:
        fail("self-test expected boolean max_iterations to fail")

    missing_alignment_ref_packets = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    del missing_alignment_ref_packets["tasks"][1]["alignment_gate_ref"]
    try:
        validate_task_packets(missing_alignment_ref_packets, "T2.6")
    except AssertionError as exc:
        if "alignment_gate_ref" not in str(exc):
            raise
    else:
        fail("self-test expected missing alignment gate ref to fail")

    missing_acceptance_item_packets = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    del missing_acceptance_item_packets["tasks"][1]["acceptance_item_ids"]
    try:
        validate_task_packets(missing_acceptance_item_packets, "T2.6")
    except AssertionError as exc:
        if "acceptance_item_ids" not in str(exc):
            raise
    else:
        fail("self-test expected missing acceptance item ids to fail")

    incomplete_worker_chain_packets = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    incomplete_worker_chain_packets["tasks"][1]["required_worker_chain"] = ["task-executor"]
    try:
        validate_task_packets(incomplete_worker_chain_packets, "T2.6")
    except AssertionError as exc:
        if "required_worker_chain" not in str(exc):
            raise
    else:
        fail("self-test expected incomplete worker chain to fail")

    disabled_lifecycle_gate_packets = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    disabled_lifecycle_gate_packets["tasks"][1]["lifecycle_gates"]["ci_required"] = False
    try:
        validate_task_packets(disabled_lifecycle_gate_packets, "T2.6")
    except AssertionError as exc:
        if "lifecycle_gates" not in str(exc):
            raise
    else:
        fail("self-test expected disabled lifecycle gate to fail")

    disconnected_packets = {"tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3-repair-001", [])]}
    try:
        validate_task_packets(disconnected_packets, "T2.6")
    except AssertionError as exc:
        if "does not depend on it" not in str(exc):
            raise
    else:
        fail("self-test expected disconnected T3 task to fail")

    slug_baseline_packets = {
        "tasks": [
            propagated_task("task-context", []),
            propagated_task("task-dev-architecture-propagation", ["task-context"]),
            propagated_task("feature-local-check", []),
        ]
    }
    try:
        validate_task_packets(slug_baseline_packets, "task-dev-architecture-propagation")
    except AssertionError as exc:
        if "does not depend on it" not in str(exc):
            raise
    else:
        fail("self-test expected disconnected task after slug baseline to fail")

    placeholder_packets = {"tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]}
    placeholder_packets["tasks"][1]["ci_lane_refs"] = [{}]
    try:
        validate_task_packets(placeholder_packets, "T2.6")
    except AssertionError as exc:
        if "ci_lane_refs" not in str(exc):
            raise
    else:
        fail("self-test expected placeholder refs to fail")

    disabled_scanner_packets = {"tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]}
    disabled_scanner_packets["tasks"][1]["security_scanner_gates"] = {"required": False, "scanner": "CodeQL"}
    try:
        validate_task_packets(disabled_scanner_packets, "T2.6")
    except AssertionError as exc:
        if "security_scanner_gates" not in str(exc):
            raise
    else:
        fail("self-test expected disabled scanner without exception to fail")

    missing_policy_packets = {"tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]}
    missing_policy_packets["tasks"][1]["engineering_policy_refs"] = missing_policy_packets["tasks"][1][
        "engineering_policy_refs"
    ][:1]
    try:
        validate_task_packets(missing_policy_packets, "T2.6")
    except AssertionError as exc:
        if "engineering_policy_refs missing" not in str(exc):
            raise
    else:
        fail("self-test expected missing engineering policy refs to fail")

    missing_ci_packets = {"tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]}
    missing_ci_packets["tasks"][1]["ci_lane_refs"] = [
        item for item in missing_ci_packets["tasks"][1]["ci_lane_refs"] if item["lane"] == "core"
    ]
    try:
        validate_task_packets(missing_ci_packets, "T2.6")
    except AssertionError as exc:
        if "ci_lane_refs missing" not in str(exc):
            raise
    else:
        fail("self-test expected missing CI lane refs to fail")

    blank_ci_packets = {"tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]}
    blank_ci_packets["tasks"][1]["ci_lane_refs"][1]["command_refs"] = [""]
    try:
        validate_task_packets(blank_ci_packets, "T2.6")
    except AssertionError as exc:
        if "missing guide propagation fields" not in str(exc):
            raise
    else:
        fail("self-test expected blank CI lane evidence to fail")

    missing_planning_packets = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    del missing_planning_packets["tasks"][1]["changelog"]
    try:
        validate_task_packets(missing_planning_packets, "T2.6")
    except AssertionError as exc:
        if "missing planning-skill fields" not in str(exc):
            raise
    else:
        fail("self-test expected missing planning-skill fields to fail")

    missing_runner_ready_field_packets = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    del missing_runner_ready_field_packets["tasks"][1]["required_proof_level"]
    try:
        validate_task_packets(missing_runner_ready_field_packets, "T2.6")
    except AssertionError as exc:
        if "missing guide propagation fields" not in str(exc):
            raise
    else:
        fail("self-test expected missing runner-ready task field to fail")

    missing_architecture_targets_packets = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    del missing_architecture_targets_packets["tasks"][1]["architecture_target_paths"]
    try:
        validate_task_packets(missing_architecture_targets_packets, "T2.6")
    except AssertionError as exc:
        if "architecture_target_paths must be a non-empty list" not in str(exc):
            raise
    else:
        fail("self-test expected missing architecture target paths to fail")

    broad_internal_allowed_packets = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    broad_internal_allowed_packets["tasks"][1]["allowed_paths"].append("internal/")
    try:
        validate_task_packets(broad_internal_allowed_packets, "T2.6")
    except AssertionError as exc:
        if "allowed_paths must not include broad internal/" not in str(exc):
            raise
    else:
        fail("self-test expected broad internal allowed path to fail")

    missing_allowed_target_packets = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    missing_allowed_target_packets["tasks"][1]["allowed_paths"] = [
        path
        for path in missing_allowed_target_packets["tasks"][1]["allowed_paths"]
        if path != "internal/source/"
    ]
    try:
        validate_task_packets(missing_allowed_target_packets, "T2.6")
    except AssertionError as exc:
        if "allowed_paths must include architecture_target_paths" not in str(exc):
            raise
    else:
        fail("self-test expected missing allowed architecture target to fail")

    repo_root_target_packets = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    repo_root_target_packets["tasks"][1]["architecture_target_paths"] = ["."]
    repo_root_target_packets["tasks"][1]["allowed_paths"] = ["."]
    try:
        validate_task_packets(repo_root_target_packets, "T2.6")
    except AssertionError as exc:
        if "must not resolve to the repository root" not in str(exc):
            raise
    else:
        fail("self-test expected repo-root architecture target to fail")

    non_string_target_packets = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    non_string_target_packets["tasks"][1]["architecture_target_paths"] = [123]
    try:
        validate_task_packets(non_string_target_packets, "T2.6")
    except AssertionError as exc:
        if "must be a string" not in str(exc):
            raise
    else:
        fail("self-test expected non-string architecture path to fail")

    behavioral_without_scorecard_packets = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    behavioral_without_scorecard_packets["tasks"][1]["required_proof_level"] = "workflow_behavior"
    behavioral_without_scorecard_packets["tasks"][1]["proof_scorecard_required"] = False
    try:
        validate_task_packets(behavioral_without_scorecard_packets, "T2.6")
    except AssertionError as exc:
        if "missing guide propagation fields" not in str(exc):
            raise
    else:
        fail("self-test expected behavioral proof without scorecard to fail")

    customer_safe_without_recursive_redaction_packets = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    customer_safe_without_recursive_redaction_packets["tasks"][1]["redaction_posture"] = {
        "classification": "customer_safe",
        "customer_safe": True,
        "redaction_notes": "Shareable evidence requires redaction.",
    }
    try:
        validate_task_packets(customer_safe_without_recursive_redaction_packets, "T2.6")
    except AssertionError as exc:
        if "missing guide propagation fields" not in str(exc):
            raise
    else:
        fail("self-test expected customer-safe redaction without recursive policy to fail")

    api_contract_without_adr = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    api_contract_without_adr["tasks"][1]["contract_impact"] = "Changes API response behavior."
    api_contract_without_adr["tasks"][1]["adr_required"] = False
    try:
        validate_task_packets(api_contract_without_adr, "T2.6")
    except AssertionError as exc:
        if "requires adr_required=true" not in str(exc):
            raise
    else:
        fail("self-test expected API contract impact without ADR to fail")

    non_contract_specific_text = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    non_contract_specific_text["tasks"][1]["contract_impact"] = "No specific user-visible behavior impact."
    non_contract_specific_text["tasks"][1]["adr_required"] = False
    validate_task_packets(non_contract_specific_text, "T2.6")

    no_public_contract_impact = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    no_public_contract_impact["tasks"][1]["contract_impact"] = "No public API or contract impact."
    no_public_contract_impact["tasks"][1]["adr_required"] = False
    validate_task_packets(no_public_contract_impact, "T2.6")

    linux_absolute_path_packets = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    linux_absolute_path_packets["tasks"][1]["artifact_refs"] = [
        "/workspace/lumyn/.factory/artifacts/run.json",
        "/root/lumyn/.factory/artifacts/run.json",
    ]
    try:
        validate_task_packets(linux_absolute_path_packets, "T2.6")
    except AssertionError as exc:
        if "machine-local absolute path" not in str(exc):
            raise
    else:
        fail("self-test expected Linux absolute path to fail")

    absolute_path_key_packets = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    absolute_path_key_packets["tasks"][1]["source_hashes"] = {
        "/workspace/lumyn/internal/result/result.go": "sha256:abc123"
    }
    try:
        validate_task_packets(absolute_path_key_packets, "T2.6")
    except AssertionError as exc:
        if "machine-local absolute path" not in str(exc):
            raise
    else:
        fail("self-test expected absolute map key path to fail")

    foundation_without_justification = {
        "tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]
    }
    foundation_without_justification["tasks"][1]["slice_type"] = "foundation"
    foundation_without_justification["tasks"][1]["slice_rationale"]["slice_type"] = "foundation"
    try:
        validate_task_packets(foundation_without_justification, "T2.6")
    except AssertionError as exc:
        if "non_vertical_justification" not in str(exc):
            raise
    else:
        fail("self-test expected non-vertical task without justification to fail")

    validate_standalone_task_packet(propagated_task("T3-repair-001", ["T2.6"]), "T2.6")

    print("repo-pack validator self-test passed")
    return 0
