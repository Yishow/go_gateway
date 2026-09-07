# measurement-semantics Specification

## Purpose

TBD - created by archiving change 'add-measurement-semantics'. Update Purpose after archive.

## Requirements

### Requirement: Per-measurement recording meaning
The system SHALL persist each measurement's quantity, unit, semantic_kind and source binding independently of protocol read type and database column.

#### Scenario: Mixed measurements on one device
- **WHEN** one device supplies temperature, pressure, flow rate and cumulative volume
- **THEN** each measurement retains its own meaning and recording eligibility
- **AND** no device-wide counter or average policy overrides the individual definitions.

#### Scenario: Unconfirmed semantics
- **WHEN** an existing numeric point has no confirmed semantic_kind
- **THEN** legacy recording remains unchanged and new usage or integration modes are blocked with an actionable reason.

---
### Requirement: Stable identity and measurement epochs
The system MUST preserve stable identity for cosmetic edits and start a new series_epoch when source identity or numerical comparability changes.

#### Scenario: Rename versus replace
- **WHEN** a measurement is renamed
- **THEN** its history and sample identity remain intact
- **AND** replacing its meter or changing its scale requires explicit epoch transition without joining old counter baselines.

---
### Requirement: Truthful typed samples
The system SHALL emit the SampleEnvelope from the shared contract, including provenance, original observation time, quality and exact numeric encoding.

#### Scenario: Partially failed acquisition
- **WHEN** pressure fails while temperature succeeds in one collection cycle
- **THEN** temperature is eligible for recording and pressure is marked missing or invalid
- **AND** an old pressure value or zero is not substituted as a new good reading.

#### Scenario: Integer above JavaScript safe range
- **WHEN** uint64 value 9007199254740993 crosses backend, JSON and UI
- **THEN** its decimal digits are preserved exactly and numeric encoding is explicit.

---
### Requirement: Stateful and non-numeric intent is explicit
The system MUST distinguish monotonic counters, signed totals, interval deltas, states, events and text; semantic conversion SHALL require validation rather than unit-name inference.

#### Scenario: Net energy falls
- **WHEN** a confirmed signed_counter decreases from 100 to 98
- **THEN** the definition permits a signed difference rather than automatically declaring a reset.

#### Scenario: Polled delta without identity
- **WHEN** a held PLC register is declared delta without a trustworthy interval or event identity
- **THEN** usage accumulation is blocked rather than counting every poll as new consumption.
