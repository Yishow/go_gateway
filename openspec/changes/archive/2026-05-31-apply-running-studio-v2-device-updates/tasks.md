## 1. 執行中設備直接套用

- [x] 1.1 交付 `Running device update application`：讓已啟動設備在合法 autosave 成功後直接套用新設定，並以 backend runtime test 驗證 save 後不需再按 Step 4。
- [x] 1.2 交付 `Not-running devices stay dormant`：讓尚未第一次啟動的設備 autosave 成功後仍維持未啟動，並以 backend service test 驗證 autosave 不會偷偷把設備切成 running。

## 2. 回應契約

- [x] 2.1 交付 `Studio v2 autosave responses expose runtime apply state` 與 `Autosave response exposes runtime apply state`：在 Step 1~4 autosave success payload 中加入 `runtime_apply_status`，並以 handler test 與前端 hook test 驗證 `not_running`, `applied`, `apply_failed` 三種結果可被辨識。
- [x] 2.2 依照 `### Running device update application` 固定 apply 失敗時的錯誤表達，不讓 UI 把 apply failure 誤認成單純 save success，並以前端 state test 驗證錯誤狀態保留。
