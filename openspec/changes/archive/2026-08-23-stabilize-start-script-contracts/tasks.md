## 1. start.ps1 拆分與工具模組

- [x] 1.1 實作「以環境變數 library 模式載入 start.ps1」：GATEWAY_START_PS1_LIBRARY_ONLY=1 時 dot-source 僅定義函數、在主流程前提前 return，不執行選單或啟動任何進程（對應 Requirement: Start scripts load in library mode for contract testing）。驗證：tests/powershell/start-dev-mode-order.ps1 可載入 start.ps1 並呼叫 Start-DevMode stub 後通過（exit 0）。
- [x] 1.2 [P] 建立 scripts/start-process-utils.ps1：「工具模組經動態作用域讀取主腳本狀態」模式，提供 Acquire-RunLock／Release-RunLock（bin/tmp/start-<port>.lock 目錄鎖 + pid 檔）與 Stop-ProcessTree（taskkill /T）／Wait-PortReady。驗證：pwsh -NoProfile -File tests/powershell/start-lock-management.ps1 與 tests/powershell/start-process-tree.ps1 全情境通過。
- [x] 1.3 [P] 建立 scripts/start-port-utils.ps1：Get-ManagedPorts 以 @(...) + (8080..8090) 串接展開 range，落實「受管端口集合防範 PowerShell range 陷阱」，null／非數字輸入略過。驗證：pwsh -NoProfile -File tests/powershell/start-managed-ports.ps1 三情境（預設 13 端口、自訂 14 端口、null 輸入）通過。
- [x] 1.4 [P] 建立 scripts/start-log-utils.ps1：Write-RuntimeLogLine 分類降噪（dashboard 輪詢計數、5xx 紅色 [HTTP]、非 Go-log 錯誤紅色、Go-log 非 Verbose 靜默）與 Show-LogNoiseSummary。驗證：pwsh -NoProfile -File tests/powershell/start-log-noise.ps1 六情境通過。

## 2. PowerShell 契約測試（Dev mode startup order / 啟動鎖 / 受管端口 / 進程樹 / 日誌降噪）

- [x] 2.1 [P] 撰寫 tests/powershell/start-dev-mode-order.ps1：「契約測試以函數覆寫 stub 驗證 orchestration 順序」，覆寫 Sync-EmbeddedFrontendIfRequested 等函數為記錄器，斷言順序 sync>clear-backend>start-backend>wait-backend>start-frontend>wait-exit>cleanup（對應 Requirement: Dev mode startup order is contractually fixed）。驗證：腳本自身 exit 0。
- [x] 2.2 [P] 撰寫 tests/powershell/start-lock-management.ps1：四情境覆蓋「啟動鎖以端口命名目錄實作」——存活持有者二次取得失敗且訊息含「已有 start.ps1 管理此 repo/port」、過期鎖清理後重取並顯示「發現過期啟動鎖」、Release 移除自身鎖、非持有者鎖不被移除（對應 Requirement: Run lock enforces a single start manager per repo and port）。驗證：腳本自身 exit 0。
- [x] 2.3 [P] 撰寫 tests/powershell/start-managed-ports.ps1：驗證預設 8080..8090+5173+5020 去重 13 端口、自訂 Port 3333 時 14 端口、null 輸入不炸（對應 Requirement: Managed port set expands ranges without duplicates）。驗證：腳本自身 exit 0。
- [x] 2.4 [P] 撰寫 tests/powershell/start-process-tree.ps1：cmd→ping 進程樹全滅、Wait-PortReady 監聽成功／未開逾時失敗／後端已退出快速失敗（對應 Requirement: Process tree termination and bounded port readiness；「進程樹終止使用 taskkill /T」）。驗證：腳本自身 exit 0。
- [x] 2.5 [P] 撰寫 tests/powershell/start-log-noise.ps1：六情境覆蓋降噪計數、5xx 紅色、錯誤行、Go-log 靜默、Verbose 不降噪、摘要輸出（對應 Requirement: Runtime log noise is reduced with counters and summary）。驗證：腳本自身 exit 0。

## 3. Shell 端對齊與 Air entrypoint

- [x] 3.1 [P] scripts/start-log-utils.sh 落實「日誌降噪實體清單單一來源」：DATALINK_DASHBOARD_ENTITIES 常數供 get_log_noise_category 與 runtime_log_line 兩處共用。驗證：bash tests/shell/start-backend-cleanup-no-log-wait.sh 與 bash tests/shell/start-backend-logfile.sh（涵蓋 runtime_log_line 行為的既有 shell 契約測試）通過。
- [x] 3.2 [P] 「Air 進入點補 .exe 副檔名」：.air.toml 的 build 輸出與 entrypoint 改為 ./bin/tmp/gateway-air.exe，start.sh 的 AIR_BIN 同步（對應 Requirement: Windows Air hot-reload uses an exe entrypoint）。驗證：.air.toml 內容審查（entrypoint 與 cmd 輸出一致）＋ bash tests/shell/start-env-load.sh、start-backend-cleanup-order.sh 等引用 AIR_BIN 的測試通過。
- [x] 3.3 全量驗證：依序執行五個 pwsh 契約測試（tests/powershell/start-{dev-mode-order,lock-management,managed-ports,process-tree,log-noise}.ps1）與 bash tests/shell/start-*.sh 全部契約測試，全部 exit 0；並確認 start.ps1 行數未超過重構前 2007 行（檔案行數規範：歷史超限檔僅可縮減）。
