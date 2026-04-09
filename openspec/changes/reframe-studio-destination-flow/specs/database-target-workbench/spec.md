## ADDED Requirements

### Requirement: Database workspace is row-planner-first

The database destination workflow SHALL present delivery groups as the primary planning unit and SHALL provide a row planner as the main workspace.

#### Scenario: Multi-tag row planning
- **WHEN** an operator selects a delivery group for database delivery
- **THEN** the workspace helps compose one row from multiple tags
- **AND** does not require the operator to bind columns one tag at a time as the primary flow
