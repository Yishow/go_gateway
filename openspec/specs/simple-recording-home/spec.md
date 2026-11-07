# simple-recording-home Specification

## Purpose

Provide an evidence-led recording home inside Studio V2 that helps operators identify the next safe action while retaining access to all existing recording configurations.

## Requirements

### Requirement: Home guidance projects canonical scoped evidence
The system SHALL derive recording home guidance from the existing workspace, source, group, lifecycle and delivery authorities. It MUST NOT introduce a second persisted configuration authority or mutate configuration through navigation.

#### Scenario: Opening the recording home
- **WHEN** an operator opens /studio/v2?view=recordings
- **THEN** the page reads the current scoped evidence and offers navigation actions without saving, activating, creating a table or restarting a recording.

#### Scenario: A previous scope returns late
- **WHEN** a request for an older workspace or group revision completes after the operator has changed scope
- **THEN** that response does not replace the current home evidence or enable an action in the new scope.




---
### Requirement: Guidance states coexist with deterministic priority
The system SHALL display loading, failed bootstrap and unprovable scope independently from empty or healthy data. For available data, unresolved operations or actionable faults SHALL precede drafts, missing devices, missing points, missing recordings and healthy summaries as the primary guidance, while every other item remains reachable.

#### Scenario: A fresh workspace has no devices
- **WHEN** a successfully loaded workspace has no devices or recording rules
- **THEN** the primary action invites creating the first recording and begins with device connection.

#### Scenario: A device lacks points
- **WHEN** a workspace has a device but no usable source points and no higher-priority issue
- **THEN** guidance opens that device's point configuration rather than requiring its connection to be recreated.

#### Scenario: Existing points are ready for a new recording
- **WHEN** usable points exist without recordings or higher-priority issues
- **THEN** the primary action selects existing points to create a recording.

#### Scenario: Draft and failed recording coexist
- **WHEN** recording A has an unresolved operation, recording B has historical committed evidence and draft C lacks a destination
- **THEN** the primary action checks A's operation, B remains inspectable, and C offers a named continuation at the missing destination section.

#### Scenario: Bootstrap fails
- **WHEN** workspace evidence cannot be loaded
- **THEN** the page shows an actionable unavailable state and never renders the workspace as empty or healthy.




---
### Requirement: Healthy summaries distinguish history from present evidence
The system SHALL label saving, start acknowledgement, collection, SQL commit and verified readback distinctly. A historical SQL commit SHALL NOT establish present freshness, and general delivery evidence SHALL NOT be presented as raw row values unless the existing readback contract supplies them.

#### Scenario: A previously successful device disconnects
- **WHEN** a recording has a historical commit at 10:32 but current source evidence is stale or disconnected
- **THEN** the home retains the last successful time and destination while exposing the current fault and a scoped repair action.

#### Scenario: First start has no SQL receipt
- **WHEN** start has been acknowledged but no matching commit evidence exists
- **THEN** the home states that it is awaiting the first write rather than claiming successful recording.




---
### Requirement: Optional home preserves existing Studio navigation
The system SHALL expose the recording home as an explicit view of /studio/v2 without changing the default entry in this change. Valid entity locators and runtime return intents SHALL retain precedence and scope. The system SHALL validate raw query encoding and uniqueness before removing only a single recognized presentation hint for the existing strict locator parser; invalid locator input SHALL remain invalid and MUST NOT fall back to the general home. Existing runtime, logs, test and advanced paths SHALL remain accessible.

#### Scenario: A deep link includes a recording home hint
- **WHEN** a valid entity locator or runtime repair focus is accompanied by view=recordings
- **THEN** the exact existing locator or repair intent is honored without losing its query, hash or scope.

#### Scenario: An unknown view is supplied
- **WHEN** the URL contains an unsupported view value without any locator intent
- **THEN** the application uses the existing entry behavior without mutating data or restoring the retired /studio surface.

#### Scenario: Presentation hints do not weaken locator validation
- **WHEN** a locator request contains malformed encoding, duplicate view values, an unknown presentation value or an invalid locator key
- **THEN** it remains invalid with no fallback to an unscoped home, while a valid locator with one recognized view retains its original scoped destination.




---
### Requirement: Theme preference resolves without inventing persisted intent
The system SHALL use one existing ThemeProvider with separate system/light/dark preference and resolved light/dark theme. The existing server theme SHALL remain global rather than introducing per-workspace preferences. Initial resolution SHALL prefer a valid actually stored server theme key, then a versioned local preference, then a valid legacy light/dark value, then system. Mapper-injected defaults MUST NOT count as stored intent. Legacy values SHALL be preserved as compatibility settings without claiming they prove manual choice. The auto setting SHALL map to system.

#### Scenario: No stored preference exists
- **WHEN** raw server settings have no theme key and no valid local preference exists
- **THEN** the provider follows the current OS theme, subscribes to subsequent OS changes, and does not persist the resolved color as an override.

#### Scenario: Explicit override and return to system
- **WHEN** an operator selects light while the OS is dark, then later selects system
- **THEN** OS changes do not override the explicit light choice, and selecting system immediately follows the current OS and future changes.

#### Scenario: Legacy value and mapper default coexist
- **WHEN** legacy local storage contains light, the raw server theme key is absent, and the settings mapper supplies dark
- **THEN** light remains the compatibility choice with an available system option rather than treating the mapper default as a stored server preference.

#### Scenario: Persistence or browser APIs are unavailable
- **WHEN** storage, matchMedia or a theme-save operation fails
- **THEN** the interface remains usable, unavailable OS resolution falls back to light without saving an override, and save failure is not represented as persisted success.

#### Scenario: An old hydration response arrives
- **WHEN** an operator changes preference or backend/app context before an earlier settings request completes
- **THEN** stale hydration does not overwrite the new choice or import another backend context's preference, and no recording state or device operation is changed.

#### Scenario: Theme-only choice does not flush business settings
- **WHEN** an operator changes theme while scheduler or Modbus Share edits remain unsaved
- **THEN** only updateKey for theme may persist, other setting keys remain unwritten, workspace navigation does not reset the global preference, and OS-only resolution changes perform no persistence at all.




---
### Requirement: Semantic design tokens govern affected surfaces
The new recording surfaces and reused affected controls SHALL implement the root DESIGN.md primitive, semantic light/dark and component-state token contract through one CSS variable authority. Typography, spacing, radii, elevation, focus, contrast and reduced motion SHALL be consistent across pages and portals. Raw color literals MUST NOT be duplicated in new component styling. Unaffected legacy pages MUST NOT require a whole-application redesign.

#### Scenario: Theme changes while a modal is open
- **WHEN** the resolved theme changes with an unsaved configuration modal or pending operation visible
- **THEN** the page, portal, native controls and status messages use the same theme while preserving field values, focus and operation identity without recording, device or database-schema commands. Explicit preference selection may only use the existing theme-key update mutation, not a full settings snapshot; OS resolution changes do not save a preference.

#### Scenario: Information must remain readable
- **WHEN** either theme renders secondary text, placeholders, disabled explanations, status or keyboard focus
- **THEN** informative text meets 4.5:1 contrast and required controls/focus meet 3:1 against their actual adjacent backgrounds, status has a non-color label, and reduced motion removes nonessential animation.




---
### Requirement: New configuration popups are centered viewport-sized modals
New simple-recording configuration popups SHALL use the modal size classifications in DESIGN.md: centered desktop form popups, approved full-viewport mobile and tablet presentation, and smaller auto-height desktop confirmations. The same responsive modal contract SHALL apply to existing modal primitives; non-modal informational popovers MUST remain excluded. The two main pages, including the recording editor, SHALL remain independently responsive and MUST NOT be converted into modals merely because they contain configuration.

#### Scenario: Opening a configuration popup
- **WHEN** an operator opens a device, point or supplemental configuration popup from the simple interface
- **THEN** the modal uses the appropriate DESIGN.md size classification, preserves the existing unsaved-change and focus-return contracts, and does not resize the main page to 75 percent.

---
### Requirement: Classified modals and contextual return state preserve safety
Normal desktop form popups SHALL retain their existing 75vw by 75dvh classification, with approved source and memory exceptions in DESIGN.md. At viewport widths up to 1024 CSS pixels, or with a coarse pointer up to 1366 CSS pixels, both form and confirmation modals SHALL use the full visible viewport with safe areas inside the modal. Low-height viewport handling SHALL keep controls reachable. Desktop short confirmations without complex forms SHALL retain the smaller auto-height classification. Modal headers and action footers SHALL remain visible while the body scrolls. Contextual help MUST NOT hide required fields, validation failures, disabled reasons, schema/history impacts, unresolved risks or consent. Return navigation SHALL preserve search, filters and scroll only within the same workspace, entity and view context and SHALL revalidate identity when revisions change.

#### Scenario: A mobile or tablet viewport opens a modal
- **WHEN** any classified form, confirmation or existing modal primitive opens on a mobile or tablet viewport
- **THEN** it uses the shared full-viewport presentation while keeping the title, actions, focused field and validation errors reachable
- **AND** orientation changes preserve the existing dirty guard, operation state and focus context.

#### Scenario: A short confirmation is displayed
- **WHEN** stop, remove or discard requires a short confirmation
- **THEN** the dialog identifies the object and consequence, retains explicit consent and does not substitute for a detailed schema-impact review
- **AND** desktop retains the smaller centered classification while mobile and tablet use the shared full-viewport presentation.

#### Scenario: Returning after an item was removed
- **WHEN** an operator returns to the same filtered list after its entity or revision changed
- **THEN** search, filter and position restore only within that context, unavailable selections are cleared with a reason, focus moves to a safe target and stale responses cannot restore old authority or action availability.

#### Scenario: Help is collapsed
- **WHEN** optional contextual help is closed
- **THEN** critical risk, errors, disabled reasons and consent remain visible and the help control remains keyboard accessible.

---
### Requirement: Development style preview is isolated from operator navigation
A development/test-only style preview SHALL reuse the actual components and semantic tokens to exercise light/dark, inputs, actions, statuses, errors, disabled, focus and modal states. It MUST NOT appear in the operator primary menu or become a third primary workflow. The existing build-mode guard SHALL exclude this development entry from production bundles. The preview MUST NOT write configuration or data, activate resources, collect credentials or fabricate production bootstrap. A mock demo MAY expose a dedicated non-primary link or hash and SHALL label its fixtures as simulated.

