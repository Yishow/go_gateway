# Counter replay：升版、crash/resend 與安全回退驗收

日期：2026-10-27。對應 change：`freeze-counter-replay-semantics`，tasks 4.3–4.4。
本紀錄是現有開發環境內的 SQLite／模擬來源驗證，不是部署紀錄。

## 驗證邊界

- 使用既有 Go 1.25.5、離線依賴及可丟棄的 local/destination SQLite。
- `lostcommit` 先執行真實 SQLite commit，再模擬 commit 回覆遺失。
  pre-commit fault 則在 SQL 寫入或 config transaction commit 前拒絕。
- 本次 task 4.3 新增 fixtures 的「restart」是關閉／重開 local DB handle，重建 Store、Pipeline、assembler；
  不是 OS kill、斷電、檔案系統損壞或真實 PLC 重啟測試。
- Legacy fixture 在 new reader 啟動前，以固定舊 wire JSON、獨立 SHA-256
  與原 SQL transaction 建立舊格式資料；不是執行舊 binary。
  不修改已接受 payload/digest，也不讓新 intake 繞過 proof-reference gate。
- Sender 寫入真實 destination SQLite，讀回 row、receipt、unit、epoch、delta。
  測試以單一 delivery partition worker 避免把 SQLite 寫入競爭當作 crash 語意。
- PostgreSQL／其他外部 DB fixture 未提供，相關 skip 是 **NOT RUN**。
  沒有新增環境、連線正式 DB、部署或執行 schema 降版。

## 本次新增的證據

### Proof apply transaction

`internal/datalink/workspace/counter_binding_apply_commit_test.go`

`TestCounterBindingProofApplyCommitCrashBoundaries` 的兩個案例：

1. 真正 commit 前失敗：沒有成功 Apply 回應，workspace、applied revision、
   runtime proof 與 journal 均未半成功。
2. 真正 commit 後回覆遺失：回傳不確定錯誤，不宣稱成功；DB 中保存完整同一版本。
   呼叫端重新讀取目前 revision 後重試，只有一份 immutable runtime descriptor。

此 proof Apply 測試使用同一 DB handle 的 SQL reread／retry；沒有關閉重開 DB。
下節 sample／closure matrix 才實際重開 local DB 並重建 runtime memory。

既有 `TestCounterBindingApplyRuntimeInsertFailureIsAtomic` 另驗證 proof INSERT
失敗不啟動 ownership；CAS／inspection-edit 測試保護 reviewed facts 的一致性。

### Sample／closure／codec 邊界

`internal/datalink/grouppipeline/counter_binding_crash_test.go`

`TestCounterBindingCrashResendAcrossJournalFormats` 涵蓋三種來源：

- versioned native uint64；
- versioned uint64 exact scale `0.1`；
- 原 snapshot 具有完整 ExactRecordingSource、但 journal 無新 reference 的 legacy exact。

每種來源驗證五個邊界：sample 寫入前拒絕、sample commit 後回覆遺失、
已 ACK 後丟棄 runtime memory、closure transaction 中途失敗、closure commit 後回覆遺失。

共同斷言：

- 未 commit 不 ACK；失敗不產生 checkpoint／baseline 或消耗原 pending input。
- 已 commit 的 c-1 重啟重送只辨識一次；legacy 原始 JSON 與 d-1 保持不變。
- Legacy 兩筆 accepted inputs 都建立 durable association，而非只驗一筆後信任整個 revision。
- 100→105 保持 `kWh / epoch-1 / delta 5`；scaled 路徑為 `100.0→105.0 / 5.0`。
- 每個 scenario 只有四筆 finalized effects／destination rows／local receipts
  （兩個 windows、兩個 members）；再次重啟／Tick／delivery 不新增或改寫 digest。
- Current-definition provider 呼叫數為零；原 runtime descriptor bytes 不變。

既有 `TestProfileClosureCrashRecovery` 保護同一 process 的 closure retry；
`TestProfileCommitAckLost` 保護 destination commit 已完成但 ACK 遺失的去重。
既有 association failure/lost-commit、codec tamper/version 與 quarantine 案例納入選測，
不重複新增替代實作。

