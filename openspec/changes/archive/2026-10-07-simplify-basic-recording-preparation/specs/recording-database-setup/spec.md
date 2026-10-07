## ADDED Requirements

### Requirement: Basic preparation retains schema and capability safety
Basic schema preparation SHALL reuse canonical group preview, action-bound tokens, explicit DDL confirmation, operation recovery and scope guards. Unsupported group database kinds SHALL be disabled with safe actionable explanations, including saved connector selection. The legacy nonpreview schema action SHALL NOT be a mainline preparation path.

#### Scenario: Cancel or invalidate basic preview
- **WHEN** a Basic schema preview is cancelled or its group, device, connector revision or workspace scope changes before confirmation
- **THEN** no implicit DDL is executed and invalid confirmation is rejected by the existing guards.

#### Scenario: Unsupported group kind
- **WHEN** MySQL or SQLServer is offered through kind or saved connector choices without group capability
- **THEN** selection is disabled with a support explanation and no new driver is used.
