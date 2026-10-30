# 記錄首頁與共用呈現：實作驗證紀錄

日期：2026-10-30。範圍：`simplify-recording-home-guidance` 任務 1.1–2.8 與 3.1 的相關驗證。
這是局部實作證據，不是七個 change 的整體完成、上線或 Browser 驗收報告。

## 實際入口與不變的安全邊界

- 正式既有 `/studio/v2` 預設仍是 advanced；`?view=recordings` 可選唯讀首頁。
  locator 與 runtime readiness return 優先；無效意圖不得降級成另一個可操作頁面。
  query/hash、同 workspace 的篩選／頁碼／返回焦點沿原 presentation owner 保留。
- 首頁讀既有 workspace、sources、groups、lifecycle、delivery、runtime evidence DTO。
  DB setup revision、Share revision、source generation、group desired/requested/effective
  各自有權威，不能互相替代。已觀測的合法空 DB revision 不等於缺失證據。
- 列表 25/50 筆、設備摘要最多 4 筆；狀態讀取固定當次頁面 demand IDs，最多 6 併發。
  已恢復項目退出篩選不可連鎖補入下一批而全表掃描。取消佇列不代表中止已送 GET。
- 舊 SQL receipt 只證明該次寫入；不推導目前健康、owner.none 或新版已生效。
  未完整觀察標未知；過去異常不能說成目前異常，也不能宣稱未觀察的異常數是零。
- ThemeProvider 仍是唯一持久偏好 owner。wire 值沿既有 auto/light/dark。
  theme-only 只 updateKey(theme)，不保存未儲存的 Share／scheduler 草稿。
  明確 Reset 是 draft-only，Cancel 零動作；舊 Save 不能把後來 Reset 當成已保存。
- 共用 Radix 對話框沿既有 controller/backend guard。UI 焦點與確認意圖 fence
  不取代命令權限；舊 revision、workspace A→B→A 回覆不得重開同意或搶焦點。
- `/__dev/recording-style` 僅 DEV lazy route，主選單沒有入口。獨立 iframe document
  建 React root、實際 home/remove-dialog/shared controls，fixture 明示模擬。
  專用 i18n instance 不裝 detector/cache，不 import 正式 i18n config。
  preview bootstrap／互動無 API、storage 或命令寫入；外層 App 原有 readonly
  settings hydration 不在「preview 零網路」宣稱內。

## 版本與驗證方法

產品基準包含正式 checkpoint `b2e7ad9161d537cdc9b57ec68028f5422cf7d070`。
這次只新增本報告與真實 theme caller 測試，沒有更動產品程式。
執行環境版本：Node 24.19.0、Vitest 5.0.5、TypeScript 7.0.2、Vite 8.3.5、ESLint 10.13.0。
命令由 repo 的 `frontend/` 執行，diff／line gate 由 repo 根目錄執行。
原始執行 logs 與 SHA manifests 留在工作證據區，不把私人路徑或大型 logs 寫入 repo。

### 本次整合選測：37 files / 484 tests PASS

選擇理由：後續列表、證據、焦點和 shared styling 改動會影響整個 readonly home；
同時覆蓋 locator hydration／接收者和 runtime/settings 返回路徑，避免只測純函式。
新 `theme-real-callers.test.tsx` 掛真 ThemeProvider、DatalinkLayout、ThemeToggle；
切換保留原 Outlet DOM 與未儲存輸入，只有 theme 寫入，不用假 useTheme consumer 代替。

```sh
node_modules/.bin/vitest run \
  tests/unit/workbench-v2/recording-home-*.test.* \
  tests/unit/workbench-v2/studio-locator-*.test.* \
  tests/unit/services/studio-recording-view.test.tsx \
  tests/unit/services/studio-entity-locator.test.ts \
  tests/unit/services/studio-v2-locator-entry.test.ts \
  tests/unit/hooks/use-studio-locator-entry.test.tsx \
  tests/unit/runtime-dashboard/runtime-settings-return-history.test.tsx \
  tests/unit/runtime-dashboard/runtime-settings-links.test.tsx \
  tests/unit/runtime-logs/runtime-log-settings-navigation.test.tsx \
  tests/unit/contexts/theme-real-callers.test.tsx
```

