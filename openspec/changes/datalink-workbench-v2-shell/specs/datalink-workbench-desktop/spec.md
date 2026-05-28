## ADDED Requirements

### Requirement: Coexisting v2 workbench at `/studio/v2`

The system SHALL accept `/studio/v2` as a parallel workbench route alongside the existing `/studio` primary route. The existing `/studio` route MUST continue to load the legacy `DatalinkWorkbenchPage` without modification. The `/studio/v2` route MUST load the Workbench v2 shell defined by capability `datalink-workbench-v2-shell`. Neither route MUST redirect to the other automatically.

#### Scenario: /studio remains the primary route

- **WHEN** the operator navigates to `/studio`
- **THEN** the system loads the pre-existing `DatalinkWorkbenchPage` shell
- **AND** does not redirect to `/studio/v2`

#### Scenario: /studio/v2 loads the v2 shell

- **WHEN** the operator navigates to `/studio/v2`
- **THEN** the system loads the Workbench v2 shell
- **AND** does not redirect to `/studio`

#### Scenario: Legacy datalink routes still redirect to /studio

- **WHEN** the operator navigates to `/datalink/workbench` or any `/datalink/workbench/*` path
- **THEN** the system redirects to `/studio` exactly as before this change
- **AND** does not redirect to `/studio/v2`
