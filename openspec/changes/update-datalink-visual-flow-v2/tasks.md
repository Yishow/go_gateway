## 1. Spec & UX Contract
- [ ] 1.1 Align with active changes and confirm no conflict with `optimize-point-configuration-ux`
- [ ] 1.2 Finalize flow model and visual states (draft/validated/active/error)
- [ ] 1.3 Freeze UI contract for Source → Grid → Tag → DB visualization

## 2. Information Architecture
- [ ] 2.1 Refactor SmartDashboard into flow-first workspace regions
- [ ] 2.2 Keep source context visible while editing mapping/tag details
- [ ] 2.3 Add responsive behavior for 1024px, 1280px, 1440px without horizontal scroll

## 3. Flow State & Observability
- [ ] 3.1 Introduce unified flow state representation in frontend state layer
- [ ] 3.2 Show latest value, quality, timestamp, and last error per flow segment
- [ ] 3.3 Surface segment-level diagnostics (source/transform/sink)

## 4. Interaction & Accessibility
- [ ] 4.1 Ensure complete keyboard navigation for all core actions
- [ ] 4.2 Ensure icon-only controls have accessible labels
- [ ] 4.3 Ensure dialog/panel focus trap and focus restore behavior is consistent
- [ ] 4.4 Ensure visible focus states and contrast pass baseline checks

## 5. Validation Workflow
- [ ] 5.1 Add explicit validate transition from draft to validated
- [ ] 5.2 Add activation transition from validated to active
- [ ] 5.3 Add error recovery paths and retry UX for each segment

## 6. Verification
- [ ] 6.1 Add/adjust component tests for flow-first interactions
- [ ] 6.2 Add/adjust integration tests for mapping + preview + activation flow
- [ ] 6.3 Run lint/typecheck/tests and record results
- [ ] 6.4 Conduct manual UX verification for operator scenarios
