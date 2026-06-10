#!/usr/bin/env python3
from __future__ import annotations

import json
import re
import sys
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parents[1]
PLAN_DIR = ROOT / ".factory" / "artifacts" / "prd-to-plan" / "lumyn-mvp"
CONTEXT_BRIEF = PLAN_DIR / "context-brief.json"
EXECUTION_PLAN = PLAN_DIR / "execution-plan.json"
TASK_PACKETS = PLAN_DIR / "task-packets.json"
VALIDATION_CONTRACT = PLAN_DIR / "validation-contract.json"
ACCEPTANCE_LEDGER = PLAN_DIR / "acceptance-ledger.json"
ACCEPTANCE_MAPPING = PLAN_DIR / "acceptance-mapping.json"
SCOPE_CLOSURE_MAP = PLAN_DIR / "scope-closure-map.json"
RISK_CLASSIFICATION = PLAN_DIR / "risk-classification.json"
FACTORYD_CONFIG = ROOT / ".factory" / "factoryd.example.json"
FACTORYD_ACTIVE_CONFIG = ROOT / ".factory" / "factoryd.json"
FACTORYD_AUTOSHIP_CONFIG = ROOT / ".factory" / "factoryd.autoship.example.json"
REQUIRED_CHECKS = ROOT / ".github" / "required-checks.json"
CODEOWNERS = ROOT / ".github" / "CODEOWNERS"
ACTION_REF_EXCEPTIONS = ROOT / ".github" / "action-ref-exceptions.yaml"
VALIDATE_WORKFLOW = ROOT / ".github" / "workflows" / "validate.yml"
CODEQL_WORKFLOW = ROOT / ".github" / "workflows" / "codeql.yml"
REPAIR_TASK_PACKETS = [
    ROOT / ".factory" / "artifacts" / "pilot" / "lumyn-mvp-slice" / "repair-loop" / "task-packet.json"
]
PRD = ROOT / "docs" / "product" / "prd.md"
TEST_MATRIX_SOURCE_BASE = "docs/dev/dev_guides.md"
COVERAGE_POLICY_SOURCE_BASE = "docs/dev/dev_guides.md"
ARCHITECTURE_GUIDE_BASE = "docs/architecture/architecture_guides.md"
RUNTIME_CONTROL_ALLOWED_RE = re.compile(
    r"^\.factory/artifacts/?$"
    r"|^\.factory/artifacts/(prd-to-plan|post-prd)(/.*)?$"
)
RUNTIME_CONTROL_FORBIDDEN_PATHS = [
    ".factory/artifacts/prd-to-plan/lumyn-mvp/",
    ".factory/artifacts/prd-to-plan/lumyn-mvp/context-brief.json",
    ".factory/artifacts/prd-to-plan/lumyn-mvp/execution-plan.json",
    ".factory/artifacts/prd-to-plan/lumyn-mvp/task-packets.json",
    ".factory/artifacts/prd-to-plan/lumyn-mvp/validation-contract.json",
    ".factory/artifacts/prd-to-plan/lumyn-mvp/acceptance-ledger.json",
    ".factory/artifacts/prd-to-plan/lumyn-mvp/acceptance-mapping.json",
    ".factory/artifacts/prd-to-plan/lumyn-mvp/scope-closure-map.json",
]

REQUIRED_GUIDES = [
    "docs/dev/dev_guides.md",
    "docs/architecture/architecture_guides.md",
]

REQUIRED_TASK_FIELDS = [
    "ci_lane_refs",
    "test_matrix_refs",
    "coverage_policy_refs",
    "security_scanner_gates",
    "engineering_policy_refs",
    "architecture_guidance_refs",
    "factory_compatibility",
    "scope_exclusions",
    "alignment_gate_ref",
    "plan_drift_policy_ref",
    "required_worker_chain",
    "lifecycle_gates",
    "allowed_paths",
    "forbidden_paths",
    "worker_type",
    "factoryd_runtime",
    "validation_commands",
    "evidence_required",
    "stop_conditions",
    "acceptance_ledger_ref",
    "acceptance_item_ids",
]


def normalize_repo_path(value: object) -> str:
    path = str(value).strip().replace("\\", "/")
    parts: list[str] = []
    for part in path.split("/"):
        if part in {"", "."}:
            continue
        if part == "..":
            if parts:
                parts.pop()
            else:
                parts.append(part)
            continue
        parts.append(part)
    return "/".join(parts)

REQUIRED_RUNNER_READY_FIELDS = [
    "worker_type",
    "factoryd_runtime",
    "validation_commands",
    "evidence_required",
    "stop_conditions",
    "allowed_paths",
    "forbidden_paths",
    "required_worker_chain",
    "lifecycle_gates",
    "scope_exclusions",
    "acceptance_ledger_ref",
    "acceptance_item_ids",
]

REQUIRED_ACCEPTANCE_ITEM_IDS = {
    "FDN-001",
    "FDN-002",
    "FDN-003",
    "FDL-001",
    "FDAP-001",
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
        "task_refs": {"T1", "T2", "T2.7", "T3", "T4.1", "T4.2", "T4.3", "T5", "T6", "T10"},
    },
    "v0.1": {
        "capability_group_id": "live_verify_boundary_ci_share",
        "task_refs": {"T7", "T8", "T9", "T10"},
    },
    "v0.2": {
        "capability_group_id": "live_agent_eval",
        "task_refs": {"T11", "T12"},
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
    "RCRR-012": {"T6"},
    "LVCIS-009": {"T8"},
    "LVCIS-010": {"T9"},
    "EVAL-011": {"T11", "T12"},
    "FR14": {"T3", "T4.1", "T4.2", "T4.3", "T5", "T6", "T7", "T8", "T9", "T10", "T11", "T12"},
    "NFR9": {"T3", "T4.1", "T4.2", "T4.3", "T5", "T6", "T7", "T8", "T9", "T10", "T11", "T12"},
    "NFR12": {"T4.1", "T4.2", "T4.3", "T5", "T6", "T7", "T8", "T9", "T10", "T11", "T12"},
    "FR9": {"T4.2", "T4.3", "T5", "T6", "T7", "T8", "T9", "T10", "T11", "T12"},
    "FR2": {"T3", "T4.1", "T4.2", "T4.3", "T5", "T6", "T10", "T11", "T12"},
    "NFR6": {"T4.1", "T4.2", "T4.3", "T5", "T6", "T7", "T10", "T11", "T12"},
}

REQUIRED_FACTORYD_RUNTIME_FIELDS = [
    "state_dir",
    "workspace_root",
    "branch_prefix",
    "worker_type",
    "worker_command",
    "approval_posture",
    "credential_posture",
    "network_posture",
    "capability_grants",
]

REQUIRED_PLAN_SKILL_REFS = [
    "factory://skills/prd-to-plan",
    "factory://skills/execution-compiler",
]

DEPRECATED_ACTIVE_WORKERS = {
    "ship-pr": "commit-push",
}

DEFAULT_REQUIRED_WORKER_CHAIN = [
    "task-executor",
    "validation-gate",
    "commit-push",
    "post-merge-monitor",
]

REVIEW_REQUIRED_WORKER_CHAIN = [
    "task-executor",
    "validation-gate",
    "code-review",
    "commit-push",
    "post-merge-monitor",
]

REQUIRED_PLAN_LEVEL_FIELDS = [
    "planning_skill_alignment",
    "factory_compatibility",
    "runtime_pins",
    "alignment_gate",
    "plan_drift_policy",
    "acceptance_ledger_coverage",
    "mvp_required_version_slices",
    "public_api_and_contract_map",
    "docs_and_oss_readiness_baseline",
    "test_matrix_wiring",
    "minimum_now_sequence",
    "explicit_non_goals",
    "definition_of_done",
]

REQUIRED_TASK_PLANNING_FIELDS = [
    "planning_skill_refs",
    "runtime_pins",
    "slice_rationale",
    "changelog",
    "contract_impact",
    "versioning_migration_impact",
    "architecture_constraints",
    "adr_required",
    "tdd_first_failing_tests",
    "cost_perf_impact",
    "chaos_failure_hypothesis",
    "semantic_invariants",
]

REQUIRED_RUNTIME_PIN_FIELDS = [
    "language",
    "go_version",
    "toolchain_version",
    "module_path",
    "module_or_package_path",
    "dependency_policy",
    "distribution_target",
    "provider_policy",
    "artifact_namespace",
    "live_work_policy",
]

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

REQUIRED_FACTORY_COMPATIBILITY_FIELDS = [
    "factory_contract_version",
    "profile_ref",
    "skill_vocabulary_version",
    "skill_inventory_ref",
    "generated_by",
    "generated_at",
    "deprecated_worker_policy",
    "deprecated_worker_aliases",
]

REQUIRED_PLAN_DRIFT_UPDATES = [
    "context_brief",
    "execution_plan",
    "task_packets",
    "validation_contract",
    "factory_compatibility",
    "acceptance_ledger",
    "acceptance_mapping",
    "scope_closure_map",
]

REQUIRED_CHANGELOG_FIELDS = [
    "impact",
    "section",
    "draft_entry",
    "semver_marker_override",
]

ADR_CONTRACT_TOKENS = [
    "public",
    "api",
    "cli",
    "command",
    "schema",
    "artifact",
    "output",
    "json",
    "contract",
    "ci",
    "review",
    "redaction",
    "share",
    "eval",
    "proof",
]

MACHINE_LOCAL_PATH_RE = re.compile(r"(?<![A-Za-z0-9+./#-])(?:/(?!/)[A-Za-z0-9._-][^\s\"'<>]*|[A-Za-z]:[\\/])")
NO_CONTRACT_IMPACT_BREAKERS = [
    " but ",
    " except ",
    " however ",
    " although ",
    " changes ",
    " adds ",
    " removes ",
    " modifies ",
]


def declares_no_contract_impact(value: str) -> bool:
    collapsed = re.sub(r"\s+", " ", value.lower()).strip(" .")
    normalized = f" {collapsed} "
    if normalized.strip() in {"none", "n/a", "not applicable"}:
        return True
    if not normalized.startswith(" no "):
        return False
    if any(breaker in normalized for breaker in NO_CONTRACT_IMPACT_BREAKERS):
        return False
    return any(token in normalized for token in [" impact ", " change ", " changes ", " effect ", " effects "])


def has_adr_contract_token(value: str) -> bool:
    if declares_no_contract_impact(value):
        return False
    normalized = value.lower()
    return any(re.search(rf"\b{re.escape(token)}\b", normalized) for token in ADR_CONTRACT_TOKENS)


REQUIRED_CI_LANES = [
    "fast",
    "core",
    "acceptance",
    "cross_platform",
    "risk",
    "release",
]

REQUIRED_ENGINEERING_POLICIES = [
    "docs_parity",
    "output_contracts",
    "release_integrity",
    "provenance_evidence",
]

REQUIRED_ARCHITECTURE_POLICIES = [
    "systems_thinking",
    "tdd",
    "adr_triggers",
    "performance",
    "reliability",
    "failure_semantics",
]

REQUIRED_STATUS_CHECKS = [
    "validate",
    "CodeQL analyze",
]

REQUIRED_ACTION_REFS = [
    "actions/checkout@v6.0.2",
    "actions/setup-go@v6.3.0",
    "github/codeql-action/init@v4",
    "github/codeql-action/autobuild@v4",
    "github/codeql-action/analyze@v4",
]

