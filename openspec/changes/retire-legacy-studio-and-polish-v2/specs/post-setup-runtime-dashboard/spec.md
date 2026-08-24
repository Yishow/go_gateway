## MODIFIED Requirements

### Requirement: Runtime dashboard layers snapshot, live stream, and degraded fallback

The runtime dashboard SHALL load runtime data by fetching a snapshot first, layering valid live stream updates on top, and degrading to snapshot polling when the live stream is unavailable. Snapshot, stream, and device-not-found failures SHALL use the typed safe-error envelope and SHALL render a visible localized state rather than a synthetic ready dashboard.

#### Scenario: Snapshot renders before live stream attaches

- **WHEN** an operator opens /studio/runtime?device_id=device-A
- **THEN** the system first fetches the runtime snapshot for device-A
- **AND** the page renders the latest available runtime summary before live stream events arrive
- **AND** the page displays a loading or waiting state until the first valid snapshot

#### Scenario: Stream disconnect degrades to snapshot polling

- **WHEN** the runtime live stream for device-A disconnects after the initial snapshot succeeds
- **THEN** the system preserves the last successful runtime summary on screen
- **AND** the page enters a visible reconnecting or degraded state that continues refreshing the snapshot until live streaming recovers
- **AND** the page does not retain a live tone while the stream is unavailable

#### Scenario: Typed runtime error is safe and actionable

- **WHEN** the runtime snapshot or stream returns runtime_stream_unavailable, runtime_snapshot_unavailable, or runtime_device_not_found
- **THEN** the page maps the code to localized operator copy and an actionable retry or context action
- **AND** raw backend exception details are not rendered

### Requirement: Post-setup runtime route shows truthful empty and degraded states

The post-setup runtime route SHALL show truthful waiting, empty, degraded, and error states when runtime data is missing or unavailable. The route SHALL use typed error codes rather than matching raw error strings and SHALL NOT present a synthetic ready dashboard just because the route resolved successfully.

#### Scenario: Route opens without real runtime data

- **WHEN** an operator opens the runtime route for a selected device that has no current runtime data yet
- **THEN** the page shows a truthful waiting, empty, or degraded state
- **AND** the page does not present a synthetic ready dashboard just because the route resolved successfully

##### Example: direct link resolves but selected device is still cold

- **GIVEN** the operator opens /studio/runtime for dev-A immediately after activation but runtime has not produced real data
- **WHEN** the route resolves successfully
- **THEN** the page shows waiting or empty state for dev-A instead of a synthetic ready dashboard

#### Scenario: Runtime failure is visible

- **WHEN** the selected device snapshot fails with a typed runtime error
- **THEN** the page shows the localized error code mapping, request identifier, and retry or device-context action
- **AND** it does not show live, ready, or healthy status

## ADDED Requirements

### Requirement: Runtime error envelope is stable

Runtime snapshot and stream endpoints SHALL return an error object containing code, message, retryable, and request_id for failures. The supported codes SHALL include runtime_device_not_found, runtime_snapshot_unavailable, and runtime_stream_unavailable. The runtime UI SHALL map the codes to en and zh-TW copy and SHALL keep raw diagnostics outside the normal dashboard.

#### Scenario: Known runtime failure maps to localized copy

- **WHEN** a runtime endpoint returns code runtime_device_not_found with retryable false
- **THEN** the dashboard shows a localized missing-device-context action
- **AND** the dashboard does not expose the raw backend message

#### Scenario: Unknown runtime failure uses safe fallback

- **WHEN** a runtime endpoint returns an unknown code
- **THEN** the dashboard shows a localized generic failure with the opaque request_id
- **AND** the dashboard does not insert raw response text into the DOM
