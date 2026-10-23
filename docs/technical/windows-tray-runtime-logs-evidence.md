# Windows Tray／runtime logs 實作驗證記錄

## 2026-10-23 Mac 未提交工作收尾

基準 `f21e2a6b4997726011e4e85eba527087b6c9c8c0`，Mac Studio darwin/arm64、Go 1.27.3（repo target 1.25.5）。本節才是這輪結果；下方歷史段落的「本輪」、PASS/FAIL/NOT RUN 都指其當日，不借用為本輪執行證據。

- 原工作樹12檔先做內容／index快照；3份新Go測試、jitter與zh-TW鍵盤測試完整保留並重驗，未整包stage未知檔或產物。
- 本輪對 runtime logs production source 沒有修改；補3個真embedded瀏覽器契約，全部成功回覆來自隔離後端。transport斷線使用原生fetch AbortSignal故障注入＋offline/online，非假成功snapshot，也不是setOffline本身保證切斷既有SSE。
- `go test -p 1 ./...`、`go vet ./...`、`golangci-lint run ./...` PASS（lint 0 issues）；7套件 race PASS：apphost、config、cmd/test_ui、diagnostics、api、api/handlers、device。
- 前端 lint、**211檔／1344項 Vitest**、含tsc的build PASS；logs focused **27項** PASS；full Playwright **18項** PASS，embedded focused **13項** PASS。
- 實際 `::1` logs listener PASS、無Skip；首次sandbox bind限制的Skip不計成功。build-command contracts **5/5** PASS，僅dry-run／腳本契約，未做Windows原生建置。
- 原稿jitter與新source/search測試經production mutation抓到退化並立即還原；resolver alias順序與config precedence經Go overlay mutation抓到，不改production檔。
- 舊48-scenario表的R16正反source/search、R23 clear前後server IDs、R25受控transport中斷後保列／cursor／重連，現新增直接browser證據。R26 hostile literal與雙容量仍是單元證據，不能冒稱真backend注入。未補的Windows原生／完整artifact情境仍保持缺口。
- 原20項任務維持10完成／10未完成（3.3已由本輪證據核對）；新增7.1–7.3只處理Mac證據收尾，已完成，合計13/23。**整個change不archive**，不因Mac PASS勾Windows驗收。
- Spectra validate PASS；analyze Critical=0、Warning=20、Suggestion=30。20 warnings為原requirement／design標題未逐字出現在task description的traceability，未改寫completed歷史來湊全清。

完整檔案盤點、命令/log/source指紋、review及Windows剩餘項目見 [本輪交付](../releases/pending-work-20261023.md) 與 [manifest](../releases/evidence/pending-work-20261023.json)。這輪沒有連現場PLC、正式DB、部署、Windows或Linux原生環境；Mac自動化不代替真人驗收。

## 2026-10-10 歷史 checkpoint

> 本節記錄「本輪」現況；「以下為先前交付歷史」標題後的 PASS／FAIL／NOT RUN 僅代表當時結果，保留原文不改寫。

### 基準與範圍

- 工作樹基準 HEAD `f752b1eae96b070dc0cf0940a9618631de72b377`；darwin/arm64；安裝 Go 1.27.2，`go.mod` target 1.25.5。**無新 commit**。
- 本輪變更僅：測試、proposal 說明釐清、task checkbox、本證據文件。**production source 未變更**。
- Proposal：`embedded-frontend-delivery` 說明自 Modified Capabilities 移到 Impact。`spectra analyze` Critical 1 → 0；其餘 20 Warning、30 Suggestion **未自動修復**，不可稱 analyze 全清。`spectra validate` JSON 有效、無 error／warning；openspec strict validate 通過。

### Tasks

- 10/20 完成。3.3 以 `spectra task done` 新完成，並明確歸屬兩個前端測試檔。
- 1.1 **維持未勾選**：可攜部分已覆蓋（macOS CLI、POSIX 權限行為 PASS），但 Windows 原生新 DB 對話框／ACL 與 Linux CLI 執行 NOT RUN。其餘 pending 任務狀態不變。

### 新增測試

| 檔案 | 行數 | 涵蓋 |
| --- | --- | --- |
| `internal/apphost/data_test.go` | 330 | resolver：mode／root、env／dotenv／default 優先序與別名、canonical path identity、歧義／缺 parent／不可讀 fail closed；以整棵目錄樹比對證明 resolver 不建立、搬移或改寫 storage |
| `internal/config/precedence_test.go` | 65 | 設定優先序 |
| `cmd/test_ui/gateway_data_selection_test.go` | 354 | 實際 `prepareHost`／`runGateway`：desktop 新庫取消（不建立 shell 與任何啟動資源、無 DB／WAL／SHM／journal、exe／cwd 目錄不變、取消成功退出）、歧義與明確路徑失敗不開資源不 fallback、CLI cwd dotenv 語意；另含 resolver 與 harness DSN identity 一致 |
| `frontend/tests/unit/runtime-logs/runtimeLogSession.test.ts` | 116 | jitter 0.8／1.2 邊界、cap、cleanup |
| `frontend/tests/e2e/embedded-frontend-delivery.spec.ts` | 324 | 真實 embedded backend、zh-TW 鍵盤流程：Studio setup → logs/runtime、filter／pause／resume |

限制：非 Windows 的 decline seam 只證明跨平台邏輯，**不證明 Windows 對話框**。

### 驗證（本輪實際執行）

