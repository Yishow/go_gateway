## 1. 入口頁定位與導流

- [ ] 1.1 落實 `Preserve existing product routes and shells`，讓新增入口頁後 `/studio`、`/gateway/quick-setup`、`/gateway/expert-workbench` 仍可直接進入且不被 redirect 成其他頁；驗證：更新路由測試，確認三個既有路徑各自渲染原有頁面或原有 route family。
- [ ] 1.2 落實 ``Reuse `/gateway/entry` as the additive setup front door`` 與 `Setup entry page presents additive route choices`，讓 `/gateway/entry` 顯示前置入口頁而不是空殼或隱藏頁；驗證：新增或更新 `GatewayEntryPage` / router 測試，確認入口頁主要區塊與 CTA 存在。
- [ ] 1.3 落實 `Quick-complete CTA composes existing setup flows`，讓 `快速完整設定` CTA 導向既有 `/gateway/quick-setup` 而不是新的重複流程；驗證：更新入口頁互動測試，確認主要 CTA 的導向目標是 `/gateway/quick-setup`。

## 2. 既有功能保護

- [ ] [P] 2.1 落實 `Feature-flag behavior only gates the new entry page, not existing workflows` 與 `Entry page rollout does not gate existing workflows`，讓 feature flag 關閉或載入失敗時既有 `/studio`、quick setup、expert workbench 仍保持可用；驗證：更新 `router/gateway` 測試，覆蓋 enabled、disabled、loading/fallback 情境。
- [ ] [P] 2.2 落實 `Guided workflow`，讓入口頁新增後既有 `/studio` 主流程仍保持可直接開啟，且不需要先進入入口頁寫入前置 state；驗證：更新 `App.tsx` 或 workbench route 測試，確認直接打開 `/studio` 仍能正常進入工作台。
- [ ] 2.3 落實 `Desktop workbench shell`，讓入口頁存在時 `/studio` 仍是主要工作台而不是被降級成相容頁；驗證：更新 workbench shell 測試，確認 `/studio` 的主要區域與既有 shell 結構仍可直接載入。

## 3. 文案與回歸驗證

- [ ] [P] 3.1 落實 `gateway-setup-entry` 的產品文案，讓入口頁同時清楚標示 `快速完整設定`、`進入工作台`、`進階設定` 三條路徑且不暗示功能被移除；驗證：檢查 i18n 文案與入口頁測試，確認三種 CTA 文案都存在。
- [ ] 3.2 落實 `Verification protects against regression of existing routes`，執行入口頁與既有 gateway/workbench 受影響的前端測試，證明新增入口頁沒有造成功能遺失；驗證：執行 `cd frontend && npm run test -- --run` 或至少執行受影響的 Vitest 測試集並確認通過。