#### Scenario: Preview is inspected in development
- **WHEN** a developer opens the dedicated development preview
- **THEN** real shared components and tokens render the selected states without backend mutations, and ordinary operator navigation has no preview menu item.

#### Scenario: Production build is generated
- **WHEN** the production frontend bundle is built
- **THEN** the development preview entry is excluded and cannot become a reachable third operator page, while the actual recording surfaces retain their shared component styling.




---
### Requirement: Large recording collections remain searchable and bounded
The home SHALL provide compact management lists with keyword search across recording, device and configured table labels, device/status/destination filters, deterministic ordering with canonical identity tie-breaks and 25-or-50-row pagination as design limits. The device summary SHALL show at most four rows with a same-page management view, not a third primary page. Attention evidence SHALL remain available independently of list filters. Client-side projection over loaded canonical data MUST NOT be described as an existing bounded cross-device backend search service.

#### Scenario: Many devices and recordings are loaded
- **WHEN** the available dataset exceeds one page
- **THEN** only the chosen page renders in the complete management list, the device summary remains bounded, and search/filter can locate a named recording without traversing an unbounded card wall.

#### Scenario: A filtered return targets a removed item
- **WHEN** an operator returns from editing and the matching scope revision has removed the selected item or reduced the number of pages
- **THEN** filters and sort restore only in the same context, selection is revalidated, the page clamps to a valid page, and focus/scroll return safely without stale authority.

#### Scenario: Search has no matches
- **WHEN** valid loaded data exists but the current search or filters match no entries
- **THEN** the page shows a no-results state with a clear-filter action rather than claiming the workspace is empty or healthy.

---
### Requirement: Compact home source observations preserve canonical evidence

The recording home SHALL use a workspace-scoped read-only compact source observation derived from the existing canonical source authorities. The endpoint MUST NOT create, repair, reconcile or activate configuration. It SHALL compare complete source projections twice and return available only when identity, membership and content checks agree. Its digest MUST NOT serve as a command or readiness token or imply an atomic snapshot.

#### Scenario: Complete unchanged sources

- **WHEN** GET /api/v1/datalink/studio-v2/recording-home/sources receives an existing workspace_id and both complete projections agree
- **THEN** it returns that workspace identity, an opaque source digest and the same device/candidate/blocker projection as the existing home logic
- **AND** it excludes credentials, connection configuration, live point values and raw backend errors

#### Scenario: Inconsistent or unavailable sources

- **WHEN** source content changes between passes, duplicate identities occur, members are missing or foreign, or a source read fails
- **THEN** the endpoint returns unknown or an actionable read error instead of a partial available summary
- **AND** navigation performs no mutation

##### Example: mismatched mapping owner

- **GIVEN** workspace W1 owns device D1 while mapping M1 claims point P1 on foreign device D2
- **WHEN** W1 source observations are requested
- **THEN** M1 does not become an available candidate and the observation is unknown

---
### Requirement: Return navigation reuses observed presentation without reusing authority

The recording home SHALL preserve its same-workspace filter, page, scroll and focus while revalidating a cached observation on return. Cached evidence SHALL be labeled observed or stale until the current visit validates it. Every visit SHALL retain a separate cancellation and late-response fence. Cache entries SHALL be isolated by QueryClient/backend context and current API origin, retain only the latest observation per workspace, and preserve the existing eight-workspace bound and oldest-scope eviction.

#### Scenario: Same workspace returns

- **WHEN** a user returns to a previously observed workspace
- **THEN** the cached presentation is available with a refreshing indicator and one new source-summary request
- **AND** pending or failed validation cannot label the sources fresh or enable commands requiring fresh proof

#### Scenario: Scope changes and old replies arrive

- **WHEN** navigation follows workspace A, workspace B, then workspace A and an earlier request finishes late
- **THEN** only the latest matching visit updates displayed evidence
- **AND** B evidence and older A evidence cannot overwrite the current A observation

#### Scenario: Refresh fails or sources mutate

- **WHEN** revalidation fails or an affected source mutation invalidates the observation
- **THEN** the last observed view is explicitly stale or failed, with read-only retry
- **AND** fresh locator, CAS, revision and readiness checks remain required by existing action owners

---
### Requirement: Home data demand stays bounded and complete

One home source revalidation SHALL use one compact summary request rather than the previous ten raw source requests. Existing bounded visible-page detail demand SHALL remain bounded, and the summary MUST NOT infer completeness from a truncated source list.

#### Scenario: Large source fixture is projected

- **WHEN** the same synthetic 20-device and 2000-point fixture is read before and after migration
- **THEN** compact source bytes are less than the original 3.06 MB replay and canonical candidates and excluded are identical
- **AND** unrelated groups do not receive new lifecycle or delivery reads merely because sources were summarized

#### Scenario: Source transport migration

- **WHEN** the new endpoint and projection parity tests pass and the home is migrated
- **THEN** the home does not run both old and new source transports
- **AND** full source APIs remain available for existing editor flows

#### Scenario: Cache context isolation and eviction

- **WHEN** two backend contexts contain the same workspace identity, or a context observes a ninth workspace
- **THEN** neither backend receives the other's cached evidence, and the oldest scope is evicted while each retained workspace has only its latest observation
- **AND** invalidation tests cover successful device, rule, point and mapping mutations without changing their save barriers

---
### Requirement: Home shows live device connection state

The recording home SHALL derive one connection state for every workspace device from a single workspace-wide `GET /api/v1/datalink/runtime/status` read, polled every 5 seconds only while the page is visible. The state SHALL be the first matching rule of: `unavailable` (no successful read within the last 15 seconds, `snapshot_state.unavailable` is true, or no collector exists for the device), `disconnected` (`breaker_state` is not `closed` or `status` is `error`), `idle` (`running` is false, availability is `unavailable`, or `points_total` is 0), `degraded` (stale or error points exist, healthy points differ from total points, `status` is `warning`, `last_error` is present, or `last_read_at` is missing or older than the maximum of 3 × `effective_interval_ms` and 5000 ms), otherwise `collecting`. Each state SHALL be shown with text, an icon and a semantic color, together with healthy/total point counts and the relative last-read time. After a failed read the last successful states MAY remain visible for at most 15 seconds marked as outdated. The home MUST NOT show `collecting` when the status read is unavailable, and MUST NOT display raw backend errors.

#### Scenario: Device states are derived from collectors

- **WHEN** the runtime status read succeeds
- **THEN** each device shows the state selected by the ordered rules and its point and freshness facts

##### Example: state derivation at now = 12:00:10, interval 1000 ms

| Collector facts | Shown state |
| --- | --- |
| running, 48/48 healthy, last read 12:00:09, breaker closed | collecting |
| running, 46/48 healthy with 2 stale, status warning | degraded |
| running, 48/48 healthy, last read 12:00:04 (6 s old, threshold 5 s) | degraded |
| running, 48/48 healthy, last read 12:00:05 (exactly 5 s old) | collecting |
| running, 47/48 healthy with 0 stale and 0 error points | degraded |
| running false, status idle | idle |
| breaker open | disconnected |
| no collector for the device | unavailable |

#### Scenario: Status read fails

- **WHEN** the runtime status read fails after a previous success
- **THEN** the last observed states remain visible for at most 15 seconds marked as outdated, after which every device shows `unavailable`

#### Scenario: Page is hidden

- **WHEN** the browser tab is hidden or the operator leaves the home
- **THEN** status polling stops and resumes with one immediate read when the home becomes visible again


<!-- @trace
source: improve-recording-home-status-and-delete
updated: 2026-11-08
code:
  - internal/datalink/dbtarget/writer_snapshot.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupLifecycleBar.tsx
  - frontend/studio-v2-diagnostic.cjs
  - internal/datalink/workspace/write_group_service.go
  - internal/datalink/runtimeview/contract.go
  - docs/validation/memory-views-20261106/skill/SKILL.md
  - frontend/src/services/runtimeLogLegacyProjection.ts
  - internal/datalink/dbtarget/writer_workers.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeResultFacts.tsx
  - frontend/src/i18n/locales/zh-TW/recording-home-status.json
  - internal/datalink/workspace/write_group_permanent_delete_audit.go
  - docs/testing.md
  - frontend/src/features/datalink/workbench-v2/recording-home/DeviceDeleteDialog.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/RecordingLifecycleHistory.tsx
  - internal/datalink/workspace/device_recording_delete.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeStatusSummary.tsx
  - internal/datalink/workspace/write_group_permanent_delete.go
  - internal/diagnostics/projection.go
  - internal/api/router.go
  - frontend/src/services/runtimeLogDomainProjection.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/recording-delete-dialog.css
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - internal/datalink/modbusshare/service_projection.go
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - frontend/src/i18n/locales/zh-TW/runtime-logs.json
  - internal/api/router_studio_v2_write_group_routes.go
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeStatusCounts.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/SavedRecordingActions.tsx
  - internal/datalink/workspace/write_group_permanent_delete_purge.go
  - internal/diagnostics/domain_event.go
  - docs/swagger/swagger.json
  - docs/swagger/docs.go
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeLifecycleMatches.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeWriteStages.tsx
  - internal/api/handlers/typed_errors.go
  - internal/datalink/device/workspace_delete.go
  - internal/datalink/groupdelivery/quota.go
  - frontend/src/services/runtimeRuleEvidenceParser.ts
  - AGENTS.md
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - internal/datalink/workspace/write_group_lifecycle_actions.go
  - frontend/src/types/runtimeRuleEvidence.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeStatusDistribution.tsx
  - internal/api/handlers/studio_v2_workspace_write_groups_handler.go
  - internal/datalink/workspace/write_group_permanent_delete_recording_start.go
  - frontend/src/i18n/locales/en/recording-home.json
  - frontend/src/services/studioV2WorkspaceWriteGroups.ts
  - internal/datalink/workspace/write_group_permanent_delete_shared.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_delete.go
  - internal/datalink/audit/types.go
  - internal/datalink/dbtarget/writer_failure_log.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/services/studioV2WorkspaceDelete.ts
  - internal/api/handlers/studio_v2_workspace_write_groups_handler_mutate.go
  - internal/datalink/dbtarget/writer_flush.go
  - frontend/src/hooks/datalink/useStudioV2WriteGroups.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeRoute.tsx
  - internal/testutil/legacyrecording/tombstone.go
  - internal/api/handlers/studio_v2_workspace_write_groups_handler_delete.go
  - frontend/src/features/datalink/workbench-v2/recording-home/deviceConnectionState.ts
  - frontend/src/i18n/locales/en/recording-home-status.json
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeRowActions.tsx
  - frontend/src/i18n/locales/zh-TW/recording-home-delete.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupDeliveryStrip.tsx
  - frontend/src/i18n/locales/zh-TW/recording-home.json
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHome.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeDestinationCell.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeWriteStageMapping.ts
  - internal/datalink/dbtarget/writer_group_flush.go
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeRows.css
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/datalink/dbtarget/writer_retry.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeUnconfirmed.tsx
  - frontend/src/features/runtime-logs/RuntimeLogLegacyEvent.tsx
  - internal/datalink/workspace/write_group_permanent_delete_operations.go
  - internal/datalink/workspace/write_group_permanent_delete_storage.go
  - frontend/src/i18n/config.ts
  - internal/datalink/workspace/write_group_permanent_delete_scope.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingDeleteDialog.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeStatus.css
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeDeviceActions.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/useRecordingHomeDeviceStatus.ts
  - docs/swagger/swagger.yaml
  - frontend/src/i18n/locales/en/recording-home-rows.json
  - frontend/src/types/studioV2Delete.ts
  - frontend/src/types/deviceConnection.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeRecordRow.tsx
  - internal/datalink/workspace/write_group_permanent_delete_guard.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeDeviceStrip.tsx
  - frontend/src/i18n/locales/zh-TW/recording-home-rows.json
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/dev/previewI18n.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/SavedRecordingRemoveDialog.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - docs/validation/memory-views-20261106/skill/references/lessons.md
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingDeleteDialogController.tsx
  - internal/datalink/grouppipeline/backlog_guard.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeCollection.tsx
  - frontend/src/i18n/locales/en/recording-home-delete.json
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingDeleteDialogError.tsx
  - frontend/src/i18n/locales/en/runtime-logs.json
