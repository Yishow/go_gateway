# 剩餘變更分批覆核與驗證

基準為迴路記錄主 change 的提交 `521364489629be24b235ac88dc70870bb1410f23`。本紀錄只處理當時其餘未提交內容，依使用者授權 review、必要修復及繁中分批提交；未實作暫存的 MySQL／MariaDB 或 SQL Server changes。

## Review 結果

規範與需求分開檢視，主代理核對來源、呼叫端、差異與失敗證據。

- 規範 P1：兩個 launcher 的埠等待可把無關程序當就緒。修正為先檢查自有程序存活，再核對 listener 為它自身或後代；Air 子程序可正常就緒，外來 listener 不接受。
- 規範 P2：鎖衝突提示引用不存在的 `--stop-all`。改為關閉原啟動器、確認原 PID 已結束再重試，不新增清除其他程序的操作。
- 規範 P2：日誌測試強制 ANSI 色彩後未還原，全域設定使 shuffle 測試確實失敗。補 cleanup，日誌延遲分級改用等價 switch 通過 lint。
- 需求 P2：映射修復文件的兩張歷史截圖不存在。移除失效證據連結，明列無法覆核，另保存此次實際版面截圖。

其餘確認範圍為既有 launcher 選單／熱重載／打包、保留新版工具鏈的相容遷移、UPDATE 404 明確重試後重新取得映射身分、短時間查詢快取／低噪音日誌，以及 legacy writer 缺表恢復。未增加 DELETE 404 政策或讓 shell 讀設定庫做 Share 埠預檢。

## 此次實際驗證

- 前端完整 Vitest：260 個檔案、1,812 項，零失敗／略過；ESLint 與 TypeScript／Vite build 通過。套件 manifest 與 npm lockfile 的根依賴一致。
- Node 建置／樣式契約：22／22 通過。
- Playwright：8／8 通過，驗證 Basic／Advanced 捲動、即時值溢出及工作台版面；[截圖與尺寸](studio-mapping-recovery-20261024.md#2026-10-25-提交前覆核)。使用真元件與 fixture／替代 API，不代表真設備或 SQL 驗收。
- Launcher 實跑：隔離副本 fresh bootstrap、Vite API proxy、host Go target、修改 lib 自動換後端 PID、編譯失敗停止舊程序／修復恢復、WebSocket HMR、TERM 清理雙埠與鎖全部通過。
- macOS PowerShell：11 支 launcher 契約／mock 測試通過；Windows 原生 `start-process-tree.ps1` 不適用於 macOS，列 NOT_RUN。
- 日誌測試：隨機順序重跑五次通過。
- Legacy writer：`go test -race` 選中 15 個 PASS 路徑、零 SKIP，涵蓋 SQLite 及 PostgreSQL 真實缺表重建、重試、並行恢復、ownership 排除及權限拒絕。不是只以 SQL 錯誤字串模擬成功。
- PostgreSQL fixture：既有 loopback 容器 PostgreSQL 17.11／TimescaleDB 2.31.0，本次 legacy writer 只驗 plain PostgreSQL；完成後 `gw_*` schema 為零，容器保持原狀。
- `go vet ./internal/api ./internal/datalink/dbtarget`、兩者的 golangci-lint、`go build ./...`、diff 與行數檢查通過。

Shell 完整契約測試 15／15 通過，零失敗／timeout。修正既有 lifecycle fixture，使 listener 回報實際子程序 PID，不用固定虛構 PID；產品 ownership 檢查維持嚴格。

## 限制與保留內容

- Windows 原生、Linux／Intel macOS 原生、實體 PLC、LAN、正式資料庫與部署皆 NOT_RUN。PowerShell mock、跨平台建置與 macOS 實跑分開計算。
- 既有 API `TestStudioV2RecordingRoutes_PostgresApplyCreatesAndVerifiesTheManagedSchema` 與 storage `TestPostgresPartitionMigration` 基線失敗未在本批修復，未宣稱全 repo Go test 通過；本批依行為選測。
- 前端 build／瀏覽器執行仍有既有 PostCSS `from` 與打包大小等警告，不宣稱零警告。
- `.pnpm-store/v11/projects/978cac21dd875c23d2063f9ee8ef9c85` 是本機生成快取 symlink，排除提交並保留原暫存內容。
- 不 push、不 archive，兩個資料庫 change 保持 parked；不更動主 change 的既有提交。
