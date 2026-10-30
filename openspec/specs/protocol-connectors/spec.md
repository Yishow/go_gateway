# protocol-connectors Specification

## Purpose
TBD - created by archiving change add-device-data-pipeline. Update Purpose after archive.

## Requirements

### Requirement: Protocol connector registry
The system SHALL provide a common protocol interface and a registry to load protocol connectors by type.

#### Scenario: Load protocol
- WHEN a device is configured with a protocol type
- THEN the system instantiates the matching connector from the registry

---
### Requirement: Modbus connectors
The system SHALL support Modbus TCP, Modbus RTU, and Modbus UDP for reading and writing registers.

#### Scenario: Read Modbus holding register
- WHEN a Modbus TCP device is polled for a holding register
- THEN the connector returns the raw value with a timestamp

---
### Requirement: FATEK FBs connector
The system SHALL support FATEK FBs protocol with ASCII frames and LRC validation.

#### Scenario: Read FATEK device
- WHEN a FATEK device is polled
- THEN the connector validates LRC and returns a decoded value

---
### Requirement: Mitsubishi MC 3E connector
The system SHALL support Mitsubishi MC Protocol 3E frame communication over TCP.

#### Scenario: Read MC 3E point
- WHEN an MC 3E device is polled
- THEN the connector returns the decoded value for the address

---
### Requirement: MQTT ingest connector
The system SHALL support MQTT subscription as a data source, mapping topics and payloads into point values.

#### Scenario: Ingest MQTT topic
- WHEN a subscribed topic receives a payload
- THEN the system emits a point value for mapping and storage

---
### Requirement: Persistent connection mode

The system SHALL support a persistent (long) connection mode for protocol connectors to reduce handshake overhead during high-frequency polling.

#### Scenario: Enable persistent connection

- GIVEN a protocol connector (e.g., Modbus TCP) is configured with `persistent_connection: true`
- WHEN the system polls the device repeatedly
- THEN the connector reuses the same TCP connection without re-handshaking

#### Scenario: Auto-reconnect on failure

- GIVEN a persistent connection is active
- WHEN the network connection is interrupted
- THEN the connector automatically attempts to reconnect on the next poll

---
### Requirement: Connection health check

The system SHALL provide a heartbeat/keep-alive mechanism for persistent connections to detect stale sockets.

#### Scenario: Keep-alive timeout

- GIVEN a persistent connection has been idle for longer than the configured `keep_alive_interval`
- WHEN no data is exchanged
- THEN the connector sends a health-check packet or closes the connection preemptively

---
### Requirement: Connection tests return phase-specific diagnostics
The system SHALL report transport-connect results and protocol-probe results as separate diagnostics for connector testing.

#### Scenario: Connect-stage failure is reported separately
- **WHEN** a connector test fails before protocol exchange because the socket cannot be established
- **THEN** the system SHALL return a connect-stage failure classification
- **AND** SHALL include the underlying network error detail

#### Scenario: Probe-stage failure is reported separately
- **WHEN** the socket connection succeeds but protocol validation or probe read fails
- **THEN** the system SHALL return a probe-stage failure classification
- **AND** SHALL preserve the successful connect result in the same test response

---
### Requirement: Connector diagnostics expose planning capability hints
Protocol connector diagnostics SHALL expose the capability hints needed by rule-driven planning and later review flows.

#### Scenario: Planning capability hints are returned with diagnostics
- **WHEN** a client requests connector diagnostics for a selected device
- **THEN** the system returns connect/probe phase results plus protocol capability hints relevant to source-rule planning
- **AND** those hints remain distinguishable from transient test failure messages

---
### Requirement: Connection probes use protocol-appropriate configuration
The system SHALL use protocol-appropriate probe configuration instead of relying on one implicit default probe behavior for all connectors.

#### Scenario: Protocol-specific probe configuration is applied
- **WHEN** a connector test runs for a configured device
- **THEN** the system SHALL use the protocol-specific probe configuration for that device
- **AND** SHALL NOT assume one shared implicit probe address or validation pattern across all protocols

