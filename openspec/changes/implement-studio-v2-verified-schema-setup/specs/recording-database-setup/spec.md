## MODIFIED Requirements

### Requirement: Capability-backed database preparation
The system SHALL expose verified adapter capabilities and actual schema metadata independently from network and authentication success. Metadata MUST identify the persisted connector revision and selected database, schema and table. Missing, forbidden and failed inspection outcomes MUST remain distinguishable. Example columns MUST NOT be returned as inspected production metadata.

#### Scenario: Connection succeeds but schema inspection is forbidden
- **WHEN** the account connects but cannot inspect the selected target
- **THEN** inspection reports forbidden without sample columns or a write-readiness claim.

#### Scenario: Actual schema differs from a template
- **WHEN** the selected table contains columns different from the application template
- **THEN** inspected columns and types reflect the actual selected table only.

#### Scenario: Empty result is ambiguous
- **WHEN** inspection does not establish whether a table exists
- **THEN** the result is failed or unconfirmed rather than automatically missing
- **AND** automatic table creation is not initiated.

#### Scenario: Target changes during inspection
- **WHEN** a response for an older connector revision or table arrives after selection changes
- **THEN** it cannot replace the current target's metadata or confirm current assignments.

### Requirement: Revision-bound schema preview and explicit creation
The system MUST bind schema creation to a fresh persisted preview token containing the server-resolved workspace, plan or custom-table configuration revision, connector identity revision, table scope and generated-statement digest. It SHALL validate these again before mutation, accept only explicit confirmation, and require the same safeguards at every public schema-mutation entry point. Preview itself MUST have no target side effects.

#### Scenario: Target changed after preview
- **WHEN** the database, connector identity or plan changes after preview
- **THEN** apply is rejected as stale before any statement executes.

#### Scenario: Preview token expired or belongs elsewhere
- **WHEN** a token is expired or belongs to another workspace
- **THEN** apply rejects it without executing statements or revealing foreign metadata.

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

## ADDED Requirements

### Requirement: Durable schema operation identity and honest recovery
Schema application SHALL persist an operation identity and serialize execution rights for the same confirmation. Repeated requests MUST return or reconcile the same operation rather than blindly executing it again. Results SHALL distinguish succeeded, partial, failed and unknown outcomes and expose safe status lookup within the owning workspace.

#### Scenario: Concurrent duplicate confirmations
- **WHEN** two requests attempt to consume the same schema token concurrently
- **THEN** at most one obtains execution rights and both refer to the same operation outcome.

#### Scenario: Connection lost after target mutation
- **WHEN** acknowledgement is lost after statements may have executed
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
