#!/usr/bin/env python3
from __future__ import annotations

import json
from copy import deepcopy
import subprocess
import sys
from typing import Any


OWNER = "Clyra-AI"
REPO = "lumyn"
BRANCH = "main"
RULESET_ID = "17433138"
REQUIRED_CONTEXTS = {"validate", "CodeQL analyze"}
FORBIDDEN_DEFAULT_BRANCH_EXCLUDES = {"~DEFAULT_BRANCH", f"refs/heads/{BRANCH}", BRANCH}


def gh_api(path: str) -> dict[str, Any]:
    completed = subprocess.run(
        ["gh", "api", path],
        check=False,
        capture_output=True,
        text=True,
    )
    if completed.returncode != 0:
        raise RuntimeError(completed.stderr.strip() or completed.stdout.strip())
    return json.loads(completed.stdout)


def require(condition: bool, message: str) -> None:
    if not condition:
        raise AssertionError(message)


def enabled(payload: dict[str, Any], key: str) -> bool:
    value = payload.get(key)
    return isinstance(value, dict) and value.get("enabled") is True


def disabled(payload: dict[str, Any], key: str) -> bool:
    value = payload.get(key)
    return isinstance(value, dict) and value.get("enabled") is False


def rule_types(ruleset: dict[str, Any]) -> set[str]:
    rules = ruleset.get("rules")
    require(isinstance(rules, list), "ruleset.rules must be a list")
    return {str(rule.get("type")) for rule in rules if isinstance(rule, dict)}


def status_contexts_from_ruleset(ruleset: dict[str, Any]) -> set[str]:
    contexts: set[str] = set()
    for rule in ruleset.get("rules", []):
        if not isinstance(rule, dict) or rule.get("type") != "required_status_checks":
            continue
        params = rule.get("parameters")
        require(isinstance(params, dict), "required_status_checks.parameters must be an object")
        require(params.get("strict_required_status_checks_policy") is True, "ruleset status checks must be strict")
        for check in params.get("required_status_checks", []):
            if isinstance(check, dict) and isinstance(check.get("context"), str):
                contexts.add(check["context"])
    return contexts


def validate_branch_protection(protection: dict[str, Any], ruleset: dict[str, Any]) -> dict[str, Any]:
    status_checks = protection.get("required_status_checks")
    require(isinstance(status_checks, dict), "branch protection must require status checks")
    require(status_checks.get("strict") is True, "branch protection status checks must be strict")
    contexts = set(status_checks.get("contexts", []))
    require(REQUIRED_CONTEXTS <= contexts, f"missing required status contexts: {sorted(REQUIRED_CONTEXTS - contexts)}")

    require(enabled(protection, "enforce_admins"), "branch protection must apply to admins")
    require(disabled(protection, "allow_force_pushes"), "force pushes must be disabled")
    require(disabled(protection, "allow_deletions"), "branch deletion must be disabled")
    require(enabled(protection, "required_conversation_resolution"), "conversation resolution must be required")

    reviews = protection.get("required_pull_request_reviews")
    require(isinstance(reviews, dict), "branch protection must declare pull request review controls")
    require(reviews.get("dismiss_stale_reviews") is True, "stale reviews must be dismissed on push")
    require(reviews.get("required_approving_review_count") == 0, "human approval count must remain explicitly zero for bootstrap")

    require(ruleset.get("name") == "protect-main-from-direct-push", "unexpected ruleset name")
    require(ruleset.get("enforcement") == "active", "ruleset must be active")
    require(ruleset.get("target") == "branch", "ruleset target must be branch")
    require(ruleset.get("current_user_can_bypass") == "never", "ruleset must not allow current-user bypass")

    ref_condition = (ruleset.get("conditions") or {}).get("ref_name") or {}
    ref_includes = ref_condition.get("include", [])
    ref_excludes = ref_condition.get("exclude", [])
    require(isinstance(ref_includes, list), "ruleset ref_name.include must be a list")
    require(isinstance(ref_excludes, list), "ruleset ref_name.exclude must be a list")
    require("~DEFAULT_BRANCH" in ref_includes, "ruleset must include the default branch")
    excluded_refs = {ref for ref in ref_excludes if isinstance(ref, str)}
    forbidden_excludes = sorted(FORBIDDEN_DEFAULT_BRANCH_EXCLUDES.intersection(excluded_refs))
    require(not forbidden_excludes, f"ruleset must not exclude default branch refs: {forbidden_excludes}")

    types = rule_types(ruleset)
    for required_rule in {"deletion", "non_fast_forward", "pull_request", "required_status_checks"}:
        require(required_rule in types, f"ruleset missing {required_rule} rule")
    require(REQUIRED_CONTEXTS <= status_contexts_from_ruleset(ruleset), "ruleset missing required status contexts")

    return {
        "required_status_contexts": sorted(REQUIRED_CONTEXTS),
        "ruleset_include_refs": ref_includes,
        "ruleset_exclude_refs": ref_excludes,
        "default_branch_excluded": False,
    }


