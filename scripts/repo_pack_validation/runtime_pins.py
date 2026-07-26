"""Canonical runtime pins for the paused Lumyn v3.1 planning generation."""

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
        "provider-originated API or SDK update delivery launched through a "
        "provider-paid, services-assisted sunset campaign; one accountably "
        "confirmed Provider Change Contract is reused by the exact invited "
        "cohort, each non-executable provider event binds its exact digest, "
        "audience, deadline, and supersession state, and neither contract nor "
        "event grants consumer repository authority; a universal registry, "
        "elaborate provider PKI, connection receipts, and receipt-backed billing "
        "are deferred"
    ),
    "agent_runner_policy": (
        "agent_execution_policy defaults to disabled so notify-only, scan-only, "
        "and deterministic-only routes require no Agent Runner credential; an "
        "agent-assisted route pauses until the API Consumer Organization "
        "configures one exact qualified adapter and version, executable source "
        "and digest, permitted noninteractive auth mode and entitlement class, "
        "and observable actual model route; Codex and Claude Code are launch "
        "targets only after common conformance and per-version live canaries, "
        "while Cursor is deferred behind the same gate; every attempt resolves "
        "the approved executable canonically, rejects repository-local PATH "
        "shadowing, and uses a clean ephemeral session with neutral home/config "
        "roots; static native instructions are disabled or explicitly selected "
        "and digest-bound as untrusted context while executable plugins, MCP "
        "servers, and hooks are prohibited; runner host isolation binds mounts, "
        "OS credentials, sockets, descriptors, child processes, egress, and "
        "cleanup to an exact qualified backend identity and hard CPU, memory, "
        "PID/process-tree, disk, and open-file quotas; each agent action freezes "
        "one local-runtime, runner-mediated, direct-model, or hybrid topology "
        "and its minimum disclosure/network/credential set; consumer_managed "
        "is the default configured funding mode, "
        "provider_sponsored_lumyn_managed is optional, credential and usage-"
        "billing owners are explicit, and no adapter, model, endpoint, credential-"
        "owner, or billing-owner fallback is silent"
    ),
    "model_policy": (
        "bounded agent use requires agent_execution_policy=configured, exact "
        "Agent Runner Vendor, adapter and version, executable source/path/digest, "
        "conformance digest, permitted noninteractive auth mode and entitlement "
        "class, clean-session identity, funding mode, credential and usage-billing "
        "owners, native-configuration digest, observable actual "
        "Model Provider, endpoint, model route and version, parameters, prompt and "
        "tool digests, selected context and request-disclosure classes, retention "
        "and regional posture, separate Agent Runner and model credential "
        "environments/scopes and network endpoints/operations, read and write "
        "paths, file, line, diff, turn, token, time, retry, concurrency, attempt, "
        "and cost budgets, plus normalized request, response, tool-call, edit, "
        "usage, error, cancellation, exit, and patch provenance; repository, "
        "provider, native-config, tool, runner, and model content is untrusted and "
        "cannot widen any boundary; managed credentials bind broker issuer, exact "
        "audience and one-hour maximum TTL; one-time broker redemption creates "
        "one attempt-scoped session whose in-attempt calls remain within hard "
        "token and cost quotas, with no refresh, post-attempt replay, or cross-"
        "attempt reuse; revocation and reconciliation use a vendor-native bound "
        "credential or approved budget-enforcing proxy"
    ),
    "verification_policy": (
        "deterministic repository and workflow verification is independent from "
        "deterministic or agent-assisted generation, runs from the exact "
        "candidate head in a fresh process and view with frozen command/config "
        "digests, no Agent Runner/model credentials, and no generation-owned "
        "evidence write handle, under a pinned qualified command-isolation "
        "backend with exact mounts/environment, hard resource quotas, inherited "
        "child limits, offline/lifecycle defaults, and cleanup evidence; a "
        "credential-bearing sandbox entrypoint uses its own exact-head, "
        "endpoint-only, teardown/orphan-evidenced isolation profile; verification "
        "preserves baseline failures and cannot be satisfied by a model "
        "completion or self-verification"
    ),
    "delivery_policy": (
        "complete local evidence plus patch, optional local branch, and PR-ready "
        "bundle is the fallback; at least one tested draft PR requires separate "
        "atomic short-lived non-default-branch-push and draft-PR actions, "
        "evidence-bound idempotency, and no auto-merge; generation, verification, "
        "branch, and PR routes never become one aggregate authority; manual "
        "delivery cannot close automated-"
        "delivery acceptance; first-campaign product proof requires one same "
        "qualifying run from authenticated event and installed preauthorization "
        "through an organically agent-assisted item on a consumer-selected "
        "qualified runner, independent exact-head verification, Lumyn-opened "
        "draft PR, and provider-received consented status projection; "
        "deterministic rerouting cannot manufacture agent proof, and status is "
        "event-bound, provenance-labeled, and never inferred from silence"
    ),
    "artifact_namespace": (
        ".factory/artifacts contains source-safe Factory controls and separately "
        "consented aggregate or digest evidence only; consumer-private code, "
        "prompts, responses, agent sessions, patches, traces, credentials, and "
        "identifiable campaign evidence use an approved external private root and "
        "are never API-provider-visible"
    ),
    "live_work_policy": (
        "Factory worker capabilities are only approval, credentials, and "
        "network and never substitute for Lumyn product action scopes; consumer "
        "installation, event-specific authorization, Agent Runner network, Agent "
        "Runner credential, model disclosure, model network, model credential, "
        "repository, command, registry, sandbox, GitHub, reporting, retention, and "
        "deletion actions are separate, expiring, revocable, fail-closed product "
        "grants whose exact route-selected union is frozen before action"
    ),
    "dispatch_policy": (
        "factoryd dispatch is paused and fail-closed until factoryd is "
        "qualified against this regenerated task, authority, model, evidence, "
        "and closure generation and one bounded task is explicitly unpaused "
        "with positive budgets and complete grants"
    ),
}


def validate_runtime_pins(value: Any, label: str) -> None:
    if value != EXPECTED_RUNTIME_PINS:
        raise AssertionError(
            f"{label} runtime pins differ from the paused Lumyn v3.1 contract"
        )
