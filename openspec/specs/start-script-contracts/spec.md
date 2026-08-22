# start-script-contracts Specification

## Purpose

TBD - created by archiving change 'stabilize-start-script-contracts'. Update Purpose after archive.

## Requirements

### Requirement: Start scripts load in library mode for contract testing

The Windows start script `start.ps1` SHALL support a library-only load mode: when the environment variable `GATEWAY_START_PS1_LIBRARY_ONLY` is set to `1` and the script is dot-sourced, it SHALL define its functions and return before executing any main flow, interactive prompt, or process launch. The bash start script counterpart SHALL retain its source-without-execution property for the shell contract tests.

#### Scenario: Library load does not execute the main flow

- **WHEN** a test sets `GATEWAY_START_PS1_LIBRARY_ONLY=1` and dot-sources `start.ps1`
- **THEN** the orchestration functions (for example `Start-DevMode`) SHALL be defined and no menu, build, or process start SHALL have executed

#### Scenario: Normal invocation ignores the library switch

- **WHEN** `start.ps1` is executed without `GATEWAY_START_PS1_LIBRARY_ONLY`
- **THEN** the script SHALL run its normal interactive main flow unchanged


<!-- @trace
source: stabilize-start-script-contracts
updated: 2026-08-23
code:
  - scripts/start-process-utils.ps1
  - tests/shell/start-backend-cleanup-order.sh
  - tests/shell/start-env-load.sh
  - Makefile
  - scripts/start-port-utils.ps1
  - scripts/run_cross_platform_contract_loop.sh
  - tests/shell/start-backend-before-frontend.sh
  - tests/shell/start-port-management.sh
  - scripts/start-log-utils.sh
  - tests/powershell/start-lock-management.ps1
  - tests/powershell/start-dev-mode-order.ps1
  - tests/powershell/start-log-noise.ps1
  - start.sh
  - tests/shell/start-lock-management.sh
  - start.ps1
  - tests/powershell/start-process-tree.ps1
  - scripts/start-log-utils.ps1
  - .air.toml
  - tests/powershell/start-managed-ports.ps1
  - tests/shell/start-backend-cleanup-no-log-wait.sh
tests:
  - internal/api/handlers/studio_v2_workspace_database_row_groups_test.go
  - cmd/test_ui/static/vite.svg
  - internal/datalink/dbtarget/writer_row_groups_test.go
  - frontend/tests/unit/workbench-v2/helpers/studioV2ServiceMocks.ts
  - frontend/tests/unit/workbench-v2/helpers/workbenchV2PageHarness.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/utils/viteConfig.test.ts
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - cmd/test_ui/static/index.html
-->

---
### Requirement: Dev mode startup order is contractually fixed

`Start-DevMode` in `start.ps1` SHALL orchestrate development mode in the order sync embedded frontend, clear backend port, start backend, wait for backend port readiness, start frontend dev server, wait for backend exit, then cleanup. This order SHALL match the bash counterpart contract verified by `tests/shell/start-backend-before-frontend.sh`, and the PowerShell contract test SHALL assert the same step sequence via function stubs.

#### Scenario: Stubs record the prescribed order

- **WHEN** the dev-mode order contract test overrides the orchestration dependencies with step recorders and invokes `Start-DevMode`
- **THEN** the recorded sequence SHALL equal sync > clear-backend > start-backend > wait-backend > start-frontend > wait-exit > cleanup


<!-- @trace
source: stabilize-start-script-contracts
updated: 2026-08-23
code:
  - scripts/start-process-utils.ps1
  - tests/shell/start-backend-cleanup-order.sh
  - tests/shell/start-env-load.sh
  - Makefile
  - scripts/start-port-utils.ps1
  - scripts/run_cross_platform_contract_loop.sh
  - tests/shell/start-backend-before-frontend.sh
  - tests/shell/start-port-management.sh
  - scripts/start-log-utils.sh
  - tests/powershell/start-lock-management.ps1
  - tests/powershell/start-dev-mode-order.ps1
  - tests/powershell/start-log-noise.ps1
  - start.sh
  - tests/shell/start-lock-management.sh
  - start.ps1
  - tests/powershell/start-process-tree.ps1
  - scripts/start-log-utils.ps1
  - .air.toml
  - tests/powershell/start-managed-ports.ps1
  - tests/shell/start-backend-cleanup-no-log-wait.sh
