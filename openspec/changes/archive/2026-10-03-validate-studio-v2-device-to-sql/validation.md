# 實作驗證紀錄

2026-10-03 更新：依 A→F 順序執行 F。**SQLite／PostgreSQL 七型別主流程通過，decimal 範圍已確認為codec／SQL層，設備採集未支援；2.1 品質矩陣通過；2.2 SQLite／PostgreSQL 各九項矩陣通過；2.3 quota／disk-full／poison／CAS及兩worker／retry上限通過；2.4 SQLite清理中斷重啟／same-operation及PostgreSQL權限案例通過，最後cleanup witness全部通過；3.1 最後Go full test／vet／lint、tagged suites及前端基準通過，Linux兩次完整mixed通過；3.2 最後審查／verify及舊案移交核對完成。** 下面保留 10-02 的結果，10-03 補驗另列；不能把歷史檢查當作最新工作樹的完整驗收。

## 環境與方法
- macOS arm64（Darwin 27.0.0），Go／Node 版本記在每個 witness JSON；原始碼基準 `b80b0431` 加上尚未提交的工作樹（witness 的 `dirty_worktree: true`）。
- 工具：`scripts/tests/f_device_to_sql/{lib,run,faults}.mjs`（README 同目錄）。真實 `cmd/test_ui`＋內嵌 UI、Playwright 操作、兩個 loopback Modbus simulator、可丟棄的 SQLite 檔／PostgreSQL 16 容器（`gw-wg-pg-test`，port 55432）。沒有 stub 最終 API，沒有直接 INSERT 取代採集。
- 證據：`docs/plans/studio-v2-write-groups/evidence-f/`（`device-to-sqlite.json`、`device-to-postgres.json`、`faults-sqlite.json` 與截圖）。

## 1.1 ProductionDeviceSQLHarness
- 腳本每次重建 binary、重建工作目錄、等前一次的 gateway 結束（`waitForPortFree`；沒有它時，連續執行會連到正在關閉的舊程序）。PostgreSQL 沒有 `POSTGRES_DSN` 時直接失敗，不借 SQLite 通過。
- 新增 `cmd/f_modbus_simulator`（loopback Modbus TCP，暫存器由 JSON 檔每 500 ms 重讀）。

## 1.2 DeviceToSQLiteMixedRows ／ 1.3 DeviceToPostgresMixedRows
- 全程由 UI 設定：兩台設備、同位址 40001／40002／40003／40007、型別 int16／uint16／uint64／bool 的規則與 Tag、SQLite 或 PostgreSQL 目的地、兩個寫入群組（儲存→套用）、首次啟用。
- 獨立 SQL 查證：A 線 `215／1013／true／9007199254740993`（uint64 無精度損失），B 線 `187／777／false／2`；每列 provenance 是四個 `ok`／`good` 的 UTC 取樣；同線沒有重複 bucket；既有「neighbor」列未被改動；UI 沒有未預期的 API 錯誤。SQLite 與 PostgreSQL（BIGINT／BOOLEAN／JSONB）各至少 3 列／線。
- 最後一次結果：SQLite 6–7 列、PostgreSQL 7 列，皆 PASS（兩次連續 PostgreSQL 都 PASS）。

## 2.x 故障與重啟（`faults.mjs`，SQLite 目的地）
| 項目 | 結果 |
|---|---|
| 2.4 確認試寫（UI preview→confirm） | PASS：「Written and read back correctly」＋「Test row removed」；殘留 `gw-test-` 列 0；neighbor 不變。**沒有做**：試寫期間 restart／同一 operation 重試（只有 D 的套件層級測試）。 |
| 2.1 設備整段靜默約 35 s | PASS：產生 bucket 缺口、沒有捏造列、`skipped`／`no_data` 計數有出現、沒有重複。**沒有做**：亂序／重複／bad／stale／late 的 UI＋SQL 逐項矩陣（B 的套件測試涵蓋語意，沒有重跑在 production 路徑）。 |
| 2.2 目的地離線 35 s → `kill -9` gateway → 重啟 → 目的地回來 | PASS：離線期間目的地列數不變、backlog 3–4 筆在本地、恢復後補齊且無重複 bucket。**沒有做**：ACK 後 kill 的精確時點、commit 回覆遺失（`lostcommit` 套件測試涵蓋）、backlog 因 endpoint 變更而不改送（service 層測試涵蓋）。 |
| 2.3 poison partition | PASS：B 線被目的地拒絕後 `quarantined=1`，A 線持續前進。**沒有做**：縮小 quota、磁碟失敗、並行 stale save／activation／test-write（前端 `writeGroupLifecycle` 與後端套件測試涵蓋）。 |
| 重啟後仍持續採集 | PASS（新增，見下方缺陷 4）。 |

