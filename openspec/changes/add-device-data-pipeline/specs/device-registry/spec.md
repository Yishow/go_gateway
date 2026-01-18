## ADDED Requirements
### Requirement: Device registry
The system SHALL allow operators to create, update, and delete devices with a selected protocol type and connection profile.

#### Scenario: Create device
- WHEN the operator saves a device with protocol type and connection fields
- THEN the device is stored and available for point assignment

### Requirement: Protocol-specific settings
The system SHALL support protocol-specific connection fields, including TCP host/port, serial port parameters, and device addressing.

#### Scenario: Configure Modbus RTU
- WHEN the operator selects Modbus RTU
- THEN the UI and API accept serial port, baud rate, parity, and slave id fields

### Requirement: Connection test
The system SHALL provide a connection test per device to validate protocol connectivity before activation.

#### Scenario: Test connection success
- WHEN the operator runs a connection test
- THEN the system reports success and stores the last test timestamp

### Requirement: Device lifecycle status
The system SHALL track device status as draft, active, or disabled and expose the current state in UI and API.

#### Scenario: Disable device
- WHEN the operator disables a device
- THEN scheduled collection for that device stops and status is updated