ARCHITECTURE_POLICY_TOKENS = {
    "systems_thinking": ["systems-thinking", "systems thinking"],
    "tdd": ["tdd", "red-first"],
    "adr_triggers": ["adr", "decision"],
    "performance": ["performance", "cost"],
    "reliability": ["reliability", "recovery"],
    "failure_semantics": ["fail-closed", "failure", "trust-mode"],
}

STOP_CONDITION_CATEGORIES = {
    "test_matrix": ["test-matrix", "test matrix", "test_matrix"],
    "ci_lanes": ["ci lane", "ci/status", "status check"],
    "scanner": ["scanner", "security"],
    "coverage_policy": ["coverage", "test-coverage"],
    "engineering_policies": ["docs parity", "output contract", "release integrity", "provenance"],
    "architecture_policies": ["architecture", "systems-thinking", "systems thinking", "adr", "fail-closed"],
    "planning_skill": ["prd-to-plan", "execution-compiler", "planning-skill", "planning skill"],
    "contract_discipline": ["changelog", "contract/api", "semantic invariants", "semantic_invariants"],
}

TASK_ORDER_RE = re.compile(r"^T(?P<version>\d+(?:\.\d+)*)(?:[^.\d].*)?$", re.IGNORECASE)


def fail(message: str) -> None:
    raise AssertionError(message)


def load_json(path: Path) -> dict[str, Any]:
    if not path.exists():
        fail(f"missing JSON artifact: {path.relative_to(ROOT)}")
    try:
        payload = json.loads(path.read_text())
    except Exception as exc:
        fail(f"{path.relative_to(ROOT)} is not valid JSON: {exc}")
    if not isinstance(payload, dict):
        fail(f"{path.relative_to(ROOT)} must contain a JSON object")
    return payload


def require_existing(relative_path: str) -> None:
    if not (ROOT / relative_path).exists():
        fail(f"missing required repo-pack file: {relative_path}")


def ref_file_exists(ref: Any) -> bool:
    if not isinstance(ref, str) or not ref.strip():
        return False
    path_part = ref.split("#", 1)[0]
    return bool(path_part) and (ROOT / path_part).exists()


def validate_coverage_policy_refs(value: Any, label: str) -> None:
    if not isinstance(value, dict):
        fail(f"{label} must be an object")
    for key in ["exception_ref", "evidence_ref"]:
        if isinstance(value.get(key), str) and value[key].strip() and not ref_file_exists(value[key]):
            fail(f"{label}.{key} points to missing file {value[key]}")
    minimums = value.get("minimums")
    if isinstance(minimums, list):
        for index, item in enumerate(minimums):
            if not isinstance(item, dict):
                continue
            if isinstance(item.get("exception_ref"), str) and item["exception_ref"].strip() and not ref_file_exists(item["exception_ref"]):
                fail(f"{label}.minimums[{index}].exception_ref points to missing file {item['exception_ref']}")


def validate_guides() -> None:
    for relative_path in REQUIRED_GUIDES:
        require_existing(relative_path)
    dev_guide = (ROOT / "docs/dev/dev_guides.md").read_text()
    dev_guide_lower = dev_guide.lower()
    tiers = set(re.findall(r"\|\s*Tier\s+(\d+)\b", dev_guide))
    expected = {str(index) for index in range(1, 13)}
    if tiers != expected:
        fail(f"docs/dev/dev_guides.md must preserve all 12 test tiers; found {sorted(tiers)}")
    for token in ["coverage gates", "make test-coverage", ">= 75%"]:
        if token not in dev_guide_lower:
            fail(f"docs/dev/dev_guides.md missing coverage token {token!r}")
    makefile = (ROOT / "Makefile").read_text()
    for token in ["test-coverage:", "check_go_coverage.py", "prepush-full: fmt lint-fast test-fast test-coverage"]:
        if token not in makefile:
            fail(f"Makefile missing coverage gate token {token!r}")
    arch_guide = (ROOT / "docs/architecture/architecture_guides.md").read_text().lower()
    for token in ["systems thinking", "tdd", "adr", "performance", "reliability", "fail-closed", "coverage gates"]:
        if token not in arch_guide:
            fail(f"docs/architecture/architecture_guides.md missing architecture token {token!r}")


def validate_ci_control_set() -> None:
    for path in [REQUIRED_CHECKS, CODEOWNERS, ACTION_REF_EXCEPTIONS, VALIDATE_WORKFLOW, CODEQL_WORKFLOW]:
        if not path.exists():
            fail(f"missing CI control file: {path.relative_to(ROOT)}")

    required_checks = load_json(REQUIRED_CHECKS).get("required_checks")
    if not isinstance(required_checks, list):
        fail(".github/required-checks.json.required_checks must be a list")
    missing_checks = [check for check in REQUIRED_STATUS_CHECKS if check not in required_checks]
    if missing_checks:
        fail(f".github/required-checks.json missing required checks: {missing_checks}")

    validate_workflow = VALIDATE_WORKFLOW.read_text()
    validate_tokens = [
        "pull_request:",
        "push:",
        "branches:",
        "- main",
        "permissions:",
        "contents: read",
        "concurrency:",
        "cancel-in-progress: true",
        "timeout-minutes:",
        "actions/checkout@v6.0.2",
        "actions/setup-go@v6.3.0",
        "go-version-file: go.mod",
        "check-latest: false",
        "cache: true",
        "make prepush-full",
    ]
    for token in validate_tokens:
        if token not in validate_workflow:
            fail(f".github/workflows/validate.yml missing CI control token {token!r}")

    codeql_workflow = CODEQL_WORKFLOW.read_text()
    codeql_tokens = [
        "pull_request:",
        "push:",
        "branches:",
        "- main",
        "permissions:",
        "security-events: write",
        "contents: read",
        "concurrency:",
        "cancel-in-progress: true",
        "timeout-minutes:",
        "actions/checkout@v6.0.2",
        "actions/setup-go@v6.3.0",
        "go-version-file: go.mod",
        "check-latest: false",
        "github/codeql-action/init@v4",
        "github/codeql-action/autobuild@v4",
        "github/codeql-action/analyze@v4",
        "languages: go",
    ]
    for token in codeql_tokens:
        if token not in codeql_workflow:
            fail(f".github/workflows/codeql.yml missing CI control token {token!r}")

    codeowners = CODEOWNERS.read_text()
    for token in ["*", "/.github/**", "/schemas/**", "/cmd/**", "/internal/**"]:
        if token not in codeowners:
            fail(f".github/CODEOWNERS missing owner token {token!r}")

    action_exceptions = ACTION_REF_EXCEPTIONS.read_text()
    for token in REQUIRED_ACTION_REFS + ["owner:", "reason:", "scope:", "expires:", "review_command:"]:
        if token not in action_exceptions:
            fail(f".github/action-ref-exceptions.yaml missing action-ref token {token!r}")


def task_id(task: dict[str, Any]) -> str:
    value = task.get("task_id")
    return value if isinstance(value, str) else ""


def depends_on(task_id_value: str, baseline_task_id: str, tasks_by_id: dict[str, dict[str, Any]], seen: set[str] | None = None) -> bool:
    if task_id_value == baseline_task_id:
        return True
    seen = seen or set()
    if task_id_value in seen:
        return False
    seen.add(task_id_value)
    task = tasks_by_id.get(task_id_value)
    if not task:
        return False
    blocked_by = task.get("blocked_by", [])
    if not isinstance(blocked_by, list):
        fail(f"{task_id_value}.blocked_by must be a list")
    for dependency in [str(value) for value in blocked_by]:
        if dependency == baseline_task_id or depends_on(dependency, baseline_task_id, tasks_by_id, seen):
            return True
    return False


def task_order_key(value: Any) -> tuple[int, ...] | None:
    if not isinstance(value, str):
        return None
    match = TASK_ORDER_RE.match(value.strip())
    if not match:
        return None
    return tuple(int(part) for part in match.group("version").split("."))


def version_gte(candidate: tuple[int, ...], baseline: tuple[int, ...]) -> bool:
    width = max(len(candidate), len(baseline))
    return candidate + (0,) * (width - len(candidate)) >= baseline + (0,) * (width - len(baseline))


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


def source_ref_base(value: Any) -> str:
    return value.split("#", 1)[0] if isinstance(value, str) else ""


def refs_include_base(task: dict[str, Any], field: str, expected_base: str) -> bool:
    value = task.get(field)
    if not isinstance(value, list):
        return False
    return any(isinstance(item, dict) and source_ref_base(item.get("source_ref")) == expected_base for item in value)


def object_source_ref_base(value: Any) -> str:
    if not isinstance(value, dict):
        return ""
    return source_ref_base(value.get("source_ref"))


def has_nonempty_list(value: Any) -> bool:
    return isinstance(value, list) and any(isinstance(item, str) and item.strip() for item in value)


def has_nonempty_string(value: Any) -> bool:
    return isinstance(value, str) and bool(value.strip())


def has_nonempty_dict(value: Any) -> bool:
    return isinstance(value, dict) and bool(value)


def is_valid_factoryd_runtime(value: Any) -> bool:
    if not isinstance(value, dict):
        return False
    for field in REQUIRED_FACTORYD_RUNTIME_FIELDS:
        if field == "capability_grants":
            if not isinstance(value.get(field), list):
                return False
            continue
        if field == "worker_command":
            if field not in value:
                return False
            continue
        if not has_nonempty_string(value.get(field)):
            return False
    if value.get("worker_type") != "codex_cli":
        return False
    credential_posture = str(value.get("credential_posture", "")).lower()
    if "no ambient" not in credential_posture:
        return False
    network_posture = str(value.get("network_posture", "")).lower()
    return "offline" in network_posture or "allowlist" in network_posture


def validate_factoryd_runtime(value: Any, label: str) -> None:
    if not isinstance(value, dict):
        fail(f"{label} must be an object")
    missing = [
        field
        for field in REQUIRED_FACTORYD_RUNTIME_FIELDS
        if field not in ["worker_command", "capability_grants"] and not has_nonempty_string(value.get(field))
    ]
    if "worker_command" not in value:
        missing.append("worker_command")
    if not isinstance(value.get("capability_grants"), list):
        missing.append("capability_grants")
    if missing:
        fail(f"{label} missing fields: {', '.join(missing)}")
    if value.get("worker_type") != "codex_cli":
        fail(f"{label}.worker_type must be codex_cli")
    credential_posture = str(value.get("credential_posture", "")).lower()
    if "no ambient" not in credential_posture:
        fail(f"{label}.credential_posture must declare no ambient secrets")
    network_posture = str(value.get("network_posture", "")).lower()
    if "offline" not in network_posture and "allowlist" not in network_posture:
        fail(f"{label}.network_posture must be offline or allowlisted")


def has_nonempty_collection(value: Any) -> bool:
    return (isinstance(value, dict) and bool(value)) or (isinstance(value, list) and bool(value))


def has_required_string_refs(value: Any, expected_refs: list[str]) -> bool:
    if not isinstance(value, list):
        return False
    present = {item for item in value if isinstance(item, str) and item.strip()}
    return all(expected in present for expected in expected_refs)


def contains_machine_local_path(value: Any) -> bool:
    if isinstance(value, str):
        return bool(MACHINE_LOCAL_PATH_RE.search(value))
    if isinstance(value, list):
        return any(contains_machine_local_path(item) for item in value)
    if isinstance(value, dict):
        return any(
            contains_machine_local_path(key) or contains_machine_local_path(item)
            for key, item in value.items()
        )
    return False


