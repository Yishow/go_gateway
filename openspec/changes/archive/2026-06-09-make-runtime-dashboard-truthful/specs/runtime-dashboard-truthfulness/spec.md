## ADDED Requirements

### Requirement: Runtime dashboard renders only backend-backed data

The runtime dashboard SHALL render only values, counters, and summaries backed by the runtime backend contract.

#### Scenario: Missing backend data produces explicit empty state

- **WHEN** the runtime backend does not have real snapshot or stream data for the selected device
- **THEN** the dashboard renders an explicit empty, degraded, or error state for that missing data
- **AND** the dashboard SHALL NOT synthesize a healthy runtime summary from defaults or unrelated state

##### Example: selected device has no snapshot yet

- **GIVEN** dev-A is selected and runtime backend has not produced any snapshot for dev-A
- **WHEN** the dashboard loads
- **THEN** the dashboard shows a waiting or empty state for dev-A instead of a healthy default summary

### Requirement: Runtime dashboard keeps selected device context truthful

The runtime dashboard SHALL keep the selected device context aligned with the workspace-scoped runtime context.

#### Scenario: Selected device does not silently fall back to another device

- **WHEN** the selected device context is missing, stale, or unavailable
- **THEN** the dashboard surfaces that specific selected-device problem
- **AND** the dashboard SHALL NOT silently switch to another device or fleet-wide synthetic summary

##### Example: dev-A missing does not fall back to dev-B

- **GIVEN** dev-A is selected, dev-A has no runtime data, and dev-B has healthy runtime data
- **WHEN** the dashboard loads for dev-A
- **THEN** the dashboard shows dev-A as missing or degraded and does not display dev-B as a substitute

### Requirement: Snapshot, stream, and degraded fallback stay semantically distinct

The runtime dashboard SHALL keep snapshot data, live stream data, and degraded fallback state semantically distinct.

#### Scenario: Stream failure does not fake fresh live data

- **WHEN** snapshot loads successfully but live stream is unavailable or broken
- **THEN** the dashboard keeps the last truthful snapshot visible with a degraded live-state indicator
- **AND** the dashboard SHALL NOT present synthetic live freshness it did not receive from the backend

##### Example: broken stream keeps the last truthful snapshot

- **GIVEN** snapshot for dev-A loaded at 10:00:00Z and the live stream disconnects at 10:00:05Z
- **WHEN** the dashboard continues rendering dev-A
- **THEN** the dashboard keeps the 10:00:00Z snapshot with a degraded live indicator instead of pretending newer live values arrived
