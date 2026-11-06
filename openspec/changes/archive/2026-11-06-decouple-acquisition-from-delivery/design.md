## Context

目前資料流：collector 排程依輪詢群組間隔產生 `CollectedValue`，經 1000 筆緩衝的 `valueChan`（滿時丟棄最舊值）送到 `runtime.consumeLoop`。consumeLoop 是單一 goroutine，逐筆以 5 秒 context 呼叫 `handleCollectedValue`，依序執行：點位狀態持久化、storage writer、`fanoutTargetWriter`（舊式資料庫輸出 `dbtarget.Writer.WriteTagValue` 與 Modbus Share），最後才廣播 SSE value 事件。寫入群組（`grouppipeline`）透過 `SampleSink.AcceptSample` 另走持久化管線，不在本次範圍。

`dbtarget.Writer` 對「已分組」mapping 已有緩衝與背景 `flushLoop`；「未分組」mapping 則每收到一筆就同步連線寫入，寫入間隔（mapping `write_interval_seconds` → 連接器 `default_write_interval_seconds` → 15 秒）只記入歷史。2026-11-04 量測（開發機，本機 MC 模擬器 + 資料庫副本）：未分組 mapping 指向連不上的 PostgreSQL 時，每 10 秒約 2 筆各耗 5.0 秒，ingest 吞吐約 3 筆/秒；正式開發機 `collected_total` 約 1.8 筆/秒，而設備端每秒完成一次 32 word 讀取（使用者提供的模擬器通訊日誌）。

所有點位共用由 `point.Service.EnsureDefaultPollingGroup` 建立的 `Auto Polling` 群組（1000ms），source rule 產生點位與 runtime 裝置同步都指派到此群組。取值優化已存在：`internal/datalink/optimizer` 依 protocol 做連續區塊聚合；MC 3E 與 FATEK 的 random 讀取僅在裝置 `connection_config.acquisition` 宣告相容設定檔時啟用（`collection-optimizer` spec 的 Capability-aware random acquisition）。

首頁（`RecordingHomeCollection`）在設備列與設備摘要各有「連線設定」「讀取點位」按鈕。UI 規範以 DESIGN.md 為 authority：置中 modal、75vw × 75dvh、語意 tokens、狀態需伴隨文字。

## Goals / Non-Goals

**Goals:**

- 任何資料庫輸出（含連不上的連接器）都不得延遲或丟失 ingest 路徑上的點位快照、SSE 事件與 Modbus Share 寫入。
- 舊式資料庫輸出依其有效寫入間隔寫入「當下快照」，每個連接器獨立 worker、失敗退避、結果如實記錄。
- 持久記錄斷線、stale、凍結三類採集異常，提供查詢 API 與首頁 modal。
- 每台設備可在連線設定調整取值間隔、啟用自動調整、選擇 acquisition 設定檔；程式啟動即依各設備設定採集。

**Non-Goals:**

- 不改寫入群組（`grouppipeline`／`groupdelivery`）的持久化與交付語意。
- 不推測未有手冊依據的機種支援 random：只開放 FX5（iQ-F）、MELSEC-Q 與 iQ-R（子指令 `0000`，上限 192）；L 系列、iQ-R 子指令 `0002`（上限 96）及 QnA／A 系列不在本次範圍。Modbus 沒有隨機讀取指令，維持連續區塊聚合。
- MC random 不支援長計時器／長計數器（LTS、LTC、LSTS、LSTC、LCS、LCC，手冊明列不可指定），也不支援 byte 單位（MC 無 byte 單位）與字串（改用連續讀取）。
- 不做異常通知推播（email、webhook）；只記錄與檢視。
- 不移除 `Auto Polling` 群組或其 API；它仍是沒有設備專屬群組時的後備。
- 不在本次做跨重啟補寫（資料庫斷線期間的快照不補寫；寫入群組才提供持久補送）。

## Decisions

### 舊式資料庫輸出改為每連接器快照 worker

ingest 路徑把每個值寫入記憶體中的「最新值表」（以 mapping 為鍵，覆寫舊值，不排隊）後立即返回；每個啟用的連接器有一個背景 worker，依每個 mapping 的有效寫入間隔，取最新值表的當下快照寫入。已分組 mapping 沿用既有分組緩衝語意，但 flush 也移到該連接器的 worker。替代方案：(a) 只對連線失敗做斷路器——仍在 ingest 路徑同步寫入，資料庫慢但沒斷時依然拖慢採集，不符合「採集與寫入分開」；(b) 每筆值排入有界佇列——佇列滿時仍需丟值且寫入量不變，與「只存當下暫存器」的需求不同。選擇最新值表因為它的記憶體上限等於 mapping 數，與讀取頻率無關。

### 連接器 worker 失敗退避