def task_slice_type(task: dict[str, Any]) -> str:
    rationale = task.get("slice_rationale")
    nested = str(rationale["slice_type"]) if isinstance(rationale, dict) and has_nonempty_string(rationale.get("slice_type")) else ""
    top_level = task.get("slice_type")
    top_level = top_level if isinstance(top_level, str) else ""
    if nested and top_level and nested != top_level:
        fail(f"{task_id(task)} has conflicting slice_type declarations")
    return nested or top_level


def iter_required_worker_chains(value: Any, path: str = "$") -> list[tuple[str, list[Any]]]:
    chains: list[tuple[str, list[Any]]] = []
    if isinstance(value, dict):
        for key, item in value.items():
            child_path = f"{path}.{key}"
            if key == "required_worker_chain" and isinstance(item, list):
                chains.append((child_path, item))
            chains.extend(iter_required_worker_chains(item, child_path))
    elif isinstance(value, list):
        for index, item in enumerate(value):
            chains.extend(iter_required_worker_chains(item, f"{path}[{index}]"))
    return chains


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


def validate_no_deprecated_active_workers(value: dict[str, Any], label: str) -> None:
    for path, chain in iter_required_worker_chains(value):
        for index, worker in enumerate(chain):
            if isinstance(worker, str) and worker in DEPRECATED_ACTIVE_WORKERS:
                replacement = DEPRECATED_ACTIVE_WORKERS[worker]
                fail(
                    f"{label}{path[1:]}[{index}] uses deprecated worker {worker!r}; "
                    f"use {replacement!r} in active required_worker_chain values"
                )


def has_factory_compatibility(value: Any) -> bool:
    if not isinstance(value, dict):
        return False
    if not all(field in value for field in REQUIRED_FACTORY_COMPATIBILITY_FIELDS):
        return False
    if value.get("deprecated_worker_policy") != "block_active_aliases":
        return False
    aliases = value.get("deprecated_worker_aliases")
    if not isinstance(aliases, list) or not aliases:
        return False
    return any(
        isinstance(alias, dict)
        and alias.get("deprecated") == "ship-pr"
        and alias.get("replacement") == "commit-push"
        for alias in aliases
    )


def validate_factory_compatibility(value: Any, label: str) -> None:
    if not has_factory_compatibility(value):
        fail(f"{label} must include current Factory compatibility metadata and ship-pr -> commit-push alias policy")


def has_runtime_pins(value: Any) -> bool:
    return (
        isinstance(value, dict)
        and all(has_nonempty_string(value.get(field)) for field in REQUIRED_RUNTIME_PIN_FIELDS)
        and has_required_string_refs(value.get("mvp_eval_providers"), REQUIRED_MVP_EVAL_PROVIDERS)
    )


def validate_runtime_pins(value: Any, label: str) -> None:
    if not isinstance(value, dict):
        fail(f"{label} must be an object")
    missing = [field for field in REQUIRED_RUNTIME_PIN_FIELDS if not has_nonempty_string(value.get(field))]
    if missing:
        fail(f"{label} missing runtime pin fields: {', '.join(missing)}")
    if not has_required_string_refs(value.get("mvp_eval_providers"), REQUIRED_MVP_EVAL_PROVIDERS):
        fail(
            f"{label}.mvp_eval_providers must include "
            f"{', '.join(REQUIRED_MVP_EVAL_PROVIDERS)}"
        )


def validate_mvp_eval_provider_adapters(value: Any, label: str) -> None:
    if not isinstance(value, dict):
        fail(f"{label} must be an object")
    if not has_required_string_refs(value.get("adapters"), REQUIRED_MVP_EVAL_ADAPTERS):
        fail(f"{label}.adapters must include {', '.join(REQUIRED_MVP_EVAL_ADAPTERS)}")
    if not has_required_string_refs(value.get("adapter_ids"), REQUIRED_MVP_EVAL_PROVIDERS):
        fail(f"{label}.adapter_ids must include {', '.join(REQUIRED_MVP_EVAL_PROVIDERS)}")
    if not has_required_string_refs(value.get("config_fields"), ["provider", "model", "temperature", "base_url", "api_key_env"]):
        fail(f"{label}.config_fields must include provider/model/temperature/base_url/api_key_env")


def has_alignment_gate(value: Any) -> bool:
    return (
        isinstance(value, dict)
        and value.get("status") == "resolved"
        and has_nonempty_string(value.get("source_context_brief_ref"))
        and has_nonempty_list(value.get("blocking_decisions"))
        and value.get("implementation_may_start") is True
    )


def validate_alignment_gate(value: Any, label: str) -> None:
    if not has_alignment_gate(value):
        fail(f"{label} must be resolved, cite the context brief, list blocking decisions, and allow implementation")


def has_plan_drift_policy(value: Any) -> bool:
    if not isinstance(value, dict):
        return False
    updates = value.get("required_updates")
    return (
        has_nonempty_list(value.get("drift_triggers"))
        and isinstance(updates, list)
        and all(required in updates for required in REQUIRED_PLAN_DRIFT_UPDATES)
        and value.get("continuation_behavior") == "block_until_artifacts_updated"
    )


def validate_plan_drift_policy(value: Any, label: str) -> None:
    if not has_plan_drift_policy(value):
        fail(
            f"{label} must require context brief, execution plan, task packets, validation contract, "
            "factory_compatibility, acceptance_ledger, acceptance_mapping, and scope_closure_map updates before continuing"
        )


def validate_acceptance_ledger_coverage(value: Any, label: str) -> None:
    if not isinstance(value, dict):
        fail(f"{label} must be an object")
    expected_refs = {
        "ledger_ref": ACCEPTANCE_LEDGER_REF,
        "acceptance_mapping_ref": ACCEPTANCE_MAPPING_REF,
        "scope_closure_map_ref": SCOPE_CLOSURE_MAP_REF,
    }
    for field, expected in expected_refs.items():
        if value.get(field) != expected:
            fail(f"{label}.{field} must cite {expected}")
    if value.get("coverage_unit") != "acceptance_item":
        fail(f"{label}.coverage_unit must be acceptance_item")
    if value.get("group_only_refs_allowed") is not False:
        fail(f"{label}.group_only_refs_allowed must be false")
    if value.get("required_item_count") != len(REQUIRED_ACCEPTANCE_ITEM_IDS):
        fail(f"{label}.required_item_count must match acceptance-ledger item count")
    if value.get("status") != "mapped":
        fail(f"{label}.status must be mapped")
    if not has_nonempty_list(value.get("required_groups")):
        fail(f"{label}.required_groups must be non-empty")
    required_slices = value.get("required_version_slices")
    if not isinstance(required_slices, list) or set(str(item) for item in required_slices) != set(REQUIRED_MVP_VERSION_SLICES):
        fail(f"{label}.required_version_slices must list v0.0, v0.1, and v0.2")


def has_lifecycle_gates(value: Any) -> bool:
    required_true = [
        "local_validation_required",
        "ci_required",
        "codex_review_required",
        "ship_pr_required",
        "post_merge_monitor_required",
        "pr_lifecycle_report_required",
    ]
    if not isinstance(value, dict):
        return False
    exception_ref = value.get("exception_ref")
    has_exception = isinstance(exception_ref, str) and bool(exception_ref.strip())
    review_gate_is_declared = isinstance(value.get("code_review_required"), bool)
    return review_gate_is_declared and all(value.get(field) is True or has_exception for field in required_true)


def expected_required_worker_chain(task: dict[str, Any]) -> list[str]:
    gates = task.get("lifecycle_gates")
    if isinstance(gates, dict) and gates.get("code_review_required") is True:
        return REVIEW_REQUIRED_WORKER_CHAIN
    return DEFAULT_REQUIRED_WORKER_CHAIN


def missing_ci_lane_refs(task: dict[str, Any]) -> list[str]:
    value = task.get("ci_lane_refs")
    if not isinstance(value, list):
        return list(REQUIRED_CI_LANES)
    present = {item.get("lane") for item in value if isinstance(item, dict)}
    return [lane for lane in REQUIRED_CI_LANES if lane not in present]


def missing_engineering_policy_refs(task: dict[str, Any]) -> list[str]:
    value = task.get("engineering_policy_refs")
    if not isinstance(value, list):
        return list(REQUIRED_ENGINEERING_POLICIES)
    present = {item.get("policy") for item in value if isinstance(item, dict)}
    return [policy for policy in REQUIRED_ENGINEERING_POLICIES if policy not in present]


def missing_architecture_policy_refs(task: dict[str, Any]) -> list[str]:
    value = task.get("architecture_guidance_refs")
    if not isinstance(value, list):
        return list(REQUIRED_ARCHITECTURE_POLICIES)
    combined = "\n".join(
        f"{item.get('source_ref', '')} {item.get('rule', '')}".lower()
        for item in value
        if isinstance(item, dict)
    )
    return [
        policy
        for policy, tokens in ARCHITECTURE_POLICY_TOKENS.items()
        if not any(token in combined for token in tokens)
    ]


def at_or_after_baseline(task: dict[str, Any], baseline_task_id: str) -> bool:
    baseline_key = task_order_key(baseline_task_id)
    if baseline_key is None:
        return False
    for value in [task_id(task), task.get("phase")]:
        candidate_key = task_order_key(value)
        if candidate_key is not None and version_gte(candidate_key, baseline_key):
            return True
    return False


def is_live_eval_dispatch_task(task: dict[str, Any]) -> bool:
    for value in [task_id(task), task.get("phase")]:
        key = task_order_key(value)
        if key is not None and key[:1] in {(11,), (12,)}:
            return True
    return False


def validate_task_guide_sources(task: dict[str, Any]) -> None:
    task_id_value = task_id(task)
    if not refs_include_base(task, "test_matrix_refs", TEST_MATRIX_SOURCE_BASE):
        fail(f"{task_id_value} test_matrix_refs must include source {TEST_MATRIX_SOURCE_BASE}")
    if object_source_ref_base(task.get("coverage_policy_refs")) != COVERAGE_POLICY_SOURCE_BASE:
        fail(f"{task_id_value} coverage_policy_refs must include source {COVERAGE_POLICY_SOURCE_BASE}")
    validate_coverage_policy_refs(task.get("coverage_policy_refs"), f"{task_id_value}.coverage_policy_refs")
    if not refs_include_base(task, "architecture_guidance_refs", ARCHITECTURE_GUIDE_BASE):
        fail(f"{task_id_value} architecture_guidance_refs must include source {ARCHITECTURE_GUIDE_BASE}")
    missing_lanes = missing_ci_lane_refs(task)
    if missing_lanes:
        fail(f"{task_id_value} ci_lane_refs missing: {', '.join(missing_lanes)}")
    missing_engineering = missing_engineering_policy_refs(task)
    if missing_engineering:
        fail(f"{task_id_value} engineering_policy_refs missing: {', '.join(missing_engineering)}")
    missing_architecture = missing_architecture_policy_refs(task)
    if missing_architecture:
        fail(f"{task_id_value} architecture_guidance_refs missing: {', '.join(missing_architecture)}")


