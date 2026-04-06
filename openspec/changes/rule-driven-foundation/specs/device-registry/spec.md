## MODIFIED Requirements

### Requirement: Device readiness check
The system SHALL provide a readiness contract that distinguishes planning eligibility from activation/apply eligibility instead of reducing device readiness to one final boolean.

The readiness contract MUST include at least:

- `connect_status`
- `probe_status`
- `planning_allowed`
- `activation_allowed`
- `apply_allowed`
- `blocking_reasons`

#### Scenario: Planning can proceed after connect success
- **WHEN** a device has successful connect diagnostics but failed probe diagnostics
- **THEN** the device remains saved and selectable for planning context
- **AND** the readiness contract reports `planning_allowed=true` while both `activation_allowed` and `apply_allowed` remain `false`

#### Scenario: Activation requires successful probe
- **WHEN** a device has successful connect and probe diagnostics
- **THEN** the readiness contract allows activation and apply according to current downstream prerequisites
- **AND** the response distinguishes probe success from any later workflow-specific blocking state

#### Scenario: Readiness reports explicit blocking reasons
- **WHEN** a client requests readiness status for a device that is not fully eligible for activation or apply
- **THEN** the system returns the relevant blocking reasons in the readiness contract
- **AND** does not require the caller to infer readiness state from one summary boolean alone
