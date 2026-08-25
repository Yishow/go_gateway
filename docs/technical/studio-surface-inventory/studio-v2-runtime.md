# `/studio/v2` 與 `/studio/runtime` Inventory

## 這兩個 surface 的產品定位

| Surface | Route | 受眾 | 定位 |
| --- | --- | --- | --- |
| Studio V2 | `/studio/v2` | 一般使用者 | 簡單、可引導、好觀察的 setup flow |
| Runtime Dashboard | `/studio/runtime` | setup 後操作者 | focused post-setup runtime monitor |

## 目前策略

- `/studio/v2`：**這一輪重點施作**
- `/studio/runtime`：**作為 V2 的 post-setup 觀察面一起施作**
- 預設入口方向：**網站服務應直接切到 `/studio/v2`**
- commit 方向：**以多台設備為正式目標，不再先假設只有單台設備**

## `/studio/v2` 現況

### 核心事實

- 已有完整 shell、四步驟、settings、summary rail、Step 4 handoff。
- 各步驟皆已透過 autosave 持久化至後端 workspace endpoint（devices / source-rules / mappings / database-config / database-targets）；reducer 僅作 UI 狀態。
- `commit` 已是真實的 backend persisted lifecycle：Step 4 走 `POST /api/v1/datalink/studio-v2/workspace/activate`（真 probe、改 device active、灌進 scheduler 採集）。request 以 `workspace_revision`、`settings_revision`、`readiness_token` 與 autosave barrier 欄位作為啟動 gate；`commitLog.ts` / `dbReducer` 的舊模擬序列為 dead code，未被使用。
- Step 4 database-config 已支援 **DB 密碼**（端到端傳遞，更新留空則保留既有密碼）。
- 已支援 **自動建立目標資料表**：手動「建立資料表」按鈕（dry-run 預覽 + 執行，`POST /studio-v2/workspace/database-schema/generate`）＋ activate 時後端自動 ensure 保險。
- `UpsertTarget` 在目標表/欄位尚未建立時改為 degraded 放行（`AllowMissingTable`），由建表流程補建，解開「先綁 target 才能建表 / 先建表才能綁 target」死結。
- 產品決策已定為直接把服務預設入口切到 `/studio/v2`。
- `commit` 的正式方向已定為多台設備；現有 handoff 與 runtime route 仍偏單台心智。

### Modbus Share startup and projection transaction semantics

- Startup first loads persisted listener settings without binding. It then restores the canonical projection from persisted source-rule candidates under the workspace revision; only `ready` hydration with non-empty durable workspace/settings revisions, readiness, and readiness token permits the configured listener to bind. A settings/workspace/revision/restore failure leaves hydration failed or pending and the listener stopped. Uncertain rollback fails closed before `MarkDirty`; even a failed dirty-marker write cannot reopen listener/readiness.
- Readiness, projection, `runtime_apply_*`, and database-connector diagnostics use fixed safe operator messages that redact endpoint, credential, DSN, and transport details. Runtime SSE unavailable/error events remain typed and retryable; the runtime page preserves the last snapshot, polls every 5 seconds, and its reconnect action rebuilds the EventSource. Preview SSE retries are bounded and explicit disconnect cancels pending retries.
- Reconcile takes a workspace lock and then a process-global projection transaction lock because the memory bank and rollback snapshot are process-global. This serializes cross-workspace swaps to prevent one workspace rollback from restoring another workspace's newer projection. It trades some cross-workspace throughput for correctness; prevalidation and CAS gates run before the global lock where possible.
- The status response is workspace-scoped: `workspace_id`, `workspace_revision`, `settings_revision`, `listener_state`, hydration/readiness, and `mapping_count` are returned from the same hydrated workspace; foreign in-memory mappings are excluded. A ready status may include the server-owned canonical plan and desired mappings, each carrying full identity/geometry and a persisted `ownership_proof`. Dirty or uncertain recovery returns safe `recovery` metadata and keeps listener/readiness fail-closed.
- Candidate review/apply requests carry `workspace_id`, `expected_workspace_revision`, and the persisted source-rule `revision_id`. The production route rejects missing hydration, foreign devices, and stale revisions before reading or mutating candidate/tag/mapping state. Local Modbus candidate snapshots are the review truth; browser labels do not replace candidate status, blocking reason, or collision state.
- Register geometry is backend-canonical: human holding-register addresses and zero-based indexes are both exposed, datatype span must fit the requested stride, and capacity/range validation occurs before projection mutation. Reconcile reports `outcome`, applied/removed/invalidated counts and spans, plus operator-safe diagnostics; `dirty_unknown` and `invalidated_unknown` are not presented as ready.

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
| Step 2 Rule | 依設備編輯規則、生成 points 視圖 | local reducer state | 無 | `frontend-local-only` |
| Step 3 Mapping | 調整 mapping / transform | local reducer state | 無 | `frontend-local-only` |
| Step 4 Database | 調整 DB target（含 DB 密碼）、建立資料表、送出設定並啟動 | autosave PUT/POST 至 workspace endpoint | `database-config`（含 password）、`database-targets`、`database-schema/generate`、`POST /api/v1/datalink/studio-v2/workspace/activate` | `wired` |
| Settings | 編輯 connector pool、scheduler、modbus-share 設定並等待 durable revision | hydrated workspace/settings state | `GET /api/v1/datalink/settings`、`PUT /api/v1/datalink/settings/:key`、`GET /api/v1/datalink/modbus-share/status` | `wired` |
| Runtime handoff | 提交成功後前往 runtime dashboard | route navigation | workspace activation 已存在；正式多台結果 handoff 仍待補齊 | `semantics-gap` |

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
| 現有 query param | `device_id` |
| 現有 guard | 僅 UUID 形態才視為 persisted backend device id |
| 現有 fallback | 取不到 persisted id 時，退回 `/studio/runtime` 不帶 `device_id` |
| 目前限制 | 整個 handoff 仍以單台 persisted device 為前提 |
| 風險 | Step 4 commit 還沒真正產出 persisted device / point / runtime context，也無法表達多台成功結果 |

