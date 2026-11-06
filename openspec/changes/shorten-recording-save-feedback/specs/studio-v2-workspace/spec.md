## ADDED Requirements

### Requirement: Configuration save acknowledgement is separate from refresh

The Studio database configuration save flow SHALL acknowledge a submitted configuration only from the existing valid server acknowledgement for its scope and revisions. Independent query refreshes SHALL NOT delay that acknowledgement. Saved configuration, fresh readback, activation, SQL commit and cleanup SHALL remain distinct states.

#### Scenario: Valid acknowledgement precedes refresh

- **WHEN** a valid configuration ACK arrives while database configuration, targets and bootstrap refreshes remain pending
- **THEN** that submitted configuration is shown as saved and dependent query state remains refreshing
- **AND** no readiness, schema, lifecycle or test-write operation is marked complete by that ACK

##### Example: deferred refreshes

- **GIVEN** a test ACK settles at 1 ms and each of three independent refresh promises takes 200 ms
- **WHEN** the valid configuration ACK is accepted
- **THEN** acknowledgement does not wait for the 600 ms serial refresh chain and all three independent refreshes can be pending together
- **AND** the UI does not label those reads fresh before they complete

#### Scenario: Acknowledgement lacks confirmation evidence

- **WHEN** a response lacks the scope or revision evidence required by the existing acknowledgement contract
- **THEN** the flow retains its confirmation-read or uncertain state
- **AND** HTTP success alone does not mark the configuration saved

### Requirement: Background refresh remains scoped and observable

The configuration save flow SHALL invalidate affected query evidence before publishing a saved result, run independent refreshes concurrently, coalesce same-key demand, and expose refresh errors. Stale responses MUST NOT overwrite newer revisions, other workspaces or newer unsaved edits.

#### Scenario: Refresh fails after confirmed save

- **WHEN** a background refresh fails after a valid configuration save ACK
- **THEN** the UI retains the saved result and separately reports stale refresh failure with read-only retry
- **AND** operations requiring fresh evidence remain unavailable

#### Scenario: New edit or workspace wins

- **WHEN** a user edits again, changes workspace, unmounts the editor, or receives an older refresh after a newer ACK
- **THEN** only the matching submitted snapshot is acknowledged and current newer edits or scope remain unchanged
- **AND** outdated data cannot hydrate over the new state

### Requirement: Save acceleration preserves uncertain and consequential operation barriers

Only the named database configuration save flow SHALL adopt the separated ACK contract. Existing target saves, write-group lifecycle, schema apply, activation and test-write flows SHALL retain their completion/readback requirements. Uncertain writes SHALL preserve the original operation identity and MUST NOT be automatically replayed.

#### Scenario: Timeout conflict or repeated action

- **WHEN** a save times out, returns a CAS conflict, loses its ACK, or receives repeated clicks while unresolved
- **THEN** existing failed or uncertain recovery and dirty-navigation barriers remain in force
- **AND** the client does not invent a successful result or automatically submit a replacement operation

#### Scenario: Another operation uses shared helpers

- **WHEN** a target, schema, lifecycle, activation or test-write flow awaits its existing reconciliation or readback
- **THEN** the configuration refresh optimization does not shorten that flow's success barrier
- **AND** Close, Cancel and Back/Forward preserve the established sent, unsent and unknown operation semantics
