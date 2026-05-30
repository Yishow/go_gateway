## 1. 第一次啟動

- [ ] 1.1 交付 `First activation targets all eligible devices`：讓 Step 4 只啟動 singleton workspace 內所有合法、可用、尚未啟動的設備，並以 backend activation test 驗證不會重覆啟動已在跑設備。
- [ ] 1.2 交付 `Activation returns per-device results`：讓 activation API 逐台回傳成功/失敗與訊息，並以前端 Step 4 test 與 handler test 驗證 partial success 結果可見。

## 2. 導頁與邊界

- [ ] 2.1 交付 `Navigation remains available after partial failure`：讓至少有一台成功時仍可前往 runtime，並以前端 integration test 驗證部分失敗不阻止導頁。
- [ ] 2.2 依照 `### First activation targets all eligible devices` 固定「無可啟動設備」的回應與畫面提示，並以 backend service test 與 UI test 驗證空 activation 結果可行動。

## 3. API 與文案

- [ ] 3.1 交付 `Studio v2 workspace first activation API`：提供 `POST /api/v1/datalink/studio-v2/workspace/activate` 契約，並以 router / handler test 驗證 payload、response 與 partial failure 行為。
- [ ] 3.2 依照 `### Activation returns per-device results` 改寫 Step 4 最後按鈕與成功/結果卡語意，不再沿用 commit 文案，並以前端 snapshot / interaction tests 驗證。
