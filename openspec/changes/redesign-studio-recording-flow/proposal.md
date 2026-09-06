# 將 Studio 改成帶人完成記錄的四步流程

## Why

目前流程讓使用者先理解接入規則、Tag、Row group、DDL 與欄位配對，才知道資料是否寫得進去。
資料庫頁也有範例欄位與實際表資訊混淆的風險；新增語意後更需要按用途引導，而不是多塞幾個下拉選單。

## What Changes

- 在現有 `/studio/v2` 依序完成連接設備、選資料、確認數值與用途、設定記錄並啟動。
- 提供三相電表與混合感測器範本、每項例外設定、批次預覽與套用。
- 最後一步先選記錄用途／目的地，再連線測試、選表或建表、預覽實際紀錄、確認啟動。
- 改為真實欄位與可確認配對；移除無意義的順序填滿和重複欄位 fallback。
- 統一儲存／生效／記錄／交付狀態、按鈕、錯誤修復、草稿接續與進階資訊。
- 保留所有既有路由身份、Share gate 與工程測試工具，不新增平行工作台。

## Capabilities

### New Capabilities
- `guided-recording-workflow`: 白話流程、範本、逐項用途、草稿和真實結果。

### Modified Capabilities
- `datalink-workbench-v2-step1-device`: 以明確選定的設備範圍判斷繼續，不把未選設備視為已通過。
- `datalink-workbench-v2-step4-database`: 連線設定與自動配對改用真實狀態和可確認建議。

## Impact

依賴 1～4 的 API 和啟動能力。主要影響 workbench-v2 shell、step1～4、settings、locales 與相關測試。
不以 mock 完成圖或前端成功訊息代替後端交付；本 change 僅起草，沒有改現行 UI 或啟動任何設備。
