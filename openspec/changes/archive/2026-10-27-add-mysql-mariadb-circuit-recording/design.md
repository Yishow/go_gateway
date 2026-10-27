## Context

使用者要求 MySQL／MariaDB 共用一份 change，SQL Server 另開，之後分批實作。本案為 feature/spec-driven；目前已授權分批實作，並追加最小精確 decimal 來源擴充。

必要前置 add-circuit-recording-profiles 已於 35390b4 封存；其 managed_long_v1、穩定 record identity、精確計算、baseline transaction、缺表 ownership 及 Basic UI 必須先實作驗證。本案唯一擁有 mysql/mariadb adapter 延伸，不複製 profile 或 counter 引擎，也不要求修改前置提案的 PostgreSQL／SQLite 範圍。

本輪已開啟查證：internal/datalink/recordingplan/capability_mysql_test.go 的 TestGetConnectorCapability_MySQLKeepsCustomTableWhileManagedSchemaIsUnverified 與 GenerateManagedSchemaDDL 明確拒絕未驗證 managed DDL；dbtarget/destination_open.go 的 OpenDestination 僅允許 SQLite／PostgreSQL；groupdelivery/sender.go 的 Sender.Deliver 使用凍結 payload、receipt 與 InsertGroupRow。KindSelector.tsx 已有 MySQL/MariaDB 顯示，但不能當作完成證據。已讀 database-output-delivery、recording-database-setup、datalink-workbench-v2-step4-database 及已封存 fix-activation-barrier-and-mysql-target-defects；新受管路徑不以舊 custom/schema 修復代替驗證。

## Goals / Non-Goals

**Goals:** MySQL 與 MariaDB 使用相同操作流程及邏輯長表；普通 InnoDB 為基本路徑、原生時間分區為可選；完整保留 schema 確認、可靠交付、缺表恢復及分別驗收。

**Non-Goals:** SQL Server、其他 storage engine、Galera／多主叢集、歷史 backfill、既有表原地分區、跨庫 migration、報表、保留刪除、壓縮或安裝 TimescaleDB。既有 custom mode 行為與未相關的 TLS／帳號機制不做全面翻修。

## Decisions

### 共用流程與分開的伺服器能力

沿用 connector kind=mysql，不建立競爭的 mariadb 設定模型。新增唯讀 adapter projection：server_family=mysql|mariadb|unknown、server_version、capability_revision、supports_managed_long、supports_native_partition、supports_group_test_write、supports_missing_repair。這些是後端探測結果，不能由 client 宣告放行；family/版本無法辨認回 unverified，不能以字串包含 MariaDB 當成所有能力已測過。

初始驗證目標為 MySQL 8.4 與 MariaDB 11.4，這是設計測試基線，不代表現場版本或其他版本已支援；交付須記錄實測 patch、engine、sql_mode 與 fixture digest。相容性以版本化 capability registry 及 operation 檢查決定，不把 MySQL 通過推及 MariaDB。連線成功、能讀 metadata、能建表與能寫入分開回報。新受管連線需既有安全傳輸選項通過，不自動降級 TLS 或開啟明碼驗證。替代方案完全獨立 adapter 會複製共同邏輯，因此只將 server-family 差異封裝於既有外部 SQL 介面。

### 固定長表與精確 MySQL 型別

沿用前置全部邏輯欄位、record_id、NULL/status 與 provenance；按 workspace＋connector＋database＋storage_mode 共用表，MySQL database/schema 為同一 scope，不能套用 main/public。SQL identifier 以 dialect quoting，值用參數綁定，避免使用者欄位拼接 SQL。

時間使用 DATETIME(6)，gateway 統一 UTC microsecond 後綁定／讀回，不能依連線時區改變窗界；IDs 使用有界且大小寫精確的字串，record_id 固定 hash 字串。主鍵 (bucket_start, record_id)，查詢索引 (workspace_id, tag_id, bucket_start)。value_int BIGINT、value_decimal/delta_value DECIMAL(65,30)、value_float DOUBLE、value_bool BOOLEAN、value_text LONGTEXT；JSON provenance 在兩種 engine 各用可驗證的 JSON 表達，digest 依 canonical application payload，不比較 JSON 儲存位元順序。

DECIMAL 的 precision/scale 是本 layout 設計界線；uint64 可精確保存，超出可表示整數位或小數位的精確值在 intake 前回 exact_value_unrepresentable，保留原始型別診斷，不默默 round 或轉 DOUBLE。型別／collation／engine／索引均納 schema digest。不同來源同名、大小寫不同 ID 與 Unicode 必須 readback 後仍可區分。改成全部字串容易保存但破壞 typed 查詢契約，因此只在共同邏輯指定的文字欄位使用文字。

### 逐句建表恢復與耐久交付

