# 共用契約 v1

本文件描述擬實作的契約，非現行 API。各 change SHALL 以此為共享資料詞彙，
修改共享欄位時須同時更新受影響的 delta、測試與 API 文件，不能各自發明相似欄位。

## 1. 三層分開

來源：設備、協議、位址、占用長度、解碼格式與現有 Point/Tag/Mapping。
量測：這個項目是什麼、單位、用途、身份與計量世代。
記錄：明細／摘要／事件／批次／最新值如何保存、算哪些結果、送到哪裡。
通訊合併不等於資料庫分組；資料型別不等於資料意義；寫入間隔不等於採樣間隔。

## 2. 量測與樣本

MeasurementDefinition 至少包含 `measurement_id`、`workspace_id`、`device_id`、
`point_id`、`tag_id`、`equipment_id`、`definition_revision`、`source_binding_revision`、
`quantity`、`unit`、`semantic_kind`、`numeric_encoding`、`counter_policy` 或 `state_map`。
`semantic_kind`：`gauge`、`rate`、`counter`、`signed_counter`、`delta`、`state`、`event`、`text`。
`quantity` 識別溫度／壓力／電流／電壓／功率／能量／體積／產量等物理意義。
同一個 unit 不足以決定用途；未確認用途只可留原值，不能自動啟用加總或積分。
更名不換 measurement_id；換表、換位址、變倍率／單位等破壞可比性時新增 `series_epoch`。
版次不覆寫舊定義；既有歷史保存所用版次，不能拿新版倍率重解舊資料冒充原紀錄。

SampleEnvelope 至少包含 `sample_id`、`workspace_id`、`measurement_id`、`series_epoch`、
`definition_revision`、`source_binding_revision`、`acquisition_id`、`source_sequence`（可空）、
`observed_at`、`received_at`、`time_origin`、`value_type`、`value`、`quality`、`quality_reason`。
有可信設備時間才採設備時間；否則由 Gateway 在取得回應時標時，不宣稱是 PLC 掃描時間。
每次真實採集產生穩定 sample_id；重送沿用，不以 timestamp 單獨去重。
品質：`good`、`missing`、`stale`、`invalid`、`uncertain`；估算是結果的獨立標記，不冒充 good 實測。
quality 非 good 不參加一般統計；原始值與原因仍保留。單項失敗不污染其他項目。
int64/uint64 及精確小數跨 JSON 使用型別標記與十進位字串，禁止經 JS Number 破壞精度。
保留解碼原值與換算後值的關聯；原始封包是否留存另有配額，不強制永久留封包。

## 3. 記錄方案

RecordingPlan：`plan_id`、`workspace_id`、`revision`、`applied_revision`、`status`、
`timezone`、`members`、`streams`、`destinations`、`retention`、`limits`。
狀態至少區分 draft、validating、ready、running、partial、blocked、paused。
member 明確引用 measurement_id 與有效定義；stream 明確引用 member 與獨立用途。
stream 模式：`raw_history`、`window_summary`、`usage_interval`、`state_changes`、
`event_log`、`batch_snapshot`、`latest_only`。一個項目可被多個 stream 重用，採集只做一次。
`raw_policy`：`every_sample`（預設）、`on_change`、`sampled`；後兩者明示不能還原完整原始歷史。
`on_change` 配絕對死區與最大心跳間隔；品質變化不得被死區隱藏。
摘要與用量使用送入耐久入口的有效樣本，不使用被降採樣後的「最後一筆」代替全部樣本。
明細降採樣時，另保留計算必要證據至 correction_horizon，容量估算要計入這份證據。
raw／summary／events／必要計算證據／本機待送佇列的保存期限各別設定。
初始建議每項 max_hold 與 integration_max_gap 都為該項有效採樣間隔的2倍；
lateness_grace 為 max(10秒,2倍採樣間隔)，correction_horizon 為24小時，且不得超過必要證據保留期限。
全部以實際 plan 設定保存並在確認摘要顯示，使用者修改後必須重新驗證。
示意預設：讀取 5 秒、摘要 1 分鐘、明細 30 天、摘要 365 天；不是設備規範或容量保證。
容量以實際長表筆數、索引、計算證據、摘要與目的地份數估算；100設備×8項每5秒，
僅raw即13,824,000筆/天。容量檢查未通過不得直接套用30天明細建議。

