## 1. 修復已重現的缺陷

- [x] 1.1「小數分類只補無界 NUMERIC」實現 Managed exact decimal schema compatibility：以正負向 preview 回歸驗證，既有身分契約不變。
- [x] 1.2「分區函式參數明確使用 DATE」實現 PostgreSQL 初始月份分區可建立：真PG migration、分區、123.45讀回驗證。
- [x] 1.3「分區回歸只使用可回滾 schema」：移除public清表，交易結束後query確認無測試schema殘留。
- [x] 1.4「手工 snapshot 對齊 Apply proof」：MySQL系列 fixture加入counter bindings，以真MariaDB unscaled/exact decimal整合驗證，產品fail-closed不變。

## 2. 驗證及交付

- [x] 2.1 真PG建表/apply與no-op回歸、recordingplan小數型別回歸及雙引擎fixture結果逐項記錄，不以skip當pass。
- [x] 2.2 verify/review、validate/analyze、相關lint/vet/行數檢查通過後archive，交由整合流程commit/push。
