## ADDED Requirements

### Requirement: Runtime dashboard backend exposes diagnostics summary

The runtime dashboard backend SHALL expose diagnostics summary and latest failure context for the selected runtime scope.

#### Scenario: Backend returns diagnostics summary for selected device

- **WHEN** a client requests runtime monitoring data for a selected device
- **THEN** the backend response can include diagnostics summary with latest success and failure context for that device scope
- **AND** the client SHALL NOT need to parse raw server logs to present that summary

##### Example: diagnostics summary returns latest failure for dev-A

- **GIVEN** dev-A last failed at 10:12:00Z during database delivery and last succeeded at 10:10:00Z
- **WHEN** the client requests runtime monitoring data for dev-A
- **THEN** the response includes both timestamps and the failure context for dev-A
