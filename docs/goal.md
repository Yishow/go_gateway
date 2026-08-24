# Studio V2 /goal 清單
> **目前 route contract（2026-08-24）**：`/studio/v2` 是唯一產品 setup 入口，`/studio/runtime`、`/test`、`/gateway/*` 維持各自 route identity；legacy `/studio` dedicated surface 已獲 owner 授權立即刪除，刪除後依 generic unknown-route policy 處理。Git rollback 只恢復 tracked source/docs 與重新建置可恢復的 assets，不恢復 runtime/deployment data、外部 DB/設備狀態、bookmark 或 localStorage。下方較早的 `/studio` 直達要求屬 pre-delete 歷史證據，不能當成現行契約。
這份清單是給 Codex CLI `/goal` 用的執行入口。目的不是重述 spec，而是把每個 change 都改寫成：

- 先讀什麼
- 只能改哪裡
- 做完要拿出哪些證據
- 哪些情況必須停下來，不准自稱完成

總目標：

- `/studio/v2` 全站對接，讓 operator 可見的主產品 surface 都以真實後端契約運作。
- 除了純開發輔助面之外，不接受 mock、noop、或 local-only fallback 冒充正式完成。
- 若 workspace 為空，必須呈現真實空狀態，而不是靠 seeded draft data 撐頁。

整體 Not complete if：

- Step 1 `執行測試` 仍是 mock animation 而非真實後端診斷。
- `/studio/v2/settings` 的 connector test 或 save bar 仍是 mock / noop。
- 空 workspace 仍靠前端預設 device / rule 開頁。
- 任何 operator 會依賴的 `/studio/v2` 主線互動仍只更新本地 state，不寫入或不讀取真實後端。

執行規則：

1. 一次只跑一個 `/goal`
2. 前一個 goal 沒貼出驗收證據，不開下一個
3. final summary 必須列出：
   - 修改檔案
   - 新增或調整的測試名稱
   - 執行命令與結果
   - 手動 demo 步驟與結果
   - 尚未驗證的風險

建議順序：

1. `make-studio-v2-default-entry`
2. `add-studio-v2-single-workspace-foundation`
3. `wire-studio-v2-step1-device-autosave`
4. `wire-studio-v2-step2-rule-autosave`
5. `wire-studio-v2-step3-mapping-autosave`
6. `wire-studio-v2-step4-database-autosave`
7. `apply-running-studio-v2-device-updates`
8. `stop-invalid-studio-v2-devices`
9. `replace-step4-with-first-activation`
10. `add-studio-runtime-workspace-device-switching`
11. `replace-step1-mock-test-with-live-diagnostics`
12. `wire-studio-v2-settings-backend`

## 1. make-studio-v2-default-entry

```text
/goal Strictly implement openspec/changes/make-studio-v2-default-entry and nothing broader.

First action: read these files and report counts before editing:
  - openspec/changes/make-studio-v2-default-entry/proposal.md
  - openspec/changes/make-studio-v2-default-entry/design.md
  - openspec/changes/make-studio-v2-default-entry/tasks.md
  - openspec/changes/make-studio-v2-default-entry/specs/datalink-workbench-v2-shell/spec.md
  - AGENTS.md
  - CLAUDE.md
Report: task count, SHALL count, affected route paths, and files expected to change.

Scope: frontend/src/App.tsx, frontend/src/features/datalink/legacyRoutes.ts, frontend/tests/unit/workbench-v2/routing.test.tsx, and only the minimum adjacent route wiring required by this change.

Constraints:
  - Do not modify Step 1-4 behavior, workspace behavior, runtime behavior, or any backend API.
  - **（歷史要求，已被 immediate legacy deletion supersede）** Preserve direct /studio access; do not force /studio to redirect to /studio/v2.
  - Do not change task-specific legacy deep links outside generic landing routes.
  - Follow AGENTS.md and CLAUDE.md exactly.
  - Do not add dependencies.

Done when:
  1. / and catch-all fallback both route to /studio/v2 in the actual route tree.
  2. Generic legacy landing routes /datalink and /datalink/workbench route to /studio/v2.
  3. **（歷史 pre-delete evidence，已失效）** Direct /studio still opened the legacy workbench path and was covered by routing tests.
  4. frontend/tests/unit/workbench-v2/routing.test.tsx covers /, unknown route, /datalink, /datalink/workbench, and /studio; cite exact test names in the final summary.
  5. Run cd frontend && npm run test -- --run frontend/tests/unit/workbench-v2/routing.test.tsx and paste the summary.
  6. Run git diff --check and confirm the diff is limited to the route-entry surface.

Stop if:
  - Preserving direct /studio requires changing runtime, workspace, or Step 1-4 logic.
  - Existing routing tests fail for reasons outside this route matrix change.
  - The change would require modifying unrelated legacy deep links.
  - A new dependency or package.json change would be required.

```

