## MODIFIED Requirements

### Requirement: Runtime dashboard renders only backend-backed data

The runtime dashboard SHALL render only values, counters, and summaries backed by the runtime backend contract.

#### Scenario: Missing backend data produces explicit empty state

- **WHEN** the runtime backend does not have real snapshot or stream data for the selected device
- **THEN** the dashboard renders an explicit empty, degraded, or error state for that missing data
- **AND** the dashboard SHALL NOT synthesize a healthy runtime summary from defaults or unrelated state

##### Example: selected device has no snapshot yet

- **GIVEN** dev-A is selected and runtime backend has not produced any snapshot for dev-A
- **WHEN** the dashboard loads
- **THEN** the dashboard shows a waiting or empty state for dev-A instead of a healthy default summary

#### Scenario: No failure is not proof of database delivery
- **WHEN** no failure_stage is present but matching committed database effect evidence is absent
- **THEN** diagnostics report only that no error was detected and distinguish acquisition, waiting for the first record and confirmed SQL writing; successful Share or reads do not imply all points reached a database.
