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
- **WHEN** confirmation is repeated or the gateway restarts after a possibly executed DDL operation
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
