## 1. Red 與失敗注入
- [x] 1.1 建立 journal/outbox/receipt 狀態機與固定時鐘測試，列出每個 transaction 提交邊界。
- [x] 1.2 Red：採集前後、checkpoint提交前後、外部commit後回應前逐點中止並重啟。
- [x] 1.3 Red：重複輸入／重送／舊結果revision、亂序、雙worker、poison record、滿磁碟與慢目標。

## 2. 綠燈耐久路徑
- [x] 2.1 實作 local journal 與成功回應前耐久保證，串接 2 的 transaction 和 3 的 calculator checkpoint。
- [x] 2.2 實作 per-destination outbox、lease/fencing、順序與有界併發。
- [x] 2.3 實作 managed receipt 與資料同交易，驗證未知commit結果的復原流程。
- [x] 2.4 實作 retry/blocked/quarantine、修復後續送與 typed error，不使用無限重試迴圈。
- [x] 2.5 實作 quota/告警/暫停/人工丟棄稽核、checkpoint/evidence/receipt 保留相依。
- [x] 2.6 實作正常關機drain、crash recovery、plan revision安全切換與不改送目的地。

## 3. 重構與驗收
- [x] 3.1 整合 workspace readiness/activation、delivery status API 與現有 Share 生命周期回歸。
- [x] 3.2 SQLite/PostgreSQL/MySQL 實際執行 receipt 去重、連線失敗和 schema變更整合測試。
- [x] 3.3 Windows/Linux/ARM 做關機、重啟與受限容量測試；未具備環境者列 pending，不假稱通過。
- [x] 3.4 執行 Go test/vet/lint、壓力與長時間斷線測試，記錄樣本數守恆與本機資源。
- [x] 3.5 strict spec 驗證及 rollout/rollback runbook完成後，才允許新 managed recording 正式啟用。

