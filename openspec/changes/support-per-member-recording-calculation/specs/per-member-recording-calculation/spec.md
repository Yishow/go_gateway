## Purpose

Allow independently configured latest, mean and delta results for recording members within one recording with individually aligned windows, while preserving legacy semantics and durable counter correctness.

## ADDED Requirements

### Requirement: Recording members own calculation intent
A versioned mixed recording profile SHALL persist explicit latest, mean or delta calculation on each canonical recording member, not on a global source Point or Tag. Reads SHALL expose configured and effective calculation, applicability, source/result types and applicable revisions. A new member SHALL initialize to latest; a group default MAY initialize new members but MUST NOT overwrite explicit choices. Unknown or omitted methods in the new mixed shape SHALL fail validation.

Independent mixed record identities SHALL permit source-reference reuse within the same workspace and destination database/schema without weakening target-table/column ownership, legacy circuit/profile create-once identity, migration source claims or lineage integrity. Mixed record scope SHALL be distinct from legacy canonical key scopes and SHALL retain idempotency, CAS and foreign-scope rejection; a source device or circuit/profile pair MUST NOT collapse all independent mixed records into one identity.

#### Scenario: One source participates in different records
- **WHEN** record A uses a source as latest every 5 seconds and independent record B in the same database/schema but another owned managed table uses that source as delta every 600 seconds with confirmed counter semantics
- **THEN** each record retains its own member intent and counter state, and editing either record does not change the source Point or the other record.

#### Scenario: Reusing a source does not transfer target ownership
- **WHEN** another mixed record references an existing source but attempts to claim a managed table or columns owned by the first record without an existing authorized shared-storage scope
- **THEN** validation rejects the target conflict without changing its owner, while a retry for the same record retains its canonical identity and lineage.

### Requirement: Member windows align independently of activation
The mixed profile SHALL persist each member's interval and align its windows to an explicit clock basis rather than counting from activation. It SHALL support the requested 5-second, 15-second, one-minute and ten-minute aligned cadences without applying the superseded shared 300-second cap. Existing UTC-epoch half-open-window alignment SHALL be reused without changing legacy profile interpretation or acquisition polling frequency. Latest, arithmetic mean and confirmed-counter delta SHALL operate on the matching member's own window. Data-window timestamps SHALL be distinct from actual SQL commit time; lateness, worker queues and retry MUST NOT be presented as punctual physical commits. A member is due at aligned boundary T and selects observations in [T-interval,T); an observation exactly at T belongs to the next window. Closure waits for the existing allowed lateness. Positive integer second/minute intervals SHALL be normalized with overflow checks and backend-advertised bounded capability validation; the required acceptance cadences SHALL be supported without inventing an unrestricted interval range.

#### Scenario: Two averages and one usage share a row
- **WHEN** one aligned boundary closes equal-length member windows containing gauge A readings 20,22,24 and gauge B readings 10,14, and a confirmed counter C advances reliably from baseline 1000 to 1008
- **THEN** the mixed wide output is one row with A=22, B=12 and C=8 under the same bucket and closure identity, not three separately scheduled records.

### Requirement: Mixed result columns are not legacy raw columns
The new member-calculation-v1 result shape SHALL store each member's selected result in its target column. Delta columns SHALL contain delta rather than raw counter readings; raw values, counter status, actual interval, coverage and provenance SHALL remain distinguishable evidence. Result types and exact SQL encoding SHALL be validated before writing. New mixed support SHALL cover existing managed wide and long layouts, with one due member per long row. Native-monthly or custom shapes without new-shape capability proof SHALL remain blocked rather than claiming one-row equivalence.

#### Scenario: A raw legacy table would change meaning
- **WHEN** an operator explicitly requests mixed delta results for sources used by a legacy cumulative recording
- **THEN** in-place Update and same-lineage successor conversion are rejected with typed errors because their canonical identities differ; the existing independent mixed-record creation flow requires a fresh RecordID and lineage and fresh schema, result-type, nullable/quality, provenance and readiness evidence, rejects legacy proofs, and leaves the old record, historical raw data and legacy create-once identity unchanged without copying baselines or automatically stopping or deleting the old record. Wide raw-to-result changes SHALL require a new owned target. Long output SHALL reuse a table only under existing shared-long ownership and row partitioning with fresh mixed result-shape, nullable, provenance and schema proofs, without reinterpreting existing rows. Otherwise it SHALL use an existing new-storage-scope/new-target flow or remain explicitly blocked, not invent a table-name bypass, alias or relaxed ownership guard. This is explicit creation of a separate mixed recording, not successful in-place conversion or splitting its members into separate recordings.

