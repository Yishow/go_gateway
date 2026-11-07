## MODIFIED Requirements

### Requirement: Recording actions follow matching server availability
The system SHALL derive saved-recording actions from matching workspace, group and revision evidence. Stop and delete SHALL require explicit scope and impact confirmation. Deleting a recording SHALL permanently delete its internal configuration and delivery records and SHALL NOT delete historical database rows or tables in the destination database.

#### Scenario: The action facts are stale
- **WHEN** lifecycle evidence does not match the displayed recording revision
- **THEN** the interface disables consequential actions, explains the need to reload evidence and does not guess availability from status.

#### Scenario: A recording is deleted
- **WHEN** an operator acknowledges that deletion cannot be undone and confirms an allowed deletion
- **THEN** only the recording deletion contract is invoked, the interface states that the recording configuration was permanently deleted while the destination table and its rows were retained, and the recording no longer appears in the home or editor lists

#### Scenario: Deletion is refused by the server
- **WHEN** the deletion contract returns `WRITE_GROUP_DELETE_ACTIVE`, `WRITE_GROUP_DELETE_BACKLOG`, `WRITE_GROUP_DELETE_OPERATION_PENDING` or `WRITE_GROUP_DELETE_HAS_SUCCESSOR`
- **THEN** the interface keeps the recording, shows an actionable explanation with the matching stop, open-recording or refresh action, and does not display raw backend text

#### Scenario: A referenced source is changed
- **WHEN** an operator attempts to delete or redefine a source used by existing recordings
- **THEN** the interface preserves existing dependency protections and identifies affected recordings or opens the scoped advanced repair flow without creating dangling references.
