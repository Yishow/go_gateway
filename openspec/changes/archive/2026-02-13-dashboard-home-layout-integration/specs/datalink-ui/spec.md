## MODIFIED Requirements

### Requirement: Sidebar navigation improvements

The system SHALL provide dashboard-first navigation aligned to the one-screen operator workflow.

#### Scenario: Pipeline Studio as primary entry
- **WHEN** a user enters datalink operations
- **THEN** the primary operation entry is the Pipeline Studio dashboard (`/datalink`)
- **AND** mapping operations are completed inside this entry

#### Scenario: Dashboard-integrated navigation replaces persistent sidebar
- **WHEN** the dashboard is rendered
- **THEN** primary section navigation is provided inside dashboard top controls
- **AND** the system SHALL NOT require a persistent external sidebar wrapper for core operations

#### Scenario: Legacy mapping pages decommissioned
- **WHEN** a user accesses `/datalink/points`, `/datalink/mappings`, or `/datalink/wizard`
- **THEN** the system redirects to `/datalink`
- **AND** shows a migration notice that the workflow has moved to Pipeline Studio

#### Scenario: Legacy section routes preserve intent
- **WHEN** a user accesses `/datalink/devices` or `/datalink/settings`
- **THEN** the system redirects to `/datalink` with section intent preserved
- **AND** the corresponding dashboard-integrated section is activated

#### Scenario: Device creation route migrated to modal flow
- **WHEN** a user accesses `/datalink/devices/new`
- **THEN** the system redirects to `/datalink` and opens device creation modal flow
- **AND** the user can complete creation without leaving dashboard workspace

#### Scenario: Device creation stepper
- **WHEN** the create/edit device modal is opened
- **THEN** the modal presents step-based flow (identity, connection, validation)
- **AND** after successful create the UI offers "switch to this device now"

#### Scenario: Local Modbus workbench route
- **WHEN** a user accesses `/datalink/local-modbus`
- **THEN** the system opens a dedicated page for local `5020` memory grid operations
- **AND** the page provides complete local sink management instead of condensed dashboard card actions
