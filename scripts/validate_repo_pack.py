#!/usr/bin/env python3
"""Validate the active Lumyn migration product and Factory control generation."""

from __future__ import annotations

import json
import re
import sys
from pathlib import Path
from typing import Any

from repo_pack_architecture import validate_architecture_budget_policy
from repo_pack_validation.acceptance_text import validate_acceptance_text
from repo_pack_validation.authority import validate_authority_grants
from repo_pack_validation.runtime_pins import validate_runtime_pins
from repo_pack_validation.self_tests import run_repo_pack_self_tests
from repo_pack_validation.task_contracts import validate_migration_task_contracts

ROOT = Path(__file__).resolve().parents[1]
PLAN_REL = ".factory/artifacts/prd-to-plan/lumyn-migration-mvp"
PLAN_DIR = ROOT / PLAN_REL
HISTORICAL_PLAN_REL = ".factory/artifacts/prd-to-plan/lumyn-mvp"
PRD = ROOT / "docs/product/prd.md"
PLAN = ROOT / "docs/product/plan.md"
CONFIG = ROOT / ".factory/factoryd.example.json"
AUTOSHIP_CONFIG = ROOT / ".factory/factoryd.autoship.example.json"
ACTIVE_CONFIG = ROOT / ".factory/factoryd.json"

ARTIFACT_PATHS = {
    "context": PLAN_DIR / "context-brief.json",
    "risk": PLAN_DIR / "risk-classification.json",
    "plan": PLAN_DIR / "execution-plan.json",
    "packets": PLAN_DIR / "task-packets.json",
    "contract": PLAN_DIR / "validation-contract.json",
    "ledger": PLAN_DIR / "acceptance-ledger.json",
    "mapping": PLAN_DIR / "acceptance-mapping.json",
    "closure": PLAN_DIR / "scope-closure-map.json",
}

ARTIFACT_REFS = {
    name: path.relative_to(ROOT).as_posix() for name, path in ARTIFACT_PATHS.items()
}

EXPECTED_TASK_IDS = ("M0", "M1", "M2", "M2.5", "M3", "M4", "M5", "M6", "M7", "M8", "M9", "M10")
EXPECTED_GROUP_IDS = {
    "retained_foundation",
    "product_rebaseline",
    "provider_change",
    "migration_benchmark",
    "authorization",
    "impact",
    "migration_plan",
    "bounded_patch",
    "verification",
    "evidence",
    "draft_pr_delivery",
    "two_sided_activation",
    "design_partner_qualification",
    "design_partner_pilot",
}
EXPECTED_SLICE_IDS = {
    "migration_foundation",
    "migration_engine",
    "design_partner_pilot",
}
ALLOWED_HISTORICAL_TASK_IDS = {"T1", "T2", "T2.5", "T2.6", "T2.7", "T3"}
ALLOWED_ITEM_STATUSES = {
    "planned",
    "implemented",
    "partial",
    "blocked",
    "deferred",
    "deferred_with_approval",
    "out_of_scope",
}
TERMINAL_EVIDENCE_STATUSES = {"implemented", "deferred_with_approval"}
EXPECTED_CARRY_FORWARD = {
    "BASE-001",
    "BASE-002",
    "BASE-004",
    "BASE-005",
    "REB-003",
}
EXPECTED_CAPABILITIES = {
    "M2.5": {"approval"},
    "M8": {"approval", "credentials", "network"},
    "M9": {"approval", "credentials", "network"},
    "M10": {"approval", "credentials", "network"},
}
EXPECTED_PRODUCT_AUTHORITIES = {
    "M8": {"command_execution", "sandbox_network", "sandbox_credential", "sandbox_request_disclosure", "artifact_retention", "artifact_deletion"},
    "M9": {"github_branch_write", "github_pr_write", "artifact_retention", "artifact_deletion"},
    "M10": {"customer_repo_read", "customer_repo_write", "command_execution", "package_registry_read", "github_branch_write", "github_pr_write", "campaign_receipt", "artifact_retention", "artifact_deletion"},
}
EXPECTED_OPTIONAL_PRODUCT_AUTHORITIES = {
    "M8": {"provider_trust_status_read"},
    "M9": {"provider_trust_status_read", "provider_attestation"},
    "M10": {"provider_trust_status_read", "sandbox_network", "sandbox_credential", "sandbox_request_disclosure", "provider_attestation"},
}
PRODUCT_AUTHORITY_CAPABILITIES = (
    set().union(*EXPECTED_PRODUCT_AUTHORITIES.values())
    | set().union(*EXPECTED_OPTIONAL_PRODUCT_AUTHORITIES.values())
)
REQUIRED_TASK_FIELDS = {
    "task_id",
    "objective",
    "risk_class",
    "required_review",
    "requires_human_approval",
    "requires_credentials",
    "requires_network",
    "blocked_by",
    "allowed_paths",
    "forbidden_paths",
    "scope_exclusions",
    "input_artifacts",
    "commands",
    "baseline_commands",
    "red_first_commands",
    "validation_commands",
    "final_validation_commands",
    "acceptance_checks",
    "acceptance_ledger_ref",
    "acceptance_item_ids",
    "validation_contract_inheritance",
    "required_worker_chain",
    "lifecycle_gates",
    "worker_evidence_required",
    "lifecycle_evidence_required",
    "test_matrix_refs",
    "ci_lane_refs",
    "security_scanner_gates",
    "engineering_policy_refs",
    "architecture_guidance_refs",
    "architecture_target_paths",
    "runtime_pins",
    "factory_compatibility",
    "alignment_gate_ref",
    "plan_drift_policy_ref",
    "worker_type",
    "factoryd_runtime",
    "requires_capabilities",
    "product_authority_requirements",
    "stop_conditions",
    "coverage_policy_refs",
    "delivery_slice_refs",
    "docs_sync_refs",
    "changelog_intent",
    "versioning_impact",
    "migration_impact",
    "max_iterations",
    "required_proof_level",
    "proof_scorecard_required",
    "failure_hypothesis",
    "semantic_invariants",
}
REQUIRED_DOCS = {
    "README.md": [
        "provider-sponsored",
        "customer-controlled", "provider-authenticated",
        "Current Implementation Status",
        HISTORICAL_PLAN_REL,
    ],
    "AGENTS.md": [
        "Two Principals, Two Authorities",
        "customer_repo_read",
        "customer_repo_write",
        "github_pr_write", "provider-signed acknowledgement",
        "Passive Codex review settle is required before merge",
    ],
    "WORKFLOW.md": [
        "provider packet",
        "explicit consumer-signed repository authorization issuance",
        "read-only impact",
        "Green CI alone is not merge-ready",
        "process escape",
    ],
    "docs/product/prd.md": [
        "Provider-Sponsored Verified API Migrations",
        "API Provider Job",
        "API Consumer Job",
        "Read Before Write",
        "Provider Change Packet", "consumer-receipt-key-binding", "fail-closed host-isolation",
        "Draft PR Delivery",
        "Falsification And Reframe Gates",
    ],
    "docs/product/plan.md": [
        "M0: Correct the command and result foundations",
        "M4: Analyze TypeScript consumer impact",
        "M9: Produce migration evidence and open an idempotent draft PR",
        "M10: Run one qualified design-partner campaign", "campaign receipt acknowledge/ack import",
    ],
    "docs/dev/dev_guides.md": [
        "Migration Corpus Policy",
        "TypeScript Impact Policy",
        "Patch And Filesystem Policy",
        "Proof-Of-Behavior Policy", "Homebrew wait for the separate approved",
        "Do not merge manually through `gh pr merge`",
    ],
    "docs/architecture/architecture_guides.md": [
        "Provider Campaign Plane",
        "Consumer Execution Plane",
        "Patch Safety Boundary",
        "Live Sandbox Boundary",
        "GitHub Boundary", "provider-authenticated consumer receipt-key bindings",
    ],
    "docs/architecture/adr-0002-provider-sponsored-customer-controlled-migrations.md": [
        "Provider-Sponsored, Customer-Controlled",
        "Provider sponsorship conveys no customer-data authority", "provider-signed, deduplicated acknowledgement",
        "Rejected Alternatives",
    ],
    "docs/factory/README.md": [
        PLAN_REL,
        "Product Authority Is Not Factory Authority",
        "customer_repo_read", "consumer signer binding and receives a provider-signed",
        "immutable historical records",
    ],
    ".factory/README.md": [
        PLAN_REL,
        "immutable records",
    ],
}

