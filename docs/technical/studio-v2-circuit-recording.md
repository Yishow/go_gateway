# Studio V2 迴路記錄：設備快照、量測長表與已審核的切換

本文件說明 Step 4 的共同操作與資料契約，包含 `unify-device-recording-workflow` 與 `add-mysql-mariadb-circuit-recording` 的長表擴充。SQL Server 的長表實作、條件閘門與驗收證據見 [SQL Server 長表說明](studio-v2-sqlserver-circuit-recording.md)。群組的 API、revision 與回復背景見 [studio-v2-write-groups.md](studio-v2-write-groups.md)；規格以 `openspec/specs/` 中的 window-average-recording、write-group-row-semantics、guided-recording-workflow 及 datalink-workbench-v2-step4-database 為準。

## 範圍與限制

- 共同用途支援 SQLite、PostgreSQL plain，以及長表專用的 PostgreSQL TimescaleDB；新寬表不支援 TimescaleDB。MySQL／MariaDB 僅延伸受管長表，不新增其受管寬表；SQL Server 長表已完成 17.0.5035.5／compat 170（ARM 模擬環境）實測驗收，並依 16.x／160 與 17.x／170 條件判定能力閘門，未知或不符版本與未驗證 scope 維持關閉。
- MySQL／MariaDB 已完成各別的本機 SQL、相容性與真瀏覽器驗收，受管長表僅對本批實測的完整版本配對開放。連線成功、driver 存在或其中一家通過，都不能替代另一家及各項操作的驗證。
- 下方本機 SQL 與真瀏覽器證據，不代表 Windows、真 PLC、LAN 或正式資料庫驗收，也不提供效能保證。
- service 模擬元件測試、正常 server 入口的真瀏覽器、SQL readback 與現場驗收是不同證據，不互相代替。

## 先選用途，再設定型態與週期

| | 整台設備／迴路快照（`managed_wide_profile_v1`） | 指定量測項目（`managed_long_v1`） |
| --- | --- | --- |
| 資料表名 | `gw_group_` + 群組 ID 的雜湊 | `gw_long_` + 16 位十六進位，由 workspace、連線、資料庫、schema 與儲存模式的雜湊決定 |
| 欄位 | 每個點位一個欄位，名稱為 `v_` + 裝置、點位、Tag 身分的雜湊；點位增加就增加欄位 | 固定 29 欄（28 個欄位 + 1 個擁有權標記欄 `_gw_owner_…`），與點位數量無關 |
| 一列代表 | 一個群組在一個週期的一列 | 一個來源在一個記錄窗口的一列 |
| 增加點位 | 需要結構變更 | 只增加列，不增加欄位、不為迴路另建表 |
| 共用 | 一個群組一張表 | 同一 scope 的所有迴路與型態共用一張表 |

兩種用途共用相同的瞬時／累積語意、來源角色及版本規則。既有 `managed_wide_v1` 與 custom 保留原布局、週期和缺值規則，不會因檢視或切換用途而自動遷移。

長表的資料表名來自 scope，不是群組名稱，也不是點位數量。主鍵是 `(bucket_start, record_id)`，並有 `idx_<資料表名>_src_time`（`workspace_id, tag_id, bucket_start`）支援依來源與時間查詢。每列帶有 workspace、群組、版本、迴路、裝置、點位、Tag 身分、型態、窗口起迄、實際取樣時間、品質、狀態與 provenance；有效數值只填入對應的型別欄位（整數、decimal、float、bool、text），其他型別欄位為 NULL。

## 型態、單位、預設值與範圍

| 型態 | 單位 | 預設 | 範圍 |
| --- | --- | --- | --- |
| 瞬時 | 秒 | 15 | 1–300 的整數 |
| 累積 | 分鐘 | 10 | 1–600 的整數 |

API 必須帶型態對應的單位；缺少或不符的單位、小數、非有限值與超出範圍的值都會被拒絕（`profile_interval_invalid`）。已儲存的舊週期不會被改動。記錄週期獨立於設備採集與傳輸寫入頻率：採集比記錄週期慢時，沒有樣本的窗口會顯示為缺值，不會重複寫入舊讀值。

## 算術平均與既有 latest

