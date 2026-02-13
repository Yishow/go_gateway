# datalink-ui Specification

## Purpose

TBD - created by archiving change add-device-data-pipeline. Update Purpose after archive.
## Requirements
### Requirement: Guided workflow

The UI SHALL provide a guided workflow that keeps a persistent visual context of the complete data flow: **Source Device Data → Memory Grid Address → Tag Mapping → Database Write Target**.

The workflow SHALL ensure:

1. Operators can always identify the currently selected source device and address.
2. Mapping and transform steps are editable without losing source context.
3. Validation and activation states are visible in the same workspace.
4. Failures are localized to a specific flow segment.

#### Scenario: End-to-end guided configuration in one workspace
- **WHEN** an operator selects a source address and configures mapping/tag/write options
- **THEN** the UI keeps source, mapping, and write target context visible
- **AND** the operator does not need to switch pages to complete activation

#### Scenario: Segment-localized error handling
- **WHEN** validation fails in transform or write stage
- **THEN** the UI marks the failed segment
- **AND** provides actionable retry or edit guidance for that segment

### Requirement: Drag-drop mapping canvas

The UI SHALL provide a drag-drop canvas to connect source points to target tags and mapping steps.

#### Scenario: Create mapping by drag

- WHEN a user drags a point onto a tag
- THEN a mapping is created and editable

#### Scenario: Edit existing mapping

- WHEN a user clicks on an existing mapping
- THEN the mapping details are shown for editing

---

### Requirement: Transform builder

The UI SHALL provide a transform builder with ordered steps, parameters, and validation.

The builder SHALL support the following transform types:

- `decode`: Decode raw bytes to typed value
- `cast`: Type conversion
- `scale`: Linear scaling with parameters (multiplier, offset)
- `lookup`: Table lookup mapping
- `conditional`: Conditional branching
- `formula`: Expression-based calculation

#### Scenario: Configure scaling step

- WHEN a user adds a scaling step with parameters
- THEN the UI validates and saves the step

#### Scenario: Reorder transform steps

- WHEN a user drags a step to a new position
- THEN the pipeline order is updated

#### Scenario: Pipeline validation

- WHEN a user saves the transform pipeline
- THEN the UI calls `/mappings/validate-pipeline` to verify

---

### Requirement: Live preview

The UI SHALL provide a live preview of raw and transformed values for a selected mapping using Server-Sent Events (SSE).

The preview panel SHALL display:

- Raw value from device
- Intermediate step results
- Final transformed value
- Quality indicator

#### Scenario: View preview

- WHEN a user opens preview
- THEN the UI displays raw, step results, and final value

#### Scenario: Live update via SSE

- WHEN device values change
- THEN the preview panel updates in real-time

#### Scenario: SSE reconnection

- WHEN the SSE connection is lost
- THEN the UI automatically reconnects

---

### Requirement: Write precision settings

The UI SHALL allow operators to configure write timestamp precision (seconds or milliseconds).

#### Scenario: Select millisecond precision

- WHEN a user selects millisecond precision
- THEN the UI saves the setting for subsequent writes

#### Scenario: Configure partition interval

- WHEN a user selects a partition interval (daily/weekly/monthly)
- THEN the UI saves the setting for time-series storage

### Requirement: Query-based state management

The UI SHALL use TanStack Query for server state management with caching, automatic refetching, and optimistic updates.

#### Scenario: Cached data loading

- WHEN a user navigates to a previously visited page
- THEN the UI displays cached data immediately while revalidating in background

#### Scenario: Optimistic status toggle

- WHEN a user toggles device status
- THEN the UI reflects the change immediately
- AND reverts if the backend operation fails

### Requirement: Device onboarding wizard

The system SHALL provide a guided onboarding wizard component to help users complete device configuration.

#### Scenario: Wizard opens after device creation

- **WHEN** a new device is created successfully
- **THEN** the onboarding wizard automatically opens

#### Scenario: Wizard guides through steps

- **WHEN** a user follows the onboarding wizard
- **THEN** the wizard guides through 7 steps: device creation, connection test, activation, point creation, polling group assignment, tag creation, and mapping creation

#### Scenario: Wizard shows progress

- **WHEN** a user progresses through the wizard
- **THEN** each step shows completion status and validation results

#### Scenario: Wizard validates each step

