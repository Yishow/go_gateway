# measurement-semantics Specification

## Purpose

TBD - created by archiving change 'add-measurement-semantics'. Update Purpose after archive.

## Requirements

### Requirement: Per-measurement recording meaning
The system SHALL persist each measurement's quantity, unit, semantic_kind and source binding independently of protocol read type and database column.

#### Scenario: Mixed measurements on one device
- **WHEN** one device supplies temperature, pressure, flow rate and cumulative volume
- **THEN** each measurement retains its own meaning and recording eligibility
- **AND** no device-wide counter or average policy overrides the individual definitions.

#### Scenario: Unconfirmed semantics
- **WHEN** an existing numeric point has no confirmed semantic_kind
- **THEN** legacy recording remains unchanged and new usage or integration modes are blocked with an actionable reason.

---
### Requirement: Stable identity and measurement epochs
The system MUST preserve stable identity for cosmetic edits and start a new series_epoch when source identity or numerical comparability changes.

#### Scenario: Rename versus replace
- **WHEN** a measurement is renamed
- **THEN** its history and sample identity remain intact
- **AND** replacing its meter or changing its scale requires explicit epoch transition without joining old counter baselines.

---
### Requirement: Truthful typed samples
The system SHALL emit the SampleEnvelope from the shared contract, including provenance, original observation time, quality and exact numeric encoding.

#### Scenario: Partially failed acquisition
- **WHEN** pressure fails while temperature succeeds in one collection cycle
- **THEN** temperature is eligible for recording and pressure is marked missing or invalid
- **AND** an old pressure value or zero is not substituted as a new good reading.

#### Scenario: Integer above JavaScript safe range
- **WHEN** uint64 value 9007199254740993 crosses backend, JSON and UI
- **THEN** its decimal digits are preserved exactly and numeric encoding is explicit.

---
### Requirement: Stateful and non-numeric intent is explicit
The system MUST distinguish monotonic counters, signed totals, interval deltas, states, events and text; semantic conversion SHALL require validation rather than unit-name inference.

#### Scenario: Net energy falls
- **WHEN** a confirmed signed_counter decreases from 100 to 98
- **THEN** the definition permits a signed difference rather than automatically declaring a reset.

#### Scenario: Polled delta without identity
- **WHEN** a held PLC register is declared delta without a trustworthy interval or event identity
- **THEN** usage accumulation is blocked rather than counting every poll as new consumption.

---
### Requirement: Confirmed template intent and recoverable failures
The system SHALL apply a real unconfirmed template preview only after an explicit current-device confirmation, invalidate previews on device/template changes, and expose safe actionable load, preview and apply errors.

#### Scenario: Explicit confirmation of a real preview
- **WHEN** a current preview has confirmed=false and the operator confirms it
- **THEN** the client SHALL submit confirmed=true and the backend SHALL continue rejecting an unconfirmed submission.

#### Scenario: Changed scope or failed operation
- **WHEN** device/template changes or load/preview/apply fails through disconnection, conflict or server failure
- **THEN** stale previews and replies SHALL NOT be applied and the operator SHALL see a safe retry action without raw exception details.

#### Scenario: Offline intent and template reload
- **WHEN** an explicit template operation is attempted offline or a failed template load is retried
- **THEN** the operation SHALL fail visibly rather than queue consent for automatic reconnect submission, and loading SHALL preserve the expanded repair section.

<!-- @trace
source: fix-studio-v2-five-review-gaps
updated: 2026-10-23
code:
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel.tsx
  - docs/releases/2026-10-23-studio-five-review-fixes.md
  - frontend/src/i18n/locales/en/studio-review.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/datalink/collector/scheduler_polling.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceMeasurements.ts
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection.tsx
  - docs/releases/evidence/studio-five-review-2026-10-23/manifest.json
  - internal/datalink/collector/scheduler_manual_poll.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor.tsx
  - internal/datalink/device/service_point_read.go
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - internal/datalink/connector/point_read.go
  - internal/datalink/recordingplan/schema_preview.go
  - docs/releases/evidence/studio-five-review-2026-10-23/observations.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupMemberTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MeasurementTemplateSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/state/writeGroup/columns.ts
  - frontend/src/i18n/locales/zh-TW/studio-review.json
tests:
  - frontend/tests/unit/workbench-v2/template-offline.test.tsx
  - frontend/tests/unit/workbench-v2/occupied-range.test.tsx
  - internal/datalink/collector/point_count_contract_test.go
  - frontend/tests/unit/workbench-v2/MeasurementTemplateSelector.test.tsx
  - frontend/tests/unit/workbench-v2/basic-required-policy.test.tsx
  - internal/datalink/recordingplan/schema_noop_json_test.go
  - scripts/tests/f_device_to_sql/five-template-errors.mjs
  - scripts/tests/f_device_to_sql/five-review-flow.mjs
  - scripts/tests/f_device_to_sql/mc-fixture.mjs
  - frontend/tests/unit/workbench-v2/writeGroupColumns.test.ts
  - frontend/tests/unit/runtime-dashboard/focused-header-evidence.test.tsx
-->