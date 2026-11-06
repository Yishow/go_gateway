# 驗證證據（decouple-acquisition-from-delivery）

最新接續結果見 [本次 apply 交接](apply-completion-evidence.md)；下文的歷史狀態保留原樣，不代表本次完成狀態。

歷史證據日期：2026-11-05。環境：開發機（macOS），只用 loopback、本機 MC 3E peer 與可丟棄 SQLite。以下數字皆為**本機單次量測**，不代表實機 PLC、正式資料庫或現場驗收；PostgreSQL／MySQL／SQL Server 真實目標、實體 PLC 與現場驗收皆**未執行**。截圖與原始紀錄保存在執行當次的 agent scratchpad（工作階段本機路徑，未提交至 repo）。

下列各節至「已知限制」保留當時版本的證據；修復後的最新結果與未完成項目見末節「2026-11-06 審查修復後覆核」。歷史截圖與全套通過結果不能代表修復後版本。

## 1.3／6.1 連接器不可達時的採集吞吐

- 1.3（production wiring）：1 台 MC 設備、32 點、1 秒週期，1 個 tag 映射到不可達 PostgreSQL（127.0.0.1:1）。暖機後 22.002 秒內 `collected_total` 由 128 增至 832（+704），平均 **31.997 筆／秒**，高於設計條件 32 × 95% = 30.4 筆／秒。測試：`GATEWAY_TASK13_THROUGHPUT=1 go test ./cmd/test_ui/ -run '^TestProductionSnapshotWorker_MCUnreachableThroughput$' -count=1 -v`。
- 6.1（端到端，2 台設備共 6 點，理論 6 筆／秒，寫入間隔 5 秒）：

| 階段（各 30 秒） | collected 增量（/秒） | SSE value（random／safe） | 目標表列數 | 寫入歷史 |
| --- | --- | --- | --- | --- |
| 可達 | 180（5.97） | 120／60 | +36 | success 36 |
| 不可達（`chmod 000` 目標 SQLite） | 186（6.17） | 124／62 | 無法開啟 | failed 5 |
| 恢復 | 180（5.97） | 120／60 | +37 | success 37 |

  不可達期間退避間隔實測約 +1、+3、+4、+9 秒（受 1 秒 flush tick 與 5 秒間隔量化影響），恢復後每 5 秒寫入一批、每 mapping 一列最新值；斷線期間不補寫（設計範圍外）。

## 6.1 採集異常（真實事件，無模擬回應）

停止 safe-v1 設備的 MC peer 後：`disconnected`（connection-failed）與 D10、D11 兩筆 `stale`（no-successful-read，門檻 5000 ms）開立，摘要 API `active_count=3`；首頁按鈕顯示「採集異常 3」（warning 樣式），modal 類型／狀態篩選正確。另一台設備不受影響（10 秒內 11 次 0403）。重啟 peer 後三筆事件同時結束並記錄持續時間。frozen 未觸發測試（門檻為 0＝停用）。

## 6.1 單一設備取值間隔

把 safe-v1 設備間隔改為 3000 ms（`runtime_apply_status=applied`，未重啟）：runtime status 該設備 configured/effective 3000/3000，另一台維持 1000/1000；peer 30 秒內收到的請求由 30 降為 10，另一台維持 30。

## 5.4／6.1 分散位址 random 讀取

- 5.4：q-iqr 設備的 mixed 規則（D39、D40、D100、D200、D1000、D2000）加既有 D0–D31，計畫為 1 個 random 請求（38 units）；模擬器實測 1 個 0403、38 個樣本皆 Good。測試：`GATEWAY_MIXED_SIMULATOR=1 go test ./cmd/test_ui/ -run TestProductionScheduledAcquisitionMixedRuleWireRequests -count=1 -v`（repo 內 wire peer 只支援 0403 word 項目）。
- 6.1：以 5.5 前端在首頁儲存的 mixed 規則（D100 int16、D200 int32、D300 uint16、D350 float32），自寫 peer 記錄每秒一次 `0403 words=[100,300] dwords=[200,350]`，共 874 次皆成功；SSE 數值與 peer 的 word 組合逐項相符（含 dword）。

## UI 截圖（1440／768／390）

- 4.2：首頁異常按鈕三態與 modal（全部、篩選、摘要失敗＝狀態未知、列表失敗＋重試）；1440 深色模式。異常資料為種入拋棄式 DB 後由真實 API 讀出，失敗狀態以 Playwright 攔截回應模擬。
- 4.3：MC／Modbus／FATEK 連線設定、欄位驗證錯誤、ⓘ 說明；真實 PUT 儲存並讀回。runtime 有效值畫面曾以攔截回應模擬，之後已補上 `/runtime/status` 的 `configured_interval_ms`／`effective_interval_ms` 輸出（見下方修正）。
- 5.5：分散位址模式初始、173／192 提醒、193／192 超過上限說明（仍可儲存）、safe-v1 停用原因；真實儲存 mixed 規則並讀回。
- 6.1：真實事件的異常按鈕與 modal（進行中、已結束）。

