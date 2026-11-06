## Context

main 0f54740 已具五視圖；add-switchable-memory-views 在 active directory 但五項歷史任務已完成，本案不重開或 archive 原案。其 recording-memory-layouts 尚未合併成 canonical spec，因此新增明確命名的 refinement capability，不冒充 MODIFY 不存在的 canonical spec。

延續 add-switchable-memory-views / 共用狀態與安全權威（避免丟狀態與重訂閱）、共用格子與型別連接（保留 physical geometry）、五種版面與共同控制（同一資料投影）。原案已拒絕複製圖片漸層，本案不反轉它。

## Design Source

2026-11-06 在雲端逐張檢視使用者先前保存的五張明暗並列 PNG。沒有 .dc.html，不能製造不存在的 HTML 來源。私人 cache 為 .spectra/design-cache/refine-memory-view-visuals/，manifest.json 記精確檔名、SHA-256、像素尺寸；tracked reference-manifest.json 同步保留完整識別與 recovery 契約。原圖不提交 Git。此雲端提案交接包包含私人 originals/，是可攜輸入來源；若 fresh checkout 只有 Git artifacts，必須向 owner 取得該交接包或逐一重附 manifest 指定原圖，核對 SHA 後放回 cache。未取得前停止視覺 apply，不從短 ID 猜下載網址。

| 版面 | 來源 ID | 原圖像素（檔案 metadata 實測） |
| --- | --- | --- |
| grid | 51e970e9 | 1821×864 |
| navigator | 86b48a15 | 1752×897 |
| map | c91df0e8 | 1752×897 |
| chapters | 943be82c | 1752×897 |
| detail | 201f2c58 | 1774×887 |

圖中 px/hex/spacing 未做幾何量測；不得把設計值冒充抽樣。可執行數值只由 DESIGN.md 的「記憶體五視圖增補」及既有 token tables 決定：該節已記 exact recipe、radius、spacing、85% 與對比。沿用既有 lucide outline icons，不增加 decorative icon system；實作前若 cache 不可讀，先取得原圖而非憑文字重畫。

## Goals / Non-Goals

**Goals:** 各視圖資訊層級可辨識，五圖構圖差異被落實，共用元件明暗一致且可鍵盤操作。
**Non-Goals:** 不新建視圖、theme engine、真實資料或健康指標；不修改後端、protocol、occupancy、live subscription authority。原圖缺少的安全控制仍是產品既有必要契約，其 placement 與溢出由本案一起驗證，不能只拿 non-goal 掩蓋畫面差異。

## Decisions

### 單一設計權威與有意取捨

DESIGN.md 是唯一 geometry/token recipe authority；本 change 只連結，不維護第二套值。原圖配色與現行 Quiet mineral 不同，選用原圖 geometry＋既有 palette，明列為有意取捨，不聲稱逐 pixel 複製。記憶體 85vw/85dvh 是既有 scoped exception，其他 modal 仍沿原契約。

### 共享卡片與分明構圖

先校正 RecordingHomeMemoryCells 的 address/type/value、unknown/conflict 與 multi-register continuation；五種 renderer 共用該 recipe。再分批調整 grid/navigator、map/chapters、detail。chapters 的完整寬度章節帶放 grid 上方，是與 navigator 左欄的具體分界。Map 不拿 detail 大卡替代摘要，也不以抽樣省略 window slots。

不用固定 reference column count；實際 CSS columns 驅動 keyboard/wrap，單次只 mount 一份 grid。所有過往 component behavior、dirty state、selection、owner chooser 與 MC proof 都要逐步回歸。

### 逐圖與互動雙重驗收

DESIGN.md 列出的基礎矩陣以及窄高、長值、unknown/conflict、跨列／頁情境全部保存產品截圖和對照說明；每種視圖記 accepted／intentional difference／defect 三種結論，不能只貼一張代表五種。

真 browser 同時驗 layout switches、keyboard、hover、touch、Close/Cancel、Back/Forward、重複點擊、observer 及唯一 parent hook。Supersedes：原雲端限定／不動 Mac 執行次序由 owner 最新指示取代：四案雲端實作及測試撰寫完成後統一在 Mac 執行 browser；之前保持待驗證，不把 jsdom 當 browser，不啟動未另行授權的 HSL binary／實體設備。

## Implementation Contract

- 只有 DESIGN.md 與 proposal 指定 memory UI/CSS、locale、tests；recording-home.css 的 memory selector 可作 regression/assertion，其他 home/modal 規則不變。
- 新增 recording-memory-visual-refinement 是原 recording-memory-layouts 的呈現後續，既有 safe actions/readValues 不改。
- 凡 source 缺失／unknown/string-width/bit proof 不完整均 fail closed；格子中的設備讀值不是 SQL receipt/readback。
- 對應單元：memory-layouts/panels/continuity/hover/keyboard/scope/recovery/observer-states/values/map/view；新 e2e 加入既有 memory-sizing/presentations/navigation-matrix/hover-matrix/recovery-matrix 範圍。
- A（optimize-embedded-route-loading）若移動來源 component owner，先完成其界線再做本案；B/C 不改 memory renderer，可於 A 後以不同檔案工作，合併時仍重跑 affected tests。
- In scope 是視覺與可達性；Out of scope 是色彩品牌重做、data model、全站重寫、API、部署。本次 owner 已授權依本案實作；browser 與 archive 仍依最終合併次序保留。
- 每批測試先行、可 build，最後一次 frontend lint/test/build、diff/line gate；browser/field evidence 分開，不借歷史 118 項結果勾新任務。

## Risks / Trade-offs

- [原圖與產品資料不同] → 圖只作 hierarchy，真 identity/typed state 優先，逐圖記有意差異。
- [柔和型別底色對比不足] → 實際相鄰背景測對比；只改 DESIGN.md recipe 並同步 CSS，不換全站 palette。
- [chapters 推開 grid 或窄屏裁切] → 有界章節滾動、保留定位與實際低高度 browser 測試。
