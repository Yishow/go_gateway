## Context

基準為 db91bbdb，保留使用者最新記錄點位橫向捲動修正。本案只修 modal 與記憶體五視圖；既有 add-switchable-memory-views 功能與歷史驗收不重寫。

## Design Source

來源是使用者本次平板截圖與 DESIGN.md 列出的五張明暗並列原圖（51e970e9 grid、86b48a15 navigator、c91df0e8 map、943be82c chapters、201f2c58 detail）。沒有 .dc.html，不能聲稱從 PNG 得到精確設計值或像素一致。

原圖決定緊湊格子、群組關係與選取資訊的層級；色彩沿 DESIGN.md Quiet mineral。以下為本案設計條件，不是 PNG 量測：手機/平板斷點 1024 CSS px；coarse pointer 平板上限1366px；既有 spacing 4/8/12/16/24px、panel radius8px、觸控目標至少44px。全螢幕外框radius0，safe area落在內部間距；桌面原75% form／85% memory與小confirm分類保留。

## Goals / Non-Goals

**Goals:**
- 全部真正 modal 在手機與平板共享全螢幕、標題/動作可達、內容內捲契約。
- 長說明移到格外，選取不再擠壓整排；所有地址、完整datatype、長值與不確定狀態仍可讀。
- 五圖構圖、明暗語意與資料權威維持一致。

**Non-Goals:**
- 不更改協議、API、後端、寫入/試寫/ownership/ready token、資料authority或存檔行為。
- 不重寫使用者剛更新的記錄點位表、不增加V3或新介面架構、不將非modal hover資訊窗全螢幕化。
- 不修改原Mac工作區、既有封存證據或假稱owner已驗收新畫面。

## Decisions

### 全螢幕分類共用但限定真正 modal

以明確的modal class覆蓋Radix與自製Dialog、ConfirmDialog、SlidePanel，不能用全域role=dialog。寬≤1024px，或pointer:coarse且寬≤1366px採全螢幕；既有低高度防遮蔽也保留。來源記憶體85%尺寸覆寫需一起調整，不能由載入次序決定實際大小。替代方案僅改semantic-modal.css不足以涵蓋其他primitives，故不用。

### 選取資訊脫離格子並維持穩定空間

在footer內設可捲動完整說明區，始終預留高度，未選取時保留簡短操作引導。格子維持地址/型別/讀值；focus、blur或hover不能重排整列或讓操作按鈕跳位。完整aria-label保留。相較改為tooltip，固定區也能供觸控與鍵盤讀完整內容；相較選取後才插入新段落，預留空間可避免點擊途中布局位移。

### 可見值訂閱保留canonical地址

每格的data-value-key使用實際值的canonical起始地址，包含跨頁continuation；footer資訊使用同一key。IntersectionObserver仍負責唯一父層可見值集合，不另加輪詢owner、不把未知/partial當成可讀值。

## Implementation Contract

- 在既有modal primitives與consumer尺寸CSS中套用同一斷點，背景阻擋、Escape、dirty guard、已送出操作、焦點trap與返回焦點不改。非modal MappingSaveInfo不受影響。
- header/footer固定在modal框內，body可捲，safe area不遮標題、關閉或儲存；表單長內容、旋轉和低高度時仍能到達錯誤與操作。
- RecordingHomeMemoryView共用area/window/focus/owner及values，RecordingHomeMemoryCells不再以description擴張列高；footer區提供完整且可捲讀的describe結果。
- 任何unknown/conflict/stale/error保持既有狀態優先與動作限制；格子點擊、keyboard箭頭、跨列/頁、五視圖切換沿原契約。
- 驗收以新增回歸與現有memory/modal相關測試，並以實際agent-browser/Playwright檢查390×844、768×1024、1024×768、1366×1024 touch與1440×900 desktop；五視圖需明暗對照。
- 前端完整lint/test/build於最終版本執行；browser不可執行時保留blocked，不拿DOM測試或舊截圖替代。此次純前端，不將後端未重跑說成已驗。

## Risks / Trade-offs

- [共同CSS被consumer覆寫] → 檢查來源host、confirm與legacy primitives的實際computed box。
- [完整描述過長] → 固定資訊區可內捲，鍵盤可到達；不能裁切唯一資料。
- [觸控判斷不是裝置型號] → 使用CSS viewport與pointer能力，明列1366px邊界，旋轉實測。
- [新增footer造成主要格子太少] → 五視圖手機/平板視覺檢查，同時整理工具列與間距，不縮小必要字體和觸控目標。