## F 找到並修掉的產品缺陷（都影響已歸檔的 A／E 範圍）
1. **第一次儲存目的地回 500**：row-group 預檢在空集合、沒有 connector ID 時失敗。`workspace/write_group_legacy_write.go` 改為回 nil；測試在 `write_group_legacy_write_test.go`。
2. **新規則預設 40001×8 與既有規則點位衝突（UNIQUE 500）、新規則無法指給第二台設備、已存規則可改設備（後端 400）**：`nextFreeStartAddress`、預設落到尚無規則的設備、已存規則鎖定設備選單（`state/sourceRule.ts`、`RuleEditor`）。
3. **群組儲存用舊 workspace revision（409）**：儲存前先 `fetchLatest`（`GroupEditor`）。
4. **規則在設備啟用前被編輯就停用了衍生點位，啟用後什麼都不採集**：`sourcerule/service.go` Update 以 `runtimeRuleEnabled`（就緒度）取代 `設備已 active`；回歸測試在 `service_readiness_gating_test.go`（修正前失敗）。
5. **重啟後所有取樣被 `revision-mismatch` 靜默丟棄**（最嚴重）。成因：Step 3 的 mappings API 建立的 mapping 沒有 rule-candidate 資訊；下次開機的 `SyncDerivedPointState` 才補上 `rule_candidate_id`／簽章，並且規則存成中性縮放（×1＋0）而 mapping 的 pipeline 為 `[]`，兩者簽章不同而被標成 `out_of_sync`。`mapping revision` 包含這些欄位，所以所有已套用的群組在重啟後變成 stale，collector 仍在跑、`write_error_total` 為 0，資料卻不再進入 SQL。修正：
   - `sourcerule/service_mapping_lifecycle.go`：簽章與「手動編輯」判斷把中性縮放視為沒有步驟（`canonicalTransformPipeline`），等價 pipeline 不再改寫。
   - `sourcerule/service_derived_sync.go`（新檔，同時避免 `service.go` 變大）：抽出 `syncDerivedRuleState`，新增 `SyncRuleDerivedState`。
   - `studio_v2_workspace_mappings_handler.go`：mapping 建立／更新後立即同步，讓群組綁定的 revision 就是重啟後的 revision。
   - 測試：`service_mapping_identity_scale_test.go`（中性縮放簽章、重啟同步保持 active；後者修正前失敗）、`studio_v2_workspace_mappings_restart_stable_test.go`（修正前失敗：重啟同步改變 mapping revision）。`SyncRuleDerivedState` 的 service 層測試是與函式同時寫的，沒有 RED。
   - 實測：`RestartResumesCollection` 在修正前 FAIL（重啟後 2 分鐘無新資料），修正後 PASS。

## 已知殘留風險（沒有修）
- **revision-mismatch 的樣本仍然靜默丟棄**：readiness 會誠實顯示 `source-revision-stale`，但交付檢視與 runtime 計數沒有顯示「被拒絕的取樣數」，操作者看到的是 collector 正常、資料不進 SQL。建議後續把拒絕原因與計數接到 delivery view。
- Step 3 的目標型別／縮放與規則不同時，mapping 會在儲存當下就成為 `out_of_sync`（先前是在重啟後才出現）。這是既有「需明確重新套用」設計；本次沒有針對它測試 UI 行為。
- 先前觀察：右側摘要與連線表單仍有範例端點；Write Strategy 版面重疊；device probe 與儲存的競態；connector 自動儲存後 UI 持有舊 identity revision（harness 以 reload 繞過）。
- 偶發：連續跑 PostgreSQL 時曾因群組編輯器未關閉而逾時，harness 已加入 `showGroupList` 重試；產品側的原因沒有調查。

