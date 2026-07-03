## ADDED Requirements

### Requirement: Database target workflow uses the live visible target set

The database target workflow SHALL use the live visible target set for planning, schema ensure, and apply behavior.

#### Scenario: Workflow excludes stale hidden mappings

- **WHEN** a database connector still contains stale hidden target mappings that are no longer visible from the current live rule-owned state
- **THEN** the workflow excludes those stale mappings from planning and delivery operations
- **AND** it SHALL NOT let a stale hidden mapping silently change or block the live target set

##### Example: planner excludes one orphaned hidden mapping

- **GIVEN** the current live planner shows targets for tags tag-A and tag-B while connector metadata still retains a hidden orphaned mapping for tag-old
- **WHEN** the workflow computes planning and delivery operations
- **THEN** tag-old is excluded from the live target set and cannot block tag-A or tag-B
