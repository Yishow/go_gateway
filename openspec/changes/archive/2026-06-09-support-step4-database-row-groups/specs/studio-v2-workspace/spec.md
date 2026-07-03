## ADDED Requirements

### Requirement: Workspace persists database row-group plans

The `/studio/v2` workspace SHALL persist and reload Step 4 database row-group plans alongside connector and target configuration.

#### Scenario: Autosave stores row-group metadata

- **WHEN** Step 4 autosave writes a database plan that includes row groups
- **THEN** the workspace persistence layer stores the row-group records and their target references in the same logical workspace
- **AND** a subsequent read returns the same row-group metadata without loss

#### Scenario: Legacy workspace remains compatible

- **WHEN** the workspace contains only legacy single-row database targets with no row-group metadata
- **THEN** load and autosave continue to work without requiring row-group fields
- **AND** the workspace SHALL NOT invent row-group ids unless the operator explicitly enables shared-column planning
