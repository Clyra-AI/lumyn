#!/usr/bin/env python3

import hashlib
import json
import subprocess
from pathlib import Path


BASE = "1c8d43e54e8b907873ba80fc1a070018fd12c5be"
ORIGINAL_HEAD = "0cd5f5e6da626cedb294fa02b9d99b3a17cbbae5"
LANDED_HEAD = "4c383c724f68a6aa5fe6ab2c59d830c64797abfe"
EXPECTED_CANDIDATE = "sha256:855f554ae6c8ff1b957c496d4d2983739caee3cdf4abcbb9b8227b66fc3694d9"
EXPECTED_MARKER_SHA = "sha256:0e37ba2322482e018b08f1cc51470cb270590cc670a21d5d6536327b283fa449"


def git_bytes(repo: Path, ref: str, path: str) -> bytes:
    return subprocess.check_output(["git", "show", f"{ref}:{path}"], cwd=repo)


def main() -> None:
    repo = Path(
        subprocess.check_output(
            ["git", "rev-parse", "--show-toplevel"], text=True
        ).strip()
    )
    marker_path = repo / ".factory/artifacts/task-runs/M0/work-proof-marker.json"
    validation_path = repo / ".factory/artifacts/task-runs/M0/validation-report.json"
    review_path = repo / ".factory/artifacts/lifecycle-evidence/M0/review-report.json"

    marker = json.loads(marker_path.read_text())
    validation = json.loads(validation_path.read_text())
    review = json.loads(review_path.read_text())
    marker_sha = "sha256:" + hashlib.sha256(marker_path.read_bytes()).hexdigest()
    verifier_sha = "sha256:" + hashlib.sha256(Path(__file__).read_bytes()).hexdigest()

    assertions = {
        "base_matches": marker["base_head"] == BASE,
        "candidate_matches_marker": marker["candidate_digest"] == EXPECTED_CANDIDATE,
        "candidate_matches_validation": validation["current_work"]["candidate_digest"]
        == EXPECTED_CANDIDATE,
        "candidate_matches_review": review["current_work"]["candidate_digest"]
        == EXPECTED_CANDIDATE,
        "marker_sha_matches": marker_sha == EXPECTED_MARKER_SHA,
        "review_marker_sha_matches": review["current_work"]["work_proof_markers"][0]["sha256"]
        == EXPECTED_MARKER_SHA,
    }
    path_matches = {
        path: git_bytes(repo, ORIGINAL_HEAD, path) == git_bytes(repo, LANDED_HEAD, path)
        for path in marker["changed_paths"]
    }
    assertions["landed_paths_match_original_head"] = all(path_matches.values())
    if not all(assertions.values()):
        raise SystemExit(
            json.dumps(
                {
                    "status": "fail",
                    "assertions": assertions,
                    "mismatched_paths": [path for path, matched in path_matches.items() if not matched],
                },
                sort_keys=True,
            )
        )

    print(
        json.dumps(
            {
                "status": "pass",
                "task_id": "M0",
                "work_item_id": "lumyn-v3-m0",
                "base_head": BASE,
                "original_pr_head": ORIGINAL_HEAD,
                "landed_main_head": LANDED_HEAD,
                "candidate_digest": EXPECTED_CANDIDATE,
                "original_work_proof_marker_sha256": EXPECTED_MARKER_SHA,
                "verifier_sha256": verifier_sha,
                "changed_path_count": len(path_matches),
                "assertions": assertions,
            },
            sort_keys=True,
        )
    )


if __name__ == "__main__":
    main()
