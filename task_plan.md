# UI/UX 收斂任務計畫

## 原始需求
- 以單人操作情境重整 datalink UI，讓主流程回到：
  1. 設定資料來源。
  2. 將來源位置 / 數值轉換可視化顯示。
  3. 綁定 Tag。
  4. 再綁到本地 Modbus server，也可寫入資料庫。
- 不要再被舊 SmartDashboard / 多頁切換 / 過多 modal 心智模型拖走。
- 優先把 UI 主線做對，再談 rollout、legacy compat 與最終入口切換。

## 目前主線策略
- 採 **混合式過渡**：以 `/datalink/workbench` 作為新主線入口，逐步取代舊 datalink UI。
- 保留既有 domain 模型與 API：`device -> point -> tag -> output`。
- 介面原則維持不變：
  - Step 1~4 要能在桌面版一路接力。
  - Step 2 以格狀可視化為主角。
  - Step 3/4 要降低不必要的額外操作與資訊噪音。
  - 不為了快速可動而犧牲主流程一致性。

## 目前狀態

### 已完成
- `master-detail polish` 主體已完成並驗證：
  - P0 readiness contract
  - P1 Step 1 device studio
  - P2 Step 2 source studio
  - P3 Step 3 guidance / Step 4 continuity
  - P4 regressions / 1920 desktop smoke / route-lock fix
- Phase 2 契約面已完成：
  - `runtime-live-value-phase2`
  - `database-target-phase2`
- 下一輪 round 2 polish 已完成：
  - Step 1：inline editor 已完成，create/edit/clone 不再走置中 modal。
  - Step 2：32/64-bit merge 顯示、工具列常駐、create flow 與 rule delete consistency 已補強。
- Step 3：已補 per-row `unbind`、unbind confirmation、template panel 預設收合，降低重複 ceremony。
- Step 4：已完成 direct surface binding：
  - Local Modbus：先選 tag chip，再直接點 register/slot 綁定或解除綁定。
  - Database：直接點 schema/column surface 綁定或解除綁定。
- `polish-round2-regression` 已完成：
  - targeted tests、lint、TypeScript、build 全數通過
  - code review 結果：`APPROVED`
  - planning files / session plan 已同步

## 當前待辦任務
1. 已完成 OpenSpec umbrella change：
   - `rework-datalink-rule-persistence-and-point-tag-flow`
2. OpenSpec implementation 已完成第一個可驗證 slice：
   - `1.1 ~ 1.3` persisted SourceRule lifecycle
   - `2.1 ~ 2.2` Step 1 connect / probe diagnostics
   - `3.1 ~ 3.3` Step 2 persisted rule + unmanaged/live state
3. 當前主線待辦：
   - `4.2` Step 3 改成 review / verification + exception handling
   - `5.1 ~ 5.3` Output / Database authoritative state model
4. rollout / embedded static / 主入口切換先暫緩
   - 目前主焦點已轉為規則持久化、Point/Tag 模型與 step 行為語意
5. 2026-03-19 使用者新鎖定決策：
   - Step 2 採 `B`：主畫面改成 **Tag-first**，`Point` 只留在 inspector debug 區
   - 起始位址策略：**優先恢復每台設備上次規劃起點，沒有才用協議預設**
   - persisted 狀態：進入 Step 2 即顯示，但要清楚標示為「既有規則 / 既有 Tag」

## 目前最重要的續作順序
1. 先修 Step 2：per-device 起始位址記憶 + protocol fallback baseline + Tag-first surface
2. 再完成 `4.2`，把 Step 3 從手動綁定改成 review-first surface
3. 再進 `5.x` output / database state model
4. Step 1 連線問題最後回頭分離：
   - UI / draft 行為問題
   - 網路拓樸 / backend host 問題
5. 完成 OpenSpec docs 後，再回頭決定 rollout / 入口切換 / embedded static

## 精簡里程碑歸檔