MACHINE_LOCAL_RE = re.compile(r"(?:^|[\s\"'])(?:/Users/|/home/|file://|[A-Za-z]:\\\\)")
ACCEPTANCE_ID_RE = re.compile(r"`([A-Z]+-\d{3})`:")


def fail(message: str) -> None:
    raise AssertionError(message)


def require(condition: bool, message: str) -> None:
    if not condition:
        fail(message)


def load_json(path: Path) -> dict[str, Any]:
    require(path.exists(), f"missing JSON artifact: {path.relative_to(ROOT)}")
    try:
        payload = json.loads(path.read_text())
    except Exception as exc:
        fail(f"{path.relative_to(ROOT)} is not valid JSON: {exc}")
    require(isinstance(payload, dict), f"{path.relative_to(ROOT)} must contain a JSON object")
    return payload


def nonempty_string(value: Any) -> bool:
    return isinstance(value, str) and bool(value.strip())


def nonempty_list(value: Any) -> bool:
    return isinstance(value, list) and bool(value)


def list_of_strings(value: Any) -> bool:
    return nonempty_list(value) and all(nonempty_string(item) for item in value)


def source_path(value: str) -> str:
    return value.split("#", 1)[0]


def require_repo_ref(value: Any, label: str, *, allow_pending: bool = False) -> None:
    require(nonempty_string(value), f"{label} must be a non-empty repo-relative reference")
    ref = str(value).strip()
    if ref.startswith(("http://", "https://")):
        return
    if allow_pending and ref.startswith("pending:"):
        return
    require(not Path(source_path(ref)).is_absolute(), f"{label} must not be absolute")
    require(".." not in Path(source_path(ref)).parts, f"{label} must stay inside the repository")
    require((ROOT / source_path(ref)).exists(), f"{label} points to missing path {source_path(ref)}")


def contains_machine_local_path(value: Any) -> bool:
    if isinstance(value, str):
        return bool(MACHINE_LOCAL_RE.search(value))
    if isinstance(value, list):
        return any(contains_machine_local_path(item) for item in value)
    if isinstance(value, dict):
        return any(contains_machine_local_path(item) for item in value.values())
    return False


def contains_true_key(value: Any, keys: set[str]) -> bool:
    if isinstance(value, dict):
        for key, item in value.items():
            if key in keys and item is True:
                return True
            if contains_true_key(item, keys):
                return True
    if isinstance(value, list):
        return any(contains_true_key(item, keys) for item in value)
    return False


def expected_acceptance_ids() -> set[str]:
    text = PRD.read_text()
    require("## Acceptance Tests" in text and "## Success Metrics" in text, "PRD acceptance section is missing")
    section = text.split("## Acceptance Tests", 1)[1].split("## Success Metrics", 1)[0]
    ids = set(ACCEPTANCE_ID_RE.findall(section))
    require(len(ids) == 62, f"PRD must define exactly 62 unique acceptance IDs; found {len(ids)}")
    return ids


def validate_docs() -> None:
    for relative, tokens in REQUIRED_DOCS.items():
        path = ROOT / relative
        require(path.exists(), f"missing required document: {relative}")
        text = path.read_text()
        for token in tokens:
            require(token in text, f"{relative} missing required product/operating token: {token}")
        require(not MACHINE_LOCAL_RE.search(text), f"{relative} contains a machine-local path")
    require((ROOT / HISTORICAL_PLAN_REL / "README.md").exists(), "historical plan must remain present")
    require((PLAN_DIR / "README.md").exists(), "active plan README is missing")
    prd = PRD.read_text()
    require("Public fixtures" in prd and "not" in prd.split("Public fixtures", 1)[1][:180], "PRD must distinguish public fixtures from demand evidence")
    require("no auto-merge" in prd.lower() or "never auto-merges" in prd.lower(), "PRD must forbid auto-merge")
    require("production credentials" in prd.lower(), "PRD must state the production-credential boundary")


