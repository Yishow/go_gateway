# 各點位分別計算：後端契約與驗證交接

更新：2026-10-31。正式產品基準：`22520a2e06b4f15425b8452979aff87cca632169`；本收尾另有精確123檔lint清理，尚待root發布。
範圍：[原 change](../../../openspec/changes/support-per-member-recording-calculation/proposal.md) 的 1.1–4.2 與 4.3 收尾，不包括其他 editor、採集批次或 Modbus 改動。
本文件不是整個產品、現場部署或瀏覽器驗收的通過證明。

## 使用與交接界線

- 同一記錄的 canonical member 各自保存 `latest`／`mean`／`delta`、cadence 與 identity；不修改 global Point／Tag，不另外建立 reducer、writer 或 baseline authority。
- 新形狀為 `mixed`／`member-calculation-v1`。既有 instantaneous/cumulative、raw 欄、create-once 與歷史資料不被升級重解。
- UTC epoch 半開窗 `[start,end)`；5／15／60／600 秒有驗證。窗口時間、到期時間、實際 SQL commit 時間不同，不能承諾整秒落盤。
- wide 在到期聯集輸出；未到期欄 `NULL + not_due`，不沿用上次值。long 僅到期 member 產列。
- 到期 required 無可靠結果：wide 跳過整列；long 只跳過該 member。optional 只在 fresh nullable/provenance proof 成立時寫 NULL，不能補 0、raw 或跨窗 span total。
- 僅 mixed 且交易內確認無 active／historical applied snapshot 的首次 Apply 使用本次 server clock；首次 partial 保留 aligned identity 與 actual coverage。歷史／resume 不走首次捷徑。mixed 換 revision 使用 checked old/new cadence 共同切界，不保存假的 group interval；legacy 普通 Apply 原語意不變。
- 同源獨立 mixed records 可在不同 owned targets 使用不同 method/cadence，獨立啟停與 baseline；canonical record、lineage、target ownership 均仍要成立。
- 仍受既有獨占來源 ownership 保護的 legacy wide/custom 情況會 typed 阻擋；本次真 fixture 驗證409後，由操作員明確 Disable，再以新 RecordID／lineage／owned target 與 fresh proof 啟用。其他既有合法雙 profile／shared-long 權限未改，不能泛稱所有 legacy wide 都必須先停。沒有自動停用、刪除、同 lineage 轉型或 baseline 複製。
- shared long 僅沿原 storage ownership／row partition，取得 fresh mixed shape、nullable、provenance 與 schema proof 後共表；既有歷史列完整保留。

### 能力旗標不是每筆記錄的啟用許可

`draft_supported=true`、`activation_supported=true`、`editor_supported=false`，原因為 `record_readiness_required`。
這表示後端有已驗證的支援路徑，不表示特定 record 已 ready、已 Apply 或正在執行；每次仍須 CAS、來源／measurement、schema consent、ownership 與 readiness。
旗標由單一 backend MarshalJSON 計算，caller 偽造值無效。

精確 counter scale 的正常 mixed delta 支援限 managed-long；managed-wide 的該組合仍有 layout-specific restriction。
精確 scale mixed latest／mean、未知 shape/layout/dialect 不可偷偷轉 raw 或繞過限制。
restriction 省略 layout 表示全部 layout；未知 layout 拒絕，不默認 long。
SQLite／PostgreSQL schema capability 不等於 PostgreSQL 實機已驗證。

### Raw、選定結果與健康狀態不可混用

SQL selected delta 是增加量；provenance raw 是既有 measurement 輸入值。
若 confirmed exact scale 在 journal 前已轉 Decimal，raw 即縮放後 Decimal；`source_type` 仍忠實保留 PLC 原 uint64。
真例：原 counter `90071992547409930 → 90071992547410010`、scale `0.1`，selected SQL TEXT `8.0`，最新 raw `9007199254741001.0`。
不能由 raw 反算／捏造未保存的原 PLC 讀值，也不能把 selected result 當 source raw。
mean 沿既有 numeric engine：exact arithmetic 的不可整除結果保留 scale／rounding／Rounded 與 sum/count 證據；不將 delta/latest 的整數精確性泛稱所有平均值沒有捨入。
`result_present`、`raw_present`、row emitted、receipt、required skip、not_due 與目前 source health 分別呈現。
ended member 可以沒有 NextDue；partial 是 coverage 狀態，不是不存在的 member result status。

