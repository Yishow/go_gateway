## ADDED Requirements

### Requirement: Studio v2 settings API reuse

The system SHALL provide a backend contract that Studio V2 settings can use for boot and persistence without mock-only fallbacks.

#### Scenario: Read settings for V2 boot

- **WHEN** a client requests settings data for `/studio/v2/settings`
- **THEN** the API returns the persisted backend values needed for the settings surface

#### Scenario: Persist changed settings

- **WHEN** a client submits changed settings from `/studio/v2/settings`
- **THEN** the backend persists those values
- **AND** the client receives success or actionable failure rather than relying on console warnings

### Requirement: Studio v2 connector pool API reuse

The system SHALL let Studio V2 settings use real connector CRUD and connector test APIs.

#### Scenario: Test one connector

- **WHEN** a client requests connector test for a settings connector
- **THEN** the backend performs a real connector test
- **AND** the response is sufficient for the UI to render success or failure without randomness
