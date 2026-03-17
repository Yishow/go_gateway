# Datalink Workbench Master-Detail Polish

- Status: Approved design draft
- Date: 2026-03-17
- Scope: `frontend/src/pages/datalink/workbench/*`, related workbench state/tests, and desktop-first operator flow for `/datalink/workbench`
- References:
  - `docs/superpowers/specs/2026-03-15-datalink-workbench-design.md`
  - `docs/superpowers/specs/2026-03-16-datalink-workbench-desktop-redesign.md`
  - `docs/superpowers/specs/2026-03-16-datalink-workbench-desktop-polish-design.md`
  - `docs/superpowers/specs/phase2-runtime-dbtarget-detail.md`

## 1. Why this polish round exists

The current workbench already restored the main datalink route, but the operator experience is still not stable enough for the user's real job on a `1920×1080` desktop.

The remaining problems are not primarily missing APIs. They are interaction and hierarchy problems:

- Step 1 still spends too much space on search/form chrome instead of device detail.
- Step 2 still makes planning feel like a tool console instead of a clear source-mapping studio.
- Point creation is still too ambiguous, so Step 3 can end up blocked by an unclear Step 2 state.
- The UI still makes the operator think too much about controls and not enough about the current source-to-output task.

This round is therefore a **master-detail polish**, not a workflow rewrite. The goal is to make the existing workbench feel like one coherent industrial studio:

1. choose/configure a source device
2. visualize source locations and transformed values
3. create points and bind them to tags
4. bind tags to local Modbus or database output

## 2. Operator profile and design intent

This workbench is a B2B operational tool for a single operator who needs high scan speed, predictable action placement, and low ambiguity.

The visual direction is intentionally restrained:

- neutral, low-noise desktop tone
- high contrast reserved for action, conflict, and readiness
- no extra decorative UI layer in this round
- desktop-first layout with responsive degradation instead of mobile-first card stacking

The intent is not to make the page look more impressive. The intent is to let the operator answer three questions quickly:

1. What am I working on now?
2. What can I do next?
3. What is blocking me?

## 3. UX psychology rules applied

This spec follows the workbench's approved direction through three key rules:

### 3.1 Hick's Law — reduce simultaneous choices

Each step should present one dominant work surface and one dominant next action.

That means:

- Step 1 must feel like browsing and selecting devices, not filling a giant form.
- Step 2 must show planning and conflict resolution as the main work, with utilities demoted.
- Step 3 must not pretend to be available before Step 2 has created at least one usable point.

### 3.2 Fitts' Law — keep primary actions close and obvious

The main CTA in each step should remain in a stable zone and stay large enough to hit quickly.

That means:

- `Create device`, `Create selected points`, `Create rule points`, and step-forward actions keep predictable placement.
- floating or hidden critical actions should be avoided unless there is a clearly visible equivalent in the main surface.

### 3.3 Miller's Law — chunk detail into scanable groups

Dense industrial data is acceptable. Unchunked dense data is not.

That means:

- health summary should collapse into a few meaningful counts
- conflict handling should become an actionable queue instead of a vague warning block
- inspector content should group detail into short operational sections instead of long mixed panels

## 4. Approved shell refinements

### 4.1 Shell posture

The existing workbench shell remains in place, but the visual posture becomes more obviously master-detail:

- left = selection / navigation
- center = primary artifact of the step
- right = explanation, trace, and side effects

The shell should not introduce another persistent summary or utility region in this round.

### 4.2 Context emphasis

The top context region should answer the current-step question in one compact sentence, while the center owns the real work.

It should not become a second dashboard again.

## 5. Step 1 — Device master-detail studio

Step 1 should become a true master-detail screen instead of a list with a large form burden.

### 5.1 Layout

- left master column: compact device list
- center detail panel: selected device overview and action surface
- right inspector: connection diagnostics, capability detail, and recent tests

The device itself must become the focus. Search and filtering should help selection, not dominate the page.

### 5.2 Master list rules

The device list should:

