## ADDED Requirements

### Requirement: Studio v2 workspace database autosave APIs

The system SHALL expose workspace-scoped Step 4 database autosave APIs.

#### Scenario: Save workspace database config

- **WHEN** a client requests `PUT /api/v1/datalink/studio-v2/workspace/database-config` with a valid connector payload
- **THEN** the API persists the workspace database config and returns the saved metadata

#### Scenario: Save one valid target row

- **WHEN** a client requests `PUT /api/v1/datalink/studio-v2/workspace/database-targets/:point_id` with a valid target payload
- **THEN** the API persists that target row and returns the saved row
- **AND** the response is scoped to that one target row rather than a bulk transaction

#### Scenario: Reject one invalid target row without touching others

- **WHEN** a client sends one invalid Step 4 target payload
- **THEN** the API returns a validation error for that target row
- **AND** previously saved workspace database rows remain unchanged
