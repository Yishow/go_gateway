## ADDED Requirements

### Requirement: Dedicated post-setup runtime dashboard route

The system SHALL expose a dedicated runtime dashboard route for post-setup monitoring without replacing the `/studio/v2` setup entry flow.

#### Scenario: Open runtime dashboard by URL

- **WHEN** an operator navigates to `/studio/runtime?device_id=device-A`
- **THEN** the system renders the runtime dashboard page
- **AND** the system does not redirect to `/studio` or `/studio/v2`

#### Scenario: Setup flow remains separate

- **WHEN** an operator navigates to `/studio/v2`
- **THEN** the system continues to render the setup workbench flow
- **AND** the runtime dashboard does not replace the v2 landing page

### Requirement: Runtime dashboard uses device-focused context

The runtime dashboard SHALL treat `device_id` as the canonical selected-device context.

#### Scenario: Missing device context shows focused empty state

- **WHEN** an operator opens `/studio/runtime` without `device_id`
- **THEN** the system shows a missing-device-context state
- **AND** the page does not silently fall back to a fleet-wide dashboard

#### Scenario: Switching device updates route context

- **WHEN** an operator switches the selected device from `device-A` to `device-B` inside the runtime dashboard
- **THEN** the system updates the route query to `device_id=device-B`
- **AND** the dashboard reloads its runtime data using `device-B` as the only selected device context

### Requirement: Runtime dashboard layers snapshot, live stream, and degraded fallback

The runtime dashboard SHALL load runtime data by fetching a snapshot first, layering live stream updates on top, and degrading to snapshot polling when the live stream is unavailable.

#### Scenario: Snapshot renders before live stream attaches

- **WHEN** an operator opens `/studio/runtime?device_id=device-A`
- **THEN** the system first fetches the runtime snapshot for `device-A`
- **AND** the page renders the latest available runtime summary before live stream events arrive

#### Scenario: Stream disconnect degrades to snapshot polling

- **WHEN** the runtime live stream for `device-A` disconnects after the initial snapshot succeeds
- **THEN** the system preserves the last successful runtime summary on screen
- **AND** the page enters a degraded state that continues refreshing the snapshot until live streaming recovers
