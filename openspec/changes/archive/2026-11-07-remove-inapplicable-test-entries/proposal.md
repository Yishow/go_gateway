## Why

使用者要求移除先前補測後剩餘的八個不執行測試項，讓同步後的手動驗收範圍明確。移除不是測試通過，必要的崩潰恢復與能力保護仍須保持驗證。

## What Changes

- 刪除真實 PLC 驗收檔與長時間排程負載矩陣入口，保留短負載 pilot。
- SQLite 不支援的僅大小寫區別欄位組合不再註冊。
- 三個子程序入口改為普通 helper，由既有父測試依專用環境參數分派，不再當獨立測試列出。
- SQL Server 關閉能力的兩個分支改用測試替身驗證跳過決策，不故意產生 skip 事件。

## Non-Goals

不改產品功能、不接觸真設備、不聲稱移除項目已驗證、不移除崩潰恢復或能力 gate 的保護。

## Impact

- Affected specs: none
- Affected code: cmd/test_ui、internal/api、internal/datalink 的八個測試檔。
- Compatibility: no capability-level observable behavior changes