## 4. 儲存與提交邊界

預設 managed 模式採可擴充的逐項長表，不依不同設備新增一堆專用表：
`gw_record_samples`、`gw_record_intervals`、`gw_record_events`、`gw_record_snapshots`、
`gw_record_definitions`、`gw_record_receipts`。實際命名由安全前綴與 schema plan 產生。
唯一身份至少包含 workspace、plan、stream、record_id；不能只用名稱或時間當主鍵。
值使用 typed columns；精確數字 canonical decimal text 是跨 adapter 的保存底線，
可另有近似數值索引／native numeric 欄，但報表計量不得取近似欄代算。
bool、enum、bitmask、text 原型別保存；無值與數字 0 不同，NaN/Inf 標 invalid。
歷史 append；摘要以穩定 result_id + calculation_revision 修訂；latest_only 條件更新，不覆寫較新資料。
索引至少支援 workspace/plan/equipment/measurement + observed_at 與結果區間查詢。
wide/grouped 表僅是進階相容輸出或轉置 view，不是第二套計算來源。
既有表不自動增刪欄位或覆寫資料；無去重能力不得標記可靠重送可用。

耐久順序：採集樣本落本機 journal → 計算 checkpoint 與輸出記錄原子提交 → 外部送出 → 回執。
2 定義儲存／transaction 介面；3 消費該介面；4 補齊 journal、outbox、worker 的生產可靠性。
SQLite/PostgreSQL/MySQL 需各自通過 managed schema、精度、transaction、receipt 整合測試才可啟用。
SQL Server 及其他目的地需能力測試通過才顯示可用；本組不承諾新增 driver。
runtime 不自動建立消失的資料庫，不把採集正常當成交付成功。

## 5. 時間與計算

時間儲存 UTC；視窗以方案 IANA timezone 切日，台灣方案預設 Asia/Taipei。
視窗使用 `[start,end)`；邊界讀值可供兩側差分，但不可重複計量。
SummaryResult 包含區間、mean/min/max/count、valid_duration、expected_duration、
quality/coverage、provisional/final、calculation_revision、evidence refs。
UsageResult 包含起訖讀值／時間、quantity_delta、known_subtotal、complete、
allocation_status、estimated、series_epoch、evidence refs、calculation_revision。
不完整期間的完整用量為 null；可另呈現已知小計，不能把小計標成整日總量。
跨分鐘但不能分配的累積差值，保留整段區間，不硬塞到最後一分鐘或自動平均分配。
有限亂序寬限與 correction_horizon 都在方案明列；超過證據保存範圍的補算拒絕並說明原因。

## 6. API 邊界與共同錯誤

重用既有 `/api/v1/datalink/studio-v2/workspace` ownership 與 API envelope。
擬新增：`/measurements`、`/recording-plans` 的查詢與 revision-aware 更新，
`/recording-plans/{id}/validate`、`/preview`、`/schema-plan`、`/schema-apply`、`/test-write`、
`/recording-history/samples`、`/intervals`、`/events`、`/export`，以及 recording delivery status。
不能新增另一條繞過 workspace activation barrier 的 `/activate`；使用既有 workspace activate，
攜帶選定 plan revisions，等待 autosave，檢查來源／目標／settings revision 與 readiness token。
狀態回應分開 configured、applied、collecting、local_durable、delivered、verified。
保存設定不代表已生效；執行中草稿需明確確認，原 applied revision 持續服務至安全切換。

schema/test preview token 綁定 ownership、connector identity revision、plan revision、
mapping signature、TTL；過期或欄位改動回 typed 409，不能拿舊預覽寫新位置。
一般欄位錯誤 422；ownership 使用既有 403/404 策略；資源不可用 503。
錯誤包含 code、field/path、retryable、action、opaque request_id；正常 UI 不露 DSN、密碼或 raw exception。
外部連線由 Gateway 主機發起；沿用既有密碼以後端 credential reference，不回傳明文。
狀態、確認預覽、啟動與 schema 寫入皆由後端驗證，不信任瀏覽器 disabled 按鈕。
