## 1. 前置及 SQL Server 普通長表

- [x] 1.1 確認 MySQL/MariaDB 5.1／5.2 已完成「已確認十進位 scale 的最小來源擴充」及 Confirmed exact decimal counter sources，保留其來源 proof／epoch 契約；落實「SQL Server 連線與能力範圍」及 Verified SQL Server connection and destination scope：確認共同 add-circuit-recording-profiles 已實作驗證、MySQL/MariaDB 批次已交接後再開始；核對既有外部 SQL adapter並重用，必要時加入 Microsoft go-mssqldb，完成 sqlserver saved connector、database/schema、唯讀 probe、TLS驗證與version/edition/operation gates。先寫 TestSQLServerConnectorScope、TestSQLServerCapabilityGate、TestSQLServerTLSRejected，確認失敗後實作，再跑 dbtarget/recordingplan focused tests；前置不足保持 blocked。
- [x] 1.2 [after: 1.1] 落實「精確長表與交易內 receipt」中的 Exact SQL Server long-table profile values：新增 managed_long_v1 的 datetime2(6)/UTC、DECIMAL(38,18)、BIT、Unicode JSON/text、case-exact identity及正確 quoting/parameters，intake前拒絕不可表示值。先寫 TestSQLServerLongExactValues、TestSQLServerDecimalRejected，覆蓋 uint64最大值、大整數差1、decimal差0.4、NULL／Unicode／UTC，不使用 rowversion 當時間，維持共同雙 profile 的界線與預設。
- [x] 1.3 [after: 1.2] 落實「確認操作與受管缺表恢復」的 Confirmed SQL Server schema operations and repair：preview/confirm/ledger綁定 workspace/group/connector/layout/digest，使用真實metadata區分 missing/forbidden/unconfirmed，DDL回應遺失按同 operation_id 查實際catalog與交易結果。先寫 TestSQLServerSchemaConfirmation、TestSQLServerSchemaAckLost、TestSQLServerForeignOwnerDenied，再實作，驗證只有已確認完整owned布局可通過readiness，不 CREATE DATABASE或改外部物件。

## 2. 耐久寫入與恢復

- [x] 2.1 [after: 1.3] 完成 Transactional SQL Server durable effects：接 production OpenDestination／InsertGroupRow／Sender.Deliver，在同 transaction寫 row＋receipt，並發效果鍵與digest安全比對，ACK遺失可查證、deadlock/outage保留資料重試，不以MERGE覆蓋歷史。先寫 TestSQLServerReceiptConcurrent、TestSQLServerCommitAckLost、TestSQLServerDeadlockRetry，再實作並回歸 groupdelivery／service wiring，確認只一列及counter不重算。
- [x] 2.2 [after: 2.1] 完成 Confirmed SQL Server schema operations and repair 的缺表部分：只有已授權 recreate_missing、有效scope/generation及目錄缺表才能協調恢復，保留committed receipts與歷史遺失說明；外來／共享partition物件不能改造。先寫 TestSQLServerMissingRepair、TestSQLServerRepairScopeDenied，驗證同scope並發恢復一次、缺權限／revision變更／異質同名表被擋、pending rows安全續送。

## 3. 原生月分區與維護的本地實作

- [x] 3.1 [after: 2.2] 落實「SQL Server 原生月分區」及 Optional SQL Server native monthly partitioning 的建置：新增明選plain/native_monthly、專屬RANGE RIGHT function/scheme、PRIMARY映射及對齊索引，確認預覽顯示UTC當月＋後續3個月／catch-all，unique含bucket_start。先寫 TestSQLServerPartitionDDL、TestSQLServerPartitionUnavailable，以本地交易／catalog fixture驗證preview、確認proof與runtime接點；原sys catalog及跨月readback實測移交3.3；不建立filegroup/檔案或修改DBA共享function。
- [x] 3.2 [after: 3.1] 完成 Optional SQL Server native monthly partitioning 的有限維護：以既有worker、scoped coordination/ledger新增NEXT USED/SPLIT boundaries，維持有限3個月視窗，停機很久不無界補齊，失敗保留history並呈現maintenance狀態。先寫 TestSQLServerPartitionRollover、TestSQLServerPartitionDowntime、TestSQLServerPartitionMaintenanceDenied，驗證catch-all安全時續寫、timeout有界、沒有DROP/TRUNCATE/MERGE或自動retention。

