# 可採取行動的執行日誌：交付與驗證

原始紀錄：2026-10-29。Change：`make-runtime-logs-actionable`。以下保留當時交付狀態；2026-10-30 的同步複核、修復與實際瀏覽器結果見末節。

## 已執行與未執行

| 項目 | 實際結果 |
| --- | --- |
| 完整前端單元測試 | 一次fresh PASS：417檔／5,117 tests，零fail／skip／todo |
| 專用E2E typecheck、10檔lint、discovery | PASS；探索70案例，不是實際執行 |
| 完整適用範圍ESLint | 1,052檔／18批，0 errors；2個既有no-console warnings：`frontend/tests/e2e/device-recording-workflow.spec.ts` 第167、199行；另6個config exclusions |
| 最終TypeScript／Vite build | PASS；1,116個frontend來源輸入未變 |
| Embedded Go build | PASS：exit 0；108個同步dist assets逐一核對來源hash及binary內bytes，未啟動binary或服務 |
| Static gates | diff/check-lines/analyze/validate PASS；本文件最後修訂後另做文件靜態複查 |
| 整合後backend | 774個選定頂層案例的分批／續跑組合證據及11個package vet；不是單次全綠或整個backend全集 |
| Backend重用 | 4.2核對1,988個Go／native／embed／module輸入hash一致，沒有重新執行backend測試 |
| 真實browser／screenshots | 70案例／72預期截圖全部NOT RUN，實際結果及圖片產出0 |
| 外部PostgreSQL | 10個巢狀subtests NOT RUN：POSTGRES_DSN未提供，不計為PASS |

原始cutover matrix曾出現 `sql: no rows in result set`。後續相同native cycles的診斷版通過，但**原始NoRows原因仍未確定**，不是產品修復證據。另有settlement測試使用明確test-only barrier修正；原失敗保留，產品timeout／backoff未放寬。

1,988輸入一致性是asset同步前的重用邊界；後續差異為embedded assets同步，不能宣稱所有backend輸入仍未變。補充有界驗證：Go embedded build通過，108個資產的來源hash及binary內bytes逐一符合。Binary為72,314,865 bytes，SHA256 `8e22090c04cf3c28991c02b0e4c6cc7eede58f17b144a62a69f8e1cfdd712752`。這只證明建置及資產嵌入，不是服務啟動、browser或端到端測試。

早期fixture假設與E2E型別檢查失敗已在新測試／fixture範圍修正；不把它們偽稱產品RED或抹除原失敗。最終完整frontend PASS是後續一次新執行，不是拼湊套件結果。

## Browser未執行原因與限制

先前觀察到的雲端browser launch、下載及managed-browser loopback限制尚未解除。本次只唯讀核對阻擋連續性，沒有重新嘗試已被拒絕的launch／download，沒有新環境或安全繞過。使用者同意在限制仍存在時保留正式案例並明列未驗證，並未免除其他測試。

這不是所有環境都無法執行的主張。實際鍵盤／焦點、390／1024／1600寬度幾何、截圖外觀、Windows、現場設備、真實PLC及production database均未驗證。

Mock HTTP／SSE驅動實際RuntimeLogsPage和原生Workbench receiver，不代表embedded backend／SQL／device端到端驗證。Native cold delivery GET是讀取，不虛稱沒有任何GET；列操作不執行重試、啟用或DDL。

## 13個scenario／example對照

Browser期待與已執行Go／unit anchors必須分開看。以下所有browser案例都NOT RUN；測試symbol可在repo搜尋。

