## Why

目前四組測試契約在 Windows 與平行執行環境下有可重現的不穩定來源：Vite 測試會讀取工作目錄中的 ambient `.env`，Studio V2 database 測試以字串插值組 JSON，連線失敗測試比對作業系統特有錯誤文字，而 row-group writer 測試把未承諾的 map／flush 順序當成結果順序。這些問題會讓測試在不同平台、暫存目錄或執行順序下失敗，卻不能反映產品行為是否正確。

## What Changes

- 讓 frontend Vite config 測試隔離 ambient `.env`，並以現有 production `vite.config.ts` 的 `loadEnv`／proxy／port 語意作為驗證基準。
- 將 Studio V2 database handler 相關測試的 JSON payload 改由 `encoding/json` 建構，涵蓋 workspace database handler、row-groups、audit handler 與 delivery-truth 測試，確保 Windows temp path 的反斜線被正確 escape。
- 將 `device_handler_extended_test.go` 的 connection-failure assertion 改為跨平台可驗證的錯誤契約；除非測試證據要求，維持 production network error 行為不變。
- 讓 `writer_row_groups_test.go` 驗證不依賴 map insertion order 的 row-group 結果，並使用 deterministic、bounded 的 flush seam；不修改 production writer ordering，因現有規格沒有承諾排序。
- 建立可重跑的 focused regression loop：frontend focused test 連跑兩次、handler package focused/full、dbtarget row-group test 以正常 GOMAXPROCS 執行 20 次，再執行 frontend full suite 與 `go test ./...`。

## Non-Goals

- 不處理九個超過行數上限的測試檔拆分、Workbench full-suite timeout、golangci v2 bootstrap／CI，或其他與本 change 無關的 baseline failures。
- 不改 production route、data、network error 或 writer ordering semantics，也不做 live DB／PLC 驗證。

## Capabilities

### New Capabilities

- `cross-platform-test-contracts`: 定義跨平台測試如何隔離環境、編碼 JSON、驗證平台中立錯誤契約，以及驗證無序／非同步 row-group 結果。

### Modified Capabilities

（無；本 change 修正測試契約與測試隔離，不改產品規格需求。）

## Impact

- Affected specs: cross-platform-test-contracts
- Affected code:
  - Modified: frontend/tests/unit/utils/viteConfig.test.ts
  - Modified: internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - Modified: internal/api/handlers/studio_v2_workspace_database_row_groups_test.go
  - Modified: internal/api/handlers/studio_v2_workspace_audit_handler_test.go
  - Modified: internal/api/handlers/studio_v2_workspace_database_delivery_truth_test.go
  - Modified: internal/api/handlers/device_handler_extended_test.go
  - Modified: internal/datalink/dbtarget/writer_row_groups_test.go
  - New: openspec/specs/cross-platform-test-contracts/spec.md
  - Removed: （無）
