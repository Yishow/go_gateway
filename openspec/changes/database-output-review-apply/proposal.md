## Why

Under umbrella `rule-centric-studio-refactor`, after `rule-driven-foundation` and `tag-review-first-flow`, Database output is still the largest remaining manual-first gap. The workflow needs a rule-scoped database review/apply phase that consumes effective tag state, preserves connector context, and remains independent from Local Modbus so the two output targets can evolve as true siblings.

## What Changes

- Execute this change as child phase 3 under umbrella `rule-centric-studio-refactor`.
- Treat `rule-driven-foundation` and `tag-review-first-flow` as prerequisites, and treat `local-modbus-review-apply` as a sibling phase rather than a prerequisite.
- Convert database output from manual per-tag binding into rule-scoped candidate review/apply.
- Bind database review to the effective post-tag-review state for a source-rule revision.
- Preserve connector, schema, table, and column context as explicit review/apply scope.
- Revalidate blocked or `out_of_sync` database candidates when connector context becomes valid again.
- Add explicit rule-scoped database apply APIs with per-item results.
- **BREAKING**: database output for the primary workflow is no longer treated as ad hoc tag-to-column binding outside the rule revision contract.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `database-target-workbench`: Change database output from manual binding to rule-scoped candidate review/apply with connector-aware revalidation.
- `datalink-api`: Add database-specific rule-scoped apply contracts.

## Impact

- **Frontend**: database output review board, connector/context invalidation handling, and apply feedback.
- **Backend services**: database candidate orchestration, connector-context revalidation, and apply execution.
- **APIs**: database output apply endpoint and response contract.
- **Tests**: connector invalidation/revalidation, revision-bound database candidates, and partial apply results.