### Requirement: Missing results obey explicit completeness policy
A mixed member without a trustworthy computed result SHALL remain unavailable rather than becoming zero, raw or a different reduction. Optional unavailable results MAY be SQL NULL only with confirmed nullable and quality/provenance capability. Any due required unavailable member SHALL skip the wide row, or only its own long member row; a required member that is not due MUST NOT cause a skip. First-baseline, uncertain discontinuity and missing-window span outcomes SHALL retain specific reasons. The mixed delta column MUST NOT mislabel a cross-window span total as current-window usage.

#### Scenario: A required delta has no first baseline
- **WHEN** the first mixed wide window has valid gauge means but its required delta member only establishes a baseline
- **THEN** no partial result row is delivered, the skip reason is recorded, and the baseline plus skipped closure advance atomically so a later window can produce a delta.

#### Scenario: An optional delta is uncertain or spans missing windows
- **WHEN** the optional delta has no proven single-window result because of reset, unexplained decrease or a gap span
- **THEN** its column is NULL with a reason and actual interval evidence where available, other eligible columns may be written under the confirmed partial policy, and no raw value, zero or span total substitutes for the missing delta.

### Requirement: Counter policy and member state remain isolated
Only delta members SHALL require calculation baselines; latest counter_raw members SHALL retain raw semantics without generating usage. Delta SHALL reuse confirmed signed-counter, rollover, reset and source-binding policies. Counter checkpoints SHALL be scoped by workspace, group/revision, member identity and measurement binding/epoch. A changed method or binding MUST NOT reuse an unproven baseline, and copying a recording MUST NOT copy its baseline or proof.

#### Scenario: A policy explains a decrease
- **WHEN** a confirmed signed counter or a proven rollover produces a valid decrease-related result under the existing policy
- **THEN** the mixed reducer follows that policy rather than universally rejecting negative changes, while an unproven discontinuity remains unavailable with its reason.

### Requirement: Calculation changes invalidate configuration proofs
Methods, result shape and types, required policy, interval and counter identity SHALL participate in persisted candidates, managed digests, schema readiness, runtime descriptors and observation/proof versions. Updates SHALL retain workspace/member ownership, expected revisions and existing idempotency protections. Old samples and outputs MUST NOT be reinterpreted under a changed calculation, even when the source Point and target column names are unchanged.

#### Scenario: An old proof is submitted after changing method
- **WHEN** a member changes from latest to delta and an earlier readiness, observation association or schema proof is submitted
- **THEN** activation or binding is rejected as stale until matching fresh evidence and required confirmation are obtained, with no partial mutation or new writer.

### Requirement: Closure and counter progress are atomic and replayable
The existing durable transaction SHALL commit mixed output or skip decisions, checkpoints, bucket progress and output intents together. A crash or rollback MUST NOT advance a baseline without its corresponding closure decision. Restart and retry SHALL replay immutable accepted runtime/output versions and stable delivery identities rather than recalculate old records with new methods. Lost SQL responses SHALL use existing receipt resolution rather than creating duplicate rows.

#### Scenario: A crash occurs around mixed closure
- **WHEN** a crash occurs before or after the local closure transaction or after SQL commit before acknowledgment
- **THEN** recovery yields the same accepted calculation and at most one managed row effect, with matching checkpoints and receipt evidence rather than duplicated counter usage.

### Requirement: Legacy profiles keep their original meaning
Existing instantaneous and cumulative profiles without the new explicit mixed shape SHALL retain their saved effective reductions, intervals, counter policies, output representation and replay semantics through upgrade, restart and unrelated edits. No migration SHALL infer arbitrary per-member method intent or reinterpret legacy cumulative raw columns as delta results. Changing from legacy to mixed SHALL use explicit independent mixed-record creation with fresh identity and current schema and lifecycle safeguards; in-place or same-lineage cross-canonical conversion MUST remain rejected. Consumers SHALL distinguish creating or copying a definition into a new recording from converting the old recording, and SHALL keep activation and stopping explicit.

