# datalink-ui Specification

## Purpose

TBD - created by archiving change add-device-data-pipeline. Update Purpose after archive.
## Requirements
### Requirement: Guided workflow

The UI SHALL provide a guided workflow that keeps a persistent visual context of the complete data flow in one workspace: **Source Plan -> Memory Grid Allocation -> Tag Mapping -> Database Commit**.

The workflow SHALL ensure:

1. Operators can define source quantity and source data type before allocation.
2. Mapping and tag operations remain visible without page switches.
3. Validation and commit states are visible in the same workspace.
4. Failures are localized to a specific flow segment.

#### Scenario: End-to-end guided configuration in one workspace
- **WHEN** an operator defines source type/count and completes allocation, tag, and commit actions
- **THEN** the UI keeps source, grid, tag, and commit context visible together
- **AND** the operator does not need to switch pages to complete DB commit

#### Scenario: Segment-localized error handling
- **WHEN** validation or commit fails
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

The system SHALL migrate legacy sidebar-aligned feature routes into dashboard modal workflows, except the independent `/test` page.

#### Scenario: Sidebar feature routes open modal workflows
- **WHEN** a user accesses legacy feature paths (`/datalink/devices`, `/datalink/settings`, `/datalink/points`, `/datalink/mappings`, `/datalink/wizard`)
- **THEN** the system redirects to `/datalink`
- **AND** opens the corresponding dashboard modal context

#### Scenario: Test page remains independent
- **WHEN** a user accesses `/test`
- **THEN** the system keeps `/test` as an independent page
- **AND** does not convert it into dashboard modal flow

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

### Requirement: Source template library

The UI SHALL allow operators to save, load, update, and delete source planning templates.

Each template SHALL store at least protocol context, source data type, source count, and naming defaults.

#### Scenario: Save and reuse source template
- **WHEN** an operator saves a source plan as a template
- **THEN** the template is persisted
- **AND** the operator can later load it to prefill planning controls

#### Scenario: Edit existing source template
- **WHEN** an operator updates a saved template
- **THEN** the updated defaults are used in subsequent planning sessions

### Requirement: Motion-guided operator flow

The UI SHALL use motion cues to guide stage transitions between source planning, grid allocation, tag linkage, and DB commit.

Animations MUST use short transitions (150-300ms) and MUST support reduced-motion preference.

#### Scenario: Guided transition after source planning
- **WHEN** an operator confirms source plan
- **THEN** the grid allocation region receives a transition cue indicating next action

#### Scenario: Reduced-motion mode
- **WHEN** user preference is `prefers-reduced-motion`
- **THEN** motion cues are replaced with static visual emphasis without animation

### Requirement: Planning intelligence and safety checks

The UI SHALL provide planning assistance and safety checks for one-screen operator execution.

#### Scenario: Naming preview and duplicate detection
- **WHEN** an operator defines batch naming rules
- **THEN** the UI previews generated names
- **AND** flags duplicates or naming conflicts before commit

#### Scenario: Two-stage validation execution
- **WHEN** an operator validates pending changes
- **THEN** the UI runs structural validation first
- **AND** runs executable validation only after structural validation succeeds

