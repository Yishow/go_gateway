## Problem

`b8c55aab` 已完成授權的 Studio V2 新手簡化與受管表安全新增欄位，但未建立專責 delta，正式規格仍要求立即全列套用及只建新表。既有 F 計畫入口也仍稱未提交。Review 另外發現多個建表卡並行時，舊 terminal 回覆可清除較新的 unknown 保護。

## Root Cause

新實作與舊 archive 的規格／文件狀態未同步；前端 recovery marker 以群組清除，沒有比對回覆所屬 operation_id。

## Proposed Solution

核對既有四步變更與證據，把已落地的選取／採集分離、明確批次來源對象及預覽、基本／進階保值、目的地優先、受管成員差異、安全加欄位、誠實寫入狀態、Retry/Skip 合約同步到既有 capabilities。以測試先重現再修正 compare-and-clear；已有效的昂貴實機模擬證據按 source hash 沿用，修改範圍重新回歸。

## Success Criteria

- 舊操作的 succeeded/failed 延遲回覆不得清除新 unknown marker；memory/storage recovery 同樣受保護。
- 正式規格與 `b8c55aab` 實際行為一致，保留 DDL token、CAS、unknown 禁止重送與自訂設定保護。
- F 入口連結與歷史／現況準確；驗證記錄區分已執行、自動化、外部未執行。
- 專責 change 通過 validate、review、受影響測試後才 archive，提交並核對遠端 main。

## Capabilities

### New Capabilities

無。

### Modified Capabilities

- `datalink-workbench-v2-step2-rule`: 點位網格點選只選取，採集使用明確勾選與批次操作。
- `runtime-dashboard-truthfulness`: 無錯誤不得冒充全部 SQL 交付。
- `datalink-workbench-v2-step3-mapping`: 基本欄位與保值展開；明確來源、勾選對象、差異確認的批次套用。
- `guided-recording-workflow`: 點選與採集分離、白話四步、目的地優先與受管成員差異、保值與寫入證據。
- `recording-database-setup`: 合併準備／唯讀預覽、operation compare-and-clear、安全 confirmed add_columns。
- `runtime-write-group-delivery`: Retry 主動作與明確次要 Skip 確認，保留逐筆稽核。

## Impact

- Affected code:
  - Modified: frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/schemaOperationRecovery.ts
  - Modified: frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel.tsx
  - Modified: frontend/tests/unit/workbench-v2/groupSchemaPanel.test.tsx
  - Modified: frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - Modified: frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - Modified: internal/datalink/recordingplan/schema_group_extend_test.go
  - New: frontend/tests/unit/workbench-v2/schemaOperationRecovery.test.ts
  - New: docs/releases/studio-v2-change-closeout-20261023.md
  - Removed: 無。
- 文件：docs/plans/studio-v2-flow-completion/README.md；六份既有 capability 由 archive 套用 delta。
- 不新增 dependency、API 或平台能力。
