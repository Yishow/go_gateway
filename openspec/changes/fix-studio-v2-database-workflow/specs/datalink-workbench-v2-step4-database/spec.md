## ADDED Requirements

### Requirement: Protected row editing and single submission
The target editor SHALL preserve each dirty draft independently from saved server state, surface same-row concurrent changes and remove obsolete row state when a row is deleted. Keyboard and blur handling MUST submit one logical edit at most once, and IME composition MUST NOT trigger premature submission.

#### Scenario: Another row is refreshed
- **WHEN** a save result for row A arrives while the operator edits row B
- **THEN** row B retains its typed text and row A can update normally.

#### Scenario: Same row changes remotely
- **WHEN** a newer server value arrives for the dirty row
- **THEN** the UI preserves the draft and offers an explicit conflict-resolution action instead of silently overwriting either value.

#### Scenario: Enter followed by blur
- **WHEN** Enter ends editing and blur follows, including repeated key events while saving
- **THEN** at most one request is issued for the same logical edit.

#### Scenario: Chinese composition or deleted row
- **WHEN** Enter is pressed during IME composition or a row is removed while a result is pending
- **THEN** composition does not submit and removed-row drafts or late results do not reappear.

### Requirement: Fresh schema results and comprehensive readonly controls
The UI SHALL associate each schema action and displayed result with the exact selected scope and setup version. Stale responses MUST NOT update or authorize a new scope. Readonly and in-flight restrictions SHALL cover all mutating controls, including schema apply, test write, connector changes and bulk edits. Ignoring a response MUST NOT be represented as confirmed backend cancellation.

#### Scenario: Old preview arrives after switching tables
- **WHEN** a preview for table A completes after the operator selects table B
- **THEN** the result cannot appear as table B's preview or enable its apply action.

#### Scenario: Activation or readonly state is active
- **WHEN** activation is in progress or the completed setup is readonly
- **THEN** schema creation, test writes and other configuration mutations are unavailable until explicit return to editing.

#### Scenario: Rapid clicks or navigation during a write
- **WHEN** the operator repeats an action or changes the view while execution remains unresolved
- **THEN** the UI prevents a duplicate write and explains the outstanding result without claiming cancellation.