def task_planning_field_has_evidence(task: dict[str, Any], field: str) -> bool:
    value = task.get(field)
    if field == "planning_skill_refs":
        return has_required_string_refs(value, REQUIRED_PLAN_SKILL_REFS)
    if field == "runtime_pins":
        return has_runtime_pins(value)
    if field == "slice_rationale":
        return (
            has_nonempty_dict(value)
            and has_nonempty_string(value.get("slice_type"))
            and has_nonempty_string(value.get("why_this_shape"))
        )
    if field == "changelog":
        if not has_nonempty_dict(value):
            return False
        return all(has_nonempty_string(value.get(changelog_field)) for changelog_field in REQUIRED_CHANGELOG_FIELDS)
    if field in ["contract_impact", "versioning_migration_impact"]:
        return has_nonempty_string(value)
    if field in ["architecture_constraints", "tdd_first_failing_tests", "semantic_invariants"]:
        return has_nonempty_list(value)
    if field == "adr_required":
        return isinstance(value, bool)
    if field == "cost_perf_impact":
        return (
            has_nonempty_dict(value)
            and has_nonempty_string(value.get("level"))
            and has_nonempty_string(value.get("measurement_expectation"))
        )
    if field == "chaos_failure_hypothesis":
        return (
            has_nonempty_dict(value)
            and has_nonempty_string(value.get("hypothesis"))
            and has_nonempty_string(value.get("expected_fail_closed_behavior"))
        )
    return False


def validate_task_planning_skill_fields(task: dict[str, Any]) -> None:
    task_id_value = task_id(task)
    missing = [
        field for field in REQUIRED_TASK_PLANNING_FIELDS if not task_planning_field_has_evidence(task, field)
    ]
    if missing:
        fail(f"{task_id_value} missing planning-skill fields: {', '.join(missing)}")
    if task_slice_type(task) != "vertical" and not has_nonempty_string(task.get("non_vertical_justification")):
        fail(f"{task_id_value} non-vertical task requires non_vertical_justification")
    contract_impact = str(task.get("contract_impact", ""))
    if has_adr_contract_token(contract_impact) and task.get("adr_required") is not True:
        fail(f"{task_id_value} public or executable contract impact requires adr_required=true")
    if contains_machine_local_path(task):
        fail(f"{task_id_value} contains a machine-local absolute path")


def validate_task_execution_compiler_fields(task: dict[str, Any]) -> None:
    task_id_value = task_id(task)
    validate_factory_compatibility(task.get("factory_compatibility"), f"{task_id_value}.factory_compatibility")
    validate_runtime_pins(task.get("runtime_pins"), f"{task_id_value}.runtime_pins")
    if not has_nonempty_list(task.get("scope_exclusions")):
        fail(f"{task_id_value}.scope_exclusions must preserve explicit PRD non-goals")
    if task.get("alignment_gate_ref") != ".factory/artifacts/prd-to-plan/lumyn-mvp/execution-plan.json#/alignment_gate":
        fail(f"{task_id_value}.alignment_gate_ref must cite the execution-plan alignment gate")
    if task.get("plan_drift_policy_ref") != ".factory/artifacts/prd-to-plan/lumyn-mvp/execution-plan.json#/plan_drift_policy":
        fail(f"{task_id_value}.plan_drift_policy_ref must cite the execution-plan drift policy")
    if task.get("acceptance_ledger_ref") != ACCEPTANCE_LEDGER_REF:
        fail(f"{task_id_value}.acceptance_ledger_ref must cite {ACCEPTANCE_LEDGER_REF}")
    item_ids = task.get("acceptance_item_ids")
    if not isinstance(item_ids, list) or not item_ids:
        fail(f"{task_id_value}.acceptance_item_ids must be non-empty")
    unknown_item_ids = sorted(str(value) for value in item_ids if str(value) not in REQUIRED_ACCEPTANCE_ITEM_IDS)
    if unknown_item_ids:
        fail(f"{task_id_value}.acceptance_item_ids references unknown ids: {unknown_item_ids}")
    validate_acceptance_item_gates(task)
    inherited = task.get("validation_contract_inheritance")
    if isinstance(inherited, dict):
        if inherited.get("acceptance_ledger_ref") != ACCEPTANCE_LEDGER_REF:
            fail(f"{task_id_value}.validation_contract_inheritance.acceptance_ledger_ref must cite {ACCEPTANCE_LEDGER_REF}")
        inherited_ids = inherited.get("acceptance_item_ids")
        if not isinstance(inherited_ids, list) or not set(str(value) for value in item_ids).issubset({str(value) for value in inherited_ids}):
            fail(f"{task_id_value}.validation_contract_inheritance.acceptance_item_ids must include task acceptance_item_ids")
    if not has_lifecycle_gates(task.get("lifecycle_gates")):
        fail(f"{task_id_value}.lifecycle_gates must enable local, CI, Codex review, ship, post-merge, and PR lifecycle gates, and explicitly declare code_review_required true only when review policy requires it")
    if task.get("required_worker_chain") != expected_required_worker_chain(task):
        fail(f"{task_id_value}.required_worker_chain must match the lifecycle gates: default validation/commit-push chain, or validation/code-review/commit-push chain when code_review_required=true")
    allowed_paths = [str(value).strip() for value in task.get("allowed_paths", [])]
    bad_allowed = [path for path in allowed_paths if RUNTIME_CONTROL_ALLOWED_RE.match(normalize_repo_path(path))]
    if bad_allowed:
        fail(f"{task_id_value}.allowed_paths includes runtime-owned control artifact paths: {bad_allowed}")
    forbidden_paths = set(str(value).strip() for value in task.get("forbidden_paths", []))
    missing_forbidden = [path for path in RUNTIME_CONTROL_FORBIDDEN_PATHS if path not in forbidden_paths]
    if missing_forbidden:
        fail(f"{task_id_value}.forbidden_paths missing runtime-owned control paths: {missing_forbidden}")


def validate_acceptance_item_gates(task: dict[str, Any]) -> None:
    task_id_value = task_id(task)
    gates = task.get("gated_by_acceptance_items")
    if gates is None:
        return
    if not isinstance(gates, list) or not gates:
        fail(f"{task_id_value}.gated_by_acceptance_items must be a non-empty list when present")
    seen: set[str] = set()
    for index, gate in enumerate(gates):
        if not isinstance(gate, dict):
            fail(f"{task_id_value}.gated_by_acceptance_items[{index}] must be an object")
        item_id_value = str(gate.get("acceptance_item_id", "")).strip()
        if not item_id_value:
            fail(f"{task_id_value}.gated_by_acceptance_items[{index}].acceptance_item_id is required")
        if item_id_value not in REQUIRED_ACCEPTANCE_ITEM_IDS:
            fail(f"{task_id_value}.gated_by_acceptance_items[{index}] references unknown acceptance item {item_id_value}")
        if item_id_value in seen:
            fail(f"{task_id_value}.gated_by_acceptance_items contains duplicate gate {item_id_value}")
        seen.add(item_id_value)
        if gate.get("required_status") not in {"implemented", "deferred_with_approval"}:
            fail(f"{task_id_value}.gated_by_acceptance_items[{index}].required_status must be implemented or deferred_with_approval")
        if not has_nonempty_string(gate.get("reason")):
            fail(f"{task_id_value}.gated_by_acceptance_items[{index}].reason is required")


def validate_live_eval_dispatch_gates(task: dict[str, Any]) -> None:
    task_id_value = task_id(task)
    gates = task.get("gated_by_acceptance_items")
    if not isinstance(gates, list):
        fail(f"{task_id_value}.gated_by_acceptance_items must gate live eval dispatch")
    by_id = {
        str(gate.get("acceptance_item_id")): gate
        for gate in gates
        if isinstance(gate, dict)
    }
    missing = sorted(REQUIRED_LIVE_EVAL_DISPATCH_GATES - set(by_id))
    if missing:
        fail(f"{task_id_value}.gated_by_acceptance_items missing live eval pull gates: {missing}")
    for required_id in REQUIRED_LIVE_EVAL_DISPATCH_GATES:
        gate = by_id[required_id]
        if gate.get("required_status") != "implemented":
            fail(f"{task_id_value}.gated_by_acceptance_items[{required_id}].required_status must be implemented")
        if gate.get("evidence_mode") != "product_signal":
            fail(f"{task_id_value}.gated_by_acceptance_items[{required_id}].evidence_mode must be product_signal")


def validate_task_version_slice_refs(task: dict[str, Any]) -> None:
    task_id_value = task_id(task)
    expected = expected_task_version_slices(task_id_value)
    if not expected:
        return
    actual = task.get("mvp_required_version_slices")
    if not isinstance(actual, list) or not actual:
        fail(f"{task_id_value}.mvp_required_version_slices must map task to required MVP version slices")
    missing = sorted(expected - {str(value) for value in actual})
    if missing:
        fail(f"{task_id_value}.mvp_required_version_slices missing {missing}")
    delivery_refs = task.get("delivery_slice_refs")
    if not isinstance(delivery_refs, list) or not delivery_refs:
        fail(f"{task_id_value}.delivery_slice_refs must map task to generic delivery slices")
    delivery_missing = sorted(expected - {str(value) for value in delivery_refs})
    if delivery_missing:
        fail(f"{task_id_value}.delivery_slice_refs missing {delivery_missing}")
    unexpected = sorted({str(value) for value in delivery_refs} - {str(value) for value in actual})
    if unexpected:
        fail(f"{task_id_value}.delivery_slice_refs has refs not present in mvp_required_version_slices: {unexpected}")


def field_has_evidence(task: dict[str, Any], field: str) -> bool:
    value = task.get(field)
    if field == "factory_compatibility":
        return has_factory_compatibility(value)
    if field == "scope_exclusions":
        return has_nonempty_list(value)
    if field in ["alignment_gate_ref", "plan_drift_policy_ref", "acceptance_ledger_ref"]:
        return has_nonempty_string(value)
    if field == "required_worker_chain":
        return value == expected_required_worker_chain(task)
    if field == "lifecycle_gates":
        return has_lifecycle_gates(value)
    if field in ["allowed_paths", "forbidden_paths"]:
        return has_nonempty_list(value)
    if field == "worker_type":
        return value == "codex_cli"
    if field == "factoryd_runtime":
        return is_valid_factoryd_runtime(value)
    if field in ["validation_commands", "evidence_required", "stop_conditions"]:
        return has_nonempty_list(value)
    if field == "security_scanner_gates":
        if not isinstance(value, dict):
            return False
        if isinstance(value.get("exception_ref"), str) and value["exception_ref"].strip():
            return True
        scanner = value.get("scanner")
        if not isinstance(scanner, str) or not scanner.strip():
            return False
        if value.get("required") is False:
            return isinstance(value.get("exception_ref"), str) and bool(value["exception_ref"].strip())
        return any(
            isinstance(value.get(key), str) and value[key].strip()
            for key in ["workflow_ref", "status_check", "evidence_ref"]
        )
    if field == "coverage_policy_refs":
        if not isinstance(value, dict):
            return False
        if value.get("required") is False and isinstance(value.get("exception_ref"), str) and value["exception_ref"].strip():
            return True
        if value.get("required") is not True:
            return False
        if has_nonempty_list(value.get("command_refs")) or has_nonempty_list(value.get("evidence_refs")):
            return True
        if isinstance(value.get("exception_ref"), str) and value["exception_ref"].strip():
            return True
        minimums = value.get("minimums")
        return isinstance(minimums, list) and any(
            isinstance(item, dict) and has_nonempty_list(item.get("command_refs"))
            for item in minimums
        )
    if not isinstance(value, list) or not value:
        return False
    if field == "acceptance_item_ids":
        return all(has_nonempty_string(item) for item in value)
    if field == "test_matrix_refs":
        return all(
            isinstance(item, dict)
            and isinstance(item.get("tier"), str)
            and item["tier"].strip()
            and isinstance(item.get("source_ref"), str)
            and item["source_ref"].strip()
            for item in value
        )
    if field == "ci_lane_refs":
        return all(
            isinstance(item, dict)
            and isinstance(item.get("lane"), str)
            and item["lane"].strip()
            and isinstance(item.get("source_ref"), str)
            and item["source_ref"].strip()
            and (
                has_nonempty_list(item.get("command_refs"))
                or has_nonempty_list(item.get("status_check_refs"))
                or (isinstance(item.get("exception_ref"), str) and item["exception_ref"].strip())
            )
            for item in value
        )
    if field == "engineering_policy_refs":
        return all(
            isinstance(item, dict)
            and isinstance(item.get("policy"), str)
            and item["policy"].strip()
            and isinstance(item.get("source_ref"), str)
            and item["source_ref"].strip()
            and (
                isinstance(item.get("rule"), str)
                and item["rule"].strip()
                or isinstance(item.get("exception_ref"), str)
                and item["exception_ref"].strip()
            )
            for item in value
        )
    if field == "architecture_guidance_refs":
        return all(
            isinstance(item, dict)
            and isinstance(item.get("source_ref"), str)
            and item["source_ref"].strip()
            and isinstance(item.get("rule"), str)
            and item["rule"].strip()
            for item in value
        )
    return False


