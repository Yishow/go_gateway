# Studio V2 SQL Server 迴路長表

SQL Server 的受管長表、耐久交付、原生月分區、群組試寫及 Studio V2 操作，已在先前 `add-sqlserver-circuit-recording` 案完成下述功能驗收。2026-10-27 在真實 SQL Server 執行普通表、分區、故障恢復與 plain／native 瀏覽器流程，均通過。`supports_managed_long`、`supports_native_partition`、`supports_group_test_write`、`supports_missing_repair` 依實際版本、scope、metadata 與各項權限分別判定，不再固定關閉；唯讀連線成功不代表可以寫入或建表。

共同操作、profile、來源語意與已審核切換見 [迴路記錄說明](studio-v2-circuit-recording.md)。先前功能案的規格與任務紀錄見 [add-sqlserver-circuit-recording](../../openspec/changes/archive/2026-10-27-add-sqlserver-circuit-recording/tasks.md)。本文件不更改 SQLite、PostgreSQL、TimescaleDB、MySQL／MariaDB 或既有 custom 群組的支援範圍。

## 本次已保存記錄管理與驗收界線

Basic／Advanced 的停止、移除、重建、三步確認與原操作復原，見 [已保存記錄操作說明](studio-v2-recording-lifecycle.md)。SQL Server successor 僅沿用受管長表、同引擎及同儲存模式的允許矩陣；不開放 wide、跨模式或跨引擎，也不搬移舊 SQL 歷史。能力與 schema proof 必須屬於此次精確目標，不能借用原目的地確認。

以下 2026-10-27 的 SQL 與瀏覽器 PASS 保留為先前功能證據，**不代表本次 `manage-recording-lifecycle-safely` 的 successor lifecycle 已通過實機驗收**。本次實際 SQL 寫入／讀回／cleanup 為 **NOT RUN**：依 2026-10-28 使用者「同意，若無法就跳過」的同機驗收例外，同機原生 private-prefix preflight 未建立官方文件支持的 workspace-only 完整 config/system root，因此安全停止，沒有下載、安裝或啟動 engine；這不是已證明 runtime 啟動失敗。

本次本地 adapter／proof／recovery 的實作與測試不得替代真實 SQL 驗收；產品 capability、TLS、schema 與 ownership 保護不變。正式 successor browser tests／截圖另屬 task 4.3，此次文件更新時尚未執行。完整例外與進度見 [lifecycle 設計](../../openspec/changes/manage-recording-lifecycle-safely/design.md)及[任務紀錄](../../openspec/changes/manage-recording-lifecycle-safely/tasks.md)。

## 已儲存的目的地與能力

- 使用已儲存的 `sqlserver` connector ID 與 identity revision，明確指定 database 與 schema。新 schema 草稿預設 `dbo`；讀取或重新載入不會覆寫既有值。
- 後端以 SQL authentication、加密且驗證憑證的連線探測；不因憑證錯誤關閉加密或信任未驗證憑證。瀏覽器不直接連資料庫。
- 能力查詢回傳 database/schema、實際 server version/edition、engine edition、compatibility level、metadata 與各項有效權限。連線成功與可檢視目錄是不同結果。
- 目錄不可見時回報 forbidden 或 unconfirmed，不能推論缺表、建立物件或使用其他 schema。修改連線身分、scope、憑證或 revision 後，舊確認不能授權新目的地。
- 候選版本為 SQL Server 2022（16.x／160）及 2025（17.x／170），版本與 compatibility level 必須配對。未知版本、scope 不符、metadata 不可見或必要權限不足仍拒絕。此次真實驗收使用 2025；16.x 的相容性回歸不等於此次跑過另一台 2022 實例。

## 同一組瞬時與累積 profile

| 型態 | 預設 | 範圍 | 資料語意 |
| --- | --- | --- | --- |
| 瞬時 | 15 秒 | 1–300 秒整數 | 依既有 gauge 平均、state 最新值與 counter_raw 語意 |
| 累積 | 10 分鐘 | 1–600 分鐘整數 | 保存原始讀值與既有 counter 計算結果，首筆無基準的 delta 為 NULL |

