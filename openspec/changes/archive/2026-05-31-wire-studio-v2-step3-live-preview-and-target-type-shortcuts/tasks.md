## 1. Live preview wiring

- [x] 1.1 補前端 RED 測試，覆蓋 Step 3 preview 的 loading、success、error 與 stale response guard。
- [x] 1.2 將 Step 3 `轉換管線預覽` 從本地 mock 改接既有 `POST /api/v1/datalink/mappings/preview`，保留 deterministic raw seed 但改由 backend 執行 transform pipeline。
- [x] 1.3 維持 `API payload preview` 為 display-only，驗證它不會因 live preview 而額外發 request。

## 2. Target type quick actions

- [x] 2.1 補前端 RED 測試，覆蓋當前列 target type 快捷設定與 `套用到全部列`。
- [x] 2.2 在 Step 3 preview 卡新增 `bool`、`int16`、`float64`、`string` 快捷操作，並沿用既有 reducer/autosave flow。
- [x] 2.3 將 `套用到全部列` 收斂為 target_type-only 的批次套用，不新增第二套批次更新機制。

## 3. 驗證

- [x] 3.1 執行 `cd frontend && npm run test -- --run tests/unit/workbench-v2/step3-mapping.test.tsx`。
- [x] 3.2 執行 `cd frontend && npm run build`。
- [x] 3.3 執行 `git diff --check`。
