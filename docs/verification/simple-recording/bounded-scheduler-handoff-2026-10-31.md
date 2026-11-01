# Bounded scheduled acquisition：task 3.2 交接（2026-10-31）

## 歸屬與範圍

- Change：`integrate-protocol-aware-acquisition-batching`；只交付 task 3.2，未自行 done、commit 或 push。[tasks](../../../openspec/changes/integrate-protocol-aware-acquisition-batching/tasks.md)、[design](../../../openspec/changes/integrate-protocol-aware-acquisition-batching/design.md)。
- 本批實作／選測完成並凍結：26 source/test paths（14 production、12 tests），加本文件共27 paths；CLI checkbox由root核對後決定。
- 基準是 task 3.1 已完成的 frozen bytes，不是 Git HEAD 的舊 scalar scheduler；每個既有檔的差異基準保存在 `bounded32-evidence/baseline`（scheduler.go 另存 `scheduler.before.go`）。
- Production 仍是原 scheduler、immutable planner、device gate、ManagedConnection 與原 wire-budget send fence；沒有新 planner、retry engine、HTTP routes、UI、SQL schema 或 durable identity。
- Manual 原 device gate 保留；manual shared planner/executor/capacity parity 留 task 3.3。公開 preview/diagnostic 呈現留 4.1；100-scan/CPU/memory 矩陣留 4.2。

## 接線與有界條件

- 原 group ticker 只捕捉及 enqueue due snapshots，不再等 `pollGroup` 全設備完成。`pollGroup`/`pollDevicePoints` 只保留既有內部直接測試入口。
- 每個 group 一 pending container，每 `(device, group)` 至多一份最新 admitted snapshot；global pending capacity 計未 dispatch tuples。同 tuple replacement 不增加深度、不改 running snapshot；滿時拒新 tuple 並計數，不 evict 他 group。
- Sorted devices 加每 ticker 的起始位置輪替，避免容量一時永遠優先相同 device。Admitted ready devices round-robin；同 device groups 保先到穩定順序。
- 單 dispatcher、原 group producers、每 device 至多一個 waiting/executing scan goroutine。Pending entries 不各自建 goroutine。
- Enqueue 在短 scheduler read lock 內重驗 captured group/run/install/config，再做 queue bookkeeping；鎖序只有 scheduler→queue。Queue 不回取 scheduler，也不執行 planning、SQL、wire I/O、observer 或 permit blocking。
- Global active permit 在 device gate 之後取得，涵 connection/read active phase；gate wait 不佔 permit。Scheduler retry 與 adapter `wirebudget.Wait` 暫放 permit但保 device gate；resume 按 FIFO 公平重取，不在實際 send admission/observer callback 阻塞。
- Request timeout during yielded backoff 不可讓後續 request 逃出容量管理；每 request/legacy point 進入原 read loop 前重取必要 permit，沿原 request/scan deadline、attempt budget。
- 每 request 的 local `invoked` 區分 pure capacity wait 與已進 connector 的原執行失敗：前者維持 Bad 但不毒化 breaker，後者包含 Take 前的 dial/validation failure，保原 health。Shared budget 的絕對 Used 不作 point health 權威。
- Stop 清本輪 pending、取消原 context並 join producers/dispatcher/workers。Restart 建新 queue/permit；舊 worker 捕捉舊 queue/lease，不能修改新版計數。

## 配置與診斷

- `SchedulerConfig.MaxActiveScans`、`PendingScanCapacity` 為 optional `*int`；nil 使用設計 defaults 32／1024，不宣稱最佳效能。顯式 <=0 由 Start 回錯。
- Constructor 解析並凍結 pointer 值；caller 後改 pointer 不改本 scheduler 容量或繞過驗證。舊 partial literals 與 constructor signature 相容。
- `DueDiagnostics()` 是本輪 native snapshot，分 Pending、Running、Coalesced、Rejected、Overruns。Overrun 單位為完成的 dispatched scan，其 capture→completion 超過 cadence 時至多一次；採 Go `time.Now`/monotonic elapsed，不用 source timestamp 或 completion Clock。
- `OutputLossDiagnostics()` 是 scheduler lifetime 累積 loss；只由 `emitValue` 真丟掉的舊值／新值計入。Identity 只取該值已捕捉的 InstallationScope/PointRevision/PointID，缺任一則 explicit Unattributed。
- Recent detail 最多 256 identity series，LRU 淘汰把整個 series loss 加進 EvictedDetailLoss；Total = Unattributed + EvictedDetailLoss + sum(Recent.Losses)。快照為獨立副本，不假稱完整歷史。
- Output drop 不混 due rejection、不改 wire health、不造下游 Good、不改 SampleEnvelope/StableSampleID。

## 回歸證據與限制

