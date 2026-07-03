## 1. Runtime projection contract

- [x] 1.1 Deliver Runtime projection is built from persisted workspace state by defining one persisted workspace projection model for activation and restart, and verify it with runtime or workspace service tests that compare activation-time and restart-time reconstruction.
- [x] 1.2 Deliver Rule lifecycle changes trigger runtime reconciliation results under Workspace projection is the only runtime source of truth by converting source-rule lifecycle mutations into explicit runtime reconciliation outcomes, and verify it with source-rule/runtime integration tests.

## 2. Targeted reconcile behavior

- [x] 2.1 Deliver Persisted config changes produce an explicit reconciliation outcome and Reconciliation is targeted and idempotent by implementing targeted reconcile flows for affected runtime scope only, and verify it with tests that repeat the same reconcile request without duplicating scheduler or binding state.
- [x] 2.2 Deliver Live apply reports runtime reconciliation outcome under Live apply returns explicit reconcile outcome by extending live-apply responses with reconcile status for running workspaces, and verify it with handler tests and frontend hook tests that assert aligned versus deferred outcomes.

## 3. Runtime alignment surfaces

- [x] 3.1 Deliver Workspace runtime view exposes projection alignment state and Runtime snapshot and stream expose projection reconciliation state under Runtime surfaces projection drift instead of hiding it by wiring projection alignment into runtime APIs and workspace/runtime UI, and verify it with runtime handler tests plus runtime dashboard state tests.