SQL Server 不另建計數器引擎。已確認整數來源若保存十進位文字 `exact_recording_scale=0.1`，`120000 → 120004` 在 journal 前精確轉成 `12000.0 → 12000.4`，沿用共同引擎得到差值 `0.4`。scale 修改建立新 epoch 並要求重新確認；已接受的 payload 與舊基準不重算。

`managed_long_v1` 使用固定欄位與 scope 衍生的 `gw_long_…` 表名，保留 immutable `gateway` 擁有權標記；有效值填入對應數值槽，其餘槽為 NULL。

| 資料 | SQL Server 表示 |
| --- | --- |
| UTC 窗口、取樣與接收時間 | `DATETIME2(6)`，保留微秒；不使用 rowversion 當時間 |
| signed integer | `BIGINT` |
| uint64、精確 decimal 與 delta | `DECIMAL(38,18)` |
| float | `FLOAT(53)` |
| bool | `BIT` |
| Unicode 文字與 JSON provenance | `NVARCHAR`，保留來源大小寫與合法 JSON |

超過精確可表示範圍時，在 intake 前以 `exact_value_unrepresentable` 拒絕，不交由 SQL 四捨五入。讀取大整數與 decimal 必須使用精確型別或文字，不能轉成 JavaScript Number。

## 準備、確認、套用與耐久交付

1. 唯讀預覽以實際目錄區分 missing、forbidden、unconfirmed 與 owned layout。
2. 明確確認綁定 workspace、group、connector revision、layout 與 digest；SQL 專用 managed executor 只執行後端產生且已確認的受管語句。通用任意 SQL executor 維持關閉。
3. DDL 回應遺失時，以同一 operation_id 查實際目錄及交易結果；未知或部分結果不能假設已 rollback，也不能用新 operation 掩蓋。
4. 套用與開始沿用共同 readiness/CAS 與生效邊界。已儲存、已連線、已檢視、已準備、已套用、已耐久接受、SQL 已提交、讀回與清理必須分別呈現。
5. 正式交付在同一交易寫 row 與 receipt；相同 effect key/digest 可查證重送，不同 digest 是衝突。ACK 遺失、deadlock、timeout 或離線保留既有重試與 backpressure；不以 MERGE 覆寫歷史。

已授權 `recreate_missing`、相同 scope/generation 且目錄證明缺少已確認 owned 表時，普通表恢復接點可重建結構並繼續 pending effects。它不保證恢復遺失的歷史資料，且保留既有 receipt。外來、不相容或無法證明擁有權的同名物件不得改造。原生分區缺 table/function/scheme 的真實修復驗收已通過。

原生恢復的本地接點涵蓋缺表、只剩 function、以及 table/function/scheme 全缺。只有經確認的獨佔物件或已驗證缺失物件可使用；執行 DDL 的交易內會重新核對 scope、權限與引用關係。重建沿用同一 generation/event CAS，保留初始 proof 與已授權的當前 boundaries，並更新 SQL Server 分區狀態；不重設 MySQL 狀態，不重建 receipt，也不修改外來或共用物件。本地測試不能代替真實 SQL Server 的交易及目錄驗收。

## 普通表與可選原生月分區

