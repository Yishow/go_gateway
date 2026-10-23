## Problem

原工作樹保存了連線必填欄位彙整實作與測試，尚未提交。舊 main 只回第一個缺欄位；目前草稿雖分別彙整 missing／empty，混合時仍以 missing 遮蔽 empty，需重送才能得知另一必填問題。

## Root Cause

驗證迴圈與 return 都依第一錯誤或錯誤類別提早結束，未把所有必填問題共同呈現。README 與計畫對 durable 路徑亦有過期或過度絕對陳述。

## Proposed Solution

保存原稿的 schema 順序、單欄訊息與 FATEK 條件語意，一次回報缺少與空值 required 欄位，再處理既有 range／enum。建立精準混合回歸與正式 delta；文件以 production wiring／crash tests 核對 durable 及 unknown 邊界，現場紀錄去識別並標示歷史來源。

## Success Criteria

Modbus／FATEK 混合 missing+empty 在一次錯誤內可見，單類及有效設定行為不變；focused RED/GREEN、完整 Go test/vet/lint、規格／行數 gate 通過後封存。沒有現場或正式資料庫操作。

## Capabilities

### New Capabilities

無。

### Modified Capabilities

- `datalink-api`: 新增 Aggregated device connection validation，既有 HTTP error shape 不變。

## Impact

- Affected code:
  - Modified: internal/datalink/device/service_validation.go
  - New: internal/datalink/device/service_validation_test.go
  - Modified: README.md
  - Modified: docs/plans/2026-10-09-incomplete-features-plan.md
  - Removed: 無。
