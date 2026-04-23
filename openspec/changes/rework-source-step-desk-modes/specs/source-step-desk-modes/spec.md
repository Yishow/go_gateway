## ADDED Requirements

### Requirement: Source desk modes share one planning skeleton
The system SHALL provide `Inspect`, `Build`, and `Triage` desks inside Step 2 while preserving one shared planning skeleton for the active source rule.

The shared skeleton SHALL keep the active rule summary, grouped Tag handoff context, diagnostic cues, and primary rule action visible without requiring route changes or mode-specific re-entry.

#### Scenario: Switching desks preserves the active rule context
- **WHEN** an operator switches between `Inspect`, `Build`, and `Triage` while the same source rule is active
- **THEN** Step 2 keeps the same active rule summary, grouped handoff copy, and blocker diagnostics visible
- **AND** the operator does not need to reselect the rule to continue planning

#### Scenario: Desk switching stays inside one Source workspace
- **WHEN** an operator changes the Source desk mode
- **THEN** the system updates the Step 2 workspace in place
- **AND** the workflow remains inside the same `/studio` Source step without opening a separate route or modal-only flow

### Requirement: Inspect and Build desks expose distinct working emphasis
The system SHALL make `Inspect` the canvas-first verification desk and SHALL make `Build` the rule-editor-first planning desk.

`Inspect` SHALL prioritize reading planned spans, live address context, and downstream readiness on the shared canvas. `Build` SHALL prioritize editing the active rule while still keeping a readable canvas region visible for immediate planning feedback.

#### Scenario: Inspect desk emphasizes canvas verification
- **WHEN** the operator opens `Inspect`
- **THEN** the main Step 2 workspace prioritizes the shared canvas and planning verification surfaces
- **AND** the rule-editing controls remain available as secondary support instead of dominating the viewport

#### Scenario: Build desk keeps the editor primary without hiding the canvas
- **WHEN** the operator opens `Build` for an active source rule
- **THEN** the rule editor becomes the primary workspace for editing
- **AND** the canvas remains visible with enough size to review planned spans, unmanaged gaps, and current address coverage without leaving the desk

##### Example: desk emphasis matrix
| Desk | Primary surface | Secondary surface |
| ---- | --------------- | ----------------- |
| Inspect | Shared source canvas | Rule details and diagnostics |
| Build | Active rule editor | Readable source canvas |

### Requirement: Triage desk is issue-first and explicit when clear
The system SHALL make `Triage` the issue-first recovery desk for conflict, blocked, and unmanaged planning states.

`Triage` SHALL surface recovery actions and the affected planning context before general inspection content. When no triage-worthy issues exist, the desk SHALL render an explicit clear-state message instead of falling back to an inspect-like layout.

#### Scenario: Triage desk elevates incident recovery
- **WHEN** conflict, blocked, or unmanaged items exist for the active source rule
- **THEN** the `Triage` desk shows those items and their recovery actions as the primary review surface
- **AND** the desk limits canvas context to the affected ranges needed to understand and repair the issue

#### Scenario: Triage desk shows a clear-state message when no issues exist
- **WHEN** the operator opens `Triage` and the active source rule has no conflict, blocked, or unmanaged items requiring intervention
- **THEN** the desk shows a clear-state explanation that no triage work is pending
- **AND** the desk guides the operator back to `Inspect` or `Build` instead of rendering a generic canvas-first workspace
