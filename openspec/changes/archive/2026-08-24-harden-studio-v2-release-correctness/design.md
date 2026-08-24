## Context

This is a release-hardening follow-up to the active `fix-protocol-address-adaptation-v2` change. That predecessor has completed Spectra artifacts but its implementation baseline is still present in the shared staged worktree, so this change must be applied only after the predecessor is reconciled. The hardening crosses the frontend Step 2 rule planner, backend source-rule validation and ownership, SQLite migration helpers, test runners, and the Windows single-binary delivery path.

The authoritative product surface is `/studio/v2`; `/studio` remains a legacy full-workbench fallback and is not a retirement target. The existing frontend uses `frontend/src/utils/addressParser.ts` and `frontend/src/features/datalink/workbench-v2/state/sourceRule.ts`; the backend persists source rules through `internal/datalink/sourcerule` and workspace source-rule handlers. SQLite is the embedded control database. The current release tree also contains an already-over-500-line legacy `frontend/src/pages/datalink/workbench/sourceCanvasModel.ts`, so the line policy permits only a non-increasing change to that file.

Source-rule ownership is deliberately indirect and single-workspace: each source rule belongs to exactly one device through `device_id`, and that device is the ownership boundary for the workspace projection. The persisted `source_rules` schema has no `workspace_id` and no tombstone. A successful normal deletion of a workspace device therefore removes its device-owned source rules through the existing foreign-key cascade; those rules become unavailable and are not persisted orphan records for later repair. Workspace listing remains scoped to the active workspace device IDs through `ListByDeviceIDs` and must not search globally beyond those IDs. Cleanup or repair of pre-existing orphan rows is explicitly outside this change.

## Goals / Non-Goals

**Goals:**

- Establish one observable protocol-address contract across frontend validation/expansion and backend source-rule validation/offsetting, including MC 3E hexadecimal `X/Y/B` behavior and explicit errors.
- Gate Step 2 continuation across every enabled rule and fail closed for unknown or deleted device references that are actually exposed by the API or retained in stale workspace state, without selecting an arbitrary device. Normal device deletion removes its device-owned rules through the existing cascade.
- Make the SQLite source-rule migration preserve existing and legacy rows, use safe defaults, and recover from a partial or repeated upgrade with a documented rollback boundary.
- Reconcile the predecessor's scope to the approved A surface, preserve unrelated dirty work, and enforce independent staged/unstaged baseline inventories plus the existing line-count policy.
- Define and execute a fail-closed release gate covering normal-worker frontend tests, frontend and Go quality checks, Spectra checks, Windows embedded EXE build, and route/asset smoke.

**Non-Goals:**

- The B-class Share lifecycle is not formalized or extended. Existing Share fields and behavior are compatibility inputs only; this design does not add new Share activation semantics.
- Legacy `/studio` retirement, redesign, or `sourceCanvasModel` behavior changes are excluded. Scope hygiene requires that this change remove accidental legacy drift and does not add a new legacy diff.
- Connector adapters, protocol transport, MQTT behavior, and broad architecture changes are excluded. A protocol address contract consumes existing adapter behavior and does not modify adapter or transport files.
- PostgreSQL control-database live support is not added and is not required for the SQLite migration gate.
- Physical PLC, database-server, GUI, network, and factory-field acceptance are external gates and are not represented as source-test success.

## Decisions

### Reconcile the predecessor and freeze two independent Git inventories

Before implementation, record an immutable fixed point (commit or equivalent tree identity) and separately capture `git diff --cached --name-status`, `git diff --name-status`, and untracked files. The active predecessor's artifact status is not implementation evidence. The apply agent SHALL classify each existing change as predecessor work, approved A work, or unrelated dirty work before editing. The final audit SHALL compare both staged and unstaged inventories to the same fixed point and SHALL leave unrelated paths untouched.

The predecessor's unresolved tasks are reconciled rather than copied. In particular, the new change does not repeat B-class Share lifecycle tasks. A predecessor task is not marked complete merely because a corresponding task exists here; the implementation and focused evidence must be rerun.

### Use a fail-closed protocol address contract

