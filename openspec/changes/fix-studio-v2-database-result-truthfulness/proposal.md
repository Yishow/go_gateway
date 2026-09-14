# 停止資料庫流程的假成功

## Why
基準 main 為 `1a0311c8e8db9c62fe4388f0701ba38afe552ff7`。記錄方案入口的 `SchemaApply`、`TestWrite` 在沒有操作目標資料庫時回傳成功；畫面又使用不一致的成功名稱。這會讓操作者誤判資料已安全送達。來源與上一輪 F01/F02/F05 對照見本批總覽的 evidence.md。

## What Changes
- 尚未接上真實操作時，建表、試寫各自明確回覆未開放，禁止假成功及自動重試。
- 依實際可用功能提供能力旗標，連線成功不代表建表、試寫或正式交付成功。
- 統一安全錯誤格式與未知結果處理；不得將逾時等同未執行，也不得等同成功。
- 取代既有 Step 4 規格中以十段計時動畫宣告成功的過時條文；保留真正的 runtime 導航。

## Capabilities
### New Capabilities
None.
### Modified Capabilities
- `recording-database-setup`: 新增未實作操作、安全錯誤與未知結果的規範。
- `datalink-workbench-v2-step4-database`: 更新 Commit sequence and animation、Commit completion card，成功須有伺服器證據。

## Impact
實作定位：`internal/api/handlers/studio_v2_workspace_recording_plans_handler.go`、`internal/datalink/recordingplan/service.go`、`frontend/src/types/recordingPlan.ts`、`frontend/src/features/datalink/workbench-v2/steps/step4/RecordingPlanSetupSection.tsx`、`CommitProgress.tsx`、`CommitSuccessCard.tsx`，及對應測試、語系與 API 文件。檔名縮寫以同一 step4 目錄為定位；tasks 使用完整路徑。

## Non-goals
不執行真正建表或試寫，不改資料庫結構，不關掉一般資料庫已存在的有效功能，不重建 `/studio`，不導入上一輪未完整驗證的 patch。

## Dependencies and Delivery
無前置 change，必須先於 C3/C4 開放真實功能。C3/C4 的實作與驗收逐項完成後，才移除相應操作的未開放狀態，不能整包一次解鎖。本次僅起草；全部 implementation tasks 保持未完成。