新建瞬時記錄預設每 15 秒取算術平均，包括切換流程中準備的新目標。只有數值型 gauge 參與平均：同一窗口的 10、20、30 得到 20；同一 SampleID 重送不重複計數，壞值、過期、窗外與已閉窗後的晚到資料不進分母。不插值、不補零、不沿用上一窗口的值，也不改採集頻率。窗口採 UTC 半開區間，邊界樣本歸下一窗；已審核切換使用實際截斷範圍。

整數來源的平均結果是獨立的 decimal，不改 Tag 型別。`9007199254740992` 與 `9007199254740993` 的平均是精確文字 `9007199254740992.5`。循環小數最多保留 18 位，以 half-even 捨入；provenance 保留 exact sum、count、scale 和 rounding。float 平均明示 `basis=float64`，不宣稱任意精確。不相容的整數結果欄位會被阻擋，不會截去小數。

既有未保存 `sample_reduction` 的群組仍讀為 `latest`。重新載入、切換設備或 Basic／Advanced 不會改成平均；改 mean 必須明確建立相容的新 revision，重新確認並套用。

## 缺值與覆蓋率

- 沒有可用觀測時保存 NULL／missing，不表示零，也不能直接推論斷線。
- 有可用值但伴隨已知失敗或截斷時，保留結果並顯示 partial 與實際原因。
- 有樣本與筆數、沒有整窗完整性證據時，覆蓋率是 unknown；不能由「沒有錯誤」推論 complete 或 100%。
- 新寬表預設來源為 optional、採 partial：部分或全部缺值均可保存支援的 NULL 欄及逐來源證據。明設 required 的來源缺值時，仍依整列規則跳列並留下原因；長表則每個來源各自記錄。
- 狀態畫面的唯讀來源證據，從同一已套用版本的凍結 layout、payload 與匹配 receipt 投影方法、筆數、型別與覆蓋率。證據缺失或不符時顯示不可用，不回傳歷史數值。

## 來源選擇

- 瞬時與累積各自獨立選來源，同一迴路的兩個型態可以使用不同來源。
- 系統不會用名稱或單位猜測來源身分，也不會把 raw Tag 自動當成計量值。未指派迴路的 raw Tag 仍可選，不會被賦予迴路或量測身分。
- 累積型態需要已確認的計量語意（counter）。Signed 與 rollover 政策在進階的來源設定確認，未確認時不會產生負差或補償圈數。
- numeric gauge 取平均；bool／text 的 state 取最新可用值；已確認的 counter_raw 也保留最新原值，不參與平均。同一計數器可同時加入瞬時與累積，各用自己的週期與基準。
- 來源預設顯示數量摘要，展開後在清單內捲動；收合不會隱藏失聯身分、來源修復或 counter 確認阻擋。

## 原始值、差值、首筆未知、重置與缺窗

累積記錄保留最新可用的原始讀值。長表使用 `delta_value`、`delta_from`／`delta_to`；新寬表的值欄保存 raw，每個來源的 `provenance.counter` 保存精確文字 `delta`、實際 `delta_from`／`delta_to`、status 與政策。兩者重用既有 counter 計算與同交易 checkpoint，不另寫差值公式。

- **首筆**：第一個可用窗口只存原始值，狀態為 baseline；長表差值為 NULL，寬表 counter 不帶可計算的 delta。最後可靠讀值成為下一窗口基準，未知差值不能讀成零。
- **重置或不確定**：非 signed 計數器的無法解釋下降，差值為 NULL（`reset_or_uncertain`），原始值與事件仍保留。例：窗口內讀到 100、10、110，保留 110，差值為 NULL。
- **Signed**：只有明確確認為 signed 的計數器才接受負差。
- **Rollover**：需要設定模數、有效範圍與足以排除多次回繞的速率／時間上限，才會補償；證據不足時差值為 NULL。
- **缺窗**：缺少的窗口維持 NULL，不推進可靠基準。跨缺窗的可信差值以實際觀測區間標示為 gap 的總和，不會平均分配到沒有支撐的子窗口。
- **換 epoch、來源綁定或 profile 版本**：重新建立 baseline。
- **精確度**：整數與 decimal 不經過 JavaScript Number。例：uint64 `9007199254740992` 到 `9007199254740993` 的差值是 1；decimal `12000.0 → 12000.4 → 12000.7` 的差值是 NULL、0.4、0.3。

## 進階與工作區的角色