def validate_ci_controls() -> None:
    required_checks = load_json(ROOT / ".github/required-checks.json")
    serialized = json.dumps(required_checks)
    for check in ["validate", "CodeQL analyze"]:
        require(check in serialized, f"required-check metadata missing {check}")
    codeowners = (ROOT / ".github/CODEOWNERS").read_text()
    for token in ["/.github/** @davidahmann", "/.factory/** @davidahmann", "/docs/product/** @davidahmann", "/schemas/** @davidahmann"]:
        require(token in codeowners, f"CODEOWNERS missing {token}")
    validate = (ROOT / ".github/workflows/validate.yml").read_text()
    codeql = (ROOT / ".github/workflows/codeql.yml").read_text()
    for label, text in [("validate.yml", validate), ("codeql.yml", codeql)]:
        for token in ["permissions:", "concurrency:", "timeout-minutes:", "actions/checkout@v6.0.2", "actions/setup-go@v6.3.0", "go-version-file: go.mod", "check-latest: false"]:
            require(token in text, f"{label} missing {token}")
    for token in ["github/codeql-action/init@v4", "github/codeql-action/autobuild@v4", "github/codeql-action/analyze@v4"]:
        require(token in codeql, f"codeql.yml missing {token}")


def validate_context(context: dict[str, Any]) -> None:
    require(context.get("artifact_type") == "context_brief", "context artifact_type must be context_brief")
    require(context.get("source_prd_ref") == "docs/product/prd.md", "context must cite the canonical PRD")
    require(context.get("alignment_decisions", {}).get("status") == "resolved", "context alignment decisions must be resolved")
    require(context.get("alignment_decisions", {}).get("implementation_may_start") is True, "context must explicitly allow bounded implementation")
    findings = "\n".join(str(item) for item in context.get("proven_findings", []))
    for token in ["not implemented", "generic pass", "OpenAPI/docs intake", "provider-status", "host isolation"]:
        require(token.lower() in findings.lower(), f"context baseline missing {token}")
    summary = str(context.get("system_model_summary", "")).lower()
    for token in ["provider campaign plane", "consumer execution plane", "explicit"]:
        require(token in summary, f"context system model missing {token}")
    decisions = json.dumps(context.get("alignment_decisions", {})).lower()
    for token in ["api provider", "api consumer", "draft github pr", "no model-provider adapter", "consumer authorization", "trust freshness", "sponsored connection meter", "distribution posture", "immutable and non-active"]:
        require(token in decisions, f"context alignment decisions missing {token}")
    require(context.get("factory_compatibility", {}).get("profile_ref") == "profiles/lumyn.yaml", "context must cite the Lumyn Factory profile")
    validate_runtime_pins(context.get("runtime_pins"), "context")


def validate_risk(risk: dict[str, Any]) -> None:
    require(risk.get("artifact_type") == "risk_classification", "risk artifact_type must be risk_classification")
    require(risk.get("default_risk_class") == "high", "migration plan default risk must be high")
    rules = json.dumps(risk.get("risk_rules", [])).lower()
    for token in ["provider change packet", "consumer authorization", "customer repository read or write", "patch generation", "host-isolated command", "sandbox credentials", "github draft pr write", "minimal campaign receipt", "provider-visible attestation", "irreversible disclosure", "distribution or oss claim", "product-signal closure"]:
        require(token in rules, f"risk classification missing high-risk surface: {token}")


def validate_ledger(ledger: dict[str, Any], required_ids: set[str]) -> dict[str, dict[str, Any]]:
    require(ledger.get("artifact_type") == "acceptance_ledger", "ledger artifact_type must be acceptance_ledger")
    require(ledger.get("source_ref") == "docs/product/prd.md", "ledger source_ref must cite the PRD")
    require(ledger.get("coverage_policy", {}).get("closure_unit") == "acceptance_item", "ledger closure unit must be acceptance_item")
    require(ledger.get("coverage_policy", {}).get("group_only_refs_allowed") is False, "group-only closure must be forbidden")
    items = ledger.get("items")
    require(isinstance(items, list), "ledger.items must be a list")
    by_id: dict[str, dict[str, Any]] = {}
    for index, item in enumerate(items):
        require(isinstance(item, dict), f"ledger.items[{index}] must be an object")
        item_id = item.get("acceptance_item_id")
        require(nonempty_string(item_id), f"ledger.items[{index}] missing acceptance_item_id")
        require(item_id not in by_id, f"duplicate ledger item {item_id}")
        by_id[str(item_id)] = item
    require(set(by_id) == required_ids, f"ledger IDs differ from PRD: missing={sorted(required_ids-set(by_id))}, extra={sorted(set(by_id)-required_ids)}")
    require(ledger.get("acceptance_item_count") == len(required_ids), "ledger acceptance_item_count is stale")
    for item_id, item in by_id.items():
        for field in ["group_id", "source_ref", "source_text", "kind", "evidence_mode", "task_refs", "validation_refs", "status", "risk_class", "notes"]:
            require(field in item, f"{item_id} missing {field}")
        require(item["status"] in ALLOWED_ITEM_STATUSES, f"{item_id} has invalid status {item['status']}")
        require_repo_ref(item["source_ref"], f"{item_id}.source_ref")
        task_refs = item.get("task_refs")
        require(list_of_strings(task_refs), f"{item_id}.task_refs must be non-empty")
        unknown_tasks = set(task_refs) - set(EXPECTED_TASK_IDS) - ALLOWED_HISTORICAL_TASK_IDS
        require(not unknown_tasks, f"{item_id} references unknown tasks {sorted(unknown_tasks)}")
        require(any(ref in EXPECTED_TASK_IDS for ref in task_refs), f"{item_id} must map to an active migration task")
        if item["status"] in TERMINAL_EVIDENCE_STATUSES:
            require(list_of_strings(item.get("evidence_refs")), f"{item_id} terminal status requires evidence_refs")
            require(nonempty_string(item.get("recorded_by")), f"{item_id} terminal status requires recorded_by")
            require(nonempty_string(item.get("recorded_at")), f"{item_id} terminal status requires recorded_at")
            for evidence_index, ref in enumerate(item["evidence_refs"]):
                require_repo_ref(ref, f"{item_id}.evidence_refs[{evidence_index}]")
    actual_carry = {item_id for item_id, item in by_id.items() if item["status"] == "implemented"}
    require(actual_carry == EXPECTED_CARRY_FORWARD, f"retained foundation status drifted: expected={sorted(EXPECTED_CARRY_FORWARD)} actual={sorted(actual_carry)}")
    return by_id