tests:
  - internal/api/handlers/studio_v2_workspace_database_row_groups_test.go
  - cmd/test_ui/static/vite.svg
  - internal/datalink/dbtarget/writer_row_groups_test.go
  - frontend/tests/unit/workbench-v2/helpers/studioV2ServiceMocks.ts
  - frontend/tests/unit/workbench-v2/helpers/workbenchV2PageHarness.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/utils/viteConfig.test.ts
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - cmd/test_ui/static/index.html
-->

---
### Requirement: Run lock enforces a single start manager per repo and port

`start.ps1` SHALL acquire a run lock implemented as a directory `bin/tmp/start-<port>.lock` containing the owner pid file before managing a repo/port. A second acquisition while the owner process is alive SHALL fail with a message identifying that a start script already manages the repo/port, and the existing lock SHALL remain intact. A lock whose owner process no longer exists SHALL be reported as stale, cleaned up, and reacquired. `Release-RunLock` SHALL remove only the lock owned by the current process.

#### Scenario: Second acquisition fails while the owner is alive

- **WHEN** the lock is held by a live process and another acquisition is attempted
- **THEN** the acquisition SHALL fail, the message SHALL state that a start script already manages the repo/port, and the original pid file SHALL remain in the lock directory

#### Scenario: Stale lock is cleaned and reacquired

- **WHEN** the lock directory exists but its pid no longer maps to a running process
- **THEN** the stale lock SHALL be reported, removed, and the acquisition SHALL succeed

#### Scenario: Release removes only the owned lock

- **WHEN** `Release-RunLock` runs against a lock directory whose pid belongs to another (nonexistent) owner
- **THEN** the lock directory SHALL NOT be removed, while releasing a lock owned by the current process SHALL remove it


<!-- @trace
source: stabilize-start-script-contracts
updated: 2026-08-23
code:
  - scripts/start-process-utils.ps1
  - tests/shell/start-backend-cleanup-order.sh
  - tests/shell/start-env-load.sh
  - Makefile
  - scripts/start-port-utils.ps1
  - scripts/run_cross_platform_contract_loop.sh
  - tests/shell/start-backend-before-frontend.sh
  - tests/shell/start-port-management.sh
  - scripts/start-log-utils.sh
  - tests/powershell/start-lock-management.ps1
  - tests/powershell/start-dev-mode-order.ps1
  - tests/powershell/start-log-noise.ps1
  - start.sh
  - tests/shell/start-lock-management.sh
  - start.ps1
  - tests/powershell/start-process-tree.ps1
  - scripts/start-log-utils.ps1
  - .air.toml
  - tests/powershell/start-managed-ports.ps1
  - tests/shell/start-backend-cleanup-no-log-wait.sh
tests:
  - internal/api/handlers/studio_v2_workspace_database_row_groups_test.go
  - cmd/test_ui/static/vite.svg
  - internal/datalink/dbtarget/writer_row_groups_test.go
  - frontend/tests/unit/workbench-v2/helpers/studioV2ServiceMocks.ts
  - frontend/tests/unit/workbench-v2/helpers/workbenchV2PageHarness.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/utils/viteConfig.test.ts
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - cmd/test_ui/static/index.html
-->

---
### Requirement: Managed port set expands ranges without duplicates

`Get-ManagedPorts` in `scripts/start-port-utils.ps1` SHALL return the deduplicated union of the backend port, the frontend dev port, the Modbus share port, and the fixed range 8080 through 8090, matching the bash managed-ports contract. The expansion SHALL avoid the PowerShell comma-before-range parsing trap, and null or non-numeric port inputs SHALL be skipped without error.

#### Scenario: Default ports expand to the deduplicated managed set

- **WHEN** the backend port is 8080, the frontend dev port is 5173, and the Modbus share port is 5020
- **THEN** the managed set SHALL contain exactly 13 unique ports including 8080, 8090, 5173, and 5020

#### Scenario: Custom backend port joins the managed set

- **WHEN** the backend port is 3333, the frontend dev port is 4173, and the Modbus share port is 15020
- **THEN** the managed set SHALL contain 14 unique ports including 3333, 4173, and 15020 alongside the 8080-8090 range

#### Scenario: Null or invalid inputs are skipped

- **WHEN** the frontend dev port or Modbus share port is null or not numeric
- **THEN** `Get-ManagedPorts` SHALL return the remaining valid ports without throwing