- **基本設定**：用設備 ID 綁定的方塊選設備，再選用途、型態、來源與週期、準備資料表並開始記錄。同名或同位址仍是不同設備；Tab、Space、Enter 可操作。切換只換檢視範圍，不建立、套用或啟停群組。草稿與未明結果依 workspace／device 隔離，可在重新載入後恢復；連線設定不是另一個週期入口。
- **既有寬表設備**：同樣先開啟上述記錄設定，並顯示舊設定摘要；「編輯既有記錄」才開啟原有週期／缺值編輯器，返回或切換設備不會修改舊群組。未套用的舊草稿不阻擋新設定；新設定接手已套用群組的相同來源時，仍須走明確的已審核切換，不能直接覆蓋舊 writer。
- **進階**：編輯同一份已儲存 profile 的 canonical 來源對應、custom 版面，以及已確認的 signed／rollover 政策。它不是另一種資料表格式，也沒有獨立的記錄週期。
- **工作區**：唯讀檢視已儲存群組的整體 readiness 與套用狀態，沒有另一個全域啟用入口。適用受管群組的 Advanced 生效動作回到同一個 Basic 設備／群組審核入口；custom 與既有工程生命週期 API 保留。
- 儲存設定不等於套用，也不等於開始記錄；已套用的版本持續寫入，直到套用新版本。取消勾選基本設定中的 profile 只改變草稿，已套用的寫入器要在進階的群組生命週期停用（依已儲存的版本檢查）。

## 準備、確認、開始與狀態

1. **準備**：預覽是唯讀的，顯示將建立的資料表與索引。
2. **確認**：每個群組各自做明確的資料表結構確認，綁定當下的 workspace、群組與連線版本。共用長表中新增來源的群組，即使預覽沒有任何 DDL，也要有自己的新確認；其他群組的確認與套用快照不受影響。
3. **開始**：確認後，只套用並啟動目前審核的設備與群組。A 成功不要求未準備的 B 一起完成，也不會啟停未選取的 B；採集已在運作仍不等於資料庫記錄已開始。
4. **回應遺失**：保留原 request ID 與群組意圖。切換設備或重新載入不會將它套到其他設備；查詢原操作或明確重送同一要求，不另建啟動意圖或版本。

狀態分開呈現：已儲存、已套用、採集中、本機已持久化、SQL 已提交。SQL receipt 不會自動變成「已做外部 SQL 讀回驗證」，也不保證 counter 差值已可計算。SQL 重送沿用凍結 payload digest 與 effect，不重新計算平均或推進 counter 基準。

## 儲存模式與缺表重建

- **plain**（預設）：PostgreSQL 與 SQLite 都可用，不需要任何擴充。
- **timescale**：只有 PostgreSQL 且必須明確選擇。後端先唯讀探測擴充、`by_range`／`create_hypertable` 功能與權限；準備時建立以非空 `bucket_start` 分區、含 `bucket_start` 與 `record_id` 複合唯一鍵的空 hypertable，預覽中顯示 7 天 chunk 間隔。它不會啟用或安裝擴充、不轉換已有資料或外來的表、不加入保留或壓縮工作，也不會悄悄退回 plain。能力不足時回傳 `timescale_unavailable` 並附修復動作；SQLite 會拒絕 timescale。
- **recreate_missing**：確認資料表結構時顯示。啟用時，只在原本已確認的受管資料表被刪除時，用先前確認的版面重建，並以冪等方式重送尚未提交的效果。只重建結構與尚未提交的資料；已儲存的歷史不保證能恢復。一般 SQL 失敗、連線版本改變、沒有確認、外來或不相容的同名表、缺少擁有權證據或 Timescale 不可用時都不會重建（`schema_repair_blocked`、`storage_owner_mismatch`）。

## MySQL／MariaDB 長表與明確試寫

沿用 `mysql` 連線種類；Step 4 的家族、版本與各項可用操作，來自後端依已儲存連線及 revision 唯讀探測的結果。修改目的地後須等待新 scope 的結果；舊回應不能替新連線放行。未驗證、metadata 禁止或能力不可用時，依畫面說明修復，不會建立資料庫或降級 TLS。

