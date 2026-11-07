## Context

`/studio/v2` 的資料記錄首頁（frontend/src/features/datalink/workbench-v2/recording-home/）以表格列出已保存記錄（write group），並在頁底以「設備摘要」列出設備。2026-11-07 在本機以 1440×900 與 390×844 截圖觀察到：

- 設備只顯示 `availability_status`（「設定可用」），沒有即時採集狀態；每筆記錄的來源欄都顯示「目前來源健康尚未確認」。
- `GET /api/v1/datalink/runtime/status`（不帶 `device_id`）一次回傳所有 collectors，每筆含 `status`、`running`、`points_total/healthy/stale/error`、`last_read_at`、`last_error`、`breaker_state`、`effective_interval_ms`，外層有 `snapshot_state` 與 `metrics`（`write_success_total`、`write_error_total`）。同日實測：MC 為 `warning`、46/48 健康、2 點過期；另兩台 `idle`、`running=false`。
- 記錄列的目的地以 `<details>`「技術識別」包住資料庫、`schema.table` 與 connector id；寫入欄由 `RecordingHomeResultFacts` 以 `<details>`「詳細證據與觀測時間」包住設定、啟用、最後寫入、owner、回讀、清理與觀測時間。全部展開時列高約為原本 3 倍。
- 列上只有「開啟記錄」。停止／移除只存在於記錄編輯器的 `SavedRecordingActions` 與 `SavedRecordingRemoveDialog`。
- `DELETE /studio-v2/workspace/write-groups/:id` 目前把記錄標成 tombstone（`WriteGroupStatusDeleted`）；對已是 tombstone 的記錄再呼叫會回 `ErrWriteGroupLifecycleBlocked`。已套用的記錄會先經 `grouppipeline.BacklogGuard.CheckWriteGroupBacklog` 檢查未送達資料的歸屬，tombstone 的待送資料之後仍由 worker 送出。
- `DELETE /studio-v2/workspace/devices/:id` 會先 `DetachDevice`、再 `device.Service.Delete`（並連帶刪除該設備擁有的 source rules），**不檢查**是否仍有記錄的 `write_group_members.device_id` 引用該設備。

使用者決定：記錄的「移除」改為直接永久刪除，不保留 tombstone；設備被引用時直接阻擋；展開內容改為直接顯示在列上。

工作目錄中已 stage 但尚未 commit 的 recording-home 變更（`RecordingHomeStatusDistribution.tsx`、`RecordingHomeUnconfirmed.tsx`、`RecordingHomeBadges.tsx` 等）是本 change 的基礎。

## Goals / Non-Goals

**Goals:**

- 操作員一眼就能看出每台設備是否正在採集；狀態取不到時明確顯示「狀態無法取得」，不能顯示成正常。
- 首頁可直接停止與刪除記錄、刪除設備，且每個破壞性動作都先說明影響並需要明確確認。
- 刪除記錄後，內部設定資料庫不留下該記錄的資料列；外部目的資料庫的資料表與資料列不受影響。
- 刪除設備不會讓任何未刪除的記錄失去來源。
- 原本要點開才看得到的目的地與寫入階段，直接顯示在列上，列高維持一致。
- 頁首只有一個主要動作（新增記錄），版面在 1440、1024、768、390 寬都不會水平捲動。

**Non-Goals:**

- 不刪除或修改外部目的資料庫的資料表、資料列或 partition；這些要由資料庫管理者另外處理。
- 不在記錄仍有待送資料時強制丟棄資料；不提供「強制刪除」。
- 不刪除設備時連帶刪除引用它的記錄。
- 不改變 `/studio/runtime` 監看頁、`/test` 或 `/gateway/*`。
- 不新增獨立的設備健康 API；沿用 `GET /runtime/status`。
- 不提供刪除記錄後的復原或垃圾桶。
- 不加入寬螢幕右側詳細面板；本次只把資訊直接放在列上。

## Decisions

### 以 runtime status collectors 推導設備連線狀態

首頁每 5 秒（頁面可見時）呼叫一次不帶 `device_id` 的 `GET /runtime/status`，以 `device_id` 對應 collectors。狀態依下列順序判斷，第一個成立者勝出：