tests:
  - internal/api/router_studio_v2_write_group_lifecycle_test.go
  - internal/datalink/workspace/write_group_permanent_delete_delivery_fixture_test.go
  - internal/api/router_studio_v2_successor_first_start_fixture_test.go
  - frontend/tests/unit/runtime-logs/legacy-event-fixtures.ts
  - frontend/tests/unit/workbench-v2/saved-recording-actions.test.tsx
  - internal/api/handlers/studio_v2_workspace_write_groups_handler_delete_test.go
  - frontend/tests/unit/workbench-v2/groupDeliveryStrip.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-accessibility.test.tsx
  - internal/api/router_studio_v2_locator_entry_test.go
  - frontend/tests/unit/dev/recording-style-bootstrap.test.tsx
  - internal/datalink/workspace/recording_start_disabled_test.go
  - internal/datalink/dbtarget/writer_failure_revision_test.go
  - cmd/test_ui/group_pipeline_outage_delete_test.go
  - frontend/tests/unit/workbench-v2/recording-editor-range-mc-prepare.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-memory-scope.test.tsx
  - frontend/tests/unit/workbench-v2/recording-guided-saved-card.test.tsx
  - internal/datalink/workspace/write_group_lifecycle_test.go
  - internal/datalink/workspace/write_group_mixed_safety_test.go
  - frontend/tests/unit/workbench-v2/recording-home-route.test.tsx
  - internal/datalink/workspace/write_group_permanent_delete_storage_test.go
  - frontend/tests/unit/workbench-v2/device-connection-state.test.ts
  - frontend/tests/unit/workbench-v2/recording-home-status-counts.test.ts
  - frontend/tests/unit/workbench-v2/recording-home-shell.test.tsx
  - frontend/tests/unit/workbench-v2/studio-locator-group-focus.test.tsx
  - cmd/test_ui/gateway_logging.go
  - internal/api/router_studio_v2_write_group_boundaries_test.go
  - frontend/tests/unit/workbench-v2/recording-home-source-modal.test.tsx
  - internal/datalink/grouppipeline/backlog_guard_lifecycle_test.go
  - internal/datalink/workspace/write_group_circuit_profile_rules_test.go
  - frontend/tests/unit/workbench-v2/recording-delete-dialog.test.tsx
  - internal/datalink/workspace/write_group_permanent_delete_test.go
  - internal/api/router_studio_v2_write_group_action_errors_test.go
  - frontend/tests/unit/runtime-logs/runtimeLogLegacyProjection.test.ts
  - cmd/test_ui/delivery_events_unavailable_test.go
  - frontend/tests/e2e/fixtures/recording-lifecycle-api.ts
  - frontend/tests/unit/runtime-logs/RuntimeLogCompatibility.test.tsx
  - internal/datalink/dbtarget/writer_failure_log_test.go
  - internal/api/router_runtime_rule_evidence_test.go
  - internal/api/router_studio_v2_write_group_stop_remove_safety_test.go
  - internal/datalink/groupdelivery/quota_reclaim_receipts_test.go
  - internal/datalink/workspace/write_group_basic_managed_test.go
  - frontend/tests/unit/workbench-v2/recording-modal-return.test.tsx
  - frontend/tests/unit/services/studioV2WorkspaceWriteGroups.test.ts
  - frontend/tests/unit/workbench-v2/recording-home-source-switching.test.tsx
  - cmd/test_ui/group_pipeline_outage_test.go
  - frontend/tests/unit/workbench-v2/recording-home-row-actions.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-legacy-delivery-safety.test.tsx
  - frontend/tests/unit/workbench-v2/device-delete-cache.test.tsx
  - internal/datalink/workspace/write_group_lifecycle_error_preservation_test.go
  - internal/datalink/workspace/write_group_reclaim_delete_test.go
  - cmd/test_ui/group_pipeline_lifecycle_crash_test.go
  - frontend/tests/unit/workbench-v2/recording-home-visual-status.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-device-status.test.tsx
  - internal/datalink/workspace/write_group_permanent_delete_legacy_test.go
  - frontend/tests/unit/services/studioV2DeleteResponses.test.ts
  - internal/api/router_studio_v2_recording_start_delete_test.go
  - frontend/tests/unit/workbench-v2/recording-home-memory-view.test.tsx
  - frontend/tests/unit/dev/recording-style-preview.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-acquisition-anomalies.test.tsx
  - internal/datalink/workspace/write_group_successor_replay_test.go
  - internal/datalink/workspace/write_group_lifecycle_actions_fixture_test.go
  - internal/datalink/workspace/write_group_lifecycle_actions_test.go
  - frontend/tests/e2e/recording-home-status-layout.spec.ts
  - frontend/tests/e2e/recording-home-tombstone-delete.spec.ts
  - frontend/tests/unit/workbench-v2/recording-home-inline-evidence.test.tsx
  - frontend/tests/unit/workbench-v2/recordingLifecycleBrowserFixture.test.ts
  - internal/datalink/workspace/write_group_service_test.go
  - internal/datalink/workspace/write_group_tombstone_remove_actions_test.go
  - cmd/test_ui/service_wiring.go
  - cmd/test_ui/group_pipeline_legacy_takeover_test.go
  - frontend/tests/unit/runtime-logs/runtimeLogDomainProjection.test.ts
  - frontend/tests/unit/workbench-v2/recording-home-delete-integration.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-device-actions.ts
  - frontend/tests/unit/workbench-v2/recording-home-mc-resave.test.tsx
  - frontend/tests/unit/workbench-v2/writeGroupLifecycle.test.tsx
  - internal/datalink/grouppipeline/backlog_guard_delete_count_test.go
  - cmd/test_ui/group_pipeline_lifecycle_test.go
  - internal/api/router_studio_v2_workspace_devices_test.go
  - internal/datalink/grouppipeline/applied_configuration_read_test.go
  - internal/datalink/workspace/write_group_permanent_delete_recording_start_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_cascade_test.go
  - internal/api/router_studio_v2_recording_delete_contract_test.go
  - frontend/tests/unit/workbench-v2/device-delete-dialog.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-management.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-mc-locator.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-source-entry.test.tsx
  - frontend/tests/e2e/recording-database-unavailable.spec.ts
  - internal/datalink/workspace/write_group_successor_guard_test.go
  - internal/datalink/dbtarget/writer_failure_log_outcome_test.go
  - internal/api/router_studio_v2_write_group_lifecycle_history_test.go
  - internal/datalink/grouppipeline/stop_remove_intake_fence_test.go
  - cmd/test_ui/gateway_logging_database_test.go
  - frontend/tests/e2e/recording-home-narrow-records.spec.ts
  - frontend/tests/unit/workbench-v2/recording-home-tombstone-delete.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler_recordings_test.go
  - internal/datalink/workspace/write_group_lifecycle_error_test.go
  - internal/datalink/workspace/write_group_permanent_delete_schema_test.go
  - internal/datalink/workspace/write_group_permanent_delete_shared_test.go
  - frontend/tests/unit/workbench-v2/recording-home-new-range-strict-ack.test.tsx
-->

---
### Requirement: Home header summarizes devices, recordings and runtime writes

The recording home header SHALL show the device total with counts per connection state, the recording total with counts per saved status computed from all loaded recordings independent of filters, and `metrics.write_success_total` and `metrics.write_error_total` labelled as runtime cumulative counts that MUST NOT be presented as SQL-confirmed writes. Below the summary a device status strip SHALL list every device with its name, protocol, connection state and point facts; activating a device item SHALL apply the existing device filter to the recording list and activating it again SHALL clear that filter. The separate bottom device summary SHALL be removed from the recordings view, and the devices view SHALL show the same connection state per device.

#### Scenario: Summary counts match loaded data

- **WHEN** the home has loaded 3 devices and 4 recordings
- **THEN** the summary counts equal the per-state and per-status counts of those loaded items regardless of active filters

##### Example: mixed workspace

- **GIVEN** devices MC (degraded), mc_new (idle), M1 (idle) and recordings with statuses draft, draft, draft, removed
- **WHEN** the home renders
- **THEN** the summary reads devices 3 (degraded 1, idle 2) and recordings 4 (draft 3, removed 1)

#### Scenario: Device strip filters recordings

- **WHEN** the operator activates device MC in the status strip
- **THEN** only recordings with a member on MC remain listed and the device filter control shows MC
- **AND** activating MC again restores the unfiltered list

#### Scenario: Snapshot is stale

- **WHEN** `snapshot_state.stale` is true
- **THEN** device states remain shown and the summary displays an outdated-status notice