## TDD／mutation 證據

此 task 的既有 production transaction 行為已成立，採 regression＋真實 mutation kill，
沒有為了產生 RED 增加多餘 production 修改。

以下五個 `/tmp` Go overlay 均以行為斷言失敗，而非編譯失敗或零測試：

| Mutation | 被抓到的違規 |
| --- | --- |
| 忽略 proof commit error | 不確定／未 commit 被回報成功 Apply |
| 忽略 AppendSample 持久化錯誤 | 未持久化卻 ACK |
| closure error 後仍前進 | fault 未回報、memory 提前提交 |
| legacy/new reference 強制相等 | 原 c-1 重送被錯誤判為 conflict |
| closure retry 改寫 finalized digest | 原 effect digest 變為 `changed` |

證據目錄：`/tmp/counter-crash-mutations/`；只修改 overlay 檔，repo production source 未被 mutation。

初版 fixture 曾錯把 long SQL 的 epoch 當作 provenance 欄位、假設一次 Tick 只關一個
window，並嘗試用新 AppendSample 建立無 reference 的舊資料。這些是 fixture 假設錯誤，
不是 production RED；修正後保留原 gate，legacy SQL 只用於升版前的 fixture。

## Current-definition fallback 查核

- `grouppipeline/reconcile.go`：counter 缺原 runtime snapshot 時回傳
  `counter-replay-proof-missing`，不重新讀 definition 製造歷史 facts。
- `buildFrozenWithCause` 將原 `CounterBindings` 交給 runtime。
- `runtime/counter_binding_frozen.go` 使用 `counterproof.ResolveSnapshot`；
  `CounterSourceFromDefinition` 的輸入是 proof 組成的 value，不是 repository 查詢。
- Store 的 reference 驗證與 legacy association 共用同一 immutable resolver。
- `Counters` 相容介面及 production wiring 仍存在，但 pipeline/runtime replay
  沒有呼叫 current-definition provider 的 production call site；未為美化差異刪掉介面。
- Public measurement PUT→restart→SQLite readback 的整合證據在
  `internal/api/router_studio_v2_counter_replay_test.go`；task 4.2 已驗證新 epoch 的獨立 NULL baseline。

## Legacy 分類與升版驗收

升版的可證明分支只接受原 immutable descriptor 的完整事實：full member identity、
revisions、epoch、unit、type、scale、policy。附加 association/audit 保存原接受來源；
不改 journal／runtime bytes，不以目前名稱、definition 或部分 checkpoint 補猜。

| 狀況 | 行為／reason |
| --- | --- |
| 完整 legacy ExactRecordingSource | 原 bytes/digest 不變，附加可追溯 association，再 replay |
| 原證據缺失 | affected revision：`counter-replay-proof-missing` |
| 原 scope、digest、member 或 audit 矛盾 | `counter-replay-proof-conflict` |
| 原 proof version 不支援 | `counter-replay-proof-unsupported` |

Quarantine 保留未 finalized inputs、snapshot、checkpoint、baseline、association 與 frozen destination。
`pending_samples` 是受影響的 accepted input 數，不是 output row 數。
健康 scope 可繼續；已 finalized 的合法 output 依原 bytes/destination retry，
SQL delivery 成功不等於 quarantine inputs 已 delivered，也不會清除 quarantine。
僅成功驗證同一原 revision 的可信證據可解除；套用新 draft／current definition 不會修復舊資料。

## 安全升版／回退操作界線

1. 切換不相容 reader 前先停止受影響 intake，保留完整 durable store 與交易一致性。
2. 使用能理解已保存 proof、journal 與 association 版本的 binary 驗證原證據；
   對無法證明的 revision 保持 quarantine，核對 scope、pending count 與安全 reason。
3. 修復只可恢復該 revision 的可信原證據，或聯絡維護者；不得用目前 definition
   重標歷史、reset baseline/checkpoint、丟棄 journal 或改寫 digest 來取得成功狀態。
