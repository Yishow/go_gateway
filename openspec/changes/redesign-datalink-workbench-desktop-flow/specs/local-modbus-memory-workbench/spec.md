## MODIFIED Requirements

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
The Local Modbus operation surface SHALL include complete local sink functions: server status, bind/start/stop, register mapping edit, conflict scan, batch mapping, test write/read, import/export, sync-from-mappings, and visual register allocation review.

#### Scenario: Full operations available in Output step
- **WHEN** the Local Modbus studio is active in workbench Output
- **THEN** the operator can execute all local sink lifecycle actions from that surface
- **AND** each action shows immediate success/failure feedback without opening a separate tool page

#### Scenario: Visual register review supports mapping decisions
- **WHEN** the operator reviews Local Modbus mappings
- **THEN** the UI shows a register allocation map with base/offset context and conflict visibility
- **AND** the operator can understand register usage before applying changes

## ADDED Requirements

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
