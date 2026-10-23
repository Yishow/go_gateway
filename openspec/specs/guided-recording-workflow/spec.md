# guided-recording-workflow Specification

## Purpose

TBD - created by archiving change 'redesign-studio-recording-flow'. Update Purpose after archive.

## Requirements

### Requirement: Four-step intent-led setup
Operators MUST be able to save and navigate offline drafts with unresolved probes; draft navigation MUST NOT imply verified readiness, and only dependent reads or activation SHALL be blocked. The existing Studio V2 route SHALL guide device connection, acquisition points, Point-to-Tag mapping and grouped database output in four steps labelled 連接設備 → 選要讀的資料 → 命名與換算 → 儲存並開始記錄. Basic/Advanced presentation changes SHALL preserve configuration values and protocol/type capabilities. Basic typed Tag writing MUST remain usable without required reporting, retention setup, aggregation or physical measurement semantics. Explicitly selected derived operations SHALL require their own confirmed semantics. The system MUST retain the Go backend and existing React framework and MUST NOT add a parallel V3 or restore a dedicated legacy /studio route. The basic path SHALL derive and persist routine mappings, show only necessary operator choices and default new recording to one managed group per device without converting existing advanced configurations.

#### Scenario: Mixed sensor configured without SQL
- **WHEN** an operator maps temperature, pressure, flow and counter source values as basic typed Tags
- **THEN** the operator can select those Tags for a managed write group without hand-written SQL or mandatory derived-value policies; choosing derived usage then requests the required confirmed roles.

#### Scenario: Reload and edited connection
- **WHEN** an operator reloads setup or changes a saved connection identity
- **THEN** the same persisted validity and revision checks determine progress, and the old probe does not keep the changed device ready.

#### Scenario: Unavailable protocol
- **WHEN** MQTT or another protocol lacks the complete V2 setup/parser/collector capability
- **THEN** the UI clearly marks it unavailable with a reason rather than letting the user enter an unsupported range flow.

#### Scenario: Same address on different devices
- **WHEN** two devices both expose address 40001
- **THEN** the address conflict checker keeps them separate while still detecting true overlapping ranges within the same device/address space.

#### Scenario: Missing advanced measurement
- **WHEN** a saved Point-to-Tag mapping has no persisted MeasurementDefinition
- **THEN** basic raw group setup proceeds without fake measurement IDs; only explicitly selected advanced semantics remain blocked.

#### Scenario: Offline draft navigation
- **WHEN** a device is offline or its saved configuration has no successful probe
- **THEN** the operator can save a draft and use Next/Back to configure later steps, including after reload, while the device remains visibly unverified and its dependent activation is blocked with a repair action.

#### Scenario: Routine mapping does not require duplicate entry
- **WHEN** the operator confirms known points and supported types in the basic path
- **THEN** the existing mapping services persist real scoped Tag identities once, display actual reads and allow correction without a second manual mapping exercise.

#### Scenario: Existing advanced configuration is reopened
- **WHEN** an existing custom-table, multi-entity or cross-device group is loaded
- **THEN** its persisted identity and layout remain intact with accessible advanced controls; the basic default neither splits nor overwrites it.

#### Scenario: New basic setup is resumed
- **WHEN** a new per-device managed draft is reloaded or its save response is retried
- **THEN** the persisted (workspace_id, device_id, canonical managed role) key reuses the same group and mappings without duplicate writers or selecting another group. A changed destination requires an explicit CAS-checked draft and new preparation/apply evidence; accepted historical data keeps its original destination. Existing advanced groups remain outside automatic conversion.


<!-- @trace
source: close-studio-v2-simplification
updated: 2026-10-23
code:
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-verification-20261023.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/schemaOperationRecovery.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023.png
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel.tsx
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-logs-20261023.json.gz
  - docs/releases/studio-v2-change-closeout-20261023.md
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023-mapping.png
  - docs/plans/studio-v2-flow-completion/README.md
tests:
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - internal/datalink/recordingplan/schema_group_extend_test.go
  - frontend/tests/unit/workbench-v2/schemaOperationRecovery.test.ts
  - frontend/tests/unit/workbench-v2/groupSchemaPanel.test.tsx
-->

---
### Requirement: Safe reusable templates and overrides
The system SHALL preview template changes per device and preserve confirmed user overrides and stable identities.

#### Scenario: Reapply a template
- **WHEN** a template is applied again after a user edits a pressure multiplier
- **THEN** the difference is shown and the multiplier is not silently overwritten or duplicated.

