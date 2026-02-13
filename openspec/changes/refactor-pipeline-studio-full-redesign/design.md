## Context

Operators must configure data ingestion in sequence:
1) define source quantity/type, 2) allocate memory positions, 3) confirm/edit global tags, 4) commit to DB.

Current UX has partial support but still requires cross-page context switching. The redesign introduces a single intent-centered workspace where every action is traceable and reversible.

## Goals / Non-Goals

- Goals:
  - One-screen completion for Source -> Grid -> Tag -> DB.
  - Persistent intent visibility and explicit state transitions.
  - Memory grid occupancy that matches selected source type and count.
  - Inline global tag edits from grid workflow.
  - Motion guidance without harming accessibility.
- Non-Goals:
  - No protocol driver rewrite.
  - No timeseries storage engine replacement.
  - No backend migration unrelated to source/template/tag workflow.

## Decisions

- Decision: Introduce `Pipeline Studio` as primary operating surface.
  - Left rail: source templates + source plan builder.
  - Center: memory grid with occupancy and linkage overlays.
  - Right rail: tag editor + commit queue + validation and commit controls.

- Decision: Persist source templates.
  - Template stores protocol context, source type, source count, naming strategy, and optional default tag naming.
  - Templates are reusable and editable.
  - Template metadata includes version and last-used timestamp to prevent stale reuse.

- Decision: Typed occupancy policy in memory grid.
  - Cell span per selected data type:
    - bool/int16/uint16: 1 cell per source.
    - int32/uint32/float32: 2 adjacent cells per source.
    - int64/uint64/float64: 4 adjacent cells per source.
  - Example:
    - `int x5` => 5 cells.
    - `float32 x10` => 20 cells shown as 10 linked pairs.
  - Conflict severity is explicit:
    - hard conflict: blocks commit.
    - soft conflict: requires confirmation and shows overwrite warning.
  - Planner supports auto-allocation strategy to find the nearest valid contiguous span.
  - Multi-cell allocations show grouped border and in-group index markers to prevent pairing ambiguity.

- Decision: Global tag edit in-place.
  - Tag edits from grid context update global tag dictionary.
  - UI must show impact summary before saving (affected mappings count).
  - Tag save requires explicit second confirmation when existing mappings are impacted.
  - Tag edits require before/after diff preview at field level before final save.

- Decision: Motion-guided workflow.
  - Use short transitions for intent handoff:
    - Source planned -> grid highlight pulse.
    - Grid allocation confirmed -> tag panel reveal.
    - Commit success -> sink confirmation animation.
  - Respect reduced-motion preference.
  - Animation is restricted to state handoff and commit feedback only.
  - Animation duration must stay in 150-300ms.

- Decision: Recovery, audit, and rollout safety.
  - Planner changes use autosave draft and one-click restore.
  - Commit executes in chunks with retry support for partial failure.
  - Every successful commit emits audit records (operator, time, changed addresses, changed tags).
  - Each commit creates a batch trace id for diagnostics and replay.
  - Initial rollout prioritizes high-frequency types (`int16`, `int32`, `float32`), then expands.
  - Snapshot restore supports reverting to selected historical planning checkpoints.

- Decision: Validation and visibility model.
  - Validation is two-stage:
    - Stage 1 structural validation (schema/type/span/conflict checks).
    - Stage 2 executable validation (pipeline simulation + sink readiness checks).
  - Grid supports conflict-only filter for rapid remediation.
  - Commit panel shows pre-commit polling load delta estimate.

- Decision: Local Modbus sharing server.
  - Add optional sink target `local_modbus_server`.
  - Selected tag values are mirrored into local virtual Modbus memory grid.
  - Local Modbus server listens on TCP port `5020` for external clients.
  - Address mapping between tags and virtual Modbus registers is explicit and inspectable in UI.

- Decision: Decommission legacy flow pages.
  - Remove navigation and default operator entry for points/mappings/wizard pages.
  - Pipeline Studio becomes canonical entry for mapping operations.

## Risks / Trade-offs

- Risk: Higher initial implementation scope.
  - Mitigation: phase-by-phase task gating with tests per phase.
- Risk: Global tag edit may surprise operators.
  - Mitigation: explicit "global impact" warning + undo snapshot + audit entry.
- Risk: Animation can reduce clarity if overused.
  - Mitigation: animation tokens with strict timing and disable path.
- Risk: Source templates become outdated after schema/rule changes.
  - Mitigation: template version check and upgrade prompt before applying.
- Risk: Chunked commit may create mixed success states.
  - Mitigation: commit report includes per-chunk status and deterministic retry queue.
- Risk: Local Modbus server may conflict with existing port usage.
  - Mitigation: startup preflight checks, explicit bind failure messaging, and optional fallback override in advanced settings.

## Migration Plan

1. Build source template model and UI shell without removing existing logic.
2. Introduce typed occupancy rendering and conflict detection.
3. Integrate global tag inline editing with impact confirmation.
4. Implement commit queue and one-click validation/commit flow.
5. Remove legacy pages from navigation and route entry.
6. Validate keyboard-only operation and reduced-motion mode.
7. Enable template versioning, commit audit logging, and chunked retry.
8. Roll out initial type subset, then extend type matrix in phase 2.
9. Enable local Modbus server sink with mirrored tag-to-register mappings on port `5020`.

## Open Questions

- Should template ownership be per-user only, or also shared across organization?
- Should `string` type span be fixed (configurable) or require explicit length per source?
