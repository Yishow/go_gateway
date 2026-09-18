## MODIFIED Requirements

### Requirement: Credential and path safety
The system SHALL reuse stored credentials by backend reference and invalidate tests after identity changes without returning secrets to the browser. The backend MUST verify access to the selected persisted connector, derive the database dialect from that connector, preserve supplied password bytes and distinguish unchanged, cleared and replaced credentials. Target names and prefixes MUST be validated and handled by the selected adapter rather than interpolated from untrusted input without validation.

#### Scenario: Reuse a saved connection
- **WHEN** the operator selects an existing unchanged connection
- **THEN** no password re-entry is required solely because the list response masks it
- **AND** changing the endpoint requires a new credential validation.

#### Scenario: Select a different database kind
- **WHEN** a persisted non-SQLite connector is selected for the plan
- **THEN** creation, preview, apply and test write use that connector and its server-derived dialect, not a fixed default connector or SQLite fallback.

#### Scenario: Password includes leading or trailing spaces
- **WHEN** the operator explicitly supplies a password containing spaces
- **THEN** its exact bytes are preserved and are not silently trimmed or exposed in normal diagnostics.

#### Scenario: Connection is absent or inaccessible
- **WHEN** the selected connection is unsaved, missing or outside the allowed workspace scope
- **THEN** preparation is blocked without guessing a connection identity or reusing credentials from another endpoint.

### Requirement: Truthful explicit test writes
The system SHALL separate data preview, confirmed write and readback verification and retain idempotent test identities. Test-write preview SHALL issue a durable token and operation_id bound to the test_write action, saved target, setup revisions, payload digest and expiry. Confirmation MUST validate that binding before the first write; a schema token or plan_id alone MUST NOT authorize test writing. Confirmed writes MUST use the selected persisted target and compare the actual test payload during readback. Results SHALL distinguish written_verified, written_unverified, failed and unknown, with an independent cleanup_status of not_attempted, cleaned, failed or unknown. The UI MUST NOT infer verified writing, cleanup or timestamps without corresponding backend evidence.

#### Scenario: Test preview and action-bound confirmation
- **WHEN** an operator requests a test-write preview for a saved, schema-verified plan
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

### Requirement: Persisted recording membership and explicit plan selection
The system SHALL create and update recording plans only from persisted measurements with verified workspace and equipment ownership. Each member MUST retain its own equipment identity. List, get, create, update and delete SHALL enforce workspace ownership on the backend; missing and foreign resources SHALL return the same safe 404 without reading or mutating foreign state. The UI SHALL require explicit plan selection when multiple plans match the active scope and distinguish query failure from an empty plan list. It MUST NOT invent measurement identifiers, assume point identifiers are measurement identifiers without verification, or silently choose the first device or plan.

#### Scenario: Plan contains measurements from two devices
- **WHEN** persisted measurements from two selected devices are added to one plan
- **THEN** each member retains its original equipment identity and only the selected measurements are included.

#### Scenario: Missing disabled or foreign measurement
- **WHEN** a requested member is missing, disabled, removed or outside the permitted workspace or selected equipment scope
- **THEN** plan creation or update is blocked with a repair action and no placeholder member is inserted.

#### Scenario: Multiple plans or plan-list failure
- **WHEN** multiple plans match or loading plans fails
- **THEN** the UI requests an explicit selection or reports the loading error respectively
- **AND** it neither operates on the first plan nor treats a failed query as permission to create a replacement.

#### Scenario: Foreign plan read update or deletion
- **GIVEN** plan plan-b belongs to workspace workspace-b
- **WHEN** workspace-a requests that plan through get, update or delete
- **THEN** each request returns the same safe 404 as a missing plan and plan-b remains unchanged.

#### Scenario: Selection changes after a test or preview
- **WHEN** the selected plan, equipment scope or connector changes
- **THEN** preview and test results for the previous selection cannot authorize or describe the new selection.

### Requirement: Consistent database setup persistence
The system SHALL save related local connector, row-group, target and workspace-reference changes atomically within their local configuration database and check expected versions before mutation. Failed local saves MUST preserve the prior consistent state and MUST NOT grant readiness. External database operations SHALL retain independent operation outcomes rather than claim cross-database atomic rollback.

#### Scenario: Second local save step fails
- **WHEN** persisting a row group or workspace target reference fails after another local mutation in the same operation
- **THEN** the local operation rolls back and reloading returns the prior consistent configuration.

#### Scenario: Stale save or unresolved external operation
- **WHEN** another session has changed the expected version or an external operation outcome remains unresolved
- **THEN** stale mutation is rejected or readiness remains blocked respectively, and the UI does not claim the current setup is fully applied.
