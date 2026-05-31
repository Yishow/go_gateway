## ADDED Requirements

### Requirement: Studio v2 workspace device autosave APIs

The system SHALL expose workspace-scoped device APIs for Step 1 autosave.

#### Scenario: List workspace devices

- **WHEN** a client requests `GET /api/v1/datalink/studio-v2/workspace/devices`
- **THEN** the API returns the devices currently attached to the singleton v2 workspace in persisted order

#### Scenario: Save one valid workspace device

- **WHEN** a client requests `POST /api/v1/datalink/studio-v2/workspace/devices` or `PUT /api/v1/datalink/studio-v2/workspace/devices/:id` with a valid device payload
- **THEN** the API persists that device and returns the saved record
- **AND** the response is scoped to that one device rather than a bulk transaction

#### Scenario: Reject one invalid workspace device without touching others

- **WHEN** a client sends one invalid Step 1 device payload
- **THEN** the API returns a validation error for that device
- **AND** previously saved devices in the same workspace remain unchanged

### Requirement: Studio v2 workspace device order API

The system SHALL expose `PUT /api/v1/datalink/studio-v2/workspace/device-order` to persist the device tab order used by V2.

#### Scenario: Persist device order

- **WHEN** a client submits an ordered list of workspace device ids
- **THEN** the API stores that order for the singleton workspace
- **AND** a later workspace bootstrap returns the same order

##### Example: reordered tabs survive reload

- **GIVEN** the current workspace device order is `[dev-A, dev-B, dev-C]`
- **WHEN** a client submits `[dev-C, dev-A, dev-B]` to `PUT /api/v1/datalink/studio-v2/workspace/device-order`
- **THEN** the next workspace bootstrap returns `[dev-C, dev-A, dev-B]`
