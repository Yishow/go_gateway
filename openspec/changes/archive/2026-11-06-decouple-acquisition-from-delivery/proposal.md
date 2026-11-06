## Why

PLC 採集與資料庫寫入目前在同一個 ingest goroutine 內同步執行：`runtime.consumeLoop` 逐筆呼叫 `handleCollectedValue`，其中舊式資料庫輸出（`dbtarget.Writer.WriteTagValue` 的未分組 mapping）每筆同步連線寫入。當資料庫連不上時，每筆等滿 5 秒 context 逾時，整條管線被拖垮。2026-11-04 在開發機實測（量測值）：設備端每秒一次 `0401` 讀取 D0 起 32 word，但 `collected_total` 只增加約 1.8 筆/秒；以本機 MC 模擬器加資料庫副本量測，點位狀態寫入約 0.2–1ms，資料庫輸出每 10 秒約 2 筆各卡 5.0 秒（`資料庫連接測試失敗: context deadline exceeded`），1000 筆 value 緩衝滿後丟棄最舊值，造成 SSE 即時值落後約 20 秒、部分點位長時間不更新，Modbus Share 也一併被拖慢。同時，操作者看不到斷線、長時間沒更新或數值凍結等採集異常的紀錄，也無法在連線設定中為單一設備調整取值間隔。

## What Changes

- 採集與資料庫寫入分離：ingest 路徑只更新點位的最新暫存器快照、SSE 與 Modbus Share；舊式資料庫輸出改由每個資料庫連接器各自的背景 worker 依其寫入間隔取「當下快照」寫入。連不上的資料庫只影響自己的 worker，不再阻塞採集。
- 連接器連線失敗時，worker 依退避節奏重試並如實記錄失敗；不以假成功補位。
- 新增採集異常紀錄：持久記錄三種事件——設備斷線、點位長時間沒更新（stale）、點位長時間同一個值（凍結）——含起訖時間、設備、點位與持續時間；門檻可設定，以每設備 5000 筆為歷史清理目標，優先刪除最舊已結束事件；進行中事件全部保留，即使超過目標。
- 首頁在每台設備的「讀取點位」按鈕旁新增「採集異常」按鈕；該設備有進行中的異常時按鈕改用警告色並附文字計數。點擊開啟依 DESIGN.md 規範的置中 modal，列出異常事件，可依類型與狀態（進行中／已結束）篩選。
- 連線設定新增「取值間隔」：每台設備擁有自己的輪詢群組，儲存後該設備點位改用此間隔；程式啟動時依各設備間隔自動開始採集。
- 新增「自動調整間隔」選項：當一輪讀取耗時超過設定間隔時，實際間隔自動拉長以避免請求堆積，設備恢復後回到設定值；畫面同時顯示設定值與實際值。
- 連線設定新增「取值優化設定檔」，沿用現有 `acquisition` 宣告：Modbus 沒有隨機讀取指令，只做連續區塊聚合；MC 3E 與 FATEK 在宣告相容設定檔時才啟用 random 讀取，optimizer 會先保留長的連續區塊、再把零散位址打包成 random 請求。
- MC 3E random 讀取（`0403`）補齊 double word 與 bit 裝置（bit 以 16 點為一個 word 單位讀取）；新增 MELSEC-Q／iQ-R 設定檔。依官方手冊 SH-080008，Q/L 與 iQ-R 使用子指令 `0000` 時 word 點數＋double word 點數合計上限 192，FX5 亦為 192；FATEK 指令 48 上限為 64 個元件且以 word 計不超過 64（32 bit 算 2）。
- random 上限顯示於連線設定，預設為協議上限，可改小但不可超過協議上限。
- 新增「分散位址」讀取範圍：可隨機讀取的設備可在同一個範圍內加入不同區域、不同型別的位址；用量達上限 90% 時提醒，超過上限時說明會自動拆成多次請求。

## Non-Goals (optional)

詳見 design.md 的 Goals / Non-Goals。本次不新增 iQ-R 子指令 `0002` 或 LTN0 支援；LTN0 回報明確不支援原因。

## Capabilities

### New Capabilities

- `acquisition-anomaly-log`: 採集異常（斷線、stale、凍結）的偵測、持久紀錄、查詢 API，以及首頁按鈕與檢視 modal。
- `device-acquisition-interval`: 每台設備的取值間隔、自動調整間隔、啟動即採集，以及在連線設定中選擇依 protocol 的取值優化設定檔。
- `random-acquisition-range`: MC 3E random 讀取支援 word、double word 與 bit（word 單位）、MELSEC-Q／iQ-R 設定檔、可調 random 上限，以及可在一個範圍內選取多個分散位址並在用量達 90% 時提醒的 UI。

### Modified Capabilities

- `database-output-delivery`: 新增「資料庫輸出不得阻塞採集」要求——舊式資料庫輸出改為每連接器背景 worker 依間隔寫入當下快照，失敗退避並如實記錄。

## Impact

- 協議與規劃：internal/protocol/mcprotocol/client.go（0403 double word 與 bit 裝置）、internal/datalink/connector/acquisition_profile.go 與 internal/datalink/connector/acquisition_profile_decode.go（新設定檔與 random 上限欄位）、internal/datalink/optimizer（double word 與 bit 的 random 打包成本）。
- 後端：internal/datalink/runtime/ingestor.go、internal/datalink/runtime/service.go、cmd/test_ui/target_writer.go、cmd/test_ui/service_wiring.go、internal/datalink/dbtarget/writer.go、internal/datalink/point/service_polling_group.go、internal/datalink/collector（排程間隔與自動調整）、新的異常偵測與儲存模組（internal/datalink 之下）。
- 資料庫：新增 migration 於 internal/datalink/schema/migrations/（異常事件資料表、設備與輪詢群組關聯）。
- API：新增異常事件查詢端點；裝置連線設定接受取值間隔、自動調整與 acquisition 設定檔欄位；需以 make gen-docs 重新產生 Swagger。
- 前端：frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx 與首頁讀取範圍表單（分散位址範圍）、frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeCollection.tsx（按鈕）、新的異常 modal 元件、frontend/src/features/datalink/workbench-v2/recording-editor 的連線設定表單、frontend/src/services 與 frontend/src/types、en／zh-TW locales。
- 行為：舊式資料庫輸出從「每讀到一次就寫一列」改為「每個寫入間隔（預設 15 秒）寫一列當下最新值」；這類資料表的列數會變少，間隔中途的瞬間變化不會寫入，需在交接說明。
