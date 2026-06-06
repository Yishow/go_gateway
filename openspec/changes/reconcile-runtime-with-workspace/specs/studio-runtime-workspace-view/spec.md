## ADDED Requirements

### Requirement: Workspace runtime view exposes projection alignment state

The workspace-scoped runtime view SHALL expose whether the current runtime projection is aligned with the latest persisted workspace state.

#### Scenario: Runtime view shows projection drift

- **WHEN** runtime is still running an older projection than the current persisted workspace state
- **THEN** the workspace runtime view reports that the projection is out of date
- **AND** the operator can see whether the next step is targeted reconcile, restart, or reactivation

##### Example: view shows stale projection after saved mapping change

- **GIVEN** mapping m-2 was persisted after runtime last reconciled projection version v12
- **WHEN** the workspace runtime view loads against still-running projection v12
- **THEN** the view reports drift between persisted state and active runtime projection
