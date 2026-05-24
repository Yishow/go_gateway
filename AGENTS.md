<!-- SPECTRA:START v1.0.2 -->

# Spectra Instructions

This project uses Spectra for Spec-Driven Development(SDD). Specs live in `openspec/specs/`, change proposals in `openspec/changes/`.

## Use `$spectra-*` skills when:

- A discussion needs structure before coding → `$spectra-discuss`
- User wants to plan, propose, or design a change → `$spectra-propose`
- Tasks are ready to implement → `$spectra-apply`
- There's an in-progress change to continue → `$spectra-ingest`
- User asks about specs or how something works → `$spectra-ask`
- Implementation is done → `$spectra-archive`
- Commit only files related to a specific change → `$spectra-commit`

## Workflow

discuss? → propose → apply ⇄ ingest → archive

- `discuss` is optional — skip if requirements are clear
- Requirements change mid-work? `ingest` → resume `apply`

## Parked Changes

Changes can be parked（暫存）— temporarily moved out of `openspec/changes/`. Parked changes won't appear in `spectra list` but can be found with `spectra list --parked`. To restore: `spectra unpark <name>`. The `$spectra-apply` and `$spectra-ingest` skills handle parked changes automatically.

<!-- SPECTRA:END -->

# Repository Guidelines

## 專案總覽
- 專案目標：工業資料採集閘道，從 PLC 協議讀取資料，經 Datalink 映射後輸出到資料儲存或訊息系統，並提供嵌入式 Web UI。
- 技術棧：後端為 Go 1.25.x（目前 `go.mod` 為 1.25.5）+ Gin；前端為 React 19 + TypeScript 5 + Vite 7。
- 部署型態：以單一可執行檔整合 API 與嵌入式前端，主要入口是 `cmd/test_ui`；前端 build 產物會嵌入 `cmd/test_ui/static`。
- 產品主線：`/studio` 是目前 datalink 主產品入口；舊 datalink legacy routes 應收斂 redirect 到 `/studio`，工程測試工具則統一收斂到 `/test`。
- 目前輸出主線以 Local Modbus 與 Database 為主，其中 Database output 正式支援 `SQLite` 與 `PostgreSQL`。

## 文件閱讀要求（強制）
- Agent 在本 repo 執行任務時，**必須同時閱讀** `AGENTS.md` 與對應 Agent 專屬文件（例如 `CLAUDE.md`）。
- 不可只讀單一文件；`CLAUDE.md` 不能取代 `AGENTS.md`，`AGENTS.md` 也不能取代 `CLAUDE.md`。
- 若兩者描述重疊或衝突，以 `AGENTS.md` 為準，再由 Agent 專屬文件補充執行脈絡。

## 專案結構與模組分工
- `cmd/`：可執行入口（如 `cmd/test_ui`、`cmd/test_all`、`cmd/fatek_test`）。
- `internal/api/`：HTTP / SSE 路由與 handlers，含 Swagger 掛載與 Datalink service 組裝。
- `internal/datalink/`：核心業務模組（device、point、tag、mapping、pollinggroup、collector、connector、runtime、dbtarget、sourcerule 等）。
- `internal/protocol/`：協議客戶端與傳輸層（Modbus、FATEK、MC Protocol 等）。
- `internal/datalink/schema/migrations/`：資料表與 migration。
- `lib/`：共用協議與演算法函式庫。
- `frontend/`：Vite + React + TypeScript 前端（`src/`、`public/`、`tests/`、ESLint / Vitest / Playwright 設定）。
- `docs/`：架構、API 與維運文件（含 `docs/swagger/`）。
- `openspec/`：需求規格、變更提案與任務追蹤；重大變更請遵循此流程。