- **plain**：預設為普通 InnoDB 長表，與 SQLite／PostgreSQL 的共同 profile、窗口及 receipt 語意一致。
- **native_monthly**：只有實際 capability 可用時才可明選。UTC `RANGE COLUMNS(bucket_start)` 初始保留當月及後續三個月，另有過去與 `MAXVALUE` 未來 catch-all；receipt 仍是普通表。已儲存的模式不因重新載入或能力暫時失效而退回 plain。
- **有限維護**：只擴充有限未來窗；長期停機不逐月無界補齊，也不自動 DROP、TRUNCATE、MERGE 或實施歷史保留政策。維護失敗顯示 `partition_maintenance_blocked`；只有已驗證 catch-all 仍可安全接收時才繼續交付。
- **部分 DDL**：逐步保存並查證實際 table／index／receipt；部分成功不表示整批 rollback。保留原 operation_id 查詢或恢復，外來同名物件不會被刪改。
- **試寫**：先預覽，再明確確認 token 與 operation_id。長表按來源產生有界的測試列，預覽顯示筆數及第一筆範例；後端保存本次精確 record IDs、來源與 receipt digest，讀回後只清除此操作的資料及收據。表的 `gateway` 擁有權標記不會改成測試標記，也不能單獨當刪除條件。
- **結果**：`written_verified` 只代表試寫及實際讀回相符，還要分別查看 cleanup。部分來源已提交時保留 partial／unknown 與清理結果；確認尚未寫入的失敗為 failed／not_attempted。完成操作重送回傳原結果，不重寫測試列。
- **相容轉接**：legacy `recording-plans/test-write` 沿用目前共同安全流程；裸 plan_id、不完整確認、過期或其他 scope 的 token 都不會授權寫入。

MySQL 的數值槽使用精確 `DECIMAL(65,30)`，時間為 UTC `DATETIME(6)`；來源識別大小寫必須保持。已確認整數 counter 若保存 `exact_recording_scale` 十進位文字，例如 `120000 → 120004` 配 `0.1`，只在 journal 前縮放一次，得到 raw `12000.0 → 12000.4` 及差值 `0.4`。SQL 可補足小數位，讀回比較保持精確數值，不經 float64；超出可表示範圍於接受前拒絕。更改 scale 開新 epoch，舊 outbox 不重算。

MySQL／MariaDB 的唯讀查詢使用反引號及各自 SQL client 的參數綁定，例如：

```sql
SELECT bucket_start, record_id, status,
       CAST(value_decimal AS CHAR) AS exact_raw,
       CAST(delta_value AS CHAR) AS exact_delta
FROM `<database>`.`<long_table>`
WHERE workspace_id = ? AND tag_id = ?
ORDER BY bucket_start;
```

## 與既有寫入器共存

既有的 wide 與 custom 群組保留原本的身分、版面、週期、完整性政策與執行方式，直到明確的已審核切換；Basic 與 Advanced 的切換、重新載入或檢視都不會新增或改寫任何設定。新長表或新寬表型態若要接手仍由執行中群組負責的同型態來源，一般「套用」會被拒絕（`cutover_required`，HTTP 409），必須走下面的切換；換 ID 或重送 Ensure 不能繞過既有 owner。

## 準備目標與已審核的切換

整個流程中，目前的寫入器持續運作，直到邊界時間。

1. **準備目標**：共同記錄設定可準備新用途；若要使用明確的只準備模式，先開啟「編輯既有記錄」，再選「為已審核的切換準備記錄目標」。這是只準備的情境：自己選寬表／長表用途、型態、來源與週期，不會複製舊群組的任何設定；不會套用、啟動、停止或移除任何群組。新瞬時目標預設 mean，舊 writer 不變。確認資料表結構是明確的資料庫步驟，與套用群組不同。關閉後目標保留為草稿。
2. **選擇與審核**：在進階的群組區塊找到「把來源切換到已準備的目標」，選擇目前的群組、已準備的受管 profile 目標與未來邊界時間，按「審核差異」。審核是唯讀的，顯示名稱、布局、儲存、週期、目的地資料表、結構確認狀態、來源清單與移動／保留的來源；新寬表目標仍須有完整的型別、來源與目的地 proof。
3. **確認**：確認時送出的是審核回傳的邊界文字（保留微秒）與審核摘要，以及審核當下的 workspace、來源、目標與連線版本。選取項目、邊界或任一已儲存群組在審核後改變，審核即過期，必須重新審核；確認鈕不會用舊的審核摘要搭配新的版本。
4. **部分來源**：未被目標涵蓋的來源不會被悄悄移走。後端要求目前群組已儲存、且恰好保留這些來源的草稿；請先在群組編輯器中自行編輯，面板只提供開啟該群組的入口，並不修改成員或政策。
5. **邊界語意**：邊界時間是不可變的。目前群組的窗口在邊界截斷並標為 partial；新窗口從邊界開始，到下一個對齊窗界，同樣標為 partial。同一來源的同一時間只屬於一個生效版本。目前群組已接受的資料依原目的地送完；新的累積型態從邊界後第一窗重新建立 baseline。
6. **失敗**：審核或套用失敗時，目前群組維持已套用的設定，畫面只顯示可執行的說明，不顯示原始 SQL 或連線字串。邊界在未來時顯示為「已排定」，不是已在收集或已提交。
7. **回應遺失**：若套用的回應遺失，畫面顯示「未知」，不假設成功或失敗。「檢查已記錄的切換」會讀取該來源已記錄的切換，並以審核摘要、來源、種類與邊界比對；找不到相符紀錄時，才可明確選擇「再送出同一個已審核的請求」。不會自動重送，也不會改送不同內容。

