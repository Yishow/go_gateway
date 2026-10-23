## Problem

main 524555f1 的實測顯示 MC 多 word 單點讀取變為陣列，阻擋首次 SQL。既有範本確認、runtime 頁首、範圍摘要與 Basic 缺值選項另有可重現誤導。

## Root Cause

單點呼叫傳 words，MC adapter 卻用型別值數解讀；範本未在確認時改 confirmed；頁首未依證據；範圍終點只算最後起址；partial 文案未揭露 required 優先。

## Proposed Solution

依 R1→R5 修復既有契約並先測失敗。Modified Capabilities: protocol-connectors、measurement-semantics、runtime-dashboard-truthfulness、datalink-workbench-v2-step2-rule、guided-recording-workflow。New Capabilities: none。

## Non-Goals

不新增前端、資料模型、報表平台或 Windows tray 功能；不碰現場 PLC、正式 DB、部署，不放寬型別、DDL/token/unknown、required、版本保護。

## Success Criteria

MC periodic/manual/read-now32/64-bit scalar及直接批次保持正確；隔離 MC D0–D7混合實際UI至少三個SQL桶。範本真preview=false確認後才套用，失敗可恢復；header無證據不報入庫；範圍包含末word；Basic揭露required並可進既有成員設定。完整Go/frontend gate、受影響browser與獨立review通過，未執行Windows/PG明列。

## Impact

- Affected code:
  - Modified: internal/datalink/collector/scheduler_polling.go, internal/datalink/collector/scheduler_manual_poll.go, internal/datalink/device/service_point_read.go
  - Modified: frontend/src/features/datalink/workbench-v2/steps/step3/MeasurementTemplateSelector.tsx, frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - Modified: frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx, frontend/src/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel.tsx
  - New: internal/datalink/connector/point_read.go, task-scoped tests and isolated simulator evidence
  - Removed: none