## 2. add-studio-v2-single-workspace-foundation

```text
/goal Strictly implement openspec/changes/add-studio-v2-single-workspace-foundation as the singleton persisted workspace foundation for /studio/v2.

First action: read these files and report counts before editing:
  - openspec/changes/add-studio-v2-single-workspace-foundation/proposal.md
  - openspec/changes/add-studio-v2-single-workspace-foundation/design.md
  - openspec/changes/add-studio-v2-single-workspace-foundation/tasks.md
  - openspec/changes/add-studio-v2-single-workspace-foundation/specs/studio-v2-workspace/spec.md
  - openspec/changes/add-studio-v2-single-workspace-foundation/specs/datalink-api/spec.md
  - AGENTS.md
  - CLAUDE.md
Report: task count, SHALL count, required API fields, and expected backend/frontend files.

Scope: internal/datalink/workspace/, internal/api/router.go, internal/api/handlers/studio_v2_workspace_handler.go, frontend/src/services/studioV2Workspace.ts, frontend/src/hooks/datalink/useStudioV2Workspace.ts, frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx, frontend/src/types/datalink.ts, and directly related tests.

Constraints:
  - Only one workspace; do not add multi-workspace switching.
  - Do not import legacy /studio data into the new workspace.
  - Do not implement Step 1-4 autosave or runtime activation in this goal.
  - Do not add dependencies or broad schema refactors outside the workspace foundation.
  - Follow AGENTS.md and CLAUDE.md exactly.

Done when:
  1. GET /api/v1/datalink/studio-v2/workspace auto-creates and returns a singleton workspace with at least id, kind, status, ordered_device_ids.
  2. The workspace persists across reloads and service restart; cite the exact backend integration test names that prove this.
  3. Bootstrap never returns a half-built workspace on failure and has an explicit error-path test.
  4. /studio/v2 boots from the workspace API instead of local-only default state; cite the exact frontend boot test names.
  5. Run the relevant backend tests with go test ./... for the touched packages and paste the package-level summary.
  6. Run the relevant frontend tests with cd frontend && npm run test -- --run and paste the summary for the workspace boot tests.
  7. When the singleton workspace has no devices or rules, /studio/v2 renders a real empty state instead of seeded local draft device/rule data.

Stop if:
  - Workspace persistence requires silently pulling in legacy /studio devices or rules.
  - Implementing bootstrap would force Step 1-4 autosave or runtime activation in the same goal.
  - Existing backend or frontend tests fail outside the workspace foundation surface.
  - go.mod or package.json changes would be required.
  - Empty-workspace boot still depends on seeded draft state to keep the page usable.

```

## 3. wire-studio-v2-step1-device-autosave

