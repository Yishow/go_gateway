## Problem

`fix-protocol-address-adaptation-v2` is a predecessor change with useful Step 2–4 work, but its current staged baseline is not release-correct: the frontend and backend do not yet have one fail-closed MC 3E radix contract, Step 2 can be masked by a valid rule, missing device ownership can fall through to an arbitrary device, and the SQLite upgrade/release gates are not proven as a coherent handoff. The predecessor also contains scope and line-limit drift that must be reconciled before this follow-up can be accepted.

## Root Cause

The predecessor mixes protocol address adaptation, Share-related persistence, and unrelated adapter/legacy work. Its frontend `addressParser` and backend source-rule paths need an explicit shared contract for MC 3E `X/Y/B` hexadecimal offsets and invalid-address errors; aggregate Step 2 readiness must inspect every enabled rule; ownership must be established from the workspace device set rather than a first-device fallback. Migration helpers and release scripts currently need evidence for preservation of existing SQLite rows, safe defaults, partial/idempotent upgrades, rollback boundaries, and a complete Windows embedded-binary gate. Staged and unstaged changes are also easy to conflate unless both inventories are captured independently.

## Proposed Solution

Create this release-hardening follow-up after first reconciling the active predecessor. The implementation SHALL:

- Align frontend and backend source-rule address parsing/offsetting: MC 3E `X/Y/B` suffixes use hexadecimal radix (`X0 + 16 = X10`), decimal areas remain decimal, and invalid input returns a rule-located error without returning the original string as a successful address.
- Make Step 2 continuation a workspace-wide gate: every enabled rule must have a valid owning device and address; one invalid rule blocks continuation and identifies that rule even when another rule is valid.
- Make unknown/deleted device references invalid and fail closed; no source rule SHALL silently bind to `devices[0]` or the first Modbus device.
- Reconcile the predecessor diff so connector adapters, protocol transport, legacy `/studio` `sourceCanvasModel`, uncalled exports, unrelated compression, and other unrelated paths are not part of this change. Preserve the existing legacy surface; retirement is out of scope.
- Harden the SQLite migration boundary to preserve existing and legacy source-rule rows, apply safety defaults, tolerate a partial prior upgrade idempotently, and document an explicit down/rollback strategy. PostgreSQL control-database live support is not a release blocker for this A change.
- Define fail-closed release gates: default-worker full Vitest (single-worker runs are diagnostic only), frontend lint/test/build, `go test ./...`, `go vet ./...`, `golangci-lint run ./...`, line-limit and staged/unstaged diff checks, strict Spectra analyze/validate, and a complete Windows embedded EXE build plus route/asset smoke. Physical PLC and factory-field acceptance remain external gates and SHALL NOT be claimed by source gates.

## Non-Goals

- Formalizing or extending the B-class Share lifecycle is not part of this change; existing Share behavior is only retained where needed to avoid regression.
- Retiring or redesigning legacy `/studio`, including `sourceCanvasModel`, is not part of this change. The hardening requirement is scope hygiene and no new legacy diff.
- No connector adapter, protocol transport, MQTT, or cross-module architecture rewrite is authorized.
- No new public API is added merely for future callers, and no unrelated reducer/callback compression is introduced to satisfy line checks.
- PostgreSQL control-database live support is not added and is not treated as a blocker for the SQLite migration contract.
- Passing source tests or an embedded smoke test is not physical PLC, database, GUI, network, or factory-field acceptance.

## Capabilities

### New Capabilities

- `studio-v2-release-correctness`: Release-hardening contracts for fail-closed source-rule validation, migration safety, scope hygiene, and executable delivery gates.

### Modified Capabilities

