## ADDED Requirements

### Requirement: Device scoped live point identity
The runtime table SHALL show only the selected device's mappings and values. Values and metadata SHALL match device and point identity; address alone SHALL NOT establish identity. Missing identity or a missing matching event SHALL display waiting rather than another device's live data.

#### Scenario: Same address on two devices
- **WHEN** A.D0=215 and B.D0=187 are mapped in one workspace and the selection changes from A to B
- **THEN** B displays only B's name, type, unit and matching value 187; A's stale event 215 is never attributed to B.
