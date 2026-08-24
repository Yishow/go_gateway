# modbus-share-lifecycle Specification

## ADDED Requirements

### Requirement: B1 Hydration-gated Share truth

The system SHALL derive Modbus Share summary, candidate state, mappings, activation inputs, and listener status from persisted workspace/settings state after bootstrap hydration. Display defaults SHALL NOT be treated as runtime truth.

#### Scenario: Share is blocked before bootstrap hydration completes

- **GIVEN** the /studio/v2 bootstrap is pending or has failed
- **WHEN** the operator opens the Share output surface or sends a Share mutation request
- **THEN** the UI shows blocked or unknown readiness rather than a ready summary
- **AND** list, delete, upsert, and activation requests are rejected with typed error modbus_share_hydration_required
- **AND** the backend does not read display defaults as mutation input

#### Scenario: Hydrated state becomes the only Share input

- **GIVEN** bootstrap returns hydration state ready, a workspace id, workspace revision, settings revision, and readiness
- **WHEN** the operator requests Share summary or activation
- **THEN** the backend uses the persisted values identified by those revisions
- **AND** the response identifies the same workspace and revisions used for the decision

### Requirement: B2 Global Share setting is an absolute gate

The system SHALL enforce settings.modbus_share.enabled before every Share summary, list, delete, upsert, activation, restore, and listener mutation. A rule-level Share flag SHALL NOT override a disabled global setting.

#### Scenario: Disabled global Share hides the surface and prevents mutation

- **GIVEN** persisted settings.modbus_share.enabled is false
- **WHEN** the operator loads /studio/v2 or calls a Share list, delete, or upsert endpoint
- **THEN** the Share summary is omitted or marked unavailable
- **AND** the API returns modbus_share_disabled for mutation/list operations
- **AND** no durable mapping or in-memory runtime mapping is created, changed, or deleted

#### Scenario: A rule flag cannot enable a disabled global Share

- **GIVEN** a persisted source rule has share_enabled=true while the global setting is disabled
- **WHEN** activation or process restart runs reconciliation
- **THEN** the rule is excluded from the desired Share projection
- **AND** no listener is started and no register mapping is restored

### Requirement: B3 Durable workspace ownership is required

The system SHALL authorize Share mapping operations only when a persisted workspace mapping, source-rule revision, and persisted tag relationship prove ownership for the requested workspace. A non-empty or browser-only tag_id SHALL NOT constitute ownership.

#### Scenario: A persisted relationship authorizes an owned mapping

- **GIVEN** a mapping, source rule revision, and tag are persisted and all three relationships identify workspace ws-a
- **WHEN** ws-a reconciles its complete desired mapping set
- **THEN** the backend can update or remove that mapping within ws-a scope
- **AND** the result records the workspace and source-rule revision used as ownership evidence

#### Scenario: A non-empty tag id cannot authorize deletion

- **GIVEN** a runtime mapping has tag_id=tag-a but no persisted workspace mapping or source-rule/tag relationship for ws-a
- **WHEN** ws-a sends a delete or reconcile request that omits the mapping
- **THEN** the backend rejects deletion with modbus_share_workspace_scope
- **AND** the existing mapping is retained and reported for recovery or operator review

#### Scenario: Cross-workspace operations are rejected

- **GIVEN** a mapping is durably owned by ws-b
- **WHEN** a caller scoped to ws-a attempts list, upsert, delete, or activation against that mapping
- **THEN** the backend rejects the operation with modbus_share_workspace_scope
- **AND** no mapping outside ws-a changes

### Requirement: B4 Register units, datatype spans, and capacity are backend contracts

The system SHALL define each holding register as two bytes, convert human register 40001 to backend zero-based register 0, and validate datatype span, stride, capacity, and range collision on the backend. Frontend checks SHALL remain advisory.

#### Scenario: Human and zero-based register values agree

- **GIVEN** a desired mapping has human share_start_register=40001 and datatype int16
- **WHEN** backend validation succeeds
- **THEN** the projection uses zero-based register 0
- **AND** the response identifies one register and two bytes for the value

#### Scenario: Multi-register values reserve their complete span

- **GIVEN** the backend uses two bytes per register
- **WHEN** an operator configures int32, float32, int64, and float64 mappings
- **THEN** the backend reserves spans of 2, 2, 4, and 4 registers respectively
- **AND** each mapping requires stride_registers greater than or equal to its datatype span

#### Scenario: Capacity and range collision fail before mutation