## 驗收發現並已修正的問題

| 問題 | 根因 | 修正與證據 |
| --- | --- | --- |
| 失敗啟用後同一進程舊式資料庫輸出不寫入（無聲） | 快照到期時間隨每筆新值往後推，持續有值時永遠不到期 | 待寫快照保留原期限、只覆寫值；所有不寫入情況記錄略過原因。`cmd/test_ui/bug_a_activation_test.go`（422 → 建表 → 啟用 → 必須寫入）、`internal/datalink/dbtarget/writer_starvation_test.go` |
| 有已存規則的設備「讀取點位」對話框打不開 | locator entry 白名單不接受後端新增的 `layout_mode`／`mixed_items` | `frontend/src/services/studioV2LocatorEntry.ts`；`frontend/tests/unit/services/studio-v2-locator-mixed-rule.test.ts` |
| `/runtime/status` 未輸出設定／有效間隔 | handler 回應未帶欄位 | `internal/api/handlers/runtime_handler.go`；`TestRuntimeCollectorResponseExposesConfiguredAndEffectiveInterval` |
| Studio v2 只改採集設定也回 restart-required | 以「請求是否帶 connection_config」判斷 | 改比對儲存前後傳輸設定（忽略 `acquisition`）；`TestStudioV2WorkspaceDevicesHandler_UpdateAcquisitionOnlyAppliesWithoutRestart` |

## spectra-review 後的修正

- 每筆舊式寫入各自 5 秒期限（原本整輪 flush 共用一個期限）：`TestWriter_HealthySlowMappingsDoNotShareOneFlushDeadline`（40 筆各 150ms 全部成功）。
- mapping 刷新重用已列出的連接器、tag 批次讀取，刷新在獨立 goroutine、不阻塞 worker 喚醒：`writer_review_refresh_test.go`。
- 舊版 `PUT /devices/:id` 依儲存前後傳輸設定判斷是否同步 runtime：`TestDeviceHandler_UpdateAcquisitionWithUnchangedConnectionResyncsRuntime`。
- 設備列表只讀一次輪詢群組：`TestServiceList_HydratesIntervalsWithOneGroupRead`。
- 移除設備時停止其專屬輪詢群組 ticker：`TestSchedulerRemoveDeviceOwnedPollingGroupsStopsOnlyOwnedTickers`。
- 前端：locator 接受 `mixed_items_error` 以便修復、mixed item 欄位與後端一致；MC 64-bit 用量計 4 點並去重；空項目與損毀項目顯示可行動錯誤；分散位址不顯示連續範圍摘要；異常原因分隔符 i18n；列錯誤以 `aria-describedby` 關聯。

## 歷史全套檢查（2026-11-05，當時版本）

- `go test ./...`：PASS
- `cd frontend && npm test -- --run`：PASS（622 檔、7216 測試）
- `cd frontend && npx tsc --noEmit`：PASS
- `bash scripts/check_file_lines.sh`：PASS（legacy 超長檔未增行）

## 已知限制

- 寫入群組擁有的輸出會先在診斷顯示 `buffer_accepted`，再由 worker 以 `write_group_owned` 略過並記錄；若要在接受前判斷，需要一個與 durable Apply 同步、無 I/O 的 ownership projection（跨 `grouppipeline`／`groupdelivery`），否則會讓 ingest 執行 SQLite 查詢、違反「採集不被寫入阻塞」。列為後續工作。
- `internal/datalink/groupdelivery` 的 `TestSenderRestoreSlowerThanSettlementBudgetStillSettles` 在 `go test ./...` 全套負載下曾偶發失敗、單獨 `-count=3` 通過；該 package 未被本 change 修改。
- `RuleEditor.tsx` 的部分欄位標籤為硬編碼中文（HEAD 既有，非本 change 引入）。

- stale 門檻儲存後無法清回預設（後端 `*int` 省略即不變）；前端已阻擋並說明。
- 設備恢復後約 26 秒才重新連線並結束異常（重連退避），未深究。
- SQLite 無法開啟時錯誤摘要顯示 driver 的「out of memory (14)」字串，具誤導性。
- FX5 與 FATEK 的 mixed 採集只驗到排除規則層；計畫與模擬器只量測 q-iqr。

## 2026-11-06 審查修復後覆核

目前為 14／18 項完成，尚不可 archive；未完成項目以 [tasks.md](tasks.md) 的 3.2、4.2、5.1、6.1 為準。本節取代歷史段落對目前版本的驗收結論，未改變既定支援邊界。

### 修復範圍

