# 六案交付核對

2026-10-03；A–F 共 54 項實作與驗收已完成，舊案移交亦核對完成。Spectra tasks／archive 狀態與各案 validation 保存交付紀錄。

## 驗收範圍與證據

使用者已確認依 A→F 實作及驗證；七種真採集型別為 int16、uint16、bool、uint64、float32、float64、string。decimal 僅 codec／SQL round-trip，設備採集未支援。沒有新增型別推測或部署。

真主流程使用 production `cmd/test_ui`、embedded frontend、實際 Studio V2 UI、兩台 loopback Modbus simulator 與 disposable targets。各 JSON 保留 run ID、source SHA／dirty worktree、binary SHA256、命令、persisted IDs、獨立 SQL 與清理結果；SQL schema 控制和 neighbor fixture 不代替實際採集。

| 實際執行 | 最新完整 run | 結果 |
| --- | --- | --- |
| macOS SQLite 七型別 | `sqlite-mac-final-20261003T0745` | PASS；A/B 各 3 列、每列 7 members、uint64 精確、neighbor 不變、unknown=0 |
| macOS PostgreSQL 七型別 | `postgres-mac-final-20261003T0750` | PASS；同樣 A/B 各 3 列，run-owned schema DROP 後獨立 absence query 為空 |
| Linux SQLite 七型別 | `sqlite-1791013947934` | PASS；A/B 各 3 列、7 members、unknown=0、owned children exit 0；pinned image `194d0b18626683c190f4598fc8750d6417288dfaa090b6b0caaeb6f8624ebdc9` |
| QualityBoundaryMatrix | `gw-f-quality-1791013596927` | 8/8 PASS |
| Capacity／poison／CAS | `gw-f-capacity-1791013739754` | 10/10 PASS |
| Worker overlap／retry cap | `gw-f-workers-1791013847850` | 2/2 PASS |
| SQLite confirmed test-write | `gw-f-test-write-sqlite-1791013951300` | 2/2 PASS；真三分鐘 lease、restart、same-operation；PG grant 案例 N/A |
| SQLite recovery | `gw-f-recovery-sqlite-1791014233199` | 9/9 PASS |
| PostgreSQL recovery | `gw-f-recovery-postgres-1791014420694` | 9/9 PASS |
| PostgreSQL test-write | `gw-f-test-write-postgres-1791014636201` | 3/3 PASS；normal／no SELECT／no DELETE；run-owned role/schema/database 查核皆 0 |

證據位於 [evidence-f](evidence-f/)，[索引](evidence-f/final-evidence-index.json) 保存各 JSON 的 SHA256。各 fault witness 的 `report_errors=[]`、`fixture_cleanup.passed=true`。以上完整 runs 先於最後「關閉錯誤回傳／cleanup false 不得報 PASS」修正；主流程、取樣、writer、SQL assertions 不變。修後驗證另外列出，不把舊 binary SHA 改成新版本。

修後 SQLite `sqlite-final-cleanup-20261003T1615` 再實跑 PASS：A4/B3、每列七種型別及 provenance、uint64 精確、UI/API/preview errors=0、全部 children exit 0。修後 PG `gw-f-test-write-postgres-1791015311503` 三案再次 PASS，cleanup 各 step 全通過，owned schema／role／database 獨立查核均為空。修前報告另存 `prior-cleanup-error-channel`，沒有改寫歷史執行身分。Linux 未為這個純關閉錯誤通道修改重跑，仍以表中 manifest／binary 的實跑為準。

## F 規格逐項核對

