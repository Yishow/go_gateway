## 1. Planning and data model
- [ ] 1.1 Confirm tag naming rules, time precision scope, and transform expression syntax
- [ ] 1.2 Define database schema for devices, points, tags, mappings, tasks, and time-series storage
- [ ] 1.3 Define Postgres partition strategy and create migration plan
- [ ] 1.4 Define SQLite test schema that mirrors the Postgres logical model
- [ ] 1.5 Define CMS integration contract (schema, indexes, and views if needed)

## 2. Protocol framework
- [ ] 2.1 Implement protocol interface, connector registry, and connection manager
- [ ] 2.2 Implement Modbus TCP/RTU/UDP connectors with read/write support
- [ ] 2.3 Implement FATEK FBs connector with LRC validation and parsing
- [ ] 2.4 Implement Mitsubishi MC 3E connector with address mapping
- [ ] 2.5 Implement MQTT ingest connector with topic and payload mapping

## 3. Device and point management
- [ ] 3.1 Implement device CRUD and connection test logic
- [ ] 3.2 Implement point CRUD with protocol-specific fields
- [ ] 3.3 Implement point status tracking (last read, error status)

## 4. Collection scheduler
- [ ] 4.1 Implement polling group scheduler with interval control
- [ ] 4.2 Implement concurrency control per device and retry/backoff
- [ ] 4.3 Attach timestamp and quality flags to collected values

## 5. Tag dictionary and mapping
- [ ] 5.1 Implement global tag dictionary with metadata and lifecycle states
- [ ] 5.2 Implement mapping storage linking points to tags
- [ ] 5.3 Implement transform pipeline steps (decode, cast, scale, lookup, conditional, formula)
- [ ] 5.4 Implement mapping validation and preview pipeline

## 6. Time-series storage
- [ ] 6.1 Implement time-series table schema and indexes
- [ ] 6.2 Implement Postgres partition creation and rollover logic
- [ ] 6.3 Implement batch writer with flush size and interval
- [ ] 6.4 Implement external read contract for CMS (indexes and optional views)

## 7. API layer
- [ ] 7.1 Implement REST APIs for devices, points, tags, and mappings
- [ ] 7.2 Implement connection test and mapping preview endpoints
- [ ] 7.3 Implement system settings endpoints for write precision and partition interval

## 8. UI layer
- [ ] 8.1 Build device and point management pages
- [ ] 8.2 Build tag dictionary management page
- [ ] 8.3 Build drag-drop mapping canvas and transform builder
- [ ] 8.4 Build mapping preview panel with step-by-step values
- [ ] 8.5 Build write precision settings panel

## 9. Validation and testing
- [ ] 9.1 Add SQL-based integration tests with SQLite for schema and insert/write
- [ ] 9.2 Add protocol connector unit tests with simulated devices
- [ ] 9.3 Add mapping pipeline tests for transform correctness
- [ ] 9.4 Add API tests for CRUD, preview, and settings endpoints
- [ ] 9.5 Add Postgres partition migration test in CI or a local container
