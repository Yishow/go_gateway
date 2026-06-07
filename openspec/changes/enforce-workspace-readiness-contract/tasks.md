## 1. Readiness domain model

- [x] 1.1 Deliver Workspace readiness evaluates all persisted setup steps by defining a normalized readiness summary and issue model over persisted Step 1 through Step 4 state, and verify it with backend unit tests that cover blocking and warning issue aggregation.
- [x] 1.2 Deliver Downstream integrity gaps surface as readiness issues before activation under Readiness evaluates step coverage and downstream integrity by turning missing rule-derived point, tag, mapping, and database-target relationships into normalized readiness issues, and verify it with source-rule/runtime-facing tests.

## 2. Activation and live-apply gating

- [x] 2.1 Deliver Activation blocks on blocking readiness issues under Backend readiness snapshot owns activation gating by making activation reject blocking issues before runtime start, and verify it with activation handler tests that assert returned issue codes and no runtime start on blockers.
- [x] 2.2 Deliver Live apply honors the workspace readiness contract under Blocking and warning issues are first-class by rejecting or deferring persisted changes that introduce blocking readiness issues, and verify it with live-apply tests that distinguish blocking from warning outcomes.

## 3. Frontend readiness surfaces

- [x] 3.1 Deliver Workspace exposes readiness summary for Studio V2 surfaces, Step 4 activation surface shows readiness blockers and warnings, and Readiness updates when persisted setup state changes under Frontend consumes the same readiness contract everywhere by wiring one readiness summary into shell, summary rail, and Step 4 activation UI, and verify it with frontend unit tests plus a manual assertion that all three surfaces refresh to the same issue set after a persisted setup change.
