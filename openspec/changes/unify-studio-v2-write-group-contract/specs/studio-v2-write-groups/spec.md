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
The system SHALL migrate simple mappings, workspace row groups and exactly representable basic recording plans in separately verifiable batches with durable provenance and stable identity maps. It MUST preserve source scope, row layout and legal shared-column behavior, and MUST block ambiguous or advanced modes instead of guessing. Repeating a migration for the same source revision SHALL be idempotent.

#### Scenario: Shared column layout is not one wide row
- **WHEN** legacy members share a column but require different row identities
- **THEN** migration preserves those identities or blocks review; it never collapses them into competing values in one row.

#### Scenario: Unsupported legacy plan
- **WHEN** a plan contains derived streams or cannot be mapped to one target safely
- **THEN** its original persisted data remains readable and unchanged, with explicit repair or advanced-mode status.

#### Scenario: Repeated reviewed migration
- **WHEN** the same source revision is reviewed and applied again
- **THEN** the same canonical group IDs are returned without duplicate outputs.

### Requirement: One production writer owner per output
The system MUST switch writer ownership only through the existing activation barrier after migration and consumer verification. Legacy APIs SHALL translate compatible writes through canonical CAS or return an actionable conflict; legacy reads SHALL project the canonical state for migrated scopes. Unmigrated scopes SHALL retain their existing path and capability label.

#### Scenario: Old client edits migrated settings
- **WHEN** a legacy request changes a migrated target
- **THEN** it updates the canonical group transactionally or is rejected, never creating a second writable authority.

#### Scenario: Activation switches writer
- **WHEN** a migrated group is applied while the legacy target exists
- **THEN** exactly one writer path owns new intake for that output and Share remains independently gated.

#### Scenario: Rollback with accepted backlog
- **WHEN** the new UI or binary is rolled back after durable records were accepted
- **THEN** new intake can be stopped but accepted records and their destination provenance are preserved.

### Requirement: Safe basic group lifecycle
The system SHALL support group rename, member and policy edits, disable and logical deletion under CAS. Rename MUST preserve the stable group and existing record/receipt identities without restarting the applied writer. Member/column/policy/destination changes SHALL save a draft and take effect only on explicit valid apply at the next bucket boundary. Disable SHALL stop new intake while retaining accepted backlog. Delete MUST tombstone the group, stop new intake and preserve immutable payloads, revisions, ownership, receipts and operation/backlog lookup until every accepted record is safely resolved; physical purge is outside this change.

#### Scenario: Rename and member edits
- **WHEN** a group is renamed or members are added or removed while it is running
- **THEN** rename keeps identity stable, semantic edits remain draft until apply, the old bucket closes under its original snapshot, and the new membership starts at the next bucket boundary.

#### Scenario: Delete with backlog
- **WHEN** the operator deletes a group with accepted pending records
- **THEN** new intake stops, the tombstone retains queryable ownership and immutable delivery evidence, and workers can complete the original accepted records without orphaning, replaying or dropping them; if the legacy format cannot preserve that ownership, deletion is blocked with a disable alternative.
