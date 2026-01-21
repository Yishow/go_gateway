# protocol-connectors Spec Delta

## MODIFIED Requirements

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

## ADDED Requirements

### Requirement: Connection health check

The system SHALL provide a heartbeat/keep-alive mechanism for persistent connections to detect stale sockets.

#### Scenario: Keep-alive timeout

- GIVEN a persistent connection has been idle for longer than the configured `keep_alive_interval`
- WHEN no data is exchanged
- THEN the connector sends a health-check packet or closes the connection preemptively

---

## Cross-References

- Related capability: `collection-scheduler` (for polling intervals)
- Related capability: `device-registry` (for device configuration)
