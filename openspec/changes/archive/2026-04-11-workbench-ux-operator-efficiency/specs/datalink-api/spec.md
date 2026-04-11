## ADDED Requirements

### Requirement: Studio V2 output apply APIs are rule-scoped and revision-aware
The API SHALL expose explicit SourceRule-scoped output apply endpoints for Studio V2 without relying on frontend mocks.

Endpoints:

- `POST /api/v1/datalink/source-rules/:id/database-outputs/apply`
- `POST /api/v1/datalink/source-rules/:id/local-modbus/apply`

The apply request MUST include `revision_id` and `candidate_ids[]`.

The apply response MUST return per-item results (`success` / `failed` / `skipped`) with actionable reasons.

#### Scenario: Revision mismatch blocks apply with explicit reason
- **WHEN** a client applies output candidates with a stale `revision_id`
- **THEN** the API rejects the request with a revision-conflict error
- **AND** returns enough context for UI to trigger candidate refresh instead of silent retry

#### Scenario: Partial apply reports per-item results
- **WHEN** a client applies multiple output candidates and only part of them succeed
- **THEN** the API returns per-candidate result details
- **AND** does not collapse the outcome into one ambiguous pass/fail flag

### Requirement: Studio V2 database tooling APIs support schema preview, dry-run validation, and write history
The API SHALL provide dedicated database tooling endpoints so Step 4 can run schema planning and verification using real backend data.

Endpoints:

- `POST /api/v1/datalink/db-targets/connectors/:id/schema/generate`
- `POST /api/v1/datalink/db-targets/connectors/:id/mappings/dry-run`
- `GET /api/v1/datalink/db-targets/connectors/:id/write-history`

The schema generate endpoint MUST support `dry_run=true|false` behavior in the request.

The dry-run endpoint MUST return candidate-level validation outcomes and blocking reasons.

The write-history endpoint MUST return recent write summaries including timestamp, result status, affected row count, and error summary when failed.

#### Scenario: Schema dry-run does not mutate connector schema
- **WHEN** a client calls schema generate with `dry_run=true`
- **THEN** the API returns SQL preview and validation feedback
- **AND** does not execute schema mutations

#### Scenario: Mapping dry-run identifies blocked candidates before apply
- **WHEN** a client requests mapping dry-run for selected candidates
- **THEN** the API returns which candidates are applyable or blocked
- **AND** includes machine-readable blocking categories (for example `schema_missing`, `connector_unavailable`, `type_conflict`)

#### Scenario: Write history supports operator diagnosis
- **WHEN** a client requests the latest write history for a connector
- **THEN** the API returns recent success/failure records in reverse chronological order
- **AND** the payload can be rendered directly in Studio Step 4 history panels
