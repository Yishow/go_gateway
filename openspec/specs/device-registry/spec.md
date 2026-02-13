# device-registry Specification

## Purpose
TBD - created by archiving change add-device-data-pipeline. Update Purpose after archive.
## Requirements
### Requirement: Device registry

The system SHALL allow operators to create, update, and delete devices with a selected protocol type and connection profile.

The system SHALL track device collection statistics including last collection time, collection count, and error count.

#### Scenario: Create device

- **WHEN** the operator saves a device with protocol type and connection fields
- **THEN** the device is stored with default collection statistics (all zeros/null) and available for point assignment

#### Scenario: Device collection statistics

- **WHEN** a device collects data successfully
- **THEN** the device's `last_collected_at` and `collection_count` are updated

#### Scenario: Device collection errors

- **WHEN** a device collection fails
- **THEN** the device's `error_count` is incremented

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

### Requirement: Device readiness check

The system SHALL provide a readiness check API to verify if a device is fully configured and ready to collect data.

#### Scenario: Check ready device

- **WHEN** a client requests readiness status for a device with all required configurations
- **THEN** the API returns `ready: true` with all checks passing

#### Scenario: Check incomplete device

- **WHEN** a client requests readiness status for a device missing required configurations
- **THEN** the API returns `ready: false` with a list of missing configurations and suggestions

#### Scenario: Readiness check includes device status

- **WHEN** a readiness check is performed
- **THEN** the system verifies the device status is `active`

#### Scenario: Readiness check includes points

- **WHEN** a readiness check is performed
- **THEN** the system verifies at least one enabled point exists for the device

#### Scenario: Readiness check includes polling groups

- **WHEN** a readiness check is performed
- **THEN** the system verifies points are assigned to enabled polling groups

#### Scenario: Readiness check includes mappings

- **WHEN** a readiness check is performed
- **THEN** the system verifies points have corresponding enabled mappings

#### Scenario: Readiness check includes scheduler status

- **WHEN** a readiness check is performed
- **THEN** the system verifies the collection scheduler is running

### Requirement: Device collection statistics

The system SHALL track collection statistics for each device including last collection time, collection count, and error count.

#### Scenario: Update collection statistics

- **WHEN** a device successfully collects data
- **THEN** the system updates `last_collected_at` and increments `collection_count`

#### Scenario: Track collection errors

- **WHEN** a device collection fails
- **THEN** the system increments `error_count`

#### Scenario: Query collection statistics

- **WHEN** a client requests device statistics
- **THEN** the API returns current collection statistics

### Requirement: Device onboarding wizard

The system SHALL provide a guided onboarding wizard to help users complete device configuration.

#### Scenario: Auto-open wizard after device creation

- **WHEN** a new device is created successfully
- **THEN** the onboarding wizard automatically opens

#### Scenario: Wizard guides through configuration steps

- **WHEN** a user follows the onboarding wizard
- **THEN** the wizard guides through: device creation, connection test, activation, point creation, polling group assignment, tag creation, and mapping creation

#### Scenario: Wizard shows completion status

- **WHEN** a user progresses through the wizard
- **THEN** each step shows completion status and validation results

#### Scenario: Wizard allows skipping steps

- **WHEN** a user wants to skip a step
- **THEN** the wizard allows skipping with a warning about incomplete configuration

### Requirement: Device status dashboard

The system SHALL provide a dashboard showing device collection status and configuration completeness.

#### Scenario: Display device statistics

- **WHEN** a user views the status dashboard
- **THEN** the dashboard displays total devices, active devices, and devices currently collecting data

#### Scenario: Display device collection status

- **WHEN** a user views a device in the dashboard
- **THEN** the dashboard shows last collection time, collection count, and error count

#### Scenario: Real-time status updates

- **WHEN** device collection status changes
- **THEN** the dashboard updates in real-time using SSE

