## MODIFIED Requirements

### Requirement: Local Modbus server sink on port 5020

The system SHALL support a local Modbus TCP server sink that mirrors committed tag values into virtual Modbus memory and exposes them to external clients on port `5020`.

#### Scenario: Start local Modbus sink server
- **WHEN** the operator enables local Modbus sharing mode
- **THEN** the gateway starts a local Modbus TCP server bound to port `5020`
- **AND** the startup status is visible in the operator workspace

#### Scenario: Mirror tag value to virtual register
- **WHEN** a committed tag value is updated in pipeline runtime
- **THEN** the mapped virtual Modbus register is updated with the latest value
- **AND** external Modbus clients reading that register receive the mirrored value

#### Scenario: Port bind conflict
- **WHEN** port `5020` is already occupied
- **THEN** the system reports a bind failure with actionable guidance
- **AND** local Modbus sharing mode remains disabled until resolved

#### Scenario: Device scope isolation
- **WHEN** pipeline data is processed for multiple devices
- **THEN** each source device mapping remains isolated by device scope
- **AND** no cross-device conflict is raised outside local sink target checks

#### Scenario: Local sink target conflict detection
- **WHEN** multiple mappings (from same or different devices) target the same local `5020` register
- **THEN** the system flags the mapping conflict in local sink scope
- **AND** prevents enabling conflicting writes until resolved
