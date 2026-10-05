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

#### Scenario: Source mutation races save
- **WHEN** a workspace save overlaps source edit, derived sync, candidate reapply or deletion
- **THEN** operations complete without deadlock and a genuinely stale source proposal cannot overwrite the accepted pipeline

#### Scenario: Another rule progresses
- **WHEN** one rule mutation is waiting and another rule is saved
- **THEN** the second rule can progress without acquiring the first rule's guard

### Requirement: Rule queued autosave retains latest drafts and actionable failures

The Step 3 frontend SHALL queue mapping saves by source rule, preserve newer local edits while an older request is in flight, and allow explicit retry after failure. It MUST retain safe error status/code/request identity and MUST NOT automatically replay a failed request or expose backend exceptions.

#### Scenario: Repeated batch and rapid edits
- **WHEN** the operator repeats batch enabled/transform actions and rapidly edits a row
- **THEN** one request per rule runs at a time and the last local draft is eventually persisted without stale response overwrite

#### Scenario: Failure and explicit retry
- **WHEN** a save receives conflict, missing row or internal failure
- **THEN** the draft remains visible with a safe actionable status/code and explicit retry; no blind automatic replay occurs

#### Scenario: Default numeric pipeline remains exact
- **WHEN** the default uint64 mapping or a scaled numeric mapping is saved through the queue
- **THEN** default 9007199254740993 remains exact and scale 0.5 plus offset 10 applied to int16 243 remains 131.5 using float64