| 命令／驗證 | 結果 |
| --- | --- |
| `env GOCACHE=/private/tmp/go-gateway-apply-cache go test ./internal/apphost ./internal/config ./cmd/test_ui -count=1` | PASS |
| 同三套件 `go test -race -count=1` | PASS；先前預設 cache／toolchain 的 race build 失敗已被此隔離 cache 實跑取代，非未解 blocker |
| `golangci-lint run ./internal/apphost ./internal/config ./cmd/test_ui`（隔離 GOCACHE／GOLANGCI_LINT_CACHE） | PASS，0 issues |
| `env GOCACHE=/private/tmp/go-gateway-apply-cache go test ./...` | PASS（Codex 最終追加檢查；darwin，不含 Windows build tag 檔） |
| `env GOCACHE=/private/tmp/go-gateway-apply-cache go vet ./...` | PASS（同上） |
| `env GOCACHE=/private/tmp/go-gateway-apply-cache GOLANGCI_LINT_CACHE=/private/tmp/go-gateway-lint-cache golangci-lint run ./...` | PASS，0 issues（Codex 最終追加檢查；僅描述本基準，下方歷史段落的 FAIL 為當時基底之記錄，不改寫） |
| `npm run lint`（frontend） | PASS |
| `npm test -- --run` | PASS，203 files／1305 tests |
| `npm run build` | PASS（於 embedded E2E fixture 內執行） |
| `npx playwright test tests/e2e/embedded-frontend-delivery.spec.ts --workers=1` | PASS 10/10；首次 sandbox 瀏覽器啟動被拒屬環境因素，同一授權重跑成功 |
| English logs embedded controls，390／768／1440 | 各 PASS（3/3）；截圖由 Codex 目視 |
| 最終 zh-TW 截圖 | 目視檢查，所檢視狀態未見版面缺陷 |
| `python3 scripts/tests/test_windows_build_contract.py` | PASS，5/5；dry-run／腳本契約，不是 Windows 原生 build |
| `go test ./internal/api -run '^TestRuntimeLogIPv6Socket$' -v -count=1`（隔離 GOCACHE） | PASS，實際 ::1 listener，無 Skip |
| `make check-lines`、`git diff --check` | PASS（最終文件編輯後重跑；300 行警告保留） |

### Mutation 驗證

- jitter：移除 jitter 後 2 個新測試失敗；還原後 focused 27/27 PASS。
- Go：以 overlay 逐項改動 ambiguity、alias 順序、CLI root、cancel 狀態、env 優先序、selected config root，各自產生預期測試失敗，還原後 PASS，未修改 production 檔。
- config-root 首次嘗試是 setup error，已判定為無效證據並改為真正的 assertion 失敗；僅計入行為層失敗。

### Review 範圍

- Spectra review snapshot 採 `touched_tracking`，僅含兩個歸屬的前端測試 diff（base=head `f752b1e…`），limitations 為空。
- Codex 另直接對照初始乾淨基準審查 3 個新 Go 測試與 proposal／task 差異。結論：本增量**未發現已確認缺陷**；這是「目前增量」審查，**不是**完整歷史實作品質、規格或 audit 認證。

### 其他事實

- 歷史上 tasks 提到的 diagnostics／cmd file-sink path 5+1 測試失敗，在本基準已不重現（先前為 macOS canonical fixture 修復）；兩套件 preflight PASS。不改寫 task 歷史，也不宣稱本輪修好 production。
- 行數：330／354 行的 Go 測試與 324 行 E2E 超過 300 行警告，已保留；無 >500 阻擋。若續增長，拆分方向為「共用 fixture helper」與「情境 case」分檔，本輪不主動重構。

### TDD test-presence audit（48 scenario）

僅為「測試是否存在、斷言是否對得上」的對照，**不是** verify 或 security audit，也不改任何 task 狀態。分類只有三種，每列恰好一種：`有測試`（已有對應行為測試；證據層級與限制另列，Windows 測試存在但本機無法跑仍屬此類，另註 NOT RUN）、`測試範圍排除`、`覆蓋缺口`（只有部分 seam，或找不到對應斷言）。本表**沒有** `測試範圍排除` 列：每項皆跨層，或含條件／預設／降級契約，且缺平台不構成排除理由。

方法限制：Claude 初稿依指定檔案讀取斷言；Codex review 再以實際檔案與搜尋補上 file faults、transport、durable recovery 的對照。`覆蓋缺口` 包含完整情境尚未對照或僅有部分測試，不等於測試不存在，更不是已確認的程式缺陷。

執行狀態（與存在與否分開）：已列 Go 測試屬 `go test ./...` PASS 的範圍（darwin／arm64）；`DO`、`DI` 帶 `//go:build windows`，**NOT RUN**，另有數個測試在 Windows 分支 Skip；前端 `FS`／`FP`／`FB`／`FT` 屬完整 Vitest PASS；`E2E` 屬 Playwright 10/10 PASS；`BC` 本輪實跑 5/5 PASS（僅 build-command contract，不是 Windows 原生建置）。PASS 只表示測試通過，不代表涵蓋完整。

檔案鍵（`::` 後為測試函式或標題）：

- `internal/api/`：`AG`＝runtime_log_guard_test.go；`AR`＝runtime_log_routes_test.go；`AQ`＝runtime_log_query_test.go；`AT`＝runtime_log_stream_test.go；`AH`＝runtime_log_http_test.go；`AX`＝runtime_log_transport_test.go
- `internal/diagnostics/`：`GB`＝broker_test.go；`GS`＝subscription_test.go；`GP`＝projection_test.go；`GF`＝file_sink_test.go；`GD`＝file_fault_test.go
- `internal/apphost/`：`AD`＝data_test.go；`AL`＝launch_test.go；`AS`＝shutdown_test.go
- `internal/desktop/`：`DT`＝desktop_test.go；`DO`＝owner_windows_test.go；`DI`＝ipc_windows_test.go
- `cmd/test_ui/`：`CS`＝gateway_startup_test.go；`CD`＝gateway_data_selection_test.go；`CR`＝gateway_readiness_test.go
- `frontend/tests/unit/runtime-logs/`：`FS`＝runtimeLogSession.test.ts；`FP`＝RuntimeLogsPage.test.tsx；`FB`＝runtimeLogBuffer.test.ts；`FT`＝runtimeLogsTransport.test.ts
- `E2E`＝frontend/tests/e2e/embedded-frontend-delivery.spec.ts；`BC`＝scripts/tests/test_windows_build_contract.py
- `internal/datalink/groupdelivery/`：`DS`＝store_test.go；`DF`＝fencing_test.go

