## MODIFIED Requirements

### Requirement: Auto-generated API documentation

The gateway SHALL serve auto-generated API documentation through Swagger UI and SHALL generate OpenAPI specs from Go source annotations. The documented contract SHALL include the typed safe-error envelope for preview, runtime snapshot/stream, workspace readiness, and activation failures, plus the immediate removal of any dedicated `/studio` route/handler/tombstone so the path follows the existing generic unknown-route policy. The contract SHALL document code, message, retryable, and request_id without exposing raw diagnostics as an operator-facing schema.

#### Scenario: Serve Swagger UI

- **WHEN** the gateway is running and the operator accesses /swagger/index.html
- **THEN** Swagger UI loads and lists all available API endpoints defined in internal/api
- **AND** the preview/runtime failure responses show the typed error schema and safe response fields

#### Scenario: Generate specs from comments

- **WHEN** swag init is executed in the project root
- **THEN** docs/swagger/docs.go, docs/swagger/swagger.yaml, and docs/swagger/swagger.json are updated from Go definitions
- **AND** the health check endpoint /health and the affected preview/runtime endpoints appear in the generated spec

#### Scenario: Document deleted route behavior

- **WHEN** an API or route documentation check evaluates the studio surfaces after `/studio` deletion
- **THEN** it documents generic unknown-route-equivalent behavior after removal of the dedicated `/studio` route/handler/tombstone
- **AND** it identifies /test and /gateway/* as non-target routes that remain unchanged

## ADDED Requirements

### Requirement: API registry and release documentation match generated contract

The implementation SHALL update the backend API registry, runtime/inventory documentation, and the release note for this change with the same endpoint ownership, typed error codes, immediate deletion, generic unknown-route equivalence, preserved route identities, pre-delete inventory, and Git rollback limitations as the generated OpenAPI contract. Inventory document edits SHALL have a matching studio_inventory_changelog entry.

#### Scenario: Registry is consistent with Swagger

- **WHEN** a reviewer compares the generated Swagger output with backend-api-registry.md
- **THEN** affected endpoint paths, response codes, typed error codes, and route ownership match
- **AND** unrelated historical Swagger gaps are explicitly outside this change

#### Scenario: Release evidence is traceable

- **WHEN** the release documentation gate runs
- **THEN** docs/releases/retire-legacy-studio-and-polish-v2.md records the deletion inventory, generic unknown-route comparison, dedicated Git commit/revert evidence, rebuild/redeploy smoke, preserved route identities, and source/data/bookmark rollback limitations
- **AND** every modified studio-surface-inventory document is named in a changelog.sqlite record containing summary, surface, files, and reason
