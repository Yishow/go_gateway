## ADDED Requirements

### Requirement: Valid-only mapping autosave

The Step 3 mapping workspace SHALL autosave only valid mapping rows to the backend workspace.

#### Scenario: Valid mapping row saves immediately

- **GIVEN** a mapping row already belongs to the singleton v2 workspace
- **WHEN** the operator edits that row into a valid state
- **THEN** the system saves that row immediately
- **AND** the persisted mapping reflects the latest valid values

#### Scenario: Invalid mapping row stays local

- **GIVEN** a mapping row already has a last successful persisted version
- **WHEN** the operator edits that row into an invalid state
- **THEN** the system does not overwrite the persisted mapping
- **AND** the UI keeps the invalid local values visible with an unsaved marker

### Requirement: Per-point mapping save isolation

The Step 3 mapping workspace SHALL isolate save success and failure per point row.

#### Scenario: One invalid row does not block another valid row

- **GIVEN** the workspace contains multiple mapping rows
- **WHEN** one row remains invalid while another row is edited into a valid state
- **THEN** the valid row still saves successfully
- **AND** the invalid row remains local with its own error state
