# Design: Studio Grouped DB Row Planning

## Context

`/studio` already runs as the single operator route, but its Source, Tag, and Database surfaces do not yet explain one continuous planning story. The current database output path also writes each tag event independently, so a logical device group such as `meter/*` cannot become one timestamped database row without manual per-tag mental assembly.

This change must preserve the real backend contract, the shared `/studio` route, and the current Local Modbus + Database Output workspace while adding grouped-row planning. Product decisions already fixed for this change are:

- prefix inference starts from the first slash in the tag key
- slashless tags and legacy mappings remain single-member by default
- connector default interval starts at 15 seconds
- either Local Modbus or Database can complete the Output mainline

## Goals / Non-Goals

**Goals:**

- Let operators move from Source point planning into Tag grouping and then Database row planning without losing context.
- Support grouped database rows with explicit compatibility rules and interval-driven writes.
- Preserve backward compatibility for existing database mappings and slashless tag keys.
- Keep Output truthful: show both targets, but mark the mainline complete once either target is configured.

**Non-Goals:**

- Do not add `/studio-v*`, `/studio/database`, or any grouped-db-only product route.
- Do not replace the existing connector/mapping model with a new group table hierarchy in this change.
- Do not require Local Modbus and Database to complete together.
- Do not introduce multi-level grouping beyond the first slash inference in this first pass.

## Decisions

### 1. Extend existing connector and mapping models instead of inventing group tables

`DatabaseConnector` gains `default_write_interval_seconds`. `DatabaseTargetMapping` gains nullable `group_key` and optional `write_interval_seconds`.

- **Why:** the current API, repository, and UI already revolve around connector + mapping records. Extending those records keeps the implementation aligned with the live contract and avoids creating a second planning system.
- **Alternative rejected:** a new `database_target_groups` / `group_members` model. It would model grouping more explicitly, but it would add a second persistence abstraction and a larger API/UI rewrite before proving the workflow value.
- **Migration rule:** existing connectors are backfilled to `default_write_interval_seconds = 15`, while existing mappings keep `group_key = null` so they continue to behave as legacy single-member rows.

### 2. Use prefix inference as the initial grouping rule, but keep overrides operator-owned

The first path segment before `/` becomes the suggested `group_key`; the remaining suffix becomes the suggested `column_name` after lowercasing it and converting any additional `/` separators to `_`. Keys without `/` keep `group_key = null` and remain single-member until explicitly regrouped.

- **Why:** prefix inference covers the confirmed `meter/*` workflow immediately and stays easy to explain in Source, Tag, and Output.
- **Alternative rejected:** deeper automatic hierarchy or mandatory manual grouping. Deeper inference is harder to verify up front, while manual-only grouping adds too much friction to the default path.

### 3. Flush grouped rows on a timer, not only on the next event

The effective interval resolves as `mapping.write_interval_seconds` when present, otherwise `connector.default_write_interval_seconds`. The writer first filters mappings into compatibility sets that share connector, schema, table, write mode, timestamp column, and effective interval. It then buffers each compatible grouped set by `(connector, schema, table, group, interval bucket)` where `interval bucket = observed_at.UTC().Truncate(effective_interval)`, and flushes due buckets on a timer plus shutdown. Mappings with `group_key = null` keep the existing single-member path.

- **Why:** the product requirement is one row every 15 seconds by default, which cannot rely on a later event to trigger flush. Compatibility gating also prevents silent cross-table or cross-mode merges.
- **Alternative rejected:** lazy bucket rollover on the next write. It leaves quiet groups unwritten until another event arrives and breaks the 15-second contract.

### 4. Keep one Output workspace with per-target truth and one-target completion

Source shows the grouped Tag handoff, Tag shows inferred groups and overrides, and the Database planner becomes grouped-row-first. Output remains one workspace that shows Local Modbus and Database side by side, but the readiness model marks the mainline complete when either target is configured.

- **Why:** this matches the existing `/studio` route ownership and the accepted product rule that only one output target is required to complete the mainline.
- **Alternative rejected:** splitting Database into its own primary route or requiring both targets before completion. Both options reintroduce route or workflow fragmentation that this repo is explicitly avoiding.

## Risks / Trade-offs

- **Grouped flush compatibility mistakes** → Validate connector/schema/table/write-mode/timestamp/interval alignment before apply and surface explicit blocking reasons in dry-run/apply.
- **Legacy mapping regression** → Keep `group_key = null` as the compatibility marker and cover slashless + existing mappings with regression tests.
- **Timer-driven writer complexity** → Use a small, explicit due-row flush path and inject clock/tick control for tests.
- **UI overload in Output** → Keep grouping logic in dedicated view-model helpers instead of scattering it across component state.

## Migration Plan

1. Add the new OpenSpec artifacts for API, workbench shell, and database planner behavior.
2. Add schema migrations so connectors default to `15` seconds and mappings can persist nullable `group_key` plus optional interval overrides.
3. Implement candidate/apply, writer, and UI changes behind the existing `/studio` route and shared APIs.
4. Validate with backend/frontend tests plus integrated browser evidence across Device → Source → Tag → Output → Shell.
5. Roll back by reverting the change set and running the database down migration if the new grouped metadata causes unacceptable regressions.

## Open Questions

- No blocking open questions remain for this change. Future follow-up may expand grouping beyond first-slash inference or introduce richer row-composition tools, but those are intentionally deferred.
