## ADDED Requirements

### Requirement: Local Modbus output uses persistent rule-owned mapping state
The system SHALL persist Local Modbus mapping ownership and revision state for rule-driven output review/apply.

#### Scenario: Restart preserves rule-owned Local Modbus state
- **WHEN** the service restarts after Local Modbus candidates or applied mappings were created from a source rule
- **THEN** the system restores the persisted rule-owned Local Modbus state
- **AND** the workbench can continue review/apply without rebuilding mappings from scratch

### Requirement: Local Modbus candidate conflicts are governed explicitly
The system SHALL detect overlapping Local Modbus register allocations across rule-derived candidates and block only the conflicting subset until the operator resolves the conflict.

#### Scenario: Overlapping register range blocks conflicting candidates
- **WHEN** two or more rule-derived Local Modbus candidates claim overlapping register ranges
- **THEN** the system marks those candidates as `blocked_conflict`
- **AND** prevents apply for only the conflicting candidates until the overlap is resolved
