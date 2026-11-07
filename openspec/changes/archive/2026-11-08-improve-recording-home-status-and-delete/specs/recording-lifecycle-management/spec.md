## MODIFIED Requirements

### Requirement: Stop and remove preserve accepted data

Stopping SHALL prevent new intake while preserving immutable versions, accepted samples, outbox rows, destination identity, and receipts. Stop SHALL remain available using the group's saved connector revision after a live connector is edited or removed.

Removing a recording SHALL permanently delete, within one configuration-database transaction, every internal row keyed by any generation of the recording's lineage: group, versions, members, revision ends, canonical and profile keys, storage scopes and events, cutovers, ownership fences, migration maps, partition-maintenance rows, generations, the lineage, runtime versions, delivery samples, checkpoints, buckets, outbox rows and receipts, and recording-plan or measurement rows owned only by that recording. Rows whose key is shared with other recordings (storage scopes and partition-maintenance rows keyed by workspace, connector, database, schema, mode and layout) SHALL be deleted only when no surviving recording's current, historical or runtime version references them. Rows with no provable single-recording owner (measurement definitions, recording plans, the delivery journal) and terminal workspace-level `recording_start` ledger entries SHALL be retained, and their readers SHALL tolerate a missing group without re-enabling it. Removal MUST NOT connect to the destination database and MUST NOT delete or alter destination tables or recorded SQL rows. Removal SHALL be refused without any change when intake is active (`WRITE_GROUP_DELETE_ACTIVE`), when any generation has an unconsumed sample or an outbox row whose state is neither `sql_committed` nor `operator_skipped` (`WRITE_GROUP_DELETE_BACKLOG`, with per-category counts), when a reserved successor, pending cutover, unfinished cross-recording handoff, an unfinished or partially executed schema or test-write operation, an unresolved test-write cleanup, or a non-terminal `recording_start` operation covering the lineage exists (`WRITE_GROUP_DELETE_OPERATION_PENDING`), when the requested group is a predecessor of a non-deleted later generation or a finished cutover, fence or handoff of the lineage is still ownership evidence for a surviving recording (`WRITE_GROUP_DELETE_HAS_SUCCESSOR`, listing up to 20 referencing recordings and their total), or when expected revisions are stale. A previously tombstoned recording SHALL be permanently deletable through the same operation. A successful removal SHALL record a `write_group_deleted` audit entry containing only the recording name, group and lineage identifiers and generation count (a legacy single-generation recording without a lineage reports an empty lineage identifier and one generation), and SHALL release the canonical key so a new recording for the same device and role can be created through the ordinary flow.

#### Scenario: Stop with queued deliveries

- **WHEN** an operator stops a group with queued or retrying rows
- **THEN** no new intake is accepted and existing workers continue delivery to each row's frozen destination without relabeling it as new-generation data.

##### Example: Stop with queued deliveries

- **GIVEN** G1 has queued row Q1 bound to R1 and destination D1
- **WHEN** stop commits for G1
- **THEN** new intake is rejected and Q1 remains eligible for delivery to D1

#### Scenario: Remove a drained recording

- **WHEN** an operator removes a stopped recording whose samples are consumed and whose outbox rows are all `sql_committed` or `operator_skipped`
- **THEN** no internal row keyed by any of its group identifiers remains, the destination table and its rows are unchanged, and a `write_group_deleted` audit entry exists

##### Example: Remove a drained two-generation recording

- **GIVEN** lineage L1 has tombstoned G1 (generation 1) and stopped head G2 (generation 2), destination table `public.gw_long_a` with 1200 rows, all outbox rows `sql_committed`
- **WHEN** remove G2 is confirmed with current revisions
- **THEN** the response reports 2 deleted generations, no internal table contains G1, G2 or L1, and `public.gw_long_a` still holds 1200 rows

#### Scenario: Removal refused by pending delivery

- **WHEN** any generation has unconsumed samples or undelivered outbox rows
- **THEN** removal returns `WRITE_GROUP_DELETE_BACKLOG` with separate counts per category and changes nothing

##### Example: pending categories

| Pending evidence | Result |
| --- | --- |
| 3 samples with consumed = 0 | 409 `WRITE_GROUP_DELETE_BACKLOG`, samples_unconsumed = 3 |
| 1 outbox row in `retrying`, 1 in `unknown` | 409 `WRITE_GROUP_DELETE_BACKLOG`, retrying = 1, unknown = 1 |
| only `sql_committed` and `operator_skipped` rows | removal succeeds |

