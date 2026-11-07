# 資料庫補測修復驗證

## 範圍

基準 ff4598e3。六個程式/測試檔案均由本 change 負責；新增工件與 recording-plan-storage 的兩個 delta requirements 一併核對。Spectra scope 屬 approximated_base，需以明確檔案清單與內容指紋確認歸屬。

## 重現與根因

- PostgreSQL 初始 migration 回 42883：next-month 表達式是 timestamp，但函式參數是 DATE。獨立 rollback schema 測試先失敗，修正 cast 後通過。
- schema apply 建立九個DDL語句後回 unknown：驗證器把產生的 val_dec NUMERIC 當成文字欄位。單元測試重現 numeric 被拒絕，修正專用分類後正常。
- MySQL系列的手工snapshot缺 immutable counter bindings；補既有 proof builder 的輸出，不改 production fail-closed。

## 規格對應

1. Managed exact decimal schema compatibility：schema_tables.go 的 decimal 分類，TestManagedDecimalColumnAcceptsExactStorageTypes/RejectsLossyStorageTypes 覆蓋全部八個具體型別。真PG route同時驗證首次apply與第二次相容no-op的 verified digest。
2. PostgreSQL 初始月份分區可建立：migration的DATE cast，TestPostgresPartitionMigration實際建立分區、建立未來分區及123.45寫入讀回。fixture的transaction-local schema在結束後查詢count為0。

## 本輪實際結果

- 最終decimal選測：10個測試事件通過，含2個父測試與8個型別子案例。text/numeric接受；real/double precision/integer/boolean/numeric(38,12)/decimal(20,6)拒絕且不儲存token。
- 真PostgreSQL：TestPostgresPartitionMigration，以及3個TestStudioV2RecordingRoutes_Postgres案例全部通過，包含建表驗證/no-op、不相容表保持不變、權限拒絕。
- 真MariaDB 11.4.13：TestMariaDBCircuitIntegration選測6個事件通過，涵蓋unscaled uint64與exact decimal。
- 最終全庫golangci-lint：0 issues；全庫go vet：exit 0。完整Go測試由整合流程作為發布前gate，不能由本報告的選測推論全套通過。
- validate/analyze：通過、0 findings；diff/行數檢查通過。

## Review

依正確性、效率、簡化重用、專案慣例四方向審查六個diff，並經獨立reviewer再次核對，未發現具體新增缺陷。

初稿曾允許有界numeric/decimal，review指出可能rounding/overflow；已縮限為生成DDL的無界NUMERIC，並將有界型別加入負向回歸。PostgreSQL inspector保留numeric(p,s)，不會誤認成無界。

隔離fixture同一transaction建立schema/設定search_path/執行migration/rollback，不再刪public資料。MySQL CanonicalApply分支仍由ResolveAppliedAt取得真正snapshot，產品proof guard未旁路。

## 限制

MySQL/MariaDB全套補測與原524項逐一覆蓋矩陣由獨立交付報告彙整；這裡只記錄本change已完成的驗證，不聲稱未跑的dialect、TimescaleDB、SQL Server或真PLC通過。