## 建置、測試與開發指令
- `make build`：安裝前端依賴、建置前端，再建置 `bin/test-ui.exe`。
- `powershell -File scripts/build.ps1`：Windows 本機完整建置（前端 build → 複製 embed static → 後端 build）。
- `make clean`：清理 `frontend/dist`、`frontend/node_modules` 與 `bin/`。
- `make gen-docs`：重新產生 Swagger 文件到 `docs/swagger/`。
- `go test ./...`：執行所有後端單元 / 整合測試。
- `go vet ./...`：執行 Go 官方靜態檢查。
- `golangci-lint run ./...`：依 `.golangci.yml` 執行 Go 靜態分析（含 `gosec`、`errorlint`、`revive` 等）。
- `make test-ui`：建置後直接執行 `bin/test-ui.exe`。
- `cd frontend; npm run dev`：啟動前端開發伺服器。
- `cd frontend; npm run lint`：執行 ESLint。
- `cd frontend; npm run test`：執行 Vitest。
- `cd frontend; npm run build`：執行 `tsc && vite build`。
- `cd frontend; npm run test:e2e`：執行 Playwright E2E。
- `cd frontend; npm run test:gateway:unit | npm run test:gateway:e2e | npm run test:gateway:gate`：執行 gateway 專用前端驗證流程。
- `make gate-smoke | make gate-final | make gate-soak | make gatev11 | make gatev12`：執行 Modbus gate / soak 類驗證腳本。
- `make points-precheck-up | make points-precheck-down | make points-migrate-up | make points-migrate-down`：執行 points unique migration precheck / apply；down 流程需明確確認策略。
- `make longtask-smoke`：驗證 long-task 提醒機制（controller start/finish + reconcile）。
- `make check-lines`：執行檔案行數規範檢查（300 警告 / 500 阻擋）。

## 程式碼樣式
- Go：使用 `gofmt` / `goimports`；package 名稱維持小寫；錯誤處理優先明確回傳並附上下文；維持 happy path 左對齊。
- Go：新增或改動 exported symbol 時補齊文件註解；避免重複 `package` 宣告；優先重用標準函式庫與既有 service / helper。
- Frontend：以 TypeScript 為主，2 空白縮排；元件檔用 `PascalCase.tsx`，hooks 用 `useX.ts`。
- Frontend：遵守 ES modules、`strict` TypeScript 與既有 ESLint 規則；避免 `any`，必要時以 `unknown` + narrowing 或既有型別重用替代。
- 前端資料流：型別優先更新 `frontend/src/types/`；API 呼叫集中於 `frontend/src/services/`；伺服器狀態透過 React Query hooks；使用者文字統一走 `frontend/src/i18n/locales/`。
- 維持 import、未使用程式碼與 lint 警告乾淨；警告視為待處理工作，不應帶入主分支。

## 命名規則
- Go package：小寫、單數語意，避免 `util` / `common` / `base` 之類泛稱。
- Go 型別 / 函式：匯出符號使用 `PascalCase`，非匯出符號使用 `camelCase`。
- Go 介面：優先以行為命名（例如 `Reader`、`Writer`）。
- Frontend 元件：檔名 `PascalCase.tsx`，hooks 檔名 `useX.ts`。
- 測試檔命名：Go 使用 `*_test.go`；前端使用 `*.test.ts` 或 `*.test.tsx`。

## Error Handling Pattern
- Go：函式呼叫後立即檢查錯誤；錯誤往上傳遞時使用 `fmt.Errorf("...: %w", err)` 附上下文。
- Go：錯誤訊息維持小寫、無句點；需語意判斷時優先 sentinel/custom error + `errors.Is` / `errors.As`。
- Go：避免同一錯誤同時 `log` 並 `return`，在最合適層級做單點處理。
- Frontend：非同步流程必須顯式處理失敗（`try/catch`、error state、toast/fallback），禁止 silent failure。
- Frontend：對使用者顯示的錯誤訊息需可行動、可理解，不洩漏敏感實作細節。

