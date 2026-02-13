## 1. Spec & UX Contract
- [x] 1.1 Align with active changes and confirm no conflict with `optimize-point-configuration-ux`
- [x] 1.2 Finalize flow model and visual states (draft/validated/active/error)
- [x] 1.3 Freeze UI contract for Source → Grid → Tag → DB visualization

## 2. Information Architecture
- [x] 2.1 Refactor SmartDashboard into flow-first workspace regions
- [x] 2.2 Keep source context visible while editing mapping/tag details
- [x] 2.3 Add responsive behavior for 1024px, 1280px, 1440px without horizontal scroll

## 3. Flow State & Observability
- [x] 3.1 Introduce unified flow state representation in frontend state layer
- [x] 3.2 Show latest value, quality, timestamp, and last error per flow segment
- [x] 3.3 Surface segment-level diagnostics (source/transform/sink)

## 4. Interaction & Accessibility
- [ ] 4.1 Ensure complete keyboard navigation for all core actions
- [x] 4.2 Ensure icon-only controls have accessible labels
- [x] 4.3 Ensure dialog/panel focus trap and focus restore behavior is consistent
- [ ] 4.4 Ensure visible focus states and contrast pass baseline checks

## 5. Validation Workflow
- [x] 5.1 Add explicit validate transition from draft to validated
- [x] 5.2 Add activation transition from validated to active
- [x] 5.3 Add error recovery paths and retry UX for each segment

## 6. Verification
- [x] 6.1 Add/adjust component tests for flow-first interactions
- [ ] 6.2 Add/adjust integration tests for mapping + preview + activation flow
- [x] 6.3 Run lint/typecheck/tests and record results
- [ ] 6.4 Conduct manual UX verification for operator scenarios
