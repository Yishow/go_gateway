# recording-memory-visual-refinement Specification

## Purpose

Define the visual refinement and accessible interaction contract for the existing five recording memory views. Preserve their shared source, occupancy and live-value authorities while making each layout serve its intended inspection task.

## Requirements

### Requirement: Memory views expose distinct inspection structures

The grid, navigator, map, chapters and detail views SHALL expose the same source information through their documented inspection structures in DESIGN.md. Range selection, current-window membership and live-data health MUST remain distinct and derived from actual evidence.

#### Scenario: A user chooses each inspection mode

- **WHEN** a user chooses grid, navigator, map, chapters or detail
- **THEN** grid prioritizes addresses, navigator provides saved-range navigation, map provides a complete current-window overview, chapters exposes expandable range summaries, and detail exposes selected-address facts
- **AND** no layout adds fictional health, timestamps, counts or addresses from a reference image

#### Scenario: Dense map or long chapter list

- **WHEN** the current map window contains 256 slots or the saved range list exceeds the visible area
- **THEN** all window slots and all saved ranges remain reachable without decorative sampling or inaccessible clipped content
- **AND** address navigation continues to target the exact canonical slot


<!-- @trace
source: refine-memory-view-visuals
updated: 2026-11-07
code:
  - DESIGN.md
-->

---
### Requirement: Typed cells preserve physical address meaning and safety

Typed cell presentation SHALL retain exact datatype text, physical-address accessibility and shared multi-register value identity. Unknown, conflicting, stale and failed evidence SHALL take priority over decorative type treatment. Existing occupancy and protocol proofs MUST remain authoritative.

#### Scenario: A multi-register value crosses a boundary

- **WHEN** a confirmed two- or four-register value crosses a rendered row or window boundary
- **THEN** the value remains identifiable by its start and continuation while each physical address keeps keyboard access
- **AND** the UI does not invent width for an unconfirmed string or unsupported bit area

#### Scenario: A cell action is requested

- **WHEN** the user acts on an empty, used, conflicting or unknown cell
- **THEN** the existing exact-address add, owner open, owner chooser or fail-closed action remains respectively in force
- **AND** live equipment readings are not presented as SQL receipt, readback or cleanup evidence


<!-- @trace
source: refine-memory-view-visuals
updated: 2026-11-07
code:
  - DESIGN.md
-->

---
### Requirement: Layout changes retain one shared interaction context

The five views SHALL preserve area, current window, keyboard focus and owner selection across switches using one shared parent live-value owner. Rendered DOM observers and keyboard navigation SHALL follow the actual active layout.

#### Scenario: Repeated layout changes and pointer inspection

- **WHEN** the user repeatedly switches layouts or inspects a value with pointer hover
- **THEN** selection and page context remain stable, observers track the newly rendered nodes, and subscriptions are not duplicated
- **AND** pointer inspection does not move the keyboard tab stop

#### Scenario: Small viewport and dismissal

- **WHEN** the modal is used at a narrow or short viewport and the user uses keyboard navigation, Close, Cancel, Escape or Back/Forward
- **THEN** essential controls, full datatype and address information remain reachable
- **AND** focus restoration and existing dirty or unresolved-operation barriers remain intact


<!-- @trace
source: refine-memory-view-visuals
updated: 2026-11-07
code:
  - DESIGN.md
-->

---
### Requirement: Visual acceptance uses current original-reference evidence

Visual refinement SHALL use DESIGN.md as the sole token and component-recipe authority and compare current product screenshots with each of the five original references. Intentional palette differences SHALL be documented without claiming pixel-identical output. Browser acceptance MUST be separate from DOM tests and historical evidence.

#### Scenario: Current visual and interaction evidence is reviewed

- **WHEN** the refinement is submitted for acceptance
- **THEN** each view has current light/dark desktop, tablet and phone evidence with product SHA, viewport, fixture and reference identity
- **AND** keyboard, touch, selection, repeated switching and dismissal checks have their own current results

#### Scenario: Browser execution is blocked

- **WHEN** the permitted final browser environment cannot start because of an environment restriction
- **THEN** browser and screenshot acceptance remain blocked and product tasks remain incomplete
- **AND** unit tests, generated images and historical Mac evidence are not substituted as a pass

<!-- @trace
source: refine-memory-view-visuals
updated: 2026-11-07
code:
  - DESIGN.md
-->