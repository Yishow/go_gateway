# go-gateway

以 Go + React 建置的工業資料採集閘道。透過 Web UI 設定設備連線、採集點位、Point-to-Tag 映射與輸出，將 PLC 資料送往 Local Modbus 或資料庫。

主程式為 `cmd/test_ui`：單一可執行檔提供 HTTP API、SSE 及嵌入式前端。技術版本以 [go.mod](go.mod) 和 [frontend/package.json](frontend/package.json) 為準，目前為 Go 1.25.5、React 19、TypeScript 5、Vite 7。

## 從哪裡開始

- `/studio/v2`：主要設定入口，依設備連線、採集點位、Tag 映射、輸出設定四步進行
- `/studio/runtime`：設定後查看 runtime 狀態、問題與最近交付資訊
- `/test`：工程測試工具
- `/gateway/*`：實驗性介面，不是主要產品流程

`/` 與未知前端路由會依現有通用規則回到 `/studio/v2`；沒有獨立 `/studio` 工作台。改善方向是在現有 V2 內簡化流程，保留 Go／React，不建立另一套 V3。

## 目前能力與限制

下列是依目前 source 核對的能力邊界，不是本次執行的驗收報告：

- **已接 runtime**：collector／scheduler、內部 SQLite storage、dbtarget writer 與 Local Modbus target fan-out，入口見 [service_wiring.go](cmd/test_ui/service_wiring.go) 和 [target_writer.go](cmd/test_ui/target_writer.go)
- **資料庫寫入**：SQLite／PostgreSQL 的實際 SQL 寫入與 grouped time-bucket buffer 已有實作。buffer 接收成功不等於外部 DB 已提交；目前 grouped buffer 位於記憶體，不能視為耐重啟佇列
- **已落地的安全設定**：已存 connector identity／revision、真實 schema metadata、preview／confirm／apply 及持久 operation 記錄；設定與啟用沿用 workspace/settings revisions 和 readiness token，不應略過
- **仍待串接**：recording-plan 設定與 measurement、aggregation、delivery 元件的存在，不代表 production entrypoint 已連成完整的 durable recording 路徑
- **明確未開放**：recording-plan test-write handler 目前回傳 501 `RECORDING_TEST_WRITE_NOT_IMPLEMENTED`，不能宣稱已有試寫、讀回與清理驗證
- **四步流程待收斂**：Step 4 仍同時承載 recording plans、measurement membership、row groups、target mappings 與 Modbus Share。規劃中的單一 write-group authority 尚未實作
- **協議可用性需逐層核對**：後端已有 Modbus、FATEK、MC 等協議元件；選單列出某協議或依賴中有其套件，不代表 V2 位址解析、採集與輸出皆已貫通

基本需求是「所選 Tag 依群組寫到資料庫」。報表、保留政策、聚合、區間電量等衍生運算不是基本寫入的前置條件。設備 running、設定已儲存、buffered、SQL committed、readback verified 與 cleanup completed 應分開理解。

後續實作見 [六份相依提案與驗收計畫](docs/plans/studio-v2-write-groups/README.md)。[原資料庫流程工作](openspec/changes/fix-studio-v2-database-workflow/handoff.md) 的完成證據保留、未完需求已明確移交；文件通過檢查不表示產品已完成。

## 本機建置與啟動

需要 Go 1.25.5 或符合 go.mod 的工具鏈、Node.js／npm，以及 Linux/macOS 的 make 與 shell，或 Windows PowerShell。前端依賴的 Node engine 範圍以 lockfile 為準；建議使用 Node 22.14 以上的 22.x 或相容更新版。

### Linux／macOS

```sh
make build
HOST=127.0.0.1 PORT=8080 ./bin/test-ui.exe
```

開啟 `http://127.0.0.1:8080/studio/v2`。Makefile 依序安裝前端依賴、build、複製 static、build Go；`.exe` 是固定檔名，不代表在 Linux/macOS 產生 Windows binary。

### Windows

```powershell
Push-Location frontend
npm ci
Pop-Location
powershell -File scripts/build.ps1
$env:HOST = "127.0.0.1"
$env:PORT = "8080"
.\bin\test-ui.exe
```

請確認 npm ci 成功後再建置，建置成功後再執行。開啟 `http://127.0.0.1:8080/studio/v2`；腳本細節見 [scripts/build.ps1](scripts/build.ps1)。

### 前端開發

```sh
cd frontend
npm ci
npm run dev
```

API backend 需另啟動；Vite 預設為 5173，API proxy 預設指向 `http://127.0.0.1:8080`。可用 `VITE_DEV_PORT`、`VITE_API_PROXY_TARGET` 或 `PORT` 調整，詳見 [vite.config.ts](frontend/vite.config.ts)。

fresh clone 只保留 embed placeholder。單獨 go build 能編譯，不代表已包含可用前端；請用完整建置流程同步 `cmd/test_ui/static`。

## 設定與資料

主要 server 從環境變數與 `.env` 讀取設定，見 [internal/config/config.go](internal/config/config.go)：

- `HOST`／`PORT`：監聽位置，預設 host 為空、port 為 8080；本機試用建議明確使用 `127.0.0.1`
- `GATEWAY_DB_PATH`：內部 SQLite 檔案位置；相容 aliases 為 `DB_PATH`、`SQLITE_PATH`，預設使用工作目錄的 `datalink.db`
- `VITE_API_PROXY_TARGET`：前端開發時的 backend 位置，與設備／外部 DB endpoint 不同

內部 SQLite 是設定與 runtime 儲存，不等於使用者選取的外部目標資料庫。對外 probe、建表與寫入均由 backend 主機發起。請先用 simulator 和可丟棄資料庫驗證，不把正式 PLC 或正式資料表當作測試環境。

不要將 `.env`、真實 DSN、帳密或設備憑證提交到版本控制。備份、還原與 migration 前先停下相關寫入並核對作用範圍；程式版本 rollback 不會自動撤銷外部建表或已寫入資料。

## 驗證指令

```sh
go test ./...
go vet ./...
golangci-lint run ./...
cd frontend
npm run lint
npm test -- --run
npm run build
npm run test:e2e
```

`golangci-lint` 與 Playwright browser dependencies 需在執行環境可用。這裡列出命令，不代表目前 checkout 已跑完全部測試。

其他工具見 [Makefile](Makefile)：`make check-lines`、`make cross-platform-loop`、Modbus gate／soak、points migration precheck／apply 與 Swagger 產生。跑 migration 或設備測試前先讀腳本，確認 target、權限、備份及可丟棄環境。

## 專案導覽

- [cmd/test_ui](cmd/test_ui)：實際 server entrypoint、wiring 與 embed
- [internal/api](internal/api)：HTTP／SSE 與 handlers
- [internal/datalink](internal/datalink)：採集、映射、runtime、輸出與資料模型
- [frontend](frontend)：介面、services、types、測試
- [docs](docs)：技術與維運文件；[Studio inventory](docs/technical/studio-surface-inventory/START_HERE.md) 提供歷史與 surface 脈絡，日期較舊的結論需對照 source
- [openspec](openspec)：規格與提案；規格存在或提案驗證通過，不代表產品實作完成

開發與 agent 共通規範見 [AGENTS.md](AGENTS.md)，Claude 補充見 [CLAUDE.md](CLAUDE.md)。README 專注使用與能力邊界，不另存一套開發規則。
