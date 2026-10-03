# A：寫入群組生命週期驗證

日期：2026-10-04；執行起點 `0feaf8a94dd4309f9bee23822c7d11d7e9a4ac10`。
本案沿用 `6676b2b0` 的草稿／恢復／鎖修補與 `0feaf8a9` 的快速重新啟用修補，沒有重寫已完成的 production 行為。
這次新增正式服務接線的回歸證據；六案按 A→B→C→D→E→F 交付，本文件不代表其他五案完成。

## 修改與正式接線

- `cmd/test_ui/group_pipeline_lifecycle_test.go`：真 SQLite／PostgreSQL 目的地、正式 workspace/Create/Update/Apply 與 frozen layout；草稿跨兩桶、重新開啟設定 DB、離線 cold intake 及獨立 SELECT。
- `cmd/test_ui/group_pipeline_lifecycle_crash_test.go`：ACK 後 Disable/Delete/replace，kill 實際子程序；重新開啟 durable 設定 DB，恢復只排空舊 revision，原目的地／effect identity、SQL 與 journal 均有斷言。另涵蓋快速重新啟用、描述缺失／版本不符／來源變動與 legacy 排他性。
- `cmd/test_ui/service_wiring.go` 將同一 WriteGroupService、TagService、DBTargetService、inspector 與 durable store 接至 pipeline；runtime SampleSink 接 pipeline，legacy writer 有 group ownership suppression 及 canonical guard。
- `cmd/test_ui/main.go` 啟動／停止 pipeline 並將它提供給 delivery API。新增測試沿用正式 `wireGatewayServices` 與同一 pipeline factory 依賴，只替換時鐘／tick cadence；目的地及 SQL 沒有 fake。
- 無 UI 原始碼修改，本案沒有用元件截圖代替 SQL／恢復測試。真 embedded UI 的四步流程、錯誤、部分成功、鍵盤及 390/768/1440 驗收仍由涉及畫面之後續案逐案執行。

Draft 是待 Apply 的設定；applied 是不可變的正在負責採集版本。Disabled 停止新 intake 資格但保留已接受資料；drain-only 只從 journal 封桶／交付，不重新接受樣本。刪除保留 tombstone、版本、ownership 與 backlog，沒有實體 purge。

## 回歸與 mutation 證據

新增測試首次即 PASS，所以不宣稱新測試在現行修補前 RED。依 TDD Regression workflow，以暫存 Go overlay 驗證斷言確實能攔截缺陷；production 工作樹從未套用 mutation。

| 暫時改壞的行為 | 指令 scope | 觀察到的必要失敗 |
| --- | --- | --- |
| lifecycle 回到 `6676b2b0` | `go test -overlay <reenable.json> ./cmd/test_ui -run '^TestProductionLifecycleRapidReenable' -count=1` | Disable→Apply 沿用相同 applied revision |
| Reconcile 不恢復 historical journal | `go test -overlay <historical.json> ./cmd/test_ui -run '^TestProductionLifecycleKilledHistoricalJournalDrainsOriginalRevision/sqlite' -count=1` | Disable/Delete/replace 各留下 2 筆未 consumed 樣本 |
| 草稿被當作退休狀態 | `go test -overlay <draft.json> ./cmd/test_ui -run '^TestProductionLifecycleDraftAndOfflineColdStartIntake/sqlite' -count=1` | 重啟後原 applied pressure member 不再接受 |
| 恢復強制重新 InspectTable | `go test -overlay <offline.json> ./cmd/test_ui -run '^TestProductionLifecycleDraftAndOfflineColdStartIntake/sqlite' -count=1` | 離線後不能恢復原 applied member intake |

移除 overlay 後，最終正式套件與 race 均 PASS。暫存 mutation JSON/source/log 只在本機 temporary directory，不列入提交。

## Requirement／Scenario 對照

