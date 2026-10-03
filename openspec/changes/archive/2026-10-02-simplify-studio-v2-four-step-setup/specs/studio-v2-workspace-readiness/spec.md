## MODIFIED Requirements

### Requirement: Workspace readiness evaluates all persisted setup steps

The system SHALL evaluate persisted Step 1, Step 2, Step 3, and Step 4 setup state into one normalized workspace readiness result. Initial navigation, reload and activation MUST use the same persisted device validation and capability facts, but unresolved readiness MUST NOT prevent saving or navigating offline drafts; it gates only dependent read/activation actions. For migrated database outputs, readiness SHALL evaluate the canonical WriteGroup membership, destination/schema/identity capability and revisions rather than competing plan or manual-target counts; existing workspace/settings/readiness-token and Share gates MUST remain effective.

#### Scenario: Readiness returns one normalized issue set

- **WHEN** a workspace has persisted devices, source rules, mappings, and database targets with a mix of complete and incomplete state
- **THEN** the readiness result returns one normalized issue set with stable issue codes, severity, and affected step or device scope
- **AND** the readiness result SHALL distinguish blocking issues from warning issues

##### Example: one blocking issue and one warning issue

| Scope | Step | Issue code | Severity | Meaning |
| ----- | ---- | ---------- | -------- | ------- |
| dev-A | Step 1 | device-probe-required | blocking | device is saved but not probe-ready for activation |
| db-main | Step 4 | database-connector-unreachable | warning | database sink is configured but currently unreachable |

#### Scenario: Concurrent stale activation
- **WHEN** the group or connector changes after readiness was issued and a client submits the old token/revisions
- **THEN** activation rejects the stale request before changing runtime writer ownership.

#### Scenario: Basic and advanced requirements stay separate
- **WHEN** a basic group has valid persisted Tags but no aggregation or report policy
- **THEN** those absent optional policies produce no basic-writing blocker; unsupported row identity or missing saved members still blocks with an actionable scoped issue.