## 3.1 基準（本機 macOS arm64）
- `go test ./...` 通過；`go vet ./...` 通過；`golangci-lint run ./...` 0 issues（既有 `service_probe.go:201` 以說明過的 `//nolint:gocritic` 處理：`%#q` 會把含反引號的名稱改成雙引號，不能直接替換）。
- 前端：`npm run lint` 通過；`npm test -- --run` 180 files／1081 tests 通過；`npm run build` 通過；repo 的 `npx playwright test` 13 passed。
- `git diff --check`、`make check-lines` 通過（`sourcerule/service.go` 1192→1169）；`spectra validate --all` 全部 valid。
- **沒有驗證**：Linux、Windows、ARM 以外平台、嵌入式瀏覽器、LAN、真實 PLC、SCADA。

## 2026-10-03：2.1 QualityBoundaryMatrix

- 實際命令：`node scripts/tests/f_device_to_sql/quality.mjs`，macOS arm64、Go 1.27.1、Node 24.15.0。使用 `f_write_group_fixture` binary、真 UI／Modbus read／runtime mapping／AcceptSample／journal／closure／outbox／SQL sender；一般 binary 的 fixture 控制另測不可用。啟動前 DB ownership 檢查在 open／migration 前執行，simulator 與 controller 均只 bind loopback。
- witness：`evidence-f/quality-sqlite.json`，run `gw-f-quality-1790998885212`，binary SHA256 `7f8f0e8a349a7b5c5e308aa7a2f387a300e83b192da33606780554c157da1e0f`；8/8 assertions PASS，`ui_problems=[]`。原始碼仍為 `b80b0431` 加 dirty 工作樹。`quality-line-a.png`／`quality-line-b.png` 已實際檢視，畫面桶數／最近原因與 SQL／delivery API 一致。
- 亂序 t=8→t=3 後 SQL 選 321，duplicate 不新增 journal，重用 actual read identity 而值不同明確拒絕，closure 後 late 拒絕且 SQL 不變；t=10 寫入下一桶的 123。uint64 獨立查回字串 `9007199254740993`。
- default missing、actual failed read 的 bad、max-age=2 的 stale、完全不 poll 的 silent 均沒有完整 SQL row。explicit partial 經 canonical Save/Apply；只寫必填 temperature=123，其他三欄均為真正 NULL，provenance 保留 missing 原因。最終 SQL 只有 A 的三列，B 無列；沒有把前桶值補入缺值欄。
- 真實 harness 首次抓到 policy Apply 後漏掉 lookahead 的雙 revision intake：同一 bucket 舊版仍寫 row、新版判 stale。失敗 witness 保留為 `quality-sqlite-first-failure.json`。現在 reconcile 依已生效 revision，對更早 boundary 設定原 interval 對齊的 cutoff；不提前退休 future prebuilt boundary，cutoff 不延長。新 `TestReconcileAfterEffective...` 修正前 journal 仍有舊版而失敗，修正後通過；完整 pipeline suite 與相關 race suite 通過，獨立唯讀 review 未發現 confirmed defect。
- 同時補齊兩個原定整合缺口：partial readiness 沿用 B 的 `NewGroupRowLayout`（nullable／provenance／entity key 能力）；UI 顯示 bounded 最近 20 桶的 safe enum 原因，讀取失敗不沿用快取宣稱已確認。nullable positive、entity key negative、原因／上限、UI counts／原因／快取失敗均保留修正前失敗的測試；workspace／groupdelivery／grouppipeline suites 通過，受影響前端 4 files／47 tests、build 與 scoped lint 通過。
- 受影響 Go scoped lint、fixture normal／tagged lint 均 0 issues；normal／tagged focused tests 與 vet 通過。初始化錯誤 clock 的新負測試修正前失敗、修正後通過，fallback 不清除原始錯誤。`git diff --check`、`make check-lines`、F `spectra validate` 通過；analyze 無 Critical／Warning，6 個 Suggestion 是其他 scenarios 尚缺具體 Example。2.2–2.4 的精確故障邊界、Linux、最新完整 Go／frontend／Playwright 基準仍待執行；此矩陣不代替這些缺項，也不證明現場設備或正式資料庫可用。