首次準備沿用 preview → explicit confirm → operation ledger → inspect → readiness → Apply。MySQL/MariaDB DDL 不當成可整批回復的交易：每個步驟存 pending/started/verified/failed，崩潰後以同 operation_id、ownership 與實際 catalog 驗證，已完成步驟不重建，未完成步驟可續行。表存在但 fingerprint 不符直接 blocked；索引按 catalog 查證後執行合法 SQL，不使用不相容的 CREATE INDEX IF NOT EXISTS，也不刪除已完成物件假裝 rollback。

資料列與 receipt 則必須同一 InnoDB transaction。沿用既有 outbox／lease／effect key；相同 key＋digest 重送回已提交，不同 digest 為衝突。不能用 INSERT IGNORE、REPLACE 或無條件 upsert 吞掉資料錯誤或覆蓋既有值。SQL commit ACK 遺失依 receipt 查證，不能重算 counter。

recreate_missing=true 僅允許已確認 scope/generation、目錄確證缺表且 connector revision 相同的受管表；scoped coordination 經 ledger 恢復 schema 後續送可確認未提交資料。table missing 與 metadata forbidden、auth failure、欄位不符分開；不 CREATE DATABASE、不改造異質同名表。receipt 不隨缺表刪除，UI 明示歷史資料不保證恢復。DDL restore 不能插入 row＋receipt 的 SQL transaction 內。

### 原生月分區與有限維護

storage_mode=plain|native_monthly，預設 plain。native_monthly 明選並經 capability/權限查核後，以 bucket_start 的 RANGE COLUMNS(DATETIME(6)) 建月分區；所有 unique keys 都含 bucket_start，不建立分區表 foreign keys。receipt 另存普通 InnoDB 表，避免以分區歷史維護影響去重。

初次確認預覽以明確 UTC 月界建立當月及後續 3 個月，前端顯示準確 boundaries，並保留 past/catch-all MAXVALUE 區間容納晚送資料。3 個月是設計維護視窗，不是效能估算。現有 worker cadence 執行可重入 maintenance，保持後續 3 個月；只在已授權 owner 範圍 REORGANIZE catch-all 新增邊界，保留既有行，不 DROP/TRUNCATE/MERGE 分區。停機很久不補出每個錯過的月，將舊 catch-all 資料保留，再從當月建立有限未來視窗，避免一次無界 DDL。

maintenance failure 留下可見 partition_maintenance_blocked 與 repair；既有 MAXVALUE 仍可正常寫入時繼續交付，不偽稱維護完成。沒有 catch-all 或布局不符則阻擋該目的地並保留 outbox。DDLs 使用 scoped lock、lease、catalog 查核與 bounded lock timeout，失敗稍後重試，不長時間鎖住所有群組。

