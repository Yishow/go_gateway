# Implementation Tasks

## 1. 安全封鎖與真實完成證據

2026-09-15 已依使用者授權開始 C1 實作，1.1–1.8 已完成本機驗證。每項以先補失敗測試、修正、整理為順序。規格中的每個 scenario 都必須有可追溯測試或人工驗收紀錄。最終版本的測試、審查、發布與回復證據見 [validation.md](validation.md)。

- [x] 1.1 在 `internal/api/` 沿用既有 recording-plan router fixture 補測試，證明格式合法的任意 token 與 plan_id 不可取得未實作成功；涵蓋 body 缺失與原有驗證失敗。 對照：Unimplemented recording operations fail closed；設計「先修後端，再修畫面」。
- [x] 1.2 [after: 1.1] 修改 `internal/api/handlers/studio_v2_workspace_recording_plans_handler.go`，讓兩個尚未接線的操作回傳各自固定 501 code；router 測試 assert 完整 error envelope 的 code、safe message、retryable=false、action=`wait_for_supported_operation` 與完整 request_id，並驗證沒有成功 payload 與任何資料庫副作用。
- [x] 1.3 [after: 1.2] 在 `internal/datalink/recordingplan/service.go` 與相關 API 能力回覆加入真實可用操作判定；全域 mask 僅取 adapter kind 能力與已接線 operation 的交集，兩個未接線操作固定為 false，connector scope 留給 workflow／C3；測試查詢可用但試寫不可用的組合，避免擴充 scope API。 對照：Capability flags describe connected implementations；設計「能力以「目前可執行的操作」為準」。
- [x] 1.4 [after: 1.3] 在 `frontend/src/types/recordingPlan.ts`、`frontend/src/utils/safeJson.ts`、`frontend/src/utils/typedErrors.ts`、對應 parser tests 及 `frontend/src/features/datalink/workbench-v2/steps/step4/RecordingPlanSetupSection.tsx` 統一 typed results，501 不自動重試、未知值不成功；測試 normalize／translation 後保留精確 code、action、retryable=false、完整 request_id，且 raw diagnostics 不進 DOM。 對照：Safe errors and uncertain outcomes remain distinguishable；設計「成功狀態與網路結果分離」。
- [x] 1.5 [after: 1.4] 在 `frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx`、`useStep4Activation.ts`、`CommitProgress.tsx`、`CommitSuccessCard.tsx` 與既有 workspace/runtime 讀取呼叫端沿用現有後端 activation 資料流，修正缺少證據的成功文案並取代主規格的固定動畫要求；remount recovery 讀既有後端 workspace/runtime 狀態，不重送 activation，無本次 revision 完成證據保持未確認，operation_id 只在既有回覆提供時保留，不新建 activation ledger 或重寫 runtime；以 BackendBackedActivationResults 測試部分成功、全部失敗、原有設備仍運作、重載、empty results 且無交付證據時不得宣稱儲存正常；runtime 導航依既有政策獨立驗證；保留 workspace/settings revision、readiness token 及 Modbus Share gate，不將本次未確認設定標成成功。 對照：Commit sequence and animation；Commit completion card；設計「完成卡與 runtime 交接」。
- [x] 1.6 [after: 1.5] 更新 `frontend/src/i18n/locales/` 與適用的 `docs/swagger/` 產物，列出安全訊息、501 及未知結果；驗證不出現原始資料庫錯誤或憑證。
- [x] 1.7 [after: 1.6] 執行相關 Go/前端測試，再依 `AGENTS.md` 執行完整測試、lint、build；記錄一般 database 建表與 Local Modbus-only 無回歸、未執行項目及原因。
- [x] 1.8 [after: 1.7] 重新核對 main，驗證本 change 規格、scenario-to-test 對照、release note 與安全 rollback；不得將「501 安全封鎖完成」記成「真實建表／試寫已完成」。

### 最終驗證狀態

續作已完成 Go 基線 lint 維護：1,110 → 0 issues。最後版本的 Go test／vet／lint、前端 test／lint／build、相關競態檢查與本機 TCP／UDP 採集寫入 smoke 均通過。已重新核對本機 main、17 個 scenario 對照、release note 與安全 rollback，1.7–1.8 完成。詳細範圍、5 項既有 Go skip、未執行的現場／E2E 驗收及兩批獨立 checkpoint 見 validation.md；未提交或封存。
