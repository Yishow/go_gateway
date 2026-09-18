# C1 實作與驗證紀錄

## 結論與範圍

2026-09-15 已完成 C1 的 1.1–1.8：**8/8 待辦、5 個 requirement、17 個 scenario 的本機實作與驗證完成。** 未發現本案範圍內尚待修正的缺陷；沒有 Example 或測試排除項目。真正 managed schema apply／recording test write 仍未接線，按規格回 501；本紀錄不是現場驗收或正式發布紀錄。

- HEAD／本機 main 均為 `569da9f98e52b2635638d85250a4af16601d4364`，未 fetch 遠端 main。提案原始證據基準仍為 `1a0311c8e8db9c62fe4388f0701ba38afe552ff7`。
- C1 保持 47 個產品、測試與 Swagger 路徑（含 11 個原 C1 新檔），並更新 proposal、design、tasks、本驗證紀錄。
- 使用者續作授權下，另外完成 **343 個 Go 路徑的基線維護，含 36 個新檔**。主要為格式、既有 API 命名相容性、context、數值界線與錯誤處理；與 C1 重疊的 3 個 Go 路徑以授權時內容為維護差異基準，沒有將維護範圍併稱為 C1 功能。
- 另外兩案與 PUBLICATION 的 12 份既有修改維持原 SHA-256；C1 delta specs、前端與 Swagger 在這輪 Go 維護期間均未改動。
- 未 stage、commit、merge、push、archive 或部署。
- 原 C1 checkpoint：`/private/tmp/go-gateway-c1-authorized-20260915/`。本輪 checkpoint：`/private/tmp/go-gateway-lint-remediation-20260915/`，包含維護專用差異、修改前備份、檔案雜湊、最終測試結果與獨立回復驗證。暫存目錄不是永久保存服務。
- 本機環境：Go 1.27.1 darwin/arm64、golangci-lint 2.13.2、Node 24.15.0、npm 12.0.2；詳細識別見 checkpoint 的 `environment.json` 與 `gate-config.sha256.json`。

## 實作結果


1. 真正 Gin router 的 recording schema-apply／test-write 對合法請求回 501，使用各自固定 code、`success:false`、`retryable:false`、`wait_for_supported_operation` 與完整 request_id；格式錯誤保留 400。無 fabricated record_id、成功 payload 或目標資料庫副作用。
2. recording capability mask 將兩個未接線操作設為 false；保留 adapter 支援資訊、schema preview 與一般 custom-table 建表功能。
3. 前端嚴格解析成功 envelope、布林值、結果列舉與必要證據，限制陣列及 DDL 大小。未知／缺欄 200、斷線及 5xx 不會變成成功或證明零副作用；安全 ID 保留，原始診斷不進 DOM。兩個 recording mutation 與 activation 均禁止自動重試。
4. Step 4 依後端確認呈現啟用結果，分開顯示設定已存、部分／全部失敗、未確認及資料交付狀態。重複點擊、reset、remount、跨 workspace 的舊回覆均受保護；recovery 只重新讀取現有 workspace/runtime 資料，不重送 activation、不建立 ledger，也不把舊 running 視為本次完成。
5. 保留 workspace/settings revision、readiness token、Modbus Share gate 與 runtime 導航。確定單一設備時帶該 ID；無法唯一判定則前往 `/studio/runtime`，callback 僅執行一次。
6. 更新繁中／英文安全文案與 Swagger 400/501 合約，不發布假的 200 成功。Go activation wire contract 未新增 outcome 欄位；前端 outcome 表示客戶端對回覆的確定程度。

## 最終版本驗證

以下均使用 `rtk proxy`；前端命令在 `frontend/` 執行。除另述者外，結果為 passed-current，來源為本輪 checkpoint 中的原始輸出。Go 維護前的 1,110 項 lint 已清為 0，沒有停用 linter 或改動檢查設定。

