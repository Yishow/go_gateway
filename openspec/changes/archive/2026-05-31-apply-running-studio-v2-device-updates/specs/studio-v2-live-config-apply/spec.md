## ADDED Requirements

### Requirement: Running device update application

After a device has already been activated once, later valid autosave updates SHALL apply directly to that running device.

#### Scenario: Valid autosave on a running device applies immediately

- **GIVEN** a workspace device is already running
- **WHEN** the operator makes a valid autosave change to that device's saved configuration
- **THEN** the system applies the new configuration directly to the running device
- **AND** the operator does not need to revisit Step 4 just to make the change take effect

### Requirement: Not-running devices stay dormant

Valid autosave updates SHALL NOT start a device that has never been activated.

#### Scenario: Valid autosave on a not-running device does not start runtime

- **GIVEN** a workspace device has not been activated yet
- **WHEN** the operator makes a valid autosave change to that device's saved configuration
- **THEN** the persisted configuration is updated
- **AND** the device remains not running until the first activation flow is invoked
