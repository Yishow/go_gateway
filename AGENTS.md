<!-- SPECTRA:START v1.3.0 -->

# Spectra Instructions

This project uses Spectra for Spec-Driven Development(SDD). Specs live in `openspec/specs/`, change proposals in `openspec/changes/`.

## Use `$spectra-*` skills when:

- An explicitly requested Spectra decision needs structure → `$spectra-discuss`
- User explicitly requests a Spectra change proposal → `$spectra-propose`
- Continue tasks for an identified change → `$spectra-apply`
- Update requirements or plans for an identified change → `$spectra-ingest`
- Check implementation matches artifacts → `$spectra-verify`
- Review implementation quality and logic issues → `$spectra-review`
- Analyze artifact consistency before coding → `$spectra-analyze`
- Audit security sharp edges → `$spectra-audit`
- Check stale changes before resuming → `$spectra-drift`
- Debug a concrete bug systematically → `$spectra-debug`
- Implementation is done → `$spectra-archive`
- Commit only files related to a specific change → `$spectra-commit`

Explicit skill invocation takes precedence. Apply existing authorization within its unchanged scope.

## Workflow

discuss? → propose → apply ⇄ ingest → verify / review → archive

- `discuss` is optional — skip if requirements are clear
- Requirements change mid-work? Plan mode (`/plan`) → `ingest` → resume `apply`

## Parked Changes

Changes can be parked（暫存）— temporarily moved out of `openspec/changes/`. Parked changes won't appear in `spectra list` but can be found with `spectra list --parked`. To restore: `spectra unpark <name>`. The `$spectra-apply` and `$spectra-ingest` skills disclose parking and restore when the named operation is already explicitly requested; respect a known refusal, otherwise ask for missing authorization.

<!-- SPECTRA:END -->

# Repository Guidelines

## 先讀與優先順序

- 執行任務前同時讀本檔與 `CLAUDE.md`；本檔是所有 agent 的共通規範，`CLAUDE.md` 只補充執行脈絡。
- 先確認目前 branch、工作目錄差異與遠端基準；保留使用者／其他工作者的未提交變更，不整包覆寫或 force push。
- 專案規則以本檔為準，再讀 agent 補充、相關 source、測試與 OpenSpec。規格與程式不符時分別記錄，不把其中一方直接當成已完成證據。
- 本版 repo **沒有 tracked `.github/instructions/` 或 `.github/workflows/`**；不要引用不存在的細部規範、CI gate 或過往機器的私人備份路徑。
- 任務若涉及 Studio surfaces，先讀 `docs/technical/studio-surface-inventory/START_HERE.md`、`context.json`、`CURRENT_STATE.md`。日期較舊的 inventory 與 release 記錄是脈絡，現況需再核對 entrypoint、service wiring、呼叫端及測試。

## 專案與產品邊界

- 工業資料採集閘道：設備連線 → 採集點位 → Point-to-Tag 映射 → 輸出。保留 Go + Gin 後端與 React + TypeScript 前端，在現有架構內改進。
- 版本來源：`go.mod` 的 Go 1.25.5、`frontend/package.json`／lockfile 的 React 19、TypeScript 5、Vite 7。不要從舊規劃文件推定版本。
- 主程式 `cmd/test_ui` 提供 API 及嵌入式 SPA；`cmd/test_ui/main.go` 宣告 embed，`internal/web/embed.go` 提供靜態服務。
- `/studio/v2` 是使用者 setup 主線；`/studio/runtime` 是 setup 後的 focused monitor；`/test` 為工程工具；`/gateway/*` 為 experimental。
- 不建立平行 V3／完整重寫，不恢復 `/studio` 專用 route、handler、tombstone 或 special redirect。其他未知路由依現有 generic policy 處理。
- 基礎目標是把所選 Tag 分組寫入資料庫；報表、保留政策、聚合與區間電量等衍生量不是基本寫入的必備前置。
- Local Modbus 與 Database 是目前 runtime 輸出主線；Database 以 SQLite／PostgreSQL 為現有支援範圍。依賴中有 driver 不等於 setup、runtime 和安全契約已全部支援。

