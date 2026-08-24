# cross-platform-test-contracts Specification

## Purpose

本規格定義跨 Windows/Linux、不同環境與平行度下的測試契約，確保環境隔離、JSON 編碼、平台中立錯誤與非同步結果驗證不因測試實作細節產生誤報。

## Requirements

### Requirement: Vite tests isolate ambient environment while preserving production semantics

The Vite configuration test in `frontend/tests/unit/utils/viteConfig.test.ts` SHALL isolate ambient `.env` and process environment state for each case while invoking the existing production `frontend/vite.config.ts` config contract. The test SHALL preserve verification of `PORT` fallback, `VITE_API_PROXY_TARGET` precedence, and `VITE_DEV_PORT` parsing.

#### Scenario: Ambient environment does not override an explicit proxy test

- **WHEN** the test sets `PORT` and `VITE_API_PROXY_TARGET` and loads the production config in a controlled environment
- **THEN** the `/api` proxy target SHALL equal the explicit `VITE_API_PROXY_TARGET` value regardless of repository `.env` contents

#### Scenario: Environment state is restored between cases

- **WHEN** a Vite config test case completes
- **THEN** the original values or absence of `PORT`, `VITE_API_PROXY_TARGET`, and `VITE_DEV_PORT` SHALL be restored before the next case


<!-- @trace
source: stabilize-cross-platform-test-contracts
updated: 2026-08-22
tests:
  - frontend/tests/unit/utils/viteConfig.test.ts
-->

---
### Requirement: Studio V2 database requests use JSON encoding for platform paths

The Studio V2 database handler tests in `internal/api/handlers/studio_v2_workspace_database_handler_test.go`, `internal/api/handlers/studio_v2_workspace_database_row_groups_test.go`, `internal/api/handlers/studio_v2_workspace_audit_handler_test.go`, and `internal/api/handlers/studio_v2_workspace_database_delivery_truth_test.go` SHALL construct request payloads with `encoding/json` or a shared typed encoder. They SHALL NOT form JSON by interpolating filesystem paths, point IDs, column names, or row-group members into raw string literals.

#### Scenario: Windows temporary SQLite paths remain valid JSON

- **WHEN** a database configuration request includes a Windows temporary path containing backslashes
- **THEN** the handler SHALL parse the request and persist the intended path without an invalid JSON escape or altered path value

#### Scenario: Row-group member values survive encoding

- **WHEN** a row-group configuration request includes the fixture point IDs and database path
- **THEN** the handler response SHALL contain the same point IDs and database configuration values that were supplied to the typed JSON encoder


<!-- @trace
source: stabilize-cross-platform-test-contracts
updated: 2026-08-22
tests:
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - internal/api/handlers/studio_v2_workspace_database_row_groups_test.go
  - internal/api/handlers/studio_v2_workspace_audit_handler_test.go
  - internal/api/handlers/studio_v2_workspace_database_delivery_truth_test.go
-->

---
### Requirement: Device connection failure assertions are platform-independent

The connection failure test in `internal/api/handlers/device_handler_extended_test.go` SHALL assert stable response state rather than an operating-system-specific error phrase. The test SHALL preserve the existing production response and SHALL verify failed connection state, disabled activation and collection, and skipped probe state.

#### Scenario: Closed listener produces a failed connection state

- **WHEN** a draft connection targets a localhost listener that has been closed
- **THEN** the response SHALL report `data.success` as false, `can_activate` and `can_collect` as false, `connect.status` as `failed`, and `probe.status` as `skipped`

#### Scenario: Error wording differs by operating system

- **WHEN** the operating system returns a platform-specific connection error string
- **THEN** the test SHALL require a non-empty connection error or an existing stable classification and SHALL NOT require the literal phrase `connection refused`


<!-- @trace
source: stabilize-cross-platform-test-contracts
updated: 2026-08-22
tests:
  - internal/api/handlers/device_handler_extended_test.go
-->

---
### Requirement: Row-group tests verify unordered outcomes with bounded flushing

The row-group writer test in `internal/datalink/dbtarget/writer_row_groups_test.go` SHALL use a deterministic, bounded flush signal and SHALL verify row-group outcomes independently of map insertion order. The test SHALL NOT change production writer ordering or require a particular SQLite auto-increment order.

#### Scenario: Two shared-column groups emit both values

- **WHEN** values `21.5` and `23.75` are written to two row groups and one controlled flush is triggered
- **THEN** exactly two rows SHALL become observable within the bounded wait, and the unordered result SHALL contain both `21.5` and `23.75`

