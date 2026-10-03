# B：交付結果記帳期限驗證

日期：2026-10-04；起點 `eb30a21f0e7aff7345ea04292bc8d286fe9c66db`。
沿用 `6676b2b0` 已完成的 context 起算修補。本次只補正式 Sender 的回歸、PostgreSQL 故障注入 fixture 與期限註解，不改預設 timeout、重試、去重、狀態、API、UI 或 schema。

## 行為與正式接線

`cmd/test_ui/service_wiring.go` 的 durable group pipeline 以 `groupdelivery.NewSender` 接正式 target resolver／SQLite store；`main.go` 啟動及有界停止同一 pipeline。不存在平行 queue 或 demo sender。

- 遠端 Resolve／InsertGroupRow 共用 attempt context，遵守 caller cancellation；預設 30 秒。
- Resolve 失敗的分支及 insert／commit 返回後，才建立新的 settlement context；預設 5 秒，以 WithoutCancel 保留提交後記帳機會。集中 `settle` 分類仍為 sql_committed／retrying／blocked／quarantined／unknown。
- 預設 lease 是 65 秒；不足 attempt＋settlement 的自訂 lease 沿用既有 30 秒 margin 調整。既有 fencing 不允許被取代的 owner 保存結果。
- worker graceful deadline 用完即取消 in-flight work，再最多等待既有 2 秒 forced grace；取消不能當作遠端未 commit 的證據。
- 本機收據交易真的失敗時保留 sending。有 destination receipt 才能安全恢復／去重，沒有 dedupe 則轉 unknown 並拒絕盲目重送。

## 回歸與刻意改壞檢查

測試直接使用正式 Sender／Store／InsertGroupRow，目標是可丟棄 SQLite 與 Docker PostgreSQL 16.14。新測試首次 PASS，依 Regression workflow 使用暫存 Go overlay，沒有宣稱現行修補前 RED，也沒有將 mutation 寫入 production 工作樹。

同步屏障在遠端已開始時通知測試，明確等候兩個 100 ms settlement budget，再釋放；attempt 為 5 秒。另以真實 driver 的 AfterCommit 屏障暫扣 commit 回覆，獨立 SELECT 先證明遠端已寫入，而 outbox 仍為 sending、local receipt 為 0。沒有用短 sleep 猜測遠端進度。

| mutation | 實際必要失敗 |
| --- | --- |
| Sender overlay 回到 `d2c22ef6`，settlement 在遠端前開始 | 兩 DB × success、transient、blocked、constraint rejection、ambiguous with/without receipt 共 12 cases 皆因 context deadline exceeded FAIL；慢 commit 回覆兩 cases 亦 FAIL |
| settlement 繼承 caller cancellation | 提交後取消的兩 DB cases 因 context canceled FAIL |
| local receipt 失敗被回報為成功 | 兩 DB × 有／無 receipt 共 4 cases，預期 ErrLocalReceiptFailed 卻得到 nil，FAIL |
| 強制保留 1 ms lease | 實際持久 claim 到期時間已過，live attempt 保護斷言 FAIL |

移除 overlay 後 focused／race PASS。兩個新測試檔各低於 300 行；PostgreSQL fixture 改用同一 lostcommit.WrapFaults，以便在真正 Commit 後取消，及重用既有回覆遺失注入。

## Requirement／Scenario 對照

唯一 delta requirement `Bounded isolated recovery and capacity` 對應 tasks 1.1–3.2；analyze 的名稱匹配警告 COV-1 不代表缺實作，以下逐情境記錄其測試。七個 scenarios 全部 covered，沒有 example 或 test-scope exclusion。

