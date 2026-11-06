## ADDED Requirements

### Requirement: Database output never blocks acquisition

Legacy database output (database target mappings handled by the database target writer) SHALL be decoupled from the runtime ingest path. Ingesting a collected value SHALL only update an in-memory latest-value table keyed by mapping and return without any database I/O. Each enabled database connector SHALL own one background worker that, for every mapping it serves, writes at most one row per effective write interval containing the latest value at write time. The effective write interval SHALL keep its existing resolution: the mapping write interval, then the connector default write interval, then 15 seconds. Grouped mappings SHALL keep their existing grouping semantics, with flushing performed by the owning connector worker. A slow or unreachable connector SHALL NOT delay point state persistence, runtime value events, Modbus Share writes, or any other connector.

#### Scenario: Unreachable connector does not slow ingest

- **GIVEN** a connector whose connection attempts never complete and 32 points mapped to it
- **WHEN** the runtime ingests one scan of 32 collected values
- **THEN** every value is ingested and its runtime value event is broadcast without waiting for the connector
- **AND** the connector status reports the failed write attempt

#### Scenario: Only the latest value per interval is written

- **GIVEN** an ungrouped mapping with a 15 second effective write interval
- **WHEN** the point reports values 10, 11 and 12 within one interval
- **THEN** the connector worker writes exactly one row with value 12 for that interval

##### Example: interval snapshots

| Values received in interval | Rows written |
| --------------------------- | ------------ |
| 10, 11, 12                  | one row, 12  |
| none                        | no row; history records the skip reason |

#### Scenario: Failed connector backs off

- **WHEN** consecutive write attempts to one connector fail
- **THEN** the worker retries after 1, 2, 4 and 8 seconds and then at most every 60 seconds
- **AND** each failure is recorded through the existing write outcome and write history
- **AND** values due during back-off are not queued; the next attempt writes only the then-latest values