- **GIVEN** capacity is 8 registers, one mapping occupies zero-based range [0,2), and another requested range overlaps it or exceeds [0,8)
- **WHEN** the desired set is validated
- **THEN** validation fails with modbus_share_range_collision or modbus_share_capacity_exceeded
- **AND** no durable or runtime mapping mutation occurs

### Requirement: B5 Reconcile is atomic, serialized, idempotent, and revision-checked

The system SHALL prevalidate the complete desired mapping set before mutation, serialize reconciliation per workspace, and use workspace revision CAS to prevent concurrent stale writes. A failed reconcile SHALL retain the prior projection or report an explicit recoverable dirty or unknown state.

#### Scenario: A valid desired set is applied as one revision

- **GIVEN** all desired mappings pass hydration, global, ownership, span, capacity, collision, and source-rule revision validation
- **WHEN** a caller submits the expected workspace revision
- **THEN** the backend commits the durable and runtime projection as one reconciled outcome
- **AND** the response contains the new workspace revision and complete changed-set summary

#### Scenario: Partial mutation failure preserves the old projection

- **GIVEN** the old projection is valid and a runtime or repository mutation fails while applying a new desired set
- **WHEN** rollback completes
- **THEN** the old durable and runtime projection remains active
- **AND** the response reports modbus_share_reconcile_failed with retryable recovery information
- **AND** the API does not report success

#### Scenario: Rollback uncertainty is explicit

- **GIVEN** applying or rolling back a desired set cannot establish the final durable/runtime state
- **WHEN** reconcile returns
- **THEN** the workspace is marked dirty_unknown
- **AND** the response identifies dirty scope, observed revision, and a recovery action
- **AND** subsequent activation is blocked until a recovery reconcile establishes alignment

#### Scenario: Double submit and stale revision are safe

- **GIVEN** two requests submit the same desired set and expected revision concurrently
- **WHEN** both requests reach the serialized reconcile service
- **THEN** one commit is applied and a repeated equivalent request returns the same idempotent outcome
- **AND** a request using an obsolete revision receives modbus_share_revision_conflict without overwriting the committed state

### Requirement: B6 Activation waits for a durable autosave barrier

The system SHALL activate Share only after all required durable saves complete successfully for the same workspace and settings revisions. Save errors, pending saves, stale save responses, and incomplete readiness SHALL block activation.

#### Scenario: Successful saves unlock activation

- **GIVEN** global Share settings, source-rule Share fields, mappings, and required candidate decisions have each returned durable success for workspace revision r
- **WHEN** the operator clicks activation
- **THEN** the frontend awaits those save results before sending activation
- **AND** the backend accepts activation only when the request carries revision r and the matching readiness token

#### Scenario: Save failure blocks activation

- **GIVEN** one required autosave returns an error or remains pending
- **WHEN** the operator requests activation
- **THEN** no activation or Share projection mutation is attempted
- **AND** the UI and API report modbus_share_save_incomplete with the failed save scope

#### Scenario: A stale activation request is rejected

- **GIVEN** the workspace changed after the save barrier produced revision r
- **WHEN** an activation request carries revision r
- **THEN** the backend returns modbus_share_revision_conflict
- **AND** the operator must hydrate the newer persisted state before retrying

### Requirement: B7 Existing candidate/apply/restore seam is authoritative

The system SHALL build Share desired mappings through the existing source-rule candidate, validation, apply, and restore seam. Backend activation SHALL be the projection authority; a direct runtime mapping API SHALL NOT bypass the Share projection contract.

#### Scenario: Candidate apply produces the only formal projection input

- **GIVEN** a source-rule candidate snapshot contains a valid persisted Share candidate and source-rule revision
- **WHEN** backend activation reconciles the workspace
- **THEN** the candidate is included in the complete desired mapping set
- **AND** the runtime mapping is created only by the Share reconcile service

#### Scenario: Direct mapping mutation is rejected

- **GIVEN** a caller sends a direct tag/register runtime mapping request without a workspace-scoped Share candidate and expected revision
- **WHEN** the API handles the request
- **THEN** the API rejects it with modbus_share_projection_required
- **AND** the formal Share desired state and runtime projection remain unchanged

#### Scenario: Restart restores persisted desired mappings

- **GIVEN** an enabled global Share has persisted source-rule candidates and workspace mappings
- **WHEN** the gateway process restarts
- **THEN** the backend reconstructs the same desired mappings and register spans through the restore/reconcile seam
- **AND** the restored projection is independent of the previous browser session

### Requirement: B8 Global settings control the listener lifecycle

The system SHALL use durable global settings as the source of truth for Share enabled, bind address, port, slave id, and register capacity. The listener SHALL be disabled by default, and lifecycle failures SHALL expose actionable typed diagnostics.

