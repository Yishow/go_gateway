## ADDED Requirements

### Requirement: Rule lifecycle changes trigger runtime reconciliation results

The system SHALL turn source-rule lifecycle changes into explicit runtime reconciliation results.

#### Scenario: Rule change reports reconcile result

- **WHEN** a persisted source rule is enabled, disabled, updated, or removed while the workspace is running
- **THEN** the system records and returns the runtime reconciliation result for that rule scope
- **AND** the operator SHALL NOT have to infer from stale runtime behavior whether the rule change actually took effect

##### Example: disabling one running rule returns immediate reconcile

- **GIVEN** rule-A is currently active in runtime for dev-A
- **WHEN** the operator persists a disable action for rule-A
- **THEN** the system returns the reconcile result for rule-A instead of leaving the operator to guess whether collection stopped
