## ADDED Requirements

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