## 2026-10-03：mixed、故障補驗與最新基準（進行中）

- 原生 string 解碼器原本落入數字 fallback。新增實際 Modbus adapter／decoder regression，修正前失敗；現在沿用 BigEndian packed ASCII 與既有尾 NUL 處理，單一 string 點位仍只有一個 register，不擴充長度或 decimal 型別。float32／float64／text 連同四種既有型別，實際 SQLite 每線至少三列，七筆 persisted member provenance／UTC／good 一致；neighbor 未改。`device-to-sqlite-mixed.json` 仍為 FAIL，原因是 Step 3 preview 400，不能因 SQL 成功就勾選 1.2。
- Step 3 中性原生 string preview 的單元 RED／GREEN 已完成；非中性 scale 保留錯誤狀態。最新前端 `npm run lint`、Vitest **181 files／1088 tests**、build、`npm run test:e2e` **13/13** 通過。E2E 首次受 sandbox Chromium 權限阻擋，升權後實跑通過；日誌在 `/private/tmp/go-gateway-f-latest-frontend-*`。完整 UI harness 仍有 preview 400，正查 actual payload，不能以單元或 E2E 基準取代 mixed 驗收。
- F2.2 新 controller 確實在真實 target Commit 後注入 lost response／hold；local receipt trigger 與 closure barrier 只作用 owned DB。review 找到 current draft destination 被誤當 active destination，具名 RED 重現後改為 `ResolveAppliedAt(fixtureClock.Now)`；無 active snapshot 安全拒絕。closure hold 時 `/state` 被 pipeline lock 阻塞，改用 tagged-only `/faults` 獨立觀察，實際 held transaction HTTP regression 通過。normal／tagged focused、race、vet／lint 已跑，但不是 F2.2 完成證據。
- SQLite run `gw-f-recovery-sqlite-1791002955364` 九個行為 assertions 曾全部通過；總列數 assertion 漏算共用 B 點位向 B/C 兩個 active group 的合法交付，失敗保留在 `recovery-sqlite-shared-source-count-failure.json`。修正 assertions 必須另外核對 B 的 effect，不能忽略該列。後續重跑 `gw-f-recovery-sqlite-1791003188513` 在 outage 恢復出現 `target-unusable`，保存 `recovery-sqlite-single-file-outage-failure.json`。現有 Stat 後 rename／open 會建立空 SQLite 的窗口已由 deterministic hook test 實際 RED 證實，正在修復 delivery-only atomic existing-file open；尚未宣稱整批 PASS。
- PostgreSQL F2.2 首次 setup FAIL：第二 connector 使用錯誤 `username` 欄位；實際 driver 契約為 `user`。Apply 正確 fail-closed。失敗保留 `recovery-postgres-b-apply-setup-failure.json`，已修 harness 欄位並補 actual readiness assertion，尚未重跑，不能當產品 regression RED。
- F2.3 低 quota 的實際修正前 binary 仍 ACK 第二個真實 Modbus sample，RED 保存在 `capacity-quota-red.json`。tagged startup quota override／固定 `/capacity disk_full` page gate 已完成 focused、normal、race、vet／lint；用 actual modernc SQLITE_FULL 測到拒絕且保留已接受資料。完整 quota／poison／CAS UI＋SQL harness 已建，**尚未執行**；dispatcher overlap／explicit retry exhaustion仍需要獨立 production witness。
- Linux 已實跑 Go 1.25.5、Node 24、Playwright 1.58.2 的 isolated UI→Modbus→SQLite，amd64 userland 在 arm64 Docker Desktop 上模擬。SQL mixed rows 成功、preview 400 導致整體 FAIL；保存 `device-to-sqlite-mixed-linux-first-failure.json`／`linux-environment-first-failure.json` 與截圖。image ID／source manifest digest 在環境 witness。不是原生 Linux 硬體、Windows／ARM deployment／embedded browser／LAN／PLC／SCADA 驗收。這個 snapshot 不取代後续原始碼的完整 Go 基準。