### Durability 的保證範圍

已 ACK／accepted 的 observations、row-or-skip、member progress、counter checkpoint、outbox 與 immutable effect 走既有 durable authority。
unknown commit 必須重建／查證原 stable identity，不清除 uncertain flag、不重送生成另一列。SQL 已 commit 回覆遺失用 receipt 解決。
這不等於未 ACK 採樣零遺失：現 production acquisition caller 的 aggregate error 不是 durable gap 或持久重送；外部 SQLite 不可用／未參與 gate 的寫入仍可能導致未 ACK 讀值缺失。
同 `*sql.DB` 的 Store `inTx` gate 只降低該入口內競爭，不涵蓋所有 direct Exec、不同 handle 或外部程序；不宣稱全 DB 無 BUSY。

## 規格對照

以下路徑相對 repo；`ws` = `internal/datalink/workspace`，`ss` = `internal/datalink/snapshot`，`gd` = `internal/datalink/groupdelivery`。
每列列出實作及可觀察斷言，不以「有測試檔」代替執行結果。執行狀態見後表。

| 原 requirement | 實作 authority | 場景／數值與測試對照 | 證據組 |
| --- | --- | --- | --- |
| Recording members own calculation intent | ws `write_group_member_calculation.go`、mixed canonical identity／CRUD | `write_group_mixed_crud_test.go` latest/5 與 delta/600、重送／foreign／same target；`write_group_mixed_active_identity_test.go`；API independent 真雙 writer 與 stop 隔離 | P13、P42 |
| Member windows align independently of activation | ss checked cadence／mixed scheduler；ws `write_group_mixed_apply_schedule.go` | schedule/history tests：12:00:07、5/15 秒、15/60→20/60、整界下一界、disabled history、overflow 零 mutation；API common boundary A20/22/24、B10/14、C1000→1008 得 22/12/8 | P21、P42、C43 |
| Mixed result columns are not legacy raw columns | dbtarget mixed EncodeRow／EncodeExactValue | API shared-long 歷史全欄值不變；fresh-wide 首次409、明確 Disable 後新 target；scaled delta SQL TEXT8.0；encoder scope／unsupported layout 拒絕 | P24、P42 |
| Missing results obey explicit completeness policy | ss mixed result／既有 counter engine | `mixed_completeness_test.go`、`mixed_gap_completeness_test.go`：required baseline skip、optional NULL、reset／gap span 不補 raw/0；真 durable restart 後下一窗有 delta | P23、P33、C43 |
| Counter policy and member state remain isolated | aggregation `CloseCounterWindow`；ss frozen policy/checkpoint copy | `mixed_counter_policy_test.go` signed10→7=-3、65530→4 wrap=10、reset tail；isolation／ownership tests、copy新record不帶baseline | P22、P34、C43 |
| Calculation changes invalidate configuration proofs | ws mixed schema/digest；counterproof；runtime frozen boundary | schema ledger oldproof+newCAS、method/cadence/required/identity/policy各別變更拒；stale source/map/measurement/scale 舊 sample 拒；正常 fresh consent 才 ready | P31、P32、P42 |
| Closure and counter progress are atomic and replayable | gd mixed Store／sqltx、runtime/grouppipeline recovery | `store_mixed*_test.go` 每寫入故障rollback、due-mask/bounds偽造、lost ACK；mixed receipt/replay tests exact payload、unknown inspect、不duplicate；disabled unknown worker重建 | P33、P34 |
| Legacy profiles keep their original meaning | migrations049–052、原 legacy normalization/encoder | full rows/typeof/index/trigger/view/FK、known schema、fault rollback/reentry、FK restore discard；legacy latest/mean/raw absence不變；API歷史row真讀回 | P12、P13、P33、P41、P42 |
| Production proof precedes editor support | canonical API／capability／runtimeview＋strict TS consumers | storage-only階段guard；後續真 Save→schema→Apply→runtime→Store→Sender→SQL＋receipt；forged capabilities 無效、editor仍false；actual JSON兩TSparser通過 | P12、P41、P42 |
| Different member cadences preserve sparse due semantics | ss typed progress；gd schedule驗證；兩種encoder | 15/60 wide NULL/not_due與long due-only；慢窗不被快窗消耗；12:00:07 first partial、exactedge歸下一窗、bounded catchup、replay不改due | P21、P33、P42、C43 |
| Persisted reduction and independent value roles | 原profile validator／legacy scheduler／encoder | `TestRecordingReductionCompatibility`、legacy reload/invalid role/cumulative mean拒；原 instantaneous counter_raw latest與 cumulative獨立；新mixed顯式契約不覆legacy | P13、P24、P32、P42 |

