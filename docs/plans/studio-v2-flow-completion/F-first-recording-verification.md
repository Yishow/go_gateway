# 空白目的地正式記錄驗收

2026-10-04；change：`validate-studio-v2-first-recording`；起點：`f8882bb6264302027b9590e6ca1510dcd05dac83`。

## 範圍與目前狀態

本案驗證 A–E 的實際組合，不新增產品功能。fresh 基本路徑必須由正常 embedded `/studio/v2` 完成所有 setup mutation、schema 明確確認及 start，API／獨立 SQL 只觀察。SQLite 起點以 stat／read-only URI 查證；PostgreSQL 僅部署可丟棄空 schema，不預建採集表或插入讀值。

fresh UI／三桶／七型別的單、雙設備及同群組多 entity 結果已有 final-v5 binary 的 PASS（見下方專節）；六次計時與恢復矩陣也已完成（見各專節與「最終核對」）；這只表示本案驗收矩陣已執行，不代表現場、真人或 release 驗收。真 UI 已重現 SQLite connector Save probe 在 schema 確認前建空檔；已完成最小 read-only probe 修正，Root 原始 HEAD overlay RED 與 focused/race GREEN 均有實際紀錄。畫面又重現 current SQL receipt 已確認卻仍顯示首桶等待；一行條件修正已先紅後綠，最新前端全套通過。managed 表診斷 payload 遺漏必要 metadata 已有 HEAD overlay 回歸與 SQLite／PostgreSQL production HTTP 整合證據，全套 Go gates 在最終 v6 工作樹亦通過；同表真 UI neighbor 驗收已完成（見「final-v6 binary 的結果」與「最終核對」）。映射與目的地 autosave 的 harness 亦補精確保存 barrier。歷史 `studio-v2-write-groups/evidence-f` 不覆寫；本次結果存於本案 `evidence-f/`，每個 run 有獨立 owned namespace。

| 已交付案 | 精確提交 |
| --- | --- |
| A | `eb30a21f0e7aff7345ea04292bc8d286fe9c66db` |
| B | `c26c8a52a1712d922015ba2d3822263c76a5dd7d` |
| C | `b560926c88ce4b477f53ec73e20e5d721964ba81` |
| D | `38c023e403efefc7601dbf8628d2fd8e8627b86f` |
| E | `f8882bb6264302027b9590e6ca1510dcd05dac83` |

正常入口的 `cmd/test_ui/service_wiring.go` 使用同一 canonical write-group service、managed/read-only inspector、production group store／pipeline／sample sink，以及 legacy ownership suppression。`target_writer.go` 的 legacy output 不把 `ErrOutputOwnedByWriteGroup` 當成成功 SQL 證據；F 的基本結果另以真正 group receipt 與獨立 SQL 核對。

逐一讀取 A–E 的實際 archive delta，共 63 個 scenarios；F 再加 11 個。本案沿用已提交的實作與既有獨立驗證，不重做產品功能；新增組合的 current witnesses 與最後 gates 已完成，見「final-v6 binary 的結果」與「最終核對」。

| 已交付規格 | Requirements／Scenarios | 原驗證對照 |
| --- | --- | --- |
| A lifecycle／durable intake | 2／11；1 example | [生命週期](A-lifecycle-verification.md)：草稿、cold start、ACK kill、frozen recovery、cutoff、reenable |
| B delivery deadline | 1／7 | [記帳期限](B-deadline-verification.md)：慢成功／失敗、取消後 commit、容量、worker、poison |
| C entity row | 2／5 | [Entity 列](C-entity-verification.md)：共欄／碰撞、獨立成列、缺值、結構錯誤復原 |
| D managed storage | 5／22 | [受管理儲存](D-managed-verification.md)：唯讀 preview、確認與 ownership、原時間／身分／精確型別 |
| E setup | 4／18 | [四步設定](E-setup-verification.md)：基本設定、scoped start、原 operation 恢復、Share-only、真正 receipt |

## 規格與必驗矩陣

兩個 requirements／11 個 scenarios，沒有 Examples；以下全部保留，沒有以時間或難度排除。

