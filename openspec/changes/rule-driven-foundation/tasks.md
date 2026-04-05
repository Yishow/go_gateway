## 1. Source-rule revision and candidate persistence

- [x] 1.1 Add source-rule revision persistence and ensure each saved rule writes a new revision id
- [x] 1.2 Persist candidate snapshots with explicit target states and restore them after restart
- [x] 1.3 Implement canonical rule-owned identity and proposed-signature generation for rule-derived candidates

## 2. Shared readiness and diagnostics contracts

- [x] 2.1 Update device readiness logic to return planning, activation, and apply eligibility separately
- [x] 2.2 Extend connector diagnostics to return planning capability hints alongside connect/probe results
- [x] 2.3 Wire source-rule activation checks to the new readiness contract without blocking planning saves

## 3. Shared mapping and API contracts

- [x] 3.1 Add `out_of_sync` lifecycle handling for rule-derived mappings using the shared identity/signature contract
- [x] 3.2 Add rule-scoped candidate query and recompute endpoints with revision-aware responses
- [x] 3.3 Ensure recompute updates the active candidate snapshot without inventing target-specific apply semantics

## 4. Validation

- [x] 4.1 Add tests for source-rule revision persistence and restart restore
- [ ] 4.2 Add tests for readiness contract transitions across connect/probe outcomes
- [ ] 4.3 Add tests for candidate query/recompute APIs and mapping `out_of_sync` transitions
