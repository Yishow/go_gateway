## Why

已套用群組的草稿儲存、離線重啟與歷史版本排空，沒有共用一致的執行資格。上一輪審查指出可能停寫、遺留未封桶 journal，以及責任時間邊界的並行讀寫；本案先用正式接線重現，再修復這些既有契約。

## What Changes

- 儲存草稿不退休仍有效的 applied revision；只有明確停用、刪除或已生效切換改變 intake。
- 已驗證的 applied layout 可從本地恢復；遠端離線阻擋交付，不阻擋可安全恢復的新樣本 journal。
- 停用、刪除與被取代版本的已接受 journal 可在重啟後只排空、不復活採集。
- 統一責任時間邊界的鎖與持久化切換證據；保留一個輸出一個 writer。

## Capabilities

### New Capabilities

無；沿用現有 capability，不建立第二套採集或交付模型。

### Modified Capabilities

- `studio-v2-write-groups`
- `runtime-write-group-delivery`

## Impact

限定 `internal/datalink/workspace/write_group_*`、`grouppipeline/`、`runtime/group_boundary.go`、必要的既有 groupdelivery 恢復查詢與相鄰測試；正式接線在 `cmd/test_ui`。必要本地 migration 只做 additive，不建立另一個 queue。

## Scope and Dependencies

只修草稿／切換／停用／重啟的既有生命週期，不改採樣協議、UTC 分桶、重試策略或目的地。來源：`write_group_service.go::Update`、`write_group_lifecycle.go::IntakeEligibility`、`grouppipeline/reconcile.go`、`runtime/group_boundary.go::AcceptSample/SetUntil`。離線恢復與歷史 journal 缺口目前是跨模組靜態證據，不能直接當作本輪 crash 實測。

前置：無；可先獨立補回歸。
本案只起草；產品 tasks 全部未完成。共同邊界與來源見 [總覽](../../../docs/plans/studio-v2-flow-completion/README.md)。