---
### Requirement: Multi-register decode byte order is configurable on Modbus-class connections
The system SHALL expose **`data_format`** in the **JSON ConfigSchema** (and persist the value in device `connection_config`) for **Modbus TCP**, **Modbus RTU**, and **Modbus UDP** protocols, using the same enumeration and semantics as existing multi-register decoding (`ABCD`, `BADC`, `CDAB`, `DCBA`) and the shared decoding library (for example `lib/hsllogic`). The Modbus connector adapters SHALL read this field when decoding multi-register integers and IEEE-754 floats.

#### Scenario: Modbus TCP device saves data_format
- **WHEN** an operator configures a Modbus TCP device with a non-default `data_format`
- **THEN** the stored configuration SHALL include the selected `data_format`
- **AND** subsequent reads of multi-register data types SHALL use that ordering

#### Scenario: Schema parity across Modbus variants
- **WHEN** the UI lists Modbus TCP, RTU, or UDP
- **THEN** each protocol’s connection schema SHALL include `data_format` with the same enum values and documented defaults

---
### Requirement: MC 3E retains data_format with documented default
The system SHALL continue to support **`data_format`** on **Mitsubishi MC 3E** connections. The default SHALL remain consistent with existing manufacturer-oriented behavior unless a breaking migration is explicitly approved.

#### Scenario: MC 3E connection still accepts data_format
- **WHEN** an operator configures MC 3E with `data_format` set
- **THEN** the connector SHALL apply that ordering to float and multi-word decode paths

---
### Requirement: Protocols without multi-register Modbus-style blocks do not require data_format
Protocols that do not use the same holding-register multi-word decode path (for example **MQTT** ingest as configured today) SHALL NOT be required to present `data_format` in their connection schema unless a future capability adds equivalent multi-word semantics.

#### Scenario: MQTT connection omits data_format
- **WHEN** an operator configures an MQTT device per current schema
- **THEN** the absence of `data_format` SHALL not block save

---
### Requirement: Scalar point reads preserve protocol count units

The system SHALL interpret single-point requests as one typed value for MC 3E, as 16-bit register counts for Modbus, and as native component counts for FATEK. FATEK scalar requests SHALL use the addressed component width to request enough components for the requested data type. Direct MC typed-value batch counts and direct FATEK native-component counts SHALL retain their existing meanings.

#### Scenario: Multiword MC point across entrypoints

- **WHEN** periodic, manual or read-now acquisition reads one uint32 point at D3 with words 34464 and 1 in CDAB order
- **THEN** it SHALL request two words and return scalar 100000 rather than an array; a single 64-bit value SHALL occupy four words and retain exact encoding.

#### Scenario: Direct batch and other protocols remain unchanged

- **WHEN** a direct MC batch requests two uint32 values, a direct FATEK request specifies a component count, or a Modbus point is polled
- **THEN** MC SHALL retain two-value batch decoding, FATEK SHALL transmit the specified native-component count, and Modbus SHALL retain its existing 16-bit register counts.

#### Scenario: FATEK scalar count follows native width across entrypoints

- **WHEN** periodic, manual or read-now acquisition reads a FATEK scalar
- **THEN** 32-bit scalars SHALL request two components in a 16-bit area and one component in a 32-bit area
- **AND** 64-bit scalars SHALL request four components in a 16-bit area and two components in a 32-bit area
- **AND** 16-bit scalars and discrete bool points SHALL retain a single-component request.

#### Scenario: Invalid FATEK scalar address does not produce a request

- **WHEN** the scalar count helper cannot resolve an existing supported FATEK component and its width
- **THEN** the caller SHALL report a read failure before sending the request
- **AND** SHALL NOT guess a successful count or publish a quality-good sample.