def validate_context_brief(context: dict[str, Any]) -> None:
    validate_no_legacy_provider_fields(context, "context-brief.json")
    validate_factory_compatibility(context.get("factory_compatibility"), "context-brief.json.factory_compatibility")
    validate_alignment_gate(context.get("alignment_gate"), "context-brief.json.alignment_gate")
    validate_plan_drift_policy(context.get("plan_drift_policy"), "context-brief.json.plan_drift_policy")
    questions = context.get("alignment_questions")
    if not isinstance(questions, list):
        fail("context-brief.json.alignment_questions must be a list")
    question_ids = {question.get("id") for question in questions if isinstance(question, dict)}
    if LEGACY_PROVIDER_FIELD in question_ids:
        fail(f"context-brief.json must not use legacy alignment question id {LEGACY_PROVIDER_FIELD!r}")
    if REQUIRED_PROVIDER_DECISION_ID not in question_ids:
        fail(f"context-brief.json must include alignment question id {REQUIRED_PROVIDER_DECISION_ID!r}")
    decision_points = context.get("decision_points")
    if not isinstance(decision_points, list):
        fail("context-brief.json.decision_points must be a list")
    if LEGACY_PROVIDER_FIELD in decision_points:
        fail(f"context-brief.json must not use legacy decision point {LEGACY_PROVIDER_FIELD!r}")
    if REQUIRED_PROVIDER_DECISION_ID not in decision_points:
        fail(f"context-brief.json must include decision point {REQUIRED_PROVIDER_DECISION_ID!r}")
    decisions = context.get("alignment_decisions")
    if not isinstance(decisions, dict):
        fail("context-brief.json missing alignment_decisions")
    validate_factoryd_runtime(
        decisions.get("factoryd_runtime"),
        "context-brief.json.alignment_decisions.factoryd_runtime",
    )
    validate_mvp_eval_provider_adapters(
        decisions.get("mvp_eval_provider_adapters"),
        "context-brief.json.alignment_decisions.mvp_eval_provider_adapters",
    )
    if contains_machine_local_path(context):
        fail("context-brief.json contains a machine-local absolute path")


def validate_execution_plan(plan: dict[str, Any]) -> str:
    validate_no_legacy_provider_fields(plan, "execution-plan.json")
    validate_no_deprecated_active_workers(plan, "execution-plan.json")
    validate_factory_compatibility(plan.get("factory_compatibility"), "execution-plan.json.factory_compatibility")
    validate_runtime_pins(plan.get("runtime_pins"), "execution-plan.json.runtime_pins")
    validate_mvp_eval_provider_adapters(
        plan.get("mvp_eval_provider_adapters"),
        "execution-plan.json.mvp_eval_provider_adapters",
    )
    validate_factoryd_runtime(plan.get("factoryd_runtime"), "execution-plan.json.factoryd_runtime")
    validate_alignment_gate(plan.get("alignment_gate"), "execution-plan.json.alignment_gate")
    validate_plan_drift_policy(plan.get("plan_drift_policy"), "execution-plan.json.plan_drift_policy")
    validate_acceptance_ledger_coverage(
        plan.get("acceptance_ledger_coverage"),
        "execution-plan.json.acceptance_ledger_coverage",
    )
    validate_mvp_version_slice_coverage(
        plan.get("mvp_required_version_slices"),
        "execution-plan.json.mvp_required_version_slices",
    )
    validate_delivery_slice_coverage(
        plan.get("delivery_slices"),
        "execution-plan.json.delivery_slices",
    )
    for field in REQUIRED_PLAN_LEVEL_FIELDS:
        value = plan.get(field)
        if not has_nonempty_collection(value):
            fail(f"execution plan missing required planning-skill section {field}")
    alignment = plan.get("planning_skill_alignment")
    if not isinstance(alignment, dict):
        fail("execution plan planning_skill_alignment must be an object")
    if alignment.get("status") != "aligned":
        fail("execution plan planning_skill_alignment.status must be aligned")
    if not has_required_string_refs(alignment.get("source_refs"), REQUIRED_PLAN_SKILL_REFS):
        fail("execution plan planning_skill_alignment.source_refs must include Factory planning skills")
    if contains_machine_local_path(plan):
        fail("execution plan contains a machine-local absolute path")

    propagation = plan.get("dev_architecture_propagation")
    if not isinstance(propagation, dict):
        fail("execution plan missing dev_architecture_propagation")
    if propagation.get("status") != "proven":
        fail("dev_architecture_propagation.status must be proven before product implementation")
    baseline_task_id = propagation.get("baseline_task_ref")
    if not isinstance(baseline_task_id, str) or not baseline_task_id:
        fail("dev_architecture_propagation.baseline_task_ref is required")
    if propagation.get("test_matrix_source_ref") != "docs/dev/dev_guides.md#12-level-test-matrix":
        fail("dev_architecture_propagation.test_matrix_source_ref must point at docs/dev/dev_guides.md")
    if propagation.get("architecture_guide_ref") != "docs/architecture/architecture_guides.md":
        fail("dev_architecture_propagation.architecture_guide_ref must point at docs/architecture/architecture_guides.md")
    coverage_policy = propagation.get("coverage_policy")
    if not isinstance(coverage_policy, dict):
        fail("dev_architecture_propagation.coverage_policy must be an object")
    if coverage_policy.get("source_ref") != "docs/dev/dev_guides.md#coverage-gates":
        fail("dev_architecture_propagation.coverage_policy.source_ref must point at docs/dev/dev_guides.md#coverage-gates")
    if coverage_policy.get("required") is not True:
        fail("dev_architecture_propagation.coverage_policy.required must be true")
    if not has_nonempty_list(coverage_policy.get("command_refs")):
        fail("dev_architecture_propagation.coverage_policy.command_refs must be non-empty")
    validate_coverage_policy_refs(coverage_policy, "dev_architecture_propagation.coverage_policy")
    requirements = propagation.get("task_packet_requirements")
    if not isinstance(requirements, list):
        fail("dev_architecture_propagation.task_packet_requirements must be a list")
    missing = [field for field in REQUIRED_TASK_FIELDS if field not in requirements]
    if missing:
        fail(f"dev_architecture_propagation.task_packet_requirements missing {missing}")
    missing_planning = [field for field in REQUIRED_TASK_PLANNING_FIELDS if field not in requirements]
    if missing_planning:
        fail(f"dev_architecture_propagation.task_packet_requirements missing planning fields {missing_planning}")
    security_scanning = propagation.get("security_scanning")
    if not isinstance(security_scanning, dict):
        fail("dev_architecture_propagation.security_scanning must be an object")
    if not isinstance(security_scanning.get("required"), bool):
        fail("dev_architecture_propagation.security_scanning.required must be boolean")
    if not isinstance(security_scanning.get("scanner"), str) or not security_scanning["scanner"].strip():
        fail("dev_architecture_propagation.security_scanning.scanner must be non-empty")
    if security_scanning.get("required") is True and not any(
        isinstance(security_scanning.get(key), str) and security_scanning[key].strip()
        for key in ["workflow_ref", "status_check", "exception_policy", "exception_ref"]
    ):
        fail("required dev_architecture_propagation.security_scanning needs workflow/status/exception evidence")
    ci_lanes = propagation.get("ci_lanes")
    if not isinstance(ci_lanes, dict):
        fail("dev_architecture_propagation.ci_lanes must be an object")
    missing_lanes = [
        lane
        for lane in REQUIRED_CI_LANES
        if not isinstance(ci_lanes.get(lane), list) or not ci_lanes[lane]
    ]
    if missing_lanes:
        fail(f"dev_architecture_propagation.ci_lanes missing non-empty lanes: {missing_lanes}")
    engineering = propagation.get("engineering_policies")
    if not isinstance(engineering, dict):
        fail("dev_architecture_propagation.engineering_policies must be an object")
    missing_engineering = [
        policy
        for policy in REQUIRED_ENGINEERING_POLICIES
        if not isinstance(engineering.get(policy), str) or not engineering[policy].strip()
    ]
    if missing_engineering:
        fail(f"dev_architecture_propagation.engineering_policies missing {missing_engineering}")
    architecture = propagation.get("architecture_policies")
    if not isinstance(architecture, dict):
        fail("dev_architecture_propagation.architecture_policies must be an object")
    missing_architecture = [
        policy
        for policy in REQUIRED_ARCHITECTURE_POLICIES
        if not isinstance(architecture.get(policy), str) or not architecture[policy].strip()
    ]
    if missing_architecture:
        fail(f"dev_architecture_propagation.architecture_policies missing {missing_architecture}")
    return baseline_task_id


def validate_task_packets(packets: dict[str, Any], baseline_task_id: str) -> None:
    validate_no_legacy_provider_fields(packets, "task-packets.json")
    validate_no_deprecated_active_workers(packets, "task-packets.json")
    tasks = packets.get("tasks")
    if not isinstance(tasks, list):
        fail("task-packets.json must contain tasks list")
    tasks_by_id = {task_id(task): task for task in tasks if isinstance(task, dict) and task_id(task)}
    if baseline_task_id not in tasks_by_id:
        fail(f"task-packets.json missing propagation baseline task {baseline_task_id}")
    scoped_tasks = []
    all_task_objects = []
    baseline_has_order_key = task_order_key(baseline_task_id) is not None
    baseline_seen = False
    for task in tasks:
        if not isinstance(task, dict):
            continue
        candidate_id = task_id(task)
        if not candidate_id:
            continue
        all_task_objects.append(task)
        if candidate_id == baseline_task_id:
            baseline_seen = True
        depends_on_baseline = depends_on(candidate_id, baseline_task_id, tasks_by_id)
        ordered_after_baseline = at_or_after_baseline(task, baseline_task_id) or (
            not baseline_has_order_key and baseline_seen
        )
        if ordered_after_baseline and not depends_on_baseline:
            fail(f"{candidate_id} is at or after propagation baseline {baseline_task_id} but does not depend on it")
        if depends_on_baseline or ordered_after_baseline:
            scoped_tasks.append(task)
    if not scoped_tasks:
        fail(f"no task packets are at or after propagation baseline {baseline_task_id}")
    for task in all_task_objects:
        current_task_id = task_id(task)
        missing = [field for field in REQUIRED_TASK_FIELDS if not field_has_evidence(task, field)]
        if missing:
            fail(f"{current_task_id} missing guide propagation fields: {', '.join(missing)}")
        validate_task_guide_sources(task)
        validate_task_planning_skill_fields(task)
        validate_task_execution_compiler_fields(task)
        validate_task_version_slice_refs(task)
        if current_task_id == "T11":
            validate_mvp_eval_provider_adapters(
                task.get("mvp_eval_provider_adapters"),
                "T11.mvp_eval_provider_adapters",
            )
            checks = "\n".join(str(value).lower() for value in task.get("acceptance_checks", []))
            if "openai-compatible" not in checks or "anthropic" not in checks:
                fail("T11 acceptance_checks must name both OpenAI-compatible and Anthropic adapter coverage")
        if is_live_eval_dispatch_task(task):
            validate_live_eval_dispatch_gates(task)
    if any(task_ref in tasks_by_id for task_ref in ["T4", "T4.1", "T4.2", "T4.3"]):
        validate_recorder_task_split(tasks_by_id)
    if "T6" in tasks_by_id:
        validate_first_session_smoke_task(tasks_by_id)


