# 本次 apply 交接（2026-11-06）

本頁是本次接續工作的最新證據；[既有紀錄](verification-evidence.md) 保留先前版本的成功、失敗與未執行狀態。schema 為 `spec-driven`，TDD／audit 開啟；本次完成 3.2、4.2、5.1、6.1，CLI 最終為 `all_done`、18／18。未 archive、commit、push 或部署。

## 決定與實作範圍

- 使用者確認：5000 筆是每設備的清理目標，全部進行中事件保留，即使超過目標；事件結束後再次清理。現有 repository 已符合，不改其行為，補預設 5000／5001 active、分頁、摘要與設備隔離的回歸測試，更新註解及 Swagger。
- 使用者確認：維持 MC `0000` 支援邊界，LTN0 明確回報不支援，不加入 `0002`。optimizer 拒絕長計時／計數目前值 LTN／LSTN／LCN；直接 word／exact-word read 在 transport 前拒絕它們。六種 long bit 保留既有非 random 規劃，不宣稱新增設備支援。
- 4.2 未改正式元件；新增隔離的 mock API fixture，實測錯誤重試、兩個首頁位置的三態、篩選、持續時間更新、Escape／焦點返回與三種寬度。
- 6.1 只用 loopback MC peer、測試 fixture 的新建 SQLite、正式服務 wiring／router、Vite 正式 App。沒有使用開發或正式資料庫，也沒有替換正常 operator 的 API 回應。
- `hc --skill` 確認沒有合適模板；三位 worker 透過互動式 `codex --no-daemon` 啟動，畫面確認 `gpt-6.1-terra`／`max`／`Fast`。主代理核對 source、diff、原始輸出及截圖，並負責 task 勾選；worker 未自行擴張任務或提交。

## 修改檔案

| 範圍 | 本次檔案 |
| --- | --- |
| MC 規劃 | `internal/datalink/optimizer/acquisition_mc_profile.go`、`acquisition_random_contract_test.go`、新增 `acquisition_random_long_timer_regression_test.go` |
| MC 直接讀取 | `internal/protocol/mcprotocol/client_word_read.go`、`client_word_read_test.go`；共用 guard 是阻止繞過 optimizer 送出 LTN `0401/0000` 的必要呼叫端 |
| 異常查詢／保留 | `internal/api/handlers/acquisition_anomaly_handler.go`、`internal/datalink/anomaly/sql_repo.go`、`types.go` 僅註解；新增 `internal/api/handlers/acquisition_anomaly_retention_regression_test.go`、`internal/datalink/anomaly/sql_repo_retention_regression_test.go` |
| API 文件 | `docs/swagger/docs.go`、`swagger.json`、`swagger.yaml`，以 `make gen-docs` 產生；另有生成器對同一 enum 的成對排序變動，未改其值或對應 |
| UI 驗證 | 新增 `frontend/tests/e2e/fixtures/recording-anomalies-host.html`、`.tsx`；本 change `evidence/` 內 `anomalies-20261106-*` 7 張與 `live-cli10-20261106-*` 8 張 PNG |
| 本 change 文件 | `proposal.md`、`design.md`、兩份 delta spec、`tasks.md`、本頁、`scenario-coverage.md` 及既有證據頁的最新入口連結 |

既有 staged WIP 保留。同期 `clarify-recording-editor-constraints` 的 tasks／驗收文件與圖片由其他工作更新，本次不歸屬、不回退；不將整個髒工作樹當作本 change 的新增修改。6.1 歸屬只列明沿用的三個驗收測試入口，不表示本次改寫它們。

## 修復前失敗與恢復證據

- optimizer：D10／D20／LTN0 在修復前沒有 LTN 錯誤，還產生 LTN contiguous request；同一新測試修復後通過。
- 直接 word／exact-word read：新測試在修復前觀察到一筆 LTN `0401/0000` 及 nil error；修復後兩個入口皆有明確錯誤、零次送出。
- 最終 review 撤回不必要的 long bit 拒絕：六種 long bit 的既有非 random 路徑先在過寬 guard 下失敗，再縮回三種目前值區域後通過。未擴充 bit、write、registry 或 `0002`。
- 異常保留原本已符合，因此採 regression workflow。以 Go overlay 移除 Open 清理、錯刪 active、移除 End 清理時，新測試分別失敗；未改共享 production source，恢復原編譯來源後通過。未把環境／cache 錯誤算作產品失敗。

## 本機串接量測

以下皆是本機單次量測，不是實體 PLC 或正式 DB 驗收。兩台 MC peer 共 12 點（9＋3）、原 1000ms；目標 SQLite 的寫入間隔為 1 秒。

| 階段 | 實測時間 | collected 增量 | SSE value（A／B） | SQLite 新增列 | MC requests（A／B） |
| --- | ---: | ---: | ---: | ---: | ---: |
| 可達 | 8.006 秒 | 96 | 72／24 | 8 | 8／8 |
| 連接器不可達 | 12.006 秒 | 144 | 108／36 | 0 | 12／12 |
| 恢復 | 8.003 秒 | 96 | 72／24 | 8 | 8／8 |
| B 改 3000ms | 12.007 秒 | 120 | 108／12 | 12 | 12／4 |

不可達以測試連接器指向 loopback 無服務埠產生，實際失敗／退避歷史可讀；恢復後驗證真 SQLite 列數。`buffer_accepted`、採集計數、SSE 與 SQL 新增列是不同觀察，不互相代替。API 更新 B 回 `applied`，runtime B configured/effective 為 3000／3000，A 保持 1000／1000。

