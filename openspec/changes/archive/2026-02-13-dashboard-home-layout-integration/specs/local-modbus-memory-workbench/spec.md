## ADDED Requirements

### Requirement: Dedicated local Modbus memory workbench page
The system SHALL provide a dedicated page at `/datalink/local-modbus` for local Modbus `5020` sink operations with full workflow coverage.

#### Scenario: Open dedicated page
- **WHEN** the operator enters `/datalink/local-modbus`
- **THEN** the UI shows a full-page workbench for local memory sink management
- **AND** the operator does not need to return to dashboard cards for core sink actions

### Requirement: Complete local sink operation set
The dedicated page SHALL include complete local sink functions: server status, bind/start/stop, register mapping edit, conflict scan, batch mapping, test write/read, import/export, and sync-from-mappings.

#### Scenario: Full operations available
- **WHEN** the workbench page is loaded
- **THEN** the operator can execute all local sink lifecycle actions from one page
- **AND** each action shows immediate success/failure feedback

### Requirement: Local 5020 conflict governance
The workbench SHALL enforce duplicate/conflict detection within local `5020` memory grid mapping scope.

#### Scenario: Detect duplicate register mapping
- **WHEN** two mappings target the same local register in `5020` sink scope
- **THEN** the page flags the conflict before apply
- **AND** blocks write activation until conflict is resolved

#### Scenario: Device isolation with local sink exception
- **WHEN** mappings come from different source devices
- **THEN** normal source-device scopes remain isolated
- **AND** conflict checks are enforced only for shared local `5020` memory grid target space

### Requirement: Verification-first write flow
The workbench SHALL require verification before enabling bulk apply to local `5020`.

#### Scenario: Verify then apply
- **WHEN** operator prepares a batch write/update on local sink mappings
- **THEN** the UI runs preflight checks (bind status, conflict, permission)
- **AND** enables apply only after all checks pass
