## ADDED Requirements

### Requirement: Durable legacy writer takeover
The system SHALL persist connector/tag ownership at the first effective canonical takeover before new canonical intake. Disable, delete, retirement, destination edits and restart SHALL NOT reactivate the old legacy writer. Before a successful takeover, legitimate legacy delivery SHALL continue. Ownership read or persistence failure SHALL fail closed for canonical takeover.

#### Scenario: Retire and restart after migration
- **WHEN** an enabled legacy target is migrated, effectively taken over, disabled and its boundary retired, then the service restarts
- **THEN** the legacy target remains suppressed and accepted canonical backlog retains its original destination.

#### Scenario: Delete or change destination after takeover
- **WHEN** a taken-over group is deleted or changed to another connector
- **THEN** the former connector/tag legacy output remains suppressed.

#### Scenario: No effective takeover
- **WHEN** a draft or failed apply has not successfully taken over an output
- **THEN** the original enabled legacy writer continues legitimate delivery.

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