The contract is semantic rather than a shared frontend/backend library. Both sides SHALL normalize protocol and address input, validate the complete address before expansion, and return a field/rule-located validation error on failure. Frontend expansion returns no derived points on invalid input; backend validation/offsetting returns an error and never returns the original invalid string as a successful address.

For MC 3E, `X`, `Y`, and `B` suffixes are hexadecimal and `D`, `W`, and `M` suffixes are decimal. Thus `X0` offset by `16` is `X10`, while `D0` offset by `16` is `D16`. FATEK supported areas and Modbus areas retain their existing radix rules. The protocol-specific formatting rule is selected from the owning device protocol, never from a generic Modbus fallback.

The frontend contract remains compatible with the current `AddressParser` shape: `validate` reports `{valid:false,error}` and `expand` returns an empty list on invalid input. The backend contract preserves `ErrValidation` classification and adds enough rule/address context for the API and UI to identify the failing rule. No connector adapter or protocol transport change is required.

### Make Step 2 continuation an aggregate rule and ownership gate

The continuation decision SHALL evaluate all enabled rules, resolve each rule's `device_id` against the current workspace device set, validate its address with that device protocol, and collect issues keyed by rule ID. An invalid address, missing device, deleted device, or unavailable protocol in a rule returned by the handler/API or retained in stale workspace state blocks continuation even if every other enabled rule is valid. The UI SHALL keep Step 2 active and identify the rule and device/address reason. This defensive check does not require a workspace list endpoint to bypass `OrderedDeviceIDs`/`ListByDeviceIDs` or expose globally stored rules that are outside the workspace projection.

New-rule creation can choose the currently selected device (or the first device only when constructing a brand-new default rule), but hydration, edit, derivation, and continuation SHALL NOT repair an existing missing/unknown `device_id` by selecting the first device or first Modbus device. `deriveAllPoints` and equivalent selectors SHALL preserve this distinction so a stale/inconsistent rule is invalid rather than silently rebound. A normal device deletion is handled by the existing device-to-source-rule foreign-key cascade; it does not require a persisted orphan or a repair search.

### Harden SQLite upgrades without rewriting source-rule data

The migration boundary is schema-only and row-preserving. Existing `source_rules` rows and all existing non-null values SHALL survive an upgrade. New boolean-like fields use a documented safe disabled default; optional numeric fields retain NULL unless an existing value is present. The SQLite path SHALL inspect table/column existence before each DDL operation, so a migration that stopped after any individual column can resume and a second run is a no-op. Tests SHALL seed legacy rows and partial schemas and compare values before and after each run.

The down strategy is explicit: take a copy/backup before applying an upgrade; use the migration down script only when the target SQLite version supports the exact column drop and the operator has confirmed that the added columns contain no required data; otherwise restore the pre-upgrade backup. Down is not an automatic production rollback and it SHALL NOT rewrite or delete legacy rows. PostgreSQL live control-database support is outside this gate; PostgreSQL SQL is syntax-reviewed only where it is already part of the migration bundle.

### Enforce scope hygiene and the legacy line boundary

The allowed application scope is the protocol parser/source-rule validation and Step 2 ownership gate, SQLite migration safety, and release/test gate evidence. The change SHALL NOT modify `internal/datalink/connector/adapters/`, `internal/protocol/`, or unrelated legacy `/studio` behavior. It SHALL remove any predecessor-only diff in `frontend/src/pages/datalink/workbench/sourceCanvasModel.ts`, uncalled exports, and unrelated compression rather than adding a retirement implementation.

The final line audit uses `scripts/check_file_lines.sh` with the repository's 300-line warning and 500-line hard limits. For the historical `sourceCanvasModel.ts`, the final count MUST be less than or equal to the recorded fixed-point count; a net increase is a blocker even though the file is already over 500 lines. No new file in this change exceeds 500 lines, and any necessary test split is behavior-preserving.

### Treat release gates as evidence-producing and fail closed

The release sequence records command, working directory, mode, exit code, duration, and relevant output for each gate. The frontend full suite runs with default/normal workers in non-watch mode (`npm run test -- --run`); a single-worker run is a diagnostic branch only and cannot produce a stability pass. A normal-worker timeout, assertion failure, missing process, or mixed evidence is `blocked`/`unknown`, not a pass.

