#!/usr/bin/env python3
"""Validate the paused Lumyn v3.1 planning and Factory control generation."""

from __future__ import annotations

import json
import re
import sys
from pathlib import Path
from typing import Any

from repo_pack_architecture import validate_architecture_budget_policy
from repo_pack_validation.acceptance_text import validate_acceptance_text
from repo_pack_validation.authority import (
    validate_active_repo_safety,
    validate_authority_grants,
)
from repo_pack_validation.markdown_refs import validate_markdown_fragment_refs
from repo_pack_validation.runtime_pins import validate_runtime_pins
from repo_pack_validation.self_tests import run_repo_pack_self_tests
from repo_pack_validation.task_contracts import (
    TASK_DEPENDENCIES,
    validate_implemented_proof,
    validate_migration_task_contracts,
)


ROOT = Path(__file__).resolve().parents[2]
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
    name: path.relative_to(ROOT).as_posix()
    for name, path in ARTIFACT_PATHS.items()
}

EXPECTED_TASK_IDS = tuple(TASK_DEPENDENCIES)
MILESTONE_REFS = {
    "M0": "docs/product/plan.md#m0-correct-command-and-result-foundations",
    "M1": "docs/product/plan.md#m1-build-deterministic-agent-assisted-blocked-and-generic-agent-benchmarks",
    "M2": "docs/product/plan.md#m2-define-update-installation-agent-verification-delivery-and-privacy-contracts",
    "M2.5": "docs/product/plan.md#m25-pre-sell-and-qualify-one-provider-campaign",
    "M3": "docs/product/plan.md#m3-normalize-provider-intent-into-a-reusable-provider-change-contract",
    "M4": "docs/product/plan.md#m4-build-the-consumer-local-typescript-repository-impact-inventory",
    "M5": "docs/product/plan.md#m5-produce-a-no-write-migration-plan-and-consumer-approval-boundary",
    "M6": "docs/product/plan.md#m6-implement-deterministic-transforms-and-the-bounded-coding-agent",
    "M7": "docs/product/plan.md#m7-verify-deterministic-agent-assisted-and-manual-candidates",
    "M8": "docs/product/plan.md#m8-add-optional-provider-sandbox-read-back",
    "M9": "docs/product/plan.md#m9-export-evidence-and-open-a-tested-draft-pr",
    "M10": "docs/product/plan.md#m10-run-one-prepaid-provider-originated-update-campaign",
}
EXPECTED_GROUP_IDS = {"retained_foundation",
    "migration_pack",
    "provider_change_event", "consumer_installation",
    "impact_and_integration_graph",
    "plan_and_routing",
    "bounded_hybrid_execution",
    "verification",
    "export_and_delivery",
    "trust_and_privacy",
    "design_partner_qualification",
    "provider_campaign_pilot",
}
EXPECTED_SLICE_IDS = {
    "migration_foundation",
    "hybrid_migration_engine",
    "provider_campaign",
}
ALLOWED_HISTORICAL_TASK_IDS = {"T1", "T2", "T2.7", "T3"}
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
EXPECTED_IMPLEMENTED_ITEMS = {
    "BASE-001",
    "BASE-002",
    "BASE-003",
    "BASE-004",
    "BASE-005",
}

EXPECTED_CAPABILITIES = {
    "M2.5": {"approval"},
    "M8": {"approval", "credentials", "network"},
    "M9": {"approval", "credentials", "network"},
    "M10": {"approval", "credentials", "network"},
}
EXPECTED_CONDITIONAL_CAPABILITIES = {
    "M6": {"approval", "credentials", "network"},
}
EXPECTED_PRODUCT_AUTHORITIES = {
    "M4": {
        "customer_repo_read",
        "artifact_retention",
        "artifact_deletion",
    },
    "M6": {
        "customer_repo_read",
        "customer_repo_write",
        "command_execution",
        "artifact_retention",
        "artifact_deletion",
    },
    "M7": {
        "customer_repo_read",
        "command_execution",
        "artifact_retention",
        "artifact_deletion",
    },
    "M8": {
        "customer_repo_read",
        "command_execution",
        "sandbox_request_disclosure",
        "sandbox_network",
        "sandbox_credential",
        "artifact_retention",
        "artifact_deletion",
    },
    "M9": {
        "customer_repo_read",
        "customer_repo_write",
        "command_execution",
        "github_branch_write",
        "github_pr_write",
        "artifact_retention",
        "artifact_deletion",
    },
    "M10": {
        "customer_repo_read",
        "customer_repo_write",
        "command_execution",
        "github_branch_write",
        "github_pr_write",
        "provider_reporting",
        "artifact_retention",
        "artifact_deletion",
    },
}
EXPECTED_OPTIONAL_PRODUCT_AUTHORITIES = {
    "M6": {
        "agent_runner_credential",
        "agent_runner_network",
        "model_request_disclosure",
        "model_network",
        "model_credential",
        "package_registry_read",
    },
    "M7": {
        "customer_repo_write",
        "agent_runner_credential",
        "agent_runner_network",
        "model_request_disclosure",
        "model_network",
        "model_credential",
        "package_registry_read",
    },
    "M8": set(),
    "M9": {
        "provider_reporting",
    },
    "M10": {
        "agent_runner_credential",
        "agent_runner_network",
        "model_request_disclosure",
        "model_network",
        "model_credential",
        "package_registry_read",
        "sandbox_request_disclosure",
        "sandbox_network",
        "sandbox_credential",
    },
}
PRODUCT_AUTHORITY_CAPABILITIES = (
    set().union(*EXPECTED_PRODUCT_AUTHORITIES.values())
    | set().union(*EXPECTED_OPTIONAL_PRODUCT_AUTHORITIES.values())
)

PRODUCT_AUTHORITY_SCOPE_SEMANTICS = (
    "task_level_capability_universe_not_runtime_grant; "
    "named_route_exact_union_governs_each_action"
)

CONDITIONAL_CAPABILITY_PREDICATES = {
    "agent_runner_credential": (
        "agent_route_topology_is_runner_mediated_or_hybrid"
    ),
    "agent_runner_network": (
        "agent_route_topology_is_runner_mediated_or_hybrid"
    ),
    "model_credential": "agent_route_topology_is_direct_model_or_hybrid",
    "model_network": "agent_route_topology_is_direct_model_or_hybrid",
    "model_request_disclosure": "agent_route_topology_is_not_local_runtime",
    "package_registry_read": (
        "approved_action_requires_registry_or_immutable_snapshot_access"
    ),
}

AGENT_ROUTE_AUTHORIZATION_TOPOLOGY_CONTRACT = {
    "selector_field": "agent_route_topology",
    "allowed_values": [
        "local_runtime",
        "runner_mediated",
        "direct_model",
        "hybrid",
    ],
    "exactly_one_topology_selected": True,
    "minimum_capability_sets": {
        "local_runtime": [],
        "runner_mediated": [
            "agent_runner_credential",
            "agent_runner_network",
            "model_request_disclosure",
        ],
        "direct_model": [
            "model_credential",
            "model_network",
            "model_request_disclosure",
        ],
        "hybrid": [
            "agent_runner_credential",
            "agent_runner_network",
            "model_credential",
            "model_network",
            "model_request_disclosure",
        ],
    },
    "selected_topology_minimum_capabilities_required": True,
    "unselected_topology_capabilities_authorized": False,
    "local_runtime_external_egress_allowed": False,
}


