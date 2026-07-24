#!/usr/bin/env python3
"""Compatibility entrypoint for the active Lumyn v3 repo-pack self-test."""

from __future__ import annotations

import sys

from repo_pack_validation.validator import run_self_test as _run_self_test


def run_self_test() -> int:
    """Run the canonical v3 negative suite with the historical return contract."""

    try:
        _run_self_test()
    except AssertionError as exc:
        print(f"repo-pack validator self-test failed: {exc}", file=sys.stderr)
        return 2
    print("repo-pack validator self-test passed")
    return 0


def main() -> int:
    if sys.argv[1:]:
        print("usage: repo_pack_self_test.py", file=sys.stderr)
        return 2
    return run_self_test()


if __name__ == "__main__":
    raise SystemExit(main())
