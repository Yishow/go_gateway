## 1. 採集與資料庫寫入分離

- [x] 1.1 依「舊式資料庫輸出改為每連接器快照 worker」實作 Database output never blocks acquisition：ingest 只更新以 mapping 為鍵的最新值表並立即返回，每個啟用連接器一個背景 worker 依有效寫入間隔寫入當下快照（已分組 mapping 的 flush 也移入該 worker）。先寫失敗測試（TDD：以永不回應的連接器假實作，斷言 32 筆值的 ingest 每筆 < 50ms 且 SSE 事件數等於輸入數；間隔內 10、11、12 只寫入一列 12），再實作讓 `go test ./internal/datalink/dbtarget/ ./internal/datalink/runtime/` 通過。
- [x] 1.2 依「連接器 worker 失敗退避」實作退避 1、2、4、8 秒、上限 60 秒，退避期間不排隊、恢復後只寫當下快照，每次失敗經既有寫入結果與寫入歷史記錄；以注入時鐘的 Go 測試驗證退避序列與「間隔內無新值不寫入並記錄略過原因」。 [after: 1.1]
- [x] 1.3 把 `cmd/test_ui` 的 production wiring 改為使用快照 worker、移除 ingest 路徑上的同步舊式寫入，Modbus Share 寫入行為不變；以 `go test ./cmd/test_ui/` 及本機 MC 模擬器加不可達連接器重跑量測，記錄 `collected_total` 每秒增長率（目標 ≥ 32 筆/秒的 95%，設計條件）於交接證據。 [after: 1.2]

## 2. 設備取值間隔

- [x] 2.1 依「設備專屬輪詢群組」實作 Each device owns its acquisition interval 的資料層：新增 migration 於 `internal/datalink/schema/migrations/`（`polling_groups.owner_device_id` 唯一可空、設備門檻欄位），啟動時冪等地為每台設備建立專屬群組（間隔沿用其點位原群組）並移動點位；Go 測試驗證升級後採集間隔不變、重跑不重複建立。
- [x] 2.2 設備更新 API 接受 `acquisition_interval_ms`（100–3600000）、`acquisition_auto_interval`、`stale_threshold_ms`、`frozen_threshold_ms`，儲存時只重新指派該設備點位並即時生效；source rule 產生點位與 runtime 裝置同步改為優先使用設備專屬群組。Go 測試：A 改 500ms、B 仍 1000ms；非法值回欄位錯誤。以 `make gen-docs` 更新 Swagger。 [after: 2.1]
- [x] 2.3 依「自動調整間隔在 scheduler 計算有效間隔」實作 Automatic interval adjustment：有效間隔 = max(設定值, 最近 5 輪最大耗時 × 1.2 取 100ms 倍數)，連續 5 輪準時逐步縮回；runtime status 回報 `configured_interval_ms`、`effective_interval_ms`。先寫失敗測試（TDD，注入掃描耗時序列：1250ms → 有效 1500ms，之後 5 輪準時回到 1000ms）再實作。 [after: 2.1]

## 3. 採集異常紀錄

- [x] 3.1 依「異常偵測在 runtime 內以週期掃描最新快照」與「異常門檻依設備設定」實作 Acquisition anomalies are detected and recorded：新增 `acquisition_anomalies` 資料表與儲存庫，runtime 以 1 秒週期掃描最新點位狀態與 collector 連線／breaker 事件，產生 stale、frozen（預設停用）、disconnected 事件的開始與結束，同類同點位只有一筆進行中。先寫失敗測試（TDD，注入時鐘：D5 12 秒無讀取開啟 stale、成功讀取結束；30 秒斷線只有一筆 disconnected）。 [after: 2.1]
- [x] 3.2 實作 Anomaly history is bounded and queryable：每設備以 5000 筆為清理目標、先刪最舊已結束事件，所有進行中事件均保留（即使超過目標），事件結束時再次清理；新增 `GET /api/v1/datalink/devices/:id/acquisition-anomalies`（state、kind、limit、offset、total）與 `GET /api/v1/datalink/acquisition-anomalies/summary`。Go handler 測試驗證篩選、分頁、保留進行中事件；`make gen-docs` 更新 Swagger。Go 測試另驗證 5001 筆進行中事件全部保留、其中一筆結束後清至 5000 筆。 [after: 3.1]

## 4. 前端

