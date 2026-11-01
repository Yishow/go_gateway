# Manual planned acquisition：task 3.3 交接（2026-10-31）

## 歸屬與範圍

- Change：`integrate-protocol-aware-acquisition-batching`，只交付 task 3.3；[tasks](../../../openspec/changes/integrate-protocol-aware-acquisition-batching/tasks.md)、[design](../../../openspec/changes/integrate-protocol-aware-acquisition-batching/design.md)。task done／commit／push 由 root 核對，不由本 owner 執行。
- 基準是 task 3.1＋3.2 已完成的工作樹 bytes，不是 Git HEAD；既有檔於首次修改前保留基準，精確差異在 repo 旁 `manual33-evidence/owned.patch`。不把前兩階段 WIP 歸成本次修改。
- 沿既有 scheduler、immutable planner、bounded executor、device gate、active capacity、ManagedConnection 與 wirebudget；沒有第二採集引擎、新 decoder、record selector authority、public API 或 UI。
- 4.1 preview／actual diagnostics 展示與 4.2 benchmark 仍未交付；本批不是速度、真 PLC、SQL durability 或部署驗收。

## 手動意圖與不可變計畫

- `PollNowContext` 在一個 scheduler 短鎖內凍結所選 canonical Point、device config、原 group/cadence、ReadProfile、revision、installation lease/scope 與 geometry；離鎖後才 planning、等 gate、取 connection 或 I/O。
- 重複 requested ID 至多一個 final；unknown／recoveryPending／空集合不擴成全部來源。原已安裝 disabled 點的 explicit manual 意圖不靠改 Enabled 或重裝 source 實現。
- 每 manual/device 一個 goroutine、一個 coherent plan、一個 AcquisitionID、一個 gate ownership；跨原 polling groups 的 Sources 保各自 group/cadence，request index 在同 plan 唯一，不重置成各 group 的 `/initial/0`。
- AcquisitionID 原本進 StableSampleID／SampleEnvelope，維持每 device 一次；RequestID／PlanFingerprint 仍只是 transient collector evidence。同 request 保真實 final completion，跨 request 保各自時間，不冒稱 PLC atomic snapshot。
- 同完整單 group 的 manual 與 scheduled 計畫身份相同；入口、caller context、時間與 observation 開關不參與 fingerprint。尚未 Start 而無 cadence 資料時保既有合法 0／未指定，不杜撰 1s 或 fake group。
- Explicit requested spans 沿原 planner：Excluded 只約束 gap，可由實際 Point span 覆蓋；Forbidden 仍拒絕該地址。沒有另造地址 parser 或移除 geometry fencing。

## 執行、回傳與生命週期

- 原 executor 抽出 final sink：scheduled 維持 `emitValue`；manual 只回傳給原 caller，不再送 ValueChannel，不造成第二份 ingestor／SQL delivery 或 output-loss。
- 原 `collectedValue` 的 observer 行為保留；manual 只對 early cancellation／supersession 剩餘的 captured Points 補一份 Bad，不用新 install metadata，也不把中間 retry 當 final。
- MC/FATEK/Modbus 使用同 planned request／fallback／deadline／actual-send fence；不為只有 scalar Read 的 fake adapter 放寬 bounded PLC 的 PlannedReader 要求。
- 非 bounded custom protocol 保原單次 scheduler-level scalar invocation；重用 scalar count/read/result helper，不把此例外誤說成 adapter 只准一次 wire attempt。
- Active permit pool 在 constructor 解析並凍結容量後建立，屬 Scheduler lifetime；Start 只新建 due queue。Explicit invalid capacity 的 manual 呼叫 fail closed，不永久卡在零容量。
- 可在 Start 前／Stop 後 explicit manual；scheduled Stop 只取消並 join 原 scheduled workers，manual 仍由 caller context 擁有。跨 restart 共用同 gate／permit，不新增 Start-busy 錯誤。
- Gate wait 不占 active permit；scheduler／adapter backoff 保 gate 並暫放 permit，恢復仍沿原公平 acquire、每 request `ResumeAfterWait`、固定 deadline 與 attempts。
- 原 local `invoked` health 分類未改：純 capacity wait 不毒化 breaker，已進 connector 的 pre-Take dial failure 仍保原 operation health；不以 inherited budget 的絕對 Used 當健康權威。
- Scheduled dispatcher 的 per-device goroutine 上界不被冒稱適用於任意數量的外部 manual callers；沒有新增內部 manual pending queue。

## 測試發現與處理