| scenario | 實作／有實際斷言的測試 |
| --- | --- |
| Outage and restart | 正式 pipeline／durable Store／Dispatcher；`TestProductionGroupOutageRecoveryOneWriterRestartAndIsolation`、`TestProductionGroupOutageRecoveryPostgresDestinationReceiptsAndRestart`；A 的正式 offline／kill recovery 證據仍保留 |
| Disk full | quota 在 durable ACK 前拒絕；`TestQuotaRejectsNewAckDiskFullIsRefusedAndKeepsAcceptedData` 真正限制 SQLite page count，原 journal/outbox 不刪除 |
| Poison row | `settle`＋partition gating；`TestBoundedRetryAndPoisonPartitionPoisonBlocksOnlyItsOwnPartition` 驗證 own successors 等候、另一 partition 繼續 |
| Overlapping workers | BeginDelivery／claim epoch／RecoverStaleClaims；`TestWorkerFencingAndShutdownStaleWorkerCannotWriteAfterTakeover`、`TwoWorkersNeverSendTheSameRowTwice`、`StopIsBoundedAndLeavesIntentRecoverable` |
| Slow successful destination attempt | `TestSenderDeadlineSlowRemoteResultsAreDurablySettled/*/success`、`TestSenderDeadlineSlowCommitResponseGetsAFreshSettlementBudget`，同 effect 的單筆 SQL／local receipt 與 payload digest |
| Slow failed destination attempt | `TestSenderDeadlineSlowRemoteResultsAreDurablySettled` 的兩 DB transient、target-blocked、row-rejected、兩 ambiguous branches，核對實際 outbox state／safe code／SQL rows／receipt absence |
| Caller cancelled after commit | `TestSenderDeadlineCancelledCallerAfterSlowCommitStillRecordsTheEffect`；`TestSenderDeadlineLocalReceiptFailureRetainsEvidenceAndRecoversSafely` 驗證真 SQLite trigger 拒絕 local receipt、有 receipt 收斂為單筆、無 dedupe unknown 不重送 |

## 實際檢查

環境：macOS darwin/arm64，Go 1.27.1（go.mod 1.25.5）、golangci-lint 2.14.0；owned PostgreSQL loopback 55432／gwtest。沒有正式 DB、PLC、LAN 或部署。cache 使用 writable temporary directory。

| 命令 | 結果／範圍 |
| --- | --- |
| `POSTGRES_DSN=<owned fixture> go test ./internal/datalink/groupdelivery -run '^TestSenderDeadline' -count=1` | PASS；兩 DB 慢成功／失敗／commit reply、取消／local receipt failure；overlay 移除後恢復 |
| `POSTGRES_DSN=<owned fixture> go test -race ./internal/datalink/groupdelivery -count=1` | PASS；包含 crash/lost response、quota、fencing、shutdown、retry defaults 與新期限回歸，無 race |
| `go test ./internal/datalink/groupdelivery -run '^TestWorkerFencingAndShutdown' -count=1` | PASS；失效 owner、hung attempt、forced shutdown 及恢復 |
| `POSTGRES_DSN=<owned fixture> go test -p 1 ./cmd/test_ui -run '^(TestProductionGroupOutageRecovery\|TestRevisionBoundBacklog)' -count=1` | PASS；正式 wiring、兩 DB outage/restart、原 destination/revision、legacy 排他性 |
| `POSTGRES_DSN= go test -p 1 ./... -count=1` | PASS；預設全套，條件式 PostgreSQL SKIP 不算 live PG pass；最終新增 helper 調整後再次執行 |
| `go vet ./...` | PASS |
| `golangci-lint run ./...` | PASS，0 issues；新測試 helper 的 context、結果命名、固定表名參數已依 gate 修正 |
| `git diff --check`、`make check-lines` | PASS |

獨立查詢 `pg_namespace` 的 `gw_wgd_%` 為 0，owned destination schemas 已清除。範圍外的舊 PostgreSQL partition migration 與 API schema verification 失敗，沿用 A 文件的已知限制；本案沒有把額外 live PostgreSQL 全套宣稱為通過。

`spectra-verify` 逐 task／七 scenario 核對 production wiring 與上述實際結果；`spectra-review` 檢查 correctness、efficiency、simplification/reuse、convention。只有期限註解／測試 helper 變更，MaxRetries=0、backoff、dedupe capability 與 transition 行為未改。

## 未執行與限制

本案沒有畫面修改。embedded UI 的主流程、錯誤／部分成功、鍵盤及 390/768/1440 witness 仍由後續涉及畫面的案件逐案完成。Windows／ARM、現場 PLC／SCADA 與部署均 NOT RUN。同步屏障證明期限生命週期，不是 100 ms 的 production 延遲或效能承諾。