---
### Requirement: Natural-language recording policies
The system SHALL distinguish acquisition, raw retention, summaries, transport batching and history retention with examples and a final readable summary.

#### Scenario: Five-second polling and one-minute summaries
- **WHEN** those settings are selected
- **THEN** the summary explains that sampled readings feed the minute statistics rather than recording only one instantaneous value per minute.

---
### Requirement: Truthful save and activation states
The system MUST separate saved draft, applied revision, collecting, durable, delivered and verified states, including bounded waiting and partial failure.

#### Scenario: Save during an active plan
- **WHEN** the operator edits a running plan
- **THEN** the current applied plan remains active until explicit validated apply and both versions are visible.

---
### Requirement: Actionable accessible failures
The system SHALL provide field-level reasons, focused repair actions and keyboard-accessible controls instead of unexplained disabled buttons or raw technical errors.

#### Scenario: Write permission denied
- **WHEN** a connected database refuses write access
- **THEN** the operator sees which destination needs permission repair and may retry it without redoing successful device setup.

---
### Requirement: Optional outputs remain independent
The workflow MUST allow a selected Local Modbus-only setup without requiring database configuration while retaining the existing Share activation gates.

#### Scenario: No database output selected
- **WHEN** the user selects only validated Local Modbus forwarding
- **THEN** database fields and tests are skipped and no database writer starts.

---
### Requirement: Complete basic preparation with one effective configuration
Basic setup SHALL complete managed schema preparation without entering Advanced. Connection controls SHALL configure only connection identity; the canonical group SHALL own recording members, storage and interval. Operator primary surfaces SHALL show names, values, types, units and truthful recording stages; technical IDs, API payloads and revisions SHALL remain available in diagnostics.

#### Scenario: Fresh basic recording
- **WHEN** an operator configures eight measurement points with type-appropriate register spans, maps tags and selects an empty managed SQLite destination
- **THEN** Basic provides schema preview and explicit confirmation and proceeds to timed one-row recording without Advanced.

#### Scenario: Inspect effective settings
- **WHEN** an operator prepares a canonical group
- **THEN** no competing connector table, INSERT/UPSERT or recording interval control is presented as effective, and technical identifiers remain inspectable through diagnostics.

#### Scenario: Destination before recording choices
- **WHEN** no database destination is saved and Share is not enabled
- **THEN** Basic first presents destination setup before recording members, frequency and start; a saved destination is shown as a summary with an edit entry preserving advanced connection and password semantics.

#### Scenario: Prepare action and focused repair
- **WHEN** Basic preparation is requested or a prepare_schema error needs repair
- **THEN** one preparation action loads the read-only preview in the Basic schema card, with explicit DDL confirmation still required and no forced switch into all Advanced settings.


<!-- @trace
source: close-studio-v2-simplification
updated: 2026-10-23
code:
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-verification-20261023.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/schemaOperationRecovery.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023.png
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel.tsx
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-logs-20261023.json.gz
  - docs/releases/studio-v2-change-closeout-20261023.md
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023-mapping.png
  - docs/plans/studio-v2-flow-completion/README.md
tests:
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - internal/datalink/recordingplan/schema_group_extend_test.go
  - frontend/tests/unit/workbench-v2/schemaOperationRecovery.test.ts
  - frontend/tests/unit/workbench-v2/groupSchemaPanel.test.tsx
-->

---
### Requirement: Selection is separate from acquisition

Step 2 point click SHALL select or inspect points without changing acquisition. Multi-selection, modifier selection and contiguous range selection SHALL remain available. Acquisition changes SHALL require an explicit collection checkbox or a clearly scoped selected-point enable/skip action. Read range descriptions SHALL distinguish item count from register/bit span, expose required types and offer editable name suggestions.

#### Scenario: Inspect and select a range
- **WHEN** the operator clicks points and uses Shift or Ctrl to select multiple points
- **THEN** only selection/detail state changes; collection remains unchanged until an explicit collection action.

#### Scenario: Multiword types are explained
- **WHEN** a range contains four float32 items
- **THEN** the read summary identifies four data items occupying eight words and retains the required type control.


<!-- @trace
source: close-studio-v2-simplification
updated: 2026-10-23
code:
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-verification-20261023.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/schemaOperationRecovery.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023.png
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel.tsx
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-logs-20261023.json.gz
  - docs/releases/studio-v2-change-closeout-20261023.md
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023-mapping.png
  - docs/plans/studio-v2-flow-completion/README.md
