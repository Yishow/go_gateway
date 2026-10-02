## ADDED Requirements

### Requirement: Production durable intake and recoverable closure
The cmd/test_ui production path SHALL acknowledge new-group samples only after local durable journal commit. Finalized payload, frozen destination/group revisions, outbox intent and consumed checkpoint MUST commit atomically. Open buckets SHALL be recoverable from durable input. The existing delivery components SHALL be reused or extended instead of introducing an independent queue.

#### Scenario: Crash after sample acknowledgement
- **WHEN** the process dies after ACK but before row finalization
- **THEN** restart restores the acknowledged input and deterministically finalizes at most one record identity.

#### Scenario: Crash during row transaction
- **WHEN** the process dies between payload/outbox/checkpoint writes
- **THEN** restart observes either the complete committed transition or the prior state, never a consumed sample without its output intent.

#### Scenario: SQL failure
- **WHEN** a destination write fails
- **THEN** the only pending copy is not removed before Exec and the durable record remains recoverable.

### Requirement: Idempotent external effects and explicit uncertainty
The sender SHALL atomically persist a target row and target-side effect identity or receipt where the verified target supports it, and reconcile the same identity after uncertain completion. Local receipts alone MUST NOT be treated as proof of external commit. A target without safe deduplication SHALL block uncertain outcomes without blind replay. No universal exactly-once claim SHALL be made.

#### Scenario: Commit acknowledgement lost
- **WHEN** the target commits a row but the response is lost
- **THEN** retry or restart finds the same durable effect identity and creates no second row.

#### Scenario: Identity payload mismatch
- **WHEN** an existing effect key has a different payload digest
- **THEN** delivery is blocked with a safe conflict and does not overwrite the existing effect.

#### Scenario: Uncertain custom table
- **WHEN** a custom target lacks verifiable deduplication and the write response is lost
- **THEN** status is unknown/blocked pending reconciliation, not success or automatic reinsert.

### Requirement: Bounded isolated recovery and capacity
Delivery SHALL bound concurrency, retry/backoff, batch sizes and shutdown time, with durable exclusive claims or lease/fencing. One destination failure MUST NOT block healthy independent destinations or independently valid Share output. Retry exhaustion SHALL retain blocked data. Quota exhaustion MUST stop new ACKs and expose affected scope and loss risk without deleting accepted pending records.

#### Scenario: Outage and restart
- **WHEN** database A is offline while database B and Share are healthy and the gateway restarts
- **THEN** A retains and resumes its backlog while healthy outputs continue within configured bounds.

#### Scenario: Disk full
- **WHEN** local durable storage reaches its configured hard limit
- **THEN** new intake is explicitly blocked without a success ACK; accepted journal/outbox evidence is retained.

#### Scenario: Poison row
- **WHEN** one row has a permanent SQL type or constraint failure
- **THEN** it is quarantined with a safe reason, its ordered group/entity successors wait for explicit resolution, and other partitions continue.

#### Scenario: Overlapping workers
- **WHEN** a stale worker and a new worker claim the same partition after restart
- **THEN** transactional ownership/fencing prevents conflicting progress and shutdown respects its deadline.

### Requirement: Revision-bound backlog and truthful stages
Accepted records SHALL retain their frozen destination, schema and group revisions after later edits or disable. Runtime status, API and SSE MUST distinguish collecting, local_durable, queued/retrying, blocked/unknown and sql_committed, with verified reserved for actual readback. A nil buffer-return error or healthy device MUST NOT imply target commit.

#### Scenario: Endpoint edited with backlog
- **WHEN** a connector is changed to a new endpoint while old records are pending
- **THEN** new intake uses the new applied revision and old records remain tied to the original destination, blocked if its credentials cannot be safely resolved.

#### Scenario: Buffered but not committed
- **WHEN** the local queue accepts a row while its database is offline
- **THEN** UI/API show local durable and pending/failure state rather than delivered.

#### Scenario: Disable group
- **WHEN** an operator disables a group after records were acknowledged
- **THEN** new intake stops while accepted backlog remains queryable and is not silently discarded or redirected.

### Requirement: Single production writer owner and lifecycle backlog
When a migrated or new group is applied, exactly one writer path SHALL own new intake for that output, with Share independently gated. Rolling back the UI or binary SHALL stop new intake while preserving accepted records and their destination provenance. Deleting a group with accepted pending records SHALL let production workers complete the original accepted records without orphaning, replaying or dropping them.

#### Scenario: Activation switches writer
- **WHEN** a migrated group is applied while the legacy target exists
- **THEN** exactly one writer path owns new intake for that output and Share remains independently gated.

#### Scenario: Rollback with accepted backlog
- **WHEN** the new UI or binary is rolled back after durable records were accepted
- **THEN** new intake can be stopped but accepted records and their destination provenance are preserved.

#### Scenario: Delete with backlog is completed by workers
- **WHEN** the operator deletes a group with accepted pending records
- **THEN** production workers deliver the original accepted records from the tombstoned group's frozen revisions and destination without orphaning, replaying or dropping them.
