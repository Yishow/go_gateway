# Datalink Workbench Desktop Polish — Scan-First

- Status: Approved design draft
- Date: 2026-03-16
- Scope: `frontend/src/pages/datalink/workbench/*`, related workbench tests, desktop information hierarchy, and entry continuity into `/datalink/workbench`
- References:
  - `docs/superpowers/specs/2026-03-15-datalink-workbench-design.md`
  - `docs/superpowers/specs/2026-03-16-datalink-workbench-desktop-redesign.md`
  - `docs/superpowers/specs/phase2-runtime-dbtarget-detail.md`

## 1. Why this polish round exists

The second-round workbench redesign restored the correct route and coverage, but the desktop still feels too busy for the user's real operating goal. On a 1920×1080 desktop, too many controls still compete at the same visual level: the context bar, summary bar, per-step toolbars, dense cards, and inspector content all ask for attention at once.

This polish round does not introduce a new workflow. It reduces noise inside the existing `device -> source -> tag -> output` path so the operator can scan the current step quickly, identify the next action without hunting, and keep secondary diagnostics nearby without letting them dominate the main work area.

## 2. Approved outcome

The approved direction is **scan-first desktop polish**.

This means:

- reduce visual noise before adding more surface detail
- keep one clear primary action per step
- let the center workspace own the main job of the step
- push secondary detail and diagnostics toward the right-side inspector or secondary controls
- preserve current capability coverage instead of rebuilding the workflow again

The approved improvement bundle for this round is:

- `01` Context bar keeps only step summary + primary CTA
- `02` Bottom summary becomes a lighter progress strip
- `05` Primary CTA stays in a consistent location
- `07` Device search/filter/create collapse into one row
- `11` Source step keeps one compact toolbar
- `13` Address canvas becomes the visual hero
- `14` Tag candidates become row-based boards instead of thick cards
- `15` Batch actions appear only after selection
- `17` Output candidates show only active-target information
- `20` Unavailable actions hide until the required selection exists

## 3. Shell-level information hierarchy

### 3.1 `WorkbenchContextBar`

The context bar should stop behaving like a second dashboard. It becomes a compact status sentence plus one primary action.

It keeps:

- current step label
- selected device / source / target summary
- the step's primary CTA

It de-emphasizes or moves out:

- long capability chip runs
- repeated readiness states already visible elsewhere
- recent test timelines
- secondary actions such as clone, template, schema, or audit shortcuts

The ideal top-line feeling is: one sentence tells the operator where they are, one button tells them what to do next.

### 3.2 `WorkbenchBottomSummaryBar`

The bottom bar becomes a lighter progress strip instead of another fully weighted control surface.

It should:

- keep counts and readiness available
- highlight the active step only
- reduce the visual weight of non-active steps to compact readiness markers
- keep the active output target visible without competing with the context bar

### 3.3 Global shell rule

Across every step, the user's eye path should be:

1. top summary
2. center work area
3. right-side detail

Not:

1. top badges
2. bottom badges
3. side badges
4. then maybe the actual workspace

## 4. Step-by-step polish

### 4.1 Step 1 — `WorkbenchDeviceStep`

Step 1 becomes easier to scan by lowering card height and flattening repeated controls.

Approved changes:

- search, protocol filter, and `Create device` live in one aligned row
- device items use a lower-height row or compact card treatment
- the list shows only identity, endpoint, connection health, and the two most important capability hints
- edit, clone, and test details move to the inspector or side panel instead of inflating every list item
- if no device is selected, the main surface emphasizes the selection task instead of showing many inactive follow-up controls

The goal is that the operator can scan the device list like a browser, not like a wall of mini dashboards.

### 4.2 Step 2 — `SourceCanvasSection`

Step 2 must visually communicate that the address lattice is the main artifact.

Approved changes:

- keep one compact toolbar only
- move template, save/load, coverage, and similar utilities out of the center attention zone
- move `jump`, `freeze`, and `snapshot` into a lighter secondary control group
- increase the visual dominance of `AddressCanvas`
- visually demote the rule list so it supports the canvas instead of competing with it

The first glance should say: "this is where source memory is planned and checked", not "this is a dense control console."

### 4.3 Step 3 — `TagBindingStudio`

Step 3 should feel like a fast review board, not a stack of heavy cards.

Approved changes:

- convert candidate presentation toward row-board / table-like density
- keep the main row focused on source identity, tag state, and a small amount of status
- move conflict and already-linked detail to the inspector
- show batch actions only after one or more rows are selected

This reduces scrolling pressure and makes multi-item comparison much easier on desktop.

### 4.4 Step 4 — `LocalModbusBoard` / `DatabaseTargetBoard`

Step 4 should stop showing both output worlds with equal weight all the time.

Approved changes:

- the center board shows only the currently active target's necessary mapping information
- non-active-target detail moves to the inspector
- Modbus metrics / server controls become lighter secondary sections
- database schema snapshot / write preview become grouped secondary sections rather than always-on equal-priority panels

The operator should feel that they are working in **one active output mode at a time**.

## 5. Interaction rules

### 5.1 Primary action consistency

Every step keeps the primary CTA in a stable position. The label may change by context, but the location should not.

### 5.2 Progressive disclosure

Controls that cannot be used yet should not dominate the screen.

Examples:

- batch actions stay hidden until rows are selected
- mapping detail stays hidden until a candidate is focused
- inspector sections open based on actual selection context

### 5.3 Entry continuity

The existing workbench deep-link behavior remains intact.

When entering from SmartDashboard or legacy routes:

- keep `step` / `target` routing
- allow a short, dismissible handoff summary if needed
- avoid adding another persistent banner layer after the first-use context is understood

## 6. Error handling and operational feedback

This round favors inline and inspector-based feedback instead of more modal interruptions.

Rules:

- errors that block the next step may appear in the main surface
- diagnostic detail belongs in the inspector
- connection test history, schema load detail, and conflict traces remain available without taking over the primary work area
- existing a11y and live-region behavior remains intact

## 7. Test strategy

This polish round should be implemented by tightening regression coverage before style-heavy edits.

Required focus:

- desktop shell remains stable at `1920×1080`
- top toolbar density does not wrap excessively
- primary CTA remains visible and predictable
- Step 2 keeps the address canvas visually primary
- Step 3 hides batch actions until selection exists
- Step 4 center content reflects only the active output target
- inspector owns the moved detail rather than duplicating it back into the main surface

## 8. Implementation notes

This is a polish pass, not a route rewrite.

Therefore:

- do not replace the workbench workflow again
- do not reintroduce legacy page layouts
- prefer surgical refactors of `WorkbenchContextBar`, `WorkbenchBottomSummaryBar`, `WorkbenchDeviceStep`, `SourceCanvasSection`, `TagBindingStudio`, `LocalModbusBoard`, and `DatabaseTargetBoard`
- keep current routing, provider state, output-target deep links, and Phase 2 functionality intact

## 9. Success criteria

This polish round is successful when:

- the operator can identify the current step and primary next action immediately
- Step 2 visually centers the address lattice
- Step 3 becomes easier to compare across many candidates
- Step 4 no longer overloads the user with inactive-target detail
- desktop workbench feels calmer without losing datalink functionality