```text
/goal Strictly implement openspec/changes/wire-studio-v2-step1-device-autosave so Step 1 becomes workspace-scoped valid-only autosave with per-device isolation.

First action: read these files and report counts before editing:
  - openspec/changes/wire-studio-v2-step1-device-autosave/proposal.md
  - openspec/changes/wire-studio-v2-step1-device-autosave/design.md
  - openspec/changes/wire-studio-v2-step1-device-autosave/tasks.md
  - openspec/changes/wire-studio-v2-step1-device-autosave/specs/datalink-workbench-v2-step1-device/spec.md
  - openspec/changes/wire-studio-v2-step1-device-autosave/specs/datalink-api/spec.md
  - AGENTS.md
  - CLAUDE.md
Report: task count, SHALL count, API endpoints to add or change, and expected frontend/backend test areas.

Scope: frontend/src/features/datalink/workbench-v2/steps/step1/, frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts, frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx, frontend/src/services/studioV2Devices.ts, frontend/src/hooks/datalink/useStudioV2Devices.ts, internal/api/router.go, internal/api/handlers/device_handler.go, internal/datalink/device/, internal/datalink/workspace/, and directly related tests.

Constraints:
  - Only valid device content may autosave.
  - Invalid device edits must stay on screen and must not overwrite the last persisted version.
  - Per-device partial success is required; one device failure must not block other valid devices.
  - Persisted device deletion must remove related V2 data.
  - Do not implement Step 2-4, runtime activation, or workspace foundation changes beyond what this autosave needs.

Done when:
  1. Step 1 create/update is valid-only autosave with exact frontend state coverage for draft-invalid, saving, saved, and save-error.
  2. A manual demo path is documented in the final summary: create two devices where A is valid and B is invalid, show that A persists, B stays unsaved, and B's input is not cleared.
  3. Workspace device order persists across reload; cite the exact backend service test and frontend reload test names.
  4. Persisted device deletion cascades through V2-related data; cite the exact backend cascade test names.
  5. Run the relevant backend tests with go test ./... for touched packages and paste the package-level summary.
  6. Run the relevant frontend tests with cd frontend && npm run test -- --run for Step 1 and routing/integration coverage and paste the summary.

Stop if:
  - Implementing device autosave would require changing Step 2-4 contracts in the same goal.
  - Invalid input would be lost from the UI after save failure.
  - Existing tests start failing and the only apparent fix is changing or skipping unrelated tests.
  - New dependencies would be required.

```

## 4. wire-studio-v2-step2-rule-autosave

```text
/goal Strictly implement openspec/changes/wire-studio-v2-step2-rule-autosave so Step 2 rules become workspace-scoped valid-only autosave with per-rule isolation.

First action: read these files and report counts before editing:
  - openspec/changes/wire-studio-v2-step2-rule-autosave/proposal.md
  - openspec/changes/wire-studio-v2-step2-rule-autosave/design.md
  - openspec/changes/wire-studio-v2-step2-rule-autosave/tasks.md
  - openspec/changes/wire-studio-v2-step2-rule-autosave/specs/datalink-workbench-v2-step2-rule/spec.md
  - openspec/changes/wire-studio-v2-step2-rule-autosave/specs/datalink-api/spec.md
  - AGENTS.md
  - CLAUDE.md
Report: task count, SHALL count, rule ownership boundaries, and expected tests.

Scope: frontend/src/features/datalink/workbench-v2/steps/step2/, frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts, frontend/src/services/studioV2Rules.ts, frontend/src/hooks/datalink/useStudioV2Rules.ts, internal/api/router.go, internal/api/handlers/source_rule_handler.go, internal/datalink/sourcerule/, internal/datalink/workspace/, and directly related tests.

Constraints:
  - Only valid rule content may autosave.
  - Invalid rules must not overwrite the last persisted version.
  - Rules must stay bound to the correct workspace_id and device_id.
  - One rule failure must not block other valid rules.
  - Do not implement mapping, database target, or activation behavior in this goal.

Done when:
  1. Step 2 rule autosave is valid-only and per-rule isolated, with exact state coverage for saving, saved, and save-error.
  2. Ownership mismatch is explicitly rejected; cite the exact handler or service test name.
  3. A manual demo path is documented in the final summary: two rules where one is valid and one is invalid, showing partial success without cross-device drift.
  4. Reload preserves workspace/device ownership of rules; cite the exact frontend reload test and backend service test names.
  5. Run the relevant backend tests with go test ./... for touched packages and paste the package-level summary.
  6. Run the relevant frontend tests with cd frontend && npm run test -- --run for Step 2 coverage and paste the summary.

Stop if:
  - Rule autosave requires changing mapping or runtime behavior in the same goal.
  - A rule save can reattach data to a different device without an explicit error.
  - Existing tests fail and the only path forward is editing unrelated tests.
  - New dependencies would be required.

```

