# 2026-11-06 審查修復交接

本輪依使用者「修」的授權處理前次 review 的 R1–R12；原定功能需求不外擴，原18項完成紀錄保留，另以單一8.1任務追蹤這次修復。所有12項都有修正前的行為失敗與修正後通過證據。主代理親自核對相對本輪凍結工作樹的差異、原規格、呼叫端與結果；未提交或封存。

## 修正與回歸

| 項目 | 最終行為與測試 | 修正前／後紀錄（位於 `/tmp/decouple-fix-20261106/`） |
| --- | --- | --- |
| R1 短 frozen 漏記 | 200ms門檻、900ms恢復、1秒掃描仍記錄200–900ms；儲存失敗可重試且不重複。`TestFrozenRecoveryBetweenScansIsRecorded` | `r1-red.log`／`r1-green.log`、最終套件與race |
| R2 部分儲存 | SQL採集欄位、設備群組與點位指派同一交易；任一側失敗都回滾，不同config DB拒絕。記憶體路徑只在interval成功後套用欄位。`TestFailedIntervalAssignmentDoesNotSaveOtherAcquisitionSettings`及同檔3項對照 | `r2-red.log`／`r2-green.log`、`owned-packages-final2.log` |
| R3 並行覆蓋 | 只更新明確提供的採集欄位，成功後重讀；Studio與一般設備API共用device.Service序列化最新影像的runtime套用。兩個入口的兩種先後順序均受保護。`TestConcurrentPartialAcquisitionSettingsPreserveBothUpdates`、`TestStudioAndDeviceAPIsApplyLatestAcquisitionSettings` | `r3-red.log`／`r3-green.log`；`r3-image-red.log`／`r3-image-green.log`；`r3-runtime-red.log`／`r3-runtime-green.log`；`cross-permanent-red.log`／`cross-permanent-green.log` |
| R4 過早釋放DB | Writer.WaitClosed觀察原始關閉，production fanout與runtime在共用DB關閉前等待真正完成，原呼叫deadline仍即時返回。`TestProductionShutdownWaitsForDatabaseTargetBeforeClosingSharedStorage` | `r4-red2.log`／`r4-green.log`、`focused-race.log` |
| R5 寫入結果時間 | `last_write_at`使用實際嘗試完成時間；history.ObservedAt與既有grouped last_flush_at保持各自語意。`TestWriter_LastWriteTimeReflectsDelayedAttempt`覆蓋成功／失敗、分組／非分組 | `r5-red.log`／`r5-green.log` |
| R6 失敗原因 | outcome顯示安全的權限／連線／結構分類、建議與退避秒數；未知原因及close failure也不回傳raw exception；history仍保留diagnostic。`TestWriter_FailedOutcomeProvidesSafeActionableReason`、`TestWriter_CloseFailureKeepsOperatorOutcomeSafe` | `r6-red.log`／`r6-green.log`；`r6-close-red.log`／`r6-close-green.log` |
| R7 FATEK保留映射 | mixed列使用穩定item ID；移除較早項目後，保留項目的mapping不進DELETE queue。真derive→reducer→autosave測試 | `frontend-r7-red.log`／`frontend-r7-green.log` |
| R8 FATEK定位 | 非首項R10/uint16顯示自身幾何；沒有匹配項目時拒絕locator／partial editor hydration，MC路徑保留 | `frontend-r8-red.log`、`frontend-r8-partial-red.log`／對應green |
| R9 超64項編輯 | mixed count依項目數，homogeneous仍1–64；已保存70項編輯可接受count70的ACK | `frontend-r9-red.log`／`frontend-r9-green.log` |
| R10 正規化ACK | 比較實際canonical request／ACK及投影；尾端空白可正常儲存；異workspace／device、舊revision與錯誤項目ACK仍鎖定 | `frontend-r10-red.log`／`frontend-r10-green.log`、`frontend-editor-guards-green.log` |
| R11 重疊與拒絕 | 儲存前重用protocol-aware重疊驗證；明確400/server validation拒絕可修正重送，模糊網路失敗仍unknown鎖定 | `frontend-r11-red.log`／`frontend-r11-green.log` |
| R12 文件URL | 兩個anomaly annotations加`/v1`並經make gen-docs產生；文件URL須存在於production NewRouter。`TestAcquisitionAnomalySwaggerURLsMatchProductionRoutes` | `swagger-red.log`／`swagger-green.log`；generated JSON其餘paths、definitions與頂層欄位語意無差異 |

R4首次測試的替身曾重複close channel；修正替身且產品尚未修改時重新取得`r4-red2.log`，只以後者作有效RED。並行跨入口早期overlay FAIL已由正式回歸取代。初次前端typecheck發現新增測試對readonly欄位賦值，改為immutable copy後通過；這些setup問題未當產品缺陷。