- keep `search`, `protocol filter`, and `Create device` in one aligned row
- use compact rows/cards with strong scanability
- show only the fields that matter for downstream work:
  - device name
  - protocol
  - endpoint or host summary
  - connection/readiness state
  - the most relevant capability hints

Secondary actions like `edit`, `clone`, and full diagnostics must not inflate every list row.

### 5.3 Detail surface

When a device is selected, the center detail area should immediately show:

- device identity
- protocol and endpoint summary
- source capability essentials
- last successful / failed connection check
- the next step CTA into source planning

When no device is selected, the detail area should explain the expected action rather than showing dead controls.

### 5.4 Create/edit interaction

The device editor should stop feeling like a cramped side utility.

Approved behavior:

- desktop: create/edit uses a large centered full-screen-style dialog surface
- the form opens with the first important fields visible without extra scrolling
- the shell behind it is visually quieted, not fully replaced
- smaller widths collapse to a stacked responsive layout without losing field grouping

This keeps creation and editing focused while still feeling like part of the workbench.

### 5.5 Responsive rules

For narrower viewports:

- master/detail may stack, but the selected device identity must remain pinned near the top
- filters should collapse before the device detail loses clarity
- the create/edit surface may switch from centered dialog posture to a full-height sheet

The desktop contract still leads this round.

## 6. Step 2 — Source planning studio

Step 2 is the heart of the workbench and must read as a planning canvas first.

### 6.1 Primary hierarchy

The visual order should be:

1. compact planning toolbar
2. large address canvas
3. lightweight rule layer and conflict queue
4. utilities behind a secondary disclosure

`More tools` must not compete with point creation or selection actions.

### 6.2 Right-top health summary

The summary is intentionally reduced to three chunks:

- `ready to create`
- `in conflict`
- `protected`

Next to those counts sit the two main actions:

- `Create selected points`
- `Create rule points`

This lets the operator understand both readiness and next action without scanning the whole canvas.

This summary is **Step 2-local readiness only**. The persistent bottom summary bar should continue to show cross-step progress and shell-level readiness, but it should not duplicate the same `ready/conflict/protected` counts again.

### 6.3 Mixed point-creation model

Step 2 should support two clear creation paths:

#### A. Rule-first creation

The operator defines a rule and creates points for the whole rule range.

Use when:

- the pattern is regular
- the rule is already correct
- the operator wants speed

#### B. Selection-first creation

The operator drag-selects one or more logical spans and creates points only for those spans.

Use when:

- only part of the rule should become points
- the operator needs local cleanup
- conflicts or exceptions need manual control

Right-click may exist as a shortcut, but it must never be the only way to reach point creation.

### 6.4 Selection toolbar

After drag-selection, a small local action strip should appear near the selection or in a consistent canvas action zone with:

- `Create points`
- `Add to rule`
- `Skip`

This provides a mouse-friendly fast path without forcing the user into hidden gestures.

### 6.5 Rule editing

Rules must support direct correction instead of forcing context switching.

Approved behavior:

- rule cards/rows support inline edit for start address, count, and data type
- the currently focused rule is visually distinct
- deleting or modifying a rule updates the canvas and conflict state immediately

### 6.6 Conflict handling

Conflicts should become an actionable queue, not just red warnings.

Each conflict row should clearly show:

- which span or rule is affected
- why it conflicts
- what the operator can do next

Supported resolution actions should include:

- move start address
- shorten range/count
- change data type
- skip the conflicting span

If a conflict blocks Step 3 readiness, that state must be explicit.

### 6.7 Logical cell model

The canvas continues to use `16-bit` as the base lattice unit, but larger widths must behave as one logical item:

- `32-bit` values render and select as one logical cell made from two base cells
- `64-bit` values render and select as one logical cell made from four base cells
- half-selection of a multi-word logical cell is not allowed

This rule is important because the user explicitly wants `32-bit` values to feel like one coherent item instead of split fragments.

### 6.8 Data type completeness

The data-type selector should be grouped by width and intent rather than exposing a random short list.