def fixture_protection() -> dict[str, Any]:
    return {
        "required_status_checks": {
            "strict": True,
            "contexts": sorted(REQUIRED_CONTEXTS),
        },
        "enforce_admins": {"enabled": True},
        "allow_force_pushes": {"enabled": False},
        "allow_deletions": {"enabled": False},
        "required_conversation_resolution": {"enabled": True},
        "required_pull_request_reviews": {
            "dismiss_stale_reviews": True,
            "required_approving_review_count": 0,
        },
    }


def fixture_ruleset() -> dict[str, Any]:
    return {
        "name": "protect-main-from-direct-push",
        "enforcement": "active",
        "target": "branch",
        "current_user_can_bypass": "never",
        "conditions": {
            "ref_name": {
                "include": ["~DEFAULT_BRANCH"],
                "exclude": [],
            },
        },
        "rules": [
            {"type": "deletion"},
            {"type": "non_fast_forward"},
            {"type": "pull_request"},
            {
                "type": "required_status_checks",
                "parameters": {
                    "strict_required_status_checks_policy": True,
                    "required_status_checks": [{"context": context} for context in sorted(REQUIRED_CONTEXTS)],
                },
            },
        ],
    }


def run_self_test() -> int:
    try:
        validate_branch_protection(fixture_protection(), fixture_ruleset())

        for excluded_ref in sorted(FORBIDDEN_DEFAULT_BRANCH_EXCLUDES):
            protection = fixture_protection()
            ruleset = deepcopy(fixture_ruleset())
            ruleset["conditions"]["ref_name"]["exclude"] = [excluded_ref]
            try:
                validate_branch_protection(protection, ruleset)
            except AssertionError as exc:
                require("must not exclude default branch refs" in str(exc), f"unexpected self-test failure: {exc}")
                continue
            raise AssertionError(f"default branch exclude {excluded_ref!r} was accepted")
    except Exception as exc:
        print(f"branch protection audit self-test failed: {exc}", file=sys.stderr)
        return 1

    print(json.dumps({"status": "pass", "self_test": "branch_protection_audit"}, sort_keys=True))
    return 0


def main() -> int:
    if len(sys.argv) > 1:
        if sys.argv[1] == "--self-test":
            return run_self_test()
        print("usage: audit_branch_protection.py [--self-test]", file=sys.stderr)
        return 2

    try:
        protection = gh_api(f"repos/{OWNER}/{REPO}/branches/{BRANCH}/protection")
        ruleset = gh_api(f"repos/{OWNER}/{REPO}/rulesets/{RULESET_ID}")
        summary = validate_branch_protection(protection, ruleset)
    except Exception as exc:
        print(f"branch protection audit failed: {exc}", file=sys.stderr)
        return 1

    print(
        json.dumps(
            {
                "status": "pass",
                "repo": f"{OWNER}/{REPO}",
                "branch": BRANCH,
                "ruleset_id": RULESET_ID,
                "required_status_contexts": summary["required_status_contexts"],
                "ruleset_include_refs": summary["ruleset_include_refs"],
                "ruleset_exclude_refs": summary["ruleset_exclude_refs"],
                "default_branch_excluded": summary["default_branch_excluded"],
                "direct_push_protection": "active",
            },
            sort_keys=True,
        )
    )
    return 0


if __name__ == "__main__":
    sys.exit(main())
