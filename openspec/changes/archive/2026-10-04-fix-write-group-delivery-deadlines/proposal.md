## Why

Sender 把 5 秒本地結果保存期限與 30 秒遠端操作期限一起啟動。遠端合法耗時超過前者後，本地即使仍可寫入，也會失去登記成功或重試的時間；本案只修期限的起算位置。

## What Changes

- 每次遠端嘗試結束後，才開始計算本地 settlement 的完整期限。
- 成功、retry、blocked、quarantined、unknown 的結果保存皆遵守相同規則。
- 保留 lease/fencing、去重、未知結果與 shutdown 的既有邊界，不改預設期限或重試次數。

## Capabilities

### New Capabilities

無；沿用現有 capability，不建立第二套採集或交付模型。

### Modified Capabilities

- `runtime-write-group-delivery`

## Impact

限定 `internal/datalink/groupdelivery/sender.go`、必要相鄰 helper 和 sender／worker 回歸。來源為 `Deliver` 的 attemptCtx、settleCtx 建立順序；不更動目的表、API 或 UI。

## Scope and Dependencies

只改交付結果記帳期限。隔離範例只能證明 context 模式，實作需以正式 Sender 接 disposable SQL 驗證；不藉此重寫 retry engine。

前置：無；可先獨立補回歸。
本案只起草；產品 tasks 全部未完成。共同邊界與來源見 [總覽](../../../docs/plans/studio-v2-flow-completion/README.md)。