The required source gates are frontend `npm run lint`, frontend default-worker full `npm run test -- --run`, frontend `npm run build`, `go test ./...`, `go vet ./...`, `golangci-lint run ./...`, the line-limit script, `git diff --check --`, `spectra analyze harden-studio-v2-release-correctness --json` with no Critical/Warning findings, and strict plus normal Spectra validation. The Windows delivery gate runs the complete `scripts/build.ps1` path with PowerShell `-NoProfile`, produces `bin/test-ui.exe`, starts a clean embedded server, and executes the existing Playwright embedded-delivery smoke for `/studio/v2`, `/studio`, `/studio/runtime`, `/test`, and the experimental route/asset graph. Any missing toolchain, blocked child process, or unavailable environment leaves the release gate incomplete.

### Keep source evidence separate from physical acceptance

Source, unit, integration, build, and route/asset smoke evidence prove only the boundaries they exercise. The final report SHALL list physical PLC communication, actual database target, GUI interaction, network reachability, and factory-field acceptance as external checks with their own owner and status. No source gate, single-worker diagnosis, or embedded EXE smoke result is described as PLC or field acceptance.

## Implementation Contract

### Observable behavior

- A valid MC 3E rule expands `X0` through `X10` with hexadecimal boundaries and a `+16` offset; an invalid MC 3E, FATEK, or Modbus address produces an explicit validation issue and no successful fallback point.
- Step 2 remains active when any enabled rule returned by the API or retained in stale workspace state has an invalid address or a missing/deleted device. The issue identifies the rule and, when available, the device and address; valid sibling rules do not mask it. A rule removed by normal device-deletion cascade is unavailable rather than a repairable orphan.
- Existing SQLite source-rule rows retain their values after a first, repeated, and partial upgrade. New optional values receive only the documented safe defaults. A rollback path is explicit and backup-based when destructive down is unsupported.
- The final diff contains only approved hardening paths, has no connector/transport or legacy-retirement work, and does not increase `sourceCanvasModel.ts` from the fixed-point count.
- Release evidence contains separate normal-worker test results, all source gates, embedded EXE build output, and route/asset smoke. An unavailable gate is reported incomplete.

### Interface and data shape

- Frontend address operations retain `AddressParser.validate(address, protocol) -> { valid: boolean, error?: string }` and `AddressParser.expand(address, count, protocol) -> string[]`; invalid expansion returns `[]`.
- Backend source-rule requests continue to use the existing `device_id`, `start_address`, `count`, and protocol-associated device record. Validation errors remain classifiable as `ErrValidation` and include rule/address context at the handler boundary.
- Step 2 readiness issues carry a stable rule identifier and a reason category for `invalid_address`, `unknown_device`, or `deleted_device` when an inconsistent record is exposed to the UI; the existing navigation callback is not invoked while any enabled-rule issue exists. The reason categories do not require the List API to enumerate rules outside `OrderedDeviceIDs`.
- SQLite upgrade semantics retain existing `source_rules` columns and values, use `share_enabled=0` only where that field is already part of the predecessor schema, and leave optional `share_start_register`/`share_stride` NULL when no value exists. This is migration compatibility, not a new Share lifecycle contract.
- Gate evidence records the fixed point, separate staged/unstaged inventories, worker mode, command, exit code, and failure classification.

### Failure modes

- Invalid or unsupported address: return a localized/rule-scoped validation issue; return no derived points; do not echo the raw invalid address as success.
- Unknown or deleted device in an API/stale-workspace record: return an ownership issue; do not choose another device; block continuation and runtime persistence for that stale rule. A normally deleted device's device-owned rules have already become unavailable through FK cascade and are not recovered by a global lookup.
- Partial/repeated SQLite migration: inspect existing columns and continue without overwriting rows; DDL failure stops the migration and preserves the pre-existing data for backup-based recovery.
- Forbidden path, line-count increase, stage/unstaged inventory mismatch, or uncalled export: fail the scope gate.
- Normal-worker full-suite assertion/timeout, toolchain failure, or embedded build/route smoke failure: leave the release status incomplete. A passing single-worker diagnostic does not override this.
- Physical PLC or field failure: report as external acceptance failure; source gates remain accurately recorded.

### Acceptance criteria

