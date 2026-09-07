# 讓曲線、用量與狀態報表使用正確算法

## Why

溫度、壓力與功率需要看變化；累積電量／水量要算差值；PLC 狀態需要看變化與持續時間。
把這些全部平均、全部相加，或每分鐘只取最後一筆，都會產生看似合理但不可信的報表。

## What Changes

- 增加後端持續運作、與瀏覽器無關的樣本統計與計量引擎。
- 涵蓋 gauge/rate 統計、counter/signed_counter 差分、已知區間 delta、state/event/text。
- 明確處理起始基準、重置、回捲、缺口、亂序、時區、重算、估算與精度。
- 支援經確認的速率積分與有限衍生公式，不把估算當電表實測或計費依據。
- 共用一套結果供資料庫與報表使用；不同目的地或畫面不得各算各的。

## Capabilities

### New Capabilities
- `measurement-aggregation`: 時間統計、用量、狀態、衍生值與可追溯重算。

### Modified Capabilities
- 無。重用 mapping 既有有效工程值，不改既有 edge-analytics FFT/PID 或 Local Modbus 計算。

## Impact

依賴 1、2；使用 2 的 transaction/checkpoint 介面，4 提供耐久交付。
擬新增 recording calculation services、checkpoint/result repositories 與測試 fixtures。
完整資料例與預期值見 `docs/plans/telemetry-recording/acceptance.md`。
