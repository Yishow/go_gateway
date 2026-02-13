# Spec: API Documentation

## ADDED Requirements

### Requirement: Auto-generated API documentation

The gateway SHALL serve auto-generated API documentation (Swagger UI) and SHALL generate OpenAPI spec from Go source annotations.

#### Scenario: Serve Swagger UI
- **WHEN** the gateway is running and the operator accesses `/swagger/index.html`
- **THEN** Swagger UI loads and lists all available API endpoints defined in `internal/api`
- **AND** the operator can execute "Try it out" requests from the browser

#### Scenario: Generate specs from comments
- **WHEN** `swag init` is executed in the project root
- **THEN** `docs/swagger.json` and `docs/swagger.yaml` are updated from Go definitions
- **AND** the health check endpoint (`/health`) appears in the generated spec
