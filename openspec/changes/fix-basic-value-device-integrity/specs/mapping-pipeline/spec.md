## ADDED Requirements

### Requirement: Source preserving exact mapping conversion
Default point-to-tag mappings SHALL preserve the source data type. Explicit casts SHALL reject invalid, non-finite, out-of-range, fractional-to-integer, or precision-losing numeric conversions with an error rather than publish zero or a rounded good value. Existing persisted mappings SHALL remain unchanged.

#### Scenario: Default uint64 reaches SQL exactly
- **WHEN** an operator maps a uint64 point with default UI options and records 9007199254740993
- **THEN** the production mapping, journal, outbox and SQLite readback retain exactly 9007199254740993.

#### Scenario: Invalid conversion fails
- **WHEN** a cast receives invalid numeric text, a negative unsigned input, an overflowing integer, NaN, infinity, a fractional integer input or an integer not exactly representable in the target float
- **THEN** it returns an error and no fabricated zero or rounded good sample is recorded.
