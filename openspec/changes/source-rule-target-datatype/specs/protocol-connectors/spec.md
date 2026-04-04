# protocol-connectors Specification Delta

## ADDED Requirements

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
