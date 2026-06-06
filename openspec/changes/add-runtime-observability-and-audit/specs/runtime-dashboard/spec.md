## ADDED Requirements

### Requirement: Runtime dashboard shows operator-facing diagnostics

The runtime dashboard SHALL show operator-facing diagnostics for recent runtime and delivery failures.

#### Scenario: Dashboard shows latest failure context

- **WHEN** runtime or database delivery has a recent failure for the selected device scope
- **THEN** the dashboard shows the latest failure context and timestamp for that selected scope
- **AND** the operator can see more than aggregate counters alone

##### Example: dashboard shows latest DB delivery failure

- **GIVEN** dev-A has a recent failed database delivery at 10:12:00Z
- **WHEN** the runtime dashboard opens for dev-A
- **THEN** the dashboard shows that latest delivery failure context and timestamp for dev-A
