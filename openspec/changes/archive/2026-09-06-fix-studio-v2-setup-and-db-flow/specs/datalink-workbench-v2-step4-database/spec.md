## ADDED Requirements

### Requirement: Select existing connector from pool in Step 4

Step 4 Database configuration SHALL allow the operator to select an existing connector configured in the Settings Connector Pool. Selecting a connector SHALL immediately populate the current step connector parameters (kind, name, host, port, database, username, schema, table) without requiring manual re-entry.

#### Scenario: Populate Step 4 from Connector Pool
- **GIVEN** the Settings Connector Pool contains a connector `conn-mysql-prod` with kind `mysql`, host `10.0.0.5`, port `3306`, database `factory_data`
- **WHEN** the operator selects `conn-mysql-prod` in Step 4
- **THEN** Step 4 connector parameters SHALL immediately reflect kind `mysql`, host `10.0.0.5`, port `3306`, and database `factory_data`

### Requirement: Step 4 kind switch auto-populates defaults

When the operator manually changes database `kind` in Step 4, the form SHALL apply the smart default preset rules aligned with the settings connector pool:
- Target standard default port and schema SHALL be assigned.
- Non-local custom host SHALL be preserved.
- Local loopback host SHALL default to `127.0.0.1`.

#### Scenario: Switching kind in Step 4
- **GIVEN** Step 4 connector has kind `postgres` and host `127.0.0.1`
- **WHEN** the operator changes kind to `mysql`
- **THEN** port SHALL be updated to `3306`
- **AND** username SHALL be updated to `root`
