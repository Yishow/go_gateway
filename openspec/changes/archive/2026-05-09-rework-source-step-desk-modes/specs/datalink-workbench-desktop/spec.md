## ADDED Requirements

### Requirement: Source workspace keeps a persistent desk skeleton
The workbench SHALL reserve a persistent Step 2 desk skeleton for the active source rule before rendering desk-specific primary surfaces.

The persistent skeleton SHALL include the active rule summary, grouped handoff strip, desk selector, blocker diagnostics, and current primary action so the desktop workspace keeps a stable reading order across desk changes.

#### Scenario: Shared desk skeleton remains anchored across desk changes
- **WHEN** the operator switches the Source workspace between `Inspect`, `Build`, and `Triage`
- **THEN** the same summary and handoff band remains anchored in the desktop shell above the mode-specific workspace
- **AND** the operator does not lose the current rule, planning status, or blocker explanation

#### Scenario: Desk skeleton keeps the primary action in a stable location
- **WHEN** the operator edits, reviews, or triages the active source rule
- **THEN** the desktop workspace keeps the current primary action in the same Step 2 shell region
- **AND** the operator does not need to scan different panels to find the next planning action after each desk switch

### Requirement: Source desk layouts allocate distinct primary regions on desktop
The workbench SHALL allocate distinct desktop-first primary regions for `Inspect`, `Build`, and `Triage` instead of reusing one nearly identical layout.

`Inspect` SHALL prioritize the source canvas, `Build` SHALL prioritize the rule editor while keeping a readable secondary canvas region, and `Triage` SHALL prioritize the incident queue and affected recovery context.

#### Scenario: Build keeps a readable secondary canvas region
- **WHEN** the operator opens `Build` on a desktop breakpoint
- **THEN** the rule editor becomes the primary panel
- **AND** the remaining canvas region stays large enough to inspect current address coverage without collapsing into a token-height preview

#### Scenario: Triage elevates incident review above general canvas browsing
- **WHEN** the operator opens `Triage` with pending planning issues
- **THEN** the incident queue and recovery controls occupy the primary desktop region
- **AND** the canvas is presented only as supporting context for the affected ranges instead of as the dominant browsing surface
