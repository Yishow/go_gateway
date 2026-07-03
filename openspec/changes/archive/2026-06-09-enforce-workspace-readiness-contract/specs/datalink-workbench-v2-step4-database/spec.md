## ADDED Requirements

### Requirement: Step 4 activation surface shows readiness blockers and warnings

The Step 4 database workspace SHALL show workspace readiness blockers and warnings before the operator starts activation.

#### Scenario: Activation control shows readiness blockers

- **WHEN** the workspace readiness result contains blocking or warning issues
- **THEN** Step 4 surfaces those readiness issues near the activation action
- **AND** the operator can distinguish whether activation is blocked or merely warned before clicking the activation button
