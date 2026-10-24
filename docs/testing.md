# 測試與驗證指引

本文件由 [AGENTS.md](../AGENTS.md) 按需載入，適用於修改檔案、驗證及合併／發布檢查。純問答與唯讀檢視不套用實作驗證流程。

## 執行時機

| 時機 | 應執行的檢查 |
| --- | --- |
| 開發、修復迭代 | 直接受影響的測試；修 bug 先補能重現問題的回歸測試，再修復並確認通過。 |
| 本次任務交付 | 受影響模組、呼叫端與相關整合測試，加上適用的 lint、型別與建置檢查。 |
| 合併、發布前 | 對最終批次執行涉及前端／後端的完整檢查；平台、資料庫、E2E 依本次驗收範圍執行。 |

- 開發迭代與局部任務交付預設選測；不因修改一個檔案就重跑前後端全套。
- 使用者明定的驗收及適用 skill 的必要檢查仍須完成，不能用選測取代。
- 檢查通過後，只有後續變更可能影響結果、檢查失敗或有具體未解風險才重跑或擴大驗證；完整檢查對最終批次執行一次即可。
- 若完整檢查後又修改相關程式、依賴或設定，先前結果不能代表最終批次；依影響重跑必要檢查。

## 選測與擴大條件

- 依變更行為、依賴與呼叫端選測，不能只測修改檔案旁邊的測試。執行前簡述受影響範圍及選測理由。
- 修改共用元件、API 契約、schema、依賴或測試設定時，擴大至相關使用端及整合測試；影響無法界定時跑受影響側的完整檢查。
- 前端 UI／UX 採測試先行，涵蓋本次變更會影響的 Studio 主流程或 TestPage 互動；協議、映射、排程、lifecycle、output binding 需相應回歸。
- 型別變更須檢查相關使用端；純型別檢查可在 `frontend/` 執行 `npx tsc --noEmit`。涉及打包、路由載入、資源或嵌入式交付時，須驗證相關 build 與交付流程。
- 選測命令必須實際選到預期測試；零測試、skip 或環境缺失不能算通過。
- 既有測試失敗先釐清與本次變更的關係；範圍外問題回報，不自動擴大實作或新增驗收條件。

## 常用選測入口

以下命令從 repo 根目錄執行，標示前端者從 `frontend/` 執行。表中為選測起點，需依變更補上呼叫端與整合測試；未列出的模組依實際 package／測試路徑選擇。

| 範圍 | 指令與對應測試 |
| --- | --- |
| Go 模組，例如 mapping | `go test ./internal/datalink/mapping`；測試維持與 package 相鄰。 |
| API handlers／router | `go test ./internal/api/handlers ./internal/api`；可在迭代時以 `-run` 限縮至受影響測試。 |
| Go 靜態檢查 | `go vet ./internal/datalink/mapping`、`golangci-lint run ./internal/datalink/mapping`；依實際影響替換或增加 package。 |
| 前端單元／整合測試 | 前端：`npm test -- --run <測試路徑>`；位於 `tests/unit/`、`tests/integration/`。 |
| 前端檔案 lint | 前端：`npx eslint <檔案路徑...>`；包含本次受影響的 source 與測試。 |
| 前端 E2E | 前端：`npm run test:e2e -- <測試路徑>`；位於 `tests/e2e/`。 |
| Gateway 專用檢查 | 前端：`npm run test:gateway:unit`、`npm run test:gateway:e2e`、`npm run test:gateway:gate`；不代替產品主線驗收。 |

例如，只迭代 mapping autosave 的 queue 行為時，可在 `frontend/` 執行：

```sh
npm test -- --run tests/unit/workbench-v2/mapping-autosave-queue
```

交付時再依影響補上 page、service 或整合測試；不能由 queue 測試通過推論整個 Studio 流程通過。

### macOS 本機瀏覽器驗證

- 首次執行 E2E 前，在 `frontend/` 執行 `npx playwright install chromium`，安裝相符的專用瀏覽器與錄影工具；headless 測試沿用既有設定的 headless shell。缺少套件時先補安裝，不改用 `channel: 'chrome'` 啟動日常使用的 Chrome。[Playwright 官方說明](https://playwright.dev/docs/browsers#chromium-headless-shell)
- 在受限執行環境中，瀏覽器啟動須使用工具提供、經核准的主機執行權限；不要先嘗試啟動系統 Chrome。曾實測 macOS 應用程式註冊階段因限制而 SIGABRT，並產生「未預期結束」彈窗。
- `agent-browser` 使用獨立工作階段及專用 headless shell，可透過 `--executable-path` 或本機設定的 `executablePath` 指定。瀏覽器版本更新後若路徑失效，重新安裝並核對設定；不退回系統 Chrome。

## 完整檢查

在合併／發布前，或符合上述擴大條件時，執行本批次涉及側的完整檢查。純文件修改適用下一節。

後端從 repo 根目錄執行：

```sh
go test ./...
go vet ./...
golangci-lint run ./...
```

前端從 `frontend/` 執行：

```sh
npm run lint
npm test -- --run
npm run build
```

E2E、真實資料庫、Windows／設備及現場驗收依本次範圍選擇，不由上述命令的成功推論通過。

## 純文件與規格修改

- 檢查內容、連結及與相關指引的一致性，至少執行 `git diff --check` 與 `make check-lines`；不必執行產品測試、lint 或 build。
- 規格任務另依適用 skill 執行 analyze／validate，未執行的產品任務維持未完成。
- 修改 Studio inventory 文件時，仍須依 AGENTS.md 同步其 changelog。

## 特定環境與驗收入口

只在本次變更或驗收需要時讀取相關腳本並執行，不將以下入口全部加入日常修改檢查：

- `make cross-platform-loop`：repo 的跨平台契約序列；名稱不能代替實際執行平台證據。
- `make gate-smoke`、`gate-final`、`gate-soak`、`gatev11`、`gatev12`：Modbus gate／soak，執行前檢查腳本的 target 與測試環境。
- `make points-precheck-up/down`、`points-migrate-up/down`：資料 migration；先檢查目標、備份、衝突與 down 策略，禁止拿正式資料試跑。
- `make longtask-smoke`：長任務提醒機制檢查，不是資料採集驗收。
- 使用 simulator 及可丟棄資料庫；真 PLC、正式 DB、SCADA、LAN 與部署操作仍需個別授權及驗收，不繞過既有確認與 ownership 保護。

## 結果與停止條件

- 回報實際命令、範圍、環境、PASS／FAIL／NOT RUN 與未執行項目的風險；工具回報 cached 時如實註明。
- 只執行 focused suite 時，不宣稱「全通過」；歷史測試數不能作為此次結果。
- 區分本機測試、真實資料庫、UI、Windows／設備與現場驗收證據，不能互相代替。
- 原定驗收與必要檢查滿足後即交付，不因存在更多可跑的測試就持續擴大驗證。
