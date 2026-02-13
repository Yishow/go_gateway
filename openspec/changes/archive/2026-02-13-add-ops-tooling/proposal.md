# Proposal: Add Operational Tools (add-ops-tooling)

## Why

As `go-gateway` moves towards production, operations teams and developers face two key challenges:

1. **API Visibility**: The REST API and WebSocket endpoints are documented only in static markdown files, which can drift from implementation. Developers need interactive, auto-generated documentation to test and integrate with the gateway.
2. **Configuration Safety**: The YAML configuration files (`gateway.yaml`, `datalink.yaml`) are complex. Currently, errors are only caught at runtime, potentially causing start-up failures in production environments.

## What Changes

We propose adding two operational tools to the `gateway` binary:

1. **Auto-generated API Docs**: Integrate `swaggo/swag` to generate OpenAPI v2 specification from Go comments and serve Swagger UI at `/swagger/index.html`.
2. **Config Validation Command**: Add a `validate` subcommand (e.g., `gateway validate -c configs/`) that parses and checks configuration logic without starting the server.

## Operational Impact

- **Ops Teams**: Can verify configuration correctness before deployment (CI/CD friendly).
- **Developers**: Can use valid Swagger UI to explore APIs without reading source code.

## Risks

- **Binary Size**: Embedding Swagger UI assets might increase binary size slightly.
- **Maintenance**: Developers must write Go comments in a specific format for `swag` to work.
