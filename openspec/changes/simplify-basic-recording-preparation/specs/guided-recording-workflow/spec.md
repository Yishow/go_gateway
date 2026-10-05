## ADDED Requirements

### Requirement: Complete basic preparation with one effective configuration
Basic setup SHALL complete managed schema preparation without entering Advanced. Connection controls SHALL configure only connection identity; the canonical group SHALL own recording members, storage and interval. Operator primary surfaces SHALL show names, values, types, units and truthful recording stages; technical IDs, API payloads and revisions SHALL remain available in diagnostics.

#### Scenario: Fresh basic recording
- **WHEN** an operator configures eight measurement points with type-appropriate register spans, maps tags and selects an empty managed SQLite destination
- **THEN** Basic provides schema preview and explicit confirmation and proceeds to timed one-row recording without Advanced.

#### Scenario: Inspect effective settings
- **WHEN** an operator prepares a canonical group
- **THEN** no competing connector table, INSERT/UPSERT or recording interval control is presented as effective, and technical identifiers remain inspectable through diagnostics.