## 4. 基本操作與交付

- [x] 4.1 [after: 3.2] 落實「分批開放與操作一致性」及 Capability-backed SQL Server workflow and acceptance：既有Step 4接SQL Server capabilities與同一雙profile，顯示database/schema、plain/native_monthly、schema修復與saved/connected/prepared/durable/committed/readback/cleanup各結果；實作group test-write的已確認scope寫入、readback與清理接點（真實instance操作驗收由4.4負責），en/zh-TW安全錯誤可修復。先寫 sqlserver-recording.test.tsx、TestSQLServerGroupProbeLifecycle，確認失敗再實作，cleanup失敗不能改寫readback成功或顯示整體成功。
- [x] 4.2 [after: 4.1] 驗證SQL Server前向目的地切換與回復沿共同CAS/boundary，舊outbox排原目標，SQLite/PostgreSQL/Timescale/MySQL/MariaDB與custom模式保持可用；新增 TestSQLServerDestinationCutover、TestSQLServerAdapterCoexistence及sqlserver-compatibility.test.tsx。只開實測版本/操作gate，文件review核對操作說明、scope與未驗證版本拒絕；保留legacy recording-plans/test-write既有preview token／operation_id安全轉接，不回退501，驗證裸請求與stale token仍拒絕。

## 5. 真實資料庫與瀏覽器驗收

- [x] 2.3 [after: 5.1] 依「實際環境的功能驗收」使用現有 SQL Server 2025／ARM 模擬的可丟棄真實 instance，新增 TestSQLServerCircuitPlainIntegration 驗證模擬來源→profiles→journal/outbox→row/receipt/readback，涵蓋精確值、baseline/gap/reset、schema partial/ACK遺失、離線重啟及缺表。記錄patch/edition/compatibility level/image digest與cleanup；記錄如實主機／容器架構，不以版本或原生 x86 作為功能驗證門檻；無法實際執行的案例保持 NOT RUN，不能以 mock／skip 勾完成。
- [x] 3.3 [after: 2.3] 新增 TestSQLServerCircuitPartitionIntegration，承接3.1的sys catalog及跨月readback驗證，在真實instance跑partition catalog/placement、boundary/late rows、並發維護、缺table/function/scheme修復與foreign ownership拒絕，對照plain結果相同及row/receipt不變；記錄每案readback與cleanup，驗證只恢復本案專屬且未被外部物件共用的布局，未實測不得開partition capability。
- [x] 4.4 [after: 3.3] 完成4.1移交的 Capability-backed SQL Server workflow and acceptance 真實group test-write驗收：同一已確認workspace/group/connector/revision及plain/native_monthly scope，實際write/readback/cleanup分別留證；執行 TestSQLServerGroupProbeLifecycle 的真實fixture入口，驗證cleanup失敗保留readback成功且不顯示整體成功。無真實instance維持NOT RUN，不以本地fake或disabled UI勾完成。
- [x] 4.3 [after: 4.4] 用真實fixture與模擬來源完成瀏覽器plain/native_monthly、雙profile/reload、schema確認/修復、錯誤與交付階段操作截圖，逐scenario對照SQL/catalog/receipt證據；依docs/testing.md執行受影響Go tests/vet/lint、前端lint/tests/build、git diff --check與make check-lines。平台/真實DB/UI分列PASS/FAIL/NOT RUN；缺fixture或skip不能勾完成，不自動commit、archive或部署。

執行與驗收分開推進（2026-10-26，使用者要求繼續完成change）：已完成五項描述與touched provenance不變。原待辦2.3／3.3／4.3保持真實驗收；原3.1的實際sys catalog及跨月readback明確移交3.3，原4.1的真實group操作移交4.4。這是執行順序與任務拆分，沒有刪除、替代或降低任何原驗收。全部真實fixture、瀏覽器及平台證據分列，缺環境維持NOT RUN，不能以mock或emulation解除原生x86驗收及operation gate。