def _product_route(
    required: tuple[str, ...],
    conditional: tuple[str, ...] = (),
    delegated_route_ref: str | None = None,
    *,
    agent_topology: bool = False,
) -> dict[str, Any]:
    route: dict[str, Any] = {
        "required_capabilities": sorted(required),
        "conditionally_selected_capabilities": sorted(conditional),
        "exact_selected_union_frozen_before_action": True,
        "unselected_capabilities_authorized": False,
    }
    if conditional:
        route["conditional_capability_predicates"] = {
            capability: CONDITIONAL_CAPABILITY_PREDICATES[capability]
            for capability in sorted(conditional)
        }
    if agent_topology:
        route["authorization_topology_contract"] = (
            AGENT_ROUTE_AUTHORIZATION_TOPOLOGY_CONTRACT
        )
    if delegated_route_ref is not None:
        route["delegated_route_ref"] = delegated_route_ref
    return route


EXPECTED_PRODUCT_ACTION_ROUTES = {
    "M4": {
        "impact_read_only": _product_route(
            ("artifact_deletion", "artifact_retention", "customer_repo_read")
        ),
    },
    "M6": {
        "deterministic_candidate": _product_route(
            (
                "artifact_deletion",
                "artifact_retention",
                "customer_repo_read",
                "customer_repo_write",
            )
        ),
        "deterministic_package_tool_candidate": _product_route(
            (
                "artifact_deletion",
                "artifact_retention",
                "command_execution",
                "customer_repo_read",
                "customer_repo_write",
            ),
            ("package_registry_read",),
        ),
        "agent_assisted_candidate": _product_route(
            (
                "artifact_deletion",
                "artifact_retention",
                "command_execution",
                "customer_repo_read",
                "customer_repo_write",
            ),
            (
                "agent_runner_credential",
                "agent_runner_network",
                "model_credential",
                "model_network",
                "model_request_disclosure",
                "package_registry_read",
            ),
            agent_topology=True,
        ),
    },
    "M7": {
        "verify": _product_route(
            (
                "artifact_deletion",
                "artifact_retention",
                "command_execution",
                "customer_repo_read",
            ),
            ("package_registry_read",),
        ),
        "repair_agent_assisted": _product_route(
            (
                "artifact_deletion",
                "artifact_retention",
                "command_execution",
                "customer_repo_read",
                "customer_repo_write",
            ),
            (
                "agent_runner_credential",
                "agent_runner_network",
                "model_credential",
                "model_network",
                "model_request_disclosure",
                "package_registry_read",
            ),
            agent_topology=True,
        ),
    },
    "M8": {
        "sandbox_read_back": _product_route(
            (
                "artifact_deletion",
                "artifact_retention",
                "command_execution",
                "customer_repo_read",
                "sandbox_credential",
                "sandbox_network",
                "sandbox_request_disclosure",
            )
        ),
    },
    "M9": {
        "local_export": _product_route(
            (
                "artifact_deletion",
                "artifact_retention",
                "command_execution",
                "customer_repo_read",
            )
        ),
        "local_branch": _product_route(
            (
                "artifact_deletion",
                "artifact_retention",
                "command_execution",
                "customer_repo_read",
                "customer_repo_write",
            )
        ),
        "remote_branch_push": _product_route(
            (
                "artifact_deletion",
                "artifact_retention",
                "command_execution",
                "customer_repo_read",
                "github_branch_write",
            )
        ),
        "draft_pr_create": _product_route(
            ("artifact_deletion", "artifact_retention", "github_pr_write")
        ),
        "provider_status_decline": _product_route(
            ("artifact_deletion", "artifact_retention")
        ),
        "provider_status_transmit": _product_route(
            ("artifact_deletion", "artifact_retention", "provider_reporting")
        ),
    },
    "M10": {
        "impact_read_only": _product_route(
            ("artifact_deletion", "artifact_retention", "customer_repo_read"),
            delegated_route_ref="M4/impact_read_only",
        ),
        "deterministic_candidate": _product_route(
            (
                "artifact_deletion",
                "artifact_retention",
                "customer_repo_read",
                "customer_repo_write",
            ),
            delegated_route_ref="M6/deterministic_candidate",
        ),
        "deterministic_package_tool_candidate": _product_route(
            (
                "artifact_deletion",
                "artifact_retention",
                "command_execution",
                "customer_repo_read",
                "customer_repo_write",
            ),
            ("package_registry_read",),
            "M6/deterministic_package_tool_candidate",
        ),
        "agent_assisted_candidate": _product_route(
            (
                "artifact_deletion",
                "artifact_retention",
                "command_execution",
                "customer_repo_read",
                "customer_repo_write",
            ),
            (
                "agent_runner_credential",
                "agent_runner_network",
                "model_credential",
                "model_network",
                "model_request_disclosure",
                "package_registry_read",
            ),
            "M6/agent_assisted_candidate",
            agent_topology=True,
        ),
        "verify": _product_route(
            (
                "artifact_deletion",
                "artifact_retention",
                "command_execution",
                "customer_repo_read",
            ),
            ("package_registry_read",),
            "M7/verify",
        ),
        "repair_agent_assisted": _product_route(
            (
                "artifact_deletion",
                "artifact_retention",
                "command_execution",
                "customer_repo_read",
                "customer_repo_write",
            ),
            (
                "agent_runner_credential",
                "agent_runner_network",
                "model_credential",
                "model_network",
                "model_request_disclosure",
                "package_registry_read",
            ),
            "M7/repair_agent_assisted",
            agent_topology=True,
        ),
        "sandbox_read_back": _product_route(
            (
                "artifact_deletion",
                "artifact_retention",
                "command_execution",
                "customer_repo_read",
                "sandbox_credential",
                "sandbox_network",
                "sandbox_request_disclosure",
            ),
            delegated_route_ref="M8/sandbox_read_back",
        ),
        "local_export": _product_route(
            (
                "artifact_deletion",
                "artifact_retention",
                "command_execution",
                "customer_repo_read",
            ),
            delegated_route_ref="M9/local_export",
        ),
        "local_branch": _product_route(
            (
                "artifact_deletion",
                "artifact_retention",
                "command_execution",
                "customer_repo_read",
                "customer_repo_write",
            ),
            delegated_route_ref="M9/local_branch",
        ),
        "remote_branch_push": _product_route(
            (
                "artifact_deletion",
                "artifact_retention",
                "command_execution",
                "customer_repo_read",
                "github_branch_write",
            ),
            delegated_route_ref="M9/remote_branch_push",
        ),
        "draft_pr_create": _product_route(
            ("artifact_deletion", "artifact_retention", "github_pr_write"),
            delegated_route_ref="M9/draft_pr_create",
        ),
        "provider_status_decline": _product_route(
            ("artifact_deletion", "artifact_retention"),
            delegated_route_ref="M9/provider_status_decline",
        ),
        "provider_status_transmit": _product_route(
            ("artifact_deletion", "artifact_retention", "provider_reporting"),
            delegated_route_ref="M9/provider_status_transmit",
        ),
    },
}
FACTORY_CAPABILITIES = {"approval", "credentials", "network"}
REQUIRED_TASK_FIELDS = {
    "task_id",
    "milestone_ref",
    "objective",
    "risk_class",
    "dispatch_status",
    "blocked_by",
    "allowed_paths",
    "forbidden_paths",
    "scope_exclusions",
    "acceptance_item_ids",
    "acceptance_checks",
    "acceptance_ledger_ref",
    "validation_contract_inheritance",
    "required_worker_chain",
    "lifecycle_gates",
    "worker_evidence_required",
    "lifecycle_evidence_required",
    "requires_capabilities",
    "conditional_factory_capabilities",
    "product_authority_requirements",
    "optional_product_action_capabilities",
    "requires_human_approval",
    "requires_credentials",
    "requires_network",
    "runtime_pins",
    "factory_compatibility",
    "factoryd_runtime",
    "validation_commands",
    "stop_conditions",
}
REQUIRED_DOCS = {
    "README.md": [
        "provider-originated change event",
        "services-assisted",
        "provider-paid",
        "bounded agent",
        "customer-selected Agent Runner",
        "Planning-only, not implemented",
        HISTORICAL_PLAN_REL,
    ],
    "AGENTS.md": [
        "Two Principals, Two Authorities",
        "model_request_disclosure",
        "agent_runner_network",
        "agent_runner_credential",
        "model_network",
        "model_credential", "provider_reporting",
        "Passive Codex review settle is required before merge",
    ],
    "WORKFLOW.md": [
        "provider-confirmed",
        "bounded agent",
        "Agent Runner Selection And Funding",
        "PR bundle",
        "Green CI alone is not merge-ready",
        "process escape", "provider_reporting",
    ],
    "docs/product/prd.md": [
        "Provider-Originated API Update Delivery",
        "Consumer Installation",
        "Provider Change Contract",
        "provider-paid",
        "bounded coding agent",
        "customer-selected",
        "provider_sponsored_lumyn_managed",
        "generic coding agent",
        "53 item-level closure units",
        "Falsification And Reframe Gates",
    ],
    "docs/product/plan.md": [
        "M0: Correct command and result foundations",
        "M6: Implement deterministic transforms and the bounded coding agent",
        "Codex and Claude Code",
        "M9: Export evidence and open a tested draft PR",
        "M10: Run one prepaid provider-originated update campaign",
        "all 53 PRD acceptance items",
    ],
    "docs/dev/dev_guides.md": [
        "services-assisted, provider-paid",
        "Bounded Agent Policy",
        "agent_runner_network",
        "agent_runner_credential",
        "deterministic verification", "provider_reporting",
        "Do not merge manually through `gh pr merge`",
    ],
    "docs/architecture/architecture_guides.md": [
        "Provider Change And Campaign Plane",
        "Consumer Execution Plane",
        "Bounded Model Plane",
        "Agent Runner Plane",
        "Patch Safety Boundary",
        "GitHub Boundary",
    ],
    "docs/architecture/adr-0003-services-led-bounded-agent-migration-execution.md": [
        "Services-Led Bounded-Agent Migration Execution",
        "Partially superseded by ADR-0004",
        "Deterministic Verification",
        "Delivery Ladder",
    ],
    "docs/architecture/adr-0004-provider-originated-api-update-delivery.md": [
        "Provider-Originated API Update Delivery",
        "Consumer Installation",
        "Provider Change Event",
        "Delivery And Status",
    ],
    "docs/architecture/adr-0005-customer-selected-agent-runners.md": [
        "Customer-Selected Agent Runners",
        "consumer_managed",
        "provider_sponsored_lumyn_managed",
        "Independent Verification",
    ],
    "docs/factory/README.md": [
        PLAN_REL,
        "Product Authority Is Not Factory Authority",
        "model_request_disclosure", "agent_runner_network",
        "agent_runner_credential", "provider_reporting",
        "immutable historical records",
    ],
    ".factory/README.md": [
        PLAN_REL,
        "factoryd dispatch remains paused",
        "immutable records", "Agent Runner", "provider_reporting",
    ],
}
MACHINE_LOCAL_RE = re.compile(
    r"(?:^|[\s\"'])(?:/Users/|/home/|file://|[A-Za-z]:\\\\)"
)
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
    require(
        isinstance(payload, dict),
        f"{path.relative_to(ROOT)} must contain a JSON object",
    )
    return payload


