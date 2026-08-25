## MODIFIED Requirements

### Requirement: Guided workflow

The UI SHALL provide a guided workflow that keeps the complete datalink product flow inside /studio/v2 as the canonical user-facing setup workspace: Device -> SourceRule -> Tag review -> Output review/apply. The workflow SHALL remain independent of the removed `/studio` workspace and SHALL NOT make `/studio` a product or compatibility workspace.

The workflow SHALL ensure:

1. Operators establish device capability context before source-rule planning.
2. Tag and output review remain bound to the active source-rule revision.
3. Manual point, tag, or mapping construction tools are secondary and SHALL NOT be required to complete the normal product path.
4. Database and Local Modbus output readiness remain visible in the same workspace.
5. Failures remain localized to the owning step, selection, or output target and use safe localized copy.
6. Step 2 SHALL surface the current active source rule as the primary planning unit, so the operator can immediately understand what is currently planned before moving into Tag review.
7. When the active output target is database, Step 2 SHALL reflect database-aware planning guidance without turning Step 2 into connector or schema setup.
8. A legacy `/studio` bookmark after deletion SHALL follow the existing generic unknown-route policy and SHALL NOT mount a legacy workspace or require query migration.

#### Scenario: End-to-end guided configuration uses v2

- **WHEN** an operator selects a device, saves a source rule, reviews tag candidates, and applies one or both output targets
- **THEN** the primary workflow stays inside /studio/v2
- **AND** the operator does not need a separate manual point-first or mapping-first route to complete the normal path

#### Scenario: Source step exposes the active rule before downstream review

- **WHEN** an operator opens Step 2 with an active source rule
- **THEN** the UI shows the current rule planning summary before the operator moves to Tag review
- **AND** the operator can identify the rule address coverage, type intent, and planning state without reading the full canvas first

#### Scenario: Deleted legacy bookmark does not re-enter legacy workspace

- **WHEN** an operator follows a `/studio` bookmark after deletion
- **THEN** the browser follows the same generic unknown-route policy as an arbitrary unknown path
- **AND** no legacy workspace, query migration, or legacy-only chunk is mounted

### Requirement: Live preview

The UI SHALL provide a live preview of raw and transformed values for a selected mapping using the server preview response and Server-Sent Events (SSE). The preview panel SHALL display raw value, intermediate step results, final transformed value, quality indicator, connection state, last successful timestamp, and actionable error or retry copy. The UI SHALL use idle, connecting, live, reconnecting, degraded, error, and stale states.

#### Scenario: View preview from a valid server response

- **WHEN** a user opens preview and the server returns a valid snapshot
- **THEN** the UI displays raw, step results, final value, quality, and live/last-success context
- **AND** the panel does not claim live before the valid server response is accepted

#### Scenario: Live update via SSE

- **WHEN** the server publishes a valid device value event
- **THEN** the preview panel updates the raw, intermediate, final, quality, and last-success values in real time
- **AND** the panel remains visibly live

#### Scenario: SSE reconnection is visible

- **WHEN** the SSE connection is lost after a valid value
- **THEN** the UI preserves the last value with stale context, enters reconnecting or degraded state, and retries with bounded backoff
- **AND** the UI returns to live only after a valid server event

#### Scenario: Server preview failure is fail-closed

- **WHEN** the server preview returns a typed failure before a valid value
- **THEN** the panel renders an actionable localized error and retry action
- **AND** it does not use local or mock values to present a normal or live preview

### Requirement: Accessible and responsive operator workspace

The UI SHALL remain fully operable by keyboard, preserve i18n compatibility, avoid horizontal overflow at supported desktop breakpoints, and expose truthful accessible state for controls that change server-backed configuration.

#### Scenario: Keyboard-only operation in v2 workbench

- **WHEN** an operator uses keyboard-only navigation in /studio/v2
- **THEN** the operator can switch steps, operate the current toolbar, review selections, activate Share/global controls, and trigger inspector actions
- **AND** all interactive controls expose visible focus and screen-reader-readable state changes

#### Scenario: Share state is announced accessibly

- **WHEN** an operator focuses the Share/global toggle
- **THEN** the control exposes a localized accessible name and aria-pressed state
- **AND** Enter and Space activation works while pending state disables repeat activation

#### Scenario: Desktop responsive stability

- **WHEN** the workbench is rendered at 1280px or above
- **THEN** the desktop shell shows all core workflow regions without horizontal scrolling
- **AND** the main working surface remains usable at 1920x1080 without overflow traps or collapsed critical controls

## ADDED Requirements

### Requirement: Commit progress reports actual result states

The v2 commit progress surface SHALL render pending, success, failed, and skipped states from the backend result model. A failed item SHALL display a safe localized error code/message and an actionable retry or recovery action. The progress surface SHALL NOT replace a failed result with a fixed HTTP 200 or success tone.

#### Scenario: Failed commit item is actionable

- **WHEN** one commit item returns failed with a typed error code
- **THEN** the item displays failed status, localized safe copy, and the retry or recovery action
- **AND** the item does not display 200 as a success result

### Requirement: Preview error diagnostics are separated

The live preview UI SHALL display safe localized operator messages and opaque request identifiers only. Raw backend exception text, stack traces, connection strings, and internal diagnostics SHALL remain outside the normal UI and SHALL be available only through the controlled diagnostics channel.

#### Scenario: Raw preview detail is not rendered

- **WHEN** a preview request fails with a raw backend detail and a typed code
- **THEN** the UI renders the typed code mapping and safe action
- **AND** the raw backend detail is absent from rendered text and DOM content
