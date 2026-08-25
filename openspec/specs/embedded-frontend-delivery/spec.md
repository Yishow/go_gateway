# embedded-frontend-delivery Specification

## Purpose

Define the delivery contract for the React frontend when it is embedded in the Go binary: the Vite build must produce a complete relative asset graph, each supported build/start path must synchronize that graph to `cmd/test_ui/static` and fail closed on missing or failed inputs, and route-level lazy loading must preserve existing route semantics. The contract is verified at source, synchronization, fresh-clone, and browser boundaries, including deterministic loading/error behavior and cross-platform Playwright/E2E coverage.

## Requirements

### Requirement: Frontend build emits a complete embedded asset graph

The frontend build SHALL emit an entry document at frontend/dist/index.html and every JavaScript, CSS, route chunk, vendor chunk, and referenced static asset required by its direct or transitive import graph. Asset references SHALL remain relative to the embedded static root. Vendor grouping SHALL use stable logical groups for the React/runtime, Material UI and Emotion, charting, router/i18n/query application dependencies, and residual dependencies.

#### Scenario: Build emits entry and lazy route assets

- **WHEN** the repository frontend build runs with its declared dependencies
- **THEN** frontend/dist/index.html exists
- **AND** the output contains the lazy assets needed by the existing /studio/v2, /studio, /studio/runtime, /test, and /gateway/* route families
- **AND** every asset reference discovered from the entry and lazy import graph resolves to a file under frontend/dist

##### Example: vendor group names

- **GIVEN** a production build includes React, Material UI, charting, and router dependencies
- **WHEN** Vite writes the production output
- **THEN** the output contains stable logical vendor groups named react-vendor, mui-vendor, charts-vendor, app-vendor, or misc-vendor for the matching dependency families
- **AND** application route chunks remain separate from those vendor groups

<!-- @trace
source: reliable-embedded-frontend-delivery
updated: 2026-08-22
code:
  - frontend/vite.config.ts
  - frontend/src/App.tsx
tests:
  - frontend/tests/unit/app-routing-lazy-load.test.tsx
  - tests/shell/embedded-frontend-delivery.sh
-->

---
### Requirement: Embedded static synchronization mirrors the complete frontend output

Each supported synchronization path SHALL ensure cmd/test_ui/static exists, SHALL remove generated children that are absent from the current frontend/dist while preserving the tracked embed-placeholder.txt, and SHALL copy the complete frontend/dist tree including nested directories and all asset types. Generated children under cmd/test_ui/static other than embed-placeholder.txt SHALL remain ignored by Git.

#### Scenario: Synchronization removes stale output and preserves the placeholder

- **GIVEN** cmd/test_ui/static contains embed-placeholder.txt, current assets, and a stale generated file that is absent from frontend/dist
- **WHEN** scripts/build.ps1, start.ps1 Build-Frontend, or start.sh build_frontend completes a successful synchronization
- **THEN** embed-placeholder.txt still exists
- **AND** the stale generated file no longer exists
- **AND** every file and directory in frontend/dist exists at the matching relative path under cmd/test_ui/static
- **AND** no generated static child outside the current frontend/dist tree remains

#### Scenario: Fresh clone provides the embed directory marker

- **GIVEN** a fresh checkout has the tracked cmd/test_ui/static/embed-placeholder.txt and has no generated frontend assets
- **WHEN** the Go embedded frontend package is compiled
- **THEN** the compiler finds a match for the static embed pattern without requiring frontend/dist
- **AND** no generated asset is required to make the checkout compile

<!-- @trace
source: reliable-embedded-frontend-delivery
updated: 2026-08-22
code:
  - .gitignore
  - cmd/test_ui/static/embed-placeholder.txt
  - scripts/build.ps1
  - start.ps1
  - start.sh
tests:
  - tests/shell/embedded-frontend-delivery.sh
-->

---
### Requirement: Build and synchronization failures are surfaced and fail closed

The frontend build and synchronization commands SHALL return a nonzero exit status when dependency installation, frontend source availability, frontend/dist generation, target-directory creation, stale-asset cleanup, or asset copy fails. A failed frontend build, missing frontend source, or missing frontend/dist SHALL NOT be followed by a backend build that is reported as a successful frontend delivery.

#### Scenario: Missing dist stops synchronization

- **GIVEN** frontend/dist does not exist after the frontend build step
- **WHEN** a supported build or start path attempts to synchronize embedded assets
- **THEN** the command returns a nonzero exit status
- **AND** the output identifies frontend/dist as unavailable
- **AND** the backend build step is not reported as successful for that invocation

#### Scenario: Missing frontend source stops the POSIX build path

- **GIVEN** the `frontend` source directory or its required build inputs are absent before `start.sh` runs its frontend build path
- **WHEN** the POSIX build path attempts to build and then invoke the backend compiler
- **THEN** the command returns a nonzero exit status
- **AND** the output identifies the missing frontend source
- **AND** the backend compiler is not invoked in that invocation

#### Scenario: Copy failure remains actionable

- **GIVEN** the target static directory cannot be cleaned or written
- **WHEN** a supported synchronization path performs the copy
- **THEN** the command returns a nonzero exit status
- **AND** the existing console or log channel contains the copy failure
- **AND** the command does not claim that embedded frontend delivery completed

<!-- @trace
source: reliable-embedded-frontend-delivery
updated: 2026-08-22
code:
  - scripts/build.ps1
  - start.ps1
  - start.sh
tests:
  - tests/shell/embedded-frontend-delivery.sh
-->

---
### Requirement: Lazy route loading has one deterministic loading and error contract

The application SHALL load heavy existing route modules through route-level dynamic imports and SHALL render one deterministic Suspense fallback while the selected route chunk is pending. A rejected route or vendor import SHALL reach the existing error boundary as an explicit error state. The implementation MUST preserve each route component's existing props and route parameters.

#### Scenario: Core route waits on the deterministic fallback

- **GIVEN** the application has started and a core route chunk has not finished loading
- **WHEN** an operator navigates to /studio/v2, /studio, /studio/runtime, /test, or a /gateway/* route
- **THEN** the route container renders the same deterministic loading fallback while the import is pending
- **AND** the selected route renders after its lazy chunk and vendor dependencies load
- **AND** no second product route or redirect is introduced

#### Scenario: Missing route chunk is explicit

- **GIVEN** a selected route chunk or vendor chunk is unavailable
- **WHEN** the browser evaluates its dynamic import
- **THEN** the existing error boundary renders an error state
- **AND** the application does not present the route as successfully loaded
- **AND** the failed asset request remains observable as a missing-resource response

<!-- @trace
source: reliable-embedded-frontend-delivery
updated: 2026-08-22
code:
  - frontend/src/App.tsx
tests:
  - frontend/tests/unit/app-routing-lazy-load.test.tsx
-->

---
### Requirement: Existing route product semantics remain unchanged during delivery refactoring

The embedded delivery implementation SHALL retain the canonical route contract: the default entry resolves to `/studio/v2`, `/studio/runtime` remains the focused monitor, `/test` remains an independent engineering tool, and `/gateway/*` remains experimental. The `/studio` route SHALL be removed immediately with its dedicated registration, handler/tombstone, and proven legacy-only assets. After deletion, `/studio` SHALL be handled exactly like an arbitrary unknown route under the repository generic unknown-route policy, without loading legacy workspace code. Lazy loading and chunk grouping SHALL NOT alter supported route parameters or non-target route identities.

#### Scenario: Default and deleted routes retain their identities

- **WHEN** an operator opens the root entry, generic datalink landing, `/studio`, and `/studio/v2` after deletion
- **THEN** root and generic landing resolve to the v2 entry contract, `/studio` follows the same generic unknown-route policy as an arbitrary unknown path, and `/studio/v2` renders the guided user-facing entry
- **AND** no dedicated `/studio` route/handler/tombstone, special redirect, legacy mount, or legacy-only asset request occurs

#### Scenario: Monitor, tool, and experimental routes remain distinct

- **WHEN** an operator opens `/studio/runtime` with its existing device context, `/test`, and one `/gateway/*` route
- **THEN** `/studio/runtime` renders the focused runtime monitor contract
- **AND** `/test` remains an independent engineering tool
- **AND** `/gateway/*` remains an experimental surface
- **AND** lazy loading introduces no redirect from these non-target routes to an unrelated product flow

##### Example: Route family identity checks

- **GIVEN** the runtime URL includes `device_id=device-A` and the gateway URL is `/gateway/quick-setup`
- **WHEN** the operator opens `/studio/runtime?device_id=device-A`, `/test`, and `/gateway/quick-setup`
- **THEN** the runtime page remains the focused monitor for device-A
- **AND** `/test` remains the independent engineering tool
- **AND** `/gateway/quick-setup` remains experimental without a redirect to `/studio/v2`

#### Scenario: Deleted route is equivalent to an arbitrary unknown path

- **WHEN** a route test requests `/studio` and `/unknown-route-for-retirement` after deletion
- **THEN** both requests produce the same response/navigation, status/fallback, error boundary, and network asset pattern
- **AND** neither request loads a legacy route chunk or legacy-only asset


<!-- @trace
source: retire-legacy-studio-and-polish-v2
updated: 2026-08-26
code:
  - internal/api/handlers/device_health_handler.go
  - internal/datalink/sourcerule/candidate_api.go
  - internal/api/router_modbus_share.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - internal/api/handlers/config.go
  - internal/datalink/sourcerule/candidate_snapshot_local_modbus_outputs.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/vite.config.ts
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - internal/datalink/migrator.go
  - internal/datalink/sourcerule/service_state.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - internal/datalink/modbusshare/service_persistence.go
  - internal/datalink/modbusshare/types.go
  - frontend/src/services/sourceRuleCandidates.ts
  - internal/datalink/sourcerule/validation.go
  - internal/datalink/dbtarget/service.go
  - internal/datalink/workspace/service_readiness_device.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsStatus.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDiagnosticsPanel.tsx
  - internal/virtual/server/modbus/server.go
  - frontend/src/utils/safeJson.ts
  - internal/datalink/collector/scheduler_lifecycle.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - internal/api/handlers/runtime_stream_handler.go
  - internal/datalink/device/service_readiness_safe.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - internal/datalink/runtime/target_delivery.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/sourcerule/share_restore_projection.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/useStep4Activation.ts
  - frontend/src/services/datalinkClient.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/api/handlers/modbus_share_handler_mapping_query.go
  - internal/datalink/sourcerule/repository_memory.go
  - frontend/src/hooks/datalink/useRuntimeStream.ts
  - internal/datalink/settings/sql_repo.go
  - internal/datalink/sourcerule/candidate_snapshot_local_modbus_conflicts.go
  - internal/datalink/sourcerule/service.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - internal/datalink/workspace/service_runtime_projection_scope.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/api/handlers/mapping_handler.go
  - frontend/src/features/datalink/workbench-v2/state/types-settings.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - docs/swagger/swagger.json
  - frontend/src/features/datalink/workbench-v2/settings/useSettingsOperations.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - internal/datalink/runtime/delivery_diagnostic.go
  - internal/api/handlers/debug.go
  - internal/datalink/sourcerule/tag_review_decision_service.go
  - internal/api/handlers/connection.go
  - internal/datalink/modbusshare/sql_revision_store_helpers.go
  - internal/datalink/tag/service_crud.go
  - internal/datalink/sourcerule/candidate_snapshot_revision.go
  - internal/datalink/collector/scheduler.go
  - internal/datalink/schema/migrations/001_initial_schema.up.sql
  - frontend/src/features/datalink/workbench-v2/steps/step4/step4DatabaseHelpers.ts
  - scripts/lib/b10_acceptance_helpers.py
  - internal/api/handlers/dashboard_handler.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/services/datalink.ts
  - internal/datalink/modbusshare/errors.go
  - internal/api/handlers/test.go
  - internal/api/handlers/typed_errors.go
  - internal/datalink/modbusshare/service_lifecycle.go
  - frontend/src/components/datalink/wizard/LivePreviewPanel.tsx
  - internal/datalink/modbusshare/reconciler_state.go
  - internal/datalink/modbusshare/service.go
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - frontend/src/services/modbusShare.ts
  - internal/datalink/runtime/service.go
  - internal/datalink/runtime/service_projection_state.go
  - internal/api/handlers/modbus_share_handler_reconcile.go
  - scripts/b10_exe_acceptance.py
  - frontend/src/features/datalink/workbench-v2/state/studioV2ShareActivation.ts
  - frontend/src/types/modbusShare.test-d.ts
  - internal/datalink/schema/migrations/001_initial_schema_sqlite.sql
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/datalink/point/service_point_crud.go
  - internal/datalink/runtime/status.go
  - internal/datalink/modbusshare/sql_revision_store.go
  - internal/datalink/sourcerule/output_apply_service.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/types/datalink.ts
  - internal/api/handlers/modbus_share_handler_gates.go
  - internal/datalink/modbusshare/reconciler.go
  - internal/datalink/workspace/service_readiness_connector_issue.go
  - internal/api/handlers/runtime_handler.go
  - docs/swagger/docs.go
  - internal/datalink/modbusshare/service_network_helpers.go
  - internal/datalink/sourcerule/share_desired_mappings.go
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - internal/datalink/sourcerule/candidate_snapshot.go
  - internal/datalink/runtime/service_source_rule_reconcile.go
  - internal/virtual/server/modbus/server_connection.go
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/services/modbusShare.test-d.ts
  - internal/api/handlers/source_rule_handler.go
  - scripts/lib/b10_projection_assertions.py
  - internal/api/handlers/modbus_share_handler.go
  - frontend/src/features/datalink/runtime-dashboard/LiveStateBanner.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - internal/datalink/sourcerule/share_ownership.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - internal/datalink/modbusshare/reconciler_transaction.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - internal/datalink/schema/schema_source_rule_models.go
  - internal/datalink/sourcerule/mutation_rollback.go
  - internal/datalink/sourcerule/interfaces.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ShareOutputSummary.tsx
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsSections.tsx
  - go.mod
  - internal/datalink/sourcerule/mutation_rollback_created.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/utils/typedErrors.ts
  - frontend/src/types/studioV2Workspace.ts
  - internal/api/handlers/response_keys.go
  - internal/datalink/modbusshare/service_write.go
  - internal/datalink/modbusshare/canonical_plan.go
  - internal/datalink/collector/scheduler_dispatch.go
  - frontend/src/types/sourceRuleCandidates.ts
  - internal/api/handlers/settings_handler.go
  - internal/datalink/sourcerule/local_modbus_mapping_reader.go
  - internal/api/handlers/studio_v2_workspace_handler.go
  - internal/datalink/sourcerule/local_modbus_restore.go
  - docs/swagger/swagger.yaml
  - frontend/src/features/datalink/workbench-v2/settings/settingsOperationOwnership.ts
  - internal/api/handlers/source_rule_handler_tag_review_decisions.go
  - internal/api/handlers/datalink_sse_handler.go
  - docs/releases/retire-legacy-studio-and-polish-v2.md
  - internal/datalink/workspace/service_readiness.go
  - internal/datalink/settings/service.go
  - internal/datalink/device/service_status.go
  - frontend/src/types/modbusShare.ts
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/api/handlers/modbus_share_handler_mapping_mutations.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3RuntimeStreams.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/LocalModbusReviewSurface.tsx
  - internal/api/handlers/source_rule_handler_candidates.go
  - internal/datalink/settings/errors.go
  - internal/datalink/sourcerule/candidate_scope.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/api/handlers/source_rule_handler_output_apply.go
  - docs/technical/studio-surface-inventory/CURRENT_STATE.md
  - frontend/src/hooks/datalink/useModbusShareCandidateReview.ts
  - frontend/src/services/sourceRuleTagReviewDecisions.ts
  - frontend/src/hooks/datalink/useModbusShareStatus.ts
  - frontend/src/components/datalink/wizard/MappingWizard.tsx
  - scripts/check_file_lines.sh
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/hooks/usePreviewStream.ts
  - internal/datalink/workspace/service_devices.go
  - frontend/src/components/datalink/wizard/steps/PreviewStep.tsx
  - internal/datalink/dbtarget/service_probe.go
  - internal/api/handlers/source_rule_handler_tag_apply.go
  - internal/datalink/runtime/modbus_share_delivery.go
  - internal/api/handlers/modbus_share_handler_status.go
  - frontend/src/hooks/previewStreamEvents.ts
  - frontend/src/pages/datalink/workbench-v2/studioV2AutosaveBarrier.ts
  - frontend/src/types/studioV2Activation.ts
  - internal/datalink/runtime/truth_state.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/step3LiveSubscription.ts
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells.tsx
  - internal/datalink/modbusshare/reconciler_helpers.go
  - internal/api/router.go
  - internal/api/handlers/transport_wrapper.go
  - internal/datalink/sourcerule/tag_apply_service.go
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - internal/datalink/modbusshare/geometry.go
  - frontend/src/types/runtimeDiagnostics.ts
  - internal/api/handlers/template.go
  - internal/datalink/collector/scheduler_config_methods.go
  - internal/datalink/modbusshare/service_projection.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - internal/datalink/sourcerule/runtime_reconcile.go
  - internal/datalink/modbusshare/settings_lifecycle.go
  - scripts/lib/b10_focused_suite.py
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - scripts/lib/b10_negative_steps.py
  - internal/virtual/server/modbus/server_register_handlers.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
tests:
  - frontend/tests/unit/utils/safeJson.test.ts
  - internal/api/handlers/source_rule_handler_database_candidates_test.go
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - internal/api/handlers/studio_v2_runtime_apply_test.go
  - cmd/test_ui/main.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - cmd/test_ui/share_startup_test.go
  - internal/datalink/modbusshare/reconciler_lifecycle_test.go
  - internal/datalink/modbusshare/settings_lifecycle_cas_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_recovery_regression_test.go
  - cmd/test_ui/target_writer_test.go
  - internal/datalink/modbusshare/service_write_hydration_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-share-activation.test.tsx
  - frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts
  - internal/api/handlers/studio_v2_workspace_activation_contract_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview-changes.test.tsx
  - internal/api/handlers/modbus_share_handler_lifecycle_test.go
  - frontend/tests/unit/workbench-v2/autosave-barrier.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-diagnostics-state.test.tsx
  - internal/api/handlers/source_rule_handler_test.go
  - internal/datalink/device/service_readiness_contract_test.go
  - cmd/test_ui/target_writer.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-truth-state.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - cmd/test_ui/server_runtime.go
  - internal/datalink/sourcerule/candidate_scope_test.go
  - frontend/tests/unit/hooks/usePreviewStream.test.tsx
  - frontend/tests/unit/workbench-v2/settings-operations-race.test.tsx
  - internal/datalink/workspace/service_readiness_connector_issue_test.go
  - internal/virtual/server/modbus/server_test.go
  - internal/datalink/workspace/service_readiness_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route-state.test.tsx
  - frontend/tests/unit/services/modbusShare.test.ts
  - internal/datalink/sourcerule/service_runtime_reconcile_test.go
  - internal/api/handlers/device_handler_extended_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/workbench-v2/commit-progress-accessibility.test.tsx
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/workbench-v2/settings-operation-ownership.test.ts
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-stream-malformed.test.tsx
  - frontend/tests/unit/workbench-v2/workbench-local-modbus-review-surface.test.tsx
  - internal/datalink/modbusshare/settings_hydration_lifecycle_test.go
  - internal/api/handlers/modbus_share_handler_parse_test.go
  - internal/datalink/modbusshare/stale_span_invalidation_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-stream-state.test.tsx
  - internal/api/handlers/runtime_workspace_setup_context_recovery_test.go
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - internal/datalink/collector/scheduler_refactor_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/locale-parity.test.ts
  - frontend/tests/unit/workbench-v2/step4-share-summary.test.tsx
  - internal/api/handlers/mapping_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/datalink/modbusshare/canonical_plan_test.go
  - frontend/tests/unit/workbench-v2/studioV2WorkspaceActivation.test.ts
  - internal/api/router_modbus_share_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.row-groups.test.tsx
  - internal/datalink/modbusshare/reconciler_test.go
  - internal/datalink/sourcerule/mutation_rollback_read_failure_test.go
  - internal/datalink/migrator_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.failure-recovery.test.tsx
  - internal/datalink/modbusshare/status_contract_test.go
  - internal/api/handlers/test_client_operations.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.fixture.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_barrier_test.go
  - frontend/tests/e2e/embedded-frontend-delivery.spec.ts
  - internal/datalink/modbusshare/reconciler_transaction_test.go
  - internal/api/handlers/test_request_helpers_test.go
  - internal/api/handlers/polling_group_handler_test.go
  - internal/api/handlers/studio_v2_runtime_apply_connector_safety_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/api/handlers/runtime_workspace_setup_context_regression_test.go
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.reconciliation.test.tsx
  - internal/datalink/workspace/service_runtime_projection_test.go
  - internal/datalink/modbusshare/reconciler_idempotency_test.go
  - internal/api/handlers/test_connection_handler.go
  - internal/api/handlers/test_script_handler.go
  - frontend/tests/unit/workbench-v2/step4-share-helpers.ts
  - internal/datalink/sourcerule/share_gates_test.go
  - cmd/test_ui/harness_config_test.go
  - cmd/test_ui/share_startup.go
  - internal/datalink/modbusshare/reconciler_empty_revision_b10_test.go
  - cmd/test_ui/harness_config.go
  - internal/api/handlers/modbus_share_handler_diagnostics_test.go
  - internal/datalink/modbusshare/settings_hydration_failure_test.go
  - internal/api/handlers/modbus_share_handler_gates_test.go
  - internal/api/handlers/test_client_factory.go
  - scripts/lib/test_b10_acceptance_helpers.py
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/datalink/modbusshare/geometry_validation_test.go
  - frontend/tests/unit/workbench-v2/step4-share.test.tsx
  - cmd/test_ui/service_wiring.go
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - internal/api/handlers/test_monitor_handler.go
  - internal/api/handlers/datalink_sse_handler_test.go
  - frontend/tests/unit/workbench-v2/step4-database-components.test.tsx
  - internal/datalink/modbusshare/reconciler_nil_store_test.go
  - internal/datalink/sourcerule/share_desired_mappings_test.go
  - internal/api/handlers/runtime_handler_test.go
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/datalink/modbusshare/service_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - cmd/test_ui/share_runtime_reconcile.go
  - internal/api/modbus_share_swagger_contract_test.go
  - internal/datalink/modbusshare/reconciler_validation_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/api/handlers/source_rule_handler_local_modbus_candidates_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - internal/datalink/modbusshare/reconciler_concurrency_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/api/handlers/modbus_share_handler_status_contract_test.go
  - internal/datalink/modbusshare/sql_revision_store_test.go
  - internal/datalink/settings/sql_repo_storage_errors_test.go
  - internal/api/handlers/settings_handler_test.go
  - internal/datalink/sourcerule/repository_memory_test.go
  - internal/datalink/sourcerule/share_restore_projection_test.go
  - scripts/lib/test_check_file_lines.py
  - frontend/tests/unit/runtime-dashboard/runtime-stream-contract.test.tsx
  - frontend/tests/unit/utils/typedErrors.test.ts
  - internal/datalink/runtime/service_source_rule_reconcile_test.go
  - internal/api/handlers/modbus_share_handler_production_gate_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/datalink/settings/sql_repo_cas_test.go
  - internal/api/handlers/typed_errors_contract_test.go
  - internal/datalink/modbusshare/settings_lifecycle_test.go
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/datalink/runtime/ingestor_target_outcomes_test.go
-->

---
### Requirement: Embedded browser smoke proves the full asset graph is loadable

A built embedded binary served with a synchronized static tree SHALL load the canonical and non-target route families without a JavaScript chunk or vendor asset returning HTTP 404. The smoke verification SHALL check network responses and generic unknown-route equivalence for `/studio/v2`, `/studio`, `/studio/runtime`, `/test`, and `/gateway/*`.

#### Scenario: Embedded binary loads all required route families

- **GIVEN** a binary built from a fresh checkout after frontend assets are synchronized
- **WHEN** browser smoke opens the five route families and waits for generic unknown-route handling to settle
- **THEN** every required JavaScript, CSS, and route asset request returns a non-404 response
- **AND** each preserved route identity and `/studio` unknown-route equivalence result matches the product contract
- **AND** the smoke check exits successfully

#### Scenario: Stale or missing asset fails the smoke check

- **GIVEN** one referenced v2, runtime, tool, experimental, or shared asset is absent from `cmd/test_ui/static`
- **WHEN** browser smoke opens the route that imports the missing asset
- **THEN** the check reports the failed asset request
- **AND** the check exits nonzero
- **AND** the missing-asset condition is not hidden by a loading fallback


<!-- @trace
source: retire-legacy-studio-and-polish-v2
updated: 2026-08-26
code:
  - internal/api/handlers/device_health_handler.go
  - internal/datalink/sourcerule/candidate_api.go
  - internal/api/router_modbus_share.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - internal/api/handlers/config.go
  - internal/datalink/sourcerule/candidate_snapshot_local_modbus_outputs.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/vite.config.ts
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - internal/datalink/migrator.go
  - internal/datalink/sourcerule/service_state.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - internal/datalink/modbusshare/service_persistence.go
  - internal/datalink/modbusshare/types.go
  - frontend/src/services/sourceRuleCandidates.ts
  - internal/datalink/sourcerule/validation.go
  - internal/datalink/dbtarget/service.go
  - internal/datalink/workspace/service_readiness_device.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsStatus.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDiagnosticsPanel.tsx
  - internal/virtual/server/modbus/server.go
  - frontend/src/utils/safeJson.ts
  - internal/datalink/collector/scheduler_lifecycle.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - internal/api/handlers/runtime_stream_handler.go
  - internal/datalink/device/service_readiness_safe.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - internal/datalink/runtime/target_delivery.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/sourcerule/share_restore_projection.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/useStep4Activation.ts
  - frontend/src/services/datalinkClient.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/api/handlers/modbus_share_handler_mapping_query.go
  - internal/datalink/sourcerule/repository_memory.go
  - frontend/src/hooks/datalink/useRuntimeStream.ts
  - internal/datalink/settings/sql_repo.go
  - internal/datalink/sourcerule/candidate_snapshot_local_modbus_conflicts.go
  - internal/datalink/sourcerule/service.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - internal/datalink/workspace/service_runtime_projection_scope.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/api/handlers/mapping_handler.go
  - frontend/src/features/datalink/workbench-v2/state/types-settings.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - docs/swagger/swagger.json
  - frontend/src/features/datalink/workbench-v2/settings/useSettingsOperations.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - internal/datalink/runtime/delivery_diagnostic.go
  - internal/api/handlers/debug.go
  - internal/datalink/sourcerule/tag_review_decision_service.go
  - internal/api/handlers/connection.go
  - internal/datalink/modbusshare/sql_revision_store_helpers.go
  - internal/datalink/tag/service_crud.go
  - internal/datalink/sourcerule/candidate_snapshot_revision.go
  - internal/datalink/collector/scheduler.go
  - internal/datalink/schema/migrations/001_initial_schema.up.sql
  - frontend/src/features/datalink/workbench-v2/steps/step4/step4DatabaseHelpers.ts
  - scripts/lib/b10_acceptance_helpers.py
  - internal/api/handlers/dashboard_handler.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/services/datalink.ts
  - internal/datalink/modbusshare/errors.go
  - internal/api/handlers/test.go
  - internal/api/handlers/typed_errors.go
  - internal/datalink/modbusshare/service_lifecycle.go
  - frontend/src/components/datalink/wizard/LivePreviewPanel.tsx
  - internal/datalink/modbusshare/reconciler_state.go
  - internal/datalink/modbusshare/service.go
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - frontend/src/services/modbusShare.ts
  - internal/datalink/runtime/service.go
  - internal/datalink/runtime/service_projection_state.go
  - internal/api/handlers/modbus_share_handler_reconcile.go
  - scripts/b10_exe_acceptance.py
  - frontend/src/features/datalink/workbench-v2/state/studioV2ShareActivation.ts
  - frontend/src/types/modbusShare.test-d.ts
  - internal/datalink/schema/migrations/001_initial_schema_sqlite.sql
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/datalink/point/service_point_crud.go
  - internal/datalink/runtime/status.go
  - internal/datalink/modbusshare/sql_revision_store.go
  - internal/datalink/sourcerule/output_apply_service.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/types/datalink.ts
  - internal/api/handlers/modbus_share_handler_gates.go
  - internal/datalink/modbusshare/reconciler.go
  - internal/datalink/workspace/service_readiness_connector_issue.go
  - internal/api/handlers/runtime_handler.go
  - docs/swagger/docs.go
  - internal/datalink/modbusshare/service_network_helpers.go
  - internal/datalink/sourcerule/share_desired_mappings.go
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - internal/datalink/sourcerule/candidate_snapshot.go
  - internal/datalink/runtime/service_source_rule_reconcile.go
  - internal/virtual/server/modbus/server_connection.go
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/services/modbusShare.test-d.ts
  - internal/api/handlers/source_rule_handler.go
  - scripts/lib/b10_projection_assertions.py
  - internal/api/handlers/modbus_share_handler.go
  - frontend/src/features/datalink/runtime-dashboard/LiveStateBanner.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - internal/datalink/sourcerule/share_ownership.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - internal/datalink/modbusshare/reconciler_transaction.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - internal/datalink/schema/schema_source_rule_models.go
  - internal/datalink/sourcerule/mutation_rollback.go
  - internal/datalink/sourcerule/interfaces.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ShareOutputSummary.tsx
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsSections.tsx
  - go.mod
  - internal/datalink/sourcerule/mutation_rollback_created.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/utils/typedErrors.ts
  - frontend/src/types/studioV2Workspace.ts
  - internal/api/handlers/response_keys.go
  - internal/datalink/modbusshare/service_write.go
  - internal/datalink/modbusshare/canonical_plan.go
  - internal/datalink/collector/scheduler_dispatch.go
  - frontend/src/types/sourceRuleCandidates.ts
  - internal/api/handlers/settings_handler.go
  - internal/datalink/sourcerule/local_modbus_mapping_reader.go
  - internal/api/handlers/studio_v2_workspace_handler.go
  - internal/datalink/sourcerule/local_modbus_restore.go
  - docs/swagger/swagger.yaml
  - frontend/src/features/datalink/workbench-v2/settings/settingsOperationOwnership.ts
  - internal/api/handlers/source_rule_handler_tag_review_decisions.go
  - internal/api/handlers/datalink_sse_handler.go
  - docs/releases/retire-legacy-studio-and-polish-v2.md
  - internal/datalink/workspace/service_readiness.go
  - internal/datalink/settings/service.go
  - internal/datalink/device/service_status.go
  - frontend/src/types/modbusShare.ts
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/api/handlers/modbus_share_handler_mapping_mutations.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3RuntimeStreams.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/LocalModbusReviewSurface.tsx
  - internal/api/handlers/source_rule_handler_candidates.go
  - internal/datalink/settings/errors.go
  - internal/datalink/sourcerule/candidate_scope.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/api/handlers/source_rule_handler_output_apply.go
  - docs/technical/studio-surface-inventory/CURRENT_STATE.md
  - frontend/src/hooks/datalink/useModbusShareCandidateReview.ts
  - frontend/src/services/sourceRuleTagReviewDecisions.ts
  - frontend/src/hooks/datalink/useModbusShareStatus.ts
  - frontend/src/components/datalink/wizard/MappingWizard.tsx
  - scripts/check_file_lines.sh
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/hooks/usePreviewStream.ts
  - internal/datalink/workspace/service_devices.go
  - frontend/src/components/datalink/wizard/steps/PreviewStep.tsx
  - internal/datalink/dbtarget/service_probe.go
  - internal/api/handlers/source_rule_handler_tag_apply.go
  - internal/datalink/runtime/modbus_share_delivery.go
  - internal/api/handlers/modbus_share_handler_status.go
  - frontend/src/hooks/previewStreamEvents.ts
  - frontend/src/pages/datalink/workbench-v2/studioV2AutosaveBarrier.ts
  - frontend/src/types/studioV2Activation.ts
  - internal/datalink/runtime/truth_state.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/step3LiveSubscription.ts
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells.tsx
  - internal/datalink/modbusshare/reconciler_helpers.go
  - internal/api/router.go
  - internal/api/handlers/transport_wrapper.go
  - internal/datalink/sourcerule/tag_apply_service.go
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - internal/datalink/modbusshare/geometry.go
  - frontend/src/types/runtimeDiagnostics.ts
  - internal/api/handlers/template.go
  - internal/datalink/collector/scheduler_config_methods.go
  - internal/datalink/modbusshare/service_projection.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - internal/datalink/sourcerule/runtime_reconcile.go
  - internal/datalink/modbusshare/settings_lifecycle.go
  - scripts/lib/b10_focused_suite.py
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - scripts/lib/b10_negative_steps.py
  - internal/virtual/server/modbus/server_register_handlers.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
tests:
  - frontend/tests/unit/utils/safeJson.test.ts
  - internal/api/handlers/source_rule_handler_database_candidates_test.go
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - internal/api/handlers/studio_v2_runtime_apply_test.go
  - cmd/test_ui/main.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - cmd/test_ui/share_startup_test.go
  - internal/datalink/modbusshare/reconciler_lifecycle_test.go
  - internal/datalink/modbusshare/settings_lifecycle_cas_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_recovery_regression_test.go
  - cmd/test_ui/target_writer_test.go
  - internal/datalink/modbusshare/service_write_hydration_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-share-activation.test.tsx
  - frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts
  - internal/api/handlers/studio_v2_workspace_activation_contract_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview-changes.test.tsx
  - internal/api/handlers/modbus_share_handler_lifecycle_test.go
  - frontend/tests/unit/workbench-v2/autosave-barrier.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-diagnostics-state.test.tsx
  - internal/api/handlers/source_rule_handler_test.go
  - internal/datalink/device/service_readiness_contract_test.go
  - cmd/test_ui/target_writer.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-truth-state.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - cmd/test_ui/server_runtime.go
  - internal/datalink/sourcerule/candidate_scope_test.go
  - frontend/tests/unit/hooks/usePreviewStream.test.tsx
  - frontend/tests/unit/workbench-v2/settings-operations-race.test.tsx
  - internal/datalink/workspace/service_readiness_connector_issue_test.go
  - internal/virtual/server/modbus/server_test.go
  - internal/datalink/workspace/service_readiness_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route-state.test.tsx
  - frontend/tests/unit/services/modbusShare.test.ts
  - internal/datalink/sourcerule/service_runtime_reconcile_test.go
  - internal/api/handlers/device_handler_extended_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/workbench-v2/commit-progress-accessibility.test.tsx
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/workbench-v2/settings-operation-ownership.test.ts
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-stream-malformed.test.tsx
  - frontend/tests/unit/workbench-v2/workbench-local-modbus-review-surface.test.tsx
  - internal/datalink/modbusshare/settings_hydration_lifecycle_test.go
  - internal/api/handlers/modbus_share_handler_parse_test.go
  - internal/datalink/modbusshare/stale_span_invalidation_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-stream-state.test.tsx
  - internal/api/handlers/runtime_workspace_setup_context_recovery_test.go
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - internal/datalink/collector/scheduler_refactor_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/locale-parity.test.ts
  - frontend/tests/unit/workbench-v2/step4-share-summary.test.tsx
  - internal/api/handlers/mapping_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/datalink/modbusshare/canonical_plan_test.go
  - frontend/tests/unit/workbench-v2/studioV2WorkspaceActivation.test.ts
  - internal/api/router_modbus_share_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.row-groups.test.tsx
  - internal/datalink/modbusshare/reconciler_test.go
  - internal/datalink/sourcerule/mutation_rollback_read_failure_test.go
  - internal/datalink/migrator_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.failure-recovery.test.tsx
  - internal/datalink/modbusshare/status_contract_test.go
  - internal/api/handlers/test_client_operations.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.fixture.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_barrier_test.go
  - frontend/tests/e2e/embedded-frontend-delivery.spec.ts
  - internal/datalink/modbusshare/reconciler_transaction_test.go
  - internal/api/handlers/test_request_helpers_test.go
  - internal/api/handlers/polling_group_handler_test.go
  - internal/api/handlers/studio_v2_runtime_apply_connector_safety_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/api/handlers/runtime_workspace_setup_context_regression_test.go
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.reconciliation.test.tsx
  - internal/datalink/workspace/service_runtime_projection_test.go
  - internal/datalink/modbusshare/reconciler_idempotency_test.go
  - internal/api/handlers/test_connection_handler.go
  - internal/api/handlers/test_script_handler.go
  - frontend/tests/unit/workbench-v2/step4-share-helpers.ts
  - internal/datalink/sourcerule/share_gates_test.go
  - cmd/test_ui/harness_config_test.go
  - cmd/test_ui/share_startup.go
  - internal/datalink/modbusshare/reconciler_empty_revision_b10_test.go
  - cmd/test_ui/harness_config.go
  - internal/api/handlers/modbus_share_handler_diagnostics_test.go
  - internal/datalink/modbusshare/settings_hydration_failure_test.go
  - internal/api/handlers/modbus_share_handler_gates_test.go
  - internal/api/handlers/test_client_factory.go
  - scripts/lib/test_b10_acceptance_helpers.py
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/datalink/modbusshare/geometry_validation_test.go
  - frontend/tests/unit/workbench-v2/step4-share.test.tsx
  - cmd/test_ui/service_wiring.go
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - internal/api/handlers/test_monitor_handler.go
  - internal/api/handlers/datalink_sse_handler_test.go
  - frontend/tests/unit/workbench-v2/step4-database-components.test.tsx
  - internal/datalink/modbusshare/reconciler_nil_store_test.go
  - internal/datalink/sourcerule/share_desired_mappings_test.go
  - internal/api/handlers/runtime_handler_test.go
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/datalink/modbusshare/service_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - cmd/test_ui/share_runtime_reconcile.go
  - internal/api/modbus_share_swagger_contract_test.go
  - internal/datalink/modbusshare/reconciler_validation_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/api/handlers/source_rule_handler_local_modbus_candidates_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - internal/datalink/modbusshare/reconciler_concurrency_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/api/handlers/modbus_share_handler_status_contract_test.go
  - internal/datalink/modbusshare/sql_revision_store_test.go
  - internal/datalink/settings/sql_repo_storage_errors_test.go
  - internal/api/handlers/settings_handler_test.go
  - internal/datalink/sourcerule/repository_memory_test.go
  - internal/datalink/sourcerule/share_restore_projection_test.go
  - scripts/lib/test_check_file_lines.py
  - frontend/tests/unit/runtime-dashboard/runtime-stream-contract.test.tsx
  - frontend/tests/unit/utils/typedErrors.test.ts
  - internal/datalink/runtime/service_source_rule_reconcile_test.go
  - internal/api/handlers/modbus_share_handler_production_gate_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/datalink/settings/sql_repo_cas_test.go
  - internal/api/handlers/typed_errors_contract_test.go
  - internal/datalink/modbusshare/settings_lifecycle_test.go
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/datalink/runtime/ingestor_target_outcomes_test.go
-->

---
### Requirement: Playwright webServer startup is cross-platform

The frontend Playwright configuration SHALL provide a webServer command and environment setup that the exact command `npm run test:e2e -- tests/e2e/embedded-frontend-delivery.spec.ts`, launched from the `frontend` directory, can start on Windows and POSIX hosts. The configuration SHALL wait for its declared readiness URL before route assertions begin and SHALL NOT depend on a POSIX-only inline environment assignment or shell-chaining syntax.

#### Scenario: Windows direct invocation starts the configured webServer

- **GIVEN** a Windows checkout with frontend dependencies installed
- **WHEN** the operator runs `npm run test:e2e -- tests/e2e/embedded-frontend-delivery.spec.ts` from `frontend`
- **THEN** Playwright starts the configured webServer without a shell syntax error
- **AND** Playwright waits for the declared readiness URL
- **AND** the embedded frontend smoke proceeds to its route assertions

#### Scenario: POSIX direct invocation starts the configured webServer

- **GIVEN** a POSIX checkout with frontend dependencies installed
- **WHEN** the operator runs `npm run test:e2e -- tests/e2e/embedded-frontend-delivery.spec.ts` from `frontend`
- **THEN** Playwright starts the same configured webServer
- **AND** Playwright waits for the declared readiness URL
- **AND** the embedded frontend smoke proceeds to its route assertions

<!-- @trace
source: reliable-embedded-frontend-delivery
updated: 2026-08-22
code:
  - frontend/playwright.config.ts
tests:
  - frontend/tests/e2e/embedded-frontend-delivery.spec.ts
-->

---
### Requirement: Historical launcher line counts do not regress

The implementation SHALL keep the line count of `start.ps1` and `start.sh` less than or equal to each file's pre-change baseline. The final gate SHALL fail when either historical launcher grows beyond its baseline, while preserving the synchronization and fail-closed behavior requirements.

#### Scenario: Final gate rejects launcher growth

- **GIVEN** pre-change line-count baselines are recorded for `start.ps1` and `start.sh`
- **WHEN** the final content and line-count gate evaluates the repaired launchers
- **THEN** each current line count is less than or equal to its corresponding baseline
- **AND** `scripts/check_file_lines.sh` reports no new hard-limit violation
- **AND** the gate returns nonzero if either current count exceeds its baseline

<!-- @trace
source: reliable-embedded-frontend-delivery
updated: 2026-08-22
code:
  - start.ps1
  - start.sh
  - scripts/check_file_lines.sh
tests:
  - tests/shell/embedded-frontend-delivery.sh
-->

---
### Requirement: Windows embedded EXE build and route smoke are release gates

The Windows release path SHALL run the complete frontend build, static synchronization, and Go embedded executable build through `scripts/build.ps1` with PowerShell profile isolation. A release candidate SHALL produce `bin/test-ui.exe` and SHALL pass a clean embedded-server route and asset smoke before the embedded delivery gate is complete. The smoke SHALL exercise `/studio/v2`, `/studio`, `/studio/runtime`, `/test`, and the existing experimental route while checking that same-origin JavaScript, module, stylesheet, and lazy route asset requests do not fail.

#### Scenario: Complete Windows build produces a loadable embedded graph

- **GIVEN** the repository has its declared frontend and Go toolchains available on Windows
- **WHEN** the release command runs `scripts/build.ps1` from the repository root
- **THEN** frontend typecheck/build completes
- **AND** the synchronized static directory contains the built entry and lazy assets
- **AND** `bin/test-ui.exe` is produced
- **AND** a clean temporary execution directory can start the binary and serve the embedded entry route

#### Scenario: Route and asset smoke covers the supported surface

- **GIVEN** a freshly built `bin/test-ui.exe` is running on a temporary local port
- **WHEN** the embedded-delivery smoke opens `/studio/v2`, `/studio`, `/studio/runtime`, `/test`, and the existing experimental route
- **THEN** each route renders its declared readiness selector and retains its route identity
- **AND** every same-origin script, module, stylesheet, and lazy route asset request returns successfully
- **AND** a deliberately missing asset returns an explicit 404 without being reported as a successful delivery

#### Scenario: Build or synchronization failure blocks delivery

- **GIVEN** frontend build output is missing, static synchronization fails, Go compilation fails, or the embedded server exits before readiness
- **WHEN** the release gate evaluates the command result
- **THEN** the embedded delivery gate is incomplete
- **AND** the captured command output identifies the failed phase
- **AND** no route-smoke pass is claimed

<!-- @trace
source: harden-studio-v2-release-correctness
updated: 2026-08-24
code:
  - scripts/build.ps1
  - frontend/playwright.config.ts
tests:
  - frontend/tests/e2e/embedded-frontend-delivery.spec.ts
-->

---
### Requirement: Final removal preserves non-target embedded delivery

After immediate deletion the embedded asset graph SHALL contain the `/studio/v2`, runtime, test, and gateway route assets required by their contracts and SHALL NOT contain a legacy-only workspace chunk, import, or asset. The `/studio` request SHALL be indistinguishable from an arbitrary unknown route and SHALL NOT load a dedicated route/handler/tombstone, legacy component, or legacy-only asset.

#### Scenario: Legacy asset is absent after removal

- **WHEN** an embedded browser requests `/studio` alongside an arbitrary unknown path after deletion
- **THEN** both requests produce the same response/navigation without a dedicated `/studio` handler or tombstone
- **AND** network evidence shows no legacy-specific redirect logic, legacy workspace component, or legacy-only asset request

#### Scenario: Non-target assets remain available

- **WHEN** an embedded browser opens `/test` and `/gateway/quick-setup` after deletion
- **THEN** their required route chunks and assets load successfully
- **AND** the removal does not redirect or delete those non-target surfaces

<!-- @trace
source: retire-legacy-studio-and-polish-v2
updated: 2026-08-26
code:
  - internal/api/handlers/device_health_handler.go
  - internal/datalink/sourcerule/candidate_api.go
  - internal/api/router_modbus_share.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - internal/api/handlers/config.go
  - internal/datalink/sourcerule/candidate_snapshot_local_modbus_outputs.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/vite.config.ts
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - internal/datalink/migrator.go
  - internal/datalink/sourcerule/service_state.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - internal/datalink/modbusshare/service_persistence.go
  - internal/datalink/modbusshare/types.go
  - frontend/src/services/sourceRuleCandidates.ts
  - internal/datalink/sourcerule/validation.go
  - internal/datalink/dbtarget/service.go
  - internal/datalink/workspace/service_readiness_device.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsStatus.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDiagnosticsPanel.tsx
  - internal/virtual/server/modbus/server.go
  - frontend/src/utils/safeJson.ts
  - internal/datalink/collector/scheduler_lifecycle.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - internal/api/handlers/runtime_stream_handler.go
  - internal/datalink/device/service_readiness_safe.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - internal/datalink/runtime/target_delivery.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/sourcerule/share_restore_projection.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/useStep4Activation.ts
  - frontend/src/services/datalinkClient.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/api/handlers/modbus_share_handler_mapping_query.go
  - internal/datalink/sourcerule/repository_memory.go
  - frontend/src/hooks/datalink/useRuntimeStream.ts
  - internal/datalink/settings/sql_repo.go
  - internal/datalink/sourcerule/candidate_snapshot_local_modbus_conflicts.go
  - internal/datalink/sourcerule/service.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - internal/datalink/workspace/service_runtime_projection_scope.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/api/handlers/mapping_handler.go
  - frontend/src/features/datalink/workbench-v2/state/types-settings.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - docs/swagger/swagger.json
  - frontend/src/features/datalink/workbench-v2/settings/useSettingsOperations.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - internal/datalink/runtime/delivery_diagnostic.go
  - internal/api/handlers/debug.go
  - internal/datalink/sourcerule/tag_review_decision_service.go
  - internal/api/handlers/connection.go
  - internal/datalink/modbusshare/sql_revision_store_helpers.go
  - internal/datalink/tag/service_crud.go
  - internal/datalink/sourcerule/candidate_snapshot_revision.go
  - internal/datalink/collector/scheduler.go
  - internal/datalink/schema/migrations/001_initial_schema.up.sql
  - frontend/src/features/datalink/workbench-v2/steps/step4/step4DatabaseHelpers.ts
  - scripts/lib/b10_acceptance_helpers.py
  - internal/api/handlers/dashboard_handler.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/services/datalink.ts
  - internal/datalink/modbusshare/errors.go
  - internal/api/handlers/test.go
  - internal/api/handlers/typed_errors.go
  - internal/datalink/modbusshare/service_lifecycle.go
  - frontend/src/components/datalink/wizard/LivePreviewPanel.tsx
  - internal/datalink/modbusshare/reconciler_state.go
  - internal/datalink/modbusshare/service.go
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - frontend/src/services/modbusShare.ts
  - internal/datalink/runtime/service.go
  - internal/datalink/runtime/service_projection_state.go
  - internal/api/handlers/modbus_share_handler_reconcile.go
  - scripts/b10_exe_acceptance.py
  - frontend/src/features/datalink/workbench-v2/state/studioV2ShareActivation.ts
  - frontend/src/types/modbusShare.test-d.ts
  - internal/datalink/schema/migrations/001_initial_schema_sqlite.sql
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/datalink/point/service_point_crud.go
  - internal/datalink/runtime/status.go
  - internal/datalink/modbusshare/sql_revision_store.go
  - internal/datalink/sourcerule/output_apply_service.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/types/datalink.ts
  - internal/api/handlers/modbus_share_handler_gates.go
  - internal/datalink/modbusshare/reconciler.go
  - internal/datalink/workspace/service_readiness_connector_issue.go
  - internal/api/handlers/runtime_handler.go
  - docs/swagger/docs.go
  - internal/datalink/modbusshare/service_network_helpers.go
  - internal/datalink/sourcerule/share_desired_mappings.go
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - internal/datalink/sourcerule/candidate_snapshot.go
  - internal/datalink/runtime/service_source_rule_reconcile.go
  - internal/virtual/server/modbus/server_connection.go
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/services/modbusShare.test-d.ts
  - internal/api/handlers/source_rule_handler.go
  - scripts/lib/b10_projection_assertions.py
  - internal/api/handlers/modbus_share_handler.go
  - frontend/src/features/datalink/runtime-dashboard/LiveStateBanner.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - internal/datalink/sourcerule/share_ownership.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - internal/datalink/modbusshare/reconciler_transaction.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - internal/datalink/schema/schema_source_rule_models.go
  - internal/datalink/sourcerule/mutation_rollback.go
  - internal/datalink/sourcerule/interfaces.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ShareOutputSummary.tsx
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsSections.tsx
  - go.mod
  - internal/datalink/sourcerule/mutation_rollback_created.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/utils/typedErrors.ts
  - frontend/src/types/studioV2Workspace.ts
  - internal/api/handlers/response_keys.go
  - internal/datalink/modbusshare/service_write.go
  - internal/datalink/modbusshare/canonical_plan.go
  - internal/datalink/collector/scheduler_dispatch.go
  - frontend/src/types/sourceRuleCandidates.ts
  - internal/api/handlers/settings_handler.go
  - internal/datalink/sourcerule/local_modbus_mapping_reader.go
  - internal/api/handlers/studio_v2_workspace_handler.go
  - internal/datalink/sourcerule/local_modbus_restore.go
  - docs/swagger/swagger.yaml
  - frontend/src/features/datalink/workbench-v2/settings/settingsOperationOwnership.ts
  - internal/api/handlers/source_rule_handler_tag_review_decisions.go
  - internal/api/handlers/datalink_sse_handler.go
  - docs/releases/retire-legacy-studio-and-polish-v2.md
  - internal/datalink/workspace/service_readiness.go
  - internal/datalink/settings/service.go
  - internal/datalink/device/service_status.go
  - frontend/src/types/modbusShare.ts
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/api/handlers/modbus_share_handler_mapping_mutations.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3RuntimeStreams.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/LocalModbusReviewSurface.tsx
  - internal/api/handlers/source_rule_handler_candidates.go
  - internal/datalink/settings/errors.go
  - internal/datalink/sourcerule/candidate_scope.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/api/handlers/source_rule_handler_output_apply.go
  - docs/technical/studio-surface-inventory/CURRENT_STATE.md
  - frontend/src/hooks/datalink/useModbusShareCandidateReview.ts
  - frontend/src/services/sourceRuleTagReviewDecisions.ts
  - frontend/src/hooks/datalink/useModbusShareStatus.ts
  - frontend/src/components/datalink/wizard/MappingWizard.tsx
  - scripts/check_file_lines.sh
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/hooks/usePreviewStream.ts
  - internal/datalink/workspace/service_devices.go
  - frontend/src/components/datalink/wizard/steps/PreviewStep.tsx
  - internal/datalink/dbtarget/service_probe.go
  - internal/api/handlers/source_rule_handler_tag_apply.go
  - internal/datalink/runtime/modbus_share_delivery.go
  - internal/api/handlers/modbus_share_handler_status.go
  - frontend/src/hooks/previewStreamEvents.ts
  - frontend/src/pages/datalink/workbench-v2/studioV2AutosaveBarrier.ts
  - frontend/src/types/studioV2Activation.ts
  - internal/datalink/runtime/truth_state.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/step3LiveSubscription.ts
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells.tsx
  - internal/datalink/modbusshare/reconciler_helpers.go
  - internal/api/router.go
  - internal/api/handlers/transport_wrapper.go
  - internal/datalink/sourcerule/tag_apply_service.go
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - internal/datalink/modbusshare/geometry.go
  - frontend/src/types/runtimeDiagnostics.ts
  - internal/api/handlers/template.go
  - internal/datalink/collector/scheduler_config_methods.go
  - internal/datalink/modbusshare/service_projection.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - internal/datalink/sourcerule/runtime_reconcile.go
  - internal/datalink/modbusshare/settings_lifecycle.go
  - scripts/lib/b10_focused_suite.py
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - scripts/lib/b10_negative_steps.py
  - internal/virtual/server/modbus/server_register_handlers.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
tests:
  - frontend/tests/unit/utils/safeJson.test.ts
  - internal/api/handlers/source_rule_handler_database_candidates_test.go
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - internal/api/handlers/studio_v2_runtime_apply_test.go
  - cmd/test_ui/main.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - cmd/test_ui/share_startup_test.go
  - internal/datalink/modbusshare/reconciler_lifecycle_test.go
  - internal/datalink/modbusshare/settings_lifecycle_cas_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_recovery_regression_test.go
  - cmd/test_ui/target_writer_test.go
  - internal/datalink/modbusshare/service_write_hydration_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-share-activation.test.tsx
  - frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts
  - internal/api/handlers/studio_v2_workspace_activation_contract_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview-changes.test.tsx
  - internal/api/handlers/modbus_share_handler_lifecycle_test.go
  - frontend/tests/unit/workbench-v2/autosave-barrier.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-diagnostics-state.test.tsx
  - internal/api/handlers/source_rule_handler_test.go
  - internal/datalink/device/service_readiness_contract_test.go
  - cmd/test_ui/target_writer.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-truth-state.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - cmd/test_ui/server_runtime.go
  - internal/datalink/sourcerule/candidate_scope_test.go
  - frontend/tests/unit/hooks/usePreviewStream.test.tsx
  - frontend/tests/unit/workbench-v2/settings-operations-race.test.tsx
  - internal/datalink/workspace/service_readiness_connector_issue_test.go
  - internal/virtual/server/modbus/server_test.go
  - internal/datalink/workspace/service_readiness_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route-state.test.tsx
  - frontend/tests/unit/services/modbusShare.test.ts
  - internal/datalink/sourcerule/service_runtime_reconcile_test.go
  - internal/api/handlers/device_handler_extended_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/workbench-v2/commit-progress-accessibility.test.tsx
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/workbench-v2/settings-operation-ownership.test.ts
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-stream-malformed.test.tsx
  - frontend/tests/unit/workbench-v2/workbench-local-modbus-review-surface.test.tsx
  - internal/datalink/modbusshare/settings_hydration_lifecycle_test.go
  - internal/api/handlers/modbus_share_handler_parse_test.go
  - internal/datalink/modbusshare/stale_span_invalidation_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-stream-state.test.tsx
  - internal/api/handlers/runtime_workspace_setup_context_recovery_test.go
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - internal/datalink/collector/scheduler_refactor_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/locale-parity.test.ts
  - frontend/tests/unit/workbench-v2/step4-share-summary.test.tsx
  - internal/api/handlers/mapping_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/datalink/modbusshare/canonical_plan_test.go
  - frontend/tests/unit/workbench-v2/studioV2WorkspaceActivation.test.ts
  - internal/api/router_modbus_share_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.row-groups.test.tsx
  - internal/datalink/modbusshare/reconciler_test.go
  - internal/datalink/sourcerule/mutation_rollback_read_failure_test.go
  - internal/datalink/migrator_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.failure-recovery.test.tsx
  - internal/datalink/modbusshare/status_contract_test.go
  - internal/api/handlers/test_client_operations.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.fixture.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_barrier_test.go
  - frontend/tests/e2e/embedded-frontend-delivery.spec.ts
  - internal/datalink/modbusshare/reconciler_transaction_test.go
  - internal/api/handlers/test_request_helpers_test.go
  - internal/api/handlers/polling_group_handler_test.go
  - internal/api/handlers/studio_v2_runtime_apply_connector_safety_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/api/handlers/runtime_workspace_setup_context_regression_test.go
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.reconciliation.test.tsx
  - internal/datalink/workspace/service_runtime_projection_test.go
  - internal/datalink/modbusshare/reconciler_idempotency_test.go
  - internal/api/handlers/test_connection_handler.go
  - internal/api/handlers/test_script_handler.go
  - frontend/tests/unit/workbench-v2/step4-share-helpers.ts
  - internal/datalink/sourcerule/share_gates_test.go
  - cmd/test_ui/harness_config_test.go
  - cmd/test_ui/share_startup.go
  - internal/datalink/modbusshare/reconciler_empty_revision_b10_test.go
  - cmd/test_ui/harness_config.go
  - internal/api/handlers/modbus_share_handler_diagnostics_test.go
  - internal/datalink/modbusshare/settings_hydration_failure_test.go
  - internal/api/handlers/modbus_share_handler_gates_test.go
  - internal/api/handlers/test_client_factory.go
  - scripts/lib/test_b10_acceptance_helpers.py
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/datalink/modbusshare/geometry_validation_test.go
  - frontend/tests/unit/workbench-v2/step4-share.test.tsx
  - cmd/test_ui/service_wiring.go
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - internal/api/handlers/test_monitor_handler.go
  - internal/api/handlers/datalink_sse_handler_test.go
  - frontend/tests/unit/workbench-v2/step4-database-components.test.tsx
  - internal/datalink/modbusshare/reconciler_nil_store_test.go
  - internal/datalink/sourcerule/share_desired_mappings_test.go
  - internal/api/handlers/runtime_handler_test.go
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/datalink/modbusshare/service_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - cmd/test_ui/share_runtime_reconcile.go
  - internal/api/modbus_share_swagger_contract_test.go
  - internal/datalink/modbusshare/reconciler_validation_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/api/handlers/source_rule_handler_local_modbus_candidates_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - internal/datalink/modbusshare/reconciler_concurrency_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/api/handlers/modbus_share_handler_status_contract_test.go
  - internal/datalink/modbusshare/sql_revision_store_test.go
  - internal/datalink/settings/sql_repo_storage_errors_test.go
  - internal/api/handlers/settings_handler_test.go
  - internal/datalink/sourcerule/repository_memory_test.go
  - internal/datalink/sourcerule/share_restore_projection_test.go
  - scripts/lib/test_check_file_lines.py
  - frontend/tests/unit/runtime-dashboard/runtime-stream-contract.test.tsx
  - frontend/tests/unit/utils/typedErrors.test.ts
  - internal/datalink/runtime/service_source_rule_reconcile_test.go
  - internal/api/handlers/modbus_share_handler_production_gate_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/datalink/settings/sql_repo_cas_test.go
  - internal/api/handlers/typed_errors_contract_test.go
  - internal/datalink/modbusshare/settings_lifecycle_test.go
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/datalink/runtime/ingestor_target_outcomes_test.go
-->