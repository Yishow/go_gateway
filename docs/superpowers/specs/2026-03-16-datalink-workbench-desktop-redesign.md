# Datalink Workbench Desktop Redesign

- Status: Approved design draft
- Date: 2026-03-16
- Scope: `frontend/src/pages/datalink/workbench/*`, related datalink hooks/services/types, desktop information architecture, and output-step UX alignment
- References:
  - `docs/superpowers/specs/2026-03-15-datalink-workbench-design.md`
  - `docs/superpowers/specs/2026-03-15-datalink-workbench-ui-detail.md`
  - `docs/superpowers/specs/phase2-runtime-dbtarget-detail.md`

## 1. Why this redesign exists

The first workbench landing proved the route direction was correct, but the desktop IA drifted away from the approved UI spec. On a 1920×1080 desktop, the current shell constrains the main area too aggressively, Step 1 embeds its own inspector inside an already narrow main column, Step 2 lacks a true planning/detail surface, and Step 3/4 do not yet communicate enough operational detail to feel production-ready.

This redesign restores the original workbench intent: one continuous route for `device -> source visualization -> tag binding -> local Modbus / database output`, with a full-width desktop shell, per-step inspector, and dense but readable operational detail.

## 2. Approved design constraints

### 2.1 Product path

The page must keep the user's single-path goal intact:

1. Configure a data source.
2. Visualize source addresses, current values, and transformation context.
3. Bind source points to tags.
4. Bind tags to local Modbus output or database output.

### 2.2 Hard constraints

- No legacy presentational component is reused.
- No old page shell is reused.
- Original datalink functional coverage must remain available.
- The redesign must not add speculative workflow branches just to make the UI look richer.
- Desktop layout must be optimized for 1920×1080 first, then degrade responsively.
- Future UI generation and implementation prompts for this redesign should explicitly target Opus-first output quality.

### 2.3 Reuse boundary

Healthy datalink domain foundations may still be reused:

- React Query hooks under `frontend/src/hooks/datalink/`
- API clients in `frontend/src/services/datalink.ts`
- shared datalink types
- pure feature helpers and validation logic

Legacy screens, cards, forms, drawers, grid shells, and panel composition are reference material only.

## 3. Desktop shell and information architecture

The desktop workbench returns to the approved five-region shell:

1. `StepRail`
2. `ContextBar`
3. `PrimaryWorkArea`
4. `InspectorPanel`
5. `BottomSummaryBar`

### 3.1 Shell rules

- The workbench should use the available app-shell width rather than a centered `max-w-7xl` content box.
- `ContextBar` stays sticky at the top of the workbench region.
- `BottomSummaryBar` stays sticky at the bottom and owns readiness/count messaging.
- `InspectorPanel` is the only persistent detail surface on desktop.
- The current right-side `ActionDock` is removed as a page-level shell region.

### 3.2 Why the current shell drifted

The earlier implementation used an interim summary dock as an always-visible page region. That made Step 1 and Step 2 compete for width inside the main column, which directly caused the 1920×1080 layout breakdown. This redesign corrects that drift by restoring the spec-approved shell instead of layering more content into the constrained main area.

## 4. Step 1 — `DeviceWorkspace`

Step 1 is not just device selection. It defines the source capability context that Step 2 must interpret correctly.

### 4.1 Primary work area

The main region is a new `DeviceBrowser` surface:

- search
- protocol filter
- status filter
- `Create`
- `Refresh`
- `Clone from selected`
- device cards / compact rows with a quick capability summary

Each device card should surface the fields that materially affect downstream source planning:

- `address base`
- `word order`
- `unit id`
- protocol badge / protocol-specific capability hints

A secondary `Recent / Pinned Devices` strip may exist beneath the browser to speed up repeated operator switching.

### 4.2 Inspector panel

The right-side inspector shows:

- device identity and connection parameters
- capability summary (`address base`, `word order`, `unit id`, protocol traits)
- recent three connection tests in a compact timeline
- primary actions: `Edit`, `Test connection`, `Clone`

### 4.3 Clone flow

Cloning uses a new workbench-specific create drawer. It pre-fills reusable connection fields from the selected device, such as host, port, unit id, protocol, and capability defaults, but still requires a new device identity and keeps the new device creation flow visually separate from legacy editors.

### 4.4 Context bar behavior

Once a device is selected, the context bar must show:

- device name
- protocol
- connection state
- last test result
- capability chips for `address base` and `word order`

## 5. Step 2 — `AddressCanvasWorkspace`