## 測試指引與測試要求
- 後端使用 Go `testing`，必要時搭配 `testify`。
- 前端使用 Vitest + Testing Library（`frontend/src/setupTests.ts`），E2E 使用 Playwright。
- 任何行為變更都要同步新增或調整測試，尤其是協議解析、映射流程、排程邏輯、runtime lifecycle 與 output binding。
- 後端先補精準單元測試，再補資料庫 / 協議邊界的整合測試。
- 前端 UI / UX 調整採 TDD 先行：先補測試，再改介面；至少覆蓋 `Studio` 主流程與 `TestPage` 的主要互動流程。
- 前端測試統一以 `frontend/tests/` 作為正式入口並分類管理：
  - `frontend/tests/unit/`：Vitest 單元與頁面互動測試。
  - `frontend/tests/integration/`：Vitest 跨模組整合測試。
  - `frontend/tests/e2e/`：Playwright 端對端測試。
- Go `*_test.go` 維持與實作檔相鄰，不搬到 root `tests/`，以符合 Go 工具鏈與 package 慣例。
- 前端本地驗證建議至少執行 `cd frontend && npm run lint && npm run test && npm run build`；CI / 非 watch 模式可使用 `npm test -- --run`。
- 後端本地驗證建議至少執行 `go test ./... && go vet ./... && golangci-lint run ./...`。

## 安全考量
- 不要把真實帳號、密碼、Token、DSN、設備憑證或其他 secrets 寫入程式碼、測試、文件、OpenSpec 或 commit message。
- 所有外部輸入（HTTP payload、protocol config、DB connector 設定、檔案內容、query / route params）都要做驗證與正規化；優先重用既有型別、validator、schema 或 service 層檢查。
- 前端顯示使用者 / 外部來源文字時，維持 React 預設 escaping；不要直接拼接不受信任的 HTML，也不要引入動態程式碼執行。
- 資料庫與儲存層修改要走既有 repository / service abstraction 或參數化查詢，避免手刻可注入的 SQL。
- Go 端靜態分析以 `golangci-lint run ./...` 為基準，其中 `gosec` 必須維持乾淨；涉及 shell、檔案路徑、序列埠、網路連線時要特別保守。
- 對外連線測試 / probe 是由 backend 主機發起；文件、錯誤訊息與 UX 說明需避免讓使用者誤判為瀏覽器端直接連線。

## 禁止事項
- 禁止把 secrets（帳密、Token、DSN、憑證）提交到程式碼、測試、文件、OpenSpec 或 commit message。
- 禁止跳過輸入驗證、正規化與參數化查詢，直接以字串拼接方式處理資料庫或指令。
- 禁止把 lint 警告、未使用 import、未使用程式碼帶入主分支。
- 禁止新增繞過 `/studio` 的平行產品入口；legacy 路由應以 redirect 收斂。
- 禁止未對齊 OpenSpec 既有規格就直接改寫需求語意；遇到衝突需先回報。

## 檔案行數規範（強制）
- 目標上限：單檔 `<= 300` 行；硬上限：單檔 `<= 500` 行。
- 檢查規則：
  - `> 300` 行：警告，PR 必須補充理由與拆分計畫。
  - `> 500` 行：視為阻擋條件（CI 失敗）。
  - 對於歷史上已超過 500 行的檔案，僅允許「不增加行數」的修改（鼓勵逐步縮減）。
- 強制工具：
  - 腳本：`scripts/check_file_lines.sh`
  - ignore 清單：`.line-limit-ignore`
  - CI workflow：`.github/workflows/file-line-limit.yml`
  - 本機 hook：`.githooks/pre-commit`（建議執行 `git config core.hooksPath .githooks` 啟用）

## 規範來源分層
- `AGENTS.md` 是專案共通規範入口，負責專案目標、工作流、測試、安全與文件要求。
- `CLAUDE.md`、`GEMINI.md` 是不同 Agent 的執行補充，負責工作方式、脈絡與執行邊界。
- `.github/instructions/` 是語言與框架的正式實作規範來源；Agent 在修改對應檔案時必須一併遵守。
- 檔案類型對應如下：
  - `*.go`、`go.mod`、`go.sum`：遵守 `.github/instructions/go.instructions.md`
  - `*.tsx`、`*.jsx`、`*.js`、`*.css`、`*.scss`：遵守 `.github/instructions/reactjs.instructions.md`
  - `*.ts`：同時遵守 `.github/instructions/reactjs.instructions.md` 與 `.github/instructions/typescript-5-es2022.instructions.md`
