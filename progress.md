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

## 2026-03-19 Rule-derived Tag/Mapping slice
- 已完成 OpenSpec `4.1`，並同步補齊 `6.1` auto-generated tag mapping regression coverage：
  - `SourceRule` create / update / enable / disable / restart restoration 都會自動建立或同步 rule-derived `Tag + Mapping`
  - `source_rule_links.tag_id / mapping_id` 會持久化回填，舊 rule link 也會在 `SyncDerivedPointState()` 補齊
  - `mapping.Service.Create` 現在會做 strict `1 Point : 1 Tag` cardinality gate，並支援 disabled mapping create
  - auto-generated tag 會寫入 rule-managed labels，rule shrink / delete 後會清理 orphan tags，不再留下無主 tag
  - `tag` / `mapping` not-found 都已改用 sentinel error，移除 fragile string matching
  - `useSourceRules` 在 source-rule mutations 後會 invalidate `tagKeys` / `mappingKeys`
- 本輪驗證已通過：
  - `go test ./internal/datalink/tag ./internal/datalink/mapping ./internal/datalink/sourcerule ./internal/api/handlers -count=1`
  - `go build ./cmd/test_ui`
  - `cd frontend && npm run test -- --run tests/unit/hooks/useSourceRules.test.tsx src/pages/datalink/workbench/__tests__/DatalinkWorkbenchSourceStep.test.tsx src/pages/datalink/workbench/__tests__/DatalinkWorkbenchTagStep.test.tsx`
  - `cd frontend && npm run lint`
  - `cd frontend && npx tsc --noEmit`
  - `cd frontend && npm run build`
  - `git diff --check`
- review 狀態：
  - focused code review：`No significant issues found`
- 下一個主題：
  - `4.2` Step 3 改成 review / verification-first + exception handling

## 2026-03-19 Step1 / Step2 bug trace（已完成這輪收斂）
- 已先定位兩個高機率根因：
  - Step 1 測試連線只會測已存檔 `selectedDevice`，不會測 inline editor 裡尚未儲存的 draft config
  - Step 2 起始位址目前硬編 `40001`，沒有依協議切換預設 address baseline
- 另有一個待驗證來源：
  - 使用者在 Step 1 改 draft 後若未 save 就進 Step 2，畫面仍會顯示舊 selectedDevice 的 persisted rule/point 狀態，看起來像「還沒規劃前就有被規劃的點位」
- 新增實機網路驗證：
  - `ping 192.168.31.62` 成功
  - 直接 TCP connect `192.168.31.62:502` 在目前執行環境同樣失敗，錯誤為 `Errno 65 No route to host`
  - 代表當前 Step 1 顯示的 connect failure 至少在這組 IP/port 上與真實網路 reachability 一致，非單純 UI 假錯
- Step 1 已完成這輪修正：
  - 新增 `POST /datalink/devices/test-draft`，允許用 unsaved `protocol + connection_config` 做一次性 connect/probe 測試
  - `WorkbenchDeviceStep` inline editor 新增 `測試目前草稿設定`
  - editor 內會顯示 draft connect / probe 結果與 activation blocked 提示
  - inspector / editor 都新增 backend-host hint，清楚說明測試是從 backend 主機發起
- Step 1 驗證已通過：
  - `go test ./internal/datalink/device ./internal/api/handlers ./internal/api -count=1`
  - `go build ./cmd/test_ui`
  - `cd frontend && npm run test -- --run tests/unit/pages/datalink/workbench-foundation.test.tsx`
  - `cd frontend && npm run lint`
  - `cd frontend && npx tsc --noEmit`
  - `cd frontend && npm run build`

## 2026-03-19 Step 2 Tag-first / baseline 修正
- 已完成 Step 2 這輪落地與 reviewer 收斂：
  - `WorkbenchProvider` 新增 per-device `plannerStartAddressByDeviceId`
  - `SourceCanvasSection` 會優先恢復每台設備最後規劃起點，否則 fallback 到協議預設
  - `addressParser.ts` 新增 `getDefaultPlannerStartAddress()`（Modbus=`40001`、FATEK=`D0`、MC3E=`D0`）
  - `addressParser.offset()` 已依協議調整 lower bound，修正 FATEK / MC3E 的 `D0` 規劃被誤偏成 `D1`
  - `AddressCanvas` / `AddressLedger` 主標題改成 tag display 優先，但仍保留 `point.name` fallback，避免未綁 Tag 的既有點位失去辨識度
  - source rule 卡片新增 `既有規則 / 草稿規則` badge
  - 主畫面中 summary / action / canvas / ledger 文案已改成較中性的 rule/source 語氣，避免 UI overclaim，但整體資訊排序仍維持 Tag-first
- 新增測試：
  - `WorkbenchProvider.test.tsx`：切設備後仍保留 per-device planner 起點記憶
  - `DatalinkWorkbenchSourceStep.test.tsx`：device-1 / device-2 切換時會恢復個別起始位址並套協議預設
  - `frontend/tests/unit/pages/datalink/workbench-provider.test.tsx`：新增正式 wrapper 入口
- 驗證已通過：
  - `cd frontend && npm run test -- --run tests/unit/pages/datalink/workbench-provider.test.tsx tests/unit/pages/datalink/workbench-source-step.test.tsx`
  - `cd frontend && npm run lint`
  - `cd frontend && npx tsc --noEmit`
  - `cd frontend && npm run build`
  - focused code review second pass：`No significant issues found`
- 下一個主題：
  - 視使用者決定回到 `4.2` Step 3 review / verification-first
  - 或繼續做下一輪實機回饋收斂

## 2026-03-19 Step 4 Local Modbus indexing / viewport 修正
- 已完成這輪 Step 4 實機回饋修正：
  - Local Modbus surface 改成 **1-based HR 顯示**，不再出現 `HR0`
  - 手動輸入 register 改為接受 `1..65536`，送 backend 時再轉回 internal 0-based register
  - tag chips、canvas、dry-run、inspector trace 與 conflict 文案已統一使用 1-based 顯示
  - canvas 不再固定只顯示 64 格；改成以目前選取 / 輸入 register 為 anchor 的 bounded viewport，讓 `HR200` 這類高位址可見，同時避免一次 render 幾萬格
  - 低位址 `HR1` 現在可直接從 canvas 點擊解除綁定
  - 補齊缺漏的 `slotBound / slotUnbound` i18n key
- 新增 / 更新測試：
  - `DatalinkWorkbenchOutputStep.test.tsx`：改以 1-based register expectation 驗證
  - 新增 `HR200` 可見 regression
  - 新增 `HR1` 可直接解除綁定 regression
- 驗證已通過：
  - `cd frontend && npm run test -- --run tests/unit/pages/datalink/workbench-output-step.test.tsx`
  - `cd frontend && npm run lint`
  - `cd frontend && npx tsc --noEmit`
  - `cd frontend && npm run build`
  - focused code review：已補上 bounded viewport，避免高位址造成 DOM freeze

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
