## MODIFIED Requirements

### Requirement: Guided workflow
The UI SHALL provide a guided workflow that keeps the complete datalink product flow inside `/studio` as the primary workspace: **Device -> SourceRule -> Tag review -> Output review/apply**.

The workflow SHALL ensure:

1. Operators establish device capability context before source-rule planning.
2. Tag and output review remain bound to the active source-rule revision.
3. Manual point, tag, or mapping construction tools are secondary and SHALL NOT be required to complete the normal product path.
4. Database and Local Modbus output readiness remain visible in the same workspace.
5. Failures remain localized to the owning step, selection, or output target.

#### Scenario: End-to-end guided configuration uses `/studio`
- **WHEN** an operator selects a device, saves a source rule, reviews tag candidates, and applies one or both output targets
- **THEN** the primary workflow stays inside `/studio`
- **AND** the operator does not need a separate manual point-first or mapping-first route to complete the normal path

### Requirement: Drag-drop mapping canvas
The UI MAY provide a drag-drop mapping canvas as a secondary engineering tool, but the primary `/studio` workflow SHALL NOT depend on manual drag-drop point-to-tag mapping.

#### Scenario: Primary workflow completes without drag-drop mapping
- **WHEN** an operator completes the normal `/studio` workflow
- **THEN** the system does not require opening a drag-drop mapping canvas to create the effective rule-driven tag or output state
