# MySQL／MariaDB review 修復證據（2026-10-27）

本輪僅修復 `add-mysql-mariadb-circuit-recording` 的三項已確認問題。未 commit、push、部署或變更正式資料；下列執行結果來自本輪 macOS、自有 loopback 模擬來源與兩家可丟棄 TLS fixture。

## 實際結果與限制

- [summary.json](summary.json)：執行命令、版本、image digest、套件、實際 pass／skip 數及環境限制；[內容指紋](checked-content-sha256.json) 記錄本輪受影響套件／前端內容。
- [SQL 案例結果](real-sql-case-results.json)：兩家 plain／native、partial DDL、ACK 遺失、missing repair、maintenance、group test-write 與來源契約的已執行案例。完整 Go JSON 原始紀錄在本輪私有暫存附件，摘要保存其 SHA-256。skip 不等於 PASS。
- [MySQL 真撤權](mysql-permissions-real.json)／[MariaDB 真撤權](mariadb-permissions-real.json)：DELETE、INDEX、ALTER 各別撤銷後的能力閘門及恢復結果。其他 active-role／partial-revoke 情境是受控 probe 測試，未宣稱有真角色驗收。
- [MySQL 瀏覽器復原](mysql-browser-recovery.json)／[MariaDB 瀏覽器復原](mariadb-browser-recovery.json)：後端真實完成試寫後截斷回應；capability 重新查詢後仍 GET 同一操作，只有一次 confirm POST，readback／cleanup 分開核對。
- [MySQL SQL／outbox／receipt](mysql-runtime-readback.json)／[MariaDB SQL／outbox／receipt](mariadb-runtime-readback.json)：正常 server 入口的 applied revision、凍結 digest 與實際 SQL 對照。MySQL 本輪瞬時選平均來源，累積為精確 counter；不將瞬時結果聲稱為 scale=0.1 的 counter 原值。原本兩模式全組合與精確來源驗收保留於操作說明，本輪真 SQL 矩陣重新驗證。
- 瀏覽器錄影保留在本輪本機暫存附件，未列為可文字審查的 Git diff。Windows、真 PLC、LAN、正式 DB 均 NOT_RUN；未擴大到其他資料庫實體驗收。

## 三項回歸證據

| 修復 | 修復前確實失敗 | 修復後結果 |
| --- | --- | --- |
| 權限不明卻開放操作 | [baseline RED](isolated-baseline-red.txt)：`git archive HEAD` 的未修改 product，同一測試實際 assertion failure | [同一測試 GREEN](isolated-final-green.txt)、兩家真撤權、完整套件 |
| DDL 超過 lease 被第二 executor 接手 | [service RED](mysql-service-red.txt)：實際 executor 次數為 2 | [service GREEN](mysql-lease-focused-final.txt)、[race](mysql-lease-focused-race-final.txt)：跨十分鐘仍同 owner／單 executor |
| 試寫未知結果被背景更新清掉 | [UI RED](ui-red.txt)：背景更新後原查詢入口消失 | [UI GREEN](ui-final-focused.txt)、兩家真瀏覽器：保留原 operation ID，GET 查回 written_verified／cleaned |

## 畫面實測

| fixture | 1024×720：回應遺失後仍可查原操作 | 1600×1000：同一操作讀回與清理完成 |
| --- | --- | --- |
| MySQL plain | [待確認](mysql-unknown-1024x720.png) | [查回結果](mysql-recovered-1600x1000.png) |
| MariaDB native_monthly | [待確認](mariadb-unknown-1024x720.png) | [查回結果](mariadb-recovered-1600x1000.png) |

主代理逐張檢視四張截圖。MySQL 的右側工作區摘要保留最初 CA 未就緒時的既有 connector 狀態快取；試寫與正常 SQL receipt 分別已實際驗證，未把右側摘要當成 SQL 結果。程序專用 CA 信任已驗證，無 skip-verify／明碼 fallback，未修改系統信任庫。

## 規格對照與複審

17／17 任務已完成（原14項加本輪三項修復紀錄）；7 requirements、15 scenarios 均有實作及測試對照。spec 未定義獨立 Example 區塊，下表的契約數值由既有真 SQL／來源測試核對。

| Requirement（依 delta 順序） | 實作 | 對應 scenario／契約值 | 本輪結果 |
| --- | --- | --- | --- |
| Family-specific capabilities | dbtarget/mysql_recording_capability.go、mysql_recording_permissions.go | Separate family、forbidden metadata：TestMySQLFamilyCapability、TestMariaDBUnverifiedGate、TestMySQLManagedInspectionForbiddenIsNotMissing；兩家真撤權 | passed-current |
| Exact MySQL long values | dbtarget/mysql_long_row.go、mysql_group_insert.go、recordingplan/schema_group_mysql.go | Exact identities／unrepresentable+baseline：TestMySQLLongExactValuesReadback、TestMySQLCircuitPlainDecimalIntegration；uint64 上限、2^53 相鄰值、0.4、NULL、Unicode、UTC | passed-current |
| Per-statement DDL | recordingplan/schema_apply_mysql.go、mysql_ddl_heartbeat.go | Crash／foreign collision：TestMySQLManagedDDLRealRecovery、TestMySQLPartialDDLRecovery、TestMySQLForeignObjectRejected；新增十分鐘租約 regression | passed-current |
| Transactional effects／repair | dbtarget/mysql_group_insert.go、recordingplan/storage_restore_mysql.go、groupdelivery/sender.go | Lost ACK／missing vs foreign：TestMySQLLongCommitAckLost、TestMySQLSenderFrozenPayloadRestartAfterCommitAckLost、TestMySQLLongMissingRepair、TestMySQLLongRepairDenied | passed-current |
| Monthly partition maintenance | recordingplan/storage_partition_maintenance.go、dbtarget/mysql_managed_partitions.go | Rollover／year downtime／permission lost：TestMySQLPartitionMaintenanceIntegration、TestMySQLPartitionPermissionLostContinuesIntegration；保留 MAXVALUE 與歷史 | passed-current |
| Independent acceptance／workflow | cmd/test_ui wiring、Step4 shared panels、grouptestwrite | Basic partial／incomplete family：TestMySQLSchemaAPIRetainsOperationOnPartialConfirmation、TestMySQLGroupProbeLifecycle、mysql-mariadb-recording.test.tsx；兩家獨立矩陣、真 browser response loss | passed-current |
| Exact decimal sources | measurement/exact_scale.go、runtime/group_boundary_exact_scale.go、grouppipeline | Scale 0.1／new epoch／unsupported or absent：TestExactRecordingScaleDefinition、TestExactRecordingScaleEpoch、TestExactRecordingScaleBoundary、TestMySQLCircuitPlainDecimalIntegration；120000→120004、12000.0→12000.4、delta 0.4 | passed-current |

程式複審依 correctness、efficiency、simplification/reuse、convention 四面向：先前完整 scope 的未變內容沿用已完成檢視，本輪重新檢視所有修復、新增測試、必要呼叫端與翻譯；未發現待修的 Critical／Warning／Suggestion。新增修復紀錄使用 explicit_files 後補歸屬；修復前 baseline 及實際失敗另由本輪證據保存，沒有改寫原14項的 provenance。檔案歸屬不代表 dirty baseline 的每個 hunk 都由本案產生；原有 35 個 preexisting_dirty 限制仍保留，未 stage 或改動其他 WIP。