def nonempty_string(value: Any) -> bool:
    return isinstance(value, str) and bool(value.strip())


def list_of_strings(value: Any, *, allow_empty: bool = False) -> bool:
    return (
        isinstance(value, list)
        and (allow_empty or bool(value))
        and all(nonempty_string(item) for item in value)
    )


def source_path(value: str) -> str:
    return value.split("#", 1)[0]


def require_repo_ref(value: Any, label: str) -> None:
    require(nonempty_string(value), f"{label} must be a non-empty reference")
    ref = str(value).strip()
    if ref.startswith(("http://", "https://", "private:", "pending:")):
        return
    relative = Path(source_path(ref))
    require(not relative.is_absolute(), f"{label} must be repo-relative")
    require(".." not in relative.parts, f"{label} must stay inside the repository")
    require((ROOT / relative).exists(), f"{label} points to missing path {relative}")


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
        return any(
            (key in keys and item is True) or contains_true_key(item, keys)
            for key, item in value.items()
        )
    if isinstance(value, list):
        return any(contains_true_key(item, keys) for item in value)
    return False


def expected_acceptance_ids() -> set[str]:
    text = PRD.read_text()
    require(
        "## Acceptance Tests" in text and "## Success Metrics" in text,
        "PRD acceptance section is missing",
    )
    section = text.split("## Acceptance Tests", 1)[1].split(
        "## Success Metrics", 1
    )[0]
    ids = set(ACCEPTANCE_ID_RE.findall(section))
    require(
        len(ids) == 53,
        f"PRD must define exactly 53 unique acceptance IDs; found {len(ids)}",
    )
    return ids


def validate_docs() -> None:
    for relative, tokens in REQUIRED_DOCS.items():
        path = ROOT / relative
        require(path.exists(), f"missing required document: {relative}")
        text = path.read_text()
        for token in tokens:
            require(
                token in text,
                f"{relative} missing required v3 operating token: {token}",
            )
        require(
            not MACHINE_LOCAL_RE.search(text),
            f"{relative} contains a machine-local path",
        )
    require(
        (ROOT / HISTORICAL_PLAN_REL / "README.md").exists(),
        "historical lumyn-mvp plan must remain present",
    )
    require((PLAN_DIR / "README.md").exists(), "active v3 plan README is missing")