- 規範衝突時，優先順序為：`AGENTS.md` 專案目標與工作流 -> Agent 專屬文件 -> `.github/instructions/` 細部實作規範。

## Commit 與 Pull Request 規範
- 近期提交慣例為簡短、祈使語氣的繁中主旨（例：`修正...`、`完成...`、`補齊...`）。
- 每個 commit 聚焦單一主題，並將對應測試一併提交。
- PR 請附：變更摘要、影響模組、測試證據（例如 `go test ./...`、`npm run test`）、相關 OpenSpec 連結；前端變更需附截圖。

## 前端 UI/UX 專案目標
- `/studio` 是 datalink 主產品介面，設計與重構都要優先對齊單人操作情境。
- `Studio` 內的四步驟主流程是唯一主線：
  1. `device`：建立 / 編輯資料來源，並分開呈現 `connect` / `probe` 診斷。
  2. `source`：建立或還原 `Source Rule`，在格狀畫布查看 `planned / used / unmanaged / conflict` 狀態與即時值。
  3. `tag`：以 review-first 方式檢查系統自動建立的 Tag / Mapping；手動 create / existing / unbind 屬於例外處理工具，不是主心智模型。
  4. `output`：直接在 Local Modbus register 或 Database schema / column 表面上綁定輸出。
- `TestPage` 是專用工程測試工具；其核心用途與操作行為應保持穩定，僅做 UI 風格一致化，不作產品主流程承載。
- 前端改版與 cleanup 時，應優先盤查並清除舊 datalink legacy 結構、舊 redirect、未引用元件與過時設計，避免新主線旁再殘留第二套產品入口。

## 文件化工作流
- 多步驟 UI / UX、架構整理或大型重構任務，預設採 SKILL `planning-with-files` 工作法。
- 專案根目錄需維持 `task_plan.md`、`findings.md`、`progress.md` 三份文件，分別記錄階段計畫、關鍵發現與執行 / 驗證過程。
- 開始執行前先做 session catchup；每完成一個 phase，要同步更新計畫狀態、測試結果、修改檔案與錯誤紀錄。
- 重要發現、legacy 清單、風險與失敗嘗試不得只留在對話上下文，必須寫入文件以利續作。

## Worktree 多版本實作規範（強制）
- 當需求包含「平行版本比較 / A-B / 2+ 變體」時，Agent 必須使用 git worktree，不得在主工作樹混做。
- 預設 worktree 目錄為 `.worktrees/`；若不存在先建立。
- 建立 worktree 前必須驗證忽略規則；若未忽略，先更新 `.gitignore` 讓 `.worktrees/` 被忽略。
- 分支命名採 `exp/<topic>-v1`、`exp/<topic>-v2`、`exp/<topic>-v3`。
- 路徑命名採 `.worktrees/<topic>-v1`、`.worktrees/<topic>-v2`、`.worktrees/<topic>-v3`。
- 每個版本必須使用不同啟動埠，避免 dev server 與 API 衝突。
- Agent 回報時必須列出：worktree 路徑、分支名稱、啟動命令、對應 URL。
- 比較完成後，保留勝出版；其餘版本使用 `git worktree remove` 清理。

## OpenSpec 與 Agent 工作注意事項
- 依規劃 phase 與步驟執行，不可跳步或僅交付最小可動版本。
- 先看現行規格 `openspec/specs/`，再處理 `openspec/changes/` 中與本次變更相關的 proposal / tasks。
- 規格導向工作需同步更新 `openspec/changes/.../tasks.md` 的實作進度；若實作與 spec 衝突，先回報差異，不要自行改寫需求。
