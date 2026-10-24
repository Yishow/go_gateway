## Context

需求來源為本次對話：瞬時值 1–300 秒／預設 15 秒、累積值 1–600 分鐘／預設 10 分鐘，依迴路簡單選擇；累積值的原始值與區間差值「兩者都要」。這些數字是需求界線，不是效能測量。長表與可選 TimescaleDB 是本提案的技術選擇，尚未部署。

已查證的來源與限制：

- internal/datalink/workspace/write_group_managed.go 的 managedGroupTableName 由群組 ID 雜湊產生 gw_group_ 名稱；prepareManagedWriteGroup 為每個來源成員建立 v_ 欄位。點位數影響欄位數，表名不編碼點位數。
- internal/datalink/workspace/write_group_schema.go 的 managedSchemaLayout 使用文字 bucket_start 及 record_id 主鍵；不能把現有布局直接宣稱為 hypertable。
- frontend/src/features/datalink/workbench-v2/state/basicRecording.ts 的 buildBasicManagedDrafts 依設備建群組；steps/step4/BasicRecordingPanel.tsx 的 DEFAULT_INTERVAL_SECONDS 為 60，尚無本案型態週期。
- internal/datalink/workspace/write_group.go 的 WriteGroup 是設定權威；internal/datalink/grouppipeline/pipeline.go 的 Pipeline 接收 typed samples 並使用 groupdelivery worker。cmd/test_ui/service_wiring.go 已接 SampleSink 及 legacy ownership suppression；本案須延伸這條路徑。
- internal/datalink/groupdelivery/types.go 的 Closure 已表達 destination、閉窗結果與 checkpoint；新增累積基準須加入同一交易，不能在外部 SQL 成功後才另存。
- internal/datalink/aggregation/counter_usage.go 的 CalculateCounterUsage 使用 float64，單一讀值會給零差值。新精確累積契約不能直接把此結果當答案；沿用既有語意及測試案例，新增精確計算路徑，不順便改寫無關報表。
- 已讀既有 guided-recording-workflow、measurement-aggregation、recording-plan-storage、timeseries-storage 及相關封存 proposal。既有 monthly partition 規格屬舊 timeseries 路徑，本案只新增 canonical managed_long_v1。
- CLI context 的舊 Go／YAML 摘要不作為實作版本依據；以 AGENTS.md、go.mod、frontend/package.json 及實際 cmd/test_ui wiring 為準。

截圖是問題證據，並非要求逐像素重製的視覺設計來源；沿用既有 Studio 元件與互動。

## Goals / Non-Goals

**Goals:**

- 讓操作人員以迴路選資料並分別設定瞬時、累積週期，簡單流程即可準備資料庫並開始記錄。
- 累積原始值與差值可追溯、保留精度，重啟／重送不重複計量。
- 新記錄共用固定 typed long schema，支援 PostgreSQL、SQLite 與明確選用的 TimescaleDB。
- 保留舊群組與歷史；缺表恢復、版本確認、操作 ledger 及可觀察交付證據仍受原有安全邊界保護。

**Non-Goals:**

- 不做歷史寬表 backfill、既有表原地轉 hypertable、自動刪除／壓縮／保留政策、continuous aggregates、報表或計費。
- 不做速率積分、跨表／跨迴路加總、電氣拓樸推導、設備自動辨識、新採集協議、平行 V3 或另一套 scheduler/outbox。
- 不重寫既有 aggregation 全能力、legacy timeseries 月分區或本次前段已修改的 legacy dbtarget 缺表修復；只驗證與新路徑的 ownership 共存。
- 不修改 Windows tray、launcher、Step 3 mapping recovery 的獨立 WIP；不把正式設備／資料庫作為開發測試環境。

## Decisions

### 迴路型態與 canonical 群組

每個儲存的來源成員可具有 workspace-scoped circuit_id 及顯示名稱。已有明確 equipment identity 時重用；未確認的設備只能先顯示「未分配迴路」，由使用者將所選點位加入具穩定 ID 的迴路，不用名稱／地址猜測，也不新增完整設備資產系統。原始 instantaneous 記錄仍可選未分配資料，不強迫建立物理語意；cumulative 必須選定數值計量來源並確認 counter 或 signed_counter、unit、series_epoch。

一個 circuit 的 instantaneous 與 cumulative 各有獨立 canonical WriteGroup，群組內為已確認的來源成員。新增 recording_profile = {kind, interval, unit}，kind 為 instantaneous 或 cumulative，unit 分別固定 second／minute；後端正規化執行秒數但保留對外單位，拒絕分數、非有限、缺單位與越界值。新草稿提供 15 秒／10 分鐘，不重寫舊 60 秒等已存設定。

