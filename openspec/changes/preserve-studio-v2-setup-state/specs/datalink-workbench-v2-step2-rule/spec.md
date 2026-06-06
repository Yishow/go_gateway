## ADDED Requirements

### Requirement: Step 2 reload restores persisted source-rule definitions and save truth

The Step 2 source-rule workspace SHALL restore persisted source-rule definitions and their persisted save truth after reload.

#### Scenario: Reload shows persisted source rules

- **WHEN** the operator reloads Studio V2 after source rules have been persisted for the workspace
- **THEN** Step 2 renders those persisted source rules with the correct owning device, enabled state, and persisted save state indicators
- **AND** Step 2 SHALL NOT regenerate only local placeholder rules when persisted rules already exist