| 結果 | 命令／檢查 | 證據與範圍 |
| --- | --- | --- |
| PASS | `go test -json ./... -count=1` | 44 個 package、1,379 項測試通過（含子測試），5 項既有 skip；另 6 個 package 無測試；`final-go-test.jsonl` |
| PASS | `go vet ./...` | `final-go-vet.log` |
| PASS | `golangci-lint run --allow-parallel-runners ./...` | 0 issues；`final-go-lint.json`、`final-go-lint.log` |
| PASS | `go test -race ./internal/datalink/runtime ./internal/datalink/delivery ./internal/datalink/collector ./internal/datalink/sourcerule -count=1` | 四個 package 通過；`final-runtime-race.log` |
| PASS | `npm test -- --run --reporter=json --outputFile=<checkpoint>/final-frontend-test.json` | 151 個檔案、880 項測試全數通過，無跳過 |
| PASS | `npm run lint`、`npm run build` | `final-frontend-lint.log`、`final-frontend-build.log`；既有 caniuse-lite 更新提醒未擴成依賴修改 |
| PASS | `bash scripts/check_file_lines.sh` | 所有歷史 >500 行來源檔維持不增行；>300 行警告清單見 `final-file-lines.log` |
| PASS | 本機 TCP／UDP 採集寫入 smoke | `go run ./cmd/loadtest_modbus -tcp 1 -udp 1 -duration 2s -interval-ms 100 -flush-ms 100 -batch-size 4 -value-base 65534 -db <checkpoint>/smoke.sqlite -report <checkpoint>/modbus-smoke.json`；38 筆 good rows，0 錯誤／數值差異 |
| PASS | Swagger 獨立再生與合約測試 | `go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/test_ui/main.go -o <checkpoint>/swagger-check/swagger`；三份產物 byte-identical，見 `final-swagger-identity.json` |
| PASS | `spectra validate fix-studio-v2-database-result-truthfulness --json` | 無 errors／warnings；`c1-validate.json` |
| PASS | `spectra analyze fix-studio-v2-database-result-truthfulness --json` | 0 Critical、0 Warning；保留 13 項建議增加 Example 的 Suggestion，17 個 scenario 已有具體測試資料 |
| PASS | `git diff --check`、雜湊、文件連結與 checkpoint 反向檢查 | 結案文件更新未改動已驗證產品內容；精確維護差異排除原 C1 與另外兩案 WIP |

### 明確跳過與未執行

- Go 自動跳過：`TestDeviceHandler_Activate`、`TestDeviceHandler_ActivateThenDisable` 缺 ConnectionManager mock；`TestPollingGroupHandler_Create_DuplicateName` 的既有 duplicate validation 尚未實作。
- Go 自動跳過：`TestRealDevice_ModbusAndMCProtocol_TelemetryPipeline` 無法連上既有設備；`TestPostgresPartitionMigration` 未設定 POSTGRES_DSN。
- NOT RUN：瀏覽器 E2E、視覺／FHD、正式部署及實際 PLC／PostgreSQL 端到端驗收。拋棄式 SQLite target witness 與本機 TCP／UDP 模擬不能替代現場驗收。
- 交付維護已阻止 receipt 儲存失敗後的自動重送；若連 blocked 狀態也無法持久化，該目的地 worker 會停止並要求核對後再啟動，跨 process 重啟後仍無去重保證。本次未新增重啟 ledger，也未接線 production delivery worker。

### 一般建表與 Local Modbus-only 回歸

- PASS：`TestStudioV2WorkspaceDatabaseHandler_GenerateSchemaCreatesTable`、`TestDatabaseTargetHandler_GenerateSchema_DryRunAndExecute`、`TestConnectorService_GenerateSchema_DryRunAndExecute`，以及 disabled／missing-tag mappings、DB error 與 last schema outcome 測試。
- PASS：[step4-share-activation.test.tsx](../../../../frontend/tests/unit/workbench-v2/step4-share-activation.test.tsx) 的 fixture 使用 `getConfig=null`、空 database targets，8 項測試涵蓋只有 Share 輸出時的 reconcile → activation，以及 Share 關閉、stale metadata、canonical plan 缺失和 reconcile 失敗；readiness／revision payload 斷言保留。
- PASS：`TestSourceRuleHTTP_ShareConfigAndTagApplyProduceReadyLocalModbusCandidate`、`TestSourceRuleHandler_ApplyLocalModbusOutputs_ReturnsPerCandidateResultsForPartialSuccess`、`TestSourceRuleHandler_ApplyDatabaseOutputs_PreservesLocalModbusSnapshot`，以及 `workbench-local-modbus-review-surface.test.tsx` 的 4 項互動測試。

