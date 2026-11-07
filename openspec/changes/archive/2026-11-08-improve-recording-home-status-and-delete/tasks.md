## 1. 前置

- [x] 1.1 確認工作目錄中既有的 recording-home 變更（`RecordingHomeStatusDistribution.tsx`、`RecordingHomeUnconfirmed.tsx`、`RecordingHomeBadges.tsx` 等）已 commit，並記錄基準：`cd frontend && npx vitest run tests/unit/workbench-v2/recording-home-shell.test.tsx tests/unit/workbench-v2/recording-home-collections.test.ts` 與 `go test ./internal/datalink/workspace/... ./internal/api/...` 的結果寫入交接紀錄；任一失敗時先回報、不開始實作。驗證：`git status --short frontend/src/features/datalink/workbench-v2/recording-home` 無未提交變更，基準輸出已記錄。

## 2. 後端：刪除記錄改為永久刪除內部設定

- [x] 2.1 先寫測試（TDD，新增 `internal/datalink/workspace/write_group_permanent_delete_test.go`），涵蓋「Stop and remove preserve accepted data」修改後的情境：兩代 lineage 已送完時永久刪除，且設計所列每張內部表都沒有該 group／lineage id 的資料列，SQLite 目的表列數不變；未消化樣本與各 outbox 狀態（`pending`、`retrying`、`unknown`、`blocked`、`quarantined`、`sending`）各自回 `WRITE_GROUP_DELETE_BACKLOG` 並分類計數；intake active 回 `WRITE_GROUP_DELETE_ACTIVE`；保留中後繼或 cutover 回 `WRITE_GROUP_DELETE_OPERATION_PENDING`；指定前代 tombstone 回 `WRITE_GROUP_DELETE_HAS_SUCCESSOR`；既有 tombstone 可清除；revision 過期回 conflict；刪除後同設備／角色可再建立記錄。驗證：`go test ./internal/datalink/workspace/ -run PermanentDelete` 在實作前失敗（預期中的失敗）。
- [x] 2.2 [after: 2.1] 在 `WriteGroupService.Delete` 實作永久刪除：同一交易內依 lineage 刪除所有 generation 的內部資料列，沿用並擴充 `grouppipeline.BacklogGuard` 的待送判斷，不連線目的資料庫；成功時以 `audit.Service.Record` 寫入 `write_group_deleted` 稽核紀錄（不含憑證或 DSN），並移除「刪除後標記為 tombstone」的分支。表清單須與 `internal/datalink/schema/migrations/` 逐表比對。驗證：`go test ./internal/datalink/workspace/ ./internal/datalink/grouppipeline/` 全部通過，含 2.1 的測試。
- [x] 2.3 [after: 2.2] `StudioV2WorkspaceWriteGroupsHandler.Delete` 回傳新的成功資料 `{deleted, group_id, lineage_id, deleted_generations}`，以及四個 409 錯誤碼（含 `action` 與分類計數 `details`），不洩漏原始錯誤；更新 Swagger 註解並以 `make gen-docs` 重產 docs/swagger/。驗證：新增 handler 測試逐一斷言每個錯誤碼的 status、code、action；在 `internal/api/router_studio_v2_recording_swagger_test.go`（或同目錄新的 contract 測試）鎖定 DELETE 的回應矩陣；`go test ./internal/api/...` 通過。

## 3. 後端：刪除設備前檢查記錄引用

- [x] 3.1 [after: 2.2] 依「Persisted device deletion cascade」修改後的需求，`StudioV2WorkspaceDevicesHandler.Delete` 在 `DetachDevice` 之前，於同一交易中檢查未刪除記錄的 `write_group_members` 與 `write_group_basic_keys` 是否引用該設備：被引用時回 409 `DEVICE_REFERENCED_BY_RECORDINGS`（附最多 20 筆記錄與總數），且不做任何變更；只被 tombstone 引用時，先依記錄刪除規則清除，任一失敗則整體失敗；成功時寫入 `device_deleted` 稽核紀錄；同步更新 Swagger。驗證：在 `internal/api/handlers/studio_v2_workspace_devices_handler_cascade_test.go` 新增被引用、只被 tombstone 引用、tombstone 有待送資料三種情境的測試；`go test ./internal/api/... ./internal/datalink/workspace/...` 通過。

