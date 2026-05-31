## 1. Diagnostics 契約

- [x] 1.1 交付 `Studio v2 Step 1 live diagnostics API`：讓 V2 可用真實 backend request 對 current draft 執行 diagnostics，並以 handler / router test 驗證 success、connect failure、probe failure。
- [x] 1.2 依照 `### Backend-produced diagnostics only` 固定 response shape，讓前端可直接渲染 connect / probe outcome，而不是自行合成 stage 成功。

## 2. 前端接線

- [x] 2.1 交付 `Replace Step 1 mock test animation`：移除 operator-visible 的 Step 1 setTimeout / Math.random success synthesis，改用 hook / mutation 驅動真實 diagnostics，並以前端測試驗證 request 發送與 loading state。
- [x] 2.2 交付 `Diagnostics-driven continue gate`：只有 backend diagnostics success 才能讓 device 進入 tested state 與通過 continue gate，並以前端 state / integration test 驗證失敗不會 auto-pass。

## 3. 驗證

- [x] 3.1 補齊 backend tests、frontend tests、手動 demo，證明沒有真實 diagnostics success 時 Step 1 不能宣告通過。
