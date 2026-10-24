# 迴路記錄實際瀏覽器驗證

對應 `add-circuit-recording-profiles` task 5.2，2026-10-24。
本頁是模擬來源與可丟棄資料庫的實跑紀錄，不是現場設備或 Windows 驗收。

## 環境與證據

- 正常 `cmd/test_ui --headless`，獨立 SQLite 設定／journal；另用獨立 SQLite 檔接收記錄。
- 真實 loopback Modbus TCP simulator，8 個映射來源：4 個 uint16、4 個 uint64。
- 專用 agent-browser Chromium headless shell；前端 Vite、API 與 simulator 都只監聽 `127.0.0.1`。
- 繁中由專用瀏覽器的 `i18nextLng=zh-TW` 偏好設定載入。沒有將設定頁下拉選單當成語言切換成功的證據。
- 實測 1600×1000 與 1280×800；兩個尺寸均未出現整頁水平溢出。
- [SQL 讀回、群組／record 身分與截圖雜湊](evidence/circuit-recording-browser-20261024.json)。
- [資料庫三模式矩陣](studio-v2-circuit-recording-db-verification.md)。
- [規格情境對照](evidence/circuit-recording-scenario-coverage.md)與[逐測試結果](evidence/circuit-recording-test-results-20261024.json)。
- [收尾 review 修正與後續 Go 重測](studio-v2-circuit-recording-final-review.md)；本頁保留實際瀏覽器操作當時的證據。

## 實際操作結果

| 操作 | 結果與可觀察證據 |
|---|---|
| 四步設定 | 從真實介面建立設備、連線測試、兩組讀取規則及8個映射；設定獨立 SQLite 目的地。 |
| 迴路A雙型態 | 瞬時只選40001，15秒；累積只選40015，1分鐘。透過既有介面明確確認 counter 與 kWh。 |
| 週期界線 | 輸入0顯示整數／範圍錯誤；恢復15後才繼續。 |
| 原生鍵盤 | 聚焦 checkbox 後按 Space，實讀 `checked=false`；再按 Space，實讀 `checked=true`，焦點仍在同一個 checkbox。 |
| 資料表準備 | preview 後目的地仍無表；分別明確確認兩個群組後，才建立共享長表與 receipts。 |
| 開始記錄 | 真實 recording-start operation 為 succeeded，兩個群組均 applied；實際採集→journal→outbox→sender→SQLite。 |
| 累積精度 | 首筆 raw=`9007199254740993`、delta=NULL；調整 simulator 一個 register 後，raw=`9007199254740994`、delta=`1`。 |
| 迴路B | 另選40002與40019；各自確認、套用並開始。與A共用長表，保有不同 canonical 群組ID。 |
| 批次修改 | 只選B的瞬時型態，先顯示15→30秒差異，再明確儲存。B瞬時ID／成員保留，revision改變；A兩組及B累積的完整已存物件逐項比較完全不變。 |
| reload | 再載入仍為B、30秒、原群組ID，沒有透過 ensure 新建替代 writer。 |
| 準備遷移目標 | 從舊寬表面板進入 prepare-only，明選新迴路、來源與週期；目標確認後仍未 applied，且畫面沒有 Start 按鈕。 |
| 切換 | 介面審核來源／目標／15秒週期／表名／來源清單，以原 digest 確認。精確邊界為 `2026-10-24T15:26:07Z`。 |
| 邊界SQL | 舊寬表保留8筆歷史與最後截斷窗；新長表首筆 `partial`，實際窗為15:26:07–15:26:15，raw=`32`。 |
| 回復 | 從同一筆切換紀錄另行審核新邊界15:28:11；long最後一窗為15:28:00–15:28:11、partial，wide新版本從15:28:11恢復。兩張表及先前歷史都保留。 |
| 缺表恢復 | 保存SQL證據後，只刪除本輪隔離目的地的受管長表；保留 receipts。下一次真實交付自動重建，generation=2、reconstructions=1。 |
| 歷史警告 | 真實 delivery 回傳 history_may_be_lost=true；繁中畫面明說結構重建與 receipts 不保證歷史復原。 |
| 重建後累積 | 再把來源加1，SQL讀回 raw=`9007199254740995`、delta=`1`，保留原基準與真實 delta_from／delta_to，沒有重設為首筆 baseline。 |

舊寬表是為遷移情境而透過正式 API 建立、確認與啟動的前置資料；切換與回復使用真實瀏覽器面板。
第一次前置資料誤設 optional member＋skip_row，依既有契約改為 required後才進行驗收，不計為產品缺陷。
CLI無法正確填入原生分段日期控制項，改用同一 input 的原生 value setter及 input/change事件；審核與確認仍走實際介面，沒有攔截或偽造回應。

