## ADDED Requirements

### Requirement: Downstream integrity gaps surface as readiness issues before activation

The system SHALL surface missing or inconsistent downstream rule-derived relationships as readiness issues before activation.

#### Scenario: Missing downstream relationships become readiness blockers

- **WHEN** a persisted source rule is missing a required derived point, tag, mapping, or required database target relationship for the chosen activation scope
- **THEN** the workspace readiness result records a normalized issue for that gap before activation starts
- **AND** the operator SHALL NOT first discover that gap from a late runtime or write-path failure

##### Example: missing tag and missing database target both surface early

- **GIVEN** rule rule-A still exists but point pt-A has no persisted tag and pt-B has no required database target
- **WHEN** readiness evaluates the workspace before activation
- **THEN** readiness returns normalized issues for rule-A or its affected points instead of waiting for runtime to fail
