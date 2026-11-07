## ADDED Requirements

### Requirement: Exact persisted successor schema confirmation replay

A completed successor schema confirmation whose operation and verified digest have already been saved SHALL replay its durable result without a fresh destination capability probe or DDL. Replay SHALL validate the canonical operation token and action, immutable reservation, current scope and revisions, connector revision, stored schema layout and complete proof through the existing confirmation authority. It SHALL NOT enable collection, move the canonical head, or weaken first-confirmation checks.

#### Scenario: Destination is offline after confirmed schema success

- **WHEN** an identical schema-confirm request is retried after its verified operation and digest were saved, and the destination is offline
- **THEN** the existing success is returned with zero destination probes or DDL and unchanged local and destination state

##### Example: Lost acknowledgement with saved proof

- **GIVEN** a reserved SQL Server successor has a succeeded schema operation and saved digest in plain or native_monthly mode
- **WHEN** the original pre-confirm request is retried while the capability inspector refuses new probes
- **THEN** both modes return the recorded success, inspector call count remains zero, and no state changes

#### Scenario: Replay identity or current scope differs

- **WHEN** a replay changes the digest, schema operation, target connector revision, lineage revision, or uses an alias token for the same recorded operation
- **THEN** confirmation is rejected with no probe, DDL, or state mutation

##### Example: Alias token cannot authorize the canonical operation

- **GIVEN** a copy of the saved preview token is stored under a different token identifier
- **WHEN** schema-confirm supplies that alias while retaining the original operation and digest
- **THEN** the request is rejected rather than treated as an exact replay
