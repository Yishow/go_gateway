## MODIFIED Requirements

### Requirement: Truthful explicit test writes
The selected canonical WriteGroup SHALL be the authority for source membership and target revisions. Legacy plan routes MUST resolve to that same group without guessing IDs or creating a second operation ledger. Basic Tag test writes MUST NOT require advanced measurement semantics.
The system SHALL separate data preview, confirmed write and readback verification and retain idempotent test identities. Test-write preview SHALL issue a durable token and operation_id bound to the test_write action, saved target, setup revisions, payload digest and expiry. Confirmation MUST validate that binding before the first write; a schema token or plan_id alone MUST NOT authorize test writing. Confirmed writes MUST use the selected persisted target and compare the actual test payload during readback. Results SHALL distinguish written_verified, written_unverified, failed and unknown, with an independent cleanup_status of not_attempted, cleaned, failed or unknown. The UI MUST NOT infer verified writing, cleanup or timestamps without corresponding backend evidence.

#### Scenario: Test preview and action-bound confirmation
- **WHEN** an operator requests a test-write preview for a saved, schema-verified group (or a legacy plan resolved to that group)
- **THEN** the backend returns previewed test content and its persisted token and operation_id without inserting or cleaning target data
- **AND** a subsequent confirmation with a schema token, stale revision or mismatched scope cannot perform a test write.

#### Scenario: Write succeeds without read permission
- **WHEN** a confirmed test write succeeds but the account cannot query the result
- **THEN** the result is written_unverified, not verified.

#### Scenario: Test action retried
- **WHEN** the same confirmed test token is retried after an unknown network outcome
- **THEN** the system checks its durable identity and does not insert a duplicate test record; a running duplicate returns 202, a retained result returns 200, and a different operation occupying the scope returns 409
- **AND** the client can query that operation through the shared workspace database-operation status endpoint; missing and foreign operations return the same safe 404.

#### Scenario: Readback payload differs
- **WHEN** the record identity is found but its payload does not match the written test values
- **THEN** the result does not claim verified writing and reports a safe mismatch reason.

#### Scenario: Cleanup cannot be confirmed
- **WHEN** writing and readback succeed but cleanup fails or has an uncertain result
- **THEN** write verification and cleanup status are displayed independently without claiming the target was cleaned.

#### Scenario: Cleanup and retry after restart
- **WHEN** a previously completed test is retried after cleanup or a process restart
- **THEN** the retained operation identity prevents a new insertion
- **AND** cleanup never deletes records outside that operation's test scope.

#### Scenario: Write outcome or status is unknown
- **WHEN** the backend cannot confirm whether a write happened or the client receives an unrecognized status
- **THEN** the UI shows an unresolved result rather than success and does not automatically repeat a potentially completed write.

## ADDED Requirements

### Requirement: Write-group test operations use production delivery
Test-write previews SHALL be mutation-free and preserve action kind, operation ID, content digest, scope, expiry and expected workspace/group/connector revisions in the existing durable operation ledger. Confirmation MUST atomically claim execution once and use the production row layout, row codec and destination insert path (including the group dedupe strategy); it MUST NOT route the test row through the delivery outbox, so it never enters group backlog, quota or partition ordering, and its result waits for destination commit evidence. The system SHALL expose saved operation results separately from cleanup status.

#### Scenario: Action-bound token and scopes
- **WHEN** a client confirms with a schema-create token, stale revision, expired unclaimed token, missing fields or foreign operation
- **THEN** the request fails before target mutation with respectively safe kind/conflict/validation/resource errors.

#### Scenario: Pending or completed resubmission
- **WHEN** the same operation is retried during execution or after a persisted outcome
- **THEN** running returns 202 and saved results return 200 without a new target write, even if its token has since expired.

#### Scenario: Real readback
- **WHEN** the row has committed and can be read
- **THEN** the system compares typed values, scope, identity and provenance before written_verified; a mismatch is written_unverified with a reason.

#### Scenario: Owned cleanup
- **WHEN** cleanup runs with neighboring production records in the same table
- **THEN** only the operation-owned test rows and metadata are removed; a durable receipt remains and neighboring production data is unchanged.

#### Scenario: Missing read or delete permission
- **WHEN** the row commits but SELECT or DELETE is forbidden
- **THEN** the response truthfully separates written_unverified from cleanup failed/not_attempted; it never claims verified and clean.

#### Scenario: Restart after ambiguous write
- **WHEN** the gateway restarts after target commit but before saving or returning the result
- **THEN** it reconciles that operation and target receipt rather than issuing another preview or inserting again.

#### Scenario: Unsafe test ownership
- **WHEN** a target cannot safely distinguish operation-owned test data
- **THEN** preview rejects the unsupported test-write capability without mutating or broad DELETE cleanup.