該次出現 5 筆 JSDOM「navigation to another Document 未實作」診斷。
測試斷言通過不等於真瀏覽器跨 document 導航已驗；此限制保留，沒有 suppress 成成功證據。

### 引用已終值且 hash 未變的相關 gates

以下是各自實際執行的批次，會互相重疊，**不可加總成獨立測試總數**。
核對相關 task manifests 與現在 bytes 一致後引用，沒有只為增加統計重跑。

| 範圍 | 終值 | 覆蓋 |
| --- | --- | --- |
| Theme/settings 最後相關批次 | 13 files / 89 tests PASS | 真 SettingsPage、theme-only、raw/OS/context/queue、reset/save race、locale |
| Reset race 最後補測 | 3 files / 21 tests PASS | 確認／取消、in-flight Save、後續明確 choice |
| Modal/return 最後修正批次 | 5 files / 91 tests PASS | closed details、真 GroupLifecycleBar、dialog、summary、home a11y |
| Shared modal/affected callers 整合批次 | 18 files / 295 tests PASS | 原2.7相關 caller gates；後續新增修正以以上批次補驗 |
| Preview 最後相關批次 | 6 files / 22 tests PASS | 真 bootstrap、互動、入口、AppRoutes/legacy regression、真 production artifact scan |

代表性測試名稱（原測試檔可核對）：

- `recording-home-management.test.tsx`：`observed-only filter cannot turn unobserved records into a health claim`。
- `recording-home-runtime-evidence.test.ts`：`preserves saved C / requested B / currently effective A and pending B boundary`。
- `recording-modal-return.test.tsx`：`an issued remove reply cannot settle a new A→B→A owner or refocus the reused trigger`。
- 同檔：`closing Advanced while remove review is loading returns to the visible Step4 heading on cancel`。
- `theme-settings.test.tsx`：`a full Save started before reset cannot silently persist the newer reset draft`。
- `recording-style-bootstrap.test.tsx`：`keeps the real iframe bootstrap and all preview interactions free of persistent or network writes`。

### Lint、型別、build 與內容核對

- 本次 scoped ESLint：81 個 TS/TSX 路徑 PASS；命令與選集見附錄。
- `node_modules/.bin/tsc --noEmit`：PASS，涵蓋新增真 caller 測試。
- 最後產品 bytes 的真 Vite production build：PASS，沿2.8終值引用。
  `recording-style-production.test.ts` 實際啟動 Vite CLI `build --mode production`
  並明確給 `NODE_ENV=production`，使用可丟棄輸出資料夾。
  掃輸出檔名和 JS/CSS/HTML/map：preview route、frame/entry、fixture marker、
  preview i18n instance 名稱及 chunks 均不存在。不是 mock ENV 後只測元件。
- `git diff --check`、`bash scripts/check_file_lines.sh`：PASS。
- 驗證前後 89 個本 change 既有檔案 SHA-256 相符；新檔另列交接 manifest。
  其他 agent 的 Go WIP 不作為本次前端驗證成功或失敗的推論。

## 獨立 review 找到並修過的問題

1. 同值 mock 遮蔽 DB setup 與 Share revision 不同權威；改真 DTO fixture，
   並保留 fresh workspace 已觀測空 revision。來源 freshness 不借 Share 造證據。
2. observed-only 異常恢復後自動補頁造成連鎖全掃；固定單次 demand IDs。
3. keyword 將 row 移出頁面後遺失 source attention；保留有 scope 的已觀察 cache。
   舊 aggregate fault 不能繞過 current effective/revision/window 判斷。
