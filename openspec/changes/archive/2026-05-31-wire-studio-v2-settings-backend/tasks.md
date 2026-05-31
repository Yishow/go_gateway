## 1. Settings boot 與 persistence

- [x] 1.1 交付 `Studio v2 settings boot from backend`：讓 settings page 啟動時讀取真實 backend settings，並以前端 boot test 與 backend handler test 驗證沒有默默退回 demo-only defaults。
- [x] 1.2 交付 `Save bar persists real settings`：讓 `儲存所有設定` 經由真實 backend APIs 寫入設定並回報失敗，並以前端 integration test 驗證不再是 noop。

## 2. Connector pool backend wiring

- [x] 2.1 交付 `Connector pool CRUD uses backend APIs`：讓 settings connector list/create/update/delete 對接既有 connector endpoints，並以前端 plus backend tests 驗證 round-trip。
- [x] 2.2 交付 `Connector test uses real backend diagnostics`：移除 `Math.random` / timer mock test，改接真實 connector test API，並以前端與 backend tests 驗證 success / failure path。

## 3. 驗證

- [x] 3.1 補齊 tests 與手動 demo，證明 settings surface 已不再依賴 operator-visible mock / noop。
