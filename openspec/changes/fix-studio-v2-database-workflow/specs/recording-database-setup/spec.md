## MODIFIED Requirements

### Requirement: Capability-backed database preparation
The system SHALL expose verified adapter capabilities and actual schema metadata, separately from network and authentication success. It MUST NOT report an unimplemented schema or test-write operation as successful. Unsupported operations SHALL be rejected without executing writes, with a stable safe code and an actionable message. A declared capability alone MUST NOT enable an operation without a verified execution path.

#### Scenario: Connection succeeds but schema inspection is forbidden
- **WHEN** the account can connect but cannot inspect the target
- **THEN** the result states that limitation without returning sample columns or claiming write readiness.

#### Scenario: Unimplemented recording operation is invoked
- **WHEN** schema apply or test write has no connected execution path
- **THEN** the backend returns HTTP 501 with success false and no success payload
- **AND** the UI states that no operation was performed and preserves the draft
- **AND** the separately implemented custom-table preparation path is not disabled by this guard.

#### Scenario: Unknown or unverified database capability
- **WHEN** the selected adapter has not passed the relevant operation verification
- **THEN** that operation is unavailable with a reason rather than optimistically executable.

### Requirement: Revision-bound schema preview and explicit creation
The system MUST bind schema creation to a fresh server-persisted preview token containing workspace, plan and connector identity revisions, target scope and setup revision, and perform backend validation again. Preview MUST NOT mutate the target schema or write test records. Apply MUST resolve the target server-side, atomically claim the operation, retain its outcome and reject stale or unauthorized requests before executing statements. Partial or unknown execution MUST NOT be reported as full success or complete rollback.

#### Scenario: Target changed after preview
- **WHEN** the database or plan changes after schema preview
- **THEN** schema apply is rejected as stale before executing statements.

#### Scenario: Another session changes the setup
- **WHEN** the expected setup revision no longer matches persisted settings at apply time
- **THEN** the backend returns a conflict before executing statements and requires a new preview.

#### Scenario: Unknown expired or foreign preview
- **WHEN** the token is missing, unknown, expired or outside the current workspace
- **THEN** apply does not execute statements and returns a safe non-success result without disclosing foreign resources.

#### Scenario: Repeated apply or lost response
- **WHEN** the same preview is applied concurrently or retried after an uncertain response
- **THEN** at most one execution is claimed and later requests inspect the retained operation result rather than rerunning statements blindly.

#### Scenario: Schema creation partially fails
- **WHEN** the adapter cannot atomically roll back all executed statements
- **THEN** the result identifies confirmed and uncertain work and prevents a full-success indication until reconciliation.

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
The system SHALL separate data preview, confirmed write and readback verification and retain idempotent test identities. Confirmed writes MUST use the selected persisted target and compare the actual test payload during readback. Results SHALL distinguish written_verified, written_unverified, failed and unknown, with an independent cleanup_status of not_attempted, cleaned, failed or unknown. The UI MUST NOT infer verified writing, cleanup or timestamps without corresponding backend evidence.

#### Scenario: Write succeeds without read permission
- **WHEN** a confirmed test write succeeds but the account cannot query the result
- **THEN** the result is written_unverified, not verified.

#### Scenario: Test action retried
- **WHEN** the same confirmed test token is retried after an unknown network outcome
- **THEN** the system checks its durable identity and does not insert a duplicate test record.

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
The system SHALL create and update recording plans only from persisted measurements with verified workspace and equipment ownership. Each member MUST retain its own equipment identity. The UI SHALL require explicit plan selection when multiple plans match the active scope and distinguish query failure from an empty plan list. It MUST NOT invent measurement identifiers, assume point identifiers are measurement identifiers without verification, or silently choose the first device or plan.

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