連線或寫入失敗時，該連接器 worker 依 1、2、4、8 秒倍增、上限 60 秒（設計條件）退避重試；退避期間到期的快照不排隊，下一次嘗試時只寫當下快照。每次失敗以既有 `recordWriteOutcome` 與寫入歷史記錄，連接器狀態可見。替代方案：固定間隔重試——連不上時每個 mapping 都打一次逾時，浪費連線資源。

### 異常偵測在 runtime 內以週期掃描最新快照

runtime 維護每個點位的「最後成功讀取時間」與「最後數值變化時間」，以 1 秒（設計條件）週期掃描並與門檻比較；設備斷線則來自 collector 既有的連線失敗與 breaker 狀態事件。事件以開始／結束兩階段寫入新資料表 `acquisition_anomalies`，同一點位同一類型同時只有一筆進行中事件。替代方案：在每筆 ingest 時判斷——stale（「沒有值進來」）無法由收到的值觸發，仍需週期掃描。

**Supersedes**: decouple-acquisition-from-delivery / 異常偵測在 runtime 內以週期掃描最新快照（原每設備嚴格 5000 筆上限）。2026-11-06 使用者決定以 5000 筆為清理目標，先刪最舊已結束事件；所有進行中事件均保留，即使它們已超過 5000 筆。事件結束時再次清理。替代方案為滿額拒絕新事件，因會遺失新的異常而未採用。

### 異常門檻依設備設定

每台設備可設 stale 門檻（預設為取值間隔的 5 倍且至少 10 秒）與凍結門檻（預設 0＝停用，因為常數點位凍結屬正常）。設定存於設備的 acquisition 相關設定，而非全域。替代方案：全域門檻——不同設備間隔不同，無法共用。

### 設備專屬輪詢群組

`polling_groups` 新增可空的 `owner_device_id`（唯一），每台設備最多一個專屬群組。儲存取值間隔時建立或更新該群組並把該設備所有點位改指派過去；之後 source rule 產生與 runtime 裝置同步為該設備指派點位時，優先使用設備專屬群組，否則沿用 `Auto Polling`。Migration 為既有設備建立專屬群組；已指派至啟用群組的點位統一採該設備原啟用群組中的最短間隔，停用群組及未指派點位維持原狀。只有單一原間隔時頻率不變，多間隔時原較慢點位會加快。替代方案：把間隔存在 `connection_config`——scheduler 讀的是群組，會出現兩個 authority。

### 既有多間隔設備升級取最短間隔

**Supersedes** 原決策「設備專屬輪詢群組」中所有既有點位升級後頻率不變的前提：2026-11-06 使用者選擇將同一設備的多個啟用群組統一為最短間隔，例如 500／2000 ms 合併為 500 ms。停用或未指派點位不因升級而開始採集；其他設備與重跑冪等性不變。替代方案為保留多個群組，但無法滿足一設備一專屬間隔。

### 自動調整間隔在 scheduler 計算有效間隔

啟用自動調整時，scheduler 以最近 5 輪（設計條件）實測掃描耗時的最大值乘 1.2、取 100ms 倍數，與設定間隔取大者作為有效間隔；連續 5 輪準時後逐步縮回設定值。設定值與有效值都經 runtime status 對外呈現。關閉時維持既有行為：超時以 overrun 回報，不偽造準時樣本。

### 取值優化設定檔沿用既有 acquisition 宣告

連線設定的下拉選單只寫入既有 `connection_config.acquisition.profile`，沿用既有裝置驗證邊界：Modbus 不顯示 random 選項並說明「連續區塊聚合」；MC 提供「未宣告（安全連續讀取）」「safe-v1」「mitsubishi-fx5-cpu-3e-binary-v1」「mitsubishi-q-iqr-3e-binary-v1」；FATEK 提供「未宣告」「safe-v1」「fatek-fbs-ascii-v1」。不新增第二個能力欄位。

### MC random 補齊 double word 與 bit，並新增 Q／iQ-R 設定檔

依官方手冊 SH-080008「Random read in word units (command: 0403)」：請求含 word 點數與 double word 點數兩段；Q/L 與 iQ-R 使用子指令 `0000` 時兩者合計 1–192，iQ-R 使用 `0002` 時 1–96；bit 裝置以 word 存取時 16 點為一單位、以 double word 存取時 32 點為一單位；長計時器與長計數器不可指定；High Performance QCPU 的 ZR 以 2 倍計。`MCClient.RandomRead` 目前把 double word 點數固定為 0 且拒絕 bit 裝置，需擴充為 word 與 double word 兩段並允許 16 點對齊的 bit 裝置。新增設定檔 `mitsubishi-q-iqr-3e-binary-v1`（子指令 `0000`，上限 192，排除 ZR 與長計時器／計數器）；FX5 設定檔沿用 192。optimizer 依 32 bit 型別使用 double word 段，計成本時 word 與 double word 各算 1 點。替代方案：一律只用 word 段並把 32 bit 拆成兩個 word——同樣佔 2 點且需要在解碼端重組，沒有好處。

