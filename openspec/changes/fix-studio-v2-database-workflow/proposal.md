# 修正 Studio v2 資料庫流程與操作介面

## Why

第 4 步必須讓使用者知道「資料存到哪裡、是否真的建立、是否真的寫入」，而不是只看到成功提示。本次依 main 的實際程式起草：記錄方案的建表／試寫入口直接回覆成功；目標連線寫死；多設備成員全部套第一台設備；示範欄位被說成真實欄位；輸入與預覽也有更新互相覆蓋的問題。來源、確定程度與尚未完成的檢查見 evidence.md。

基準：Yishow/go_gateway，main @ 1a0311c8e8db9c62fe4388f0701ba38afe552ff7。分支：openspec/fix-studio-v2-database-workflow。日期：2026-09-15（Asia/Taipei）。

初稿只交付提案、設計、需求與待辦。2026-09-15 已另獲授權開始實作，依下列分階段交付；目前實際完成狀態以 tasks.md 與 validation.md 為準，不套用未驗證的上一輪修補包，不更動正式資料或啟用設備。

## Dependencies and Delivery

本案與 `fix-studio-v2-database-result-truthfulness`、`implement-studio-v2-verified-schema-setup` 共同交付。每項需求只有一個實作及 delta spec 負責者；完整順序與舊待辦移交見 design.md 的 Implementation Contract。

1. `fix-studio-v2-database-result-truthfulness` 完成安全封鎖、能力／錯誤及進度／完成卡。
2. 本案先完成 1.2、1.3、2.1、2.2、2.4，交付可核對版本的已存連線、成員與原子儲存。
3. `implement-studio-v2-verified-schema-setup` 完成真正 metadata、preview、operation ledger 與所有建表入口的確認保護。
4. 本案才接續 2.3、3.3、3.4、4.1、4.3，交付配對、試寫／清理及引導驗收。

## What Changes

- 沿用 `fix-studio-v2-database-result-truthfulness` 的安全封鎖及能力／錯誤契約；本案不重做其 handler 或完成卡。
- 使用已儲存的真實連線、測量項目及其設備歸屬；方案要明確選取，不以第一筆清單或猜測編號代替。
- 讀取真正資料表與欄位；建議範本、待確認配對及已確認配對分開，保留合法的同筆資料分組規則。
- 消費 `implement-studio-v2-verified-schema-setup` 的真實欄位、版本綁定預覽及建表結果；沿用其持久操作機制擴充試寫預覽、明確確認及狀態查詢，再實作試寫、讀回及清理，未知結果不盲目重做。
- 設定儲存不能只成功一半卻顯示已完成；新增故障情境測試。
- 保留正在編輯的文字，Enter 只送一次，中文選字不提前送出，舊預覽不更新新畫面，唯讀時不能建表或試寫。
- 第 4 步改成「存到哪裡 → 使用現成表或建立新表 → 資料放哪個欄位 → 檢查並開始」，補上易懂文案、錯誤、鍵盤操作與窄畫面。
- 整合 `fix-studio-v2-database-result-truthfulness` 以後端為準的進度與完成卡，保留設備啟用與 Modbus Share 保護。

## Capabilities

### New Capabilities

無。沿用既有能力，不另建平行資料庫設定入口。

### Modified Capabilities

- `recording-database-setup`：補齊憑證／目標與設備歸屬、真實試寫及一致儲存的可驗收規則。
- `datalink-workbench-v2-step4-database`：補齊真實欄位的配對、分組衝突、編輯／預覽保護與引導操作。

## Impact

`workspace-database-row-groups` 的合法共用欄位規則保留，不改成所有同名欄位一律禁止。`openspec/specs/` 主規格本次不直接改寫；差異規格留在本 change 供後續驗證與同步。

實作允許範圍（實際修改及驗證依 tasks.md／validation.md 記錄）：

- `internal/api/handlers/studio_v2_workspace_recording_plans_handler.go`
- `internal/api/router_studio_v2_recording_routes.go` 與 schema change 提供的 operation status handler
- `internal/api/handlers/studio_v2_workspace_database_handler.go`
- `internal/api/handlers/studio_v2_workspace_database_handler_helpers.go`
- `internal/datalink/recordingplan/`、`internal/datalink/dbtarget/`、`internal/datalink/workspace/`
- `frontend/src/features/datalink/workbench-v2/steps/step4/`
- `frontend/src/features/datalink/workbench-v2/state/`
- `frontend/src/hooks/datalink/useStudioV2WorkspaceRecordingPlans.ts`
- `frontend/src/services/studioV2WorkspaceDatabase.ts`、`frontend/src/services/studioV2WorkspaceRecordingPlans.ts`、`frontend/src/services/datalink.ts`
- `frontend/src/types/recordingPlan.ts`、`frontend/src/i18n/locales/`
- 相鄰 Go 測試、`frontend/tests/unit/`、`frontend/tests/integration/`、`frontend/tests/e2e/`。

**相容性變更**：尚未執行的記錄方案操作不能繼續回覆成功；舊客戶端必須處理未開放及不確定結果。前端與後端結果格式要成套部署。儲存資料如需新增版本／操作記錄，先新增相容欄位與讀取方式，再移轉；不能直接刪除既有設定。

**範圍外**：整個倉庫全面重構、新資料庫驅動、PLC 協議、採集排程、整套歷史報表、`/studio/runtime` 重設計、重新建立 `/studio`、現場部署。不能用局部測試宣稱全案或現場驗收完成。

## Draft Status

產品實作進行中：前置 C1 已完成並封存，階段 B（1.2／1.3／2.1／2.2／2.4）與 C3（implement-studio-v2-verified-schema-setup 3.1–3.10）皆已實作並驗證；剩餘 2.3 的 metadata 消費驗收、3.3／3.4 的真實試寫與第 4 節收斂。初稿環境限制、文件 review 與產品測試分開記錄於 validation.md；只有文件通過檢查不代表產品完成。