<!-- @trace
source: improve-recording-home-status-and-delete
updated: 2026-11-08
code:
  - internal/datalink/dbtarget/writer_snapshot.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupLifecycleBar.tsx
  - frontend/studio-v2-diagnostic.cjs
  - internal/datalink/workspace/write_group_service.go
  - internal/datalink/runtimeview/contract.go
  - docs/validation/memory-views-20261106/skill/SKILL.md
  - frontend/src/services/runtimeLogLegacyProjection.ts
  - internal/datalink/dbtarget/writer_workers.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeResultFacts.tsx
  - frontend/src/i18n/locales/zh-TW/recording-home-status.json
  - internal/datalink/workspace/write_group_permanent_delete_audit.go
  - docs/testing.md
  - frontend/src/features/datalink/workbench-v2/recording-home/DeviceDeleteDialog.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/RecordingLifecycleHistory.tsx
  - internal/datalink/workspace/device_recording_delete.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeStatusSummary.tsx
  - internal/datalink/workspace/write_group_permanent_delete.go
  - internal/diagnostics/projection.go
  - internal/api/router.go
  - frontend/src/services/runtimeLogDomainProjection.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/recording-delete-dialog.css
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - internal/datalink/modbusshare/service_projection.go
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - frontend/src/i18n/locales/zh-TW/runtime-logs.json
  - internal/api/router_studio_v2_write_group_routes.go
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeStatusCounts.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/SavedRecordingActions.tsx
  - internal/datalink/workspace/write_group_permanent_delete_purge.go
  - internal/diagnostics/domain_event.go
  - docs/swagger/swagger.json
  - docs/swagger/docs.go
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeLifecycleMatches.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeWriteStages.tsx
  - internal/api/handlers/typed_errors.go
  - internal/datalink/device/workspace_delete.go
  - internal/datalink/groupdelivery/quota.go
  - frontend/src/services/runtimeRuleEvidenceParser.ts
  - AGENTS.md
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - internal/datalink/workspace/write_group_lifecycle_actions.go
  - frontend/src/types/runtimeRuleEvidence.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeStatusDistribution.tsx
  - internal/api/handlers/studio_v2_workspace_write_groups_handler.go
  - internal/datalink/workspace/write_group_permanent_delete_recording_start.go
  - frontend/src/i18n/locales/en/recording-home.json
  - frontend/src/services/studioV2WorkspaceWriteGroups.ts
  - internal/datalink/workspace/write_group_permanent_delete_shared.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_delete.go
  - internal/datalink/audit/types.go
  - internal/datalink/dbtarget/writer_failure_log.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/services/studioV2WorkspaceDelete.ts
  - internal/api/handlers/studio_v2_workspace_write_groups_handler_mutate.go
  - internal/datalink/dbtarget/writer_flush.go
  - frontend/src/hooks/datalink/useStudioV2WriteGroups.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeRoute.tsx
  - internal/testutil/legacyrecording/tombstone.go
  - internal/api/handlers/studio_v2_workspace_write_groups_handler_delete.go
  - frontend/src/features/datalink/workbench-v2/recording-home/deviceConnectionState.ts
  - frontend/src/i18n/locales/en/recording-home-status.json
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeRowActions.tsx
  - frontend/src/i18n/locales/zh-TW/recording-home-delete.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupDeliveryStrip.tsx
  - frontend/src/i18n/locales/zh-TW/recording-home.json
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHome.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeDestinationCell.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeWriteStageMapping.ts
  - internal/datalink/dbtarget/writer_group_flush.go
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeRows.css
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/datalink/dbtarget/writer_retry.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeUnconfirmed.tsx
  - frontend/src/features/runtime-logs/RuntimeLogLegacyEvent.tsx
  - internal/datalink/workspace/write_group_permanent_delete_operations.go
  - internal/datalink/workspace/write_group_permanent_delete_storage.go
  - frontend/src/i18n/config.ts
  - internal/datalink/workspace/write_group_permanent_delete_scope.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingDeleteDialog.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeStatus.css
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeDeviceActions.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/useRecordingHomeDeviceStatus.ts
  - docs/swagger/swagger.yaml
  - frontend/src/i18n/locales/en/recording-home-rows.json
  - frontend/src/types/studioV2Delete.ts
  - frontend/src/types/deviceConnection.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeRecordRow.tsx
  - internal/datalink/workspace/write_group_permanent_delete_guard.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeDeviceStrip.tsx
  - frontend/src/i18n/locales/zh-TW/recording-home-rows.json
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/dev/previewI18n.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/SavedRecordingRemoveDialog.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - docs/validation/memory-views-20261106/skill/references/lessons.md
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingDeleteDialogController.tsx
  - internal/datalink/grouppipeline/backlog_guard.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeCollection.tsx
  - frontend/src/i18n/locales/en/recording-home-delete.json
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingDeleteDialogError.tsx
  - frontend/src/i18n/locales/en/runtime-logs.json
tests:
  - internal/api/router_studio_v2_write_group_lifecycle_test.go
  - internal/datalink/workspace/write_group_permanent_delete_delivery_fixture_test.go
  - internal/api/router_studio_v2_successor_first_start_fixture_test.go
  - frontend/tests/unit/runtime-logs/legacy-event-fixtures.ts
  - frontend/tests/unit/workbench-v2/saved-recording-actions.test.tsx
  - internal/api/handlers/studio_v2_workspace_write_groups_handler_delete_test.go
  - frontend/tests/unit/workbench-v2/groupDeliveryStrip.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-accessibility.test.tsx
  - internal/api/router_studio_v2_locator_entry_test.go
  - frontend/tests/unit/dev/recording-style-bootstrap.test.tsx
  - internal/datalink/workspace/recording_start_disabled_test.go
  - internal/datalink/dbtarget/writer_failure_revision_test.go
  - cmd/test_ui/group_pipeline_outage_delete_test.go
  - frontend/tests/unit/workbench-v2/recording-editor-range-mc-prepare.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-memory-scope.test.tsx
  - frontend/tests/unit/workbench-v2/recording-guided-saved-card.test.tsx
  - internal/datalink/workspace/write_group_lifecycle_test.go
  - internal/datalink/workspace/write_group_mixed_safety_test.go
  - frontend/tests/unit/workbench-v2/recording-home-route.test.tsx
  - internal/datalink/workspace/write_group_permanent_delete_storage_test.go
  - frontend/tests/unit/workbench-v2/device-connection-state.test.ts
  - frontend/tests/unit/workbench-v2/recording-home-status-counts.test.ts
  - frontend/tests/unit/workbench-v2/recording-home-shell.test.tsx
  - frontend/tests/unit/workbench-v2/studio-locator-group-focus.test.tsx
  - cmd/test_ui/gateway_logging.go
  - internal/api/router_studio_v2_write_group_boundaries_test.go
  - frontend/tests/unit/workbench-v2/recording-home-source-modal.test.tsx
  - internal/datalink/grouppipeline/backlog_guard_lifecycle_test.go
  - internal/datalink/workspace/write_group_circuit_profile_rules_test.go
  - frontend/tests/unit/workbench-v2/recording-delete-dialog.test.tsx
  - internal/datalink/workspace/write_group_permanent_delete_test.go
  - internal/api/router_studio_v2_write_group_action_errors_test.go
  - frontend/tests/unit/runtime-logs/runtimeLogLegacyProjection.test.ts
  - cmd/test_ui/delivery_events_unavailable_test.go
  - frontend/tests/e2e/fixtures/recording-lifecycle-api.ts
  - frontend/tests/unit/runtime-logs/RuntimeLogCompatibility.test.tsx
  - internal/datalink/dbtarget/writer_failure_log_test.go
  - internal/api/router_runtime_rule_evidence_test.go
  - internal/api/router_studio_v2_write_group_stop_remove_safety_test.go
  - internal/datalink/groupdelivery/quota_reclaim_receipts_test.go
  - internal/datalink/workspace/write_group_basic_managed_test.go
  - frontend/tests/unit/workbench-v2/recording-modal-return.test.tsx
  - frontend/tests/unit/services/studioV2WorkspaceWriteGroups.test.ts
  - frontend/tests/unit/workbench-v2/recording-home-source-switching.test.tsx
  - cmd/test_ui/group_pipeline_outage_test.go
  - frontend/tests/unit/workbench-v2/recording-home-row-actions.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-legacy-delivery-safety.test.tsx
  - frontend/tests/unit/workbench-v2/device-delete-cache.test.tsx
  - internal/datalink/workspace/write_group_lifecycle_error_preservation_test.go
  - internal/datalink/workspace/write_group_reclaim_delete_test.go
  - cmd/test_ui/group_pipeline_lifecycle_crash_test.go
  - frontend/tests/unit/workbench-v2/recording-home-visual-status.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-device-status.test.tsx
  - internal/datalink/workspace/write_group_permanent_delete_legacy_test.go
  - frontend/tests/unit/services/studioV2DeleteResponses.test.ts
  - internal/api/router_studio_v2_recording_start_delete_test.go
  - frontend/tests/unit/workbench-v2/recording-home-memory-view.test.tsx
  - frontend/tests/unit/dev/recording-style-preview.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-acquisition-anomalies.test.tsx
  - internal/datalink/workspace/write_group_successor_replay_test.go
  - internal/datalink/workspace/write_group_lifecycle_actions_fixture_test.go
  - internal/datalink/workspace/write_group_lifecycle_actions_test.go
  - frontend/tests/e2e/recording-home-status-layout.spec.ts
  - frontend/tests/e2e/recording-home-tombstone-delete.spec.ts
  - frontend/tests/unit/workbench-v2/recording-home-inline-evidence.test.tsx
  - frontend/tests/unit/workbench-v2/recordingLifecycleBrowserFixture.test.ts
  - internal/datalink/workspace/write_group_service_test.go
  - internal/datalink/workspace/write_group_tombstone_remove_actions_test.go
  - cmd/test_ui/service_wiring.go
  - cmd/test_ui/group_pipeline_legacy_takeover_test.go
  - frontend/tests/unit/runtime-logs/runtimeLogDomainProjection.test.ts
  - frontend/tests/unit/workbench-v2/recording-home-delete-integration.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-device-actions.ts
  - frontend/tests/unit/workbench-v2/recording-home-mc-resave.test.tsx
  - frontend/tests/unit/workbench-v2/writeGroupLifecycle.test.tsx
  - internal/datalink/grouppipeline/backlog_guard_delete_count_test.go
  - cmd/test_ui/group_pipeline_lifecycle_test.go
  - internal/api/router_studio_v2_workspace_devices_test.go
  - internal/datalink/grouppipeline/applied_configuration_read_test.go
  - internal/datalink/workspace/write_group_permanent_delete_recording_start_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_cascade_test.go
  - internal/api/router_studio_v2_recording_delete_contract_test.go
  - frontend/tests/unit/workbench-v2/device-delete-dialog.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-management.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-mc-locator.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-source-entry.test.tsx
  - frontend/tests/e2e/recording-database-unavailable.spec.ts
  - internal/datalink/workspace/write_group_successor_guard_test.go
  - internal/datalink/dbtarget/writer_failure_log_outcome_test.go
  - internal/api/router_studio_v2_write_group_lifecycle_history_test.go
  - internal/datalink/grouppipeline/stop_remove_intake_fence_test.go
  - cmd/test_ui/gateway_logging_database_test.go
  - frontend/tests/e2e/recording-home-narrow-records.spec.ts
  - frontend/tests/unit/workbench-v2/recording-home-tombstone-delete.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler_recordings_test.go
  - internal/datalink/workspace/write_group_lifecycle_error_test.go
  - internal/datalink/workspace/write_group_permanent_delete_schema_test.go
  - internal/datalink/workspace/write_group_permanent_delete_shared_test.go
  - frontend/tests/unit/workbench-v2/recording-home-new-range-strict-ack.test.tsx
