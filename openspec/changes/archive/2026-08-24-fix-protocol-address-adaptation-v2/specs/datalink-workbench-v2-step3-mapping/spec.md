# datalink-workbench-v2-step3-mapping Specification Delta

## ADDED Requirements

### Requirement: Protocol-aware default tag key generation

The system SHALL generate default tag keys and display names in Step 3 based on the point's protocol-adapted address. The generated tag key MUST retain the register prefix identifier for non-Modbus devices (e.g. `dev.sensor.rd0` for MC/FATEK `D0` instead of stripped numeric index).

#### Scenario: Default tag key for MC 3E point
- **WHEN** entering Step 3 with an MC 3E point at address `D100`
- **THEN** the suggested tag key contains `rd100` and accurately reflects the alphanumeric address

#### Scenario: Default tag key for Modbus point
- **WHEN** entering Step 3 with a Modbus point at address `40001`
- **THEN** the suggested tag key contains `r40001`