## Why

目前 Studio V2 的 Step 1 到 Step 4 以前端 reducer 記憶體狀態搭配非同步 autosave 運作，重載頁面後只能部分重新 hydrate 已保存資料，且無法明確區分「已持久化的後端真相」與「尚未保存的本地草稿」。這會讓操作員在重載或服務重啟後看到設定回退、狀態誤判，無法信任 Step 1 到 Step 4 的設定是否真的保存。

## What Changes

- 建立 Studio V2 reload recovery 契約，要求 Step 1 到 Step 4 在 shell 顯示前先以已持久化資料完成 hydrate。
- 將 Step 1、Step 2、Step 3、Step 4 的 reload restore 行為與 save truth 顯式化，避免 factory defaults 偽裝成已保存狀態。
- 明確定義未完成 autosave 的本地草稿在 hard refresh 後的結果，系統必須提示草稿遺失，而不是沉默覆蓋。
- 讓 shell 上層的 committed/runtime 指示器在 reload 後以已持久化 activation truth 為準，而不是回到本地預設 idle。

## Capabilities

### New Capabilities

- `studio-v2-reload-recovery`: 定義 Step 1 到 Step 4 的 persisted truth、reload recovery、draft loss signaling 與 restore 順序。

### Modified Capabilities

- `studio-v2-workspace`: workspace bootstrap 必須支援組合出完整且一致的 persisted setup snapshot。
- `datalink-workbench-v2-step1-device`: Step 1 reload 後必須回到 persisted device state 與正確 save truth。
- `datalink-workbench-v2-step2-rule`: Step 2 reload 後必須回到 persisted source-rule state 與正確 save truth。
- `datalink-workbench-v2-step3-mapping`: Step 3 reload 後必須以目前 point identity 重建 persisted mappings。
- `datalink-workbench-v2-step4-database`: Step 4 reload 後必須回到 persisted database connector/targets 與正確 save truth。
- `datalink-workbench-v2-shell`: shell 的 runtime/committed 指示器必須反映 restore 後的 persisted activation truth。

## Impact

- Affected specs:
  - New: studio-v2-reload-recovery
  - Modified: studio-v2-workspace, datalink-workbench-v2-step1-device, datalink-workbench-v2-step2-rule, datalink-workbench-v2-step3-mapping, datalink-workbench-v2-step4-database, datalink-workbench-v2-shell
- Affected code:
  - Modified: frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - Modified: frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - Modified: frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - Modified: frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - Modified: frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - Modified: frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - Modified: frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - Modified: frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - Modified: frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - Modified: frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - Modified: frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - Modified: frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - Modified: internal/api/handlers/studio_v2_workspace_handler.go
  - Modified: internal/datalink/workspace/service.go
  - Modified: internal/datalink/workspace/service_devices.go
  - Modified: internal/datalink/workspace/service_database.go
