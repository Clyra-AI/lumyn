#!/usr/bin/env python3
from __future__ import annotations

import json
import subprocess
import sys
from typing import Any


OWNER = "Clyra-AI"
REPO = "lumyn"
BRANCH = "main"
RULESET_ID = "17433138"
REQUIRED_CONTEXTS = {"validate", "CodeQL analyze"}


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


def main() -> int:
    try:
        protection = gh_api(f"repos/{OWNER}/{REPO}/branches/{BRANCH}/protection")
        ruleset = gh_api(f"repos/{OWNER}/{REPO}/rulesets/{RULESET_ID}")

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
        ref_names = ((ruleset.get("conditions") or {}).get("ref_name") or {}).get("include", [])
        require("~DEFAULT_BRANCH" in ref_names, "ruleset must include the default branch")
        types = rule_types(ruleset)
        for required_rule in {"deletion", "non_fast_forward", "pull_request", "required_status_checks"}:
            require(required_rule in types, f"ruleset missing {required_rule} rule")
        require(REQUIRED_CONTEXTS <= status_contexts_from_ruleset(ruleset), "ruleset missing required status contexts")
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
                "required_status_contexts": sorted(REQUIRED_CONTEXTS),
                "direct_push_protection": "active",
            },
            sort_keys=True,
        )
    )
    return 0


if __name__ == "__main__":
    sys.exit(main())
