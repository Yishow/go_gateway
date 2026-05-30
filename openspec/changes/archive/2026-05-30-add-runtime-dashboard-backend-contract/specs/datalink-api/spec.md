## ADDED Requirements

### Requirement: Runtime status snapshot API

The system SHALL expose `GET /api/v1/datalink/runtime/status` as the polling-compatible snapshot endpoint for runtime monitoring clients.

The endpoint SHALL accept an optional `device_id` query parameter.

Each successful response SHALL include top-level fields `running`, `uptime_seconds`, `metrics`, and `collectors`.

The `metrics` object SHALL include at least:

- `collected_total`
- `write_success_total`
- `write_error_total`
- `mapping_error_total`
- `point_state_error_total`

Each collector object SHALL include at least:

- `device_id`
- `device_name`
- `protocol`
- `status`
- `points_total`
- `points_healthy`
- `points_stale`
- `points_error`
- `last_read_at`
- `last_error`
- `breaker_state`

#### Scenario: Request one device snapshot

- **WHEN** a client requests `GET /api/v1/datalink/runtime/status?device_id=device-A`
- **THEN** the response includes runtime top-level metrics
- **AND** the `collectors` array contains exactly the collector summary for `device-A`

#### Scenario: Request fleet snapshot without device filter

- **WHEN** a client requests `GET /api/v1/datalink/runtime/status` without `device_id`
- **THEN** the response includes runtime top-level metrics
- **AND** the `collectors` array SHALL contain all device collector summaries visible to the runtime service

### Requirement: Runtime value stream API

The system SHALL expose `GET /api/v1/datalink/runtime/stream` as the live runtime monitoring stream for one device.

The endpoint SHALL require `device_id` and SHALL accept an optional comma-separated `point_ids` query parameter.

The stream SHALL emit the following SSE event types:

- `value`
- `status`
- `heartbeat`

A `value` event payload SHALL include at least `device_id`, `point_id`, `address`, `raw_value`, `transformed_value`, `quality`, `stale`, and `timestamp`.

A `status` event payload SHALL include at least `device_id`, `status`, `points_total`, `points_healthy`, `points_stale`, `points_error`, `last_read_at`, `last_error`, and `breaker_state`.

A `heartbeat` event payload SHALL include `ts`.

#### Scenario: Reject stream request without device id

- **WHEN** a client requests `GET /api/v1/datalink/runtime/stream` without `device_id`
- **THEN** the server returns HTTP 400
- **AND** the request does not create a runtime subscription

#### Scenario: Stream only selected points for one device

- **WHEN** a client requests `GET /api/v1/datalink/runtime/stream?device_id=device-A&point_ids=pt-1,pt-2`
- **THEN** the stream emits `value` events only for `pt-1` and `pt-2`
- **AND** the stream continues to emit `status` and `heartbeat` events for `device-A`

#### Scenario: Emit device status change on runtime degradation

- **WHEN** runtime monitoring detects that `device-A` transitions from healthy collection to stale, point-error, or breaker-open state
- **THEN** the stream emits a `status` event for `device-A`
- **AND** the payload reflects the same derived status model used by the runtime snapshot endpoint
