## ADDED Requirements

### Requirement: Step 1 reload restores persisted device definitions and save truth

The Step 1 device workspace SHALL restore persisted device definitions and their persisted save truth after reload.

#### Scenario: Reload shows persisted device state instead of defaults

- **WHEN** the operator reloads Studio V2 after one or more devices have already been persisted
- **THEN** Step 1 renders those persisted device definitions, protocol settings, and persisted save state indicators
- **AND** Step 1 SHALL NOT show the default single-device form for devices that already exist in persisted storage