-->

---
### Requirement: Recording rows show destination identity and write stages inline

Each recording row SHALL show its destination `schema.table` directly in a monospace secondary line with a copy control, and SHALL show write progress as four labelled stages — configuration, activation, SQL write, readback — each in exactly one of complete, unconfirmed or attention, derived from the existing home result summary. Historical or previous-revision SQL commits SHALL be labelled historical and treated as unconfirmed for the current revision. Stage details formerly behind disclosure controls (observed time, effective revision, requested boundary, owner inspection, test-write cleanup) SHALL be reachable by hover and keyboard focus on the stage. Recording rows SHALL NOT contain disclosure (`details`) elements for destination identity or result evidence. The connector identifier SHALL NOT be rendered as visible row text.

#### Scenario: Stage mapping

- **WHEN** a row renders with a result summary
- **THEN** each stage shows the state given by the mapping and an accessible text label

##### Example: stage mapping

| Summary fields | configuration | activation | SQL write | readback |
| --- | --- | --- | --- | --- |
| configuration=not_requested, effective=unknown, sql=unknown, readback=unknown | unconfirmed | unconfirmed | unconfirmed | unconfirmed |
| configuration=requested, effective=requested_effective, sql=revision_receipt, readback=test_write_verified | complete | complete | complete | complete |
| configuration=requested, effective=requested_not_effective, sql=historical_commit | complete | attention | unconfirmed (historical) | unconfirmed |
| configuration=requested, effective=requested_effective, current=delivery_attention | complete | complete | attention | unconfirmed |

#### Scenario: Copy destination table

- **WHEN** the operator activates the copy control for `public.gw_long_3fe66f8c5c98acb3`
- **THEN** that exact text is written to the clipboard and a status message confirms it, or when the clipboard is refused the text is selected and a status message states it could not be copied


<!-- @trace
source: improve-recording-home-status-and-delete
updated: 2026-11-08
code:
  - internal/datalink/dbtarget/writer_snapshot.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupLifecycleBar.tsx
  - frontend/studio-v2-diagnostic.cjs
  - internal/datalink/workspace/write_group_service.go
  - internal/datalink/runtimeview/contract.go
  - docs/validation/memory-views-20261106/skill/SKILL.md
  - frontend/src/services/runtimeLogLegacyProjection.ts
  - internal/datalink/dbtarget/writer_workers.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeResultFacts.tsx
  - frontend/src/i18n/locales/zh-TW/recording-home-status.json
  - internal/datalink/workspace/write_group_permanent_delete_audit.go
  - docs/testing.md
  - frontend/src/features/datalink/workbench-v2/recording-home/DeviceDeleteDialog.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/RecordingLifecycleHistory.tsx
  - internal/datalink/workspace/device_recording_delete.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeStatusSummary.tsx
  - internal/datalink/workspace/write_group_permanent_delete.go
  - internal/diagnostics/projection.go
  - internal/api/router.go
  - frontend/src/services/runtimeLogDomainProjection.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/recording-delete-dialog.css
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - internal/datalink/modbusshare/service_projection.go
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - frontend/src/i18n/locales/zh-TW/runtime-logs.json
  - internal/api/router_studio_v2_write_group_routes.go
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeStatusCounts.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/SavedRecordingActions.tsx
  - internal/datalink/workspace/write_group_permanent_delete_purge.go
  - internal/diagnostics/domain_event.go
  - docs/swagger/swagger.json
  - docs/swagger/docs.go
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeLifecycleMatches.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeWriteStages.tsx
  - internal/api/handlers/typed_errors.go
  - internal/datalink/device/workspace_delete.go
  - internal/datalink/groupdelivery/quota.go
  - frontend/src/services/runtimeRuleEvidenceParser.ts
  - AGENTS.md
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - internal/datalink/workspace/write_group_lifecycle_actions.go
  - frontend/src/types/runtimeRuleEvidence.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeStatusDistribution.tsx
  - internal/api/handlers/studio_v2_workspace_write_groups_handler.go
  - internal/datalink/workspace/write_group_permanent_delete_recording_start.go
  - frontend/src/i18n/locales/en/recording-home.json
  - frontend/src/services/studioV2WorkspaceWriteGroups.ts
  - internal/datalink/workspace/write_group_permanent_delete_shared.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_delete.go
  - internal/datalink/audit/types.go
  - internal/datalink/dbtarget/writer_failure_log.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/services/studioV2WorkspaceDelete.ts
  - internal/api/handlers/studio_v2_workspace_write_groups_handler_mutate.go
  - internal/datalink/dbtarget/writer_flush.go
  - frontend/src/hooks/datalink/useStudioV2WriteGroups.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeRoute.tsx
  - internal/testutil/legacyrecording/tombstone.go
  - internal/api/handlers/studio_v2_workspace_write_groups_handler_delete.go
  - frontend/src/features/datalink/workbench-v2/recording-home/deviceConnectionState.ts
  - frontend/src/i18n/locales/en/recording-home-status.json
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeRowActions.tsx
  - frontend/src/i18n/locales/zh-TW/recording-home-delete.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupDeliveryStrip.tsx
  - frontend/src/i18n/locales/zh-TW/recording-home.json
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHome.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeDestinationCell.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeWriteStageMapping.ts
  - internal/datalink/dbtarget/writer_group_flush.go
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeRows.css
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/datalink/dbtarget/writer_retry.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeUnconfirmed.tsx
  - frontend/src/features/runtime-logs/RuntimeLogLegacyEvent.tsx
  - internal/datalink/workspace/write_group_permanent_delete_operations.go
  - internal/datalink/workspace/write_group_permanent_delete_storage.go
  - frontend/src/i18n/config.ts
  - internal/datalink/workspace/write_group_permanent_delete_scope.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingDeleteDialog.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeStatus.css
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeDeviceActions.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/useRecordingHomeDeviceStatus.ts
  - docs/swagger/swagger.yaml
  - frontend/src/i18n/locales/en/recording-home-rows.json
  - frontend/src/types/studioV2Delete.ts
  - frontend/src/types/deviceConnection.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeRecordRow.tsx
  - internal/datalink/workspace/write_group_permanent_delete_guard.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeDeviceStrip.tsx
  - frontend/src/i18n/locales/zh-TW/recording-home-rows.json
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/dev/previewI18n.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/SavedRecordingRemoveDialog.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - docs/validation/memory-views-20261106/skill/references/lessons.md
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingDeleteDialogController.tsx
  - internal/datalink/grouppipeline/backlog_guard.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeCollection.tsx
  - frontend/src/i18n/locales/en/recording-home-delete.json
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingDeleteDialogError.tsx
  - frontend/src/i18n/locales/en/runtime-logs.json
