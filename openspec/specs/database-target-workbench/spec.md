# database-target-workbench Specification

## Purpose
Define the canonical Studio-side database output workflow so connector selection, schema context, and tag-to-column mapping remain scoped, coherent, and limited to the supported SQLite and PostgreSQL connectors.
## Requirements
### Requirement: Database output workflow is layered by connector, schema context, and mapping
The workbench SHALL present database output as four linked layers: source-rule revision context, connector selection, schema/table context, and grouped row planner review/apply.

#### Scenario: Connector selection scopes schema context and candidate board
- **WHEN** an operator selects a database connector in the Output step for a source rule
- **THEN** the available schema and table context SHALL scope to that connector only
- **AND** the rule-scoped database candidate board SHALL update using only that connector context

#### Scenario: Apply uses current layered context
- **WHEN** an operator applies a database output candidate
- **THEN** the mapping SHALL be applied against the current source-rule revision and the selected connector, schema, table, column, group key, and effective interval context
- **AND** the UI SHALL NOT apply the action to stale prior selections

### Requirement: Database connector scope for this workflow is SQLite and PostgreSQL
The system SHALL support SQLite and PostgreSQL as the database connector kinds for this workbench flow in this change.

#### Scenario: Supported connectors are first-class options
- **WHEN** an operator creates or edits a database target connector in the workbench
- **THEN** the workbench SHALL present SQLite and PostgreSQL as supported connector kinds
- **AND** SHALL NOT advertise unsupported connector kinds as first-class options in this workflow

### Requirement: Database output is rule-scoped and connector-aware
The workbench SHALL generate database output candidates from rule-owned tag state using persisted connector and table context.

The system SHALL revalidate blocked or `out_of_sync` database candidates whenever the referenced connector context becomes valid again.

#### Scenario: Connector invalidation blocks database candidate apply
- **WHEN** the selected connector, schema, table, or column becomes invalid for an existing database output candidate
- **THEN** the candidate or applied mapping is marked as `blocked` or `out_of_sync`
- **AND** database apply SHALL NOT silently rebind to a different connector or table

#### Scenario: Restored connector context re-enables database apply
- **WHEN** a previously invalid connector, schema, table, and column context becomes valid again for an existing database output candidate
- **THEN** the system revalidates the candidate against the restored context
- **AND** returns the candidate to an apply-eligible state without requiring the operator to recreate it from scratch

### Requirement: Database planner organizes compatible tags into grouped rows
The database workbench SHALL organize compatible database candidates into grouped row plans keyed by connector context, table context, group key, and interval semantics.

#### Scenario: Prefix-inferred tags appear as one row plan
- **WHEN** the candidate set includes `meter/A1`, `meter/A2`, `meter/A3`, and `meter/kw`
- **THEN** the database planner shows one `meter` row group with member columns `a1`, `a2`, `a3`, and `kw`
- **AND** the operator is not forced to bind each tag one by one as the primary first-pass workflow

### Requirement: Database grouped rows expose editable grouping and interval controls
The database workbench SHALL let the operator review and edit grouped row metadata before apply, including `group_key`, `column_name`, and effective interval.

#### Scenario: Operator overrides grouped row metadata before apply
- **WHEN** an inferred grouped row needs a different group key, column name, or interval
- **THEN** the planner lets the operator edit those values before apply
- **AND** the chosen overrides are preserved in the apply request

#### Scenario: Interval edits validate grouped compatibility
- **WHEN** the operator changes the effective interval for a grouped row
- **THEN** the planner validates whether every member still shares the same effective interval
- **AND** the row stays blocked until the interval mismatch is resolved

#### Scenario: Incompatible grouped members stay blocked
- **WHEN** a proposed grouped row mixes members with different connector, schema, table, write mode, timestamp column, or effective interval
- **THEN** the planner marks that row as blocked with an explicit reason
- **AND** dry-run/apply do not silently split or merge the incompatible members

### Requirement: Legacy single-member mappings remain first-class
The database workbench SHALL preserve legacy mappings and slashless tag keys as single-member rows when no group key is set.

#### Scenario: Slashless tag stays single-member
- **WHEN** a database candidate has no slash in its tag key or an existing mapping persists with `group_key = null`
- **THEN** the planner keeps it as a single-member row
- **AND** grouped row logic does not merge it with unrelated mappings
