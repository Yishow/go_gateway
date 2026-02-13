## 1. Spec and Contract Lock
- [ ] 1.1 Confirm source template schema and storage contract.
- [ ] 1.2 Confirm typed occupancy span table and validation rules.
- [ ] 1.3 Confirm global tag inline-edit guardrails (impact warning + confirmation).
- [ ] 1.4 Confirm legacy page decommission scope (`points`, `mappings`, `wizard`).

## 2. Source Planning and Templates
- [ ] 2.1 Add source planner UI section (type, count, naming rules).
- [ ] 2.2 Add template CRUD (save, load, update, delete).
- [ ] 2.3 Wire template selection to planner defaults.
- [ ] 2.4 Add tests for planner-to-template persistence.

## 3. Memory Grid Typed Occupancy
- [ ] 3.1 Implement span calculation by type (1/2/4 cells).
- [ ] 3.2 Render grouped cells for multi-cell types (pair/quad visualization).
- [ ] 3.3 Show occupancy states (`available`, `planned`, `used`, `conflict`, `linked`).
- [ ] 3.4 Add tests for occupancy math and edge collision handling.

## 4. Tag Linkage and Global Editing
- [ ] 4.1 Build inline tag linkage panel bound to selected grid allocation.
- [ ] 4.2 Support creating or selecting existing global tags in context.
- [ ] 4.3 Enable global tag metadata edit in-place with impact summary.
- [ ] 4.4 Add tests for global tag edit propagation behavior.

## 5. Commit Workflow (Tag -> DB)
- [ ] 5.1 Introduce commit queue view for pending allocations.
- [ ] 5.2 Add `validate` and `commit` actions with segment-level feedback.
- [ ] 5.3 Add rollback/retry UX for partial failure.
- [ ] 5.4 Add integration tests for validate/commit/error-recovery lifecycle.

## 6. Motion-Guided Interaction
- [ ] 6.1 Define motion tokens and transition rules (150-300ms).
- [ ] 6.2 Add stage transitions and focus handoff animations.
- [ ] 6.3 Add reduced-motion fallback and tests.
- [ ] 6.4 Run accessibility checks for keyboard and focus visibility.

## 7. Route and Navigation Decommission
- [ ] 7.1 Remove deprecated nav entries and route entry points.
- [ ] 7.2 Redirect deprecated routes to Pipeline Studio with migration notice.
- [ ] 7.3 Update docs and help text to the one-screen workflow.

## 8. Verification and Release Gate
- [ ] 8.1 Run `npm run lint`, `npm run typecheck`, `npm test`.
- [ ] 8.2 Run manual scenario checks for type-count occupancy examples.
- [ ] 8.3 Record UX acceptance notes for source template reuse and global tag edits.
