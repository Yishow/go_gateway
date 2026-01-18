# Project: go-gateway

## Overview

`go-gateway` is an industrial data collector designed for cross-platform deployment (Windows, Linux, Embedded ARM). It acts as a bridge between industrial devices and IT databases/cloud systems.

## Key Features

- **Multi-Protocol Support**: Modbus RTU/TCP, FATEK FBs, Mitsubishi MC Protocol (3E), MQTT.
- **Flexible Data Mapping**: Configurable "Datalink" engine to map PLC registers directly to Database tables/columns.
- **Cross-Platform**: Built with Go for single-binary deployment on Windows, Linux, and Raspberry Pi.
- **Multi-Database**: Supports SQLite, MySQL, PostgreSQL, and SQL Server.

## Current Status

**Phase:** Design & Prototyping
**Primary Language:** Go (Planned), Python (Prototyping)

The project is currently in the detailed design phase. A comprehensive architecture specification and a Python prototype for the FATEK driver have been developed.

## Key Documentation

The core design documents are located in the `docs/` directory. **Read these first to understand the system architecture.**

- **[`docs/snazzy-herding-firefly.md`](docs/snazzy-herding-firefly.md)**: **Main Design Document**. Contains the system overview, module structure, interface definitions (Go), and database schema.
- **[`docs/plan.md`](docs/plan.md)**: Work Breakdown Structure (WBS) and detailed implementation phases.
- **[`docs/fatek.md`](docs/fatek.md)**: Technical specification for the FATEK Protocol (ASCII/TCP), including frame structure and checksum algorithms.

## Prototypes

- **[`fatek_driver.py`](fatek_driver.py)**: A fully functional Python implementation of the FATEK PLC driver. It supports both Serial and TCP modes, auto-checksum calculation, and high-level commands (Read/Write Status, Registers, Mixed Read). Use this as a reference when implementing the Go version.

## Architecture Highlights

- **Protocol Abstraction**: A common `Protocol` interface allows easy addition of new device drivers.
- **Datalink Engine**: An ETL-like engine that handles `Source (PLC) -> Transform -> Target (DB)` logic.
- **Task Engine**: A scheduler that manages concurrent data collection tasks with different intervals and priorities.

## Usage (Prototype)

### FATEK Driver (Python)

To use the Python driver prototype:

```bash
# Install dependencies
pip install pyserial

# Import in your python script
from fatek_driver import create_tcp_client
client = create_tcp_client('192.168.1.5')
client.connect()
print(client.read_registers('D', 0, 10))
```

## Planned Go Workflow

Once the Go implementation starts, the standard workflow will be:

```bash
# Initialize dependencies
go mod tidy

# Run tests
go test ./...

# Build binary
go build -o gateway ./cmd/gateway

# Run with configuration
./gateway -c configs/gateway.yaml
```

繁體中文回答所有的問題。
依規畫的階段逐步執行，不得跳過任何一個階段及步驟，不能以最小實作，而是要完整。
