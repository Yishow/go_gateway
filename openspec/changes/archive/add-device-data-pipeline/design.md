## Context
The project aims to collect industrial device data over multiple protocols, transform values into globally defined tags, and store them as a time-series record with a UI to manage mapping and preview. Historical queries are handled by an external Next.js CMS. There are no existing specs or active changes in OpenSpec yet.

## Goals / Non-Goals
- Goals:
  - Provide an end-to-end data path from protocol collection to time-series storage
  - Provide a global tag dictionary with metadata and lifecycle states
  - Provide a drag-drop mapping UI with transform steps and preview
  - Provide a durable write path with Postgres partitions for external history queries
  - Support SQLite as a test backend with the same logical schema
- Non-Goals:
  - Role-based access control or multi-tenant isolation
  - CSV or Excel import/export for mappings
  - Built-in history query API or UI
  - Advanced analytics beyond mapping preview
  - Automatic data retention or deletion policies
  - Mandatory external time-series extensions (deferred)

## Decisions
- Decision: Use a global tag dictionary with unique tag keys and optional namespace segments to avoid per-device tag duplication.
- Decision: Tag keys are not required to follow a strict hierarchy, but use a safe ASCII character set for portability and are unique case-insensitively.
- Decision: Use an ordered transform pipeline with explicit input/output types, including a formula step for advanced rules.
- Decision: Formula expressions use infix operators with snake_case function names and boolean operators and/or/not.
- Decision: Persist raw and final values with a quality flag for traceability and troubleshooting.
- Decision: Use a narrow time-series table with typed value columns (value_num, value_text, value_bool) and tag_id + ts indexes.
- Decision: Use Postgres range partitions for time-series storage; default to monthly partitions and allow configuration changes if needed.
- Decision: Use a non-partitioned table in SQLite for tests.
- Decision: Timestamp precision is a global setting configured via UI and applied at write time.
- Decision: Use a batch writer with size and interval flush to balance latency and throughput.
- Decision: Use a guided UI workflow and a drag-drop mapping canvas with preview.
- Decision: Do not provide built-in history query APIs; the CMS reads from the database using the published schema.

## Alternatives considered
- Alternative: Wide tables per device or per line. Rejected due to schema churn and operational overhead.
- Alternative: Store values as JSON blobs. Rejected due to query complexity and weak typing.
- Alternative: Use TimescaleDB. Deferred to keep dependencies minimal in the first release.

## Risks / Trade-offs
- Partition management adds operational complexity; mitigate with automated creation and monitoring.
- Transform rules can become complex; mitigate with safe expression evaluation and clear UI validation.
- High write volume can overwhelm storage; mitigate with batching and backpressure.
- UI mapping complexity can impact usability; mitigate with templates and progressive disclosure.

## Migration Plan
- Create baseline schema for device registry, tag dictionary, mapping rules, and time-series storage.
- Apply non-partitioned migrations for SQLite in test environments.
- Apply partitioned migrations for Postgres in production environments.
- Rollback by disabling schedulers and stopping new writes; schema changes are additive.

## Open Questions
