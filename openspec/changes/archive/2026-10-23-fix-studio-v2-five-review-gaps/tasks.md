## 1. 五項修復

- [x] 1.1 交付 Scalar point reads preserve protocol count units 與單點協定數量與直接批次分離；先失敗再修，驗periodic/manual/read-now、32/64-bit、direct batch及Modbus/FATEK契約。
- [x] 1.2 [after: 1.1] 交付 Confirmed template intent and recoverable failures 與既有UI依真實證據與明確意圖顯示；真preview false確認、安全錯誤、scope變更，先失敗再修及測409/500/斷線。
- [x] 1.3 [after: 1.2] 交付 Focused header claims only evidenced stages；以header測試覆蓋等待首筆、Share-only、停止、失敗、unavailable且無入庫宣稱。
- [x] 1.4 [after: 1.3] 交付 Complete occupied range without changing points；以range summary測試驗MC/Modbus32/64-bit、單word/bit、進位與stride/span。
- [x] 1.5 [after: 1.4] 交付 Basic missing-value policy reveals required gating；以Basic卡與snapshot測試驗required skip、optional partial NULL原因、入口與恢復，禁止暗改optional。

## 2. 整合與交付

- [x] 2.1 [after: 1.5] 以隔離MC D0–D7模擬器及temp SQLite真UI完成三桶SQL核對縮放、單位與精度；真API驗範本成功與未確認拒絕、失敗重試；保留失敗和成功證據。
- [x] 2.2 [after: 2.1] 執行完整Go test/vet/lint、frontend lint/typecheck/tests/build及Playwright，獨立review必要修復，核對五項scenario覆蓋與既有安全保護。
- [x] 2.3 [after: 2.2] 整理實測證據、限制與歸檔前核對，確認正式同步規格及archive條件；Git發佈與原工作樹同步結果另於交付回報核實，本任務勾選不代表已push。
