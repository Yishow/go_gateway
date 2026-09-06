## ADDED Requirements

### Requirement: Automatic database creation on probe and initialization

When probing a database connector or establishing a target connection, if the target database does not exist (such as MySQL Error 1049 `ER_BAD_DB_ERROR` or PostgreSQL error `3D000`), the backend probe service SHALL attempt to automatically create the database on the server:
- For MySQL, connect without specifying the target database and execute `CREATE DATABASE IF NOT EXISTS <dbname>`.
- For PostgreSQL, connect to the administrative database (`postgres` or `template1`) and execute `CREATE DATABASE <dbname>`.
- Upon successful database creation, the backend SHALL re-establish the connection to the newly created database and complete verification.

#### Scenario: MySQL connection test with non-existent database
- **GIVEN** a MySQL server where database `gateway_metrics` does not yet exist
- **WHEN** the backend probes the connector configuration
- **THEN** the backend SHALL automatically execute database creation and return status `ready`

#### Scenario: Database creation error due to missing privileges
- **GIVEN** a database user without `CREATE DATABASE` privileges on the target server
- **WHEN** the backend attempts to auto-create the database and receives an access denied error
- **THEN** the probe SHALL fail with a clear authorization error message and SHALL NOT panic
