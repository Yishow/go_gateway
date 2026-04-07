## MODIFIED Requirements

### Requirement: Device readiness API
The system SHALL expose an API endpoint to check device readiness status using the explicit planning, activation, and apply eligibility contract.

Endpoint:

- `GET /api/v1/datalink/devices/:id/readiness` - Check device readiness

#### Scenario: Get readiness status
- **WHEN** a client requests readiness status for a device
- **THEN** the API returns `connect_status`, `probe_status`, `planning_allowed`, `activation_allowed`, `apply_allowed`, and blocking metadata

#### Scenario: Readiness response format
- **WHEN** a readiness check is performed
- **THEN** the API returns a response with explicit planning/activation/apply eligibility fields
- **AND** does not require clients to infer workflow state from one final `ready` boolean alone

## ADDED Requirements

### Requirement: Rule-scoped candidate query and recompute APIs
The API SHALL expose rule-scoped endpoints to query and recompute candidate snapshots for a source-rule revision.

Endpoints:

- `GET /api/v1/datalink/source-rules/:id/candidates` - Get the current candidate snapshot for a source rule
- `POST /api/v1/datalink/source-rules/:id/candidates/recompute` - Recompute the candidate snapshot for the latest source-rule revision

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

#### Scenario: Recompute candidate snapshot
- **WHEN** a client requests candidate recompute for a source rule
- **THEN** the API regenerates the current revision's candidate snapshot using the shared identity/signature contract
- **AND** returns the recomputed snapshot without introducing target-specific apply side effects