## Requirement 對照摘要

| Requirement | 實作位置 | 有測試的 scenario | Example 值 | 執行證據 |
| --- | --- | --- | --- | --- |
| Unimplemented recording operations fail closed | recording handler 的 SchemaApply／TestWrite | 3/3；router contract、target witness、400 validation | 無 Example | passed-current；完整 Go test |
| Capability flags describe connected implementations | recordingplan capability mask、RecordingPlanSetupSection | 2/2；capability contract、UI gates、一般 GenerateSchema 回歸 | 無 Example | passed-current；完整 Go／前端 test |
| Safe errors and uncertain outcomes remain distinguishable | typedErrors、recordingPlanJson、services／mutation hooks | 3/3；真實語系整合、retry policy、unknown／transport results | 無 Example | passed-current；完整前端 test |
| Commit sequence and animation | useStep4Activation、Step4Database、page recovery | 3/3；evidence boundaries、BackendBackedActivationResults | 無 Example | passed-current；完整前端 test |
| Commit completion card | CommitSuccessCard、neutral summary、runtime resolver | 6/6；evidence boundaries、commit、resolver、partial results | 無 Example | passed-current；完整前端 test |

## Scenario-to-test 對照


每列皆為 covered-by-test／passed-current；來源為本次 `final-go-test.jsonl` 與 `final-frontend-test.json`。下列簡名指向同一 C1 的實作與測試，不代表現場測試。

| Spec scenario | 實作位置 | 測試證據 |
| --- | --- | --- |
| Nonempty schema token reaches an unimplemented operation | recording handler `SchemaApply` | `TestStudioV2WorkspaceRecordingPlans_SchemaApply_UnimplementedRejectsAnyTokenWithoutRepoMutation`；target witness |
| Well-formed test request reaches an unimplemented operation | recording handler `TestWrite` | `TestStudioV2WorkspaceRecordingPlans_TestWrite_UnimplementedRejectsSeededAndMissingPlansWithoutRepoMutation`；target witness |
| Invalid input remains invalid | 既有 DTO／ShouldBindJSON | `TestStudioV2WorkspaceRecordingPlans_MutationValidationRemains400` 的 10 個子情境 |
| Connection succeeds while test writes are unimplemented | recordingplan `GetConnectorCapability`、RecordingPlanSetupSection | `TestRecordingPlanService_Capabilities`、`TestStudioV2WorkspaceRecordingPlans_CapabilitiesMaskUnconnectedMutations`；UI capability gates |
| Existing custom-table operation remains implemented | recording capability mask；既有一般 GenerateSchema | preview API／UI 測試及上列三個 GenerateSchema 回歸 |
| Operation is not implemented | recording service、typedErrors、mutation hooks、語系 | recording-errors integration 的兩種語系 × 兩個 501；retry policy 的 schema／test write 501 |
| Connection is lost after dispatch | recording／activation services 與 hooks | retry policy 的兩個 transport-loss；RecordingPlanSetupSection 的 transport loss；activation-response-truth 的 response loss／opid |
| Malformed or unfamiliar result | recordingPlanJson、safeJson、service parsers | service 的 unknown／incomplete nominal 200；UI unknown results；integration 的 unfamiliar 200／operation ID |
| Presentation timer completes without a response | useStep4Activation、Step4Database | step4-evidence-boundaries：`does not infer success when all presentation timers elapse before a response` |
| Acknowledged partial activation | CommitSuccessCard、CommitProgress | BackendBackedActivationResults：`keeps a partial activation card amber and hands off only confirmed device ids once` |
| Remount during activation | Page recovery、useStep4Activation | step4-evidence-boundaries：`reads server state on remount without replaying activation or accepting the old acknowledgement`；BackendBackedActivationResults recovery／workspace isolation |
| Configuration saved without activation | Step4Database、CommitSummary、ActivationNeutralSummary | step4-evidence-boundaries：`labels saved configuration without claiming activation or database delivery` |
| Activation succeeds without a database receipt | CommitSuccessCard、語系 | step4-evidence-boundaries：`keeps delivery unconfirmed after an acknowledgement and does not rewrite saved device status` |
| Empty results without delivery evidence | Step4Database、ActivationNeutralSummary | BackendBackedActivationResults 的 empty response；step4-evidence-boundaries 的 earlier active／empty new attempt |
| Previously running device and unconfirmed database change | recovery／neutral summary、runtime handoff | step4-evidence-boundaries：`keeps an earlier active device navigable without confirming an empty new attempt` |
| Handoff includes resolved device id | WorkbenchV2Shell、resolveRuntimeDashboardDevice | BackendBackedActivationResults 的 partial activation handoff once；resolver 的 confirmed backend device |
| Handoff falls back when device cannot be resolved | WorkbenchV2Shell、resolveRuntimeDashboardDevice | resolver 的 ambiguous／local confirmations；step4-commit 的 runtime callback／generic fallback |

