## MODIFIED Requirements

### Requirement: Scheduled collection
The system SHALL schedule reads at configured polling-group intervals using the existing scheduler and immutable acquisition plans. Each group's due Point set SHALL remain independent of recording membership and other group cadences. Slow device completion SHALL NOT block admission of subsequent due scans for unrelated devices in that group.

Scheduled acquisition SHALL capture one coherent installed device/source/group snapshot before planning, without holding the scheduler lock during planning, connection acquisition or wire I/O. Its deterministic plan identity SHALL exclude acquisition time, scan/request identifiers, observation enablement and ephemeral connection tokens. Original PointRevision, immutable ReadProfile, InstallationLease/Scope, configuration fingerprints and frozen raw-to-mapped bindings SHALL retain their existing semantics; no SourceEpoch or durable sample identifier SHALL be introduced.


#### Scenario: Interval read
- **WHEN** a polling group has a 5s interval
- **THEN** its Point scans become due every 5s, with measured overruns reported rather than fabricated on-time samples.

#### Scenario: Multiple groups on one device
- **WHEN** two groups become due while their device is busy
- **THEN** bounded admitted pending work preserves both groups' latest due requirements without overwriting either Point set or creating an unbounded backlog; new entries rejected by an exhausted global capacity are explicitly counted instead of silently evicting another group's admitted work.

### Requirement: Retry and backoff
The system SHALL enforce configured retry, deadline and cancellation limits across scheduler, adapter and transport attempts. Actual reconnect resends SHALL consume the same attempt budget. Only explicit unsupported-command rejection SHALL trigger one safe fallback plan per scan; malformed responses and timeouts SHALL NOT trigger unbounded scalar fanout retries.

#### Scenario: Fixed scan deadline and actual-send admission
- **GIVEN** an installed snapshot, its original request timeout and a precomputed upper bound of initial plus at most one safe fallback plan's requests
- **WHEN** the scan begins
- **THEN** checked arithmetic establishes a fixed deadline no later than the caller deadline, and no retry, reconnect, queue wait or fallback extends it
- **AND** each request retains its original timeout and shared `1+retry-count` actual-send budget
- **AND** every actual send first checks cancellation and the captured install/config/stop eligibility at the existing wire-budget boundary; successful admission is the linearization point, so later disable cannot retract that admitted send but rejects subsequent retry/fallback sends.

#### Scenario: Explicit rejection and connection-scoped downgrade
- **WHEN** a complete matching random response explicitly rejects the command, including FATEK command 48 with error 9 rather than error 4
- **THEN** at most one same-scan safe fallback replans only unresolved Points within the original deadline and authorized attempt bound
- **AND** the downgrade SHALL be retained for the next scan only with a valid final-wire-attempt successful-connection token and the same capability revision
- **AND** failed or late connection completions cannot reset a newer token's downgrade, while timeouts, malformed responses and unrelated errors authorize neither fallback nor capability inference.


#### Scenario: Retry after failure
- **WHEN** a read fails
- **THEN** total actual wire attempts obey the configured limit and deadline, and exhausted results are Bad without duplicate successful emissions.

#### Scenario: Safe capability fallback
- **WHEN** a configured random command receives an explicit unsupported-command response
- **THEN** one validated fallback plan is executed within the scan budget and the downgrade reason is exposed for that connection capability revision.

### Requirement: Concurrency control
The system SHALL serialize device operations and enforce bounded fair ready work, global active capacity and per-device/per-group coalesced pending state. Connection dialing and backoff SHALL NOT hold a global manager lock or an unrelated device's execution slot. Shared physical serial-port aliases SHALL either share proven serialization or be rejected before opening incompatible connections. Cancellation SHALL interrupt bounded transport waits and prevent new sends.

Pending capacity SHALL count admitted, undispatched (device, polling-group) entries. Each configured group SHALL own at most one pending container, with at most one entry per device. Replacement coalesces only the same tuple without consuming another slot or modifying a running scan's captured set/deadline. Active execution permits SHALL exclude device-gate waits and retry backoff; device scan/gate ownership SHALL remain serialized during backoff. Resuming scans and newly ready devices SHALL share fair capacity admission. Waiting/executing scheduled-scan goroutines SHALL be bounded by configured device count, never tick or pending-entry count. Unspecified capacity settings SHALL preserve existing constructor/default compatibility; explicitly invalid values SHALL be rejected rather than silently treated as usable capacity. Blocking permit acquisition SHALL NOT run inside nonblocking actual-send admission or observer callbacks. Task 3.2 applies this to scheduled work, retaining the existing manual device gate; task 3.3 completes manual shared-executor capacity behavior.