1. **狀態無法取得**（`unavailable`）：請求失敗、`snapshot_state.unavailable` 為 true，或該設備沒有對應 collector。
2. **中斷**（`disconnected`）：`breaker_state` 不是 `closed`，或 `status` 為 `error`。
3. **未採集**（`idle`）：`running` 為 false、`availability_status` 為 `unavailable`，或 `points_total` 為 0。
4. **部分異常**（`degraded`）：`points_error` 或 `points_stale` 大於 0、`points_healthy` 不等於 `points_total`（健康數與總數對不上，代表有點位狀態未被分類）、`status` 為 `warning`、`last_error` 非空，或 `last_read_at` 為空或早於 `max(3 × effective_interval_ms, 5000ms)`。
5. **採集中**（`collecting`）：其餘情況，也就是正在執行、點位全部健康、讀值夠新。

`snapshot_state.stale` 為 true 時仍顯示各設備狀態，但在頁首摘要加上「狀態可能過期」提示。

其他做法：改用既有的 `GET /devices/health`（只有斷路器資訊）或為每台設備各發一次請求。不採用，因為前者資訊不足，後者違反「Home data demand stays bounded」。5 秒與 5000ms 下限是設計條件，不是量測結果。

### 頁首狀態摘要與設備狀態列取代頁底設備摘要

頁首摘要一行列出三組數字：設備（總數，以及採集中／部分異常／中斷／未採集的數量）、記錄（總數與各保存狀態的數量，沿用既有 `RecordingHomeStatusDistribution` 的全量計算）、寫入（`metrics.write_success_total` 與 `write_error_total`，標籤寫明「runtime 累計」，不能寫成 SQL 已確認）。

摘要下方是設備狀態列，每台設備一張小卡：名稱、協議、狀態燈號（顏色加文字加圖示）、`健康數/總數 點`、相對讀值時間（例如「3 秒前」）。點卡片會套用既有的「設備」篩選；再點一次取消。手機版可水平捲動，但頁面本身不水平捲動。原本頁底的 `<aside className="gw-summary">` 設備摘要移除，「設備」分頁保留完整的連線設定、讀取點位、採集異常操作，並加上同樣的狀態燈號欄。

其他做法：保留頁底摘要並另外加燈號。不採用，因為內容與「設備」分頁重複，而且在手機版排版最亂。

### 目的地與寫入階段直接顯示在列上

- 目的地欄：第二行以等寬小字顯示 `schema.table`（過長時從中間截斷，`title` 為完整值），旁邊有複製按鈕，複製成功時以 `role="status"` 告知。資料庫名稱維持第三行。connector id 只放在複製按鈕旁圖示的 tooltip，不直接顯示。移除 `gw-destination-identity` 的 `<details>`。
- 寫入欄：以 `recordingHomeResultSummary` 的既有結果，畫出四格進度條「設定 → 啟用 → SQL 寫入 → 回讀」。每格有三種狀態：完成、未確認、需要注意，每格都附文字標籤（螢幕閱讀器可讀），不只靠顏色。

| 格 | 完成 | 未確認 | 需要注意 |
| --- | --- | --- | --- |
| 設定 | `configuration=requested` | `not_requested`（顯示「草稿未要求啟用」） | — |
| 啟用 | `effective=requested_effective` | `unknown` | `requested_not_effective` |
| SQL 寫入 | `sql` 為 `effective_revision_receipt` 或 `revision_receipt` | `unknown`、`awaiting_first_write`；`previous_revision_receipt` 與 `historical_commit` 顯示為「歷史寫入」並標成未確認 | `current=delivery_attention` |
| 回讀 | `readback=test_write_verified`（標籤寫「試寫回讀」） | `unknown` | — |

進度條下方一行顯示最重要的一句（沿用目前的 `result.sql.*` 文字）與 `lastCommittedAt` 的相對時間。滑過或聚焦某一格時，以 tooltip 顯示該格的觀測時間、版本與 owner、清理狀態等原本在 `<details>` 裡的內容。移除「詳細證據與觀測時間」的 `<details>`。

其他做法：右側詳細面板。不採用，使用者選擇直接顯示在列上。

### 記錄列操作選單提供停止與刪除

每列保留「開啟記錄」按鈕，旁邊加「⋯」選單（具 `aria-haspopup="menu"`，鍵盤可操作，Esc 關閉並把焦點還給按鈕）。選單項目：