| Scenario／契約 | 本次狀態 |
| --- | --- |
| SQLite stat／read-only 不建立目的地；fresh file/schema | 11 個前置負測試 PASS；正常早期 UI Save／preview 均維持缺檔，confirmation 才建檔／表 |
| 空白 SQLite／PostgreSQL 經明確 UI confirmation 可持續記錄三桶 | 單設備兩 DB 三桶 PASS；PG 兩次完整 run 首列 143.891／141.611 秒，SQLite single1520 126.153 秒；final-v5 binary 的 PG A/B 同群組 matrix 與 SQLite 雙設備、同群組 A/B 各自 PASS，見下方「final-v5 binary 的雙設備與多 entity 結果」 |
| 七型別及 uint64 9007199254740993／9223372036854775808／18446744073709551615 | 兩 DB 單設備完整 run 九個精確值 PASS；final-v5 binary 的 PG matrix（A/B 各九個 tag）、SQLite 雙設備與 SQLite 同群組 A/B 均 PASS，SQL 驗證 failures 為空 |
| 單／雙設備同位址，identity／UTC／quality／revision 真 SQL 匹配 | 單設備 PG 三桶獨立比對通過；SQLite 雙設備先前 run 因 harness 字串 expected 索引錯誤失敗，該結果與 revalidation 保留；final-v5 binary 的乾淨重跑 `dual-sqlite-20261004-final-rerun` PASS，兩 Basic 群組各三桶，cleanup 通過 |
| 同一群組 A/B entity 同／不同業務欄位 | PASS；PG `pg-matrix-1791086682274`（2 Basic + 1 advanced 同群組 A/B，6 列）與 SQLite `f3-sqlite-multi-20261004-rerun2`（一個 A/B 群組、共用值欄位加各自欄位）各三桶、SQL 驗證 failures 為空；沒有以兩群組結果代替 |
| test-write readback／owned cleanup／neighbor 不變 | PASS（兩 DB）；`written_verified`／`cleaned`，既有 runtime 列逐位元組不變，額外列須有 current durable outbox／receipt，owned test row 清理後不存在；各 run 的 schema／namespace、process、browser cleanup 均通過 |
| 草稿不 Apply，跨兩桶及 restart 維持舊會員 | PASS-current；27 個 production lifecycle baseline tests 包含兩 DB |
| offline cold start new ACK；ACK 後 Disable/Delete/supersede 排空 | PASS-current；兩 DB lifecycle／kill historical baseline |
| slow sender success/failure；cutoff race | PASS-current；8 個新兩 DB cases；Root tagged race 通過，frozen payload/revision、SQL/receipt、checkpoint/journal 斷言 |
| 延遲補送保留 acquisition time；schema/start replay、stale、partial | PASS；final-v5 binary 的 `fresh-recovery-{sqlite,postgres}-…-recovery-20261004-v5` 兩 DB 均 PASS（schema／start lost response、延遲補送、main UI），原 operation 與 frozen payload／原採集時間一致，cleanup 通過；main、partial、error 在 390／768／1440 的截圖都已產生；我只實際檢視 SQLite 768 寬的 partial 與 error 圖，其餘寬度與 PostgreSQL 的圖僅確認存在、未逐張檢視；所見 partial 顯示「Recording started only partly／Device … not activated／Retry this request」，error 顯示「The result is not final. Check the same operation…」且 Apply 被擋 |
| custom／Share-only、quality／capacity／poison／worker／unknown／ownership | 既有 focused regression 與 final-v4 repository gates PASS；完整 command／skip／legacy PG failure 見 focused-regressions-root-20261004.json |
| SQLite／PostgreSQL 各三次 fresh 60 秒預設首列 ≤300 秒 | PASS（自動化計時）；各三次、每次獨立 namespace，最慢 SQLite 141.7 秒、PostgreSQL 117.0 秒（final-v6），見下方「首筆 SQL 計時」；真人操作觀察仍 NOT RUN |

## 計時與品質判準

使用同一 harness monotonic clock 的 `t_open`、`t_start`、`t_closed`、`t_sql`，各記 setup、完整 UTC 桶等待、交付及端到端耗時；poll interval／觀察延遲另列，SQL acquisition／bucket UTC 時間不拿來算跨主機耗時。第一個 run 保留產品 60 秒預設；故障案例如使用既有 10 秒區間，須從 UI 明確設定並獨立標示。