#### Scenario: Backoff yields active capacity without yielding device ownership
- **WHEN** global active capacity is one and device A enters scheduler or adapter retry backoff while device B is ready
- **THEN** B can acquire available active capacity, while another group or manual operation cannot interleave commands on A
- **AND** A reacquires fairly within its original deadline and retry budget, or cancellation prevents another send.


#### Scenario: Serial device lock
- **WHEN** two tasks target the same serial device
- **THEN** commands are serialized and unsupported shared-port aliases are rejected instead of opening parallel conflicting sessions.

#### Scenario: Slow device and queue saturation
- **WHEN** one device times out and others remain ready
- **THEN** available execution capacity visits ready devices fairly, every blocking read is deadline-bounded, pending work never exceeds configured bounds, and missed/coalesced scans are counted.

#### Scenario: Connection revision changes
- **WHEN** protocol or connection configuration changes during an in-flight read
- **THEN** the manager does not reuse the old connection as the new revision and late old results cannot refresh new-revision state.

### Requirement: Quality and timestamp
The system SHALL emit at most one final result per Point per scan, with immutable Point identity, effective decoder, configuration fingerprint and acquisition lease. Complete validated raw responses SHALL be scattered per Point; corrupt alignment SHALL invalidate every affected Point, and a locally decodable Point failure SHALL NOT invalidate unrelated successful requests. Without trusted source time, completion of the actual successful read SHALL determine observation time. Shared batch identity SHALL NOT imply an atomic PLC snapshot.

The latest scheduled scan evidence SHALL be bounded by configured device/group tuples and retain original captured identities. It SHALL not become a history registry or change acquisition, retry, emission, or sample-identity behavior. Overrun SHALL describe only the matching scan's captured-due elapsed time relative to its cadence; global due counters SHALL not be projected as scoped counters. Output-loss projection SHALL reuse existing counters, expose only provably matching retained detail as partial, and never attribute global evicted or unattributed loss to the requested scope.

#### Scenario: Final request evidence preserves durable identity
- **WHEN** retries or fallback finish a scheduled scan
- **THEN** each Point emits at most once with its captured binding and the actual final request completion time
- **AND** request ID and plan fingerprint are transient final collector/observation evidence, not durable SampleEnvelope or StableSampleID inputs
- **AND** stale installation/revision results cannot refresh replacement configuration or be rebound to its mappings, and different requests do not claim one atomic PLC snapshot.

#### Scenario: Bad quality
- **WHEN** a read times out or a random response has invalid length, status or component encoding
- **THEN** affected Points are Bad with safe reasons, never Good zero or stale reused values.

#### Scenario: Same request and separate requests
- **WHEN** Points are obtained in one request
- **THEN** they share that request's gateway completion time but retain separate Point identities; Points from another request use its own completion time even within the same acquisition scan.

#### Scenario: Buffer pressure
- **WHEN** collector delivery capacity drops a result
- **THEN** counted loss and affected freshness or quality diagnostics are observable, and the drop is not reported as a newly accepted Good observation downstream.
- **AND** output-result losses are counted separately from due-admission rejection using the dropped value's captured scope/revision/Point identity, never a later lookup; missing legacy identity is explicitly unattributed
- **AND** cumulative loss remains available while recent identity detail is bounded to at most 256 series, with evicted-detail loss explicitly aggregated rather than represented as a complete history. These native diagnostics are implemented once in task 3.2 and consumed by task 4.1 without adding a public diagnostics route in task 3.2.


#### Scenario: Manual acquisition parity
- **WHEN** manual polling requests selected Points
- **THEN** it uses the same immutable planner and bounded executor on exactly that requested subset, preserving cancellation, decoding and identity without creating duplicate scheduled work.

#### Scenario: Manual return-only and captured identity
- **WHEN** manual polling selects installed, non-recovery Points, including multiple original polling groups on one device or repeated IDs
- **THEN** one coherent device/config/Point/profile/lease snapshot preserves the original group context, one acquisition ID per manual device intent and unique request identifiers within the shared plan
- **AND** each canonical Point has at most one final result returned through the existing caller/observer path, with no extra emission to the scheduled ValueChannel or changes to unrequested Points
- **AND** unknown, recovery-pending and empty selections retain their existing exclusion behavior; cancellation of known requested Points retains the existing Bad return contract.

#### Scenario: Manual capacity before and across scheduled runs
- **WHEN** explicit manual polling overlaps Start, scheduled Stop or restart
- **THEN** a constructor-owned scheduler-lifetime active permit pool and the existing device gates enforce the same global capacity without resetting active manual ownership or requiring scheduled activation
- **AND** scheduled Stop only cancels and joins its own workers while the original caller context owns manual cancellation; no new Start-busy condition or duplicate scheduled work is introduced
- **AND** invalid explicit capacity fails closed rather than waiting forever, and legacy non-bounded custom protocols retain their single scheduler-level scalar invocation and original adapter budget behavior.