## 2026-10-03：本次最新實測（取代上段的進行中狀態）

- Step3真實400 payload為decode(string)→scale(1,0)→cast(float64)：native string的預設target仍float64。mappingDefaults改為string，numeric不變；default regression RED/GREEN後，SQLite七型別主流程 `device-to-sqlite-mixed.json` PASS6rows、7-member provenance、UTC/good、neighbor不變，mapping_preview_requests=[]。PostgreSQL `device-to-postgres-mixed.json` 亦PASS6rows，float32/float64/text與uint649007199254740993精確保留。decimal未新增或刪去驗收承諾。
- 最新前端lint／Vitest **181files／1089tests**／build／Playwright **13/13** PASS。最新normal `go test ./...`、`go vet ./...` PASS，完整lint首次指出simulator的evalOrder，改為先StartWithConfig再return後完整lint **0issues**、完整Go test再PASS；logs `/private/tmp/go-gateway-f-final-go-{test2,vet,lint2}.log`。這不是Linux/現場驗收；後續tag-only secondary identity另有focused證據，final tagged gates尚待重跑。
- SQLite `recovery-sqlite.json` run **gw-f-recovery-sqlite-1791005821557**、binary **9ff044ca100601899c094f185c1d37a9f548535036c21230831adc60e09dda5e**：9/9 PASS、ui_problems=[]，實際A6/B5/retarget0 rows、A6/B2target receipts。先前empty-file窗口已修為delivery-only atomic existing-file open；ro、memory URI/plain filename、unknown mode另有RED/GREEN，沒有擴改schema/probe。
- `capacity-sqlite.json` run **gw-f-capacity-1791005793308**、相同binary：10/10 PASS，group/global quota均拒絕新ACK且暴露scope/loss-risk，真正SQLITE_FULL無新增journal；恢復後SQL/provenance/receipts證明accepted buckets交付，poison保留payload只阻同partition，CAS真實200/409、staleApply/preview拒絕、UI409保留草稿/重載。ui_problems只有預期409。第一個empty poll是target離線重啟無metadata可建boundary的harness时序缺口，改owned SQLite BEGIN IMMEDIATE reserved lock保持metadata可讀；另一失敗是錯誤要求成功交付後永久保留consumed journal，已改依snapshot semantics查SQL provenance。兩個failure JSON均保留。
- `worker-overlap-sqlite.json` run **gw-f-workers-1791006569118**、binary **c4e2bc3102a8e1f15047eaf96202707c5b9fbc4439fb7c2737f217b03bac0870**：兩個actual production pipelines共用owned store；不同node/incarnation下live claim不被取走、另一worker真實交付healthy partition、同group successor持續pending；release後每effect只一個SQL/receipt。tag-only F_FIXTURE_MAX_RETRIES=2真實兩次target-unavailable後blocked，payload保留、B持續、UI誠實。normal build忽略ENV。第一輪同node兩process被既有restart recovery當成舊incarnation，屬harness身份錯置，保留failure；只增加白名單secondary的test identity，沒有改productionfencing或新建sender。
- `test-write-sqlite.json` run **gw-f-test-write-sqlite-1791005883773**：UIpreview零targetmutation，confirm written_verified/cleaned，actual採集建立的A/B neighbor rows不變、completed same-operation retained200/GET不改SQL。精確owner_value外部BEFORE DELETE trigger只叫既有tagged scalar；ledger phaseCleaning＋target test row存在時SIGKILL，DELETE回滾、livelease回202，真實三分鐘後same-operation只續清理、不重寫；獨立SQL以TEXT確認兩欄9007199254740993、最終SQL neighbor不變，UIretained-result清理成功。SQLite權限模型不具有PG grants，noSELECT/noDELETE在這個witness為N/A；必須另有PG對照才算F2.4完成。
- Cleanup缺陷RED：DELETE後SELECT失敗原誤報cleaned；現在unknown＋safe readback reason，完整grouptestwrite/race/vet/lint PASS。日志 `/private/tmp/go-gateway-f2-4-cleanup-readback-{red,green}.log`。PG權限首次harness獨立connector不符saved Step4 destination，UI正確disabled；保留setup failure，未作產品RED，owned orphan角色/schema/database已精確清除，改由真UI先保存role destination再確認。
- Linux最新pinned容器有真實UI/SQL執行，preview0errors；但B sql_committed1、queued7、unknown1，整體FAIL，正在保存/診斷實際delivery/commit證據。前兩次.git metadata的純harness錯誤亦保留。不能從macOS/PG PASS或amd64模擬container推論native Linux／Windows／ARM deployment／embedded browser／LAN／PLC／SCADA。

