## 1. Step 4 database autosave

- [x] 1.1 交付 `Valid-only database autosave`：讓 Step 4 connector / target 編輯只有在內容合法時才 autosave，並以前端 Step 4 測試與後端 validation test 驗證不合法內容不覆蓋最後成功版本。
- [x] 1.2 交付 `Per-target save isolation`：讓 database target row 逐項保存、逐項顯示錯誤，不因單項失敗阻塞其它 target，並以前端 integration test 驗證 partial success。

## 2. 啟動邊界

- [x] 2.1 依照 `### Database autosave remains pre-activation` 固定 Step 4 autosave 成功不會直接啟動 runtime 的 contract，並以後端 service test 與前端 UI test 驗證 save 成功後 scheduler 狀態不自動改變。
- [x] 2.2 依照 `### Valid-only database autosave` 建立 connector metadata 與 target rows 的分離 save state，並以前端 state test 驗證 metadata 成功不會掩蓋 target row 失敗。

## 3. API 契約

- [x] 3.1 交付 `Studio v2 workspace database autosave APIs`：提供 workspace-scoped database config / target API，並以 router / handler test 驗證 payload、response 與 validation error。
- [x] 3.2 依照 `### Per-target save isolation` 固定 target row validation 與錯誤訊息，不讓單項失敗演變成整批失敗，並以 backend service test 驗證。
