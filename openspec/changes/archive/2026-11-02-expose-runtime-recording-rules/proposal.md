## Why

執行監看目前把來源規則縮成地址／數量，把資料庫縮成連接器／表名，操作者無法回答「這個 Tag 現在依哪個版本、什麼方式、多久寫到哪裡」。已有 durable 群組交付證據未在此頁串起，必須讓設定、實際生效與 SQL 證據可以對照，而不是增加更多綠色狀態。

## What Changes

- 新增 workspace／selected-device 範圍的唯讀規則總覽與完整 drilldown：設備、來源地址、有效輪詢週期、解碼與轉換、Tag、寫入群組、目的表／欄位、寬長表、寫入週期與 latest／mean／delta。
- 明確區分已儲存草稿、實際生效版本及切換時間；設定缺失、runtime projection 不可證明及跨版本讀取皆安全顯示未知，不以 draft 冒充 active。
- 重用 canonical WriteGroup、不可變 applied snapshot、DeliveryView 與 SQL receipts；顯示 accepted、queued、retrying、blocked、unknown、SQL committed、舊 revision backlog、窗口與最近交付時間。
- 規則與錯誤可精準回到同 workspace 的 device／rule／mapping／group 設定，刪除或過期 identity 有明確 fallback。
- 將安全規則與交付證據置於 focused monitor 主視線；歷史報表保持次要，不擴張報表功能；修正無條件 LIVE 與跨設備設定摘要。
- 本 change 唯一擁有共用安全讀取模型、identity／revision／stage／action 與定位契約；make-runtime-logs-actionable 消費此契約。
- 本次只建立提案，不實作、部署、連線私人主機或宣稱現場驗收。

## Capabilities

### New Capabilities

- `runtime-recording-rule-view`: 可追溯設備與資料庫規則、active/draft 真相、交付證據、設定定位及安全降級的 focused runtime view。

### Modified Capabilities

無。既有 runtime-dashboard-truthfulness、studio-runtime-workspace-view、runtime-write-group-delivery 的要求保持，本案新增其上的規則讀取能力。

## Impact

- 後端讀取契約：internal/api/handlers/runtime_workspace_setup_context.go、internal/api/handlers/runtime_workspace_setup_types.go、internal/api/handlers/runtime_handler.go、internal/api/router.go；重用 internal/datalink/workspace/write_group_lifecycle.go、internal/datalink/grouppipeline/delivery_view.go 與 delivery_sources.go。
- 前端：frontend/src/features/datalink/runtime-dashboard/、frontend/src/types/studioV2RuntimeContext.ts、frontend/src/services/studioV2RuntimeContext.ts、frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts；重用 write-group delivery parser／hook 與 en／zh-TW 文案。
- API 以既有 workspace/runtime-context、runtime/status、write-groups delivery 讀取鏈擴充安全投影；不加入寫入、DDL、測試寫入或遠端程序日誌代理。
- 依賴 fix-step4-recording-authoring 的正確來源／記錄方式、manage-recording-lifecycle-safely 的狀態與 tombstone；依賴 guide-step4-recording-workflow 先提供 canonical scope-selection interface，再由本案擁有 query／resolver／deeplink 並加入 Step 4 binding，導航依賴不反向；安全讀取 provider 依分期先交付 guided。
- 不變更 add-windows-tray-runtime-logs 的 local-only guard／broker／Windows 工作與既有完成證據。實際來源與驗收界線見 design.md。

## 執行順序調整

**Supersedes:** 原先整個 change 必須等待前案全部完成的執行安排，僅改為可隔離準備先行、正式啟用與整合驗收維持原依賴順序。可先固定安全 typed contract 與未掛載 parser；實際安全 reader/projection 依下述補正，在 freeze→cutover→lifecycle 與 backend proof revisions 交接後先行；定位與最終驗收仍待 guided，不能以 fixture 填入實況。 七個 change 與本案需求／驗收不變；本次只更新 artifacts，不代表產品已實作。詳見 design.md「執行分期與交接」。

**Supersedes（安全讀取順序補正）:** 上述整案等待僅保留於 runtime dashboard／locator／最終驗收；runtime 1.3–1.4 的單一安全 projection、全舊新 redaction、applied-config-provider 與共用 service/hook 在 freeze→cutover→lifecycle 後先交付，guided 2.2 與 authoring 前置一併消費。區分 saved、scheduled/latest requested applied 與 as-of effective revision／effective_at，retained snapshot 不代表 intake permission；詳見 design「安全生效證據交接」。已完成歷史不變，以 pending compensation 修正 caller 契約。