## 5. wire-studio-v2-step3-mapping-autosave

```text
/goal Strictly implement openspec/changes/wire-studio-v2-step3-mapping-autosave so Step 3 mappings become workspace-scoped valid-only autosave with per-row isolation.

First action: read these files and report counts before editing:
  - openspec/changes/wire-studio-v2-step3-mapping-autosave/proposal.md
  - openspec/changes/wire-studio-v2-step3-mapping-autosave/design.md
  - openspec/changes/wire-studio-v2-step3-mapping-autosave/tasks.md
  - openspec/changes/wire-studio-v2-step3-mapping-autosave/specs/datalink-workbench-v2-step3-mapping/spec.md
  - openspec/changes/wire-studio-v2-step3-mapping-autosave/specs/datalink-api/spec.md
  - AGENTS.md
  - CLAUDE.md
Report: task count, SHALL count, save-state fields, and expected frontend/backend tests.

Scope: frontend/src/features/datalink/workbench-v2/steps/step3/, frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts, frontend/src/services/studioV2Mappings.ts, frontend/src/hooks/datalink/useStudioV2Mappings.ts, internal/api/router.go, internal/api/handlers/mapping_handler.go, internal/datalink/mapping/, internal/datalink/workspace/, and directly related tests.

Constraints:
  - Only valid mapping rows may autosave.
  - Invalid rows must keep local input and must not overwrite persisted values.
  - Row-level partial success is required.
  - Frontend must preserve and show local_value, persisted_value, and save_state distinctly.
  - Do not implement database-target autosave or runtime behavior in this goal.

Done when:
  1. Step 3 autosave is valid-only and row-isolated; cite the exact Step 3 test and row editor/state test names.
  2. Reload still distinguishes unsaved local value from persisted value for a row.
  3. A manual demo path is documented in the final summary: two rows where one persists and one fails while retaining local input.
  4. Backend validation cleanly splits valid and invalid mapping rows; cite the exact service or handler test names.
  5. Run the relevant backend tests with go test ./... for touched packages and paste the package-level summary.
  6. Run the relevant frontend tests with cd frontend && npm run test -- --run for Step 3 coverage and paste the summary.

Stop if:
  - Mapping autosave would require changing Step 4 or activation flow in the same goal.
  - Save failure clears a user's local row input.
  - Existing tests fail and the only path forward is muting or rewriting unrelated tests.
  - New dependencies would be required.

```

## 6. wire-studio-v2-step4-database-autosave

