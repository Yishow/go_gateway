## ADDED Requirements

### Requirement: Studio v2 workspace rule autosave APIs

The system SHALL expose workspace-scoped source-rule APIs for Step 2 autosave.

#### Scenario: List workspace rules

- **WHEN** a client requests `GET /api/v1/datalink/studio-v2/workspace/source-rules`
- **THEN** the API returns the persisted rules currently attached to the singleton v2 workspace

#### Scenario: Save one valid workspace rule

- **WHEN** a client requests `POST /api/v1/datalink/studio-v2/workspace/source-rules` or `PUT /api/v1/datalink/studio-v2/workspace/source-rules/:id` with a valid rule payload
- **THEN** the API persists that rule and returns the saved record including `revision_id`
- **AND** the response is scoped to that one rule rather than a bulk transaction

#### Scenario: Reject one invalid workspace rule without touching others

- **WHEN** a client sends one invalid Step 2 rule payload
- **THEN** the API returns a validation error for that rule
- **AND** previously saved rules in the same workspace remain unchanged
