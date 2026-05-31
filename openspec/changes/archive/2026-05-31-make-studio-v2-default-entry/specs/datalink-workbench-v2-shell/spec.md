## ADDED Requirements

### Requirement: Default v2 entry route

The system SHALL treat `/studio/v2` as the default guided entry for the datalink setup flow.

#### Scenario: Root path enters v2

- **WHEN** an operator opens `/`
- **THEN** the system redirects to `/studio/v2`
- **AND** the first rendered setup surface is the v2 shell rather than the legacy `/studio` page

#### Scenario: Unknown path falls back to v2

- **WHEN** an operator opens an unknown route that does not match any explicit surface
- **THEN** the system redirects to `/studio/v2`
- **AND** the fallback does not land on `/studio`

#### Scenario: Direct legacy mainline remains available

- **WHEN** an operator opens `/studio`
- **THEN** the system still renders the legacy workbench page
- **AND** the route does not auto-redirect to `/studio/v2`

### Requirement: Legacy landing paths converge on /studio/v2

The system SHALL route generic legacy datalink landing paths to `/studio/v2` while preserving task-specific deep-link behavior outside this change.

#### Scenario: Generic legacy datalink root lands on v2

- **WHEN** an operator opens `/datalink`
- **THEN** the system redirects to `/studio/v2`
- **AND** the operator does not first land on `/studio`

#### Scenario: Generic legacy workbench path lands on v2

- **WHEN** an operator opens `/datalink/workbench`
- **THEN** the system redirects to `/studio/v2`
- **AND** the operator does not first land on `/studio`
