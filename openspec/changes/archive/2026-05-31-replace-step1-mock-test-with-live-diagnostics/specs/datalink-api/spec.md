## ADDED Requirements

### Requirement: Studio v2 Step 1 diagnostics API

The system SHALL provide a real backend diagnostics contract that Studio V2 Step 1 can use for the current device draft.

#### Scenario: Successful diagnostics response

- **WHEN** a client requests Step 1 diagnostics for the current device draft
- **THEN** the API returns connect / probe outcomes produced by backend logic
- **AND** the response is sufficient for the UI to render success without local synthesis

#### Scenario: Failed diagnostics response

- **WHEN** backend diagnostics fails at connect or probe stage
- **THEN** the API returns the failing stage and actionable message
- **AND** the client does not need to guess failure state from a generic exception alone
