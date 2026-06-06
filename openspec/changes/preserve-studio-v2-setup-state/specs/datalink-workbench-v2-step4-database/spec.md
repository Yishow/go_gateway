## ADDED Requirements

### Requirement: Step 4 reload restores persisted database setup and save truth

The Step 4 database workspace SHALL restore the persisted database connector, database targets, and their persisted save truth after reload.

#### Scenario: Reload shows persisted database setup

- **WHEN** the operator reloads Studio V2 after persisting database connector settings and database targets
- **THEN** Step 4 renders the persisted connector fields, target bindings, and persisted save state indicators
- **AND** Step 4 SHALL NOT silently replace persisted connector or target state with default placeholders