**Supersedes**: decouple-acquisition-from-delivery / MC random 補齊 double word 與 bit，並新增 Q／iQ-R 設定檔（LTN0 非 random 範例）。2026-11-06 使用者確認維持 `0000` 支援邊界；LTN0 明確回報不支援原因，不送出 `0000` 或 `0002` 讀取。含 D10、D20、LTN0 的規劃可拒絕不支援輸入並附原因，不承諾部分成功。替代方案為加入 `0002`，屬範圍外而不採用。

### random 上限可調但不可放寬

`acquisition.limits` 新增 `random_points`（正整數），作為 random 請求的點數上限，只能小於或等於設定檔的協議上限（MC 192、FATEK 64 以 word 加權）。連線設定以協議上限為預設值顯示並可改小；超過協議上限由既有驗證拒絕。替代方案：允許放寬——會送出 PLC 拒絕的請求。

### 分散位址讀取範圍需補齊 mixed 規則的後端持久化

**Supersedes** 原決策「分散位址讀取範圍沿用 mixed 版面」中「後端已支援 `layout_mode: mixed` 與 `mixed_items`」的前提。2026-11-05 實作 5.3 前查證：`schema.SourceRule` 雖有 `LayoutMode`／`MixedItems` 欄位且 migration 已加欄，但 `sourcerule.SQLRepository` 的 create／update 不寫入這兩欄，非測試程式沒有任何寫入路徑；`acquisitionRuleExclusions` 對 mixed 規則一律回報 `source-layout-unsupported`。因此只改前端會造成「看似可存、實際不存或不採集」的假成功。使用者選擇先補後端：

- 持久化與展開（5.3）：source rule repository 與 Studio v2 工作區規則 API 讀寫 `layout_mode`／`mixed_items`，項目格式沿用 `internal/datalink/measurement.MixedItem`（每項自帶 `address` 與 `data_type`，可跨區域）；儲存前以 `ValidateMixedLayout` 與逐項位址解析驗證，失敗回欄位錯誤；mixed 規則每項展開為一個點位與規則連結並保留項目識別。homogeneous 規則不變。
- 採集（5.4）：宣告 random 設定檔的設備，`acquisitionRuleExclusions` 接受 mixed 規則並交給 optimizer；未宣告、`safe-v1` 與 Modbus 維持 `source-layout-unsupported`。
- 範圍外：量測範本（meter／mixed-sensor templates）的 per-device review 流程、mixed 規則的 Modbus 支援。

### 分散位址讀取範圍沿用 mixed 版面

在 5.3／5.4 的後端支援之上，讀取範圍表單新增「分散位址」模式（5.5）：設備的設定檔支援 random 時才可選；每列輸入位址與型別，可跨區域。表單即時計算 random 用量（MC：word 型別 1 點、32 bit 型別 1 點；bit 以 16 點對齊的 word 計；FATEK：16 bit 算 1、32 bit 算 2），達上限 90% 顯示提醒，超過時說明會拆成多次請求，不阻擋儲存。是否實際用 random 仍由 optimizer 決定（可能把其中的連續段落改為連續讀取）。替代方案：新增第二種規則模型——會與既有 mixed 版面形成兩個 authority。

## Implementation Contract

**Behavior**
- 一個連接器連不上時，ingest 吞吐、SSE 事件延遲與 Modbus Share 寫入不受影響；同時連接器狀態顯示失敗與退避中。
- 舊式資料庫輸出每個 mapping 在每個有效寫入間隔最多寫一列，內容為該時刻最新值；在該間隔內沒有新值時不寫入，並在寫入歷史標記略過原因。
- 首頁每台設備的「讀取點位」按鈕旁有「採集異常」按鈕；有進行中異常時以 warning tokens 呈現並顯示數量文字（例如「採集異常 2」），沒有時為一般次要按鈕。點擊開啟分類為 form 的置中 modal，列出該設備的異常事件，可依類型（斷線、未更新、數值凍結）與狀態（進行中、已結束）篩選，顯示開始、結束、持續時間、點位位址與最後數值。
- 連線設定新增：取值間隔（毫秒，下限 100、上限 3600000，設計條件）、自動調整開關、stale 門檻、凍結門檻、取值優化設定檔；每欄有 DESIGN.md 規定的 ⓘ 說明。儲存後不需重啟即生效；程式啟動時依各設備設定開始採集。

