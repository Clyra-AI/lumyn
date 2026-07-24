#!/usr/bin/env python3
"""Enforce overall and stable-package Go statement coverage."""

from __future__ import annotations

import argparse
import re
import subprocess
import sys
from pathlib import Path


TOTAL_RE = re.compile(r"^total:\s+\(statements\)\s+([0-9]+(?:\.[0-9]+)?)%$")
PROFILE_RE = re.compile(r"^(?P<file>.+):[^ ]+\s+(?P<statements>\d+)\s+(?P<count>\d+)$")


def overall_coverage(coverprofile: Path) -> float:
    result = subprocess.run(
        ["go", "tool", "cover", f"-func={coverprofile}"],
        check=False,
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )
    if result.returncode != 0:
        raise ValueError(result.stderr.strip() or "go tool cover failed")
    for line in result.stdout.splitlines():
        match = TOTAL_RE.match(line.strip())
        if match:
            return float(match.group(1))
    raise ValueError("could not find total coverage in go tool cover output")


def repo_relative_source(source: str) -> str:
    normalized = source.replace("\\", "/")
    marker = "github.com/Clyra-AI/lumyn/"
    if marker in normalized:
        return normalized.split(marker, 1)[1]
    return normalized.lstrip("./")


def stable_coverage(coverprofile: Path, packages: list[str]) -> tuple[float, int]:
    total = 0
    covered = 0
    prefixes = tuple(f"{value.strip('/')}/" for value in packages)
    for line in coverprofile.read_text().splitlines():
        if not line or line.startswith("mode:"):
            continue
        match = PROFILE_RE.match(line)
        if not match:
            raise ValueError(f"invalid coverprofile row: {line}")
        if not repo_relative_source(match.group("file")).startswith(prefixes):
            continue
        statements = int(match.group("statements"))
        total += statements
        if int(match.group("count")) > 0:
            covered += statements
    if total == 0:
        raise ValueError(
            "stable-package selection matched no statements: " + ", ".join(packages)
        )
    return (covered / total) * 100.0, total


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument("coverprofile", type=Path)
    parser.add_argument("overall_minimum", type=float)
    parser.add_argument("--stable-minimum", type=float, required=True)
    parser.add_argument("--stable-package", action="append", required=True)
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    if not args.coverprofile.is_file():
        print(f"coverage profile not found: {args.coverprofile}", file=sys.stderr)
        return 2
    try:
        overall = overall_coverage(args.coverprofile)
        stable, statements = stable_coverage(
            args.coverprofile, args.stable_package
        )
    except ValueError as exc:
        print(str(exc), file=sys.stderr)
        return 2

    print(
        f"go coverage total: {overall:.1f}% "
        f"(minimum {args.overall_minimum:.1f}%)"
    )
    print(
        f"go stable-package coverage: {stable:.1f}% across {statements} statements "
        f"(minimum {args.stable_minimum:.1f}%)"
    )
    failures = []
    if overall + 1e-9 < args.overall_minimum:
        failures.append(
            f"overall {overall:.1f}% < {args.overall_minimum:.1f}%"
        )
    if stable + 1e-9 < args.stable_minimum:
        failures.append(
            f"stable packages {stable:.1f}% < {args.stable_minimum:.1f}%"
        )
    if failures:
        print("go coverage below minimum: " + "; ".join(failures), file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
