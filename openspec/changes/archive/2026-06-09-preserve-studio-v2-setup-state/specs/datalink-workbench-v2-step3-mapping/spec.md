## ADDED Requirements

### Requirement: Step 3 reload restores persisted mapping rows against current point identities

The Step 3 mapping workspace SHALL restore persisted mappings against the current derived point identities after reload.

#### Scenario: Reload reattaches persisted mappings to current points

- **WHEN** the operator reloads Studio V2 after mappings have already been persisted for derived points
- **THEN** Step 3 rehydrates those persisted mappings onto the correct current point rows
- **AND** Step 3 SHALL NOT silently drop persisted mappings only because point rows were reconstructed during bootstrap
