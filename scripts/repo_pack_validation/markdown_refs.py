"""Validate repo-local Markdown fragments used by compiled planning artifacts."""

from __future__ import annotations

import html
import re
from collections.abc import Iterator
from pathlib import Path
from typing import Any


_ATX_HEADING_RE = re.compile(r"^\s{0,3}#{1,6}\s+(.+?)\s*$")
_CLOSING_HASHES_RE = re.compile(r"\s+#+\s*$")
_FENCE_RE = re.compile(r"^\s{0,3}(`{3,}|~{3,})")
_INLINE_LINK_RE = re.compile(r"!?\[([^\]]*)\]\([^)]*\)")
_REFERENCE_LINK_RE = re.compile(r"!?\[([^\]]*)\]\[[^\]]*\]")
_AUTOLINK_RE = re.compile(r"<((?:https?://|mailto:)[^>]+)>")
_HTML_TAG_RE = re.compile(r"<[^>]*>")
_NON_SLUG_RE = re.compile(r"[^\w\s-]", re.UNICODE)


def _rendered_inline_text(heading: str) -> str:
    rendered = _INLINE_LINK_RE.sub(r"\1", heading)
    rendered = _REFERENCE_LINK_RE.sub(r"\1", rendered)
    rendered = _AUTOLINK_RE.sub(r"\1", rendered)
    rendered = _HTML_TAG_RE.sub("", rendered)
    return html.unescape(rendered).replace("`", "")


def _github_heading_slug(heading: str) -> str:
    """Return the GitHub-style slug used by the repo's plain ATX headings."""

    normalized = _rendered_inline_text(heading).lower()
    normalized = _NON_SLUG_RE.sub("", normalized)
    return re.sub(r"\s", "-", normalized)


def _markdown_anchors(text: str) -> set[str]:
    anchors: set[str] = set()
    fence_character = ""
    fence_length = 0
    for line in text.splitlines():
        fence = _FENCE_RE.match(line)
        if fence_character:
            if (
                fence
                and fence.group(1)[0] == fence_character
                and len(fence.group(1)) >= fence_length
            ):
                fence_character = ""
                fence_length = 0
            continue
        if fence:
            fence_character = fence.group(1)[0]
            fence_length = len(fence.group(1))
            continue
        heading = _ATX_HEADING_RE.match(line)
        if not heading:
            continue
        title = _CLOSING_HASHES_RE.sub("", heading.group(1))
        base_slug = _github_heading_slug(title)
        slug = base_slug
        duplicate_index = 0
        while slug in anchors:
            duplicate_index += 1
            slug = f"{base_slug}-{duplicate_index}"
        anchors.add(slug)
    return anchors


def _strings(value: Any, label: str) -> Iterator[tuple[str, str]]:
    if isinstance(value, str):
        yield label, value
    elif isinstance(value, list):
        for index, item in enumerate(value):
            yield from _strings(item, f"{label}[{index}]")
    elif isinstance(value, dict):
        for key, item in value.items():
            yield from _strings(item, f"{label}.{key}")


def validate_markdown_fragment_refs(
    root: Path,
    payload: Any,
    label: str,
) -> None:
    """Reject local ``*.md#fragment`` refs that do not resolve to headings."""

    anchor_cache: dict[Path, set[str]] = {}
    resolved_root = root.resolve()
    for value_label, value in _strings(payload, label):
        if "#" not in value:
            continue
        relative, fragment = value.split("#", 1)
        if (
            not relative.endswith(".md")
            or relative.startswith(("http://", "https://"))
        ):
            continue
        relative_path = Path(relative)
        if relative_path.is_absolute() or ".." in relative_path.parts:
            raise AssertionError(
                f"{value_label} must stay inside the repository: {relative}"
            )
        markdown_path = (resolved_root / relative_path).resolve()
        if not markdown_path.is_relative_to(resolved_root):
            raise AssertionError(
                f"{value_label} resolves outside the repository: {relative}"
            )
        if not markdown_path.is_file():
            raise AssertionError(
                f"{value_label} points to missing Markdown path {relative}"
            )
        if markdown_path not in anchor_cache:
            anchor_cache[markdown_path] = _markdown_anchors(
                markdown_path.read_text()
            )
        anchors = anchor_cache[markdown_path]
        if not fragment or fragment not in anchors:
            raise AssertionError(
                f"{value_label} points to missing Markdown anchor "
                f"#{fragment} in {relative}"
            )