每個 success run 至少查三個正式連續桶，核對原 device/point/tag/group、source/mapping/applied revisions、精確業務值、record/effect/receipt identity 與 quality。SQL committed、真正 readback、owned cleanup 和 child exit 分別記錄；任何 cleanup／child failure 都不發布 PASS。

自動化計時不等於人因驗收。一般操作流程的輸入、卡點、等待與修復行為另記；沒有真人 usability/sign-off 時明列 NOT RUN。

## 執行證據

darwin/arm64、Go 1.27.1（module 1.25.5）、golangci-lint 2.14；loopback simulator、可丟棄 SQLite／自有 PostgreSQL 16.14 容器。沒有真 PLC、LAN、SCADA、正式 DB 或部署操作。

| 命令／範圍 | 結果 |
| --- | --- |
| 四份現有 `node --test`：HTTP deadline、mixed identity、fixture cleanup、CAS witness | PASS-current；24 tests；`gw-F-existing-harness-unit.log` |
| 正常 build 的 fixture hooks／容量環境變數保護 | PASS-current；兩個相鄰 tests；`gw-F-normal-no-fixture.log`；正式時鐘與容量不受 fixture 設定取代 |
| fresh 目的地前置保護 | PASS-current；11 tests；owned aliases／失敗觀察 cleanup 負測試先紅後綠，mutation 確認意外開啟 SQLite 建檔會失敗，unknown／非 boolean absence proof 亦失敗；還原後全綠 |
| 空白 PostgreSQL owned schema 與 cleanup | PASS-current；Root 獨立 SELECT 查證 0 個 recording objects、沒有預建採集表，再刪除該唯一 schema 並查證不存在 |
| 早期單設備 SQLite UI→首列 | PASS；[early-root-1791073252766.json](evidence-f/early-root-1791073252766.json)；兩個 D/E 點位、default 60 秒，首列 128.6 秒，setup 17.1／完整桶等待 110.4／交付觀察 1.0 秒；不代替完整矩陣 |
| production recovery／fault matrix | focused recordingplan/groupdelivery/grouppipeline/API/cmd 與 live PG cmd/dbtarget 已 PASS；新 wired deadline/cutoff 與 UI 恢復已在 final-v6 完成，見「final-v6 binary 的結果」 |
| repository Go／frontend 最低 gates | 最後正常 v4 snapshot：Go test／vet／lint PASS，49 個測試 packages、另 7 個無測試檔；1107 module/source/config inputs 前後一致，go list -test 的 1082 個 Go inputs 均在其中；[實際 gate 紀錄](evidence-f/default-go-gates-final-v4-root-20261004.json)。兩 DB slow/cutoff 加 metadata HTTP 的 10 個 tagged race 子案例 PASS（16.398 秒）；[整合與 race](evidence-f/production-regressions-final-v4-root-20261004.json)。最新 UI lint／190 files、1211 tests／tsc+Vite build PASS，見 [先紅後綠與全套結果](evidence-f/basic-pending-regression-and-ui-gates-root-20261004.json) |
| 正常 F embedded binary | 修正後 build PASS-current；592 production source identities 前後一致；無 fixture tags 的 `go build ./cmd/test_ui`；99 個 dist files 與 embedded static 逐一 hash 相同 |

本次初始 source/test gate identity manifest 涵蓋 1701 個 tracked Go／frontend files，存於自有本機暫存紀錄；最後交付須核對實際差異，不能只憑 E 的綠燈沿用。

本次修正前正常 binary SHA-256 為 `13e02bed932e56a08ef992f31c3106e698527ccd97bdf695587d0cf49229ac45`；592 個 production Go／go.mod／go.sum 的 identity digest 為 `9d5ccd55b90bef4ba51627ef58c9264e4c266343a775e46892781a963537c575`。上述產物僅保留早期失敗證據；最終 UI／SQL 驗收須指定修正後正常產物；fault fixtures 與 integration build tags 另列，不能替代此基本路徑。

