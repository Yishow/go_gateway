## Why

`start.ps1`（2007 行）與 `start.sh` 聲稱跨平台對等，但行為契約只在 shell 端有測試（`tests/shell/`），Windows 端任何修改都無法驗證回歸：啟動鎖、受管端口展開、進程樹終止、啟動順序與日誌降噪全是黑箱。同時 Windows 上 Air 只執行 `.exe` 副檔名進入點，現行 `.air.toml` 設定導致 hot-reload 模式靜默失效。本次 `start.ps1` 大規模重構（抽出工具模組）已發生，需要補上規格與契約測試錨點，防止後續漂移。

## What Changes

- `start.ps1` 支援 library 載入模式（`GATEWAY_START_PS1_LIBRARY_ONLY=1` 時僅定義函數不執行主流程），使契約測試可 stub 依賴後注入驗證
- 抽出三個與 `start.sh` 對等的功能模組：`scripts/start-process-utils.ps1`（啟動鎖、進程樹終止、端口就緒等待）、`scripts/start-port-utils.ps1`（受管端口集合與清理）、`scripts/start-log-utils.ps1`（執行期日誌降噪）
- 新增五個 PowerShell 契約測試（`tests/powershell/`）：啟動順序、啟動鎖管理、受管端口、進程樹終止、日誌降噪，與 `tests/shell/` 對應測試一一對齊
- 對齊既有 shell 端契約：`DATALINK_DASHBOARD_ENTITIES` 抽為單一來源常數、managed ports range 展開回歸補測
- 修正 Windows Air hot-reload：`.air.toml` 與 `start.sh` 的 build 輸出與 entrypoint 改帶 `.exe`

## Non-Goals

- 不改動 `start.ps1`／`start.sh` 的使用者互動流程與選單行為
- `scripts/run_cross_platform_contract_loop.sh` 與 `Makefile` 的 `cross-platform-loop` target 屬既有 `cross-platform-test-contracts` 規格「Regression loops prove cross-platform test stability」需求的交付，不納入本變更
- 前端測試重構（workbench-v2 harness/mocks）與嵌入式前端資產清理屬其他變更批次，不在本範圍
- 不為 macOS / ARM 上的 PowerShell 行為立約（僅 Windows pwsh + Linux bash 兩軸）

## Capabilities

### New Capabilities

- `start-script-contracts`: 啟動腳本跨平台行為契約——開發模式啟動順序、啟動鎖取得/釋放/過期清理、受管端口集合展開、進程樹終止不留孤兒、執行期日誌降噪——在 Windows pwsh 與 Linux bash 兩軸以可執行契約測試釘住

### Modified Capabilities

(none)

## Impact

- Affected specs: `start-script-contracts`（新增）
- Affected code:
  - Modified: start.ps1、start.sh、.air.toml、scripts/start-log-utils.sh、tests/shell/start-backend-before-frontend.sh、tests/shell/start-backend-cleanup-no-log-wait.sh、tests/shell/start-backend-cleanup-order.sh、tests/shell/start-env-load.sh、tests/shell/start-lock-management.sh、tests/shell/start-port-management.sh
  - New: scripts/start-process-utils.ps1、scripts/start-port-utils.ps1、scripts/start-log-utils.ps1、tests/powershell/start-dev-mode-order.ps1、tests/powershell/start-lock-management.ps1、tests/powershell/start-log-noise.ps1、tests/powershell/start-managed-ports.ps1、tests/powershell/start-process-tree.ps1
