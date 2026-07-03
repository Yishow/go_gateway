## ADDED Requirements

### Requirement: Workspace audit history records critical lifecycle events

The system SHALL record critical workspace lifecycle events in audit history.

#### Scenario: Activation writes audit history entry

- **WHEN** an operator starts workspace activation and the system returns success, partial success, or failure
- **THEN** the system records an audit history entry for that activation attempt with result and scope
- **AND** the entry remains queryable after browser refresh and service restart

##### Example: partial activation is still recorded

- **GIVEN** activation succeeds for dev-A and fails for dev-B in workspace ws-1
- **WHEN** the activation response is returned
- **THEN** audit history records one activation event for ws-1 with partial-success result and device scope details

### Requirement: Workspace audit history records delivery-impacting configuration changes

The system SHALL record persisted configuration changes that affect runtime or database delivery.

#### Scenario: Persisted Step 4 change records audit entry

- **WHEN** an operator persists a Step 4 database connector or target change that affects delivery behavior
- **THEN** the system records an audit history entry for that persisted change with result and affected scope
- **AND** later diagnostics can correlate delivery issues with that configuration change
