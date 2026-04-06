## ADDED Requirements

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
