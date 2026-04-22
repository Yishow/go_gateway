## 1. Source-step information hierarchy

- [ ] 1.1 調整 `MuiSourceCommandDeck.tsx`，將 active rule summary 提升為 Source step 的前景區塊
- [ ] 1.2 調整 `SourceCanvasSection.tsx`，讓地址覆蓋、規劃狀態、type / scale / naming 摘要可在不讀完整 canvas 的情況下被辨識
- [ ] 1.3 弱化 Source step 的 mode tabs / 次級控制層級，確保 primary reading order 先落在 rule list 與 active rule summary

## 2. Primary action and immediate confirmation

- [ ] 2.1 收斂 Source step 的 primary CTA 為「套用規劃到畫布」，保留既有 apply 邏輯但更新文案與階層
- [ ] 2.2 在 Source step 補上套用後的即時回饋，明確顯示 rule 已成形與目前覆蓋摘要
- [ ] 2.3 更新相關 i18n 文案（`zh-TW` / `en`），讓 Source-step CTA、confirmation、handoff copy 對齊新語意

## 3. Rule-scoped Tag preview continuity

- [ ] 3.1 重用 `useSourceRuleCandidatesQuery` 與 `sourceRuleSelection.ts`，在 Source step 建立 rule-scoped Tag preview 資料流
- [ ] 3.2 支援 click 固定 preview 與 hover 暫時 preview，且離開 hover 後恢復原本 focused rule
- [ ] 3.3 確保 Source-step preview 與 Tag step 實際 candidate scope 一致，且返回 Source step 時 continuity 不丟失

## 4. Database-aware planning defaults and hints

- [ ] 4.1 將 `activeOutputTarget === 'database'` 接入 Source-step planning defaults，補上 target data type / scale / naming 的 advisory defaults
- [ ] 4.2 重用 `databaseGroupingSuggestions.ts` 與既有 row planner semantics，在 Source step 顯示 table/row/column 導向的 grouping hints
- [ ] 4.3 確保 database-aware 變化僅限於 defaults 與 hints，不引入 connector / schema / table setup controls 到 Step 2

## 5. Verification and regression coverage

- [ ] 5.1 更新 `DatalinkWorkbenchSourceStep.test.tsx`，覆蓋 active rule summary、primary CTA 與 mode hierarchy 變更
- [ ] 5.2 新增或擴充 Source-step preview / rule-summary 測試，覆蓋 hover/click preview continuity
- [ ] 5.3 更新 Tag / Database 相關測試，驗證 Source preview 與 Tag review scope、database-aware defaults / hints 的契約
- [ ] 5.4 執行 `cd frontend && npm run lint`
- [ ] 5.5 執行 `cd frontend && npm run test -- --run`
- [ ] 5.6 執行 `cd frontend && npm run build`
- [ ] 5.7 啟動並手動驗證 `/studio?target=database&step=source`：選 device → 規劃 rule → 套用 → 預覽 Tag review → 前往 Tag → 返回 Source