### V2 必補的正式契約

1. commit request 必須支援一次送出多台設備設定
2. commit response 必須逐台回傳 persisted device 結果
3. commit response 必須回 runtime 啟動狀態與最小 summary
4. runtime handoff 必須同時支援單台與多台結果，不能只靠單一 `device_id`
5. route handoff 不能再從 local draft state 猜 runtime device
6. `/` 與主要入口 redirect 應可直接落到 V2

## `/studio/runtime` 現況

### 使用者可見動作

- 目前仍是依 `device_id` 聚焦單台設備
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
| runtime snapshot | `useRuntimeStatus()` | `GET /api/v1/datalink/runtime/status?device_id=...` | `wired`；truthful top-level `running`, collector configured/status/running fields, nullable availability/last-read/last-error fields, and `modbus_share_delivery` diagnostics |
| runtime stream | `useRuntimeDashboardStream()` | `GET /api/v1/datalink/runtime/stream?device_id=...&point_ids=...` | `wired`；truthful value/status/heartbeat events and typed failure envelope；stream-state failures include `retryable`, opaque `request_id`, and actionable `action` |
| status events | SSE `status` | same stream endpoint | `wired` |
| heartbeat | SSE `heartbeat` | same stream endpoint | `wired` |
| Mapping preview stream | `usePreviewStream()` | `GET /api/v1/datalink/preview/stream?mapping_id=...&workspace_id=...` | `wired`；persisted mapping identity + required workspace ownership proof before subscription，typed `connected`/`preview`/`heartbeat`/`error` SSE events with `retryable`/opaque `request_id`/actionable `action` metadata，失敗不合成成功值；cross-workspace subscription is rejected before listener registration |
| Production target writer | embedded runtime target writer | collector target delivery fan-out | `wired`；database + authoritative Modbus Share projection only after hydration, records per point/tag delivery diagnostics; diagnostic sync endpoints are not the runtime delivery path |
| B10 EXE/TCP witness | `scripts/b10_exe_acceptance.py` | pure HTTP configure/save + EXE Modbus TCP read/restart and SQL/negative matrix | `13/13 PASS`；current pair `155805Z` + `155901Z`；focused Go `22/22`、Python `16/16`、`py_compile` PASS；SHA `8398ffdf4b8977848c7dffd22a7fd3df367643b2867df795d3a5f61c3d0b2de7`；155314Z failure is historical and its retained witness was deleted, so no artifact link is listed；external/embedded Chromium、Windows/LAN/real PLC/SCADA/field sign-off and deployment-owner acceptance remain pending |
| Modbus Share status | `useStudioV2ShareActivation()` | `GET /api/v1/datalink/modbus-share/status` | `wired`；workspace-scoped `mapping_count`、`listener_state`、revisions、canonical desired mappings/ownership proof 與 dirty recovery 由 backend 提供；bind failure 僅保存 safe typed code/retryable/action，success/disabled 清除；settings/reconcile 共用 lifecycle serialization/epoch，candidate snapshot 與 device membership 以 durable revision CAS 保護；hydration/settings/listener failures use typed envelope；disabled or non-ready workspaces stay fail-closed |
| Modbus Share candidate review | `useModbusShareCandidateReview()` | `GET /api/v1/datalink/source-rules/:id/candidates?workspace_id=...&expected_workspace_revision=...&revision_id=...` | `wired`；Step 4 以 persisted candidate snapshot 的 status/blocking reason/geometry 作為 review truth；scope、hydration、ownership、revision failure 不回退 browser mock |
| Modbus Share reconcile | `useStudioV2ShareActivation()` | `POST /api/v1/datalink/modbus-share/reconcile` | `wired`；server-owned canonical plan covers full mapping identity/geometry and signature，workspace/settings revisions plus `readiness_token` are required，stride/capacity/collision validation precedes mutation，response outcome/invalidated spans/diagnostics feeds the UI；candidate/apply/reconcile is authoritative，direct mapping mutation returns typed `modbus_share_projection_required`，typed 409/422 failures use `APIErrorResponse` |
| Modbus Share handler response contract | generated Swagger + `internal/api/modbus_share_swagger_contract_test.go` | Start/Stop、mapping query/mutations、reconcile | `wired`；annotations match only reachable typed `400/403/409/422/500/503` paths per operation，three generated Swagger artifacts are regenerated from the same source |
| Settings update | settings autosave | `PUT /api/v1/datalink/settings/:key` | `wired`；200 update or typed 409/422 listener bind/revision failures |
| Workspace activation | Step 4 activation controller | `POST /api/v1/datalink/studio-v2/workspace/activate` | `wired`；required JSON body with `workspace_revision`、`settings_revision`、`readiness_token` plus autosave barrier fields；typed 400/409/422/500 failures |

