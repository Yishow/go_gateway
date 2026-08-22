## 1. 前端環境隔離契約

- [x] 1.1 [P] Isolate Vite environment loading without changing production config，並完成 Vite tests isolate ambient environment while preserving production semantics 的 RED-GREEN-REFACTOR：先加入可在衝突 ambient `.env` 下重現錯誤的 regression case，再讓 `frontend/tests/unit/utils/viteConfig.test.ts` 控制 `PORT`、`VITE_API_PROXY_TARGET`、`VITE_DEV_PORT` 的快照與清理；以 `cd frontend && npm test -- --run tests/unit/utils/viteConfig.test.ts` 連續執行兩次，確認 proxy fallback、explicit target precedence 與 dev port 仍通過。

## 2. Studio V2 JSON request 契約

- [x] 2.1 [P] Encode handler request payloads with typed values，讓 Studio V2 database requests use JSON encoding for platform paths：在 `internal/api/handlers/studio_v2_workspace_database_handler_test.go`、`internal/api/handlers/studio_v2_workspace_database_row_groups_test.go`、`internal/api/handlers/studio_v2_workspace_audit_handler_test.go`、`internal/api/handlers/studio_v2_workspace_database_delivery_truth_test.go` 的 request／fixture 邊界使用 `encoding/json.Marshal` 或共享 typed encoder，先以 Windows-style temp path regression case 驗證 RED，再完成 GREEN；以 `go test ./internal/api/handlers -run 'TestStudioV2Workspace(Database|Audit)'` 及 handler package full test 驗證合法 JSON、path round-trip、row-group members 與既有 audit/delivery response。

## 3. Device connection failure 契約

- [x] 3.1 [P] Assert portable connection failure semantics，完成 Device connection failure assertions are platform-independent：保留已關閉 localhost listener 與 production response，將 `internal/api/handlers/device_handler_extended_test.go` 從 literal `connection refused` 改為 stable state／classification assertion，驗證 `data.success=false`、`can_activate=false`、`can_collect=false`、`connect.status=failed`、`probe.status=skipped` 與 error presence；以該 focused test 和 `go test ./internal/api/handlers` 驗證 Windows/Linux 可接受的錯誤差異。

## 4. Row-group outcome 契約

- [x] 4.1 [P] Separate row-group outcome semantics from flush scheduling，完成 Row-group tests verify unordered outcomes with bounded flushing：在 `internal/datalink/dbtarget/writer_row_groups_test.go` 以 deterministic、bounded flush tick／signal 觸發 flush，將結果視為無序集合並確認兩筆值 `21.5`、`23.75`，不得改 production writer ordering；以正常 GOMAXPROCS 執行 `go test ./internal/datalink/dbtarget -run RowGroups -count=20`，確認無 order-dependent 或 unbounded-wait failure。

## 5. 回歸與交付檢查

- [x] 5.1 Use a staged regression sequence with normal parallelism，完成 Regression loops prove cross-platform test stability：依序執行 Vite focused test 兩次、handler focused/full、dbtarget row-group `-count=20`、frontend full suite 與 `go test ./...`，記錄每個命令的通過／失敗 scope，並將非四項契約的既有 baseline failure 明確分離。
- [x] 5.2 完成 change-only diff review：確認只修改七個核准 test files 與 Spectra artifacts，production `frontend/vite.config.ts`、device handler、dbtarget writer ordering、route/data semantics 均未改動；以 `git diff --check`、`spectra analyze stabilize-cross-platform-test-contracts --json`、`spectra validate stabilize-cross-platform-test-contracts` 與 staged baseline hash 比對驗證。