<!-- @trace
source: fix-plc-read-decoding-and-address-consistency
updated: 2026-10-30
code:
  - lib/hsllogic/register_conversion.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/datalink/device/service_point_read.go
  - frontend/src/i18n/locales/en/workbench-v2-plc.json
  - frontend/src/pages/datalink/workbench-v2/useStudioV2PersistedMCPoints.ts
  - lib/hsllogic/bytetransform.go
  - internal/datalink/collector/scheduler_polling.go
  - internal/datalink/sourcerule/service.go
  - internal/datalink/sourcerule/runtime_reconcile.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/ReadinessIssuesWarning.tsx
  - internal/datalink/sourcerule/mc_address_plan.go
  - internal/datalink/collector/scheduler_manual_poll.go
  - internal/datalink/connector/adapters/fatek_helpers.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/pages/datalink/workbench-v2/studioV2LocatorHydration.ts
  - frontend/src/features/datalink/workbench-v2/state/nextFreeStartAddress.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - internal/datalink/point/service_batch_validation.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2-plc.json
  - internal/datalink/connector/point_read.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/state/mcPersistedPoints.ts
  - internal/datalink/connector/adapters/fatek_operations.go
  - internal/datalink/sourcerule/protocol_address_planner.go
  - frontend/src/i18n/config.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - internal/protocol/fatek/client.go
  - frontend/src/utils/addressParser.ts
  - internal/datalink/connector/adapters/modbus_parse_helpers.go
  - internal/protocol/fatek/address_parse.go
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
tests:
  - internal/datalink/sourcerule/mc_validation_boundary_test.go
  - internal/datalink/sourcerule/mc_address_planner_regression_test.go
  - internal/protocol/fatek/status_payload_test.go
  - internal/datalink/collector/point_count_contract_test.go
  - internal/datalink/connector/adapters/fatek_refactor_test.go
  - internal/datalink/connector/adapters/modbus_numeric_bounds_test.go
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/mc-persisted-draft-preservation.test.ts
  - frontend/tests/unit/workbench-v2/mc-persisted-hydration.test.tsx
  - frontend/tests/unit/workbench-v2/studio-locator-mc-hydration.test.tsx
  - frontend/tests/unit/workbench-v2/mc-address-bounds.test.tsx
  - internal/datalink/collector/fatek_status_acquisition_test.go
  - frontend/tests/unit/utils/addressParser.test.ts
  - frontend/tests/unit/workbench-v2/mc-next-free-address.test.ts
  - frontend/tests/unit/workbench-v2/mc-persisted-hydration.fixtures.ts
  - frontend/tests/unit/workbench-v2/occupied-range.test.tsx
  - internal/datalink/sourcerule/protocol_address_test.go
  - internal/datalink/connector/adapters/plc_fatek_decoding_regression_test.go
  - internal/datalink/collector/fatek_point_count_contract_test.go
  - internal/datalink/sourcerule/mc_address_bounds_test.go
  - frontend/tests/unit/workbench-v2/mc-persisted-view.test.tsx
  - internal/datalink/connector/adapters/plc_fatek_status_regression_test.go
  - internal/datalink/connector/adapters/plc_modbus_address_regression_test.go
  - internal/datalink/sourcerule/mc_radix_compatibility_test.go
  - frontend/tests/unit/workbench-v2/mc-persisted-mapping-queue.test.tsx
  - internal/datalink/connector/point_read_test.go
  - internal/datalink/point/mc_hex_address_test.go
  - internal/datalink/sourcerule/mc_radix_runtime_reconcile_test.go
  - internal/datalink/connector/adapters/plc_read_decoding_regression_test.go
  - lib/hsllogic/plc_byte_order_regression_test.go
  - frontend/tests/unit/workbench-v2/mc-persisted-readiness-regression.test.ts
-->

---
### Requirement: Multiword numeric decoding follows all configured byte orders

The shared decoder used by Modbus-class and MC 3E connectors SHALL decode int32, uint32, float32, int64, uint64 and float64 according to the configured ABCD, BADC, CDAB or DCBA order. BADC SHALL swap bytes inside each word, CDAB SHALL reverse word order across the scalar, and DCBA SHALL perform both operations. Signed values and IEEE-754 values SHALL preserve their complete bit patterns. Existing connection defaults, point override precedence and single-word int16/uint16 decoding SHALL remain unchanged.

#### Scenario: Four byte orders decode the same float32

- **WHEN** the shared decoder or a Modbus/MC connector reads a float32 using the word sequence for its configured order
- **THEN** each sequence in the following table SHALL produce typed float32 1.

##### Example: Float32 word sequences

| Order | Words in response order | Expected |
| --- | --- | --- |
| ABCD | 3F80 0000 | float32 1 |
| BADC | 803F 0000 | float32 1 |
| CDAB | 0000 3F80 | float32 1 |
| DCBA | 0000 803F | float32 1 |

#### Scenario: Four byte orders preserve a complete uint64

