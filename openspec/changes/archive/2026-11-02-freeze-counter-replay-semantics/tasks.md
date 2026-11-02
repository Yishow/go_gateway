## 1. 擴充 immutable proof 與確認流程

以下均是待實作任務，文件完成不代表產品完成。每項先寫會因舊行為失敗的測試，再做最小修復並重構；每項結束須可建置且相關測試通過。此 change 先於 `fix-recording-cutover-boundaries`，兩者完成後才開放 `manage-recording-lifecycle-safely` 的新切換能力。

- [x] 1.1 落實「凍結完整 counter proof」的版本化資料契約：以 `counter_binding_proof_test.go` 先鎖住完整 member identity、epoch、unit、kind、scale/type、rollover/rate 與 canonical digest 的 round-trip／tamper 行為，再加入 proof codec，確保 unknown version 被 typed 拒絕且非 counter snapshot 相容。
- [x] 1.2 [after: 1.1] 完成 Immutable counter proof for accepted samples 的未縮放 apply 路徑：在 `saveRuntimeVersionInTx` 前凍結確認時的 definition/source-binding revisions，先以 apply/edit 競爭測試重現缺口，再讓 CAS 保證完整 proof 或無新 applied revision；以 workspace focused tests 驗證不可半成功。
- [x] 1.3 [after: 1.2] 把既有 `ExactRecordingSource` 接到同一 proof 契約而不重讀 definition、不重複縮放；先新增已縮放／未縮放 golden fixtures，再保留 `TestFrozenExactRecordingSourceDoesNotRefetchDefinition` 的零呼叫與 exact digits 斷言，確認既有支援 type 組合未擴張。

## 2. 接受與重播只使用被接受的版本

- [x] 2.1 [after: 1.3] 完成「Acceptance 與 proof identity 原子關聯」：新 sample 只在 proof reference 與 journal 同時耐久後 ACK，失敗不改 memory；以 codec/resend tests 先驗證舊 digest 不變、新 proof reference 納入新版本 identity，再修復 `AppendSample`／decode 的相容行為。
- [x] 2.2 [after: 2.1] 拒絕新 sample 自帶與 active proof 不符的 measurement／epoch／definition／binding claims，空 claims 僅能繼承 applied proof；先以 `counter_binding_replay_test.go` 斷言 typed refusal、無 journal append／baseline advance，再完成 acceptance gate。
- [x] 2.3 [after: 2.2] 讓 `buildFrozen`、`assembleBoundary`、reconcile 與 restart replay 僅使用原 immutable proof：先以 definition edit/delete/outage 的 fixtures 重現 100→105 被改標 MWh，再修復並斷言舊 epoch-1/kWh/delta 5、frozen destination 及零 current-definition read。

## 3. 保留舊資料並安全遷移

- [x] 3.1 [after: 2.3] 完成 Evidence-preserving legacy counter replay 與「Legacy proof 分類遷移與隔離」的可證明分支：先以舊 codec＋完整 exact-source snapshot fixture 驗證，再加入附加 proof association／audit；原 payload、digest 與 duplicate resend 結果必須不變。
- [x] 3.2 [after: 3.1] 對 missing／conflicting／unsupported legacy proof 分別隔離 affected group/revision，提供安全 reason、pending count 及修復指引；先加入 fault fixtures，再完成 gate，驗證不借用 current definition、不 reset checkpoint、不刪除 journal，其他 scope 不受阻。
- [x] 3.3 [after: 3.2] 分離未 finalized replay quarantine 與 finalized delivery：先用混合 backlog fixture 斷言 healthy group 繼續、已 finalized effects 可按 frozen bytes/destination retry，再修正狀態及 delivery gating，確認沒有重算或把 quarantine 顯示為 delivered。

## 4. 跨 epoch、重啟與公開 API 驗收

- [x] 4.1 [after: 3.3] 完成 Replay proof preserves epoch and idempotency boundaries 與「重播與新 epoch 的驗收隔離」：先寫舊 backlog＋新明確 applied epoch fixture，再保證 old proof drain、新 owner 首筆 NULL delta／baseline，counter_raw 與 cumulative 不共享錯誤 baseline。
- [x] 4.2 [after: 4.1] 將公開 measurement PUT→applied snapshot 不變→restart→SQLite readback 加入 `router_studio_v2_counter_replay_test.go`；先重現錯誤，再斷言 100→105 的舊單位／epoch／delta 及新 epoch 隔離，不以 handler 200 或 memory-only 結果代替 SQL 證據。
- [x] 4.3 [after: 4.2] 以 proof persistence、sample commit、closure commit 前後的故障注入測試完成 crash/resend 契約，再移除未 finalized replay 的 current-definition fallback；驗證未 commit 不 ACK、已 ACK 可恢復一次、finalized digest 不變，記錄升版／安全回退與 legacy quarantine 驗收結果。
- [x] 4.4 [after: 4.3] 執行實際涵蓋新增案例的 `go test ./internal/datalink/workspace ./internal/datalink/groupdelivery ./internal/datalink/grouppipeline ./internal/datalink/runtime ./internal/api`、受影響 package 的 `go test -race` 與 `go vet`、`git diff --check`、`make check-lines`；列明 PASS／FAIL／SKIP、SQLite readback 及未測外部資料庫，不把 skip 或歷史結果算通過。
