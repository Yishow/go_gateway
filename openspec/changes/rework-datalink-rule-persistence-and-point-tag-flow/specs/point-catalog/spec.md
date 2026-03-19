## ADDED Requirements

### Requirement: Rule-derived points are first-class runtime assets
The system SHALL allow persisted source rules to derive and manage points as runtime-facing assets in the primary workbench flow.

#### Scenario: Rule creation derives points
- **WHEN** an operator creates or enables a source rule
- **THEN** the system SHALL create or synchronize the derived points needed for that rule
- **AND** SHALL associate those points with the rule's device and address model context

#### Scenario: Unmanaged points remain distinguishable
- **WHEN** a point exists without a current persisted rule association
- **THEN** the system SHALL preserve the point record
- **AND** SHALL expose it as unmanaged or legacy state in planning contexts instead of silently treating it as a rule-derived point

### Requirement: Point collection follows rule lifecycle
The system SHALL stop or resume runtime collection for rule-derived points according to the owning rule's enabled state.

#### Scenario: Disable rule preserves point definition
- **WHEN** an operator disables a persisted rule
- **THEN** the derived point definitions SHALL remain stored
- **AND** runtime collection for those points SHALL stop until the rule is re-enabled
