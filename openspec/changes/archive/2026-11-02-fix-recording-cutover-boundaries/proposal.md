## Problem

已確認的 cutover 邊界目前可能只截短輸出 metadata，沒有排除切換後才收到、卻先進入舊 boundary 的值。基準 `913cead9` 的 API＋SQLite 重現以相對時間 22.5 秒切換、延遲 reconcile，舊表寫出 23 秒的值 7，而非 22.499 秒的值 6；另一 long fixture 的 row end 為 12.5 秒、observed_at 卻為 13 秒且值 200，違反原有不雙寫、精確歸屬契約。

## Root Cause

`Pipeline.reconcileGroup` 只在輪詢後呼叫 `GroupBoundary.SetEndsAt`；`Assembler.SetEndsAt` 更新邊界卻未重新分配已入列 observation，`markTruncation` 只標記實際窗界。cutover 的 persisted ownership 與 intake／closure 之間缺少共同的線性化屏障；只加快輪詢或只改 row end 不能消除競爭。

## Proposed Solution

- 用已持久化的 cutover intent／revision end 作為唯一 owner authority，讓 cutover commit、sample acceptance 與 closure 共同檢查 ownership fence，不以 reconcile 時點決定樣本歸屬。
- 固定 UTC 半開區間：observed_at < boundary 屬舊 owner，observed_at = boundary 或更晚屬新 owner；僅在目前已支援的 cutover／rollback 路徑內修正。
- 對已入舊 buffer／journal 而尚未 finalized 的跨界 sample，依原 sample identity 原子轉交正確 owner 並重建 eligible set；舊窗 filter 與新 owner durable 接收必須同時成立，不能只濾掉而遺失、複製到兩邊或篡改 finalized payload。
- `freeze-counter-replay-semantics` 先落地以保護移交時的 counter proof；兩項修正均為 `manage-recording-lifecycle-safely` 的前置，新增 layout／profile 轉換矩陣與 successor generation 不在本 change。

## Success Criteria

- 上述固定 fixture 的舊窗分別保留 6 與 100；7／200 只能由新 owner 接收；恰好位於邊界的 observation 永遠只屬新 owner。
- reconcile 延遲、cutover 與 intake／closure 競爭、重啟、失敗後重試及現有 rollback，不造成任何已接受樣本無主或同一 owner scope 重複消費。
- legacy wide／custom、managed wide profile 與 managed long 的既有支援路徑，同樣以 actual window bounds 選出 latest／mean／counter 輸入，保留既有缺值政策、partial evidence 與新 cumulative baseline。
- 已 finalized 的舊效果維持原 payload／目的地；無法證明安全 ownership 時阻擋該 scope 並保留資料，不默默丟棄。本次只有規格交接，未實作修正。

## Impact

- Affected specs: `long-table-group-recording`；新增 cutover ownership barrier／buffer handoff 的細化要求，補足既有 `Explicit forward cutover preserves wide history`，並遵守 circuit/window-average 的半開窗要求。
- Affected code:
  - Modified: `internal/datalink/workspace/write_group_cutover_apply.go`, `internal/datalink/grouppipeline/reconcile.go`, `internal/datalink/grouppipeline/pipeline.go`, `internal/datalink/grouppipeline/recovery.go`, `internal/datalink/runtime/group_boundary.go`, `internal/datalink/runtime/group_boundary_build.go`, `internal/datalink/snapshot/assembler.go`, `internal/datalink/groupdelivery/ledger.go`。
  - New: `internal/datalink/grouppipeline/cutover_ownership_test.go`, `internal/datalink/runtime/group_boundary_cutover_replay_test.go`, `internal/datalink/snapshot/cutover_eligibility_test.go`, `internal/api/router_studio_v2_cutover_barrier_test.go`；必要的 ownership-fence persistence 使用既有 ledger／新的向前 internal migration。
  - Removed: 無；不刪除來源表、歷史、journal 或既有 receipt。
- 保留 reviewed destination／member diff、schema confirmation、workspace／group／connector CAS、ownership、防重複 legacy writer 與既有 late-after-close 契約。
