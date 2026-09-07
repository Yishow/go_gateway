## MODIFIED Requirements

### Requirement: Connector configuration form
The Step 4 connector section SHALL present saved connections or a new connection with only relevant fields, verified kind capabilities and explicit connection-test state. Managed recording SHALL generate storage context from confirmed plan intent; custom-table mode SHALL retain table and write-policy controls in an advanced section. The form MUST NOT use example endpoints, tables or ready status as real configured state.

#### Scenario: Default connector on first render
- **WHEN** Step 4 mounts without a persisted connector
- **THEN** it shows an unconfigured draft with clearly labelled placeholders and supported-kind choices
- **AND** it does not claim connection or write readiness.

#### Scenario: Kind switch triggers auto-assign
- **WHEN** the operator changes connector kind
- **THEN** stale tests and schema previews are invalidated and actual metadata is reloaded
- **AND** bindings remain only when valid for the new identity; missing or ambiguous matches require confirmation instead of sample-column reassignment.

#### Scenario: Write strategy persistence
- **WHEN** an operator selects a compatible recording strategy
- **THEN** the plan persists and explains whether it appends history or updates latest state
- **AND** an update strategy without a valid unique identity is blocked.

### Requirement: Tag-to-column auto-assignment
The system SHALL provide confirmed assignments, reviewable suggestions and unmatched results using the actual selected schema and compatible measurement semantics. Existing confirmed targets MUST be revalidated. Name similarity alone SHALL NOT confirm uncertain semantics; index fallback, wraparound reuse and sample columns MUST NOT produce production assignments.

#### Scenario: Eight points to nine columns
- **GIVEN** eight confirmed measurements and sufficient real compatible columns
- **WHEN** matching runs
- **THEN** valid unambiguous assignments and reviewable suggestions are shown without assigning any column twice within a row.

#### Scenario: Existing target preserved
- **WHEN** an existing confirmed target remains compatible with the current connector, table, column and measurement definition
- **THEN** it is preserved; otherwise it is marked for repair rather than silently rebound.

#### Scenario: Fewer columns than points causes wrap and conflict
- **GIVEN** eight enabled measurements and only four compatible columns
- **WHEN** matching runs
- **THEN** extra measurements remain unmatched and the UI offers create-column-plan, different-table or explicit exclusion options
- **AND** no wraparound or duplicate assignment is generated.