## 2026-10-03：PostgreSQL 權限實測與測試證據修正

- `test-write-postgres.json` run **gw-f-test-write-postgres-1791008317225**、binary **eb57bbf822ce43fcf7c6f68515560d0b00b0e53e7d1494deb3b508aed4d71efe**：normal preview/confirm/readback/owned cleanup、noSELECT、noDELETE 三案實跑通過，neighbor不變與retained200/GET相等。真UI缺SELECT顯示written_unverified/readback-denied與cleanup failed/cleanup-denied；缺DELETE顯示written_verified但cleanup failed/cleanup-denied。兩張權限截圖已實際檢視。
- 權限 RED 使用真UI與lib/pq driver：初始回generic cleanup-failed；typed regression確認PQ SQLSTATE42501／row／schema／authentication錯誤原誤判transient。加入pq.Error分支共用既有SQLSTATE分類，dbtarget／grouptestwrite／groupdelivery完整package GREEN；不改None unknown保護。日誌 `/private/tmp/go-gateway-f-pq-classification-{red,green}.log`。
- PG先前四個setup failures（destination gate、schema locator、legacy ownership guard、pool inventory）均保留，不能當產品RED。最終setup從初始真UI使用run專用非superuser role，只撤本run readings SELECT/DELETE，connector identity不變。
- 獨立review發現secondary listener未綁owned PID及witness先於fixture cleanup寫出；已增加雙listener PID核對與cleanup逐step結果，成功報告在cleanup之後才寫。四項helper測試通過，刻意讓cleanup回假成功時三項測試失敗，還原後GREEN；日誌 `/private/tmp/go-gateway-f-cleanup-witness-{mutation,green}.log`。上述PG實跑在此helper修改之前，最新cleanup witness仍待重跑，不能用helper單元測試代替實跑。
- `recovery-postgres.json` run **gw-f-recovery-postgres-1791008708406**、binary **567a53c6fd4cf0ae8b1d17ffbbd7e326b624678ed5063d54513129696a04041e**：九項實際ACK／closure／target commit／local receipt故障、endpoint revision與None unknown皆通過；不是借SQLite通過。PG recovery與test-write完成後唯讀核對exact-run schemas／other database／permission role剩餘數皆0，兩witness加入post_run_cleanup_verification。

## 2026-10-03：使用者確認型別範圍與SQLite提交鎖修正

