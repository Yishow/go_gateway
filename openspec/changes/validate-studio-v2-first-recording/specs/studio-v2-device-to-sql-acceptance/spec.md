## MODIFIED Requirements

### Requirement: Production device-to-SQL witness
The acceptance harness SHALL run the production cmd/test_ui binary with its embedded frontend, a controlled loopback device simulator and disposable SQLite and PostgreSQL targets. It MUST drive the actual Studio V2 UI, verify persisted device/point/tag/group identities and independently query resulting SQL rows. The basic witness SHALL start without recording tables and perform their preparation through confirmed product UI. Mock success, direct setup API mutation, fixture DDL or fixture INSERT MUST NOT substitute for the path under test.

#### Scenario: Mixed multi-device happy path
- **WHEN** the UI creates two simulated devices sharing addresses and a group with heterogeneous typed values
- **THEN** the persisted IDs and actual SQL payload, identity, UTC times and quality match the simulator fixture without address collisions or precision loss.

#### Scenario: Supported acquisition types and decimal boundary
- **WHEN** the device-to-SQL harness runs the existing acquisition types int16, uint16, bool, uint64, float32, float64 and string
- **THEN** all seven are verified through actual UI, acquisition and independent SQL queries, including exact uint64 9007199254740993.
- **AND** decimal evidence remains limited to the existing codec and SQL round-trip layer; device decimal acquisition is reported unsupported, not executed.

#### Scenario: Operation cleanup
- **WHEN** the UI confirms a test-write with a neighboring production fixture row
- **THEN** readback evidence matches the operation payload, cleanup removes only owned test data, and the production neighbor is unchanged.

#### Scenario: Empty destination becomes usable through UI
- **GIVEN** a fresh workspace, known loopback source points and no recording tables in the SQLite or PostgreSQL target
- **WHEN** the operator completes setup, explicitly confirms schema preparation and starts recording entirely through the UI
- **THEN** independent SQL observes the original source identities, exact values, UTC acquisition/bucket times and quality, and at least three consecutive production buckets are checked.

#### Scenario: One group contains multiple entities
- **WHEN** a single saved group has A and B entities with shared or distinct business columns
- **THEN** both entity rows are independently verified; two separate single-entity groups do not satisfy this scenario.

#### Scenario: Controlled first-use timing
- **GIVEN** the documented ready test environment, known connection/point data and the unchanged 60-second new-draft default
- **WHEN** each destination runs three fresh UI setups
- **THEN** every first-page-to-first-independent-SQL-row duration is recorded and MUST be at most 300 seconds to pass the controlled first-use criterion, without fixture-created recording tables, direct setup mutations or reduced correctness checks.

## ADDED Requirements

### Requirement: Review-regression witnesses cover actual recovery combinations
Acceptance SHALL exercise the reviewed lifecycle, deadline and entity failures through production services, with deterministic fault boundaries and durable/SQL assertions. Existing quality, capacity, uncertainty, cleanup and ownership safeguards MUST remain covered. An isolated context or race reproduction is supporting evidence, not a substitute for repository integration tests.

#### Scenario: Draft does not stop existing recording
- **WHEN** a running group's name or member draft is saved without Apply, two buckets pass and the gateway restarts
- **THEN** the old applied membership continues recording with no draft-induced gap.

#### Scenario: Offline cold start and historical journal
- **WHEN** the gateway restarts while the target is offline, and separate cases restart after ACK followed by Disable, Delete or supersession
- **THEN** eligible new input is locally durable and historical accepted input drains under its frozen version without new historical intake, redirection or duplicate effects.

#### Scenario: Slow result settlement and concurrent cutoff
- **WHEN** a destination attempt outlasts the settlement budget but remains within its own limit, or acquisition races with a cutoff
- **THEN** the result receives its full post-attempt local budget and the actual repository race tests report no unsynchronized cutoff access.

#### Scenario: Start recovery and unchanged acquisition time
- **WHEN** UI start or schema confirmation loses its response, or accepted rows are delivered after an outage
- **THEN** operation identity/progress remains recoverable, only the confirmed scope is affected, original acquisition times remain in SQL and UI distinguishes local acceptance from target commit.
