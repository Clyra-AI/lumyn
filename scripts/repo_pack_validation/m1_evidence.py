"""Fail-closed validation for attended M1 implementation evidence."""

from __future__ import annotations

import copy
import hashlib
import json
import os
import re
import stat
import subprocess
import tempfile
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Callable

from . import factory_schema_core


EVIDENCE_ROOT = ".factory/artifacts/task-runs/M1-IMPLEMENTATION"
M1_BASE_GIT_SHA = "d7ae311c391775a2517d56add6d57148d5891ef3"
M1_WORK_ITEM_ID = "lumyn-v3-m1-attended-implementation"
REPORT_REF = f"{EVIDENCE_ROOT}/validation-report.json"
SCORECARD_REF = f"{EVIDENCE_ROOT}/proof-of-behavior-scorecard.json"
RED_MARKER_REF = f"{EVIDENCE_ROOT}/red-first/work-proof-marker.json"
HOLDOUT_RESULT_REF = ".factory/artifacts/lifecycle-evidence/M1/holdout-result.json"
REVIEW_REPORT_REF = ".factory/artifacts/lifecycle-evidence/M1/review-report.json"
HOLDOUT_SCHEMA_REF = ".factory/contracts/factory/holdout-result.schema.json"
HOLDOUT_SCHEMA_DIGEST = "sha256:436c0fb514ea904f2b9c0f304f66bc5b50970a520844b14e9ace03e3213cf4bd"
REVIEW_SCHEMA_REF = ".factory/contracts/factory/review-report.schema.json"
REVIEW_SCHEMA_DIGEST = "sha256:7228314ce33630338e0028e0f4b7df2166972f8837aaa832b96b395f4c4dfbf2"
M1_REQUIRED_WORKER_CHAIN = [
    "task-executor",
    "validation-gate",
    "code-review",
    "holdout-evaluator",
    "commit-push",
    "post-merge-monitor",
]
M1_REQUIRED_LIFECYCLE_EVIDENCE = {
    "ship_packet",
    "pr_lifecycle_report",
    "post_merge_report",
    "scope_closure_report",
    "review_report",
    "holdout_result",
}
RED_COMMAND = (
    "LUMYN_M1_VERIFICATION_TARGET=baseline npm --prefix "
    "examples/consumer-repos/det-operation-rename test --silent"
)
RED_CWD = "."
RED_WORKSPACE_REF = "baseline:complete-harness-pre-fix-candidate"
RED_FAILURE_REASON_CODE = "target_contract_mismatch"
RED_PROOF_LINE = b"det-operation-rename: baseline target contract rejected\n"
RED_SUMMARY = (
    "The complete representative fixture and verifier rejected the pre-fix baseline "
    "with the typed det-operation-rename target-contract mismatch; missing, parse, "
    "link, and top-level failures are not accepted as red evidence."
)
COMMANDS = [
    "go test ./tests -run '^TestM1' -count=1",
    "make lint-fast",
    "make test-fast",
    "make test-coverage",
    "python3 scripts/validate_repo_pack.py",
    "make prepush-full",
]
MARKER_REFS = [
    f"{EVIDENCE_ROOT}/validation-{index:03d}/work-proof-marker.json"
    for index in range(1, len(COMMANDS) + 1)
]
CANDIDATE_EXCLUDED_ROOTS = [
    f"{EVIDENCE_ROOT}/",
    ".factory/artifacts/lifecycle-evidence/M1/",
    ".factory/artifacts/pr-lifecycle/lumyn-v3-m1/",
]
PREREQUISITE_ROOTS = [
    "AGENTS.md",
    "WORKFLOW.md",
    "docs/product/prd.md",
    "docs/product/plan.md",
    ".factory/artifacts/prd-to-plan/lumyn-migration-mvp/",
]


def _require(condition: bool, message: str) -> None:
    if not condition:
        raise AssertionError(message)


def _sha256(payload: bytes) -> str:
    return f"sha256:{hashlib.sha256(payload).hexdigest()}"


def _git_output(root: Path, *arguments: str) -> str:
    return subprocess.check_output(
        ["git", *arguments], cwd=root, text=True
    ).strip()


def _candidate_paths(root: Path, base_git_sha: str) -> list[str]:
    tracked = subprocess.run(
        ["git", "diff", "--name-only", "--relative", "--no-renames", base_git_sha],
        cwd=root,
        text=True,
        capture_output=True,
        check=True,
    ).stdout
    untracked = subprocess.run(
        ["git", "ls-files", "--others", "--exclude-standard"],
        cwd=root,
        text=True,
        capture_output=True,
        check=True,
    ).stdout
    paths = {
        path.strip().replace("\\", "/")
        for path in [*tracked.splitlines(), *untracked.splitlines()]
        if path.strip()
    }
    return sorted(
        path
        for path in paths
        if not any(
            path == excluded.rstrip("/") or path.startswith(excluded)
            for excluded in CANDIDATE_EXCLUDED_ROOTS
        )
    )


def _candidate_entry(root: Path, relative: str) -> dict[str, object]:
    path = Path(relative)
    _require(
        relative and not path.is_absolute() and ".." not in path.parts,
        f"M1 candidate path escapes repository: {relative}",
    )
    source = root / path
    if source.is_symlink():
        return {
            "path": relative,
            "kind": "symlink",
            "mode": "120000",
            "digest": _sha256(os.readlink(source).encode()),
        }
    if source.is_file():
        executable = bool(source.stat().st_mode & stat.S_IXUSR)
        return {
            "path": relative,
            "kind": "file",
            "mode": "100755" if executable else "100644",
            "digest": _sha256(source.read_bytes()),
        }
    _require(not source.exists(), f"unsupported M1 candidate path: {relative}")
    return {"path": relative, "kind": "deleted"}


