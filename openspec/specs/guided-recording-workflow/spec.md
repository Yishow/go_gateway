# guided-recording-workflow Specification

## Purpose

TBD - created by archiving change 'redesign-studio-recording-flow'. Update Purpose after archive.

## Requirements

### Requirement: Four-step intent-led setup
The Studio V2 workflow SHALL guide device connection, source selection, value-and-purpose confirmation and recording setup within the existing route.

#### Scenario: Mixed sensor configured without SQL
- **WHEN** an operator assigns temperature, pressure, flow and counter roles to valid source items
- **THEN** recording options are generated from those confirmed roles and managed storage requires no hand-written SQL.

---
### Requirement: Safe reusable templates and overrides
The system SHALL preview template changes per device and preserve confirmed user overrides and stable identities.

#### Scenario: Reapply a template
- **WHEN** a template is applied again after a user edits a pressure multiplier
- **THEN** the difference is shown and the multiplier is not silently overwritten or duplicated.

---
### Requirement: Natural-language recording policies
The system SHALL distinguish acquisition, raw retention, summaries, transport batching and history retention with examples and a final readable summary.

#### Scenario: Five-second polling and one-minute summaries
- **WHEN** those settings are selected
- **THEN** the summary explains that sampled readings feed the minute statistics rather than recording only one instantaneous value per minute.

---
### Requirement: Truthful save and activation states
The system MUST separate saved draft, applied revision, collecting, durable, delivered and verified states, including bounded waiting and partial failure.

#### Scenario: Save during an active plan
- **WHEN** the operator edits a running plan
- **THEN** the current applied plan remains active until explicit validated apply and both versions are visible.

---
### Requirement: Actionable accessible failures
The system SHALL provide field-level reasons, focused repair actions and keyboard-accessible controls instead of unexplained disabled buttons or raw technical errors.

#### Scenario: Write permission denied
- **WHEN** a connected database refuses write access
- **THEN** the operator sees which destination needs permission repair and may retry it without redoing successful device setup.

---
### Requirement: Optional outputs remain independent
The workflow MUST allow a selected Local Modbus-only setup without requiring database configuration while retaining the existing Share activation gates.

#### Scenario: No database output selected
- **WHEN** the user selects only validated Local Modbus forwarding
- **THEN** database fields and tests are skipped and no database writer starts.
