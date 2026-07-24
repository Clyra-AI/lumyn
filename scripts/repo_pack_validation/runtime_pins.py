"""Canonical runtime pins for the paused Lumyn v3 planning generation."""

from __future__ import annotations

from typing import Any


EXPECTED_RUNTIME_PINS = {
    "language": "go_with_parser_backed_typescript_analysis",
    "toolchain_version": (
        "go1.26.5; exact Node, npm, registry-or-snapshot, package-integrity, "
        "and toolchain pins before package-lock mutation"
    ),
    "module_or_package_path": "github.com/Clyra-AI/lumyn",
    "dependency_policy": (
        "standard library first; pinned Go modules, parser runtime, model "
        "adapter, or tool dependency only with task evidence and license review"
    ),
    "distribution_target": (
        "explicitly licensed and integrity-signed design-partner binary or "
        "source package for consumer-local or consumer-controlled CI execution; "
        "public OSS, self-serve, and Homebrew require a separate approved gate"
    ),
    "provider_intent_policy": (
        "services-led provider-paid API or SDK sunset campaign with "
        "provider-confirmed source and target semantics; a signed declarative "
        "provider packet is authoritative when supplied and confirmed, remains "
        "data, cannot execute code, and grants no consumer repository authority; "
        "mandatory packet PKI, continuous status, connection receipts, and "
        "receipt-backed billing are deferred"
    ),
    "model_policy": (
        "bounded agent use requires exact model provider, endpoint, model and "
        "version, parameters, prompt and tool digests, selected context and "
        "request-disclosure classes, retention and regional posture, credential "
        "environment and scopes, network endpoint and operations, read and write "
        "paths, file, line, diff, turn, token, time, retry, concurrency, attempt, "
        "and cost budgets, plus request, response, tool-call, usage, and patch "
        "provenance; repository, provider, tool, and model content is untrusted "
        "and cannot widen any boundary"
    ),
    "verification_policy": (
        "deterministic repository and workflow verification is independent from "
        "deterministic or agent-assisted generation, runs from the exact "
        "candidate head, preserves baseline failures, and cannot be satisfied by "
        "a model completion or self-verification"
    ),
    "delivery_policy": (
        "complete local evidence plus patch, optional local branch, and PR-ready "
        "bundle is the baseline; optional remote branch and draft PR require "
        "separate exact authorization; default-branch writes and auto-merge are "
        "forbidden"
    ),
    "artifact_namespace": (
        ".factory/artifacts contains source-safe Factory controls and separately "
        "consented aggregate or digest evidence only; consumer-private code, "
        "prompts, responses, patches, traces, credentials, and identifiable "
        "campaign evidence use an approved external private root"
    ),
    "live_work_policy": (
        "Factory worker capabilities are only approval, credentials, and "
        "network and never substitute for Lumyn product action scopes; model "
        "disclosure, model network, model credential, repository, command, "
        "registry, sandbox, GitHub, reporting, retention, and deletion actions "
        "are separate, expiring, revocable, fail-closed product grants"
    ),
    "dispatch_policy": (
        "factoryd dispatch is paused and fail-closed until the external "
        "profiles/lumyn.yaml contract is updated for v3 and factoryd is "
        "qualified against this regenerated task, authority, model, evidence, "
        "and closure generation"
    ),
}


def validate_runtime_pins(value: Any, label: str) -> None:
    if value != EXPECTED_RUNTIME_PINS:
        raise AssertionError(
            f"{label} runtime pins differ from the paused Lumyn v3 contract"
        )
