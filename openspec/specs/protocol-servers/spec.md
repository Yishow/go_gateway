# protocol-servers Specification

## Purpose
TBD - created by archiving change p3-enable-virtual-device-mode. Update Purpose after archive.
## Requirements
### Requirement: Virtual server hosting

The system SHALL support hosting one or more protocol servers (listeners) to accept incoming connections from external clients (SCADA/HMI).

#### Scenario: Host Modbus TCP server

- GIVEN the gateway is configured to start a Modbus TCP server on port 10502
- WHEN the gateway starts
- THEN it listens on port 10502 and accepts Modbus queries

### Requirement: Virtual memory mapping

The system SHALL provide a strictly addressed virtual memory space that serves as the data source for hosted servers.

#### Scenario: Read from virtual memory

- GIVEN a value 12345 is stored at virtual address D100
- WHEN an external Modbus client reads Holding Register 100
- THEN the server responds with value 12345

### Requirement: Data bridging

The system SHALL support mapping collected data from real devices into the virtual memory space.

#### Scenario: Bridge data to virtual server

- GIVEN a real device tag `PLC1.D200` is collected
- AND a mapping exists from `PLC1.D200` to Virtual Address `V.D100`
- WHEN `PLC1.D200` changes value
- THEN `V.D100` is automatically updated

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