### 全部 25 個具名 scenario 的定位

| Scenario | 斷言所在（檔名相對上表 package） |
| --- | --- |
| One source participates in different records | ws `write_group_mixed_crud_test.go`、ss `mixed_counter_isolation_test.go`、API `router_mixed_recording_acceptance_independent_test.go` |
| Reusing a source does not transfer target ownership | ws `write_group_mixed_active_identity_test.go`、`write_group_mixed_crud_test.go` |
| Two averages and one usage share a row | API `router_mixed_recording_acceptance_matrix_test.go`，22/12/8 真 SQL |
| Never-applied mixed recording starts at the transaction time | ws `write_group_mixed_apply_schedule_test.go`；獨立 reviewer Save 後前進 clock 再 Apply |
| Active member cadence revisions share a safe cutover | ws `write_group_mixed_apply_schedule_test.go`、`write_group_mixed_apply_history_test.go` |
| Disabled history is not a first activation | ws `write_group_mixed_apply_history_test.go` |
| Unrepresentable revision boundary fails atomically | schedule tests；獨立 reviewer checked LCM overflow 後逐表逐值零變更 |
| A raw legacy table would change meaning | API `router_mixed_recording_acceptance_shared_test.go`、`_wide_history_test.go`；ws transition negatives |
| A required delta has no first baseline | ss `mixed_completeness_test.go`；runtime durable baseline skip/restart tests |
| An optional delta is uncertain or spans missing windows | ss `mixed_gap_completeness_test.go`、counter policy tests |
| A policy explains a decrease | ss `mixed_counter_policy_test.go`，signed及confirmed rollover |
| An old proof is submitted after changing method | ws `write_group_mixed_schema_ledger_test.go`、runtime frozen binding tests |
| A crash occurs around mixed closure | gd `store_mixed_failure_test.go`、`store_mixed_retention_test.go`；runtime mixed receipt/replay tests |
| Expansion fails while replacing the configuration table | datalink `migrator_*mixed*`／member expansion／record profile key safety tests |
| Existing profiles are reloaded after upgrade | migration full legacy row tests、`TestRecordingReductionCompatibility`、legacy encoder regressions |
| Mixed storage exists before cadence validation | 歷史1.2 checkpoint expand guard測試；正常draft於1.3完成後才開，不要求最終版永久拒絕 |
| Only a helper supports mixed methods | 階段guards＋最終API真SQL鏈；editor capability仍false，不把helper當support |
| Starting between clock boundaries | ss `assembler_mixed_windows_test.go`、ws schedule tests、API firstpartial矩陣 |
| Mixed cadence and delayed delivery | ss `assembler_mixed_windows_test.go`（5/15/60/600）；runtime wide receipt restart，due與receipt時間分離 |
| Fifteen-second means and a minute delta | API matrix wide/long及ss windows tests |
| Activation truncates the first aligned window | API matrix／scaled tests，aligned metadata與actualcoverage分離 |
| Legacy reload and explicit edit | workspace reduction compatibility、legacy migration/reload tests |
| Invalid combinations | workspace member applicability／counterconfirmation／legacy validation tests |
| One source in both profiles | legacy profile/counter tests；mixed未改其identity或counter_raw規則 |
| An explicit mixed recording is requested | API common22/12/8真鏈及legacy JSON對照 |

