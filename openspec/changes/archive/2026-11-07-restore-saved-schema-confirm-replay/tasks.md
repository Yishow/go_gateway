## 1. 修復與回歸

- [x] 1.1「已保存proof沿用既有確認權威」實現 Exact persisted successor schema confirmation replay：plain/native_monthly離線精確重放回成功且零probe/DDL，local/external狀態不變。
- [x] 1.2「非精確重放保持拒絕」：兩模式各5個digest/operation/connector/lineage/alias負向回歸拒絕且零probe/變更。

## 2. 驗證

- [x] 2.1 核對handler/既有service所有scope保護與原未保存proof路徑，完成真SQLServer正負向驗證、verify/review與validate/analyze後archive，最終API/lint/vet由整合發布gate負責。
