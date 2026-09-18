## Why

2026-09-18 對暫存之 Studio v2 step4 資料庫工作（staged diff vs `87ee5392`）的雙軸程式碼審查發現多項缺陷：i18n 缺 key 使英文介面漏出 hardcoded 中文、規格要求的 `unknown` 試寫結果狀態被前端型別與 parser 丟棄、建表就緒判定以單一啟發式可能誤報、錯誤碼清單與 Go 錯誤呈現重複、命名與測試檔名混淆。這些程式碼尚未提交，此刻修復成本最低。

## What Changes

- i18n 修復：於 zh-TW 與 en locale 補上缺失的 `step4.conflict_tooltip`；將 step4 元件中約 99 處 inline 中文 `defaultValue`（以 rg 樣式 `t\('[^']+',\s*'[^']+'` 於 `frontend/src/features/datalink/workbench-v2/steps/step4/*.tsx` 計數，2026-09-18 量測）分批收斂至 locale；修正 `step4.target_remote_conflict` 預設值與 locale 文字不一致。
- `unknown` 結果狀態對齊：`TestWriteResultStatus` 加入 `unknown`、`parseRecordingTestWriteResult` 接受 `unknown`，UI 顯示「無法確認是否寫入」而非視為非法回應，對齊 `fix-studio-v2-database-workflow` delta spec 情境「Write outcome or status is unknown」；僅修正狀態模型，不接線試寫。
- 建表就緒判定補強：`EnsureWorkspaceSchema` 除 dry-run 語句數外，加入 table inspection status 佐證，避免 `GenerateSchema` 回傳非缺表語句時誤報 `WORKSPACE_SCHEMA_PREPARATION_REQUIRED`。
- 錯誤碼清單收斂：`typedErrors.ts` 的 `BACKEND_ERROR_CODES` 與 `safeJson.ts` 的 `SAFE_ERROR_CODES` 合併為單一來源（含 `RECORDING_*` 與 `modbus_*` 碼）。
- Go 錯誤呈現收斂：`studio_v2_workspace_database_errors.go` 的 `renderStudioV2WorkspaceDatabaseError` 與 `renderStudioV2WorkspaceSchemaEnsureError` 去除重複（ErrValidation case、500 訊息），envelope 統一為 `TypedAPIErrorEnvelope`。
- 常數對齊：`internal/datalink/recordingplan/schema_apply.go` 中 inline 的 `"check the target tables, then preview again"` 改用既有 `schemaNextActionCheck` 常數或語義正確的具名常數。
- 前綴同步防護：為前端 `MANAGED_TABLE_PREFIX`（`gw_record_`）新增單元測試鎖定值，並在前後端兩處註解互相引用後端 `defaultManagedTablePrefix`；伺服端持有 canonical 前綴的改造不在本次範圍。
- 命名釐清：`optionalBoundedText`（256 字元上限）改名以揭露界線；`boundedRecordingText`（16KB 上限）補文件註解；`TargetMappingRow` 的 `witnessIndex`／`isPendingWitness` 以行為命名重構；同名異質測試檔 `step4-target-metadata.test.ts`（單元）與 `step4-target-metadata.test.tsx`（元件）改名區分。

## Non-Goals (optional)

- 不接線 test-write preview/apply 與回讀清理（`fix-studio-v2-database-workflow` tasks 3.3／3.4 範圍）。
- 不實作 custom-table 建表、metadata 自動配對完整驗收、引導介面與整合／瀏覽器 e2e（同 change 既有未完成任務）。
- 不修改任何 spec、delta spec 或 canonical specs；本 change 為 no-spec。
- 不處理「workspace/generic 建表入口唯讀 vs task 3.8 語意」與「canonical specs 重寫程序矛盾」，留待 `$spectra-ingest` 更新既有 change。
- 不改變任何 API 路由識別，不新增 `/studio` surface。

## Impact

- Affected specs: none
- Affected code:
  - New: 無新增子系統；僅新增單元測試與 locale key
  - Modified: `frontend/src/features/datalink/workbench-v2/steps/step4/`（TargetMappingRow.tsx、useRecordingPlanActions.ts 及 step4 元件）、`frontend/src/i18n/locales/zh-TW/workbench-v2.json`、`frontend/src/i18n/locales/en/workbench-v2.json`、`frontend/src/utils/typedErrors.ts`、`frontend/src/utils/safeJson.ts`、`frontend/src/utils/recordingPlanJson.ts`、`frontend/src/types/recordingPlan.ts`、`frontend/tests/unit/workbench-v2/`、`internal/api/handlers/studio_v2_workspace_database_errors.go`、`internal/api/handlers/studio_v2_workspace_database_handler.go`、`internal/datalink/recordingplan/schema_apply.go`
  - Removed: 無
- Compatibility: `unknown` 狀態由「非法回應」改為顯示「無法確認」屬 UI 呈現對齊既有規格情境；其餘為內部收斂與文案修正，無 capability-level 可觀察行為變更。