Delta specs 無 `##### Example` 區塊；scenario 內數值仍依上表測試追蹤。沒有以視覺細節為由排除任何條件／錯誤分支。

## 執行證據與版本

證據存於本次執行環境 `/tmp/`，不是部署機的必要檔案；正式測試本身已隨 checkpoint 保存。
`passed-prior` 只用於原來源／測試 hash 仍吻合或已由後續受影響 gate 覆蓋者，不把早期失敗改名為通過。
清理前 `/tmp/member43-evidence-identities.json` 對照14份歷次 manifests；清理後320個本案paths的目前hash在 `/tmp/member43-current-sha256.txt`。
與清理前比較：精確123個owned paths改變、其餘197個未變、範圍外變動0（`member43-final-identity-delta.json`）；123採新gate，其餘適用原有效證據，不能說原final35或舊320hash仍全部吻合。
較早檔案已變更時，採後續同路徑驗證；4.2最終35檔先核對發布一致，再以本次123檔選測／vet／lint覆蓋清理變動。
Spectra scope 初次 `c67099ac25fc2ead7481fb6d58971f03a19ec595` 因別案發布而失效；已重新捕捉／查核。最終記錄為 `/tmp/member43-scope-final.json`／`member43-scope-check-final.json`（320paths），獨立lint仍使用published22520a2＋本次exact123檔，不讀別案WIP。
`migrator.go` 有 preexisting-dirty 歸屬限制；共享生成Swagger也含別案已發布內容，不能把整個diff全算本案。

| 組別 | 範圍／有效終態 | 結果來源 |
| --- | --- | --- |
| P12/P13 | persistence、migration、nullablelegacy、CAS與cadence；相關normal/race/vet passed-prior | `member12-final-*`、`member13-review-packet.md`，後續migration051/052再驗 |
| P21/P22/P23 | snapshot typed progress、counter policy、missing；race1.085/1.086/1.116s passed-prior；最新normal另見C43 | `member24-fixed-related-race.log`、`member22-snapshot-race.log`、`member23-snapshot-race.log` |
| P24 | encoder scope/rounded/type/exactSQL，dbtarget race17.930s；後續4.2 partial修復再驗 | `member24-fixed-related-race.log`、P42 |
| P31/P32 | proof/readiness、frozen runtime、cutoff；相關race/vet passed-prior，不重跑舊274s workspace大包 | `member31-review-fix-race.log`、`member32-final-race.log`、`member32-cutoff-race.log` |
| P33 | atomic/due-mask/fence/restart、sameDBgate；Store99.064s/sqltx1.026s及pipeline recovery相關race/vet passed-prior | `member33a-review-packet.md`、`member33-gate-review-packet.md`、`member33-gate-isolated-race.log` |
| P34 | immutable receipt/replay、copy不帶baseline；runtime40.584s/workspace4.452s race，isolatedvet EXIT0 | `member34-isolated-race.log`、`member34-isolated-vet.log` |
| P41 | strict Go/TS source/result DTO、migration052實際query EXPLAIN、API/generatedschema；已發布hash與後續修復對照 | `member41-review-packet.md`、`member41-fix-schema-api-normal.log`、`member41-index-race.log` |
| P42 | publishedbase＋own overlay，真canonical intake/quota；workspace37.876s/API89.989s race EXIT0，四packagevet EXIT0 | `member42-final-canonical-ownership-race.log`、`member42-final-vet.log` |
| P42 encoder/runtime | dbtarget18.834s、grouppipeline137.317s；相關workspace125.882s/API58.335s race EXIT0；後續scope變更以P42新gate補 | `member42-related-race.log` |
| P42最後cap修正 | capability/caller forgery/layout/HTTP/Swagger：workspace0.029s、API0.391s，scopevet EXIT0；原兩caller獨立複驗0.260s | `member42-cap-all-contract-final.log`、`member42-cap-all-contract-vet.log`、`member42-independent/final35-old-contract-green.log` |
| P42前端 | 44tests、fulltsc、scopedESLint EXIT0；root獨立reviewclean；真producer兩個strict TS consumers均PASS | `member42-root-frontend-review.md`、`member42-consumer-review/*-result-final.json` |
| C43 | 最新22520a2隔離後完整snapshot0.022s／aggregation0.005s normal EXIT0，passed-current | `member43-reducer-normal.log` |
| C43 lint | 正式golangci-lint2.6.2、repo設定、8相關package、本案patch範圍；passed-current：修後正式8package／本案patch範圍0 issues，EXIT0；另9package fresh vet EXIT0 | `member43-fixed-lint-final.log`、`member43-lint-final-vet.log` |