- 真 loopback MC/FATEK/Modbus 新測先確認 baseline scalar 為 Good，再重現 duplicate IDs 的 4→3 final、相同 wire span 的 2→1 requests、空 RequestID 與 capacity-one 跨生命週期失守。
- 初版雙 format 測試誤用 legacy DataFormat；MC/FATEK legacy 原本不採該 override。改用正式 ResolveReadProfile／FormatOverride／lease，以 task-start source overlay 再現真正 2→1 RED，沒有修改 decoder。
- 舊 fake point-count/read-profile fixtures 補 PlannedReader 接點，仍保原 native count、immutable profile、scalar/direct 正向斷言；不把 fake 限制變成 production fallback。
- 舊 concurrent-format fixture 固定等待四個 scalar requests；遷移後兩 caller 各一個 raw read，改先 assert 兩 requests 再讀取，原值／format／fingerprint／lease 斷言全保留。
- Connection config 變更會使舊 Point install 失效：測試新增 stale Bad／無新 connect，明確 reinstall 後仍驗原新 config fingerprint。該完整 test 搬到獨立小檔，避免歷史長檔增行，未刪除原正向驗收。
- 獨立 review 抓到同單 group、同 Sources 僅因 manual 入口產不同 PlanFingerprint；pure parity RED 後收斂單 group 的原 proof，未改 scheduled fingerprint 或多 group request index。
- 真 persisted recording draft 選 D0/D10/D20，新增 explicit manual D1 只讀 D1、只有一次 observer、無 scheduled channel/SampleSink；之後原 21-source scheduled rounds 與 selector 縮減驗收全保留。
- 廣選發現 RTU matrix 在 baseline scheduled 路徑即失敗，而舊 scalar manual 可成功。獨立診斷確認 raw peer 缺 `SendReceiveExactContext`；production exact RTU guard 正確 fail closed。只補既有完整 CRC/frame simulator 的 explicit capability，先 `wirebudget.Take(ctx)` 再原 SendReceive；不動 adapter、native serial transport 或 exact-frame safety guard。
- 所有中間 FAIL／中斷與 baseline overlay 證據保留，不把 setup／fixture 前提錯誤冒充產品 RED。

## 最終相關驗證結果

- Final related normal：5 packages，250 top-level／600 含 subtests PASS，0 FAIL／0 SKIP。Collector 148／322、cmd/test_ui 39／154、runtime 37／96、API handlers 25／27、exact RTU guard 1／1。
- Final related race：相同 5 packages、250 top-level／600 含 subtests PASS，0 FAIL／0 SKIP／data-race warning。最終 cmd normal／race 使用落地 23-path bytes；不變 collector／runtime／API／guard 沿用既有同 production bytes 的 PASS，未重跑無關套件。
- 廣選 cmd 初輪在舊 RTU fixture 失敗而中斷，原 FAIL 不抹除。補 exact-context capability 後完整 RTU matrix／ZeroShortAndLegacy／registry scope／budget-cancel 窄 GREEN，以及原 production negative guard 均 PASS；最後 cmd normal 98.517s PASS。
- Legacy single-invocation 原行為 regression PASS；test-only overlay 故意採 scheduler retry 後，測試因應為 Bad 卻得到 Good 而 FAIL。Production 未被 mutation overlay 改動，final collector normal 再 PASS。
- 原 cmd race 的 1 秒 loopback fixture 在 24-source 同 device gate 競爭時回 deadline Bad；test-only 5 秒 overlay 全格式 PASS。Root 核准兩個 peer config 1→5＋註解，保所有原斷言及失敗證據；落地重跑 cmd 39／154 normal（114.707s）及 race（521.825s）PASS。未改 production timeout／retry／guard，這不是 deadline 效能證明。
- Official lint 僅測試 import grouping findings，已整理 read-profile 與 RTU fixture；最後重跑 0 issues PASS。Vet 五個受影響 packages PASS。
- 獨立唯讀 review：23-path source hashes 全匹配、0 remaining source blockers；已獨立解析 10 份 normal／race JSONL 與 stderr，確認上述各 250／600 PASS；最後 23-path test-config delta official lint 亦為 0 issues PASS。
- 23-path source/test SHA manifest digest：`79d386626ac95489822b51c82f514694ca3fa260248c4ed241e15781acb7f9dd`；兩個 loopback peer 的 5 秒 test config 經 root 核准；最終 normal／race gates 前後核對相同 bytes，不含本文件。
- Linux/amd64、Go 1.25.5；GOMAXPROCS=1、GOGC=50、Go `-p 1`。Go／lint／browser／frontend compiler 共用一 heavy slot；等待 UI RED 時不併跑。
- 只清本 owner 兩個已中斷且 terminal 的 GOTMP build 目錄，共 112 MiB；保留主 GOCACHE、GOPATH、DB、原資料與全部證據。

## 精確 source/test paths

下列 23 paths 加本交接文件；最終 bytes／命令／log 與官方 scoped line gate 由 root 檢查 `manual33-evidence` 後封版。

```text
cmd/test_ui/acquisition_manual_planned_test.go
cmd/test_ui/acquisition_manual_rtu_fixture_test.go
cmd/test_ui/acquisition_scheduled_membership_test.go
cmd/test_ui/word_format_rtu_transport_test.go
cmd/test_ui/word_format_serial_datagram_test.go
cmd/test_ui/word_format_wire_test.go
internal/datalink/collector/point_count_contract_test.go
internal/datalink/collector/scheduler.go
internal/datalink/collector/scheduler_acquisition_dispatch.go
internal/datalink/collector/scheduler_acquisition_execute.go
internal/datalink/collector/scheduler_acquisition_facts_test.go
internal/datalink/collector/scheduler_acquisition_snapshot.go
internal/datalink/collector/scheduler_active_capacity.go
internal/datalink/collector/scheduler_format_concurrency_test.go
internal/datalink/collector/scheduler_lifecycle.go
internal/datalink/collector/scheduler_manual_admission_test.go
internal/datalink/collector/scheduler_manual_backoff_test.go
internal/datalink/collector/scheduler_manual_capacity_test.go
internal/datalink/collector/scheduler_manual_connection_identity_test.go
internal/datalink/collector/scheduler_manual_poll.go
internal/datalink/collector/scheduler_manual_snapshot_test.go
internal/datalink/collector/scheduler_polling.go
internal/datalink/collector/scheduler_read_profile_test.go
```
