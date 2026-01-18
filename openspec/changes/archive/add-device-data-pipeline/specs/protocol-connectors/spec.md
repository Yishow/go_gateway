## ADDED Requirements
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
