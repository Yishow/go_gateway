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

