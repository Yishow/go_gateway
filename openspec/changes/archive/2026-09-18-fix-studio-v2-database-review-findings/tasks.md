## 1. 前端錯誤碼與解析收斂

- [x] 1.1 建立 `frontend/src/utils/backendErrorCodes.ts` 作為後端錯誤碼唯一字面值清單（取 `typedErrors.ts` 的 `BACKEND_ERROR_CODES` 與 `safeJson.ts` 的 `SAFE_ERROR_CODES` 聯集，含僅存在於 safeJson 的 `preview_stream_closed`）；`typedErrors.ts` 與 `safeJson.ts` 改由該清單建構各自的集合，行為上兩處可接受的 code 集合不縮小。驗證：更新前端單元測試，斷言兩個集合皆由同一清單建構且含 `preview_stream_closed`；`cd frontend && npm test -- --run` 全綠。
- [x] 1.2 `TestWriteResultStatus`（`frontend/src/types/recordingPlan.ts`）加入 `'unknown'`，`recordingPlanJson.ts` 的 `TEST_WRITE_STATUSES` 與 `parseRecordingTestWriteResult` 接受 `unknown` 並照常解析 `record_id`／`table`／`observed_at` 等欄位。驗證：`recordingPlanJson` 單元測試新增 unknown payload 案例，斷言回傳結果 status 為 `unknown` 且欄位保留。

## 2. unknown 結果的 UI 呈現

- [x] 2.1 [after: 1.2, 3.1] step4 試寫結果呈現對 `unknown` 狀態顯示「無法確認是否寫入」（新 i18n key，zh-TW 與 en 皆提供），不出現成功樣式、不顯示錯誤樣式。驗證：step4 元件測試以 `status: 'unknown'` 的結果渲染，斷言出現無法確認文案且成功／錯誤樣式不出現。

## 3. i18n 文案收斂

- [x] 3.1 於 `frontend/src/i18n/locales/zh-TW/workbench-v2.json` 與 `en/workbench-v2.json` 補上 `step4.conflict_tooltip`，並使 `step4.target_remote_conflict` 的呼叫端預設值與 locale 文字一致（以 locale 為準）；至少一個 step4 測試改用真實 i18n 實例（不再 mock 成回傳 key），使缺失 key 會讓測試失敗。驗證：該真實 i18n 測試在 key 齊備時通過；暫時移除 key 可使其失敗（手動驗證後還原）。
- [x] 3.2 [after: 3.1] 將 step4 衝突／錯誤相關元件（`TargetMappingRow.tsx`、`TargetMappingTable.tsx`、`TargetMetadataStatus.tsx`）的 inline 中文 defaultValue 收斂為 locale key（zh-TW 與 en 同步補齊），呼叫端僅傳 key。驗證：`npm test -- --run` 全綠，且以 rg 樣式 `t\('[^']+',\s*'[^']+'` 檢查該批檔案無殘留 inline defaultValue。
- [x] 3.3 [after: 3.2] 將 step4 其餘元件（`ActivationNeutralSummary.tsx`、`CommitProgress.tsx`、`CommitSuccessCard.tsx`、`CommitSummary.tsx`、`LocalModbusReviewSurface.tsx`、`SchemaSetupSection.tsx`、`RecordingPlanSetupSection.tsx` 等）的 inline defaultValue 分批完成同樣收斂。驗證：`npm run lint && npm test -- --run && npm run build` 全綠，step4 目錄不再有 inline 中文 defaultValue 樣式。

## 4. 後端錯誤呈現與建表判定

- [x] 4.1 於 `internal/api/handlers/studio_v2_workspace_database_errors.go` 檔內抽出共用 render helper，`renderStudioV2WorkspaceDatabaseError` 與 `renderStudioV2WorkspaceSchemaEnsureError` 僅保留各自 status／code／訊息差異，envelope 統一為 `TypedAPIErrorEnvelope`，對外回應語義不變。驗證：`go test ./internal/api/...` 全綠（既有 handler 測試為回歸閘）。
- [x] 4.2 重寫 `EnsureWorkspaceSchema`（`internal/api/handlers/studio_v2_workspace_database_handler.go`）判定：dry-run 語句為空且 managed 目標表 inspection（`internal/datalink/dbtarget` 的 `TableInspection`）全為 `exists` ⇒ 回 nil；語句非空或 inspection 有 `missing` ⇒ 回 `errWorkspaceSchemaPreparationRequired`；inspection 為 `forbidden`／`failed`／無法確認 ⇒ 回安全錯誤。驗證：新增 handler 測試涵蓋四種狀態各至少一例，並回歸 `studio_v2_workspace_activation_handler` 啟動閘測試；`go test ./internal/api/... ./internal/datalink/...` 全綠。
- [x] 4.3 `internal/datalink/recordingplan/schema_apply.go` 中 `:192` 的 inline 字串 `"check the target tables, then preview again"` 改為使用 `schemaNextActionCheck` 常數或語義正確的具名常數，行為輸出不變。驗證：`go test ./internal/datalink/recordingplan/...` 全綠。

## 5. 命名釐清與測試檔區分

- [x] 5.1 [after: 3.2] `recordingPlanJson.ts` 的 `optionalBoundedText` 改名 `optionalShortText`（揭露 256 字元上限語義）、`boundedRecordingText` 補「16KB 上限（`MAX_RECORDING_STATEMENT_LENGTH`）」doc 註解；`TargetMappingRow.tsx` 的 `witnessIndex` 改名 `settledCommitIndex`、`isPendingWitness` 改名 `isCommitSettled`（語義：已回報 commit 在 pending 佇列中的位置，用於截斷已確定項目）。行為不變，僅命名與註解。驗證：`npm run lint && npm test -- --run` 全綠。
- [x] 5.2 將 `frontend/tests/unit/workbench-v2/step4-target-metadata.test.ts`（單元測試）改名為 `step4-target-metadata-parsing.test.ts`，與同名元件測試 `.test.tsx` 區分。驗證：`npm test -- --run` 全綠且兩檔皆被執行。

## 6. 前綴同步防護

- [x] 6.1 新增前端單元測試斷言 `useRecordingPlanActions.ts` 的 `MANAGED_TABLE_PREFIX === 'gw_record_'`；前端常數註解指向後端 `defaultManagedTablePrefix`，後端該常數宣告處反向指回前端路徑。驗證：新增測試通過；`npm test -- --run` 與 `go test ./internal/datalink/...` 全綠。
