## ADDED Requirements

### Requirement: Device readiness API

The system SHALL expose an API endpoint to check device readiness status.

Endpoint:

- `GET /api/v1/datalink/devices/:id/readiness` - Check device readiness

#### Scenario: Get readiness status

- **WHEN** a client requests readiness status for a device
- **THEN** the API returns readiness status with checks and suggestions

#### Scenario: Readiness response format

- **WHEN** a readiness check is performed
- **THEN** the API returns a response with `ready` boolean, `checks` object, `missing` array, and `suggestions` array

### Requirement: Device collection statistics API

The system SHALL expose an API endpoint to retrieve device collection statistics.

Endpoint:

- `GET /api/v1/datalink/devices/:id/stats` - Get device collection statistics

#### Scenario: Get device statistics

- **WHEN** a client requests device statistics
- **THEN** the API returns collection statistics including last collection time, collection count, and error count

### Requirement: Device status stream API

The system SHALL expose an SSE endpoint for real-time device status updates.

Endpoint:

- `GET /api/v1/datalink/devices/stream` - Server-Sent Events stream for device status

#### Scenario: Subscribe to device status stream

- **WHEN** a client connects to the device status stream
- **THEN** the server pushes device status updates in real-time