### runtime route state machine

| route state | 來源 | 備註 |
| --- | --- | --- |
| `missing-device-context` | 沒有可用的 selected device，或 typed code `runtime_device_not_found` | 不依賴 raw backend error string；顯示回 Studio V2 的 context action |
| `loading` | 尚未拿到首個 snapshot | 正常首屏 |
| `live` | snapshot 已就緒且 SSE connected | 正常監看 |
| `degraded` | 有最後成功 snapshot，但 SSE error/unavailable | 保留 snapshot、停止 live tone，退回 5 秒 polling |
| `reconnecting` | retryable SSE unavailable/error 且操作者要求 reconnect | 實際重建 EventSource；成功後回 `live`，不合成成功值 |
| `error` | snapshot 失敗且非 `runtime_device_not_found` | 顯示 typed code 的 localized copy、opaque request ID 與 retry action |
| `empty` | snapshot `snapshot_state.empty=true` | 顯示等待/空資料，不合成 ready dashboard |

### Runtime 契約缺口

| 類型 | 問題 | 影響 |
| --- | --- | --- |
| typed-error | snapshot/stream failure 使用 `error.code`, `error.message`, `error.retryable`, `error.request_id` | 前端以 code 映射 en/zh-TW safe copy，raw backend detail 不進正常 dashboard DOM |
| missing-api | runtime handoff 還沒有正式的多台結果參數或上下文 | V2 commit 成功後無法穩定交接多台設備結果 |
| semantics-gap | 缺 `starting` / `committed but not running` / `not_committed` 類狀態 | setup 後無法正確解釋等待期 |
| semantics-gap | runtime 與 setup commit 沒有正式 lifecycle handoff | V2 不能穩定直達 focused runtime |
| exists-not-wired | runtime snapshot 只有 post-setup route 用 | `/studio/v2` 內部還不能當正式 commit result source |

## 使用者期望與現況的差距

### 使用者想要的 V2

- 簡單
- 有引導
- 容易觀察
- 作為服務預設入口
- 支援多台設備一起送出

### 目前缺的不是 UI，而是 contract

1. V2 沒有正式的多台設備 commit API
2. V2 沒有逐台 persisted output summary
3. Runtime 還沒有正式接住多台 handoff
4. Runtime 沒有產品化 lifecycle semantics
5. Router 還沒改成以 V2 為預設入口

## V2 近期應優先完成的順序

1. `make-studio-v2-default-entry`
   - 直接把產品入口切到 `/studio/v2`
2. `integrate-studio-v2-multi-device-commit`
   - 讓 V2 不再只是 local-only setup shell，且能一次送出多台設備
3. `adapt-runtime-handoff-for-multi-device`
   - 讓 runtime 接得住單台或多台 commit 結果
4. `normalize-runtime-lifecycle-status-contract`
   - 讓 runtime 成為可靠的 post-setup observer

## V2 暫時不該膨脹成的樣子

1. 不要變成第二個 `/studio` 完整版。
2. 不要把 fleet-wide dashboard、全局 review desk、工程工具全塞進來。
3. 不要在 commit contract 未落地前，繼續用更多模擬 log 或 reducer state 掩蓋真實 lifecycle 缺口。