```text
/goal Strictly implement openspec/changes/wire-studio-v2-step4-database-autosave so Step 4 database config and target rows become valid-only autosave while remaining pre-activation.

First action: read these files and report counts before editing:
  - openspec/changes/wire-studio-v2-step4-database-autosave/proposal.md
  - openspec/changes/wire-studio-v2-step4-database-autosave/design.md
  - openspec/changes/wire-studio-v2-step4-database-autosave/tasks.md
  - openspec/changes/wire-studio-v2-step4-database-autosave/specs/datalink-workbench-v2-step4-database/spec.md
  - openspec/changes/wire-studio-v2-step4-database-autosave/specs/datalink-api/spec.md
  - AGENTS.md
  - CLAUDE.md
Report: task count, SHALL count, pre-activation boundary, and expected tests.

Scope: frontend/src/features/datalink/workbench-v2/steps/step4/, frontend/src/features/datalink/workbench-v2/state/dbReducer.ts, frontend/src/services/studioV2DatabaseTargets.ts, frontend/src/hooks/datalink/useStudioV2DatabaseTargets.ts, internal/api/router.go, internal/api/handlers/db_target_handler.go, internal/datalink/dbtarget/, internal/datalink/workspace/, and directly related tests.

Constraints:
  - Only valid connector and target content may autosave.
  - Target rows must save independently.
  - Autosave success must not start runtime or scheduler behavior.
  - Metadata save state and target-row save state must remain separate.
  - Do not implement first activation in this goal.

Done when:
  1. Step 4 autosave is valid-only and per-target isolated; cite the exact Step 4 and integration test names.
  2. A manual demo path is documented in the final summary: valid connector metadata plus one valid target row and one invalid target row, showing partial success without auto-start.
  3. There is an explicit test proving autosave success does not change runtime/scheduler state.
  4. Frontend state tests prove metadata success does not hide target-row failure.
  5. Run the relevant backend tests with go test ./... for touched packages and paste the package-level summary.
  6. Run the relevant frontend tests with cd frontend && npm run test -- --run for Step 4 coverage and paste the summary.

Stop if:
  - Implementing autosave requires activation or runtime handoff changes in the same goal.
  - Autosave success starts runtime implicitly.
  - Existing tests fail and the only apparent fix is changing unrelated tests.
  - New dependencies would be required.

```

## 7. apply-running-studio-v2-device-updates

```text
/goal Strictly implement openspec/changes/apply-running-studio-v2-device-updates so valid autosaves apply directly to already-running devices, while not-running devices remain dormant.

First action: read these files and report counts before editing:
  - openspec/changes/apply-running-studio-v2-device-updates/proposal.md
  - openspec/changes/apply-running-studio-v2-device-updates/design.md
  - openspec/changes/apply-running-studio-v2-device-updates/tasks.md
  - openspec/changes/apply-running-studio-v2-device-updates/specs/studio-v2-live-config-apply/spec.md
  - openspec/changes/apply-running-studio-v2-device-updates/specs/datalink-api/spec.md
  - AGENTS.md
  - CLAUDE.md
Report: task count, SHALL count, runtime_apply_status values, and expected tests.

Scope: internal/datalink/runtime/, internal/datalink/device/, internal/datalink/sourcerule/, internal/datalink/mapping/, internal/datalink/dbtarget/, internal/api/handlers/, frontend/src/hooks/datalink/, frontend/src/features/datalink/workbench-v2/, frontend/src/types/studioV2RuntimeApply.ts, and directly related tests.

Constraints:
  - Already-running devices must apply valid autosaves directly.
  - Not-running devices must remain not running after autosave.
  - Autosave responses must expose runtime_apply_status with not_running, applied, and apply_failed.
  - UI must not treat apply_failed as simple save success.
  - Do not implement invalid-stop behavior in this goal.

Done when:
  1. Backend runtime tests prove running devices apply valid autosaves without pressing Step 4 again; cite exact test names.
  2. Backend service tests prove not-running devices stay dormant after autosave; cite exact test names.
  3. Step 1-4 autosave success payloads include runtime_apply_status and frontend hook tests recognize all three values.
  4. Frontend state tests preserve explicit apply-failure visibility instead of flattening it into save success.
  5. Run the relevant backend tests with go test ./... for touched packages and paste the package-level summary.
  6. Run the relevant frontend tests with cd frontend && npm run test -- --run for runtime-apply coverage and paste the summary.

Stop if:
  - Direct apply for running devices would require changing first-activation semantics in this goal.
  - A not-running device starts running because of autosave alone.
  - Existing tests fail and the only path forward is muting or rewriting unrelated tests.
  - New dependencies would be required.

```

