## 1. Tag candidate and review-decision persistence

- [x] 1.1 Change rule-generated tag handling to remain candidate-only until explicit tag apply
- [x] 1.2 Persist rename, skip, and override decisions against the canonical rule-owned tag identity
- [x] 1.3 Mark prior review decisions stale when a later revision changes the tag identity materially

## 2. Step 3 review/apply flow

- [x] 2.1 Update Step 3 to load tag candidates and pending mapping intent from the active source-rule revision
- [x] 2.2 Add stale revision detection and block tag apply until the review surface refreshes
- [x] 2.3 Add exception-handling actions for rename, skip, override, and batch review feedback

## 3. Tag apply APIs

- [ ] 3.1 Add rule-scoped tag apply endpoint with revision-aware request validation
- [ ] 3.2 Return per-item apply results without auto-applying database or Local Modbus outputs
- [ ] 3.3 Recompute downstream candidate state after successful tag apply without redefining the shared candidate contract

## 4. Validation

- [ ] 4.1 Add tests for rename/skip/override carry-forward across compatible revisions
- [ ] 4.2 Add tests for stale revision blocking in Step 3
- [ ] 4.3 Add tests for tag apply partial-result responses and post-apply candidate refresh
