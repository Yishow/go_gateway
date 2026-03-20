# point-catalog Specification

## Purpose
TBD - created by archiving change add-device-data-pipeline. Update Purpose after archive.
## Requirements
### Requirement: Point definition
The system SHALL allow operators to define points per device with address, function, and data type details required by the protocol.

#### Scenario: Add holding register point
- WHEN a point is created with address, function, and data type
- THEN the point is stored and available for scheduling

### Requirement: Read and write modes
The system SHALL allow points to be marked as read-only or read-write.

#### Scenario: Read-only point
- WHEN a point is set to read-only
- THEN write operations for that point are rejected

### Requirement: Polling group configuration
The system SHALL allow points to be assigned to polling groups with a configured interval.

#### Scenario: Assign to fast group
- WHEN a point is assigned to a group with a 1s interval
- THEN the scheduler reads the point at the group interval

### Requirement: Point status and health
The system SHALL track last read time, last value, and error status per point.

#### Scenario: Read failure
- WHEN a point read fails
- THEN the point status records the error and timestamp

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