已分配迴路的建立鍵為 (workspace_id, circuit_id, profile_kind, canonical_role)。未分配的原始來源以既有 device_id 作內部 selection_scope，使用 (workspace_id, selection_scope, instantaneous, unassigned_role) 建立鍵；circuit_id 留空，畫面顯示「未分配迴路」，此內部分組不代表已確認的實體迴路。重載、重送與改名重用同群組；目的地修改沿用 CAS 和新準備／Apply，不新增第二個 writer。同 circuit 可同時含功率瞬時與電量累積；同一來源在同目的地的重複同型態群組需阻擋，使用者明確啟用兩種型態時以不同 profile 身分區隔。替代方案「整台設備只能選一種」會使同設備多電表或 kW/kWh 混合無法表達，因此不採用。

### 瞬時選樣與累積精確計算

瞬時記錄採既有 UTC 半開窗 [start,end)，取窗內 observation time 最新且品質有效、未過期的樣本；相同 observation time 以穩定 sample_id 排序。沒有合格值則寫一筆 typed NULL 與 missing/bad/stale 原因，不跨窗拿舊值補新資料。這是各成員獨立長列，不沿用寬表某一 required 成員擋住全組的規則；舊寬表 skip_row/partial 行為保留。此模式不是平均，也不調整 acquisition。

累積 raw 為窗內最新合格原始讀值；差值計算須看進入耐久 journal 的完整有效樣本序列及不連續證據，而非只比較最後兩個顯示值。baseline key 包含 workspace、group、來源 identity、definition/binding revision、series_epoch 與 profile revision。初次有值時 raw 保留、delta NULL、status baseline；這個窗不發出部分差值冒充完整結果，末筆讀值作下一窗基準。

後續窗從上次可靠基準到本窗末筆逐段計算，保存 delta_from、delta_to、baseline_sample_id、end_sample_id。只有相同 epoch／綁定、有效順序、允許的差值才累加；正常完全覆蓋輸出 ok。缺窗後同 epoch 的可靠正差可保存為跨實際時間的 span_total，status gap，不能標為單一窗完整用量或拆分填入缺窗。原始值缺失則 raw/delta 都為 NULL，baseline 不前移。

monotonic counter 任一不明下降使本窗 delta NULL、status reset_or_uncertain，保留 raw／事件並以本窗最後可靠值建立新基準；換表、epoch 或語意修訂建立 baseline，不跨界相減。signed_counter 只有經明確確認才允許負差。rollover 需 modulus、數值範圍與 max_rate × 實際時間能排除多圈的證據；只有 modulus 不足，證據不足時為 uncertain。相同觀察時間不同內容不得產生可計量差值。

整數／decimal 用精確字串與 Go 精確數值計算；來源 float 保留原 float 型別與可用精度，不偽稱能重建儀表未提供的小數。超過 JavaScript 安全整數的值不經 Number。新增精確路徑與 focused tests，避免把現有 float64 helper 全域改動。既有 lateness gate 保留；此 profile 的已閉窗 correction horizon 為零，太晚資料留診斷、不重開已交付區間，避免新增歷史重算範圍。

### 長表布局與 TimescaleDB

新增 layout_version=managed_long_v1，與缺省 managed_wide_v1 並存。每個已確認 (workspace, connector, database, schema, storage_mode) 共用一張受管長表，不按迴路、點數或週期增建資料表。table identity 與 storage generation 存在 canonical storage metadata，命名不依顯示文字。owner scope 為 workspace＋destination＋layout，所有引用群組各自持有確認版本；共享布局的變更不能替其他群組偷偷確認。

固定欄位如下，命名是本案介面契約：

| 欄位 | 意義／PostgreSQL 型別 |
| --- | --- |
| bucket_start, bucket_end | 記錄窗，TIMESTAMPTZ NOT NULL |
| record_id | 含 workspace/group/revision/source/profile/window 的穩定身分，TEXT |
| workspace_id, group_id, group_revision, circuit_id | 工作區、設定及迴路身分，TEXT |
| device_id, point_id, tag_id, measurement_id, series_epoch | 來源及計量世代；非計量原始資料允許 measurement_id／epoch 為 NULL |
| profile_kind, value_type, unit, status, quality | 記錄語意、型別及品質，TEXT |
| observed_at, received_at | 真實樣本時間，缺值列允許 NULL |
| value_int, value_decimal, value_float, value_bool, value_text | BIGINT、NUMERIC、DOUBLE PRECISION、BOOLEAN、TEXT；有效原始值只使用與 value_type 對應的一欄 |
| delta_value, delta_from, delta_to | NUMERIC、TIMESTAMPTZ、TIMESTAMPTZ；沒有可信差值即 NULL |
| provenance | JSONB，包含型別、精度、sample IDs、定義／綁定版本與異常原因 |