- **WHEN** a 64-bit scalar is decoded using the word sequence for its configured order
- **THEN** each sequence in the following table SHALL produce uint64 0x0102030405060708 without conversion through floating point.

##### Example: Uint64 word sequences

| Order | Words in response order |
| --- | --- |
| ABCD | 0102 0304 0506 0708 |
| BADC | 0201 0403 0605 0807 |
| CDAB | 0708 0506 0304 0102 |
| DCBA | 0807 0605 0403 0201 |

#### Scenario: Signed and float64 values retain bits

- **WHEN** each configured order decodes negative signed 32/64-bit values or float64 with canonical bits 3FF0000000000000
- **THEN** signed values SHALL retain their two's-complement bit pattern and float64 SHALL equal 1
- **AND** existing ABCD/CDAB values and single-word results SHALL remain unchanged.


<!-- @trace
source: fix-plc-read-decoding-and-address-consistency
updated: 2026-10-30
code:
  - lib/hsllogic/register_conversion.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/datalink/device/service_point_read.go
  - frontend/src/i18n/locales/en/workbench-v2-plc.json
  - frontend/src/pages/datalink/workbench-v2/useStudioV2PersistedMCPoints.ts
  - lib/hsllogic/bytetransform.go
  - internal/datalink/collector/scheduler_polling.go
  - internal/datalink/sourcerule/service.go
  - internal/datalink/sourcerule/runtime_reconcile.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/ReadinessIssuesWarning.tsx
  - internal/datalink/sourcerule/mc_address_plan.go
  - internal/datalink/collector/scheduler_manual_poll.go
  - internal/datalink/connector/adapters/fatek_helpers.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/pages/datalink/workbench-v2/studioV2LocatorHydration.ts
  - frontend/src/features/datalink/workbench-v2/state/nextFreeStartAddress.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - internal/datalink/point/service_batch_validation.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2-plc.json
  - internal/datalink/connector/point_read.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/state/mcPersistedPoints.ts
  - internal/datalink/connector/adapters/fatek_operations.go
  - internal/datalink/sourcerule/protocol_address_planner.go
  - frontend/src/i18n/config.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - internal/protocol/fatek/client.go
  - frontend/src/utils/addressParser.ts
  - internal/datalink/connector/adapters/modbus_parse_helpers.go
  - internal/protocol/fatek/address_parse.go
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
tests:
  - internal/datalink/sourcerule/mc_validation_boundary_test.go
  - internal/datalink/sourcerule/mc_address_planner_regression_test.go
  - internal/protocol/fatek/status_payload_test.go
  - internal/datalink/collector/point_count_contract_test.go
  - internal/datalink/connector/adapters/fatek_refactor_test.go
  - internal/datalink/connector/adapters/modbus_numeric_bounds_test.go
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/mc-persisted-draft-preservation.test.ts
  - frontend/tests/unit/workbench-v2/mc-persisted-hydration.test.tsx
  - frontend/tests/unit/workbench-v2/studio-locator-mc-hydration.test.tsx
  - frontend/tests/unit/workbench-v2/mc-address-bounds.test.tsx
  - internal/datalink/collector/fatek_status_acquisition_test.go
  - frontend/tests/unit/utils/addressParser.test.ts
  - frontend/tests/unit/workbench-v2/mc-next-free-address.test.ts
  - frontend/tests/unit/workbench-v2/mc-persisted-hydration.fixtures.ts
  - frontend/tests/unit/workbench-v2/occupied-range.test.tsx
  - internal/datalink/sourcerule/protocol_address_test.go
  - internal/datalink/connector/adapters/plc_fatek_decoding_regression_test.go
  - internal/datalink/collector/fatek_point_count_contract_test.go
  - internal/datalink/sourcerule/mc_address_bounds_test.go
  - frontend/tests/unit/workbench-v2/mc-persisted-view.test.tsx
  - internal/datalink/connector/adapters/plc_fatek_status_regression_test.go
  - internal/datalink/connector/adapters/plc_modbus_address_regression_test.go
  - internal/datalink/sourcerule/mc_radix_compatibility_test.go
  - frontend/tests/unit/workbench-v2/mc-persisted-mapping-queue.test.tsx
  - internal/datalink/connector/point_read_test.go
  - internal/datalink/point/mc_hex_address_test.go
  - internal/datalink/sourcerule/mc_radix_runtime_reconcile_test.go
  - internal/datalink/connector/adapters/plc_read_decoding_regression_test.go
  - lib/hsllogic/plc_byte_order_regression_test.go
  - frontend/tests/unit/workbench-v2/mc-persisted-readiness-regression.test.ts
