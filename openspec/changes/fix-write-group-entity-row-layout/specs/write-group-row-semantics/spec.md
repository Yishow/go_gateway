## ADDED Requirements

### Requirement: Entity-scoped column layout agreement
Column validation and row encoding SHALL use the same persisted entity-to-member partition. Distinct identified rows within a group MAY share a destination column, but competing values within the same row MUST be rejected before Apply. UI guidance, readiness and runtime MUST agree for the same saved layout.

#### Scenario: Shared column in distinct rows
- **GIVEN** one group has entities A and B, each with a member mapped to value and a verified distinct row identity
- **WHEN** the group is saved, checked and applied
- **THEN** all layers accept the layout and retain both entities rather than reporting a global duplicate-column conflict.

#### Scenario: Competing values in one row
- **WHEN** two members in entity A map to the same effective SQL column, including equivalent case variants for that destination
- **THEN** readiness and Apply reject the collision with an actionable explanation before intake begins.

### Requirement: Entity-local encoding and recoverable structural failure
Each finalized row SHALL encode only its own entity members with the existing exact-value and quality policies. Missing members in other entities MUST NOT invalidate that row. A structural layout/encoding mismatch MUST preserve unresolved accepted input and checkpoint evidence rather than consuming it as a normal quality skip.

#### Scenario: Independent values in one group
- **WHEN** entities A and B finalize using shared or different destination column names
- **THEN** SQL contains each entity's own values and identity without requiring the other entity's members.

#### Scenario: One entity lacks data
- **WHEN** A has a missing required member while B has all required good members
- **THEN** A records its existing scoped incomplete outcome while B produces its own valid row; explicit partial and silent-bucket rules remain unchanged.

#### Scenario: Structural mismatch after durable acceptance
- **WHEN** an acknowledged sample cannot be encoded because the saved layout is structurally inconsistent
- **THEN** its unresolved closure is retained with a safe diagnostic and is not consumed as a successful row or ordinary missing-data skip.