| requirement / scenario | 實作 | 有實際斷言的測試 |
| --- | --- | --- |
| Safe basic group lifecycle：Rename and member edits | workspace lifecycle/version repository；pipeline 使用 applied snapshot | `TestGroupLifecycleApplySwitchesAtNextBucketWithoutChangingStableID`、`TestProductionLifecycleDraftAndOfflineColdStartIntake` |
| Delete with backlog | workspace tombstone/backlog guard；pipeline historical recovery；既有 sender | `TestGroupLifecycleDeleteWithOutboxPreservesPayloadAndReceipt`、`TestProductionLifecycleKilledHistoricalJournalDrainsOriginalRevision/*/delete` |
| Saved draft leaves running revision active | IntakeEligibility 不以 draft 退休；Reconcile 解原 applied version | 新正式兩 DB 草稿跨桶及 offline restart；既有 `TestDraftSaveKeepsAppliedIntakeAcrossBucketsAndRestart` |
| Concurrent cutoff and acquisition | GroupBoundary 的 until 檢查與 AcceptSample/SetUntil 共用 mutex | `TestGroupBoundarySetUntilAndAcceptSampleSynchronizeCutoff`，實際 runtime package `-race` PASS |
| Re-enable during or after retirement | disabled Apply 產生新 applied version，舊 revision 排空 | `TestProductionLifecycleRapidReenableJournalsAfterOldCutoff`、`TestGroupLifecycleReenableSameSemanticsCreatesFreshAppliedRevision` |
| Rapid disable/re-enable example | 同上 | 正式 fixture 的 10 秒 interval、t=12 Disable、t=15 Apply、t=21 雙 member journal，重啟後兩 revision 各一筆 outbox |
| Production durable intake：Crash after ACK | durable sample commit；frozen versions；journal recovery | 新正式兩 DB × 三狀態 kill matrix；既有 `TestCrashKillAfterAckLosesNothingAndDeliversOnce` |
| Crash during row transaction | Store.CloseBucket 的 payload/outbox/checkpoint 同一 SQL transaction | `TestAtomicRowOutboxCheckpointFailureAtEveryBoundaryRollsBackEverything`、`TestAtomicRowOutboxCheckpointBoundaryCommitFailureKeepsBucketOpenForReplay` |
| SQL failure | sender 保留 durable intent，receipt/fencing | 正式兩 DB outage/recovery suite、`TestProductionLifecycleKilledHistoricalJournalDrainsOriginalRevision` 與 groupdelivery suite |
| Destination offline during cold start | Apply 持久 frozen layout；buildFrozen 不查遠端 | 新正式兩 DB draft/cold intake、`TestVerifiedAppliedLayoutRestoresIntakeWhileRemoteIsOffline` |
| Restart drains historical accepted input | OpenJournals 按 revision 發現；drainOnly 不 AcceptSample | 正式 SQLite/PostgreSQL Disable/Delete/replace kill matrix，舊 2 樣本全 consumed、單 effect、單 SQL row |
| Recovery metadata cannot be proven | frozen key validation；legacy source/mapping revisions 驗證；typed blocked | `TestProductionLifecycleUnprovableRecoveryKeepsJournalAndBlocksLegacy`、`TestMissingLegacyRecoveryDescriptorRetainsAcceptedJournalWhenOffline` |

全部 11 個 scenario 與 1 個 example 都有測試；沒有 test-scope exclusion。Migration replay／交易 rollback／receipt capability 由 `write_group_runtime_layout_test.go` 四項回歸涵蓋；connector 改 identity 的排他及原 destination 留存由 `TestRevisionBoundBacklogEndpointEditAndDisableNeverRetargetOldRows` 與 pipeline recovery 回歸涵蓋。

## 本次實際執行

環境：macOS darwin/arm64，Go 1.27.1（go.mod 1.25.5），golangci-lint 2.14.0，owned Docker `postgres:16`／PostgreSQL 16.14，loopback `55432/gwtest`。未使用正式 DB、PLC 或 LAN。SQLite destination 與 gateway configuration/journal 各自獨立。