修正後正常 binary SHA-256 為 `055dcbd0c536f5fe96bc66c15910c71d4ed835183654005bf3b716affeb8805a`；592 個 production source 的 identity digest 為 `552dd703fa62198cd99455c9cf737d95855d3fa37ab42f8a37b5fdcab2522076`。目前只有 `service_probe.go` 的 production bytes 改變；99 個 embedded static 仍與 dist 逐一一致。見 [build](evidence-f/normal-build-sqlite-probe-v2-root-20261004.json) 與 [先紅後綠回歸](evidence-f/sqlite-probe-regression-root-20261004.json)。

最新等待提示修正的正常 UI-v3 binary 為 `b91ef7e6d0fcff6d4c71b12d68acdef4646fec421578049054fd5e02942f4039`，Go recording source 保持 probe-v2 identity，99 個真正 dist 資產與 embedded 檔逐一一致；tracked embed placeholder 另保留。此 build 先於診斷 metadata 修正，最後 matrix 必須指定後續正常 build。見 [UI-v3 build](evidence-f/normal-build-ui-evidence-v3-root-20261004.json)。

PG 首次完整 main 三寬度：390 的 Basic bounds 為 100–374、768 為 104–748、1440 為 112–1132；Basic controls 無水平越界，advanced summary 的 Enter 開合／恢復通過。390 的既有長 connector label 讓 document 寬度為 465，原始觀察保留；沿用 E 的 Basic 控制項與側欄／connector 收合驗收邊界，沒有新增整頁 responsive 重設條件。該 run 的 3 PNG 主代理已檢視；最新 UI 另補 app 內捲動後的底部控制項與 receipt 圖。

最後正常 final-v4 binary SHA-256 為 `c81c8706b17d6cb0e630a713eb7e8e25a1c4937100d3ffe3c50d61113a882a0a`；592 production Go/mod/sum inputs digest 為 `290b0b8f0309567abd67dc242b8a455fda482f4b58c936c66c0f743e193ec355`，99 個 dist 資產與 embedded static 逐一相同，tracked placeholder 保留 HEAD bytes。見 [最後 build](evidence-f/normal-build-final-v4-root-20261004.json)。後續只有驗收 harness／文件變動，不重用不相符的產品 gate。

managed metadata 修補沿用既有 row layout、preview digest、readback 與 owned cleanup，沒有解除缺 entity owner 的 fail-closed 422，也沒有新增共欄診斷模式。回歸的 unit／HTTP fixture neighbor 與真 UI matrix 的 runtime neighbor 分別記錄；見 [原版回歸與修補結果](evidence-f/managed-metadata-regression-root-20261004.json)。

獨立 verifier 增加 SQL provenance 與 frozen bucket 的 sample ID/status/quality/原採集 UTC 比對，保留奈秒精度及等值的小數尾零。原 verifier 對錯誤 sample ID 的負測試失敗，修正後五 tests 通過；再對兩個 PG single 與 SQLite single1520 真實留存觀察核對。這是純資料重驗，不宣稱新的 UI、SQL SELECT 或 clock；見 [verifier 回歸](evidence-f/frozen-provenance-verifier-regression-root-20261004.json)、[實際資料重驗](evidence-f/current-frozen-provenance-verification-root-20261004.json)。

## 首筆 SQL 計時

最終 binary 為 final-v6（`/private/tmp/gw-F-test-ui-production-final-v6`，SHA-256 `97d64d8844373037ecdd0da3a502c18988565e814a2c6bd743e1fd326407ccc9`，99 個內嵌資產與 `frontend/dist` 逐一相同）。單設備、九個 tag、群組 `row_policy.interval_seconds = 60`（三個桶的 `bucket_start` 相隔 1 分鐘；connector 目的地的 `write_interval_seconds = 5` 是另一個欄位，不是記錄區間）。每次由 harness 從首次打開空白 setup 起算；每個 run 獨立 namespace，結束後獨立查證 namespace／schema 不存在、child 與 browser 正常結束；證據內的 `binary` 欄位均為 v6。