- **plain** 為預設；能力不足不會把使用者已儲存的 native 選擇改成 plain。
- **native_monthly** 必須明選且有實測操作能力。使用本案獨佔的 `DATETIME2(6) RANGE RIGHT` partition function 與 scheme，映射現有 `PRIMARY`；唯一鍵含 `bucket_start`，索引對齊。
- 初始預覽包含 UTC 當月與後續三個完整月，以及過去／未來 catch-all。receipt 為普通表。
- 同一既有 worker 只新增有限未來 boundaries，透過 scoped coordination/ledger 執行 `NEXT USED`／`SPLIT`；長期停機不逐月無界補齊。不新增另一個 scheduler，不建立 database、file、filegroup，不 DROP/TRUNCATE/MERGE 歷史或加入 retention。
- 已確認的初始 layout/digest 不會被維護改寫；當前擴充目錄須由同一 scope/generation 的授權 ledger 查證。未記錄或外來的分區漂移不能冒充初始 proof。
- 維護錯誤以 `partition_maintenance_blocked` 保留狀態。已驗證的 catch-all 仍安全時，既有交付接點可保留原凍結效果；這不代表可忽略目錄、row/receipt 權限或 operation gate。
- 分區建立權限不足時，明確 native 選擇回 `partition_unavailable` 與選 plain 的修復動作，不悄悄 fallback 或改 DBA 共用物件。

## 明確試寫與清理

Step 4 試寫沿用「預覽 → 明確確認 token/operation_id → 寫入 → 讀回 → 清理」。能力與 scope 通過檢查才允許明確試寫，畫面依真實結果分別呈現，不顯示固定成功結果。

接點只接受同一已確認 workspace/group/connector/revision 的受管長表。每個來源產生有界測試列，先保存本次 record IDs、sample provenance、effect keys 與 payload digests，再進行寫入。讀回比對所有 canonical 欄位、NULL、UTC、精確型別及 receipt digest。

清理交易只查本次精確 ID 集合，先驗證 scope、擁有權、sample identity 與 receipt digest，再以相同條件刪除；同 ID 的外來資料或衝突 receipt 會中止整批清理。不可只依表的 `gateway` 標記刪列。

| 讀回 | 清理 | 應呈現的結果 |
| --- | --- | --- |
| verified | cleaned | 試寫與清理各自通過 |
| verified | failed | 保留讀回成功與清理失敗，提供清理修復，不顯示無條件完成 |
| verified | unknown | 保留讀回成功，清理結果待查證 |
| partial／unknown | 任一狀態 | 分別呈現，不補成成功 |

完成操作重送回傳保留結果，不重寫測試列。legacy `recording-plans/test-write` 保留原 preview token／operation_id 安全轉接；裸請求、錯誤 action、scope/revision 改變或 stale token 不會授權寫入，也不回退為固定 501。

## 目的地切換與回復

儲存新連線或切換畫面不會自動改寫已儲存群組。明確選擇新目的地後，仍需共同 reviewed/CAS 與新生效邊界；受管表名由後端依新 scope 產生，不能沿用別的 scope 的確認。

已接受的 outbox 保存原 connector ID/revision、database/schema/table、payload 與 digest，依原目的地排送。修改連線或重試不能將舊效果改送新目標，也不重新計算平均或 counter。回復使用新的前向邊界，保留兩邊歷史與各自 pending effects，不刪除表或歷史。

## 2026-10-27 先前功能驗證與 fixture 門檻

單元測試使用 scripted SQL/catalog 與可丟棄 SQLite 設定庫，與真實 SQL Server 驗收分列。此次真實實例為 SQL Server 2025 `17.0.5035.5`、Enterprise Developer Edition、engine edition 3、compatibility level 170；Docker 主機為 ARM，執行 amd64 image。TLS 驗證專用測試 CA，沒有停用加密或憑證驗證。

依使用者確認，功能驗收使用現有可丟棄實例，不以版本或原生 x86 作為執行門檻。保留原本全部功能及錯誤案例：

- **2.3 PASS**：來源→profile→journal/outbox→row/receipt/readback，精確值、counter、故障、重啟及缺表。
- **3.3 PASS**：真實 sys catalog、分區 placement、並發維護、缺失分區物件修復與外來 ownership 拒絕。
- **4.4 PASS**：plain／native_monthly 相同已確認 scope 的真實 group write/readback/cleanup；清理失敗仍保留已驗證讀回、精確 row／receipt 與重送結果，解除故障後只清理該操作資料。缺表在預覽階段明確拒絕。
- **4.3 PASS**：plain／native_monthly 各自獨立工作區的實際瀏覽器流程；雙 profile、reload、schema 確認／修復、明確試寫／清理、正式記錄與 SQL／receipt 讀回，並保存 1024×720、1600×1000 截圖。兩輪各執行對應 mode，另一個 mode 由獨立 fixture 的另一輪驗證。

