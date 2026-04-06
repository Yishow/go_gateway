## Why

Under umbrella `rule-centric-studio-refactor`, once `rule-driven-foundation` exists, the next blocking workflow gap is Step 3. The current product still treats tags and mappings as if rule generation should become active downstream state immediately, but the approved workflow requires Step 3 to be a review/apply surface that preserves rename, skip, and override decisions across compatible revisions.

## What Changes

- Execute this change as child phase 2 under umbrella `rule-centric-studio-refactor`.
- Treat `rule-driven-foundation` as the prerequisite child change for this phase.
- Convert rule-generated tags and pending mappings into explicit review/apply candidates instead of immediate active downstream records.
- Make Step 3 of the workbench a dedicated tag review/apply surface with exception handling.
- Persist tag review decisions (`rename`, `skip`, `override`) across compatible revisions using the shared identity contract from the foundation phase.
- Add explicit rule-scoped tag apply APIs with per-item results.
- **BREAKING**: rule-generated tags and mappings are no longer considered active merely because a source rule was saved or enabled.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `tag-dictionary`: Change rule-generated tag behavior from immediate active creation to candidate review/apply with preserved review decisions.
- `datalink-workbench-desktop`: Change Step 3 from a manual-first binding surface into a candidate review/apply surface.
- `datalink-api`: Add explicit rule-scoped tag apply contracts.

## Impact

- **Frontend**: Step 3 review board, stale-state handling, exception actions, and batch/apply feedback.
- **Backend services**: tag decision persistence, rule-generated tag candidate handling, and tag apply orchestration.
- **APIs**: rule-scoped tag apply endpoint and response contract.
- **Tests**: rename/skip/override carry-forward, stale revision blocking, and tag apply result coverage.
