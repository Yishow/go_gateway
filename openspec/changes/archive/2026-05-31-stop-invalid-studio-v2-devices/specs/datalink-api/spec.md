## ADDED Requirements

### Requirement: Studio v2 device availability status

The system SHALL include availability status in workspace and runtime-facing device payloads used by V2.

#### Scenario: Available device payload

- **WHEN** a workspace device is currently valid and usable
- **THEN** the payload includes `availability_status = "available"`

#### Scenario: Unavailable device payload

- **WHEN** a workspace device becomes invalid after editing
- **THEN** the payload includes `availability_status = "unavailable"`
- **AND** the payload includes `availability_reason`
