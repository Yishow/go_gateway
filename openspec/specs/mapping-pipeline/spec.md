# mapping-pipeline Specification

## Purpose
TBD - created by archiving change add-device-data-pipeline. Update Purpose after archive.

## Requirements

### Requirement: Mapping definition
The system SHALL define mappings from source points to target tags with an ordered transform pipeline.

#### Scenario: Create mapping
- WHEN an operator maps a point to a tag
- THEN the mapping is stored with its transform steps

---
### Requirement: Transform steps
The system SHALL support transform steps including decode, type cast, scaling, lookup table, conditional rules, and formula expressions.

#### Scenario: Scale and offset
- WHEN a raw value is processed with scale and offset
- THEN the final value is computed and stored

---
### Requirement: Formula function set
The system SHALL support formula expressions with operators (+, -, *, /, %), comparisons (>, <, >=, <=, ==, !=), boolean operators (and, or, not), conditional function (if), math functions (abs, min, max, clamp, round, floor, ceil, sqrt, pow), bitwise functions (bitand, bitor, bitxor, shiftl, shiftr), and conversion functions (int, float, bool, string, coalesce).

#### Scenario: Clamp and round
- WHEN a formula applies clamp and round to a value
- THEN the system returns the expected numeric result

---
### Requirement: Formula expression syntax
The system SHALL accept formula expressions using infix operators, snake_case function names, and boolean operators and/or/not.

#### Scenario: Evaluate infix expression
- WHEN a formula uses infix operators with snake_case functions
- THEN the system evaluates the expression without syntax errors

---
### Requirement: Validation of transforms
The system SHALL validate mapping steps for type compatibility and required parameters.

#### Scenario: Invalid cast
- WHEN a mapping attempts to cast text to a numeric type without a parser
- THEN the system rejects the mapping with a validation error

---
### Requirement: Raw and final value capture
The system SHALL preserve raw values alongside final transformed values.

#### Scenario: Preserve raw
- WHEN a value is transformed
- THEN the raw and final values are available for preview and storage

---
### Requirement: Mapping preview
The system SHALL provide a preview result for a mapping using live or sample values.

#### Scenario: Preview with sample
- WHEN an operator runs preview
- THEN the system shows raw, step-by-step, and final values

---
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

---
### Requirement: Rule-derived mappings use shared identity and signature contracts
The mapping pipeline SHALL compare rule-derived mappings using the shared rule-owned identity and proposed-signature contracts defined by source-rule orchestration.

#### Scenario: Same identity with new signature preserves applied mapping
- **WHEN** the system recomputes a rule-derived mapping with the same rule-owned identity and a different proposed signature
- **THEN** the existing applied mapping remains intact
- **AND** the recomputed mapping is surfaced as a reviewable `out_of_sync` candidate

---
### Requirement: Segment-level diagnostic output

The system SHALL expose diagnostic output for each flow segment (`source`, `transform`, `sink`) to support UI visualization.

#### Scenario: Diagnostic payload for preview
- **WHEN** preview or live evaluation is requested
- **THEN** the system returns segment-level latest value, quality, timestamp, and error information

#### Scenario: Sink write visibility
- **WHEN** transformed values are written to storage
- **THEN** the system provides sink segment status indicating write success or failure
- **AND** includes last successful write timestamp

---
### Requirement: Source rule MAY seed a default transform pipeline
The system SHALL allow a mapping’s initial `transform_pipeline` to be populated automatically when created from a source rule that declares a target data type different from the point read type. The automatic content SHALL be limited to validated steps required to satisfy that intent (typically a single `cast` step) unless a broader policy is explicitly defined elsewhere.

#### Scenario: Default cast is validatable
- **WHEN** a mapping is created or updated from source-rule synchronization with differing point and target types
- **THEN** the automatically inserted steps SHALL pass the same transform validation as manually authored pipelines
- **AND** execution SHALL produce values compatible with the tag’s declared data type

---
### Requirement: Automatic pipeline updates SHALL NOT silently destroy unrelated manual steps
The system SHALL define and enforce a deterministic policy when a mapping already contains non-empty transform steps that were not derived from the rule (for example manual scale or formula). The policy MUST either preserve operator intent with explicit conflict resolution or refuse automatic updates until resolved.

#### Scenario: Conflict is visible to the operator
- **WHEN** automatic rule synchronization would replace or prepend transforms on a mapping that has non-default manual steps
- **THEN** the system SHALL NOT silently drop manual steps
- **AND** SHALL surface a conflict or required confirmation according to the chosen policy

---
### Requirement: Source rule MAY seed a default scale step
The system SHALL allow a mapping’s `transform_pipeline` to include an automatically generated **`scale`** step when the source rule declares non-default scale parameters. The step SHALL use the same `scale` transform type and validation rules as manually authored pipelines.

#### Scenario: Scale-only rule without cast
- **WHEN** a rule keeps protocol read type equal to target type but sets scale parameters
- **THEN** the derived mapping MAY contain only a `scale` step (or `scale` combined with other allowed steps)
- **AND** validation SHALL succeed

---
### Requirement: Source preserving exact mapping conversion
Neutral default point-to-tag mappings SHALL preserve the source data type. Explicit non-neutral numeric scaling SHALL default to float64 and reject integer inputs that cannot be represented without loss. New mappings SHALL cast the scaled result to their declared target type; saved pipelines SHALL retain their stored step order. Explicit casts SHALL reject invalid, non-finite, out-of-range, fractional-to-integer, or precision-losing numeric conversions with an error rather than publish zero or a rounded good value. Existing persisted mappings SHALL remain unchanged.

#### Scenario: Default uint64 reaches SQL exactly
- **WHEN** an operator maps a uint64 point with default UI options and records 9007199254740993
- **THEN** the production mapping, journal, outbox and SQLite readback retain exactly 9007199254740993.

#### Scenario: Invalid conversion fails
- **WHEN** a cast receives invalid numeric text, a negative unsigned input, an overflowing integer, NaN, infinity, a fractional integer input or an integer not exactly representable in the target float
- **THEN** it returns an error and no fabricated zero or rounded good sample is recorded.

#### Scenario: Compatible checked scaling
- **WHEN** a default int16 mapping explicitly uses scale 0.5 and offset 10, or a new integer target scales a fractional result
- **THEN** the former declares float64 and records its legitimate scaled value, while the latter rejects the fraction; stored pipeline step ordering is unchanged.

#### Scenario: Confirmed save survives source synchronization
- **WHEN** an operator explicitly saves a scaled mapping, another mapping on the same source rule is edited without confirmation, and synchronization runs twice
- **THEN** only the successful save retains its declared type and accepted pipeline signature; the other mapping remains out of sync, and a later source edit requires explicit candidate reapplication.

#### Scenario: Exact live value and failed conversion quality
- **WHEN** runtime streams an integer outside JavaScript safe range or a cast fails, including nonzero numeric underflow
- **THEN** the large integer is represented exactly as decimal text, while the failed sample has bad quality and no transformed good value or SQL output.
