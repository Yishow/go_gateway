# Backend API Registry For Frontend Surfaces

## 使用方式

這份表不是完整 Swagger 替代品，而是「哪些 backend APIs 對目前前端 surface 有關」的維護索引。  
欄位中的 `主用頁面` 指的是目前實際或預期會依賴這組 API 的前端 surface。  
這一輪解讀時請優先看 `/studio/v2` 與 `/studio/runtime`。表格中只標示 `/studio` 的項目是 legacy pre-delete historical owner，不代表刪除 backend API；API 是否仍由 V2、runtime 或其他保留 surface 使用，必須回讀現行 source。`/studio` 不再是 fallback，也不得因刪除 dedicated frontend route 而猜測刪除 backend domain capability。

## This change: affected contract

- The dedicated `/studio` frontend surface was removed in Git commit `efa355f8c7943941b246a750cf7650b5ba211f88`; `/studio` labels below are historical pre-delete ownership, not active route registrations.
- `/studio/v2`, `/studio/runtime`, `/test`, and `/gateway/*` retain their route identities. Runtime APIs remain backend domain capabilities and are actively owned by `/studio/runtime`.
- Typed failures use `handlers.APIErrorResponse` with `error.code`, `error.message`, `error.retryable`, and opaque `error.request_id`. Runtime codes are `runtime_device_not_found`, `runtime_snapshot_unavailable`, and `runtime_stream_unavailable`; preview uses `preview_invalid_request`, `preview_unavailable`, and fail-closed workspace scope; readiness and activation use `workspace_not_ready`, `readiness_blocked`, `activation_request_invalid`, `activation_failed`, and the Modbus Share save/revision codes. Runtime status additionally exposes latest `modbus_share_delivery` diagnostics; nullable availability/last-read/last-error fields remain nullable. Modbus Share projection/status/settings routes are documented in the output section below.
- Startup semantics are deliberately ordered: durable listener settings hydrate without binding, persisted source-rule candidates restore the canonical workspace projection, and only a ready hydration state with non-empty workspace/settings revisions, readiness, and readiness token may bind the configured listener. Failed settings/workspace/revision/projection hydration remains failed/stopped; no defaults or fixed-port fallback are reported as ready. Uncertain rollback fails closed before attempting the durable dirty marker; a `MarkDirty` failure does not reopen listener/readiness.
- Readiness, projection, `runtime_apply_*`, and connector messages are normalized operator-safe copy with actionable categories and no endpoint, credential, DSN, or transport detail. Runtime SSE failures retain typed action metadata; runtime snapshot fallback polls every 5 seconds, and its reconnect control creates a new EventSource. Preview reconnect attempts are bounded and cancelled by explicit disconnect.
- Reconcile uses a per-workspace lock plus a process-global projection transaction lock because the memory bank and rollback snapshot are process-global. Cross-workspace projection transactions therefore serialize for rollback isolation (a correctness/throughput trade-off); prevalidation and CAS occur before the global lock where possible.
- Current B10 EXE/TCP acceptance pair is `b10-exe-acceptance-20260825T155805Z.json` and `b10-exe-acceptance-20260825T155901Z.json`, both 13/13 PASS with focused Go 22/22, Python 16/16, `py_compile` PASS, and binary SHA-256 `8398ffdf4b8977848c7dffd22a7fd3df367643b2867df795d3a5f61c3d0b2de7`. The 155314Z failure and earlier pairs are historical only; the 155314Z witness was removed by retention and is not linked here. Windows, LAN, real PLC, SCADA, external/embedded Chromium, field sign-off, and deployment-owner witnesses remain pending.
- Review18 activation contract: `POST /api/v1/datalink/studio-v2/workspace/activate` requires a JSON body with non-empty `workspace_revision`, `settings_revision`, and `readiness_token`; generated Swagger and Go binding tags are synchronized. Missing body is typed `422` (`modbus_share_save_incomplete`), missing required fields are typed `400` (`activation_request_invalid`), pending/error saves remain typed `422`, and stale revisions remain typed `409`.
- Review20 Modbus Share handler contract: generated Swagger and the handler annotation contract test are synchronized with reachable typed responses. Start is `400/409/422/503`; Stop is `400/409/422/500/503`; workspace mapping list is `403/422/503`; mapping upsert is `400/403/422/503`; mapping delete is `403/422/503`; reconcile is `400/403/409/422/503`. Statuses not reachable from the typed gate/handler paths are intentionally omitted.
- Review22 contract refresh: candidate review/apply/recompute and Modbus Share status/reconcile now document workspace/revision scope, workspace-only mapping counts, persisted ownership proof, stride-versus-span rejection, invalidated/dirty outcomes, and the generated `modbusshare.OwnershipProof` schema. Lifecycle settings/reconcile share serialization and epoch guards (a staged reconcile after disable is rejected without durable/runtime mutation); status exposes safe typed bind diagnostics and clears them on success/disable; candidate snapshots use repository revision CAS; device attach/detach advances the durable workspace scope through the existing revision store and rolls back the workspace record on failure. The three Swagger artifacts remain synchronized with zero JSON/YAML drift; current B10 is the 155805Z/155901Z PASS pair.

