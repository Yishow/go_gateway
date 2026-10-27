## Purpose

Extend canonical circuit recording to verified MySQL and MariaDB destinations using the shared typed long-table and durable delivery contracts. Keep server-family capabilities, schema recovery and optional native partition maintenance explicit and independently verified.

## ADDED Requirements

### Requirement: Family-specific verified MySQL and MariaDB capabilities

The backend SHALL retain connector kind mysql and expose actual server_family, server_version and operation-specific capabilities for each saved destination. It MUST verify family, target database/schema, transactional engine and permissions independently from connection success. Managed long recording SHALL initially target verified MySQL 8.4 and MariaDB 11.4 baselines; unverified families, versions or operations SHALL remain unavailable. Probe SHALL be read-only, use the saved connector revision and preserve credential/TLS protections. A changed identity or capability revision SHALL invalidate dependent confirmation. MariaDB MUST NOT inherit verification solely because MySQL uses the same driver.

#### Scenario: Separate server-family evidence
- **WHEN** MySQL 8.4 has passed managed-long verification but the selected MariaDB destination has not
- **THEN** the backend reports the actual MariaDB family and blocks its unverified managed operation without claiming MySQL evidence covers it.

#### Scenario: Connected but forbidden metadata
- **WHEN** a saved account authenticates but cannot inspect the selected table
- **THEN** inspection reports forbidden or unconfirmed rather than missing, no automatic DDL runs, and the connector is not write-ready.

### Requirement: Shared profiles with exact MySQL long values

MySQL and MariaDB SHALL consume the common canonical profiles and finalized payloads without implementing a separate counter or scheduler. Instantaneous intervals SHALL retain 1–300 seconds and a 15-second default; cumulative intervals SHALL retain 1–600 minutes and a 10-minute default with both raw readings and differences. Managed long storage SHALL use InnoDB, UTC DATETIME(6), exact integer/DECIMAL slots, independent source identity and the common NULL/status semantics. Values outside the declared DECIMAL(65,30) representation SHALL be rejected before intake with exact_value_unrepresentable, never rounded or converted to floating point. The composite key SHALL include bucket_start and record_id, with case-exact identities and source/time indexes.

#### Scenario: Exact values and source identity
- **WHEN** confirmed profiles deliver uint64 18446744073709551615, integer readings 9007199254740992 then 9007199254740993, decimal 12000.0 then 12000.4, Unicode text and case-distinct source IDs
- **THEN** raw values and identities round-trip unchanged, corresponding differences are 1 and 0.4, and UTC microsecond timestamps retain their window identity.

#### Scenario: Unrepresentable decimal and first baseline
- **WHEN** a value exceeds the confirmed decimal representation or a counter has its first usable window
- **THEN** the former is rejected with an actionable exactness error and the latter stores raw with NULL delta and baseline status, without fabricating zero consumption.

### Requirement: Recoverable per-statement managed DDL

Preparation SHALL use revision-bound preview, explicit confirmation and the existing operation ledger. Each MySQL/MariaDB DDL step SHALL retain its started and verified result because schema batches MUST NOT be represented as atomically reversible. Recovery SHALL reuse the operation identity, inspect real catalog state and resume only missing compatible owned steps. Index syntax SHALL match the detected family; incompatible same-name objects or stale confirmation SHALL block readiness without destructive repair. No operation SHALL create a database or reinterpret forbidden inspection as missing storage.

#### Scenario: Crash after table creation
- **WHEN** the process stops after an owned table is created but before its index step or confirmation response completes
- **THEN** recovery verifies the existing table, resumes the missing index through the same operation identity and marks ready only after all required results are verified; it does not claim the table was rolled back.

#### Scenario: Foreign object collision
- **WHEN** the intended table name already holds an incompatible or foreign object
- **THEN** confirmation fails with ownership or schema mismatch and the object is neither dropped nor altered.

### Requirement: Transactional MySQL effects and safe missing-table recovery

Rows and matching receipts SHALL commit in one InnoDB transaction through the existing production outbox. Repeated effect identity with the same digest SHALL be recognized as already committed; a different digest SHALL be a conflict. The adapter MUST NOT use INSERT IGNORE, REPLACE or unconditional overwrite to conceal errors. Uncertain commit acknowledgement SHALL be resolved through durable receipts without recalculating counters. Catalog-confirmed missing managed tables SHALL be recreated only under the saved recreate_missing policy, valid frozen layout/generation and unchanged destination scope. Concurrent repair SHALL serialize and reuse the ledger; committed receipts SHALL remain intact and lost history MUST NOT be claimed recovered.

