## Why

已保存的寬表／長表記錄其實可以停用、移除與部分編輯，但入口隱藏在進階群組，固定布局又沒有明確的後繼設定流程。移除留下的 canonical key 不能直接重建；需要讓使用者安全管理設定，而不是誤以為要刪 SQL 歷史才能重做。

## What Changes

- 由伺服器回傳可編輯、需建立後繼版本與不可執行的操作，以及具體阻擋原因。
- 將「停用／刪除」呈現為「停止新記錄／移除記錄設定」，保留 accepted backlog、原目的地、歷史、收據與操作紀錄。
- 新增同一邏輯記錄身分下、不同群組世代的後繼建立與 reviewed transition；支援矩陣限定 SQLite／PostgreSQL 的 managed wide profile 與 plain managed long 雙向轉換，以及同引擎同布局同 storage mode 的目的地換接。
- 移除後重新建立必須明確建立新世代，不復活 tombstone、不刪 canonical key、不重播既有資料；來源／型態 ownership 交接保持唯一。
- 沿用既有 revision、readiness、schema confirmation、operation ledger，依賴 `freeze-counter-replay-semantics` 與 `fix-recording-cutover-boundaries` 的正確性驗收後才開放新 live transition。

## Capabilities

### New Capabilities

- `recording-lifecycle-management`: 記錄設定的操作分類、停止、移除、後繼世代與能力受限的安全轉換。

### Modified Capabilities

無；以新增生命週期契約補充既有 write-group authority，不改寫既有歷史交付語意。

## Impact

- 後端：`internal/datalink/workspace/` 的 write-group service、profile key、cutover 與 immutable versions；`internal/datalink/grouppipeline/backlog_guard.go`、`internal/datalink/schema/migrations/`、`internal/api/handlers/` 與 generated Swagger。
- 前端：`frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/`、`frontend/src/services/studioV2WorkspaceWriteGroups.ts`、`frontend/src/types/studioV2WriteGroup.ts` 與 en／zh-TW 文案。
- 相依分工：本 change 唯一擁有新增轉換矩陣與 successor lifecycle；既有 cutover boundary 修正歸 `fix-recording-cutover-boundaries`，引導式版面歸 `guide-step4-recording-workflow`。

## 本次驗收例外

2026-10-28 使用者同意同機 SQL Server Developer 測試並指示「若無法就跳過」。同機完整隔離的官方配置方式未能確立，故實機 SQL Server write/readback／cleanup 保留 NOT RUN；保留本地 adapter／guards 與測試，不放寬產品矩陣或安全檢查，後續操作復原與 UI 可繼續。發布報告須明列未驗證風險。

2026-10-28 15:27 UTC 使用者另明示「跳過並commit到github main後繼續」，批准本次實際瀏覽器操作與截圖驗收的獨立例外。已備妥並 discovery 42 個正式案例，但實際執行及截圖仍為 SKIP／NOT RUN；315 個前端測試檔、2,942 項測試、完整 lint、型別檢查與打包通過，不能代替真實 browser 驗收。保留測試與環境阻擋證據，允許完成本次交付並承接後續 change；不擴大任何產品能力或安全權限。一般開始記錄在 A→B→A 快速切換後可能持續顯示 unknown 的既有競態已在舊 HEAD 與目前來源重現，仍未修復，須在交接揭露；正常回歸通過不表示該競態已消失。
