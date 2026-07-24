"""Keep compiled acceptance text identical to the authored PRD."""

from __future__ import annotations

import re
from typing import Any


_ITEM = re.compile(
    r"^\d+\. `([A-Z]+-\d{3})`: (.*?)(?=^\d+\. `[A-Z]+-\d{3}`: |^---\s*$|^##|\Z)",
    re.MULTILINE | re.DOTALL,
)


def _normalize(value: str) -> str:
    return " ".join(value.split())


def validate_acceptance_text(prd_text: str, ledger: dict[str, Any]) -> None:
    authored = {item_id: _normalize(text) for item_id, text in _ITEM.findall(prd_text)}
    compiled = {
        str(item.get("acceptance_item_id")): _normalize(str(item.get("source_text", "")))
        for item in ledger.get("items", [])
        if isinstance(item, dict)
    }
    if authored != compiled:
        missing = sorted(set(authored) - set(compiled))
        extra = sorted(set(compiled) - set(authored))
        changed = sorted(item_id for item_id in set(authored) & set(compiled) if authored[item_id] != compiled[item_id])
        raise AssertionError(
            f"compiled acceptance text differs from PRD: missing={missing} extra={extra} changed={changed}"
        )
