## Problem
預設 UI 將非字串 point 映射為 float64，uint64 9007199254740993 可在 SQL 前失去精度。Runtime 同址雙設備可能互相顯示讀值。

## Root Cause
buildDefaultMapping 選 float64；executeCast 經 float64 且忽略失敗。LivePointsTable 遍歷 workspace mappings，地址 fallback 未驗設備身份。

## Proposed Solution
保持來源型別；整數轉換檢查精確性、範圍與有效性，失敗回傳錯誤。Runtime 依 selected device、point 身份顯示，舊 SSE 不得落到新設備。同批只將值／型別單位與名稱做直接必要的可讀調整。

## Success Criteria
- 預設 UI 到 production mapping／journal／outbox／SQLite 仍保有 uint64 9007199254740993；非法值不變 0。
- A.D0=215、B.D0=187 切換不混用值、名稱與單位；等待新設備時不可冒稱 live。
- 有失敗回歸、修正、scope review 與 exact HEAD 檢查證據。

## Capabilities
### Modified Capabilities
- mapping-pipeline: 保留型別與精確失敗契約。
- post-setup-runtime-dashboard: 設備隔離顯示。

## Impact
- Affected code:
  - Modified: frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts, frontend/src/features/datalink/runtime-dashboard/LivePointsTable.tsx, frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx, internal/datalink/mapping/pipeline_execution.go, internal/datalink/mapping/pipeline_value_helpers.go
  - New: 對應 Go／Vitest focused regressions 與驗證紀錄。
