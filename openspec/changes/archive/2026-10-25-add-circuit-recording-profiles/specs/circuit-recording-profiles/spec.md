## Purpose

Provide circuit-oriented recording choices with independent instantaneous and cumulative schedules. Preserve both original counter readings and trustworthy interval differences through the existing canonical group and durable delivery lifecycle.

## ADDED Requirements

### Requirement: Stable circuit membership and independent profiles

Studio V2 SHALL persist explicit circuit membership using stable workspace-scoped identity, reusing confirmed equipment identity where available. It MUST NOT infer circuit membership or counter semantics from display names or addresses. Instantaneous and cumulative profiles SHALL be independently selectable for the same circuit, with one canonical group per circuit and profile kind. Unassigned raw Tags SHALL remain recordable without mandatory physical semantics. Duplicate create requests, reloads and label changes MUST reuse the canonical group. Two enabled profiles for one source MUST retain distinct output identities; duplicate ownership within the same destination and profile kind SHALL be rejected.

#### Scenario: Power and energy on one circuit
- **WHEN** a circuit contains a power Tag and a confirmed cumulative-energy Tag
- **THEN** the operator can select instantaneous power and cumulative energy together, with independent schedules and stable group identities after reload.

#### Scenario: One device contains several circuits
- **WHEN** an operator assigns selected source members from the same device to two named circuits
- **THEN** their memberships remain separate, renaming a circuit preserves its identity, and unselected sources are unchanged.

#### Scenario: Unknown semantics and duplicate requests
- **WHEN** a raw numeric Tag without confirmed counter semantics is selected for cumulative recording or the same create request is resent
- **THEN** cumulative activation reports counter_semantics_required without blocking independent raw recording, and a resent valid create request creates no duplicate writer.

### Requirement: Bounded recording intervals independent of acquisition

The backend and UI SHALL accept integer instantaneous intervals from 1 through 300 seconds, defaulting new instantaneous profiles to 15 seconds. They SHALL accept integer cumulative intervals from 1 through 600 minutes, defaulting new cumulative profiles to 10 minutes. API values SHALL include the profile-specific unit and reject missing or mismatched units, fractional, non-finite and out-of-range values. Existing saved intervals SHALL remain unchanged. Recording selection MUST NOT change device acquisition or transport flush frequency.

#### Scenario: Interval boundaries
- **WHEN** intervals are validated
- **THEN** instantaneous 1, 15 and 300 seconds and cumulative 1, 10 and 600 minutes are accepted; instantaneous 0 or 301, cumulative 0 or 601, fractional values and mismatched units return profile_interval_invalid.

#### Scenario: Slow acquisition
- **WHEN** collection produces fewer samples than the selected recording frequency
- **THEN** acquisition remains unchanged and uncovered recording windows expose missing data instead of duplicated old readings.

### Requirement: Deterministic per-source instantaneous records

The runtime SHALL close UTC half-open recording windows, choose the newest eligible observation within each window and retain actual observed and received times. Equal observation times SHALL use stable sample identity ordering. Every enabled member SHALL produce an independently identified long record: an exact typed value when valid or typed NULL with a missing, bad or stale reason. Values MUST NOT be carried from an older window. One missing member MUST NOT suppress another member's record. Existing wide-group completeness policies SHALL remain unchanged. Inputs received after the configured closure deadline SHALL be diagnosed without reopening a finalized profile window.

#### Scenario: Boundary and missing member
- **WHEN** a 15-second window ends at 00:00:15, a valid source sample has observed_at 00:00:15, and another member has no eligible sample
- **THEN** the boundary sample belongs to the next window and the missing member's current row contains NULL plus its reason; no previous-window value is substituted.

#### Scenario: Reordering and late arrivals
- **WHEN** the same accepted samples arrive in different orders before closure and an older sample arrives after closure
- **THEN** the finalized rows are identical for the accepted set and the late sample does not mutate or duplicate the closed rows.

### Requirement: Cumulative raw values and exact differences

Each cumulative record SHALL retain its latest eligible raw reading together with delta_value, actual delta_from and delta_to, sample provenance and a calculation status. Calculations SHALL inspect pre-decimation accepted inputs, preserve exact integer and decimal differences and avoid JavaScript Number coercion for exact values. Source floats SHALL retain their declared precision without claiming reconstructed precision. A first usable window SHALL store its raw value with NULL delta and baseline status, then checkpoint its last reliable reading. Later trustworthy differences SHALL reference comparable readings in the same semantic binding and series epoch.

#### Scenario: First reading and decimal progression
- **WHEN** successive usable windows end with exact readings 12000.0, 12000.4 and 12000.7 kWh
- **THEN** all raw readings are stored and the corresponding deltas are NULL, 0.4 and 0.3, with the first row marked baseline.

#### Scenario: Exact large integer
- **WHEN** a confirmed counter advances from 9007199254740992 to 9007199254740993
- **THEN** both readings round-trip exactly and the delta equals 1 in local state, API encoding and destination SQL.

### Requirement: Counter discontinuities and gaps remain explicit

An unexplained monotonic decrease, contradictory same-time readings or insufficient rollover evidence SHALL yield NULL delta and an uncertainty status while retaining raw values and events. Negative differences SHALL be accepted only for explicitly confirmed signed counters. Rollover compensation SHALL require a configured modulus, valid ranges and a rate/time bound sufficient to exclude multiple wraps. Missing windows SHALL retain NULL values and SHALL NOT advance a reliable baseline. A trustworthy difference across missing windows SHALL be retained as a gap-labelled span total over its actual observation interval, never allocated to unsupported subwindows. Changed epoch, source binding or profile revision SHALL start a new baseline. After an uncertain window, its last reliable raw reading SHALL establish the next baseline without counting the uncertain change.

#### Scenario: Reset inside a recording window
- **WHEN** a monotonic counter has accepted readings 100, 10 and 110 within the calculation span
- **THEN** raw 110 is retained, delta is NULL with reset_or_uncertain, and the decrease is not hidden by comparing only the endpoints.

#### Scenario: Signed and rollover evidence
- **WHEN** a confirmed signed counter changes from 10 to 7, or a monotonic counter changes from 65530 to 4 with modulus 65536
- **THEN** the signed delta is -3; rollover delta is 10 only with a verified single-wrap rate/time bound, otherwise it is NULL with uncertainty.

#### Scenario: Missing window and changed meter
- **WHEN** trustworthy readings 1000 and 1002 at 10:00 and 10:20 surround a missing 10-minute window, followed by a meter epoch change
- **THEN** the missing window remains NULL, the 2-unit span has its actual 10:00–10:20 bounds and gap status, and the new epoch's first raw reading has NULL delta and baseline status.

### Requirement: Atomic baseline and delivery progress

Accepted inputs SHALL use the existing durable journal. The runtime MUST atomically commit counter checkpoints, consumed input progress, finalized typed records and outbox intents with frozen revisions and destination. Recovery and SQL retries MUST reuse finalized payloads without recomputing consumption. Input identity conflicts SHALL fail visibly, and destination unavailability SHALL retain pending data under the existing bounded queue and backpressure contract.

#### Scenario: Crash at closure and uncertain SQL acknowledgement
- **WHEN** the runtime crashes before or after the closure transaction, or after SQL commit before receiving acknowledgement
- **THEN** recovery yields one committed destination effect per record identity, the baseline advances once, and no additional consumption is generated by retry.

#### Scenario: Offline restart
- **WHEN** the gateway restarts with a saved baseline and pending rows while the destination is offline
- **THEN** it restores the durable profile state and later delivers unchanged payloads to their frozen destination rather than resetting the counter or silently dropping pending records.
