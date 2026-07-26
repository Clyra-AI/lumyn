# Lumyn V3 CLI Contract

Status: M2 grammar and compatibility contract; commands remain unavailable
until their owning runtime milestones land.

## Compatibility Rules

- `lumyn init` and `lumyn check` retain their implemented M0 behavior.
- Every state-returning command emits `lumyn.command_result` JSON. `--json` is
  accepted before the command and on the retained M0 commands.
- An unavailable v3 production command returns a nonzero typed result. A help
  entry, schema, or fixture never makes runtime behavior available.
- File arguments identify one exact artifact. Implementations resolve and
  digest the bytes before use; a later mutation makes dependent evidence stale.
- Commands never accept inline credentials, reusable token files, default-
  branch write, auto-merge, wildcard paths, or implicit current-session agent
  authority.
- A composed command dereferences and revalidates its source artifacts; it does
  not synthesize their aggregate authority.

## Grammar

Square brackets are optional; `...` means one or more repetitions of the
immediately preceding option.

```text
lumyn [--json] init [--config PATH] [--openapi PATH] [--docs PATH]
lumyn [--json] check [--config PATH] [--strict]

lumyn [--json] pack build --evidence PATH... --output PATH
lumyn [--json] pack validate --pack PATH

lumyn [--json] install create \
  --event-url HTTPS_URL --campaign-key PATH --repository PATH \
  --policy PATH --output PATH
lumyn [--json] install validate --installation PATH
lumyn [--json] install revoke --installation PATH --reason TEXT

lumyn [--json] update --event PATH_OR_PINNED_URL \
  --installation PATH --private-root PATH
lumyn [--json] impact --event PATH --installation PATH \
  --execution-manifest PATH --output PATH
lumyn [--json] plan --event PATH --installation PATH \
  --impact PATH --output PATH
lumyn [--json] apply --authorization PATH --plan PATH --private-root PATH
lumyn [--json] candidate import --manual --authorization PATH \
  --plan PATH --patch PATH --output PATH
lumyn [--json] verify --authorization PATH --candidate PATH \
  --execution-manifest PATH --output PATH
lumyn [--json] repair --authorization PATH --verification PATH \
  --output PATH

lumyn [--json] export --authorization PATH --candidate PATH \
  --verification PATH --mode patch|local_branch|pr_bundle --output PATH
lumyn [--json] pr create --draft --authorization PATH \
  --export-result PATH --head REF --base REF
lumyn [--json] trace --artifact PATH [--output PATH]
lumyn [--json] outcome record --event PATH --installation PATH \
  --state accepted|merged|closed|reverted|corrected --evidence PATH
```

`update` is the only composed product command. It stops at the installation's
action ceiling and emits the same individually schema-valid artifacts as the
corresponding commands. It does not create a branch or PR unless the immutable
event authorization separately contains the exact remote-branch and draft-PR
grants and the runtime obtains distinct short-lived credentials for each
action.

## Typed Error And Exit Compatibility

Errors use stable snake-case `code`, a non-secret `message`, and the retained
exit-code range. The artifact reference causing a failure is reported through
the command result, not embedded as an unstructured stack trace.

| Exit | Contract class | Representative v3 errors |
|---:|---|---|
| `0` | success | none |
| `1` | internal invariant | `internal_error` |
| `2` | usage, local input, or product authorization | `command_not_implemented`, `migration_pack_invalid`, `authorization_missing`, `authorization_expired`, `read_scope_exceeded`, `write_scope_exceeded`, `agent_not_authorized`, `human_input_required` |
| `3` | source completeness | `source_target_mismatch`, `repo_context_insufficient`, `impact_uncertain`, `dynamic_usage_uncertain` |
| `4` | artifact or workflow contract | `provider_event_invalid`, `provider_confirmation_missing`, `plan_stale`, `candidate_stale`, `evidence_stale`, `agent_runner_contract_violation` |
| `5` | independent verification | `compile_failed`, `typecheck_failed`, `tests_failed`, `replay_failed`, `workflow_proof_gap` |
| `6` | reserved compatibility slot | no v3 migration meaning assigned in M2 |
| `7` | credential, auth, entitlement, or isolation environment | `agent_runner_auth_failed`, `agent_runner_entitlement_invalid`, `agent_runner_executable_untrusted`, `redaction_uncertain` |
| `8` | approved dependency or network route | `agent_runner_unavailable`, `model_unavailable`, `sandbox_unavailable`, `export_failed` when remote delivery is unavailable |
| `9` | persisted evidence or replay integrity | `candidate_conflict`, `duplicate_pr`, `trace_integrity_failed` |

One error may have several contributing findings, but its exit class is chosen
from the first boundary that prevents safe progress. No error code changes an
installation or event authorization, and a model or Agent Runner message never
becomes a typed Lumyn result without normalization and policy validation.
