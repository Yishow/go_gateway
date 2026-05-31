## 1. Step 3 mapping autosave

- [x] 1.1 交付 `Valid-only mapping autosave`：讓每列 mapping 只有在內容合法時才 autosave 到後端，並以前端 Step 3 測試與後端 validation test 驗證不合法列不覆蓋最後成功版本。
- [x] 1.2 交付 `Per-point mapping save isolation`：讓 mapping row 逐列保存、逐列顯示錯誤，不因一列失敗阻塞其它列，並以前端 integration test 驗證 partial success。

## 2. Save state 呈現

- [x] 2.1 依照 `### Save-state visibility` 為每列建立 `local_value`, `persisted_value`, `save_state` 契約，並以前端 state test 驗證 reload 後仍能區分本地未存與後端已存版本。
- [x] 2.2 依照 `### Per-point mapping save isolation` 固定 save-error 不會清空使用者本地輸入，並以前端 row editor test 驗證錯誤後仍保留本地值。

## 3. API 契約

- [x] 3.1 交付 `Studio v2 workspace mapping autosave APIs`：提供 workspace-scoped mapping list/create/update/delete API，並以 router / handler test 驗證 payload、response 與 validation error。
- [x] 3.2 依照 `### Valid-only mapping autosave` 把既有 mapping validation 邊界接到 autosave contract，並以 backend service test 驗證合法/不合法分流。