## 最終驗證

所有Go命令均使用本輪macOS arm64的Go 1.27.3、`GOPROXY=off GOFLAGS=-mod=readonly GOCACHE=/tmp/decouple-20261106/random-go-cache`及`-count=1`；go.mod目標仍為1.25.5。以下都是本輪結果，沒有把先前全套PASS冒充此次重跑。

| 結果 | 指令與範圍 | 證據 |
| --- | --- | --- |
| PASS | `go test ./internal/datalink/anomaly ./internal/datalink/device ./internal/datalink/point -count=1` | `owned-packages-final2.log` |
| PASS | `go test ./internal/datalink/dbtarget ./internal/datalink/runtime ./cmd/test_ui ./internal/api -count=1 -timeout 180s`中的dbtarget/runtime套件 | `delivery-api-aggregate.log`；整體命令因cmd套件總逾時而非PASS |
| PASS | `go test ./cmd/test_ui -count=1 -timeout 360s` | `cmd-full-final.log`：176.875秒 |
| PASS | `go test ./internal/api ./internal/api/handlers -run 'Acquisition|DeviceHandler_Update|StudioV2WorkspaceDevicesHandler_Update|StudioAndDevice|SameTransportConnectionConfig' -count=1 -timeout 120s` | `api-final-selection.log`：兩個套件均有選到測試並通過 |
| PASS | device／Studio handler更新與跨入口相關選測 | `device-handlers-final2.log`，包含原有失敗／not-running／restart-required行為 |
| PASS | anomaly/device/handler的凍結與並行更新`-race`；production關閉`-race` | `owned-race-final2.log`、`focused-race.log` |
| PASS | 受影響anomaly/device/point/handlers/dbtarget/runtime/cmd/test_ui/api的`go vet` | `vet-final2.log`、`delivery-api-vet.log` |
| PASS | `npx --no-install vitest run tests/unit/workbench-v2 tests/unit/scatteredRangeForm.test.tsx tests/unit/scatteredRangeReviewFixes.test.tsx tests/unit/services/studio-v2-locator-mixed-rule.test.ts tests/unit/services/studio-v2-locator-acquisition.test.ts` | `frontend-consumers.log`：449檔、4734測試 |
| PASS | `npx --no-install tsc --noEmit`，包含新增瀏覽器fixture | `tsc-with-browser-host.log` |
| PASS | 修改前端source／test的ESLint及新增browser fixture ESLint | `frontend-eslint.log`、`frontend-eslint-final-test.log`、`browser-fixture-eslint.log` |
| PASS | `npm run build` | `frontend-build.log`；保留既有大chunk警告，未為此擴大拆包，未同步embed或部署 |
| PASS | `make check-lines`、`git diff --check` | `check-lines.log`；歷史超500行檔不增行；後續文件加入後另做最後檢查 |
| FAIL，既有問題 | `golangci-lint run ./internal/datalink/anomaly ./internal/datalink/device ./internal/datalink/point ./internal/api/handlers` | `lint-final2.log`：11項，七個涉及檔案均與本輪起始雜湊相同，無一指向此次修改的檔案；未擴大修復 |
| NOT RUN | 整個repository Go／前端全套、真PLC、正式外部DB、Windows、部署／現場 | 本輪只執行上述受影響範圍；先前本機實測仍屬歷史證據 |

第一次cmd套件總逾時為180秒，當時正執行`TestProductionWordFormatUDPZeroShortAndLegacy`的第4秒；先前該套件整體曾需212.555秒。未把這個總預算不足誤當產品失敗。首次aggregate的API雖通過，但device相關source在執行中更新，因此最終API證據採新版本的獨立選測。

## 實際瀏覽器與截圖

使用專用agent-browser session及Playwright headless shell，Vite只綁127.0.0.1:43169，proxy指向本機不可用埠。新`recording-mixed-review-host`使用正式已保存編輯器、reducer、autosave hook與模擬Axios adapter；所有API呼叫在記憶體截取，並不是完整App、真API或設備驗收。fixture的acknowledged observer允許測試回覆，完整權限／revision保護另由正式單元與整合測試驗證。

- `Pump `經實際輸入、儲存，送出與ACK均為`Pump`，畫面顯示已儲存；[1440截圖](evidence/review-fix-20261106-normalized-saved-1440.png)。
- 70項FATEK仍送70項、ACK count70，saved=true、unknown=false、輸入可編輯；[上方用量](evidence/review-fix-20261106-large-saved-1440.png)與[底部已儲存](evidence/review-fix-20261106-large-saved-bottom-1440.png)。
- D3 uint32／D4 int16顯示重疊錯誤，Save disabled、輸入可修正、API請求0；[1440截圖](evidence/review-fix-20261106-overlap-editable-1440.png)。
- 模擬400 validation顯示安全錯誤、可編輯與重送，無unknown；第二次改`Corrected`成功，實際PUT兩次；[錯誤可修正](evidence/review-fix-20261106-validation-retry-1440.png)與[390成功截圖](evidence/review-fix-20261106-retry-saved-390.png)。
- 主代理已開圖檢視上述成功、錯誤與手機畫面；沒有修改排版、CSS或擴充設計驗收。

