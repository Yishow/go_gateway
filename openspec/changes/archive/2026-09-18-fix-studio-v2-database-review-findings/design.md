## Context

2026-09-18 對暫存之 Studio v2 step4 資料庫工作（staged diff vs `87ee5392`）的雙軸審查產生本 change。已核實的現況事實：

- `TargetMappingRow.tsx:224` 使用 `t('step4.conflict_tooltip', '欄位與其他點位重複')`，但 zh-TW 與 en 兩份 `workbench-v2.json` 皆無此 key；`step4-database.test.tsx:143` 因 mock `t` 而掩蓋問題。step4 目錄 `*.tsx` 以樣式 `t\('[^']+',\s*'[^']+'` 計數約 99 處 inline 中文 defaultValue。
- `frontend/src/types/recordingPlan.ts:141` 的 `TestWriteResultStatus` 僅含 `written_verified | written_unverified | failed`；`recordingPlanJson.ts:75` 的 `parseRecordingTestWriteResult` 將 `unknown` 判為非法回應。
- `EnsureWorkspaceSchema`（`studio_v2_workspace_database_handler.go:130`）僅以 dry-run 語句數判定：語句非空即回 `errWorkspaceSchemaPreparationRequired`；啟動流程經 `schemaEnsurer` 介面（`studio_v2_workspace_activation_handler.go:22`、`:133`）呼叫。
- `typedErrors.ts` 的 `BACKEND_ERROR_CODES` 與 `safeJson.ts` 的 `SAFE_ERROR_CODES` 內容重複且已漂移：後端 `typed_errors.go` 的 `ErrCodePreviewStreamClosed = "preview_stream_closed"` 只存在於 safeJson 清單。
- `internal/datalink/recordingplan/schema_apply.go:51` 已有常數 `schemaNextActionCheck`，`:192` 卻 inline 異體字串 `"check the target tables, then preview again"`。
- `recordingPlanJson.ts` 的 `optionalBoundedText`（256 字元上限 `boundedString`）與 `boundedRecordingText`（16KB 上限）命名未揭露界線；`TargetMappingRow.tsx:59-71` 的 `witnessIndex`／`isPendingWitness` 指向 pendingCommits 佇列中已回報項目的位置，需註解才能理解。
- `frontend/tests/unit/workbench-v2/step4-target-metadata.test.ts`（單元）與 `step4-target-metadata.test.tsx`（元件）同名異質。
- `useRecordingPlanActions.ts:23` 的 `MANAGED_TABLE_PREFIX = 'gw_record_'` 與後端 `defaultManagedTablePrefix` 僅靠註解同步；後端經 `table_prefix` 參數接受前端值並驗證，無 canonical 前綴暴露管道。

限制：不得修改任何 spec／delta spec；不得改變 route identity；前端 `npm run lint && npm test && npm run build` 與後端 `go test ./... && go vet ./... && golangci-lint run ./...` 基準必須維持。

## Goals / Non-Goals

Goals：

1. i18n 文案收斂至 locale 檔，補齊缺失 key，消除英文介面漏中文。
2. 前端接受並顯示 `unknown` 試寫結果狀態。
3. `EnsureWorkspaceSchema` 以 inspection status 佐證，不再單靠 dry-run 語句數誤判。
4. 錯誤碼清單單一來源化；Go 錯誤呈現去重；`schema_apply.go` 常數對齊。
5. 命名釐清與測試檔改名；前綴同步防護測試。

Non-Goals：見 proposal「Non-Goals」；另加：不導入伺服端 canonical 前綴機制、不製作文案抽取工具（僅手工分批收斂）。

## Decisions

