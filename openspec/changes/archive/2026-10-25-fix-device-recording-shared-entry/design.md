## Context

原案 unify-device-recording-workflow 已以 68e18d0e 提交並歸檔。後續使用者以兩台設備的截圖指出舊寬表會自動把第一台帶回另一個主面板；截圖是行為問題證據，不是像素設計來源。本案只收斂已授權的入口修復，不重開原案或改寫其完成紀錄。

修正前新增兩個回歸測試，確實因第一台沒有共同設定面板而失敗；修正後通過。相容性測試改成先明確開啟舊編輯器，保留其值、草稿與不寫入斷言；既有 ownership 與只準備保護以原入口驗證，不捏造額外的測試先行失敗。

## Goals / Non-Goals

**Goals:**
- 舊寬表 A 與新設備 B 預設顯示相同的記錄用途與瞬時／累積欄位。
- 新設定草稿、原有編輯草稿與設備／工作區作用域保持獨立。
- 舊設定的編輯、只準備與已審核交接保持明確，切換檢視不改寫資源。

**Non-Goals:**
- 不自動遷移、刪除、停用或重啟舊群組，不變更歷史、來源歸屬、schema 確認或開始 API。
- 不新增 Go／資料庫功能、writer、全域狀態、依賴、頁面或報表。
- 不操作使用者截圖的真設備、LAN 或資料庫，不增加現場／Windows／SQL 矩陣驗收。

## Decisions

### 共同入口與獨立編輯

BasicRecordingPanel 預設呈現 CircuitProfileBasicPanel。舊群組只供 ExistingDeviceRecording 摘要顯示；必須按「編輯既有記錄」才進入舊編輯器。編輯／只準備檢視以 workspace＋device 為範圍，選設備或 Advanced 審核請求會回到共同設定。新表單沿用既有 session 草稿，舊編輯器沿用原草稿機制，不建立第二份 canonical 設定。

拒絕只依「有未刪除舊群組」把整台設備留在舊主流程，也不把舊設定自動改成新的平均／長表。

### 沿用伺服器權威

applied_revision 只能證明有套用過的版本，摘要不冒稱目前正在採集。新設定使用原 scoped recording-start；相同來源若仍有 owner，保留後端的拒絕與明確 reviewed cutover。開始請求只含審核的新群組，不帶入未選取的舊群組；未套用的舊草稿不構成前端的整台設備封鎖。既有只準備模式仍不恢復或發送開始操作。工作進行中不允許以切換編輯模式繞過原本的操作鎖。

### 有界驗收與交付

使用前端 service 模擬測試檢查共同入口、草稿隔離、舊值不變、scoped request 與 typed 拒絕；使用獨立 Chromium 在 loopback 元件 fixture 驗證 1024×720／1600×1000。這兩個尺寸是本次驗收條件，不是所有尺寸保證。瀏覽器 fixture 的服務為模擬，寫入一律拒絕，不宣稱 SQL readback 或現場驗收。

## Implementation Contract

| 行為與介面 | 失敗／保護 | 驗證 |
| --- | --- | --- |
| 兩設備都預設呈現用途、瞬時、累積與同一個 scoped start | 舊群組不觸發自動 create/update/apply/start | device-recording-tiles.test.tsx 的共同入口與 A/B 草稿案例 |
| 舊值在明確開啟後可編輯、返回與重載可恢復 | 60 秒舊值不被新 15 秒預設覆蓋；舊編輯草稿 73 與新草稿 47 分離 | tiles／compatibility／migration-prepare 測試 |
| 新開始只帶已審核的新群組，保留 server ownership | typed cutover 拒絕不改、停或刪舊群組；prepareOnly 不開始 | recording-workflow-unified／migration-prepare／migration-timescale 測試 |
| UI 與操作指南同步，沒有偽造 SQL 證據 | service 模擬與實際資料庫驗收分開 | 真瀏覽器 fixture、前端 tests／tsc／lint／build、diff 與行數檢查 |

## Risks / Trade-offs

- [曾套用但已交接的舊群組仍保留] → 使用條件式說明，來源是否仍被占用由原伺服器判斷，不用存在與否推斷 active owner。
- [一般開始復原與只準備模式不同] → 一般設定保留原操作查詢；只準備模式不再查詢或發送該操作，保留明確分界。
- [工作樹有另一案正在寫入] → 只以明列的 15 個修正檔歸屬；提交前比較內容與索引，排除所有外案內容。
