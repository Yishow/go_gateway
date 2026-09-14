# Implementation Tasks

以下全為待執行；本次只提交規格。每項以先補失敗測試、修正、整理為順序。規格中的每個 scenario 都必須有可追溯測試或人工驗收紀錄。

- [ ] 1.1 在 `internal/api/handlers/` 補 router 層測試，證明格式合法的任意 token 與 plan_id 不可取得未實作成功；涵蓋 body 缺失與原有驗證失敗。
- [ ] 1.2 [after: 1.1] 修改 `internal/api/handlers/studio_v2_workspace_recording_plans_handler.go`，讓兩個尚未接線的操作回傳各自固定 501 code；驗證沒有成功 payload 與任何資料庫副作用。
- [ ] 1.3 [after: 1.2] 在 `internal/datalink/recordingplan/service.go` 與相關 API 能力回覆加入真實可用操作判定；測試查詢可用但試寫不可用的組合。
- [ ] 1.4 [after: 1.3] 在 `frontend/src/types/recordingPlan.ts` 及 `frontend/src/features/datalink/workbench-v2/steps/step4/RecordingPlanSetupSection.tsx` 統一 typed results，501 不自動重試、未知值不成功；補 Testing Library 測試。
- [ ] 1.5 [after: 1.4] 在 `frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx`、`CommitSuccessCard.tsx` 與呼叫端移除由時間推導成功的顯示；測試部分成功與 runtime 導航仍可用。
- [ ] 1.6 [after: 1.5] 更新 `frontend/src/i18n/locales/` 與適用的 `docs/swagger/` 產物，列出安全訊息、501 及未知結果；驗證不出現原始資料庫錯誤或憑證。
- [ ] 1.7 [after: 1.6] 執行相關 Go/前端測試，再依 `AGENTS.md` 執行完整測試、lint、build；記錄一般 database 建表與 Local Modbus-only 無回歸、未執行項目及原因。
- [ ] 1.8 [after: 1.7] 重新核對 main，驗證本 change 規格、scenario-to-test 對照、release note 與安全 rollback；不得將「501 安全封鎖完成」記成「真實建表／試寫已完成」。