#### Scenario: Flush does not complete in the bounded interval

- **WHEN** the expected row count is not observable before the configured timeout
- **THEN** the test SHALL fail with the observed count or query error and SHALL terminate without an unbounded wait


<!-- @trace
source: stabilize-cross-platform-test-contracts
updated: 2026-08-22
tests:
  - internal/datalink/dbtarget/writer_row_groups_test.go
-->

---
### Requirement: Regression loops prove cross-platform test stability

The change SHALL provide a repeatable validation sequence that runs the Vite focused test twice, the relevant handler package focused and full tests, the row-group test with `-count=20` under normal GOMAXPROCS, the frontend full suite, and `go test ./...`. The sequence SHALL distinguish failures in these scopes from unrelated baseline gates.

#### Scenario: Focused fixes pass before broad suites

- **WHEN** the focused frontend, handler, and row-group loops are executed in the prescribed order
- **THEN** each loop SHALL pass before the frontend full suite and Go full suite are started

#### Scenario: Existing baseline failures remain outside this contract

- **WHEN** frontend full suite or `go test ./...` reports a failure unrelated to the four contracts
- **THEN** the report SHALL identify the failing scope and SHALL NOT attribute the failure to this capability without reproducing one of its scenarios

<!-- @trace
source: stabilize-cross-platform-test-contracts
updated: 2026-08-22
tests:
  - frontend/tests/unit/utils/viteConfig.test.ts
  - internal/api/handlers/device_handler_extended_test.go
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - internal/api/handlers/studio_v2_workspace_database_row_groups_test.go
  - internal/datalink/dbtarget/writer_row_groups_test.go
-->

---
### Requirement: Staged and unstaged release baselines are checked independently

The release audit SHALL record one fixed point and SHALL inspect staged, unstaged, and untracked path inventories as separate datasets. A combined diff or a single clean working-tree result SHALL NOT substitute for the three inventories. The audit SHALL identify ownership conflicts and forbidden paths before any gate result is accepted.

#### Scenario: Independent inventories expose a staged-only path

- **GIVEN** a file exists only in `git diff --cached --name-status` and not in `git diff --name-status`
- **WHEN** the fixed-point release audit runs
- **THEN** the staged dataset records that file independently
- **AND** the unstaged dataset remains separately recorded
- **AND** the audit does not silently treat the file as an unstaged change or as absent

#### Scenario: Inventory mismatch blocks acceptance

- **GIVEN** the active predecessor or another worker changes a path after the fixed point
- **WHEN** the final staged/unstaged comparison no longer matches the recorded ownership inventory
- **THEN** the audit reports the path and mismatch category
- **AND** the release-hardening acceptance remains blocked until the owning agent reconciles it


<!-- @trace
source: harden-studio-v2-release-correctness
updated: 2026-08-24
-->

---
### Requirement: Cross-platform and environment failures are classified fail closed

The release evidence SHALL distinguish source defects, assertion failures, process or worker timeouts, missing toolchains, permission failures, and unavailable external systems. A mixed or insufficient signal SHALL be classified as `blocked` or `unknown`; it SHALL NOT be converted into a pass by changing worker count, rerunning without evidence, or relying on a different platform result.

#### Scenario: Toolchain unavailability blocks the gate

- **GIVEN** a required Windows, Go, Node, Playwright, or Spectra tool is unavailable or cannot start
- **WHEN** the gate runner records the result
- **THEN** the affected gate is marked incomplete with the command and environment error
- **AND** a pass is not inferred from a focused or different-platform run

#### Scenario: Mixed assertion and timeout evidence remains unknown

- **GIVEN** one run reports assertion failures and an external watchdog reports a timeout or process cleanup uncertainty
- **WHEN** the stability matrix classifies the run
- **THEN** the result is `blocked` or `unknown`
- **AND** a single-worker pass or a later partial rerun does not replace the original evidence

<!-- @trace
source: harden-studio-v2-release-correctness
updated: 2026-08-24
code:
  - scripts/diagnose-frontend-vitest-matrix.ps1
  - scripts/lib/FrontendVitestEvidence.psm1
  - scripts/lib/FrontendVitestProcess.psm1
tests:
  - scripts/tests/diagnose-frontend-vitest-matrix.Tests.ps1
  - scripts/tests/diagnose-frontend-vitest-process-result.Tests.ps1
-->