| Requirement／Scenario／Example | 實作或測試證據 | 判定 |
| --- | --- | --- |
| Production device-to-SQL witness | `run.mjs`、`mixed-identity.mjs`、真 UI＋SQL 三路完整 JSON | 已實跑 |
| Mixed multi-device happy path | persisted device／point／Tag／group、同位址不同 device、UTC／quality／values／neighbor assertions | 已實跑 |
| Supported acquisition types and decimal boundary | 七型別三路 JSON；`TestExact*` 真 SQLite／PG decimal codec/SQL tests | 採集七型別與 codec 分開通過；設備 decimal NOT RUN／unsupported |
| Operation cleanup | `test-write.mjs`、grouptestwrite execute／ledger；neighbor／owned owner value／retained operation assertions | 已實跑 |
| Failure recovery and revision acceptance matrix | tagged fixed fault controls＋production journal／outbox／sender，quality／capacity／worker／recovery witnesses | 已實跑 |
| Outage recovery | 兩 DB 的 ACK kill、closure rollback、target commit loss、local receipt、restart、effect key／SQL／UI | 已實跑 |
| Concurrent settings | capacity 真 CAS 200/409、stale Apply／confirm、UI held draft；recovery endpoint revision blocked | 已實跑 |
| Incomplete bucket | default missing／bad／stale／silent 無完整列；explicit partial 真 NULL 與 reasons | 已實跑 |
| deterministic production quality boundary example | t=8→t=3、duplicate、t=10 next bucket、late release unchanged SQL；quality 的具名八項 assertions | 已實跑，並非 fixture 直接造 sample |
| Evidence distinguishes actual execution and field limits | witness 的 source/build/commands/platform/limits；Linux manifest/environment | 已核對 |
| PostgreSQL unavailable | harness preflight 不接受缺 DSN／非 owned target，mixed-harness-regressions 負測試；未以 SQLite 冒充 PG | 防護測試通過；本次 PG 實際可用 |
| Linux simulator passes | Linux amd64 container 在 arm64 host 模擬，有自己的 binary／source manifest／UI／SQL | 已實跑；原生硬體與現場另列未驗 |

## 舊案保留的 18 個 scenarios

唯一 owner、原五項完成文字、migration／rollback 與需求實作對照見 [handoff-closeout.md](handoff-closeout.md)。下列區分目前可執行情境和 E 已退役的 UI，沒有替舊 UI 補上虛構通過數。

| 舊 Scenario | 現行證據與界線 |
| --- | --- |
| Another row is refreshed | `writeGroupLifecycle` late Save A→close→edit B 測試；修正前 B 被 A 覆蓋，修正後保留 B |
| Same row changes remotely | writeGroupSection／writeGroupLifecycle 的 local draft、refreshed revision、409 conflict／explicit reload |
| Enter followed by blur | 現行 GroupEditor explicit Save；已移除逐列 Enter／blur autosave，該舊操作情境不適用，沒有聲稱當次執行退役 UI |
| Chinese composition or deleted row | 現行群組沒有 IME／blur mutation handler；已卸載 editor 的遲到 Save 不重開，late callback 回歸通過；退役 row keyboard 情境不適用 |
| Old preview arrives after switching tables | SchemaSetupSection scope signature 與 step4-database 的 stale preview／scope change tests；群組試寫綁 immutable revision |
| Activation or readonly state is active | step4-database-controls／writeGroupLifecycle readonly、schema/test-write／mutation guards |
| Rapid clicks or navigation during a write | double-click／lost response／running／retained operation tests；same-operation 真 UI／SQL，沒有宣稱 backend cancel |
| Reuse a saved connection | saved connector identity resolver、connector_identity／workspace_database_identity tests；masked password 不變成新密碼 |
| Select a different database kind | server-derived dialect／connector guards與真 SQLite／PG UI→SQL，各自獨立結果 |
| Password includes leading or trailing spaces | dbtarget `service_postgres_test`／connector identity tests 保留 exact bytes、replace／clear／endpoint revision |
| Connection is absent or inaccessible | unsaved／foreign／stale connector target resolver 與 readiness tests，拒絕猜 identity |
| Plan contains measurements from two devices | recording plans membership／members route tests，persisted measurement 的各自 equipment identity |
| Missing disabled or foreign measurement | 同 route tests／membership validator；缺測量、disabled、foreign、unwired reader 安全拒絕 |
| Multiple plans or plan-list failure | E 已移除 basic UI 的 advanced plan picker 與自動建立／選第一筆；此 basic 主線沒有可觸發的 plan-list 操作，舊 picker 情境未實跑、不適用；advanced API ownership 契約保留 |
| Foreign plan read update or deletion | workspace-scoped recordingplan get/update/delete 回同 safe404，foreign state 未變 |
| Selection changes after a test or preview | schema／test-write scope signatures、connector/group revision 防護；已移除 basic advanced picker，未聲稱進階 UI 驗收 |
| Second local save step fails | AtomicSetupSave actual SQLite row-group／target-ref 第二段失敗 rollback、first-save／reload tests |
| Stale save or unresolved external operation | AtomicSetupSave expected_setup_revision，operation unresolved／readiness fail-closed tests與 F CAS |