## 現況判斷與完成標準

- 實際 runtime wiring 先查 `cmd/test_ui/service_wiring.go`、`target_writer.go` 及 `main.go`，不可只看 package、route 或 UI 按鈕存在。
- 實作須完成已同意的驗收範圍，不只交付最小可動版本；錯誤處理與邊界條件屬任務範圍，不是默認後續工作。
- 區分：已有程式、已接 production entrypoint、已有測試檔、此次實際測試通過、現場驗收。不同層級不能互相代替。
- 目前 dbtarget writer 可實際寫入 SQL；grouped buffer、delivery／aggregation 元件與 recording-plan API 不可因此宣稱已有端到端 durable recording。
- `recording-plans/test-write` 目前回傳明確 501；直到真實 write/readback/cleanup 流程完成並驗證前，必須維持安全失敗，不以範例資料或固定成功動畫補位。
- 保留已落地的 workspace/settings/connector revisions、readiness token、schema preview/confirm/apply、operation ledger、scope ownership 與 Modbus Share fail-closed 保護。
- 修改設定、測試成功、runtime running、buffer accepted、SQL committed、readback verified 與 cleanup completed 必須分別呈現。

## 目錄定位

- `cmd/`：可執行入口與 wiring；`cmd/test_ui` 是主要入口。
- `internal/api/`：HTTP／SSE 路由、handlers 與 API 組裝。
- `internal/datalink/`：device、point、tag、mapping、collector、runtime、workspace、dbtarget、measurement、recordingplan、aggregation、delivery 等模組。
- `internal/datalink/schema/migrations/`：內部設定資料庫 migration。
- `internal/protocol/`、`lib/`：協議與共用實作。
- `frontend/src/`：React、types、services、i18n；`frontend/tests/`：正式前端測試入口。
- `docs/`：技術、API、release 證據；`openspec/specs/`：已合併規格；`openspec/changes/`：進行中提案。
- `.agents/skills/`：repo-local Spectra skills；工具目錄未列出時直接讀對應 `SKILL.md`，不臆造等效流程。

## Spectra 與規格文件

- 明確要求 Spectra proposal 時使用 repo-local `spectra-propose`；指定其他 workflow 時先核對實際 skill 與 CLI，不能悄悄替換。先確認真實 CLI 可用，再依 CLI template／locale／dependency order 產生 artifacts。
- 提案不等於實作；所有尚未執行的產品任務保持未完成，不因文件交付而 archive 或勾選。
- 先核對 active、parked、archived changes 與相關 specs；處理重疊時留下唯一負責 change、依賴與移交原因，保留原有完成證據。
- 不直接改寫既有規格語意而不留 delta／衝突說明。擴充採 expand → migrate → contract，每批可建置與驗證。
- 路徑為 repo-relative；工作項目需寫明可觀察行為和驗收，不能只有「修改某檔」。數字區分測量、設計條件與估算。
- 文件任務至少執行 `git diff --check` 與 `make check-lines`；規格任務另依 skill 執行 analyze／validate，明列未執行或失敗的檢查。
- 修改 `docs/technical/studio-surface-inventory/` 內文件時，必須以 `go run ./cmd/studio_inventory_changelog ...` 同步其 `changelog.sqlite`，至少記錄 summary、surface、files、reason。

## 建置與開發

- `make build`：npm install → frontend build → 複製至 `cmd/test_ui/static` → Go build 到 `bin/test-ui.exe`；副檔名不表示 cross-compile。
- `powershell -File scripts/build.ps1`：Windows 完整建置。輸出路徑與參數以腳本為準。
- `make test-ui`：建置並執行；`make clean` 會刪除 frontend/dist、frontend/node_modules、bin。
- `make gen-docs`：用已安裝的 swag 依 annotations 重產 `docs/swagger/`；不要手改 generated Swagger。
- `cd frontend && npm ci`：按 lockfile 安裝；`npm run dev`：Vite dev server；API backend 需另啟動。
- 設定讀取以 `internal/config/config.go` 的環境變數／dotenv 為準，不把其他入口的 Viper/YAML 說成主要 server 的設定方式。
- fresh clone 只有 embed placeholder；需要可用 UI 時先完成前端 build 和 static sync，不能只 go build 就宣稱已有網頁。