def validate_ci_controls() -> None:
    required_checks = load_json(ROOT / ".github/required-checks.json")
    serialized = json.dumps(required_checks)
    for check in ("validate", "CodeQL analyze"):
        require(check in serialized, f"required-check metadata missing {check}")
    codeowners = (ROOT / ".github/CODEOWNERS").read_text()
    for token in (
        "/.github/** @davidahmann",
        "/.factory/** @davidahmann",
        "/docs/product/** @davidahmann",
        "/schemas/** @davidahmann",
    ):
        require(token in codeowners, f"CODEOWNERS missing {token}")


def validate_pause_contract(value: Any, label: str) -> None:
    require(isinstance(value, dict), f"{label} must be an object")
    require(value.get("dispatch_enabled") is False, f"{label} dispatch must be disabled")
    require(value.get("fail_closed") is True, f"{label} must fail closed")
    require(
        value.get("profile_compatibility_status")
        == "blocked_v3_1_profile_update_required",
        f"{label} must record the v2 Factory profile incompatibility",
    )
    require(
        value.get("runtime_qualification_status")
        == "blocked_factoryd_v3_1_qualification_required",
        f"{label} must require factoryd v3.1 qualification",
    )
    dependencies = value.get("unblock_dependencies")
    require(
        isinstance(dependencies, list)
        and any("profiles/lumyn.yaml" in str(item) for item in dependencies)
        and any("factoryd" in str(item).lower() for item in dependencies),
        f"{label} must name exact profile and factoryd unblock dependencies",
    )


def validate_context(context: dict[str, Any]) -> None:
    require(
        context.get("artifact_type") == "context_brief",
        "context artifact_type must be context_brief",
    )
    require(
        context.get("source_prd_ref") == "docs/product/prd.md",
        "context must cite the canonical PRD",
    )
    alignment = context.get("alignment_decisions", {})
    require(alignment.get("status") == "blocked", "context alignment must be blocked")
    require(
        alignment.get("implementation_may_start") is False,
        "context implementation must remain blocked",
    )
    resolved = json.dumps(alignment.get("resolved", [])).lower()
    for token in (
        "provider-paid",
        "services-assisted",
        "provider change contract",
        "consumer installation",
        "bounded agent",
        "customer-selected qualified bounded agent runner",
        "codex and claude code",
        "agent_execution_policy",
        "deterministic verification",
        "local fallback",
        "tested draft pr",
        "consented rollout status",
    ):
        require(token in resolved, f"context decisions missing {token}")
    baseline = json.dumps(context.get("proven_findings", [])).lower()
    for token in (
        "openapi",
        "not implemented",
        "generic pass",
        "v2 factory profile",
        "factoryd",
    ):
        require(token in baseline, f"context baseline missing {token}")
    compatibility = context.get("factory_compatibility", {})
    require(
        compatibility.get("profile_ref") == "profiles/lumyn.yaml",
        "context Factory profile ref is stale",
    )
    require(
        compatibility.get("status") == "blocked_profile_and_runtime_unqualified",
        "context must block incompatible Factory dispatch",
    )
    validate_pause_contract(context.get("factoryd_runtime"), "context.factoryd_runtime")
    validate_runtime_pins(context.get("runtime_pins"), "context")


def validate_risk(risk: dict[str, Any]) -> None:
    require(
        risk.get("artifact_type") == "risk_classification",
        "risk artifact_type must be risk_classification",
    )
    require(risk.get("default_risk_class") == "high", "default risk must be high")
    serialized = json.dumps(risk).lower()
    for token in (
        "provider-confirmed provider change contract",
        "provider event authenticity",
        "consumer installation",
        "consumer repository read or write",
        "model request disclosure",
        "model endpoint network",
        "model credential",
        "prompt injection",
        "bounded agent",
        "agent runner qualification",
        "agent runner authentication",
        "usage-billing ownership",
        "deterministic verification",
        "github draft pr",
        "event-bound provider-visible reporting",
        "generic-agent baseline",
        "paid campaign",
    ):
        require(token in serialized, f"risk classification missing {token}")


def validate_ledger(
    ledger: dict[str, Any],
    required_ids: set[str],
) -> dict[str, dict[str, Any]]:
    require(
        ledger.get("artifact_type") == "acceptance_ledger",
        "ledger artifact_type must be acceptance_ledger",
    )
    require(
        ledger.get("source_ref") == "docs/product/prd.md",
        "ledger source_ref must cite the PRD",
    )
    policy = ledger.get("coverage_policy", {})
    require(
        policy.get("closure_unit") == "acceptance_item",
        "ledger closure unit must be acceptance_item",
    )
    require(
        policy.get("group_only_refs_allowed") is False,
        "group-only closure must be forbidden",
    )
    items = ledger.get("items")
    require(isinstance(items, list), "ledger.items must be a list")
    by_id: dict[str, dict[str, Any]] = {}
    for index, item in enumerate(items):
        require(isinstance(item, dict), f"ledger.items[{index}] must be an object")
        item_id = item.get("acceptance_item_id")
        require(nonempty_string(item_id), f"ledger.items[{index}] missing ID")
        require(item_id not in by_id, f"duplicate ledger item {item_id}")
        by_id[str(item_id)] = item
    require(
        set(by_id) == required_ids,
        "ledger IDs differ from the exact PRD acceptance set",
    )
    require(
        ledger.get("acceptance_item_count") == len(required_ids),
        "ledger acceptance count is stale",
    )
    for item_id, item in by_id.items():
        for field in (
            "group_id",
            "source_ref",
            "source_text",
            "kind",
            "evidence_mode",
            "task_refs",
            "validation_refs",
            "status",
            "risk_class",
            "notes",
        ):
            require(field in item, f"{item_id} missing {field}")
        require(
            item["status"] in ALLOWED_ITEM_STATUSES,
            f"{item_id} has invalid status {item['status']}",
        )
        require_repo_ref(item["source_ref"], f"{item_id}.source_ref")
        task_refs = item.get("task_refs")
        require(list_of_strings(task_refs), f"{item_id}.task_refs must be non-empty")
        unknown = (
            set(task_refs)
            - set(EXPECTED_TASK_IDS)
            - ALLOWED_HISTORICAL_TASK_IDS
        )
        require(not unknown, f"{item_id} references unknown tasks {sorted(unknown)}")
        require(
            any(ref in EXPECTED_TASK_IDS for ref in task_refs),
            f"{item_id} must have an active task owner",
        )
        if item["status"] in TERMINAL_EVIDENCE_STATUSES:
            require(
                list_of_strings(item.get("evidence_refs")),
                f"{item_id} terminal status requires evidence refs",
            )
            for index, ref in enumerate(item["evidence_refs"]):
                require_repo_ref(ref, f"{item_id}.evidence_refs[{index}]")
    actual_implemented = {
        item_id
        for item_id, item in by_id.items()
        if item["status"] == "implemented"
    }
    require(
        actual_implemented == EXPECTED_IMPLEMENTED_ITEMS,
        "implemented status must match directly evidenced completed M0 items",
    )
    validate_implemented_proof(by_id)
    return by_id


