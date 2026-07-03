## ADDED Requirements

### Requirement: Workspace database row groups capture shared-column intent

The workspace SHALL persist database row groups as first-class Step 4 planning entities whenever an operator wants multiple points to reuse the same business columns in one table.

#### Scenario: Persist one shared-column row group

- **WHEN** an operator creates a row group for one connector and one table, assigns multiple points to it, and saves Step 4
- **THEN** the workspace persists a stable row-group record with its id, connector/table scope, member points, and group-key metadata
- **AND** the member database targets persist their referenced row-group id instead of being treated as unrelated single-row targets

#### Scenario: Reload row-group planning on next workspace load

- **WHEN** `/studio/v2` reloads a workspace that already contains persisted row groups
- **THEN** Step 4 restores the same row-group structure and shared-column memberships
- **AND** the workspace SHALL NOT silently flatten those members back into global unique-column bindings

### Requirement: Row groups define the legality of shared business columns

The system SHALL treat repeated business columns as legal only when every repeated binding belongs to the same row group and that row group defines a valid row identity contract.

#### Scenario: Shared column stays legal inside one row group

- **WHEN** two points inside the same row group both bind to `temperature_c`
- **THEN** validation marks the plan as legal if the row group contains a valid row identity contract
- **AND** the plan SHALL NOT surface a generic duplicate-column blocker

#### Scenario: Shared column remains blocking across row groups

- **WHEN** two points in different row groups bind to `temperature_c` for the same connector/table scope
- **THEN** validation reports a blocking conflict
- **AND** the conflict SHALL identify the cross-group reuse rather than a generic single-point column error

### Requirement: Row-group write mode guardrails protect unsafe upsert reuse

The system SHALL block row-group configurations that reuse business columns under `upsert` unless the row group defines a stable uniqueness key that can distinguish row instances.

#### Scenario: Insert-mode shared-column row group is allowed

- **WHEN** a row group reuses business columns and the connector write mode is `insert`
- **THEN** readiness allows the plan to proceed if all other validation passes

#### Scenario: Upsert without uniqueness key is rejected

- **WHEN** a row group reuses business columns and the connector write mode is `upsert` without a stable uniqueness key
- **THEN** readiness reports a blocking issue
- **AND** apply SHALL NOT continue until the operator supplies a valid uniqueness contract or changes the write mode
