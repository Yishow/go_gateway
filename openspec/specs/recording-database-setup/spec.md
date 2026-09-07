# recording-database-setup Specification

## Purpose

TBD - created by archiving change 'add-recording-plans'. Update Purpose after archive.

## Requirements

### Requirement: Capability-backed database preparation
The system SHALL expose verified adapter capabilities and actual schema metadata, separately from network and authentication success.

#### Scenario: Connection succeeds but schema inspection is forbidden
- **WHEN** the account can connect but cannot inspect the target
- **THEN** the result states that limitation without returning sample columns or claiming write readiness.

---
### Requirement: Revision-bound schema preview and explicit creation
The system MUST bind schema creation to a fresh preview token containing workspace, plan and connector identity revisions and perform backend validation again.

#### Scenario: Target changed after preview
- **WHEN** the database or plan changes after schema preview
- **THEN** schema apply is rejected as stale before executing statements.

---
### Requirement: Credential and path safety
The system SHALL reuse stored credentials by backend reference and invalidate tests after identity changes without returning secrets to the browser.

#### Scenario: Reuse a saved connection
- **WHEN** the operator selects an existing unchanged connection
- **THEN** no password re-entry is required solely because the list response masks it
- **AND** changing the endpoint requires a new credential validation.

---
### Requirement: Truthful explicit test writes
The system SHALL separate data preview, confirmed write and readback verification and retain idempotent test identities.

#### Scenario: Write succeeds without read permission
- **WHEN** a confirmed test write succeeds but the account cannot query the result
- **THEN** the result is written_unverified, not verified.

#### Scenario: Test action retried
- **WHEN** the same confirmed test token is retried after an unknown network outcome
- **THEN** the system checks its durable identity and does not insert a duplicate test record.