## 驗證

- 後端行為變更最低基準：`go test ./...`、`go vet ./...`、`golangci-lint run ./...`。
- 前端行為變更最低基準：`cd frontend && npm run lint && npm test -- --run && npm run build`。
- 前端 UI／UX 採測試先行；至少涵蓋受影響的 Studio 主流程與 TestPage 互動。協議、映射、排程、lifecycle、output binding 都需相應回歸。
- Go 測試維持與 package 相鄰；前端分為 `frontend/tests/unit/`、`integration/`、`e2e/`。
- `cd frontend && npm run test:e2e`：Playwright；gateway 專用 `test:gateway:unit`、`test:gateway:e2e`、`test:gateway:gate` 不代替產品主線驗收。
- `make cross-platform-loop`：repo 的跨平台契約序列；名稱不能代替實際執行平台證據。
- `make gate-smoke`、`gate-final`、`gate-soak`、`gatev11`、`gatev12`：Modbus gate／soak，執行前檢查腳本的 target 與測試環境。
- `make points-precheck-up/down`、`points-migrate-up/down`：資料 migration；先檢查目標、備份、衝突與 down 策略，禁止拿正式資料試跑。
- `make longtask-smoke`：長任務提醒機制檢查，不是資料採集驗收。
- 不宣稱「全通過」而只執行 focused suite；明列命令、結果、環境、未執行項目與對應風險。歷史測試數不能作為此次結果。

## 程式碼與安全

- Go 使用 gofmt／goimports、小寫 package、明確錯誤回傳及 `fmt.Errorf("context: %w", err)`；避免同一錯誤重複 log／return，exported symbol 補文件。
- TypeScript 使用 strict、ES modules、2 空白縮排；元件 PascalCase、hooks useX；避免 any，採 unknown + narrowing。
- types 在 `frontend/src/types/`、API 在 services、server state 用既有 React Query hooks、文字進 en／zh-TW locales；先重用既有元件和 service。
- 非同步失敗必須顯示可行動的 error state，不能 silent failure；普通 operator DOM 不洩漏 endpoint／credential／raw backend exception。
- 所有外部輸入需驗證與正規化；SQL 走既有 abstraction 及參數化查詢，識別字採明確驗證／dialect quoting。
- 不提交真實帳密、Token、DSN、設備憑證、私人路徑或其他 secrets；密碼不任意 trim，遮蔽值不是新密碼。
- 對外 probe／寫入由 backend 主機發起；不得誤導成瀏覽器直連。
- 測試使用 simulator 及可丟棄資料庫；真 PLC、正式 DB、SCADA、LAN 與部署操作需個別授權及驗收。
- 憑證、權限、破壞性 migration 或外部建表／試寫不得繞過既有確認與 ownership 保護。

## 行數、提交與交接

- `scripts/check_file_lines.sh` 是實際執行規則，搭配 `.line-limit-ignore`；包含 Markdown 等文字檔，沒有通用 MD 豁免。
- 300 行以上警告，500 行以上阻擋；歷史超長檔只准不增行，需記錄拆分計畫。產生檔與略過目錄依腳本／ignore 清單。
- `make check-lines` 預設看目前差異；已 commit 的完整批次用 `bash scripts/check_file_lines.sh --base <base-sha>`。
- `.githooks/pre-commit` 可透過 `git config core.hooksPath .githooks` 啟用；本 repo 未提供 tracked GitHub Actions workflow，不能宣稱 CI 已執行此 gate。
- commit 主旨使用簡短繁中祈使句；只 stage 此次範圍。未授權不得 commit／push／merge／部署。
- PR／交接列出修改檔案、影響、OpenSpec 連結、實際驗證、風險及後續；有 UI 變更附實測截圖。
- 發布前重新檢查遠端進度，保留並整合 concurrent updates；發布後核對 remote SHA 與該 SHA 的 CI 狀態。
