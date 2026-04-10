## ADDED Requirements

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
