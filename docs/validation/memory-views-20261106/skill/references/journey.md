# 新增設備到完整記錄

由空工作區在 `/studio/v2` 透過使用者路線新增設備、來源與記錄，終點是實際目的地的同版本 SQL 資料，再重載與進入 `/studio/runtime` 確認。

讀 repo 當下的 `simple-recording-home`、`simple-recording-editor`、`guided-recording-workflow`、`studio-v2-workspace-readiness`、`studio-v2-device-to-sql-acceptance` specs。規格、程式與 UI 不符時分別記錄，不拿舊路線當現況。

每個協議各做來源階段，每個協議 × 資料庫組合各做記錄階段。清單只代表宣告範圍。

## 0. 版本與環境

清單 ID：`build-provenance`。

- 固定 HEAD、dirty diff 與候選檔案雜湊；dirty build 不標成 clean HEAD。
- 前端 build → static sync → production `cmd/test_ui` binary；核對實際提供的 SPA 資源。embed placeholder 不算可用 UI。
- 留 backend、前端、simulator 各自埠、hc 名稱、PID、DB 及 ownership marker。
- 空設定 DB 與空目的地分開；缺 SQLite 檔先 stat，不能用一般 connection 開出檔案假驗空。已有檔使用唯讀 URI。
- 專用 Chromium 及相符版本；零測試、skip、fixture 缺失都不是通過。

## 1. UI 新增設備

ID：`source/<protocol>/device-create`。

1. 首頁切「設備」→「新增設備」，或目前 UI 明示的新建入口；不可 POST 預建待驗設備。
2. 填名稱、協議、loopback host／port、必要站號與 timeout，核對欄值沒有串接。先取消一次，確認未寫入。
3. 重填保存，核對一次 create、設備 ID、workspace membership 與 canonical config。
4. F5 重開，確認名稱、協議、host／port、站號與格式。用 ID 對身份，不能只用名稱。

同位址的兩台設備維持不同 ID；正式 multi-entity 要驗同一群組中的 A、B，不能用兩個單設備群組取代。

## 2. 就緒與點位

ID：`source/<protocol>/readiness`、`range-create-edit`、`range-activation`。

- 設定保存、connect、probe／read、planning／activation／apply eligibility 分開核對。
- 新增範圍後保存；再開既有範圍修改並保存。兩路都有 payload、ACK、canonical GET 與重載，不能只驗新增。
- 未設定資料庫時 database setup revision 可明確為空字串；缺失／純空白不是合法空版本，不能用 Modbus Share revision 代替。
- 停用規則可保存離線設定。啟用未就緒設備時應有可行動說明，不是 generic 400，也不繞過 readiness。
- 核對「測試連線並儲存」：按之前不偷偷 connect/probe，按後成功才保存啟用；失敗保留草稿與原停用保存值。
- workspace、rule revision 或 connection config 在測試／查詢間改變時拒絕過期保存；重複點擊不多建來源。
- 追 Point → Tag → mapping IDs；有範圍卡片不等於可記錄。

## 3. 記憶體與真實數值

ID：`source/<protocol>/memory-views-values`、`source-failure-recovery`。

- 來源準備、設備採集、記錄啟動是不同階段。未開始採集時顯示未知值可屬正常；依目前 UI 明確啟動採集，或在記錄啟動後回驗本段。不為硬套順序而直接 API 啟用，也不把尚未採集當成數值失敗。
- 五視圖：格子專注、範圍導航、水平地圖、範圍章節、格子＋詳情。
- 各驗地址、型別寬度、跨列、占用／空位、已用格編輯、空格預填後實際保存與重載。
- 預期值來自 simulator 種值／獨立協議讀回，不由 UI 本身計算。覆蓋 0、負數及目前協議可表示的型別；正式 spec 另含 bool、string、精確 uint64。
- MC／Modbus 位址與 word order 分別核對，不因兩邊 order 名稱相同就認定 wire 相同。
- 改來源值，確認 API／UI 收到新值。暫停 owned simulator，觀察 error/stale，再恢復並種另一值，排除舊快取假復原。
- 有原圖則逐點列結構、資訊、外觀與互動，附對照。功能通過不能掩蓋缺少區塊。
- 響應式／明暗矩陣依範圍列明；沒有手機原圖就只驗可用性，不虛構視覺一致。

## 4. UI 建立記錄與目的地

ID：`recording/<protocol>/<database>/record-create`。

