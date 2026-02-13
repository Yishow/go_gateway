## 1. Spec and Contract Lock
- [x] 1.1 Confirm source template schema and storage contract.
- [x] 1.2 Confirm typed occupancy span table and validation rules.
- [x] 1.3 Confirm global tag inline-edit guardrails (impact warning + confirmation).
- [x] 1.4 Confirm legacy page decommission scope (`points`, `mappings`, `wizard`).

## 2. Source Planning and Templates
- [x] 2.1 Add source planner UI section (type, count, naming rules).
- [x] 2.2 Add template CRUD (save, load, update, delete).
- [x] 2.3 Wire template selection to planner defaults.
- [x] 2.4 Add tests for planner-to-template persistence.
- [x] 2.5 Add template version and last-used timestamp handling.
- [x] 2.6 Add stale-template warning and upgrade prompt.
- [x] 2.7 Add batch naming preview with duplicate/conflict detection.

## 3. Memory Grid Typed Occupancy
- [x] 3.1 Implement span calculation by type (1/2/4 cells).
- [x] 3.2 Render grouped cells for multi-cell types (pair/quad visualization).
- [x] 3.3 Show occupancy states (`available`, `planned`, `used`, `conflict`, `linked`).
- [x] 3.4 Add tests for occupancy math and edge collision handling.
- [x] 3.5 Add hard/soft conflict severity rendering and behavior.
- [x] 3.6 Add auto-allocation strategy for nearest valid contiguous span.
- [x] 3.7 Add grouped-border and in-group index markers for multi-cell types.
- [x] 3.8 Add conflict-only filter mode in grid controls.

## 4. Tag Linkage and Global Editing
- [x] 4.1 Build inline tag linkage panel bound to selected grid allocation.
- [x] 4.2 Support creating or selecting existing global tags in context.
- [x] 4.3 Enable global tag metadata edit in-place with impact summary.
- [x] 4.4 Add tests for global tag edit propagation behavior.
- [x] 4.5 Add affected-mappings summary and second-confirmation flow for global edits.
- [x] 4.6 Add before/after field-level diff preview before global tag save.

## 5. Commit Workflow (Tag -> DB)
- [x] 5.1 Introduce commit queue view for pending allocations.
- [x] 5.2 Add `validate` and `commit` actions with segment-level feedback.
- [x] 5.3 Add rollback/retry UX for partial failure.
- [x] 5.4 Add integration tests for validate/commit/error-recovery lifecycle.
- [x] 5.5 Execute commit in chunks with per-chunk result visibility.
- [x] 5.6 Add deterministic retry queue for failed chunks only.
- [x] 5.7 Add commit impact summary (new points, global tag updates, conflicts).
- [x] 5.8 Add audit logging payload and UI trace link.
- [x] 5.9 Add commit batch trace id creation and display.
- [x] 5.10 Add two-stage validation flow (structural then executable).
- [x] 5.11 Add pre-commit polling load delta estimate.

## 6. Motion-Guided Interaction
- [ ] 6.1 Define motion tokens and transition rules (150-300ms).
- [ ] 6.2 Add stage transitions and focus handoff animations.
- [ ] 6.3 Add reduced-motion fallback and tests.
- [ ] 6.4 Run accessibility checks for keyboard and focus visibility.
- [ ] 6.5 Restrict animations to intent transitions; no decorative long animations.
- [ ] 6.6 Add animation readability checklist and QA gate for visual clarity.

## 7. Local Modbus Sharing Mode
- [x] 7.1 Add optional sink target `local_modbus_server`.
- [x] 7.2 Implement tag-to-virtual-register mapping configuration UI.
- [x] 7.3 Mirror committed tag values to local virtual Modbus memory grid.
- [x] 7.4 Start local Modbus TCP server on port `5020`.
- [x] 7.5 Add bind preflight check and actionable port-conflict error handling.
- [x] 7.6 Add integration tests for external read verification against mirrored values.

## 8. Route and Navigation Decommission
- [x] 8.1 Remove deprecated nav entries and route entry points.
- [x] 8.2 Redirect deprecated routes to Pipeline Studio with migration notice.
- [x] 8.3 Update docs and help text to the one-screen workflow.

## 9. Verification and Release Gate
- [x] 9.1 Run `npm run lint`, `npm run typecheck`, `npm test`.
- [x] 9.2 Run manual scenario checks for type-count occupancy examples.
- [x] 9.3 Record UX acceptance notes for source template reuse and global tag edits.
- [x] 9.4 Verify autosave/restore behavior under abrupt page reload.
- [x] 9.5 Verify phased rollout with initial type subset (`int16`, `int32`, `float32`).
- [x] 9.6 Verify Modbus server behavior on port `5020` using an external Modbus client.