def validate_slices(
    values: Any,
    required_ids: set[str],
    label: str,
) -> dict[str, dict[str, Any]]:
    require(isinstance(values, list), f"{label} must be a list")
    by_id: dict[str, dict[str, Any]] = {}
    covered: list[str] = []
    for item in values:
        require(isinstance(item, dict), f"{label} entries must be objects")
        slice_id = item.get("slice_id")
        require(nonempty_string(slice_id), f"{label} slice_id is required")
        require(slice_id not in by_id, f"{label} duplicate slice {slice_id}")
        by_id[str(slice_id)] = item
        require(
            item.get("required_for_completion") is True,
            f"{label}.{slice_id} must be required as a coverage lens",
        )
        require(
            list_of_strings(item.get("task_refs")),
            f"{label}.{slice_id}.task_refs must be non-empty",
        )
        require(
            list_of_strings(item.get("acceptance_item_ids")),
            f"{label}.{slice_id}.acceptance_item_ids must be non-empty",
        )
        covered.extend(item["acceptance_item_ids"])
    require(set(by_id) == EXPECTED_SLICE_IDS, f"{label} slice IDs drifted")
    require(
        set(covered) == required_ids and len(covered) == len(required_ids),
        f"{label} must cover each acceptance item exactly once",
    )
    return by_id


def validate_mapping(mapping: dict[str, Any], required_ids: set[str]) -> None:
    require(
        mapping.get("artifact_type") == "acceptance_mapping",
        "mapping artifact_type must be acceptance_mapping",
    )
    require(
        mapping.get("acceptance_ledger_ref") == ARTIFACT_REFS["ledger"],
        "mapping ledger ref is stale",
    )
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
        require(
            list_of_strings(ids),
            f"mapping group {group_id} must contain acceptance IDs",
        )
        mapped.extend(ids)
    require(group_ids == EXPECTED_GROUP_IDS, "mapping group taxonomy drifted")
    require(
        set(mapped) == required_ids and len(mapped) == len(required_ids),
        "mapping must cover the exact PRD acceptance set once",
    )
    validate_slices(mapping.get("delivery_slices"), required_ids, "mapping.slices")


def validate_plan(plan: dict[str, Any], required_ids: set[str]) -> None:
    require(
        plan.get("artifact_type") == "execution_plan",
        "execution plan artifact_type must be execution_plan",
    )
    for field, expected in (
        ("source_prd_ref", "docs/product/prd.md"),
        ("authored_plan_ref", "docs/product/plan.md"),
        ("context_brief_ref", ARTIFACT_REFS["context"]),
        ("risk_classification_ref", ARTIFACT_REFS["risk"]),
        ("acceptance_ledger_ref", ARTIFACT_REFS["ledger"]),
        ("acceptance_mapping_ref", ARTIFACT_REFS["mapping"]),
        ("validation_contract_ref", ARTIFACT_REFS["contract"]),
    ):
        require(plan.get(field) == expected, f"execution plan {field} is stale")
    require(plan.get("risk_class") == "high", "execution plan risk must be high")
    validate_runtime_pins(plan.get("runtime_pins"), "execution plan")
    alignment = plan.get("alignment_gate", {})
    require(alignment.get("status") == "blocked", "execution alignment must be blocked")
    require(
        alignment.get("implementation_may_start") is False,
        "execution implementation must remain blocked",
    )
    require(
        alignment.get("blocked_tasks") == list(EXPECTED_TASK_IDS),
        "execution alignment must block every v3 task",
    )
    validate_pause_contract(plan.get("factoryd_runtime"), "execution factoryd runtime")
    require(
        plan.get("dependency_graph") == TASK_DEPENDENCIES,
        "execution dependency graph differs from M0-M10",
    )
    require(
        plan.get("tasks") == list(EXPECTED_TASK_IDS),
        "execution task order differs from M0-M10",
    )
    retained = plan.get("retained_evidence_policy", {})
    require(
        retained.get("historical_plan_ref") == f"{HISTORICAL_PLAN_REL}/"
        and retained.get("status") == "immutable_non_active",
        "historical plan must remain immutable and non-active",
    )
    coverage = plan.get("acceptance_ledger_coverage", {})
    require(
        coverage.get("coverage_unit") == "acceptance_item"
        and coverage.get("required_item_count") == len(required_ids)
        and coverage.get("group_only_refs_allowed") is False,
        "execution acceptance coverage is stale",
    )
    locked_decisions = plan.get("locked_decisions")
    require(
        isinstance(locked_decisions, list)
        and all(nonempty_string(value) for value in locked_decisions)
        and len(locked_decisions) == len(set(locked_decisions)),
        "execution locked decisions must be unique non-empty strings",
    )
    locked = json.dumps(locked_decisions).lower()
    for token in (
        "provider-originated",
        "services-assisted",
        "provider-paid",
        "consumer-local",
        "customer-selected qualified bounded agent runner",
        "codex and claude code",
        "agent_execution_policy",
        "clean ephemeral session",
        "consumer_managed",
        "no adapter",
        "consumer installation",
        "deterministic verification",
        "model request disclosure",
        "generic-agent baseline",
        "same run",
        "organically agent-assisted",
        "consumer-selected qualified runner",
        "tested lumyn-opened draft pr",
        "provider-received consented status projection",
        "silence is unknown",
        "no auto-merge",
    ):
        require(token in locked, f"execution locked decisions missing {token}")
    validate_slices(
        plan.get("delivery_slices"),
        required_ids,
        "execution_plan.delivery_slices",
    )