`R`＝runtime-log-viewer、`W`＝windows-desktop-lifecycle（兩份 delta spec 路徑見 `openspec/changes/add-windows-tray-runtime-logs/specs/`）。

| # | Scenario | 分類 | 現有測試（已讀斷言） | 缺／限制 |
| --- | --- | --- | --- | --- |
| R01 | Production sources feed the managed view | 覆蓋缺口 | `AH::TestManagedHTTPProjectionPreservesTemplatesAndSuppressesSecrets`（http.access、http.recovered）；`CS::TestProductionManagedWritersAndFileOnlyPolicy`（log／slog／Gin writer 不留 fixture 文字）；`CS::TestRuntimeFailureProjectsDegradedWithoutRawDetails`（runtime.degraded）；`GP::TestUnknownSecretInputIsSuppressedBeforeStorage`（raw.suppressed） | 缺：log／slog／Gin writer 經 host 實際產生固定 `raw.suppressed` 的斷言、HTTP-server ErrorLog 來源、startup fault 穩定 code 對照 |
| R02 | Secret-bearing input never enters a new sink | 覆蓋缺口 | `GP::TestUnknownSecretInputIsSuppressedBeforeStorage`（ring／snapshot）；`AH::TestManagedHTTPProjection…`（snapshot 無 `PRIVATE_`）；`AR::TestRuntimeLogSnapshotFiltersAndSecrets`（snapshot body） | 缺：檔案 sink、SSE frame、operator DOM 的秘密斷言；regex 縱深防禦無獨立測試 |
| R03 | Safe HTTP logs do not inherit full query logging | 有測試 | `AH::TestManagedHTTPProjection…`（query／Authorization／Cookie／錯誤文字／panic 值皆不入 snapshot；有 route template 與 request_id）；`GB::TestHTTPProjectionAllowsOnlyScalarTemplates`（僅允許欄位）；`CS::TestProductionManagedWritersAndFileOnlyPolicy`（console mirror 相容另測） | status／duration 值僅間接驗證 |
| R04 | Log reads cannot amplify themselves | 覆蓋缺口 | `AR::TestProductionRouterWiresManagedLogRoutes`（讀 snapshot 後 ring 仍僅 1 筆）；`AH::TestManagedHTTPProjection…`（`/system/logs` 讀取不產生 access） | 缺：SSE 讀取與 heartbeat 不產生 access 的斷言；sink 失敗回報不遞迴寫回失敗 sink |
| R05 | A burst exceeds capacity | 有測試 | `GB::TestAdmissionHasIndependentCountAndByteCaps`、`TestRingHasIndependentCountAndByteCaps`（count／byte 各自界線，drop／evict >0）；`TestEncodedEventUTF8Limit`（UTF-8 安全截斷＋flag）；`TestConcurrentSequencesAndDistinctInstances` | drop／evict 斷言在 broker 內部計數，上限以常數比對 |
| R06 | One reader stops consuming | 有測試 | `GB::TestFilteredProgressSlowReaderAndCapacity`、`AT::TestRuntimeLogCapacityAndShutdown`、`TestRuntimeLogOverflowGapStaysWithItsOwnStream`；`AX::TestRuntimeLogStalledTransportDeadlineDoesNotLockBroker`（net.Pipe 不讀取：producer 與 snapshot 續行、4–6 秒退出、8 個訂閱 slot 可再次取得） | 真實阻塞 transport 有斷言；Windows native 不屬此 API 情境 |
| R07 | Rotate at the size boundary | 有測試 | `GD::TestRotationIsThreeBoundedOwnedFiles`（連續寫入跨輪替、三檔各 <=5 MiB、總量 <=15 MiB、owned reopen）；`GF::TestAllProtectedAliasesRejectedBeforeModification`、`TestRejectProtectedAndUnownedFileSet`（保護資料與非自有檔案不改寫） | Windows 原生 reparse 與磁碟驗收仍屬 task 2.2 |
| R08 | Diagnostic storage fails after startup | 覆蓋缺口 | `GD::TestDiskFailureRecoveryAndSafeMemoryContinue`（disk error、ring 持續、degraded／loss／recovery interval、檔案無 raw error）、`TestBlockedWriterCannotBlockCaptureAndCloseDeadline`（stall 不擋 capture、Close deadline 誠實）、`TestReadOnlyDirectoryFailsWithoutFallback`；`FS`／`FP`／`FT` 有 file-loss／degraded 顯示 | 缺完整 tray 狀態及 Windows 原生故障驗收；後端 portable 故障注入已有測試，不應重做 |
| R09 | Separate databases do not share a rotation set | 有測試 | `GF::TestDefaultNamespacesOverlapLocksAndOwnedReopen`（不同 DatabaseID 預設路徑不同；active／.1／.2 任一重疊的第二 sink 被拒；關閉後可重開） | 以 sink 錯誤為準，host 層 safe startup code 未斷言 |
| R10 | A configured log path aliases protected data | 覆蓋缺口 | `GF::TestAllProtectedAliasesRejectedBeforeModification`、`TestRejectProtectedAndUnownedFileSet`、`TestDirectorySymlinkAndWindowsLeafSyntaxRejected`、`TestMissingProtectedLeafThroughDirectoryAliasNotCreated`；`GD::TestPathSwapCannotRedirectHeldWritesOrLieAboutSavedPath`（開啟後替換成 DB symlink，不改 DB、不謊稱已保存） | portable link-swap 已覆蓋；該測試 Windows Skip，Windows reparse 競態仍待驗 |
| R11 | Initial diagnostics cannot be saved | 覆蓋缺口 | `DT::TestNativeErrorAdvertisesOnlySavedDiagnostic`（未保存時明說、不顯示路徑） | 缺：host 層 sink 開啟失敗→acquisition 前 native 錯誤、不換 DB、不宣稱 log 頁可用 |
| R12 | Normal embedded same-origin reads work | 覆蓋缺口 | `AG::TestLocalLogGuardConnectionBoundary`、`AR::TestProductionRouterWiresManagedLogRoutes`；`E2E` 真實 embedded IPv4；`AX::TestRuntimeLogIPv6Socket`（真實 ::1 listener＋同源 Origin、200） | 成功 snapshot／SSE 的 no-store／nosniff／same-origin policy 逐項斷言尚未對照；IPv6 測試存在，實跑結果見驗證表 |
| R13 | Spoofed proxy and LAN requests are denied | 有測試 | `AG::TestLocalLogGuardConnectionBoundary`（LAN peer、LAN destination）、`TestLocalLogGuardRejectsProxyHeaders`；`AR::TestRuntimeLogGuardPrecedesStoreAndWildcardCORS`（兩端點固定 403 `logs_local_only`，store 為 nil 仍先拒） | proxy header 案例僅在 guard 單元，未逐端點經 router |
| R14 | A hostile browser origin or Host is denied | 有測試 | `AG::TestLocalLogGuardConnectionBoundary`（rebinding Host、錯 port、foreign／null／重複 Origin、same-site 無 Origin、cross-site、畸形 Fetch Metadata）；`AR::TestRuntimeLogGuardPrecedesStoreAndWildcardCORS`（全域 `*` CORS 下兩端點仍 403、無 ACAO／ACAC、preflight 無 ACAO） | — |
| R15 | Remote-only binding does not expand log exposure | 覆蓋缺口 | `AL::TestVerifiedListenerURLs`（LAN、127.0.0.2、fe80 的 logs 旗標為 false）；`FP::translates denial and scope controls…`（403 時說明僅限本機） | 缺：tray／page 顯示此綁定無法檢視日誌的文字斷言；不新增 listener、不改 HOST 的斷言 |
| R16 | Filter a retained window | 有測試 | `AR::TestRuntimeLogSnapshotFiltersAndSecrets`（source＋level＋q 僅回符合的 1 筆）；`GB::TestQueryValidationAndUnicodeFolding`（空結果為 `"records":[]`）；`FP::distinguishes empty…`、`shows safe failed state rather than empty…`；`E2E`（篩選後 No matching retained records） | 跨級最低嚴重度排序僅以單一 warn 案例間接支持 |
| R17 | Reject an unbounded query | 有測試 | `AQ::TestRuntimeLogQueryStrictBounds`（limit 0／501／-1、q>256、regex、重複參數、壞 cursor 逐項拒絕）；`AR::TestRuntimeLogSnapshotFiltersAndSecrets`（limit=501→400 固定 `logs_invalid_query`）；`GB::TestQueryValidationAndUnicodeFolding` | — |
| R18 | Append between snapshot and subscription | 有測試 | `GS::TestSnapshotToSubscribeConcurrentAppendHasNoGap`（並行 append，水位後依序無缺、progress 到 250）；`AT::TestRuntimeLogReplayLiveAndFilteredProgress`；`FB::deduplicates without coercing…`（穩定 ID 去重） | — |
| R19 | Resume after eviction or process restart | 有測試 | `GB::TestRetentionReplayLimitResetAndFutureCursor`（retention／replay_limit gap、舊 instance→reset、未來 cursor→ErrInvalidCursor）；`FP::distinguishes empty, capture gaps, file-only gaps and reset`；`FS::distinguishes file-only loss and process reset…` | — |
| R20 | A filter matches no new records | 有測試 | `AT::TestRuntimeLogReplayLiveAndFilteredProgress`（被濾掉的事件只推進 heartbeat progress，frame 序列無 gap）；`GB::TestFilteredProgressSlowReaderAndCapacity`（level=error 訂閱 progress 到 300）；`TestReplaySeparateAndProgressRequiresAck`（未 ack 不前進） | — |
| R21 | Replay is larger than the live queue | 有測試 | `GB::TestReplaySeparateAndProgressRequiresAck`（replay 400、live 0）、`TestRetentionReplayLimitResetAndFutureCursor`（500 筆上限＋replay_limit）；`GS::TestReplayAndLiveIndependentByteLimits`（1 MiB 位元組上限＋gap）；`AT::TestRuntimeLogOverflowGapStaysWithItsOwnStream` | producer 不被網路 I/O 阻擋以 live 寫入阻塞間接證明 |
| R22 | Close subscriptions safely | 覆蓋缺口 | `GB::TestFilteredProgressSlowReaderAndCapacity`（並行 Close）、`AT::TestRuntimeLogCapacityAndShutdown`、`FS` stop 清除 timer、`FP` unmount 中止 stream；`AX::TestRuntimeLogStalledTransportDeadlineDoesNotLockBroker`（handler 結束後 slot 全釋放） | stalled transport 釋放已有斷言；已連線 SSE 在 gateway shutdown deadline 內結束的完整整合證據仍待對照 |
| R23 | Pause, clear, and resume | 有測試 | `FS::pauses/clears locally and resumes from processed progress…`；`FP::renders literal hostile text…`（pause→clear→resume，snapshot 1 次、stream 2 次）；`E2E`「/studio/logs connects locally…」（clear 後伺服器 snapshot 仍 ok，resume 回 Connected） | E2E 僅斷言 records 為陣列，未比對 clear 前後伺服器筆數 |
| R24 | Filters change during an outstanding response | 有測試 | `FS::ignores canceled filter A snapshot and old stream callbacks`（abort、無舊 stream、timer 為 0）；`FP::ignores late filter A and cancels filter B…` | — |
| R25 | Reconnect truthfully after a dropped connection | 有測試 | `FS::requires handshake, retains rows on disconnect and owns one retry timer`、`backs off 1/2/4/8/16/30 seconds…`、`jitters retries to 0.8x/1.2x…`、`stops retries on %s 403…`；`FP::translates denial and scope controls…` | 真實瀏覽器斷線重連未在 E2E |
| R26 | Malicious text is not executable | 有測試 | `FP::renders literal hostile text…`（img／script 不成為元素）；`FB::evicts oldest rows at the independent 1000-record limit`、`enforces 2 MiB UTF-8 bytes…` | — |
| W01 | Launch the standalone desktop artifact | 覆蓋缺口 | `BC::test_desktop_flags_and_artifact`（dry-run：`-tags desktop`、`-H=windowsgui`）；`CR::TestReadinessRequiresOwnedEmbeddedResponse`；`E2E`（console 產物載入 setup／logs／runtime，無 404） | 缺：Windows 互動桌面——無 console 與 taskbar 視窗、tray 實際建立、無 Go／Node、開啟頁面不啟動 command shell（task 4.2／5.3） |
| W02 | Preserve console and headless execution | 有測試 | `AL::TestLaunchModesAndExplicitOverride`（`--headless` 在 AUTO_OPEN_BROWSER 為 true／1 仍不開；無旗標 console 保留 opt-in）；`CS::TestInfoCommandsDoNotInitializeDatabase`、`TestProductionBindFailsBeforeDatabaseInitialization`；`AD::TestCLIModesResolveFromWorkingDirectoryWithoutConfirmation`；`CD::TestCLIModesSelectDatabaseFromWorkingDirectoryDotenv`；`BC::test_console_and_nonwindows_preserved` | redirect 輸出僅由 console mirror 測試間接；`CD` 該測試於 Windows Skip |
| W03 | Capture a Windows GUI artifact in headless mode | 覆蓋缺口 | `AL::TestLaunchModesAndExplicitOverride`（headless 解析）；`BC`（僅 build flag） | 缺：GUI 子系統 exe 於 cmd／PowerShell／supervisor 的 redirect、wait、exit status 與不配置 console 的 harness（task 5.3） |
| W04 | Windows-only release flags | 覆蓋缺口 | `BC::test_desktop_flags_and_artifact`、`test_console_and_nonwindows_preserved`、`test_invalid_platform_or_mode_fails`（make -n）、`test_frontend_must_precede_backend`（build.ps1 字串順序與 `index.html` 檢查） | 缺：Makefile 路徑前端失敗／缺 index 的失敗閉鎖、Go embed 缺資產失敗斷言；`BC` 本輪 5/5 PASS；仍不代表實際平台建置（task 5.1） |
| W05 | An explicit existing database wins | 有測試 | `AD::TestExplicitAbsoluteDatabaseWinsOverCompetingLocations`（process／dotenv 別名順序、空白別名 fall through、DataRoot、整棵樹不變）、`TestExplicitDatabaseUsesCanonicalIdentity`；`AL::TestDataSelectionNeverCreatesOrGuesses`；`CD::TestDatabaseAliasOrderAgreesBetweenResolverAndHarnessDSN` | owner guard 取鎖的 host 接線僅由 `DO`（Windows，NOT RUN）間接 |
| W06 | A legacy working directory is ambiguous | 覆蓋缺口 | `AD::TestDesktopWorkingDirectoryAmbiguityStopsBeforeStorage`（4 案 `startup.path_ambiguous`、樹不變）；`CD::TestDesktopAmbiguousWorkingDirectoryStopsBeforeAnyResource`（無資源、`nativeNextStep` 含 GATEWAY_DB_PATH、legacy 未搬移） | 缺：native 說明含競爭位置的文字與 Windows 呈現未斷言 |
| W07 | Confirm a genuinely new desktop database | 覆蓋缺口 | `AD`／`AL`（ConfirmCreate 真假）；`CD::TestDesktopNewDatabaseCancelLeavesStorageAndServicesUntouched`、`TestDesktopNewDatabaseCancelExitsSuccessfully`（非 Windows 恆拒絕 seam：無資源、無 DB／WAL／SHM／journal、目錄不變、成功退出） | 缺：Windows 原生確認框（顯示位置、預設 Cancel）與接受路徑；兩測試在 Windows Skip（task 1.1） |
| W08 | An unwritable root fails visibly | 覆蓋缺口 | `AD::TestPathFailuresDoNotFallBackOrCreate`、`TestUnreadableLocationsFailInsteadOfLookingEmpty`（POSIX 權限）；`CD::TestDesktopExplicitPathFailureDoesNotFallBackToDefaultDatabase` | 缺：初始診斷目錄不可用案例、native 對話框與非零 exit 斷言、Windows ACL（task 1.1） |
| W09 | Concurrent launch cannot duplicate acquisition | 覆蓋缺口 | `DO::TestOwnerCaseAliasesShareOneActualFileGuard`、`TestHardlinksAreIdentifiedButFreshWritableOwnershipIsRejected`；`DI::TestRepeatedDesktopHandoff`（Windows，NOT RUN） | 缺：兩個真實行程並行啟動，斷言僅一個開啟／migrate、無第二組 collector／delivery／Share listener，secondary 開頁或回報後退出（task 4.1） |
| W10 | A stale owner marker does not block recovery | 有測試 | `DO::TestOwnerDeathReleasesOSLockDespiteStalePID`（子行程持鎖時第二 owner 被拒；kill 後寫入舊 `.pid`，OS 釋放後即可取得） | Windows build tag，**NOT RUN**；「不因 PID 終止行程」無負向斷言 |
| W11 | Handoff is unavailable or untrusted | 覆蓋缺口 | `DI::TestMalformedOrIdleClientCannotControlGateway`、`TestPeerRequiresUserAndSession`（SID／session 不符被拒）；`DO::TestOwnerReservationDoesNotExtendDatabase`（headless owner 的 `RequestOpenSetup` 回錯誤）；皆 NOT RUN | 缺：secondary 在 headless／異 session／驗證失敗時回報已在執行、不啟動 runtime、不開未驗證服務的 cmd 層斷言 |
| W12 | A different application owns the configured port | 有測試 | `CS::TestProductionBindFailsBeforeDatabaseInitialization`（`startup.bind_failed`、bind 前不建 DB）；`CD::TestCLIModesSelectDatabaseFromWorkingDirectoryDotenv`（佔用 port 時停在 bind、樹不變）；`CR::TestReadinessRequiresOwnedEmbeddedResponse`（他人回應不算 ready） | 未涵蓋 desktop 模式的 tray／browser 面；Windows 分支 NOT RUN |
| W13 | HTTP serves while acquisition is unavailable | 覆蓋缺口 | `CS::TestRuntimeFailureProjectsDegradedWithoutRawDetails`（runtime.degraded、無 runtime.started）；`DT::TestMinimalMenuProjectsTruthfulState`（starting／stopping／stop-timeout 選單） | 缺：serving 與 acquisition degraded 在 tray state／選單分列、不以 HTTP running 顯示健康 collector 的斷言 |
| W14 | Tray creation fails before service startup | 覆蓋缺口 | `CS::TestInitialTrayFailureCleansUpBeforeAcquisition`（注入 failing shell：`startup.tray_failed`、未開 db／runtime／listener、tray 已 Close；Windows Skip） | 缺：原生 tray 建立失敗的 native 錯誤對話框與非零 exit（task 4.2） |
| W15 | An early fatal occurs without a console | 覆蓋缺口 | `DT::TestNativeErrorAdvertisesOnlySavedDiagnostic`；`CD`／`CS` 的啟動錯誤碼斷言 | 缺：config／storage／migration／service 失敗→native 對話框，且不顯示憑證、DSN、raw 例外的斷言；無 console 原生呈現（task 1.2／5.3） |
| W16 | An operator uses the menu | 覆蓋缺口 | `DT::TestMinimalMenuProjectsTruthfulState`（六項、ready 前 setup 停用、stopping 時 quit 停用、stop-timeout 顯示 pending phase）、`TestProductURLBoundary`、`TestNativeVersionUsesActualMetadataAndLocalPath`；`AL::TestVerifiedListenerURLs` | 缺：選單動作接到 `/studio/v2`、`/studio/logs`、日誌不可用說明、滑鼠與鍵盤實際開啟（task 4.2） |
| W17 | The browser is closed or minimized | 覆蓋缺口 | 已讀測試中無對應斷言 | 缺：browser 關閉後 HTTP／採集／日誌續行、tray 可重開的 harness 或手動證據（task 5.3） |
| W18 | Explorer rebuilds the notification area | 覆蓋缺口 | `DI` 有 `trayWindowProcedure` 注入點（session-end 測試），無 TaskbarCreated 斷言 | 缺：TaskbarCreated 重新註冊、不重啟 gateway、不複製 worker 的原生驗證（task 4.2／5.3） |
| W19 | Cancel then confirm exit | 覆蓋缺口 | `AS::TestShutdownIsOnceAndUsesOneDeadline`（32 次並行 Begin，每 phase 一次、共用 deadline）；`DT::TestMinimalMenuProjectsTruthfulState`（stopping 時 quit 停用）；`CS::TestStoppingGateAndLateStartupCannotAdmitNewWork` | 缺：取消後服務不變；確認後 tray 與 owner 於資源清理後移除；原生確認框（task 4.3） |
| W20 | The shutdown deadline expires | 覆蓋缺口 | `AS::TestSharedDeadlineRetainsStorageUntilWorkerActuallyEnds`（stop-timeout、DB 不在 worker 下關閉、原協調器完成後回 ErrShutdownTimeout）、`TestTimeoutCannotBeOverwrittenByDelayedNotification`；`CS::TestShutdownTruthAndHandledExitPreserveFailure`；`DT::TestMinimalMenuProjectsTruthfulState` | 缺：繼續等待與二次確認 force termination 的原生流程及警語；headless 不彈窗且非成功回報（task 4.3） |
| W21 | Durable backlog survives exit | 覆蓋缺口 | `DS::TestDurableSampleAckCrashSurvivesReopeningTheDatabase`（SQLite close/reopen 保留 ACKed sample）；`DF::TestWorkerFencingAndShutdownRecoveredRowWithoutDedupeIsUnknownNotResent`（unknown 不重送、不被 stale worker 升格 committed）、`TestWorkerFencingAndShutdownRestartRecoversPriorIncarnationBeforeLeaseExpiry`（重啟 recovery 冪等、不接管活節點） | 低層 durable recovery 已有測試；同一 Windows artifact 的完整退出／重啟及 operator 顯示證據仍缺（task 4.3／5.3） |
| W22 | Multiple shutdown causes race | 覆蓋缺口 | `AS::TestShutdownIsOnceAndUsesOneDeadline`；`DI::TestSessionQueryAndCancellationDoNotStopGateway`（session end 只轉發一次，NOT RUN）；`CS::TestStoppingGateAndLateStartupCannotAdmitNewWork` | 缺：tray 確認與 signal 並行且有 SSE client 時，單一 coordinator 取消 streams、無 double close 或重複 cleanup 的整合斷言 |

