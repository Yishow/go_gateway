# Tasks: Studio Grouped DB Row Planning

## 1. Backend contract and persistence

- [x] 1.1 Add connector and mapping fields for `default_write_interval_seconds`, nullable `group_key`, and optional `write_interval_seconds`, including migrations and repository coverage
- [x] 1.2 Expose grouped metadata through rule-scoped database candidate/apply APIs and preserve legacy `group_key = null` behavior
- [x] 1.3 Implement timer-driven grouped database writes plus grouped write-history semantics without breaking single-member mappings

## 2. Studio workflow surfaces

- [x] 2.1 Add Source-step handoff context for grouped Tag review
- [x] 2.2 Make Tag review group-aware with editable overrides for inferred database row groups
- [x] 2.3 Rework the Database Output surface into a grouped row planner with interval controls and per-target readiness

## 3. Flow alignment and validation

- [x] 3.1 Align step rail, context bar, and shell recovery cues to the grouped Output readiness truth
- [ ] 3.2 Add or update backend/frontend regression coverage for grouped planning, compatibility rules, and one-target completion
- [ ] 3.3 Run integrated `/studio` browser validation and capture evidence for Device → Source → Tag → Output → Shell
