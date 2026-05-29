## ADDED Requirements

### Requirement: Focused runtime dashboard surface
The system SHALL render a focused post-setup runtime dashboard surface for one selected device, without replacing the `/studio/v2` setup shell or introducing a fleet-wide card wall.

#### Scenario: Render focused device monitoring page
- **WHEN** the runtime dashboard receives a selected device context and a successful runtime snapshot
- **THEN** the page renders a focused device header, runtime summary, collector health panel, and live point values table
- **AND** the page does not require queue backlog, diagnostic logs, or a global navigation card wall to become usable

##### Example: Focused page for one committed device
- **GIVEN** the selected device is `device-A` and the runtime snapshot for `device-A` loaded successfully
- **WHEN** the runtime dashboard page renders
- **THEN** the operator can read one page centered on `device-A` instead of a multi-device overview

### Requirement: Runtime summary presents backend-supported metrics
The runtime dashboard SHALL present only metrics and summary fields guaranteed by the runtime backend contract.

#### Scenario: Render summary from runtime snapshot
- **WHEN** the runtime dashboard receives a runtime snapshot for the selected device
- **THEN** the summary panel shows the backend-supported runtime metrics and collector summary values
- **AND** the page does not infer unsupported backlog or log counters from unrelated endpoints

##### Example: Summary uses only formal contract fields
- **GIVEN** the runtime snapshot includes `collected_total=128`, `write_success_total=128`, `write_error_total=0`, and collector `status=healthy`
- **WHEN** the summary panel renders
- **THEN** it shows those runtime metrics and collector summary values without inventing queue backlog or log counters

### Requirement: Collector health panel
The runtime dashboard SHALL show collector health for the selected device using the runtime status fields defined by the backend contract.

#### Scenario: Show collector health counters and status
- **WHEN** the runtime snapshot or aligned status update reports `status`, `points_total`, `points_healthy`, `points_stale`, `points_error`, `last_read_at`, and `breaker_state`
- **THEN** the collector health panel renders those fields for the selected device
- **AND** the panel uses the same status vocabulary as the backend contract instead of introducing page-only health states

### Requirement: Live point values table
The runtime dashboard SHALL present a live point values table for the selected device, including raw values, transformed values, quality or stale state, and timestamps.

#### Scenario: Show live values after snapshot-first load
- **WHEN** the selected device snapshot has loaded but live point events have not arrived yet
- **THEN** the page keeps the runtime summary visible
- **AND** the live point values table shows a waiting or placeholder state instead of disappearing

##### Example: Snapshot visible before first live point event
- **GIVEN** `device-A` snapshot loaded at `10:00:00Z` and no live point events have arrived yet
- **WHEN** the page enters its initial post-snapshot state
- **THEN** the runtime summary remains visible and the live values area shows a waiting placeholder

#### Scenario: Update table with live values
- **WHEN** live point value updates arrive for the selected device
- **THEN** the table updates the corresponding rows with raw value, transformed value, quality or stale state, and timestamp

##### Example: Raw and transformed values update one row
- **GIVEN** the table already contains point `pt-1` for `device-A`
- **WHEN** a live update arrives with `raw_value=150`, `transformed_value=20.0`, `quality=good`, and timestamp `2026-05-29T10:00:05Z`
- **THEN** the `pt-1` row shows `150`, `20.0`, `good`, and `2026-05-29T10:00:05Z`

### Requirement: Degraded runtime monitoring banner
The runtime dashboard SHALL clearly indicate when live runtime monitoring is degraded while preserving the last successful runtime data on screen.

#### Scenario: Live stream disconnect preserves last data
- **WHEN** the live stream disconnects after the page has already rendered runtime summary and point values
- **THEN** the page shows a degraded monitoring banner
- **AND** the last successful runtime summary and point values remain visible until fresher data arrives

##### Example: Disconnect after live data keeps last snapshot on screen
- **GIVEN** the page is already showing runtime summary and the latest values for `device-A`
- **WHEN** the live stream disconnects
- **THEN** the page switches to degraded monitoring while keeping the last successful summary and point values visible
