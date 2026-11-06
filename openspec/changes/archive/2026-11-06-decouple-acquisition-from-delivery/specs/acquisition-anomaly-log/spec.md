## Purpose

Detect, persist and present acquisition anomalies — device disconnection, points that stop updating, and points whose value stays frozen — so operators can see when and for how long collected data was unreliable.

## ADDED Requirements

### Requirement: Acquisition anomalies are detected and recorded

The runtime SHALL record three anomaly kinds as persisted events with a start and an end: `disconnected` for a device whose connection fails or whose breaker opens, `stale` for a point with no successful read for longer than its device's stale threshold, and `frozen` for a point whose value has not changed for longer than its device's frozen threshold. At most one active event SHALL exist per device, point and kind. An event SHALL end when its condition clears, recording the end time. The stale threshold SHALL default to five times the device acquisition interval and at least 10 seconds; the frozen threshold SHALL default to disabled. Detection SHALL read the runtime's latest point state and SHALL NOT issue device reads.

#### Scenario: Stale point opens and closes an event

- **GIVEN** a device with a 1000 ms interval and default stale threshold of 10 seconds
- **WHEN** point D5 has no successful read for 12 seconds and then reads successfully
- **THEN** one `stale` event for D5 starts when the threshold is exceeded and ends at the successful read

#### Scenario: Frozen detection is opt-in

- **WHEN** a device has no frozen threshold configured and a point keeps value 0 for an hour
- **THEN** no `frozen` event is recorded
- **AND** after the device sets a 600 second frozen threshold, a point unchanged for longer than 600 seconds opens one `frozen` event

#### Scenario: Disconnect is recorded once

- **WHEN** a device's connection fails repeatedly for 30 seconds and then recovers
- **THEN** exactly one `disconnected` event covers that period

### Requirement: Anomaly history is bounded and queryable

Anomaly events SHALL be stored with device, optional point, kind, start time, optional end time, last observed value and detail. Each device SHALL use 5000 events as a retention target. When the total exceeds that target, the oldest ended events SHALL be deleted first. Active events SHALL never be deleted, even when active events alone exceed 5000; ending an event SHALL trigger pruning again. The API SHALL list a device's events filtered by state (`active`, `ended`) and kind with pagination and total count, and SHALL return a per-device count of active events.

#### Scenario: Filtered listing

- **WHEN** a client requests a device's events with state `active` and kind `stale`
- **THEN** only active stale events of that device are returned with the total count

#### Scenario: Retention keeps active events

- **GIVEN** a device with 5000 stored events of which 3 are active
- **WHEN** a new event is recorded
- **THEN** the oldest ended event is removed and all 3 active events remain

#### Scenario: Active events exceed the retention target

- **GIVEN** a device with 5001 active events and no ended events
- **WHEN** retention runs
- **THEN** all 5001 active events remain
- **AND** when one of those events ends, that ended event is removed and the 5000 remaining active events stay

### Requirement: Home exposes anomalies beside the device points action

The recording home SHALL show an "acquisition anomalies" button next to each device's read-points action in both the device list and the device summary. When the device has active anomalies, the button SHALL use warning styling and show the active count as text; when it has none, it SHALL be a normal secondary button; when the summary cannot be read, it SHALL state that the status is unknown. Activating it SHALL open a classified form modal listing that device's events with kind and state filters, start, end, duration, point address and last value. A failed read SHALL show an actionable error with retry and SHALL NOT be shown as an empty list.

#### Scenario: Warning state with count

- **WHEN** device MC has 2 active anomalies
- **THEN** its anomalies button shows warning styling with the text count 2 next to the read-points action

#### Scenario: Read failure is not an empty state

- **WHEN** the anomaly list request fails
- **THEN** the modal shows an error with a retry action instead of a "no anomalies" message
