## ADDED Requirements

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