def validate_recorder_task_split(tasks_by_id: dict[str, dict[str, Any]]) -> None:
    if "T4" in tasks_by_id:
        fail("task-packets.json must split broad T4 into T4.1, T4.2, and T4.3 before recorder dispatch")
    required = {
        "T4.1": {"RCRR-003", "RCRR-009", "FR8", "NFR5"},
        "T4.2": {"RCRR-004", "ACT-004", "FR5", "FR7", "NFR7", "NFR13"},
        "T4.3": {"REC-QUALITY-001", "ACT-003", "NFR2"},
    }
    for task_id_value, item_ids in required.items():
        task = tasks_by_id.get(task_id_value)
        if not isinstance(task, dict):
            fail(f"task-packets.json missing recorder split task {task_id_value}")
        actual_ids = {str(value) for value in task.get("acceptance_item_ids", [])}
        missing = sorted(item_ids - actual_ids)
        if missing:
            fail(f"{task_id_value}.acceptance_item_ids missing recorder split ids: {missing}")
        if "v0.0" not in {str(value) for value in task.get("delivery_slice_refs", [])}:
            fail(f"{task_id_value}.delivery_slice_refs must include v0.0")
    t42_deps = {str(value) for value in tasks_by_id["T4.2"].get("blocked_by", [])}
    if "T4.1" not in t42_deps:
        fail("T4.2 must depend on T4.1")
    t43_deps = {str(value) for value in tasks_by_id["T4.3"].get("blocked_by", [])}
    if "T4.2" not in t43_deps:
        fail("T4.3 must depend on T4.2")
    t5_deps = {str(value) for value in tasks_by_id["T5"].get("blocked_by", [])}
    if "T4.3" not in t5_deps or "T4" in t5_deps:
        fail("T5 must depend on T4.3, not the removed broad T4 packet")
    quality_checks = "\n".join(str(value).lower() for value in tasks_by_id["T4.3"].get("acceptance_checks", []))
    for token in ["strong-spec", "weak-spec", "non-claims", "70 percent"]:
        if token not in quality_checks:
            fail(f"T4.3 acceptance_checks must define REC-QUALITY-001 measurement token {token!r}")


def validate_first_session_smoke_task(tasks_by_id: dict[str, dict[str, Any]]) -> None:
    task = tasks_by_id.get("T6")
    if not isinstance(task, dict):
        fail("task-packets.json missing T6")
    smoke = task.get("first_session_smoke")
    if not isinstance(smoke, dict):
        fail("T6.first_session_smoke is required")
    if smoke.get("required") is not True:
        fail("T6.first_session_smoke.required must be true")
    if smoke.get("command") != "make smoke-first-session":
        fail("T6.first_session_smoke.command must be make smoke-first-session")
    if smoke.get("report_ref") != ".factory/artifacts/task-runs/T6/first-session-smoke.json":
        fail("T6.first_session_smoke.report_ref must be task-scoped")
    smoke_ids = {str(value) for value in smoke.get("acceptance_item_ids", [])}
    if not {"ACT-001", "ACT-002", "ACT-003"}.issubset(smoke_ids):
        fail("T6.first_session_smoke.acceptance_item_ids must cover ACT-001, ACT-002, and ACT-003")
    if "first_session_smoke_report" not in {str(value) for value in task.get("evidence_required", [])}:
        fail("T6.evidence_required must include first_session_smoke_report")
    if "make smoke-first-session" not in {str(value) for value in task.get("validation_commands", [])}:
        fail("T6.validation_commands must include make smoke-first-session")
    allowed_paths = {str(value) for value in task.get("allowed_paths", [])}
    missing_paths = sorted({"Makefile", "scripts/"} - allowed_paths)
    if missing_paths:
        fail(f"T6.allowed_paths must include smoke target implementation paths: {missing_paths}")


def validate_standalone_task_packet(packet: dict[str, Any], baseline_task_id: str) -> None:
    validate_no_legacy_provider_fields(packet, "standalone task packet")
    validate_no_deprecated_active_workers(packet, "standalone task packet")
    task_id_value = task_id(packet)
    if not task_id_value:
        fail("standalone task packet missing task_id")
    if not at_or_after_baseline(packet, baseline_task_id):
        return
    blocked_by = packet.get("blocked_by", [])
    if not isinstance(blocked_by, list) or baseline_task_id not in [str(value) for value in blocked_by]:
        fail(f"{task_id_value} is at or after propagation baseline {baseline_task_id} but does not depend on it")
    missing = [field for field in REQUIRED_TASK_FIELDS if not field_has_evidence(packet, field)]
    if missing:
        fail(f"{task_id_value} missing guide propagation fields: {', '.join(missing)}")
    validate_task_guide_sources(packet)
    validate_task_planning_skill_fields(packet)
    validate_task_execution_compiler_fields(packet)
    validate_task_version_slice_refs(packet)
    if is_live_eval_dispatch_task(packet):
        validate_live_eval_dispatch_gates(packet)


def validate_validation_contract(contract: dict[str, Any]) -> None:
    validate_no_legacy_provider_fields(contract, "validation-contract.json")
    validate_factory_compatibility(contract.get("factory_compatibility"), "validation-contract.json.factory_compatibility")
    validate_runtime_pins(contract.get("runtime_pins"), "validation-contract.json.runtime_pins")
    if contract.get("acceptance_ledger_ref") != ACCEPTANCE_LEDGER_REF:
        fail("validation-contract.json must cite acceptance-ledger.json")
    if contract.get("acceptance_item_count") != len(REQUIRED_ACCEPTANCE_ITEM_IDS):
        fail("validation-contract.json acceptance_item_count must match acceptance-ledger item count")
    if not has_nonempty_list(contract.get("acceptance_criteria")):
        fail("validation-contract.json must include itemized acceptance_criteria")
    validate_mvp_version_slice_coverage(
        contract.get("mvp_required_version_slices"),
        "validation-contract.json.mvp_required_version_slices",
    )
    validate_delivery_slice_coverage(
        contract.get("delivery_slices"),
        "validation-contract.json.delivery_slices",
    )
    validate_mvp_eval_provider_adapters(
        contract.get("mvp_eval_provider_adapters"),
        "validation-contract.json.mvp_eval_provider_adapters",
    )
    validate_plan_drift_policy(contract.get("plan_drift_policy"), "validation-contract.json.plan_drift_policy")
    validate_coverage_policy_refs(contract.get("coverage_policy"), "validation-contract.json.coverage_policy")
    alignment = contract.get("planning_skill_alignment")
    if not isinstance(alignment, dict):
        fail("validation-contract.json missing planning_skill_alignment")
    if not has_required_string_refs(alignment.get("source_refs"), REQUIRED_PLAN_SKILL_REFS):
        fail("validation-contract.json planning_skill_alignment.source_refs must include Factory planning skills")
    required_plan_sections = alignment.get("required_plan_sections")
    if not isinstance(required_plan_sections, list) or not all(
        section in required_plan_sections for section in REQUIRED_PLAN_LEVEL_FIELDS
    ):
        fail("validation-contract.json planning_skill_alignment.required_plan_sections is incomplete")
    required_task_fields = alignment.get("required_task_fields")
    if not isinstance(required_task_fields, list) or not all(
        field in required_task_fields for field in REQUIRED_TASK_PLANNING_FIELDS
    ):
        fail("validation-contract.json planning_skill_alignment.required_task_fields is incomplete")
    missing_execution_fields = [field for field in REQUIRED_TASK_FIELDS if field not in required_task_fields]
    if missing_execution_fields:
        fail(
            "validation-contract.json planning_skill_alignment.required_task_fields missing "
            f"execution-compiler fields: {missing_execution_fields}"
        )
    factoryd_requirements = contract.get("factoryd_runtime_requirements")
    if not isinstance(factoryd_requirements, dict):
        fail("validation-contract.json missing factoryd_runtime_requirements")
    missing_runner_ready = [
        field
        for field in REQUIRED_RUNNER_READY_FIELDS
        if field not in factoryd_requirements.get("runner_ready_fields", [])
    ]
    if missing_runner_ready:
        fail(f"validation-contract.json.factoryd_runtime_requirements.runner_ready_fields missing {missing_runner_ready}")
    validate_factoryd_runtime(
        factoryd_requirements.get("runtime"),
        "validation-contract.json.factoryd_runtime_requirements.runtime",
    )
    if contains_machine_local_path(contract):
        fail("validation-contract.json contains a machine-local absolute path")
    stop_conditions = contract.get("stop_conditions")
    if not isinstance(stop_conditions, list):
        fail("validation-contract.json must contain stop_conditions list")
    combined = "\n".join(str(value).lower() for value in stop_conditions)
    missing = [
        category
        for category, tokens in STOP_CONDITION_CATEGORIES.items()
        if not any(token in combined for token in tokens)
    ]
    if missing:
        fail(f"validation-contract.json stop_conditions missing guide categories: {missing}")


