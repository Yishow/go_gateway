## ADDED Requirements

### Requirement: Runtime page is workspace-scoped first

The runtime page SHALL load the singleton v2 workspace device set before choosing the current device view.

#### Scenario: Runtime opens from workspace context

- **WHEN** an operator opens `/studio/runtime`
- **THEN** the system loads the singleton workspace runtime context
- **AND** the page is not limited to a single-device-only bootstrap contract

### Requirement: Default device follows v2 order

The runtime page SHALL choose the first available device from the v2 workspace ordering as the default observed device.

#### Scenario: First available device becomes default

- **GIVEN** the workspace order is `device-A`, `device-B`, `device-C`
- **AND** `device-A` is available
- **WHEN** the runtime page opens
- **THEN** the default observed device is `device-A`

#### Scenario: Unavailable first device is skipped for default observation

- **GIVEN** the workspace order is `device-A`, `device-B`
- **AND** `device-A` is unavailable while `device-B` is available
- **WHEN** the runtime page opens
- **THEN** the default observed device is `device-B`
- **AND** `device-A` still remains visible in the device list as unavailable

### Requirement: Unavailable devices remain visible

The runtime page SHALL keep unavailable workspace devices visible in the device list.

#### Scenario: Unavailable device stays listed

- **WHEN** a workspace device is unavailable
- **THEN** it still appears in the runtime device list
- **AND** the UI marks it as unavailable with a reason

##### Example: failed device remains visible beside a healthy device

- **GIVEN** the workspace device list contains `device-A` as available and `device-B` as unavailable with reason `invalid Step 1 configuration`
- **WHEN** the runtime page renders the device list
- **THEN** both `device-A` and `device-B` remain visible
- **AND** `device-B` is labeled unavailable instead of disappearing

### Requirement: Empty runtime state stays on page

If no workspace device is currently available, the runtime page SHALL remain on `/studio/runtime` and show an empty-state message.

#### Scenario: No available devices

- **WHEN** the workspace runtime context contains zero available devices
- **THEN** the runtime page stays open
- **AND** the page shows an empty-state message with guidance to return to `/studio/v2`
