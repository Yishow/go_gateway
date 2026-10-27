## Context

本案依使用者要求獨立規劃 SQL Server，與 MySQL/MariaDB 各有 change，後續分批執行。必要共同前置 add-circuit-recording-profiles 已於 35390b4 封存；本案不能取代它的 profile/counter/baseline 工作。

已查證 internal/datalink/recordingplan/service.go 的 GetConnectorCapability 對 sqlserver 走不支援分支；dbtarget/destination_open.go 的 OpenDestination 只接受 SQLite/PostgreSQL；recordingplan/schema_generator.go 未提供 SQL Server managed 方言；groupdelivery/sender.go 的 Sender.Deliver 已有 frozen payload、receipt 與 InsertGroupRow 接點。Step 4 KindSelector.tsx 有 SQL Server 選項，不能推定實際支援。相關既有規格為 recording-database-setup、database-output-delivery 與 datalink-workbench-v2-step4-database，沒有把其所有資料庫能力視為已完成。

## Goals / Non-Goals

**Goals:** SQL Server 經真實查證的連線、metadata、受管長表、精確 raw/delta、schema confirmation、production durable delivery、缺表恢復及可選原生時間分區。

**Non-Goals:** 新計量引擎、MySQL/MariaDB 實作、Azure SQL／Synapse／Fabric 額外相容承諾、Windows Integrated／Entra 驗證、SSIS、報表、跨庫寫入、歷史 backfill、舊表原地分區、自動保留刪除、partition switch/archive、建立資料庫／filegroup／檔案或全域 DBA 排程。

## Decisions

### SQL Server 連線與能力範圍