| DB | Run | setup 秒 | 完整桶等待 秒 | 交付觀察 秒 | 端到端 秒 |
| --- | --- | --- | --- | --- | --- |
| SQLite | `v6b-timing-sqlite-1` | 40.9 | 100.8 | 0.0（毫秒級） | 141.7 |
| SQLite | `v6b-timing-sqlite-2` | 41.2 | 76.6 | 0.0（毫秒級） | 117.8 |
| SQLite | `v6b-timing-sqlite-3` | 40.8 | 76.6 | 0.0（毫秒級） | 117.4 |
| PostgreSQL | `pg-root-1791098946101` | 41.6 | 72.0 | 1.2 | 114.7 |
| PostgreSQL | `pg-root-1791099184115` | 41.6 | 74.1 | 1.2 | 116.9 |
| PostgreSQL | `pg-root-1791099424097` | 41.6 | 74.3 | 1.2 | 117.0 |

最慢值：SQLite 141.7 秒、PostgreSQL 117.0 秒，皆低於 300 秒目標。每個 run 也驗證三個連續正式桶，沒有預建採集表或用直接 API 取代 UI；SQLite 三次的 connector 狀態為 `ready`。

限制：
- 計時是 harness 的腳本速度；setup 約 41 秒是腳本的固定等待，不代表一般操作者的輸入、理解或卡點，真人觀察仍 NOT RUN。
- 交付觀察包含最多一個 poll interval（SQLite 2 秒、PostgreSQL 1 秒）與 SQL 命令時間；跨主機時鐘不納入。
- 環境是本機 darwin/arm64、loopback simulator，不外推到真設備或現場。

先前以 final-v5 binary（SHA `413668d4…`）各三次的結果保留作歷史參考：SQLite 117.5／135.8／117.4 秒、PostgreSQL 116.7／117.0／118.1 秒；另有名為 `v6-timing-sqlite-1..3`、`pg-root-1791096606587／…844491／…084543` 的 run，實際仍是 v5 binary（當時三支 harness 寫死 binary 路徑，忽略 `F_GATEWAY_BINARY`，證據的 `binary` 欄位如實記錄），不作為 v6 證據。更早最初一批三次（SQLite 兩次、PostgreSQL 一次）的 JSON／PNG 在讀完數字後從 `evidence-f/` 消失，原因不明，不列入也不當證據。

## 恢復與互動矩陣（歷史：final-v5）

| Run | 結果 | 備註 |
| --- | --- | --- |
| `fresh-recovery-sqlite-sqlite-recovery-20261004-v5` | PASS | binary `413668d4…`；schema／start lost response、延遲補送、main UI；process、namespace cleanup 通過 |
| `fresh-recovery-postgres-postgres-recovery-20261004-v5` | PASS | 同上；owned schema 已 DROP 並查證不存在 |

既有保護在目前工作樹重跑（production source 自 final-v5 gate 後沒有變動）：

- `POSTGRES_DSN=<owned loopback> go test -tags f_write_group_fixture -race -p 1 ./cmd/test_ui -run '^TestProductionF(SlowSenderSuccessAndFailureDurable|ParallelCutoffKeepsDurableIdentity)$|^TestProductionTestWriteKeepsManagedMetadataAndNeighborRows$' -count=1`：三個 test PASS（13.4 秒），涵蓋 SQLite／PostgreSQL 的慢 sender 成功／失敗、並行 cutoff、managed metadata／neighbor／owned cleanup。
- `env -u POSTGRES_DSN go test ./internal/datalink/{recordingplan,groupdelivery,grouppipeline,grouptestwrite,dbtarget} ./internal/api/... -count=1`：全部 ok。品質、容量、poison、worker、unknown commit、custom／Share-only 與 ownership 的既有 focused tests 在這些套件內，沒有刪減。
- 舊版 `quality.mjs`／`capacity.mjs`／`worker-overlap.mjs` 等 harness 會寫入歷史 `studio-v2-write-groups/evidence-f`，為不覆蓋歷史證據，本次沒有重跑；其保護以上述 Go focused tests 與先前 final-v5 紀錄為準。

## 歷史：final-v5 binary 的雙設備與多 entity 結果