-->

---
### Requirement: FATEK typed register decoding preserves complete values

The FATEK adapter SHALL assemble typed numeric scalars from native 16-bit or 32-bit components using its existing high-word-first integer convention. It SHALL reinterpret floating-point bits as IEEE-754, preserve all 64 integer bits, and return the requested Go numeric type. Insufficient data and unsupported conversions SHALL return an error with bad quality and no successful value rather than returning the first component. This requirement SHALL NOT introduce a configurable FATEK byte order or change write encoding.

#### Scenario: Float32 has the same bits in word and doubleword areas

- **WHEN** D0 returns components 3F80 and 0000, or DR0 returns component 3F800000, for a float32 request
- **THEN** both reads SHALL return typed float32 1
- **AND** neither read SHALL numerically cast the encoded integer into a float.

#### Scenario: Complete 64-bit values in both native widths

- **WHEN** D0 returns 0102 0304 0506 0708, or DR0 returns 01020304 05060708, for a uint64 request
- **THEN** both reads SHALL return typed uint64 0x0102030405060708
- **AND** int64 and float64 requests SHALL use the complete corresponding bit pattern.

#### Scenario: A typed scalar cannot be completed

- **WHEN** a FATEK float32 request has fewer than two 16-bit components, or a 64-bit request has fewer than four 16-bit or two 32-bit components
- **THEN** the read SHALL fail with bad quality and no successful scalar
- **AND** no missing bits SHALL be fabricated from zero or a first-component fallback.


<!-- @trace
source: fix-plc-read-decoding-and-address-consistency
updated: 2026-10-30
code:
  - lib/hsllogic/register_conversion.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/datalink/device/service_point_read.go
  - frontend/src/i18n/locales/en/workbench-v2-plc.json
  - frontend/src/pages/datalink/workbench-v2/useStudioV2PersistedMCPoints.ts
  - lib/hsllogic/bytetransform.go
  - internal/datalink/collector/scheduler_polling.go
  - internal/datalink/sourcerule/service.go
  - internal/datalink/sourcerule/runtime_reconcile.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/ReadinessIssuesWarning.tsx
  - internal/datalink/sourcerule/mc_address_plan.go
  - internal/datalink/collector/scheduler_manual_poll.go
  - internal/datalink/connector/adapters/fatek_helpers.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/pages/datalink/workbench-v2/studioV2LocatorHydration.ts
  - frontend/src/features/datalink/workbench-v2/state/nextFreeStartAddress.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - internal/datalink/point/service_batch_validation.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2-plc.json
  - internal/datalink/connector/point_read.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/state/mcPersistedPoints.ts
  - internal/datalink/connector/adapters/fatek_operations.go
  - internal/datalink/sourcerule/protocol_address_planner.go
  - frontend/src/i18n/config.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - internal/protocol/fatek/client.go
  - frontend/src/utils/addressParser.ts
  - internal/datalink/connector/adapters/modbus_parse_helpers.go
  - internal/protocol/fatek/address_parse.go
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
tests:
  - internal/datalink/sourcerule/mc_validation_boundary_test.go
  - internal/datalink/sourcerule/mc_address_planner_regression_test.go
  - internal/protocol/fatek/status_payload_test.go
  - internal/datalink/collector/point_count_contract_test.go
  - internal/datalink/connector/adapters/fatek_refactor_test.go
  - internal/datalink/connector/adapters/modbus_numeric_bounds_test.go
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/mc-persisted-draft-preservation.test.ts
  - frontend/tests/unit/workbench-v2/mc-persisted-hydration.test.tsx
  - frontend/tests/unit/workbench-v2/studio-locator-mc-hydration.test.tsx
  - frontend/tests/unit/workbench-v2/mc-address-bounds.test.tsx
  - internal/datalink/collector/fatek_status_acquisition_test.go
  - frontend/tests/unit/utils/addressParser.test.ts
  - frontend/tests/unit/workbench-v2/mc-next-free-address.test.ts
  - frontend/tests/unit/workbench-v2/mc-persisted-hydration.fixtures.ts
  - frontend/tests/unit/workbench-v2/occupied-range.test.tsx
  - internal/datalink/sourcerule/protocol_address_test.go
  - internal/datalink/connector/adapters/plc_fatek_decoding_regression_test.go
  - internal/datalink/collector/fatek_point_count_contract_test.go
  - internal/datalink/sourcerule/mc_address_bounds_test.go
  - frontend/tests/unit/workbench-v2/mc-persisted-view.test.tsx
  - internal/datalink/connector/adapters/plc_fatek_status_regression_test.go
  - internal/datalink/connector/adapters/plc_modbus_address_regression_test.go
  - internal/datalink/sourcerule/mc_radix_compatibility_test.go
  - frontend/tests/unit/workbench-v2/mc-persisted-mapping-queue.test.tsx
  - internal/datalink/connector/point_read_test.go
  - internal/datalink/point/mc_hex_address_test.go
  - internal/datalink/sourcerule/mc_radix_runtime_reconcile_test.go
  - internal/datalink/connector/adapters/plc_read_decoding_regression_test.go
  - lib/hsllogic/plc_byte_order_regression_test.go
  - frontend/tests/unit/workbench-v2/mc-persisted-readiness-regression.test.ts
