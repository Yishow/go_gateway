## ADDED Requirements

### Requirement: Shell status reflects restored activation truth after reload

The Studio V2 shell SHALL derive its committed and runtime status indicators from restored persisted activation truth after reload.

#### Scenario: Reload preserves committed shell truth

- **WHEN** the workspace has already been activated and the operator reloads /studio/v2
- **THEN** the shell restores the committed/runtime indicator from persisted activation truth
- **AND** the shell SHALL NOT reset the top-bar scheduler indicator to idle only because the browser session restarted
