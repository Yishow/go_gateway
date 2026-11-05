## Context

基準 a9539351 的 RecordingHomeMemoryView 已顯示即時讀值，父層 RecordingHomeSourceContent 持有單一 useHomeMemoryValues。buildMemoryMap 與 memoryWindow 是占用及有界分頁權威，並保留 MC revision proof、string-width 與 bit-area fail-closed。

## Design Source

來源：使用者於 2026-11-04 確認保留的五張明暗並列 PNG，識別為 51e970e9、86b48a15、c91df0e8、943be82c、201f2c58，依序為格子專注、範圍導航、水平地圖、範圍章節、格子＋詳情。原圖與 SHA256 保存在 repo 外私人驗收目錄；沒有 .dc.html，不製造不存在的來源。若需移交圖片，須由私人驗收資料取得，不將圖片提交 Git。

圖中精確幾何尚未測量，以下皆為沿用 DESIGN.md 的設計條件：spacing 4/8/12/16/24px、control radius 6px、panel 8px、dialog 12px；light ground #F7F8F5 / dark #191A19；light text #203332 / dark #EFEFEB；控制高度桌面40px、手機44px；沿用 Inter/Noto Sans TC 與現有 lucide outline。完整型別徽章、單格位址及相連多格為共同視覺語彙。圖內假的位址／計數／重複字樣不複製，主題不抄圖中漸層。

## Goals / Non-Goals

**Goals:** 五種可切換表面、相同資料與動作、安全不退化、明暗與響應式、逐功能及逐圖驗收。

**Non-Goals:** 不改後端、協議型別寬度、訂閱設計、資料持久化、現場 PLC／DB、其他首頁或全站換皮。

## Decisions

### 共用狀態與安全權威

在同一 RecordingHomeMemoryView 保留 areaKey、starts、focus、owners、locate，新增 layout 與 inspect selection。單次只渲染一個實際 grid；不用五個隱藏副本或重掛父層。替代的每視圖獨立 hook 會丟狀態與重訂閱，拒絕。

### 共用格子與型別連接

以每個 physical register 保留可訪問位址／鍵盤焦點，相同 value owner 的相鄰格透過群組外觀相連；跨列／頁顯示起始位址及續接。型別徽章及可辨識低飽和底色共用，無左色條；保留 unknown、partial、skipped、conflict。不把 string 的未確認寬度畫成已知十格。

### 五種版面與共同控制

頂部一列 layout 切換，區域、定位、分頁與狀態仍共用。grid 為滿寬；navigator 加左側範圍清單；map 加當前視窗占用條與選取摘要；chapters 加範圍可展開摘要與跳轉；detail 加選取格資料與明確開啟動作。保留格子原 click 契約，詳情由 focus/hover 與明確選取控制更新，不把查看改成額外必按一次。手機側欄折成上方可收合區或下方詳情，所有資訊仍可到達。

## Implementation Contract

- 介面延用 MemoryArea、MemoryCell、MemoryValues、onOpen/onAdd；layout 僅 UI state，不寫 server 或新訂閱。
- empty click 傳精確 address；used 開 rule；conflict 顯 owner chooser；unconfirmed 不可新增。recovery/local/stale 仍封鎖新增，已知 owner 仍可開。
- bool/int16/uint16 一 register，int32/uint32/float32 二 register，int64/uint64/float64 四 register；partial 不假定連續。
- 各 renderer 保持 data-offset/data-value-key、焦點與 visible reporting，layout 改變重綁 observer 和測量實際 columns，不能重建 live hook。
- 範圍導航定位到 shared window，不產生第二份來源資料；所有數字由實際 map/input/value 推導。
- Loading/error/checking/unsupported 與 retry 沿用，不能改作空格。
- 驗收以五版面參數化單元測試、現有 memory map/scope/values/keyboard 測試、每個實際按鈕 E2E、明暗 1440/768/390 screenshots；瀏覽器拒絕時清楚標 blocked，不宣稱通過。
- 開發只跑相關測試；完成後一次完整 frontend lint/test/build、diff/line gates。未改 backend 不擴大後端實作。

使用者補充希望使用已保留 HSL simulator 驗證數值。原使用者上傳 binary 的實際執行須取得本輪逐次同意；僅在授權且 sandbox 允許時啟動 owned loopback MC／Modbus，使用可丟棄資料與 typed seed/readback，驗證即時值更新及失聯狀態。未取得同意或 browser/socket 被拒則保留 pending，不用 mock 充當 simulator E2E。

## Risks / Trade-offs

- [較窄主 grid 改變 wrap] → 使用 computed columns 並測 W hex 與四格值跨列。
- [layout DOM 改變使 visible IDs 失效] → effect 明確依賴 layout，清理 observer，保留單一 hook。
- [圖中樣本失真] → 真 map 與 typed value 勝過生成圖的位址與計數；圖只作層級與排版參考。
