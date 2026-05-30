## ADDED Requirements

### Requirement: Device-scoped runtime monitoring contract

The system SHALL provide a device-scoped runtime monitoring contract so a post-setup client can monitor one committed device without subscribing to fleet-wide point traffic.

#### Scenario: Request snapshot for one device

- **WHEN** a client requests the runtime snapshot with `device_id=device-A`
- **THEN** the system returns runtime metrics plus exactly one collector summary for `device-A`
- **AND** the response does not require the client to filter unrelated devices locally

#### Scenario: Subscribe with point filter inside one device

- **WHEN** a client opens the runtime stream with `device_id=device-A&point_ids=pt-1,pt-2`
- **THEN** the system emits live `value` events only for `pt-1` and `pt-2`
- **AND** the stream remains scoped to `device-A`

### Requirement: Snapshot and stream use one derived device status model

The system SHALL derive runtime device status consistently across the snapshot response and stream `status` events using the same freshness, point error, and breaker-state rules.

#### Scenario: Stale points produce warning status consistently

- **WHEN** the latest read time for one or more enabled points exceeds the effective freshness window for a device
- **THEN** the runtime snapshot reports the device status as `warning`
- **AND** the next emitted `status` event for that device reports the same `warning` status with matching stale-point counts

##### Example: stale counts stay aligned

| Source | points_total | points_healthy | points_stale | points_error | breaker_state | derived status |
| ----- | ------------ | -------------- | ------------ | ------------ | ------------- | -------------- |
| snapshot | 4 | 3 | 1 | 0 | closed | warning |
| status event | 4 | 3 | 1 | 0 | closed | warning |

#### Scenario: Open breaker produces error status consistently

- **WHEN** a device breaker state becomes `open`
- **THEN** the runtime snapshot reports the device status as `error`
- **AND** the stream emits a `status` event with `breaker_state=open` and `status=error`

### Requirement: Live monitoring stream remains loss-tolerant

The system SHALL protect runtime collection progress when a monitoring subscriber is slow, even if some intermediate live events cannot be delivered to that subscriber.

#### Scenario: Slow subscriber does not block collection

- **WHEN** a monitoring subscriber stops consuming events and its stream buffer is saturated
- **THEN** runtime collection continues without waiting for that subscriber
- **AND** the system preserves connection availability for other subscribers and future snapshot requests

##### Example: saturated subscriber buffer

| Subscriber state | Incoming value events | Expected collector behavior | Expected slow-subscriber behavior |
| ----- | --------------------- | --------------------------- | --------------------------------- |
| buffer full | 10 new values for `device-A` | collection loop continues processing all 10 values | some intermediate SSE deliveries are dropped |
