## MODIFIED Requirements

### Requirement: Traceable scoped delivery head resolution
The system SHALL expose safe blocked/quarantined row attention and explicit retry or confirmed skip within the owning workspace and group. Resolution SHALL atomically record decision provenance and validate expected state, existing claim epoch and payload digest. Retry SHALL preserve effect identity, payload and frozen destination. Skip SHALL retain the record and identify it as undelivered. Unknown effects SHALL reject retry and skip.

#### Scenario: Repair blocked head
- **WHEN** a missing table or permission failure is repaired and the operator retries that blocked head
- **THEN** the same effect is retried against its original destination and following rows become deliverable after settlement.

#### Scenario: Skip poison row
- **WHEN** an operator explicitly confirms skipping a quarantined bad row
- **THEN** its payload and resolution provenance remain available and the next row is unblocked.

#### Scenario: Unsafe or stale resolution
- **WHEN** a foreign workspace/group effect, stale expected state or epoch (including blocked → retry → blocked), or unknown effect is submitted for resolution
- **THEN** the mutation is rejected without resetting delivery or changing provenance.

#### Scenario: Repeated decision and new blocked attempt
- **WHEN** a decision reply is uncertain, or that effect is blocked again after a new delivery attempt
- **THEN** the identical decision ID reads its atomic audit without resending, while a new decision requires the latest state revision and the UI permits a fresh repair note.

#### Scenario: Repair action is primary and Skip impact is explicit
- **WHEN** a blocked or quarantined row is presented
- **THEN** retry after repair is the primary action; Skip is secondary and expands to explain that this record will no longer be delivered, require a reason and explicit confirmation, and preserve per-record audit and decision_id.