def validate_task(
    task: dict[str, Any],
    required_ids: set[str],
    ledger_by_id: dict[str, dict[str, Any]],
) -> None:
    task_id = str(task.get("task_id", ""))
    missing = REQUIRED_TASK_FIELDS - set(task)
    require(not missing, f"{task_id} missing required fields {sorted(missing)}")
    require(
        task.get("milestone_ref") == MILESTONE_REFS[task_id],
        f"{task_id} milestone ref is stale",
    )
    require(
        task.get("dispatch_status")
        == "attended_task_requires_explicit_approval_factoryd_paused",
        f"{task_id} must require attended approval while factoryd is paused",
    )
    ids = task.get("acceptance_item_ids")
    require(list_of_strings(ids), f"{task_id} acceptance IDs must be non-empty")
    require(
        not (set(ids) - required_ids),
        f"{task_id} references unknown acceptance IDs",
    )
    checks = task.get("acceptance_checks")
    require(
        isinstance(checks, list) and len(checks) == len(ids),
        f"{task_id} acceptance checks must match its IDs",
    )
    expected_checks = {
        f"{item_id}: {ledger_by_id[item_id]['source_text']}" for item_id in ids
    }
    require(set(checks) == expected_checks, f"{task_id} acceptance checks drifted")
    require(
        task.get("acceptance_ledger_ref") == ARTIFACT_REFS["ledger"],
        f"{task_id} ledger ref is stale",
    )
    validate_runtime_pins(task.get("runtime_pins"), f"{task_id}")
    factory = task.get("factory_compatibility", {})
    require(
        factory.get("profile_ref") == "profiles/lumyn.yaml"
        and factory.get("status") == "blocked_profile_and_runtime_unqualified",
        f"{task_id} Factory compatibility must remain blocked",
    )
    required_capabilities = set(task.get("requires_capabilities", []))
    require(
        required_capabilities == EXPECTED_CAPABILITIES.get(task_id, set()),
        f"{task_id} Factory capabilities differ",
    )
    conditional = set(task.get("conditional_factory_capabilities", []))
    require(
        conditional
        == EXPECTED_CONDITIONAL_CAPABILITIES.get(task_id, set()),
        f"{task_id} conditional Factory capabilities differ",
    )
    require(
        (required_capabilities | conditional).issubset(FACTORY_CAPABILITIES),
        f"{task_id} uses a product action as a Factory capability",
    )
    product = set(task.get("product_authority_requirements", []))
    require(
        product == EXPECTED_PRODUCT_AUTHORITIES.get(task_id, set()),
        f"{task_id} product authority requirements differ",
    )
    optional_product = set(task.get("optional_product_action_capabilities", []))
    require(
        optional_product
        == EXPECTED_OPTIONAL_PRODUCT_AUTHORITIES.get(task_id, set()),
        f"{task_id} optional product action scopes differ",
    )
    route_contract = task.get("product_action_route_contract", {})
    expected_routes = EXPECTED_PRODUCT_ACTION_ROUTES.get(task_id, {})
    if product or optional_product:
        require(
            task.get("product_authority_scope_semantics")
            == PRODUCT_AUTHORITY_SCOPE_SEMANTICS,
            f"{task_id} product authority arrays must be capability universes",
        )
        require(
            bool(route_contract),
            f"{task_id} must define at least one exact product action route",
        )
    declared_route_capabilities: set[str] = set()
    for route_id, route in route_contract.items():
        required_route_capabilities = route.get("required_capabilities", [])
        conditional_route_capabilities = route.get(
            "conditionally_selected_capabilities", []
        )
        require(
            isinstance(required_route_capabilities, list)
            and all(
                nonempty_string(capability)
                for capability in required_route_capabilities
            )
            and len(required_route_capabilities)
            == len(set(required_route_capabilities)),
            f"{task_id}/{route_id} required route capabilities must be unique",
        )
        require(
            isinstance(conditional_route_capabilities, list)
            and all(
                nonempty_string(capability)
                for capability in conditional_route_capabilities
            )
            and len(conditional_route_capabilities)
            == len(set(conditional_route_capabilities)),
            f"{task_id}/{route_id} conditional route capabilities must be unique",
        )
        required_route_set = set(required_route_capabilities)
        conditional_route_set = set(conditional_route_capabilities)
        require(
            required_route_set.isdisjoint(conditional_route_set),
            f"{task_id}/{route_id} route capability classes must be disjoint",
        )
        predicates = route.get("conditional_capability_predicates", {})
        require(
            (
                conditional_route_set
                and isinstance(predicates, dict)
                and set(predicates) == conditional_route_set
                and all(
                    predicates[capability]
                    == CONDITIONAL_CAPABILITY_PREDICATES.get(capability)
                    and nonempty_string(predicates[capability])
                    for capability in conditional_route_set
                )
            )
            or (not conditional_route_set and predicates in ({}, None)),
            f"{task_id}/{route_id} conditional capability predicates differ",
        )
        agent_route_capabilities = {
            "agent_runner_credential",
            "agent_runner_network",
            "model_credential",
            "model_network",
            "model_request_disclosure",
        }
        topology = route.get("authorization_topology_contract")
        if conditional_route_set.intersection(agent_route_capabilities):
            require(
                topology == AGENT_ROUTE_AUTHORIZATION_TOPOLOGY_CONTRACT
                and all(
                    set(capabilities).issubset(conditional_route_set)
                    for capabilities in topology[
                        "minimum_capability_sets"
                    ].values()
                ),
                f"{task_id}/{route_id} agent authorization topology differs",
            )
        else:
            require(
                topology is None,
                f"{task_id}/{route_id} unexpected agent authorization topology",
            )
        route_union = required_route_set | conditional_route_set
        declared_route_capabilities.update(route_union)
        require(
            route_union.issubset(product | optional_product)
            and route.get("exact_selected_union_frozen_before_action") is True
            and route.get("unselected_capabilities_authorized") is False,
            f"{task_id}/{route_id} product action scope union is not exact",
        )
    require(
        declared_route_capabilities == product | optional_product,
        f"{task_id} routes must cover its exact product capability universe",
    )
    require(
        route_contract == expected_routes,
        f"{task_id} exact product action route contract differs",
    )
    require(
        task.get("requires_credentials") is ("credentials" in required_capabilities),
        f"{task_id} credential flag differs from required Factory capabilities",
    )
    require(
        task.get("requires_network") is ("network" in required_capabilities),
        f"{task_id} network flag differs from required Factory capabilities",
    )
    require(
        task.get("requires_human_approval")
        is bool(
            required_capabilities
            or conditional
            or product
            or optional_product
        ),
        f"{task_id} human approval flag differs from external/action scope",
    )
    require(
        list_of_strings(task.get("allowed_paths")),
        f"{task_id} allowed paths must be explicit",
    )
    require(
        list_of_strings(task.get("forbidden_paths")),
        f"{task_id} forbidden paths must be explicit",
    )
    require(
        all(not Path(path).is_absolute() and ".." not in Path(path).parts for path in task["allowed_paths"]),
        f"{task_id} allowed paths must be repo-relative",
    )
    require(
        f".factory/artifacts/pr-lifecycle/{task_id}/" in task["forbidden_paths"],
        f"{task_id} must forbid PR-lifecycle evidence writes",
    )
    require(
        "python3 scripts/validate_repo_pack.py" in task.get("validation_commands", []),
        f"{task_id} must run repo-pack validation",
    )
    require(
        not contains_machine_local_path(task),
        f"{task_id} contains a machine-local path",
    )