def candidate_binding(
    root: Path,
    declared_paths: list[str] | None = None,
    *,
    base_git_sha: str = M1_BASE_GIT_SHA,
) -> dict[str, Any]:
    root = root.resolve()
    try:
        _git_output(root, "cat-file", "-e", f"{base_git_sha}^{{commit}}")
    except subprocess.CalledProcessError as exc:
        raise AssertionError(
            f"M1 immutable base commit is unavailable: {base_git_sha}; CI must fetch full history"
        ) from exc
    head = _git_output(root, "rev-parse", "HEAD")
    _require(
        _git_output(root, "merge-base", base_git_sha, head) == base_git_sha,
        "M1 immutable base is not an ancestor of the current checkout",
    )
    discovered_paths = _candidate_paths(root, base_git_sha)
    if declared_paths is None:
        paths = discovered_paths
    else:
        _require(
            declared_paths == sorted(set(declared_paths))
            and all(isinstance(path, str) and path for path in declared_paths),
            "M1 declared candidate paths must be sorted, unique, nonempty strings",
        )
        missing = sorted(set(declared_paths) - set(discovered_paths))
        unexpected = sorted(set(discovered_paths) - set(declared_paths))
        _require(
            not missing and not unexpected,
            f"M1 candidate path set drifted from the immutable base; missing={missing}, unexpected={unexpected}",
        )
        paths = declared_paths
    entries = [_candidate_entry(root, path) for path in paths]
    payload = json.dumps(
        {"base_git_sha": base_git_sha, "entries": entries},
        sort_keys=True,
        separators=(",", ":"),
    ).encode()
    return {
        "base_git_sha": base_git_sha,
        "candidate_digest": _sha256(payload),
        "changed_paths": paths,
        "entries": entries,
    }


def _current_candidate_binding(root: Path, declared_paths: list[str]) -> dict[str, Any]:
    return candidate_binding(root, declared_paths)


def validate_evidence_ownership(
    packet: dict[str, Any], require: Callable[[bool, str], None]
) -> None:
    ownership = packet.get("evidence_ownership", {})
    executor = ownership.get("task-executor", {})
    validation = ownership.get("validation-gate", {})
    require(
        executor.get("allowed_paths_ref") == "#/allowed_paths"
        and executor.get("forbidden_paths") == [f"{EVIDENCE_ROOT}/"]
        and validation.get("allowed_paths") == [f"{EVIDENCE_ROOT}/"]
        and validation.get("forbidden_paths_ref") == "#/allowed_paths",
        "M1 executor and validation-gate evidence ownership must remain disjoint",
    )


def _load_json(path: Path) -> tuple[dict[str, Any], bytes]:
    _require(path.is_file(), f"missing M1 evidence artifact: {path}")
    raw = path.read_bytes()
    try:
        value = json.loads(raw)
    except Exception as exc:
        raise AssertionError(f"invalid M1 evidence JSON {path}: {exc}") from exc
    _require(isinstance(value, dict), f"M1 evidence must be an object: {path}")
    return value, raw


def _validate_vendored_schema(
    root: Path,
    schema_ref: str,
    expected_digest: str,
    payload: dict[str, Any],
    label: str,
) -> None:
    schema, raw = _load_json(root / schema_ref)
    _require(
        _sha256(raw) == expected_digest,
        f"{label} vendored Factory schema digest drifted",
    )

    def fail(message: str) -> None:
        raise AssertionError(message)

    def load_json(path: Path) -> Any:
        return json.loads(path.read_text(encoding="utf-8"))

    factory_schema_core.validate_schema(
        schema,
        payload,
        label,
        root=root,
        fail=fail,
        load_json=load_json,
        validation_error_type=AssertionError,
    )


def load_evidence(root: Path) -> tuple[dict[str, dict[str, Any]], dict[str, bytes]]:
    payloads: dict[str, dict[str, Any]] = {}
    artifacts: dict[str, bytes] = {}
    for ref in [
        REPORT_REF,
        SCORECARD_REF,
        RED_MARKER_REF,
        *MARKER_REFS,
        HOLDOUT_RESULT_REF,
        REVIEW_REPORT_REF,
    ]:
        payloads[ref], artifacts[ref] = _load_json(root / ref)
    for marker_ref in [RED_MARKER_REF, *MARKER_REFS]:
        parent = marker_ref.rsplit("/", 1)[0]
        for suffix in ("stdout", "stderr"):
            raw_ref = f"{parent}/work-proof-marker.{suffix}.log"
            path = root / raw_ref
            _require(path.is_file(), f"missing M1 raw log: {raw_ref}")
            artifacts[raw_ref] = path.read_bytes()
    return payloads, artifacts


def _timestamp(value: object, label: str) -> datetime:
    _require(isinstance(value, str) and value, f"{label} timestamp is required")
    try:
        parsed = datetime.fromisoformat(str(value).replace("Z", "+00:00"))
    except ValueError as exc:
        raise AssertionError(f"{label} timestamp is invalid: {exc}") from exc
    _require(parsed.tzinfo is not None, f"{label} timestamp requires timezone")
    return parsed.astimezone(timezone.utc)


def _path_allowed(path: str, allowed: list[str]) -> bool:
    return any(
        path.startswith(rule) if rule.endswith("/") else path == rule
        for rule in allowed
    )


def _contains_key(value: object, forbidden: set[str]) -> bool:
    if isinstance(value, dict):
        return any(
            key in forbidden or _contains_key(child, forbidden)
            for key, child in value.items()
        )
    if isinstance(value, list):
        return any(_contains_key(child, forbidden) for child in value)
    return False


def _validate_marker(
    marker: dict[str, Any],
    artifacts: dict[str, bytes],
    *,
    ref: str,
    command: str,
    marker_id: str,
    validation_checkout_sha: str,
    should_pass: bool,
    cwd: str,
    workspace_ref: str,
    validation_run_id: str | None = None,
    candidate_digest: str | None = None,
) -> tuple[datetime, datetime]:
    _require(marker.get("command") == command, f"M1 marker command drifted: {ref}")
    _require(marker.get("marker_id") == marker_id, f"M1 marker identity drifted: {ref}")
    _require(marker.get("generated_by") == "trusted_runner", f"M1 marker is not trusted: {ref}")
    _require(
        marker.get("git_sha") == validation_checkout_sha,
        f"M1 marker validation checkout SHA drifted: {ref}",
    )
    _require(marker.get("cwd") == cwd, f"M1 marker cwd drifted: {ref}")
    _require(marker.get("workspace_ref") == workspace_ref, f"M1 marker workspace drifted: {ref}")
    _require(marker.get("runner_id") == "trusted-runner:factory-reference-validation", f"M1 marker runner drifted: {ref}")
    _require("workspace_proof" not in marker, f"M1 marker overclaims isolated workspace: {ref}")
    if validation_run_id is not None or candidate_digest is not None:
        _require(
            marker.get("validation_run_id") == validation_run_id
            and marker.get("candidate_digest") == candidate_digest,
            f"M1 marker current-work binding drifted: {ref}",
        )
    expected_status = "pass" if should_pass else "fail"
    _require(marker.get("execution_status") == expected_status, f"M1 marker status drifted: {ref}")
    exit_code = marker.get("exit_code")
    _require(type(exit_code) is int, f"M1 marker exit code missing: {ref}")
    _require((exit_code == 0) is should_pass, f"M1 marker exit code contradicts status: {ref}")
    raw_refs = marker.get("raw_log_refs")
    parent = ref.rsplit("/", 1)[0]
    expected_raw = [
        f"{parent}/work-proof-marker.stdout.log",
        f"{parent}/work-proof-marker.stderr.log",
    ]
    _require(raw_refs == expected_raw, f"M1 marker raw-log refs drifted: {ref}")
    _require(
        marker.get("stdout_digest") == _sha256(artifacts[expected_raw[0]])
        and marker.get("stderr_digest") == _sha256(artifacts[expected_raw[1]]),
        f"M1 marker raw-log digest drifted: {ref}",
    )
    started = _timestamp(marker.get("started_at"), f"{ref}.started_at")
    finished = _timestamp(marker.get("finished_at"), f"{ref}.finished_at")
    _require(finished >= started, f"M1 marker timestamps are reversed: {ref}")
    return started, finished


