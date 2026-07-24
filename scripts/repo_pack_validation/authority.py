"""Validate Factory worker grants without conflating Lumyn product authority."""

from __future__ import annotations

import hashlib
import json
import re
from ipaddress import ip_address
from datetime import datetime, timezone
from pathlib import PurePosixPath
from typing import Any


FACTORY_CAPABILITIES = {"approval", "credentials", "network"}
_RFC3339 = re.compile(
    r"^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})$"
)
_SHA256 = re.compile(r"^sha256:[0-9a-f]{64}$")
_DNS_NAME = re.compile(
    r"^(?=.{1,253}$)(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)*"
    r"[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$",
    re.IGNORECASE,
)
_SEMANTIC_WILDCARDS = {
    "all",
    "any",
    "default",
    "everything",
    "global",
    "unrestricted",
}
_PREFLIGHT_DIGEST_FIELDS = (
    "enforcement_point",
    "owner",
    "purpose",
    "allowed_private_fields",
    "approved_private_storage_boundary",
    "participant_consent_required",
    "retention_ttl_and_expiry_required",
    "deletion_on_revocation_required",
    "deletion_receipt_and_orphan_owner_required",
    "public_fields",
    "public_disclosure_requires_separate_consent",
    "prohibited_actions",
    "product_runtime_authority",
    "failure_behavior",
)


def manual_preflight_scope_digest(preflight: dict[str, Any]) -> str:
    """Hash the exact governed M2.5 preflight projection deterministically."""

    canonical = json.dumps(
        {field: preflight.get(field) for field in _PREFLIGHT_DIGEST_FIELDS},
        ensure_ascii=False,
        separators=(",", ":"),
        sort_keys=True,
    )
    return f"sha256:{hashlib.sha256(canonical.encode()).hexdigest()}"


def _require(condition: bool, message: str) -> None:
    if not condition:
        raise AssertionError(message)


def _text(value: Any) -> bool:
    return isinstance(value, str) and bool(value.strip())


def _exact_text(value: Any, label: str, *, allow_brackets: bool = False) -> None:
    _require(_text(value), f"{label} must be a non-empty string")
    text = str(value)
    forbidden = r"[\r\n*\\]" if allow_brackets else r"[\r\n*\\\[\]]"
    _require(text == text.strip() and not re.search(forbidden, text), f"{label} must be literal and wildcard-free")
    _require(text.casefold() not in _SEMANTIC_WILDCARDS, f"{label} must not use a semantic wildcard")


def _exact_list(value: Any, label: str) -> None:
    _require(isinstance(value, list) and bool(value), f"{label} must be a non-empty list")
    normalized = [str(item).casefold() for item in value]
    _require(len(normalized) == len(set(normalized)), f"{label} must not contain case-insensitive duplicates")
    for index, item in enumerate(value):
        _exact_text(item, f"{label}[{index}]")


def _network_endpoint(value: Any, label: str) -> None:
    """Require an exact host or host:port, never a URL, CIDR, or broad bind."""

    _exact_text(value, label, allow_brackets=True)
    endpoint = str(value)
    _require("://" not in endpoint and "/" not in endpoint and "?" not in endpoint and "#" not in endpoint, f"{label} must be an exact host or host:port")

    host = endpoint
    port: str | None = None
    if endpoint.startswith("["):
        close = endpoint.find("]")
        _require(close > 1, f"{label} has an invalid bracketed IPv6 host")
        host = endpoint[1:close]
        suffix = endpoint[close + 1 :]
        _require(not suffix or suffix.startswith(":"), f"{label} has an invalid IPv6 endpoint suffix")
        port = suffix[1:] if suffix else None
    elif endpoint.count(":") == 1:
        host, port = endpoint.rsplit(":", 1)
    elif endpoint.count(":") > 1:
        # A bare IPv6 address is exact; bracket notation is required when a port is present.
        host = endpoint

    _require(bool(host), f"{label} host must be non-empty")
    if port is not None:
        _require(port.isdigit() and 1 <= int(port) <= 65535, f"{label} port must be in 1..65535")

    try:
        parsed_ip = ip_address(host)
    except ValueError:
        _require(_DNS_NAME.fullmatch(host) is not None, f"{label} host must be a literal DNS name or IP address")
        _require(host.casefold() not in {"localhost", "local", "invalid"}, f"{label} host must identify an approved remote endpoint")
    else:
        _require(not parsed_ip.is_unspecified, f"{label} must not use an unspecified address")


def _scoped_ref(value: Any, label: str) -> None:
    _exact_text(value, label)
    ref = str(value)
    _require(not ref.startswith(("/", "~", "file:")), f"{label} must not be machine-local")
    if ":" not in ref:
        _require(".." not in PurePosixPath(ref.split("#", 1)[0]).parts, f"{label} must stay inside its declared artifact root")


