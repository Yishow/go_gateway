## ADDED Requirements

### Requirement: Valid-only rule autosave

The Step 2 source-rule workspace SHALL autosave only valid rule content to the backend workspace.

#### Scenario: Valid rule change saves immediately

- **GIVEN** a rule already belongs to the singleton v2 workspace
- **WHEN** the operator edits that rule into a valid state
- **THEN** the system saves the rule immediately
- **AND** the persisted rule revision reflects the latest valid values

#### Scenario: Invalid rule change stays local

- **GIVEN** a rule already has a last successful persisted revision
- **WHEN** the operator edits that rule into an invalid state
- **THEN** the system does not overwrite the persisted revision
- **AND** the UI keeps the invalid local values visible with an unsaved marker

### Requirement: Per-rule save isolation

The Step 2 source-rule workspace SHALL isolate save success and failure per rule.

#### Scenario: One invalid rule does not block another valid rule

- **GIVEN** the workspace contains multiple rules
- **WHEN** one rule remains invalid while another rule is edited into a valid state
- **THEN** the valid rule still saves successfully
- **AND** the invalid rule remains local with its own error state

### Requirement: Workspace-scoped rule ownership

Every Step 2 source rule SHALL preserve its owning workspace and device relationship across reloads.

#### Scenario: Reload restores device ownership

- **WHEN** the operator reloads `/studio/v2`
- **THEN** each persisted rule reappears under the same workspace and owning device as before reload
- **AND** the system does not silently reassign that rule to another device