- **停止**：記錄正在執行時才出現，沿用既有的 disable 契約與確認內容。
- **刪除記錄**：開啟改版後的 `SavedRecordingRemoveDialog`。

動作是否可用仍依既有「Recording actions follow matching server availability」由伺服器的動作事實決定；事實與顯示版本不符時，選單項目停用並提示重新整理。

### 刪除記錄改為永久刪除內部設定

**Supersedes**: 無（修改 `recording-lifecycle-management` 中「Stop and remove preserve accepted data」的移除語意）

`DELETE /studio-v2/workspace/write-groups/:id` 改為在同一個設定資料庫交易中永久刪除整個邏輯記錄：

- 刪除範圍：該記錄所在 lineage 的所有 generation（目前 head 與已結束的前代），以及以這些 group id 為鍵的內部資料列：`write_groups`、`write_group_versions`、`write_group_members`、`write_group_revision_ends`、`write_group_basic_keys`、`write_group_profile_keys`、`write_group_storage_scopes`、`write_group_storage_events`、`write_group_cutovers`、`write_group_ownership_fences`、`write_group_migration_maps`、`write_group_mysql_partition_maintenance`、`write_group_sqlserver_partition_maintenance`、`write_group_generations`、`write_group_lineages`、`wg_runtime_versions`、`wg_delivery_samples`、`wg_delivery_checkpoints`、`wg_delivery_buckets`、`wg_delivery_outbox`、`wg_delivery_receipts`，以及 recording plan／measurement 中只屬於該記錄的資料列。實作時須以 migration 清單逐表確認，測試斷言刪除後這些表中不再有該 group id 的資料列。
- 拒絕條件（回 409，什麼都不改）：
  - 任一 generation 的 `wg_delivery_outbox` 有 `state` 不是 `sql_committed` 或 `operator_skipped` 的資料列（待送、重試中、結果未知、被擋下）；或 `wg_delivery_samples` 有 `consumed = 0` 的資料列。錯誤碼 `WRITE_GROUP_DELETE_BACKLOG`，`action` 為 `stop_and_wait`，回應附上依類別分開的數量（樣本數、各 outbox 狀態的數量），不加總成單一數字。
  - 記錄仍在接收資料（intake active）。錯誤碼 `WRITE_GROUP_DELETE_ACTIVE`，`action` 為 `stop_first`。
  - 有保留中的後繼 generation、進行中的 cutover 或尚未結束的後繼操作。錯誤碼 `WRITE_GROUP_DELETE_OPERATION_PENDING`，`action` 為 `open_recording`。
  - 指定的 group 是某個未刪除的後續 generation 的前代（例如重建後留下的 tombstone）。錯誤碼 `WRITE_GROUP_DELETE_HAS_SUCCESSOR`，`action` 為 `open_recording`；操作員要從目前的 head 刪除整筆記錄。
  - 既有的 workspace／group／connector revision 不符時，維持既有 409 revision conflict。
- 已是 tombstone 的記錄，同一個 API 也可永久刪除（不再回 `ErrWriteGroupLifecycleBlocked`）。
- 不連線外部目的資料庫，不執行任何 DDL 或 DML。
- 成功後以既有的 `audit.Service.Record` 寫入一筆稽核紀錄：新增 `EventType` `write_group_deleted`，`ReferenceID` 為 group id，`Details` 只含記錄名稱、lineage id 與被刪的 generation 數，不含 connector 憑證或 DSN。稽核寫入失敗不回滾刪除，但要記錄 log（只記一次）。
- 刪除後同一設備／角色可以用一般流程重新建立記錄（canonical key 已釋放）。

其他做法：

- 保留 tombstone、首頁預設隱藏。使用者否決。
- 刪除時一併丟棄待送資料。不採用，會無聲地遺失已收進來的資料。
- 只刪 head、保留前代 generation。不採用，會留下操作員看不到也刪不掉的殘餘資料。

### 刪除對話框說明影響並要求明確確認

`SavedRecordingRemoveDialog` 改為「刪除記錄」對話框：

