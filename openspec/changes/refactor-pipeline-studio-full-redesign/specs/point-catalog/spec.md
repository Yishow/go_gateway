## ADDED Requirements
### Requirement: Typed allocation planning

The system SHALL support typed allocation planning so each planned source reserves memory span according to its data type before point creation.

#### Scenario: Plan int16 sources
- **WHEN** an operator plans 5 sources of `int16`
- **THEN** the planner reserves 5 address cells

#### Scenario: Plan float32 sources
- **WHEN** an operator plans 10 sources of `float32`
- **THEN** the planner reserves 20 address cells
- **AND** each source is represented by a linked 2-cell group

#### Scenario: Reject overlapping planned spans
- **WHEN** planned spans overlap existing points or reserved spans
- **THEN** the system reports conflict details before commit
