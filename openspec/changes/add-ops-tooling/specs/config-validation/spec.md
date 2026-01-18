# Spec: Configuration Validation

## ADDED Requirements

#### Scenario: Valid Configuration

- Given a valid `configs/` directory.
- When `gateway validate -c configs/gateway.yaml` is run.
- Then the output should say "Configuration is valid." and exit with code 0.

#### Scenario: Invalid Protocol

- Given a `devices.yaml` where a device uses an unsupported protocol (e.g., "modbus_bad").
- When `gateway validate` is run.
- Then the output should contain "Error: unsupported protocol 'modbus_bad' for device '...'" and exit with code 1.

#### Scenario: Broken Datalink Reference

- Given a `datalink.yaml` referencing a non-existent device ID.
- When `gateway validate` is run.
- Then the output should contain "Error: mapping '...' references unknown device ID '...'" and exit with code 1.

## Scope

- Validate `gateway.yaml` (Server settings).
- Validate `devices.yaml` (Protocol/Connection settings).
- Validate `datalink.yaml` (Mapping rules and references).