| 規格scenario／example | Browser family（未執行） | Go／unit anchor |
| --- | --- | --- |
| Recording group fails and recovers | RL-E01 | `cmd/test_ui/delivery_events_test.go::TestProductionDeliveryEventsRetryThenExactReceiptRecovery` |
| No proof of SQL commitment | RL-E02 | `internal/datalink/groupdelivery/transition_observation_test.go::TestDeliveryObservationRequiresExactMatchingReceipt`, `TestDeliveryObservationQueuedUnknownNeverClaimCommit`; `cmd/test_ui/delivery_events_test.go::TestProductionDeliveryEventsReceiptConflictDoesNotClaimCommit` |
| Event capture fails | RL-E03 | `cmd/test_ui/delivery_events_test.go::TestProductionDeliveryEventsFullOrClosedCapturePreservesNativeCommit`; `internal/diagnostics/domain_coalescer_test.go::TestDomainCoalescerCapacityAndDropSeparation` |
| Hostile and secret-bearing fields | RL-E04 | `internal/api/runtime_log_domain_test.go::TestRuntimeLogDomainMemoryFileSnapshotReplayAndLive`; `frontend/tests/unit/runtime-logs/RuntimeLogCompatibility.test.tsx` hostile/nested candidate cases |
| Remote viewing is denied without a proxy bypass | RL-T01 | `internal/api/runtime_log_domain_test.go::TestRuntimeLogDomainGuardPreservesLocalOnly`; `runtime_log_guard_test.go::TestLocalLogGuardConnectionBoundary`, `TestLocalLogGuardRejectsProxyHeaders`; Compatibility snapshot/stream 403 cases |
| Repeated failure and revision change | RL-E05 | `internal/diagnostics/domain_coalescer_delivery_test.go::TestDomainDeliveryRevisionCountsAndLateSummaryRemainSeparate`, `TestDomainDeliveryRecoveryMatchesExactEffectAndRevision`; `domain_broker_sink_test.go::TestDomainBrokerSourceWorkspaceRevisionIsolatesOccurrencesAndRecovery` |
| High-cardinality burst | RL-E06 | `internal/diagnostics/domain_coalescer_test.go::TestDomainCoalescerAggregateRecordBudget`, `TestDomainCoalescerOutputQueueBounds`, `TestDomainCoalescerInFlightOutputConsumesSharedSlot`; `domain_broker_test.go::TestDomainBrokerCoalescerCapsAndRevisionIsolation` |
| Restart after a previous commit | RL-E07 | `cmd/test_ui/delivery_events_test.go::TestProductionDeliveryEventsRestartDoesNotReplayOldReceipt`; `runtimeLogSession.test.ts` process-reset/duplicate-cleared-record case |
| Filters change with an old response in flight | RL-T02, RL-T03 | `runtimeLogDomainFiltersSession.test.tsx` exact A 9007199254740993 → B case and late A snapshot; `internal/api/runtime_log_domain_filters_test.go::TestRuntimeLogDomainFiltersActualSnapshotReplayAndLive` |
| Retention gap prevents complete fault history | RL-T04 (also RL-E03/RL-E07) | `runtimeLogDomainFiltersSession.test.tsx` retained fault without recovery case; `RuntimeLogCompatibility.test.tsx` file-loss/capture-gap/eviction/reset case |
| Operator opens the exact failed group | RL-N01, RL-N05 (390px) | `runtime-log-settings-navigation.test.tsx` “spec W/G/R1 opens exact current G/R2 without mutation or historical relabeling” |
| Configuration was deleted or renamed | RL-N02, RL-N03, RL-N04, RL-N06 | `runtime-log-settings-navigation.test.tsx` deleted G1/same-name G2, rename/Back, foreign response and late entry cases |
| Blocked browser execution is reported without a false pass | 全部70案例／72預期截圖：NOT RUN | 專用E2E typecheck、lint、discovery及本目錄逐案例JSON |

[完整逐案例JSON](make-runtime-logs-actionable-browser-cases.json)列出70案例、語系、寬度、重跑命令與72預期截圖名稱；全部 `generated: false`。Discovery報告的skipped不是執行結果。第13個阻擋揭露scenario關聯全部70筆，不虛構第71個測試。

## 重跑方式

前提：依repo既有package-lock準備相符Node依賴與Playwright要求的Chromium，並允許本機127.0.0.1:4173。不可停用安全限制來替代環境前提。既有 `frontend/playwright.config.ts` 會啟動Vite；確認該port不是不相關服務。

從repo根目錄進入 `frontend/`：

```sh
./node_modules/.bin/tsc --noEmit -p ../.task-evidence/logs-4.2/tsconfig.json
./node_modules/.bin/playwright test tests/e2e/runtime-logs-evidence.spec.ts tests/e2e/runtime-logs-transport.spec.ts tests/e2e/runtime-logs-navigation.spec.ts --list
./node_modules/.bin/playwright test tests/e2e/runtime-logs-evidence.spec.ts tests/e2e/runtime-logs-transport.spec.ts tests/e2e/runtime-logs-navigation.spec.ts --project=chromium
```

`--list`只探索，最後一行才實際執行；單例命令在JSON的 `rerun_cwd`／`rerun_command`。en／zh-TW與1024／1600及額外390px期待由spec定義。`*-mock-data.png`透過 `testInfo.outputPath()` 寫入Playwright test-results；HTML報告在 `frontend/playwright-report/`。重跑結果須另記日期／來源版本，不回填成本次已執行。

## QA歸屬與可攜性

