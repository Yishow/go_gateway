# protocol-connectors Specification

## Purpose
TBD - created by archiving change add-device-data-pipeline. Update Purpose after archive.
## Requirements
### Requirement: Protocol connector registry
The system SHALL provide a common protocol interface and a registry to load protocol connectors by type.

#### Scenario: Load protocol
- WHEN a device is configured with a protocol type
- THEN the system instantiates the matching connector from the registry

### Requirement: Modbus connectors
The system SHALL support Modbus TCP, Modbus RTU, and Modbus UDP for reading and writing registers.

#### Scenario: Read Modbus holding register
- WHEN a Modbus TCP device is polled for a holding register
- THEN the connector returns the raw value with a timestamp

### Requirement: FATEK FBs connector
The system SHALL support FATEK FBs protocol with ASCII frames and LRC validation.

#### Scenario: Read FATEK device
- WHEN a FATEK device is polled
- THEN the connector validates LRC and returns a decoded value

### Requirement: Mitsubishi MC 3E connector
The system SHALL support Mitsubishi MC Protocol 3E frame communication over TCP.

#### Scenario: Read MC 3E point
- WHEN an MC 3E device is polled
- THEN the connector returns the decoded value for the address

### Requirement: MQTT ingest connector
The system SHALL support MQTT subscription as a data source, mapping topics and payloads into point values.

#### Scenario: Ingest MQTT topic
- WHEN a subscribed topic receives a payload
- THEN the system emits a point value for mapping and storage

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

### Requirement: Connector diagnostics expose planning capability hints
Protocol connector diagnostics SHALL expose the capability hints needed by rule-driven planning and later review flows.

#### Scenario: Planning capability hints are returned with diagnostics
- **WHEN** a client requests connector diagnostics for a selected device
- **THEN** the system returns connect/probe phase results plus protocol capability hints relevant to source-rule planning
- **AND** those hints remain distinguishable from transient test failure messages

### Requirement: Connection probes use protocol-appropriate configuration
The system SHALL use protocol-appropriate probe configuration instead of relying on one implicit default probe behavior for all connectors.

#### Scenario: Protocol-specific probe configuration is applied
- **WHEN** a connector test runs for a configured device
- **THEN** the system SHALL use the protocol-specific probe configuration for that device
- **AND** SHALL NOT assume one shared implicit probe address or validation pattern across all protocols

### Requirement: Multi-register decode byte order is configurable on Modbus-class connections
The system SHALL expose **`data_format`** in the **JSON ConfigSchema** (and persist the value in device `connection_config`) for **Modbus TCP**, **Modbus RTU**, and **Modbus UDP** protocols, using the same enumeration and semantics as existing multi-register decoding (`ABCD`, `BADC`, `CDAB`, `DCBA`) and the shared decoding library (for example `lib/hsllogic`). The Modbus connector adapters SHALL read this field when decoding multi-register integers and IEEE-754 floats.

#### Scenario: Modbus TCP device saves data_format
- **WHEN** an operator configures a Modbus TCP device with a non-default `data_format`
- **THEN** the stored configuration SHALL include the selected `data_format`
- **AND** subsequent reads of multi-register data types SHALL use that ordering

#### Scenario: Schema parity across Modbus variants
- **WHEN** the UI lists Modbus TCP, RTU, or UDP
- **THEN** each protocol’s connection schema SHALL include `data_format` with the same enum values and documented defaults

### Requirement: MC 3E retains data_format with documented default
The system SHALL continue to support **`data_format`** on **Mitsubishi MC 3E** connections. The default SHALL remain consistent with existing manufacturer-oriented behavior unless a breaking migration is explicitly approved.

#### Scenario: MC 3E connection still accepts data_format
- **WHEN** an operator configures MC 3E with `data_format` set
- **THEN** the connector SHALL apply that ordering to float and multi-word decode paths

### Requirement: Protocols without multi-register Modbus-style blocks do not require data_format
Protocols that do not use the same holding-register multi-word decode path (for example **MQTT** ingest as configured today) SHALL NOT be required to present `data_format` in their connection schema unless a future capability adds equivalent multi-word semantics.

#### Scenario: MQTT connection omits data_format
- **WHEN** an operator configures an MQTT device per current schema
- **THEN** the absence of `data_format` SHALL not block save
