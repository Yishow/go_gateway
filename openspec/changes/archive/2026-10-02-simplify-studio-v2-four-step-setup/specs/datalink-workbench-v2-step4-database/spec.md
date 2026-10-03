## MODIFIED Requirements

### Requirement: Connector configuration form
The group editor SHALL use one canonical WriteGroup and saved destination identity for managed and custom storage. Managed readiness MUST NOT require a separate manually populated target list. If every existing plan is incompatible, the operator SHALL still be able to create a new group or explicitly review a migration; the system MUST distinguish load errors, empty results and mismatched scope.
The Step 4 connector section SHALL present saved connections or a new connection with only relevant fields, verified kind capabilities and explicit connection-test state. Managed recording SHALL generate storage context from confirmed group intent; custom-table mode SHALL retain table and write-policy controls in an advanced section. The form MUST NOT use example endpoints, tables or ready status as real configured state. Equipment scope, selected group (or explicitly chosen advanced plan), persisted connector and save status SHALL be visible before dependent preparation actions.

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
- **THEN** the group configuration persists and explains whether it appends history or updates latest state
- **AND** an update strategy without a valid unique identity is blocked.

#### Scenario: Unsaved destination
- **WHEN** the destination configuration has not been saved successfully
- **THEN** dependent create, apply and test-write actions are blocked with an explanation, not a default target fallback.

### Requirement: Tag-to-column auto-assignment
The system SHALL provide confirmed assignments, reviewable suggestions and unmatched results using the actual selected schema and compatible source value types; confirmed measurement semantics SHALL be required only for explicitly selected advanced operations. Existing confirmed targets MUST be revalidated. Name similarity alone SHALL NOT confirm uncertain semantics; index fallback, wraparound reuse and sample columns MUST NOT produce production assignments. Metadata state SHALL distinguish not queried, loading, found, missing and unavailable; proposed new-table columns MUST be identified as proposals rather than observed metadata.

#### Scenario: Eight points to nine columns
- **GIVEN** eight confirmed typed Tags and sufficient real compatible columns
- **WHEN** matching runs
- **THEN** valid unambiguous assignments and reviewable suggestions are shown without assigning any column twice within a row; intentional reuse across distinct rows in a valid persisted row group remains legal.

#### Scenario: Existing target preserved
- **WHEN** an existing confirmed target remains compatible with the current connector, table, column and source type (plus any explicitly selected measurement definition)
- **THEN** it is preserved; otherwise it is marked for repair rather than silently rebound.

#### Scenario: Fewer columns than points causes wrap and conflict
- **GIVEN** eight enabled ungrouped typed Tags requiring distinct columns and only four compatible columns
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
The system SHALL validate repeated column bindings within the persisted connector and table scope using the workspace-database-row-groups identity contract. Reuse inside one valid row group SHALL remain legal when values belong to distinct identified rows. Across canonical groups, reuse SHALL be allowed only when verified distinct entity/record identities prevent row collision; unverified legacy cross-group reuse, competing values in one row, or unsafe upsert identity SHALL block submission. Invalid rows SHALL show a visible and accessible explanation identifying the affected bindings, not color alone.

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
- **WHEN** repeated bindings cross row groups without proven distinct row identities, or an upsert lacks the required stable identity
- **THEN** the specific invalid reuse is reported and cannot be applied.

#### Scenario: Verified distinct groups reuse a business column
- **WHEN** two canonical groups share a table and business column but their verified entity/record keys identify distinct rows
- **THEN** they remain valid and do not overwrite one another, while competing values for the same row and column remain blocked.

## ADDED Requirements

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

### Requirement: Write-group completion follows persisted readiness
The UI SHALL derive submit eligibility and completion from backend group readiness and applied workspace/settings/group/connector revisions, preserving readiness tokens and Share gates. Saved draft, applied revision, collecting, local durable, SQL committed, verified and cleanup states MUST remain distinct. Share-only setup MUST NOT require database fields or start a database writer.

#### Scenario: Managed output without manual targets
- **WHEN** a managed group is ready but the old db.targets collection is empty
- **THEN** the UI does not block on enabledTargetCount from the unrelated model.

#### Scenario: Saved or collecting but database pending
- **WHEN** a group is saved or its device is running while target delivery is queued
- **THEN** the completion view shows the actual stage and does not claim SQL success.

#### Scenario: Keyboard and filter scope
- **WHEN** the operator filters rows and uses bulk actions at 390, 768 or 1440 CSS-pixel test viewports
- **THEN** focus and error repair remain reachable, the table scrolls independently, and the explicit affected scope/count prevents hidden edits outside the selected bulk scope.

#### Scenario: Share only
- **WHEN** only a valid Local Modbus output is selected
- **THEN** database validation is skipped while existing Share readiness barriers remain enforced.