Step 2 becomes the primary planning and live-validation workspace.

### 5.1 Core spatial model

The address canvas uses one stable model:

- A **rule** is one source-planning directive with:
  - start address
  - data type / bit width
  - item count
  - optional naming or template metadata
- Each rule expands into one or more merged blocks on the shared 16-bit lattice.
- base unit = `16-bit` cell
- `32-bit` / `float32` = merge 2 adjacent 16-bit cells
- `64-bit` = merge 4 adjacent 16-bit cells
- multiple rules may coexist at once
- the rendered canvas spans from the minimum rule start address to the maximum rule end address
- gaps remain visible as empty address space rather than disappearing

This means the operator can understand both occupied ranges and unused holes across the same continuous memory segment.

### 5.2 Primary work area structure

#### `PlannerToolbar`

The toolbar owns the planning and interpretation controls:

- view switch: `Plan / Live / Link`
- rule create / apply controls
- value format switch: `decimal / hex / binary / float`
- `freeze live`
- `snapshot compare`
- `jump to address`
- `save as template`
- `load template`

#### `RuleLayerBar`

Multiple rules must be visible and manageable, not hidden in a modal:

- enable / disable
- lock / unlock
- reorder
- focus a single rule
- show start address, bit width, count, and covered range

#### `AddressCanvas`

The canvas is the central spatial surface. It should remain geometrically stable while the information layer changes by view mode.

#### `CoverageOverview`

A compact overview strip shows the entire address coverage at a glance:

- planned ranges
- gaps
- conflicts
- quick navigation to a chosen sub-range

### 5.3 Three view modes

#### `Plan`

Planning view emphasizes:

- address readability
- merged spans
- rule boundaries
- coverage state
- gaps
- conflicts

#### `Live`

Live view keeps the same geometry but changes the information density:

- raw value
- transformed value
- value quality / timestamp when available
- value-format-dependent rendering
- freeze / snapshot comparison while planning

#### `Link`

Link view is the validation lens into downstream work:

- point status
- tag binding status
- output readiness state

It exists to help transitions into Step 3 and Step 4, not to become a fourth workflow.

### 5.4 Template / recipe support

Rule groups can be saved as reusable templates.

For this redesign round, template persistence is defined as:

- browser-local persisted templates only
- backed by workbench-specific local storage
- no new backend template API in this round

Each template contains:

- rule list
- rule ordering / lock state
- preferred Step 2 view mode
- optional operator-facing template name

Each template explicitly excludes:

- live values
- temporary freeze snapshots
- device-specific runtime results

Applying a template to another device must warn when capability assumptions no longer match, especially around `address base` and `word order`.

### 5.5 Inspector panel

The inspector changes by selection type:

- selected rule -> coverage, bit width, count, template origin, quick actions
- selected span/block -> address, bit width, raw value, transformed value, point state
- multi-selection -> batch point creation from selected address spans and batch apply summary

### 5.6 Auxiliary audit surface

Dense list/table auditing remains available as a secondary companion surface rather than a competing fourth primary mode.

The audit surface is opened from a Step 2 toolbar toggle and appears as a secondary drawer/panel, not as another persistent desktop shell column.

It covers dense review tasks such as:

- address sorting and filtering
- status-based audit
- batch point selection confirmation
- quick comparison between adjacent spans and their derived point state

## 6. Step 3 — `TagBindingBoard`

Step 3 becomes an information-dense source-to-tag review and creation workspace.

### 6.1 Main board

Each row/card should show enough detail to review a binding decision without opening the inspector first:

- tag key naming-rule preview
- source address / span snapshot
- raw value vs transformed value
- data type / bit width / merge rule
- status layer: `bound / unbound / partially complete`

### 6.2 Required management actions

The board must support:

- create new tag
- bind existing tag
- switch between the two flows explicitly
- batch naming
- batch field application
- filter/search by status, rule, or keyword
- preview diffs before applying a batch change
- review a post-apply result summary (`created / linked / skipped / failed`)

### 6.3 Inspector panel

The inspector owns deeper editing rather than repeating the full board:

- full tag metadata edit
- conflict explanation
- naming-template breakdown
- single-item override actions
- batch summary detail when multiple items are selected

## 7. Step 4 — `OutputWorkspace`

Step 4 keeps Local Modbus and Database output inside one output workspace.

### 7.1 Shared structure

The primary area contains:

- `OutputTargetSwitcher`
- `OutputCandidateBoard`
- active target studio (`Local Modbus` or `Database`)

Each output candidate row should show:

