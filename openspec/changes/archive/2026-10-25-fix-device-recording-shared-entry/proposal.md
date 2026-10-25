## Why

Step 4 依設備是否保留舊受管寬表，自動分派「基本記錄」或「迴路記錄設定」。同一工作區的兩台設備因此無法從相同入口設定瞬時／累積記錄；舊群組即使未套用，也會讓第一台留在舊面板。使用者已明確要求修正這項分流。

## What Changes

- 所有資料庫設備預設開啟相同的用途、瞬時／累積與準備／開始介面。
- 舊寬表顯示已儲存設定摘要，以明確的「編輯既有記錄」進入原編輯器；舊週期、缺值政策與草稿不自動轉換。
- 未套用的舊草稿不攔住新設定；已有套用版本的舊設定提供條件式切換說明，仍由既有伺服器 ownership／reviewed cutover 判斷來源能否交接。
- 保留明確的只準備模式、Share-only、操作復原與 Advanced 同一份設定。
- 本案追蹤已完成修復的後續驗證與交付，不重新歸檔已提交的 unify-device-recording-workflow。

## Capabilities

### New Capabilities

無。

### Modified Capabilities

- `guided-recording-workflow`：Purpose-led layout with one advanced authority 明訂舊寬表設備也從共同設定入口開始，舊編輯為明確的次要操作。

## Impact

只包含 Step 4 的面板選擇、舊設定摘要、en／zh-TW 文案、對應前端測試與隔離瀏覽器 fixture，以及既有操作指南。沒有 API、Go、資料模型、資料 migration、SQL writer 或新依賴變更。MySQL／MariaDB、SQL Server 與其他 WIP 不屬於本案。
