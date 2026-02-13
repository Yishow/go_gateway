## MODIFIED Requirements
### Requirement: Guided workflow

The UI SHALL provide a guided workflow that keeps a persistent visual context of the complete data flow: **Source Device Data → Memory Grid Address → Tag Mapping → Database Write Target**.

The workflow SHALL ensure:

1. Operators can always identify the currently selected source device and address.
2. Mapping and transform steps are editable without losing source context.
3. Validation and activation states are visible in the same workspace.
4. Failures are localized to a specific flow segment.

#### Scenario: End-to-end guided configuration in one workspace
- **WHEN** an operator selects a source address and configures mapping/tag/write options
- **THEN** the UI keeps source, mapping, and write target context visible
- **AND** the operator does not need to switch pages to complete activation

#### Scenario: Segment-localized error handling
- **WHEN** validation fails in transform or write stage
- **THEN** the UI marks the failed segment
- **AND** provides actionable retry or edit guidance for that segment

## ADDED Requirements
### Requirement: Flow-first workspace visualization

The UI SHALL provide a flow-first workspace with four explicit sections:
- Source (device/protocol/connection health)
- Memory Grid (address selection and occupancy)
- Tag Linkage (point-to-tag mapping)
- Write Target (storage status and write readiness)

#### Scenario: Persistent flow visualization
- **WHEN** the operator changes selected devices or addresses
- **THEN** the workspace updates all four sections cohesively
- **AND** preserves a consistent flow reading order

#### Scenario: Memory-to-tag linkage visibility
- **WHEN** an address is linked to a tag
- **THEN** the UI shows the linkage immediately in both grid context and tag context
- **AND** displays current linkage status (draft/validated/active/error)

### Requirement: Accessible and responsive operator workspace

The UI SHALL remain fully operable by keyboard and avoid horizontal overflow at supported desktop breakpoints.

#### Scenario: Keyboard-only operation
- **WHEN** an operator uses keyboard-only navigation
- **THEN** the operator can complete address selection, mapping, validation, and activation
- **AND** all interactive controls have visible focus state

#### Scenario: Desktop responsive stability
- **WHEN** viewport width is 1024px or above
- **THEN** the workspace shows all core flow functions without horizontal scrolling
- **AND** critical actions remain visible without hidden overflow traps
