## ADDED Requirements

### Requirement: Backend-backed V2 settings surface

The `/studio/v2/settings` surface SHALL boot from real backend state instead of demo-only local defaults.

#### Scenario: Settings boot from backend

- **WHEN** the operator opens `/studio/v2/settings`
- **THEN** the page loads settings and connector data from backend APIs
- **AND** the page does not silently treat local defaults as a completed backend load

### Requirement: Real connector pool operations

The connector pool SHALL use real backend CRUD and test APIs.

#### Scenario: Real connector test

- **GIVEN** a connector row exists in the pool
- **WHEN** the operator clicks the test button
- **THEN** the system sends a real backend connector test request
- **AND** the result shown in the UI comes from backend response, not Math.random or a mock timer

### Requirement: Save bar persists settings

The `儲存所有設定` action SHALL persist settings through backend APIs.

#### Scenario: Save settings

- **GIVEN** the operator has changed one or more settings fields
- **WHEN** the operator clicks `儲存所有設定`
- **THEN** the system writes the changed settings to backend
- **AND** save failure is reported as an actionable error instead of a warning-only noop