tests:
  - internal/api/router_studio_v2_write_group_lifecycle_test.go
  - internal/datalink/workspace/write_group_permanent_delete_delivery_fixture_test.go
  - internal/api/router_studio_v2_successor_first_start_fixture_test.go
  - frontend/tests/unit/runtime-logs/legacy-event-fixtures.ts
  - frontend/tests/unit/workbench-v2/saved-recording-actions.test.tsx
  - internal/api/handlers/studio_v2_workspace_write_groups_handler_delete_test.go
  - frontend/tests/unit/workbench-v2/groupDeliveryStrip.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-accessibility.test.tsx
  - internal/api/router_studio_v2_locator_entry_test.go
  - frontend/tests/unit/dev/recording-style-bootstrap.test.tsx
  - internal/datalink/workspace/recording_start_disabled_test.go
  - internal/datalink/dbtarget/writer_failure_revision_test.go
  - cmd/test_ui/group_pipeline_outage_delete_test.go
  - frontend/tests/unit/workbench-v2/recording-editor-range-mc-prepare.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-memory-scope.test.tsx
  - frontend/tests/unit/workbench-v2/recording-guided-saved-card.test.tsx
  - internal/datalink/workspace/write_group_lifecycle_test.go
  - internal/datalink/workspace/write_group_mixed_safety_test.go
  - frontend/tests/unit/workbench-v2/recording-home-route.test.tsx
  - internal/datalink/workspace/write_group_permanent_delete_storage_test.go
  - frontend/tests/unit/workbench-v2/device-connection-state.test.ts
  - frontend/tests/unit/workbench-v2/recording-home-status-counts.test.ts
  - frontend/tests/unit/workbench-v2/recording-home-shell.test.tsx
  - frontend/tests/unit/workbench-v2/studio-locator-group-focus.test.tsx
  - cmd/test_ui/gateway_logging.go
  - internal/api/router_studio_v2_write_group_boundaries_test.go
  - frontend/tests/unit/workbench-v2/recording-home-source-modal.test.tsx
  - internal/datalink/grouppipeline/backlog_guard_lifecycle_test.go
  - internal/datalink/workspace/write_group_circuit_profile_rules_test.go
  - frontend/tests/unit/workbench-v2/recording-delete-dialog.test.tsx
  - internal/datalink/workspace/write_group_permanent_delete_test.go
  - internal/api/router_studio_v2_write_group_action_errors_test.go
  - frontend/tests/unit/runtime-logs/runtimeLogLegacyProjection.test.ts
  - cmd/test_ui/delivery_events_unavailable_test.go
  - frontend/tests/e2e/fixtures/recording-lifecycle-api.ts
  - frontend/tests/unit/runtime-logs/RuntimeLogCompatibility.test.tsx
  - internal/datalink/dbtarget/writer_failure_log_test.go
  - internal/api/router_runtime_rule_evidence_test.go
  - internal/api/router_studio_v2_write_group_stop_remove_safety_test.go
  - internal/datalink/groupdelivery/quota_reclaim_receipts_test.go
  - internal/datalink/workspace/write_group_basic_managed_test.go
  - frontend/tests/unit/workbench-v2/recording-modal-return.test.tsx
  - frontend/tests/unit/services/studioV2WorkspaceWriteGroups.test.ts
  - frontend/tests/unit/workbench-v2/recording-home-source-switching.test.tsx
  - cmd/test_ui/group_pipeline_outage_test.go
  - frontend/tests/unit/workbench-v2/recording-home-row-actions.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-legacy-delivery-safety.test.tsx
  - frontend/tests/unit/workbench-v2/device-delete-cache.test.tsx
  - internal/datalink/workspace/write_group_lifecycle_error_preservation_test.go
  - internal/datalink/workspace/write_group_reclaim_delete_test.go
  - cmd/test_ui/group_pipeline_lifecycle_crash_test.go
  - frontend/tests/unit/workbench-v2/recording-home-visual-status.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-device-status.test.tsx
  - internal/datalink/workspace/write_group_permanent_delete_legacy_test.go
  - frontend/tests/unit/services/studioV2DeleteResponses.test.ts
  - internal/api/router_studio_v2_recording_start_delete_test.go
  - frontend/tests/unit/workbench-v2/recording-home-memory-view.test.tsx
  - frontend/tests/unit/dev/recording-style-preview.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-acquisition-anomalies.test.tsx
  - internal/datalink/workspace/write_group_successor_replay_test.go
  - internal/datalink/workspace/write_group_lifecycle_actions_fixture_test.go
  - internal/datalink/workspace/write_group_lifecycle_actions_test.go
  - frontend/tests/e2e/recording-home-status-layout.spec.ts
  - frontend/tests/e2e/recording-home-tombstone-delete.spec.ts
  - frontend/tests/unit/workbench-v2/recording-home-inline-evidence.test.tsx
  - frontend/tests/unit/workbench-v2/recordingLifecycleBrowserFixture.test.ts
  - internal/datalink/workspace/write_group_service_test.go
  - internal/datalink/workspace/write_group_tombstone_remove_actions_test.go
  - cmd/test_ui/service_wiring.go
  - cmd/test_ui/group_pipeline_legacy_takeover_test.go
  - frontend/tests/unit/runtime-logs/runtimeLogDomainProjection.test.ts
  - frontend/tests/unit/workbench-v2/recording-home-delete-integration.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-device-actions.ts
  - frontend/tests/unit/workbench-v2/recording-home-mc-resave.test.tsx
  - frontend/tests/unit/workbench-v2/writeGroupLifecycle.test.tsx
  - internal/datalink/grouppipeline/backlog_guard_delete_count_test.go
  - cmd/test_ui/group_pipeline_lifecycle_test.go
  - internal/api/router_studio_v2_workspace_devices_test.go
  - internal/datalink/grouppipeline/applied_configuration_read_test.go
  - internal/datalink/workspace/write_group_permanent_delete_recording_start_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_cascade_test.go
  - internal/api/router_studio_v2_recording_delete_contract_test.go
  - frontend/tests/unit/workbench-v2/device-delete-dialog.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-management.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-mc-locator.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-source-entry.test.tsx
  - frontend/tests/e2e/recording-database-unavailable.spec.ts
  - internal/datalink/workspace/write_group_successor_guard_test.go
  - internal/datalink/dbtarget/writer_failure_log_outcome_test.go
  - internal/api/router_studio_v2_write_group_lifecycle_history_test.go
  - internal/datalink/grouppipeline/stop_remove_intake_fence_test.go
  - cmd/test_ui/gateway_logging_database_test.go
  - frontend/tests/e2e/recording-home-narrow-records.spec.ts
  - frontend/tests/unit/workbench-v2/recording-home-tombstone-delete.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler_recordings_test.go
  - internal/datalink/workspace/write_group_lifecycle_error_test.go
  - internal/datalink/workspace/write_group_permanent_delete_schema_test.go
  - internal/datalink/workspace/write_group_permanent_delete_shared_test.go
  - frontend/tests/unit/workbench-v2/recording-home-new-range-strict-ack.test.tsx
-->

---
### Requirement: Home offers recording and device actions with explicit confirmation

Each recording row SHALL keep its open action and add a keyboard-operable actions menu offering stop when the recording is running and delete when server action facts allow it; stale action facts SHALL disable both with a reload hint. Each device in the status strip and devices view SHALL offer delete through the same menu pattern. The recording delete dialog SHALL state what is permanently deleted, what external data is retained, that the action cannot be undone, and SHALL require an explicit acknowledgement before submission. The device delete dialog SHALL list recordings referencing the device; when any exist it SHALL block confirmation and explain they must be deleted first, otherwise it SHALL require typing the exact device name. Server conflicts SHALL keep the dialog open with an actionable explanation per error code and SHALL NOT be shown as success.

#### Scenario: Delete a recording from the home

- **WHEN** the operator opens the row menu, chooses delete, acknowledges and confirms, and the server returns success
- **THEN** the recording disappears from the list and the summary counts update without a page reload

#### Scenario: Delete is refused because of pending delivery

- **WHEN** the server returns `WRITE_GROUP_DELETE_BACKLOG` with 3 unconsumed samples and 2 pending outbox rows
- **THEN** the dialog shows the samples and outbox rows as separately labelled counts, offers stop and refresh actions, and the recording remains listed

#### Scenario: Device still referenced

- **WHEN** the operator chooses delete for a device that 2 non-deleted recordings reference
- **THEN** the dialog names both recordings, the confirm action is disabled, and no delete request is sent

#### Scenario: Device name mismatch

- **WHEN** an unreferenced device named `mc_new` is being deleted and the typed name is `mc-new`
- **THEN** the confirm action remains disabled


<!-- @trace
source: improve-recording-home-status-and-delete
updated: 2026-11-08
code:
  - internal/datalink/dbtarget/writer_snapshot.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupLifecycleBar.tsx
  - frontend/studio-v2-diagnostic.cjs
  - internal/datalink/workspace/write_group_service.go
  - internal/datalink/runtimeview/contract.go
  - docs/validation/memory-views-20261106/skill/SKILL.md
  - frontend/src/services/runtimeLogLegacyProjection.ts
  - internal/datalink/dbtarget/writer_workers.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeResultFacts.tsx
  - frontend/src/i18n/locales/zh-TW/recording-home-status.json
  - internal/datalink/workspace/write_group_permanent_delete_audit.go
  - docs/testing.md
  - frontend/src/features/datalink/workbench-v2/recording-home/DeviceDeleteDialog.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/RecordingLifecycleHistory.tsx
  - internal/datalink/workspace/device_recording_delete.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeStatusSummary.tsx
  - internal/datalink/workspace/write_group_permanent_delete.go
  - internal/diagnostics/projection.go
  - internal/api/router.go
  - frontend/src/services/runtimeLogDomainProjection.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/recording-delete-dialog.css
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - internal/datalink/modbusshare/service_projection.go
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - frontend/src/i18n/locales/zh-TW/runtime-logs.json
  - internal/api/router_studio_v2_write_group_routes.go
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeStatusCounts.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/SavedRecordingActions.tsx
  - internal/datalink/workspace/write_group_permanent_delete_purge.go
  - internal/diagnostics/domain_event.go
  - docs/swagger/swagger.json
  - docs/swagger/docs.go
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeLifecycleMatches.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeWriteStages.tsx
  - internal/api/handlers/typed_errors.go
  - internal/datalink/device/workspace_delete.go
  - internal/datalink/groupdelivery/quota.go
  - frontend/src/services/runtimeRuleEvidenceParser.ts
  - AGENTS.md
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - internal/datalink/workspace/write_group_lifecycle_actions.go
  - frontend/src/types/runtimeRuleEvidence.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeStatusDistribution.tsx
  - internal/api/handlers/studio_v2_workspace_write_groups_handler.go
  - internal/datalink/workspace/write_group_permanent_delete_recording_start.go
  - frontend/src/i18n/locales/en/recording-home.json
  - frontend/src/services/studioV2WorkspaceWriteGroups.ts
  - internal/datalink/workspace/write_group_permanent_delete_shared.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_delete.go
  - internal/datalink/audit/types.go
  - internal/datalink/dbtarget/writer_failure_log.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/services/studioV2WorkspaceDelete.ts
  - internal/api/handlers/studio_v2_workspace_write_groups_handler_mutate.go
  - internal/datalink/dbtarget/writer_flush.go
  - frontend/src/hooks/datalink/useStudioV2WriteGroups.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeRoute.tsx
  - internal/testutil/legacyrecording/tombstone.go
  - internal/api/handlers/studio_v2_workspace_write_groups_handler_delete.go
  - frontend/src/features/datalink/workbench-v2/recording-home/deviceConnectionState.ts
  - frontend/src/i18n/locales/en/recording-home-status.json
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeRowActions.tsx
  - frontend/src/i18n/locales/zh-TW/recording-home-delete.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupDeliveryStrip.tsx
  - frontend/src/i18n/locales/zh-TW/recording-home.json
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHome.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeDestinationCell.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeWriteStageMapping.ts
  - internal/datalink/dbtarget/writer_group_flush.go
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeRows.css
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/datalink/dbtarget/writer_retry.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeUnconfirmed.tsx
  - frontend/src/features/runtime-logs/RuntimeLogLegacyEvent.tsx
  - internal/datalink/workspace/write_group_permanent_delete_operations.go
  - internal/datalink/workspace/write_group_permanent_delete_storage.go
  - frontend/src/i18n/config.ts
  - internal/datalink/workspace/write_group_permanent_delete_scope.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingDeleteDialog.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeStatus.css
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeDeviceActions.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/useRecordingHomeDeviceStatus.ts
  - docs/swagger/swagger.yaml
  - frontend/src/i18n/locales/en/recording-home-rows.json
  - frontend/src/types/studioV2Delete.ts
  - frontend/src/types/deviceConnection.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeRecordRow.tsx
  - internal/datalink/workspace/write_group_permanent_delete_guard.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeDeviceStrip.tsx
  - frontend/src/i18n/locales/zh-TW/recording-home-rows.json
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/dev/previewI18n.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/SavedRecordingRemoveDialog.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - docs/validation/memory-views-20261106/skill/references/lessons.md
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingDeleteDialogController.tsx
  - internal/datalink/grouppipeline/backlog_guard.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeCollection.tsx
  - frontend/src/i18n/locales/en/recording-home-delete.json
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingDeleteDialogError.tsx
  - frontend/src/i18n/locales/en/runtime-logs.json
