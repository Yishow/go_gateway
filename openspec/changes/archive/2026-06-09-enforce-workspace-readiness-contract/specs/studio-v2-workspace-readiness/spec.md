## ADDED Requirements

### Requirement: Workspace readiness evaluates all persisted setup steps

The system SHALL evaluate persisted Step 1, Step 2, Step 3, and Step 4 setup state into one normalized workspace readiness result.

#### Scenario: Readiness returns one normalized issue set

- **WHEN** a workspace has persisted devices, source rules, mappings, and database targets with a mix of complete and incomplete state
- **THEN** the readiness result returns one normalized issue set with stable issue codes, severity, and affected step or device scope
- **AND** the readiness result SHALL distinguish blocking issues from warning issues

##### Example: one blocking issue and one warning issue

| Scope | Step | Issue code | Severity | Meaning |
| ----- | ---- | ---------- | -------- | ------- |
| dev-A | Step 1 | device-probe-required | blocking | device is saved but not probe-ready for activation |
| db-main | Step 4 | database-connector-unreachable | warning | database sink is configured but currently unreachable |

### Requirement: Activation blocks on blocking readiness issues

The system SHALL reject activation when the workspace readiness result contains blocking issues.

#### Scenario: Activation returns blocking readiness issues

- **WHEN** an operator attempts activation while the readiness result still contains one or more blocking issues
- **THEN** activation is rejected before runtime start
- **AND** the response surfaces the blocking readiness issue codes and scope instead of failing later inside runtime

##### Example: activation rejects one device blocker

- **GIVEN** workspace readiness includes blocking issue device-probe-required for dev-A
- **WHEN** the operator starts activation
- **THEN** activation returns that blocking issue for dev-A and runtime does not start for the workspace

### Requirement: Readiness updates when persisted setup state changes

The system SHALL recompute workspace readiness from persisted state whenever persisted Step 1 through Step 4 setup data changes.

#### Scenario: Saving setup data updates readiness

- **WHEN** the operator persists a device, source rule, mapping, or database target change
- **THEN** the next readiness result reflects the newly persisted state
- **AND** the shell SHALL NOT continue showing a stale readiness conclusion from the prior persisted snapshot

##### Example: saving a missing mapping clears one blocker

- **GIVEN** readiness currently reports mapping-missing as a blocking issue for point pt-A
- **WHEN** the operator persists the missing mapping for pt-A
- **THEN** the next readiness result removes mapping-missing for pt-A and updates the shell summary accordingly