- 列出會刪除的內容（記錄設定、版本與內部交付紀錄）與會保留的內容（外部資料庫 `schema.table` 及其中已寫入的資料）。
- 註明「此動作無法復原」。
- 必須勾選「我了解此記錄的設定會被永久刪除」才能按下「刪除記錄」。
- 伺服器回 409 時，依錯誤碼顯示可行動的說明（例如「仍有 3 筆樣本、2 筆待送資料，請先停止記錄並等待送完」），並提供對應的按鈕（停止、開啟記錄、重新整理）。
- 不顯示原始後端例外訊息。

### 刪除設備前檢查記錄引用

`StudioV2WorkspaceDevicesHandler.Delete` 在 `DetachDevice` 之前，於同一交易中查詢 `write_group_members` 與 `write_group_basic_keys` 是否有 status 不是 tombstone 的記錄引用該 `device_id`。有引用時回 409，錯誤碼 `DEVICE_REFERENCED_BY_RECORDINGS`，`action` 為 `remove_recordings_first`，資料附上受影響記錄的 `id` 與名稱（最多 20 筆，另附總數）。只被 tombstone 引用時，先依記錄刪除規則永久刪除這些 tombstone，再刪除設備；任一 tombstone 無法刪除時（例如仍有待送資料），整個設備刪除失敗且不做任何變更；設備刪除成功時寫入 `device_deleted` 稽核紀錄（`ReferenceID` 為 device id，`Details` 只含設備名稱與協議）。

前端的「刪除設備」放在設備狀態卡與「設備」分頁的「⋯」選單中。對話框先以首頁已載入的記錄資料列出引用該設備的記錄與點位數；有引用時直接說明需先刪除這些記錄，並停用確認按鈕；沒有引用時，要輸入設備名稱才能確認。伺服器的 409 仍是最終依據，前端顯示伺服器回傳的記錄清單。

其他做法：只在前端檢查。不採用，因為 Step 1 編輯器或其他呼叫端仍可能直接呼叫 API，會留下懸空引用。

### 頁首版面、最後更新時間與可收合的未確認區

- 頁首一行：左邊是標題「資料記錄」，右邊是次要工具列（開啟設定、日誌、重新整理，使用次要按鈕樣式），最右邊是唯一的主要按鈕「新增記錄」。
- 「繼續編輯草稿」改為橫幅，寫出草稿名稱與「繼續編輯」按鈕；沒有草稿時不顯示。
- 「來源已於本次瀏覽重新觀察」與「僅查詢本頁詳細狀態…」兩行灰字，改為重新整理按鈕旁的「最後更新 HH:mm:ss」，旁邊的說明圖示以 tooltip 顯示原本的完整說明。
- 「尚未確認或未完成的記錄」改為可收合的提醒列，標題寫「N 筆需要處理」；有需要注意或結果未知的項目時預設展開，只有 tombstone 時預設收合。
- 互動回饋（取自 design-taste-frontend 的通用原則）：所有按鈕在 hover 時有可見變化、按下時有按壓回饋，有清楚的 focus ring；只用一組主要／次要按鈕樣式、一種卡片樣式、一組圓角；不加入捲動觸發的動畫；`prefers-reduced-motion` 時不播放動畫。顏色一律使用既有的 semantic design tokens（「Semantic design tokens govern affected surfaces」），淺色與深色主題都要符合文字對比。
- 手機版（≤767px）：每筆記錄是一張卡片，「開啟記錄」為全寬按鈕，其他操作在「⋯」選單裡；設備狀態列可水平捲動。

### 實作期間確認的刪除邊界

**Supersedes**: 補充並收窄「刪除記錄改為永久刪除內部設定」與「刪除設備前檢查記錄引用」中「刪除範圍」與「同一交易」的敘述；原決策其餘內容不變。

實作時逐一比對 `internal/datalink/schema/migrations/` 後，發現部分資料表不是以單一 group 為鍵，或同時是其他記錄的所有權證據。直接刪除會破壞其他記錄，因此 2026-11-08 補充下列規則：

