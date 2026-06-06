## ADDED Requirements

### Requirement: Runtime backend expresses degraded and empty states explicitly

The runtime backend SHALL express empty, degraded, unavailable, and stale states explicitly in snapshot and stream responses.

#### Scenario: Backend marks stream unavailable explicitly

- **WHEN** runtime stream is unavailable for the selected device
- **THEN** the backend response or event model marks that unavailability explicitly
- **AND** the client SHALL NOT need to infer it from missing events alone

##### Example: selected stream outage is explicit

- **GIVEN** dev-A snapshot is available but the live stream channel for dev-A cannot attach
- **WHEN** the backend responds to the client
- **THEN** the contract marks stream unavailable for dev-A instead of silently returning no events
