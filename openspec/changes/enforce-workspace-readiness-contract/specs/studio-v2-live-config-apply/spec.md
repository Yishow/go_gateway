## ADDED Requirements

### Requirement: Live apply honors the workspace readiness contract

The system SHALL allow live apply only when the affected persisted workspace state satisfies the readiness contract for the requested change scope.

#### Scenario: Live apply rejects readiness blockers

- **WHEN** a running workspace receives a persisted config change that introduces a blocking readiness issue
- **THEN** the system rejects or defers live apply for that change scope
- **AND** the operator receives the normalized readiness issue instead of a silent runtime drift

##### Example: live apply defers a blocking source-rule change

- **GIVEN** a running workspace and a persisted rule change that removes the last valid mapping for pt-A
- **WHEN** the system evaluates live apply for that change
- **THEN** the response returns the blocking readiness issue for pt-A and marks the runtime reconcile outcome as deferred or rejected