- 使用者確認維持既有七型別，取代F1.2/共用fixture的decimal真採集承諾；proposal/design/spec/tasks與acceptance/contracts同步。已完成checkbox描述與B codec歷史證據保留，取消設備decimal不假勾為已實作。ingest analyze零Critical/Warning，六項Suggestion是舊抽象scenario可補examples建議；validate PASS。
- 真production OpenDestination→InsertGroupRow，外部DELETE-journal reader SHARED lock造成commit phase rawSQLITE_BUSY code5（behavioral RED），修正僅delivery DSN未明示busy_timeout時預設15000ms；約195ms等讀取鎖釋放後成功。explicit0、ro/memory/unknown mode、atomic existing-file open保留，不改WAL/pool/retry/dedupe或None unknown。normal/tagged package/race/vet/lint PASS。日誌 `/private/tmp/go-gateway-f2-linux-reader-{red,green}-moved.log` 與 `/private/tmp/go-gateway-f2-linux-lock-*-final2.log`；兩次完整Linux實跑仍待新witness。
- `worker-overlap-sqlite.json` run **gw-f-workers-1791009081533**、binary **567a53c6fd4cf0ae8b1d17ffbbd7e326b624678ed5063d54513129696a04041e**：兩worker與retry上限2/2通過；secondary雙listener核對owned PID75270，fixture_cleanup逐step全部PASS，報告在owned children已退出之後寫出。
- Decimal codec/SQL層此次重新執行：`POSTGRES_DSN=<owned loopback> go test ./internal/datalink/dbtarget ./internal/datalink/measurement -run ^TestExact -count=1 -v` PASS，含真SQLite及PostgreSQL decimal round-trip，值1234567890.123456789012345678精確；不是設備採集。既有Mac七型別witness僅追加scope_decision_annotation，原run/source/build/actual values/pass未改。
- 最後Go full test初次FAIL在既有outage SQL observer（busy_timeout=0），不是資料比對錯誤；production提交改有限等待後，觀察端也需等待短鎖。新增real exclusive-lock讀取回歸直接重現SQLITE_BUSY，RED日誌 `/private/tmp/go-gateway-f-observer-read-red.log`。只對這個已失敗test observer設定有限busy_timeout15000，保留全部rows/isolation assertions；GREEN/最後full gates待結果。quality-lib同pattern未有本次failure，保持不改。
- 新Linux run1為/work noexec tmpfs的esbuild EACCES setup failure；run2即時SQLite CLI observer lock失敗，不能宣稱PASS。run2 shutdown後A/B各三列、七型別exact、6committed/receipts、unknown0，已保存與原FAIL分開；正修有實際證據的只讀CLI有限等待，重跑兩次完整鏈，不借post-shutdown查詢升格為整體PASS。
- 最新SQLite cleanup完整UI實跑 **gw-f-test-write-sqlite-1791009564180**、binary **567a53c6fd4cf0ae8b1d17ffbbd7e326b624678ed5063d54513129696a04041e** PASS：精確uint64、真三分鐘lease接手與owned cleanup各step全部通過；permission仍N/A，由已實跑PG counterpart補證。

## 2026-10-03：最後執行結果

- Linux fresh run3 **sqlite-linux-20261003T064744743Z** A4/B3、run4 **sqlite-linux-20261003T065155156Z** A3/B3完整PASS；gateway **777754f1ae2520f83e659847b3d20463332b4c9718bbb2610f11b2a4d5fe3849**、simulator **da3fb8e2d8e037d1fce6b62f959e898df9212e264980aa6678778b6d8a94f363**。七型別exact／7-member UTC/good／neighbor不變／UI API preview errors0／unknown0／queued0，run3/run4 raw與compact witness、manifest/environment皆保留。Linuxx64 Dockeramd64 emulation onarm64 host；不是native hardware或ARM部署驗收。
- SQLiteCLI reader timeout回歸：無等待rawbusy exit1、timeout15000真讀A|215 exit0、timeout100仍失敗，沒有吞錯作成功。原run1 noexec setupfail、run2 observerfail與postshutdown SQL保留，不改成PASS。
- 最後 `go test ./...`、`go vet ./...`、`golangci-lint run ./...` PASS（0issues）；`go test -tags f_write_group_fixture ./cmd/test_ui ./internal/datalink/dbtarget ./internal/datalink/runtime -count=1` PASS。日誌 `/private/tmp/go-gateway-f-final3-*`。先前final2 observerFAIL未刪；新exclusive-lock觀察測試及原outage場景連續3次GREEN。
- 最後PG權限run **gw-f-test-write-postgres-1791009949489**、binary **567a53c6fd4cf0ae8b1d17ffbbd7e326b624678ed5063d54513129696a04041e**：normal/noSELECT/noDELETE三案PASS，fixture_cleanup browser／ownedchildren／port／run-role／ownedtargets全部PASS，報告在清理後寫出。SQLite counterpart **gw-f-test-write-sqlite-1791009564180**同binary，restart/same-operation/precision/UI與cleanup全PASS（permissionN/A由PG實跑補證）。
- 前端維持此次最後mappingDefaults修正後的lint／181files1089tests／build／13/13Playwright PASS，未再改frontend inputs；保留既有當次日誌，不用Go或Linux witness代替frontend suites。Windows/nativeARM/nativeLinux hardware/embeddedbrowser/LAN/PLC/SCADA/productionDB/deployment/長時間soak仍NOT RUN。


