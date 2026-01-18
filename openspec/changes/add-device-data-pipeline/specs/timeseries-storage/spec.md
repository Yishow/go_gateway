## ADDED Requirements
### Requirement: Time-series schema
The system SHALL store values in a time-series table with columns for timestamp, tag id, typed values, raw value, and quality.

#### Scenario: Insert numeric value
- WHEN a numeric value is collected
- THEN the system inserts a row with value_num and the timestamp

### Requirement: Timestamp precision
The system SHALL store timestamps at a configured precision (seconds or milliseconds) for time-series writes.

#### Scenario: Store second precision
- WHEN the write precision is set to seconds
- THEN stored timestamps are truncated to second resolution

### Requirement: Postgres partitioning
The system SHALL create time-range partitions for the time-series table in PostgreSQL, defaulting to monthly partitions unless configured otherwise.

#### Scenario: Create monthly partition
- WHEN the current month is detected and a partition is missing
- THEN the system creates a partition for the month

### Requirement: SQLite test support
The system SHALL support SQLite for testing using the same logical schema without partitions.

#### Scenario: Run test insert
- WHEN running in SQLite test mode
- THEN inserts succeed on the non-partitioned table

### Requirement: External query support
The system SHALL provide indexes and a stable schema to support external history queries by tag and time range.

#### Scenario: CMS query by tag and time range
- WHEN the CMS queries by tag id and time range
- THEN the schema and indexes support ordered retrieval

### Requirement: No automatic retention
The system SHALL not delete historical data unless explicitly requested by an operator.

#### Scenario: No retention job
- WHEN the system runs scheduled jobs
- THEN no data retention deletion is executed by default
