## 1. Step 2 規則 autosave

- [ ] 1.1 交付 `Valid-only rule autosave`：讓 Step 2 規則只有在內容合法時才觸發 save，並以前端 Step 2 測試與後端 validation test 驗證不合法編輯不覆蓋最後成功版本。
- [ ] 1.2 交付 `Per-rule save isolation`：讓多條規則逐條保存、逐條顯示錯誤，不因單條失敗擋住其它合法規則，並以前端 integration test 驗證 partial success。

## 2. Workspace 與 device 關聯

- [ ] 2.1 交付 `Workspace-scoped rule ownership`：讓每條規則都能重建 `workspace_id` 與 `device_id` 關係，並以後端 service test 與前端 reload test 驗證重整後 ownership 不漂移。
- [ ] 2.2 依照 `### Workspace-scoped rule ownership` 固定跨設備 rule save boundary，不讓 rule 在 save 時偷偷改掛到別的 device，並以 handler test 驗證 ownership mismatch error。

## 3. API 契約

- [ ] 3.1 交付 `Studio v2 workspace rule autosave APIs`：提供 workspace-scoped rule list/create/update/delete API，並以 router / handler test 驗證 payload、response 與 validation error。
- [ ] 3.2 依照 `### Valid-only rule autosave` 固定前端每條規則的 `save_state` 呈現，並以前端 state test 驗證 `saving`、`saved`、`save-error` 不互相污染。
