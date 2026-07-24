"""Canonical runtime pins shared by the active Lumyn planning generation."""

from __future__ import annotations

from typing import Any


EXPECTED_RUNTIME_PINS = {
    "language": "go_with_parser_backed_typescript_analysis",
    "toolchain_version": "go1.26.5; exact Node, npm, registry-or-snapshot, package-integrity, and toolchain pins before package-lock mutation",
    "module_or_package_path": "github.com/Clyra-AI/lumyn",
    "dependency_policy": "standard library first; pinned Go modules or a bounded TypeScript parser runtime only with task evidence and license review",
    "distribution_target": "explicitly licensed and integrity-signed design-partner binary/source package for customer-controlled local or CI execution; public OSS/self-serve and Homebrew only after the separate license, security, contribution, support, vulnerability-response, and release-integrity gate",
    "provider_policy": "API-provider change authority is a signed declarative packet whose current trust state is revalidated before every local, sandbox, or remote side effect against a consumer-pinned provider trust root and verified provider-to-package ownership binding, including issuer key, issue time, audience, expiry, rotation, revocation, withdrawal, and replay checks; freshness comes only from a signed offline provider-status snapshot inside the pinned maximum age or an exact endpoint read under provider_trust_status_read, with no repository or consumer data in the request; no model-provider adapter is required for the MVP",
    "artifact_namespace": ".factory/artifacts contains only Factory evidence and separately consented aggregate/hash-only public pilot evidence; private product runtime artifacts use an explicitly configured consumer-controlled root outside the checkout and public source repository; provider export and public commit are irreversible disclosure boundaries, so revocation blocks future sharing and deletes only Lumyn-controlled private copies",
    "live_work_policy": "offline by default; Factory worker approval, credential, and network grants govern implementation work only and never substitute for separate consumer-signed Lumyn authorization; Lumyn revalidates the private bundle and current provider-status evidence immediately before each repository write, host-isolated command, provider-status read, package-registry read, sandbox payload/network/credential action, minimal campaign receipt, provider attestation, remote branch write, GitHub PR, artifact retention, or artifact deletion side effect",
}


def validate_runtime_pins(value: Any, label: str) -> None:
    if value != EXPECTED_RUNTIME_PINS:
        raise AssertionError(f"{label} runtime pins differ from the Lumyn Factory profile contract")
