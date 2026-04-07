# datalink-api Specification

## Purpose
TBD - created by archiving change add-device-data-pipeline. Update Purpose after archive.
## Requirements
### Requirement: Management APIs

The system SHALL expose APIs for devices, points, tags, and mappings with persistent storage.

All data SHALL be persisted to SQLite (for development/testing) or PostgreSQL (for production) using a Repository pattern with factory-based selection.

#### Scenario: Create mapping via API

- WHEN a client posts a mapping
- THEN the mapping is saved to database and returned with an id

#### Scenario: Data persistence across restarts

- WHEN the server restarts
- THEN all previously saved data is available

---

### Requirement: Connection and preview APIs
The system SHALL expose APIs to test device connections and preview mapping results.

#### Scenario: Mapping preview API
- WHEN a client requests a preview for a mapping
- THEN the API returns raw and transformed values

### Requirement: Write settings API

The system SHALL expose APIs to read and update write settings, including timestamp precision and partition interval.

Settings SHALL be persisted to database (not hardcoded).

#### Scenario: Update write precision

- WHEN a client updates the write precision to milliseconds
- THEN subsequent writes use millisecond precision

#### Scenario: Get settings

- WHEN a client requests settings
- THEN the API returns current settings from database

### Requirement: Polling Groups API

The system SHALL expose APIs for managing polling groups.

Endpoints:

- `GET /datalink/polling-groups` - List all polling groups
- `POST /datalink/polling-groups` - Create polling group
- `GET /datalink/polling-groups/:id` - Get polling group by ID
- `PUT /datalink/polling-groups/:id` - Update polling group
- `DELETE /datalink/polling-groups/:id` - Delete polling group

#### Scenario: Create polling group

- WHEN a client posts a polling group with name and interval
- THEN the polling group is saved and returned with an id

#### Scenario: List polling groups

- WHEN a client requests polling groups
- THEN all polling groups are returned

---

### Requirement: Protocols API

The system SHALL expose an API to retrieve supported protocols and their configuration schemas.

Endpoint:

- `GET /datalink/protocols` - List supported protocols

#### Scenario: List protocols

- WHEN a client requests protocols
- THEN the API returns protocol list with type, name, description, and config_schema

---

### Requirement: Health API

The system SHALL expose a health check API for the Datalink service.

Endpoint:

- `GET /datalink/health` - Check service health

#### Scenario: Health check

- WHEN a client requests health status
- THEN the API returns status and service name

---

### Requirement: Tag lifecycle APIs

The system SHALL expose APIs for tag lifecycle management.

Endpoints:

- `POST /datalink/tags/:id/activate` - Activate a draft tag
- `POST /datalink/tags/:id/retire` - Retire an active tag
- `POST /datalink/tags/batch` - Batch create tags
- `POST /datalink/tags/validate-key` - Validate tag key format

#### Scenario: Activate tag

- WHEN a client activates a draft tag
- THEN the tag status changes to active

#### Scenario: Retire tag

- WHEN a client retires an active tag
- THEN the tag status changes to retired

#### Scenario: Batch create tags

- WHEN a client posts multiple tags
- THEN all valid tags are created and errors are reported

#### Scenario: Validate tag key

- WHEN a client validates a tag key
- THEN the API returns validity, normalized key, and existence status

---

### Requirement: Device batch test API

The system SHALL expose an API to batch test device connections.

Endpoint:

- `POST /datalink/devices/test-batch` - Test multiple device connections

#### Scenario: Batch test connections

- WHEN a client posts device IDs for testing
- THEN the API tests all devices and returns individual results

---

### Requirement: Point polling APIs

The system SHALL expose APIs for on-demand point polling.

Endpoints:

- `POST /datalink/points/:id/poll` - Poll single point
- `POST /datalink/points/poll` - Batch poll multiple points

#### Scenario: Poll single point

- WHEN a client requests a point poll
- THEN the API returns the current value and timestamp

#### Scenario: Batch poll points

- WHEN a client posts point IDs for polling
- THEN the API polls all points and returns results

---

### Requirement: Pipeline validation API

The system SHALL expose an API to validate transform pipelines.

Endpoint:

- `POST /datalink/mappings/validate-pipeline` - Validate pipeline

#### Scenario: Validate valid pipeline

- WHEN a client submits a valid pipeline
- THEN the API returns valid=true

#### Scenario: Validate invalid pipeline

- WHEN a client submits an invalid pipeline
- THEN the API returns valid=false with error details

---

### Requirement: Preview stream API

The system SHALL expose an SSE endpoint for live mapping preview.

Endpoint:

- `GET /datalink/preview/stream` - Server-Sent Events stream

#### Scenario: Subscribe to preview stream

- WHEN a client connects to the preview stream with a mapping ID
- THEN the server pushes raw and transformed values in real-time

---

### Requirement: Device readiness API

The system SHALL expose an API endpoint to check device readiness status using the explicit planning, activation, and apply eligibility contract.

Endpoint:

- `GET /api/v1/datalink/devices/:id/readiness` - Check device readiness

#### Scenario: Get readiness status

- **WHEN** a client requests readiness status for a device
- **THEN** the API returns `connect_status`, `probe_status`, `planning_allowed`, `activation_allowed`, `apply_allowed`, and blocking metadata

#### Scenario: Readiness response format

- **WHEN** a readiness check is performed
- **THEN** the API returns a response with explicit planning/activation/apply eligibility fields
- **AND** does not require clients to infer workflow state from one final `ready` boolean alone

### Requirement: Rule-scoped candidate query and recompute APIs
The API SHALL expose rule-scoped endpoints to query and recompute candidate snapshots for a source-rule revision.

Endpoints:

- `GET /api/v1/datalink/source-rules/:id/candidates` - Get the current candidate snapshot for a source rule
- `POST /api/v1/datalink/source-rules/:id/candidates/recompute` - Recompute the candidate snapshot for the latest source-rule revision

The candidate query response MUST include at least:

- `source_rule_id`
- `revision_id`
- `tags`
- `database_outputs`
- `local_modbus_outputs`

Each returned candidate object MUST include at least:

- `id`
- `status`
- `proposed_signature`
- `last_applied_signature`
- `blocking_reason`

#### Scenario: Query rule candidates
- **WHEN** a client requests candidates for a source rule
- **THEN** the API returns the current tag, database, and Local Modbus candidate sets with revision and status metadata

#### Scenario: Recompute candidate snapshot
- **WHEN** a client requests candidate recompute for a source rule
- **THEN** the API regenerates the current revision's candidate snapshot using the shared identity/signature contract
- **AND** returns the recomputed snapshot without introducing target-specific apply side effects

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
