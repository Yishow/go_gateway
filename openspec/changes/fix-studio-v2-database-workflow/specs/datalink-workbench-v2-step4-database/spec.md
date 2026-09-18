## MODIFIED Requirements

### Requirement: Connector configuration form
The Step 4 connector section SHALL present saved connections or a new connection with only relevant fields, verified kind capabilities and explicit connection-test state. Managed recording SHALL generate storage context from confirmed plan intent; custom-table mode SHALL retain table and write-policy controls in an advanced section. The form MUST NOT use example endpoints, tables or ready status as real configured state. Equipment scope, selected plan, persisted connector and save status SHALL be visible before dependent preparation actions.

#### Scenario: Default connector on first render
- **WHEN** Step 4 mounts without a persisted connector
- **THEN** it shows an unconfigured draft with clearly labelled placeholders and supported-kind choices
- **AND** it does not claim connection or write readiness.

#### Scenario: Kind switch triggers auto-assign
- **WHEN** the operator changes connector kind
- **THEN** stale tests and schema previews are invalidated and actual metadata is reloaded
- **AND** bindings remain only when valid for the new identity; missing or ambiguous matches require confirmation instead of sample-column reassignment.

#### Scenario: Write strategy persistence
- **WHEN** an operator selects a compatible recording strategy
- **THEN** the plan persists and explains whether it appends history or updates latest state
- **AND** an update strategy without a valid unique identity is blocked.

#### Scenario: Unsaved destination
- **WHEN** the destination configuration has not been saved successfully
- **THEN** dependent create, apply and test-write actions are blocked with an explanation, not a default target fallback.

### Requirement: Tag-to-column auto-assignment
The system SHALL provide confirmed assignments, reviewable suggestions and unmatched results using the actual selected schema and compatible measurement semantics. Existing confirmed targets MUST be revalidated. Name similarity alone SHALL NOT confirm uncertain semantics; index fallback, wraparound reuse and sample columns MUST NOT produce production assignments. Metadata state SHALL distinguish not queried, loading, found, missing and unavailable; proposed new-table columns MUST be identified as proposals rather than observed metadata.

#### Scenario: Eight points to nine columns
- **GIVEN** eight confirmed measurements and sufficient real compatible columns
- **WHEN** matching runs
- **THEN** valid unambiguous assignments and reviewable suggestions are shown without assigning any column twice within a row; intentional reuse across distinct rows in a valid persisted row group remains legal.

#### Scenario: Existing target preserved
- **WHEN** an existing confirmed target remains compatible with the current connector, table, column and measurement definition
- **THEN** it is preserved; otherwise it is marked for repair rather than silently rebound.

#### Scenario: Fewer columns than points causes wrap and conflict
- **GIVEN** eight enabled ungrouped measurements requiring distinct columns and only four compatible columns
- **WHEN** matching runs
- **THEN** extra measurements remain unmatched and the UI offers create-column-plan, different-table or explicit exclusion options
- **AND** no wraparound or duplicate assignment is generated.

#### Scenario: Schema lookup fails
- **WHEN** metadata cannot be read due to permissions, connection failure or an invalid response
- **THEN** the UI explains the unavailable state and offers retry without substituting sample columns or claiming the table is absent.

#### Scenario: New table is proposed
- **WHEN** the operator selects managed table creation rather than an existing table
- **THEN** generated columns are labelled as a proposal requiring confirmation and are not shown as already present.

### Requirement: Column conflict detection
The system SHALL validate repeated column bindings within the persisted connector and table scope using the workspace-database-row-groups identity contract. Reuse inside one valid row group SHALL remain legal. Reuse across row groups, outside a valid group or under an unsafe upsert identity SHALL block submission. Invalid rows SHALL show a visible and accessible explanation identifying the affected bindings, not color alone.

#### Scenario: Two enabled ungrouped targets collide
- **GIVEN** two enabled targets share a column without a valid shared row-group identity
- **WHEN** validation runs
- **THEN** both affected rows show a conflict explanation and submission is blocked.

#### Scenario: Disabling one side resolves conflict
- **WHEN** one conflicting target is disabled and no other blocker remains
- **THEN** the column conflict clears and submission is enabled.

#### Scenario: Legal shared-column group
- **WHEN** repeated bindings share one valid row group with the required row identity
- **THEN** they are not rejected by a generic global duplicate-column check.

#### Scenario: Cross-group or unsafe upsert reuse
- **WHEN** repeated bindings belong to different row groups or an upsert lacks the required stable identity
- **THEN** the specific invalid reuse is reported and cannot be applied.

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

### Requirement: Guided database setup and accessible controls
Step 4 SHALL guide the operator through destination, existing-versus-managed table preparation, bindings or managed-plan review, and readiness. Managed and custom-table paths MUST be explicit rather than silently activating duplicate outputs. State and errors SHALL use the configured language, name a corrective action and distinguish loading failure from empty data. Controls MUST be labelled, keyboard reachable and usable without relying solely on color.

#### Scenario: New operator follows the setup
- **WHEN** no device, saved connection, table or plan has been selected
- **THEN** the relevant empty state explains the missing prerequisite before enabling dependent actions.

#### Scenario: Managed versus custom table
- **WHEN** the operator selects managed recording
- **THEN** the UI reviews members and planned structure rather than requiring unrelated manual column bindings
- **AND** custom-table mode presents actual metadata and explicit bindings without creating a second output implicitly.

#### Scenario: Narrow viewport and keyboard navigation
- **WHEN** the page is checked at 390, 768 and 1440 CSS-pixel widths using the keyboard
- **THEN** primary actions remain reachable, focus is visible, table overflow is contained in its own scrollable region and errors identify their controls.

#### Scenario: Filtering and bulk edits
- **WHEN** a search or problems-only filter is active and a bulk action is offered
- **THEN** the action states its affected scope and count and does not silently modify rows outside that scope.

#### Scenario: Independent evidence statuses
- **WHEN** configuration, schema, test-write, cleanup and device-runtime outcomes differ
- **THEN** the UI shows those results separately, keeps unavailable facts unconfirmed and never combines them into a fabricated overall success.