def validate_mapping(mapping: dict[str, Any], required_ids: set[str]) -> None:
    require(mapping.get("artifact_type") == "acceptance_mapping", "mapping artifact_type must be acceptance_mapping")
    require(mapping.get("acceptance_ledger_ref") == ARTIFACT_REFS["ledger"], "mapping must cite active ledger")
    groups = mapping.get("groups")
    require(isinstance(groups, list), "mapping.groups must be a list")
    group_ids: set[str] = set()
    mapped: list[str] = []
    for group in groups:
        require(isinstance(group, dict), "mapping group must be an object")
        group_id = group.get("group_id")
        require(nonempty_string(group_id), "mapping group_id is required")
        require(group_id not in group_ids, f"duplicate mapping group {group_id}")
        group_ids.add(str(group_id))
        ids = group.get("acceptance_item_ids")
        require(list_of_strings(ids), f"mapping group {group_id} must contain acceptance_item_ids")
        mapped.extend(str(item) for item in ids)
    require(group_ids == EXPECTED_GROUP_IDS, f"mapping groups are incomplete: missing={sorted(EXPECTED_GROUP_IDS-group_ids)}, extra={sorted(group_ids-EXPECTED_GROUP_IDS)}")
    require(set(mapped) == required_ids, "mapping does not cover the exact PRD acceptance set")
    require(len(mapped) == len(required_ids), "each acceptance item must appear in exactly one mapping group")


def validate_slices(values: Any, required_ids: set[str], label: str) -> dict[str, dict[str, Any]]:
    require(isinstance(values, list), f"{label} must be a list")
    by_id: dict[str, dict[str, Any]] = {}
    covered: set[str] = set()
    for item in values:
        require(isinstance(item, dict), f"{label} entries must be objects")
        slice_id = item.get("slice_id")
        require(nonempty_string(slice_id), f"{label} slice_id is required")
        require(slice_id not in by_id, f"{label} duplicate slice {slice_id}")
        by_id[str(slice_id)] = item
        require(item.get("required_for_completion") is True, f"{label}.{slice_id} must be required")
        require(list_of_strings(item.get("task_refs")), f"{label}.{slice_id}.task_refs must be non-empty")
        require(list_of_strings(item.get("acceptance_item_ids")), f"{label}.{slice_id}.acceptance_item_ids must be non-empty")
        unknown_items = set(item["acceptance_item_ids"]) - required_ids
        require(not unknown_items, f"{label}.{slice_id} references unknown acceptance items {sorted(unknown_items)}")
        covered.update(item["acceptance_item_ids"])
    require(set(by_id) == EXPECTED_SLICE_IDS, f"{label} slice IDs differ from the active plan")
    require(covered == required_ids, f"{label} delivery slices do not cover all acceptance items")
    return by_id


def validate_plan(plan: dict[str, Any], required_ids: set[str]) -> None:
    require(plan.get("artifact_type") == "execution_plan", "execution plan artifact_type must be execution_plan")
    require(plan.get("source_prd_ref") == "docs/product/prd.md", "execution plan must cite the PRD")
    require(plan.get("authored_plan_ref") == "docs/product/plan.md", "execution plan must cite the authored plan")
    require(plan.get("context_brief_ref") == ARTIFACT_REFS["context"], "execution plan context ref is stale")
    require(plan.get("risk_classification_ref") == ARTIFACT_REFS["risk"], "execution plan risk ref is stale")
    require(plan.get("acceptance_ledger_ref") == ARTIFACT_REFS["ledger"], "execution plan ledger ref is stale")
    require(plan.get("acceptance_mapping_ref") == ARTIFACT_REFS["mapping"], "execution plan mapping ref is stale")
    require(plan.get("validation_contract_ref") == ARTIFACT_REFS["contract"], "execution plan contract ref is stale")
    require(plan.get("risk_class") == "high", "execution plan must be high risk")
    validate_runtime_pins(plan.get("runtime_pins"), "execution plan")
    require(plan.get("initial_slice", {}).get("task_id") == "M0", "M0 must be the initial slice")
    require(plan.get("alignment_gate", {}).get("status") == "resolved", "execution-plan alignment gate must be resolved")
    require(plan.get("alignment_gate", {}).get("blocked_live_tasks") == ["M2.5", "M8", "M9", "M10"], "execution plan must name live-gated tasks")
    retained = plan.get("retained_evidence_policy", {})
    require(retained.get("historical_plan_ref") == f"{HISTORICAL_PLAN_REL}/", "execution plan historical ref is stale")
    require(retained.get("status") == "immutable_non_active", "historical plan must be immutable and non-active")
    coverage = plan.get("acceptance_ledger_coverage", {})
    require(coverage.get("coverage_unit") == "acceptance_item", "execution-plan coverage unit must be acceptance_item")
    require(coverage.get("required_item_count") == len(required_ids), "execution-plan acceptance count is stale")
    require(coverage.get("group_only_refs_allowed") is False, "execution plan must forbid group-only closure")
    validate_slices(plan.get("delivery_slices"), required_ids, "execution_plan.delivery_slices")
    locked = json.dumps(plan.get("locked_decisions", {})).lower()
    for token in ["api provider", "api consumer", "draft github pr", "no model-provider adapter", "consumer authorization", "trust freshness", "sponsored connection meter", "distribution posture"]:
        require(token in locked, f"execution plan locked decisions missing {token}")
    rollback = " ".join(plan.get("rollback_or_deletion_test", [])).lower()
    require("future side effects and disclosures block" in rollback and "lumyn-controlled private" in rollback and "non-recallable" in rollback and "no public or provider-facing copy survives" not in rollback, "execution plan rollback must preserve irreversible external-disclosure semantics")
    require(not contains_true_key(plan, {"auto_merge", "default_branch_write", "provider_raw_repo_access"}), "execution plan enables a forbidden product capability")
    serialized = json.dumps(plan).lower()
    for stale in ["mvp_eval_provider_adapters", "mvp_required_version_slices"]:
        require(stale not in serialized, f"active execution plan contains stale field or worker {stale}")