def validate_packets(
    packets: dict[str, Any],
    required_ids: set[str],
    ledger_by_id: dict[str, dict[str, Any]],
) -> dict[str, dict[str, Any]]:
    require(
        packets.get("artifact_type") == "task_packets",
        "task-packets artifact_type must be task_packets",
    )
    validate_runtime_pins(packets.get("runtime_pins"), "task packet set")
    for field, expected in (
        ("source_prd_ref", "docs/product/prd.md"),
        ("source_ref", ARTIFACT_REFS["plan"]),
        ("acceptance_ledger_ref", ARTIFACT_REFS["ledger"]),
        ("acceptance_mapping_ref", ARTIFACT_REFS["mapping"]),
    ):
        require(packets.get(field) == expected, f"task packet {field} is stale")
    require(
        packets.get("acceptance_item_count") == len(required_ids),
        "task packet acceptance count is stale",
    )
    tasks_value = packets.get("tasks")
    require(isinstance(tasks_value, list), "task-packets.tasks must be a list")
    tasks: dict[str, dict[str, Any]] = {}
    for task in tasks_value:
        require(isinstance(task, dict), "task entries must be objects")
        task_id = task.get("task_id")
        require(nonempty_string(task_id), "task packet missing task_id")
        require(task_id not in tasks, f"duplicate task packet {task_id}")
        tasks[str(task_id)] = task
    require(
        tuple(tasks) == EXPECTED_TASK_IDS,
        "task order/IDs differ from authored M0-M10",
    )
    for task in tasks.values():
        validate_task(task, required_ids, ledger_by_id)
    covered = set().union(
        *(set(task["acceptance_item_ids"]) for task in tasks.values())
    )
    require(covered == required_ids, "active tasks do not cover exact acceptance set")
    validate_migration_task_contracts(tasks)
    validate_slices(
        packets.get("delivery_slices"),
        required_ids,
        "task_packets.delivery_slices",
    )
    return tasks


def validate_contract(contract: dict[str, Any], required_ids: set[str]) -> None:
    require(
        contract.get("artifact_type") == "validation_contract",
        "validation contract artifact_type must be validation_contract",
    )
    require(
        contract.get("acceptance_ledger_ref") == ARTIFACT_REFS["ledger"]
        and contract.get("acceptance_mapping_ref") == ARTIFACT_REFS["mapping"],
        "validation contract acceptance refs are stale",
    )
    require(
        contract.get("acceptance_item_count") == len(required_ids),
        "validation contract acceptance count is stale",
    )
    validate_runtime_pins(contract.get("runtime_pins"), "validation contract")
    criteria = contract.get("acceptance_criteria")
    require(
        isinstance(criteria, list) and len(criteria) == len(required_ids),
        "validation contract must contain one criterion per acceptance item",
    )
    require(
        {str(value).split(":", 1)[0] for value in criteria} == required_ids,
        "validation criteria IDs differ from the PRD",
    )
    checks = contract.get("required_checks")
    for check in (
        "make lint-fast",
        "make test-fast",
        "make test-coverage",
        "make test-contracts",
        "make prepush-full",
        "GitHub Actions validate",
        "GitHub Actions CodeQL analyze",
        "passive Codex review settle",
    ):
        require(check in checks, f"validation contract missing {check}")
    authority = contract.get("authority_requirements", {})
    require(
        set(authority.get("factory_worker_capabilities", []))
        == FACTORY_CAPABILITIES,
        "validation contract Factory capabilities drifted",
    )
    require(
        set(authority.get("separate_product_action_scopes", []))
        == PRODUCT_AUTHORITY_CAPABILITIES,
        "validation contract product action scopes drifted",
    )
    require(
        "never substitute" in authority.get("separation_rule", "").lower(),
        "validation contract must separate Factory and product authority",
    )
    require(
        authority.get("task_and_campaign_arrays")
        == "capability_universe_only_not_runtime_grant"
        and authority.get(
            "named_product_action_route_required_when_scopes_declared"
        ) is True
        and authority.get(
            "exact_required_and_selected_conditional_union_frozen_before_action"
        ) is True
        and authority.get("route_scope_classes_disjoint_and_unique") is True
        and authority.get(
            "routes_cover_exact_task_capability_universe"
        ) is True
        and authority.get(
            "campaign_aggregate_union_authorized_per_installation"
        ) is False
        and authority.get(
            "remote_branch_and_draft_pr_are_atomic_actions"
        ) is True
        and authority.get("aggregate_composed_action_route_allowed") is False
        and authority.get(
            "delegated_route_semantics_must_match_source"
        ) is True,
        "validation contract must enforce exact per-action route authority",
    )
    model = contract.get("model_control_requirements", {})
    for field in (
        "agent_execution_policy",
        "agent_runner",
        "executable",
        "auth_and_entitlement",
        "onboarding_preflight",
        "route_identity",
        "m1_live_execution",
        "commercial_prequalification",
        "session",
        "runner_host_isolation",
        "funding",
        "managed_credential",
        "ownership",
        "native_configuration",
        "fallback",
        "verification_boundary",
        "verifier_isolation",
        "repository_command_isolation",
        "sandbox_entrypoint_isolation",
        "disclosure",
        "endpoint",
        "credential",
        "network",
        "context",
        "token_budget",
        "time_budget",
        "cost_budget",
        "attempt_budget",
        "path_budget",
        "diff_budget",
        "provenance",
    ):
        require(field in model, f"validation contract model controls missing {field}")
    conditional = contract.get("conditional_acceptance_rules", {})
    install_rule = str(conditional.get("INSTALL-001", "")).lower()
    require(
        all(
            token in install_rule
            for token in (
                "disabled",
                "configured",
                "only configured",
                "exact agent runner route",
            )
        ),
        "validation contract INSTALL-001 must preserve conditional Agent Runner setup",
    )
    pilot_rule = str(conditional.get("PILOT-003", "")).lower()
    require(
        all(
            token in pilot_rule
            for token in (
                "real consumer-selected",
                "agent-assisted",
                "independent exact-head",
                "installed-preauthorization",
                "same run",
                "organically",
                "provider-received",
                "deterministic rerouting",
            )
        ),
        "validation contract PILOT-003 must preserve both agent and composed-delivery proofs",
    )
    agent_rule = str(conditional.get("AGENT-007", "")).lower()
    require(
        all(
            token in agent_rule
            for token in (
                "managed-credential broker",
                "runner-host isolation",
                "blocks before launch",
            )
        ),
        "validation contract AGENT-007 must enforce broker and host isolation",
    )
    verifier_rule = str(conditional.get("VER-002", "")).lower()
    require(
        all(
            token in verifier_rule
            for token in (
                "fresh exact-head",
                "separate process",
                "frozen command/configuration",
                "no agent runner/model credentials",
                "evidence writer unavailable to generation",
            )
        ),
        "validation contract VER-002 must enforce verifier isolation",
    )
    evidence_rule = str(conditional.get("VER-006", "")).lower()
    require(
        all(
            token in evidence_rule
            for token in (
                "explicit agent_execution_policy",
                "every candidate and evidence outcome",
                "only when policy=configured",
            )
        ),
        "validation contract VER-006 must make runner provenance conditional",
    )
    verification = contract.get("verification_requirements", {})
    for field in (
        "independent_from_generation",
        "exact_candidate_head",
        "non_mutating_verify",
        "separate_repair_action",
        "model_self_verification_forbidden",
        "automated_draft_pr_required_for_exp_003",
        "manual_delivery_cannot_close_automated_delivery",
        "fresh_verification_view",
        "separate_verifier_process",
        "agent_runner_and_model_credentials_absent",
        "generation_session_cannot_write_verification_evidence",
        "verification_command_digest_frozen",
        "verification_configuration_digest_frozen",
        "repair_requires_configured_agent_execution_policy",
    ):
        require(
            verification.get(field) is True,
            f"validation contract verification requirement {field} must be true",
        )
    require(
        verification.get("agent_reported_tests_qualify") is False
        and verification.get("verification_evidence_writer")
        == "independent_verifier_evidence_boundary"
        and verification.get("repair_route") == "agent_assisted_only",
        "validation contract verifier or repair boundary drifted",
    )
    validate_pause_contract(
        contract.get("factoryd_runtime_requirements"),
        "validation contract factoryd runtime",
    )
    validate_slices(
        contract.get("delivery_slices"),
        required_ids,
        "validation_contract.delivery_slices",
    )


