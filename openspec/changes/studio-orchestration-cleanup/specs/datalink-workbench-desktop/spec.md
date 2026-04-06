## MODIFIED Requirements

### Requirement: Desktop workbench shell
The system SHALL provide a desktop-first workbench shell at `/studio` as the primary product route with persistent workbench regions for step navigation, context, main review surfaces, inspection, and summary feedback.

Any legacy `/datalink/workbench` entry SHALL redirect to `/studio` instead of remaining a separate primary workflow.

#### Scenario: Primary workbench opens through `/studio`
- **WHEN** the operator opens the primary datalink workbench
- **THEN** the system loads the desktop-first shell at `/studio`
- **AND** does not require a separate primary `/datalink/workbench` route for the same workflow

### Requirement: Unified output workspace supports both targets
The system SHALL provide one Output workspace inside `/studio` that keeps Local Modbus and Database outputs in the same workbench step.

#### Scenario: Operator changes output target without leaving `/studio`
- **WHEN** the operator switches between `Local Modbus` and `Database` in the Output step
- **THEN** the output review surface remains in the same `/studio` route
- **AND** the target-specific studio changes without requiring navigation to a separate primary page

#### Scenario: Output target switching preserves one primary workflow
- **WHEN** legacy navigation or bookmarks would previously open a separate workbench route for output work
- **THEN** the system keeps the operator inside `/studio`
- **AND** treats any retained legacy entry as a compatibility redirect rather than a second primary workflow