### 保留的失敗，不混稱全部通過

- `member33-store-final-race.log`：600.054s aggregate timeout，原命令 FAIL；後續分片與gate修復各自有獨立證據，不把它改成 PASS。
- 初期 same-Store concurrent writer BUSY 在原baseline也重現；後續context-aware gate原neverbusy斷言未放寬，race57.008s/count3與99.064s通過。外部未ACK限制仍存在。
- `member42-source-reuse-and-canonical-normal.log` 整命令 FAIL（新增legacy測試fixture錯誤）；修正後正式窄race通過，API原該次也PASS，不冒充原命令全綠。
- `member42-independent/final34-old-contract-red.log` 舊caller還期待activation=false而FAIL；同步既定cap契約及嚴格optional-layout schema斷言後final35複驗PASS。
- 其他有意mutation反例、review發現的due mask／cutoff／scope／proof／parser失敗及測試fixture錯誤完整保留於各packet；獨立review最終未解產品finding為0。

## 未執行與移交

- NOT RUN：產品全套tests/build/lint，依使用者要求待七案完整功能最後統一執行；本文件只證相關範圍。
- NOT RUN：本案Browser端到端、實體PLC、正式DB、Windows、PostgreSQL實機。沒有以SQLite/helper或driver存在替代。
- editor/adoption：消費canonical member method/cadence/resultshape/source/resulttype/revisions、能力restriction與typedresult；不可把activation_supported當recordready，也不可把只讀mixed草稿降級寫回legacy或補60秒。
- 以「各點位分別計算」說明用途；明示新增獨立記錄、明確啟停，不能將另建說成歷史原地轉換。
- 4.3 checkbox／發布由root處理；本文件不自動archive。整案最終總測仍需後續完成。

## 本次收尾命令與程序狀態

```sh
# published 22520a2 + 排除其他 WIP 的 overlay
go test -overlay=/tmp/member43-published-overlay.json -p 1 \
  ./internal/datalink/snapshot ./internal/datalink/aggregation -count=1
# 純 published source 實體樹，正式 repo .golangci.yml，非 vet 冒充 lint
golangci-lint run -j 1 --timeout 15m --new-from-patch=/tmp/member43-owned.patch \
  ./internal/datalink/workspace ./internal/datalink/snapshot \
  ./internal/datalink/aggregation ./internal/datalink/runtime \
  ./internal/datalink/dbtarget ./internal/datalink/groupdelivery \
  ./internal/datalink/grouppipeline ./internal/api
spxa validate support-per-member-recording-calculation
bash scripts/check_file_lines.sh
```

