## Purpose

Extend canonical circuit recording to SQL Server through verified destination capabilities and an exact typed long-table adapter. Preserve schema ownership, durable effects and optional native time partition management without changing the shared acquisition or counter engine.

## ADDED Requirements

### Requirement: Verified SQL Server connection and destination scope

The backend SHALL use saved sqlserver connector identity, explicit database/schema and backend-originated read-only connection and metadata probes. New schema drafts SHALL default to dbo without rewriting saved values. Connections SHALL use SQL authentication with encrypted certificate-verified transport and retain existing credential masking/revision protections. Operation capabilities SHALL reflect actual server version, edition, target scope and permissions. Functional verification SHALL use the available real disposable SQL Server instance and record its actual version, edition, compatibility level and host/container architecture. The server version and emulated architecture alone SHALL NOT prevent running the functional acceptance cases. Unverified versions or operations SHALL remain gated for product release; test execution SHALL NOT falsify metadata or change production capability gates. Metadata invisibility SHALL not be treated as a missing table. SQL Server SHALL use its own identifier quoting and parameter syntax rather than PostgreSQL or MySQL syntax.

#### Scenario: Connected but not prepared
- **WHEN** a SQL Server account connects to its saved database but cannot inspect the selected schema
- **THEN** the connection result is preserved, inspection reports forbidden or unconfirmed, dependent preparation stays blocked, and no schema object is created.

#### Scenario: Changed destination and invalid certificate
- **WHEN** a connector identity changes or certificate validation fails
- **THEN** stale preparation cannot authorize the new scope, the security failure is actionable, and the adapter does not silently disable encryption or certificate verification.

### Requirement: Exact SQL Server long-table profile values

The adapter SHALL consume the shared canonical profile and finalized raw/delta results without a second counter engine. Instantaneous 1–300 seconds/default 15 seconds and cumulative 1–600 minutes/default 10 minutes SHALL remain identical to the common contract. Long storage SHALL map UTC times to datetime2(6), signed integers to BIGINT, uint64/exact decimals and deltas to DECIMAL(38,18), floats to FLOAT(53), bool to BIT and Unicode text to NVARCHAR. Provenance SHALL retain valid JSON and source identity SHALL remain case-exact. Values outside the exact decimal range SHALL be rejected before intake with exact_value_unrepresentable, never silently rounded. SQL timestamp/rowversion MUST NOT be used for event times.

#### Scenario: Large integers and Unicode
- **WHEN** fixture profiles write uint64 18446744073709551615, integer readings 9007199254740992 then 9007199254740993, decimal 12000.0 then 12000.4, Chinese text and case-distinct IDs
- **THEN** readback preserves the exact raw values and identities, differences are 1 and 0.4, and timestamps retain UTC microsecond window identity.

#### Scenario: Unknown delta and excessive scale
- **WHEN** the first cumulative record has no comparable baseline or an exact decimal has more fractional digits than the confirmed representation supports
- **THEN** the first raw reading is saved with NULL delta, while the unrepresentable value is rejected before intake instead of being rounded by SQL Server.

### Requirement: Confirmed SQL Server schema operations and repair

Managed schema SHALL use read-only preview, explicit confirmation, revision-bound tokens, actual catalog inspection and the existing operation ledger. Owned table, indexes and optional partition objects SHALL be included in the confirmed digest. Recovery after lost acknowledgement SHALL query the same operation identity and inspect actual transactional results before retrying. Existing foreign or incompatible objects SHALL not be changed. Missing-table recreation SHALL require saved recreate_missing authorization, valid layout/generation, unchanged connector scope and catalog proof; receipt history SHALL remain intact and lost historical values SHALL not be claimed restored. Missing or foreign partition objects SHALL block repair unless exclusive ownership and the frozen layout are proven.

#### Scenario: DDL committed but response lost
- **WHEN** schema DDL commits and its response is lost
- **THEN** the existing operation is inspected and resolved from actual catalog state without issuing an unrelated operation or assuming rollback.