-->

---
### Requirement: Explicit raw Modbus addresses preserve requested function

For Modbus TCP, RTU and UDP, the parser SHALL preserve the existing HR/IR/C/CO/DI aliases and recognized five-digit traditional area addresses. Aliases SHALL keep zero-based offsets; traditional addresses SHALL use a leading 0/1/3/4 followed by a four-digit number from 0001 through 9999 and SHALL subtract one from that number. Their explicit prefixes SHALL retain priority over Function. For any other decimal numeric address with an explicit Function, the parser SHALL use that Function and the complete numeric address as a zero-based offset in 0 through 65535. It SHALL reject negative, malformed and overflowing raw offsets before wire access. The same classification SHALL apply to read and write callers of the shared parser.

#### Scenario: Explicit raw holding-register request stays a holding-register request

- **WHEN** Address is 100 and Function is 03
- **THEN** the wire request SHALL use FC03 and offset 100
- **AND** SHALL NOT use FC01 or offset 99.

#### Scenario: Raw boundary offsets are preserved

- **WHEN** explicit Function 04 is used with numeric Address 0, 65535 or 65536
- **THEN** the first two requests SHALL use FC04 at offsets 0 and 65535 respectively
- **AND** Address 65536 SHALL fail before wire access.

#### Scenario: Traditional and alias addresses retain prefix semantics

- **WHEN** a caller uses the following recognized addresses
- **THEN** the parser SHALL produce the listed function and offset even when Function contains a different default value.

##### Example: Preserved prefixed addresses

| Address | Function code | Offset |
| --- | --- | ---: |
| 00001 | 01 | 0 |
| 10001 | 02 | 0 |
| 30001 | 04 | 0 |
| 40001 | 03 | 0 |
| HR100 | 03 | 100 |
| IR100 | 04 | 100 |
| C100 or CO100 | 01 | 100 |
| DI100 | 02 | 100 |

#### Scenario: No-function fallback remains compatible

- **WHEN** a direct caller sends nontraditional numeric Address 100 without Function
- **THEN** the parser SHALL retain its historical FC01, offset 99 fallback
- **AND** other no-function cases SHALL retain their existing default and range checks.

#### Scenario: Persisted short numeric points expose the corrected precedence

- **GIVEN** a persisted point has Address 100 and stored Function 03
- **WHEN** the corrected reader polls it
- **THEN** it SHALL read FC03 offset 100 and SHALL leave stored Address, Function, point ID and mappings unchanged
- **AND** release evidence SHALL identify this wire-target change as breaking rather than claiming physical-target compatibility.