選測完整regex（224個top-level＋terminal cleanup／unknown／Storegate）保存在 `/tmp/member43-tests-regex.txt`；命令 `go test -p 1 <上述8package＋sqltx> -run "$(cat /tmp/member43-tests-regex.txt)" -count=1`。
Fresh vet命令為 `go vet -p 1 <上述8package＋sqltx>`。最後lint使用更新patch `/tmp/member43-fixed-owned.patch`。

Go1.25.5；GOMAXPROCS=1、GOGC=50，沿同GOPATH/GOCACHE。官方 golangci-lint2.6.2 linux-amd64 archive SHA256：
`499c864b5fd9841c4fa8e80b5e2be30f73f085cf186f1b111ff81a2783b7de12`，與官方checksums相同。
首次lint的GOFLAGS overlay不能保護analyzer直接讀shared AST，遇別案wirebudget WIP三個typecheck而EXIT1；不是member產品finding，也不是PASS。以22MB純published實體樹重驗，沒有借用shared未發布依賴。

查核時任務13/14完成，4.3本身仍由root持有；依verify規則這是尚待關閉的task程序項，不是已證實缺少產品實作。11個requirement、25個scenario均有實作／斷言映射；未解產品review finding為0。原lint357 issues已按精確scope清理，修後0 issues、相關tests與fresh vet均EXIT0；程序項待root核證勾選與發布，不自行archive。

容量紀錄：07:03 UTC 可用空間降至1.3 GiB，依root要求只對owner session61904送Ctrl-C，工具確認exit130。未清cache或停止他人作業；等回收後僅補同lint，不重跑既有tests。

07:06–07:11 重跑終態：正式lint EXIT1，357 issues（180個test noctx、79個test goimports；其餘有context cleanup、typed degradation、ignored error及style待核）。這些不能以已有vet/review PASS抵銷；分類清單 `/tmp/member43-lint-issues.json`，修復已獲root授權並完成，原證據保留，不把357 FAIL改寫成原先PASS。

### Lint 清理的獨立批次（4.3，不新增功能）

357診斷逐行git blame：355行來自本案checkpoint；另2行是word69整合觸碰同一member API contract test。沒有把無關歷史整個package納入修改。
範圍精確123檔：93tests與30production，與當時其他owner WIP交集0。baseline與final manifests：`member43-lint-baseline-sha256.txt`、`member43-lint-final-sha256.txt`。
獨立靜態review確認hash `5632f051d077ab5033e731d28d09e0d82f6d74a96bc9b88fbdad821c48d48ab6`，結果clean，但不代替執行。

- tests：明確SQL Context、import/style與未使用參數；傳入transaction callback ctx，不借已取消的測試ctx執行terminalcleanup。原assertions、fault injection與driver fallback測試保留。
- production：wire literal重用、lexical命名／宣告、明確檢查Marshal錯誤；沒有新engine、schema、API型別或業務分支。Marshal錯誤為防禦性faultguard，不捏造可達behavior RED。
- `MemberNotDue` 前後都是typed `MemberStatus`；保持untyped的是 `ReasonNotDue`／`ReasonWindowArithmetic`，沒有變更值或assignability。
- sqltx terminalcleanup沿原獨立bounded context，typed unavailable/schema-not-ready沿原結構回覆；僅逐行註明lint理由，沒有關閉全域規則或放寬安全驗證。
- 修後normal選224個改檔top-level tests，再加原cleanup/unknown/panic/cancel/Storegate/baselineproof回歸：8個有選中測試的package全PASS，EXIT0，`member43-lint-related-normal.log`。aggregation此regex未選到測試，不算此命令的產品驗證；另有C43完整package PASS。
- 首輪修後lint剩1個test unlambda，`member43-fixed-lint.log`仍是FAIL；改成直接傳原callback後，原測試PASS0.162s，9package fresh vet EXIT0，正式8package lint **0 issues／EXIT0**（`member43-unlambda-normal.log`、`member43-lint-final-vet.log`、`member43-fixed-lint-final.log`）。其餘122檔不變；不因機械格式重跑全產品或舊大race。
