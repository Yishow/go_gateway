## 1. Step 1 自動存

- [x] 1.1 交付 `Valid-only device autosave`：讓 Step 1 只有在設備表單合法時才觸發 create / update，並以前端 `step1` 測試與後端 handler test 驗證不合法輸入不會覆蓋後端最後成功版本。
- [x] 1.2 交付 `Per-device save isolation`：讓多台設備逐台保存、逐台顯示失敗，不因一台錯誤阻塞其它合法設備，並以前端 integration test 驗證單台失敗時其它合法設備仍可保存。

## 2. 排序與刪除

- [x] 2.1 交付 `Studio v2 workspace device order API` 並依照 `### Workspace-owned order` 提供 workspace device order 更新契約，讓 V2 排列順序可持久化，並以 backend service test 與前端 route reload test 驗證重整後順序一致。
- [x] 2.2 交付 `Persisted device deletion cascade`：刪除已保存設備時同步刪除該設備的 V2 關聯資料，並以後端 cascade test 驗證設備、workspace 關聯與衍生 V2 資料不再存在。

## 3. API 契約

- [x] 3.1 交付 `Studio v2 workspace device autosave APIs`：提供 workspace-scoped device list/create/update/delete API，並以 router / handler test 驗證路徑、payload 與錯誤訊息。
- [x] 3.2 依照 `### Valid-only device autosave` 固定前端 save state 與錯誤標示，不讓 save-error 把使用者輸入清掉，並以前端 state test 驗證 `draft-invalid`、`saving`、`saved`、`save-error` 轉換。
