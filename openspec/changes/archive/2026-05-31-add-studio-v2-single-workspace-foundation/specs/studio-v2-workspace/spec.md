## ADDED Requirements

### Requirement: Singleton studio v2 workspace

The system SHALL maintain exactly one persisted workspace for `/studio/v2`.

#### Scenario: First load creates the singleton workspace

- **WHEN** the operator opens `/studio/v2` and no v2 workspace exists yet
- **THEN** the system creates one persisted workspace record
- **AND** the same request returns that workspace instead of requiring a separate creation step

#### Scenario: Restart keeps the same workspace

- **WHEN** the service restarts after a v2 workspace has been created
- **THEN** the next `/studio/v2` load returns the same workspace record
- **AND** the workspace id is preserved across restart

### Requirement: First-open workspace bootstrap

The system SHALL bootstrap the singleton workspace on page open without showing a separate create-workspace screen.

#### Scenario: Page open returns bootstrap metadata

- **WHEN** the operator opens `/studio/v2`
- **THEN** the system returns workspace bootstrap metadata including id, status, and ordered device skeleton
- **AND** the page can continue booting from that response alone

### Requirement: Workspace excludes legacy /studio data

The singleton workspace SHALL start isolated from the legacy `/studio` mainline.

#### Scenario: Legacy data is not imported on bootstrap

- **GIVEN** legacy `/studio` data already exists in the service
- **WHEN** the operator opens `/studio/v2` for the first time
- **THEN** the new singleton workspace does not automatically include those legacy devices, rules, or outputs
- **AND** the operator starts from a V2-owned workspace boundary
