## Why

前端允許不同 entity 共用資料欄位，但 runtime 的 layout 做全群組欄位碰撞檢查；單列編碼又要求所有 entity 的 members。這會讓既有合法配置被擋，或在採集後變成 skipped，本案修到 UI、readiness 與 SQL 用同一個資料列範圍。

## What Changes

- 欄位衝突與 SQL layout 以 group 中的 entity 為範圍；同一列仍禁止競爭欄位。
- 每列只編碼自己 entity 的 members，不要求其他列的資料。
- 既有合法多 entity 需實際落庫驗收；只有不合法配置才在 Apply 前拒絕。
- 結構性編碼錯誤不得被當成正常 missing/skipped 而消耗已接受資料。

## Capabilities

### New Capabilities

無；沿用現有 capability，不建立第二套採集或交付模型。

### Modified Capabilities

- `write-group-row-semantics`

## Impact

限定 `internal/datalink/dbtarget/group_row_layout.go`、`runtime/group_boundary.go`、既有 snapshot/readiness 的布局銜接與前端 writeGroup 欄位驗證及回歸。`snapshot/assembler.go::finalize` 與 `columns.ts::findColumnConflicts` 是對照來源。

## Scope and Dependencies

是修復目前已允許的 distinct-row 行為，不增加跨設備同步、聚合、join 或新 row policy。基本 UI 預設每台設備一組由後續操作案負責；不能以隱藏進階入口代替本案修復。

前置：無；可先獨立補回歸。
本案只起草；產品 tasks 全部未完成。共同邊界與來源見 [總覽](../../../docs/plans/studio-v2-flow-completion/README.md)。