<!-- @trace
source: stabilize-start-script-contracts
updated: 2026-08-23
code:
  - scripts/start-process-utils.ps1
  - tests/shell/start-backend-cleanup-order.sh
  - tests/shell/start-env-load.sh
  - Makefile
  - scripts/start-port-utils.ps1
  - scripts/run_cross_platform_contract_loop.sh
  - tests/shell/start-backend-before-frontend.sh
  - tests/shell/start-port-management.sh
  - scripts/start-log-utils.sh
  - tests/powershell/start-lock-management.ps1
  - tests/powershell/start-dev-mode-order.ps1
  - tests/powershell/start-log-noise.ps1
  - start.sh
  - tests/shell/start-lock-management.sh
  - start.ps1
  - tests/powershell/start-process-tree.ps1
  - scripts/start-log-utils.ps1
  - .air.toml
  - tests/powershell/start-managed-ports.ps1
  - tests/shell/start-backend-cleanup-no-log-wait.sh
tests:
  - internal/api/handlers/studio_v2_workspace_database_row_groups_test.go
  - cmd/test_ui/static/vite.svg
  - internal/datalink/dbtarget/writer_row_groups_test.go
  - frontend/tests/unit/workbench-v2/helpers/studioV2ServiceMocks.ts
  - frontend/tests/unit/workbench-v2/helpers/workbenchV2PageHarness.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/utils/viteConfig.test.ts
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - cmd/test_ui/static/index.html
-->

---
### Requirement: Process tree termination and bounded port readiness

`Stop-ProcessTree` in `scripts/start-process-utils.ps1` SHALL terminate a process and its entire child process tree without leaving orphan processes. `Wait-PortReady` SHALL succeed when the target port is listening, SHALL fail with a timeout message when it is not, and SHALL fail fast without consuming all attempts when the backend process has already exited.

#### Scenario: Child processes die with the parent

- **WHEN** a parent shell process spawns a long-running child and `Stop-ProcessTree` terminates the parent
- **THEN** neither the parent nor the child process SHALL remain observable

#### Scenario: Port readiness outcomes are bounded

- **WHEN** the target port is listening and `Wait-PortReady` is invoked
- **THEN** it SHALL return success within its configured attempts

- **WHEN** the target port never opens
- **THEN** `Wait-PortReady` SHALL return failure with a message stating startup did not succeed in time

- **WHEN** the backend process has already exited while waiting for the port
- **THEN** `Wait-PortReady` SHALL return failure promptly without consuming all attempts, with a message stating the backend process exited


<!-- @trace
source: stabilize-start-script-contracts
updated: 2026-08-23
code:
  - scripts/start-process-utils.ps1
  - tests/shell/start-backend-cleanup-order.sh
  - tests/shell/start-env-load.sh
  - Makefile
  - scripts/start-port-utils.ps1
  - scripts/run_cross_platform_contract_loop.sh
  - tests/shell/start-backend-before-frontend.sh
  - tests/shell/start-port-management.sh
  - scripts/start-log-utils.sh
  - tests/powershell/start-lock-management.ps1
  - tests/powershell/start-dev-mode-order.ps1
  - tests/powershell/start-log-noise.ps1
  - start.sh
  - tests/shell/start-lock-management.sh
  - start.ps1
  - tests/powershell/start-process-tree.ps1
  - scripts/start-log-utils.ps1
  - .air.toml
  - tests/powershell/start-managed-ports.ps1
  - tests/shell/start-backend-cleanup-no-log-wait.sh
tests:
  - internal/api/handlers/studio_v2_workspace_database_row_groups_test.go
  - cmd/test_ui/static/vite.svg
  - internal/datalink/dbtarget/writer_row_groups_test.go
  - frontend/tests/unit/workbench-v2/helpers/studioV2ServiceMocks.ts
  - frontend/tests/unit/workbench-v2/helpers/workbenchV2PageHarness.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/utils/viteConfig.test.ts
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - cmd/test_ui/static/index.html
-->

---
### Requirement: Runtime log noise is reduced with counters and summary

The Windows log utility `scripts/start-log-utils.ps1` SHALL classify backend runtime log lines the same way as `scripts/start-log-utils.sh`: dashboard polling requests (GET 200 on the datalink entity endpoints) SHALL be counted under a noise counter and not emitted; 5xx responses SHALL be emitted in red with an HTTP-formatted line; non-Go-format error lines SHALL be emitted in red; Go-format log lines SHALL be silent unless verbose mode is enabled; verbose mode SHALL emit original lines without reduction; and the summary output SHALL list per-category counts. The datalink entity list SHALL come from a single shared source per platform.

