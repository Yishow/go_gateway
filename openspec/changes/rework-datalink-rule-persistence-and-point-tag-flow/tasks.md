## 1. Rule lifecycle foundation

- [x] 1.1 Add persisted source-rule schema, models, and repositories for rule definition, enabled state, and derived linkage metadata
- [x] 1.2 Implement backend rule services and APIs for create/update, enable/disable, and restart restoration
- [x] 1.3 Gate rule activation on device probe success and align runtime bootstrap with persisted rule state

## 2. Step 1 device diagnostics

- [x] 2.1 Extend device test responses to report separate connect/probe results with actionable error classification
- [x] 2.2 Update Step 1 UI to show phase-aware diagnostics and allow save-without-activation when probe fails

## 3. Step 2 rule-driven source planning

- [x] 3.1 Bind Step 2 planner behavior to selected device capability context and persisted rule data
- [x] 3.2 Render clear `planned`, `used`, `conflict`, `unmanaged`, and drift explanations on the source grid
- [x] 3.3 Feed runtime live values and restored rule state back into the Step 2 grid and inspector

## 4. Step 3 Point/Tag review flow

- [x] 4.1 Auto-create or synchronize rule-derived Tag and Mapping records with strict `1 Point : 1 Tag` enforcement
- [x] 4.2 Convert Step 3 from manual first-pass binding into review/verification plus exception handling

## 5. Output and database state model

- [x] 5.1 Refactor Local Modbus output state so each target keeps isolated selection state and one authoritative bind/unbind model
- [x] 5.2 Rebuild Database output flow into explicit Connector / Schema / Mapping layers for SQLite and PostgreSQL
- [x] 5.3 Preserve direct binding feedback, conflict handling, and compatibility behavior across Output targets

## 6. Validation and migration

- [x] 6.1 Add regression coverage for connect/probe gating, rule persistence, restart restoration, and auto-generated tag mappings
- [x] 6.2 Add migration and UX coverage for unmanaged legacy points, output-state drift, and layered database mapping flows
- [x] 6.3 Update operator-facing docs and OpenSpec-linked guidance for the new rule-driven workflow