- **WHEN** a user attempts to proceed to the next step
- **THEN** the wizard validates the current step before allowing progression

### Requirement: Device status dashboard

The system SHALL provide a dashboard page showing device collection status and configuration completeness.

#### Scenario: Display device statistics

- **WHEN** a user views the status dashboard
- **THEN** the dashboard displays summary statistics: total devices, active devices, devices collecting data

#### Scenario: Display device list with status

- **WHEN** a user views the status dashboard
- **THEN** the dashboard displays a list of devices with their collection status, last collection time, and error count

#### Scenario: Real-time status updates

- **WHEN** device collection status changes
- **THEN** the dashboard updates automatically using SSE

#### Scenario: Filter devices by status

- **WHEN** a user filters devices by status
- **THEN** the dashboard shows only devices matching the filter criteria

### Requirement: Device readiness indicator

The system SHALL display readiness status indicators in the device list and device detail pages.

#### Scenario: Show readiness badge

- **WHEN** a user views the device list
- **THEN** each device displays a readiness badge (ready/not ready)

#### Scenario: Show readiness details

- **WHEN** a user clicks on a device readiness badge
- **THEN** the system shows detailed readiness check results and missing configurations

#### Scenario: Show completion progress

- **WHEN** a user views a device detail page
- **THEN** the page displays a configuration completion progress bar

### Requirement: Sidebar navigation improvements

The system SHALL provide improved sidebar navigation with all core features accessible.

#### Scenario: Points navigation item

- **WHEN** a user views the sidebar
- **THEN** the sidebar displays a Points navigation item linking to `/datalink/points`

#### Scenario: Navigation order reflects workflow

- **WHEN** a user views the sidebar
- **THEN** navigation items are ordered to reflect the data pipeline workflow: Dashboard → Devices → Points → Tags → Mappings → Mapping Wizard → Settings

#### Scenario: Visual grouping

- **WHEN** a user views the sidebar
- **THEN** related navigation items are visually grouped (e.g., data collection group, data processing group, system group)

#### Scenario: Responsive sidebar

- **WHEN** a user views the application on a mobile device
- **THEN** the sidebar is displayed as a drawer that can be opened/closed with gestures

### Requirement: Points management page

The system SHALL provide a dedicated Points management page with full CRUD functionality.

#### Scenario: Display points list

- **WHEN** a user navigates to the Points page
- **THEN** the page displays a list of all points with filtering and sorting capabilities

#### Scenario: Filter points by device

- **WHEN** a user filters points by device
- **THEN** the page shows only points belonging to the selected device

#### Scenario: Create point from Points page

- **WHEN** a user creates a new point from the Points page
- **THEN** the point is created and added to the list

#### Scenario: Edit point

- **WHEN** a user edits a point
- **THEN** the point details are updated and reflected in the list

#### Scenario: Delete point

- **WHEN** a user deletes a point
- **THEN** the point is removed from the list after confirmation

#### Scenario: Batch operations

- **WHEN** a user selects multiple points
- **THEN** the page provides batch operations (enable/disable, assign to polling group, delete)

### Requirement: Flow-first workspace visualization

The UI SHALL provide a flow-first workspace with four explicit sections:
- Source (device/protocol/connection health)
- Memory Grid (address selection and occupancy)
- Tag Linkage (point-to-tag mapping)
- Write Target (storage status and write readiness)

#### Scenario: Persistent flow visualization
- **WHEN** the operator changes selected devices or addresses
- **THEN** the workspace updates all four sections cohesively
- **AND** preserves a consistent flow reading order

#### Scenario: Memory-to-tag linkage visibility
- **WHEN** an address is linked to a tag
- **THEN** the UI shows the linkage immediately in both grid context and tag context
- **AND** displays current linkage status (draft/validated/active/error)

### Requirement: Accessible and responsive operator workspace

The UI SHALL remain fully operable by keyboard and avoid horizontal overflow at supported desktop breakpoints.

#### Scenario: Keyboard-only operation
- **WHEN** an operator uses keyboard-only navigation
- **THEN** the operator can complete address selection, mapping, validation, and activation
- **AND** all interactive controls have visible focus state

#### Scenario: Desktop responsive stability
- **WHEN** viewport width is 1024px or above
- **THEN** the workspace shows all core flow functions without horizontal scrolling
- **AND** critical actions remain visible without hidden overflow traps