Expected groups:

- `single-word`: `bool`, `int16`, `uint16`, `string`
- `32-bit`: `int32`, `uint32`, `float32`
- `64-bit`: `int64`, `uint64`, `float64`

If some combinations are not yet supported by the current protocol/runtime or current source context, the option should remain visible but disabled with a reason. This round does **not** introduce a new `bitset` type unless the underlying type system is expanded in a separate change.

### 6.9 Rename lock semantics

The current lock semantics are too abstract for this workflow.

Approved wording:

- rename visible user-facing language toward `Protect plan` / `Protected`
- explain clearly that protected spans:
  - cannot be casually dragged
  - cannot be overwritten by a new rule
  - are excluded from accidental bulk creation

The behavior does not need to change radically if the current model already matches this meaning. The wording and explanation must.

## 7. Step 3 — Tag binding unlock and guidance

Step 3 should stop feeling broken when Step 2 has not yet produced usable points.

### 7.1 Unlock rule

Step 3 becomes available as soon as at least one point has been created and confirmed by the current point model. A rule or temporary selection by itself does not unlock Step 3.

In occupancy terms, `planned` or `selected` spans are not enough. The unlock threshold should align with at least one persisted/usable point reflected as `used` or later.

### 7.2 Empty state

If no points exist yet, the empty state should explicitly state:

- no points have been created yet
- how many eligible spans are currently available in Step 2
- a direct action to return to Step 2

This makes the blockage understandable instead of mysterious.

### 7.3 Candidate board

The current row-board direction remains valid, but Step 3 should visually reinforce that it is a review-and-bind surface:

- row density stays compact
- batch actions stay hidden until selection exists
- row focus opens richer detail in the inspector
- already-linked or conflict states remain visible but do not dominate every row

## 8. Step 4 — Output continuity

This round does not reinvent output mapping. It preserves the active-target approach already approved in prior work.

The Step 4 contract remains:

- one shared tag-selection surface
- one active output target at a time
- target-specific detail in the workspace and inspector

The only refinement required here is continuity:

- Step 3 and Step 4 should feel unlocked by successful point/tag progress
- output should not duplicate selection UI that already exists elsewhere

## 9. Error handling and operational feedback

Operational feedback should be direct and local to the operator's current task.

Rules:

- blocking issues appear near the relevant work surface
- rich diagnostic detail belongs in the inspector or conflict queue
- error language should tell the operator what to do next, not only what failed
- empty, blocked, and ready states must all look intentionally different

## 10. Testing strategy

This polish round should be implemented with regression-first discipline.

Required focus:

- Step 1 master-detail layout remains stable at `1920×1080`
- device search/filter/create stays in one aligned control row at desktop width
- create/edit device dialog opens in the approved large centered posture
- Step 2 keeps one dominant planning toolbar and a visibly dominant address canvas
- Step 2 summary shows `ready`, `conflict`, and `protected` counts
- point creation works through both rule-first and selection-first paths
- conflict queue actions update readiness and visible guidance
- multi-word logical cells cannot be half-selected
- Step 3 empty state explains why it is blocked and links back to Step 2
- Step 3 batch actions stay hidden until rows are selected
- Step 4 retains shared tag-selection continuity without duplicating selection controls

## 11. Implementation notes

This document is a targeted polish layer over the existing workbench, not a license to restart the redesign.

Therefore:

- keep current route and provider flow intact
- prefer surgical refactors inside existing workbench screens
- do not reintroduce legacy SmartDashboard or legacy Local Modbus page structure
- do not let utility controls reclaim primary visual hierarchy

## 12. Success criteria

This polish round is successful when:

- Step 1 reads as a device browser plus focused detail studio, not a crowded setup screen
- Step 2 makes point creation and conflict resolution obvious without hunting through secondary controls
- `32-bit` and larger values feel like single logical items in the canvas
- Step 3 no longer feels blocked without explanation
- the operator can move from source planning to tag binding to output mapping with less hesitation and fewer dead-end states