## 8. stop-invalid-studio-v2-devices

```text
/goal Strictly implement openspec/changes/stop-invalid-studio-v2-devices so invalid edits immediately stop running devices and expose explicit unavailable status.

First action: read these files and report counts before editing:
  - openspec/changes/stop-invalid-studio-v2-devices/proposal.md
  - openspec/changes/stop-invalid-studio-v2-devices/design.md
  - openspec/changes/stop-invalid-studio-v2-devices/tasks.md
  - openspec/changes/stop-invalid-studio-v2-devices/specs/studio-v2-device-validity/spec.md
  - openspec/changes/stop-invalid-studio-v2-devices/specs/datalink-api/spec.md
  - AGENTS.md
  - CLAUDE.md
Report: task count, SHALL count, availability fields, and expected tests.

Scope: internal/datalink/runtime/, internal/datalink/device/, internal/api/handlers/, frontend/src/features/datalink/workbench-v2/, frontend/src/features/datalink/runtime-dashboard/, frontend/src/types/studioV2Availability.ts, and directly related tests.

Constraints:
  - Invalid edit on a running device must stop it immediately.
  - No silent last-good fallback is allowed.
  - availability_status and availability_reason must be explicit in payloads.
  - Unavailable devices must remain visible in collections; they must not disappear.
  - Do not implement workspace runtime page switching in this goal.

Done when:
  1. Backend runtime tests prove a running device transitions from running to stopped/unavailable on invalid edit; cite exact test names.
  2. Backend service tests prove invalidation does not continue reporting the device as available.
  3. Workspace and runtime payloads expose availability_status and availability_reason, covered by handler tests and frontend type tests.
  4. Frontend state tests prove unavailable devices remain visible rather than being dropped.
  5. Run the relevant backend tests with go test ./... for touched packages and paste the package-level summary.
  6. Run the relevant frontend tests with cd frontend && npm run test -- --run for availability-state coverage and paste the summary.

Stop if:
  - Implementing invalid-stop would require the workspace runtime page rewrite in the same goal.
  - The only way to pass is to hide unavailable devices from the UI.
  - Existing tests fail and the only apparent fix is changing unrelated tests.
  - New dependencies would be required.

```

## 9. replace-step4-with-first-activation

```text
/goal Strictly implement openspec/changes/replace-step4-with-first-activation so Step 4 becomes first activation only, with per-device results and partial-failure navigation.

First action: read these files and report counts before editing:
  - openspec/changes/replace-step4-with-first-activation/proposal.md
  - openspec/changes/replace-step4-with-first-activation/design.md
  - openspec/changes/replace-step4-with-first-activation/tasks.md
  - openspec/changes/replace-step4-with-first-activation/specs/datalink-workbench-v2-step4-database/spec.md
  - openspec/changes/replace-step4-with-first-activation/specs/datalink-api/spec.md
  - AGENTS.md
  - CLAUDE.md
Report: task count, SHALL count, eligible-device rule, activation API contract, and expected tests.

Scope: frontend/src/features/datalink/workbench-v2/steps/step4/, frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx, internal/api/router.go, internal/api/handlers/runtime_handler.go, internal/datalink/runtime/, internal/datalink/device/, frontend/src/types/studioV2Activation.ts, and directly related tests.

Constraints:
  - Step 4 must activate only valid, available, not-yet-running devices in the singleton workspace.
  - Activation must return per-device success/failure results.
  - At least one successful activation must still allow navigation to runtime.
  - No eligible devices must produce an actionable empty result instead of fake success.
  - Do not reintroduce commit-style data persistence into Step 4.

Done when:
  1. Backend activation tests prove Step 4 does not re-activate already-running devices; cite exact test names.
  2. Handler tests prove POST /api/v1/datalink/studio-v2/workspace/activate returns per-device results under partial success.
  3. Frontend Step 4 tests prove the button and result cards are activation-oriented, not commit-oriented.
  4. Frontend integration tests prove partial activation failure does not block navigation when at least one device succeeds.
  5. A manual demo path is documented in the final summary: at least two eligible devices where one starts and one fails, followed by successful navigation to runtime.
  6. Run the relevant backend tests with go test ./... and the relevant frontend tests with cd frontend && npm run test -- --run, then paste both summaries.

Stop if:
  - Activation requires redesigning runtime workspace switching in the same goal.
  - Step 4 needs to persist data again to appear successful.
  - Existing tests fail and the only apparent fix is changing unrelated tests.
  - New dependencies would be required.

```

