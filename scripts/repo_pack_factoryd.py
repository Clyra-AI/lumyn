#!/usr/bin/env python3
from __future__ import annotations

import re
from pathlib import Path
from typing import Any

from repo_pack_architecture import validate_architecture_budget_policy
from repo_pack_contracts import ACCEPTANCE_LEDGER_REF, fail, has_nonempty_string


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

MACHINE_LOCAL_PATH_RE = re.compile(r"(?<![A-Za-z0-9+./#-])(?:/(?!/)[A-Za-z0-9._-][^\s\"'<>]*|[A-Za-z]:[\\/])")


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


def validate_factoryd_config(
    root: Path,
    config: dict[str, Any],
    active_config: dict[str, Any],
    autoship_config: dict[str, Any],
    repo_key: str = "lumyn",
) -> None:
    if contains_machine_local_path(config):
        fail(".factory/factoryd.example.json contains a machine-local absolute path")
    if active_config and contains_machine_local_path(active_config):
        fail(".factory/factoryd.json contains a machine-local absolute path")
    if contains_machine_local_path(autoship_config):
        fail(".factory/factoryd.autoship.example.json contains a machine-local absolute path")
    repos = config.get("repos")
    if not isinstance(repos, dict) or repo_key not in repos:
        fail(f".factory/factoryd.example.json must define repos.{repo_key}")
    lumyn = repos[repo_key]
    if not isinstance(lumyn, dict):
        fail(f".factory/factoryd.example.json repos.{repo_key} must be an object")
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
            fail(f".factory/factoryd.example.json repos.{repo_key}.{key} must be {expected!r}")
    validate_factoryd_runtime(lumyn, f".factory/factoryd.example.json repos.{repo_key}")
    validate_architecture_budget_policy(root, lumyn, f".factory/factoryd.example.json repos.{repo_key}")
    commands = lumyn.get("validation_commands")
    if not isinstance(commands, list) or "python3 scripts/validate_repo_pack.py" not in commands:
        fail(".factory/factoryd.example.json must run validate_repo_pack.py")
    shipping = lumyn.get("shipping")
    if not isinstance(shipping, dict):
        fail(f".factory/factoryd.example.json repos.{repo_key} must declare shipping block")
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
    if active_config:
        active_repos = active_config.get("repos")
        if not isinstance(active_repos, dict) or repo_key not in active_repos:
            fail(f".factory/factoryd.json must define repos.{repo_key}")
        active_lumyn = active_repos[repo_key]
        if not isinstance(active_lumyn, dict):
            fail(f".factory/factoryd.json repos.{repo_key} must be an object")
        for key, expected in expected_paths.items():
            if active_lumyn.get(key) != expected:
                fail(f".factory/factoryd.json repos.{repo_key}.{key} must be {expected!r}")
        validate_factoryd_runtime(active_lumyn, f".factory/factoryd.json repos.{repo_key}")
        validate_architecture_budget_policy(root, active_lumyn, f".factory/factoryd.json repos.{repo_key}")
        active_commands = active_lumyn.get("validation_commands")
        if not isinstance(active_commands, list) or "python3 scripts/validate_repo_pack.py" not in active_commands:
            fail(".factory/factoryd.json must run validate_repo_pack.py")
        active_shipping = active_lumyn.get("shipping")
        if not isinstance(active_shipping, dict):
            fail(f".factory/factoryd.json repos.{repo_key} must declare shipping block")
        if active_lumyn.get("auto_ship") is not False or active_shipping.get("enabled") is not False:
            fail(".factory/factoryd.json must remain safe-attended; use factoryd.autoship.example.json for full-loop shipping")
        active_factory = active_config.get("factory")
        if not isinstance(active_factory, dict):
            fail(".factory/factoryd.json must define factory")
        if not has_nonempty_string(active_factory.get("repo_path")):
            fail(".factory/factoryd.json factory.repo_path must be non-empty")
        if active_factory.get("profile_path") != "profiles/lumyn.yaml":
            fail(".factory/factoryd.json factory.profile_path must be profiles/lumyn.yaml")
    autoship_repos = autoship_config.get("repos")
    if not isinstance(autoship_repos, dict) or repo_key not in autoship_repos:
        fail(f".factory/factoryd.autoship.example.json must define repos.{repo_key}")
    autoship_lumyn = autoship_repos[repo_key]
    if not isinstance(autoship_lumyn, dict):
        fail(f".factory/factoryd.autoship.example.json repos.{repo_key} must be an object")
    for key, expected in expected_paths.items():
        if autoship_lumyn.get(key) != expected:
            fail(f".factory/factoryd.autoship.example.json repos.{repo_key}.{key} must be {expected!r}")
    validate_factoryd_runtime(autoship_lumyn, f".factory/factoryd.autoship.example.json repos.{repo_key}")
    validate_architecture_budget_policy(root, autoship_lumyn, f".factory/factoryd.autoship.example.json repos.{repo_key}")
    autoship_shipping = autoship_lumyn.get("shipping")
    if not isinstance(autoship_shipping, dict):
        fail(f".factory/factoryd.autoship.example.json repos.{repo_key} must declare shipping block")
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
    if ".factoryd/" not in (root / ".gitignore").read_text(encoding="utf-8"):
        fail(".gitignore must ignore .factoryd/")
