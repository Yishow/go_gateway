## ADDED Requirements

### Requirement: Invalid edit stops running device

If a running v2 workspace device becomes invalid due to later editing, the system SHALL stop that device immediately.

#### Scenario: Running device becomes invalid

- **GIVEN** a workspace device is currently running
- **WHEN** the operator edits that device into an invalid state
- **THEN** the system stops runtime collection for that device immediately
- **AND** the device is no longer treated as available

### Requirement: Availability state is explicit

The system SHALL expose explicit availability state for v2 workspace devices.

#### Scenario: Unavailable device reports explicit reason

- **WHEN** a running workspace device becomes invalid
- **THEN** the system reports `availability_status = "unavailable"`
- **AND** the system includes an `availability_reason` describing why the device is no longer usable