#### Scenario: Commit acknowledgement lost
- **WHEN** SQL commits a row and receipt but the acknowledgement is lost and the gateway restarts
- **THEN** the same frozen effect is resolved once through its receipt, with one destination row and no extra counter delta.

#### Scenario: Missing table versus unsafe recovery
- **WHEN** a confirmed fixture table disappears with repair enabled, or a same-name foreign table replaces it
- **THEN** the missing owned layout is restored and eligible pending effects resume once, while the foreign replacement blocks repair and remains unchanged.

### Requirement: Optional native monthly partition maintenance

Storage SHALL default to plain InnoDB and offer explicit native_monthly only for verified family capabilities. Partitioned long tables SHALL use UTC monthly RANGE COLUMNS on bucket_start, include that column in every unique key, omit foreign keys and retain a separate unpartitioned receipt table. Initial confirmation SHALL show the current month and following three months plus past/future catch-all coverage. Authorized maintenance SHALL idempotently extend only the bounded future horizon, preserve rows and record actual results. It MUST NOT automatically delete, truncate or merge historical partitions or claim retention is required. Maintenance failure SHALL remain visible; delivery SHALL continue only where validated catch-all coverage still safely accepts the row.

#### Scenario: Month rollover and prolonged downtime
- **WHEN** the clock crosses a UTC month boundary or the gateway returns after a year offline
- **THEN** maintenance safely establishes the current bounded future horizon, accepted old data remains queryable, and no unbounded sequence of missed monthly partitions or historical deletion is performed.

#### Scenario: Maintenance permission lost
- **WHEN** partition maintenance is denied but the verified MAXVALUE catch-all still accepts records
- **THEN** recording continues with partition_maintenance_blocked visible; incompatible or unverifiable coverage instead blocks the affected destination and retains pending data.

### Requirement: Verified workflow and independent database acceptance

Studio SHALL expose the actual family, operation availability, plain/native_monthly choice and separate saved, connected, prepared, applied, locally durable, SQL committed, readback and cleanup results. Errors SHALL have localized actionable codes without credentials or raw SQL. Managed group test-write SHALL perform real scoped write/readback/cleanup. Existing custom groups and frozen old destinations SHALL remain valid. Full support SHALL require separate real disposable MySQL and MariaDB acceptance for plain/partitioned storage, partial DDL recovery, receipts and missing-table repair; mock-only, skipped or driver-only evidence MUST NOT release the gate.

#### Scenario: Partial operation shown in Basic
- **WHEN** Basic schema preparation stops after one successful DDL step
- **THEN** the operator sees the partial result and resumes the same operation without redoing device setup or receiving a false complete-recording claim.

#### Scenario: Acceptance matrix incomplete
- **WHEN** MySQL passes real acceptance but MariaDB tests are skipped
- **THEN** MariaDB remains unverified, its tests are recorded NOT RUN, and the combined change is not declared complete.

### Requirement: Confirmed exact decimal counter sources

The recording path SHALL consume an explicitly saved positive finite decimal-text exact_recording_scale only for confirmed integer counter sources. It MUST validate the original integer type, measurement identity, definition revision, source-binding revision and series epoch before converting exactly once to ExactDecimal before journal intake. Scale persistence and arithmetic MUST NOT pass through floating point. A scale edit SHALL create a new series epoch and invalidate the old confirmation; frozen accepted effects and baselines SHALL retain their original proof without recalculation. Absent scale SHALL preserve existing integer and floating behavior. Invalid scale, unsupported source encoding, stale proof and unrepresentable results SHALL fail closed without creating an accepted decimal sample. The shared implementation SHALL be owned by add-mysql-mariadb-circuit-recording and consumed unchanged by the subsequent SQL Server batch.

#### Scenario: Confirmed integer source produces exact decimal raw and delta
- **WHEN** confirmed integer readings 120000 then 120004 carry exact_recording_scale `0.1` through instantaneous counter_raw and cumulative profiles
- **THEN** journal, frozen outbox and destination readback preserve decimal raw 12000.0 then 12000.4 with exact cumulative difference 0.4 using the existing counter engine.

#### Scenario: Changed scale preserves accepted history
- **WHEN** a confirmed source scale changes while accepted effects or a comparable baseline exist
- **THEN** the new definition has a new epoch, stale confirmation is rejected, new intake requires the existing forward confirmation boundary, and old effects retain their original scale and payload.

#### Scenario: Unsupported or absent scale
- **WHEN** scale evidence is invalid, the envelope has a noninteger original value, or no scale was saved
- **THEN** invalid or mismatched evidence cannot authorize exact decimal intake, while the absent-scale path retains its existing exact integer or floating provenance without rounding.
