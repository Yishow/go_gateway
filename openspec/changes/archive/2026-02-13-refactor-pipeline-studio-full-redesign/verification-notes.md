# Verification Notes

## 2026-02-13

### 9.2 Type-count occupancy examples
- Verified by automated UI regression tests:
  - `frontend/src/components/datalink/__tests__/MemoryGrid.test.tsx`
  - Scenario: `int16 x5` renders exactly 5 planned cells.
  - Scenario: `float32 x10` renders 20 planned cells with pair markers (`1/2`, `2/2`).

### 9.4 Autosave / restore after reload
- Verified by storage module round-trip test:
  - `frontend/src/features/datalink/__tests__/sourceTemplateStorage.test.ts`
  - `saveSourceTemplates()` persists templates to local storage.
  - `loadSourceTemplates()` restores persisted templates and handles invalid JSON safely.

### 9.6 External Modbus client readback
- Verified by backend integration test:
  - `internal/datalink/modbusshare/service_test.go`
  - Scenario: start local Modbus share server, mirror tag value, read holding registers through TCP client (`internal/protocol/modbus`).

### 9.5 Initial phased rollout subset
- Verified in Pipeline Studio source planner UI:
  - `frontend/src/pages/datalink/SmartDashboard.tsx`
  - Current selectable source data types are limited to `int16`, `int32`, and `float32` for phased rollout.

### 9.3 UX acceptance notes (current status)
- Source template reuse:
  - Save/load/delete flow is available in the Source Planner panel.
  - Persistence across reload is covered by `sourceTemplateStorage` tests.
- Global tag inline-edit:
  - Current Pipeline Studio flow supports tag linkage and modbus mapping actions.
  - Dedicated global tag metadata impact-summary + second-confirmation acceptance flow still requires browser walkthrough and explicit UX sign-off.

### 2.5 / 2.6 Template versioning and stale upgrade prompt
- Added schema version and last-used timestamp fields in:
  - `frontend/src/features/datalink/sourceTemplateStorage.ts`
- Added stale detection and one-click template upgrade prompt in:
  - `frontend/src/pages/datalink/SmartDashboard.tsx`
- Regression tests:
  - `frontend/src/features/datalink/__tests__/sourceTemplateStorage.test.ts`
  - Covers legacy template normalization (missing version/lastUsedAt) and upgrade behavior.

### 2.7 Batch naming preview and conflict detection
- Added batch naming preview utility:
  - `frontend/src/features/datalink/batchNaming.ts`
- Added planner UI for prefix input and preview chips with conflict highlighting:
  - `frontend/src/pages/datalink/SmartDashboard.tsx`
- Regression tests:
  - `frontend/src/features/datalink/__tests__/batchNaming.test.ts`
  - Covers sequence generation, case-insensitive conflicts, and blank-prefix fallback.

### 1.1 Source template schema and storage contract
- Contract helper:
  - `frontend/src/features/datalink/sourcePlannerContract.ts`
  - Locks template payload shape (`id`, `name`, `dataType`, `count`, `startAddress`, `updatedAt`, `lastUsedAt`, `version`) and storage-validity checks.
- Regression tests:
  - `frontend/src/features/datalink/__tests__/sourcePlannerContract.test.ts`
  - `frontend/src/features/datalink/__tests__/sourceTemplateStorage.test.ts`

### 1.2 Typed occupancy span table and validation rules
- Span and validation rules:
  - `frontend/src/features/datalink/typedOccupancy.ts`
  - Added typed plan validation (`count` range + span validity) with explicit errors.
- Regression tests:
  - `frontend/src/features/datalink/__tests__/typedOccupancy.test.ts`

### 1.3 Global tag inline-edit guardrails
- Guardrail helper:
  - `frontend/src/features/datalink/globalTagGuardrails.ts`
  - Provides affected-mappings warning and second-confirmation requirement signal.
- Integrated in UI:
  - `frontend/src/pages/datalink/SmartDashboard.tsx`
- Regression tests:
  - `frontend/src/features/datalink/__tests__/globalTagGuardrails.test.ts`

### 1.4 Legacy decommission scope lock
- Scope lock helper:
  - `frontend/src/features/datalink/legacyRoutes.ts`
  - Fixed allowed legacy scope to `points`, `mappings`, `wizard`.
- Route integration:
  - `frontend/src/App.tsx` uses shared redirect builder to avoid scope drift.
- Regression tests:
  - `frontend/src/features/datalink/__tests__/legacyRoutes.test.ts`

### 2.1 / 2.2 / 2.3 Source planner and template flow
- Source Planner already in Pipeline Studio and now fully contract-driven:
  - `frontend/src/pages/datalink/SmartDashboard.tsx`
  - Added naming-rule hint, count clamping (`1..200`), normalized prefix, and template CRUD helper wiring.
- Planner/template contract helpers:
  - `frontend/src/features/datalink/sourcePlannerContract.ts`

### 2.4 Planner-to-template persistence
- Added planner-template persistence regression:
  - `frontend/src/features/datalink/__tests__/sourcePlannerContract.test.ts`
  - Covers planner draft -> template create/upsert -> localStorage save/load -> planner defaults restore.
