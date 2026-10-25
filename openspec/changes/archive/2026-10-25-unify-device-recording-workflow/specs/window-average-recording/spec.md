## Purpose

Provide opt-in arithmetic-mean recording over accepted acquisition windows without changing acquisition frequency or reinterpreting existing snapshots. Preserve exact result types, missing-data evidence, counter semantics and durable delivery across managed wide and long layouts.

## ADDED Requirements

### Requirement: Persisted reduction and independent value roles

The system SHALL persist recording_profile.sample_reduction as latest or mean and each instantaneous member's value_role as gauge, counter_raw or state. An omitted reduction SHALL decode as latest for existing groups and snapshots. New instantaneous recording SHALL default to mean with a 15-second interval and retain the existing integer 1–300-second bounds. Cumulative recording SHALL retain the existing integer 1–600-minute bounds and 10-minute default and SHALL reject mean. Gauge mean SHALL require a numeric source and explicit operator selection without requiring a fabricated measurement identity. Counter roles SHALL require the existing confirmed counter binding. Bool/text states SHALL retain latest eligible values and SHALL NOT enter an average. Names and units MUST NOT infer roles.

#### Scenario: Legacy reload and explicit edit
- **WHEN** an existing group without sample_reduction is loaded, or its interval is outside new instantaneous bounds
- **THEN** its saved latest semantics and interval remain unchanged; only an explicit compatible new revision enters mean recording.

#### Scenario: Invalid combinations
- **WHEN** mean is requested for a cumulative profile, bool/text gauge, unknown role or unconfirmed counter role
- **THEN** validation rejects the incompatible intent with an actionable typed error before Apply and no writer starts.

#### Scenario: One source in both profiles
- **WHEN** a confirmed energy counter is selected as counter_raw in an instantaneous snapshot and as cumulative in its separate profile
- **THEN** the instantaneous record stores its latest eligible raw value, the cumulative record uses its own existing interval and baseline, and neither averages the counter.

### Requirement: Arithmetic mean uses unique eligible observations

For each numeric gauge, the runtime SHALL calculate the arithmetic mean of eligible observations whose observed_at is inside the actual UTC half-open window. Each accepted SampleID SHALL contribute once with equal weight. Repeated identical content SHALL not change sum or count; conflicting content SHALL retain the existing identity-conflict failure. Bad, invalid, stale, out-of-window and closed-window late samples SHALL not contribute to numerator or denominator. Reduction MUST NOT interpolate, zero-fill, carry values between windows, use flush time as observation time or modify acquisition cadence.

#### Scenario: Equal weights and excluded data
- **WHEN** a 15-second window contains eligible values 10, 20 and 30, a retransmission of the 10 sample, and a bad value 1000
- **THEN** the mean is 20, valid_count is 3, rejected_count reflects the distinct excluded bad observation, and the retransmission contributes nothing additional.

#### Scenario: Boundary and delivery order
- **WHEN** the same accepted observations arrive in a different order before closure and another observation is exactly at the window end
- **THEN** the accepted set produces the same mean and evidence, the boundary observation belongs to the next window, and an observation arriving after the existing closure deadline cannot reopen the result.

#### Scenario: Actual cutover bounds
- **WHEN** a reviewed activation or cutover starts an effective window at second 7 inside the nominal [0,15) window
- **THEN** only observations in [7,15) contribute, the row retains actual and nominal bounds, and truncated coverage is partial.

### Requirement: Mean results have an explicit numeric contract

Exact int, uint and decimal inputs SHALL use an exact sum and integer count without float64 or JavaScript Number coercion. Exact means SHALL have decimal result type independent of source_value_type. Finite decimal results SHALL be projected to at most 18 fractional digits using half-even rounding; provenance SHALL retain exact sum, valid_count, output scale and rounding so a rounded result is not claimed infinitely exact. Float sources SHALL use a deterministic observed_at plus SampleID ordering and explicitly declare basis=float64. Arithmetic or verified SQL representation failures SHALL be surfaced without substituting zero, truncating a fraction or changing the source Tag type.

#### Scenario: Exact large integer and fraction
- **WHEN** exact observations 9007199254740992 and 9007199254740993 close one window
- **THEN** the stored decimal mean is exactly 9007199254740992.5, the original source type remains unchanged, and API/SQLite/PostgreSQL round-trip the result and sum/count without numeric coercion.

#### Scenario: Repeating decimal
- **WHEN** exact observations 0, 0 and 1 close one window
- **THEN** the result is 0.333333333333333333 with sum 1, count 3, scale 18 and half-even metadata.

