## ADDED Requirements

### Requirement: Rule revisions produce downstream candidate snapshots
The system SHALL persist a dedicated downstream candidate snapshot for each saved source-rule revision before tag or output apply is allowed.

Source-rule runtime SHALL own candidate generation and persistence for downstream review/apply workflows.

Each persisted candidate snapshot MUST include at least:

- `source_rule_id`
- `revision_id`
- `candidate_type`
- `payload`
- `status`
- `generated_at`

#### Scenario: Candidate snapshot blocks apply until ready
- **WHEN** a source rule is created or updated
- **THEN** the system persists the new revision and synchronizes runtime point state
- **AND** blocks downstream apply until the candidate snapshot for that revision is complete

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

An applied downstream asset is `out_of_sync` only when the same rule-owned identity exists but its proposed signature differs from the last applied signature.

#### Scenario: Same identity with changed payload becomes out of sync
- **WHEN** a later rule revision preserves a downstream object's rule-owned identity but changes its effective payload
- **THEN** the system preserves the existing applied asset
- **AND** marks the new candidate as `out_of_sync` because the proposed signature changed

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
- **WHEN** a saved source rule revision changes the derived downstream result
- **THEN** the system preserves the previously applied asset
- **AND** marks it as `out_of_sync` until the operator explicitly reapplies the new candidate

### Requirement: Rule revision rollback restores the last applied downstream state
The system SHALL treat rollback as a state-based restore of the last successfully applied downstream state per target.

#### Scenario: Unapplied revision is rolled back
- **WHEN** a newer source-rule revision is abandoned or rolled back before downstream apply is completed
- **THEN** the system removes unapplied candidate snapshots for the rolled-back revision
- **AND** preserves the last successfully applied downstream state for tags, database output, and Local Modbus output independently
