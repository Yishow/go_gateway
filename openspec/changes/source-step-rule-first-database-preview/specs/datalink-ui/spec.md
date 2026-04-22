## MODIFIED Requirements

### Requirement: Guided workflow
The UI SHALL provide a guided workflow that keeps the complete datalink product flow inside `/studio` as the primary workspace: **Device -> SourceRule -> Tag review -> Output review/apply**.

The workflow SHALL ensure:

1. Operators establish device capability context before source-rule planning.
2. Tag and output review remain bound to the active source-rule revision.
3. Manual point, tag, or mapping construction tools are secondary and SHALL NOT be required to complete the normal product path.
4. Database and Local Modbus output readiness remain visible in the same workspace.
5. Failures remain localized to the owning step, selection, or output target.
6. Step 2 SHALL surface the current active source rule as the primary planning unit, so the operator can immediately understand what is currently planned before moving into Tag review.
7. When the active output target is `database`, Step 2 SHALL reflect database-aware planning guidance without turning Step 2 into connector or schema setup.

#### Scenario: End-to-end guided configuration uses `/studio`
- **WHEN** an operator selects a device, saves a source rule, reviews tag candidates, and applies one or both output targets
- **THEN** the primary workflow stays inside `/studio`
- **AND** the operator does not need a separate manual point-first or mapping-first route to complete the normal path

#### Scenario: Source step exposes the active rule before downstream review
- **WHEN** an operator opens Step 2 with an active source rule
- **THEN** the UI shows the current rule’s planning summary before the operator moves to Tag review
- **AND** the operator can identify the rule’s address coverage, type intent, and current planning state without reading the full canvas first
