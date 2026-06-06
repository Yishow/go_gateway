## ADDED Requirements

### Requirement: Rule-derived database target drift does not block unrelated live runtime

The system SHALL keep stale rule-derived database target drift from blocking unrelated live runtime and database delivery scope.

#### Scenario: Orphaned database target is isolated from live runtime scope

- **WHEN** a rule-derived database target becomes orphaned from its live rule-owned relationship
- **THEN** the system isolates that orphaned target from live schema ensure and write delivery
- **AND** unrelated live rule-owned targets continue to activate and deliver normally

##### Example: removed rule leaves one orphaned target behind

- **GIVEN** rule-A was removed, target row-A is now orphaned, and rule-B still owns live target row-B
- **WHEN** activation and delivery evaluate current runtime scope
- **THEN** row-A is isolated and row-B continues to activate and deliver normally
