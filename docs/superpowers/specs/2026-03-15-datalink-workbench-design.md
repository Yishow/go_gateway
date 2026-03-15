# Datalink Workbench Design

- Status: Approved for specification review
- Date: 2026-03-15
- Scope: `frontend/src/pages/datalink/*`, related datalink hooks/services/types, and required backend contract additions for runtime/live value and database target support

## 1. Summary

This design replaces the current fragmented datalink UI workflow with a single, production-grade workbench focused on one continuous user journey:

1. Configure the data source.
2. Visualize source addresses and transformation context.
3. Bind points to tags.
4. Bind tags to local Modbus output and, in the next stage of the same workflow, to database targets.

The chosen delivery strategy is a **hybrid transition**:

- Introduce a new route: `/datalink/workbench`
- Keep the existing `SmartDashboard` and `LocalModbusWorkbenchPage` as fallback during rollout
- Deliver the first production release with a complete `Source -> Tag -> Local Modbus` workflow
- Add runtime/live value feedback and database target support as the next release once the backend contracts are available

## 2. Problem Statement

The current datalink UI does not match the actual user goal.

### Current issues

1. **The main workflow is fragmented.**
   - `SmartDashboard` spreads work across planning, tag, modbus, and commit tabs.
   - `LocalModbusWorkbenchPage` is a separate page.
   - Device management still relies on modal-heavy flows.

2. **`SmartDashboardPage.tsx` is a page-level bottleneck.**
   - It acts as a page shell, workflow coordinator, modal manager, data mapper, and action router in one file.
   - The current structure makes safe iteration expensive.

3. **The grid is not treated as the primary workflow surface.**
   - The address grid exists, but the user loses context when switching tabs, overlays, and pages.

4. **Transformation and current value feedback are not visible enough.**
   - Users need stronger feedback for “what address am I reading”, “what value do I have now”, and “what does the transform produce”.

5. **Local Modbus output is operationally useful but visually disconnected.**
   - The domain logic exists and is usable.
   - The workflow placement is wrong.

6. **Database target output is not yet modeled as a first-class UI destination.**
   - It requires explicit contract design rather than ad-hoc page growth.

## 3. Goals

### Product goals

- Provide one clean, modern, production-ready page for the datalink workflow.
- Keep the grid as the primary visual surface.
- Make the workflow understandable to first-time users without relying on modal-driven exploration.
- Preserve operational depth: batch actions, conflict handling, validation, sync, import/export, and status visibility.

### UX goals

- Grid-first interface, table-assisted editing.
- Minimal context switching.
- Strong inline feedback for success, warning, blocked, and error states.
- A polished, modern industrial control UI that feels purposeful rather than experimental.

### Engineering goals

- Reuse healthy datalink foundations: hooks, services, types, `MemoryGrid`, and shared styles.
- Avoid a risky deep rewrite of `SmartDashboardPage.tsx`.
- Define clear units so the implementation plan can be executed incrementally.

## 4. Non-Goals

- Rewriting existing datalink domain types or React Query data hooks without need.
- Deleting the legacy `SmartDashboard` during the first release.
- Implementing speculative database target behavior before backend contracts exist.
- Introducing new workflow concepts unrelated to source, visualization, tag binding, or output binding.

## 5. Chosen Approach

Three approaches were evaluated:

1. Incrementally restructure `SmartDashboard`
2. Build a completely new standalone workbench immediately
3. **Hybrid transition**: build the new workbench while keeping legacy entry points during rollout

The chosen approach is **Hybrid transition** because it:

- Matches the approved user direction
- Avoids high-risk surgery on `SmartDashboardPage.tsx`
- Reuses most of the healthy domain and UI foundations
- Produces a production-quality first release faster than a full page replacement strategy

## 6. Information Architecture

The new page is `DatalinkWorkbenchPage` mounted at `/datalink/workbench`.

### Layout

The page is composed of five stable regions:

1. **Step Rail**
   - A persistent left-side navigation rail
   - Steps:
     - Device
     - Source Grid
     - Tag Binding
     - Output Targets

2. **Context Bar**
   - Persistent header within the page
   - Shows:
     - Selected device
     - Protocol
     - Connection status
     - Last connection test result
     - Quick actions: test connection, refresh, switch device

3. **Primary Work Area**
   - Keeps the source grid visible as the central visual anchor
   - Supports a grid/table toggle without leaving the page
   - Never relies on full-screen modals for the main path

4. **Inspector / Action Surface**
   - Shows the currently selected address span or point selection
   - Displays value, transform context, validation status, conflicts, and next actions
   - Adapts by step rather than spawning separate tabs

