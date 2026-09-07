# recording-history-reports Specification

## Purpose

TBD - created by archiving change 'add-recording-history-reports'. Update Purpose after archive.

## Requirements

### Requirement: Authoritative bounded history queries
The system SHALL query one explicit plan history source with ownership validation, bounded ranges and consistent revision cutoffs.

#### Scenario: Replicated destinations
- **WHEN** a plan writes the same records to two databases
- **THEN** a history query uses its selected primary source and does not double-count both copies.

---
### Requirement: Truthful trends and resolution
The system MUST display actual sample or summary resolution, units, quality, sampled extrema and gaps without fabricated values.

#### Scenario: Raw history expired
- **WHEN** only minute summaries remain
- **THEN** the UI labels minute resolution and does not present reconstructed five-second measurements.

---
### Requirement: Usage reports preserve evidence and completeness
The system SHALL use persisted calculation results and expose boundary readings, known subtotals, estimated status and calculation revisions.

#### Scenario: Incomplete daily usage
- **WHEN** a daily result has incomplete boundaries
- **THEN** the report does not show the known subtotal as a complete day and provides a route to its evidence.

---
### Requirement: Mixed quantities and states remain distinct
The system MUST support separate unit-aware trends, cumulative readings, interval quantities and state/event timelines within a plan.

#### Scenario: MC mixed layout report
- **WHEN** a plan contains temperature, pressure, rate, volume counter, state, bitmask and production counter
- **THEN** each receives the appropriate chart or timeline without averaging counters or text.

---
### Requirement: Multi-equipment comparisons do not hide missing data
The system SHALL allow compatible equipment comparison and explicit scope-aware totals while exposing missing members and preventing known double counting.

#### Scenario: One selected meter is missing
- **WHEN** a group report lacks a member's usage
- **THEN** its completeness states the missing member and no zero value is invented.

---
### Requirement: Safe precise exports
The system SHALL export consistent-cutoff records with identity, time, unit, quality, precision and revision metadata, and sanitize untrusted textual cells for spreadsheet use.

#### Scenario: Formula-like label and negative numeric reading
- **WHEN** an export includes an untrusted label starting with = and a numeric temperature -12.3
- **THEN** the label cannot execute as a formula and the legitimate numeric value remains -12.3.

---
### Requirement: Runtime repair returns to exact context
The system SHALL show collection and delivery independently and return the operator to the relevant setup context for repair.

#### Scenario: Delivery permissions fail
- **WHEN** a plan's destination is blocked by permissions
- **THEN** its repair action opens the matching plan/connection settings without discarding other successful setup.

---
### Requirement: End-to-end acceptance is evidence-based
The system MUST pass the documented electrical-meter, mixed-sensor, precision, fault-recovery and capacity matrix before the complete recording feature is released.

#### Scenario: Platform not tested
- **WHEN** an implementation has no real PLC or ARM execution evidence
- **THEN** those acceptance entries remain pending and are not reported as passed.
