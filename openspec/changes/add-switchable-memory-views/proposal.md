## Why

來源設定的記憶體視圖目前只有單一密集格子排版，資料型別與跨暫存器值不易一眼辨認。讓操作員在相同資料、區域與位址下切換五種閱讀方式，依保存的五套明暗參考圖校正，同時維持既有安全與即時讀值契約。

## What Changes

- 頂部新增格子專注、範圍導航、水平地圖、範圍章節、格子＋詳情五種視圖，切換保留區域、位址、頁面及選取狀態。
- 型別以文字及一致背景／輪廓區分；單暫存器值一格，多暫存器值連接呈現，跨列與頁仍可辨識其起始位址及分段。卡片保留間距，禁止格子左色條。
- 五種視圖共用現有 occupancy、MC proof、讀值與安全控制；空格新增、已用查看、重疊 owner 選擇及未知不可新增皆維持。
- 明暗沿現有主題，驗證桌面／平板／手機與鍵盤；保留五圖私人證據，逐套實際瀏覽器對照。

## Capabilities

### New Capabilities

- `recording-memory-layouts`: 記憶體五視圖切換、共用狀態、型別連續性及響應式呈現。

### Modified Capabilities

無。既有 simple-recording-home 的資料來源與安全規則不變。

## Impact

- frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeMemoryView.tsx 與相鄰呈現元件、樣式。
- frontend/src/i18n/locales/en/recording-home.json、frontend/src/i18n/locales/zh-TW/recording-home.json。
- frontend/tests/unit/workbench-v2 與相關瀏覽器驗證。
- 不變更後端 API、協議解析、訂閱生命週期、依賴與部署。