**Interface / data shape**
- 新資料表 `acquisition_anomalies(id, device_id, point_id NULL, kind ∈ {disconnected, stale, frozen}, started_at, ended_at NULL, last_value NULL, detail)`，保存以每台設備 5000 筆為清理目標（設計條件），超出先刪除最舊已結束事件；進行中事件不刪除，即使數量超過目標，事件結束時再次清理。
- `GET /api/v1/datalink/devices/:id/acquisition-anomalies?state=active|ended&kind=...&limit&offset` 回傳事件列表與總數；`GET /api/v1/datalink/acquisition-anomalies/summary` 回傳每台設備進行中事件數。
- 設備更新請求接受 `acquisition_interval_ms`、`acquisition_auto_interval`、`stale_threshold_ms`、`frozen_threshold_ms`；acquisition 設定檔沿用 `connection_config.acquisition`。
- runtime status 的 collector 項目新增 `configured_interval_ms` 與 `effective_interval_ms`。

**Failure modes**
- LTN0 在目前 `0000` 邊界明確回報不支援原因，不送出非法或範圍外的請求。
- 異常 API 讀取失敗時 modal 顯示可行動錯誤與重試，不顯示為「沒有異常」。
- 摘要讀取失敗時按鈕顯示「狀態未知」而非一般樣式。
- 非法間隔或門檻由後端驗證拒絕並回傳欄位錯誤；不以預設值默默取代。
- 取值設定檔與 protocol 不相容時沿用既有 `ErrAcquisitionConfig` 拒絕。

**Acceptance criteria**
- Go 測試：以永不回應的連接器假實作驗證 ingest 處理 32 筆值的耗時與連接器無關（每筆 < 50ms，設計條件），且 SSE 事件數量等於輸入數量。
- Go 測試：快照 worker 在間隔內多次更新只寫入最後一筆；失敗退避序列為 1、2、4、8…60 秒。
- Go 測試：stale、frozen、disconnected 的開始與結束轉換；同點位同類型只有一筆進行中。
- Go 測試：儲存取值間隔後該設備點位全部改屬設備專屬群組，其他設備不變；migration 後單一原間隔不變；多個啟用間隔統一為最短值，停用／未指派點位與其他設備不變，重跑不重複。
- 前端測試：按鈕三態（無異常、有異常、未知）、modal 篩選、錯誤狀態；連線設定欄位驗證與 protocol 對應的設定檔選項。
- 以本機模擬器重跑量測：連接器不可達時 `collected_total` 增長率不低於每秒 32 筆的 95%（設計條件）。

**Scope boundaries**
- 範圍內：MC `0403` double word 與 bit、Q／iQ-R 設定檔、random 上限欄位、mixed 規則的持久化／展開／採集、分散位址範圍表單；舊式資料庫輸出路徑、runtime ingest 與異常偵測、輪詢群組指派、scheduler 有效間隔、相關 API、首頁按鈕與 modal、連線設定表單。
- 範圍外：寫入群組交付、通知推播、Modbus 隨機讀取、L 系列與 iQ-R 子指令 `0002`、資料庫斷線期間的補寫。

## Risks / Trade-offs

- [未分組 append 型 mapping 的資料列數由「每次讀取一列」變為「每間隔一列」，屬可見行為變更] → 在交接與 release 說明中明示；有效間隔沿用既有 mapping／連接器設定，不另設新預設。
- [最新值表只保留最後值，間隔內的短暫峰值不會寫入舊式輸出] → 需要逐筆或聚合的使用者改用寫入群組（已支援 per-member 計算）；在 ⓘ 說明。
- [凍結偵測對常數點位會產生雜訊] → 預設停用，由使用者為設備設定門檻。
- [異常表持續成長] → 以每設備 5000 筆為清理目標，刪除最舊已結束事件；進行中事件超過目標時仍保留，事件結束後再清理。
- [既有多間隔設備統一為最短間隔可能提高原較慢點位的採集頻率] → 依使用者決策採最短值；驗證採集與設備負載仍須使用者授權的現場驗收。
- [自動調整可能掩蓋設備變慢] → 有效間隔與設定值同時顯示，且拉長時記錄一筆說明事件。

## Migration Plan

1. 新增 migration：`polling_groups.owner_device_id`、`acquisition_anomalies` 資料表、設備門檻欄位。
2. 啟動時為沒有專屬群組的設備建立群組，依該設備既有啟用群組的最短間隔移動啟用點位；停用／未指派點位維持原狀，重跑冪等。
3. 快照 worker 與 ingest 解耦後，舊式同步路徑移除；回滾方式為還原 binary，新欄位為可空，舊版會忽略。

## Open Questions

- 「長時間沒更新」門檻預設為取值間隔 × 5 且至少 10 秒，已由使用者採納，可在連線設定修改。
- 2026-11-06 已釐清：保留所有進行中異常，5000 筆為清理目標；LTN0 明確不支援，本次不加入 `0002`。
