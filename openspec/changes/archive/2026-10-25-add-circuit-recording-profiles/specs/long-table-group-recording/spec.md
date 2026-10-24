## Purpose

Store canonical circuit recording outputs in a fixed typed long-table layout with optional TimescaleDB time partitioning. Preserve existing ownership, schema confirmation and durable delivery guarantees while retaining legacy wide-table histories and execution.

## ADDED Requirements

### Requirement: Versioned shared typed long storage

New basic recording profiles SHALL use managed_long_v1, with one shared managed table per confirmed workspace, connector, database, schema and storage mode. Adding a circuit or source SHALL add records rather than value columns or per-circuit tables. Records MUST carry stable workspace, group, revision, circuit and source identities, profile kind, window bounds, actual sample times, quality, status, exact typed raw value, optional delta and provenance. Valid values SHALL populate exactly the matching typed value slot; invalid or missing values SHALL remain NULL with reasons. PostgreSQL SHALL preserve uint64 and decimal through NUMERIC and UTC times through TIMESTAMPTZ. SQLite SHALL preserve the same logical values using exact decimal text and normalized UTC timestamp text. Missing layout metadata in existing groups MUST retain legacy wide semantics.

#### Scenario: New point without column migration
- **WHEN** a saved circuit adds another supported typed source to a prepared long layout
- **THEN** the next confirmed group revision records that source with its identity without adding a value column, while requiring the existing revision-bound confirmation even for an empty DDL preview.

#### Scenario: Same Tag label on different circuits
- **WHEN** two circuits expose equally named Tags with different source identities
- **THEN** their rows remain distinct and exact uint64, decimal, bool and text values survive write and readback without coercion to float or text-only semantics.

### Requirement: Optional explicit TimescaleDB mode

PostgreSQL destinations SHALL default to plain long storage and SHALL offer timescale only as an explicit selection. The backend SHALL probe installed extension, required by_range/create_hypertable functionality and permissions without changing the database. Timescale preparation SHALL create an owned empty hypertable partitioned by non-null bucket_start, with a composite unique key including bucket_start and record_id, and a previewed seven-day chunk interval. It MUST NOT convert populated or foreign tables, enable extensions, combine native and Timescale partitioning, add retention jobs or silently fall back to plain mode. SQLite SHALL reject timescale selection. Unsupported capability SHALL return timescale_unavailable with a repair action.

#### Scenario: Plain PostgreSQL without extension
- **WHEN** an operator selects plain storage on PostgreSQL without TimescaleDB
- **THEN** confirmed long-table recording and readback succeed without extension installation.

#### Scenario: Hypertable and unique effect
- **WHEN** timescale is explicitly selected on a capable disposable destination and schema confirmation succeeds
- **THEN** catalog inspection proves an owned hypertable with the requested partition key, records spanning chunk boundaries are queryable, and replaying the same record identity creates no additional row.

#### Scenario: Unavailable extension
- **WHEN** timescale is selected but its required capability or permission is absent
- **THEN** preparation is blocked with an actionable reason, no database objects are changed, and the selected mode is not silently rewritten.

### Requirement: Scoped confirmation and truthful delivery evidence

Long-table preparation SHALL use the existing read-only preview, revision-bound confirmation operation, readiness and apply sequence. Shared storage ownership SHALL bind workspace, destination and layout generation; each referencing group SHALL retain its own confirmed revision. Stale workspace, connector, group or schema versions SHALL invalidate dependent readiness. Delivery SHALL reuse the canonical journal, outbox and receipt mechanism; every effect SHALL retain its frozen payload and destination. Saved, applied, collecting, locally durable, SQL committed and readback verified states MUST remain distinguishable.

#### Scenario: Shared table with stale group confirmation
- **WHEN** one group's source membership changes while another group already uses the shared long table
- **THEN** the changed group requires fresh confirmation without altering the other group's applied snapshot or claiming it confirmed the change.

#### Scenario: Confirmation response lost
- **WHEN** the schema-confirmation response is lost after execution
- **THEN** recovery queries the existing operation identity and actual result rather than issuing unrelated DDL or claiming success from runtime state.

### Requirement: Restore only previously confirmed missing managed storage

A saved recreate_missing policy SHALL be shown during explicit schema confirmation. When enabled, the backend SHALL recover a catalog-confirmed missing managed table using its previously confirmed frozen layout, generation and ownership evidence, then retry eligible uncommitted effects idempotently. Recovery SHALL serialize concurrent attempts through scoped coordination and the operation ledger. It MUST NOT run for generic SQL failures, changed destination revisions, absent confirmation, foreign or incompatible same-name tables, missing ownership evidence or unavailable Timescale capability. Recovery MUST NOT alter existing incompatible objects, create a database, install extensions or recreate a missing SQLite database file. Committed receipt history SHALL remain intact; recovered structure MUST NOT be presented as recovery of lost historical data.

#### Scenario: Confirmed table removed during delivery
- **WHEN** a disposable fixture removes a previously confirmed table with recreate_missing enabled and uncommitted records await delivery
- **THEN** one effective reconstruction restores the confirmed plain or hypertable layout, pending records are retried without duplicates, and diagnostics disclose that previously stored history is not guaranteed restored.

#### Scenario: Unsafe or disabled recovery
- **WHEN** recovery is disabled, a foreign table occupies the name, the connector revision changes, or ownership evidence is missing
- **THEN** delivery remains blocked with schema_repair_blocked or storage_owner_mismatch and no foreign data or existing receipt is altered.

#### Scenario: Permissions or SQLite file missing
- **WHEN** writing fails for permissions or the configured SQLite file is missing
- **THEN** the backend reports the actual actionable failure and does not treat it as permission to create a table or an empty replacement database.

### Requirement: Explicit forward cutover preserves wide history

Existing managed wide and custom groups SHALL retain saved identity, layout, intervals, completeness policies and runtime behavior. Switching selected sources to long profiles SHALL require a reviewed destination and member difference, fresh schema confirmation and CAS-checked apply. An immutable effective boundary SHALL assign each source observation to one active revision, mark boundary-truncated windows partial and reset cumulative baseline. Accepted old records SHALL drain to their frozen destination. Failed cutover SHALL leave the prior applied revision active. Rollback SHALL use a new explicit boundary and MUST NOT delete history, rewrite frozen payloads, backfill old tables or enable duplicate legacy writers.

#### Scenario: Successful switch with pending old rows
- **WHEN** a wide group switches to long profiles while old outbox rows remain pending
- **THEN** pre-boundary data retains its old layout and destination, post-boundary data uses the confirmed long profile, boundary windows are partial, and the first cumulative record has NULL delta.

#### Scenario: Failed switch and rollback
- **WHEN** a switch has stale revisions or preparation fails, or a completed switch is explicitly rolled back
- **THEN** failed preparation preserves the original applied group; rollback establishes a new boundary, drains already accepted long records and preserves all existing tables and history.

### Requirement: Stable history access without automatic deletion

The long layout SHALL expose a documented fixed schema and indexes for workspace, source identity and time-range retrieval. Plain and Timescale modes SHALL return equivalent logical records for the same accepted inputs. This capability MUST NOT add automatic historical deletion, compression, reports or retention configuration as a prerequisite to basic recording.

#### Scenario: Equivalent query across storage modes
- **WHEN** identical fixture inputs are recorded into plain PostgreSQL and TimescaleDB
- **THEN** ordered source/time queries return equivalent typed raw values, deltas and quality states, and no retention or compression policy is created.
