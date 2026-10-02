## Purpose

Preserve typed acquisition facts through basic write-group snapshots, with deterministic UTC closure, explicit quality and stable scoped row identity. Durable acceptance and destination delivery are owned by the dependent delivery change.

## ADDED Requirements

### Requirement: Typed acquisition facts reach the writer
The production sample path SHALL preserve sample/acquisition identity, workspace/device/point/tag identity, source and mapping revisions, observed_at, received_at, time_origin, exact typed value, quality and reason through mapping and group assembly. When a device provides no trusted timestamp, observed_at MUST be the actual gateway acquisition time with that origin, not a later flush time.


#### Scenario: Heterogeneous values
- **WHEN** bool, text, decimal and uint64 9007199254740993 pass through API and group encoding
- **THEN** their types and exact values remain intact and no JavaScript Number rounding occurs.

##### Example: exact JSON and SQLite value roundtrip
- **GIVEN** values bool `true`, text `batch-001`, int64 `-9223372036854775808`, uint64 `9007199254740993` and decimal `1234567890.123456789012345678`, with SQLite INTEGER columns for bool/signed/representable unsigned values and an explicit TEXT exact strategy for decimal
- **WHEN** the opt-in exact codec encodes and decodes JSON, binds parameters into a disposable SQLite table and decodes actual SELECT results using each expected type
- **THEN** every type and exact value remains intact, uint64 `18446744073709551615` uses an explicit TEXT strategy, and INTEGER overflow, decimal NUMERIC/REAL storage, PostgreSQL NUMERIC(16,4) rounding, NaN, infinity and malformed encoding are rejected without a substitute value.

#### Scenario: Failed or non-finite read
- **WHEN** a read is missing, bad, NaN, infinite or cannot fit the verified SQL type
- **THEN** it is excluded as a good value and retains an explicit missing/bad/type reason, never zero or an old reading.

#### Scenario: Identity collision
- **WHEN** two devices expose the same address or tag display name
- **THEN** their persisted IDs keep the samples distinct.

##### Example: preserve gateway acquisition time and exact identity
- **GIVEN** workspace `workspace-A`, devices `device-A` and `device-B` both use address `40001`, mapped points `point-A`/`point-B` and Tags `tag-A`/`tag-B` share a display name, acquisition `acquisition-A` completes at `2026-01-01T00:00:08Z`, and its exact value is uint64 `9007199254740993` with no trusted source timestamp
- **WHEN** runtime processes that acquisition at `2026-01-01T00:00:20Z` and sends it to the optional typed sample sink
- **THEN** observed_at and received_at remain `2026-01-01T00:00:08Z` with gateway acquisition origin, persisted IDs and source/mapping revisions remain distinct, repeated acquisition yields the same sample ID, and no measurement ID, zero value, current flush time or credential is fabricated.

### Requirement: Deterministic UTC snapshot closure
The system SHALL create basic snapshots in UTC half-open intervals [start,end), selected by observed_at rather than arrival. Each member SHALL choose the greatest observed_at within the bucket, using lexicographically greatest stable sample_id to break a time tie. A bucket SHALL close once at end plus the persisted nonnegative allowed_lateness_seconds, default zero; samples arriving after closure MUST NOT reopen it. Applied groups MUST close buckets from a time-driven scheduler even if no sample arrived, emitting one scoped no_data/skipped outcome for a silent bucket. Basic snapshots MUST NOT be described as interval averages, usage or every-sample history.

#### Scenario: Out-of-order arrival
- **WHEN** a sample observed at second 8 arrives before one observed at second 3 in the same ten-second fixture bucket
- **THEN** the second-8 sample remains selected regardless of arrival order.

#### Scenario: Boundary and late sample
- **WHEN** a sample is observed exactly at bucket end or arrives after the previous bucket is closed
- **THEN** the boundary sample belongs to the next bucket and the late sample is reported without rewriting the closed row.

#### Scenario: Repeated sample ID
- **WHEN** the same sample ID is received again
- **THEN** identical content is a no-op and different content is a safe identity conflict.

#### Scenario: Future clock skew
- **WHEN** observed_at exceeds the configured maximum future clock skew
- **THEN** the sample is invalid with a reason and cannot keep a bucket open indefinitely.

### Requirement: Explicit freshness and completeness
The default group policy SHALL require all enabled members, skip a row if a required member is missing/bad/stale, and record the bucket and cause. Freshness SHALL compare end minus observed_at against persisted max_age_seconds, defaulting to the interval, without carrying any value from another bucket. Explicit partial mode SHALL be accepted only when verified storage supports nullable values and member-level observation/quality metadata.

#### Scenario: Missing member
- **WHEN** temperature is valid but pressure has no good value in the bucket
- **THEN** the default group emits no SQL row and reports the incomplete bucket; it does not claim complete delivery.

#### Scenario: Partial mode
- **WHEN** a supported group explicitly uses partial mode and one member is missing
- **THEN** the missing member is SQL NULL with its reason and observation metadata, not zero or a stale value.

#### Scenario: Unsupported partial target
- **WHEN** a custom table cannot store NULL or required quality/provenance metadata
- **THEN** partial configuration is rejected before activation.

### Requirement: Scoped stable row identity
The system SHALL derive a stable record ID from workspace, group ID, group revision, entity key and bucket start, and scope destination effects by the frozen destination. Managed storage MUST enforce this identity. Per-member provenance for all-good complete rows MAY remain in the local durable envelope linked to record and destination identity; such custom tables MUST NOT be forced to add external per-field metadata. Explicit partial mode or the chosen deduplication strategy SHALL still require its verified storage capabilities. Custom modes MUST verify exact typed storage from real metadata. Upsert and automatic retry after uncertain completion SHALL require proven destination uniqueness; timestamp alone SHALL NOT be sufficient across distinct groups/entities. Custom append without target deduplication SHALL be explicitly labelled limited and block uncertain outcomes rather than require unrelated metadata tables or blindly retry.

#### Scenario: Same timestamp different entity
- **WHEN** two entities or groups write to the same table in one bucket
- **THEN** distinct identities preserve both rows without timestamp-only overwrite.

#### Scenario: Repeated finalized row
- **WHEN** a finalized row is submitted twice to an identity-capable target
- **THEN** the same effect key is used rather than generating a new row identity.

#### Scenario: Revision edit
- **WHEN** a scale, membership or destination revision changes
- **THEN** new rows use the new revision while accepted old rows keep their original encoding and identity.

#### Scenario: Completely silent bucket
- **WHEN** an applied group receives no samples during an entire configured bucket
- **THEN** the time-driven closure emits one scoped no_data/skipped outcome without an SQL data row, and the UI does not imply healthy complete recording.