1. **錯誤碼單一來源**：新增 `frontend/src/utils/backendErrorCodes.ts`，匯出 `BACKEND_ERROR_CODES`（`as const` 陣列，唯一字面值清單）；`typedErrors.ts` 與 `safeJson.ts` 皆由此匯入建構各自集合。清單內容取兩檔聯集（含 `preview_stream_closed`）。替代方案「safeJson 直接 import typedErrors」否決：safeJson 是更底層 util，反向依賴語意錯置。
2. **Go 錯誤呈現去重**：於 `studio_v2_workspace_database_errors.go` 檔內抽出一個共用 render helper（非匯出 API），`renderStudioV2WorkspaceDatabaseError` 與 `renderStudioV2WorkspaceSchemaEnsureError` 僅提供各自的 status／code／訊息差異；envelope 統一為 `TypedAPIErrorEnvelope`。既有 handler 測試為回歸閘。
3. **unknown 對齊**：`TestWriteResultStatus` 加入 `'unknown'`；`parseRecordingTestWriteResult` 接受 `unknown` 並照常解析其餘欄位；UI 對 `unknown` 顯示「無法確認是否寫入」（非成功、非錯誤），對齊 `fix-studio-v2-database-workflow` delta spec 情境「Write outcome or status is unknown」。不新增 route、不接線仍回 501 的後端試寫。
4. **EnsureWorkspaceSchema 佐證**：判定改為——dry-run 語句為空且 managed 目標表 inspection 全為 `exists` ⇒ 就緒（nil）；語句非空或 inspection 有 `missing` ⇒ `errWorkspaceSchemaPreparationRequired`；inspection 為 `forbidden`／`failed`／無法確認 ⇒ 回傳安全錯誤（不得誤判 ready 或 preparation required）。inspection 重用 `internal/datalink/dbtarget` 的 `TableInspection`（`exists/missing/forbidden/failed`，`table_inspection.go:21-24`）路徑，由 handler 既有 connector service 依賴接線；驗收以 handler 測試釘選四種狀態。
5. **i18n 收斂**：locale 檔為唯一文案來源，呼叫處保留 key、移除對 inline defaultValue 的依賴。分批進行：第一批衝突／錯誤相關文案（含 `conflict_tooltip`），其後一般文案。`step4.target_remote_conflict` 以 locale 文字為準。
6. **命名**：`optionalBoundedText` 改名 `optionalShortText`（揭露 256 上限語義）；`boundedRecordingText` 補「16KB 上限」doc 註解；`witnessIndex` 改名 `settledCommitIndex`、`isPendingWitness` 改名 `isCommitSettled`（語義：該筆已回報 commit 在 pending 佇列中的位置，用於截斷已確定項目；實作時若發現更貼切的行為名可調整，但不得保留 witness 隱喻）。測試檔：`.test.ts` 改名 `step4-target-metadata-parsing.test.ts`，`.test.tsx` 維持元件測試名。
7. **前綴同步防護**：新增前端單元測試斷言 `MANAGED_TABLE_PREFIX === 'gw_record_'`；前端常數註解指向後端 `defaultManagedTablePrefix`，後端該常數處反向指回前端路徑。

## Implementation Contract

- 能力級可觀察行為不變：除 `unknown` 顯示與文案鍵修正外，所有 API 回應、路由、狀態碼不變。
- 驗收命令：`cd frontend && npm run lint && npm test -- --run && npm run build`；`go test ./... && go vet ./... && golangci-lint run ./...`。
- locale key 缺失由「以真實 i18n 實例（非 mock 回傳 key）」的前端測試捕捉。
- `EnsureWorkspaceSchema` 四態驗收：語句空＋全 exists ⇒ nil；語句非空或有 missing ⇒ preparation required；forbidden／failed／無法確認 ⇒ 安全錯誤；每態至少一個 handler 測試。
- 範圍邊界——In scope：proposal「What Changes」八組；Out of scope：proposal「Non-Goals」。

## Risks / Trade-offs

- [約 99 處文案搬移觸及大量 step4 元件] -> 分批小步，每批跑 vitest；先補真實 i18n 實例測試再搬移。
- [EnsureWorkspaceSchema 收緊後，原本誤判 ready 的情境將改報錯] -> 以四態 handler 測試釘選行為，並回歸啟動閘（activation）測試。
- [改名觸及多檔] -> 限 step4 與 utils 範圍，rename 後立即 lint+test。
- [錯誤碼清單聯集可能放行新代碼] -> 僅加入有後端來源依據的 `preview_stream_closed`，其餘逐一比對後端 `typed_errors.go`。

## Migration Plan

純前端與 handler 內部重構：無資料遷移、無部署順序依賴；單一變更可整體 `git revert` 回滾。

## Open Questions

無。`EnsureWorkspaceSchema` 的 inspection 接線細節由實作決定，以四態 handler 測試驗收。