小計（僅為上表列數，不是涵蓋率）：`R` 17 列 `有測試`、9 列 `覆蓋缺口`；`W` 4 列 `有測試`、18 列 `覆蓋缺口`；`測試範圍排除` 0 列。

未解的涵蓋證據：

- 檔案 sink 輪替、disk failure／writer stall／恢復、portable pathname swap 已補上 `GD` 對照；仍缺 Windows 原生故障／reparse、完整 tray 顯示，以及 R04 的 sink 失敗不遞迴斷言對照。
- 標準 log／slog／Gin 實際 `raw.suppressed`、ErrorLog（R01）、檔案與 SSE 的秘密斷言（R02）：待對照。
- 成功回應安全標頭（R12）的逐項斷言仍待對照；真實 IPv4 embedded 與 IPv6 socket 測試已有本輪證據。
- SSE 於 shutdown 取消（R22、W22）、host 層 sink 開啟失敗與 tray 狀態分列（R11、W13、W15）、選單動作接線（W16）：待對照。W21 的 SQLite reopen／unknown recovery 已補上低層測試，仍缺完整 Windows artifact 退出與重啟證據。
- Windows 原生項目（W01、W03、W06–W09、W11、W14–W20）全部屬既有 pending task 1.1、1.2、4.1、4.2、4.3、5.1、5.3 的原生 harness 或手動證據；建議在這些任務內完成，不另增功能，也不因本表變更任何 task 狀態。

