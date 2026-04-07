## ADDED Requirements

### Requirement: Rule-scoped candidate APIs
The API SHALL expose rule-scoped endpoints to query and recompute tag and output candidates for a source-rule revision.

The candidate query response MUST include at least:

- `source_rule_id`
- `revision_id`
- `tags`
- `database_outputs`
- `local_modbus_outputs`

Each returned candidate object MUST include at least:

- `id`
- `status`
- `proposed_signature`
- `last_applied_signature`
- `blocking_reason`

#### Scenario: Query rule candidates
- **WHEN** a client requests candidates for a source rule
- **THEN** the API returns the current tag, database, and Local Modbus candidate sets with revision and status metadata

### Requirement: Rule-scoped apply APIs
The API SHALL expose explicit per-target apply endpoints for rule-driven tags, database outputs, and Local Modbus outputs.

An apply request MUST identify the target type and the candidate ids being applied.

An apply response MUST report per-item success or failure so partial success is explicit.

#### Scenario: Apply one target without affecting the other
- **WHEN** a client applies database output candidates for a source rule
- **THEN** the API applies only the database target changes
- **AND** preserves Local Modbus candidate and apply state for the same rule

#### Scenario: Partial apply returns per-item results
- **WHEN** a client applies multiple candidates and only some succeed
- **THEN** the API returns per-candidate success and failure results
- **AND** does not collapse the response into one ambiguous pass/fail result
