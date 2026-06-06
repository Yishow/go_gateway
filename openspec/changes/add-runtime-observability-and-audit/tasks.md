## 1. Diagnostics stage model

- [ ] 1.1 Deliver Runtime diagnostics expose collector to database delivery stages by defining a stage-oriented diagnostics model over collector, mapping, runtime projection, and database delivery, and verify it with backend tests that assert the failing stage is reported for broken flows.
- [ ] 1.2 Deliver Runtime diagnostics keep recent success and failure context under Diagnostics use a stage-oriented model by recording latest success and failure timestamps and reasons for selected scope, and verify it with runtime or writer tests that update diagnostics on both success and failure paths.

## 2. Audit history contract

- [ ] 2.1 Deliver Workspace audit history records critical lifecycle events and Workspace audit history records delivery-impacting configuration changes under Audit history records state transitions and delivery-impacting events by persisting activation, runtime transition, and delivery-impacting persisted changes into audit history, and verify it with handler or service tests that query recorded entries after writes.

## 3. Operator-facing surfaces

- [ ] 3.1 Deliver Runtime dashboard backend exposes diagnostics summary and Runtime dashboard shows operator-facing diagnostics under Operator surfaces consume summaries, not raw logs by adding diagnostics summaries to runtime-facing APIs and panels, and verify it with runtime dashboard state/page tests plus handler tests.
- [ ] 3.2 Deliver Workspace surfaces recent activation and audit history by exposing recent activation outcomes and audit summaries inside Studio V2 workspace surfaces, and verify it with workspace-facing handler tests and a manual assertion that the operator can inspect recent lifecycle events after refresh.
