# Project Context

## Purpose

`go-gateway` is a cross-platform (Windows/Linux/ARM) industrial data collector designed to bridge OT (Operation Technology) and IT systems. It collects data from PLCs via various protocols (Modbus, FATEK, MC Protocol) and syncs it to databases or MQTT brokers using a configurable ETL engine (Datalink).

## Tech Stack

- **Language**: Go 1.22+
- **Core Libraries**:
  - `github.com/gin-gonic/gin` (HTTP/API)
  - `github.com/spf13/viper` (Configuration)
  - `github.com/spf13/cobra` (CLI)
  - `go.uber.org/zap` (Logging)
  - `modernc.org/sqlite`, `gorm.io/gorm` (Database)
  - `github.com/eclipse/paho.mqtt.golang` (MQTT)

## Project Conventions

### Code Style

- Standard Go formatting (`gofmt`).
- Error handling: explicit error returns, wrapping with context.
- Configuration: YAML-based, loaded via Viper.

### Architecture Patterns

- **Hexagonal Architecture**: Core logic (Datalink, Task) decoupled from adapters (Protocol, Database).
- **Single Binary**: All components (Backend, Web UI, CLI) packaged into one executable.

### Testing Strategy

- Unit tests for core logic (Datalink, Protocol parsers).
- Integration tests for Database and Protocol adapters.
- **TDD Mandate**: For high-complexity logic (e.g., Algorithms, State Machines, Concurrency Control), the **Red-Green-Refactor** TDD workflow is MANDATORY in `tasks.md`.

## Domain Context

- **Tags**: Data points on a device (e.g., "D100").
- **Datalink**: The mapping rule from Source (Device Tag) to Target (Database Column).
- **Driver**: Implementation of a specific communication protocol.

## Important Constraints

- **Performance**: High concurrency for polling multiple devices.
- **Reliability**: Auto-reconnection and buffering are critical.
- **Environment**: Must run on low-resource hardware (Raspberry Pi) and Windows Servers.