- [x] 4.1 在 `frontend/src/services` 與 `frontend/src/types` 新增異常列表與摘要的 service／型別，以既有 React Query hooks 讀取；vitest 驗證回應解析與錯誤處理（失敗不轉成空列表）。 [after: 3.2]
- [x] 4.2 實作 Home exposes anomalies beside the device points action：`RecordingHomeCollection` 的設備列與設備摘要在「讀取點位」旁加入「採集異常」按鈕，三態（無異常一般次要按鈕、有異常 warning tokens＋文字數量、摘要失敗顯示狀態未知）；依 DESIGN.md 新增 form 分類置中 modal，列出事件並提供類型／狀態篩選、開始、結束、持續時間、位址、最後數值，錯誤時提供重試。vitest 驗證三態、篩選與錯誤；以 1440/768/390 寬度實際截圖並附於交接。修復持續時間後的畫面尚未重驗：沙箱內瀏覽器無法啟動，先前截圖僅作歷史證據。 [after: 4.1]
- [x] 4.3 依「取值優化設定檔沿用既有 acquisition 宣告」實作 Connection settings expose protocol-appropriate acquisition optimization：連線設定表單新增取值間隔、自動調整、stale 門檻、凍結門檻、設定檔下拉與 random 上限（預設協議上限、可改小）（Modbus 只顯示連續聚合說明；MC 3E：未宣告／safe-v1／mitsubishi-fx5-cpu-3e-binary-v1／mitsubishi-q-iqr-3e-binary-v1；FATEK：未宣告／safe-v1／fatek-fbs-ascii-v1），每欄附 ⓘ 說明，並顯示 runtime 的設定值與有效值；en／zh-TW 文字。vitest 驗證各 protocol 的選項、欄位驗證與儲存 payload；實際截圖附於交接。 [after: 2.2, 2.3, 5.2]

## 5. 分散位址 random 讀取

- [x] 5.1 依「MC random 補齊 double word 與 bit，並新增 Q／iQ-R 設定檔」實作 MC random read covers word, double word and bit devices：`MCClient.RandomRead` 送出 word 與 double word 兩段、允許 16 點對齊的 bit 裝置並解出指定 bit、拒絕長計時器／長計數器與字串；optimizer 依 32 bit 型別放入 double word 段並以合計點數計成本。先寫失敗測試（TDD：D39 int16、D100 int32、M5 bool 產生一筆 0403，位元組與手冊 SH-080008 範例格式一致；LTN0 明確回報不支援且不送出請求），再實作讓 `go test ./internal/protocol/mcprotocol/ ./internal/datalink/optimizer/` 通過。依使用者決定維持 `0000` 支援邊界；以 D10、D20、LTN0 驗證明確不支援原因，不新增 `0002` 或部分成功契約。
- [x] 5.2 依「random 上限可調但不可放寬」實作 Random profiles carry documented limits：設定檔目錄新增 `mitsubishi-q-iqr-3e-binary-v1`（192，排除 ZR），`acquisition.limits.random_points` 只能小於或等於設定檔上限，FATEK 維持 64 元件／64 word 加權；Go 測試涵蓋 120 接受、300 拒絕、FATEK 48 接受，且 registry schema 與驗證同源。 [after: 5.1]
- [x] 5.3 依「分散位址讀取範圍需補齊 mixed 規則的後端持久化」實作 Scattered mixed rules are persisted and acquired（持久化與展開部分）：source rule repository 的 create／update 與 Studio v2 工作區規則 API 讀寫 `layout_mode` 與 `mixed_items`（沿用 `internal/datalink/measurement` 的 `MixedItem`，每項自帶 `address` 與 `data_type`，可跨區域），儲存前以 `ValidateMixedLayout` 加上逐項位址解析驗證並回傳欄位錯誤（含重疊項目與位址）；mixed 規則依項目展開為點位與規則連結，每項一個點位且保留項目識別，homogeneous 規則行為不變。先寫失敗測試（TDD：MC 設備儲存 D39 int16、D100 int32、M5 bool 三項後重新讀回內容一致並產生 3 個點位；D3 uint32 與 D4 int16 重疊時拒絕並指出兩項），再實作讓 `go test ./internal/datalink/sourcerule/ ./internal/api/handlers/` 通過；`make gen-docs` 更新 Swagger。 [after: 5.2]
- [x] 5.4 依「分散位址讀取範圍需補齊 mixed 規則的後端持久化」實作 Scattered mixed rules are persisted and acquired（採集部分）：`acquisitionRuleExclusions` 對宣告 random 設定檔（`mitsubishi-fx5-cpu-3e-binary-v1`、`mitsubishi-q-iqr-3e-binary-v1`、`fatek-fbs-ascii-v1`）的設備接受 mixed 規則並把各項目交給 optimizer 規劃（未宣告、`safe-v1` 與 Modbus 仍回報 `source-layout-unsupported`），runtime projection、readiness 與 Modbus Share 保護沿用既有路徑。先寫失敗測試（TDD：q-iqr 設備的 mixed 規則 D39、D40、D100、D200、D1000、D2000 加既有 D0–D31 homogeneous 規則，計畫請求數 ≤ 2；safe-v1 設備的同一規則仍被排除並附原因），再實作讓 `go test ./internal/datalink/sourcerule/ ./internal/datalink/optimizer/ ./internal/datalink/workspace/` 通過；以本機 MC 模擬器確認實際請求數 ≤ 2，無法執行時如實記錄。 [after: 5.3]
- [x] 5.5 依「分散位址讀取範圍沿用 mixed 版面」實作 Scattered address range with usage warning：讀取範圍表單在設備設定檔支援 random 時提供「分散位址」模式，經 5.3 的 API 儲存為 `layout_mode: mixed` 與 `mixed_items`；即時顯示 random 用量對上限（MC 16／32 bit 各 1 點、bit 以 16 點單位計；FATEK 32 bit 算 2），≥ 90% 顯示提醒、超過時說明會拆成多次請求但不阻擋儲存，不支援 random 的設備顯示原因；有效上限重用 `frontend/src/features/datalink/workbench-v2/state/acquisitionSettings.ts`；en／zh-TW 文字。vitest 驗證 173／192 出現提醒、用量計算、停用原因與儲存 payload；實際截圖附於交接。 [after: 5.4, 4.3]

