# Spec: API Documentation

## ADDED Requirements

#### Scenario: Serve Swagger UI

- When the gateway is running, accessing `http://localhost:8080/swagger/index.html` should load the Swagger UI.
- The UI should list all available API endpoints defined in `internal/api`.
- Users should be able to execute "Try it out" requests directly from the browser.

#### Scenario: Generate Specs from Comments

- When `swag init` is executed in the project root, `docs/swagger.json` and `docs/swagger.yaml` should be updated based on Go definitions.
- Health check endpoint (`/health`) should appear in the generated spec.

## Dependencies

- `github.com/swaggo/swag/cmd/swag` (Tool)
- `github.com/swaggo/gin-swagger` (Library)
- `github.com/swaggo/files` (Library)
