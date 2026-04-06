## MODIFIED Requirements

### Requirement: Flow state machine for mapping lifecycle
The system SHALL track rule-derived mapping pipelines with explicit lifecycle states: `draft`, `validated`, `active`, `out_of_sync`, and `error`.

`out_of_sync` means the same rule-owned mapping identity still exists but its current proposed signature no longer matches the last applied signature.

#### Scenario: Draft to validated transition
- **WHEN** an operator completes mapping configuration and runs validation
- **THEN** the mapping state transitions from `draft` to `validated`

#### Scenario: Validated to active transition
- **WHEN** an operator activates a validated mapping
- **THEN** the mapping state transitions to `active`

#### Scenario: Active mapping becomes out of sync after rule revision
- **WHEN** a rule revision changes the derived transform pipeline for a mapping that already has applied state
- **THEN** the system marks the mapping as `out_of_sync`
- **AND** requires explicit operator reapply instead of silently replacing the applied pipeline

#### Scenario: Runtime failure transition
- **WHEN** source read, transform, or write execution fails
- **THEN** the mapping state transitions to `error`
- **AND** error metadata includes failed segment and reason

## ADDED Requirements

### Requirement: Rule-derived mappings use shared identity and signature contracts
The mapping pipeline SHALL compare rule-derived mappings using the shared rule-owned identity and proposed-signature contracts defined by source-rule orchestration.

#### Scenario: Same identity with new signature preserves applied mapping
- **WHEN** the system recomputes a rule-derived mapping with the same rule-owned identity and a different proposed signature
- **THEN** the existing applied mapping remains intact
- **AND** the recomputed mapping is surfaced as a reviewable `out_of_sync` candidate