## Datalink Core

| API | 主用頁面 | 用途 |
| --- | --- | --- |
| `GET /api/v1/datalink/devices` | `/studio`, `/studio/runtime` | 列出設備、解析 focused device |
| `POST /api/v1/datalink/devices` | `/studio` | 建立設備 |
| `PUT /api/v1/datalink/devices/:id` | `/studio` | 更新設備 |
| `DELETE /api/v1/datalink/devices/:id` | `/studio` | 刪除設備 |
| `POST /api/v1/datalink/devices/test-draft` | `/studio` | 測試 draft 設備連線 |
| `POST /api/v1/datalink/devices/:id/test` | `/studio` | 測試 persisted 設備連線 |
| `POST /api/v1/datalink/devices/:id/readiness` | 未全面接線 | 檢查設備 readiness |
| `POST /api/v1/datalink/devices/:id/activate` | `/studio` | 啟用設備 / collector |
| `POST /api/v1/datalink/devices/:id/disable` | `/studio` | 停用設備 / collector |
| `GET /api/v1/datalink/protocols` | 預期 `/studio` | 列 protocol 清單 |

## Source / Point / Rule

| API | 主用頁面 | 用途 |
| --- | --- | --- |
| `GET /api/v1/datalink/source-rules` | `/studio` | 列出設備的 source rules |
| `POST /api/v1/datalink/source-rules` | `/studio` | 建立 source rule |
| `PUT /api/v1/datalink/source-rules/:id` | `/studio` | 更新 source rule |
| `DELETE /api/v1/datalink/source-rules/:id` | `/studio` | 刪除 source rule |
| `POST /api/v1/datalink/source-rules/:id/enable` | `/studio` | 啟用 rule |
| `POST /api/v1/datalink/source-rules/:id/disable` | `/studio` | 停用 rule |
| `GET /api/v1/datalink/points` | `/studio`, `/studio/runtime` | 列 points |
| `POST /api/v1/datalink/points` | `/studio` | 建立 point |
| `DELETE /api/v1/datalink/points/:id` | `/studio` | 刪除 point |
| `POST /api/v1/datalink/points/:id/poll` | 預期 `/studio` | 單點即時輪詢 |
| `POST /api/v1/datalink/points/poll` | 預期 `/studio` | 批次輪詢 |

## Tag / Mapping / Review

| API | 主用頁面 | 用途 |
| --- | --- | --- |
| `GET /api/v1/datalink/tags` | `/studio` | 列出 tags |
| `POST /api/v1/datalink/tags` | `/studio` | 建立 tag |
| `DELETE /api/v1/datalink/tags/:id` | `/studio` | 刪除 tag |
| `POST /api/v1/datalink/tags/batch` | `/studio` | 批次建立 tags |
| `POST /api/v1/datalink/tags/validate-key` | 預期 `/studio` | 驗證 key |
| `GET /api/v1/datalink/mappings` | `/studio` | 列出 mappings |
| `POST /api/v1/datalink/mappings` | `/studio` | 建立 mapping |
| `PUT /api/v1/datalink/mappings/:id` | `/studio` | 更新 mapping |
| `DELETE /api/v1/datalink/mappings/:id` | `/studio` | 刪除 mapping |
| `POST /api/v1/datalink/mappings/preview` | `/studio/v2` Step 3 | pipeline 預覽；typed failures: `preview_invalid_request`, `preview_unavailable` |
| `POST /api/v1/datalink/mappings/validate-pipeline` | 預期 `/studio` | pipeline 驗證 |
| `GET /api/v1/datalink/source-rules/:id/candidates` | `/studio/v2` Step 4 | 讀 persisted candidate snapshot；production request 必須帶 `workspace_id`、`expected_workspace_revision`、`revision_id`，回傳 candidate set status/blocking reason 與 Local Modbus candidate geometry |
| `POST /api/v1/datalink/source-rules/:id/candidates/recompute` | `/studio/v2` Step 4 | 以同一 workspace/revision scope 重算 persisted candidate snapshot；hydration、ownership 或 revision 不符時 fail closed |
| `GET /api/v1/datalink/source-rules/:id/tag-review-decisions` | `/studio/v2` Step 4 | 讀 workspace/revision-scoped tag review decisions |
| `POST /api/v1/datalink/source-rules/:id/tag-review-decisions` | `/studio/v2` Step 4 | 寫入帶 `workspace_id`、`expected_workspace_revision`、`revision_id` 的 tag review decision |
| `POST /api/v1/datalink/source-rules/:id/tags/apply` | `/studio/v2` Step 4 | 將 persisted review 決策正式套用；同一 candidate scope gate 與 typed workspace/revision failures 生效 |

