## ADDED Requirements

### Requirement: Source desk switching preserves guided-flow continuity
The UI SHALL preserve guided-flow continuity when the operator switches between Step 2 Source desk modes.

Step 2 desk switching SHALL keep the same active source-rule context, grouped Tag handoff cues, and blocker diagnostics so the operator can continue toward Tag review and Output review without reconstructing planning state.

#### Scenario: Desk switching keeps Source-to-Tag continuity visible
- **WHEN** an operator changes Step 2 from `Inspect` to `Build` or `Triage`
- **THEN** the UI keeps the active rule summary and grouped Tag handoff cues visible for the same rule
- **AND** the operator can still understand what Step 3 will review next without reopening or reloading the workspace

#### Scenario: Desk switching does not behave like a separate workflow branch
- **WHEN** an operator uses the Step 2 desk selector repeatedly during one planning session
- **THEN** the workflow remains a single `Device -> SourceRule -> Tag review -> Output review/apply` path inside `/studio`
- **AND** changing desks SHALL NOT clear progress, hide the primary next-step call to action, or create a second primary workflow inside Step 2
