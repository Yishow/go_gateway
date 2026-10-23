## MODIFIED Requirements

### Requirement: Managed preparation preserves operation and table ownership
Group schema preparation MUST reuse durable operation identity and enforce current scope, revisions, digest, expiry and explicit confirmation. It SHALL create confirmed owned structures, verify an owned compatible no-op, or explicitly confirm the narrowly scoped additive extension defined below. Unknown ownership, incompatible existing tables and unresolved effects MUST NOT cause automatic ALTER, DROP, rename, replay or data movement.

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

### Requirement: Durable schema operation identity and honest recovery
Schema preview SHALL return a server-issued operation_id durably bound to its token. Schema application SHALL serialize execution rights for that confirmation, reject a substituted operation identity and retain the same operation across retries. Repeated requests MUST return or reconcile the same operation rather than blindly executing it again. Results SHALL distinguish succeeded, partial, failed and unknown outcomes and expose safe status lookup within the owning workspace.

#### Scenario: Concurrent duplicate confirmations
- **WHEN** two requests attempt to consume the same schema token concurrently
- **THEN** at most one obtains execution rights and both refer to the same operation outcome: a running duplicate returns 202 and a terminal duplicate returns 200 with the retained result. A different operation occupying the same scope returns 409 without executing another batch.

#### Scenario: Connection lost after target mutation
- **WHEN** acknowledgement is lost after statements have potentially executed
- **THEN** the operation remains unknown until actual schema or durable execution evidence resolves it
- **AND** retry does not blindly repeat the statement batch.

#### Scenario: Adapter cannot roll back all statements
- **WHEN** a later statement fails after earlier statements committed
- **THEN** the result reports partial with verified completed scope and a repair action
- **AND** it does not claim full rollback or complete success.

#### Scenario: Process restarts during execution
- **WHEN** the process restarts with a running or unknown schema operation
- **THEN** recovery uses persisted scope and execution evidence before any further mutation
- **AND** unresolved operations remain blocked from blind replay.

#### Scenario: Unverified adapter selected
- **WHEN** managed schema behavior has not been verified for the selected adapter
- **THEN** that operation remains unavailable without claiming that separate existing custom-table features are unsupported.

#### Scenario: Delayed terminal result cannot release a newer fence
- **GIVEN** one open schema panel is querying operation A while another panel retains a newer unresolved operation B for the same workspace/group
- **WHEN** A later returns failed or succeeded
- **THEN** the client clears a recovery marker only if its retained operation_id equals A; B stays fenced in storage and memory, and another already-open panel cannot preview or apply a new batch.


## ADDED Requirements

### Requirement: Confirmed additive extension of owned managed tables

For verified SQLite and PostgreSQL adapters, managed preview SHALL expose add_columns only for absent nullable, non-primary-key, server-generated v_[0-9a-f]{24} columns in a table whose ownership, identity/metadata, types and primary key match the canonical group. Compatible retained nullable historical managed columns SHALL remain valid and preserved. Preview SHALL remain read-only. Apply SHALL recompute the exact pending statements against the actual destination before claiming the existing operation/token and executing explicitly confirmed DDL. The managed receipt table SHALL retain its full strict schema contract and SHALL NOT use this additive relaxation. Custom or foreign tables, required/primary-key additions, unrelated extra columns, ownership/type/metadata mismatches and stale destination structures SHALL fail closed. Existing records SHALL remain intact; no DROP, rename, destructive conversion or data movement is authorized.

#### Scenario: Add a managed point column
- **GIVEN** a compatible owned table with historical records and a saved canonical group containing a new nullable managed member column
- **WHEN** the operator previews add_columns and explicitly confirms the token-bound operation
- **THEN** only the missing managed column is added, old records remain intact and new records can contain the new member.

#### Scenario: Retain a historical nullable column
- **WHEN** an old nullable managed value column remains after membership removal
- **THEN** readiness accepts the compatible retained column without dropping it or historical values.

#### Scenario: Reject unsafe or changed structure
- **WHEN** the target is foreign/custom, a new column is required or primary-key, or ownership/types/metadata/actual pending statements differ from the confirmed preview
- **THEN** no unconfirmed DDL executes and the operator receives the existing blocked/stale repair response.

#### Scenario: Lost additive apply reply
- **WHEN** the additive apply response is lost after execution
- **THEN** status lookup uses the original operation_id and no unknown operation is resent as a new ALTER batch.

