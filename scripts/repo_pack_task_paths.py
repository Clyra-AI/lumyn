#!/usr/bin/env python3
from __future__ import annotations

import re
from typing import Any

from repo_pack_contracts import fail


ARCHITECTURE_PATH_PLANNING_METHOD = "architecture_target_paths_v1"
BROAD_ARCHITECTURE_PATHS = {"cmd", "internal"}


def task_id(task: dict[str, Any]) -> str:
    return str(task.get("task_id") or task.get("id") or "").strip()


def normalize_repo_path(value: object) -> str:
    path = str(value).strip().replace("\\", "/")
    parts: list[str] = []
    for part in path.split("/"):
        if part in {"", "."}:
            continue
        if part == "..":
            if parts:
                parts.pop()
            else:
                parts.append(part)
            continue
        parts.append(part)
    return "/".join(parts)


def repo_relative_path_error(value: object) -> str | None:
    if not isinstance(value, str):
        return "must be a string"
    path = value.strip().replace("\\", "/")
    if not path:
        return "must be non-empty"
    if path.startswith("/") or re.match(r"^[A-Za-z]:", path):
        return "must be repo-relative"
    if any(part == ".." for part in path.split("/")):
        return "must not traverse outside the repository"
    if not normalize_repo_path(path):
        return "must not resolve to the repository root"
    return None


def path_covers_target(scope_path: str, target_path: str) -> bool:
    scope = normalize_repo_path(scope_path)
    target = normalize_repo_path(target_path)
    if scope == target:
        return True
    if scope and scope_path.strip().endswith("/"):
        return target.startswith(scope.rstrip("/") + "/")
    return False


def validate_architecture_target_paths(task: dict[str, Any]) -> None:
    task_id_value = task_id(task)
    targets = task.get("architecture_target_paths")
    if not isinstance(targets, list) or not any(str(value).strip() for value in targets):
        fail(f"{task_id_value}.architecture_target_paths must be a non-empty list")
    if task.get("path_planning_method") != ARCHITECTURE_PATH_PLANNING_METHOD:
        fail(f"{task_id_value}.path_planning_method must be {ARCHITECTURE_PATH_PLANNING_METHOD}")

    normalized_targets: list[str] = []
    for index, target in enumerate(targets):
        path_error = repo_relative_path_error(target)
        if path_error:
            fail(f"{task_id_value}.architecture_target_paths[{index}] {path_error}")
        normalized = normalize_repo_path(target)
        if normalized in BROAD_ARCHITECTURE_PATHS:
            fail(
                f"{task_id_value}.architecture_target_paths[{index}] must name a bounded path, "
                "not broad cmd/ or internal/"
            )
        normalized_targets.append(normalized)

    allowed_paths = task.get("allowed_paths")
    if not isinstance(allowed_paths, list) or not any(str(value).strip() for value in allowed_paths):
        fail(f"{task_id_value}.allowed_paths must be a non-empty list")
    normalized_allowed: list[str] = []
    for index, path in enumerate(allowed_paths):
        path_error = repo_relative_path_error(path)
        if path_error:
            fail(f"{task_id_value}.allowed_paths[{index}] {path_error}")
        normalized = normalize_repo_path(path)
        if normalized in BROAD_ARCHITECTURE_PATHS:
            fail(f"{task_id_value}.allowed_paths must not include broad {normalized}/")
        normalized_allowed.append(normalized)

    missing = sorted(
        target
        for target in normalized_targets
        if not any(path_covers_target(allowed, target) for allowed in normalized_allowed)
    )
    if missing:
        fail(f"{task_id_value}.allowed_paths must include architecture_target_paths: {missing}")
