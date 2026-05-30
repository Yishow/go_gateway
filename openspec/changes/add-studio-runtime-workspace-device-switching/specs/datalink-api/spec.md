## ADDED Requirements

### Requirement: Studio v2 workspace runtime context API

The system SHALL expose `GET /api/v1/datalink/studio-v2/workspace/runtime-context` as the workspace-scoped runtime context endpoint.

#### Scenario: Ordered runtime context payload

- **WHEN** a client requests `GET /api/v1/datalink/studio-v2/workspace/runtime-context`
- **THEN** the API returns the workspace id, ordered device list, per-device availability fields, and `default_device_id`

#### Scenario: Empty runtime context payload

- **WHEN** the singleton workspace currently has no available devices
- **THEN** the API still returns HTTP 200
- **AND** the device list is empty or contains only unavailable devices
- **AND** the payload still includes enough metadata for the runtime page to render an empty state
