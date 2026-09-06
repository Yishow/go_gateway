## ADDED Requirements

### Requirement: Generated schema statements use dialect-valid syntax

Schema generation SHALL emit statements accepted by the target database dialect. For MySQL the unique-index statement SHALL NOT use the IF NOT EXISTS clause, which MySQL does not accept for index creation. For SQLite and PostgreSQL the existing IF NOT EXISTS form SHALL be retained. Duplicate index creation SHALL continue to be prevented by the inspected column metadata — a column already reported as a primary key or unique SHALL produce no index statement — and by deduplication of identical index names within one generation batch.

#### Scenario: MySQL unique index statement omits IF NOT EXISTS

- **GIVEN** a MySQL connector with an upsert mapping whose timestamp column exists but is neither a primary key nor unique
- **WHEN** schema generation runs
- **THEN** the generated statement creates a unique index without the IF NOT EXISTS clause
- **AND** executing the generated statements against MySQL succeeds

#### Scenario: SQLite and PostgreSQL keep the guarded form

- **GIVEN** a SQLite or PostgreSQL connector with an upsert mapping whose timestamp column is neither a primary key nor unique
- **WHEN** schema generation runs
- **THEN** the generated statement creates a unique index using the IF NOT EXISTS clause

#### Scenario: An already-unique timestamp column produces no index statement

- **GIVEN** a connector whose inspected timestamp column is reported as a primary key or unique
- **WHEN** schema generation runs
- **THEN** no unique-index statement is generated for that column

##### Example: unique index statement by dialect and column state

| Dialect | Timestamp column state | Index statement generated |
| ------- | ---------------------- | ------------------------- |
| MySQL | plain column | CREATE UNIQUE INDEX without IF NOT EXISTS |
| MySQL | primary key or unique | none |
| PostgreSQL | plain column | CREATE UNIQUE INDEX with IF NOT EXISTS |
| SQLite | plain column | CREATE UNIQUE INDEX with IF NOT EXISTS |

### Requirement: Upsert mappings require a guaranteed unique key on MySQL

MySQL upsert delivery relies on a unique or primary key covering the timestamp column; without it the upsert clause silently degrades to append-only inserts. Schema generation SHALL guarantee that key for every enabled MySQL upsert mapping, either by reporting the column as already unique or by emitting a statement that creates the unique index. When the index cannot be created, schema generation SHALL fail and SHALL surface the underlying database error rather than completing successfully.

#### Scenario: MySQL upsert mapping ends with a unique key present

- **GIVEN** an enabled MySQL upsert mapping against an existing table whose timestamp column has no unique key
- **WHEN** schema generation completes successfully
- **THEN** the timestamp column is covered by a unique index

#### Scenario: Unique index creation failure fails schema generation

- **GIVEN** an existing MySQL table containing duplicate values in the timestamp column
- **WHEN** schema generation attempts to create the unique index
- **THEN** schema generation is recorded as failed
- **AND** the underlying database error is preserved in the reported failure

### Requirement: Schema generation resolves the default schema from the connector

When a target mapping record carries an empty table schema, schema generation SHALL resolve the default schema from the connector, using the same resolution the mapping create and update paths use. For MySQL that resolution SHALL yield the database name from the connector connection configuration. Schema generation SHALL NOT fall back to a kind-only default that cannot match the inspected table set.

#### Scenario: MySQL mapping with an empty table schema targets the connection database

- **GIVEN** a MySQL connector whose connection configuration names database gateway_metrics
- **AND** an enabled mapping whose table schema is empty
- **WHEN** schema generation runs
- **THEN** the generated statements qualify the table with gateway_metrics
- **AND** no generated statement qualifies a table with the kind-only default schema name

#### Scenario: Mapping with an empty schema matches existing inspected tables

- **GIVEN** a MySQL connector whose target table already exists in the connection database
- **AND** an enabled mapping for that table whose table schema is empty
- **WHEN** schema generation runs
- **THEN** the existing table is matched and no create-table statement is generated for it

##### Example: resolved schema by connector kind

| Connector kind | Connection configuration | Mapping table schema | Resolved schema |
| -------------- | ------------------------ | -------------------- | --------------- |
| MySQL | database=gateway_metrics | empty | gateway_metrics |
| MySQL | database=gateway_metrics | reporting | reporting |
| PostgreSQL | not applicable | empty | public |

### Requirement: MySQL connections fail closed rather than sending cleartext passwords

MySQL connection assembly SHALL NOT enable cleartext password authentication by default, and SHALL NOT silently fall back to an unencrypted connection when cleartext authentication is enabled. The same rules SHALL apply to both the target connection and the administrative connection used by connector probing, from a single shared resolution.

The resolution SHALL be:

- An explicit TLS value in the connection configuration SHALL be used as given.
- Otherwise, when TLS is requested through the use-TLS option, the TLS mode SHALL be an encrypting mode that does not permit falling back to plaintext.
- Otherwise, when cleartext authentication is explicitly enabled, the TLS mode SHALL be an encrypting mode that does not permit falling back to plaintext.
- Otherwise, no TLS mode SHALL be forced.

#### Scenario: Default configuration does not enable cleartext authentication

- **GIVEN** a MySQL connector configuration that does not set the cleartext password option
- **WHEN** the connection descriptor is assembled
- **THEN** cleartext password authentication is disabled

#### Scenario: Explicitly enabled cleartext requires an encrypted transport

- **GIVEN** a MySQL connector configuration that explicitly enables cleartext password authentication and sets no explicit TLS value
- **WHEN** the connection descriptor is assembled
- **THEN** the TLS mode is an encrypting mode
- **AND** falling back to an unencrypted connection is not permitted
- **AND** a server without TLS causes the connection to fail with the existing connection error classification

#### Scenario: An explicit TLS value is honoured

- **GIVEN** a MySQL connector configuration that sets an explicit TLS value
- **WHEN** the connection descriptor is assembled
- **THEN** that TLS value is used unchanged

#### Scenario: Probe and target connections share the resolution

- **GIVEN** any MySQL connector configuration
- **WHEN** both the target connection descriptor and the probe administrative connection descriptor are assembled
- **THEN** both carry the same cleartext and TLS resolution

##### Example: TLS and cleartext resolution

| Explicit TLS value | use-TLS requested | cleartext explicitly enabled | Cleartext | TLS mode | Plaintext fallback |
| ------------------ | ----------------- | ---------------------------- | --------- | -------- | ------------------ |
| none | no | no | disabled | none forced | not applicable |
| none | yes | no | disabled | encrypting | not permitted |
| none | no | yes | enabled | encrypting | not permitted |
| set | any | any | as configured | the explicit value | as implied by that value |

### Requirement: Connector probing attempts database creation once per kind

Connector probing SHALL attempt automatic database creation at most once per probe for a given connector kind. For MySQL, where the connection helper already performs create-then-retry internally, the probe SHALL NOT contain a second create-and-retry branch.

#### Scenario: MySQL probe does not retry database creation twice

- **GIVEN** a MySQL connector whose target database does not exist
- **WHEN** the connector is probed
- **THEN** database creation is attempted once
- **AND** the connection is retried once after creation

#### Scenario: A non-missing-database MySQL error is not treated as a creation candidate

- **GIVEN** a MySQL connector that fails to connect for a reason other than a missing database
- **WHEN** the connector is probed
- **THEN** no database creation is attempted
- **AND** the original failure is classified and reported
