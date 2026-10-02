## MODIFIED Requirements

### Requirement: Workspace database row groups capture shared-column intent
The workspace SHALL persist database row groups as first-class Step 4 planning entities whenever an operator wants multiple points to reuse business columns in one table. For migrated outputs, this representation SHALL be a compatibility projection of the canonical WriteGroup rather than an independent authority. Stable IDs, connector/table scope, member points, group-key and unique-key metadata MUST be preserved. Shared-column reuse is valid only when the row identity contract prevents competing values in the same row.

#### Scenario: Persist one shared-column row group
- **WHEN** an operator creates a row group for one connector and one table, assigns multiple points, and saves Step 4
- **THEN** the workspace persists its stable group identity, scope, membership and key metadata; member targets reference that group, with one canonical owner for migrated output.

#### Scenario: Reload row-group planning on next workspace load
- **WHEN** Studio V2 reloads persisted row groups
- **THEN** it restores the same structure and shared-column memberships instead of flattening them to unrelated globally unique-column bindings.

#### Scenario: Migration cannot preserve row identity
- **WHEN** two shared-column members would target the same destination row and column after conversion
- **THEN** the migration is blocked with an actionable reason and the original group remains unchanged.