## 10. add-studio-runtime-workspace-device-switching

```text
/goal Strictly implement openspec/changes/add-studio-runtime-workspace-device-switching so /studio/runtime becomes workspace-scoped first with device switching, unavailable visibility, and empty-state staying on page.

First action: read these files and report counts before editing:
  - openspec/changes/add-studio-runtime-workspace-device-switching/proposal.md
  - openspec/changes/add-studio-runtime-workspace-device-switching/design.md
  - openspec/changes/add-studio-runtime-workspace-device-switching/tasks.md
  - openspec/changes/add-studio-runtime-workspace-device-switching/specs/studio-runtime-workspace-view/spec.md
  - openspec/changes/add-studio-runtime-workspace-device-switching/specs/datalink-api/spec.md
  - AGENTS.md
  - CLAUDE.md
Report: task count, SHALL count, runtime context payload fields, default selection rule, and expected tests.

Scope: frontend/src/features/datalink/runtime-dashboard/, frontend/src/App.tsx, frontend/src/services/studioV2RuntimeContext.ts, frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts, internal/api/router.go, internal/api/handlers/runtime_handler.go, internal/datalink/runtime/, internal/datalink/workspace/, and directly related tests.

Constraints:
  - /studio/runtime must load the singleton workspace device set first, not rely solely on one device_id.
  - Default selection must follow V2 order and choose the first available device.
  - Unavailable devices must remain visible with reasons.
  - If there are no available devices, stay on the runtime page and show an empty state; do not redirect.
  - device_id query may remain as an optional override, not the only source of truth.

Done when:
  1. Backend runtime-context tests prove GET /api/v1/datalink/studio-v2/workspace/runtime-context returns ordered devices plus default_device_id; cite exact test names.
  2. Frontend route boot tests prove runtime loads from workspace context without requiring device_id.
  3. Frontend runtime tests prove default device follows V2 order rather than event timing.
  4. Frontend UI tests prove unavailable devices remain visible and empty runtime state stays on page without redirect.
  5. A manual demo path is documented in the final summary: one available device, one unavailable device, and one empty-state case.
  6. Run the relevant backend tests with go test ./... and the relevant frontend tests with cd frontend && npm run test -- --run, then paste both summaries.

Stop if:
  - Implementing runtime workspace view requires multi-workspace support.
  - The only way to pass is to hide unavailable devices or force redirect on empty state.
  - Existing tests fail and the only apparent fix is changing unrelated tests.
  - New dependencies would be required.

```

## 11. replace-step1-mock-test-with-live-diagnostics

