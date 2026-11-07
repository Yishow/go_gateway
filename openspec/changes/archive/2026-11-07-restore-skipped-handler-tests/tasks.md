## 1. 重現與修復

- [x] 1.1 將原始 524 skip 逐一分類並保留 test/package/reason，數量必須相加為 524。
- [x] 1.2 實現「輪詢群組名稱唯一性保持一致」與「名稱檢查與寫入持有同一把鎖」：用 handler 與 MemoryRepository 回歸重現同名被接受，修正後拒絕且原資料不變。
- [x] 1.3 兩個設備 handler「使用真實 loopback Modbus peer」 實際 activate/disable 並驗證持久狀態，不再無條件 skip。
- [x] 1.4「硬體驗收須明確選入」：真 PLC 測試在 opt-in 前不執行網路操作；以 gate 回歸及預設 skip 執行結果驗證。

- [x] 1.5「區分過期 CAS 與損毀 lineage」：資料庫 proof fixture 分別驗證 stale 409 與損毀 500，執行結果明列各 dialect。

## 2. 驗證與交付

- [x] 2.1 受影響 handlers/pollinggroup 完整測試、race 選測與 lint/vet 無新增問題。
- [x] 2.2 Spectra verify/review、規格 validate 與行數檢查通過，外部 DB/硬體未執行項目如實記錄。
