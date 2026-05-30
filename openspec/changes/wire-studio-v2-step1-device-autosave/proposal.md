## Why

Step 1 目前只是前端本地設備編輯器，合法與不合法內容都停留在 reducer 內。若不先把設備資料變成真正的後端自動存，後面的規則、mapping、第一次啟動與 runtime 工作區視角都無法建立在真實資料之上。

## What Changes

- 讓 Step 1 的設備建立、修改、排序與刪除改成 workspace-scoped 真實後端操作
- 只對合法設備內容做自動存；不合法內容留在畫面上但不寫入後端
- 多台設備採逐台保存與逐台錯誤回報，不因單台失敗阻塞其它合法設備
- 刪除已保存設備時，同步刪除該設備的 V2 關聯資料

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `datalink-workbench-v2-step1-device`: 增加逐台自動存、未存成功標示與刪除後端資料的要求
- `datalink-api`: 新增 workspace-scoped device autosave API 與 device order API

## Impact

- Affected specs: `datalink-workbench-v2-step1-device`, `datalink-api`
- Affected code:
  - Modified: `frontend/src/features/datalink/workbench-v2/steps/step1/`, `frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts`, `frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx`
  - New: `frontend/src/services/studioV2Devices.ts`, `frontend/src/hooks/datalink/useStudioV2Devices.ts`
  - Modified: `internal/api/router.go`, `internal/api/handlers/device_handler.go`, `internal/datalink/device/`, `internal/datalink/workspace/`
  - Removed: none