- 首輪兩 slow/group fixtures 使用非法非 JSON connection config，屬 setup 假 RED，原失敗保留但不當證據。修合法 JSON 後，以精確 task 3.1 `scheduler_dispatch.go` overlay 再現兩項舊 fast device 1<3，之後新 production GREEN。
- 獨立 review 曾抓出 late old capture 覆新 pending、request timeout 後漏重取 permit、容量一固定 sorted admission starvation；各有 deterministic RED→GREEN，不抹除原 findings。
- Membership 真 SQLite/wire fixture 已適配新 queue：排除 cutoff 後兩個完整21-point rounds（最多 running＋一 pending），再選第三 distinct AcquisitionID；原21 sources、單request、相同PlanFP assertions保留，30s只限制race SQLite觀測backlog，不延長產品deadline。
- Controlled/synctest adapters只證明scheduler契約；loopback wire peers才計實際協議requests，也不是 physical PLC。SampleSink、wire Good與SQL commit不互相冒充。
- Linux/amd64 only；沒有 Mac、真 PLC或正式資料庫驗收。完整前後端套件待七changes總完成後處理，這批只做相關選測。

## 最終實際檢查

環境：Linux/amd64、Go 1.25.5；`GOMAXPROCS=1`、`GOGC=50`、Go `-p 1`。所有 compiler/race/vet/lint/browser 共用單 heavy slot，期間曾讓出給獨立 UI QA；沒有同時跑重工作。

| 檢查 | 結果與證據界線 |
|---|---|
| Final related normal | `final-normal.*`：collector、connector、wirebudget、cmd/test_ui 4 packages，104 top-level／186含subtests PASS，0 FAIL／0 SKIP。包含9個真 production scheduled fixtures。 |
| 同 selection race | `final-race.*`：同104／186 PASS，0 FAIL／0 SKIP，無 race report。cmd package 158.808s；membership fixture整體33.78s含setup，不是產品scan deadline。 |
| Vet | `final-vet.*`：同4 packages `go vet -p 1` PASS。 |
| Official lint 2.6.2 | `lint-first.log` 的21 findings FAIL保留；精確本批patch、原repo設定修後 `lint-second.log` 0 issues PASS。只有dispatcher單行contextcheck例外，理由是snapshot.ctx保原run取消／admission identity；無blanket suppression。 |
| 獨立唯讀 review | `review.md` 保留原findings、health pre-Take／inherited-budget追查及修復ledger。最後26檔0 remaining blockers；review未冒充執行測試。 |
| Static gates | `git diff --check` PASS；官方未修改line script在精確27份copy上檢查 PASS（27 checked，0 warning／failure），未stage/commit產品工作樹。 |
| Source identity | 26行 source SHA manifest digest：`b7dfc4df1665dfd13fd98ab0ad3180fa660113a717ab63d83907eed15d7b6aa0`；final normal/race/vet/lint及post-lint review對應同bytes。 |

Top-level與含subtests不是兩組測試，不相加；normal/race是分開執行，也不製造雙倍測試數。這不是全產品359項或全repo驗證。

工作區證據 `bounded32-evidence` 位於repo旁，不列入Git source manifest；command、JSONL、RED/FAIL、frozen baseline、review與SHA由root保留。27-path SHA與官方scoped line結果由root交接記錄核對。

## 精確 source/test 路徑

```text
internal/datalink/collector/scheduler.go
internal/datalink/collector/scheduler_dispatch.go
internal/datalink/collector/scheduler_lifecycle.go
internal/datalink/collector/scheduler_acquisition_dispatch.go
internal/datalink/collector/scheduler_acquisition_read.go
internal/datalink/collector/scheduler_polling.go
internal/datalink/collector/scheduler_bounded_dispatch.go
internal/datalink/collector/scheduler_due_queue.go
internal/datalink/collector/scheduler_active_capacity.go
internal/datalink/collector/scheduler_output_loss.go
internal/datalink/collector/scheduler_output_diagnostics.go
internal/datalink/collector/scheduler_bounded_test.go
internal/datalink/collector/scheduler_bounded_backoff_test.go
internal/datalink/collector/scheduler_due_queue_test.go
internal/datalink/collector/scheduler_capacity_test.go
internal/datalink/collector/scheduler_output_loss_test.go
internal/datalink/collector/scheduler_due_runtime_test.go
internal/datalink/collector/scheduler_capacity_deadline_test.go
internal/protocol/wirebudget/budget.go
internal/protocol/wirebudget/wait_cooperation.go
internal/protocol/wirebudget/wait_cooperation_test.go
cmd/test_ui/acquisition_scheduled_membership_test.go
internal/datalink/collector/scheduler_due_admission_test.go
internal/datalink/collector/scheduler_acquisition_execute.go
internal/datalink/collector/scheduler_bounded_wire_test.go
internal/datalink/collector/scheduler_capacity_connection_failure_test.go
```
