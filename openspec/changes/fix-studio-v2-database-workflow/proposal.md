# 修正 Studio v2 資料庫流程與操作介面

## Why

第 4 步必須讓使用者知道「資料存到哪裡、是否真的建立、是否真的寫入」，而不是只看到成功提示。本次依 main 的實際程式起草：記錄方案的建表／試寫入口直接回覆成功；目標連線寫死；多設備成員全部套第一台設備；示範欄位被說成真實欄位；輸入與預覽也有更新互相覆蓋的問題。來源、確定程度與尚未完成的檢查見 evidence.md。

基準：Yishow/go_gateway，main @ 1a0311c8e8db9c62fe4388f0701ba38afe552ff7。分支：openspec/fix-studio-v2-database-workflow。日期：2026-09-15（Asia/Taipei）。

本次只建立提案、設計、需求與待辦，不套用上一輪修補包，不修改產品程式、不更動正式資料、不啟用設備。這不是整個倉庫已完成逐檔審查，也不是修復已完成。

## What Changes

- 先停止假成功：未接上真實操作的記錄方案建表／試寫，明確回覆未開放；不關閉另一條已呼叫服務的一般建表功能。
- 使用已儲存的真實連線、測量項目及其設備歸屬；方案要明確選取，不以第一筆清單或猜測編號代替。
- 讀取真正資料表與欄位；建議範本、待確認配對及已確認配對分開，保留合法的同筆資料分組規則。
- 建表預覽綁定工作區、方案、連線與設定版本；執行前再次核對。試寫、讀回、清理分別記錄，未知結果不盲目重做。
- 設定儲存不能只成功一半卻顯示已完成；新增故障情境測試。
- 保留正在編輯的文字，Enter 只送一次，中文選字不提前送出，舊預覽不更新新畫面，唯讀時不能建表或試寫。
- 第 4 步改成「存到哪裡 → 使用現成表或建立新表 → 資料放哪個欄位 → 檢查並開始」，補上易懂文案、錯誤、鍵盤操作與窄畫面。
- 修正規格內仍要求固定動畫與無後端呼叫的舊成功定義，以實際後端結果為準；不重寫既有設備啟用與 Modbus Share 保護。

## Capabilities

### New Capabilities

無。沿用既有能力，不另建平行資料庫設定入口。

### Modified Capabilities

- `recording-database-setup`：補齊未開放狀態、實際目標與設備歸屬、版本核對、真實試寫及一致儲存的可驗收規則。
- `datalink-workbench-v2-step4-database`：補齊真實欄位、分組衝突、編輯／預覽保護、操作順序與以後端為準的結果顯示。

`workspace-database-row-groups` 的合法共用欄位規則保留，不改成所有同名欄位一律禁止。`openspec/specs/` 主規格本次不直接改寫；差異規格留在本 change 供後續驗證與同步。

## Impact

後續預計修改範圍，不代表本次已修改：

- `internal/api/handlers/studio_v2_workspace_recording_plans_handler.go`
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

草案已成形，產品實作尚未開始。執行環境沒有 OpenSpec／Spectra CLI，未執行正式 analyze、validate 或 park。須先完成 validation.md 的正式文件檢查，再開始 tasks.md；不以本次簡單格式檢查取代官方驗證。