- **舊記錄沒有 lineage**：未建立 lineage 的單代舊記錄刪除成功時回 `lineage_id: ""`、`deleted_generations: 1`；稽核紀錄也記空字串。前端只在 `deleted_generations` 為 1 時接受空 `lineage_id`，其他不符合的回應一律視為「結果未知」。
- **共用 key 的資料表**：`write_group_storage_scopes`（032/034）、`write_group_mysql_partition_maintenance`（041）、`write_group_sqlserver_partition_maintenance`（042）以 workspace、connector、database、schema、mode、layout 為共用 key。只有可證明不再被任何存活記錄的目前、歷史或 runtime 版本引用時才刪除；版本無法解碼時視為無法證明，保留。
- **沒有記錄 owner 的資料表**：`measurement_definitions`（018/039）、`recording_plans`（019）、`gw_delivery_journal`（020）沒有可證明只屬於單一記錄的 owner，一律保留。`gw_delivery_outbox`／`gw_delivery_receipts`（020）以 `destination_id` 等於 group id 清除，未送達時拒絕。
- **跨記錄的 cutover、fence、handoff**（036/037/044/045/046）：存活記錄驗證所有權時會讀取 cutover 兩側的 `write_group_versions`；且 `write_group_versions.group_id` 有 `ON DELETE RESTRICT`。因此不能只保留 cutover 而刪除本記錄的 group 與 version。尚未完成時回 409 `WRITE_GROUP_DELETE_OPERATION_PENDING`。已完成但仍被其他存活記錄當作證據時，回 409 `WRITE_GROUP_DELETE_HAS_SUCCESSOR`（`action` 為 `open_recording`），`details` 為 `{ "referenced_by": [{ "id", "name" }], "total": <n> }`（最多列 20 筆）。兩種情況都不做任何變更。
- **schema／試寫操作紀錄**（`managed_schema_operations`，023/026）：
  - 與該記錄相關的操作若為 pending、running、partial、unknown，或為 `failed` 且 `executed_statements > 0`，或試寫清理尚未完成，回 409 `WRITE_GROUP_DELETE_OPERATION_PENDING`。
  - `failed` 且 `executed_statements = 0`、沒有待清理工作者視為已結束，與記錄一起刪除。
  - 狀態組合互相矛盾或無法判斷時一律拒絕。
- **`recording_start` 操作紀錄**：屬於 workspace 層紀錄，同時記錄設備啟用，也可能包含多筆記錄。若 `detail.request.groups` 涉及本 lineage 且狀態為 pending、running、partial、unknown 或未知，在刪除交易內回 409 `WRITE_GROUP_DELETE_OPERATION_PENDING`；範圍無法解析時同樣拒絕。已結束（succeeded／failed）的紀錄一律保留。讀取端須容忍其中的 group 已不存在：
  - 首頁只依目前存在的記錄組裝內容。
  - 成功操作重播只回歷史結果，不會重新啟用。
  - 失敗操作重播在 group 不存在時拒絕。
- **與 recording start 併行**：在 SQLite 上，若 start 在刪除檢查之後先提交，刪除交易的寫入會因過期快照得到 `database is locked` 而完整回滾，重試時回 409。若 start 在刪除提交之後才 claim，它在套用與啟用前會重新驗證 group，找不到時拒絕。
- **刪除設備與 Modbus Share**：引用檢查、tombstone 清除、detach、設備與 source rule 刪除在同一個設定資料庫交易內完成。之後才依既有流程呼叫 Modbus Share 的 `InvalidateWorkspaceMembership`，因為它會另開 durable revision 讀寫，不能放進單一連線的 SQLite 交易。
  - 若這一步失敗，設備已刪除且無法回滾。API 回 500 `DEVICE_DELETED_SHARE_SCOPE_REFRESH_FAILED`，Share 維持 fail-closed，錯誤只記一次 log。
  - 提交後關閉設備連線失敗時，回 500 `DEVICE_DELETED_CONNECTION_CLOSE_FAILED`。
  - 交易提交到 Share 更新之間若程序中斷，目前沒有驗證。
- **交付回條（`wg_delivery_receipts`，025）**：這張表只有 `effect_key`、`payload_digest`、`committed_at`，只能透過 outbox 找到所屬記錄。`groupdelivery` 的 `Reclaim` 改為在同一交易內，先刪除符合條件的回條，再刪除過期的已送達 outbox，之後不會再產生無法歸屬的回條。修正部署前已被回收 outbox 留下的舊回條無法歸屬到任何記錄，永久刪除不會清除它們；如需清理，要另外決定保留期限或資料修復方式。
- **多台設備共用的 tombstone**：無法隨單一設備安全刪除，回 409 `DEVICE_REFERENCED_BY_RECORDINGS`，不做任何變更。
- **設備刪除遇到 tombstone 待送資料**：若 tombstone 仍有待送資料，設備刪除回記錄刪除的 409 `WRITE_GROUP_DELETE_BACKLOG`。前端顯示「設備尚未刪除」與分類數量，不當成結果未知。

