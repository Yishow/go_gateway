## ADDED Requirements

### Requirement: Runtime snapshot and stream expose projection reconciliation state

The runtime snapshot and stream contract SHALL expose whether each selected device is running the latest persisted workspace projection.

#### Scenario: Snapshot and stream carry projection status

- **WHEN** a client requests runtime snapshot or subscribes to runtime stream for a selected device
- **THEN** the response includes whether the device projection is aligned, deferred, or stale against the latest persisted workspace state
- **AND** the client SHALL NOT need to invent a projection-truth guess from unrelated counters

##### Example: stale snapshot reports its own projection status

- **GIVEN** runtime snapshot for dev-A still runs projection v12 while persisted workspace is at projection v13
- **WHEN** the client fetches snapshot or stream status for dev-A
- **THEN** the response marks dev-A as stale or equivalent instead of implying aligned truth
