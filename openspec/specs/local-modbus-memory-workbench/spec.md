# local-modbus-memory-workbench Specification

## Purpose
TBD - created by archiving change dashboard-home-layout-integration. Update Purpose after archive.
## Requirements
### Requirement: Dedicated local Modbus memory workbench page

The system SHALL provide Local Modbus operations as a first-class surface inside the workbench `Output` step, while `/datalink/local-modbus` may remain as a compatibility or fallback entry.

#### Scenario: Workbench output hosts primary Local Modbus operations
- **WHEN** a user enters Step 4 of `/datalink/workbench` and selects `Local Modbus`
- **THEN** the UI exposes the complete Local Modbus operation surface inside the workbench
- **AND** the operator does not need to leave the workbench route to complete the primary local sink workflow

#### Scenario: Compatibility entry remains discoverable
- **WHEN** a user accesses `/datalink/local-modbus`
- **THEN** the system provides a compatibility path to the same Local Modbus operation model or a clear return into workbench Output
- **AND** preserves operator context as much as possible

### Requirement: Complete local sink operation set
The Local Modbus operation surface SHALL include rule-scoped candidate review/apply, server status, verification-first apply checks, register conflict review, batch apply, test write/read, and visual register allocation review.

#### Scenario: Full rule-scoped operations available in Output step
- **WHEN** the Local Modbus studio is active in the workbench Output step
- **THEN** the operator can review rule-derived Local Modbus candidates, verify them, and apply them from that surface
- **AND** each action shows immediate success, conflict, or failure feedback without opening a separate manual-first tool page

#### Scenario: Visual register review supports candidate decisions
- **WHEN** the operator reviews Local Modbus candidates
- **THEN** the UI shows a register allocation map with base/offset context and conflict visibility
- **AND** the operator can understand register usage before apply

### Requirement: Local 5020 conflict governance
The workbench SHALL enforce duplicate/conflict detection within local `5020` memory grid mapping scope for rule-derived Local Modbus candidates.

#### Scenario: Detect duplicate register mapping
- **WHEN** two or more rule-derived Local Modbus candidates target overlapping local register ranges in `5020` sink scope
- **THEN** the system flags those candidates as `blocked_conflict` before apply
- **AND** blocks activation only for the conflicting candidates until the overlap is resolved

#### Scenario: Non-conflicting candidates remain applyable
- **WHEN** Local Modbus candidates share the same review session but only some overlap in register space
- **THEN** the workbench keeps non-conflicting candidates apply-eligible
- **AND** does not block the entire candidate set because of one conflict group

### Requirement: Verification-first write flow
The workbench SHALL require verification before enabling Local Modbus apply for the selected candidate set.

#### Scenario: Verify then apply
- **WHEN** an operator prepares a batch apply on Local Modbus candidates
- **THEN** the UI runs preflight checks for bind status, conflicts, and write readiness
- **AND** enables apply only after all selected candidates pass verification

### Requirement: Auto-map and dry-run support
The Local Modbus operation surface SHALL support operator-selectable auto-map strategies and a dry-run style write verification flow.

#### Scenario: Choose auto-map strategy
- **WHEN** the operator runs auto-map for Local Modbus candidates
- **THEN** the UI allows selection of an auto-map strategy
- **AND** applies the chosen strategy consistently to the current candidate set

#### Scenario: Validate write without destructive commit
- **WHEN** the operator executes a test write or dry-run validation
- **THEN** the UI shows whether the mapping/write path is valid before final apply
- **AND** surfaces any blocking error inline

### Requirement: Health summary and sync result visibility
The Local Modbus operation surface SHALL display server health and the latest sync result summary.

#### Scenario: Show server health summary
- **WHEN** the Local Modbus studio is active
- **THEN** the UI shows server health information and current operational status
- **AND** the operator can tell whether the local sink is available for output

#### Scenario: Show latest sync outcome
- **WHEN** a sync or apply action completes
- **THEN** the UI shows the latest success or failure summary
- **AND** keeps enough detail for the operator to understand what needs follow-up

### Requirement: Local Modbus selection state is isolated from other output targets
The workbench SHALL isolate Local Modbus selection state from other Output target states.

#### Scenario: Switching targets does not reuse stale selection
- **WHEN** an operator switches from Database output back to Local Modbus
- **THEN** the Local Modbus studio SHALL use its own selection context
- **AND** SHALL NOT inherit stale selected tag or binding state from the Database target

### Requirement: Local Modbus binding uses one authoritative state model
The workbench SHALL derive selected tag, selected register target, bind readiness, and unbind readiness from one authoritative state model.

#### Scenario: Register target and selected tag stay consistent
- **WHEN** an operator changes the selected tag or selected register target
- **THEN** the visible binding controls and values SHALL update from one shared state model
- **AND** manual field input SHALL NOT silently drift away from the actual pending binding target

#### Scenario: Bind and unbind are explicit operations with inline feedback
- **WHEN** an operator binds or unbinds a Local Modbus register
- **THEN** the workbench SHALL perform one explicit operation
- **AND** SHALL show immediate conflict, success, or failure feedback inline in the same flow

### Requirement: Local Modbus output uses persistent rule-owned mapping state
The system SHALL persist Local Modbus mapping ownership and revision state for rule-driven output review/apply.

#### Scenario: Restart preserves rule-owned Local Modbus state
- **WHEN** the service restarts after Local Modbus candidates or applied mappings were created from a source rule
- **THEN** the system restores the persisted rule-owned Local Modbus state
- **AND** the workbench can continue review/apply without rebuilding mappings from scratch

