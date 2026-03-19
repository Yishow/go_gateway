# Progress

## 當前狀態摘要
- `task_plan.md`、`findings.md`、`progress.md` 已重構為精簡續作版本。
- 原始需求已保留：datalink UI 仍以「來源設定 -> 可視化 -> Tag -> Local Modbus / Database」為唯一主線。
- 目前最新狀態：
  - Step 1 round 2 polish：完成
  - Step 2 round 2 polish：完成
  - Step 3 round 2 polish：完成
  - Step 4 round 2 polish：完成
  - final regression / review：完成

## 最新一輪已完成事項

### round 2 polish：Step 1
- `WorkbenchDeviceStep.tsx`
  - create/edit/clone 改成 inline editor，整合進中央主工作區。
  - 不再使用大型置中 modal 作為主要編輯姿態。
- `DatalinkWorkbenchFoundation.test.tsx`
  - 已調整 Step 1 測試 seam，對齊 inline editor。

### round 2 polish：Step 2
- `SourceCanvasSection.tsx`
  - 移除 `More tools` toggle，改為常駐工具列。
  - `Create rule points` 改為 per-span safety，不再被單一 conflict 全面鎖死。
  - `handleDeleteRule` 會一併處理 orphaned points cleanup。
- `AddressCanvas.tsx`
  - 32/64-bit merge 顯示補強。
- tests
  - source step 與 `sourceCanvasModel` 測試已擴充。

### round 2 polish：Step 3
- `TagBindingStudio.tsx`
  - 已補 per-row `unbind`
  - 已補 unbind confirmation
  - template panel 預設收合，降低規則命名後的重複 ceremony
- tests
  - tag step targeted tests 已通過

### round 2 polish：Step 4
- `LocalModbusBoard.tsx`
  - 改成 compact tag chips + register/slot direct click-to-bind
- `DatabaseTargetBoard.tsx`
  - 改成 schema/column surface direct click-to-bind / unbind
  - 移除舊的 form-first 綁定流程
- tests
  - output step targeted tests 已通過

### 已完成的 master-detail polish / regression 基礎
- commits：
  - `668543b 完成 Step 3 空狀態引導與解鎖語意對齊`
  - `0e239b7 修正 countEligibleSpans 語意對齊 Step 2 readyToCreateCount`
  - `3827d39 接通 Step 3→Step 4 跨步驟 tag 焦點交接`
  - `efe0e0b 修正 workbench 深連結後步驟導航被鎖定的回歸問題`
  - `20f3c02 同步 workbench master-detail polish 進度文件`
- 驗證已做過：
  - workbench suite `134/134` passed
  - `npm run lint`
  - `npx tsc --noEmit`
  - `npm run build`
  - 1920×1080 desktop smoke

## 待續作 / blockers
1. 視結果決定 rollout / 入口切換 / embedded static 交付
2. 若使用者實機驗收後仍有阻力，再開小範圍 polish

## 重要驗證摘要
- Step 1 agent 回報：23/24 tests pass，唯一失敗是 Step 2 mock 缺口造成的 pre-existing failure。
- Step 2 agent回報：57/57 tests pass，ESLint clean，TSC clean。
- Step 3 targeted validation：20/20 tests pass，ESLint clean。
- Step 4 targeted validation：29/29 tests pass，TypeScript / ESLint / build clean。
- round 2 validation pass：
  - 163/163 targeted tests pass
  - `npm run lint`
  - `npx tsc --noEmit`
  - `npm run build`
- round 2 code review：
  - `APPROVED`
  - 確認 Step 4 已真正達成 direct surface binding，不只是縮成 compact selector。
- 四步驟除錯修正驗證：
  - workbench suite `141/141` passed（含修復 7 個 pre-existing shell-ui mock 缺口）
  - Go tests: `tag/...`、`handlers/...`、`device/...` 全部通過
  - `go build ./...`、`npx tsc --noEmit`、`npm run lint` 全部通過

## 2026-03-18 OpenSpec / 選項分析進度
- 已完成 OpenSpec 前置盤查：
  - Step1 connect/probe 與診斷分流
  - Step2 rule drift / address model / persistence 方向
  - Step3/4/Database state drift 與 connector scope
  - rule persistence / point-tag model 兩條 domain audit
- 使用者已明確拍板：
  - connect success 與 probe success 要分開顯示
  - connect 成功但 probe 失敗時，可存設備但不可啟動規則/採集
  - rule disable 只停止採集，不拆 Point/Tag/Output 關聯
  - 建立規則後自動建立 Tag + Mapping，Step3 改成檢查/覆核
  - Database 這輪 scope 鎖 `SQLite + PostgreSQL`