### NOT RUN／未涵蓋

- 整個 change 的最終 verify、security audit、archive：**NOT RUN**（Windows 與相依任務仍 pending）。上方 test-presence audit 只是測試存在與斷言對照，不等於 scenario verify，也不是 security audit；各列的執行狀態已在該節分開說明。
- Windows 原生 owner／tray／IPC／對話框／Explorer／headless／portable 機器，以及 Linux CLI 本輪執行：**NOT RUN**；無實際 Windows 平台。
- 未接觸 PLC 或正式資料庫；無 deploy、push、commit。

---

## 以下為先前交付歷史

以下內容為先前交付的歷史記錄，其 PASS／FAIL／NOT RUN 僅代表當時狀態；本輪現況以上方 checkpoint 為準。

## 來源與交付邊界

- Change：`openspec/changes/add-windows-tray-runtime-logs`
- 實作基底：`eea0b3e5b4e85bd9868cb2cbceb119dcefcbb189`
- 本輪為依已發布規格重新實作的 source；先前暫存環境失去存取的 commit 不屬本交付，也不借用其測試結果
- 最終 commit、parent、tree、artifact SHA256 以交付 ZIP 的 `DELIVERY.json` 為準，避免文件自我引用尚未存在的 commit hash
- 只建立一個本機 implementation commit；未 push、merge、archive、部署，也未存取真 PLC／正式資料庫
- 本文件的 PASS 限明列的自動化或 source review；Windows 互動桌面、macOS 原生及真正 browser 驗收不在已通過範圍