- source address / span
- tag key
- latest value / state
- local Modbus mapping status
- database mapping status
- overall output readiness state

The workspace must also support filtering by `mapped / unmapped / conflict`.

Local Modbus and Database outputs are not mutually exclusive. A single candidate may be ready for one target, both targets, or neither target at the same time.

### 7.2 Local Modbus studio

Local Modbus needs to feel spatial and operational, not just tabular.

Required surfaces:

- `RegisterMapCanvas` for register allocation visualization
- base / offset visualization
- auto-map strategy selection:
  - sequential fill: assign from the current base upward without skipping free space
  - gap-aware fill: preserve existing occupied ranges and fit new mappings into remaining gaps first
  - data-type-aligned fill: prefer start positions aligned to the candidate data width
- conflict highlighting and conflict repair suggestions
- server health summary
- test write / dry-run validation
- final sync / failure status visibility

### 7.3 Database studio

Database output must make schema alignment visible.

Required surfaces:

- connector health summary
- schema snapshot with type badges
- required/missing field highlighting
- mapping UI for target columns
- write-row preview before commit
- final sync / failure status visibility

### 7.4 Inspector panel

The inspector shows the full source -> tag -> output chain for the selected candidate, including reasons for `partial` or `blocked` readiness.

## 8. Cross-step behavior

### 8.1 Selection and context retention

- The selected device remains visible in `ContextBar` across all steps.
- Source context must remain understandable when the user moves from Step 2 to Step 3.
- Output candidates must preserve traceability back to source spans and tag definitions.

### 8.2 Readiness model

Every main object shown in the shell should map to explicit readiness states:

- `draft`
- `ready`
- `partial`
- `blocked`
- `applied`

`WorkbenchReadiness` is a UI-derived cross-step state, not a replacement for existing backend/domain enums.

Minimum mapping rules:

| Domain object | Existing status | Workbench readiness |
| --- | --- | --- |
| Device | `ready` | `ready` |
| Device | `not_ready` | `draft` |
| Device | `error` | `blocked` |
| Device | `disabled` | `blocked` |
| Tag | `draft` | `draft` |
| Tag | `active` with missing output mapping | `partial` |
| Tag | `active` with required mappings present | `ready` |
| Tag | `retired` | `blocked` |

`applied` is reserved for UI flows that have already completed a target-side apply/sync action, such as confirmed output mappings or a completed batch apply result.

### 8.3 Bottom summary bar

The bottom bar owns the compact global summary:

- selected step state
- points configured
- tags bound
- Local Modbus mappings
- Database mappings
- validation / output readiness
- latest actionable issue count

### 8.4 Loading and error states

Each workspace must explicitly handle:

- initial loading
- partial data load
- API error
- connection-lost / health-degraded state

Errors must remain inline and actionable rather than falling back to generic empty states.

### 8.5 Localization and accessibility baseline

- All user-facing labels, status text, empty-state copy, and validation messaging must remain compatible with the existing i18n contract.
- Keyboard navigation must cover step switching, primary toolbar actions, list selection, and inspector actions.
- The redesign must preserve the project's accessibility direction for focus order, live regions, and screen-reader-readable state changes.

## 9. New component boundaries

The redesign should be implemented as new workbench-first units:

- `WorkbenchFrame`
- `WorkbenchContextBar`
- `WorkbenchBottomSummaryBar`
- `DeviceWorkspace`
- `DeviceInspectorPanel`
- `AddressCanvasWorkspace`
- `PlannerToolbar`
- `RuleLayerBar`
- `CoverageOverview`
- `TagBindingBoard`
- `TagBindingInspectorPanel`
- `OutputWorkspace`
- `RegisterMapCanvas`
- `DatabaseSchemaBoard`
- `OutputInspectorPanel`

These names are functional boundaries, not a requirement to mirror the old spec's exact file names.

## 10. Validation targets

The redesign is not considered ready until it proves all of the following:

- desktop shell remains legible on 1920×1080
- Step 1 no longer embeds a second desktop inspector inside the constrained main column
- Step 2 supports multi-rule planning with merged cells and visible gaps
- Step 2 can switch between `Plan`, `Live`, and `Link` without changing geometry
- Step 3 shows source-to-tag detail densely enough for batch review
- Step 4 presents both Local Modbus and Database as one continuous output workspace
- no legacy UI components are reused in the new workbench surface

## 11. OpenSpec sync intent

This redesign should be mirrored into a new active OpenSpec change so the delta can be translated into proposal, design, spec, and tasks artifacts without relying only on the informal session context.
