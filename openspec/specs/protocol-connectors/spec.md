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
The system SHALL interpret single-point requests as one typed value for MC 3E, while preserving register-count semantics for Modbus and FATEK and the existing typed-value count of direct MC batch requests.

#### Scenario: Multiword MC point across entrypoints
- **WHEN** periodic, manual or read-now acquisition reads one uint32 point at D3 with words 34464 and 1 in CDAB order
- **THEN** it SHALL request two words and return scalar 100000 rather than an array; a single 64-bit value SHALL occupy four words and retain exact encoding.

#### Scenario: Direct batch and other protocols remain unchanged
- **WHEN** a direct MC batch requests two uint32 values or a Modbus/FATEK point is polled
- **THEN** MC SHALL retain two-value batch decoding and other protocols SHALL retain their existing word counts.

<!-- @trace
source: fix-studio-v2-five-review-gaps
updated: 2026-10-23
code:
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel.tsx
  - docs/releases/2026-10-23-studio-five-review-fixes.md
  - frontend/src/i18n/locales/en/studio-review.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/datalink/collector/scheduler_polling.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceMeasurements.ts
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection.tsx
  - docs/releases/evidence/studio-five-review-2026-10-23/manifest.json
  - internal/datalink/collector/scheduler_manual_poll.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor.tsx
  - internal/datalink/device/service_point_read.go
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - internal/datalink/connector/point_read.go
  - internal/datalink/recordingplan/schema_preview.go
  - docs/releases/evidence/studio-five-review-2026-10-23/observations.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupMemberTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MeasurementTemplateSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/state/writeGroup/columns.ts
  - frontend/src/i18n/locales/zh-TW/studio-review.json
tests:
  - frontend/tests/unit/workbench-v2/template-offline.test.tsx
  - frontend/tests/unit/workbench-v2/occupied-range.test.tsx
  - internal/datalink/collector/point_count_contract_test.go
  - frontend/tests/unit/workbench-v2/MeasurementTemplateSelector.test.tsx
  - frontend/tests/unit/workbench-v2/basic-required-policy.test.tsx
  - internal/datalink/recordingplan/schema_noop_json_test.go
  - scripts/tests/f_device_to_sql/five-template-errors.mjs
  - scripts/tests/f_device_to_sql/five-review-flow.mjs
  - scripts/tests/f_device_to_sql/mc-fixture.mjs
  - frontend/tests/unit/workbench-v2/writeGroupColumns.test.ts
  - frontend/tests/unit/runtime-dashboard/focused-header-evidence.test.tsx
-->