## 6. 驗收

- [x] 6.1 端到端驗收：以本機 MC 模擬器與可丟棄資料庫，確認連接器不可達時採集與 SSE 不受影響、恢復後依間隔寫入、首頁異常按鈕與 modal 顯示斷線與 stale 事件、修改單一設備間隔只影響該設備、分散位址範圍以 random 讀取；執行 `go test ./...`、`npm test -- --run`、`npx tsc --noEmit`、`bash scripts/check_file_lines.sh`，結果與未執行項目記入交接。最新審查修復後 Go 全套通過；曾於十分鐘切換時有一項範圍外測試失敗，單項重跑三次通過；前端全套有一項既有不穩定測試失敗，畫面重驗未執行。 [after: 1.3, 4.2, 5.5]

## 7. 審查後補償

- [x] 7.1 依「既有多間隔設備升級取最短間隔」補償已完成的 2.1：同一設備原有啟用 500／2000 ms 群組時遷移至一個 500 ms 專屬群組；停用群組所屬及未指派點位與其他設備不變，停用點位不被啟用，重跑不重複。SQLite migration 測試在修正前原始快照重現失敗，修正後 `go test ./internal/datalink/ -run DeviceAcquisitionBackfill` 通過。 [after: 2.1]

## 本次接續證據

2026-11-06 本輪完成 3.2、4.2、5.1、6.1，最終 18／18。上方保留已完成任務的歷史描述；其中前次前端失敗／畫面未重驗的狀態已由本輪前端全套通過、三寬截圖與真 API 串接證據取代。最新結果、修改範圍及未執行項目見 [本次 apply 交接](apply-completion-evidence.md)，非阻擋的情境測試細節見 [測試對照](scenario-coverage.md)。未 archive 或提交。

## 8. 本次審查修復

- [x] 8.1 依使用者「修」的授權處理 review R1–R12：保留兩次掃描間結束的 frozen 事件；採集欄位、群組與點位原子儲存並防止部分更新／跨 API runtime 影像互相覆蓋；關閉時等 target worker 真正完成後才釋放 DB；寫入結果使用實際時間及安全失敗分類；修正 FATEK mixed 映射識別與定位、超過64項編輯、正規化 ACK、重疊與明確拒絕後的可編輯狀態；修正兩個 anomaly Swagger 路徑。各項先重現再修正，完成相關 Go／前端回歸、race、vet、型別、修改檔案 lint、build、行數與模擬 API 瀏覽器檢查，結果及未執行項目記於 [審查修復交接](review-fixes-evidence.md)。原18項完成紀錄保留；不擴充協議或功能需求。