## 修改檔案

本輪相對凍結工作樹的38個source／test／generated檔如下，另更新tasks.md的修復追蹤，新增本交接、6張截圖與證據manifest。保留原有staged／unstaged／untracked WIP，未整包覆寫。

```text
cmd/test_ui/target_writer.go
cmd/test_ui/target_writer_shutdown_regression_test.go
docs/swagger/docs.go
docs/swagger/swagger.json
docs/swagger/swagger.yaml
frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeSavedRangeEditor.tsx
frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
frontend/src/pages/datalink/workbench-v2/studioV2LocatorHydration.ts
frontend/tests/e2e/fixtures/recording-mixed-review-host.html
frontend/tests/e2e/fixtures/recording-mixed-review-host.tsx
frontend/tests/unit/workbench-v2/scattered-locator-regression.test.ts
frontend/tests/unit/workbench-v2/scattered-mapping-identity-regression.test.tsx
frontend/tests/unit/workbench-v2/scattered-saved-editor-regression.test.tsx
internal/api/handlers/acquisition_anomaly_handler.go
internal/api/handlers/device_handler.go
internal/api/handlers/studio_v2_device_acquisition_concurrency_regression_test.go
internal/api/handlers/studio_v2_workspace_devices_handler.go
internal/api/router_acquisition_anomaly_swagger_test.go
internal/datalink/anomaly/detector.go
internal/datalink/anomaly/detector_frozen_between_scans_regression_test.go
internal/datalink/dbtarget/delivery_outcome_test.go
internal/datalink/dbtarget/writer_group_flush.go
internal/datalink/dbtarget/writer_outcome_time_regression_test.go
internal/datalink/dbtarget/writer_retry.go
internal/datalink/dbtarget/writer_retry_summary_regression_test.go
internal/datalink/dbtarget/writer_snapshot.go
internal/datalink/dbtarget/writer_workers.go
internal/datalink/device/service.go
internal/datalink/device/service_acquisition_concurrency_regression_test.go
internal/datalink/device/service_acquisition_update.go
internal/datalink/device/service_crud.go
internal/datalink/device/sql_repo_acquisition.go
internal/datalink/device/sql_repo_acquisition_patch.go
internal/datalink/point/service_device_acquisition_settings_regression_test.go
internal/datalink/point/service_device_polling_group.go
internal/datalink/point/sql_repo_device_acquisition.go
internal/datalink/runtime/service_lifecycle.go
```

## 限制與範圍核對

- R3共用runtime套用鎖以production共用的一個device.Service為邊界，不主張跨獨立服務實例或啟停流程皆已序列化；同服務不同設備的套用也會依序進行。
- R2原子性涵蓋此次採集欄位、群組及點位；未擴大重構所有設備欄位的既有儲存契約。
- 仍維持所有進行中異常保留、5000為歷史清理目標，LTN0不支援且不加入0002；durable groupdelivery設計未變。
- 前次verify唯一部分覆蓋仍是字面值120的request split斷言，未據此新增產品完成條件。
- 範圍外靜態觀察：既有schema ensure failure的raw diagnostic投影另有風險，見暫存delivery報告；不屬此次write／flush outcome修復，未自動實作。
- 原始proposal／design／delta specs及原18項任務語意保留；[原始交接](apply-completion-evidence.md)、[情境測試對照](scenario-coverage.md)提供先前背景。本文件只補充本輪修復結果，不將先前完整驗收改寫成此次結果。

## 完成紀錄

12項修復的必要檢查已完成，原18項加上本次修復追蹤共19項。CLI拒絕對已完成任務補登檔案，因此以8.1的明確檔案清單記錄此次修復，不重開原任務或改寫完成證據。主代理核對本輪38個source／test／generated檔的差異，generated Swagger語意只改兩個路徑，並檢視6張實測截圖。所有source在最終檢查後保持相同指紋；原index與HEAD保持不變，起始檔案中除列出的修復與追蹤範圍外沒有修改。最終檔案與log指紋見 [manifest](evidence/review-fix-20261106-manifest.json)。

未archive、commit、push或部署。可接續 `$spectra-verify decouple-acquisition-from-delivery` 與 `$spectra-review decouple-acquisition-from-delivery`，準備妥當後才封存；不要把仍存在的既有lint問題或未執行的現場驗收列為PASS。