## 回復

在同一個面板選「回復一次切換」，選擇一筆已記錄且尚未回復的切換，輸入新的未來邊界，審核後確認。回復使用新的邊界與新的審核摘要，不重用舊的；不刪除或改寫歷史、不補寫舊資料表、不刪除目標表。目標已接受的資料仍寫入原目的地，兩邊歷史都保留。被恢復的群組是否就緒，由後端判斷（restored_group_not_ready）。

## 查詢範例

以下使用佔位符，請代入自己的 schema、資料表名與識別碼。長表的時間為 UTC；SQLite 以固定微秒的文字儲存（例如 `2026-10-24T00:00:00.000000Z`），數值以精確的文字儲存。

```sql
SELECT bucket_start, bucket_end, status, value_float, value_decimal, value_int, delta_value
FROM "<schema>"."<long_table>"
WHERE workspace_id = '<workspace-id>'
  AND tag_id = '<tag-id>'
  AND bucket_start >= '<start-utc>' AND bucket_start < '<end-utc>'
ORDER BY bucket_start;
```

uint64 與 decimal 請以資料庫的精確型別或文字讀出（例如 `value_int::text`），不要在 JavaScript 以 Number 轉換；那會失去精確度。

新寬表的每個來源各占一個 provenance 元素。以下是 SQLite 的唯讀範例；baseline 的 delta 路徑不存在時回傳 NULL，不能補成零：

```sql
SELECT g.record_id, g.bucket_start,
       json_extract(m.value, '$.member') AS source_identity,
       json_extract(m.value, '$.mean.coverage') AS mean_coverage,
       json_extract(m.value, '$.counter.status') AS counter_status,
       json_extract(m.value, '$.counter.delta') AS exact_delta,
       json_extract(m.value, '$.counter.delta_from') AS delta_from,
       json_extract(m.value, '$.counter.delta_to') AS delta_to
FROM "<wide_table>" AS g, json_each(g.provenance) AS m
WHERE g.group_id = '<group-id>'
ORDER BY g.bucket_start;
```

## 驗證證據與限制

- SQL：`TestDeviceRecordingIntegration` 已實跑新 wide/plain SQLite、PostgreSQL，新 mean long/plain SQLite、PostgreSQL，以及既有 long Timescale。核對 20／1.5／9007199254740992.5、NULL／partial／unknown、counter baseline／delta／reset／gap、匹配 receipt、凍結 effect 與 cleanup；Timescale 另查實際 hypertable、時間維度及 chunks，沒有新增 wide Timescale。
- 耐久與歷史：`TestAverageCommitAckLost`、`TestProfileCommitAckLost`、`TestReviewedWideMeanTransition` 與 `TestWideToLongCutover` 提供實際重送及歷史目的地證據。這些是分開的代表案例，不宣稱每種故障與每個 SQL 模式都做了全組合測試。
- 真瀏覽器：`frontend/tests/e2e/device-recording-workflow.spec.ts` 使用正常 server、loopback 模擬設備及自有 SQLite。已實跑 1024×720／1600×1000 的原生鍵盤、同名設備草稿／迴路隔離、32 來源收合／內捲動、counter 確認、兩用途、明確 DDL 確認及單一 scoped 開始。後端真提交後只中斷 HTTP 回應，切 B／reload 後重送原 request，核對同一操作／版本；三群組的 SQL、outbox、receipt 對同一 applied revision，並產出兩尺寸截圖。
- 重跑 E2E 必須明確提供 `DEVICE_RECORDING_E2E_FIXTURE`：manifest 指定全新、自有、預先由正式 Step 1–3 API 設定的 workspace、兩設備、來源、SQLite 連線與資料庫路徑。初始記錄群組必須為空；未提供 fixture 會標示 NOT RUN，不會重設現有資料。不得指向正式 DB／設備。
- 未完成的環境驗收：Windows、真 PLC、LAN 與正式資料庫。元件測試不替代真瀏覽器，SQL receipt 不替代外部 readback，以上也不代表現場驗收。