測試來源索引：

- 後端：[router contract](../../../../internal/api/router_studio_v2_workspace_recording_plans_contract_test.go)、[target witness](../../../../internal/api/router_studio_v2_workspace_recording_plans_target_witness_test.go)、[capability contract](../../../../internal/api/router_studio_v2_workspace_recording_plans_capability_test.go)、[recording service](../../../../internal/datalink/recordingplan/service_test.go)、[Swagger](../../../../internal/api/router_studio_v2_recording_swagger_test.go)。
- 前端：[recording service](../../../../frontend/tests/unit/services/studioV2WorkspaceRecordingPlans.test.ts)、[retry policy](../../../../frontend/tests/unit/hooks/useStudioV2WorkspaceRecordingPlans.retry.test.tsx)、[recording UI](../../../../frontend/tests/unit/workbench-v2/RecordingPlanSetupSection.test.tsx)、[真實語系整合](../../../../frontend/tests/integration/workbench-v2/recording-errors.integration.test.tsx)。
- 啟用：[BackendBackedActivationResults](../../../../frontend/tests/unit/workbench-v2/BackendBackedActivationResults.test.tsx)、[response truth](../../../../frontend/tests/unit/workbench-v2/activation-response-truth.test.tsx)、[evidence boundaries](../../../../frontend/tests/unit/workbench-v2/step4-evidence-boundaries.test.tsx)、[first activation](../../../../frontend/tests/unit/workbench-v2/step4-first-activation.test.tsx)、[commit](../../../../frontend/tests/unit/workbench-v2/step4-commit.test.tsx)、[resolver](../../../../frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts)。

## 最終 review／verify

主代理核對本案 proposal、design、specs、tasks 與測試斷言。沿用原 checkpoint 的 47 檔完整審查紀錄，再以雜湊確認其中 44 檔未變，並完整回讀另外 3 個 Go 路徑的維護差異及新增依賴。已依正確性、效率、簡化／重用與慣例四個面向回讀本輪實作差異；Luna 回報僅作為輔助證據。

本輪修正吞掉 JSON／SQL／receipt 錯誤、整數溢位、SourceRule rollback 漏還原 DataFormat、回復失敗仍同步 runtime、候選失敗未清理先前結果，以及 receipt 儲存失敗後可能自動重送等問題。維護沒有更動 C1 的 501、安全錯誤、未知結果、workspace revision／readiness／Share guards 或 runtime 導航契約。前次已修正的 malformed 200 metadata、跨 workspace 舊回覆、重複請求與 recovery 覆蓋問題仍有本輪完整前端測試證據。

