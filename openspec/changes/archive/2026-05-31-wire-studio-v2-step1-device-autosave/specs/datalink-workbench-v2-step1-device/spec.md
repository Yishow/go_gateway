## ADDED Requirements

### Requirement: Valid-only device autosave

The Step 1 device workspace SHALL autosave only valid device content to the backend workspace.

#### Scenario: Valid device change saves immediately

- **GIVEN** a device row already belongs to the v2 workspace
- **WHEN** the operator edits the device into a valid form state
- **THEN** the system sends an immediate save request for that device
- **AND** the saved backend version reflects the latest valid values

#### Scenario: Invalid device change stays local

- **GIVEN** a device row already has a last successful backend version
- **WHEN** the operator edits the device into an invalid form state
- **THEN** the system does not overwrite the backend version
- **AND** the UI keeps the invalid local values visible with an unsaved marker

### Requirement: Per-device save isolation

The Step 1 device workspace SHALL isolate save success and failure per device tab.

#### Scenario: One invalid device does not block another valid device

- **GIVEN** the workspace shows multiple devices
- **WHEN** one device remains invalid while another device is edited into a valid state
- **THEN** the valid device still saves successfully
- **AND** the invalid device remains local with its own error state

### Requirement: Persisted device deletion cascade

Deleting a persisted Step 1 device from the v2 workspace SHALL remove that device and its v2-owned related data from the backend.

#### Scenario: Delete persisted device

- **GIVEN** a persisted device belongs to the singleton v2 workspace
- **WHEN** the operator deletes that device from Step 1
- **THEN** the backend removes the device from the workspace
- **AND** the device no longer appears after a page reload