官方依據：[MySQL 分區鍵](https://dev.mysql.com/doc/refman/8.4/en/partitioning-limitations-partitioning-keys-unique-keys.html)、[MySQL DDL 隱式提交](https://dev.mysql.com/doc/refman/8.4/en/implicit-commit.html)、[MariaDB 分區限制](https://mariadb.com/docs/server/server-usage/partitioning-tables/partitioning-limitations)。兩個 family 各自驗收，不能套用 Timescale 語法。

### 分批開放與既有模式相容

**Supersedes**: add-mysql-mariadb-circuit-recording / 分批開放與既有模式相容（原 legacy 501 保留條款）

2026-10-25 使用者核准銜接已合併的 unify-device-recording-workflow 安全轉接契約；此決策只替換原 501 條款，其他範圍與驗收保持。

批次順序：共同迴路提案 → 本案 MySQL/MariaDB → SQL Server。後兩者因共用 capability/schema/UI 檔案採序列，不形成 SQL 方言相互依賴。本案先 expand plain adapter，再驗證 MariaDB family 與 native partition，最後解除有實測證據的 managed-long gates。group test-write 使用既有狹窄 group 流程做真實 write/readback/cleanup；legacy recording-plans/test-write 保留新 HEAD 既有的安全轉接契約：沿用共同 preview token／operation_id 確認、scope／revision 查證及各別 write/readback/cleanup 結果，不回退為 501。

Step 4 顯示 MySQL/MariaDB family 與可用操作，選迴路後使用共同雙 profile 與相同預設。schema 部分成功以同 operation_id 修復；錯誤是可翻譯安全 code，不輸出 DSN。切換目的地走共同 reviewed/CAS boundary，舊 pending rows 排原目的地，不雙寫；停用新 intake 後可排空，不能刪表當作回復。

### 已確認十進位 scale 的最小來源擴充

依使用者追加授權，保留原精確 decimal raw/delta 驗收。MeasurementDefinition 新增可選 `exact_recording_scale` 十進位文字；空值維持舊行為。只接受正值有限十進位且 NumericEncoding 與實際 Tag 均為整數來源的已確認 counter/signed_counter；不由 float、名稱或單位推導，不改寫既有 mixed-layout float scale，也不新增設備 decimal datatype、offset 或第二套 counter 引擎。

保存及讀回 scale 使用文字與既有精確運算，無 float64 中介。變更 scale 須提升 definition revision、開新 series epoch，並使舊確認失效；runtime 凍結 scale、原始整數型別、measurement/binding revisions 與 epoch。sampleFrom 先依原始型別驗證 envelope，再僅一次精確縮放為 ExactDecimal，於 journal 前通過目的地可表示性檢查；瞬時 counter_raw 與累積共用此路徑。沒有 scale 的整數、float 與 Modbus 輸出維持既有語意。

已接受 payload 與 baseline 保留原 proof，不因定義編輯重算；新的定義走既有 reviewed/CAS forward boundary，禁止跨 epoch 比較。整數 120000→120004 搭配 scale `0.1` 是契約案例，產生精確 raw 12000.0→12000.4、delta 0.4；舊 float 誤差不可四捨五入冒充本驗收。共同程式與驗證唯一歸 add-mysql-mariadb-circuit-recording 的 5.1／5.2；SQL Server 批次只消費同一契約與自己的真實資料庫證據。

## Implementation Contract

- 追加來源契約：`exact_recording_scale` 為可選正值十進位文字；僅已確認整數 counter 可用。非法 scale／非整數 encoding／原始型別錯誤／stale revisions 不得產生可用 decimal；未設定保持舊行為。measurement 保存、epoch transition、group 型別驗證、frozen proof、pre-journal conversion 與回歸測試納入範圍，不擴充設備協議或另建計算引擎。驗證 TestExactRecordingScaleDefinition、TestExactRecordingScaleEpoch、TestExactRecordingScaleBoundary；原兩案真實 DB decimal raw/delta 仍須各別驗收。

- 介面：既有 mysql kind＋上述 server capability projection、managed_long_v1、plain/native_monthly、recreate_missing；必須保存 family/version-bound schema digest 與 group/connector/workspace revision，不接受 client 自稱已驗證。
- 失敗：adapter_unverified、exact_value_unrepresentable、schema_operation_partial、storage_owner_mismatch、partition_maintenance_blocked；沿用既有 revision conflict／target blocked/outage，區分 missing/forbidden/failed。
- 驗收：兩種真實 fixture 各跑 profile 週期界線、raw/delta、uint64 18446744073709551615、9007199254740992→9007199254740993、decimal 12000.0→12000.4、NULL/bool/Unicode／UTC readback，拒絕超界精確值。
- 故障驗收：每步 DDL 後中斷、operation 回應遺失、SQL ACK 遺失、並發重送／digest 衝突、permission denied、缺表／同名異質表、跨月／停機超出預建窗；比較 catalog、rows、receipt、outbox 的 scope 與實際結果。
- 產品門檻：兩 family 的 plain/native_monthly、schema恢復與 group test-write 均有實際證據後才開完整能力；未驗證版本維持 gate，不將 skipped 或 driver 存在算 PASS。
- 範圍內只延伸資料庫 adapter、安全確認與必要 UI；除追加授權的 scale 來源接點外，共同計量與舊 custom 設定不重寫。所有 proposed Test* 由 tasks 新增，非現有測試聲明。
- 任務估計各 30–60 分鐘，實際從開始至 focused verification 記時；超過 60 分鐘未能交付可建置切片則拆分剩餘任務。這是規劃估計，不能當完成證據。

## Risks / Trade-offs

- [同 protocol 但 DDL/JSON/collation 不同] → family/version capability 與各別真實 fixture 矩陣。
- [分區重整有鎖成本] → 可選、有限視窗、catalog 查核與 timeout；不宣稱分區必然較快。
- [DDL 部分完成不可整批 rollback] → 以既有 ledger 記錄並恢復，readiness 直到所有必要步驟實際 verified 才通過。
- [exact decimal 範圍小於 PostgreSQL 無界 NUMERIC] → intake 前檢查可表示性，明確錯誤、不靜默 round。

## Migration Plan

Expand 新增 family/capabilities 與 plain layout，保留未驗證 gate；Migrate 以新 canonical groups 明確選用、經 preview/confirm/apply 啟動，驗證後加入 native_monthly；Contract 只解除本案已驗證版本／操作的 gate，保留舊 custom editor、parser 與 frozen outbox。每批可建置並通過 focused tests。回復只停新 intake 並排空已接受資料，依共同 forward boundary 回到先前配置，不刪表或轉換歷史。

## Open Questions

沒有需要使用者再選的需求。後續實作前須準備兩種可丟棄 fixture 並記錄版本／digest；未具環境時資料庫驗收為 NOT RUN，本輪依授權與 tasks 繼續實作。
