# local-modbus-memory-workbench Specification

## ADDED Requirements

### Requirement: Hydrated workspace-scoped Local Modbus operations

The Local Modbus operation surface SHALL remain available inside the workbench Output step, but every summary, candidate review, apply, delete, and register review action SHALL use a ready persisted workspace scope. Browser-only selections SHALL NOT grant runtime ownership.

#### Scenario: Output surface waits for hydration

- **GIVEN** the /studio/v2 workspace bootstrap is pending or failed
- **WHEN** the operator opens Local Modbus in the Output step
- **THEN** the surface shows blocked or unknown readiness
- **AND** candidate review and mapping mutation controls remain unavailable until persisted workspace hydration succeeds

#### Scenario: Workspace scope is sent with each operation

- **GIVEN** workspace ws-a is hydrated with a persisted revision
- **WHEN** the operator reviews or applies Local Modbus candidates
- **THEN** each backend request carries the workspace identity and expected revision
- **AND** the response identifies the same scope and revision used by the backend

### Requirement: Backend-authoritative Local Modbus projection

The workbench SHALL obtain Local Modbus mapping truth from the source-rule candidate/apply/restore projection and SHALL NOT treat direct runtime mappings or frontend display state as the durable desired set.

#### Scenario: Candidate review reflects persisted source-rule state

- **GIVEN** a source rule has a persisted Local Modbus candidate snapshot
- **WHEN** the operator opens register allocation review
- **THEN** the UI renders the candidate status, datatype, span, ownership, and backend validation outcome from that snapshot
- **AND** an unpersisted frontend row is not shown as applied

#### Scenario: Direct runtime mapping is not a formal apply path

- **GIVEN** a caller attempts to bind a tag directly to a runtime register without a workspace-scoped candidate revision
- **WHEN** the request is handled
- **THEN** the backend rejects the operation with a projection-required error
- **AND** the workbench keeps the last durable projection unchanged

### Requirement: Span-aware review and recovery feedback

The register allocation review SHALL display human register, zero-based register, datatype span, stride, capacity, collision state, invalidated state, and actionable recovery diagnostics for the complete desired set.

#### Scenario: Review exposes multi-register allocation

- **GIVEN** a candidate contains an int64 or float32 datatype
- **WHEN** the operator inspects the allocation map
- **THEN** the UI displays the complete four-register or two-register span
- **AND** it displays the configured stride and capacity used by backend validation

#### Scenario: Reconcile failure remains visible

- **GIVEN** backend prevalidation, commit, or rollback returns a typed failure
- **WHEN** the operation completes
- **THEN** the surface shows failed, dirty, or unknown state with the backend error code and retry action
- **AND** it does not show a success or ready state for a partially applied mapping
