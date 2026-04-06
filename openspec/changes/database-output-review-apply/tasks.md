## 1. Database candidate orchestration

- [x] 1.1 Bind database candidates to the effective post-tag-review state for the active source-rule revision
- [ ] 1.2 Persist connector, schema, table, and column context as the authoritative review/apply scope
- [ ] 1.3 Mark database candidates `blocked` or `out_of_sync` when referenced context becomes invalid

## 2. Database review/apply flow

- [ ] 2.1 Update the database output workspace to show rule-scoped candidates and current connector context together
- [ ] 2.2 Revalidate blocked database candidates automatically when the original connector context becomes valid again
- [ ] 2.3 Keep Local Modbus candidate/apply state untouched during database review and apply actions

## 3. Database apply APIs

- [ ] 3.1 Add rule-scoped database output apply endpoint with revision-aware request validation
- [ ] 3.2 Return per-item apply results for partial success and failure handling
- [ ] 3.3 Ensure apply does not silently rebind candidates to a different connector or table

## 4. Validation

- [ ] 4.1 Add tests for connector invalidation and revalidation
- [ ] 4.2 Add tests for revision-bound database candidate generation and apply
- [ ] 4.3 Add tests proving database apply does not affect Local Modbus state