### MySQL／MariaDB 本批證據

- 受管 SQL 操作只接受本批實測的完整版本字串：MySQL `8.4.13`、MariaDB `11.4.13-MariaDB-ubu2404`。其他 patch／發行版後綴仍可顯示唯讀版本資訊，但不得以主次版號相同授權建表、試寫、恢復或交付。
- `TestMariaDBCircuitIntegration` 獨立驗證普通表、精確 decimal、跨月分區、長期停機及故障注入；MySQL 的結果不替代 MariaDB。
- `TestMySQLGroupProbeLifecycle` 已在 MySQL 8.4.13、MariaDB 11.4.13 的專用可丟棄 fixture，各自實跑 plain／native_monthly。經共同 API 預覽、確認、實際寫入、讀回及清理，核對原有 SQL 資料與 receipt 未改動；完成後重送沿用原操作結果。
- `TestMySQLNativeExplicitReconfirmationAfterConnectorRevision` 已在兩家 fixture 驗證：連線 identity revision 改變後，以新能力 proof 明確重新確認，保留已授權的 7 個分區與 generation；舊自動恢復 proof 拒絕，已接受 outbox 不重算。單純開放已驗證操作不改變實體能力 hash，既有結構確認不因此失效。
- 分區試寫在有限維護擴充後使用目前授權的實體布局，原確認的分區邊界與已接受 payload 不變。權限拒絕的部分案例使用真實 catalog 加故障注入，未宣稱已實測受限角色。
- `TestMySQLDestinationCutover` 已在兩家 fixture 經真實 Update／Apply 與 SQL 讀回驗證：過期 CAS 拒絕，舊 outbox 的 payload bytes 與目的地欄位不變且排入原表，新版本只寫新表，累積 profile 保持原目的地。task 4.2 初驗使用僅移除開放閘門的 Go overlay；task 4.3 已以正常產品程式、無 overlay，在兩家 fixture 重跑通過。
- 本批真瀏覽器已分別使用正常 macOS 後端、獨立設定資料庫、loopback Modbus 模擬器及兩家真實 SQL fixture，完成 plain／native_monthly 的明確建表、雙 profile 開始記錄及四組試寫。每家均讀回精確 raw `12000.0 → 12000.4`、首筆 baseline／NULL delta 與後續 delta `0.4`，並對上同一 applied revision 的 outbox、receipt 及實際 SQL 列；原生表各有 6 個預覽分區。這些是本機模擬來源與真實 SQL 證據，不代表 PLC／Windows／LAN／正式資料庫驗收。
- 兩家各四次瀏覽器試寫均為 `written_verified`＋`cleaned`，與正常記錄列分開；schema 確認後若顯示結果待確認，使用「查詢操作狀態」查同一筆 operation，不重送建表。1024×720／1600×1000 截圖與 WebM 錄影已保存於本輪本機驗收附件。上方 unify 的 SQLite 瀏覽器證據不替代本批驗收。
- 前端完整回應現在保留 backend 確認的精確來源 proof，包含可省略的 counter policy；試寫預覽僅對 `values[].value` 沿用 1024 字元上限，以容納實際 provenance，其餘欄位與通用 JSON 預設界線不變。
- 受驗證的兩個完整 family／版本配對才開放 managed long、native partition、group test-write 與 missing repair；仍需實際連線、metadata、InnoDB 與目標 scope 查證。未知 patch／發行版與未確認 scope 保持拒絕；舊 custom／wide 能力契約不因新增 long 支援而放寬。
- 本批最後檢查：受影響 12 個 Go packages tests／vet／lint、來源與分區 migration focused tests、前端 lint／280 檔 2057 tests／build 通過。前端 lint 有兩個既有 E2E console warning，build 保留 chunk size 提醒；其他資料庫 fixture 與實體平台未執行項目分列，未以 skip 當通過。
