## ADDED Requirements

### Requirement: Canonical group managed preparation
A saved canonical WriteGroup SHALL be sufficient input for managed schema preparation without a RecordingPlan, fabricated measurement or manual target list. For verified SQLite and PostgreSQL adapters, the operator SHALL be able to prepare an initially empty recording target through the product UI using the existing revision-bound preview and explicitly confirmed operation safeguards.

#### Scenario: Empty managed destination
- **GIVEN** saved device/point/Tag/group identities and no recording tables in the selected destination
- **WHEN** the operator previews and explicitly confirms managed preparation through Studio V2
- **THEN** the required owned data and receipt structures are created and independently inspected without user-written SQL or a separate advanced plan.

#### Scenario: SQLite preview remains read-only
- **WHEN** a proposed managed SQLite destination file does not yet exist
- **THEN** preview explains the proposed backend-owned destination without creating the file; only explicit confirmation can create it.

#### Scenario: Internal database is not a recording destination
- **WHEN** a requested destination resolves to the gateway configuration or journal database
- **THEN** managed creation is refused before mutation and a distinct application-owned destination is offered.

### Requirement: Managed preparation preserves operation and table ownership
Group schema preparation MUST reuse durable operation identity and enforce current scope, revisions, digest, expiry and explicit confirmation. It SHALL create only the confirmed owned structures or verify an owned compatible no-op. Unknown ownership, incompatible existing tables and unresolved effects MUST NOT cause automatic ALTER, DROP, rename, replay or data movement.

#### Scenario: Duplicate or lost confirmation response
- **WHEN** confirmation is repeated or the gateway restarts after a DDL operation whose execution is unresolved
- **THEN** the same operation is returned or reconciled from actual evidence; unresolved results stay unknown and no new statement batch is blindly executed.

#### Scenario: Scope changes before confirmation
- **WHEN** group/source/connector revision, destination scope or preview digest no longer matches
- **THEN** confirmation is rejected before any target mutation and a fresh preview is required.

#### Scenario: Existing table conflict or insufficient privilege
- **WHEN** the intended name is unowned/incompatible or the account cannot inspect or create the required structures
- **THEN** the product reports the specific safe blocked state, preserves existing objects and offers a new-table or advanced-table repair path rather than assuming absence.

#### Scenario: Existing custom table remains unchanged
- **WHEN** an operator retains a custom table whose all-good complete-row contract does not require external per-member metadata
- **THEN** this change neither adds metadata/receipt columns nor converts that target into managed storage; its existing capability and uncertainty rules remain in force.

#### Scenario: Metadata belongs to the saved group table
- **GIVEN** the workspace connector defaults to table A and saved group B targets table B
- **WHEN** the operator loads B with its current group and connector revisions
- **THEN** metadata is inspected for persisted table B; unknown, deleted, foreign or stale group scope is rejected before inspection and no client-supplied table replaces the saved target.

## MODIFIED Requirements

### Requirement: Revision-bound schema preview and explicit creation
A canonical WriteGroup SHALL serve as the custom-table configuration: PlanID maps to group.ID and PlanRevision maps to group.Revision in the existing preview token. The group adapter MUST also bind and validate destination schema revision/digest, member/source revisions and the exact layout digest; an empty pre-Apply schema revision is not proof of confirmation.

The system MUST bind schema creation to a fresh persisted preview token containing the server-resolved workspace, plan or custom-table configuration revision, connector identity revision, table scope and generated-statement digest. It SHALL validate these again before mutation, accept only explicit confirmation, and require the same safeguards at every public schema-mutation entry point. Preview itself MUST have no target side effects. The backend SHALL reject missing or unknown confirmation, a client dialect that differs from the server-resolved connector dialect, and stale setup revisions before mutation; legacy tokens lacking required persisted scope, revisions, digest or expiry SHALL require a new preview.

#### Scenario: Target changed after preview
- **WHEN** the database, connector identity or plan changes after preview
- **THEN** apply is rejected as stale before any statement executes.

#### Scenario: Preview token expired or belongs elsewhere
- **WHEN** a token is expired or belongs to another workspace
- **THEN** apply rejects it without executing statements or revealing foreign metadata.

#### Scenario: Missing unknown or legacy confirmation
- **WHEN** confirmation is missing, unknown or lacks the required persisted protection fields
- **THEN** apply rejects it without mutation and requires a new preview rather than guessing valid revisions.

#### Scenario: Another session changes the setup
- **WHEN** the expected workspace or custom-table setup revision differs from persisted state at apply time
- **THEN** the backend returns 409 before executing statements and requires a new preview.

#### Scenario: Client alters statements
- **WHEN** a client submits statements different from the stored preview
- **THEN** no client-supplied statement executes and the request is rejected.

#### Scenario: Compatible existing schema
- **WHEN** inspection confirms the requested structure already exists unchanged
- **THEN** preview explains that no change is needed
- **AND** any successful no-op apply is supported by actual compatibility verification.

#### Scenario: Existing incompatible table
- **WHEN** a requested managed table name collides with an incompatible existing table
- **THEN** preparation is blocked without altering, renaming or deleting that table or its data.

#### Scenario: Legacy endpoint bypass attempt
- **WHEN** a client requests non-preview generation without the required valid confirmation
- **THEN** the operation is rejected even when invoked through a legacy or generic schema endpoint.

#### Scenario: Activation without explicit schema confirmation
- **WHEN** workspace activation finds missing required database structure
- **THEN** activation reports the required preparation step rather than silently creating structure
- **AND** a Local Modbus-only setup does not require database preparation.

#### Scenario: Canonical group or layout changes after preview
- **WHEN** the saved group revision, source revision or destination schema revision/digest changes after preview
- **THEN** confirmation rejects the stale canonical configuration before any target mutation and requires a new preview; no fabricated RecordingPlan or second token family is introduced.