def validate_task_dependencies(tasks: dict[str, dict[str, Any]]) -> None:
    for task_id, task in tasks.items():
        deps = task.get("blocked_by")
        require(isinstance(deps, list), f"{task_id}.blocked_by must be a list")
        unknown = set(str(dep) for dep in deps) - set(tasks)
        require(not unknown, f"{task_id}.blocked_by references unknown tasks {sorted(unknown)}")
        require(task_id not in deps, f"{task_id} cannot depend on itself")
    visiting: set[str] = set()
    visited: set[str] = set()

    def visit(task_id: str) -> None:
        if task_id in visited:
            return
        require(task_id not in visiting, f"task dependency cycle contains {task_id}")
        visiting.add(task_id)
        for dependency in tasks[task_id]["blocked_by"]:
            visit(str(dependency))
        visiting.remove(task_id)
        visited.add(task_id)

    for task_id in tasks:
        visit(task_id)


def validate_task(task: dict[str, Any], required_ids: set[str]) -> None:
    task_id = str(task.get("task_id", ""))
    missing = sorted(REQUIRED_TASK_FIELDS - set(task))
    require(not missing, f"{task_id} missing runner/planning fields: {missing}")
    require(task.get("worker_type") == "codex_cli", f"{task_id}.worker_type must be codex_cli")
    require(task.get("risk_class") in {"medium", "high"}, f"{task_id} risk class must be medium or high")
    validate_runtime_pins(task.get("runtime_pins"), task_id)
    high_risk = task.get("risk_class") == "high"
    required_review = task.get("required_review", {})
    require(required_review.get("required") is high_risk, f"{task_id}.required_review must match risk class")
    if high_risk:
        require(required_review.get("review_type") in {"code", "security", "architecture"}, f"{task_id} high-risk review_type is invalid")
        require(nonempty_string(required_review.get("reviewer_class")), f"{task_id} high-risk reviewer_class is required")
    for field in ["allowed_paths", "forbidden_paths", "scope_exclusions", "input_artifacts", "commands", "baseline_commands", "red_first_commands", "validation_commands", "final_validation_commands", "acceptance_checks", "worker_evidence_required", "lifecycle_evidence_required", "stop_conditions", "architecture_target_paths", "semantic_invariants"]:
        require(list_of_strings(task.get(field)), f"{task_id}.{field} must be a non-empty string list")
    docs_sync_refs = task.get("docs_sync_refs")
    require(nonempty_list(docs_sync_refs) and all(isinstance(item, dict) and nonempty_string(item.get("path")) and nonempty_string(item.get("reason")) for item in docs_sync_refs), f"{task_id}.docs_sync_refs must contain path/reason objects")
    require(f"{PLAN_REL}/" in task["forbidden_paths"], f"{task_id} must forbid active planning artifacts")
    require(f"{HISTORICAL_PLAN_REL}/" in task["forbidden_paths"], f"{task_id} must forbid historical planning artifacts")
    require(".git/" in task["forbidden_paths"], f"{task_id} must forbid .git")
    require(f".factory/artifacts/pr-lifecycle/{task_id}/" in task["forbidden_paths"], f"{task_id} must forbid PR lifecycle evidence")
    require(f".factory/artifacts/lifecycle-evidence/{task_id}/" in task["forbidden_paths"], f"{task_id} must forbid independent lifecycle evidence")
    require(all(not str(path).startswith(f"{PLAN_REL}/") for path in task["allowed_paths"]), f"{task_id} allows writes to active control truth")
    require(
        all(
            not str(path).startswith((".factory/artifacts/pr-lifecycle/", ".factory/artifacts/lifecycle-evidence/"))
            for path in task["allowed_paths"]
        ),
        f"{task_id} allows implementation writes to lifecycle-owned evidence",
    )
    ids = task.get("acceptance_item_ids")
    require(list_of_strings(ids), f"{task_id}.acceptance_item_ids must be non-empty")
    require(len(ids) <= 15, f"{task_id} has more than 15 acceptance items; split the task")
    require(not (set(ids) - required_ids), f"{task_id} references unknown acceptance items")
    require(task.get("acceptance_ledger_ref") == ARTIFACT_REFS["ledger"], f"{task_id} ledger ref is stale")
    inherited = task.get("validation_contract_inheritance", {})
    require(inherited.get("acceptance_ledger_ref") == ARTIFACT_REFS["ledger"], f"{task_id} inherited ledger ref is stale")
    require(set(inherited.get("acceptance_item_ids", [])) == set(ids), f"{task_id} inherited acceptance IDs must equal task IDs")
    require(inherited.get("required_review") == required_review, f"{task_id} inherited review requirement must exactly match the task")
    require("proof-of-behavior-scorecard" in task["worker_evidence_required"], f"{task_id} requires a proof scorecard")
    for evidence in ["ship_packet", "pr_lifecycle_report", "post_merge_report", "scope_closure_report", "factoryd_run_once_report"]:
        require(evidence in task["lifecycle_evidence_required"], f"{task_id} lifecycle evidence missing {evidence}")
    gates = task.get("lifecycle_gates", {})
    for gate in ["local_validation_required", "ci_required", "codex_review_required", "commit_push_required", "post_merge_monitor_required", "pr_lifecycle_report_required"]:
        require(gates.get(gate) is True, f"{task_id} lifecycle gate {gate} must be true")
    expected_chain = ["task-executor", "validation-gate"]
    if gates.get("code_review_required") is True:
        expected_chain.append("code-review")
    for gate, worker in [
        ("holdout_provisioning_required", "holdout-evaluator"),
        ("holdout_evaluation_required", "holdout-evaluator"),
        ("trace_grading_required", "trace-grader"),
        ("evidence_attestation_required", "evidence-attestor"),
    ]:
        if gates.get(gate) is True and worker not in expected_chain:
            expected_chain.append(worker)
    expected_chain.extend(["commit-push", "post-merge-monitor"])
    require(task.get("required_worker_chain") == expected_chain, f"{task_id} worker chain does not match review posture")
    require("ship-pr" not in task["required_worker_chain"], f"{task_id} uses deprecated ship-pr worker")
    require(task.get("alignment_gate_ref") == f"{ARTIFACT_REFS['plan']}#/alignment_gate", f"{task_id} alignment ref is stale")
    require(task.get("plan_drift_policy_ref") == f"{ARTIFACT_REFS['plan']}#/plan_drift_policy", f"{task_id} drift ref is stale")
    require(task.get("max_iterations") == 2, f"{task_id}.max_iterations must be 2")
    require(task.get("proof_scorecard_required") is True, f"{task_id} must require a proof scorecard")
    require(task.get("factory_compatibility", {}).get("profile_ref") == "profiles/lumyn.yaml", f"{task_id} profile ref is stale")
    runtime = task.get("factoryd_runtime", {})
    require(runtime.get("worker_type") == "codex_cli", f"{task_id} runtime worker must be codex_cli")
    require(runtime.get("network_posture", "").startswith("offline by default"), f"{task_id} runtime must be offline by default")
    required_capabilities = set(task.get("requires_capabilities", []))
    expected_capabilities = EXPECTED_CAPABILITIES.get(task_id, set())
    require(required_capabilities == expected_capabilities, f"{task_id} capability requirements differ: {sorted(required_capabilities)}")
    product_authorities = set(task.get("product_authority_requirements", []))
    expected_product_authorities = EXPECTED_PRODUCT_AUTHORITIES.get(task_id, set())
    require(product_authorities == expected_product_authorities, f"{task_id} product authority requirements differ: {sorted(product_authorities)}")
    optional_product_authorities = set(task.get("optional_product_action_capabilities", []))
    expected_optional_product_authorities = EXPECTED_OPTIONAL_PRODUCT_AUTHORITIES.get(task_id, set())
    require(optional_product_authorities == expected_optional_product_authorities, f"{task_id} optional product authority requirements differ: {sorted(optional_product_authorities)}")
    require(task.get("requires_human_approval") is bool(expected_capabilities or expected_product_authorities), f"{task_id}.requires_human_approval must match live capability scope")
    expected_credentials = "credentials" in expected_capabilities
    expected_network = "network" in expected_capabilities
    require(task.get("requires_credentials") is expected_credentials, f"{task_id}.requires_credentials must match live capability scope")
    require(task.get("requires_network") is expected_network, f"{task_id}.requires_network must match live capability scope")
    grants = runtime.get("capability_grants")
    require(isinstance(grants, list), f"{task_id}.factoryd_runtime.capability_grants must be a list")
    grant_capabilities = {grant.get("capability") for grant in grants if isinstance(grant, dict)}
    require(grant_capabilities == expected_capabilities, f"{task_id} placeholder grants differ from required capabilities")
    for index, grant in enumerate(grants):
        require(grant.get("task_id") == task_id, f"{task_id} grant {index} task_id mismatch")
        require(grant.get("approved") is False, f"{task_id} planning-time capability grants must remain unapproved")
        require_repo_ref(grant.get("evidence_ref"), f"{task_id}.grant[{index}].evidence_ref", allow_pending=True)
    control = runtime.get("runtime_control", {})
    require(control.get("autonomy_level") == "trusted", f"{task_id} Factory lifecycle runtime control must be trusted and path-bounded")
    require(set(control.get("max_write_scope_paths", [])) == set(task["allowed_paths"]) - {f".factory/artifacts/task-runs/{task_id}/"}, f"{task_id} runtime control paths must equal implementation paths")
    require(not contains_true_key(task, {"auto_merge", "default_branch_write", "provider_raw_repo_access"}), f"{task_id} enables a forbidden product capability")
    for field in ["test_matrix_refs", "ci_lane_refs", "engineering_policy_refs", "architecture_guidance_refs"]:
        require(nonempty_list(task.get(field)) and all(isinstance(item, dict) for item in task[field]), f"{task_id}.{field} must contain policy objects")
    matrix = json.dumps(task["test_matrix_refs"])
    require("docs/dev/dev_guides.md#12-level-test-matrix" in matrix, f"{task_id} lacks test-matrix propagation")
    lanes = {str(item.get("lane")) for item in task["ci_lane_refs"]}
    require({"fast", "core", "acceptance", "risk"}.issubset(lanes), f"{task_id} lacks CI-lane coverage")
    guide_refs = json.dumps(task["architecture_guidance_refs"])
    for token in ["trust-and-data-planes", "patch-safety-boundary", "command-execution-boundary", "systems-thinking-map", "adr-0002"]:
        require(token in guide_refs, f"{task_id} architecture guidance missing {token}")
    policy_refs = json.dumps(task["engineering_policy_refs"])
    for token in ["docs_parity", "structured_data", "migration_corpus", "patch_boundary", "proof", "provenance"]:
        require(token in policy_refs, f"{task_id} engineering policy missing {token}")
    coverage = task.get("coverage_policy_refs", {})
    require(coverage.get("required") is True, f"{task_id} coverage policy must be required")
    minimums = coverage.get("minimums")
    require(isinstance(minimums, list) and {item.get("threshold") for item in minimums if isinstance(item, dict)} == {75, 85}, f"{task_id} coverage thresholds must be 75 and 85")
    require(not contains_machine_local_path(task), f"{task_id} contains a machine-local path")


