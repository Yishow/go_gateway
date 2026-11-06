## Summary

在既有五種記憶體視圖上校正層級、型別卡片、間距與各版面辨識度，延伸 DESIGN.md 為唯一設計 authority。

## Motivation

五張原始明暗參考已保存，現有功能完成紀錄保留，但新的 source／圖片檢視發現 chapters 使用側欄，與參考的上方全寬章節不同。此案處理具體視覺落差，不把歷史功能驗收當成此次實作的逐圖驗收。

## Proposed Solution

- 延用 grid、navigator、map、chapters、detail 共用資料與互動；明確定義每種構圖、type badge／surface／outline 與跨格連續性。
- 原圖提供幾何與資訊層級；配色沿 DESIGN.md Quiet mineral semantic tokens，不引入全站藍色、漸層、玻璃或左色條。
- 保留 memory modal 85vw／85dvh 的設計條件及窄螢幕可達性，不為對圖硬塞桌面固定尺寸。
- 對五圖逐一記錄 light/dark 與桌面／平板／手機對照，真瀏覽器阻擋必須明列，不能以 DOM 或生成圖替代。

## Impact

- Affected specs: recording-memory-visual-refinement
- Affected code:
  - Modified: DESIGN.md, frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeMemoryView.tsx, frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeMemoryCells.tsx, frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeMemoryPanels.tsx, frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeMemory.css, frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeMemoryLayouts.css, frontend/tests/unit/workbench-v2/recording-home-memory-layouts.test.tsx, frontend/src/i18n/locales/en/recording-home.json, frontend/src/i18n/locales/zh-TW/recording-home.json
  - New: frontend/tests/e2e/recording-memory-visuals.spec.ts
  - Removed: 無。
- add-switchable-memory-views 已有五項歷史完成任務，本案是視覺後續責任，保留其 artifacts／checks。不 archive、重勾或重寫原證據；不變更 API、協議、subscription authority 或全站 theme。