本節為歷史紀錄（最終證據以 final-v6 為準），沿用同一個正常 binary `/private/tmp/gw-F-test-ui-production-final-v5`（SHA-256 `413668d4db31b34a9bed09a2303ab56746caff1526af038155012307d991f9f1`），PG 為重建的可丟棄 `gw-wg-pg-test`（PostgreSQL 16.14，127.0.0.1:55432）。本段新增的只有 harness 修正，沒有改產品 source，也沒有改 service 不外露 `SchemaOperation.Detail` 的契約。

| Run | 結果 | 備註 |
| --- | --- | --- |
| `pg-matrix-1791086682274` | PASS | A/B 同群組、test-write；setup 84.6／完整桶等待 112.4／交付觀察 1.2 秒；端到端 198.2 秒；schema 與 owned work 已清除，容器內無殘留 `gw_*` namespace |
| `f3-sqlite-multi-20261004-rerun2` | PASS | 一個 A/B 群組、explicit test-write；setup 79.5／等待 75.0／交付觀察 1.0 秒；端到端 155.6 秒；namespace 已移除 |
| `dual-sqlite-20261004-final-rerun` | PASS | 兩 Basic 群組、同位址各三桶；端到端 150.7 秒；namespace 已移除 |

本輪修掉的兩個 harness 缺陷，均先有失敗測試再修：

1. PG matrix 的 `test-write-draft` 失敗為 `"undefined" is not valid JSON`：harness 從 API 讀 `ledger.detail`，但該欄位在產品是 `json:"-"`。改為從本地 `managed_schema_operations` 讀 detail，並新增檢查：API 若外露 `detail` 即失敗（`readLedgerDetail`）。先前失敗證據 `pg-matrix-final-v5-20261004-1005.json` 保留。
2. SQLite 多 entity 在 `explicit-test-write` 連續兩次失敗（`f3-sqlite-multi-20261004-v5-final`、`-rerun1`，record_id 不同）：`readSQLiteDelivery` 的 outbox SELECT 少了 `entity_key`，使 `assertRowsPreserved` 的 entity 比對永遠為 `undefined`。補欄位後 PASS。

harness 單元測試最新為 60 tests 全過。上述雙設備與多 entity 每次 run 仍只是單次觀察；各三次預設 60 秒的正式計時見上方「首筆 SQL 計時」。

## final-v6 binary 的結果

產品 source 在審查後改了一處（見「審查後修正」），所以 binary 重建為 final-v6，核心 run 在 v6 上重跑。PG 為重建的可丟棄 `gw-wg-pg-test`（PostgreSQL 16.14）。

| Run | 結果 | 備註 |
| --- | --- | --- |
| `pg-matrix-1791097324520` | PASS | A/B 同群組、test-write；setup 84.1／等待 91.8／交付 1.2 秒；端到端 177.1 秒；驗證 failures 為空，cleanup 通過 |
| `v6-sqlite-multi` | PASS | 一個 A/B 群組、explicit test-write；端到端 175.3 秒；connector `ready` |
| `v6c-dual-sqlite` | PASS | 兩 Basic 群組同位址；兩裝置端到端 162.7／160.7 秒；connector `ready` |
| `fresh-recovery-postgres-postgres-recovery-20261004-v6` | PASS | 恢復矩陣，cleanup 通過 |
| `fresh-recovery-sqlite-…-v6s1／v6s2／v6s3` | PASS ×3 | 恢復矩陣；gateway 均正常結束、未被 SIGKILL |
| `fresh-recovery-sqlite-…-v6r2／v6r3` | PASS ×2 | 加入停止餘裕之前的重跑 |
| `fresh-recovery-sqlite-…-v6` | FAIL | 四個情境結果皆有，但 cleanup 失敗：`gateway did not exit cleanly`（SIGTERM 逾時後被 SIGKILL）；失敗證據與 gateway 日誌 `fresh-recovery-sqlite-v6-failed-gateway-shutdown.log.txt` 保留 |

