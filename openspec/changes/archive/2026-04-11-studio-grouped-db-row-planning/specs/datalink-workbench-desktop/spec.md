## ADDED Requirements

### Requirement: Source planning surfaces grouped-tag handoff context

The workbench SHALL show, in Step 2, enough Source planning context for the operator to understand how planned points will move into grouped Tag review and downstream Database planning.

#### Scenario: Source summary explains grouped handoff
- **WHEN** an operator plans source points for a rule
- **THEN** Step 2 shows the planned point count and current naming-prefix context
- **AND** the handoff copy explains that the next step will review grouped tags before database planning

### Requirement: Tag review exposes grouped database suggestions before Output

The workbench SHALL show, in Step 3, the inferred database grouping context for rule-derived tags before the operator moves into Database planning.

#### Scenario: Tag review shows inferred grouping metadata
- **WHEN** the operator opens Step 3 for slash-based tags such as `meter/A1` and `meter/kw`
- **THEN** the Tag review surface shows inferred `group_key`, suggested `column_name`, and group membership for each candidate
- **AND** the operator can understand which tags are expected to land in the same database row

#### Scenario: Tag overrides carry forward to Output planning
- **WHEN** the operator edits grouped row suggestions in Step 3
- **THEN** the chosen group and column overrides are preserved for the Database planner in Output
- **AND** the operator does not have to re-enter the same override state from scratch

### Requirement: Output mainline completes when either target is configured

The workbench SHALL treat the Output step as mainline-complete when either Local Modbus or Database is configured, while still keeping both target states visible in the same workspace.

#### Scenario: Database apply satisfies the mainline
- **WHEN** at least one database output mapping has been applied for the active source rule and Local Modbus is still incomplete
- **THEN** the Output step is reported as complete for mainline progression
- **AND** the UI still shows Local Modbus as incomplete instead of hiding it

#### Scenario: Local Modbus apply satisfies the mainline
- **WHEN** at least one Local Modbus binding has been applied for the active source rule and Database is still incomplete
- **THEN** the Output step is reported as complete for mainline progression
- **AND** the UI still shows Database planning status and required follow-up work

### Requirement: Recovery cues point to the true blocked surface

The workbench SHALL align the step rail, context bar, and shell recovery action to the same blocker truth after grouped Output planning changes.

#### Scenario: Locked steps explain blockers without acting as a second router
- **WHEN** a downstream step is blocked by missing grouped Tag review or Output configuration
- **THEN** the step rail explains the blocking reason without behaving like an alternate forward-action router
- **AND** the shell return action points to the actual recovery surface
