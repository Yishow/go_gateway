# source-rule-runtime Specification

## Purpose
Define the canonical runtime contract for persisted source rules so rule state survives restarts, activation is gated by probe success, derived relationships remain intact across enable/disable transitions, and live values surface back into the source-planning grid.
## Requirements
### Requirement: Persisted source-rule lifecycle
The system SHALL persist each source rule as a first-class record with device association, planning definition, revision metadata, and enabled/disabled state.

#### Scenario: Save rule definition
- **WHEN** an operator creates or updates a source rule
- **THEN** the system stores the rule in the database with its device association, planning inputs, current enabled/disabled state, and a new revision id
- **AND** the saved rule can be reloaded independently of the current browser session

#### Scenario: Restore rule state after restart
- **WHEN** the service restarts
- **THEN** the system reloads persisted source rules from the database together with the latest persisted revision metadata
- **AND** restores each rule's enabled/disabled state before resuming runtime behavior

### Requirement: Rule revisions produce downstream candidate snapshots
The system SHALL persist a dedicated downstream candidate snapshot for each saved source-rule revision before any later review/apply workflow can consume that revision.

Source-rule runtime SHALL own candidate generation and persistence for downstream review/apply workflows.

Each persisted candidate snapshot MUST include at least:

- `source_rule_id`
- `revision_id`
- `candidate_type`
- `payload`
- `status`
- `generated_at`

#### Scenario: Save rule writes candidate snapshot
- **WHEN** a source rule is created or updated
- **THEN** the system persists the new revision and synchronizes runtime rule-owned point state
- **AND** persists the candidate snapshot for that revision before later review flows can query it

#### Scenario: Candidate snapshot survives restart
- **WHEN** the service restarts after a source-rule revision has generated candidate snapshots
- **THEN** the system restores the persisted candidate snapshots for that revision
- **AND** downstream review/apply workflows can query them without regenerating an incompatible contract

### Requirement: Rule-owned identity and signature are canonical across revisions
The system SHALL compare rule-derived tags, mappings, and output candidates using one canonical rule-owned identity and one deterministic proposed-signature model.

A rule-owned identity MUST include at least:

- `source_rule_id`
- `candidate_type`
- `derived_from_rule_address`
- target binding scope fields needed to distinguish one downstream object from another

The proposed signature MUST be derived from the effective candidate payload for that identity.

#### Scenario: Same identity with changed payload is detectable
- **WHEN** a later rule revision preserves a downstream object's rule-owned identity but changes its effective payload
- **THEN** the system preserves the same rule-owned identity for comparison
- **AND** computes a different proposed signature for the changed candidate

### Requirement: Snapshot completeness records ready, blocked, and deferred target state explicitly
The system SHALL treat a candidate snapshot as complete only when it contains the current tag candidate set and an explicit state entry for each supported downstream target.

Each downstream target entry MUST be present as one of `ready`, `blocked`, or `deferred`, and blocked or deferred entries MUST include the reason.

Snapshot completeness SHALL NOT require every downstream target to be apply-ready.

#### Scenario: Missing database context still yields a complete snapshot
- **WHEN** a source-rule revision generates tag candidates but the database target lacks connector context
- **THEN** the snapshot is still marked complete
- **AND** the database target entry is persisted as `blocked` or `deferred` with its reason instead of being omitted

### Requirement: Rule changes preserve applied downstream state until explicit reapply
The system SHALL mark changed downstream assets as `out_of_sync` rather than silently replacing applied tag or output state.

#### Scenario: Applied downstream asset becomes out of sync
- **WHEN** a saved source-rule revision changes the derived downstream result
- **THEN** the system preserves the previously applied asset
- **AND** marks it as `out_of_sync` until the operator explicitly reapplies the new candidate

### Requirement: Rule revision rollback restores the last applied downstream state
The system SHALL treat rollback as a state-based restore of the last successfully applied downstream state per target.

#### Scenario: Unapplied revision is rolled back
- **WHEN** a newer source-rule revision is abandoned or rolled back before downstream apply is completed
- **THEN** the system removes unapplied candidate snapshots for the rolled-back revision
- **AND** preserves the last successfully applied downstream state for tags, database output, and Local Modbus output independently

### Requirement: Rule enablement controls collection without deleting derived relationships
The system SHALL treat rule enablement as collection control only and SHALL preserve Point, Tag, Mapping, and Output relationships when a rule is disabled.

#### Scenario: Disable rule preserves derived records
- **WHEN** an operator disables a rule
- **THEN** the system stops collection for that rule
- **AND** preserves the derived Point, Tag, Mapping, and Output relationships in storage

#### Scenario: Re-enable rule resumes existing definition
- **WHEN** an operator re-enables a previously disabled rule
- **THEN** the system resumes collection using the stored rule definition and derived relationships
- **AND** does not require the operator to recreate the rule or downstream bindings

### Requirement: Rule activation is gated by protocol probe success
The system SHALL use the device readiness contract to distinguish planning from runtime activation and collection.

#### Scenario: Connect succeeds but probe fails
- **WHEN** a device test reports transport connect success and protocol probe failure
- **THEN** the operator MAY save the device and save source-rule planning state
- **AND** the system SHALL reject rule activation and data collection until probe succeeds

### Requirement: Rule-driven live values appear in the source grid
The system SHALL display parsed live values for enabled rules in the Step 2 source grid.

#### Scenario: Enabled rule updates grid
- **WHEN** runtime collection succeeds for an enabled source rule
- **THEN** the Step 2 grid shows the parsed value, timestamp, and error state for the rule-derived span
- **AND** the displayed value follows the rule's configured data-type and merge semantics