其他做法：保留已刪記錄的 group、version 與 revision end 作為其他記錄的證據。不採用，因為等於保留 tombstone，與使用者決定的永久刪除相違。

## Implementation Contract

**行為：**

- 在 `/studio/v2` 首頁，每台設備都顯示以下其中一種狀態：採集中、部分異常、中斷、未採集、狀態無法取得。狀態以文字、圖示與顏色同時呈現，並附健康點數與相對讀值時間；狀態 API 失敗時，最多保留上次成功的狀態 15 秒並標示可能過期，之後一律顯示「狀態無法取得」。
- 頁首摘要的三組數字，與首頁已載入的設備、記錄以及 runtime metrics 一致；寫入數字標示為「runtime 累計」。
- 記錄列不再有 `<details>`。目的地第二行顯示 `schema.table` 並可複製；寫入欄顯示四格進度條。
- 列上「⋯」選單提供停止與刪除記錄；設備卡與設備分頁的「⋯」選單提供刪除設備。

**介面與資料形狀：**

- `DELETE /api/v1/datalink/studio-v2/workspace/write-groups/{id}`：request body 維持既有的 `WriteGroupMutationRequest`（含 expected revisions）。成功回 200，`data` 為 `{ "deleted": true, "group_id": "<id>", "lineage_id": "<id>", "deleted_generations": <n> }`；沒有 lineage 的舊記錄回 `lineage_id: ""` 且 `deleted_generations: 1`。失敗時回 `APIErrorResponse`：
  - 409 `WRITE_GROUP_DELETE_BACKLOG`（`action` 為 `stop_and_wait`；`details` 為 `{ "samples_unconsumed": <n>, "outbox": { "<state>": <n> } }`，`outbox` 只列數量大於 0 的狀態，未知狀態也原樣列出）
  - 409 `WRITE_GROUP_DELETE_ACTIVE`
  - 409 `WRITE_GROUP_DELETE_OPERATION_PENDING`
  - 409 `WRITE_GROUP_DELETE_HAS_SUCCESSOR`（跨記錄證據時 `details` 為 `{ "referenced_by": [{ "id", "name" }], "total": <n> }`）
  - 既有的 409 revision conflict、404、503
- `DELETE /api/v1/datalink/studio-v2/workspace/devices/{id}`：新增 409 `DEVICE_REFERENCED_BY_RECORDINGS`，`details` 為 `{ "recordings": [{ "id", "name" }], "total": <n> }`，`action` 為 `remove_recordings_first`；提交後的錯誤為 500 `DEVICE_DELETED_SHARE_SCOPE_REFRESH_FAILED` 與 500 `DEVICE_DELETED_CONNECTION_CLOSE_FAILED`。
- 兩個 API 的 Swagger 註解同步更新，並以 `make gen-docs` 重產 docs/swagger/。
- 前端：在 frontend/src/types/ 新增連線狀態型別 `DeviceConnectionState = 'collecting' | 'degraded' | 'disconnected' | 'idle' | 'unavailable'`；新增純函式 `deriveDeviceConnectionState(collector, nowMs)`；`studioV2WorkspaceWriteGroups.remove` 解析新的成功資料形狀與錯誤碼。

**失敗模式：**

- 狀態 API 失敗：保留上次成功的狀態最多 15 秒，並標示「狀態可能過期」；超過 15 秒顯示「狀態無法取得」。不顯示原始錯誤訊息。
- 刪除遇到 409：對話框保持開啟，依錯誤碼顯示可行動的說明與按鈕，不關閉、也不顯示成功。
- 刪除遇到網路錯誤或結果未知：顯示「無法確認是否已刪除」，提供重新整理按鈕；重新整理後依伺服器實際結果顯示。
- 複製失敗（瀏覽器不允許）：以 `role="status"` 顯示「無法複製」，並選取該文字讓使用者手動複製。

