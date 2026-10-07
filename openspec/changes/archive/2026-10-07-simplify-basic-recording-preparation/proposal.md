## Summary
Basic 四步設定可完成空白資料庫 preparation，採唯一有效 canonical group 設定，主畫面保留名稱／值／型別單位／記錄狀態。

## Motivation
BasicRecordingPanel 將 schema preview/confirm 留在 Advanced，fresh-ui 驗收也必須往返。ConnectorSection 暴露舊 table／INSERT-UPSERT／write interval，舊 SchemaSetupSection 未有 group preview token；KindSelector 仍能選 group 不支援 driver。

## Proposed Solution
在 Basic 重用 GroupSchemaPanel 的既有 preview/token/confirm/operation，不新增 schema API；明確 DDL 確認後才建表。Connector 表單與 SummaryRail 只顯示連線，group 掌握 storage/interval/members；移除主線舊 nonpreview 建表入口。依既有 capability 禁用不支援 MySQL/SQLServer 並解釋。原始 IDs/API payload/revision 等進既有 diagnostics 或 disclosure，只做直接必要 UI 調整。

## Capabilities
### Modified Capabilities
- guided-recording-workflow: Basic preparation 與 operator 主畫面。
- recording-database-setup: 同一安全 schema 與 capability 支援契約。

## Impact
- Affected specs: guided-recording-workflow, recording-database-setup
- Affected code:
  - Modified: frontend/src/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel.tsx, frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx, frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx, frontend/src/features/datalink/workbench-v2/steps/step4/KindSelector.tsx, frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx, scripts/tests/f_device_to_sql/fresh-ui.mjs
  - New: Basic schema/capability/diagnostics focused Vitest regressions。
