# mapping-pipeline Specification

## Purpose
TBD - created by archiving change add-device-data-pipeline. Update Purpose after archive.
## Requirements
### Requirement: Mapping definition
The system SHALL define mappings from source points to target tags with an ordered transform pipeline.

#### Scenario: Create mapping
- WHEN an operator maps a point to a tag
- THEN the mapping is stored with its transform steps

### Requirement: Transform steps
The system SHALL support transform steps including decode, type cast, scaling, lookup table, conditional rules, and formula expressions.

#### Scenario: Scale and offset
- WHEN a raw value is processed with scale and offset
- THEN the final value is computed and stored

### Requirement: Formula function set
The system SHALL support formula expressions with operators (+, -, *, /, %), comparisons (>, <, >=, <=, ==, !=), boolean operators (and, or, not), conditional function (if), math functions (abs, min, max, clamp, round, floor, ceil, sqrt, pow), bitwise functions (bitand, bitor, bitxor, shiftl, shiftr), and conversion functions (int, float, bool, string, coalesce).

#### Scenario: Clamp and round
- WHEN a formula applies clamp and round to a value
- THEN the system returns the expected numeric result

### Requirement: Formula expression syntax
The system SHALL accept formula expressions using infix operators, snake_case function names, and boolean operators and/or/not.

#### Scenario: Evaluate infix expression
- WHEN a formula uses infix operators with snake_case functions
- THEN the system evaluates the expression without syntax errors

### Requirement: Validation of transforms
The system SHALL validate mapping steps for type compatibility and required parameters.

#### Scenario: Invalid cast
- WHEN a mapping attempts to cast text to a numeric type without a parser
- THEN the system rejects the mapping with a validation error

### Requirement: Raw and final value capture
The system SHALL preserve raw values alongside final transformed values.

#### Scenario: Preserve raw
- WHEN a value is transformed
- THEN the raw and final values are available for preview and storage

### Requirement: Mapping preview
The system SHALL provide a preview result for a mapping using live or sample values.

#### Scenario: Preview with sample
- WHEN an operator runs preview
- THEN the system shows raw, step-by-step, and final values

### Requirement: Flow state machine for mapping lifecycle

The system SHALL track each mapping pipeline with explicit lifecycle states: `draft`, `validated`, `active`, `out_of_sync`, and `error`.

For rule-derived mappings, `out_of_sync` means the same rule-owned mapping identity still exists but its current proposed signature no longer matches the last applied signature.

#### Scenario: Draft to validated transition
- **WHEN** an operator completes mapping configuration and runs validation
- **THEN** the mapping state transitions from `draft` to `validated`

#### Scenario: Validated to active transition
- **WHEN** an operator activates a validated mapping
- **THEN** the mapping state transitions to `active`

#### Scenario: Active mapping becomes out of sync after rule revision
- **WHEN** a rule revision changes the derived transform pipeline for a mapping that already has applied state
- **THEN** the system marks the mapping as `out_of_sync`
- **AND** requires explicit operator reapply instead of silently replacing the applied pipeline

#### Scenario: Runtime failure transition
- **WHEN** source read, transform, or write execution fails
- **THEN** the mapping state transitions to `error`
- **AND** error metadata includes failed segment and reason

### Requirement: Rule-derived mappings use shared identity and signature contracts
The mapping pipeline SHALL compare rule-derived mappings using the shared rule-owned identity and proposed-signature contracts defined by source-rule orchestration.

#### Scenario: Same identity with new signature preserves applied mapping
- **WHEN** the system recomputes a rule-derived mapping with the same rule-owned identity and a different proposed signature
- **THEN** the existing applied mapping remains intact
- **AND** the recomputed mapping is surfaced as a reviewable `out_of_sync` candidate

### Requirement: Segment-level diagnostic output

The system SHALL expose diagnostic output for each flow segment (`source`, `transform`, `sink`) to support UI visualization.

#### Scenario: Diagnostic payload for preview
- **WHEN** preview or live evaluation is requested
- **THEN** the system returns segment-level latest value, quality, timestamp, and error information

#### Scenario: Sink write visibility
- **WHEN** transformed values are written to storage
- **THEN** the system provides sink segment status indicating write success or failure
- **AND** includes last successful write timestamp