tests:
  - internal/api/router_studio_v2_write_group_lifecycle_test.go
  - internal/datalink/workspace/write_group_permanent_delete_delivery_fixture_test.go
  - internal/api/router_studio_v2_successor_first_start_fixture_test.go
  - frontend/tests/unit/runtime-logs/legacy-event-fixtures.ts
  - frontend/tests/unit/workbench-v2/saved-recording-actions.test.tsx
  - internal/api/handlers/studio_v2_workspace_write_groups_handler_delete_test.go
  - frontend/tests/unit/workbench-v2/groupDeliveryStrip.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-accessibility.test.tsx
  - internal/api/router_studio_v2_locator_entry_test.go
  - frontend/tests/unit/dev/recording-style-bootstrap.test.tsx
  - internal/datalink/workspace/recording_start_disabled_test.go
  - internal/datalink/dbtarget/writer_failure_revision_test.go
  - cmd/test_ui/group_pipeline_outage_delete_test.go
  - frontend/tests/unit/workbench-v2/recording-editor-range-mc-prepare.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-memory-scope.test.tsx
  - frontend/tests/unit/workbench-v2/recording-guided-saved-card.test.tsx
  - internal/datalink/workspace/write_group_lifecycle_test.go
  - internal/datalink/workspace/write_group_mixed_safety_test.go
  - frontend/tests/unit/workbench-v2/recording-home-route.test.tsx
  - internal/datalink/workspace/write_group_permanent_delete_storage_test.go
  - frontend/tests/unit/workbench-v2/device-connection-state.test.ts
  - frontend/tests/unit/workbench-v2/recording-home-status-counts.test.ts
  - frontend/tests/unit/workbench-v2/recording-home-shell.test.tsx
  - frontend/tests/unit/workbench-v2/studio-locator-group-focus.test.tsx
  - cmd/test_ui/gateway_logging.go
  - internal/api/router_studio_v2_write_group_boundaries_test.go
  - frontend/tests/unit/workbench-v2/recording-home-source-modal.test.tsx
  - internal/datalink/grouppipeline/backlog_guard_lifecycle_test.go
  - internal/datalink/workspace/write_group_circuit_profile_rules_test.go
  - frontend/tests/unit/workbench-v2/recording-delete-dialog.test.tsx
  - internal/datalink/workspace/write_group_permanent_delete_test.go
  - internal/api/router_studio_v2_write_group_action_errors_test.go
  - frontend/tests/unit/runtime-logs/runtimeLogLegacyProjection.test.ts
  - cmd/test_ui/delivery_events_unavailable_test.go
  - frontend/tests/e2e/fixtures/recording-lifecycle-api.ts
  - frontend/tests/unit/runtime-logs/RuntimeLogCompatibility.test.tsx
  - internal/datalink/dbtarget/writer_failure_log_test.go
  - internal/api/router_runtime_rule_evidence_test.go
  - internal/api/router_studio_v2_write_group_stop_remove_safety_test.go
  - internal/datalink/groupdelivery/quota_reclaim_receipts_test.go
  - internal/datalink/workspace/write_group_basic_managed_test.go
  - frontend/tests/unit/workbench-v2/recording-modal-return.test.tsx
  - frontend/tests/unit/services/studioV2WorkspaceWriteGroups.test.ts
  - frontend/tests/unit/workbench-v2/recording-home-source-switching.test.tsx
  - cmd/test_ui/group_pipeline_outage_test.go
  - frontend/tests/unit/workbench-v2/recording-home-row-actions.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-legacy-delivery-safety.test.tsx
  - frontend/tests/unit/workbench-v2/device-delete-cache.test.tsx
  - internal/datalink/workspace/write_group_lifecycle_error_preservation_test.go
  - internal/datalink/workspace/write_group_reclaim_delete_test.go
  - cmd/test_ui/group_pipeline_lifecycle_crash_test.go
  - frontend/tests/unit/workbench-v2/recording-home-visual-status.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-device-status.test.tsx
  - internal/datalink/workspace/write_group_permanent_delete_legacy_test.go
  - frontend/tests/unit/services/studioV2DeleteResponses.test.ts
  - internal/api/router_studio_v2_recording_start_delete_test.go
  - frontend/tests/unit/workbench-v2/recording-home-memory-view.test.tsx
  - frontend/tests/unit/dev/recording-style-preview.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-acquisition-anomalies.test.tsx
  - internal/datalink/workspace/write_group_successor_replay_test.go
  - internal/datalink/workspace/write_group_lifecycle_actions_fixture_test.go
  - internal/datalink/workspace/write_group_lifecycle_actions_test.go
  - frontend/tests/e2e/recording-home-status-layout.spec.ts
  - frontend/tests/e2e/recording-home-tombstone-delete.spec.ts
  - frontend/tests/unit/workbench-v2/recording-home-inline-evidence.test.tsx
  - frontend/tests/unit/workbench-v2/recordingLifecycleBrowserFixture.test.ts
  - internal/datalink/workspace/write_group_service_test.go
  - internal/datalink/workspace/write_group_tombstone_remove_actions_test.go
  - cmd/test_ui/service_wiring.go
  - cmd/test_ui/group_pipeline_legacy_takeover_test.go
  - frontend/tests/unit/runtime-logs/runtimeLogDomainProjection.test.ts
  - frontend/tests/unit/workbench-v2/recording-home-delete-integration.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-device-actions.ts
  - frontend/tests/unit/workbench-v2/recording-home-mc-resave.test.tsx
  - frontend/tests/unit/workbench-v2/writeGroupLifecycle.test.tsx
  - internal/datalink/grouppipeline/backlog_guard_delete_count_test.go
  - cmd/test_ui/group_pipeline_lifecycle_test.go
  - internal/api/router_studio_v2_workspace_devices_test.go
  - internal/datalink/grouppipeline/applied_configuration_read_test.go
  - internal/datalink/workspace/write_group_permanent_delete_recording_start_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_cascade_test.go
  - internal/api/router_studio_v2_recording_delete_contract_test.go
  - frontend/tests/unit/workbench-v2/device-delete-dialog.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-management.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-mc-locator.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-source-entry.test.tsx
  - frontend/tests/e2e/recording-database-unavailable.spec.ts
  - internal/datalink/workspace/write_group_successor_guard_test.go
  - internal/datalink/dbtarget/writer_failure_log_outcome_test.go
  - internal/api/router_studio_v2_write_group_lifecycle_history_test.go
  - internal/datalink/grouppipeline/stop_remove_intake_fence_test.go
  - cmd/test_ui/gateway_logging_database_test.go
  - frontend/tests/e2e/recording-home-narrow-records.spec.ts
  - frontend/tests/unit/workbench-v2/recording-home-tombstone-delete.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler_recordings_test.go
  - internal/datalink/workspace/write_group_lifecycle_error_test.go
  - internal/datalink/workspace/write_group_permanent_delete_schema_test.go
  - internal/datalink/workspace/write_group_permanent_delete_shared_test.go
  - frontend/tests/unit/workbench-v2/recording-home-new-range-strict-ack.test.tsx
-->

---
### Requirement: Home header has one primary action and visible freshness

The recording home header SHALL present add recording as the only primary action, with settings, logs and refresh as secondary actions. The two observation explanation lines SHALL be replaced by a last-updated time next to refresh, with the full explanation available from an adjacent help control. A resume-draft banner SHALL name the draft when one exists and SHALL be absent otherwise. The unconfirmed or incomplete recordings region SHALL be collapsible with a heading stating the item count, expanded by default when any item is not a tombstone and collapsed when all items are tombstones. At widths of 1440, 1024, 768 and 390 pixels the page SHALL NOT scroll horizontally, and at 390 pixels each recording SHALL render as a card whose secondary actions are in the actions menu.

#### Scenario: Only tombstones are unconfirmed

- **WHEN** the only unconfirmed item is a removed recording
- **THEN** the region is collapsed and its heading reads 1 item needs attention

#### Scenario: Narrow viewport

- **WHEN** the home renders at 390 pixels wide
- **THEN** the document scroll width equals the viewport width and recordings render as cards

