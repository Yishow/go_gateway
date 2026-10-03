## Purpose

Provide repeatable production device-to-SQL evidence for basic write groups, including deterministic fault boundaries and explicit limits on platform and field acceptance.

## ADDED Requirements

### Requirement: Production device-to-SQL witness
The acceptance harness SHALL run the production cmd/test_ui binary with its embedded frontend, a controlled loopback device simulator and disposable SQLite and PostgreSQL targets. It MUST drive the actual Studio V2 UI, verify persisted device/point/tag/group identities and independently query the resulting SQL rows. Mock API success or direct fixture INSERT SHALL NOT substitute for the end-to-end path.

#### Scenario: Mixed multi-device happy path
- **WHEN** the UI creates two simulated devices sharing addresses and a group with heterogeneous typed values
- **THEN** the persisted IDs and actual SQL payload, identity, UTC times and quality match the simulator fixture without address collisions or precision loss.

#### Scenario: Supported acquisition types and decimal boundary
- **WHEN** the device-to-SQL harness runs the existing acquisition types `int16`, `uint16`, `bool`, `uint64`, `float32`, `float64` and `string`
- **THEN** it SHALL verify all seven through the actual UI, acquisition and independent SQL queries, including exact `uint64` value `9007199254740993`.
- **AND** decimal evidence MUST remain explicitly limited to the opt-in codec and SQL round-trip layer; the witness MUST state that device decimal acquisition is unsupported and MUST NOT imply it was executed.

#### Scenario: Operation cleanup
- **WHEN** the UI confirms a test-write with a neighboring production fixture row
- **THEN** readback evidence matches the operation payload, cleanup removes only owned test data, and the production neighbor is unchanged.

### Requirement: Failure recovery and revision acceptance matrix
Acceptance MUST exercise missing/bad/stale/late values, database outage and restart, duplicate submission, lost commit acknowledgement, capacity exhaustion, poison data, edited settings and concurrent revisions with deterministic fault boundaries. Each case SHALL assert actual durable and SQL state as well as UI truth, and preserve a replayable witness.

#### Scenario: Outage recovery
- **WHEN** the database is disconnected across local acknowledgement and gateway restart, then restored
- **THEN** accepted records are recovered to their original target without duplicate effect keys and the UI transitions from pending to committed only with SQL evidence.

#### Scenario: Concurrent settings
- **WHEN** one client edits or replaces the destination while another saves/applies/tests an older revision
- **THEN** stale actions are rejected, existing backlog stays on its original destination and no successful state is fabricated.

#### Scenario: Incomplete bucket
- **WHEN** a required member is missing or bad at the closure boundary
- **THEN** the default group records a scoped incomplete outcome and emits no misleading complete row.

##### Example: deterministic production quality boundary
- **GIVEN** a real UI-created applied group, persisted source IDs, a loopback Modbus simulator and a disposable SQLite destination, with ten-second UTC buckets starting at `2026-01-01T00:00:00Z`
- **WHEN** the acceptance build captures actual reads at second `8` and then second `3`, releases them to the production pipeline, replays a duplicate and closes the bucket at second `10`
- **THEN** actual SQL preserves the second-8 values exactly once, the duplicate is a no-op, post-closure release is late with unchanged SQL, and missing/bad/stale/silent cases yield durable scoped skipped/no_data outcomes with no complete SQL row and matching UI truth.

### Requirement: Evidence distinguishes actual execution and field limits
Every reported pass SHALL include source/build identity, command, environment, fixture IDs, assertions and sanitized SQL/UI evidence. Unexecuted database, browser or operating-system checks MUST be reported as not run or blocked. Simulator tests MUST NOT imply live PLC, LAN, SCADA, Windows/ARM or field sign-off.

#### Scenario: PostgreSQL unavailable
- **WHEN** only SQLite can be exercised in the current validation environment
- **THEN** PostgreSQL is reported blocked/not run rather than passed by analogy.

#### Scenario: Linux simulator passes
- **WHEN** the simulator suite passes on Linux only
- **THEN** the result remains Linux simulator evidence and separately lists unexecuted platform and field checks.
