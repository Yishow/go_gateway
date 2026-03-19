## ADDED Requirements

### Requirement: Database output workflow is layered by connector, schema context, and mapping
The workbench SHALL present database output as three linked layers: connector selection, schema/table context, and tag-to-column mapping.

#### Scenario: Connector selection scopes schema context
- **WHEN** an operator selects a database connector in the Output step
- **THEN** the available schema and table context SHALL scope to that connector only
- **AND** stale schema or mapping state from previously selected connectors SHALL NOT remain active

#### Scenario: Mapping action uses current layered context
- **WHEN** an operator binds a tag to a database column
- **THEN** the mapping SHALL be created against the currently selected connector, schema, and table context
- **AND** the UI SHALL NOT apply the action to stale prior selections

### Requirement: Database connector scope for this workflow is SQLite and PostgreSQL
The system SHALL support SQLite and PostgreSQL as the database connector kinds for this workbench flow in this change.

#### Scenario: Supported connectors are first-class options
- **WHEN** an operator creates or edits a database target connector in the workbench
- **THEN** the workbench SHALL present SQLite and PostgreSQL as supported connector kinds
- **AND** SHALL NOT advertise unsupported connector kinds as first-class options in this workflow
