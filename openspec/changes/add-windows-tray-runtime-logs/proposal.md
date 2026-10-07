## Why

Windows 單執行檔部署目前仍以 console 視窗承載日誌與結束入口；直接隱藏視窗會失去狀態、操作及啟動失敗的診斷方式。這次規劃以系統列圖示管理本機程序，並在既有嵌入式網頁提供安全、有容量限制的執行日誌。

## What Changes

- Windows 桌面發行版以 GUI subsystem 啟動，不出現黑色 console，也不因開瀏覽器短暫跳出 console；Go backend、前端資產與圖示一起內嵌
- 右鍵選單提供「開啟設定頁」、「查看執行日誌」、唯讀程序／採集狀態、「版本資訊」及需確認的「結束」；關閉瀏覽器不會停止閘道
- 定義桌面／headless 模式、資料路徑辨識、同資料庫單實例、port 衝突、tray 初始化失敗、Explorer 重啟及有期限的關閉協調
- 新增 `/studio/logs` 與受本機限制的 log snapshot／SSE；提供等級、來源、搜尋、暫停、繼續及跟隨最新紀錄
- 統一安全事件投影，遮蔽先於新 memory／file／SSE sink；限制記憶體、檔案輪替、訂閱數及慢讀者，明示掉落、截斷、重連與歷史缺口
- 保留現有 console 開發流程、Linux/macOS CLI、既有環境變數優先序、Studio routes、durable group delivery 與 Modbus Share fail-closed 契約
- 本次只提交規格；不實作、不打包或部署、不 archive、不勾選產品任務

## Capabilities

### New Capabilities

- `windows-desktop-lifecycle`: Windows 無 console 發行模式、tray 選單、資料來源與單實例保護、啟動診斷、狀態及確認關閉
- `runtime-log-viewer`: 受本機限制且不洩漏秘密的程序日誌、有限保留與串流、網頁互動及可觀察的遺失／斷線狀態

### Modified Capabilities

無。此案新增兩項能力；既有 `embedded-frontend-delivery`、`start-script-contracts`、`runtime-dashboard-truthfulness` 的產品邊界維持不變。既有 `/debug/logs` 仍是工程工具紀錄，不被改成程序日誌；console HTTP formatter 的既有 query 相容測試不代表新安全 sink 可保留 query 值。

## Impact

- 主程式／建置：`cmd/test_ui/main.go`、`cmd/test_ui/server_runtime.go`、`cmd/test_ui/service_wiring.go`、`cmd/test_ui/harness_config.go`、`Makefile`、`scripts/build.ps1` 及 Windows 啟動腳本的模式選擇
- 新增平台隔離的 desktop/lifecycle 與安全 logging 模組；以現有 `golang.org/x/sys` 的 Windows API 能力優先評估，不引入 WebView、Electron 或另一套產品 UI
- API／UI：`internal/api/router.go`、新增 log handlers、`internal/config/config.go`、`internal/web/embed.go`、`frontend/src/App.tsx`、新 log page/service/hooks/types 及 en／zh-TW 文案
- 新增 GET `${API_BASE_PATH}/system/logs` 與 GET `${API_BASE_PATH}/system/logs/stream`；不新增遠端控制、登入平台、log DB、export／刪除 server logs API
- Windows 實機、embedded-browser、跨平台 CLI 及 durable shutdown 回歸是後續實作驗收；本次文件檢查不構成產品驗證
- 固定來源、重疊盤點、工具來源與未驗界線見 [evidence.md](evidence.md)；完整決策見 [design.md](design.md)
