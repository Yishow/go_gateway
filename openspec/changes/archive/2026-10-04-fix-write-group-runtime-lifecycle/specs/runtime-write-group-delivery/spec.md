## MODIFIED Requirements

### Requirement: Production durable intake and recoverable closure
The cmd/test_ui production path SHALL acknowledge new-group samples only after local durable journal commit. Finalized payload, frozen destination/group revisions, outbox intent and consumed checkpoint MUST commit atomically. Open buckets SHALL be recoverable from durable input, including disabled, tombstoned and superseded revisions. Verified applied layouts MUST support local intake recovery while the destination is offline. The existing delivery components SHALL be reused or extended instead of introducing an independent queue.

#### Scenario: Crash after sample acknowledgement
- **WHEN** the process dies after ACK but before row finalization
- **THEN** restart restores the acknowledged input and deterministically finalizes at most one record identity.

#### Scenario: Crash during row transaction
- **WHEN** the process dies between payload/outbox/checkpoint writes
- **THEN** restart observes either the complete committed transition or the prior state, never a consumed sample without its output intent.

#### Scenario: SQL failure
- **WHEN** a destination write fails
- **THEN** the only pending copy is not removed before Exec and the durable record remains recoverable.

#### Scenario: Destination offline during cold start
- **GIVEN** an eligible applied version with a locally persisted verified layout and available local quota
- **WHEN** the gateway restarts while the destination cannot be reached
- **THEN** new matching samples can be durably acknowledged without a remote schema query; delivery waits and later uses the original frozen destination.

#### Scenario: Restart drains historical accepted input
- **GIVEN** acknowledged journal input that has not closed, followed by Disable, Delete or an applied revision replacement
- **WHEN** the process dies and restarts
- **THEN** each historical version resumes closure using its frozen layout and persisted cutoff, accepts no new samples outside its responsibility and remains queryable until its accepted data is resolved.

#### Scenario: Recovery metadata cannot be proven
- **WHEN** an old revision has no trustworthy locally persisted layout and its original source or destination cannot be verified
- **THEN** it remains explicitly blocked with retained journal and repair guidance; no current draft, changed source type or different endpoint is substituted.
