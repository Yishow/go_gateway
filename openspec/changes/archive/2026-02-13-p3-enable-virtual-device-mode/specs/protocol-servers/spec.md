# protocol-servers Spec Delta

## ADDED Requirements

### Requirement: Virtual server hosting

The system SHALL support hosting one or more protocol servers (listeners) to accept incoming connections from external clients (SCADA/HMI).

#### Scenario: Host Modbus TCP server

- GIVEN the gateway is configured to start a Modbus TCP server on port 10502
- WHEN the gateway starts
- THEN it listens on port 10502 and accepts Modbus queries

### Requirement: Virtual memory mapping

The system SHALL provide a strictly addressed virtual memory space that serves as the data source for hosted servers.

#### Scenario: Read from virtual memory

- GIVEN a value 12345 is stored at virtual address D100
- WHEN an external Modbus client reads Holding Register 100
- THEN the server responds with value 12345

### Requirement: Data bridging

The system SHALL support mapping collected data from real devices into the virtual memory space.

#### Scenario: Bridge data to virtual server

- GIVEN a real device tag `PLC1.D200` is collected
- AND a mapping exists from `PLC1.D200` to Virtual Address `V.D100`
- WHEN `PLC1.D200` changes value
- THEN `V.D100` is automatically updated
