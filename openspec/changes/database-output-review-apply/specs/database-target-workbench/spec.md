## MODIFIED Requirements

### Requirement: Database output workflow is layered by connector, schema context, and mapping
The workbench SHALL present database output as four linked layers: source-rule revision context, connector selection, schema/table context, and candidate-to-column review/apply.

#### Scenario: Connector selection scopes schema context and candidate board
- **WHEN** an operator selects a database connector in the Output step for a source rule
- **THEN** the available schema and table context SHALL scope to that connector only
- **AND** the rule-scoped database candidate board SHALL update using only that connector context

#### Scenario: Apply uses current layered context
- **WHEN** an operator applies a database output candidate
- **THEN** the mapping SHALL be applied against the current source-rule revision and the selected connector, schema, table, and column context
- **AND** the UI SHALL NOT apply the action to stale prior selections

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
