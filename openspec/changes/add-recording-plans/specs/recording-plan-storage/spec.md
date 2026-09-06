## ADDED Requirements

### Requirement: Versioned recording plans with independent streams
The system SHALL persist RecordingPlan revisions and separate source acquisition, raw retention, summary intervals and destination flush policies.

#### Scenario: One source serves curve and usage
- **WHEN** a meter plan selects raw history, one-minute summaries and hourly usage
- **THEN** collection is shared while outputs keep independent stream identities and policies.

### Requirement: Record mode semantics
The system MUST support raw history, window summary, usage interval, state changes, events, batch snapshots and latest-only with truthful capabilities.

#### Scenario: Latest-only has no history
- **WHEN** an operator chooses only latest_only
- **THEN** the UI and query capability state that past curves cannot be reconstructed.

#### Scenario: Flush delay does not reduce readings
- **WHEN** twelve samples are acquired and buffered before one transport flush
- **THEN** every_sample retains twelve samples rather than only the last value.

### Requirement: Managed typed storage and identity
The system SHALL store source identity, observation time, definition revision, quality and exact typed value in managed tables with stable record keys.

#### Scenario: Many devices share the same field name
- **WHEN** two meters both expose active_power
- **THEN** their samples remain distinct by identity without creating new per-device tables.

### Requirement: Explicit reduction without corrupting statistics
The system MUST calculate summaries and usage from pre-decimation inputs and preserve necessary calculation evidence for the configured correction horizon.

#### Scenario: Short spike under on-change reduction
- **WHEN** a valid sampled spike reaches the calculation input
- **THEN** its sampled maximum remains in the summary even when the raw display stream is reduced.

### Requirement: Batch snapshot freshness and trigger identity
The system SHALL bind a batch snapshot to a stable trigger identity and record member-level observation times and completeness.

#### Scenario: Held trigger with one missing member
- **WHEN** the same trigger is polled repeatedly and pressure is missing
- **THEN** one partial batch is recorded with missing pressure, not repeated batches filled with old pressure.

### Requirement: Safe compatibility with existing tables
The system MUST preserve existing single-point/grouped mappings and require explicit review before creating managed tables or changing output mode.

#### Scenario: Existing table name collision
- **WHEN** the requested managed table name already belongs to an incompatible table
- **THEN** preparation is blocked and no existing table or data is altered.
