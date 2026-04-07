## MODIFIED Requirements

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

## ADDED Requirements

### Requirement: Local Modbus output uses persistent rule-owned mapping state
The system SHALL persist Local Modbus mapping ownership and revision state for rule-driven output review/apply.

#### Scenario: Restart preserves rule-owned Local Modbus state
- **WHEN** the service restarts after Local Modbus candidates or applied mappings were created from a source rule
- **THEN** the system restores the persisted rule-owned Local Modbus state
- **AND** the workbench can continue review/apply without rebuilding mappings from scratch