原始13項加上本輪授權的補償任務5.1，共14項；完成狀態以上方勾選與 CLI 回報為準。每項原估計30–60分鐘，實際修復與驗收時間以執行紀錄為準。交易、並發／恢復及UI採red → green → refactor，每切片結束需可建置且通過focused tests。

Test*與前端測試名稱是本案驗證入口，非PASS聲明。Expand增加plain及gates，migrate加入已確認新groups/partition，contract僅解除實測操作gate；保留其他資料庫、既有preview token／operation_id安全轉接與frozen outbox。前批MySQL/MariaDB artifacts及touched evidence不改。

## 本機驗證紀錄（2026-10-26）

使用者要求透過 Herdr 分配多個 AGY 平行測試，主代理只調配及最後驗收、不擴權。本輪由 `sql-be`、`sql-fe`、`sql-fixture` 分別執行後端、前端與驗收前置檢查；主代理回讀原始輸出及核對工作樹。以下是本機驗證結果，不是整案或真實 SQL Server 驗收通過。

| 檢查 | 結果 | 範圍與限制 |
| --- | --- | --- |
| Go SQL Server 選測 | PASS；真實 DB 部分 NOT RUN | 232 個頂層測試通過、0 失敗、10 個真實 DB 測試跳過；`cmd/test_ui` 無符合選測條件的測試，不計為測試通過 |
| Go vet／golangci-lint | PASS | 下列 10 個後端 packages，exit code 均為 0 |
| SQL Server 前端選測 | PASS | 5 個檔案、47 個測試通過 |
| 相關前端回歸 | PASS | 5 個檔案、68 個測試通過；有 React `act(...)` 測試警告 |
| 前端 lint | PASS | `npm run lint`，0 errors、2 個 `no-console` warnings，位於既有 `device-recording-workflow.spec.ts` |
| 前端型別與建置 | PASS | `npm run build`（tsc + Vite）；有 chunk size 警告，未同步嵌入式 static |
| 差異與行數檢查 | PASS | `git diff --check` exit 0；`make check-lines` 檢查 439 個檔案通過，有行數警告及既有超長檔未增行；涵蓋共同 dirty tree，不代表全部檔案歸本案 |
| 瀏覽器測試入口 | 僅清單可載入 | Playwright `--list` 列出 plain／native_monthly 共 2 項；沒有實際瀏覽器驗收或截圖 |
| 真實資料庫／瀏覽器驗收 | NOT RUN | 本機 Darwin arm64，未設定 `GW_SQLSERVER_TEST_FIXTURE`，且 PATH 無 `sqlcmd`；沿用暫不驗真實 SQL Server 的限制，未建立環境或連線試寫 |

後端以清空 `GW_SQLSERVER_TEST_FIXTURE` 的環境執行 `go test -count=1 -v -run SQLServer`，並對相同範圍執行 `go vet`、`golangci-lint run`：`./internal/datalink/dbtarget`、`./internal/datalink/recordingplan`、`./internal/datalink/groupdelivery`、`./internal/datalink/grouppipeline`、`./internal/datalink/runtime`、`./internal/datalink/workspace`、`./internal/datalink/grouptestwrite`、`./internal/api/handlers`、`./internal/api`、`./cmd/test_ui`。沒有執行 `go test ./...`，不得宣稱後端全套通過。

前端同樣清空 fixture，從 `frontend/` 以 `npm test -- --run` 選取下列檔案，未執行前端全套：

- SQL Server：`tests/unit/workbench-v2/sqlserver-recording.test.tsx`、`sqlserver-compatibility.test.tsx`、`sqlserver-circuit-fixture.test.ts`（後兩者同目錄）、`tests/unit/services/sqlserverRecordingCapabilityContract.test.ts`、`tests/unit/utils/sqlserverErrors.test.ts`。
- 相關回歸：`tests/unit/workbench-v2/circuit-profile-actions-reduction.test.ts`、`locale-parity.test.ts`、`writeGroupLifecycle.test.tsx`（後兩者同目錄）、`tests/unit/services/studioV2WriteGroupTestWrite.test.ts`、`tests/unit/utils/safeJson.test.ts`。

