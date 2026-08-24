# protocol-servers Specification

## MODIFIED Requirements

### Requirement: Local Modbus server sink on port 5020

The system SHALL support a local Modbus TCP server sink that mirrors committed tag values into virtual Modbus memory and exposes them to external clients using the persisted Share bind address, port, slave id, and register capacity. Port 5020 SHALL be treated as a persisted or explicitly selected value, not as an unconditional handler fallback.

#### Scenario: Start local Modbus sink from durable settings

- **GIVEN** persisted Share settings have enabled=true, bind address 127.0.0.1, port 15020, valid slave id, and register capacity
- **WHEN** the gateway starts or the operator enables local Modbus sharing
- **THEN** the gateway starts a Modbus TCP server bound to 127.0.0.1:15020
- **AND** the startup status exposes the configured address, port, slave id, capacity, and running state

#### Scenario: Disabled Share has no listener

- **GIVEN** persisted settings.modbus_share.enabled is false
- **WHEN** the gateway starts, hydrates workspace state, or receives a rule-level share_enabled=true
- **THEN** no local Modbus TCP listener is bound
- **AND** the rule-level flag does not create a listener or mutate virtual memory
- **AND** the status reports disabled rather than ready

#### Scenario: Mirror tag value across its complete register span

- **GIVEN** a committed workspace-owned tag mapping has a validated register start, datatype span, and stride
- **WHEN** the mapped tag value is updated in pipeline runtime
- **THEN** the corresponding virtual Modbus registers are updated using two bytes per register and the declared encoding
- **AND** external Modbus clients reading that span receive the mirrored value

#### Scenario: Port or address bind conflict

- **GIVEN** the configured bind address or port is already occupied or unavailable
- **WHEN** local Modbus sharing starts
- **THEN** no partially started listener remains
- **AND** the system reports a typed bind failure with the configured endpoint and an operator action
- **AND** local Modbus sharing remains failed or disabled until the durable setting is corrected and retried

#### Scenario: Device scope isolation

- **GIVEN** pipeline data is processed for multiple devices and workspaces
- **WHEN** each workspace reconciles its persisted Share desired set
- **THEN** each source device mapping remains isolated by workspace and source-rule ownership
- **AND** no cross-device or cross-workspace deletion occurs outside the authorized scope

#### Scenario: Local sink target conflict detection

- **GIVEN** multiple mappings target overlapping register ranges in the same local Share capacity
- **WHEN** the complete desired set is validated
- **THEN** the system flags the range collision before mutation
- **AND** prevents enabling or activating the conflicting set
- **AND** retains the last valid projection until the collision is resolved
