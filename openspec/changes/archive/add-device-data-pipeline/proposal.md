# Change: Add device data pipeline with protocol ingestion, tagging, and time-series storage

## Why
The project needs a complete path to collect device data through multiple protocols, map it to global tags with transforms, and store it in a time-series database with a usable mapping UI. Historical queries will be handled by an external Next.js CMS.

## What Changes
- Add device and point registry with protocol-specific connection settings
- Add protocol connectors for Modbus (TCP/RTU/UDP), FATEK FBs, Mitsubishi MC 3E, and MQTT ingest
- Add collection scheduler with retries, batching, and quality flags
- Add global tag dictionary and mapping pipeline with full transform rules
- Add time-series storage schema with Postgres partitions and SQLite test support
- Add API endpoints for management, mapping preview, and write settings (no history query API)
- Add UI workflow for drag-drop mapping, transform builder, preview, and write precision settings
- Add database schema and indexes to support external CMS history queries

## Impact
- Affected specs: device-registry, point-catalog, protocol-connectors, collection-scheduler, tag-dictionary, mapping-pipeline, timeseries-storage, datalink-api, datalink-ui
- Affected code: internal/protocol, internal/collector, internal/tag, internal/mapping, internal/storage, internal/api, web/ui
