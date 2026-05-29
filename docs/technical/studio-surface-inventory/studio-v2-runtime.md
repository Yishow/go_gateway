# `/studio/v2` 與 `/studio/runtime` Inventory

## 這兩個 surface 的產品定位

| Surface | Route | 受眾 | 定位 |
| --- | --- | --- | --- |
| Studio V2 | `/studio/v2` | 一般使用者 | 簡單、可引導、好觀察的 setup flow |
| Runtime Dashboard | `/studio/runtime` | setup 後操作者 | focused post-setup runtime monitor |

## 目前策略

- `/studio/v2`：**這一輪重點施作**
- `/studio/runtime`：**作為 V2 的 post-setup 觀察面一起施作**
- 預設入口方向：**網站服務應預設導向 `/studio/v2`**

## `/studio/v2` 現況

### 核心事實

- 已有完整 shell、四步驟、settings、summary rail、Step 4 handoff。
- 大多數資料仍存在 `useWorkbenchV2State()` reducer 中。
- `commit` 目前是前端模擬流程，不是正式的 backend persisted lifecycle。
- 使用者期待未來「服務預設入口」是 `/studio/v2`，但目前 router 預設仍不是它。

## V2 的產品要求

1. 簡單：不要暴露 `/studio` 那種完整但高密度的 operator 心智模型。
2. 有指引：每一步要能告訴使用者現在在做什麼、下一步是什麼。
3. 好觀察：setup 後應自然 handoff 到 focused runtime 觀察面。
4. 可當預設入口：不能再依賴純 mock / draft commit 才能走完主流程。

### Page / Step matrix

| 區塊 | 主要使用者行為 | 目前前端資料來源 | 實際 API wiring | 狀態 |
| --- | --- | --- | --- | --- |
| Shell / step rail / summary rail | 切步驟、看摘要、收合 UI | `useWorkbenchV2State()` | 無 | `frontend-local-only` |
| Step 1 Device | 編輯設備草稿、跑 staged test、看 connect/probe 階段 | local reducer state | 無正式 datalink API | `frontend-local-only` |
| Step 2 Rule | 編輯規則、生成 points 視圖 | local reducer state | 無 | `frontend-local-only` |
| Step 3 Mapping | 調整 mapping / transform | local reducer state | 無 | `frontend-local-only` |
| Step 4 Database | 調整 DB target、看 commit progress | local reducer + `buildCommitLogSequence()` | 目前僅模擬 endpoint label | `frontend-local-only` |
| Settings | 編輯 connector pool、scheduler、modbus-share 設定草稿 | local reducer state | 無 | `frontend-local-only` |
| Runtime handoff | 提交成功後前往 runtime dashboard | route navigation | 無 commit API；僅 URL handoff | `semantics-gap` |

### 重要檔案

- `frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx`
- `frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx`
- `frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts`
- `frontend/src/features/datalink/workbench-v2/state/commitLog.ts`
- `frontend/src/features/datalink/workbench-v2/shell/resolveRuntimeDashboardDevice.ts`

### V2 現有 handoff 契約

| 項目 | 現況 |
| --- | --- |
| handoff target | `/studio/runtime` |
| query param | `device_id` |
| guard | 僅 UUID 形態才視為 persisted backend device id |
| fallback | 取不到 persisted id 時，退回 `/studio/runtime` 不帶 `device_id` |
| 風險 | Step 4 commit 還沒真正產出 persisted device / point / runtime context |

### V2 必補的正式契約

1. commit success response 必須回 persisted `device_id`
2. commit success response 必須回 runtime 啟動狀態
3. 若 commit 會建立 points / mappings / outputs，最好回最小 summary
4. route handoff 不能再從 local draft state 猜 runtime device
5. `/` 與主要入口 redirect 應可直接落到 V2

## `/studio/runtime` 現況

### 使用者可見動作

- 依 `device_id` 聚焦單台設備
- 讀取 snapshot
- 開啟 SSE
- 看 summary / collector health / live points / logs
- 重新選擇 device
- 在 stream 失敗時退回 degraded polling

### 前端主要檔案

- `frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardRoute.tsx`
- `frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts`
- `frontend/src/features/datalink/runtime-dashboard/useRuntimeStatus.ts`
- `frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts`
- `frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx`

### 對應 hooks / API

| 類型 | 前端 | 後端 API | 狀態 |
| --- | --- | --- | --- |
| device list | `useDevicesQuery()` | `GET /api/v1/datalink/devices` | `wired` |
| device-scoped points | `usePointsQuery({ device_id })` | `GET /api/v1/datalink/points` | `wired` |
| runtime snapshot | `useRuntimeStatus()` | `GET /api/v1/datalink/runtime/status?device_id=...` | `wired` |
| runtime stream | `useRuntimeDashboardStream()` | `GET /api/v1/datalink/runtime/stream?device_id=...&point_ids=...` | `wired` |
| status events | SSE `status` | same stream endpoint | `wired` |
| heartbeat | SSE `heartbeat` | same stream endpoint | `wired` |

### runtime route state machine

| route state | 來源 | 備註 |
| --- | --- | --- |
| `missing-device-context` | 沒有 `device_id`，或 snapshot 報 `device not found` | 目前部分判斷仍依賴錯誤字串 |
| `loading` | 尚未拿到首個 snapshot | 正常首屏 |
| `live` | snapshot 已就緒且 SSE connected | 正常監看 |
| `degraded` | 有 snapshot，但 SSE error | 退回 polling |
| `error` | snapshot 失敗且非 missing-device-context | 真正錯誤 |

### Runtime 契約缺口

| 類型 | 問題 | 影響 |
| --- | --- | --- |
| semantics-gap | `device not found` 目前從 500 + 錯誤字串推導 | 前端必須做字串判斷 |
| semantics-gap | 缺 `starting` / `committed but not running` / `not_committed` 類狀態 | setup 後無法正確解釋等待期 |
| semantics-gap | runtime 與 setup commit 沒有正式 lifecycle handoff | V2 不能穩定直達 focused runtime |
| exists-not-wired | runtime snapshot 只有 post-setup route 用 | `/studio/v2` 內部還不能當正式 commit result source |

## 使用者期望與現況的差距

### 使用者想要的 V2

- 簡單
- 有引導
- 容易觀察
- 作為服務預設入口

### 目前缺的不是 UI，而是 contract

1. V2 沒有正式 commit API
2. V2 沒有 persisted output summary
3. Runtime 沒有產品化 lifecycle semantics
4. Router 還沒改成以 V2 為預設入口

## V2 近期應優先完成的順序

1. `make-studio-v2-default-entry`
   - 先讓產品入口與使用者方向一致
2. `integrate-studio-v2-commit-with-runtime-lifecycle`
   - 讓 V2 不再只是 local-only setup shell
3. `normalize-runtime-lifecycle-status-contract`
   - 讓 runtime 成為可靠的 post-setup observer

## V2 暫時不該膨脹成的樣子

1. 不要變成第二個 `/studio` 完整版。
2. 不要把 fleet-wide dashboard、全局 review desk、工程工具全塞進來。
3. 不要在 commit contract 未落地前，繼續用更多模擬 log 或 reducer state 掩蓋真實 lifecycle 缺口。
