# Change: Refactor Datalink to Single-Screen Pipeline Studio

## Why

Current Datalink operations are still split across Smart Dashboard, Points, Mappings, and Wizard pages. This causes context loss for operators when configuring source quantity/type, visualizing memory occupancy, updating tags, and committing to DB.

The new workflow must be fully intent-driven in one screen, with clear visual guidance and predictable state transitions.

## What Changes

- Replace fragmented UI flow with a single-screen Pipeline Studio workspace.
- Add source planning with persisted templates (save/reuse source quantity and type presets).
- Add typed memory occupancy rules (cell consumption derived from selected type and count).
- Add inline global tag editing directly from Memory Grid linkage workflow.
- Add guided motion system for step transitions and commit feedback.
- Deprecate legacy workflow pages (`/datalink/points`, `/datalink/mappings`, `/datalink/wizard`) from operator navigation.
- Keep full flow intent visible: Source Plan -> Grid Allocation -> Tag Linkage -> DB Commit.

## Impact

- Affected specs:
  - `datalink-ui`
  - `datalink-smart-dashboard`
  - `point-catalog`
  - `tag-dictionary`
- Affected code:
  - `frontend/src/pages/datalink/SmartDashboard.tsx`
  - `frontend/src/components/datalink/MemoryGrid.tsx`
  - `frontend/src/components/datalink/QuickActions.tsx`
  - `frontend/src/components/datalink/PointDetailPanel.tsx`
  - `frontend/src/components/datalink/TagForm.tsx`
  - `frontend/src/layouts/DatalinkLayout.tsx`
  - `frontend/src/App.tsx`
  - related hooks and API contracts for source templates and commit UX