def _validate_red_marker(
    marker: dict[str, Any],
    artifacts: dict[str, bytes],
    *,
    validation_checkout_sha: str,
    validation_run_id: str,
    candidate_digest: str,
) -> tuple[datetime, datetime]:
    started, finished = _validate_marker(
        marker,
        artifacts,
        ref=RED_MARKER_REF,
        command=RED_COMMAND,
        marker_id="red-first",
        validation_checkout_sha=validation_checkout_sha,
        should_pass=False,
        cwd=RED_CWD,
        workspace_ref=RED_WORKSPACE_REF,
        validation_run_id=validation_run_id,
        candidate_digest=candidate_digest,
    )
    stdout_ref, stderr_ref = marker["raw_log_refs"]
    _require(marker.get("exit_code") == 1, "M1 red-first exit code is not the typed target mismatch")
    _require(
        marker.get("failure_reason_code") == RED_FAILURE_REASON_CODE,
        "M1 red-first failure reason is not the typed target mismatch",
    )
    _require(
        artifacts[stdout_ref] == RED_PROOF_LINE and artifacts[stderr_ref] == b"",
        "M1 red-first output is not the exact typed target-contract mismatch",
    )
    return started, finished


def _validate_lifecycle_reports(
    root: Path,
    lifecycle_task: dict[str, Any],
    payloads: dict[str, dict[str, Any]],
    artifacts: dict[str, bytes],
    *,
    validation_run_id: str,
    candidate_digest: str,
    validation_generated_at: datetime,
) -> None:
    gates = lifecycle_task.get("lifecycle_gates", {})
    declared_refs = lifecycle_task.get("lifecycle_evidence_refs", {})
    _require(
        lifecycle_task.get("task_id") == "M1"
        and lifecycle_task.get("work_item_id") == "lumyn-v3-m1"
        and lifecycle_task.get("required_worker_chain") == M1_REQUIRED_WORKER_CHAIN
        and set(lifecycle_task.get("lifecycle_evidence_required", []))
        == M1_REQUIRED_LIFECYCLE_EVIDENCE
        and gates.get("local_validation_required") is True
        and gates.get("ci_required") is True
        and gates.get("code_review_required") is True
        and gates.get("codex_review_required") is True
        and gates.get("holdout_provisioning_required") is True
        and gates.get("holdout_evaluation_required") is False
        and gates.get("commit_push_required") is True
        and gates.get("post_merge_monitor_required") is True
        and gates.get("pr_lifecycle_report_required") is True
        and declared_refs.get("holdout_result") == HOLDOUT_RESULT_REF
        and declared_refs.get("review_report") == REVIEW_REPORT_REF,
        "M1 task-required lifecycle report declarations drifted",
    )
    _require(
        HOLDOUT_RESULT_REF in payloads and REVIEW_REPORT_REF in payloads,
        "M1 task-required lifecycle reports are missing",
    )
    marker_bindings = [
        {"ref": ref, "sha256": _sha256(artifacts[ref])}
        for ref in MARKER_REFS
    ]
    holdout = payloads[HOLDOUT_RESULT_REF]
    review = payloads[REVIEW_REPORT_REF]
    validation_report = payloads[REPORT_REF]
    _validate_vendored_schema(
        root,
        HOLDOUT_SCHEMA_REF,
        HOLDOUT_SCHEMA_DIGEST,
        holdout,
        "M1 holdout-result",
    )
    _validate_vendored_schema(
        root,
        REVIEW_SCHEMA_REF,
        REVIEW_SCHEMA_DIGEST,
        review,
        "M1 review-report",
    )

    holdout_policy = lifecycle_task.get("holdout_policy", {})
    canonical_policy = {
        key: value
        for key, value in holdout_policy.items()
        if key != "policy_digest"
    }
    expected_policy_digest = _sha256(
        json.dumps(
            canonical_policy,
            sort_keys=True,
            separators=(",", ":"),
        ).encode()
    )
    _require(
        canonical_policy
        == {
            "mode": "provision",
            "suite_namespace": "private://lumyn/m1/v1",
            "commitment_algorithm": "hmac-sha256",
        }
        and holdout_policy.get("policy_digest") == expected_policy_digest,
        "M1 holdout policy is not the exact digest-bound Factory provision contract",
    )
    holdout_manifest, _ = _load_json(root / "examples/holdout-manifest.json")
    provisioning_contract = lifecycle_task.get("holdout_provisioning_contract", {})
    _require(
        holdout_manifest.get("suite_namespace") == holdout_policy.get("suite_namespace")
        and holdout_manifest.get("commitment_algorithm")
        == holdout_policy.get("commitment_algorithm")
        and holdout_manifest.get("provisioning_result_ref")
        == provisioning_contract.get("result_ref")
        and holdout_manifest.get("committed_fields")
        == provisioning_contract.get("committed_fields")
        and set(holdout_manifest.get("prohibited_fields", []))
        == set(provisioning_contract.get("prohibited_committed_fields", [])),
        "M1 holdout policy and source-safe provisioning manifest drifted",
    )
    implementation_producer = validation_report.get("implementation_producer", {})
    executor_id = implementation_producer.get("producer_id")
    _require(
        implementation_producer.get("worker") == "task-executor"
        and implementation_producer.get("producer_class") == "agent"
        and isinstance(executor_id, str)
        and bool(executor_id)
        and executor_id != "trusted-runner:factory-reference-validation",
        "M1 validation report omits the bound task-executor identity",
    )
    holdout_producer = holdout.get("producer", {})
    evaluator_id = holdout.get("evaluator_id")
    evaluator_class = holdout.get("evaluator_class")
    _require(
        holdout.get("artifact_type") == "holdout_result"
        and holdout.get("schema_version") == "0.1"
        and holdout.get("task_id") == "M1"
        and holdout.get("work_item_id") == "lumyn-v3-m1"
        and holdout.get("policy_mode") == holdout_policy.get("mode") == "provision"
        and isinstance(holdout.get("run_id"), str)
        and holdout["run_id"].startswith("holdout-provision:")
        and isinstance(evaluator_id, str)
        and bool(evaluator_id)
        and evaluator_class in {"independent", "human"}
        and evaluator_id not in {"trusted-runner:factory-reference-validation", executor_id}
        and holdout_producer.get("worker") == "holdout-evaluator"
        and holdout_producer.get("producer_class") == evaluator_class
        and holdout_producer.get("producer_id") == evaluator_id,
        "M1 holdout provisioning identity or independence drifted",
    )
    _require(
        holdout.get("current_work")
        == {
            "validation_run_id": validation_run_id,
            "candidate_digest": candidate_digest,
            "work_proof_markers": marker_bindings,
        }
        and holdout.get("work_proof_marker_refs") == MARKER_REFS,
        "M1 holdout provisioning current-work binding drifted",
    )
    _require(
        holdout.get("task_policy_digest") == expected_policy_digest
        and holdout.get("suite_ref") == holdout_policy.get("suite_namespace")
        and isinstance(holdout.get("suite_commitment"), str)
        and holdout["suite_commitment"].startswith(
            holdout_policy.get("commitment_algorithm") + ":"
        )
        and len(holdout["suite_commitment"])
        == len(holdout_policy.get("commitment_algorithm") + ":") + 64
        and all(
            character in "0123456789abcdef"
            for character in holdout["suite_commitment"][
                len(holdout_policy.get("commitment_algorithm") + ":") :
            ]
        )
        and type(holdout.get("cases_run")) is int
        and holdout["cases_run"] >= 1
        and holdout.get("failing_cases") == []
        and holdout.get("promotion_decision") == "pass",
        "M1 holdout provisioning result or policy binding drifted",
    )
    prohibited = set(provisioning_contract.get("prohibited_committed_fields", []))
    _require(
        prohibited and not _contains_key(holdout, prohibited),
        "M1 holdout provisioning result contains prohibited private fields",
    )
    holdout_created = _timestamp(holdout.get("created_at"), "M1 holdout created_at")
    _require(
        holdout_created >= validation_generated_at,
        "M1 holdout provisioning predates the bound validation run",
    )

    review_producer = review.get("producer", {})
    reviewer_id = review.get("reviewer_id")
    reviewer_class = review.get("reviewer_class")
    expected_producer_class = {
        "peer_agent": "peer_agent",
        "security_reviewer": "security_reviewer",
        "human_reviewer": "human_reviewer",
    }.get(reviewer_class)
    _require(
        review.get("artifact_type") == "review_report"
        and review.get("schema_version") == "0.1"
        and review.get("task_id") == "M1"
        and review.get("work_item_id") == "lumyn-v3-m1"
        and review.get("review_type") == "code"
        and review.get("risk_class") == lifecycle_task.get("risk_class") == "medium"
        and isinstance(reviewer_id, str)
        and bool(reviewer_id)
        and expected_producer_class is not None
        and reviewer_id
        not in {"trusted-runner:factory-reference-validation", executor_id, evaluator_id}
        and review_producer.get("worker") == "code-review"
        and review_producer.get("producer_id") == reviewer_id
        and review_producer.get("producer_class") == expected_producer_class,
        "M1 review identity or independence drifted",
    )
    _require(
        review.get("current_work")
        == {
            "validation_run_id": validation_run_id,
            "candidate_digest": candidate_digest,
            "work_proof_markers": marker_bindings,
        }
        and review.get("work_proof_marker_refs") == MARKER_REFS,
        "M1 review current-work binding drifted",
    )
    findings = review.get("findings")
    approval = review.get("approval_effect", {})
    review_scope = review.get("review_scope", [])
    scope_paths = [
        item
        for item in review_scope
        if isinstance(item, str) and not item.startswith("git diff ")
    ]
    _require(
        isinstance(review_scope, list)
        and f"git diff {M1_BASE_GIT_SHA}" in review_scope
        and all(
            _path_allowed(path, scope_paths)
            for path in validation_report.get("changed_paths", [])
        )
        and isinstance(findings, list)
        and not any(item.get("status") == "open" for item in findings if isinstance(item, dict))
        and not any(
            item.get("severity") in {"P0", "P1"}
            and (
                item.get("status") != "resolved"
                or item.get("verification_status") != "independently_verified"
                or item.get("claim_source") != "independently_observed"
            )
            for item in findings
            if isinstance(item, dict)
        )
        and review.get("verdict") == "approved"
        and review.get("required_fixes") == []
        and approval.get("promotion_decision") == "ready_for_pr"
        and approval.get("approvals_granted")
        == [
            "code_review_passed",
            "medium_risk_structured_review_passed",
            "independent_candidate_binding_verified",
        ]
        and approval.get("blocks_promotion") is False,
        "M1 independent review does not approve the exact candidate",
    )
    evidence_refs = review.get("evidence_refs", [])
    isolation = review.get("context_isolation", {})
    allowed_sources = set(isolation.get("allowed_sources", []))
    disallowed_sources = set(isolation.get("disallowed_sources", []))
    _require(
        REPORT_REF in evidence_refs
        and HOLDOUT_RESULT_REF not in evidence_refs
        and all(ref in evidence_refs for ref in MARKER_REFS)
        and not {
            "private_builder_notes",
            "unpromoted_chat_reasoning",
            "private_holdout_inputs",
            "private_holdout_answers",
        }.intersection(allowed_sources)
        and {
            "private_builder_notes",
            "unpromoted_chat_reasoning",
            "private_holdout_inputs",
            "private_holdout_answers",
        }.issubset(disallowed_sources)
        and isolation.get("builder_notes_authoritative") is False,
        "M1 review evidence refs or context isolation drifted",
    )
    systems_questions = review.get("systems_questions", {})
    _require(
        all(
            isinstance(systems_questions.get(key), list)
            and bool(systems_questions[key])
            and all(
                isinstance(item, str) and bool(item.strip())
                for item in systems_questions[key]
            )
            for key in (
                "state_impact",
                "source_of_truth_changes",
                "feedback_plan",
                "blast_radius",
                "rollback_or_deletion_test",
            )
        )
        and isinstance(systems_questions.get("system_model_summary"), str)
        and bool(systems_questions["system_model_summary"].strip()),
        "M1 review systems analysis is empty or malformed",
    )
    residuals = " ".join(str(item).lower() for item in review.get("residual_risks", []))
    _require(
        "transition or archive" in residuals
        and "before unrelated milestone" in residuals,
        "M1 review omits the required post-main historical transition residual",
    )
    review_created = _timestamp(review.get("created_at"), "M1 review created_at")
    _require(
        review_created >= validation_generated_at,
        "M1 review predates the bound validation run",
    )
    _require(
        holdout_created >= review_created,
        "M1 holdout provisioning predates canonical code-review completion",
    )