4. summary 可見文字與 accessible name 不一致；兩語言完整含觀測時間且每筆唯一。
5. explicit Reset 被舊 theme owner 忽略；沿原 choice/save fence 協調 draft-only intent。
6. modal 外恢復按鈕在窄畫面仍40px；受影響入口採既有 mobile44 token。
7. remove consent 跨 group/revision/workspace 被重用；captured intent/owner fence
   與晚回覆保護，不改 backend guard。closed native details 內容不能返回焦點，
   第一直接 summary 仍可見，外層同 workspace heading 可安全 fallback。
8. preview 正式 i18n detector 在 bootstrap 先寫 storage，原 mount 後 spy 漏掉；
   改 disposable instance，spy 在真 entry dynamic import 前安裝，正式 RED→GREEN。

## 仍未完成／不能由本報告宣稱通過

- **Browser 全部 NOT RUN**：沒有真 viewport、Tab/Shift+Tab、portal/frame focus、
  light/dark 視覺對比、手機／短螢幕幾何、真 OS 行為或實測截圖。
  JSDOM/CSS contract 不能代替這些驗收；環境限制與使用者授權不豁免全案 Browser gate。
- 編輯器完整 route 掛載、lifecycle 簡化流程整合、adoption 預設切換仍由各自後續任務處理；
  本報告不把已發布 foundation／唯讀導航說成端到端可用。
- 本次沒有跑整個產品前後端全套、真 PLC/protocol server 或 SQL E2E；
  七案整體完成時仍要最終全套與真場景驗收。沒有新增或宣稱不存在的 CI 成功。
- 3.2 獨立整體 review 已完成，沒有已確認而尚未修復的缺陷；累計最終 91 個路徑 hash 相符，
  獨立窄整合 6 files / 28 tests PASS。既有逐批結果未加總或冒充 Browser 驗收。
- 本 change 暫不 archive；後續 editor／lifecycle／adoption 與七案整體驗收仍未完成。

## 附錄：本次 scoped ESLint 選集

