## ADDED Requirements

### Requirement: Step 4 row-group planner exposes shared-table structure

Step 4 SHALL expose row groups as explicit planning units whenever operators configure multiple rules to write into one table with reused business columns.

#### Scenario: Rule groups are visible before apply

- **WHEN** multiple rules target the same connector and table in Step 4
- **THEN** the UI shows which points belong to each row group before apply
- **AND** the operator can distinguish "same table, different row group" from "same table, same row group"

#### Scenario: Shared-column plan stays scoped to one table

- **WHEN** an operator edits a row group in Step 4
- **THEN** the planner keeps every member scoped to the currently selected connector and table
- **AND** the planner SHALL NOT reuse a row-group definition across a different connector or table silently

### Requirement: Step 4 shared-column validation distinguishes row-group conflicts

Step 4 SHALL validate repeated business columns by row-group scope instead of treating every repeated column as a generic conflict.

#### Scenario: One row group can reuse one business column

- **WHEN** several points inside the same row group bind to the same business column
- **THEN** the planner marks the configuration as valid if the row group identity contract is complete

#### Scenario: Cross-group duplicate remains blocking

- **WHEN** repeated business columns span different row groups
- **THEN** the planner shows a blocking conflict
- **AND** the conflict message identifies the row-group boundary that makes the reuse illegal
