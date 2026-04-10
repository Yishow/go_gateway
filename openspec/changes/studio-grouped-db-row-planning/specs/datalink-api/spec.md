## ADDED Requirements

### Requirement: Database connector and mapping APIs expose grouped-row metadata

The system SHALL expose grouped database-row planning metadata through connector and mapping APIs so `/studio` can plan interval-driven row writes without inventing a version-specific backend contract.

#### Scenario: Connector returns default write interval
- **WHEN** a client creates, updates, or reads a database connector for this workflow
- **THEN** the payload includes `default_write_interval_seconds`
- **AND** the field defaults to `15` unless explicitly changed

#### Scenario: Mapping preserves grouped compatibility fields
- **WHEN** a client creates, updates, or reads a database target mapping
- **THEN** the payload includes nullable `group_key` and optional `write_interval_seconds`
- **AND** a `null` `group_key` remains legacy single-member behavior rather than being coerced into a shared group

#### Scenario: Existing records keep compatible migration behavior
- **WHEN** existing database connectors and mappings are read after the grouped-row migration
- **THEN** connectors expose `default_write_interval_seconds = 15` unless a different value was explicitly persisted
- **AND** existing mappings keep `group_key = null` so they remain single-member rows

#### Scenario: Effective interval resolves predictably
- **WHEN** a client reads or validates a database mapping
- **THEN** the effective interval equals `write_interval_seconds` when that override is present
- **AND** otherwise the effective interval equals the connector `default_write_interval_seconds`

### Requirement: Rule-scoped database candidate and apply APIs carry grouped planning metadata

The system SHALL return grouped database planning metadata in rule-scoped candidate, dry-run, and apply payloads so the frontend can review row groups before commit.

#### Scenario: Candidate payload exposes inferred group planning
- **WHEN** a client requests rule-scoped database output candidates for tags such as `meter/A1` and `meter/kw`
- **THEN** each candidate includes the inferred `group_key`, normalized `column_name`, and effective `write_interval_seconds`
- **AND** slashless tags continue to return `group_key = null`

#### Scenario: Column-name normalization is deterministic
- **WHEN** a client requests rule-scoped database output candidates for slash-based tag keys
- **THEN** the remaining suffix after the first slash is normalized to lowercase for `column_name`
- **AND** any additional `/` separators are converted to `_`

#### Scenario: Apply rejects incompatible grouped members
- **WHEN** a client dry-runs or applies grouped database candidates whose connector, schema, table, write mode, timestamp column, or effective interval do not match
- **THEN** the API returns an explicit blocking reason
- **AND** the system does not silently merge or remap incompatible members

### Requirement: Database write history reflects grouped row flushes

The system SHALL expose grouped database write history through `GET /api/v1/datalink/db-targets/connectors/:id/write-history` and SHALL report row outcomes per grouped flush rather than per raw tag event.

#### Scenario: Write history query returns grouped flush records
- **WHEN** a client requests grouped database write history for a connector
- **THEN** the response returns recent grouped flush records in reverse chronological order
- **AND** each record includes at least `observed_at`, `status`, `row_count`, `group_key`, `table_name`, `effective_interval_seconds`, and `error_summary` when failed

#### Scenario: Grouped flush records one row outcome
- **WHEN** four compatible `meter/*` tags flush into one database row for the same interval bucket
- **THEN** write history records one row outcome for that grouped flush
- **AND** the history remains directly renderable in Studio diagnostics

#### Scenario: Interval buckets use wall-clock alignment
- **WHEN** grouped database writes are bucketed for an effective interval
- **THEN** the bucket start is aligned by `observed_at.UTC().Truncate(effective_interval)`
- **AND** all writes within the same aligned bucket contribute to the same grouped row outcome