- 採集與排程：Q／iQ-R 位址進位、FATEK random 加權與禁止範圍、safe-v1／未宣告／Modbus mixed 實際排除、逐項格式與自動間隔回縮。
- 資料與識別：同設備 500／2000 ms 群組升級取 500 ms 且不啟用停用點位；mixed `item_id` 移址、交換及同位址替換維持身分，更新與失敗回復保護並行審查決策及 SQLite mapping。
- 異常與輸出：掃描間短暫異常、frozen 重啟回復、結束事件清理；worker 取消、退避、連線逾時及群組 flush 邊界。
- 前端：分散項目中繼資料往返、異常持續時間；記錄編輯器錯誤文字及「到來源表格」清除篩選後聚焦正確項目。

### 實際執行結果

以下為同一工作階段的既有執行紀錄；覆核時兩個 Spectra scope 快照仍為 current，補充 65 個檔案的 SHA-256 全部相符，因此沿用結果，未為文件更新重跑產品全套。

| 結果 | 指令與範圍 | 證據及限制 |
| --- | --- | --- |
| PASS | `go test ./... -count=1` | 最後一輪 Go 全套通過；之後只新增記憶體決策回復測試，已另跑 race。先前一次全套的 `TestWideCounterParity` 在十分鐘切換時失敗，該範圍外測試重跑三次通過。 |
| PASS | `go test -race ./internal/api/handlers -run '^TestWorkspaceMixed' -count=1` | mixed 項目識別、替換、審查決策及回復測試；不代表其他 handler 的 race 驗收。 |
| PASS | `go test -race ./internal/datalink/sourcerule -run '^TestMixedReviewDecisionMemoryRollback$' -count=1` | 記憶體審查決策回復與並行編輯保護。 |
| PASS | `GATEWAY_TASK13_THROUGHPUT=1 GATEWAY_MIXED_SIMULATOR=1 go test ./cmd/test_ui -run '^(TestProductionSnapshotWorker_MCUnreachableThroughput\|TestProductionScheduledAcquisitionMixedRuleWireRequests\|TestProductionScheduledUnsupportedMixedRuleIsNotRead)$' -count=1 -v` | 本機 MC peer：32 點、22.002 秒增加 704 筆、31.997 筆／秒；不可達連接器拒絕 5 次。mixed 為一個 0403、38 units，unsupported mixed 未發出讀取。採集計數與 buffer acceptance 不等於 SQL committed。 |
| FAIL | `cd frontend && npm test -- --run` | 兩輪皆 633 檔中 632 通過、1 失敗；7277 項中 7276 通過。`recording-shared-configuration-navigation.test.tsx:108` 在修復前隔離快照亦可重現，依範圍限制未修改。 |
| PASS | `cd frontend && npx tsc --noEmit`、相關檔案 ESLint；相關 Go 套件 `go vet` | 局部靜態檢查，不推論完整 lint／build 通過。 |
| NOT RUN | 修復後 1440／768／390 瀏覽器截圖與完整 UI 端到端重驗 | 沙箱內瀏覽器因 Mach port 權限無法啟動；未提權。歷史截圖只代表舊版本。 |

Go 全套紀錄 SHA-256：`539240e2db6cb724bba4917c00103ce35c571fe0e6f4bf19014cf4b15e3ea9ec`；前端第二輪紀錄：`3c1ddc75f5661fd4d6c3d601a09b40e19c0bf011ff1aeb467c1e6de607626635`；最後 loopback 紀錄：`a12a693b3cb6a45661131298710b6f41812b5712440cfb4a60393147e74bdfae`。原始紀錄留於本工作階段本機暫存，未加入 repo。

### 未解與範圍限制

- 3.2：超過 5000 筆全部為進行中事件時，嚴格上限與不得刪除進行中事件衝突；目前保留進行中事件，未自行改規格。
- 5.1：LTN0 範例與 iQ-R `0002` 排除邊界衝突；依使用者決定維持目前支援範圍並暫列未解。
- 4.2／6.1：修復後 UI 截圖與完整端到端未完成，前端全套仍有上述失敗；真實 PLC、正式 DB、Windows 及現場驗收未執行。
- 唯讀審查另指出既有通用項目刪除／失敗回復的並行修改風險，以及獨立記憶體儲存庫的關聯競態；原始 HEAD 已有相同通用流程，本次未擴大修改。記憶體關聯競態仍屬靜態推導，不能記為已重現或已修復。
- 本次 product scope 快照：decouple `4e617ee87dd1bdbb03cc36dc1ff2c98379095759`（208 檔）、clarify `2dcbd9e8021d8c289a0bb267800e8ec4355c1387`（28 檔），均為 touched tracking；既有髒檔不能逐 hunk 確認歸屬。另有 65 個必要依賴／修復檔以工作階段 SHA-256 清單補充；追蹤資料庫唯讀，未提權補登。