## 工具與流程

Go 1.25.5、Modern Go Guidelines 0.1.1、golangci-lint 2.5.0、官方 `@fission-ai/openspec@1.2.0` 均隔離使用，不修改系統 Go 或使用者電腦設定。前端依 repo lockfile 還原，React 19／TypeScript 5／Vite 7。

依使用者明示的原 OpenSpec apply workflow 實作；目前 repo 的 Spectra CLI 不可用，**正式 Spectra verify/review/analyze/audit/drift 為 NOT RUN**，不以普通 review 冒充。已執行兩個獨立 source review，分別聚焦 Windows／lifecycle 與 diagnostics／API／UI。

## 實作落點

| 範圍 | Source |
| --- | --- |
| 模式、資料選擇、共用 shutdown budget | `internal/apphost/` |
| Win32 tray／shell、OS owner lock、same-user/session IPC | `internal/desktop/` |
| 固定 safe events、bounded broker／file sink | `internal/diagnostics/` |
| log-only guard、snapshot、SSE | `internal/api/runtime_log_*.go` |
| lazy log viewer、transport、雙容量 model、i18n | `frontend/src/features/runtime-logs/`、`services/runtimeLogs.ts` |
| 真實 entrypoint、bind-first、readiness、safe adapters、cleanup | `cmd/test_ui/gateway_*.go`、`main.go`、`server_runtime.go` |
| 原 worker 的 completion／deadline seams | `internal/datalink/{collector,runtime,groupdelivery,grouppipeline,storage}/` |
| 完整 SPA／GUI-console 建置分流 | `Makefile`、`scripts/build.ps1` |

