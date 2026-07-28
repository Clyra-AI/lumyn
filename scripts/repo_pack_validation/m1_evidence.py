"""Fail-closed validation for attended M1 implementation evidence."""

from __future__ import annotations

import copy
import hashlib
import json
import os
import stat
import subprocess
import tempfile
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Callable


EVIDENCE_ROOT = ".factory/artifacts/task-runs/M1-IMPLEMENTATION"
M1_BASE_GIT_SHA = "d7ae311c391775a2517d56add6d57148d5891ef3"
M1_WORK_ITEM_ID = "lumyn-v3-m1-attended-implementation"
REPORT_REF = f"{EVIDENCE_ROOT}/validation-report.json"
SCORECARD_REF = f"{EVIDENCE_ROOT}/proof-of-behavior-scorecard.json"
RED_MARKER_REF = f"{EVIDENCE_ROOT}/red-first/work-proof-marker.json"
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


def load_evidence(root: Path) -> tuple[dict[str, dict[str, Any]], dict[str, bytes]]:
    payloads: dict[str, dict[str, Any]] = {}
    artifacts: dict[str, bytes] = {}
    for ref in [REPORT_REF, SCORECARD_REF, RED_MARKER_REF, *MARKER_REFS]:
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
    head: str,
    should_pass: bool,
    cwd: str,
    workspace_ref: str,
    validation_run_id: str | None = None,
    candidate_digest: str | None = None,
) -> tuple[datetime, datetime]:
    _require(marker.get("command") == command, f"M1 marker command drifted: {ref}")
    _require(marker.get("marker_id") == marker_id, f"M1 marker identity drifted: {ref}")
    _require(marker.get("generated_by") == "trusted_runner", f"M1 marker is not trusted: {ref}")
    _require(marker.get("git_sha") == head, f"M1 marker base SHA drifted: {ref}")
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
    head: str,
    validation_run_id: str,
    candidate_digest: str,
) -> tuple[datetime, datetime]:
    started, finished = _validate_marker(
        marker,
        artifacts,
        ref=RED_MARKER_REF,
        command=RED_COMMAND,
        marker_id="red-first",
        head=head,
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


def validate_evidence_bundle(
    root: Path,
    packet: dict[str, Any],
    payloads: dict[str, dict[str, Any]],
    artifacts: dict[str, bytes],
    binding: dict[str, Any],
) -> None:
    report = payloads[REPORT_REF]
    scorecard = payloads[SCORECARD_REF]
    head = str(binding["base_git_sha"])
    digest = str(binding["candidate_digest"])
    validation_run_id = report.get("validation_run_id")
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
        declared_binding.get("base_git_sha") == head
        and declared_binding.get("candidate_digest") == digest
        and declared_binding.get("excluded_roots") == CANDIDATE_EXCLUDED_ROOTS,
        "M1 validation report candidate binding drifted",
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
        head=head,
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
            head=head,
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
            and row.get("repo_ref") == head,
            f"M1 validation report row drifted: {command}",
        )
        _require(started <= finished, f"M1 validation command timestamps drifted: {command}")
    _require(_timestamp(report.get("generated_at"), "M1 report generated_at") >= latest, "M1 report predates green evidence")

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


def validate_m1_evidence(root: Path, packet: dict[str, Any]) -> None:
    payloads, artifacts = load_evidence(root)
    declared_paths = payloads[REPORT_REF].get("changed_paths")
    _require(isinstance(declared_paths, list), "M1 validation report changed paths are required")
    validate_evidence_bundle(
        root,
        packet,
        payloads,
        artifacts,
        _current_candidate_binding(root, declared_paths),
    )


def run_self_tests(root: Path, packet: dict[str, Any]) -> None:
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
        ("validator marker failure", lambda values, _: values[MARKER_REFS[4]].update({"execution_status": "fail", "exit_code": 2})),
        ("full-gate marker failure", lambda values, _: values[MARKER_REFS[5]].update({"execution_status": "fail", "exit_code": 2})),
        ("changed paths", lambda values, _: values[REPORT_REF]["changed_paths"].append("outside-task.txt")),
        ("parent acceptance closure", lambda values, _: values[SCORECARD_REF].__setitem__("acceptance_item_ids", ["PILOT-001"])),
        ("red-first status", lambda values, _: values[RED_MARKER_REF].__setitem__("execution_status", "pass")),
        ("red-first reason", lambda values, _: values[RED_MARKER_REF].__setitem__("failure_reason_code", "command_failed")),
        ("red-first candidate replay", lambda values, _: values[RED_MARKER_REF].__setitem__("candidate_digest", "sha256:" + "8" * 64)),
        ("red-first wrong class", replace_red_output_with_missing_manifest),
        ("green marker replay order", reorder_green_markers),
    ]
    for label, mutate in cases:
        candidate_payloads = copy.deepcopy(payloads)
        candidate_artifacts = copy.deepcopy(artifacts)
        mutate(candidate_payloads, candidate_artifacts)
        try:
            validate_evidence_bundle(
                root,
                packet,
                candidate_payloads,
                candidate_artifacts,
                binding,
            )
        except AssertionError:
            continue
        raise AssertionError(f"M1 evidence self-test mutation did not fail: {label}")
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
        refs = [REPORT_REF, SCORECARD_REF, RED_MARKER_REF, *MARKER_REFS]
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