#### Scenario: Disabled-by-default startup has no listener

- **GIVEN** persisted Share settings have enabled=false
- **WHEN** the gateway starts or settings hydration completes
- **THEN** no Modbus listener is bound
- **AND** the status reports disabled
- **AND** Share activation does not mutate runtime mappings

#### Scenario: Enabled settings start the configured listener

- **GIVEN** persisted settings have enabled=true, an explicit bind address, port, valid slave id, and capacity
- **WHEN** startup or an enable action runs
- **THEN** the listener binds exactly that address and port with that slave id
- **AND** status reports starting followed by running with the configured values

#### Scenario: Bind failure is actionable and non-partial

- **GIVEN** the configured address or port cannot be bound
- **WHEN** the listener starts or the setting is applied
- **THEN** no new listener remains bound
- **AND** status reports failed with modbus_share_listener_bind_failed, retryability, and an operator action
- **AND** persisted settings remain available for correction without silently falling back to port 5020

#### Scenario: Disable stops the listener

- **GIVEN** the listener is running
- **WHEN** the operator persists enabled=false and the save succeeds
- **THEN** the listener stops before the lifecycle action completes
- **AND** subsequent status is disabled with no Share mutation caused by the disabled rule flags

### Requirement: B9 Mapping changes invalidate obsolete datatype spans

The system SHALL clear obsolete register spans or explicitly mark them invalidated when a mapping is deleted, moved, or retyped. A mapping identity SHALL include workspace, source rule, tag, datatype, zero-based start, span, stride, and source revision.

#### Scenario: Delete clears the complete old span

- **GIVEN** an owned int64 mapping occupies four registers
- **WHEN** the operator deletes it and reconcile succeeds
- **THEN** all four old registers are cleared or recorded as invalidated before the new projection is committed
- **AND** the response lists the removed span

#### Scenario: Move or retype does not leave stale words

- **GIVEN** an owned float32 mapping occupies two registers at zero-based start 0
- **WHEN** the operator moves it or changes it to int64
- **THEN** the old two-register span is removed or invalidated and the new four-register span is validated before commit
- **AND** a collision or capacity error preserves the old projection instead of exposing a partial span

#### Scenario: Unknown ownership blocks silent cleanup

- **GIVEN** an obsolete runtime span lacks persisted ownership evidence
- **WHEN** the desired set no longer contains that span
- **THEN** the backend does not silently clear it
- **AND** it reports invalidated_unknown with modbus_share_workspace_scope and a recovery action
- **AND** the workspace is not reported as fully ready

### Requirement: B10 End-to-end acceptance and field boundary

The system SHALL provide an acceptance path that verifies configure, durable save, activation, Modbus client read, process restart, mapping restore, and post-restart read. Disabled mode, collision, capacity, partial failure, ownership, and double-submit tests SHALL be included. Real PLC/SCADA hardware acceptance SHALL remain a separate gate.

#### Scenario: Configure through restart and read

- **GIVEN** an operator configures explicit Share settings and valid source-rule/tag mappings
- **WHEN** the acceptance flow saves settings, activates the workspace, reads the mapped holding registers with a Modbus TCP client, restarts the EXE, and reads again
- **THEN** the first and second reads return the expected encoded values
- **AND** the restored listener address/port/slave id, mapping spans, and workspace revision are observable in backend status

#### Scenario: Disabled mode proves no listener and no mutation

- **GIVEN** global Share is disabled and a rule-level Share flag is enabled
- **WHEN** the acceptance flow starts the EXE, requests summary/list/upsert/activation, and probes the configured endpoint
- **THEN** the endpoint has no listener
- **AND** the API returns modbus_share_disabled
- **AND** durable mappings and runtime memory are unchanged

#### Scenario: Required negative tests are executable

- **GIVEN** fixtures for range collision, capacity overflow, partial repository/runtime failure, missing ownership, and two concurrent revisions
- **WHEN** the focused backend/API test suite runs
- **THEN** each fixture asserts its typed error and verifies no forbidden partial mutation
- **AND** the double-submit fixture asserts idempotent convergence

#### Scenario: Field hardware acceptance is explicit

- **GIVEN** repo HTTP, runtime, EXE, and real Modbus client checks pass
- **WHEN** a deployment team performs acceptance with the target SCADA/PLC network and operator
- **THEN** the result records endpoint, bind policy, EXE identity/hash, register read evidence, restart evidence, and operator sign-off
- **AND** repo-only validation does not claim the field gate is complete
