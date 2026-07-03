## ADDED Requirements

### Requirement: Database delivery operates only on the live enabled mapping set

The system SHALL build schema ensure and write delivery from the live enabled database-target mapping set only.

#### Scenario: Hidden stale target does not block live delivery

- **WHEN** a connector still has an old hidden stale database target mapping that is no longer tied to a live rule-owned tag or point
- **THEN** schema ensure and runtime write delivery ignore that stale mapping
- **AND** the stale mapping SHALL NOT block unrelated live database targets from activation or write delivery

##### Example: one stale hidden target is ignored while one live target proceeds

- **GIVEN** connector db-main still contains stale target row-old for missing tag tag-old and live target row-live for tag-live
- **WHEN** schema ensure and writer evaluate db-main
- **THEN** row-old is ignored and row-live remains eligible for activation and delivery

### Requirement: Database delivery records last schema and write outcomes

The system SHALL record the last schema ensure outcome and the last write delivery outcome for each relevant connector or target scope.

#### Scenario: Failed write updates delivery outcome

- **WHEN** runtime attempts a database write and the write fails
- **THEN** the system records the failed delivery outcome with timestamp and failure reason for the affected connector or target scope
- **AND** the next operator-facing surface can show that failure without scraping raw logs

##### Example: permission failure updates connector delivery truth

- **GIVEN** connector db-main is selected and a write fails with a permission error at 10:03:00Z
- **WHEN** the failed write outcome is recorded
- **THEN** db-main reports the failed delivery outcome with 10:03:00Z and the latest failure reason

### Requirement: Runtime write delivery is diagnosable end to end

The system SHALL make the path from collected value to database write result diagnosable for the affected connector and target scope.

#### Scenario: Collected value reaches a failed write path

- **WHEN** runtime collects a value for a point whose database target write later fails
- **THEN** the system can report that the value was collected, mapped, and rejected by the database delivery path
- **AND** the operator SHALL NOT have to infer the failure only from missing rows in the external database

##### Example: collected point value fails only at DB write stage

- **GIVEN** pt-A is collected at 10:05:00Z, mapping succeeds, and dbtarget writer later fails for connector db-main
- **WHEN** the operator inspects delivery truth
- **THEN** the system reports pt-A as collected and mapped before failing in database delivery