## 驗證結果

本節以最終收斂時的命令結果為準；原始 log、source manifest、review 結果與可重跑 smoke harness 置於交付 ZIP 的 `evidence/`。

| 命令／驗證 | 結果與限制 |
| --- | --- |
| `go test ./...` | PASS；包含新增 production seams 與原本回歸，不代表外部 PostgreSQL／現場測試已跑 |
| `go vet ./...` | PASS |
| `golangci-lint run ./...` | FAIL；基底獨立重跑有 90 項既有問題，未弱化 config／移除檢查 |
| `golangci-lint run --new-from-rev <base> ./...` | PASS，0 個新增行問題；不替代 full lint |
| 受影響 Go package `-race` | PASS；具體 package／命令见 evidence manifest |
| lifecycle 重複回歸 | PASS，24 個新低層 regression，shutdown 10 次；後續 startup cancellation race 有獨立重跑 |
| 前端 `npm run lint` | PASS |
| 前端 `npm test -- --run` | PASS，203 files／1301 tests，jsdom 不等於 browser |
| 前端 `npm run build` | PASS；logs 為獨立 lazy chunk，既有大 chunk 警告保留 |
| Windows build contract tests | PASS，5 項；未執行 Windows PowerShell native build |
| `make check-lines`／`git diff --check` | PASS；>300 行警告保留，沒有新增 >500 行檔案 |
| 官方 OpenSpec strict validate | PASS；只驗 artifact 結構／規格，不證明產品通過 |
| Linux production smoke | PASS：own embedded logs、99 個 assets 位元組、snapshot／SSE reconnect、403、literal DB path、SIGTERM 0、occupied port fail-before-DB |
| 無 Go/Node runtime PATH | Linux smoke PASS；不是 Windows portable acceptance |
| Windows GUI／console cross-build、PE 檢查 | 結果與產物 hash 見 DELIVERY；不能代替原生操作 |
| Windows native owner／IPC／tray／Explorer／headless handles | NOT RUN，無連線的 Windows 互動環境 |
| Embedded Playwright／真實 browser UI | NOT RUN；已知 browser 執行／loopback 路徑遭環境阻擋，未繞過；未提供假截圖 |
| macOS native、真 PLC、正式 DB、field/deploy | NOT RUN |

一輪合併的低層 race 曾在既有 `TestSenderSlowRemoteFailureStartsSettlementTimeoutAfterAttempt` 的 5ms settlement 預算失敗，沒有 race report；其後隔離的完整 groupdelivery race 與該測試 5 次重跑 PASS。保留原 FAIL 與 rerun，未放寬 timeout、刪測試或改資料語意。

## Review 發現與修正

