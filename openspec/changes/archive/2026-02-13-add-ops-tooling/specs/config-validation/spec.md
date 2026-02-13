# Spec: Configuration Validation

## ADDED Requirements

### Requirement: Configuration validation command

The gateway SHALL provide a `validate` subcommand that parses and checks configuration files without starting the server, and SHALL report validation errors with clear messages.

#### Scenario: Valid configuration
- **WHEN** the operator runs `gateway validate -c configs/gateway.yaml` with a valid config directory
- **THEN** the output states "Configuration is valid." and the process exits with code 0

#### Scenario: Invalid protocol
- **WHEN** a device uses an unsupported protocol (e.g. "modbus_bad") in config and `gateway validate` is run
- **THEN** the output contains an error referencing the unsupported protocol and device
- **AND** the process exits with code 1

#### Scenario: Broken datalink reference
- **WHEN** `datalink.yaml` references a non-existent device ID and `gateway validate` is run
- **THEN** the output contains an error referencing the unknown device ID
- **AND** the process exits with code 1
