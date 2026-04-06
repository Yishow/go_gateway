## MODIFIED Requirements

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

### Requirement: Rule activation is gated by protocol probe success
The system SHALL use the device readiness contract to distinguish planning from runtime activation and collection.

#### Scenario: Connect succeeds but probe fails
- **WHEN** a device test reports transport connect success and protocol probe failure
- **THEN** the operator MAY save the device and save source-rule planning state
- **AND** the system SHALL reject rule activation and data collection until probe succeeds

## ADDED Requirements

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
