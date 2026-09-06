## ADDED Requirements

### Requirement: Database kind switch auto-populates defaults

When the operator changes the `kind` of a connector in the Connector Pool row, the system SHALL update the connection attributes according to the smart default synchronization rules:
- `port` SHALL be updated to the default port of the target database kind (PostgreSQL: 5432, MySQL: 3306, SQL Server: 1433, SQLite: 0).
- `host` SHALL default to `127.0.0.1` if empty or previously set to a default loopback address (`127.0.0.1`, `localhost`, `tsdb.internal`); custom non-local host values SHALL be preserved.
- `schema` SHALL update to the default schema of the target database kind (PostgreSQL: `public`, SQL Server: `dbo`, MySQL/SQLite: empty string).
- `username` SHALL update to the default username of the target database kind if currently empty or set to a standard default username (PostgreSQL: `postgres`, MySQL: `root`, SQL Server: `sa`, SQLite: empty string).
- For SQLite, `host`, `port`, `username`, and `password` SHALL be cleared, and `database` SHALL default to `gateway.db` if empty.

#### Scenario: Switching from PostgreSQL to MySQL
- **GIVEN** a connector with kind `postgres`, host `127.0.0.1`, port `5432`, schema `public`, username `postgres`
- **WHEN** the operator changes kind to `mysql`
- **THEN** port SHALL become `3306`
- **AND** host SHALL remain `127.0.0.1`
- **AND** schema SHALL become `""`
- **AND** username SHALL become `root`

#### Scenario: Switching from PostgreSQL to SQL Server with custom remote host
- **GIVEN** a connector with kind `postgres`, host `192.168.1.50`, port `5432`
- **WHEN** the operator changes kind to `sqlserver`
- **THEN** host SHALL remain `192.168.1.50`
- **AND** port SHALL become `1433`
- **AND** schema SHALL become `dbo`

### Requirement: Protected connector password input

The connector password input field SHALL maintain isolated controlled state during user typing, and SHALL NOT be cleared or overwritten by parent component re-renders or background autosave operations.

#### Scenario: Typing password without intermediate purge
- **GIVEN** a connector row in the settings page
- **WHEN** the operator enters characters into the password field
- **THEN** the entered characters SHALL persist in the input and state without being cleared on debounce or autosave ticks
