## MODIFIED Requirements

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