<!-- @trace
source: fix-plc-read-decoding-and-address-consistency
updated: 2026-10-30
code:
  - lib/hsllogic/register_conversion.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/datalink/device/service_point_read.go
  - frontend/src/i18n/locales/en/workbench-v2-plc.json
  - frontend/src/pages/datalink/workbench-v2/useStudioV2PersistedMCPoints.ts
  - lib/hsllogic/bytetransform.go
  - internal/datalink/collector/scheduler_polling.go
  - internal/datalink/sourcerule/service.go
  - internal/datalink/sourcerule/runtime_reconcile.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/ReadinessIssuesWarning.tsx
  - internal/datalink/sourcerule/mc_address_plan.go
  - internal/datalink/collector/scheduler_manual_poll.go
  - internal/datalink/connector/adapters/fatek_helpers.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/pages/datalink/workbench-v2/studioV2LocatorHydration.ts
  - frontend/src/features/datalink/workbench-v2/state/nextFreeStartAddress.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - internal/datalink/point/service_batch_validation.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2-plc.json
  - internal/datalink/connector/point_read.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/state/mcPersistedPoints.ts
  - internal/datalink/connector/adapters/fatek_operations.go
  - internal/datalink/sourcerule/protocol_address_planner.go
  - frontend/src/i18n/config.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - internal/protocol/fatek/client.go
  - frontend/src/utils/addressParser.ts
  - internal/datalink/connector/adapters/modbus_parse_helpers.go
  - internal/protocol/fatek/address_parse.go
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
tests:
  - internal/datalink/sourcerule/mc_validation_boundary_test.go
  - internal/datalink/sourcerule/mc_address_planner_regression_test.go
  - internal/protocol/fatek/status_payload_test.go
  - internal/datalink/collector/point_count_contract_test.go
  - internal/datalink/connector/adapters/fatek_refactor_test.go
  - internal/datalink/connector/adapters/modbus_numeric_bounds_test.go
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/mc-persisted-draft-preservation.test.ts
  - frontend/tests/unit/workbench-v2/mc-persisted-hydration.test.tsx
  - frontend/tests/unit/workbench-v2/studio-locator-mc-hydration.test.tsx
  - frontend/tests/unit/workbench-v2/mc-address-bounds.test.tsx
  - internal/datalink/collector/fatek_status_acquisition_test.go
  - frontend/tests/unit/utils/addressParser.test.ts
  - frontend/tests/unit/workbench-v2/mc-next-free-address.test.ts
  - frontend/tests/unit/workbench-v2/mc-persisted-hydration.fixtures.ts
  - frontend/tests/unit/workbench-v2/occupied-range.test.tsx
  - internal/datalink/sourcerule/protocol_address_test.go
  - internal/datalink/connector/adapters/plc_fatek_decoding_regression_test.go
  - internal/datalink/collector/fatek_point_count_contract_test.go
  - internal/datalink/sourcerule/mc_address_bounds_test.go
  - frontend/tests/unit/workbench-v2/mc-persisted-view.test.tsx
  - internal/datalink/connector/adapters/plc_fatek_status_regression_test.go
  - internal/datalink/connector/adapters/plc_modbus_address_regression_test.go
  - internal/datalink/sourcerule/mc_radix_compatibility_test.go
  - frontend/tests/unit/workbench-v2/mc-persisted-mapping-queue.test.tsx
  - internal/datalink/connector/point_read_test.go
  - internal/datalink/point/mc_hex_address_test.go
  - internal/datalink/sourcerule/mc_radix_runtime_reconcile_test.go
  - internal/datalink/connector/adapters/plc_read_decoding_regression_test.go
  - lib/hsllogic/plc_byte_order_regression_test.go
  - frontend/tests/unit/workbench-v2/mc-persisted-readiness-regression.test.ts
-->

---
### Requirement: FATEK discrete payloads are validated before success

After frame, command status and LRC validation, FATEK Cmd 44 SHALL accept a success payload only when its byte length equals the requested count and every byte is ASCII 0 or 1. Invalid payloads SHALL return an error without a partial value; the adapter SHALL return bad quality and no valid Value. Existing count limits from 1 through 255 SHALL remain enforced.

#### Scenario: Malformed discrete data is rejected

- **WHEN** a count-2 request receives an otherwise valid success frame with an empty payload, 1, 101 or 1X
- **THEN** each response SHALL fail with bad quality
- **AND** neither nil, false nor a shortened slice SHALL be reported as a successful sample.

#### Scenario: Valid scalar and batch responses remain supported

- **WHEN** a valid count-1 payload is 1, a valid count-2 payload is 01, or a valid count-255 payload contains exactly 255 binary characters
- **THEN** the adapter SHALL return true, the false/true slice, or the complete 255-value slice respectively
- **AND** counts 0 and 256 SHALL remain rejected before wire access.

