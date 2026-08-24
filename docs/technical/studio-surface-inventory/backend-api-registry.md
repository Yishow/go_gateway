# Backend API Registry For Frontend Surfaces

## 使用方式

這份表不是完整 Swagger 替代品，而是「哪些 backend APIs 對目前前端 surface 有關」的維護索引。  
欄位中的 `主用頁面` 指的是目前實際或預期會依賴這組 API 的前端 surface。  
這一輪解讀時請優先看 `/studio/v2` 與 `/studio/runtime`。表格中只標示 `/studio` 的項目是 legacy pre-delete historical owner，不代表刪除 backend API；API 是否仍由 V2、runtime 或其他保留 surface 使用，必須回讀現行 source。`/studio` 不再是 fallback，也不得因刪除 dedicated frontend route 而猜測刪除 backend domain capability。

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
| `POST /api/v1/datalink/mappings/preview` | 預期 `/studio` | pipeline 預覽 |
| `POST /api/v1/datalink/mappings/validate-pipeline` | 預期 `/studio` | pipeline 驗證 |
| `GET /api/v1/datalink/source-rules/:id/candidates` | `/studio` | review candidate snapshot |
| `POST /api/v1/datalink/source-rules/:id/candidates/recompute` | 預期 `/studio` | 重算候選項 |
| `GET /api/v1/datalink/source-rules/:id/tag-review-decisions` | `/studio` | 讀 tag review decisions |
| `POST /api/v1/datalink/source-rules/:id/tag-review-decisions` | `/studio` | 寫 tag review decision |
| `POST /api/v1/datalink/source-rules/:id/tags/apply` | 預期 `/studio` | 將 review 決策正式套用 |

## Runtime / Monitoring

| API | 主用頁面 | 用途 |
| --- | --- | --- |
| `GET /api/v1/datalink/runtime/status` | `/studio`, `/studio/runtime` | runtime snapshot |
| `GET /api/v1/datalink/runtime/stream` | `/studio`, `/studio/runtime` | value / status / heartbeat SSE |
| `GET /api/v1/datalink/dashboard/stats` | 舊 dashboard 契約 | fleet 型統計 |
| `GET /api/v1/datalink/dashboard/device-statuses` | 舊 dashboard 契約 | fleet 型裝置狀態 |
| `GET /api/v1/datalink/preview/stream` | 預覽流 | 頁面實際使用需另整理 |

## Output: Local Modbus

| API | 主用頁面 | 用途 |
| --- | --- | --- |
| `GET /api/v1/datalink/modbus-share/status` | `/studio` | 讀 local modbus share 狀態 |
| `POST /api/v1/datalink/modbus-share/start` | `/studio` | 啟動 share |
| `POST /api/v1/datalink/modbus-share/stop` | `/studio` | 停止 share |
| `GET /api/v1/datalink/modbus-share/mappings` | `/studio` | 讀 mapping |
| `PUT /api/v1/datalink/modbus-share/mappings/:tagId` | `/studio` | 設定 mapping |
| `DELETE /api/v1/datalink/modbus-share/mappings/:tagId` | `/studio` | 刪除 mapping |
| `POST /api/v1/datalink/modbus-share/write-tag-value` | `/studio` | 寫測試值 |
| `POST /api/v1/datalink/modbus-share/sync` | `/studio` | 將 mapping 同步至 share |
| `POST /api/v1/datalink/source-rules/:id/local-modbus/apply` | 預期 `/studio` | 將 rule review 套用為 modbus output |

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
| `GET /api/v1/datalink/settings` | `/gateway/*`、未來可供 `/studio/v2` | 讀 feature flags / system settings |
| `PUT /api/v1/datalink/settings/:key` | 系統設定頁潛在使用 | 改設定 |

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
