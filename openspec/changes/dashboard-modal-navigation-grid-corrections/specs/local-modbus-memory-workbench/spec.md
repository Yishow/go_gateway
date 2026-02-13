## MODIFIED Requirements

### Requirement: Local Modbus workbench route

The system SHALL provide explicit and discoverable links from dashboard to `/datalink/local-modbus`.

#### Scenario: Dashboard has persistent workbench entry
- **WHEN** a user is on `/datalink`
- **THEN** the UI shows a persistent `Server Memory Grid` entry in dashboard controls
- **AND** selecting it navigates to `/datalink/local-modbus`

#### Scenario: Workbench supports dashboard return
- **WHEN** a user is on `/datalink/local-modbus`
- **THEN** the UI provides a return action to `/datalink`
- **AND** preserves section context if present