def validate_closure(
    closure: dict[str, Any],
    ledger_by_id: dict[str, dict[str, Any]],
    required_ids: set[str],
) -> None:
    require(
        closure.get("artifact_type") == "scope_closure_map",
        "closure artifact_type must be scope_closure_map",
    )
    require(
        closure.get("acceptance_ledger_ref") == ARTIFACT_REFS["ledger"]
        and closure.get("acceptance_mapping_ref") == ARTIFACT_REFS["mapping"],
        "closure acceptance refs are stale",
    )
    require(
        closure.get("acceptance_item_count") == len(required_ids),
        "closure acceptance count is stale",
    )
    items = closure.get("items")
    require(isinstance(items, list), "closure.items must be a list")
    by_id = {
        str(item.get("scope_item_id")): item
        for item in items
        if isinstance(item, dict)
    }
    require(
        len(by_id) == len(items) and set(by_id) == required_ids,
        "closure map must cover exact PRD IDs once",
    )
    for item_id, item in by_id.items():
        ledger_item = ledger_by_id[item_id]
        expected = "implemented" if ledger_item["status"] == "implemented" else "missing"
        require(
            item.get("status") == expected,
            f"{item_id} closure status differs from ledger",
        )
        require(
            item.get("acceptance_item_ids") == [item_id],
            f"{item_id} closure must be item-granular",
        )
        if expected == "implemented":
            require(
                item.get("implemented_task_refs") == ["M0"],
                f"{item_id} implemented closure must bind to M0",
            )
            require(
                item.get("remaining_task_refs") == [],
                f"{item_id} implemented closure cannot retain M0 as remaining scope",
            )
            require(
                item.get("evidence_refs") == ledger_item.get("evidence_refs"),
                f"{item_id} closure evidence differs from ledger",
            )
        else:
            require(
                item.get("implemented_task_refs") == [],
                f"{item_id} cannot claim an unexecuted v3 task",
            )
            require(
                list_of_strings(item.get("remaining_task_refs"))
                and item.get("remaining_task_refs") == item.get("task_refs"),
                f"{item_id} must retain every active v3 task as remaining scope",
            )
    slices = validate_slices(
        closure.get("delivery_slices"),
        required_ids,
        "scope_closure_map.delivery_slices",
    )
    for item in slices.values():
        expected_remaining = set(item["acceptance_item_ids"]) - EXPECTED_IMPLEMENTED_ITEMS
        require(
            set(item.get("remaining_acceptance_item_ids", []))
            == expected_remaining,
            f"{item['slice_id']} remaining acceptance IDs are stale",
        )
    require(closure.get("closure_status") == "partial", "v3 closure must be partial")


def validate_config_payload(config: dict[str, Any], label: str) -> None:
    require(
        config.get("factory", {}).get("profile_path") == "profiles/lumyn.yaml",
        f"{label} profile ref is stale",
    )
    repo = config.get("repos", {}).get("lumyn")
    require(isinstance(repo, dict), f"{label} missing repos.lumyn")
    for field, ref_key in (
        ("task_packets", "packets"),
        ("scope_closure_map", "closure"),
        ("validation_contract", "contract"),
        ("acceptance_ledger", "ledger"),
    ):
        require(
            repo.get(field) == ARTIFACT_REFS[ref_key],
            f"{label} {field} ref is stale",
        )
    validate_active_repo_safety(repo, ARTIFACT_REFS)
    posture = " ".join(
        str(repo.get(field, ""))
        for field in (
            "approval_posture",
            "credential_posture",
            "network_posture",
            "dispatch_posture",
        )
    ).lower()
    for token in (
        "factory worker capabilities",
        "approval, credentials, and network",
        "product action",
        "no ambient",
        "offline by default",
        "v3.1 profile",
        "factoryd qualification",
    ):
        require(token in posture, f"{label} posture missing {token}")
    require(
        "python3 scripts/validate_repo_pack.py" in repo.get(
            "validation_commands", []
        ),
        f"{label} must run repo-pack validation",
    )
    validate_architecture_budget_policy(ROOT, repo, label)


def validate_loaded(
    data: dict[str, dict[str, Any]],
    *,
    validate_configs: bool = True,
) -> dict[str, dict[str, Any]]:
    required_ids = expected_acceptance_ids()
    validate_acceptance_text(PRD.read_text(), data["ledger"])
    for name, payload in data.items():
        require(
            not contains_machine_local_path(payload),
            f"{name} contains a machine-local path",
        )
        require(
            not contains_true_key(
                payload,
                {
                    "auto_merge",
                    "default_branch_write",
                    "provider_raw_repo_access",
                    "provider_code_access",
                },
            ),
            f"{name} enables a forbidden product capability",
        )
        validate_markdown_fragment_refs(ROOT, payload, name)
    validate_context(data["context"])
    validate_risk(data["risk"])
    ledger_by_id = validate_ledger(data["ledger"], required_ids)
    validate_mapping(data["mapping"], required_ids)
    validate_plan(data["plan"], required_ids)
    tasks = validate_packets(data["packets"], required_ids, ledger_by_id)
    validate_contract(data["contract"], required_ids)
    validate_closure(data["closure"], ledger_by_id, required_ids)
    plan_slices = {
        item["slice_id"]: set(item["acceptance_item_ids"])
        for item in data["plan"]["delivery_slices"]
    }
    for name in ("packets", "contract", "mapping", "closure"):
        other = {
            item["slice_id"]: set(item["acceptance_item_ids"])
            for item in data[name]["delivery_slices"]
        }
        require(other == plan_slices, f"{name} delivery slices differ from plan")
    if validate_configs:
        validate_config_payload(data["config"], ".factory/factoryd.example.json")
        validate_config_payload(
            data["autoship"], ".factory/factoryd.autoship.example.json"
        )
    return tasks


def validate_active_config(
    config: dict[str, Any],
    tasks: dict[str, dict[str, Any]],
) -> None:
    repo = config.get("repos", {}).get("lumyn")
    validate_active_repo_safety(repo, ARTIFACT_REFS)
    validate_authority_grants(
        repo.get("capability_grants"),
        tasks,
        EXPECTED_CAPABILITIES,
    )


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
            validate_active_config(load_json(ACTIVE_CONFIG), tasks)
    except AssertionError as exc:
        print(f"repo-pack validation failed: {exc}", file=sys.stderr)
        return 2
    print("repo-pack validation passed")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