4. 已寫入新格式後，不可讓不理解格式的舊 binary 重播。
   不能假設舊 reader 一定會拒絕未知欄位；必要回退仍須停止 affected intake，
   保留 store，改用理解其格式的相容 binary。
5. 不執行會刪除 proof/reference/audit 的 schema rollback。
   新明確 applied epoch 是新 owner，不是消除舊 quarantine 的工具。

以上是安全操作指引；本次沒有執行真實升版／降版部署或現場事故演練。
原有 finalized consumed journal 可按既有 GC 回收；此事不授權回收 pending inputs。

## 命令與結果

環境：`GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off`，既有 Go 1.25.5；
`GOCACHE=/workspace/scratch/b2a0a232be39/go-build-cache`，
`GOTMPDIR=/workspace/scratch/b2a0a232be39/go-test-tmp`，`GOPATH=/tmp/gateway-gopath`。
下列 `go` 指 `/tmp/gateway-go/go/bin/go`。

```sh
go test -mod=readonly -json ./internal/datalink/workspace ./internal/datalink/grouppipeline \
  -run '^TestCounterBinding(ProofApplyCommitCrashBoundaries|CrashResendAcrossJournalFormats)$' -count=3

go test -mod=readonly -json ./internal/datalink/workspace ./internal/datalink/grouppipeline \
  ./internal/datalink/runtime ./internal/datalink/groupdelivery ./internal/datalink/snapshot \
  ./internal/datalink/counterproof \
  -run '^Test(CounterBinding|CounterProof|CounterReplay|CounterFinalized|LegacyCounterProof|ProfileClosureCrashRecovery|ProfileCommitAckLost|DurableSampleAckCrash|SampleCodec|SelectorCounterProof|FrozenExact)' -count=1

go test -mod=readonly -race -json ./internal/datalink/workspace ./internal/datalink/grouppipeline \
  -run '^TestCounterBinding(ProofApplyCommitCrashBoundaries|CrashResendAcrossJournalFormats)$' -count=1

go test -mod=readonly -json ./internal/datalink/counterproof -count=1

go vet -mod=readonly ./internal/datalink/workspace ./internal/datalink/grouppipeline
```

- 新 regression 三次：57 PASS events，零失敗／skip；`/tmp/counter-crash-regression-green.jsonl`。
- 最終選測：414 PASS events、6 PostgreSQL fixture SKIP、零失敗；`/tmp/counter-crash-final.jsonl`。
  counterproof package 此 selector 未選中測試，不計入測試 PASS 數；另列完整 package run 補測其 storage projection。
- Counterproof package 全測：9 PASS events、零失敗／skip；`/tmp/counter-crash-counterproof-final.jsonl`。
- Race：19 PASS events、零失敗／skip；`/tmp/counter-crash-race-final.jsonl`。
- Vet 通過：`/tmp/counter-crash-vet.log`。
- 五個 mutation 各自選到實際案例且失敗；不是 production failures。
- 不宣稱全 workspace suite 綠燈：先前在原 HEAD 亦重現的 concurrent EnsureBasic
  `SQLITE_BUSY` baseline 未修復、未隱藏，task 4.3 未重跑該 full-suite gate；task 4.4 的實跑結果另列於下。
- `/tmp` linker 空間限制已在先前 task 透過同一 workspace 的 cache/temp relocation 解決；
  不構成新環境，也沒有刪除別人的 cache 或證據。

## Spec scenario → 實際測試入口

對照 `openspec/changes/freeze-counter-replay-semantics/specs/durable-recording-delivery/spec.md`。
下列測試名是驗收入口；是否通過以本節後的 task 4.4 實跑結果為準，不能把名稱當結果。

