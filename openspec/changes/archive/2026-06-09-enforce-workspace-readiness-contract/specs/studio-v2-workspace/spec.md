## ADDED Requirements

### Requirement: Workspace exposes readiness summary for Studio V2 surfaces

The system SHALL expose a workspace readiness summary that Studio V2 shell surfaces can consume consistently.

#### Scenario: Shell reads workspace readiness summary

- **WHEN** the Studio V2 shell loads a persisted workspace
- **THEN** it can retrieve a readiness summary with blocking counts, warning counts, and normalized issue entries
- **AND** the same readiness summary can be used by shell, summary rail, and activation entry points without re-deriving ad hoc logic in each surface