另以原 32 點、1 秒週期的 production fixture 重跑不可達連接器量測：最終 22.001 秒增加 684 筆，31.090 筆／秒，高於設計條件 30.4 筆／秒，失敗歷史 5 次。mixed 規則 D0–D31 加 D39／D40／D100／D200／D1000／D2000，當次實際一個 0403、38 units／Good samples；safe-v1／未宣告／Modbus 的 mixed 仍不讀取。

暫停 B 的真 MC TCP 回應後，API 與首頁均觀察到 3 stale＋1 disconnected，A 仍為 0。UI 類型／狀態篩選與真 GET 對上；恢復後同 4 個 event ID 均已結束，摘要歸零。三筆 stale 持續 76001ms、斷線 67997ms，UI 顯示 1 分 16 秒／1 分 7 秒。此串接全部使用 D 位址；最後 guard 縮限只保留原 long bit 路徑，不改這批 D 位址流程，並另跑最終版本 Go 全套與 loopback 選測。

## 截圖

- Mock API 三寬：[1440](evidence/anomalies-20261106-1440.png)、[768](evidence/anomalies-20261106-768.png)、[768 右側欄](evidence/anomalies-20261106-768-scrolled.png)、[390](evidence/anomalies-20261106-390.png)。含進行中時間更新、固定已結束時間及可重試錯誤。
- 真 API：[首頁警告 4](evidence/live-cli10-20261106-active-home-1440.png)、[進行中事件](evidence/live-cli10-20261106-active-1440.png)、[stale 篩選](evidence/live-cli10-20261106-active-stale-1440.png)、[斷線篩選](evidence/live-cli10-20261106-active-disconnected-1440.png)、[已結束](evidence/live-cli10-20261106-ended-1440.png)、[390](evidence/live-cli10-20261106-ended-390.png)、[390 清單內捲動](evidence/live-cli10-20261106-ended-390-scrolled.png)、[恢復首頁](evidence/live-cli10-20261106-recovered-home-1440.png)。

主代理已檢視三寬 mock 與真 API 1440／390 截圖。768 表格在 modal 內橫向捲動，390 堆疊欄位且清單可內捲；文件未橫向溢出。真 App 的 390 畫面有測試環境浮動控制，來源未查證、未遮住 Close；本輪未擴大修改它。

## 檢查與限制

環境為 macOS arm64、Go 1.27.3；程式碼風格依 go.mod 的 Go 1.25.5 指引。原始 logs／worker reports／scope hashes 位於工作階段暫存 `/tmp/decouple-20261106/`，未將完整暫存目錄或其他 change 納入 repo。

- PASS：最終 `go test ./... -count=1`，`go-full-final.log`。
- PASS：`cd frontend && npm test -- --run`，633 檔、7277 項，`frontend-full.log`；本次前端正式 source 未變。
- PASS：`cd frontend && npx tsc --noEmit`；新增 UI fixture 另以隔離 tsconfig 型別檢查及 ESLint 通過。
- PASS：異常模組／handlers 的相關 race，MC／optimizer／adapter／collector 的最終相關 race 與 focused suites；最終相關 `go vet`。
- PASS：`make gen-docs`；規格 `spectra analyze` 與 `spectra validate`。
- PASS：最終版本的三項 loopback 驗收（`GATEWAY_TASK13_THROUGHPUT=1 GATEWAY_MIXED_SIMULATOR=1 go test ./cmd/test_ui -run '^(TestProductionSnapshotWorker_MCUnreachableThroughput|TestProductionScheduledAcquisitionMixedRuleWireRequests|TestProductionScheduledUnsupportedMixedRuleIsNotRead)$' -count=1 -v`）。
- PASS：`bash scripts/check_file_lines.sh`、staged／unstaged 的 `git diff --check`。現有歷史超長檔依規則未增行；本輪文件均未超過 300 行。Swagger JSON 比對確認只有要求的 description 語意變更與等價 enum 成對排序。
- 兩次真瀏覽器工作階段已關閉；本次 Vite 已停止。Go acceptance host 結束為 PASS，fixture 的 SQLite／MC peers 隨清理關閉，不留對正式資料的修改。
- 真實 PLC、正式 PostgreSQL／MySQL／SQL Server、Windows、嵌入式發行包、部署與現場驗收：NOT RUN，不在本次本機驗收主張內。

情境覆蓋分類見 [scenario-coverage.md](scenario-coverage.md)：23／24 個 Scenario 與 2／2 個 Example 有對應測試；有效上限 120 的精確請求拆分仍為非阻擋部分覆蓋，未因此新增功能或完成條件。

最終 log SHA-256（完整記錄留於上述暫存目錄）：

| 紀錄 | SHA-256 |
| --- | --- |
| `go-full-final.log` | `ba4b133e385d673edcc8d8b6ec0111a7d5384983dd420c138046b0d4f588e892` |
| `frontend-full.log` | `2c75af17c9d387ef308407d672e064f6073c9aaa32315464778b84bc058d154f` |
| `loopback-final.log` | `bf95f27bea39389b9f9c9b309ccea60f8e6e46e2f7371996272552ad7612fb27` |
| `live-acceptance.log` | `f0672560274ebf16343d9b5cb75080e0d9d6adc57db0d17df6792d6d4c7c3433` |
| `acceptance-host.log` | `1c6dafecea4009ac1bbddfbb107be8f7f00c63abdcef26f7ac3231b7e7fc7c30` |

index 的 path／mode／blob／stage 清單與開始時逐位元組相同；本次未 stage 任何檔案。後續先執行 `spectra-verify decouple-acquisition-from-delivery` 與 `spectra-review decouple-acquisition-from-delivery`，再決定 archive。
