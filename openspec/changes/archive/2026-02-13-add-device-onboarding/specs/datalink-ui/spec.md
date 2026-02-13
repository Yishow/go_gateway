## ADDED Requirements

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