def _future_expiry(value: Any, label: str) -> None:
    _exact_text(value, label)
    raw = str(value)
    _require(_RFC3339.fullmatch(raw) is not None, f"{label} must be RFC3339")
    parsed = datetime.fromisoformat(raw.replace("Z", "+00:00"))
    _require(parsed.tzinfo is not None and parsed > datetime.now(timezone.utc), f"{label} must be a future instant")


def validate_active_repo_safety(repo: Any, artifact_refs: dict[str, str]) -> None:
    """Keep the gitignored active daemon config attended and locally bounded."""

    _require(isinstance(repo, dict), ".factory/factoryd.json missing repos.lumyn")
    ref_fields = {
        "task_packets": "packets",
        "scope_closure_map": "closure",
        "validation_contract": "contract",
        "acceptance_ledger": "ledger",
    }
    for field, ref_key in ref_fields.items():
        _require(
            repo.get(field) == artifact_refs[ref_key],
            f".factory/factoryd.json {field} ref is stale",
        )
    shipping = repo.get("shipping")
    _require(
        repo.get("auto_ship") is False
        and isinstance(shipping, dict)
        and shipping.get("enabled") is False,
        ".factory/factoryd.json must remain safe-attended; use "
        "factoryd.autoship.example.json for full-loop shipping",
    )


def validate_authority_grants(
    grants: Any,
    tasks: dict[str, dict[str, Any]],
    expected_capabilities: dict[str, set[str]],
) -> None:
    """Validate Factory execution grants.

    Exact Lumyn product permissions are separate private authorization artifacts
    named by each task's ``product_authority_requirements``. Factory grants only
    authorize the worker to consume an already validated product bundle.
    """

    _require(isinstance(grants, list), "active capability_grants must be a list")
    seen: set[tuple[str, str]] = set()
    approved: dict[str, set[str]] = {}
    for index, grant in enumerate(grants):
        label = f"active capability grant {index}"
        _require(isinstance(grant, dict), f"{label} must be an object")
        task_id = str(grant.get("task_id", "")).strip()
        capability = str(grant.get("capability", "")).strip()
        _require(task_id in tasks and task_id != "*", f"{label} references unknown or wildcard task {task_id}")
        _require(capability in FACTORY_CAPABILITIES, f"{label} references unknown or product capability {capability}")
        pair = (task_id, capability)
        _require(pair not in seen, f"{label} duplicates grant {task_id}/{capability}")
        seen.add(pair)
        declared = set(tasks[task_id].get("requires_capabilities", []))
        _require(capability in declared, f"{label} capability is not declared by task {task_id}")
        if grant.get("approved") is not True:
            continue
        _scoped_ref(grant.get("evidence_ref"), f"{label}.evidence_ref")
        _future_expiry(grant.get("expires_at"), f"{label}.expires_at")
        preflight = tasks[task_id].get("manual_external_evidence_preflight")
        if capability == "approval" and isinstance(preflight, dict):
            expected_digest = manual_preflight_scope_digest(preflight)
            _require(
                preflight.get("approval_scope_digest") == expected_digest,
                f"{task_id} manual preflight approval_scope_digest must match its canonical content",
            )
            supplied_digest = grant.get("approval_scope_digest")
            _require(_SHA256.fullmatch(expected_digest) is not None, f"{task_id} manual preflight must declare a canonical SHA-256 approval_scope_digest")
            _require(
                supplied_digest == expected_digest,
                f"{label}.approval_scope_digest must exactly bind the current manual preflight scope",
            )
            _require(
                grant.get("evidence_ref") == preflight.get("approval_evidence_ref"),
                f"{label}.evidence_ref must cite the current manual preflight",
            )
        if capability == "credentials":
            _exact_list(grant.get("credential_scopes"), f"{label}.credential_scopes")
            _exact_text(grant.get("credential_environment"), f"{label}.credential_environment")
        elif capability == "network":
            _exact_list(grant.get("network_allowlist"), f"{label}.network_allowlist")
            for endpoint_index, endpoint in enumerate(grant["network_allowlist"]):
                _network_endpoint(endpoint, f"{label}.network_allowlist[{endpoint_index}]")
        approved.setdefault(task_id, set()).add(capability)

    for task_id, expected in expected_capabilities.items():
        task_grants = approved.get(task_id, set())
        if not task_grants:
            continue
        missing = expected - task_grants
        _require(not missing, f"{task_id} active grants missing exact Factory capabilities {sorted(missing)}")