def validate_evidence_bundle(
    root: Path,
    packet: dict[str, Any],
    lifecycle_task: dict[str, Any],
    payloads: dict[str, dict[str, Any]],
    artifacts: dict[str, bytes],
    binding: dict[str, Any],
) -> None:
    report = payloads[REPORT_REF]
    scorecard = payloads[SCORECARD_REF]
    base_git_sha = str(binding["base_git_sha"])
    digest = str(binding["candidate_digest"])
    validation_run_id = report.get("validation_run_id")
    validation_checkout_sha = report.get("validation_checkout_sha")
    changed_paths = binding["changed_paths"]
    declared_binding = report.get("candidate_binding", {})
    _require(
        report.get("artifact_type") == "validation_report"
        and report.get("schema_version") == "1.0"
        and report.get("worker") == "validation-gate"
        and report.get("task_id") == "M1-IMPLEMENTATION"
        and report.get("work_item_id") == M1_WORK_ITEM_ID,
        "M1 validation report identity drifted",
    )
    run_prefix = "validation:"
    run_suffix = ":lumyn-m1-implementation-r2"
    _require(
        isinstance(validation_run_id, str)
        and validation_run_id.startswith(run_prefix)
        and validation_run_id.endswith(run_suffix)
        and len(validation_run_id) > len(run_prefix) + len(run_suffix),
        "M1 validation run identity is missing or malformed",
    )
    _timestamp(
        validation_run_id[len(run_prefix) : -len(run_suffix)],
        "M1 validation run identity",
    )
    _require(
        declared_binding.get("base_git_sha") == base_git_sha
        and declared_binding.get("candidate_digest") == digest
        and declared_binding.get("excluded_roots") == CANDIDATE_EXCLUDED_ROOTS,
        "M1 validation report candidate binding drifted",
    )
    _require(
        isinstance(validation_checkout_sha, str)
        and re.fullmatch(r"[0-9a-f]{40}", validation_checkout_sha) is not None,
        "M1 validation checkout SHA is missing or malformed",
    )
    _require(report.get("changed_paths") == changed_paths, "M1 validation report changed paths drifted")
    allowed = packet.get("allowed_paths", [])
    planning_prerequisites = packet.get("planning_prerequisite_paths", [])
    _require(
        isinstance(allowed, list)
        and planning_prerequisites == PREREQUISITE_ROOTS
        and all(
            _path_allowed(path, allowed) or _path_allowed(path, planning_prerequisites)
            for path in changed_paths
        ),
        "M1 release candidate path falls outside implementation or prerequisite scope",
    )
    context = report.get("changed_path_context", {})
    _require(
        context.get("allowed_paths") == allowed
        and context.get("forbidden_paths") == packet.get("forbidden_paths")
        and context.get("path_violations") == []
        and context.get("prerequisite_control_repairs_included_in_release_candidate") == planning_prerequisites,
        "M1 validation report path context drifted",
    )

    _require(report.get("work_proof_marker_refs") == MARKER_REFS, "M1 work-proof marker refs drifted")
    marker_bindings = report.get("work_proof_markers")
    _require(
        marker_bindings == [
            {"ref": ref, "sha256": _sha256(artifacts[ref])}
            for ref in MARKER_REFS
        ],
        "M1 work-proof marker file binding drifted",
    )
    rows = report.get("validation_commands")
    checks = report.get("checks")
    _require(isinstance(rows, list) and len(rows) == len(COMMANDS), "M1 validation command count drifted")
    _require(
        checks == [{"name": command, "status": "pass"} for command in COMMANDS],
        "M1 validation check list drifted",
    )
    red = payloads[RED_MARKER_REF]
    _, red_finished = _validate_red_marker(
        red,
        artifacts,
        validation_checkout_sha=validation_checkout_sha,
        validation_run_id=str(validation_run_id),
        candidate_digest=digest,
    )

    latest = red_finished
    for index, (command, ref) in enumerate(zip(COMMANDS, MARKER_REFS)):
        marker = payloads[ref]
        started, finished = _validate_marker(
            marker,
            artifacts,
            ref=ref,
            command=command,
            marker_id=f"validation-{index + 1:03d}",
            validation_checkout_sha=validation_checkout_sha,
            should_pass=True,
            cwd=".",
            workspace_ref="branch:codex/lumyn-m1-runner-contract",
            validation_run_id=str(validation_run_id),
            candidate_digest=digest,
        )
        _require(
            started >= latest,
            f"M1 green marker order overlaps or replays prior evidence: {ref}",
        )
        latest = max(latest, finished)
        row = rows[index]
        _require(
            row.get("command") == command
            and row.get("exit_code") == 0
            and row.get("started_at") == marker.get("started_at")
            and row.get("completed_at") == marker.get("finished_at")
            and row.get("output_ref") == marker.get("raw_log_refs")[0]
            and row.get("output_sha256") == marker.get("stdout_digest")
            and row.get("runner_identity") == marker.get("runner_id")
            and row.get("repo_ref") == validation_checkout_sha,
            f"M1 validation report row drifted: {command}",
        )
        _require(started <= finished, f"M1 validation command timestamps drifted: {command}")
    report_generated_at = _timestamp(report.get("generated_at"), "M1 report generated_at")
    _require(report_generated_at >= latest, "M1 report predates green evidence")

    red_evidence = report.get("red_evidence")
    _require(
        red_evidence
        == [
            {
                "check": RED_COMMAND,
                "summary": RED_SUMMARY,
                "evidence_refs": [
                    RED_MARKER_REF,
                    f"{EVIDENCE_ROOT}/red-first/work-proof-marker.stdout.log",
                    f"{EVIDENCE_ROOT}/red-first/work-proof-marker.stderr.log",
                    "examples/verification/run.mjs",
                    "examples/consumer-repos/det-operation-rename/package.json",
                ],
            }
        ],
        "M1 red-first evidence is not trusted-marker bound",
    )

    _require(
        scorecard.get("task_id") == "M1-IMPLEMENTATION"
        and scorecard.get("work_item_id") == M1_WORK_ITEM_ID
        and scorecard.get("candidate_digest") == digest
        and scorecard.get("overall_status") == "pass"
        and scorecard.get("required_proof_level") == "workflow_behavior",
        "M1 proof scorecard candidate or status drifted",
    )
    _require(
        not _contains_key(scorecard, {"acceptance_item_ids", "acceptance_results", "acceptance_result_requirements"})
        and scorecard.get("coverage_boundary", {}).get("parent_acceptance_closure_claimed") is False,
        "M1 proof scorecard claims parent acceptance closure",
    )
    authorization_result = report.get("authorization_result", {})
    _require(
        report.get("result") == "pass"
        and report.get("promotion_decision") == "blocked"
        and authorization_result
        == {
            "sovereignty_mode": "local_only",
            "network_used": False,
            "credentials_used": False,
            "live_agent_or_model_used": False,
            "external_writes_performed": False,
            "github_actions_performed": False,
            "acceptance_closure_claimed": False,
        },
        "M1 evidence widens authority or lifecycle promotion",
    )
    _validate_lifecycle_reports(
        root,
        lifecycle_task,
        payloads,
        artifacts,
        validation_run_id=str(validation_run_id),
        candidate_digest=digest,
        validation_generated_at=report_generated_at,
    )