<!-- @trace
source: improve-recording-home-status-and-delete
updated: 2026-11-08
code:
  - internal/datalink/dbtarget/writer_snapshot.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupLifecycleBar.tsx
  - frontend/studio-v2-diagnostic.cjs
  - internal/datalink/workspace/write_group_service.go
  - internal/datalink/runtimeview/contract.go
  - docs/validation/memory-views-20261106/skill/SKILL.md
  - frontend/src/services/runtimeLogLegacyProjection.ts
  - internal/datalink/dbtarget/writer_workers.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeResultFacts.tsx
  - frontend/src/i18n/locales/zh-TW/recording-home-status.json
  - internal/datalink/workspace/write_group_permanent_delete_audit.go
  - docs/testing.md
  - frontend/src/features/datalink/workbench-v2/recording-home/DeviceDeleteDialog.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/RecordingLifecycleHistory.tsx
  - internal/datalink/workspace/device_recording_delete.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeStatusSummary.tsx
  - internal/datalink/workspace/write_group_permanent_delete.go
  - internal/diagnostics/projection.go
  - internal/api/router.go
  - frontend/src/services/runtimeLogDomainProjection.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/recording-delete-dialog.css
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - internal/datalink/modbusshare/service_projection.go
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - frontend/src/i18n/locales/zh-TW/runtime-logs.json
  - internal/api/router_studio_v2_write_group_routes.go
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeStatusCounts.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/SavedRecordingActions.tsx
  - internal/datalink/workspace/write_group_permanent_delete_purge.go
  - internal/diagnostics/domain_event.go
  - docs/swagger/swagger.json
  - docs/swagger/docs.go
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeLifecycleMatches.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeWriteStages.tsx
  - internal/api/handlers/typed_errors.go
  - internal/datalink/device/workspace_delete.go
  - internal/datalink/groupdelivery/quota.go
  - frontend/src/services/runtimeRuleEvidenceParser.ts
  - AGENTS.md
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - internal/datalink/workspace/write_group_lifecycle_actions.go
  - frontend/src/types/runtimeRuleEvidence.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeStatusDistribution.tsx
  - internal/api/handlers/studio_v2_workspace_write_groups_handler.go
  - internal/datalink/workspace/write_group_permanent_delete_recording_start.go
  - frontend/src/i18n/locales/en/recording-home.json
  - frontend/src/services/studioV2WorkspaceWriteGroups.ts
  - internal/datalink/workspace/write_group_permanent_delete_shared.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_delete.go
  - internal/datalink/audit/types.go
  - internal/datalink/dbtarget/writer_failure_log.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/services/studioV2WorkspaceDelete.ts
  - internal/api/handlers/studio_v2_workspace_write_groups_handler_mutate.go
  - internal/datalink/dbtarget/writer_flush.go
  - frontend/src/hooks/datalink/useStudioV2WriteGroups.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeRoute.tsx
  - internal/testutil/legacyrecording/tombstone.go
  - internal/api/handlers/studio_v2_workspace_write_groups_handler_delete.go
  - frontend/src/features/datalink/workbench-v2/recording-home/deviceConnectionState.ts
  - frontend/src/i18n/locales/en/recording-home-status.json
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeRowActions.tsx
  - frontend/src/i18n/locales/zh-TW/recording-home-delete.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupDeliveryStrip.tsx
  - frontend/src/i18n/locales/zh-TW/recording-home.json
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHome.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeDestinationCell.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeWriteStageMapping.ts
  - internal/datalink/dbtarget/writer_group_flush.go
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeRows.css
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/datalink/dbtarget/writer_retry.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeUnconfirmed.tsx
  - frontend/src/features/runtime-logs/RuntimeLogLegacyEvent.tsx
  - internal/datalink/workspace/write_group_permanent_delete_operations.go
  - internal/datalink/workspace/write_group_permanent_delete_storage.go
  - frontend/src/i18n/config.ts
  - internal/datalink/workspace/write_group_permanent_delete_scope.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingDeleteDialog.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeStatus.css
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeDeviceActions.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/useRecordingHomeDeviceStatus.ts
  - docs/swagger/swagger.yaml
  - frontend/src/i18n/locales/en/recording-home-rows.json
  - frontend/src/types/studioV2Delete.ts
  - frontend/src/types/deviceConnection.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeRecordRow.tsx
  - internal/datalink/workspace/write_group_permanent_delete_guard.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeDeviceStrip.tsx
  - frontend/src/i18n/locales/zh-TW/recording-home-rows.json
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/dev/previewI18n.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/SavedRecordingRemoveDialog.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - docs/validation/memory-views-20261106/skill/references/lessons.md
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingDeleteDialogController.tsx
  - internal/datalink/grouppipeline/backlog_guard.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeCollection.tsx
  - frontend/src/i18n/locales/en/recording-home-delete.json
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingDeleteDialogError.tsx
  - frontend/src/i18n/locales/en/runtime-logs.json
tests:
  - internal/api/router_studio_v2_write_group_lifecycle_test.go
  - internal/datalink/workspace/write_group_permanent_delete_delivery_fixture_test.go
  - internal/api/router_studio_v2_successor_first_start_fixture_test.go
  - frontend/tests/unit/runtime-logs/legacy-event-fixtures.ts
  - frontend/tests/unit/workbench-v2/saved-recording-actions.test.tsx
  - internal/api/handlers/studio_v2_workspace_write_groups_handler_delete_test.go
  - frontend/tests/unit/workbench-v2/groupDeliveryStrip.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-accessibility.test.tsx
  - internal/api/router_studio_v2_locator_entry_test.go
  - frontend/tests/unit/dev/recording-style-bootstrap.test.tsx
  - internal/datalink/workspace/recording_start_disabled_test.go
  - internal/datalink/dbtarget/writer_failure_revision_test.go
  - cmd/test_ui/group_pipeline_outage_delete_test.go
  - frontend/tests/unit/workbench-v2/recording-editor-range-mc-prepare.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-memory-scope.test.tsx
  - frontend/tests/unit/workbench-v2/recording-guided-saved-card.test.tsx
  - internal/datalink/workspace/write_group_lifecycle_test.go
  - internal/datalink/workspace/write_group_mixed_safety_test.go
  - frontend/tests/unit/workbench-v2/recording-home-route.test.tsx
  - internal/datalink/workspace/write_group_permanent_delete_storage_test.go
  - frontend/tests/unit/workbench-v2/device-connection-state.test.ts
  - frontend/tests/unit/workbench-v2/recording-home-status-counts.test.ts
  - frontend/tests/unit/workbench-v2/recording-home-shell.test.tsx
  - frontend/tests/unit/workbench-v2/studio-locator-group-focus.test.tsx
  - cmd/test_ui/gateway_logging.go
  - internal/api/router_studio_v2_write_group_boundaries_test.go
  - frontend/tests/unit/workbench-v2/recording-home-source-modal.test.tsx
  - internal/datalink/grouppipeline/backlog_guard_lifecycle_test.go
  - internal/datalink/workspace/write_group_circuit_profile_rules_test.go
  - frontend/tests/unit/workbench-v2/recording-delete-dialog.test.tsx
  - internal/datalink/workspace/write_group_permanent_delete_test.go
  - internal/api/router_studio_v2_write_group_action_errors_test.go
  - frontend/tests/unit/runtime-logs/runtimeLogLegacyProjection.test.ts
  - cmd/test_ui/delivery_events_unavailable_test.go
  - frontend/tests/e2e/fixtures/recording-lifecycle-api.ts
  - frontend/tests/unit/runtime-logs/RuntimeLogCompatibility.test.tsx
  - internal/datalink/dbtarget/writer_failure_log_test.go
  - internal/api/router_runtime_rule_evidence_test.go
  - internal/api/router_studio_v2_write_group_stop_remove_safety_test.go
  - internal/datalink/groupdelivery/quota_reclaim_receipts_test.go
  - internal/datalink/workspace/write_group_basic_managed_test.go
  - frontend/tests/unit/workbench-v2/recording-modal-return.test.tsx
  - frontend/tests/unit/services/studioV2WorkspaceWriteGroups.test.ts
  - frontend/tests/unit/workbench-v2/recording-home-source-switching.test.tsx
  - cmd/test_ui/group_pipeline_outage_test.go
  - frontend/tests/unit/workbench-v2/recording-home-row-actions.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-legacy-delivery-safety.test.tsx
  - frontend/tests/unit/workbench-v2/device-delete-cache.test.tsx
  - internal/datalink/workspace/write_group_lifecycle_error_preservation_test.go
  - internal/datalink/workspace/write_group_reclaim_delete_test.go
  - cmd/test_ui/group_pipeline_lifecycle_crash_test.go
  - frontend/tests/unit/workbench-v2/recording-home-visual-status.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-device-status.test.tsx
  - internal/datalink/workspace/write_group_permanent_delete_legacy_test.go
  - frontend/tests/unit/services/studioV2DeleteResponses.test.ts
  - internal/api/router_studio_v2_recording_start_delete_test.go
  - frontend/tests/unit/workbench-v2/recording-home-memory-view.test.tsx
  - frontend/tests/unit/dev/recording-style-preview.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-acquisition-anomalies.test.tsx
  - internal/datalink/workspace/write_group_successor_replay_test.go
  - internal/datalink/workspace/write_group_lifecycle_actions_fixture_test.go
  - internal/datalink/workspace/write_group_lifecycle_actions_test.go
  - frontend/tests/e2e/recording-home-status-layout.spec.ts
  - frontend/tests/e2e/recording-home-tombstone-delete.spec.ts
  - frontend/tests/unit/workbench-v2/recording-home-inline-evidence.test.tsx
  - frontend/tests/unit/workbench-v2/recordingLifecycleBrowserFixture.test.ts
  - internal/datalink/workspace/write_group_service_test.go
  - internal/datalink/workspace/write_group_tombstone_remove_actions_test.go
  - cmd/test_ui/service_wiring.go
  - cmd/test_ui/group_pipeline_legacy_takeover_test.go
  - frontend/tests/unit/runtime-logs/runtimeLogDomainProjection.test.ts
  - frontend/tests/unit/workbench-v2/recording-home-delete-integration.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-device-actions.ts
  - frontend/tests/unit/workbench-v2/recording-home-mc-resave.test.tsx
  - frontend/tests/unit/workbench-v2/writeGroupLifecycle.test.tsx
  - internal/datalink/grouppipeline/backlog_guard_delete_count_test.go
  - cmd/test_ui/group_pipeline_lifecycle_test.go
  - internal/api/router_studio_v2_workspace_devices_test.go
  - internal/datalink/grouppipeline/applied_configuration_read_test.go
  - internal/datalink/workspace/write_group_permanent_delete_recording_start_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_cascade_test.go
  - internal/api/router_studio_v2_recording_delete_contract_test.go
  - frontend/tests/unit/workbench-v2/device-delete-dialog.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-management.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-mc-locator.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-source-entry.test.tsx
  - frontend/tests/e2e/recording-database-unavailable.spec.ts
  - internal/datalink/workspace/write_group_successor_guard_test.go
  - internal/datalink/dbtarget/writer_failure_log_outcome_test.go
  - internal/api/router_studio_v2_write_group_lifecycle_history_test.go
  - internal/datalink/grouppipeline/stop_remove_intake_fence_test.go
  - cmd/test_ui/gateway_logging_database_test.go
  - frontend/tests/e2e/recording-home-narrow-records.spec.ts
  - frontend/tests/unit/workbench-v2/recording-home-tombstone-delete.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler_recordings_test.go
  - internal/datalink/workspace/write_group_lifecycle_error_test.go
  - internal/datalink/workspace/write_group_permanent_delete_schema_test.go
  - internal/datalink/workspace/write_group_permanent_delete_shared_test.go
  - frontend/tests/unit/workbench-v2/recording-home-new-range-strict-ack.test.tsx
-->