**驗收：**

- Go 測試：
  - 永久刪除後，上述各表不再有該 group id 的資料列。
  - 有未消化樣本或各種待送 outbox 狀態時回 409，且資料完全不變。
  - intake active 時回 409。
  - 有進行中的後繼操作或 cutover 時回 409。
  - 指定前代 tombstone 而後代仍存在時回 409 `WRITE_GROUP_DELETE_HAS_SUCCESSOR`。
  - 可永久刪除既有 tombstone。
  - 刪除後可用同一設備／角色重新建立記錄。
  - 設備被未刪除記錄引用時回 409 並附記錄清單；只被 tombstone 引用時，刪除成功且 tombstone 一併清除。
  - Swagger contract 測試鎖定新的回應狀態。
- Vitest：
  - `deriveDeviceConnectionState` 涵蓋五種狀態與每條判斷邊界（例如 `last_read_at` 剛好等於門檻值）。
  - 頁首摘要的數字與資料一致。
  - 四格進度條對照上表的每種輸入。
  - ⋯ 選單的鍵盤操作與焦點還原。
  - 刪除對話框在 200 與每種 409 時的顯示。
  - 刪除設備對話框在有引用、無引用、名稱不符時的狀態。
- Playwright：在 1440、1024、768、390 寬截圖首頁，確認沒有水平捲動，記錄列內沒有 `<details>`；附上截圖。
- `make check-lines` 通過。

**範圍：**

- 範圍內：recording-home 元件、樣式與 i18n；`SavedRecordingRemoveDialog`；write group 刪除服務與 handler；workspace device 刪除 handler；Swagger。
- 範圍外：外部資料庫的清理；`/studio/runtime`；記錄編輯器其他流程（編輯器內的移除按鈕改為呼叫同一個刪除對話框，但不另外改版）；新的健康 API。

## Risks / Trade-offs

- [永久刪除無法復原，操作員可能誤刪] → 對話框列出刪除與保留的內容、需勾選確認、伺服器以 revision 防止過期操作，並寫入 `write_group_deleted` 稽核紀錄。
- [刪除範圍漏掉某張以 group id 為鍵的內部表，留下殘餘資料] → 實作時以 migrations 目錄逐表比對，測試對每張表斷言 0 筆；新增表時須更新同一份清單與測試。
- [規格原本要求保留 tombstone 與操作歷史，改為永久刪除會失去稽核脈絡] → 以 `write_group_deleted` 稽核紀錄保留「何時刪除哪筆記錄」；不保留其版本與交付明細。此取捨已由使用者決定。
- [重建已移除記錄（recreate tombstone）的流程失去用途] → 該流程與其規格不變，仍適用於尚未清除的既有 tombstone；新的刪除不再產生 tombstone，所以之後不會再有新的可重建對象。
- [每 5 秒輪詢 runtime status 增加負載] → 只在頁面可見時輪詢、只發一個請求；離開頁面或分頁隱藏時停止。這個負載影響尚未量測，可用 network panel 觀察每分鐘的請求數與回應大小來確認。
- [狀態門檻（3 倍週期、最少 5 秒）可能讓慢速設備誤報為部分異常] → 門檻基於 `effective_interval_ms`；現場若誤報，再調整倍數。
- [本 change 疊在尚未 commit 的 recording-home 變更上] → apply 前先確認那批已 commit，再開始實作。

## Migration Plan

1. 先 commit 工作目錄中既有的 recording-home 變更。
2. 後端：依 expand → migrate → contract 推進。
   - 先加入永久刪除服務與新錯誤碼，以及設備引用檢查。
   - 再讓 handler 改用永久刪除。
   - 最後移除「刪除後產生 tombstone」的分支。
3. 前端：先加入連線狀態與摘要；再把展開內容改為直接顯示；最後加入 ⋯ 選單與刪除對話框，並移除頁底設備摘要與 `<details>`。
4. 不需要資料庫 migration，因為只刪除資料列、不變更 schema。
5. 回復方式：回退程式版本即可恢復 tombstone 行為；但已永久刪除的記錄無法還原，需從設定資料庫備份還原。

## Open Questions

（無；使用者已決定刪除語意、設備引用處理與資訊呈現方式。）