5. **Bottom Summary Bar**
   - Persistent status row
   - Shows counts and readiness:
     - configured points
     - bound tags
     - local Modbus mappings
     - database mappings
     - validation readiness

## 7. Visual Direction

The approved UI direction is:

- Modern and clean
- Practical first, but visually polished
- Grid-first, with reduced clutter
- Dark control-plane theme with clear state colors

### Visual rules

- Use a dark graphite/slate base to make states and values legible.
- Use color intentionally:
  - Blue: in progress / interactive
  - Green: valid / healthy / applied
  - Amber: warning / conflict
  - Red: blocked / failed
- Avoid stacked modal shells and scattered icon bars.
- Use larger, calmer layout regions with clear spacing and hierarchy.
- Keep the grid visually dominant over side controls.

## 8. Workflow Design

### Step 1: Device

This step replaces device-selection-driven modal workflows.

#### User can:

- Search devices
- Filter by protocol / status
- Switch device
- Test connection
- Create or edit a device through an embedded drawer/panel flow

#### UI structure

- Device selector cards or compact list
- Inline create/edit drawer
- Context bar updates immediately on selection

### Step 2: Source Grid

This is the primary working step.

#### User can:

- Define start address, data type, and count
- Apply a source plan to the grid
- Select one or multiple address spans
- Batch create points
- Toggle between:
  - Grid view
  - Table view

#### Inspector shows

- Selected address span
- Address metadata
- Point shape and data type
- Conflict state
- Current value or last poll value when available
- Transform preview status

#### Grid and table mode distinction

- **Grid mode** is the primary operational view for spatial understanding, contiguous span selection, quick conflict spotting, and multi-address planning.
- **Table mode** is the auxiliary auditing view for dense review and batch editing, showing rows with address, data type, span length, current value, transform state, tag binding status, and output readiness.

### Step 3: Tag Binding

This step is designed for batch usability rather than one-cell-at-a-time interaction.

#### User can:

- Select multiple points
- Batch generate tag keys
- Edit display names, units, descriptions
- Link to existing tags
- See affected mappings and conflict warnings before commit

#### Required behavior

- Preserve source context while editing tags
- Show selected points summary
- Highlight unbound, partially bound, and fully bound states

### Step 4: Output Targets

This step brings all output work into the same page.

#### Local Modbus

First production release must fully support:

- Server start/stop
- Port configuration
- Tag to register mapping
- Register conflict detection
- Mapping sync
- Test write
- Import/export

#### Database target

Database output lives in the same step but is delivered after its backend contract is implemented.

The final target workflow includes:

- Connector selection
- Target table selection
- Column mapping
- Missing-field detection
- Validation status

Until backend support lands, the database panel is present in the IA but marked as unavailable with precise reason text rather than hidden.

## 9. Component Boundaries

The new page must be composed of focused units with clear responsibilities.

### New page-level units

- `DatalinkWorkbenchPage`
  - Route entry
  - Page shell

- `WorkbenchProvider`
  - Owns local workflow state
  - Coordinates selected device, selected spans, selection summaries, and step transitions

- `WorkbenchStepRail`
  - Renders the persistent step navigation

- `WorkbenchContextBar`
  - Displays active device and quick actions

- `SourceCanvasSection`
  - Hosts `MemoryGrid`
  - Hosts grid/table mode switch
  - Hosts source planner controls

- `SourceInspectorPanel`
  - Displays selection details, transform preview, and conflicts

- `BatchTagBinder`
  - Handles multi-point tag creation and linking

- `OutputTargetsPanel`
  - Contains:
    - `LocalModbusTargetCard`
    - `DatabaseTargetCard`

- `WorkbenchBottomBar`
  - Shows readiness and validation counts

### Existing modules to reuse

- `MemoryGrid`
- `hooks/datalink/*`
- `services/datalink.ts`
- `types/datalink.ts`
- `frontend/src/styles/designSystem.ts`
- Pure feature helpers under `frontend/src/features/datalink/`

### Existing modules to extract logic from

- `LocalModbusWorkbenchPage`
  - Extract domain interactions into reusable hooks/sections

- `useSmartDashboardTagLinking`
  - Reuse domain behavior while removing SmartDashboard-specific assumptions

- `useSmartDashboardCommitFlow`
  - Reuse only if the validation model still fits the workbench

## 10. State Model

The workbench uses one dedicated workflow store/provider.

### Required frontend state

- Release 1:
  - current step
  - selected device id
  - selected address spans
  - grid view mode (`grid` or `table`)
  - planned source settings
  - selected points
  - tag binding draft state
  - local Modbus output draft state
  - validation summaries
