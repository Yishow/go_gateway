## ADDED Requirements

### Requirement: Workspace bootstrap returns a complete persisted setup snapshot

The system SHALL let the Studio V2 bootstrap compose a complete persisted setup snapshot for Step 1 through Step 4 from workspace-owned records.

#### Scenario: Bootstrap reads a complete persisted setup snapshot

- **WHEN** the Studio V2 client loads a workspace with persisted devices, source rules, mappings, database connector settings, and database targets
- **THEN** the backend data model supports composing one consistent persisted snapshot for those records
- **AND** the composed snapshot preserves stable record identity across service restarts