def validate_safety_corpus_ready_plan(
    plan: dict[str, Any],
    packets: dict[str, Any],
    contract: dict[str, Any],
    ledger: dict[str, Any],
    mapping: dict[str, Any],
    scope: dict[str, Any],
) -> None:
    prd_text = PRD.read_text().lower()
    for token in [
        "safety and corpus-ready evidence",
        "corpus_eligible: false",
        "boundary violations",
    ]:
        if token not in prd_text:
            fail(f"docs/product/prd.md missing safety/corpus-ready token: {token}")

    required_by_task = {
        "T2.7": {"FDN-003"},
        "T6": {"RCRR-012"},
        "T8": {"LVCIS-009"},
        "T9": {"LVCIS-010"},
        "T11": {"EVAL-011"},
        "T12": {"EVAL-011"},
    }
    required_ids = set().union(*required_by_task.values())
    ledger_items = {
        str(item.get("acceptance_item_id")): item
        for item in ledger.get("items", [])
        if isinstance(item, dict)
    }
    missing_ledger_ids = sorted(required_ids - set(ledger_items))
    if missing_ledger_ids:
        fail(f"acceptance-ledger.json missing safety/corpus ids: {missing_ledger_ids}")
    for item_id in required_ids:
        source_ref = str(ledger_items[item_id].get("source_ref", ""))
        if source_ref != "docs/product/prd.md#safety-and-corpus-ready-evidence":
            fail(f"acceptance-ledger.json {item_id} must cite the safety/corpus PRD section")

    criteria = "\n".join(str(value) for value in contract.get("acceptance_criteria", []))
    for item_id in sorted(required_ids):
        if item_id not in criteria:
            fail(f"validation-contract.json acceptance_criteria missing {item_id}")

    def delivery_slice_ids(document: dict[str, Any], slice_id: str) -> set[str]:
        slice_item = next(
            (
                item
                for item in document.get("delivery_slices", [])
                if isinstance(item, dict) and item.get("slice_id") == slice_id
            ),
            None,
        )
        if not isinstance(slice_item, dict):
            fail(f"{slice_id} delivery slice missing")
        return {str(value) for value in slice_item.get("acceptance_item_ids", [])}

    for label, document in [
        ("execution-plan.json", plan),
        ("validation-contract.json", contract),
        ("acceptance-mapping.json", mapping),
        ("scope-closure-map.json", scope),
    ]:
        if not {"FDN-003", "RCRR-012"}.issubset(delivery_slice_ids(document, "v0.0")):
            fail(f"{label} v0.0 delivery slice missing safety/corpus record/report ids")
        if not {"LVCIS-009", "LVCIS-010"}.issubset(delivery_slice_ids(document, "v0.1")):
            fail(f"{label} v0.1 delivery slice missing boundary/CI safety ids")
        if "EVAL-011" not in delivery_slice_ids(document, "v0.2"):
            fail(f"{label} v0.2 delivery slice missing eval failure-event id")

    tasks = {
        str(task.get("task_id")): task
        for task in packets.get("tasks", [])
        if isinstance(task, dict)
    }
    if "T2.7" not in tasks:
        fail("task-packets.json missing T2.7 safety/corpus contract task")
    if "T2.7" not in set(str(value) for value in tasks["T3"].get("blocked_by", [])):
        fail("T3 must depend on T2.7 before product implementation resumes")
    for task_id_value, ids in required_by_task.items():
        task = tasks.get(task_id_value)
        if not isinstance(task, dict):
            fail(f"task-packets.json missing {task_id_value}")
        task_ids = {str(value) for value in task.get("acceptance_item_ids", [])}
        missing = sorted(ids - task_ids)
        if missing:
            fail(f"{task_id_value}.acceptance_item_ids missing safety/corpus ids: {missing}")
        checks = "\n".join(str(value).lower() for value in task.get("acceptance_checks", []))
        if not any(token in checks for token in ["safety", "corpus", "boundary"]):
            fail(f"{task_id_value}.acceptance_checks must describe safety/corpus-ready evidence")
        refs = task.get("safety_corpus_ready_evidence_refs")
        if not isinstance(refs, list) or not refs:
            fail(f"{task_id_value}.safety_corpus_ready_evidence_refs must be non-empty")
        ref_text = "\n".join(
            f"{item.get('source_ref', '')} {item.get('rule', '')}".lower()
            for item in refs
            if isinstance(item, dict)
        )
        if "safety-and-corpus-ready-evidence" not in ref_text:
            fail(f"{task_id_value}.safety_corpus_ready_evidence_refs must cite the PRD safety/corpus section")

    plan_text = "\n".join(
        str(value).lower()
        for value in plan.get("definition_of_done", []) + plan.get("explicit_non_goals", [])
    )
    for token in ["corpus_eligible", "boundary", "hosted corpus"]:
        if token not in plan_text:
            fail(f"execution-plan.json missing safety/corpus plan token: {token}")


def validate_factoryd_config(config: dict[str, Any], active_config: dict[str, Any], autoship_config: dict[str, Any]) -> None:
    if contains_machine_local_path(config):
        fail(".factory/factoryd.example.json contains a machine-local absolute path")
    if contains_machine_local_path(active_config):
        fail(".factory/factoryd.json contains a machine-local absolute path")
    if contains_machine_local_path(autoship_config):
        fail(".factory/factoryd.autoship.example.json contains a machine-local absolute path")
    repos = config.get("repos")
    if not isinstance(repos, dict) or "lumyn" not in repos:
        fail(".factory/factoryd.example.json must define repos.lumyn")
    lumyn = repos["lumyn"]
    if not isinstance(lumyn, dict):
        fail(".factory/factoryd.example.json repos.lumyn must be an object")
    expected_paths = {
        "repo_path": "..",
        "acceptance_ledger": ACCEPTANCE_LEDGER_REF,
        "task_packets": ".factory/artifacts/prd-to-plan/lumyn-mvp/task-packets.json",
        "scope_closure_map": ".factory/artifacts/prd-to-plan/lumyn-mvp/scope-closure-map.json",
        "validation_contract": ".factory/artifacts/prd-to-plan/lumyn-mvp/validation-contract.json",
        "state_dir": "../.factoryd",
        "workspace_root": "../.factoryd/workspaces",
    }
    for key, expected in expected_paths.items():
        if lumyn.get(key) != expected:
            fail(f".factory/factoryd.example.json repos.lumyn.{key} must be {expected!r}")
    validate_factoryd_runtime(lumyn, ".factory/factoryd.example.json repos.lumyn")
    commands = lumyn.get("validation_commands")
    if not isinstance(commands, list) or "python3 scripts/validate_repo_pack.py" not in commands:
        fail(".factory/factoryd.example.json must run validate_repo_pack.py")
    shipping = lumyn.get("shipping")
    if not isinstance(shipping, dict):
        fail(".factory/factoryd.example.json repos.lumyn must declare shipping block")
    if shipping.get("enabled") is not False or lumyn.get("auto_ship") is not False:
        fail(".factory/factoryd.example.json must keep auto shipping disabled until remote lifecycle hooks are approved")
    for key in [
        "push_required",
        "pr_required",
        "ci_required",
        "codex_review_required",
        "merge_required",
        "post_merge_required",
        "scope_closure_required",
    ]:
        if shipping.get(key) is not False:
            fail(f".factory/factoryd.example.json shipping.{key} must be false until hooks are approved")
    for key in [
        "push_command",
        "open_pr_command",
        "ci_command",
        "codex_review_command",
        "merge_command",
        "post_merge_command",
        "scope_closure_command",
    ]:
        if shipping.get(key) != "":
            fail(f".factory/factoryd.example.json shipping.{key} must be empty until hooks are approved")
    active_repos = active_config.get("repos")
    if not isinstance(active_repos, dict) or "lumyn" not in active_repos:
        fail(".factory/factoryd.json must define repos.lumyn")
    if active_repos["lumyn"] != lumyn:
        fail(".factory/factoryd.json repos.lumyn must match the safe attended example config")
    active_factory = active_config.get("factory")
    if not isinstance(active_factory, dict) or active_factory.get("repo_path") != "../../factory":
        fail(".factory/factoryd.json factory.repo_path must point to sibling ../../factory")
    autoship_repos = autoship_config.get("repos")
    if not isinstance(autoship_repos, dict) or "lumyn" not in autoship_repos:
        fail(".factory/factoryd.autoship.example.json must define repos.lumyn")
    autoship_lumyn = autoship_repos["lumyn"]
    if not isinstance(autoship_lumyn, dict):
        fail(".factory/factoryd.autoship.example.json repos.lumyn must be an object")
    for key, expected in expected_paths.items():
        if autoship_lumyn.get(key) != expected:
            fail(f".factory/factoryd.autoship.example.json repos.lumyn.{key} must be {expected!r}")
    validate_factoryd_runtime(autoship_lumyn, ".factory/factoryd.autoship.example.json repos.lumyn")
    autoship_shipping = autoship_lumyn.get("shipping")
    if not isinstance(autoship_shipping, dict):
        fail(".factory/factoryd.autoship.example.json repos.lumyn must declare shipping block")
    if autoship_lumyn.get("auto_ship") is not True or autoship_shipping.get("enabled") is not True:
        fail(".factory/factoryd.autoship.example.json must explicitly enable auto shipping")
    if autoship_shipping.get("provider") != "github_cli":
        fail(".factory/factoryd.autoship.example.json shipping.provider must be github_cli")
    for key in [
        "push_required",
        "pr_required",
        "ci_required",
        "codex_review_required",
        "merge_required",
        "post_merge_required",
        "scope_closure_required",
    ]:
        if autoship_shipping.get(key) is not True:
            fail(f".factory/factoryd.autoship.example.json shipping.{key} must be true")
    if autoship_shipping.get("scope_closure_mode") != "semantic":
        fail(".factory/factoryd.autoship.example.json shipping.scope_closure_mode must be semantic")
    if ".factoryd/" not in (ROOT / ".gitignore").read_text():
        fail(".gitignore must ignore .factoryd/")


def validate_acceptance_ledger(ledger: dict[str, Any]) -> set[str]:
    validate_no_legacy_provider_fields(ledger, "acceptance-ledger.json")
    if ledger.get("artifact_type") != "acceptance_ledger":
        fail("acceptance-ledger.json artifact_type must be acceptance_ledger")
    if ledger.get("source_prd_ref") != "docs/product/prd.md":
        fail("acceptance-ledger.json source_prd_ref must point at docs/product/prd.md")
    policy = ledger.get("coverage_policy")
    if not isinstance(policy, dict):
        fail("acceptance-ledger.json coverage_policy must be an object")
    if policy.get("enumerated_items_required") is not True:
        fail("acceptance-ledger.json coverage_policy.enumerated_items_required must be true")
    if policy.get("group_only_refs_allowed") is not False:
        fail("acceptance-ledger.json coverage_policy.group_only_refs_allowed must be false")
    if policy.get("closure_unit") != "acceptance_item":
        fail("acceptance-ledger.json coverage_policy.closure_unit must be acceptance_item")
    items = ledger.get("items")
    if not isinstance(items, list) or not items:
        fail("acceptance-ledger.json must contain items")
    seen: set[str] = set()
    task_refs_by_item_id: dict[str, set[str]] = {}
    for index, item in enumerate(items):
        if not isinstance(item, dict):
            fail(f"acceptance-ledger.json items[{index}] must be an object")
        item_id = str(item.get("acceptance_item_id", "")).strip()
        if not item_id:
            fail(f"acceptance-ledger.json items[{index}] missing acceptance_item_id")
        if item_id in seen:
            fail(f"acceptance-ledger.json duplicate acceptance_item_id {item_id}")
        seen.add(item_id)
        for key in ["group_id", "source_ref", "source_text", "kind", "evidence_mode", "status", "risk_class"]:
            if not has_nonempty_string(item.get(key)):
                fail(f"acceptance-ledger.json {item_id} missing {key}")
        if not has_nonempty_list(item.get("closure_required_for")):
            fail(f"acceptance-ledger.json {item_id} missing closure_required_for")
        if not has_nonempty_list(item.get("task_refs")):
            fail(f"acceptance-ledger.json {item_id} missing task_refs")
        task_refs_by_item_id[item_id] = {str(value) for value in item.get("task_refs", [])}
        if item.get("status") == "implemented" and not has_nonempty_list(item.get("validation_refs")):
            fail(f"acceptance-ledger.json {item_id} implemented item missing validation_refs")
    missing = sorted(REQUIRED_ACCEPTANCE_ITEM_IDS - seen)
    if missing:
        fail(f"acceptance-ledger.json missing required item ids: {missing}")
    if "REC-QUALITY-001" not in seen:
        fail("acceptance-ledger.json must include recorder 70 percent quality gate")
    recorder_quality = next(
        (item for item in items if isinstance(item, dict) and item.get("acceptance_item_id") == "REC-QUALITY-001"),
        None,
    )
    if not isinstance(recorder_quality, dict):
        fail("acceptance-ledger.json missing REC-QUALITY-001")
    if recorder_quality.get("source_ref") != "docs/product/prd.md#phase-1-recorder-spike":
        fail("REC-QUALITY-001 source_ref must point at the Phase 1 recorder spike PRD anchor")
    for item_id, required_task_refs in REQUIRED_ACCEPTANCE_TASK_REFS.items():
        actual_task_refs = task_refs_by_item_id.get(item_id, set())
        missing_task_refs = sorted(required_task_refs - actual_task_refs)
        if missing_task_refs:
            fail(f"acceptance-ledger.json {item_id} missing required task_refs: {missing_task_refs}")
    return seen