| Scenario | Named test（package） |
| --- | --- |
| Definition changes before replay | `TestStudioV2CounterReplayPublicMeasurementEditPreservesAcceptedEpoch`（api）；`TestCounterBindingReplayIgnoresDefinitionEditDeleteAndOutage`（grouppipeline） |
| Definition is unavailable after durable acceptance | `TestCounterBindingReplayIgnoresDefinitionEditDeleteAndOutage`（grouppipeline）；`TestFrozenExactRecordingSourceDoesNotRefetchDefinition`（grouppipeline） |
| Concurrent definition edit and apply | `TestCounterBindingApplyRejectsDefinitionEditDuringInspection`、`TestCounterBindingApplyCASAfterReview`、`TestCounterBindingApplyRuntimeInsertFailureIsAtomic`（workspace） |
| Incoming claims name another epoch | `TestCounterBindingClaimsRefuseWithoutAdvancingAcceptedState`、`TestCounterBindingClaimsCopyAppliedIdentityNotCallerConfig`（runtime） |
| Proven legacy exact source | `TestLegacyCounterProofAssociationKeepsOriginalAcceptance`、`TestLegacyCounterProofAssociationChecksEveryOriginalSample`（groupdelivery）；`TestCounterBindingLegacyExactJournalUsesOriginalBytes`（grouppipeline） |
| Legacy proof cannot be established | `TestCounterReplayQuarantineReportsExactPendingInputsAndRetainsEvidence`、`TestCounterReplayQuarantinePreservesOriginalLegacyJournalAndAssociationAudit`（grouppipeline） |
| Quarantine does not rewrite finalized effects | `TestCounterFinalizedDeliveryRetainsMixedBacklogAndIndependentReceipts`、`TestCounterFinalizedLegacyBacklogNeedsNoUnfinalizedProof`（grouppipeline） |
| New epoch after old backlog | `TestCounterEpochReplayKeepsOldBacklogAndNewOwnerIndependent`（grouppipeline）；`TestCounterBindingExplicitEpochApplyCreatesNewProofAndRetainsAcceptedOwner`（workspace）；上述 public API 測試 |
| Crash and resend across codec versions | `TestCounterBindingProofApplyCommitCrashBoundaries`（workspace）；`TestCounterBindingCrashResendAcrossJournalFormats`（grouppipeline）；`TestLegacyCounterProofAssociationLostCommitResponseRetriesWithoutRewriting`（groupdelivery） |

公開 API fixture 不是只驗 HTTP 200：它走 measurement PUT、原 pending journal 不變、
重開 local SQLite／重建 runtime，再讀 destination SQL 的 `kWh / epoch-1 / delta 5`，
以及新明確 applied epoch 的 NULL delta。Mixed backlog／epoch／crash fixtures 另外核對
frozen outbox bytes/digest、destination rows 與 local receipt，不把記憶體狀態當 SQL 證據。

## Task 4.4 最終 gate（2026-10-27）

沿用上列 toolchain、offline module、同 workspace cache/temp，`-count=1` 不使用 cached result。
完整原始 logs 保留於 `/workspace/scratch/b2a0a232be39/counter-final-verification/`。

```sh
packages='./internal/datalink/workspace ./internal/datalink/groupdelivery ./internal/datalink/grouppipeline ./internal/datalink/runtime ./internal/api ./internal/datalink/counterproof ./internal/datalink/snapshot ./internal/datalink/schema'
go test -mod=readonly -json $packages -count=1
go test -mod=readonly -race -json $packages -count=1
go vet -mod=readonly $packages
git diff --check
make check-lines
```

### 首次完整執行與有界修復

`full.jsonl`：2520 PASS、4 FAIL、230 SKIP test/subtest terminal events；
其中 1038 PASS、2 FAIL、34 SKIP 是 top-level tests。
四個 FAIL events 只來自兩個 API leaf failures（含其父層 fail），不是四項獨立缺陷。

- `TestMeanLongRecording` 的 pre-Apply 不完整 epoch 診斷，被新 unsupported-type
  early return 遮住。HEAD 原行為會回報 `counter-measurement-unconfirmed`。
  `write_group_counter_readiness.go` 僅把 ConfirmCounter 診斷移到 type abort 之前；
  unknown／missing／string／bool type 仍不得產生 proof。
  `diagnostics-unproven-type-red.jsonl` 先證明原 epoch/non-counter 原因被遮蔽，再修復。