- `datalink-workbench-v2-step2-rule`: Require protocol-radix correctness, explicit invalid-address errors, aggregate enabled-rule continuation gating, and no arbitrary-device fallback.
- `studio-v2-device-validity`: Require unknown or deleted device ownership to be unavailable and fail closed for source-rule continuation.
- `embedded-frontend-delivery`: Require complete Windows embedded EXE build and route/asset smoke evidence as a release gate.
- `frontend-test-suite-stability`: Require default-worker full Vitest evidence and forbid single-worker success from being reported as stability acceptance.
- `cross-platform-test-contracts`: Require independent staged/unstaged baseline checks and fail-closed classification of unavailable gates.

## Success Criteria

- The resulting implementation has contract tests that demonstrate MC 3E `X0 + 16 = X10`, reject invalid protocol addresses explicitly, and prove that the original invalid string is never returned as a successful point address.
- A Step 2 workspace with one valid enabled rule and one invalid enabled rule remains on Step 2 and identifies the invalid rule; a rule referring to a deleted or unknown device is likewise blocked without selecting another device.
- SQLite migration tests preserve pre-existing and legacy source-rule values, produce only documented safe defaults, succeed when run twice, and fail without data loss when a partial column upgrade is present. The proposal documents the approved down/rollback boundary and does not require PostgreSQL live support.
- The final change has no forbidden predecessor paths or unrelated compression/export drift, does not increase the already-over-500-line `frontend/src/pages/datalink/workbench/sourceCanvasModel.ts`, and verifies staged and unstaged inventories separately from one fixed point.
- Release evidence records default-worker full Vitest, frontend lint/test/build, all required Go gates, line-limit, diff-check, strict Spectra analyze/validate, Windows embedded EXE build, and route/asset smoke. Any toolchain or environment block remains incomplete and is reported as such; physical PLC/field acceptance is listed separately.

## Impact

- Affected specs:
  - `studio-v2-release-correctness`
  - `datalink-workbench-v2-step2-rule`
  - `studio-v2-device-validity`
  - `embedded-frontend-delivery`
  - `frontend-test-suite-stability`
  - `cross-platform-test-contracts`
- Affected code and validation seams:
  - Frontend source-rule/address and Step 2 continuation: `frontend/src/utils/addressParser.ts`, `frontend/src/features/datalink/workbench-v2/state/sourceRule.ts`, `frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx`, `frontend/src/features/datalink/workbench-v2/state/selectors.ts`, and their `frontend/tests/unit/utils/addressParser.test.ts`, `frontend/tests/unit/workbench-v2/sourceRule.test.ts`, and `frontend/tests/unit/workbench-v2/step2-rule.test.tsx` coverage.
  - Backend source-rule validation, ownership, and migration: `internal/datalink/sourcerule/validation.go`, `internal/datalink/sourcerule/service.go`, `internal/datalink/sourcerule/sql_repo.go`, `internal/api/handlers/studio_v2_workspace_source_rules_handler.go`, `internal/api/handlers/source_rule_handler.go`, `internal/datalink/migrator.go`, `internal/datalink/migrator_test.go`, and the existing source-rule migration SQL under `internal/datalink/schema/migrations/`.
  - Release and test gates: `scripts/build.ps1`, `cmd/test_ui/main.go`, `frontend/tests/e2e/embedded-frontend-delivery.spec.ts`, `scripts/diagnose-frontend-vitest-matrix.ps1`, `scripts/lib/FrontendVitest*.psm1`, `scripts/check_file_lines.sh`, `Makefile`, and `.github/workflows/file-line-limit.yml` as required by the implementation and evidence path.
  - Scope/line audit reference: `frontend/src/pages/datalink/workbench/sourceCanvasModel.ts` is inspected for predecessor drift and must not gain lines; it is not a legacy feature-retirement target.
- Explicitly excluded from modification: `internal/datalink/connector/adapters/`, `internal/protocol/`, unrelated legacy `/studio` behavior, and the B-class Share lifecycle beyond compatibility preservation.

This is a bug-fix/hardening follow-up. Before implementation, reconcile the active `fix-protocol-address-adaptation-v2` change and avoid duplicate tasks; unfinished predecessor work is not completion evidence for this change.
