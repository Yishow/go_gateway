# studio-v2-release-correctness Specification

## Purpose

Define the auditable Studio V2 release-hardening contract for fail-closed source-rule validation, Step 2 device ownership and readiness, SQLite source-rule migration safety, and embedded frontend delivery. The contract requires bounded scope and independent evidence for frontend, Go, Spectra, and Windows gates, records backup-based recovery when SQLite down is unsupported, and keeps PLC, actual database, GUI, network, and factory-field acceptance as separate external checks.

## Requirements

### Requirement: Release-hardening scope is auditable and bounded

The release-hardening change SHALL contain only the approved A-class source-rule validation, Step 2 readiness, device ownership, SQLite migration safety, and release-gate work. The change SHALL NOT contain connector-adapter files, protocol-transport files, unrelated legacy `/studio` behavior, B-class Share lifecycle formalization, uncalled exports, or unrelated compression. The historical `frontend/src/pages/datalink/workbench/sourceCanvasModel.ts` file SHALL have a final line count less than or equal to its fixed-point count; its existing over-500-line status SHALL NOT be used to justify a net increase or a legacy-retirement implementation.

#### Scenario: Forbidden predecessor path blocks the scope gate

- **GIVEN** the fixed-point audit lists a change under `internal/datalink/connector/adapters/`, `internal/protocol/`, or an unrelated legacy `/studio` path
- **WHEN** the release-hardening path inventory is evaluated
- **THEN** the scope gate fails with the exact forbidden path
- **AND** no release pass is recorded until the path is removed from this change or explicitly assigned to its owning predecessor

#### Scenario: Historical over-limit legacy file does not gain lines

- **GIVEN** `frontend/src/pages/datalink/workbench/sourceCanvasModel.ts` is already over 500 lines at the recorded fixed point
- **WHEN** the final line-limit and diff audit compares the file with that fixed point
- **THEN** an equal or lower line count is accepted for the legacy boundary
- **AND** any net increase fails the gate
- **AND** the audit does not treat this result as authorization to retire or redesign `/studio`

#### Scenario: Uncalled export and unrelated compression are rejected

- **GIVEN** a proposed diff adds an exported symbol with no repository caller or compresses an unrelated reducer/callback solely to change line counts
- **WHEN** the scope audit reviews source references and the fixed-point diff
- **THEN** the change is rejected as out of scope
- **AND** the approved source-rule and release behavior remains otherwise unchanged


<!-- @trace
source: harden-studio-v2-release-correctness
updated: 2026-08-24
-->

---
### Requirement: SQLite source-rule migration preserves legacy data and supports safe recovery

The SQLite source-rule migration SHALL preserve every existing source-rule row and existing non-null value, SHALL assign only documented safe defaults to newly added fields, SHALL resume from any partial column state, and SHALL be idempotent on repeated execution. A down or rollback operation SHALL use an operator-confirmed backup boundary and SHALL NOT rewrite or delete legacy source-rule rows. PostgreSQL control-database live support SHALL NOT be a prerequisite for this SQLite release gate.

#### Scenario: Existing legacy rows survive the first upgrade

- **GIVEN** a legacy SQLite `source_rules` table contains rows with existing address, count, type, and optional values
- **WHEN** the migration adds the source-rule columns required by the predecessor schema
- **THEN** every legacy row remains present with the same pre-existing values
- **AND** newly introduced boolean-like fields use the documented disabled default only when no prior value exists
- **AND** newly introduced optional numeric fields remain NULL when no prior value exists

#### Scenario: A partial upgrade resumes without overwriting data

- **GIVEN** a migration stopped after adding any one subset of the new source-rule columns
- **WHEN** the migration runs again
- **THEN** it detects existing columns, adds only missing columns, and preserves all row values
- **AND** a third execution completes successfully with no additional schema or data change

#### Scenario: Rollback is backup-bounded and does not require PostgreSQL live support

- **GIVEN** an operator requests rollback after a SQLite upgrade
- **WHEN** the target SQLite version supports the tested down operation and a pre-upgrade backup is available
- **THEN** the operator-confirmed down operation removes only the added schema columns and retains legacy rows
- **AND** when the capability check fails, the documented recovery is restoration from the backup
- **AND** PostgreSQL control-database live support is not reported as a blocker for this SQLite gate


<!-- @trace
source: harden-studio-v2-release-correctness
updated: 2026-08-24
code:
  - internal/datalink/migrator.go
tests:
  - internal/datalink/migrator_test.go
  - internal/datalink/migrator_source_rule_partial_test.go
  - internal/datalink/migrator_source_rule_recovery_test.go
-->

---
### Requirement: Release gates produce complete fail-closed evidence

A release-hardening pass SHALL require default-worker full Vitest, frontend lint/test/build, `go test ./...`, `go vet ./...`, `golangci-lint run ./...`, the line-limit gate, independent diff checks, strict Spectra analysis/validation, a complete Windows embedded EXE build, and route/asset smoke. Every gate record SHALL include its command, working directory, execution mode, exit result, and failure classification. A blocked toolchain or external environment SHALL leave the gate incomplete.

#### Scenario: Single-worker success does not close the full-suite gate

- **GIVEN** a single-worker Vitest diagnostic passes while the normal/default-worker full suite times out, asserts, or has mixed evidence
- **WHEN** release evidence is classified
- **THEN** the full-suite gate remains `blocked` or `unknown`
- **AND** the single-worker result is retained as diagnosis-only evidence
- **AND** the release is not reported as stable

#### Scenario: All source gates and embedded routes pass

- **GIVEN** default-worker full Vitest, frontend lint/test/build, all Go gates, line/diff gates, Spectra checks, Windows EXE build, and route/asset smoke each pass
- **WHEN** the source release gate is summarized
- **THEN** the source gate is marked complete with one evidence record per command and route
- **AND** `/studio/v2`, `/studio`, `/studio/runtime`, `/test`, and the existing experimental route retain their declared route identities
- **AND** physical PLC, database, network, GUI, and factory-field acceptance remains listed as a separate external gate

#### Scenario: Build or environment failure remains incomplete

- **GIVEN** the Windows build, child-process startup, route smoke, lint, test, or Go toolchain cannot run or exits non-zero
- **WHEN** the release gate is summarized
- **THEN** the affected gate is marked incomplete with the captured error and environment
- **AND** no aggregate release pass is claimed

<!-- @trace
source: harden-studio-v2-release-correctness
updated: 2026-08-24
code:
  - scripts/build.ps1
  - scripts/diagnose-frontend-vitest-matrix.ps1
  - scripts/lib/FrontendVitestEvidence.psm1
  - frontend/playwright.config.ts
tests:
  - scripts/tests/diagnose-frontend-vitest-matrix.Tests.ps1
  - scripts/tests/diagnose-frontend-vitest-process-result.Tests.ps1
  - frontend/tests/e2e/embedded-frontend-delivery.spec.ts
-->