## 最新基準與最後審查

- 最後清理修正後 `go test -p 1 ./...`、`go test -p 1 -tags f_write_group_fixture ./...`、一般及 tagged `go vet ./...`／`golangci-lint run ./...` 均通過；兩完整 lint 皆 0 issues。日誌為 `/private/tmp/go-gateway-f-final5-*`，[115 份程式／測試／設定輸入 SHA256](evidence-f/final-code-inputs.json) 與執行後工作樹一致。兩 full suites 前一輪並行碰到既有測試固定 5020 port，串行重跑通過，不改測試 assertions。
- 前端 `npm run lint`、`npm test -- --run` **181 files／1090 tests**、`npm run build`、`npm run test:e2e` **13/13** 通過。late Save 修正後沒有再改前端；不是以 focused suite 代替 full gate。
- harness 共同 HTTP deadline／token witness／cleanup／run namespace／identity 四個測試檔 **24/24 PASS**；never-response/body deadline 真的逾時，caller cancellation 保留。新增 API 未知值不當成空成功。
- Review 範圍使用 Spectra touched tracking＋Git base `b80b0431e6424aece6675d31ee46726c89d19478`，檢查 correctness、efficiency、reuse、convention；task-start dirty paths 只能證明檔案範圍，不能證明所有 hunk 屬 F。未提交 E／無關 WIP 保留。
- 最後確認的缺陷是 cleanup 回 false 卻可能出版 passed=true，及 tagged cleanup errors 未傳到 child exit。兩者已限共同 helper／fixture shutdown 修正：真行為 RED 原報 PASS、閉 DB 後 cleanup 原回零值；GREEN 明確回 error／非零 exit／report false，正常關閉及 held-closure 清理仍回 nil。root 親讀最終 patch，獨立 reviewer 的 closeout 核對為零具體缺陷／零重大未移交需求。
- reviewer 提到 unmounted editor 的遲到 error state，但 keyed editor 已卸載，沒有可重現的新 editor 覆寫，未列為缺陷或追加功能。
- 實際檢視 Mac／Linux setup、quality 與 PG 權限截圖。setup screenshot 的 cached 0／not-yet-confirmed 只證明設定及 activation；SQL count 由獨立 SQL/API assertion 證明，quality／recovery 另重讀 UI 核對，不把 setup screenshot 當最終交付數。
- 安全檢查限本次變更的 typed input、loopback binding、owned disposable DB、literal secret bytes、bounded errors、fail-closed、child/resource ownership；沒有宣稱外部正式環境安全稽核。
- F／舊案 `spectra validate` 均 PASS；analyze 零 Critical／Warning，分別六／十七項 Suggestion 為抽象 scenarios 可補具體 example。行數／diff／本次文件連結檢查通過，沒有以建議增開產品功能。

## 移轉、回復與未執行範圍

A 024 additive contract migration、C 025 durable recording、D 026 operation ledger 的 current suites 與 legacy adapters／restart／ownership 規則通過。F 不對正式資料做 migration；回退 binary 不撤銷外部 DDL／committed effects，先停新 intake，保留 immutable revisions／accepted journal／outbox／receipt，unknown 不盲目重送，endpoint 改變後 backlog 不轉送。

`unique_key` 正向 production 場景 NOT RUN：canonical RecordKeyColumn 尚未提供；不新增未授權功能以湊驗收。none unknown 與 target_receipt 路徑各有實跑。

Windows、原生 ARM／Linux 硬體、embedded browser、LAN、真 PLC、SCADA、正式 DB、正式 backup／restore、部署與長時間 soak 均 NOT RUN。已通過的 simulator／Docker／本機 suites 不升格為現場驗收。沒有 stage、commit、push、merge 或部署。