| 檢查面向 | 結果 |
| --- | --- |
| 待辦完整性 | 8/8 完成；全庫 lint 阻擋已解除 |
| 規格與測試 | 5/5 requirement、17/17 scenario 有測試並通過；沒有 Example 或排除項目 |
| 設計符合性 | 保留 C1 安全封鎖與既有 activation 資料流，未新增 ledger 或接線真正建表／試寫 |
| 程式品質 | 本次範圍未發現尚待修正缺陷；完整 Go lint 為 0 issues |

Spectra 的 47 個產品路徑皆可檢視、scope 為 resolved，與獨立 C1 manifest 一致。結案文件更新後的 scope snapshot、current 檢查結果與維護清單存於本輪 checkpoint 的 `c1-final-review-scope.json`、`c1-final-scope-check.json`、`maintenance-files.json`。本機驗證完成；封存與提交尚未執行。

### 行數與產物範圍說明


- `WorkbenchV2Shell.tsx` 481、`RecordingPlanSetupSection.tsx` 367、`Step4Database.tsx` 319 行：本案修改現有流程邊界。若後續 workflow change 改寫流程，可把 recording form／操作結果及 shell 的 workflow adapter 分離；C1 不做相鄰 UI 重構。
- `safeJson.ts` 322 行：已把本次 recording 專屬 parser 拆至 `recordingPlanJson.ts`；剩餘為既有共享型別解析，後續新增其他領域 parser 時依其邊界分離。
- `BackendBackedActivationResults.test.tsx` 321、`RecordingPlanSetupSection.test.tsx` 348、`step4-commit.test.tsx` 432 行：各自沿用 cohesive fixture。若新增 lifecycle 或 recording 表單情境，再按情境群組拆檔及抽取同群 fixture。
- 兩份語系 JSON 保持既有 502 行，僅改相關 step4／errors 文案。
- Swagger 除兩個 recording endpoint 與 request DTO，還包含 generator 對既有 `mapping.StepResult` 的 `input_value`／`output_value` 兩個 property 所產生的既有漂移。三份產物均與 pinned generator 的獨立輸出一致，沒有手改 generated output。

## Release note

記錄方案的建表／試寫尚未接上可驗證的資料庫操作，現在明確顯示未開放並回 HTTP 501，舊客戶端不能再依 200 誤判完成。能力旗標同步停用這兩個操作；一般資料庫建表及 Local Modbus 流程保留。Step 4 只依實際回覆顯示啟用結果，部分失敗、斷線、重載及缺少交付證據時持續顯示未確認；runtime 觀察入口仍可使用。

這是刻意的安全相容性中斷。真正 managed schema／test write 須由後續 change 完成服務、scope 與驗證後開放。本紀錄尚不是部署或正式發布紀錄。

## 安全 rollback

- 尚未部署；未對正式環境執行 migration 或目標資料庫變更。本機測試只使用拋棄式 SQLite 與本機模擬器，因此沒有正式資料的回復步驟。
- 需要局部回復時，先保留當時工作樹。原 C1 使用 `go-gateway-c1-authorized-20260915/authorized-change.patch`；本輪維護使用 `go-gateway-lint-remediation-20260915/maintenance-only.patch`，結案文件使用同目錄的 `closeout-documents.patch`。兩批差異有先後相依，須依 manifest 逐 hunk 核對；授權前既有 WIP 以原始備份與雜湊為準，禁止整棵工作樹 reset。維護差異已做反向檢查與隔離副本還原比對，沒有回復實際工作樹。
- 不可直接部署回舊的假成功 handler。若回退 UI 或其他整合，必須同時保留兩個 501 handler、capability=false 及未知結果不成功的封鎖；沒有安全的舊版本時，停用這兩個操作並修補向前。
- 回復後重跑 router 400/501、target witness、capability、未知結果／不重試及一般建表／Local Modbus 回歸，再執行所影響層的必要完整檢查。