def validate_packets(packets: dict[str, Any], required_ids: set[str]) -> dict[str, dict[str, Any]]:
    require(packets.get("artifact_type") == "task_packets", "task-packets artifact_type must be task_packets")
    validate_runtime_pins(packets.get("runtime_pins"), "task packet set")
    require(packets.get("source_prd_ref") == "docs/product/prd.md", "task packets must cite the PRD")
    require(packets.get("source_ref") == ARTIFACT_REFS["plan"], "task packets must cite the execution plan")
    require(packets.get("acceptance_ledger_ref") == ARTIFACT_REFS["ledger"], "task packets ledger ref is stale")
    require(packets.get("acceptance_mapping_ref") == ARTIFACT_REFS["mapping"], "task packets mapping ref is stale")
    require(packets.get("acceptance_item_count") == len(required_ids), "task packet item count is stale")
    tasks_value = packets.get("tasks")
    require(isinstance(tasks_value, list), "task-packets.tasks must be a list")
    tasks: dict[str, dict[str, Any]] = {}
    for task in tasks_value:
        require(isinstance(task, dict), "task packet entries must be objects")
        task_id = task.get("task_id")
        require(nonempty_string(task_id), "task packet missing task_id")
        require(task_id not in tasks, f"duplicate task packet {task_id}")
        tasks[str(task_id)] = task
    require(tuple(tasks) == EXPECTED_TASK_IDS, f"task order/IDs must be {list(EXPECTED_TASK_IDS)}")
    validate_task_dependencies(tasks)
    for task in tasks.values():
        validate_task(task, required_ids)
    covered = set().union(*(set(task["acceptance_item_ids"]) for task in tasks.values()))
    require(covered == required_ids, f"active tasks do not cover exact acceptance set: missing={sorted(required_ids-covered)}")
    validate_migration_task_contracts(tasks)
    return tasks