證據基準為 `accfa1963db122ce483212a4b127ba39c8838941` 上的本輪未提交工作樹。原始證據包含 `gotest_sqlserver.log`、`govet.log`、`golangci-lint.log`、`sqlserver-unit-tests.log`、`regression-unit-tests.log`、`npm-run-lint.log`、`npm-run-build.log`、`env_check.log`、`playwright_list.log`、`check_lines.log`，保留於本輪暫存證據目錄。主代理在更新本紀錄前比對 4,294 個 tracked／untracked 檔案內容、Git index 與 HEAD，均未變動，並直接執行 `git diff --check` 通過。

完成狀態維持 **9/13**。2.3 → 3.3 → 4.4 → 4.3 的真實驗收仍須依序執行；4.3 的本機檢查部分已有上述證據，實際瀏覽器與 SQL/catalog/receipt 證據仍缺。四項 checkbox 不變，未解除 `supports_managed_long`、`supports_native_partition`、`supports_group_test_write`、`supports_missing_repair`，未執行修碼、提交、封存或部署。本次僅補登驗證紀錄，不變更需求、驗收條件或既有 touched provenance。

## 真實 SQL Server 儲存驗證（2026-10-27）

使用者澄清本輪只需驗證「repo 可存 SQL Server」，不以版本或 x86／ARM 作為此項儲存驗證的門檻。本輪使用現有 image，由兩個 Herdr AGY 分別準備可丟棄容器與執行 repo writer 測試，主代理檢視暫存測試及原始結果。這項範圍說明取代本輪前置檢查對 2022／原生 x86 的要求；不取消原有分區、故障恢復或完整瀏覽器驗收。

- 實際環境：SQL Server 2025 `17.0.5035.5`、Enterprise Developer Edition、engine edition 3、compatibility level 170；Docker daemon 為 linux/aarch64，執行 linux/amd64 image，架構如實記為模擬環境。
- Image：`mcr.microsoft.com/mssql/server:2025-latest`，digest `sha256:245f6b1c8903edae52cb20fa31419cefcde1a49938bb8c3447099011fc871ca3`。使用專用可丟棄資料庫 `gw_circuit_smoke`，只綁定 loopback，TLS 驗證臨時 CA，未跳過憑證驗證。
- 執行方式：`go test -overlay <暫存 overlay.json> ./internal/datalink/dbtarget -v -run TestSQLServerStorageSmoke`，以額外暫存測試直接呼叫未修改的 production `openExternalDBManager`、`InsertGroupRow`，catalog proof 取自真實 metadata。沒有覆寫既有 fixture loader、偽造版本／架構或更改產品能力開關。
- **PASS**：1 個真實資料庫測試，包含 2 筆不同型別資料；精確 decimal `12000.4` 與 delta `0.4`、uint64 最大值 `18446744073709551615`、Unicode、NULL 及 UTC 微秒讀回均符合斷言；row 與 receipt 寫入成功。
- **PASS**：相同 effect key／digest 重送回報 `AlreadyCommitted`，不同 digest 拒絕為 `ErrReceiptDigestMismatch`；清理刪除 2 筆測試 row 與 2 筆 receipt，讀回剩餘筆數均為 0，測試表移除亦無錯誤。
- 原始證據：本輪暫存目錄內 `runner/sqlserver_smoke_test.go`、`runner/overlay.json`、`runner/smoke_test.log`；結果為 `TestSQLServerStorageSmoke` PASS，實測測試時間 0.47 秒。日誌 SHA-256：`ce10b5ba4da9f29caef98f3a3631c4a0b882f4906e265252d48124be4f22db6a`。
- 主代理核對：執行前後 repo 檔案內容、Git index 與 HEAD 均未變動；本次只補登此驗證紀錄。未驗證 int64／bool 等其他型別，不由暫存測試註解推論已通過。
- 環境清理：已移除本輪建立的容器及臨時帳密、TLS 憑證與私鑰，保留使用者原有 SQL Server image；無秘密的測試程式與日誌留存於本輪暫存證據目錄。