Storage expansion SHALL preserve nullable absence and all legacy payloads, revisions and relationships. Expanding the SQLite profile-kind constraint SHALL accept only recognized schema shapes, preserve complete rows and dependent schema objects, validate foreign keys before commit, roll back on failure and restore the connection foreign-key setting; a connection whose setting cannot be restored MUST NOT return to the pool. Existing initialized lineage corruption MUST fail rather than be guessed or repaired into new history.

#### Scenario: Expansion fails while replacing the configuration table
- **WHEN** a fault occurs during copy, table replacement, index restoration or foreign-key validation
- **THEN** the migration rolls back without changing legacy members, snapshots, mappings, lineage or revisions, and rerun does not infer calculation or cadence.

#### Scenario: Existing profiles are reloaded after upgrade
- **WHEN** a legacy mean profile with gauge/counter_raw/state roles and a legacy cumulative profile are reloaded or replayed
- **THEN** gauges retain their previous mean behavior, counter_raw/state retain latest behavior, cumulative raw-plus-provenance retains its prior representation, and neither profile silently becomes mixed.

### Requirement: Production proof precedes editor support
Completion SHALL include tests through actual scoped persistence, readiness, frozen runtime, mixed snapshot dispatch, typed row encoding, durable ledger and disposable SQL write/readback. Pure helper or mock UI success SHALL NOT establish production support. The simple editor SHALL consume the canonical member contract without a second calculation authority and show unsupported capabilities explicitly.

Storage expansion SHALL establish internal mixed-snapshot persistence and scoped readback without enabling normal mixed mutation or activation. Normal mixed draft CRUD SHALL require complete member cadence, applicability and record identity validation; it MUST NOT fabricate a group interval to satisfy a legacy validator. Saving a draft SHALL NOT establish readiness, runtime or editor capability before their existing production gates pass.

#### Scenario: Mixed storage exists before cadence validation
- **WHEN** the expansion phase can round-trip an internal mixed snapshot but normal member cadence validation is not implemented
- **THEN** normal mixed create/update and activation remain rejected, legacy CRUD/CAS stay unchanged, and foreign-scope reads remain rejected.

#### Scenario: Only a helper supports mixed methods
- **WHEN** a helper computes means and delta but the production writer or frozen proof still rejects or drops the member method
- **THEN** acceptance remains incomplete and the editor does not claim that the new combination is supported.

#### Scenario: Starting between clock boundaries
- **WHEN** recording starts at 12:00:07 with a 15-second member and a 5-second member
- **THEN** future boundaries remain on their declared clock-aligned phases rather than 12:00:22 or a new activation-relative schedule, and any partial startup window is identified rather than claimed complete.

#### Scenario: Mixed cadence and delayed delivery
- **WHEN** one member uses 15 seconds, another one minute and another ten minutes under UTC-epoch alignment displayed at UTC+08:00
- **THEN** their boundary previews show seconds 00/15/30/45, minute boundaries, and minutes 00/10/20/30/40/50 respectively, while actual SQL commit time is separately reported and cannot rewrite those data-window identities.

### Requirement: Different member cadences preserve sparse due semantics
Wide output SHALL use the union of member due boundaries. A not-due member SHALL be NULL with not_due quality and SHALL NOT carry a previous result, substitute zero, advance its baseline or trigger required-member failure. Long output SHALL emit rows only for due members, distinguishing due failure, no baseline, gap and reset from not_due. Each due result SHALL retain its own method, aligned window, actual coverage and provenance. Fresh schema proof SHALL establish nullable not-due columns and quality support even for members required when due. Durable closures SHALL retain the due set, per-member progress and immutable result versions; replay SHALL preserve wide row identity or long member/window identity without duplicate counter usage.

#### Scenario: Fifteen-second means and a minute delta
- **WHEN** mean A and B have 15-second intervals and delta C has a 60-second interval
- **THEN** wide boundaries :15/:30/:45 contain means and C=NULL with not_due, while :00 contains the preceding 15-second means and preceding 60-second delta; long output contains only the members due at each boundary.

#### Scenario: Activation truncates the first aligned window
- **WHEN** a 15-second member starts at 12:00:07 and first closes at 12:00:15 plus allowed lateness
- **THEN** its identity remains [12:00:00,12:00:15), actual coverage starts at 12:00:07, eligible latest/mean may retain a boundary-truncated partial result under existing safeguards, and a delta without a proven baseline remains unavailable with NULL or required skip and atomic baseline closure, not fabricated full-window usage.
