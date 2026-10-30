## Purpose

Make existing recording lifecycle and recovery operations understandable from the simple Studio interface while retaining server authority, forward-only storage changes and truthful outcomes.

## ADDED Requirements

### Requirement: Recording actions follow matching server availability
The system SHALL derive saved-recording actions from matching workspace, group and revision evidence. Stop and remove SHALL require explicit scope and impact confirmation. Removing a recording SHALL NOT delete historical database rows or tables.

#### Scenario: The action facts are stale
- **WHEN** lifecycle evidence does not match the displayed recording revision
- **THEN** the interface disables consequential actions, explains the need to reload evidence and does not guess availability from status.

#### Scenario: A recording is removed
- **WHEN** an operator confirms an allowed removal
- **THEN** only the existing configuration-removal contract is invoked and the interface states that historical database data is retained.

#### Scenario: A referenced source is changed
- **WHEN** an operator attempts to delete or redefine a source used by existing recordings
- **THEN** the interface preserves existing dependency protections and identifies affected recordings or opens the scoped advanced repair flow without creating dangling references.

### Requirement: Resume revalidates and copying creates only a new draft
The system SHALL resume through existing scoped review and start contracts. Copy SHALL prefill a new inactive draft without old identity, operation, readiness, receipt or ownership proof. Unsupported representations or managed-profile deduplication SHALL be explained rather than reported as successful new recordings.

#### Scenario: A stopped recording is resumed
- **WHEN** an operator requests resume
- **THEN** current source, destination, revision and readiness evidence is checked and explicit start confirmation is required before waiting for a new write result.

#### Scenario: A supported recording is copied
- **WHEN** an operator copies a representable recording
- **THEN** an inactive draft is opened with editable descriptive settings, the original is unchanged, and no command authority or success evidence is copied.

#### Scenario: Managed ownership prevents a second copy
- **WHEN** submitting a copied managed profile would resolve to the existing canonical role
- **THEN** the interface explains that a distinct duplicate is not available and does not claim a new group or writer was created.

### Requirement: Table changes use reviewed forward-only successors
The system SHALL use existing successor review, reservation, schema planning, confirmation and scoped start or cutover for supported layout or destination changes. The interface SHALL show affected sources, old and new destinations, retained history and the boundary. It MUST NOT perform destructive in-place layout conversion or historical data migration.

#### Scenario: A supported wide recording is changed to long layout
- **WHEN** the operator selects an allowed transition
- **THEN** the interface reviews a new target table and cutover impact through existing guards, retains the old table and history, and distinguishes preparing, scheduled and effective states.

#### Scenario: The reverse layout is not supported
- **WHEN** an operator requests a layout or destination combination that the existing backend rejects
- **THEN** the interface displays the matching blocker and repair or advanced path without claiming conversion or expanding backend capabilities.

#### Scenario: A first start differs from a cutover
- **WHEN** the successor has no prior effective recording
- **THEN** the interface uses the existing first-start review instead of presenting a completed historical cutover.

### Requirement: Unknown operations recover without replacement commands
The system SHALL retain submitted operation identity and truthful pending or unknown states when dialogs close or navigation changes. Closing a dialog SHALL NOT imply cancellation. Recovery SHALL query the original operation before a replacement mutation is permitted by the existing contract. Existing exact-retry or repreview actions SHALL remain available only after fresh matching status evidence satisfies their original guards and the operator explicitly requests the action.

#### Scenario: A response is lost and the dialog closes
- **WHEN** a schema or activation request has an unknown outcome and the operator closes then reopens the dialog
- **THEN** the same operation is shown for status lookup without automatically resending it or creating a new operation.

#### Scenario: A late result belongs to another scope
- **WHEN** an old operation response arrives after switching the selected recording
- **THEN** it does not update the new recording or enable its actions.

#### Scenario: Fresh evidence permits exact retry or repreview
- **WHEN** recovery returns fresh matching evidence that satisfies the existing not-submitted or preview-not-submitted retry guards
- **THEN** the interface offers explicit exact retry with the original request identity, or the existing repreview flow with its permitted new preview identity and retained root-operation lineage and prior previews; neither action weakens its original proof requirements.

#### Scenario: Recovery does not prove retry is safe
- **WHEN** the operation status remains unknown, stale or scope-mismatched
- **THEN** exact retry and repreview remain unavailable and no replacement command is submitted.


### Requirement: Results and repair paths are actionable without invented values
The system SHALL display source and database writing rules, current scoped state, last committed evidence and scoped repair actions before optional technical detail. Partial, stale, failed and unknown outcomes SHALL remain distinct. General delivery evidence SHALL NOT be represented as raw database rows. Unsupported existing configurations SHALL retain read-only summaries and advanced access.

#### Scenario: One source is stale while the database is available
- **WHEN** current evidence identifies partial or stale source data
- **THEN** the interface identifies the affected source and opens its point or device repair context without labelling the database as the failing component.

#### Scenario: SQL writing fails after earlier success
- **WHEN** a recording has a previous receipt and a current database-write failure
- **THEN** the interface retains the historical timestamp, explains the current database issue and offers the scoped destination or logs action without displaying sample rows as actual results.

#### Scenario: Advanced-only configuration is inspected
- **WHEN** the simple view cannot represent an existing configuration
- **THEN** its rules are summarized without changing them and its advanced locator remains reachable.
