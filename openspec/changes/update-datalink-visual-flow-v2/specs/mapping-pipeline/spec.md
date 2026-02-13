## ADDED Requirements
### Requirement: Flow state machine for mapping lifecycle

The system SHALL track each mapping pipeline with explicit lifecycle states: `draft`, `validated`, `active`, and `error`.

#### Scenario: Draft to validated transition
- **WHEN** an operator completes mapping configuration and runs validation
- **THEN** the mapping state transitions from `draft` to `validated`

#### Scenario: Validated to active transition
- **WHEN** an operator activates a validated mapping
- **THEN** the mapping state transitions to `active`

#### Scenario: Runtime failure transition
- **WHEN** source read, transform, or write execution fails
- **THEN** the mapping state transitions to `error`
- **AND** error metadata includes failed segment and reason

### Requirement: Segment-level diagnostic output

The system SHALL expose diagnostic output for each flow segment (`source`, `transform`, `sink`) to support UI visualization.

#### Scenario: Diagnostic payload for preview
- **WHEN** preview or live evaluation is requested
- **THEN** the system returns segment-level latest value, quality, timestamp, and error information

#### Scenario: Sink write visibility
- **WHEN** transformed values are written to storage
- **THEN** the system provides sink segment status indicating write success or failure
- **AND** includes last successful write timestamp