def validate_candidate_scope(root: Path, packet: dict[str, Any]) -> dict[str, Any]:
    binding = candidate_binding(root)
    allowed = packet.get("allowed_paths", [])
    planning_prerequisites = packet.get("planning_prerequisite_paths", [])
    _require(isinstance(allowed, list), "M1 candidate allowed paths are required")
    _require(
        planning_prerequisites == PREREQUISITE_ROOTS,
        "M1 candidate planning prerequisites drifted",
    )
    _require(
        all(
            _path_allowed(path, allowed)
            or _path_allowed(path, planning_prerequisites)
            for path in binding["changed_paths"]
        ),
        "M1 current candidate falls outside implementation or planning-prerequisite scope",
    )
    return binding


def validate_m1_evidence(
    root: Path,
    packet: dict[str, Any],
    lifecycle_task: dict[str, Any],
) -> None:
    payloads, artifacts = load_evidence(root)
    declared_paths = payloads[REPORT_REF].get("changed_paths")
    _require(isinstance(declared_paths, list), "M1 validation report changed paths are required")
    validate_evidence_bundle(
        root,
        packet,
        lifecycle_task,
        payloads,
        artifacts,
        _current_candidate_binding(root, declared_paths),
    )


def run_self_tests(
    root: Path,
    packet: dict[str, Any],
    lifecycle_task: dict[str, Any],
) -> None:
    payloads, artifacts = load_evidence(root)
    declared_paths = payloads[REPORT_REF].get("changed_paths")
    _require(isinstance(declared_paths, list), "M1 validation report changed paths are required")
    binding = _current_candidate_binding(root, declared_paths)
    def replace_red_output_with_missing_manifest(
        values: dict[str, dict[str, Any]], blobs: dict[str, bytes]
    ) -> None:
        stdout_ref = values[RED_MARKER_REF]["raw_log_refs"][0]
        raw = b"read examples/migration-packs/benchmark-manifest.json: file does not exist\n"
        blobs[stdout_ref] = raw
        values[RED_MARKER_REF]["stdout_digest"] = _sha256(raw)

    def reorder_green_markers(
        values: dict[str, dict[str, Any]], _: dict[str, bytes]
    ) -> None:
        replayed = values[MARKER_REFS[0]]["started_at"]
        values[MARKER_REFS[1]]["started_at"] = replayed
        values[REPORT_REF]["validation_commands"][1]["started_at"] = replayed

    cases: list[tuple[str, Callable[[dict[str, dict[str, Any]], dict[str, bytes]], None]]] = [
        ("candidate binding", lambda values, _: values[SCORECARD_REF].__setitem__("candidate_digest", "sha256:" + "0" * 64)),
        ("raw-log digest", lambda values, blobs: blobs.__setitem__(values[MARKER_REFS[0]]["raw_log_refs"][0], b"mutated")),
        ("marker command", lambda values, _: values[MARKER_REFS[0]].__setitem__("command", "true")),
        ("marker candidate replay", lambda values, _: values[MARKER_REFS[0]].__setitem__("candidate_digest", "sha256:" + "9" * 64)),
        ("marker run replay", lambda values, _: values[MARKER_REFS[0]].__setitem__("validation_run_id", "validation:replayed")),
        ("missing validation run identity", lambda values, _: values[REPORT_REF].pop("validation_run_id")),
        ("missing validation checkout identity", lambda values, _: values[REPORT_REF].pop("validation_checkout_sha")),
        ("marker checkout substitution", lambda values, _: values[MARKER_REFS[0]].__setitem__("git_sha", "f" * 40)),
        ("validator marker failure", lambda values, _: values[MARKER_REFS[4]].update({"execution_status": "fail", "exit_code": 2})),
        ("full-gate marker failure", lambda values, _: values[MARKER_REFS[5]].update({"execution_status": "fail", "exit_code": 2})),
        ("changed paths", lambda values, _: values[REPORT_REF]["changed_paths"].append("outside-task.txt")),
        ("parent acceptance closure", lambda values, _: values[SCORECARD_REF].__setitem__("acceptance_item_ids", ["PILOT-001"])),
        ("red-first status", lambda values, _: values[RED_MARKER_REF].__setitem__("execution_status", "pass")),
        ("red-first reason", lambda values, _: values[RED_MARKER_REF].__setitem__("failure_reason_code", "command_failed")),
        ("red-first candidate replay", lambda values, _: values[RED_MARKER_REF].__setitem__("candidate_digest", "sha256:" + "8" * 64)),
        ("red-first wrong class", replace_red_output_with_missing_manifest),
        ("green marker replay order", reorder_green_markers),
        ("missing holdout result", lambda values, _: values.pop(HOLDOUT_RESULT_REF)),
        ("holdout candidate replay", lambda values, _: values[HOLDOUT_RESULT_REF]["current_work"].__setitem__("candidate_digest", "sha256:" + "7" * 64)),
        ("holdout marker replay", lambda values, _: values[HOLDOUT_RESULT_REF]["current_work"]["work_proof_markers"][0].__setitem__("sha256", "sha256:" + "6" * 64)),
        ("holdout policy widening", lambda values, _: values[HOLDOUT_RESULT_REF].__setitem__("policy_mode", "evaluate")),
        ("missing review report", lambda values, _: values.pop(REVIEW_REPORT_REF)),
        ("review candidate replay", lambda values, _: values[REVIEW_REPORT_REF]["current_work"].__setitem__("candidate_digest", "sha256:" + "5" * 64)),
        ("review marker replay", lambda values, _: values[REVIEW_REPORT_REF]["current_work"]["work_proof_markers"][0].__setitem__("sha256", "sha256:" + "4" * 64)),
        ("review rejection", lambda values, _: values[REVIEW_REPORT_REF].__setitem__("verdict", "changes_required")),
        ("holdout missing evaluator", lambda values, _: values[HOLDOUT_RESULT_REF].__setitem__("evaluator_id", None)),
        ("holdout suite substitution", lambda values, _: values[HOLDOUT_RESULT_REF].__setitem__("suite_ref", "private://unrelated/suite")),
        ("holdout non-hex commitment", lambda values, _: values[HOLDOUT_RESULT_REF].__setitem__("suite_commitment", "hmac-sha256:" + "z" * 64)),
        ("holdout predates validation", lambda values, _: values[HOLDOUT_RESULT_REF].__setitem__("created_at", "2020-01-01T00:00:00Z")),
        ("review missing reviewer", lambda values, _: values[REVIEW_REPORT_REF].__setitem__("reviewer_id", None)),
        ("review lens substitution", lambda values, _: values[REVIEW_REPORT_REF].__setitem__("review_type", "architecture")),
        ("review risk downgrade", lambda values, _: values[REVIEW_REPORT_REF].__setitem__("risk_class", "low")),
        ("review producer mismatch", lambda values, _: values[REVIEW_REPORT_REF]["producer"].__setitem__("producer_class", "security_reviewer")),
        ("review malformed finding", lambda values, _: values[REVIEW_REPORT_REF].__setitem__("findings", [None])),
        ("review predates validation", lambda values, _: values[REVIEW_REPORT_REF].__setitem__("created_at", "2020-01-01T00:00:00Z")),
        ("review incomplete scope", lambda values, _: values[REVIEW_REPORT_REF].__setitem__("review_scope", ["README.md"])),
        ("review private source allowed", lambda values, _: values[REVIEW_REPORT_REF]["context_isolation"]["allowed_sources"].append("private_holdout_inputs")),
        ("review cites future holdout", lambda values, _: values[REVIEW_REPORT_REF]["evidence_refs"].append(HOLDOUT_RESULT_REF)),
        ("review private exclusions removed", lambda values, _: values[REVIEW_REPORT_REF]["context_isolation"].__setitem__("disallowed_sources", [])),
        ("review approvals omitted", lambda values, _: values[REVIEW_REPORT_REF]["approval_effect"].__setitem__("approvals_granted", [])),
        ("review systems analysis omitted", lambda values, _: values[REVIEW_REPORT_REF]["systems_questions"].__setitem__("system_model_summary", "")),
        ("review self-authored", lambda values, _: values[REVIEW_REPORT_REF].update({"reviewer_id": values[REPORT_REF]["implementation_producer"]["producer_id"], "producer": {"worker": "code-review", "producer_id": values[REPORT_REF]["implementation_producer"]["producer_id"], "producer_class": "peer_agent"}})),
        ("holdout self-provisioned", lambda values, _: values[HOLDOUT_RESULT_REF].update({"evaluator_id": values[REPORT_REF]["implementation_producer"]["producer_id"], "producer": {"worker": "holdout-evaluator", "producer_id": values[REPORT_REF]["implementation_producer"]["producer_id"], "producer_class": "independent"}})),
        ("implementation identity omitted", lambda values, _: values[REPORT_REF].pop("implementation_producer")),
        ("holdout before review", lambda values, _: values[HOLDOUT_RESULT_REF].__setitem__("created_at", values[REPORT_REF]["generated_at"])),
        ("review accepted critical risk", lambda values, _: values[REVIEW_REPORT_REF].__setitem__("findings", [{"id": "P1-accepted", "severity": "P1", "summary": "critical risk", "status": "accepted_risk", "claim_source": "independently_observed", "verification_status": "independently_verified", "evidence_refs": [REPORT_REF]}])),
        ("review unverified resolved critical risk", lambda values, _: values[REVIEW_REPORT_REF].__setitem__("findings", [{"id": "P1-unverified", "severity": "P1", "summary": "claimed resolution", "status": "resolved", "claim_source": "builder_provided", "verification_status": "unverified_claim", "evidence_refs": [REPORT_REF]}])),
    ]
    for label, mutate in cases:
        candidate_payloads = copy.deepcopy(payloads)
        candidate_artifacts = copy.deepcopy(artifacts)
        mutate(candidate_payloads, candidate_artifacts)
        try:
            validate_evidence_bundle(
                root,
                packet,
                lifecycle_task,
                candidate_payloads,
                candidate_artifacts,
                binding,
            )
        except AssertionError:
            continue
        raise AssertionError(f"M1 evidence self-test mutation did not fail: {label}")

    def jointly_forge_holdout_policy(
        task: dict[str, Any], values: dict[str, dict[str, Any]]
    ) -> None:
        task["holdout_policy"]["suite_namespace"] = "private://unrelated/v1"
        canonical = {
            key: value
            for key, value in task["holdout_policy"].items()
            if key != "policy_digest"
        }
        digest = _sha256(
            json.dumps(canonical, sort_keys=True, separators=(",", ":")).encode()
        )
        task["holdout_policy"]["policy_digest"] = digest
        values[HOLDOUT_RESULT_REF]["task_policy_digest"] = digest
        values[HOLDOUT_RESULT_REF]["suite_ref"] = "private://unrelated/v1"

    task_cases: list[
        tuple[
            str,
            Callable[[dict[str, Any], dict[str, dict[str, Any]]], None],
        ]
    ] = [
        ("parent work item", lambda task, _: task.__setitem__("work_item_id", "other")),
        ("parent worker chain", lambda task, _: task.__setitem__("required_worker_chain", [])),
        ("parent lifecycle evidence", lambda task, _: task.__setitem__("lifecycle_evidence_required", [])),
        ("jointly forged holdout policy", jointly_forge_holdout_policy),
    ]
    for label, mutate in task_cases:
        candidate_task = copy.deepcopy(lifecycle_task)
        candidate_payloads = copy.deepcopy(payloads)
        mutate(candidate_task, candidate_payloads)
        try:
            validate_evidence_bundle(
                root,
                packet,
                candidate_task,
                candidate_payloads,
                copy.deepcopy(artifacts),
                binding,
            )
        except AssertionError:
            continue
        raise AssertionError(f"M1 lifecycle-task self-test mutation did not fail: {label}")
    candidate_payloads = copy.deepcopy(payloads)
    candidate_artifacts = copy.deepcopy(artifacts)
    candidate_payloads[MARKER_REFS[4]]["execution_status"] = "fail"
    candidate_payloads[MARKER_REFS[4]]["exit_code"] = 2
    previous = os.environ.get("LUMYN_M1_EVIDENCE_RECORDING_REF")
    os.environ["LUMYN_M1_EVIDENCE_RECORDING_REF"] = MARKER_REFS[4]
    try:
        try:
            validate_evidence_bundle(
                root,
                packet,
                lifecycle_task,
                candidate_payloads,
                candidate_artifacts,
                binding,
            )
        except AssertionError:
            pass
        else:
            raise AssertionError(
                "ambient recording environment bypassed a failing M1 marker"
            )
    finally:
        if previous is None:
            os.environ.pop("LUMYN_M1_EVIDENCE_RECORDING_REF", None)
        else:
            os.environ["LUMYN_M1_EVIDENCE_RECORDING_REF"] = previous
    _run_raw_log_loader_confinement_self_test(root)
    _run_candidate_binding_state_self_test()