- 5個 `frontend/tests/e2e/fixtures/runtime-logs-*`：HTTP／SSE／state／實際host準備。
- `frontend/tests/unit/runtime-logs/runtimeLogBrowserFixture.test.ts`、`RuntimeLogBrowserHost.test.tsx`：fixture與host可用性。
- `frontend/tests/e2e/runtime-logs.helpers.ts`及3個 `runtime-logs-*.spec.ts`：正式案例與截圖期待。
- `.task-evidence/logs-4.2/tsconfig.json`保留原路徑，因相對extends／typeRoots／include依賴該位置。附件記錄全部12個QA來源hash。

本交付紀錄保留repo-relative測試位置、來源hash、重跑指令及實際驗證狀態，供讀者核對與重跑；不含設備或資料庫憑證。

## 2026-10-30 macOS 同步複核

整合本機 `2d8ee159` 與遠端 `b6028a40`，保留雙方歷史。兩個共用 Go 檔自動合併，沒有文字衝突；審查同時確認 PLC scalar count 與 MC 更新前驗證仍保留。本節結果對應本節所在提交的產品與測試來源，未將舊雲端紀錄改寫為已執行。

修正皆限於本 change：

- 日誌佇列在空佇列檢查與交出排出權限間可能遺留已接受事件。新並行回歸在未改產品碼時，以 `-race -count=10` 重現 `queued=1 admitted=1 dropped=0`；修正後相同十輪通過。交接改在同一 mutex 下完成，回呼維持鎖外，每批上限不變。
- 每筆 SSE 更新重繪所有舊列，原 RL-E06 單獨執行仍超過 60 秒。重用內容未變的事件列後，原四個大量資料案例均通過。
- 390px 訊息欄原僅占表格約 17%／22%，文字難以閱讀。窄版改為單欄事件，保留表格語意、欄名與桌面版；RL-N05 補實際訊息欄寬度檢查，先失敗再修正。
- 兩個新 Go fixture 改用既有 `canonicalTempDir`，解決 macOS 系統暫存路徑的 symlink 檢查；未放寬日誌檔所有權保護。窄版 E2E 的下拉選單改依可及名稱與 combobox role 定位。

| 檢查 | 本次實際結果與界線 |
| --- | --- |
| 完整後端套件集合 | 分批涵蓋全部 61 個 package：56 PASS，5 個無測試；初輪兩個新 fixture 失敗後已修正，`cmd/test_ui` 完整重跑 PASS。外部環境及平台案例仍依原條件 SKIP，含子測試共 519 個，不計入 PASS。 |
| 併發檢查 | 10 個受影響 package 的選定測試 `-race` PASS；佇列修正後另重跑 diagnostics、API、主程式相關測試，213 個頂層／671 個含子測試 PASS。新交接回歸十輪 PASS。 |
| 完整前端單元測試 | 426 檔／5,181 tests：5,178 PASS、3 FAIL；同三個失敗在同步前 `2d8ee159` 隔離副本可重現，不擴修。 |
| 最終日誌單元測試 | 14 檔／456 tests PASS；完整前端既有失敗位於 `address-conflicts.test.ts` 及 `step4-first-activation.test.tsx`。 |
| 最終真實瀏覽器 | 原 70 案例全數 PASS、零 skip／flaky；實際產出 72 張截圖，涵蓋 en／zh-TW、390／1024／1600px、鍵盤焦點、歷史版本定位、離線／重連／缺漏及大量紀錄。 |
| 瀏覽器證據限制 | 使用實際 RuntimeLogsPage／原生設定 receiver 與 mock HTTP／SSE；不是 embedded backend／SQL／PLC 端到端證據。原逐案例 JSON 仍是 2026-10-29 的 NOT RUN 紀錄。 |
| 靜態檢查與建置 | 完整 Go vet、完整前端 ESLint、最終變更檔 ESLint、TypeScript、專用 E2E typecheck、Vite 與 embedded Go build PASS；完整 ESLint 保留原兩個 no-console warnings。108 個 dist 資產逐一核對同步內容與 binary 內 bytes。 |
| 全域 Go lint | 初輪 `golangci-lint run ./...` 失敗，回報 1,752 issues，含存量與本批診斷；未宣稱全綠或進行跨範圍整理。最終修改 Go 檔的 gofmt 與受影響 package vet 另檢查。 |
| 未執行 | Windows、真實 PLC／現場、外部資料庫與 embedded browser 到 SQL 的端到端驗收仍為 NOT RUN。沒有 push 或部署。 |

瀏覽器重跑沿用上節三個 spec，另指定 `--workers=4`；未增加替代產品路由、proxy 或操作按鈕。日誌測試、瀏覽器截圖與 native simulator／SQLite 的 production identity 證據分別判讀。
