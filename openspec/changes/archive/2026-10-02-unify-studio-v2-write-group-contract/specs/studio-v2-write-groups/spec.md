## Purpose

Define one revision-safe database write-group authority, with explicit migration review and ownership-preserving lifecycle. Basic Tag writing remains independent of advanced measurement semantics.

## ADDED Requirements

### Requirement: One persisted write-group authority
The system SHALL use one versioned WriteGroup as the sole configuration authority for each migrated database output, containing persisted workspace/device/point/tag membership, connector identity revision, table and column bindings, row policy and applied revision. Managed and custom storage SHALL be strategies of that group, not independently writable parallel plans. Basic typed Tag writing MUST NOT require reporting, retention configuration, physical measurement semantics, aggregation or derived energy.

#### Scenario: Basic raw Tag without measurement semantics
- **WHEN** an operator saves a group of valid persisted Tags without MeasurementDefinitions
- **THEN** the group is valid for basic writing without fabricating measurement IDs, quantities or semantic kinds; explicitly selected derived operations remain gated on their own confirmed semantics.

#### Scenario: Foreign or missing source
- **WHEN** a member references an absent or foreign workspace/device/point/tag
- **THEN** the mutation returns a safe resource or validation error and saves no group or projection.

#### Scenario: Managed group has no manual targets
- **WHEN** a managed group has validated members and storage bindings
- **THEN** its readiness uses that group and does not require a second manually populated target-mapping collection.

### Requirement: Atomic revision-safe group saves
The system MUST persist a group, compatibility projections and workspace revision in one local transaction, gated by expected workspace/group/connector revisions. Saving SHALL create a draft; activation MUST retain existing settings revision and readiness-token checks. External DDL SHALL remain a separate confirmed operation.

#### Scenario: Second save stage fails
- **WHEN** projection persistence fails after the group write in the same save
- **THEN** all local changes roll back and reload returns the prior coherent revisions.

#### Scenario: Concurrent edits
- **WHEN** two clients save against the same group revision
- **THEN** only one change commits and the other receives a conflict without losing its draft.

#### Scenario: Save differs from apply
- **WHEN** a running group is edited and saved
- **THEN** the applied revision remains unchanged until explicit validated activation, and accepted backlog is not reinterpreted.

### Requirement: Reviewable repeatable compatibility migration
The system SHALL migrate simple mappings through explicitly reviewed snapshot conversion, and workspace row groups in separately verifiable batches with durable provenance and stable identity maps. Recording plans without proven equivalent snapshot layout SHALL retain their persisted intent and expose a blocked migration preview with repair guidance; new basic writing SHALL use canonical WriteGroups. It MUST preserve source scope, row layout and legal shared-column behavior, and MUST block ambiguous or advanced modes instead of guessing. Repeating a migration for the same source revision SHALL be idempotent.

#### Scenario: Explicit single-mapping snapshot conversion
- **WHEN** a legacy single mapping writes every observed sample and the operator previews conversion
- **THEN** the preview lists every-sample to periodic snapshot, value selection, lateness, incomplete-row and timestamp differences; saving requires explicit confirmation bound to the current source revision and review digest, and preserves the legacy configuration and writer until valid Apply.

#### Scenario: Stale or unconfirmed conversion
- **WHEN** confirmation is absent, a source is blocked, or workspace, connector, source revision or preview digest changed
- **THEN** no canonical group, identity map or projection is saved; stale review returns conflict and unsupported or unconfirmed conversion returns validation error.

#### Scenario: Shared column layout is not one wide row
- **WHEN** legacy members share a column but require different row identities
- **THEN** migration preserves those identities or blocks review; it never collapses them into competing values in one row.


##### Example: distinct persisted buffer rows
- **GIVEN** row group `legacy-g` has points `p-a` and `p-b`, target column `value`, GroupKeys `legacy-g:p-a` and `legacy-g:p-b`, and unique-key metadata `[entity_id, observed_at]`
- **WHEN** the operator confirms a current preview
- **THEN** the saved members retain separate entity keys and the original key metadata; no point ID or GroupKey is fabricated as an SQL business value, the original target mappings remain unchanged, and no writer is activated.

#### Scenario: Row-group collision or missing identity
- **WHEN** shared-column targets lack unique-key metadata or a persisted GroupKey, or two members use the same entity and column, or destination, interval or source ownership is ambiguous
- **THEN** preview marks the row group blocked and review saves no group, identity map or projection.

#### Scenario: Reviewed row-group differences and compatibility reads
- **WHEN** a valid row-group preview is confirmed and its canonical draft is later edited
- **THEN** the review retains original row-group and target IDs, source provenance and key metadata, explicitly states snapshot selection and epoch-alignment differences, and legacy reads reflect each corresponding canonical member or return an actionable conflict when that representation is no longer possible.

#### Scenario: Unsupported legacy plan
- **WHEN** a plan lacks an equivalent periodic snapshot layout or contains derived streams, multiple targets or unresolved measurement sources
- **THEN** preview returns its complete original intent, source revision and explicit blocking reasons with repair guidance; no candidate is fabricated and the original persisted plan remains readable and unchanged.

