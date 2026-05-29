# `/test` Tooling Inventory

## 定位

- Route: `/test`
- 受眾: 工程測試使用者
- 定位: 協議測試 / debug / 工程工具
- 目前策略: **先暫停**

## 暫停說明

`/test` 對你很重要，但它不是這一輪主產品交付的阻塞點。  
因此這份文件會完整記錄它目前能做什麼、依賴哪些 API、還有哪些未完全接齊的工具面，但不把它列成近期主要 change 目標。

## 使用者可見動作

- connect / disconnect
- get status
- read / write
- batch request
- start / stop monitor
- refresh / clear debug packets and logs
- profile / config 切換與暫存

## 主要前端實作

- `frontend/src/pages/TestPage.tsx`
- `frontend/src/services/api.ts`
- `frontend/src/components/ConfigForm.tsx`
- `frontend/src/components/TestOperations.tsx`
- `frontend/src/components/MonitorControl.tsx`
- `frontend/src/components/DebugPanel.tsx`

## 對應 API

| 類型 | API | 狀態 | 備註 |
| --- | --- | --- | --- |
| connect | `POST /api/v1/test/connect` | `wired` | 也是 `/gateway/*` prototype submit 的底層依賴 |
| disconnect | `POST /api/v1/test/disconnect` | `wired` | 工程斷線 |
| status | `GET /api/v1/test/status` | `wired` | 測試連線狀態 |
| read | `POST /api/v1/test/read` | `wired` | 單次讀取 |
| write | `POST /api/v1/test/write` | `wired` | 單次寫入 |
| batch | `POST /api/v1/test/batch` | `wired` | 批次操作 |
| script | `POST /api/v1/test/script` | `exists-not-wired` | API 存在，頁面實際用法待後續整理 |
| list scripts | `GET /api/v1/test/scripts` | `exists-not-wired` | 同上 |
| save script | `POST /api/v1/test/scripts` | `exists-not-wired` | 同上 |
| delete script | `DELETE /api/v1/test/scripts/:id` | `exists-not-wired` | 同上 |
| monitor start | `POST /api/v1/test/monitor/start` | `wired` | 工程監看 |
| monitor stop | `POST /api/v1/test/monitor/stop` | `wired` | 工程監看 |
| monitor stream | `GET /api/v1/test/monitor/stream` | `exists-not-wired` | 頁面目前未全面盤進此流 |
| debug packets | `GET /api/v1/debug/packets` | `wired` | debug 視圖 |
| debug logs | `GET /api/v1/debug/logs` | `wired` | debug 視圖 |
| debug clear | `DELETE /api/v1/debug/clear` | `wired` | 清理 debug buffer |
| send raw | `POST /api/v1/debug/send-raw` | `exists-not-wired` | API 存在，主要給工程調試 |

## `/test` 對主產品的間接影響

1. `/gateway/quick-setup` 與 `/gateway/expert-workbench` 目前提交時其實是打 `/test/connect`。
2. 因此 `/gateway/*` 若不轉正，`/test` 就仍是它們的底層依賴。
3. 但 `/studio` / `/studio/v2` 不應以 `/test/*` 取代正式 datalink persisted contract。

## 後續再回來做時建議補的面

1. script / monitor stream 的頁面 wiring
2. 工程 profile / preset / history 的統一治理
3. 更明確的 protocol-specific diagnostics
4. 與 `/gateway/*` 的關係切清楚：共用 tooling 還是正式產品 contract

## 暫停期間維護重點

1. 若新增 `/test/*` API，仍應更新這份文件。
2. 若 `/gateway/*` 繼續依賴 `/test/connect`，要在這份文件保留註記。
3. 不要把 `/test` 的工程便利能力直接視為 `/studio/v2` 可依賴的正式產品契約。
