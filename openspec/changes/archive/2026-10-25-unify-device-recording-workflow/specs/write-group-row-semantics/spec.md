## MODIFIED Requirements

### Requirement: Deterministic UTC snapshot closure
The system SHALL create basic snapshots in UTC half-open intervals [start,end), selected by observed_at rather than arrival. Legacy latest members SHALL choose the greatest observed_at within the bucket, using lexicographically greatest stable sample_id to break a time tie. Explicit mean gauge members SHALL instead use the eligible accepted set under window-average-recording, with separate source/result types; state and counter_raw members retain latest eligible semantics. A bucket SHALL close once at end plus the persisted nonnegative allowed_lateness_seconds, default zero; samples arriving after closure MUST NOT reopen it. Applied groups MUST close buckets from a time-driven scheduler even if no sample arrived, emitting one scoped no_data/skipped outcome for a silent bucket. Latest snapshots MUST NOT be described as interval averages, usage or every-sample history. Only an explicitly persisted mean revision SHALL be described as an arithmetic mean; cumulative usage requires the existing confirmed counter contract.

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

#### Scenario: Explicit mean does not rewrite latest
- **WHEN** one applied latest group and another explicitly selected mean revision receive the same eligible values 10 and 30
- **THEN** latest retains 30 and mean stores 20 with sum/count evidence, and replay never reinterprets either frozen payload.

### Requirement: Explicit freshness and completeness
Legacy wide/custom group defaults SHALL require all enabled members, skip a row if a required member is missing/bad/stale, and record the bucket and cause. New managed mean wide groups SHALL default to optional partial members with verified support and SHALL retain NULL/evidence even for an all-missing window; an explicitly required member SHALL still gate the row. Long profile records remain independent per member. Existing saved required flags SHALL not be changed by this new default. Freshness SHALL compare end minus observed_at against persisted max_age_seconds, defaulting to the interval, without carrying any value from another bucket. Explicit partial mode SHALL be accepted only when verified storage supports nullable values and member-level observation/quality metadata.

#### Scenario: Missing member
- **WHEN** temperature is valid but pressure has no good value in the bucket
- **THEN** the default group emits no SQL row and reports the incomplete bucket; it does not claim complete delivery.

#### Scenario: Partial mode
- **WHEN** a supported group explicitly uses partial mode and one member is missing
- **THEN** the missing member is SQL NULL with its reason and observation metadata, not zero or a stale value.

#### Scenario: Unsupported partial target
- **WHEN** a custom table cannot store NULL or required quality/provenance metadata
- **THEN** partial configuration is rejected before activation.

#### Scenario: New wide all-missing mean window
- **WHEN** every optional member in a supported new mean wide window lacks an eligible observation
- **THEN** its row stores NULL fields and per-member missing evidence; the same condition in an existing latest wide group retains its saved no-data/skip behavior.
