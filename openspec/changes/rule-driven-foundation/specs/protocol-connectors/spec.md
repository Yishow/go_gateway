## ADDED Requirements

### Requirement: Connector diagnostics expose planning capability hints
Protocol connector diagnostics SHALL expose the capability hints needed by rule-driven planning and later review flows.

#### Scenario: Planning capability hints are returned with diagnostics
- **WHEN** a client requests connector diagnostics for a selected device
- **THEN** the system returns connect/probe phase results plus protocol capability hints relevant to source-rule planning
- **AND** those hints remain distinguishable from transient test failure messages
