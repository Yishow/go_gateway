## ADDED Requirements

### Requirement: Studio v2 autosave responses expose runtime apply state

The system SHALL include runtime apply state in successful V2 autosave responses for device-affecting updates.

#### Scenario: Running device response reports applied

- **WHEN** a valid autosave change succeeds for a running workspace device
- **THEN** the success payload includes `runtime_apply_status = "applied"`

#### Scenario: Not-running device response reports not running

- **WHEN** a valid autosave change succeeds for a workspace device that has never been activated
- **THEN** the success payload includes `runtime_apply_status = "not_running"`

#### Scenario: Apply failure is distinct from save success

- **WHEN** the persisted save succeeds but runtime apply fails for a running workspace device
- **THEN** the success payload includes `runtime_apply_status = "apply_failed"`
- **AND** the payload includes a failure message explaining why runtime apply did not complete