## Runtime / Monitoring

| API | 主用頁面 | 用途 |
| --- | --- | --- |
| `GET /api/v1/datalink/runtime/status` | `/studio/runtime` | truthful runtime snapshot；top-level `running`、collector `status`/`running`/availability fields、nullable availability/last-read/last-error fields，以及 per point/tag `modbus_share_delivery` diagnostics；typed failure: `runtime_snapshot_unavailable` |
| `GET /api/v1/datalink/runtime/stream` | `/studio/runtime` | truthful value / status / heartbeat SSE；stream-state failures include typed code, `retryable`, opaque `request_id`, and actionable `action`; typed failures: `runtime_device_not_found`, `runtime_stream_unavailable` |
| `GET /api/v1/datalink/studio-v2/workspace/runtime-context` | `/studio/runtime` | persisted device/setup context；200 response plus typed 404/503 failures (`runtime_device_not_found`, `workspace_not_ready`, `runtime_snapshot_unavailable`) |
| `POST /api/v1/datalink/studio-v2/workspace/activate` | `/studio/v2` Step 4 | persisted workspace activation；required JSON body carries non-empty `workspace_revision`, `settings_revision`, and `readiness_token`, plus optional `pending_saves` and `save_error`; typed 400/409/422/500 failures |
| `GET /api/v1/datalink/dashboard/stats` | 舊 dashboard 契約 | fleet 型統計 |
| `GET /api/v1/datalink/dashboard/device-statuses` | 舊 dashboard 契約 | fleet 型裝置狀態 |
| `GET /api/v1/datalink/preview/stream` | `/studio/v2` Step 3 | persisted mapping 的 truthful runtime-value SSE；required `mapping_id` + `workspace_id` ownership proof before subscription；`connected`/`preview`/`heartbeat`/`error` events carry typed code, `retryable`, opaque `request_id`, and actionable `action` metadata，不合成成功值；typed failures: `preview_invalid_request`, `preview_unavailable`, `modbus_share_workspace_scope` |

## Output: Local Modbus

| API | 主用頁面 | 用途 |
| --- | --- | --- |
| `GET /api/v1/datalink/modbus-share/status` | `/studio/v2` Settings / Step 4 | 讀 workspace-scoped status：`workspace_id`、`workspace_revision`、`settings_revision`、`listener_state`、hydration/readiness、workspace mapping count、canonical plan/desired mappings、ownership proof 與 dirty recovery；hydration 未 ready 或 global gate 關閉時 fail-closed；200 response 或 typed 503 `APIErrorResponse` |
| `POST /api/v1/datalink/modbus-share/start` | `/studio/v2` Settings | 依持久化 listener settings 啟動 share；hydration/global gate 與 typed listener diagnostics 生效，未 ready 不 bind |
| `POST /api/v1/datalink/modbus-share/stop` | `/studio/v2` Settings | 依持久化 listener lifecycle 停止 share；hydration/global gate 與 typed listener diagnostics 生效 |
| `GET /api/v1/datalink/modbus-share/mappings` | `/studio/v2` Step 4 | 讀 workspace-scoped backend projection；mapping count 與 rows 不包含其他 workspace，hydration/global gate 與 typed diagnostics 生效 |
| `PUT /api/v1/datalink/modbus-share/mappings/:tagId` | `/studio/v2` Step 4 | legacy direct mapping mutation；422 `modbus_share_projection_required` diagnostic，正式變更必須走 candidate/apply/reconcile |
| `DELETE /api/v1/datalink/modbus-share/mappings/:tagId` | `/studio/v2` Step 4 | legacy direct mapping mutation；403 ownership 或 422 `modbus_share_projection_required` typed diagnostic，正式變更必須走 candidate/apply/reconcile |
| `POST /api/v1/datalink/modbus-share/write-tag-value` | `/studio/v2` Step 4 diagnostics | 寫測試值的 gated diagnostic path；不得取代 candidate/apply/reconcile projection |
| `POST /api/v1/datalink/modbus-share/sync` | `/studio/v2` Step 4 diagnostics | 將既有 mapping values 同步至 share 的 gated diagnostic path；正式 projection 仍由 reconcile 負責 |
| `POST /api/v1/datalink/modbus-share/reconcile` | `/studio/v2` Step 4 | server-owned canonical plan validates full mapping identity/geometry and signature (`workspace_id`, revisions, `mapping_id`, `tag_id`, source-rule revision, datatype/span/stride/capacity); stride smaller than datatype span is rejected; response reports `outcome`, applied/removed/invalidated spans and safe diagnostics; stale revisions use CAS and return typed 409/422 failures (`modbus_share_revision_conflict`, `modbus_share_reconcile_failed`, `modbus_share_dirty_unknown`, `modbus_share_projection_required`) |
| `POST /api/v1/datalink/source-rules/:id/local-modbus/apply` | `/studio/v2` Step 4 | 將 scoped persisted Local Modbus candidate review 套用為 output；正式 projection 仍只能由 candidate/apply/reconcile seam 產生 |

