#!/usr/bin/env python3

import hashlib
import json
import subprocess
import tempfile
from pathlib import Path


BASE = "d7ae311c391775a2517d56add6d57148d5891ef3"
VALIDATION_CHECKOUT = "bc17b58f128bc15c4e715aba15c505c20b224e35"
ORIGINAL_HEAD = "656c1f0bbb61cd558500c2f1b91a5a8f084f4f29"
LANDED_HEAD = "702b5be8d53b46a8c2a394f0b00770f626a8bbdd"
BUNDLE_REF = "refs/heads/codex/lumyn-m1-runner-contract"
BUNDLE_LANDED_REF = "refs/remotes/origin/main"
BUNDLE_PATH = (
    ".factory/artifacts/pr-lifecycle/lumyn-v3-m1/implementation/"
    "pr74-original-head.bundle"
)
EXPECTED_BUNDLE = (
    "sha256:0fa0bb1199ba167d2f97b015a37300c72f9e462bea08495ef90275e91f199f4d"
)
EXPECTED_VALIDATION_RUN = (
    "validation:2026-07-28T17:05:04Z:lumyn-m1-implementation-r2"
)
EXPECTED_CANDIDATE = (
    "sha256:971af9dd40e66ea47b220bcb4a817c62705e9764926eeb84ab4787564b8b2ed4"
)
EXPECTED_MARKERS = {
    ".factory/artifacts/task-runs/M1-IMPLEMENTATION/validation-001/work-proof-marker.json":
        "sha256:f6526b604667f712631f3ea64ac206dcbddbcb6ae53b1683a95b30c20be8f690",
    ".factory/artifacts/task-runs/M1-IMPLEMENTATION/validation-002/work-proof-marker.json":
        "sha256:469314cd909950a2f9653bc221dd233cdcf5414f302443c3d92da13e9d35a36c",
    ".factory/artifacts/task-runs/M1-IMPLEMENTATION/validation-003/work-proof-marker.json":
        "sha256:3ab2e06cfc85ecc72f74ba5bbbeb71094e94199f75f277a79659df511f526474",
    ".factory/artifacts/task-runs/M1-IMPLEMENTATION/validation-004/work-proof-marker.json":
        "sha256:9044b98a1cd40df6f5f50942d747e8656cc27062fcb2c369f0e67b0d28863015",
    ".factory/artifacts/task-runs/M1-IMPLEMENTATION/validation-005/work-proof-marker.json":
        "sha256:b04c54c6b2c68d0932534c70b4274377efe878d54911bc5c84ea2585b0c893b3",
    ".factory/artifacts/task-runs/M1-IMPLEMENTATION/validation-006/work-proof-marker.json":
        "sha256:ac249655d16994f03b00b2e173c7e59008931c4d5d1c1915f1203604d08aa09f",
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
        [
            "git",
            "fetch",
            "--quiet",
            str(bundle),
            f"{BUNDLE_REF}:refs/evidence/original",
            f"{BUNDLE_LANDED_REF}:refs/evidence/landed",
        ],
        check=True,
        cwd=archive,
    )


def main() -> None:
    repo = Path(
        subprocess.check_output(
            ["git", "rev-parse", "--show-toplevel"], text=True
        ).strip()
    )
    validation_path = (
        repo / ".factory/artifacts/task-runs/M1-IMPLEMENTATION/validation-report.json"
    )
    scorecard_path = (
        repo
        / ".factory/artifacts/task-runs/M1-IMPLEMENTATION/proof-of-behavior-scorecard.json"
    )
    review_path = repo / ".factory/artifacts/lifecycle-evidence/M1/review-report.json"
    holdout_path = repo / ".factory/artifacts/lifecycle-evidence/M1/holdout-result.json"
    validation = json.loads(validation_path.read_text())
    scorecard = json.loads(scorecard_path.read_text())
    review = json.loads(review_path.read_text())
    holdout = json.loads(holdout_path.read_text())
    binding = validation["candidate_binding"]
    bundle = repo / BUNDLE_PATH
    bundle_digest = sha256(bundle)

    marker_digests = {ref: sha256(repo / ref) for ref in EXPECTED_MARKERS}
    review_markers = {
        item["ref"]: item["sha256"]
        for item in review["current_work"]["work_proof_markers"]
    }
    holdout_markers = {
        item["ref"]: item["sha256"]
        for item in holdout["current_work"]["work_proof_markers"]
    }
    with tempfile.TemporaryDirectory(prefix="lumyn-m1-binding-") as archive_dir:
        archive = Path(archive_dir)
        populate_retained_repository(repo, bundle, archive)
        assertions = {
            "base_matches": binding["base_git_sha"] == BASE,
            "validation_checkout_matches": (
                validation["validation_checkout_sha"] == VALIDATION_CHECKOUT
            ),
            "candidate_matches_validation": (
                binding["candidate_digest"] == EXPECTED_CANDIDATE
            ),
            "candidate_matches_scorecard": (
                scorecard["candidate_digest"] == EXPECTED_CANDIDATE
            ),
            "candidate_matches_review": (
                review["current_work"]["candidate_digest"] == EXPECTED_CANDIDATE
            ),
            "candidate_matches_holdout": (
                holdout["current_work"]["candidate_digest"] == EXPECTED_CANDIDATE
            ),
            "validation_run_matches_validation": (
                validation["validation_run_id"] == EXPECTED_VALIDATION_RUN
            ),
            "validation_run_matches_review": (
                review["current_work"]["validation_run_id"]
                == EXPECTED_VALIDATION_RUN
            ),
            "validation_run_matches_holdout": (
                holdout["current_work"]["validation_run_id"]
                == EXPECTED_VALIDATION_RUN
            ),
            "marker_digests_match": marker_digests == EXPECTED_MARKERS,
            "review_marker_digests_match": review_markers == EXPECTED_MARKERS,
            "holdout_marker_digests_match": holdout_markers == EXPECTED_MARKERS,
            "retained_bundle_digest_matches": bundle_digest == EXPECTED_BUNDLE,
            "retained_bundle_head_matches": subprocess.check_output(
                ["git", "rev-parse", "refs/evidence/original"],
                cwd=archive,
                text=True,
            ).strip()
            == ORIGINAL_HEAD,
            "retained_bundle_landed_head_matches": subprocess.check_output(
                ["git", "rev-parse", "refs/evidence/landed"],
                cwd=archive,
                text=True,
            ).strip()
            == LANDED_HEAD,
        }
        path_matches = {
            path: git_bytes(archive, ORIGINAL_HEAD, path)
            == git_bytes(archive, LANDED_HEAD, path)
            for path in validation["changed_paths"]
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
                "task_id": "M1",
                "work_item_id": "lumyn-v3-m1",
                "base_head": BASE,
                "validation_checkout_sha": VALIDATION_CHECKOUT,
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
