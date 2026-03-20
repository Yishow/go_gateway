## Why

The current datalink workbench still behaves like a set of loosely connected steps instead of one coherent operating model. Device connectivity diagnostics, source rule execution, Point/Tag lifecycle, Local Modbus binding, and Database output state all drift in ways that make real-world operation hard to trust.

This change is needed now because the product direction is no longer just UI polish. The selected flow requires persisted source rules, connect/probe diagnostics, automatic Point→Tag linkage, and stable Output/Database state semantics before further rollout or embedded delivery can be trusted.

## What Changes

- Rework Step 1 device testing so `transport connect` and `protocol probe` are displayed and judged separately.
- Allow devices to be saved when transport connect succeeds but protocol probe fails, while blocking rule activation and data collection until probe succeeds.
- Introduce persisted source-rule lifecycle behavior so rules are stored in the database, can be enabled/disabled, drive runtime collection, and restore their enabled state after restart.
- Align Step 2 source planning with the selected device capability context, including address model, merge/data-type behavior, rule options, warning/block semantics, and clearer explanations for `planned`, `used`, `conflict`, and drift states.
- Shift the Point/Tag flow toward automatic linkage: creating a rule should automatically create the derived Point/Tag/Mapping relationship, and Step 3 should become a review/verification surface rather than a mandatory manual binding step.
- Preserve strict `1 Point : 1 Tag` semantics so a Tag cannot be bound to multiple Points.
- Rework Step 4 Local Modbus output state so each target has isolated selection state and all register/tag/binding state uses a single source of truth.
- Rework the Database output workflow for SQLite and PostgreSQL around explicit `Connector / Schema / Mapping` layering instead of one tightly coupled screen state.

## Capabilities

### New Capabilities
- `source-rule-runtime`: Persisted source-rule lifecycle, runtime execution, enable/disable semantics, restart restoration, and derived collection behavior.
- `database-target-workbench`: Workbench database target model for SQLite and PostgreSQL with explicit connector/schema/mapping layering and predictable binding behavior.

### Modified Capabilities
- `datalink-workbench-desktop`: Step 1–4 workbench behavior changes, including diagnostics, rule-driven flow, Step 3 review semantics, and stable output-state handling.
- `protocol-connectors`: Connection testing must distinguish transport connectivity from protocol probe outcomes and return more actionable diagnostics.
- `point-catalog`: Points become more explicitly rule-derived/runtime-facing, and source planning conflict semantics must align with persisted rule behavior.
- `tag-dictionary`: Tag linkage changes toward automatic rule-driven creation with strict `1 Point : 1 Tag` behavior and review-first editing.
- `local-modbus-memory-workbench`: Output-step Local Modbus binding behavior changes to isolate target state and use a consistent bind/unbind source of truth.

## Impact

- Affected frontend flows:
  - `WorkbenchDeviceStep`
  - `SourceCanvasSection`
  - `TagBindingStudio`
  - `LocalModbusBoard`
  - `DatabaseTargetBoard`
  - `WorkbenchProvider` and related workbench state models
- Affected backend/domain areas:
  - device connection testing and protocol diagnostics
  - source-rule persistence and runtime hooks
  - point/tag/mapping lifecycle behavior
  - database target configuration and mapping behavior
- Affected APIs:
  - device test result payloads
  - new or revised rule lifecycle endpoints
  - automatic mapping/tag creation flows
  - database target layering and validation flows
- Affected systems:
  - runtime collection / scheduler bootstrap
  - restart recovery semantics
  - Local Modbus output state management
  - SQLite and PostgreSQL database target handling
