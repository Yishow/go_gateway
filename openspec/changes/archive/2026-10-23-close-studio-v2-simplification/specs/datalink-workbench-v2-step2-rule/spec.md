## MODIFIED Requirements

### Requirement: Point grid with modifier-key interactions

The current rule's point grid SHALL display expanded addresses and point names with responsive layout. Plain click SHALL select or inspect one point; Shift+click SHALL select a contiguous range from lastClickedIdx; Cmd/Ctrl+click SHALL toggle membership in a per-rule point-ID selection set. None of these selection gestures SHALL change skipped_addresses. Collection changes SHALL require an explicit collection checkbox or the existing clearly scoped selected-point enable/skip action. Switching the selected rule SHALL reset selection and lastClickedIdx.

#### Scenario: Plain click selects without skipping
- **GIVEN** cells 40001-40008 are collected
- **WHEN** the operator clicks 40003
- **THEN** only 40003 is selected and skipped_addresses remains empty.

#### Scenario: Shift+click range
- **GIVEN** lastClickedIdx is 0 at 40001
- **WHEN** the operator Shift+clicks index 3 at 40004
- **THEN** 40001 through 40004 are selected and collection is unchanged.

#### Scenario: Cmd+click multi-select
- **WHEN** the operator Cmd/Ctrl+clicks 40002, 40004 and 40006
- **THEN** those point IDs toggle in the selection set and the selected count is displayed without changing collection.

#### Scenario: Explicit collection change
- **WHEN** the operator unchecks a point's collection checkbox or explicitly skips the selected points
- **THEN** only the requested collection scope changes and it can be explicitly enabled again.

#### Scenario: Switching rule resets selection
- **WHEN** the operator switches from rule A with selected points to rule B
- **THEN** selection is empty and lastClickedIdx is null without changing either rule's collection configuration.
