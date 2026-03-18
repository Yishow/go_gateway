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
1. 規劃後續 rollout
   - 是否把 `/datalink/workbench` 升格成 datalink 主入口。
   - 何時處理 embedded static / build 交付。
2. 規劃是否還需要下一輪純 UX polish
   - 以實機驗收為準，不再先預設新一輪大改。

## 目前最重要的續作順序
1. 視實機驗收結果決定 rollout / 入口切換
2. 若要交付桌面版 build，再處理 embedded static
3. 若使用者仍有局部操作阻力，再開新一輪小範圍 polish

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