- Release 2:
  - database target draft state
  - runtime/live value subscription state

### State requirements

- Preserve drafts while moving between steps
- Restore workbench state on route-local navigation
- Avoid mixing route query params with core workflow state
- Keep legacy deep links separate from the new workbench state model

## 11. Data Flow

### Step flow

1. Device is selected
2. Source grid plan is applied
3. Points are created or updated
4. Points are linked to tags
5. Tags are mapped to output targets
6. Validation is calculated and shown inline

### Important rule

The user never has to leave the workbench route to continue the main path.

## 12. Error Handling and Validation

Every step must expose one of these states:

- `draft`
- `valid`
- `warning`
- `blocked`
- `applied`

### Required error cases

- device connection failed
- invalid address range
- incompatible data type span
- overlapping grid allocation
- duplicate tag key
- tag already linked with conflicting mapping
- duplicate Modbus register
- unavailable runtime/live value feed
- unavailable database target contract

### UX requirements

- Errors must be inline and specific.
- Blocked states must explain the next fix action.
- Batch operations must return both summary and per-item detail.
- No silent failure paths.

## 13. Testing Requirements

The design is not considered ready for implementation without a full test plan.

### Unit tests

- workflow store/provider
- pure helpers and derived selectors
- output conflict helpers
- transform preview helpers

### Integration tests

- Device -> Source Grid -> Tag Binding -> Local Modbus main flow
- grid/table mode transitions
- step-to-step state preservation
- conflict and blocked state handling

### UI regression tests

- grid cell states
- inspector rendering
- output target cards
- summary bar readiness states

### End-to-end tests

For the production rollout, add E2E coverage for:

- creating a device
- planning points
- binding tags
- mapping to Local Modbus
- validating visible state transitions

## 14. Backend Contract Additions

The first production release does not require all backend additions, but the full design does.

### Runtime / live value

Add backend contracts for:

- `GET /api/v1/datalink/runtime/status`
- `GET /api/v1/datalink/runtime/stream`

These support:

- collector runtime status
- live value feedback
- per-device/per-point health indicators

Suggested response skeletons:

```ts
type RuntimeStatusResponse = {
  running: boolean;
  collectors: Array<{
    deviceId: string;
    deviceName: string;
    status: 'idle' | 'running' | 'warning' | 'error';
    lastReadAt: string | null;
    lastError: string | null;
  }>;
};

type RuntimeStreamEvent = {
  deviceId: string;
  pointId: string;
  address: string;
  value: string | number | boolean | null;
  transformedValue?: string | number | boolean | null;
  quality: 'good' | 'stale' | 'bad';
  timestamp: string;
};
```

### Point polling

`POST /api/v1/datalink/points/:id/poll` must return reliable current value data usable by the inspector and table mode.

### Database target contracts

Add backend contracts for:

- listing available output connectors
- listing available target tables
- creating/updating tag-to-column mappings
- validating target completeness

Suggested response skeletons:

```ts
type DatabaseConnectorSummary = {
  id: string;
  name: string;
  kind: 'sqlite' | 'postgres' | 'mysql' | 'sqlserver';
  status: 'ready' | 'warning' | 'error';
};

type DatabaseTargetMapping = {
  tagId: string;
  connectorId: string;
  table: string;
  column: string;
  writeMode: 'insert' | 'upsert';
};
```

The frontend must not ship fake database output flows. If the contract is unavailable, the database panel remains visibly unavailable with explicit reason text.

## 15. Rollout Plan

### Release 1: Production-ready main path

Deliver:

- `/datalink/workbench`
- Device selection
- Source planning and grid/table workflow
- Batch point creation
- Batch tag binding
- Local Modbus output panel with full operational behavior
- summary bar and inline validation
- legacy fallback link

This release is production-grade for the approved first scope.

### Release 2: Runtime visibility and database target

Deliver:

- live value visibility
- runtime status surface
- database target panel activation
- database validation states

### Release 3: Legacy reduction

Deliver:

- promote workbench to primary entry
- downgrade `SmartDashboard` to fallback or redirect
- retain rollback path until confidence is established

## 16. Migration Rules

- Do not delete the legacy page during Release 1.
- Do not route legacy query-intent paths into unfinished workbench flows.
- Add a clear “Try new workbench” or equivalent entry from the legacy page during rollout.
- Keep rollback simple: the legacy route remains operational until the new route is validated.

## 17. Final Decision

Proceed with a **hybrid transition** and treat the first workbench release as a production release for the approved scope, not as a throwaway MVP.

The final page is a clean, modern, grid-first workbench that keeps the user in one route and one mental model from source selection through output binding.