def validate_acceptance_mapping(mapping: dict[str, Any], ledger_ids: set[str], contract: dict[str, Any]) -> None:
    validate_no_legacy_provider_fields(mapping, "acceptance-mapping.json")
    if mapping.get("acceptance_ledger_ref") != ACCEPTANCE_LEDGER_REF:
        fail("acceptance-mapping.json must cite acceptance-ledger.json")
    validate_mvp_version_slice_coverage(
        mapping.get("mvp_required_version_slices"),
        "acceptance-mapping.json.mvp_required_version_slices",
    )
    validate_delivery_slice_coverage(
        mapping.get("delivery_slices"),
        "acceptance-mapping.json.delivery_slices",
    )
    groups = mapping.get("groups")
    if not isinstance(groups, list):
        fail("acceptance-mapping.json must contain groups list")
    mapped_ids: set[str] = set()
    mapped_groups: set[str] = set()
    for group in groups:
        if not isinstance(group, dict):
            continue
        group_id = str(group.get("group_id", "")).strip()
        if not group_id:
            fail("acceptance-mapping.json group missing group_id")
        mapped_groups.add(group_id)
        item_ids = group.get("acceptance_item_ids")
        if not isinstance(item_ids, list) or not item_ids:
            fail(f"acceptance group {group.get('group_id')} missing acceptance_item_ids")
        unknown = sorted(str(value) for value in item_ids if str(value) not in ledger_ids)
        if unknown:
            fail(f"acceptance group {group.get('group_id')} references unknown acceptance item ids: {unknown}")
        mapped_ids.update(str(value) for value in item_ids)
    live_eval = next((group for group in groups if isinstance(group, dict) and group.get("group_id") == "live_agent_eval"), None)
    if not isinstance(live_eval, dict):
        fail("acceptance-mapping.json missing live_agent_eval group")
    if not has_required_string_refs(live_eval.get("provider_adapter_coverage"), REQUIRED_MVP_EVAL_PROVIDERS):
        fail("live_agent_eval.provider_adapter_coverage must include both MVP provider adapters")
    approvals = "\n".join(str(value).lower() for value in live_eval.get("requires_human_approval", []))
    if "openai" not in approvals or "anthropic" not in approvals:
        fail("live_agent_eval.requires_human_approval must name both provider credential postures")
    for required_id in ["EVAL-001", "PULL-001", "PULL-004"]:
        if required_id not in set(str(value) for value in live_eval.get("acceptance_item_ids", [])):
            fail(f"live_agent_eval acceptance mapping missing {required_id}")
    contract_groups = contract.get("acceptance_groups")
    if not isinstance(contract_groups, list) or not contract_groups:
        fail("validation-contract.json must declare acceptance_groups")
    missing_contract_groups = sorted(str(group_id) for group_id in contract_groups if str(group_id) not in mapped_groups)
    if missing_contract_groups:
        fail(f"acceptance-mapping.json missing validation-contract groups: {missing_contract_groups}")
    nfr_group = next((group for group in groups if isinstance(group, dict) and group.get("group_id") == "nonfunctional_requirements"), None)
    if not isinstance(nfr_group, dict):
        fail("acceptance-mapping.json missing nonfunctional_requirements group")
    if nfr_group.get("source_ref") != "docs/product/prd.md#non-functional-requirements":
        fail("nonfunctional_requirements source_ref must match PRD heading anchor")
    missing = sorted(ledger_ids - mapped_ids)
    if missing:
        fail(f"acceptance-mapping.json does not map ledger ids: {missing}")


def validate_scope_closure_map(scope: dict[str, Any], ledger_ids: set[str]) -> None:
    validate_no_legacy_provider_fields(scope, "scope-closure-map.json")
    if scope.get("acceptance_ledger_ref") != ACCEPTANCE_LEDGER_REF:
        fail("scope-closure-map.json must cite acceptance-ledger.json")
    validate_mvp_version_slice_coverage(
        scope.get("mvp_required_version_slices"),
        "scope-closure-map.json.mvp_required_version_slices",
    )
    validate_delivery_slice_coverage(
        scope.get("delivery_slices"),
        "scope-closure-map.json.delivery_slices",
    )
    items = scope.get("items")
    if not isinstance(items, list):
        fail("scope-closure-map.json must contain items list")
    closure_ids: set[str] = set()
    for item in items:
        if not isinstance(item, dict):
            continue
        item_ids = item.get("acceptance_item_ids")
        statuses = item.get("acceptance_item_statuses")
        if not isinstance(item_ids, list) or not item_ids:
            fail(f"scope item {item.get('scope_item')} missing acceptance_item_ids")
        if not isinstance(statuses, list) or not statuses:
            fail(f"scope item {item.get('scope_item')} missing acceptance_item_statuses")
        scope_item_id = str(item.get("scope_item_id", "")).strip()
        for slice_id, spec in REQUIRED_MVP_VERSION_SLICES.items():
            if scope_item_id != spec["capability_group_id"]:
                continue
            actual_slices = {str(value) for value in item.get("mvp_required_version_slices", [])}
            if slice_id not in actual_slices:
                fail(f"scope item {scope_item_id} missing mvp_required_version_slices entry {slice_id}")
            delivery_slices = {str(value) for value in item.get("delivery_slice_refs", [])}
            if slice_id not in delivery_slices:
                fail(f"scope item {scope_item_id} missing delivery_slice_refs entry {slice_id}")
        status_ids = {str(status.get("acceptance_item_id")) for status in statuses if isinstance(status, dict)}
        missing_status = sorted(str(value) for value in item_ids if str(value) not in status_ids)
        if missing_status:
            fail(f"scope item {item.get('scope_item')} missing item statuses for {missing_status}")
        closure_ids.update(str(value) for value in item_ids)
    live_eval = next((item for item in items if isinstance(item, dict) and item.get("scope_item") == "Live agent eval"), None)
    if not isinstance(live_eval, dict):
        fail("scope-closure-map.json missing Live agent eval item")
    if not has_required_string_refs(live_eval.get("provider_adapter_coverage"), REQUIRED_MVP_EVAL_PROVIDERS):
        fail("Live agent eval provider_adapter_coverage must include both MVP provider adapters")
    blockers = "\n".join(str(value).lower() for value in live_eval.get("blockers", []))
    if "openai" not in blockers or "anthropic" not in blockers:
        fail("Live agent eval blockers must name both provider credential postures")
    missing = sorted(ledger_ids - closure_ids)
    if missing:
        fail(f"scope-closure-map.json does not cover ledger ids: {missing}")
    nfr_item = next((item for item in items if isinstance(item, dict) and item.get("scope_item_id") == "nonfunctional_requirements"), None)
    if not isinstance(nfr_item, dict):
        fail("scope-closure-map.json missing nonfunctional_requirements scope item")
    if nfr_item.get("source_ref") != "docs/product/prd.md#non-functional-requirements":
        fail("nonfunctional_requirements scope source_ref must match PRD heading anchor")


def validate_risk_classification(risk: dict[str, Any]) -> None:
    validate_no_legacy_provider_fields(risk, "risk-classification.json")
    rules = risk.get("risk_rules")
    if not isinstance(rules, list):
        fail("risk-classification.json must contain risk_rules list")
    high = next((rule for rule in rules if isinstance(rule, dict) and rule.get("risk_class") == "high"), None)
    if not isinstance(high, dict):
        fail("risk-classification.json missing high risk rule")
    applies = "\n".join(str(value).lower() for value in high.get("applies_to", []))
    if "openai-compatible" not in applies or "anthropic" not in applies:
        fail("high risk rule must name both OpenAI-compatible and Anthropic provider key surfaces")


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
            "ship_pr_required": True,
            "post_merge_monitor_required": True,
            "pr_lifecycle_report_required": True,
            "skip_policy": "approved_exception_required",
        },
        "allowed_paths": [
            "cmd/",
            "internal/",
            "schemas/",
            "tests/",
            "docs/",
            f".factory/artifacts/task-runs/{task_id_value}/",
            f".factory/artifacts/pr-lifecycle/{task_id_value}/",
        ],
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
        "evidence_required": ["validation_report", "work_proof_marker", "factoryd_run_once_report"],
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


def run_self_test() -> int:
    valid_packets = {"tasks": [propagated_task("T2.6", ["T2.5"]), propagated_task("T3", ["T2.6"])]}
    validate_task_packets(valid_packets, "T2.6")
    if expected_task_version_slices("T11.1") != {"v0.2"}:
        fail("self-test expected dotted live-eval task to inherit parent v0.2 delivery slice")
    if expected_task_version_slices("T2.5"):
        fail("self-test expected T2.5 lifecycle baseline task to remain delivery-slice exempt")
    if contains_machine_local_path(".factory/artifacts/prd-to-plan/lumyn-mvp/execution-plan.json#/alignment_gate"):
        fail("self-test expected repo-relative JSON pointer to remain portable")

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
        ".factory/artifacts/task-runs/T3/../../prd-to-plan/lumyn-mvp/scope-closure-map.json"
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


def main() -> int:
    if sys.argv[1:] == ["--self-test"]:
        try:
            return run_self_test()
        except AssertionError as exc:
            print(f"repo-pack validator self-test failed: {exc}", file=sys.stderr)
            return 2
    if sys.argv[1:]:
        print("usage: validate_repo_pack.py [--self-test]", file=sys.stderr)
        return 2
    try:
        validate_guides()
        validate_ci_control_set()
        context = load_json(CONTEXT_BRIEF)
        plan = load_json(EXECUTION_PLAN)
        packets = load_json(TASK_PACKETS)
        contract = load_json(VALIDATION_CONTRACT)
        factoryd_config = load_json(FACTORYD_CONFIG)
        factoryd_active_config = load_json(FACTORYD_ACTIVE_CONFIG)
        factoryd_autoship_config = load_json(FACTORYD_AUTOSHIP_CONFIG)
        acceptance_ledger = load_json(ACCEPTANCE_LEDGER)
        acceptance_mapping = load_json(ACCEPTANCE_MAPPING)
        scope_closure_map = load_json(SCOPE_CLOSURE_MAP)
        risk_classification = load_json(RISK_CLASSIFICATION)
        validate_context_brief(context)
        baseline_task_id = validate_execution_plan(plan)
        ledger_ids = validate_acceptance_ledger(acceptance_ledger)
        validate_task_packets(packets, baseline_task_id)
        for packet_path in REPAIR_TASK_PACKETS:
            validate_standalone_task_packet(load_json(packet_path), baseline_task_id)
        validate_validation_contract(contract)
        validate_factoryd_config(factoryd_config, factoryd_active_config, factoryd_autoship_config)
        validate_acceptance_mapping(acceptance_mapping, ledger_ids, contract)
        validate_scope_closure_map(scope_closure_map, ledger_ids)
        validate_safety_corpus_ready_plan(
            plan,
            packets,
            contract,
            acceptance_ledger,
            acceptance_mapping,
            scope_closure_map,
        )
        validate_risk_classification(risk_classification)
    except AssertionError as exc:
        print(f"repo-pack validation failed: {exc}", file=sys.stderr)
        return 2
    print("repo-pack validation passed")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
