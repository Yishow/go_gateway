# C 方案功能對齊落差矩陣（Phase 0）

更新時間：2026-03-02  
盤點範圍：`frontend/src/pages/datalink/*`（舊流程）對照 `frontend/src/pages/gateway/*` + `frontend/src/features/gateway/*`（C 版）  
硬限制：**不修改後端 API 契約**

## 狀態定義
- **已完成**：C 版已有可用實作，且有測試/守門。
- **部分完成**：已有骨架或局部流程，尚未達舊版可用深度。
- **未完成**：C 版尚未覆蓋。

## 功能對照（原 datalink → C 版）

| # | 原 datalink 功能（現況） | C 版狀態 | 補齊方案（可執行） | 主要風險 | 驗收條件 |
|---|---|---|---|---|---|
| 1 | 建立入口導流（`/datalink/devices/new` 依旗標導到新入口）`frontend/src/App.tsx`, `frontend/src/router/gateway.tsx` | 已完成 | 保持 `GatewayCreateEntryRedirect`，僅在前端調整導流，不碰 API。 | 旗標讀取失敗時誤導流 | `cd frontend && npm run test:gateway:unit` 綠燈，手動驗證 flag on/off 導流正確。 |
| 2 | 雙入口選擇頁（Quick/Expert）`GatewayEntryPage.tsx` | 已完成 | 保留入口卡結構，後續只補可觀測性與文案 i18n。 | 卡片文案與實際能力不一致 | `npm run test:gateway:e2e` 案例可從 entry 進入兩條路徑。 |
| 3 | Quick Setup 五步引導（設備→點位→Tag→驗證→Commit）`SmartDashboardPage.tsx` | 部分完成 | 將 `GatewayQuickSetupPage.tsx` 拆成五步 state machine，復用 datalink 查詢 hooks。 | 步驟切換遺失草稿 | E2E 覆蓋 5 步 happy/error path，切換步驟資料不丟失。 |
| 4 | Expert Workbench 全域編輯（路由/插件/原始設定）`GatewayExpertWorkbenchPage.tsx` | 部分完成 | 以現有 route table + raw manifest 為基底，接入 datalink 的裝置/點位/映射資料。 | 高密度畫面效能退化 | 100+ 列資料仍可互動，且提交 payload 符合既有 API。 |
| 5 | Adapter 三向轉換（quick/expert/payload）`gatewayAdapter.ts` | 部分完成 | 擴充欄位映射規則，加入 datalink 常見欄位保留策略（未知欄位不得覆寫）。 | 降級時洗掉進階欄位 | `gatewayAdapter.test.ts` 增加複雜 payload 回填案例並通過。 |
| 6 | `X-UI-Version` 追蹤 header（`dual_quick/dual_expert`）`services/api.ts`, `services/datalink.ts` | 已完成 | 維持 interceptor，並補 E2E 斷言 header。 | 漏標 header 導致灰度指標失真 | Playwright 攔截請求可看到 `X-UI-Version`。 |
| 7 | 設備清單搜尋/狀態篩選 `SmartDashboardWorkflowModal.tsx` | 未完成 | 在 Quick Step1 增加設備搜尋與狀態篩選，重用 `useDevicesQuery`。 | 新手入口被過度複雜化 | Unit 測試篩選邏輯；Manual 驗證 3 種狀態篩選。 |
| 8 | 設備 CRUD + 啟停 + 切換保護 `SmartDashboardPage.tsx` | 未完成 | 抽共用 DeviceOps hook，給 Quick/Expert 共用。 | 切換設備造成髒資料 | E2E：切換設備時有確認提示，確認後狀態一致。 |
| 9 | Source Planner（型別/數量/起始位址）`SmartDashboardWorkspaceContent.tsx` | 未完成 | 在 Quick Step2 引入簡化版 planner，保留 typed occupancy 驗證。 | 位址規劃錯誤導致後續衝突 | Unit：`typedOccupancy` 邊界值；Manual：輸入非法位址顯示錯誤。 |
|10| 自動連續位址配置 `findNearestValidContiguousSpan` | 未完成 | Quick Step2 加「自動配置」按鈕，直接復用策略函式。 | 大範圍掃描延遲 | 單元測試 + 200 筆規劃手動測時可接受。 |
|11| 模板儲存/載入/升級 `sourceTemplateStorage` | 未完成 | Quick/Expert 共用模板抽屜，沿用現有 schema version。 | 模板版本升級不相容 | 單元測試 upgrade path；Manual 驗證舊模板可升級載入。 |
|12| Memory Grid 衝突可視化與選取 `MemoryGrid` | 未完成 | Expert 佈局嵌入 `MemoryGrid`，保留衝突過濾與選取回饋。 | 網格與新佈局交互衝突 | E2E：衝突格顯示、選取、右鍵操作可用。 |
|13| 批次建立點位與命名預覽 `SmartDashboardPanels.tsx` | 未完成 | 在 Expert 增加 Batch 建立面板，沿用 `BatchPointCreator`。 | 批次建立失敗回滾不完整 | Integration：建立 N 筆點位後資料一致。 |
|14| Tag 連結/新建/全域編輯 guardrail `useSmartDashboardTagLinking.ts` | 未完成 | 把 Tag 流程接入 Quick Step3 與 Expert 側欄。 | 全域 Tag 被誤改 | Unit：guardrail；Manual：多映射 Tag 編輯需二次確認。 |
|15| Local Modbus Share 快捷操作 `SmartDashboardModbusPanel.tsx` | 未完成 | Expert 側欄保留綁定/推送/同步按鈕，串 `useSmartDashboardModbusActions`。 | 寫入時機錯誤 | Manual：bind/push/sync 三按鈕全通。 |
|16| Local Modbus 完整工作台（5020）`LocalModbusWorkbenchPage.tsx` | 未完成 | C 版保留 deep-link；先不重寫，直接導向既有頁。 | 入口切換上下文遺失 | 點擊「完整工作台」可帶 `section` 回跳。 |
|17| 匯入/匯出點位、Undo/Redo `ImportDialog/ExportDialog/usePointHistory` | 未完成 | 在 Expert 工具列接回 import/export/undo/redo。 | 匯入大量資料卡頓 | E2E：匯入後可撤銷/重做，匯出檔可重新匯入。 |
|18| 流程驗證（結構 + pipeline）`runStructuralValidation` + `validatePipeline` | 未完成 | Quick Step4/Expert Commit 前統一走 validation gate。 | 驗證規則分叉 | Unit：validation case；Manual：錯誤訊息一致。 |
|19| Commit 生命週期（chunk/retry/rollback/audit）`useSmartDashboardCommitFlow.ts` | 未完成 | 在 Quick Step5 與 Expert 共用 commit hook。 | 重試/回滾語意不一致 | E2E：partial fail→retry→rollback 流程可重現。 |
|20| 工作流 intent modal / legacy notice / 快捷鍵 `SmartDashboardIntentNotices`, `useKeyboardShortcuts` | 未完成 | C 版先補最小 notices + 快捷鍵，再補完整 modal。 | 快捷鍵衝突 | 手動驗證 Ctrl+Enter / Ctrl+Shift+Enter / Alt+R。 |
|21| i18n 完整化（zh-TW/en）`SmartDashboard*` 使用 `t()` | 未完成 | 將 Gateway 三頁 hardcode 文案改為 i18n key。 | 文案回歸與 key 漏補 | `npm run lint && npm run test:gateway:unit`；語系切換無硬編碼殘留。 |

## Phase 0 結論
- C 版已完成「入口、旗標、最小 Quick/Expert 骨架、API header 埋點」。
- 與原 datalink 的核心作業能力仍有明顯落差，重點缺口集中在：**設備管理、規劃/網格、Tag/Mapping、Commit 生命周期**。
- 建議先以 P0 補齊可提交閉環（Step1~Step5 + commit/retry/rollback），再做 P1/P2 擴展。

