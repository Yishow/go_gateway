## ADDED Requirements

### Requirement: Live apply reports runtime reconciliation outcome

The system SHALL report runtime reconciliation outcome for persisted changes applied to a running workspace.

#### Scenario: Live apply response includes reconcile outcome

- **WHEN** a running workspace accepts a persisted config change that affects runtime state
- **THEN** the live apply response includes the reconciliation outcome for that change scope
- **AND** the operator can tell whether runtime is already aligned, deferred, or requires restart or reactivation

##### Example: device host change returns restart-required

- **GIVEN** dev-A is running and the operator persists a host or port change for dev-A
- **WHEN** the live apply response is returned
- **THEN** the response includes restart-required for dev-A if runtime cannot safely reconnect in place
