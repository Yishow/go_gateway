## MODIFIED Requirements

### Requirement: Commit sequence and animation
Step 4 SHALL display activation progress derived from acknowledged backend operations for the submitted workspace revision. Timers MAY pace presentation but MUST NOT create successful operation evidence. Repeated activation clicks SHALL be blocked while an attempt is in flight. Remounting SHALL recover server state rather than resume a simulated success sequence.

#### Scenario: Presentation timer completes without a response
- **WHEN** all presentation timers elapse while the backend has not confirmed activation
- **THEN** the UI remains pending or unconfirmed
- **AND** no collecting, delivered or verified state is inferred from elapsed time.

#### Scenario: Acknowledged partial activation
- **WHEN** the server confirms some devices and rejects others
- **THEN** progress identifies each confirmed and rejected outcome
- **AND** it does not report the entire workspace as successful.

#### Scenario: Remount during activation
- **WHEN** the operator returns to Step 4 during an outstanding activation
- **THEN** the client reads authoritative status and preserves the operation identity when available
- **AND** it does not dispatch duplicate activation solely because it remounted.

### Requirement: Commit completion card
Step 4 SHALL distinguish saved configuration, applied revision, active collection and database delivery evidence. A successful activation card MUST require a confirmed activation outcome and MUST NOT promise a first write solely from an interval. Partial and unconfirmed outcomes SHALL remain visible alongside any available runtime-navigation action.

#### Scenario: Configuration saved without activation
- **WHEN** a draft was saved but no activation was confirmed
- **THEN** the UI reports saved configuration rather than collecting or delivered data.

#### Scenario: Activation succeeds without a database receipt
- **WHEN** activation is confirmed but no database delivery evidence exists
- **THEN** the UI may report active collection but SHALL display database delivery as not yet confirmed.

#### Scenario: Handoff includes resolved device id
- **WHEN** runtime navigation is available and one confirmed device is resolved
- **THEN** the navigation callback runs once and opens /studio/runtime with that device_id.

#### Scenario: Handoff falls back when device cannot be resolved
- **WHEN** runtime navigation is available without one unambiguous device
- **THEN** the callback opens /studio/runtime without inventing a device identifier
- **AND** the action is not a no-op or console-only side effect.
