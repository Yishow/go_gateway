## MODIFIED Requirements

### Requirement: Four-step intent-led setup
Operators MUST be able to save and navigate offline drafts with unresolved probes; draft navigation MUST NOT imply verified readiness, and only dependent reads or activation SHALL be blocked. The existing Studio V2 route SHALL guide device connection, acquisition points, Point-to-Tag mapping and grouped database output in four steps. Basic typed Tag writing MUST remain usable without required reporting, retention setup, aggregation or physical measurement semantics. Explicitly selected derived operations SHALL require their own confirmed semantics. The system MUST retain the Go backend and existing React framework and MUST NOT add a parallel V3 or restore a dedicated legacy /studio route.

#### Scenario: Mixed sensor configured without SQL
- **WHEN** an operator maps temperature, pressure, flow and counter source values as basic typed Tags
- **THEN** the operator can select those Tags for a managed write group without hand-written SQL or mandatory derived-value policies; choosing derived usage then requests the required confirmed roles.

#### Scenario: Reload and edited connection
- **WHEN** an operator reloads setup or changes a saved connection identity
- **THEN** the same persisted validity and revision checks determine progress, and the old probe does not keep the changed device ready.

#### Scenario: Unavailable protocol
- **WHEN** MQTT or another protocol lacks the complete V2 setup/parser/collector capability
- **THEN** the UI clearly marks it unavailable with a reason rather than letting the user enter an unsupported range flow.

#### Scenario: Same address on different devices
- **WHEN** two devices both expose address 40001
- **THEN** the address conflict checker keeps them separate while still detecting true overlapping ranges within the same device/address space.

#### Scenario: Missing advanced measurement
- **WHEN** a saved Point-to-Tag mapping has no persisted MeasurementDefinition
- **THEN** basic raw group setup proceeds without fake measurement IDs; only explicitly selected advanced semantics remain blocked.

#### Scenario: Offline draft navigation
- **WHEN** a device is offline or its saved configuration has no successful probe
- **THEN** the operator can save a draft and use Next/Back to configure later steps, including after reload, while the device remains visibly unverified and its dependent activation is blocked with a repair action.
