## MODIFIED Requirements

### Requirement: Desktop workbench shell
The system SHALL provide a desktop-first workbench shell at `/studio` as the primary product workbench route with persistent workbench regions for step navigation, context, main review surfaces, inspection, and summary feedback.

Any legacy `/datalink/workbench` entry SHALL redirect to `/studio` instead of remaining a separate primary workflow.
The system SHALL allow a setup entry page to exist before the workbench, but that entry page SHALL NOT replace direct access to `/studio`.

#### Scenario: Primary workbench opens through `/studio`
- **WHEN** the operator opens the primary datalink workbench
- **THEN** the system loads the desktop-first shell at `/studio`
- **AND** does not require a separate primary `/datalink/workbench` route for the same workflow

#### Scenario: Setup entry page does not block direct workbench access
- **WHEN** the operator opens `/studio` directly while the setup entry page is enabled
- **THEN** the system loads the desktop-first shell at `/studio`
- **AND** does not redirect the operator back to `/gateway/entry`