#### Scenario: Removed table and shared foreign partition function
- **WHEN** an authorized owned table is removed, or its expected partition function is replaced by a foreign/shared object
- **THEN** safe missing-table repair resumes eligible pending data, while the foreign/shared object blocks repair and remains unchanged; neither outcome claims lost history was recovered.

### Requirement: Transactional SQL Server durable effects

Production delivery SHALL reuse the canonical journal, frozen outbox payload and destination revisions. SQL row and receipt SHALL commit in one transaction with concurrency-safe effect identity and digest comparison. The same identity/digest SHALL resolve once and a different digest SHALL be a conflict; unconditional MERGE or overwriting historical values MUST NOT mask conflicts. Commit ambiguity SHALL be resolved through receipts. Deadlocks, timeouts and outages SHALL retain data under existing retry/backpressure controls and SHALL NOT reset counter baselines or redirect old records after connector edits.

#### Scenario: Concurrent delivery and commit ambiguity
- **WHEN** two attempts deliver the same effect and one loses its commit acknowledgement
- **THEN** receipt resolution produces one committed row, identical retries are harmless and a conflicting digest is visibly blocked.

#### Scenario: Offline restart and changed connector
- **WHEN** SQL Server is offline during restart and its connector is later edited
- **THEN** accepted records retain their original destination snapshot, pending data remains recoverable and no new consumption is computed from retry.

### Requirement: Optional SQL Server native monthly partitioning

Storage SHALL default to plain and offer native_monthly only after explicit selection and verified capability. Partitioned storage SHALL use an owned datetime2(6) RANGE RIGHT partition function, an owned scheme on the existing PRIMARY filegroup and aligned indexes with bucket_start in unique keys. Confirmation SHALL expose current-month and next-three-month UTC boundaries with past/future catch-all coverage. Authorized maintenance SHALL add bounded future boundaries idempotently using the saved scope and preserve all historical rows. It MUST NOT create databases, files or filegroups, alter foreign/shared functions, drop/truncate/merge history or install retention jobs. Maintenance errors SHALL remain visible and delivery SHALL proceed only where validated coverage remains safe.

#### Scenario: Cross-month rows and catch-all
- **WHEN** records span a UTC month boundary or arrive after a prolonged maintenance outage
- **THEN** sys catalog and readback prove correct range placement and unchanged history, maintenance adds only the bounded future horizon and no missing-month deletion or unbounded backfill occurs.

#### Scenario: Partition permission denied
- **WHEN** plain-table operations are available but partition creation is forbidden
- **THEN** plain recording remains selectable and an explicit partition choice returns partition_unavailable without silent fallback or changes to DBA objects.

### Requirement: Capability-backed SQL Server workflow and acceptance

Studio SHALL use the same circuit profile controls and display saved, connected, inspected, prepared, applied, durable, SQL committed, readback and cleanup results separately. Group test-write SHALL perform scoped real write/readback/cleanup after confirmation. Errors SHALL be localized actionable codes without DSNs, credentials or raw SQL. Existing database modes SHALL remain available. Completion SHALL require a real disposable SQL Server instance for plain/partitioned recording, exact values, failure recovery, concurrency and schema ownership plus browser evidence. A real SQL Server engine running under architecture emulation SHALL provide functional evidence for the cases actually executed, without implying native-platform performance or stability acceptance. Cases that cannot execute, skipped tests and mocked SQL SHALL be recorded NOT RUN and MUST NOT release the operation gate.

#### Scenario: Real write and cleanup differ
- **WHEN** a group test commits and readback succeeds but cleanup fails
- **THEN** the UI preserves those distinct results and provides cleanup repair instead of presenting an unqualified completed test.

#### Scenario: SQL Server fixture unavailable
- **WHEN** unit tests pass but no real SQL Server instance is available
- **THEN** real-database acceptance remains NOT RUN and the change is not declared complete or fully supported.

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
