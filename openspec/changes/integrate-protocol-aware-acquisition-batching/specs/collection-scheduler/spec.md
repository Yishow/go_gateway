## MODIFIED Requirements

### Requirement: Scheduled collection
The system SHALL schedule reads at configured polling-group intervals using the existing scheduler and immutable acquisition plans. Each group's due Point set SHALL remain independent of recording membership and other group cadences. Slow device completion SHALL NOT block admission of subsequent due scans for unrelated devices in that group.

#### Scenario: Interval read
- **WHEN** a polling group has a 5s interval
- **THEN** its Point scans become due every 5s, with measured overruns reported rather than fabricated on-time samples.

#### Scenario: Multiple groups on one device
- **WHEN** two groups become due while their device is busy
- **THEN** bounded pending work preserves both groups' latest due requirements without overwriting either Point set or creating an unbounded backlog.

### Requirement: Retry and backoff
The system SHALL enforce configured retry, deadline and cancellation limits across scheduler, adapter and transport attempts. Actual reconnect resends SHALL consume the same attempt budget. Only explicit unsupported-command rejection SHALL trigger one safe fallback plan per scan; malformed responses and timeouts SHALL NOT trigger unbounded scalar fanout retries.

#### Scenario: Retry after failure
- **WHEN** a read fails
- **THEN** total actual wire attempts obey the configured limit and deadline, and exhausted results are Bad without duplicate successful emissions.

#### Scenario: Safe capability fallback
- **WHEN** a configured random command receives an explicit unsupported-command response
- **THEN** one validated fallback plan is executed within the scan budget and the downgrade reason is exposed for that connection capability revision.

### Requirement: Concurrency control
The system SHALL serialize device operations and enforce bounded fair ready work, global active capacity and per-device/per-group coalesced pending state. Connection dialing and backoff SHALL NOT hold a global manager lock or an unrelated device's execution slot. Shared physical serial-port aliases SHALL either share proven serialization or be rejected before opening incompatible connections. Cancellation SHALL interrupt bounded transport waits and prevent new sends.

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

#### Scenario: Bad quality
- **WHEN** a read times out or a random response has invalid length, status or component encoding
- **THEN** affected Points are Bad with safe reasons, never Good zero or stale reused values.

#### Scenario: Same request and separate requests
- **WHEN** Points are obtained in one request
- **THEN** they share that request's gateway completion time but retain separate Point identities; Points from another request use its own completion time even within the same acquisition scan.

#### Scenario: Buffer pressure
- **WHEN** collector delivery capacity drops a result
- **THEN** counted loss and affected freshness or quality diagnostics are observable, and the drop is not reported as a newly accepted Good observation downstream.

#### Scenario: Manual acquisition parity
- **WHEN** manual polling requests selected Points
- **THEN** it uses the same immutable planner and bounded executor on exactly that requested subset, preserving cancellation, decoding and identity without creating duplicate scheduled work.
