## Problem
group 退休 boundary 清除後 legacy writer 可能恢復寫入。blocked／quarantine head 後的已接受 rows 沒有 operator repair 路徑。

## Root Cause
service_wiring 使用 groupPipe.Owns 抑制仍 enabled 的 legacy mappings；Owns 只看 active boundaries。Store.ResolveQuarantine 已有 retry／skip，但 router/UI 只讀 delivery 統計。

## Proposed Solution
將首次有效接管的 durable ownership 與 active worker 分離，停用／刪除／換目的地／重啟均不恢復舊 writer；首次接管前 legacy 繼續合法運作。增加 group scoped、可追溯 retry／skip 與最小 UI；保留 frozen destination 及 unknown safety。

## Success Criteria
- production chain 回歸涵蓋 migration、成功接管、disable/delete/destination change/restart，以及未接管 legacy 正常寫入。
- missing table／permission／bad row 修好後 retry 或明確 skip 可釋放 head，舊目的地身份不改，unknown 不可重送。
- 回歸先失敗，再 GREEN、review→fix→rereview；不重寫 migration。

## Capabilities
### Modified Capabilities
- runtime-write-group-delivery: durable takeover 與 operator 恢復。

## Impact
- Affected code:
  - Modified: cmd/test_ui/service_wiring.go, internal/datalink/grouppipeline/pipeline.go, internal/datalink/grouppipeline/reconcile.go, internal/datalink/groupdelivery/outbox.go, internal/api/router_studio_v2_write_group_routes.go, frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupDeliveryStrip.tsx
  - New: ownership/recovery persistence、group scoped handler/service/hooks/types 與 focused regressions。
