## 1. 持久化 owner authority 與 cutover 屏障

前置 change：`freeze-counter-replay-semantics` 的 immutable proof／legacy gate 必須先完成。本 change 不新增 lifecycle 轉換矩陣；`manage-recording-lifecycle-safely` 在本 change 後承接。每項先寫會失敗的 deterministic 測試，再完成最小修復及重構；每項結束須可建置、相關測試通過。

共同範圍依 design「既有切換矩陣界線澄清」保留原已支援 managed long／managed wide profile targets；各項 fence、handoff、rollback 測試沿用既有矩陣，不新增轉換能力。

- [x] 1.1 落實「持久化 ownership fence 統一線性化點」：先以 source end／target start／fence 同 transaction 的故障測試鎖住 Cutover ownership is enforced at durable acceptance and closure，再擴充既有 cutover intent，驗證成功 commit 才發布 owner epoch，失敗保留舊 owner。
- [x] 1.2 [after: 1.1] 使 acceptance 以 durable owner fence 而非 reconcile cache 決定歸屬：先在 `cutover_ownership_test.go` 暫停 reconcile 並送 boundary 前／等於／後 samples，再實作 stale-fence retry／typed refusal，新 worker 未就緒只能 durable pending handoff 後 ACK，不能舊 owner 代收。
- [x] 1.3 [after: 1.1] 完成「並行階段與唯一寫入責任」的 workspace cutover-side finalized-overlap transaction gate：先寫已 finalized future-observation、closure-wins 競爭、post-reclaim、合法 boundary equality／nonoverlap fixtures，再在 `applyCutoverInTx` 依全部 `plan.ended` 與既有 checkpoint／effect authority 拒絕衝突，提供安全 `recording-cutover-finalized-overlap` 和 fresh-review action；驗證失敗不改 history、source revision、baseline 或既有 effects，不因回收後缺 rows 推論無歷史、不拒絕所有 checkpoint。僅沿用既有矩陣；本項不承擔 cutover-wins stale `CommitClosure`，完成不代表完整 intake／closure barrier。

## 2. 已接受資料的 filter 與完整交接

- [x] 2.1 [after: 1.2, 1.3] 完成 Buffered cutover inputs are transferred without loss or reinterpretation 的「半開窗 filter 與原子 durable handoff」：先以同一 sample 重試及 transaction 中斷 fixtures，再原子保存 source disposition、target acceptance／pending intent、stable transfer ID，驗證原 payload/digest 保留且兩端消費交集為空。
- [x] 2.2 [after: 2.1] 重建 legacy/custom wide 與 managed long/wide-profile 的 latest eligible set，不只修改 row end：先以 22.499 秒值 6／23 秒值 7、12.5 秒邊界前值 100／13 秒值 200 fixtures，再修復 selector rebuild，驗證舊窗 6／100、新 owner 7／200，較早有效候選不因移除最新值而遺失。 同步啟用 live handoff／Reconcile：source 與 target 兩端都依 actual bounds 重建 eligible set 並刷新同一 read-tx 捕捉的 receipt／input／exclusion view；先驗證 stale source 仍含 transferred input、stale target 漏 input 時 closure 均安全拒絕，再驗證兩端一致重建後才解除 transfer guard，不允許錯誤或 empty closure。保留原 6／7、100／200、較早有效 latest 候選及 actual bounds 驗收；不取代 2.3 reducer 矩陣、2.4 proof 衝突矩陣或 3.1 一般 owner epoch fence。
- [x] 2.3 [after: 2.2] 完成 Actual cutover bounds govern all supported reductions and recovery 的 mean/state/counter_raw/cumulative 分支：先寫同一 accepted set 的 profile matrix，再只以 actual bounds 重建 mean sum/count 與 counter pre-decimation set，驗證 partial evidence、新 owner NULL delta／baseline 及既有缺值政策。
- [x] 2.4 [after: 2.3] 讓 handoff 只改 owner、不改 immutable measurement meaning：先用 incompatible epoch/unit/proof 與 unreadable authority fixtures，再回傳 `recording-handoff-proof-conflict`／`recording-ownership-unavailable` 並保留 pending data，驗證無 current-definition reinterpretation、無未持久化 ACK、healthy scope 繼續。

## 3. Closure、重啟與期限契約

- [x] 3.1 [after: 2.4] 完成「Closure 與 recovery 不跨 ownership 進度」的 closure fence：先暫停 transfer 並競爭 `CommitClosure`，再阻止 stale plan 及有待交接 accepted inputs 的提前 closure，驗證 checkpoint／consumed set 只在有效 fence 下原子前進。 明確承擔 cutover-wins 次序：cutover commit 後 stale `CommitClosure` 必須安全拒絕／重試，checkpoint、consumed set、baseline、outbox 均不得前進；以 deterministic barrier 同時驗證 closure-wins 與 cutover-wins，完整雙次序競爭於本項及 4.3 整合驗收。
- [x] 3.2 [after: 3.1] 在 `group_boundary_cutover_replay_test.go` 先注入 handoff commit 前／後及 closure 前／後 crash，再依 durable disposition／proof 恢復；驗證 sample union 不減、同 scope 只有一個 consuming owner、重試無新增效果、finalized bytes/destination 不變。
- [x] 3.3 [after: 3.2] 分開原本按時 accepted 的延後 handoff 與真正 late-after-close：先寫延遲 reconcile／restart 超過 deadline 及首次接收已逾 deadline fixtures，再保留原 acceptance-time 資格，驗證前者不丟失、後者明確診斷且不重開 finalized window。

## 4. 公開流程與競爭驗收

- [x] 4.1 [after: 3.3] 完成「以 deterministic 競爭 fixture 驗證 cutover」的公開 cutover-review/apply→延遲 reconcile→SQLite readback 路徑；在 `router_studio_v2_cutover_barrier_test.go` 先鎖住原失敗，再斷言 old/new journal、transfer audit、SQL observed_at/value/bounds 與 exactly-boundary 歸屬，不只驗 metadata。
- [x] 4.2 [after: 4.1] 先以現有支援 rollback、部分搬移／保留 source、同 source instantaneous＋cumulative fixtures 驗證，再修正 scope propagation，保證每個 profile-kind ownership scope 唯一 owner、合法獨立 profiles 不被去重且 rollback 不改歷史。
- [x] 4.3 [after: 4.2] 執行包含新增案例的 `go test ./internal/datalink/workspace ./internal/datalink/groupdelivery ./internal/datalink/grouppipeline ./internal/datalink/runtime ./internal/datalink/snapshot ./internal/api`、受影響 package 的 deterministic race tests／`go test -race`／`go vet`、`git diff --check`、`make check-lines`；確認沒有未加 fence 的 acceptance／closure fallback，交接實際 PASS／FAIL／SKIP 與未執行環境，不宣稱外部 DB／現場驗收完成。