def _run_raw_log_loader_confinement_self_test(root: Path) -> None:
    with tempfile.TemporaryDirectory(prefix="lumyn-m1-log-loader-") as directory:
        test_root = Path(directory)
        refs = [
            REPORT_REF,
            SCORECARD_REF,
            RED_MARKER_REF,
            *MARKER_REFS,
            HOLDOUT_RESULT_REF,
            REVIEW_REPORT_REF,
        ]
        for marker_ref in [RED_MARKER_REF, *MARKER_REFS]:
            parent = marker_ref.rsplit("/", 1)[0]
            refs.extend(
                [
                    f"{parent}/work-proof-marker.stdout.log",
                    f"{parent}/work-proof-marker.stderr.log",
                ]
            )
        for ref in refs:
            destination = test_root / ref
            destination.parent.mkdir(parents=True, exist_ok=True)
            destination.write_bytes((root / ref).read_bytes())

        marker_path = test_root / RED_MARKER_REF
        marker = json.loads(marker_path.read_text(encoding="utf-8"))
        marker["raw_log_refs"] = ["/dev/zero", "../../outside-repository"]
        marker_path.write_text(json.dumps(marker, indent=2) + "\n", encoding="utf-8")

        payloads, artifacts = load_evidence(test_root)
        _require(
            payloads[RED_MARKER_REF]["raw_log_refs"]
            == ["/dev/zero", "../../outside-repository"]
            and "/dev/zero" not in artifacts
            and "../../outside-repository" not in artifacts,
            "M1 evidence loader followed marker-declared out-of-root log refs",
        )