#### Scenario: Malformed payload cannot become a valid acquisition sample

- **WHEN** periodic, manual or read-now acquisition receives a malformed Cmd 44 response
- **THEN** it SHALL propagate the read failure and bad quality without publishing a quality-good value
- **AND** existing failure evidence SHALL remain distinguishable from a successful new sample.

<!-- @trace
source: fix-plc-read-decoding-and-address-consistency
updated: 2026-10-30
code:
  - lib/hsllogic/register_conversion.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/datalink/device/service_point_read.go
  - frontend/src/i18n/locales/en/workbench-v2-plc.json
  - frontend/src/pages/datalink/workbench-v2/useStudioV2PersistedMCPoints.ts
  - lib/hsllogic/bytetransform.go
  - internal/datalink/collector/scheduler_polling.go
  - internal/datalink/sourcerule/service.go
  - internal/datalink/sourcerule/runtime_reconcile.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/ReadinessIssuesWarning.tsx
  - internal/datalink/sourcerule/mc_address_plan.go
  - internal/datalink/collector/scheduler_manual_poll.go
  - internal/datalink/connector/adapters/fatek_helpers.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/pages/datalink/workbench-v2/studioV2LocatorHydration.ts
  - frontend/src/features/datalink/workbench-v2/state/nextFreeStartAddress.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - internal/datalink/point/service_batch_validation.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2-plc.json
  - internal/datalink/connector/point_read.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/state/mcPersistedPoints.ts
  - internal/datalink/connector/adapters/fatek_operations.go
  - internal/datalink/sourcerule/protocol_address_planner.go
  - frontend/src/i18n/config.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - internal/protocol/fatek/client.go
  - frontend/src/utils/addressParser.ts
  - internal/datalink/connector/adapters/modbus_parse_helpers.go
  - internal/protocol/fatek/address_parse.go
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
tests:
  - internal/datalink/sourcerule/mc_validation_boundary_test.go
  - internal/datalink/sourcerule/mc_address_planner_regression_test.go
  - internal/protocol/fatek/status_payload_test.go
  - internal/datalink/collector/point_count_contract_test.go
  - internal/datalink/connector/adapters/fatek_refactor_test.go
  - internal/datalink/connector/adapters/modbus_numeric_bounds_test.go
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/mc-persisted-draft-preservation.test.ts
  - frontend/tests/unit/workbench-v2/mc-persisted-hydration.test.tsx
  - frontend/tests/unit/workbench-v2/studio-locator-mc-hydration.test.tsx
  - frontend/tests/unit/workbench-v2/mc-address-bounds.test.tsx
  - internal/datalink/collector/fatek_status_acquisition_test.go
  - frontend/tests/unit/utils/addressParser.test.ts
  - frontend/tests/unit/workbench-v2/mc-next-free-address.test.ts
  - frontend/tests/unit/workbench-v2/mc-persisted-hydration.fixtures.ts
  - frontend/tests/unit/workbench-v2/occupied-range.test.tsx
  - internal/datalink/sourcerule/protocol_address_test.go
  - internal/datalink/connector/adapters/plc_fatek_decoding_regression_test.go
  - internal/datalink/collector/fatek_point_count_contract_test.go
  - internal/datalink/sourcerule/mc_address_bounds_test.go
  - frontend/tests/unit/workbench-v2/mc-persisted-view.test.tsx
  - internal/datalink/connector/adapters/plc_fatek_status_regression_test.go
  - internal/datalink/connector/adapters/plc_modbus_address_regression_test.go
  - internal/datalink/sourcerule/mc_radix_compatibility_test.go
  - frontend/tests/unit/workbench-v2/mc-persisted-mapping-queue.test.tsx
  - internal/datalink/connector/point_read_test.go
  - internal/datalink/point/mc_hex_address_test.go
  - internal/datalink/sourcerule/mc_radix_runtime_reconcile_test.go
  - internal/datalink/connector/adapters/plc_read_decoding_regression_test.go
  - lib/hsllogic/plc_byte_order_regression_test.go
  - frontend/tests/unit/workbench-v2/mc-persisted-readiness-regression.test.ts
-->