#### Scenario: Integer column cannot accept the result
- **WHEN** source integers 1 and 2 require result 1.5 but the requested destination only provides an incompatible integer column
- **THEN** preparation or runtime-proof validation blocks the incompatible layout before recording and does not alter the existing column or round to an integer.

### Requirement: Missing data and coverage remain evidence based

Each new reduction result SHALL retain aggregation version/method, source/result types, basis, unique valid_count, rejected_count, actual first/last observation times, coverage and stable reasons. Coverage SHALL be complete, partial, unknown or missing. No eligible observations SHALL produce NULL with missing coverage and an evidence-backed no-samples/bad/stale/invalid reason. A known failed observation, proven gap or truncated window with usable observations SHALL produce the available mean with partial coverage. Usable observations without sufficient whole-window coverage proof SHALL use unknown. Complete SHALL require existing saved coverage proof, never absence of errors. The existing scoped delivery response SHALL expose verified value-free method/count/coverage evidence for both new managed layouts from frozen payloads and matching receipts; invalid or unavailable proof SHALL remain unavailable rather than complete. It SHALL not become a historical-value query endpoint. The system MUST NOT fabricate expected sample counts, coverage percentages or a disconnection cause and MUST NOT add a new scheduler or health-monitoring subsystem to obtain them.

#### Scenario: Partial window and recovery
- **WHEN** a window has valid values 10 and 30 plus a known failed read, followed by a window with valid values 20 and 40
- **THEN** the first mean is 20 with partial coverage and the known failure reason; the next mean is 30 using only its own observations and available coverage proof.

#### Scenario: Silent outage
- **WHEN** no eligible observation exists during a scheduled recording window and there is no saved network-disconnection event
- **THEN** the record contains NULL/no-samples and missing coverage, does not claim zero or proven disconnection, and the time-driven closure still completes.

#### Scenario: Counts do not prove completeness
- **WHEN** a window has a valid value but no authoritative expected-observation or continuity evidence
- **THEN** it stores the computed value and actual counts with unknown coverage, not 100 percent or complete.

### Requirement: Managed layouts preserve independent reductions and counters

The system SHALL support new managed_wide_profile_v1 device/circuit snapshots and managed_long_v1 per-source records using the same persisted roles, profile intervals and reduction rules. Existing managed_wide_v1 and custom groups SHALL remain unchanged. New wide value columns SHALL use frozen result types; cumulative wide columns SHALL retain raw values and per-member provenance.counter SHALL retain exact delta, actual bounds, status and policy. Long records SHALL retain the existing typed value and delta slots with explicit source/result metadata. Counter calculation SHALL reuse the existing confirmed binding, CloseCounterWindow and checkpoint policies; a gauge mean MUST NOT supply a counter baseline. Layout selection MUST NOT create a second same-scope same-kind owner.

#### Scenario: Complete meter with mixed values
- **WHEN** a meter has numeric voltage/current/power gauges, a text state and a confirmed cumulative kWh source
- **THEN** both layouts preserve gauge means, latest state/raw counter values and independently selected cumulative results with equivalent logical values and truthful types; no meter field is converted into another semantic kind by its name.

#### Scenario: New wide missing fields
- **WHEN** a new mean wide group uses its default optional partial policy and some or all members have no usable values
- **THEN** it writes the supported NULL columns and per-member evidence, including an all-missing window, without zero-fill; an explicitly required missing member still invokes the selected whole-row gate and records its skip cause.

#### Scenario: Reset and missing cumulative window
- **WHEN** the existing counter engine receives an unexplained decrease or observations spanning a missing cumulative window
- **THEN** both layouts preserve its NULL uncertainty or gap-labelled span total with actual bounds, do not average or abs the raw readings, and do not invent per-subwindow usage.

### Requirement: Durable average closure reuses delivery identity

Accepted observations SHALL use the existing durable journal. Closure SHALL atomically commit finalized average/counter rows, consumed input progress, checkpoints and outbox effects through the existing transaction. Restart SHALL replay the same accepted set; destination retries SHALL deliver frozen results and evidence without recomputation. Existing bounded backlog/backpressure and per-tick work limits SHALL remain in force. Unknown reductions, incomplete new-layout proof or unsupported restored metadata SHALL fail closed without a fallback writer.

#### Scenario: Closure failure and restart
- **WHEN** the closure transaction fails before commit or the process restarts after journal acknowledgement
- **THEN** recovery produces the same mean, counts and counter baseline once, no uncommitted input is consumed, and unrelated groups retain their results.

#### Scenario: Lost SQL acknowledgement
- **WHEN** a destination commits a finalized mean row but the reply is lost
- **THEN** resend uses the same payload digest and effect identity, SQL contains one row/effect, and sum/count and counter checkpoints do not advance again.