uint64 與精確 decimal 使用 value_decimal；SQLite 對應精確 decimal/TEXT、bool/INTEGER、時間 UTC TEXT，不能落到 REAL 丟失整數。新 long layout 的時間統一為 UTC microsecond 精度，先正規化再序列化並納入 identity；不能用送出 SQL 的當下時間取代。主鍵為 (bucket_start, record_id)；查詢索引以 (workspace_id, tag_id, bucket_start) 為基礎。receipt 與 payload digest 沿用現有交付機制。

PostgreSQL 預設 plain；明選 timescale 才讀 pg_extension、所需 function signature 及權限，準備流程使用 create_hypertable 的 by_range 介面，僅處理本次新建空表，migrate_data=false。chunk interval 初版採 7 日設計預設，在預覽中呈現並凍結，無吞吐量／最佳化保證；部署資料量若不適用，另案調整，不加入調校介面。SQLite 不接受 timescale。明選模式而能力不足回 timescale_unavailable，不自動裝 extension／改用 plain。

Timescale 時間分區需要非 NULL 時間欄位，唯一鍵須包含分區欄位；這是複合主鍵的理由，並非舊 record_id 無法維持冪等。官方參考：[create_hypertable](https://www.tigerdata.com/docs/reference/timescaledb/hypertables/create_hypertable)。測試紀錄須列實測 PG／extension 版本與 fixture image digest，不能從服務名稱宣稱版本相容。

替代方案保留每點一欄容易一起查詢固定欄位，但新增點位需要 DDL；長表避免這個耦合，代價是更多資料列、查詢顯示寬欄時要 pivot。本案不保證查詢或磁碟效能提升，也不新增 native partition 與 Timescale 雙重分區。

### 耐久閉窗與缺表恢復

沿用 Pipeline.AcceptSample → 本地 journal → boundary closure → groupdelivery outbox／receipt → 目的地 SQL。長列展開、精確計算結果、baseline checkpoint、consumed input progress、凍結 destination/schema/profile revisions 與 outbox intent 在同一交易提交。外部離線照常保留接受的資料；重啟從 checkpoint 恢復，重送用凍結 payload，不重算差值。長表每個 member 的 missing 列也是有穩定 identity 的結果。可恢復的 counter baseline 不隨已交付 journal/outbox 的本地清理一併刪除，必須保留到對應 profile 生命週期不再需要。

首次建表保持 preview → confirm operation → readiness → apply；無確認不能寫 SQL。確認時清楚顯示 recreate_missing=true 的受管表恢復政策，沿用使用者「沒有表自動重建」需求，允許明確關閉。之後只有目錄查證表缺失且該凍結布局已有有效確認、目的地版本未漂移、scope ownership 可驗證時才恢復。單一 scoped lock／operation ledger 執行 idempotent DDL、檢查實際布局，再重試 pending effect；不因一般權限／欄位錯誤觸發建表。

同名但異質表、缺失 owner 證據、extension 不可用、connector revision 變更皆 blocked 並有 repair action；不 DROP／ALTER 他人的表、不 CREATE DATABASE、不自動裝 extension、不重建消失的 SQLite 檔案。受管表 generation 記錄在本地確認證據，缺表重建後標示歷史可能遺失；已 committed receipts 不刪除或假裝重新取得已清除的歷史 payload，僅續送可確認未提交的 pending rows。

### 基本操作與前向切換

Step 4 基本區依序為目的地摘要、迴路／點位勾選、瞬時與累積開關及週期、準備預覽、確認、開始記錄。可在多選迴路套用型態／週期；提交前列出有個別覆寫的迴路並確認覆蓋，沒有勾選的迴路不改。型態不確定只提示候選，不以名稱或單位自動把 raw Tag 變 counter。

進階群組處理來源對應、signed／rollover 與既有 custom/wide 設定；新 long group 在進階顯示相同 canonical 資料，沒有第二組互相競爭的週期。工作區處理整體 readiness 與啟動範圍，並非另一種資料表。錯誤需指出是哪個迴路、來源或目的地及下一步，不重複堆疊相同 blocker 卡片。狀態分開呈現 saved、applied、collecting、locally durable、SQL committed、readback verified；SQL 寫入成功不保證 delta 已可計算。

舊群組繼續既有編輯與執行；需轉 long 時明確選來源群組、確認新布局並執行 CAS cutover。cutover time 為已存不可變邊界，舊窗在邊界截斷並標 partial，新窗由邊界開始至下一對齊窗界亦標 partial；同一來源時間只歸一個 active revision。舊接受資料依原 destination 送完，新 cumulative 第一窗重新建立 baseline。切換失敗保留舊 applied revision，不雙寫；回復也以另一個明確邊界恢復舊布局，不能覆寫過去 snapshots。

## Implementation Contract

- 範圍：上述新 profile、長表、可選 Timescale、缺表恢復與 Studio 操作；非目標不得因改共用模組而自動擴增。
- JSON：WriteGroup 擴充 recording_profile、layout_version、storage_mode 與 member circuit_id；storage metadata 保存 layout digest、generation、recreate_missing、confirmed scopes。舊 payload 缺少新增欄位仍讀為既有 wide，不偷偷套入新預設。API／TS types／hydration 同步支援。
- 錯誤：使用穩定代碼 profile_interval_invalid、counter_semantics_required、timescale_unavailable、storage_owner_mismatch、schema_repair_blocked 及原有 revision conflict；UI 翻譯為欄位／迴路修復文字，raw SQL、DSN 留遮蔽診斷。
- 時間驗收：瞬時 1/15/300 秒通過，0/301/小數拒絕；累積 1/10/600 分鐘通過，0/601/小數拒絕；慢採集時缺值可見且不變更 scheduler；邊界樣本、亂序、late 與重送有決定性結果。
- 計量驗收：整數 9007199254740992→9007199254740993 差為 1；decimal 12000.0→12000.4→12000.7 差為 0.4／0.3；首窗 NULL、reset NULL、已確認 signed 負差、可證單圈 rollover、缺窗 span_total、新 epoch baseline 均有固定預期值。
- 耐久驗收：journal ACK 後、閉窗交易前後、SQL commit 後 ACK 前各插入 crash；同一 effect 最終只有一列且 baseline 不重複前移。與舊 wide、custom、legacy ownership 共存測試必須過。
- 儲存驗收：可丟棄 SQLite 與普通 PostgreSQL 做建表、寫入、readback、刪除 fixture 表後恢復；Timescale fixture 做 extension capability、hypertable catalog、跨 chunk 寫入及複合唯一鍵驗證。無 extension、權限不足、同名異質表、operation 回應遺失有明確失敗／恢復測試，不把 skip 算 PASS。
- UI 驗收：先補失敗案例再實作；實際瀏覽器驗證迴路批次選擇、雙型態、reload、越界錯誤、確認及切換，交付截圖／操作證據，區分 local durable 與 SQL commit。
- 每個 implementation task 的預算為估計 30–60 分鐘，以該 task 開始至 focused verification 的實際時間驗證；達 60 分鐘尚未可驗收就停在可建置切片並細分剩餘工作，不以預估作效能或完成保證。

## Risks / Trade-offs

- [更多資料列與離線積壓] → 批次 transport 不改 record semantics，沿用現有 queue 上限及可觀察 backpressure，不宣稱採集間隔等於 flush。
- [歷史規格有較廣 aggregation 承諾] → 本案只承接明確 counter 契約與零 correction horizon；不宣稱報表、速率積分或所有 recording-plan modes 已完成。
- [長表共享 ownership 及缺表 receipts] → scope storage generation 與群組個別確認並存；只恢復表結構與未提交資料，明示舊歷史不保證復原。
- [部署 Timescale 版本未知] → 開發使用可丟棄具 by_range 能力版本並記錄 image digest；正式環境先唯讀能力查核，未通過就不啟用該模式。
- [迴路名稱不等於電氣拓樸] → 僅由明確來源 membership 定義，沒有自動總表／分表加總。

## Migration Plan

1. Expand：增加 nullable metadata／版本化 layout，舊 payload、snapshots、DB rows 可讀，保留舊基本入口與所有舊 writer 契約。
2. Migrate：先在 fixture 驗證新建雙 profile 的 plain 長表，再增加 Timescale 與缺表恢復；所有批次保持建置與 focused tests 通過。
3. Migrate：提供明確、可回復的 forward cutover，驗證舊 pending rows 排空與新 baseline；不轉歷史表。
4. Contract：新建基本流程只提供 profile 選擇的長表預設，移除該新流程重複的舊週期控制；既有 wide/custom 編輯器、parser、runtime 保留。
5. 回復：停用新 profile，依新 effective boundary 恢復先前可用群組；已接受新資料繼續按其 frozen destination 排空，不 down-migrate 或刪表。

## Open Questions

無需使用者再決定的需求。fixture 的 PostgreSQL／TimescaleDB 版本、extension 能力與實際驗收結果已記錄於 docs/technical/studio-v2-circuit-recording-db-verification.md 及瀏覽器驗收文件；提案時尚未執行驗收，不代表後續任務仍未完成，也不授權啟用正式環境。
