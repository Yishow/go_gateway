## ADDED Requirements

### Requirement: Workspace bootstrap hydrates persisted setup state before rendering steps

The system SHALL fetch and compose the persisted Step 1, Step 2, Step 3, and Step 4 workspace state before the Studio V2 shell renders editable step content.

#### Scenario: Reload restores persisted setup state

- **WHEN** an operator reloads /studio/v2 after saving devices, source rules, mappings, database connector settings, and database targets
- **THEN** the shell renders those persisted values from storage before showing editable step panels
- **AND** the shell SHALL NOT fall back to factory defaults for any field that already has a persisted value

### Requirement: Reload makes unrecovered drafts explicit

The system SHALL distinguish persisted setup state from unsaved in-memory drafts and SHALL surface when a hard refresh discards unrecovered local drafts.

#### Scenario: Hard refresh discards unsaved local draft

- **WHEN** an operator edits a Step 1, Step 2, Step 3, or Step 4 field that has not reached persisted state and reloads the page
- **THEN** the reloaded shell restores the last persisted value for that field
- **AND** the shell surfaces that the unrecovered draft was discarded instead of silently pretending it was saved

### Requirement: Hydration preserves cross-step identities

The system SHALL restore persisted Step 1 through Step 4 state in a dependency-safe order so rule, point, mapping, and database-target relationships remain aligned after reload.

#### Scenario: Mapping hydration keeps current point identity alignment

- **WHEN** the workspace has persisted source rules, derived points, persisted mappings, and persisted database targets
- **THEN** reload hydrates devices and rules before reconstructing mappings and database targets
- **AND** each restored mapping and database target remains bound to the correct current point identity instead of a stale in-memory index

##### Example: two persisted mappings survive point reconstruction

- **GIVEN** persisted rules generate points pt-A and pt-B, persisted mappings target pt-A and pt-B, and persisted database targets also bind pt-A and pt-B
- **WHEN** the shell reloads and rebuilds points from persisted rules before rehydrating mappings and database targets
- **THEN** the restored mapping for pt-A binds the rebuilt pt-A row, the restored mapping for pt-B binds the rebuilt pt-B row, and no mapping is reassigned only by array position
