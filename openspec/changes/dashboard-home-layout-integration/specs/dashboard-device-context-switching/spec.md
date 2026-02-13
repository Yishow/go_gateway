## ADDED Requirements

### Requirement: Device context bar as the primary device state surface
The dashboard SHALL show a persistent device context bar that displays current selected device identity and status, and SHALL provide a single "switch device" entry point.

#### Scenario: Show selected device context
- **WHEN** a device is selected for dashboard operations
- **THEN** the context bar shows device name, id snippet, protocol, and connection/readiness status
- **AND** the context bar shows last switch timestamp

#### Scenario: Status visual levels
- **WHEN** device state changes between active, offline, and read-only
- **THEN** the context bar uses distinct visual levels for each state
- **AND** the current state meaning is visible without opening details

#### Scenario: No device selected state
- **WHEN** no device has been selected yet
- **THEN** the context bar shows a clear empty-state message
- **AND** key write/commit actions remain disabled until a device is selected

### Requirement: Device switch drawer on demand
The dashboard SHALL open a device-selection drawer only when the user explicitly requests switching device.

#### Scenario: Open drawer from context bar
- **WHEN** the user clicks "switch device" in the context bar
- **THEN** the UI opens the device drawer
- **AND** the main workspace remains visible in the background
- **AND** on mobile the drawer is rendered as side-sliding panel (not full-screen sheet)

#### Scenario: Apply selected device from drawer
- **WHEN** the user chooses a device and confirms
- **THEN** the context bar updates to the new device
- **AND** source planner, memory grid, and commit panel reload device-scoped data

#### Scenario: Switching in progress
- **WHEN** a switch request is in-flight
- **THEN** switch controls are temporarily locked
- **AND** a loading state is shown to prevent duplicate submissions

### Requirement: Device switch exception handling
The dashboard SHALL provide deterministic behavior for switch-related exceptions.

#### Scenario: Device offline
- **WHEN** the selected device is offline or not ready
- **THEN** the context bar shows offline/blocked state
- **AND** commit actions are disabled with a visible reason

#### Scenario: Switch request fails
- **WHEN** device switching fails due to API or connectivity error
- **THEN** the UI keeps the previous selected device
- **AND** shows an error toast with retry and details actions

#### Scenario: Unsaved changes before switching
- **WHEN** the user attempts to switch device while unsaved planning changes exist
- **THEN** the switch button shows unsaved-change indicator before user action
- **AND** the UI shows a confirmation dialog with options: save then switch, discard then switch, or cancel
- **AND** default focus in the dialog is set to cancel

#### Scenario: Insufficient permission
- **WHEN** current user lacks write permission for selected device scope
- **THEN** the UI enters read-only mode for affected actions
- **AND** presents a clear permission message in workspace controls
