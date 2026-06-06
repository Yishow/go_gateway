## ADDED Requirements

### Requirement: Runtime diagnostics expose collector to database delivery stages

The system SHALL expose runtime diagnostics that cover collector, mapping, runtime projection, and database delivery stages for the affected scope.

#### Scenario: Diagnostics show the stage where delivery stopped

- **WHEN** a point value is collected but later fails during mapping, runtime projection, or database delivery
- **THEN** the diagnostics result identifies the latest successful stage and the failing stage for that scope
- **AND** the operator SHALL NOT have to infer the break only from missing output data

##### Example: diagnostics stop at database delivery

- **GIVEN** pt-A was collected and mapped successfully but db-main rejected the write
- **WHEN** diagnostics are queried for pt-A or its selected device scope
- **THEN** diagnostics report collector and mapping as successful and database delivery as the failing stage

### Requirement: Runtime diagnostics keep recent success and failure context

The system SHALL keep recent success and failure context for the operator-facing diagnostics contract.

#### Scenario: Diagnostics show last success and last failure timestamps

- **WHEN** runtime experiences successful and failed processing over time
- **THEN** the diagnostics contract can report the last success timestamp, last failure timestamp, and latest failure reason for the relevant scope
- **AND** the operator can tell whether the problem is current or historical

##### Example: diagnostics distinguish historical and current failures

- **GIVEN** dev-A last succeeded at 10:00:00Z and last failed at 10:03:00Z because of write timeout
- **WHEN** diagnostics are requested for dev-A
- **THEN** the response includes both timestamps and the latest failure reason so the operator can judge freshness of the fault
