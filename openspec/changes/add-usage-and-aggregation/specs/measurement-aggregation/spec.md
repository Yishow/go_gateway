## ADDED Requirements

### Requirement: Sample-based window statistics
The system SHALL aggregate valid pre-decimation samples into time-weighted mean, sampled extrema, count, valid duration and coverage using an explicit bounded-hold policy.

#### Scenario: Unequal sampling intervals
- **WHEN** value 10 covers 10 seconds and value 20 covers 50 seconds
- **THEN** the mean is 1100/60 rather than 15, with 60 valid seconds and sampled extrema 10 and 20.

#### Scenario: No valid data
- **WHEN** a window has no valid coverage
- **THEN** mean/min/max are null and coverage is zero, not a zero-valued measurement.

### Requirement: Counter usage preserves baseline and precision
The system MUST derive usage from comparable consecutive readings in the same series_epoch and retain original readings and exact differences.

#### Scenario: Meter example
- **WHEN** readings are 12000.0, 12000.4 and 12000.7 kWh
- **THEN** the first reading establishes the baseline, deltas are 0.4 and 0.3, and their total is 0.7 kWh.

### Requirement: Reset and rollover handling is explicit
The system SHALL flag discontinuities and apply rollover compensation only with a configured modulus and evidence sufficient to rule out multiple wraps.

#### Scenario: Confirmed single rollover
- **WHEN** readings change from 65530 to 4 with modulus 65536 and a verified single-wrap bound
- **THEN** usage is 10 and the rollover event is retained.

#### Scenario: Unexplained negative difference
- **WHEN** a monotonic counter decreases without reset or rollover evidence
- **THEN** the interval is uncertain and is not repaired with absolute value or an assumed zero baseline.

### Requirement: Signed totals and identified deltas
The system MUST allow negative differences only for confirmed signed totals and sum delta values only once per trustworthy event or non-overlapping interval identity.

#### Scenario: Repeated delta delivery
- **WHEN** two deliveries carry the same source interval identity
- **THEN** the interval is counted once even if delivery ids differ.

#### Scenario: Equal independent delta events
- **WHEN** two different non-overlapping source intervals each report 1 unit
- **THEN** the usage total is 2 units rather than deduplicating by equal value.

### Requirement: Gaps and boundaries remain visible
The system SHALL retain unallocatable interval totals and MUST NOT assign them to unsupported subperiods or claim incomplete daily data is complete.

#### Scenario: Ten-minute outage
- **WHEN** trustworthy continuous-epoch readings at 10:00 and 10:10 increase from 1000 to 1002
- **THEN** the ten-minute total is 2 and minute allocations remain unknown without additional evidence.

#### Scenario: Missing midnight boundary
- **WHEN** a daily counter report lacks trustworthy boundary information
- **THEN** complete daily usage is unknown and any known subtotal is separately labelled.

### Requirement: Explicit estimated rate integration
The system SHALL integrate compatible rate units only by an explicitly chosen method across valid bounded gaps and label outputs estimated.

#### Scenario: Constant power estimate
- **WHEN** valid 2 kW readings cover 1800 seconds under the configured integration policy
- **THEN** estimated energy is 1 kWh and it is not added again to authoritative meter usage.

### Requirement: State and non-numeric recording
The system MUST preserve state transitions, unknown durations, bitmask changes and text values without numeric averaging or fabricated events.

#### Scenario: Initial running state followed by a gap
- **WHEN** the first observed state is running and later readings stop beyond max_hold
- **THEN** an initial snapshot is recorded, no prior start event is invented, and the uncovered duration is unknown.

### Requirement: Derived measurement validation
The system SHALL evaluate only validated typed expressions with explicit units, dependencies and freshness limits.

#### Scenario: Invalid sum or dependency cycle
- **WHEN** an expression adds incompatible units or creates a cycle
- **THEN** activation is blocked with the offending inputs rather than producing a number.

### Requirement: Timezone and deterministic correction
The system MUST store UTC timestamps, apply plan-local calendar boundaries, and revise affected results deterministically when retained evidence permits correction.

#### Scenario: Late sample after provisional result
- **WHEN** an older valid sample arrives within the configured correction horizon
- **THEN** affected results receive higher calculation revisions and old deliveries cannot overwrite them.

### Requirement: Atomic calculation progress
The system SHALL commit accepted input progress, calculator checkpoint, derived results and delivery intent atomically through the local recording transaction.

#### Scenario: Crash between calculation and delivery
- **WHEN** the process crashes after local commit but before external delivery
- **THEN** restart resumes the pending output without recomputing it as extra consumption.

### Requirement: Aggregation preserves metering scope
The system MUST require explicit non-overlapping scope for cross-equipment usage totals and expose member completeness.

#### Scenario: Parent and child meters selected
- **WHEN** a total includes a meter and a known submeter under it
- **THEN** double-counting is blocked or requires an explicit reviewed expression rather than an automatic sum.