結論：本輪「repo 原有低層 writer 能實際寫入並讀回 SQL Server」已驗證通過。此測試未經 Studio setup、高階 capability 開放或瀏覽器，也未驗證來源→journal/outbox 全流程、分區或故障恢復；因此不將它等同 2.3／3.3／4.4／4.3 全部完成，原 9/13 與已完成項目 provenance 保持不變。

## 完整驗收續跑（2026-10-27）

使用者指出原已要求多 AGY 測試，不能只以低層 writer smoke 收尾。依 design「實際環境的功能驗收」續跑四項完整驗收；前述 2025／ARM 儲存 smoke 是已有部分證據，並非 2.3／3.3／4.4／4.3 的替代。固定版本／原生架構的執行門檻已被本輪決定取代，其餘驗收與正式能力 gate 保留。測試失敗或產品路徑阻擋必須留下實際證據，不自動修碼或勾選。

## 6. 完整驗收阻擋修正

- [x] 5.1 [after: 4.2] 完成「正式能力依實測條件開放」：修正 production probe 在 connected／metadata／權限均有效時仍固定關閉四項能力，以及 schema／fixture 對本輪 17.x／170、如實 ARM 環境的必要接線阻擋。依實際版本／engine／scope／權限分別導出能力，未知版本、metadata／權限不足、TLS 錯誤仍拒絕，不放寬 token／revision／ownership。保留既有 16.x／160 契約。先以 TestSQLServerProductionCapabilityGateEvidence 的真 SQL FAIL 及 focused regression 取得未改產品的失敗證據，再修正並跑 capability／schema／fixture 的 focused tests 與真 probe，確認完整 2.3 能開始；前端以 sqlserver-circuit-fixture.test.ts 的真實環境案例先失敗再通過。此項是已完成 1.1／1.3／4.1 的補償，不改寫其完成描述及 touched provenance。不得直接全部設 true、不 commit／archive／部署。

- [x] 5.2 [after: 5.1] 修復本輪 verify／review 的五項既有需求缺口：SQL Server schema 重新確認後同步分區維護狀態；native capability 要求實際所需目錄可見權限；已讀回但清理失敗的試寫提供原操作限定清理修復；native fixture 能力不足在 DDL 前標示 NOT RUN；同步共同文件支援狀態。先補產品未修正時失敗的回歸，再完成最小修復及受影響測試、靜態檢查與必要真實 SQL／瀏覽器驗證。保留 scope／revision／token／CAS／ownership，不重寫試寫資料、不改其他問題。此項僅承接原需求的審查修正，不擴大需求；封存另依使用者本輪授權執行。

## 完整功能驗收結果（2026-10-27）

依使用者「授權修正本案阻擋並續測」，由 Herdr 三個 AGY 分別負責 Go 整合、SQL Server 專用試寫修正與瀏覽器；主代理負責分工、逐次 diff／原始結果覆核、Spectra 勾選與最後驗收。固定版本／原生架構門檻依 design 的 Supersedes 決策調整，其餘功能與失敗案例全部保留。

實際環境為 SQL Server 2025 `17.0.5035.5`、Enterprise Developer Edition、engine edition 3、compatibility 170；Docker ARM 主機執行 amd64 image，並使用驗證專用 CA 的 TLS。image digest 沿用上方實測紀錄。Go 與瀏覽器使用不同的可丟棄 database，plain／native 使用不同設定庫與工作區；不以模擬能力或 skip 代替對應 mode 的真實執行。

本輪必要產品修正限於能力條件、17.x／170 接線、同 scope 多群組各自已確認的 ledger 證據、SQL Server 原生分區來源入口，以及 SQL Server 試寫的 member reduction 與合法 profile 時間窗。gauge 才產生 mean，counter_raw／state 保留既有語意；未知或非法 profile 拒絕，不猜測時間窗。writer、scope／revision／token／generation／ownership 與 TLS 拒絕保護保留。

