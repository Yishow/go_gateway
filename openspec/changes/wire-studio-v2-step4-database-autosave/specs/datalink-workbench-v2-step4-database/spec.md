## ADDED Requirements

### Requirement: Valid-only database autosave

The Step 4 database workspace SHALL autosave only valid connector and database target content to the backend workspace.

#### Scenario: Valid connector edit saves immediately

- **GIVEN** the workspace already has persisted Step 4 database metadata
- **WHEN** the operator edits the connector into a valid state
- **THEN** the system saves that connector metadata immediately
- **AND** the persisted backend version reflects the latest valid values

#### Scenario: Invalid connector edit stays local

- **GIVEN** the connector already has a last successful persisted version
- **WHEN** the operator edits the connector into an invalid state
- **THEN** the system does not overwrite the persisted connector version
- **AND** the UI keeps the invalid local values visible with an unsaved marker

### Requirement: Per-target save isolation

The Step 4 database workspace SHALL isolate save success and failure per target row.

#### Scenario: One invalid target row does not block another valid target row

- **GIVEN** the workspace contains multiple database target rows
- **WHEN** one target row remains invalid while another target row is edited into a valid state
- **THEN** the valid target row still saves successfully
- **AND** the invalid target row remains local with its own error state

### Requirement: Database autosave remains pre-activation

Saving Step 4 database content SHALL NOT by itself start runtime collection.

#### Scenario: Autosave does not start runtime

- **WHEN** the operator successfully autosaves connector or target edits in Step 4
- **THEN** the persisted database content is updated
- **AND** runtime collection is still not started solely by that autosave
