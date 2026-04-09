## ADDED Requirements

### Requirement: Local Modbus is a share destination workspace

The system SHALL treat Local Modbus as a share destination under `/studio/share/local-modbus`, not as a required sub-surface of the primary product workflow completion path.

#### Scenario: Local Modbus is reached through share flow
- **WHEN** an operator chooses share delivery from the destination hub
- **THEN** the operator can enter the Local Modbus workspace through the share path
- **AND** the operator is not required to configure Local Modbus to complete a database-first workflow

### Requirement: Local Modbus uses block-first planning

The Local Modbus workspace SHALL use delivery groups and block-oriented planning as the primary interaction model.

#### Scenario: One delivery group becomes one register block plan
- **WHEN** an operator selects a delivery group for Local Modbus
- **THEN** the UI plans a contiguous or strategy-driven register block
- **AND** does not require per-tag register assignment as the primary default flow