| 驗收 | 結果與原始證據 |
| --- | --- |
| 2.3 普通表完整流程 | PASS：`02_plain_integration.log`（36.17 秒）、`02_raw_integration.log`、`02_fault_acklost.log`、`02_circuit_offline_restart.log`、`02_missing_table.log`、`02_schema_fault.log`。涵蓋精確值、baseline／gap／reset、單次效果／receipt、payload 不重算、DDL 回應遺失、離線重啟與缺表拒絕。 |
| 3.3 原生分區 | PASS：`03_partition_integration.log`（9.16 秒）；真實 catalog／placement、並發維護、三種缺失物件恢復、初始 proof 不變、current boundaries 保留及外來擁有權拒絕。 |
| 4.4 群組試寫 | PASS：`04_groupprobe_integration.log`（25.55 秒）；plain／native 真實寫入、精確讀回、限定清理、清理失敗保留結果與資料、重送及缺表拒絕。 |
| 4.3 瀏覽器 | PASS：`browser-plain-success/e2e-plain.log`（1.7 分鐘）與 `browser-native-success/e2e-native.log`（1.8 分鐘）；各自完成完整測試，保留 SQL／receipt 證據與 1024／1600 截圖。每輪另一個 mode 因使用獨立 fixture 排除，其實際 PASS 由另一輪提供。 |
| Go 選測與靜態檢查 | PASS：10 個既有受影響 packages 的 SQLServer 選測；新 runtime boundary、試寫及 lostcommit 回歸另重驗。受影響產品 packages 與最後 API 檔案的 vet／golangci-lint 均通過，見 `local-*`、`final-*` 日誌。缺 fixture 的本機跳過結果未計入真實 SQL 驗收。 |
| 前端 | PASS：10 檔 78 tests、lint、TypeScript／Vite build。lint 保留既有 2 個 no-console warnings，build 保留 chunk size 提醒；未同步嵌入式 static。 |
| 規格測試對照 | 15 個 Scenario 各有具體測試對照，見 `scenario-coverage.md`；此為靜態覆蓋盤點，執行結果以上列實際日誌為準。 |
| 工作樹檢查 | `git diff --check` 與 `make check-lines` PASS；依使用者指示保留現有暫存，不 commit／archive／部署。 |

能力固定關閉、native 來源入口、混合 member 試寫與零長時間窗均保留修正前的真實失敗及修正後通過證據。測試 fixture 修正只處理真實版本／架構、凍結證據讀取、唯一測試 ID、SQL Server catalog 物件識別與操作資料清理，沒有刪除必要斷言或放寬正式 writer。

本節為本輪完整驗收結果；先前 9/13、僅 smoke 或 NOT RUN 的段落保留為歷史紀錄。任務完成狀態以目前 checkbox 與 CLI 的 14 項為準。

最後核對：CLI 回報 `all_done`、14/14，`spectra validate` 為 valid 且沒有 error／warning；`analyze` 沒有 Critical／Warning，保留 13 項既有的具體範例建議。工作樹除本案明列 27 個既有檔案與 3 個新增測試檔外，其餘 4,267 個既有檔案內容不變，HEAD 與使用者要求保留的暫存項目皆未改動，見 `final-scope-audit.json`。

本輪 backend／Vite／兩個 simulator 已停止；專用容器及臨時帳密、私鑰、設定庫已移除，保留使用者原有 image、manifest metadata、全部驗收日誌與截圖，見 `cleanup-summary.json`。未 commit、archive 或部署。


## 審查修正與封存前驗收（2026-10-27）

使用者明確授權修復前輪 verify／review 的五項問題並封存，另要求禁止擴大範圍。本輪由 Herdr 三個 AGY 分工，主代理核對完整修復差異及原始證據；AGY 額度用完後，主代理接手最後的專用瀏覽器驗收。只處理本節五項及必要的測試、產生文件與可丟棄 fixture 前置，未修復範圍外問題。

