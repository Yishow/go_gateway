# 開發熱重載與單檔打包檢視

範圍：`start.sh`、`start.ps1` 及直接依賴的啟動／建置工具。原選單收斂為兩項：開發熱重載與單檔打包。未修改 Go／React 產品行為、既有資料或設備設定；此文件記錄當時尚未提交的驗證；目前提交前覆核見 [分批驗證紀錄](remaining-wip-review-20261025.md)。

## 檢視與修正

- 原 `go run` 開發模式需手動重啟，改為 Air + Vite；保留舊開發旗標為別名。
- Air v1.63.0 使用 `build.bin`；原 `entrypoint` 被忽略，實測會嘗試執行不存在的 `tmp/main`。已改正並固定工具版本。
- Air 原本排除 `lib`，已納入；SQL、go.mod、go.sum 也會觸發重編。編譯失敗停止過期後端，修正儲存後自動恢復。
- macOS 內建 Bash 3.2 的關聯陣列、`BASHPID` 與部分中文變數插值問題已移除或修正。
- 根目錄依腳本位置解析；使用本機 npm／Vite，不再要求 pnpm 或透過 npx 隱式下載工具。
- fresh clone 只有 embed placeholder，但後端啟動會檢查內嵌首頁；缺少首頁時先建置一次前端供後端啟動，之後使用 Vite 即時開發。
- 不再自動終止占用連接埠的其他程序。監控前後端，退出時只清理本次啟動的程序及鎖。
- 移除混在啟動器內的 lint、測試、格式化、診斷等選項；需要時仍可使用原本專用命令。
- 單檔建置共用 `scripts/build-single.mjs`：前端建置成功且首頁有效後，同步完整 static，再以指定 GOOS／GOARCH 與 `CGO_ENABLED=0` 編譯。
- 產物名稱標示平台與架構；後端編譯失敗不覆蓋上一份成功產物。

## 本次實際驗證

環境：macOS ARM64、Go 1.27.3、Node 24.22.1、系統 Bash 3.2；PowerShell 檢查使用暫存的 macOS PowerShell 7.5.4。

- PASS：`node --test scripts/tests/build-single.test.mjs`，20 個測試；涵蓋六個目標、完整資產同步、過期檔案移除、placeholder、缺少輸入、安裝／建置失敗及開發 bootstrap。
- PASS：以 `/bin/bash` 執行全部 13 支 `tests/shell/start-*.sh`。
- PASS：PowerShell parser 與 9 組 focused suites，涵蓋 bootstrap、順序、鎖、埠占用、log、建置失敗退出碼及環境還原；從 `/tmp` 實際執行 `start.ps1 -BuildSingle -OS mac -Arch arm64` 成功。
- PASS：`python3 -m unittest scripts/tests/test_windows_build_contract.py`，5 個既有建置契約測試。
- PASS：實際產生 Windows、macOS、Linux 各 amd64／arm64 共六份執行檔，並以 `file` 核對 PE、Mach-O、ELF 及 CPU 架構。產物與建置 log 在忽略的 `bin/` 內。
- PASS：只複製 macOS ARM64 執行檔至空白暫存目錄，使用可丟棄 SQLite；health API、內嵌首頁與兩個首頁資產回傳 200，SIGTERM 後 exit 0 並釋放埠。
- PASS：`scripts/tests/launcher-live-smoke.py` 的實測流程；隔離副本從 placeholder 啟動、使用前端代理、覆寫錯誤的交叉編譯環境、修改 lib 自動更換後端 PID、編譯失敗後修正恢復、WebSocket 收到 Vite HMR update，以及 TERM 後釋放雙埠與鎖。

PASS：`git diff --check`、`make check-lines`。兩支啟動腳本僅有超過 300 行的警告，未超過 500 行阻擋限制。PowerShell mock／parser 不能替代 Windows 原生執行。

## 驗證限制

- NOT_RUN：Windows 原生 PowerShell／程序樹、Linux 原生啟動、Intel macOS 原生啟動。
- NOT_RUN：整個產品的 Go test／vet／lint、前端 lint／Vitest／Playwright；本次未修改產品程式碼，已實際執行完整前端建置與六個 Go 目標編譯。
- NOT_RUN：真實 PLC、正式資料庫與現場採集驗收。
- 單一執行檔包含前後端程式與網頁資產；執行時的 SQLite、log 與使用者設定仍是外部資料檔。

使用命令見 [README](../../README.md#本機開發與單檔打包)。相關既有契約見 [embedded-frontend-delivery](../../openspec/specs/embedded-frontend-delivery/spec.md)。

## 選單與瀏覽器詢問補充

- 無參數執行兩支腳本皆顯示「開發模式／建置單一執行檔／離開」選單。
- 開發模式詢問「是否開啟瀏覽器？[y/N]」；Enter、N 或無可用輸入預設不開，Y 才在前後端就緒後開啟 Vite 頁面。啟動失敗不開啟。
- Air 子程序固定 `AUTO_OPEN_BROWSER=false`，避免原設定或後端熱重載造成重複開啟；保留父程序原環境。
- PASS：14 組系統 Bash 3.2 測試、10 組 macOS PowerShell 測試、兩支實際無參數選單入口、diff 與行數檢查。瀏覽器開啟使用 mock 驗證命令與時機，未實際操作桌面瀏覽器，Windows 原生執行仍為 NOT_RUN。
