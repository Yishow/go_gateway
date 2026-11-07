## Problem

使用者平板截圖顯示，記憶體選取格內的長說明把整排格子撐高，與原五張參考圖的緊湊位址結構不符。手機和平板的 modal 仍保留桌面外緣，減少可用空間；使用者要求此類裝置全部 modal 改全螢幕。

## Root Cause

RecordingHomeMemoryCells 的 gw-memory-focus-information 在 data-expanded=true 時以一般文字顯示於窄格，recordingHomeMemory.css 同時增加上距，造成 grid row 依最長內容增高。semantic-modal.css 只有小於 640px 的近滿版規則，來源 modal 還有獨立 85% 尺寸，其他 dialog primitives 未共用相同響應式契約。

## Proposed Solution

- 明確更新 form、confirm、來源/記憶體與既有 modal primitives 的手機/平板全螢幕契約，保留桌面尺寸、safe area、內部捲動與操作保護。
- 選取資訊改為獨立可見且可鍵盤到達的區域，不再撐大記憶體格；完整位址、型別、值、狀態仍可取得。
- 五種視圖保持原構圖差異與 DESIGN.md 的配色/字體權威，整理平板工具列、間距和選取層級。

## Capabilities

### Modified Capabilities

- simple-recording-home: 手機與平板 modal 改全螢幕並維持桌面尺寸分類。
- recording-memory-visual-refinement: 選取詳情與固定密度的記憶體格分離，補平板與觸控驗收。

## Success Criteria

- 設計條件：CSS viewport 寬度 ≤1024px，或 coarse pointer 且寬度 ≤1366px 時，所有 modal 填滿可見 viewport；1440px 一般桌面維持現有尺寸。
- 選取含長名稱/長值的格子後，相鄰格不因完整描述而被撐高；完整內容可見可達，未知與衝突狀態不被美化成成功。
- 實測手機、直/橫平板、桌面上的五視圖與 modal primitives，包含選取、旋轉、關閉、Escape、dirty guard；未執行不得標通過。

## Impact

- Affected code:
  - Modified: DESIGN.md、frontend/src/styles/semantic-modal.css、frontend/src/components/ui/dialog.tsx、frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeMemoryView.tsx、frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeMemoryCells.tsx、frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeMemory.css、frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeMemoryLayouts.css，以及實際已存在的 ConfirmDialog/SlidePanel 與來源 modal 樣式。
  - New: frontend/tests/e2e/tablet-modal-layout.spec.ts（本案回歸驗收）。
  - Removed: 無產品模組移除。
