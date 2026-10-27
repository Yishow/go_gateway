## 1. 前置確認與普通長表擴充

- [x] 1.1 落實「共用流程與分開的伺服器能力」及 Family-specific verified MySQL and MariaDB capabilities：先確認 add-circuit-recording-profiles 的共同 profiles/long schema/durable closure 已實作驗證，否則停留 blocked；在既有 mysql kind 增加後端 family/version 與逐操作 capability projection，連線、metadata、managed 與 partition 分開 gate。先寫 TestMySQLFamilyCapability、TestMariaDBUnverifiedGate，確認失敗後實作／整理並跑 recordingplan/dbtarget focused tests；定位 GetConnectorCapability 與 OpenDestination，不以現有 driver 或 UI 選項宣稱支援。
- [x] 1.2 [after: 1.1] 落實「固定長表與精確 MySQL 型別」及 Shared profiles with exact MySQL long values：沿用 managed_long_v1，新增 InnoDB、DATETIME(6)/UTC、DECIMAL(65,30)、case-exact IDs、typed NULL 與 JSON/digest codec，實作 intake 前精確值可表示性檢查及 database/schema scope。先寫 TestMySQLLongExactValues、TestMySQLLongScope，驗證 uint64最大值、大整數差1、decimal差0.4、Unicode／UTC、超界拒絕，保留前置 profiles 的界線與預設。
- [x] 1.3 [after: 1.2] 實現「逐句建表恢復與耐久交付」中的 Recoverable per-statement managed DDL：preview/confirm/operation ledger 對每步 DDL 保存進度，實際 catalog 查證 ownership/索引與同 operation_id 恢復；不能以整批 rollback 當結果。先新增 TestMySQLPartialDDLRecovery、TestMySQLDDLResponseLost、TestMySQLForeignObjectRejected，再實作，更新 capability_mysql_test.go 的新受管分支與未驗證拒絕案例，確認只有完整查證後 readiness 才通過。

## 2. 可靠交付與已確認缺表恢復

- [x] 2.1 [after: 1.3] 實現 Transactional MySQL effects and safe missing-table recovery 的交付部分：接入 production OpenDestination／InsertGroupRow／Sender.Deliver，InnoDB transaction 內寫 row＋receipt，frozen payload 重送以 digest 比對，identity conflict 不覆寫。先寫 TestMySQLLongReceiptReplay、TestMySQLLongCommitAckLost、TestMySQLLongConcurrentConflict，確認 SQL ACK 遺失及重啟只有一列、不重算 baseline，再跑 groupdelivery／production wiring focused 回歸。
- [x] 2.2 [after: 2.1] 完成 Transactional MySQL effects and safe missing-table recovery 的 schema repair：僅在已確認 recreate_missing/generation/scope、catalog 缺表與相同 connector revision 下依 ledger 協調恢復；保留 receipt，不改造外來同名表。先寫 TestMySQLLongMissingRepair、TestMySQLLongRepairDenied，驗證並發恢復一次、異質 owner／權限／revision 失敗可見、DDL 不混入 row transaction，且歷史不宣稱復原。
- [x] 2.3 [after: 5.2] 在可丟棄 MySQL 8.4 fixture 驗證普通長表垂直切片：模擬輸入→雙 profile→journal/outbox→SQL/readback，包含精確值、baseline、reset/gap、失敗重送、DDL partial及缺表。新增 TestMySQLCircuitPlainIntegration 入口並實際選中所有案例，記錄 patch/engine/sql_mode/image digest、receipt與cleanup；缺環境標 NOT RUN，不勾完成。

## 3. 可選月分區與 MariaDB 各別驗收

- [x] 3.1 [after: 2.3] 落實「原生月分區與有限維護」及 Optional native monthly partition maintenance 的建表部分：增加 plain/native_monthly 明選、UTC RANGE COLUMNS、含 bucket_start 的 unique keys、獨立普通 receipt 表與預覽當月＋後續3個月／catch-all；能力不足明確拒絕。先寫 TestMySQLNativeMonthlyDDL、TestMySQLPartitionCapabilityGate，再以真實 MySQL catalog／跨月 readback 驗證，不建 foreign keys／retention jobs。
- [x] 3.2 [after: 3.1] 完成 Optional native monthly partition maintenance 的有限維護：使用既有 worker 加入 scoped lease/ledger 的 future horizon 擴充與 MAXVALUE REORGANIZE，停機很久只補有限未來窗，maintenance failed 可見而安全 catch-all 可續寫。先寫 TestMySQLPartitionRollover、TestMySQLPartitionLongDowntime、TestMySQLPartitionPermissionLost，驗證 history/receipt不變、timeout有界、沒有 DROP/TRUNCATE/MERGE 或逐月無界補齊。
- [x] 3.3 [after: 3.2] 在獨立 MariaDB 11.4 fixture 實跑相同 plain/native_monthly、partial DDL、receipt/ACK遺失、missing repair 及 maintenance 矩陣，僅修正本案 family-specific 差異；新增 TestMariaDBCircuitIntegration，記錄 patch/engine/sql_mode/digest及每案readback，確認 MySQL 通過不能替代 MariaDB證據，任一家未跑則整案仍未完成。

## 4. 操作流程與批次交付

