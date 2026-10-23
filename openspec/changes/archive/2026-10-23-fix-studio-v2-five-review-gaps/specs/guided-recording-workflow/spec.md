## ADDED Requirements

### Requirement: Basic missing-value policy reveals required gating
Basic recording SHALL state that required missing or invalid members still block a partial row, show the effective required state, and provide a narrow entry to existing per-member required controls without silently changing any member.

#### Scenario: Required gap remains blocked
- **WHEN** a basic group selects partial and a required member is missing or invalid
- **THEN** the UI SHALL explain that the row is skipped and SHALL expose the existing member settings without changing required automatically.

#### Scenario: Optional gap and recovery
- **WHEN** an explicitly optional member is missing in a partial-capable destination and subsequently recovers
- **THEN** the system SHALL preserve NULL and its reason for the partial row and write its actual value after recovery; skip_row and required protection SHALL remain unchanged.

#### Scenario: Edit existing exact managed members
- **WHEN** an operator opens required settings for saved SQLite managed uint16 or uint32 members in canonical generated TEXT columns
- **THEN** the editor SHALL retain their exact storage compatibility and require fresh schema confirmation after a group revision change without weakening custom-table checks.

#### Scenario: Confirm unchanged schema after member policy edit
- **WHEN** a saved managed group changes required or completeness policy without changing its physical columns
- **THEN** a compatible preview SHALL return an empty statements array, retain its revision-bound token and operation identity, and require explicit confirmation before applying the new group revision.
