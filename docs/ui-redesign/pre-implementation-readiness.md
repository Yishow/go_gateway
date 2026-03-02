# 實作前準備清單（Pre-Implementation Readiness）

日期：2026-03-02

---

## A. 規格與決策確認
- [x] 決策鎖定：1A / 2B / 3C
- [x] 已有 SPEC：`two-entry-ui-spec.md`
- [x] 已有 PLAN：`two-entry-ui-implementation-plan.md`
- [ ] 產品/工程/測試三方評審完成

## B. 技術前置
- [ ] 路由與旗標命名定稿（`feature.datalink.twoEntryUI`）
- [ ] store 與 adapter 介面定稿
- [ ] 舊路徑相容策略（導向/保留）定稿

## C. 測試前置
- [ ] Quick Setup 五步流程測試案例建立
- [ ] Workbench 對照案例建立
- [ ] 入口切換不丟草稿案例建立
- [ ] E2E happy path + error path 案例建立

## D. 交付前置
- [ ] PR 模板加入「雙入口檢查清單」
- [ ] 監控指標儀表（完成率/錯誤率/回滾率）準備
- [ ] 灰度開關與回滾腳本驗證

## E. 視覺決策前置（Gemini 生圖）
- [x] 流程預覽圖已產出：`preview_flow.svg`
- [x] 最終線框預覽圖已產出：`final_preview_wireframe.svg`
- [ ] 使用者最終採用版本確認

---

## 啟動條件（Go/No-Go）
- Go：A~E 區塊至少 90% 完成，且 API 不破壞策略確認
- No-Go：任何核心流程缺測、或旗標/回滾方案未完成