##### Example: every-sample plan cannot become a snapshot implicitly
- **GIVEN** plan `legacy-plan-A` has one persisted measurement `measurement-A`, one `raw_history` / `every_sample` stream and destination `connector-A`, with no explicit snapshot column or row identity bindings
- **WHEN** the operator previews migration and submits confirmation using its current revisions and digest
- **THEN** preview is blocked with no candidate, review returns validation failure, the original plan payload and revision remain unchanged, and no group, identity map, projection, writer or external SQL operation is created.

#### Scenario: Recording-plan preview revision changes
- **WHEN** the plan, referenced measurement/source or saved destination changes after preview, or the source belongs to a foreign workspace
- **THEN** review rejects stale revisions or digest with conflict and foreign or unknown resources with the same safe not-found response; no persisted data changes.

#### Scenario: Repeated reviewed migration
- **WHEN** the same source revision is reviewed and applied again
- **THEN** the same canonical group IDs are returned without duplicate outputs, workspace revision increments or overwriting subsequent canonical edits. The durable map retains the reviewed source revision, original intent, adapter version and digest.

### Requirement: One production writer owner per output
The system MUST switch writer ownership only through the existing activation barrier after migration and consumer verification. Legacy APIs SHALL translate compatible writes through canonical CAS or return an actionable conflict; legacy reads SHALL project the canonical state for migrated scopes. Unmigrated scopes SHALL retain their existing path and capability label.

#### Scenario: Old client edits migrated settings
- **WHEN** a legacy request changes a migrated target
- **THEN** it updates the canonical group transactionally or is rejected, never creating a second writable authority.

##### Example: original and current target scopes remain owned
- **GIVEN** target mapping `legacy-A` for Tag `tag-A` and connector `connector-A` has migrated to group `group-A`, whose current destination is `connector-B`
- **WHEN** an old client updates or deletes `legacy-A`, or creates a target for `tag-A` at `connector-B` using a different table or column
- **THEN** the request returns HTTP409 with `WRITE_GROUP_LEGACY_WRITE_CONFLICT` and `open_write_groups`; legacy rows, canonical group, identity map and workspace revision remain unchanged.

#### Scenario: Legacy row-group replacement removes an owned group
- **WHEN** a current-revision Studio database-config request explicitly replaces row groups and removes or changes a canonical-owned row group
- **THEN** it returns an actionable conflict before connector or projection persistence; retaining an equal owned group or omitting replacement preserves the existing path.

##### Example: explicit empty replacement cannot erase migrated intent
- **GIVEN** workspace row group `legacy-g` belongs to canonical `group-A` and the request has the current setup revision
- **WHEN** the client sends `row_groups: []`
- **THEN** HTTP409 preserves the connector, workspace row group, target references, canonical group and identity map.

#### Scenario: Unmigrated legacy target writes remain available
- **WHEN** a target write uses neither an owned legacy mapping identity nor a canonical-owned Tag/connector scope
- **THEN** existing CRUD remains available without creating a canonical group or activating a new writer.

#### Scenario: Prepared owner transition shares the Apply transaction
- **WHEN** an installed transaction-scoped owner barrier validates a new semantic Apply, rejects it, or the later workspace save fails
- **THEN** immutable version, applied identity, owner projection and workspace revision commit together or all roll back; stale local revisions or invalid sources do not invoke the owner transition, and a metadata-only rename does not restart it.

##### Example: owner write rolls back with the workspace save
- **GIVEN** group `group-A` has draft revision `group-revision-1`, a 15-second interval, clock `2026-10-02T12:00:07Z`, and the test owner projection is `legacy`
- **WHEN** the injected barrier writes the proposed owner for effective time `2026-10-02T12:00:15Z` and the subsequent workspace save fails
- **THEN** reload retains owner `legacy`, an empty applied revision, no immutable version and the prior workspace revision; original legacy targets and accepted payload/receipt identities remain unchanged.

#### Scenario: Preparation keeps production consumers unchanged
- **WHEN** only the domain Apply preparation seam has shipped and the production consumer has not been wired
- **THEN** no new runtime or writer is enabled; applied metadata is not treated as production owner proof, and the existing Share settings, hydration, readiness token and independent activation checks remain in effect.

### Requirement: Safe basic group lifecycle
The system SHALL support group rename, member and policy edits, disable and logical deletion under CAS. Rename MUST preserve the stable group and existing record/receipt identities without restarting the applied writer. Member/column/policy/destination changes SHALL save a draft and take effect only on explicit valid apply at the next bucket boundary. Disable SHALL stop new intake while retaining accepted backlog. Delete MUST tombstone the group, stop new intake and preserve immutable payloads, revisions, ownership, receipts and operation/backlog lookup until every accepted record is safely resolved; physical purge is outside this change.

#### Scenario: Rename and member edits
- **WHEN** a group is renamed or members are added or removed while it is running
- **THEN** rename keeps identity stable, semantic edits remain draft until apply, the old bucket closes under its original snapshot, and the new membership starts at the next bucket boundary.

#### Scenario: Delete with backlog
- **WHEN** the operator deletes a group with accepted pending records
- **THEN** new intake stops and the tombstone retains queryable ownership and immutable delivery evidence; if the legacy format cannot preserve that ownership, deletion is blocked with a disable alternative. Completing the accepted records through a production worker is owned by the durable delivery change.