tests:
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - internal/datalink/recordingplan/schema_group_extend_test.go
  - frontend/tests/unit/workbench-v2/schemaOperationRecovery.test.ts
  - frontend/tests/unit/workbench-v2/groupSchemaPanel.test.tsx
-->

---
### Requirement: Narrow basic managed membership updates

Basic SHALL present a reviewed member-addition/removal difference for eligible canonical managed groups and save it through existing workspace/group/connector version checks. Existing columns, required flags, deduplication, entity/provenance rules, frozen historical data and custom groups SHALL NOT be overwritten. Customized or unsaved-ineligible mappings SHALL be blocked or directed to existing Advanced controls. Required storage changes SHALL use the existing explicit DDL preview rather than an implicit mutation.

#### Scenario: Add a newly mapped point
- **WHEN** an eligible basic managed group has a newly saved point mapping
- **THEN** Basic offers the narrow member difference and explicit update, preserves existing group policies/columns, and links to confirmed schema preparation when a column is missing.

#### Scenario: Existing customization or stale version
- **WHEN** the group has custom columns/entity/deduplication semantics or the expected version is stale
- **THEN** Basic does not overwrite the configuration; existing advanced editing or version repair remains available.


<!-- @trace
source: close-studio-v2-simplification
updated: 2026-10-23
code:
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-verification-20261023.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/schemaOperationRecovery.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023.png
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel.tsx
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-logs-20261023.json.gz
  - docs/releases/studio-v2-change-closeout-20261023.md
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023-mapping.png
  - docs/plans/studio-v2-flow-completion/README.md
tests:
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - internal/datalink/recordingplan/schema_group_extend_test.go
  - frontend/tests/unit/workbench-v2/schemaOperationRecovery.test.ts
  - frontend/tests/unit/workbench-v2/groupSchemaPanel.test.tsx
-->

---
### Requirement: Recording claims follow committed evidence

Runtime diagnostics with no failure_stage SHALL state only that no error was detected. Acquisition, waiting for the first record and confirmed SQL commit SHALL remain distinguishable. Basic recording success SHALL require the existing workspace/group/applied revision/connector-scoped committed-effect evidence. Share-only success or acquisition alone SHALL NOT claim all points reached a database.

#### Scenario: Collection without SQL evidence
- **WHEN** reads or Share delivery succeed but no matching committed SQL effect is observed
- **THEN** diagnostics do not claim database delivery and Basic continues to show the applicable waiting or evidence state.

#### Scenario: Scoped database commit
- **WHEN** matching committed-effect evidence exists for the effective group and connector
- **THEN** Basic can show confirmed database writing using that evidence without treating missing errors as proof.

<!-- @trace
source: close-studio-v2-simplification
updated: 2026-10-23
code:
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-verification-20261023.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/schemaOperationRecovery.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023.png
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel.tsx
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-logs-20261023.json.gz
  - docs/releases/studio-v2-change-closeout-20261023.md
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023-mapping.png
  - docs/plans/studio-v2-flow-completion/README.md
tests:
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - internal/datalink/recordingplan/schema_group_extend_test.go
  - frontend/tests/unit/workbench-v2/schemaOperationRecovery.test.ts
  - frontend/tests/unit/workbench-v2/groupSchemaPanel.test.tsx
-->

---
### Requirement: Basic missing-value policy reveals required gating
Basic recording SHALL state that required missing or invalid members still block a partial row, show the effective required state, and provide a narrow entry to existing per-member required controls without silently changing any member.

#### Scenario: Required gap remains blocked
- **WHEN** a basic group selects partial and a required member is missing or invalid
- **THEN** the UI SHALL explain that the row is skipped and SHALL expose the existing member settings without changing required automatically.

#### Scenario: Optional gap and recovery
- **WHEN** an explicitly optional member is missing in a partial-capable destination and subsequently recovers
- **THEN** the system SHALL preserve NULL and its reason for the partial row and write its actual value after recovery; skip_row and required protection SHALL remain unchanged.

#### Scenario: Edit existing exact managed members
- **WHEN** an operator opens required settings for saved SQLite managed uint16 or uint32 members in canonical generated TEXT columns
- **THEN** the editor SHALL retain their exact storage compatibility and require fresh schema confirmation after a group revision change without weakening custom-table checks.

#### Scenario: Confirm unchanged schema after member policy edit
- **WHEN** a saved managed group changes required or completeness policy without changing its physical columns
- **THEN** a compatible preview SHALL return an empty statements array, retain its revision-bound token and operation identity, and require explicit confirmation before applying the new group revision.

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