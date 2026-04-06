## Why

The current datalink workflow is split across two incompatible operating models: the backend already lets `SourceRule` derive runtime points and tag mappings, but the UI and API still expose `point`, `tag`, `mapping`, database output, and local Modbus output as separate manual-first workflows. That split causes state drift, inconsistent operator expectations, and makes `/studio` harder to turn into a clear product path.

This change is needed now because the project already has enough rule-aware foundations to support a rule-centric workflow, but it lacks one umbrella contract that aligns UI, API, runtime, and output targets around the same operating model. Without a master refactor change, later phase work will keep re-litigating scope, boundaries, and sequencing.

## What Changes

- Establish one umbrella refactor contract for a rule-centric `/studio` workflow driven by `Device -> SourceRule -> Tag review -> Output review/apply`.
- Define `Device` and `SourceRule` as the primary workflow inputs, while pushing `Point` into runtime/internal responsibility.
- Define a candidate/review/apply model for rule-derived tags and for both output targets, so downstream state can be regenerated safely after rule changes.
- Make Database and Local Modbus equal first-class output targets inside the same workflow, with target-isolated conflict and readiness behavior.
- Define phase sequencing for the five child changes:
  1. `rule-driven-foundation`
  2. `tag-review-first-flow`
  3. `database-output-review-apply`
  4. `local-modbus-review-apply`
  5. `studio-orchestration-cleanup`
- **BREAKING**: The primary `/studio` product flow will no longer treat manual `point/tag/mapping` construction as the main operator path.
- **BREAKING**: Local Modbus output moves from runtime-only manual behavior toward persistent, rule-owned review/apply behavior.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `datalink-ui`: The product-level datalink UI requirements change to require one rule-centric `/studio` flow instead of multiple manual-first workflow paths.
- `datalink-workbench-desktop`: The workbench requirements change to make `SourceRule` the primary source-planning object, Step 3 review-first, and Step 4 dual-target review/apply.
- `source-rule-runtime`: Source rules change from immediate downstream side effects to revision-aware orchestration with candidate regeneration and apply gating.
- `tag-dictionary`: Tags change from primarily manual binding targets to rule-derived review/apply assets with explicit rule ownership and revision tracking.
- `database-target-workbench`: Database output changes from per-tag manual binding to rule-derived candidate review/apply with connector-scoped context invalidation.
- `local-modbus-memory-workbench`: Local Modbus changes from manual/runtime-only mapping to persistent, rule-derived candidate review/apply with allocator and conflict governance.
- `device-registry`: Device requirements change to distinguish `connect` and `probe`, preserve saved devices after partial success, and expose capability context to downstream workflow steps.
- `protocol-connectors`: Connector requirements change to support explicit probe readiness semantics and capability hints consumed by rule-driven planning.
- `mapping-pipeline`: Mapping requirements change to align transform ownership with rule-derived mappings, revision diffs, and override-safe apply behavior.
- `datalink-api`: API requirements change to add rule-scoped candidate/recompute/apply contracts instead of relying on lower-level CRUD as the workflow backbone.

## Superseded Legacy Assumptions

- `/studio` supersedes `/datalink/workbench` as the primary product workflow route.
- Manual `point -> tag -> mapping` assembly is no longer the primary product path.
- Source-rule save generates reviewable candidates and pending mapping intent, not immediate active downstream persistence.
- Rule-derived mappings and outputs gain explicit `out_of_sync` behavior instead of silent overwrite.

## Impact

- **Frontend**: `/studio`, workbench step composition, review/apply state models, and route-level orchestration.
- **Backend services**: `device`, `sourcerule`, `tag`, `mapping`, `dbtarget`, `modbusshare`, and any orchestration added around them.
- **Persistence**: rule ownership metadata, revision tracking, candidate storage, and persistent Local Modbus mappings.
- **APIs**: new rule-scoped candidate and apply endpoints plus updated readiness semantics.
- **Tests**: unit, integration, and E2E coverage across device gating, candidate regeneration, tag review, and both output targets.
- **Delivery model**: this umbrella change becomes the parent contract for the five executable child changes listed above.
