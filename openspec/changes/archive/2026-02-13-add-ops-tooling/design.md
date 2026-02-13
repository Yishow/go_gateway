# Design: Operational Tools

## Architecture

### API Documentation (Swag)

- **Library**: `github.com/swaggo/swag` for generation, `github.com/swaggo/gin-swagger` for serving.
- **Workflow**:
  1. Developers add declarative comments to API handlers.
  2. `swag init` is run (via Makefile) to generate `docs/docs.go`, `docs/swagger.json`, `docs/swagger.yaml`.
  3. The `gateway` server imports `docs` package and mounts a route `/swagger/*any`.
- **Constraint**: The `docs` package _must_ be regenerated whenever API code changes.

### Configuration Validation (CLI)

- **Library**: Reuse existing `spf13/viper` and `spf13/cobra` setup.
- **Command Structure**:
  - `gateway root` (starts server by default)
  - `gateway validate` (new subcommand)
- **Validation Logic**:
  - The validation logic currently embedded in `app.Initialize` needs to be extracted into a reusable `config.LoadAndValidate()` function.
  - The `validate` command will call this function and print specific errors (e.g., "Device 'plc-01' references unknown protocol 'modbus_xyz'").
  - **Exit Codes**: 0 for success, 1 for failure (script-friendly).

## Trade-offs

- **Swag vs OpenAPIV3**: `swaggo/swag` primarily supports OpenAPI 2.0 (Swagger 2.0). While OpenAPI 3.0 is newer, `swag` ecosystem is more mature for Go/Gin. We stick to Sway for simplicity.
