## 1. Persistent Local Modbus state

- [x] 1.1 Persist rule-owned Local Modbus mapping ownership and revision metadata
- [x] 1.2 Restore Local Modbus candidate and applied mapping state after restart
- [x] 1.3 Bind Local Modbus candidates to the active source-rule revision and effective tag state

## 2. Conflict-governed review/apply flow

- [ ] 2.1 Detect overlapping register allocations and mark only the conflicting subset as `blocked_conflict`
- [ ] 2.2 Update the Local Modbus review surface to show conflict reasons, verification state, and health feedback inline
- [ ] 2.3 Keep database candidate/apply state untouched during Local Modbus review and apply actions

## 3. Local Modbus apply APIs

- [ ] 3.1 Add rule-scoped Local Modbus apply endpoint with revision-aware request validation
- [ ] 3.2 Return per-item apply results for partial success and failure handling
- [ ] 3.3 Keep verification-first checks inside the Local Modbus apply flow instead of a separate manual-first path

## 4. Validation

- [ ] 4.1 Add tests for restart restore of Local Modbus rule-owned state
- [ ] 4.2 Add tests for conflict detection and selective blocking behavior
- [ ] 4.3 Add tests proving Local Modbus apply does not affect Database state
