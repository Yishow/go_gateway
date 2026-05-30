## Why

Step 2 的來源規則現在仍停留在前端本地編輯，無法成為 workspace 內真正可持久化的規劃資料。若不先把規則層變成合法即自動存，後面的 point/mapping/output 都只能建立在假資料上，無法支撐多台設備逐步保存。

## What Changes

- 讓 Step 2 的規則建立、修改、刪除改成 workspace-scoped 真實後端操作
- 只有合法規則內容才寫入後端；不合法內容留在畫面上但不覆蓋最後成功版本
- 多條規則逐條保存、逐條回報錯誤，不做整批回滾
- 固定規則與所屬設備的 workspace 關係，讓多台設備規劃可持久化

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `datalink-workbench-v2-step2-rule`: 增加逐條 autosave、逐條錯誤隔離與 workspace-scoped rule ownership
- `datalink-api`: 新增 workspace-scoped rule autosave API 契約

## Impact

- Affected specs: `datalink-workbench-v2-step2-rule`, `datalink-api`
- Affected code:
  - Modified: `frontend/src/features/datalink/workbench-v2/steps/step2/`, `frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts`
  - New: `frontend/src/services/studioV2Rules.ts`, `frontend/src/hooks/datalink/useStudioV2Rules.ts`
  - Modified: `internal/api/router.go`, `internal/api/handlers/source_rule_handler.go`, `internal/datalink/sourcerule/`, `internal/datalink/workspace/`
  - Removed: none