- 同一 mean test 原 post-Apply 期待 nil provider／目前失效 definition 讓有效原 proof
  停機，已被此 change 的 immutable 合約取代。保留全部 pre-Apply refusal 與 mean SQL
  assertions，改為在目前 definition 仍失效時完成 intake、restart 與 SQL readback。
- `TestCumulativeLongGroupRolloverFromSavedPolicy` 的舊子案例只改目前 policy、未新 Apply，
  卻期待 baseline reset。新名 `a current policy edit cannot reset the applied baseline after restart`
  斷言原 epoch／policy 100、`9 - 4 = 5`，保留 descriptor、已 closed rows/effect bytes/digests。
  真正新 applied owner 的 NULL baseline 仍由 named public API／eight-way epoch tests 驗證。
- `TestCounterBindingLegacyUnscaledJournalCannotBorrowCurrentDefinition` 僅把 editable provider
  的 revision 明列 `def-2` 以對應規格 example；原 journal／immutable bytes 與零呼叫斷言不變。

`repair-focused.jsonl`：29 PASS、7 PostgreSQL fixture SKIP、零 FAIL。
首輪 full 失敗不是因為 known EnsureBasic SQLITE_BUSY；該測試在首輪確實通過。
不以重跑取代修復：第二輪完整 gate 是在上述有界 source/test 變更後重新驗證。

### 最終結果

- Normal full（`full-final.jsonl`）：八個 package 均 PASS；2526 PASS／0 FAIL／230 SKIP events。
  Top-level：1040 PASS／0 FAIL／34 SKIP；leaf：2197 PASS／0 FAIL／230 SKIP。
  Parent/subtest events 不應相加宣稱為獨立情境數。
- 230 SKIP 細分：179 PostgreSQL fixture、38 MySQL/MariaDB fixture、10 SQL Server fixture
  未提供；2 SQL Server operation gate 關閉；1 僅供 subprocess 呼叫的 helper。
  外部資料庫皆 **NOT RUN**，skip 不計作成功。
- 本次 full 的 `TestWriteGroupServiceEnsureBasicManagedConcurrentAcrossServicesConverges`
  通過；先前原 HEAD 重現的 SQLITE_BUSY baseline 仍記錄，未修復也未用 rerun 隱藏。
- Vet：八個 package 通過，`vet-final.log` 為空，exit 0。
- 原 full race（`race-final.jsonl`）未通過：1594 PASS／125 SKIP／11 FAIL terminal
  test/subtest events（四個 assertion leaf failures）；另有四個 package 在預設 10 分鐘到期，
  因此不能視為完整 coverage。runtime／counterproof／snapshot／schema 已完整 PASS。
  未出現 race detector warning，不等於 full race PASS。
- Timeout stacks 顯示當時 active tests 才運行 0–2 秒，處於 runnable SQLite fixture migration
  或純 proof marshal；沒有顯示卡在新 managed gate 的鎖。只對 workspace／groupdelivery／
  grouppipeline／api 各做一次 `-race -timeout=30m` 完整重驗，其他已完成 package 不重跑。
  `race-adjusted.jsonl` 於 UTC 12:04:35 完成（22m24s），沒有再次 timeout。
- Adjusted race：2090 PASS／5 FAIL／230 SKIP events；top-level 811 PASS／2 FAIL／34 SKIP；
  leaf 1798 PASS／2 FAIL／230 SKIP。Workspace 694 PASS；Store 281 PASS／19 SKIP；
  Pipeline 398 PASS／3 FAIL／20 SKIP；API 717 PASS／2 FAIL／191 SKIP。
  兩個 leaf failures 僅為下述 timestamp substring／coexistence fixture；原 full race 仍非 PASS。
  兩輪均無 DATA RACE warning，但不宣稱所有 race assertions 通過。
