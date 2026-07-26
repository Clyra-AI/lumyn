#!/usr/bin/env python3

import hashlib
import json
import subprocess
import tempfile
from pathlib import Path


BASE = "7609e5c49c0776c1028c1aeb3e2e2ee942b613b6"
ORIGINAL_HEAD = "9345f3392ec98eb0e10345fe7941fd9d1450e55b"
LANDED_HEAD = "f89bc82490ffb6df908df6f8572054ee051ed6c6"
BUNDLE_REF = "refs/heads/codex/lumyn-m2-contracts"
BUNDLE_PATH = (
    ".factory/artifacts/pr-lifecycle/lumyn-v3-m2/implementation/"
    "pr72-original-head.bundle"
)
EXPECTED_BUNDLE = (
    "sha256:1052339f87c09336ef87c890f3e6f7cd20184ec3dd27bb9ea62d7d3ab1dec2e6"
)
EXPECTED_VALIDATION_RUN = (
    "validation:2026-07-26T16:32:25Z:50773146792248c5907105189e125e8a"
)
EXPECTED_CANDIDATE = (
    "sha256:8509fbe0980c17d6e6ad3d9cd2b8f04e122d9d6699cb7162b33492198e040324"
)
EXPECTED_MARKERS = {
    ".factory/artifacts/task-runs/M2/validation-001/work-proof-marker.json":
        "sha256:7605a554352fdc6b1abd8e523916bedefcc468067b6c5064aacfb662c34006dc",
    ".factory/artifacts/task-runs/M2/validation-002/work-proof-marker.json":
        "sha256:c4d30505c686940995d58abc7ebb6d2c1ecc9ff3e42ae92370fcdd58fc2d4331",
    ".factory/artifacts/task-runs/M2/validation-003/work-proof-marker.json":
        "sha256:845df89efb744c6978b9e55b12703b5f48ea3a8cfae22f13451262850d50e1ca",
    ".factory/artifacts/task-runs/M2/validation-004/work-proof-marker.json":
        "sha256:a011989ecd5fdb0f3178897faef8913ed08a481151a9ea1f9ec685f3fd1454c7",
}


def sha256(path: Path) -> str:
    return "sha256:" + hashlib.sha256(path.read_bytes()).hexdigest()


def git_bytes(repo: Path, ref: str, path: str) -> bytes:
    return subprocess.check_output(["git", "show", f"{ref}:{path}"], cwd=repo)


def populate_retained_repository(repo: Path, bundle: Path, archive: Path) -> None:
    subprocess.run(
        ["git", "init", "--quiet", "--bare", str(archive)],
        check=True,
        cwd=repo,
    )
    subprocess.run(
        ["git", "fetch", "--quiet", str(repo), LANDED_HEAD],
        check=True,
        cwd=archive,
    )
    subprocess.run(
        ["git", "fetch", "--quiet", str(bundle), f"{BUNDLE_REF}:refs/evidence/original"],
        check=True,
        cwd=archive,
    )
def main() -> None:
    repo = Path(
        subprocess.check_output(
            ["git", "rev-parse", "--show-toplevel"], text=True
        ).strip()
    )
    summary_path = repo / ".factory/artifacts/task-runs/M2/validation-run-summary.json"
    validation_path = repo / ".factory/artifacts/task-runs/M2/validation-report.json"
    review_path = repo / ".factory/artifacts/lifecycle-evidence/M2/review-report.json"
    summary = json.loads(summary_path.read_text())
    validation = json.loads(validation_path.read_text())
    review = json.loads(review_path.read_text())
    binding = summary["candidate_binding"]
    bundle = repo / BUNDLE_PATH
    bundle_digest = sha256(bundle)

    marker_digests = {
        ref: sha256(repo / ref)
        for ref in EXPECTED_MARKERS
    }
    review_markers = {
        item["ref"]: item["sha256"]
        for item in review["current_work"]["work_proof_markers"]
    }
    with tempfile.TemporaryDirectory(prefix="lumyn-m2-binding-") as archive_dir:
        archive = Path(archive_dir)
        populate_retained_repository(repo, bundle, archive)
        assertions = {
            "base_matches": binding["base_git_sha"] == BASE,
            "candidate_matches_summary": summary["candidate_digest"] == EXPECTED_CANDIDATE,
            "candidate_matches_binding": binding["candidate_digest"] == EXPECTED_CANDIDATE,
            "candidate_matches_validation": validation["candidate_digest"] == EXPECTED_CANDIDATE,
            "candidate_matches_review": review["current_work"]["candidate_digest"] == EXPECTED_CANDIDATE,
            "validation_run_matches_summary": summary["validation_run_id"] == EXPECTED_VALIDATION_RUN,
            "validation_run_matches_validation": validation["validation_run_id"] == EXPECTED_VALIDATION_RUN,
            "validation_run_matches_review": review["current_work"]["validation_run_id"] == EXPECTED_VALIDATION_RUN,
            "marker_digests_match": marker_digests == EXPECTED_MARKERS,
            "review_marker_digests_match": review_markers == EXPECTED_MARKERS,
            "retained_bundle_digest_matches": bundle_digest == EXPECTED_BUNDLE,
            "retained_bundle_head_matches": subprocess.check_output(
                ["git", "rev-parse", "refs/evidence/original"],
                cwd=archive,
                text=True,
            ).strip() == ORIGINAL_HEAD,
        }
        path_matches = {
            path: git_bytes(archive, ORIGINAL_HEAD, path)
            == git_bytes(repo, LANDED_HEAD, path)
            for path in binding["changed_paths"]
        }
    assertions["landed_paths_match_original_head"] = all(path_matches.values())

    if not all(assertions.values()):
        raise SystemExit(
            json.dumps(
                {
                    "status": "fail",
                    "assertions": assertions,
                    "mismatched_paths": [
                        path for path, matched in path_matches.items() if not matched
                    ],
                },
                sort_keys=True,
            )
        )

    print(
        json.dumps(
            {
                "status": "pass",
                "task_id": "M2",
                "work_item_id": "lumyn-v3-m2",
                "base_head": BASE,
                "original_pr_head": ORIGINAL_HEAD,
                "retained_bundle_ref": BUNDLE_PATH,
                "retained_bundle_sha256": bundle_digest,
                "landed_main_head": LANDED_HEAD,
                "validation_run_id": EXPECTED_VALIDATION_RUN,
                "candidate_digest": EXPECTED_CANDIDATE,
                "changed_path_count": len(path_matches),
                "marker_digests": marker_digests,
                "verifier_sha256": sha256(Path(__file__)),
                "assertions": assertions,
            },
            sort_keys=True,
        )
    )


if __name__ == "__main__":
    main()
