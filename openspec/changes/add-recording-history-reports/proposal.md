# 以同一份可信資料提供曲線、用量與事件報表

## Why

留下資料還不等於能用：使用者需要從設備直接看到曲線、每日用量與異常，並知道資料是否完整。
報表不能重算另一份答案，也不能在斷線時用零或平滑線掩飾缺口。

## What Changes

- 在 /studio/runtime 加入相同計畫的歷史、趨勢、用量、狀態／事件檢視與 CSV 匯出。
- 提供範圍、粒度、資料來源、完整度與計算版次一致的後端查詢 API。
- 圖表以實際摘要降採樣，保留採樣最高最低與缺口；不同單位不混成一條數字。
- 支援同方案多設備比較、明確計量範圍的合計、區間用量回查原值。
- 補齊六份 change 的跨平台、資料庫、故障恢復與使用者完整驗收。

## Capabilities

### New Capabilities
- `recording-history-reports`: 歷史查詢、品質呈現、事件、匯出與完整驗收。

### Modified Capabilities
- 無。沿用 post-setup focused monitor，不另建 fleet-first dashboard 或退休路由。

## Impact

依賴 1～5。新增 history query service、adapter range queries、frontend runtime tabs與匯出。
以同一套持久化結果為來源；不在瀏覽器計算用量，不引入需要使用者另裝的報表服務。