- [x] 4.1 [after: 3.3] 落實「分批開放與既有模式相容」及 Verified workflow and independent database acceptance：Step 4 使用共同 profiles，顯示實際 family/capabilities、plain/native_monthly、schema partial/恢復與各個交付狀態，en/zh-TW 同步；group test-write 實際寫入/readback/cleanup且保留各結果。先寫 mysql-mariadb-recording.test.tsx、TestMySQLGroupProbeLifecycle，確認失敗後實作，保留 legacy recording-plans/test-write 既有 preview token／operation_id 安全轉接契約，不回退 501；驗證裸請求及 stale token 仍拒絕。
- [x] 4.2 [after: 4.1] 驗證新 mysql 受管能力只對已測版本/操作開放，既有 custom、SQLite/PostgreSQL 與 Modbus-only 仍可用，目的地編輯採共同 CAS boundary 且舊 outbox 排原目標。以新增 TestMySQLLongLegacyCoexistence、TestMySQLDestinationCutover 及前端 mysql-mariadb-compatibility.test.tsx 驗證，保留舊 capability 安全拒絕情境；更新相關操作說明，content review 對照 specs 與 scope。
- [x] 4.3 [after: 4.2] 完成兩 family 的實際瀏覽器錄製／截圖及整合證據對照，依 docs/testing.md 執行受影響 Go packages tests/vet/lint、前端 lint/tests/build、git diff --check、make check-lines；每個 spec scenario 對應真實結果，資料庫／UI／平台分開記 PASS/FAIL/NOT RUN。只在兩家 plain/partition/repair/group test-write 均驗證後才解除對應 gates，不自動 commit、archive 或部署。

原12項保留原編號與已完成證據；新增下列2項經使用者授權的來源前置，合計14項，以各 checkbox 為完成依據。每項估計30–60分鐘，由實際起迄及 focused verification 驗證；超過60分鐘尚無可驗收切片則在可建置狀態拆細剩餘工作。演算法、交易／恢復與UI採 red → green → refactor，每切片須可建置及 focused tests 通過。

Test* 與前端測試檔名為本案要求新增的驗證入口，並非現有或已執行聲明。[after] 序列也涵蓋 shared schema/capability/UI 檔案衝突；先 expand plain、再 migrate family/partition、最後只解除已驗證操作 gates。不得同時改寫另案 SQL Server；本案完成後才開始該批。

## 5. 追加授權的精確 decimal 來源前置

- [x] 5.1 [after: 2.2] 落實「已確認十進位 scale 的最小來源擴充」及 Confirmed exact decimal counter sources 的保存契約：MeasurementDefinition 以可選 exact_recording_scale 十進位文字保存正值有限倍率，只允許已確認整數 counter 語意，空值相容；scale 編輯提升 definition revision 並開新 series epoch，非法值不得保存，不重解 mixed-layout float scale。先寫 TestExactRecordingScaleDefinition、TestExactRecordingScaleEpoch 與 repository roundtrip 測試，red → green → refactor 後跑 measurement/schema focused tests；本切片可建置，不先開新 runtime intake。
- [x] 5.2 [after: 5.1] 完成「已確認十進位 scale 的最小來源擴充」及 Confirmed exact decimal counter sources 的凍結及 intake：group backend 推導原始整數／ExactDecimal 結果型別並凍結 scale、measurement/binding revisions、epoch，sampleFrom 在 journal 前驗證原始型別且只縮放一次，沿用既有 counter、representability、CAS boundary；接受後 payload 不重算，scale 切換不跨 epoch 比較。先寫 TestExactRecordingScaleBoundary 與 stale proof／wrong source／unscaled integer及float回歸測試，red → green → refactor；驗證 120000→120004＋0.1 的瞬時 raw 及累積 raw/delta 0.4 經 journal/outbox，不取代2.3兩家真實 SQL readback。

追加決策來源：使用者明確選擇「授權最小精確 decimal 來源擴充，保留兩案原驗收」。5.1／5.2 接在完成的2.2之後、2.3完成之前；附加於檔尾以保留原 CLI checkbox 索引及已完成 task_desc/provenance。原2.3已開始的獨立整數測試可保留，但須完成兩個前置及原全部驗收才勾選。

銜接決策（2026-10-25 使用者核准）：4.1 原 legacy 501 保留條款由 design「分批開放與既有模式相容」的 Supersedes 決策取代，沿用新 HEAD 已完成的共同安全轉接流程；不更改已完成任務描述或其 provenance。

## 6. 本次 review 確認問題的修復（不新增需求）

- [x] 6.1 修復已連線但未驗證操作權限仍開放的問題：唯讀判定各操作直接授權，保持未知／作用中角色未驗證／部分撤權 fail-closed，完整授權保留既有 capability revision；以未修改 product 的失敗測試、修復後測試及兩家可丟棄 fixture 實際撤權驗證。
- [x] 6.2 修復長時間 DDL 超過 operation lease 後可能重複執行的問題：建表及缺表修復共用續期，失去 owner 時取消執行並保留 catalog 查證；以 service regression 跨原十分鐘 lease、單 executor、失敗取消及 race 驗證，不宣稱真 DB 曾持鎖十分鐘。
- [x] 6.3 修復試寫未知結果在 capability refetch／無關 workspace 更新後遺失的問題：保留同群組 revision 的原 operation ID／GET 查詢，DDL 預覽仍隨 workspace 更新；補齊修復前失敗、前端回歸、兩家真瀏覽器回應遺失、既定完整檢查與去敏證據。

本節僅記錄使用者要求修復的三項 review findings 與必要驗證，保留前14項 checkbox／描述／既有 provenance，不新增產品能力、其他資料庫驗收或現場平台條件。
