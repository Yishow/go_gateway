## ADDED Requirements

### Requirement: First activation targets all eligible devices

The Step 4 final action SHALL activate all devices in the singleton v2 workspace that are currently valid, available, and not yet running.

#### Scenario: Activate all eligible devices

- **GIVEN** the singleton workspace contains multiple devices
- **AND** some devices are valid, available, and not yet running
- **WHEN** the operator invokes the final Step 4 action
- **THEN** the system activates every eligible device in that set
- **AND** already running devices are not activated a second time

### Requirement: Activation returns per-device results

The Step 4 final action SHALL report activation success and failure per device.

#### Scenario: Partial activation success

- **GIVEN** two eligible devices are activated and one succeeds while one fails
- **WHEN** the Step 4 action completes
- **THEN** the result surface reports one success and one failure separately
- **AND** the operator can tell which device failed and why

### Requirement: Navigation remains available after partial failure

The operator SHALL still be allowed to proceed to runtime after partial Step 4 activation failure.

#### Scenario: Proceed after partial activation failure

- **GIVEN** at least one device activated successfully and at least one device failed
- **WHEN** the operator chooses to continue
- **THEN** the system allows navigation to `/studio/runtime`
- **AND** the failed devices remain reported in Step 4
