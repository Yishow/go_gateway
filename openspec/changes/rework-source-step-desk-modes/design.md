## Context

目前 Step 2 的 Source workspace 已經具備 active rule summary、grouped-tag handoff 與 database-aware hints，但 desk mode 的層次還不夠清楚。`巡檢` 主要呈現畫布，容易讓 operator 忽略它其實是在驗證 rule planning；`建造` 讓 rule builder 佔據太多垂直空間，導致畫布只剩極小的輔助區；`分診` 則只在有 incident queue 時稍微不同，沒有問題時幾乎退化成巡檢版面。這三種 mode 缺乏可辨識的工作姿勢，讓 Source -> Tag -> Output 的手off 主線被 mode 切換打斷。

目前前端責任也偏混雜：`MuiSourceCommandDeck.tsx` 管 desk 切換與 shell 包裝，`SourceCanvasSection.tsx` 同時承擔共享骨架、畫布、conflict queue 與 mode-specific 內容，`MuiWorkbenchSourceStyles.tsx` 以大量條件樣式支撐布局差異。若直接在現況上繼續疊條件分支，後續很難維持可讀性與測試穩定性。

## Goals / Non-Goals

**Goals:**

- 讓 `巡檢 / 建造 / 分診` 各自對應清楚的工作目的與可辨識主畫面。
- 在 mode 切換時保留同一條 Source-step planning skeleton：active rule summary、handoff strip、desk selector、diagnostic cues、primary CTA。
- 讓 `建造` 保持 editor-primary，但同時保留可閱讀的畫布視窗作為即時規劃回饋。
- 讓 `分診` 成為真正的 issue-first workspace，聚焦 conflict、blocked、unmanaged 項目與修復路徑。
- 降低 `SourceCanvasSection.tsx` 的 layout 分支複雜度，並用 targeted tests 鎖定 desk contracts。

**Non-Goals:**

- 不改變 Source rule、Tag candidate、Database preview 的後端契約或資料模型。
- 不重做 `Plan / Live / Link` 畫布 information layer，也不改變其底層 address geometry。
- 不重新定義 Step 3 與 Step 4 的主流程，只調整它們在 Source step 中的 handoff 呈現方式。
- 不處理這次提案以外的 startup、embed build、port override 或 shell routing 問題。

## Decisions

### Preserve one shared source-step planning skeleton

所有 desk mode 都共用同一條 Source-step 骨架：active rule summary、handoff strip、desk mode tabs、diagnostic cues 與 primary CTA 固定留在 Source workspace 的上層區域。mode 切換只改變主工作區與次要輔助區的優先順序，不改變 operator 對「目前是哪條 rule、下一步要交給誰、目前被什麼擋住」的閱讀位置。

這個做法比把三個 mode 做成三套幾乎獨立的頁面更安全，因為它保留單一路徑的 guided workflow，也避免把 Step 2 再做成隱性的子路由。

### Split mode-specific panels out of SourceCanvasSection

`SourceCanvasSection.tsx` 保留 orchestrator 角色，持有共享資料、選取狀態與 canvas coordination；mode-specific 內容拆到獨立 panel，例如 `SourceRuleLayerPanel.tsx` 與 `SourceTriagePanel.tsx`。`MuiSourceCommandDeck.tsx` 保留 mode state 與 shell 邊界，`MuiWorkbenchSourceStyles.tsx` 改為宣告 desk layout contract，而不是承載所有條件式內容。

替代方案是繼續把 `inspect/build/triage` 分支留在同一檔案內，但這會讓已經偏大的 Step 2 組件再度膨脹，也讓測試更難針對 mode 合約定位。

### Keep build mode editor-primary without collapsing the canvas

`建造` 模式會把 rule editor / planning layer 放在主要焦點區，但畫布必須保持可用且可讀，至少能讓 operator 同步看到規劃 span、unmanaged 區間與 live address context。實作上會避免目前 `maxHeight: 180px` 這種近乎壓扁的次要畫布，改成保留具有最小可讀高度的 secondary canvas region，並讓 builder 區與畫布在桌面寬度下維持穩定比例。

替代方案是把 builder 做成全螢幕 modal 或完全分頁，但那會把即時規劃回饋切離，與 rule-first planning 的工作方式衝突。

### Make triage mode issue-first with filtered recovery cues

`分診` 模式改成真正的 issue-first workspace：優先顯示 conflict / blocked / unmanaged queue、對應 recovery action 與必要的 filtered canvas context。當目前沒有需分診項目時，分診模式必須出現明確 empty state，說明目前無待處理衝突，而不是退化成與巡檢幾乎相同的畫面。

替代方案是只替 triage 套不同色彩或把 conflict box 往上移，但這種做法無法建立清楚的行為差異，仍會讓 operator 以為 triage 只是 inspect 的另一個皮膚。

### Lock desk contracts with targeted source-step tests

這次重構以 mode contract 為測試核心：檢查共享 skeleton 在 desk 切換後仍存在、`建造` 的 editor/canvas 佈局 contract、`分診` 的 issue queue 與 empty state、以及 incident recovery 導回正確 surface。測試會集中在 `DatalinkWorkbenchSourceStep.test.tsx`、incident desk reopen tests、以及 `frontend/tests/unit/pages/datalink/workbench-source-step.test.tsx`。

替代方案是只保留 snapshot 或單一路徑 smoke tests，但那不足以防止未來把 triage 再改回 inspect 變體，或讓 build desk 再次把畫布壓縮到不可讀。

## Risks / Trade-offs

- [Risk] Source step shell 與樣式責任重新分配後，可能造成既有 desk-specific CSS 選擇器失效。 → Mitigation：同步調整 `MuiWorkbenchSourceStyles.tsx` 與 mode tests，以 DOM contract 而不是脆弱的細節 class 名稱作為驗證。
- [Risk] 抽離 panel 元件可能讓共享 state props 變多。 → Mitigation：維持 `SourceCanvasSection.tsx` 作為資料協調器，只把展示責任移出，避免在多個 panel 之間複製 domain 狀態。
- [Risk] `分診` 若沒有 issues，新增的 empty state 可能被誤解成功能空白。 → Mitigation：empty state 必須明確說明 triage 的用途、當前無阻塞項目，以及何時應回到巡檢或建造。
- [Risk] `建造` 保留較大畫布後，可能讓 editor 區首屏資訊減少。 → Mitigation：將 builder 的摘要、欄位群組與 CTA 緊縮成更高密度 layout，優先保留 active editing 所需的欄位與回饋。

## Migration Plan

1. 先以測試定義新的 desk contracts，讓現況失敗並明確描述三種 mode 的預期主區域。
2. 重構共享 skeleton 與 mode-specific panels，完成 `巡檢`、`建造`、`分診` 三種 layout contract。
3. 更新 i18n 文案與 incident/recovery copy，確保新的 handoff 與 empty state 可被辨識。
4. 重新執行前端 lint、unit tests、build，並以實機 UI 檢查 `/studio` 的 Step 2 mode 切換。

## Open Questions

- 無。這次提案沿用既有 `巡檢 / 建造 / 分診` 命名，只調整其工作語意與畫面責任。
