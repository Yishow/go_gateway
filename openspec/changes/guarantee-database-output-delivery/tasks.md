## 1. Live delivery target set

- [ ] 1.1 Deliver Database delivery operates only on the live enabled mapping set and Database target workflow uses the live visible target set by defining one filtered live target projection for schema ensure and runtime writer, and verify it with dbtarget service tests that exclude stale hidden mappings.
- [ ] 1.2 Deliver Rule-derived database target drift does not block unrelated live runtime under Delivery operates only on the live visible mapping set by isolating orphaned rule-derived targets from live activation and delivery scope, and verify it with source-rule/database integration tests.

## 2. Delivery truth and write outcome recording

- [ ] 2.1 Deliver Database delivery records last schema and write outcomes and Writer persists last-write status and flush outcome by persisting the latest schema ensure and write delivery result for connector or target scope, and verify it with dbtarget writer/tooling tests that assert success and failure state updates.
- [ ] 2.2 Deliver Runtime write delivery is diagnosable end to end under Schema ensure and write path share the same filtered target projection by surfacing where collected values fail inside the database delivery path, and verify it with integration tests that assert collected, mapped, and failed-write stages.

## 3. Step 4 delivery visibility

- [ ] 3.1 Deliver Step 4 shows database connector readiness and delivery truth under Step 4 surfaces delivery truth, not only config truth by wiring connector readiness and last delivery outcome into Step 4 UI, and verify it with frontend Step 4 tests plus a manual assertion that failed delivery truth is visible without opening raw logs.