SQLite 恢復矩陣在 v6 共 6 次：1 次失敗、5 次 PASS。失敗那次的日誌顯示收到關閉信號後 5 秒出現「關閉 datalink runtime 失敗: context deadline exceeded」，這個 5 秒是 `cmd/test_ui/main.go` 的固定關閉期限；當時 harness 的 SIGTERM 逾時也剛好是 5 秒，且 harness 的目的地寫鎖在 gateway 停止之後才釋放，所以我推斷是 harness 逾時與產品期限賽跑。這是推論，不是直接證明；adapter 的 gateway 停止逾時改為 8 秒後連跑 3 次均 PASS，被 SIGKILL 仍判失敗。

雙設備在 v6 的前兩次（`v6-dual-sqlite`、`v6b-dual-sqlite`）因計時判準失敗：harness 先觀察完第一組的三個桶才輪到第二組，第二組的 `t_sql` 被推遲（第二裝置 301.5 秒但只輪詢 1 次）。已改為同時觀察兩組，判準仍是 300 秒；失敗時也會保存實測計時。失敗的兩份證據保留。

### 審查後修正

`/spectra-review` 提出的項目與處理：

1. **缺檔 SQLite connector 停在 `unreachable`**（`service_probe.go`）：讀 readiness 程式確認它會產生 Warning「database connector is currently unreachable」，且 DDL 建檔後不會更新；這是先前「Save probe 改唯讀」的副作用。修正：檔案不存在但上層目錄存在且可寫時回報 `ready`（重用 `managedSQLitePath`／`validateManagedMissingParent`），檔案仍不提前建立；上層目錄不存在仍是 `unreachable`，且不建立任何目錄。HEAD 原本的 handler 測試即斷言此情境為 `ready`，上一輪修補才改成 `unreachable`，這次改回並保留不建檔的斷言。
2. **harness cleanup 例外會吞掉證據**：新增 `removeOwnedWork`（不丟錯、記錄 `owned_work_cleanup_error`、失敗即不 PASS），套用於 `fresh-early-root`、`fresh-postgres-root`、`fresh-postgres-matrix-root`；`fresh-cleanup.test.mjs` 以唯讀上層目錄驗證。
3. `groupDeviceID` 與另外兩處重複的 Suggestion，以及 `waitFCommitHold` 的時間斷言偏弱（只能在排程停頓時失敗）的 Suggestion：未處理，列為後續建議。
4. 第二次 review 的 Warning（`fresh-ui.mjs` 預設 binary 路徑寫死 v5、simulator 沒有覆蓋變數）：預設值改讀 `F_GATEWAY_BINARY`／`F_SIMULATOR_BINARY`，未設定時指向 final-v6；`fresh-cleanup.test.mjs` 先紅後綠。仍未改為每次從 source 建置，也未統一 `fresh-early-root` 專用的 `F_EARLY_ROOT_BINARY`。

重跑過程中另外修正的 harness 問題：`fresh-single-sqlite`／`fresh-dual-sqlite`／`fresh-postgres-root` 改為尊重 `F_GATEWAY_BINARY`；雙設備併行觀察；雙設備失敗時保存計時；恢復 adapter 的停止逾時。

## 最終核對（工作樹目前狀態）

起點 `f8882bb6264302027b9590e6ca1510dcd05dac83`，工作樹含未提交的產品修補與驗收 harness；最終 binary 為 final-v6（SHA-256 `97d64d88…`）。darwin/arm64、Node 24.15、loopback simulator、可丟棄 SQLite 與 PostgreSQL 16.14（容器 `gw-wg-pg-test`，已重建）。

repository 全套最低檢查（產品修正後的目前工作樹，`go test` 使用 `-count=1`，未帶 `POSTGRES_DSN`）：

| 檢查 | 結果 |
| --- | --- |
| `go test -p 1 ./... -count=1` | PASS；49 個套件 ok、7 個無測試檔 |
| `go vet ./...`、`golangci-lint run ./...` | PASS；lint 0 issues |
| 前端 `npm run lint`、`npx vitest run`、`npm run build` | PASS（190 個測試檔、1213 tests）；此結果取自產品修正前，之後沒有任何前端檔案變動 |
| `go test -tags f_write_group_fixture -race -p 1 ./cmd/test_ui -count=1`（帶 owned PostgreSQL） | PASS；131 個 PASS、0 FAIL、2 SKIP（`TestProductionLifecycleCrashChild` 是無環境變數時的 helper 子行程；`sqlite/quoted-distinct` 在測試中明確只對 PostgreSQL 執行） |
| harness `node --test scripts/tests/f_device_to_sql/*.test.mjs` | PASS；61 tests |