def validate_contract(contract: dict[str, Any], required_ids: set[str]) -> None:
    require(contract.get("artifact_type") == "validation_contract", "validation contract artifact_type must be validation_contract")
    require(contract.get("acceptance_ledger_ref") == ARTIFACT_REFS["ledger"], "validation contract ledger ref is stale")
    require(contract.get("acceptance_mapping_ref") == ARTIFACT_REFS["mapping"], "validation contract mapping ref is stale")
    require(contract.get("acceptance_item_count") == len(required_ids), "validation contract item count is stale")
    validate_runtime_pins(contract.get("runtime_pins"), "validation contract")
    criteria = contract.get("acceptance_criteria")
    require(isinstance(criteria, list) and len(criteria) == len(required_ids), "validation contract must contain one criterion per acceptance item")
    criteria_ids = {str(value).split(":", 1)[0] for value in criteria}
    require(criteria_ids == required_ids, "validation criteria IDs differ from the PRD")
    checks = contract.get("required_checks")
    for check in ["make lint-fast", "make test-fast", "make test-coverage", "make test-contracts", "make prepush-full", "GitHub Actions validate", "GitHub Actions CodeQL analyze", "passive Codex review settle"]:
        require(check in checks, f"validation contract missing required check {check}")
    require(contract.get("required_review", {}).get("required") is True, "validation contract must require risk-based independent review")
    for evidence in ["validation_report", "work_proof_marker", "proof-of-behavior-scorecard", "ship_packet", "pr_lifecycle_report", "post_merge_report", "scope_closure_report"]:
        require(evidence in contract.get("evidence_required", []), f"validation contract missing evidence {evidence}")
    authority = contract.get("authority_requirements", {})
    require(set(authority.get("factory_worker_capabilities", [])) == {"approval", "credentials", "network"}, "validation contract Factory capabilities drifted")
    require(set(authority.get("separate_product_capabilities", [])) == PRODUCT_AUTHORITY_CAPABILITIES, "validation contract product authority capabilities drifted")
    require("never substitutes" in authority.get("separation_rule", ""), "validation contract must separate Factory and product authority")
    forbidden = set(authority.get("forbidden", []))
    require({"provider_implied_customer_authority", "default_branch_write", "auto_merge", "arbitrary_provider_script", "default_raw_customer_data_sharing"}.issubset(forbidden), "validation contract forbidden authority states are incomplete")
    validate_slices(contract.get("delivery_slices"), required_ids, "validation_contract.delivery_slices")


def validate_closure(closure: dict[str, Any], ledger_by_id: dict[str, dict[str, Any]], required_ids: set[str]) -> None:
    require(closure.get("artifact_type") == "scope_closure_map", "closure artifact_type must be scope_closure_map")
    require(closure.get("acceptance_ledger_ref") == ARTIFACT_REFS["ledger"], "closure ledger ref is stale")
    require(closure.get("acceptance_mapping_ref") == ARTIFACT_REFS["mapping"], "closure mapping ref is stale")
    require(closure.get("acceptance_item_count") == len(required_ids), "closure item count is stale")
    items = closure.get("items")
    require(isinstance(items, list), "closure.items must be a list")
    by_id: dict[str, dict[str, Any]] = {}
    for item in items:
        require(isinstance(item, dict), "closure item must be an object")
        item_id = item.get("scope_item_id")
        require(nonempty_string(item_id), "closure item missing scope_item_id")
        require(item_id not in by_id, f"duplicate closure item {item_id}")
        by_id[str(item_id)] = item
    require(set(by_id) == required_ids, "closure map does not cover exact PRD acceptance IDs")
    for item_id, item in by_id.items():
        ledger_item = ledger_by_id[item_id]
        expected_item_status = "missing" if ledger_item.get("status") == "planned" else ledger_item.get("status")
        require(item.get("status") == expected_item_status, f"{item_id} closure status differs from ledger")
        require(item.get("acceptance_item_ids") == [item_id], f"{item_id} closure must be item-granular")
        statuses = item.get("acceptance_item_statuses")
        require(isinstance(statuses, list) and len(statuses) == 1 and statuses[0].get("acceptance_item_id") == item_id, f"{item_id} nested status is invalid")
        require(statuses[0].get("status") == expected_item_status, f"{item_id} nested status differs from ledger")
        if ledger_item["status"] in TERMINAL_EVIDENCE_STATUSES:
            require(item.get("evidence_refs") == ledger_item.get("evidence_refs"), f"{item_id} closure evidence differs from ledger")
            if ledger_item["status"] == "implemented":
                require(item.get("implemented_task_refs") == [], f"{item_id} cannot claim an unexecuted active task as implemented")
                require(set(item.get("remaining_task_refs", [])) == set(item.get("task_refs", [])), f"{item_id} carried evidence must leave active preservation tasks remaining")
            else:
                require(item.get("remaining_task_refs") == [], f"{item_id} approved terminal disposition cannot retain remaining tasks")
        else:
            require(nonempty_list(item.get("remaining_task_refs")), f"{item_id} non-terminal closure requires remaining tasks")
    slice_map = validate_slices(closure.get("delivery_slices"), required_ids, "scope_closure_map.delivery_slices")
    for slice_id, item in slice_map.items():
        expected_remaining = set(item["acceptance_item_ids"]) - {item_id for item_id, ledger_item in ledger_by_id.items() if ledger_item["status"] in TERMINAL_EVIDENCE_STATUSES}
        require(set(item.get("remaining_acceptance_item_ids", [])) == expected_remaining, f"{slice_id} remaining acceptance IDs are stale")
    terminal = all(item["status"] in {"implemented", "deferred_with_approval", "out_of_scope"} for item in ledger_by_id.values())
    expected_closure = "closed" if terminal else "partial"
    require(closure.get("closure_status") == expected_closure, f"closure_status must be {expected_closure}")