```sh
node_modules/.bin/eslint \
  src/App.tsx \
  src/components/ThemeToggle.tsx \
  src/contexts/ThemeContext.tsx \
  src/contexts/themePreference.ts \
  src/dev/RecordingStyleFrame.tsx \
  src/dev/RecordingStylePreview.tsx \
  src/dev/previewI18n.ts \
  src/dev/recording-style-entry.tsx \
  src/features/datalink/workbench-v2/recording-home/RecordingHome.tsx \
  src/features/datalink/workbench-v2/recording-home/RecordingHomeCollection.tsx \
  src/features/datalink/workbench-v2/recording-home/RecordingHomeResultFacts.tsx \
  src/features/datalink/workbench-v2/recording-home/RecordingHomeRoute.tsx \
  src/features/datalink/workbench-v2/recording-home/recordingHomeCachedReads.ts \
  src/features/datalink/workbench-v2/recording-home/recordingHomeCollections.ts \
  src/features/datalink/workbench-v2/recording-home/recordingHomeGuidance.ts \
  src/features/datalink/workbench-v2/recording-home/recordingHomeOwnerReads.ts \
  src/features/datalink/workbench-v2/recording-home/recordingHomePresentation.ts \
  src/features/datalink/workbench-v2/recording-home/recordingHomeProjection.ts \
  src/features/datalink/workbench-v2/recording-home/recordingHomeReadQueue.ts \
  src/features/datalink/workbench-v2/recording-home/recordingHomeResultGuidance.ts \
  src/features/datalink/workbench-v2/recording-home/recordingHomeResultSummary.ts \
  src/features/datalink/workbench-v2/recording-home/recordingHomeRuntimeEvidence.ts \
  src/features/datalink/workbench-v2/recording-home/recordingHomeSources.ts \
  src/features/datalink/workbench-v2/recording-home/useRecordingHomeEvidence.ts \
  src/features/datalink/workbench-v2/recording-home/useRecordingHomeFacts.ts \
  src/features/datalink/workbench-v2/settings/SettingsPage.tsx \
  src/features/datalink/workbench-v2/settings/UiSection.tsx \
  src/features/datalink/workbench-v2/settings/useSettingsOperations.ts \
  src/features/datalink/workbench-v2/state/basicRecording.ts \
  src/features/datalink/workbench-v2/state/writeGroup/candidates.ts \
  src/features/datalink/workbench-v2/steps/step4/GuidedRecordingPreparation.tsx \
  src/features/datalink/workbench-v2/steps/step4/RecordingRuleSummary.tsx \
  src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx \
  src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupLifecycleBar.tsx \
  src/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/RecordingChangeDialog.tsx \
  src/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/SavedRecordingRemoveDialog.tsx \
  src/features/datalink/workbench-v2/tokens.ts \
  src/hooks/datalink/studioV2RecordingReadScope.ts \
  src/hooks/datalink/useStudioV2Workspace.ts \
  src/hooks/datalink/useStudioV2WriteGroupLifecycle.ts \
  src/hooks/datalink/useStudioV2WriteGroups.ts \
  src/i18n/config.ts \
  src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx \
  src/services/studioRecordingView.ts \
  src/services/studioV2WorkspaceWriteGroups.ts \
  src/services/studioV2WriteGroupLifecycle.ts \
  src/styles/tokens.ts \
  src/utils/returnFocus.ts \
  tests/unit/contexts/theme-i18n.test.ts \
  tests/unit/contexts/theme-preference.test.tsx \
  tests/unit/contexts/theme-settings.test.tsx \
  tests/unit/dev/recording-style-bootstrap.test.tsx \
  tests/unit/dev/recording-style-entry.test.tsx \
  tests/unit/dev/recording-style-preview.test.tsx \
  tests/unit/dev/recording-style-production.test.ts \
  tests/unit/services/studio-recording-view.test.tsx \
  tests/unit/utils/return-focus.test.ts \
  tests/unit/workbench-v2/quiet-mineral-tokens.test.ts \
  tests/unit/workbench-v2/recording-change-dialog.test.tsx \
  tests/unit/workbench-v2/recording-guided-preparation-host.test.tsx \
  tests/unit/workbench-v2/recording-guided-summary.test.tsx \
  tests/unit/workbench-v2/recording-home-accessibility.test.tsx \
  tests/unit/workbench-v2/recording-home-cached-reads.test.ts \
  tests/unit/workbench-v2/recording-home-collections.test.ts \
  tests/unit/workbench-v2/recording-home-facts.test.tsx \
  tests/unit/workbench-v2/recording-home-hooks.test.tsx \
  tests/unit/workbench-v2/recording-home-management.test.tsx \
  tests/unit/workbench-v2/recording-home-owner-reads.test.ts \
  tests/unit/workbench-v2/recording-home-projection.test.ts \
  tests/unit/workbench-v2/recording-home-read-queue.test.ts \
  tests/unit/workbench-v2/recording-home-read-signals.test.ts \
  tests/unit/workbench-v2/recording-home-result-guidance.test.ts \
  tests/unit/workbench-v2/recording-home-result-summary.test.ts \
  tests/unit/workbench-v2/recording-home-route.test.tsx \
  tests/unit/workbench-v2/recording-home-runtime-evidence.test.ts \
  tests/unit/workbench-v2/recording-home-shell.test.tsx \
  tests/unit/workbench-v2/recording-home-sources.test.ts \
  tests/unit/workbench-v2/recording-home-start-inspection.test.ts \
  tests/unit/workbench-v2/recording-modal-return.test.tsx \
  tests/unit/workbench-v2/semantic-modal-contract.test.tsx \
  tests/unit/contexts/theme-real-callers.test.tsx
```
