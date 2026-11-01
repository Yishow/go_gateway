## Summary

先把原Step4接成固定v9 demo方向的單頁簡易呈現與總覽新增／修改／返回；「已選demo主旅程＋其必要生命週期安全」經過驗證與既有使用者頁面斷點確認後，將同一 Studio V2 的記錄首頁作為無明確意圖時的預設入口；保留原進階流程與所有明確深連結。驗收未完成時維持可選入口。

## Motivation

只新增連結而不定義採用條件，會讓使用者永久面對兩套入口；過早切換又會把尚未驗證的流程變成主線。需要清楚的完整驗收、可回退呈現與證據交付，不能以demo可點或單元測試通過宣稱正式完成。

## Proposed Solution

- 新增 simple-recording-entry-adoption 能力：預設採用有明確gate，root/wildcard仍收斂到 /studio/v2，進階用明確view保留。
- 先接首次建立、已有來源、保存／reload、修改／返回、明確Start與真實結果；最終仍逐項驗已選修改／停止／恢復／複製／移除／表格及目的地前向變更與必要失敗／中斷。lifecycle按既有server availability走原owner，真正unsupported明示原因／原advanced交接，不造命令、不增新引擎；必要動作pending不得隱藏以提前default。
- 完成相關測試後對最終功能批次跑完整前端lint/test/build與嵌入式交付檢查；後端僅在本批有改動時跑完整對應檢查。
- UI瀏覽器及真設備/SQL證據分列；環境限制須呈現NOT RUN，不繼承歷史waiver。切換預設與發布都須遵守現有授權，不擅自部署。

## Impact

- Affected specs: simple-recording-entry-adoption (new)
- Affected code:
  - Modified: `frontend/src/App.tsx`; `frontend/src/features/datalink/legacyRoutes.ts`; `frontend/src/services/studioEntityLocator.ts`; `frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx`; `frontend/src/services/studioRecordingView.ts`; existing `recording-home` and Step4 `WriteGroupSection` / `GroupEditor` presentation
  - New: frontend/tests/e2e/simple-recording-journeys.spec.ts; docs/releases/simple-recording-entry.md
  - Removed: none
- 接線使用 simplify-recording-home-guidance 與 streamline-recording-task-editor 的既有契約；simplify-recording-lifecycle-recovery 及採集優化保留原owner與pending項，不作未用到能力的整包接線前置。本案負責簡易host整合及所選主流程的最終預設gate；不把Settings全頁或大型效能矩陣綁入。

本次更新納入 `DESIGN.md` 語意 tokens、system/light/dark 主題、RWD、簡潔文案與字型階層；新彈出設定為 viewport 寬高各75%的置中 modal。已選手機近滿版／小確認框較小例外，分類依 DESIGN.md；舊 demo 不構成此次確認。

最終採用另依賴 support-per-member-recording-calculation 與editor consumer；新增同記錄mean/mean/delta各member自己的對齊窗口、到期同列驗收，不以拆多record或demo效果代替。


本輪已確認逐member interval與UTC epoch對齊：5秒00/05…55、15秒00/15/30/45、10分鐘00/10…50；wide採due邊界聯集，:15/:30/:45的60秒delta留NULL＋not_due，:00同列各自前15秒mean／前60秒delta；long只新增到期member列。不carry、不補0；required-not-due不skip，到期失敗／無baseline／gap/reset分別標原因。顯示各member的[start,end)、partial coverage與next due，SQL commit時間另列；first partial、nullable schema、安全gate與durable驗收依support-per-member-recording-calculation。此語意確認不是新版demo或apply核准。

採集次數預覽只消費 integrate-protocol-aware-acquisition-batching 的同planner、scoped/revision-bound read-only API；計畫估算與實際協議attempt／retry分開，不能用demo或選取點數宣稱效能改善。此依賴只限預覽／批次能力，不阻擋獨立member時間窗決策。