def validate_config_payload(config: dict[str, Any], label: str, *, autoship: bool) -> None:
    require(config.get("factory", {}).get("profile_path") == "profiles/lumyn.yaml", f"{label} profile ref is stale")
    repo = config.get("repos", {}).get("lumyn")
    require(isinstance(repo, dict), f"{label} missing repos.lumyn")
    require(repo.get("task_packets") == ARTIFACT_REFS["packets"], f"{label} task_packets ref is stale")
    require(repo.get("scope_closure_map") == ARTIFACT_REFS["closure"], f"{label} closure ref is stale")
    require(repo.get("validation_contract") == ARTIFACT_REFS["contract"], f"{label} contract ref is stale")
    require(repo.get("acceptance_ledger") == ARTIFACT_REFS["ledger"], f"{label} ledger ref is stale")
    require(repo.get("capability_grants") == [], f"{label} checked-in template must not grant live product capabilities")
    require(repo.get("auto_ship") is autoship, f"{label} auto_ship must be {autoship}")
    posture = " ".join(str(repo.get(field, "")) for field in ["approval_posture", "credential_posture", "network_posture"]).lower()
    for token in ["factory worker grants", "factory dispatch does not validate or confer private lumyn product authority", "no ambient", "offline by default", "task-scoped"]:
        require(token in posture, f"{label} posture missing {token}")
    commands = repo.get("validation_commands")
    require(isinstance(commands, list) and "python3 scripts/validate_repo_pack.py" in commands, f"{label} must run repo-pack validation")
    if autoship:
        shipping = repo.get("shipping", {})
        for field in ["enabled", "push_required", "pr_required", "ci_required", "codex_review_required", "merge_required", "post_merge_required", "scope_closure_required"]:
            require(shipping.get(field) is True, f"{label}.shipping.{field} must be true")
    validate_architecture_budget_policy(ROOT, repo, label)


def validate_loaded(data: dict[str, dict[str, Any]], *, validate_configs: bool = True) -> dict[str, dict[str, Any]]:
    required_ids = expected_acceptance_ids()
    validate_acceptance_text(PRD.read_text(), data["ledger"])
    for name, payload in data.items():
        require(not contains_machine_local_path(payload), f"{name} contains a machine-local path")
        require(not contains_true_key(payload, {"auto_merge", "default_branch_write", "provider_raw_repo_access"}), f"{name} enables a forbidden product capability")
    validate_context(data["context"])
    validate_risk(data["risk"])
    ledger_by_id = validate_ledger(data["ledger"], required_ids)
    validate_mapping(data["mapping"], required_ids)
    validate_plan(data["plan"], required_ids)
    tasks = validate_packets(data["packets"], required_ids)
    validate_contract(data["contract"], required_ids)
    validate_closure(data["closure"], ledger_by_id, required_ids)
    plan_slices = {item["slice_id"]: set(item["acceptance_item_ids"]) for item in data["plan"]["delivery_slices"]}
    for name in ["packets", "contract", "mapping", "closure"]:
        other = {item["slice_id"]: set(item["acceptance_item_ids"]) for item in data[name]["delivery_slices"]}
        require(other == plan_slices, f"{name} delivery slices differ from execution plan")
    if validate_configs:
        validate_config_payload(data["config"], ".factory/factoryd.example.json", autoship=False)
        validate_config_payload(data["autoship"], ".factory/factoryd.autoship.example.json", autoship=True)
    return tasks


def validate_active_config(config: dict[str, Any], tasks: dict[str, dict[str, Any]]) -> None:
    repo = config.get("repos", {}).get("lumyn")
    require(isinstance(repo, dict), ".factory/factoryd.json missing repos.lumyn")
    require(repo.get("task_packets") == ARTIFACT_REFS["packets"], ".factory/factoryd.json task_packets ref is stale")
    require(repo.get("scope_closure_map") == ARTIFACT_REFS["closure"], ".factory/factoryd.json closure ref is stale")
    require(repo.get("validation_contract") == ARTIFACT_REFS["contract"], ".factory/factoryd.json validation contract ref is stale")
    require(repo.get("acceptance_ledger") == ARTIFACT_REFS["ledger"], ".factory/factoryd.json ledger ref is stale")
    validate_authority_grants(repo.get("capability_grants"), tasks, EXPECTED_CAPABILITIES)


def load_all() -> dict[str, dict[str, Any]]:
    data = {name: load_json(path) for name, path in ARTIFACT_PATHS.items()}
    data["config"] = load_json(CONFIG)
    data["autoship"] = load_json(AUTOSHIP_CONFIG)
    return data


def run_self_test() -> None:
    run_repo_pack_self_tests(
        load_all(),
        validate_loaded=validate_loaded,
        validate_config_payload=validate_config_payload,
        validate_active_config=validate_active_config,
        historical_plan_rel=HISTORICAL_PLAN_REL,
        expected_capabilities=EXPECTED_CAPABILITIES,
    )


def main() -> int:
    if sys.argv[1:] == ["--self-test"]:
        try:
            run_self_test()
        except AssertionError as exc:
            print(f"repo-pack validator self-test failed: {exc}", file=sys.stderr)
            return 2
        print("repo-pack validator self-test passed")
        return 0
    if sys.argv[1:]:
        print("usage: validate_repo_pack.py [--self-test]", file=sys.stderr)
        return 2
    try:
        validate_docs()
        validate_ci_controls()
        tasks = validate_loaded(load_all())
        if ACTIVE_CONFIG.exists():
            active = load_json(ACTIVE_CONFIG)
            validate_active_config(active, tasks)
    except AssertionError as exc:
        print(f"repo-pack validation failed: {exc}", file=sys.stderr)
        return 2
    print("repo-pack validation passed")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
