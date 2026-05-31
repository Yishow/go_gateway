## ADDED Requirements

### Requirement: Studio v2 workspace bootstrap API

The system SHALL expose `GET /api/v1/datalink/studio-v2/workspace` as the bootstrap endpoint for the singleton v2 workspace.

#### Scenario: Bootstrap returns existing workspace

- **WHEN** a client requests `GET /api/v1/datalink/studio-v2/workspace` and the singleton workspace already exists
- **THEN** the API returns HTTP 200 with the persisted workspace metadata
- **AND** the payload includes at least `id`, `kind`, `status`, `ordered_device_ids`, `created_at`, and `updated_at`

#### Scenario: Bootstrap auto-creates missing workspace

- **WHEN** a client requests `GET /api/v1/datalink/studio-v2/workspace` and the singleton workspace does not exist yet
- **THEN** the API creates the singleton workspace before responding
- **AND** the response still returns HTTP 200 with the created workspace metadata
