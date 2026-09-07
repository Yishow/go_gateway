# durable-recording-delivery Specification

## Purpose

TBD - created by archiving change 'add-durable-recording-delivery'. Update Purpose after archive.

## Requirements

### Requirement: Durable acceptance and atomic calculation
The system MUST acknowledge recording intake only after durable local commit and atomically persist calculation progress and its output intents.

#### Scenario: Process terminates after intake acknowledgement
- **WHEN** the process restarts after acknowledging a sample but before delivery
- **THEN** that sample or its pending output remains recoverable without duplicate usage.

---
### Requirement: Destination-specific idempotent effects
The system SHALL use stable record identities and transactional destination receipts where supported, without claiming universal exactly-once behavior.

#### Scenario: External commit response lost
- **WHEN** a destination commits a record but its response is lost
- **THEN** retry checks its durable identity and produces no duplicate managed record.

---
### Requirement: Independent bounded delivery workers
The system SHALL isolate destination failures, preserve per-series ordering and bound concurrency, retry delays and queues.

#### Scenario: One database is offline
- **WHEN** destination A is offline and destination B is healthy
- **THEN** B continues within configured capacity and A retains its own backlog and failure status.

---
### Requirement: Visible capacity exhaustion
The system MUST warn before quota exhaustion and MUST NOT silently discard unacknowledged records or fabricate complete history.

#### Scenario: Durable storage reaches its hard limit
- **WHEN** the local backlog cannot accept more records
- **THEN** affected recording intake becomes blocked with data-loss risk and exact affected scope
- **AND** any still-running device collection or Share output is reported separately.

---
### Requirement: Revision-safe delivery and apply
The system SHALL freeze destination and plan semantics for accepted records and use the existing workspace activation barrier for changes.

#### Scenario: Database endpoint changes with backlog
- **WHEN** a new plan revision targets a different endpoint
- **THEN** existing pending records are not automatically redirected and the operator receives a migration choice.

---
### Requirement: Recovery and retention protect evidence
The system MUST retain journal entries, checkpoints, receipts and calculation evidence according to all pending consumers and replay horizons.

#### Scenario: Cleanup while one consumer is behind
- **WHEN** a retention pass finds unprocessed or unacknowledged records
- **THEN** it does not delete them as ordinary expired history and exposes quota implications.

---
### Requirement: Delivery status represents actual stages
The system SHALL report collecting, local_durable, delivered and verified independently, including pending count, oldest age and actionable errors.

#### Scenario: Device healthy but database blocked
- **WHEN** acquisition is successful but target permissions fail
- **THEN** the UI does not label end-to-end recording successful and provides the target-specific repair action.

---
### Requirement: Safe shutdown and worker ownership
The system MUST use bounded shutdown, durable recovery and exclusive worker ownership for each delivery partition.

#### Scenario: Two workers attempt the same partition
- **WHEN** a recovering worker and a stale worker overlap
- **THEN** lease/fencing or equivalent transactional ownership prevents simultaneous conflicting progress.