沿用 connector kind=sqlserver 與既有 saved connector identity；實作前核對實際外部 adapter，能重用就擴充，否則在既有 database/sql abstraction 加入 Microsoft go-mssqldb 的必要依賴，不新增另一套 connection manager。官方 driver 參考：[Microsoft Go driver](https://learn.microsoft.com/en-us/sql/connect/golang/microsoft-go-mssqldb-driver)。

基本介面為 host/port、database、schema（新草稿 dbo）、SQL authentication 帳密與既有 TLS 選項；預設加密並驗證 certificate，測試 fixture 也使用明確 CA，不偷偷信任任意憑證。連線測試由 backend 發起且唯讀；密碼遮蔽及 revision 比較沿用，不能 trim 真實密碼。database 與 schema 分開解析，identifier 用 SQL Server quoting，值用 @pN 或 named parameters，禁止使用 PG placeholder 或 MySQL quoting。

原先設計測試基線為 SQL Server 2022；依下方「實際環境的功能驗收」決定，本輪使用可用真實 SQL Server 環境驗證功能，如實記錄版本與架構，不代表未測版本或其他 edition 自動相容。probe 取得 product version、engine edition、selected DB/schema、operation permissions 與 capability_revision。連線成功、metadata可讀、DDL可用、transaction/receipt、partition可用分別呈現；未驗證版本/操作返回 adapter_unverified，不由 UI 勾選宣告支援。

### 精確長表與交易內 receipt

同一 managed_long_v1 邏輯欄位與固定 source/profile/window identity，scope 為 workspace＋connector＋database＋schema＋storage_mode。SQL Server 映射：時間 datetime2(6) 並約定 UTC；signed integer BIGINT、uint64及decimal用 DECIMAL(38,18)、浮點 FLOAT(53)、bool BIT、文字 NVARCHAR(MAX)、provenance NVARCHAR(MAX) 且驗證 JSON。record_id 為固定長度 hash、source IDs 為有界 Unicode 文字並採明確大小寫／binary identity比較，不依預設不分大小寫 collation 合併來源。主鍵 (bucket_start, record_id)，source/time index 依共同契約且限制索引鍵寬。

DECIMAL(38,18) 保留 20 位整數，因此能保存 uint64 最大值；超過整數／小數表示範圍在 intake 前回 exact_value_unrepresentable，不能讓 DB 默默 round。應用程式精確計算並以可驗證參數型別綁定，不以 SQL 浮點差分代替共同 counter 計算。UTC 先正規化到 microsecond 再納 identity，禁止用 SQL 的 timestamp/rowversion 表示事件時間。官方 precision 依據：[Microsoft decimal precision](https://learn.microsoft.com/en-us/sql/t-sql/data-types/precision-scale-and-length-transact-sql)。

同一 database transaction 寫資料列＋receipt；按 effect key 與 payload digest 區別重送／衝突，transaction-safe 鎖定及唯一鍵處理並發。不能單靠 IF NOT EXISTS 檢查就假設沒有競爭，也不以無條件 MERGE 覆蓋既有歷史。SQL commit ACK 遺失由 receipt 查證；重試用凍結 payload，counter baseline 不重算。deadlock/timeout 分類為保留資料後有限 backoff，型別／identity 衝突隔離，不阻塞別的目的地。

### 確認操作與受管缺表恢復

依前置 preview → explicit confirm → operation ledger → actual inspection → readiness → Apply；token 绑定 workspace/group/connector/layout 與 schema digest。metadata 使用 sys catalog／INFORMATION_SCHEMA，object invisible 不等於不存在：查證權限與 scope，無證據回 forbidden/unconfirmed。

本案建立的 schema/table/index/partition 物件使用 owner metadata，preview 完整列出新增物件；使用 SQL Server 支援的交易執行 DDL，但仍保存同 operation_id 與前後實際狀態，處理 client 回應遺失、deadlock與不確定 commit，不能只靠 client rollback 推定全部復原。部分物件已存在先驗證 ownership及完整 fingerprint，不 DROP/ALTER 異質物件；token stale 需重新預覽。

recreate_missing=true 且既有確認、generation、connector revision 皆有效，目錄確證本案 table 缺失時才由 scoped coordination／ledger 恢復；receipt 與資料列分開檢查，不丟棄 committed receipts。若 partition function/scheme 遺失或異質，只能恢復本案專屬且可證明未被外部物件共用的布局，否則 blocked。僅續送可確認未提交的 outbox，顯示原有歷史可能遺失；不 CREATE DATABASE 或變更外部 DBA 物件。

### SQL Server 原生月分區

storage_mode=plain|native_monthly，預設 plain。明選 native_monthly 後使用專屬 partition function（datetime2(6)、RANGE RIGHT）＋partition scheme，初版全部映射既有 PRIMARY filegroup，不管理實體檔案。clustered primary key 與相關索引對齊 partition scheme，unique keys 含 bucket_start。receipt 保持普通表，不隨 history 分區管理。官方參考：[partitioned tables and indexes](https://learn.microsoft.com/en-us/sql/relational-databases/partitions/partitioned-tables-and-indexes?view=sql-server-ver16)。

首次確認列出當月及後續 3 個月 UTC 月界，保留過去／未來 catch-all；3 個月為設計維護視窗，不是效能估計。既有 worker 執行有限、可重入的維護，只 SPLIT 新邊界並設定 NEXT USED 指向 PRIMARY，以預建後續 3 個月。停機很久只補當下有限未來視窗，舊廣域區間原樣保留，不無界補齊中間月。任何已有資料的 split 必須 timeout-bound，保留資料，失敗留可見 maintenance 狀態與下一次修復，不重建全表。

無分區權限仍能選 plain；明選而不可用回 partition_unavailable，不降級或安裝其他元件。分區維護失敗但 catch-all 可寫時，保持資料交付並顯示 partition_maintenance_blocked；目錄不符或不能證實安全寫入時保留 outbox並阻擋此 scope。不 DROP/TRUNCATE/MERGE、無自動刪除／壓縮政策，不借用／修改 DBA 的共享 partition function。

### 分批開放與操作一致性

**Supersedes**: add-sqlserver-circuit-recording / 分批開放與操作一致性（原 legacy 501 保留條款）

沿用使用者 2026-10-25 已核准銜接新 HEAD 的共同安全轉接決定；本批只替換失效的 501 條款，保留原 SQL Server 範圍與全部資料庫／瀏覽器驗收。

實作順序為 add-circuit-recording-profiles → add-mysql-mariadb-circuit-recording → 本案；順序用來避免共同 capability/schema/Step 4 檔案衝突，SQL Server 技術上只依賴共同長表契約。前批實作完成後再核對實際介面與版本，不預先複製尚未存在的 helper。

Step 4 使用同一 circuit profile UI，瞬時 1–300 秒／15 秒及累積 1–600 分鐘／10 分鐘預設不變，累積同存 raw/delta。UI 明確呈現 connected/inspected/prepared/applied/durable/SQL committed/readback，各操作按實際 capability gate；group test-write 用狹窄已確認 scope 真實寫入、readback、cleanup並呈現各結果。legacy recording-plans/test-write 沿用現有共同安全轉接契約：preview token／operation_id、scope／revision 查證與各別 write/readback/cleanup 結果；本案不回退 501，裸請求與 stale token 仍拒絕。

目的地修改採共同 reviewed/CAS forward cutover，已接受資料排原目標；停用新 intake 仍保留 outbox，不刪表回復。普通與分區皆完整驗證後才開對應 capability，先有 driver 或 static SQL 測試不算完成。

### 已確認十進位 scale 的最小來源擴充

依使用者追加授權，保留原精確 decimal raw/delta 驗收。MeasurementDefinition 新增可選 `exact_recording_scale` 十進位文字；空值維持舊行為。只接受正值有限十進位且 NumericEncoding 與實際 Tag 均為整數來源的已確認 counter/signed_counter；不由 float、名稱或單位推導，不改寫既有 mixed-layout float scale，也不新增設備 decimal datatype、offset 或第二套 counter 引擎。

保存及讀回 scale 使用文字與既有精確運算，無 float64 中介。變更 scale 須提升 definition revision、開新 series epoch，並使舊確認失效；runtime 凍結 scale、原始整數型別、measurement/binding revisions 與 epoch。sampleFrom 先依原始型別驗證 envelope，再僅一次精確縮放為 ExactDecimal，於 journal 前通過目的地可表示性檢查；瞬時 counter_raw 與累積共用此路徑。沒有 scale 的整數、float 與 Modbus 輸出維持既有語意。

已接受 payload 與 baseline 保留原 proof，不因定義編輯重算；新的定義走既有 reviewed/CAS forward boundary，禁止跨 epoch 比較。整數 120000→120004 搭配 scale `0.1` 是契約案例，產生精確 raw 12000.0→12000.4、delta 0.4；舊 float 誤差不可四捨五入冒充本驗收。共同程式與驗證唯一歸 add-mysql-mariadb-circuit-recording 的 5.1／5.2；SQL Server 批次只消費同一契約與自己的真實資料庫證據。

### 執行與驗收分開推進

**Supersedes**: add-sqlserver-circuit-recording / 分批開放與操作一致性（待辦任務須先完成真實fixture才可繼續本地實作的執行順序）

依使用者2026-10-26「change做完」指示繼續原案：先完成本地原生分區、有限維護與既有Step 4整合，再執行真實普通／分區／group操作與瀏覽器驗收。此決定只調整待辦順序及拆分，不變更原需求、已完成五項描述或其touched provenance。原3.1的sys catalog與跨月readback實測由3.3負責；原4.1的真實scope write/readback/cleanup由新增4.4負責。2.3、3.3、4.4、4.3仍各自持有完整驗收，不以本地catalog fixture、mock或skip代替真實 SQL 實測；本輪可使用如實記錄的 ARM 模擬真實引擎，不提前解除能力開關。真實fixture尚未取得時繼續不依賴它的本地實作，真實驗收保持NOT RUN且不宣稱整案完成。

### 實際環境的功能驗收

**Supersedes**: add-sqlserver-circuit-recording / SQL Server 連線與能力範圍（固定 2022 驗證環境）、執行與驗收分開推進（必須原生 x86 的驗收環境條件）

使用者 2026-10-27 明確表示只需驗證 repo 可寫入 SQL Server，不用在意版本或 x86／ARM，並要求多 AGY 完成原案剩餘驗收。使用現有 SQL Server 2025 image，在 ARM 主機模擬運行的真實 SQL 引擎執行 2.3 → 3.3 → 4.4 → 4.3；版本、edition、compatibility level、image digest 與主機／容器架構均保存實際值。可執行的模擬環境功能實測計入功能證據，不以它推論原生平台效能、長期穩定性或其他版本皆支援。

保持所有精確值、來源到 outbox、schema ownership、分區／修復、write/readback/cleanup 與真瀏覽器條件。測試入口可接受如實的環境資料；不能偽造 native_amd64、版本或 compatibility，不用 mock／skip 代替實測，也不為了測試改正式產品 gate。若未改動產品的真實路徑仍阻擋，保留實際 HTTP／SQL／UI 證據與未驗範圍，不能勾完成。已完成九項描述與 touched provenance 不變；這是驗證環境調整，不是修碼、部署或新版本支援授權。

### 正式能力依實測條件開放

**Supersedes**: add-sqlserver-circuit-recording / SQL Server 連線與能力範圍（驗收期間所有 operation flags 固定關閉的實作狀態）

使用者 2026-10-27 在真實 production probe 證明連線、目錄與權限正常但四項 operation flags 固定為 false 後，明確授權「修正本案阻擋並續測」。允許必要的 SQL Server 能力判定、schema 接線與測試入口修正，完成原完整驗收；主代理負責調配與最後驗收，不提交或部署。

正式 capability 依實際 connected、verified metadata、已驗證 engine/version/compatibility 與各操作權限導出，不能無條件設定所有 flags。納入本輪真實 17.x／170 與原 16.x／160 的一致判定；未知版本、metadata／scope 不足、權限不足或 TLS 錯誤仍拒絕。長表、分區、試寫與修復各自符合所需條件，保留 scope/revision/readiness/token/ledger/ownership 與故障保護。所有原驗收條件保留，實測失敗不能以開關開啟宣稱完成。新增 5.1 承接已完成 1.1／1.3／4.1 的必要補償，已完成描述與 attribution 不變。

## Implementation Contract

- 追加來源契約：`exact_recording_scale` 為可選正值十進位文字；僅已確認整數 counter 可用。非法 scale／非整數 encoding／原始型別錯誤／stale revisions 不得產生可用 decimal；未設定保持舊行為。measurement 保存、epoch transition、group 型別驗證、frozen proof、pre-journal conversion 與回歸測試納入範圍，不擴充設備協議或另建計算引擎。驗證 TestExactRecordingScaleDefinition、TestExactRecordingScaleEpoch、TestExactRecordingScaleBoundary；原兩案真實 DB decimal raw/delta 仍須各別驗收。

- 資料形狀：sqlserver connector saved identity、database/schema scope、server version/edition、operation-specific capabilities、managed_long_v1、plain/native_monthly、recreate_missing；儲存版本及 layout digest 與確認 revision 一致。
- 失敗 code：adapter_unverified、exact_value_unrepresentable、schema_operation_partial、storage_owner_mismatch、partition_unavailable、partition_maintenance_blocked；沿用 revision conflict、outage、auth failure，UI 使用 en/zh-TW安全文字與修復入口。
- 驗收值：uint64 18446744073709551615、9007199254740992→9007199254740993 差 1、12000.0→12000.4 差 0.4、中文與大小寫不同來源、NULL、bool、UTC microsecond；超過 DECIMAL(38,18) 拒絕，不能 round。數字為契約案例，不是效能測量。
- 故障驗收：DDL/confirm 回應遺失、SQL ACK 遺失、並發同 key/digest conflict、deadlock、offline restart、permission denial、缺表／異質 owner、缺 partition 物件與跨月維護；readback、receipt、outbox、sys catalogs 必須對到相同 scope。
- 實測範圍：依「實際環境的功能驗收」使用可丟棄真實 SQL Server instance，記錄 edition/patch/fixture digest/SQL compatibility level 及如實的主機／容器架構；實際執行成功的模擬環境案例可計入功能證據，無法執行或 skip 維持 NOT RUN，不把功能證據當作原生平台效能驗收。實際瀏覽器涵蓋 plain/native_monthly、雙 profile、修復與狀態並留截圖。
- 範圍內為本資料庫 adapter及必要共同接點；不用此案修改無關設備、UI shell、legacy mapper 或前批 SQL方言。任務 Test* 為需新增的驗收，不表示已存在。
- 任務各估計 30–60 分鐘，開始到 focused verification 實際計時；超過 60 分鐘尚無可建置切片就拆細剩餘範圍，不以估計宣稱完成。

## Risks / Trade-offs

- [SQL Server 名稱涵蓋多種引擎／edition] → 首版固定測試基線與 operation gate，不承諾未測雲端產品。
- [decimal precision/scale 有限] → intake 前比對可表示性，維持精確值或明確拒絕。
- [分區 SPLIT 有鎖及資料搬移成本] → 可選、預建有限視窗、timeout及catch-all保護，不保證效能提升。
- [macOS ARM fixture 可用性] → 真實 instance 是完成門檻；缺環境保留 NOT RUN，不能只靠 mock解除gate。

## Migration Plan

Expand 新增必要 driver、connector/capabilities 與 plain layout，保留未驗證 gates；Migrate 以新 group、明確確認啟動 SQL Server，驗證 receipt/repair 後加入 native_monthly；Contract 只開啟已實測的操作，保留其他資料庫與 legacy 模式。每批可建置並通過 focused tests。回復關閉新 intake、排空 frozen outbox，再以新 effective boundary 回原配置，不刪物件或歷史。

## Open Questions

沒有需要使用者再選的需求；後續實作須準備真實可丟棄 SQL Server fixture 並確認 driver 與版本，環境缺失時相應驗收 NOT RUN。本輪依既定順序繼續實作。
