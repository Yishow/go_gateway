## ADDED Requirements

### Requirement: Studio uses destination-oriented workflow

The UI SHALL provide `/studio` as a workflow-oriented product route using the step model:

- Device
- Source
- Tag
- Destination

#### Scenario: Destination replaces old Output framing
- **WHEN** an operator reaches the final workflow step
- **THEN** the UI presents `Destination` instead of a unified old Output framing
- **AND** the operator first lands on a destination hub before entering a specific destination workspace

### Requirement: Test route remains a first-class engineering console

The system SHALL keep `/test` as an independent engineering route and SHALL NOT fold it into the `/studio` product workflow.

#### Scenario: Test remains independent
- **WHEN** a field engineer opens `/test`
- **THEN** the UI presents an engineering debug console
- **AND** does not require entering or completing `/studio`
