## ADDED Requirements

### Requirement: Concurrent owned mapping saves preserve rule integrity

The system SHALL coordinate same-rule workspace mapping saves and source mutations while atomically replacing rule links. Different rules SHALL be able to progress independently. Persisted source revision, accepted hash and exact-pipeline confirmation guards MUST remain enforced.

#### Scenario: Eight same-rule rows save concurrently
- **GIVEN** eight persisted int16 mappings for addresses 40001 through 40008 with scale 1 and offset 0
- **WHEN** their workspace PUT requests run concurrently using real SQL repositories with connection pools 1 and 4
- **THEN** every request succeeds and all eight point/tag/mapping links, pipelines and accepted hashes remain consistent

#### Scenario: Atomic replacement fails
- **WHEN** replacement link insertion violates a unique constraint or fails after deleting old links inside the transaction
- **THEN** the request fails and readers retain the complete previous links without an empty intermediate view or false success

##### Example: Second insert fails
- **GIVEN** diag-rule has complete links for 40001 and 40002
- **WHEN** a trigger rejects the replacement insert for 40002 after the transaction deleted its old links
- **THEN** the save returns 500, both previous links and accepted mapping fields remain, and a concurrent reader sees the previous complete set

#### Scenario: Source mutation races save
- **WHEN** a workspace save overlaps source edit, derived sync, candidate reapply or deletion
- **THEN** operations complete without deadlock and a genuinely stale source proposal cannot overwrite the accepted pipeline

##### Example: Source scale changes during a save
- **GIVEN** a PUT for diag-rule/40001 is paused inside link replacement
- **WHEN** an update of the rule scale multiplier from 1 to 2 is requested and the paused PUT is released
- **THEN** the PUT and source update complete in order, and a subsequent old workspace save returns 409 workspace_mapping_conflict

#### Scenario: Another rule progresses
- **WHEN** one rule mutation is waiting and another rule is saved
- **THEN** the second rule can progress without acquiring the first rule's guard

##### Example: Different rule keeps saving and deleting
- **GIVEN** diag-rule is held by its mutation guard and other-rule owns address 40101
- **WHEN** other-rule receives a PUT and then DELETE
- **THEN** both return 200 before the diag-rule guard is released

### Requirement: Rule queued autosave retains latest drafts and actionable failures

The Step 3 frontend SHALL queue mapping saves by source rule, preserve newer local edits while an older request is in flight, and allow explicit retry after failure. It MUST retain safe error status/code/request identity and MUST NOT automatically replay a failed request or expose backend exceptions.

#### Scenario: Repeated batch and rapid edits
- **WHEN** the operator repeats batch enabled/transform actions and rapidly edits a row
- **THEN** one request per rule runs at a time and the last local draft is eventually persisted without stale response overwrite

##### Example: Eight rows and latest display name
- **GIVEN** eight int16 rows at 40001 through 40008 with first raw value 243
- **WHEN** enable-all is repeated three batches, float64/scale0.5/offset10 is applied to all, and the first display name changes First edit → Second edit → Final rapid edit
- **THEN** all eight remain linked to their original IDs and Saved, and SQL stores Final rapid edit with preview value 131.5

#### Scenario: Failure and explicit retry
- **WHEN** a save receives conflict, missing row or internal failure
- **THEN** the draft remains visible with a safe actionable status/code and explicit retry; no blind automatic replay occurs

##### Example: Owned database failure and operator retry
- **GIVEN** the first mapping's UPDATE is rejected by an owned test trigger
- **WHEN** the operator edits its display name to Recovered draft
- **THEN** the draft remains visible with HTTP500/workspace_mapping_save_failed and a bounded request ID, no SQL exception appears, and no automatic replay occurs
- **WHEN** the trigger is removed and the operator selects Retry
- **THEN** the same mapping identity is Saved and SQL stores Recovered draft

#### Scenario: Default numeric pipeline remains exact
- **WHEN** the default uint64 mapping or a scaled numeric mapping is saved through the queue
- **THEN** default 9007199254740993 remains exact and scale 0.5 plus offset 10 applied to int16 243 remains 131.5 using float64