| 命令 | 結果／範圍 |
| --- | --- |
| `POSTGRES_DSN=<owned fixture> go test ./cmd/test_ui -run '^TestProductionLifecycle' -count=1` | PASS；兩 DB 草稿／cold intake、六種 kill/recovery、快速啟用；helper 只由 parent spawn，正常 suite 的 helper SKIP 不是案例缺跑 |
| `go test ./cmd/test_ui -run '^TestProductionLifecycleUnprovable' -count=1` | PASS；缺描述離線、revision mismatch、source 改變，journal 保留且 legacy mutation 被拒 |
| `POSTGRES_DSN=<owned fixture> go test -race -p 2 ./internal/datalink/workspace ./internal/datalink/grouppipeline ./internal/datalink/runtime ./cmd/test_ui -count=1` | PASS；實際四個 package，包含兩 DB 新回歸，沒有 race |
| `POSTGRES_DSN= go test -p 1 ./... -count=1` | PASS；預設環境全套最低檢查，條件式 PostgreSQL tests 的 SKIP 不算 PostgreSQL 實測 |
| `go vet ./...` | PASS；最終來源 |
| `golangci-lint run ./...` | PASS，0 issues；Go/lint cache 使用 writable temporary directory |
| `git diff --check`、`make check-lines` | PASS；新增兩個測試各低於 300 行 |

Review 找到 crash helper 在 readiness 失敗時只有 Kill、沒有 Wait 的資源回收問題；已補上條件式 Kill/Wait。修後 `POSTGRES_DSN=<owned fixture> go test -race ./cmd/test_ui -run '^TestProductionLifecycle' -count=1`、`go vet ./...`、`golangci-lint run ./...` 再次 PASS。獨立查核 `pg_namespace` 的 `gw_lifecycle_%` 為 0；本案 owned schemas 已清除。

`spectra-verify`：9/9 tasks、2 requirements、11 scenarios、1 example，沒有未覆蓋情境或 design 不符。`spectra-review`：以 `touched_tracking` 的 `6676b2b0→0feaf8a9` 加本次差異審查；檢查 correctness、efficiency、simplification/reuse、convention，修後無剩餘 Critical/Warning。既有 `6676b2b0` 恢復實作另外逐檔與 production wiring 核對；兩個 snapshot check 都通過。全套 PostgreSQL 的下列既有失敗及平台限制不因審查通過而消失。

### 額外全套 PostgreSQL 的既有失敗

`POSTGRES_DSN=<owned fixture> go test -p 4 ./... -count=1` FAIL：

1. `internal/api/TestStudioV2RecordingRoutes_PostgresApplyCreatesAndVerifiesTheManagedSchema`，schema operation 回 `unknown/verification_unavailable`；單獨重跑同樣 FAIL。該原始碼／測試與本案起點相同，與 A 新測試沒有共同 namespace；不得算成 A 的恢復能力失敗或 PostgreSQL 全套通過。
2. `internal/datalink/storage/TestPostgresPartitionMigration`，既有 migration 缺 `create_timeseries_partition_monthly(timestamp)`，`timeseries` 不存在。這是舊 storage migration，A 不修改其行為。
3. `modbusshare/TestService_Status_AfterStartAndStop` 在 package 平行時遇 5020 bind 衝突；預設全套改為 `-p 1` PASS，沒有修改 Share。

上述兩個 PostgreSQL 非 A 套件問題保持已知失敗，不新增 A 驗收需求或順手改其他模組；後續 D/F 如觸及相同 production path，須用該案範圍重新判斷必要修復，不能沿用本文件宣稱其全套已過。

## 未執行與限制

- 真 embedded UI、Windows/ARM、PLC/SCADA/現場與部署：NOT RUN；本案沒有畫面實作，後續涉及畫面之案仍須執行使用者指定矩陣。
- 子程序實際跑正式服務接線與 pipeline；不是整台 gateway binary 的嵌入網頁／設備採集 witness，兩者不得混稱。
- runtime frozen layout 由本地設定 DB 的既有信任邊界保護，不新增簽章系統、跨程序協調或第二套 queue。
