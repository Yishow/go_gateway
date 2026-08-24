# embedded-frontend-delivery Specification Delta

## ADDED Requirements

### Requirement: Windows embedded EXE build and route smoke are release gates

The Windows release path SHALL run the complete frontend build, static synchronization, and Go embedded executable build through `scripts/build.ps1` with PowerShell profile isolation. A release candidate SHALL produce `bin/test-ui.exe` and SHALL pass a clean embedded-server route and asset smoke before the embedded delivery gate is complete. The smoke SHALL exercise `/studio/v2`, `/studio`, `/studio/runtime`, `/test`, and the existing experimental route while checking that same-origin JavaScript, module, stylesheet, and lazy route asset requests do not fail.

#### Scenario: Complete Windows build produces a loadable embedded graph

- **GIVEN** the repository has its declared frontend and Go toolchains available on Windows
- **WHEN** the release command runs `scripts/build.ps1` from the repository root
- **THEN** frontend typecheck/build completes
- **AND** the synchronized static directory contains the built entry and lazy assets
- **AND** `bin/test-ui.exe` is produced
- **AND** a clean temporary execution directory can start the binary and serve the embedded entry route

#### Scenario: Route and asset smoke covers the supported surface

- **GIVEN** a freshly built `bin/test-ui.exe` is running on a temporary local port
- **WHEN** the embedded-delivery smoke opens `/studio/v2`, `/studio`, `/studio/runtime`, `/test`, and the existing experimental route
- **THEN** each route renders its declared readiness selector and retains its route identity
- **AND** every same-origin script, module, stylesheet, and lazy route asset request returns successfully
- **AND** a deliberately missing asset returns an explicit 404 without being reported as a successful delivery

#### Scenario: Build or synchronization failure blocks delivery

- **GIVEN** frontend build output is missing, static synchronization fails, Go compilation fails, or the embedded server exits before readiness
- **WHEN** the release gate evaluates the command result
- **THEN** the embedded delivery gate is incomplete
- **AND** the captured command output identifies the failed phase
- **AND** no route-smoke pass is claimed
