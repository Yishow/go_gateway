## ADDED Requirements

### Requirement: Studio v2 workspace mapping autosave APIs

The system SHALL expose workspace-scoped mapping APIs for Step 3 autosave.

#### Scenario: List workspace mappings

- **WHEN** a client requests `GET /api/v1/datalink/studio-v2/workspace/mappings`
- **THEN** the API returns the persisted mappings currently attached to the singleton v2 workspace

#### Scenario: Save one valid workspace mapping

- **WHEN** a client requests `POST /api/v1/datalink/studio-v2/workspace/mappings` or `PUT /api/v1/datalink/studio-v2/workspace/mappings/:id` with a valid mapping payload
- **THEN** the API persists that mapping and returns the saved record
- **AND** the response is scoped to that one row rather than a bulk transaction

#### Scenario: Reject one invalid workspace mapping without touching others

- **WHEN** a client sends one invalid Step 3 mapping payload
- **THEN** the API returns a validation error for that mapping row
- **AND** previously saved mappings in the same workspace remain unchanged
