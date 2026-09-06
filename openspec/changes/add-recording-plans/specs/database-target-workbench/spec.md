## MODIFIED Requirements

### Requirement: Database output workflow is layered by connector, schema context, and mapping
The workbench SHALL support a managed recording-plan workflow and an explicit legacy/custom-table workflow. Both SHALL preserve source revision, connector identity, schema context and apply review; managed mode SHALL generate storage layout from confirmed recording intent rather than require manual columns.

#### Scenario: Connector selection scopes schema context and candidate board
- **WHEN** an operator selects a database connector for a rule or recording plan
- **THEN** schema/table candidates are scoped only to that connector and current source/plan revisions.

#### Scenario: Apply uses current layered context
- **WHEN** an operator applies a database output plan
- **THEN** connector, source/plan revision, table, key, interval and mapping signatures are revalidated
- **AND** stale selections cannot be applied.

### Requirement: Database connector scope for this workflow is SQLite and PostgreSQL
The system SHALL present connector kinds by verified capabilities rather than a hard-coded two-kind list. SQLite and PostgreSQL SHALL remain supported; MySQL SHALL retain its existing inspection/schema behavior and pass the managed recording capability suite before that mode is enabled. Other kinds MUST NOT be advertised as ready without equivalent evidence.

#### Scenario: Supported connectors are first-class options
- **WHEN** an operator creates or edits a target connector
- **THEN** only currently verified modes for that kind can be selected and unsupported modes explain their limitation.

### Requirement: Database planner organizes compatible tags into grouped rows
The database workbench SHALL automatically organize managed records by explicit equipment and measurement identity. Legacy/custom-table grouped mode SHALL retain compatible row-plan inference as a reviewable suggestion, never as an unconfirmed semantic binding.

#### Scenario: Prefix-inferred tags appear as one row plan
- **WHEN** legacy grouped candidates contain meter/A1, meter/A2, meter/A3 and meter/kw
- **THEN** one suggested meter row with a1, a2, a3 and kw is shown for confirmation
- **AND** the user is not forced to bind each field as the first-pass workflow.

#### Scenario: Managed plan contains mixed sensor meanings
- **WHEN** temperature and cumulative volume belong to the same equipment
- **THEN** the managed plan groups them for operator review but preserves different recording and calculation policies.
