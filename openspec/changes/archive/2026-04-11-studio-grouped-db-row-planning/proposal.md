# Proposal: Studio Grouped DB Row Planning

## Why

`/studio` currently makes Source, Tag, and Database work feel like separate local optimizations instead of one operator mainline. At the same time, the database path still behaves like one-tag-to-one-write, which does not match the product decision that related tags such as `meter/A1`, `meter/A2`, `meter/A3`, and `meter/kw` should land in one database row every 15 seconds by default.

## What Changes

- Add grouped database-row planning for related tags, starting with prefix inference from slash-based tag keys and operator-editable overrides.
- Extend database connector and mapping contracts with default and per-group interval metadata while preserving legacy single-member mappings.
- Change the database writer to flush one row per compatible group and interval bucket instead of writing one row per tag event.
- Make `/studio` clarify the Source → Tag → Database handoff, keep Local Modbus and Database visible together, and treat Output as mainline-complete when either target is configured.
- Keep one real `/studio` route, one real backend contract, and no grouped-db-only parallel route.

## Capabilities

### New Capabilities

- None.

### Modified Capabilities

- `datalink-api`: expose grouped database connector/mapping metadata, grouped candidate/apply payloads, and grouped write-history semantics.
- `datalink-workbench-desktop`: clarify Source → Tag → Database handoff, one-target Output completion, and blocker recovery cues tied to the true recovery surface.
- `database-target-workbench`: make the database planner grouped-row-first with interval controls and legacy single-member compatibility.

## Impact

- OpenSpec artifacts for `datalink-api`, `datalink-workbench-desktop`, and `database-target-workbench`
- Backend schema, dbtarget service/repository/writer, source-rule candidate/apply flow, and API handlers
- Frontend Source, Tag, Output, shell/readiness surfaces, plus i18n and regression coverage
