## Why

Step 4 同時呈現舊快照、新瞬時／累積、進階群組套用及工作區啟用，操作人員難以判斷有效設定與啟動範圍；完整電錶快照與異質點位記錄也尚未形成清楚的兩條路徑。使用者已確認瞬時資料按有效樣本算術平均保存，需辨識缺資料，並以設備方塊快速切換。

## What Changes

- 在既有 Step 4 提供「整台設備快照（寬表）」與「指定量測項目（長表）」兩條用途路徑；共用 canonical WriteGroup、來源身分、準備與開始流程。寬／長布局不等同瞬時／累積，兩種布局均可保存瞬時量、累積原值及其可信增量。
- 新建瞬時記錄預設每 15 秒做有效樣本算術平均；沿用可設 1–300 秒的範圍。數字均為需求設計條件。每筆有效讀值等權，拒絕補零、跨窗沿用及把累積計數器平均化；既有未選平均的群組繼續原本最新值語意。
- 累積量保留既有原值、增量、基準、精確數字及 gap/reset/epoch 政策；可設 1–600 分鐘、預設 10 分鐘（既有契約），不新增第二個計數演算法或採集器。布林／文字狀態採明示的最新有效值，不能冒充平均。
- 每窗保存有效／排除樣本數、實際樣本起訖與完整度證據。部分缺值、已知失敗或截斷需標示；沒有足夠證據判定完整時標示未知，整窗無有效值保存 NULL／原因，不杜撰斷線或完整率。
- Step 4 設備篩選改用具選取狀態、可用鍵盤操作的小方塊，禁止設備下拉；切換不改已存值、不丟草稿、不啟停設備，異步回覆只更新原範圍。
- 統一目的地 → 設備／用途／來源與週期 → 預覽及明確結構確認 →「套用並開始記錄」。進階編輯同一份設定；工作區只做總覽，不另放要求重做的啟用按鈕。既有 API 與 Share-only 的安全界線保留。
- 寬表才顯示整列完整性規則，說清楚必填缺值仍會略過整列；長表逐來源顯示結果。來源明細預設收合為項目數，診斷及錯誤按需展開，保留已落地的顯示名稱與 ⓘ 提示。
- **BREAKING（僅新明選平均的版本）**：瞬時整數的平均輸出可為小數，須經凍結的輸出型別與目的地相容性確認；不得把既有整數欄、歷史列或未明選的 latest 群組改成平均。
- 實作授權上限為本提案列出的用途、算法、設備篩選與啟用整合。任何新增功能、全面整理或額外驗收條件，必須指出原需求失敗的具體證據並先取得使用者授權；review 建議與代理分工不得自行擴大範圍。

## Capabilities

### New Capabilities

- `window-average-recording`: 明選且版本化的有效樣本算術平均、精確輸出、缺資料證據及兩種布局的耐久閉窗。

### Modified Capabilities

- `guided-recording-workflow`: 以設備方塊與用途引導兩種布局、同一份進階設定、精簡明細及禁止外擴的操作範圍。
- `write-group-row-semantics`: 區分既有 latest 與明選 mean，保留半開窗、去重、精度與舊 required 政策；新受管寬表允許按明示政策保存缺值證據及累積結果。
- `datalink-workbench-v2-step4-database`: 一個具明確設備／群組範圍的開始流程取代平行的 UI 啟用入口，保留 scoped recovery、Share-only 與既有全域 API 的相容性。

## Impact

- 後端設定與流程：internal/datalink/workspace/、internal/api/router_studio_v2_recording_start.go、internal/api/handlers/、必要的 additive 設定 migration；只延伸既有版本、確認、operation 與 ownership，沒有第二套設定權威。
- 記錄執行：internal/datalink/snapshot/、internal/datalink/runtime/、internal/datalink/dbtarget/、internal/datalink/groupdelivery/、internal/datalink/grouppipeline/、internal/datalink/recordingplan/；沿用既有 journal/outbox/receipt 與 counter 引擎。
- 畫面：frontend/src/features/datalink/workbench-v2/steps/step4/、相關 state/hooks/types/services、en／zh-TW locales 與對應測試；不修改 Step 1–3 主流程或 Runtime 監看頁。
- 驗證僅涵蓋受影響功能、可丟棄 SQLite／PostgreSQL 與既有長表 Timescale 模式，文件和產生的 API 說明依實際契約變更同步；不宣稱效能、正式設備或現場驗收。
- 唯一負責 change：unify-device-recording-workflow；依賴已實作且本機提交的 add-circuit-recording-profiles。前案仍 active、15/15，不能把它的完成證據當成本案證據，也不得在此改寫或封存前案。前案 newest-eligible 及 legacy-wide 主面板規則，只由本案的明選新版本與設備作用域規則取代；舊資料保持原解讀。
- 不包含 CMS 報表／查詢引擎、費率／告警／圖表、每樣本永久留存、新協議／採集排程、寬表 Timescale 化、保留／壓縮政策、歷史回填或自動雙向轉表。MySQL／MariaDB、SQL Server 維持原 parked changes；既有 ⓘ 的 staged 修改及其他 WIP 原樣保留。完整禁止事項與停止條件寫入 design.md。