1. 首頁新增記錄，選同一已驗設備與點位，確認 Tag／mapping／workspace 正確。
2. 透過 UI 設定本輪可丟棄 connector，核對 ID／revision。
3. 填名稱、所需欄位、週期與 layout，保存後讀回 group ID、revision、members、destination。
4. 切 A／B 設備不可互借 interval 或 members；重載後選取／設定不變。
5. 草稿保存不是 DDL、啟用或 SQL 成功。不要偷偷縮短預設週期加速；若獲准用短週期，明示不能代替原週期規格驗收。

## 5. Schema 與試寫

ID：`recording/<protocol>/<database>/schema-consent`、`test-write`。

- 比對 preview 前後 catalog／目標檔狀態，預覽不得 DDL。
- 核對同一 connector／group／revision，從 UI 明確確認 apply，不由腳本建表補過。
- 正式試寫 UI 的 preview／confirm，串 operation ID、write、獨立 readback、cleanup。
- cleanup 僅移除 operation 擁有的測試資料，保留鄰接資料。若 spec 要鄰列 fixture，按其 owned fixture 規則準備；不把試寫當正常採集。
- capability 不可用、權限不足、cleanup 未決都如實回報；不以固定成功／SQLite mock 補其他 provider。
- 正確的 unsupported 不是授權新增能力的理由。例：目前 mixed 受管寬表不提供現行試寫計畫所需 entity key，安全拒絕不代表正常記錄失敗；保留 preview 無寫入的 negative 證據，試寫執行／讀回／清理列 `BLOCKED`，正常記錄主線可繼續。不得將拒絕標為試寫 PASS，或改表／移除 ownership guard。

## 6. 啟用、SQL 與重載

ID：`recording/<protocol>/<database>/record-start`、`sql-readback`、`reload`。

1. 從 UI 檢視 scope 與影響後明確啟動，留 request／operation ID、group／connector revision。
2. running、buffer accepted、pending、SQL committed 分別觀察；不能只看綠燈。
3. 等完整窗口，用同一 monotonic clock 記開頁、啟動、closure、首次 SQL 與 polling 開銷，不用固定 sleep 假定成功。
4. 查至少三個連續完整正常 bucket；核對每個 entity／member 的值、NULL、型別、UTC 時間、quality、provenance。uint64 用十進位字串，避免 JS Number 捨入。
5. 目的地獨立唯讀 SQL：按 record ID、effect key、digest 對上 outbox、receipt、applied group revision、connector revision 與 frozen runtime version。
   - 群組交付採 `wg_delivery_outbox`／`wg_delivery_receipts`，不是 legacy `gw_delivery_receipts`；目的庫 `gw_effect_receipts` 的 effect key／digest 應對上當前真正 authority。停止或編輯後 current revision 可改變，歷史資料應核對當時 frozen runtime 與 outbox 身份，不要求現在的 applied revision 永遠等於舊值。
6. F5、關閉重開、進 runtime，身份與版本仍相同，並持續產生新的 SQL 資料。
7. ACK 遺失／partial 的額外案例，查原 operation 或使用契約允許的同一 request 恢復，不建新 request 猜成功。

## 7. 首頁、停止及 cleanup

ID：`home-stability`、`recording/<protocol>/<database>/record-stop`、`cleanup`。

- 保存前後保留首頁 DOM reference、filters／scroll／focus；計 React 重掛與主 frame navigation，不能混為 F5。
- 分開量 HTTP、保存到表單、connect/probe；一次快跑不能否定使用者較慢環境的延遲。
- 正式 UI 停止本輪記錄：新 intake 停止，已接受 backlog 按契約 drain，原 SQL／鄰列保留；停止不是刪歷史。
- 關閉自己建立的 browser、hc 服務及 worker，核對 PID／namespace；不動使用者既有服務，保留證據。

## 現有 runner 的界線

- `scripts/tests/f_device_to_sql/fresh-*.mjs`：先讀現行 source／README。歷史版本有私人 `/tmp` binary fallback，明確指定並驗證 `F_GATEWAY_BINARY`、`F_SIMULATOR_BINARY`。不能只因檔名有 fresh 就認定走目前首頁路線。
- 某些 runner 自行 spawn child；hc 管理入口，仍需核對內部埠、PID 與 cleanup ownership。
- `frontend/tests/e2e/device-recording-workflow.spec.ts` 需要 `DEVICE_RECORDING_E2E_FIXTURE`，Step 1–3 先由 API 建好；有真 SQL，仍不能證明 UI 新增設備與點位。
- guided／memory fixture E2E 可驗操作、版面與錯誤狀態，mocked API 不是真實連線、DDL 或 SQL 證據。
