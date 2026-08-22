## Context

這是測試契約修正，涵蓋 frontend Vite config、Studio V2 HTTP handler、device connection probe 與 database row-group writer 四個邊界。現有測試已能在部分環境通過，但有四個可定位的非產品因素：`frontend/tests/unit/utils/viteConfig.test.ts` 匯入 production config 時會受 ambient `.env` 影響；`internal/api/handlers/studio_v2_workspace_database_handler_test.go` 與 `internal/api/handlers/studio_v2_workspace_database_row_groups_test.go` 將 Windows temp path 直接插入 JSON 字串；`internal/api/handlers/device_handler_extended_test.go` 將 connection failure 固定比對 `connection refused`；`internal/datalink/dbtarget/writer_row_groups_test.go` 將 row insertion 順序固定成 `[21.5, 23.75]`，但 writer 會從 map／非同步 flush 產生結果。

需求是提高測試對 Windows、Linux、不同工作目錄與正常平行度的可重現性，同時不改變 production route、network error、writer ordering 或 database data semantics。

## Goals / Non-Goals

**Goals:**

- 使 Vite 測試能控制 ambient environment，且仍驗證 `vite.config.ts` 的正式 `loadEnv`、proxy target 與 dev port 行為。
- 使 Studio V2 database handler request payload 在 Windows path 含反斜線時仍是合法 JSON，並涵蓋 database handler、row-groups、audit handler、delivery-truth 測試使用的 fixture／request 建構。
- 使 connection failure test 驗證跨平台穩定的失敗狀態、connect status、probe skipped 狀態與可判斷的 error presence，而不把 OS-specific wording 當成產品契約。
- 使 row-group test 驗證所有預期 row／value 成果而非 map insertion order，並以可控、有限的 flush signal 避免開放式 sleep 或不受控 Eventually。
- 以 focused tests 先形成 feedback loop，再執行 handler package、dbtarget repeated test、frontend full suite 與 Go full suite。

**Non-Goals:**

- 不修改 production `frontend/vite.config.ts`、device handler network error 產生邏輯或 dbtarget writer 的 ordering。
- 不拆分九個超過行數上限的測試檔，不修 Workbench full-suite timeout，不處理 golangci v2 bootstrap／CI。
- 不改 route／data semantics，不新增 live DB／PLC 依賴，也不將 OS-specific error text 重新包裝成新的 production API contract。

## Decisions

### Isolate Vite environment loading without changing production config

測試只在 `frontend/tests/unit/utils/viteConfig.test.ts` 控制 `process.env` 與 Vite config 的讀取邊界，採用既有 production `vite.config.ts` 的 callable config contract，不複製或重構 production config。測試固定清理與還原 `PORT`、`VITE_API_PROXY_TARGET`、`VITE_DEV_PORT`，並建立不依賴 repo ambient `.env` 的 execution context；預設 proxy、環境變數優先級與 dev port 仍依 production implementation 驗證。

替代方案是修改 production config 以忽略 `.env`，但那會改變開發／部署語意且超出 bug fix；故不採用。

### Encode handler request payloads with typed values

Studio V2 database test requests 以 Go struct 或 map 搭配 `encoding/json.Marshal` 建構，再以 bytes buffer 傳給 httptest request。database path、point IDs、column names 與 row-group members 都是資料欄位，不再以 raw string interpolation 形成 JSON。共享 fixture helper 仍建立 temp SQLite path，但 escape 責任交給 JSON encoder；audit 與 delivery-truth 測試保持既有 observable assertions，僅採用同一安全 payload／fixture 建構邊界。

替代方案是對 Windows path 手工替換反斜線或改用 Unix-like temp path；前者容易漏掉 JSON escaping，後者不能代表 production Windows path，故不採用。

### Assert portable connection failure semantics

`device_handler_extended_test.go` 保留使用已關閉 localhost listener 觸發 connection failure，驗證 response success envelope、`data.success == false`、`can_activate == false`、`can_collect == false`、`connect.status == failed` 與 `probe.status == skipped`，並只要求 error 欄位非空或包含可穩定識別的連線失敗訊息。不得因測試改動 production network error mapping；若現有 API 已提供 stable classification，優先使用該 classification，否則不比對 Windows／Linux 特有文字。

替代方案是將 production error 統一改寫成固定英文；這會擴大 API 行為變更，且不是本次測試契約需要，故不採用。

### Separate row-group outcome semantics from flush scheduling

`writer_row_groups_test.go` 以明確的 flush tick／signal 或等價 bounded seam 觸發一次可觀察 flush，等待 row count 到達預期上限後讀取資料，將結果轉成與順序無關的集合或排序後比較。測試必須確認兩個 row-group 都輸出且 values 為 `21.5` 與 `23.75`，但不得要求資料庫 `id` 順序；不得以無限等待、固定 sleep 或提高 GOMAXPROCS 來掩蓋競態。

替代方案是修改 production writer 對 map keys 排序；既有規格沒有 ordering promise，會不必要地新增產品語意，故不採用。

### Use a staged regression sequence with normal parallelism

先分別執行 Vite focused test 兩次，再執行 handler package 的 targeted tests 與 full package，接著以正常 GOMAXPROCS 執行 dbtarget row-group test `-count=20`。上述 loop 通過後才執行 frontend full suite 與 `go test ./...`，以區分本 change 修復與既有全域 baseline failures。

## Implementation Contract

- **Behavior:** 測試不再因 ambient `.env`、Windows temp path JSON escaping、作業系統 connection error wording 或 row-group 非保證順序而非預期失敗；被測 production 行為與 response schema 不變。
- **Interface / data shape:** handler request body 是由 `encoding/json.Marshal` 產生的合法 JSON bytes；connection failure response 仍包含既有 `data.success`、`can_activate`、`can_collect`、`connect.status`、`probe.status` 與 error 欄位；row-group outcome 是兩筆值的無序集合 `{21.5, 23.75}`。
- **Failure modes:** ambient env 不能覆寫測試設定；Windows path 的 `\` 必須由 encoder escape；connection failure 只在 stable state／classification 不符時失敗，不因 OS-specific wording 失敗；flush 若在 bounded timeout 內沒有兩筆 row，測試明確失敗並保留診斷，而非無限等待。
- **Acceptance criteria:** Vite focused test 連續兩次通過；四個 Studio V2 handler test 檔的 focused／handler package tests 通過；device connection failure test 在 Windows／Linux 可用的 assertion 下通過；`go test ./internal/datalink/dbtarget -run RowGroups -count=20` 在正常 GOMAXPROCS 下通過；再執行 frontend full suite 與 `go test ./...` 並記錄非本 change baseline failures。
- **Scope boundaries:** 只修改列出的七個 test files 與本 change artifacts；不修改 production source、既有 route/data/network/writer semantics、其他 line-limit／Workbench／lint bootstrap 工作或 live external systems。

## Risks / Trade-offs

- [Risk] 測試若只檢查 error 非空，可能漏掉 handler 將不同錯誤誤判成 connection failure。→ 優先驗證既有 stable status／classification 欄位，只有 wording 才放寬；若沒有 classification，保留 connect failed 與 probe skipped 的狀態契約。
- [Risk] 使用 bounded flush signal 仍可能受資料庫初始化時間影響。→ 等待明確 row-count predicate，設定有限 timeout，並在 timeout 時回報 query error／觀測 count。
- [Risk] ambient `.env` 的隔離若改變 Vite config 執行上下文，測試可能不再代表 production。→ 不複製 config；透過同一 callable config 與既有 env key 驗證 production semantics。
