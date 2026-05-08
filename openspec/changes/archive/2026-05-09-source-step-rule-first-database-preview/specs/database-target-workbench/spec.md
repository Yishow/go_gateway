## ADDED Requirements

### Requirement: Database target context informs Source-step planning hints
The workbench SHALL let the active `database` output target influence Source-step planning guidance before the operator enters the Output step.

This guidance SHALL remain limited to planning defaults and grouped-row hints. It SHALL NOT require the operator to configure connector, schema, table, or column setup in Step 2.

#### Scenario: Database context changes Source-step grouping guidance
- **WHEN** the active output target is `database` while the operator is planning a source rule in Step 2
- **THEN** the workbench shows grouping-oriented hints that help the operator understand how the rule is likely to map into future table, row, or column groupings
- **AND** the workbench does not expose connector or schema setup controls in Step 2

#### Scenario: Database-aware preview remains rule-scoped
- **WHEN** Step 2 previews downstream Tag or Database implications for the focused source rule
- **THEN** the preview remains scoped to the active rule only
- **AND** the operator does not have to infer grouped-row implications from unrelated rules or device-wide state