image digest 為 `sha256:245f6b1c8903edae52cb20fa31419cefcde1a49938bb8c3447099011fc871ca3`。正式 probe、普通表、分區與瀏覽器均使用真實 SQL Server。原始命令、退出碼及讀回／清理日誌保留在 `sqlserver-acceptance-20261027` 本輪證據中；完整結果見 change 的任務紀錄。測試只使用專用容器、loopback 及模擬來源。

### 真實驗收 fixture 的入口

Go 驗收以 `GW_SQLSERVER_TEST_FIXTURE` 指向 JSON manifest 的絕對路徑；它不是 DSN。manifest 不放帳密，`connection_config_file` 指向既有後端連線設定的絕對路徑，檔案權限不得寬於 `0600`。設定沿用 `user`，密碼逐字保留；不接受改用 `username` 或跳過 TLS 憑證驗證。

manifest 必須包含 `schema_version=1`、`disposable=true`、如實的 `host_os`（linux／windows／darwin）、`host_arch`（amd64／arm64）與 `native_amd64`。宣告原生 amd64 時架構必須為 amd64；ARM 模擬環境記為 arm64／false，不冒充原生。另需精確的 `database`／`schema`、`engine_edition`（2、3 或 4），以及 16.x／160 或 17.x／170 的版本配對。`server_version_prefix` 可使用更完整的實際 patch 前綴，不支援萬用字元。`server_edition` 可指定完整預期值；容器另記錄 `image_digest`，其他 fixture 以 `provenance` 記錄來源。實際 probe 仍須核對 scope、版本、edition、相容性與權限，manifest 宣告不能代替目錄證據。

| 入口 | 驗收目標 |
| --- | --- |
| `TestSQLServerCircuitPlainIntegration` | 真實普通表來源、雙 profile、journal/outbox、精確讀回及 counter |
| `TestSQLServerCircuitPartitionIntegration` | 原生分區目錄、placement、維護與缺失物件恢復 |
| `TestSQLServerGroupProbeLifecycleIntegration` | 已確認 scope 的試寫、讀回、限定清理與清理失敗 |
| `frontend/tests/e2e/sqlserver-circuit-recording.spec.ts` | 搭配 fixture 的瀏覽器流程與 SQL 證據 |

manifest 缺失時，入口以 `NOT_RUN` 跳過；operation gate 關閉時也不執行 schema DDL。編譯、列出測試與跳過測試不代表驗收通過。任務依 2.3 → 3.3 → 4.4 → 4.3 驗收完成；使用獨立可丟棄資料庫時，可平行收集互不影響的瀏覽器與 Go 證據，並保留各案讀回與清理結果。

瀏覽器 manifest 另需 `browser_storage_mode`（`plain` 或 `native_monthly`）、loopback 的 `backend`／`frontend`、既有 `config_db`、`workspace_id`、兩台 simulator 的 `devices` 與其 `mappings`。每種 mode 使用獨立且未開始記錄的 fixture，不重設或重用另一輪的工作區。讀回使用 `sqlcmd` 與可驗證的測試憑證信任鏈；本輪使用只連專用容器的暫存 wrapper，沒有安裝主機工具或關閉 TLS 驗證。

普通表另有 `TestSQLServerCircuitRawIntegration`、`TestSQLServerCircuitPlainCommitAckLostRetriesOneEffect`、`TestSQLServerCircuitSchemaFaultIntegration` 及 `TestSQLServerCircuitOfflineRestartIntegration`，分別承接精確型別、提交 ACK 遺失、部分 DDL 與離線重啟。分區真實驗收直接使用正式 probe 的能力結果，不複製後強制開啟 operation flags。