- Shutdown refresh 可覆蓋 stopping：以同一 state lock 串行 native projection，quit／coordinator 狀態在鎖內重查；controlled interleaving regression 先 RED 再 race GREEN
- Quit 後 HTTP gate 尚未關閉：admission 直接觀察 quit，既有 handler 仍須真正結束後才 close DB
- file-only 在 logger restore 後洩漏 raw fatal／遲到 ready 或 opener error：command exit 與晚到 producer 保留原選定 policy；新增 runCommand 與 post-shutdown regression
- Native post-start fault 無 safe event／錯誤退出：固定 callback 先記錄穩定 code，再請求 quit；handled sentinel 與晚到 fault 不掩蓋非成功結果
- Startup completion 等待 browser：完整發布 HTTP resource fields 後即結束 startup phase；native opener 最多一個未完成工作，主流程可取消等待，不讓 ShellExecute 卡住 worker cleanup
- Windows canonical 祖先可被重新綁定：owner 保留 root-first directory handles，拒絕 reparse、deny write/delete sharing，整鏈固定後重驗 FileID；原生實測仍待補
- Hardlink DB 與 WAL 名称不是同一個 identity：可丟棄 crash probe 證實 alias 可能看不到原 WAL 的已接受資料。Windows writable owner 對無法證明同一 sidecar namespace 的 multi-hardlink DB fail closed；不搬移／刪 WAL。參考 [SQLite 官方說明](https://www.sqlite.org/howtocorrupt.html#_multiple_links_to_the_same_file)
- UI metadata 原先可額外保留 snapshot records：移除這份 hidden array，確保 browser 雙容量上限；file loss interval 與較晚 admission drop 以 sequence 正確排序

Producer inventory：standard log／slog 任意文字只投影固定 `raw.suppressed`；Gin access 用 route template 等允許欄位，panic recovery／HTTP ErrorLog／startup、runtime error callback、native fault、shutdown 有固定 code。未 tee 任意 stdout/stderr；OS/runtime fatal、直接 `fmt` 或第三方自行寫輸出不保證捕捉。舊 `/debug/logs` 與 raw console／query formatter 分開保留。

## 15 項任務：4 完成，11 保留

分類：A＝尚有已知 source／unit 缺項；B＝已實作並有部分實測，但該整項驗收仍缺平台／環境或完整 gate；C＝整項要求的目標平台驗收未跑。勾選以完整任務為單位，不把 source 存在當成全部驗收完成。

| Task | 狀態 | 尚未勾選原因／已具證據 |
| --- | --- | --- |
| 1.1 | B | mode／alias／ambiguity／無寫入 resolver 自動測試 PASS；Windows native 確認取消、真權限與 macOS CLI 未驗 |
| 1.2 | B | bind-first、own HTTP ready、assets、degraded、tray failure、native fault seams PASS；GUI 無 console 的 early fatal 原生呈現未驗 |
| 2.1 | 完成 | secret／未知 raw／8KiB／雙界線／sequence／race PASS |
| 2.2 | B | rotation、alias、disk error/stall、recovery、flush 自動故障注入 PASS；Windows reparse／native disk fault 未驗 |
| 2.3 | 完成 | production adapters 與必要 safe producers 已接線，secret／file-only／late producer／舊 console 回歸 PASS，未捕捉範圍已記錄 |
| 3.1 | 完成 | 實際 IPv4／IPv6、本機 guard、Host／Origin／Fetch Metadata／proxy／CORS 負測試與 bounded snapshot PASS |
| 3.2 | 完成 | replay-live boundary、progress／gap／reset、上限、stalled transport deadline、cancel／race PASS |
| 3.3 | B | lazy route、控制、XSS、i18n／鍵盤 unit、完整前端套件 PASS；embedded Playwright 未跑 |
| 4.1 | B | OS-backed owner／ACL IPC source、portable contract／cross-compile 有證據；Windows concurrent launch／owner death／hardlink failclosed／session 測試尚未原生跑 |
| 4.2 | B | native tray／shell／Explorer source 與注入 seams；無 console 閃窗、keyboard/overflow、browser/Explorer 原生驗收未跑 |
| 4.3 | B | once／shared deadline／真 completion／HTTP gate／可丟棄 SQLite durable reopen/race PASS；native cancel/confirm/force/session 行為未驗 |
| 5.1 | B | scripts／portable contracts／embed smoke／cross-build；Windows fresh-machine PowerShell／standalone 實測未跑 |
| 5.2 | B | 明列的大部分自動 gates PASS；full lint 保留基底 FAIL，Playwright 未跑，不能宣稱全綠 |
| 5.3 | C | 本項整套 Windows 互動桌面／GUI-headless supervisor／同 artifact durable/browser acceptance 未跑 |
| 5.4 | B | 使用與交接文件、scenario matrix 已整理；依賴 5.3 的最終現場交接尚未完成，不 archive |

目前 review 未保留已知可在此環境修正卻未修正的 A 項；这不代表 B/C 的平台行為已被證明。

## 可重跑與後續驗收

一般 build/run、模式、資料選擇與回退請看[使用交接](windows-tray-runtime-logs.md)。交付內的 `implementation.patch` 可在上述 base 使用 `git am`；`implementation.bundle` 是帶 base prerequisite 的 Git bundle，先核對 base 再匯入。套用前保持工作目錄乾淨，不把 patch 套到未知 production tree。

Windows 驗收須使用 DELIVERY 指定 artifact 或從同 commit 按完整 build script 重建，逐 scenario 記錄 OS、session、exe hash、設定／可丟棄 DB、步驟、觀察與截圖。至少涵蓋 owner alias／ancestor pin 與 WAL reopen、Explorer restart、磁碟故障、正常與逾時退出、GUI-headless cmd/PowerShell/supervisor handles。不可拿 cross-build 或 Linux smoke 勾掉 5.3。

48 個 scenario 的逐列 source／test／限制對照在交付 ZIP：`evidence/runtime-log-scenarios.md`（26 項）與 `evidence/windows-lifecycle-scenarios.md`（22 項）。沒有修改 Studio inventory 文件，故本次不觸發其 changelog SQLite 更新。
