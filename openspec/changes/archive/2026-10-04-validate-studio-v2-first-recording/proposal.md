## Why

既有驗收能證明預先建表後的採集與交付，不能完整證明一般使用者從空白目的地開始的流程。本案為 A-E 收口，只加入缺少的正式 UI／SQL 與恢復組合驗收，不新開產品能力。

## What Changes

- 在正式 embedded Studio V2、loopback 模擬設備、全新 workspace、空白 SQLite／PostgreSQL 目的地，全程 UI 設定與確認，獨立 SQL 驗證。
- 納入六個審查問題的具體 production 回歸、預設群組與同群組多 entity 的差異，以及延遲補送時間。
- 度量受控首次設定到首筆 SQL 的時間；保留七型別、精確 uint64、品質與既有故障證據。
- 保存 sanitized 可重跑 witness、source/build identity、實際平台與未執行範圍。

## Capabilities

### New Capabilities

無；沿用現有 capability，不建立第二套採集或交付模型。

### Modified Capabilities

- `studio-v2-device-to-sql-acceptance`

## Impact

限定既有 `scripts/tests/f_device_to_sql/`、frontend e2e、相鄰整合測試與驗收文件。需要故障點時沿用既有受 build tag 限制且預設停用的 fixture；正常 binary 不加入故障開關。

## Scope and Dependencies

不新增平台、效能專案、正式環境部署或測試框架。A-E 各自先交付 focused tests，F 的系統驗收不得代替它們。發現本範圍 bug 回唯一 owner 修，不藉驗收擴充產品模式。

前置：`streamline-studio-v2-recording-setup`。
本案只起草；產品 tasks 全部未完成。共同邊界與來源見 [總覽](../../../docs/plans/studio-v2-flow-completion/README.md)。
