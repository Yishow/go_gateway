# 停止資料庫流程的假成功

## Why
基準 main 為 `1a0311c8e8db9c62fe4388f0701ba38afe552ff7`。記錄方案入口的 `SchemaApply`、`TestWrite` 在沒有操作目標資料庫時回傳成功；畫面又使用不一致的成功名稱。這會讓操作者誤判資料已安全送達。目前可核對的來源見 [workflow evidence](../../fix-studio-v2-database-workflow/evidence.md) 的 E01／E03／E04；舊 F01/F02/F05 附件不在 repo，不能當作已讀證據。

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
實作定位：`internal/api/handlers/studio_v2_workspace_recording_plans_handler.go`、`internal/datalink/recordingplan/service.go`、`frontend/src/types/recordingPlan.ts`、`frontend/src/utils/safeJson.ts`、`frontend/src/utils/typedErrors.ts`、對應 parser tests、`frontend/src/features/datalink/workbench-v2/steps/step4/RecordingPlanSetupSection.tsx`、`Step4Database.tsx`、`useStep4Activation.ts`、`CommitProgress.tsx`、`CommitSuccessCard.tsx`、既有 workspace/runtime 讀取呼叫端，及對應測試、語系與 API 文件。檔名縮寫以同一 step4 目錄為定位；tasks 使用完整路徑。

## Non-goals
不執行真正建表或試寫，不改資料庫結構，不關掉一般資料庫已存在的有效功能，不重建 `/studio`，不導入上一輪未完整驗證的 patch。

## Dependencies and Delivery
無前置 change；本案唯一負責假成功、能力／錯誤及進度／完成卡。`fix-studio-v2-database-workflow` 先完成其 1.2、1.3、2.1、2.2、2.4，才由 `implement-studio-v2-verified-schema-setup` 開放通過驗證的建表；試寫直到 `fix-studio-v2-database-workflow` 3.3、3.4 通過才解鎖。完整跨案順序見 [workflow design](../../fix-studio-v2-database-workflow/design.md) 的 Implementation Contract。不依賴未發布的 C2/C4 附件。

## Implementation Checkpoint

2026-09-15 已完成 C1 的 1.1–1.8 本機實作與驗證；HEAD 與本機 main 均為 `569da9f98e52b2635638d85250a4af16601d4364`。續作的 Go 基線維護已解除完整 lint 阻擋，詳細結果與範圍分界見 [validation.md](validation.md)。本紀錄不代表真正建表／試寫完成；另外兩案仍未實作，未提交或封存。