#### Scenario: Dashboard polling is counted, not printed

- **WHEN** a runtime log line records GET on a datalink entity endpoint with status 200
- **THEN** the dashboard-refresh counter SHALL increment and the line SHALL NOT be emitted

#### Scenario: 5xx and error lines stay visible in red

- **WHEN** a runtime log line records a 5xx response or a non-Go-format error
- **THEN** the line SHALL be emitted in red, with 5xx lines in HTTP format

#### Scenario: Verbose mode disables reduction

- **WHEN** verbose mode is enabled
- **THEN** runtime log lines SHALL be emitted in original form without counting-based suppression


<!-- @trace
source: stabilize-start-script-contracts
updated: 2026-08-23
code:
  - scripts/start-process-utils.ps1
  - tests/shell/start-backend-cleanup-order.sh
  - tests/shell/start-env-load.sh
  - Makefile
  - scripts/start-port-utils.ps1
  - scripts/run_cross_platform_contract_loop.sh
  - tests/shell/start-backend-before-frontend.sh
  - tests/shell/start-port-management.sh
  - scripts/start-log-utils.sh
  - tests/powershell/start-lock-management.ps1
  - tests/powershell/start-dev-mode-order.ps1
  - tests/powershell/start-log-noise.ps1
  - start.sh
  - tests/shell/start-lock-management.sh
  - start.ps1
  - tests/powershell/start-process-tree.ps1
  - scripts/start-log-utils.ps1
  - .air.toml
  - tests/powershell/start-managed-ports.ps1
  - tests/shell/start-backend-cleanup-no-log-wait.sh
tests:
  - internal/api/handlers/studio_v2_workspace_database_row_groups_test.go
  - cmd/test_ui/static/vite.svg
  - internal/datalink/dbtarget/writer_row_groups_test.go
  - frontend/tests/unit/workbench-v2/helpers/studioV2ServiceMocks.ts
  - frontend/tests/unit/workbench-v2/helpers/workbenchV2PageHarness.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/utils/viteConfig.test.ts
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - cmd/test_ui/static/index.html
-->

---
### Requirement: Windows Air hot-reload uses an exe entrypoint

The Air configuration SHALL build the hot-reload binary to `./bin/tmp/gateway-air.exe` and set the matching `entrypoint`, and the bash start script SHALL reference the same exe path, so that Air hot-reload works on Windows while remaining functional on Linux.

#### Scenario: Air rebuilds and restarts the exe entrypoint on Windows

- **WHEN** Air runs on Windows with the repository configuration and a watched Go file changes
- **THEN** Air SHALL build to `./bin/tmp/gateway-air.exe` and successfully start the rebuilt binary

<!-- @trace
source: stabilize-start-script-contracts
updated: 2026-08-23
code:
  - scripts/start-process-utils.ps1
  - tests/shell/start-backend-cleanup-order.sh
  - tests/shell/start-env-load.sh
  - Makefile
  - scripts/start-port-utils.ps1
  - scripts/run_cross_platform_contract_loop.sh
  - tests/shell/start-backend-before-frontend.sh
  - tests/shell/start-port-management.sh
  - scripts/start-log-utils.sh
  - tests/powershell/start-lock-management.ps1
  - tests/powershell/start-dev-mode-order.ps1
  - tests/powershell/start-log-noise.ps1
  - start.sh
  - tests/shell/start-lock-management.sh
  - start.ps1
  - tests/powershell/start-process-tree.ps1
  - scripts/start-log-utils.ps1
  - .air.toml
  - tests/powershell/start-managed-ports.ps1
  - tests/shell/start-backend-cleanup-no-log-wait.sh
tests:
  - internal/api/handlers/studio_v2_workspace_database_row_groups_test.go
  - cmd/test_ui/static/vite.svg
  - internal/datalink/dbtarget/writer_row_groups_test.go
  - frontend/tests/unit/workbench-v2/helpers/studioV2ServiceMocks.ts
  - frontend/tests/unit/workbench-v2/helpers/workbenchV2PageHarness.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/utils/viteConfig.test.ts
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - cmd/test_ui/static/index.html
-->