## Output: Database

| API | 主用頁面 | 用途 |
| --- | --- | --- |
| `GET /api/v1/datalink/db-targets/connectors` | `/studio` | 列 DB connectors |
| `POST /api/v1/datalink/db-targets/connectors` | `/studio` | 建 connector |
| `PUT /api/v1/datalink/db-targets/connectors/:id` | `/studio` | 更新 connector |
| `DELETE /api/v1/datalink/db-targets/connectors/:id` | `/studio` | 刪 connector |
| `POST /api/v1/datalink/db-targets/connectors/:id/test` | `/studio` | 測試 connector |
| `GET /api/v1/datalink/db-targets/connectors/:id/tables` | `/studio` | 列 tables |
| `GET /api/v1/datalink/db-targets/connectors/:id/validate` | `/studio` | 驗證 connector |
| `POST /api/v1/datalink/db-targets/connectors/:id/schema/generate` | `/studio` | 產 schema |
| `POST /api/v1/datalink/db-targets/connectors/:id/mappings/dry-run` | `/studio` | dry-run mappings |
| `GET /api/v1/datalink/db-targets/mappings` | `/studio` | 列 output mappings |
| `POST /api/v1/datalink/db-targets/mappings` | `/studio` | 建 mapping |
| `PUT /api/v1/datalink/db-targets/mappings/:id` | `/studio` | 更新 mapping |
| `DELETE /api/v1/datalink/db-targets/mappings/:id` | `/studio` | 刪 mapping |
| `POST /api/v1/datalink/source-rules/:id/database-outputs/apply` | 預期 `/studio` | 將 rule review 套用為 DB output |

## Settings / Flags

| API | 主用頁面 | 用途 |
| --- | --- | --- |
| `GET /api/v1/datalink/settings` | `/studio/v2` Settings、`/gateway/*` | 讀 feature flags / system settings；V2 為現行 owner |
| `PUT /api/v1/datalink/settings/:key` | `/studio/v2` Settings | durable settings / Modbus Share listener；200 update response or typed 409/422 failures (`settings_invalid`, `settings_update_failed`, `modbus_share_revision_conflict`, `modbus_share_listener_bind_failed`) |

## Engineering / Test APIs

| API | 主用頁面 | 用途 |
| --- | --- | --- |
| `POST /api/v1/test/connect` | `/gateway/quick-setup`, `/gateway/expert-workbench`, `/test` | connect-only 測試 |
| `POST /api/v1/test/disconnect` | `/test` | 中斷測試連線 |
| `GET /api/v1/test/status` | `/test` | 讀測試連線狀態 |
| `POST /api/v1/test/read` | `/test` | 單次讀取 |
| `POST /api/v1/test/write` | `/test` | 單次寫入 |
| `POST /api/v1/test/batch` | `/test` | 批次測試 |
| `POST /api/v1/test/monitor/start` / `stop` | `/test` | 啟停 monitor |
| `GET /api/v1/debug/packets` / `logs` | `/test` | debug 視圖 |
| `DELETE /api/v1/debug/clear` | `/test` | 清 debug 資料 |

## 這份 registry 要怎麼維護

1. 新頁面接新 API 時，同步把 `主用頁面` 改掉。
2. API 已存在但尚未前端使用時，保留在 registry 並在對應 surface 文件標記 `exists-not-wired`。
3. 如果 page 開始依賴錯誤字串或隱含語意，不能只改前端，必須在 gap roadmap 裡補一條契約 change。