## 2026-10-03：最終核對與清理錯誤通道

- [final-verification.md](../../../../docs/plans/studio-v2-write-groups/final-verification.md) 完整對照 F 三條 requirements／八個 scenarios／一個具體 example，以及舊案五條保留需求／十八個 scenarios；已退休 row autosave、advanced picker 與現場未驗情境分開列示。原五項完成 task 文字未改寫。
- 最新完整 runs：Mac SQLite `sqlite-mac-final-20261003T0745`、PG `postgres-mac-final-20261003T0750`、Linux `sqlite-1791013947934` 均真七型別 UI→SQL PASS。quality 8、capacity/CAS 10、worker 2、SQLite cleanup 2、兩 DB recovery 各9、PG permission 3 全部 PASS；完整 run IDs／JSON hashes／binary identity 見 evidence-f/final-evidence-index.json。
- 獨立 review 最後找到 helper 返回 cleanup false 未設定 exitCode，及 Go tagged cleanup error 只進記憶體 counter 的假成功通道。前者真 RED 產 passing report，後者閉 DB 後 callback 回零值；現在各 cleanup step 嘗試後 errors.Join、main defer 回 error、owned child 非零退出即 failure、publish 強制 passed=false。正常 default／held closure cleanup 仍 nil，root 親讀全部新 patch，沒有新增控制 endpoint。
- 修後 Mac SQLite `sqlite-final-cleanup-20261003T1615` A4/B3 七型別 PASS，binary `1fb4d02d611ef1ca9604528e885f75c800d29b1901506ffd481847dfb33cecab`；PG cleanup `gw-f-test-write-postgres-1791015311503` 三案 PASS，binary `4b1fec85e5fd9f188eaa731bede465d79dd9719302e938d61848538ed8f5a37a`，全部 children exit0、owned schema/role/database absence 查核為空。修前 witnesses 另存 prior-cleanup-error-channel；Linux 與其餘矩陣保留修前實際 source/build，不冒稱修後 binary。
- 清理修正後最後一般／tagged full `go test -p 1 ./...`、兩版本 vet／golangci-lint 全 PASS（兩 lint 0issues）；harness 四測試檔 24/24 PASS，logs `/private/tmp/go-gateway-f-final5-*`。final-code-inputs.json 115份 source/config/test SHA256 執行後未變。前一輪 full suites 同時跑造成 fixed5020 port collision，改串行通過；sandbox cache／Chromium 阻擋另存，不當產品 RED。
- 前端最後 lint／Vitest181files1090tests／build／Playwright13/13 PASS；late Save A 回覆覆寫正在編輯 B 的真 RED 已修，之後前端來源未變。setup 截圖只證明設定／activation，cached0/not-yet-confirmed 不冒充最後 SQL 數；quality/recovery 另重讀 UI 與 SQL/API 一致，權限截圖如實分開 write/readback/cleanup。
- review 四維度與 scope 邊界核對；preexisting dirty E/F shared paths 不能僅以檔名證明全部 hunk 歸屬，binary 診斷內容不當程式碼審查。closeout 獨立核對零具體缺陷／零重大未移交需求。完整 scope 在最後 task tracking 更新後重新 capture/check；validate/analyze/diff/line/links 通過，analyze六項抽象example建議保留。
- Migration Plan／Open Questions 已核對：024/025/026 additive/current suites、immutable backlog、safe unknown/no retarget 保留；不對正式資料試 migration、rollback或部署。unique_key正向因canonicalRecordKeyColumn未提供而NOT RUN，未加未授權功能。Windows/nativeARM/nativeLinux硬體/embedded/LAN/PLC/SCADA/productionDB/longsoak仍NOT RUN。
