## ADDED Requirements

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

### Requirement: Connection probes use protocol-appropriate configuration
The system SHALL use protocol-appropriate probe configuration instead of relying on one implicit default probe behavior for all connectors.

#### Scenario: Protocol-specific probe configuration is applied
- **WHEN** a connector test runs for a configured device
- **THEN** the system SHALL use the protocol-specific probe configuration for that device
- **AND** SHALL NOT assume one shared implicit probe address or validation pattern across all protocols