### 2026-03-15：方向定案
- 完成 datalink UI 深度盤查。
- 確認應以 `/datalink/workbench` 承接新主線，而不是繼續疊在 `SmartDashboardPage` 上。
- 設計 spec：
  - `docs/superpowers/specs/2026-03-15-datalink-workbench-design.md`
- commit：`5afef10 新增 datalink workbench 設計規格`

### 2026-03-16：桌面 redesign / quiet desktop / scan-first 收斂
- 完成 workbench shell、source、tag、output 的桌面重整。
- 建立第二輪 desktop redesign 與 scan-first / quiet desktop v2 設計基礎。
- 完成 phase 2 runtime / database target 契約細化與落地。

### 2026-03-17：master-detail polish 與回歸
- 完成 Step 3 guidance 與 `countEligibleSpans` 對齊：
  - `668543b`
  - `0e239b7`
- 完成 Step 3 -> Step 4 continuity：
  - `3827d39`
- 完成 route-lock regression fix：
  - `efe0e0b`
- 完成 planning/progress 文件同步：
  - `20f3c02`

### 2026-03-17：round 2 使用者回饋
- 使用者要求：
  - Step 1 不要 modal。
  - Step 2 修正 merge / create / delete consistency。
  - Step 3 要降低重複綁定感，且綁定後可取消。
  - Step 4 要拿掉大型候選區，改成直接點表面綁定，並加強資料庫頁面。
- 目前已完成 Step 1 / Step 2 / Step 3 / Step 4 / regression。

### 2026-03-19：OpenSpec 4.1 完成
- `SourceRule` 現在會在 create / update / enable / disable / restart restoration 時自動建立或同步 rule-derived `Tag + Mapping`。
- `mapping.Service.Create` 已改成 strict `1 Point : 1 Tag` gate，並支援建立 disabled mapping，供 disabled rule 預先持有關聯。
- `source_rule_links.tag_id / mapping_id` 會回填，legacy rule link 也會在 `SyncDerivedPointState()` 啟動時補齊。
- rule shrink / delete 會清除 auto-managed orphan tags；tag / mapping not-found 都已改成 sentinel error，避免 string matching drift。
- `useSourceRules` 會在 source-rule mutations 後同步 invalidates `tagKeys` / `mappingKeys`，讓 Step 3/4 即時看到自動產生的關聯。

### 2026-03-19：Step 2 設計方向鎖定
- 使用者已批准 `B`：
  - Step 2 主畫面改用 Tag-first 語意
  - Point 不再作為主流程名詞，只保留在 inspector 少量 debug 細節
- 起始位址行為已定案：
  - 同設備優先恢復上次規劃起點
  - 無歷史值時依協議 fallback（例如 Modbus / Fatek / MC3E 各自預設）
- persisted 狀態語意已定案：
  - 直接顯示既有狀態
  - 但必須明確標示「既有規則 / 既有 Tag」，避免與本次新規劃混淆

### 2026-03-19：Step 2 baseline / Tag-first 實作完成
- 已落地 per-device 規劃起點記憶：
  - `WorkbenchSourcePlanningState.plannerStartAddressByDeviceId`
  - 同設備切回 Step 2 時會恢復上次規劃起點
- 已落地 protocol fallback baseline：
  - Modbus=`40001`
  - FATEK=`D0`
  - MC3E=`D0`
- 已修正 address offset lower bound：
  - Modbus 維持從 `1` 起算
  - FATEK / MC3E 允許 `D0`
- 已落地 Tag-first surface：
  - `AddressCanvas` / `AddressLedger` 改成 `tagDisplayName -> tagKey -> point.name -> generic label`
  - source rule 卡片會標 `既有規則 / 草稿規則`
- 已完成 reviewer 收斂：
  - 文案改回中性 `rule/source` 語氣，避免 overclaim
  - second-pass review：`No significant issues found`
- 已完成驗證：
  - `cd frontend && npm run lint`
  - `cd frontend && npx tsc --noEmit`
  - `cd frontend && npm run build`
- 下一步回到 Step 1：
  - 釐清 draft config test flow 是否要補 draft-aware diagnostics
  - 把 UI 行為問題與 backend host / network topology 問題分開處理
