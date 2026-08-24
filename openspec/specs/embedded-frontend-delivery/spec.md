# embedded-frontend-delivery Specification

## Purpose

Define the delivery contract for the React frontend when it is embedded in the Go binary: the Vite build must produce a complete relative asset graph, each supported build/start path must synchronize that graph to `cmd/test_ui/static` and fail closed on missing or failed inputs, and route-level lazy loading must preserve existing route semantics. The contract is verified at source, synchronization, fresh-clone, and browser boundaries, including deterministic loading/error behavior and cross-platform Playwright/E2E coverage.

## Requirements

### Requirement: Frontend build emits a complete embedded asset graph

The frontend build SHALL emit an entry document at frontend/dist/index.html and every JavaScript, CSS, route chunk, vendor chunk, and referenced static asset required by its direct or transitive import graph. Asset references SHALL remain relative to the embedded static root. Vendor grouping SHALL use stable logical groups for the React/runtime, Material UI and Emotion, charting, router/i18n/query application dependencies, and residual dependencies.

#### Scenario: Build emits entry and lazy route assets

- **WHEN** the repository frontend build runs with its declared dependencies
- **THEN** frontend/dist/index.html exists
- **AND** the output contains the lazy assets needed by the existing /studio/v2, /studio, /studio/runtime, /test, and /gateway/* route families
- **AND** every asset reference discovered from the entry and lazy import graph resolves to a file under frontend/dist

##### Example: vendor group names

- **GIVEN** a production build includes React, Material UI, charting, and router dependencies
- **WHEN** Vite writes the production output
- **THEN** the output contains stable logical vendor groups named react-vendor, mui-vendor, charts-vendor, app-vendor, or misc-vendor for the matching dependency families
- **AND** application route chunks remain separate from those vendor groups

<!-- @trace
source: reliable-embedded-frontend-delivery
updated: 2026-08-22
code:
  - frontend/vite.config.ts
  - frontend/src/App.tsx
tests:
  - frontend/tests/unit/app-routing-lazy-load.test.tsx
  - tests/shell/embedded-frontend-delivery.sh
-->

---
### Requirement: Embedded static synchronization mirrors the complete frontend output

Each supported synchronization path SHALL ensure cmd/test_ui/static exists, SHALL remove generated children that are absent from the current frontend/dist while preserving the tracked embed-placeholder.txt, and SHALL copy the complete frontend/dist tree including nested directories and all asset types. Generated children under cmd/test_ui/static other than embed-placeholder.txt SHALL remain ignored by Git.

#### Scenario: Synchronization removes stale output and preserves the placeholder

- **GIVEN** cmd/test_ui/static contains embed-placeholder.txt, current assets, and a stale generated file that is absent from frontend/dist
- **WHEN** scripts/build.ps1, start.ps1 Build-Frontend, or start.sh build_frontend completes a successful synchronization
- **THEN** embed-placeholder.txt still exists
- **AND** the stale generated file no longer exists
- **AND** every file and directory in frontend/dist exists at the matching relative path under cmd/test_ui/static
- **AND** no generated static child outside the current frontend/dist tree remains

#### Scenario: Fresh clone provides the embed directory marker

- **GIVEN** a fresh checkout has the tracked cmd/test_ui/static/embed-placeholder.txt and has no generated frontend assets
- **WHEN** the Go embedded frontend package is compiled
- **THEN** the compiler finds a match for the static embed pattern without requiring frontend/dist
- **AND** no generated asset is required to make the checkout compile

<!-- @trace
source: reliable-embedded-frontend-delivery
updated: 2026-08-22
code:
  - .gitignore
  - cmd/test_ui/static/embed-placeholder.txt
  - scripts/build.ps1
  - start.ps1
  - start.sh
tests:
  - tests/shell/embedded-frontend-delivery.sh
-->

---
### Requirement: Build and synchronization failures are surfaced and fail closed

The frontend build and synchronization commands SHALL return a nonzero exit status when dependency installation, frontend source availability, frontend/dist generation, target-directory creation, stale-asset cleanup, or asset copy fails. A failed frontend build, missing frontend source, or missing frontend/dist SHALL NOT be followed by a backend build that is reported as a successful frontend delivery.

#### Scenario: Missing dist stops synchronization

- **GIVEN** frontend/dist does not exist after the frontend build step
- **WHEN** a supported build or start path attempts to synchronize embedded assets
- **THEN** the command returns a nonzero exit status
- **AND** the output identifies frontend/dist as unavailable
- **AND** the backend build step is not reported as successful for that invocation

#### Scenario: Missing frontend source stops the POSIX build path

- **GIVEN** the `frontend` source directory or its required build inputs are absent before `start.sh` runs its frontend build path
- **WHEN** the POSIX build path attempts to build and then invoke the backend compiler
- **THEN** the command returns a nonzero exit status
- **AND** the output identifies the missing frontend source
- **AND** the backend compiler is not invoked in that invocation

#### Scenario: Copy failure remains actionable

- **GIVEN** the target static directory cannot be cleaned or written
- **WHEN** a supported synchronization path performs the copy
- **THEN** the command returns a nonzero exit status
- **AND** the existing console or log channel contains the copy failure
- **AND** the command does not claim that embedded frontend delivery completed

<!-- @trace
source: reliable-embedded-frontend-delivery
updated: 2026-08-22
code:
  - scripts/build.ps1
  - start.ps1
  - start.sh
tests:
  - tests/shell/embedded-frontend-delivery.sh
-->

---
### Requirement: Lazy route loading has one deterministic loading and error contract

The application SHALL load heavy existing route modules through route-level dynamic imports and SHALL render one deterministic Suspense fallback while the selected route chunk is pending. A rejected route or vendor import SHALL reach the existing error boundary as an explicit error state. The implementation MUST preserve each route component's existing props and route parameters.

#### Scenario: Core route waits on the deterministic fallback

- **GIVEN** the application has started and a core route chunk has not finished loading
- **WHEN** an operator navigates to /studio/v2, /studio, /studio/runtime, /test, or a /gateway/* route
- **THEN** the route container renders the same deterministic loading fallback while the import is pending
- **AND** the selected route renders after its lazy chunk and vendor dependencies load
- **AND** no second product route or redirect is introduced

#### Scenario: Missing route chunk is explicit

- **GIVEN** a selected route chunk or vendor chunk is unavailable
- **WHEN** the browser evaluates its dynamic import
- **THEN** the existing error boundary renders an error state
- **AND** the application does not present the route as successfully loaded
- **AND** the failed asset request remains observable as a missing-resource response

<!-- @trace
source: reliable-embedded-frontend-delivery
updated: 2026-08-22
code:
  - frontend/src/App.tsx
tests:
  - frontend/tests/unit/app-routing-lazy-load.test.tsx
-->

---
### Requirement: Existing route product semantics remain unchanged during delivery refactoring

The embedded delivery implementation SHALL retain the current route contract: the default entry resolves to /studio/v2, /studio remains the full-workbench fallback, /studio/runtime remains the focused monitor, /test remains an independent engineering tool, and /gateway/* remains experimental. Lazy loading and chunk grouping SHALL NOT alter redirects, route parameters, or the identity of the existing route components.

#### Scenario: Default and fallback routes retain their identities

- **WHEN** an operator opens the root entry, generic datalink landing, /studio, and /studio/v2
- **THEN** root and generic landing resolve according to the existing /studio/v2 redirect contract
- **AND** /studio still renders the full-workbench fallback
- **AND** /studio/v2 still renders the guided user-facing entry
- **AND** all route asset requests return successful responses

#### Scenario: Monitor, tool, and experimental routes remain distinct

- **WHEN** an operator opens /studio/runtime with its existing device context, /test, and one /gateway/* route
- **THEN** /studio/runtime renders the focused runtime monitor contract
- **AND** /test remains an independent engineering tool
- **AND** /gateway/* remains an experimental surface
- **AND** lazy loading introduces no redirect to a different product flow

##### Example: Route family identity checks

- **GIVEN** the runtime URL includes device_id=device-A and the gateway URL is /gateway/quick-setup
- **WHEN** the operator opens /studio/runtime?device_id=device-A, /test, and /gateway/quick-setup
- **THEN** the runtime page remains the focused monitor for device-A
- **AND** /test remains the independent engineering tool
- **AND** /gateway/quick-setup remains an experimental surface without a redirect to /studio/v2

<!-- @trace
source: reliable-embedded-frontend-delivery
updated: 2026-08-22
code:
  - frontend/src/App.tsx
tests:
  - frontend/tests/unit/app-routing-lazy-load.test.tsx
  - frontend/tests/e2e/embedded-frontend-delivery.spec.ts
-->

---
### Requirement: Embedded browser smoke proves the full asset graph is loadable

A built embedded binary served with a synchronized static tree SHALL load the existing core route families without a JavaScript chunk or vendor asset returning HTTP 404. The smoke verification SHALL check both network responses and the route identity or redirect assertions for /studio/v2, /studio, /studio/runtime, /test, and /gateway/*.

#### Scenario: Embedded binary loads all core route families

- **GIVEN** a binary built from a fresh checkout after frontend assets are synchronized
- **WHEN** browser smoke opens the five route families and waits for lazy imports to settle
- **THEN** every JavaScript, CSS, and route asset request returns a non-404 response
- **AND** each route identity or redirect matches the existing product contract
- **AND** the smoke check exits successfully

#### Scenario: Stale or missing asset fails the smoke check

- **GIVEN** one referenced route or vendor asset is absent from cmd/test_ui/static
- **WHEN** browser smoke opens the route that imports the missing asset
- **THEN** the check reports the failed asset request
- **AND** the check exits nonzero
- **AND** the missing-asset condition is not hidden by the loading fallback

<!-- @trace
source: reliable-embedded-frontend-delivery
updated: 2026-08-22
code:
  - frontend/playwright.config.ts
tests:
  - frontend/tests/e2e/embedded-frontend-delivery.spec.ts
-->

---
### Requirement: Playwright webServer startup is cross-platform

The frontend Playwright configuration SHALL provide a webServer command and environment setup that the exact command `npm run test:e2e -- tests/e2e/embedded-frontend-delivery.spec.ts`, launched from the `frontend` directory, can start on Windows and POSIX hosts. The configuration SHALL wait for its declared readiness URL before route assertions begin and SHALL NOT depend on a POSIX-only inline environment assignment or shell-chaining syntax.

#### Scenario: Windows direct invocation starts the configured webServer

- **GIVEN** a Windows checkout with frontend dependencies installed
- **WHEN** the operator runs `npm run test:e2e -- tests/e2e/embedded-frontend-delivery.spec.ts` from `frontend`
- **THEN** Playwright starts the configured webServer without a shell syntax error
- **AND** Playwright waits for the declared readiness URL
- **AND** the embedded frontend smoke proceeds to its route assertions

#### Scenario: POSIX direct invocation starts the configured webServer

- **GIVEN** a POSIX checkout with frontend dependencies installed
- **WHEN** the operator runs `npm run test:e2e -- tests/e2e/embedded-frontend-delivery.spec.ts` from `frontend`
- **THEN** Playwright starts the same configured webServer
- **AND** Playwright waits for the declared readiness URL
- **AND** the embedded frontend smoke proceeds to its route assertions

<!-- @trace
source: reliable-embedded-frontend-delivery
updated: 2026-08-22
code:
  - frontend/playwright.config.ts
tests:
  - frontend/tests/e2e/embedded-frontend-delivery.spec.ts
-->

---
### Requirement: Historical launcher line counts do not regress

The implementation SHALL keep the line count of `start.ps1` and `start.sh` less than or equal to each file's pre-change baseline. The final gate SHALL fail when either historical launcher grows beyond its baseline, while preserving the synchronization and fail-closed behavior requirements.

#### Scenario: Final gate rejects launcher growth

- **GIVEN** pre-change line-count baselines are recorded for `start.ps1` and `start.sh`
- **WHEN** the final content and line-count gate evaluates the repaired launchers
- **THEN** each current line count is less than or equal to its corresponding baseline
- **AND** `scripts/check_file_lines.sh` reports no new hard-limit violation
- **AND** the gate returns nonzero if either current count exceeds its baseline

<!-- @trace
source: reliable-embedded-frontend-delivery
updated: 2026-08-22
code:
  - start.ps1
  - start.sh
  - scripts/check_file_lines.sh
tests:
  - tests/shell/embedded-frontend-delivery.sh
-->

---
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

<!-- @trace
source: harden-studio-v2-release-correctness
updated: 2026-08-24
code:
  - scripts/build.ps1
  - frontend/playwright.config.ts
tests:
  - frontend/tests/e2e/embedded-frontend-delivery.spec.ts
-->