本次 spec 11 個 scenarios 對照（證據以 v6 為準）：

| Scenario | 結果與證據 |
| --- | --- |
| Mixed multi-device happy path | PASS：PG `pg-matrix-1791097324520`、SQLite `v6c-dual-sqlite` |
| Supported acquisition types and decimal boundary | PASS：兩 DB 九個精確 tag（含三個 uint64 邊界）；decimal 僅 codec／SQL 層，設備採集 NOT RUN（unsupported） |
| Operation cleanup | PASS：兩 DB test-write `written_verified`／`cleaned`，neighbor 不變 |
| Empty destination becomes usable through UI | PASS：兩 DB、各三個以上連續桶 |
| SQLite emptiness checks do not create the destination | PASS：destination-preflight 負測試、UI 前後 stat／read-only 查證，以及缺上層目錄不建立目錄的新測試 |
| One group contains multiple entities | PASS：PG 與 SQLite 各自單一 A/B 群組，沒有以兩群組替代 |
| Controlled first-use timing | PASS（自動化，v6）：各三次，最慢 SQLite 141.7、PostgreSQL 117.0 秒 |
| Draft does not stop existing recording | PASS：tagged `cmd/test_ui` 套件（含 PostgreSQL） |
| Offline cold start and historical journal | PASS：同上，含 killed historical journal、unprovable recovery 保留 journal |
| Slow result settlement and concurrent cutoff | PASS：tagged race 兩 DB |
| Start recovery and unchanged acquisition time | PASS：恢復矩陣兩 DB（v6）與 API replay focused tests；SQLite 有 1 次因 harness 逾時賽跑失敗的紀錄，見上節 |

NOT RUN／限制：Windows、原生 Linux／ARM 設備、LAN、真 PLC、SCADA、正式 DB、部署、長時間 soak、真人 usability／sign-off、decimal 設備採集。此外：
- 本輪重跑未使用 `quality.mjs`／`capacity.mjs`／`worker-overlap.mjs` 等會寫入歷史 evidence 的舊 harness，相關保護以 Go focused tests 與先前紀錄為準。
- 前面的重跑證實有些 harness 曾寫死 binary 路徑；v5 與 v6 的 run 已用證據內的 `binary` 欄位區分。
- 一批最初的計時證據檔與 41 個已追蹤檔案曾在工作樹消失，原因不明；已追蹤檔案由 `git checkout` 還原，未追蹤的計時證據以補跑取代。
- 提交的證據範圍：`evidence-f/` 全部 JSON（含失敗與歷史 run）、失敗 gateway 日誌的 `.log.txt` 副本，以及 SQLite `v6s3` 與 PostgreSQL `v6` 兩次最終恢復矩陣各一組截圖；其餘開發中途的截圖只留在本機、沒有提交，文件中引用到的 run id 以其 JSON 為準。
- 本案完成只表示驗收矩陣已執行，不表示 release。

## 尚未執行與限制

Windows、原生 Linux／ARM 設備、LAN、真 PLC、SCADA、正式 DB、部署、長時間 soak、真人 usability/sign-off 均 NOT RUN。decimal 不新增設備採集，僅沿用 existing codec／SQL round-trip 證據；其設備採集明列 unsupported／NOT RUN。以上不以此次本機 pass 替代。

本次 live PostgreSQL focused API 仍重現 A 已記錄的舊 recording-plan `TestStudioV2RecordingRoutes_PostgresApplyCreatesAndVerifiesTheManagedSchema`：2 passed／1 failed，結果為 `unknown/verification_unavailable`，不是 confirmed success。此 endpoint 不等於本案 canonical group schema route；未改寫該範圍外功能，也不宣稱帶 POSTGRES_DSN 的全 repo suite 通過。本案 fresh PostgreSQL 路徑仍須獨立完成，不能用這個區分跳過其驗收。
