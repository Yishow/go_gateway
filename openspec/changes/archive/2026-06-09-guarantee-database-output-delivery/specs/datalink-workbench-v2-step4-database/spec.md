## ADDED Requirements

### Requirement: Step 4 shows database connector readiness and delivery truth

The Step 4 database workspace SHALL show connector readiness, schema ensure outcome, and delivery truth for the current database output scope.

#### Scenario: Step 4 surfaces delivery failure truth

- **WHEN** the selected database connector or target scope has a failed schema ensure or failed write delivery outcome
- **THEN** Step 4 shows that failed delivery truth with the affected scope
- **AND** the operator can distinguish setup configuration from actual delivery health