#### Scenario: Removal refused while recording

- **WHEN** removal is requested for a group whose intake is active
- **THEN** it returns `WRITE_GROUP_DELETE_ACTIVE` with action `stop_first` and changes nothing

#### Scenario: Purge an existing tombstone

- **WHEN** removal is requested for an already tombstoned recording with no backlog and no later generation
- **THEN** it is permanently deleted instead of returning a lifecycle-blocked error

#### Scenario: Ownership proof fails during removal

- **WHEN** an accepted sample or undelivered row cannot be resolved to its immutable owner
- **THEN** removal fails without changing the group, backlog, or historical SQL data.

##### Example: Ownership proof fails during removal

- **GIVEN** pending Q1 names a missing immutable revision R0
- **WHEN** remove G1 is requested
- **THEN** removal fails and G1 plus Q1 remain unchanged

#### Scenario: Removal refused by another recording's ownership evidence

- **WHEN** a finished cutover or handoff of the recording is still read by a surviving recording's ownership proof
- **THEN** removal returns `WRITE_GROUP_DELETE_HAS_SUCCESSOR` with action `open_recording` and the referencing recordings, and changes nothing

##### Example: cross-lineage cutover evidence

- **GIVEN** recording A's lineage handed over to surviving recording B through a finished cutover whose proof reads A's versions
- **WHEN** removal of A is requested
- **THEN** the response is 409 `WRITE_GROUP_DELETE_HAS_SUCCESSOR` with `referenced_by` naming B and total 1, and A's group, versions and cutover rows are unchanged

#### Scenario: Shared destination scope survives removal

- **WHEN** a removed recording's storage scope key is still referenced by a surviving recording's version
- **THEN** that scope and its partition-maintenance rows remain and every row owned only by the removed lineage is deleted

#### Scenario: Schema operation state decides removal

- **WHEN** removal is requested for a recording with a related schema or test-write operation
- **THEN** a failed operation with zero executed statements and no pending cleanup is deleted with the recording, and any pending, running, partial, unknown, partially executed or cleanup-pending operation refuses removal with `WRITE_GROUP_DELETE_OPERATION_PENDING`

##### Example: schema operation states

| Operation evidence | Result |
| --- | --- |
| failed, executed_statements = 0, no cleanup pending | removal succeeds |
| failed, executed_statements = 2 | 409 `WRITE_GROUP_DELETE_OPERATION_PENDING` |
| running | 409 `WRITE_GROUP_DELETE_OPERATION_PENDING` |
| succeeded test write with cleanup not completed | 409 `WRITE_GROUP_DELETE_OPERATION_PENDING` |

#### Scenario: Recording start ledger after removal

- **WHEN** a recording whose terminal `recording_start` ledger entry names it is removed
- **THEN** the ledger entry is retained, operation reads still return its historical result, replaying a succeeded start does not re-enable the recording, and replaying a failed start is refused because the group no longer exists

#### Scenario: Legacy recording without lineage

- **WHEN** a legacy single-generation recording without a lineage is removed
- **THEN** the success data reports `lineage_id` as an empty string and `deleted_generations` as 1

### Requirement: Historical destinations and pending units remain explicit

The lifecycle review SHALL list predecessor and successor destination identities, affected sources, generations, boundary, and pending evidence grouped by immutable revision. Samples, outbox rows, and blocked or unknown delivery states SHALL be separately labeled and SHALL NOT be summed as a single record count. Review, draft saved, schema confirmed, intake enabled, SQL committed, readback verified, and cleanup completed SHALL remain distinct facts.

#### Scenario: Remove or transition with mixed pending evidence

- **WHEN** a group contains collecting samples and queued or unknown outbox rows from older revisions
- **THEN** a transition confirmation displays their distinct units and original destinations and does not imply that transition retargets finalized rows, and a removal attempt is refused with the same distinct units so the operator can stop and wait for delivery; any unfinalized observation handoff is separately identified with its original provenance and consuming owner.

##### Example: Remove or transition with mixed pending evidence

- **GIVEN** R1 has 3 collecting samples and 2 queued outbox rows including 1 unknown outcome
- **WHEN** the operator opens impact review or requests removal
- **THEN** the counts remain separately labeled, the UI does not claim 5 homogeneous records, and removal is refused
