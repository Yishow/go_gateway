## MODIFIED Requirements

### Requirement: Source planning surfaces grouped-tag handoff context

The workbench SHALL show, in Step 2, enough Source planning context for the operator to understand how planned points will move into grouped Tag review and downstream Database planning.

Step 2 SHALL treat the active source rule as the primary planning unit and SHALL present a rule-first summary before the larger canvas surface. The canvas MAY remain available for detailed address inspection, but it SHALL NOT be the only place where the operator can understand the current plan.

#### Scenario: Source summary explains grouped handoff
- **WHEN** an operator plans source points for a rule
- **THEN** Step 2 shows the planned point count and current naming-prefix context
- **AND** the handoff copy explains that the next step will review grouped tags before database planning

#### Scenario: Active rule summary precedes canvas reading
- **WHEN** an operator focuses a source rule in Step 2
- **THEN** the workbench shows that rule’s address coverage, rule-scoped planning status, and data-type intent before the operator reads the full lattice canvas
- **AND** the operator does not need to infer the current rule only from canvas cells or secondary tabs

## ADDED Requirements

### Requirement: Source step previews rule-scoped Tag review before navigation
The workbench SHALL provide, within Step 2, a rule-scoped preview of the next Tag review result for the currently focused source rule.

The preview SHALL reuse the same active-rule continuity that Step 3 consumes, so the operator can trust that the Source-step preview and Tag-step review are describing the same candidate scope.

#### Scenario: Clicked rule previews next Tag review scope
- **WHEN** an operator selects a source rule in Step 2
- **THEN** the workbench shows a preview of the Tag candidates that the next step will review for that rule
- **AND** the preview remains bound to the same focused rule when the operator navigates to Step 3

#### Scenario: Hovered rule temporarily previews downstream scope
- **WHEN** an operator hovers a different source rule in Step 2
- **THEN** the workbench MAY temporarily preview that rule’s Tag-review scope
- **AND** leaving the hover state restores the previously selected rule preview unless the operator explicitly changes focus

### Requirement: Source step prioritizes planning actions over mode switching
The workbench SHALL make applying the current rule to the planning surface the primary action in Step 2, while keeping mode switching as a secondary control.

#### Scenario: Apply planning is the primary action
- **WHEN** an operator is editing or reviewing a source rule in Step 2
- **THEN** the most prominent action communicates applying the rule to the planning surface
- **AND** the workbench gives immediate confirmation that the rule has taken shape for the next review step

#### Scenario: Mode tabs remain secondary
- **WHEN** Step 2 renders `plan`, `live`, or `link` modes
- **THEN** those controls remain available without overpowering the active-rule summary or primary apply action
- **AND** the operator can still identify the current plan without switching modes first
