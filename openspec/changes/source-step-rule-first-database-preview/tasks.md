## 1. Guided workflow and Source-step information hierarchy

- [x] [P] 1.1 對齊 `Guided workflow` 與「維持既有 active-rule contract，只重組 Source step UI 結構」：調整 `frontend/src/pages/datalink/workbench/SourceCanvasSection.tsx`，將 active rule summary 放在畫布上方前景區塊，先回答目前正在編哪一條 rule
- [x] [P] 1.2 對齊 `Source planning surfaces grouped-tag handoff context` 與「將 rule summary + preview 提升到 canvas 之前」：調整 `frontend/src/pages/datalink/workbench/SourceCanvasSection.tsx`，讓 summary 至少顯示起始位址、count、data type 與 planned / used / conflict 覆蓋摘要
- [x] 1.3 對齊 `Source step prioritizes planning actions over mode switching` 與「mode tabs 改為次級控制，不再佔主導」：調整 `frontend/src/pages/datalink/workbench/SourceCanvasSection.tsx`，弱化 Source step 的 mode tabs / 次級控制層級，確保 primary reading order 先落在 rule list 與 active rule summary
- [x] 1.4 調整 Source-step 文案與區塊說明，明確把畫布定位為驗證與診斷 surface，而不是主要規劃入口

## 2. Primary action and immediate confirmation

- [x] 2.1 收斂 Source step 的 primary CTA 為「套用規劃到畫布」，保留既有 apply 邏輯但更新文案與階層
- [x] 2.2 在 Source step 補上套用後的即時回饋，明確顯示 rule 已成形與目前覆蓋摘要
- [x] 2.3 更新相關 i18n 文案（`zh-TW` / `en`），讓 Source-step CTA、confirmation、handoff copy 對齊新語意

## 3. Rule-scoped Tag preview continuity

- [x] [P] 3.1 對齊 `Source step previews rule-scoped Tag review before navigation` 與「Tag preview 採用 rule-scoped candidate snapshot，而不是複製 Tag step state」：重用 `useSourceRuleCandidatesQuery` 與 `sourceRuleSelection.ts`，在 Source step 建立 rule-scoped Tag preview 資料流
- [x] 3.2 支援 click 固定 preview 與 hover 暫時 preview，且離開 hover 後恢復原本 focused rule
- [x] 3.3 確保 Source-step preview 與 Tag step 實際 candidate scope 一致，且返回 Source step 時 continuity 不丟失

## 4. Database-aware planning defaults and hints

- [x] [P] 4.1 對齊 `Database context MAY influence Source-rule defaults without overriding operator intent`：將 `activeOutputTarget === 'database'` 接入 Source-step planning defaults，補上 target data type / scale / naming 的 advisory defaults，且保留使用者顯式覆寫
- [x] [P] 4.2 對齊 `Database target context informs Source-step planning hints` 與「`target=database` 僅影響 planning defaults 與 grouping hints」：重用 `databaseGroupingSuggestions.ts` 與既有 row planner semantics，在 Source step 顯示 table/row/column 導向的 grouping hints
- [x] 4.3 確保 database-aware 變化僅限於 defaults 與 hints，不引入 connector / schema / table setup controls 到 Step 2

## 5. Verification and regression coverage

- [x] [P] 5.1 更新 `frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchSourceStep.test.tsx`，覆蓋 active rule summary、primary CTA、畫布前景順序與 mode hierarchy 變更
- [x] [P] 5.2 新增或擴充 Source-step preview / rule-summary 測試，覆蓋 hover/click preview continuity 與畫布作為診斷 surface 的文案契約
- [x] [P] 5.3 更新 Tag / Database 相關測試，驗證 Source preview 與 Tag review scope、database-aware defaults / hints 的契約
- [x] 5.4 手動驗證 `/studio?target=database&step=source` 是否能在 5 秒內辨識 active rule、覆蓋摘要與 primary CTA
- [x] 5.5 執行 `cd frontend && npm run lint`
- [x] 5.6 執行 `cd frontend && npm run test -- --run`
- [x] 5.7 執行 `cd frontend && npm run build`
- [x] 5.8 啟動並手動驗證 `/studio?target=database&step=source`：選 device → 規劃 rule → 套用 → 預覽 Tag review → 前往 Tag → 返回 Source