- R1：SQL Server schema 重新確認時，在同一設定庫交易以原 generation／digest／revision 查證及 CAS 同步分區維護狀態；同表保留已授權擴充邊界，重建世代重設至該次確認布局，外來或過期狀態拒絕並回滾。三項回歸均取得未修產品時失敗、修復後通過證據。
- R2：native capability 同時要求 database VIEW DEFINITION 與 ALTER ANY DATASPACE；schema-only 可見性仍保留 plain 能力。production projection 與 probe 回歸先失敗再通過，未放寬版本、scope 或其他權限。
- R3：既有 group test-write API 新增明確 action `cleanup_repair`，僅 SQL Server 受管長表可用。查證原 token／內容 digest／workspace／group／connector revisions 與清理 proof，原操作以 scope claim 防止重複執行，只清理原 record IDs／effect keys／digests；不重寫。未確認操作、異質 scope、未知 action 與不完整 proof 安全拒絕。既有無 action 的確認與 legacy 路徑維持重播行為；UI 分別呈現寫入、讀回與清理，失敗或待查證時提供限定修復入口，其他資料庫不顯示此入口。
- R4：native fixture 在外部 schema DDL 前檢查 native 能力，關閉時明確 NOT RUN；四種 plain／native 能力組合及實際 skip 分支有回歸。這些分支單元測試不冒充真實 DB 通過。
- V1：共同操作文件同步 SQL Server 已驗收環境與條件 gate。Swagger 由 `make gen-docs` 正式產生，僅同步本案 optional action 及已存在的 SQL Server current-readiness／partition 模型依賴，沒有新增產品行為。

本輪實際驗證結果：

| 檢查 | 結果與限制 |
| --- | --- |
| 受影響 Go modules／API 選測 | PASS：684 個頂層測試；39 個缺其他外部 fixture 的項目明列 NOT RUN，不計通過。未跑 `go test ./...`。 |
| Go vet／golangci-lint | PASS：dbtarget、recordingplan、grouptestwrite、api、handlers、cmd/test_ui 六個 packages；修正本輪新增檔案造成的 lint 問題，未清理其他舊檔。 |
| 前端選測／lint／型別與建置 | PASS：7 個測試檔、79 個測試；受影響檔案 ESLint 與 `npm run build` 通過。保留既有 chunk-size 警告，未同步嵌入式 static。 |
| 真實 SQL 清理修復 | PASS：`TestSQLServerGroupProbeLifecycleCleanupDenied`；解除 DELETE 故障後，以同一 token／operation 執行修復，INSERT 被故障守衛拒絕仍成功清理；原 row／receipt 歸零，其他 receipt 保留。最後產品版本再次實跑通過。 |
| 真實 SQL 原生分區 | PASS：`TestSQLServerCircuitPartitionIntegration`，含缺失物件修復及外來 ownership 拒絕。原測試會從並發返回值選到舊快照；僅改為在維護完成後讀取正式保存的 current state，保留精確邊界、歷史及 immutable proof 斷言，未更改產品維護流程。保留失敗與修正後通過日誌。 |
| 真實瀏覽器 | PASS：agent-browser 專用工作階段，一個 plain／instantaneous gauge 群組完成 UI 預覽、確認、真實讀回、DELETE 故障、修復按鈕及清理；同一 operation 保留 `written_verified`，清理由 failed 變 cleaned。SQL 原操作 1 row／1 receipt 均歸零，修復及正常確認重送期間啟用 INSERT 故障守衛，重送保留同一 operation 與保存時間。已保存 1024×720、1600×1000 各自失敗與修復後截圖。這是 R3 的聚焦驗收，未重跑前輪完整雙 profile／plain／native 瀏覽器矩陣。 |

本輪環境仍為 SQL Server 2025 `17.0.5035.5`、Enterprise Developer Edition、engine edition 3、compatibility 170，ARM 主機執行 amd64 image；TLS 驗證專用 CA。證據保留於 `sqlserver-review-repair` 的修復報告、原始測試日誌、source fingerprints 與四張截圖。未將此功能證據延伸為原生平台效能、長期穩定性或其他版本驗收。

前輪 14 項已完成工作與證據保留；本輪 5.2 僅承接同一需求的審查修正。13 項既有文件格式 Suggestion 未擴大處理。封存依本輪使用者明確授權執行，未授權提交、推送或部署；原暫存內容保留。