- Coexistence 的四個 nested cases 在 adjusted run 因 setup fail 未執行，不計 PASS；
  同 source 的原 race run 四個皆 PASS。這是分次結果披露，不是合併拼成全綠。

### Full race 已觀察到的 fixture 限制（原失敗仍保留）

- `TestCircuitRecordingIntegration/shared_frozen_reconstruction/sqlite`：連續 Apply
  剛好跨 UTC 分鐘邊界，既有 fixture 要求全部 group 的 next-boundary effective time 相同，
  實際為 11:32／11:33。這是 unchanged wall-clock setup assertion，不是 race detector report。
- `TestDeviceRecordingIntegration/reviewed_long_history_pending/sqlite`：期待 pending=4，
  實際 pending=2、sql_committed=3。既有 Worker 在 Start 後先立即執行第一個 delivery cycle，
  之後才等 interval；`restart(time.Hour)` 並非暫停 delivery 的同步保證。
  `worker.go` 與此 fixture 皆和 HEAD 相同；這是 source 證明的 scheduling hazard，
  **不是**本次已執行 HEAD 並重現的 baseline 證據。

- `TestReviewedWideMeanTransition/a_rollback_is_a_reviewed_return_that_deletes_neither_history`：
  預期 7 committed effects，實際 6。獨立 review 用 diagnostic-only `anchor=wideEffective`
  重現原 race 的分鐘碰撞（目前 code 3/3 與 HEAD 913cead overlay 均為同樣六筆）；舊 history
  已 finalized 相同 bucket，fixture 倒退時鐘卻期待再產生 old cut row。
  原 HTTP 時序亦顯示 Create 在 11:38:59、Apply 在 11:39:00；不是新增 replay 遺失資料。
  證據：`/tmp/counter-review-rollback/{head-collision.jsonl,head-overlay.json,collision-only.jsonl}`。
- `TestWideCounterParity/sqlite/gas_is_an_exact_integer_difference`：raw／delta／status 一致，
  只有 independently applied long/wide 的絕對 delta bounds 差 600 秒；未固定共同 Apply clock
  的既有 fixture 跨過 10 分鐘邊界。這是 source diagnosis，不是 HEAD execution reproduction。

上述失敗不以修改 production cadence、削弱 assertion、skip 或反覆重跑隱藏。
完整 race 將保留原 FAIL；獲准的一次性 exact-case focused race 僅查明目前行為，
不會使原 full race 變為 PASS。此 change 不接管下一個 cutover change 的 fixture hardening。

兩個已核准 exact-case isolated race 都只執行一次並通過（各三個 test/subtest events）：
reconstruction：UTC 11:41:54→11:42:05，elapsed 11s；pending：11:42:05→11:42:11，6s。
Logs：`race-isolated-reconstruction.jsonl`、`race-isolated-pending.jsonl`。
所有 fixture、production clock/cadence 均未為此重驗而改寫；原 full failure 仍保留。

Adjusted race 另記錄 `TestAverageDeliveryEvidence/mean_long/...` 的既有字串 assertion：
禁止整份 JSON 含 `15.5`，但它命中的是 `committed_at` 的秒數 `15.520737158`／`15.511611668`，
不是 value/sum 欄位。保留原失敗；未因 timestamp collision 改寫產品或弱化資料遮蔽契約。

Adjusted race 的 `TestProfileLegacyCoexistence/sqlite` 另有 `late-after-close`：既有 fixture
先以所有 group 的 latest effective+5s 啟動 background Tick，再注入各自較早 effective+2s 的 sample。
15s／60s windows 與非同步注入可讓 earlier window 先關閉；記錄此 source diagnosis，未稱已執行 HEAD 重現。

最終 `git diff --check`、untracked scope whitespace、`make check-lines` 與 gofmt check 結果見
`diff-check-final.log`、`untracked-check-final.log`、`check-lines-final.log`、`gofmt-final.log`；
精確本 task 六檔清單與 hash：`task44-touched.txt`、`task44-final.sha256`。
