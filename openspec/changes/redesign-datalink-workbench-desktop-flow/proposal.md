## Why

目前 `/datalink/workbench` 的方向是對的，但 desktop 資訊架構已經偏離先前核准的工作台骨架：1920×1080 下主工作區被過度壓縮、Step 1 在主欄內再塞 inline inspector、Step 2 缺少真正的詳細規畫與值/轉換呈現、Step 3/4 的資訊密度不足以支撐真實作業。這讓 workbench 雖然已經能走通基本流程，卻還沒有達到可作為 datalink 主工作台的可用性與可維運性。

現在需要用一個正式的 OpenSpec change，把第二輪 desktop redesign 收斂成可規劃、可拆任務、可驗證的變更契約，避免後續又回到局部補丁式修 UI。

## What Changes

1. 將 `/datalink/workbench` 的 desktop shell 收斂回核准骨架：`StepRail + ContextBar + PrimaryWorkArea + InspectorPanel + BottomSummaryBar`，移除目前偏離 spec 的常駐 page-level `ActionDock` 角色。
2. 以全新 workbench-first UI 重做 Step 1~4，不沿用舊 SmartDashboard / Local Modbus presentational components，只重用健康的 hooks、services、types 與純 helper。
3. 重新定義 Step 1 `DeviceWorkspace`，補齊設備能力摘要、clone flow、連線測試歷史，讓後續 Source 規畫有正確上下文。
4. 重新定義 Step 2 `AddressCanvasWorkspace`，明確定義 rule、16-bit lattice / merged span 規則、多組 rule 圖層管理、coverage/gap overview、Plan/Live/Link 三種資訊視圖、value format/freeze/snapshot/template 等能力。
5. 重新定義 Step 3 `TagBindingBoard`，把 source→tag 的命名規則、地址/span、raw/transformed value、bit width/merge、綁定狀態、既有/新建分流、diff preview、結果總表放進高資訊密度主區。
6. 重新定義 Step 4 `OutputWorkspace`，把 Local Modbus 與 Database 輸出整合到同一個工作區，補齊 register map/schema snapshot、required/missing、output readiness、filter/search、auto-map、dry-run、health summary 與 preview/result surfaces。
7. 更新 datalink UI 規格，使 workbench desktop 路徑成為正式可規劃能力，而不是只存在於非正式設計文件。

## Capabilities

### New Capabilities
- `datalink-workbench-desktop`: 定義 `/datalink/workbench` 的 desktop shell、Step 1~4 工作區、共享 inspector、bottom summary 與跨步驟 readiness/selection 行為。

### Modified Capabilities
- `datalink-ui`: 更新 guided workflow、workspace visualization、responsive/accessibility、source template 等需求，使其與新的 workbench desktop shell 一致。
- `local-modbus-memory-workbench`: 將 Local Modbus 的主要作業面重新定位為 workbench Output step 的一部分，並保留必要的能力契約與導流/相容策略。

## Impact

- Affected specs:
  - `openspec/specs/datalink-ui/spec.md`
  - `openspec/specs/local-modbus-memory-workbench/spec.md`
  - new `openspec/changes/redesign-datalink-workbench-desktop-flow/specs/datalink-workbench-desktop/spec.md`
- Affected design docs:
  - `docs/superpowers/specs/2026-03-16-datalink-workbench-desktop-redesign.md`
- Affected frontend areas:
  - `frontend/src/pages/datalink/workbench/*`
  - related datalink hooks/services/types/i18n/tests
- Affected UX contracts:
  - desktop 1920×1080 layout stability
  - source planning and live visualization fidelity
  - tag management information density
  - integrated Local Modbus / Database output workflow