## 4. 前端：以 runtime status collectors 推導設備連線狀態

- [x] 4.1 新增 `DeviceConnectionState` 型別與純函式 `deriveDeviceConnectionState(collector, nowMs)`，依「Home shows live device connection state」的順序判斷五種狀態，門檻為 `max(3 × effective_interval_ms, 5000)`。驗證：新增 `frontend/tests/unit/workbench-v2/device-connection-state.test.ts`，以規格的範例表（含剛好 5 秒、6 秒、沒有 collector、breaker open）逐列斷言，`npx vitest run` 該檔通過。
- [x] 4.2 [after: 4.1] 實作「頁首狀態摘要與設備狀態列取代頁底設備摘要」，滿足「Home header summarizes devices, recordings and runtime writes」：只在頁面可見時每 5 秒讀一次不帶 `device_id` 的 runtime status（失敗時保留上次結果最多 15 秒並標示可能過期，之後顯示「狀態無法取得」）；頁首摘要三組數字（寫入標示 runtime 累計）；設備狀態列可點選篩選、再點取消；移除頁底 `gw-summary` 設備摘要；「設備」分頁加入狀態欄；新增 en／zh-TW 文字。驗證：新增 `recording-home-device-status.test.tsx`，涵蓋摘要數字、篩選切換、失敗 15 秒前後、分頁隱藏時停止輪詢（fake timers 並斷言請求次數）。

## 5. 前端：目的地與寫入階段直接顯示在列上

- [x] 5.1 依「Recording rows show destination identity and write stages inline」，目的地欄第二行以等寬字顯示 `schema.table` 並提供複製（剪貼簿被拒時改為選取文字並以 `role="status"` 告知），connector id 只放在 tooltip；寫入欄改為四格進度條，依 design 的對照表從 `recordingHomeResultSummary` 推導，每格有文字標籤，hover 與鍵盤聚焦可看到原本 `<details>` 的內容；移除列內兩個 `<details>`。新樣式放在新的 CSS 檔，不讓 `recording-home.css` 超過 300 行。驗證：新增 `recording-home-inline-evidence.test.tsx`，以規格範例表斷言每格狀態與標籤、複製成功與失敗、列內 `querySelectorAll('details').length === 0`。

## 6. 前端：記錄列操作選單提供停止與刪除

