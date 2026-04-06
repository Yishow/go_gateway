## MODIFIED Requirements

### Requirement: Device readiness check
The system SHALL provide a readiness contract that distinguishes planning eligibility from activation/apply eligibility instead of reducing device readiness to one final boolean.

The readiness contract MUST include at least:

- `connect_status`: `ready` or `failed`
- `probe_status`: `ready`, `pending`, or `failed`
- `planning_allowed`
- `activation_allowed`
- `apply_allowed`

#### Scenario: Partial-success device remains available for planning
- **WHEN** a device has successful connect diagnostics but failed probe diagnostics
- **THEN** the device remains saved and selectable for planning context
- **AND** the readiness contract reports `connect_status=ready`, `probe_status=failed`, `planning_allowed=true`, and both `activation_allowed` and `apply_allowed` as `false`

#### Scenario: Planning readiness no longer depends on downstream assets
- **WHEN** a client checks readiness for a saved device that has connect success but no applied rule-driven tags or outputs yet
- **THEN** the contract can still report `planning_allowed=true`
- **AND** activation/apply readiness remains gated until probe and downstream workflow requirements pass