```text
/goal Strictly implement openspec/changes/replace-step1-mock-test-with-live-diagnostics so Step 1 `執行測試` uses real backend diagnostics instead of mock animation.

First action: read these files and report counts before editing:
  - openspec/changes/replace-step1-mock-test-with-live-diagnostics/proposal.md
  - openspec/changes/replace-step1-mock-test-with-live-diagnostics/design.md
  - openspec/changes/replace-step1-mock-test-with-live-diagnostics/tasks.md
  - openspec/changes/replace-step1-mock-test-with-live-diagnostics/specs/datalink-workbench-v2-step1-device/spec.md
  - openspec/changes/replace-step1-mock-test-with-live-diagnostics/specs/datalink-api/spec.md
  - AGENTS.md
  - CLAUDE.md
Report: task count, SHALL count, chosen diagnostics endpoint(s), and expected frontend/backend test areas.

Scope: frontend/src/features/datalink/workbench-v2/steps/step1/, frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts, frontend/src/hooks/datalink/, frontend/src/services/, internal/api/router.go, internal/api/handlers/device_handler.go, internal/datalink/device/, and directly related tests.

Constraints:
  - Step 1 test must use a real backend request for the current draft or persisted device diagnostics.
  - The UI must reflect backend-produced connect / probe outcomes, not local timers or Math.random.
  - Test failure must remain visible and must not mark the device as tested.
  - Do not broaden this goal into Step 2-4 autosave, first activation, or settings backend wiring.

Done when:
  1. Clicking `執行測試` issues a real backend request instead of a mock animation; cite the exact frontend test name.
  2. Backend tests prove success and failure paths return diagnostics V2 can render without synthesizing fake stages.
  3. Frontend Step 1 tests prove success, connect failure, and probe failure all remain actionable and do not auto-pass the continue gate.
  4. The implementation removes operator-visible setTimeout/Math.random success synthesis from Step 1.
  5. Run the relevant backend tests with go test ./... and the relevant frontend tests with cd frontend && npm run test -- --run, then paste both summaries.

Stop if:
  - The only available implementation path is keeping the mock animation and merely relabeling it.
  - The chosen diagnostics contract would require redesigning workspace autosave semantics in the same goal.
  - Existing tests fail and the only apparent fix is changing unrelated tests.
  - New dependencies would be required.

```

## 12. wire-studio-v2-settings-backend

```text
/goal Strictly implement openspec/changes/wire-studio-v2-settings-backend so /studio/v2/settings boots from backend state, persists edits, and uses real connector testing.

First action: read these files and report counts before editing:
  - openspec/changes/wire-studio-v2-settings-backend/proposal.md
  - openspec/changes/wire-studio-v2-settings-backend/design.md
  - openspec/changes/wire-studio-v2-settings-backend/tasks.md
  - openspec/changes/wire-studio-v2-settings-backend/specs/datalink-workbench-v2-settings/spec.md
  - openspec/changes/wire-studio-v2-settings-backend/specs/datalink-api/spec.md
  - AGENTS.md
  - CLAUDE.md
Report: task count, SHALL count, reused backend endpoints, and expected frontend/backend test areas.

Scope: frontend/src/features/datalink/workbench-v2/settings/, frontend/src/features/datalink/workbench-v2/state/settingsReducer.ts, frontend/src/hooks/datalink/, frontend/src/services/, internal/api/router.go, internal/api/handlers/settings_handler.go, internal/api/handlers/dbtarget_handler.go, internal/datalink/settings/, internal/datalink/dbtarget/, and directly related tests.

Constraints:
  - Settings page must read initial values from backend, not only in-memory defaults.
  - Connector pool CRUD/test must use real backend APIs; no Math.random, no mock timer.
  - `儲存所有設定` must persist or return actionable failure; it cannot remain a noop.
  - Do not redesign the settings visual layout in this goal.

Done when:
  1. /studio/v2/settings boots from backend data and frontend tests cover the initial load contract.
  2. Connector pool create/update/delete/test all use real backend APIs and have frontend plus backend coverage.
  3. Clicking `儲存所有設定` persists changed settings through backend APIs and no longer logs a noop warning.
  4. Operator-visible mock/noop behavior is removed from the settings surface.
  5. Run the relevant backend tests with go test ./... and the relevant frontend tests with cd frontend && npm run test -- --run, then paste both summaries.

Stop if:
  - The only path forward is leaving connector test or save bar as mock/noop.
  - Existing backend contracts are insufficient and the change would silently drop user settings.
  - Existing tests fail and the only apparent fix is changing unrelated tests.
  - New dependencies would be required.

```
