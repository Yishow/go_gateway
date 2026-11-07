## Why

`/studio/v2` 資料記錄首頁目前看不出設備是否正在連線採集（只顯示「設定可用」，每筆記錄都寫「目前來源健康尚未確認」），而 `GET /api/v1/datalink/runtime/status` 早已回傳每台設備的採集狀態、點位健康數、最後讀值時間與斷路器狀態。首頁也沒有停止、刪除記錄或刪除設備的入口；已移除的記錄以 tombstone 長期佔據「尚未確認」卡片；資料表名稱與寫入證據藏在兩層點一下展開裡，展開後列高變成三倍，操作員難以快速掃視。

## What Changes

- 首頁顯示每台設備的即時連線狀態（採集中、部分異常、中斷、未採集、狀態無法取得），附健康點位比與最後讀值時間；status 讀取失敗時顯示「狀態無法取得」，不預設為正常。
- 頁首加入狀態摘要列（設備數與各連線狀態數、記錄數與各保存狀態數、runtime 累計寫入成功／失敗並標明層級），並以頁首設備狀態列取代頁底重複的「設備摘要」；點設備可篩選記錄。
- 記錄列的「技術識別」與「詳細證據與觀測時間」改為直接顯示：目的地欄以等寬字顯示 `schema.table` 並可複製；寫入欄以「設定 → 啟用 → SQL 寫入 → 回讀」四格進度條分別呈現各階段（完成／未確認／失敗），滑過顯示觀測時間。
- 記錄列新增「⋯」操作選單：停止（沿用既有 stop 契約）與刪除。
- **BREAKING**：刪除記錄由「留下 tombstone」改為永久刪除該記錄在內部設定資料庫的全部設定與內部交付紀錄；外部目的資料庫的資料表與已寫入資料列一律保留。仍有未送達、重試中或結果未知的交付資料、或有進行中的轉換／後繼操作時，後端拒絕刪除（409）並指出需先停止並等待送完。既有 tombstone 也可用同一操作清除。
- 刪除設備：對話框先列出受影響的記錄與點位數，需輸入設備名稱確認；後端在設備仍被未刪除記錄引用時回 409 並列出記錄，不留下懸空引用。
- 排版：「新增記錄」為唯一主要按鈕；「開啟設定／日誌／重新整理」降為次要工具列；兩行灰色說明改為「最後更新 時:分:秒」加說明提示；「尚未確認或未完成的記錄」可收合；手機版每筆記錄為卡片、操作收進「⋯」。

## Non-Goals

（詳見 design.md 的 Goals / Non-Goals。）

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `simple-recording-home`: 新增設備即時連線狀態、頁首狀態摘要與設備狀態列、目的地與寫入階段直接顯示、記錄列操作選單與刪除確認、頁首版面與可收合未確認區的需求。
- `recording-lifecycle-management`: 「移除」改為永久刪除記錄設定（保留外部資料表與資料列、有待交付資料時拒絕）；既有 tombstone 可被清除；待交付資料的確認情境改為「刪除會被拒絕並分類列出」。
- `simple-recording-lifecycle`: 簡易介面的移除動作改為刪除，確認內容改為說明永久刪除設定並保留外部歷史資料。
- `datalink-workbench-v2-step1-device`: 刪除已保存設備前須檢查是否仍被未刪除記錄引用，被引用時拒絕並回報受影響記錄。

## Impact

- Affected specs: `simple-recording-home`, `recording-lifecycle-management`, `simple-recording-lifecycle`, `datalink-workbench-v2-step1-device`
- 後端：internal/datalink/workspace/write_group_service.go（刪除改為永久刪除與待交付檢查）、internal/api/handlers/studio_v2_workspace_write_groups_handler.go（錯誤碼與 Swagger 註解）、internal/api/handlers/studio_v2_workspace_devices_handler.go 與 internal/datalink/workspace/service_devices.go（設備引用檢查）、docs/swagger/（以 make gen-docs 重產）。
- 前端：frontend/src/features/datalink/workbench-v2/recording-home/（RecordingHome.tsx、RecordingHomeCollection.tsx、RecordingHomeUnconfirmed.tsx、RecordingHomeResultFacts.tsx、recording-home.css 與新元件）、frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/SavedRecordingRemoveDialog.tsx、frontend/src/services/studioV2WorkspaceWriteGroups.ts、frontend/src/services/studioV2WorkspaceDevices.ts、frontend/src/i18n/locales/en 與 zh-TW 的 recording-home.json。
- 資料：內部設定資料庫中屬於被刪記錄的 write_group 相關資料列會被永久刪除，無法復原；外部目的資料庫不受影響。
- 依賴：工作目錄中尚未 commit 的 recording-home 變更（RecordingHomeStatusDistribution.tsx、RecordingHomeUnconfirmed.tsx 等）須先 commit，本 change 在其上實作。