- 目前進到第二輪「各 step 可優化處選項分析」：
  - `step2-optimization-analysis`、`step3-optimization-analysis` 已收斂為 done
  - `step1-optimization-analysis` 補跑完成，推薦 `connect/probe 分流 + connect 成功可存設備但不可啟動規則`
  - `step4-optimization-analysis` 補跑完成，推薦 `target selection 隔離 + 單一狀態來源`
  - `database-optimization-analysis` 補跑完成，推薦 `connector/schema/mapping 分層`
  - `workflow-optimization-synthesis` 已完成
  - 使用者已選擇：
    - Step1：`connect/probe 分流`、`connect 成功可存設備但 probe 失敗禁採集`、`更細錯誤診斷`
    - Step2：`規則持久化/啟停/重啟恢復`、`依設備能力自動切換`、`格子語意講清楚`、`drift 偵測與提示`
    - Step3：`建立規則後自動建立 Tag + Mapping`、`Step3 改成檢查/覆核`、`Point 內部化 / Tag 對外化`
    - Step4：`target selection 狀態分離`、`單一狀態來源`、`綁定/解除流程更明確`、`即時衝突提示與高亮`
    - Database：`Connector / Schema / Mapping 分層`
  - OpenSpec umbrella change 已建立完成：
    - `openspec/changes/rework-datalink-rule-persistence-and-point-tag-flow/`
    - artifacts：`proposal.md`、`design.md`、`specs/**/*.md`、`tasks.md`
  - `openspec status --change rework-datalink-rule-persistence-and-point-tag-flow` 已顯示 `All artifacts complete!`

## 2026-03-19 SourceRule persistence slice
- 已完成 OpenSpec `rework-datalink-rule-persistence-and-point-tag-flow` 的 `1.1 ~ 3.3`：
  - backend 新增 persisted `SourceRule` schema / migrations / memory + SQL repositories
  - 新增 `sourcerule.Service` 與 `/datalink/source-rules` CRUD + enable/disable API
  - runtime bootstrap 會在啟動前同步 persisted rule enable 狀態回 derived points
  - Step 2 `Create rule points` 已改走 source-rule API，不再直接 batch create unmanaged points
  - Step 2 現在會 merge persisted rules 與 local drafts，並把 `unmanaged` legacy points 與 `used` rule-derived points 分開顯示
  - Step 2 inspector 已補 unmanaged notice；runtime live value / restored rule state 會回灌到 grid 與 inspector
- 本輪驗證已通過：
  - `go test ./internal/api/handlers ./internal/api ./internal/datalink ./internal/datalink/sourcerule ./cmd/test_ui -count=1`
  - `cd frontend && npm run test -- --run tests/unit/pages/datalink/workbench-source-step.test.tsx tests/unit/pages/datalink/workbench-tag-step.test.tsx tests/unit/pages/datalink/workbench-runtime-phase.test.tsx tests/unit/pages/datalink/workbench-foundation.test.tsx tests/unit/features/datalink/workbench-locale.test.ts`
  - `cd frontend && npm run lint`
  - `cd frontend && npx tsc --noEmit`
  - `cd frontend && npm run build`
- 下一個主題：
  - `4.1` auto-create / sync Tag + Mapping（strict `1 Point : 1 Tag`）
  - `4.2` Step 3 review / verification flow

## 精簡歷史里程碑

### 2026-03-15
- 完成 datalink UI 盤查、設計 spec 與 implementation planning。
- 決定新主線入口為 `/datalink/workbench`。

### 2026-03-16
- 完成 workbench shell / source / tag / output 基底與 phase 2 契約收尾。
- quiet desktop / scan-first 方向定案。

### 2026-03-17（上半段）
- 完成 master-detail polish、route-lock regression、1920 desktop smoke 與文件同步。

### 2026-03-17（下半段）
- 使用者提出 round 2 polish 回饋。
- Step 1、Step 2、Step 3、Step 4 與 final regression 均已完成。

### 2026-03-18
- 完成四步驟除錯修正：
  - Step 1: 連線測試 15 秒逾時
  - Step 2: sourceCanvasModel merge 特徵衝突偵測
  - Step 3: tagAPI.batchCreate + 後端 BatchCreate 交易包裝
  - Step 4: 驗證已對齊，無需修改
- 修復 7 個 pre-existing shell-ui test mock 缺口
- commit: `bc84416`
