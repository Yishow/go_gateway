## MODIFIED Requirements

### Requirement: Database delivery operates only on the live enabled mapping set
The system SHALL build new schema plans and accept new recording inputs only from the current live enabled target set. Legacy immediate writes SHALL continue to use that live set. Records already accepted by the durable recording path SHALL retain their validated destination and definition snapshots and remain deliverable after a later disable or edit unless explicitly cancelled with an auditable affected-record decision. Hidden stale mappings MUST NOT create new delivery records or block unrelated live targets.

#### Scenario: Hidden stale target does not block live delivery
- **WHEN** a connector contains an old hidden mapping no longer tied to a live rule-owned point or tag
- **THEN** new schema planning and intake ignore that mapping
- **AND** unrelated live targets remain eligible for activation and delivery.

#### Scenario: Accepted records survive later disable
- **WHEN** valid durable records are pending and their source rule is subsequently disabled
- **THEN** no new samples are accepted for that disabled scope and previously accepted records remain pending for their original destination.

#### Scenario: Explicit cancellation
- **WHEN** an operator explicitly cancels pending records after reviewing scope and count
- **THEN** cancellation is audited and affected reports expose incomplete coverage rather than claiming successful delivery.