def _run_candidate_binding_state_self_test() -> None:
    with tempfile.TemporaryDirectory(prefix="lumyn-m1-binding-") as directory:
        root = Path(directory)

        def git(*arguments: str) -> str:
            result = subprocess.run(
                ["git", *arguments],
                cwd=root,
                text=True,
                capture_output=True,
                check=True,
                env={
                    "PATH": os.environ.get("PATH", "/usr/bin:/bin"),
                    "HOME": str(root / "home"),
                    "GIT_CONFIG_NOSYSTEM": "1",
                    "GIT_CONFIG_GLOBAL": os.devnull,
                    "GIT_AUTHOR_DATE": "2026-07-28T12:00:00Z",
                    "GIT_COMMITTER_DATE": "2026-07-28T12:00:00Z",
                    "LC_ALL": "C",
                },
            )
            return result.stdout.strip()

        git("init", "--quiet", "--initial-branch=main")
        (root / "app.txt").write_text("before\n", encoding="utf-8")
        git("add", "app.txt")
        git("-c", "commit.gpgsign=false", "-c", "user.name=Lumyn M1", "-c", "user.email=m1@example.invalid", "commit", "--quiet", "--no-verify", "-m", "base")
        base = git("rev-parse", "HEAD")

        (root / "app.txt").write_text("candidate\n", encoding="utf-8")
        before_commit = candidate_binding(root, base_git_sha=base)
        git("add", "app.txt")
        git("-c", "commit.gpgsign=false", "-c", "user.name=Lumyn M1", "-c", "user.email=m1@example.invalid", "commit", "--quiet", "--no-verify", "-m", "candidate")
        after_commit = candidate_binding(root, before_commit["changed_paths"], base_git_sha=base)
        _require(
            before_commit["candidate_digest"] == after_commit["candidate_digest"],
            "M1 candidate binding changed after committing the exact candidate",
        )

        (root / "unrelated.txt").write_text("later task\n", encoding="utf-8")
        git("add", "unrelated.txt")
        git("-c", "commit.gpgsign=false", "-c", "user.name=Lumyn M1", "-c", "user.email=m1@example.invalid", "commit", "--quiet", "--no-verify", "-m", "later")
        try:
            candidate_binding(root, before_commit["changed_paths"], base_git_sha=base)
        except AssertionError:
            pass
        else:
            raise AssertionError("M1 candidate binding accepted an undeclared later path")

        git("rm", "--quiet", "unrelated.txt")
        git("-c", "commit.gpgsign=false", "-c", "user.name=Lumyn M1", "-c", "user.email=m1@example.invalid", "commit", "--quiet", "--no-verify", "-m", "remove unrelated")

        (root / "app.txt").write_text("mutated later\n", encoding="utf-8")
        after_mutation = candidate_binding(root, before_commit["changed_paths"], base_git_sha=base)
        _require(
            before_commit["candidate_digest"] != after_mutation["candidate_digest"],
            "M1 candidate binding did not detect a later mutation to a bound path",
        )
