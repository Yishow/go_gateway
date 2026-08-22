## Context

`start.sh` 已有 `tests/shell/` 行為契約測試，`start.ps1`（重構前 2007 行）完全沒有對等驗證。本次重構將 `start.ps1` 拆為主腳本 + 三個工具模組，趁結構落地時補上 Windows 軸契約測試，使兩軸行為對等可執行驗證。約束：PowerShell 函式庫載入機制、與 `start.sh` 的語意對齊、歷史超限檔僅可縮減行數。

## Goals / Non-Goals

**Goals**
- Windows pwsh 軸對五個行為契約（啟動順序、啟動鎖、受管端口、進程樹、日誌降噪）有可執行測試
- `start.ps1` 可被測試載入而不觸發主流程
- 工具模組邊界清楚，與 `scripts/start-*-utils.sh` 語意對齊
- 修正 Windows Air hot-reload entrypoint

**Non-Goals**
- 不改使用者互動流程、選單行為
- 不引入 Pester 測試框架（維持 plain pwsh 腳本 + exit code 慣例，與 tests/shell 對稱）
- 不覆蓋 macOS／ARM PowerShell 行為

## Decisions

### 以環境變數 library 模式載入 start.ps1

測試以 `GATEWAY_START_PS1_LIBRARY_ONLY=1` dot-source `start.ps1`：主腳本定義所有函數後，在進入主流程前檢查此變數並提前 return。替代方案：(a) 拆成 `.psm1` 模組——需要改變載入語意且與 start.sh 單檔慣例不一致；(b) Pester——引入新框架依賴，違反 repo 現行測試慣例。Library 模式改動最小，且與 shell 測試「source 後 stub」的模式對稱。

### 工具模組經動態作用域讀取主腳本狀態

`scripts/start-{process,port,log}-utils.ps1` 的函數直接讀取主腳本的 `$Port`、`$script:FrontendDevPort`、`$script:TMP_DIR` 與輸出函數（`Write-Info` 等），不經參數傳遞。理由：PowerShell dot-source 後函數在呼叫端作用域執行，此模式與 `start.sh` 內函數讀取全域變數的語意完全對稱，測試也能以覆寫函數／變數的方式注入 stub。代價是模組不可獨立使用——可接受，因為它們本就是 start.ps1 的內部結構。

### 契約測試以函數覆寫 stub 驗證 orchestration 順序

`tests/powershell/start-dev-mode-order.ps1` 在 library 載入後覆寫 `Start-BackendProcess` 等函數為記錄器，斷言 `Start-DevMode` 的呼叫順序為 `sync>clear-backend>start-backend>wait-backend>start-frontend>wait-exit>cleanup`。此模式移植自 `tests/shell/start-backend-before-frontend.sh` 的 stub 手法，兩軸斷言字串一致。

### 啟動鎖以端口命名目錄實作

`Acquire-RunLock` 在 `bin/tmp/start-<port>.lock` 建目錄並寫入 pid：第二次取得失敗且訊息說明持有者；持有者進程已死時清理過期鎖後重新取得；`Release-RunLock` 僅移除自己持有的鎖。與 `start-lock-management.sh` 四情境一一對齊。

### 受管端口集合防範 PowerShell range 陷阱

`Get-ManagedPorts` 以 `@(...) + (8080..8090)` 串接展開 range——PowerShell 中逗號優先於 `..`，`@($a, 8080..8090)` 會被解析為錯誤範圍。此陷阱已在 `tests/powershell/start-managed-ports.ps1` 立為回歸情境。

### 進程樹終止使用 taskkill /T

`Stop-ProcessTree` 以 `taskkill /PID <id> /T /F` 終止整棵樹，避免 cmd → node、air → go → exe 鏈留下孤兒進程佔用端口。`Wait-PortReady` 提供有界重試，並在後端進程已退出時快速失敗（不等滿重試次數）。

### 日誌降噪實體清單單一來源

dashboard 輪詢的 datalink 實體清單（devices、polling-groups、points、mappings、tags）抽為常數：shell 端 `DATALINK_DASHBOARD_ENTITIES`、pwsh 端對應變數，兩處比對邏輯共用同一清單，避免新增實體時只改一軸。

### Air 進入點補 .exe 副檔名

Windows 上 Air 只執行 `.exe` 副檔名進入點。`.air.toml` 的 `entrypoint` 與 build 輸出、`start.sh` 的 `AIR_BIN` 統一改為 `gateway-air.exe`，Linux 軸不受影響（Air 對帶 .exe 的輸出仍正常執行）。

## Implementation Contract

- **Library 載入**：`GATEWAY_START_PS1_LIBRARY_ONLY=1` 時 dot-source `start.ps1` 僅定義函數，不執行任何主流程、不出現互動提示。
- **啟動鎖語意**：鎖目錄 `bin/tmp/start-<port>.lock` 含 pid 檔；重複取得失敗訊息含「已有 start.ps1 管理此 repo/port」；過期鎖（pid 對應進程不存在）清理後可重新取得並顯示「發現過期啟動鎖」；非持有者的鎖不可被 Release 移除。
- **受管端口**：預設為 8080..8090 共 11 個 + 前端 dev port + Modbus share port（5020），去重；null／非數字輸入略過不炸。
- **進程樹**：`Stop-ProcessTree` 終止父進程及其全部子進程；`Wait-PortReady` 端口監聽時成功、未監聽時逾時失敗（訊息含「在預期時間內未成功啟動」）、後端進程已退出時快速失敗（訊息含「後端進程已退出」）。
- **日誌降噪**：GET 200 的 dashboard 輪詢行計數不輸出；5xx 以紅色 [HTTP] 格式輸出；非 Go-log 格式錯誤行紅色輸出；Go-log 格式行非 Verbose 靜默；`Show-LogNoiseSummary` 列出各分類計數。
- **Air entrypoint**：`.air.toml` build 輸出與 entrypoint 均為 `./bin/tmp/gateway-air.exe`；`start.sh` 的 `AIR_BIN` 同步。
- **In scope**：上述行為契約與其測試。**Out of scope**：`cross-platform-loop` 驗證序列（屬既有規格）、前端測試重構、互動選單行為。

## Risks / Trade-offs

- [動態作用域耦合：主腳本變數改名會靜默破壞 utils 模組] → 契約測試覆蓋各函數的主要路徑，重命名後測試立即失敗
- [library 模式環境變數名誤用於生產環境] → 變數僅在 start.ps1 內部判讀，文件註明測試用途
- [taskkill /F 強制終止可能造成資料寫入中斷] → 僅用於本腳本自己啟動的 dev 進程樹，與 start.sh cleanup 語意一致
- [PowerShell 7 與 Windows PowerShell 5.x 行為差異] → 測試統一以 pwsh（PowerShell 7+）執行，與 repo 開發環境一致

## Migration Plan

純開發工具鏈變更，無部署影響。退回方式：revert 單一 commit 即可，`start.ps1` 回到 monolith 版本。

## Open Questions

(none)