- Focused frontend and backend tests cover `X0 + 16 = X10`, invalid address fail-closed behavior, one-valid/one-invalid enabled-rule continuation, and unknown/deleted device ownership for API/stale-workspace records. The per-connection FK contract test at `internal/datalink/db_sqlite_fk_test.go` and backend cascade contract test at `internal/api/handlers/studio_v2_workspace_devices_handler_cascade_test.go` prove that normal device deletion removes its device-owned source rules; latest action-time runs pass with a task-specific writable `GOCACHE`.
- SQLite tests cover a legacy row with non-default values, each partial-column upgrade state, two consecutive migration runs, and backup/down documentation. Tests assert row equality and safe defaults.
- A fixed-point audit proves that forbidden directories, `sourceCanvasModel.ts` net growth, unrelated compression, and uncalled exports are absent; staged and unstaged path inventories are checked separately.
- The normal-worker full Vitest command, frontend lint/build, Go test/vet/lint, line gate, diff-check, Spectra analyze/validate, complete Windows embedded EXE build, and route/asset smoke each have a recorded pass. Any unavailable gate is explicitly incomplete.
- Physical PLC/field acceptance is listed separately and is not inferred from the above evidence.

### Scope boundaries

In scope: the frontend address/Step 2 validation seam, backend source-rule validation/ownership seam, existing SQLite source-rule migration safety, gate scripts/evidence needed to run the specified release checks, and change-local Spectra artifacts.

Out of scope: connector adapters, protocol transport, new Share lifecycle semantics, MQTT, PostgreSQL live control-database support, legacy `/studio` retirement/redesign, unrelated dirty work, broad refactoring, and physical field acceptance.

## Risks / Trade-offs

- **[Risk] Frontend and backend parser rules diverge again.** → **Mitigation:** keep the concrete radix examples in both focused test suites and require the same invalid-address/fail-closed assertions at each boundary.
- **[Risk] Aggregate validation breaks a workspace that previously relied on a first-device fallback.** → **Mitigation:** preserve first-device selection only for a newly constructed default rule and surface an actionable ownership issue for hydrated/stale rules returned by the API; normal device deletion is resolved by the existing source-rule cascade rather than repair lookup.
- **[Risk] SQLite column-drop behavior differs across deployed SQLite versions.** → **Mitigation:** require a backup and an explicit capability check; restore the backup when the exact down operation is unsupported instead of rewriting data.
- **[Risk] Default-worker Vitest remains dynamically unstable.** → **Mitigation:** preserve per-run evidence, classify mixed/timeout results as blocked or unknown, and do not convert a single-worker pass into a stability claim.
- **[Risk] Existing staged work belongs to another agent.** → **Mitigation:** freeze and compare independent staged/unstaged inventories, preserve unrelated paths, and stop implementation when ownership cannot be proven.
- **[Risk] Embedded EXE smoke passes while physical equipment is unavailable.** → **Mitigation:** report route/asset smoke as source-level evidence and leave PLC/field acceptance as a separate external gate.

## Migration Plan

1. Before implementation, capture the fixed point, staged and unstaged inventories, untracked paths, active predecessor status, and the current `sourceCanvasModel.ts` line count.
2. Reconcile the predecessor by separating approved A work from B Share lifecycle work and unrelated dirty changes. Remove only predecessor drift within the authorized ownership boundary.
3. Add focused tests first for protocol radix, invalid-address propagation, aggregate Step 2 continuation, ownership, and SQLite row preservation/partial upgrades; implement the smallest matching source changes.
4. Run focused tests, then the normal-worker full frontend suite and required frontend/Go/Spectra/line/diff gates. Keep single-worker runs in a diagnostic record only.
5. Run the complete Windows embedded build and route/asset smoke from a clean temporary runtime directory. Preserve the generated binary as validation evidence only; do not claim physical acceptance.
6. If migration rollback is required, stop the application, preserve the database backup, use the tested down script only when the SQLite capability check passes, and otherwise restore the backup. Never delete or rewrite legacy rows to force a down migration.

## Open Questions

None. The implementation uses the existing frontend/backend contracts and migration bundle; any unavailable toolchain or external system is a fail-closed validation result, not a request to widen scope.