- [x] 6.1 [after: 2.3, 3.1] `studioV2WorkspaceWriteGroups.remove` 解析新的成功資料與四個 409 錯誤碼；`studioV2WorkspaceDevices` 的刪除解析 `DEVICE_REFERENCED_BY_RECORDINGS` 與其記錄清單；未知回應視為「結果未知」而不是成功。驗證：services 單元測試以 mock 回應涵蓋 200、每個 409、網路錯誤，`npx vitest run` 通過。
- [x] 6.2 [after: 5.1] 每列新增「⋯」選單（`aria-haspopup="menu"`、方向鍵移動、Esc 關閉並把焦點還給按鈕），執行中的記錄才出現「停止」（沿用既有 disable 契約），伺服器動作事實允許時才出現「刪除記錄」；事實與版本不符時兩者停用並提示重新整理，符合「Recording actions follow matching server availability」。驗證：新增 `recording-home-row-actions.test.tsx`，涵蓋鍵盤操作、焦點還原、停用條件。
- [x] 6.3 [after: 6.1, 6.2] 依「刪除對話框說明影響並要求明確確認」改版 `SavedRecordingRemoveDialog` 為刪除記錄對話框：列出刪除與保留的內容、註明無法復原、需勾選確認；成功時清單與摘要不需重新載入就更新；每個 409 依錯誤碼顯示可行動說明與按鈕，待送資料依「Historical destinations and pending units remain explicit」分類列出、不加總；網路錯誤顯示「無法確認是否已刪除」與重新整理；記錄編輯器原有的移除入口改用同一對話框。驗證：新增 `recording-delete-dialog.test.tsx`，涵蓋未勾選時停用、200、四個 409、網路錯誤；既有引用 `SavedRecordingRemoveDialog` 的測試同步更新後通過。
- [x] 6.4 [after: 6.1, 4.2] 依「刪除設備前檢查記錄引用」與「Home offers recording and device actions with explicit confirmation」，在設備狀態卡與「設備」分頁加入「⋯ → 刪除設備」：對話框以已載入記錄列出引用該設備的記錄，有引用時停用確認並說明需先刪除；無引用時需輸入完全相同的設備名稱；伺服器 409 時改顯示伺服器回傳的記錄清單。驗證：新增 `device-delete-dialog.test.tsx`，涵蓋有引用、名稱不符（`mc-new` 對 `mc_new`）、成功、伺服器 409。

## 7. 前端：頁首版面、最後更新時間與可收合的未確認區

- [x] 7.1 [after: 4.2] 依「Home header has one primary action and visible freshness」調整頁首：「新增記錄」為唯一主要按鈕，設定／日誌／重新整理為次要按鈕；兩行灰色說明改為「最後更新 HH:mm:ss」加說明提示；有草稿時顯示寫出草稿名稱的繼續編輯橫幅；未確認區可收合，標題寫數量，只有 tombstone 時預設收合；手機版（≤767px）記錄為卡片、次要操作在「⋯」；按鈕有 hover、按下與 focus 回饋，`prefers-reduced-motion` 時不播放動畫，顏色只用既有 semantic tokens。驗證：更新 `recording-home-shell.test.tsx`，斷言主要按鈕數量為 1、未確認區預設展開與收合條件、草稿橫幅出現條件。

## 8. 驗證與文件

- [x] 8.1 [after: 2.3, 3.1, 6.3, 6.4, 7.1] 以 Playwright 對 `/studio/v2` 在 1440×900、1024×768、768×1024、390×844 截圖（淺色與深色），斷言 `document.documentElement.scrollWidth === innerWidth`，且記錄列內沒有 `details` 元素；新增 `frontend/tests/e2e/recording-home-status-layout.spec.ts`。另執行 `cd frontend && npx vitest run`、`go test ./...`、`make check-lines`，並確認 en／zh-TW 的 key 一致。驗證：所有指令通過，截圖附在交接紀錄；未通過或未執行的項目如實列出。
- [x] 8.2 [after: 8.1] 更新 `docs/technical/studio-surface-inventory/studio-v2-runtime.md` 與 `backend-api-registry.md`，記錄首頁新的狀態來源、刪除契約與錯誤碼，並以 `go run ./cmd/studio_inventory_changelog add` 寫入 changelog（含 summary、surface、files、reason）。驗證：`go run ./cmd/studio_inventory_changelog list -limit 1` 顯示這筆紀錄。

## 變更說明（非任務）

- 2026-11-08 實作期間補充 design 的「實作期間確認的刪除邊界」：共用 key 的儲存範圍與分割維護表、沒有記錄 owner 的資料表，以及已結束的 `recording_start` 紀錄都保留不刪；跨記錄的 cutover／handoff、schema／試寫操作與 `recording_start` 的拒絕條件；沒有 lineage 的舊記錄回空 `lineage_id`；設備刪除在提交後才更新 Modbus Share，失敗時回錯誤碼。這些規則已在 2.2、2.3、3.1 實作並由其測試涵蓋，6.1、6.3、6.4 的前端解析與對話框也已對應，不另外新增任務。