## 實測發現與必要修正

每項修正都有先觀察到的行為失敗；正向對照及後補案例另列，沒有把全部測試都稱為先失敗。

1. SQLite schema在後端正規化成 `main`，前端空值比較卻一直判為dirty。
   改為僅在SQLite比較空值與main的同義性，保留PostgreSQL的嚴格比較。
2. 新型態啟動把資料庫setup revision當成Share workspace revision。
   分開兩種版本；由既有宿主傳遞完整Share證明，不用setup revision代填。
3. 長表partial誤走寬表每成員欄位檢查，已建表仍不能開始。
   長表沿用其固定布局／owner／主鍵／索引檢查，寬表partial驗證原樣保留。
4. 首次連線保存與query hydration競態可能把saved覆蓋回saving。
   僅在相同回應ID與非空identity revision已被hydrate時承認完成；真正切換目的地或較新編輯仍保留。
   此交錯由回歸測試重現；沒有宣稱觀察到瀏覽器當下的確切排程順序。
5. 缺表重建旗標已有API契約，但新型態畫面沒有歷史遺失提示。
   補上群組身分核對後的世代與明確警告；legacy缺欄及其他群組回應不顯示。

新增回歸：`circuit-profile-browser-wiring.test.tsx`、`database-connector-save-race.test.ts`、
`circuit-profile-storage-notice.test.tsx`、`TestCircuitProfilesPartialPolicyUsesTheLongLayout`。另補 `TestManagedWidePolicyEditRequiresNoChangeConfirmation`，實跑SQLite／PG，確認policy-only修改仍需新revision的空DDL明確確認；此案例為既有行為回歸，沒有宣稱先紅。

## 檢查結果

- 前端：完整260個測試檔／1812測試 PASS、0 FAIL、0 SKIP；完整ESLint、TypeScript及Vite build PASS。
- Go：workspace、measurement、recordingplan、dbtarget、runtime、grouppipeline、groupdelivery、snapshot、aggregation、cmd/test_ui及internal/api，共11個套件 PASS。
- API明確排除已確認的無關基線 `TestStudioV2RecordingRoutes_PostgresApplyCreatesAndVerifiesTheManagedSchema`；不將它記為PASS。
- Go四個SKIP：兩個crash-child helper；SQLite不支援只差大小寫的識別字案例；無Timescale擴充案例在本次有擴充fixture跳過，另以task5.2 PG18無擴充fixture實跑同一API測試PASS；該phase五個需要擴充的子案例跳過，均有TS全套回歸的對應PASS。
- `go vet` PASS；golangci-lint只剩無關WIP `internal/api/router_logger.go:109` 的gocritic，未宣稱整體lint全綠。
- Swagger以 `make gen-docs`重產；inventory由正式changelog CLI同步。
- `git diff --check`及行數檢查 PASS。本輪未執行stage或重設index；收尾觀察到其他操作更新暫存區，已保留現狀。

## 截圖

初次7張截圖在收尾連結檢查時已不在原路徑，原因未確認；實際操作與SQL證據未遺失。
下列4張為同一份實際已保存資料的補存畫面，已獨立位元組備份，並在關閉瀏覽器與測試服務後再次核對檔案及SHA-256。

- [只準備目標的明確入口與既存設定](evidence/circuit-recording-20261024/final-prepare-only.png)
- [只選B的唯讀預覽](evidence/circuit-recording-20261024/final-batch-preview.png)：補存畫面為30→45秒，拍攝後取消，沒有儲存；原先15→30秒的套用斷言見JSON。
- [forward與rollback的真實已記錄歷史](evidence/circuit-recording-20261024/final-cutover-history.png)：非尚未套用的審核畫面。
- [重建第2代、歷史警告與SQL receipt分階段證據](evidence/circuit-recording-20261024/final-reconstruction-warning.png)

## 限制

瀏覽器此次目的地為SQLite；PG17／PG18／Timescale的真實SQL與catalog證據由5.1及後端整合套件提供，沒有宣稱各資料庫都跑過瀏覽器。
Runtime串流曾顯示degraded及停止文案，但後端採集與SQL已持續產生資料；沒有把串流畫面當成SQL成功證據，也未擴大修復該監看surface。
Step 3的數值預覽曾將超過2^53的整數顯示為相鄰值；這不影響本次已讀回的精確raw／delta，亦未宣稱該預覽的顯示精度已修正。
沒有使用正式DB／PLC／LAN，沒有Windows實機驗收，沒有commit、archive、push或部署。
