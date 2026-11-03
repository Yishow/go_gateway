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

- 本檔是所有 agent 的唯一共通規範；`CLAUDE.md` 只以 `@AGENTS.md` 匯入本檔，不另寫規則。agent 專屬脈絡確有需要時才放各自檔案，且不得與本檔衝突。
- 先確認任務層級（研究、proposal、實作、文件整理），只做獲授權的層級。
- 先確認目前 branch、工作目錄差異與遠端基準；保留使用者／其他工作者的未提交變更，不整包覆寫或 force push。
- 專案規則以本檔為準，再讀相關 source、測試與 OpenSpec。規格與程式不符時分別記錄，不把其中一方直接當成已完成證據。
- 本版 repo **沒有 tracked `.github/instructions/` 或 `.github/workflows/`**；不要引用不存在的工具、細部規範、CI gate 或過往機器的私人備份路徑。
- 任務若涉及 Studio surfaces，先讀 `docs/technical/studio-surface-inventory/START_HERE.md`、`context.json`、`CURRENT_STATE.md`，再開實際 route、handler、service wiring 與測試。日期較舊的 inventory、release 記錄與 archived tasks 的勾選只是脈絡，不能據以推論 production 已接線。

## 專案與產品邊界

- 工業資料採集閘道：設備連線 → 採集點位 → Point-to-Tag 映射 → 輸出。保留 Go + Gin 後端與 React + TypeScript 前端，在現有架構內改進。
- 版本來源：`go.mod` 的 Go 1.25.5、`frontend/package.json`／lockfile 的 React 19、TypeScript 7、Vite 8；TypeScript 6 alias 僅供 lint API 使用。不要從舊規劃文件推定版本。
- 主程式 `cmd/test_ui` 提供 API 及嵌入式 SPA；`cmd/test_ui/main.go` 宣告 embed，`internal/web/embed.go` 提供靜態服務。
- `/studio/v2` 是使用者 setup 主線；`/studio/runtime` 是 setup 後的 focused monitor；`/test` 為工程工具；`/gateway/*` 為 experimental。
- 不建立平行 V3／完整重寫，不恢復 `/studio` 專用 route、handler、tombstone 或 special redirect。其他未知路由依現有 generic policy 處理。
- 基礎目標是把所選 Tag 分組寫入資料庫；報表、保留政策、聚合與區間電量等衍生量不是基本寫入的必備前置。
- Local Modbus 與 Database 是目前 runtime 輸出主線；Database 目標涵蓋 SQLite、PostgreSQL、MySQL／MariaDB、SQL Server，各 dialect 以 `internal/datalink/dbtarget` 的 capability 檢查為準。有 driver 或 dialect 分支不等於該資料庫的 setup、runtime 與安全契約已驗收。

## 現況判斷與完成標準

- 實際 runtime wiring 先查 `cmd/test_ui/service_wiring.go`、`target_writer.go` 及 `main.go`，不可只看 package、route 或 UI 按鈕存在。
- 實作須完成已同意的驗收範圍，不只交付最小可動版本；錯誤處理與邊界條件屬任務範圍，不是默認後續工作。
- 區分：已有程式、已接 production entrypoint、已有測試檔、此次實際測試通過、現場驗收。不同層級不能互相代替。
- Write group 的 durable delivery（`groupdelivery`／`grouppipeline`）與明確試寫（`grouptestwrite`）已接到 `service_wiring.go`；wiring 存在不等於特定資料庫或設備組合已通過現場驗收。
- 試寫必須真實執行 write/readback/cleanup，並如實呈現未決或失敗狀態；service 未接時回報 capability unavailable，不以範例資料或固定成功動畫補位。
- 保留已落地的 workspace/settings/connector revisions、readiness token、schema preview/confirm/apply、operation ledger、scope ownership 與 Modbus Share fail-closed 保護；不能以簡化 UI 為理由移除。
- 修改設定、測試成功、runtime running、buffer accepted、SQL committed、readback verified 與 cleanup completed 必須分別呈現。

## 調查方式

- 沿使用者操作追到 persisted identity，再追 runtime target writer；UI、API 與 SQL 證據要能對上同一 workspace／revision。
- 遇到 competing model 時，先寫明誰是 authority、哪些只是 projection／相容 adapter，避免新增一套平行狀態。
- 測試與 source 不一致時保留事實及風險，不為了讓完成報告好看而改掉預期。

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

- 明確要求 Spectra proposal 時使用 repo-local `spectra-propose`；指定其他 workflow 時先核對實際 skill、CLI 與既有 change，不能悄悄替換，也不用舊機器記憶或自製格式代替。先確認真實 CLI 可用，再依 CLI template／locale／dependency order 產生 artifacts。
- 提案不等於實作；proposal 完成只表示規格可交接。所有尚未執行的產品任務保持未完成，不因文件交付而 archive 或勾選。
- 先核對 active、parked、archived changes 與相關 specs；處理重疊時留下唯一負責 change、依賴與移交原因，保留原有完成證據。
- 不直接改寫既有規格語意而不留 delta／衝突說明。擴充採 expand → migrate → contract，每批可建置與驗證。
- 路徑為 repo-relative；工作項目需寫明可觀察行為和驗收，不能只有「修改某檔」。數字區分測量、設計條件與估算。
- 規格任務另依 skill 執行 analyze／validate，明列未執行或失敗的檢查。
- 修改 `docs/technical/studio-surface-inventory/` 內文件時，必須以 `go run ./cmd/studio_inventory_changelog ...` 同步其 `changelog.sqlite`，至少記錄 summary、surface、files、reason。

## 建置與開發

- `make build`：npm install → frontend build → 複製至 `cmd/test_ui/static` → Go build 到 `bin/test-ui.exe`；副檔名不表示 cross-compile。
- `powershell -File scripts/build.ps1`：Windows 完整建置。輸出路徑與參數以腳本為準。
- `make test-ui`：建置並執行；`make clean` 會刪除 frontend/dist、frontend/node_modules、bin。
- `make gen-docs`：用已安裝的 swag 依 annotations 重產 `docs/swagger/`；不要手改 generated Swagger。
- `cd frontend && npm ci`：按 lockfile 安裝；`npm run dev`：Vite dev server；API backend 需另啟動。
- 設定讀取以 `internal/config/config.go` 的環境變數／dotenv 為準，不把其他入口的 Viper/YAML 說成主要 server 的設定方式。
- fresh clone 只有 embed placeholder；需要可用 UI 時先完成前端 build 和 static sync，不能只 go build 就宣稱已有網頁。

## 驗證（按需載入）

- 修改程式、設定、測試或文件前，或任務要求驗證、合併／發布檢查時，先讀 [測試與驗證指引](docs/testing.md)，依任務選用對應章節。
- 純問答與唯讀檢視不必載入；只有選測或查證測試結果需要時才讀取。

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
- `.githooks/pre-commit` 可透過 `git config core.hooksPath .githooks` 啟用；沒有 CI，不能宣稱此 gate 已在 CI 執行。
- commit 主旨使用簡短繁中祈使句；只 stage 此次範圍。未授權不得 commit／push／merge／部署。
- 每次完成回報與 PR／交接列出修改檔案、影響、OpenSpec 連結、已執行／未執行的驗證、風險及後續建議；有 UI 變更附實測截圖。
- 發布前重新檢查遠端進度，保留並整合 concurrent updates；發布後核對 remote SHA。
