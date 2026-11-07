## 1. 修復 fixture

- [x] 1.1「Timescale 清理使用新連線」：真資料庫 recovery/cleanup 回歸通過且結束 extension count 為 1。
- [x] 1.2「同批 Apply 固定時鐘」：同批群組生效時間一致，恢復案例重跑3次通過且defer還原clock。
- [x] 1.3「NativeMonthly 證據完整傳遞」：SQL Server native cleanup不再因converter漏欄位而拒絕，保留原斷言。
- [x] 1.4「時間以瞬間而非字串比較」：SQL Server runtime boundaries接受同瞬間表示，解析錯誤與不同瞬間仍失敗。

## 2. 驗證

- [x] 2.1 記錄真DB回歸與靜態review，完成validate/analyze/diff/行數檢查並archive；最終API/lint/vet交由整合發布gate。
