## ADDED Requirements

### Requirement: Runtime projection is built from persisted workspace state

The system SHALL build the runtime device, rule, mapping, and output projection from persisted workspace state during activation and service restart.

#### Scenario: Activation and restart build the same projection

- **WHEN** an operator activates a persisted workspace and later the service restarts
- **THEN** runtime rebuilds the same device, rule, mapping, and output projection from persisted workspace state in both cases
- **AND** runtime SHALL NOT depend on transient browser session state to restore collection behavior

##### Example: one workspace yields the same two-device projection twice

- **GIVEN** persisted workspace ws-1 contains devices dev-A and dev-B, rules rule-A and rule-B, and enabled mappings for both devices
- **WHEN** ws-1 is activated and later the service restarts
- **THEN** runtime projects dev-A and dev-B with the same persisted rule and mapping set in both runs

### Requirement: Persisted config changes produce an explicit reconciliation outcome

The system SHALL return an explicit reconciliation outcome whenever a persisted config change affects running runtime state.

#### Scenario: Persisted change reports reconcile outcome

- **WHEN** a persisted device, source rule, mapping, or database target change affects a running workspace
- **THEN** the system returns whether the change was reconciled immediately, deferred, or requires restart or reactivation
- **AND** the system SHALL NOT silently accept the persisted change while continuing to run an older projection

##### Example: mapping change returns deferred outcome

- **GIVEN** workspace ws-1 is running and a persisted mapping change affects a hot path that cannot be reloaded in place
- **WHEN** the system applies that persisted change
- **THEN** the response marks the reconcile outcome as deferred or restart-required instead of pretending runtime is already aligned

### Requirement: Reconciliation is targeted and idempotent

The system SHALL reconcile only the affected runtime scope and SHALL keep repeated reconciliation requests idempotent.

#### Scenario: Same reconcile request does not duplicate runtime state

- **WHEN** the same persisted change or reconciliation request is processed more than once
- **THEN** runtime converges to one consistent projection for the affected scope
- **AND** the system SHALL NOT duplicate scheduler registrations, mapping bindings, or output bindings

##### Example: repeated reconcile for dev-A does not double-register

- **GIVEN** dev-A already has one active scheduler registration and one active mapping projection
- **WHEN** the same reconcile request for dev-A is replayed twice
- **THEN** runtime keeps one scheduler registration and one mapping projection for dev-A
